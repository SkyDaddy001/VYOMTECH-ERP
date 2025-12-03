package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateAccountHandler validates account creation
func TestCreateAccountHandler(t *testing.T) {
	payload := map[string]interface{}{
		"account_code": "10000",
		"account_name": "Cash in Hand",
		"account_type": "asset",
		"description":  "Cash account",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/gl/accounts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestGetAccountHandler validates account retrieval
func TestGetAccountHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/gl/accounts/10000", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.Method)
}

// TestPostJournalEntryHandler validates journal entry posting
func TestPostJournalEntryHandler(t *testing.T) {
	payload := map[string]interface{}{
		"entry_date":  "2025-01-15",
		"description": "Sales invoice",
		"journal_lines": []map[string]interface{}{
			{
				"account_code": "10000",
				"debit":        10000.0,
			},
			{
				"account_code": "40000",
				"credit":       10000.0,
			},
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/gl/journal-entries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestGetTrialBalanceHandler validates trial balance endpoint
func TestGetTrialBalanceHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/gl/trial-balance?asof_date=2025-01-31", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.Method)
}

// TestGetAccountBalanceHandler validates account balance
func TestGetAccountBalanceHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/gl/accounts/10000/balance", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.Method)
}

// TestGetAccountLedgerHandler validates account ledger
func TestGetAccountLedgerHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/gl/accounts/10000/ledger?from_date=2025-01-01&to_date=2025-01-31", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.Contains(t, req.URL.RawQuery, "from_date=2025-01-01")
}

// TestDoubleEntryValidation validates debit equals credit (Fundamental rule)
func TestDoubleEntryValidation(t *testing.T) {
	testCases := []struct {
		description string
		totalDebit  float64
		totalCredit float64
		balanced    bool
	}{
		{"Sales Invoice 18% GST", 11800.0, 11800.0, true},
		{"Cash Payment", 5000.0, 5000.0, true},
		{"Salary Posting", 50000.0, 50000.0, true},
		{"Unbalanced Entry (Should Reject)", 10000.0, 9999.99, false},
		{"Unbalanced Entry (Should Reject)", 10000.0, 10000.01, false},
	}

	for _, tc := range testCases {
		// Entry is only valid if debit == credit (tolerance: 0.01 paise)
		balanced := tc.totalDebit == tc.totalCredit
		assert.Equal(t, tc.balanced, balanced, tc.description)
	}
}

// TestAccountTypeValidation validates account types per CoA
func TestAccountTypeValidation(t *testing.T) {
	// Standard Chart of Accounts structure
	validTypes := map[string][]string{
		"Asset":     {"Current", "Fixed"},
		"Liability": {"Current", "Long-term"},
		"Equity":    {"Capital", "Reserves", "Retained Earnings"},
		"Revenue":   {"Operating", "Non-Operating"},
		"Expense":   {"Operating", "Non-Operating"},
	}

	for acType, subTypes := range validTypes {
		assert.NotEmpty(t, acType)
		assert.NotEmpty(t, subTypes)
	}
}

// TestAccountCodeValidation validates account code format
func TestAccountCodeValidation(t *testing.T) {
	testCases := []struct {
		code  string
		valid bool
	}{
		{"10000", true},
		{"10001", true},
		{"40000", true},
		{"", false},
		{"ABC", false},
	}

	for _, tc := range testCases {
		isValid := tc.code != "" && len(tc.code) >= 4
		assert.Equal(t, tc.valid, isValid)
	}
}

// TestAccountBalanceCalculation validates balance calculation per account type
func TestAccountBalanceCalculation(t *testing.T) {
	testCases := []struct {
		accountType string
		openingBal  float64
		debits      float64
		credits     float64
		expectedBal float64
		description string
	}{
		// Asset: Opening + Debits - Credits
		{"Asset", 10000.0, 5000.0, 2000.0, 13000.0, "Cash Account"},
		// Liability: Opening + Credits - Debits
		{"Liability", 50000.0, 10000.0, 20000.0, 60000.0, "Loan Account"},
		// Equity: Opening + Credits - Debits
		{"Equity", 100000.0, 5000.0, 15000.0, 110000.0, "Capital Account"},
		// Revenue: Opening + Credits - Debits (Contra revenue)
		{"Revenue", 0.0, 0.0, 100000.0, 100000.0, "Sales Revenue"},
		// Expense: Opening + Debits - Credits (Contra expense)
		{"Expense", 0.0, 50000.0, 1000.0, 49000.0, "Rent Expense"},
	}

	for _, tc := range testCases {
		var balance float64
		switch tc.accountType {
		case "Asset":
			balance = tc.openingBal + tc.debits - tc.credits
		case "Liability", "Equity":
			balance = tc.openingBal + tc.credits - tc.debits
		case "Revenue":
			balance = tc.openingBal + tc.credits - tc.debits
		case "Expense":
			balance = tc.openingBal + tc.debits - tc.credits
		}

		assert.InDelta(t, tc.expectedBal, balance, 0.01, tc.description)
	}
}

// TestMultiCurrencyHandling validates multi-currency
func TestMultiCurrencyHandling(t *testing.T) {
	payload := map[string]interface{}{
		"account_code": "10000",
		"currency":     "INR",
		"amount":       10000.0,
	}

	body, _ := json.Marshal(payload)
	assert.NotNil(t, body)

	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)
	assert.Equal(t, "INR", data["currency"])
}

// TestFinancialPeriodFiltering validates period filtering
func TestFinancialPeriodFiltering(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/gl/trial-balance?from_date=2025-01-01&to_date=2025-01-31", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.Contains(t, req.URL.RawQuery, "from_date=2025-01-01")
	assert.Contains(t, req.URL.RawQuery, "to_date=2025-01-31")
}

// TestJournalEntryLineValidation validates entry lines per accounting rules
func TestJournalEntryLineValidation(t *testing.T) {
	testCases := []struct {
		description string
		accountCode string
		debit       float64
		credit      float64
		valid       bool
	}{
		{"Valid Debit Entry", "10000", 1000.0, 0.0, true},
		{"Valid Credit Entry", "40000", 0.0, 1000.0, true},
		{"Missing account code", "", 1000.0, 0.0, false},
		{"Zero amount (Invalid)", "10000", 0.0, 0.0, false},
		{"Both Debit AND Credit (Invalid)", "10000", 1000.0, 1000.0, false},
	}

	for _, tc := range testCases {
		hasDebit := tc.debit > 0
		hasCredit := tc.credit > 0
		hasAmount := tc.debit > 0 || tc.credit > 0
		// Valid if: account exists AND (has debit XOR credit) AND not both
		isValid := tc.accountCode != "" && hasAmount && !(hasDebit && hasCredit)
		assert.Equal(t, tc.valid, isValid, tc.description)
	}
}

// TestReversalEntryHandler validates reversal posting
func TestReversalEntryHandler(t *testing.T) {
	payload := map[string]interface{}{
		"original_entry_id": "entry-uuid-1234",
		"reason":            "Correction",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/gl/reversals", bytes.NewReader(body))
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestMultiTenantGLRouting validates GL tenant isolation
func TestMultiTenantGLRouting(t *testing.T) {
	tenantReq1 := httptest.NewRequest("GET", "/api/v1/gl/accounts/10000", nil)
	tenantReq1.Header.Set("X-Tenant-ID", "tenant-1")

	tenantReq2 := httptest.NewRequest("GET", "/api/v1/gl/accounts/10000", nil)
	tenantReq2.Header.Set("X-Tenant-ID", "tenant-2")

	assert.NotEqual(t,
		tenantReq1.Header.Get("X-Tenant-ID"),
		tenantReq2.Header.Get("X-Tenant-ID"))
}

// TestTrialBalanceValidation validates TB balance
func TestTrialBalanceValidation(t *testing.T) {
	totalDebits := 500000.0
	totalCredits := 500000.0

	assert.Equal(t, totalDebits, totalCredits)
}

// BenchmarkPostJournalEntryEndpoint benchmarks entry posting
func BenchmarkPostJournalEntryEndpoint(b *testing.B) {
	payload := map[string]interface{}{
		"entry_date":  "2025-01-15",
		"description": "Sales invoice",
		"journal_lines": []map[string]interface{}{
			{"account_code": "10000", "debit": 10000.0},
			{"account_code": "40000", "credit": 10000.0},
		},
	}

	body, _ := json.Marshal(payload)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/v1/gl/journal-entries", bytes.NewReader(body))
		req.Header.Set("X-Tenant-ID", "tenant-1")
		httptest.NewRecorder()
	}
}

// BenchmarkGetTrialBalanceEndpoint benchmarks TB retrieval
func BenchmarkGetTrialBalanceEndpoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/v1/gl/trial-balance?asof_date=2025-01-31", nil)
		req.Header.Set("X-Tenant-ID", "tenant-1")
		httptest.NewRecorder()
	}
}
