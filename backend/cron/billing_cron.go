package cron

import (
	"log"
	"teras-vps/backend/config"
	"teras-vps/backend/services"
	"time"

	"gorm.io/gorm"
)

// BillingCron handles billing-related cron jobs
type BillingCron struct {
	db             *gorm.DB
	billingService *services.BillingService
}

// NewBillingCron creates a new billing cron
func NewBillingCron(db *gorm.DB, cfg *config.Config) *BillingCron {
	return &BillingCron{
		db:             db,
		billingService: services.NewBillingService(db, cfg),
	}
}

// CheckOverdueInvoices checks for overdue invoices and suspends VMs
func (c *BillingCron) CheckOverdueInvoices() {
	log.Println("Checking overdue invoices...")

	suspended, err := c.billingService.SuspendOverdueVMs()
	if err != nil {
		log.Printf("Error suspending overdue VMs: %v", err)
	} else {
		log.Printf("Suspended %d VMs with overdue invoices", suspended)
	}
}

// CheckTerminatedVMs deletes VMs with very overdue invoices
func (c *BillingCron) CheckTerminatedVMs() {
	log.Println("Checking terminated VMs...")

	deleted, err := c.billingService.DeleteTerminatedVMs()
	if err != nil {
		log.Printf("Error deleting terminated VMs: %v", err)
	} else {
		log.Printf("Deleted %d VMs with very overdue invoices", deleted)
	}
}

// Run starts the billing cron jobs
func (c *BillingCron) Run() {
	// Run every hour
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			// Check overdue invoices
			c.CheckOverdueInvoices()

			// Check terminated VMs (every 24 hours, handled by time check)
			if time.Now().Hour() == 0 {
				c.CheckTerminatedVMs()
			}
		}
	}()

	log.Println("Billing cron jobs started")
}
