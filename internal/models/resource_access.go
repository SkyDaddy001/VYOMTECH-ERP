package models

import "time"

// ResourceAccess defines access to specific resources by users
type ResourceAccess struct {
	ID           string     `json:"id" db:"id"`
	TenantID     string     `json:"tenant_id" db:"tenant_id"`
	UserID       int64      `json:"user_id" db:"user_id"`
	ResourceType string     `json:"resource_type" db:"resource_type"` // "lead", "customer", "project", "employee"
	ResourceID   string     `json:"resource_id" db:"resource_id"`     // specific resource ID
	AccessLevel  string     `json:"access_level" db:"access_level"`   // "view", "edit", "delete", "admin"
	ExpiresAt    *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}

// TimeBasedPermission defines permissions with time windows
type TimeBasedPermission struct {
	ID            string    `json:"id" db:"id"`
	TenantID      string    `json:"tenant_id" db:"tenant_id"`
	RoleID        string    `json:"role_id" db:"role_id"`
	PermissionID  string    `json:"permission_id" db:"permission_id"`
	EffectiveFrom time.Time `json:"effective_from" db:"effective_from"`
	ExpiresAt     time.Time `json:"expires_at" db:"expires_at"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	Reason        *string   `json:"reason" db:"reason"`           // temporary access reason
	ApprovedBy    *string   `json:"approved_by" db:"approved_by"` // user ID who approved
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// FieldLevelPermission defines field visibility/editability by role
type FieldLevelPermission struct {
	ID          string    `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	RoleID      string    `json:"role_id" db:"role_id"`
	ModuleName  string    `json:"module_name" db:"module_name"` // "sales", "hr", "accounting"
	EntityName  string    `json:"entity_name" db:"entity_name"` // "Lead", "Customer", "Employee"
	FieldName   string    `json:"field_name" db:"field_name"`   // "salary", "phone", "email"
	CanView     bool      `json:"can_view" db:"can_view"`
	CanEdit     bool      `json:"can_edit" db:"can_edit"`
	IsMasked    bool      `json:"is_masked" db:"is_masked"`       // hide sensitive data
	MaskPattern *string   `json:"mask_pattern" db:"mask_pattern"` // "XXXX-XXXX" for phone
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RoleDelegation allows managers to delegate role creation to subordinates
type RoleDelegation struct {
	ID              string     `json:"id" db:"id"`
	TenantID        string     `json:"tenant_id" db:"tenant_id"`
	ParentRoleID    string     `json:"parent_role_id" db:"parent_role_id"`     // manager role
	SubRoleID       string     `json:"sub_role_id" db:"sub_role_id"`           // delegated role
	PermissionBound string     `json:"permission_bound" db:"permission_bound"` // max permissions delegator can assign
	DelegatedBy     string     `json:"delegated_by" db:"delegated_by"`         // user ID who set delegation
	IsActive        bool       `json:"is_active" db:"is_active"`
	ExpiresAt       *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// BulkPermissionLog tracks bulk permission operations
type BulkPermissionLog struct {
	ID             string    `json:"id" db:"id"`
	TenantID       string    `json:"tenant_id" db:"tenant_id"`
	AssignmentType string    `json:"assignment_type" db:"assignment_type"` // assign, revoke, update
	TargetType     string    `json:"target_type" db:"target_type"`         // user, role, resource
	TotalCount     int       `json:"total_count" db:"total_count"`
	SuccessCount   int       `json:"success_count" db:"success_count"`
	FailedCount    int       `json:"failed_count" db:"failed_count"`
	ExecutedBy     string    `json:"executed_by" db:"executed_by"`
	ExecutionDate  time.Time `json:"execution_date" db:"execution_date"`
	Details        string    `json:"details" db:"details"` // JSON details
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
