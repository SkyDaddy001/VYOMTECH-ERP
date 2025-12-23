package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
)

// RoleTemplateHandler handles role template operations
type RoleTemplateHandler struct {
	db     *sql.DB
	logger *log.Logger
	svc    *services.RoleTemplateService
}

// NewRoleTemplateHandler creates a new role template handler
func NewRoleTemplateHandler(db *sql.DB, logger *log.Logger) *RoleTemplateHandler {
	return &RoleTemplateHandler{
		db:     db,
		logger: logger,
		svc:    services.NewRoleTemplateService(db, logger),
	}
}

// CreateRoleTemplateRequest represents the request to create a role template
type CreateRoleTemplateRequest struct {
	Name             string                 `json:"name" binding:"required"`
	Description      string                 `json:"description"`
	Category         string                 `json:"category"`
	IsSystemTemplate bool                   `json:"is_system_template"`
	PermissionIDs    []string               `json:"permission_ids"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// CreateRoleTemplate creates a new role template
func (h *RoleTemplateHandler) CreateRoleTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)

	var req CreateRoleTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metadataJSON, _ := json.Marshal(req.Metadata)
	template := &models.RoleTemplate{
		TenantID:         tenantID,
		Name:             req.Name,
		Description:      req.Description,
		Category:         req.Category,
		IsSystemTemplate: req.IsSystemTemplate,
		PermissionIDs:    req.PermissionIDs,
		Metadata:         string(metadataJSON),
		IsActive:         true,
	}

	if err := h.svc.CreateTemplate(r.Context(), template); err != nil {
		h.logger.Printf("Error creating role template: %v", err)
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
		h.logger.Printf("Error getting role template: %v", err)
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
		h.logger.Printf("Error listing role templates: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

// UpdateRoleTemplate updates a role template
func (h *RoleTemplateHandler) UpdateRoleTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)

	var req CreateRoleTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	templateID := r.URL.Query().Get("id")
	metadataJSON, _ := json.Marshal(req.Metadata)

	template := &models.RoleTemplate{
		ID:               templateID,
		TenantID:         tenantID,
		Name:             req.Name,
		Description:      req.Description,
		Category:         req.Category,
		IsSystemTemplate: req.IsSystemTemplate,
		PermissionIDs:    req.PermissionIDs,
		Metadata:         string(metadataJSON),
		IsActive:         true,
	}

	if err := h.svc.UpdateTemplate(r.Context(), template); err != nil {
		h.logger.Printf("Error updating role template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

// DeleteRoleTemplate deletes a role template
func (h *RoleTemplateHandler) DeleteRoleTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)
	templateID := r.URL.Query().Get("id")

	if err := h.svc.DeleteTemplate(r.Context(), tenantID, templateID); err != nil {
		h.logger.Printf("Error deleting role template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// CreateRoleFromTemplate creates a new role from a template
type CreateRoleFromTemplateRequest struct {
	TemplateID     string                 `json:"template_id" binding:"required"`
	RoleName       string                 `json:"role_name" binding:"required"`
	Customizations map[string]interface{} `json:"customizations"`
}

// CreateRoleFromTemplate creates a new role from a template
func (h *RoleTemplateHandler) CreateRoleFromTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)
	userID := r.Context().Value("user_id").(string)

	var req CreateRoleFromTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customizationsJSON, _ := json.Marshal(req.Customizations)
	instance := &models.TemplateInstance{
		TenantID:       tenantID,
		TemplateID:     req.TemplateID,
		Customizations: string(customizationsJSON),
	}

	// Parse user ID and convert to int64
	var createdBy int64
	json.Unmarshal([]byte(`"`+userID+`"`), &createdBy)
	instance.CreatedBy = createdBy

	if err := h.svc.CreateInstanceFromTemplate(r.Context(), instance, req.RoleName); err != nil {
		h.logger.Printf("Error creating role from template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instance)
}

// ListTemplateInstances lists all template instances for a tenant
func (h *RoleTemplateHandler) ListTemplateInstances(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value("tenant_id").(string)

	instances, err := h.svc.ListInstances(r.Context(), tenantID)
	if err != nil {
		h.logger.Printf("Error listing template instances: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instances)
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
