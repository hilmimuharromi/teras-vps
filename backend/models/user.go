package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole represents user role
type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

// User represents a user account
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email        string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Phone        string    `gorm:"size:20" json:"phone,omitempty"`
	Role         UserRole  `gorm:"size:20;default:customer" json:"role"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	VMs      []VM       `gorm:"foreignKey:UserID" json:"vms,omitempty"`
	Invoices []Invoice  `gorm:"foreignKey:UserID" json:"invoices,omitempty"`
	SSHKeys  []SSHKey   `gorm:"foreignKey:UserID" json:"ssh_keys,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = RoleCustomer
	}
	return nil
}
