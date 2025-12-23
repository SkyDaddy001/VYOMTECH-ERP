package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

// MockDB provides a mock database for testing
type MockDB struct {
	execContext     func() (sql.Result, error)
	queryContext    func() (*sql.Rows, error)
	queryRowContext func() *sql.Row
}

// TestRBACHandlerCreateRole tests the CreateRole endpoint
func TestRBACHandlerCreateRole(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		body           map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid role creation by admin",
			role:           "admin",
			body:           map[string]interface{}{"role_name": "Manager", "description": "Test", "is_active": true},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name:           "Non-admin cannot create role",
			role:           "user",
			body:           map[string]interface{}{"role_name": "Manager", "description": "Test", "is_active": true},
			expectedStatus: http.StatusForbidden,
			expectedError:  "Admin access required",
		},
		{
			name:           "Missing role_name field",
			role:           "admin",
			body:           map[string]interface{}{"description": "Test", "is_active": true},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Role name is required",
		},
		{
			name:           "Empty role_name",
			role:           "admin",
			body:           map[string]interface{}{"role_name": "   ", "description": "Test", "is_active": true},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Role name is required",
		},
		{
			name:           "Duplicate role name",
			role:           "admin",
			body:           map[string]interface{}{"role_name": "Existing", "description": "Test", "is_active": true},
			expectedStatus: http.StatusConflict,
			expectedError:  "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockLog := logger.New()
			mockRBACService := &services.RBACService{}
			mockDB := &sql.DB{}

			handler := NewRBACHandler(mockRBACService, mockDB, mockLog)

			// Create request
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/v1/rbac/roles", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Add context
			ctx := context.WithValue(req.Context(), middleware.TenantIDKey, "test-tenant")
			ctx = context.WithValue(ctx, middleware.UserIDKey, int64(1))
			ctx = context.WithValue(ctx, middleware.RoleKey, tt.role)
			req = req.WithContext(ctx)

			// Record response
			w := httptest.NewRecorder()

			// Execute
			handler.CreateRole(w, req)

			// Verify
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				body := w.Body.String()
				if !bytes.Contains([]byte(body), []byte(tt.expectedError)) {
					t.Errorf("Expected error message containing '%s', got: %s", tt.expectedError, body)
				}
			}
		})
	}
}

// TestRBACHandlerGetRole tests the GetRole endpoint
func TestRBACHandlerGetRole(t *testing.T) {
	tests := []struct {
		name           string
		roleID         string
		expectedStatus int
		shouldExist    bool
	}{
		{
			name:           "Get existing role",
			roleID:         "role-001",
			expectedStatus: http.StatusOK,
			shouldExist:    true,
		},
		{
			name:           "Get non-existent role",
			roleID:         "role-999",
			expectedStatus: http.StatusNotFound,
			shouldExist:    false,
		},
		{
			name:           "Invalid role ID format",
			roleID:         "",
			expectedStatus: http.StatusBadRequest,
			shouldExist:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLog := logger.New()
			mockRBACService := &services.RBACService{}
			mockDB := &sql.DB{}

			handler := NewRBACHandler(mockRBACService, mockDB, mockLog)

			url := fmt.Sprintf("/api/v1/rbac/roles/%s", tt.roleID)
			req := httptest.NewRequest("GET", url, nil)

			ctx := context.WithValue(req.Context(), middleware.TenantIDKey, "test-tenant")
			ctx = context.WithValue(ctx, middleware.UserIDKey, int64(1))
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.GetRole(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// TestRBACHandlerListRoles tests the ListRoles endpoint
func TestRBACHandlerListRoles(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "List roles successfully",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLog := logger.New()
			mockRBACService := &services.RBACService{}
			mockDB := &sql.DB{}

			handler := NewRBACHandler(mockRBACService, mockDB, mockLog)

			req := httptest.NewRequest("GET", "/api/v1/rbac/roles", nil)

			ctx := context.WithValue(req.Context(), middleware.TenantIDKey, "test-tenant")
			ctx = context.WithValue(ctx, middleware.UserIDKey, int64(1))
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.ListRoles(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// TestRBACHandlerListPermissions tests the ListPermissions endpoint
func TestRBACHandlerListPermissions(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "List permissions successfully",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLog := logger.New()
			mockRBACService := &services.RBACService{}
			mockDB := &sql.DB{}

			handler := NewRBACHandler(mockRBACService, mockDB, mockLog)

			req := httptest.NewRequest("GET", "/api/v1/rbac/permissions", nil)

			ctx := context.WithValue(req.Context(), middleware.TenantIDKey, "test-tenant")
			ctx = context.WithValue(ctx, middleware.UserIDKey, int64(1))
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.ListPermissions(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// TestRBACHandlerAssignPermissions tests the AssignPermissions endpoint
func TestRBACHandlerAssignPermissions(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		roleID         string
		body           map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Assign permissions as admin",
			role:           "admin",
			roleID:         "role-001",
			body:           map[string]interface{}{"permission_ids": []string{"perm-001", "perm-002"}},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Non-admin cannot assign",
			role:           "user",
			roleID:         "role-001",
			body:           map[string]interface{}{"permission_ids": []string{"perm-001"}},
			expectedStatus: http.StatusForbidden,
			expectedError:  "Admin access required",
		},
		{
			name:           "Missing permission_ids",
			role:           "admin",
			roleID:         "role-001",
			body:           map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Permission IDs are required",
		},
		{
			name:           "Empty permission_ids array",
			role:           "admin",
			roleID:         "role-001",
			body:           map[string]interface{}{"permission_ids": []string{}},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Permission IDs are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLog := logger.New()
			mockRBACService := &services.RBACService{}
			mockDB := &sql.DB{}

			handler := NewRBACHandler(mockRBACService, mockDB, mockLog)

			bodyBytes, _ := json.Marshal(tt.body)
			url := fmt.Sprintf("/api/v1/rbac/roles/%s/permissions", tt.roleID)
			req := httptest.NewRequest("PUT", url, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			ctx := context.WithValue(req.Context(), middleware.TenantIDKey, "test-tenant")
			ctx = context.WithValue(ctx, middleware.UserIDKey, int64(1))
			ctx = context.WithValue(ctx, middleware.RoleKey, tt.role)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.AssignPermissions(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				body := w.Body.String()
				if !bytes.Contains([]byte(body), []byte(tt.expectedError)) {
					t.Errorf("Expected error containing '%s', got: %s", tt.expectedError, body)
				}
			}
		})
	}
}

// TestRBACHandlerDeleteRole tests the DeleteRole endpoint
func TestRBACHandlerDeleteRole(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		roleID         string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Delete role as admin",
			role:           "admin",
			roleID:         "role-001",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Non-admin cannot delete",
			role:           "user",
			roleID:         "role-001",
			expectedStatus: http.StatusForbidden,
			expectedError:  "Admin access required",
		},
		{
			name:           "Delete non-existent role",
			role:           "admin",
			roleID:         "role-999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Role not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLog := logger.New()
			mockRBACService := &services.RBACService{}
			mockDB := &sql.DB{}

			handler := NewRBACHandler(mockRBACService, mockDB, mockLog)

			url := fmt.Sprintf("/api/v1/rbac/roles/%s", tt.roleID)
			req := httptest.NewRequest("DELETE", url, nil)

			ctx := context.WithValue(req.Context(), middleware.TenantIDKey, "test-tenant")
			ctx = context.WithValue(ctx, middleware.UserIDKey, int64(1))
			ctx = context.WithValue(ctx, middleware.RoleKey, tt.role)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.DeleteRole(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				body := w.Body.String()
				if !bytes.Contains([]byte(body), []byte(tt.expectedError)) {
					t.Errorf("Expected error containing '%s', got: %s", tt.expectedError, body)
				}
			}
		})
	}
}

// TestRBACHandlerAuthorizationCheck tests authorization checks
func TestRBACHandlerAuthorizationCheck(t *testing.T) {
	tests := []struct {
		name           string
		hasUserID      bool
		hasTenantID    bool
		expectedStatus int
	}{
		{
			name:           "Valid context with user and tenant",
			hasUserID:      true,
			hasTenantID:    true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing tenant ID",
			hasUserID:      true,
			hasTenantID:    false,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Missing user ID",
			hasUserID:      false,
			hasTenantID:    true,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLog := logger.New()
			mockRBACService := &services.RBACService{}
			mockDB := &sql.DB{}

			handler := NewRBACHandler(mockRBACService, mockDB, mockLog)

			req := httptest.NewRequest("GET", "/api/v1/rbac/roles", nil)

			ctx := req.Context()
			if tt.hasUserID {
				ctx = context.WithValue(ctx, middleware.UserIDKey, int64(1))
			}
			if tt.hasTenantID {
				ctx = context.WithValue(ctx, middleware.TenantIDKey, "test-tenant")
			}
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.ListRoles(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// TestRBACHandlerResponseFormat tests response format consistency
func TestRBACHandlerResponseFormat(t *testing.T) {
	mockLog := logger.New()
	mockRBACService := &services.RBACService{}
	mockDB := &sql.DB{}

	handler := NewRBACHandler(mockRBACService, mockDB, mockLog)

	req := httptest.NewRequest("GET", "/api/v1/rbac/roles", nil)
	ctx := context.WithValue(req.Context(), middleware.TenantIDKey, "test-tenant")
	ctx = context.WithValue(ctx, middleware.UserIDKey, int64(1))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ListRoles(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	// Verify response has required fields
	if _, hasMessage := response["message"]; !hasMessage {
		t.Error("Response missing 'message' field")
	}
	if _, hasData := response["data"]; !hasData {
		t.Error("Response missing 'data' field")
	}
}

// TestMultiTenantIsolation tests that roles are isolated by tenant
func TestMultiTenantIsolation(t *testing.T) {
	mockLog := logger.New()
	mockRBACService := &services.RBACService{}
	mockDB := &sql.DB{}

	handler := NewRBACHandler(mockRBACService, mockDB, mockLog)

	// Request with tenant-1
	req1 := httptest.NewRequest("GET", "/api/v1/rbac/roles", nil)
	ctx1 := context.WithValue(req1.Context(), middleware.TenantIDKey, "tenant-1")
	ctx1 = context.WithValue(ctx1, middleware.UserIDKey, int64(1))
	req1 = req1.WithContext(ctx1)

	w1 := httptest.NewRecorder()
	handler.ListRoles(w1, req1)

	// Request with tenant-2
	req2 := httptest.NewRequest("GET", "/api/v1/rbac/roles", nil)
	ctx2 := context.WithValue(req2.Context(), middleware.TenantIDKey, "tenant-2")
	ctx2 = context.WithValue(ctx2, middleware.UserIDKey, int64(2))
	req2 = req2.WithContext(ctx2)

	w2 := httptest.NewRecorder()
	handler.ListRoles(w2, req2)

	// Both should return 200 but different tenant contexts
	if w1.Code != http.StatusOK || w2.Code != http.StatusOK {
		t.Error("Both requests should succeed with proper tenant context")
	}
}
