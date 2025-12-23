package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

type ConstructionHandler struct {
	service     *services.ConstructionService
	RBACService *services.RBACService
}

// NewConstructionHandler creates a new construction handler
func NewConstructionHandler(service *services.ConstructionService, rbacService *services.RBACService) *ConstructionHandler {
	return &ConstructionHandler{
		service:     service,
		RBACService: rbacService,
	}
}

// GetDashboardMetrics - GET /api/v1/construction/dashboard
func (h *ConstructionHandler) GetDashboardMetrics(w http.ResponseWriter, r *http.Request) {
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

// CreateProject - POST /api/v1/construction/projects
func (h *ConstructionHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	// Extract user and tenant from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, `{"error": "Tenant ID not found in context"}`, http.StatusForbidden)
		return
	}

	// Verify permission
	if err := h.RBACService.VerifyPermission(r.Context(), tenantID, userID, "construction.projects.create"); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Permission denied: %s"}`, err.Error()), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Project created"})
}

// ListProjects - GET /api/v1/construction/projects
func (h *ConstructionHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  []interface{}{},
		"total": 0,
		"page":  1,
		"limit": 20,
	})
}

// GetProject - GET /api/v1/construction/projects/{id}
func (h *ConstructionHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Project details"})
}

// UpdateProject - PUT /api/v1/construction/projects/{id}
func (h *ConstructionHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Project updated"})
}

// DeleteProject - DELETE /api/v1/construction/projects/{id}
func (h *ConstructionHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// CreateBOQItem - POST /api/v1/construction/boq
func (h *ConstructionHandler) CreateBOQItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "BOQ item created"})
}

// ListBOQItems - GET /api/v1/construction/projects/{projectId}/boq
func (h *ConstructionHandler) ListBOQItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  []interface{}{},
		"total": 0,
		"page":  1,
		"limit": 20,
	})
}

// LogProgress - POST /api/v1/construction/progress
func (h *ConstructionHandler) LogProgress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Progress logged"})
}

// GetProgressHistory - GET /api/v1/construction/projects/{projectId}/progress
func (h *ConstructionHandler) GetProgressHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  []interface{}{},
		"total": 0,
		"page":  1,
		"limit": 20,
	})
}

// CreateQualityInspection - POST /api/v1/construction/quality
func (h *ConstructionHandler) CreateQualityInspection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Inspection created"})
}

// ListQualityInspections - GET /api/v1/construction/projects/{projectId}/quality
func (h *ConstructionHandler) ListQualityInspections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  []interface{}{},
		"total": 0,
		"page":  1,
		"limit": 20,
	})
}

// RegisterConstructionRoutes registers construction routes
func RegisterConstructionRoutes(router *mux.Router, service *services.ConstructionService, rbacService *services.RBACService) {
	handler := NewConstructionHandler(service, rbacService)

	// Dashboard
	router.HandleFunc("/api/v1/construction/dashboard", handler.GetDashboardMetrics).Methods("GET")

	// Projects
	router.HandleFunc("/api/v1/construction/projects", handler.CreateProject).Methods("POST")
	router.HandleFunc("/api/v1/construction/projects", handler.ListProjects).Methods("GET")
	router.HandleFunc("/api/v1/construction/projects/{id}", handler.GetProject).Methods("GET")
	router.HandleFunc("/api/v1/construction/projects/{id}", handler.UpdateProject).Methods("PUT")
	router.HandleFunc("/api/v1/construction/projects/{id}", handler.DeleteProject).Methods("DELETE")

	// Bill of Quantities
	router.HandleFunc("/api/v1/construction/boq", handler.CreateBOQItem).Methods("POST")
	router.HandleFunc("/api/v1/construction/projects/{projectId}/boq", handler.ListBOQItems).Methods("GET")

	// Progress Tracking
	router.HandleFunc("/api/v1/construction/progress", handler.LogProgress).Methods("POST")
	router.HandleFunc("/api/v1/construction/projects/{projectId}/progress", handler.GetProgressHistory).Methods("GET")

	// Quality Control
	router.HandleFunc("/api/v1/construction/quality", handler.CreateQualityInspection).Methods("POST")
	router.HandleFunc("/api/v1/construction/projects/{projectId}/quality", handler.ListQualityInspections).Methods("GET")
}
