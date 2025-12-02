package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
)

// GLService handles General Ledger operations
type GLService struct {
	DB *sql.DB
}

// NewGLService creates a new GL service
func NewGLService(db *sql.DB) *GLService {
	return &GLService{DB: db}
}

// ============================================================================
// CHART OF ACCOUNTS MANAGEMENT
// ============================================================================

// CreateAccount creates a new account in chart of accounts
func (s *GLService) CreateAccount(tenantID string, account *models.ChartOfAccount) error {
	account.TenantID = tenantID
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	query := `INSERT INTO chart_of_accounts (
		id, tenant_id, account_code, account_name, account_type, sub_account_type,
		parent_account_id, description, opening_balance, current_balance, is_active,
		is_header, is_default, currency, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		account.ID, account.TenantID, account.AccountCode, account.AccountName, account.AccountType,
		account.SubAccountType, account.ParentAccountID, account.Description, account.OpeningBalance,
		account.CurrentBalance, account.IsActive, account.IsHeader, account.IsDefault, account.Currency,
		account.CreatedAt, account.UpdatedAt,
	)

	return err
}

// GetAccount retrieves an account by ID
func (s *GLService) GetAccount(tenantID, accountID string) (*models.ChartOfAccount, error) {
	var account models.ChartOfAccount
	query := `SELECT id, tenant_id, account_code, account_name, account_type, sub_account_type,
		parent_account_id, description, opening_balance, current_balance, is_active, is_header,
		is_default, currency, created_at, updated_at, deleted_at
		FROM chart_of_accounts WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	err := s.DB.QueryRow(query, accountID, tenantID).Scan(
		&account.ID, &account.TenantID, &account.AccountCode, &account.AccountName, &account.AccountType,
		&account.SubAccountType, &account.ParentAccountID, &account.Description, &account.OpeningBalance,
		&account.CurrentBalance, &account.IsActive, &account.IsHeader, &account.IsDefault, &account.Currency,
		&account.CreatedAt, &account.UpdatedAt, &account.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	return &account, err
}

// ListAccounts retrieves all accounts
func (s *GLService) ListAccounts(tenantID string, accountType string) ([]models.ChartOfAccount, error) {
	var accounts []models.ChartOfAccount

	var query string
	var args []interface{}

	if accountType != "" {
		query = `SELECT id, tenant_id, account_code, account_name, account_type, sub_account_type,
			parent_account_id, description, opening_balance, current_balance, is_active, is_header,
			is_default, currency, created_at, updated_at, deleted_at
			FROM chart_of_accounts WHERE tenant_id = ? AND account_type = ? AND deleted_at IS NULL
			ORDER BY account_code ASC`
		args = []interface{}{tenantID, accountType}
	} else {
		query = `SELECT id, tenant_id, account_code, account_name, account_type, sub_account_type,
			parent_account_id, description, opening_balance, current_balance, is_active, is_header,
			is_default, currency, created_at, updated_at, deleted_at
			FROM chart_of_accounts WHERE tenant_id = ? AND deleted_at IS NULL
			ORDER BY account_code ASC`
		args = []interface{}{tenantID}
	}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var acc models.ChartOfAccount
		err := rows.Scan(
			&acc.ID, &acc.TenantID, &acc.AccountCode, &acc.AccountName, &acc.AccountType,
			&acc.SubAccountType, &acc.ParentAccountID, &acc.Description, &acc.OpeningBalance,
			&acc.CurrentBalance, &acc.IsActive, &acc.IsHeader, &acc.IsDefault, &acc.Currency,
			&acc.CreatedAt, &acc.UpdatedAt, &acc.DeletedAt,
		)
		if err != nil {
			log.Printf("Error scanning account: %v", err)
			continue
		}
		accounts = append(accounts, acc)
	}

	return accounts, rows.Err()
}

// ============================================================================
// JOURNAL ENTRY MANAGEMENT
// ============================================================================

// CreateJournalEntry creates a new journal entry
func (s *GLService) CreateJournalEntry(tenantID string, entry *models.JournalEntry) error {
	entry.TenantID = tenantID
	entry.EntryStatus = "Draft"
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()

	query := `INSERT INTO journal_entries (
		id, tenant_id, entry_date, reference_number, reference_type, reference_id,
		description, amount, narration, entry_status, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		entry.ID, entry.TenantID, entry.EntryDate, entry.ReferenceNumber, entry.ReferenceType,
		entry.ReferenceID, entry.Description, entry.Amount, entry.Narration, entry.EntryStatus,
		entry.CreatedAt, entry.UpdatedAt,
	)

	return err
}

// AddJournalEntryDetail adds a debit/credit line to an entry
func (s *GLService) AddJournalEntryDetail(detail *models.JournalEntryDetail) error {
	query := `INSERT INTO journal_entry_details (
		id, tenant_id, journal_entry_id, account_id, account_code, debit_amount, credit_amount,
		description, line_number, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		detail.ID, detail.TenantID, detail.JournalEntryID, detail.AccountID, detail.AccountCode,
		detail.DebitAmount, detail.CreditAmount, detail.Description, detail.LineNumber,
		detail.CreatedAt, detail.UpdatedAt,
	)

	return err
}

// PostJournalEntry posts a draft entry (moves from Draft to Posted)
func (s *GLService) PostJournalEntry(tenantID, entryID, postedBy string) error {
	// Validate debit/credit balance
	var totalDebit, totalCredit float64

	query := `SELECT SUM(debit_amount) as debit, SUM(credit_amount) as credit
		FROM journal_entry_details WHERE journal_entry_id = ? AND tenant_id = ?`

	err := s.DB.QueryRow(query, entryID, tenantID).Scan(&totalDebit, &totalCredit)
	if err != nil {
		return fmt.Errorf("failed to calculate totals: %v", err)
	}

	if totalDebit != totalCredit {
		return fmt.Errorf("journal entry is not balanced: debit %.2f != credit %.2f", totalDebit, totalCredit)
	}

	// Post the entry
	updateQuery := `UPDATE journal_entries SET entry_status = 'Posted', posted_by = ?, posted_at = NOW(), updated_at = NOW()
		WHERE id = ? AND tenant_id = ?`

	_, err = s.DB.Exec(updateQuery, postedBy, entryID, tenantID)
	if err != nil {
		return err
	}

	// Update account balances
	detailsQuery := `SELECT account_id, debit_amount, credit_amount FROM journal_entry_details
		WHERE journal_entry_id = ? AND tenant_id = ?`

	rows, err := s.DB.Query(detailsQuery, entryID, tenantID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var accountID string
		var debit, credit float64

		if err := rows.Scan(&accountID, &debit, &credit); err != nil {
			return err
		}

		// Update current balance
		balanceQuery := `UPDATE chart_of_accounts SET current_balance = current_balance + ? - ?
			WHERE id = ? AND tenant_id = ?`
		_, err = s.DB.Exec(balanceQuery, debit, credit, accountID, tenantID)
		if err != nil {
			return err
		}
	}

	return rows.Err()
}

// GetJournalEntry retrieves an entry with its details
func (s *GLService) GetJournalEntry(tenantID, entryID string) (*models.JournalEntry, error) {
	var entry models.JournalEntry

	query := `SELECT id, tenant_id, entry_date, reference_number, reference_type, reference_id,
		description, amount, narration, entry_status, posted_by, posted_at, created_at, updated_at, deleted_at
		FROM journal_entries WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	err := s.DB.QueryRow(query, entryID, tenantID).Scan(
		&entry.ID, &entry.TenantID, &entry.EntryDate, &entry.ReferenceNumber, &entry.ReferenceType,
		&entry.ReferenceID, &entry.Description, &entry.Amount, &entry.Narration, &entry.EntryStatus,
		&entry.PostedBy, &entry.PostedAt, &entry.CreatedAt, &entry.UpdatedAt, &entry.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("journal entry not found")
	}
	if err != nil {
		return nil, err
	}

	// Get details
	detailsQuery := `SELECT id, tenant_id, journal_entry_id, account_id, account_code, debit_amount,
		credit_amount, description, line_number, created_at, updated_at
		FROM journal_entry_details WHERE journal_entry_id = ? AND tenant_id = ?
		ORDER BY line_number ASC`

	rows, err := s.DB.Query(detailsQuery, entryID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var detail models.JournalEntryDetail
		err := rows.Scan(
			&detail.ID, &detail.TenantID, &detail.JournalEntryID, &detail.AccountID, &detail.AccountCode,
			&detail.DebitAmount, &detail.CreditAmount, &detail.Description, &detail.LineNumber,
			&detail.CreatedAt, &detail.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning detail: %v", err)
			continue
		}
		entry.Details = append(entry.Details, detail)
	}

	return &entry, rows.Err()
}

// ListJournalEntries retrieves entries for a date range
func (s *GLService) ListJournalEntries(tenantID string, fromDate, toDate time.Time) ([]models.JournalEntry, error) {
	var entries []models.JournalEntry

	query := `SELECT id, tenant_id, entry_date, reference_number, reference_type, reference_id,
		description, amount, narration, entry_status, posted_by, posted_at, created_at, updated_at, deleted_at
		FROM journal_entries WHERE tenant_id = ? AND entry_date BETWEEN ? AND ? AND deleted_at IS NULL
		ORDER BY entry_date DESC`

	rows, err := s.DB.Query(query, tenantID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.JournalEntry
		err := rows.Scan(
			&entry.ID, &entry.TenantID, &entry.EntryDate, &entry.ReferenceNumber, &entry.ReferenceType,
			&entry.ReferenceID, &entry.Description, &entry.Amount, &entry.Narration, &entry.EntryStatus,
			&entry.PostedBy, &entry.PostedAt, &entry.CreatedAt, &entry.UpdatedAt, &entry.DeletedAt,
		)
		if err != nil {
			log.Printf("Error scanning entry: %v", err)
			continue
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

// ============================================================================
// REPORTING
// ============================================================================

// GetTrialBalance calculates trial balance for a period
func (s *GLService) GetTrialBalance(tenantID string, periodStart, periodEnd time.Time) ([]models.TrialBalance, error) {
	var balances []models.TrialBalance

	query := `SELECT DISTINCT
		CONCAT(YEAR(je.entry_date), '-', MONTH(je.entry_date)) as period_id,
		coa.id as account_id,
		coa.account_code,
		coa.account_name,
		COALESCE(SUM(CASE WHEN jed.debit_amount > 0 THEN jed.debit_amount ELSE 0 END), 0) as debit_balance,
		COALESCE(SUM(CASE WHEN jed.credit_amount > 0 THEN jed.credit_amount ELSE 0 END), 0) as credit_balance
	FROM chart_of_accounts coa
	LEFT JOIN journal_entry_details jed ON coa.id = jed.account_id AND coa.tenant_id = jed.tenant_id
	LEFT JOIN journal_entries je ON jed.journal_entry_id = je.id AND je.entry_status = 'Posted'
	WHERE coa.tenant_id = ? AND (je.entry_date IS NULL OR (je.entry_date BETWEEN ? AND ?))
	GROUP BY coa.id, coa.account_code, coa.account_name
	ORDER BY coa.account_code ASC`

	rows, err := s.DB.Query(query, tenantID, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tb models.TrialBalance
		err := rows.Scan(&tb.PeriodID, &tb.AccountID, &tb.AccountCode, &tb.AccountName, &tb.DebitBalance, &tb.CreditBalance)
		if err != nil {
			log.Printf("Error scanning trial balance: %v", err)
			continue
		}
		balances = append(balances, tb)
	}

	return balances, rows.Err()
}

// GetAccountLedger retrieves all transactions for an account
func (s *GLService) GetAccountLedger(tenantID, accountID string, fromDate, toDate time.Time) ([]struct {
	EntryDate time.Time
	Debit     float64
	Credit    float64
	Balance   float64
	Reference string
}, error) {
	var ledger []struct {
		EntryDate time.Time
		Debit     float64
		Credit    float64
		Balance   float64
		Reference string
	}

	query := `SELECT je.entry_date, jed.debit_amount, jed.credit_amount, 0 as balance, je.reference_number
		FROM journal_entry_details jed
		JOIN journal_entries je ON jed.journal_entry_id = je.id
		WHERE jed.account_id = ? AND jed.tenant_id = ? AND je.entry_status = 'Posted'
			AND je.entry_date BETWEEN ? AND ?
		ORDER BY je.entry_date ASC`

	rows, err := s.DB.Query(query, accountID, tenantID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runningBalance float64
	for rows.Next() {
		var entry struct {
			EntryDate time.Time
			Debit     float64
			Credit    float64
			Balance   float64
			Reference string
		}

		err := rows.Scan(&entry.EntryDate, &entry.Debit, &entry.Credit, &entry.Balance, &entry.Reference)
		if err != nil {
			log.Printf("Error scanning ledger: %v", err)
			continue
		}

		runningBalance += entry.Debit - entry.Credit
		entry.Balance = runningBalance
		ledger = append(ledger, entry)
	}

	return ledger, rows.Err()
}

// ============================================================================
// PERIOD MANAGEMENT
// ============================================================================

// CreateFinancialPeriod creates a new financial period
func (s *GLService) CreateFinancialPeriod(tenantID string, period *models.FinancialPeriod) error {
	period.TenantID = tenantID
	period.CreatedAt = time.Now()
	period.UpdatedAt = time.Now()

	query := `INSERT INTO financial_periods (
		id, tenant_id, period_name, period_type, start_date, end_date, is_closed, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		period.ID, period.TenantID, period.PeriodName, period.PeriodType, period.StartDate,
		period.EndDate, period.IsClosed, period.CreatedAt, period.UpdatedAt,
	)

	return err
}

// GetFinancialPeriod retrieves a period
func (s *GLService) GetFinancialPeriod(tenantID, periodID string) (*models.FinancialPeriod, error) {
	var period models.FinancialPeriod

	query := `SELECT id, tenant_id, period_name, period_type, start_date, end_date, is_closed,
		closed_by, closed_at, created_at, updated_at, deleted_at
		FROM financial_periods WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	err := s.DB.QueryRow(query, periodID, tenantID).Scan(
		&period.ID, &period.TenantID, &period.PeriodName, &period.PeriodType, &period.StartDate,
		&period.EndDate, &period.IsClosed, &period.ClosedBy, &period.ClosedAt, &period.CreatedAt,
		&period.UpdatedAt, &period.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("financial period not found")
	}
	return &period, err
}

// ClosePeriod closes a financial period
func (s *GLService) ClosePeriod(tenantID, periodID, closedBy string) error {
	query := `UPDATE financial_periods SET is_closed = TRUE, closed_by = ?, closed_at = NOW(), updated_at = NOW()
		WHERE id = ? AND tenant_id = ?`

	_, err := s.DB.Exec(query, closedBy, periodID, tenantID)
	return err
}
