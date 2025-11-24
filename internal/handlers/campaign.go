package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/internal/services"
	"multi-tenant-ai-callcenter/pkg/logger"
)

// CampaignHandler handles campaign-related HTTP requests
type CampaignHandler struct {
	campaignService *services.CampaignService
	logger          *logger.Logger
}

// NewCampaignHandler creates a new CampaignHandler
func NewCampaignHandler(campaignService *services.CampaignService, logger *logger.Logger) *CampaignHandler {
	return &CampaignHandler{
		campaignService: campaignService,
		logger:          logger,
	}
}

// CreateCampaignRequest is the request body for creating a campaign
type CreateCampaignRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status,omitempty"`
	TargetLeads int       `json:"target_leads"`
	Budget      float64   `json:"budget"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

// UpdateCampaignRequest is the request body for updating a campaign
type UpdateCampaignRequest struct {
	Name           string     `json:"name,omitempty"`
	Description    string     `json:"description,omitempty"`
	Status         string     `json:"status,omitempty"`
	TargetLeads    int        `json:"target_leads,omitempty"`
	GeneratedLeads int        `json:"generated_leads,omitempty"`
	ConvertedLeads int        `json:"converted_leads,omitempty"`
	Budget         float64    `json:"budget,omitempty"`
	SpentBudget    float64    `json:"spent_budget,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
}

// GetCampaigns retrieves all campaigns for the tenant
// GET /api/v1/campaigns
func (ch *CampaignHandler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	filter := &models.CampaignFilter{
		Status: r.URL.Query().Get("status"),
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

	campaigns, err := ch.campaignService.GetCampaigns(ctx, tenantID, filter)
	if err != nil {
		ch.logger.Error("Failed to get campaigns", "error", err)
		http.Error(w, "failed to get campaigns", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(campaigns)
}

// GetCampaign retrieves a specific campaign
// GET /api/v1/campaigns/{id}
func (ch *CampaignHandler) GetCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	campaignID := r.URL.Query().Get("id")
	if campaignID == "" {
		http.Error(w, "campaign id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(campaignID, 10, 64)
	if err != nil {
		http.Error(w, "invalid campaign id", http.StatusBadRequest)
		return
	}

	campaign, err := ch.campaignService.GetCampaign(ctx, id, tenantID)
	if err != nil {
		ch.logger.Error("Failed to get campaign", "campaignID", id, "error", err)
		http.Error(w, "campaign not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(campaign)
}

// CreateCampaign creates a new campaign
// POST /api/v1/campaigns
func (ch *CampaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		ch.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Budget <= 0 || req.TargetLeads <= 0 {
		http.Error(w, "name, budget, and target_leads are required", http.StatusBadRequest)
		return
	}

	status := req.Status
	if status == "" {
		status = "planned"
	}

	campaign := &models.Campaign{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
		TargetLeads: req.TargetLeads,
		Budget:      req.Budget,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := ch.campaignService.CreateCampaign(ctx, campaign); err != nil {
		ch.logger.Error("Failed to create campaign", "userID", userID, "error", err)
		http.Error(w, "failed to create campaign", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(campaign)
}

// UpdateCampaign updates an existing campaign
// PUT /api/v1/campaigns/{id}
func (ch *CampaignHandler) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		ch.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	campaignID := r.URL.Query().Get("id")
	if campaignID == "" {
		http.Error(w, "campaign id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(campaignID, 10, 64)
	if err != nil {
		http.Error(w, "invalid campaign id", http.StatusBadRequest)
		return
	}

	// Get existing campaign
	campaign, err := ch.campaignService.GetCampaign(ctx, id, tenantID)
	if err != nil {
		ch.logger.Error("Failed to get campaign", "campaignID", id, "error", err)
		http.Error(w, "campaign not found", http.StatusNotFound)
		return
	}

	var req UpdateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Update fields that were provided
	if req.Name != "" {
		campaign.Name = req.Name
	}
	if req.Description != "" {
		campaign.Description = req.Description
	}
	if req.Status != "" {
		campaign.Status = req.Status
	}
	if req.TargetLeads > 0 {
		campaign.TargetLeads = req.TargetLeads
	}
	if req.GeneratedLeads > 0 {
		campaign.GeneratedLeads = req.GeneratedLeads
	}
	if req.ConvertedLeads > 0 {
		campaign.ConvertedLeads = req.ConvertedLeads
	}
	if req.Budget > 0 {
		campaign.Budget = req.Budget
	}
	if req.SpentBudget > 0 {
		campaign.SpentBudget = req.SpentBudget
	}
	if req.EndDate != nil {
		campaign.EndDate = *req.EndDate
	}

	// Recalculate metrics
	if campaign.GeneratedLeads > 0 {
		campaign.CostPerLead = campaign.SpentBudget / float64(campaign.GeneratedLeads)
	}
	if campaign.GeneratedLeads > 0 {
		campaign.ConversionRate = float64(campaign.ConvertedLeads) / float64(campaign.GeneratedLeads)
	}

	if err := ch.campaignService.UpdateCampaign(ctx, campaign); err != nil {
		ch.logger.Error("Failed to update campaign", "campaignID", id, "userID", userID, "error", err)
		http.Error(w, "failed to update campaign", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(campaign)
}

// DeleteCampaign deletes a campaign
// DELETE /api/v1/campaigns/{id}
func (ch *CampaignHandler) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		ch.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	campaignID := r.URL.Query().Get("id")
	if campaignID == "" {
		http.Error(w, "campaign id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(campaignID, 10, 64)
	if err != nil {
		http.Error(w, "invalid campaign id", http.StatusBadRequest)
		return
	}

	if err := ch.campaignService.DeleteCampaign(ctx, id, tenantID); err != nil {
		ch.logger.Error("Failed to delete campaign", "campaignID", id, "userID", userID, "error", err)
		http.Error(w, "campaign not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "campaign deleted successfully"})
}

// GetCampaignStats retrieves campaign statistics
// GET /api/v1/campaigns/stats
func (ch *CampaignHandler) GetCampaignStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	stats, err := ch.campaignService.GetCampaignStats(ctx, tenantID)
	if err != nil {
		ch.logger.Error("Failed to get campaign stats", "error", err)
		http.Error(w, "failed to get campaign stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
