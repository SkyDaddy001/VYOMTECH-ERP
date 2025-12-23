package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/logger"
)

// RoleTemplateService handles role templates and template-based role creation
type RoleTemplateService struct {
	db     *sql.DB
	logger *logger.Logger
	rbac   *RBACService
}

// NewRoleTemplateService creates a new role template service
func NewRoleTemplateService(db *sql.DB, rbacService *RBACService, log *logger.Logger) *RoleTemplateService {
	return &RoleTemplateService{
		db:     db,
		logger: log,
		rbac:   rbacService,
	}
}

// InitializeDefaultTemplates creates default role templates for a tenant
func (ts *RoleTemplateService) InitializeDefaultTemplates(ctx context.Context, tenantID string) error {
	// Define default permissions for each template
	templatePermissions := map[string][]string{
		"admin": {
			"sales.leads.create", "sales.leads.read", "sales.leads.update", "sales.leads.delete",
			"hr.employees.create", "hr.employees.read", "hr.employees.update", "hr.employees.delete",
			"accounting.create", "accounting.read", "accounting.update", "accounting.delete",
			"rbac.roles.create", "rbac.roles.read", "rbac.roles.update", "rbac.roles.delete",
			"admin.access",
		},
		"manager": {
			"sales.leads.create", "sales.leads.read", "sales.leads.update",
			"hr.employees.read", "hr.employees.update",
			"accounting.read",
			"reports.access",
		},
		"sales": {
			"sales.leads.create", "sales.leads.read", "sales.leads.update",
			"sales.customers.create", "sales.customers.read",
		},
		"hr": {
			"hr.employees.create", "hr.employees.read", "hr.employees.update",
			"hr.payroll.read",
		},
		"finance": {
			"accounting.create", "accounting.read", "accounting.update",
			"reports.financial",
		},
	}

	for templateKey, template := range models.DefaultTemplates {
		// Check if template already exists
		query := `
			SELECT id FROM role_template 
			WHERE tenant_id = ? AND name = ? AND is_system_template = TRUE
		`
		var existingID string
		err := ts.db.QueryRowContext(ctx, query, tenantID, template.Name).Scan(&existingID)
		if err == nil {
			continue // Template already exists
		}
		if err != sql.ErrNoRows {
			ts.logger.Error("Failed to check template existence", "error", err)
			continue
		}

		// Create template
		insertQuery := `
			INSERT INTO role_template (
				id, tenant_id, name, description, category, is_system_template, 
				is_active, created_at, updated_at
			) VALUES (
				UUID(), ?, ?, ?, ?, ?, ?, NOW(), NOW()
			)
		`

		result, err := ts.db.ExecContext(ctx, insertQuery,
			tenantID, template.Name, template.Description,
			template.Category, template.IsSystemTemplate, template.IsActive,
		)
		if err != nil {
			ts.logger.Error("Failed to create template", "error", err, "template", templateKey)
			continue
		}

		// Get the inserted template ID
		lastID, _ := result.LastInsertId()
		templateID := fmt.Sprintf("template-%d", lastID)

		// Store template-permission mappings in metadata
		metadataJSON, _ := json.Marshal(map[string]interface{}{
			"permissions": templatePermissions[templateKey],
			"created_at":  time.Now(),
		})

		updateQuery := `
			UPDATE role_template SET metadata = ? WHERE id = ?
		`
		ts.db.ExecContext(ctx, updateQuery, string(metadataJSON), templateID)

		ts.logger.Info("Default template created", "template", templateKey, "tenant_id", tenantID)
	}

	return nil
}

// GetTemplate retrieves a role template
func (ts *RoleTemplateService) GetTemplate(ctx context.Context, tenantID string, templateID string) (*models.RoleTemplate, error) {
	query := `
		SELECT id, tenant_id, name, description, category, is_system_template, 
		       is_active, metadata, created_at, updated_at
		FROM role_template
		WHERE id = ? AND tenant_id = ?
	`

	template := &models.RoleTemplate{}
	var metadata sql.NullString

	err := ts.db.QueryRowContext(ctx, query, templateID, tenantID).Scan(
		&template.ID, &template.TenantID, &template.Name, &template.Description,
		&template.Category, &template.IsSystemTemplate, &template.IsActive,
		&metadata, &template.CreatedAt, &template.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found")
		}
		return nil, err
	}

	// Parse metadata for permissions
	if metadata.Valid {
		var meta map[string]interface{}
		if err := json.Unmarshal([]byte(metadata.String), &meta); err == nil {
			if perms, ok := meta["permissions"].([]interface{}); ok {
				for _, p := range perms {
					if pStr, ok := p.(string); ok {
						template.PermissionIDs = append(template.PermissionIDs, pStr)
					}
				}
			}
		}
	}

	return template, nil
}

// ListTemplates lists all role templates for a tenant
func (ts *RoleTemplateService) ListTemplates(ctx context.Context, tenantID string) ([]models.RoleTemplate, error) {
	query := `
		SELECT id, tenant_id, name, description, category, is_system_template, 
		       is_active, metadata, created_at, updated_at
		FROM role_template
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY is_system_template DESC, created_at DESC
	`

	rows, err := ts.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		ts.logger.Error("Failed to list templates", "error", err, "tenant_id", tenantID)
		return nil, err
	}
	defer rows.Close()

	var templates []models.RoleTemplate
	for rows.Next() {
		template := models.RoleTemplate{}
		var metadata sql.NullString

		err := rows.Scan(
			&template.ID, &template.TenantID, &template.Name, &template.Description,
			&template.Category, &template.IsSystemTemplate, &template.IsActive,
			&metadata, &template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			continue
		}

		// Parse metadata
		if metadata.Valid {
			var meta map[string]interface{}
			if err := json.Unmarshal([]byte(metadata.String), &meta); err == nil {
				if perms, ok := meta["permissions"].([]interface{}); ok {
					for _, p := range perms {
						if pStr, ok := p.(string); ok {
							template.PermissionIDs = append(template.PermissionIDs, pStr)
						}
					}
				}
			}
		}

		templates = append(templates, template)
	}

	return templates, nil
}

// CreateRoleFromTemplate creates a role based on a template
func (ts *RoleTemplateService) CreateRoleFromTemplate(
	ctx context.Context,
	tenantID string,
	templateID string,
	roleName string,
	customPermissions []string,
) (string, error) {
	// Get template
	template, err := ts.GetTemplate(ctx, tenantID, templateID)
	if err != nil {
		ts.logger.Error("Failed to get template", "error", err, "template_id", templateID)
		return "", err
	}

	// Use custom permissions if provided, otherwise use template permissions
	permissions := template.PermissionIDs
	if len(customPermissions) > 0 {
		permissions = customPermissions
	}

	// Create role
	roleQuery := `
		INSERT INTO role (id, tenant_id, role_name, description, is_active, created_at, updated_at)
		VALUES (UUID(), ?, ?, ?, TRUE, NOW(), NOW())
	`

	description := fmt.Sprintf("Created from template: %s", template.Name)
	result, err := ts.db.ExecContext(ctx, roleQuery, tenantID, roleName, description)
	if err != nil {
		ts.logger.Error("Failed to create role from template", "error", err)
		return "", err
	}

	roleID, _ := result.LastInsertId()
	roleIDStr := fmt.Sprintf("role-%d", roleID)

	// Assign permissions
	for _, permName := range permissions {
		// Get permission ID
		permQuery := `SELECT id FROM permission WHERE tenant_id = ? AND permission_name = ?`
		var permID string
		err := ts.db.QueryRowContext(ctx, permQuery, tenantID, permName).Scan(&permID)
		if err != nil {
			ts.logger.Warn("Permission not found, skipping", "permission", permName)
			continue
		}

		// Create role-permission mapping
		assignQuery := `
			INSERT INTO role_permission (id, tenant_id, role_id, permission_id, created_at)
			VALUES (UUID(), ?, ?, ?, NOW())
		`
		ts.db.ExecContext(ctx, assignQuery, tenantID, roleIDStr, permID)
	}

	ts.logger.Info("Role created from template",
		"template_id", templateID,
		"role_id", roleIDStr,
		"role_name", roleName,
		"tenant_id", tenantID,
	)

	return roleIDStr, nil
}

// CreateCustomTemplate creates a custom role template
func (ts *RoleTemplateService) CreateCustomTemplate(
	ctx context.Context,
	tenantID string,
	name string,
	description string,
	permissions []string,
) (*models.RoleTemplate, error) {
	template := &models.RoleTemplate{
		TenantID:         tenantID,
		Name:             name,
		Description:      description,
		Category:         "custom",
		IsSystemTemplate: false,
		IsActive:         true,
		PermissionIDs:    permissions,
	}

	// Create metadata
	metadata := map[string]interface{}{
		"permissions": permissions,
		"created_at":  time.Now(),
	}
	metadataJSON, _ := json.Marshal(metadata)

	query := `
		INSERT INTO role_template (
			id, tenant_id, name, description, category, is_system_template,
			is_active, metadata, created_at, updated_at
		) VALUES (
			UUID(), ?, ?, ?, ?, ?, ?, ?, NOW(), NOW()
		)
	`

	result, err := ts.db.ExecContext(ctx, query,
		tenantID, name, description, template.Category,
		template.IsSystemTemplate, template.IsActive, string(metadataJSON),
	)
	if err != nil {
		ts.logger.Error("Failed to create custom template", "error", err)
		return nil, err
	}

	lastID, _ := result.LastInsertId()
	template.ID = fmt.Sprintf("template-%d", lastID)

	ts.logger.Info("Custom template created", "template_id", template.ID, "tenant_id", tenantID)

	return template, nil
}

// SaveAsTemplate saves a role as a reusable template
func (ts *RoleTemplateService) SaveAsTemplate(
	ctx context.Context,
	tenantID string,
	roleID string,
	templateName string,
) (*models.RoleTemplate, error) {
	// Get role permissions
	query := `
		SELECT DISTINCT p.permission_name
		FROM role_permission rp
		JOIN permission p ON rp.permission_id = p.id
		WHERE rp.role_id = ? AND rp.tenant_id = ?
	`

	rows, err := ts.db.QueryContext(ctx, query, roleID, tenantID)
	if err != nil {
		ts.logger.Error("Failed to get role permissions", "error", err)
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			continue
		}
		permissions = append(permissions, perm)
	}

	// Create template from role
	return ts.CreateCustomTemplate(ctx, tenantID, templateName,
		fmt.Sprintf("Created from role %s", roleID), permissions)
}
