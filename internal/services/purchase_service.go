package services

import (
	"database/sql"
	"fmt"
	"time"

	"vyomtech-backend/internal/models"
)

// PurchaseService handles Purchase Management and GL Integration
type PurchaseService struct {
	DB *sql.DB
}

// NewPurchaseService creates a new Purchase service
func NewPurchaseService(db *sql.DB) *PurchaseService {
	return &PurchaseService{DB: db}
}

// ============================================================================
// VENDOR MANAGEMENT
// ============================================================================

// CreateVendor creates a new vendor
func (s *PurchaseService) CreateVendor(tenantID string, vendor *models.Vendor) error {
	vendor.TenantID = tenantID
	vendor.CreatedAt = time.Now()
	vendor.UpdatedAt = time.Now()

	query := `INSERT INTO vendors (
		id, tenant_id, vendor_code, name, vendor_type, email, phone,
		address, city, state, country, postal_code, tax_id, payment_terms,
		rating, is_active, created_at, updated_at, status
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		vendor.ID, vendor.TenantID, vendor.VendorCode, vendor.Name, vendor.VendorType,
		vendor.Email, vendor.Phone, vendor.Address, vendor.City, vendor.State,
		vendor.Country, vendor.PostalCode, vendor.TaxID, vendor.PaymentTerms,
		vendor.Rating, vendor.IsActive, vendor.CreatedAt, vendor.UpdatedAt, vendor.Status,
	)

	return err
}

// GetVendor retrieves a vendor by ID
func (s *PurchaseService) GetVendor(tenantID, vendorID string) (*models.Vendor, error) {
	var vendor models.Vendor
	query := `SELECT id, tenant_id, vendor_code, name, vendor_type, email, phone,
		address, city, state, country, postal_code, tax_id, payment_terms,
		rating, is_active, is_blocked, created_at, updated_at, deleted_at, status
		FROM vendors WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	err := s.DB.QueryRow(query, vendorID, tenantID).Scan(
		&vendor.ID, &vendor.TenantID, &vendor.VendorCode, &vendor.Name, &vendor.VendorType,
		&vendor.Email, &vendor.Phone, &vendor.Address, &vendor.City, &vendor.State,
		&vendor.Country, &vendor.PostalCode, &vendor.TaxID, &vendor.PaymentTerms,
		&vendor.Rating, &vendor.IsActive, &vendor.IsBlocked, &vendor.CreatedAt, &vendor.UpdatedAt, &vendor.DeletedAt, &vendor.Status,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("vendor not found")
	}
	return &vendor, err
}

// ListVendors retrieves all active vendors
func (s *PurchaseService) ListVendors(tenantID string) ([]models.Vendor, error) {
	var vendors []models.Vendor
	query := `SELECT id, tenant_id, vendor_code, name, vendor_type, email, phone,
		address, city, state, country, postal_code, tax_id, payment_terms,
		rating, is_active, is_blocked, created_at, updated_at, deleted_at, status
		FROM vendors WHERE tenant_id = ? AND deleted_at IS NULL ORDER BY name`

	rows, err := s.DB.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vendor models.Vendor
		err := rows.Scan(
			&vendor.ID, &vendor.TenantID, &vendor.VendorCode, &vendor.Name, &vendor.VendorType,
			&vendor.Email, &vendor.Phone, &vendor.Address, &vendor.City, &vendor.State,
			&vendor.Country, &vendor.PostalCode, &vendor.TaxID, &vendor.PaymentTerms,
			&vendor.Rating, &vendor.IsActive, &vendor.IsBlocked, &vendor.CreatedAt, &vendor.UpdatedAt, &vendor.DeletedAt, &vendor.Status,
		)
		if err != nil {
			return nil, err
		}
		vendors = append(vendors, vendor)
	}

	return vendors, rows.Err()
}

// ============================================================================
// PURCHASE ORDER MANAGEMENT
// ============================================================================

// CreatePurchaseOrder creates a new PO
func (s *PurchaseService) CreatePurchaseOrder(tenantID string, po *models.PurchaseOrder) error {
	po.TenantID = tenantID
	po.CreatedAt = time.Now()
	po.UpdatedAt = time.Now()

	query := `INSERT INTO purchase_orders (
		id, tenant_id, po_number, vendor_id, po_date, delivery_date,
		total_amount, tax_amount, shipping_amount, discount_amount, net_amount,
		payment_terms, delivery_location, special_instructions, status,
		created_at, updated_at, created_by
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		po.ID, po.TenantID, po.PONumber, po.VendorID, po.PODate, po.DeliveryDate,
		po.TotalAmount, po.TaxAmount, po.ShippingAmount, po.DiscountAmount, po.NetAmount,
		po.PaymentTerms, po.DeliveryLocation, po.SpecialInstructions, po.Status,
		po.CreatedAt, po.UpdatedAt, po.CreatedBy,
	)

	return err
}

// GetPurchaseOrder retrieves a PO by ID
func (s *PurchaseService) GetPurchaseOrder(tenantID, poID string) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	query := `SELECT id, tenant_id, po_number, vendor_id, po_date, delivery_date,
		total_amount, tax_amount, shipping_amount, discount_amount, net_amount,
		payment_terms, delivery_location, special_instructions, status,
		created_at, updated_at, deleted_at
		FROM purchase_orders WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	err := s.DB.QueryRow(query, poID, tenantID).Scan(
		&po.ID, &po.TenantID, &po.PONumber, &po.VendorID, &po.PODate, &po.DeliveryDate,
		&po.TotalAmount, &po.TaxAmount, &po.ShippingAmount, &po.DiscountAmount, &po.NetAmount,
		&po.PaymentTerms, &po.DeliveryLocation, &po.SpecialInstructions, &po.Status,
		&po.CreatedAt, &po.UpdatedAt, &po.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("purchase order not found")
	}

	return &po, err
}

// ============================================================================
// PURCHASE INVOICE MANAGEMENT
// ============================================================================

// CreateVendorInvoice creates a new vendor invoice
func (s *PurchaseService) CreateVendorInvoice(tenantID string, invoice *models.VendorInvoice) error {
	invoice.TenantID = tenantID
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	query := `INSERT INTO vendor_invoices (
		id, tenant_id, invoice_number, vendor_id, po_id, grn_id, invoice_date,
		due_date, invoice_amount, tax_amount, discount_amount, total_payable,
		status, matched_status, three_way_match, created_at, updated_at, created_by
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		invoice.ID, invoice.TenantID, invoice.InvoiceNumber, invoice.VendorID,
		invoice.POID, invoice.GRNID, invoice.InvoiceDate, invoice.DueDate,
		invoice.InvoiceAmount, invoice.TaxAmount, invoice.DiscountAmount,
		invoice.TotalPayable, invoice.Status, invoice.MatchedStatus, invoice.ThreeWayMatch,
		invoice.CreatedAt, invoice.UpdatedAt, invoice.CreatedBy,
	)

	return err
}

// GetVendorInvoice retrieves an invoice by ID
func (s *PurchaseService) GetVendorInvoice(tenantID, invoiceID string) (*models.VendorInvoice, error) {
	var invoice models.VendorInvoice
	query := `SELECT id, tenant_id, invoice_number, vendor_id, po_id, grn_id,
		invoice_date, due_date, invoice_amount, tax_amount, discount_amount,
		total_payable, status, matched_status, three_way_match, approved_at,
		approved_by, rejection_reason, created_at, updated_at
		FROM vendor_invoices WHERE id = ? AND tenant_id = ?`

	err := s.DB.QueryRow(query, invoiceID, tenantID).Scan(
		&invoice.ID, &invoice.TenantID, &invoice.InvoiceNumber, &invoice.VendorID,
		&invoice.POID, &invoice.GRNID, &invoice.InvoiceDate, &invoice.DueDate,
		&invoice.InvoiceAmount, &invoice.TaxAmount, &invoice.DiscountAmount,
		&invoice.TotalPayable, &invoice.Status, &invoice.MatchedStatus, &invoice.ThreeWayMatch,
		&invoice.ApprovedAt, &invoice.ApprovedBy, &invoice.RejectionReason,
		&invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invoice not found")
	}

	return &invoice, err
}

// ============================================================================
// PURCHASE INVOICE TO GL INTEGRATION
// ============================================================================

// PostInvoiceToGL posts purchase invoice to General Ledger
// Creates debit/credit entries for:
// - Expense/Asset account (Debit)
// - Accounts Payable (Credit)
// - Tax Payable (Credit)
func (s *PurchaseService) PostInvoiceToGL(tenantID, invoiceID string, glService *GLService, postedBy string) (string, error) {
	// Get invoice details
	invoice, err := s.GetVendorInvoice(tenantID, invoiceID)
	if err != nil {
		return "", fmt.Errorf("failed to get invoice: %w", err)
	}

	// Get vendor details for reference
	vendor, err := s.GetVendor(tenantID, invoice.VendorID)
	if err != nil {
		return "", fmt.Errorf("failed to get vendor: %w", err)
	}

	// Create journal entry header for invoice posting
	journalEntry := &models.JournalEntry{
		ID:              fmt.Sprintf("JE-PO-INVOICE-%s", invoiceID),
		TenantID:        tenantID,
		EntryDate:       time.Now(),
		ReferenceNumber: &invoice.InvoiceNumber,
		ReferenceType:   "Purchase_Invoice",
		ReferenceID:     &invoiceID,
		Description:     fmt.Sprintf("Purchase invoice from %s", vendor.Name),
		Amount:          invoice.TotalPayable,
		Narration:       fmt.Sprintf("Invoice %s dated %s", invoice.InvoiceNumber, invoice.InvoiceDate.Format("02-Jan-2006")),
		EntryStatus:     "Draft",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Create journal entry in GL service
	if err := glService.CreateJournalEntry(tenantID, journalEntry); err != nil {
		return "", fmt.Errorf("failed to create journal entry: %w", err)
	}

	// Add debit line: Expense/Asset account (based on PO)
	// For purchase of materials, debit inventory/asset account
	// For purchase of services, debit expense account
	lineAmount := invoice.InvoiceAmount - invoice.DiscountAmount
	expenseDetail := &models.JournalEntryDetail{
		ID:             fmt.Sprintf("JED-%s-EXP", invoiceID),
		TenantID:       tenantID,
		JournalEntryID: journalEntry.ID,
		AccountID:      "ACC-PURCHASE-EXPENSE", // Should be configured per tenant
		DebitAmount:    lineAmount,
		CreditAmount:   0,
		Description:    "Purchase expense/inventory",
		LineNumber:     1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := glService.AddJournalEntryDetail(expenseDetail); err != nil {
		return "", fmt.Errorf("failed to add expense detail: %w", err)
	}

	lineNum := 2

	// Add debit line for Input Tax (if applicable - for GST input credit, etc.)
	// DR: Input Tax Receivable (GST can be claimed as credit)
	// This reduces the amount payable to government
	if invoice.TaxAmount > 0 {
		taxDetail := &models.JournalEntryDetail{
			ID:             fmt.Sprintf("JED-%s-TAX", invoiceID),
			TenantID:       tenantID,
			JournalEntryID: journalEntry.ID,
			AccountID:      "ACC-INPUT-TAX", // GST Input Tax Receivable
			DebitAmount:    invoice.TaxAmount,
			CreditAmount:   0,
			Description:    "Input GST/Tax (recoverable)",
			LineNumber:     lineNum,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := glService.AddJournalEntryDetail(taxDetail); err != nil {
			return "", fmt.Errorf("failed to add tax detail: %w", err)
		}
		lineNum++
	}

	// Add credit line: Accounts Payable (main liability)
	// CR: AP = Purchase amount + Tax amount (vendor owes this total)
	totalPayable := lineAmount + invoice.TaxAmount
	apDetail := &models.JournalEntryDetail{
		ID:             fmt.Sprintf("JED-%s-AP", invoiceID),
		TenantID:       tenantID,
		JournalEntryID: journalEntry.ID,
		AccountID:      "ACC-ACCOUNTS-PAYABLE", // Should be configured per tenant
		DebitAmount:    0,
		CreditAmount:   totalPayable,
		Description:    fmt.Sprintf("Accounts payable - %s", vendor.Name),
		LineNumber:     lineNum,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := glService.AddJournalEntryDetail(apDetail); err != nil {
		return "", fmt.Errorf("failed to add AP detail: %w", err)
	}

	// Post the journal entry (validates debit=credit balance)
	if err := glService.PostJournalEntry(tenantID, journalEntry.ID, postedBy); err != nil {
		return "", fmt.Errorf("failed to post journal entry: %w", err)
	}

	// Update invoice status to indicate GL posting
	updateQuery := `UPDATE vendor_invoices SET status = 'posted_to_gl', updated_at = ? 
		WHERE id = ? AND tenant_id = ?`
	_, err = s.DB.Exec(updateQuery, time.Now(), invoiceID, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to update invoice status: %w", err)
	}

	return journalEntry.ID, nil
}

// ============================================================================
// PURCHASE PAYMENT TO GL INTEGRATION
// ============================================================================

// PostPaymentToGL posts purchase payment to General Ledger
// When we pay vendor invoice:
// - DR: Accounts Payable (reduction in liability)
// - CR: Cash/Bank (cash outflow)
func (s *PurchaseService) PostPaymentToGL(tenantID, paymentID string, glService *GLService, postedBy string) (string, error) {
	// Get payment details from database
	var invoiceID, paymentNumber string
	var paymentAmount float64
	var paymentDate time.Time

	query := `SELECT id, invoice_id, payment_number, payment_amount, payment_date 
		FROM purchase_payments WHERE id = ? AND tenant_id = ?`

	err := s.DB.QueryRow(query, paymentID, tenantID).Scan(
		&invoiceID, &paymentNumber, &paymentAmount, &paymentDate,
	)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("payment not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to get payment: %w", err)
	}

	// Get invoice to find vendor name
	invoice, err := s.GetVendorInvoice(tenantID, invoiceID)
	if err != nil {
		return "", fmt.Errorf("failed to get invoice: %w", err)
	}

	vendor, err := s.GetVendor(tenantID, invoice.VendorID)
	if err != nil {
		return "", fmt.Errorf("failed to get vendor: %w", err)
	}

	// Create journal entry header for payment posting
	journalEntry := &models.JournalEntry{
		ID:              fmt.Sprintf("JE-PO-PAY-%s", paymentID),
		TenantID:        tenantID,
		EntryDate:       time.Now(),
		ReferenceNumber: &paymentNumber,
		ReferenceType:   "Purchase_Payment",
		ReferenceID:     &paymentID,
		Description:     fmt.Sprintf("Payment to %s", vendor.Name),
		Amount:          paymentAmount,
		Narration:       fmt.Sprintf("Payment %s made on %s for ₹%.2f", paymentNumber, paymentDate.Format("02-Jan-2006"), paymentAmount),
		EntryStatus:     "Draft",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Create journal entry in GL service
	if err := glService.CreateJournalEntry(tenantID, journalEntry); err != nil {
		return "", fmt.Errorf("failed to create journal entry: %w", err)
	}

	// Add debit line: Accounts Payable
	// DR: AP = reduction in amount we owe vendor
	apDetail := &models.JournalEntryDetail{
		ID:             fmt.Sprintf("JED-%s-AP", paymentID),
		TenantID:       tenantID,
		JournalEntryID: journalEntry.ID,
		AccountID:      "ACC-ACCOUNTS-PAYABLE", // AP account
		DebitAmount:    paymentAmount,
		CreditAmount:   0,
		Description:    fmt.Sprintf("Accounts payable payment to %s", vendor.Name),
		LineNumber:     1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := glService.AddJournalEntryDetail(apDetail); err != nil {
		return "", fmt.Errorf("failed to add AP detail: %w", err)
	}

	// Add credit line: Cash/Bank account
	// CR: Cash/Bank = cash outflow to pay vendor
	cashDetail := &models.JournalEntryDetail{
		ID:             fmt.Sprintf("JED-%s-CASH", paymentID),
		TenantID:       tenantID,
		JournalEntryID: journalEntry.ID,
		AccountID:      "ACC-BANK-CASH", // Cash/Bank account
		DebitAmount:    0,
		CreditAmount:   paymentAmount,
		Description:    "Cash/Bank payment to vendor",
		LineNumber:     2,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := glService.AddJournalEntryDetail(cashDetail); err != nil {
		return "", fmt.Errorf("failed to add cash detail: %w", err)
	}

	// Post the journal entry (validates debit=credit balance)
	// Validates: DR AP = CR Cash
	if err := glService.PostJournalEntry(tenantID, journalEntry.ID, postedBy); err != nil {
		return "", fmt.Errorf("failed to post journal entry: %w", err)
	}

	// Update payment status to indicate GL posting
	updateQuery := `UPDATE purchase_payments SET payment_status = 'posted_to_gl', updated_at = ? 
		WHERE id = ? AND tenant_id = ?`
	_, err = s.DB.Exec(updateQuery, time.Now(), paymentID, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to update payment status: %w", err)
	}

	return journalEntry.ID, nil
}
