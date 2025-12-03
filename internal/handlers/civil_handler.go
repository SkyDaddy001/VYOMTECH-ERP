package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

type CivilHandler struct {
	service *services.CivilService
}

// NewCivilHandler creates a new civil handler
func NewCivilHandler(service *services.CivilService) *CivilHandler {
	return &CivilHandler{
		service: service,
	}
}

// GetDashboardMetrics - GET /api/v1/civil/dashboard
func (h *CivilHandler) GetDashboardMetrics(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, `{"error": "X-Tenant-ID header required"}`, http.StatusBadRequest)
		return
	}

	metrics, err := h.service.GetDashboardMetrics(tenantID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

// CreateSite - POST /api/v1/civil/sites
func (h *CivilHandler) CreateSite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Site created"})
}

// ListSites - GET /api/v1/civil/sites
func (h *CivilHandler) ListSites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  []interface{}{},
		"total": 0,
		"page":  1,
		"limit": 20,
	})
}

// GetSite - GET /api/v1/civil/sites/{id}
func (h *CivilHandler) GetSite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Site details"})
}

// UpdateSite - PUT /api/v1/civil/sites/{id}
func (h *CivilHandler) UpdateSite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Site updated"})
}

// DeleteSite - DELETE /api/v1/civil/sites/{id}
func (h *CivilHandler) DeleteSite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// ReportIncident - POST /api/v1/civil/incidents
func (h *CivilHandler) ReportIncident(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Incident reported"})
}

// ListIncidents - GET /api/v1/civil/sites/{siteId}/incidents
func (h *CivilHandler) ListIncidents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  []interface{}{},
		"total": 0,
		"page":  1,
		"limit": 20,
	})
}

// UpdateIncidentStatus - PUT /api/v1/civil/incidents/{id}/status
func (h *CivilHandler) UpdateIncidentStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// ListCompliance - GET /api/v1/civil/sites/{siteId}/compliance
func (h *CivilHandler) ListCompliance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  []interface{}{},
		"total": 0,
		"page":  1,
		"limit": 20,
	})
}

// UpdateComplianceStatus - PUT /api/v1/civil/compliance/{id}/status
func (h *CivilHandler) UpdateComplianceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// CreatePermit - POST /api/v1/civil/permits
func (h *CivilHandler) CreatePermit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Permit created"})
}

// ListPermits - GET /api/v1/civil/sites/{siteId}/permits
func (h *CivilHandler) ListPermits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  []interface{}{},
		"total": 0,
		"page":  1,
		"limit": 20,
	})
}

// UpdatePermitStatus - PUT /api/v1/civil/permits/{id}/status
func (h *CivilHandler) UpdatePermitStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// RegisterCivilRoutes registers civil engineering routes
func RegisterCivilRoutes(router *mux.Router, service *services.CivilService) {
	handler := NewCivilHandler(service)

	// Dashboard
	router.HandleFunc("/api/v1/civil/dashboard", handler.GetDashboardMetrics).Methods("GET")

	// Sites
	router.HandleFunc("/api/v1/civil/sites", handler.CreateSite).Methods("POST")
	router.HandleFunc("/api/v1/civil/sites", handler.ListSites).Methods("GET")
	router.HandleFunc("/api/v1/civil/sites/{id}", handler.GetSite).Methods("GET")
	router.HandleFunc("/api/v1/civil/sites/{id}", handler.UpdateSite).Methods("PUT")
	router.HandleFunc("/api/v1/civil/sites/{id}", handler.DeleteSite).Methods("DELETE")

	// Incidents
	router.HandleFunc("/api/v1/civil/incidents", handler.ReportIncident).Methods("POST")
	router.HandleFunc("/api/v1/civil/sites/{siteId}/incidents", handler.ListIncidents).Methods("GET")
	router.HandleFunc("/api/v1/civil/incidents/{id}/status", handler.UpdateIncidentStatus).Methods("PUT")

	// Compliance
	router.HandleFunc("/api/v1/civil/sites/{siteId}/compliance", handler.ListCompliance).Methods("GET")
	router.HandleFunc("/api/v1/civil/compliance/{id}/status", handler.UpdateComplianceStatus).Methods("PUT")

	// Permits
	router.HandleFunc("/api/v1/civil/permits", handler.CreatePermit).Methods("POST")
	router.HandleFunc("/api/v1/civil/sites/{siteId}/permits", handler.ListPermits).Methods("GET")
	router.HandleFunc("/api/v1/civil/permits/{id}/status", handler.UpdatePermitStatus).Methods("PUT")
}
