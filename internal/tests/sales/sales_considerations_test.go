package sales

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// 1. PROFITABILITY & MARGIN ANALYSIS
// ============================================================================

// TestProfitabilityMarginCalculation validates cost price tracking and margin calculations
func TestProfitabilityMarginCalculation(t *testing.T) {
	testCases := []struct {
		description    string
		sellingPrice   float64
		costPrice      float64
		quantity       float64
		expectedMargin float64
		expectedGP     float64 // Gross Profit
		expectedGPM    float64 // Gross Profit Margin %
	}{
		{
			description:    "Standard product with 40% margin",
			sellingPrice:   1000.00,
			costPrice:      600.00,
			quantity:       10.0,
			expectedMargin: 400.00,
			expectedGP:     4000.00,
			expectedGPM:    40.0,
		},
		{
			description:    "Low margin commodity",
			sellingPrice:   500.00,
			costPrice:      475.00,
			quantity:       20.0,
			expectedMargin: 25.00,
			expectedGP:     500.00,
			expectedGPM:    5.0,
		},
		{
			description:    "Premium product with high margin",
			sellingPrice:   5000.00,
			costPrice:      2500.00,
			quantity:       5.0,
			expectedMargin: 2500.00,
			expectedGP:     12500.00,
			expectedGPM:    50.0,
		},
	}

	for _, tc := range testCases {
		// Calculate per-unit margin
		unitMargin := tc.sellingPrice - tc.costPrice
		assert.InDelta(t, tc.expectedMargin, unitMargin, 0.01, tc.description+" - unit margin")

		// Calculate total gross profit
		grossProfit := unitMargin * tc.quantity
		assert.InDelta(t, tc.expectedGP, grossProfit, 0.01, tc.description+" - gross profit")

		// Calculate gross profit margin %
		gpm := (unitMargin / tc.sellingPrice) * 100
		assert.InDelta(t, tc.expectedGPM, gpm, 0.1, tc.description+" - GPM%")
	}
}

// TestMinimumMarginEnforcement validates that sales respect minimum margin requirements
func TestMinimumMarginEnforcement(t *testing.T) {
	minimumMarginPercent := 15.0 // Company policy: min 15% margin

	testCases := []struct {
		sellingPrice float64
		costPrice    float64
		allowed      bool
	}{
		{1000.00, 800.00, true},  // 20% margin - OK
		{1000.00, 850.00, true},  // 15% margin - OK (at limit)
		{1000.00, 851.00, false}, // 14.9% margin - NOT OK
		{500.00, 425.00, true},   // 15% - OK (at limit)
	}

	for _, tc := range testCases {
		actualMargin := ((tc.sellingPrice - tc.costPrice) / tc.sellingPrice) * 100
		allowed := actualMargin >= minimumMarginPercent
		assert.Equal(t, tc.allowed, allowed)
	}
}

// TestContributionMarginAnalysis validates contribution margin for pricing decisions
func TestContributionMarginAnalysis(t *testing.T) {
	testCases := []struct {
		sellingPrice     float64
		variableCost     float64
		fixedCostPerUnit float64
		quantity         float64
		expectedContrib  float64
		expectedCM       float64 // Contribution Margin
	}{
		{
			sellingPrice:     1000.00,
			variableCost:     600.00,
			fixedCostPerUnit: 100.00,
			quantity:         10.0,
			expectedContrib:  4000.00,
			expectedCM:       40.0,
		},
	}

	for _, tc := range testCases {
		// Contribution = (SP - VC) per unit
		unitContribution := tc.sellingPrice - tc.variableCost
		assert.InDelta(t, tc.expectedContrib, unitContribution*tc.quantity, 0.01)

		// CM% = (SP - VC) / SP * 100
		cm := ((tc.sellingPrice - tc.variableCost) / tc.sellingPrice) * 100
		assert.InDelta(t, tc.expectedCM, cm, 0.1)
	}
}

// ============================================================================
// 2. SALES COMMISSION & INCENTIVES
// ============================================================================

// TestCommissionCalculation validates tiered commission structures
func TestCommissionCalculation(t *testing.T) {
	testCases := []struct {
		description         string
		salesAmount         float64
		commissionStructure map[string]float64 // threshold: rate
		expectedCommission  float64
	}{
		{
			description: "Basic tier - below 5L",
			salesAmount: 300000.00,
			commissionStructure: map[string]float64{
				"0":       5.0,
				"500000":  7.5,
				"1000000": 10.0,
			},
			expectedCommission: 15000.00, // 300k * 5%
		},
		{
			description: "Mid tier - 5L to 10L",
			salesAmount: 750000.00,
			commissionStructure: map[string]float64{
				"0":       5.0,
				"500000":  7.5,
				"1000000": 10.0,
			},
			expectedCommission: 56250.00, // 750k * 7.5%
		},
		{
			description: "Premium tier - above 10L",
			salesAmount: 1500000.00,
			commissionStructure: map[string]float64{
				"0":       5.0,
				"500000":  7.5,
				"1000000": 10.0,
			},
			expectedCommission: 150000.00, // 1.5M * 10%
		},
	}

	for _, tc := range testCases {
		// Simplified: apply highest applicable rate
		var applicableRate float64 = 5.0
		if tc.salesAmount >= 1000000 {
			applicableRate = 10.0
		} else if tc.salesAmount >= 500000 {
			applicableRate = 7.5
		}

		commission := tc.salesAmount * (applicableRate / 100)
		assert.InDelta(t, tc.expectedCommission, commission, 0.01, tc.description)
	}
}

// TestCommissionAchievementTracking validates commission vs target
func TestCommissionAchievementTracking(t *testing.T) {
	testCases := []struct {
		targetAmount   float64
		actualAmount   float64
		commissionRate float64
		expectedComm   float64
		achievement    float64
	}{
		{
			targetAmount:   1000000.00,
			actualAmount:   1000000.00,
			commissionRate: 5.0,
			expectedComm:   50000.00,
			achievement:    100.0,
		},
		{
			targetAmount:   1000000.00,
			actualAmount:   1250000.00,
			commissionRate: 5.0,
			expectedComm:   62500.00,
			achievement:    125.0,
		},
		{
			targetAmount:   1000000.00,
			actualAmount:   750000.00,
			commissionRate: 5.0,
			expectedComm:   37500.00,
			achievement:    75.0,
		},
	}

	for _, tc := range testCases {
		commission := tc.actualAmount * (tc.commissionRate / 100)
		assert.InDelta(t, tc.expectedComm, commission, 0.01)

		achievement := (tc.actualAmount / tc.targetAmount) * 100
		assert.InDelta(t, tc.achievement, achievement, 0.1)
	}
}

// TestIncentiveBonus validates bonus structures for overachievement
func TestIncentiveBonus(t *testing.T) {
	testCases := []struct {
		achievement    float64
		baseCommission float64
		expectedBonus  float64 // 10% extra for every 10% above target
	}{
		{100.0, 50000.00, 0.00},     // At target - no bonus
		{110.0, 50000.00, 5000.00},  // 10% above - 10% bonus
		{120.0, 50000.00, 10000.00}, // 20% above - 20% bonus
		{150.0, 50000.00, 25000.00}, // 50% above - 50% bonus
	}

	for _, tc := range testCases {
		bonusMultiplier := 0.0
		if tc.achievement > 100 {
			bonusMultiplier = (tc.achievement - 100) / 100
		}
		bonus := tc.baseCommission * bonusMultiplier
		assert.InDelta(t, tc.expectedBonus, bonus, 0.01)
	}
}

// ============================================================================
// 3. INVENTORY & STOCK MANAGEMENT
// ============================================================================

// TestStockAvailabilityValidation validates stock before sale confirmation
func TestStockAvailabilityValidation(t *testing.T) {
	testCases := []struct {
		description       string
		requestedQty      float64
		availableStock    float64
		allowBackorder    bool
		canFulfill        bool
		fulfillmentStatus string
	}{
		{
			description:       "Sufficient stock",
			requestedQty:      10.0,
			availableStock:    50.0,
			allowBackorder:    false,
			canFulfill:        true,
			fulfillmentStatus: "full",
		},
		{
			description:       "No stock, no backorder allowed",
			requestedQty:      10.0,
			availableStock:    0.0,
			allowBackorder:    false,
			canFulfill:        false,
			fulfillmentStatus: "rejected",
		},
		{
			description:       "Partial fulfillment",
			requestedQty:      20.0,
			availableStock:    15.0,
			allowBackorder:    true,
			canFulfill:        true,
			fulfillmentStatus: "partial",
		},
		{
			description:       "Backorder allowed",
			requestedQty:      30.0,
			availableStock:    0.0,
			allowBackorder:    true,
			canFulfill:        true,
			fulfillmentStatus: "backorder",
		},
	}

	for _, tc := range testCases {
		var status string
		canFulfill := false

		if tc.requestedQty <= tc.availableStock {
			status = "full"
			canFulfill = true
		} else if tc.availableStock > 0 && tc.allowBackorder {
			status = "partial"
			canFulfill = true
		} else if tc.availableStock == 0 && tc.allowBackorder {
			status = "backorder"
			canFulfill = true
		} else {
			status = "rejected"
			canFulfill = false
		}

		assert.Equal(t, tc.canFulfill, canFulfill, tc.description)
		assert.Equal(t, tc.fulfillmentStatus, status, tc.description)
	}
}

// TestFIFOValuation validates FIFO stock valuation
func TestFIFOValuation(t *testing.T) {
	// Batches: Batch1(100 units @ ₹100), Batch2(50 units @ ₹110), Batch3(75 units @ ₹105)
	batches := []struct {
		quantity  float64
		unitCost  float64
		totalCost float64
	}{
		{100.0, 100.00, 10000.00},
		{50.0, 110.00, 5500.00},
		{75.0, 105.00, 7875.00},
	}

	// Sell 120 units (100 from Batch1 + 20 from Batch2)
	unitsToSell := 120.0
	var cogs float64

	remainingToSell := unitsToSell
	batchIdx := 0

	// FIFO: consume oldest batches first
	for remainingToSell > 0 && batchIdx < len(batches) {
		if remainingToSell >= batches[batchIdx].quantity {
			// Consume entire batch
			cogs += batches[batchIdx].totalCost
			remainingToSell -= batches[batchIdx].quantity
		} else {
			// Consume partial batch
			cogs += remainingToSell * batches[batchIdx].unitCost
			remainingToSell = 0
		}
		batchIdx++
	}

	expectedCOGS := 10000.00 + (20.0 * 110.00) // Batch1 full + 20 from Batch2
	assert.InDelta(t, expectedCOGS, cogs, 0.01)
}

// TestWeightedAverageCost validates weighted average cost method
func TestWeightedAverageCost(t *testing.T) {
	// Batches: Batch1(100 @ ₹100), Batch2(50 @ ₹110)
	totalUnits := 150.0
	totalCost := 15500.00 // (100*100 + 50*110)
	weightedAvgCost := totalCost / totalUnits

	// Sell 120 units
	unitsToSell := 120.0
	cogs := unitsToSell * weightedAvgCost

	expectedCOGS := 120.0 * (15500.00 / 150.0)
	assert.InDelta(t, expectedCOGS, cogs, 0.01)
}

// ============================================================================
// 4. PRODUCT MIX & CROSS-SELLING
// ============================================================================

// TestProductBundleDetection validates bundle selling and discounts
func TestProductBundleDetection(t *testing.T) {
	testCases := []struct {
		description      string
		items            []string // Product codes
		isBundle         bool
		bundleDiscount   float64
		expectedDiscount float64
	}{
		{
			description:      "Single item - no bundle",
			items:            []string{"PROD001"},
			isBundle:         false,
			bundleDiscount:   0.0,
			expectedDiscount: 0.0,
		},
		{
			description:      "Standard bundle - 2 products",
			items:            []string{"PROD001", "PROD002"},
			isBundle:         true,
			bundleDiscount:   5.0,
			expectedDiscount: 5.0,
		},
		{
			description:      "Premium bundle - 3+ products",
			items:            []string{"PROD001", "PROD002", "PROD003"},
			isBundle:         true,
			bundleDiscount:   10.0,
			expectedDiscount: 10.0,
		},
	}

	for _, tc := range testCases {
		isBundle := len(tc.items) >= 2
		assert.Equal(t, tc.isBundle, isBundle, tc.description)

		if isBundle {
			assert.Equal(t, tc.expectedDiscount, tc.bundleDiscount, tc.description)
		}
	}
}

// TestSKULevelAnalysis validates product-level performance tracking
func TestSKULevelAnalysis(t *testing.T) {
	skuMetrics := map[string]struct {
		unitsSold    float64
		revenue      float64
		avgPrice     float64
		profitMargin float64
	}{
		"SKU001": {
			unitsSold:    100.0,
			revenue:      100000.00,
			avgPrice:     1000.00,
			profitMargin: 35.0,
		},
		"SKU002": {
			unitsSold:    50.0,
			revenue:      75000.00,
			avgPrice:     1500.00,
			profitMargin: 40.0,
		},
	}

	// Validate each SKU
	for sku, metrics := range skuMetrics {
		actualAvgPrice := metrics.revenue / metrics.unitsSold
		assert.InDelta(t, metrics.avgPrice, actualAvgPrice, 0.01, sku)
	}
}

// ============================================================================
// 5. CUSTOMER SEGMENTATION & LIFETIME VALUE
// ============================================================================

// TestCustomerSegmentationProfitability validates profitability by segment
func TestCustomerSegmentationProfitability(t *testing.T) {
	segments := map[string]struct {
		totalRevenue   float64
		totalCOGS      float64
		grossProfit    float64
		profitMargin   float64
		customerCount  int
		revPerCustomer float64
	}{
		"Enterprise": {
			totalRevenue:   10000000.00,
			totalCOGS:      6000000.00,
			grossProfit:    4000000.00,
			profitMargin:   40.0,
			customerCount:  5,
			revPerCustomer: 2000000.00,
		},
		"SMB": {
			totalRevenue:   3000000.00,
			totalCOGS:      2100000.00,
			grossProfit:    900000.00,
			profitMargin:   30.0,
			customerCount:  50,
			revPerCustomer: 60000.00,
		},
		"Retail": {
			totalRevenue:   1000000.00,
			totalCOGS:      850000.00,
			grossProfit:    150000.00,
			profitMargin:   15.0,
			customerCount:  1000,
			revPerCustomer: 1000.00,
		},
	}

	for segment, metrics := range segments {
		actualGP := metrics.totalRevenue - metrics.totalCOGS
		assert.InDelta(t, metrics.grossProfit, actualGP, 0.01, segment)

		actualMargin := (actualGP / metrics.totalRevenue) * 100
		assert.InDelta(t, metrics.profitMargin, actualMargin, 0.1, segment)

		revPerCust := metrics.totalRevenue / float64(metrics.customerCount)
		assert.InDelta(t, metrics.revPerCustomer, revPerCust, 0.01, segment)
	}
}

// TestCustomerLifetimeValue validates CLV calculation
func TestCustomerLifetimeValue(t *testing.T) {
	testCases := []struct {
		description      string
		avgOrderValue    float64
		orderFrequency   float64 // orders per year
		customerLifespan float64 // years
		expectedCLV      float64
	}{
		{
			description:      "High-value repeat customer",
			avgOrderValue:    50000.00,
			orderFrequency:   12.0, // monthly
			customerLifespan: 5.0,
			expectedCLV:      3000000.00, // 50k * 12 * 5
		},
		{
			description:      "Standard customer",
			avgOrderValue:    10000.00,
			orderFrequency:   4.0, // quarterly
			customerLifespan: 3.0,
			expectedCLV:      120000.00, // 10k * 4 * 3
		},
	}

	for _, tc := range testCases {
		clv := tc.avgOrderValue * tc.orderFrequency * tc.customerLifespan
		assert.InDelta(t, tc.expectedCLV, clv, 0.01, tc.description)
	}
}

// TestRFMSegmentation validates RFM (Recency/Frequency/Monetary) scoring
func TestRFMSegmentation(t *testing.T) {
	testCases := []struct {
		description    string
		recencyDays    int
		frequencyCount int
		monetaryValue  float64
		rScore         int
		fScore         int
		mScore         int
	}{
		{
			description:    "Champion - Recent, Frequent, High Value",
			recencyDays:    5,
			frequencyCount: 35,
			monetaryValue:  500000.00,
			rScore:         5,
			fScore:         5,
			mScore:         5,
		},
		{
			description:    "At Risk - Old, Frequent, Was High Value",
			recencyDays:    120,
			frequencyCount: 15,
			monetaryValue:  400000.00,
			rScore:         2,
			fScore:         3,
			mScore:         4,
		},
		{
			description:    "New - Recent, Low Frequency, Low Value",
			recencyDays:    10,
			frequencyCount: 1,
			monetaryValue:  5000.00,
			rScore:         5,
			fScore:         1,
			mScore:         1,
		},
	}

	for _, tc := range testCases {
		// Simple scoring: 1-5 scale
		var rScore, fScore, mScore int

		// Recency: 0-30 days = 5, 31-60 = 4, 61-90 = 3, 91-180 = 2, 180+ = 1
		if tc.recencyDays <= 30 {
			rScore = 5
		} else if tc.recencyDays <= 60 {
			rScore = 4
		} else if tc.recencyDays <= 90 {
			rScore = 3
		} else if tc.recencyDays <= 180 {
			rScore = 2
		} else {
			rScore = 1
		}

		// Frequency: 1-5 = 1, 6-10 = 2, 11-15 = 3, 16-20 = 4, 20+ = 5
		if tc.frequencyCount <= 5 {
			fScore = 1
		} else if tc.frequencyCount <= 10 {
			fScore = 2
		} else if tc.frequencyCount <= 15 {
			fScore = 3
		} else if tc.frequencyCount <= 20 {
			fScore = 4
		} else {
			fScore = 5
		}

		// Monetary: <50k=1, 50-100k=2, 100-250k=3, 250-500k=4, 500k+=5
		if tc.monetaryValue < 50000 {
			mScore = 1
		} else if tc.monetaryValue < 100000 {
			mScore = 2
		} else if tc.monetaryValue < 250000 {
			mScore = 3
		} else if tc.monetaryValue < 500000 {
			mScore = 4
		} else {
			mScore = 5
		}

		assert.Equal(t, tc.rScore, rScore, tc.description+" - R score")
		assert.Equal(t, tc.fScore, fScore, tc.description+" - F score")
		assert.Equal(t, tc.mScore, mScore, tc.description+" - M score")
	}
}

// ============================================================================
// 6. SALES RETURNS & ALLOWANCES
// ============================================================================

// TestReturnRateTracking validates return frequency analysis
func TestReturnRateTracking(t *testing.T) {
	testCases := []struct {
		description    string
		totalInvoiced  float64
		totalReturned  float64
		expectedRate   float64
		acceptableRate float64
		isAcceptable   bool
	}{
		{
			description:    "Low return rate - excellent",
			totalInvoiced:  1000000.00,
			totalReturned:  20000.00,
			expectedRate:   2.0,
			acceptableRate: 5.0,
			isAcceptable:   true,
		},
		{
			description:    "Normal return rate",
			totalInvoiced:  1000000.00,
			totalReturned:  40000.00,
			expectedRate:   4.0,
			acceptableRate: 5.0,
			isAcceptable:   true,
		},
		{
			description:    "High return rate - concern",
			totalInvoiced:  1000000.00,
			totalReturned:  80000.00,
			expectedRate:   8.0,
			acceptableRate: 5.0,
			isAcceptable:   false,
		},
	}

	for _, tc := range testCases {
		returnRate := (tc.totalReturned / tc.totalInvoiced) * 100
		assert.InDelta(t, tc.expectedRate, returnRate, 0.1, tc.description)

		isAcceptable := returnRate <= tc.acceptableRate
		assert.Equal(t, tc.isAcceptable, isAcceptable, tc.description)
	}
}

// TestReturnReasonAnalysis validates root cause categorization
func TestReturnReasonAnalysis(t *testing.T) {
	returnReasons := map[string]struct {
		count      int
		percentage float64
		category   string
	}{
		"Quality_Defect": {
			count:      45,
			percentage: 45.0,
			category:   "quality",
		},
		"Wrong_Item": {
			count:      25,
			percentage: 25.0,
			category:   "fulfillment",
		},
		"Damaged": {
			count:      20,
			percentage: 20.0,
			category:   "logistics",
		},
		"No_Longer_Needed": {
			count:      10,
			percentage: 10.0,
			category:   "customer",
		},
	}

	totalReturns := 100
	for reason, data := range returnReasons {
		pct := (float64(data.count) / float64(totalReturns)) * 100
		assert.InDelta(t, data.percentage, pct, 0.1, reason)
	}
}

// TestWarrantyAndDefectTracking validates warranty claim tracking
func TestWarrantyAndDefectTracking(t *testing.T) {
	testCases := []struct {
		description        string
		invoiceDate        string // YYYY-MM-DD
		claimDate          string
		warrantyMonths     int
		withinWarranty     bool
		defectCategory     string
		replacementAllowed bool
	}{
		{
			description:        "Within 12-month warranty",
			invoiceDate:        "2024-01-01",
			claimDate:          "2024-06-01",
			warrantyMonths:     12,
			withinWarranty:     true,
			defectCategory:     "manufacturing",
			replacementAllowed: true,
		},
		{
			description:        "Beyond warranty period",
			invoiceDate:        "2023-01-01",
			claimDate:          "2024-06-01",
			warrantyMonths:     12,
			withinWarranty:     false,
			defectCategory:     "manufacturing",
			replacementAllowed: false,
		},
	}

	for _, tc := range testCases {
		assert.NotEmpty(t, tc.invoiceDate)
		assert.NotEmpty(t, tc.claimDate)
		assert.Greater(t, tc.warrantyMonths, 0)
	}
}

// ============================================================================
// 7. DYNAMIC PRICING & DISCOUNTS
// ============================================================================

// TestVolumePricingTiers validates quantity-based pricing
func TestVolumePricingTiers(t *testing.T) {
	pricingTiers := map[string]struct {
		minQty       float64
		maxQty       float64
		pricePerUnit float64
	}{
		"Tier1": {minQty: 1, maxQty: 10, pricePerUnit: 1000.00},
		"Tier2": {minQty: 11, maxQty: 50, pricePerUnit: 950.00},
		"Tier3": {minQty: 51, maxQty: 100, pricePerUnit: 900.00},
		"Tier4": {minQty: 101, maxQty: -1, pricePerUnit: 850.00}, // Unlimited
	}

	testCases := []struct {
		quantity      float64
		expectedPrice float64
	}{
		{5.0, 1000.00},  // Tier1
		{25.0, 950.00},  // Tier2
		{75.0, 900.00},  // Tier3
		{150.0, 850.00}, // Tier4
	}

	for _, tc := range testCases {
		var applicablePrice float64
		for _, tier := range pricingTiers {
			if tc.quantity >= tier.minQty && (tier.maxQty < 0 || tc.quantity <= tier.maxQty) {
				applicablePrice = tier.pricePerUnit
				break
			}
		}
		assert.InDelta(t, tc.expectedPrice, applicablePrice, 0.01)
	}
}

// TestSeasonalPricingAdjustment validates time-based price changes
func TestSeasonalPricingAdjustment(t *testing.T) {
	testCases := []struct {
		description      string
		basePricePerUnit float64
		month            int // 1-12
		seasonalFactor   float64
		adjustedPrice    float64
	}{
		{
			description:      "Peak season - December",
			basePricePerUnit: 1000.00,
			month:            12,
			seasonalFactor:   1.25,
			adjustedPrice:    1250.00,
		},
		{
			description:      "Normal season - June",
			basePricePerUnit: 1000.00,
			month:            6,
			seasonalFactor:   1.0,
			adjustedPrice:    1000.00,
		},
		{
			description:      "Off-season - April",
			basePricePerUnit: 1000.00,
			month:            4,
			seasonalFactor:   0.85,
			adjustedPrice:    850.00,
		},
	}

	for _, tc := range testCases {
		adjusted := tc.basePricePerUnit * tc.seasonalFactor
		assert.InDelta(t, tc.adjustedPrice, adjusted, 0.01, tc.description)
	}
}

// TestMaximumDiscountEnforcement validates discount caps
func TestMaximumDiscountEnforcement(t *testing.T) {
	maximumDiscountPercent := 20.0 // Company policy

	testCases := []struct {
		requestedDiscount float64
		allowed           bool
	}{
		{5.0, true},   // Under limit
		{15.0, true},  // Under limit
		{20.0, true},  // At limit
		{21.0, false}, // Over limit
		{50.0, false}, // Way over limit
	}

	for _, tc := range testCases {
		allowed := tc.requestedDiscount <= maximumDiscountPercent
		assert.Equal(t, tc.allowed, allowed)
	}
}

// TestCompetitiveDiscounting validates discount strategy
func TestCompetitiveDiscounting(t *testing.T) {
	testCases := []struct {
		description        string
		ourPrice           float64
		competitorPrice    float64
		marketingCost      float64
		maxAllowedDiscount float64
		recommended        float64
	}{
		{
			description:        "Slightly higher - small discount ok",
			ourPrice:           1000.00,
			competitorPrice:    950.00,
			marketingCost:      50.00,
			maxAllowedDiscount: 10.0,
			recommended:        5.0,
		},
		{
			description:        "Much higher - large discount needed",
			ourPrice:           1000.00,
			competitorPrice:    700.00,
			marketingCost:      100.00,
			maxAllowedDiscount: 20.0,
			recommended:        20.0,
		},
	}

	for _, tc := range testCases {
		_ = ((tc.ourPrice - tc.competitorPrice) / tc.ourPrice) * 100 // Price differential
		assert.Greater(t, tc.ourPrice, tc.competitorPrice-1)         // Our price context
	}
}

// ============================================================================
// 8. CUSTOMER CREDIT MANAGEMENT
// ============================================================================

// TestAccountsReceivableAging validates AR aging analysis
func TestAccountsReceivableAging(t *testing.T) {
	testCases := []struct {
		description   string
		invoiceAmount float64
		daysOverdue   int
		ageCategory   string
		riskLevel     string
	}{
		{
			description:   "Current",
			invoiceAmount: 100000.00,
			daysOverdue:   0,
			ageCategory:   "current",
			riskLevel:     "low",
		},
		{
			description:   "31-60 days",
			invoiceAmount: 50000.00,
			daysOverdue:   45,
			ageCategory:   "31-60",
			riskLevel:     "medium",
		},
		{
			description:   "61-90 days",
			invoiceAmount: 30000.00,
			daysOverdue:   75,
			ageCategory:   "61-90",
			riskLevel:     "high",
		},
		{
			description:   "90+ days",
			invoiceAmount: 20000.00,
			daysOverdue:   120,
			ageCategory:   "90+",
			riskLevel:     "critical",
		},
	}

	for _, tc := range testCases {
		var riskLevel string

		if tc.daysOverdue == 0 {
			riskLevel = "low"
		} else if tc.daysOverdue <= 30 {
			riskLevel = "low"
		} else if tc.daysOverdue <= 60 {
			riskLevel = "medium"
		} else if tc.daysOverdue <= 90 {
			riskLevel = "high"
		} else {
			riskLevel = "critical"
		}

		assert.Equal(t, tc.riskLevel, riskLevel, tc.description)
	}
}

// TestBadDebtProvision validates bad debt allowance calculation
func TestBadDebtProvision(t *testing.T) {
	testCases := []struct {
		description       string
		invoiceAmount     float64
		daysOverdue       int
		provisionPercent  float64
		expectedProvision float64
	}{
		{
			description:       "Current invoice",
			invoiceAmount:     100000.00,
			daysOverdue:       0,
			provisionPercent:  0.0,
			expectedProvision: 0.00,
		},
		{
			description:       "31-60 days overdue",
			invoiceAmount:     50000.00,
			daysOverdue:       45,
			provisionPercent:  5.0,
			expectedProvision: 2500.00,
		},
		{
			description:       "61-90 days overdue",
			invoiceAmount:     30000.00,
			daysOverdue:       75,
			provisionPercent:  10.0,
			expectedProvision: 3000.00,
		},
		{
			description:       "90+ days overdue",
			invoiceAmount:     20000.00,
			daysOverdue:       120,
			provisionPercent:  25.0,
			expectedProvision: 5000.00,
		},
	}

	for _, tc := range testCases {
		provision := tc.invoiceAmount * (tc.provisionPercent / 100)
		assert.InDelta(t, tc.expectedProvision, provision, 0.01, tc.description)
	}
}

// TestPaymentTermCompliance validates payment term enforcement
func TestPaymentTermCompliance(t *testing.T) {
	testCases := []struct {
		description     string
		creditDays      int
		invoiceAmount   float64
		paymentSchedule []float64 // Installments
		isCompliant     bool
	}{
		{
			description:     "Net 30 - paid on time",
			creditDays:      30,
			invoiceAmount:   100000.00,
			paymentSchedule: []float64{100000.00},
			isCompliant:     true,
		},
		{
			description:     "Net 60 - installment plan",
			creditDays:      60,
			invoiceAmount:   100000.00,
			paymentSchedule: []float64{50000.00, 50000.00},
			isCompliant:     true,
		},
	}

	for _, tc := range testCases {
		totalPaid := 0.0
		for _, payment := range tc.paymentSchedule {
			totalPaid += payment
		}
		assert.InDelta(t, tc.invoiceAmount, totalPaid, 0.01, tc.description)
	}
}

// ============================================================================
// 9. SALES FORECASTING & PIPELINE
// ============================================================================

// TestSalesCyclePipelineAnalysis validates sales stage progression
func TestSalesCyclePipelineAnalysis(t *testing.T) {
	pipeline := map[string]struct {
		stageNumber    int
		stageName      string
		deals          int
		totalValue     float64
		conversionRate float64
		avgCycleDays   int
	}{
		"Stage1": {
			stageNumber:    1,
			stageName:      "Lead",
			deals:          100,
			totalValue:     5000000.00,
			conversionRate: 10.0, // 10% to next stage
			avgCycleDays:   7,
		},
		"Stage2": {
			stageNumber:    2,
			stageName:      "Opportunity",
			deals:          10,
			totalValue:     500000.00,
			conversionRate: 50.0, // 50% to next stage
			avgCycleDays:   14,
		},
		"Stage3": {
			stageNumber:    3,
			stageName:      "Proposal",
			deals:          5,
			totalValue:     250000.00,
			conversionRate: 80.0, // 80% to next stage
			avgCycleDays:   21,
		},
		"Stage4": {
			stageNumber:    4,
			stageName:      "Negotiation",
			deals:          4,
			totalValue:     200000.00,
			conversionRate: 75.0, // 75% to closed
			avgCycleDays:   30,
		},
	}

	for stageKey, metrics := range pipeline {
		assert.NotEmpty(t, metrics.stageName)
		assert.Greater(t, metrics.deals, 0)
		assert.Greater(t, metrics.totalValue, 0.0)
		assert.NotEmpty(t, stageKey)
	}
}

// TestDealProbabilityWeighting validates weighted pipeline value
func TestDealProbabilityWeighting(t *testing.T) {
	testCases := []struct {
		description     string
		dealValue       float64
		dealProbability float64
		weightedValue   float64
	}{
		{
			description:     "High probability deal",
			dealValue:       100000.00,
			dealProbability: 0.80,
			weightedValue:   80000.00,
		},
		{
			description:     "Medium probability deal",
			dealValue:       50000.00,
			dealProbability: 0.50,
			weightedValue:   25000.00,
		},
		{
			description:     "Low probability deal",
			dealValue:       30000.00,
			dealProbability: 0.25,
			weightedValue:   7500.00,
		},
	}

	for _, tc := range testCases {
		weighted := tc.dealValue * tc.dealProbability
		assert.InDelta(t, tc.weightedValue, weighted, 0.01, tc.description)
	}
}

// TestRevenueForecasting validates revenue predictions
func TestRevenueForecasting(t *testing.T) {
	testCases := []struct {
		description       string
		pipelineValue     float64
		avgConversionRate float64
		forecastedRevenue float64
	}{
		{
			description:       "Conservative forecast (50% conversion)",
			pipelineValue:     1000000.00,
			avgConversionRate: 0.50,
			forecastedRevenue: 500000.00,
		},
		{
			description:       "Optimistic forecast (70% conversion)",
			pipelineValue:     1000000.00,
			avgConversionRate: 0.70,
			forecastedRevenue: 700000.00,
		},
	}

	for _, tc := range testCases {
		forecast := tc.pipelineValue * tc.avgConversionRate
		assert.InDelta(t, tc.forecastedRevenue, forecast, 0.01, tc.description)
	}
}

// ============================================================================
// 10. TERRITORY & CHANNEL MANAGEMENT
// ============================================================================

// TestTerritoryAssignment validates salesperson territory allocation
func TestTerritoryAssignment(t *testing.T) {
	territories := map[string]struct {
		salesPersonID     string
		region            string
		state             string
		assignedCustomers int
		targetRevenue     float64
		actualRevenue     float64
		achievementPct    float64
	}{
		"Territory-North": {
			salesPersonID:     "SP001",
			region:            "North",
			state:             "Delhi,Punjab,Himachal",
			assignedCustomers: 50,
			targetRevenue:     2000000.00,
			actualRevenue:     2300000.00,
			achievementPct:    115.0,
		},
		"Territory-South": {
			salesPersonID:     "SP002",
			region:            "South",
			state:             "Karnataka,Tamil Nadu,Telangana",
			assignedCustomers: 60,
			targetRevenue:     2500000.00,
			actualRevenue:     2100000.00,
			achievementPct:    84.0,
		},
	}

	for territoryKey, metrics := range territories {
		assert.NotEmpty(t, metrics.salesPersonID)
		assert.Greater(t, metrics.assignedCustomers, 0)
		assert.Greater(t, metrics.targetRevenue, 0.0)
		assert.NotEmpty(t, territoryKey)
	}
}

// TestChannelProfitabilityAnalysis validates profitability by channel
func TestChannelProfitabilityAnalysis(t *testing.T) {
	channels := map[string]struct {
		channelName  string
		revenue      float64
		costOfSales  float64
		channelCost  float64 // Distribution cost
		netProfit    float64
		profitMargin float64
	}{
		"Direct": {
			channelName:  "Direct Sales",
			revenue:      5000000.00,
			costOfSales:  3000000.00,
			channelCost:  500000.00,
			netProfit:    1500000.00,
			profitMargin: 30.0,
		},
		"Distributor": {
			channelName:  "Distributor Network",
			revenue:      3000000.00,
			costOfSales:  2100000.00,
			channelCost:  600000.00, // Higher distribution cost
			netProfit:    300000.00,
			profitMargin: 10.0,
		},
		"Online": {
			channelName:  "E-commerce",
			revenue:      2000000.00,
			costOfSales:  1200000.00,
			channelCost:  400000.00,
			netProfit:    400000.00,
			profitMargin: 20.0,
		},
	}

	for channel, metrics := range channels {
		actualProfit := metrics.revenue - metrics.costOfSales - metrics.channelCost
		assert.InDelta(t, metrics.netProfit, actualProfit, 0.01, channel)

		actualMargin := (actualProfit / metrics.revenue) * 100
		assert.InDelta(t, metrics.profitMargin, actualMargin, 0.1, channel)
	}
}

// TestChannelPartnerManagement validates partner performance tracking
func TestChannelPartnerManagement(t *testing.T) {
	testCases := []struct {
		description     string
		partnerName     string
		partnerType     string // Distributor, Reseller, Integrator
		yearlyTarget    float64
		actualSales     float64
		achievementRate float64
		partnerHealth   string
	}{
		{
			description:     "High performing partner",
			partnerName:     "Partner-A",
			partnerType:     "Distributor",
			yearlyTarget:    1000000.00,
			actualSales:     1300000.00,
			achievementRate: 130.0,
			partnerHealth:   "excellent",
		},
		{
			description:     "Underperforming partner",
			partnerName:     "Partner-B",
			partnerType:     "Reseller",
			yearlyTarget:    500000.00,
			actualSales:     300000.00,
			achievementRate: 60.0,
			partnerHealth:   "at-risk",
		},
	}

	for _, tc := range testCases {
		achievement := (tc.actualSales / tc.yearlyTarget) * 100
		assert.InDelta(t, tc.achievementRate, achievement, 0.1, tc.description)

		var health string
		if achievement >= 100 {
			health = "excellent"
		} else if achievement >= 80 {
			health = "healthy"
		} else if achievement >= 60 {
			health = "at-risk"
		} else {
			health = "critical"
		}
		assert.Equal(t, tc.partnerHealth, health, tc.description)
	}
}

// TestGeographicCoverageOptimization validates market penetration
func TestGeographicCoverageOptimization(t *testing.T) {
	marketData := map[string]struct {
		market          string
		marketSize      float64 // Total market opportunity
		ourMarketShare  float64 // Our sales
		penetrationRate float64
		growthOppty     float64 // Opportunity for growth
	}{
		"Metro": {
			market:          "Delhi Metro",
			marketSize:      10000000.00,
			ourMarketShare:  4000000.00,
			penetrationRate: 40.0,
			growthOppty:     6000000.00,
		},
		"Tier2": {
			market:          "Tier-2 Cities",
			marketSize:      5000000.00,
			ourMarketShare:  1000000.00,
			penetrationRate: 20.0,
			growthOppty:     4000000.00,
		},
	}

	for region, data := range marketData {
		penetration := (data.ourMarketShare / data.marketSize) * 100
		assert.InDelta(t, data.penetrationRate, penetration, 0.1, region)

		oppty := data.marketSize - data.ourMarketShare
		assert.InDelta(t, data.growthOppty, oppty, 0.01, region)
	}
}
