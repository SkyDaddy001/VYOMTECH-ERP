package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"vyomtech-backend/internal/models"
)

// PartnerLeadService handles partner lead submissions and referrals
type PartnerLeadService interface {
	// Lead Submission
	SubmitPartnerLead(ctx context.Context, tenantID string, lead *models.PartnerLead) (*models.PartnerLead, error)
	GetPartnerLead(ctx context.Context, tenantID string, leadID int64) (*models.PartnerLead, error)
	GetPartnerLeads(ctx context.Context, tenantID string, filter *models.PartnerLeadFilter) ([]models.PartnerLead, int64, error)
	UpdatePartnerLeadStatus(ctx context.Context, tenantID string, leadID int64, status string, notes string) error

	// Lead Review & Approval
	ApprovePartnerLead(ctx context.Context, tenantID string, leadID int64, approvedBy int64, actualLeadID int64) error
	RejectPartnerLead(ctx context.Context, tenantID string, leadID int64, rejectionReason string, rejectedBy int64) error
	GetPendingLeadsForReview(ctx context.Context, tenantID string, limit int, offset int) ([]models.PartnerLead, error)

	// Lead Credit Management
	GetLeadCredits(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerLeadCredit, error)
	SubmitLeadCreditApprovalRequest(ctx context.Context, tenantID string, credit *models.PartnerLeadCredit) (*models.PartnerLeadCredit, error)
	ApproveLeadCredit(ctx context.Context, tenantID string, creditID int64, approvedBy int64) error
	RejectLeadCredit(ctx context.Context, tenantID string, creditID int64, rejectionReason string) error

	// Activity Tracking
	LogPartnerActivity(ctx context.Context, tenantID string, activity *models.PartnerActivity) error
	GetPartnerActivity(ctx context.Context, tenantID string, partnerID int64, limit int, offset int) ([]models.PartnerActivity, error)
}

type partnerLeadService struct {
	db             *sql.DB
	partnerService PartnerService
}

// NewPartnerLeadService creates a new partner lead service
func NewPartnerLeadService(db *sql.DB, partnerService PartnerService) PartnerLeadService {
	return &partnerLeadService{
		db:             db,
		partnerService: partnerService,
	}
}

// SubmitPartnerLead submits a lead from a partner
func (s *partnerLeadService) SubmitPartnerLead(ctx context.Context, tenantID string, lead *models.PartnerLead) (*models.PartnerLead, error) {
	if lead.PartnerID == 0 || lead.SubmittedBy == 0 {
		return nil, errors.New("partner_id and submitted_by are required")
	}

	lead.TenantID = tenantID
	lead.Status = "submitted"
	lead.CreatedAt = time.Now()
	lead.UpdatedAt = time.Now()

	// Calculate quality score
	qualityService := s.partnerService
	lead.QualityScore = qualityService.CalculateLeadQualityScore(ctx, &lead.LeadData)

	query := `
		INSERT INTO partner_leads (
			tenant_id, partner_id, submission_type, status, lead_data, quality_score,
			submitted_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		lead.TenantID, lead.PartnerID, lead.SubmissionType, lead.Status, lead.LeadData, lead.QualityScore,
		lead.SubmittedBy, lead.CreatedAt, lead.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to submit partner lead: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get lead ID: %w", err)
	}

	lead.ID = id

	// Log activity
	s.LogPartnerActivity(ctx, tenantID, &models.PartnerActivity{
		PartnerID:  lead.PartnerID,
		UserID:     &lead.SubmittedBy,
		Action:     "lead_submitted",
		Resource:   "lead",
		ResourceID: id,
		TenantID:   tenantID,
		CreatedAt:  time.Now(),
	})

	return lead, nil
}

// GetPartnerLead retrieves a partner lead
func (s *partnerLeadService) GetPartnerLead(ctx context.Context, tenantID string, leadID int64) (*models.PartnerLead, error) {
	lead := &models.PartnerLead{}

	query := `
		SELECT id, tenant_id, partner_id, lead_id, submission_type, status,
		lead_data, quality_score, rejection_reason, reviewed_by, reviewed_at,
		submitted_by, conversion_date, credit_amount, credit_status,
		created_at, updated_at
		FROM partner_leads
		WHERE id = ? AND tenant_id = ?
	`

	err := s.db.QueryRowContext(ctx, query, leadID, tenantID).Scan(
		&lead.ID, &lead.TenantID, &lead.PartnerID, &lead.LeadID, &lead.SubmissionType, &lead.Status,
		&lead.LeadData, &lead.QualityScore, &lead.RejectionReason, &lead.ReviewedBy, &lead.ReviewedAt,
		&lead.SubmittedBy, &lead.ConversionDate, &lead.CreditAmount, &lead.CreditStatus,
		&lead.CreatedAt, &lead.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("lead not found")
		}
		return nil, fmt.Errorf("failed to get partner lead: %w", err)
	}

	return lead, nil
}

// GetPartnerLeads retrieves partner leads with filtering
func (s *partnerLeadService) GetPartnerLeads(ctx context.Context, tenantID string, filter *models.PartnerLeadFilter) ([]models.PartnerLead, int64, error) {
	if filter == nil {
		filter = &models.PartnerLeadFilter{Limit: 50}
	}
	if filter.Limit == 0 {
		filter.Limit = 50
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	leads := []models.PartnerLead{}
	query := `
		SELECT id, tenant_id, partner_id, lead_id, submission_type, status,
		lead_data, quality_score, rejection_reason, reviewed_by, reviewed_at,
		submitted_by, conversion_date, credit_amount, credit_status,
		created_at, updated_at
		FROM partner_leads
		WHERE tenant_id = ?`

	args := []interface{}{tenantID}

	if filter.PartnerID != 0 {
		query += " AND partner_id = ?"
		args = append(args, filter.PartnerID)
	}
	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}
	if filter.SubmissionType != "" {
		query += " AND submission_type = ?"
		args = append(args, filter.SubmissionType)
	}
	if filter.QualityScoreMin > 0 {
		query += " AND quality_score >= ?"
		args = append(args, filter.QualityScoreMin)
	}
	if filter.QualityScoreMax > 0 {
		query += " AND quality_score <= ?"
		args = append(args, filter.QualityScoreMax)
	}

	// Get count
	countQuery := strings.Replace(query, "SELECT id, tenant_id, partner_id, lead_id, submission_type, status, lead_data, quality_score, rejection_reason, reviewed_by, reviewed_at, submitted_by, conversion_date, credit_amount, credit_status, created_at, updated_at", "SELECT COUNT(*)", 1)
	var total int64
	s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get partner leads: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lead models.PartnerLead
		err := rows.Scan(
			&lead.ID, &lead.TenantID, &lead.PartnerID, &lead.LeadID, &lead.SubmissionType, &lead.Status,
			&lead.LeadData, &lead.QualityScore, &lead.RejectionReason, &lead.ReviewedBy, &lead.ReviewedAt,
			&lead.SubmittedBy, &lead.ConversionDate, &lead.CreditAmount, &lead.CreditStatus,
			&lead.CreatedAt, &lead.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan partner lead: %w", err)
		}
		leads = append(leads, lead)
	}

	return leads, total, rows.Err()
}

// UpdatePartnerLeadStatus updates a partner lead status
func (s *partnerLeadService) UpdatePartnerLeadStatus(ctx context.Context, tenantID string, leadID int64, status string, notes string) error {
	query := `
		UPDATE partner_leads SET status = ?, updated_at = ? WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query, status, time.Now(), leadID, tenantID)
	return err
}

// ApprovePartnerLead approves a partner lead
func (s *partnerLeadService) ApprovePartnerLead(ctx context.Context, tenantID string, leadID int64, approvedBy int64, actualLeadID int64) error {
	now := time.Now()
	query := `
		UPDATE partner_leads SET
			status = ?, lead_id = ?, reviewed_by = ?, reviewed_at = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		"approved", actualLeadID, approvedBy, now, now,
		leadID, tenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to approve partner lead: %w", err)
	}

	// Update partner stats
	s.db.ExecContext(ctx, `
		UPDATE partners SET
			approved_leads = approved_leads + 1,
			current_month_leads = current_month_leads + 1
		WHERE id = ? AND tenant_id = ?
	`, leadID, tenantID)

	return nil
}

// RejectPartnerLead rejects a partner lead
func (s *partnerLeadService) RejectPartnerLead(ctx context.Context, tenantID string, leadID int64, rejectionReason string, rejectedBy int64) error {
	now := time.Now()
	query := `
		UPDATE partner_leads SET
			status = ?, rejection_reason = ?, reviewed_by = ?, reviewed_at = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		"rejected", rejectionReason, rejectedBy, now, now,
		leadID, tenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to reject partner lead: %w", err)
	}

	// Update partner stats
	s.db.ExecContext(ctx, `
		UPDATE partners SET rejected_leads = rejected_leads + 1 WHERE id = ? AND tenant_id = ?
	`, leadID, tenantID)

	return nil
}

// GetPendingLeadsForReview retrieves pending leads for management review
func (s *partnerLeadService) GetPendingLeadsForReview(ctx context.Context, tenantID string, limit int, offset int) ([]models.PartnerLead, error) {
	if limit == 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	leads := []models.PartnerLead{}
	query := `
		SELECT id, tenant_id, partner_id, lead_id, submission_type, status,
		lead_data, quality_score, rejection_reason, reviewed_by, reviewed_at,
		submitted_by, conversion_date, credit_amount, credit_status,
		created_at, updated_at
		FROM partner_leads
		WHERE tenant_id = ? AND status = 'submitted'
		ORDER BY quality_score DESC, created_at ASC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending leads: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lead models.PartnerLead
		err := rows.Scan(
			&lead.ID, &lead.TenantID, &lead.PartnerID, &lead.LeadID, &lead.SubmissionType, &lead.Status,
			&lead.LeadData, &lead.QualityScore, &lead.RejectionReason, &lead.ReviewedBy, &lead.ReviewedAt,
			&lead.SubmittedBy, &lead.ConversionDate, &lead.CreditAmount, &lead.CreditStatus,
			&lead.CreatedAt, &lead.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lead: %w", err)
		}
		leads = append(leads, lead)
	}

	return leads, rows.Err()
}

// GetLeadCredits retrieves pending lead credits for a partner
func (s *partnerLeadService) GetLeadCredits(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerLeadCredit, error) {
	credits := []models.PartnerLeadCredit{}

	query := `
		SELECT id, tenant_id, partner_lead_id, partner_id, credit_amount, calculation_type,
		status, approved_by, approved_at, rejection_reason, notes,
		created_at, updated_at
		FROM partner_lead_credits
		WHERE partner_id = ? AND tenant_id = ? AND status = 'pending_approval'
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, partnerID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead credits: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var credit models.PartnerLeadCredit
		err := rows.Scan(
			&credit.ID, &credit.TenantID, &credit.PartnerLeadID, &credit.PartnerID, &credit.CreditAmount, &credit.CalculationType,
			&credit.Status, &credit.ApprovedBy, &credit.ApprovedAt, &credit.RejectionReason, &credit.Notes,
			&credit.CreatedAt, &credit.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credit: %w", err)
		}
		credits = append(credits, credit)
	}

	return credits, rows.Err()
}

// SubmitLeadCreditApprovalRequest submits a lead credit for approval
func (s *partnerLeadService) SubmitLeadCreditApprovalRequest(ctx context.Context, tenantID string, credit *models.PartnerLeadCredit) (*models.PartnerLeadCredit, error) {
	if credit.PartnerLeadID == 0 || credit.PartnerID == 0 {
		return nil, errors.New("partner_lead_id and partner_id are required")
	}

	credit.TenantID = tenantID
	credit.Status = "pending_approval"
	credit.CreatedAt = time.Now()
	credit.UpdatedAt = time.Now()

	query := `
		INSERT INTO partner_lead_credits (
			tenant_id, partner_lead_id, partner_id, credit_amount, calculation_type,
			status, notes, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		credit.TenantID, credit.PartnerLeadID, credit.PartnerID, credit.CreditAmount, credit.CalculationType,
		credit.Status, credit.Notes, credit.CreatedAt, credit.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to submit lead credit: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get credit ID: %w", err)
	}

	credit.ID = id
	return credit, nil
}

// ApproveLeadCredit approves a lead credit request
func (s *partnerLeadService) ApproveLeadCredit(ctx context.Context, tenantID string, creditID int64, approvedBy int64) error {
	now := time.Now()
	query := `
		UPDATE partner_lead_credits SET
			status = 'approved', approved_by = ?, approved_at = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	result, err := s.db.ExecContext(ctx, query,
		approvedBy, now, now,
		creditID, tenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to approve lead credit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("credit not found or already processed")
	}

	// Get credit details to update partner balance
	credit := &models.PartnerLeadCredit{}
	s.db.QueryRowContext(ctx, `
		SELECT partner_id, credit_amount FROM partner_lead_credits WHERE id = ? AND tenant_id = ?
	`, creditID, tenantID).Scan(&credit.PartnerID, &credit.CreditAmount)

	// Update partner available balance
	s.db.ExecContext(ctx, `
		UPDATE partners SET
			available_balance = available_balance + ?,
			pending_payout_amount = pending_payout_amount + ?
		WHERE id = ? AND tenant_id = ?
	`, credit.CreditAmount, credit.CreditAmount, credit.PartnerID, tenantID)

	return nil
}

// RejectLeadCredit rejects a lead credit request
func (s *partnerLeadService) RejectLeadCredit(ctx context.Context, tenantID string, creditID int64, rejectionReason string) error {
	query := `
		UPDATE partner_lead_credits SET
			status = 'rejected', rejection_reason = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		rejectionReason, time.Now(),
		creditID, tenantID,
	)

	return err
}

// LogPartnerActivity logs a partner activity
func (s *partnerLeadService) LogPartnerActivity(ctx context.Context, tenantID string, activity *models.PartnerActivity) error {
	activity.TenantID = tenantID
	activity.CreatedAt = time.Now()

	query := `
		INSERT INTO partner_activities (
			tenant_id, partner_id, user_id, action, resource, resource_id, details, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		activity.TenantID, activity.PartnerID, activity.UserID, activity.Action,
		activity.Resource, activity.ResourceID, activity.Details, activity.CreatedAt,
	)

	return err
}

// GetPartnerActivity retrieves partner activities
func (s *partnerLeadService) GetPartnerActivity(ctx context.Context, tenantID string, partnerID int64, limit int, offset int) ([]models.PartnerActivity, error) {
	if limit == 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	activities := []models.PartnerActivity{}

	query := `
		SELECT id, tenant_id, partner_id, user_id, action, resource, resource_id, details, created_at
		FROM partner_activities
		WHERE tenant_id = ? AND partner_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, partnerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get partner activities: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var activity models.PartnerActivity
		err := rows.Scan(
			&activity.ID, &activity.TenantID, &activity.PartnerID, &activity.UserID,
			&activity.Action, &activity.Resource, &activity.ResourceID, &activity.Details, &activity.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, activity)
	}

	return activities, rows.Err()
}
