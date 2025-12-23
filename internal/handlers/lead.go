package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"vyomtech-backend/internal/middleware"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
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
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	Phone           string  `json:"phone"`
	CompanyName     string  `json:"company_name"`
	Industry        string  `json:"industry,omitempty"`
	Status          string  `json:"status,omitempty"`
	DetailedStatus  string  `json:"detailed_status,omitempty"`
	PipelineStage   string  `json:"pipeline_stage,omitempty"`
	Probability     float64 `json:"probability,omitempty"`
	Source          string  `json:"source"`
	AssignedTo      string  `json:"assigned_to,omitempty"`
	NextActionDate  string  `json:"next_action_date,omitempty"`
	NextActionNotes string  `json:"next_action_notes,omitempty"`
}

// UpdateLeadRequest is the request body for updating a lead
type UpdateLeadRequest struct {
	FirstName       string  `json:"first_name,omitempty"`
	LastName        string  `json:"last_name,omitempty"`
	Email           string  `json:"email,omitempty"`
	Phone           string  `json:"phone,omitempty"`
	CompanyName     string  `json:"company_name,omitempty"`
	Industry        string  `json:"industry,omitempty"`
	Status          string  `json:"status,omitempty"`
	DetailedStatus  string  `json:"detailed_status,omitempty"`
	PipelineStage   string  `json:"pipeline_stage,omitempty"`
	Probability     float64 `json:"probability,omitempty"`
	Source          string  `json:"source,omitempty"`
	AssignedTo      string  `json:"assigned_to,omitempty"`
	NextActionDate  string  `json:"next_action_date,omitempty"`
	NextActionNotes string  `json:"next_action_notes,omitempty"`
}

// UpdateLeadStatusRequest is the request body for updating a lead's status
type UpdateLeadStatusRequest struct {
	Status          string  `json:"status" binding:"required"`
	DetailedStatus  string  `json:"detailed_status,omitempty"`
	PipelineStage   string  `json:"pipeline_stage,omitempty"`
	Notes           string  `json:"notes,omitempty"`
	CaptureDate     *string `json:"capture_date,omitempty"`      // capture_date_a, b, c, or d
	CaptureDateType *string `json:"capture_date_type,omitempty"` // a, b, c, or d
}

// LeadStatusResponse is the response for status update operations
type LeadStatusResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
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
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
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
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		lh.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	// Get leadID directly as string from URL
	leadID := r.PathValue("id")
	if leadID == "" {
		http.Error(w, "invalid lead id", http.StatusBadRequest)
		return
	}

	lead, err := lh.leadService.GetLead(ctx, leadID, tenantID)
	if err != nil {
		lh.logger.Error("Failed to get lead", "leadID", leadID, "userID", userID, "error", err)
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
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		lh.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
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
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Phone == "" || req.Source == "" {
		http.Error(w, "first_name, last_name, email, phone, and source are required", http.StatusBadRequest)
		return
	}

	status := req.Status
	if status == "" {
		status = "new"
	}

	// Parse NextActionDate if provided
	var nextActionDate *time.Time
	if req.NextActionDate != "" {
		t, err := time.Parse(time.RFC3339, req.NextActionDate)
		if err != nil {
			http.Error(w, "invalid next_action_date format", http.StatusBadRequest)
			return
		}
		nextActionDate = &t
	}

	lead := &models.Lead{
		TenantID:    tenantID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Phone:       req.Phone,
		CompanyName: req.CompanyName,
		Industry: func(s string) *string {
			if s == "" {
				return nil
			}
			return &s
		}(req.Industry),
		Status:      status,
		Probability: req.Probability,
		Source:      req.Source,
		AssignedTo: func(s string) *string {
			if s == "" {
				return nil
			}
			return &s
		}(req.AssignedTo),
		NextActionDate:  nextActionDate,
		NextActionNotes: req.NextActionNotes,
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
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		lh.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	leadID := r.PathValue("id")
	if leadID == "" {
		http.Error(w, "lead id required", http.StatusBadRequest)
		return
	}

	// Get existing lead
	lead, err := lh.leadService.GetLead(ctx, leadID, tenantID)
	if err != nil {
		lh.logger.Error("Failed to get lead", "leadID", leadID, "error", err)
		http.Error(w, "lead not found", http.StatusNotFound)
		return
	}

	var req UpdateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Update fields that were provided
	if req.FirstName != "" {
		lead.FirstName = req.FirstName
	}
	if req.LastName != "" {
		lead.LastName = req.LastName
	}
	if req.Email != "" {
		lead.Email = req.Email
	}
	if req.Phone != "" {
		lead.Phone = req.Phone
	}
	if req.CompanyName != "" {
		lead.CompanyName = req.CompanyName
	}
	if req.Industry != "" {
		lead.Industry = &req.Industry
	}
	if req.Status != "" {
		lead.Status = req.Status
	}
	if req.Probability > 0 {
		lead.Probability = req.Probability
	}
	if req.Source != "" {
		lead.Source = req.Source
	}
	if req.AssignedTo != "" {
		lead.AssignedTo = &req.AssignedTo
	}
	if req.NextActionDate != "" {
		t, err := time.Parse(time.RFC3339, req.NextActionDate)
		if err != nil {
			http.Error(w, "invalid next_action_date format", http.StatusBadRequest)
			return
		}
		lead.NextActionDate = &t
	}
	if req.NextActionNotes != "" {
		lead.NextActionNotes = req.NextActionNotes
	}

	if err := lh.leadService.UpdateLead(ctx, lead); err != nil {
		lh.logger.Error("Failed to update lead", "leadID", leadID, "userID", userID, "error", err)
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
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		lh.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	leadID := r.PathValue("id")
	if leadID == "" {
		http.Error(w, "lead id required", http.StatusBadRequest)
		return
	}

	if err := lh.leadService.DeleteLead(ctx, leadID, tenantID); err != nil {
		lh.logger.Error("Failed to delete lead", "leadID", leadID, "userID", userID, "error", err)
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
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
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

// UpdateLeadStatus updates a lead's status with validation and pipeline tracking
// PUT /api/v1/leads/{id}/status
func (lh *LeadHandler) UpdateLeadStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		lh.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	leadID := r.PathValue("id")
	if leadID == "" {
		http.Error(w, "lead id required", http.StatusBadRequest)
		return
	}

	var req UpdateLeadStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate status
	if !models.IsValidLeadStatus(req.Status) {
		http.Error(w, "invalid lead status", http.StatusBadRequest)
		return
	}

	// Get existing lead
	lead, err := lh.leadService.GetLead(ctx, leadID, tenantID)
	if err != nil {
		lh.logger.Error("Failed to get lead", "leadID", leadID, "error", err)
		http.Error(w, "lead not found", http.StatusNotFound)
		return
	}

	// Update status
	lead.Status = req.Status
	if req.Notes != "" {
		lead.NextActionNotes = req.Notes
	}

	if err := lh.leadService.UpdateLead(ctx, lead); err != nil {
		lh.logger.Error("Failed to update lead status", "leadID", leadID, "userID", userID, "error", err)
		http.Error(w, "failed to update lead status", http.StatusInternalServerError)
		return
	}

	response := LeadStatusResponse{
		ID:      lead.ID,
		Status:  lead.Status,
		Message: "status updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetLeadsByPipelineStage retrieves leads by pipeline stage
// GET /api/v1/leads/pipeline/{stage}
func (lh *LeadHandler) GetLeadsByPipelineStage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	stage := r.URL.Query().Get("stage")
	if stage == "" {
		http.Error(w, "pipeline stage required", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	filter := &models.LeadFilter{
		Limit:  25,
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

	leads, err := lh.leadService.GetLeadsByPipelineStage(ctx, tenantID, stage, filter)
	if err != nil {
		lh.logger.Error("Failed to get leads by pipeline stage", "stage", stage, "error", err)
		http.Error(w, "failed to get leads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leads)
}

// GetLeadsByStatus retrieves leads by detailed status
// GET /api/v1/leads/status/{status}
func (lh *LeadHandler) GetLeadsByStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		lh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		http.Error(w, "status required", http.StatusBadRequest)
		return
	}

	// Validate status
	if !models.IsValidLeadStatus(status) {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	filter := &models.LeadFilter{
		Status: status,
		Limit:  25,
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
		lh.logger.Error("Failed to get leads by status", "status", status, "error", err)
		http.Error(w, "failed to get leads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leads)
}
