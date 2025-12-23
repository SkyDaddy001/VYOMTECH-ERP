package models

import (
	"database/sql"
	"time"
)

// ============================================================================
// JointApplicant Model
// ============================================================================
type JointApplicant struct {
	ID                           int64          `gorm:"primaryKey" json:"id"`
	TenantID                     int64          `gorm:"index" json:"tenant_id"`
	BookingID                    int64          `gorm:"index" json:"booking_id"`
	
	// Primary Applicant Flag
	IsPrimaryApplicant           bool           `json:"is_primary_applicant"`
	
	// Applicant Details
	ApplicantName                string         `json:"applicant_name"`
	DateOfBirth                  *time.Time     `json:"date_of_birth"`
	Gender                       *string        `json:"gender"`
	MaritalStatus                *string        `json:"marital_status"`
	Nationality                  *string        `json:"nationality"`
	
	// Contact Information
	Email                        *string        `json:"email"`
	PhoneNumber                  *string        `json:"phone_number"`
	AlternatePhone               *string        `json:"alternate_phone"`
	AddressLine1                 *string        `json:"address_line1"`
	AddressLine2                 *string        `json:"address_line2"`
	City                         *string        `json:"city"`
	StateProvince                *string        `json:"state_province"`
	PostalCode                   *string        `json:"postal_code"`
	Country                      *string        `json:"country"`
	
	// Identification
	IDType                       *string        `json:"id_type"`
	IDNumber                     *string        `json:"id_number"`
	IDIssueDate                  *time.Time     `json:"id_issue_date"`
	IDExpiryDate                 *time.Time     `json:"id_expiry_date"`
	
	// Occupation & Income
	Occupation                   *string        `json:"occupation"`
	EmployerName                 *string        `json:"employer_name"`
	AnnualIncome                 *float64       `json:"annual_income"`
	IncomeSource                 *string        `json:"income_source"`
	BankName                     *string        `json:"bank_name"`
	AccountNumber                *string        `json:"account_number"`
	AccountType                  *string        `json:"account_type"`
	
	// Co-Ownership Details
	OwnershipSharePercentage     *float64       `json:"ownership_share_percentage"`
	OwnershipType                *string        `json:"ownership_type"`
	
	// Verification Status
	KYCStatus                    *string        `json:"kyc_status"`
	KYCVerifiedDate              *time.Time     `json:"kyc_verified_date"`
	IncomeVerified               bool           `json:"income_verified"`
	IncomeVerifiedDate           *time.Time     `json:"income_verified_date"`
	DocumentVerificationStatus   *string        `json:"document_verification_status"`
	
	// Legal Consent
	ConsentGiven                 bool           `json:"consent_given"`
	ConsentGivenDate             *time.Time     `json:"consent_given_date"`
	LegalNoticeAcknowledged      bool           `json:"legal_notice_acknowledged"`
	
	// Metadata
	CreatedBy                    *int64         `json:"created_by"`
	CreatedAt                    time.Time      `json:"created_at"`
	UpdatedBy                    *int64         `json:"updated_by"`
	UpdatedAt                    time.Time      `json:"updated_at"`
	DeletedAt                    sql.NullTime   `gorm:"index" json:"deleted_at"`
	
	// Relations
	Documents                    []JointApplicantDocument       `gorm:"foreignKey:JointApplicantID" json:"documents,omitempty"`
	IncomeVerification           *JointApplicantIncomeVerification `gorm:"foreignKey:JointApplicantID" json:"income_verification,omitempty"`
	Liabilities                  []JointApplicantLiability      `gorm:"foreignKey:JointApplicantID" json:"liabilities,omitempty"`
}

func (JointApplicant) TableName() string {
	return "joint_applicant"
}

// ============================================================================
// JointApplicantDocument Model
// ============================================================================
type JointApplicantDocument struct {
	ID                  int64          `gorm:"primaryKey" json:"id"`
	TenantID            int64          `gorm:"index" json:"tenant_id"`
	JointApplicantID    int64          `gorm:"index" json:"joint_applicant_id"`
	
	// Document Details
	DocumentType        *string        `json:"document_type"`
	DocumentName        string         `json:"document_name"`
	DocumentURL         string         `json:"document_url"`
	DocumentHash        *string        `json:"document_hash"`
	
	// File Information
	FileSize            *int64         `json:"file_size"`
	FileType            *string        `json:"file_type"`
	UploadDate          time.Time      `json:"upload_date"`
	
	// Verification
	DocumentStatus      *string        `json:"document_status"`
	VerificationDate    *time.Time     `json:"verification_date"`
	VerificationNotes   *string        `json:"verification_notes"`
	
	// Metadata
	CreatedBy           *int64         `json:"created_by"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedBy           *int64         `json:"updated_by"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           sql.NullTime   `gorm:"index" json:"deleted_at"`
}

func (JointApplicantDocument) TableName() string {
	return "joint_applicant_document"
}

// ============================================================================
// CoOwnershipAgreement Model
// ============================================================================
type CoOwnershipAgreement struct {
	ID                          int64          `gorm:"primaryKey" json:"id"`
	TenantID                    int64          `gorm:"index" json:"tenant_id"`
	BookingID                   int64          `gorm:"index" json:"booking_id"`
	
	// Agreement Details
	AgreementReferenceNo        string         `gorm:"uniqueIndex" json:"agreement_reference_no"`
	AgreementDate               *time.Time     `json:"agreement_date"`
	AgreementType               *string        `json:"agreement_type"`
	
	// Property Details
	PropertyDescription         *string        `json:"property_description"`
	AgreedPurchasePrice         *float64       `json:"agreed_purchase_price"`
	
	// Ownership Distribution
	TotalSharePercentage        *float64       `json:"total_share_percentage"`
	
	// Status & Verification
	AgreementStatus             string         `gorm:"index" json:"agreement_status"`
	Notarized                   bool           `json:"notarized"`
	NotarizationDate            *time.Time     `json:"notarization_date"`
	RegisteredWithAuthorities   bool           `json:"registered_with_authorities"`
	RegistrationNumber          *string        `json:"registration_number"`
	RegistrationDate            *time.Time     `json:"registration_date"`
	
	// Approvals
	AllPartiesAgreed            bool           `json:"all_parties_agreed"`
	LegalReviewCompleted        bool           `json:"legal_review_completed"`
	LegalReviewerNotes          *string        `json:"legal_reviewer_notes"`
	
	// Document Storage
	AgreementDocumentURL        *string        `json:"agreement_document_url"`
	NotarizedCopyURL            *string        `json:"notarized_copy_url"`
	RegistrationCertificateURL  *string        `json:"registration_certificate_url"`
	
	// Metadata
	CreatedBy                   *int64         `json:"created_by"`
	CreatedAt                   time.Time      `json:"created_at"`
	UpdatedBy                   *int64         `json:"updated_by"`
	UpdatedAt                   time.Time      `json:"updated_at"`
	DeletedAt                   sql.NullTime   `gorm:"index" json:"deleted_at"`
	
	// Relations
	Signatories                 []CoOwnershipSignatory `gorm:"foreignKey:CoOwnershipAgreementID" json:"signatories,omitempty"`
}

func (CoOwnershipAgreement) TableName() string {
	return "co_ownership_agreement"
}

// ============================================================================
// CoOwnershipSignatory Model
// ============================================================================
type CoOwnershipSignatory struct {
	ID                      int64          `gorm:"primaryKey" json:"id"`
	TenantID                int64          `gorm:"index" json:"tenant_id"`
	CoOwnershipAgreementID  int64          `gorm:"index" json:"co_ownership_agreement_id"`
	JointApplicantID        int64          `gorm:"index" json:"joint_applicant_id"`
	
	// Signatory Details
	SignatoryName           string         `json:"signatory_name"`
	SignatoryEmail          *string        `json:"signatory_email"`
	SignatoryPhone          *string        `json:"signatory_phone"`
	SignatureType           *string        `json:"signature_type"`
	
	// Signature Status
	SignatureStatus         string         `gorm:"index" json:"signature_status"`
	SignatureDate           *time.Time     `json:"signature_date"`
	SignatureIPAddress      *string        `json:"signature_ip_address"`
	SignatureDeviceInfo     *string        `json:"signature_device_info"`
	
	// Signature Document
	SignatureImageURL       *string        `json:"signature_image_url"`
	SignedDocumentURL       *string        `json:"signed_document_url"`
	
	// Witness Information
	WitnessName             *string        `json:"witness_name"`
	WitnessIDType           *string        `json:"witness_id_type"`
	WitnessIDNumber         *string        `json:"witness_id_number"`
	WitnessSignatureURL     *string        `json:"witness_signature_url"`
	
	// Metadata
	CreatedBy               *int64         `json:"created_by"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedBy               *int64         `json:"updated_by"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               sql.NullTime   `gorm:"index" json:"deleted_at"`
}

func (CoOwnershipSignatory) TableName() string {
	return "co_ownership_signatory"
}

// ============================================================================
// JointApplicantIncomeVerification Model
// ============================================================================
type JointApplicantIncomeVerification struct {
	ID                     int64          `gorm:"primaryKey" json:"id"`
	TenantID               int64          `gorm:"index" json:"tenant_id"`
	JointApplicantID       int64          `gorm:"index" json:"joint_applicant_id"`
	
	// Income Details
	AnnualIncome           *float64       `json:"annual_income"`
	MonthlyIncome          *float64       `json:"monthly_income"`
	IncomeYear             *int           `json:"income_year"`
	IncomeTaxSlab          *string        `json:"income_tax_slab"`
	
	// Source of Income
	PrimaryOccupation      *string        `json:"primary_occupation"`
	PrimaryEmployer        *string        `json:"primary_employer"`
	EmploymentType         *string        `json:"employment_type"`
	
	// Financial Status
	TotalAssets            *float64       `json:"total_assets"`
	TotalLiabilities       *float64       `json:"total_liabilities"`
	NetWorth               *float64       `json:"net_worth"`
	
	// Bank Details
	PrimaryBankName        *string        `json:"primary_bank_name"`
	PrimaryAccountNumber   *string        `json:"primary_account_number"`
	PrimaryAccountType     *string        `json:"primary_account_type"`
	BankVerificationStatus *string        `json:"bank_verification_status"`
	
	// Income Tax Information
	PANNumber              *string        `gorm:"index" json:"pan_number"`
	AadharNumber           *string        `json:"aadhar_number"`
	LastITRFiledYear       *int           `json:"last_itr_filed_year"`
	LastITRAmount          *float64       `json:"last_itr_amount"`
	
	// Verification Documents
	SalarySlipURL          *string        `json:"salary_slip_url"`
	ITRDocumentURL         *string        `json:"itr_document_url"`
	BankStatementURL       *string        `json:"bank_statement_url"`
	EmployerLetterURL      *string        `json:"employer_letter_url"`
	
	// Verification Status
	VerificationStatus     string         `gorm:"index" json:"verification_status"`
	VerificationDate       *time.Time     `json:"verification_date"`
	VerifierID             *int64         `json:"verifier_id"`
	VerificationNotes      *string        `json:"verification_notes"`
	
	// Metadata
	CreatedBy              *int64         `json:"created_by"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedBy              *int64         `json:"updated_by"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              sql.NullTime   `gorm:"index" json:"deleted_at"`
}

func (JointApplicantIncomeVerification) TableName() string {
	return "joint_applicant_income_verification"
}

// ============================================================================
// JointApplicantLiability Model
// ============================================================================
type JointApplicantLiability struct {
	ID                  int64          `gorm:"primaryKey" json:"id"`
	TenantID            int64          `gorm:"index" json:"tenant_id"`
	JointApplicantID    int64          `gorm:"index" json:"joint_applicant_id"`
	
	// Liability Details
	LiabilityType       *string        `json:"liability_type"`
	CreditorName        string         `json:"creditor_name"`
	
	// Loan/Credit Details
	OutstandingAmount   *float64       `json:"outstanding_amount"`
	MonthlyPayment      *float64       `json:"monthly_payment"`
	InterestRate        *float64       `json:"interest_rate"`
	LoanMaturityDate    *time.Time     `json:"loan_maturity_date"`
	
	// Liability Status
	Status              *string        `json:"status"`
	
	// Documents
	DocumentURL         *string        `json:"document_url"`
	
	// Metadata
	CreatedBy           *int64         `json:"created_by"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedBy           *int64         `json:"updated_by"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           sql.NullTime   `gorm:"index" json:"deleted_at"`
}

func (JointApplicantLiability) TableName() string {
	return "joint_applicant_liability"
}

// ============================================================================
// DTO Models
// ============================================================================

// CreateJointApplicantRequest request to create a new joint applicant
type CreateJointApplicantRequest struct {
	BookingID                int64   `json:"booking_id" binding:"required"`
	IsPrimaryApplicant       bool    `json:"is_primary_applicant"`
	ApplicantName            string  `json:"applicant_name" binding:"required"`
	DateOfBirth              *string `json:"date_of_birth"`
	Gender                   *string `json:"gender"`
	Email                    *string `json:"email"`
	PhoneNumber              *string `json:"phone_number"`
	IDType                   *string `json:"id_type"`
	IDNumber                 *string `json:"id_number"`
	Occupation               *string `json:"occupation"`
	AnnualIncome             *float64 `json:"annual_income"`
	OwnershipSharePercentage *float64 `json:"ownership_share_percentage"`
	OwnershipType            *string `json:"ownership_type"`
}

// UpdateJointApplicantRequest request to update a joint applicant
type UpdateJointApplicantRequest struct {
	ApplicantName                *string  `json:"applicant_name"`
	Email                        *string  `json:"email"`
	PhoneNumber                  *string  `json:"phone_number"`
	Occupation                   *string  `json:"occupation"`
	AnnualIncome                 *float64 `json:"annual_income"`
	KYCStatus                    *string  `json:"kyc_status"`
	IncomeVerified               *bool    `json:"income_verified"`
	DocumentVerificationStatus   *string  `json:"document_verification_status"`
	ConsentGiven                 *bool    `json:"consent_given"`
}

// CreateCoOwnershipAgreementRequest request to create a co-ownership agreement
type CreateCoOwnershipAgreementRequest struct {
	BookingID                int64    `json:"booking_id" binding:"required"`
	AgreementType            string   `json:"agreement_type" binding:"required"`
	PropertyDescription      *string  `json:"property_description"`
	AgreedPurchasePrice      *float64 `json:"agreed_purchase_price"`
	TotalSharePercentage     *float64 `json:"total_share_percentage"`
}

// JointApplicantSummaryResponse response with joint applicant summary
type JointApplicantSummaryResponse struct {
	BookingID              int64   `json:"booking_id"`
	TotalApplicants        int     `json:"total_applicants"`
	PrimaryApplicantName   string  `json:"primary_applicant_name"`
	TotalOwnershipShare    float64 `json:"total_ownership_share"`
	KYCCompletionPercent   float64 `json:"kyc_completion_percent"`
	IncomeVerifiedCount    int     `json:"income_verified_count"`
	DocumentsVerifiedCount int     `json:"documents_verified_count"`
	AgreementStatus        string  `json:"agreement_status"`
}
