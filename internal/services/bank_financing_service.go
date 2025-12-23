package services

import (
	"context"
	"database/sql"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/logger"
)

// ============================================================================
// BANK FINANCING SERVICE
// ============================================================================

type BankFinancingService struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewBankFinancingService creates new bank financing service
func NewBankFinancingService(db *sql.DB, log *logger.Logger) *BankFinancingService {
	return &BankFinancingService{
		db:     db,
		logger: log,
	}
}

// ============================================================================
// BANK FINANCING METHODS
// ============================================================================

// CreateBankFinancing creates new financing record
func (s *BankFinancingService) CreateBankFinancing(ctx context.Context, financing *models.BankFinancing) (*models.BankFinancing, error) {
	financing.ID = generateBankFinancingUUID()
	financing.CreatedAt = time.Now()
	financing.UpdatedAt = time.Now()

	query := `
		INSERT INTO bank_financing (
			id, tenant_id, booking_id, bank_id, loan_amount, sanctioned_amount,
			disbursed_amount, outstanding_amount, loan_type, interest_rate,
			tenure_months, emi_amount, status, application_date, approval_date,
			sanction_date, expected_completion_date, application_ref_no,
			sanction_letter_url, created_by, created_at, updated_by, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		financing.ID, financing.TenantID, financing.BookingID, financing.BankID,
		financing.LoanAmount, financing.SanctionedAmount, financing.DisbursedAmount,
		financing.OutstandingAmount, financing.LoanType, financing.InterestRate,
		financing.TenureMonths, financing.EMIAmount, financing.Status,
		financing.ApplicationDate, financing.ApprovalDate, financing.SanctionDate,
		financing.ExpectedCompletionDate, financing.ApplicationRefNo,
		financing.SanctionLetterURL, financing.CreatedBy, financing.CreatedAt,
		financing.UpdatedBy, financing.UpdatedAt,
	)

	if err != nil {
		s.logger.Error("Failed to create bank financing", "error", err)
		return nil, err
	}

	return financing, nil
}

// GetBankFinancing retrieves a financing record
func (s *BankFinancingService) GetBankFinancing(ctx context.Context, tenantID, financingID string) (*models.BankFinancing, error) {
	query := `
		SELECT id, tenant_id, booking_id, bank_id, loan_amount, sanctioned_amount,
		       disbursed_amount, outstanding_amount, loan_type, interest_rate,
		       tenure_months, emi_amount, status, application_date, approval_date,
		       sanction_date, expected_completion_date, application_ref_no,
		       sanction_letter_url, created_by, created_at, updated_by, updated_at, deleted_at
		FROM bank_financing
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	financing := &models.BankFinancing{}
	err := s.db.QueryRowContext(ctx, query, financingID, tenantID).Scan(
		&financing.ID, &financing.TenantID, &financing.BookingID, &financing.BankID,
		&financing.LoanAmount, &financing.SanctionedAmount, &financing.DisbursedAmount,
		&financing.OutstandingAmount, &financing.LoanType, &financing.InterestRate,
		&financing.TenureMonths, &financing.EMIAmount, &financing.Status,
		&financing.ApplicationDate, &financing.ApprovalDate, &financing.SanctionDate,
		&financing.ExpectedCompletionDate, &financing.ApplicationRefNo,
		&financing.SanctionLetterURL, &financing.CreatedBy, &financing.CreatedAt,
		&financing.UpdatedBy, &financing.UpdatedAt, &financing.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		s.logger.Error("Failed to get bank financing", "error", err)
		return nil, err
	}

	return financing, nil
}

// ListBankFinancing lists financing records
func (s *BankFinancingService) ListBankFinancing(ctx context.Context, tenantID string) ([]models.BankFinancing, error) {
	query := `
		SELECT id, tenant_id, booking_id, bank_id, loan_amount, sanctioned_amount,
		       disbursed_amount, outstanding_amount, loan_type, interest_rate,
		       tenure_months, emi_amount, status, application_date, approval_date,
		       sanction_date, expected_completion_date, application_ref_no,
		       sanction_letter_url, created_by, created_at, updated_by, updated_at, deleted_at
		FROM bank_financing
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		s.logger.Error("Failed to list bank financing", "error", err)
		return nil, err
	}
	defer rows.Close()

	var financings []models.BankFinancing
	for rows.Next() {
		financing := models.BankFinancing{}
		err := rows.Scan(
			&financing.ID, &financing.TenantID, &financing.BookingID, &financing.BankID,
			&financing.LoanAmount, &financing.SanctionedAmount, &financing.DisbursedAmount,
			&financing.OutstandingAmount, &financing.LoanType, &financing.InterestRate,
			&financing.TenureMonths, &financing.EMIAmount, &financing.Status,
			&financing.ApplicationDate, &financing.ApprovalDate, &financing.SanctionDate,
			&financing.ExpectedCompletionDate, &financing.ApplicationRefNo,
			&financing.SanctionLetterURL, &financing.CreatedBy, &financing.CreatedAt,
			&financing.UpdatedBy, &financing.UpdatedAt, &financing.DeletedAt,
		)
		if err != nil {
			continue
		}
		financings = append(financings, financing)
	}

	return financings, nil
}

// ============================================================================
// BANK DISBURSEMENT METHODS
// ============================================================================

// CreateBankDisbursement creates new disbursement record
func (s *BankFinancingService) CreateBankDisbursement(ctx context.Context, disbursement *models.BankDisbursement) (*models.BankDisbursement, error) {
	disbursement.ID = generateBankFinancingUUID()
	disbursement.CreatedAt = time.Now()
	disbursement.UpdatedAt = time.Now()

	query := `
		INSERT INTO bank_disbursement (
			id, tenant_id, financing_id, disbursement_number, scheduled_amount,
			actual_amount, milestone_id, milestone_percentage, status,
			scheduled_date, actual_date, bank_reference_no, claim_document_url,
			release_approval_by, release_approval_date, created_by, created_at,
			updated_by, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		disbursement.ID, disbursement.TenantID, disbursement.FinancingID,
		disbursement.DisbursementNumber, disbursement.ScheduledAmount,
		disbursement.ActualAmount, disbursement.MilestoneID, disbursement.MilestonePercentage,
		disbursement.Status, disbursement.ScheduledDate, disbursement.ActualDate,
		disbursement.BankReferenceNo, disbursement.ClaimDocumentURL,
		disbursement.ReleaseApprovalBy, disbursement.ReleaseApprovalDate,
		disbursement.CreatedBy, disbursement.CreatedAt, disbursement.UpdatedBy, disbursement.UpdatedAt,
	)

	if err != nil {
		s.logger.Error("Failed to create bank disbursement", "error", err)
		return nil, err
	}

	return disbursement, nil
}

// ============================================================================
// BANK NOC METHODS
// ============================================================================

// CreateBankNOC creates new NOC record
func (s *BankFinancingService) CreateBankNOC(ctx context.Context, noc *models.BankNOC) (*models.BankNOC, error) {
	noc.ID = generateBankFinancingUUID()
	noc.CreatedAt = time.Now()
	noc.UpdatedAt = time.Now()

	query := `
		INSERT INTO bank_noc (
			id, tenant_id, financing_id, noc_type, noc_request_date,
			noc_received_date, noc_document_url, noc_amount, status,
			issued_by_bank, valid_till_date, remarks, created_by, created_at,
			updated_by, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		noc.ID, noc.TenantID, noc.FinancingID, noc.NOCType, noc.NOCRequestDate,
		noc.NOCReceivedDate, noc.NOCDocumentURL, noc.NOCAmount, noc.Status,
		noc.IssuedByBank, noc.ValidTillDate, noc.Remarks, noc.CreatedBy,
		noc.CreatedAt, noc.UpdatedBy, noc.UpdatedAt,
	)

	if err != nil {
		s.logger.Error("Failed to create bank NOC", "error", err)
		return nil, err
	}

	return noc, nil
}

// ============================================================================
// BANK COLLECTION METHODS
// ============================================================================

// CreateBankCollection creates new collection record
func (s *BankFinancingService) CreateBankCollection(ctx context.Context, collection *models.BankCollectionTracking) (*models.BankCollectionTracking, error) {
	collection.ID = generateBankFinancingUUID()
	collection.CreatedAt = time.Now()
	collection.UpdatedAt = time.Now()

	query := `
		INSERT INTO bank_collection_tracking (
			id, tenant_id, financing_id, collection_type, collection_amount,
			collection_date, payment_mode, payment_reference_no, emi_month,
			emi_number, principal_amount, interest_amount, status,
			bank_confirmation_date, created_by, created_at, updated_by, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		collection.ID, collection.TenantID, collection.FinancingID,
		collection.CollectionType, collection.CollectionAmount, collection.CollectionDate,
		collection.PaymentMode, collection.PaymentReferenceNo, collection.EMIMonth,
		collection.EMINumber, collection.PrincipalAmount, collection.InterestAmount,
		collection.Status, collection.BankConfirmationDate, collection.CreatedBy,
		collection.CreatedAt, collection.UpdatedBy, collection.UpdatedAt,
	)

	if err != nil {
		s.logger.Error("Failed to create bank collection", "error", err)
		return nil, err
	}

	return collection, nil
}

// ============================================================================
// BANK MASTER METHODS
// ============================================================================

// CreateBank creates new bank record
func (s *BankFinancingService) CreateBank(ctx context.Context, bank *models.Bank) (*models.Bank, error) {
	bank.ID = generateBankFinancingUUID()
	bank.CreatedAt = time.Now()
	bank.UpdatedAt = time.Now()

	query := `
		INSERT INTO bank (
			id, tenant_id, bank_name, branch_name, ifsc_code, branch_contact,
			branch_email, relationship_manager_name, relationship_manager_phone,
			relationship_manager_email, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		bank.ID, bank.TenantID, bank.BankName, bank.BranchName, bank.IFSCCode,
		bank.BranchContact, bank.BranchEmail, bank.RelationshipManagerName,
		bank.RelationshipManagerPhone, bank.RelationshipManagerEmail,
		bank.Status, bank.CreatedAt, bank.UpdatedAt,
	)

	if err != nil {
		s.logger.Error("Failed to create bank", "error", err)
		return nil, err
	}

	return bank, nil
}

// ListBanks lists all bank records
func (s *BankFinancingService) ListBanks(ctx context.Context, tenantID string) ([]models.Bank, error) {
	query := `
		SELECT id, tenant_id, bank_name, branch_name, ifsc_code, branch_contact,
		       branch_email, relationship_manager_name, relationship_manager_phone,
		       relationship_manager_email, status, created_at, updated_at
		FROM bank
		WHERE tenant_id = ?
		ORDER BY bank_name
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		s.logger.Error("Failed to list banks", "error", err)
		return nil, err
	}
	defer rows.Close()

	var banks []models.Bank
	for rows.Next() {
		bank := models.Bank{}
		err := rows.Scan(
			&bank.ID, &bank.TenantID, &bank.BankName, &bank.BranchName, &bank.IFSCCode,
			&bank.BranchContact, &bank.BranchEmail, &bank.RelationshipManagerName,
			&bank.RelationshipManagerPhone, &bank.RelationshipManagerEmail,
			&bank.Status, &bank.CreatedAt, &bank.UpdatedAt,
		)
		if err != nil {
			continue
		}
		banks = append(banks, bank)
	}

	return banks, nil
}

// Helper function to generate UUID - using time-based UUID for now
// TODO: Replace with proper UUID generation (github.com/google/uuid)
func generateBankFinancingUUID() string {
	// This should be replaced with actual UUID generation
	return "uuid-" + time.Now().Format("20060102150405")
}
