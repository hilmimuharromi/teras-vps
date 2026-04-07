package models

import (
	"time"
)

// TicketStatus represents support ticket status
type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "open"
	TicketStatusInProgress TicketStatus = "in_progress"
	TicketStatusClosed     TicketStatus = "closed"
)

// TicketPriority represents support ticket priority
type TicketPriority string

const (
	TicketPriorityLow    TicketPriority = "low"
	TicketPriorityMedium TicketPriority = "medium"
	TicketPriorityHigh   TicketPriority = "high"
	TicketPriorityUrgent TicketPriority = "urgent"
)

// Ticket represents a support ticket
type Ticket struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Subject   string         `gorm:"size:255;not null" json:"subject"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Status    TicketStatus   `gorm:"size:20;default:open" json:"status"`
	Priority  TicketPriority `gorm:"size:20;default:medium" json:"priority"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Ticket model
func (Ticket) TableName() string {
	return "tickets"
}

// TicketMessage represents a message in a support ticket
type TicketMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `gorm:"not null;index" json:"ticket_id"`
	UserID    *uint     `gorm:"index" json:"user_id,omitempty"` // Nullable for admin responses
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Ticket Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	User   *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for TicketMessage model
func (TicketMessage) TableName() string {
	return "ticket_messages"
}
