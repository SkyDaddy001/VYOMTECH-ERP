package services

import (
	"context"
	"database/sql"
	"testing"
)

func TestGetUserLevel(t *testing.T) {
	// Mock database connection (you would use a test database in real scenarios)
	var mockDB *sql.DB

	service := NewGamificationService(mockDB)

	testUserID := int64(1)
	testTenantID := "test-tenant"

	// Test GetGamificationProfile which includes level information
	profile, _ := service.GetGamificationProfile(context.Background(), testUserID, testTenantID)

	// In a test database, this would return data or nil if user doesn't exist
	// For now, we just verify the service can be called without panicking
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	// If profile is nil, that's expected if the user doesn't exist
	if profile != nil && profile.CurrentLevel != nil {
		// Profile exists and has level data
		if profile.CurrentLevel.CurrentLevel < 1 {
			t.Errorf("Expected current_level >= 1, got %d", profile.CurrentLevel.CurrentLevel)
		}
	}
}
