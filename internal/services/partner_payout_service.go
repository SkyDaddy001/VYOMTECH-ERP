package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"vyomtech-backend/internal/models"
)

// PartnerPayoutService handles partner payout approvals and management
type PartnerPayoutService interface {
	// Payout Generation
	GeneratePayoutPeriod(ctx context.Context, tenantID string, partnerID int64, periodStart time.Time, periodEnd time.Time) (*models.PartnerPayout, error)
	CreatePayout(ctx context.Context, tenantID string, payout *models.PartnerPayout) (*models.PartnerPayout, error)
	GetPayout(ctx context.Context, tenantID string, payoutID int64) (*models.PartnerPayout, error)
	GetPayouts(ctx context.Context, tenantID string, partnerID int64, limit int, offset int) ([]models.PartnerPayout, error)
	GetPendingPayouts(ctx context.Context, tenantID string, limit int, offset int) ([]models.PartnerPayout, error)

	// Payout Review & Approval
	ApprovePayout(ctx context.Context, tenantID string, payoutID int64, approvedAmount float64, approvedBy int64) error
	RejectPayout(ctx context.Context, tenantID string, payoutID int64, rejectionNotes string) error
	PartiallyApprovePayout(ctx context.Context, tenantID string, payoutID int64, approvedAmount float64, approvedBy int64) error

	// Payout Details
	GetPayoutDetails(ctx context.Context, tenantID string, payoutID int64, limit int, offset int) ([]models.PartnerPayoutDetail, error)
	AddPayoutDetail(ctx context.Context, tenantID string, detail *models.PartnerPayoutDetail) (*models.PartnerPayoutDetail, error)
	ApprovePayoutDetail(ctx context.Context, tenantID string, detailID int64) error
	RejectPayoutDetail(ctx context.Context, tenantID string, detailID int64, notes string) error

	// Payout Processing
	MarkPayoutAsPaid(ctx context.Context, tenantID string, payoutID int64, paymentDate time.Time, referenceNumber string) error
	GetPayoutStats(ctx context.Context, tenantID string, partnerID int64) (*PayoutStats, error)
}

type partnerPayoutService struct {
	db             *sql.DB
	partnerService PartnerService
	leadService    PartnerLeadService
}

// PayoutStats represents payout statistics
type PayoutStats struct {
	TotalPayoutsGenerated  int64   `json:"total_payouts_generated"`
	ApprovedPayouts        int64   `json:"approved_payouts"`
	RejectedPayouts        int64   `json:"rejected_payouts"`
	TotalAmountGenerated   float64 `json:"total_amount_generated"`
	TotalAmountApproved    float64 `json:"total_amount_approved"`
	TotalAmountPaid        float64 `json:"total_amount_paid"`
	AverageApprovalRate    float64 `json:"average_approval_rate"`
	PendingApprovalAmount  float64 `json:"pending_approval_amount"`
	AverageDaysToApproval  float64 `json:"average_days_to_approval"`
}

// NewPartnerPayoutService creates a new payout service
func NewPartnerPayoutService(db *sql.DB, partnerService PartnerService, leadService PartnerLeadService) PartnerPayoutService {
	return &partnerPayoutService{
		db:             db,
		partnerService: partnerService,
		leadService:    leadService,
	}
}

// GeneratePayoutPeriod generates a payout for a period
func (s *partnerPayoutService) GeneratePayoutPeriod(ctx context.Context, tenantID string, partnerID int64, periodStart time.Time, periodEnd time.Time) (*models.PartnerPayout, error) {
	// Get all approved leads for the period
	query := `
		SELECT COUNT(*), SUM(credit_amount)
		FROM partner_leads
		WHERE partner_id = ? AND tenant_id = ? AND status = 'approved'
		AND created_at >= ? AND created_at <= ?
	`

	var totalLeads int64
	var totalAmount sql.NullFloat64

	err := s.db.QueryRowContext(ctx, query, partnerID, tenantID, periodStart, periodEnd).Scan(&totalLeads, &totalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to generate payout period: %w", err)
	}

	// Get converted leads count
	convertedQuery := `
		SELECT COUNT(*) FROM partner_leads
		WHERE partner_id = ? AND tenant_id = ? AND status = 'converted'
		AND created_at >= ? AND created_at <= ?
	`

	var convertedLeads int64
	s.db.QueryRowContext(ctx, convertedQuery, partnerID, tenantID, periodStart, periodEnd).Scan(&convertedLeads)

	payout := &models.PartnerPayout{
		TenantID:        tenantID,
		PartnerID:       partnerID,
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		TotalLeadsCount: totalLeads,
		ConvertedLeads:  convertedLeads,
		Status:          "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if totalAmount.Valid {
		payout.TotalAmount = totalAmount.Float64
	}

	return payout, nil
}

// CreatePayout creates a new payout record
func (s *partnerPayoutService) CreatePayout(ctx context.Context, tenantID string, payout *models.PartnerPayout) (*models.PartnerPayout, error) {
	if payout.PartnerID == 0 {
		return nil, errors.New("partner_id is required")
	}

	payout.TenantID = tenantID
	payout.CreatedAt = time.Now()
	payout.UpdatedAt = time.Now()

	if payout.Status == "" {
		payout.Status = "pending"
	}

	query := `
		INSERT INTO partner_payouts (
			tenant_id, partner_id, period_start, period_end, total_leads_count,
			approved_leads, converted_leads, total_amount, status,
			notes, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		payout.TenantID, payout.PartnerID, payout.PeriodStart, payout.PeriodEnd,
		payout.TotalLeadsCount, payout.ApprovedLeads, payout.ConvertedLeads, payout.TotalAmount,
		payout.Status, payout.Notes, payout.CreatedAt, payout.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get payout ID: %w", err)
	}

	payout.ID = id
	return payout, nil
}

// GetPayout retrieves a payout
func (s *partnerPayoutService) GetPayout(ctx context.Context, tenantID string, payoutID int64) (*models.PartnerPayout, error) {
	payout := &models.PartnerPayout{}

	query := `
		SELECT id, tenant_id, partner_id, period_start, period_end, total_leads_count,
		approved_leads, converted_leads, total_amount, approved_amount, rejected_amount,
		status, payment_method, payment_date, reference_number, reviewed_by, approved_by,
		approved_at, rejection_notes, notes, created_at, updated_at
		FROM partner_payouts
		WHERE id = ? AND tenant_id = ?
	`

	err := s.db.QueryRowContext(ctx, query, payoutID, tenantID).Scan(
		&payout.ID, &payout.TenantID, &payout.PartnerID, &payout.PeriodStart, &payout.PeriodEnd, &payout.TotalLeadsCount,
		&payout.ApprovedLeads, &payout.ConvertedLeads, &payout.TotalAmount, &payout.ApprovedAmount, &payout.RejectedAmount,
		&payout.Status, &payout.PaymentMethod, &payout.PaymentDate, &payout.ReferenceNumber, &payout.ReviewedBy, &payout.ApprovedBy,
		&payout.ApprovedAt, &payout.RejectionNotes, &payout.Notes, &payout.CreatedAt, &payout.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payout not found")
		}
		return nil, fmt.Errorf("failed to get payout: %w", err)
	}

	return payout, nil
}

// GetPayouts retrieves payouts for a partner
func (s *partnerPayoutService) GetPayouts(ctx context.Context, tenantID string, partnerID int64, limit int, offset int) ([]models.PartnerPayout, error) {
	if limit == 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	payouts := []models.PartnerPayout{}

	query := `
		SELECT id, tenant_id, partner_id, period_start, period_end, total_leads_count,
		approved_leads, converted_leads, total_amount, approved_amount, rejected_amount,
		status, payment_method, payment_date, reference_number, reviewed_by, approved_by,
		approved_at, rejection_notes, notes, created_at, updated_at
		FROM partner_payouts
		WHERE tenant_id = ? AND partner_id = ?
		ORDER BY period_end DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, partnerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get payouts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.PartnerPayout
		err := rows.Scan(
			&p.ID, &p.TenantID, &p.PartnerID, &p.PeriodStart, &p.PeriodEnd, &p.TotalLeadsCount,
			&p.ApprovedLeads, &p.ConvertedLeads, &p.TotalAmount, &p.ApprovedAmount, &p.RejectedAmount,
			&p.Status, &p.PaymentMethod, &p.PaymentDate, &p.ReferenceNumber, &p.ReviewedBy, &p.ApprovedBy,
			&p.ApprovedAt, &p.RejectionNotes, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payout: %w", err)
		}
		payouts = append(payouts, p)
	}

	return payouts, rows.Err()
}

// GetPendingPayouts retrieves payouts pending approval
func (s *partnerPayoutService) GetPendingPayouts(ctx context.Context, tenantID string, limit int, offset int) ([]models.PartnerPayout, error) {
	if limit == 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	payouts := []models.PartnerPayout{}

	query := `
		SELECT id, tenant_id, partner_id, period_start, period_end, total_leads_count,
		approved_leads, converted_leads, total_amount, approved_amount, rejected_amount,
		status, payment_method, payment_date, reference_number, reviewed_by, approved_by,
		approved_at, rejection_notes, notes, created_at, updated_at
		FROM partner_payouts
		WHERE tenant_id = ? AND status = 'pending'
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending payouts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.PartnerPayout
		err := rows.Scan(
			&p.ID, &p.TenantID, &p.PartnerID, &p.PeriodStart, &p.PeriodEnd, &p.TotalLeadsCount,
			&p.ApprovedLeads, &p.ConvertedLeads, &p.TotalAmount, &p.ApprovedAmount, &p.RejectedAmount,
			&p.Status, &p.PaymentMethod, &p.PaymentDate, &p.ReferenceNumber, &p.ReviewedBy, &p.ApprovedBy,
			&p.ApprovedAt, &p.RejectionNotes, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payout: %w", err)
		}
		payouts = append(payouts, p)
	}

	return payouts, rows.Err()
}

// ApprovePayout approves a payout
func (s *partnerPayoutService) ApprovePayout(ctx context.Context, tenantID string, payoutID int64, approvedAmount float64, approvedBy int64) error {
	now := time.Now()
	query := `
		UPDATE partner_payouts SET
			status = 'approved', approved_amount = ?, approved_by = ?, approved_at = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		approvedAmount, approvedBy, now, now,
		payoutID, tenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to approve payout: %w", err)
	}

	// Update partner balance
	payout, _ := s.GetPayout(ctx, tenantID, payoutID)
	s.db.ExecContext(ctx, `
		UPDATE partners SET
			withdrawn_amount = withdrawn_amount + ?,
			available_balance = available_balance - ?
		WHERE id = ? AND tenant_id = ?
	`, approvedAmount, approvedAmount, payout.PartnerID, tenantID)

	return nil
}

// RejectPayout rejects a payout
func (s *partnerPayoutService) RejectPayout(ctx context.Context, tenantID string, payoutID int64, rejectionNotes string) error {
	query := `
		UPDATE partner_payouts SET
			status = 'rejected', rejection_notes = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		rejectionNotes, time.Now(),
		payoutID, tenantID,
	)

	return err
}

// PartiallyApprovePayout partially approves a payout
func (s *partnerPayoutService) PartiallyApprovePayout(ctx context.Context, tenantID string, payoutID int64, approvedAmount float64, approvedBy int64) error {
	now := time.Now()
	query := `
		UPDATE partner_payouts SET
			status = 'partially_approved', approved_amount = ?, rejected_amount = total_amount - ?,
			approved_by = ?, approved_at = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		approvedAmount, approvedAmount,
		approvedBy, now, now,
		payoutID, tenantID,
	)

	return err
}

// GetPayoutDetails retrieves payout line items
func (s *partnerPayoutService) GetPayoutDetails(ctx context.Context, tenantID string, payoutID int64, limit int, offset int) ([]models.PartnerPayoutDetail, error) {
	if limit == 0 {
		limit = 50
	}

	details := []models.PartnerPayoutDetail{}

	query := `
		SELECT id, payout_id, partner_lead_id, lead_submission_id, amount, status, approval_notes, created_at
		FROM partner_payout_details
		WHERE payout_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, payoutID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get payout details: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d models.PartnerPayoutDetail
		err := rows.Scan(
			&d.ID, &d.PayoutID, &d.PartnerLeadID, &d.LeadSubmissionID, &d.Amount, &d.Status, &d.ApprovalNotes, &d.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payout detail: %w", err)
		}
		details = append(details, d)
	}

	return details, rows.Err()
}

// AddPayoutDetail adds a line item to a payout
func (s *partnerPayoutService) AddPayoutDetail(ctx context.Context, tenantID string, detail *models.PartnerPayoutDetail) (*models.PartnerPayoutDetail, error) {
	detail.CreatedAt = time.Now()

	query := `
		INSERT INTO partner_payout_details (
			payout_id, partner_lead_id, lead_submission_id, amount, status, created_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		detail.PayoutID, detail.PartnerLeadID, detail.LeadSubmissionID, detail.Amount, detail.Status, detail.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to add payout detail: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get detail ID: %w", err)
	}

	detail.ID = id
	return detail, nil
}

// ApprovePayoutDetail approves a payout line item
func (s *partnerPayoutService) ApprovePayoutDetail(ctx context.Context, tenantID string, detailID int64) error {
	query := `UPDATE partner_payout_details SET status = 'approved' WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, detailID)
	return err
}

// RejectPayoutDetail rejects a payout line item
func (s *partnerPayoutService) RejectPayoutDetail(ctx context.Context, tenantID string, detailID int64, notes string) error {
	query := `UPDATE partner_payout_details SET status = 'rejected', approval_notes = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, notes, detailID)
	return err
}

// MarkPayoutAsPaid marks a payout as paid
func (s *partnerPayoutService) MarkPayoutAsPaid(ctx context.Context, tenantID string, payoutID int64, paymentDate time.Time, referenceNumber string) error {
	query := `
		UPDATE partner_payouts SET
			status = 'paid', payment_date = ?, reference_number = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		paymentDate, referenceNumber, time.Now(),
		payoutID, tenantID,
	)

	return err
}

// GetPayoutStats retrieves payout statistics
func (s *partnerPayoutService) GetPayoutStats(ctx context.Context, tenantID string, partnerID int64) (*PayoutStats, error) {
	stats := &PayoutStats{}

	query := `
		SELECT
			COUNT(*) as total_payouts,
			SUM(CASE WHEN status = 'approved' OR status = 'paid' THEN 1 ELSE 0 END) as approved_payouts,
			SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END) as rejected_payouts,
			SUM(total_amount) as total_amount,
			SUM(CASE WHEN status = 'approved' OR status = 'paid' THEN approved_amount ELSE 0 END) as approved_amount,
			SUM(CASE WHEN status = 'paid' THEN approved_amount ELSE 0 END) as paid_amount,
			SUM(CASE WHEN status = 'pending' THEN total_amount ELSE 0 END) as pending_amount
		FROM partner_payouts
		WHERE tenant_id = ? AND partner_id = ?
	`

	var totalPayouts, approvedPayouts, rejectedPayouts sql.NullInt64
	var totalAmount, approvedAmount, paidAmount, pendingAmount sql.NullFloat64

	err := s.db.QueryRowContext(ctx, query, tenantID, partnerID).Scan(
		&totalPayouts, &approvedPayouts, &rejectedPayouts,
		&totalAmount, &approvedAmount, &paidAmount, &pendingAmount,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get payout stats: %w", err)
	}

	if totalPayouts.Valid {
		stats.TotalPayoutsGenerated = totalPayouts.Int64
	}
	if approvedPayouts.Valid {
		stats.ApprovedPayouts = approvedPayouts.Int64
	}
	if rejectedPayouts.Valid {
		stats.RejectedPayouts = rejectedPayouts.Int64
	}
	if totalAmount.Valid {
		stats.TotalAmountGenerated = totalAmount.Float64
	}
	if approvedAmount.Valid {
		stats.TotalAmountApproved = approvedAmount.Float64
	}
	if paidAmount.Valid {
		stats.TotalAmountPaid = paidAmount.Float64
	}
	if pendingAmount.Valid {
		stats.PendingApprovalAmount = pendingAmount.Float64
	}

	if stats.TotalPayoutsGenerated > 0 {
		stats.AverageApprovalRate = (float64(stats.ApprovedPayouts) / float64(stats.TotalPayoutsGenerated)) * 100
	}

	return stats, nil
}
