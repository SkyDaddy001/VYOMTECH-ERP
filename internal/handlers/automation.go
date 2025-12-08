package handlers

import (
	"vyomtech-backend/internal/middleware"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

// AutomationHandler handles automation and routing requests
type AutomationHandler struct {
	automationService *services.AutomationService
	logger            *logger.Logger
}

// NewAutomationHandler creates a new AutomationHandler
func NewAutomationHandler(automationService *services.AutomationService, logger *logger.Logger) *AutomationHandler {
	return &AutomationHandler{
		automationService: automationService,
		logger:            logger,
	}
}

// CreateRoutingRuleRequest is the request body for creating a routing rule
type CreateRoutingRuleRequest struct {
	Name        string                 `json:"name"`
	Priority    int                    `json:"priority"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	ActionType  string                 `json:"action_type"` // assign_to_agent, assign_to_team, round_robin
	ActionValue string                 `json:"action_value"`
	Enabled     bool                   `json:"enabled"`
}

// ScheduleCampaignRequest is the request body for scheduling a campaign
type ScheduleCampaignRequest struct {
	CampaignID  int64  `json:"campaign_id"`
	ScheduledAt string `json:"scheduled_at"` // ISO8601 format
}

// CalculateLeadScore calculates the score for a lead
// POST /api/v1/automation/leads/{id}/score
func (ah *AutomationHandler) CalculateLeadScore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		ah.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	leadIDStr := r.URL.Query().Get("id")
	if leadIDStr == "" {
		http.Error(w, "lead id required", http.StatusBadRequest)
		return
	}

	leadID, err := strconv.ParseInt(leadIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid lead id", http.StatusBadRequest)
		return
	}

	score, err := ah.automationService.CalculateLeadScore(ctx, tenantID, leadID)
	if err != nil {
		ah.logger.Error("Failed to calculate lead score", "error", err)
		http.Error(w, "failed to calculate score", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(score)
}

// RankLeads retrieves ranked leads by score
// GET /api/v1/automation/leads/ranked?limit=20
func (ah *AutomationHandler) RankLeads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		ah.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	scores, err := ah.automationService.RankLeads(ctx, tenantID, limit)
	if err != nil {
		ah.logger.Error("Failed to rank leads", "error", err)
		http.Error(w, "failed to rank leads", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total": len(scores),
		"leads": scores,
	})
}

// RouteLeadToAgent routes a lead to an agent
// POST /api/v1/automation/leads/{id}/route
func (ah *AutomationHandler) RouteLeadToAgent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		ah.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	leadIDStr := r.URL.Query().Get("id")
	if leadIDStr == "" {
		http.Error(w, "lead id required", http.StatusBadRequest)
		return
	}

	leadID, err := strconv.ParseInt(leadIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid lead id", http.StatusBadRequest)
		return
	}

	agentID, err := ah.automationService.RouteLeadToAgent(ctx, tenantID, leadID)
	if err != nil {
		ah.logger.Error("Failed to route lead", "error", err)
		http.Error(w, "failed to route lead", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"lead_id":  leadID,
		"agent_id": agentID,
	})
}

// CreateRoutingRule creates a new routing rule
// POST /api/v1/automation/routing-rules
func (ah *AutomationHandler) CreateRoutingRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		ah.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req CreateRoutingRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	rule := &services.RoutingRule{
		TenantID:    tenantID,
		Name:        req.Name,
		Priority:    req.Priority,
		Conditions:  req.Conditions,
		ActionType:  req.ActionType,
		ActionValue: req.ActionValue,
		Enabled:     req.Enabled,
	}

	err := ah.automationService.CreateRoutingRule(ctx, rule)
	if err != nil {
		ah.logger.Error("Failed to create routing rule", "error", err)
		http.Error(w, "failed to create rule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rule)
}

// ScheduleCampaign schedules a campaign for future execution
// POST /api/v1/automation/schedule-campaign
func (ah *AutomationHandler) ScheduleCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		ah.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req ScheduleCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		http.Error(w, "invalid scheduled_at format (use RFC3339)", http.StatusBadRequest)
		return
	}

	err = ah.automationService.ScheduleCampaign(ctx, req.CampaignID, tenantID, scheduledTime)
	if err != nil {
		ah.logger.Error("Failed to schedule campaign", "error", err)
		http.Error(w, "failed to schedule campaign", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Campaign scheduled successfully",
		"campaign_id":  req.CampaignID,
		"scheduled_at": scheduledTime,
	})
}

// GetLeadScoringMetrics gets lead scoring metrics for the tenant
// GET /api/v1/automation/metrics
func (ah *AutomationHandler) GetLeadScoringMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		ah.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	metrics, err := ah.automationService.GetLeadScoringMetrics(ctx, tenantID)
	if err != nil {
		ah.logger.Error("Failed to get metrics", "error", err)
		http.Error(w, "failed to get metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}
