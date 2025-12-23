package models

import (
	"time"
)

// CustomerProfile represents a customer's profile information
type CustomerProfile struct {
	ID                          int64      `gorm:"primaryKey" json:"id"`
	TenantID                    int64      `json:"tenant_id"`
	UserID                      int64      `json:"user_id"`
	BookingID                   *int64     `json:"booking_id"`
	PhoneNumber                 *string    `json:"phone_number"`
	AlternatePhone              *string    `json:"alternate_phone"`
	EmailAddress                *string    `json:"email_address"`
	AlternateEmail              *string    `json:"alternate_email"`
	DateOfBirth                 *time.Time `json:"date_of_birth"`
	Gender                      *string    `json:"gender"`
	CurrentAddress              *string    `json:"current_address"`
	PermanentAddress            *string    `json:"permanent_address"`
	City                        *string    `json:"city"`
	State                       *string    `json:"state"`
	PostalCode                  *string    `json:"postal_code"`
	Country                     *string    `json:"country"`
	IDProofType                 *string    `json:"id_proof_type"`
	IDProofNumber               *string    `json:"id_proof_number"`
	IDProofFileURL              *string    `json:"id_proof_file_url"`
	S3Bucket                    *string    `json:"s3_bucket"`
	S3Key                       *string    `json:"s3_key"`
	CommunicationPreference     *string    `json:"communication_preference"`
	LanguagePreference          *string    `json:"language_preference"`
	Timezone                    *string    `json:"timezone"`
	NotificationEnabled         bool       `json:"notification_enabled"`
	EmailUpdatesEnabled         bool       `json:"email_updates_enabled"`
	SMSUpdatesEnabled           bool       `json:"sms_updates_enabled"`
	PushNotificationsEnabled    bool       `json:"push_notifications_enabled"`
	ProfileCompletionPercentage float64    `json:"profile_completion_percentage"`
	IsVerified                  bool       `json:"is_verified"`
	VerifiedAt                  *time.Time `json:"verified_at"`
	VerificationNotes           *string    `json:"verification_notes"`
	Metadata                    *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy                   *int64     `json:"created_by"`
	UpdatedBy                   *int64     `json:"updated_by"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
	DeletedAt                   *time.Time `json:"deleted_at"`
}

// CustomerNotification represents a notification sent to customer
type CustomerNotification struct {
	ID                int64      `gorm:"primaryKey" json:"id"`
	TenantID          int64      `json:"tenant_id"`
	UserID            int64      `json:"user_id"`
	NotificationType  string     `json:"notification_type"`
	Title             string     `json:"title"`
	Message           string     `json:"message"`
	Description       *string    `json:"description"`
	Category          string     `json:"category"`
	Priority          string     `json:"priority"`
	RelatedEntityType *string    `json:"related_entity_type"`
	RelatedEntityID   *int64     `json:"related_entity_id"`
	IsRead            bool       `json:"is_read"`
	ReadAt            *time.Time `json:"read_at"`
	IsArchived        bool       `json:"is_archived"`
	ArchivedAt        *time.Time `json:"archived_at"`
	DeliveryStatus    string     `json:"delivery_status"`
	EmailSent         bool       `json:"email_sent"`
	EmailSentAt       *time.Time `json:"email_sent_at"`
	SMSSent           bool       `json:"sms_sent"`
	SMSSentAt         *time.Time `json:"sms_sent_at"`
	PushSent          bool       `json:"push_sent"`
	PushSentAt        *time.Time `json:"push_sent_at"`
	InAppSent         bool       `json:"in_app_sent"`
	InAppSentAt       *time.Time `json:"in_app_sent_at"`
	ActionURL         *string    `json:"action_url"`
	CTAText           *string    `json:"cta_text"`
	Metadata          *JSONMap   `gorm:"type:json" json:"metadata"`
	ExpiresAt         *time.Time `json:"expires_at"`
	CreatedBy         *int64     `json:"created_by"`
	UpdatedBy         *int64     `json:"updated_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

// CustomerConversation represents a support conversation
type CustomerConversation struct {
	ID               int64      `gorm:"primaryKey" json:"id"`
	TenantID         int64      `json:"tenant_id"`
	CustomerUserID   int64      `json:"customer_user_id"`
	SupportUserID    *int64     `json:"support_user_id"`
	BookingID        *int64     `json:"booking_id"`
	Subject          *string    `json:"subject"`
	ConversationType string     `json:"conversation_type"`
	Status           string     `json:"status"`
	Priority         string     `json:"priority"`
	LastMessageAt    *time.Time `json:"last_message_at"`
	LastMessageFrom  *string    `json:"last_message_from"`
	AssignedTo       *int64     `json:"assigned_to"`
	AssignedAt       *time.Time `json:"assigned_at"`
	ResolutionNotes  *string    `json:"resolution_notes"`
	ResolvedAt       *time.Time `json:"resolved_at"`
	ResolvedBy       *int64     `json:"resolved_by"`
	Metadata         *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy        *int64     `json:"created_by"`
	UpdatedBy        *int64     `json:"updated_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

// CustomerMessage represents a message within a conversation
type CustomerMessage struct {
	ID                 int64      `gorm:"primaryKey" json:"id"`
	TenantID           int64      `json:"tenant_id"`
	ConversationID     int64      `json:"conversation_id"`
	SenderUserID       int64      `json:"sender_user_id"`
	SenderType         string     `json:"sender_type"`
	MessageText        string     `json:"message_text"`
	MessageType        string     `json:"message_type"`
	AttachmentURL      *string    `json:"attachment_url"`
	AttachmentFileName *string    `json:"attachment_file_name"`
	AttachmentFileSize *int64     `json:"attachment_file_size"`
	AttachmentFileType *string    `json:"attachment_file_type"`
	S3Bucket           *string    `json:"s3_bucket"`
	S3Key              *string    `json:"s3_key"`
	IsRead             bool       `json:"is_read"`
	ReadAt             *time.Time `json:"read_at"`
	Metadata           *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
}

// CustomerDocumentUpload represents a customer document upload
type CustomerDocumentUpload struct {
	ID                       int64      `gorm:"primaryKey" json:"id"`
	TenantID                 int64      `json:"tenant_id"`
	UserID                   int64      `json:"user_id"`
	BookingID                *int64     `json:"booking_id"`
	DocumentType             *string    `json:"document_type"`
	DocumentName             string     `json:"document_name"`
	DocumentDescription      *string    `json:"document_description"`
	FileURL                  *string    `json:"file_url"`
	FileName                 string     `json:"file_name"`
	FileSize                 int64      `json:"file_size"`
	FileExtension            *string    `json:"file_extension"`
	FileMimeType             *string    `json:"file_mime_type"`
	S3Bucket                 *string    `json:"s3_bucket"`
	S3Key                    *string    `json:"s3_key"`
	UploadStatus             string     `json:"upload_status"`
	UploadProgressPercentage float64    `json:"upload_progress_percentage"`
	VerificationStatus       string     `json:"verification_status"`
	VerificationNotes        *string    `json:"verification_notes"`
	VerifiedBy               *int64     `json:"verified_by"`
	VerifiedAt               *time.Time `json:"verified_at"`
	IsRequiredDocument       bool       `json:"is_required_document"`
	IsVerifiedByAdmin        bool       `json:"is_verified_by_admin"`
	Metadata                 *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy                *int64     `json:"created_by"`
	UpdatedBy                *int64     `json:"updated_by"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at"`
}

// CustomerBookingTracking tracks booking progress
type CustomerBookingTracking struct {
	ID                     int64      `gorm:"primaryKey" json:"id"`
	TenantID               int64      `json:"tenant_id"`
	BookingID              int64      `json:"booking_id"`
	UserID                 int64      `json:"user_id"`
	CurrentStatus          *string    `json:"current_status"`
	PreviousStatus         *string    `json:"previous_status"`
	StatusChangedAt        *time.Time `json:"status_changed_at"`
	StatusChangeReason     *string    `json:"status_change_reason"`
	PropertyID             *int64     `json:"property_id"`
	PropertyName           *string    `json:"property_name"`
	PropertyLocation       *string    `json:"property_location"`
	UnitNumber             *string    `json:"unit_number"`
	BookingDate            *time.Time `json:"booking_date"`
	PossessionDate         *time.Time `json:"possession_date"`
	EstimatedHandoverDate  *time.Time `json:"estimated_handover_date"`
	TotalAmount            float64    `json:"total_amount"`
	AmountPaid             float64    `json:"amount_paid"`
	AmountPending          float64    `json:"amount_pending"`
	PaymentPercentage      float64    `json:"payment_percentage"`
	TotalMilestones        *int       `json:"total_milestones"`
	CompletedMilestones    int        `json:"completed_milestones"`
	PendingMilestones      int        `json:"pending_milestones"`
	RequiredDocumentsCount *int       `json:"required_documents_count"`
	UploadedDocumentsCount int        `json:"uploaded_documents_count"`
	VerifiedDocumentsCount int        `json:"verified_documents_count"`
	LastUpdateAt           *time.Time `json:"last_update_at"`
	LastUpdateType         *string    `json:"last_update_type"`
	Metadata               *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy              *int64     `json:"created_by"`
	UpdatedBy              *int64     `json:"updated_by"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at"`
}

// CustomerPaymentTracking tracks payment information
type CustomerPaymentTracking struct {
	ID                 int64      `gorm:"primaryKey" json:"id"`
	TenantID           int64      `json:"tenant_id"`
	BookingID          int64      `json:"booking_id"`
	UserID             int64      `json:"user_id"`
	InvoiceNumber      *string    `json:"invoice_number"`
	InvoiceDate        *time.Time `json:"invoice_date"`
	DueDate            *time.Time `json:"due_date"`
	InvoiceAmount      float64    `json:"invoice_amount"`
	TaxAmount          float64    `json:"tax_amount"`
	TotalAmount        float64    `json:"total_amount"`
	PaymentStatus      string     `json:"payment_status"`
	PaymentMethod      *string    `json:"payment_method"`
	TransactionID      *string    `json:"transaction_id"`
	PaymentDate        *time.Time `json:"payment_date"`
	AmountPaid         float64    `json:"amount_paid"`
	IsRefunded         bool       `json:"is_refunded"`
	RefundAmount       *float64   `json:"refund_amount"`
	RefundDate         *time.Time `json:"refund_date"`
	RefundReason       *string    `json:"refund_reason"`
	ReminderSentCount  int        `json:"reminder_sent_count"`
	LastReminderSentAt *time.Time `json:"last_reminder_sent_at"`
	InvoiceFileURL     *string    `json:"invoice_file_url"`
	S3Bucket           *string    `json:"s3_bucket"`
	S3Key              *string    `json:"s3_key"`
	Metadata           *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy          *int64     `json:"created_by"`
	UpdatedBy          *int64     `json:"updated_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
}

// CustomerFeedback represents customer feedback and ratings
type CustomerFeedback struct {
	ID                  int64      `gorm:"primaryKey" json:"id"`
	TenantID            int64      `json:"tenant_id"`
	UserID              int64      `json:"user_id"`
	BookingID           *int64     `json:"booking_id"`
	FeedbackType        string     `json:"feedback_type"`
	Subject             *string    `json:"subject"`
	Message             string     `json:"message"`
	OverallRating       *float32   `json:"overall_rating"`
	ServiceRating       *float32   `json:"service_rating"`
	CommunicationRating *float32   `json:"communication_rating"`
	DocumentationRating *float32   `json:"documentation_rating"`
	FeedbackStatus      string     `json:"feedback_status"`
	ResponseMessage     *string    `json:"response_message"`
	RespondedBy         *int64     `json:"responded_by"`
	RespondedAt         *time.Time `json:"responded_at"`
	AttachmentsCount    int        `json:"attachments_count"`
	Metadata            *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy           *int64     `json:"created_by"`
	UpdatedBy           *int64     `json:"updated_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

// CustomerActivityLog represents activity tracking
type CustomerActivityLog struct {
	ID                  int64     `gorm:"primaryKey" json:"id"`
	TenantID            int64     `json:"tenant_id"`
	UserID              int64     `json:"user_id"`
	ActivityType        string    `json:"activity_type"`
	ActivityDescription *string   `json:"activity_description"`
	EntityType          *string   `json:"entity_type"`
	EntityID            *int64    `json:"entity_id"`
	ActionTaken         *string   `json:"action_taken"`
	OldValue            *string   `json:"old_value"`
	NewValue            *string   `json:"new_value"`
	IPAddress           *string   `json:"ip_address"`
	UserAgent           *string   `json:"user_agent"`
	DeviceType          *string   `json:"device_type"`
	Metadata            *JSONMap  `gorm:"type:json" json:"metadata"`
	CreatedAt           time.Time `json:"created_at"`
}

// CustomerPreferences represents customer preferences
type CustomerPreferences struct {
	ID                       int64      `gorm:"primaryKey" json:"id"`
	TenantID                 int64      `json:"tenant_id"`
	UserID                   int64      `json:"user_id"`
	EmailNotifications       bool       `json:"email_notifications"`
	SMSNotifications         bool       `json:"sms_notifications"`
	PushNotifications        bool       `json:"push_notifications"`
	InAppNotifications       bool       `json:"in_app_notifications"`
	ReceiveBookingUpdates    bool       `json:"receive_booking_updates"`
	ReceivePaymentReminders  bool       `json:"receive_payment_reminders"`
	ReceiveDocumentRequests  bool       `json:"receive_document_requests"`
	ReceivePossessionUpdates bool       `json:"receive_possession_updates"`
	ReceivePromotionalEmails bool       `json:"receive_promotional_emails"`
	ReceiveNewsletter        bool       `json:"receive_newsletter"`
	PreferredLanguage        *string    `json:"preferred_language"`
	PreferredTimezone        *string    `json:"preferred_timezone"`
	IsProfilePublic          bool       `json:"is_profile_public"`
	AllowContactByPhone      bool       `json:"allow_contact_by_phone"`
	AllowContactByEmail      bool       `json:"allow_contact_by_email"`
	AllowContactBySMS        bool       `json:"allow_contact_by_sms"`
	AllowMarketing           bool       `json:"allow_marketing"`
	DashboardLayout          *string    `json:"dashboard_layout"`
	ThemePreference          *string    `json:"theme_preference"`
	Metadata                 *JSONMap   `gorm:"type:json" json:"metadata"`
	CreatedBy                *int64     `json:"created_by"`
	UpdatedBy                *int64     `json:"updated_by"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at"`
}

// ============================================
// DTO Models for Request/Response Handling
// ============================================

// CreateCustomerPortalProfileRequest represents a request to create customer profile
type CreateCustomerPortalProfileRequest struct {
	PhoneNumber             *string    `json:"phone_number"`
	EmailAddress            *string    `json:"email_address"`
	DateOfBirth             *time.Time `json:"date_of_birth"`
	Gender                  *string    `json:"gender"`
	CurrentAddress          *string    `json:"current_address"`
	City                    *string    `json:"city"`
	State                   *string    `json:"state"`
	PostalCode              *string    `json:"postal_code"`
	Country                 *string    `json:"country"`
	CommunicationPreference *string    `json:"communication_preference"`
	LanguagePreference      *string    `json:"language_preference"`
	Timezone                *string    `json:"timezone"`
}

// UpdateCustomerPortalProfileRequest represents a request to update customer profile
type UpdateCustomerPortalProfileRequest struct {
	PhoneNumber        *string `json:"phone_number"`
	AlternatePhone     *string `json:"alternate_phone"`
	EmailAddress       *string `json:"email_address"`
	CurrentAddress     *string `json:"current_address"`
	City               *string `json:"city"`
	State              *string `json:"state"`
	PostalCode         *string `json:"postal_code"`
	LanguagePreference *string `json:"language_preference"`
	Timezone           *string `json:"timezone"`
}

// CreatePortalNotificationRequest represents a request to create notification
type CreatePortalNotificationRequest struct {
	NotificationType string `json:"notification_type" binding:"required"`
	Title            string `json:"title" binding:"required"`
	Message          string `json:"message" binding:"required"`
	Category         string `json:"category"`
	Priority         string `json:"priority"`
}

// CreatePortalConversationRequest represents a request to create a conversation
type CreatePortalConversationRequest struct {
	Subject          *string `json:"subject"`
	ConversationType string  `json:"conversation_type"`
	BookingID        *int64  `json:"booking_id"`
}

// SendPortalMessageRequest represents a request to send a message
type SendPortalMessageRequest struct {
	MessageText string `json:"message_text" binding:"required"`
	MessageType string `json:"message_type"`
}

// CreatePortalDocumentUploadRequest represents a request to upload a document
type CreatePortalDocumentUploadRequest struct {
	DocumentType        string  `json:"document_type"`
	DocumentName        string  `json:"document_name" binding:"required"`
	DocumentDescription *string `json:"document_description"`
	IsRequiredDocument  bool    `json:"is_required_document"`
}

// UpdatePortalPaymentStatusRequest represents a request to update payment status
type UpdatePortalPaymentStatusRequest struct {
	PaymentStatus string     `json:"payment_status"`
	AmountPaid    *float64   `json:"amount_paid"`
	PaymentDate   *time.Time `json:"payment_date"`
}

// CreatePortalFeedbackRequest represents a request to create feedback
type CreatePortalFeedbackRequest struct {
	FeedbackType        string   `json:"feedback_type" binding:"required"`
	Subject             *string  `json:"subject"`
	Message             string   `json:"message" binding:"required"`
	OverallRating       *float32 `json:"overall_rating"`
	ServiceRating       *float32 `json:"service_rating"`
	CommunicationRating *float32 `json:"communication_rating"`
	DocumentationRating *float32 `json:"documentation_rating"`
}

// UpdatePortalPreferencesRequest represents a request to update preferences
type UpdatePortalPreferencesRequest struct {
	EmailNotifications      *bool   `json:"email_notifications"`
	SMSNotifications        *bool   `json:"sms_notifications"`
	PushNotifications       *bool   `json:"push_notifications"`
	ReceiveBookingUpdates   *bool   `json:"receive_booking_updates"`
	ReceivePaymentReminders *bool   `json:"receive_payment_reminders"`
	PreferredLanguage       *string `json:"preferred_language"`
	ThemePreference         *string `json:"theme_preference"`
}

// CustomerDashboardResponse represents customer dashboard data
type CustomerDashboardResponse struct {
	Profile                    *CustomerProfile           `json:"profile"`
	BookingTracking            *CustomerBookingTracking   `json:"booking_tracking"`
	RecentPayments             []*CustomerPaymentTracking `json:"recent_payments"`
	PendingDocuments           []*CustomerDocumentUpload  `json:"pending_documents"`
	UnreadNotifications        int                        `json:"unread_notifications"`
	CompletionPercentage       float64                    `json:"completion_percentage"`
	NextPaymentDue             *time.Time                 `json:"next_payment_due"`
	DocumentVerificationStatus string                     `json:"document_verification_status"`
	UpdatedAt                  time.Time                  `json:"updated_at"`
}
