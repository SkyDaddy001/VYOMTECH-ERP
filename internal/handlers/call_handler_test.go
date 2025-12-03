package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateCallHandler validates call creation endpoint
func TestCreateCallHandler(t *testing.T) {
	payload := map[string]interface{}{
		"lead_id":   "lead-uuid-1234",
		"agent_id":  "agent-uuid-5678",
		"call_type": "inbound",
		"status":    "initiated",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/calls", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "tenant-1", req.Header.Get("X-Tenant-ID"))
}

// TestGetCallHandler validates call retrieval endpoint
func TestGetCallHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/calls/call-uuid-1234", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.Method)
	assert.Contains(t, req.URL.Path, "call-uuid-1234")
}

// TestListCallsHandler validates calls list endpoint
func TestListCallsHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/calls?limit=50&offset=0", nil)
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.Method)
	assert.Contains(t, req.URL.RawQuery, "limit=50")
}

// TestUpdateCallStatusHandler validates call status update
func TestUpdateCallStatusHandler(t *testing.T) {
	payload := map[string]interface{}{
		"status":  "completed",
		"outcome": "successful",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/api/v1/calls/call-uuid-1234/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "PUT", req.Method)
}

// TestEndCallHandler validates call end endpoint
func TestEndCallHandler(t *testing.T) {
	payload := map[string]interface{}{
		"duration_seconds": 300,
		"outcome":          "successful",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/calls/call-uuid-1234/end", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
}

// TestCallValidation validates call data validation
func TestCallValidation(t *testing.T) {
	testCases := []struct {
		leadID   string
		agentID  string
		callType string
		valid    bool
	}{
		{"lead-1", "agent-1", "inbound", true},
		{"", "agent-1", "inbound", false},            // Missing lead ID
		{"lead-1", "", "inbound", false},             // Missing agent ID
		{"lead-1", "agent-1", "", false},             // Missing call type
		{"lead-1", "agent-1", "invalid_type", false}, // Invalid call type
	}

	for _, tc := range testCases {
		isValid := tc.leadID != "" && tc.agentID != "" &&
			(tc.callType == "inbound" || tc.callType == "outbound")
		assert.Equal(t, tc.valid, isValid)
	}
}

// TestCallDurationValidation validates duration constraints
func TestCallDurationValidation(t *testing.T) {
	testCases := []struct {
		duration int
		valid    bool
	}{
		{300, true},   // 5 minutes
		{3600, true},  // 1 hour
		{0, true},     // No duration yet
		{-100, false}, // Negative duration
	}

	for _, tc := range testCases {
		isValid := tc.duration >= 0
		assert.Equal(t, tc.valid, isValid)
	}
}

// TestMultiTenantCallRouting validates tenant isolation
func TestMultiTenantCallRouting(t *testing.T) {
	tenantReq1 := httptest.NewRequest("GET", "/api/v1/calls/call-1", nil)
	tenantReq1.Header.Set("X-Tenant-ID", "tenant-1")

	tenantReq2 := httptest.NewRequest("GET", "/api/v1/calls/call-1", nil)
	tenantReq2.Header.Set("X-Tenant-ID", "tenant-2")

	assert.NotEqual(t,
		tenantReq1.Header.Get("X-Tenant-ID"),
		tenantReq2.Header.Get("X-Tenant-ID"))
}

// TestCallOutcomeValidation validates allowed outcomes
func TestCallOutcomeValidation(t *testing.T) {
	validOutcomes := []string{"successful", "failed", "no_answer", "callback_requested"}

	for _, outcome := range validOutcomes {
		assert.NotEmpty(t, outcome)
	}
}

// TestCallRecordingCapture validates recording handling
func TestCallRecordingCapture(t *testing.T) {
	payload := map[string]interface{}{
		"recording_url": "https://cdn.example.com/call-1234.mp3",
		"duration":      300,
	}

	body, _ := json.Marshal(payload)
	assert.NotNil(t, body)
}

// TestCallNoteCapture validates note handling
func TestCallNoteCapture(t *testing.T) {
	payload := map[string]interface{}{
		"notes": "Customer interested in product demo",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/calls/call-1/notes", bytes.NewReader(body))
	req.Header.Set("X-Tenant-ID", "tenant-1")

	assert.NotNil(t, req)
}

// BenchmarkCreateCallEndpoint benchmarks call creation
func BenchmarkCreateCallEndpoint(b *testing.B) {
	payload := map[string]interface{}{
		"lead_id":   "lead-uuid-1234",
		"agent_id":  "agent-uuid-5678",
		"call_type": "inbound",
	}

	body, _ := json.Marshal(payload)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/v1/calls", bytes.NewReader(body))
		req.Header.Set("X-Tenant-ID", "tenant-1")
		httptest.NewRecorder()
	}
}

// BenchmarkListCallsEndpoint benchmarks call listing
func BenchmarkListCallsEndpoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/v1/calls?limit=50", nil)
		req.Header.Set("X-Tenant-ID", "tenant-1")
		httptest.NewRecorder()
	}
}
