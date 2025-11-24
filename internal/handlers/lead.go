package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/internal/services"
	"multi-tenant-ai-callcenter/pkg/logger"
)

// LeadHandler handles lead-related HTTP requests
type LeadHandler struct {
	leadService *services.LeadService
	logger      *logger.Logger
}

// NewLeadHandler creates a new LeadHandler
func NewLeadHandler(leadService *services.LeadService, logger *logger.Logger) *LeadHandler {
	return &LeadHandler{
		leadService: leadService,
		logger:      logger,
	}
}

// CreateLeadRequest is the request body for creating a lead
type CreateLeadRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Status        string `json:"status,omitempty"`
	Source        string `json:"source"`
	CampaignID    *int64 `json:"campaign_id,omitempty"`
	AssignedAgent *int64 `json:"assigned_agent_id,omitempty"`
	Notes         string `json:"notes,omitempty"`
}

// UpdateLeadRequest is the request body for updating a lead
type UpdateLeadRequest struct {
	Name          string `json:"name,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Status        string `json:"status,omitempty"`
	Source        string `json:"source,omitempty"`
	CampaignID    *int64 `json:"campaign_id,omitempty"`
	AssignedAgent *int64 `json:"assigned_agent_id,omitempty"`
	Notes         string `json:"notes,omitempty"`
}

// ListLeadsRequest is the query parameters for listing leads
type ListLeadsRequest struct {
	Status     string
	Source     string
	CampaignID int64
	Limit      int
	Offset     int
}

// GetLeads retrieves all leads for the tenant
// GET /api/v1/leads
func (lh *LeadHandler) GetLeads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	filter := &models.LeadFilter{
		Status: r.URL.Query().Get("status"),
		Source: r.URL.Query().Get("source"),
		Limit:  10,
		Offset: 0,
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if parsed, err := strconv.Atoi(limit); err == nil && parsed > 0 {
			filter.Limit = parsed
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if parsed, err := strconv.Atoi(offset); err == nil && parsed >= 0 {
			filter.Offset = parsed
		}
	}

	leads, err := lh.leadService.GetLeads(ctx, tenantID, filter)
	if err != nil {
		lh.logger.Error("Failed to get leads", "error", err)
		http.Error(w, "failed to get leads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leads)
}

// GetLead retrieves a specific lead
// GET /api/v1/leads/{id}
func (lh *LeadHandler) GetLead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		lh.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	leadID := r.URL.Query().Get("id")
	if leadID == "" {
		http.Error(w, "lead id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(leadID, 10, 64)
	if err != nil {
		http.Error(w, "invalid lead id", http.StatusBadRequest)
		return
	}

	lead, err := lh.leadService.GetLead(ctx, id, tenantID)
	if err != nil {
		lh.logger.Error("Failed to get lead", "leadID", id, "userID", userID, "error", err)
		http.Error(w, "lead not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lead)
}

// CreateLead creates a new lead
// POST /api/v1/leads
func (lh *LeadHandler) CreateLead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		lh.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req CreateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" || req.Phone == "" || req.Source == "" {
		http.Error(w, "name, email, phone, and source are required", http.StatusBadRequest)
		return
	}

	status := req.Status
	if status == "" {
		status = "new"
	}

	lead := &models.Lead{
		TenantID:      tenantID,
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Status:        status,
		Source:        req.Source,
		CampaignID:    req.CampaignID,
		AssignedAgent: req.AssignedAgent,
		Notes:         req.Notes,
	}

	if err := lh.leadService.CreateLead(ctx, lead); err != nil {
		lh.logger.Error("Failed to create lead", "userID", userID, "error", err)
		http.Error(w, "failed to create lead", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lead)
}

// UpdateLead updates an existing lead
// PUT /api/v1/leads/{id}
func (lh *LeadHandler) UpdateLead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		lh.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	leadID := r.URL.Query().Get("id")
	if leadID == "" {
		http.Error(w, "lead id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(leadID, 10, 64)
	if err != nil {
		http.Error(w, "invalid lead id", http.StatusBadRequest)
		return
	}

	// Get existing lead
	lead, err := lh.leadService.GetLead(ctx, id, tenantID)
	if err != nil {
		lh.logger.Error("Failed to get lead", "leadID", id, "error", err)
		http.Error(w, "lead not found", http.StatusNotFound)
		return
	}

	var req UpdateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Update fields that were provided
	if req.Name != "" {
		lead.Name = req.Name
	}
	if req.Email != "" {
		lead.Email = req.Email
	}
	if req.Phone != "" {
		lead.Phone = req.Phone
	}
	if req.Status != "" {
		lead.Status = req.Status
	}
	if req.Source != "" {
		lead.Source = req.Source
	}
	if req.CampaignID != nil {
		lead.CampaignID = req.CampaignID
	}
	if req.AssignedAgent != nil {
		lead.AssignedAgent = req.AssignedAgent
	}
	if req.Notes != "" {
		lead.Notes = req.Notes
	}

	if err := lh.leadService.UpdateLead(ctx, lead); err != nil {
		lh.logger.Error("Failed to update lead", "leadID", id, "userID", userID, "error", err)
		http.Error(w, "failed to update lead", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lead)
}

// DeleteLead deletes a lead
// DELETE /api/v1/leads/{id}
func (lh *LeadHandler) DeleteLead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		lh.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	leadID := r.URL.Query().Get("id")
	if leadID == "" {
		http.Error(w, "lead id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(leadID, 10, 64)
	if err != nil {
		http.Error(w, "invalid lead id", http.StatusBadRequest)
		return
	}

	if err := lh.leadService.DeleteLead(ctx, id, tenantID); err != nil {
		lh.logger.Error("Failed to delete lead", "leadID", id, "userID", userID, "error", err)
		http.Error(w, "lead not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "lead deleted successfully"})
}

// GetLeadStats retrieves lead statistics
// GET /api/v1/leads/stats
func (lh *LeadHandler) GetLeadStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	stats, err := lh.leadService.GetLeadStats(ctx, tenantID)
	if err != nil {
		lh.logger.Error("Failed to get lead stats", "error", err)
		http.Error(w, "failed to get lead stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
