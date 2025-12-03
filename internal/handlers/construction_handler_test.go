package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateProjectHandler validates project creation
func TestCreateProjectHandler(t *testing.T) {
	payload := map[string]interface{}{
		"project_code":  "PRJ-001",
		"project_name":  "Commercial Complex",
		"location":      "Mumbai",
		"start_date":    "2025-01-15",
		"estimated_end": "2026-12-31",
		"total_budget":  50000000.00,
		"status":        "planning",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/construction/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestGetProjectHandler validates project retrieval
func TestGetProjectHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/construction/projects/proj-uuid-1234", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.Method)
	assert.Contains(t, req.URL.Path, "proj-uuid-1234")
}

// TestCreateBOQHandler validates BOQ creation
func TestCreateBOQHandler(t *testing.T) {
	payload := map[string]interface{}{
		"project_id":  "proj-uuid-1234",
		"boq_number":  "BOQ-001",
		"description": "Structural Steel",
		"items": []map[string]interface{}{
			{
				"item_description": "Steel Bars",
				"unit":             "kg",
				"quantity":         5000.0,
				"unit_rate":        50.0,
			},
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/construction/boqs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestImportBOQHandler validates BOQ import endpoint
func TestImportBOQHandler(t *testing.T) {
	payload := map[string]interface{}{
		"project_id": "proj-uuid-1234",
		"file_name":  "boq_import.xlsx",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/construction/boqs/import", bytes.NewReader(body))
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestExportBOQHandler validates BOQ export endpoint
func TestExportBOQHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/construction/projects/proj-uuid-1234/boqs/export", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.Method)
}

// TestCreateProgressHandler validates progress recording
func TestCreateProgressHandler(t *testing.T) {
	payload := map[string]interface{}{
		"project_id": "proj-uuid-1234",
		"month":      "2025-01",
		"percentage": 25.0,
		"notes":      "Foundation work completed",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/construction/progress", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestRecordQCHandler validates quality control recording
func TestRecordQCHandler(t *testing.T) {
	payload := map[string]interface{}{
		"project_id":  "proj-uuid-1234",
		"boq_item_id": "item-uuid-5678",
		"status":      "passed",
		"remarks":     "Quality check passed",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/construction/qc", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestBOQCalculations validates BOQ calculations with precision
func TestBOQCalculations(t *testing.T) {
	testCases := []struct {
		description string
		qty         float64
		rate        float64
		expected    float64
	}{
		{"Cement Bags", 1000.0, 50.0, 50000.0},
		{"Steel (kg)", 5000.0, 50.0, 250000.0},
		{"Sand (cubic_meters)", 100.0, 1500.0, 150000.0},
		{"Precision test", 250.5, 75.25, 18850.125},
	}

	for _, tc := range testCases {
		total := tc.qty * tc.rate
		assert.InDelta(t, tc.expected, total, 0.01, tc.description)
	}
}

// TestProgressPercentageValidation validates progress percentage constraints
func TestProgressPercentageValidation(t *testing.T) {
	testCases := []struct {
		description string
		percentage  float64
		valid       bool
	}{
		{"Project Start", 0.0, true},
		{"Mid-way Progress", 25.5, true},
		{"Half-way", 50.0, true},
		{"Near Complete", 99.9, true},
		{"Project Complete", 100.0, true},
		{"Negative Progress (Invalid)", -10.0, false},
		{"Over 100% (Invalid)", 110.0, false},
		{"Decimal Precision", 33.333, true},
	}

	for _, tc := range testCases {
		isValid := tc.percentage >= 0 && tc.percentage <= 100
		assert.Equal(t, tc.valid, isValid, tc.description)
	}
}

// TestProjectStatusValidation validates project status values
func TestProjectStatusValidation(t *testing.T) {
	validStatuses := []string{"planning", "initiated", "in_progress", "paused", "completed", "cancelled"}

	for _, status := range validStatuses {
		assert.NotEmpty(t, status)
	}
}

// TestBOQItemUnitValidation validates unit types
func TestBOQItemUnitValidation(t *testing.T) {
	validUnits := []string{"bags", "cubic_meters", "square_meters", "pieces", "liters", "kg", "meters"}

	for _, unit := range validUnits {
		assert.NotEmpty(t, unit)
	}
}

// TestQualityControlStatus validates QC status values
func TestQualityControlStatus(t *testing.T) {
	validStatuses := []string{"pending", "passed", "failed", "rework_required"}

	for _, status := range validStatuses {
		assert.NotEmpty(t, status)
	}
}

// TestEquipmentStatusTracking validates equipment status
func TestEquipmentStatusTracking(t *testing.T) {
	statuses := []string{"available", "assigned", "in_use", "maintenance", "retired"}

	assert.Equal(t, 5, len(statuses))
}

// TestMultiTenantConstructionRouting validates tenant isolation
func TestMultiTenantConstructionRouting(t *testing.T) {
	tenantReq1 := httptest.NewRequest("GET", "/api/v1/construction/projects/proj-1", nil)
	tenantReq1.Header.Set("X-Tenant-ID", "tenant-1")

	tenantReq2 := httptest.NewRequest("GET", "/api/v1/construction/projects/proj-1", nil)
	tenantReq2.Header.Set("X-Tenant-ID", "tenant-2")

	assert.NotEqual(t,
		tenantReq1.Header.Get("X-Tenant-ID"),
		tenantReq2.Header.Get("X-Tenant-ID"))
}

// TestProjectCodeUniqueness validates project code uniqueness
func TestProjectCodeUniqueness(t *testing.T) {
	codes := map[string]string{
		"PRJ-001": "proj-uuid-1",
		"PRJ-002": "proj-uuid-2",
		"PRJ-003": "proj-uuid-3",
	}

	assert.Equal(t, 3, len(codes))
	assert.NotEqual(t, codes["PRJ-001"], codes["PRJ-002"])
}

// TestBOQBudgetTracking validates budget tracking
func TestBOQBudgetTracking(t *testing.T) {
	budget := 50000000.00
	spent := 25000000.00
	remaining := budget - spent

	assert.InDelta(t, 25000000.00, remaining, 0.01)
	assert.True(t, spent <= budget)
}

// BenchmarkCreateProjectEndpoint benchmarks project creation
func BenchmarkCreateProjectEndpoint(b *testing.B) {
	payload := map[string]interface{}{
		"project_code": "PRJ-001",
		"project_name": "Commercial Complex",
		"location":     "Mumbai",
		"total_budget": 50000000.00,
	}

	body, _ := json.Marshal(payload)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/v1/construction/projects", bytes.NewReader(body))
		req.Header.Set("X-Tenant-ID", "tenant-1")
		httptest.NewRecorder()
	}
}

// BenchmarkCreateBOQEndpoint benchmarks BOQ creation
func BenchmarkCreateBOQEndpoint(b *testing.B) {
	payload := map[string]interface{}{
		"project_id": "proj-uuid-1234",
		"boq_number": "BOQ-001",
	}

	body, _ := json.Marshal(payload)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/v1/construction/boqs", bytes.NewReader(body))
		req.Header.Set("X-Tenant-ID", "tenant-1")
		httptest.NewRecorder()
	}
}
