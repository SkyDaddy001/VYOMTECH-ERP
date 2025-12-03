package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

type DashboardHandler struct {
	service *services.DashboardService
	logger  *logger.Logger
}

func NewDashboardHandler(router *mux.Router, service *services.DashboardService, logger *logger.Logger) *DashboardHandler {
	h := &DashboardHandler{
		service: service,
		logger:  logger,
	}

	// Register routes
	dashboardRoutes := router.PathPrefix("/dashboard").Subrouter()
	dashboardRoutes.HandleFunc("/analytics", h.GetAnalytics).Methods(http.MethodGet)
	dashboardRoutes.HandleFunc("/activity-logs", h.GetActivityLogs).Methods(http.MethodGet)
	dashboardRoutes.HandleFunc("/users", h.GetUsers).Methods(http.MethodGet)
	dashboardRoutes.HandleFunc("/usage", h.GetUsage).Methods(http.MethodGet)

	return h
}

// GetAnalytics returns tenant analytics
func (h *DashboardHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, `{"error":"Tenant ID required"}`, http.StatusBadRequest)
		return
	}

	analytics, err := h.service.GetTenantAnalytics(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("error getting analytics", map[string]interface{}{"error": err})
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"Failed to get analytics"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// GetActivityLogs returns activity logs for tenant
func (h *DashboardHandler) GetActivityLogs(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, `{"error":"Tenant ID required"}`, http.StatusBadRequest)
		return
	}

	// Parse query parameters
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	logs, err := h.service.GetActivityLogs(r.Context(), tenantID, limit, offset)
	if err != nil {
		h.logger.Error("error getting activity logs", map[string]interface{}{"error": err})
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"Failed to get activity logs"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"logs":   logs,
		"count":  len(logs),
		"limit":  limit,
		"offset": offset,
	}
	json.NewEncoder(w).Encode(response)
}

// GetUsers returns tenant users with stats
func (h *DashboardHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, `{"error":"Tenant ID required"}`, http.StatusBadRequest)
		return
	}

	users, err := h.service.GetTenantUsers(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("error getting users", map[string]interface{}{"error": err})
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"Failed to get users"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"users": users,
		"count": len(users),
	}
	json.NewEncoder(w).Encode(response)
}

// GetUsage returns usage metrics for today
func (h *DashboardHandler) GetUsage(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, `{"error":"Tenant ID required"}`, http.StatusBadRequest)
		return
	}

	metrics, err := h.service.GetUsageMetrics(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("error getting usage metrics", map[string]interface{}{"error": err})
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"Failed to get usage metrics"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
