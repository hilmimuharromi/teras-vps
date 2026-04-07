package controllers

import (
	"teras-vps/backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

// ListSSHKeys returns all SSH keys for the user
func (c *UserController) ListSSHKeys(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    []models.SSHKey{},
	})
}

// AddSSHKey adds a new SSH key
func (c *UserController) AddSSHKey(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Add SSH key endpoint - TODO: Implement SSH key addition",
	})
}

// DeleteSSHKey deletes an SSH key
func (c *UserController) DeleteSSHKey(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Delete SSH key endpoint - TODO: Implement SSH key deletion",
	})
}

// GetProfile returns user profile
func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    models.User{},
	})
}

// UpdateProfile updates user profile
func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Update profile endpoint - TODO: Implement profile update",
	})
}

// ChangePassword changes user password
func (c *UserController) ChangePassword(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Change password endpoint - TODO: Implement password change",
	})
}
