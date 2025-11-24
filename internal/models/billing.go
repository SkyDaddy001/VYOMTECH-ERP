package models

import (
	"database/sql/driver"
	"time"
)

// BillingCycleType represents the billing cycle frequency
type BillingCycleType string

const (
	BillingCycleMonthly   BillingCycleType = "monthly"
	BillingCycleQuarterly BillingCycleType = "quarterly"
	BillingCycleAnnual    BillingCycleType = "annual"
	BillingCycleCustom    BillingCycleType = "custom"
)

// Scan implements the sql.Scanner interface
func (bct BillingCycleType) Scan(value interface{}) error {
	*(*string)(&bct) = value.(string)
	return nil
}

// Value implements the driver.Valuer interface
func (bct BillingCycleType) Value() (driver.Value, error) {
	return string(bct), nil
}

// Invoice represents a billing invoice
type Invoice struct {
	ID                 string     `json:"id" db:"id"`
	TenantID           string     `json:"tenant_id" db:"tenant_id"`
	InvoiceNumber      string     `json:"invoice_number" db:"invoice_number"`
	BillingPeriodStart time.Time  `json:"billing_period_start" db:"billing_period_start"`
	BillingPeriodEnd   time.Time  `json:"billing_period_end" db:"billing_period_end"`
	SubtotalAmount     float64    `json:"subtotal_amount" db:"subtotal_amount"`
	TaxAmount          float64    `json:"tax_amount" db:"tax_amount"`
	DiscountAmount     float64    `json:"discount_amount" db:"discount_amount"`
	TotalAmount        float64    `json:"total_amount" db:"total_amount"`
	Status             string     `json:"status" db:"status"`                 // "draft", "pending", "paid", "overdue", "cancelled"
	PaymentMethod      string     `json:"payment_method" db:"payment_method"` // "credit_card", "bank_transfer", "invoice"
	IssuedAt           time.Time  `json:"issued_at" db:"issued_at"`
	DueAt              time.Time  `json:"due_at" db:"due_at"`
	PaidAt             *time.Time `json:"paid_at" db:"paid_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// InvoiceLineItem represents a line in an invoice
type InvoiceLineItem struct {
	ID          string    `json:"id" db:"id"`
	InvoiceID   string    `json:"invoice_id" db:"invoice_id"`
	ModuleID    *string   `json:"module_id" db:"module_id"`
	Description string    `json:"description" db:"description"`
	Quantity    int       `json:"quantity" db:"quantity"`
	UnitPrice   float64   `json:"unit_price" db:"unit_price"`
	TotalPrice  float64   `json:"total_price" db:"total_price"`
	TaxRate     float64   `json:"tax_rate" db:"tax_rate"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// UsageMetrics represents daily usage tracking for billing
type UsageMetrics struct {
	ID              string    `json:"id" db:"id"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	CompanyID       *string   `json:"company_id" db:"company_id"`
	ProjectID       *string   `json:"project_id" db:"project_id"`
	Date            time.Time `json:"date" db:"date"`
	ActiveUsers     int       `json:"active_users" db:"active_users"`
	NewUsers        int       `json:"new_users" db:"new_users"`
	APICallsUsed    int       `json:"api_calls_used" db:"api_calls_used"`
	StorageUsedMB   float64   `json:"storage_used_mb" db:"storage_used_mb"`
	ModuleUsageData string    `json:"module_usage_data" db:"module_usage_data"` // JSON
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// Billing represents the billing configuration for a tenant
type Billing struct {
	ID               string           `json:"id" db:"id"`
	TenantID         string           `json:"tenant_id" db:"tenant_id"`
	BillingEmail     string           `json:"billing_email" db:"billing_email"`
	BillingCycle     BillingCycleType `json:"billing_cycle" db:"billing_cycle"`
	NextBillingDate  time.Time        `json:"next_billing_date" db:"next_billing_date"`
	AutomaticPayment bool             `json:"automatic_payment" db:"automatic_payment"`
	PaymentMethodID  *string          `json:"payment_method_id" db:"payment_method_id"`
	TaxRate          float64          `json:"tax_rate" db:"tax_rate"`
	TaxID            *string          `json:"tax_id" db:"tax_id"`
	BillingAddress   string           `json:"billing_address" db:"billing_address"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
}

// PricingPlan represents a predefined pricing package
type PricingPlan struct {
	ID                    string    `json:"id" db:"id"`
	Name                  string    `json:"name" db:"name"` // "Startup", "Professional", "Enterprise"
	Description           string    `json:"description" db:"description"`
	MonthlyPrice          float64   `json:"monthly_price" db:"monthly_price"`
	AnnualPrice           float64   `json:"annual_price" db:"annual_price"`
	MaxUsers              int       `json:"max_users" db:"max_users"`
	MaxCompanies          int       `json:"max_companies" db:"max_companies"`
	MaxProjectsPerCompany int       `json:"max_projects_per_company" db:"max_projects_per_company"`
	IncludedModules       string    `json:"included_modules" db:"included_modules"`     // JSON array
	AdditionalModules     string    `json:"additional_modules" db:"additional_modules"` // JSON array
	Features              string    `json:"features" db:"features"`                     // JSON array
	SortOrder             int       `json:"sort_order" db:"sort_order"`
	IsActive              bool      `json:"is_active" db:"is_active"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// TenantPlanSubscription represents a tenant's subscription to a plan
type TenantPlanSubscription struct {
	ID              string           `json:"id" db:"id"`
	TenantID        string           `json:"tenant_id" db:"tenant_id"`
	PricingPlanID   string           `json:"pricing_plan_id" db:"pricing_plan_id"`
	StartDate       time.Time        `json:"start_date" db:"start_date"`
	NextBillingDate time.Time        `json:"next_billing_date" db:"next_billing_date"`
	EndDate         *time.Time       `json:"end_date" db:"end_date"`
	Status          string           `json:"status" db:"status"` // "active", "paused", "cancelled"
	BillingCycle    BillingCycleType `json:"billing_cycle" db:"billing_cycle"`
	IsAutoRenew     bool             `json:"is_auto_renew" db:"is_auto_renew"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at" db:"updated_at"`
}
