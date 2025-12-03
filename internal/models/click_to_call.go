package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// VoIPProvider represents a VoIP service provider configuration
type VoIPProvider struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	ProviderName    string    `json:"provider_name"`
	ProviderType    string    `json:"provider_type"` // ASTERISK, SIP, MCUBE, EXOTEL, TWILIO, VONAGE, CUSTOM
	APIKey          string    `json:"api_key,omitempty"`
	APISecret       string    `json:"api_secret,omitempty"`
	APIURL          string    `json:"api_url"`
	WebhookURL      string    `json:"webhook_url"`
	CallbackURL     string    `json:"callback_url"`
	AuthToken       string    `json:"auth_token,omitempty"`
	PhoneNumber     string    `json:"phone_number"`
	CallerID        string    `json:"caller_id"`
	DialPlanPrefix  string    `json:"dial_plan_prefix"`
	IsActive        bool      `json:"is_active"`
	RetryCount      int       `json:"retry_count"`
	TimeoutSeconds  int       `json:"timeout_seconds"`
	Priority        int       `json:"priority"`
	ConfigJSON      string    `json:"config_json,omitempty"`
	Notes           string    `json:"notes"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ClickToCallSession represents a single click-to-call session/call
type ClickToCallSession struct {
	ID                string    `json:"id"`
	TenantID          string    `json:"tenant_id"`
	InitiatedBy       string    `json:"initiated_by"`
	FromPhone         string    `json:"from_phone"`
	ToPhone           string    `json:"to_phone"`
	PhoneType         string    `json:"phone_type"` // AGENT, LEAD, CUSTOMER, INTERNAL, EXTERNAL
	ContactName       string    `json:"contact_name"`
	ContactEmail      string    `json:"contact_email"`
	ContactID         string    `json:"contact_id"`
	LeadID            string    `json:"lead_id"`
	AgentID           string    `json:"agent_id"`
	AccountID         string    `json:"account_id"`
	CampaignID        string    `json:"campaign_id"`
	ProviderID        string    `json:"provider_id"`
	ProviderType      string    `json:"provider_type"`
	SessionID         string    `json:"session_id"`
	CorrelationID     string    `json:"correlation_id"`
	Status            string    `json:"status"` // INITIATED, CONNECTING, RINGING, CONNECTED, DISCONNECTED, FAILED, COMPLETED
	Direction         string    `json:"direction"`
	CallStartedAt     *time.Time `json:"call_started_at"`
	CallEndedAt       *time.Time `json:"call_ended_at"`
	DurationSeconds   int       `json:"duration_seconds"`
	RingTimeSeconds   int       `json:"ring_time_seconds"`
	AnswerTimeSeconds int       `json:"answer_time_seconds"`
	DisconnectReason  string    `json:"disconnect_reason"`
	ErrorCode         string    `json:"error_code"`
	ErrorMessage      string    `json:"error_message"`
	RecordingURL      string    `json:"recording_url"`
	TranscriptURL     string    `json:"transcript_url"`
	CallQualityScore  *float32  `json:"call_quality_score"`
	IsRecorded        bool      `json:"is_recorded"`
	IsTransferred     bool      `json:"is_transferred"`
	TransferToAgent   string    `json:"transfer_to_agent"`
	Notes             string    `json:"notes"`
	Metadata          string    `json:"metadata"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

// CallRoutingRule represents a call routing rule
type CallRoutingRule struct {
	ID                 string    `json:"id"`
	TenantID           string    `json:"tenant_id"`
	RuleName           string    `json:"rule_name"`
	RuleType           string    `json:"rule_type"` // AGENT_AVAILABILITY, SKILL_BASED, LOAD_BALANCING, TIME_BASED, PRIORITY, FAILOVER
	Priority           int       `json:"priority"`
	IsActive           bool      `json:"is_active"`
	ConditionJSON      string    `json:"condition_json"`
	ActionJSON         string    `json:"action_json"`
	ProviderID         string    `json:"provider_id"`
	FallbackProviderID string    `json:"fallback_provider_id"`
	RetryOnFailure     bool      `json:"retry_on_failure"`
	RetryCount         int       `json:"retry_count"`
	RetryDelaySeconds  int       `json:"retry_delay_seconds"`
	Notes              string    `json:"notes"`
	CreatedBy          string    `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// IVRMenu represents an IVR menu item
type IVRMenu struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	MenuName      string    `json:"menu_name"`
	MenuCode      string    `json:"menu_code"`
	PromptText    string    `json:"prompt_text"`
	PromptFileURL string    `json:"prompt_file_url"`
	TimeoutSeconds int      `json:"timeout_seconds"`
	MaxRetries    int       `json:"max_retries"`
	ParentMenuID  string    `json:"parent_menu_id"`
	IsActive      bool      `json:"is_active"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Options       []IVRMenuOption `json:"options,omitempty"`
}

// IVRMenuOption represents an IVR menu option
type IVRMenuOption struct {
	ID           string    `json:"id"`
	MenuID       string    `json:"menu_id"`
	TenantID     string    `json:"tenant_id"`
	OptionDigit  string    `json:"option_digit"`
	OptionLabel  string    `json:"option_label"`
	NextMenuID   string    `json:"next_menu_id"`
	ActionType   string    `json:"action_type"` // ROUTE_AGENT, ROUTE_DEPARTMENT, PLAY_MESSAGE, COLLECT_DTMF, HANGUP, TRANSFER, VOICEMAIL
	ActionTarget string    `json:"action_target"`
	Priority     int       `json:"priority"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CallRecordingConfig represents call recording configuration
type CallRecordingConfig struct {
	ID                   string    `json:"id"`
	TenantID             string    `json:"tenant_id"`
	RecordAllCalls       bool      `json:"record_all_calls"`
	RecordOnConsentOnly  bool      `json:"record_on_consent_only"`
	RetentionDays        int       `json:"retention_days"`
	StorageLocation      string    `json:"storage_location"` // local, s3, gcs, azure
	StorageBucket        string    `json:"storage_bucket"`
	EncryptionEnabled    bool      `json:"encryption_enabled"`
	EncryptionMethod     string    `json:"encryption_method"`
	QualityFormat        string    `json:"quality_format"` // MP3, WAV, OGG, M4A
	BitrateBkps          int       `json:"bitrate_kbps"`
	SampleRateHz         int       `json:"sample_rate_hz"`
	AutoTranscription    bool      `json:"auto_transcription"`
	TranscriptionProvider string   `json:"transcription_provider"`
	TranscriptionLang    string    `json:"transcription_lang"`
	SentimentAnalysis    bool      `json:"sentiment_analysis"`
	CreatedBy            string    `json:"created_by"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// AgentActivityLog represents agent activity
type AgentActivityLog struct {
	ID               string    `json:"id"`
	TenantID         string    `json:"tenant_id"`
	AgentID          string    `json:"agent_id"`
	ActivityType     string    `json:"activity_type"` // LOGIN, LOGOUT, ON_CALL, ON_BREAK, ON_ADMIN, IDLE, AWAY, READY
	StatusValue      string    `json:"status_value"`
	SessionID        string    `json:"session_id"`
	IsAvailable      bool      `json:"is_available"`
	ActivityTimestamp time.Time `json:"activity_timestamp"`
	DurationSeconds  int       `json:"duration_seconds"`
	Notes            string    `json:"notes"`
}

// PhoneNumberList represents whitelist/blacklist entries
type PhoneNumberList struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	PhoneNumber  string     `json:"phone_number"`
	ListType     string     `json:"list_type"` // WHITELIST, BLACKLIST, VIP, SPAM
	ContactName  string     `json:"contact_name"`
	ContactEmail string     `json:"contact_email"`
	Reason       string     `json:"reason"`
	Priority     int        `json:"priority"`
	IsActive     bool       `json:"is_active"`
	ExpiresAt    *time.Time `json:"expires_at"`
	CreatedBy    string     `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// CallerIDProfile represents caller ID profiles
type CallerIDProfile struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	ProfileName       string     `json:"profile_name"`
	DisplayName       string     `json:"display_name"`
	PhoneNumber       string     `json:"phone_number"`
	PhoneCountryCode  string     `json:"phone_country_code"`
	IsDefault         bool       `json:"is_default"`
	IsActive          bool       `json:"is_active"`
	VerificationStatus string    `json:"verification_status"` // UNVERIFIED, PENDING, VERIFIED, REJECTED
	VerifiedAt        *time.Time `json:"verified_at"`
	VerifiedBy        string     `json:"verified_by"`
	Notes             string     `json:"notes"`
	CreatedBy         string     `json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// CallWebhookLog represents webhook event logs
type CallWebhookLog struct {
	ID                 string    `json:"id"`
	TenantID           string    `json:"tenant_id"`
	ProviderID         string    `json:"provider_id"`
	WebhookEventType   string    `json:"webhook_event_type"`
	WebhookPayload     string    `json:"webhook_payload"`
	WebhookSignature   string    `json:"webhook_signature"`
	IsValid            bool      `json:"is_valid"`
	ProcessingStatus   string    `json:"processing_status"` // PENDING, PROCESSED, FAILED, SKIPPED
	ErrorMessage       string    `json:"error_message"`
	ProcessedAt        *time.Time `json:"processed_at"`
	ReceivedAt         time.Time `json:"received_at"`
}

// CallDTMFInteraction represents DTMF key presses
type CallDTMFInteraction struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id"`
	TenantID      string    `json:"tenant_id"`
	DTMFDigit     string    `json:"dtmf_digit"`
	DTMFSequence  string    `json:"dtmf_sequence"`
	ReceivedAt    time.Time `json:"received_at"`
	ActionTriggered string  `json:"action_triggered"`
	Notes         string    `json:"notes"`
}

// CallQualityMetric represents call quality metrics
type CallQualityMetric struct {
	ID             string          `json:"id"`
	SessionID      string          `json:"session_id"`
	TenantID       string          `json:"tenant_id"`
	MetricName     string          `json:"metric_name"`
	MetricValue    sql.NullFloat64 `json:"metric_value"`
	Unit           string          `json:"unit"`
	MetricTimestamp time.Time      `json:"metric_timestamp"`
}

// CallTransfer represents call transfer details
type CallTransfer struct {
	ID                  string    `json:"id"`
	SessionID           string    `json:"session_id"`
	TenantID            string    `json:"tenant_id"`
	TransferredFromAgent string   `json:"transferred_from_agent"`
	TransferredToAgent  string    `json:"transferred_to_agent"`
	TransferType        string    `json:"transfer_type"` // BLIND, ATTENDED, WARM, COLD
	TransferReason      string    `json:"transfer_reason"`
	TransferAt          time.Time `json:"transfer_at"`
	NewSessionID        string    `json:"new_session_id"`
	TransferSuccess     bool      `json:"transfer_success"`
	Notes               string    `json:"notes"`
}

// CreateClickToCallRequest is the request payload for initiating a click-to-call
type CreateClickToCallRequest struct {
	ToPhone          string            `json:"to_phone" binding:"required"`
	PhoneType        string            `json:"phone_type"` // AGENT, LEAD, CUSTOMER, INTERNAL, EXTERNAL
	LeadID           string            `json:"lead_id,omitempty"`
	AgentID          string            `json:"agent_id,omitempty"`
	AccountID        string            `json:"account_id,omitempty"`
	CampaignID       string            `json:"campaign_id,omitempty"`
	ProviderID       string            `json:"provider_id,omitempty"`
	Direction        string            `json:"direction"` // INBOUND, OUTBOUND
	CallMetadata     map[string]interface{} `json:"metadata,omitempty"`
	ContactName      string            `json:"contact_name,omitempty"`
	ContactEmail     string            `json:"contact_email,omitempty"`
	Notes            string            `json:"notes,omitempty"`
}

// CallWebhookPayload represents incoming webhook from VoIP provider
type CallWebhookPayload struct {
	EventType       string                 `json:"event_type"`
	SessionID       string                 `json:"session_id"`
	CorrelationID   string                 `json:"correlation_id"`
	Timestamp       time.Time              `json:"timestamp"`
	FromPhone       string                 `json:"from_phone"`
	ToPhone         string                 `json:"to_phone"`
	Status          string                 `json:"status"`
	Duration        int                    `json:"duration"`
	RecordingURL    string                 `json:"recording_url"`
	ErrorCode       string                 `json:"error_code"`
	ErrorMessage    string                 `json:"error_message"`
	Metadata        map[string]interface{} `json:"metadata"`
	RawData         json.RawMessage        `json:"raw_data"`
}

// CallRoutingRequest represents internal routing request
type CallRoutingRequest struct {
	FromPhone      string
	ToPhone        string
	AgentID        string
	LeadID         string
	Direction      string
	CallMetadata   map[string]interface{}
}

// CallRoutingResponse represents routing decision
type CallRoutingResponse struct {
	SelectedProviderID string
	FallbackProviderID string
	RouteRuleID        string
	RoutingDecision    string
	Metadata           map[string]interface{}
}
