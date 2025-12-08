package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/logger"
)

type TenantAdminHandler struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewTenantAdminHandler(db *sql.DB, logger *logger.Logger) *TenantAdminHandler {
	return &TenantAdminHandler{
		db:     db,
		logger: logger,
	}
}

// ListTenants handles GET /api/v1/tenants
func (h *TenantAdminHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(r.Context(),
		"SELECT id, name, domain, status, max_users, max_concurrent_calls, ai_budget_monthly, created_at, updated_at FROM tenant ORDER BY created_at DESC")
	if err != nil {
		h.logger.Error("Failed to list tenants", "error", err)
		http.Error(w, "Failed to list tenants", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tenants []models.Tenant
	for rows.Next() {
		var tenant models.Tenant
		if err := rows.Scan(&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.Status, &tenant.MaxUsers, &tenant.MaxConcurrentCalls, &tenant.AIBudgetMonthly, &tenant.CreatedAt, &tenant.UpdatedAt); err != nil {
			h.logger.Error("Failed to scan tenant", "error", err)
			continue
		}
		tenants = append(tenants, tenant)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenants)
}

// GetTenant handles GET /api/v1/tenants/:id
func (h *TenantAdminHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	var tenant models.Tenant
	err := h.db.QueryRowContext(r.Context(),
		"SELECT id, name, domain, status, max_users, max_concurrent_calls, ai_budget_monthly, created_at, updated_at FROM tenant WHERE id = ?",
		tenantID).Scan(&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.Status, &tenant.MaxUsers, &tenant.MaxConcurrentCalls, &tenant.AIBudgetMonthly, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Tenant not found", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to get tenant", "error", err)
			http.Error(w, "Failed to get tenant", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

// CreateTenantRequest defines the request structure for creating a tenant
type CreateTenantRequest struct {
	Name               string  `json:"name"`
	Domain             string  `json:"domain"`
	Status             string  `json:"status"`
	MaxUsers           int     `json:"max_users"`
	MaxConcurrentCalls int     `json:"max_concurrent_calls"`
	AIBudgetMonthly    float64 `json:"ai_budget_monthly"`
}

// CreateTenant handles POST /api/v1/tenants
func (h *TenantAdminHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var req CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Tenant name is required", http.StatusBadRequest)
		return
	}

	// Default values
	if req.Status == "" {
		req.Status = "active"
	}
	if req.MaxUsers == 0 {
		req.MaxUsers = 100
	}
	if req.MaxConcurrentCalls == 0 {
		req.MaxConcurrentCalls = 50
	}
	if req.AIBudgetMonthly == 0 {
		req.AIBudgetMonthly = 10000.00
	}

	// Generate tenant ID from name (simple approach)
	tenantID := "tenant_" + req.Name[:3] + "_" + req.Domain[:3]

	// Check if tenant already exists
	var existingID string
	err := h.db.QueryRowContext(r.Context(),
		"SELECT id FROM tenant WHERE id = ?",
		tenantID).Scan(&existingID)
	if err == nil {
		http.Error(w, "Tenant already exists", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		h.logger.Error("Database error", "error", err)
		http.Error(w, "Failed to create tenant", http.StatusInternalServerError)
		return
	}

	// Insert tenant
	_, err = h.db.ExecContext(r.Context(),
		"INSERT INTO tenant (id, name, domain, status, max_users, max_concurrent_calls, ai_budget_monthly, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())",
		tenantID, req.Name, req.Domain, req.Status, req.MaxUsers, req.MaxConcurrentCalls, req.AIBudgetMonthly)
	if err != nil {
		h.logger.Error("Failed to create tenant", "error", err)
		http.Error(w, "Failed to create tenant", http.StatusInternalServerError)
		return
	}

	tenant := models.Tenant{
		ID:                 tenantID,
		Name:               req.Name,
		Domain:             req.Domain,
		Status:             req.Status,
		MaxUsers:           req.MaxUsers,
		MaxConcurrentCalls: req.MaxConcurrentCalls,
		AIBudgetMonthly:    req.AIBudgetMonthly,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tenant)
}

// UpdateTenantRequest defines the request structure for updating a tenant
type UpdateTenantRequest struct {
	Name               string  `json:"name"`
	Domain             string  `json:"domain"`
	Status             string  `json:"status"`
	MaxUsers           int     `json:"max_users"`
	MaxConcurrentCalls int     `json:"max_concurrent_calls"`
	AIBudgetMonthly    float64 `json:"ai_budget_monthly"`
}

// UpdateTenant handles PUT /api/v1/tenants/:id
func (h *TenantAdminHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	var req UpdateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify tenant exists
	var existingID string
	err := h.db.QueryRowContext(r.Context(),
		"SELECT id FROM tenant WHERE id = ?",
		tenantID).Scan(&existingID)
	if err == sql.ErrNoRows {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error("Database error", "error", err)
		http.Error(w, "Failed to update tenant", http.StatusInternalServerError)
		return
	}

	// Update tenant
	_, err = h.db.ExecContext(r.Context(),
		"UPDATE tenant SET name = ?, domain = ?, status = ?, max_users = ?, max_concurrent_calls = ?, ai_budget_monthly = ?, updated_at = NOW() WHERE id = ?",
		req.Name, req.Domain, req.Status, req.MaxUsers, req.MaxConcurrentCalls, req.AIBudgetMonthly, tenantID)
	if err != nil {
		h.logger.Error("Failed to update tenant", "error", err)
		http.Error(w, "Failed to update tenant", http.StatusInternalServerError)
		return
	}

	// Return updated tenant
	var tenant models.Tenant
	err = h.db.QueryRowContext(r.Context(),
		"SELECT id, name, domain, status, max_users, max_concurrent_calls, ai_budget_monthly, created_at, updated_at FROM tenant WHERE id = ?",
		tenantID).Scan(&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.Status, &tenant.MaxUsers, &tenant.MaxConcurrentCalls, &tenant.AIBudgetMonthly, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		h.logger.Error("Failed to get updated tenant", "error", err)
		http.Error(w, "Failed to update tenant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

// DeleteTenant handles DELETE /api/v1/tenants/:id
func (h *TenantAdminHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	// Verify tenant exists
	var existingID string
	err := h.db.QueryRowContext(r.Context(),
		"SELECT id FROM tenant WHERE id = ?",
		tenantID).Scan(&existingID)
	if err == sql.ErrNoRows {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error("Database error", "error", err)
		http.Error(w, "Failed to delete tenant", http.StatusInternalServerError)
		return
	}

	// Delete tenant
	_, err = h.db.ExecContext(r.Context(),
		"DELETE FROM tenant WHERE id = ?",
		tenantID)
	if err != nil {
		h.logger.Error("Failed to delete tenant", "error", err)
		http.Error(w, "Failed to delete tenant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetTenantUsers handles GET /api/v1/tenants/:id/users
func (h *TenantAdminHandler) GetTenantUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["id"]

	rows, err := h.db.QueryContext(r.Context(),
		"SELECT id, email, role, tenant_id, created_at, updated_at FROM user WHERE tenant_id = ? ORDER BY created_at DESC",
		tenantID)
	if err != nil {
		h.logger.Error("Failed to get tenant users", "error", err)
		http.Error(w, "Failed to get tenant users", http.StatusInternalServerError)
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
