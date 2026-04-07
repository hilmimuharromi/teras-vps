package models

import (
	"time"

	"gorm.io/gorm"
)

// Plan represents a pricing plan
type Plan struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Cores        int       `gorm:"not null" json:"cores"`
	Memory       int       `gorm:"not null" json:"memory"`      // in MB
	Disk         int       `gorm:"not null" json:"disk"`        // in GB
	PriceMonthly int       `gorm:"not null" json:"price_monthly"` // in IDR
	PriceDaily   int       `gorm:"not null" json:"price_daily"`   // in IDR
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	VMs []VM `gorm:"foreignKey:PlanID" json:"vms,omitempty"`
}

// TableName specifies the table name for Plan model
func (Plan) TableName() string {
	return "plans"
}
