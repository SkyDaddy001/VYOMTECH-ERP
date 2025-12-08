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

// CommunicationHandler handles communication-related requests
type CommunicationHandler struct {
	commService *services.CommunicationService
	logger      *logger.Logger
}

// NewCommunicationHandler creates a new CommunicationHandler
func NewCommunicationHandler(commService *services.CommunicationService, logger *logger.Logger) *CommunicationHandler {
	return &CommunicationHandler{
		commService: commService,
		logger:      logger,
	}
}

// RegisterProviderRequest is the request body for registering a provider
type RegisterProviderRequest struct {
	ProviderType string                 `json:"provider_type"` // sms, email, whatsapp, slack
	APIKey       string                 `json:"api_key"`
	APISecret    string                 `json:"api_secret,omitempty"`
	BaseURL      string                 `json:"base_url,omitempty"`
	Config       map[string]interface{} `json:"config,omitempty"`
}

// CreateTemplateRequest is the request body for creating a template
type CreateTemplateRequest struct {
	Name         string   `json:"name"`
	Category     string   `json:"category"` // welcome, reminder, follow_up, promotion
	ProviderType string   `json:"provider_type"`
	Subject      string   `json:"subject,omitempty"`
	Body         string   `json:"body"`
	Variables    []string `json:"variables,omitempty"`
}

// SendMessageRequest is the request body for sending a message
type SendMessageRequest struct {
	Recipient    string `json:"recipient"`
	ProviderType string `json:"provider_type"`
	TemplateID   *int64 `json:"template_id,omitempty"`
	Subject      string `json:"subject,omitempty"`
	Body         string `json:"body,omitempty"`
}

// RegisterProvider registers a new communication provider
// POST /api/v1/communication/providers
func (ch *CommunicationHandler) RegisterProvider(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req RegisterProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cred := &services.ProviderCredential{
		TenantID:     tenantID,
		ProviderType: services.CommunicationProviderType(req.ProviderType),
		APIKey:       req.APIKey,
		APISecret:    req.APISecret,
		BaseURL:      req.BaseURL,
		Config:       req.Config,
		IsActive:     true,
	}

	err := ch.commService.RegisterProvider(ctx, cred)
	if err != nil {
		ch.logger.Error("Failed to register provider", "error", err)
		http.Error(w, "failed to register provider", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cred)
}

// CreateTemplate creates a message template
// POST /api/v1/communication/templates
func (ch *CommunicationHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	template := &services.MessageTemplate{
		TenantID:     tenantID,
		Name:         req.Name,
		Category:     req.Category,
		ProviderType: services.CommunicationProviderType(req.ProviderType),
		Subject:      req.Subject,
		Body:         req.Body,
		Variables:    req.Variables,
		IsActive:     true,
	}

	err := ch.commService.CreateTemplate(ctx, template)
	if err != nil {
		ch.logger.Error("Failed to create template", "error", err)
		http.Error(w, "failed to create template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(template)
}

// SendMessage sends a message
// POST /api/v1/communication/messages
func (ch *CommunicationHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	msg := &services.Message{
		TenantID:     tenantID,
		Recipient:    req.Recipient,
		ProviderType: services.CommunicationProviderType(req.ProviderType),
		TemplateID:   req.TemplateID,
		Subject:      req.Subject,
		Body:         req.Body,
	}

	err := ch.commService.SendMessage(ctx, msg)
	if err != nil {
		ch.logger.Error("Failed to send message", "error", err)
		http.Error(w, "failed to send message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

// GetMessageStatus retrieves message status
// GET /api/v1/communication/messages/{id}
func (ch *CommunicationHandler) GetMessageStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	msgIDStr := r.URL.Query().Get("id")
	if msgIDStr == "" {
		http.Error(w, "message id required", http.StatusBadRequest)
		return
	}

	msgID, err := strconv.ParseInt(msgIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	msg, err := ch.commService.GetMessageStatus(ctx, tenantID, msgID)
	if err != nil {
		ch.logger.Error("Failed to get message status", "error", err)
		http.Error(w, "message not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

// GetMessageStats retrieves communication statistics
// GET /api/v1/communication/stats?start_date=2024-01-01&end_date=2024-01-31
func (ch *CommunicationHandler) GetMessageStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		ch.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = parsed
		}
	}

	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = parsed
		}
	}

	stats, err := ch.commService.GetMessageStats(ctx, tenantID, startDate, endDate)
	if err != nil {
		ch.logger.Error("Failed to get statistics", "error", err)
		http.Error(w, "failed to get statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
