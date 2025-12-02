package models

import "time"

// Site represents a construction site
type Site struct {
	ID              int64  `gorm:"primaryKey"`
	TenantID        string `gorm:"index"`
	SiteName        string
	Location        string
	ProjectID       string
	SiteManager     string
	StartDate       time.Time
	ExpectedEndDate time.Time
	CurrentStatus   string // planning, active, paused, completed, closed
	SiteAreaSqm     float64
	WorkforceCount  int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// SafetyIncident represents a safety incident or near-miss
type SafetyIncident struct {
	ID             int64  `gorm:"primaryKey"`
	TenantID       string `gorm:"index"`
	SiteID         int64
	IncidentType   string // accident, near_miss, hazard, violation
	Severity       string // low, medium, high, critical
	IncidentDate   time.Time
	Description    string `gorm:"type:text"`
	ReportedBy     string
	Status         string // open, investigating, resolved, closed
	IncidentNumber string `gorm:"uniqueIndex"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Compliance represents compliance tracking record
type Compliance struct {
	ID             int64  `gorm:"primaryKey"`
	TenantID       string `gorm:"index"`
	SiteID         int64
	ComplianceType string // safety, environmental, labor, regulatory
	Requirement    string
	DueDate        time.Time
	Status         string // compliant, non_compliant, in_progress, not_applicable
	LastAuditDate  time.Time
	AuditResult    string // pass, fail, pending
	Notes          string `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Permit represents construction permits
type Permit struct {
	ID               int64  `gorm:"primaryKey"`
	TenantID         string `gorm:"index"`
	SiteID           int64
	PermitType       string
	PermitNumber     string `gorm:"uniqueIndex"`
	IssuedDate       time.Time
	ExpiryDate       time.Time
	IssuingAuthority string
	Status           string // active, expired, cancelled, pending
	DocumentURL      string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// TableName specifies the table name for Site
func (Site) TableName() string {
	return "sites"
}

// TableName specifies the table name for SafetyIncident
func (SafetyIncident) TableName() string {
	return "safety_incidents"
}

// TableName specifies the table name for Compliance
func (Compliance) TableName() string {
	return "compliance_records"
}

// TableName specifies the table name for Permit
func (Permit) TableName() string {
	return "permits"
}
