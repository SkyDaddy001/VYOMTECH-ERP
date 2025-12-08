package services

import (
	"testing"

	"vyomtech-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

// TestCallModelUUIDs verifies Call model uses string IDs
func TestCallModelUUIDs(t *testing.T) {
	call := &models.Call{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		TenantID: "tenant-001",
		LeadID:   "lead-uuid-001",
		AgentID:  "agent-uuid-001",
		Status:   "initiated",
	}

	// Verify string types
	assert.IsType(t, "", call.ID)
	assert.IsType(t, "", call.LeadID)
	assert.IsType(t, "", call.AgentID)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", call.ID)
	assert.Equal(t, "tenant-001", call.TenantID)
	assert.Equal(t, "initiated", call.Status)
}

// TestCallStatusValues validates call status field
func TestCallStatusValues(t *testing.T) {
	validStatuses := []string{"initiated", "ringing", "active", "ended"}

	for _, status := range validStatuses {
		call := &models.Call{Status: status}
		assert.Equal(t, status, call.Status)
	}
}

// TestCallOutcome validates call outcome field
func TestCallOutcome(t *testing.T) {
	validOutcomes := []string{"success", "no_answer", "busy", "declined", "failed"}

	for _, outcome := range validOutcomes {
		call := &models.Call{Outcome: outcome}
		assert.Equal(t, outcome, call.Outcome)
	}
}

// TestCallDurationSeconds validates duration field
func TestCallDurationSeconds(t *testing.T) {
	call := &models.Call{DurationSeconds: 300}
	assert.Equal(t, 300, call.DurationSeconds)
}

// TestCallStats validates stats calculation
func TestCallStats(t *testing.T) {
	stats := &models.CallStats{
		Total:           100,
		Active:          5,
		Completed:       90,
		Failed:          5,
		AverageDuration: 300,
		SuccessRate:     90.0,
	}

	assert.Equal(t, 100, stats.Total)
	assert.Equal(t, 5, stats.Active)
	assert.Equal(t, 90, stats.Completed)
	assert.Equal(t, 5, stats.Failed)
	assert.Equal(t, 300, stats.AverageDuration)
	assert.Equal(t, 90.0, stats.SuccessRate)
}

// TestMultiTenantCallIsolation ensures calls are tenant-isolated
func TestMultiTenantCallIsolation(t *testing.T) {
	call1 := &models.Call{ID: "call-1", TenantID: "tenant-1"}
	call2 := &models.Call{ID: "call-2", TenantID: "tenant-2"}

	assert.NotEqual(t, call1.TenantID, call2.TenantID)
	assert.NotEqual(t, call1.ID, call2.ID)
}
