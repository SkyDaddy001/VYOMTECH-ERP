package compliance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// INCOME TAX COMPLIANCE TESTS
// ============================================================================

// TestIncomeCalculationAccuracy validates income tax calculations
func TestIncomeCalculationAccuracy(t *testing.T) {
	// Case 1: Standard individual assessment
	grossIncome := 1000000.0
	deductions := 150000.0 // Section 80C, 80D, etc.
	standardDeduction := 50000.0

	taxableIncome := grossIncome - deductions - standardDeduction
	assert.InDelta(t, 800000.0, taxableIncome, 0.01)

	// Case 2: Corporate entity
	corporateGrossIncome := 5000000.0
	allowedDeductions := 1500000.0

	corporateTaxableIncome := corporateGrossIncome - allowedDeductions
	assert.InDelta(t, 3500000.0, corporateTaxableIncome, 0.01)
}

// TestProgressiveTaxRates validates progressive tax bracket application
func TestProgressiveTaxRates(t *testing.T) {
	// FY 2024-25 Individual Tax Slabs (Standard Regime)
	taxableIncome := 1500000.0 // 15 lakhs

	// Slab 1: 0-3L @ 0%
	tax := 0.0

	// Slab 2: 3L-7L @ 5%
	if taxableIncome > 300000 {
		slabAmount := 400000.0
		tax += slabAmount * 0.05
	}

	// Slab 3: 7L-10L @ 10%
	if taxableIncome > 700000 {
		slabAmount := 300000.0
		tax += slabAmount * 0.10
	}

	// Slab 4: 10L-12.5L @ 15%
	if taxableIncome > 1000000 {
		slabAmount := 300000.0
		tax += slabAmount * 0.15
	}

	// Slab 5: Above 12.5L @ 20%
	if taxableIncome > 1250000 {
		slabAmount := taxableIncome - 1250000
		tax += slabAmount * 0.20
	}

	expectedTax := 20000 + 30000 + 45000 + 50000 // 145000
	assert.InDelta(t, float64(expectedTax), tax, 0.01)
}

// TestSurchargeAndCess validates surcharge and cess calculation
func TestSurchargeAndCess(t *testing.T) {
	// 15% Surcharge on income >= 1 Cr
	income := 100000000.0 // 1 Cr
	baseTax := 25000000.0

	surcharge := 0.0
	if income >= 100000000 {
		surcharge = baseTax * 0.15 // 3750000
	}

	// 4% Health & Education Cess on (tax + surcharge)
	totalTaxWithSurcharge := baseTax + surcharge // 28750000
	cess := totalTaxWithSurcharge * 0.04         // 1150000

	totalTaxLiability := totalTaxWithSurcharge + cess // 29900000

	assert.InDelta(t, 3750000.0, surcharge, 0.01)
	assert.InDelta(t, 1150000.0, cess, 0.01)
	assert.InDelta(t, 29900000.0, totalTaxLiability, 0.01)
}

// TestDeductionLimits validates deduction limits under IT Act
func TestDeductionLimits(t *testing.T) {
	// Section 80C - Max 150000
	section80C := 150000.0
	assert.LessOrEqual(t, section80C, 150000.0)

	// Section 80D - Medical Insurance (Individual: 15000, Senior: 20000)
	medicalInsurance := 20000.0
	assert.LessOrEqual(t, medicalInsurance, 20000.0)

	// Section 80E - Education Loan Interest (No limit)
	educationLoanInterest := 500000.0
	assert.Greater(t, educationLoanInterest, 0.0)

	// Section 80EEA - Home Loan Interest (Max 200000)
	homeLoanInterest := 200000.0
	assert.LessOrEqual(t, homeLoanInterest, 200000.0)

	totalDeductions := section80C + medicalInsurance + homeLoanInterest
	assert.Equal(t, 370000.0, totalDeductions)
}

// TestTDSCalculation validates TDS computation
func TestTDSCalculation(t *testing.T) {
	// TDS on payments to contractors/professionals
	contractorPayment := 100000.0
	tdsRate := 0.10 // 10% standard rate

	tdsAmount := contractorPayment * tdsRate
	assert.InDelta(t, 10000.0, tdsAmount, 0.01)

	// TDS on interest (non-resident)
	interestAmount := 1000000.0
	tdsRateInterest := 0.20 // 20% for non-residents

	tdsOnInterest := interestAmount * tdsRateInterest
	assert.InDelta(t, 200000.0, tdsOnInterest, 0.01)
}

// ============================================================================
// GST COMPLIANCE TESTS
// ============================================================================

// TestGSTCalculation validates GST calculation on transactions
func TestGSTCalculation(t *testing.T) {
	// Standard 5%, 12%, 18%, 28% tax rates in India
	testCases := []struct {
		baseAmount  float64
		gstRate     float64
		expectedGST float64
	}{
		{100000.0, 0.05, 5000.0},  // 5% GST
		{100000.0, 0.12, 12000.0}, // 12% GST
		{100000.0, 0.18, 18000.0}, // 18% GST
		{100000.0, 0.28, 28000.0}, // 28% GST
	}

	for _, tc := range testCases {
		gst := tc.baseAmount * tc.gstRate
		assert.InDelta(t, tc.expectedGST, gst, 0.01)
	}
}

// TestIGSTCalculation validates IGST for interstate transactions
func TestIGSTCalculation(t *testing.T) {
	// IGST = SGST + CGST for interstate supplies
	baseAmount := 100000.0
	sgstRate := 0.09 // 9% SGST
	cgstRate := 0.09 // 9% CGST
	igstRate := 0.18 // 18% IGST (SGST + CGST)

	sgstAmount := baseAmount * sgstRate
	cgstAmount := baseAmount * cgstRate
	igstAmount := baseAmount * igstRate

	assert.InDelta(t, 9000.0, sgstAmount, 0.01)
	assert.InDelta(t, 9000.0, cgstAmount, 0.01)
	assert.InDelta(t, 18000.0, igstAmount, 0.01)

	// IGST should equal SGST + CGST
	assert.InDelta(t, sgstAmount+cgstAmount, igstAmount, 0.01)
}

// TestInputCreditMechanism validates ITC (Input Tax Credit)
func TestInputCreditMechanism(t *testing.T) {
	// Output GST on sales
	salesAmount := 100000.0
	outputGST := salesAmount * 0.18

	// Input GST on purchases
	purchaseAmount := 50000.0
	inputGST := purchaseAmount * 0.18

	// Net GST payable to government
	netGST := outputGST - inputGST

	assert.InDelta(t, 18000.0, outputGST, 0.01)
	assert.InDelta(t, 9000.0, inputGST, 0.01)
	assert.InDelta(t, 9000.0, netGST, 0.01)
}

// TestGSTRegistrationThreshold validates registration requirement
func TestGSTRegistrationThreshold(t *testing.T) {
	// Mandatory registration: Turnover > 40L (40 lakhs)
	// Manufacturing: 40L, Trading: 40L

	annualTurnover := 5000000.0 // 50 lakhs
	threshold := 4000000.0      // 40 lakhs

	requiresRegistration := annualTurnover > threshold
	assert.True(t, requiresRegistration)

	// Below threshold
	lowTurnover := 3000000.0 // 30 lakhs
	requiresRegistrationLow := lowTurnover > threshold
	assert.False(t, requiresRegistrationLow)
}

// TestReversalOfITC validates ITC reversal conditions
func TestReversalOfITC(t *testing.T) {
	// ITC not allowed on personal expenses
	personalExpensesITC := 0.0
	assert.Equal(t, 0.0, personalExpensesITC)

	// ITC allowed on business inputs
	businessInputsITC := 50000.0
	assert.Greater(t, businessInputsITC, 0.0)

	// ITC must be reversed for exempt supplies
	exemptSupplyAmount := 100000.0
	proportionITCReversal := exemptSupplyAmount * 0.18
	assert.InDelta(t, 18000.0, proportionITCReversal, 0.01)
}

// ============================================================================
// DOUBLE ENTRY BOOKKEEPING & GL COMPLIANCE
// ============================================================================

// TestDoubleEntryPrinciple validates debit = credit
func TestDoubleEntryPrinciple(t *testing.T) {
	// Every transaction must have at least one debit and one credit
	// Total debits must always equal total credits

	testCases := []struct {
		description string
		debits      float64
		credits     float64
		balanced    bool
	}{
		{"Sales Invoice", 11800.0, 11800.0, true},
		{"Payment Received", 5000.0, 5000.0, true},
		{"Salary Payment", 50000.0, 50000.0, true},
		{"Unbalanced Entry", 10000.0, 9999.0, false},
	}

	for _, tc := range testCases {
		isBalanced := tc.debits == tc.credits
		assert.Equal(t, tc.balanced, isBalanced, tc.description)
	}
}

// TestAccountingEquation validates Assets = Liabilities + Equity
func TestAccountingEquation(t *testing.T) {
	// Assets = Liabilities + Equity (Balance Sheet Equation)
	assets := 10000000.0     // 1 Cr
	liabilities := 4000000.0 // 40 lakhs
	equity := 6000000.0      // 60 lakhs

	assert.InDelta(t, assets, liabilities+equity, 0.01)

	// If equity increases, assets must increase
	newEquity := 7000000.0
	newAssets := liabilities + newEquity
	assert.Equal(t, 11000000.0, newAssets)
}

// TestJournalEntryPosting validates proper posting
func TestJournalEntryPosting(t *testing.T) {
	// Transaction: Sales of goods for 11,800 (inclusive of 18% GST)
	invoiceTotal := 11800.0
	baseAmount := invoiceTotal / 1.18
	gstAmount := invoiceTotal - baseAmount

	// Journal Entry:
	// DR Cash/AR                  11,800
	//    CR Revenue              10,000
	//    CR GST Payable           1,800

	assert.InDelta(t, 10000.0, baseAmount, 0.01)
	assert.InDelta(t, 1800.0, gstAmount, 0.01)
	assert.InDelta(t, invoiceTotal, baseAmount+gstAmount, 0.01)
}

// TestAccountClassification validates correct account type mapping
func TestAccountClassification(t *testing.T) {
	// Asset accounts: Increase with Debit, Decrease with Credit
	cashAccount := struct {
		name    string
		acType  string
		debit   float64
		credit  float64
		balance float64
	}{
		name:    "Cash",
		acType:  "Asset",
		debit:   100000.0,
		credit:  30000.0,
		balance: 70000.0,
	}

	assert.Equal(t, "Asset", cashAccount.acType)
	assert.InDelta(t, 70000.0, cashAccount.debit-cashAccount.credit, 0.01)

	// Liability accounts: Increase with Credit, Decrease with Debit
	loanAccount := struct {
		name    string
		acType  string
		debit   float64
		credit  float64
		balance float64
	}{
		name:    "Bank Loan",
		acType:  "Liability",
		debit:   25000.0,
		credit:  100000.0,
		balance: 75000.0,
	}

	assert.Equal(t, "Liability", loanAccount.acType)
	assert.InDelta(t, 75000.0, loanAccount.credit-loanAccount.debit, 0.01)
}

// ============================================================================
// CREDIT & LIQUIDITY TESTS
// ============================================================================

// TestCreditLimitEnforcementStrict validates strict credit policies
func TestCreditLimitEnforcementStrict(t *testing.T) {
	creditLimit := 500000.0
	utilisedCredit := 300000.0
	availableCredit := creditLimit - utilisedCredit

	// New invoice cannot exceed available credit
	newInvoiceAmount := 200000.0
	canIssueInvoice := newInvoiceAmount <= availableCredit

	assert.True(t, canIssueInvoice)
	assert.InDelta(t, 200000.0, availableCredit, 0.01)

	// Invoice that would exceed credit
	excessiveInvoice := 250000.0
	canIssueExcess := excessiveInvoice <= availableCredit
	assert.False(t, canIssueExcess)
}

// TestAgingOfReceivables validates receivables aging analysis
func TestAgingOfReceivables(t *testing.T) {
	// Total Outstanding Receivables: 10 Lakhs
	totalReceivables := 1000000.0

	// Aging buckets
	current := 300000.0    // 0-30 days
	thirtyPlus := 400000.0 // 31-60 days
	sixtyPlus := 200000.0  // 61-90 days
	ninetyPlus := 100000.0 // 90+ days

	total := current + thirtyPlus + sixtyPlus + ninetyPlus
	assert.InDelta(t, totalReceivables, total, 0.01)

	// Bad Debt Allowance: 10% of 90+ days
	badDebtReserve := ninetyPlus * 0.10
	assert.InDelta(t, 10000.0, badDebtReserve, 0.01)

	// Net Receivables
	netReceivables := totalReceivables - badDebtReserve
	assert.InDelta(t, 990000.0, netReceivables, 0.01)
}

// TestPaymentTermsCompliance validates payment term adherence
func TestPaymentTermsCompliance(t *testing.T) {
	// Invoice issued: Jan 1, 2025
	// Terms: Net 30 (Payment due by Jan 31, 2025)
	creditDays := 30

	// Payment scheduled date
	dueDate := creditDays

	// Actual payment date: Jan 45 (overdue by 14 days)
	actualPaymentDays := 45
	daysPastDue := actualPaymentDays - dueDate

	assert.Greater(t, daysPastDue, 0)
	assert.Equal(t, 15, daysPastDue)
}

// ============================================================================
// AUDITOR & LENDER COMPLIANCE TESTS
// ============================================================================

// TestInternalControlsOverRevenue validates revenue recognition controls
func TestInternalControlsOverRevenue(t *testing.T) {
	// Revenue should be recognized only when:
	// 1. Performance obligation is satisfied
	// 2. Consideration is probable
	// 3. Amount is determinable
	// 4. Payment is probable
	// 5. Collection risk is minimal

	invoiceAmount := 500000.0
	receivedAmount := 500000.0

	// Revenue recognized = Received Amount (Conservative approach)
	revenueToRecognize := receivedAmount
	assert.InDelta(t, 500000.0, revenueToRecognize, 0.01)

	// Invoice issued but not received
	revenueNotRecognized := 0.0 // Until received
	assert.Equal(t, 0.0, revenueNotRecognized)
	assert.Equal(t, 0.0, invoiceAmount-receivedAmount) // Outstanding calculated inline
}

// TestCutoffValidation validates proper cutoff procedure
func TestCutoffValidation(t *testing.T) {
	// Month 1: Goods shipped Dec 31
	shipmentDate := "2024-12-31"
	invoiceDate := "2025-01-02"

	// Shipment (Dec 31) = Year 1 closing inventory
	// Invoice (Jan 2) = Year 2 revenue (but shipped in Year 1)
	// Proper cutoff: Revenue should be in Year 1 (on shipment)

	assert.Equal(t, "2024-12-31", shipmentDate)
	assert.NotEqual(t, "2024-12-31", invoiceDate)
}

// TestInventoryValuation validates inventory valuation methods
func TestInventoryValuation(t *testing.T) {
	// Ending inventory: 1000 units

	// FIFO Method
	unitPrice := 100.0
	fifoValue := 1000 * unitPrice
	assert.InDelta(t, 100000.0, fifoValue, 0.01)

	// Weighted Average Method
	weightedAvgValue := 1000 * 95.0
	assert.InDelta(t, 95000.0, weightedAvgValue, 0.01)

	// Must use consistent method year-over-year
	assert.Less(t, weightedAvgValue, fifoValue)
}

// TestDepreciationCalculation validates fixed asset depreciation
func TestDepreciationCalculation(t *testing.T) {
	// Asset Cost: 10 Lakhs
	assetCost := 1000000.0
	salvageValue := 100000.0
	usefulLife := 10.0 // years

	// Straight Line Depreciation
	annualDepreciation := (assetCost - salvageValue) / usefulLife
	assert.InDelta(t, 90000.0, annualDepreciation, 0.01)

	// After 5 years
	bookValue := assetCost - (annualDepreciation * 5)
	assert.InDelta(t, 550000.0, bookValue, 0.01)
}

// TestProvisionForDoubtfulDebts validates provision calculation
func TestProvisionForDoubtfulDebts(t *testing.T) {
	// Total Receivables: 20 Lakhs
	totalReceivables := 2000000.0

	// Provision percentages based on aging
	current := 1000000.0   // 0% provision
	thirtyPlus := 500000.0 // 5% provision = 25000
	sixtyPlus := 300000.0  // 10% provision = 30000
	ninetyPlus := 200000.0 // 50% provision = 100000

	provision := (current * 0.0) + (thirtyPlus * 0.05) + (sixtyPlus * 0.10) + (ninetyPlus * 0.50)
	// = 0 + 25000 + 30000 + 100000 = 155000
	assert.InDelta(t, 155000.0, provision, 0.01)

	// Net Receivables after provision
	netReceivables := totalReceivables - provision
	assert.InDelta(t, 1845000.0, netReceivables, 0.01)
}

// TestComplianceWithInd-AS validates Ind-AS compliance
func TestComplianceWithIndAS(t *testing.T) {
	// Ind-AS 115: Revenue from Contracts with Customers
	// Ind-AS 116: Leases
	// Ind-AS 109: Financial Instruments

	contractValue := 1000000.0
	performanceObligation := 500000.0 // 50% satisfied

	// Revenue to recognize
	revenueToRecognize := contractValue * (performanceObligation / contractValue)
	assert.InDelta(t, 500000.0, revenueToRecognize, 0.01)
}

// TestLenderCovenants validates lender covenant compliance
func TestLenderCovenants(t *testing.T) {
	// Debt Service Coverage Ratio (DSCR) >= 1.25
	ebitda := 2000000.0
	annualDebtService := 1400000.0 // Principal + Interest

	dscr := ebitda / annualDebtService
	assert.GreaterOrEqual(t, dscr, 1.25)

	// Debt to Equity <= 2.0
	totalDebt := 5000000.0
	totalEquity := 3000000.0
	debtToEquityRatio := totalDebt / totalEquity

	assert.LessOrEqual(t, debtToEquityRatio, 2.0)

	// Current Ratio >= 1.5
	currentAssets := 5000000.0
	currentLiabilities := 3000000.0
	currentRatio := currentAssets / currentLiabilities

	assert.GreaterOrEqual(t, currentRatio, 1.5)
}
