package services

import (
	"context"
	"database/sql"
	"fmt"

	"multi-tenant-ai-callcenter/internal/models"
)

// CampaignService handles all campaign-related operations
type CampaignService struct {
	db *sql.DB
}

// NewCampaignService creates a new CampaignService
func NewCampaignService(db *sql.DB) *CampaignService {
	return &CampaignService{
		db: db,
	}
}

// CreateCampaign creates a new campaign
func (cs *CampaignService) CreateCampaign(ctx context.Context, campaign *models.Campaign) error {
	query := `
		INSERT INTO campaign (tenant_id, name, description, status, target_leads, generated_leads, converted_leads,
		                       budget, spent_budget, cost_per_lead, conversion_rate, start_date, end_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := cs.db.ExecContext(ctx, query,
		campaign.TenantID, campaign.Name, campaign.Description, campaign.Status, campaign.TargetLeads,
		campaign.GeneratedLeads, campaign.ConvertedLeads, campaign.Budget, campaign.SpentBudget,
		campaign.CostPerLead, campaign.ConversionRate, campaign.StartDate, campaign.EndDate,
	)
	if err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get campaign id: %w", err)
	}

	campaign.ID = id
	return nil
}

// GetCampaign retrieves a campaign by ID
func (cs *CampaignService) GetCampaign(ctx context.Context, id int64, tenantID string) (*models.Campaign, error) {
	query := `
		SELECT id, tenant_id, name, description, status, target_leads, generated_leads, converted_leads,
		       budget, spent_budget, cost_per_lead, conversion_rate, start_date, end_date, created_at, updated_at
		FROM campaign
		WHERE id = ? AND tenant_id = ?
	`

	campaign := &models.Campaign{}
	err := cs.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&campaign.ID, &campaign.TenantID, &campaign.Name, &campaign.Description, &campaign.Status,
		&campaign.TargetLeads, &campaign.GeneratedLeads, &campaign.ConvertedLeads, &campaign.Budget,
		&campaign.SpentBudget, &campaign.CostPerLead, &campaign.ConversionRate, &campaign.StartDate,
		&campaign.EndDate, &campaign.CreatedAt, &campaign.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("campaign not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	return campaign, nil
}

// UpdateCampaign updates an existing campaign
func (cs *CampaignService) UpdateCampaign(ctx context.Context, campaign *models.Campaign) error {
	query := `
		UPDATE campaign
		SET name = ?, description = ?, status = ?, target_leads = ?, generated_leads = ?,
		    converted_leads = ?, budget = ?, spent_budget = ?, cost_per_lead = ?, conversion_rate = ?,
		    start_date = ?, end_date = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	result, err := cs.db.ExecContext(ctx, query,
		campaign.Name, campaign.Description, campaign.Status, campaign.TargetLeads,
		campaign.GeneratedLeads, campaign.ConvertedLeads, campaign.Budget, campaign.SpentBudget,
		campaign.CostPerLead, campaign.ConversionRate, campaign.StartDate, campaign.EndDate,
		campaign.ID, campaign.TenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("campaign not found")
	}

	return nil
}

// DeleteCampaign deletes a campaign
func (cs *CampaignService) DeleteCampaign(ctx context.Context, id int64, tenantID string) error {
	query := `DELETE FROM campaign WHERE id = ? AND tenant_id = ?`

	result, err := cs.db.ExecContext(ctx, query, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("campaign not found")
	}

	return nil
}

// GetCampaigns retrieves campaigns with filtering and pagination
func (cs *CampaignService) GetCampaigns(ctx context.Context, tenantID string, filter *models.CampaignFilter) ([]*models.Campaign, error) {
	query := `
		SELECT id, tenant_id, name, description, status, target_leads, generated_leads, converted_leads,
		       budget, spent_budget, cost_per_lead, conversion_rate, start_date, end_date, created_at, updated_at
		FROM campaign
		WHERE tenant_id = ?
	`

	args := []interface{}{tenantID}

	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, filter.Limit, filter.Offset)
	}

	rows, err := cs.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []*models.Campaign
	for rows.Next() {
		campaign := &models.Campaign{}
		err := rows.Scan(
			&campaign.ID, &campaign.TenantID, &campaign.Name, &campaign.Description, &campaign.Status,
			&campaign.TargetLeads, &campaign.GeneratedLeads, &campaign.ConvertedLeads, &campaign.Budget,
			&campaign.SpentBudget, &campaign.CostPerLead, &campaign.ConversionRate, &campaign.StartDate,
			&campaign.EndDate, &campaign.CreatedAt, &campaign.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign: %w", err)
		}
		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

// GetCampaignStats retrieves campaign statistics
func (cs *CampaignService) GetCampaignStats(ctx context.Context, tenantID string) (*models.CampaignStats, error) {
	query := `
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN status = 'paused' THEN 1 ELSE 0 END) as paused,
			AVG(conversion_rate) as avg_conversion,
			SUM(budget) as total_budget,
			SUM(spent_budget) as total_spent,
			AVG(cost_per_lead) as avg_cost_per_lead
		FROM campaign
		WHERE tenant_id = ?
	`

	stats := &models.CampaignStats{}
	err := cs.db.QueryRowContext(ctx, query, tenantID).Scan(
		&stats.Total, &stats.Active, &stats.Completed, &stats.Paused,
		&stats.AverageConversion, &stats.TotalBudget, &stats.TotalSpent, &stats.AverageCostPerLead,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign stats: %w", err)
	}

	return stats, nil
}
