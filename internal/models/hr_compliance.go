package models

import "time"

// ============================================================================
// HR COMPLIANCE MODELS - LABOUR LAWS
// ============================================================================

// HRComplianceRule represents configurable compliance rules
type HRComplianceRule struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	RuleName        string `json:"rule_name"`
	RuleType        string `json:"rule_type"` // ESI_Threshold, EPF_Threshold, Bonus_Calculation, etc.
	RuleDescription string `json:"rule_description"`

	// Applicable Period
	ApplicableFrom time.Time  `json:"applicable_from"`
	ApplicableTo   *time.Time `json:"applicable_to"`

	// Rule Parameters (JSON)
	RuleParameters string `json:"rule_parameters"` // JSON string

	IsActive bool `json:"is_active"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedBy *string    `json:"created_by"`
}

// ESICompliance represents ESI (Employee State Insurance) compliance record
type ESICompliance struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	EmployeeID string `json:"employee_id"`

	// ESI Registration
	ESINumber           string     `json:"esi_number"`
	ESIRegistrationDate *time.Time `json:"esi_registration_date"`
	ESIOffice           string     `json:"esi_office"`
	ESIState            string     `json:"esi_state"`

	// Contribution Details
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`

	// Contribution Rates
	EmployeeContributionRate float64 `json:"employee_contribution_rate"` // 0.75%
	EmployerContributionRate float64 `json:"employer_contribution_rate"` // 3.25%
	WageLimit                float64 `json:"wage_limit"`                 // Usually ₹21,000

	// Benefits
	SickLeaveBalance         int  `json:"sick_leave_balance"`
	DisabilityBenefitAvailed bool `json:"disability_benefit_availed"`
	MedicalBenefitAvailed    bool `json:"medical_benefit_availed"`

	// Compliance Status
	IsCompliant         bool   `json:"is_compliant"`
	NonComplianceReason string `json:"non_compliance_reason"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// EPFCompliance represents EPF/PF (Employee Provident Fund) compliance
type EPFCompliance struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	EmployeeID string `json:"employee_id"`

	// PF Account Details
	PFAccountNumber    string     `json:"pf_account_number"`
	PFRegistrationDate *time.Time `json:"pf_registration_date"`
	PFOffice           string     `json:"pf_office"`
	PFState            string     `json:"pf_state"`

	// Contribution Rates
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`

	EmployeeContributionRate float64 `json:"employee_contribution_rate"` // 12%
	EmployerContributionRate float64 `json:"employer_contribution_rate"` // 12%
	VPSContributionRate      float64 `json:"vps_contribution_rate"`

	// Exemption Status
	PFExempt     bool   `json:"pf_exempt"`
	ExemptReason string `json:"exempt_reason"`

	// Account Details
	AccumulatedBalance float64 `json:"accumulated_balance"`
	CurrentBalance     float64 `json:"current_balance"`

	// Withdrawal & Settlement
	PartialWithdrawalAllowed bool       `json:"partial_withdrawal_allowed"`
	FinalSettlementRequested bool       `json:"final_settlement_requested"`
	SettlementDate           *time.Time `json:"settlement_date"`
	SettlementAmount         float64    `json:"settlement_amount"`

	// Compliance
	IsCompliant         bool   `json:"is_compliant"`
	NonComplianceReason string `json:"non_compliance_reason"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// ProfessionalTaxCompliance represents State-wise Professional Tax compliance
type ProfessionalTaxCompliance struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	EmployeeID string `json:"employee_id"`

	// PT Registration
	PTNumber            string     `json:"pt_number"`
	PTRegistrationState string     `json:"pt_registration_state"`
	PTRegistrationDate  *time.Time `json:"pt_registration_date"`

	// PT Slab
	MonthlySalaryThreshold float64    `json:"monthly_salary_threshold"`
	PTAmount               float64    `json:"pt_amount"`
	ApplicableFrom         time.Time  `json:"applicable_from"`
	ApplicableTo           *time.Time `json:"applicable_to"`

	IsCompliant bool `json:"is_compliant"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// GratuityCompliance represents Gratuity Act compliance (1972)
// Gratuity = 15 days' salary for first 5 years + 30 days for subsequent years
type GratuityCompliance struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	EmployeeID string `json:"employee_id"`

	// Eligibility
	JoiningDate      time.Time `json:"joining_date"`
	GratuityEligible bool      `json:"gratuity_eligible"` // Only if service > 5 years
	YearsOfService   int       `json:"years_of_service"`

	// Gratuity Calculation
	LastDrawnSalary   float64 `json:"last_drawn_salary"`
	GratuityAccrued   float64 `json:"gratuity_accrued"`
	GratuityLiability float64 `json:"gratuity_liability"`

	// Payment
	GratuityPaid       bool       `json:"gratuity_paid"`
	GratuityPaidDate   *time.Time `json:"gratuity_paid_date"`
	GratuityPaidAmount float64    `json:"gratuity_paid_amount"`

	// Fund/Insurance
	FundType     string  `json:"fund_type"` // Gratuity_Fund, Insurance_Policy, Direct_Payment
	PolicyNumber string  `json:"policy_number"`
	FundBalance  float64 `json:"fund_balance"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// BonusCompliance represents Bonus & Variable Pay compliance
type BonusCompliance struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	EmployeeID string `json:"employee_id"`

	// Bonus Configuration
	BonusType string `json:"bonus_type"` // Festival_Bonus, Annual_Bonus, Performance_Bonus, Diwali_Bonus
	BonusName string `json:"bonus_name"`

	// Eligibility (30 days continuous service)
	Eligible            bool `json:"eligible"`
	DaysWorked          int  `json:"days_worked"`
	MinimumDaysRequired int  `json:"minimum_days_required"`

	// Calculation
	ApplicableSalaryComponent string  `json:"applicable_salary_component"` // Basic, Basic_DA, CTC, Gross
	BonusPercentage           float64 `json:"bonus_percentage"`            // e.g., 8.33% for 1 month
	BonusAmount               float64 `json:"bonus_amount"`

	// Payment
	BonusYear     int        `json:"bonus_year"`
	DueDate       *time.Time `json:"due_date"`
	PaymentDate   *time.Time `json:"payment_date"`
	PaymentStatus string     `json:"payment_status"` // Not_Applicable, Pending, Paid, Waived

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// LeaveCompliance represents leave entitlement and compliance
type LeaveCompliance struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	EmployeeID  string `json:"employee_id"`
	LeaveTypeID string `json:"leave_type_id"`

	// Fiscal Year
	FiscalYearFrom time.Time `json:"fiscal_year_from"`
	FiscalYearTo   time.Time `json:"fiscal_year_to"`

	// Entitlement
	AnnualEntitlement int `json:"annual_entitlement"`

	// Opening Balance (carry forward)
	OpeningBalance    int `json:"opening_balance"`
	CarryForwardLimit int `json:"carry_forward_limit"`

	// Current Year Usage
	Utilized  int `json:"utilized"`
	Available int `json:"available"`

	// Encashment
	EncashmentAllowed bool    `json:"encashment_allowed"`
	EncashmentAmount  float64 `json:"encashment_amount"`

	IsCompliant bool `json:"is_compliant"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// StatutoryComplianceAudit tracks compliance checks and violations
type StatutoryComplianceAudit struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	EmployeeID string `json:"employee_id"`

	// Compliance Check
	ComplianceType   string     `json:"compliance_type"` // ESI, EPF, PT, Leave, Gratuity, etc.
	ComplianceItem   string     `json:"compliance_item"`
	CompliancePeriod *time.Time `json:"compliance_period"`

	// Check Result
	IsCompliant    bool   `json:"is_compliant"`
	ViolationFound string `json:"violation_found"`
	Severity       string `json:"severity"` // Critical, High, Medium, Low

	// Action
	ActionRequired string     `json:"action_required"`
	ActionTaken    string     `json:"action_taken"`
	ActionDate     *time.Time `json:"action_date"`
	ActionTakenBy  *string    `json:"action_taken_by"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// LabourComplianceDocument represents compliance-related documents
type LabourComplianceDocument struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	EmployeeID string `json:"employee_id"`

	// Document Details
	DocumentType   string     `json:"document_type"` // Offer_Letter, Appointment_Letter, ESI_Certificate, etc.
	DocumentName   string     `json:"document_name"`
	DocumentDate   *time.Time `json:"document_date"`
	DocumentNumber string     `json:"document_number"`

	// Validity
	IssueDate    *time.Time `json:"issue_date"`
	ExpiryDate   *time.Time `json:"expiry_date"`
	DocumentPath string     `json:"document_path"` // Storage path

	// Verification
	VerifiedBy         *string    `json:"verified_by"`
	VerifiedDate       *time.Time `json:"verified_date"`
	VerificationStatus string     `json:"verification_status"` // Not_Verified, Verified, Expired, Invalid

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// WorkingHoursCompliance represents overtime and working hours compliance
type WorkingHoursCompliance struct {
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	EmployeeID string `json:"employee_id"`

	// Period
	FiscalYearFrom time.Time `json:"fiscal_year_from"`
	FiscalYearTo   time.Time `json:"fiscal_year_to"`

	// Working Hours Configuration
	NormalWorkingHoursPerDay float64 `json:"normal_working_hours_per_day"` // 8 hours
	WorkingDaysPerWeek       int     `json:"working_days_per_week"`        // 5 or 6

	// Overtime Policy
	OvertimeAllowed          bool    `json:"overtime_allowed"`
	OvertimeRateMultiplier   float64 `json:"overtime_rate_multiplier"` // 1.5x or 2x
	MaxOvertimeHoursPerMonth int     `json:"max_overtime_hours_per_month"`

	// Actual Hours
	TotalHoursWorked    float64 `json:"total_hours_worked"`
	TotalOvertimeHours  float64 `json:"total_overtime_hours"`
	TotalUndertimeHours float64 `json:"total_undtime_hours"`

	// Compliance
	IsCompliant         bool   `json:"is_compliant"`
	NonComplianceReason string `json:"non_compliance_reason"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// InternalComplaint represents POSH Act complaints
type InternalComplaint struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	// Complaint
	ComplaintNumber string    `json:"complaint_number"`
	ComplaintDate   time.Time `json:"complaint_date"`

	// Parties
	ComplainantID string `json:"complainant_id"`
	AccusedID     string `json:"accused_id"`

	// Details
	ComplaintType string `json:"complaint_type"` // Sexual_Harassment, Discrimination, Harassment, Other
	Description   string `json:"description"`

	// Investigation
	InvestigationInitiated      bool       `json:"investigation_initiated"`
	InvestigationStartDate      *time.Time `json:"investigation_start_date"`
	InvestigationCompletionDate *time.Time `json:"investigation_completion_date"`
	InvestigationReportPath     string     `json:"investigation_report_path"`

	// Resolution
	Status           string     `json:"status"` // Open, Under_Investigation, Resolved, Closed, Escalated
	ResolutionDate   *time.Time `json:"resolution_date"`
	ResolutionAction string     `json:"resolution_action"`
	ActionTakenBy    *string    `json:"action_taken_by"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// StatutoryFormsFiling represents statutory form submissions
type StatutoryFormsFiling struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`

	// Form Details
	FormName string `json:"form_name"` // Form 5, Form 1, Form 12AA, etc.
	FormType string `json:"form_type"` // EPF_Form, ESI_Form, Income_Tax_Form, Labour_Form

	// Period
	ApplicableFrom time.Time `json:"applicable_from"`
	ApplicableTo   time.Time `json:"applicable_to"`

	// Filing
	FilingDate            *time.Time `json:"filing_date"`
	FilingReferenceNumber string     `json:"filing_reference_number"`
	FiledBy               *string    `json:"filed_by"`

	// Status
	SubmissionStatus string     `json:"submission_status"` // Not_Due, Draft, Submitted, Accepted, Rejected, Pending_Response
	SubmissionDate   *time.Time `json:"submission_date"`

	// Documents
	FormDocumentPath            string `json:"form_document_path"`
	AcknowledgementDocumentPath string `json:"acknowledgement_document_path"`

	Remarks string `json:"remarks"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// ============================================================================
// REQUEST/RESPONSE MODELS
// ============================================================================

// ComplianceCheckResult represents the result of a compliance check
type ComplianceCheckResult struct {
	EmployeeID   string `json:"employee_id"`
	EmployeeName string `json:"employee_name"`

	ESICompliance      *ComplianceStatus `json:"esi_compliance"`
	EPFCompliance      *ComplianceStatus `json:"epf_compliance"`
	PTCompliance       *ComplianceStatus `json:"pt_compliance"`
	GratuityCompliance *ComplianceStatus `json:"gratuity_compliance"`
	LeaveCompliance    *ComplianceStatus `json:"leave_compliance"`

	OverallStatus      string   `json:"overall_status"` // Compliant, Non-Compliant, Action_Required
	CriticalViolations int      `json:"critical_violations"`
	Recommendations    []string `json:"recommendations"`
}

// ComplianceStatus represents status of a single compliance check
type ComplianceStatus struct {
	IsCompliant     bool       `json:"is_compliant"`
	LastCheckedDate *time.Time `json:"last_checked_date"`
	Violations      []string   `json:"violations"`
	ActionRequired  string     `json:"action_required"`
	DueDate         *time.Time `json:"due_date"`
}
