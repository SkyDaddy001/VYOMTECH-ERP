package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/logger"
)

// AuditService handles audit logging and compliance tracking
type AuditService struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewAuditService creates a new audit service
func NewAuditService(db *sql.DB, log *logger.Logger) *AuditService {
	return &AuditService{
		db:     db,
		logger: log,
	}
}

// LogAction logs an action for audit trail
func (as *AuditService) LogAction(ctx context.Context, log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (tenant_id, user_id, action, resource, details, ip_address, user_agent, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := as.db.ExecContext(ctx, query,
		log.TenantID, log.UserID, log.Action, log.Resource, log.Details,
		log.IPAddress, log.UserAgent, log.Status, time.Now())

	if err != nil {
		as.logger.Error("Failed to log action", "error", err, "action", log.Action)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	log.ID = id
	log.CreatedAt = time.Now()

	return nil
}

// GetAuditLogs retrieves audit logs with optional filters
func (as *AuditService) GetAuditLogs(ctx context.Context, tenantID string, filters map[string]interface{}, limit, offset int) ([]models.AuditLog, error) {
	query := `
		SELECT id, tenant_id, user_id, action, resource, details, ip_address, user_agent, status, created_at
		FROM audit_logs
		WHERE tenant_id = ?
	`

	args := []interface{}{tenantID}
	paramIndex := 2

	// Apply filters
	if userID, ok := filters["user_id"]; ok {
		query += ` AND user_id = $` + strconv.Itoa(paramIndex)
		args = append(args, userID)
		paramIndex++
	}

	if action, ok := filters["action"]; ok {
		query += ` AND action = $` + strconv.Itoa(paramIndex)
		args = append(args, action)
		paramIndex++
	}

	if resource, ok := filters["resource"]; ok {
		query += ` AND resource = $` + strconv.Itoa(paramIndex)
		args = append(args, resource)
		paramIndex++
	}

	if status, ok := filters["status"]; ok {
		query += ` AND status = $` + strconv.Itoa(paramIndex)
		args = append(args, status)
		paramIndex++
	}

	if startDate, ok := filters["start_date"]; ok {
		query += ` AND created_at >= $` + strconv.Itoa(paramIndex)
		args = append(args, startDate)
		paramIndex++
	}

	if endDate, ok := filters["end_date"]; ok {
		query += ` AND created_at <= $` + strconv.Itoa(paramIndex)
		args = append(args, endDate)
		paramIndex++
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(paramIndex) + ` OFFSET $` + strconv.Itoa(paramIndex+1)
	args = append(args, limit, offset)

	rows, err := as.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		err := rows.Scan(&log.ID, &log.TenantID, &log.UserID, &log.Action, &log.Resource,
			&log.Details, &log.IPAddress, &log.UserAgent, &log.Status, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// GetAuditLogsByUser retrieves all audit logs for a specific user
func (as *AuditService) GetAuditLogsByUser(ctx context.Context, tenantID string, userID int64, limit, offset int) ([]models.AuditLog, error) {
	return as.GetAuditLogs(ctx, tenantID, map[string]interface{}{"user_id": userID}, limit, offset)
}

// GetAuditLogsByResource retrieves all audit logs for a specific resource
func (as *AuditService) GetAuditLogsByResource(ctx context.Context, tenantID, resource string, limit, offset int) ([]models.AuditLog, error) {
	return as.GetAuditLogs(ctx, tenantID, map[string]interface{}{"resource": resource}, limit, offset)
}

// GetAuditLogsByDateRange retrieves audit logs within a date range
func (as *AuditService) GetAuditLogsByDateRange(ctx context.Context, tenantID string, startDate, endDate time.Time, limit, offset int) ([]models.AuditLog, error) {
	return as.GetAuditLogs(ctx, tenantID, map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDate,
	}, limit, offset)
}

// GetAuditSummary generates a summary of audit activities
func (as *AuditService) GetAuditSummary(ctx context.Context, tenantID string, days int) (map[string]interface{}, error) {
	query := `
		SELECT
			action,
			COUNT(*) as count,
			SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success_count,
			SUM(CASE WHEN status = 'failure' THEN 1 ELSE 0 END) as failure_count
		FROM audit_logs
		WHERE tenant_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY action
		ORDER BY count DESC
	`

	rows, err := as.db.QueryContext(ctx, query, tenantID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := make(map[string]interface{})
	actionStats := make([]map[string]interface{}, 0)

	totalActions := 0
	totalSuccesses := 0
	totalFailures := 0

	for rows.Next() {
		var action string
		var count, successCount, failureCount int

		err := rows.Scan(&action, &count, &successCount, &failureCount)
		if err != nil {
			return nil, err
		}

		actionStats = append(actionStats, map[string]interface{}{
			"action":        action,
			"count":         count,
			"success_count": successCount,
			"failure_count": failureCount,
			"success_rate":  float64(successCount) / float64(count),
		})

		totalActions += count
		totalSuccesses += successCount
		totalFailures += failureCount
	}

	summary["total_actions"] = totalActions
	summary["success_count"] = totalSuccesses
	summary["failure_count"] = totalFailures
	summary["success_rate"] = float64(totalSuccesses) / float64(totalActions)
	summary["action_stats"] = actionStats
	summary["days"] = days

	return summary, rows.Err()
}

// LogSecurityEvent logs security-related events
func (as *AuditService) LogSecurityEvent(ctx context.Context, event *models.SecurityEvent) error {
	query := `
		INSERT INTO security_events (tenant_id, user_id, event_type, severity, description, ip_address, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := as.db.ExecContext(ctx, query,
		event.TenantID, event.UserID, event.EventType, event.Severity, event.Description,
		event.IPAddress, time.Now())

	if err != nil {
		as.logger.Error("Failed to log security event", "error", err, "event_type", event.EventType)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	event.ID = id
	event.CreatedAt = time.Now()

	return nil
}

// GetSecurityEvents retrieves security events with optional filters
func (as *AuditService) GetSecurityEvents(ctx context.Context, tenantID string, filters map[string]interface{}, limit, offset int) ([]models.SecurityEvent, error) {
	query := `
		SELECT id, tenant_id, user_id, event_type, severity, description, ip_address, resolved_at, created_at
		FROM security_events
		WHERE tenant_id = ?
	`

	args := []interface{}{tenantID}
	paramIndex := 2

	if eventType, ok := filters["event_type"]; ok {
		query += ` AND event_type = $` + strconv.Itoa(paramIndex)
		args = append(args, eventType)
		paramIndex++
	}

	if severity, ok := filters["severity"]; ok {
		query += ` AND severity = $` + strconv.Itoa(paramIndex)
		args = append(args, severity)
		paramIndex++
	}

	if unresolved, ok := filters["unresolved"].(bool); ok && unresolved {
		query += ` AND resolved_at IS NULL`
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(paramIndex) + ` OFFSET $` + strconv.Itoa(paramIndex+1)
	args = append(args, limit, offset)

	rows, err := as.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.SecurityEvent
	for rows.Next() {
		var event models.SecurityEvent
		err := rows.Scan(&event.ID, &event.TenantID, &event.UserID, &event.EventType,
			&event.Severity, &event.Description, &event.IPAddress, &event.ResolvedAt, &event.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

// ResolveSecurityEvent marks a security event as resolved
func (as *AuditService) ResolveSecurityEvent(ctx context.Context, tenantID string, eventID int64) error {
	query := `
		UPDATE security_events
		SET resolved_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := as.db.ExecContext(ctx, query, time.Now(), eventID, tenantID)
	if err != nil {
		as.logger.Error("Failed to resolve security event", "error", err)
		return err
	}

	return nil
}

// ArchiveOldAuditLogs archives audit logs older than specified days
func (as *AuditService) ArchiveOldAuditLogs(ctx context.Context, tenantID string, retentionDays int) (int64, error) {
	query := `
		DELETE FROM audit_logs
		WHERE tenant_id = ? AND created_at < DATE_SUB(NOW(), INTERVAL ? DAY)
	`

	result, err := as.db.ExecContext(ctx, query, tenantID, retentionDays)
	if err != nil {
		as.logger.Error("Failed to archive audit logs", "error", err)
		return 0, err
	}

	return result.RowsAffected()
}

// GetComplianceReport generates a compliance report
func (as *AuditService) GetComplianceReport(ctx context.Context, tenantID string, startDate, endDate time.Time) (map[string]interface{}, error) {
	report := make(map[string]interface{})

	// Get audit log statistics
	auditSummary, err := as.GetAuditLogsByDateRange(ctx, tenantID, startDate, endDate, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	// Get unresolved security events
	securityEvents, err := as.GetSecurityEvents(ctx, tenantID, map[string]interface{}{"unresolved": true}, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get security events: %w", err)
	}

	// Count events by severity
	severityCount := make(map[string]int)
	for _, event := range securityEvents {
		severityCount[event.Severity]++
	}

	// Count audit actions
	actionCount := make(map[string]int)
	successCount := 0
	failureCount := 0

	for _, log := range auditSummary {
		actionCount[log.Action]++
		if log.Status == "success" {
			successCount++
		} else {
			failureCount++
		}
	}

	report["period"] = map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDate,
	}
	report["audit_logs"] = map[string]interface{}{
		"total":        len(auditSummary),
		"success":      successCount,
		"failure":      failureCount,
		"success_rate": float64(successCount) / float64(len(auditSummary)),
		"actions":      actionCount,
	}
	report["security_events"] = map[string]interface{}{
		"total":            len(securityEvents),
		"severity_summary": severityCount,
		"unresolved":       len(securityEvents),
	}

	return report, nil
}

// LogUserAction logs a specific user action with context
func (as *AuditService) LogUserAction(ctx context.Context, tenantID string, userID int64, action, resource string, details map[string]interface{}, ipAddress, userAgent, status string) error {
	detailsJSON, _ := json.Marshal(details)

	log := &models.AuditLog{
		TenantID:  tenantID,
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Details:   string(detailsJSON),
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Status:    status,
	}

	return as.LogAction(ctx, log)
}
