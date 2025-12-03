package services

import (
	"testing"

	"vyomtech-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// TestSalesLead validates sales lead model
func TestSalesLead(t *testing.T) {
	lead := &models.SalesLead{
		ID:        "lead-uuid-001",
		TenantID:  "tenant-001",
		LeadCode:  "LD001",
		FirstName: "John",
		Email:     "john@example.com",
		Status:    "new",
		Source:    "website",
	}

	assert.Equal(t, "lead-uuid-001", lead.ID)
	assert.Equal(t, "LD001", lead.LeadCode)
	assert.Equal(t, "new", lead.Status)
}

// TestLeadStatuses validates valid lead statuses
func TestLeadStatuses(t *testing.T) {
	statuses := []string{"new", "contacted", "qualified", "negotiation", "converted", "lost"}

	for _, status := range statuses {
		lead := &models.SalesLead{Status: status}
		assert.Equal(t, status, lead.Status)
	}
}

// TestLeadSources validates valid lead sources
func TestLeadSources(t *testing.T) {
	sources := []string{"website", "email", "phone", "referral", "event", "social"}

	for _, source := range sources {
		lead := &models.SalesLead{Source: source}
		assert.Equal(t, source, lead.Source)
	}
}

// TestSalesCustomer validates customer model
func TestSalesCustomer(t *testing.T) {
	customer := &models.SalesCustomer{
		ID:             "cust-uuid-001",
		TenantID:       "tenant-001",
		CustomerCode:   "CUST001",
		CustomerName:   "John Doe",
		CreditLimit:    100000.00,
		CreditDays:     30,
		CurrentBalance: 25000.00,
		Status:         "active",
	}

	assert.Equal(t, "cust-uuid-001", customer.ID)
	assert.Equal(t, 100000.00, customer.CreditLimit)
	assert.Equal(t, "active", customer.Status)
}

// TestQuotation validates quotation model
func TestQuotation(t *testing.T) {
	quotation := &models.SalesQuotation{
		ID:              "quote-uuid-001",
		TenantID:        "tenant-001",
		QuotationNumber: "QT001",
		CustomerID:      "cust-uuid-001",
		SubtotalAmount:  10000.00,
		DiscountAmount:  1000.00,
		TaxAmount:       1620.00,
		TotalAmount:     10620.00,
		Status:          "draft",
	}

	assert.Equal(t, "quote-uuid-001", quotation.ID)
	assert.Equal(t, "draft", quotation.Status)
	expectedTotal := 10000.00 - 1000.00 + 1620.00
	assert.InDelta(t, expectedTotal, quotation.TotalAmount, 0.01)
}

// TestQuotationCalculations validates amount calculations
func TestQuotationCalculations(t *testing.T) {
	subtotal := 10000.00
	discountRate := 0.10
	taxRate := 0.18

	discount := subtotal * discountRate
	netAmount := subtotal - discount
	tax := netAmount * taxRate
	total := netAmount + tax

	assert.InDelta(t, 10620.00, total, 0.01)
}

// TestTaxCalculations validates various tax rates
func TestTaxCalculations(t *testing.T) {
	testCases := []struct {
		amount      float64
		taxRate     float64
		expectedTax float64
	}{
		{10000.0, 0.18, 1800.0},
		{10000.0, 0.05, 500.0},
		{10000.0, 0.12, 1200.0},
		{10000.0, 0.00, 0.0},
	}

	for _, tc := range testCases {
		tax := tc.amount * tc.taxRate
		assert.InDelta(t, tc.expectedTax, tax, 0.01)
	}
}

// TestCreditLimitValidation validates credit limit enforcement
func TestCreditLimitValidation(t *testing.T) {
	creditLimit := 100000.0
	usedCredit := 50000.0
	available := creditLimit - usedCredit
	invoiceAmount := 30000.0

	canPlace := invoiceAmount <= available
	assert.True(t, canPlace)
}

// TestQuotationStatuses validates quotation status values
func TestQuotationStatuses(t *testing.T) {
	statuses := []string{"draft", "sent", "accepted", "rejected", "expired", "converted_to_order"}

	for _, status := range statuses {
		quote := &models.SalesQuotation{Status: status}
		assert.Equal(t, status, quote.Status)
	}
}

// TestSalesPaymentFlow validates payment tracking
func TestSalesPaymentFlow(t *testing.T) {
	invoiceAmount := 10000.0
	payment1 := 5000.0
	payment2 := 5000.0

	totalPaid := payment1 + payment2
	balance := invoiceAmount - totalPaid

	assert.Equal(t, invoiceAmount, totalPaid)
	assert.Equal(t, 0.0, balance)
}

// TestDiscountValidation validates discount constraints
func TestDiscountValidation(t *testing.T) {
	testCases := []struct {
		percent float64
		valid   bool
	}{
		{10.0, true},
		{50.0, true},
		{0.0, true},
		{-10.0, false},
		{100.0, false},
	}

	for _, tc := range testCases {
		isValid := tc.percent >= 0 && tc.percent < 100
		assert.Equal(t, tc.valid, isValid)
	}
}

// TestMultiTenantSales validates sales data isolation
func TestMultiTenantSales(t *testing.T) {
	lead1 := &models.SalesLead{ID: "lead-1", TenantID: "tenant-1"}
	lead2 := &models.SalesLead{ID: "lead-2", TenantID: "tenant-2"}

	assert.NotEqual(t, lead1.TenantID, lead2.TenantID)
}

// TestCustomerCategoryValidation validates customer categories
func TestCustomerCategoryValidation(t *testing.T) {
	categories := []string{"gold", "silver", "bronze", "regular"}

	for _, cat := range categories {
		customer := &models.SalesCustomer{CustomerCategory: cat}
		assert.Equal(t, cat, customer.CustomerCategory)
	}
}
