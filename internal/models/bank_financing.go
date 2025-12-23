package models

import "time"

// ============================================================================
// BANK FINANCING MODELS
// ============================================================================

type BankFinancing struct {
	ID                     string     `json:"id" db:"id"`
	TenantID               string     `json:"tenant_id" db:"tenant_id"`
	BookingID              string     `json:"booking_id" db:"booking_id"`
	BankID                 *string    `json:"bank_id" db:"bank_id"`
	LoanAmount             float64    `json:"loan_amount" db:"loan_amount"`
	SanctionedAmount       float64    `json:"sanctioned_amount" db:"sanctioned_amount"`
	DisbursedAmount        float64    `json:"disbursed_amount" db:"disbursed_amount"`
	OutstandingAmount      float64    `json:"outstanding_amount" db:"outstanding_amount"`
	LoanType               string     `json:"loan_type" db:"loan_type"` // Home Loan, Construction Loan, Bridge Loan
	InterestRate           *float64   `json:"interest_rate" db:"interest_rate"`
	TenureMonths           *int       `json:"tenure_months" db:"tenure_months"`
	EMIAmount              *float64   `json:"emi_amount" db:"emi_amount"`
	Status                 string     `json:"status" db:"status"` // pending, approved, sanctioned, disbursing, completed, rejected
	ApplicationDate        *time.Time `json:"application_date" db:"application_date"`
	ApprovalDate           *time.Time `json:"approval_date" db:"approval_date"`
	SanctionDate           *time.Time `json:"sanction_date" db:"sanction_date"`
	ExpectedCompletionDate *time.Time `json:"expected_completion_date" db:"expected_completion_date"`
	ApplicationRefNo       *string    `json:"application_ref_no" db:"application_ref_no"`
	SanctionLetterURL      *string    `json:"sanction_letter_url" db:"sanction_letter_url"`
	CreatedBy              *string    `json:"created_by" db:"created_by"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedBy              *string    `json:"updated_by" db:"updated_by"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at" db:"deleted_at"`
}

type BankDisbursement struct {
	ID                  string     `json:"id" db:"id"`
	TenantID            string     `json:"tenant_id" db:"tenant_id"`
	FinancingID         string     `json:"financing_id" db:"financing_id"`
	DisbursementNumber  int        `json:"disbursement_number" db:"disbursement_number"`
	ScheduledAmount     float64    `json:"scheduled_amount" db:"scheduled_amount"`
	ActualAmount        *float64   `json:"actual_amount" db:"actual_amount"`
	MilestoneID         *string    `json:"milestone_id" db:"milestone_id"`
	MilestonePercentage *int       `json:"milestone_percentage" db:"milestone_percentage"`
	Status              string     `json:"status" db:"status"` // pending, released, credited, delayed, cancelled
	ScheduledDate       time.Time  `json:"scheduled_date" db:"scheduled_date"`
	ActualDate          *time.Time `json:"actual_date" db:"actual_date"`
	BankReferenceNo     *string    `json:"bank_reference_no" db:"bank_reference_no"`
	ClaimDocumentURL    *string    `json:"claim_document_url" db:"claim_document_url"`
	ReleaseApprovalBy   *string    `json:"release_approval_by" db:"release_approval_by"`
	ReleaseApprovalDate *time.Time `json:"release_approval_date" db:"release_approval_date"`
	CreatedBy           *string    `json:"created_by" db:"created_by"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedBy           *string    `json:"updated_by" db:"updated_by"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at" db:"deleted_at"`
}

type BankNOC struct {
	ID              string     `json:"id" db:"id"`
	TenantID        string     `json:"tenant_id" db:"tenant_id"`
	FinancingID     string     `json:"financing_id" db:"financing_id"`
	NOCType         string     `json:"noc_type" db:"noc_type"` // Pre-sanction, Post-completion, Full-settlement
	NOCRequestDate  time.Time  `json:"noc_request_date" db:"noc_request_date"`
	NOCReceivedDate *time.Time `json:"noc_received_date" db:"noc_received_date"`
	NOCDocumentURL  *string    `json:"noc_document_url" db:"noc_document_url"`
	NOCAmount       *float64   `json:"noc_amount" db:"noc_amount"`
	Status          string     `json:"status" db:"status"` // requested, issued, expired, cancelled
	IssuedByBank    *string    `json:"issued_by_bank" db:"issued_by_bank"`
	ValidTillDate   *time.Time `json:"valid_till_date" db:"valid_till_date"`
	Remarks         *string    `json:"remarks" db:"remarks"`
	CreatedBy       *string    `json:"created_by" db:"created_by"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedBy       *string    `json:"updated_by" db:"updated_by"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at" db:"deleted_at"`
}

type BankCollectionTracking struct {
	ID                   string     `json:"id" db:"id"`
	TenantID             string     `json:"tenant_id" db:"tenant_id"`
	FinancingID          string     `json:"financing_id" db:"financing_id"`
	CollectionType       string     `json:"collection_type" db:"collection_type"` // EMI, Prepayment, Partial, Full-Settlement
	CollectionAmount     float64    `json:"collection_amount" db:"collection_amount"`
	CollectionDate       time.Time  `json:"collection_date" db:"collection_date"`
	PaymentMode          *string    `json:"payment_mode" db:"payment_mode"` // Bank Transfer, Cheque, NEFT, RTGS
	PaymentReferenceNo   *string    `json:"payment_reference_no" db:"payment_reference_no"`
	EMIMonth             *string    `json:"emi_month" db:"emi_month"` // e.g., 'Jan-2025'
	EMINumber            *int       `json:"emi_number" db:"emi_number"`
	PrincipalAmount      *float64   `json:"principal_amount" db:"principal_amount"`
	InterestAmount       *float64   `json:"interest_amount" db:"interest_amount"`
	Status               string     `json:"status" db:"status"` // pending, verified, credited, failed
	BankConfirmationDate *time.Time `json:"bank_confirmation_date" db:"bank_confirmation_date"`
	CreatedBy            *string    `json:"created_by" db:"created_by"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedBy            *string    `json:"updated_by" db:"updated_by"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at" db:"deleted_at"`
}

type Bank struct {
	ID                       string    `json:"id" db:"id"`
	TenantID                 string    `json:"tenant_id" db:"tenant_id"`
	BankName                 string    `json:"bank_name" db:"bank_name"`
	BranchName               string    `json:"branch_name" db:"branch_name"`
	IFSCCode                 *string   `json:"ifsc_code" db:"ifsc_code"`
	BranchContact            *string   `json:"branch_contact" db:"branch_contact"`
	BranchEmail              *string   `json:"branch_email" db:"branch_email"`
	RelationshipManagerName  *string   `json:"relationship_manager_name" db:"relationship_manager_name"`
	RelationshipManagerPhone *string   `json:"relationship_manager_phone" db:"relationship_manager_phone"`
	RelationshipManagerEmail *string   `json:"relationship_manager_email" db:"relationship_manager_email"`
	Status                   string    `json:"status" db:"status"` // active, inactive
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`
}
