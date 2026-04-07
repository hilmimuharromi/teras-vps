package controllers

import (
	"teras-vps/backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BillingController struct {
	db *gorm.DB
}

func NewBillingController(db *gorm.DB) *BillingController {
	return &BillingController{db: db}
}

// ListInvoices returns all invoices for the authenticated user
func (c *BillingController) ListInvoices(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get invoices
	var invoices []models.Invoice
	if err := c.db.Where("user_id = ?", user.ID).Order("created_at DESC").Find(&invoices).Error; err != nil {
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
		},
	})
}

// GetInvoice returns invoice details
func (c *BillingController) GetInvoice(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get invoice ID
	invoiceID := ctx.Params("id")

	// Find invoice
	var invoice models.Invoice
	if err := c.db.Where("id = ? AND user_id = ?", invoiceID, user.ID).First(&invoice).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Invoice not found",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"invoice": invoice,
		},
	})
}

// PayInvoice initiates payment for an invoice
func (c *BillingController) PayInvoice(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get invoice ID
	invoiceID := ctx.Params("id")

	// Find invoice
	var invoice models.Invoice
	if err := c.db.Where("id = ? AND user_id = ?", invoiceID, user.ID).First(&invoice).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Invoice not found",
			},
		})
	}

	// Check if invoice is already paid
	if invoice.Status == models.InvoiceStatusPaid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "ALREADY_PAID",
				"message": "Invoice is already paid",
			},
		})
	}

	// TODO: Integrate payment gateway (coming soon)
	// For now, return mock response
	return ctx.JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "PAYMENT_NOT_IMPLEMENTED",
			"message": "Payment gateway not implemented yet. Coming soon!",
		},
	})
}

// ListPlans returns all available plans
func (c *BillingController) ListPlans(ctx *fiber.Ctx) error {
	// Get active plans
	var plans []models.Plan
	if err := c.db.Where("is_active = ?", true).Find(&plans).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to fetch plans",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"plans": plans,
		},
	})
}
