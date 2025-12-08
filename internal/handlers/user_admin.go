package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/logger"
)

type UserAdminHandler struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewUserAdminHandler(db *sql.DB, logger *logger.Logger) *UserAdminHandler {
	return &UserAdminHandler{
		db:     db,
		logger: logger,
	}
}

// ListUsers handles GET /api/v1/users
func (h *UserAdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || userRole == "" {
		http.Error(w, "User role not found", http.StatusUnauthorized)
		return
	}

	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	// Master admin can see users from all tenants, tenant admin can only see their own
	var query string
	var args []interface{}

	if userRole == "master_admin" {
		query = "SELECT id, email, role, tenant_id, created_at, updated_at FROM user ORDER BY created_at DESC"
	} else {
		// Tenant admin/regular admin can only see their own tenant's users
		query = "SELECT id, email, role, tenant_id, created_at, updated_at FROM user WHERE tenant_id = ? ORDER BY created_at DESC"
		args = append(args, tenantID)
	}

	rows, err := h.db.QueryContext(r.Context(), query, args...)
	if err != nil {
		h.logger.Error("Failed to list users", "error", err)
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Role, &user.TenantID, &user.CreatedAt, &user.UpdatedAt); err != nil {
			h.logger.Error("Failed to scan user", "error", err)
			continue
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUser handles GET /api/v1/users/:id
func (h *UserAdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || userRole == "" {
		http.Error(w, "User role not found", http.StatusUnauthorized)
		return
	}

	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	var user models.User
	var query string
	var args []interface{}

	if userRole == "master_admin" {
		// Master admin can get any user
		query = "SELECT id, email, role, tenant_id, created_at, updated_at FROM user WHERE id = ?"
		args = append(args, userID)
	} else {
		// Tenant admin can only get users from their tenant
		query = "SELECT id, email, role, tenant_id, created_at, updated_at FROM user WHERE id = ? AND tenant_id = ?"
		args = append(args, userID, tenantID)
	}

	err := h.db.QueryRowContext(r.Context(), query, args...).Scan(&user.ID, &user.Email, &user.Role, &user.TenantID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found or access denied", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to get user", "error", err)
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUserRequest defines the request structure for creating a user
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	TenantID string `json:"tenant_id"`
}

// CreateUser handles POST /api/v1/users
func (h *UserAdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || userRole == "" {
		http.Error(w, "User role not found", http.StatusUnauthorized)
		return
	}

	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Default role
	if req.Role == "" {
		req.Role = "agent"
	}

	// Authorization check: master_admin can create users for any tenant, tenant admin only for their own
	targetTenantID := tenantID
	if req.TenantID != "" {
		if userRole != "master_admin" {
			// Non-master admins cannot create users for other tenants
			if req.TenantID != tenantID {
				http.Error(w, "You can only create users for your own tenant", http.StatusForbidden)
				return
			}
		}
		targetTenantID = req.TenantID
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("Failed to hash password", "error", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Check if user already exists in the target tenant
	var existingID int
	err = h.db.QueryRowContext(r.Context(),
		"SELECT id FROM user WHERE email = ? AND tenant_id = ?",
		req.Email, targetTenantID).Scan(&existingID)
	if err == nil {
		http.Error(w, "User already exists in this tenant", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		h.logger.Error("Database error checking user existence", "error", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Insert user
	result, err := h.db.ExecContext(r.Context(),
		"INSERT INTO user (email, password_hash, role, tenant_id, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())",
		req.Email, string(hashedPassword), req.Role, targetTenantID)
	if err != nil {
		h.logger.Error("Failed to create user", "error", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		h.logger.Error("Failed to get user ID", "error", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:       int(userID),
		Email:    req.Email,
		Role:     req.Role,
		TenantID: targetTenantID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UpdateUserRequest defines the request structure for updating a user
type UpdateUserRequest struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// UpdateUser handles PUT /api/v1/users/:id
func (h *UserAdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || userRole == "" {
		http.Error(w, "User role not found", http.StatusUnauthorized)
		return
	}

	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify user exists and belongs to tenant (or master admin can update any user)
	var userTenantID string
	var existingID int
	query := "SELECT id, tenant_id FROM user WHERE id = ?"
	args := []interface{}{userID}

	if userRole != "master_admin" {
		query += " AND tenant_id = ?"
		args = append(args, tenantID)
	}

	row := h.db.QueryRowContext(r.Context(), query, args...)
	err := row.Scan(&existingID, &userTenantID)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found or access denied", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error("Database error", "error", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Update user
	_, err = h.db.ExecContext(r.Context(),
		"UPDATE user SET role = ?, updated_at = NOW() WHERE id = ?",
		req.Role, userID)
	if err != nil {
		h.logger.Error("Failed to update user", "error", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Return updated user
	var user models.User
	err = h.db.QueryRowContext(r.Context(),
		"SELECT id, email, role, tenant_id, created_at, updated_at FROM user WHERE id = ?",
		userID).Scan(&user.ID, &user.Email, &user.Role, &user.TenantID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		h.logger.Error("Failed to get updated user", "error", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUser handles DELETE /api/v1/users/:id
func (h *UserAdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || userRole == "" {
		http.Error(w, "User role not found", http.StatusUnauthorized)
		return
	}

	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	// Verify user exists and belongs to tenant (or master admin can delete any user)
	var userTenantID string
	var existingID int
	query := "SELECT id, tenant_id FROM user WHERE id = ?"
	args := []interface{}{userID}

	if userRole != "master_admin" {
		query += " AND tenant_id = ?"
		args = append(args, tenantID)
	}

	row := h.db.QueryRowContext(r.Context(), query, args...)
	err := row.Scan(&existingID, &userTenantID)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found or access denied", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error("Database error", "error", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Check authorization for non-master admins
	if userRole != "master_admin" && userTenantID != tenantID {
		http.Error(w, "You can only delete users from your own tenant", http.StatusForbidden)
		return
	}

	// Delete user
	_, err = h.db.ExecContext(r.Context(),
		"DELETE FROM user WHERE id = ?",
		userID)
	if err != nil {
		h.logger.Error("Failed to delete user", "error", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateUserRoleRequest defines the request structure for updating user role
type UpdateUserRoleRequest struct {
	Role string `json:"role"`
}

// UpdateUserRole handles PUT /api/v1/users/:id/role
func (h *UserAdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || userRole == "" {
		http.Error(w, "User role not found", http.StatusUnauthorized)
		return
	}

	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	var req UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify user exists and belongs to tenant (or master admin can update any user)
	var userTenantID string
	var existingID int
	query := "SELECT id, tenant_id FROM user WHERE id = ?"
	args := []interface{}{userID}

	if userRole != "master_admin" {
		query += " AND tenant_id = ?"
		args = append(args, tenantID)
	}

	row := h.db.QueryRowContext(r.Context(), query, args...)
	err := row.Scan(&existingID, &userTenantID)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found or access denied", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error("Database error", "error", err)
		http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		return
	}

	// Check authorization for non-master admins
	if userRole != "master_admin" && userTenantID != tenantID {
		http.Error(w, "You can only update users from your own tenant", http.StatusForbidden)
		return
	}

	// Update user role
	_, err = h.db.ExecContext(r.Context(),
		"UPDATE user SET role = ?, updated_at = NOW() WHERE id = ?",
		req.Role, userID)
	if err != nil {
		h.logger.Error("Failed to update user role", "error", err)
		http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Role updated successfully"})
}

// ResetPasswordRequest defines the request structure for resetting password
type ResetPasswordRequest struct {
	Password string `json:"password"`
}

// ResetPassword handles POST /api/v1/users/:id/reset-password
func (h *UserAdminHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || userRole == "" {
		http.Error(w, "User role not found", http.StatusUnauthorized)
		return
	}

	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Verify user exists and belongs to tenant (or master admin can reset any user)
	var userTenantID string
	var existingID int
	query := "SELECT id, tenant_id FROM user WHERE id = ?"
	args := []interface{}{userID}

	if userRole != "master_admin" {
		query += " AND tenant_id = ?"
		args = append(args, tenantID)
	}

	row := h.db.QueryRowContext(r.Context(), query, args...)
	err := row.Scan(&existingID, &userTenantID)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found or access denied", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error("Database error", "error", err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	// Check authorization for non-master admins
	if userRole != "master_admin" && userTenantID != tenantID {
		http.Error(w, "You can only reset passwords for users in your own tenant", http.StatusForbidden)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("Failed to hash password", "error", err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	// Update password
	result, err := h.db.ExecContext(r.Context(),
		"UPDATE user SET password_hash = ?, updated_at = NOW() WHERE id = ?",
		string(hashedPassword), userID)
	if err != nil {
		h.logger.Error("Failed to reset password", "error", err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		h.logger.Error("Failed to get rows affected", "error", err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successfully"})
}
