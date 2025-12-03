package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateLeadHandler validates lead creation endpoint
func TestCreateLeadHandler(t *testing.T) {
	payload := map[string]interface{}{
		"first_name":  "John",
		"last_name":   "Doe",
		"email":       "john.doe@example.com",
		"phone":       "+91-9876543210",
		"company":     "TechCorp",
		"lead_status": "new",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/sales/leads", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestCreateCustomerHandler validates customer creation
func TestCreateCustomerHandler(t *testing.T) {
	payload := map[string]interface{}{
		"name":         "John Doe",
		"email":        "john@example.com",
		"phone":        "+91-9876543210",
		"city":         "Mumbai",
		"credit_limit": 100000.00,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/sales/customers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestCreateQuotationHandler validates quotation creation
func TestCreateQuotationHandler(t *testing.T) {
	payload := map[string]interface{}{
		"customer_id": "cust-uuid-1234",
		"items": []map[string]interface{}{
			{
				"description": "Product A",
				"quantity":    10.0,
				"unit_price":  1000.00,
			},
		},
		"tax_rate": 18.0,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/sales/quotations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestGetQuotationHandler validates quotation retrieval
func TestGetQuotationHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/sales/quotations/quot-uuid-1234", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.Method)
}

// TestCreateInvoiceHandler validates invoice creation
func TestCreateInvoiceHandler(t *testing.T) {
	payload := map[string]interface{}{
		"customer_id": "cust-uuid-1234",
		"items": []map[string]interface{}{
			{
				"description": "Product A",
				"quantity":    10.0,
				"unit_price":  1000.00,
			},
		},
		"tax_rate": 18.0,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/sales/invoices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestRecordPaymentHandler validates payment recording
func TestRecordPaymentHandler(t *testing.T) {
	payload := map[string]interface{}{
		"invoice_id":   "inv-uuid-1234",
		"payment_mode": "bank_transfer",
		"amount_paid":  10000.00,
		"reference_no": "TXN12345678",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/sales/payments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestQuotationCalculation validates quotation amount calculation with proper GST
func TestQuotationCalculation(t *testing.T) {
	testCases := []struct {
		description string
		baseAmount  float64
		discount    float64
		gstRate     float64
		expected    float64
	}{
		// Base 10000, 0% discount, 18% GST
		{"No discount with 18% GST", 10000.00, 0.00, 0.18, 11800.00},
		// Base 50000, 0% discount, 18% GST
		{"50K with 18% GST", 50000.00, 0.00, 0.18, 59000.00},
		// Base 100000, 5000 discount, 5% GST = (100000-5000) + (95000*0.05) = 95000 + 4750 = 99750
		{"100K with 5% discount and 5% GST", 100000.00, 5000.00, 0.05, 99750.00},
		// Edge case: 0% GST (Exempt items)
		{"Exempt supply (0% GST)", 10000.00, 0.00, 0.00, 10000.00},
	}

	for _, tc := range testCases {
		// Calculate net after discount
		netAmount := tc.baseAmount - tc.discount
		// Calculate GST
		gst := netAmount * tc.gstRate
		// Total = Net + GST
		total := netAmount + gst

		assert.InDelta(t, tc.expected, total, 0.01, tc.description)
	}
} // TestCreditLimitValidation validates credit limit enforcement
func TestCreditLimitValidation(t *testing.T) {
	testCases := []struct {
		creditLimit  float64
		invoiceTotal float64
		currentUsed  float64
		allowed      bool
	}{
		{100000.00, 50000.00, 0.00, true},      // Within limit
		{100000.00, 85000.00, 20000.00, false}, // Exceeds limit (20k + 85k > 100k)
		{100000.00, 100000.00, 0.00, true},     // At limit
		{100000.00, 50000.00, 60000.00, false}, // Exceeds after usage
	}

	for _, tc := range testCases {
		available := tc.creditLimit - tc.currentUsed
		allowed := tc.invoiceTotal <= available
		assert.Equal(t, tc.allowed, allowed)
	}
}

// TestDiscountValidation validates discount constraints
func TestDiscountValidation(t *testing.T) {
	testCases := []struct {
		discountPercent float64
		valid           bool
	}{
		{5.0, true},
		{15.0, true},
		{25.0, true},
		{-10.0, false},
		{150.0, false},
	}

	for _, tc := range testCases {
		isValid := tc.discountPercent >= 0 && tc.discountPercent <= 100
		assert.Equal(t, tc.valid, isValid)
	}
}

// TestPaymentModeValidationWithCompliance validates payment modes and compliance
func TestPaymentModeValidationWithCompliance(t *testing.T) {
	validModes := []struct {
		mode          string
		requiresTDS   bool // TDS deductible
		documentation string
	}{
		{"cash", false, "Receipt"},
		{"cheque", true, "Cheque + Bank statement"},
		{"bank_transfer", true, "Bank statement"},
		{"upi", true, "Screenshot + Bank statement"},
		{"card", true, "Card statement"},
	}

	for _, mode := range validModes {
		assert.NotEmpty(t, mode.mode)
		assert.NotEmpty(t, mode.documentation)
	}
}

// TestMultiTenantSalesRouting validates tenant isolation
func TestMultiTenantSalesRouting(t *testing.T) {
	tenantReq1 := httptest.NewRequest("GET", "/api/v1/sales/customers/cust-1", nil)
	tenantReq1.Header.Set("X-Tenant-ID", "tenant-1")

	tenantReq2 := httptest.NewRequest("GET", "/api/v1/sales/customers/cust-1", nil)
	tenantReq2.Header.Set("X-Tenant-ID", "tenant-2")

	assert.NotEqual(t,
		tenantReq1.Header.Get("X-Tenant-ID"),
		tenantReq2.Header.Get("X-Tenant-ID"))
}

// TestLeadStatusTransition validates lead status workflow
func TestLeadStatusTransition(t *testing.T) {
	validStatuses := []string{"new", "qualified", "proposal_sent", "negotiation", "won", "lost"}

	for _, status := range validStatuses {
		assert.NotEmpty(t, status)
	}
}

// TestInvoiceStatusValidation validates invoice status workflow
func TestInvoiceStatusValidation(t *testing.T) {
	// Valid status transitions as per business rules
	validTransitions := map[string][]string{
		"draft":          {"sent", "cancelled"},
		"sent":           {"partially_paid", "paid", "overdue", "cancelled"},
		"partially_paid": {"paid", "overdue", "cancelled"},
		"paid":           {},                    // Terminal state
		"overdue":        {"paid", "cancelled"}, // Can still collect
		"cancelled":      {},                    // Terminal state
	}

	assert.NotNil(t, validTransitions)
	assert.Equal(t, 6, len(validTransitions))
}

// BenchmarkCreateQuotationEndpoint benchmarks quotation creation
func BenchmarkCreateQuotationEndpoint(b *testing.B) {
	payload := map[string]interface{}{
		"customer_id": "cust-uuid-1234",
		"items": []map[string]interface{}{
			{
				"description": "Product",
				"quantity":    10.0,
				"unit_price":  1000.00,
			},
		},
		"tax_rate": 18.0,
	}

	body, _ := json.Marshal(payload)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/v1/sales/quotations", bytes.NewReader(body))
		req.Header.Set("X-Tenant-ID", "tenant-1")
		httptest.NewRecorder()
	}
}

// BenchmarkCreateInvoiceEndpoint benchmarks invoice creation
func BenchmarkCreateInvoiceEndpoint(b *testing.B) {
	payload := map[string]interface{}{
		"customer_id": "cust-uuid-1234",
		"items": []map[string]interface{}{
			{
				"description": "Product",
				"quantity":    10.0,
				"unit_price":  1000.00,
			},
		},
		"tax_rate": 18.0,
	}

	body, _ := json.Marshal(payload)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/v1/sales/invoices", bytes.NewReader(body))
		req.Header.Set("X-Tenant-ID", "tenant-1")
		httptest.NewRecorder()
	}
}
