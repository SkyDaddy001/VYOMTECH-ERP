package models

import "time"

// Campaign represents a marketing or sales campaign
type Campaign struct {
	ID             int64     `json:"id" db:"id"`
	TenantID       string    `json:"tenant_id" db:"tenant_id"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	Status         string    `json:"status" db:"status"` // planned, active, paused, completed
	TargetLeads    int       `json:"target_leads" db:"target_leads"`
	GeneratedLeads int       `json:"generated_leads" db:"generated_leads"`
	ConvertedLeads int       `json:"converted_leads" db:"converted_leads"`
	Budget         float64   `json:"budget" db:"budget"`
	SpentBudget    float64   `json:"spent_budget" db:"spent_budget"`
	CostPerLead    float64   `json:"cost_per_lead" db:"cost_per_lead"`
	ConversionRate float64   `json:"conversion_rate" db:"conversion_rate"`
	StartDate      time.Time `json:"start_date" db:"start_date"`
	EndDate        time.Time `json:"end_date" db:"end_date"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// CampaignFilter for filtering campaigns
type CampaignFilter struct {
	Status string
	Limit  int
	Offset int
}

// CampaignStats contains campaign statistics
type CampaignStats struct {
	Total              int     `json:"total"`
	Active             int     `json:"active"`
	Completed          int     `json:"completed"`
	Paused             int     `json:"paused"`
	AverageConversion  float64 `json:"average_conversion_rate"`
	TotalBudget        float64 `json:"total_budget"`
	TotalSpent         float64 `json:"total_spent"`
	AverageCostPerLead float64 `json:"average_cost_per_lead"`
}

// CampaignPerformance tracks daily performance
type CampaignPerformance struct {
	ID             int64     `json:"id" db:"id"`
	CampaignID     int64     `json:"campaign_id" db:"campaign_id"`
	Date           time.Time `json:"date" db:"date"`
	LeadsGenerated int       `json:"leads_generated" db:"leads_generated"`
	Conversions    int       `json:"conversions" db:"conversions"`
	SpentBudget    float64   `json:"spent_budget" db:"spent_budget"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
