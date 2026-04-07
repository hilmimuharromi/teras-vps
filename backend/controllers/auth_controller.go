package controllers

import (
	"teras-vps/backend/models"
	"teras-vps/backend/utils"

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
	// Parse input
	var input utils.RegisterInput
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

	// Check if email already exists
	var existingUser models.User
	if err := c.db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "EMAIL_EXISTS",
				"message": "Email already registered",
			},
		})
	}

	// Check if username already exists
	if err := c.db.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "USERNAME_EXISTS",
				"message": "Username already taken",
			},
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "HASH_ERROR",
				"message": "Failed to hash password",
			},
		})
	}

	// Create user
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Phone:        input.Phone,
		Role:         models.RoleCustomer,
		IsActive:     true,
	}

	if err := c.db.Create(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to create user",
			},
		})
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "TOKEN_ERROR",
				"message": "Failed to generate token",
			},
		})
	}

	// Return success response
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Registration successful",
		"data": fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
		},
	})
}

// Login handles user login
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	// Parse input
	var input utils.LoginInput
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

	// Find user by email
	var user models.User
	if err := c.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
	}

	// Check password
	if !utils.ComparePassword(input.Password, user.PasswordHash) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
	}

	// Check if account is active
	if !user.IsActive {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "ACCOUNT_SUSPENDED",
				"message": "Your account has been suspended",
			},
		})
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "TOKEN_ERROR",
				"message": "Failed to generate token",
			},
		})
	}

	// Return success response
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data": fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
		},
	})
}

// Logout handles user logout
func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	// Note: JWT is stateless, so we don't need to do anything on the server side.
	// The client should remove the token from storage.
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Logout successful",
	})
}

// GetMe returns current user info
func (c *AuthController) GetMe(ctx *fiber.Ctx) error {
	// Get user from context (set by Auth middleware)
	user := ctx.Locals("user").(*models.User)

	// Return user info (without password)
	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": fiber.Map{
				"id":        user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"phone":     user.Phone,
				"role":      user.Role,
				"is_active": user.IsActive,
				"created_at": user.CreatedAt,
				"updated_at": user.UpdatedAt,
			},
		},
	})
}
