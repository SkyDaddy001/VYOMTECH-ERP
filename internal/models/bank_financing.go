package models

import (
	"database/sql"
	"time"
)

// ============================================================================
// BankFinancing Model
// ============================================================================
type BankFinancing struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	TenantID  int64  `gorm:"index" json:"tenant_id"`
	BookingID int64  `gorm:"index" json:"booking_id"`
	BankID    *int64 `gorm:"index" json:"bank_id"`

	// Financing Details
	LoanAmount        float64 `json:"loan_amount"`
	SanctionedAmount  float64 `json:"sanctioned_amount"`
	DisbursedAmount   float64 `json:"disbursed_amount"`
	OutstandingAmount float64 `json:"outstanding_amount"`

	// Loan Details
	LoanType     string   `json:"loan_type"` // Home Loan, Construction Loan, Bridge Loan
	InterestRate *float64 `json:"interest_rate"`
	TenureMonths *int     `json:"tenure_months"`
	EMIAmount    *float64 `json:"emi_amount"`

	// Status & Dates
	Status                 string     `gorm:"index" json:"status"` // pending, approved, sanctioned, disbursing, completed, rejected
	ApplicationDate        *time.Time `json:"application_date"`
	ApprovalDate           *time.Time `json:"approval_date"`
	SanctionDate           *time.Time `json:"sanction_date"`
	ExpectedCompletionDate *time.Time `json:"expected_completion_date"`

	// Documents
	ApplicationRefNo  *string `gorm:"uniqueIndex" json:"application_ref_no"`
	SanctionLetterURL *string `json:"sanction_letter_url"`

	// Metadata
	CreatedBy *int64       `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedBy *int64       `json:"updated_by"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`

	// Relations
	Bank          *Bank                    `gorm:"foreignKey:BankID" json:"bank,omitempty"`
	Booking       *CustomerBooking         `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Disbursements []BankDisbursement       `gorm:"foreignKey:FinancingID" json:"disbursements,omitempty"`
	NOCs          []BankNOC                `gorm:"foreignKey:FinancingID" json:"nocs,omitempty"`
	Collections   []BankCollectionTracking `gorm:"foreignKey:FinancingID" json:"collections,omitempty"`
}

// TableName specifies the table name
func (BankFinancing) TableName() string {
	return "bank_financing"
}

// ============================================================================
// BankDisbursement Model
// ============================================================================
type BankDisbursement struct {
	ID          int64 `gorm:"primaryKey" json:"id"`
	TenantID    int64 `gorm:"index" json:"tenant_id"`
	FinancingID int64 `gorm:"index" json:"financing_id"`

	// Disbursement Details
	DisbursementNumber int      `json:"disbursement_number"`
	ScheduledAmount    float64  `json:"scheduled_amount"`
	ActualAmount       *float64 `json:"actual_amount"`

	// Milestone Linked
	MilestoneID         *int64 `gorm:"index" json:"milestone_id"`
	MilestonePercentage *int   `json:"milestone_percentage"`

	// Status & Dates
	Status          string     `gorm:"index" json:"status"` // pending, released, credited, delayed, cancelled
	ScheduledDate   time.Time  `gorm:"index" json:"scheduled_date"`
	ActualDate      *time.Time `json:"actual_date"`
	BankReferenceNo *string    `json:"bank_reference_no"`

	// Documentation
	ClaimDocumentURL    *string    `json:"claim_document_url"`
	ReleaseApprovalBy   *int64     `json:"release_approval_by"`
	ReleaseApprovalDate *time.Time `json:"release_approval_date"`

	// Metadata
	CreatedBy *int64       `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedBy *int64       `json:"updated_by"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`

	// Relations
	Financing *BankFinancing `gorm:"foreignKey:FinancingID" json:"financing,omitempty"`
}

// TableName specifies the table name
func (BankDisbursement) TableName() string {
	return "bank_disbursement"
}

// ============================================================================
// BankNOC Model (No Objection Certificate)
// ============================================================================
type BankNOC struct {
	ID          int64 `gorm:"primaryKey" json:"id"`
	TenantID    int64 `gorm:"index" json:"tenant_id"`
	FinancingID int64 `gorm:"index" json:"financing_id"`

	// NOC Details
	NOCType         string     `gorm:"index" json:"noc_type"` // Pre-sanction, Post-completion, Full-settlement
	NOCRequestDate  time.Time  `json:"noc_request_date"`
	NOCReceivedDate *time.Time `json:"noc_received_date"`

	// Documents
	NOCDocumentURL *string  `json:"noc_document_url"`
	NOCAmount      *float64 `json:"noc_amount"`

	// Status
	Status        string     `gorm:"index" json:"status"` // requested, issued, expired, cancelled
	IssuedByBank  *string    `json:"issued_by_bank"`
	ValidTillDate *time.Time `json:"valid_till_date"`

	// Remarks
	Remarks *string `json:"remarks"`

	// Metadata
	CreatedBy *int64       `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedBy *int64       `json:"updated_by"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`

	// Relations
	Financing *BankFinancing `gorm:"foreignKey:FinancingID" json:"financing,omitempty"`
}

// TableName specifies the table name
func (BankNOC) TableName() string {
	return "bank_noc"
}

// ============================================================================
// BankCollectionTracking Model
// ============================================================================
type BankCollectionTracking struct {
	ID          int64 `gorm:"primaryKey" json:"id"`
	TenantID    int64 `gorm:"index" json:"tenant_id"`
	FinancingID int64 `gorm:"index" json:"financing_id"`

	// Collection Details
	CollectionType   string    `json:"collection_type"` // EMI, Prepayment, Partial, Full-Settlement
	CollectionAmount float64   `json:"collection_amount"`
	CollectionDate   time.Time `gorm:"index" json:"collection_date"`

	// Payment Mode
	PaymentMode        *string `json:"payment_mode"` // Bank Transfer, Cheque, NEFT, RTGS
	PaymentReferenceNo *string `json:"payment_reference_no"`

	// EMI Details
	EMIMonth        *string  `gorm:"index" json:"emi_month"` // Jan-2025
	EMINumber       *int     `json:"emi_number"`
	PrincipalAmount *float64 `json:"principal_amount"`
	InterestAmount  *float64 `json:"interest_amount"`

	// Status
	Status               string     `gorm:"index" json:"status"` // pending, verified, credited, failed
	BankConfirmationDate *time.Time `json:"bank_confirmation_date"`

	// Metadata
	CreatedBy *int64       `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedBy *int64       `json:"updated_by"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`

	// Relations
	Financing *BankFinancing `gorm:"foreignKey:FinancingID" json:"financing,omitempty"`
}

// TableName specifies the table name
func (BankCollectionTracking) TableName() string {
	return "bank_collection_tracking"
}

// ============================================================================
// Bank Model
// ============================================================================
type Bank struct {
	ID       int64 `gorm:"primaryKey" json:"id"`
	TenantID int64 `gorm:"index" json:"tenant_id"`

	BankName      string  `json:"bank_name"`
	BranchName    *string `json:"branch_name"`
	IFSCCode      *string `json:"ifsc_code"`
	BranchContact *string `json:"branch_contact"`
	BranchEmail   *string `json:"branch_email"`

	// Contact Person
	RelationshipManagerName  *string `json:"relationship_manager_name"`
	RelationshipManagerPhone *string `json:"relationship_manager_phone"`
	RelationshipManagerEmail *string `json:"relationship_manager_email"`

	Status string `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name
func (Bank) TableName() string {
	return "bank"
}

// ============================================================================
// DTO Models for API Requests/Responses
// ============================================================================

// CreateFinancingRequest request to create new financing
type CreateFinancingRequest struct {
	BookingID        int64    `json:"booking_id" binding:"required"`
	BankID           *int64   `json:"bank_id"`
	LoanAmount       float64  `json:"loan_amount" binding:"required"`
	SanctionedAmount float64  `json:"sanctioned_amount" binding:"required"`
	LoanType         string   `json:"loan_type" binding:"required"`
	InterestRate     *float64 `json:"interest_rate"`
	TenureMonths     *int     `json:"tenure_months"`
	ApplicationRefNo *string  `json:"application_ref_no"`
}

// UpdateFinancingRequest request to update financing
type UpdateFinancingRequest struct {
	SanctionedAmount *float64   `json:"sanctioned_amount"`
	LoanType         *string    `json:"loan_type"`
	InterestRate     *float64   `json:"interest_rate"`
	TenureMonths     *int       `json:"tenure_months"`
	Status           *string    `json:"status"`
	ApprovalDate     *time.Time `json:"approval_date"`
	SanctionDate     *time.Time `json:"sanction_date"`
}

// CreateDisbursementRequest request to create disbursement
type CreateDisbursementRequest struct {
	FinancingID         int64     `json:"financing_id" binding:"required"`
	ScheduledAmount     float64   `json:"scheduled_amount" binding:"required"`
	MilestoneID         *int64    `json:"milestone_id"`
	MilestonePercentage *int      `json:"milestone_percentage"`
	ScheduledDate       time.Time `json:"scheduled_date" binding:"required"`
}

// CreateNOCRequest request to create NOC
type CreateNOCRequest struct {
	FinancingID    int64     `json:"financing_id" binding:"required"`
	NOCType        string    `json:"noc_type" binding:"required"`
	NOCRequestDate time.Time `json:"noc_request_date" binding:"required"`
	NOCAmount      *float64  `json:"noc_amount"`
}

// CreateCollectionRequest request to create collection
type CreateCollectionRequest struct {
	FinancingID      int64     `json:"financing_id" binding:"required"`
	CollectionType   string    `json:"collection_type" binding:"required"`
	CollectionAmount float64   `json:"collection_amount" binding:"required"`
	CollectionDate   time.Time `json:"collection_date" binding:"required"`
	PaymentMode      *string   `json:"payment_mode"`
	EMIMonth         *string   `json:"emi_month"`
	EMINumber        *int      `json:"emi_number"`
}

// BankFinancingResponse response with financing details
type BankFinancingResponse struct {
	*BankFinancing
	RemainingBalance       float64 `json:"remaining_balance"`
	DisbursementPercentage float64 `json:"disbursement_percentage"`
	CollectionPercentage   float64 `json:"collection_percentage"`
}
