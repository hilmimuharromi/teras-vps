package middleware

import (
	"teras-vps/backend/models"
	"teras-vps/backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Auth middleware for protecting routes
func Auth(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "MISSING_TOKEN",
					"message": "Authorization header is required",
				},
			})
		}

		// Extract Bearer token
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Invalid token format. Expected: Bearer <token>",
				},
			})
		}

		tokenString := authHeader[7:]

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired token",
				},
			})
		}

		// Get user from database
		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "USER_NOT_FOUND",
					"message": "User not found",
				},
			})
		}

		// Check if user is active
		if !user.IsActive {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "ACCOUNT_SUSPENDED",
					"message": "Your account has been suspended",
				},
			})
		}

		// Store user in context
		c.Locals("user", &user)
		c.Locals("user_id", user.ID)
		c.Locals("user_email", user.Email)
		c.Locals("user_role", user.Role)

		// Continue to next handler
		return c.Next()
	}
}

// AdminOnly middleware for admin-only routes
func AdminOnly(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from context (set by Auth middleware)
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "UNAUTHORIZED",
					"message": "Authentication required",
				},
			})
		}

		// Check if user is admin
		u := user.(*models.User)
		if u.Role != models.RoleAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "FORBIDDEN",
					"message": "Admin access required",
				},
			})
		}

		// Continue to next handler
		return c.Next()
	}
}
