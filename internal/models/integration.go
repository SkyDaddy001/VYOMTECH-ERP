package models

import (
	"encoding/json"
	"time"
)

// IntegrationProvider represents a third-party integration provider
type IntegrationProvider struct {
	ID                int64            `json:"id"`
	TenantID          string           `json:"tenant_id"`
	Name              string           `json:"name"`
	Type              string           `json:"type"`
	APIBaseURL        string           `json:"api_base_url"`
	IsActive          bool             `json:"is_active"`
	RateLimit         int              `json:"rate_limit"`
	RetryCount        int              `json:"retry_count"`
	TimeoutSeconds    int              `json:"timeout_seconds"`
	WebhookSecret     *string          `json:"webhook_secret"`
	OAuthClientID     *string          `json:"oauth_client_id"`
	OAuthClientSecret *string          `json:"oauth_client_secret"`
	Metadata          *json.RawMessage `json:"metadata"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         *time.Time       `json:"deleted_at"`
}

// IntegrationCredential represents API credentials for a provider
type IntegrationCredential struct {
	ID                int64      `json:"id"`
	TenantID          string     `json:"tenant_id"`
	ProviderID        int64      `json:"provider_id"`
	UserID            int64      `json:"user_id"`
	APIKey            *string    `json:"api_key"`
	APISecret         *string    `json:"api_secret"`
	AccessToken       *string    `json:"access_token"`
	RefreshToken      *string    `json:"refresh_token"`
	TokenExpiresAt    *time.Time `json:"token_expires_at"`
	IsValid           bool       `json:"is_valid"`
	LastVerifiedAt    *time.Time `json:"last_verified_at"`
	VerificationError *string    `json:"verification_error"`
	EncryptionKeyID   *string    `json:"encryption_key_id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

// IntegrationWebhook represents a webhook configuration
type IntegrationWebhook struct {
	ID               int64            `json:"id"`
	TenantID         string           `json:"tenant_id"`
	ProviderID       int64            `json:"provider_id"`
	EventType        string           `json:"event_type"`
	WebhookURL       string           `json:"webhook_url"`
	IsActive         bool             `json:"is_active"`
	RetryPolicy      string           `json:"retry_policy"`
	MaxRetries       int              `json:"max_retries"`
	TimeoutSeconds   int              `json:"timeout_seconds"`
	Headers          *json.RawMessage `json:"headers"`
	FilterConditions *json.RawMessage `json:"filter_conditions"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	DeletedAt        *time.Time       `json:"deleted_at"`
}

// IntegrationWebhookEvent represents a webhook event log
type IntegrationWebhookEvent struct {
	ID               int64           `json:"id"`
	TenantID         string          `json:"tenant_id"`
	WebhookID        int64           `json:"webhook_id"`
	EventType        string          `json:"event_type"`
	EventData        json.RawMessage `json:"event_data"`
	Status           string          `json:"status"`
	DeliveryAttempts int             `json:"delivery_attempts"`
	ErrorMessage     *string         `json:"error_message"`
	DeliveredAt      *time.Time      `json:"delivered_at"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// IntegrationFieldMapping represents field mapping for data sync
type IntegrationFieldMapping struct {
	ID                   int64            `json:"id"`
	TenantID             string           `json:"tenant_id"`
	ProviderID           int64            `json:"provider_id"`
	SourceEntity         string           `json:"source_entity"`
	SourceField          string           `json:"source_field"`
	TargetEntity         string           `json:"target_entity"`
	TargetField          string           `json:"target_field"`
	TransformationType   string           `json:"transformation_type"`
	TransformationConfig *json.RawMessage `json:"transformation_config"`
	IsBidirectional      bool             `json:"is_bidirectional"`
	SyncEnabled          bool             `json:"sync_enabled"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
	DeletedAt            *time.Time       `json:"deleted_at"`
}

// IntegrationSyncJob represents a sync job
type IntegrationSyncJob struct {
	ID               int64            `json:"id"`
	TenantID         string           `json:"tenant_id"`
	ProviderID       int64            `json:"provider_id"`
	SyncType         string           `json:"sync_type"`
	Status           string           `json:"status"`
	LastSyncAt       *time.Time       `json:"last_sync_at"`
	NextSyncAt       *time.Time       `json:"next_sync_at"`
	RecordsSynced    int              `json:"records_synced"`
	RecordsFailed    int              `json:"records_failed"`
	SyncDurationSecs *int             `json:"sync_duration_seconds"`
	ErrorLog         *string          `json:"error_log"`
	SyncConfig       *json.RawMessage `json:"sync_config"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

// IntegrationRateLimit represents rate limit tracking
type IntegrationRateLimit struct {
	ID               int64      `json:"id"`
	TenantID         string     `json:"tenant_id"`
	ProviderID       int64      `json:"provider_id"`
	UserID           *int64     `json:"user_id"`
	RequestsCount    int        `json:"requests_count"`
	LimitWindowStart time.Time  `json:"limit_window_start"`
	LimitWindowEnd   *time.Time `json:"limit_window_end"`
	IsRateLimited    bool       `json:"is_rate_limited"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// IntegrationErrorLog represents an error log
type IntegrationErrorLog struct {
	ID              int64            `json:"id"`
	TenantID        string           `json:"tenant_id"`
	ProviderID      int64            `json:"provider_id"`
	ErrorCode       *string          `json:"error_code"`
	ErrorMessage    string           `json:"error_message"`
	ErrorDetails    *json.RawMessage `json:"error_details"`
	Endpoint        *string          `json:"endpoint"`
	RequestPayload  *json.RawMessage `json:"request_payload"`
	ResponsePayload *json.RawMessage `json:"response_payload"`
	Severity        string           `json:"severity"`
	Resolved        bool             `json:"resolved"`
	ResolvedAt      *time.Time       `json:"resolved_at"`
	CreatedAt       time.Time        `json:"created_at"`
}

// IntegrationAuditLog represents an audit log entry
type IntegrationAuditLog struct {
	ID           int64            `json:"id"`
	TenantID     string           `json:"tenant_id"`
	ProviderID   int64            `json:"provider_id"`
	UserID       *int64           `json:"user_id"`
	Action       string           `json:"action"`
	ResourceType *string          `json:"resource_type"`
	ResourceID   *int64           `json:"resource_id"`
	OldValues    *json.RawMessage `json:"old_values"`
	NewValues    *json.RawMessage `json:"new_values"`
	IPAddress    *string          `json:"ip_address"`
	UserAgent    *string          `json:"user_agent"`
	CreatedAt    time.Time        `json:"created_at"`
}

// Request/Response DTOs

// CreateProviderRequest represents a request to create a provider
type CreateProviderRequest struct {
	Name              string           `json:"name" validate:"required"`
	Type              string           `json:"type" validate:"required,oneof=CRM ERP ACCOUNTING COMMUNICATION ANALYTICS PAYMENT LOGISTICS"`
	APIBaseURL        string           `json:"api_base_url" validate:"required,url"`
	RateLimit         int              `json:"rate_limit" validate:"min=1"`
	RetryCount        int              `json:"retry_count" validate:"min=0"`
	TimeoutSeconds    int              `json:"timeout_seconds" validate:"min=1"`
	WebhookSecret     *string          `json:"webhook_secret"`
	OAuthClientID     *string          `json:"oauth_client_id"`
	OAuthClientSecret *string          `json:"oauth_client_secret"`
	Metadata          *json.RawMessage `json:"metadata"`
}

// UpdateProviderRequest represents a request to update a provider
type UpdateProviderRequest struct {
	Name           *string          `json:"name"`
	APIBaseURL     *string          `json:"api_base_url"`
	IsActive       *bool            `json:"is_active"`
	RateLimit      *int             `json:"rate_limit"`
	RetryCount     *int             `json:"retry_count"`
	TimeoutSeconds *int             `json:"timeout_seconds"`
	WebhookSecret  *string          `json:"webhook_secret"`
	Metadata       *json.RawMessage `json:"metadata"`
}

// CreateWebhookRequest represents a request to create a webhook
type CreateWebhookRequest struct {
	ProviderID       int64            `json:"provider_id" validate:"required"`
	EventType        string           `json:"event_type" validate:"required"`
	WebhookURL       string           `json:"webhook_url" validate:"required,url"`
	RetryPolicy      string           `json:"retry_policy" validate:"oneof=LINEAR EXPONENTIAL FIBONACCI"`
	MaxRetries       int              `json:"max_retries" validate:"min=0"`
	TimeoutSeconds   int              `json:"timeout_seconds" validate:"min=1"`
	Headers          *json.RawMessage `json:"headers"`
	FilterConditions *json.RawMessage `json:"filter_conditions"`
}

// TriggerSyncRequest represents a request to trigger a sync job
type TriggerSyncRequest struct {
	ProviderID int64            `json:"provider_id" validate:"required"`
	SyncType   string           `json:"sync_type" validate:"required,oneof=FULL INCREMENTAL DELTA"`
	SyncConfig *json.RawMessage `json:"sync_config"`
}

// VerifyCredentialRequest represents a request to verify credentials
type VerifyCredentialRequest struct {
	ProviderID  int64   `json:"provider_id" validate:"required"`
	APIKey      *string `json:"api_key"`
	APISecret   *string `json:"api_secret"`
	AccessToken *string `json:"access_token"`
}
