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

// CallHandler handles call-related HTTP requests
type CallHandler struct {
	callService *services.CallService
	logger      *logger.Logger
}

// NewCallHandler creates a new CallHandler
func NewCallHandler(callService *services.CallService, logger *logger.Logger) *CallHandler {
	return &CallHandler{
		callService: callService,
		logger:      logger,
	}
}

// CreateCallRequest is the request body for creating a call
type CreateCallRequest struct {
	LeadID       int64  `json:"lead_id"`
	AgentID      int64  `json:"agent_id"`
	Status       string `json:"status,omitempty"`
	RecordingURL string `json:"recording_url,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

// EndCallRequest is the request body for ending a call
type EndCallRequest struct {
	Outcome  string `json:"outcome"`
	Duration int    `json:"duration"`
	Notes    string `json:"notes,omitempty"`
}

// GetCalls retrieves all calls for the tenant
// GET /api/v1/calls
func (ch *CallHandler) GetCalls(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	filter := &models.CallFilter{
		Status:  r.URL.Query().Get("status"),
		Outcome: r.URL.Query().Get("outcome"),
		Limit:   10,
		Offset:  0,
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

	calls, err := ch.callService.GetCalls(ctx, tenantID, filter)
	if err != nil {
		ch.logger.Error("Failed to get calls", "error", err)
		http.Error(w, "failed to get calls", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calls)
}

// GetCall retrieves a specific call
// GET /api/v1/calls/{id}
func (ch *CallHandler) GetCall(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	callID := r.URL.Query().Get("id")
	if callID == "" {
		http.Error(w, "call id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(callID, 10, 64)
	if err != nil {
		http.Error(w, "invalid call id", http.StatusBadRequest)
		return
	}

	call, err := ch.callService.GetCall(ctx, id, tenantID)
	if err != nil {
		ch.logger.Error("Failed to get call", "callID", id, "error", err)
		http.Error(w, "call not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(call)
}

// CreateCall creates a new call
// POST /api/v1/calls
func (ch *CallHandler) CreateCall(w http.ResponseWriter, r *http.Request) {
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

	var req CreateCallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.LeadID == 0 || req.AgentID == 0 {
		http.Error(w, "lead_id and agent_id are required", http.StatusBadRequest)
		return
	}

	status := req.Status
	if status == "" {
		status = "initiated"
	}

	now := time.Now()
	call := &models.Call{
		TenantID:     tenantID,
		LeadID:       req.LeadID,
		AgentID:      req.AgentID,
		Status:       status,
		RecordingURL: req.RecordingURL,
		Notes:        req.Notes,
		StartedAt:    &now,
	}

	if err := ch.callService.CreateCall(ctx, call); err != nil {
		ch.logger.Error("Failed to create call", "userID", userID, "error", err)
		http.Error(w, "failed to create call", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(call)
}

// EndCall ends an active call
// POST /api/v1/calls/{id}/end
func (ch *CallHandler) EndCall(w http.ResponseWriter, r *http.Request) {
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

	callID := r.URL.Query().Get("id")
	if callID == "" {
		http.Error(w, "call id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(callID, 10, 64)
	if err != nil {
		http.Error(w, "invalid call id", http.StatusBadRequest)
		return
	}

	var req EndCallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Outcome == "" || req.Duration == 0 {
		http.Error(w, "outcome and duration are required", http.StatusBadRequest)
		return
	}

	if err := ch.callService.EndCall(ctx, id, tenantID, req.Outcome, req.Duration); err != nil {
		ch.logger.Error("Failed to end call", "callID", id, "userID", userID, "error", err)
		http.Error(w, "call not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "call ended successfully"})
}

// GetCallStats retrieves call statistics
// GET /api/v1/calls/stats
func (ch *CallHandler) GetCallStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	stats, err := ch.callService.GetCallStats(ctx, tenantID)
	if err != nil {
		ch.logger.Error("Failed to get call stats", "error", err)
		http.Error(w, "failed to get call stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
