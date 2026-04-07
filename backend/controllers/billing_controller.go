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
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    []models.Invoice{},
	})
}

// GetInvoice returns invoice details
func (c *BillingController) GetInvoice(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    models.Invoice{},
	})
}

// PayInvoice initiates payment for an invoice
func (c *BillingController) PayInvoice(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Pay invoice endpoint - TODO: Implement payment logic",
	})
}

// ListPlans returns all available plans
func (c *BillingController) ListPlans(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    []models.Plan{},
	})
}
