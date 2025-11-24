package models

import (
	"time"
)

type Tenant struct {
	ID                 string    `json:"id" db:"id"`
	Name               string    `json:"name" db:"name"`
	Domain             string    `json:"domain" db:"domain"`
	Status             string    `json:"status" db:"status"`
	MaxUsers           int       `json:"max_users" db:"max_users"`
	MaxConcurrentCalls int       `json:"max_concurrent_calls" db:"max_concurrent_calls"`
	AIBudgetMonthly    float64   `json:"ai_budget_monthly" db:"ai_budget_monthly"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}
