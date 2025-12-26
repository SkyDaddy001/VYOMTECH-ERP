package services

import (
	"database/sql"
	"fmt"
	"time"

	"lms/cmd/api/models"
)

// ===== CHART OF ACCOUNTS SERVICE =====

type ChartOfAccountsService struct {
	db *sql.DB
}

func NewChartOfAccountsService(db *sql.DB) *ChartOfAccountsService {
	return &ChartOfAccountsService{db: db}
}

func (s *ChartOfAccountsService) CreateAccount(tenantID string, account *models.ChartOfAccounts) (string, error) {
	account.ID = generateID("coa")
	account.TenantID = tenantID
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()
	account.CurrentBalance = account.OpeningBalance
	account.IsActive = true

	query := `
		INSERT INTO chart_of_accounts (id, tenant_id, account_code, account_name, account_type, sub_type, 
			description, is_active, opening_balance, current_balance, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := s.db.Exec(query,
		account.ID, account.TenantID, account.AccountCode, account.AccountName, account.AccountType,
		account.SubType, account.Description, account.IsActive, account.OpeningBalance,
		account.CurrentBalance, account.Currency, account.CreatedAt, account.UpdatedAt,
	)

	return account.ID, err
}

func (s *ChartOfAccountsService) GetAccount(tenantID, accountID string) (*models.ChartOfAccounts, error) {
	account := &models.ChartOfAccounts{}
	query := `SELECT id, tenant_id, account_code, account_name, account_type, sub_type, 
		description, is_active, opening_balance, current_balance, currency, created_at, updated_at
		FROM chart_of_accounts WHERE id = $1 AND tenant_id = $2`

	err := s.db.QueryRow(query, accountID, tenantID).Scan(
		&account.ID, &account.TenantID, &account.AccountCode, &account.AccountName, &account.AccountType,
		&account.SubType, &account.Description, &account.IsActive, &account.OpeningBalance,
		&account.CurrentBalance, &account.Currency, &account.CreatedAt, &account.UpdatedAt,
	)

	return account, err
}

func (s *ChartOfAccountsService) GetAccountByCode(tenantID, accountCode string) (*models.ChartOfAccounts, error) {
	account := &models.ChartOfAccounts{}
	query := `SELECT id, tenant_id, account_code, account_name, account_type, sub_type, 
		description, is_active, opening_balance, current_balance, currency, created_at, updated_at
		FROM chart_of_accounts WHERE account_code = $1 AND tenant_id = $2`

	err := s.db.QueryRow(query, accountCode, tenantID).Scan(
		&account.ID, &account.TenantID, &account.AccountCode, &account.AccountName, &account.AccountType,
		&account.SubType, &account.Description, &account.IsActive, &account.OpeningBalance,
		&account.CurrentBalance, &account.Currency, &account.CreatedAt, &account.UpdatedAt,
	)

	return account, err
}

func (s *ChartOfAccountsService) ListAccounts(tenantID string, accountType models.ChartOfAccountsType, limit, offset int) ([]*models.ChartOfAccounts, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM chart_of_accounts WHERE tenant_id = $1`
	args := []interface{}{tenantID}

	if accountType != "" {
		countQuery += ` AND account_type = $2`
		args = append(args, accountType)
	}

	s.db.QueryRow(countQuery, args...).Scan(&total)

	query := `SELECT id, tenant_id, account_code, account_name, account_type, sub_type, 
		description, is_active, opening_balance, current_balance, currency, created_at, updated_at
		FROM chart_of_accounts WHERE tenant_id = $1`

	if accountType != "" {
		query += ` AND account_type = $2`
	}

	query += ` ORDER BY account_code LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()

	var accounts []*models.ChartOfAccounts
	for rows.Next() {
		account := &models.ChartOfAccounts{}
		err := rows.Scan(
			&account.ID, &account.TenantID, &account.AccountCode, &account.AccountName, &account.AccountType,
			&account.SubType, &account.Description, &account.IsActive, &account.OpeningBalance,
			&account.CurrentBalance, &account.Currency, &account.CreatedAt, &account.UpdatedAt,
		)
		if err == nil {
			accounts = append(accounts, account)
		}
	}

	return accounts, total, nil
}

// ===== JOURNAL ENTRY SERVICE =====

type JournalEntryService struct {
	db *sql.DB
}

func NewJournalEntryService(db *sql.DB) *JournalEntryService {
	return &JournalEntryService{db: db}
}

func (s *JournalEntryService) CreateJournalEntry(tenantID string, entry *models.JournalEntry, lines []*models.JournalEntryLine) (string, error) {
	entry.ID = generateID("je")
	entry.TenantID = tenantID
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	entry.Status = models.JournalEntryStatusDraft

	// Generate entry number
	var maxNumber int
	s.db.QueryRow(`SELECT COALESCE(MAX(CAST(SUBSTRING(entry_number, 3) AS INTEGER)), 0) FROM journal_entries WHERE tenant_id = $1`, tenantID).Scan(&maxNumber)
	entry.EntryNumber = fmt.Sprintf("JE%06d", maxNumber+1)

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Insert journal entry
	query := `
		INSERT INTO journal_entries (id, tenant_id, entry_number, entry_date, reference, description, status, prepared_by_id, total_debit, total_credit, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = tx.Exec(query,
		entry.ID, entry.TenantID, entry.EntryNumber, entry.EntryDate, entry.Reference, entry.Description,
		entry.Status, entry.PreparedByID, entry.TotalDebit, entry.TotalCredit, entry.CreatedAt, entry.UpdatedAt,
	)

	if err != nil {
		return "", err
	}

	// Insert lines
	for _, line := range lines {
		line.ID = generateID("jel")
		line.TenantID = tenantID
		line.JournalEntryID = entry.ID
		line.CreatedAt = time.Now()

		lineQuery := `
			INSERT INTO journal_entry_lines (id, tenant_id, journal_entry_id, account_id, account_code, account_name, 
				debit_amount, credit_amount, description, reference, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`

		_, err = tx.Exec(lineQuery,
			line.ID, line.TenantID, line.JournalEntryID, line.AccountID, line.AccountCode, line.AccountName,
			line.DebitAmount, line.CreditAmount, line.Description, line.Reference, line.CreatedAt,
		)

		if err != nil {
			return "", err
		}

		// Update chart of accounts balance
		if line.DebitAmount > 0 {
			_, err = tx.Exec(`UPDATE chart_of_accounts SET current_balance = current_balance + $1, updated_at = $2 WHERE id = $3`,
				line.DebitAmount, time.Now(), line.AccountID)
		} else if line.CreditAmount > 0 {
			_, err = tx.Exec(`UPDATE chart_of_accounts SET current_balance = current_balance - $1, updated_at = $2 WHERE id = $3`,
				line.CreditAmount, time.Now(), line.AccountID)
		}

		if err != nil {
			return "", err
		}
	}

	return entry.ID, tx.Commit()
}

func (s *JournalEntryService) GetJournalEntry(tenantID, entryID string) (*models.JournalEntry, []*models.JournalEntryLine, error) {
	entry := &models.JournalEntry{}
	query := `SELECT id, tenant_id, entry_number, entry_date, reference, description, status, prepared_by_id, approved_by_id, total_debit, total_credit, created_at, updated_at, posted_at
		FROM journal_entries WHERE id = $1 AND tenant_id = $2`

	err := s.db.QueryRow(query, entryID, tenantID).Scan(
		&entry.ID, &entry.TenantID, &entry.EntryNumber, &entry.EntryDate, &entry.Reference, &entry.Description,
		&entry.Status, &entry.PreparedByID, &entry.ApprovedByID, &entry.TotalDebit, &entry.TotalCredit,
		&entry.CreatedAt, &entry.UpdatedAt, &entry.PostedAt,
	)

	if err != nil {
		return nil, nil, err
	}

	// Get lines
	lineQuery := `SELECT id, tenant_id, journal_entry_id, account_id, account_code, account_name, 
		debit_amount, credit_amount, description, reference, created_at
		FROM journal_entry_lines WHERE journal_entry_id = $1 ORDER BY created_at`

	rows, err := s.db.Query(lineQuery, entryID)
	if err != nil {
		return entry, nil, err
	}
	defer rows.Close()

	var lines []*models.JournalEntryLine
	for rows.Next() {
		line := &models.JournalEntryLine{}
		err := rows.Scan(
			&line.ID, &line.TenantID, &line.JournalEntryID, &line.AccountID, &line.AccountCode, &line.AccountName,
			&line.DebitAmount, &line.CreditAmount, &line.Description, &line.Reference, &line.CreatedAt,
		)
		if err == nil {
			lines = append(lines, line)
		}
	}

	return entry, lines, nil
}

func (s *JournalEntryService) PostJournalEntry(tenantID, entryID, approvedByID string) error {
	query := `UPDATE journal_entries SET status = $1, approved_by_id = $2, posted_at = $3, updated_at = $4 
		WHERE id = $5 AND tenant_id = $6 AND status = $7`

	_, err := s.db.Exec(query,
		models.JournalEntryStatusPosted, approvedByID, time.Now(), time.Now(),
		entryID, tenantID, models.JournalEntryStatusDraft,
	)

	return err
}

func (s *JournalEntryService) ListJournalEntries(tenantID string, status models.JournalEntryStatus, limit, offset int) ([]*models.JournalEntry, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM journal_entries WHERE tenant_id = $1`
	args := []interface{}{tenantID}

	if status != "" {
		countQuery += ` AND status = $2`
		args = append(args, status)
	}

	s.db.QueryRow(countQuery, args...).Scan(&total)

	query := `SELECT id, tenant_id, entry_number, entry_date, reference, description, status, prepared_by_id, approved_by_id, total_debit, total_credit, created_at, updated_at, posted_at
		FROM journal_entries WHERE tenant_id = $1`

	if status != "" {
		query += ` AND status = $2`
	}

	query += ` ORDER BY entry_date DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()

	var entries []*models.JournalEntry
	for rows.Next() {
		entry := &models.JournalEntry{}
		err := rows.Scan(
			&entry.ID, &entry.TenantID, &entry.EntryNumber, &entry.EntryDate, &entry.Reference, &entry.Description,
			&entry.Status, &entry.PreparedByID, &entry.ApprovedByID, &entry.TotalDebit, &entry.TotalCredit,
			&entry.CreatedAt, &entry.UpdatedAt, &entry.PostedAt,
		)
		if err == nil {
			entries = append(entries, entry)
		}
	}

	return entries, total, nil
}

// ===== GENERAL LEDGER SERVICE =====

type GeneralLedgerService struct {
	db *sql.DB
}

func NewGeneralLedgerService(db *sql.DB) *GeneralLedgerService {
	return &GeneralLedgerService{db: db}
}

func (s *GeneralLedgerService) GetAccountLedger(tenantID, accountID string, fromDate, toDate time.Time, limit, offset int) ([]*models.GeneralLedger, error) {
	query := `SELECT id, tenant_id, account_id, account_code, journal_entry_id, entry_date, 
		reference, debit_amount, credit_amount, running_balance, description, created_at
		FROM general_ledger WHERE tenant_id = $1 AND account_id = $2 AND entry_date BETWEEN $3 AND $4
		ORDER BY entry_date, created_at LIMIT $5 OFFSET $6`

	rows, err := s.db.Query(query, tenantID, accountID, fromDate, toDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ledgers []*models.GeneralLedger
	for rows.Next() {
		ledger := &models.GeneralLedger{}
		err := rows.Scan(
			&ledger.ID, &ledger.TenantID, &ledger.AccountID, &ledger.AccountCode, &ledger.JournalEntryID,
			&ledger.EntryDate, &ledger.Reference, &ledger.DebitAmount, &ledger.CreditAmount,
			&ledger.RunningBalance, &ledger.Description, &ledger.CreatedAt,
		)
		if err == nil {
			ledgers = append(ledgers, ledger)
		}
	}

	return ledgers, nil
}

// ===== TRIAL BALANCE SERVICE =====

type TrialBalanceService struct {
	db *sql.DB
}

func NewTrialBalanceService(db *sql.DB) *TrialBalanceService {
	return &TrialBalanceService{db: db}
}

func (s *TrialBalanceService) GetTrialBalance(tenantID string, asOfDate time.Time) ([]*models.TrialBalance, error) {
	query := `SELECT coa.account_code, coa.account_name, coa.account_type, 
		SUM(CASE WHEN gl.debit_amount > 0 THEN gl.debit_amount ELSE 0 END) as debit_amount,
		SUM(CASE WHEN gl.credit_amount > 0 THEN gl.credit_amount ELSE 0 END) as credit_amount
		FROM chart_of_accounts coa
		LEFT JOIN general_ledger gl ON coa.id = gl.account_id AND gl.tenant_id = $1 AND gl.entry_date <= $2
		WHERE coa.tenant_id = $1 AND coa.is_active = true
		GROUP BY coa.id, coa.account_code, coa.account_name, coa.account_type
		ORDER BY coa.account_code`

	rows, err := s.db.Query(query, tenantID, asOfDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trialBalance []*models.TrialBalance
	for rows.Next() {
		tb := &models.TrialBalance{}
		err := rows.Scan(&tb.AccountCode, &tb.AccountName, &tb.AccountType, &tb.DebitAmount, &tb.CreditAmount)
		if err == nil {
			trialBalance = append(trialBalance, tb)
		}
	}

	return trialBalance, nil
}

// ===== ACCOUNTS RECEIVABLE SERVICE =====

type AccountsReceivableService struct {
	db *sql.DB
}

func NewAccountsReceivableService(db *sql.DB) *AccountsReceivableService {
	return &AccountsReceivableService{db: db}
}

func (s *AccountsReceivableService) CreateAR(tenantID string, ar *models.AccountsReceivable) (string, error) {
	ar.ID = generateID("ar")
	ar.TenantID = tenantID
	ar.CreatedAt = time.Now()
	ar.UpdatedAt = time.Now()

	query := `
		INSERT INTO accounts_receivable (id, tenant_id, customer_id, customer_name, invoice_id, invoice_number, 
			invoice_amount, amount_paid, outstanding_amount, invoice_date, due_date, status, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := s.db.Exec(query,
		ar.ID, ar.TenantID, ar.CustomerID, ar.CustomerName, ar.InvoiceID, ar.InvoiceNumber,
		ar.InvoiceAmount, ar.AmountPaid, ar.OutstandingAmount, ar.InvoiceDate, ar.DueDate,
		ar.Status, ar.Currency, ar.CreatedAt, ar.UpdatedAt,
	)

	return ar.ID, err
}

func (s *AccountsReceivableService) GetOutstandingAR(tenantID string) (float64, int64, error) {
	var total float64
	var count int64

	err := s.db.QueryRow(
		`SELECT SUM(outstanding_amount), COUNT(*) FROM accounts_receivable 
		WHERE tenant_id = $1 AND status NOT IN ('paid', 'written_off')`,
		tenantID,
	).Scan(&total, &count)

	if err != nil && err != sql.ErrNoRows {
		return 0, 0, err
	}

	return total, count, nil
}

// ===== ACCOUNTS PAYABLE SERVICE =====

type AccountsPayableService struct {
	db *sql.DB
}

func NewAccountsPayableService(db *sql.DB) *AccountsPayableService {
	return &AccountsPayableService{db: db}
}

func (s *AccountsPayableService) CreateAP(tenantID string, ap *models.AccountsPayable) (string, error) {
	ap.ID = generateID("ap")
	ap.TenantID = tenantID
	ap.CreatedAt = time.Now()
	ap.UpdatedAt = time.Now()

	query := `
		INSERT INTO accounts_payable (id, tenant_id, vendor_id, vendor_name, bill_id, bill_number, 
			bill_amount, amount_paid, outstanding_amount, bill_date, due_date, status, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := s.db.Exec(query,
		ap.ID, ap.TenantID, ap.VendorID, ap.VendorName, ap.BillID, ap.BillNumber,
		ap.BillAmount, ap.AmountPaid, ap.OutstandingAmount, ap.BillDate, ap.DueDate,
		ap.Status, ap.Currency, ap.CreatedAt, ap.UpdatedAt,
	)

	return ap.ID, err
}

func (s *AccountsPayableService) GetOutstandingAP(tenantID string) (float64, int64, error) {
	var total float64
	var count int64

	err := s.db.QueryRow(
		`SELECT SUM(outstanding_amount), COUNT(*) FROM accounts_payable 
		WHERE tenant_id = $1 AND status NOT IN ('paid', 'written_off')`,
		tenantID,
	).Scan(&total, &count)

	if err != nil && err != sql.ErrNoRows {
		return 0, 0, err
	}

	return total, count, nil
}

// ===== FINANCIAL STATEMENTS SERVICE =====

type FinancialStatementsService struct {
	db *sql.DB
}

func NewFinancialStatementsService(db *sql.DB) *FinancialStatementsService {
	return &FinancialStatementsService{db: db}
}

func (s *FinancialStatementsService) GetIncomeStatement(tenantID string, fromDate, toDate time.Time) (*models.FinancialStatement, error) {
	statement := &models.FinancialStatement{
		TenantID:      tenantID,
		StatementType: "income_statement",
		Period:        fmt.Sprintf("%s-%s", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")),
		Sections:      make([]models.FinancialStatementSection, 0),
	}

	// Get revenues
	revQuery := `SELECT coa.account_code, coa.account_name, SUM(gl.credit_amount - gl.debit_amount) as amount
		FROM chart_of_accounts coa
		LEFT JOIN general_ledger gl ON coa.id = gl.account_id AND gl.tenant_id = $1 AND gl.entry_date BETWEEN $2 AND $3
		WHERE coa.tenant_id = $1 AND coa.account_type = $4 AND coa.is_active = true
		GROUP BY coa.id, coa.account_code, coa.account_name
		ORDER BY coa.account_code`

	revRows, _ := s.db.Query(revQuery, tenantID, fromDate, toDate, models.AccountTypeRevenue)
	defer revRows.Close()

	var revLineItems []models.FinancialStatementLineItem
	var revTotal float64
	for revRows.Next() {
		var accountCode, accountName string
		var amount float64
		revRows.Scan(&accountCode, &accountName, &amount)
		revLineItems = append(revLineItems, models.FinancialStatementLineItem{
			AccountCode: accountCode,
			AccountName: accountName,
			Amount:      amount,
		})
		revTotal += amount
	}

	statement.Sections = append(statement.Sections, models.FinancialStatementSection{
		SectionName:  "Revenue",
		SectionType:  "Revenue",
		LineItems:    revLineItems,
		SectionTotal: revTotal,
	})

	// Get expenses
	expQuery := `SELECT coa.account_code, coa.account_name, SUM(gl.debit_amount - gl.credit_amount) as amount
		FROM chart_of_accounts coa
		LEFT JOIN general_ledger gl ON coa.id = gl.account_id AND gl.tenant_id = $1 AND gl.entry_date BETWEEN $2 AND $3
		WHERE coa.tenant_id = $1 AND coa.account_type = $4 AND coa.is_active = true
		GROUP BY coa.id, coa.account_code, coa.account_name
		ORDER BY coa.account_code`

	expRows, _ := s.db.Query(expQuery, tenantID, fromDate, toDate, models.AccountTypeExpense)
	defer expRows.Close()

	var expLineItems []models.FinancialStatementLineItem
	var expTotal float64
	for expRows.Next() {
		var accountCode, accountName string
		var amount float64
		expRows.Scan(&accountCode, &accountName, &amount)
		expLineItems = append(expLineItems, models.FinancialStatementLineItem{
			AccountCode: accountCode,
			AccountName: accountName,
			Amount:      amount,
		})
		expTotal += amount
	}

	statement.Sections = append(statement.Sections, models.FinancialStatementSection{
		SectionName:  "Expenses",
		SectionType:  "Expenses",
		LineItems:    expLineItems,
		SectionTotal: expTotal,
	})

	statement.TotalRevenue = revTotal
	statement.TotalExpenses = expTotal
	statement.NetIncome = revTotal - expTotal

	return statement, nil
}

func (s *FinancialStatementsService) GetBalanceSheet(tenantID string, asOfDate time.Time) (*models.FinancialStatement, error) {
	statement := &models.FinancialStatement{
		TenantID:      tenantID,
		StatementType: "balance_sheet",
		Period:        asOfDate.Format("2006-01-02"),
		Sections:      make([]models.FinancialStatementSection, 0),
	}

	// Get assets
	assetQuery := `SELECT coa.account_code, coa.account_name, SUM(gl.debit_amount - gl.credit_amount) as amount
		FROM chart_of_accounts coa
		LEFT JOIN general_ledger gl ON coa.id = gl.account_id AND gl.tenant_id = $1 AND gl.entry_date <= $2
		WHERE coa.tenant_id = $1 AND coa.account_type = $3 AND coa.is_active = true
		GROUP BY coa.id, coa.account_code, coa.account_name
		ORDER BY coa.account_code`

	assetRows, _ := s.db.Query(assetQuery, tenantID, asOfDate, models.AccountTypeAsset)
	defer assetRows.Close()

	var assetLineItems []models.FinancialStatementLineItem
	var assetTotal float64
	for assetRows.Next() {
		var accountCode, accountName string
		var amount float64
		assetRows.Scan(&accountCode, &accountName, &amount)
		assetLineItems = append(assetLineItems, models.FinancialStatementLineItem{
			AccountCode: accountCode,
			AccountName: accountName,
			Amount:      amount,
		})
		assetTotal += amount
	}

	statement.Sections = append(statement.Sections, models.FinancialStatementSection{
		SectionName:  "Assets",
		SectionType:  "Assets",
		LineItems:    assetLineItems,
		SectionTotal: assetTotal,
	})

	// Similar for liabilities and equity...
	statement.TotalAssets = assetTotal

	return statement, nil
}
