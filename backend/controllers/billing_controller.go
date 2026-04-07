package controllers

import (
	"fmt"
	"teras-vps/backend/config"
	"teras-vps/backend/models"
	"teras-vps/backend/services"
	"teras-vps/backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BillingController struct {
	db             *gorm.DB
	billingService *services.BillingService
}

func NewBillingController(db *gorm.DB) *BillingController {
	cfg := config.Load()
	billingService := services.NewBillingService(db, cfg)
	return &BillingController{
		db:             db,
		billingService: billingService,
	}
}

// ListInvoices returns all invoices for the authenticated user
func (c *BillingController) ListInvoices(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get invoices
	var invoices []models.Invoice
	if err := c.db.Preload("Payment").Where("user_id = ?", user.ID).Order("created_at DESC").Find(&invoices).Error; err != nil {
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
	if err := c.db.Preload("Payment").Preload("VM").Preload("VM.Plan").Where("id = ? AND user_id = ?", invoiceID, user.ID).First(&invoice).Error; err != nil {
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

// PayInvoiceRequest represents payment request
type PayInvoiceRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=qris transfer stripe"`
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

	// Parse input
	var input PayInvoiceRequest
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

	// TODO: Integrate payment gateway (Xendit, Midtrans, Stripe, etc.)
	// For now, return placeholder response
	return ctx.JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "PAYMENT_NOT_IMPLEMENTED",
			"message": "Payment gateway integration coming soon! (Week 4 placeholder)",
			"details": fiber.Map{
				"invoice_id":     invoice.ID,
				"invoice_number": invoice.InvoiceNumber,
				"amount":        invoice.Amount,
				"due_date":      invoice.DueDate,
			},
		},
	})
}

// AdminPayInvoice allows admin to manually mark invoice as paid
func (c *BillingController) AdminPayInvoice(ctx *fiber.Ctx) error {
	// Get invoice ID
	invoiceID := ctx.Params("id")

	// Parse input
	var input struct {
		Reference string `json:"reference" validate:"required"`
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

	// Find invoice
	var invoice models.Invoice
	if err := c.db.First(&invoice, invoiceID).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Invoice not found",
			},
		})
	}

	// Process payment
	if err := c.billingService.ProcessPayment(invoice.ID, invoice.Amount, input.Reference); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "PAYMENT_ERROR",
				"message": fmt.Sprintf("Failed to process payment: %v", err),
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Invoice marked as paid successfully",
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
