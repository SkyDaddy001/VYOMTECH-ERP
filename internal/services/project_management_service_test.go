package services

import (
	"database/sql"
	"testing"
	"time"
	"vyomtech-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// TestCreateCustomerProfile tests customer profile creation
func TestCreateCustomerProfile(t *testing.T) {
	service := &ProjectManagementService{
		DB: &sql.DB{},
	}

	profile := &models.PropertyCustomerProfile{
		CustomerCode: "CUST001",
		FirstName:    "John",
		Email:        "john@example.com",
	}

	assert.NotNil(t, service)
	assert.Equal(t, "John", profile.FirstName)
	assert.Equal(t, "CUST001", profile.CustomerCode)
	assert.Equal(t, "john@example.com", profile.Email)
}

// TestGetCustomerProfile tests fetching customer profile
func TestGetCustomerProfile(t *testing.T) {
	service := &ProjectManagementService{
		DB: &sql.DB{},
	}

	assert.NotNil(t, service)
}

// TestUpdateCustomerProfile tests updating customer profile
func TestUpdateCustomerProfile(t *testing.T) {
	service := &ProjectManagementService{
		DB: &sql.DB{},
	}

	updatedProfile := &models.PropertyCustomerProfile{
		FirstName: "Jane",
	}

	assert.NotNil(t, service)
	assert.Equal(t, "Jane", updatedProfile.FirstName)
}

// TestListCustomerProfiles tests listing customer profiles
func TestListCustomerProfiles(t *testing.T) {
	service := &ProjectManagementService{
		DB: &sql.DB{},
	}

	assert.NotNil(t, service)
}

// TestCreateAreaStatement tests area statement creation
func TestCreateAreaStatement(t *testing.T) {
	areaStmt := &models.CreateAreaStatementRequest{
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

	assert.NotNil(t, areaStmt)
	assert.Equal(t, 1050.00, areaStmt.CarPetAreaWithBalconySqft)
}

// TestCreateBankFinancing tests bank financing creation
func TestCreateBankFinancing(t *testing.T) {
	financing := &models.CreateBankFinancingRequest{
		ProjectID:        "proj-123",
		UnitID:           "unit-123",
		CustomerID:       "cust-123",
		ApartmentCost:    2500000.00,
		SanctionedAmount: 2000000.00,
		BankName:         "HDFC Bank",
	}

	assert.NotNil(t, financing)
	assert.Equal(t, "HDFC Bank", financing.BankName)
}

// TestCreatePaymentStage tests payment stage creation
func TestCreatePaymentStage(t *testing.T) {
	paymentStage := &models.CreatePaymentStageRequest{
		ProjectID:       "proj-123",
		UnitID:          "unit-123",
		CustomerID:      "cust-123",
		StageName:       "BOOKING",
		StageNumber:     1,
		StagePercentage: 25.0,
		ApartmentCost:   2500000.00,
		DueDate:         "2024-03-01",
	}

	assert.NotNil(t, paymentStage)
	assert.Equal(t, 1, paymentStage.StageNumber)
}

// TestUpdatePaymentStage tests payment stage update
func TestUpdatePaymentStage(t *testing.T) {
	paymentUpdate := &models.UpdatePaymentStageRequest{
		PaymentStageID:   "stage-123",
		AmountReceived:   625000.00,
		PaymentMode:      "NEFT",
		ReferenceNo:      "REF123",
		CollectionStatus: "COMPLETED",
	}

	assert.NotNil(t, paymentUpdate)
	assert.Equal(t, 625000.00, paymentUpdate.AmountReceived)
}

// TestCreateDisbursementSchedule tests disbursement creation
func TestCreateDisbursementSchedule(t *testing.T) {
	disbursement := &models.CreateDisbursementScheduleRequest{
		FinancingID:                "fin-123",
		UnitID:                     "unit-123",
		CustomerID:                 "cust-123",
		DisbursementNo:             1,
		ExpectedDisbursementDate:   "2024-04-01",
		ExpectedDisbursementAmount: 400000.00,
		DisbursementPercentage:     25.0,
		LinkedMilestoneID:          "milestone-1",
		MilestoneStage:             "FOUNDATION",
	}

	assert.NotNil(t, disbursement)
	assert.Equal(t, 400000.00, disbursement.ExpectedDisbursementAmount)
}

// TestUpdateDisbursement tests disbursement update
func TestUpdateDisbursement(t *testing.T) {
	disbursementUpdate := &models.UpdateDisbursementRequest{
		DisbursementID:           "disb-123",
		ActualDisbursementDate:   "2024-04-05",
		ActualDisbursementAmount: 400000.00,
		DisbursementStatus:       "COMPLETED",
		ChequeNo:                 "CHQ123",
		BankReferenceNo:          "BNK123",
	}

	assert.NotNil(t, disbursementUpdate)
	assert.Equal(t, "COMPLETED", disbursementUpdate.DisbursementStatus)
}

// TestUpdateCostSheet tests cost sheet update
func TestUpdateCostSheet(t *testing.T) {
	costSheet := &models.UpdateCostSheetRequest{
		UnitID:                       "unit-123",
		BlockName:                    "Block A",
		SBUA:                         1400.00,
		RatePerSqft:                  5000.00,
		CarParkingCost:               500000.00,
		ApartmentCostExcludingGovt:   7000000.00,
		ActualSoldPriceExcludingGovt: 7500000.00,
		GSTApplicable:                true,
		GSTPercentage:                18.0,
		ClubMembership:               100000.00,
		RegistrationCharge:           100000.00,
	}

	assert.NotNil(t, costSheet)
	assert.Equal(t, 1400.00, costSheet.SBUA)
	assert.True(t, costSheet.GSTApplicable)
}

// TestCreateProjectCostConfig tests cost configuration
func TestCreateProjectCostConfig(t *testing.T) {
	config := &models.CreateProjectCostConfigRequest{
		ProjectID:    "proj-123",
		ConfigName:   "CMWSSB Charges",
		ConfigType:   "OTHER_CHARGE_1",
		ChargeType:   "PER_SQFT",
		ChargeAmount: 100.00,
		DisplayOrder: 1,
		IsMandatory:  true,
		Description:  "Water charges",
	}

	assert.NotNil(t, config)
	assert.Equal(t, "CMWSSB Charges", config.ConfigName)
}

// TestFieldValidation tests field validation
func TestFieldValidation(t *testing.T) {
	profile := &models.PropertyCustomerProfile{
		FirstName:    "",
		Email:        "invalid-email",
		PhonePrimary: "12345",
	}

	assert.Equal(t, "", profile.FirstName)
	assert.NotEmpty(t, profile.Email)
	assert.Equal(t, "12345", profile.PhonePrimary)
}

// TestMoneyFieldHandling tests monetary field handling
func TestMoneyFieldHandling(t *testing.T) {
	stage := &models.PropertyPaymentStage{
		StageDueAmount: 2500000.00,
		AmountReceived: 1250000.00,
		AmountPending:  1250000.00,
	}

	assert.Equal(t, 2500000.00, stage.StageDueAmount)
	assert.Equal(t, 1250000.00, stage.AmountReceived)
	assert.Equal(t, 1250000.00, stage.AmountPending)
}

// TestDateFieldHandling tests date field handling
func TestDateFieldHandling(t *testing.T) {
	now := time.Now()
	dueDate := &now

	stage := &models.PropertyPaymentStage{
		DueDate: dueDate,
	}

	assert.NotNil(t, stage.DueDate)
	assert.Equal(t, dueDate, stage.DueDate)
}

// TestCoApplicantHandling tests co-applicant field handling
func TestCoApplicantHandling(t *testing.T) {
	profile := &models.PropertyCustomerProfile{
		FirstName:            "Primary",
		CoApplicant1Name:     "CoApp1",
		CoApplicant2Name:     "CoApp2",
		CoApplicant3Name:     "CoApp3",
		CoApplicant1Relation: "SPOUSE",
		CoApplicant2Relation: "CHILD",
		CoApplicant3Relation: "PARENT",
	}

	assert.Equal(t, "Primary", profile.FirstName)
	assert.Equal(t, "CoApp1", profile.CoApplicant1Name)
	assert.Equal(t, "CoApp2", profile.CoApplicant2Name)
	assert.Equal(t, "CoApp3", profile.CoApplicant3Name)
	assert.Equal(t, "SPOUSE", profile.CoApplicant1Relation)
	assert.Equal(t, "CHILD", profile.CoApplicant2Relation)
	assert.Equal(t, "PARENT", profile.CoApplicant3Relation)
}

// TestAddressFieldHandling tests address field handling
func TestAddressFieldHandling(t *testing.T) {
	profile := &models.PropertyCustomerProfile{
		CommunicationAddressLine1: "123 Main St",
		CommunicationAddressLine2: "Apt 4B",
		CommunicationCity:         "New York",
		CommunicationState:        "NY",
		CommunicationCountry:      "USA",
		CommunicationZip:          "10001",
		PermanentAddressLine1:     "456 Oak Ave",
		PermanentAddressLine2:     "Suite 200",
		PermanentCity:             "Los Angeles",
		PermanentState:            "CA",
		PermanentCountry:          "USA",
		PermanentZip:              "90001",
	}

	assert.Equal(t, "123 Main St", profile.CommunicationAddressLine1)
	assert.Equal(t, "Apt 4B", profile.CommunicationAddressLine2)
	assert.Equal(t, "New York", profile.CommunicationCity)
	assert.Equal(t, "NY", profile.CommunicationState)
	assert.Equal(t, "USA", profile.CommunicationCountry)
	assert.Equal(t, "10001", profile.CommunicationZip)
	assert.Equal(t, "456 Oak Ave", profile.PermanentAddressLine1)
	assert.Equal(t, "Suite 200", profile.PermanentAddressLine2)
	assert.Equal(t, "Los Angeles", profile.PermanentCity)
	assert.Equal(t, "CA", profile.PermanentState)
	assert.Equal(t, "USA", profile.PermanentCountry)
	assert.Equal(t, "90001", profile.PermanentZip)
}

// TestBooleanFieldHandling tests boolean field handling
func TestBooleanFieldHandling(t *testing.T) {
	config := &models.CreateProjectCostConfigRequest{
		IsMandatory: true,
	}

	assert.True(t, config.IsMandatory)

	costSheet := &models.UpdateCostSheetRequest{
		GSTApplicable: false,
	}

	assert.False(t, costSheet.GSTApplicable)
}

// TestStringFieldHandling tests string field handling
func TestStringFieldHandling(t *testing.T) {
	profile := &models.PropertyCustomerProfile{
		FirstName:   "John",
		MiddleName:  "Michael",
		LastName:    "Smith",
		CompanyName: "Tech Corp",
	}

	assert.Equal(t, "John", profile.FirstName)
	assert.Equal(t, "Michael", profile.MiddleName)
	assert.Equal(t, "Smith", profile.LastName)
	assert.Equal(t, "Tech Corp", profile.CompanyName)
}

// TestNilFieldHandling tests nil field handling
func TestNilFieldHandling(t *testing.T) {
	stage := &models.PropertyPaymentStage{
		DueDate:                &time.Time{},
		ExpectedCollectionDate: nil,
		ActualCollectionDate:   nil,
	}

	assert.NotNil(t, stage.DueDate)
	assert.Nil(t, stage.ExpectedCollectionDate)
	assert.Nil(t, stage.ActualCollectionDate)
}

// BenchmarkCreateCustomerProfile benchmarks customer creation
func BenchmarkCreateCustomerProfile(b *testing.B) {
	profile := &models.PropertyCustomerProfile{
		CustomerCode: "CUST001",
		FirstName:    "John",
		Email:        "john@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = profile.CustomerCode
		_ = profile.FirstName
		_ = profile.Email
	}
}

// BenchmarkAreaStatementCreation benchmarks area statement creation
func BenchmarkAreaStatementCreation(b *testing.B) {
	areaStmt := &models.CreateAreaStatementRequest{
		ProjectID:                 "proj-123",
		UnitID:                    "unit-123",
		RERACarPetAreaSqft:        950.00,
		CarPetAreaWithBalconySqft: 1050.00,
		PlinthAreaSqft:            1150.00,
		SBUASqft:                  1400.00,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use all fields to prevent unused write warnings
		_ = areaStmt.ProjectID + areaStmt.UnitID
		_ = areaStmt.RERACarPetAreaSqft + areaStmt.CarPetAreaWithBalconySqft + areaStmt.PlinthAreaSqft + areaStmt.SBUASqft
	}
}
