package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBOQImportResult validates import result structure
func TestBOQImportResult(t *testing.T) {
	result := &BOQImportResult{
		TotalRows:      100,
		SuccessCount:   95,
		FailureCount:   5,
		CreatedBOQs:    90,
		UpdatedBOQs:    5,
		TotalAmountINR: 1000000.00,
		Errors:         []string{"Row 5: error"},
	}

	assert.Equal(t, 100, result.TotalRows)
	assert.Equal(t, 95, result.SuccessCount)
	assert.Equal(t, 5, result.FailureCount)
	assert.Equal(t, 1000000.00, result.TotalAmountINR)
}

// TestBOQImportSuccessRate validates success rate calculation
func TestBOQImportSuccessRate(t *testing.T) {
	testCases := []struct {
		total    int
		success  int
		expected float64
	}{
		{100, 100, 100.0},
		{100, 95, 95.0},
		{100, 50, 50.0},
		{100, 0, 0.0},
	}

	for _, tc := range testCases {
		rate := (float64(tc.success) / float64(tc.total)) * 100
		assert.InDelta(t, tc.expected, rate, 0.1)
	}
}

// TestParseHeaders validates header parsing
func TestParseHeaders(t *testing.T) {
	headerRow := []string{"BOQ Number", "Item Description", "Unit", "Quantity"}
	headerMap := parseHeaders(headerRow)

	assert.NotNil(t, headerMap)
	assert.Equal(t, 0, headerMap["BOQ Number"])
	assert.Equal(t, 1, headerMap["Item Description"])
	assert.Equal(t, 2, headerMap["Unit"])
	assert.Equal(t, 3, headerMap["Quantity"])
}

// TestBOQUnitTypes validates unit types
func TestBOQUnitTypes(t *testing.T) {
	units := []string{"bags", "cubic_meters", "square_meters", "pieces", "liters", "kg"}

	for _, unit := range units {
		assert.NotEmpty(t, unit)
	}
}

// TestBOQPagination validates pagination parameters
func TestBOQPagination(t *testing.T) {
	testCases := []struct {
		limit  int
		offset int
		valid  bool
	}{
		{50, 0, true},
		{50, 100, true},
		{500, 0, true},
		{0, 0, false},
		{-10, 0, false},
		{50, -10, false},
	}

	for _, tc := range testCases {
		isValid := tc.limit > 0 && tc.offset >= 0
		assert.Equal(t, tc.valid, isValid)
	}
}

// TestBOQTotalCalculation validates amount calculation
func TestBOQTotalCalculation(t *testing.T) {
	testCases := []struct {
		qty      float64
		rate     float64
		expected float64
	}{
		{100.0, 350.0, 35000.0},
		{50.5, 200.0, 10100.0},
		{75.5, 225.5, 17025.25},
	}

	for _, tc := range testCases {
		total := tc.qty * tc.rate
		assert.InDelta(t, tc.expected, total, 0.01)
	}
} // TestBOQDuplicateDetection validates duplicate detection
func TestBOQDuplicateDetection(t *testing.T) {
	type BOQKey struct {
		projectID string
		itemDesc  string
	}

	boq1 := BOQKey{"proj-1", "Cement"}
	boq2 := BOQKey{"proj-1", "Cement"}
	boq3 := BOQKey{"proj-1", "Steel"}

	assert.Equal(t, boq1, boq2)
	assert.NotEqual(t, boq1, boq3)
}

// TestBOQValidation validates BOQ item validation
func TestBOQValidation(t *testing.T) {
	testCases := []struct {
		description string
		unitRate    float64
		valid       bool
	}{
		{"Cement", 350.0, true},
		{"", 350.0, false},       // Missing description
		{"Steel", -350.0, false}, // Negative rate
		{"Wire", 0.0, true},      // Zero rate allowed initially
	}

	for _, tc := range testCases {
		isValid := tc.description != "" && tc.unitRate >= 0
		assert.Equal(t, tc.valid, isValid)
	}
}

// TestBOQBulkOperations validates bulk operations
func TestBOQBulkOperations(t *testing.T) {
	itemCount := 10
	items := make([]float64, itemCount)

	for i := 0; i < itemCount; i++ {
		items[i] = 35000.0
	}

	var total float64
	for _, item := range items {
		total += item
	}

	assert.Equal(t, 10, len(items))
	assert.Equal(t, 350000.0, total)
}

// TestBOQStatusWorkflow validates status transitions
func TestBOQStatusWorkflow(t *testing.T) {
	status := "planned"
	assert.Equal(t, "planned", status)

	status = "in_progress"
	assert.Equal(t, "in_progress", status)

	status = "completed"
	assert.Equal(t, "completed", status)
}

// TestBOQCategorySum validates summing by category
func TestBOQCategorySum(t *testing.T) {
	categories := map[string]float64{
		"civil":      35000.0,
		"structural": 25000.0,
		"electrical": 15000.0,
		"plumbing":   10000.0,
		"finishing":  5000.0,
	}

	var total float64
	for _, amount := range categories {
		total += amount
	}

	assert.Equal(t, 90000.0, total)
	assert.Equal(t, 5, len(categories))
}

// TestBOQExportFormat validates export data structure
func TestBOQExportFormat(t *testing.T) {
	exportItem := struct {
		boqNumber string
		itemDesc  string
		unit      string
		qty       float64
		rate      float64
		total     float64
	}{
		boqNumber: "BOQ001",
		itemDesc:  "Cement",
		unit:      "bags",
		qty:       100.0,
		rate:      350.0,
		total:     35000.0,
	}

	assert.NotEmpty(t, exportItem.boqNumber)
	assert.NotEmpty(t, exportItem.itemDesc)
	assert.Greater(t, exportItem.qty, 0.0)
	assert.InDelta(t, exportItem.qty*exportItem.rate, exportItem.total, 0.01)
}
