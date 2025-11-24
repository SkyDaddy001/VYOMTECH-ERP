package models

import "time"

// ============================================================================
// PHASE 1: AGENT AVAILABILITY & LEAD SCORING
// ============================================================================

// AgentAvailability represents current agent status and availability
type AgentAvailability struct {
	ID                         int64     `json:"id"`
	UserID                     int64     `json:"user_id"`
	TenantID                   string    `json:"tenant_id"`
	Status                     string    `json:"status"` // available, busy, on_break, offline, in_meeting, away
	BreakReason                *string   `json:"break_reason,omitempty"`
	IsAcceptingLeads           bool      `json:"is_accepting_leads"`
	TotalCallsToday            int       `json:"total_calls_today"`
	CurrentCallDurationSeconds int       `json:"current_call_duration_seconds"`
	LastStatusChange           time.Time `json:"last_status_change"`
	LastActivity               time.Time `json:"last_activity"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

// Agent Status Constants
const (
	AgentStatusAvailable = "available"
	AgentStatusBusy      = "busy"
	AgentStatusOnBreak   = "on_break"
	AgentStatusOffline   = "offline"
	AgentStatusInMeeting = "in_meeting"
	AgentStatusAway      = "away"
)

// LeadScore represents calculated lead scoring
type LeadScore struct {
	ID                    int64     `json:"id"`
	LeadID                int64     `json:"lead_id"`
	TenantID              string    `json:"tenant_id"`
	SourceQualityScore    float64   `json:"source_quality_score"`
	EngagementScore       float64   `json:"engagement_score"`
	ConversionProbability float64   `json:"conversion_probability"`
	UrgencyScore          float64   `json:"urgency_score"`
	OverallScore          float64   `json:"overall_score"`
	ScoreCategory         string    `json:"score_category"` // hot, warm, cold, nurture
	PreviousScore         *float64  `json:"previous_score,omitempty"`
	ScoreChange           *float64  `json:"score_change,omitempty"`
	ReasonText            *string   `json:"reason_text,omitempty"`
	CalculationMethod     string    `json:"calculation_method"`
	LastCalculated        time.Time `json:"last_calculated"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// Lead Score Categories
const (
	ScoreCategoryHot     = "hot"
	ScoreCategoryWarm    = "warm"
	ScoreCategoryCold    = "cold"
	ScoreCategoryNurture = "nurture"
)

// ============================================================================
// PHASE 1: LEAD ACTIVITIES & TASKS
// ============================================================================

// LeadActivity represents an activity on a lead
type LeadActivity struct {
	ID           int64      `json:"id"`
	LeadID       int64      `json:"lead_id"`
	TenantID     string     `json:"tenant_id"`
	ActivityType string     `json:"activity_type"` // call, email, meeting, follow_up, note, status_change
	Description  string     `json:"description"`
	CreatedBy    *int64     `json:"created_by,omitempty"`
	ActivityDate time.Time  `json:"activity_date"`
	DurationMin  *int       `json:"duration_minutes,omitempty"`
	Outcome      *string    `json:"outcome,omitempty"`
	NextAction   *string    `json:"next_action,omitempty"`
	NextFollowUp *time.Time `json:"next_follow_up,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Activity Type Constants
const (
	ActivityTypeCall         = "call"
	ActivityTypeEmail        = "email"
	ActivityTypeMeeting      = "meeting"
	ActivityTypeFollowUp     = "follow_up"
	ActivityTypeNote         = "note"
	ActivityTypeStatusChange = "status_change"
	ActivityTypeAssignment   = "assignment"
)

// Task represents a task assigned to a user
type Task struct {
	ID          int64      `json:"id"`
	AssignedTo  int64      `json:"assigned_to"`
	CreatedBy   int64      `json:"created_by"`
	LeadID      *int64     `json:"lead_id,omitempty"`
	TenantID    string     `json:"tenant_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Priority    string     `json:"priority"` // critical, high, normal, low
	Status      string     `json:"status"`   // pending, in_progress, completed, overdue, cancelled
	DueDate     time.Time  `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Task Status Constants
const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
	TaskStatusOverdue    = "overdue"
	TaskStatusCancelled  = "cancelled"
)

// Task Priority Constants
const (
	TaskPriorityCritical = "critical"
	TaskPriorityHigh     = "high"
	TaskPriorityNormal   = "normal"
	TaskPriorityLow      = "low"
)

// ============================================================================
// PHASE 1: NOTIFICATIONS & COMMUNICATION
// ============================================================================

// Notification represents a user notification
type Notification struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`
	TenantID        string     `json:"tenant_id"`
	Type            string     `json:"type"` // lead_assigned, call_missed, deadline_reminder, task_completed
	Title           string     `json:"title"`
	Message         string     `json:"message"`
	RelatedEntityID *string    `json:"related_entity_id,omitempty"`
	IsRead          bool       `json:"is_read"`
	Priority        string     `json:"priority"` // critical, high, normal, low
	ReadAt          *time.Time `json:"read_at,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// Notification Type Constants
const (
	NotificationTypeLeadAssigned     = "lead_assigned"
	NotificationTypeCallMissed       = "call_missed"
	NotificationTypeDeadlineReminder = "deadline_reminder"
	NotificationTypeTaskCompleted    = "task_completed"
	NotificationTypeLeadScoredHot    = "lead_scored_hot"
)

// CommunicationTemplate represents a reusable message template
type CommunicationTemplate struct {
	ID        int64     `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Channel   string    `json:"channel"` // email, sms, whatsapp, slack
	Content   string    `json:"content"`
	Variables []byte    `json:"variables,omitempty"` // JSON
	CreatedBy int64     `json:"created_by"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommunicationLog represents a sent message
type CommunicationLog struct {
	ID          int64      `json:"id"`
	LeadID      *int64     `json:"lead_id,omitempty"`
	UserID      *int64     `json:"user_id,omitempty"`
	TenantID    string     `json:"tenant_id"`
	Channel     string     `json:"channel"`
	Recipient   string     `json:"recipient"`
	Message     string     `json:"message"`
	TemplateID  *int64     `json:"template_id,omitempty"`
	Status      string     `json:"status"` // sent, failed, delivered, read, bounced
	SentAt      time.Time  `json:"sent_at"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Communication Status Constants
const (
	CommStatusSent      = "sent"
	CommStatusFailed    = "failed"
	CommStatusDelivered = "delivered"
	CommStatusRead      = "read"
	CommStatusBounced   = "bounced"
)

// Communication Channel Constants
const (
	CommChannelEmail    = "email"
	CommChannelSMS      = "sms"
	CommChannelWhatsapp = "whatsapp"
	CommChannelSlack    = "slack"
)

// ============================================================================
// PHASE 1: ANALYTICS & SECURITY
// ============================================================================

// AnalyticsDaily represents daily metrics
type AnalyticsDaily struct {
	ID                int64     `json:"id"`
	Date              time.Time `json:"date"`
	TenantID          string    `json:"tenant_id"`
	TeamID            *int64    `json:"team_id,omitempty"`
	TotalLeadsCreated int       `json:"total_leads_created"`
	TotalCallsMade    int       `json:"total_calls_made"`
	TotalConversions  int       `json:"total_conversions"`
	AvgCallDuration   int       `json:"average_call_duration"`
	ConversionRate    float64   `json:"conversion_rate"`
	AvgLeadScore      float64   `json:"average_lead_score"`
	HotLeadsCount     int       `json:"hot_leads_count"`
	WarmLeadsCount    int       `json:"warm_leads_count"`
	ColdLeadsCount    int       `json:"cold_leads_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TwoFactorCode represents a 2FA code
type TwoFactorCode struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	IsUsed    bool      `json:"is_used"`
	CreatedAt time.Time `json:"created_at"`
}

// TwoFactorSettings represents user's 2FA settings
type TwoFactorSettings struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	IsEnabled  bool       `json:"is_enabled"`
	Method     *string    `json:"method,omitempty"` // sms, email, authenticator
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TwoFA Method Constants
const (
	TwoFAMethodSMS           = "sms"
	TwoFAMethodEmail         = "email"
	TwoFAMethodAuthenticator = "authenticator"
)
