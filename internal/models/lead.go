package models

import "time"

// Lead represents a potential customer or prospect
type Lead struct {
	ID            int64     `json:"id" db:"id"`
	TenantID      string    `json:"tenant_id" db:"tenant_id"`
	Name          string    `json:"name" db:"name"`
	Email         string    `json:"email" db:"email"`
	Phone         string    `json:"phone" db:"phone"`
	Status        string    `json:"status" db:"status"`           // new, contacted, qualified, converted, lost
	Source        string    `json:"source" db:"source"`           // campaign, manual, import
	CampaignID    *int64    `json:"campaign_id" db:"campaign_id"` // FK to campaign
	AssignedAgent *int64    `json:"assigned_agent_id" db:"assigned_agent_id"`
	Notes         string    `json:"notes" db:"notes"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
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
