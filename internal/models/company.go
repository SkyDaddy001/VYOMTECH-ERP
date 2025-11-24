package models

import (
	"time"
)

// Company represents an organization within a tenant
type Company struct {
	ID                  string    `json:"id" db:"id"`
	TenantID            string    `json:"tenant_id" db:"tenant_id"`
	Name                string    `json:"name" db:"name"`
	Description         string    `json:"description" db:"description"`
	Status              string    `json:"status" db:"status"` // "active", "inactive", "suspended"
	IndustryType        string    `json:"industry_type" db:"industry_type"`
	EmployeeCount       *int      `json:"employee_count" db:"employee_count"`
	Website             *string   `json:"website" db:"website"`
	MaxProjects         int       `json:"max_projects" db:"max_projects"`
	MaxUsers            int       `json:"max_users" db:"max_users"`
	CurrentUserCount    int       `json:"current_user_count" db:"current_user_count"`
	CurrentProjectCount int       `json:"current_project_count" db:"current_project_count"`
	BillingEmail        string    `json:"billing_email" db:"billing_email"`
	BillingAddress      string    `json:"billing_address" db:"billing_address"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// Project represents a project under a company
type Project struct {
	ID               string     `json:"id" db:"id"`
	CompanyID        string     `json:"company_id" db:"company_id"`
	TenantID         string     `json:"tenant_id" db:"tenant_id"`
	Name             string     `json:"name" db:"name"`
	Description      string     `json:"description" db:"description"`
	Status           string     `json:"status" db:"status"`             // "active", "inactive", "archived"
	ProjectType      string     `json:"project_type" db:"project_type"` // "sales", "support", "marketing", "custom"
	MaxUsers         int        `json:"max_users" db:"max_users"`
	CurrentUserCount int        `json:"current_user_count" db:"current_user_count"`
	BudgetAllocated  float64    `json:"budget_allocated" db:"budget_allocated"`
	BudgetSpent      float64    `json:"budget_spent" db:"budget_spent"`
	StartDate        time.Time  `json:"start_date" db:"start_date"`
	EndDate          *time.Time `json:"end_date" db:"end_date"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// CompanyMember represents a user's membership in a company
type CompanyMember struct {
	ID         string    `json:"id" db:"id"`
	CompanyID  string    `json:"company_id" db:"company_id"`
	UserID     int       `json:"user_id" db:"user_id"`
	TenantID   string    `json:"tenant_id" db:"tenant_id"`
	Role       string    `json:"role" db:"role"` // "owner", "admin", "manager", "member", "viewer"
	Department string    `json:"department" db:"department"`
	IsActive   bool      `json:"is_active" db:"is_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// ProjectMember represents a user's membership in a project
type ProjectMember struct {
	ID        string    `json:"id" db:"id"`
	ProjectID string    `json:"project_id" db:"project_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	CompanyID string    `json:"company_id" db:"company_id"`
	TenantID  string    `json:"tenant_id" db:"tenant_id"`
	Role      string    `json:"role" db:"role"` // "lead", "member", "viewer", "analyst"
	JoinedAt  time.Time `json:"joined_at" db:"joined_at"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserRole represents a role and its permissions (cross-company/project)
type UserRole struct {
	ID          string    `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name"` // "admin", "accountant", "manager"
	Description string    `json:"description" db:"description"`
	Permissions string    `json:"permissions" db:"permissions"` // JSON array of permission strings
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
