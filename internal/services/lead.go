package services

import (
	"context"
	"database/sql"
	"fmt"

	"vyomtech-backend/internal/models"
)

// LeadService handles all lead-related operations
type LeadService struct {
	db *sql.DB
}

// NewLeadService creates a new LeadService
func NewLeadService(db *sql.DB) *LeadService {
	return &LeadService{
		db: db,
	}
}

// CreateLead creates a new lead
func (ls *LeadService) CreateLead(ctx context.Context, lead *models.Lead) error {
	query := `
		INSERT INTO lead (tenant_id, name, email, phone, status, source, campaign_id, assigned_agent_id, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := ls.db.ExecContext(ctx, query,
		lead.TenantID, lead.Name, lead.Email, lead.Phone, lead.Status, lead.Source,
		lead.CampaignID, lead.AssignedAgent, lead.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to create lead: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get lead id: %w", err)
	}

	lead.ID = id
	return nil
}

// GetLead retrieves a lead by ID
func (ls *LeadService) GetLead(ctx context.Context, id int64, tenantID string) (*models.Lead, error) {
	query := `
		SELECT id, tenant_id, name, email, phone, status, source, campaign_id, assigned_agent_id, notes, created_at, updated_at
		FROM lead
		WHERE id = ? AND tenant_id = ?
	`

	lead := &models.Lead{}
	err := ls.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&lead.ID, &lead.TenantID, &lead.Name, &lead.Email, &lead.Phone, &lead.Status,
		&lead.Source, &lead.CampaignID, &lead.AssignedAgent, &lead.Notes, &lead.CreatedAt, &lead.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("lead not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get lead: %w", err)
	}

	return lead, nil
}

// UpdateLead updates an existing lead
func (ls *LeadService) UpdateLead(ctx context.Context, lead *models.Lead) error {
	query := `
		UPDATE lead
		SET name = ?, email = ?, phone = ?, status = ?, source = ?, campaign_id = ?,
		    assigned_agent_id = ?, notes = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	result, err := ls.db.ExecContext(ctx, query,
		lead.Name, lead.Email, lead.Phone, lead.Status, lead.Source, lead.CampaignID,
		lead.AssignedAgent, lead.Notes, lead.ID, lead.TenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update lead: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("lead not found")
	}

	return nil
}

// DeleteLead deletes a lead
func (ls *LeadService) DeleteLead(ctx context.Context, id int64, tenantID string) error {
	query := `DELETE FROM lead WHERE id = ? AND tenant_id = ?`

	result, err := ls.db.ExecContext(ctx, query, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete lead: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("lead not found")
	}

	return nil
}

// GetLeads retrieves leads with filtering and pagination
func (ls *LeadService) GetLeads(ctx context.Context, tenantID string, filter *models.LeadFilter) ([]*models.Lead, error) {
	query := `
		SELECT id, tenant_id, name, email, phone, status, source, campaign_id, assigned_agent_id, notes, created_at, updated_at
		FROM lead
		WHERE tenant_id = ?
	`

	args := []interface{}{tenantID}

	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}
	if filter.Source != "" {
		query += " AND source = ?"
		args = append(args, filter.Source)
	}
	if filter.CampaignID > 0 {
		query += " AND campaign_id = ?"
		args = append(args, filter.CampaignID)
	}
	if filter.AssignedTo > 0 {
		query += " AND assigned_agent_id = ?"
		args = append(args, filter.AssignedTo)
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, filter.Limit, filter.Offset)
	}

	rows, err := ls.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query leads: %w", err)
	}
	defer rows.Close()

	var leads []*models.Lead
	for rows.Next() {
		lead := &models.Lead{}
		err := rows.Scan(
			&lead.ID, &lead.TenantID, &lead.Name, &lead.Email, &lead.Phone, &lead.Status,
			&lead.Source, &lead.CampaignID, &lead.AssignedAgent, &lead.Notes, &lead.CreatedAt, &lead.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lead: %w", err)
		}
		leads = append(leads, lead)
	}

	return leads, nil
}

// GetLeadStats retrieves statistics for leads
func (ls *LeadService) GetLeadStats(ctx context.Context, tenantID string) (*models.LeadStats, error) {
	query := `
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN status = 'new' THEN 1 ELSE 0 END) as new,
			SUM(CASE WHEN status = 'contacted' THEN 1 ELSE 0 END) as contacted,
			SUM(CASE WHEN status = 'qualified' THEN 1 ELSE 0 END) as qualified,
			SUM(CASE WHEN status = 'converted' THEN 1 ELSE 0 END) as converted,
			SUM(CASE WHEN status = 'lost' THEN 1 ELSE 0 END) as lost
		FROM lead
		WHERE tenant_id = ?
	`

	stats := &models.LeadStats{}
	err := ls.db.QueryRowContext(ctx, query, tenantID).Scan(
		&stats.Total, &stats.New, &stats.Contacted, &stats.Qualified, &stats.Converted, &stats.Lost,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead stats: %w", err)
	}

	if stats.Total > 0 {
		stats.ConvRate = float64(stats.Converted) / float64(stats.Total)
	}

	return stats, nil
}

// GetLeadsByPipelineStage retrieves leads by pipeline stage
func (ls *LeadService) GetLeadsByPipelineStage(ctx context.Context, tenantID string, stage string, filter *models.LeadFilter) ([]*models.Lead, error) {
	query := `
		SELECT id, tenant_id, name, email, phone, status, source, campaign_id, assigned_agent_id, notes, created_at, updated_at
		FROM lead
		WHERE tenant_id = ? AND (pipeline_stage = ? OR status IN (SELECT status FROM lead_pipeline_config WHERE tenant_id = ? AND pipeline_stage = ?))
	`

	args := []interface{}{tenantID, stage, tenantID, stage}

	if filter.AssignedTo > 0 {
		query += " AND assigned_agent_id = ?"
		args = append(args, filter.AssignedTo)
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, filter.Limit, filter.Offset)
	}

	rows, err := ls.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query leads by pipeline stage: %w", err)
	}
	defer rows.Close()

	var leads []*models.Lead
	for rows.Next() {
		lead := &models.Lead{}
		err := rows.Scan(
			&lead.ID, &lead.TenantID, &lead.Name, &lead.Email, &lead.Phone, &lead.Status,
			&lead.Source, &lead.CampaignID, &lead.AssignedAgent, &lead.Notes, &lead.CreatedAt, &lead.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lead: %w", err)
		}
		leads = append(leads, lead)
	}

	return leads, nil
}

// LogStatusChange logs a lead status change for audit purposes
func (ls *LeadService) LogStatusChange(ctx context.Context, leadID int64, tenantID string, oldStatus, newStatus string, userID *int64, reason string) error {
	query := `
		INSERT INTO lead_status_log (id, tenant_id, lead_id, old_status, new_status, old_pipeline_stage, new_pipeline_stage, changed_by, change_reason, created_at)
		VALUES (UUID(), ?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`

	oldStage := models.GetPipelineStage(oldStatus)
	newStage := models.GetPipelineStage(newStatus)

	_, err := ls.db.ExecContext(ctx, query, tenantID, leadID, oldStatus, newStatus, oldStage, newStage, userID, reason)
	if err != nil {
		return fmt.Errorf("failed to log status change: %w", err)
	}

	return nil
}

// GetLeadStatusHistory retrieves the history of status changes for a lead
func (ls *LeadService) GetLeadStatusHistory(ctx context.Context, leadID int64, tenantID string) ([]map[string]interface{}, error) {
	query := `
		SELECT id, old_status, new_status, old_pipeline_stage, new_pipeline_stage, changed_by, change_reason, created_at
		FROM lead_status_log
		WHERE lead_id = ? AND tenant_id = ?
		ORDER BY created_at DESC
		LIMIT 50
	`

	rows, err := ls.db.QueryContext(ctx, query, leadID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query status history: %w", err)
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var id, oldStatus, newStatus, oldStage, newStage, changeReason string
		var changedBy *int64
		var createdAt string

		err := rows.Scan(&id, &oldStatus, &newStatus, &oldStage, &newStage, &changedBy, &changeReason, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}

		entry := map[string]interface{}{
			"id":                 id,
			"old_status":         oldStatus,
			"new_status":         newStatus,
			"old_pipeline_stage": oldStage,
			"new_pipeline_stage": newStage,
			"changed_by":         changedBy,
			"change_reason":      changeReason,
			"created_at":         createdAt,
		}
		history = append(history, entry)
	}

	return history, nil
}
