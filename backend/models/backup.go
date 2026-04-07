package models

import (
	"time"
)

// Backup represents a VM backup
type Backup struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VMID        uint      `gorm:"not null;index" json:"vm_id"`
	BackupName  string    `gorm:"size:100;not null" json:"backup_name"`
	StoragePath string    `gorm:"size:255" json:"storage_path"`
	SizeGB      *float64  `json:"size_gb,omitempty"`
	CreatedAt   time.Time `json:"created_at"`

	// Relations
	VM VM `gorm:"foreignKey:VMID" json:"vm,omitempty"`
}

// TableName specifies the table name for Backup model
func (Backup) TableName() string {
	return "vm_backups"
}
