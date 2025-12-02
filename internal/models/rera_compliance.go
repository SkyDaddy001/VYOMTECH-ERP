package models

import "time"

// ============================================================================
// PROJECT COLLECTION ACCOUNTS (RERA COMPLIANCE)
// ============================================================================

// ProjectCollectionAccount represents a segregated bank account for each project
// As per RERA regulations, each project must have a separate collection account
type ProjectCollectionAccount struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	ProjectID string `json:"project_id"`

	// Account Details
	AccountName   string `json:"account_name"`
	AccountCode   string `json:"account_code"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	IFSCCode      string `json:"ifsc_code"`
	AccountType   string `json:"account_type"` // Current, Savings

	// Balance Management
	OpeningBalance float64 `json:"opening_balance"`
	CurrentBalance float64 `json:"current_balance"`
	MinimumBalance float64 `json:"minimum_balance"`

	// RERA Compliance
	RERACompliant       bool    `json:"rera_compliant"`
	RegulatedAccount    bool    `json:"regulated_account"`     // Must be segregated
	MaxBorrowingAllowed float64 `json:"max_borrowing_allowed"` // 10% of collections
	CurrentBorrowing    float64 `json:"current_borrowing"`

	// Status
	Status     string     `json:"status"` // Active, Inactive, Closed
	OpenedDate time.Time  `json:"opened_date"`
	ClosedDate *time.Time `json:"closed_date"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedBy *string    `json:"created_by"`
}

// ProjectCollectionLedger tracks all collections received for a project
type ProjectCollectionLedger struct {
	ID                  string `json:"id"`
	TenantID            string `json:"tenant_id"`
	ProjectID           string `json:"project_id"`
	CollectionAccountID string `json:"collection_account_id"`

	// Collection Details
	CollectionDate   time.Time `json:"collection_date"`
	CollectionNumber string    `json:"collection_number"`
	BookingID        string    `json:"booking_id"`
	CustomerID       string    `json:"customer_id"`
	UnitID           string    `json:"unit_id"`

	// Payment Information
	PaymentMode     string  `json:"payment_mode"` // Cash, Cheque, NEFT, RTGS, IFT, Demand_Draft
	AmountCollected float64 `json:"amount_collected"`

	// Cheque Details
	ChequeNumber string     `json:"cheque_number"`
	ChequeDate   *time.Time `json:"cheque_date"`
	ChequeStatus string     `json:"cheque_status"` // Not_Applicable, Pending, Cleared, Bounced, Post_Dated

	// Reference
	ReferenceNumber string `json:"reference_number"`
	PaidBy          string `json:"paid_by"`
	Remarks         string `json:"remarks"`

	// GL Integration
	GLPosted       bool       `json:"gl_posted"`
	JournalEntryID string     `json:"journal_entry_id"`
	PostedAt       *time.Time `json:"posted_at"`
	PostedBy       *string    `json:"posted_by"`

	// Status
	Status string `json:"status"` // Pending, Approved, Reversed, Cancelled

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedBy *string    `json:"created_by"`
}

// CollectionAgainstMilestone links collections to payment milestones
type CollectionAgainstMilestone struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	ProjectID string `json:"project_id"`

	CollectionID      string `json:"collection_id"`
	BookingID         string `json:"booking_id"`
	PaymentScheduleID string `json:"payment_schedule_id"`

	// Amount Mapping
	ScheduledAmount float64 `json:"scheduled_amount"`
	CollectedAmount float64 `json:"collected_amount"`
	ExcessAmount    float64 `json:"excess_amount"`
	ShortfallAmount float64 `json:"shortfall_amount"`

	// Dates
	ScheduledDate  *time.Time `json:"scheduled_date"`
	CollectionDate time.Time  `json:"collection_date"`
	DaysDelayed    int        `json:"days_delayed"`

	// Interest on Delayed Collection
	InterestApplicable bool    `json:"interest_applicable"`
	InterestRate       float64 `json:"interest_rate"` // % per month
	InterestAmount     float64 `json:"interest_amount"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// ProjectFundUtilization tracks how collection funds are used (RERA requirement)
type ProjectFundUtilization struct {
	ID                  string `json:"id"`
	TenantID            string `json:"tenant_id"`
	ProjectID           string `json:"project_id"`
	CollectionAccountID string `json:"collection_account_id"`

	// Utilization Details
	UtilizationDate time.Time `json:"utilization_date"`
	UtilizationType string    `json:"utilization_type"` // Construction, Land_Cost, Statutory_Approval, Admin, Interest
	Description     string    `json:"description"`

	// Amount
	AmountUtilized float64 `json:"amount_utilized"`

	// Supporting Documents
	BillNumber string     `json:"bill_number"`
	BillDate   *time.Time `json:"bill_date"`
	InvoiceID  string     `json:"invoice_id"`
	BillAmount float64    `json:"bill_amount"`

	// Approval
	ApprovedBy      *string    `json:"approved_by"`
	ApprovalDate    *time.Time `json:"approval_date"`
	ApprovalRemarks string     `json:"approval_remarks"`

	// GL Integration
	GLPosted       bool   `json:"gl_posted"`
	JournalEntryID string `json:"journal_entry_id"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// ProjectAccountBorrowing tracks borrowings against project collections (RERA - Max 10%)
type ProjectAccountBorrowing struct {
	ID                  string `json:"id"`
	TenantID            string `json:"tenant_id"`
	ProjectID           string `json:"project_id"`
	CollectionAccountID string `json:"collection_account_id"`

	// Loan Details
	LoanNumber       string    `json:"loan_number"`
	LoanDate         time.Time `json:"loan_date"`
	PrincipalAmount  float64   `json:"principal_amount"`
	InterestRate     float64   `json:"interest_rate"`
	LoanTenureMonths int       `json:"loan_tenure_months"`

	// Collections-Based Limits
	TotalCollectionsTillDate float64 `json:"total_collections_till_date"`
	MaxBorrowingAllowed      float64 `json:"max_borrowing_allowed"` // 10% of collections

	// Repayment
	MaturityDate      *time.Time `json:"maturity_date"`
	RepaidAmount      float64    `json:"repaid_amount"`
	OutstandingAmount float64    `json:"outstanding_amount"`

	// Status
	Status string `json:"status"` // Active, Repaid, Overdue, Waived

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// ProjectAccountReconciliation represents monthly/quarterly reconciliation
type ProjectAccountReconciliation struct {
	ID                  string `json:"id"`
	TenantID            string `json:"tenant_id"`
	ProjectID           string `json:"project_id"`
	CollectionAccountID string `json:"collection_account_id"`

	// Period
	ReconciliationDate time.Time `json:"reconciliation_date"`
	PeriodFrom         time.Time `json:"period_from"`
	PeriodTo           time.Time `json:"period_to"`

	// Collections
	OpeningBalance    float64 `json:"opening_balance"`
	TotalCollections  float64 `json:"total_collections"`
	TotalUtilizations float64 `json:"total_utilizations"`
	Borrowings        float64 `json:"borrowings"`
	InterestAccrued   float64 `json:"interest_accrued"`
	ClosingBalance    float64 `json:"closing_balance"`

	// Bank Reconciliation
	BankBalance            *float64 `json:"bank_balance"`
	BookBalance            *float64 `json:"book_balance"`
	Reconciled             bool     `json:"reconciled"`
	ReconciliationVariance float64  `json:"reconciliation_variance"`

	// Approval
	VerifiedBy   *string    `json:"verified_by"`
	VerifiedDate *time.Time `json:"verified_date"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// ============================================================================
// REQUEST/RESPONSE MODELS
// ============================================================================

// CreateProjectCollectionAccountRequest for creating a new project collection account
type CreateProjectCollectionAccountRequest struct {
	ProjectID      string  `json:"project_id" validate:"required"`
	AccountName    string  `json:"account_name" validate:"required"`
	AccountCode    string  `json:"account_code" validate:"required"`
	BankName       string  `json:"bank_name" validate:"required"`
	AccountNumber  string  `json:"account_number" validate:"required"`
	IFSCCode       string  `json:"ifsc_code" validate:"required"`
	AccountType    string  `json:"account_type"`
	MinimumBalance float64 `json:"minimum_balance"`
}

// ProjectCollectionResponse for returning project collection data
type ProjectCollectionResponse struct {
	ID                   string                    `json:"id"`
	ProjectName          string                    `json:"project_name"`
	CollectionAccount    *ProjectCollectionAccount `json:"collection_account"`
	RecentCollections    []ProjectCollectionLedger `json:"recent_collections"`
	TotalCollected       float64                   `json:"total_collected"`
	TotalUtilized        float64                   `json:"total_utilized"`
	AvailableBalance     float64                   `json:"available_balance"`
	RERAComplianceStatus string                    `json:"rera_compliance_status"`
}
