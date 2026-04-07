package models

import (
	"time"
)

// SSHKey represents an SSH public key
type SSHKey struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Name        string    `gorm:"size:100" json:"name"`
	PublicKey   string    `gorm:"type:text;not null" json:"public_key"`
	Fingerprint string    `gorm:"size:255" json:"fingerprint"`
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for SSHKey model
func (SSHKey) TableName() string {
	return "ssh_keys"
}
