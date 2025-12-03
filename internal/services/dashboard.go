package services

import (
	"context"
	"database/sql"
	"time"

	"vyomtech-backend/pkg/logger"
)

type DashboardService struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewDashboardService(db *sql.DB, logger *logger.Logger) *DashboardService {
	return &DashboardService{
		db:     db,
		logger: logger,
	}
}

// ============================================================================
// ANALYTICS
// ============================================================================

type TenantAnalytics struct {
	TenantID         string    `json:"tenant_id"`
	TotalLeads       int       `json:"total_leads"`
	ActiveLeads      int       `json:"active_leads"`
	ConvertedLeads   int       `json:"converted_leads"`
	TotalUsers       int       `json:"total_users"`
	ActiveUsers      int       `json:"active_users"`
	AverageLeadScore float64   `json:"average_lead_score"`
	LastUpdated      time.Time `json:"last_updated"`
}

// GetTenantAnalytics returns analytics for a specific tenant
func (ds *DashboardService) GetTenantAnalytics(ctx context.Context, tenantID string) (*TenantAnalytics, error) {
	analytics := &TenantAnalytics{
		TenantID:    tenantID,
		LastUpdated: time.Now(),
	}

	// Get total leads
	err := ds.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM `lead` WHERE tenant_id = ?",
		tenantID,
	).Scan(&analytics.TotalLeads)
	if err != nil && err != sql.ErrNoRows {
		ds.logger.Error("error getting total leads", map[string]interface{}{"error": err})
		return nil, err
	}

	// Get active leads (status = 'new' or 'contacted')
	err = ds.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM `lead` WHERE tenant_id = ? AND status IN ('new', 'contacted')",
		tenantID,
	).Scan(&analytics.ActiveLeads)
	if err != nil && err != sql.ErrNoRows {
		ds.logger.Error("error getting active leads", map[string]interface{}{"error": err})
		return nil, err
	}

	// Get converted leads
	err = ds.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM `lead` WHERE tenant_id = ? AND status = 'converted'",
		tenantID,
	).Scan(&analytics.ConvertedLeads)
	if err != nil && err != sql.ErrNoRows {
		ds.logger.Error("error getting converted leads", map[string]interface{}{"error": err})
		return nil, err
	}

	// Get total users
	err = ds.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM `user` WHERE tenant_id = ?",
		tenantID,
	).Scan(&analytics.TotalUsers)
	if err != nil && err != sql.ErrNoRows {
		ds.logger.Error("error getting total users", map[string]interface{}{"error": err})
		return nil, err
	}

	// Get active users (status = 'active')
	err = ds.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM `user` WHERE tenant_id = ? AND status = 'active'",
		tenantID,
	).Scan(&analytics.ActiveUsers)
	if err != nil && err != sql.ErrNoRows {
		ds.logger.Error("error getting active users", map[string]interface{}{"error": err})
		return nil, err
	}

	// Get average lead score
	err = ds.db.QueryRowContext(ctx,
		"SELECT COALESCE(AVG(overall_score), 0) FROM lead_scores WHERE tenant_id = ?",
		tenantID,
	).Scan(&analytics.AverageLeadScore)
	if err != nil && err != sql.ErrNoRows {
		ds.logger.Error("error getting average lead score", map[string]interface{}{"error": err})
		return nil, err
	}

	return analytics, nil
}

// ============================================================================
// ACTIVITY LOGS
// ============================================================================

type ActivityLog struct {
	ID           int64     `json:"id"`
	TenantID     string    `json:"tenant_id"`
	UserID       int64     `json:"user_id"`
	Action       string    `json:"action"`
	EntityType   string    `json:"entity_type"`
	EntityID     string    `json:"entity_id"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetActivityLogs retrieves activity logs for a tenant
func (ds *DashboardService) GetActivityLogs(ctx context.Context, tenantID string, limit int, offset int) ([]ActivityLog, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}

	query := `
		SELECT id, tenant_id, user_id, action, entity_type, entity_id, status, error_message, created_at
		FROM audit_logs
		WHERE tenant_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := ds.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		ds.logger.Error("error querying activity logs", map[string]interface{}{"error": err})
		return nil, err
	}
	defer rows.Close()

	var logs []ActivityLog
	for rows.Next() {
		var log ActivityLog
		err := rows.Scan(
			&log.ID, &log.TenantID, &log.UserID, &log.Action,
			&log.EntityType, &log.EntityID, &log.Status, &log.ErrorMessage, &log.CreatedAt,
		)
		if err != nil {
			ds.logger.Error("error scanning activity log", map[string]interface{}{"error": err})
			continue
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// ============================================================================
// USER MANAGEMENT
// ============================================================================

type UserStats struct {
	TenantID      string    `json:"tenant_id"`
	UserID        int64     `json:"user_id"`
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	LeadsAssigned int       `json:"leads_assigned"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetTenantUsers retrieves all users in a tenant with stats
func (ds *DashboardService) GetTenantUsers(ctx context.Context, tenantID string) ([]UserStats, error) {
	query := `
		SELECT u.id, u.email, u.first_name, u.last_name, u.role, u.status, u.created_at,
		       COUNT(l.id) as leads_assigned
		FROM ` + "`user`" + ` u
		LEFT JOIN ` + "`lead`" + ` l ON u.id = l.assigned_agent_id AND l.tenant_id = u.tenant_id
		WHERE u.tenant_id = ?
		GROUP BY u.id, u.email, u.first_name, u.last_name, u.role, u.status, u.created_at
		ORDER BY u.created_at DESC
	`

	rows, err := ds.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		ds.logger.Error("error querying tenant users", map[string]interface{}{"error": err})
		return nil, err
	}
	defer rows.Close()

	var users []UserStats
	for rows.Next() {
		var user UserStats
		user.TenantID = tenantID

		err := rows.Scan(
			&user.UserID, &user.Email, &user.FirstName, &user.LastName,
			&user.Role, &user.Status, &user.CreatedAt, &user.LeadsAssigned,
		)
		if err != nil {
			ds.logger.Error("error scanning user stats", map[string]interface{}{"error": err})
			continue
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// ============================================================================
// USAGE TRACKING
// ============================================================================

type UsageMetrics struct {
	TenantID          string    `json:"tenant_id"`
	LeadsCreatedToday int       `json:"leads_created_today"`
	CallsToday        int       `json:"calls_today"`
	TasksCreatedToday int       `json:"tasks_created_today"`
	NotificationsSent int       `json:"notifications_sent"`
	ApiCallsToday     int       `json:"api_calls_today"`
	AverageResponseMs float64   `json:"average_response_ms"`
	LastUpdated       time.Time `json:"last_updated"`
}

// GetUsageMetrics returns usage metrics for a tenant today
func (ds *DashboardService) GetUsageMetrics(ctx context.Context, tenantID string) (*UsageMetrics, error) {
	metrics := &UsageMetrics{
		TenantID:    tenantID,
		LastUpdated: time.Now(),
	}

	today := time.Now().Format("2006-01-02")

	// Leads created today
	err := ds.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM `lead` WHERE tenant_id = ? AND DATE(created_at) = ?",
		tenantID, today,
	).Scan(&metrics.LeadsCreatedToday)
	if err != nil && err != sql.ErrNoRows {
		ds.logger.Error("error getting leads created today", map[string]interface{}{"error": err})
	}

	// Tasks created today
	err = ds.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM tasks WHERE tenant_id = ? AND DATE(created_at) = ?",
		tenantID, today,
	).Scan(&metrics.TasksCreatedToday)
	if err != nil && err != sql.ErrNoRows {
		ds.logger.Error("error getting tasks created today", map[string]interface{}{"error": err})
	}

	// Notifications sent today
	err = ds.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM notifications WHERE tenant_id = ? AND DATE(created_at) = ?",
		tenantID, today,
	).Scan(&metrics.NotificationsSent)
	if err != nil && err != sql.ErrNoRows {
		ds.logger.Error("error getting notifications sent", map[string]interface{}{"error": err})
	}

	return metrics, nil
}

// LogActivity logs an action to the audit trail
func (ds *DashboardService) LogActivity(ctx context.Context, tenantID string, userID int64, action string, entityType string, entityID int64, description string) error {
	query := `
		INSERT INTO audit_logs (tenant_id, user_id, action, entity_type, entity_id, description, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`

	_, err := ds.db.ExecContext(ctx, query, tenantID, userID, action, entityType, entityID, description)
	if err != nil {
		ds.logger.Error("error logging activity", map[string]interface{}{"error": err, "action": action})
		return err
	}

	return nil
}
