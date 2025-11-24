package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/pkg/logger"
)

// RBACService handles role-based access control
type RBACService struct {
	db     *sql.DB
	logger *logger.Logger
	// Cache for permissions by role
	permCache map[string][]string
}

// NewRBACService creates a new RBAC service
func NewRBACService(db *sql.DB, log *logger.Logger) *RBACService {
	return &RBACService{
		db:        db,
		logger:    log,
		permCache: make(map[string][]string),
	}
}

// CreateRole creates a new role with permissions
func (rs *RBACService) CreateRole(ctx context.Context, role *models.Role) error {
	query := `
		INSERT INTO roles (tenant_id, name, description, permissions, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	permJSON, _ := json.Marshal(role.Permissions)
	result, err := rs.db.ExecContext(ctx, query, role.TenantID, role.Name, role.Description, string(permJSON), role.IsActive, time.Now(), time.Now())
	if err != nil {
		rs.logger.Error("Failed to create role", "error", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	role.ID = id
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	// Clear cache
	rs.permCache = make(map[string][]string)

	return nil
}

// GetRole retrieves a role by ID
func (rs *RBACService) GetRole(ctx context.Context, tenantID string, roleID int64) (*models.Role, error) {
	query := `
		SELECT id, tenant_id, name, description, permissions, is_active, created_at, updated_at
		FROM roles
		WHERE id = ? AND tenant_id = ?
	`

	role := &models.Role{}
	var permJSON string

	err := rs.db.QueryRowContext(ctx, query, roleID, tenantID).Scan(
		&role.ID, &role.TenantID, &role.Name, &role.Description, &permJSON,
		&role.IsActive, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	_ = json.Unmarshal([]byte(permJSON), &role.Permissions)
	return role, nil
}

// ListRoles retrieves all roles for a tenant
func (rs *RBACService) ListRoles(ctx context.Context, tenantID string) ([]models.Role, error) {
	query := `
		SELECT id, tenant_id, name, description, permissions, is_active, created_at, updated_at
		FROM roles
		WHERE tenant_id = ? AND is_active = true
		ORDER BY name
	`

	rows, err := rs.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		var permJSON string

		err := rows.Scan(&role.ID, &role.TenantID, &role.Name, &role.Description, &permJSON,
			&role.IsActive, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal([]byte(permJSON), &role.Permissions)
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

// UpdateRole updates a role and its permissions
func (rs *RBACService) UpdateRole(ctx context.Context, role *models.Role) error {
	query := `
		UPDATE roles
		SET name = ?, description = ?, permissions = ?, is_active = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	permJSON, _ := json.Marshal(role.Permissions)
	_, err := rs.db.ExecContext(ctx, query, role.Name, role.Description, string(permJSON), role.IsActive, time.Now(), role.ID, role.TenantID)
	if err != nil {
		rs.logger.Error("Failed to update role", "error", err)
		return err
	}

	// Clear cache
	rs.permCache = make(map[string][]string)

	return nil
}

// DeleteRole soft-deletes a role
func (rs *RBACService) DeleteRole(ctx context.Context, tenantID string, roleID int64) error {
	query := `
		UPDATE roles
		SET is_active = false, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := rs.db.ExecContext(ctx, query, time.Now(), roleID, tenantID)
	if err != nil {
		rs.logger.Error("Failed to delete role", "error", err)
		return err
	}

	// Clear cache
	rs.permCache = make(map[string][]string)

	return nil
}

// CreatePermission creates a new permission
func (rs *RBACService) CreatePermission(ctx context.Context, permission *models.Permission) error {
	query := `
		INSERT INTO permissions (code, description, category, created_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := rs.db.ExecContext(ctx, query, permission.Code, permission.Description, permission.Category, time.Now())
	if err != nil {
		rs.logger.Error("Failed to create permission", "error", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	permission.ID = id
	permission.CreatedAt = time.Now()

	return nil
}

// GetUserPermissions retrieves all permissions for a user
func (rs *RBACService) GetUserPermissions(ctx context.Context, tenantID string, userID int64) ([]string, error) {
	query := `
		SELECT DISTINCT p.code
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = ? AND ur.tenant_id = ?
	`

	rows, err := rs.db.QueryContext(ctx, query, userID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, rows.Err()
}

// HasPermission checks if a user has a specific permission
func (rs *RBACService) HasPermission(ctx context.Context, tenantID string, userID int64, permissionCode string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = ? AND ur.tenant_id = ? AND p.code = ?
	`

	var count int
	err := rs.db.QueryRowContext(ctx, query, userID, tenantID, permissionCode).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// AssignRoleToUser assigns a role to a user
func (rs *RBACService) AssignRoleToUser(ctx context.Context, tenantID string, userID, roleID int64) error {
	query := `
		INSERT INTO user_roles (tenant_id, user_id, role_id, created_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := rs.db.ExecContext(ctx, query, tenantID, userID, roleID, time.Now())
	if err != nil {
		if sql.ErrNoRows != err {
			rs.logger.Error("Failed to assign role to user", "error", err)
		}
		return err
	}

	// Clear cache
	rs.permCache = make(map[string][]string)

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (rs *RBACService) RemoveRoleFromUser(ctx context.Context, tenantID string, userID, roleID int64) error {
	query := `
		DELETE FROM user_roles
		WHERE tenant_id = ? AND user_id = ? AND role_id = ?
	`

	_, err := rs.db.ExecContext(ctx, query, tenantID, userID, roleID)
	if err != nil {
		rs.logger.Error("Failed to remove role from user", "error", err)
		return err
	}

	// Clear cache
	rs.permCache = make(map[string][]string)

	return nil
}

// GetUserRoles retrieves all roles assigned to a user
func (rs *RBACService) GetUserRoles(ctx context.Context, tenantID string, userID int64) ([]models.Role, error) {
	query := `
		SELECT DISTINCT r.id, r.tenant_id, r.name, r.description, r.permissions, r.is_active, r.created_at, r.updated_at
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND ur.tenant_id = ?
		ORDER BY r.name
	`

	rows, err := rs.db.QueryContext(ctx, query, userID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		var permJSON string

		err := rows.Scan(&role.ID, &role.TenantID, &role.Name, &role.Description, &permJSON,
			&role.IsActive, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal([]byte(permJSON), &role.Permissions)
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

// SetupDefaultRoles creates standard roles for a tenant
func (rs *RBACService) SetupDefaultRoles(ctx context.Context, tenantID string) error {
	defaultRoles := []struct {
		name        string
		description string
		permissions []string
	}{
		{
			name:        "admin",
			description: "Full system access",
			permissions: []string{
				"leads.create", "leads.read", "leads.update", "leads.delete",
				"calls.create", "calls.read", "calls.update", "calls.delete", "calls.record",
				"campaigns.create", "campaigns.read", "campaigns.update", "campaigns.delete",
				"agents.create", "agents.read", "agents.update", "agents.delete",
				"users.create", "users.read", "users.update", "users.delete",
				"roles.create", "roles.read", "roles.update", "roles.delete",
				"reports.export", "reports.view", "audit.view",
				"settings.manage",
			},
		},
		{
			name:        "manager",
			description: "Team and campaign management",
			permissions: []string{
				"leads.read", "leads.update",
				"calls.read", "calls.record",
				"campaigns.read", "campaigns.update",
				"agents.read", "agents.update",
				"reports.export", "reports.view",
				"audit.view",
			},
		},
		{
			name:        "agent",
			description: "Day-to-day operations",
			permissions: []string{
				"leads.read", "leads.update",
				"calls.create", "calls.read", "calls.record",
				"campaigns.read",
				"reports.view",
			},
		},
		{
			name:        "supervisor",
			description: "Monitoring and quality assurance",
			permissions: []string{
				"leads.read",
				"calls.read", "calls.record",
				"campaigns.read",
				"agents.read",
				"reports.view",
				"audit.view",
			},
		},
	}

	for _, r := range defaultRoles {
		role := &models.Role{
			TenantID:    tenantID,
			Name:        r.name,
			Description: r.description,
			Permissions: r.permissions,
			IsActive:    true,
		}

		err := rs.CreateRole(ctx, role)
		if err != nil {
			rs.logger.Warn("Failed to create default role", "role", r.name, "error", err)
			// Continue with other roles
		}
	}

	return nil
}

// VerifyPermission checks permission and returns error if denied
func (rs *RBACService) VerifyPermission(ctx context.Context, tenantID string, userID int64, permissionCode string) error {
	hasPermission, err := rs.HasPermission(ctx, tenantID, userID, permissionCode)
	if err != nil {
		return fmt.Errorf("permission check failed: %w", err)
	}

	if !hasPermission {
		return fmt.Errorf("permission denied: %s", permissionCode)
	}

	return nil
}
