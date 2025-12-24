package models

import (
	"encoding/json"
	"time"
)

// ===== SALES & PRESALES REP MODELS =====

// LeadSource represents where a lead came from
type LeadSource string

const (
	LeadSourceWebsite     LeadSource = "website"
	LeadSourceReferral    LeadSource = "referral"
	LeadSourceColdCall    LeadSource = "cold_call"
	LeadSourceEmail       LeadSource = "email"
	LeadSourceTrade       LeadSource = "trade_show"
	LeadSourcePartner     LeadSource = "partner"
	LeadSourceSocial      LeadSource = "social_media"
	LeadSourceAdvertising LeadSource = "advertising"
	LeadSourceOther       LeadSource = "other"
)

// LeadStatus represents the stage of a lead
type LeadStatus string

const (
	LeadStatusNew       LeadStatus = "new"
	LeadStatusQualified LeadStatus = "qualified"
	LeadStatusNurturing LeadStatus = "nurturing"
	LeadStatusConverted LeadStatus = "converted"
	LeadStatusLost      LeadStatus = "lost"
	LeadStatusOnHold    LeadStatus = "on_hold"
)

// OpportunityStage represents stages in sales pipeline
type OpportunityStage string

const (
	OpportunityStageLead          OpportunityStage = "lead"
	OpportunityStageQualified     OpportunityStage = "qualified"
	OpportunityStageDemonstration OpportunityStage = "demonstration"
	OpportunityStageProposal      OpportunityStage = "proposal"
	OpportunityStageNegotiation   OpportunityStage = "negotiation"
	OpportunityStageClosed        OpportunityStage = "closed"
	OpportunityStageClosedWon     OpportunityStage = "closed_won"
	OpportunityStageClosedLost    OpportunityStage = "closed_lost"
)

// Lead represents a sales lead
type Lead struct {
	ID              string          `db:"id" json:"id"`
	TenantID        string          `db:"tenant_id" json:"tenant_id"`
	AssignedToID    string          `db:"assigned_to_id" json:"assigned_to_id"`
	AssignedToName  string          `db:"assigned_to_name" json:"assigned_to_name"`
	FirstName       string          `db:"first_name" json:"first_name"`
	LastName        string          `db:"last_name" json:"last_name"`
	CompanyName     string          `db:"company_name" json:"company_name"`
	Email           string          `db:"email" json:"email"`
	Phone           string          `db:"phone" json:"phone"`
	MobilePhone     string          `db:"mobile_phone" json:"mobile_phone"`
	Website         string          `db:"website" json:"website"`
	Industry        string          `db:"industry" json:"industry"`
	CompanySize     string          `db:"company_size" json:"company_size"` // small, medium, large, enterprise
	Location        string          `db:"location" json:"location"`
	City            string          `db:"city" json:"city"`
	State           string          `db:"state" json:"state"`
	Country         string          `db:"country" json:"country"`
	PostalCode      string          `db:"postal_code" json:"postal_code"`
	Source          LeadSource      `db:"source" json:"source"`
	Status          LeadStatus      `db:"status" json:"status"`
	Budget          float64         `db:"budget" json:"budget"`
	Currency        string          `db:"currency" json:"currency"`
	Description     string          `db:"description" json:"description"`
	Rating          int             `db:"rating" json:"rating"` // 1-5 stars
	Metadata        json.RawMessage `db:"metadata" json:"metadata"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
	ConvertedAt     *time.Time      `db:"converted_at" json:"converted_at"`
	LastContactedAt *time.Time      `db:"last_contacted_at" json:"last_contacted_at"`
}

// Opportunity represents a sales opportunity/deal
type Opportunity struct {
	ID              string           `db:"id" json:"id"`
	TenantID        string           `db:"tenant_id" json:"tenant_id"`
	LeadID          string           `db:"lead_id" json:"lead_id"`
	AccountID       string           `db:"account_id" json:"account_id"` // Customer account
	AssignedToID    string           `db:"assigned_to_id" json:"assigned_to_id"`
	AssignedToName  string           `db:"assigned_to_name" json:"assigned_to_name"`
	Name            string           `db:"name" json:"name"`
	Description     string           `db:"description" json:"description"`
	Stage           OpportunityStage `db:"stage" json:"stage"`
	Amount          float64          `db:"amount" json:"amount"`
	Currency        string           `db:"currency" json:"currency"`
	CloseDate       time.Time        `db:"close_date" json:"close_date"`
	ExpectedRevenue float64          `db:"expected_revenue" json:"expected_revenue"`
	Probability     int              `db:"probability" json:"probability"` // 0-100
	Source          LeadSource       `db:"source" json:"source"`
	CompetitorInfo  string           `db:"competitor_info" json:"competitor_info"`
	NextAction      string           `db:"next_action" json:"next_action"`
	NextActionDate  *time.Time       `db:"next_action_date" json:"next_action_date"`
	CreatedAt       time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time        `db:"updated_at" json:"updated_at"`
	WonAt           *time.Time       `db:"won_at" json:"won_at"`
	LostAt          *time.Time       `db:"lost_at" json:"lost_at"`
	LostReason      string           `db:"lost_reason" json:"lost_reason"`
}

// Activity represents a sales activity (call, email, meeting)
type Activity struct {
	ID             string     `db:"id" json:"id"`
	TenantID       string     `db:"tenant_id" json:"tenant_id"`
	LeadID         *string    `db:"lead_id" json:"lead_id"`
	OpportunityID  *string    `db:"opportunity_id" json:"opportunity_id"`
	AccountID      *string    `db:"account_id" json:"account_id"`
	ActivityType   string     `db:"activity_type" json:"activity_type"` // call, email, meeting, task
	Subject        string     `db:"subject" json:"subject"`
	Description    string     `db:"description" json:"description"`
	CreatedByID    string     `db:"created_by_id" json:"created_by_id"`
	CreatedByName  string     `db:"created_by_name" json:"created_by_name"`
	AssignedToID   string     `db:"assigned_to_id" json:"assigned_to_id"`
	AssignedToName string     `db:"assigned_to_name" json:"assigned_to_name"`
	ActivityDate   time.Time  `db:"activity_date" json:"activity_date"`
	DueDate        *time.Time `db:"due_date" json:"due_date"`
	Status         string     `db:"status" json:"status"`     // completed, pending, cancelled
	Priority       string     `db:"priority" json:"priority"` // high, medium, low
	Duration       int        `db:"duration" json:"duration"` // in minutes
	Outcome        string     `db:"outcome" json:"outcome"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}

// ===== POST-SALES/CRM MODELS =====

// Account represents a customer account
type Account struct {
	ID              string    `db:"id" json:"id"`
	TenantID        string    `db:"tenant_id" json:"tenant_id"`
	Name            string    `db:"name" json:"name"`
	Website         string    `db:"website" json:"website"`
	Industry        string    `db:"industry" json:"industry"`
	CompanySize     string    `db:"company_size" json:"company_size"`
	BillingAddress  string    `db:"billing_address" json:"billing_address"`
	ShippingAddress string    `db:"shipping_address" json:"shipping_address"`
	Phone           string    `db:"phone" json:"phone"`
	Email           string    `db:"email" json:"email"`
	AccountManager  string    `db:"account_manager_id" json:"account_manager_id"`
	AnnualRevenue   float64   `db:"annual_revenue" json:"annual_revenue"`
	Employees       int       `db:"employees" json:"employees"`
	Rating          int       `db:"rating" json:"rating"` // 1-5
	Type            string    `db:"type" json:"type"`     // prospect, customer, partner
	Description     string    `db:"description" json:"description"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// Contact represents a contact at a customer account
type Contact struct {
	ID                string    `db:"id" json:"id"`
	TenantID          string    `db:"tenant_id" json:"tenant_id"`
	AccountID         string    `db:"account_id" json:"account_id"`
	FirstName         string    `db:"first_name" json:"first_name"`
	LastName          string    `db:"last_name" json:"last_name"`
	Title             string    `db:"title" json:"title"` // Job title
	Department        string    `db:"department" json:"department"`
	Email             string    `db:"email" json:"email"`
	Phone             string    `db:"phone" json:"phone"`
	MobilePhone       string    `db:"mobile_phone" json:"mobile_phone"`
	Role              string    `db:"role" json:"role"` // Decision maker, influencer, user
	IsPrimary         bool      `db:"is_primary" json:"is_primary"`
	LinkedIn          string    `db:"linkedin" json:"linkedin"`
	Twitter           string    `db:"twitter" json:"twitter"`
	PreferredLanguage string    `db:"preferred_language" json:"preferred_language"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// Interaction represents customer interactions (support, feedback)
type Interaction struct {
	ID              string     `db:"id" json:"id"`
	TenantID        string     `db:"tenant_id" json:"tenant_id"`
	AccountID       string     `db:"account_id" json:"account_id"`
	ContactID       *string    `db:"contact_id" json:"contact_id"`
	InteractionType string     `db:"interaction_type" json:"interaction_type"` // support, feedback, complaint
	Channel         string     `db:"channel" json:"channel"`                   // email, phone, chat, social
	Subject         string     `db:"subject" json:"subject"`
	Message         string     `db:"message" json:"message"`
	Status          string     `db:"status" json:"status"`     // open, resolved, closed
	Priority        string     `db:"priority" json:"priority"` // high, medium, low
	AssignedToID    string     `db:"assigned_to_id" json:"assigned_to_id"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	ResolvedAt      *time.Time `db:"resolved_at" json:"resolved_at"`
}

// ===== REQUEST/RESPONSE MODELS =====

type CreateLeadRequest struct {
	FirstName   string     `json:"first_name" binding:"required"`
	LastName    string     `json:"last_name" binding:"required"`
	CompanyName string     `json:"company_name" binding:"required"`
	Email       string     `json:"email" binding:"required,email"`
	Phone       string     `json:"phone" binding:"required"`
	Source      LeadSource `json:"source" binding:"required"`
	Budget      float64    `json:"budget"`
	Currency    string     `json:"currency" default:"INR"`
	Industry    string     `json:"industry"`
	Description string     `json:"description"`
}

type CreateOpportunityRequest struct {
	LeadID      string           `json:"lead_id" binding:"required"`
	Name        string           `json:"name" binding:"required"`
	Amount      float64          `json:"amount" binding:"required,gt=0"`
	Currency    string           `json:"currency" default:"INR"`
	CloseDate   time.Time        `json:"close_date" binding:"required"`
	Stage       OpportunityStage `json:"stage"`
	Probability int              `json:"probability"`
	Description string           `json:"description"`
}

type UpdateOpportunityStageRequest struct {
	Stage OpportunityStage `json:"stage" binding:"required"`
	Notes string           `json:"notes"`
}

type CreateActivityRequest struct {
	LeadID        *string   `json:"lead_id"`
	OpportunityID *string   `json:"opportunity_id"`
	ActivityType  string    `json:"activity_type" binding:"required"`
	Subject       string    `json:"subject" binding:"required"`
	Description   string    `json:"description"`
	ActivityDate  time.Time `json:"activity_date"`
	AssignedToID  string    `json:"assigned_to_id"`
	Priority      string    `json:"priority"`
}

// SalesDashboard aggregates sales metrics
type SalesDashboard struct {
	TotalLeads           int64          `json:"total_leads"`
	QualifiedLeads       int64          `json:"qualified_leads"`
	ConversionRate       float64        `json:"conversion_rate"`
	TotalOpportunities   int64          `json:"total_opportunities"`
	PipelineValue        float64        `json:"pipeline_value"`
	WonDeals             int64          `json:"won_deals"`
	WonValue             float64        `json:"won_value"`
	LostDeals            int64          `json:"lost_deals"`
	AvgDealSize          float64        `json:"avg_deal_size"`
	SalesCycle           int            `json:"sales_cycle"` // in days
	TopSalesReps         []RepMetrics   `json:"top_sales_reps"`
	LeadsBySource        map[string]int `json:"leads_by_source"`
	OpportunitiesByStage map[string]int `json:"opportunities_by_stage"`
}

type RepMetrics struct {
	RepID         string  `json:"rep_id"`
	RepName       string  `json:"rep_name"`
	LeadsCount    int     `json:"leads_count"`
	OpptCount     int     `json:"oppt_count"`
	WonDealsCount int     `json:"won_deals_count"`
	TotalWon      float64 `json:"total_won"`
	AvgDealSize   float64 `json:"avg_deal_size"`
	WinRate       float64 `json:"win_rate"`
}
