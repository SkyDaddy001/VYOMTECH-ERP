package services

import (
	"database/sql"
	"fmt"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
)

type SalesService struct {
	DB *sql.DB
}

// NewSalesService creates a new sales service
func NewSalesService(db *sql.DB) *SalesService {
	return &SalesService{
		DB: db,
	}
}

// ============================================================================
// SALES INVOICE TO GL INTEGRATION
// ============================================================================

// PostInvoiceToGL posts sales invoice to General Ledger with proper double-entry
// Creates debit/credit entries for:
// - Accounts Receivable (Debit) - customer owes us this amount
// - Sales Revenue (Credit) - income earned from sale
// - Output Tax Payable (Credit) - tax we owe to government
// Formula: AR (Debit) = Revenue (Credit) + Tax (Credit)
// This properly balances: DR AR = CR Revenue + CR Tax
func (s *SalesService) PostInvoiceToGL(tenantID, invoiceID string, glService *GLService, postedBy string) (string, error) {
	// Get invoice details from database
	var invoiceAmount, taxAmount, discountAmount float64
	var invoiceNumber, customerName string
	var invoiceDate time.Time

	query := `SELECT invoice_number, customer_id, invoice_date, 
		invoice_amount, tax_amount, discount_amount 
		FROM sales_invoices WHERE id = ? AND tenant_id = ?`

	err := s.DB.QueryRow(query, invoiceID, tenantID).Scan(
		&invoiceNumber, &customerName, &invoiceDate,
		&invoiceAmount, &taxAmount, &discountAmount,
	)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("invoice not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to get invoice: %w", err)
	}

	// Calculate net revenue (after discount)
	netRevenue := invoiceAmount - discountAmount

	// Create journal entry header for sales posting
	journalEntry := &models.JournalEntry{
		ID:              fmt.Sprintf("JE-SALES-INV-%s", invoiceID),
		TenantID:        tenantID,
		EntryDate:       time.Now(),
		ReferenceNumber: &invoiceNumber,
		ReferenceType:   "Sales_Invoice",
		ReferenceID:     &invoiceID,
		Description:     fmt.Sprintf("Sales invoice to %s", customerName),
		Amount:          (netRevenue + taxAmount), // Total AR amount
		Narration:       fmt.Sprintf("Invoice %s dated %s for ₹%.2f", invoiceNumber, invoiceDate.Format("02-Jan-2006"), netRevenue),
		EntryStatus:     "Draft",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Create journal entry in GL service
	if err := glService.CreateJournalEntry(tenantID, journalEntry); err != nil {
		return "", fmt.Errorf("failed to create journal entry: %w", err)
	}

	lineNum := 1

	// Add debit line: Accounts Receivable
	// DR: AR = what customer owes us (revenue + tax)
	arDetail := &models.JournalEntryDetail{
		ID:             fmt.Sprintf("JED-%s-AR", invoiceID),
		TenantID:       tenantID,
		JournalEntryID: journalEntry.ID,
		AccountID:      "ACC-ACCOUNTS-RECEIVABLE", // AR account
		DebitAmount:    (netRevenue + taxAmount),
		CreditAmount:   0,
		Description:    fmt.Sprintf("Accounts receivable - %s", customerName),
		LineNumber:     lineNum,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := glService.AddJournalEntryDetail(arDetail); err != nil {
		return "", fmt.Errorf("failed to add AR detail: %w", err)
	}
	lineNum++

	// Add credit line: Sales Revenue (the income earned)
	// CR: Sales Revenue = net amount earned from the sale
	revenueDetail := &models.JournalEntryDetail{
		ID:             fmt.Sprintf("JED-%s-REV", invoiceID),
		TenantID:       tenantID,
		JournalEntryID: journalEntry.ID,
		AccountID:      "ACC-SALES-REVENUE", // Revenue account
		DebitAmount:    0,
		CreditAmount:   netRevenue, // Credit the actual revenue earned
		Description:    "Sales revenue earned",
		LineNumber:     lineNum,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := glService.AddJournalEntryDetail(revenueDetail); err != nil {
		return "", fmt.Errorf("failed to add revenue detail: %w", err)
	}
	lineNum++

	// Add credit line: Output Tax Payable (if applicable - for GST output tax, etc.)
	// CR: Output Tax Payable = tax we owe to government
	if taxAmount > 0 {
		outputTaxDetail := &models.JournalEntryDetail{
			ID:             fmt.Sprintf("JED-%s-OTAX", invoiceID),
			TenantID:       tenantID,
			JournalEntryID: journalEntry.ID,
			AccountID:      "ACC-OUTPUT-TAX", // Output Tax/GST Payable
			DebitAmount:    0,
			CreditAmount:   taxAmount,
			Description:    "Output GST/Sales Tax payable to government",
			LineNumber:     lineNum,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := glService.AddJournalEntryDetail(outputTaxDetail); err != nil {
			return "", fmt.Errorf("failed to add output tax detail: %w", err)
		}
	}

	// Post the journal entry (validates debit=credit balance)
	// Validates: DR AR = CR Revenue + CR Tax
	if err := glService.PostJournalEntry(tenantID, journalEntry.ID, postedBy); err != nil {
		return "", fmt.Errorf("failed to post journal entry: %w", err)
	}

	// Update invoice status to indicate GL posting
	updateQuery := `UPDATE sales_invoices SET invoice_status = 'posted_to_gl', updated_at = ? 
		WHERE id = ? AND tenant_id = ?`
	_, err = s.DB.Exec(updateQuery, time.Now(), invoiceID, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to update invoice status: %w", err)
	}

	return journalEntry.ID, nil
}

// PostPaymentToGL posts sales payment/receipt to GL
// When customer pays an invoice, we reduce AR and increase Cash/Bank
// DR: Cash/Bank (what we receive)
// CR: Accounts Receivable (what customer paid down)
func (s *SalesService) PostPaymentToGL(tenantID, paymentID string, glService *GLService, postedBy string) (string, error) {
	// Get payment details from database
	var invoiceID, paymentNumber string
	var paymentAmount float64
	var paymentDate time.Time

	query := `SELECT id, invoice_id, payment_number, payment_amount, payment_date 
		FROM sales_payments WHERE id = ? AND tenant_id = ?`

	err := s.DB.QueryRow(query, paymentID, tenantID).Scan(
		&invoiceID, &paymentNumber, &paymentAmount, &paymentDate,
	)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("payment not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to get payment: %w", err)
	}

	// Create journal entry header for payment posting
	journalEntry := &models.JournalEntry{
		ID:              fmt.Sprintf("JE-SALES-PAY-%s", paymentID),
		TenantID:        tenantID,
		EntryDate:       time.Now(),
		ReferenceNumber: &paymentNumber,
		ReferenceType:   "Sales_Payment",
		ReferenceID:     &paymentID,
		Description:     "Payment received against invoice",
		Amount:          paymentAmount,
		Narration:       fmt.Sprintf("Payment %s received on %s for ₹%.2f", paymentNumber, paymentDate.Format("02-Jan-2006"), paymentAmount),
		EntryStatus:     "Draft",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Create journal entry in GL service
	if err := glService.CreateJournalEntry(tenantID, journalEntry); err != nil {
		return "", fmt.Errorf("failed to create journal entry: %w", err)
	}

	// Add debit line: Cash/Bank account
	// DR: Cash/Bank = what we received from customer
	cashDetail := &models.JournalEntryDetail{
		ID:             fmt.Sprintf("JED-%s-CASH", paymentID),
		TenantID:       tenantID,
		JournalEntryID: journalEntry.ID,
		AccountID:      "ACC-BANK-CASH", // Cash/Bank account
		DebitAmount:    paymentAmount,
		CreditAmount:   0,
		Description:    "Cash/Bank receipt from customer",
		LineNumber:     1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := glService.AddJournalEntryDetail(cashDetail); err != nil {
		return "", fmt.Errorf("failed to add cash detail: %w", err)
	}

	// Add credit line: Accounts Receivable
	// CR: AR = reduction in amount customer owes
	arDetail := &models.JournalEntryDetail{
		ID:             fmt.Sprintf("JED-%s-AR", paymentID),
		TenantID:       tenantID,
		JournalEntryID: journalEntry.ID,
		AccountID:      "ACC-ACCOUNTS-RECEIVABLE", // AR account
		DebitAmount:    0,
		CreditAmount:   paymentAmount,
		Description:    "Accounts receivable collection",
		LineNumber:     2,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := glService.AddJournalEntryDetail(arDetail); err != nil {
		return "", fmt.Errorf("failed to add AR detail: %w", err)
	}

	// Post the journal entry (validates debit=credit balance)
	// Validates: DR Cash = CR AR
	if err := glService.PostJournalEntry(tenantID, journalEntry.ID, postedBy); err != nil {
		return "", fmt.Errorf("failed to post journal entry: %w", err)
	}

	// Update payment status to indicate GL posting
	updateQuery := `UPDATE sales_payments SET payment_status = 'posted_to_gl', updated_at = ? 
		WHERE id = ? AND tenant_id = ?`
	_, err = s.DB.Exec(updateQuery, time.Now(), paymentID, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to update payment status: %w", err)
	}

	return journalEntry.ID, nil
}
