package controllers

import (
	"fmt"
	"teras-vps/backend/config"
	"teras-vps/backend/models"
	"teras-vps/backend/proxmox"
	"teras-vps/backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type VMController struct {
	db      *gorm.DB
	cfg     *config.Config
	proxmox *proxmox.Client
}

func NewVMController(db *gorm.DB) *VMController {
	cfg := config.Load()
	return &VMController{
		db:  db,
		cfg: cfg,
	}
}

// getProxmoxClient gets or creates a Proxmox client
func (c *VMController) getProxmoxClient() (*proxmox.Client, error) {
	if c.proxmox == nil {
		client, err := proxmox.NewClient()
		if err != nil {
			return nil, err
		}
		c.proxmox = client
	}
	return c.proxmox, nil
}

// CreateVMRequest represents VM creation request
type CreateVMRequest struct {
	PlanID   uint   `json:"plan_id" validate:"required"`
	Hostname string `json:"hostname" validate:"required,min=3,max=100"`
	Template string `json:"template" validate:"omitempty"`
}

// ListVMs returns all VMs for the authenticated user
func (c *VMController) ListVMs(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get VMs for user
	var vms []models.VM
	if err := c.db.Where("user_id = ?", user.ID).Find(&vms).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to fetch VMs",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"vms": vms,
		},
	})
}

// CreateVM creates a new VM
func (c *VMController) CreateVM(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Check VM limit
	var vmCount int64
	c.db.Model(&models.VM{}).Where("user_id = ?", user.ID).Count(&vmCount)
	if vmCount >= int64(c.cfg.MaxVMsPerUser) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VM_LIMIT_REACHED",
				"message": fmt.Sprintf("Maximum %d VMs allowed per user", c.cfg.MaxVMsPerUser),
			},
		})
	}

	// Parse input
	var input CreateVMRequest
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INVALID_INPUT",
				"message": "Invalid input data",
			},
		})
	}

	// Validate input
	if err := utils.Validate(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
	}

	// Get plan
	var plan models.Plan
	if err := c.db.First(&plan, input.PlanID).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PLAN_NOT_FOUND",
				"message": "Plan not found",
			},
		})
	}

	// Get Proxmox client
	proxmoxClient, err := c.getProxmoxClient()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": "Failed to connect to Proxmox",
			},
		})
	}

	// Get next VM ID
	vmid, err := proxmoxClient.GetNextVMID()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": "Failed to get VM ID",
			},
		})
	}

	// Create VM in Proxmox
	if err := proxmoxClient.CreateVM(vmid, input.Hostname, plan.Cores, plan.Memory, plan.Disk, input.Template); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": fmt.Sprintf("Failed to create VM: %v", err),
			},
		})
	}

	// Set expiry date (1 month from now)
	expiresAt := time.Now().AddDate(0, 1, 0)

	// Create VM record in database
	vm := models.VM{
		UserID:    user.ID,
		VMID:      vmid,
		Name:      input.Hostname,
		Hostname:  input.Hostname,
		Status:    models.VMStatusStopped,
		Cores:     plan.Cores,
		Memory:    plan.Memory,
		Disk:      plan.Disk,
		PlanID:    plan.ID,
		Template:  input.Template,
		ExpiresAt: expiresAt,
	}

	if err := c.db.Create(&vm).Error; err != nil {
		// Try to rollback: delete VM from Proxmox
		proxmoxClient.DeleteVM(vmid)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to create VM record",
			},
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "VM created successfully",
		"data": fiber.Map{
			"vm": vm,
		},
	})
}

// GetVM returns VM details
func (c *VMController) GetVM(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get VM ID
	vmID := ctx.Params("id")

	// Find VM
	var vm models.VM
	if err := c.db.Where("id = ? AND user_id = ?", vmID, user.ID).First(&vm).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "VM not found",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"vm": vm,
		},
	})
}

// UpdateVM updates VM configuration
func (c *VMController) UpdateVM(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get VM ID
	vmID := ctx.Params("id")

	// Find VM
	var vm models.VM
	if err := c.db.Where("id = ? AND user_id = ?", vmID, user.ID).First(&vm).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "VM not found",
			},
		})
	}

	// Parse input
	var input struct {
		Name string `json:"name" validate:"required"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INVALID_INPUT",
				"message": "Invalid input data",
			},
		})
	}

	// Update VM
	vm.Name = input.Name
	if err := c.db.Save(&vm).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to update VM",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "VM updated successfully",
		"data": fiber.Map{
			"vm": vm,
		},
	})
}

// DeleteVM deletes a VM
func (c *VMController) DeleteVM(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get VM ID
	vmID := ctx.Params("id")

	// Find VM
	var vm models.VM
	if err := c.db.Where("id = ? AND user_id = ?", vmID, user.ID).First(&vm).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "VM not found",
			},
		})
	}

	// Delete VM from Proxmox
	proxmoxClient, err := c.getProxmoxClient()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": "Failed to connect to Proxmox",
			},
		})
	}

	if err := proxmoxClient.DeleteVM(vm.VMID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": fmt.Sprintf("Failed to delete VM: %v", err),
			},
		})
	}

	// Delete VM from database
	if err := c.db.Delete(&vm).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to delete VM record",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "VM deleted successfully",
	})
}

// StartVM starts a VM
func (c *VMController) StartVM(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get VM ID
	vmID := ctx.Params("id")

	// Find VM
	var vm models.VM
	if err := c.db.Where("id = ? AND user_id = ?", vmID, user.ID).First(&vm).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "VM not found",
			},
		})
	}

	// Get Proxmox client
	proxmoxClient, err := c.getProxmoxClient()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": "Failed to connect to Proxmox",
			},
		})
	}

	// Start VM in Proxmox
	if err := proxmoxClient.StartVM(vm.VMID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": fmt.Sprintf("Failed to start VM: %v", err),
			},
		})
	}

	// Update VM status in database
	vm.Status = models.VMStatusRunning
	c.db.Save(&vm)

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "VM started successfully",
	})
}

// StopVM stops a VM
func (c *VMController) StopVM(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get VM ID
	vmID := ctx.Params("id")

	// Find VM
	var vm models.VM
	if err := c.db.Where("id = ? AND user_id = ?", vmID, user.ID).First(&vm).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "VM not found",
			},
		})
	}

	// Get Proxmox client
	proxmoxClient, err := c.getProxmoxClient()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": "Failed to connect to Proxmox",
			},
		})
	}

	// Stop VM in Proxmox
	if err := proxmoxClient.StopVM(vm.VMID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": fmt.Sprintf("Failed to stop VM: %v", err),
			},
		})
	}

	// Update VM status in database
	vm.Status = models.VMStatusStopped
	c.db.Save(&vm)

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "VM stopped successfully",
	})
}

// RebootVM reboots a VM
func (c *VMController) RebootVM(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get VM ID
	vmID := ctx.Params("id")

	// Find VM
	var vm models.VM
	if err := c.db.Where("id = ? AND user_id = ?", vmID, user.ID).First(&vm).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "VM not found",
			},
		})
	}

	// Get Proxmox client
	proxmoxClient, err := c.getProxmoxClient()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": "Failed to connect to Proxmox",
			},
		})
	}

	// Reboot VM in Proxmox
	if err := proxmoxClient.RebootVM(vm.VMID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": fmt.Sprintf("Failed to reboot VM: %v", err),
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "VM reboot initiated successfully",
	})
}

// GetVMStats returns VM statistics
func (c *VMController) GetVMStats(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get VM ID
	vmID := ctx.Params("id")

	// Find VM
	var vm models.VM
	if err := c.db.Where("id = ? AND user_id = ?", vmID, user.ID).First(&vm).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "VM not found",
			},
		})
	}

	// Get Proxmox client
	proxmoxClient, err := c.getProxmoxClient()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": "Failed to connect to Proxmox",
			},
		})
	}

	// Get VM stats from Proxmox
	stats, err := proxmoxClient.GetVMStats(vm.VMID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": fmt.Sprintf("Failed to get VM stats: %v", err),
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"stats": stats,
		},
	})
}

// CreateBackup creates a VM backup
func (c *VMController) CreateBackup(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get VM ID
	vmID := ctx.Params("id")

	// Find VM
	var vm models.VM
	if err := c.db.Where("id = ? AND user_id = ?", vmID, user.ID).First(&vm).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "VM not found",
			},
		})
	}

	// Parse input
	var input struct {
		Name string `json:"name" validate:"required"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INVALID_INPUT",
				"message": "Invalid input data",
			},
		})
	}

	// Validate input
	if err := utils.Validate(input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
	}

	// Check backup limit
	var backupCount int64
	c.db.Model(&models.Backup{}).Where("vm_id = ?", vm.ID).Count(&backupCount)
	if backupCount >= int64(c.cfg.MaxBackupsPerVM) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "BACKUP_LIMIT_REACHED",
				"message": fmt.Sprintf("Maximum %d backups allowed per VM", c.cfg.MaxBackupsPerVM),
			},
		})
	}

	// Get Proxmox client
	proxmoxClient, err := c.getProxmoxClient()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": "Failed to connect to Proxmox",
			},
		})
	}

	// Create snapshot in Proxmox
	description := fmt.Sprintf("Backup created at %s", time.Now().Format(time.RFC3339))
	if err := proxmoxClient.CreateSnapshot(vm.VMID, input.Name, description); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PROXMOX_ERROR",
				"message": fmt.Sprintf("Failed to create backup: %v", err),
			},
		})
	}

	// Create backup record in database
	backup := models.Backup{
		VMID:        vm.ID,
		BackupName:  input.Name,
		StoragePath: fmt.Sprintf("/nodes/%s/qemu/%d/snapshot/%s", c.cfg.ProxmoxNode, vm.VMID, input.Name),
	}

	if err := c.db.Create(&backup).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to create backup record",
			},
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Backup created successfully",
		"data": fiber.Map{
			"backup": backup,
		},
	})
}

// ListBackups lists all backups for a VM
func (c *VMController) ListBackups(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get VM ID
	vmID := ctx.Params("id")

	// Find VM
	var vm models.VM
	if err := c.db.Where("id = ? AND user_id = ?", vmID, user.ID).First(&vm).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "VM not found",
			},
		})
	}

	// Get backups
	var backups []models.Backup
	if err := c.db.Where("vm_id = ?", vm.ID).Find(&backups).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to fetch backups",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"backups": backups,
		},
	})
}
