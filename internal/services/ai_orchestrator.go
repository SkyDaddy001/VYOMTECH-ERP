package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/pkg/logger"
)

type AIOrchestrator struct {
	providers map[string]models.AIProvider
	db        *sql.DB
	cache     map[string]*models.AIResponse
	cacheMu   sync.RWMutex
	logger    *logger.Logger
}

func NewAIOrchestrator(db *sql.DB, logger *logger.Logger) *AIOrchestrator {
	return &AIOrchestrator{
		providers: make(map[string]models.AIProvider),
		db:        db,
		cache:     make(map[string]*models.AIResponse),
		logger:    logger,
	}
}

func (o *AIOrchestrator) RegisterProvider(name string, provider models.AIProvider) {
	o.providers[name] = provider
	o.logger.Info("AI provider registered", "provider", name)
}

func (o *AIOrchestrator) ProcessQuery(ctx context.Context, req *models.AIRequest) (*models.AIResponse, error) {
	// Generate cache key
	cacheKey := o.generateCacheKey(req)

	// Check cache first
	if cached := o.getCachedResponse(cacheKey); cached != nil {
		o.logger.Info("Returning cached AI response", "tenant_id", req.TenantID)
		return cached, nil
	}

	// Select best provider
	provider := o.selectProvider(req.Priority)
	if provider == nil {
		return nil, fmt.Errorf("no available AI providers")
	}

	// Call the provider
	startTime := time.Now()
	response, err := provider.Call(req)
	processingTime := time.Since(startTime)

	if err != nil {
		o.logger.Error("AI provider call failed", "error", err, "provider", response.Provider, "tenant_id", req.TenantID)
		return nil, fmt.Errorf("AI provider error: %w", err)
	}

	response.ProcessingTime = processingTime

	// Cache the response
	o.cacheResponse(cacheKey, response)

	// Log the request for analytics
	if err := o.logAIRequest(req, response); err != nil {
		o.logger.Warn("Failed to log AI request", "error", err)
	}

	o.logger.Info("AI request processed successfully",
		"tenant_id", req.TenantID,
		"provider", response.Provider,
		"tokens_used", response.TokensUsed,
		"processing_time", processingTime,
		"cost", response.Cost)

	return response, nil
}

func (o *AIOrchestrator) selectProvider(priority string) models.AIProvider {
	availableProviders := o.GetAvailableProviders()

	if len(availableProviders) == 0 {
		return nil
	}

	// For high priority requests, prefer faster/cheaper providers
	if priority == "high" {
		// Sort by cost per token (ascending)
		sort.Slice(availableProviders, func(i, j int) bool {
			return availableProviders[i].GetCostPerToken() < availableProviders[j].GetCostPerToken()
		})
	} else {
		// For normal priority, use round-robin or load balancing
		// For simplicity, just return the first available
	}

	return availableProviders[0]
}

func (o *AIOrchestrator) GetAvailableProviders() []models.AIProvider {
	var available []models.AIProvider
	for _, provider := range o.providers {
		if provider.IsAvailable() {
			available = append(available, provider)
		}
	}
	return available
}

func (o *AIOrchestrator) generateCacheKey(req *models.AIRequest) string {
	// Create a hash of tenant_id, query, and context
	hashInput := fmt.Sprintf("%s:%s:%v", req.TenantID, req.Query, req.Context)
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

func (o *AIOrchestrator) getCachedResponse(key string) *models.AIResponse {
	o.cacheMu.RLock()
	defer o.cacheMu.RUnlock()

	if response, exists := o.cache[key]; exists {
		// Check if cache is still valid (e.g., 1 hour)
		if time.Since(time.Now().Add(-time.Hour)) < time.Hour {
			response.Cached = true
			return response
		}
		// Remove expired cache
		delete(o.cache, key)
	}
	return nil
}

func (o *AIOrchestrator) cacheResponse(key string, response *models.AIResponse) {
	o.cacheMu.Lock()
	defer o.cacheMu.Unlock()

	// Limit cache size (simple implementation)
	if len(o.cache) >= 1000 {
		// Remove oldest entry (simplified)
		for k := range o.cache {
			delete(o.cache, k)
			break
		}
	}

	o.cache[key] = response
}

func (o *AIOrchestrator) logAIRequest(req *models.AIRequest, resp *models.AIResponse) error {
	_, err := o.db.Exec(`
        INSERT INTO ai_request_log (
            tenant_id, query, provider, tokens_used, processing_time_ms,
            cost, priority, created_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`,
		req.TenantID,
		truncateString(req.Query, 1000),
		resp.Provider,
		resp.TokensUsed,
		resp.ProcessingTime.Milliseconds(),
		resp.Cost,
		req.Priority,
	)
	return err
}

func (o *AIOrchestrator) GetProviderStats() map[string]interface{} {
	stats := make(map[string]interface{})

	for name, provider := range o.providers {
		stats[name] = map[string]interface{}{
			"available":      provider.IsAvailable(),
			"cost_per_token": provider.GetCostPerToken(),
		}
	}

	return stats
}

func (o *AIOrchestrator) GetTenantUsage(tenantID string, startDate, endDate time.Time) (map[string]interface{}, error) {
	var totalRequests, totalTokens int
	var totalCost float64

	err := o.db.QueryRow(`
        SELECT
            COUNT(*) as total_requests,
            COALESCE(SUM(tokens_used), 0) as total_tokens,
            COALESCE(SUM(cost), 0) as total_cost
        FROM ai_request_log
        WHERE tenant_id = ? AND created_at BETWEEN ? AND ?`,
		tenantID, startDate, endDate).Scan(&totalRequests, &totalTokens, &totalCost)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_requests": totalRequests,
		"total_tokens":   totalTokens,
		"total_cost":     totalCost,
		"period_start":   startDate,
		"period_end":     endDate,
	}, nil
}

// Helper function to truncate strings for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
