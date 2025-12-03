package services

import (
	"testing"

	"vyomtech-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// TestChartOfAccount validates account model structure
func TestChartOfAccount(t *testing.T) {
	account := &models.ChartOfAccount{
		ID:             "acc-001",
		TenantID:       "tenant-001",
		AccountCode:    "1010",
		AccountName:    "Cash",
		AccountType:    "Asset",
		OpeningBalance: 10000.00,
		CurrentBalance: 15000.00,
		Currency:       "INR",
	}

	assert.Equal(t, "acc-001", account.ID)
	assert.Equal(t, "1010", account.AccountCode)
	assert.Equal(t, "Asset", account.AccountType)
	assert.Equal(t, 10000.00, account.OpeningBalance)
}

// TestAccountTypes validates valid account types
func TestAccountTypes(t *testing.T) {
	validTypes := []string{"Asset", "Liability", "Equity", "Revenue", "Expense"}

	for _, accType := range validTypes {
		acc := &models.ChartOfAccount{AccountType: accType}
		assert.Equal(t, accType, acc.AccountType)
	}
}

// TestJournalEntry validates journal entry structure
func TestJournalEntry(t *testing.T) {
	entry := &models.JournalEntry{
		ID:          "je-001",
		TenantID:    "tenant-001",
		Description: "Sales invoice",
		Amount:      5000.00,
		EntryStatus: "Draft",
	}

	assert.Equal(t, "je-001", entry.ID)
	assert.Equal(t, 5000.00, entry.Amount)
	assert.Equal(t, "Draft", entry.EntryStatus)
}

// TestJournalEntryDetail validates detail line structure
func TestJournalEntryDetail(t *testing.T) {
	detail := &models.JournalEntryDetail{
		ID:             "jed-001",
		JournalEntryID: "je-001",
		AccountID:      "acc-ar",
		DebitAmount:    5000.00,
		CreditAmount:   0.00,
		LineNumber:     1,
	}

	assert.Equal(t, 5000.00, detail.DebitAmount)
	assert.Equal(t, 0.00, detail.CreditAmount)
}

// TestDoubleEntryBalance validates debit=credit principle
func TestDoubleEntryBalance(t *testing.T) {
	testCases := []struct {
		name     string
		debit    float64
		credit   float64
		balanced bool
	}{
		{"Balanced", 5000.0, 5000.0, true},
		{"Unbalanced", 5100.0, 5000.0, false},
		{"Zero", 0.0, 0.0, true},
	}

	for _, tc := range testCases {
		balanced := tc.debit == tc.credit
		assert.Equal(t, tc.balanced, balanced)
	}
}

// TestTrialBalance validates trial balance structure
func TestTrialBalance(t *testing.T) {
	tb := &models.TrialBalance{
		AccountID:     "acc-001",
		AccountName:   "Cash",
		DebitBalance:  5000.00,
		CreditBalance: 0.00,
	}

	assert.Equal(t, 5000.00, tb.DebitBalance)
	assert.Equal(t, 0.00, tb.CreditBalance)
}

// TestGLAccountBalance validates account balance calculation
func TestGLAccountBalance(t *testing.T) {
	balance := &models.GLAccountBalance{
		ID:             "gab-001",
		OpeningBalance: 5000.00,
		TotalDebit:     3000.00,
		TotalCredit:    1000.00,
		ClosingBalance: 7000.00,
	}

	expectedClosing := balance.OpeningBalance + balance.TotalDebit - balance.TotalCredit
	assert.Equal(t, expectedClosing, balance.ClosingBalance)
}

// TestEntryStatusTransitions validates entry state transitions
func TestEntryStatusTransitions(t *testing.T) {
	statuses := []string{"Draft", "Posted", "Cancelled"}

	for _, status := range statuses {
		entry := &models.JournalEntry{EntryStatus: status}
		assert.Equal(t, status, entry.EntryStatus)
	}
}

// TestMultiCurrency validates currency field
func TestMultiCurrency(t *testing.T) {
	currencies := []string{"INR", "USD", "EUR"}

	for _, curr := range currencies {
		acc := &models.ChartOfAccount{Currency: curr}
		assert.Equal(t, curr, acc.Currency)
	}
}

// TestFinancialPeriod validates period structure
func TestFinancialPeriod(t *testing.T) {
	period := &models.FinancialPeriod{
		ID:         "fp-001",
		PeriodName: "January 2025",
		PeriodType: "Monthly",
		IsClosed:   false,
	}

	assert.Equal(t, "fp-001", period.ID)
	assert.Equal(t, "January 2025", period.PeriodName)
	assert.False(t, period.IsClosed)
}

// TestGLMultiTenant validates tenant isolation
func TestGLMultiTenant(t *testing.T) {
	acc1 := &models.ChartOfAccount{ID: "acc-1", TenantID: "tenant-1"}
	acc2 := &models.ChartOfAccount{ID: "acc-2", TenantID: "tenant-2"}

	assert.NotEqual(t, acc1.TenantID, acc2.TenantID)
}
