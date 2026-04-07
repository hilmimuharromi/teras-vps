package controllers

import (
	"teras-vps/backend/models"
	"teras-vps/backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

// GetProfile returns user profile
func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Return user profile
	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": fiber.Map{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"phone":      user.Phone,
				"role":       user.Role,
				"is_active":  user.IsActive,
				"created_at": user.CreatedAt,
				"updated_at": user.UpdatedAt,
			},
		},
	})
}

// UpdateProfile updates user profile
func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Parse input
	var input utils.UpdateProfileInput
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

	// Update user
	updates := map[string]interface{}{}
	if input.Phone != "" {
		updates["phone"] = input.Phone
	}

	if len(updates) > 0 {
		if err := c.db.Model(user).Updates(updates).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "DB_ERROR",
					"message": "Failed to update profile",
				},
			})
		}
	}

	// Get updated user
	c.db.First(user, user.ID)

	// Return success response
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
		"data": fiber.Map{
			"user": fiber.Map{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"phone":      user.Phone,
				"role":       user.Role,
				"is_active":  user.IsActive,
				"updated_at": user.UpdatedAt,
			},
		},
	})
}

// ChangePassword changes user password
func (c *UserController) ChangePassword(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Parse input
	var input utils.ChangePasswordInput
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

	// Verify old password
	if !utils.ComparePassword(input.OldPassword, user.PasswordHash) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INVALID_PASSWORD",
				"message": "Current password is incorrect",
			},
		})
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "HASH_ERROR",
				"message": "Failed to hash password",
			},
		})
	}

	// Update password
	if err := c.db.Model(user).Update("password_hash", hashedPassword).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to update password",
			},
		})
	}

	// Return success response
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Password changed successfully",
	})
}

// ListSSHKeys returns all SSH keys for the user
func (c *UserController) ListSSHKeys(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get SSH keys
	var sshKeys []models.SSHKey
	if err := c.db.Where("user_id = ?", user.ID).Find(&sshKeys).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to fetch SSH keys",
			},
		})
	}

	// Return SSH keys
	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"ssh_keys": sshKeys,
		},
	})
}

// AddSSHKey adds a new SSH key
func (c *UserController) AddSSHKey(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Parse input
	var input struct {
		Name      string `json:"name" validate:"required,max=100"`
		PublicKey string `json:"public_key" validate:"required"`
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

	// Check if SSH key already exists
	var existingKey models.SSHKey
	if err := c.db.Where("user_id = ? AND public_key = ?", user.ID, input.PublicKey).First(&existingKey).Error; err == nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "KEY_EXISTS",
				"message": "SSH key already exists",
			},
		})
	}

	// Create SSH key
	sshKey := models.SSHKey{
		UserID:    user.ID,
		Name:      input.Name,
		PublicKey: input.PublicKey,
	}

	if err := c.db.Create(&sshKey).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to add SSH key",
			},
		})
	}

	// Return success response
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "SSH key added successfully",
		"data": fiber.Map{
			"ssh_key": sshKey,
		},
	})
}

// DeleteSSHKey deletes an SSH key
func (c *UserController) DeleteSSHKey(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get SSH key ID
	keyID := ctx.Params("id")

	// Find SSH key
	var sshKey models.SSHKey
	if err := c.db.Where("id = ? AND user_id = ?", keyID, user.ID).First(&sshKey).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "SSH key not found",
			},
		})
	}

	// Delete SSH key
	if err := c.db.Delete(&sshKey).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to delete SSH key",
			},
		})
	}

	// Return success response
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "SSH key deleted successfully",
	})
}
