package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

// IntegrationHandler handles integration HTTP requests
type IntegrationHandler struct {
	Service *services.IntegrationService
	Logger  *log.Logger
}

// NewIntegrationHandler creates a new integration handler
func NewIntegrationHandler(service *services.IntegrationService, logger *log.Logger) *IntegrationHandler {
	return &IntegrationHandler{Service: service, Logger: logger}
}

// getTenantIDFromHeader extracts tenant ID from headers
func (h *IntegrationHandler) getTenantIDFromHeader(r *http.Request) (string, error) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		return "", http.ErrNoCookie
	}
	return tenantID, nil
}

// respondWithIntegrationJSON sends JSON response
func (h *IntegrationHandler) respondWithIntegrationJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// Provider Endpoints

// CreateProvider creates a new provider
func (h *IntegrationHandler) CreateProvider(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	var req models.CreateProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Printf("Invalid request body: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	provider, err := h.Service.CreateProvider(r.Context(), tenantID, &req)
	if err != nil {
		h.Logger.Printf("Error creating provider: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusCreated, provider)
}

// GetProvider gets a provider by ID
func (h *IntegrationHandler) GetProvider(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	providerIDStr := vars["providerID"]
	providerID, err := strconv.ParseInt(providerIDStr, 10, 64)
	if err != nil {
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid provider ID"})
		return
	}

	provider, err := h.Service.GetProvider(r.Context(), tenantID, providerID)
	if err != nil {
		h.Logger.Printf("Error getting provider: %v", err)
		if err.Error() == "provider not found" {
			h.respondWithIntegrationJSON(w, http.StatusNotFound, map[string]string{"error": "Provider not found"})
		} else {
			h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusOK, provider)
}

// ListProviders lists providers
func (h *IntegrationHandler) ListProviders(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	limit, offset := h.getPagination(r)

	providers, err := h.Service.ListProviders(r.Context(), tenantID, limit, offset)
	if err != nil {
		h.Logger.Printf("Error listing providers: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusOK, providers)
}

// UpdateProvider updates a provider
func (h *IntegrationHandler) UpdateProvider(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	providerIDStr := vars["providerID"]
	providerID, err := strconv.ParseInt(providerIDStr, 10, 64)
	if err != nil {
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid provider ID"})
		return
	}

	var req models.UpdateProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Printf("Invalid request body: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	err = h.Service.UpdateProvider(r.Context(), tenantID, providerID, &req)
	if err != nil {
		h.Logger.Printf("Error updating provider: %v", err)
		if err.Error() == "provider not found" {
			h.respondWithIntegrationJSON(w, http.StatusNotFound, map[string]string{"error": "Provider not found"})
		} else {
			h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	provider, err := h.Service.GetProvider(r.Context(), tenantID, providerID)
	if err != nil {
		h.Logger.Printf("Error retrieving updated provider: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusOK, provider)
}

// DeleteProvider deletes a provider
func (h *IntegrationHandler) DeleteProvider(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	providerIDStr := vars["providerID"]
	providerID, err := strconv.ParseInt(providerIDStr, 10, 64)
	if err != nil {
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid provider ID"})
		return
	}

	err = h.Service.DeleteProvider(r.Context(), tenantID, providerID)
	if err != nil {
		h.Logger.Printf("Error deleting provider: %v", err)
		if err.Error() == "provider not found" {
			h.respondWithIntegrationJSON(w, http.StatusNotFound, map[string]string{"error": "Provider not found"})
		} else {
			h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusOK, map[string]string{"message": "Provider deleted successfully"})
}

// Webhook Endpoints

// CreateWebhook creates a new webhook
func (h *IntegrationHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	var req models.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Printf("Invalid request body: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	webhook, err := h.Service.CreateWebhook(r.Context(), tenantID, &req)
	if err != nil {
		h.Logger.Printf("Error creating webhook: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusCreated, webhook)
}

// GetWebhook gets a webhook by ID
func (h *IntegrationHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	webhookIDStr := vars["webhookID"]
	webhookID, err := strconv.ParseInt(webhookIDStr, 10, 64)
	if err != nil {
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid webhook ID"})
		return
	}

	webhook, err := h.Service.GetWebhook(r.Context(), tenantID, webhookID)
	if err != nil {
		h.Logger.Printf("Error getting webhook: %v", err)
		if err.Error() == "webhook not found" {
			h.respondWithIntegrationJSON(w, http.StatusNotFound, map[string]string{"error": "Webhook not found"})
		} else {
			h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusOK, webhook)
}

// ListWebhooks lists webhooks for a provider
func (h *IntegrationHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	providerIDStr := vars["providerID"]
	providerID, err := strconv.ParseInt(providerIDStr, 10, 64)
	if err != nil {
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid provider ID"})
		return
	}

	webhooks, err := h.Service.ListWebhooks(r.Context(), tenantID, providerID)
	if err != nil {
		h.Logger.Printf("Error listing webhooks: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusOK, webhooks)
}

// Sync Job Endpoints

// TriggerSync triggers a sync job
func (h *IntegrationHandler) TriggerSync(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	var req models.TriggerSyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Printf("Invalid request body: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	job, err := h.Service.CreateSyncJob(r.Context(), tenantID, &req)
	if err != nil {
		h.Logger.Printf("Error triggering sync: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusAccepted, job)
}

// GetSyncJob gets a sync job by ID
func (h *IntegrationHandler) GetSyncJob(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	jobIDStr := vars["jobID"]
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid job ID"})
		return
	}

	job, err := h.Service.GetSyncJob(r.Context(), tenantID, jobID)
	if err != nil {
		h.Logger.Printf("Error getting sync job: %v", err)
		if err.Error() == "sync job not found" {
			h.respondWithIntegrationJSON(w, http.StatusNotFound, map[string]string{"error": "Sync job not found"})
		} else {
			h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusOK, job)
}

// ListSyncJobs lists sync jobs for a provider
func (h *IntegrationHandler) ListSyncJobs(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	providerIDStr := vars["providerID"]
	providerID, err := strconv.ParseInt(providerIDStr, 10, 64)
	if err != nil {
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid provider ID"})
		return
	}

	jobs, err := h.Service.ListSyncJobs(r.Context(), tenantID, providerID)
	if err != nil {
		h.Logger.Printf("Error listing sync jobs: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusOK, jobs)
}

// Error Log Endpoints

// ListErrorLogs lists error logs
func (h *IntegrationHandler) ListErrorLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	providerIDStr := vars["providerID"]
	providerID, err := strconv.ParseInt(providerIDStr, 10, 64)
	if err != nil {
		h.respondWithIntegrationJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid provider ID"})
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}

	errors, err := h.Service.ListErrorLogs(r.Context(), tenantID, providerID, limit)
	if err != nil {
		h.Logger.Printf("Error listing error logs: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusOK, errors)
}

// Statistics Endpoint

// GetIntegrationStats gets integration statistics
func (h *IntegrationHandler) GetIntegrationStats(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantIDFromHeader(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	stats, err := h.Service.GetIntegrationStats(r.Context(), tenantID)
	if err != nil {
		h.Logger.Printf("Error getting integration stats: %v", err)
		h.respondWithIntegrationJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.respondWithIntegrationJSON(w, http.StatusOK, stats)
}

// Helper functions

func (h *IntegrationHandler) getPagination(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limit = 20
	} else {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		} else {
			limit = 20
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offset = 0
	} else {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		} else {
			offset = 0
		}
	}

	return limit, offset
}
