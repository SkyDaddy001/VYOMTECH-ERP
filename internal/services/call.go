package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"vyomtech-backend/internal/models"

	"github.com/google/uuid"
)

// CallService handles all call-related operations
type CallService struct {
	db *sql.DB
}

// NewCallService creates a new CallService
func NewCallService(db *sql.DB) *CallService {
	return &CallService{
		db: db,
	}
}

// CreateCall creates a new call record
func (cs *CallService) CreateCall(ctx context.Context, call *models.Call) error {
	// Generate UUID if not provided
	if call.ID == "" {
		call.ID = uuid.New().String()
	}

	query := `
		INSERT INTO call (id, tenant_id, lead_id, agent_id, status, duration_seconds, recording_url, notes, outcome, started_at, ended_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := cs.db.ExecContext(ctx, query,
		call.ID, call.TenantID, call.LeadID, call.AgentID, call.Status, call.DurationSeconds,
		call.RecordingURL, call.Notes, call.Outcome, call.StartedAt, call.EndedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create call: %w", err)
	}

	return nil
}

// GetCall retrieves a call by ID
func (cs *CallService) GetCall(ctx context.Context, id string, tenantID string) (*models.Call, error) {
	query := `
		SELECT id, tenant_id, lead_id, agent_id, status, duration_seconds, recording_url, notes, outcome, started_at, ended_at, created_at, updated_at
		FROM ` + "`call`" + `
		WHERE id = ? AND tenant_id = ?
	`

	call := &models.Call{}
	err := cs.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&call.ID, &call.TenantID, &call.LeadID, &call.AgentID, &call.Status, &call.DurationSeconds,
		&call.RecordingURL, &call.Notes, &call.Outcome, &call.StartedAt, &call.EndedAt, &call.CreatedAt, &call.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("call not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get call: %w", err)
	}

	return call, nil
}

// UpdateCall updates a call record
func (cs *CallService) UpdateCall(ctx context.Context, call *models.Call) error {
	query := `
		UPDATE call
		SET lead_id = ?, agent_id = ?, status = ?, duration_seconds = ?, recording_url = ?,
		    notes = ?, outcome = ?, started_at = ?, ended_at = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	result, err := cs.db.ExecContext(ctx, query,
		call.LeadID, call.AgentID, call.Status, call.DurationSeconds, call.RecordingURL,
		call.Notes, call.Outcome, call.StartedAt, call.EndedAt, call.ID, call.TenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update call: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("call not found")
	}

	return nil
}

// EndCall marks a call as ended
func (cs *CallService) EndCall(ctx context.Context, id string, tenantID string, outcome string, duration int) error {
	now := time.Now()
	query := `
		UPDATE call
		SET status = 'ended', outcome = ?, duration_seconds = ?, ended_at = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	result, err := cs.db.ExecContext(ctx, query, outcome, duration, now, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to end call: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("call not found")
	}

	return nil
}

// GetCalls retrieves calls with filtering and pagination
func (cs *CallService) GetCalls(ctx context.Context, tenantID string, filter *models.CallFilter) ([]*models.Call, error) {
	query := `
		SELECT id, tenant_id, lead_id, agent_id, call_status, duration_seconds, recording_url, notes, outcome, started_at, ended_at, created_at, updated_at
		FROM ` + "`call`" + `
		WHERE tenant_id = ?
	`

	args := []interface{}{tenantID}

	if filter.Status != "" {
		query += " AND `call_status` = ?"
		args = append(args, filter.Status)
	}
	if filter.Outcome != "" {
		query += " AND outcome = ?"
		args = append(args, filter.Outcome)
	}
	if filter.AgentID != "" {
		query += " AND agent_id = ?"
		args = append(args, filter.AgentID)
	}
	if filter.LeadID != "" {
		query += " AND lead_id = ?"
		args = append(args, filter.LeadID)
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, filter.Limit, filter.Offset)
	}

	rows, err := cs.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query calls: %w", err)
	}
	defer rows.Close()

	var calls []*models.Call
	for rows.Next() {
		call := &models.Call{}
		err := rows.Scan(
			&call.ID, &call.TenantID, &call.LeadID, &call.AgentID, &call.Status, &call.DurationSeconds,
			&call.RecordingURL, &call.Notes, &call.Outcome, &call.StartedAt, &call.EndedAt, &call.CreatedAt, &call.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan call: %w", err)
		}
		calls = append(calls, call)
	}

	return calls, nil
}

// GetCallStats retrieves call statistics
func (cs *CallService) GetCallStats(ctx context.Context, tenantID string) (*models.CallStats, error) {
	query := `
		SELECT
			COUNT(*) as total,
		SUM(CASE WHEN call_status = 'active' THEN 1 ELSE 0 END) as active,
		SUM(CASE WHEN call_status = 'ended' THEN 1 ELSE 0 END) as completed,
		SUM(CASE WHEN outcome = 'failed' THEN 1 ELSE 0 END) as failed,
		AVG(CASE WHEN duration_seconds > 0 THEN duration_seconds ELSE 0 END) as avg_duration,
		SUM(duration_seconds) as total_duration
		FROM ` + "`call`" + `
		WHERE tenant_id = ?
	`

	stats := &models.CallStats{}
	err := cs.db.QueryRowContext(ctx, query, tenantID).Scan(
		&stats.Total, &stats.Active, &stats.Completed, &stats.Failed, &stats.AverageDuration, &stats.TotalDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get call stats: %w", err)
	}

	if stats.Completed > 0 {
		stats.SuccessRate = float64(stats.Completed-stats.Failed) / float64(stats.Completed)
	}

	return stats, nil
}
