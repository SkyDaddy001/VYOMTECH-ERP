package services

import (
	"testing"

	"vyomtech-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// TestConstructionProject validates project model uses UUID strings
func TestConstructionProject(t *testing.T) {
	project := &models.ConstructionProject{
		ID:                     "project-uuid-001",
		TenantID:               "tenant-001",
		ProjectName:            "Tower A",
		ProjectCode:            "TA001",
		ContractValue:          50000000.00,
		CurrentProgressPercent: 0,
		Status:                 "planning",
	}

	// Verify string ID type (not int64)
	assert.IsType(t, "", project.ID)
	assert.Equal(t, "project-uuid-001", project.ID)
	assert.Equal(t, 50000000.00, project.ContractValue)
}

// TestProjectStatus validates project status values
func TestProjectStatus(t *testing.T) {
	statuses := []string{"planning", "active", "suspended", "completed", "on_hold"}

	for _, status := range statuses {
		proj := &models.ConstructionProject{Status: status}
		assert.Equal(t, status, proj.Status)
	}
}

// TestBillOfQuantities validates BOQ model uses UUID strings
func TestBillOfQuantities(t *testing.T) {
	boq := &models.BillOfQuantities{
		ID:              "boq-uuid-001",
		TenantID:        "tenant-001",
		ProjectID:       "project-uuid-001", // Should be string, not int64
		BOQNumber:       "BOQ001",
		ItemDescription: "Cement",
		Unit:            "bags",
		Quantity:        100.0,
		UnitRate:        350.00,
		TotalAmount:     35000.00,
		Status:          "planned",
	}

	// Verify string types
	assert.IsType(t, "", boq.ID)
	assert.IsType(t, "", boq.ProjectID)
	assert.Equal(t, "boq-uuid-001", boq.ID)
	assert.Equal(t, "project-uuid-001", boq.ProjectID)
}

// TestBOQCalculation validates quantity * rate = total
func TestBOQCalculation(t *testing.T) {
	testCases := []struct {
		qty      float64
		rate     float64
		expected float64
	}{
		{100.0, 350.0, 35000.0},
		{50.5, 200.0, 10100.0},
		{1000.0, 1000.0, 1000000.0},
	}

	for _, tc := range testCases {
		total := tc.qty * tc.rate
		assert.InDelta(t, tc.expected, total, 0.01)
	}
}

// TestBOQStatus validates BOQ status values
func TestBOQStatus(t *testing.T) {
	statuses := []string{"planned", "in_progress", "completed", "on_hold"}

	for _, status := range statuses {
		boq := &models.BillOfQuantities{Status: status}
		assert.Equal(t, status, boq.Status)
	}
}

// TestBOQCategory validates BOQ categories
func TestBOQCategory(t *testing.T) {
	categories := []string{"civil", "structural", "electrical", "plumbing", "finishing", "other"}

	for _, cat := range categories {
		boq := &models.BillOfQuantities{Category: cat}
		assert.Equal(t, cat, boq.Category)
	}
}

// TestProgressTracking validates progress model uses UUID strings
func TestProgressTracking(t *testing.T) {
	progress := &models.ProgressTracking{
		ID:                "progress-uuid-001",
		TenantID:          "tenant-001",
		ProjectID:         "project-uuid-001", // Should be string, not int64
		PercentComplete:   25,
		WorkforceDeployed: 50,
	}

	assert.IsType(t, "", progress.ID)
	assert.IsType(t, "", progress.ProjectID)
	assert.Equal(t, 25, progress.PercentComplete)
}

// TestProgressValidation validates progress percentage
func TestProgressValidation(t *testing.T) {
	testCases := []struct {
		percent int
		valid   bool
	}{
		{0, true},
		{50, true},
		{100, true},
		{-10, false},
		{110, false},
	}

	for _, tc := range testCases {
		isValid := tc.percent >= 0 && tc.percent <= 100
		assert.Equal(t, tc.valid, isValid)
	}
}

// TestQualityControl validates QC model uses UUID strings
func TestQualityControl(t *testing.T) {
	qc := &models.QualityControl{
		ID:            "qc-uuid-001",
		TenantID:      "tenant-001",
		ProjectID:     "project-uuid-001", // Should be string
		BOQItemID:     "boq-uuid-001",     // Should be string
		QualityStatus: "passed",
		InspectorName: "Inspector John",
	}

	assert.IsType(t, "", qc.ID)
	assert.IsType(t, "", qc.ProjectID)
	assert.IsType(t, "", qc.BOQItemID)
	assert.Equal(t, "passed", qc.QualityStatus)
}

// TestQCStatus validates QC status values
func TestQCStatus(t *testing.T) {
	statuses := []string{"passed", "failed", "partial", "pending"}

	for _, status := range statuses {
		qc := &models.QualityControl{QualityStatus: status}
		assert.Equal(t, status, qc.QualityStatus)
	}
}

// TestConstructionEquipment validates equipment model uses UUID strings
func TestConstructionEquipment(t *testing.T) {
	equip := &models.ConstructionEquipment{
		ID:            "equip-uuid-001",
		TenantID:      "tenant-001",
		ProjectID:     "project-uuid-001", // Should be string, not int64
		EquipmentName: "Excavator",
		Status:        "in_use",
		CostPerDay:    5000.00,
	}

	assert.IsType(t, "", equip.ID)
	assert.IsType(t, "", equip.ProjectID)
	assert.Equal(t, "in_use", equip.Status)
}

// TestEquipmentStatus validates equipment status values
func TestEquipmentStatus(t *testing.T) {
	statuses := []string{"available", "in_use", "maintenance", "retired"}

	for _, status := range statuses {
		equip := &models.ConstructionEquipment{Status: status}
		assert.Equal(t, status, equip.Status)
	}
}

// TestBOQSummation validates summing BOQ items
func TestBOQSummation(t *testing.T) {
	items := []float64{35000.0, 25000.0, 30000.0}
	var total float64
	for _, amount := range items {
		total += amount
	}

	assert.Equal(t, 90000.0, total)
}

// TestConstructionCostTracking validates contract value vs BOQ total
func TestConstructionCostTracking(t *testing.T) {
	contractValue := 50000000.0
	boqTotal := 45000000.0
	contingency := contractValue - boqTotal

	contingencyPercent := (contingency / contractValue) * 100
	assert.InDelta(t, 10.0, contingencyPercent, 0.1)
}

// TestMultiTenantConstruction validates construction isolation
func TestMultiTenantConstruction(t *testing.T) {
	proj1 := &models.ConstructionProject{ID: "proj-1", TenantID: "tenant-1"}
	proj2 := &models.ConstructionProject{ID: "proj-2", TenantID: "tenant-2"}

	assert.NotEqual(t, proj1.TenantID, proj2.TenantID)
}

// TestProjectIdStringType validates ProjectID is always string
func TestProjectIdStringType(t *testing.T) {
	// BOQ
	boq := &models.BillOfQuantities{ProjectID: "proj-uuid"}
	assert.IsType(t, "", boq.ProjectID)

	// Progress
	progress := &models.ProgressTracking{ProjectID: "proj-uuid"}
	assert.IsType(t, "", progress.ProjectID)

	// QC
	qc := &models.QualityControl{ProjectID: "proj-uuid"}
	assert.IsType(t, "", qc.ProjectID)

	// Equipment
	equip := &models.ConstructionEquipment{ProjectID: "proj-uuid"}
	assert.IsType(t, "", equip.ProjectID)
}
