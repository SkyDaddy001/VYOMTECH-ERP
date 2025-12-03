package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
)

// MultiChannelCommunicationHandler handles multi-channel communication endpoints
type MultiChannelCommunicationHandler struct {
	service services.MultiChannelCommunicationService
	logger  services.Logger
}

// NewMultiChannelCommunicationHandler creates a new handler
func NewMultiChannelCommunicationHandler(service services.MultiChannelCommunicationService, logger services.Logger) *MultiChannelCommunicationHandler {
	return &MultiChannelCommunicationHandler{
		service: service,
		logger:  logger,
	}
}

// CreateCommunicationChannel creates a new communication channel
// POST /api/v1/communication/channels
func (h *MultiChannelCommunicationHandler) CreateCommunicationChannel(w http.ResponseWriter, r *http.Request) {
	var req models.CreateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantID := r.Context().Value("tenant_id").(string)

	channel, err := h.service.CreateCommunicationChannel(r.Context(), &req, tenantID)
	if err != nil {
		h.logger.Error("Failed to create channel: " + err.Error())
		http.Error(w, "Failed to create channel", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(channel)
}

// GetCommunicationChannel retrieves a channel by ID
// GET /api/v1/communication/channels/{id}
func (h *MultiChannelCommunicationHandler) GetCommunicationChannel(w http.ResponseWriter, r *http.Request) {
	channelID := r.PathValue("id")
	if channelID == "" {
		http.Error(w, "Channel ID required", http.StatusBadRequest)
		return
	}

	tenantID := r.Context().Value("tenant_id").(string)

	channel, err := h.service.GetCommunicationChannel(r.Context(), tenantID, channelID)
	if err != nil {
		h.logger.Error("Channel not found: " + err.Error())
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channel)
}

// ListCommunicationChannels lists all channels for a tenant
// GET /api/v1/communication/channels
func (h *MultiChannelCommunicationHandler) ListCommunicationChannels(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	tenantID := r.Context().Value("tenant_id").(string)

	channels, err := h.service.ListCommunicationChannels(r.Context(), tenantID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list channels: " + err.Error())
		http.Error(w, "Failed to list channels", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"channels": channels,
		"limit":    limit,
		"offset":   offset,
	})
}

// SendMessage sends a message via specified channel
// POST /api/v1/communication/messages
func (h *MultiChannelCommunicationHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req models.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantID := r.Context().Value("tenant_id").(string)

	message, err := h.service.SendMessage(r.Context(), tenantID, &req)
	if err != nil {
		h.logger.Error("Failed to send message: " + err.Error())
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// SendBulkMessages sends bulk messages
// POST /api/v1/communication/bulk-send
func (h *MultiChannelCommunicationHandler) SendBulkMessages(w http.ResponseWriter, r *http.Request) {
	var req models.BulkSendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantID := r.Context().Value("tenant_id").(string)

	campaign, err := h.service.SendBulkMessages(r.Context(), tenantID, &req)
	if err != nil {
		h.logger.Error("Failed to create bulk campaign: " + err.Error())
		http.Error(w, "Failed to create bulk campaign", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(campaign)
}

// CreateMessageTemplate creates a new message template
// POST /api/v1/communication/templates
func (h *MultiChannelCommunicationHandler) CreateMessageTemplate(w http.ResponseWriter, r *http.Request) {
	var req models.MessageTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantID := r.Context().Value("tenant_id").(string)

	template, err := h.service.CreateMessageTemplate(r.Context(), tenantID, &req)
	if err != nil {
		h.logger.Error("Failed to create template: " + err.Error())
		http.Error(w, "Failed to create template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(template)
}

// GetMessageTemplate retrieves a template
// GET /api/v1/communication/templates/{id}
func (h *MultiChannelCommunicationHandler) GetMessageTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := r.PathValue("id")
	if templateID == "" {
		http.Error(w, "Template ID required", http.StatusBadRequest)
		return
	}

	tenantID := r.Context().Value("tenant_id").(string)

	template, err := h.service.GetMessageTemplate(r.Context(), tenantID, templateID)
	if err != nil {
		h.logger.Error("Template not found: " + err.Error())
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

// ListMessageTemplates lists all templates
// GET /api/v1/communication/templates
func (h *MultiChannelCommunicationHandler) ListMessageTemplates(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	tenantID := r.Context().Value("tenant_id").(string)

	templates, err := h.service.ListMessageTemplates(r.Context(), tenantID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list templates: " + err.Error())
		http.Error(w, "Failed to list templates", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"templates": templates,
		"limit":     limit,
		"offset":    offset,
	})
}

// UpdateContactPreference updates contact preferences
// POST /api/v1/communication/contacts/{id}/preferences
func (h *MultiChannelCommunicationHandler) UpdateContactPreference(w http.ResponseWriter, r *http.Request) {
	contactID := r.PathValue("id")
	if contactID == "" {
		http.Error(w, "Contact ID required", http.StatusBadRequest)
		return
	}

	var req models.UpdateContactPreferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantID := r.Context().Value("tenant_id").(string)

	pref, err := h.service.UpdateContactPreference(r.Context(), tenantID, contactID, &req)
	if err != nil {
		h.logger.Error("Failed to update preferences: " + err.Error())
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pref)
}

// GetContactPreference gets contact preferences
// GET /api/v1/communication/contacts/{id}/preferences
func (h *MultiChannelCommunicationHandler) GetContactPreference(w http.ResponseWriter, r *http.Request) {
	contactID := r.PathValue("id")
	if contactID == "" {
		http.Error(w, "Contact ID required", http.StatusBadRequest)
		return
	}

	tenantID := r.Context().Value("tenant_id").(string)

	pref, err := h.service.GetContactPreference(r.Context(), tenantID, contactID)
	if err != nil {
		h.logger.Error("Contact preferences not found: " + err.Error())
		http.Error(w, "Contact preferences not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pref)
}

// HandleWebhookEvent handles incoming webhooks from providers
// POST /api/v1/communication/webhooks/{provider-type}
func (h *MultiChannelCommunicationHandler) HandleWebhookEvent(w http.ResponseWriter, r *http.Request) {
	providerType := r.PathValue("provider-type")
	if providerType == "" {
		http.Error(w, "Provider type required", http.StatusBadRequest)
		return
	}

	tenantID := r.Context().Value("tenant_id").(string)
	signature := r.Header.Get("X-Webhook-Signature")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ProcessWebhookEvent(r.Context(), tenantID, strings.ToUpper(providerType), signature, body); err != nil {
		h.logger.Error("Failed to process webhook: " + err.Error())
		http.Error(w, "Failed to process webhook", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "processed"})
}

// GetCommunicationSession retrieves a session
// GET /api/v1/communication/sessions/{id}
func (h *MultiChannelCommunicationHandler) GetCommunicationSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	tenantID := r.Context().Value("tenant_id").(string)

	session, err := h.service.GetCommunicationSession(r.Context(), tenantID, sessionID)
	if err != nil {
		h.logger.Error("Session not found: " + err.Error())
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// ListCommunicationSessions lists sessions
// GET /api/v1/communication/sessions
func (h *MultiChannelCommunicationHandler) ListCommunicationSessions(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	tenantID := r.Context().Value("tenant_id").(string)

	sessions, err := h.service.ListCommunicationSessions(r.Context(), tenantID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list sessions: " + err.Error())
		http.Error(w, "Failed to list sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sessions": sessions,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetCommunicationStats retrieves communication statistics
// GET /api/v1/communication/stats
func (h *MultiChannelCommunicationHandler) GetCommunicationStats(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	tenantID := r.Context().Value("tenant_id").(string)

	// Parse dates and get analytics (simple implementation)
	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	} else {
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	} else {
		endDate = time.Now()
	}

	analytics, err := h.service.GetCommunicationAnalytics(r.Context(), tenantID, startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get analytics: " + err.Error())
		http.Error(w, "Failed to get analytics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"analytics":  analytics,
		"start_date": startDate,
		"end_date":   endDate,
	})
}

// Error response helper
type ErrorResponseHelper struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	TraceID string `json:"trace_id"`
}
