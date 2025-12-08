package handlers

import (
	"encoding/json"
	"net/http"

	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

type TenantHandler struct {
	tenantService *services.TenantService
	logger        *logger.Logger
}

func NewTenantHandler(tenantService *services.TenantService, logger *logger.Logger) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
		logger:        logger,
	}
}

// ListTenants lists all tenants (admin only)
func (h *TenantHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenants, err := h.tenantService.ListTenants(ctx)
	if err != nil {
		h.logger.Error("Failed to list tenants", "error", err)
		http.Error(w, "Failed to list tenants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenants)
}

// GetTenantInfo gets current tenant information
func (h *TenantHandler) GetTenantInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get tenant ID from context
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "Tenant ID not found in context", http.StatusBadRequest)
		return
	}

	tenant, err := h.tenantService.GetTenant(ctx, tenantID)
	if err != nil {
		h.logger.Warn("Tenant not found", "tenant_id", tenantID)
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenant)
}

// GetTenantUserCount gets the number of users in a tenant
func (h *TenantHandler) GetTenantUserCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "Tenant ID not found in context", http.StatusBadRequest)
		return
	}

	count, err := h.tenantService.GetTenantUserCount(ctx, tenantID)
	if err != nil {
		h.logger.Error("Failed to get user count", "error", err)
		http.Error(w, "Failed to get user count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tenant_id":  tenantID,
		"user_count": count,
	})
}

// SwitchTenant switches the user's active tenant
func (h *TenantHandler) SwitchTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// Get user ID and new tenant ID from context
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}

	// Extract tenant ID from URL path
	tenantID := r.PathValue("id")
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Switch tenant
	err := h.tenantService.SwitchUserTenant(ctx, userID, tenantID)
	if err != nil {
		h.logger.Warn("Failed to switch tenant", "user_id", userID, "tenant_id", tenantID, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get tenant info to return
	tenant, err := h.tenantService.GetTenant(ctx, tenantID)
	if err != nil {
		h.logger.Error("Failed to get tenant", "error", err)
		http.Error(w, "Failed to get tenant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tenant switched successfully",
		"tenant":  tenant,
	})
}

// AddTenantMember adds a member to a tenant
func (h *TenantHandler) AddTenantMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// Get user ID for authorization check
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}

	tenantID := r.PathValue("id")
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// Check if user is admin
	isAdmin, err := h.tenantService.UserIsTenantAdmin(ctx, userID, tenantID)
	if err != nil {
		h.logger.Error("Failed to check admin status", "error", err)
		http.Error(w, "Failed to check permissions", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Only admins can add members", http.StatusForbidden)
		return
	}

	// Parse request
	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Role == "" {
		http.Error(w, "Email and role are required", http.StatusBadRequest)
		return
	}

	// For now, we'll create the member record
	// In production, you'd look up the user by email and add them properly
	// This is a simplified implementation
	member, err := h.tenantService.AddTenantMember(ctx, tenantID, 0, req.Email, req.Role)
	if err != nil {
		h.logger.Error("Failed to add tenant member", "error", err)
		http.Error(w, "Failed to add member", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User added to tenant",
		"member":  member,
	})
}

// RemoveTenantMember removes a member from a tenant
func (h *TenantHandler) RemoveTenantMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// Get user ID for authorization check
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}

	tenantID := r.PathValue("id")
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	email := r.PathValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Check if user is admin
	isAdmin, err := h.tenantService.UserIsTenantAdmin(ctx, userID, tenantID)
	if err != nil {
		h.logger.Error("Failed to check admin status", "error", err)
		http.Error(w, "Failed to check permissions", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Only admins can remove members", http.StatusForbidden)
		return
	}

	// Get user ID by email
	memberUserID, err := h.tenantService.GetUserIDByEmail(ctx, email)
	if err != nil {
		h.logger.Warn("Failed to find user by email", "email", email, "error", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Remove the member
	err = h.tenantService.RemoveTenantMember(ctx, tenantID, memberUserID)
	if err != nil {
		h.logger.Error("Failed to remove tenant member", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Member removed successfully",
		"email":   email,
	})
}

// GetUserTenants gets all tenants for the current user
func (h *TenantHandler) GetUserTenants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}

	tenants, err := h.tenantService.GetUserTenants(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to get user tenants", "error", err)
		http.Error(w, "Failed to get tenants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenants)
}

// CreateTenant creates a new tenant (used during registration)
func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type CreateTenantRequest struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
	}

	var req CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Tenant name is required", http.StatusBadRequest)
		return
	}

	tenant, err := h.tenantService.CreateTenant(ctx, req.Name, req.Domain)
	if err != nil {
		h.logger.Error("Failed to create tenant", "error", err)
		http.Error(w, "Failed to create tenant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tenant)
}
