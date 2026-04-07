package controllers

import (
	"teras-vps/backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AdminController struct {
	db *gorm.DB
}

func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{db: db}
}

// StatsResponse represents platform statistics
type StatsResponse struct {
	TotalUsers      int64 `json:"total_users"`
	TotalVMs        int64 `json:"total_vms"`
	RunningVMs      int64 `json:"running_vms"`
	StoppedVMs      int64 `json:"stopped_vms"`
	SuspendedVMs    int64 `json:"suspended_vms"`
	TotalInvoices   int64 `json:"total_invoices"`
	PaidInvoices    int64 `json:"paid_invoices"`
	UnpaidInvoices  int64 `json:"unpaid_invoices"`
	OverdueInvoices int64 `json:"overdue_invoices"`
	TotalRevenue    int64 `json:"total_revenue"`
	MonthlyRevenue  int64 `json:"monthly_revenue"`
}

// GetStats returns platform statistics
func (c *AdminController) GetStats(ctx *fiber.Ctx) error {
	var stats StatsResponse

	// Count users
	c.db.Model(&models.User{}).Count(&stats.TotalUsers)

	// Count VMs
	c.db.Model(&models.VM{}).Count(&stats.TotalVMs)
	c.db.Model(&models.VM{}).Where("status = ?", models.VMStatusRunning).Count(&stats.RunningVMs)
	c.db.Model(&models.VM{}).Where("status = ?", models.VMStatusStopped).Count(&stats.StoppedVMs)
	c.db.Model(&models.VM{}).Where("status = ?", models.VMStatusSuspended).Count(&stats.SuspendedVMs)

	// Count invoices
	c.db.Model(&models.Invoice{}).Count(&stats.TotalInvoices)
	c.db.Model(&models.Invoice{}).Where("status = ?", models.InvoiceStatusPaid).Count(&stats.PaidInvoices)
	c.db.Model(&models.Invoice{}).Where("status = ?", models.InvoiceStatusUnpaid).Count(&stats.UnpaidInvoices)
	c.db.Model(&models.Invoice{}).Where("status = ?", models.InvoiceStatusOverdue).Count(&stats.OverdueInvoices)

	// Calculate revenue
	c.db.Model(&models.Invoice{}).Select("COALESCE(SUM(amount), 0)").Where("status = ?", models.InvoiceStatusPaid).Scan(&stats.TotalRevenue)

	// Calculate monthly revenue (this month)
	c.db.Model(&models.Invoice{}).Select("COALESCE(SUM(amount), 0)").Where("status = ? AND created_at >= NOW() - INTERVAL '1 month'", models.InvoiceStatusPaid).Scan(&stats.MonthlyRevenue)

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"stats": stats,
		},
	})
}

// ListAllUsers returns all users (admin only)
func (c *AdminController) ListAllUsers(ctx *fiber.Ctx) error {
	// Parse pagination
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)
	offset := (page - 1) * limit

	var users []models.User
	var total int64

	// Count total users
	c.db.Model(&models.User{}).Count(&total)

	// Get users with pagination
	if err := c.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to fetch users",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"users": users,
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// ListUsers is a compatibility wrapper for ListAllUsers so routes can call adminController.ListUsers
func (c *AdminController) ListUsers(ctx *fiber.Ctx) error {
	return c.ListAllUsers(ctx)
}

// ListAllVMs returns all VMs (admin only)
func (c *AdminController) ListAllVMs(ctx *fiber.Ctx) error {
	// Parse pagination
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)
	offset := (page - 1) * limit

	var vms []models.VM
	var total int64

	// Count total VMs
	c.db.Model(&models.VM{}).Count(&total)

	// Get VMs with user relation
	if err := c.db.Preload("User").Preload("Plan").Limit(limit).Offset(offset).Order("created_at DESC").Find(&vms).Error; err != nil {
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
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// SuspendUserRequest represents suspend user request
type SuspendUserRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
}

// SuspendUser suspends a user account (admin only)
func (c *AdminController) SuspendUser(ctx *fiber.Ctx) error {
	// Get user ID
	userID := ctx.Params("id")

	// Parse input
	var input SuspendUserRequest
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INVALID_INPUT",
				"message": "Invalid input data",
			},
		})
	}

	// Find user
	var user models.User
	if err := c.db.First(&user, userID).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "User not found",
			},
		})
	}

	// Check if user is admin
	if user.Role == models.UserRoleAdmin {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "CANNOT_SUSPEND_ADMIN",
				"message": "Cannot suspend admin user",
			},
		})
	}

	// Suspend user and all VMs
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Suspend user
	user.IsActive = false
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to suspend user",
			},
		})
	}

	// Stop all user VMs
	if err := tx.Model(&models.VM{}).Where("user_id = ?", user.ID).Update("status", models.VMStatusSuspended).Error; err != nil {
		tx.Rollback()
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to suspend VMs",
			},
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to complete suspension",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "User suspended successfully",
	})
}

// UnsuspendUser unsuspends a user account (admin only)
func (c *AdminController) UnsuspendUser(ctx *fiber.Ctx) error {
	// Get user ID
	userID := ctx.Params("id")

	// Find user
	var user models.User
	if err := c.db.First(&user, userID).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "User not found",
			},
		})
	}

	// Check if user is already active
	if user.IsActive {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "USER_ALREADY_ACTIVE",
				"message": "User is already active",
			},
		})
	}

	// Unsuspend user
	user.IsActive = true
	if err := c.db.Save(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to unsuspend user",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "User unsuspended successfully",
	})
}

// ListAllInvoices returns all invoices (admin only)
func (c *AdminController) ListAllInvoices(ctx *fiber.Ctx) error {
	// Parse pagination
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)
	offset := (page - 1) * limit

	var invoices []models.Invoice
	var total int64

	// Count total invoices
	c.db.Model(&models.Invoice{}).Count(&total)

	// Get invoices with relations
	if err := c.db.Preload("User").Preload("VM").Preload("Payment").Limit(limit).Offset(offset).Order("created_at DESC").Find(&invoices).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to fetch invoices",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"invoices": invoices,
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}
