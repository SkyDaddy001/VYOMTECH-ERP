package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// NotificationService defines notification management operations
type NotificationService interface {
	// Notification CRUD
	CreateNotification(ctx context.Context, tenantID string, notification *Notification) (*Notification, error)
	GetNotification(ctx context.Context, tenantID string, notificationID int64) (*Notification, error)
	GetUserNotifications(ctx context.Context, tenantID string, userID int64, limit int, offset int) ([]Notification, error)
	GetUnreadNotifications(ctx context.Context, tenantID string, userID int64) ([]Notification, error)
	MarkAsRead(ctx context.Context, tenantID string, notificationID int64) error
	MarkAllAsRead(ctx context.Context, tenantID string, userID int64) error
	ArchiveNotification(ctx context.Context, tenantID string, notificationID int64) error
	DeleteNotification(ctx context.Context, tenantID string, notificationID int64) error
	DeleteOldNotifications(ctx context.Context, tenantID string, days int) error

	// Notification preferences
	GetPreferences(ctx context.Context, tenantID string, userID int64) (*NotificationPreferences, error)
	UpdatePreferences(ctx context.Context, tenantID string, userID int64, prefs *NotificationPreferences) error

	// Statistics
	GetNotificationStats(ctx context.Context, tenantID string, userID int64) (*NotificationStats, error)
}

// Notification represents a notification
type Notification struct {
	ID                int64      `json:"id"`
	UserID            int64      `json:"user_id"`
	TenantID          string     `json:"tenant_id"`
	Type              string     `json:"type"` // lead_assigned, call_missed, deadline_reminder, task_completed, task_assigned
	Title             string     `json:"title"`
	Message           string     `json:"message"`
	Priority          string     `json:"priority"` // critical, high, normal, low
	Category          string     `json:"category"`
	RelatedEntityType *string    `json:"related_entity_type,omitempty"`
	RelatedEntityID   *int64     `json:"related_entity_id,omitempty"`
	IsRead            bool       `json:"is_read"`
	ReadAt            *time.Time `json:"read_at,omitempty"`
	ActionURL         *string    `json:"action_url,omitempty"`
	IsArchived        bool       `json:"is_archived"`
	ArchivedAt        *time.Time `json:"archived_at,omitempty"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// NotificationPreferences represents user notification preferences
type NotificationPreferences struct {
	ID                         int64     `json:"id"`
	UserID                     int64     `json:"user_id"`
	TenantID                   string    `json:"tenant_id"`
	EnableEmailNotifications   bool      `json:"enable_email_notifications"`
	EnableSmsNotifications     bool      `json:"enable_sms_notifications"`
	EnableInAppNotifications   bool      `json:"enable_in_app_notifications"`
	EnableDesktopNotifications bool      `json:"enable_desktop_notifications"`
	NotifyTaskAssigned         bool      `json:"notify_task_assigned"`
	NotifyTaskCompleted        bool      `json:"notify_task_completed"`
	NotifyDeadlineReminder     bool      `json:"notify_deadline_reminder"`
	NotifyLeadAssigned         bool      `json:"notify_lead_assigned"`
	NotifyCallMissed           bool      `json:"notify_call_missed"`
	NotifyMessageReceived      bool      `json:"notify_message_received"`
	NotifySystemAlerts         bool      `json:"notify_system_alerts"`
	QuietHoursEnabled          bool      `json:"quiet_hours_enabled"`
	QuietHoursStart            *string   `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd              *string   `json:"quiet_hours_end,omitempty"`
	EmailBatchFrequency        string    `json:"email_batch_frequency"`
	SmsBatchFrequency          string    `json:"sms_batch_frequency"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

// NotificationStats represents notification statistics
type NotificationStats struct {
	TotalNotifications int `json:"total_notifications"`
	UnreadCount        int `json:"unread_count"`
	CriticalCount      int `json:"critical_count"`
	HighPriorityCount  int `json:"high_priority_count"`
	TaskNotifications  int `json:"task_notifications"`
	LeadNotifications  int `json:"lead_notifications"`
	SystemAlerts       int `json:"system_alerts"`
}

// Notification type constants
const (
	NotificationTypeLeadAssigned     = "lead_assigned"
	NotificationTypeCallMissed       = "call_missed"
	NotificationTypeDeadlineReminder = "deadline_reminder"
	NotificationTypeTaskCompleted    = "task_completed"
	NotificationTypeTaskAssigned     = "task_assigned"
	NotificationTypeMessageReceived  = "message_received"
)

// Notification priority constants
const (
	NotificationPriorityCritical = "critical"
	NotificationPriorityHigh     = "high"
	NotificationPriorityNormal   = "normal"
	NotificationPriorityLow      = "low"
)

// notificationService implements NotificationService
type notificationService struct {
	db *sql.DB
}

// NewNotificationService creates a new notification service
func NewNotificationService(db *sql.DB) NotificationService {
	return &notificationService{db: db}
}

// CreateNotification creates a new notification
func (s *notificationService) CreateNotification(ctx context.Context, tenantID string, notification *Notification) (*Notification, error) {
	query := `
		INSERT INTO notifications (user_id, tenant_id, type, title, message, priority, category,
		                           related_entity_type, related_entity_id, action_url, expires_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		notification.UserID, tenantID, notification.Type, notification.Title, notification.Message,
		notification.Priority, notification.Category, notification.RelatedEntityType,
		notification.RelatedEntityID, notification.ActionURL, notification.ExpiresAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get notification ID: %w", err)
	}

	notification.ID = id
	notification.IsRead = false
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	return notification, nil
}

// GetNotification retrieves a notification by ID
func (s *notificationService) GetNotification(ctx context.Context, tenantID string, notificationID int64) (*Notification, error) {
	query := `
		SELECT id, user_id, tenant_id, type, title, message, priority, category,
		       related_entity_type, related_entity_id, is_read, read_at, action_url,
		       is_archived, archived_at, expires_at, created_at, updated_at
		FROM notifications
		WHERE id = ? AND tenant_id = ?
	`

	notification := &Notification{}
	err := s.db.QueryRowContext(ctx, query, notificationID, tenantID).Scan(
		&notification.ID, &notification.UserID, &notification.TenantID, &notification.Type,
		&notification.Title, &notification.Message, &notification.Priority, &notification.Category,
		&notification.RelatedEntityType, &notification.RelatedEntityID, &notification.IsRead,
		&notification.ReadAt, &notification.ActionURL, &notification.IsArchived,
		&notification.ArchivedAt, &notification.ExpiresAt, &notification.CreatedAt, &notification.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("notification not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	return notification, nil
}

// GetUserNotifications retrieves notifications for a user
func (s *notificationService) GetUserNotifications(ctx context.Context, tenantID string, userID int64, limit int, offset int) ([]Notification, error) {
	query := `
		SELECT id, user_id, tenant_id, type, title, message, priority, category,
		       related_entity_type, related_entity_id, is_read, read_at, action_url,
		       is_archived, archived_at, expires_at, created_at, updated_at
		FROM notifications
		WHERE user_id = ? AND tenant_id = ? AND is_archived = FALSE
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, userID, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		notification := Notification{}
		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.TenantID, &notification.Type,
			&notification.Title, &notification.Message, &notification.Priority, &notification.Category,
			&notification.RelatedEntityType, &notification.RelatedEntityID, &notification.IsRead,
			&notification.ReadAt, &notification.ActionURL, &notification.IsArchived,
			&notification.ArchivedAt, &notification.ExpiresAt, &notification.CreatedAt, &notification.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// GetUnreadNotifications retrieves unread notifications for a user
func (s *notificationService) GetUnreadNotifications(ctx context.Context, tenantID string, userID int64) ([]Notification, error) {
	query := `
		SELECT id, user_id, tenant_id, type, title, message, priority, category,
		       related_entity_type, related_entity_id, is_read, read_at, action_url,
		       is_archived, archived_at, expires_at, created_at, updated_at
		FROM notifications
		WHERE user_id = ? AND tenant_id = ? AND is_read = FALSE AND is_archived = FALSE
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread notifications: %w", err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		notification := Notification{}
		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.TenantID, &notification.Type,
			&notification.Title, &notification.Message, &notification.Priority, &notification.Category,
			&notification.RelatedEntityType, &notification.RelatedEntityID, &notification.IsRead,
			&notification.ReadAt, &notification.ActionURL, &notification.IsArchived,
			&notification.ArchivedAt, &notification.ExpiresAt, &notification.CreatedAt, &notification.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read
func (s *notificationService) MarkAsRead(ctx context.Context, tenantID string, notificationID int64) error {
	now := time.Now()
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query, now, notificationID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (s *notificationService) MarkAllAsRead(ctx context.Context, tenantID string, userID int64) error {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = NOW(), updated_at = NOW()
		WHERE user_id = ? AND tenant_id = ? AND is_read = FALSE
	`

	_, err := s.db.ExecContext(ctx, query, userID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

// ArchiveNotification archives a notification
func (s *notificationService) ArchiveNotification(ctx context.Context, tenantID string, notificationID int64) error {
	now := time.Now()
	query := `
		UPDATE notifications
		SET is_archived = TRUE, archived_at = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query, now, notificationID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to archive notification: %w", err)
	}

	return nil
}

// DeleteNotification deletes a notification
func (s *notificationService) DeleteNotification(ctx context.Context, tenantID string, notificationID int64) error {
	query := `DELETE FROM notifications WHERE id = ? AND tenant_id = ?`

	_, err := s.db.ExecContext(ctx, query, notificationID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}

// DeleteOldNotifications deletes notifications older than specified days
func (s *notificationService) DeleteOldNotifications(ctx context.Context, tenantID string, days int) error {
	query := `
		DELETE FROM notifications
		WHERE tenant_id = ? AND created_at < DATE_SUB(NOW(), INTERVAL ? DAY)
		AND (is_archived = TRUE OR is_read = TRUE)
	`

	_, err := s.db.ExecContext(ctx, query, tenantID, days)
	if err != nil {
		return fmt.Errorf("failed to delete old notifications: %w", err)
	}

	return nil
}

// GetPreferences retrieves notification preferences for a user
func (s *notificationService) GetPreferences(ctx context.Context, tenantID string, userID int64) (*NotificationPreferences, error) {
	query := `
		SELECT id, user_id, tenant_id, enable_email_notifications, enable_sms_notifications,
		       enable_in_app_notifications, enable_desktop_notifications, notify_task_assigned,
		       notify_task_completed, notify_deadline_reminder, notify_lead_assigned,
		       notify_call_missed, notify_message_received, notify_system_alerts,
		       quiet_hours_enabled, quiet_hours_start, quiet_hours_end,
		       email_batch_frequency, sms_batch_frequency, created_at, updated_at
		FROM notification_preferences
		WHERE user_id = ? AND tenant_id = ?
	`

	prefs := &NotificationPreferences{}
	err := s.db.QueryRowContext(ctx, query, userID, tenantID).Scan(
		&prefs.ID, &prefs.UserID, &prefs.TenantID, &prefs.EnableEmailNotifications,
		&prefs.EnableSmsNotifications, &prefs.EnableInAppNotifications, &prefs.EnableDesktopNotifications,
		&prefs.NotifyTaskAssigned, &prefs.NotifyTaskCompleted, &prefs.NotifyDeadlineReminder,
		&prefs.NotifyLeadAssigned, &prefs.NotifyCallMissed, &prefs.NotifyMessageReceived,
		&prefs.NotifySystemAlerts, &prefs.QuietHoursEnabled, &prefs.QuietHoursStart,
		&prefs.QuietHoursEnd, &prefs.EmailBatchFrequency, &prefs.SmsBatchFrequency,
		&prefs.CreatedAt, &prefs.UpdatedAt)

	if err == sql.ErrNoRows {
		// Create default preferences
		return s.createDefaultPreferences(ctx, tenantID, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	return prefs, nil
}

// UpdatePreferences updates notification preferences
func (s *notificationService) UpdatePreferences(ctx context.Context, tenantID string, userID int64, prefs *NotificationPreferences) error {
	query := `
		UPDATE notification_preferences
		SET enable_email_notifications = ?, enable_sms_notifications = ?,
		    enable_in_app_notifications = ?, enable_desktop_notifications = ?,
		    notify_task_assigned = ?, notify_task_completed = ?,
		    notify_deadline_reminder = ?, notify_lead_assigned = ?,
		    notify_call_missed = ?, notify_message_received = ?,
		    notify_system_alerts = ?, quiet_hours_enabled = ?,
		    quiet_hours_start = ?, quiet_hours_end = ?,
		    email_batch_frequency = ?, sms_batch_frequency = ?,
		    updated_at = NOW()
		WHERE user_id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		prefs.EnableEmailNotifications, prefs.EnableSmsNotifications,
		prefs.EnableInAppNotifications, prefs.EnableDesktopNotifications,
		prefs.NotifyTaskAssigned, prefs.NotifyTaskCompleted,
		prefs.NotifyDeadlineReminder, prefs.NotifyLeadAssigned,
		prefs.NotifyCallMissed, prefs.NotifyMessageReceived,
		prefs.NotifySystemAlerts, prefs.QuietHoursEnabled,
		prefs.QuietHoursStart, prefs.QuietHoursEnd,
		prefs.EmailBatchFrequency, prefs.SmsBatchFrequency,
		userID, tenantID)

	if err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	return nil
}

// GetNotificationStats retrieves notification statistics
func (s *notificationService) GetNotificationStats(ctx context.Context, tenantID string, userID int64) (*NotificationStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN is_read = FALSE THEN 1 ELSE 0 END) as unread,
			SUM(CASE WHEN priority = ? THEN 1 ELSE 0 END) as critical,
			SUM(CASE WHEN priority = ? THEN 1 ELSE 0 END) as high_priority,
			SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as task_notif,
			SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) as lead_notif,
			SUM(CASE WHEN category = 'system' THEN 1 ELSE 0 END) as system_alerts
		FROM notifications
		WHERE user_id = ? AND tenant_id = ? AND is_archived = FALSE
	`

	stats := &NotificationStats{}
	var total, unread, critical, highPriority, taskNotif, leadNotif, systemAlerts sql.NullInt64

	err := s.db.QueryRowContext(ctx, query,
		NotificationPriorityCritical, NotificationPriorityHigh,
		NotificationTypeTaskAssigned, NotificationTypeLeadAssigned,
		userID, tenantID).Scan(
		&total, &unread, &critical, &highPriority, &taskNotif, &leadNotif, &systemAlerts)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get notification stats: %w", err)
	}

	if total.Valid {
		stats.TotalNotifications = int(total.Int64)
	}
	if unread.Valid {
		stats.UnreadCount = int(unread.Int64)
	}
	if critical.Valid {
		stats.CriticalCount = int(critical.Int64)
	}
	if highPriority.Valid {
		stats.HighPriorityCount = int(highPriority.Int64)
	}
	if taskNotif.Valid {
		stats.TaskNotifications = int(taskNotif.Int64)
	}
	if leadNotif.Valid {
		stats.LeadNotifications = int(leadNotif.Int64)
	}
	if systemAlerts.Valid {
		stats.SystemAlerts = int(systemAlerts.Int64)
	}

	return stats, nil
}

// createDefaultPreferences creates default notification preferences for a user
func (s *notificationService) createDefaultPreferences(ctx context.Context, tenantID string, userID int64) (*NotificationPreferences, error) {
	query := `
		INSERT INTO notification_preferences (user_id, tenant_id, enable_email_notifications,
		                                      enable_sms_notifications, enable_in_app_notifications,
		                                      enable_desktop_notifications, notify_task_assigned,
		                                      notify_task_completed, notify_deadline_reminder,
		                                      notify_lead_assigned, notify_call_missed,
		                                      notify_message_received, notify_system_alerts,
		                                      email_batch_frequency, sms_batch_frequency)
		VALUES (?, ?, TRUE, FALSE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, 'immediate', 'immediate')
	`

	_, err := s.db.ExecContext(ctx, query, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to create default preferences: %w", err)
	}

	// Return the created preferences
	return s.GetPreferences(ctx, tenantID, userID)
}
