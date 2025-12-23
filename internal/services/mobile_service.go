package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"vyomtech-backend/internal/models"
)

// MobileService handles all mobile app-related operations
type MobileService struct {
	db *sql.DB
}

// NewMobileService creates a new MobileService
func NewMobileService(db *sql.DB) *MobileService {
	return &MobileService{
		db: db,
	}
}

// ============================================================
// MOBILE APP METHODS
// ============================================================

// CreateApp creates a new mobile app configuration
func (ms *MobileService) CreateApp(ctx context.Context, app *models.MobileApp) error {
	query := `
		INSERT INTO mobile_app (id, tenant_id, app_name, app_type, bundle_identifier, version, build_number, status, description, api_key, api_secret, min_supported_version, max_supported_version, store_url_ios, store_url_android, support_email, support_phone, privacy_policy_url, terms_of_service_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := ms.db.ExecContext(ctx, query,
		app.ID, app.TenantID, app.AppName, app.AppType, app.BundleIdentifier, app.Version, app.BuildNumber,
		app.Status, app.Description, app.APIKey, app.APISecret, app.MinSupportedVersion, app.MaxSupportedVersion,
		app.StoreURLIOS, app.StoreURLAndroid, app.SupportEmail, app.SupportPhone, app.PrivacyPolicyURL, app.TermsOfServiceURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create mobile app: %w", err)
	}

	return nil
}

// GetApp retrieves an app by ID
func (ms *MobileService) GetApp(ctx context.Context, appID, tenantID string) (*models.MobileApp, error) {
	query := `
		SELECT id, tenant_id, app_name, app_type, bundle_identifier, version, build_number, status, description, api_key, api_secret, min_supported_version, max_supported_version, store_url_ios, store_url_android, app_icon_url, app_banner_url, support_email, support_phone, support_chat_url, privacy_policy_url, terms_of_service_url, changelog_url, feature_flags, metadata, created_at, updated_at, deleted_at
		FROM mobile_app
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	app := &models.MobileApp{}
	err := ms.db.QueryRowContext(ctx, query, appID, tenantID).Scan(
		&app.ID, &app.TenantID, &app.AppName, &app.AppType, &app.BundleIdentifier, &app.Version, &app.BuildNumber,
		&app.Status, &app.Description, &app.APIKey, &app.APISecret, &app.MinSupportedVersion, &app.MaxSupportedVersion,
		&app.StoreURLIOS, &app.StoreURLAndroid, &app.AppIconURL, &app.AppBannerURL, &app.SupportEmail, &app.SupportPhone,
		&app.SupportChatURL, &app.PrivacyPolicyURL, &app.TermsOfServiceURL, &app.ChangelogURL, &app.FeatureFlags,
		&app.Metadata, &app.CreatedAt, &app.UpdatedAt, &app.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("app not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get app: %w", err)
	}

	return app, nil
}

// ListApps retrieves all apps for a tenant
func (ms *MobileService) ListApps(ctx context.Context, tenantID string, limit, offset int) ([]*models.MobileApp, error) {
	query := `
		SELECT id, tenant_id, app_name, app_type, bundle_identifier, version, build_number, status, description, api_key, api_secret, min_supported_version, max_supported_version, store_url_ios, store_url_android, app_icon_url, support_email, support_phone, created_at, updated_at, deleted_at
		FROM mobile_app
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := ms.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list apps: %w", err)
	}
	defer rows.Close()

	var apps []*models.MobileApp
	for rows.Next() {
		app := &models.MobileApp{}
		if err := rows.Scan(
			&app.ID, &app.TenantID, &app.AppName, &app.AppType, &app.BundleIdentifier, &app.Version, &app.BuildNumber,
			&app.Status, &app.Description, &app.APIKey, &app.APISecret, &app.MinSupportedVersion, &app.MaxSupportedVersion,
			&app.StoreURLIOS, &app.StoreURLAndroid, &app.AppIconURL, &app.SupportEmail, &app.SupportPhone,
			&app.CreatedAt, &app.UpdatedAt, &app.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan app: %w", err)
		}
		apps = append(apps, app)
	}

	return apps, nil
}

// UpdateApp updates an existing app
func (ms *MobileService) UpdateApp(ctx context.Context, appID, tenantID string, updates *models.UpdateMobileAppRequest) error {
	query := `UPDATE mobile_app SET `
	args := []interface{}{}
	hasUpdates := false

	if updates.AppName != nil {
		if hasUpdates {
			query += ", "
		}
		query += "app_name = ?"
		args = append(args, *updates.AppName)
		hasUpdates = true
	}

	if updates.Status != nil {
		if hasUpdates {
			query += ", "
		}
		query += "status = ?"
		args = append(args, *updates.Status)
		hasUpdates = true
	}

	if !hasUpdates {
		return nil
	}

	query += ", updated_at = NOW() WHERE id = ? AND tenant_id = ?"
	args = append(args, appID, tenantID)

	result, err := ms.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update app: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("app not found")
	}

	return nil
}

// ============================================================
// DEVICE MANAGEMENT METHODS
// ============================================================

// RegisterDevice registers a new mobile device
func (ms *MobileService) RegisterDevice(ctx context.Context, device *models.MobileDevice) error {
	query := `
		INSERT INTO mobile_device (id, tenant_id, user_id, app_id, device_id, device_name, device_model, device_manufacturer, os_type, os_version, app_version, app_build, push_token, device_uuid, imei, device_serial, screen_resolution, screen_density, locale, timezone, battery_optimization_enabled, biometric_enabled, biometric_type, fcm_token, apns_token, device_status, first_seen_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), NOW())
	`

	_, err := ms.db.ExecContext(ctx, query,
		device.ID, device.TenantID, device.UserID, device.AppID, device.DeviceID, device.DeviceName, device.DeviceModel,
		device.DeviceManufacturer, device.OSType, device.OSVersion, device.AppVersion, device.AppBuild, device.PushToken,
		device.DeviceUUID, device.IMEI, device.DeviceSerial, device.ScreenResolution, device.ScreenDensity, device.Locale,
		device.Timezone, device.BatteryOptimizationEnabled, device.BiometricEnabled, device.BiometricType, device.FCMToken,
		device.APNSToken, device.DeviceStatus,
	)
	if err != nil {
		return fmt.Errorf("failed to register device: %w", err)
	}

	return nil
}

// GetDevice retrieves a device by ID
func (ms *MobileService) GetDevice(ctx context.Context, deviceID, tenantID string) (*models.MobileDevice, error) {
	query := `
		SELECT id, tenant_id, user_id, app_id, device_id, device_name, device_model, device_manufacturer, os_type, os_version, app_version, app_build, push_token, push_token_updated_at, device_uuid, imei, device_serial, screen_resolution, screen_density, locale, timezone, battery_optimization_enabled, biometric_enabled, biometric_type, fcm_token, apns_token, device_status, last_activity_at, first_seen_at, last_seen_at, created_at, updated_at, deleted_at
		FROM mobile_device
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	device := &models.MobileDevice{}
	err := ms.db.QueryRowContext(ctx, query, deviceID, tenantID).Scan(
		&device.ID, &device.TenantID, &device.UserID, &device.AppID, &device.DeviceID, &device.DeviceName, &device.DeviceModel,
		&device.DeviceManufacturer, &device.OSType, &device.OSVersion, &device.AppVersion, &device.AppBuild, &device.PushToken,
		&device.PushTokenUpdatedAt, &device.DeviceUUID, &device.IMEI, &device.DeviceSerial, &device.ScreenResolution,
		&device.ScreenDensity, &device.Locale, &device.Timezone, &device.BatteryOptimizationEnabled, &device.BiometricEnabled,
		&device.BiometricType, &device.FCMToken, &device.APNSToken, &device.DeviceStatus, &device.LastActivityAt,
		&device.FirstSeenAt, &device.LastSeenAt, &device.CreatedAt, &device.UpdatedAt, &device.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("device not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	return device, nil
}

// ListUserDevices retrieves all devices for a user
func (ms *MobileService) ListUserDevices(ctx context.Context, userID int64, tenantID string) ([]*models.MobileDevice, error) {
	query := `
		SELECT id, tenant_id, user_id, app_id, device_id, device_name, device_model, os_type, os_version, app_version, push_token, device_status, last_activity_at, first_seen_at, created_at, updated_at
		FROM mobile_device
		WHERE user_id = ? AND tenant_id = ? AND deleted_at IS NULL
		ORDER BY last_activity_at DESC
	`

	rows, err := ms.db.QueryContext(ctx, query, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user devices: %w", err)
	}
	defer rows.Close()

	var devices []*models.MobileDevice
	for rows.Next() {
		device := &models.MobileDevice{}
		if err := rows.Scan(
			&device.ID, &device.TenantID, &device.UserID, &device.AppID, &device.DeviceID, &device.DeviceName, &device.DeviceModel,
			&device.OSType, &device.OSVersion, &device.AppVersion, &device.PushToken, &device.DeviceStatus,
			&device.LastActivityAt, &device.FirstSeenAt, &device.CreatedAt, &device.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// UpdateDevicePushToken updates device push token
func (ms *MobileService) UpdateDevicePushToken(ctx context.Context, deviceID, tenantID, pushToken string) error {
	query := `UPDATE mobile_device SET push_token = ?, push_token_updated_at = NOW(), updated_at = NOW() WHERE id = ? AND tenant_id = ?`

	result, err := ms.db.ExecContext(ctx, query, pushToken, deviceID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to update push token: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("device not found")
	}

	return nil
}

// ============================================================
// SESSION MANAGEMENT METHODS
// ============================================================

// CreateSession creates a new mobile session
func (ms *MobileService) CreateSession(ctx context.Context, session *models.MobileSession) error {
	query := `
		INSERT INTO mobile_session (id, tenant_id, user_id, device_id, app_id, session_token, refresh_token, token_expires_at, refresh_token_expires_at, session_status, ip_address, user_agent, app_version, login_method, login_timestamp, device_trusted, two_factor_verified, app_in_foreground, background_activity_allowed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := ms.db.ExecContext(ctx, query,
		session.ID, session.TenantID, session.UserID, session.DeviceID, session.AppID, session.SessionToken,
		session.RefreshToken, session.TokenExpiresAt, session.RefreshTokenExpiresAt, session.SessionStatus,
		session.IPAddress, session.UserAgent, session.AppVersion, session.LoginMethod,
		session.DeviceTrusted, session.TwoFactorVerified, session.AppInForeground, session.BackgroundActivityAllowed,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetSession retrieves a session by token
func (ms *MobileService) GetSession(ctx context.Context, sessionToken, tenantID string) (*models.MobileSession, error) {
	query := `
		SELECT id, tenant_id, user_id, device_id, app_id, session_token, refresh_token, token_expires_at, refresh_token_expires_at, session_status, ip_address, user_agent, app_version, login_method, login_timestamp, last_activity_timestamp, logout_timestamp, session_duration_seconds, device_trusted, two_factor_verified, network_type, app_in_foreground, created_at, updated_at, deleted_at
		FROM mobile_session
		WHERE session_token = ? AND tenant_id = ? AND session_status = 'active' AND deleted_at IS NULL
	`

	session := &models.MobileSession{}
	err := ms.db.QueryRowContext(ctx, query, sessionToken, tenantID).Scan(
		&session.ID, &session.TenantID, &session.UserID, &session.DeviceID, &session.AppID, &session.SessionToken,
		&session.RefreshToken, &session.TokenExpiresAt, &session.RefreshTokenExpiresAt, &session.SessionStatus,
		&session.IPAddress, &session.UserAgent, &session.AppVersion, &session.LoginMethod, &session.LoginTimestamp,
		&session.LastActivityTimestamp, &session.LogoutTimestamp, &session.SessionDurationSeconds, &session.DeviceTrusted,
		&session.TwoFactorVerified, &session.NetworkType, &session.AppInForeground, &session.CreatedAt, &session.UpdatedAt,
		&session.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// ============================================================
// NOTIFICATION METHODS
// ============================================================

// SendNotification sends a push notification
func (ms *MobileService) SendNotification(ctx context.Context, notification *models.AppNotification) error {
	query := `
		INSERT INTO app_notification (id, tenant_id, app_id, user_id, device_id, notification_type, notification_category, title, body, content_type, image_url, action_url, priority, custom_data, send_time, delivery_status, delivery_attempts, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, ?, NOW(), NOW())
	`

	_, err := ms.db.ExecContext(ctx, query,
		notification.ID, notification.TenantID, notification.AppID, notification.UserID, notification.DeviceID,
		notification.NotificationType, notification.NotificationCategory, notification.Title, notification.Body,
		notification.ContentType, notification.ImageURL, notification.ActionURL, notification.Priority,
		notification.CustomData, notification.DeliveryStatus, notification.DeliveryAttempts,
	)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}

// GetUserNotifications retrieves notifications for a user
func (ms *MobileService) GetUserNotifications(ctx context.Context, userID int64, tenantID string, limit, offset int) ([]*models.AppNotification, error) {
	query := `
		SELECT id, tenant_id, app_id, user_id, device_id, notification_type, title, body, delivery_status, read_at, clicked_at, created_at
		FROM app_notification
		WHERE user_id = ? AND tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := ms.db.QueryContext(ctx, query, userID, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*models.AppNotification
	for rows.Next() {
		notif := &models.AppNotification{}
		if err := rows.Scan(
			&notif.ID, &notif.TenantID, &notif.AppID, &notif.UserID, &notif.DeviceID, &notif.NotificationType,
			&notif.Title, &notif.Body, &notif.DeliveryStatus, &notif.ReadAt, &notif.ClickedAt, &notif.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notif)
	}

	return notifications, nil
}

// MarkNotificationAsRead marks a notification as read
func (ms *MobileService) MarkNotificationAsRead(ctx context.Context, notificationID, tenantID string) error {
	query := `UPDATE app_notification SET read_at = NOW(), delivery_status = 'read', updated_at = NOW() WHERE id = ? AND tenant_id = ?`

	result, err := ms.db.ExecContext(ctx, query, notificationID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// ============================================================
// FEATURE FLAG METHODS
// ============================================================

// GetAppFeatures retrieves enabled features for an app
func (ms *MobileService) GetAppFeatures(ctx context.Context, appID, tenantID string) ([]*models.AppFeature, error) {
	query := `
		SELECT id, tenant_id, app_id, feature_name, feature_code, feature_description, feature_category, is_enabled, rollout_percentage, config
		FROM app_feature
		WHERE app_id = ? AND tenant_id = ? AND is_enabled = 1 AND deleted_at IS NULL
		ORDER BY feature_name
	`

	rows, err := ms.db.QueryContext(ctx, query, appID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get features: %w", err)
	}
	defer rows.Close()

	var features []*models.AppFeature
	for rows.Next() {
		feature := &models.AppFeature{}
		if err := rows.Scan(
			&feature.ID, &feature.TenantID, &feature.AppID, &feature.FeatureName, &feature.FeatureCode,
			&feature.FeatureDescription, &feature.FeatureCategory, &feature.IsEnabled, &feature.RolloutPercentage,
			&feature.Config,
		); err != nil {
			return nil, fmt.Errorf("failed to scan feature: %w", err)
		}
		features = append(features, feature)
	}

	return features, nil
}

// IsFeatureEnabled checks if a feature is enabled for a user
func (ms *MobileService) IsFeatureEnabled(ctx context.Context, featureCode, appID, tenantID string, userID *int64) (bool, error) {
	query := `
		SELECT is_enabled, rollout_percentage, enabled_for_users, disabled_for_users
		FROM app_feature
		WHERE feature_code = ? AND app_id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	var isEnabled bool
	var rolloutPercent int
	var enabledForUsers sql.NullString
	var disabledForUsers sql.NullString

	err := ms.db.QueryRowContext(ctx, query, featureCode, appID, tenantID).Scan(&isEnabled, &rolloutPercent, &enabledForUsers, &disabledForUsers)
	if err == sql.ErrNoRows {
		return false, fmt.Errorf("feature not found")
	}
	if err != nil {
		return false, fmt.Errorf("failed to check feature: %w", err)
	}

	if !isEnabled {
		return false, nil
	}

	// Check if user is in disabled list
	if disabledForUsers.Valid && userID != nil {
		var disabled []int64
		if err := json.Unmarshal([]byte(disabledForUsers.String), &disabled); err == nil {
			for _, id := range disabled {
				if id == *userID {
					return false, nil
				}
			}
		}
	}

	// Check rollout percentage
	if rolloutPercent < 100 && userID != nil {
		userHash := int64(0)
		for _, b := range []byte(fmt.Sprintf("%d", *userID)) {
			userHash = (userHash * 31) + int64(b)
		}
		if (userHash % 100) >= int64(rolloutPercent) {
			return false, nil
		}
	}

	return true, nil
}

// ============================================================
// OFFLINE DATA METHODS
// ============================================================

// SaveOfflineData saves data for offline sync
func (ms *MobileService) SaveOfflineData(ctx context.Context, offlineData *models.OfflineData) error {
	query := `
		INSERT INTO offline_data (id, tenant_id, device_id, app_id, user_id, data_type, data_key, cached_data, data_hash, compression_type, sync_status, needs_sync, cache_priority, size_bytes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE cached_data = VALUES(cached_data), data_hash = VALUES(data_hash), sync_status = 'pending', needs_sync = 1, updated_at = NOW()
	`

	_, err := ms.db.ExecContext(ctx, query,
		offlineData.ID, offlineData.TenantID, offlineData.DeviceID, offlineData.AppID, offlineData.UserID,
		offlineData.DataType, offlineData.DataKey, offlineData.CachedData, offlineData.DataHash, offlineData.CompressionType,
		offlineData.SyncStatus, offlineData.NeedsSync, offlineData.CachePriority, offlineData.SizeBytes,
	)
	if err != nil {
		return fmt.Errorf("failed to save offline data: %w", err)
	}

	return nil
}

// GetOfflineData retrieves offline data for sync
func (ms *MobileService) GetOfflineData(ctx context.Context, deviceID, appID, tenantID string) ([]*models.OfflineData, error) {
	query := `
		SELECT id, tenant_id, device_id, app_id, user_id, data_type, data_key, data_hash, sync_status, needs_sync, size_bytes
		FROM offline_data
		WHERE device_id = ? AND app_id = ? AND tenant_id = ? AND needs_sync = 1 AND deleted_at IS NULL
		ORDER BY cache_priority DESC, created_at ASC
	`

	rows, err := ms.db.QueryContext(ctx, query, deviceID, appID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get offline data: %w", err)
	}
	defer rows.Close()

	var offlineDataList []*models.OfflineData
	for rows.Next() {
		od := &models.OfflineData{}
		if err := rows.Scan(
			&od.ID, &od.TenantID, &od.DeviceID, &od.AppID, &od.UserID, &od.DataType, &od.DataKey, &od.DataHash,
			&od.SyncStatus, &od.NeedsSync, &od.SizeBytes,
		); err != nil {
			return nil, fmt.Errorf("failed to scan offline data: %w", err)
		}
		offlineDataList = append(offlineDataList, od)
	}

	return offlineDataList, nil
}

// ============================================================
// CRASH REPORTING METHODS
// ============================================================

// ReportCrash reports an app crash
func (ms *MobileService) ReportCrash(ctx context.Context, crash *models.AppCrash) error {
	query := `
		INSERT INTO app_crash (id, tenant_id, app_id, device_id, user_id, crash_timestamp, app_version, app_build, os_version, device_model, crash_type, crash_reason, exception_type, exception_message, stack_trace, breadcrumbs, user_report, severity, memory_used_mb, memory_available_mb, battery_level, storage_available_mb, network_type, session_id, crash_status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := ms.db.ExecContext(ctx, query,
		crash.ID, crash.TenantID, crash.AppID, crash.DeviceID, crash.UserID, crash.CrashTimestamp, crash.AppVersion,
		crash.AppBuild, crash.OSVersion, crash.DeviceModel, crash.CrashType, crash.CrashReason, crash.ExceptionType,
		crash.ExceptionMessage, crash.StackTrace, crash.Breadcrumbs, crash.UserReport, crash.Severity, crash.MemoryUsedMB,
		crash.MemoryAvailableMB, crash.BatteryLevel, crash.StorageAvailableMB, crash.NetworkType, crash.SessionID, crash.CrashStatus,
	)
	if err != nil {
		return fmt.Errorf("failed to report crash: %w", err)
	}

	return nil
}

// GetCrashes retrieves recent crashes for an app
func (ms *MobileService) GetCrashes(ctx context.Context, appID, tenantID string, limit int) ([]*models.AppCrash, error) {
	query := `
		SELECT id, tenant_id, app_id, device_id, user_id, crash_timestamp, app_version, exception_type, exception_message, severity, crash_status
		FROM app_crash
		WHERE app_id = ? AND tenant_id = ? AND deleted_at IS NULL
		ORDER BY crash_timestamp DESC
		LIMIT ?
	`

	rows, err := ms.db.QueryContext(ctx, query, appID, tenantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get crashes: %w", err)
	}
	defer rows.Close()

	var crashes []*models.AppCrash
	for rows.Next() {
		crash := &models.AppCrash{}
		if err := rows.Scan(
			&crash.ID, &crash.TenantID, &crash.AppID, &crash.DeviceID, &crash.UserID, &crash.CrashTimestamp,
			&crash.AppVersion, &crash.ExceptionType, &crash.ExceptionMessage, &crash.Severity, &crash.CrashStatus,
		); err != nil {
			return nil, fmt.Errorf("failed to scan crash: %w", err)
		}
		crashes = append(crashes, crash)
	}

	return crashes, nil
}

// ============================================================
// ANALYTICS METHODS
// ============================================================

// TrackEvent tracks a user analytics event
func (ms *MobileService) TrackEvent(ctx context.Context, event *models.DeviceAnalytic) error {
	query := `
		INSERT INTO device_analytic (id, tenant_id, app_id, device_id, user_id, event_type, event_name, event_category, event_action, event_label, event_value, event_id, event_timestamp, screen_name, session_id, utm_source, utm_campaign, custom_params, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := ms.db.ExecContext(ctx, query,
		event.ID, event.TenantID, event.AppID, event.DeviceID, event.UserID, event.EventType, event.EventName,
		event.EventCategory, event.EventAction, event.EventLabel, event.EventValue, event.EventID, event.EventTimestamp,
		event.ScreenName, event.SessionID, event.UTMSource, event.UTMCampaign, event.CustomParams,
	)
	if err != nil {
		return fmt.Errorf("failed to track event: %w", err)
	}

	return nil
}

// ============================================================
// APP UPDATE METHODS
// ============================================================

// CreateUpdate creates a new app version update
func (ms *MobileService) CreateUpdate(ctx context.Context, update *models.AppUpdate) error {
	query := `
		INSERT INTO app_update (id, tenant_id, app_id, version, build_number, release_date, update_type, is_mandatory, update_title, update_description, changelog, update_size_mb, min_os_version, download_url_ios, download_url_android, checksum_sha256, update_stage, rollout_percentage, update_status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := ms.db.ExecContext(ctx, query,
		update.ID, update.TenantID, update.AppID, update.Version, update.BuildNumber, update.ReleaseDate,
		update.UpdateType, update.IsMandatory, update.UpdateTitle, update.UpdateDescription, update.Changelog,
		update.UpdateSizeMB, update.MinOSVersion, update.DownloadURLIOS, update.DownloadURLAndroid,
		update.ChecksumSHA256, update.UpdateStage, update.RolloutPercentage, update.UpdateStatus,
	)
	if err != nil {
		return fmt.Errorf("failed to create update: %w", err)
	}

	return nil
}

// GetLatestUpdate retrieves the latest app update
func (ms *MobileService) GetLatestUpdate(ctx context.Context, appID, tenantID string) (*models.AppUpdate, error) {
	query := `
		SELECT id, tenant_id, app_id, version, build_number, release_date, update_type, is_mandatory, update_title, changelog, download_url_ios, download_url_android, checksum_sha256, update_stage, rollout_percentage
		FROM app_update
		WHERE app_id = ? AND tenant_id = ? AND update_status = 'active' AND deleted_at IS NULL
		ORDER BY build_number DESC
		LIMIT 1
	`

	update := &models.AppUpdate{}
	err := ms.db.QueryRowContext(ctx, query, appID, tenantID).Scan(
		&update.ID, &update.TenantID, &update.AppID, &update.Version, &update.BuildNumber, &update.ReleaseDate,
		&update.UpdateType, &update.IsMandatory, &update.UpdateTitle, &update.Changelog, &update.DownloadURLIOS,
		&update.DownloadURLAndroid, &update.ChecksumSHA256, &update.UpdateStage, &update.RolloutPercentage,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no update available")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest update: %w", err)
	}

	return update, nil
}

// RecordDeviceUpdate records device update history
func (ms *MobileService) RecordDeviceUpdate(ctx context.Context, updateHistory *models.DeviceUpdateHistory) error {
	query := `
		INSERT INTO device_update_history (id, tenant_id, device_id, app_id, user_id, update_id, previous_version, new_version, update_start_time, update_status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, NOW(), NOW())
	`

	_, err := ms.db.ExecContext(ctx, query,
		updateHistory.ID, updateHistory.TenantID, updateHistory.DeviceID, updateHistory.AppID, updateHistory.UserID,
		updateHistory.UpdateID, updateHistory.PreviousVersion, updateHistory.NewVersion, updateHistory.UpdateStatus,
	)
	if err != nil {
		return fmt.Errorf("failed to record device update: %w", err)
	}

	return nil
}

// ============================================================
// SETTINGS METHODS
// ============================================================

// SaveUserSetting saves a user app setting
func (ms *MobileService) SaveUserSetting(ctx context.Context, setting *models.AppSetting) error {
	query := `
		INSERT INTO app_setting (id, tenant_id, user_id, app_id, device_id, setting_key, setting_value, setting_type, category, is_user_editable, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE setting_value = VALUES(setting_value), updated_at = NOW()
	`

	_, err := ms.db.ExecContext(ctx, query,
		setting.ID, setting.TenantID, setting.UserID, setting.AppID, setting.DeviceID, setting.SettingKey,
		setting.SettingValue, setting.SettingType, setting.Category, setting.IsUserEditable,
	)
	if err != nil {
		return fmt.Errorf("failed to save setting: %w", err)
	}

	return nil
}

// GetUserSettings retrieves user app settings
func (ms *MobileService) GetUserSettings(ctx context.Context, userID int64, appID, tenantID string) ([]*models.AppSetting, error) {
	query := `
		SELECT id, tenant_id, user_id, app_id, device_id, setting_key, setting_value, setting_type, category
		FROM app_setting
		WHERE user_id = ? AND app_id = ? AND tenant_id = ? AND deleted_at IS NULL
		ORDER BY category, setting_key
	`

	rows, err := ms.db.QueryContext(ctx, query, userID, appID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}
	defer rows.Close()

	var settings []*models.AppSetting
	for rows.Next() {
		setting := &models.AppSetting{}
		if err := rows.Scan(
			&setting.ID, &setting.TenantID, &setting.UserID, &setting.AppID, &setting.DeviceID, &setting.SettingKey,
			&setting.SettingValue, &setting.SettingType, &setting.Category,
		); err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings = append(settings, setting)
	}

	return settings, nil
}
