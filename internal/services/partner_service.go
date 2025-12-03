package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"vyomtech-backend/internal/models"
)

// PartnerService handles partner operations
type PartnerService interface {
	// Partner Management
	CreatePartner(ctx context.Context, tenantID string, partner *models.Partner) (*models.Partner, error)
	GetPartner(ctx context.Context, tenantID string, partnerID int64) (*models.Partner, error)
	GetPartnerByCode(ctx context.Context, tenantID string, partnerCode string) (*models.Partner, error)
	GetPartners(ctx context.Context, tenantID string, filter *models.PartnerFilter) ([]models.Partner, int64, error)
	UpdatePartner(ctx context.Context, tenantID string, partner *models.Partner) (*models.Partner, error)
	UpdatePartnerStatus(ctx context.Context, tenantID string, partnerID int64, status models.PartnerStatus, reason string, approvedBy int64) error
	DeactivatePartner(ctx context.Context, tenantID string, partnerID int64, reason string) error
	SuspendPartner(ctx context.Context, tenantID string, partnerID int64, reason string) error

	// Partner Users
	CreatePartnerUser(ctx context.Context, tenantID string, user *models.PartnerUser) (*models.PartnerUser, error)
	GetPartnerUser(ctx context.Context, tenantID string, userID int64) (*models.PartnerUser, error)
	GetPartnerUserByEmail(ctx context.Context, tenantID string, email string) (*models.PartnerUser, error)
	GetPartnerUsers(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerUser, error)
	UpdatePartnerUser(ctx context.Context, tenantID string, user *models.PartnerUser) (*models.PartnerUser, error)
	UpdatePartnerUserPassword(ctx context.Context, tenantID string, userID int64, newPasswordHash string) error
	DeactivatePartnerUser(ctx context.Context, tenantID string, userID int64) error

	// Partner Statistics
	GetPartnerStats(ctx context.Context, tenantID string, partnerID int64) (*models.PartnerStats, error)
	GetPartnerMonthlyStats(ctx context.Context, tenantID string, partnerID int64, year int, month int) (*models.PartnerStats, error)

	// Quality Scoring
	CalculateLeadQualityScore(ctx context.Context, leadData *models.LeadData) float64
}

type partnerService struct {
	db *sql.DB
}

// NewPartnerService creates a new partner service
func NewPartnerService(db *sql.DB) PartnerService {
	return &partnerService{db: db}
}

// CreatePartner creates a new partner
func (s *partnerService) CreatePartner(ctx context.Context, tenantID string, partner *models.Partner) (*models.Partner, error) {
	if partner.OrganizationName == "" || partner.ContactEmail == "" {
		return nil, errors.New("organization name and contact email are required")
	}

	// Generate partner code if not provided
	if partner.PartnerCode == "" {
		partner.PartnerCode = fmt.Sprintf("PARTNER_%d_%d", time.Now().Unix(), int64(partner.ID))
	}

	// Set defaults
	if partner.Status == "" {
		partner.Status = models.PartnerStatusPending
	}
	partner.TenantID = tenantID
	partner.CreatedAt = time.Now()
	partner.UpdatedAt = time.Now()

	query := `
		INSERT INTO partners (
			tenant_id, partner_code, organization_name, partner_type, status,
			contact_email, contact_phone, contact_person, website, description,
			address, city, state, country, zip_code, tax_id,
			banking_details, commission_percentage, lead_price, monthly_quota,
			document_urls, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		partner.TenantID, partner.PartnerCode, partner.OrganizationName, partner.PartnerType, partner.Status,
		partner.ContactEmail, partner.ContactPhone, partner.ContactPerson, partner.Website, partner.Description,
		partner.Address, partner.City, partner.State, partner.Country, partner.ZipCode, partner.TaxID,
		partner.BankingDetails, partner.CommissionPercentage, partner.LeadPrice, partner.MonthlyQuota,
		partner.DocumentURLs, partner.CreatedBy, partner.CreatedAt, partner.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create partner: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get partner ID: %w", err)
	}

	partner.ID = id
	return partner, nil
}

// GetPartner retrieves a partner by ID
func (s *partnerService) GetPartner(ctx context.Context, tenantID string, partnerID int64) (*models.Partner, error) {
	partner := &models.Partner{}

	query := `
		SELECT id, tenant_id, partner_code, organization_name, partner_type, status,
		contact_email, contact_phone, contact_person, website, description,
		address, city, state, country, zip_code, tax_id,
		banking_details, commission_percentage, lead_price, monthly_quota,
		current_month_leads, total_leads_submitted, approved_leads, rejected_leads,
		converted_leads, total_earnings, pending_payout_amount, withdrawn_amount,
		available_balance, approved_by, approved_at, rejection_reason,
		suspension_reason, suspended_at, document_urls,
		created_by, created_at, updated_at, deleted_at
		FROM partners
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := s.db.QueryRowContext(ctx, query, partnerID, tenantID).Scan(
		&partner.ID, &partner.TenantID, &partner.PartnerCode, &partner.OrganizationName, &partner.PartnerType, &partner.Status,
		&partner.ContactEmail, &partner.ContactPhone, &partner.ContactPerson, &partner.Website, &partner.Description,
		&partner.Address, &partner.City, &partner.State, &partner.Country, &partner.ZipCode, &partner.TaxID,
		&partner.BankingDetails, &partner.CommissionPercentage, &partner.LeadPrice, &partner.MonthlyQuota,
		&partner.CurrentMonthLeads, &partner.TotalLeadsSubmitted, &partner.ApprovedLeads, &partner.RejectedLeads,
		&partner.ConvertedLeads, &partner.TotalEarnings, &partner.PendingPayoutAmount, &partner.WithdrawnAmount,
		&partner.AvailableBalance, &partner.ApprovedBy, &partner.ApprovedAt, &partner.RejectionReason,
		&partner.SuspensionReason, &partner.SuspendedAt, &partner.DocumentURLs,
		&partner.CreatedBy, &partner.CreatedAt, &partner.UpdatedAt, &partner.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("partner not found")
		}
		return nil, fmt.Errorf("failed to get partner: %w", err)
	}

	return partner, nil
}

// GetPartnerByCode retrieves a partner by code
func (s *partnerService) GetPartnerByCode(ctx context.Context, tenantID string, partnerCode string) (*models.Partner, error) {
	partner := &models.Partner{}

	query := `
		SELECT id, tenant_id, partner_code, organization_name, partner_type, status,
		contact_email, contact_phone, contact_person, website, description,
		address, city, state, country, zip_code, tax_id,
		banking_details, commission_percentage, lead_price, monthly_quota,
		current_month_leads, total_leads_submitted, approved_leads, rejected_leads,
		converted_leads, total_earnings, pending_payout_amount, withdrawn_amount,
		available_balance, approved_by, approved_at, rejection_reason,
		suspension_reason, suspended_at, document_urls,
		created_by, created_at, updated_at, deleted_at
		FROM partners
		WHERE partner_code = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := s.db.QueryRowContext(ctx, query, partnerCode, tenantID).Scan(
		&partner.ID, &partner.TenantID, &partner.PartnerCode, &partner.OrganizationName, &partner.PartnerType, &partner.Status,
		&partner.ContactEmail, &partner.ContactPhone, &partner.ContactPerson, &partner.Website, &partner.Description,
		&partner.Address, &partner.City, &partner.State, &partner.Country, &partner.ZipCode, &partner.TaxID,
		&partner.BankingDetails, &partner.CommissionPercentage, &partner.LeadPrice, &partner.MonthlyQuota,
		&partner.CurrentMonthLeads, &partner.TotalLeadsSubmitted, &partner.ApprovedLeads, &partner.RejectedLeads,
		&partner.ConvertedLeads, &partner.TotalEarnings, &partner.PendingPayoutAmount, &partner.WithdrawnAmount,
		&partner.AvailableBalance, &partner.ApprovedBy, &partner.ApprovedAt, &partner.RejectionReason,
		&partner.SuspensionReason, &partner.SuspendedAt, &partner.DocumentURLs,
		&partner.CreatedBy, &partner.CreatedAt, &partner.UpdatedAt, &partner.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("partner not found")
		}
		return nil, fmt.Errorf("failed to get partner: %w", err)
	}

	return partner, nil
}

// GetPartners retrieves partners with filtering
func (s *partnerService) GetPartners(ctx context.Context, tenantID string, filter *models.PartnerFilter) ([]models.Partner, int64, error) {
	if filter == nil {
		filter = &models.PartnerFilter{Limit: 50}
	}
	if filter.Limit == 0 {
		filter.Limit = 50
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	partners := []models.Partner{}
	query := `SELECT id, tenant_id, partner_code, organization_name, partner_type, status,
		contact_email, contact_phone, contact_person, website, description,
		address, city, state, country, zip_code, tax_id,
		banking_details, commission_percentage, lead_price, monthly_quota,
		current_month_leads, total_leads_submitted, approved_leads, rejected_leads,
		converted_leads, total_earnings, pending_payout_amount, withdrawn_amount,
		available_balance, approved_by, approved_at, rejection_reason,
		suspension_reason, suspended_at, document_urls,
		created_by, created_at, updated_at, deleted_at
		FROM partners WHERE tenant_id = ? AND deleted_at IS NULL`

	args := []interface{}{tenantID}

	if filter.PartnerType != "" {
		query += " AND partner_type = ?"
		args = append(args, filter.PartnerType)
	}
	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}
	if filter.Search != "" {
		query += " AND (organization_name LIKE ? OR partner_code LIKE ? OR contact_email LIKE ?)"
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	// Get count
	countQuery := strings.Replace(query, "SELECT id, tenant_id, partner_code, organization_name, partner_type, status, contact_email, contact_phone, contact_person, website, description, address, city, state, country, zip_code, tax_id, banking_details, commission_percentage, lead_price, monthly_quota, current_month_leads, total_leads_submitted, approved_leads, rejected_leads, converted_leads, total_earnings, pending_payout_amount, withdrawn_amount, available_balance, approved_by, approved_at, rejection_reason, suspension_reason, suspended_at, document_urls, created_by, created_at, updated_at, deleted_at", "SELECT COUNT(*)", 1)
	var total int64
	s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get partners: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Partner
		err := rows.Scan(
			&p.ID, &p.TenantID, &p.PartnerCode, &p.OrganizationName, &p.PartnerType, &p.Status,
			&p.ContactEmail, &p.ContactPhone, &p.ContactPerson, &p.Website, &p.Description,
			&p.Address, &p.City, &p.State, &p.Country, &p.ZipCode, &p.TaxID,
			&p.BankingDetails, &p.CommissionPercentage, &p.LeadPrice, &p.MonthlyQuota,
			&p.CurrentMonthLeads, &p.TotalLeadsSubmitted, &p.ApprovedLeads, &p.RejectedLeads,
			&p.ConvertedLeads, &p.TotalEarnings, &p.PendingPayoutAmount, &p.WithdrawnAmount,
			&p.AvailableBalance, &p.ApprovedBy, &p.ApprovedAt, &p.RejectionReason,
			&p.SuspensionReason, &p.SuspendedAt, &p.DocumentURLs,
			&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan partner: %w", err)
		}
		partners = append(partners, p)
	}

	return partners, total, rows.Err()
}

// UpdatePartner updates a partner
func (s *partnerService) UpdatePartner(ctx context.Context, tenantID string, partner *models.Partner) (*models.Partner, error) {
	partner.UpdatedAt = time.Now()

	query := `
		UPDATE partners SET
			organization_name = ?, contact_email = ?, contact_phone = ?, contact_person = ?,
			website = ?, description = ?, address = ?, city = ?, state = ?, country = ?,
			zip_code = ?, tax_id = ?, banking_details = ?, commission_percentage = ?,
			lead_price = ?, monthly_quota = ?, document_urls = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		partner.OrganizationName, partner.ContactEmail, partner.ContactPhone, partner.ContactPerson,
		partner.Website, partner.Description, partner.Address, partner.City, partner.State, partner.Country,
		partner.ZipCode, partner.TaxID, partner.BankingDetails, partner.CommissionPercentage,
		partner.LeadPrice, partner.MonthlyQuota, partner.DocumentURLs, partner.UpdatedAt,
		partner.ID, tenantID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update partner: %w", err)
	}

	return partner, nil
}

// UpdatePartnerStatus updates partner status
func (s *partnerService) UpdatePartnerStatus(ctx context.Context, tenantID string, partnerID int64, status models.PartnerStatus, reason string, approvedBy int64) error {
	query := `
		UPDATE partners SET
			status = ?, rejection_reason = ?, approved_by = ?, approved_at = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, query,
		status, reason, approvedBy, now, now,
		partnerID, tenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update partner status: %w", err)
	}

	return nil
}

// DeactivatePartner deactivates a partner
func (s *partnerService) DeactivatePartner(ctx context.Context, tenantID string, partnerID int64, reason string) error {
	return s.UpdatePartnerStatus(ctx, tenantID, partnerID, models.PartnerStatusInactive, reason, 0)
}

// SuspendPartner suspends a partner
func (s *partnerService) SuspendPartner(ctx context.Context, tenantID string, partnerID int64, reason string) error {
	query := `
		UPDATE partners SET
			status = ?, suspension_reason = ?, suspended_at = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, query,
		models.PartnerStatusSuspended, reason, now, now,
		partnerID, tenantID,
	)

	return err
}

// CreatePartnerUser creates a new partner user
func (s *partnerService) CreatePartnerUser(ctx context.Context, tenantID string, user *models.PartnerUser) (*models.PartnerUser, error) {
	if user.Email == "" || user.PartnerID == 0 {
		return nil, errors.New("email and partner_id are required")
	}

	user.TenantID = tenantID
	user.IsActive = true
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO partner_users (
			partner_id, tenant_id, email, first_name, last_name, phone,
			password_hash, role, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		user.PartnerID, user.TenantID, user.Email, user.FirstName, user.LastName, user.Phone,
		user.PasswordHash, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create partner user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	user.ID = id
	return user, nil
}

// GetPartnerUser retrieves a partner user by ID
func (s *partnerService) GetPartnerUser(ctx context.Context, tenantID string, userID int64) (*models.PartnerUser, error) {
	user := &models.PartnerUser{}

	query := `
		SELECT id, partner_id, tenant_id, email, first_name, last_name, phone,
		password_hash, role, is_active, last_login, created_at, updated_at, deleted_at
		FROM partner_users
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := s.db.QueryRowContext(ctx, query, userID, tenantID).Scan(
		&user.ID, &user.PartnerID, &user.TenantID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
		&user.PasswordHash, &user.Role, &user.IsActive, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get partner user: %w", err)
	}

	return user, nil
}

// GetPartnerUserByEmail retrieves a partner user by email
func (s *partnerService) GetPartnerUserByEmail(ctx context.Context, tenantID string, email string) (*models.PartnerUser, error) {
	user := &models.PartnerUser{}

	query := `
		SELECT id, partner_id, tenant_id, email, first_name, last_name, phone,
		password_hash, role, is_active, last_login, created_at, updated_at, deleted_at
		FROM partner_users
		WHERE email = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := s.db.QueryRowContext(ctx, query, email, tenantID).Scan(
		&user.ID, &user.PartnerID, &user.TenantID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
		&user.PasswordHash, &user.Role, &user.IsActive, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get partner user: %w", err)
	}

	return user, nil
}

// GetPartnerUsers retrieves all users for a partner
func (s *partnerService) GetPartnerUsers(ctx context.Context, tenantID string, partnerID int64) ([]models.PartnerUser, error) {
	users := []models.PartnerUser{}

	query := `
		SELECT id, partner_id, tenant_id, email, first_name, last_name, phone,
		password_hash, role, is_active, last_login, created_at, updated_at, deleted_at
		FROM partner_users
		WHERE partner_id = ? AND tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, partnerID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get partner users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u models.PartnerUser
		err := rows.Scan(
			&u.ID, &u.PartnerID, &u.TenantID, &u.Email, &u.FirstName, &u.LastName, &u.Phone,
			&u.PasswordHash, &u.Role, &u.IsActive, &u.LastLogin, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan partner user: %w", err)
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

// UpdatePartnerUser updates a partner user
func (s *partnerService) UpdatePartnerUser(ctx context.Context, tenantID string, user *models.PartnerUser) (*models.PartnerUser, error) {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE partner_users SET
			email = ?, first_name = ?, last_name = ?, phone = ?,
			role = ?, is_active = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		user.Email, user.FirstName, user.LastName, user.Phone,
		user.Role, user.IsActive, user.UpdatedAt,
		user.ID, tenantID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update partner user: %w", err)
	}

	return user, nil
}

// UpdatePartnerUserPassword updates a partner user password
func (s *partnerService) UpdatePartnerUserPassword(ctx context.Context, tenantID string, userID int64, newPasswordHash string) error {
	query := `
		UPDATE partner_users SET
			password_hash = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		newPasswordHash, time.Now(),
		userID, tenantID,
	)

	return err
}

// DeactivatePartnerUser deactivates a partner user
func (s *partnerService) DeactivatePartnerUser(ctx context.Context, tenantID string, userID int64) error {
	query := `
		UPDATE partner_users SET
			is_active = ?, deleted_at = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, query,
		false, now, now,
		userID, tenantID,
	)

	return err
}

// GetPartnerStats retrieves statistics for a partner
func (s *partnerService) GetPartnerStats(ctx context.Context, tenantID string, partnerID int64) (*models.PartnerStats, error) {
	stats := &models.PartnerStats{}

	query := `
		SELECT
			total_leads_submitted, approved_leads, rejected_leads, converted_leads,
			total_earnings, available_balance, current_month_leads, monthly_quota,
			pending_payout_amount
		FROM partners
		WHERE id = ? AND tenant_id = ?
	`

	err := s.db.QueryRowContext(ctx, query, partnerID, tenantID).Scan(
		&stats.TotalLeadsSubmitted, &stats.ApprovedLeads, &stats.RejectedLeads, &stats.ConvertedLeads,
		&stats.TotalEarnings, &stats.AvailableBalance, &stats.CurrentMonthLeads, &stats.MonthlyQuota,
		&stats.PendingPayout,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get partner stats: %w", err)
	}

	// Calculate rates
	if stats.TotalLeadsSubmitted > 0 {
		stats.ApprovalRate = (float64(stats.ApprovedLeads) / float64(stats.TotalLeadsSubmitted)) * 100
		stats.ConversionRate = (float64(stats.ConvertedLeads) / float64(stats.ApprovedLeads)) * 100
	}

	// Get average lead quality
	avgQualityQuery := `
		SELECT COALESCE(AVG(quality_score), 0) FROM partner_leads
		WHERE partner_id = ? AND tenant_id = ? AND deleted_at IS NULL
	`
	s.db.QueryRowContext(ctx, avgQualityQuery, partnerID, tenantID).Scan(&stats.AverageLeadQuality)

	return stats, nil
}

// GetPartnerMonthlyStats retrieves monthly statistics for a partner
func (s *partnerService) GetPartnerMonthlyStats(ctx context.Context, tenantID string, partnerID int64, year int, month int) (*models.PartnerStats, error) {
	stats := &models.PartnerStats{}

	startDate := fmt.Sprintf("%04d-%02d-01", year, month)
	endDate := fmt.Sprintf("%04d-%02d-31", year, month)

	query := `
		SELECT
			COUNT(*) as total_leads,
			SUM(CASE WHEN status = 'approved' THEN 1 ELSE 0 END) as approved_leads,
			SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END) as rejected_leads,
			SUM(CASE WHEN status = 'converted' THEN 1 ELSE 0 END) as converted_leads,
			AVG(quality_score) as avg_quality
		FROM partner_leads
		WHERE partner_id = ? AND tenant_id = ? AND created_at >= ? AND created_at <= ?
	`

	err := s.db.QueryRowContext(ctx, query, partnerID, tenantID, startDate, endDate).Scan(
		&stats.TotalLeadsSubmitted, &stats.ApprovedLeads, &stats.RejectedLeads,
		&stats.ConvertedLeads, &stats.AverageLeadQuality,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get monthly stats: %w", err)
	}

	if stats.TotalLeadsSubmitted > 0 {
		stats.ApprovalRate = (float64(stats.ApprovedLeads) / float64(stats.TotalLeadsSubmitted)) * 100
		stats.ConversionRate = (float64(stats.ConvertedLeads) / float64(stats.ApprovedLeads)) * 100
	}

	return stats, nil
}

// CalculateLeadQualityScore calculates quality score for a lead (0-100)
func (s *partnerService) CalculateLeadQualityScore(ctx context.Context, leadData *models.LeadData) float64 {
	score := 0.0

	// Check required fields (30 points)
	if leadData.FirstName != "" {
		score += 5
	}
	if leadData.LastName != "" {
		score += 5
	}
	if leadData.Email != "" {
		score += 10
	}
	if leadData.Phone != "" {
		score += 10
	}

	// Check additional fields (40 points)
	if leadData.Company != "" {
		score += 10
	}
	if leadData.JobTitle != "" {
		score += 10
	}
	if leadData.Industry != "" {
		score += 10
	}
	if leadData.Address != "" {
		score += 10
	}

	// Check lead qualification fields (30 points)
	if leadData.LeadType != "" && (leadData.LeadType == "warm_lead" || leadData.LeadType == "prospect") {
		score += 10
	}
	if leadData.BudgetRange != "" {
		score += 10
	}
	if leadData.TimelineDays > 0 && leadData.TimelineDays <= 180 {
		score += 10
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	// Round to 2 decimal places
	score = math.Round(score*100) / 100

	return score
}
