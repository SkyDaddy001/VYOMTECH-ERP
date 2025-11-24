package models

import "time"

type AIRequest struct {
    TenantID    string                 `json:"tenant_id"`
    Query       string                 `json:"query"`
    Context     map[string]interface{} `json:"context,omitempty"`
    Priority    string                 `json:"priority"` // low, medium, high
    MaxTokens   int                    `json:"max_tokens,omitempty"`
    Temperature float64                `json:"temperature,omitempty"`
}

type AIResponse struct {
    Response    string    `json:"response"`
    Provider    string    `json:"provider"`
    TokensUsed  int       `json:"tokens_used"`
    ProcessingTime time.Duration `json:"processing_time"`
    Cost        float64   `json:"cost"`
    Cached      bool      `json:"cached"`
}

type AIProvider interface {
    Call(req *AIRequest) (*AIResponse, error)
    GetCostPerToken() float64
    IsAvailable() bool
}

type ProviderConfig struct {
    APIKey      string  `json:"api_key"`
    BaseURL     string  `json:"base_url,omitempty"`
    Model       string  `json:"model"`
    MaxTokens   int     `json:"max_tokens"`
    Temperature float64 `json:"temperature"`
    CostPerToken float64 `json:"cost_per_token"`
}
