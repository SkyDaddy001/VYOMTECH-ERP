package models

import (
	"encoding/json"
	"time"
)

// CommunicationChannel represents a configured communication channel
type CommunicationChannel struct {
	ID             string    `json:"id" db:"id"`
	TenantID       string    `json:"tenant_id" db:"tenant_id"`
	ChannelType    string    `json:"channel_type" db:"channel_type"` // TELEGRAM, WHATSAPP, SMS, EMAIL
	ChannelName    string    `json:"channel_name" db:"channel_name"`
	APIProvider    string    `json:"api_provider" db:"api_provider"`
	APIKey         string    `json:"api_key" db:"api_key"`
	APISecret      string    `json:"api_secret" db:"api_secret"`
	APIURL         string    `json:"api_url" db:"api_url"`
	WebhookURL     string    `json:"webhook_url" db:"webhook_url"`
	CallbackURL    string    `json:"callback_url" db:"callback_url"`
	AuthToken      string    `json:"auth_token" db:"auth_token"`
	AccountID      string    `json:"account_id" db:"account_id"`
	SenderID       string    `json:"sender_id" db:"sender_id"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	RetryCount     int       `json:"retry_count" db:"retry_count"`
	TimeoutSeconds int       `json:"timeout_seconds" db:"timeout_seconds"`
	Priority       int       `json:"priority" db:"priority"`
	ConfigJSON     string    `json:"config_json" db:"config_json"`
	Notes          string    `json:"notes" db:"notes"`
	CreatedBy      string    `json:"created_by" db:"created_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// MessageTemplate represents a template for messages
type MessageTemplate struct {
	ID                string          `json:"id" db:"id"`
	TenantID          string          `json:"tenant_id" db:"tenant_id"`
	TemplateName      string          `json:"template_name" db:"template_name"`
	ChannelType       string          `json:"channel_type" db:"channel_type"`
	TemplateBody      string          `json:"template_body" db:"template_body"`
	TemplateVariables json.RawMessage `json:"template_variables" db:"template_variables"`
	Subject           string          `json:"subject" db:"subject"`
	Language          string          `json:"language" db:"language"`
	IsActive          bool            `json:"is_active" db:"is_active"`
	UsageCount        int             `json:"usage_count" db:"usage_count"`
	LastUsedAt        *time.Time      `json:"last_used_at" db:"last_used_at"`
	CreatedBy         string          `json:"created_by" db:"created_by"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
}

// CommunicationSession represents a communication conversation
type CommunicationSession struct {
	ID              string          `json:"id" db:"id"`
	TenantID        string          `json:"tenant_id" db:"tenant_id"`
	InitiatedBy     string          `json:"initiated_by" db:"initiated_by"`
	ChannelType     string          `json:"channel_type" db:"channel_type"`
	ChannelID       string          `json:"channel_id" db:"channel_id"`
	SenderID        string          `json:"sender_id" db:"sender_id"`
	RecipientID     string          `json:"recipient_id" db:"recipient_id"`
	RecipientName   string          `json:"recipient_name" db:"recipient_name"`
	RecipientEmail  string          `json:"recipient_email" db:"recipient_email"`
	ContactID       string          `json:"contact_id" db:"contact_id"`
	LeadID          string          `json:"lead_id" db:"lead_id"`
	AgentID         string          `json:"agent_id" db:"agent_id"`
	AccountID       string          `json:"account_id" db:"account_id"`
	CampaignID      string          `json:"campaign_id" db:"campaign_id"`
	ConversationID  string          `json:"conversation_id" db:"conversation_id"`
	CorrelationID   string          `json:"correlation_id" db:"correlation_id"`
	Status          string          `json:"status" db:"status"`       // INITIATED, SENT, DELIVERED, READ, FAILED, COMPLETED
	Direction       string          `json:"direction" db:"direction"` // INBOUND, OUTBOUND, INTERNAL
	MessageCount    int             `json:"message_count" db:"message_count"`
	LastMessageAt   *time.Time      `json:"last_message_at" db:"last_message_at"`
	FirstMessageAt  *time.Time      `json:"first_message_at" db:"first_message_at"`
	LastResponseAt  *time.Time      `json:"last_response_at" db:"last_response_at"`
	ResponseTimeMin int             `json:"response_time_minutes" db:"response_time_minutes"`
	IsArchived      bool            `json:"is_archived" db:"is_archived"`
	IsStarred       bool            `json:"is_starred" db:"is_starred"`
	Priority        int             `json:"priority" db:"priority"`
	Notes           string          `json:"notes" db:"notes"`
	Metadata        json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time      `json:"deleted_at" db:"deleted_at"`
}

// CommunicationMessage represents individual messages
type CommunicationMessage struct {
	ID                string          `json:"id" db:"id"`
	SessionID         string          `json:"session_id" db:"session_id"`
	TenantID          string          `json:"tenant_id" db:"tenant_id"`
	ChannelType       string          `json:"channel_type" db:"channel_type"`
	FromAddress       string          `json:"from_address" db:"from_address"`
	ToAddress         string          `json:"to_address" db:"to_address"`
	MessageType       string          `json:"message_type" db:"message_type"` // TEXT, IMAGE, VIDEO, AUDIO, FILE, LOCATION, TEMPLATE
	MessageSubject    string          `json:"message_subject" db:"message_subject"`
	MessageBody       string          `json:"message_body" db:"message_body"`
	MessageHTML       string          `json:"message_html" db:"message_html"`
	MediaURL          string          `json:"media_url" db:"media_url"`
	MediaSizeBytes    int             `json:"media_size_bytes" db:"media_size_bytes"`
	MediaType         string          `json:"media_type" db:"media_type"`
	Attachments       json.RawMessage `json:"attachments" db:"attachments"`
	TemplateID        string          `json:"template_id" db:"template_id"`
	TemplateVariables json.RawMessage `json:"template_variables" db:"template_variables"`
	ExternalMessageID string          `json:"external_message_id" db:"external_message_id"`
	Status            string          `json:"status" db:"status"` // QUEUED, SENT, DELIVERED, READ, FAILED, BOUNCED
	ErrorCode         string          `json:"error_code" db:"error_code"`
	ErrorMessage      string          `json:"error_message" db:"error_message"`
	RetryCount        int             `json:"retry_count" db:"retry_count"`
	MaxRetries        int             `json:"max_retries" db:"max_retries"`
	SentAt            *time.Time      `json:"sent_at" db:"sent_at"`
	DeliveredAt       *time.Time      `json:"delivered_at" db:"delivered_at"`
	ReadAt            *time.Time      `json:"read_at" db:"read_at"`
	ReadBy            string          `json:"read_by" db:"read_by"`
	Cost              float64         `json:"cost" db:"cost"`
	CostCurrency      string          `json:"cost_currency" db:"cost_currency"`
	Metadata          json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time      `json:"deleted_at" db:"deleted_at"`
}

// ContactCommunicationPreference represents user communication preferences
type ContactCommunicationPreference struct {
	ID                string     `json:"id" db:"id"`
	TenantID          string     `json:"tenant_id" db:"tenant_id"`
	ContactID         string     `json:"contact_id" db:"contact_id"`
	LeadID            string     `json:"lead_id" db:"lead_id"`
	EmailAddress      string     `json:"email_address" db:"email_address"`
	PhoneNumber       string     `json:"phone_number" db:"phone_number"`
	TelegramID        string     `json:"telegram_id" db:"telegram_id"`
	WhatsappNumber    string     `json:"whatsapp_number" db:"whatsapp_number"`
	PreferredChannel  string     `json:"preferred_channel" db:"preferred_channel"`
	AllowTelegram     bool       `json:"allow_telegram" db:"allow_telegram"`
	AllowWhatsapp     bool       `json:"allow_whatsapp" db:"allow_whatsapp"`
	AllowSMS          bool       `json:"allow_sms" db:"allow_sms"`
	AllowEmail        bool       `json:"allow_email" db:"allow_email"`
	OptInSMS          bool       `json:"opt_in_sms" db:"opt_in_sms"`
	OptInMarketing    bool       `json:"opt_in_marketing" db:"opt_in_marketing"`
	OptOutDate        *time.Time `json:"opt_out_date" db:"opt_out_date"`
	DoNotContact      bool       `json:"do_not_contact" db:"do_not_contact"`
	DoNotContactUntil *time.Time `json:"do_not_contact_until" db:"do_not_contact_until"`
	Notes             string     `json:"notes" db:"notes"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// CommunicationWebhookLog represents webhook events
type CommunicationWebhookLog struct {
	ID               string          `json:"id" db:"id"`
	TenantID         string          `json:"tenant_id" db:"tenant_id"`
	ChannelID        string          `json:"channel_id" db:"channel_id"`
	ChannelType      string          `json:"channel_type" db:"channel_type"`
	WebhookEventType string          `json:"webhook_event_type" db:"webhook_event_type"`
	WebhookPayload   json.RawMessage `json:"webhook_payload" db:"webhook_payload"`
	WebhookSignature string          `json:"webhook_signature" db:"webhook_signature"`
	IsValid          bool            `json:"is_valid" db:"is_valid"`
	ProcessingStatus string          `json:"processing_status" db:"processing_status"`
	ErrorMessage     string          `json:"error_message" db:"error_message"`
	ProcessedAt      *time.Time      `json:"processed_at" db:"processed_at"`
	ReceivedAt       time.Time       `json:"received_at" db:"received_at"`
}

// BulkMessageCampaign represents bulk messaging campaigns
type BulkMessageCampaign struct {
	ID              string          `json:"id" db:"id"`
	TenantID        string          `json:"tenant_id" db:"tenant_id"`
	CampaignName    string          `json:"campaign_name" db:"campaign_name"`
	ChannelType     string          `json:"channel_type" db:"channel_type"`
	CampaignDesc    string          `json:"campaign_description" db:"campaign_description"`
	TemplateID      string          `json:"template_id" db:"template_id"`
	RecipientFilter json.RawMessage `json:"recipient_filter" db:"recipient_filter"`
	TotalRecipients int             `json:"total_recipients" db:"total_recipients"`
	SentCount       int             `json:"sent_count" db:"sent_count"`
	DeliveredCount  int             `json:"delivered_count" db:"delivered_count"`
	FailedCount     int             `json:"failed_count" db:"failed_count"`
	CampaignStatus  string          `json:"campaign_status" db:"campaign_status"` // DRAFT, SCHEDULED, RUNNING, COMPLETED, CANCELLED
	ScheduledAt     *time.Time      `json:"scheduled_at" db:"scheduled_at"`
	StartedAt       *time.Time      `json:"started_at" db:"started_at"`
	CompletedAt     *time.Time      `json:"completed_at" db:"completed_at"`
	EstimatedCost   float64         `json:"estimated_cost" db:"estimated_cost"`
	ActualCost      float64         `json:"actual_cost" db:"actual_cost"`
	CreatedBy       string          `json:"created_by" db:"created_by"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// BulkMessageRecipient represents recipients in bulk campaigns
type BulkMessageRecipient struct {
	ID                string          `json:"id" db:"id"`
	CampaignID        string          `json:"campaign_id" db:"campaign_id"`
	TenantID          string          `json:"tenant_id" db:"tenant_id"`
	RecipientAddress  string          `json:"recipient_address" db:"recipient_address"`
	RecipientName     string          `json:"recipient_name" db:"recipient_name"`
	ContactID         string          `json:"contact_id" db:"contact_id"`
	LeadID            string          `json:"lead_id" db:"lead_id"`
	TemplateVariables json.RawMessage `json:"template_variables" db:"template_variables"`
	SendStatus        string          `json:"send_status" db:"send_status"` // PENDING, SENT, DELIVERED, FAILED, SKIPPED
	ErrorCode         string          `json:"error_code" db:"error_code"`
	ErrorMessage      string          `json:"error_message" db:"error_message"`
	MessageID         string          `json:"message_id" db:"message_id"`
	SentAt            *time.Time      `json:"sent_at" db:"sent_at"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
}

// MessageAutomationRule represents automation rules
type MessageAutomationRule struct {
	ID             string          `json:"id" db:"id"`
	TenantID       string          `json:"tenant_id" db:"tenant_id"`
	RuleName       string          `json:"rule_name" db:"rule_name"`
	RuleType       string          `json:"rule_type" db:"rule_type"` // TRIGGER_ON_EVENT, SCHEDULED, WORKFLOW, DRIP_CAMPAIGN
	TriggerEvent   string          `json:"trigger_event" db:"trigger_event"`
	ChannelType    string          `json:"channel_type" db:"channel_type"`
	TemplateID     string          `json:"template_id" db:"template_id"`
	ConditionJSON  json.RawMessage `json:"condition_json" db:"condition_json"`
	ActionJSON     json.RawMessage `json:"action_json" db:"action_json"`
	IsActive       bool            `json:"is_active" db:"is_active"`
	Priority       int             `json:"priority" db:"priority"`
	ExecutionCount int             `json:"execution_count" db:"execution_count"`
	Notes          string          `json:"notes" db:"notes"`
	CreatedBy      string          `json:"created_by" db:"created_by"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

// CommunicationAnalytics represents analytics data
type CommunicationAnalytics struct {
	ID                     string    `json:"id" db:"id"`
	TenantID               string    `json:"tenant_id" db:"tenant_id"`
	ChannelType            string    `json:"channel_type" db:"channel_type"`
	MetricDate             time.Time `json:"metric_date" db:"metric_date"`
	TotalMessages          int       `json:"total_messages" db:"total_messages"`
	SentMessages           int       `json:"sent_messages" db:"sent_messages"`
	DeliveredMessages      int       `json:"delivered_messages" db:"delivered_messages"`
	FailedMessages         int       `json:"failed_messages" db:"failed_messages"`
	ReadMessages           int       `json:"read_messages" db:"read_messages"`
	TotalCost              float64   `json:"total_cost" db:"total_cost"`
	AvgDeliveryTimeSeconds int       `json:"avg_delivery_time_seconds" db:"avg_delivery_time_seconds"`
	AvgResponseTimeMins    int       `json:"avg_response_time_minutes" db:"avg_response_time_minutes"`
	UniqueRecipients       int       `json:"unique_recipients" db:"unique_recipients"`
	EngagementRate         float64   `json:"engagement_rate" db:"engagement_rate"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
}

// ScheduledMessage represents scheduled messages
type ScheduledMessage struct {
	ID                string     `json:"id" db:"id"`
	TenantID          string     `json:"tenant_id" db:"tenant_id"`
	ChannelType       string     `json:"channel_type" db:"channel_type"`
	FromAddress       string     `json:"from_address" db:"from_address"`
	ToAddress         string     `json:"to_address" db:"to_address"`
	TemplateID        string     `json:"template_id" db:"template_id"`
	MessageBody       string     `json:"message_body" db:"message_body"`
	ScheduledFor      time.Time  `json:"scheduled_for" db:"scheduled_for"`
	ScheduledTimezone string     `json:"scheduled_timezone" db:"scheduled_timezone"`
	RecurrencePattern string     `json:"recurrence_pattern" db:"recurrence_pattern"` // daily, weekly, monthly, once
	RecurrenceEndDate *time.Time `json:"recurrence_end_date" db:"recurrence_end_date"`
	Status            string     `json:"status" db:"status"` // SCHEDULED, SENT, CANCELLED, FAILED
	LastSentAt        *time.Time `json:"last_sent_at" db:"last_sent_at"`
	NextSendAt        *time.Time `json:"next_send_at" db:"next_send_at"`
	Notes             string     `json:"notes" db:"notes"`
	CreatedBy         string     `json:"created_by" db:"created_by"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// CommunicationAttachment represents message attachments
type CommunicationAttachment struct {
	ID              string    `json:"id" db:"id"`
	MessageID       string    `json:"message_id" db:"message_id"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	AttachmentType  string    `json:"attachment_type" db:"attachment_type"` // IMAGE, VIDEO, AUDIO, FILE, DOCUMENT
	FileName        string    `json:"file_name" db:"file_name"`
	FileSizeBytes   int       `json:"file_size_bytes" db:"file_size_bytes"`
	MimeType        string    `json:"mime_type" db:"mime_type"`
	FileURL         string    `json:"file_url" db:"file_url"`
	StoragePath     string    `json:"storage_path" db:"storage_path"`
	StorageLocation string    `json:"storage_location" db:"storage_location"` // local, s3, gcs, azure
	Description     string    `json:"description" db:"description"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// UserCommunicationPermission represents user permissions
type UserCommunicationPermission struct {
	ID                   string    `json:"id" db:"id"`
	TenantID             string    `json:"tenant_id" db:"tenant_id"`
	UserID               string    `json:"user_id" db:"user_id"`
	CanSendTelegram      bool      `json:"can_send_telegram" db:"can_send_telegram"`
	CanSendWhatsapp      bool      `json:"can_send_whatsapp" db:"can_send_whatsapp"`
	CanSendSMS           bool      `json:"can_send_sms" db:"can_send_sms"`
	CanSendEmail         bool      `json:"can_send_email" db:"can_send_email"`
	CanViewConversations bool      `json:"can_view_conversations" db:"can_view_conversations"`
	CanManageTemplates   bool      `json:"can_manage_templates" db:"can_manage_templates"`
	CanManageCampaigns   bool      `json:"can_manage_campaigns" db:"can_manage_campaigns"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// API Request/Response types

// SendMessageRequest represents a request to send a message
type SendMessageRequest struct {
	ChannelType       string                 `json:"channel_type" binding:"required"`
	ToAddress         string                 `json:"to_address" binding:"required"`
	MessageBody       string                 `json:"message_body" binding:"required"`
	MessageType       string                 `json:"message_type"`
	TemplateID        string                 `json:"template_id"`
	TemplateVariables map[string]interface{} `json:"template_variables"`
	Subject           string                 `json:"subject"`
	Attachments       []string               `json:"attachments"`
	ScheduleFor       *time.Time             `json:"schedule_for"`
	ContactID         string                 `json:"contact_id"`
	LeadID            string                 `json:"lead_id"`
	Priority          int                    `json:"priority"`
}

// BulkSendRequest represents a request for bulk messaging
type BulkSendRequest struct {
	CampaignName      string          `json:"campaign_name" binding:"required"`
	ChannelType       string          `json:"channel_type" binding:"required"`
	TemplateID        string          `json:"template_id"`
	Recipients        []BulkRecipient `json:"recipients" binding:"required"`
	ScheduleFor       *time.Time      `json:"schedule_for"`
	RecurrencePattern string          `json:"recurrence_pattern"`
}

// BulkRecipient represents a recipient in bulk send
type BulkRecipient struct {
	Address           string                 `json:"address" binding:"required"`
	Name              string                 `json:"name"`
	ContactID         string                 `json:"contact_id"`
	LeadID            string                 `json:"lead_id"`
	TemplateVariables map[string]interface{} `json:"template_variables"`
}

// CreateChannelRequest represents request to create a channel
type CreateChannelRequest struct {
	ChannelType string `json:"channel_type" binding:"required"`
	ChannelName string `json:"channel_name" binding:"required"`
	APIProvider string `json:"api_provider" binding:"required"`
	APIKey      string `json:"api_key"`
	APISecret   string `json:"api_secret"`
	APIURL      string `json:"api_url"`
	WebhookURL  string `json:"webhook_url"`
	AuthToken   string `json:"auth_token"`
	AccountID   string `json:"account_id"`
	SenderID    string `json:"sender_id"`
	ConfigJSON  string `json:"config_json"`
}

// MessageTemplateRequest represents request to create a template
type MessageTemplateRequest struct {
	TemplateName      string                 `json:"template_name" binding:"required"`
	ChannelType       string                 `json:"channel_type" binding:"required"`
	TemplateBody      string                 `json:"template_body" binding:"required"`
	TemplateVariables map[string]interface{} `json:"template_variables"`
	Subject           string                 `json:"subject"`
	Language          string                 `json:"language"`
}

// UpdateContactPreferenceRequest represents request to update contact preferences
type UpdateContactPreferenceRequest struct {
	EmailAddress      string     `json:"email_address"`
	PhoneNumber       string     `json:"phone_number"`
	TelegramID        string     `json:"telegram_id"`
	WhatsappNumber    string     `json:"whatsapp_number"`
	PreferredChannel  string     `json:"preferred_channel"`
	AllowTelegram     *bool      `json:"allow_telegram"`
	AllowWhatsapp     *bool      `json:"allow_whatsapp"`
	AllowSMS          *bool      `json:"allow_sms"`
	AllowEmail        *bool      `json:"allow_email"`
	OptInSMS          *bool      `json:"opt_in_sms"`
	OptInMarketing    *bool      `json:"opt_in_marketing"`
	DoNotContact      *bool      `json:"do_not_contact"`
	DoNotContactUntil *time.Time `json:"do_not_contact_until"`
}
