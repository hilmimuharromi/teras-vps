package services

import (
	"fmt"
	"teras-vps/backend/config"
	"teras-vps/backend/models"
	"time"

	"gorm.io/gorm"
)

// BillingService handles billing operations
type BillingService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewBillingService creates a new billing service
func NewBillingService(db *gorm.DB, cfg *config.Config) *BillingService {
	return &BillingService{
		db:  db,
		cfg: cfg,
	}
}

// GenerateInvoice generates a new invoice
func (s *BillingService) GenerateInvoice(userID uint, vmID *uint, amount int, description string) (*models.Invoice, error) {
	// Generate invoice number
	invoiceNumber := fmt.Sprintf("INV-%s-%04d", time.Now().Format("20060102"), s.getNextInvoiceNumber())

	// Calculate due date (7 days from now)
	dueDate := time.Now().AddDate(0, 0, s.cfg.SuspendDays)

	// Create invoice
	invoice := models.Invoice{
		UserID:        userID,
		VMID:          vmID,
		InvoiceNumber: invoiceNumber,
		Amount:        amount,
		Status:        models.InvoiceStatusUnpaid,
		DueDate:       dueDate,
		PaymentMethod: "qris", // Default payment method
	}

	if err := s.db.Create(&invoice).Error; err != nil {
		return nil, err
	}

	// Create transaction record (pending)
	transaction := models.Transaction{
		InvoiceID:     invoice.ID,
		Amount:        amount,
		Status:        models.TransactionStatusPending,
		PaymentMethod: "qris",
	}

	if err := s.db.Create(&transaction).Error; err != nil {
		return nil, err
	}

	return &invoice, nil
}

// GenerateMonthlyInvoice generates a monthly invoice for a VM
func (s *BillingService) GenerateMonthlyInvoice(vm *models.VM) (*models.Invoice, error) {
	// Get plan
	var plan models.Plan
	if err := s.db.First(&plan, vm.PlanID).Error; err != nil {
		return nil, err
	}

	// Generate invoice with monthly price
	invoice, err := s.GenerateInvoice(vm.UserID, &vm.ID, plan.PriceMonthly, fmt.Sprintf("Monthly invoice for VM: %s", vm.Name))
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

// ProcessPayment processes a payment
func (s *BillingService) ProcessPayment(invoiceID uint, amount int, paymentReference string) error {
	// Start transaction
	tx := s.db.Begin()

	// Get invoice
	var invoice models.Invoice
	if err := tx.First(&invoice, invoiceID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Check if invoice is already paid
	if invoice.Status == models.InvoiceStatusPaid {
		tx.Rollback()
		return fmt.Errorf("invoice already paid")
	}

	// Update invoice
	now := time.Now()
	invoice.Status = models.InvoiceStatusPaid
	invoice.PaidDate = &now
	invoice.PaymentReference = paymentReference

	if err := tx.Save(&invoice).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update transaction
	if err := tx.Model(&models.Transaction{}).Where("invoice_id = ? AND status = ?", invoiceID, models.TransactionStatusPending).Updates(map[string]interface{}{
		"status":            models.TransactionStatusSuccess,
		"gateway_reference": paymentReference,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Extend VM expiry date if VM exists
	if invoice.VMID != nil {
		var vm models.VM
		if err := tx.First(&vm, *invoice.VMID).Error; err == nil {
			// Extend by 1 month
			vm.ExpiresAt = vm.ExpiresAt.AddDate(0, 1, 0)
			// Unsuspend if was suspended
			if vm.Status == models.VMStatusSuspended {
				vm.Status = models.VMStatusRunning
			}
			tx.Save(&vm)
		}
	}

	// Commit transaction
	return tx.Commit().Error
}

// CheckOverdueInvoices checks for overdue invoices
func (s *BillingService) CheckOverdueInvoices() ([]models.Invoice, error) {
	var invoices []models.Invoice

	now := time.Now()
	suspendDate := now.AddDate(0, 0, -s.cfg.SuspendDays)

	if err := s.db.Where("due_date < ? AND status != ?", suspendDate, models.InvoiceStatusPaid).Find(&invoices).Error; err != nil {
		return nil, err
	}

	return invoices, nil
}

// SuspendOverdueVMs suspends VMs with overdue invoices
func (s *BillingService) SuspendOverdueVMs() (int, error) {
	// Get overdue invoices
	invoices, err := s.CheckOverdueInvoices()
	if err != nil {
		return 0, err
	}

	suspendedCount := 0

	for _, invoice := range invoices {
		// Skip if invoice is already overdue
		if invoice.Status == models.InvoiceStatusOverdue {
			continue
		}

		// Update invoice status
		invoice.Status = models.InvoiceStatusOverdue
		s.db.Save(&invoice)

		// Suspend VM if exists
		if invoice.VMID != nil {
			var vm models.VM
			if err := s.db.First(&vm, *invoice.VMID).Error; err == nil {
				// Only suspend if still running
				if vm.Status == models.VMStatusRunning {
					vm.Status = models.VMStatusSuspended
					s.db.Save(&vm)
					suspendedCount++
				}
			}
		}
	}

	return suspendedCount, nil
}

// DeleteTerminatedVMs deletes VMs with very overdue invoices
func (s *BillingService) DeleteTerminatedVMs() (int, error) {
	// Get very overdue invoices (14+ days)
	var invoices []models.Invoice

	now := time.Now()
	deleteDate := now.AddDate(0, 0, -s.cfg.DeleteDays)

	if err := s.db.Where("due_date < ? AND status = ?", deleteDate, models.InvoiceStatusOverdue).Find(&invoices).Error; err != nil {
		return 0, err
	}

	deletedCount := 0

	for _, invoice := range invoices {
		// Delete VM if exists
		if invoice.VMID != nil {
			var vm models.VM
			if err := s.db.First(&vm, *invoice.VMID).Error; err == nil {
				// Delete VM (Proxmox integration would go here)
				// For now, just mark as deleted in database
				s.db.Delete(&vm)
				deletedCount++

				// Update invoice status
				invoice.Status = models.InvoiceStatusCancelled
				s.db.Save(&invoice)
			}
		}
	}

	return deletedCount, nil
}

// getNextInvoiceNumber gets the next invoice number for the day
func (s *BillingService) getNextInvoiceNumber() int {
	var count int64
	s.db.Model(&models.Invoice{}).Where("created_at >= ?", time.Now().Truncate(24*time.Hour)).Count(&count)
	return int(count) + 1
}
