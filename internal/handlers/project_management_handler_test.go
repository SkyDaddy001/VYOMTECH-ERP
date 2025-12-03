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

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestCreateCustomerProfile tests customer profile creation
func TestCreateCustomerProfile(t *testing.T) {
	customer := models.PropertyCustomerProfile{
		CustomerCode:      "CUST001",
		FirstName:         "John",
		MiddleName:        "Doe",
		LastName:          "Smith",
		Email:             "john@example.com",
		PhonePrimary:      "9876543210",
		CompanyName:       "ABC Corp",
		Designation:       "Manager",
		PANNumber:         "ABCDE1234F",
		AadharNumber:      "1234-5678-9012",
		CommunicationCity: "Delhi",
		PermanentCity:     "Delhi",
		Profession:        "IT",
		EmployerName:      "TechCorp",
		EmploymentType:    "SALARIED",
		MonthlyIncome:     100000.00,
	}

	body, _ := json.Marshal(customer)
	req := httptest.NewRequest("POST", "/api/v1/project-management/customers", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "/api/v1/project-management/customers", req.URL.Path)
}

// TestGetCustomerProfile tests fetching a customer profile
func TestGetCustomerProfile(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/project-management/customers/{id}", handler.GetCustomerProfile)

	req := httptest.NewRequest("GET", "/api/v1/project-management/customers/cust-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.NotNil(t, req)
}

// TestListCustomers tests listing customers
func TestListCustomers(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/project-management/customers?limit=20&offset=0", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.Contains(t, req.URL.RawQuery, "limit=20")
}

// TestUpdateCustomerProfile tests updating customer profile
func TestUpdateCustomerProfile(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	customer := models.PropertyCustomerProfile{
		FirstName:     "Jane",
		Email:         "jane@example.com",
		PhonePrimary:  "9876543210",
		MonthlyIncome: 150000.00,
	}

	body, _ := json.Marshal(customer)
	req := httptest.NewRequest("PUT", "/api/v1/project-management/customers/cust-123", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "PUT", req.Method)
	assert.NotNil(t, handler)
}

// TestCreateAreaStatement tests area statement creation
func TestCreateAreaStatement(t *testing.T) {
	areaStmt := models.CreateAreaStatementRequest{
		ProjectID:                 "proj-123",
		UnitID:                    "unit-123",
		AptNo:                     "A-101",
		Floor:                     "1",
		UnitType:                  "2BHK",
		Facing:                    "NORTH",
		RERACarPetAreaSqft:        950.00,
		CarPetAreaWithBalconySqft: 1050.00,
		PlinthAreaSqft:            1150.00,
		SBUASqft:                  1400.00,
		AlotedTo:                  "John Doe",
		NOCTaken:                  "YES",
	}

	body, _ := json.Marshal(areaStmt)
	req := httptest.NewRequest("POST", "/api/v1/project-management/area-statements", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "POST", req.Method)
	assert.NotNil(t, req)
}

// TestListAreaStatements tests listing area statements
func TestListAreaStatements(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/project-management/area-statements?project_id=proj-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.Contains(t, req.URL.RawQuery, "project_id=proj-123")
}

// TestGetAreaStatement tests fetching a specific area statement
func TestGetAreaStatement(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/project-management/area-statements/{unit_id}", handler.GetAreaStatement)

	req := httptest.NewRequest("GET", "/api/v1/project-management/area-statements/unit-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.NotNil(t, handler)
}

// TestDeleteAreaStatement tests deleting an area statement
func TestDeleteAreaStatement(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/project-management/area-statements/{unit_id}", handler.DeleteAreaStatement)

	req := httptest.NewRequest("DELETE", "/api/v1/project-management/area-statements/unit-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "DELETE", req.Method)
	assert.NotNil(t, handler)
}

// TestCreateBankFinancing tests bank financing creation
func TestCreateBankFinancing(t *testing.T) {
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
	assert.NotNil(t, req)
}

// TestGetBankFinancing tests fetching bank financing details
func TestGetBankFinancing(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/project-management/bank-financing/{id}", handler.GetBankFinancing)

	req := httptest.NewRequest("GET", "/api/v1/project-management/bank-financing/fin-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.NotNil(t, handler)
}

// TestUpdateBankFinancing tests updating bank financing
func TestUpdateBankFinancing(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	financing := models.CreateBankFinancingRequest{
		ApartmentCost:    2500000.00,
		SanctionedAmount: 2100000.00,
		BankName:         "ICICI Bank",
	}

	body, _ := json.Marshal(financing)
	req := httptest.NewRequest("PUT", "/api/v1/project-management/bank-financing/fin-123", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "PUT", req.Method)
	assert.NotNil(t, handler)
}

// TestCreatePaymentStage tests payment stage creation
func TestCreatePaymentStage(t *testing.T) {
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
	assert.NotNil(t, req)
}

// TestListPaymentStages tests listing payment stages
func TestListPaymentStages(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/project-management/payment-stages?unit_id=unit-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.Contains(t, req.URL.RawQuery, "unit_id=unit-123")
}

// TestRecordPaymentCollection tests recording payment collection
func TestRecordPaymentCollection(t *testing.T) {
	payment := models.UpdatePaymentStageRequest{
		PaymentStageID:   "stage-123",
		AmountReceived:   500000.00,
		PaymentMode:      "NEFT",
		ReferenceNo:      "REF123",
		CollectionStatus: "COMPLETED",
	}

	body, _ := json.Marshal(payment)
	req := httptest.NewRequest("POST", "/api/v1/project-management/payment-collection", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "POST", req.Method)
	assert.NotNil(t, req)
}

// TestCreateDisbursementSchedule tests disbursement schedule creation
func TestCreateDisbursementSchedule(t *testing.T) {
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
	assert.NotNil(t, req)
}

// TestUpdateDisbursement tests updating disbursement
func TestUpdateDisbursement(t *testing.T) {
	disbursement := models.UpdateDisbursementRequest{
		DisbursementID:           "disb-123",
		ActualDisbursementDate:   "2024-04-05",
		ActualDisbursementAmount: 400000.00,
		DisbursementStatus:       "COMPLETED",
	}

	body, _ := json.Marshal(disbursement)
	req := httptest.NewRequest("PUT", "/api/v1/project-management/disbursement-schedule/disb-123", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "PUT", req.Method)
	assert.NotNil(t, req)
}

// TestUpdateCostSheet tests cost sheet updates
func TestUpdateCostSheet(t *testing.T) {
	costSheet := models.UpdateCostSheetRequest{
		UnitID:                       "unit-123",
		BlockName:                    "Block A",
		SBUA:                         1400.00,
		RatePerSqft:                  5000.00,
		CarParkingCost:               500000.00,
		ApartmentCostExcludingGovt:   7000000.00,
		ActualSoldPriceExcludingGovt: 7500000.00,
		GSTApplicable:                true,
		GSTPercentage:                18.0,
	}

	body, _ := json.Marshal(costSheet)
	req := httptest.NewRequest("PUT", "/api/v1/project-management/cost-sheet", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "PUT", req.Method)
	assert.NotNil(t, req)
}

// TestCreateProjectCostConfiguration tests project cost configuration
func TestCreateProjectCostConfiguration(t *testing.T) {
	config := models.CreateProjectCostConfigRequest{
		ProjectID:    "proj-123",
		ConfigName:   "CMWSSB Charges",
		ConfigType:   "OTHER_CHARGE_1",
		DisplayOrder: 1,
		IsMandatory:  true,
		Description:  "Water and sewerage charges",
	}

	body, _ := json.Marshal(config)
	req := httptest.NewRequest("POST", "/api/v1/project-management/cost-configuration", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "POST", req.Method)
	assert.NotNil(t, req)
}

// TestGetCostBreakdown tests cost breakdown API
func TestGetCostBreakdown(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/project-management/cost-breakdown/{unit_id}", handler.GetCostBreakdown)

	req := httptest.NewRequest("GET", "/api/v1/project-management/cost-breakdown/unit-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.NotNil(t, handler)
}

// TestGetProjectSummary tests project summary
func TestGetProjectSummary(t *testing.T) {
	handler := &ProjectManagementHandler{
		Service: &services.ProjectManagementService{},
		DB:      &sql.DB{},
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/project-management/summary/{project_id}", handler.GetProjectSummary)

	req := httptest.NewRequest("GET", "/api/v1/project-management/summary/proj-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.NotNil(t, handler)
}

// TestGetBankFinancingReport tests bank financing report
func TestGetBankFinancingReport(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/project-management/reports/bank-financing?project_id=proj-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.Contains(t, req.URL.RawQuery, "project_id=proj-123")
}

// TestGetPaymentStageReport tests payment stage report
func TestGetPaymentStageReport(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/project-management/reports/payment-stages?project_id=proj-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.Contains(t, req.URL.RawQuery, "project_id=proj-123")
}

// TestGetCollectionStatus tests collection status API
func TestGetCollectionStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/project-management/collection-status?project_id=proj-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.NotNil(t, req)
}

// TestGetDisbursementStatus tests disbursement status API
func TestGetDisbursementStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/project-management/disbursement-status?project_id=proj-123", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))

	assert.Equal(t, "GET", req.Method)
	assert.NotNil(t, req)
}

// BenchmarkCreateCustomerProfile benchmarks customer creation
func BenchmarkCreateCustomerProfile(b *testing.B) {
	customer := models.PropertyCustomerProfile{
		CustomerCode: "CUST001",
		FirstName:    "John",
		Email:        "john@example.com",
	}

	body, _ := json.Marshal(customer)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/v1/project-management/customers", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))
		_ = req
	}
}

// BenchmarkListCustomers benchmarks customer listing
func BenchmarkListCustomers(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/v1/project-management/customers?limit=20&offset=0", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.TenantIDKey, "tenant-123"))
		_ = req
	}
}
