package models

import "time"

// ============================================================================
// ACCOUNTS (GENERAL LEDGER) MODELS
// ============================================================================

// ChartOfAccount represents an account in the chart of accounts
type ChartOfAccount struct {
	ID              string  `json:"id"`
	TenantID        string  `json:"tenant_id"`
	AccountCode     string  `json:"account_code"`
	AccountName     string  `json:"account_name"`
	AccountType     string  `json:"account_type"` // Asset, Liability, Equity, Revenue, Expense
	SubAccountType  string  `json:"sub_account_type"`
	ParentAccountID *string `json:"parent_account_id"`
	Description     string  `json:"description"`
	OpeningBalance  float64 `json:"opening_balance"`
	CurrentBalance  float64 `json:"current_balance"`
	IsActive        bool    `json:"is_active"`
	IsHeader        bool    `json:"is_header"`
	IsDefault       bool    `json:"is_default"`
	Currency        string  `json:"currency"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// JournalEntry represents a transaction entry
type JournalEntry struct {
	ID              string     `json:"id"`
	TenantID        string     `json:"tenant_id"`
	EntryDate       time.Time  `json:"entry_date"`
	ReferenceNumber *string    `json:"reference_number"`
	ReferenceType   string     `json:"reference_type"` // Manual, HR_Payroll, Sales_Invoice, etc.
	ReferenceID     *string    `json:"reference_id"`
	Description     string     `json:"description"`
	Amount          float64    `json:"amount"`
	Narration       string     `json:"narration"`
	EntryStatus     string     `json:"entry_status"` // Draft, Posted, Cancelled
	PostedBy        *string    `json:"posted_by"`
	PostedAt        *time.Time `json:"posted_at"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	// Details for UI
	Details []JournalEntryDetail `json:"details,omitempty"`
}

// JournalEntryDetail represents a debit/credit line in a journal entry
type JournalEntryDetail struct {
	ID             string  `json:"id"`
	TenantID       string  `json:"tenant_id"`
	JournalEntryID string  `json:"journal_entry_id"`
	AccountID      string  `json:"account_id"`
	AccountCode    string  `json:"account_code"`
	DebitAmount    float64 `json:"debit_amount"`
	CreditAmount   float64 `json:"credit_amount"`
	Description    string  `json:"description"`
	LineNumber     int     `json:"line_number"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GLAccountBalance represents cached balance for an account in a period
type GLAccountBalance struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	AccountID      string    `json:"account_id"`
	FiscalPeriod   time.Time `json:"fiscal_period"`
	OpeningBalance float64   `json:"opening_balance"`
	TotalDebit     float64   `json:"total_debit"`
	TotalCredit    float64   `json:"total_credit"`
	ClosingBalance float64   `json:"closing_balance"`
}

// FinancialPeriod represents an accounting period
type FinancialPeriod struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	PeriodName string     `json:"period_name"`
	PeriodType string     `json:"period_type"` // Monthly, Quarterly, Annual
	StartDate  time.Time  `json:"start_date"`
	EndDate    time.Time  `json:"end_date"`
	IsClosed   bool       `json:"is_closed"`
	ClosedBy   *string    `json:"closed_by"`
	ClosedAt   *time.Time `json:"closed_at"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// TrialBalance represents the trial balance for a period
type TrialBalance struct {
	ID            string  `json:"id"`
	TenantID      string  `json:"tenant_id"`
	PeriodID      string  `json:"period_id"`
	AccountID     string  `json:"account_id"`
	AccountCode   string  `json:"account_code"`
	AccountName   string  `json:"account_name"`
	DebitBalance  float64 `json:"debit_balance"`
	CreditBalance float64 `json:"credit_balance"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IncomeStatement represents the P&L statement
type IncomeStatement struct {
	ID                string  `json:"id"`
	TenantID          string  `json:"tenant_id"`
	PeriodID          string  `json:"period_id"`
	RevenueTotal      float64 `json:"revenue_total"`
	CostOfGoodsSold   float64 `json:"cost_of_goods_sold"`
	GrossProfit       float64 `json:"gross_profit"`
	OperatingExpenses float64 `json:"operating_expenses"`
	OperatingIncome   float64 `json:"operating_income"`
	OtherIncome       float64 `json:"other_income"`
	OtherExpenses     float64 `json:"other_expenses"`
	IncomeBeforeTax   float64 `json:"income_before_tax"`
	TaxExpense        float64 `json:"tax_expense"`
	NetIncome         float64 `json:"net_income"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BalanceSheet represents the balance sheet (Assets = Liabilities + Equity)
type BalanceSheet struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	PeriodID string `json:"period_id"`

	// Assets
	CurrentAssets float64 `json:"current_assets"`
	FixedAssets   float64 `json:"fixed_assets"`
	OtherAssets   float64 `json:"other_assets"`
	TotalAssets   float64 `json:"total_assets"`

	// Liabilities
	CurrentLiabilities  float64 `json:"current_liabilities"`
	LongTermLiabilities float64 `json:"long_term_liabilities"`
	OtherLiabilities    float64 `json:"other_liabilities"`
	TotalLiabilities    float64 `json:"total_liabilities"`

	// Equity
	PaidUpCapital    float64 `json:"paid_up_capital"`
	RetainedEarnings float64 `json:"retained_earnings"`
	TotalEquity      float64 `json:"total_equity"`

	// Validation
	IsBalanced bool `json:"is_balanced"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Reconciliation represents a bank/cash reconciliation
type Reconciliation struct {
	ID                   string     `json:"id"`
	TenantID             string     `json:"tenant_id"`
	AccountID            string     `json:"account_id"`
	ReconciliationType   string     `json:"reconciliation_type"` // Bank, Cash, Receivables, etc.
	PeriodFrom           time.Time  `json:"period_from"`
	PeriodTo             time.Time  `json:"period_to"`
	SystemBalance        *float64   `json:"system_balance"`
	ActualBalance        *float64   `json:"actual_balance"`
	Difference           *float64   `json:"difference"`
	ReconciliationStatus string     `json:"reconciliation_status"` // Pending, In Progress, Completed, Discrepancy
	ReconciledBy         *string    `json:"reconciled_by"`
	ReconciledAt         *time.Time `json:"reconciled_at"`
	Notes                string     `json:"notes"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// GLAuditLog tracks all GL changes
type GLAuditLog struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Action     string    `json:"action"`
	OldValues  string    `json:"old_values"`
	NewValues  string    `json:"new_values"`
	ChangedBy  string    `json:"changed_by"`
	ChangedAt  time.Time `json:"changed_at"`
}

// JournalEntryRequest is the request body for creating entries
type JournalEntryRequest struct {
	EntryDate       time.Time `json:"entry_date"`
	ReferenceNumber string    `json:"reference_number,omitempty"`
	ReferenceType   string    `json:"reference_type"`
	Description     string    `json:"description"`
	Narration       string    `json:"narration,omitempty"`
	Details         []struct {
		AccountID    string  `json:"account_id"`
		DebitAmount  float64 `json:"debit_amount"`
		CreditAmount float64 `json:"credit_amount"`
		Description  string  `json:"description,omitempty"`
	} `json:"details"`
}

// PostJournalEntryRequest is for posting an entry
type PostJournalEntryRequest struct {
	PostedBy string `json:"posted_by"`
}
