package models

import (
	"time"
)

// ===== FINANCE/ACCOUNTING MODELS =====

// ChartOfAccountsType represents account types
type ChartOfAccountsType string

const (
	AccountTypeAsset     ChartOfAccountsType = "asset"
	AccountTypeLiability ChartOfAccountsType = "liability"
	AccountTypeEquity    ChartOfAccountsType = "equity"
	AccountTypeRevenue   ChartOfAccountsType = "revenue"
	AccountTypeExpense   ChartOfAccountsType = "expense"
)

// JournalEntryStatus represents status of journal entry
type JournalEntryStatus string

const (
	JournalEntryStatusDraft    JournalEntryStatus = "draft"
	JournalEntryStatusPosted   JournalEntryStatus = "posted"
	JournalEntryStatusCanceled JournalEntryStatus = "canceled"
)

// ChartOfAccounts represents the Chart of Accounts
type ChartOfAccounts struct {
	ID             string              `db:"id" json:"id"`
	TenantID       string              `db:"tenant_id" json:"tenant_id"`
	AccountCode    string              `db:"account_code" json:"account_code"`
	AccountName    string              `db:"account_name" json:"account_name"`
	AccountType    ChartOfAccountsType `db:"account_type" json:"account_type"`
	SubType        string              `db:"sub_type" json:"sub_type"` // asset, current_asset, fixed_asset, etc
	Description    string              `db:"description" json:"description"`
	IsActive       bool                `db:"is_active" json:"is_active"`
	OpeningBalance float64             `db:"opening_balance" json:"opening_balance"`
	CurrentBalance float64             `db:"current_balance" json:"current_balance"`
	Currency       string              `db:"currency" json:"currency"`
	CreatedAt      time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time           `db:"updated_at" json:"updated_at"`
}

// JournalEntry represents an accounting journal entry
type JournalEntry struct {
	ID           string             `db:"id" json:"id"`
	TenantID     string             `db:"tenant_id" json:"tenant_id"`
	EntryNumber  string             `db:"entry_number" json:"entry_number"`
	EntryDate    time.Time          `db:"entry_date" json:"entry_date"`
	Reference    string             `db:"reference" json:"reference"`
	Description  string             `db:"description" json:"description"`
	Status       JournalEntryStatus `db:"status" json:"status"`
	PreparedByID string             `db:"prepared_by_id" json:"prepared_by_id"`
	ApprovedByID *string            `db:"approved_by_id" json:"approved_by_id"`
	TotalDebit   float64            `db:"total_debit" json:"total_debit"`
	TotalCredit  float64            `db:"total_credit" json:"total_credit"`
	CreatedAt    time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `db:"updated_at" json:"updated_at"`
	PostedAt     *time.Time         `db:"posted_at" json:"posted_at"`
}

// JournalEntryLine represents a line item in a journal entry
type JournalEntryLine struct {
	ID             string    `db:"id" json:"id"`
	TenantID       string    `db:"tenant_id" json:"tenant_id"`
	JournalEntryID string    `db:"journal_entry_id" json:"journal_entry_id"`
	AccountID      string    `db:"account_id" json:"account_id"`
	AccountCode    string    `db:"account_code" json:"account_code"`
	AccountName    string    `db:"account_name" json:"account_name"`
	DebitAmount    float64   `db:"debit_amount" json:"debit_amount"`
	CreditAmount   float64   `db:"credit_amount" json:"credit_amount"`
	Description    string    `db:"description" json:"description"`
	Reference      string    `db:"reference" json:"reference"` // Invoice, bill, etc
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

// GeneralLedger represents account ledger entries
type GeneralLedger struct {
	ID             string    `db:"id" json:"id"`
	TenantID       string    `db:"tenant_id" json:"tenant_id"`
	AccountID      string    `db:"account_id" json:"account_id"`
	AccountCode    string    `db:"account_code" json:"account_code"`
	JournalEntryID string    `db:"journal_entry_id" json:"journal_entry_id"`
	EntryDate      time.Time `db:"entry_date" json:"entry_date"`
	Reference      string    `db:"reference" json:"reference"`
	DebitAmount    float64   `db:"debit_amount" json:"debit_amount"`
	CreditAmount   float64   `db:"credit_amount" json:"credit_amount"`
	RunningBalance float64   `db:"running_balance" json:"running_balance"`
	Description    string    `db:"description" json:"description"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

// TrialBalance represents trial balance summary
type TrialBalance struct {
	AccountCode  string  `db:"account_code" json:"account_code"`
	AccountName  string  `db:"account_name" json:"account_name"`
	AccountType  string  `db:"account_type" json:"account_type"`
	DebitAmount  float64 `db:"debit_amount" json:"debit_amount"`
	CreditAmount float64 `db:"credit_amount" json:"credit_amount"`
}

// AccountsReceivable represents customer invoices (receivable)
type AccountsReceivable struct {
	ID                string    `db:"id" json:"id"`
	TenantID          string    `db:"tenant_id" json:"tenant_id"`
	CustomerID        string    `db:"customer_id" json:"customer_id"`
	CustomerName      string    `db:"customer_name" json:"customer_name"`
	InvoiceID         string    `db:"invoice_id" json:"invoice_id"`
	InvoiceNumber     string    `db:"invoice_number" json:"invoice_number"`
	InvoiceAmount     float64   `db:"invoice_amount" json:"invoice_amount"`
	AmountPaid        float64   `db:"amount_paid" json:"amount_paid"`
	OutstandingAmount float64   `db:"outstanding_amount" json:"outstanding_amount"`
	InvoiceDate       time.Time `db:"invoice_date" json:"invoice_date"`
	DueDate           time.Time `db:"due_date" json:"due_date"`
	Status            string    `db:"status" json:"status"` // invoiced, partial_paid, paid, overdue
	Currency          string    `db:"currency" json:"currency"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// AccountsPayable represents vendor bills (payable)
type AccountsPayable struct {
	ID                string    `db:"id" json:"id"`
	TenantID          string    `db:"tenant_id" json:"tenant_id"`
	VendorID          string    `db:"vendor_id" json:"vendor_id"`
	VendorName        string    `db:"vendor_name" json:"vendor_name"`
	BillID            string    `db:"bill_id" json:"bill_id"`
	BillNumber        string    `db:"bill_number" json:"bill_number"`
	BillAmount        float64   `db:"bill_amount" json:"bill_amount"`
	AmountPaid        float64   `db:"amount_paid" json:"amount_paid"`
	OutstandingAmount float64   `db:"outstanding_amount" json:"outstanding_amount"`
	BillDate          time.Time `db:"bill_date" json:"bill_date"`
	DueDate           time.Time `db:"due_date" json:"due_date"`
	Status            string    `db:"status" json:"status"` // received, partial_paid, paid, overdue
	Currency          string    `db:"currency" json:"currency"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// FinancialStatement represents P&L or Balance Sheet
type FinancialStatement struct {
	TenantID         string                      `json:"tenant_id"`
	StatementType    string                      `json:"statement_type"` // income_statement, balance_sheet
	Period           string                      `json:"period"`         // Month/Year
	Sections         []FinancialStatementSection `json:"sections"`
	TotalAssets      float64                     `json:"total_assets,omitempty"`
	TotalLiabilities float64                     `json:"total_liabilities,omitempty"`
	TotalEquity      float64                     `json:"total_equity,omitempty"`
	TotalRevenue     float64                     `json:"total_revenue,omitempty"`
	TotalExpenses    float64                     `json:"total_expenses,omitempty"`
	NetIncome        float64                     `json:"net_income,omitempty"`
}

// FinancialStatementSection represents sections in financial statements
type FinancialStatementSection struct {
	SectionName  string                       `json:"section_name"`
	SectionType  string                       `json:"section_type"` // Assets, Liabilities, etc
	LineItems    []FinancialStatementLineItem `json:"line_items"`
	SectionTotal float64                      `json:"section_total"`
}

// FinancialStatementLineItem represents line items in financial statements
type FinancialStatementLineItem struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Amount      float64 `json:"amount"`
}

// REQUEST/RESPONSE MODELS

type CreateChartOfAccountsRequest struct {
	AccountCode    string              `json:"account_code" binding:"required"`
	AccountName    string              `json:"account_name" binding:"required"`
	AccountType    ChartOfAccountsType `json:"account_type" binding:"required"`
	SubType        string              `json:"sub_type"`
	Description    string              `json:"description"`
	OpeningBalance float64             `json:"opening_balance"`
}

type CreateJournalEntryRequest struct {
	EntryDate   time.Time                       `json:"entry_date" binding:"required"`
	Reference   string                          `json:"reference" binding:"required"`
	Description string                          `json:"description"`
	LineItems   []CreateJournalEntryLineRequest `json:"line_items" binding:"required,min=2"`
}

type CreateJournalEntryLineRequest struct {
	AccountID    string  `json:"account_id" binding:"required"`
	DebitAmount  float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
	Description  string  `json:"description"`
	Reference    string  `json:"reference"`
}

type PostJournalEntryRequest struct {
	ApprovedByID string `json:"approved_by_id" binding:"required"`
	Comment      string `json:"comment"`
}

// DashboardMetrics represents financial dashboard metrics
type FinancialDashboard struct {
	TotalAssets          float64 `json:"total_assets"`
	TotalLiabilities     float64 `json:"total_liabilities"`
	TotalEquity          float64 `json:"total_equity"`
	CurrentMonthRevenue  float64 `json:"current_month_revenue"`
	CurrentMonthExpenses float64 `json:"current_month_expenses"`
	NetProfit            float64 `json:"net_profit"`
	ProfitMargin         float64 `json:"profit_margin"`
	AccountsReceivable   float64 `json:"accounts_receivable"`
	AccountsPayable      float64 `json:"accounts_payable"`
	CashBalance          float64 `json:"cash_balance"`
	OutstandingInvoices  int64   `json:"outstanding_invoices"`
	OutstandingBills     int64   `json:"outstanding_bills"`
	QuickRatio           float64 `json:"quick_ratio"`
	CurrentRatio         float64 `json:"current_ratio"`
	DebtToEquity         float64 `json:"debt_to_equity"`
}
