package models

import (
	"time"
)

// TransactionStatus represents transaction status
type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

// Transaction represents a payment transaction
type Transaction struct {
	ID               uint             `gorm:"primaryKey" json:"id"`
	InvoiceID        uint             `gorm:"not null;index" json:"invoice_id"`
	Amount           int              `gorm:"not null" json:"amount"` // in IDR
	Status           TransactionStatus `gorm:"size:20;default:pending" json:"status"`
	PaymentMethod    string           `gorm:"size:50" json:"payment_method,omitempty"`
	GatewayReference string           `gorm:"size:100" json:"gateway_reference,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`

	// Relations
	Invoice Invoice `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
}

// TableName specifies the table name for Transaction model
func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate hook
func (t *Transaction) BeforeCreate() error {
	if t.Status == "" {
		t.Status = TransactionStatusPending
	}
	return nil
}
