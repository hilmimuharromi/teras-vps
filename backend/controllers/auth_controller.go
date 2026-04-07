package controllers

import (
	"teras-vps/backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

// Register handles user registration
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Registration endpoint - TODO: Implement registration logic",
	})
}

// Login handles user login
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Login endpoint - TODO: Implement login logic",
	})
}

// Logout handles user logout
func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Logout endpoint - TODO: Implement logout logic",
	})
}

// GetMe returns current user info
func (c *AuthController) GetMe(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": models.User{},
		},
	})
}
