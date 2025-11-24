package models

import (
    "time"
)

type Agent struct {
    UserID             int       `json:"user_id" db:"user_id"`
    Status             string    `json:"status" db:"status"` // active, inactive
    Availability       string    `json:"availability" db:"availability"` // online, offline, busy
    Skills             []string  `json:"skills" db:"skills"` // JSON array
    MaxConcurrentCalls int       `json:"max_concurrent_calls" db:"max_concurrent_calls"`
    CurrentCalls       int       `json:"current_calls" db:"current_calls"`
    TotalCalls         int       `json:"total_calls" db:"total_calls"`
    AvgHandleTime      float64   `json:"avg_handle_time" db:"avg_handle_time"`
    SatisfactionScore  float64   `json:"satisfaction_score" db:"satisfaction_score"`
    LastActive         time.Time `json:"last_active" db:"last_active"`
    TenantID           string    `json:"tenant_id" db:"tenant_id"`
}

type AgentStats struct {
    OnlineAgents int `json:"online_agents"`
    BusyAgents   int `json:"busy_agents"`
    TotalAgents  int `json:"total_agents"`
}
