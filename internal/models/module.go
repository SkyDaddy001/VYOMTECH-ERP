package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ModuleStatus represents the status of a module
type ModuleStatus string

const (
	ModuleStatusActive     ModuleStatus = "active"
	ModuleStatusInactive   ModuleStatus = "inactive"
	ModuleStatusBeta       ModuleStatus = "beta"
	ModuleStatusDeprecated ModuleStatus = "deprecated"
)

// PricingModel represents how a module is priced
type PricingModel string

const (
	PricingModelFree       PricingModel = "free"
	PricingModelPerUser    PricingModel = "per_user"
	PricingModelPerProject PricingModel = "per_project"
	PricingModelPerCompany PricingModel = "per_company"
	PricingModelFlat       PricingModel = "flat"
	PricingModelTiered     PricingModel = "tiered"
)

// Module represents a feature module in the system
type Module struct {
	ID               string          `json:"id" db:"id"`     // Unique identifier (e.g., "gamification", "workflow", "analytics")
	Name             string          `json:"name" db:"name"` // Display name
	Description      string          `json:"description" db:"description"`
	Category         string          `json:"category" db:"category"` // "core", "analytics", "automation", "communication", "ai"
	Status           ModuleStatus    `json:"status" db:"status"`
	Version          string          `json:"version" db:"version"` // Current version
	PricingModel     PricingModel    `json:"pricing_model" db:"pricing_model"`
	BaseCost         float64         `json:"base_cost" db:"base_cost"` // Monthly base cost
	CostPerUser      float64         `json:"cost_per_user" db:"cost_per_user"`
	CostPerProject   float64         `json:"cost_per_project" db:"cost_per_project"`
	CostPerCompany   float64         `json:"cost_per_company" db:"cost_per_company"`
	MaxUsers         *int            `json:"max_users" db:"max_users"` // null = unlimited
	MaxProjects      *int            `json:"max_projects" db:"max_projects"`
	MaxCompanies     *int            `json:"max_companies" db:"max_companies"`
	IsDependentOn    json.RawMessage `json:"is_dependent_on" db:"is_dependent_on"` // Array of module IDs it depends on
	IsCore           bool            `json:"is_core" db:"is_core"`                 // Cannot be disabled
	RequiresApproval bool            `json:"requires_approval" db:"requires_approval"`
	TrialDaysAllowed int             `json:"trial_days_allowed" db:"trial_days_allowed"`
	Features         json.RawMessage `json:"features" db:"features"` // JSON array of feature flags
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

// Scan implements the sql.Scanner interface
func (ms ModuleStatus) Scan(value interface{}) error {
	*(*string)(&ms) = value.(string)
	return nil
}

// Value implements the driver.Valuer interface
func (ms ModuleStatus) Value() (driver.Value, error) {
	return string(ms), nil
}

// Scan implements the sql.Scanner interface
func (pm PricingModel) Scan(value interface{}) error {
	*(*string)(&pm) = value.(string)
	return nil
}

// Value implements the driver.Valuer interface
func (pm PricingModel) Value() (driver.Value, error) {
	return string(pm), nil
}

// ModuleSubscription represents a tenant's subscription to a module
type ModuleSubscription struct {
	ID                   string          `json:"id" db:"id"`
	TenantID             string          `json:"tenant_id" db:"tenant_id"`
	CompanyID            *string         `json:"company_id" db:"company_id"` // null = tenant-level subscription
	ProjectID            *string         `json:"project_id" db:"project_id"` // null = company/tenant-level subscription
	ModuleID             string          `json:"module_id" db:"module_id"`
	Status               string          `json:"status" db:"status"` // "active", "inactive", "trial", "suspended", "expired"
	SubscriptionStarted  time.Time       `json:"subscription_started" db:"subscription_started"`
	SubscriptionEnded    *time.Time      `json:"subscription_ended" db:"subscription_ended"`
	TrialStartedAt       *time.Time      `json:"trial_started_at" db:"trial_started_at"`
	TrialEndsAt          *time.Time      `json:"trial_ends_at" db:"trial_ends_at"`
	MaxUsersAllowed      *int            `json:"max_users_allowed" db:"max_users_allowed"`
	CurrentUserCount     int             `json:"current_user_count" db:"current_user_count"`
	MonthlyBudget        *float64        `json:"monthly_budget" db:"monthly_budget"`
	AmountSpentThisMonth float64         `json:"amount_spent_this_month" db:"amount_spent_this_month"`
	Configuration        json.RawMessage `json:"configuration" db:"configuration"` // Module-specific config
	IsEnabled            bool            `json:"is_enabled" db:"is_enabled"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
}

// ModuleUsage tracks usage metrics for monetization
type ModuleUsage struct {
	ID             string          `json:"id" db:"id"`
	SubscriptionID string          `json:"subscription_id" db:"subscription_id"`
	TenantID       string          `json:"tenant_id" db:"tenant_id"`
	CompanyID      *string         `json:"company_id" db:"company_id"`
	ProjectID      *string         `json:"project_id" db:"project_id"`
	ModuleID       string          `json:"module_id" db:"module_id"`
	UserCount      int             `json:"user_count" db:"user_count"`
	ProjectCount   int             `json:"project_count" db:"project_count"`
	CompanyCount   int             `json:"company_count" db:"company_count"`
	CustomMetrics  json.RawMessage `json:"custom_metrics" db:"custom_metrics"` // Module-specific usage
	UsageDate      time.Time       `json:"usage_date" db:"usage_date"`         // Daily snapshot
	EstimatedCost  float64         `json:"estimated_cost" db:"estimated_cost"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

// ModuleLicense represents license details for a tenant
type ModuleLicense struct {
	ID                    string          `json:"id" db:"id"`
	TenantID              string          `json:"tenant_id" db:"tenant_id"`
	LicenseKey            string          `json:"license_key" db:"license_key"`
	LicenseType           string          `json:"license_type" db:"license_type"` // "startup", "enterprise", "custom"
	MaxCompanies          int             `json:"max_companies" db:"max_companies"`
	MaxProjectsPerCompany int             `json:"max_projects_per_company" db:"max_projects_per_company"`
	MaxUsersPerProject    int             `json:"max_users_per_project" db:"max_users_per_project"`
	TotalMaxUsers         int             `json:"total_max_users" db:"total_max_users"`
	EnabledModules        json.RawMessage `json:"enabled_modules" db:"enabled_modules"` // Array of module IDs
	DisabledModules       json.RawMessage `json:"disabled_modules" db:"disabled_modules"`
	ExpiresAt             time.Time       `json:"expires_at" db:"expires_at"`
	IssuedAt              time.Time       `json:"issued_at" db:"issued_at"`
	IssuedBy              string          `json:"issued_by" db:"issued_by"`
	Status                string          `json:"status" db:"status"` // "active", "expired", "revoked"
	CreatedAt             time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at" db:"updated_at"`
}
