package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/internal/services"
	"multi-tenant-ai-callcenter/pkg/logger"

	"github.com/gorilla/mux"
)

// Phase1Handler handles all Phase 1 feature endpoints
type Phase1Handler struct {
	leadScoringService *services.LeadScoringService
	logger             *logger.Logger
}

func NewPhase1Handler(router *mux.Router, leadScoringService *services.LeadScoringService, logger *logger.Logger) *Phase1Handler {
	handler := &Phase1Handler{
		leadScoringService: leadScoringService,
		logger:             logger,
	}

	// Lead Scoring Routes - registered with authentication middleware in router
	router.HandleFunc("/leads/{id}/score", handler.GetLeadScore).Methods("GET")
	router.HandleFunc("/leads/{id}/score/calculate", handler.CalculateLeadScore).Methods("POST")
	router.HandleFunc("/leads/scores/category/{category}", handler.GetLeadsByCategory).Methods("GET")
	router.HandleFunc("/leads/scores/batch-calculate", handler.BatchCalculateScores).Methods("POST")

	return handler
}

// GetLeadScore retrieves the current score for a lead
func (h *Phase1Handler) GetLeadScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leadID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid lead ID"})
		return
	}

	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Tenant ID required"})
		return
	}

	score, err := h.leadScoringService.GetLeadScore(r.Context(), leadID, tenantID)
	if err != nil {
		h.logger.Error("failed to get lead score", map[string]interface{}{
			"lead_id": leadID,
			"error":   err,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve lead score"})
		return
	}

	if score == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Lead score not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(score)
}

// CalculateLeadScore calculates and saves a lead score
func (h *Phase1Handler) CalculateLeadScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leadID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid lead ID"})
		return
	}

	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Tenant ID required"})
		return
	}

	score, err := h.leadScoringService.UpdateLeadScore(r.Context(), leadID, tenantID)
	if err != nil {
		h.logger.Error("failed to calculate lead score", map[string]interface{}{
			"lead_id": leadID,
			"error":   err,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to calculate lead score"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Lead score calculated successfully",
		"score":   score,
	})
}

// GetLeadsByCategory retrieves leads in a specific score category
func (h *Phase1Handler) GetLeadsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	// Validate category
	validCategories := map[string]bool{
		models.ScoreCategoryHot:     true,
		models.ScoreCategoryWarm:    true,
		models.ScoreCategoryCold:    true,
		models.ScoreCategoryNurture: true,
	}

	if !validCategories[category] {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid category. Must be: hot, warm, cold, or nurture"})
		return
	}

	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Tenant ID required"})
		return
	}

	// Get limit from query params
	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	scores, err := h.leadScoringService.GetLeadsByCategory(r.Context(), tenantID, category, limit)
	if err != nil {
		h.logger.Error("failed to get leads by category", map[string]interface{}{
			"category": category,
			"error":    err,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve leads"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"category": category,
		"count":    len(scores),
		"leads":    scores,
	})
}

// BatchCalculateScores calculates scores for all recent leads
func (h *Phase1Handler) BatchCalculateScores(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Tenant ID required"})
		return
	}

	err := h.leadScoringService.BatchCalculateScores(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("batch calculate scores failed", map[string]interface{}{
			"error": err,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Batch calculation failed"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Batch score calculation started",
		"status":  "processing",
	})
}
