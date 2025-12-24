package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// ============================================================
// OAUTH INTEGRATION TESTS
// ============================================================

// TestGoogleOAuthFlow tests the complete Google OAuth flow
func TestGoogleOAuthFlow(t *testing.T) {
	// Setup
	router := setupTestRouter()
	platform := "google"
	clientID := "mock-google-client-id"
	clientSecret := "mock-google-client-secret"
	redirectURI := "http://localhost:3000/oauth/callback/google"
	scope := "https://www.googleapis.com/auth/adwords"

	// Step 1: Authorize
	t.Run("Step 1: Authorization", func(t *testing.T) {
		authReq := MockOAuthRequest{
			ClientID:    clientID,
			RedirectURI: redirectURI,
			Scope:       scope,
			State:       "test_state_123",
			Platform:    platform,
		}

		body, _ := json.Marshal(authReq)
		req, _ := http.NewRequest("POST", "/mock/oauth/google/authorize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if resp["code"] == nil {
			t.Error("Expected authorization code in response")
		}
		if resp["state"] != "test_state_123" {
			t.Error("State mismatch")
		}
	})

	// Step 2: Token Exchange
	t.Run("Step 2: Token Exchange", func(t *testing.T) {
		// First get an auth code
		authReq := MockOAuthRequest{
			ClientID:    clientID,
			RedirectURI: redirectURI,
			Scope:       scope,
			State:       "test_state_123",
			Platform:    platform,
		}
		body, _ := json.Marshal(authReq)
		req, _ := http.NewRequest("POST", "/mock/oauth/google/authorize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var authResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &authResp)
		authCode := authResp["code"].(string)

		// Exchange for tokens
		tokenReq := MockOAuthTokenRequest{
			Code:         authCode,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURI:  redirectURI,
			GrantType:    "authorization_code",
		}

		body, _ = json.Marshal(tokenReq)
		req, _ = http.NewRequest("POST", "/mock/oauth/google/token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var tokenResp MockOAuthToken
		json.Unmarshal(w.Body.Bytes(), &tokenResp)

		if tokenResp.AccessToken == "" {
			t.Error("Expected access token in response")
		}
		if tokenResp.RefreshToken == "" {
			t.Error("Expected refresh token in response")
		}
		if tokenResp.AccountID != "123-456-7890" {
			t.Error("Expected Google account ID in response")
		}
	})

	// Step 3: Token Validation
	t.Run("Step 3: Token Validation", func(t *testing.T) {
		// Get valid token first
		authReq := MockOAuthRequest{
			ClientID:    clientID,
			RedirectURI: redirectURI,
			Scope:       scope,
			State:       "test_state_123",
			Platform:    platform,
		}
		body, _ := json.Marshal(authReq)
		req, _ := http.NewRequest("POST", "/mock/oauth/google/authorize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var authResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &authResp)
		authCode := authResp["code"].(string)

		tokenReq := MockOAuthTokenRequest{
			Code:         authCode,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURI:  redirectURI,
			GrantType:    "authorization_code",
		}
		body, _ = json.Marshal(tokenReq)
		req, _ = http.NewRequest("POST", "/mock/oauth/google/token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var tokenResp MockOAuthToken
		json.Unmarshal(w.Body.Bytes(), &tokenResp)

		// Now validate the token
		req, _ = http.NewRequest("GET", "/mock/oauth/google/validate?access_token="+tokenResp.AccessToken, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var validateResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &validateResp)

		if validateResp["valid"] != true {
			t.Error("Expected valid=true")
		}
	})

	// Step 4: Token Refresh
	t.Run("Step 4: Token Refresh", func(t *testing.T) {
		// Get tokens first
		authReq := MockOAuthRequest{
			ClientID:    clientID,
			RedirectURI: redirectURI,
			Scope:       scope,
			State:       "test_state_123",
			Platform:    platform,
		}
		body, _ := json.Marshal(authReq)
		req, _ := http.NewRequest("POST", "/mock/oauth/google/authorize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var authResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &authResp)
		authCode := authResp["code"].(string)

		tokenReq := MockOAuthTokenRequest{
			Code:         authCode,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURI:  redirectURI,
			GrantType:    "authorization_code",
		}
		body, _ = json.Marshal(tokenReq)
		req, _ = http.NewRequest("POST", "/mock/oauth/google/token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var tokenResp MockOAuthToken
		json.Unmarshal(w.Body.Bytes(), &tokenResp)

		// Refresh token
		refreshReq := struct {
			RefreshToken string `json:"refresh_token"`
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		}{
			RefreshToken: tokenResp.RefreshToken,
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}

		body, _ = json.Marshal(refreshReq)
		req, _ = http.NewRequest("POST", "/mock/oauth/google/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var refreshResp MockOAuthToken
		json.Unmarshal(w.Body.Bytes(), &refreshResp)

		if refreshResp.AccessToken == "" {
			t.Error("Expected new access token")
		}
	})

	// Step 5: Token Revocation
	t.Run("Step 5: Token Revocation", func(t *testing.T) {
		// Get token first
		authReq := MockOAuthRequest{
			ClientID:    clientID,
			RedirectURI: redirectURI,
			Scope:       scope,
			State:       "test_state_123",
			Platform:    platform,
		}
		body, _ := json.Marshal(authReq)
		req, _ := http.NewRequest("POST", "/mock/oauth/google/authorize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var authResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &authResp)
		authCode := authResp["code"].(string)

		tokenReq := MockOAuthTokenRequest{
			Code:         authCode,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURI:  redirectURI,
			GrantType:    "authorization_code",
		}
		body, _ = json.Marshal(tokenReq)
		req, _ = http.NewRequest("POST", "/mock/oauth/google/token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var tokenResp MockOAuthToken
		json.Unmarshal(w.Body.Bytes(), &tokenResp)

		// Revoke token
		revokeReq := struct {
			AccessToken string `json:"access_token"`
		}{
			AccessToken: tokenResp.AccessToken,
		}

		body, _ = json.Marshal(revokeReq)
		req, _ = http.NewRequest("POST", "/mock/oauth/google/revoke", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var revokeResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &revokeResp)

		if revokeResp["success"] != true {
			t.Error("Expected success=true")
		}
	})
}

// TestMetaOAuthFlow tests the complete Meta OAuth flow
func TestMetaOAuthFlow(t *testing.T) {
	// Setup
	router := setupTestRouter()
	platform := "meta"
	clientID := "mock-meta-client-id"
	clientSecret := "mock-meta-client-secret"
	redirectURI := "http://localhost:3000/oauth/callback/meta"
	scope := "ads_management,business_management"

	t.Run("Meta: Full OAuth Flow", func(t *testing.T) {
		// Authorize
		authReq := MockOAuthRequest{
			ClientID:    clientID,
			RedirectURI: redirectURI,
			Scope:       scope,
			State:       "test_state_meta_123",
			Platform:    platform,
		}

		body, _ := json.Marshal(authReq)
		req, _ := http.NewRequest("POST", "/mock/oauth/meta/authorize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var authResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &authResp)
		authCode := authResp["code"].(string)

		// Exchange for tokens
		tokenReq := MockOAuthTokenRequest{
			Code:         authCode,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURI:  redirectURI,
			GrantType:    "authorization_code",
		}

		body, _ = json.Marshal(tokenReq)
		req, _ = http.NewRequest("POST", "/mock/oauth/meta/token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var tokenResp MockOAuthToken
		json.Unmarshal(w.Body.Bytes(), &tokenResp)

		if tokenResp.AccountID != "act_1234567890" {
			t.Error("Expected Meta account ID in response")
		}
	})
}

// TestOAuthErrorHandling tests error scenarios
func TestOAuthErrorHandling(t *testing.T) {
	router := setupTestRouter()

	t.Run("Invalid Client ID", func(t *testing.T) {
		authReq := MockOAuthRequest{
			ClientID:    "wrong-client-id",
			RedirectURI: "http://localhost:3000/oauth/callback/google",
			Scope:       "https://www.googleapis.com/auth/adwords",
			State:       "test_state",
			Platform:    "google",
		}

		body, _ := json.Marshal(authReq)
		req, _ := http.NewRequest("POST", "/mock/oauth/google/authorize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Invalid Redirect URI", func(t *testing.T) {
		authReq := MockOAuthRequest{
			ClientID:    "mock-google-client-id",
			RedirectURI: "http://wrong-uri.com",
			Scope:       "https://www.googleapis.com/auth/adwords",
			State:       "test_state",
			Platform:    "google",
		}

		body, _ := json.Marshal(authReq)
		req, _ := http.NewRequest("POST", "/mock/oauth/google/authorize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Invalid Authorization Code", func(t *testing.T) {
		tokenReq := MockOAuthTokenRequest{
			Code:         "invalid_code",
			ClientID:     "mock-google-client-id",
			ClientSecret: "mock-google-client-secret",
			RedirectURI:  "http://localhost:3000/oauth/callback/google",
			GrantType:    "authorization_code",
		}

		body, _ := json.Marshal(tokenReq)
		req, _ := http.NewRequest("POST", "/mock/oauth/google/token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

// TestTestDataEndpoints tests the test data retrieval endpoints
func TestTestDataEndpoints(t *testing.T) {
	router := setupTestRouter()

	t.Run("Get All Test Data", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test/data", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("Get Google Ads Test Data", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test/data/google-ads/campaign", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("Get Meta Ads Test Data", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test/data/meta-ads/campaign", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("Invalid Test Data Key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test/data/google-ads/invalid-key", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

// ============================================================
// TEST HELPER FUNCTIONS
// ============================================================

// setupTestRouter creates a test router with all routes registered
func setupTestRouter() *gin.Engine {
	// Use test mode for Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Register routes
	RegisterMockOAuthRoutes(router)
	RegisterTestDataRoutes(router)

	return router
}

// Benchmarks
func BenchmarkOAuthAuthorize(b *testing.B) {
	router := setupTestRouter()

	authReq := MockOAuthRequest{
		ClientID:    "mock-google-client-id",
		RedirectURI: "http://localhost:3000/oauth/callback/google",
		Scope:       "https://www.googleapis.com/auth/adwords",
		State:       "test_state",
		Platform:    "google",
	}

	body, _ := json.Marshal(authReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/mock/oauth/google/authorize", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkOAuthToken(b *testing.B) {
	router := setupTestRouter()

	// First get auth code
	authReq := MockOAuthRequest{
		ClientID:    "mock-google-client-id",
		RedirectURI: "http://localhost:3000/oauth/callback/google",
		Scope:       "https://www.googleapis.com/auth/adwords",
		State:       "test_state",
		Platform:    "google",
	}
	body, _ := json.Marshal(authReq)
	req, _ := http.NewRequest("POST", "/mock/oauth/google/authorize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var authResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &authResp)
	authCode := authResp["code"].(string)

	tokenReq := MockOAuthTokenRequest{
		Code:         authCode,
		ClientID:     "mock-google-client-id",
		ClientSecret: "mock-google-client-secret",
		RedirectURI:  "http://localhost:3000/oauth/callback/google",
		GrantType:    "authorization_code",
	}

	body, _ = json.Marshal(tokenReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/mock/oauth/google/token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
