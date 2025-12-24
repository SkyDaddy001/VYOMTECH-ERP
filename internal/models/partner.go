package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// PartnerType represents the type of external partner
type PartnerType string

const (
	PartnerTypePortal         PartnerType = "portal"          // White-label portal
	PartnerTypeChannelPartner PartnerType = "channel_partner" // Channel/reseller partner
	PartnerTypeVendor         PartnerType = "vendor"          // Vendor/supplier
	PartnerTypeCustomer       PartnerType = "customer"        // Direct customer/organization
)

// PartnerStatus represents the status of a partner
type PartnerStatus string

const (
	PartnerStatusPending   PartnerStatus = "pending"   // Awaiting approval
	PartnerStatusActive    PartnerStatus = "active"    // Approved and active
	PartnerStatusInactive  PartnerStatus = "inactive"  // Deactivated
	PartnerStatusSuspended PartnerStatus = "suspended" // Suspended for violations
	PartnerStatusRejected  PartnerStatus = "rejected"  // Rejected during approval
)

// Partner represents an external partner organization
type Partner struct {
	ID                   string         `json:"id" db:"id"`
	TenantID             string         `json:"tenant_id" db:"tenant_id"`
	PartnerCode          string         `json:"partner_code" db:"partner_code"` // Unique identifier
	OrganizationName     string         `json:"organization_name" db:"organization_name"`
	PartnerType          PartnerType    `json:"partner_type" db:"partner_type"` // portal, channel_partner, vendor, customer
	Status               PartnerStatus  `json:"status" db:"status"`             // pending, active, inactive, suspended
	ContactEmail         string         `json:"contact_email" db:"contact_email"`
	ContactPhone         string         `json:"contact_phone" db:"contact_phone"`
	ContactPerson        string         `json:"contact_person" db:"contact_person"`
	Website              string         `json:"website" db:"website"`
	Description          string         `json:"description" db:"description"`
	Address              string         `json:"address" db:"address"`
	City                 string         `json:"city" db:"city"`
	State                string         `json:"state" db:"state"`
	Country              string         `json:"country" db:"country"`
	ZipCode              string         `json:"zip_code" db:"zip_code"`
	TaxID                string         `json:"tax_id" db:"tax_id"`
	BankingDetails       BankingDetails `json:"banking_details" db:"banking_details"`             // JSON
	CommissionPercentage float64        `json:"commission_percentage" db:"commission_percentage"` // % of lead value
	LeadPrice            float64        `json:"lead_price" db:"lead_price"`                       // Fixed price per lead
	MonthlyQuota         int            `json:"monthly_quota" db:"monthly_quota"`                 // Max leads per month
	CurrentMonthLeads    int            `json:"current_month_leads" db:"current_month_leads"`
	TotalLeadsSubmitted  int64          `json:"total_leads_submitted" db:"total_leads_submitted"`
	ApprovedLeads        int64          `json:"approved_leads" db:"approved_leads"`
	RejectedLeads        int64          `json:"rejected_leads" db:"rejected_leads"`
	ConvertedLeads       int64          `json:"converted_leads" db:"converted_leads"`
	TotalEarnings        float64        `json:"total_earnings" db:"total_earnings"`
	PendingPayoutAmount  float64        `json:"pending_payout_amount" db:"pending_payout_amount"`
	WithdrawnAmount      float64        `json:"withdrawn_amount" db:"withdrawn_amount"`
	AvailableBalance     float64        `json:"available_balance" db:"available_balance"`
	ApprovedBy           *string        `json:"approved_by" db:"approved_by"`
	ApprovedAt           *time.Time     `json:"approved_at" db:"approved_at"`
	RejectionReason      string         `json:"rejection_reason" db:"rejection_reason"`
	SuspensionReason     string         `json:"suspension_reason" db:"suspension_reason"`
	SuspendedAt          *time.Time     `json:"suspended_at" db:"suspended_at"`
	DocumentURLs         DocumentURLs   `json:"document_urls" db:"document_urls"` // JSON - KYC docs
	CreatedBy            string         `json:"created_by" db:"created_by"`
	CreatedAt            time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time     `json:"deleted_at" db:"deleted_at"`
}

// BankingDetails stores partner banking information
type BankingDetails struct {
	BankName      string     `json:"bank_name"`
	AccountHolder string     `json:"account_holder"`
	AccountNumber string     `json:"account_number"`
	RoutingNumber string     `json:"routing_number"`
	IBAN          string     `json:"iban"`
	SWIFT         string     `json:"swift"`
	Currency      string     `json:"currency"`
	PaymentMethod string     `json:"payment_method"` // bank_transfer, paypal, check, wire
	IsVerified    bool       `json:"is_verified"`
	VerifiedAt    *time.Time `json:"verified_at"`
}

// Scan implements sql.Scanner interface
func (bd *BankingDetails) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, &bd)
}

// Value implements driver.Valuer interface
func (bd BankingDetails) Value() (driver.Value, error) {
	return json.Marshal(bd)
}

// DocumentURLs stores partner document URLs
type DocumentURLs struct {
	BusinessRegistration string `json:"business_registration"`
	TaxCertificate       string `json:"tax_certificate"`
	IdentityProof        string `json:"identity_proof"`
	AddressProof         string `json:"address_proof"`
	BankStatements       string `json:"bank_statements"`
	AgreementSigned      string `json:"agreement_signed"`
	Other                string `json:"other"`
}

// Scan implements sql.Scanner interface
func (du *DocumentURLs) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, &du)
}

// Value implements driver.Valuer interface
func (du DocumentURLs) Value() (driver.Value, error) {
	return json.Marshal(du)
}

// PartnerUser represents a user account for a partner
type PartnerUser struct {
	ID           string     `json:"id" db:"id"`
	PartnerID    string     `json:"partner_id" db:"partner_id"`
	TenantID     string     `json:"tenant_id" db:"tenant_id"`
	Email        string     `json:"email" db:"email"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	Phone        string     `json:"phone" db:"phone"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Role         string     `json:"role" db:"role"` // admin, lead_manager, viewer
	IsActive     bool       `json:"is_active" db:"is_active"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}

// PartnerLead represents a lead submitted by a partner
type PartnerLead struct {
	ID              string     `json:"id" db:"id"`
	TenantID        string     `json:"tenant_id" db:"tenant_id"`
	PartnerID       string     `json:"partner_id" db:"partner_id"`
	LeadID          *string    `json:"lead_id" db:"lead_id"`                 // FK to actual lead (if approved)
	SubmissionType  string     `json:"submission_type" db:"submission_type"` // new_lead, referral, import_batch
	Status          string     `json:"status" db:"status"`                   // submitted, under_review, approved, rejected, converted
	LeadData        LeadData   `json:"lead_data" db:"lead_data"`             // JSON - all lead info
	QualityScore    float64    `json:"quality_score" db:"quality_score"`     // Auto-calculated (0-100)
	RejectionReason string     `json:"rejection_reason" db:"rejection_reason"`
	ReviewedBy      *string    `json:"reviewed_by" db:"reviewed_by"`
	ReviewedAt      *time.Time `json:"reviewed_at" db:"reviewed_at"`
	SubmittedBy     string     `json:"submitted_by" db:"submitted_by"` // Partner user ID
	ConversionDate  *time.Time `json:"conversion_date" db:"conversion_date"`
	CreditAmount    float64    `json:"credit_amount" db:"credit_amount"` // Amount earned for this lead
	CreditStatus    string     `json:"credit_status" db:"credit_status"` // pending, approved, rejected, paid
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// LeadData stores submitted lead information
type LeadData struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	Company          string `json:"company"`
	Industry         string `json:"industry"`
	JobTitle         string `json:"job_title"`
	Address          string `json:"address"`
	City             string `json:"city"`
	State            string `json:"state"`
	Country          string `json:"country"`
	ZipCode          string `json:"zip_code"`
	LeadType         string `json:"lead_type"`      // prospect, customer, warm_lead
	BudgetRange      string `json:"budget_range"`   // low, medium, high
	TimelineDays     int    `json:"timeline_days"`  // Days to make decision
	InterestAreas    string `json:"interest_areas"` // CSV
	AdditionalInfo   string `json:"additional_info"`
	VerificationCode string `json:"verification_code"` // For verification
}

// Scan implements sql.Scanner interface
func (ld *LeadData) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, &ld)
}

// Value implements driver.Valuer interface
func (ld LeadData) Value() (driver.Value, error) {
	return json.Marshal(ld)
}

// PartnerPayout represents a payout approval/rejection
type PartnerPayout struct {
	ID              string     `json:"id" db:"id"`
	TenantID        string     `json:"tenant_id" db:"tenant_id"`
	PartnerID       string     `json:"partner_id" db:"partner_id"`
	PeriodStart     time.Time  `json:"period_start" db:"period_start"`
	PeriodEnd       time.Time  `json:"period_end" db:"period_end"`
	TotalLeadsCount int64      `json:"total_leads_count" db:"total_leads_count"`
	ApprovedLeads   int64      `json:"approved_leads" db:"approved_leads"`
	ConvertedLeads  int64      `json:"converted_leads" db:"converted_leads"`
	TotalAmount     float64    `json:"total_amount" db:"total_amount"`       // Before approval
	ApprovedAmount  float64    `json:"approved_amount" db:"approved_amount"` // After approval
	RejectedAmount  float64    `json:"rejected_amount" db:"rejected_amount"`
	Status          string     `json:"status" db:"status"`                 // pending, approved, rejected, paid, partially_paid
	PaymentMethod   string     `json:"payment_method" db:"payment_method"` // bank_transfer, paypal, check, wire
	PaymentDate     *time.Time `json:"payment_date" db:"payment_date"`
	ReferenceNumber string     `json:"reference_number" db:"reference_number"`
	ReviewedBy      *string    `json:"reviewed_by" db:"reviewed_by"`
	ApprovedBy      *string    `json:"approved_by" db:"approved_by"`
	ApprovedAt      *time.Time `json:"approved_at" db:"approved_at"`
	RejectionNotes  string     `json:"rejection_notes" db:"rejection_notes"`
	Notes           string     `json:"notes" db:"notes"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// PartnerPayoutDetail tracks individual lead credits in a payout
type PartnerPayoutDetail struct {
	ID               string    `json:"id" db:"id"`
	PayoutID         string    `json:"payout_id" db:"payout_id"`
	PartnerLeadID    string    `json:"partner_lead_id" db:"partner_lead_id"`
	LeadSubmissionID string    `json:"lead_submission_id" db:"lead_submission_id"`
	Amount           float64   `json:"amount" db:"amount"`
	Status           string    `json:"status" db:"status"` // approved, rejected
	ApprovalNotes    string    `json:"approval_notes" db:"approval_notes"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// PartnerLeadCredit represents credit approval for individual leads
type PartnerLeadCredit struct {
	ID              string     `json:"id" db:"id"`
	TenantID        string     `json:"tenant_id" db:"tenant_id"`
	PartnerLeadID   string     `json:"partner_lead_id" db:"partner_lead_id"`
	PartnerID       string     `json:"partner_id" db:"partner_id"`
	CreditAmount    float64    `json:"credit_amount" db:"credit_amount"`
	CalculationType string     `json:"calculation_type" db:"calculation_type"` // percentage, fixed_price
	Status          string     `json:"status" db:"status"`                     // pending_approval, approved, rejected
	ApprovedBy      *string    `json:"approved_by" db:"approved_by"`
	ApprovedAt      *time.Time `json:"approved_at" db:"approved_at"`
	RejectionReason string     `json:"rejection_reason" db:"rejection_reason"`
	Notes           string     `json:"notes" db:"notes"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// PartnerActivity logs partner activities
type PartnerActivity struct {
	ID         string    `json:"id" db:"id"`
	TenantID   string    `json:"tenant_id" db:"tenant_id"`
	PartnerID  string    `json:"partner_id" db:"partner_id"`
	UserID     *string   `json:"user_id" db:"user_id"`
	Action     string    `json:"action" db:"action"`     // lead_submitted, lead_approved, lead_rejected, payout_requested
	Resource   string    `json:"resource" db:"resource"` // lead, payout, partner_profile
	ResourceID string    `json:"resource_id" db:"resource_id"`
	Details    string    `json:"details" db:"details"` // JSON
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// PartnerFilter for filtering partners
type PartnerFilter struct {
	PartnerType PartnerType
	Status      PartnerStatus
	Search      string // Searches organization_name, partner_code, contact_email
	Limit       int
	Offset      int
}

// PartnerLeadFilter for filtering partner leads
type PartnerLeadFilter struct {
	PartnerID       int64
	Status          string
	SubmissionType  string
	QualityScoreMin float64
	QualityScoreMax float64
	Limit           int
	Offset          int
}

// PartnerStats contains statistics for a partner
type PartnerStats struct {
	TotalLeadsSubmitted int64   `json:"total_leads_submitted"`
	ApprovedLeads       int64   `json:"approved_leads"`
	RejectedLeads       int64   `json:"rejected_leads"`
	ConvertedLeads      int64   `json:"converted_leads"`
	ApprovalRate        float64 `json:"approval_rate"`   // %
	ConversionRate      float64 `json:"conversion_rate"` // %
	TotalEarnings       float64 `json:"total_earnings"`
	AvailableBalance    float64 `json:"available_balance"`
	PendingPayout       float64 `json:"pending_payout"`
	CurrentMonthLeads   int     `json:"current_month_leads"`
	MonthlyQuota        int     `json:"monthly_quota"`
	AverageLeadQuality  float64 `json:"average_lead_quality"`
}

// PartnerSourceType defines partner categorization for lead source tracking
type PartnerSourceType string

const (
	PartnerSourceCustomerReference PartnerSourceType = "customer_reference" // Direct customer referral
	PartnerSourceVendorReference   PartnerSourceType = "vendor_reference"   // Vendor referral
	PartnerSourceChannelPartner    PartnerSourceType = "channel_partner"    // Channel partner
	PartnerSourcePropertyPortal    PartnerSourceType = "property_portal"    // Property portal
)

// PartnerSource maps partner to lead source for tracking
type PartnerSource struct {
	ID             int64             `json:"id" db:"id"`
	TenantID       string            `json:"tenant_id" db:"tenant_id"`
	PartnerID      int64             `json:"partner_id" db:"partner_id"`
	SourceType     PartnerSourceType `json:"source_type" db:"source_type"` // customer_reference, vendor_reference, channel_partner, property_portal
	SourceCode     string            `json:"source_code" db:"source_code"` // e.g., "PORTAL_001"
	SourceName     string            `json:"source_name" db:"source_name"` // Display name
	Description    string            `json:"description" db:"description"`
	IsActive       bool              `json:"is_active" db:"is_active"`
	LeadsGenerated int64             `json:"leads_generated" db:"leads_generated"`
	LeadsConverted int64             `json:"leads_converted" db:"leads_converted"`
	TotalRevenue   float64           `json:"total_revenue" db:"total_revenue"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time        `json:"deleted_at" db:"deleted_at"`
}

// CreditPolicyType defines policy calculation methodology
type CreditPolicyType string

const (
	CreditPolicyTypeTimeBased     CreditPolicyType = "time_based"     // Monthly, quarterly, annual
	CreditPolicyTypeProjectBased  CreditPolicyType = "project_based"  // Per project/campaign
	CreditPolicyTypeCampaignBased CreditPolicyType = "campaign_based" // Per marketing campaign
)

// CreditPolicyCalculation defines how credits are calculated
type CreditPolicyCalculation string

const (
	CreditPolicyCalcPercentage   CreditPolicyCalculation = "percentage"    // % of lead value
	CreditPolicyCalcFixedPrice   CreditPolicyCalculation = "fixed_price"   // Fixed amount per lead
	CreditPolicyCalcTiered       CreditPolicyCalculation = "tiered"        // Volume-based tiers
	CreditPolicyCalcConversion   CreditPolicyCalculation = "conversion"    // Based on actual conversion
	CreditPolicyCalcRevenueshare CreditPolicyCalculation = "revenue_share" // % of deal revenue
)

// PartnerCreditPolicy defines credit calculation rules per partner
type PartnerCreditPolicy struct {
	ID              int64                   `json:"id" db:"id"`
	TenantID        string                  `json:"tenant_id" db:"tenant_id"`
	PartnerID       int64                   `json:"partner_id" db:"partner_id"`
	PolicyCode      string                  `json:"policy_code" db:"policy_code"` // e.g., "POLICY_2025_Q1"
	PolicyName      string                  `json:"policy_name" db:"policy_name"`
	PolicyType      CreditPolicyType        `json:"policy_type" db:"policy_type"` // time_based, project_based, campaign_based
	CalculationType CreditPolicyCalculation `json:"calculation_type" db:"calculation_type"`

	// Time-Based Policy
	TimeUnitType    string     `json:"time_unit_type" db:"time_unit_type"`   // monthly, quarterly, annual
	TimeUnitValue   int        `json:"time_unit_value" db:"time_unit_value"` // 1 for monthly, 3 for quarterly
	PolicyStartDate time.Time  `json:"policy_start_date" db:"policy_start_date"`
	PolicyEndDate   *time.Time `json:"policy_end_date" db:"policy_end_date"`

	// Project-Based Policy
	ProjectID   *int64 `json:"project_id" db:"project_id"` // FK to project
	ProjectName string `json:"project_name" db:"project_name"`

	// Campaign-Based Policy
	CampaignID   *int64 `json:"campaign_id" db:"campaign_id"` // FK to campaign
	CampaignName string `json:"campaign_name" db:"campaign_name"`

	// Credit Details
	BaseCredit      float64 `json:"base_credit" db:"base_credit"`           // Base amount or percentage
	MinimumCredit   float64 `json:"minimum_credit" db:"minimum_credit"`     // Floor value
	MaximumCredit   float64 `json:"maximum_credit" db:"maximum_credit"`     // Ceiling value
	BonusPercentage float64 `json:"bonus_percentage" db:"bonus_percentage"` // Bonus % if conditions met

	// Tier Configuration (JSON)
	TierConfig TierConfig `json:"tier_config" db:"tier_config"`

	// Conditions
	MinLeadQualityScore float64 `json:"min_lead_quality_score" db:"min_lead_quality_score"` // Only credits if quality >= this
	RequiresApproval    bool    `json:"requires_approval" db:"requires_approval"`           // Manual approval needed
	AutoApprove         bool    `json:"auto_approve" db:"auto_approve"`                     // Auto-approve if criteria met

	// Status
	IsActive         bool `json:"is_active" db:"is_active"`
	ApprovalRequired bool `json:"approval_required" db:"approval_required"`

	// Statistics
	TotalLeadsUnderPolicy int64   `json:"total_leads_under_policy" db:"total_leads_under_policy"`
	TotalCreditsAllocated float64 `json:"total_credits_allocated" db:"total_credits_allocated"`

	// Audit
	CreatedBy  int64      `json:"created_by" db:"created_by"`
	ApprovedBy *int64     `json:"approved_by" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at" db:"approved_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" db:"deleted_at"`
}

// TierConfig stores tiered credit configuration
type TierConfig struct {
	Tiers []CreditTier `json:"tiers"`
}

// CreditTier represents a volume-based credit tier
type CreditTier struct {
	TierLevel    int     `json:"tier_level"`    // 1, 2, 3, etc.
	MinLeads     int     `json:"min_leads"`     // Minimum leads in period
	MaxLeads     int     `json:"max_leads"`     // Maximum leads (or 0 for unlimited)
	CreditAmount float64 `json:"credit_amount"` // Credit per lead at this tier
	BonusPercent float64 `json:"bonus_percent"` // Bonus % at this tier
}

// Scan implements sql.Scanner interface
func (tc *TierConfig) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, &tc)
}

// Value implements driver.Valuer interface
func (tc TierConfig) Value() (driver.Value, error) {
	return json.Marshal(tc)
}

// PartnerCreditPolicyMapping maps leads to applicable policies
type PartnerCreditPolicyMapping struct {
	ID               int64     `json:"id" db:"id"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	PartnerLeadID    int64     `json:"partner_lead_id" db:"partner_lead_id"`
	PolicyID         int64     `json:"policy_id" db:"policy_id"`
	CalculatedCredit float64   `json:"calculated_credit" db:"calculated_credit"`
	Reason           string    `json:"reason" db:"reason"` // Why this policy was applied
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// PartnerSourceStats tracks performance by source
type PartnerSourceStats struct {
	SourceType            PartnerSourceType `json:"source_type"`
	SourceCode            string            `json:"source_code"`
	TotalLeads            int64             `json:"total_leads"`
	ApprovedLeads         int64             `json:"approved_leads"`
	ConvertedLeads        int64             `json:"converted_leads"`
	ApprovalRate          float64           `json:"approval_rate"`
	ConversionRate        float64           `json:"conversion_rate"`
	TotalCreditsAllocated float64           `json:"total_credits_allocated"`
	AverageQualityScore   float64           `json:"average_quality_score"`
}
