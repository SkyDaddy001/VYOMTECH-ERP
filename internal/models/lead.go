package models

import "time"

// Lead represents a potential customer or prospect
type Lead struct {
	ID                  string     `json:"id" db:"id"`
	TenantID            string     `json:"tenant_id" db:"tenant_id"`
	LeadCode            string     `json:"lead_code" db:"lead_code"`
	FirstName           string     `json:"first_name" db:"first_name"`
	LastName            string     `json:"last_name" db:"last_name"`
	Email               string     `json:"email" db:"email"`
	Phone               string     `json:"phone" db:"phone"`
	CompanyName         string     `json:"company_name" db:"company_name"`
	Industry            string     `json:"industry" db:"industry"`
	Status              string     `json:"status" db:"status"` // new, contacted, qualified, negotiation, converted, lost
	Probability         float64    `json:"probability" db:"probability"`
	Source              string     `json:"source" db:"source"`           // campaign, manual, import
	CampaignID          string     `json:"campaign_id" db:"campaign_id"` // FK to campaign
	AssignedTo          string     `json:"assigned_to" db:"assigned_to"` // FK to user
	AssignedDate        *time.Time `json:"assigned_date" db:"assigned_date"`
	ConvertedToCustomer bool       `json:"converted_to_customer" db:"converted_to_customer"`
	CustomerID          string     `json:"customer_id" db:"customer_id"`
	NextActionDate      *time.Time `json:"next_action_date" db:"next_action_date"`
	NextActionNotes     string     `json:"next_action_notes" db:"next_action_notes"`
	CreatedBy           string     `json:"created_by" db:"created_by"`
	GLCustomerAccountID string     `json:"gl_customer_account_id" db:"gl_customer_account_id"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

// LeadFilter for filtering leads
type LeadFilter struct {
	Status     string
	Source     string
	CampaignID int64
	AssignedTo int64
	Limit      int
	Offset     int
}

// LeadStats contains statistics for leads
type LeadStats struct {
	Total     int     `json:"total"`
	New       int     `json:"new"`
	Contacted int     `json:"contacted"`
	Qualified int     `json:"qualified"`
	Converted int     `json:"converted"`
	Lost      int     `json:"lost"`
	ConvRate  float64 `json:"conversion_rate"`
}
