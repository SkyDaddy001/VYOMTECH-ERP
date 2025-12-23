package models

import (
	"time"
)

// TitleClearance represents a title clearance record
type TitleClearance struct {
	ID                  int64      `gorm:"primaryKey" json:"id"`
	TenantID            int64      `json:"tenant_id"`
	BookingID           int64      `json:"booking_id"`
	PropertyID          *int64     `json:"property_id"`
	Status              string     `json:"status"`         // pending, in_progress, cleared, issues_found, rejected, expired
	ClearanceType       string     `json:"clearance_type"` // full_clearance, encumbrance_check, mutation, legal_opinion, boundary_verification
	StartDate           *time.Time `json:"start_date"`
	TargetClearanceDate *time.Time `json:"target_clearance_date"`
	ActualClearanceDate *time.Time `json:"actual_clearance_date"`
	ClearancePercentage float64    `json:"clearance_percentage"`
	IssuesCount         int        `json:"issues_count"`
	ResolvedIssuesCount int        `json:"resolved_issues_count"`
	Priority            string     `json:"priority"` // low, normal, high, critical
	Notes               *string    `json:"notes"`
	CreatedBy           *int64     `json:"created_by"`
	UpdatedBy           *int64     `json:"updated_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

// TitleIssue represents a title issue or encumbrance
type TitleIssue struct {
	ID               int64      `gorm:"primaryKey" json:"id"`
	TenantID         int64      `json:"tenant_id"`
	ClearanceID      int64      `json:"clearance_id"`
	IssueType        string     `json:"issue_type"` // lien, encumbrance, dispute, boundary_issue, mortgage, legal_claim, tax_issue
	IssueTitle       string     `json:"issue_title"`
	IssueDescription *string    `json:"issue_description"`
	Severity         string     `json:"severity"` // low, medium, high, critical
	Status           string     `json:"status"`   // open, under_review, escalated, resolved, deferred, invalid
	ReportedDate     *time.Time `json:"reported_date"`
	SourceDocument   *string    `json:"source_document"`
	AffectedParties  *string    `json:"affected_parties"`
	ResolutionNotes  *string    `json:"resolution_notes"`
	ResolvedDate     *time.Time `json:"resolved_date"`
	ResolvedBy       *int64     `json:"resolved_by"`
	ResolutionMethod *string    `json:"resolution_method"` // legal_opinion, court_order, mutual_agreement, rectification, insurance
	Metadata         *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy        *int64     `json:"created_by"`
	UpdatedBy        *int64     `json:"updated_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

// TitleSearchReport represents a title search report
type TitleSearchReport struct {
	ID                 int64      `gorm:"primaryKey" json:"id"`
	TenantID           int64      `json:"tenant_id"`
	ClearanceID        int64      `json:"clearance_id"`
	SearchType         string     `json:"search_type"` // government_records, registry_search, municipal_search, court_search, property_search
	SearchDate         *time.Time `json:"search_date"`
	SearchAuthority    *string    `json:"search_authority"`
	SearchReferenceNum *string    `json:"search_reference_number"`
	ReportURL          *string    `json:"report_url"`
	ReportFileName     *string    `json:"report_file_name"`
	ReportFileSize     *int64     `json:"report_file_size"`
	S3Bucket           *string    `json:"s3_bucket"`
	S3Key              *string    `json:"s3_key"`
	EncumbrancesFound  int        `json:"encumbrances_found"`
	SearchStatus       string     `json:"search_status"` // pending, completed, issues_found, verified
	VerifiedBy         *int64     `json:"verified_by"`
	VerifiedAt         *time.Time `json:"verified_at"`
	VerificationNotes  *string    `json:"verification_notes"`
	SearchCost         float64    `json:"search_cost"`
	Metadata           *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy          *int64     `json:"created_by"`
	UpdatedBy          *int64     `json:"updated_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
}

// TitleVerificationChecklist represents a verification checklist item
type TitleVerificationChecklist struct {
	ID                int64      `gorm:"primaryKey" json:"id"`
	TenantID          int64      `json:"tenant_id"`
	ClearanceID       int64      `json:"clearance_id"`
	ItemName          string     `json:"item_name"`
	ItemCategory      *string    `json:"item_category"` // ownership, encumbrance, mutation, boundary, legal, tax, litigation
	Description       *string    `json:"description"`
	IsMandatory       bool       `json:"is_mandatory"`
	Status            string     `json:"status"` // pending, verified, not_applicable, issue_found
	VerifiedBy        *int64     `json:"verified_by"`
	VerifiedAt        *time.Time `json:"verified_at"`
	VerificationNotes *string    `json:"verification_notes"`
	SequenceOrder     int        `json:"sequence_order"`
	Metadata          *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy         *int64     `json:"created_by"`
	UpdatedBy         *int64     `json:"updated_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

// TitleLegalOpinion represents a legal opinion or expert review
type TitleLegalOpinion struct {
	ID                  int64      `gorm:"primaryKey" json:"id"`
	TenantID            int64      `json:"tenant_id"`
	ClearanceID         int64      `json:"clearance_id"`
	OpinionType         string     `json:"opinion_type"` // legal_opinion, expert_review, boundary_survey, environmental_review
	ExpertName          *string    `json:"expert_name"`
	ExpertOrganization  *string    `json:"expert_organization"`
	ExpertLicenseNumber *string    `json:"expert_license_number"`
	OpinionDate         *time.Time `json:"opinion_date"`
	OpinionStatus       string     `json:"opinion_status"` // pending, received, under_review, approved, concerns_noted, rejected
	OpinionURL          *string    `json:"opinion_url"`
	OpinionFileName     *string    `json:"opinion_file_name"`
	FileSize            *int64     `json:"file_size"`
	S3Bucket            *string    `json:"s3_bucket"`
	S3Key               *string    `json:"s3_key"`
	OpinionSummary      *string    `json:"opinion_summary"`
	Recommendations     *string    `json:"recommendations"`
	RiskAssessment      *string    `json:"risk_assessment"` // low_risk, medium_risk, high_risk
	ReviewByLawyer      *int64     `json:"review_by_lawyer"`
	ReviewNotes         *string    `json:"review_notes"`
	ReviewedAt          *time.Time `json:"reviewed_at"`
	Cost                float64    `json:"cost"`
	Metadata            *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy           *int64     `json:"created_by"`
	UpdatedBy           *int64     `json:"updated_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

// TitleClearanceApproval represents an approval in the clearance process
type TitleClearanceApproval struct {
	ID                      int64      `gorm:"primaryKey" json:"id"`
	TenantID                int64      `json:"tenant_id"`
	ClearanceID             int64      `json:"clearance_id"`
	ApprovalType            string     `json:"approval_type"` // issue_resolution_approval, legal_opinion_approval, final_clearance_approval
	ApproverID              int64      `json:"approver_id"`
	ApproverRole            *string    `json:"approver_role"`
	ApprovalStatus          string     `json:"approval_status"` // pending, approved, rejected, conditional
	ApprovalNotes           *string    `json:"approval_notes"`
	ConditionalRequirements *string    `json:"conditional_requirements"`
	ApprovalDate            *time.Time `json:"approval_date"`
	ValidFrom               *time.Time `json:"valid_from"`
	ValidTill               *time.Time `json:"valid_till"`
	SequenceOrder           int        `json:"sequence_order"`
	IsFinalApproval         bool       `json:"is_final_approval"`
	Metadata                *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy               *int64     `json:"created_by"`
	UpdatedBy               *int64     `json:"updated_by"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
	DeletedAt               *time.Time `json:"deleted_at"`
}

// TitleClearanceAuditLog represents audit logs for title clearance actions
type TitleClearanceAuditLog struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	TenantID      int64     `json:"tenant_id"`
	ClearanceID   int64     `json:"clearance_id"`
	Action        string    `json:"action"`
	EntityType    string    `json:"entity_type"` // clearance, issue, search_report, legal_opinion, approval
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

// CreateTitleClearanceRequest represents a request to create a title clearance
type CreateTitleClearanceRequest struct {
	BookingID           int64      `json:"booking_id" binding:"required"`
	PropertyID          *int64     `json:"property_id"`
	ClearanceType       string     `json:"clearance_type" binding:"required"`
	TargetClearanceDate *time.Time `json:"target_clearance_date"`
	Priority            string     `json:"priority"`
	Notes               *string    `json:"notes"`
}

// UpdateTitleClearanceRequest represents a request to update a title clearance
type UpdateTitleClearanceRequest struct {
	Status              *string    `json:"status"`
	TargetClearanceDate *time.Time `json:"target_clearance_date"`
	ClearancePercentage *float64   `json:"clearance_percentage"`
	Priority            *string    `json:"priority"`
	Notes               *string    `json:"notes"`
}

// CreateTitleIssueRequest represents a request to create a title issue
type CreateTitleIssueRequest struct {
	IssueType        string  `json:"issue_type" binding:"required"`
	IssueTitle       string  `json:"issue_title" binding:"required"`
	IssueDescription *string `json:"issue_description"`
	Severity         string  `json:"severity" binding:"required"`
	SourceDocument   *string `json:"source_document"`
}

// UpdateTitleIssueRequest represents a request to update a title issue
type UpdateTitleIssueRequest struct {
	Status           *string `json:"status"`
	Severity         *string `json:"severity"`
	ResolutionNotes  *string `json:"resolution_notes"`
	ResolutionMethod *string `json:"resolution_method"`
}

// ResolveTitleIssueRequest represents a request to resolve a title issue
type ResolveTitleIssueRequest struct {
	Status           string  `json:"status" binding:"required"`
	ResolutionMethod string  `json:"resolution_method" binding:"required"`
	ResolutionNotes  *string `json:"resolution_notes"`
}

// CreateSearchReportRequest represents a request to create a search report
type CreateSearchReportRequest struct {
	SearchType         string   `json:"search_type" binding:"required"`
	SearchAuthority    *string  `json:"search_authority"`
	SearchReferenceNum *string  `json:"search_reference_number"`
	SearchCost         *float64 `json:"search_cost"`
}

// VerifySearchReportRequest represents a request to verify a search report
type VerifySearchReportRequest struct {
	SearchStatus      string  `json:"search_status" binding:"required"`
	EncumbrancesFound *int    `json:"encumbrances_found"`
	VerificationNotes *string `json:"verification_notes"`
}

// CreateLegalOpinionRequest represents a request to create a legal opinion
type CreateLegalOpinionRequest struct {
	OpinionType         string  `json:"opinion_type" binding:"required"`
	ExpertName          *string `json:"expert_name"`
	ExpertOrganization  *string `json:"expert_organization"`
	ExpertLicenseNumber *string `json:"expert_license_number"`
}

// ReviewLegalOpinionRequest represents a request to review a legal opinion
type ReviewLegalOpinionRequest struct {
	OpinionStatus  string  `json:"opinion_status" binding:"required"`
	RiskAssessment string  `json:"risk_assessment" binding:"required"`
	ReviewNotes    *string `json:"review_notes"`
}

// CreateClearanceApprovalRequest represents a request to create an approval
type CreateClearanceApprovalRequest struct {
	ApprovalType    string `json:"approval_type" binding:"required"`
	ApproverID      int64  `json:"approver_id" binding:"required"`
	SequenceOrder   int    `json:"sequence_order"`
	IsFinalApproval bool   `json:"is_final_approval"`
}

// ApproveClearanceRequest represents a request to approve a clearance
type ApproveClearanceRequest struct {
	ApprovalStatus          string  `json:"approval_status" binding:"required"`
	ApprovalNotes           *string `json:"approval_notes"`
	ConditionalRequirements *string `json:"conditional_requirements"`
	IsFinalApproval         bool    `json:"is_final_approval"`
}

// TitleClearanceSummaryResponse represents a summary of title clearance status
type TitleClearanceSummaryResponse struct {
	ClearanceID            int64      `json:"clearance_id"`
	BookingID              int64      `json:"booking_id"`
	Status                 string     `json:"status"`
	ClearanceType          string     `json:"clearance_type"`
	ClearancePercentage    float64    `json:"clearance_percentage"`
	TargetClearanceDate    *time.Time `json:"target_clearance_date"`
	ActualClearanceDate    *time.Time `json:"actual_clearance_date"`
	TotalIssues            int        `json:"total_issues"`
	OpenIssues             int        `json:"open_issues"`
	ResolvedIssues         int        `json:"resolved_issues"`
	SearchReports          int        `json:"search_reports"`
	LegalOpinions          int        `json:"legal_opinions"`
	PendingApprovals       int        `json:"pending_approvals"`
	ChecklistItems         int        `json:"checklist_items"`
	VerifiedChecklistItems int        `json:"verified_checklist_items"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}
