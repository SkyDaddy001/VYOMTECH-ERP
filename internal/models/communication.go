package models

import "time"

type MessageType string

const (
	MessageTypeEmail    MessageType = "email"
	MessageTypeSMS      MessageType = "sms"
	MessageTypeWhatsApp MessageType = "whatsapp"
	MessageTypeSlack    MessageType = "slack"
)

type Message struct {
	ID          string                 `json:"id" db:"id"`
	TenantID    string                 `json:"tenant_id" db:"tenant_id"`
	Type        MessageType            `json:"type" db:"type"`
	To          string                 `json:"to" db:"to"`
	Subject     string                 `json:"subject,omitempty" db:"subject"`
	Body        string                 `json:"body" db:"body"`
	Attachments []MessageAttachment    `json:"attachments,omitempty" db:"attachments"`
	Status      MessageStatus          `json:"status" db:"status"`
	Provider    string                 `json:"provider,omitempty" db:"provider"`
	ProviderID  string                 `json:"provider_id,omitempty" db:"provider_id"`
	Error       string                 `json:"error,omitempty" db:"error"`
	SentAt      *time.Time             `json:"sent_at,omitempty" db:"sent_at"`
	DeliveredAt *time.Time             `json:"delivered_at,omitempty" db:"delivered_at"`
	ReadAt      *time.Time             `json:"read_at,omitempty" db:"read_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

type MessageAttachment struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
	MimeType string `json:"mime_type"`
}

type MessageStatus string

const (
	MessageStatusPending   MessageStatus = "pending"
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

type CommunicationCampaign struct {
	ID             string               `json:"id" db:"id"`
	TenantID       string               `json:"tenant_id" db:"tenant_id"`
	Name           string               `json:"name" db:"name"`
	Description    string               `json:"description,omitempty" db:"description"`
	Type           MessageType          `json:"type" db:"type"`
	Status         CampaignStatus       `json:"status" db:"status"`
	TargetAudience CampaignAudience     `json:"target_audience" db:"target_audience"`
	Content        CampaignContent      `json:"content" db:"content"`
	Schedule       CampaignSchedule     `json:"schedule,omitempty" db:"schedule"`
	Metrics        CommunicationMetrics `json:"metrics" db:"metrics"`
	Budget         float64              `json:"budget,omitempty" db:"budget"`
	Cost           float64              `json:"cost" db:"cost"`
	CreatedBy      string               `json:"created_by" db:"created_by"`
	CreatedAt      time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at" db:"updated_at"`
}

type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "draft"
	CampaignStatusScheduled CampaignStatus = "scheduled"
	CampaignStatusRunning   CampaignStatus = "running"
	CampaignStatusPaused    CampaignStatus = "paused"
	CampaignStatusCompleted CampaignStatus = "completed"
	CampaignStatusCancelled CampaignStatus = "cancelled"
)

type CampaignAudience struct {
	LeadIDs       []string               `json:"lead_ids,omitempty"`
	Segments      []AudienceSegment      `json:"segments,omitempty"`
	Filters       map[string]interface{} `json:"filters,omitempty"`
	TotalContacts int                    `json:"total_contacts"`
}

type AudienceSegment struct {
	Name  string `json:"name"`
	Query string `json:"query"`
	Count int    `json:"count"`
}

type CampaignContent struct {
	Subject     string                 `json:"subject,omitempty"`
	Body        string                 `json:"body"`
	Attachments []MessageAttachment    `json:"attachments,omitempty"`
	TemplateID  string                 `json:"template_id,omitempty"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
}

type CampaignSchedule struct {
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	TimeZone  string     `json:"time_zone"`
	SendTimes []string   `json:"send_times,omitempty"` // HH:MM format
	Frequency string     `json:"frequency,omitempty"`  // once, daily, weekly
}

type CommunicationMetrics struct {
	TotalSent      int     `json:"total_sent"`
	TotalDelivered int     `json:"total_delivered"`
	TotalOpened    int     `json:"total_opened"`
	TotalClicked   int     `json:"total_clicked"`
	TotalConverted int     `json:"total_converted"`
	TotalFailed    int     `json:"total_failed"`
	OpenRate       float64 `json:"open_rate"`
	ClickRate      float64 `json:"click_rate"`
	ConversionRate float64 `json:"conversion_rate"`
	BounceRate     float64 `json:"bounce_rate"`
}

type CampaignRecipient struct {
	ID         string                 `json:"id" db:"id"`
	CampaignID string                 `json:"campaign_id" db:"campaign_id"`
	LeadID     string                 `json:"lead_id,omitempty" db:"lead_id"`
	Contact    string                 `json:"contact" db:"contact"` // email or phone
	Status     MessageStatus          `json:"status" db:"status"`
	SentAt     *time.Time             `json:"sent_at,omitempty" db:"sent_at"`
	OpenedAt   *time.Time             `json:"opened_at,omitempty" db:"opened_at"`
	ClickedAt  *time.Time             `json:"clicked_at,omitempty" db:"clicked_at"`
	Error      string                 `json:"error,omitempty" db:"error"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
}

type CommunicationProvider interface {
	Name() string
	SendMessage(msg *Message) error
	GetStatus(messageID string) (MessageStatus, error)
	IsAvailable() bool
	GetCost() float64
}

type CommunicationProviderConfig struct {
	Name      string                 `json:"name"`
	Type      MessageType            `json:"type"`
	APIKey    string                 `json:"api_key"`
	APISecret string                 `json:"api_secret,omitempty"`
	BaseURL   string                 `json:"base_url,omitempty"`
	Settings  map[string]interface{} `json:"settings,omitempty"`
	IsActive  bool                   `json:"is_active"`
}
