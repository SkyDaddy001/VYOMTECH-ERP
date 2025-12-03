package models

import (
	"time"
)

type Agent struct {
	ID                 string    `json:"id" db:"id"`
	TenantID           string    `json:"tenant_id" db:"tenant_id"`
	AgentCode          string    `json:"agent_code" db:"agent_code"`
	FirstName          string    `json:"first_name" db:"first_name"`
	LastName           string    `json:"last_name" db:"last_name"`
	Email              string    `json:"email" db:"email"`
	Phone              string    `json:"phone" db:"phone"`
	Status             string    `json:"status" db:"status"`         // available, offline, busy
	AgentType          string    `json:"agent_type" db:"agent_type"` // inbound, outbound, blended
	Skills             []string  `json:"skills" db:"skills"`         // JSON array
	MaxConcurrentCalls int       `json:"max_concurrent_calls" db:"max_concurrent_calls"`
	Available          bool      `json:"available" db:"available"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type AgentStats struct {
	OnlineAgents int `json:"online_agents"`
	BusyAgents   int `json:"busy_agents"`
	TotalAgents  int `json:"total_agents"`
}
