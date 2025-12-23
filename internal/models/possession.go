package models

import (
	"time"
)

// PossessionStatus represents the possession status of a booking
type PossessionStatus struct {
	ID                      int64      `gorm:"primaryKey" json:"id"`
	TenantID                int64      `json:"tenant_id"`
	BookingID               int64      `json:"booking_id"`
	Status                  string     `json:"status"` // pending, in_progress, completed, cancelled
	PossessionDate          *time.Time `json:"possession_date"`
	EstimatedPossessionDate *time.Time `json:"estimated_possession_date"`
	PossessionReason        *string    `json:"possession_reason"`
	PossessionType          string     `json:"possession_type"` // normal, partial, interim, final
	IsComplete              bool       `json:"is_complete"`
	CompletionPercentage    float64    `json:"completion_percentage"`
	Notes                   *string    `json:"notes"`
	CreatedBy               *int64     `json:"created_by"`
	UpdatedBy               *int64     `json:"updated_by"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
	DeletedAt               *time.Time `json:"deleted_at"`
}

// PossessionDocument represents a document in the possession process
type PossessionDocument struct {
	ID                int64      `gorm:"primaryKey" json:"id"`
	TenantID          int64      `json:"tenant_id"`
	PossessionID      int64      `json:"possession_id"`
	DocumentType      string     `json:"document_type"` // possession_letter, handover_checklist, keys, utilities_list, final_statement, insurance_doc
	DocumentName      string     `json:"document_name"`
	DocumentURL       *string    `json:"document_url"`
	FileName          *string    `json:"file_name"`
	FileSize          *int64     `json:"file_size"`
	FileFormat        *string    `json:"file_format"`
	S3Bucket          *string    `json:"s3_bucket"`
	S3Key             *string    `json:"s3_key"`
	DocumentStatus    string     `json:"document_status"` // pending, submitted, verified, rejected, approved
	VerificationNotes *string    `json:"verification_notes"`
	VerifiedBy        *int64     `json:"verified_by"`
	VerifiedAt        *time.Time `json:"verified_at"`
	IsMandatory       bool       `json:"is_mandatory"`
	UploadedBy        *int64     `json:"uploaded_by"`
	UploadedAt        *time.Time `json:"uploaded_at"`
	Metadata          *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy         *int64     `json:"created_by"`
	UpdatedBy         *int64     `json:"updated_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

// PossessionRegistration represents the registration process
type PossessionRegistration struct {
	ID                     int64      `gorm:"primaryKey" json:"id"`
	TenantID               int64      `json:"tenant_id"`
	PossessionID           int64      `json:"possession_id"`
	RegistrationType       string     `json:"registration_type"` // registration, name_transfer, title_transfer, mortgage_release
	RegistrationNumber     *string    `json:"registration_number"`
	RegistrationOffice     *string    `json:"registration_office"`
	RegistrationDate       *time.Time `json:"registration_date"`
	RegistrationStatus     string     `json:"registration_status"` // pending, in_progress, completed, rejected, appealed
	AmountPaid             float64    `json:"amount_paid"`
	AmountPending          float64    `json:"amount_pending"`
	PaymentMode            *string    `json:"payment_mode"` // online, cheque, demand_draft, cash, bank_transfer
	ReferenceNumber        *string    `json:"reference_number"`
	SubmissionDate         *time.Time `json:"submission_date"`
	ExpectedCompletionDate *time.Time `json:"expected_completion_date"`
	ActualCompletionDate   *time.Time `json:"actual_completion_date"`
	Remarks                *string    `json:"remarks"`
	ApprovedBy             *int64     `json:"approved_by"`
	ApprovedAt             *time.Time `json:"approved_at"`
	CreatedBy              *int64     `json:"created_by"`
	UpdatedBy              *int64     `json:"updated_by"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at"`
}

// PossessionCertificate represents possession certificates
type PossessionCertificate struct {
	ID                 int64      `gorm:"primaryKey" json:"id"`
	TenantID           int64      `json:"tenant_id"`
	PossessionID       int64      `json:"possession_id"`
	CertificateType    string     `json:"certificate_type"` // possession_certificate, occupancy_certificate, completion_certificate, no_dues_certificate
	CertificateNumber  *string    `json:"certificate_number"`
	IssuingAuthority   *string    `json:"issuing_authority"`
	IssueDate          *time.Time `json:"issue_date"`
	ValidityDate       *time.Time `json:"validity_date"`
	CertificateURL     *string    `json:"certificate_url"`
	FileName           *string    `json:"file_name"`
	FileSize           *int64     `json:"file_size"`
	S3Bucket           *string    `json:"s3_bucket"`
	S3Key              *string    `json:"s3_key"`
	CertificateStatus  string     `json:"certificate_status"`  // pending, issued, verified, expired, cancelled
	VerificationStatus string     `json:"verification_status"` // pending, verified, rejected
	VerifiedBy         *int64     `json:"verified_by"`
	VerifiedAt         *time.Time `json:"verified_at"`
	VerificationNotes  *string    `json:"verification_notes"`
	Metadata           *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy          *int64     `json:"created_by"`
	UpdatedBy          *int64     `json:"updated_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
}

// PossessionApproval represents approvals in the possession process
type PossessionApproval struct {
	ID                 int64      `gorm:"primaryKey" json:"id"`
	TenantID           int64      `json:"tenant_id"`
	PossessionID       int64      `json:"possession_id"`
	ApprovalType       string     `json:"approval_type"` // possession_approval, document_approval, registration_approval, final_approval
	ApproverID         int64      `json:"approver_id"`
	ApproverRole       *string    `json:"approver_role"`
	ApprovalStatus     string     `json:"approval_status"` // pending, approved, rejected, conditional
	ApprovalNotes      *string    `json:"approval_notes"`
	ConditionalRemarks *string    `json:"conditional_remarks"`
	ApprovalDate       *time.Time `json:"approval_date"`
	ValidFrom          *time.Time `json:"valid_from"`
	ValidTill          *time.Time `json:"valid_till"`
	SequenceOrder      int        `json:"sequence_order"`
	IsFinalApproval    bool       `json:"is_final_approval"`
	Metadata           *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy          *int64     `json:"created_by"`
	UpdatedBy          *int64     `json:"updated_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
}

// PossessionAuditLog represents audit logs for possession actions
type PossessionAuditLog struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	TenantID      int64     `json:"tenant_id"`
	PossessionID  int64     `json:"possession_id"`
	Action        string    `json:"action"`
	EntityType    string    `json:"entity_type"` // possession, document, registration, certificate, approval
	EntityID      int64     `json:"entity_id"`
	OldValue      *string   `json:"old_value"`
	NewValue      *string   `json:"new_value"`
	ChangeReason  *string   `json:"change_reason"`
	PerformedBy   *int64    `json:"performed_by"`
	UserIPAddress *string   `json:"user_ip_address"`
	UserAgent     *string   `json:"user_agent"`
	Metadata      *JSONMap  `gorm:"type:json" json:"metadata"`
	CreatedAt     time.Time `json:"created_at"`
}

// ============================================
// DTO Models for Request/Response Handling
// ============================================

// CreatePossessionRequest represents a request to create a possession status
type CreatePossessionRequest struct {
	BookingID               int64      `json:"booking_id" binding:"required"`
	Status                  string     `json:"status" binding:"required"`
	EstimatedPossessionDate *time.Time `json:"estimated_possession_date"`
	PossessionReason        *string    `json:"possession_reason"`
	PossessionType          string     `json:"possession_type" binding:"required"`
	Notes                   *string    `json:"notes"`
}

// UpdatePossessionRequest represents a request to update a possession status
type UpdatePossessionRequest struct {
	Status                  *string    `json:"status"`
	PossessionDate          *time.Time `json:"possession_date"`
	EstimatedPossessionDate *time.Time `json:"estimated_possession_date"`
	PossessionReason        *string    `json:"possession_reason"`
	IsComplete              *bool      `json:"is_complete"`
	CompletionPercentage    *float64   `json:"completion_percentage"`
	Notes                   *string    `json:"notes"`
}

// CreatePossessionDocumentRequest represents a request to create a possession document
type CreatePossessionDocumentRequest struct {
	DocumentType string   `json:"document_type" binding:"required"`
	DocumentName string   `json:"document_name" binding:"required"`
	DocumentURL  *string  `json:"document_url"`
	IsMandatory  bool     `json:"is_mandatory"`
	Metadata     *JSONMap `json:"metadata"`
}

// UpdatePossessionDocumentRequest represents a request to update a possession document
type UpdatePossessionDocumentRequest struct {
	DocumentStatus    *string  `json:"document_status"`
	VerificationNotes *string  `json:"verification_notes"`
	Metadata          *JSONMap `json:"metadata"`
}

// VerifyPossessionDocumentRequest represents a request to verify a document
type VerifyPossessionDocumentRequest struct {
	DocumentStatus    string  `json:"document_status" binding:"required"`
	VerificationNotes *string `json:"verification_notes"`
}

// CreatePossessionRegistrationRequest represents a request to create a registration
type CreatePossessionRegistrationRequest struct {
	RegistrationType       string     `json:"registration_type" binding:"required"`
	RegistrationNumber     *string    `json:"registration_number"`
	RegistrationOffice     *string    `json:"registration_office"`
	ExpectedCompletionDate *time.Time `json:"expected_completion_date"`
	PaymentMode            *string    `json:"payment_mode"`
}

// UpdatePossessionRegistrationRequest represents a request to update a registration
type UpdatePossessionRegistrationRequest struct {
	RegistrationStatus     *string    `json:"registration_status"`
	AmountPaid             *float64   `json:"amount_paid"`
	AmountPending          *float64   `json:"amount_pending"`
	ExpectedCompletionDate *time.Time `json:"expected_completion_date"`
	Remarks                *string    `json:"remarks"`
}

// ApproveRegistrationRequest represents a request to approve a registration
type ApproveRegistrationRequest struct {
	ApprovalStatus string  `json:"approval_status" binding:"required"` // approved, rejected, conditional
	ApprovalNotes  *string `json:"approval_notes"`
}

// CreatePossessionApprovalRequest represents a request to create an approval
type CreatePossessionApprovalRequest struct {
	ApprovalType    string  `json:"approval_type" binding:"required"`
	ApproverID      int64   `json:"approver_id" binding:"required"`
	ApprovalNotes   *string `json:"approval_notes"`
	IsFinalApproval bool    `json:"is_final_approval"`
	SequenceOrder   int     `json:"sequence_order"`
}

// ApprovePossessionRequest represents a request to approve a possession
type ApprovePossessionRequest struct {
	ApprovalStatus  string  `json:"approval_status" binding:"required"` // approved, rejected, conditional
	ApprovalNotes   *string `json:"approval_notes"`
	IsFinalApproval bool    `json:"is_final_approval"`
}

// PossessionSummaryResponse represents a summary response of possession status
type PossessionSummaryResponse struct {
	PossessionID            int64      `json:"possession_id"`
	BookingID               int64      `json:"booking_id"`
	Status                  string     `json:"status"`
	PossessionType          string     `json:"possession_type"`
	CompletionPercentage    float64    `json:"completion_percentage"`
	PossessionDate          *time.Time `json:"possession_date"`
	EstimatedPossessionDate *time.Time `json:"estimated_possession_date"`
	TotalDocuments          int        `json:"total_documents"`
	VerifiedDocuments       int        `json:"verified_documents"`
	PendingDocuments        int        `json:"pending_documents"`
	Registrations           int        `json:"registrations"`
	Certificates            int        `json:"certificates"`
	PendingApprovals        int        `json:"pending_approvals"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}
