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

// PartnerSourceService handles partner source tracking and mapping
type PartnerSourceService interface {
	// Source Management
	CreatePartnerSource(ctx context.Context, tenantID string, source *models.PartnerSource) (*models.PartnerSource, error)
	GetPartnerSource(ctx context.Context, tenantID string, sourceID int64) (*models.PartnerSource, error)
	GetPartnerSources(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerSource, error)
	GetSourceByCode(ctx context.Context, tenantID string, sourceCode string) (*models.PartnerSource, error)
	UpdatePartnerSource(ctx context.Context, tenantID string, source *models.PartnerSource) (*models.PartnerSource, error)
	DeactivatePartnerSource(ctx context.Context, tenantID string, sourceID int64) error

	// Statistics
	GetSourceStats(ctx context.Context, tenantID string, sourceID int64) (*models.PartnerSourceStats, error)
	GetPartnerSourceStats(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerSourceStats, error)
}

type partnerSourceService struct {
	db *sql.DB
}

// NewPartnerSourceService creates a new partner source service
func NewPartnerSourceService(db *sql.DB) PartnerSourceService {
	return &partnerSourceService{db: db}
}

// CreatePartnerSource creates a new partner source mapping
func (s *partnerSourceService) CreatePartnerSource(ctx context.Context, tenantID string, source *models.PartnerSource) (*models.PartnerSource, error) {
	if source.PartnerID == 0 || source.SourceType == "" {
		return nil, errors.New("partner_id and source_type are required")
	}

	source.TenantID = tenantID
	source.IsActive = true
	source.CreatedAt = time.Now()
	source.UpdatedAt = time.Now()

	query := `
		INSERT INTO partner_sources (
			tenant_id, partner_id, source_type, source_code, source_name,
			description, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		source.TenantID, source.PartnerID, source.SourceType, source.SourceCode, source.SourceName,
		source.Description, source.IsActive, source.CreatedAt, source.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create partner source: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get source ID: %w", err)
	}

	source.ID = id
	return source, nil
}

// GetPartnerSource retrieves a partner source
func (s *partnerSourceService) GetPartnerSource(ctx context.Context, tenantID string, sourceID int64) (*models.PartnerSource, error) {
	source := &models.PartnerSource{}

	query := `
		SELECT id, tenant_id, partner_id, source_type, source_code, source_name,
		description, is_active, leads_generated, leads_converted, total_revenue,
		created_at, updated_at, deleted_at
		FROM partner_sources
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := s.db.QueryRowContext(ctx, query, sourceID, tenantID).Scan(
		&source.ID, &source.TenantID, &source.PartnerID, &source.SourceType, &source.SourceCode, &source.SourceName,
		&source.Description, &source.IsActive, &source.LeadsGenerated, &source.LeadsConverted, &source.TotalRevenue,
		&source.CreatedAt, &source.UpdatedAt, &source.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("source not found")
		}
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	return source, nil
}

// GetPartnerSources retrieves all sources for a partner
func (s *partnerSourceService) GetPartnerSources(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerSource, error) {
	sources := []models.PartnerSource{}

	query := `
		SELECT id, tenant_id, partner_id, source_type, source_code, source_name,
		description, is_active, leads_generated, leads_converted, total_revenue,
		created_at, updated_at, deleted_at
		FROM partner_sources
		WHERE tenant_id = ? AND partner_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, partnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get partner sources: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var source models.PartnerSource
		err := rows.Scan(
			&source.ID, &source.TenantID, &source.PartnerID, &source.SourceType, &source.SourceCode, &source.SourceName,
			&source.Description, &source.IsActive, &source.LeadsGenerated, &source.LeadsConverted, &source.TotalRevenue,
			&source.CreatedAt, &source.UpdatedAt, &source.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan source: %w", err)
		}
		sources = append(sources, source)
	}

	return sources, rows.Err()
}

// GetSourceByCode retrieves a source by code
func (s *partnerSourceService) GetSourceByCode(ctx context.Context, tenantID string, sourceCode string) (*models.PartnerSource, error) {
	source := &models.PartnerSource{}

	query := `
		SELECT id, tenant_id, partner_id, source_type, source_code, source_name,
		description, is_active, leads_generated, leads_converted, total_revenue,
		created_at, updated_at, deleted_at
		FROM partner_sources
		WHERE tenant_id = ? AND source_code = ? AND deleted_at IS NULL
	`

	err := s.db.QueryRowContext(ctx, query, tenantID, sourceCode).Scan(
		&source.ID, &source.TenantID, &source.PartnerID, &source.SourceType, &source.SourceCode, &source.SourceName,
		&source.Description, &source.IsActive, &source.LeadsGenerated, &source.LeadsConverted, &source.TotalRevenue,
		&source.CreatedAt, &source.UpdatedAt, &source.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("source not found")
		}
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	return source, nil
}

// UpdatePartnerSource updates a partner source
func (s *partnerSourceService) UpdatePartnerSource(ctx context.Context, tenantID string, source *models.PartnerSource) (*models.PartnerSource, error) {
	source.UpdatedAt = time.Now()

	query := `
		UPDATE partner_sources SET
			source_code = ?, source_name = ?, description = ?,
			is_active = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		source.SourceCode, source.SourceName, source.Description,
		source.IsActive, source.UpdatedAt,
		source.ID, tenantID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update source: %w", err)
	}

	return source, nil
}

// DeactivatePartnerSource deactivates a source
func (s *partnerSourceService) DeactivatePartnerSource(ctx context.Context, tenantID string, sourceID int64) error {
	query := `
		UPDATE partner_sources SET is_active = FALSE, updated_at = ? WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query, time.Now(), sourceID, tenantID)
	return err
}

// GetSourceStats retrieves statistics for a source
func (s *partnerSourceService) GetSourceStats(ctx context.Context, tenantID string, sourceID int64) (*models.PartnerSourceStats, error) {
	stats := &models.PartnerSourceStats{}

	source, err := s.GetPartnerSource(ctx, tenantID, sourceID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			COUNT(*) as total_leads,
			SUM(CASE WHEN pl.status = 'approved' THEN 1 ELSE 0 END) as approved_leads,
			SUM(CASE WHEN pl.status = 'converted' THEN 1 ELSE 0 END) as converted_leads,
			AVG(pl.quality_score) as avg_quality,
			SUM(pl.credit_amount) as total_credits
		FROM partner_leads pl
		WHERE pl.tenant_id = ? AND pl.partner_id = ?
		AND EXISTS (
			SELECT 1 FROM partner_sources ps
			WHERE ps.id = ? AND ps.tenant_id = pl.tenant_id
		)
	`

	var totalLeads, approvedLeads, convertedLeads sql.NullInt64
	var avgQuality, totalCredits sql.NullFloat64

	err = s.db.QueryRowContext(ctx, query, tenantID, source.PartnerID, sourceID).Scan(
		&totalLeads, &approvedLeads, &convertedLeads, &avgQuality, &totalCredits,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get source stats: %w", err)
	}

	stats.SourceType = source.SourceType
	stats.SourceCode = source.SourceCode

	if totalLeads.Valid {
		stats.TotalLeads = totalLeads.Int64
	}
	if approvedLeads.Valid {
		stats.ApprovedLeads = approvedLeads.Int64
	}
	if convertedLeads.Valid {
		stats.ConvertedLeads = convertedLeads.Int64
	}
	if avgQuality.Valid {
		stats.AverageQualityScore = avgQuality.Float64
	}
	if totalCredits.Valid {
		stats.TotalCreditsAllocated = totalCredits.Float64
	}

	if stats.TotalLeads > 0 {
		stats.ApprovalRate = (float64(stats.ApprovedLeads) / float64(stats.TotalLeads)) * 100
		if stats.ApprovedLeads > 0 {
			stats.ConversionRate = (float64(stats.ConvertedLeads) / float64(stats.ApprovedLeads)) * 100
		}
	}

	return stats, nil
}

// GetPartnerSourceStats retrieves stats for all sources of a partner
func (s *partnerSourceService) GetPartnerSourceStats(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerSourceStats, error) {
	statsSlice := []models.PartnerSourceStats{}

	sources, err := s.GetPartnerSources(ctx, tenantID, partnerID)
	if err != nil {
		return nil, err
	}

	for _, source := range sources {
		stats, err := s.GetSourceStats(ctx, tenantID, source.ID)
		if err == nil {
			statsSlice = append(statsSlice, *stats)
		}
	}

	return statsSlice, nil
}

// PartnerCreditPolicyService handles credit policy management
type PartnerCreditPolicyService interface {
	// Policy Management
	CreateCreditPolicy(ctx context.Context, tenantID string, policy *models.PartnerCreditPolicy) (*models.PartnerCreditPolicy, error)
	GetCreditPolicy(ctx context.Context, tenantID string, policyID int64) (*models.PartnerCreditPolicy, error)
	GetPartnerCreditPolicies(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerCreditPolicy, error)
	GetActiveCreditPolicies(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerCreditPolicy, error)
	UpdateCreditPolicy(ctx context.Context, tenantID string, policy *models.PartnerCreditPolicy) (*models.PartnerCreditPolicy, error)
	ApproveCreditPolicy(ctx context.Context, tenantID string, policyID int64, approvedBy int64) error
	DeactivateCreditPolicy(ctx context.Context, tenantID string, policyID int64) error

	// Credit Calculation
	CalculateLeadCredit(ctx context.Context, tenantID string, partnerLeadID int64) (float64, []int64, error)
	GetApplicablePolicies(ctx context.Context, tenantID string, partnerID int64, leadData *models.LeadData) ([]models.PartnerCreditPolicy, error)

	// Policy Mapping
	GetPolicyMappings(ctx context.Context, tenantID string, partnerLeadID int64) ([]models.PartnerCreditPolicyMapping, error)
}

type partnerCreditPolicyService struct {
	db *sql.DB
}

// NewPartnerCreditPolicyService creates a new credit policy service
func NewPartnerCreditPolicyService(db *sql.DB) PartnerCreditPolicyService {
	return &partnerCreditPolicyService{db: db}
}

// CreateCreditPolicy creates a new credit policy
func (s *partnerCreditPolicyService) CreateCreditPolicy(ctx context.Context, tenantID string, policy *models.PartnerCreditPolicy) (*models.PartnerCreditPolicy, error) {
	if policy.PartnerID == 0 || policy.PolicyType == "" {
		return nil, errors.New("partner_id and policy_type are required")
	}

	policy.TenantID = tenantID
	policy.IsActive = true
	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()

	tierJSON, _ := json.Marshal(policy.TierConfig)

	query := `
		INSERT INTO partner_credit_policies (
			tenant_id, partner_id, policy_code, policy_name, policy_type,
			calculation_type, time_unit_type, time_unit_value, policy_start_date,
			policy_end_date, project_id, project_name, campaign_id, campaign_name,
			base_credit, minimum_credit, maximum_credit, bonus_percentage,
			tier_config, min_lead_quality_score, requires_approval, auto_approve,
			is_active, approval_required, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		policy.TenantID, policy.PartnerID, policy.PolicyCode, policy.PolicyName, policy.PolicyType,
		policy.CalculationType, policy.TimeUnitType, policy.TimeUnitValue, policy.PolicyStartDate,
		policy.PolicyEndDate, policy.ProjectID, policy.ProjectName, policy.CampaignID, policy.CampaignName,
		policy.BaseCredit, policy.MinimumCredit, policy.MaximumCredit, policy.BonusPercentage,
		tierJSON, policy.MinLeadQualityScore, policy.RequiresApproval, policy.AutoApprove,
		policy.IsActive, policy.ApprovalRequired, policy.CreatedBy, policy.CreatedAt, policy.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create credit policy: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get policy ID: %w", err)
	}

	policy.ID = id
	return policy, nil
}

// GetCreditPolicy retrieves a credit policy
func (s *partnerCreditPolicyService) GetCreditPolicy(ctx context.Context, tenantID string, policyID int64) (*models.PartnerCreditPolicy, error) {
	policy := &models.PartnerCreditPolicy{}

	query := `
		SELECT id, tenant_id, partner_id, policy_code, policy_name, policy_type,
		calculation_type, time_unit_type, time_unit_value, policy_start_date,
		policy_end_date, project_id, project_name, campaign_id, campaign_name,
		base_credit, minimum_credit, maximum_credit, bonus_percentage,
		tier_config, min_lead_quality_score, requires_approval, auto_approve,
		is_active, approval_required, total_leads_under_policy, total_credits_allocated,
		created_by, approved_by, approved_at, created_at, updated_at, deleted_at
		FROM partner_credit_policies
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	var tierJSON sql.NullString
	err := s.db.QueryRowContext(ctx, query, policyID, tenantID).Scan(
		&policy.ID, &policy.TenantID, &policy.PartnerID, &policy.PolicyCode, &policy.PolicyName, &policy.PolicyType,
		&policy.CalculationType, &policy.TimeUnitType, &policy.TimeUnitValue, &policy.PolicyStartDate,
		&policy.PolicyEndDate, &policy.ProjectID, &policy.ProjectName, &policy.CampaignID, &policy.CampaignName,
		&policy.BaseCredit, &policy.MinimumCredit, &policy.MaximumCredit, &policy.BonusPercentage,
		&tierJSON, &policy.MinLeadQualityScore, &policy.RequiresApproval, &policy.AutoApprove,
		&policy.IsActive, &policy.ApprovalRequired, &policy.TotalLeadsUnderPolicy, &policy.TotalCreditsAllocated,
		&policy.CreatedBy, &policy.ApprovedBy, &policy.ApprovedAt, &policy.CreatedAt, &policy.UpdatedAt, &policy.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("policy not found")
		}
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}

	if tierJSON.Valid {
		json.Unmarshal([]byte(tierJSON.String), &policy.TierConfig)
	}

	return policy, nil
}

// GetPartnerCreditPolicies retrieves all policies for a partner
func (s *partnerCreditPolicyService) GetPartnerCreditPolicies(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerCreditPolicy, error) {
	policies := []models.PartnerCreditPolicy{}

	query := `
		SELECT id, tenant_id, partner_id, policy_code, policy_name, policy_type,
		calculation_type, time_unit_type, time_unit_value, policy_start_date,
		policy_end_date, project_id, project_name, campaign_id, campaign_name,
		base_credit, minimum_credit, maximum_credit, bonus_percentage,
		tier_config, min_lead_quality_score, requires_approval, auto_approve,
		is_active, approval_required, total_leads_under_policy, total_credits_allocated,
		created_by, approved_by, approved_at, created_at, updated_at, deleted_at
		FROM partner_credit_policies
		WHERE tenant_id = ? AND partner_id = ? AND deleted_at IS NULL
		ORDER BY policy_start_date DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, partnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get policies: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.PartnerCreditPolicy
		var tierJSON sql.NullString

		err := rows.Scan(
			&p.ID, &p.TenantID, &p.PartnerID, &p.PolicyCode, &p.PolicyName, &p.PolicyType,
			&p.CalculationType, &p.TimeUnitType, &p.TimeUnitValue, &p.PolicyStartDate,
			&p.PolicyEndDate, &p.ProjectID, &p.ProjectName, &p.CampaignID, &p.CampaignName,
			&p.BaseCredit, &p.MinimumCredit, &p.MaximumCredit, &p.BonusPercentage,
			&tierJSON, &p.MinLeadQualityScore, &p.RequiresApproval, &p.AutoApprove,
			&p.IsActive, &p.ApprovalRequired, &p.TotalLeadsUnderPolicy, &p.TotalCreditsAllocated,
			&p.CreatedBy, &p.ApprovedBy, &p.ApprovedAt, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan policy: %w", err)
		}

		if tierJSON.Valid {
			json.Unmarshal([]byte(tierJSON.String), &p.TierConfig)
		}

		policies = append(policies, p)
	}

	return policies, rows.Err()
}

// GetActiveCreditPolicies retrieves active policies for a partner
func (s *partnerCreditPolicyService) GetActiveCreditPolicies(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerCreditPolicy, error) {
	policies := []models.PartnerCreditPolicy{}

	query := `
		SELECT id, tenant_id, partner_id, policy_code, policy_name, policy_type,
		calculation_type, time_unit_type, time_unit_value, policy_start_date,
		policy_end_date, project_id, project_name, campaign_id, campaign_name,
		base_credit, minimum_credit, maximum_credit, bonus_percentage,
		tier_config, min_lead_quality_score, requires_approval, auto_approve,
		is_active, approval_required, total_leads_under_policy, total_credits_allocated,
		created_by, approved_by, approved_at, created_at, updated_at, deleted_at
		FROM partner_credit_policies
		WHERE tenant_id = ? AND partner_id = ? AND is_active = TRUE AND deleted_at IS NULL
		AND policy_start_date <= NOW()
		AND (policy_end_date IS NULL OR policy_end_date >= NOW())
		ORDER BY policy_start_date DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, partnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active policies: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.PartnerCreditPolicy
		var tierJSON sql.NullString

		err := rows.Scan(
			&p.ID, &p.TenantID, &p.PartnerID, &p.PolicyCode, &p.PolicyName, &p.PolicyType,
			&p.CalculationType, &p.TimeUnitType, &p.TimeUnitValue, &p.PolicyStartDate,
			&p.PolicyEndDate, &p.ProjectID, &p.ProjectName, &p.CampaignID, &p.CampaignName,
			&p.BaseCredit, &p.MinimumCredit, &p.MaximumCredit, &p.BonusPercentage,
			&tierJSON, &p.MinLeadQualityScore, &p.RequiresApproval, &p.AutoApprove,
			&p.IsActive, &p.ApprovalRequired, &p.TotalLeadsUnderPolicy, &p.TotalCreditsAllocated,
			&p.CreatedBy, &p.ApprovedBy, &p.ApprovedAt, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan policy: %w", err)
		}

		if tierJSON.Valid {
			json.Unmarshal([]byte(tierJSON.String), &p.TierConfig)
		}

		policies = append(policies, p)
	}

	return policies, rows.Err()
}

// UpdateCreditPolicy updates a credit policy
func (s *partnerCreditPolicyService) UpdateCreditPolicy(ctx context.Context, tenantID string, policy *models.PartnerCreditPolicy) (*models.PartnerCreditPolicy, error) {
	policy.UpdatedAt = time.Now()

	tierJSON, _ := json.Marshal(policy.TierConfig)

	query := `
		UPDATE partner_credit_policies SET
			policy_name = ?, base_credit = ?, minimum_credit = ?, maximum_credit = ?,
			bonus_percentage = ?, tier_config = ?, min_lead_quality_score = ?,
			requires_approval = ?, auto_approve = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		policy.PolicyName, policy.BaseCredit, policy.MinimumCredit, policy.MaximumCredit,
		policy.BonusPercentage, tierJSON, policy.MinLeadQualityScore,
		policy.RequiresApproval, policy.AutoApprove, policy.UpdatedAt,
		policy.ID, tenantID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update policy: %w", err)
	}

	return policy, nil
}

// ApproveCreditPolicy approves a credit policy
func (s *partnerCreditPolicyService) ApproveCreditPolicy(ctx context.Context, tenantID string, policyID int64, approvedBy int64) error {
	now := time.Now()
	query := `
		UPDATE partner_credit_policies SET
			approval_required = FALSE, approved_by = ?, approved_at = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query, approvedBy, now, now, policyID, tenantID)
	return err
}

// DeactivateCreditPolicy deactivates a credit policy
func (s *partnerCreditPolicyService) DeactivateCreditPolicy(ctx context.Context, tenantID string, policyID int64) error {
	query := `
		UPDATE partner_credit_policies SET
			is_active = FALSE, deleted_at = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, query, now, now, policyID, tenantID)
	return err
}

// CalculateLeadCredit calculates credit for a lead based on applicable policies
func (s *partnerCreditPolicyService) CalculateLeadCredit(ctx context.Context, tenantID string, partnerLeadID int64) (float64, []int64, error) {
	totalCredit := 0.0
	policyIDs := []int64{}

	// Get lead details
	var partnerID int64
	var qualityScore float64
	var leadData models.LeadData

	query := `
		SELECT partner_id, quality_score, lead_data
		FROM partner_leads
		WHERE id = ? AND tenant_id = ?
	`

	var leadDataJSON sql.NullString
	err := s.db.QueryRowContext(ctx, query, partnerLeadID, tenantID).Scan(
		&partnerID, &qualityScore, &leadDataJSON,
	)

	if err != nil {
		return 0, nil, fmt.Errorf("failed to get lead: %w", err)
	}

	if leadDataJSON.Valid {
		json.Unmarshal([]byte(leadDataJSON.String), &leadData)
	}

	// Get applicable policies
	policies, err := s.GetActiveCreditPolicies(ctx, tenantID, partnerID)
	if err != nil {
		return 0, nil, err
	}

	for _, policy := range policies {
		// Check if lead quality meets minimum
		if qualityScore < policy.MinLeadQualityScore {
			continue
		}

		credit := s.calculateCreditAmount(policy, qualityScore)
		if credit > 0 {
			totalCredit += credit
			policyIDs = append(policyIDs, policy.ID)
		}
	}

	return totalCredit, policyIDs, nil
}

// calculateCreditAmount calculates credit based on policy type
func (s *partnerCreditPolicyService) calculateCreditAmount(policy models.PartnerCreditPolicy, qualityScore float64) float64 {
	credit := 0.0

	switch policy.CalculationType {
	case models.CreditPolicyCalcPercentage:
		credit = policy.BaseCredit // Already a percentage
	case models.CreditPolicyCalcFixedPrice:
		credit = policy.BaseCredit
	case models.CreditPolicyCalcTiered:
		// Find applicable tier
		for _, tier := range policy.TierConfig.Tiers {
			if policy.TotalLeadsUnderPolicy >= int64(tier.MinLeads) {
				credit = tier.CreditAmount
				break
			}
		}
	case models.CreditPolicyCalcConversion:
		// No credit on submission, only on conversion
		credit = 0
	case models.CreditPolicyCalcRevenueshare:
		credit = policy.BaseCredit // % of revenue
	}

	// Apply bonus if quality score high
	if qualityScore >= 80 {
		credit += (credit * policy.BonusPercentage / 100)
	}

	// Apply minimum/maximum bounds
	if credit > 0 && policy.MinimumCredit > 0 && credit < policy.MinimumCredit {
		credit = policy.MinimumCredit
	}
	if credit > 0 && policy.MaximumCredit > 0 && credit > policy.MaximumCredit {
		credit = policy.MaximumCredit
	}

	return credit
}

// GetApplicablePolicies retrieves policies applicable to a lead
func (s *partnerCreditPolicyService) GetApplicablePolicies(ctx context.Context, tenantID string, partnerID int64, leadData *models.LeadData) ([]models.PartnerCreditPolicy, error) {
	policies, err := s.GetActiveCreditPolicies(ctx, tenantID, partnerID)
	if err != nil {
		return nil, err
	}

	return policies, nil
}

// GetPolicyMappings retrieves policy mappings for a lead
func (s *partnerCreditPolicyService) GetPolicyMappings(ctx context.Context, tenantID string, partnerLeadID int64) ([]models.PartnerCreditPolicyMapping, error) {
	mappings := []models.PartnerCreditPolicyMapping{}

	query := `
		SELECT id, tenant_id, partner_lead_id, policy_id, calculated_credit, reason, created_at
		FROM partner_credit_policy_mappings
		WHERE tenant_id = ? AND partner_lead_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, partnerLeadID)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy mappings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m models.PartnerCreditPolicyMapping
		err := rows.Scan(
			&m.ID, &m.TenantID, &m.PartnerLeadID, &m.PolicyID, &m.CalculatedCredit, &m.Reason, &m.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mapping: %w", err)
		}
		mappings = append(mappings, m)
	}

	return mappings, rows.Err()
}
