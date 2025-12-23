package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"vyomtech-backend/internal/models"
)

// SiteVisitService handles site visit operations using sql.DB
type SiteVisitService struct {
	db *sql.DB
}

// NewSiteVisitService creates a new SiteVisitService
func NewSiteVisitService(db *sql.DB) *SiteVisitService {
	return &SiteVisitService{db: db}
}

// ScheduleFilter represents filters for listing schedules
type ScheduleFilter struct {
	Status      *string
	StartDate   *time.Time
	EndDate     *time.Time
	LeadID      *string
	ScheduledBy *string
}

// VisitStats represents visit statistics for a user
type VisitStats struct {
	TotalScheduled int64 `json:"total_scheduled"`
	TotalCompleted int64 `json:"total_completed"`
	TotalCancelled int64 `json:"total_cancelled"`
	TotalNoShow    int64 `json:"total_no_show"`
}

// VisitSummary represents a summary of visits for a project
type VisitSummary struct {
	TotalVisits     int64      `json:"total_visits"`
	CompletedVisits int64      `json:"completed_visits"`
	UpcomingVisits  int64      `json:"upcoming_visits"`
	LastVisitDate   *time.Time `json:"last_visit_date"`
}

// ScheduleVisitRequest represents a request to schedule a visit
type ScheduleVisitRequest struct {
	TenantID      string    `json:"tenant_id"`
	LeadID        *string   `json:"lead_id"`
	VisitorName   string    `json:"visitor_name"`
	VisitorPhone  *string   `json:"visitor_phone"`
	VisitorEmail  *string   `json:"visitor_email"`
	ScheduledDate time.Time `json:"scheduled_date"`
	ScheduledBy   string    `json:"scheduled_by"`
	Status        string    `json:"status"`
}

// CreateVisitLogRequest represents a request to create a visit log
type CreateVisitLogRequest struct {
	TenantID         string     `json:"tenant_id"`
	VisitScheduleID  string     `json:"visit_schedule_id"`
	CheckInTime      *time.Time `json:"check_in_time"`
	CheckOutTime     *time.Time `json:"check_out_time"`
	VisitedBy        string     `json:"visited_by"`
	UnitsViewed      []string   `json:"units_viewed"`
	Feedback         *string    `json:"feedback"`
	FollowUpRequired bool       `json:"follow_up_required"`
	NextFollowupDate *time.Time `json:"next_followup_date"`
}

// UpdateVisitLogRequest represents a request to update a visit log
type UpdateVisitLogRequest struct {
	CheckInTime      *time.Time `json:"check_in_time"`
	CheckOutTime     *time.Time `json:"check_out_time"`
	UnitsViewed      []string   `json:"units_viewed"`
	Feedback         *string    `json:"feedback"`
	FollowUpRequired *bool      `json:"follow_up_required"`
	NextFollowupDate *time.Time `json:"next_followup_date"`
}

// CreateSchedule creates a new site visit schedule
func (s *SiteVisitService) CreateSchedule(ctx context.Context, req *ScheduleVisitRequest) (*models.SiteVisitSchedule, error) {
	if err := s.validateScheduleRequest(req); err != nil {
		return nil, err
	}

	scheduleID := generateUUID() // Assume you have a UUID generator

	query := `
		INSERT INTO site_visit_schedule (
			id, tenant_id, lead_id, visitor_name, visitor_phone, visitor_email,
			scheduled_date, scheduled_by, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := s.db.ExecContext(ctx, query,
		scheduleID, req.TenantID, req.LeadID, req.VisitorName, req.VisitorPhone,
		req.VisitorEmail, req.ScheduledDate, req.ScheduledBy, req.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return s.GetSchedule(ctx, scheduleID, req.TenantID)
}

// GetSchedule retrieves a site visit schedule by ID and tenant
func (s *SiteVisitService) GetSchedule(ctx context.Context, scheduleID, tenantID string) (*models.SiteVisitSchedule, error) {
	query := `
		SELECT id, tenant_id, lead_id, visitor_name, visitor_phone, visitor_email,
			   scheduled_date, scheduled_by, status, created_at, updated_at
		FROM site_visit_schedule
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	var schedule models.SiteVisitSchedule
	err := s.db.QueryRowContext(ctx, query, scheduleID, tenantID).Scan(
		&schedule.ID, &schedule.TenantID, &schedule.LeadID, &schedule.VisitorName,
		&schedule.VisitorPhone, &schedule.VisitorEmail, &schedule.ScheduledDate,
		&schedule.ScheduledBy, &schedule.Status, &schedule.CreatedAt, &schedule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("schedule not found")
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	return &schedule, nil
}

// ListSchedules lists site visit schedules with filters and pagination
func (s *SiteVisitService) ListSchedules(ctx context.Context, tenantID string, filter ScheduleFilter, limit, offset int) ([]models.SiteVisitSchedule, error) {
	query := `
		SELECT id, tenant_id, lead_id, visitor_name, visitor_phone, visitor_email,
			   scheduled_date, scheduled_by, status, created_at, updated_at
		FROM site_visit_schedule
		WHERE tenant_id = ? AND deleted_at IS NULL
	`
	args := []interface{}{tenantID}

	if filter.Status != nil {
		query += " AND status = ?"
		args = append(args, *filter.Status)
	}
	if filter.StartDate != nil {
		query += " AND scheduled_date >= ?"
		args = append(args, *filter.StartDate)
	}
	if filter.EndDate != nil {
		query += " AND scheduled_date <= ?"
		args = append(args, *filter.EndDate)
	}
	if filter.LeadID != nil {
		query += " AND lead_id = ?"
		args = append(args, *filter.LeadID)
	}
	if filter.ScheduledBy != nil {
		query += " AND scheduled_by = ?"
		args = append(args, *filter.ScheduledBy)
	}

	query += " ORDER BY scheduled_date DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list schedules: %w", err)
	}
	defer rows.Close()

	var schedules []models.SiteVisitSchedule
	for rows.Next() {
		var schedule models.SiteVisitSchedule
		err := rows.Scan(
			&schedule.ID, &schedule.TenantID, &schedule.LeadID, &schedule.VisitorName,
			&schedule.VisitorPhone, &schedule.VisitorEmail, &schedule.ScheduledDate,
			&schedule.ScheduledBy, &schedule.Status, &schedule.CreatedAt, &schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

// UpdateScheduleStatus updates the status of a site visit schedule
func (s *SiteVisitService) UpdateScheduleStatus(ctx context.Context, scheduleID, tenantID, status string) error {
	validStatuses := []string{"scheduled", "completed", "cancelled", "no_show"}
	valid := false
	for _, s := range validStatuses {
		if status == s {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("invalid status")
	}

	query := `
		UPDATE site_visit_schedule
		SET status = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, status, scheduleID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to update schedule status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("schedule not found")
	}

	return nil
}

// CancelSchedule cancels a site visit schedule
func (s *SiteVisitService) CancelSchedule(ctx context.Context, scheduleID, tenantID string) error {
	return s.UpdateScheduleStatus(ctx, scheduleID, tenantID, "cancelled")
}

// CreateVisitLog creates a new site visit log
func (s *SiteVisitService) CreateVisitLog(ctx context.Context, req *CreateVisitLogRequest) (*models.SiteVisitLog, error) {
	logID := generateUUID() // Assume you have a UUID generator

	unitsViewedJSON, err := json.Marshal(req.UnitsViewed)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal units viewed: %w", err)
	}

	query := `
		INSERT INTO site_visit_log (
			id, tenant_id, visit_schedule_id, check_in_time, check_out_time,
			visited_by, units_viewed, feedback, follow_up_required, next_followup_date,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err = s.db.ExecContext(ctx, query,
		logID, req.TenantID, req.VisitScheduleID, req.CheckInTime, req.CheckOutTime,
		req.VisitedBy, string(unitsViewedJSON), req.Feedback, req.FollowUpRequired, req.NextFollowupDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create visit log: %w", err)
	}

	return s.GetVisitLog(ctx, logID, req.TenantID)
}

// GetVisitLog retrieves a site visit log by ID and tenant
func (s *SiteVisitService) GetVisitLog(ctx context.Context, logID, tenantID string) (*models.SiteVisitLog, error) {
	query := `
		SELECT id, tenant_id, visit_schedule_id, check_in_time, check_out_time,
			   visited_by, units_viewed, feedback, follow_up_required, next_followup_date,
			   created_at, updated_at
		FROM site_visit_log
		WHERE id = ? AND tenant_id = ?
	`

	var log models.SiteVisitLog
	var unitsViewedJSON string
	err := s.db.QueryRowContext(ctx, query, logID, tenantID).Scan(
		&log.ID, &log.TenantID, &log.VisitScheduleID, &log.CheckInTime, &log.CheckOutTime,
		&log.VisitedBy, &unitsViewedJSON, &log.Feedback, &log.FollowUpRequired,
		&log.NextFollowupDate, &log.CreatedAt, &log.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("visit log not found")
		}
		return nil, fmt.Errorf("failed to get visit log: %w", err)
	}

	unitsViewed, err := s.parseUnitsViewed(unitsViewedJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse units viewed: %w", err)
	}
	uv := models.UnitsViewed(unitsViewed)
	log.UnitsViewed = &uv

	return &log, nil
}

// ListVisitLogs lists site visit logs for a schedule
func (s *SiteVisitService) ListVisitLogs(ctx context.Context, scheduleID, tenantID string) ([]models.SiteVisitLog, error) {
	query := `
		SELECT id, tenant_id, visit_schedule_id, check_in_time, check_out_time,
			   visited_by, units_viewed, feedback, follow_up_required, next_followup_date,
			   created_at, updated_at
		FROM site_visit_log
		WHERE visit_schedule_id = ? AND tenant_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, scheduleID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list visit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.SiteVisitLog
	for rows.Next() {
		var log models.SiteVisitLog
		var unitsViewedJSON string
		err := rows.Scan(
			&log.ID, &log.TenantID, &log.VisitScheduleID, &log.CheckInTime, &log.CheckOutTime,
			&log.VisitedBy, &unitsViewedJSON, &log.Feedback, &log.FollowUpRequired,
			&log.NextFollowupDate, &log.CreatedAt, &log.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan visit log: %w", err)
		}

		unitsViewed, err := s.parseUnitsViewed(unitsViewedJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to parse units viewed: %w", err)
		}
		uv := models.UnitsViewed(unitsViewed)
		log.UnitsViewed = &uv

		logs = append(logs, log)
	}

	return logs, nil
}

// UpdateVisitLog updates a site visit log
func (s *SiteVisitService) UpdateVisitLog(ctx context.Context, logID, tenantID string, req *UpdateVisitLogRequest) error {
	unitsViewedJSON, err := json.Marshal(req.UnitsViewed)
	if err != nil {
		return fmt.Errorf("failed to marshal units viewed: %w", err)
	}

	query := `
		UPDATE site_visit_log
		SET check_in_time = ?, check_out_time = ?, units_viewed = ?, feedback = ?,
			follow_up_required = ?, next_followup_date = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	result, err := s.db.ExecContext(ctx, query,
		req.CheckInTime, req.CheckOutTime, string(unitsViewedJSON), req.Feedback,
		req.FollowUpRequired, req.NextFollowupDate, logID, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update visit log: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("visit log not found")
	}

	return nil
}

// GetLeadVisitHistory gets the visit history for a lead
func (s *SiteVisitService) GetLeadVisitHistory(ctx context.Context, leadID, tenantID string) ([]models.SiteVisitSchedule, error) {
	query := `
		SELECT id, tenant_id, lead_id, visitor_name, visitor_phone, visitor_email,
			   scheduled_date, scheduled_by, status, created_at, updated_at
		FROM site_visit_schedule
		WHERE lead_id = ? AND tenant_id = ? AND deleted_at IS NULL
		ORDER BY scheduled_date DESC
	`

	rows, err := s.db.QueryContext(ctx, query, leadID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead visit history: %w", err)
	}
	defer rows.Close()

	var schedules []models.SiteVisitSchedule
	for rows.Next() {
		var schedule models.SiteVisitSchedule
		err := rows.Scan(
			&schedule.ID, &schedule.TenantID, &schedule.LeadID, &schedule.VisitorName,
			&schedule.VisitorPhone, &schedule.VisitorEmail, &schedule.ScheduledDate,
			&schedule.ScheduledBy, &schedule.Status, &schedule.CreatedAt, &schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

// GetUserVisitStats gets visit statistics for a user
func (s *SiteVisitService) GetUserVisitStats(ctx context.Context, userID, tenantID string, startDate, endDate time.Time) (*VisitStats, error) {
	query := `
		SELECT
			COUNT(CASE WHEN status = 'scheduled' THEN 1 END) as total_scheduled,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as total_completed,
			COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as total_cancelled,
			COUNT(CASE WHEN status = 'no_show' THEN 1 END) as total_no_show
		FROM site_visit_schedule
		WHERE scheduled_by = ? AND tenant_id = ? AND scheduled_date BETWEEN ? AND ?
			AND deleted_at IS NULL
	`

	var stats VisitStats
	err := s.db.QueryRowContext(ctx, query, userID, tenantID, startDate, endDate).Scan(
		&stats.TotalScheduled, &stats.TotalCompleted, &stats.TotalCancelled, &stats.TotalNoShow,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user visit stats: %w", err)
	}

	return &stats, nil
}

// GetProjectVisitSummary gets a summary of visits for a project (assuming project_id is in lead or schedule)
func (s *SiteVisitService) GetProjectVisitSummary(ctx context.Context, projectID, tenantID string) (*VisitSummary, error) {
	// This assumes there's a way to link schedules to projects, e.g., through leads
	// You may need to adjust this query based on your actual schema
	query := `
		SELECT
			COUNT(*) as total_visits,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_visits,
			COUNT(CASE WHEN status = 'scheduled' AND scheduled_date > NOW() THEN 1 END) as upcoming_visits,
			MAX(scheduled_date) as last_visit_date
		FROM site_visit_schedule
		WHERE tenant_id = ? AND deleted_at IS NULL
		-- Add project filtering logic here
	`

	var summary VisitSummary
	err := s.db.QueryRowContext(ctx, query, tenantID).Scan(
		&summary.TotalVisits, &summary.CompletedVisits, &summary.UpcomingVisits, &summary.LastVisitDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get project visit summary: %w", err)
	}

	return &summary, nil
}

// ListPendingFollowups lists visit logs that require follow-up
func (s *SiteVisitService) ListPendingFollowups(ctx context.Context, tenantID string) ([]models.SiteVisitLog, error) {
	query := `
		SELECT id, tenant_id, visit_schedule_id, check_in_time, check_out_time,
			   visited_by, units_viewed, feedback, follow_up_required, next_followup_date,
			   created_at, updated_at
		FROM site_visit_log
		WHERE tenant_id = ? AND follow_up_required = true
		ORDER BY next_followup_date ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list pending followups: %w", err)
	}
	defer rows.Close()

	var logs []models.SiteVisitLog
	for rows.Next() {
		var log models.SiteVisitLog
		var unitsViewedJSON string
		err := rows.Scan(
			&log.ID, &log.TenantID, &log.VisitScheduleID, &log.CheckInTime, &log.CheckOutTime,
			&log.VisitedBy, &unitsViewedJSON, &log.Feedback, &log.FollowUpRequired,
			&log.NextFollowupDate, &log.CreatedAt, &log.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan visit log: %w", err)
		}

		unitsViewed, err := s.parseUnitsViewed(unitsViewedJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to parse units viewed: %w", err)
		}
		uv := models.UnitsViewed(unitsViewed)
		log.UnitsViewed = &uv

		logs = append(logs, log)
	}

	return logs, nil
}

// validateScheduleRequest validates a schedule visit request
func (s *SiteVisitService) validateScheduleRequest(req *ScheduleVisitRequest) error {
	if req.TenantID == "" {
		return errors.New("tenant_id is required")
	}
	if req.VisitorName == "" {
		return errors.New("visitor_name is required")
	}
	if req.ScheduledBy == "" {
		return errors.New("scheduled_by is required")
	}
	if req.ScheduledDate.IsZero() {
		return errors.New("scheduled_date is required")
	}
	validStatuses := []string{"scheduled", "completed", "cancelled", "no_show"}
	valid := false
	for _, status := range validStatuses {
		if req.Status == status {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("invalid status")
	}
	return nil
}

// parseUnitsViewed parses the units viewed JSON string
func (s *SiteVisitService) parseUnitsViewed(unitsJSON string) ([]string, error) {
	var units []string
	if unitsJSON == "" || unitsJSON == "null" {
		return units, nil
	}
	err := json.Unmarshal([]byte(unitsJSON), &units)
	return units, err
}

// generateUUID generates a new UUID (placeholder - implement actual UUID generation)
func generateUUID() string {
	// Implement UUID generation logic here
	return "placeholder-uuid"
}
