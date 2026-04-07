package models

import (
	"time"

	"gorm.io/gorm"
)

// VMStatus represents VM status
type VMStatus string

const (
	VMStatusRunning   VMStatus = "running"
	VMStatusStopped   VMStatus = "stopped"
	VMStatusSuspended VMStatus = "suspended"
)

// VM represents a virtual machine
type VM struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	VMID      int       `gorm:"not null;index" json:"vmid"` // Proxmox VM ID
	Name      string    `gorm:"size:100;not null" json:"name"`
	Hostname  string    `gorm:"size:100" json:"hostname"`
	Status    VMStatus  `gorm:"size:20;default:stopped" json:"status"`
	Cores     int       `gorm:"not null" json:"cores"`
	Memory    int       `gorm:"not null" json:"memory"`      // in MB
	Disk      int       `gorm:"not null" json:"disk"`        // in GB
	PublicIP  string    `gorm:"type:inet" json:"public_ip,omitempty"`
	SSHPort   int       `gorm:"default:22" json:"ssh_port"`
	VNCPort   int       `gorm:"default:5900" json:"vnc_port"`
	Template  string    `gorm:"size:255" json:"template"`
	PlanID    uint      `gorm:"index" json:"plan_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Plan      Plan      `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	Backups   []Backup  `gorm:"foreignKey:VMID" json:"backups,omitempty"`
	Invoices  []Invoice `gorm:"foreignKey:VMID" json:"invoices,omitempty"`
}

// TableName specifies the table name for VM model
func (VM) TableName() string {
	return "vms"
}

// BeforeCreate hook
func (v *VM) BeforeCreate(tx *gorm.DB) error {
	if v.Status == "" {
		v.Status = VMStatusStopped
	}
	return nil
}
