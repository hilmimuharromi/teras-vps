package models

import (
	"time"

	"gorm.io/gorm"
)

// InvoiceStatus represents invoice status
type InvoiceStatus string

const (
	InvoiceStatusUnpaid   InvoiceStatus = "unpaid"
	InvoiceStatusPaid     InvoiceStatus = "paid"
	InvoiceStatusOverdue  InvoiceStatus = "overdue"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

// Invoice represents a billing invoice
type Invoice struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"not null;index" json:"user_id"`
	VMID            *uint          `gorm:"index" json:"vm_id,omitempty"` // Nullable for adjustment invoices
	InvoiceNumber   string         `gorm:"uniqueIndex;size:50;not null" json:"invoice_number"`
	Amount          int            `gorm:"not null" json:"amount"`        // in IDR
	Status          InvoiceStatus  `gorm:"size:20;default:unpaid" json:"status"`
	DueDate         time.Time      `gorm:"not null" json:"due_date"`
	PaidDate        *time.Time     `json:"paid_date,omitempty"`
	PaymentMethod   string         `gorm:"size:50" json:"payment_method,omitempty"`
	PaymentReference string         `gorm:"size:100" json:"payment_reference,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	// Relations
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	VM    *VM   `gorm:"foreignKey:VMID" json:"vm,omitempty"`
	Payment *Transaction `gorm:"foreignKey:InvoiceID" json:"payment,omitempty"`
}

// TableName specifies the table name for Invoice model
func (Invoice) TableName() string {
	return "invoices"
}

// BeforeCreate hook
func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.Status == "" {
		i.Status = InvoiceStatusUnpaid
	}
	return nil
}
