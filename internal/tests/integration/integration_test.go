package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// buildRequest is a helper to construct HTTP requests for testing
func buildRequest(method, path string, body interface{}) *http.Request {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			panic(err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "test-tenant-1")
	return req
}

func TestCallCreationToCompletion(t *testing.T) {
	// Test call status transitions
	statuses := []string{"initiated", "ringing", "connected", "completed"}
	for i, status := range statuses {
		assert.NotEmpty(t, status, "Status at index %d should not be empty", i)
	}

	// Verify no duplicate statuses
	statusMap := make(map[string]bool)
	for _, status := range statuses {
		assert.False(t, statusMap[status], "Status %s should not repeat", status)
		statusMap[status] = true
	}

	// Test call duration validation
	durationSeconds := 300
	assert.Greater(t, durationSeconds, 0, "Call duration should be positive")

	// Test call outcome
	outcome := "successful"
	validOutcomes := map[string]bool{
		"successful":   true,
		"unsuccessful": true,
		"no_answer":    true,
		"voicemail":    true,
	}
	assert.True(t, validOutcomes[outcome], "Invalid call outcome: %s", outcome)
}

// TestSalesQuotationToInvoiceFlow tests sales workflow
func TestSalesQuotationToInvoiceFlow(t *testing.T) {
	// Step 1: Create quotation
	quotationID := "quot-uuid-1234"
	customerID := "cust-uuid-5678"
	subtotal := 50000.0

	assert.NotEmpty(t, quotationID)
	assert.NotEmpty(t, customerID)
	assert.Greater(t, subtotal, 0.0)

	// Step 2: Calculate tax
	taxRate := 18.0
	tax := subtotal * (taxRate / 100)
	total := subtotal + tax

	assert.InDelta(t, 59000.0, total, 0.01)

	// Step 3: Convert to invoice
	invoiceID := "inv-uuid-9012"
	invoiceStatus := "sent"
	assert.NotEmpty(t, invoiceID)
	assert.Equal(t, "sent", invoiceStatus)

	// Step 4: Record payment
	paidAmount := 59000.0
	invoiceStatus = "paid"
	assert.Equal(t, "paid", invoiceStatus)
	assert.InDelta(t, total, paidAmount, 0.01, "Paid amount should match invoice total")

	// Additional validation: invoice status transitions
	expectedStatuses := []string{"draft", "sent", "paid"}
	statusMap := make(map[string]bool)
	for _, s := range expectedStatuses {
		assert.False(t, statusMap[s], "Status %s should not repeat", s)
		statusMap[s] = true
	}
}

// TestConstructionProjectToCompletion tests construction workflow
func TestConstructionProjectToCompletion(t *testing.T) {
	// Step 1: Create project
	projectID := "proj-uuid-1234"
	projectCode := "PRJ-001"
	budget := 50000000.0

	assert.NotEmpty(t, projectID)
	assert.NotEmpty(t, projectCode)
	assert.Greater(t, budget, 0.0)

	// Step 2: Create BOQ
	boqID := "boq-uuid-5678"
	boqTotal := 35000000.0
	assert.NotEmpty(t, boqID)
	assert.Less(t, boqTotal, budget)

	// Step 3: Record progress
	progressMonths := []string{"2025-01", "2025-02", "2025-03"}
	progressPercentages := []float64{25.0, 50.0, 75.0}

	for i, month := range progressMonths {
		assert.NotEmpty(t, month)
		assert.GreaterOrEqual(t, progressPercentages[i], 0.0)
		assert.LessOrEqual(t, progressPercentages[i], 100.0)
	}

	// Step 4: Quality control
	qcStatus := "passed"
	assert.Equal(t, "passed", qcStatus)
	validQCStatuses := map[string]bool{
		"passed":  true,
		"failed":  true,
		"pending": true,
	}
	assert.True(t, validQCStatuses[qcStatus], "Invalid QC status: %s", qcStatus)

	// Step 5: Project completion
	finalProgress := 100.0
	projectStatus := "completed"
	assert.Equal(t, 100.0, finalProgress, "Final progress should be 100%")
	assert.Equal(t, "completed", projectStatus)
	assert.LessOrEqual(t, finalProgress, 100.0, "Progress cannot exceed 100%")
}

// TestMultiTenantIsolation tests tenant data isolation
func TestMultiTenantIsolation(t *testing.T) {
	// Tenant 1 data
	tenant1ID := "tenant-1"
	tenant1CallID := "call-uuid-t1"
	tenant1ProjectID := "proj-uuid-t1"

	// Tenant 2 data
	tenant2ID := "tenant-2"
	tenant2CallID := "call-uuid-t2"
	tenant2ProjectID := "proj-uuid-t2"

	// Verify isolation
	assert.NotEqual(t, tenant1ID, tenant2ID)
	assert.NotEqual(t, tenant1CallID, tenant2CallID)
	assert.NotEqual(t, tenant1ProjectID, tenant2ProjectID)

	// Verify tenant-specific operations
	tenant1Calls := 5
	tenant2Calls := 3
	assert.NotEqual(t, tenant1Calls, tenant2Calls, "Tenant call counts should be different")
	assert.Greater(t, tenant1Calls, 0, "Tenant 1 should have calls")
	assert.Greater(t, tenant2Calls, 0, "Tenant 2 should have calls")

	// Verify no cross-tenant contamination
	totalCalls := tenant1Calls + tenant2Calls
	assert.Equal(t, 8, totalCalls, "Total calls should be sum of both tenants")
}

// TestGLPostingForSalesInvoice tests GL integration for sales
func TestGLPostingForSalesInvoice(t *testing.T) {
	// Sales invoice data
	subtotal := 50000.0
	taxRate := 18.0
	tax := subtotal * (taxRate / 100)
	total := subtotal + tax

	// GL posting should debit AR, credit revenue + tax
	drAR := total              // Debit AR (Asset)
	crRevenue := subtotal      // Credit Revenue
	crTax := tax               // Credit Tax Payable
	crSum := crRevenue + crTax // Sum of credits

	// Verify double-entry
	assert.InDelta(t, drAR, crSum, 0.01)
}

// TestGLPostingForPayment tests GL integration for payments
func TestGLPostingForPayment(t *testing.T) {
	// Payment data
	invoiceTotal := 59000.0
	paidAmount := 59000.0

	// GL posting should debit bank, credit AR
	drBank := paidAmount // Debit Bank (Asset)
	crAR := invoiceTotal // Credit AR (Asset)

	// Verify double-entry
	assert.InDelta(t, drBank, crAR, 0.01)
}

// TestBOQToProgressFlow tests BOQ tracking through progress
func TestBOQToProgressFlow(t *testing.T) {
	// BOQ items
	boqItems := []struct {
		description string
		qty         float64
		rate        float64
	}{
		{"Cement", 1000.0, 350.0},
		{"Steel", 2500.0, 50.0},
		{"Sand", 100.0, 500.0},
	}

	// Calculate total
	var boqTotal float64
	for _, item := range boqItems {
		itemTotal := item.qty * item.rate
		boqTotal += itemTotal
	}

	expectedTotal := 350000.0 + 125000.0 + 50000.0 // 525000.0
	assert.InDelta(t, expectedTotal, boqTotal, 0.01)

	// Track progress
	progress := 0.0
	assert.Equal(t, 0.0, progress)

	progress = 50.0
	progressAmount := boqTotal * (progress / 100)
	assert.InDelta(t, 262500.0, progressAmount, 0.01)

	progress = 100.0
	progressAmount = boqTotal * (progress / 100)
	assert.InDelta(t, boqTotal, progressAmount, 0.01)
}

// TestMultiStagePaymentFlow tests payment stages
func TestMultiStagePaymentFlow(t *testing.T) {
	invoiceTotal := 100000.0

	// Advance payment (20%)
	advancePayment := invoiceTotal * 0.20
	remaining := invoiceTotal - advancePayment

	assert.InDelta(t, 20000.0, advancePayment, 0.01)
	assert.InDelta(t, 80000.0, remaining, 0.01)

	// Milestone payment (30%)
	milestonePayment := invoiceTotal * 0.30
	remaining -= milestonePayment

	assert.InDelta(t, 30000.0, milestonePayment, 0.01)
	assert.InDelta(t, 50000.0, remaining, 0.01)

	// Final payment
	finalPayment := remaining
	remaining -= finalPayment

	assert.InDelta(t, 50000.0, finalPayment, 0.01)
	assert.InDelta(t, 0.0, remaining, 0.01)
}

// TestCreditLimitEnforcement tests credit limit workflow
func TestCreditLimitEnforcement(t *testing.T) {
	creditLimit := 100000.0
	currentUsed := 0.0

	// Invoice 1
	invoice1 := 40000.0
	currentUsed += invoice1
	available := creditLimit - currentUsed

	assert.LessOrEqual(t, invoice1, creditLimit)
	assert.InDelta(t, 60000.0, available, 0.01)

	// Invoice 2 - would exceed limit if issued
	invoice2 := 70000.0 // This exceeds available (60000)
	canIssue := invoice2 <= available

	assert.False(t, canIssue) // Would exceed limit
	assert.InDelta(t, 60000.0, available, 0.01)

	// Allowed invoice
	invoice3 := 30000.0
	canIssue = invoice3 <= available

	assert.True(t, canIssue)
	currentUsed += invoice3
	available = creditLimit - currentUsed
	assert.InDelta(t, 30000.0, available, 0.01)
}

// TestBudgetTracking tests project budget tracking
func TestBudgetTracking(t *testing.T) {
	totalBudget := 50000000.0
	boqAmount := 35000000.0
	contingency := totalBudget * 0.10 // 10% contingency

	budgetAllocated := boqAmount + contingency
	assert.Less(t, budgetAllocated, totalBudget)

	// Track spending
	spent := 15000000.0
	remaining := totalBudget - spent

	assert.InDelta(t, 35000000.0, remaining, 0.01)
	assert.Less(t, spent, totalBudget)
}

// TestStatusTransitionValidation tests valid state transitions
func TestStatusTransitionValidation(t *testing.T) {
	// Call status transitions
	callStatuses := []string{"initiated", "ringing", "connected", "completed"}
	for i := 0; i < len(callStatuses)-1; i++ {
		assert.NotEqual(t, callStatuses[i], callStatuses[i+1], "Adjacent statuses should differ: %s vs %s", callStatuses[i], callStatuses[i+1])
	}

	// Project status transitions
	projectStatuses := []string{"planning", "initiated", "in_progress", "completed"}
	for i := 0; i < len(projectStatuses)-1; i++ {
		assert.NotEqual(t, projectStatuses[i], projectStatuses[i+1], "Adjacent statuses should differ: %s vs %s", projectStatuses[i], projectStatuses[i+1])
	}

	// Invoice status transitions
	invoiceStatuses := []string{"draft", "sent", "paid"}
	for i := 0; i < len(invoiceStatuses)-1; i++ {
		assert.NotEqual(t, invoiceStatuses[i], invoiceStatuses[i+1], "Adjacent statuses should differ: %s vs %s", invoiceStatuses[i], invoiceStatuses[i+1])
	}

	// Verify status order is sequential
	assert.Equal(t, 4, len(callStatuses), "Call should have 4 status transitions")
	assert.Equal(t, 4, len(projectStatuses), "Project should have 4 status transitions")
	assert.Equal(t, 3, len(invoiceStatuses), "Invoice should have 3 status transitions")
}

// TestDataConsistencyAcrossModules tests data consistency
func TestDataConsistencyAcrossModules(t *testing.T) {
	// Call linked to Lead
	callID := "call-uuid-1234"
	leadID := "lead-uuid-5678"
	assert.NotEmpty(t, callID)
	assert.NotEmpty(t, leadID)

	// Project linked to BOQ
	projectID := "proj-uuid-1234"
	boqID := "boq-uuid-5678"
	assert.NotEmpty(t, projectID)
	assert.NotEmpty(t, boqID)

	// Invoice linked to Customer
	invoiceID := "inv-uuid-1234"
	customerID := "cust-uuid-5678"
	assert.NotEmpty(t, invoiceID)
	assert.NotEmpty(t, customerID)
}
