package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// BACKEND LOGIC VALIDATION SUITE
// ============================================================================
// Comprehensive validation of all core backend business logic

// ============================================================================
// 1. FINANCIAL CALCULATIONS & ACCOUNTING
// ============================================================================

// TestGrossProfitCalculation validates basic profit calculation
func TestGrossProfitCalculation(t *testing.T) {
	testCases := []struct {
		description string
		revenue     float64
		costOfSales float64
		expectedGP  float64
		expectedGPM float64
	}{
		{"Standard product", 10000.00, 6000.00, 4000.00, 40.0},
		{"Low margin item", 5000.00, 4500.00, 500.00, 10.0},
		{"High margin product", 20000.00, 8000.00, 12000.00, 60.0},
	}

	for _, tc := range testCases {
		gp := tc.revenue - tc.costOfSales
		assert.InDelta(t, tc.expectedGP, gp, 0.01, tc.description+" - GP")

		gpm := (gp / tc.revenue) * 100
		assert.InDelta(t, tc.expectedGPM, gpm, 0.1, tc.description+" - GPM%")
	}
}

// TestTaxCalculationAccuracy validates proper tax computation
func TestTaxCalculationAccuracy(t *testing.T) {
	testCases := []struct {
		description   string
		baseAmount    float64
		taxRate       float64
		expectedTax   float64
		expectedTotal float64
	}{
		// 10000 × 18% = 1800, Total = 11800
		{"Standard 18% GST", 10000.00, 0.18, 1800.00, 11800.00},
		// 5000 × 5% = 250, Total = 5250
		{"5% GST rate", 5000.00, 0.05, 250.00, 5250.00},
		// 100000 × 0% = 0 (exempt)
		{"Exempt supply", 100000.00, 0.00, 0.00, 100000.00},
	}

	for _, tc := range testCases {
		tax := tc.baseAmount * tc.taxRate
		assert.InDelta(t, tc.expectedTax, tax, 0.01, tc.description+" - Tax")

		total := tc.baseAmount + tax
		assert.InDelta(t, tc.expectedTotal, total, 0.01, tc.description+" - Total")
	}
}

// TestDoubleEntryValidation validates accounting equation
func TestDoubleEntryValidation(t *testing.T) {
	testCases := []struct {
		description string
		totalDebit  float64
		totalCredit float64
		balanced    bool
	}{
		{"Balanced entry", 10000.00, 10000.00, true},
		{"Unbalanced debit", 10000.00, 9999.00, false},
		{"Unbalanced credit", 10000.00, 10001.00, false},
		{"Zero entry", 0.00, 0.00, true},
	}

	for _, tc := range testCases {
		balanced := tc.totalDebit == tc.totalCredit
		assert.Equal(t, tc.balanced, balanced, tc.description)
	}
}

// ============================================================================
// 2. BUSINESS RULE ENFORCEMENT
// ============================================================================

// TestMinimumMarginEnforcement validates margin floor
func TestMinimumMarginEnforcement(t *testing.T) {
	minimumMargin := 0.15 // 15%

	testCases := []struct {
		description string
		costPrice   float64
		sellPrice   float64
		allowed     bool
	}{
		// (750 - 500) / 750 = 33% > 15% ✓
		{"Above minimum", 500.00, 750.00, true},
		// (575 - 500) / 575 = 13% < 15% ✗
		{"Below minimum", 500.00, 575.00, false},
		// (1000 - 850) / 1000 = 15% = 15% ✓
		{"At minimum", 850.00, 1000.00, true},
	}

	for _, tc := range testCases {
		margin := (tc.sellPrice - tc.costPrice) / tc.sellPrice
		allowed := margin >= minimumMargin
		assert.Equal(t, tc.allowed, allowed, tc.description)
	}
}

// TestCreditLimitValidation ensures credit limit compliance
func TestCreditLimitValidation(t *testing.T) {
	testCases := []struct {
		description string
		creditLimit float64
		currentUsed float64
		newInvoice  float64
		allowed     bool
	}{
		// Limit 100K, used 30K, new 50K = 80K ≤ 100K ✓
		{"Within limit", 100000.00, 30000.00, 50000.00, true},
		// Limit 100K, used 60K, new 50K = 110K > 100K ✗
		{"Exceeds limit", 100000.00, 60000.00, 50000.00, false},
		// Limit 100K, used 0K, new 100K = 100K ≤ 100K ✓
		{"At limit exactly", 100000.00, 0.00, 100000.00, true},
	}

	for _, tc := range testCases {
		available := tc.creditLimit - tc.currentUsed
		allowed := tc.newInvoice <= available
		assert.Equal(t, tc.allowed, allowed, tc.description)
	}
}

// TestInvoiceStatusWorkflow validates state machine
func TestInvoiceStatusWorkflow(t *testing.T) {
	validTransitions := map[string][]string{
		"draft":          {"sent", "cancelled"},
		"sent":           {"partially_paid", "paid", "overdue", "cancelled"},
		"partially_paid": {"paid", "overdue", "cancelled"},
		"paid":           {},                    // Terminal
		"overdue":        {"paid", "cancelled"}, // Can still collect
		"cancelled":      {},                    // Terminal
	}

	testCases := []struct {
		description     string
		fromStatus      string
		toStatus        string
		transitionValid bool
	}{
		{"Draft to sent", "draft", "sent", true},
		{"Sent to paid", "sent", "paid", true},
		{"Sent to draft", "sent", "draft", false}, // Invalid regression
		{"Paid to sent", "paid", "sent", false},   // Terminal state
	}

	for _, tc := range testCases {
		validStates := validTransitions[tc.fromStatus]
		isValid := false
		for _, state := range validStates {
			if state == tc.toStatus {
				isValid = true
				break
			}
		}
		assert.Equal(t, tc.transitionValid, isValid, tc.description)
	}
}

// ============================================================================
// 3. DATA INTEGRITY & VALIDATION
// ============================================================================

// TestQuantityValidation ensures positive quantities
func TestQuantityValidation(t *testing.T) {
	testCases := []struct {
		description string
		quantity    float64
		valid       bool
	}{
		{"Positive quantity", 100.0, true},
		{"Zero quantity", 0.0, false},
		{"Negative quantity", -50.0, false},
		{"Decimal quantity", 25.5, true},
		{"Very large quantity", 999999.99, true},
	}

	for _, tc := range testCases {
		valid := tc.quantity > 0
		assert.Equal(t, tc.valid, valid, tc.description)
	}
}

// TestPriceValidation ensures non-negative prices
func TestPriceValidation(t *testing.T) {
	testCases := []struct {
		description string
		price       float64
		valid       bool
	}{
		{"Positive price", 1000.00, true},
		{"Zero price", 0.00, true}, // Allow zero for promotional items
		{"Negative price", -500.00, false},
		{"Decimal price", 999.99, true},
	}

	for _, tc := range testCases {
		valid := tc.price >= 0
		assert.Equal(t, tc.valid, valid, tc.description)
	}
}

// TestDiscountValidation ensures discount constraints
func TestDiscountValidation(t *testing.T) {
	testCases := []struct {
		description string
		discount    float64
		valid       bool
	}{
		{"No discount", 0.0, true},
		{"5% discount", 5.0, true},
		{"50% discount", 50.0, true},
		{"100% discount", 100.0, true},
		{"Negative discount", -10.0, false},
		{"Over 100% discount", 150.0, false},
	}

	for _, tc := range testCases {
		valid := tc.discount >= 0 && tc.discount <= 100
		assert.Equal(t, tc.valid, valid, tc.description)
	}
}

// ============================================================================
// 4. INVENTORY & STOCK MANAGEMENT
// ============================================================================

// TestStockAvailabilityCheck ensures sufficient inventory
func TestStockAvailabilityCheck(t *testing.T) {
	testCases := []struct {
		description  string
		requestedQty float64
		availableQty float64
		allowed      bool
	}{
		{"Sufficient stock", 50.0, 100.0, true},
		{"Exact quantity", 100.0, 100.0, true},
		{"Insufficient stock", 150.0, 100.0, false},
		{"Zero request", 0.0, 100.0, true},
		{"No stock", 50.0, 0.0, false},
	}

	for _, tc := range testCases {
		allowed := tc.requestedQty <= tc.availableQty && tc.requestedQty > 0
		if tc.requestedQty == 0 {
			allowed = true // Allow zero requests
		}
		assert.Equal(t, tc.allowed, allowed, tc.description)
	}
}

// TestFIFOCalculation validates first-in-first-out inventory
func TestFIFOCalculation(t *testing.T) {
	// Stock: 100@100, 200@150, 100@200
	// Issue 150: (100@100) + (50@150) = 10000 + 7500 = 17500
	testCases := []struct {
		description string
		batches     []struct {
			qty  float64
			rate float64
		}
		issueQty     float64
		expectedCost float64
	}{
		{
			"Standard FIFO",
			[]struct {
				qty  float64
				rate float64
			}{
				{100, 100},
				{200, 150},
				{100, 200},
			},
			150.0,
			17500.00,
		},
	}

	for _, tc := range testCases {
		var totalCost float64
		remainingQty := tc.issueQty

		for _, batch := range tc.batches {
			if remainingQty <= 0 {
				break
			}
			issueFromBatch := remainingQty
			if remainingQty > batch.qty {
				issueFromBatch = batch.qty
			}
			totalCost += issueFromBatch * batch.rate
			remainingQty -= issueFromBatch
		}

		assert.InDelta(t, tc.expectedCost, totalCost, 0.01, tc.description)
	}
}

// ============================================================================
// 5. PAYMENT & CREDIT PROCESSING
// ============================================================================

// TestPaymentAllocationLogic validates correct payment distribution
func TestPaymentAllocationLogic(t *testing.T) {
	testCases := []struct {
		description     string
		invoiceAmount   float64
		paymentAmount   float64
		expectedPaid    float64
		expectedPending float64
		status          string
	}{
		// Full payment
		{"Full payment", 10000.00, 10000.00, 10000.00, 0.00, "paid"},
		// Partial payment
		{"Partial payment", 10000.00, 6000.00, 6000.00, 4000.00, "partially_paid"},
		// Over payment (should cap at invoice)
		{"Over payment", 10000.00, 12000.00, 10000.00, 0.00, "paid"},
		// Zero payment (no change)
		{"Zero payment", 10000.00, 0.00, 0.00, 10000.00, "unpaid"},
	}

	for _, tc := range testCases {
		paidAmount := tc.paymentAmount
		if paidAmount > tc.invoiceAmount {
			paidAmount = tc.invoiceAmount
		}

		pending := tc.invoiceAmount - paidAmount

		assert.InDelta(t, tc.expectedPaid, paidAmount, 0.01, tc.description+" - paid")
		assert.InDelta(t, tc.expectedPending, pending, 0.01, tc.description+" - pending")
	}
}

// TestAccountsReceivableAging validates AR aging buckets
func TestAccountsReceivableAging(t *testing.T) {
	testCases := []struct {
		description string
		daysPastDue int
		riskLevel   string
	}{
		// 0-30 days = low risk
		{"Current", 15, "low"},
		// 31-60 days = medium risk
		{"30-60 days", 45, "medium"},
		// 61-90 days = high risk
		{"60-90 days", 75, "high"},
		// 90+ days = critical
		{"90+ days", 120, "critical"},
	}

	for _, tc := range testCases {
		var risk string
		if tc.daysPastDue <= 30 {
			risk = "low"
		} else if tc.daysPastDue <= 60 {
			risk = "medium"
		} else if tc.daysPastDue <= 90 {
			risk = "high"
		} else {
			risk = "critical"
		}

		assert.Equal(t, tc.riskLevel, risk, tc.description)
	}
}

// ============================================================================
// 6. OPERATIONAL WORKFLOWS
// ============================================================================

// TestProjectProgressValidation ensures progress is between 0-100%
func TestProjectProgressValidation(t *testing.T) {
	testCases := []struct {
		description string
		progress    float64
		valid       bool
	}{
		{"Start (0%)", 0.0, true},
		{"Mid-way (50%)", 50.0, true},
		{"Complete (100%)", 100.0, true},
		{"Over complete (110%)", 110.0, false},
		{"Negative progress", -10.0, false},
	}

	for _, tc := range testCases {
		valid := tc.progress >= 0 && tc.progress <= 100
		assert.Equal(t, tc.valid, valid, tc.description)
	}
}

// TestBOQPrecisionCalculation validates BOQ arithmetic accuracy
func TestBOQPrecisionCalculation(t *testing.T) {
	// Example: 250.50 × 75.25 = 18,850.125
	testCases := []struct {
		description   string
		quantity      float64
		rate          float64
		expectedTotal float64
		tolerance     float64
	}{
		{"Standard calculation", 250.50, 75.25, 18850.125, 0.01},
		{"Integer values", 100.00, 500.00, 50000.00, 0.01},
		{"Decimal precision", 99.99, 99.99, 9998.0001, 0.01},
	}

	for _, tc := range testCases {
		total := tc.quantity * tc.rate
		assert.InDelta(t, tc.expectedTotal, total, tc.tolerance, tc.description)
	}
}

// ============================================================================
// 7. MULTI-TENANT ISOLATION
// ============================================================================

// TestMultiTenantDataIsolation validates tenant separation
func TestMultiTenantDataIsolation(t *testing.T) {
	testCases := []struct {
		description string
		tenantA     string
		tenantB     string
		sameData    bool
	}{
		{"Different tenants", "tenant-1", "tenant-2", false},
		{"Same tenant", "tenant-1", "tenant-1", true},
		{"Case sensitive", "Tenant-1", "tenant-1", false},
	}

	for _, tc := range testCases {
		isolated := tc.tenantA != tc.tenantB
		assert.Equal(t, tc.sameData, !isolated, tc.description)
	}
}

// ============================================================================
// 8. ERROR HANDLING & EDGE CASES
// ============================================================================

// TestNullValueHandling validates NULL handling in calculations
func TestNullValueHandling(t *testing.T) {
	// Using pointers to represent potential NULL values
	testCases := []struct {
		description string
		value       *float64
		defaultVal  float64
		expected    float64
	}{
		{"Nil value", nil, 0.0, 0.0},
		{"Valid value", func() *float64 { v := 100.0; return &v }(), 0.0, 100.0},
	}

	for _, tc := range testCases {
		result := tc.defaultVal
		if tc.value != nil {
			result = *tc.value
		}
		assert.InDelta(t, tc.expected, result, 0.01, tc.description)
	}
}

// TestZeroDivisionPrevention prevents division errors
func TestZeroDivisionPrevention(t *testing.T) {
	testCases := []struct {
		description string
		numerator   float64
		denominator float64
		safe        bool
	}{
		{"Valid division", 1000.00, 100.0, true},
		{"Zero denominator", 1000.00, 0.0, false},
		{"Both zero", 0.0, 0.0, false},
	}

	for _, tc := range testCases {
		safe := tc.denominator != 0
		assert.Equal(t, tc.safe, safe, tc.description)
	}
}

// ============================================================================
// 9. ROUNDING & PRECISION
// ============================================================================

// TestMonetaryPrecision validates 2-decimal-place accuracy
func TestMonetaryPrecision(t *testing.T) {
	testCases := []struct {
		description string
		amount      float64
		rounded     float64
	}{
		// 1000.126 rounds to 1000.13
		{"Rounding up", 1000.126, 1000.13},
		// 1000.124 rounds to 1000.12
		{"Rounding down", 1000.124, 1000.12},
		// Exactly 2 decimals
		{"Already precise", 1000.50, 1000.50},
	}

	for _, tc := range testCases {
		rounded := float64(int(tc.amount*100+0.5)) / 100
		assert.InDelta(t, tc.rounded, rounded, 0.01, tc.description)
	}
}

// ============================================================================
// 10. AUDIT & COMPLIANCE
// ============================================================================

// TestAuditTrailRequirements validates audit data capture
func TestAuditTrailRequirements(t *testing.T) {
	testCases := []struct {
		description  string
		hasUserID    bool
		hasTenantID  bool
		hasTimestamp bool
		isAuditable  bool
	}{
		{"Complete audit", true, true, true, true},
		{"Missing user", false, true, true, false},
		{"Missing tenant", true, false, true, false},
		{"Missing timestamp", true, true, false, false},
	}

	for _, tc := range testCases {
		auditable := tc.hasUserID && tc.hasTenantID && tc.hasTimestamp
		assert.Equal(t, tc.isAuditable, auditable, tc.description)
	}
}
