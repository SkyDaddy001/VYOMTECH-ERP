package models

import "time"

// RoleTemplate represents a pre-built role template
type RoleTemplate struct {
	ID               string    `json:"id" db:"id"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	Category         string    `json:"category" db:"category"` // e.g., "built-in", "custom"
	IsSystemTemplate bool      `json:"is_system_template" db:"is_system_template"`
	PermissionIDs    []string  `json:"permission_ids"`
	Metadata         string    `json:"metadata" db:"metadata"` // JSON metadata
	IsActive         bool      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// TemplateInstance represents a role created from a template
type TemplateInstance struct {
	ID             string    `json:"id" db:"id"`
	TenantID       string    `json:"tenant_id" db:"tenant_id"`
	TemplateID     string    `json:"template_id" db:"template_id"`
	RoleID         string    `json:"role_id" db:"role_id"`
	CreatedBy      int64     `json:"created_by" db:"created_by"`
	Customizations string    `json:"customizations" db:"customizations"` // JSON customizations
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// Pre-defined role templates
var DefaultTemplates = map[string]RoleTemplate{
	"admin": {
		Name:             "Admin",
		Description:      "Full system access with all permissions",
		Category:         "built-in",
		IsSystemTemplate: true,
		IsActive:         true,
	},
	"manager": {
		Name:             "Manager",
		Description:      "Management access for team and operations",
		Category:         "built-in",
		IsSystemTemplate: true,
		IsActive:         true,
	},
	"sales": {
		Name:             "Sales",
		Description:      "Sales team member with lead and customer access",
		Category:         "built-in",
		IsSystemTemplate: true,
		IsActive:         true,
	},
	"hr": {
		Name:             "HR Manager",
		Description:      "Human resources management access",
		Category:         "built-in",
		IsSystemTemplate: true,
		IsActive:         true,
	},
	"finance": {
		Name:             "Finance",
		Description:      "Finance and accounting access",
		Category:         "built-in",
		IsSystemTemplate: true,
		IsActive:         true,
	},
}
