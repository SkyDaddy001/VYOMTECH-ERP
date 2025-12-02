package models

import "time"

// ============================================================================
// TAX COMPLIANCE MODELS - INCOME TAX & GST
// ============================================================================

// TaxConfiguration represents the overall tax setup for an organization
type TaxConfiguration struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	// Financial Year
	FinancialYearStart time.Time `json:"financial_year_start"`
	FinancialYearEnd   time.Time `json:"financial_year_end"`
	FiscalYear         int       `json:"fiscal_year"`

	// Income Tax Configuration
	ITAssessmentYear     int    `json:"it_assessment_year"`
	ITRegistrationNumber string `json:"it_registration_number"`
	ITPAN                string `json:"it_pan"`
	ITCircle             string `json:"it_circle"`
	ITAssessingOfficer   string `json:"it_assessing_officer"`

	// GST Configuration
	GSTRegistrationNumber    string     `json:"gst_registration_number"`
	GSTRegistrationDate      *time.Time `json:"gst_registration_date"`
	GSTStatus                string     `json:"gst_status"`                  // Active, Inactive, Cancelled, Suspended
	GSTReturnFilingFrequency string     `json:"gst_return_filing_frequency"` // Monthly, Quarterly, Annual

	// Filing Status
	ITFilingStatus  string `json:"it_filing_status"`
	GSTFilingStatus string `json:"gst_filing_status"`

	// Compliance Contact
	ComplianceOfficerName  string `json:"compliance_officer_name"`
	ComplianceOfficerEmail string `json:"compliance_officer_email"`
	ComplianceOfficerPhone string `json:"compliance_officer_phone"`

	// Audit Info
	IsActive      bool       `json:"is_active"`
	LastAuditDate *time.Time `json:"last_audit_date"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// IncomeTexCompliance represents annual income tax compliance record
type IncomeTaxCompliance struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	// Financial Year
	FiscalYear     int `json:"fiscal_year"`
	AssessmentYear int `json:"assessment_year"`

	// Income Calculation
	TotalGrossIncome float64 `json:"total_gross_income"`
	TotalDeductions  float64 `json:"total_deductions"`
	TaxableIncome    float64 `json:"taxable_income"`

	// Income Breakup
	SalaryIncome       float64 `json:"salary_income"`
	BusinessIncome     float64 `json:"business_income"`
	CapitalGainsIncome float64 `json:"capital_gains_income"`
	RentalIncome       float64 `json:"rental_income"`
	OtherIncome        float64 `json:"other_income"`

	// Deductions (Section 80)
	Section80CDeduction   float64 `json:"section_80c_deduction"`
	Section80DDeduction   float64 `json:"section_80d_deduction"`
	Section80EDeduction   float64 `json:"section_80e_deduction"`
	Section80GDeduction   float64 `json:"section_80g_deduction"`
	Section80GGCDeduction float64 `json:"section_80ggc_deduction"`
	OtherDeductions       float64 `json:"other_deductions"`

	// Tax Calculation
	TaxBeforeSurcharge float64 `json:"tax_before_surcharge"`
	SurchargeAmount    float64 `json:"surcharge_amount"`
	CessAmount         float64 `json:"cess_amount"`
	TotalTaxLiability  float64 `json:"total_tax_liability"`

	// Tax Payment
	AdvanceTaxPaid       float64 `json:"advance_tax_paid"`
	TDSCredit            float64 `json:"tds_credit"`
	TotalTaxPaid         float64 `json:"total_tax_paid"`
	TaxPayableRefundable float64 `json:"tax_payable_refundable"`

	// ITR Filing
	ITRFormType               string     `json:"itr_form_type"` // ITR-1 to ITR-7
	ReturnFilingDate          *time.Time `json:"return_filing_date"`
	ReturnFiledStatus         string     `json:"return_filed_status"`
	ACKNumber                 string     `json:"ack_number"`
	FilingAcknowledgementDate *time.Time `json:"filing_acknowledgement_date"`

	// Scrutiny/Assessment
	ScrutinyInitiated     bool       `json:"scrutiny_initiated"`
	ScrutinyDate          *time.Time `json:"scrutiny_date"`
	ScrutinyResponseDate  *time.Time `json:"scrutiny_response_date"`
	AssessmentOrderDate   *time.Time `json:"assessment_order_date"`
	AdditionalTaxAssessed float64    `json:"additional_tax_assessed"`

	// Compliance
	IsCompliant         bool   `json:"is_compliant"`
	NonComplianceReason string `json:"non_compliance_reason"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// TDSCompliance represents TDS (Tax Deducted at Source) compliance
type TDSCompliance struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	FiscalYear int `json:"fiscal_year"`

	// TDS Category
	TDSSection  string `json:"tds_section"` // 192, 193, 194A, 194C, etc.
	TDSCategory string `json:"tds_category"`

	// TDS Calculation
	TotalPaymentsMade float64 `json:"total_payments_made"`
	TDSRate           float64 `json:"tds_rate"`
	TDSAmountDeducted float64 `json:"tds_amount_deducted"`

	// TDS Deposit
	TDSDepositDate *time.Time `json:"tds_deposit_date"`
	ChallanNumber  string     `json:"challan_number"`
	BankName       string     `json:"bank_name"`

	// TDS Return Filing
	TDSReturnPeriod     string     `json:"tds_return_period"` // Q1, Q2, Q3, Q4
	TDSReturnFilingDate *time.Time `json:"tds_return_filing_date"`
	TDSReturnStatus     string     `json:"tds_return_status"`

	// Reconciliation
	PayeeDetailsReconciled bool       `json:"payee_details_reconciled"`
	ReconciliationDate     *time.Time `json:"reconciliation_date"`
	ReconciliationVariance float64    `json:"reconciliation_variance"`

	// Annual Summary
	AnnualTDSSummaryFiled   bool       `json:"annual_tds_summary_filed"`
	AnnualSummaryFilingDate *time.Time `json:"annual_summary_filing_date"`

	IsCompliant bool `json:"is_compliant"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// GSTCompliance represents GST compliance for a return period
type GSTCompliance struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	FiscalYear   int       `json:"fiscal_year"`
	ReturnPeriod string    `json:"return_period"` // Monthly, Quarterly
	MonthYear    time.Time `json:"month_year"`

	// Outward Supplies
	TotalOutwardSupplies float64 `json:"total_outward_supplies"`

	GST5Sales  float64 `json:"gst_5_sales"`
	GST12Sales float64 `json:"gst_12_sales"`
	GST18Sales float64 `json:"gst_18_sales"`
	GST28Sales float64 `json:"gst_28_sales"`

	IntraStateSales float64 `json:"intra_state_sales"`
	InterStateSales float64 `json:"inter_state_sales"`
	Exports         float64 `json:"exports"`

	// Output GST
	OutputGST5     float64 `json:"output_gst_5"`
	OutputGST12    float64 `json:"output_gst_12"`
	OutputGST18    float64 `json:"output_gst_18"`
	OutputGST28    float64 `json:"output_gst_28"`
	OutputGSTCess  float64 `json:"output_gst_cess"`
	TotalOutputGST float64 `json:"total_output_gst"`

	// Inward Supplies
	TotalInwardSupplies float64 `json:"total_inward_supplies"`

	// Input GST
	InputGST5     float64 `json:"input_gst_5"`
	InputGST12    float64 `json:"input_gst_12"`
	InputGST18    float64 `json:"input_gst_18"`
	InputGST28    float64 `json:"input_gst_28"`
	InputGSTCess  float64 `json:"input_gst_cess"`
	TotalInputGST float64 `json:"total_input_gst"`

	// Interest Free Credit
	InterestFreeCreditCarried float64 `json:"interest_free_credit_carried"`

	// Net GST
	NetGSTPay    float64 `json:"net_gst_payable"`
	RefundClaim  float64 `json:"refund_claim"`
	RefundStatus string  `json:"refund_status"`

	AdvanceTaxPaid float64 `json:"advance_tax_paid"`

	// GSTR Filing
	GSTR1Filed  bool `json:"gstr_1_filed"`
	GSTR2Filed  bool `json:"gstr_2_filed"`
	GSTR3Filed  bool `json:"gstr_3_filed"`
	GSTR4Filed  bool `json:"gstr_4_filed"`
	GSTR5Filed  bool `json:"gstr_5_filed"`
	GSTR6Filed  bool `json:"gstr_6_filed"`
	GSTR7Filed  bool `json:"gstr_7_filed"`
	GSTR8Filed  bool `json:"gstr_8_filed"`
	GSTR9Filed  bool `json:"gstr_9_filed"`
	GSTR10Filed bool `json:"gstr_10_filed"`

	GSTRFilingDate   *time.Time `json:"gstr_filing_date"`
	GSTRFilingStatus string     `json:"gstr_filing_status"`

	// Reconciliation
	GSTR1GSTR3Reconciled bool   `json:"gstr_1_gstr_3_reconciled"`
	ReconciliationNotes  string `json:"reconciliation_notes"`

	IsCompliant         bool   `json:"is_compliant"`
	NonComplianceReason string `json:"non_compliance_reason"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// GSTInvoiceTracking tracks GST on individual invoices
type GSTInvoiceTracking struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	InvoiceID     string     `json:"invoice_id"`
	InvoiceNumber string     `json:"invoice_number"`
	InvoiceDate   *time.Time `json:"invoice_date"`
	InvoiceAmount float64    `json:"invoice_amount"`

	CustomerID    string `json:"customer_id"`
	CustomerGSTIN string `json:"customer_gstin"`

	GSTRate   float64 `json:"gst_rate"`
	GSTAmount float64 `json:"gst_amount"`

	InvoiceRaisedDate *time.Time `json:"invoice_raised_date"`
	InvoiceCancelled  bool       `json:"invoice_cancelled"`
	CancellationDate  *time.Time `json:"cancellation_date"`

	GSTR1Reported    bool `json:"gstr_1_reported"`
	GSTR1FilingMonth int  `json:"gstr_1_filing_month"`
	GSTR1FilingYear  int  `json:"gstr_1_filing_year"`

	ITCEligible bool `json:"itc_eligible"`
	ITCClaimed  bool `json:"itc_claimed"`

	Reconciled         bool       `json:"reconciled"`
	ReconciliationDate *time.Time `json:"reconciliation_date"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// GSTInputCredit tracks input GST on purchases
type GSTInputCredit struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	PurchaseInvoiceID string `json:"purchase_invoice_id"`
	VendorID          string `json:"vendor_id"`
	VendorGSTIN       string `json:"vendor_gstin"`

	InvoiceNumber string     `json:"invoice_number"`
	InvoiceDate   *time.Time `json:"invoice_date"`
	InvoiceAmount float64    `json:"invoice_amount"`

	GSTRate   float64 `json:"gst_rate"`
	GSTAmount float64 `json:"gst_amount"`

	ITCEligible         bool   `json:"itc_eligible"`
	ITCIneligibleReason string `json:"itc_ineligible_reason"`

	ITCClaimed    bool `json:"itc_claimed"`
	ITCClaimMonth int  `json:"itc_claim_month"`
	ITCClaimYear  int  `json:"itc_claim_year"`

	BlockedPercentage float64 `json:"blocked_percentage"`
	BlockedAmount     float64 `json:"blocked_amount"`

	GSTR2Reported    bool `json:"gstr_2_reported"`
	GSTR2FilingMonth int  `json:"gstr_2_filing_month"`
	GSTR2FilingYear  int  `json:"gstr_2_filing_year"`

	VendorGSTR1Reconciled bool       `json:"vendor_gstr_1_reconciled"`
	ReconciliationDate    *time.Time `json:"reconciliation_date"`
	DiscrepancyFound      bool       `json:"discrepancy_found"`
	DiscrepancyNotes      string     `json:"discrepancy_notes"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// AdvanceTaxSchedule tracks quarterly advance tax payments
type AdvanceTaxSchedule struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	FiscalYear     int `json:"fiscal_year"`
	AssessmentYear int `json:"assessment_year"`

	// Q1 (June 15)
	Q1DueDate       *time.Time `json:"q1_due_date"`
	Q1AmountDue     float64    `json:"q1_amount_due"`
	Q1AmountPaid    float64    `json:"q1_amount_paid"`
	Q1PaymentDate   *time.Time `json:"q1_payment_date"`
	Q1ChallanNumber string     `json:"q1_challan_number"`
	Q1Status        string     `json:"q1_status"`

	// Q2 (September 15)
	Q2DueDate       *time.Time `json:"q2_due_date"`
	Q2AmountDue     float64    `json:"q2_amount_due"`
	Q2AmountPaid    float64    `json:"q2_amount_paid"`
	Q2PaymentDate   *time.Time `json:"q2_payment_date"`
	Q2ChallanNumber string     `json:"q2_challan_number"`
	Q2Status        string     `json:"q2_status"`

	// Q3 (December 15)
	Q3DueDate       *time.Time `json:"q3_due_date"`
	Q3AmountDue     float64    `json:"q3_amount_due"`
	Q3AmountPaid    float64    `json:"q3_amount_paid"`
	Q3PaymentDate   *time.Time `json:"q3_payment_date"`
	Q3ChallanNumber string     `json:"q3_challan_number"`
	Q3Status        string     `json:"q3_status"`

	// Q4 (March 15)
	Q4DueDate       *time.Time `json:"q4_due_date"`
	Q4AmountDue     float64    `json:"q4_amount_due"`
	Q4AmountPaid    float64    `json:"q4_amount_paid"`
	Q4PaymentDate   *time.Time `json:"q4_payment_date"`
	Q4ChallanNumber string     `json:"q4_challan_number"`
	Q4Status        string     `json:"q4_status"`

	// Summary
	TotalAdvanceTaxDue  float64 `json:"total_advance_tax_due"`
	TotalAdvanceTaxPaid float64 `json:"total_advance_tax_paid"`

	IsCompliant bool `json:"is_compliant"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// TaxAuditTrail represents an audit of tax records
type TaxAuditTrail struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	AuditType   string     `json:"audit_type"` // Income_Tax, GST, TDS, Statutory, Internal
	AuditDate   *time.Time `json:"audit_date"`
	AuditPeriod string     `json:"audit_period"`

	AuditorName          string `json:"auditor_name"`
	AuditorFirm          string `json:"auditor_firm"`
	AuditReportReference string `json:"audit_report_reference"`

	AuditReportPath  string  `json:"audit_report_path"`
	TotalAdjustments float64 `json:"total_adjustments"`
	Recommendations  string  `json:"recommendations"`

	AuditStatus      string     `json:"audit_status"`
	FollowUpRequired bool       `json:"follow_up_required"`
	FollowUpDate     *time.Time `json:"follow_up_date"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// TaxComplianceDocument represents a tax document
type TaxComplianceDocument struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	DocumentType    string     `json:"document_type"`
	DocumentName    string     `json:"document_name"`
	DocumentDate    *time.Time `json:"document_date"`
	ReferenceNumber string     `json:"reference_number"`

	FiscalYear   int    `json:"fiscal_year"`
	FilingPeriod string `json:"filing_period"`

	DocumentPath string `json:"document_path"`
	FileSize     int    `json:"file_size"`
	FileHash     string `json:"file_hash"`

	Verified     bool       `json:"verified"`
	VerifiedBy   *string    `json:"verified_by"`
	VerifiedDate *time.Time `json:"verified_date"`

	Archived    bool       `json:"archived"`
	ArchiveDate *time.Time `json:"archive_date"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// TaxComplianceChecklist represents compliance checklist items
type TaxComplianceChecklist struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	ComplianceItem string `json:"compliance_item"`
	ComplianceType string `json:"compliance_type"` // Income_Tax, GST, TDS, Other

	DueDate                 time.Time  `json:"due_date"`
	EstimatedCompletionDate *time.Time `json:"estimated_completion_date"`
	ActualCompletionDate    *time.Time `json:"actual_completion_date"`

	Status string `json:"status"` // Not_Started, In_Progress, Completed, Overdue, N/A

	AssignedTo *string `json:"assigned_to"`

	Description       string `json:"description"`
	DocumentsRequired string `json:"documents_required"`
	Remarks           string `json:"remarks"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// ============================================================================
// REQUEST/RESPONSE MODELS
// ============================================================================

// TaxComplianceStatusResponse represents overall tax compliance status
type TaxComplianceStatusResponse struct {
	TenantID   string `json:"tenant_id"`
	FiscalYear int    `json:"fiscal_year"`

	IncomeTaxStatus  *IncomeTaxComplianceStatus `json:"income_tax_status"`
	GSTStatus        *GSTComplianceStatus       `json:"gst_status"`
	TDSStatus        *TDSComplianceStatus       `json:"tds_status"`
	AdvanceTaxStatus *AdvanceTaxStatus          `json:"advance_tax_status"`

	OverallCompliance string   `json:"overall_compliance"` // Compliant, Non-Compliant, Action_Required
	PendingActions    []string `json:"pending_actions"`
	UpcomingDeadlines []string `json:"upcoming_deadlines"`
}

// IncomeTaxComplianceStatus represents IT filing status
type IncomeTaxComplianceStatus struct {
	ITRFiled     bool       `json:"itr_filed"`
	FilingDate   *time.Time `json:"filing_date"`
	Status       string     `json:"status"`
	TaxLiability float64    `json:"tax_liability"`
	RefundAmount float64    `json:"refund_amount"`
}

// GSTComplianceStatus represents GST filing status
type GSTComplianceStatus struct {
	GSTRFiled    bool       `json:"gstr_filed"`
	FilingDate   *time.Time `json:"filing_date"`
	Status       string     `json:"status"`
	TotalTax     float64    `json:"total_tax"`
	RefundStatus string     `json:"refund_status"`
}

// TDSComplianceStatus represents TDS filing status
type TDSComplianceStatus struct {
	TDSDeducted    float64 `json:"tds_deducted"`
	TDSDeposited   float64 `json:"tds_deposited"`
	TDSReturnFiled bool    `json:"tds_return_filed"`
	Status         string  `json:"status"`
}

// AdvanceTaxStatus represents advance tax payment status
type AdvanceTaxStatus struct {
	Q1Status  string  `json:"q1_status"`
	Q2Status  string  `json:"q2_status"`
	Q3Status  string  `json:"q3_status"`
	Q4Status  string  `json:"q4_status"`
	TotalPaid float64 `json:"total_paid"`
	TotalDue  float64 `json:"total_due"`
}
