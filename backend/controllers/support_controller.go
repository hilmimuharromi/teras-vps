package controllers

import (
	"teras-vps/backend/models"
	"teras-vps/backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SupportController struct {
	db *gorm.DB
}

func NewSupportController(db *gorm.DB) *SupportController {
	return &SupportController{db: db}
}

// CreateTicketRequest represents ticket creation request
type CreateTicketRequest struct {
	Subject  string `json:"subject" validate:"required,max=255"`
	Message  string `json:"message" validate:"required"`
	Priority string `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
}

// ListTickets returns all tickets for the authenticated user
func (c *SupportController) ListTickets(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get tickets
	var tickets []models.Ticket
	if err := c.db.Where("user_id = ?", user.ID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to fetch tickets",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"tickets": tickets,
		},
	})
}

// CreateTicket creates a new support ticket
func (c *SupportController) CreateTicket(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Parse input
	var input CreateTicketRequest
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

	// Set default priority
	if input.Priority == "" {
		input.Priority = "medium"
	}

	// Create ticket
	ticket := models.Ticket{
		UserID:   user.ID,
		Subject:  input.Subject,
		Message:  input.Message,
		Priority: models.TicketPriority(input.Priority),
		Status:   models.TicketStatusOpen,
	}

	if err := c.db.Create(&ticket).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to create ticket",
			},
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Ticket created successfully",
		"data": fiber.Map{
			"ticket": ticket,
		},
	})
}

// GetTicket returns ticket details with messages
func (c *SupportController) GetTicket(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get ticket ID
	ticketID := ctx.Params("id")

	// Find ticket
	var ticket models.Ticket
	if err := c.db.Preload("User").Where("id = ? AND user_id = ?", ticketID, user.ID).First(&ticket).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Ticket not found",
			},
		})
	}

	// Get messages
	var messages []models.TicketMessage
	if err := c.db.Preload("User").Where("ticket_id = ?", ticket.ID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to fetch messages",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"ticket":   ticket,
			"messages": messages,
		},
	})
}

// AddMessageRequest represents adding a message to a ticket
type AddMessageRequest struct {
	Message string `json:"message" validate:"required"`
}

// AddMessage adds a message to a ticket
func (c *SupportController) AddMessage(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get ticket ID
	ticketID := ctx.Params("id")

	// Verify ticket ownership
	var ticket models.Ticket
	if err := c.db.Where("id = ? AND user_id = ?", ticketID, user.ID).First(&ticket).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Ticket not found",
			},
		})
	}

	// Parse input
	var input AddMessageRequest
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

	// Create message
	message := models.TicketMessage{
		TicketID: ticket.ID,
		UserID:   &user.ID,
		IsAdmin:  false,
		Message:  input.Message,
	}

	if err := c.db.Create(&message).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to add message",
			},
		})
	}

	// Update ticket status if was closed
	if ticket.Status != models.TicketStatusOpen {
		ticket.Status = models.TicketStatusOpen
		c.db.Save(&ticket)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Message added successfully",
		"data": fiber.Map{
			"message": message,
		},
	})
}

// CloseTicket closes a support ticket
func (c *SupportController) CloseTicket(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*models.User)

	// Get ticket ID
	ticketID := ctx.Params("id")

	// Find ticket
	var ticket models.Ticket
	if err := c.db.Where("id = ? AND user_id = ?", ticketID, user.ID).First(&ticket).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Ticket not found",
			},
		})
	}

	// Update ticket status
	ticket.Status = models.TicketStatusClosed
	if err := c.db.Save(&ticket).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "DB_ERROR",
				"message": "Failed to close ticket",
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Ticket closed successfully",
	})
}
