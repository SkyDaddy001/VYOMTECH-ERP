package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

// RoleTemplateHandler handles role template operations
type RoleTemplateHandler struct {
	db     *sql.DB
	logger *logger.Logger
	svc    *services.RoleTemplateService
	rbac   *services.RBACService
}

// NewRoleTemplateHandler creates a new role template handler
func NewRoleTemplateHandler(db *sql.DB, log *logger.Logger, rbacService *services.RBACService) *RoleTemplateHandler {
	return &RoleTemplateHandler{
		db:     db,
		logger: log,
		svc:    services.NewRoleTemplateService(db, rbacService, log),
		rbac:   rbacService,
	}
}

// CreateRoleTemplateRequest represents the request to create a custom role template
type CreateRoleTemplateRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Permissions []string `json:"permissions"`
}

// CreateRoleTemplate creates a new custom role template
func (h *RoleTemplateHandler) CreateRoleTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)

	var req CreateRoleTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	template := &models.RoleTemplate{
		TenantID:      tenantID,
		Name:          req.Name,
		Description:   req.Description,
		Category:      req.Category,
		PermissionIDs: req.Permissions,
		IsActive:      true,
	}

	if err := h.svc.CreateCustomTemplate(r.Context(), template); err != nil {
		h.logger.Error("Error creating role template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

// GetRoleTemplate retrieves a role template by ID
func (h *RoleTemplateHandler) GetRoleTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)
	templateID := r.URL.Query().Get("id")

	template, err := h.svc.GetTemplate(r.Context(), tenantID, templateID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Error getting role template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

// ListRoleTemplates lists all role templates for a tenant
func (h *RoleTemplateHandler) ListRoleTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)

	templates, err := h.svc.ListTemplates(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("Error listing role templates", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

// CreateRoleFromTemplateRequest for creating a role from a template
type CreateRoleFromTemplateRequest struct {
	TemplateID string `json:"template_id" binding:"required"`
	RoleName   string `json:"role_name" binding:"required"`
}

// CreateRoleFromTemplate creates a new role from a template
func (h *RoleTemplateHandler) CreateRoleFromTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)

	var req CreateRoleFromTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	roleID, err := h.svc.CreateRoleFromTemplate(r.Context(), tenantID, req.TemplateID, req.RoleName)
	if err != nil {
		h.logger.Error("Error creating role from template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"role_id": roleID,
		"message": "Role created successfully from template",
	})
}

// InitializeDefaultTemplates initializes default templates for a tenant
func (h *RoleTemplateHandler) InitializeDefaultTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)

	if err := h.svc.InitializeDefaultTemplates(r.Context(), tenantID); err != nil {
		h.logger.Error("Error initializing default templates", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Default templates initialized successfully",
	})
}

// SaveAsTemplateRequest for saving a role as a template
type SaveAsTemplateRequest struct {
	RoleID       string `json:"role_id" binding:"required"`
	TemplateName string `json:"template_name" binding:"required"`
	Description  string `json:"description"`
	Category     string `json:"category"`
}

// SaveAsTemplate saves an existing role as a template
func (h *RoleTemplateHandler) SaveAsTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)

	var req SaveAsTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	templateID, err := h.svc.SaveAsTemplate(r.Context(), tenantID, req.RoleID, req.TemplateName, req.Description, req.Category)
	if err != nil {
		h.logger.Error("Error saving role as template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"template_id": templateID,
		"message":     "Role saved as template successfully",
	})
}

// GetDefaultTemplates returns the default role templates
func (h *RoleTemplateHandler) GetDefaultTemplates(w http.ResponseWriter, r *http.Request) {
	templates := []models.RoleTemplate{}
	for _, template := range models.DefaultTemplates {
		template.IsSystemTemplate = true
		template.IsActive = true
		templates = append(templates, template)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}
