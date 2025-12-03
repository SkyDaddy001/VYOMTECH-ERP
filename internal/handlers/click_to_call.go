package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

// ClickToCallHandler handles click-to-call HTTP requests
type ClickToCallHandler struct {
	clickToCallService *services.ClickToCallService
	logger             *logger.Logger
}

// NewClickToCallHandler creates a new ClickToCallHandler
func NewClickToCallHandler(clickToCallService *services.ClickToCallService, logger *logger.Logger) *ClickToCallHandler {
	return &ClickToCallHandler{
		clickToCallService: clickToCallService,
		logger:             logger,
	}
}

// InitiateCall initiates a click-to-call
// POST /api/v1/click-to-call/initiate
func (h *ClickToCallHandler) InitiateCall(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		h.logger.Error("Failed to extract user ID")
		http.Error(w, "user id not found", http.StatusBadRequest)
		return
	}

	var req models.CreateClickToCallRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		h.logger.Error("Failed to parse request", "error", err)
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	if req.ToPhone == "" {
		http.Error(w, "to_phone is required", http.StatusBadRequest)
		return
	}

	if req.Direction == "" {
		req.Direction = "OUTBOUND"
	}

	if req.PhoneType == "" {
		req.PhoneType = "EXTERNAL"
	}

	session, err := h.clickToCallService.CreateClickToCallSession(ctx, tenantID, &req)
	if err != nil {
		h.logger.Error("Failed to create session", "error", err)
		http.Error(w, "failed to initiate call", http.StatusInternalServerError)
		return
	}

	session.InitiatedBy = userID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"session": session,
	})
}

// GetCallSession retrieves a call session
// GET /api/v1/click-to-call/sessions/{id}
func (h *ClickToCallHandler) GetCallSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	sessionID := r.URL.Query().Get("id")
	if sessionID == "" {
		http.Error(w, "session id required", http.StatusBadRequest)
		return
	}

	session, err := h.clickToCallService.GetSession(ctx, sessionID, tenantID)
	if err != nil {
		h.logger.Error("Failed to get session", "error", err)
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// ListCallSessions lists call sessions with filters
// GET /api/v1/click-to-call/sessions
func (h *ClickToCallHandler) ListCallSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	filters := make(map[string]interface{})
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}
	if agentID := r.URL.Query().Get("agent_id"); agentID != "" {
		filters["agent_id"] = agentID
	}
	if leadID := r.URL.Query().Get("lead_id"); leadID != "" {
		filters["lead_id"] = leadID
	}

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	sessions, total, err := h.clickToCallService.ListSessions(ctx, tenantID, filters, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list sessions", "error", err)
		http.Error(w, "failed to list sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":    total,
		"limit":    limit,
		"offset":   offset,
		"sessions": sessions,
	})
}

// UpdateCallStatus updates call status
// PATCH /api/v1/click-to-call/sessions/{id}/status
func (h *ClickToCallHandler) UpdateCallStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	sessionID := r.URL.Query().Get("id")
	if sessionID == "" {
		http.Error(w, "session id required", http.StatusBadRequest)
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	status, ok := req["status"]
	if !ok || status == "" {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	if err := h.clickToCallService.UpdateSessionStatus(ctx, sessionID, tenantID, status); err != nil {
		h.logger.Error("Failed to update session status", "error", err)
		http.Error(w, "failed to update status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Status updated successfully",
	})
}

// EndCall ends a call
// POST /api/v1/click-to-call/sessions/{id}/end
func (h *ClickToCallHandler) EndCall(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	sessionID := r.URL.Query().Get("id")
	if sessionID == "" {
		http.Error(w, "session id required", http.StatusBadRequest)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	reason := ""
	if r, ok := req["reason"].(string); ok {
		reason = r
	}

	if err := h.clickToCallService.EndSession(ctx, sessionID, tenantID, reason, "COMPLETED"); err != nil {
		h.logger.Error("Failed to end call", "error", err)
		http.Error(w, "failed to end call", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Call ended successfully",
	})
}

// CreateVoIPProvider creates a new VoIP provider configuration
// POST /api/v1/voip-providers
func (h *ClickToCallHandler) CreateVoIPProvider(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		userID = "system"
	}

	var provider models.VoIPProvider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		h.logger.Error("Failed to parse request", "error", err)
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	if provider.ProviderName == "" || provider.ProviderType == "" {
		http.Error(w, "provider_name and provider_type are required", http.StatusBadRequest)
		return
	}

	provider.CreatedBy = userID
	if err := h.clickToCallService.CreateVoIPProvider(ctx, tenantID, &provider); err != nil {
		h.logger.Error("Failed to create provider", "error", err)
		http.Error(w, "failed to create provider", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(provider)
}

// GetVoIPProvider gets a VoIP provider by ID
// GET /api/v1/voip-providers/{id}
func (h *ClickToCallHandler) GetVoIPProvider(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	providerID := r.URL.Query().Get("id")
	if providerID == "" {
		http.Error(w, "provider id required", http.StatusBadRequest)
		return
	}

	provider, err := h.clickToCallService.GetVoIPProvider(ctx, providerID, tenantID)
	if err != nil {
		h.logger.Error("Failed to get provider", "error", err)
		http.Error(w, "provider not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(provider)
}

// ListVoIPProviders lists all VoIP providers for tenant
// GET /api/v1/voip-providers
func (h *ClickToCallHandler) ListVoIPProviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	providers, err := h.clickToCallService.ListVoIPProviders(ctx, tenantID)
	if err != nil {
		h.logger.Error("Failed to list providers", "error", err)
		http.Error(w, "failed to list providers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":     len(providers),
		"providers": providers,
	})
}

// HandleWebhookEvent handles incoming webhook from VoIP provider
// POST /api/v1/webhooks/voip/{provider-type}
func (h *ClickToCallHandler) HandleWebhookEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := r.URL.Query().Get("tenant")
	if tenantID == "" {
		// Try to extract from header
		tenantID = r.Header.Get("X-Tenant-ID")
	}

	if tenantID == "" {
		h.logger.Warn("No tenant ID in webhook request")
		http.Error(w, "tenant id required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Save webhook log
	webhookLog := &models.CallWebhookLog{
		WebhookPayload:   string(body),
		WebhookSignature: r.Header.Get("X-Signature"),
		IsValid:          true,
	}

	if err := h.clickToCallService.SaveWebhookLog(ctx, tenantID, webhookLog); err != nil {
		h.logger.Error("Failed to save webhook log", "error", err)
	}

	// Parse webhook payload
	var payload models.CallWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		h.logger.Error("Failed to parse webhook payload", "error", err)
		http.Error(w, "invalid webhook format", http.StatusBadRequest)
		return
	}

	// Process webhook event
	if err := h.clickToCallService.ProcessWebhookEvent(ctx, tenantID, &payload); err != nil {
		h.logger.Error("Failed to process webhook event", "error", err)
		http.Error(w, "failed to process webhook", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Webhook processed successfully",
	})
}

// LogAgentActivity logs agent activity
// POST /api/v1/agent-activity
func (h *ClickToCallHandler) LogAgentActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var activity models.AgentActivityLog
	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		h.logger.Error("Failed to parse request", "error", err)
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	activity.TenantID = tenantID
	if activity.ActivityType == "" {
		http.Error(w, "activity_type is required", http.StatusBadRequest)
		return
	}

	if activity.AgentID == "" {
		http.Error(w, "agent_id is required", http.StatusBadRequest)
		return
	}

	if err := h.clickToCallService.LogAgentActivity(ctx, &activity); err != nil {
		h.logger.Error("Failed to log activity", "error", err)
		http.Error(w, "failed to log activity", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Activity logged successfully",
	})
}

// GetCallStats gets call statistics for tenant
// GET /api/v1/click-to-call/stats
func (h *ClickToCallHandler) GetCallStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	filters := make(map[string]interface{})
	sessions, _, err := h.clickToCallService.ListSessions(ctx, tenantID, filters, 1000, 0)
	if err != nil {
		h.logger.Error("Failed to get sessions", "error", err)
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	// Calculate statistics
	stats := map[string]interface{}{
		"total_calls":      len(sessions),
		"completed_calls":  0,
		"failed_calls":     0,
		"total_duration":   0,
		"average_duration": 0,
		"by_status":        make(map[string]int),
	}

	totalDuration := 0
	for _, s := range sessions {
		stats["by_status"].(map[string]int)[s.Status]++
		if s.Status == "COMPLETED" {
			stats["completed_calls"] = stats["completed_calls"].(int) + 1
			totalDuration += s.DurationSeconds
		}
		if s.Status == "FAILED" {
			stats["failed_calls"] = stats["failed_calls"].(int) + 1
		}
	}

	stats["total_duration"] = totalDuration
	if completed := stats["completed_calls"].(int); completed > 0 {
		stats["average_duration"] = totalDuration / completed
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
