package models

import "time"

// Call represents a phone or video call
type Call struct {
	ID              int64      `json:"id" db:"id"`
	TenantID        string     `json:"tenant_id" db:"tenant_id"`
	LeadID          int64      `json:"lead_id" db:"lead_id"`
	AgentID         int64      `json:"agent_id" db:"agent_id"`
	Status          string     `json:"status" db:"status"` // initiated, ringing, active, ended
	DurationSeconds int        `json:"duration_seconds" db:"duration_seconds"`
	RecordingURL    string     `json:"recording_url" db:"recording_url"`
	Notes           string     `json:"notes" db:"notes"`
	Outcome         string     `json:"outcome" db:"outcome"` // success, no_answer, busy, declined, failed
	StartedAt       *time.Time `json:"started_at" db:"started_at"`
	EndedAt         *time.Time `json:"ended_at" db:"ended_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// CallFilter for filtering calls
type CallFilter struct {
	Status  string
	Outcome string
	AgentID int64
	LeadID  int64
	Limit   int
	Offset  int
}

// CallStats contains call statistics
type CallStats struct {
	Total           int     `json:"total"`
	Active          int     `json:"active"`
	Completed       int     `json:"completed"`
	Failed          int     `json:"failed"`
	AverageDuration int     `json:"average_duration"`
	TotalDuration   int     `json:"total_duration"`
	SuccessRate     float64 `json:"success_rate"`
}

// CallLog tracks call history
type CallLog struct {
	ID        int64     `json:"id" db:"id"`
	CallID    int64     `json:"call_id" db:"call_id"`
	Event     string    `json:"event" db:"event"` // ringing, connected, transferred, ended
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Details   string    `json:"details" db:"details"`
}
