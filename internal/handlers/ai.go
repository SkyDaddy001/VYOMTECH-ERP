package handlers

import (
	"encoding/json"
	"net/http"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/internal/services"
	"multi-tenant-ai-callcenter/pkg/logger"
)

// AIHandler handles AI-related HTTP requests
type AIHandler struct {
	aiOrchestrator *services.AIOrchestrator
	logger         *logger.Logger
}

// NewAIHandler creates a new AIHandler
func NewAIHandler(aiOrchestrator *services.AIOrchestrator, logger *logger.Logger) *AIHandler {
	return &AIHandler{
		aiOrchestrator: aiOrchestrator,
		logger:         logger,
	}
}

// ProcessAIQueryRequest is the request body for processing AI queries
type ProcessAIQueryRequest struct {
	Query       string  `json:"query"`
	Context     string  `json:"context,omitempty"`
	ProviderID  string  `json:"provider_id,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

// AIQueryResponse is the response for AI queries
type AIQueryResponse struct {
	Response string `json:"response"`
	Provider string `json:"provider"`
	Tokens   int    `json:"tokens,omitempty"`
}

// AIProvider represents an available AI provider
type AIProvider struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
}

// ProcessAIQuery processes an AI query
// POST /api/v1/ai/query
func (ah *AIHandler) ProcessAIQuery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		ah.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ah.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req ProcessAIQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Query == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	// Create AI request
	aiReq := &models.AIRequest{
		TenantID:    tenantID,
		Query:       req.Query,
		Priority:    "medium",
		Temperature: req.Temperature,
	}

	// Process query through AI orchestrator
	result, err := ah.aiOrchestrator.ProcessQuery(ctx, aiReq)
	if err != nil {
		ah.logger.Error("Failed to process AI query", "userID", userID, "error", err)
		http.Error(w, "failed to process query", http.StatusInternalServerError)
		return
	}

	response := &AIQueryResponse{
		Response: result.Response,
		Provider: result.Provider,
		Tokens:   result.TokensUsed,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListAIProviders lists available AI providers
// GET /api/v1/ai/providers
func (ah *AIHandler) ListAIProviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := ctx.Value("userID").(int64)
	if !ok {
		ah.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get available providers from AI orchestrator
	providers := ah.aiOrchestrator.GetAvailableProviders()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"providers": providers})
}
