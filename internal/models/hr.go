package models

import "time"

// ============================================================================
// HR & PAYROLL MODELS
// ============================================================================

// Employee represents an employee in the HR system
type Employee struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Country     string    `json:"country"`
	PostalCode  string    `json:"postal_code"`

	// Employment Details
	EmployeeID     string     `json:"employee_id"`
	Designation    string     `json:"designation"`
	Department     string     `json:"department"`
	ReportTo       string     `json:"report_to"`
	EmploymentType string     `json:"employment_type"`
	JoiningDate    time.Time  `json:"joining_date"`
	ExitDate       *time.Time `json:"exit_date"`
	Status         string     `json:"status"`

	// Bank Details
	BankAccountNumber string `json:"bank_account_number"`
	BankIFSCCode      string `json:"bank_ifsc_code"`
	BankName          string `json:"bank_name"`
	AccountHolderName string `json:"account_holder_name"`

	// Salary Structure
	BaseSalary          float64 `json:"base_salary"`
	DAAllowance         float64 `json:"dearness_allowance"`
	HRAAllowance        float64 `json:"house_rent_allowance"`
	SpecialAllowance    float64 `json:"special_allowance"`
	ConveyanceAllowance float64 `json:"conveyance_allowance"`
	MedicalAllowance    float64 `json:"medical_allowance"`
	OtherAllowances     float64 `json:"other_allowances"`

	// Deductions
	EPFDeduction     float64 `json:"epf_deduction"`
	ESIDeduction     float64 `json:"esi_deduction"`
	ProfessionalTax  float64 `json:"professional_tax"`
	IncomeTax        float64 `json:"income_tax"`
	LoanDeduction    float64 `json:"loan_deduction"`
	AdvanceDeduction float64 `json:"advance_deduction"`
	OtherDeductions  float64 `json:"other_deductions"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Attendance tracks employee attendance
type Attendance struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	EmployeeID     string     `json:"employee_id"`
	AttendanceDate time.Time  `json:"attendance_date"`
	CheckInTime    *time.Time `json:"check_in_time"`
	CheckOutTime   *time.Time `json:"check_out_time"`
	WorkingHours   float64    `json:"working_hours"`
	Status         string     `json:"status"`
	Notes          string     `json:"notes"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// LeaveType represents different types of leaves
type LeaveType struct {
	ID                string `json:"id"`
	TenantID          string `json:"tenant_id"`
	LeaveTypeName     string `json:"leave_type_name"`
	Description       string `json:"description"`
	AnnualEntitlement int    `json:"annual_entitlement"`
	IsPaid            bool   `json:"is_paid"`
	CarryForwardLimit int    `json:"carry_forward_limit"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// LeaveRequest represents a leave request from an employee
type LeaveRequest struct {
	ID              string     `json:"id"`
	TenantID        string     `json:"tenant_id"`
	EmployeeID      string     `json:"employee_id"`
	LeaveTypeID     string     `json:"leave_type_id"`
	FromDate        time.Time  `json:"from_date"`
	ToDate          time.Time  `json:"to_date"`
	NumberOfDays    int        `json:"number_of_days"`
	Reason          string     `json:"reason"`
	Status          string     `json:"status"`
	ApprovedBy      string     `json:"approved_by"`
	ApprovalDate    *time.Time `json:"approval_date"`
	RejectionReason string     `json:"rejection_reason"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// PayrollRecord represents a monthly payroll record for an employee
type PayrollRecord struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	EmployeeID    string    `json:"employee_id"`
	PayrollMonth  time.Time `json:"payroll_month"`
	PayrollStatus string    `json:"payroll_status"`

	// Earnings
	BasicSalary      float64 `json:"basic_salary"`
	DAAllowance      float64 `json:"dearness_allowance"`
	HRAAllowance     float64 `json:"house_rent_allowance"`
	SpecialAllowance float64 `json:"special_allowance"`
	ConveyanceAllow  float64 `json:"conveyance_allowance"`
	MedicalAllow     float64 `json:"medical_allowance"`
	OtherAllowances  float64 `json:"other_allowances"`
	TotalEarnings    float64 `json:"total_earnings"`

	// Deductions
	EPFDeduction     float64 `json:"epf_deduction"`
	ESIDeduction     float64 `json:"esi_deduction"`
	ProfessionalTax  float64 `json:"professional_tax"`
	IncomeTax        float64 `json:"income_tax"`
	LoanDeduction    float64 `json:"loan_deduction"`
	AdvanceDeduction float64 `json:"advance_deduction"`
	OtherDeductions  float64 `json:"other_deductions"`
	TotalDeductions  float64 `json:"total_deductions"`

	// Net Salary
	NetSalary float64 `json:"net_salary"`

	// Additional Info
	WorkingDays int     `json:"working_days"`
	LeaveDays   int     `json:"leave_days"`
	PaidDays    float64 `json:"paid_days"`
	Notes       string  `json:"notes"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ProcessedAt *time.Time `json:"processed_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// EmployeeLoan represents a loan given to an employee
type EmployeeLoan struct {
	ID               string     `json:"id"`
	TenantID         string     `json:"tenant_id"`
	EmployeeID       string     `json:"employee_id"`
	LoanAmount       float64    `json:"loan_amount"`
	LoanType         string     `json:"loan_type"`
	InterestRate     float64    `json:"interest_rate"`
	TenureMonths     int        `json:"tenure_months"`
	EMIAmount        float64    `json:"emi_amount"`
	LoanStatus       string     `json:"loan_status"`
	DisbursementDate *time.Time `json:"disbursement_date"`
	ClosureDate      *time.Time `json:"closure_date"`
	RemainingAmount  float64    `json:"remaining_amount"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// LoanRepayment represents a loan repayment installment
type LoanRepayment struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	LoanID            string     `json:"loan_id"`
	InstallmentNumber int        `json:"installment_number"`
	EMIAmount         float64    `json:"emi_amount"`
	PrincipalAmount   float64    `json:"principal_amount"`
	InterestAmount    float64    `json:"interest_amount"`
	DueDate           time.Time  `json:"due_date"`
	PaidDate          *time.Time `json:"paid_date"`
	PaymentStatus     string     `json:"payment_status"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// PerformanceAppraisal represents an employee's performance review
type PerformanceAppraisal struct {
	ID                  string    `json:"id"`
	TenantID            string    `json:"tenant_id"`
	EmployeeID          string    `json:"employee_id"`
	AppraisalPeriodFrom time.Time `json:"appraisal_period_from"`
	AppraisalPeriodTo   time.Time `json:"appraisal_period_to"`
	AppraiserID         string    `json:"appraiser_id"`
	OverallRating       float64   `json:"overall_rating"`
	CommunicationRating float64   `json:"communication_rating"`
	PerformanceRating   float64   `json:"performance_rating"`
	AttendanceRating    float64   `json:"attendance_rating"`
	TeamworkRating      float64   `json:"teamwork_rating"`
	InitiativeRating    float64   `json:"initiative_rating"`
	Comments            string    `json:"comments"`
	AppraisalStatus     string    `json:"appraisal_status"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// EmployeeBenefit represents a benefit provided to an employee
type EmployeeBenefit struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	EmployeeID    string     `json:"employee_id"`
	BenefitType   string     `json:"benefit_type"`
	Provider      string     `json:"provider"`
	PolicyNumber  string     `json:"policy_number"`
	SumInsured    float64    `json:"sum_insured"`
	EffectiveDate time.Time  `json:"effective_date"`
	ExpiryDate    *time.Time `json:"expiry_date"`
	Status        string     `json:"status"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// EmployeeDocument represents important documents for an employee
type EmployeeDocument struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	EmployeeID     string     `json:"employee_id"`
	DocumentType   string     `json:"document_type"`
	DocumentNumber string     `json:"document_number"`
	IssueDate      *time.Time `json:"issue_date"`
	ExpiryDate     *time.Time `json:"expiry_date"`
	DocumentPath   string     `json:"document_path"`
	Status         string     `json:"status"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// HRAuditLog tracks changes to HR entities
type HRAuditLog struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Action     string    `json:"action"`
	OldValues  string    `json:"old_values"`
	NewValues  string    `json:"new_values"`
	ChangedBy  string    `json:"changed_by"`
	ChangedAt  time.Time `json:"changed_at"`
}
