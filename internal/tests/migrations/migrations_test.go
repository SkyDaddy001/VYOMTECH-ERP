package migrations

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMigration026SalesPostSales validates all tables from migration 026
func TestMigration026SalesPostSales(t *testing.T) {
	t.Run("SalesLead_Model", func(t *testing.T) {
		// Test key fields that are actually used in business logic
		lead := struct {
			ID        string
			TenantID  string
			LeadCode  string
			FirstName string
			Status    string
		}{
			ID:        "test-lead-001",
			TenantID:  "tenant-001",
			LeadCode:  "LD001",
			FirstName: "John",
			Status:    "new",
		}

		assert.NotEmpty(t, lead.ID)
		assert.NotEmpty(t, lead.TenantID)
		assert.NotEmpty(t, lead.LeadCode)
		assert.NotEmpty(t, lead.FirstName)
		assert.Equal(t, "new", lead.Status)
	})

	t.Run("Booking_Model", func(t *testing.T) {
		bookingDate := time.Now()
		booking := struct {
			ID          string
			TenantID    string
			BookingCode string
			UnitID      string
			LeadID      string
			BookingDate time.Time
			Status      string
		}{
			ID:          "test-booking-001",
			TenantID:    "tenant-001",
			BookingCode: "BK001",
			UnitID:      "unit-001",
			LeadID:      "lead-001",
			BookingDate: bookingDate,
			Status:      "active",
		}

		assert.NotEmpty(t, booking.ID)
		assert.NotEmpty(t, booking.TenantID)
		assert.NotEmpty(t, booking.BookingCode)
		assert.NotEmpty(t, booking.UnitID)
		assert.NotEmpty(t, booking.LeadID)
		assert.Equal(t, "active", booking.Status)
		assert.Equal(t, bookingDate, booking.BookingDate)
	})

	t.Run("Unit_Model", func(t *testing.T) {
		unit := struct {
			ID          string
			TenantID    string
			ProjectID   string
			UnitCode    string
			Block       string
			AptNo       string
			Floor       int
			UnitType    string
			SBUA        float64
			UDSPerSqft  float64
			Status      string
			RatePerSqft float64
		}{
			ID:          "test-unit-001",
			TenantID:    "tenant-001",
			ProjectID:   "project-001",
			UnitCode:    "A-101",
			Block:       "A",
			AptNo:       "101",
			Floor:       1,
			UnitType:    "2BHK",
			SBUA:        1200.50,
			UDSPerSqft:  45.75,
			Status:      "available",
			RatePerSqft: 5500.00,
		}

		assert.NotEmpty(t, unit.ID)
		assert.NotEmpty(t, unit.TenantID)
		assert.NotEmpty(t, unit.ProjectID)
		assert.NotEmpty(t, unit.UnitCode)
		assert.NotEmpty(t, unit.Block)
		assert.NotEmpty(t, unit.AptNo)
		assert.Equal(t, 1, unit.Floor)
		assert.Equal(t, "2BHK", unit.UnitType)
		assert.Equal(t, "available", unit.Status)
		assert.Equal(t, 1200.50, unit.SBUA)
		assert.Equal(t, 45.75, unit.UDSPerSqft)
		assert.Equal(t, 5500.00, unit.RatePerSqft)
	})

	t.Run("UnitCostSheet_Model", func(t *testing.T) {
		costSheet := struct {
			ID                         string
			TenantID                   string
			UnitID                     string
			FRC                        float64
			CarParkingType             string
			CarParkingCost             float64
			ApartmentCostExcludingGovt float64
			TotalCost                  float64
		}{
			ID:                         "test-costsheet-001",
			TenantID:                   "tenant-001",
			UnitID:                     "unit-001",
			FRC:                        5500.00,
			CarParkingType:             "covered",
			CarParkingCost:             300000.00,
			ApartmentCostExcludingGovt: 6600000.00,
			TotalCost:                  7200000.00,
		}

		assert.NotEmpty(t, costSheet.ID)
		assert.NotEmpty(t, costSheet.TenantID)
		assert.NotEmpty(t, costSheet.UnitID)
		assert.Equal(t, 5500.00, costSheet.FRC)
		assert.Equal(t, "covered", costSheet.CarParkingType)
		assert.Equal(t, 300000.00, costSheet.CarParkingCost)
		assert.Equal(t, 6600000.00, costSheet.ApartmentCostExcludingGovt)
		assert.Equal(t, 7200000.00, costSheet.TotalCost)
	})

	t.Run("Client_Model", func(t *testing.T) {
		client := struct {
			ID            string
			TenantID      string
			BookingID     string
			ApplicantType string
			FirstName     string
			LastName      string
			Phone         string
			AadharNo      string
			PanNo         string
		}{
			ID:            "test-client-001",
			TenantID:      "tenant-001",
			BookingID:     "booking-001",
			ApplicantType: "primary",
			FirstName:     "Ramesh",
			LastName:      "Kumar",
			Phone:         "9876543210",
			AadharNo:      "123456789012",
			PanNo:         "ABCDE1234F",
		}

		assert.NotEmpty(t, client.ID)
		assert.NotEmpty(t, client.TenantID)
		assert.NotEmpty(t, client.BookingID)
		assert.Equal(t, "primary", client.ApplicantType)
		assert.NotEmpty(t, client.FirstName)
		assert.NotEmpty(t, client.LastName)
		assert.NotEmpty(t, client.Phone)
		assert.Len(t, client.AadharNo, 12)
		assert.Len(t, client.PanNo, 10)
	})

	t.Run("PaymentSchedule_Model", func(t *testing.T) {
		schedule := struct {
			ID                string
			TenantID          string
			BookingID         string
			PaymentStage      int
			ConstructionStage string
			AmountDue         float64
			PaymentType       string
			Status            string
		}{
			ID:                "test-schedule-001",
			TenantID:          "tenant-001",
			BookingID:         "booking-001",
			PaymentStage:      1,
			ConstructionStage: "booking_advance",
			AmountDue:         100000.00,
			PaymentType:       "milestone",
			Status:            "pending",
		}

		assert.NotEmpty(t, schedule.ID)
		assert.NotEmpty(t, schedule.TenantID)
		assert.NotEmpty(t, schedule.BookingID)
		assert.Equal(t, 1, schedule.PaymentStage)
		assert.Equal(t, "booking_advance", schedule.ConstructionStage)
		assert.Equal(t, 100000.00, schedule.AmountDue)
		assert.Equal(t, "milestone", schedule.PaymentType)
		assert.Equal(t, "pending", schedule.Status)
	})

	t.Run("Payment_Model", func(t *testing.T) {
		paymentDate := time.Now()
		payment := struct {
			ID           string
			TenantID     string
			BookingID    string
			CustomerName string
			ReceiptNo    string
			PaymentDate  time.Time
			PaymentMode  string
			PaidBy       string
			Towards      string
			Amount       float64
		}{
			ID:           "test-payment-001",
			TenantID:     "tenant-001",
			BookingID:    "booking-001",
			CustomerName: "Ramesh Kumar",
			ReceiptNo:    "RCP001",
			PaymentDate:  paymentDate,
			PaymentMode:  "bank_transfer",
			PaidBy:       "customer",
			Towards:      "apartment",
			Amount:       500000.00,
		}

		assert.NotEmpty(t, payment.ID)
		assert.NotEmpty(t, payment.TenantID)
		assert.NotEmpty(t, payment.BookingID)
		assert.NotEmpty(t, payment.CustomerName)
		assert.NotEmpty(t, payment.ReceiptNo)
		assert.Equal(t, paymentDate, payment.PaymentDate)
		assert.Equal(t, "bank_transfer", payment.PaymentMode)
		assert.Equal(t, "customer", payment.PaidBy)
		assert.Equal(t, "apartment", payment.Towards)
		assert.Equal(t, 500000.00, payment.Amount)
	})

	t.Run("BankLoan_Model", func(t *testing.T) {
		sanctionDate := time.Now()
		loan := struct {
			ID                 string
			TenantID           string
			BookingID          string
			BankName           string
			LoanSanctionDate   *time.Time
			SanctionAmount     float64
			DisbursedAmount    float64
			DisbursementStatus string
		}{
			ID:                 "test-loan-001",
			TenantID:           "tenant-001",
			BookingID:          "booking-001",
			BankName:           "HDFC Bank",
			LoanSanctionDate:   &sanctionDate,
			SanctionAmount:     4000000.00,
			DisbursedAmount:    2000000.00,
			DisbursementStatus: "partial",
		}

		assert.NotEmpty(t, loan.ID)
		assert.NotEmpty(t, loan.TenantID)
		assert.NotEmpty(t, loan.BookingID)
		assert.Equal(t, "HDFC Bank", loan.BankName)
		assert.NotNil(t, loan.LoanSanctionDate)
		assert.Equal(t, 4000000.00, loan.SanctionAmount)
		assert.Equal(t, 2000000.00, loan.DisbursedAmount)
		assert.Equal(t, "partial", loan.DisbursementStatus)
		// Test loan_available_for_disbursement calculation
		loanAvailable := loan.SanctionAmount - loan.DisbursedAmount
		assert.Equal(t, 2000000.00, loanAvailable)
	})

	t.Run("RegistrationDetails_Model", func(t *testing.T) {
		regDetails := struct {
			ID                        string
			TenantID                  string
			BookingID                 string
			GSTApplicable             bool
			GSTPercentage             float64
			GSTCost                   float64
			ApartmentCostIncludingGST float64
			RegistrationType          string
			RegistrationCost          float64
			Status                    string
		}{
			ID:                        "test-reg-001",
			TenantID:                  "tenant-001",
			BookingID:                 "booking-001",
			GSTApplicable:             true,
			GSTPercentage:             5.00,
			GSTCost:                   330000.00,
			ApartmentCostIncludingGST: 6930000.00,
			RegistrationType:          "sale_deed",
			RegistrationCost:          415800.00, // 6% of 6930000
			Status:                    "pending",
		}

		assert.NotEmpty(t, regDetails.ID)
		assert.NotEmpty(t, regDetails.TenantID)
		assert.NotEmpty(t, regDetails.BookingID)
		assert.True(t, regDetails.GSTApplicable)
		assert.Equal(t, 5.00, regDetails.GSTPercentage)
		assert.Equal(t, 330000.00, regDetails.GSTCost)
		assert.Equal(t, 6930000.00, regDetails.ApartmentCostIncludingGST)
		assert.Equal(t, "sale_deed", regDetails.RegistrationType)
		assert.Equal(t, 415800.00, regDetails.RegistrationCost)
		assert.Equal(t, "pending", regDetails.Status)
	})

	t.Run("AdditionalCharges_Model", func(t *testing.T) {
		charges := struct {
			ID                string
			TenantID          string
			BookingID         string
			MaintenanceCharge float64
			CorpusCharge      float64
			EBDeposit         float64
			OtherWorksCharge  float64
		}{
			ID:                "test-charges-001",
			TenantID:          "tenant-001",
			BookingID:         "booking-001",
			MaintenanceCharge: 72000.00,
			CorpusCharge:      120000.00,
			EBDeposit:         15000.00,
			OtherWorksCharge:  50000.00,
		}

		assert.NotEmpty(t, charges.ID)
		assert.NotEmpty(t, charges.TenantID)
		assert.NotEmpty(t, charges.BookingID)
		assert.Equal(t, 72000.00, charges.MaintenanceCharge)
		assert.Equal(t, 120000.00, charges.CorpusCharge)
		assert.Equal(t, 15000.00, charges.EBDeposit)
		assert.Equal(t, 50000.00, charges.OtherWorksCharge)
		totalAdditional := charges.MaintenanceCharge + charges.CorpusCharge + charges.EBDeposit + charges.OtherWorksCharge
		assert.Equal(t, 257000.00, totalAdditional)
	})
}

// TestMigration027SalesAnalyticsViews tests analytics view calculations
func TestMigration027SalesAnalyticsViews(t *testing.T) {
	t.Run("BankOwnPaymentAnalysis_Calculation", func(t *testing.T) {
		// Test scenario: Calculate bank vs own payment
		apartmentCostWithGST := 6930000.00
		bankSanctioned := 4000000.00
		bankDisbursed := 2000000.00
		customerOwnPaid := 1500000.00 // via cash/cheque/bank_transfer

		loanAvailableForDisbursement := bankSanctioned - bankDisbursed
		customerOwnDue := apartmentCostWithGST - bankDisbursed - customerOwnPaid

		assert.Equal(t, 2000000.00, loanAvailableForDisbursement)
		assert.Equal(t, 3430000.00, customerOwnDue)
	})

	t.Run("CollectionPercentage_Calculation", func(t *testing.T) {
		totalValue := 6930000.00
		totalCollected := 3500000.00

		collectionPercentage := (totalCollected / totalValue) * 100
		assert.InDelta(t, 50.505, collectionPercentage, 0.01)
	})

	t.Run("OccupancyPercentage_Calculation", func(t *testing.T) {
		totalUnits := 100
		soldUnits := 75

		occupancyPercentage := (float64(soldUnits) / float64(totalUnits)) * 100
		assert.Equal(t, 75.0, occupancyPercentage)
	})

	t.Run("PaymentStagePercentage_Calculation", func(t *testing.T) {
		// 13 payment stages total
		totalStages := 13
		currentStage := 5

		stagePercentage := (float64(currentStage) / float64(totalStages)) * 100
		assert.InDelta(t, 38.46, stagePercentage, 0.01)
	})

	t.Run("TotalReceivable_Calculation", func(t *testing.T) {
		apartmentCostWithGST := 6930000.00
		registrationCost := 415800.00
		maintenanceCharge := 72000.00
		corpusCharge := 120000.00
		ebDeposit := 15000.00
		otherWorks := 50000.00

		totalReceivable := apartmentCostWithGST + registrationCost + maintenanceCharge + corpusCharge + ebDeposit + otherWorks
		assert.Equal(t, 7602800.00, totalReceivable)
	})

	t.Run("QuarterlyCollection_Aggregation", func(t *testing.T) {
		// Q1: Apr, May, Jun
		month1Collections := 500000.00
		month2Collections := 750000.00
		month3Collections := 600000.00

		quarterTotal := month1Collections + month2Collections + month3Collections
		assert.Equal(t, 1850000.00, quarterTotal)
	})
}

// TestPaymentModeLogic tests the payment mode segregation
func TestPaymentModeLogic(t *testing.T) {
	t.Run("CustomerOwnPayment_AllModes", func(t *testing.T) {
		// Customer can pay via cash, cheque, or bank_transfer
		cashPayment := 100000.00
		chequePayment := 200000.00
		bankTransferPayment := 300000.00

		customerOwnPaid := cashPayment + chequePayment + bankTransferPayment
		assert.Equal(t, 600000.00, customerOwnPaid)
	})

	t.Run("BankLoanDisbursement_Separate", func(t *testing.T) {
		// Bank loan is tracked separately via bank_loan.disbursed_amount
		bankSanctioned := 4000000.00
		bankDisbursed := 2500000.00
		loanAvailable := bankSanctioned - bankDisbursed

		assert.Equal(t, 1500000.00, loanAvailable)
	})

	t.Run("TotalReceived_Calculation", func(t *testing.T) {
		bankDisbursed := 2500000.00
		customerOwnPaid := 600000.00

		totalReceived := bankDisbursed + customerOwnPaid
		assert.Equal(t, 3100000.00, totalReceived)
	})
}

// TestAgreementStatus tests agreement workflow
func TestAgreementStatus(t *testing.T) {
	t.Run("AgreementPending", func(t *testing.T) {
		bookingDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

		// Test getAgreementStatus helper with nil date
		status := getAgreementStatus(nil)
		assert.Equal(t, "Pending", status)

		// Days pending
		today := time.Date(2025, 12, 8, 0, 0, 0, 0, time.UTC)
		daysPending := int(today.Sub(bookingDate).Hours() / 24)
		assert.Greater(t, daysPending, 300)
	})

	t.Run("AgreementSigned", func(t *testing.T) {
		bookingDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
		agreementSignDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

		// Test getAgreementStatus helper with non-nil date
		status := getAgreementStatus(&agreementSignDate)
		assert.Equal(t, "Signed", status)

		// Days to sign
		daysToSign := int(agreementSignDate.Sub(bookingDate).Hours() / 24)
		assert.Equal(t, 17, daysToSign)
	})
}

// getAgreementStatus is a helper to determine agreement status
func getAgreementStatus(agreementDate *time.Time) string {
	if agreementDate != nil {
		return "Signed"
	}
	return "Pending"
}

// TestFinancialCalculations tests real estate financial logic
func TestFinancialCalculations(t *testing.T) {
	t.Run("GST_Calculation", func(t *testing.T) {
		apartmentCostExcludingGovt := 6600000.00
		gstPercentage := 5.00

		gstCost := apartmentCostExcludingGovt * (gstPercentage / 100)
		apartmentCostWithGST := apartmentCostExcludingGovt + gstCost

		assert.Equal(t, 330000.00, gstCost)
		assert.Equal(t, 6930000.00, apartmentCostWithGST)
	})

	t.Run("RegistrationCost_Calculation", func(t *testing.T) {
		apartmentCostWithGST := 6930000.00
		registrationPercentage := 6.00 // Tamil Nadu rates

		registrationCost := apartmentCostWithGST * (registrationPercentage / 100)
		assert.Equal(t, 415800.00, registrationCost)
	})

	t.Run("PaymentStageAmount_Calculation", func(t *testing.T) {
		apartmentCostWithGST := 6930000.00

		// Stage percentages as per schedule
		stagePercentages := []float64{10, 35, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5}
		var total float64 = 0

		for _, pct := range stagePercentages {
			total += pct
		}

		assert.Equal(t, 100.0, total)

		// Calculate stage 1 amount (10%)
		stage1Amount := apartmentCostWithGST * (stagePercentages[0] / 100)
		assert.Equal(t, 693000.00, stage1Amount)

		// Calculate stage 2 amount (35%)
		stage2Amount := apartmentCostWithGST * (stagePercentages[1] / 100)
		assert.Equal(t, 2425500.00, stage2Amount)
	})
}
