package models

import (
	"time"

	"gorm.io/datatypes"
)

// ==================== TEAM CHAT MODELS ====================

// TeamChatChannel represents a team chat channel or group
type TeamChatChannel struct {
	ID          string     `gorm:"primaryKey" json:"id"`
	TenantID    string     `json:"tenant_id"`
	ChannelName string     `json:"channel_name"`
	ChannelType string     `json:"channel_type"` // DIRECT, GROUP, ANNOUNCEMENT, DEPARTMENT, PROJECT
	Description string     `json:"description"`
	AvatarURL   string     `json:"avatar_url"`
	IsArchived  int        `json:"is_archived"`
	IsPrivate   int        `json:"is_private"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
}

// TeamChatMember represents a member in a chat channel
type TeamChatMember struct {
	ID                string     `gorm:"primaryKey" json:"id"`
	ChannelID         string     `json:"channel_id"`
	TenantID          string     `json:"tenant_id"`
	UserID            string     `json:"user_id"`
	Role              string     `json:"role"` // OWNER, MODERATOR, MEMBER
	IsMuted           int        `json:"is_muted"`
	LastReadMessageID string     `json:"last_read_message_id"`
	LastReadAt        *time.Time `json:"last_read_at"`
	JoinedAt          time.Time  `json:"joined_at"`
	LeftAt            *time.Time `json:"left_at"`
}

// TeamChatMessage represents a message in a chat channel
type TeamChatMessage struct {
	ID            string                     `gorm:"primaryKey" json:"id"`
	ChannelID     string                     `json:"channel_id"`
	TenantID      string                     `json:"tenant_id"`
	SenderID      string                     `json:"sender_id"`
	MessageType   string                     `json:"message_type"` // TEXT, IMAGE, FILE, VIDEO, LINK, MENTION, SYSTEM
	MessageBody   string                     `json:"message_body"`
	MessageHTML   string                     `json:"message_html"`
	FileURL       string                     `json:"file_url"`
	FileName      string                     `json:"file_name"`
	FileSizeBytes int                        `json:"file_size_bytes"`
	FileMimeType  string                     `json:"file_mime_type"`
	Mentions      datatypes.JSONType[string] `json:"mentions"`
	Reactions     datatypes.JSONType[string] `json:"reactions"`
	IsEdited      int                        `json:"is_edited"`
	IsPinned      int                        `json:"is_pinned"`
	EditedAt      *time.Time                 `json:"edited_at"`
	RepliedToID   string                     `json:"replied_to_message_id"`
	CreatedAt     time.Time                  `json:"created_at"`
	DeletedAt     *time.Time                 `gorm:"index" json:"deleted_at"`
}

// TeamChatReaction represents a reaction to a message
type TeamChatReaction struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	MessageID    string    `json:"message_id"`
	UserID       string    `json:"user_id"`
	ReactionType string    `json:"reaction_type"`
	CreatedAt    time.Time `json:"created_at"`
}

// ==================== WEBRTC VOICE/VIDEO CALL MODELS ====================

// VoiceVideoCall represents a WebRTC voice or video call
type VoiceVideoCall struct {
	ID              string                       `gorm:"primaryKey" json:"id"`
	TenantID        string                       `json:"tenant_id"`
	CallType        string                       `json:"call_type"` // ONE_TO_ONE, GROUP, CONFERENCE
	InitiatorID     string                       `json:"initiator_id"`
	CallStatus      string                       `json:"call_status"`    // RINGING, ACCEPTED, IN_PROGRESS, MISSED, REJECTED, ENDED
	CallDirection   string                       `json:"call_direction"` // INBOUND, OUTBOUND
	IsAudioEnabled  int                          `json:"is_audio_enabled"`
	IsVideoEnabled  int                          `json:"is_video_enabled"`
	IsScreenShared  int                          `json:"is_screen_shared"`
	IsRecording     int                          `json:"is_recording"`
	RecordingURL    string                       `json:"recording_url"`
	CallDuration    int                          `json:"call_duration_seconds"`
	StartedAt       *time.Time                   `json:"started_at"`
	EndedAt         *time.Time                   `json:"ended_at"`
	MissedReason    string                       `json:"missed_reason"`
	WebRTCRoomID    string                       `json:"webrtc_room_id"`
	SignalingServer string                       `json:"signaling_server"`
	STUNServers     datatypes.JSONType[[]string] `json:"stun_servers"`
	TURNServers     datatypes.JSONType[[]string] `json:"turn_servers"`
	ICECandidates   datatypes.JSONType[[]string] `json:"ice_candidates"`
	Metadata        datatypes.JSONType[string]   `json:"metadata"`
	CreatedAt       time.Time                    `json:"created_at"`
	UpdatedAt       time.Time                    `json:"updated_at"`
}

// VoiceVideoCallParticipant represents a participant in a call
type VoiceVideoCallParticipant struct {
	ID                string     `gorm:"primaryKey" json:"id"`
	CallID            string     `json:"call_id"`
	TenantID          string     `json:"tenant_id"`
	UserID            string     `json:"user_id"`
	ParticipantStatus string     `json:"participant_status"` // INVITED, RINGING, JOINED, ON_HOLD, LEFT
	IsAudioMuted      int        `json:"is_audio_muted"`
	IsVideoOff        int        `json:"is_video_off"`
	JoinedAt          *time.Time `json:"joined_at"`
	LeftAt            *time.Time `json:"left_at"`
	DurationSeconds   int        `json:"duration_seconds"`
	AudioQualityScore int        `json:"audio_quality_score"`
	VideoQualityScore int        `json:"video_quality_score"`
	PacketLossPercent float64    `json:"packet_loss_percent"`
	LatencyMs         int        `json:"latency_ms"`
	CreatedAt         time.Time  `json:"created_at"`
}

// ==================== MEETING ROOM MODELS ====================

// MeetingRoom represents a virtual meeting room
type MeetingRoom struct {
	ID                  string                     `gorm:"primaryKey" json:"id"`
	TenantID            string                     `json:"tenant_id"`
	RoomName            string                     `json:"room_name"`
	RoomCode            string                     `json:"room_code"`
	Description         string                     `json:"description"`
	RoomAvatarURL       string                     `json:"room_avatar_url"`
	MaxParticipants     int                        `json:"max_participants"`
	CurrentParticipants int                        `json:"current_participants"`
	RoomStatus          string                     `json:"room_status"` // AVAILABLE, IN_USE, MAINTENANCE, ARCHIVED
	IsPasswordProtected int                        `json:"is_password_protected"`
	PasswordHash        string                     `json:"password_hash"`
	RoomType            string                     `json:"room_type"` // PERMANENT, TEMPORARY, RECURRING
	OwnerID             string                     `json:"owner_id"`
	AllowRecording      int                        `json:"allow_recording"`
	AllowScreenShare    int                        `json:"allow_screen_share"`
	AllowChat           int                        `json:"allow_chat"`
	WebRTCConfig        datatypes.JSONType[string] `json:"webrtc_config"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
	ArchivedAt          *time.Time                 `json:"archived_at"`
}

// MeetingRoomAccess represents access and permissions for a meeting room
type MeetingRoomAccess struct {
	ID                    string    `gorm:"primaryKey" json:"id"`
	RoomID                string    `json:"room_id"`
	TenantID              string    `json:"tenant_id"`
	UserID                string    `json:"user_id"`
	UserRole              string    `json:"user_role"`
	AccessType            string    `json:"access_type"` // OWNER, MODERATOR, PRESENTER, PARTICIPANT, VIEWER
	CanMuteOthers         int       `json:"can_mute_others"`
	CanRemoveParticipants int       `json:"can_remove_participants"`
	CanRecord             int       `json:"can_record"`
	CanShareScreen        int       `json:"can_share_screen"`
	IsActive              int       `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
}

// ==================== CALENDAR MODELS ====================

// CalendarEvent represents a calendar event or appointment
type CalendarEvent struct {
	ID                string                     `gorm:"primaryKey" json:"id"`
	TenantID          string                     `json:"tenant_id"`
	EventTitle        string                     `json:"event_title"`
	EventDescription  string                     `json:"event_description"`
	EventType         string                     `json:"event_type"` // MEETING, CALL, TASK, REMINDER, APPOINTMENT, CONFERENCE
	CreatorID         string                     `json:"creator_id"`
	AssignedTo        string                     `json:"assigned_to"`
	LinkedRoomID      string                     `json:"linked_room_id"`
	LinkedCallID      string                     `json:"linked_call_id"`
	Location          string                     `json:"location"`
	StartTime         time.Time                  `json:"start_time"`
	EndTime           time.Time                  `json:"end_time"`
	DurationMinutes   int                        `json:"duration_minutes"`
	Timezone          string                     `json:"timezone"`
	IsAllDay          int                        `json:"is_all_day"`
	ReminderMinutes   int                        `json:"reminder_minutes"`
	Status            string                     `json:"status"` // SCHEDULED, IN_PROGRESS, COMPLETED, CANCELLED, RESCHEDULED
	IsRecurring       int                        `json:"is_recurring"`
	RecurrencePattern string                     `json:"recurrence_pattern"`
	RecurrenceEndDate *time.Time                 `json:"recurrence_end_date"`
	IsBusy            int                        `json:"is_busy"`
	IsPrivate         int                        `json:"is_private"`
	CalendarID        string                     `json:"calendar_id"`
	ColorCode         string                     `json:"color_code"`
	Attachments       datatypes.JSONType[string] `json:"attachments"`
	Metadata          datatypes.JSONType[string] `json:"metadata"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
	DeletedAt         *time.Time                 `gorm:"index" json:"deleted_at"`
}

// CalendarAttendee represents an attendee for a calendar event
type CalendarAttendee struct {
	ID               string     `gorm:"primaryKey" json:"id"`
	EventID          string     `json:"event_id"`
	TenantID         string     `json:"tenant_id"`
	UserID           string     `json:"user_id"`
	AttendanceStatus string     `json:"attendance_status"` // INVITED, ACCEPTED, DECLINED, TENTATIVE, NO_RESPONSE
	ReminderSent     int        `json:"reminder_sent"`
	IsOrganizer      int        `json:"is_organizer"`
	RespondedAt      *time.Time `json:"responded_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

// ==================== AUTO-DIALER MODELS ====================

// DialerCampaign represents an auto-dialing campaign
type DialerCampaign struct {
	ID                     string     `gorm:"primaryKey" json:"id"`
	TenantID               string     `json:"tenant_id"`
	CampaignName           string     `json:"campaign_name"`
	CampaignType           string     `json:"campaign_type"`   // OUTBOUND, PREVIEW, PREDICTIVE, PROGRESSIVE
	CampaignStatus         string     `json:"campaign_status"` // DRAFT, SCHEDULED, ACTIVE, PAUSED, COMPLETED, CANCELLED
	Description            string     `json:"description"`
	ScriptID               string     `json:"script_id"`
	CallListID             string     `json:"call_list_id"`
	DialStrategy           string     `json:"dial_strategy"` // SEQUENTIAL, RANDOM, PRIORITY_BASED, SKILL_BASED
	MaxConcurrentCalls     int        `json:"max_concurrent_calls"`
	MaxRetries             int        `json:"max_retries"`
	RetryIntervalMinutes   int        `json:"retry_interval_minutes"`
	AbandonedCallThreshold float64    `json:"abandoned_call_threshold_percent"`
	DoNotCallListID        string     `json:"do_not_call_list_id"`
	CallerIDNumber         string     `json:"caller_id_number"`
	VoicemailDetection     int        `json:"voicemail_detection"`
	RecordingEnabled       int        `json:"recording_enabled"`
	AMDEnabled             int        `json:"amd_enabled"`
	ScheduledStartTime     *time.Time `json:"scheduled_start_time"`
	ScheduledEndTime       *time.Time `json:"scheduled_end_time"`
	ActualStartTime        *time.Time `json:"actual_start_time"`
	ActualEndTime          *time.Time `json:"actual_end_time"`
	TotalContacts          int        `json:"total_contacts"`
	ContactedCount         int        `json:"contacted_count"`
	ConnectedCount         int        `json:"connected_count"`
	FailedCount            int        `json:"failed_count"`
	AbandonedCount         int        `json:"abandoned_count"`
	CreatedBy              string     `json:"created_by"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

// CallPriorityQueue represents a call in the priority queue
type CallPriorityQueue struct {
	ID                    string                     `gorm:"primaryKey" json:"id"`
	TenantID              string                     `json:"tenant_id"`
	CampaignID            string                     `json:"campaign_id"`
	ContactPhoneNumber    string                     `json:"contact_phone_number"`
	ContactName           string                     `json:"contact_name"`
	ContactID             string                     `json:"contact_id"`
	LeadID                string                     `json:"lead_id"`
	PriorityLevel         int                        `json:"priority_level"`
	PriorityReason        string                     `json:"priority_reason"`
	QueueStatus           string                     `json:"queue_status"` // PENDING, ASSIGNED, CALLING, COMPLETED, FAILED, RESCHEDULED
	AssignedAgentID       string                     `json:"assigned_agent_id"`
	AssignedAt            *time.Time                 `json:"assigned_at"`
	CallAttemptCount      int                        `json:"call_attempt_count"`
	LastCallTime          *time.Time                 `json:"last_call_time"`
	NextCallTime          *time.Time                 `json:"next_call_time"`
	CallResult            string                     `json:"call_result"` // CONNECTED, VOICEMAIL, DISCONNECTED, NO_ANSWER, BUSY, INVALID, DO_NOT_CALL
	CallNotes             string                     `json:"call_notes"`
	CallDurationSeconds   int                        `json:"call_duration_seconds"`
	IsCallback            int                        `json:"is_callback"`
	CallbackRequestedTime *time.Time                 `json:"callback_requested_time"`
	Metadata              datatypes.JSONType[string] `json:"metadata"`
	CreatedAt             time.Time                  `json:"created_at"`
	UpdatedAt             time.Time                  `json:"updated_at"`
}

// DialerScript represents a dialing script or call flow
type DialerScript struct {
	ID                    string                     `gorm:"primaryKey" json:"id"`
	TenantID              string                     `json:"tenant_id"`
	ScriptName            string                     `json:"script_name"`
	ScriptDescription     string                     `json:"script_description"`
	ScriptType            string                     `json:"script_type"` // GREETING, QUALIFICATION, OBJECTION_HANDLING, CLOSING, FOLLOW_UP
	ScriptContent         datatypes.JSONType[string] `json:"script_content"`
	VoiceGuidanceAudioURL string                     `json:"voice_guidance_audio_url"`
	CreatedBy             string                     `json:"created_by"`
	IsActive              int                        `json:"is_active"`
	Version               int                        `json:"version"`
	CreatedAt             time.Time                  `json:"created_at"`
	UpdatedAt             time.Time                  `json:"updated_at"`
}

// ==================== WORK TRACKING MODELS ====================

// WorkItem represents a task or work item
type WorkItem struct {
	ID                  string                     `gorm:"primaryKey" json:"id"`
	TenantID            string                     `json:"tenant_id"`
	WorkTitle           string                     `json:"work_title"`
	WorkDescription     string                     `json:"work_description"`
	WorkType            string                     `json:"work_type"` // TASK, BUG, FEATURE, IMPROVEMENT, DOCUMENTATION
	Status              string                     `json:"status"`    // TODO, IN_PROGRESS, IN_REVIEW, BLOCKED, COMPLETED, CANCELLED
	Priority            string                     `json:"priority"`  // CRITICAL, HIGH, MEDIUM, LOW
	AssignedTo          string                     `json:"assigned_to"`
	CreatedBy           string                     `json:"created_by"`
	ParentItemID        string                     `json:"parent_item_id"`
	EstimatedHours      float64                    `json:"estimated_hours"`
	ActualHours         float64                    `json:"actual_hours"`
	DueDate             *time.Time                 `json:"due_date"`
	CompletedDate       *time.Time                 `json:"completed_date"`
	PercentageComplete  int                        `json:"percentage_complete"`
	Tags                datatypes.JSONType[string] `json:"tags"`
	Attachments         datatypes.JSONType[string] `json:"attachments"`
	LinkedChatChannelID string                     `json:"linked_chat_channel_id"`
	LinkedCallID        string                     `json:"linked_call_id"`
	LinkedEventID       string                     `json:"linked_event_id"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
	DeletedAt           *time.Time                 `gorm:"index" json:"deleted_at"`
}

// WorkItemComment represents a comment on a work item
type WorkItemComment struct {
	ID          string                     `gorm:"primaryKey" json:"id"`
	WorkItemID  string                     `json:"work_item_id"`
	TenantID    string                     `json:"tenant_id"`
	CommenterID string                     `json:"commenter_id"`
	CommentText string                     `json:"comment_text"`
	CommentType string                     `json:"comment_type"` // COMMENT, STATUS_UPDATE, ATTACHMENT, MENTION
	Mentions    datatypes.JSONType[string] `json:"mentions"`
	Attachments datatypes.JSONType[string] `json:"attachments"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
}

// WorkItemTimeLog represents time spent on a work item
type WorkItemTimeLog struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	WorkItemID   string    `json:"work_item_id"`
	TenantID     string    `json:"tenant_id"`
	UserID       string    `json:"user_id"`
	TimeSpentMin int       `json:"time_spent_minutes"`
	LogDate      time.Time `json:"log_date"`
	LogNotes     string    `json:"log_notes"`
	IsBillable   int       `json:"is_billable"`
	CreatedAt    time.Time `json:"created_at"`
}

// ==================== NOTIFICATION MODELS ====================

// UserNotification represents a real-time notification
type UserNotification struct {
	ID               string     `gorm:"primaryKey" json:"id"`
	TenantID         string     `json:"tenant_id"`
	UserID           string     `json:"user_id"`
	NotificationType string     `json:"notification_type"` // MESSAGE, CALL, TASK, EVENT, MENTION, SYSTEM
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	ReferenceID      string     `json:"reference_id"`
	ReferenceType    string     `json:"reference_type"`
	IsRead           int        `json:"is_read"`
	ReadAt           *time.Time `json:"read_at"`
	ActionURL        string     `json:"action_url"`
	CreatedAt        time.Time  `json:"created_at"`
	ExpiresAt        *time.Time `json:"expires_at"`
}

// ==================== PRESENCE MODELS ====================

// UserPresence represents user online status and presence
type UserPresence struct {
	ID              string     `gorm:"primaryKey" json:"id"`
	TenantID        string     `json:"tenant_id"`
	UserID          string     `json:"user_id"`
	Status          string     `json:"status"` // ONLINE, AWAY, BUSY, DND, OFFLINE
	CurrentActivity string     `json:"current_activity"`
	LastSeen        *time.Time `json:"last_seen"`
	SessionID       string     `json:"session_id"`
	DeviceType      string     `json:"device_type"`
	IPAddress       string     `json:"ip_address"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TeamChatSendMessageRequest is the request to send a chat message
type TeamChatSendMessageRequest struct {
	ChannelID   string   `json:"channel_id" binding:"required"`
	MessageType string   `json:"message_type" binding:"required"`
	MessageBody string   `json:"message_body" binding:"required"`
	FileURL     string   `json:"file_url"`
	FileName    string   `json:"file_name"`
	Mentions    []string `json:"mentions"`
}

// ==================== REQUEST/RESPONSE TYPES ====================

// InitiateCallRequest is the request to initiate a call
type InitiateCallRequest struct {
	CallType       string                 `json:"call_type" binding:"required"`
	ParticipantIDs []string               `json:"participant_ids" binding:"required"`
	IsAudioEnabled bool                   `json:"is_audio_enabled"`
	IsVideoEnabled bool                   `json:"is_video_enabled"`
	LinkedRoomID   string                 `json:"linked_room_id"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// CreateEventRequest is the request to create a calendar event
type CreateEventRequest struct {
	EventTitle        string    `json:"event_title" binding:"required"`
	EventDescription  string    `json:"event_description"`
	EventType         string    `json:"event_type" binding:"required"`
	StartTime         time.Time `json:"start_time" binding:"required"`
	EndTime           time.Time `json:"end_time" binding:"required"`
	AttendeeIDs       []string  `json:"attendee_ids"`
	LinkedRoomID      string    `json:"linked_room_id"`
	Location          string    `json:"location"`
	IsRecurring       bool      `json:"is_recurring"`
	RecurrencePattern string    `json:"recurrence_pattern"`
}

// CreateDialerCampaignRequest is the request to create a dialer campaign
type CreateDialerCampaignRequest struct {
	CampaignName       string    `json:"campaign_name" binding:"required"`
	CampaignType       string    `json:"campaign_type" binding:"required"`
	Description        string    `json:"description"`
	ScriptID           string    `json:"script_id"`
	DialStrategy       string    `json:"dial_strategy"`
	MaxConcurrentCalls int       `json:"max_concurrent_calls"`
	CallerIDNumber     string    `json:"caller_id_number"`
	ScheduledStartTime time.Time `json:"scheduled_start_time"`
	ScheduledEndTime   time.Time `json:"scheduled_end_time"`
}

// CreateWorkItemRequest is the request to create a work item
type CreateWorkItemRequest struct {
	WorkTitle       string    `json:"work_title" binding:"required"`
	WorkDescription string    `json:"work_description"`
	WorkType        string    `json:"work_type" binding:"required"`
	Priority        string    `json:"priority"`
	AssignedTo      string    `json:"assigned_to"`
	DueDate         time.Time `json:"due_date"`
	ParentItemID    string    `json:"parent_item_id"`
	EstimatedHours  float64   `json:"estimated_hours"`
	LinkedChannelID string    `json:"linked_channel_id"`
}

// MessageResponse is the response after sending a message
type MessageResponse struct {
	MessageID string    `json:"message_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Error     string    `json:"error,omitempty"`
}

// CallResponse is the response after initiating a call
type CallResponse struct {
	CallID          string                   `json:"call_id"`
	Status          string                   `json:"status"`
	WebRTCRoomID    string                   `json:"webrtc_room_id"`
	SignalingServer string                   `json:"signaling_server"`
	STUNServers     []string                 `json:"stun_servers"`
	TURNServers     []map[string]interface{} `json:"turn_servers"`
	CreatedAt       time.Time                `json:"created_at"`
	Error           string                   `json:"error,omitempty"`
}

// EventResponse is the response after creating an event
type EventResponse struct {
	EventID         string    `json:"event_id"`
	Status          string    `json:"status"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	CreatedAt       time.Time `json:"created_at"`
	InvitationsSent int       `json:"invitations_sent"`
	Error           string    `json:"error,omitempty"`
}

// DialerCampaignResponse is the response after creating a dialer campaign
type DialerCampaignResponse struct {
	CampaignID    string    `json:"campaign_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	TotalContacts int       `json:"total_contacts"`
	Error         string    `json:"error,omitempty"`
}

// WorkItemResponse is the response after creating a work item
type WorkItemResponse struct {
	WorkItemID string    `json:"work_item_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	Error      string    `json:"error,omitempty"`
}
