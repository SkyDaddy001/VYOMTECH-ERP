package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"

	"github.com/google/uuid"
)

// RBACHandler manages role-based access control operations
type RBACHandler struct {
	rbacService *services.RBACService
	db          *sql.DB
	logger      *logger.Logger
}

// NewRBACHandler creates a new RBAC handler
func NewRBACHandler(rbacService *services.RBACService, db *sql.DB, logger *logger.Logger) *RBACHandler {
	return &RBACHandler{
		rbacService: rbacService,
		db:          db,
		logger:      logger,
	}
}

// generateID generates a new UUID
func (h *RBACHandler) generateID() string {
	return uuid.New().String()
}

// RBACCreateRoleRequest represents a request to create a new role
type RBACCreateRoleRequest struct {
	RoleName    string `json:"role_name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// RoleResponse represents a role with its permissions
type RoleResponse struct {
	ID          string               `json:"id"`
	TenantID    string               `json:"tenant_id"`
	RoleName    string               `json:"role_name"`
	Description string               `json:"description"`
	IsActive    bool                 `json:"is_active"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

// PermissionResponse represents a permission
type PermissionResponse struct {
	ID             string `json:"id"`
	PermissionName string `json:"permission_name"`
	Description    string `json:"description"`
	Resource       string `json:"resource"`
	Action         string `json:"action"`
}

// RBACAssignPermissionsRequest represents a request to assign permissions to a role
type RBACAssignPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required"`
}

// RBACErrorResponse represents an error response
type RBACErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// RBACSuccessResponse represents a success response
type RBACSuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// respondError sends an error response
func (h *RBACHandler) respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(RBACErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	})
}

// respondSuccess sends a success response
func (h *RBACHandler) respondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(RBACSuccessResponse{
		Message: message,
		Data:    data,
	})
}

// CreateRole creates a new role
// POST /api/v1/rbac/roles
func (h *RBACHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	// Extract context
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Verify admin permission (using a generic admin check)
	// In production, this should check for a specific "rbac.roles.create" permission
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || !strings.Contains(strings.ToLower(role), "admin") {
		h.respondError(w, http.StatusForbidden, "Admin access required")
		return
	}

	// Parse request
	var req RBACCreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if strings.TrimSpace(req.RoleName) == "" {
		h.respondError(w, http.StatusBadRequest, "Role name is required")
		return
	}

	// Create role in database
	query := `
		INSERT INTO role (id, tenant_id, role_name, description, is_active, created_at, updated_at)
		VALUES (UUID(), ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := h.db.ExecContext(r.Context(), query, tenantID, req.RoleName, req.Description, req.IsActive)
	if err != nil {
		h.logger.Error("Failed to create role", "error", err, "tenant_id", tenantID)
		if strings.Contains(err.Error(), "Duplicate entry") {
			h.respondError(w, http.StatusConflict, "Role with this name already exists in this tenant")
		} else {
			h.respondError(w, http.StatusInternalServerError, "Failed to create role")
		}
		return
	}

	// Get the created role ID
	var roleID string
	err = h.db.QueryRowContext(r.Context(), `
		SELECT id FROM role WHERE tenant_id = ? AND role_name = ? ORDER BY created_at DESC LIMIT 1
	`, tenantID, req.RoleName).Scan(&roleID)
	if err != nil {
		h.logger.Error("Failed to retrieve created role", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Role created but failed to retrieve")
		return
	}

	h.logger.Info("Role created successfully", "role_id", roleID, "tenant_id", tenantID, "created_by", userID)

	roleResp := RoleResponse{
		ID:          roleID,
		TenantID:    tenantID,
		RoleName:    req.RoleName,
		Description: req.Description,
		IsActive:    req.IsActive,
		Permissions: []PermissionResponse{},
	}

	h.respondSuccess(w, http.StatusCreated, "Role created successfully", roleResp)
}

// GetRole retrieves a role with its permissions
// GET /api/v1/rbac/roles/:id
func (h *RBACHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	// Extract context
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	// Get role ID from URL path
	roleID := strings.TrimPrefix(r.URL.Path, "/api/v1/rbac/roles/")

	if roleID == "" {
		h.respondError(w, http.StatusBadRequest, "Role ID is required")
		return
	}

	// Query role
	var role RoleResponse
	query := `
		SELECT id, tenant_id, role_name, description, is_active, created_at, updated_at
		FROM role
		WHERE id = ? AND tenant_id = ?
	`

	err := h.db.QueryRowContext(r.Context(), query, roleID, tenantID).Scan(
		&role.ID, &role.TenantID, &role.RoleName, &role.Description, &role.IsActive,
		&role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			h.respondError(w, http.StatusNotFound, "Role not found")
		} else {
			h.logger.Error("Failed to retrieve role", "error", err, "role_id", roleID)
			h.respondError(w, http.StatusInternalServerError, "Failed to retrieve role")
		}
		return
	}

	// Get permissions for role
	permQuery := `
		SELECT p.id, p.permission_name, p.description, p.resource, p.action
		FROM permission p
		INNER JOIN role_permission rp ON p.id = rp.permission_id
		WHERE rp.role_id = ? AND rp.tenant_id = ?
	`

	rows, err := h.db.QueryContext(r.Context(), permQuery, roleID, tenantID)
	if err != nil {
		h.logger.Error("Failed to retrieve role permissions", "error", err, "role_id", roleID)
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve permissions")
		return
	}
	defer rows.Close()

	role.Permissions = []PermissionResponse{}
	for rows.Next() {
		var perm PermissionResponse
		if err := rows.Scan(&perm.ID, &perm.PermissionName, &perm.Description, &perm.Resource, &perm.Action); err != nil {
			h.logger.Error("Failed to scan permission", "error", err)
			continue
		}
		role.Permissions = append(role.Permissions, perm)
	}

	h.respondSuccess(w, http.StatusOK, "Role retrieved successfully", role)
}

// AssignPermissions assigns permissions to a role
// PUT /api/v1/rbac/roles/:id/permissions
func (h *RBACHandler) AssignPermissions(w http.ResponseWriter, r *http.Request) {
	// Extract context
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Verify admin permission
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || !strings.Contains(strings.ToLower(role), "admin") {
		h.respondError(w, http.StatusForbidden, "Admin access required")
		return
	}

	// Get role ID from URL path
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/rbac/roles/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		h.respondError(w, http.StatusBadRequest, "Role ID is required")
		return
	}
	roleID := parts[0]

	// Verify role exists and belongs to tenant
	var exists bool
	err := h.db.QueryRowContext(r.Context(), `
		SELECT EXISTS(SELECT 1 FROM role WHERE id = ? AND tenant_id = ?)
	`, roleID, tenantID).Scan(&exists)
	if err != nil || !exists {
		h.respondError(w, http.StatusNotFound, "Role not found")
		return
	}

	// Parse request
	var req RBACAssignPermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.PermissionIDs) == 0 {
		h.respondError(w, http.StatusBadRequest, "Permission IDs are required")
		return
	}

	// Start transaction
	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		h.logger.Error("Failed to start transaction", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update permissions")
		return
	}
	defer tx.Rollback()

	// Delete existing role-permission mappings
	_, err = tx.ExecContext(r.Context(), `
		DELETE FROM role_permission WHERE role_id = ? AND tenant_id = ?
	`, roleID, tenantID)
	if err != nil {
		h.logger.Error("Failed to delete existing permissions", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update permissions")
		return
	}

	// Insert new role-permission mappings
	for _, permID := range req.PermissionIDs {
		// Verify permission exists and belongs to tenant
		var permExists bool
		err := tx.QueryRowContext(r.Context(), `
			SELECT EXISTS(SELECT 1 FROM permission WHERE id = ? AND tenant_id = ?)
		`, permID, tenantID).Scan(&permExists)
		if err != nil || !permExists {
			h.respondError(w, http.StatusBadRequest, fmt.Sprintf("Permission %s not found", permID))
			return
		}

		_, err = tx.ExecContext(r.Context(), `
			INSERT INTO role_permission (id, tenant_id, role_id, permission_id, created_at)
			VALUES (UUID(), ?, ?, ?, NOW())
		`, tenantID, roleID, permID)
		if err != nil {
			h.logger.Error("Failed to assign permission", "error", err, "perm_id", permID)
			h.respondError(w, http.StatusInternalServerError, "Failed to assign permission")
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		h.logger.Error("Failed to commit transaction", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update permissions")
		return
	}

	h.logger.Info("Permissions assigned to role", "role_id", roleID, "perm_count", len(req.PermissionIDs), "assigned_by", userID)

	h.respondSuccess(w, http.StatusOK, "Permissions assigned successfully", map[string]interface{}{
		"role_id":        roleID,
		"permission_ids": req.PermissionIDs,
	})
}

// DeleteRole deletes a role (soft delete)
// DELETE /api/v1/rbac/roles/:id
func (h *RBACHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Extract context
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "User ID not found")
		return
	}

	// Verify admin permission
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || !strings.Contains(strings.ToLower(role), "admin") {
		h.respondError(w, http.StatusForbidden, "Admin access required")
		return
	}

	// Get role ID from URL path
	roleID := strings.TrimPrefix(r.URL.Path, "/api/v1/rbac/roles/")

	if roleID == "" {
		h.respondError(w, http.StatusBadRequest, "Role ID is required")
		return
	}

	// Check if role exists
	var exists bool
	err := h.db.QueryRowContext(r.Context(), `
		SELECT EXISTS(SELECT 1 FROM role WHERE id = ? AND tenant_id = ?)
	`, roleID, tenantID).Scan(&exists)
	if err != nil || !exists {
		h.respondError(w, http.StatusNotFound, "Role not found")
		return
	}

	// Soft delete role (set is_active to false)
	_, err = h.db.ExecContext(r.Context(), `
		UPDATE role SET is_active = FALSE, updated_at = NOW() WHERE id = ? AND tenant_id = ?
	`, roleID, tenantID)
	if err != nil {
		h.logger.Error("Failed to delete role", "error", err, "role_id", roleID)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete role")
		return
	}

	h.logger.Info("Role deleted (soft delete)", "role_id", roleID, "tenant_id", tenantID, "deleted_by", userID)

	h.respondSuccess(w, http.StatusOK, "Role deleted successfully", map[string]string{
		"role_id": roleID,
	})
}

// ListRoles lists all roles for a tenant
// GET /api/v1/rbac/roles
func (h *RBACHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	// Extract context
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	// Query roles
	query := `
		SELECT id, tenant_id, role_name, description, is_active, created_at, updated_at
		FROM role
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY created_at DESC
	`

	rows, err := h.db.QueryContext(r.Context(), query, tenantID)
	if err != nil {
		h.logger.Error("Failed to list roles", "error", err, "tenant_id", tenantID)
		h.respondError(w, http.StatusInternalServerError, "Failed to list roles")
		return
	}
	defer rows.Close()

	roles := []RoleResponse{}
	for rows.Next() {
		var role RoleResponse
		if err := rows.Scan(&role.ID, &role.TenantID, &role.RoleName, &role.Description, &role.IsActive, &role.CreatedAt, &role.UpdatedAt); err != nil {
			h.logger.Error("Failed to scan role", "error", err)
			continue
		}

		// Get permissions for this role
		permQuery := `
			SELECT p.id, p.permission_name, p.description, p.resource, p.action
			FROM permission p
			INNER JOIN role_permission rp ON p.id = rp.permission_id
			WHERE rp.role_id = ? AND rp.tenant_id = ?
		`

		permRows, err := h.db.QueryContext(r.Context(), permQuery, role.ID, tenantID)
		if err != nil {
			h.logger.Error("Failed to retrieve role permissions", "error", err)
			continue
		}

		role.Permissions = []PermissionResponse{}
		for permRows.Next() {
			var perm PermissionResponse
			if err := permRows.Scan(&perm.ID, &perm.PermissionName, &perm.Description, &perm.Resource, &perm.Action); err != nil {
				continue
			}
			role.Permissions = append(role.Permissions, perm)
		}
		permRows.Close()

		roles = append(roles, role)
	}

	h.respondSuccess(w, http.StatusOK, "Roles listed successfully", roles)
}

// ListPermissions lists all available permissions for a tenant
// GET /api/v1/rbac/permissions
func (h *RBACHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	// Extract context
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	// Query permissions
	query := `
		SELECT id, permission_name, description, resource, action
		FROM permission
		WHERE tenant_id = ?
		ORDER BY resource, action
	`

	rows, err := h.db.QueryContext(r.Context(), query, tenantID)
	if err != nil {
		h.logger.Error("Failed to list permissions", "error", err, "tenant_id", tenantID)
		h.respondError(w, http.StatusInternalServerError, "Failed to list permissions")
		return
	}
	defer rows.Close()

	permissions := []PermissionResponse{}
	for rows.Next() {
		var perm PermissionResponse
		if err := rows.Scan(&perm.ID, &perm.PermissionName, &perm.Description, &perm.Resource, &perm.Action); err != nil {
			h.logger.Error("Failed to scan permission", "error", err)
			continue
		}
		permissions = append(permissions, perm)
	}

	h.respondSuccess(w, http.StatusOK, "Permissions listed successfully", permissions)
}

// GrantResourceAccess grants user access to a specific resource (Phase 4.1)
// POST /api/v1/rbac/resource-access
func (h *RBACHandler) GrantResourceAccess(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusForbidden, "User ID not found")
		return
	}

	// Check admin permission
	err := h.rbacService.VerifyPermission(r.Context(), tenantID, userID, "rbac.admin")
	if err != nil {
		h.respondError(w, http.StatusForbidden, "Insufficient permissions for resource access management")
		return
	}

	var req struct {
		UserID       int64   `json:"user_id"`
		ResourceType string  `json:"resource_type"`
		ResourceID   string  `json:"resource_id"`
		AccessLevel  string  `json:"access_level"`
		ExpiresAt    *string `json:"expires_at"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UserID == 0 || req.ResourceType == "" || req.ResourceID == "" || req.AccessLevel == "" {
		h.respondError(w, http.StatusBadRequest, "Missing required fields: user_id, resource_type, resource_id, access_level")
		return
	}

	// Validate access level
	validLevels := map[string]bool{"view": true, "edit": true, "delete": true, "admin": true}
	if !validLevels[req.AccessLevel] {
		h.respondError(w, http.StatusBadRequest, "Invalid access level. Must be: view, edit, delete, or admin")
		return
	}

	// Insert resource access
	id := h.generateID()
	query := `
		INSERT INTO resource_access (id, tenant_id, user_id, resource_type, resource_id, access_level, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE access_level = VALUES(access_level), expires_at = VALUES(expires_at), updated_at = NOW()
	`

	_, err = h.db.ExecContext(r.Context(), query, id, tenantID, req.UserID, req.ResourceType, req.ResourceID, req.AccessLevel, req.ExpiresAt)
	if err != nil {
		h.logger.Error("Failed to grant resource access", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to grant resource access")
		return
	}

	h.respondSuccess(w, http.StatusCreated, "Resource access granted successfully", map[string]string{"id": id})
}

// CreateTimeBasedPermission creates a permission with time window (Phase 4.2)
// POST /api/v1/rbac/time-based-permissions
func (h *RBACHandler) CreateTimeBasedPermission(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusForbidden, "User ID not found")
		return
	}

	// Check admin permission
	err := h.rbacService.VerifyPermission(r.Context(), tenantID, userID, "rbac.admin")
	if err != nil {
		h.respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var req struct {
		RoleID        string  `json:"role_id"`
		PermissionID  string  `json:"permission_id"`
		EffectiveFrom string  `json:"effective_from"`
		ExpiresAt     string  `json:"expires_at"`
		Reason        *string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RoleID == "" || req.PermissionID == "" || req.EffectiveFrom == "" || req.ExpiresAt == "" {
		h.respondError(w, http.StatusBadRequest, "Missing required fields: role_id, permission_id, effective_from, expires_at")
		return
	}

	id := h.generateID()
	query := `
		INSERT INTO time_based_permission (id, tenant_id, role_id, permission_id, effective_from, expires_at, is_active, reason, approved_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, true, ?, ?, NOW(), NOW())
	`

	_, err = h.db.ExecContext(r.Context(), query, id, tenantID, req.RoleID, req.PermissionID, req.EffectiveFrom, req.ExpiresAt, req.Reason, userID)
	if err != nil {
		h.logger.Error("Failed to create time-based permission", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create time-based permission")
		return
	}

	h.respondSuccess(w, http.StatusCreated, "Time-based permission created successfully", map[string]string{"id": id})
}

// SetFieldLevelPermission sets field visibility/editability for a role (Phase 4.3)
// POST /api/v1/rbac/field-permissions
func (h *RBACHandler) SetFieldLevelPermission(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusForbidden, "User ID not found")
		return
	}

	// Check admin permission
	err := h.rbacService.VerifyPermission(r.Context(), tenantID, userID, "rbac.admin")
	if err != nil {
		h.respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var req struct {
		RoleID      string  `json:"role_id"`
		ModuleName  string  `json:"module_name"`
		EntityName  string  `json:"entity_name"`
		FieldName   string  `json:"field_name"`
		CanView     bool    `json:"can_view"`
		CanEdit     bool    `json:"can_edit"`
		IsMasked    bool    `json:"is_masked"`
		MaskPattern *string `json:"mask_pattern"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RoleID == "" || req.ModuleName == "" || req.EntityName == "" || req.FieldName == "" {
		h.respondError(w, http.StatusBadRequest, "Missing required fields: role_id, module_name, entity_name, field_name")
		return
	}

	id := h.generateID()
	query := `
		INSERT INTO field_level_permission (id, tenant_id, role_id, module_name, entity_name, field_name, can_view, can_edit, is_masked, mask_pattern, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE can_view = VALUES(can_view), can_edit = VALUES(can_edit), is_masked = VALUES(is_masked), mask_pattern = VALUES(mask_pattern), updated_at = NOW()
	`

	_, err = h.db.ExecContext(r.Context(), query, id, tenantID, req.RoleID, req.ModuleName, req.EntityName, req.FieldName, req.CanView, req.CanEdit, req.IsMasked, req.MaskPattern)
	if err != nil {
		h.logger.Error("Failed to set field-level permission", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to set field-level permission")
		return
	}

	h.respondSuccess(w, http.StatusCreated, "Field-level permission set successfully", map[string]string{"id": id})
}

// DelegateRole creates role delegation for sub-role creation (Phase 4.5)
// POST /api/v1/rbac/delegations
func (h *RBACHandler) DelegateRole(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusForbidden, "User ID not found")
		return
	}

	// Check admin permission
	err := h.rbacService.VerifyPermission(r.Context(), tenantID, userID, "rbac.admin")
	if err != nil {
		h.respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var req struct {
		ParentRoleID    string  `json:"parent_role_id"`
		SubRoleID       string  `json:"sub_role_id"`
		PermissionBound string  `json:"permission_bound"`
		ExpiresAt       *string `json:"expires_at"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ParentRoleID == "" || req.SubRoleID == "" {
		h.respondError(w, http.StatusBadRequest, "Missing required fields: parent_role_id, sub_role_id")
		return
	}

	id := h.generateID()
	query := `
		INSERT INTO role_delegation (id, tenant_id, parent_role_id, sub_role_id, permission_bound, delegated_by, is_active, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, true, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE permission_bound = VALUES(permission_bound), expires_at = VALUES(expires_at), updated_at = NOW()
	`

	_, err = h.db.ExecContext(r.Context(), query, id, tenantID, req.ParentRoleID, req.SubRoleID, req.PermissionBound, userID, req.ExpiresAt)
	if err != nil {
		h.logger.Error("Failed to create role delegation", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create role delegation")
		return
	}

	h.respondSuccess(w, http.StatusCreated, "Role delegation created successfully", map[string]string{"id": id})
}

// BulkAssignPermissions handles bulk permission assignment (Phase 4.4)
// POST /api/v1/rbac/bulk-assign
func (h *RBACHandler) BulkAssignPermissions(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusForbidden, "User ID not found")
		return
	}

	// Check admin permission
	err := h.rbacService.VerifyPermission(r.Context(), tenantID, userID, "rbac.admin")
	if err != nil {
		h.respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var req struct {
		TargetType    string   `json:"target_type"` // "user" or "role"
		TargetIDs     []int64  `json:"target_ids"`
		PermissionIDs []string `json:"permission_ids"`
		Action        string   `json:"action"` // "assign" or "revoke"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.TargetType == "" || len(req.TargetIDs) == 0 || len(req.PermissionIDs) == 0 || (req.Action != "assign" && req.Action != "revoke") {
		h.respondError(w, http.StatusBadRequest, "Missing or invalid fields: target_type, target_ids, permission_ids, action")
		return
	}

	logID := h.generateID()
	successCount := 0
	failedCount := 0

	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback()

	// Process bulk assignment
	for _, targetID := range req.TargetIDs {
		for _, permID := range req.PermissionIDs {
			if req.Action == "assign" {
				_, err = tx.ExecContext(r.Context(),
					`INSERT INTO role_permission (role_id, permission_id, tenant_id) VALUES (?, ?, ?)`,
					targetID, permID, tenantID)
			} else {
				_, err = tx.ExecContext(r.Context(),
					`DELETE FROM role_permission WHERE role_id = ? AND permission_id = ? AND tenant_id = ?`,
					targetID, permID, tenantID)
			}

			if err != nil {
				h.logger.Error("Failed to process permission", "error", err)
				failedCount++
			} else {
				successCount++
			}
		}
	}

	// Log bulk operation
	detailsJSON := fmt.Sprintf(`{"action": "%s", "target_ids": %d, "permission_ids": %d}`, req.Action, len(req.TargetIDs), len(req.PermissionIDs))
	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO bulk_permission_log (id, tenant_id, assignment_type, target_type, total_count, success_count, failed_count, executed_by, details)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		logID, tenantID, req.Action, req.TargetType, len(req.TargetIDs)*len(req.PermissionIDs), successCount, failedCount, userID, detailsJSON)

	if err != nil {
		h.logger.Error("Failed to log bulk operation", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to log bulk operation")
		return
	}

	if err = tx.Commit(); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	h.respondSuccess(w, http.StatusOK, "Bulk operation completed", map[string]interface{}{
		"log_id":  logID,
		"success": successCount,
		"failed":  failedCount,
		"total":   len(req.TargetIDs) * len(req.PermissionIDs),
	})
}

// ============================================================
// PHASE 3.6: USER ROLE ASSIGNMENT & MEMBERSHIP
// ============================================================

// AssignRoleToUser assigns a role to a user (Phase 3.6)
// POST /api/v1/rbac/users/{user_id}/roles
func (h *RBACHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusForbidden, "User ID not found")
		return
	}

	// Check admin permission
	err := h.rbacService.VerifyPermission(r.Context(), tenantID, userID, "rbac.admin")
	if err != nil {
		h.respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var req struct {
		UserID    int64   `json:"user_id"`
		RoleID    string  `json:"role_id"`
		ExpiresAt *string `json:"expires_at"` // optional expiration date
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UserID == 0 || req.RoleID == "" {
		h.respondError(w, http.StatusBadRequest, "Missing required fields: user_id, role_id")
		return
	}

	// Verify role exists and belongs to tenant
	var roleExists bool
	err = h.db.QueryRowContext(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM role WHERE id = ? AND tenant_id = ?)`,
		req.RoleID, tenantID).Scan(&roleExists)

	if err != nil || !roleExists {
		h.respondError(w, http.StatusBadRequest, "Role not found or doesn't belong to this tenant")
		return
	}

	// Check if assignment already exists
	var assignmentExists bool
	err = h.db.QueryRowContext(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM user_role WHERE user_id = ? AND role_id = ? AND tenant_id = ?)`,
		req.UserID, req.RoleID, tenantID).Scan(&assignmentExists)

	if assignmentExists {
		h.respondError(w, http.StatusConflict, "User already has this role assigned")
		return
	}

	// Assign role to user
	query := `
		INSERT INTO user_role (tenant_id, user_id, role_id, expires_at, created_at)
		VALUES (?, ?, ?, ?, NOW())
	`

	_, err = h.db.ExecContext(r.Context(), query, tenantID, req.UserID, req.RoleID, req.ExpiresAt)
	if err != nil {
		h.logger.Error("Failed to assign role to user", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to assign role to user")
		return
	}

	h.respondSuccess(w, http.StatusCreated, "Role assigned to user successfully", map[string]interface{}{
		"user_id": req.UserID,
		"role_id": req.RoleID,
	})
}

// RemoveRoleFromUser removes a role from a user (Phase 3.6)
// DELETE /api/v1/rbac/users/{user_id}/roles/{role_id}
func (h *RBACHandler) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusForbidden, "User ID not found")
		return
	}

	// Check admin permission
	err := h.rbacService.VerifyPermission(r.Context(), tenantID, userID, "rbac.admin")
	if err != nil {
		h.respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var req struct {
		UserID int64  `json:"user_id"`
		RoleID string `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UserID == 0 || req.RoleID == "" {
		h.respondError(w, http.StatusBadRequest, "Missing required fields: user_id, role_id")
		return
	}

	// Remove role from user
	result, err := h.db.ExecContext(r.Context(),
		`DELETE FROM user_role WHERE user_id = ? AND role_id = ? AND tenant_id = ?`,
		req.UserID, req.RoleID, tenantID)

	if err != nil {
		h.logger.Error("Failed to remove role from user", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to remove role from user")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		h.respondError(w, http.StatusNotFound, "User-role assignment not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, "Role removed from user successfully", map[string]interface{}{
		"user_id": req.UserID,
		"role_id": req.RoleID,
	})
}

// GetUserRoles retrieves all roles assigned to a user (Phase 3.6)
// GET /api/v1/rbac/users/{user_id}/roles
func (h *RBACHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	_, ok = r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusForbidden, "User ID not found")
		return
	}

	// Parse target user ID from request
	var targetUserID int64
	targetIDStr := strings.TrimPrefix(r.URL.Path, "/api/v1/rbac/users/")
	targetIDStr = strings.Split(targetIDStr, "/")[0]
	fmt.Sscanf(targetIDStr, "%d", &targetUserID)

	if targetUserID == 0 {
		h.respondError(w, http.StatusBadRequest, "Invalid user ID in path")
		return
	}

	query := `
		SELECT r.id, r.name, r.description, r.is_active, ur.expires_at, ur.created_at
		FROM role r
		JOIN user_role ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND ur.tenant_id = ?
		ORDER BY r.name
	`

	rows, err := h.db.QueryContext(r.Context(), query, targetUserID, tenantID)
	if err != nil {
		h.logger.Error("Failed to get user roles", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve user roles")
		return
	}
	defer rows.Close()

	type UserRoleResponse struct {
		RoleID      string  `json:"role_id"`
		RoleName    string  `json:"role_name"`
		Description string  `json:"description"`
		IsActive    bool    `json:"is_active"`
		ExpiresAt   *string `json:"expires_at"`
		AssignedAt  string  `json:"assigned_at"`
	}

	roles := []UserRoleResponse{}
	for rows.Next() {
		var role UserRoleResponse
		if err := rows.Scan(&role.RoleID, &role.RoleName, &role.Description, &role.IsActive, &role.ExpiresAt, &role.AssignedAt); err != nil {
			h.logger.Error("Failed to scan user role", "error", err)
			continue
		}
		roles = append(roles, role)
	}

	h.respondSuccess(w, http.StatusOK, "User roles retrieved successfully", map[string]interface{}{
		"user_id": targetUserID,
		"roles":   roles,
	})
}

// GetRoleMembers retrieves all users assigned to a role (Phase 3.6)
// GET /api/v1/rbac/roles/{role_id}/members
func (h *RBACHandler) GetRoleMembers(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	_, ok = r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusForbidden, "User ID not found")
		return
	}

	// Parse role ID from path
	var roleID string
	roleIDStr := strings.TrimPrefix(r.URL.Path, "/api/v1/rbac/roles/")
	roleID = strings.Split(roleIDStr, "/")[0]

	if roleID == "" {
		h.respondError(w, http.StatusBadRequest, "Invalid role ID in path")
		return
	}

	// Verify role exists
	var roleExists bool
	err := h.db.QueryRowContext(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM role WHERE id = ? AND tenant_id = ?)`,
		roleID, tenantID).Scan(&roleExists)

	if err != nil || !roleExists {
		h.respondError(w, http.StatusNotFound, "Role not found")
		return
	}

	query := `
		SELECT ur.user_id, ur.expires_at, ur.created_at
		FROM user_role ur
		WHERE ur.role_id = ? AND ur.tenant_id = ?
		ORDER BY ur.created_at DESC
	`

	rows, err := h.db.QueryContext(r.Context(), query, roleID, tenantID)
	if err != nil {
		h.logger.Error("Failed to get role members", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve role members")
		return
	}
	defer rows.Close()

	type RoleMemberResponse struct {
		UserID     int64   `json:"user_id"`
		ExpiresAt  *string `json:"expires_at"`
		AssignedAt string  `json:"assigned_at"`
	}

	members := []RoleMemberResponse{}
	for rows.Next() {
		var member RoleMemberResponse
		if err := rows.Scan(&member.UserID, &member.ExpiresAt, &member.AssignedAt); err != nil {
			h.logger.Error("Failed to scan role member", "error", err)
			continue
		}
		members = append(members, member)
	}

	h.respondSuccess(w, http.StatusOK, "Role members retrieved successfully", map[string]interface{}{
		"role_id": roleID,
		"count":   len(members),
		"members": members,
	})
}

// UpdateUserRole updates a user's role assignment (Phase 3.6)
// PUT /api/v1/rbac/users/{user_id}/roles/{role_id}
func (h *RBACHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusForbidden, "User ID not found")
		return
	}

	// Check admin permission
	err := h.rbacService.VerifyPermission(r.Context(), tenantID, userID, "rbac.admin")
	if err != nil {
		h.respondError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var req struct {
		UserID    int64   `json:"user_id"`
		RoleID    string  `json:"role_id"`
		ExpiresAt *string `json:"expires_at"` // update expiration date
		IsActive  *bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UserID == 0 || req.RoleID == "" {
		h.respondError(w, http.StatusBadRequest, "Missing required fields: user_id, role_id")
		return
	}

	// Update user-role assignment
	updateQuery := `UPDATE user_role SET expires_at = ?, updated_at = NOW() WHERE user_id = ? AND role_id = ? AND tenant_id = ?`

	result, err := h.db.ExecContext(r.Context(), updateQuery, req.ExpiresAt, req.UserID, req.RoleID, tenantID)
	if err != nil {
		h.logger.Error("Failed to update user role", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		h.respondError(w, http.StatusNotFound, "User-role assignment not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, "User role updated successfully", map[string]interface{}{
		"user_id": req.UserID,
		"role_id": req.RoleID,
	})
}
