package models

import (
	"time"
)

// Role represents a user role with permissions
type Role struct {
	ID          int64     `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name"` // admin, manager, agent, supervisor
	Description string    `json:"description" db:"description"`
	Permissions []string  `json:"permissions" db:"permissions"` // JSON array
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Permission represents an action that can be granted to a role
type Permission struct {
	ID          int64     `json:"id" db:"id"`
	Code        string    `json:"code" db:"code"` // e.g., "leads.view", "calls.edit", "reports.export"
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"` // leads, calls, agents, campaigns, reports
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// RolePermission maps roles to permissions
type RolePermission struct {
	ID           int64     `json:"id" db:"id"`
	RoleID       int64     `json:"role_id" db:"role_id"`
	PermissionID int64     `json:"permission_id" db:"permission_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID        int64     `json:"id" db:"id"`
	TenantID  string    `json:"tenant_id" db:"tenant_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"`     // CREATE, READ, UPDATE, DELETE, LOGIN, LOGOUT
	Resource  string    `json:"resource" db:"resource"` // lead, call, campaign, user, etc.
	Details   string    `json:"details" db:"details"`   // JSON details of the action
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Status    string    `json:"status" db:"status"` // success, failure
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// DataEncryption represents encrypted sensitive data
type DataEncryption struct {
	ID           int64     `json:"id" db:"id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	FieldName    string    `json:"field_name" db:"field_name"`       // phone, ssn, account_number
	FieldValue   string    `json:"field_value" db:"field_value"`     // encrypted value
	ResourceType string    `json:"resource_type" db:"resource_type"` // lead, agent, user
	ResourceID   int64     `json:"resource_id" db:"resource_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// GDPRRequest represents a GDPR data request
type GDPRRequest struct {
	ID        int64     `json:"id" db:"id"`
	TenantID  string    `json:"tenant_id" db:"tenant_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Type      string    `json:"type" db:"type"`     // access, deletion, portability
	Status    string    `json:"status" db:"status"` // pending, approved, rejected, completed
	Reason    string    `json:"reason" db:"reason"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"` // deletion request expires in 30 days
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ConsentRecord represents user consent for data processing
type ConsentRecord struct {
	ID          int64      `json:"id" db:"id"`
	TenantID    string     `json:"tenant_id" db:"tenant_id"`
	UserID      int64      `json:"user_id" db:"user_id"`
	Type        string     `json:"type" db:"type"` // marketing, analytics, third_party
	Given       bool       `json:"given" db:"given"`
	ConsentedAt time.Time  `json:"consented_at" db:"consented_at"`
	ExpiresAt   *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// DataClassification represents how sensitive data should be handled
type DataClassification struct {
	ID                 int64     `json:"id" db:"id"`
	FieldName          string    `json:"field_name" db:"field_name"`         // phone, email, ssn, account_number
	Classification     string    `json:"classification" db:"classification"` // public, internal, confidential, restricted
	EncryptionRequired bool      `json:"encryption_required" db:"encryption_required"`
	AuditRequired      bool      `json:"audit_required" db:"audit_required"`
	RetentionDays      int       `json:"retention_days" db:"retention_days"` // -1 for indefinite
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// SecurityEvent represents security-related events
type SecurityEvent struct {
	ID          int64      `json:"id" db:"id"`
	TenantID    string     `json:"tenant_id" db:"tenant_id"`
	UserID      *int64     `json:"user_id" db:"user_id"`
	EventType   string     `json:"event_type" db:"event_type"` // failed_login, suspicious_activity, permission_denied
	Severity    string     `json:"severity" db:"severity"`     // low, medium, high, critical
	Description string     `json:"description" db:"description"`
	IPAddress   string     `json:"ip_address" db:"ip_address"`
	ResolvedAt  *time.Time `json:"resolved_at" db:"resolved_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}
