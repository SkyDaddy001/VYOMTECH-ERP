package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// DocumentCategory represents a document category
type DocumentCategory struct {
	ID           int64      `gorm:"primaryKey" json:"id"`
	TenantID     int64      `json:"tenant_id"`
	CategoryName string     `json:"category_name"`
	CategoryCode string     `json:"category_code"`
	Description  *string    `json:"description"`
	IconURL      *string    `json:"icon_url"`
	DisplayOrder int        `json:"display_order"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

// DocumentType represents a type of document
type DocumentType struct {
	ID              int64      `gorm:"primaryKey" json:"id"`
	TenantID        int64      `json:"tenant_id"`
	CategoryID      int64      `json:"category_id"`
	TypeName        string     `json:"type_name"`
	TypeCode        string     `json:"type_code"`
	Description     *string    `json:"description"`
	IsMandatory     bool       `json:"is_mandatory"`
	IsIdentityProof bool       `json:"is_identity_proof"`
	IsPropertyDoc   bool       `json:"is_property_doc"`
	IsFinancialDoc  bool       `json:"is_financial_doc"`
	ExpiryRequired  bool       `json:"expiry_required"`
	FileFormats     *string    `json:"file_formats"`
	MaxFileSize     int        `json:"max_file_size"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

// Document represents a document
type Document struct {
	ID                  int64      `gorm:"primaryKey" json:"id"`
	TenantID            int64      `json:"tenant_id"`
	DocumentTypeID      int64      `json:"document_type_id"`
	EntityType          string     `json:"entity_type"`
	EntityID            int64      `json:"entity_id"`
	DocumentName        string     `json:"document_name"`
	DocumentDescription *string    `json:"document_description"`
	DocumentURL         string     `json:"document_url"`
	FileName            string     `json:"file_name"`
	FileSize            *int64     `json:"file_size"`
	FileFormat          *string    `json:"file_format"`
	S3Bucket            *string    `json:"s3_bucket"`
	S3Key               *string    `json:"s3_key"`
	ThumbnailURL        *string    `json:"thumbnail_url"`
	DocumentStatus      string     `json:"document_status"`
	VerificationStatus  string     `json:"verification_status"`
	VerificationNotes   *string    `json:"verification_notes"`
	VerifierID          *int64     `json:"verifier_id"`
	VerifiedAt          *time.Time `json:"verified_at"`
	IssueDate           *time.Time `json:"issue_date"`
	ExpiryDate          *time.Time `json:"expiry_date"`
	IsPrimary           bool       `json:"is_primary"`
	Metadata            *JSONMap   `gorm:"type:json" json:"metadata"`
	Tags                *string    `json:"tags"`
	CreatedBy           *int64     `json:"created_by"`
	UpdatedBy           *int64     `json:"updated_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

// DocumentTemplate represents a document template
type DocumentTemplate struct {
	ID                  int64      `gorm:"primaryKey" json:"id"`
	TenantID            int64      `json:"tenant_id"`
	TemplateName        string     `json:"template_name"`
	TemplateDescription *string    `json:"template_description"`
	TemplateType        string     `json:"template_type"`
	TemplateContent     *string    `gorm:"type:longtext" json:"template_content"`
	RequiredDocuments   *JSONArray `gorm:"type:json" json:"required_documents"`
	IsActive            bool       `json:"is_active"`
	CreatedBy           *int64     `json:"created_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

// DocumentVerification represents a verification of a document
type DocumentVerification struct {
	ID                          int64     `gorm:"primaryKey" json:"id"`
	TenantID                    int64     `json:"tenant_id"`
	DocumentID                  int64     `json:"document_id"`
	VerificationType            string    `json:"verification_type"`
	VerifierID                  *int64    `json:"verifier_id"`
	VerificationNotes           *string   `json:"verification_notes"`
	VerificationResult          string    `json:"verification_result"`
	RejectionReason             *string   `json:"rejection_reason"`
	AIConfidenceScore           *float64  `json:"ai_confidence_score"`
	AIVerificationData          *JSONMap  `gorm:"type:json" json:"ai_verification_data"`
	ThirdPartyResult            *JSONMap  `gorm:"type:json" json:"third_party_result"`
	VerificationDurationMinutes *int      `json:"verification_duration_minutes"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

// DocumentCompliance tracks compliance status
type DocumentCompliance struct {
	ID                    int64      `gorm:"primaryKey" json:"id"`
	TenantID              int64      `json:"tenant_id"`
	DocumentID            int64      `json:"document_id"`
	EntityType            string     `json:"entity_type"`
	EntityID              int64      `json:"entity_id"`
	ComplianceRequirement string     `json:"compliance_requirement"`
	IsCompliant           bool       `json:"is_compliant"`
	LastCheckedAt         *time.Time `json:"last_checked_at"`
	ComplianceNotes       *string    `json:"compliance_notes"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// DocumentAuditLog tracks document actions
type DocumentAuditLog struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	TenantID      int64     `json:"tenant_id"`
	DocumentID    int64     `json:"document_id"`
	Action        string    `json:"action"`
	ActionBy      *int64    `json:"action_by"`
	ActionDetails *string   `json:"action_details"`
	IPAddress     *string   `json:"ip_address"`
	UserAgent     *string   `json:"user_agent"`
	CreatedAt     time.Time `json:"created_at"`
}

// DocumentShare represents a document share
type DocumentShare struct {
	ID               int64      `gorm:"primaryKey" json:"id"`
	TenantID         int64      `json:"tenant_id"`
	DocumentID       int64      `json:"document_id"`
	SharedWithUserID *int64     `json:"shared_with_user_id"`
	SharedWithRole   *string    `json:"shared_with_role"`
	SharePermission  string     `json:"share_permission"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	IsActive         bool       `json:"is_active"`
	SharedBy         *int64     `json:"shared_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// DocumentCollection represents a collection of documents
type DocumentCollection struct {
	ID                    int64      `gorm:"primaryKey" json:"id"`
	TenantID              int64      `json:"tenant_id"`
	CollectionName        string     `json:"collection_name"`
	CollectionDescription *string    `json:"collection_description"`
	CollectionType        string     `json:"collection_type"`
	EntityType            string     `json:"entity_type"`
	EntityID              int64      `json:"entity_id"`
	Status                string     `json:"status"`
	CompletionPercentage  int        `json:"completion_percentage"`
	CreatedBy             *int64     `json:"created_by"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at"`
}

// DocumentCollectionItem represents an item in a collection
type DocumentCollectionItem struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	CollectionID int64     `json:"collection_id"`
	DocumentID   int64     `json:"document_id"`
	IsMandatory  bool      `json:"is_mandatory"`
	Status       string    `json:"status"`
	AddedAt      time.Time `json:"added_at"`
}

// JSONMap is a custom type for JSON data
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// JSONArray is a custom type for JSON arrays
type JSONArray []interface{}

func (ja JSONArray) Value() (driver.Value, error) {
	return json.Marshal(ja)
}

// DTOs

// CreateDocumentRequest represents a request to create a document
type CreateDocumentRequest struct {
	DocumentTypeID int64      `json:"document_type_id"`
	EntityType     string     `json:"entity_type"`
	EntityID       int64      `json:"entity_id"`
	DocumentName   string     `json:"document_name"`
	DocumentURL    string     `json:"document_url"`
	FileName       string     `json:"file_name"`
	FileSize       *int64     `json:"file_size"`
	FileFormat     *string    `json:"file_format"`
	IssueDate      *time.Time `json:"issue_date"`
	ExpiryDate     *time.Time `json:"expiry_date"`
	IsPrimary      bool       `json:"is_primary"`
}

// UpdateDocumentRequest represents a request to update a document
type UpdateDocumentRequest struct {
	DocumentName string     `json:"document_name"`
	ExpiryDate   *time.Time `json:"expiry_date"`
	IsPrimary    bool       `json:"is_primary"`
}

// VerifyDocumentRequest represents a verification request
type VerifyDocumentRequest struct {
	VerificationResult string   `json:"verification_result"`
	VerificationNotes  *string  `json:"verification_notes"`
	RejectionReason    *string  `json:"rejection_reason"`
	AIConfidenceScore  *float64 `json:"ai_confidence_score"`
}

// CreateDocumentCollectionRequest represents a request to create a collection
type CreateDocumentCollectionRequest struct {
	CollectionName string `json:"collection_name"`
	CollectionType string `json:"collection_type"`
	EntityType     string `json:"entity_type"`
	EntityID       int64  `json:"entity_id"`
}

// DocumentSummaryResponse provides a summary of documents
type DocumentSummaryResponse struct {
	TotalDocuments   int64     `json:"total_documents"`
	VerifiedCount    int64     `json:"verified_count"`
	PendingCount     int64     `json:"pending_count"`
	RejectedCount    int64     `json:"rejected_count"`
	ExpiredCount     int64     `json:"expired_count"`
	ComplianceStatus bool      `json:"compliance_status"`
	LastUpdated      time.Time `json:"last_updated"`
}
