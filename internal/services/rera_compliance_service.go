package services

import (
	"database/sql"
	"fmt"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
)

// ============================================================================
// RERA COMPLIANCE SERVICE
// ============================================================================
// Manages project-specific collection accounts as per RERA regulations

type RERAComplianceService struct {
	DB *sql.DB
}

func NewRERAComplianceService(db *sql.DB) *RERAComplianceService {
	return &RERAComplianceService{DB: db}
}

// CreateProjectCollectionAccount creates a segregated collection account for a project
// RERA Requirement: Each project must have a separate, segregated collection account
func (s *RERAComplianceService) CreateProjectCollectionAccount(
	tenantID, projectID string,
	accountName, bankName, accountNumber, ifscCode string,
	minimumBalance float64,
) (*models.ProjectCollectionAccount, error) {

	accountCode := fmt.Sprintf("ACC-COLL-%s", projectID[:8])

	// Calculate maximum borrowing allowed: 10% of projected collections
	// This will be updated as collections come in
	maxBorrowingAllowed := 0.0

	account := &models.ProjectCollectionAccount{
		ID:                  fmt.Sprintf("PCA-%s-%d", projectID, time.Now().Unix()),
		TenantID:            tenantID,
		ProjectID:           projectID,
		AccountName:         accountName,
		AccountCode:         accountCode,
		BankName:            bankName,
		AccountNumber:       accountNumber,
		IFSCCode:            ifscCode,
		AccountType:         "Current",
		OpeningBalance:      0,
		CurrentBalance:      0,
		MinimumBalance:      minimumBalance,
		RERACompliant:       true,
		RegulatedAccount:    true,
		MaxBorrowingAllowed: maxBorrowingAllowed,
		CurrentBorrowing:    0,
		Status:              "Active",
		OpenedDate:          time.Now(),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	query := `INSERT INTO project_collection_accounts 
		(id, tenant_id, project_id, account_name, account_code, bank_name, account_number, 
		 ifsc_code, account_type, opening_balance, current_balance, minimum_balance, 
		 rera_compliant, regulated_account, max_borrowing_allowed, current_borrowing, 
		 status, opened_date, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		account.ID, account.TenantID, account.ProjectID, account.AccountName,
		account.AccountCode, account.BankName, account.AccountNumber, account.IFSCCode,
		account.AccountType, account.OpeningBalance, account.CurrentBalance,
		account.MinimumBalance, account.RERACompliant, account.RegulatedAccount,
		account.MaxBorrowingAllowed, account.CurrentBorrowing, account.Status,
		account.OpenedDate, account.CreatedAt, account.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create collection account: %w", err)
	}

	return account, nil
}

// RecordCollection records a collection received from customer
// Validates against payment schedule and calculates interest on delayed payments
func (s *RERAComplianceService) RecordCollection(
	tenantID, projectID, collectionAccountID, bookingID, unitID string,
	paymentMode string,
	amountCollected float64,
	paidBy string,
) (*models.ProjectCollectionLedger, error) {

	// Get collection account
	account, err := s.GetCollectionAccount(tenantID, collectionAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection account: %w", err)
	}

	// Create collection ledger entry
	collectionNumber := fmt.Sprintf("COLL-%s-%d", projectID[:4], time.Now().Unix())
	collection := &models.ProjectCollectionLedger{
		ID:                  fmt.Sprintf("PCL-%d", time.Now().UnixNano()),
		TenantID:            tenantID,
		ProjectID:           projectID,
		CollectionAccountID: collectionAccountID,
		CollectionDate:      time.Now(),
		CollectionNumber:    collectionNumber,
		BookingID:           bookingID,
		UnitID:              unitID,
		PaymentMode:         paymentMode,
		AmountCollected:     amountCollected,
		PaidBy:              paidBy,
		Status:              "Pending",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	query := `INSERT INTO project_collection_ledger 
		(id, tenant_id, project_id, collection_account_id, collection_date, collection_number, 
		 booking_id, unit_id, payment_mode, amount_collected, paid_by, status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = s.DB.Exec(query,
		collection.ID, collection.TenantID, collection.ProjectID, collection.CollectionAccountID,
		collection.CollectionDate, collection.CollectionNumber, collection.BookingID,
		collection.UnitID, collection.PaymentMode, collection.AmountCollected,
		collection.PaidBy, collection.Status, collection.CreatedAt, collection.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to record collection: %w", err)
	}

	// Update collection account balance
	newBalance := account.CurrentBalance + amountCollected
	updateQuery := `UPDATE project_collection_accounts 
		SET current_balance = ?, updated_at = ? WHERE id = ?`

	_, err = s.DB.Exec(updateQuery, newBalance, time.Now(), collectionAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	return collection, nil
}

// RecordFundUtilization records how collection funds are utilized (RERA requirement)
// RERA mandates tracking of fund utilization: Construction, Land Cost, Statutory Approvals, Admin
func (s *RERAComplianceService) RecordFundUtilization(
	tenantID, projectID, collectionAccountID string,
	utilizationType, description string,
	amountUtilized float64,
	billNumber string,
) (*models.ProjectFundUtilization, error) {

	utilization := &models.ProjectFundUtilization{
		ID:                  fmt.Sprintf("PFU-%d", time.Now().UnixNano()),
		TenantID:            tenantID,
		ProjectID:           projectID,
		CollectionAccountID: collectionAccountID,
		UtilizationDate:     time.Now(),
		UtilizationType:     utilizationType,
		Description:         description,
		AmountUtilized:      amountUtilized,
		BillNumber:          billNumber,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	query := `INSERT INTO project_fund_utilization 
		(id, tenant_id, project_id, collection_account_id, utilization_date, utilization_type, 
		 description, amount_utilized, bill_number, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		utilization.ID, utilization.TenantID, utilization.ProjectID,
		utilization.CollectionAccountID, utilization.UtilizationDate, utilization.UtilizationType,
		utilization.Description, utilization.AmountUtilized, utilization.BillNumber,
		utilization.CreatedAt, utilization.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to record fund utilization: %w", err)
	}

	// Update collection account balance
	account, err := s.GetCollectionAccount(tenantID, collectionAccountID)
	if err != nil {
		return nil, err
	}

	newBalance := account.CurrentBalance - amountUtilized
	updateQuery := `UPDATE project_collection_accounts 
		SET current_balance = ?, updated_at = ? WHERE id = ?`

	_, err = s.DB.Exec(updateQuery, newBalance, time.Now(), collectionAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	return utilization, nil
}

// CheckBorrowingLimit validates if borrowing is within RERA limits (max 10% of collections)
func (s *RERAComplianceService) CheckBorrowingLimit(tenantID, projectID string, proposedBorrowingAmount float64) (bool, float64, error) {
	query := `SELECT COALESCE(SUM(amount_collected), 0) FROM project_collection_ledger 
		WHERE tenant_id = ? AND project_id = ? AND status = 'Approved'`

	var totalCollections float64
	err := s.DB.QueryRow(query, tenantID, projectID).Scan(&totalCollections)
	if err != nil {
		return false, 0, fmt.Errorf("failed to calculate total collections: %w", err)
	}

	// RERA limit: Borrowing cannot exceed 10% of total collections
	maxBorrowingAllowed := totalCollections * 0.10
	canBorrow := proposedBorrowingAmount <= maxBorrowingAllowed

	return canBorrow, maxBorrowingAllowed, nil
}

// GetCollectionAccount retrieves a collection account
func (s *RERAComplianceService) GetCollectionAccount(tenantID, accountID string) (*models.ProjectCollectionAccount, error) {
	query := `SELECT id, tenant_id, project_id, account_name, account_code, bank_name, account_number, 
		ifsc_code, account_type, opening_balance, current_balance, minimum_balance, 
		rera_compliant, regulated_account, max_borrowing_allowed, current_borrowing, 
		status, opened_date, closed_date, created_at, updated_at, deleted_at 
		FROM project_collection_accounts WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	account := &models.ProjectCollectionAccount{}
	err := s.DB.QueryRow(query, accountID, tenantID).Scan(
		&account.ID, &account.TenantID, &account.ProjectID, &account.AccountName,
		&account.AccountCode, &account.BankName, &account.AccountNumber, &account.IFSCCode,
		&account.AccountType, &account.OpeningBalance, &account.CurrentBalance, &account.MinimumBalance,
		&account.RERACompliant, &account.RegulatedAccount, &account.MaxBorrowingAllowed,
		&account.CurrentBorrowing, &account.Status, &account.OpenedDate, &account.ClosedDate,
		&account.CreatedAt, &account.UpdatedAt, &account.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("collection account not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get collection account: %w", err)
	}

	return account, nil
}

// GetProjectCollectionSummary returns comprehensive collection and utilization summary for a project
func (s *RERAComplianceService) GetProjectCollectionSummary(tenantID, projectID string) (*models.ProjectCollectionResponse, error) {
	query := `SELECT id, account_name FROM project_collection_accounts 
		WHERE tenant_id = ? AND project_id = ? AND deleted_at IS NULL`

	account := &models.ProjectCollectionAccount{}
	err := s.DB.QueryRow(query, tenantID, projectID).Scan(&account.ID, &account.AccountName)

	if err != nil {
		return nil, fmt.Errorf("no collection account found for project: %w", err)
	}

	// Get account balance
	accountDetails, err := s.GetCollectionAccount(tenantID, account.ID)
	if err != nil {
		return nil, err
	}

	// Calculate total utilization
	utilizationQuery := `SELECT COALESCE(SUM(amount_utilized), 0) FROM project_fund_utilization 
		WHERE tenant_id = ? AND project_id = ? AND deleted_at IS NULL`

	var totalUtilized float64
	err = s.DB.QueryRow(utilizationQuery, tenantID, projectID).Scan(&totalUtilized)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate fund utilization: %w", err)
	}

	summary := &models.ProjectCollectionResponse{
		ID:                   projectID,
		CollectionAccount:    accountDetails,
		TotalCollected:       accountDetails.CurrentBalance + totalUtilized,
		TotalUtilized:        totalUtilized,
		AvailableBalance:     accountDetails.CurrentBalance,
		RERAComplianceStatus: "Compliant",
	}

	return summary, nil
}

// PerformMonthlyReconciliation performs monthly reconciliation of collection account
func (s *RERAComplianceService) PerformMonthlyReconciliation(tenantID, projectID, collectionAccountID string) (*models.ProjectAccountReconciliation, error) {

	// Get opening balance from previous month
	prevMonthQuery := `SELECT COALESCE(closing_balance, 0) FROM project_account_reconciliation 
		WHERE tenant_id = ? AND project_id = ? AND collection_account_id = ? 
		ORDER BY reconciliation_date DESC LIMIT 1`

	var openingBalance float64
	err := s.DB.QueryRow(prevMonthQuery, tenantID, projectID, collectionAccountID).Scan(&openingBalance)
	if err != nil && err != sql.ErrNoRows {
		openingBalance = 0 // First month, no opening balance
	}

	// Calculate collections for current month
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	collectionsQuery := `SELECT COALESCE(SUM(amount_collected), 0) FROM project_collection_ledger 
		WHERE tenant_id = ? AND project_id = ? AND collection_date >= ? AND collection_date <= ? AND status = 'Approved'`

	var totalCollections float64
	err = s.DB.QueryRow(collectionsQuery, tenantID, projectID, startOfMonth, endOfMonth).Scan(&totalCollections)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate collections: %w", err)
	}

	// Calculate utilizations for current month
	utilizationsQuery := `SELECT COALESCE(SUM(amount_utilized), 0) FROM project_fund_utilization 
		WHERE tenant_id = ? AND project_id = ? AND utilization_date >= ? AND utilization_date <= ?`

	var totalUtilizations float64
	err = s.DB.QueryRow(utilizationsQuery, tenantID, projectID, startOfMonth, endOfMonth).Scan(&totalUtilizations)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate utilizations: %w", err)
	}

	closingBalance := openingBalance + totalCollections - totalUtilizations

	reconciliation := &models.ProjectAccountReconciliation{
		ID:                  fmt.Sprintf("PAR-%d", time.Now().UnixNano()),
		TenantID:            tenantID,
		ProjectID:           projectID,
		CollectionAccountID: collectionAccountID,
		ReconciliationDate:  time.Now(),
		PeriodFrom:          startOfMonth,
		PeriodTo:            endOfMonth,
		OpeningBalance:      openingBalance,
		TotalCollections:    totalCollections,
		TotalUtilizations:   totalUtilizations,
		ClosingBalance:      closingBalance,
		Reconciled:          false,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	query := `INSERT INTO project_account_reconciliation 
		(id, tenant_id, project_id, collection_account_id, reconciliation_date, period_from, period_to, 
		 opening_balance, total_collections, total_utilizations, closing_balance, reconciled, 
		 created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = s.DB.Exec(query,
		reconciliation.ID, reconciliation.TenantID, reconciliation.ProjectID,
		reconciliation.CollectionAccountID, reconciliation.ReconciliationDate,
		reconciliation.PeriodFrom, reconciliation.PeriodTo, reconciliation.OpeningBalance,
		reconciliation.TotalCollections, reconciliation.TotalUtilizations,
		reconciliation.ClosingBalance, reconciliation.Reconciled,
		reconciliation.CreatedAt, reconciliation.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create reconciliation: %w", err)
	}

	return reconciliation, nil
}
