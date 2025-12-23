package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

// AIRecommendationsHandler handles AI recommendations-related HTTP requests
type AIRecommendationsHandler struct {
	aiService *services.AIService
	logger    *logger.Logger
}

// NewAIRecommendationsHandler creates a new AIRecommendationsHandler
func NewAIRecommendationsHandler(aiService *services.AIService, logger *logger.Logger) *AIRecommendationsHandler {
	return &AIRecommendationsHandler{
		aiService: aiService,
		logger:    logger,
	}
}

// generateAIUUID generates a UUID for AI operations
func generateAIUUID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().UnixNano())
}

// Helper functions
func (h *AIRecommendationsHandler) getAITenantID(r *http.Request) string {
	return r.Header.Get("X-Tenant-ID")
}

func (h *AIRecommendationsHandler) getAIUserID(r *http.Request) (int, error) {
	userIDStr := r.Header.Get("X-User-ID")
	return strconv.Atoi(userIDStr)
}

// ============================================================
// AI MODEL HANDLERS
// ============================================================

// CreateAIModel creates a new AI model
func (h *AIRecommendationsHandler) CreateAIModel(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	var req models.CreateAIModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	model := &models.AIModel{
		ID:            generateAIUUID(),
		TenantID:      tenantID,
		ModelName:     req.ModelName,
		ModelType:     req.ModelType,
		ModelVersion:  req.ModelVersion,
		Description:   req.Description,
		AlgorithmName: req.AlgorithmName,
		InputFeatures: req.InputFeatures,
		OutputFormat:  req.OutputFormat,
		Status:        "draft",
		IsProduction:  false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.aiService.CreateAIModel(r.Context(), model); err != nil {
		h.logger.Error("Failed to create AI model", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create AI model")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "AI model created successfully",
		Data:    model,
	})
}

// GetAIModel retrieves an AI model
func (h *AIRecommendationsHandler) GetAIModel(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	modelID := strings.TrimPrefix(r.URL.Path, "/api/v1/ai/models/")

	model, err := h.aiService.GetAIModel(r.Context(), modelID, tenantID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "AI model not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "AI model retrieved successfully",
		Data:    model,
	})
}

// ListAIModels lists all AI models for a tenant
func (h *AIRecommendationsHandler) ListAIModels(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	modelType := r.URL.Query().Get("model_type")
	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	var modelTypePtr *string
	if modelType != "" {
		modelTypePtr = &modelType
	}

	modelsList, total, err := h.aiService.ListAIModels(r.Context(), tenantID, modelTypePtr, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list AI models", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to list AI models")
		return
	}

	respondWithJSON(w, http.StatusOK, PaginatedResponse{
		Items:      modelsList,
		Page:       (offset / limit) + 1,
		Limit:      limit,
		Total:      int64(total),
		TotalPages: int((int64(total) + int64(limit) - 1) / int64(limit)),
	})
}

// ============================================================
// RECOMMENDATION ENGINE HANDLERS
// ============================================================

// CreateRecommendationEngine creates a new recommendation engine
func (h *AIRecommendationsHandler) CreateRecommendationEngine(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	var req models.CreateRecommendationEngineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	engine := &models.RecommendationEngine{
		ID:                           generateAIUUID(),
		TenantID:                     tenantID,
		EngineName:                   req.EngineName,
		Description:                  req.Description,
		RecommendationType:           req.RecommendationType,
		ModelID:                      req.ModelID,
		StrategyType:                 req.StrategyType,
		ScoringAlgorithm:             req.ScoringAlgorithm,
		MinConfidenceThreshold:       req.MinConfidenceThreshold,
		MaxRecommendationsPerUser:    req.MaxRecommendationsPerUser,
		MaxRecommendationsPerSession: 5,
		RecommendationTTLHours:       intPtr(24),
		EnablePersonalization:        true,
		EnableRealTimeUpdates:        true,
		EnableABTesting:              false,
		FeedbackEnabled:              true,
		Status:                       "draft",
		CreatedAt:                    time.Now(),
		UpdatedAt:                    time.Now(),
	}

	if err := h.aiService.CreateRecommendationEngine(r.Context(), engine); err != nil {
		h.logger.Error("Failed to create recommendation engine", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create recommendation engine")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Recommendation engine created successfully",
		Data:    engine,
	})
}

// GetRecommendationEngine retrieves a recommendation engine
func (h *AIRecommendationsHandler) GetRecommendationEngine(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	engineID := strings.TrimPrefix(r.URL.Path, "/api/v1/ai/engines/")

	engine, err := h.aiService.GetRecommendationEngine(r.Context(), engineID, tenantID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Recommendation engine not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Recommendation engine retrieved successfully",
		Data:    engine,
	})
}

// ============================================================
// RECOMMENDATIONS HANDLERS
// ============================================================

// GenerateRecommendations generates recommendations for a user
func (h *AIRecommendationsHandler) GenerateRecommendations(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	_, err := h.getAIUserID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		EngineID string `json:"engine_id"`
		ItemType string `json:"item_type"`
		Limit    int    `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Generate recommendations (simplified - in production would call ML model)
	recommendations := make([]*models.UserRecommendation, 0)

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Recommendations generated successfully",
		Data:    recommendations,
	})
}

// GetUserRecommendations retrieves recommendations for a user
func (h *AIRecommendationsHandler) GetUserRecommendations(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	userID, err := h.getAIUserID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	engineID := r.URL.Query().Get("engine_id")
	if engineID == "" {
		respondWithError(w, http.StatusBadRequest, "Engine ID required")
		return
	}

	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	recommendations, err := h.aiService.GetUserRecommendations(r.Context(), userID, engineID, tenantID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get user recommendations", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get recommendations")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Recommendations retrieved successfully",
		Data:    recommendations,
	})
}

// ============================================================
// FEEDBACK HANDLERS
// ============================================================

// SubmitRecommendationFeedback submits feedback on a recommendation
func (h *AIRecommendationsHandler) SubmitRecommendationFeedback(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	userID, err := h.getAIUserID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req models.SubmitFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	feedback := &models.RecommendationFeedback{
		ID:               generateAIUUID(),
		TenantID:         tenantID,
		UserID:           userID,
		RecommendationID: req.RecommendationID,
		FeedbackType:     "explicit",
		FeedbackValue:    req.FeedbackValue,
		RatingScore:      req.RatingScore,
		DetailedFeedback: req.DetailedFeedback,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := h.aiService.SubmitRecommendationFeedback(r.Context(), feedback); err != nil {
		h.logger.Error("Failed to submit feedback", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to submit feedback")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Feedback submitted successfully",
		Data:    feedback,
	})
}

// GetRecommendationFeedback retrieves feedback for a recommendation
func (h *AIRecommendationsHandler) GetRecommendationFeedback(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	recommendationID := r.URL.Query().Get("recommendation_id")
	if recommendationID == "" {
		respondWithError(w, http.StatusBadRequest, "Recommendation ID required")
		return
	}

	feedback, err := h.aiService.GetRecommendationFeedback(r.Context(), recommendationID, tenantID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Feedback not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Feedback retrieved successfully",
		Data:    feedback,
	})
}

// ============================================================
// PREDICTION HANDLERS
// ============================================================

// MakePrediction creates a prediction
func (h *AIRecommendationsHandler) MakePrediction(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	var req struct {
		ModelID     string          `json:"model_id"`
		InputDataID string          `json:"input_data_id"`
		InputData   json.RawMessage `json:"input_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	prediction := &models.PredictionResult{
		ID:          generateAIUUID(),
		TenantID:    tenantID,
		ModelID:     req.ModelID,
		InputDataID: &req.InputDataID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.aiService.MakePrediction(r.Context(), prediction); err != nil {
		h.logger.Error("Failed to make prediction", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to make prediction")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Prediction created successfully",
		Data:    prediction,
	})
}

// GetPredictionResult retrieves a prediction result
func (h *AIRecommendationsHandler) GetPredictionResult(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	predictionID := strings.TrimPrefix(r.URL.Path, "/api/v1/ai/predictions/")

	prediction, err := h.aiService.GetPredictionResult(r.Context(), predictionID, tenantID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Prediction not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Prediction retrieved successfully",
		Data:    prediction,
	})
}

// ============================================================
// ANOMALY DETECTION HANDLERS
// ============================================================

// DetectAnomalies detects anomalies
func (h *AIRecommendationsHandler) DetectAnomalies(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	var req models.DetectAnomaliesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	anomaly := &models.AnomalyDetection{
		ID:            generateAIUUID(),
		TenantID:      tenantID,
		AnomalyType:   req.AnomalyType,
		SeverityLevel: "medium",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.aiService.DetectAnomaly(r.Context(), anomaly); err != nil {
		h.logger.Error("Failed to detect anomaly", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to detect anomaly")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Anomaly detected successfully",
		Data:    anomaly,
	})
}

// GetAnomalies retrieves anomalies
func (h *AIRecommendationsHandler) GetAnomalies(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	severityLevel := r.URL.Query().Get("severity_level")
	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	var severityPtr *string
	if severityLevel != "" {
		severityPtr = &severityLevel
	}

	anomalies, total, err := h.aiService.GetAnomalies(r.Context(), tenantID, severityPtr, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get anomalies", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get anomalies")
		return
	}

	respondWithJSON(w, http.StatusOK, PaginatedResponse{
		Items:      anomalies,
		Page:       (offset / limit) + 1,
		Limit:      limit,
		Total:      int64(total),
		TotalPages: int((int64(total) + int64(limit) - 1) / int64(limit)),
	})
}

// ============================================================
// AI INSIGHTS HANDLERS
// ============================================================

// GenerateInsights generates AI insights
func (h *AIRecommendationsHandler) GenerateInsights(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	var req models.GenerateInsightsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	insight := &models.AIInsight{
		ID:              generateAIUUID(),
		TenantID:        tenantID,
		InsightType:     req.InsightType,
		InsightCategory: req.InsightCategory,
		Title:           "AI Insight",
		Description:     "Generated insight",
		ConfidenceScore: floatPtr(0.85),
		ImpactScore:     floatPtr(0.75),
		RelevanceScore:  floatPtr(0.80),
		TimePeriod:      req.TimePeriod,
		PriorityLevel:   stringPtr("medium"),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.aiService.GenerateInsight(r.Context(), insight); err != nil {
		h.logger.Error("Failed to generate insight", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate insight")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Insight generated successfully",
		Data:    insight,
	})
}

// GetInsights retrieves AI insights
func (h *AIRecommendationsHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	insightType := r.URL.Query().Get("insight_type")
	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	var insightTypePtr *string
	if insightType != "" {
		insightTypePtr = &insightType
	}

	insights, total, err := h.aiService.GetInsights(r.Context(), tenantID, insightTypePtr, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get insights", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get insights")
		return
	}

	respondWithJSON(w, http.StatusOK, PaginatedResponse{
		Items:      insights,
		Page:       (offset / limit) + 1,
		Limit:      limit,
		Total:      int64(total),
		TotalPages: int((int64(total) + int64(limit) - 1) / int64(limit)),
	})
}

// ============================================================
// MODEL PERFORMANCE HANDLERS
// ============================================================

// GetModelPerformance retrieves model performance metrics
func (h *AIRecommendationsHandler) GetModelPerformance(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	modelID := r.URL.Query().Get("model_id")
	if modelID == "" {
		respondWithError(w, http.StatusBadRequest, "Model ID required")
		return
	}

	performance, err := h.aiService.GetModelPerformance(r.Context(), modelID, tenantID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Performance metrics not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Model performance retrieved successfully",
		Data:    performance,
	})
}

// ============================================================
// STATISTICS HANDLERS
// ============================================================

// GetAIStats retrieves AI system statistics
func (h *AIRecommendationsHandler) GetAIStats(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	stats, err := h.aiService.GetAIStats(r.Context(), tenantID)
	if err != nil {
		h.logger.Error("Failed to get AI stats", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get AI stats")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "AI stats retrieved successfully",
		Data:    stats,
	})
}

// GetRecommendationStats retrieves recommendation statistics
func (h *AIRecommendationsHandler) GetRecommendationStats(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getAITenantID(r)
	if tenantID == "" {
		respondWithError(w, http.StatusBadRequest, "Tenant ID required")
		return
	}

	_, err := h.getAIUserID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	engineID := r.URL.Query().Get("engine_id")
	if engineID == "" {
		respondWithError(w, http.StatusBadRequest, "Engine ID required")
		return
	}

	stats, err := h.aiService.GetRecommendationStats(r.Context(), tenantID, engineID)
	if err != nil {
		h.logger.Error("Failed to get recommendation stats", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get recommendation stats")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Recommendation stats retrieved successfully",
		Data:    stats,
	})
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}
