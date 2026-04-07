package models

import (
	"time"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id,omitempty"` // Nullable for system actions
	Action    string    `gorm:"size:50;not null" json:"action"`
	Details   string    `gorm:"type:text" json:"details,omitempty"`
	IPAddress string    `gorm:"type:inet" json:"ip_address,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}
