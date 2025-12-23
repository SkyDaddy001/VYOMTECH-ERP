package models

import (
	"time"
)

// ============================================================
// MOBILE APP
// ============================================================

// MobileApp represents a mobile application configuration
type MobileApp struct {
	ID                   string     `json:"id" db:"id"`
	TenantID             string     `json:"tenant_id" db:"tenant_id"`
	AppName              string     `json:"app_name" db:"app_name"`
	AppType              string     `json:"app_type" db:"app_type"` // ios, android, cross-platform
	BundleIdentifier     string     `json:"bundle_identifier" db:"bundle_identifier"`
	Version              string     `json:"version" db:"version"`
	BuildNumber          int        `json:"build_number" db:"build_number"`
	Status               string     `json:"status" db:"status"` // active, inactive, deprecated, beta
	Description          *string    `json:"description" db:"description"`
	APIKey               string     `json:"api_key" db:"api_key"`
	APISecret            string     `json:"-" db:"api_secret"` // never expose in response
	MinSupportedVersion  *string    `json:"min_supported_version" db:"min_supported_version"`
	MaxSupportedVersion  *string    `json:"max_supported_version" db:"max_supported_version"`
	StoreURLIOS          *string    `json:"store_url_ios" db:"store_url_ios"`
	StoreURLAndroid      *string    `json:"store_url_android" db:"store_url_android"`
	AppIconURL           *string    `json:"app_icon_url" db:"app_icon_url"`
	AppBannerURL         *string    `json:"app_banner_url" db:"app_banner_url"`
	SupportEmail         *string    `json:"support_email" db:"support_email"`
	SupportPhone         *string    `json:"support_phone" db:"support_phone"`
	SupportChatURL       *string    `json:"support_chat_url" db:"support_chat_url"`
	PrivacyPolicyURL     *string    `json:"privacy_policy_url" db:"privacy_policy_url"`
	TermsOfServiceURL    *string    `json:"terms_of_service_url" db:"terms_of_service_url"`
	ChangelogURL         *string    `json:"changelog_url" db:"changelog_url"`
	FeatureFlags         *string    `json:"feature_flags" db:"feature_flags"` // JSON
	Metadata             *string    `json:"metadata" db:"metadata"`           // JSON
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at" db:"deleted_at"`
}

// CreateMobileAppRequest is the DTO for creating a mobile app
type CreateMobileAppRequest struct {
	AppName              string  `json:"app_name" binding:"required"`
	AppType              string  `json:"app_type" binding:"required"` // ios, android, cross-platform
	BundleIdentifier     string  `json:"bundle_identifier" binding:"required"`
	Version              string  `json:"version" binding:"required"`
	BuildNumber          int     `json:"build_number" binding:"required"`
	Description          *string `json:"description"`
	MinSupportedVersion  *string `json:"min_supported_version"`
	MaxSupportedVersion  *string `json:"max_supported_version"`
	StoreURLIOS          *string `json:"store_url_ios"`
	StoreURLAndroid      *string `json:"store_url_android"`
	SupportEmail         *string `json:"support_email"`
	SupportPhone         *string `json:"support_phone"`
	PrivacyPolicyURL     *string `json:"privacy_policy_url"`
	TermsOfServiceURL    *string `json:"terms_of_service_url"`
}

// UpdateMobileAppRequest is the DTO for updating a mobile app
type UpdateMobileAppRequest struct {
	AppName              *string `json:"app_name"`
	Version              *string `json:"version"`
	BuildNumber          *int    `json:"build_number"`
	Status               *string `json:"status"`
	Description          *string `json:"description"`
	MinSupportedVersion  *string `json:"min_supported_version"`
	MaxSupportedVersion  *string `json:"max_supported_version"`
	SupportEmail         *string `json:"support_email"`
	SupportPhone         *string `json:"support_phone"`
}

// ============================================================
// MOBILE DEVICE
// ============================================================

// MobileDevice represents a registered mobile device
type MobileDevice struct {
	ID                        string     `json:"id" db:"id"`
	TenantID                  string     `json:"tenant_id" db:"tenant_id"`
	UserID                    int64      `json:"user_id" db:"user_id"`
	AppID                     string     `json:"app_id" db:"app_id"`
	DeviceID                  string     `json:"device_id" db:"device_id"`
	DeviceName                *string    `json:"device_name" db:"device_name"`
	DeviceModel               *string    `json:"device_model" db:"device_model"`
	DeviceManufacturer        *string    `json:"device_manufacturer" db:"device_manufacturer"`
	OSType                    string     `json:"os_type" db:"os_type"` // ios, android, windows
	OSVersion                 *string    `json:"os_version" db:"os_version"`
	AppVersion                *string    `json:"app_version" db:"app_version"`
	AppBuild                  *int       `json:"app_build" db:"app_build"`
	PushToken                 *string    `json:"push_token" db:"push_token"`
	PushTokenUpdatedAt        *time.Time `json:"push_token_updated_at" db:"push_token_updated_at"`
	DeviceUUID                *string    `json:"device_uuid" db:"device_uuid"`
	IMEI                      *string    `json:"imei" db:"imei"`
	DeviceSerial              *string    `json:"device_serial" db:"device_serial"`
	ScreenResolution          *string    `json:"screen_resolution" db:"screen_resolution"`
	ScreenDensity             *string    `json:"screen_density" db:"screen_density"`
	Locale                    *string    `json:"locale" db:"locale"`
	Timezone                  *string    `json:"timezone" db:"timezone"`
	BatteryOptimizationEnabled bool      `json:"battery_optimization_enabled" db:"battery_optimization_enabled"`
	BiometricEnabled          bool      `json:"biometric_enabled" db:"biometric_enabled"`
	BiometricType             *string    `json:"biometric_type" db:"biometric_type"` // fingerprint, face, iris
	FCMToken                  *string    `json:"fcm_token" db:"fcm_token"`
	APNSToken                 *string    `json:"apns_token" db:"apns_token"`
	DeviceStatus              string     `json:"device_status" db:"device_status"` // active, inactive, suspended, lost
	LastActivityAt            *time.Time `json:"last_activity_at" db:"last_activity_at"`
	FirstSeenAt               time.Time  `json:"first_seen_at" db:"first_seen_at"`
	LastSeenAt                *time.Time `json:"last_seen_at" db:"last_seen_at"`
	Metadata                  *string    `json:"metadata" db:"metadata"` // JSON
	CreatedAt                 time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                 *time.Time `json:"deleted_at" db:"deleted_at"`
}

// RegisterDeviceRequest is the DTO for registering a device
type RegisterDeviceRequest struct {
	DeviceID           string  `json:"device_id" binding:"required"`
	DeviceName         *string `json:"device_name"`
	DeviceModel        *string `json:"device_model"`
	DeviceManufacturer *string `json:"device_manufacturer"`
	OSType             string  `json:"os_type" binding:"required"`
	OSVersion          *string `json:"os_version"`
	AppVersion         *string `json:"app_version"`
	AppBuild           *int    `json:"app_build"`
	PushToken          *string `json:"push_token"`
	BiometricEnabled   *bool   `json:"biometric_enabled"`
	BiometricType      *string `json:"biometric_type"`
}

// UpdateDeviceRequest is the DTO for updating device info
type UpdateDeviceRequest struct {
	PushToken                 *string `json:"push_token"`
	DeviceStatus              *string `json:"device_status"`
	BiometricEnabled          *bool   `json:"biometric_enabled"`
	BatteryOptimizationEnabled *bool  `json:"battery_optimization_enabled"`
}

// ============================================================
// MOBILE SESSION
// ============================================================

// MobileSession represents a user session on mobile
type MobileSession struct {
	ID                      string     `json:"id" db:"id"`
	TenantID                string     `json:"tenant_id" db:"tenant_id"`
	UserID                  int64      `json:"user_id" db:"user_id"`
	DeviceID                string     `json:"device_id" db:"device_id"`
	AppID                   string     `json:"app_id" db:"app_id"`
	SessionToken            string     `json:"session_token" db:"session_token"`
	RefreshToken            *string    `json:"refresh_token" db:"refresh_token"`
	TokenExpiresAt          time.Time  `json:"token_expires_at" db:"token_expires_at"`
	RefreshTokenExpiresAt   *time.Time `json:"refresh_token_expires_at" db:"refresh_token_expires_at"`
	SessionStatus           string     `json:"session_status" db:"session_status"` // active, inactive, expired, revoked
	IPAddress               *string    `json:"ip_address" db:"ip_address"`
	UserAgent               *string    `json:"user_agent" db:"user_agent"`
	AppVersion              *string    `json:"app_version" db:"app_version"`
	LoginMethod             *string    `json:"login_method" db:"login_method"` // password, sso, biometric, otp
	LoginTimestamp          time.Time  `json:"login_timestamp" db:"login_timestamp"`
	LastActivityTimestamp   *time.Time `json:"last_activity_timestamp" db:"last_activity_timestamp"`
	LogoutTimestamp         *time.Time `json:"logout_timestamp" db:"logout_timestamp"`
	SessionDurationSeconds  *int       `json:"session_duration_seconds" db:"session_duration_seconds"`
	DeviceTrusted           bool       `json:"device_trusted" db:"device_trusted"`
	TwoFactorVerified       bool       `json:"two_factor_verified" db:"two_factor_verified"`
	NetworkType             *string    `json:"network_type" db:"network_type"` // wifi, cellular, unknown
	NetworkProvider         *string    `json:"network_provider" db:"network_provider"`
	AppInForeground         bool       `json:"app_in_foreground" db:"app_in_foreground"`
	BackgroundActivityAllowed bool     `json:"background_activity_allowed" db:"background_activity_allowed"`
	Metadata                *string    `json:"metadata" db:"metadata"` // JSON
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt               *time.Time `json:"deleted_at" db:"deleted_at"`
}

// CreateSessionRequest is the DTO for creating a session
type CreateSessionRequest struct {
	DeviceID             string  `json:"device_id" binding:"required"`
	AppID                string  `json:"app_id" binding:"required"`
	LoginMethod          string  `json:"login_method"`
	IPAddress            *string `json:"ip_address"`
	UserAgent            *string `json:"user_agent"`
	DeviceTrusted        *bool   `json:"device_trusted"`
	BiometricVerified    *bool   `json:"biometric_verified"`
}

// ============================================================
// APP NOTIFICATION
// ============================================================

// AppNotification represents a push notification
type AppNotification struct {
	ID                string     `json:"id" db:"id"`
	TenantID          string     `json:"tenant_id" db:"tenant_id"`
	AppID             string     `json:"app_id" db:"app_id"`
	UserID            *int64     `json:"user_id" db:"user_id"`
	DeviceID          *string    `json:"device_id" db:"device_id"`
	NotificationType  string     `json:"notification_type" db:"notification_type"`
	NotificationCategory *string `json:"notification_category" db:"notification_category"`
	Title             string     `json:"title" db:"title"`
	Body              string     `json:"body" db:"body"`
	ContentType       *string    `json:"content_type" db:"content_type"`
	ImageURL          *string    `json:"image_url" db:"image_url"`
	ActionURL         *string    `json:"action_url" db:"action_url"`
	ActionType        *string    `json:"action_type" db:"action_type"`
	Priority          *string    `json:"priority" db:"priority"`
	Sound             *string    `json:"sound" db:"sound"`
	BadgeCount        *int       `json:"badge_count" db:"badge_count"`
	CustomData        *string    `json:"custom_data" db:"custom_data"` // JSON
	ScheduledTime     *time.Time `json:"scheduled_time" db:"scheduled_time"`
	SendTime          *time.Time `json:"send_time" db:"send_time"`
	DeliveryStatus    string     `json:"delivery_status" db:"delivery_status"`
	DeliveryAttempts  int        `json:"delivery_attempts" db:"delivery_attempts"`
	NextRetryAt       *time.Time `json:"next_retry_at" db:"next_retry_at"`
	ReadAt            *time.Time `json:"read_at" db:"read_at"`
	ClickedAt         *time.Time `json:"clicked_at" db:"clicked_at"`
	DismissedAt       *time.Time `json:"dismissed_at" db:"dismissed_at"`
	ExpiryTime        *time.Time `json:"expiry_time" db:"expiry_time"`
	IsCampaign        bool       `json:"is_campaign" db:"is_campaign"`
	CampaignID        *string    `json:"campaign_id" db:"campaign_id"`
	SegmentID         *string    `json:"segment_id" db:"segment_id"`
	ABTestVariant     *string    `json:"ab_test_variant" db:"ab_test_variant"`
	RetryCount        int        `json:"retry_count" db:"retry_count"`
	MaxRetries        int        `json:"max_retries" db:"max_retries"`
	Metadata          *string    `json:"metadata" db:"metadata"` // JSON
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at" db:"deleted_at"`
}

// SendNotificationRequest is the DTO for sending notifications
type SendNotificationRequest struct {
	NotificationType  string   `json:"notification_type" binding:"required"`
	Title             string   `json:"title" binding:"required"`
	Body              string   `json:"body" binding:"required"`
	UserIDs           *[]int64 `json:"user_ids"` // if nil, broadcast to all
	DeviceIDs         *[]string `json:"device_ids"`
	Priority          *string  `json:"priority"`
	ScheduledTime     *time.Time `json:"scheduled_time"`
	ActionURL         *string  `json:"action_url"`
	CustomData        *string  `json:"custom_data"`
}

// ============================================================
// APP FEATURE
// ============================================================

// AppFeature represents a feature flag
type AppFeature struct {
	ID                  string     `json:"id" db:"id"`
	TenantID            string     `json:"tenant_id" db:"tenant_id"`
	AppID               string     `json:"app_id" db:"app_id"`
	FeatureName         string     `json:"feature_name" db:"feature_name"`
	FeatureCode         string     `json:"feature_code" db:"feature_code"`
	FeatureDescription  *string    `json:"feature_description" db:"feature_description"`
	FeatureCategory     *string    `json:"feature_category" db:"feature_category"`
	IsEnabled           bool       `json:"is_enabled" db:"is_enabled"`
	MinAppVersion       *string    `json:"min_app_version" db:"min_app_version"`
	MaxAppVersion       *string    `json:"max_app_version" db:"max_app_version"`
	MinOSVersion        *string    `json:"min_os_version" db:"min_os_version"`
	RequiredPermissions *string    `json:"required_permissions" db:"required_permissions"` // JSON
	EnabledForUsers     *string    `json:"enabled_for_users" db:"enabled_for_users"`       // JSON
	DisabledForUsers    *string    `json:"disabled_for_users" db:"disabled_for_users"`     // JSON
	RolloutPercentage   int        `json:"rollout_percentage" db:"rollout_percentage"`
	ABTestVariant       *string    `json:"ab_test_variant" db:"ab_test_variant"`
	AnalyticsTrack      bool       `json:"analytics_track" db:"analytics_track"`
	RequiresConsent     bool       `json:"requires_consent" db:"requires_consent"`
	ConsentType         *string    `json:"consent_type" db:"consent_type"`
	FeatureURL          *string    `json:"feature_url" db:"feature_url"`
	DocumentationURL    *string    `json:"documentation_url" db:"documentation_url"`
	SupportEmail        *string    `json:"support_email" db:"support_email"`
	LaunchDate          *time.Time `json:"launch_date" db:"launch_date"`
	SunsetDate          *time.Time `json:"sunset_date" db:"sunset_date"`
	Config              *string    `json:"config" db:"config"` // JSON
	Metadata            *string    `json:"metadata" db:"metadata"` // JSON
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at" db:"deleted_at"`
}

// CreateFeatureRequest is the DTO for creating a feature
type CreateFeatureRequest struct {
	FeatureName        string  `json:"feature_name" binding:"required"`
	FeatureCode        string  `json:"feature_code" binding:"required"`
	FeatureDescription *string `json:"feature_description"`
	FeatureCategory    *string `json:"feature_category"`
	IsEnabled          *bool   `json:"is_enabled"`
	MinAppVersion      *string `json:"min_app_version"`
	RolloutPercentage  *int    `json:"rollout_percentage"`
}

// ============================================================
// APP SETTING
// ============================================================

// AppSetting represents a user app preference
type AppSetting struct {
	ID               string     `json:"id" db:"id"`
	TenantID         string     `json:"tenant_id" db:"tenant_id"`
	UserID           int64      `json:"user_id" db:"user_id"`
	AppID            string     `json:"app_id" db:"app_id"`
	DeviceID         *string    `json:"device_id" db:"device_id"`
	SettingKey       string     `json:"setting_key" db:"setting_key"`
	SettingValue     *string    `json:"setting_value" db:"setting_value"`
	SettingType      *string    `json:"setting_type" db:"setting_type"`
	DisplayName      *string    `json:"display_name" db:"display_name"`
	Description      *string    `json:"description" db:"description"`
	Category         *string    `json:"category" db:"category"`
	IsUserEditable   bool       `json:"is_user_editable" db:"is_user_editable"`
	IsDeviceSpecific bool       `json:"is_device_specific" db:"is_device_specific"`
	DefaultValue     *string    `json:"default_value" db:"default_value"`
	ValidationRules  *string    `json:"validation_rules" db:"validation_rules"` // JSON
	Metadata         *string    `json:"metadata" db:"metadata"`                 // JSON
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at" db:"deleted_at"`
}

// UpdateSettingRequest is the DTO for updating settings
type UpdateSettingRequest struct {
	SettingValue *string `json:"setting_value" binding:"required"`
}

// ============================================================
// OFFLINE DATA
// ============================================================

// OfflineData represents cached data for offline access
type OfflineData struct {
	ID                        string     `json:"id" db:"id"`
	TenantID                  string     `json:"tenant_id" db:"tenant_id"`
	DeviceID                  string     `json:"device_id" db:"device_id"`
	AppID                     string     `json:"app_id" db:"app_id"`
	UserID                    int64      `json:"user_id" db:"user_id"`
	DataType                  string     `json:"data_type" db:"data_type"`
	DataKey                   string     `json:"data_key" db:"data_key"`
	DataSourceTable           *string    `json:"data_source_table" db:"data_source_table"`
	SourceRecordID            *string    `json:"source_record_id" db:"source_record_id"`
	CachedData                []byte     `json:"-" db:"cached_data"` // BLOB
	DataHash                  *string    `json:"data_hash" db:"data_hash"`
	CompressionType           *string    `json:"compression_type" db:"compression_type"`
	SyncStatus                string     `json:"sync_status" db:"sync_status"`
	LastSyncAt                *time.Time `json:"last_sync_at" db:"last_sync_at"`
	NeedsSync                 bool       `json:"needs_sync" db:"needs_sync"`
	LocalOnly                 bool       `json:"local_only" db:"local_only"`
	CachePriority             *string    `json:"cache_priority" db:"cache_priority"`
	ExpiryTime                *time.Time `json:"expiry_time" db:"expiry_time"`
	SizeBytes                 *int64     `json:"size_bytes" db:"size_bytes"`
	ConflictResolutionStrategy *string   `json:"conflict_resolution_strategy" db:"conflict_resolution_strategy"`
	Conflicts                 *string    `json:"conflicts" db:"conflicts"` // JSON
	Metadata                  *string    `json:"metadata" db:"metadata"`   // JSON
	CreatedAt                 time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                 *time.Time `json:"deleted_at" db:"deleted_at"`
}

// SyncOfflineDataRequest is the DTO for syncing offline data
type SyncOfflineDataRequest struct {
	DataKey                   string       `json:"data_key" binding:"required"`
	DataType                  string       `json:"data_type" binding:"required"`
	DataSourceTable           *string      `json:"data_source_table"`
	SourceRecordID            *string      `json:"source_record_id"`
	ConflictResolutionStrategy *string     `json:"conflict_resolution_strategy"`
}

// ============================================================
// APP CRASH
// ============================================================

// AppCrash represents a crash report
type AppCrash struct {
	ID                  string     `json:"id" db:"id"`
	TenantID            string     `json:"tenant_id" db:"tenant_id"`
	AppID               string     `json:"app_id" db:"app_id"`
	DeviceID            string     `json:"device_id" db:"device_id"`
	UserID              *int64     `json:"user_id" db:"user_id"`
	CrashTimestamp      time.Time  `json:"crash_timestamp" db:"crash_timestamp"`
	AppVersion          *string    `json:"app_version" db:"app_version"`
	AppBuild            *int       `json:"app_build" db:"app_build"`
	OSVersion           *string    `json:"os_version" db:"os_version"`
	DeviceModel         *string    `json:"device_model" db:"device_model"`
	CrashType           *string    `json:"crash_type" db:"crash_type"`
	CrashReason         *string    `json:"crash_reason" db:"crash_reason"`
	ExceptionType       *string    `json:"exception_type" db:"exception_type"`
	ExceptionMessage    *string    `json:"exception_message" db:"exception_message"`
	StackTrace          *string    `json:"stack_trace" db:"stack_trace"`
	Breadcrumbs         *string    `json:"breadcrumbs" db:"breadcrumbs"` // JSON
	UserReport          *string    `json:"user_report" db:"user_report"`
	Severity            *string    `json:"severity" db:"severity"`
	IsReproducible      bool       `json:"is_reproducible" db:"is_reproducible"`
	ReproductionSteps   *string    `json:"reproduction_steps" db:"reproduction_steps"`
	MemoryUsedMB        *int       `json:"memory_used_mb" db:"memory_used_mb"`
	MemoryAvailableMB   *int       `json:"memory_available_mb" db:"memory_available_mb"`
	BatteryLevel        *int       `json:"battery_level" db:"battery_level"`
	BatteryHealth       *string    `json:"battery_health" db:"battery_health"`
	StorageAvailableMB  *int       `json:"storage_available_mb" db:"storage_available_mb"`
	NetworkType         *string    `json:"network_type" db:"network_type"`
	CPUUsagePercent     *int       `json:"cpu_usage_percent" db:"cpu_usage_percent"`
	TemperatureCelsius  *int       `json:"temperature_celsius" db:"temperature_celsius"`
	SessionID           *string    `json:"session_id" db:"session_id"`
	CrashStatus         string     `json:"crash_status" db:"crash_status"`
	AssignedTo          *int64     `json:"assigned_to" db:"assigned_to"`
	AssignedAt          *time.Time `json:"assigned_at" db:"assigned_at"`
	ResolvedAt          *time.Time `json:"resolved_at" db:"resolved_at"`
	ResolutionNotes     *string    `json:"resolution_notes" db:"resolution_notes"`
	RelatedIssueID      *string    `json:"related_issue_id" db:"related_issue_id"`
	Metadata            *string    `json:"metadata" db:"metadata"` // JSON
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at" db:"deleted_at"`
}

// ReportCrashRequest is the DTO for reporting a crash
type ReportCrashRequest struct {
	CrashType          *string `json:"crash_type"`
	ExceptionType      *string `json:"exception_type"`
	ExceptionMessage   *string `json:"exception_message"`
	StackTrace         *string `json:"stack_trace" binding:"required"`
	UserReport         *string `json:"user_report"`
	DeviceState        *string `json:"device_state"`
	MemoryUsedMB       *int    `json:"memory_used_mb"`
	MemoryAvailableMB  *int    `json:"memory_available_mb"`
	BatteryLevel       *int    `json:"battery_level"`
	StorageAvailableMB *int    `json:"storage_available_mb"`
	NetworkType        *string `json:"network_type"`
}

// ============================================================
// DEVICE ANALYTIC
// ============================================================

// DeviceAnalytic represents device analytics event
type DeviceAnalytic struct {
	ID                    string     `json:"id" db:"id"`
	TenantID              string     `json:"tenant_id" db:"tenant_id"`
	AppID                 string     `json:"app_id" db:"app_id"`
	DeviceID              string     `json:"device_id" db:"device_id"`
	UserID                int64      `json:"user_id" db:"user_id"`
	EventType             string     `json:"event_type" db:"event_type"`
	EventName             *string    `json:"event_name" db:"event_name"`
	EventCategory         *string    `json:"event_category" db:"event_category"`
	EventAction           *string    `json:"event_action" db:"event_action"`
	EventLabel            *string    `json:"event_label" db:"event_label"`
	EventValue            *float64   `json:"event_value" db:"event_value"`
	ScreenName            *string    `json:"screen_name" db:"screen_name"`
	ScreenClass           *string    `json:"screen_class" db:"screen_class"`
	PageTitle             *string    `json:"page_title" db:"page_title"`
	PagePath              *string    `json:"page_path" db:"page_path"`
	SessionID             *string    `json:"session_id" db:"session_id"`
	EventID               string     `json:"event_id" db:"event_id"`
	EventTimestamp        time.Time  `json:"event_timestamp" db:"event_timestamp"`
	SessionStartTime      *time.Time `json:"session_start_time" db:"session_start_time"`
	TimeOnScreenSeconds   *int       `json:"time_on_screen_seconds" db:"time_on_screen_seconds"`
	EngagementTimeSeconds *int       `json:"engagement_time_seconds" db:"engagement_time_seconds"`
	ScrollDepthPercent    *int       `json:"scroll_depth_percent" db:"scroll_depth_percent"`
	ClickCount            *int       `json:"click_count" db:"click_count"`
	FormCompletionPercent *int       `json:"form_completion_percent" db:"form_completion_percent"`
	ErrorsEncountered     *int       `json:"errors_encountered" db:"errors_encountered"`
	Referrer              *string    `json:"referrer" db:"referrer"`
	UTMSource             *string    `json:"utm_source" db:"utm_source"`
	UTMMedium             *string    `json:"utm_medium" db:"utm_medium"`
	UTMCampaign           *string    `json:"utm_campaign" db:"utm_campaign"`
	UTMContent            *string    `json:"utm_content" db:"utm_content"`
	UTMTerm               *string    `json:"utm_term" db:"utm_term"`
	CustomParams          *string    `json:"custom_params" db:"custom_params"`     // JSON
	DeviceInfo            *string    `json:"device_info" db:"device_info"`         // JSON
	PerformanceMetrics    *string    `json:"performance_metrics" db:"performance_metrics"` // JSON
	Metadata              *string    `json:"metadata" db:"metadata"`               // JSON
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at" db:"deleted_at"`
}

// TrackAnalyticsRequest is the DTO for tracking analytics
type TrackAnalyticsRequest struct {
	EventType   string                 `json:"event_type" binding:"required"`
	EventName   *string                `json:"event_name"`
	EventValue  *float64               `json:"event_value"`
	ScreenName  *string                `json:"screen_name"`
	EventParams map[string]interface{} `json:"event_params"`
}

// ============================================================
// APP UPDATE
// ============================================================

// AppUpdate represents an app version release
type AppUpdate struct {
	ID                      string     `json:"id" db:"id"`
	TenantID                string     `json:"tenant_id" db:"tenant_id"`
	AppID                   string     `json:"app_id" db:"app_id"`
	Version                 string     `json:"version" db:"version"`
	BuildNumber             int        `json:"build_number" db:"build_number"`
	ReleaseDate             time.Time  `json:"release_date" db:"release_date"`
	UpdateType              string     `json:"update_type" db:"update_type"`
	IsMandatory             bool       `json:"is_mandatory" db:"is_mandatory"`
	IsRollbackAvailable     bool       `json:"is_rollback_available" db:"is_rollback_available"`
	RollbackVersion         *string    `json:"rollback_version" db:"rollback_version"`
	UpdateTitle             *string    `json:"update_title" db:"update_title"`
	UpdateDescription       *string    `json:"update_description" db:"update_description"`
	Changelog               *string    `json:"changelog" db:"changelog"`
	ReleaseNotes            *string    `json:"release_notes" db:"release_notes"`
	UpdateSizeMB            *float64   `json:"update_size_mb" db:"update_size_mb"`
	MinOSVersion            *string    `json:"min_os_version" db:"min_os_version"`
	MaxOSVersion            *string    `json:"max_os_version" db:"max_os_version"`
	DownloadURLIOS          *string    `json:"download_url_ios" db:"download_url_ios"`
	DownloadURLAndroid      *string    `json:"download_url_android" db:"download_url_android"`
	DownloadURLWindows      *string    `json:"download_url_windows" db:"download_url_windows"`
	ChecksumSHA256          *string    `json:"checksum_sha256" db:"checksum_sha256"`
	InstallInstructions     *string    `json:"install_instructions" db:"install_instructions"`
	BreakingChanges         *string    `json:"breaking_changes" db:"breaking_changes"`
	DeprecatedFeatures      *string    `json:"deprecated_features" db:"deprecated_features"` // JSON
	NewFeatures             *string    `json:"new_features" db:"new_features"`               // JSON
	BugFixes                *string    `json:"bug_fixes" db:"bug_fixes"`                     // JSON
	PerformanceImprovements *string    `json:"performance_improvements" db:"performance_improvements"` // JSON
	SecurityUpdates         *string    `json:"security_updates" db:"security_updates"`       // JSON
	UpdateStage             string     `json:"update_stage" db:"update_stage"`
	RolloutPercentage       int        `json:"rollout_percentage" db:"rollout_percentage"`
	RolloutStartTime        *time.Time `json:"rollout_start_time" db:"rollout_start_time"`
	RolloutEndTime          *time.Time `json:"rollout_end_time" db:"rollout_end_time"`
	InstallCount            int        `json:"install_count" db:"install_count"`
	UninstallCount          int        `json:"uninstall_count" db:"uninstall_count"`
	RollbackCount           int        `json:"rollback_count" db:"rollback_count"`
	UpdateStatus            string     `json:"update_status" db:"update_status"`
	DeprecationReason       *string    `json:"deprecation_reason" db:"deprecation_reason"`
	Metadata                *string    `json:"metadata" db:"metadata"` // JSON
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt               *time.Time `json:"deleted_at" db:"deleted_at"`
}

// CreateUpdateRequest is the DTO for creating an update
type CreateUpdateRequest struct {
	Version              string  `json:"version" binding:"required"`
	BuildNumber          int     `json:"build_number" binding:"required"`
	UpdateType           string  `json:"update_type" binding:"required"`
	IsMandatory          *bool   `json:"is_mandatory"`
	UpdateTitle          *string `json:"update_title"`
	UpdateDescription    *string `json:"update_description"`
	Changelog            *string `json:"changelog"`
	UpdateSizeMB         *float64 `json:"update_size_mb"`
	MinOSVersion         *string `json:"min_os_version"`
}

// ============================================================
// DEVICE UPDATE HISTORY
// ============================================================

// DeviceUpdateHistory represents update history for a device
type DeviceUpdateHistory struct {
	ID                  string     `json:"id" db:"id"`
	TenantID            string     `json:"tenant_id" db:"tenant_id"`
	DeviceID            string     `json:"device_id" db:"device_id"`
	AppID               string     `json:"app_id" db:"app_id"`
	UserID              int64      `json:"user_id" db:"user_id"`
	UpdateID            string     `json:"update_id" db:"update_id"`
	PreviousVersion     *string    `json:"previous_version" db:"previous_version"`
	NewVersion          string     `json:"new_version" db:"new_version"`
	UpdateStartTime     *time.Time `json:"update_start_time" db:"update_start_time"`
	UpdateCompletionTime *time.Time `json:"update_completion_time" db:"update_completion_time"`
	UpdateDurationSeconds *int     `json:"update_duration_seconds" db:"update_duration_seconds"`
	UpdateMethod        *string    `json:"update_method" db:"update_method"`
	UpdateStatus        string     `json:"update_status" db:"update_status"`
	FailureReason       *string    `json:"failure_reason" db:"failure_reason"`
	InstallationLog     *string    `json:"installation_log" db:"installation_log"`
	DeviceRestarted     bool       `json:"device_restarted" db:"device_restarted"`
	RestartTime         *time.Time `json:"restart_time" db:"restart_time"`
	PostInstallErrors   *string    `json:"post_install_errors" db:"post_install_errors"` // JSON
	UserFeedback        *string    `json:"user_feedback" db:"user_feedback"`
	UserRating          *int       `json:"user_rating" db:"user_rating"`
	Metadata            *string    `json:"metadata" db:"metadata"` // JSON
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at" db:"deleted_at"`
}

// UpdateStatusRequest is the DTO for updating update status
type UpdateStatusRequest struct {
	UpdateStatus        string  `json:"update_status" binding:"required"`
	UpdateDurationSeconds *int  `json:"update_duration_seconds"`
	FailureReason       *string `json:"failure_reason"`
	UserFeedback        *string `json:"user_feedback"`
	UserRating          *int    `json:"user_rating"`
}
