package services

import (
	"database/sql"
	"fmt"
	"time"

	"vyomtech-backend/internal/models"
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

// ============================================================================
// SALES DASHBOARD QUERY METHODS
// ============================================================================

// GetSalesOverviewMetrics returns aggregated sales overview metrics
func (s *SalesService) GetSalesOverviewMetrics(tenantID string) (map[string]interface{}, error) {
	metrics := map[string]interface{}{
		"ytd_revenue":           0.0,
		"current_month_revenue": 0.0,
		"total_opportunities":   0,
		"pipeline_value":        0.0,
		"conversion_rate":       0.0,
		"top_customers":         []map[string]interface{}{},
	}

	// Query YTD revenue from invoices
	query := `
		SELECT COALESCE(SUM(invoice_amount), 0) as ytd_revenue
		FROM sales_invoices
		WHERE tenant_id = ? AND deleted_at IS NULL 
		AND YEAR(invoice_date) = YEAR(NOW())
	`

	var ytdRevenue float64
	err := s.DB.QueryRow(query, tenantID).Scan(&ytdRevenue)
	if err != nil && err != sql.ErrNoRows {
		return metrics, fmt.Errorf("failed to query YTD revenue: %w", err)
	}
	metrics["ytd_revenue"] = ytdRevenue

	// Query current month revenue
	query = `
		SELECT COALESCE(SUM(invoice_amount), 0) as month_revenue
		FROM sales_invoices
		WHERE tenant_id = ? AND deleted_at IS NULL 
		AND YEAR(invoice_date) = YEAR(NOW()) 
		AND MONTH(invoice_date) = MONTH(NOW())
	`

	var monthRevenue float64
	err = s.DB.QueryRow(query, tenantID).Scan(&monthRevenue)
	if err != nil && err != sql.ErrNoRows {
		return metrics, fmt.Errorf("failed to query monthly revenue: %w", err)
	}
	metrics["current_month_revenue"] = monthRevenue

	// Query pipeline opportunities
	query = `
		SELECT COUNT(*) as total_ops, COALESCE(SUM(expected_value), 0) as pipeline_value
		FROM sales_opportunities
		WHERE tenant_id = ? AND deleted_at IS NULL AND stage NOT IN ('Closed-Won', 'Closed-Lost')
	`

	var totalOps int
	var pipelineValue float64
	err = s.DB.QueryRow(query, tenantID).Scan(&totalOps, &pipelineValue)
	if err != nil && err != sql.ErrNoRows {
		return metrics, fmt.Errorf("failed to query opportunities: %w", err)
	}
	metrics["total_opportunities"] = totalOps
	metrics["pipeline_value"] = pipelineValue

	// Calculate conversion rate
	if totalOps > 0 {
		var closedWon int
		query = `
			SELECT COUNT(*) FROM sales_opportunities
			WHERE tenant_id = ? AND deleted_at IS NULL AND stage = 'Closed-Won'
		`
		_ = s.DB.QueryRow(query, tenantID).Scan(&closedWon)
		metrics["conversion_rate"] = (float64(closedWon) / float64(totalOps)) * 100
	}

	return metrics, nil
}

// GetPipelineAnalysisMetrics returns sales pipeline by stage analysis
func (s *SalesService) GetPipelineAnalysisMetrics(tenantID string) (map[string]interface{}, error) {
	metrics := map[string]interface{}{
		"by_stage": map[string]interface{}{},
	}

	// Query opportunities by stage
	query := `
		SELECT stage, COUNT(*) as count, COALESCE(SUM(expected_value), 0) as total_value,
		       AVG(DATEDIFF(NOW(), created_at)) as avg_days_in_stage
		FROM sales_opportunities
		WHERE tenant_id = ? AND deleted_at IS NULL
		GROUP BY stage
	`

	rows, err := s.DB.Query(query, tenantID)
	if err != nil {
		return metrics, fmt.Errorf("failed to query pipeline: %w", err)
	}
	defer rows.Close()

	byStage := make(map[string]interface{})
	for rows.Next() {
		var stage string
		var count int
		var totalValue float64
		var avgDays int

		if err := rows.Scan(&stage, &count, &totalValue, &avgDays); err != nil {
			continue
		}

		byStage[stage] = map[string]interface{}{
			"opportunity_count": count,
			"total_value":       totalValue,
			"avg_age_days":      avgDays,
		}
	}

	metrics["by_stage"] = byStage
	return metrics, nil
}

// GetSalesMetricsForPeriod returns sales metrics for a specified period
func (s *SalesService) GetSalesMetricsForPeriod(tenantID string, startDate, endDate time.Time) (map[string]interface{}, error) {
	metrics := map[string]interface{}{
		"period_revenue":        0.0,
		"invoice_count":         0,
		"average_invoice_value": 0.0,
		"collection_rate":       0.0,
		"top_customers":         []map[string]interface{}{},
		"top_products":          []map[string]interface{}{},
	}

	// Query invoices for period
	query := `
		SELECT COUNT(*) as invoice_count, 
		       COALESCE(SUM(invoice_amount), 0) as total_revenue,
		       COALESCE(AVG(invoice_amount), 0) as avg_invoice
		FROM sales_invoices
		WHERE tenant_id = ? AND deleted_at IS NULL 
		AND invoice_date BETWEEN ? AND ?
	`

	var invoiceCount int
	var totalRevenue, avgInvoice float64

	err := s.DB.QueryRow(query, tenantID, startDate, endDate).Scan(&invoiceCount, &totalRevenue, &avgInvoice)
	if err != nil && err != sql.ErrNoRows {
		return metrics, fmt.Errorf("failed to query period metrics: %w", err)
	}

	metrics["invoice_count"] = invoiceCount
	metrics["period_revenue"] = totalRevenue
	metrics["average_invoice_value"] = avgInvoice

	// Calculate collection rate from payments
	query = `
		SELECT COALESCE(SUM(payment_amount), 0) as total_collected
		FROM sales_payments
		WHERE tenant_id = ? AND deleted_at IS NULL 
		AND payment_date BETWEEN ? AND ?
	`

	var totalCollected float64
	err = s.DB.QueryRow(query, tenantID, startDate, endDate).Scan(&totalCollected)
	if err == nil && totalRevenue > 0 {
		metrics["collection_rate"] = (totalCollected / totalRevenue) * 100
	}

	return metrics, nil
}

// GetInvoiceStatusMetrics returns invoice aging and status metrics
func (s *SalesService) GetInvoiceStatusMetrics(tenantID string) (map[string]interface{}, error) {
	metrics := map[string]interface{}{
		"total_outstanding": 0.0,
		"overdue_count":     0,
		"average_dso":       0.0,
		"aging_buckets": map[string]interface{}{
			"current": 0.0,
			"30_days": 0.0,
			"60_days": 0.0,
			"90_days": 0.0,
			"over_90": 0.0,
		},
	}

	// Query outstanding invoices
	query := `
		SELECT COALESCE(SUM(i.invoice_amount - COALESCE(p.paid_amount, 0)), 0) as outstanding,
		       COUNT(DISTINCT CASE WHEN i.due_date < NOW() THEN i.id END) as overdue_count
		FROM sales_invoices i
		LEFT JOIN (
			SELECT invoice_id, SUM(payment_amount) as paid_amount
			FROM sales_payments
			WHERE tenant_id = ? AND payment_status = 'Completed'
			GROUP BY invoice_id
		) p ON i.id = p.invoice_id
		WHERE i.tenant_id = ? AND i.deleted_at IS NULL AND i.invoice_status NOT IN ('Cancelled', 'Paid')
	`

	var outstanding float64
	var overdueCount int

	err := s.DB.QueryRow(query, tenantID, tenantID).Scan(&outstanding, &overdueCount)
	if err != nil && err != sql.ErrNoRows {
		return metrics, fmt.Errorf("failed to query invoice status: %w", err)
	}

	metrics["total_outstanding"] = outstanding
	metrics["overdue_count"] = overdueCount

	// Calculate DSO (Days Sales Outstanding)
	query = `
		SELECT AVG(DATEDIFF(NOW(), invoice_date)) as avg_dso
		FROM sales_invoices
		WHERE tenant_id = ? AND deleted_at IS NULL 
		AND invoice_status NOT IN ('Cancelled', 'Paid')
		AND invoice_date > DATE_SUB(NOW(), INTERVAL 90 DAY)
	`

	var avgDSO int
	err = s.DB.QueryRow(query, tenantID).Scan(&avgDSO)
	if err == nil {
		metrics["average_dso"] = avgDSO
	}

	return metrics, nil
}
