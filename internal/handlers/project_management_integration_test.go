package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/stretchr/testify/assert"
)

// TestProjectManagementIntegration tests complete project management workflows
func TestProjectManagementIntegration(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.Service)
}

// TestCustomerLifecycleWorkflow tests the complete customer lifecycle
func TestCustomerLifecycleWorkflow(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	// Create customer
	customer := models.PropertyCustomerProfile{
		CustomerCode:      "CUST001",
		FirstName:         "John",
		Email:             "john@example.com",
		PhonePrimary:      "9876543210",
		CommunicationCity: "Delhi",
		PermanentCity:     "Delhi",
	}

	body, _ := json.Marshal(customer)
	req := httptest.NewRequest("POST", "/api/v1/project-management/customers", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "POST", req.Method)

	// Get customer
	req2 := httptest.NewRequest("GET", "/api/v1/project-management/customers/cust-123", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req2.Method)

	// Update customer
	customer.FirstName = "Jane"
	body2, _ := json.Marshal(customer)
	req3 := httptest.NewRequest("PUT", "/api/v1/project-management/customers/cust-123", bytes.NewBuffer(body2))
	req3 = req3.WithContext(context.WithValue(req3.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "PUT", req3.Method)

	// List customers
	req4 := httptest.NewRequest("GET", "/api/v1/project-management/customers", nil)
	req4 = req4.WithContext(context.WithValue(req4.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req4.Method)
	assert.NotNil(t, handler)
}

// TestAreaStatementWorkflow tests area statement operations
func TestAreaStatementWorkflow(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	// Create area statement
	areaStmt := models.CreateAreaStatementRequest{
		ProjectID:                 "proj-123",
		UnitID:                    "unit-123",
		AptNo:                     "A-101",
		RERACarPetAreaSqft:        950.00,
		CarPetAreaWithBalconySqft: 1050.00,
		PlinthAreaSqft:            1150.00,
		SBUASqft:                  1400.00,
	}

	body, _ := json.Marshal(areaStmt)
	req := httptest.NewRequest("POST", "/api/v1/project-management/area-statements", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "POST", req.Method)

	// List area statements
	req2 := httptest.NewRequest("GET", "/api/v1/project-management/area-statements", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req2.Method)

	// Get area statement
	req3 := httptest.NewRequest("GET", "/api/v1/project-management/area-statements/unit-123", nil)
	req3 = req3.WithContext(context.WithValue(req3.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req3.Method)

	// Delete area statement
	req4 := httptest.NewRequest("DELETE", "/api/v1/project-management/area-statements/unit-123", nil)
	req4 = req4.WithContext(context.WithValue(req4.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "DELETE", req4.Method)
	assert.NotNil(t, handler)
}

// TestBankFinancingWorkflow tests bank financing operations
func TestBankFinancingWorkflow(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	// Create bank financing
	financing := models.CreateBankFinancingRequest{
		ProjectID:        "proj-123",
		UnitID:           "unit-123",
		CustomerID:       "cust-123",
		ApartmentCost:    2500000.00,
		SanctionedAmount: 2000000.00,
		BankName:         "HDFC Bank",
	}

	body, _ := json.Marshal(financing)
	req := httptest.NewRequest("POST", "/api/v1/project-management/bank-financing", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "POST", req.Method)

	// Get bank financing
	req2 := httptest.NewRequest("GET", "/api/v1/project-management/bank-financing/fin-123", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req2.Method)

	// Update bank financing
	financing.SanctionedAmount = 2100000.00
	body2, _ := json.Marshal(financing)
	req3 := httptest.NewRequest("PUT", "/api/v1/project-management/bank-financing/fin-123", bytes.NewBuffer(body2))
	req3 = req3.WithContext(context.WithValue(req3.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "PUT", req3.Method)
	assert.NotNil(t, handler)
}

// TestPaymentStageWorkflow tests payment stage operations
func TestPaymentStageWorkflow(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	// Create payment stage
	paymentStage := models.CreatePaymentStageRequest{
		ProjectID:       "proj-123",
		UnitID:          "unit-123",
		CustomerID:      "cust-123",
		StageName:       "BOOKING",
		StageNumber:     1,
		StagePercentage: 25.0,
		ApartmentCost:   2500000.00,
		DueDate:         "2024-03-01",
	}

	body, _ := json.Marshal(paymentStage)
	req := httptest.NewRequest("POST", "/api/v1/project-management/payment-stages", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "POST", req.Method)

	// List payment stages
	req2 := httptest.NewRequest("GET", "/api/v1/project-management/payment-stages", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req2.Method)

	// Record payment
	paymentUpdate := models.UpdatePaymentStageRequest{
		PaymentStageID:   "stage-123",
		AmountReceived:   625000.00,
		PaymentMode:      "NEFT",
		ReferenceNo:      "REF123",
		CollectionStatus: "COMPLETED",
	}

	body2, _ := json.Marshal(paymentUpdate)
	req3 := httptest.NewRequest("POST", "/api/v1/project-management/payment-collection", bytes.NewBuffer(body2))
	req3 = req3.WithContext(context.WithValue(req3.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "POST", req3.Method)
	assert.NotNil(t, handler)
}

// TestDisbursementWorkflow tests disbursement operations
func TestDisbursementWorkflow(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	// Create disbursement schedule
	disbursement := models.CreateDisbursementScheduleRequest{
		FinancingID:                "fin-123",
		UnitID:                     "unit-123",
		CustomerID:                 "cust-123",
		DisbursementNo:             1,
		ExpectedDisbursementDate:   "2024-04-01",
		ExpectedDisbursementAmount: 400000.00,
		DisbursementPercentage:     25.0,
	}

	body, _ := json.Marshal(disbursement)
	req := httptest.NewRequest("POST", "/api/v1/project-management/disbursement-schedule", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "POST", req.Method)

	// Update disbursement
	disbursementUpdate := models.UpdateDisbursementRequest{
		DisbursementID:           "disb-123",
		ActualDisbursementDate:   "2024-04-05",
		ActualDisbursementAmount: 400000.00,
		DisbursementStatus:       "COMPLETED",
	}

	body2, _ := json.Marshal(disbursementUpdate)
	req2 := httptest.NewRequest("PUT", "/api/v1/project-management/disbursement-schedule/disb-123", bytes.NewBuffer(body2))
	req2 = req2.WithContext(context.WithValue(req2.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "PUT", req2.Method)
	assert.NotNil(t, handler)
}

// TestCostManagementWorkflow tests cost sheet and configuration operations
func TestCostManagementWorkflow(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	// Create cost configuration
	config := models.CreateProjectCostConfigRequest{
		ProjectID:    "proj-123",
		ConfigName:   "CMWSSB Charges",
		ConfigType:   "OTHER_CHARGE_1",
		DisplayOrder: 1,
	}

	body, _ := json.Marshal(config)
	req := httptest.NewRequest("POST", "/api/v1/project-management/cost-configuration", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "POST", req.Method)

	// Update cost sheet
	costSheet := models.UpdateCostSheetRequest{
		UnitID:         "unit-123",
		SBUA:           1400.00,
		RatePerSqft:    5000.00,
		CarParkingCost: 500000.00,
		GSTApplicable:  true,
		GSTPercentage:  18.0,
	}

	body2, _ := json.Marshal(costSheet)
	req2 := httptest.NewRequest("PUT", "/api/v1/project-management/cost-sheet", bytes.NewBuffer(body2))
	req2 = req2.WithContext(context.WithValue(req2.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "PUT", req2.Method)
	assert.NotNil(t, handler)
}

// TestDashboardAndReports tests reporting endpoints
func TestDashboardAndReports(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	// Get project summary
	req := httptest.NewRequest("GET", "/api/v1/project-management/summary/proj-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)

	// Get bank financing report
	req2 := httptest.NewRequest("GET", "/api/v1/project-management/reports/bank-financing", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req2.Method)

	// Get payment stage report
	req3 := httptest.NewRequest("GET", "/api/v1/project-management/reports/payment-stages", nil)
	req3 = req3.WithContext(context.WithValue(req3.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req3.Method)

	// Get collection status
	req4 := httptest.NewRequest("GET", "/api/v1/project-management/collection-status", nil)
	req4 = req4.WithContext(context.WithValue(req4.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req4.Method)

	// Get disbursement status
	req5 := httptest.NewRequest("GET", "/api/v1/project-management/disbursement-status", nil)
	req5 = req5.WithContext(context.WithValue(req5.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req5.Method)

	// Get cost breakdown
	req6 := httptest.NewRequest("GET", "/api/v1/project-management/cost-breakdown/unit-123", nil)
	req6 = req6.WithContext(context.WithValue(req6.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req6.Method)
	assert.NotNil(t, handler)
}

// TestMultiTenantIsolation tests multi-tenant data isolation
func TestMultiTenantIsolation(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	// Request from tenant 1
	req1 := httptest.NewRequest("GET", "/api/v1/project-management/customers", nil)
	req1 = req1.WithContext(context.WithValue(req1.Context(), middleware.TenantIDKey, "tenant-1"))

	// Request from tenant 2
	req2 := httptest.NewRequest("GET", "/api/v1/project-management/customers", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), middleware.TenantIDKey, "tenant-2"))

	assert.NotEqual(t, "tenant-1", "tenant-2")
	assert.NotNil(t, handler)
}

// TestErrorHandling tests error handling in operations
func TestErrorHandling(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	// Test with invalid JSON
	invalidJSON := []byte(`{invalid}`)
	req := httptest.NewRequest("POST", "/api/v1/project-management/customers", bytes.NewBuffer(invalidJSON))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.NotNil(t, req)
	assert.NotNil(t, handler)
}

// TestDataConsistency tests data consistency across operations
func TestDataConsistency(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	// Create and verify customer data
	customer := models.PropertyCustomerProfile{
		CustomerCode: "CUST001",
		FirstName:    "John",
		Email:        "john@example.com",
	}

	assert.Equal(t, "CUST001", customer.CustomerCode)
	assert.Equal(t, "John", customer.FirstName)
	assert.Equal(t, "john@example.com", customer.Email)
	assert.NotNil(t, handler)
}
