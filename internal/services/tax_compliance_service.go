package services

import (
	"database/sql"
	"fmt"
	"time"

	"vyomtech-backend/internal/models"
)

// ============================================================================
// TAX COMPLIANCE SERVICE
// ============================================================================
// Manages Income Tax, GST, TDS, and other tax compliance

type TaxComplianceService struct {
	DB *sql.DB
}

func NewTaxComplianceService(db *sql.DB) *TaxComplianceService {
	return &TaxComplianceService{DB: db}
}

// SetupTaxConfiguration initializes tax configuration for an organization
func (s *TaxComplianceService) SetupTaxConfiguration(
	tenantID string,
	pan, gstNumber string,
	gstRegistrationDate time.Time,
) (*models.TaxConfiguration, error) {

	// Determine financial year based on current date
	currentYear := time.Now().Year()
	currentMonth := time.Now().Month()
	var fiscalYear int
	var fyStart, fyEnd time.Time

	if currentMonth >= time.April {
		fiscalYear = currentYear
		fyStart = time.Date(currentYear, time.April, 1, 0, 0, 0, 0, time.UTC)
		fyEnd = time.Date(currentYear+1, time.March, 31, 23, 59, 59, 0, time.UTC)
	} else {
		fiscalYear = currentYear - 1
		fyStart = time.Date(currentYear-1, time.April, 1, 0, 0, 0, 0, time.UTC)
		fyEnd = time.Date(currentYear, time.March, 31, 23, 59, 59, 0, time.UTC)
	}

	config := &models.TaxConfiguration{
		ID:                       fmt.Sprintf("TAX-CFG-%s", tenantID),
		TenantID:                 tenantID,
		FinancialYearStart:       fyStart,
		FinancialYearEnd:         fyEnd,
		FiscalYear:               fiscalYear,
		ITAssessmentYear:         fiscalYear + 1,
		ITPAN:                    pan,
		GSTRegistrationNumber:    gstNumber,
		GSTRegistrationDate:      &gstRegistrationDate,
		GSTStatus:                "Active",
		GSTReturnFilingFrequency: "Monthly",
		ITFilingStatus:           "Not_Filed",
		GSTFilingStatus:          "Not_Filed",
		IsActive:                 true,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	query := `INSERT INTO tax_configuration 
		(id, tenant_id, financial_year_start, financial_year_end, fiscal_year, 
		 it_assessment_year, it_pan, gst_registration_number, gst_registration_date, 
		 gst_status, gst_return_filing_frequency, it_filing_status, gst_filing_status, 
		 is_active, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		config.ID, config.TenantID, config.FinancialYearStart, config.FinancialYearEnd,
		config.FiscalYear, config.ITAssessmentYear, config.ITPAN, config.GSTRegistrationNumber,
		config.GSTRegistrationDate, config.GSTStatus, config.GSTReturnFilingFrequency,
		config.ITFilingStatus, config.GSTFilingStatus, config.IsActive,
		config.CreatedAt, config.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to setup tax configuration: %w", err)
	}

	return config, nil
}

// InitializeIncomeTaxCompliance creates initial IT compliance record for fiscal year
func (s *TaxComplianceService) InitializeIncomeTaxCompliance(
	tenantID string,
	fiscalYear int,
) (*models.IncomeTaxCompliance, error) {

	assessmentYear := fiscalYear + 1

	itCompliance := &models.IncomeTaxCompliance{
		ID:                 fmt.Sprintf("ITC-%s-%d", tenantID, fiscalYear),
		TenantID:           tenantID,
		FiscalYear:         fiscalYear,
		AssessmentYear:     assessmentYear,
		TotalGrossIncome:   0,
		TotalDeductions:    0,
		TaxableIncome:      0,
		TaxBeforeSurcharge: 0,
		SurchargeAmount:    0,
		CessAmount:         0,
		TotalTaxLiability:  0,
		AdvanceTaxPaid:     0,
		TDSCredit:          0,
		TotalTaxPaid:       0,
		ReturnFiledStatus:  "Not_Filed",
		IsCompliant:        true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	query := `INSERT INTO income_tax_compliance 
		(id, tenant_id, fiscal_year, assessment_year, total_gross_income, total_deductions, 
		 taxable_income, tax_before_surcharge, surcharge_amount, cess_amount, 
		 total_tax_liability, advance_tax_paid, tds_credit, total_tax_paid, 
		 return_filed_status, is_compliant, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		itCompliance.ID, itCompliance.TenantID, itCompliance.FiscalYear,
		itCompliance.AssessmentYear, itCompliance.TotalGrossIncome, itCompliance.TotalDeductions,
		itCompliance.TaxableIncome, itCompliance.TaxBeforeSurcharge, itCompliance.SurchargeAmount,
		itCompliance.CessAmount, itCompliance.TotalTaxLiability, itCompliance.AdvanceTaxPaid,
		itCompliance.TDSCredit, itCompliance.TotalTaxPaid, itCompliance.ReturnFiledStatus,
		itCompliance.IsCompliant, itCompliance.CreatedAt, itCompliance.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize IT compliance: %w", err)
	}

	return itCompliance, nil
}

// InitializeGSTCompliance creates GST compliance record for a return period
func (s *TaxComplianceService) InitializeGSTCompliance(
	tenantID string,
	fiscalYear int,
	returnPeriod string,
	monthYear time.Time,
) (*models.GSTCompliance, error) {

	gstCompliance := &models.GSTCompliance{
		ID:                        fmt.Sprintf("GST-%s-%d", tenantID, monthYear.UnixNano()),
		TenantID:                  tenantID,
		FiscalYear:                fiscalYear,
		ReturnPeriod:              returnPeriod,
		MonthYear:                 monthYear,
		TotalOutwardSupplies:      0,
		TotalInwardSupplies:       0,
		TotalOutputGST:            0,
		TotalInputGST:             0,
		NetGSTPay:                 0,
		GSTRFilingStatus:          "Not_Filed",
		InterestFreeCreditCarried: 0,
		IsCompliant:               true,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}

	query := `INSERT INTO gst_compliance 
		(id, tenant_id, fiscal_year, return_period, month_year, total_outward_supplies, 
		 total_inward_supplies, total_output_gst, total_input_gst, net_gst_payable, 
		 gstr_filing_status, interest_free_credit_carried, is_compliant, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		gstCompliance.ID, gstCompliance.TenantID, gstCompliance.FiscalYear,
		gstCompliance.ReturnPeriod, gstCompliance.MonthYear, gstCompliance.TotalOutwardSupplies,
		gstCompliance.TotalInwardSupplies, gstCompliance.TotalOutputGST, gstCompliance.TotalInputGST,
		gstCompliance.NetGSTPay, gstCompliance.GSTRFilingStatus,
		gstCompliance.InterestFreeCreditCarried, gstCompliance.IsCompliant,
		gstCompliance.CreatedAt, gstCompliance.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize GST compliance: %w", err)
	}

	return gstCompliance, nil
}

// TrackGSTInvoice records GST on individual sales invoice
func (s *TaxComplianceService) TrackGSTInvoice(
	tenantID, invoiceID, invoiceNumber, customerGSTIN string,
	invoiceAmount, gstRate, gstAmount float64,
	invoiceDate time.Time,
) (*models.GSTInvoiceTracking, error) {

	tracking := &models.GSTInvoiceTracking{
		ID:                fmt.Sprintf("GSTI-%s-%d", invoiceID, time.Now().Unix()),
		TenantID:          tenantID,
		InvoiceID:         invoiceID,
		InvoiceNumber:     invoiceNumber,
		InvoiceDate:       &invoiceDate,
		InvoiceAmount:     invoiceAmount,
		CustomerGSTIN:     customerGSTIN,
		GSTRate:           gstRate,
		GSTAmount:         gstAmount,
		InvoiceRaisedDate: &invoiceDate,
		ITCEligible:       true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	query := `INSERT INTO gst_invoice_tracking 
		(id, tenant_id, invoice_id, invoice_number, invoice_date, invoice_amount, 
		 customer_gstin, gst_rate, gst_amount, invoice_raised_date, itc_eligible, 
		 created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		tracking.ID, tracking.TenantID, tracking.InvoiceID, tracking.InvoiceNumber,
		tracking.InvoiceDate, tracking.InvoiceAmount, tracking.CustomerGSTIN,
		tracking.GSTRate, tracking.GSTAmount, tracking.InvoiceRaisedDate,
		tracking.ITCEligible, tracking.CreatedAt, tracking.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to track GST invoice: %w", err)
	}

	return tracking, nil
}

// TrackGSTInputCredit records GST on purchase invoices (ITC eligibility)
func (s *TaxComplianceService) TrackGSTInputCredit(
	tenantID, purchaseInvoiceID, vendorGSTIN string,
	invoiceNumber string,
	invoiceAmount, gstRate, gstAmount float64,
	invoiceDate time.Time,
	itcEligible bool,
) (*models.GSTInputCredit, error) {

	credit := &models.GSTInputCredit{
		ID:                fmt.Sprintf("GSTC-%s-%d", purchaseInvoiceID, time.Now().Unix()),
		TenantID:          tenantID,
		PurchaseInvoiceID: purchaseInvoiceID,
		VendorGSTIN:       vendorGSTIN,
		InvoiceNumber:     invoiceNumber,
		InvoiceDate:       &invoiceDate,
		InvoiceAmount:     invoiceAmount,
		GSTRate:           gstRate,
		GSTAmount:         gstAmount,
		ITCEligible:       itcEligible,
		BlockedPercentage: 0,
		BlockedAmount:     0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	query := `INSERT INTO gst_input_credit 
		(id, tenant_id, purchase_invoice_id, vendor_gstin, invoice_number, invoice_date, 
		 invoice_amount, gst_rate, gst_amount, itc_eligible, blocked_percentage, 
		 blocked_amount, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		credit.ID, credit.TenantID, credit.PurchaseInvoiceID, credit.VendorGSTIN,
		credit.InvoiceNumber, credit.InvoiceDate, credit.InvoiceAmount,
		credit.GSTRate, credit.GSTAmount, credit.ITCEligible, credit.BlockedPercentage,
		credit.BlockedAmount, credit.CreatedAt, credit.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to track GST input credit: %w", err)
	}

	return credit, nil
}

// InitializeAdvanceTaxSchedule sets up quarterly advance tax dues
func (s *TaxComplianceService) InitializeAdvanceTaxSchedule(
	tenantID string,
	fiscalYear int,
	estimatedTaxLiability float64,
) (*models.AdvanceTaxSchedule, error) {

	// Advance tax is paid quarterly: June, September, December, March
	// Usually 30%, 60%, 100%, 0% (self-assessment at March)
	q1Due := estimatedTaxLiability * 0.30
	q2Due := estimatedTaxLiability * 0.30
	q3Due := estimatedTaxLiability * 0.40

	schedule := &models.AdvanceTaxSchedule{
		ID:                 fmt.Sprintf("ATS-%s-%d", tenantID, fiscalYear+1),
		TenantID:           tenantID,
		FiscalYear:         fiscalYear,
		AssessmentYear:     fiscalYear + 1,
		Q1DueDate:          &time.Time{}, // Will be June 15
		Q1AmountDue:        q1Due,
		Q1Status:           "Pending",
		Q2DueDate:          &time.Time{}, // Will be September 15
		Q2AmountDue:        q2Due,
		Q2Status:           "Pending",
		Q3DueDate:          &time.Time{}, // Will be December 15
		Q3AmountDue:        q3Due,
		Q3Status:           "Pending",
		Q4DueDate:          &time.Time{}, // Will be March 15
		Q4AmountDue:        0,
		Q4Status:           "Not_Due",
		TotalAdvanceTaxDue: q1Due + q2Due + q3Due,
		IsCompliant:        true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Set actual dates
	year := time.Now().Year()
	q1Date := time.Date(year, time.June, 15, 0, 0, 0, 0, time.UTC)
	q2Date := time.Date(year, time.September, 15, 0, 0, 0, 0, time.UTC)
	q3Date := time.Date(year, time.December, 15, 0, 0, 0, 0, time.UTC)
	q4Date := time.Date(year+1, time.March, 15, 0, 0, 0, 0, time.UTC)

	schedule.Q1DueDate = &q1Date
	schedule.Q2DueDate = &q2Date
	schedule.Q3DueDate = &q3Date
	schedule.Q4DueDate = &q4Date

	query := `INSERT INTO advance_tax_schedule 
		(id, tenant_id, fiscal_year, assessment_year, 
		 q1_due_date, q1_amount_due, q1_status, 
		 q2_due_date, q2_amount_due, q2_status, 
		 q3_due_date, q3_amount_due, q3_status, 
		 q4_due_date, q4_amount_due, q4_status, 
		 total_advance_tax_due, is_compliant, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		schedule.ID, schedule.TenantID, schedule.FiscalYear, schedule.AssessmentYear,
		schedule.Q1DueDate, schedule.Q1AmountDue, schedule.Q1Status,
		schedule.Q2DueDate, schedule.Q2AmountDue, schedule.Q2Status,
		schedule.Q3DueDate, schedule.Q3AmountDue, schedule.Q3Status,
		schedule.Q4DueDate, schedule.Q4AmountDue, schedule.Q4Status,
		schedule.TotalAdvanceTaxDue, schedule.IsCompliant,
		schedule.CreatedAt, schedule.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize advance tax schedule: %w", err)
	}

	return schedule, nil
}

// RecordAdvanceTaxPayment records advance tax payment for a quarter
func (s *TaxComplianceService) RecordAdvanceTaxPayment(
	tenantID string,
	fiscalYear, quarter int,
	amountPaid float64,
	challanNumber string,
) error {

	var queryStr string
	now := time.Now()

	switch quarter {
	case 1:
		queryStr = `UPDATE advance_tax_schedule 
			SET q1_amount_paid = ?, q1_payment_date = ?, q1_challan_number = ?, 
			    q1_status = 'Paid', updated_at = ? 
			WHERE tenant_id = ? AND fiscal_year = ?`
	case 2:
		queryStr = `UPDATE advance_tax_schedule 
			SET q2_amount_paid = ?, q2_payment_date = ?, q2_challan_number = ?, 
			    q2_status = 'Paid', updated_at = ? 
			WHERE tenant_id = ? AND fiscal_year = ?`
	case 3:
		queryStr = `UPDATE advance_tax_schedule 
			SET q3_amount_paid = ?, q3_payment_date = ?, q3_challan_number = ?, 
			    q3_status = 'Paid', updated_at = ? 
			WHERE tenant_id = ? AND fiscal_year = ?`
	case 4:
		queryStr = `UPDATE advance_tax_schedule 
			SET q4_amount_paid = ?, q4_payment_date = ?, q4_challan_number = ?, 
			    q4_status = 'Paid', updated_at = ? 
			WHERE tenant_id = ? AND fiscal_year = ?`
	default:
		return fmt.Errorf("invalid quarter: %d", quarter)
	}

	_, err := s.DB.Exec(queryStr, amountPaid, now, challanNumber, now, tenantID, fiscalYear)
	if err != nil {
		return fmt.Errorf("failed to record advance tax payment: %w", err)
	}

	return nil
}

// GetTaxComplianceStatus returns overall tax compliance status
func (s *TaxComplianceService) GetTaxComplianceStatus(
	tenantID string,
	fiscalYear int,
) (*models.TaxComplianceStatusResponse, error) {

	response := &models.TaxComplianceStatusResponse{
		TenantID:   tenantID,
		FiscalYear: fiscalYear,
	}

	// Get Income Tax Status
	var itFiled bool
	var itStatus string
	itQuery := `SELECT COALESCE(return_filed_status = 'Filed', false), 
		COALESCE(return_filed_status, 'Not_Filed') 
		FROM income_tax_compliance 
		WHERE tenant_id = ? AND fiscal_year = ? AND deleted_at IS NULL`

	err := s.DB.QueryRow(itQuery, tenantID, fiscalYear).Scan(&itFiled, &itStatus)
	if err == nil {
		response.IncomeTaxStatus = &models.IncomeTaxComplianceStatus{
			ITRFiled: itFiled,
			Status:   itStatus,
		}
	}

	// Get GST Status
	var gstFiled bool
	var gstStatus string
	gstQuery := `SELECT COALESCE(gstr_filing_status = 'Filed', false), 
		COALESCE(gstr_filing_status, 'Not_Filed') 
		FROM gst_compliance 
		WHERE tenant_id = ? AND fiscal_year = ? AND deleted_at IS NULL LIMIT 1`

	err = s.DB.QueryRow(gstQuery, tenantID, fiscalYear).Scan(&gstFiled, &gstStatus)
	if err == nil {
		response.GSTStatus = &models.GSTComplianceStatus{
			GSTRFiled: gstFiled,
			Status:    gstStatus,
		}
	}

	// Get Advance Tax Status
	var q1Status, q2Status, q3Status, q4Status string
	var totalPaid, totalDue float64
	atQuery := `SELECT COALESCE(q1_status, 'Not_Due'), COALESCE(q2_status, 'Not_Due'), 
		COALESCE(q3_status, 'Not_Due'), COALESCE(q4_status, 'Not_Due'), 
		COALESCE(total_advance_tax_paid, 0), COALESCE(total_advance_tax_due, 0) 
		FROM advance_tax_schedule 
		WHERE tenant_id = ? AND fiscal_year = ? AND deleted_at IS NULL`

	err = s.DB.QueryRow(atQuery, tenantID, fiscalYear).Scan(
		&q1Status, &q2Status, &q3Status, &q4Status, &totalPaid, &totalDue,
	)
	if err == nil {
		response.AdvanceTaxStatus = &models.AdvanceTaxStatus{
			Q1Status:  q1Status,
			Q2Status:  q2Status,
			Q3Status:  q3Status,
			Q4Status:  q4Status,
			TotalPaid: totalPaid,
			TotalDue:  totalDue,
		}
	}

	// Determine overall compliance
	compliant := true
	var pendingActions []string

	if response.IncomeTaxStatus != nil && response.IncomeTaxStatus.Status == "Not_Filed" {
		compliant = false
		pendingActions = append(pendingActions, "File ITR for FY "+fmt.Sprintf("%d", fiscalYear))
	}

	if response.GSTStatus != nil && response.GSTStatus.Status == "Not_Filed" {
		compliant = false
		pendingActions = append(pendingActions, "File GSTR returns")
	}

	if response.AdvanceTaxStatus != nil && (response.AdvanceTaxStatus.Q1Status == "Pending" ||
		response.AdvanceTaxStatus.Q2Status == "Pending" ||
		response.AdvanceTaxStatus.Q3Status == "Pending") {
		compliant = false
		pendingActions = append(pendingActions, "Pay pending advance tax installments")
	}

	if compliant {
		response.OverallCompliance = "Compliant"
	} else {
		response.OverallCompliance = "Non-Compliant"
	}

	response.PendingActions = pendingActions

	return response, nil
}

// ============================================================================
// TAX COMPLIANCE DASHBOARD QUERY METHODS
// ============================================================================

// GetTaxComplianceMetrics returns aggregated tax compliance metrics for dashboard
func (s *TaxComplianceService) GetTaxComplianceMetrics(tenantID string) (map[string]interface{}, error) {
	metrics := map[string]interface{}{
		"income_tax_status":  "Compliant",
		"gst_status":         "Compliant",
		"tds_status":         "Compliant",
		"advance_tax_status": "Compliant",
		"itr_filed":          false,
		"gstr_filed":         0,
		"tds_collected":      0.0,
		"advance_tax_paid":   0.0,
		"violations":         0,
	}

	// Query Income Tax compliance
	query := `
		SELECT COUNT(*) as count,
		       SUM(CASE WHEN status = 'Filed' THEN 1 ELSE 0 END) as filed,
		       SUM(CASE WHEN status = 'Not Filed' THEN 1 ELSE 0 END) as not_filed
		FROM income_tax_compliance
		WHERE tenant_id = ? AND deleted_at IS NULL
	`

	var itCount, itFiled, itNotFiled int
	err := s.DB.QueryRow(query, tenantID).Scan(&itCount, &itFiled, &itNotFiled)
	if err != nil && err != sql.ErrNoRows {
		return metrics, fmt.Errorf("failed to query IT metrics: %w", err)
	}
	if itFiled > 0 {
		metrics["itr_filed"] = true
		metrics["income_tax_status"] = "Filed"
	} else if itNotFiled > 0 {
		metrics["income_tax_status"] = "Pending"
	}

	// Query GST compliance
	query = `
		SELECT COUNT(*) as count,
		       SUM(CASE WHEN filing_status = 'Filed' THEN 1 ELSE 0 END) as filed
		FROM gst_compliance
		WHERE tenant_id = ? AND deleted_at IS NULL
	`

	var gstCount, gstFiled int
	err = s.DB.QueryRow(query, tenantID).Scan(&gstCount, &gstFiled)
	if err != nil && err != sql.ErrNoRows {
		return metrics, fmt.Errorf("failed to query GST metrics: %w", err)
	}
	metrics["gstr_filed"] = gstFiled
	if gstFiled < gstCount {
		metrics["gst_status"] = "Partial - " + fmt.Sprintf("%d/%d returns filed", gstFiled, gstCount)
	}

	// Query TDS collections
	query = `
		SELECT COALESCE(SUM(tds_amount), 0) as total_tds,
		       COUNT(*) as return_count
		FROM tds_compliance
		WHERE tenant_id = ? AND deleted_at IS NULL
	`

	var tdsTotalAmount float64
	var tdsReturnCount int
	err = s.DB.QueryRow(query, tenantID).Scan(&tdsTotalAmount, &tdsReturnCount)
	if err != nil && err != sql.ErrNoRows {
		return metrics, fmt.Errorf("failed to query TDS metrics: %w", err)
	}
	metrics["tds_collected"] = tdsTotalAmount

	// Query Advance Tax
	query = `
		SELECT COALESCE(SUM(amount_paid), 0) as total_paid
		FROM advance_tax_compliance
		WHERE tenant_id = ? AND deleted_at IS NULL AND payment_status = 'Paid'
	`

	var advanceTaxPaid float64
	err = s.DB.QueryRow(query, tenantID).Scan(&advanceTaxPaid)
	if err != nil && err != sql.ErrNoRows {
		return metrics, fmt.Errorf("failed to query advance tax metrics: %w", err)
	}
	metrics["advance_tax_paid"] = advanceTaxPaid

	return metrics, nil
}
