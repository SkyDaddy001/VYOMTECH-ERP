package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ============================================================
// MOCK OAUTH SERVICE FOR DEVELOPMENT & TESTING
// ============================================================
// This file provides mock OAuth implementations for testing
// without connecting to real Google Ads or Meta APIs
// Perfect for development, testing, and CI/CD pipelines

// MockOAuthProvider represents a mock OAuth provider
type MockOAuthProvider struct {
	ProviderName string // "google" or "meta"
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Tokens       map[string]*MockOAuthToken // state -> token
}

// MockOAuthToken represents a mock OAuth token response
type MockOAuthToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	IssuedAt     time.Time `json:"issued_at"`
	Scope        string    `json:"scope"`
	AccountID    string    `json:"account_id"`
	AccountName  string    `json:"account_name"`
}

// MockOAuthRequest represents an OAuth authorization request
type MockOAuthRequest struct {
	ClientID    string `json:"client_id" binding:"required"`
	RedirectURI string `json:"redirect_uri" binding:"required"`
	Scope       string `json:"scope" binding:"required"`
	State       string `json:"state" binding:"required"`
	Platform    string `json:"platform" binding:"required"` // "google" or "meta"
}

// MockOAuthTokenRequest represents a token exchange request
type MockOAuthTokenRequest struct {
	Code         string `json:"code" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	RedirectURI  string `json:"redirect_uri" binding:"required"`
	GrantType    string `json:"grant_type" binding:"required"`
}

// MockOAuthCallbackRequest represents the callback data
type MockOAuthCallbackRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}

// ============================================================
// GLOBAL MOCK PROVIDERS
// ============================================================

var mockProviders = map[string]*MockOAuthProvider{
	"google": {
		ProviderName: "google",
		ClientID:     "mock-google-client-id",
		ClientSecret: "mock-google-client-secret",
		RedirectURI:  "http://localhost:3000/oauth/callback/google",
		Tokens:       make(map[string]*MockOAuthToken),
	},
	"meta": {
		ProviderName: "meta",
		ClientID:     "mock-meta-client-id",
		ClientSecret: "mock-meta-client-secret",
		RedirectURI:  "http://localhost:3000/oauth/callback/meta",
		Tokens:       make(map[string]*MockOAuthToken),
	},
}

// ============================================================
// MOCK OAUTH ENDPOINTS FOR TESTING
// ============================================================

// mockOAuthAuthorize simulates the authorization endpoint
// Used for: GET /mock/oauth/{platform}/authorize
// Returns: Authorization code and state
func mockOAuthAuthorize(c *gin.Context) {
	platform := c.Param("platform")
	provider, exists := mockProviders[platform]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_platform",
			"message": fmt.Sprintf("Platform '%s' not supported", platform),
		})
		return
	}

	var req MockOAuthRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate request
	if req.ClientID != provider.ClientID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_client_id",
		})
		return
	}

	if req.RedirectURI != provider.RedirectURI {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "redirect_uri_mismatch",
			"expected": provider.RedirectURI,
			"got":      req.RedirectURI,
		})
		return
	}

	// Generate authorization code
	authCode := "auth_" + uuid.New().String()

	// Store state for verification
	provider.Tokens[authCode] = &MockOAuthToken{
		IssuedAt: time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       authCode,
		"state":      req.State,
		"expires_in": 600, // 10 minutes
		"message":    "Authorization successful. Use this code to exchange for tokens.",
	})
}

// mockOAuthToken simulates the token endpoint
// Used for: POST /mock/oauth/{platform}/token
// Returns: Access token, refresh token, and metadata
func mockOAuthToken(c *gin.Context) {
	platform := c.Param("platform")
	provider, exists := mockProviders[platform]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_platform",
		})
		return
	}

	var req MockOAuthTokenRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate credentials
	if req.ClientID != provider.ClientID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_client_id",
		})
		return
	}

	if req.ClientSecret != provider.ClientSecret {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_client_secret",
		})
		return
	}

	// Check if authorization code exists
	_, exists = provider.Tokens[req.Code]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_authorization_code",
		})
		return
	}

	// Generate tokens
	accessToken := "access_" + uuid.New().String()
	refreshToken := "refresh_" + uuid.New().String()

	// Determine account details based on platform
	var accountID, accountName, scope string

	switch platform {
	case "google":
		accountID = "123-456-7890" // Mock Google Customer ID
		accountName = "Mock Google Ads Account"
		scope = "https://www.googleapis.com/auth/adwords"
	case "meta":
		accountID = "act_1234567890" // Mock Meta Ad Account ID
		accountName = "Mock Meta Ads Account"
		scope = "ads_management,business_management"
	}

	token := &MockOAuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
		TokenType:    "Bearer",
		IssuedAt:     time.Now(),
		Scope:        scope,
		AccountID:    accountID,
		AccountName:  accountName,
	}

	// Store token for refresh operations
	provider.Tokens[accessToken] = token

	c.JSON(http.StatusOK, token)
}

// mockOAuthRefresh simulates refreshing an access token
// Used for: POST /mock/oauth/{platform}/refresh
// Returns: New access token
func mockOAuthRefresh(c *gin.Context) {
	platform := c.Param("platform")
	provider, exists := mockProviders[platform]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_platform",
		})
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
		ClientID     string `json:"client_id" binding:"required"`
		ClientSecret string `json:"client_secret" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate credentials
	if req.ClientID != provider.ClientID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_client_id",
		})
		return
	}

	if req.ClientSecret != provider.ClientSecret {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_client_secret",
		})
		return
	}

	// Find token with matching refresh token
	var foundToken *MockOAuthToken
	for _, token := range provider.Tokens {
		if token.RefreshToken == req.RefreshToken {
			foundToken = token
			break
		}
	}

	if foundToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_refresh_token",
		})
		return
	}

	// Generate new access token
	newAccessToken := "access_" + uuid.New().String()

	newToken := &MockOAuthToken{
		AccessToken:  newAccessToken,
		RefreshToken: foundToken.RefreshToken, // Keep same refresh token
		ExpiresIn:    3600,
		TokenType:    "Bearer",
		IssuedAt:     time.Now(),
		Scope:        foundToken.Scope,
		AccountID:    foundToken.AccountID,
		AccountName:  foundToken.AccountName,
	}

	// Store new token
	provider.Tokens[newAccessToken] = newToken

	c.JSON(http.StatusOK, newToken)
}

// mockOAuthValidate validates an access token
// Used for: GET /mock/oauth/{platform}/validate
// Returns: Token info if valid
func mockOAuthValidate(c *gin.Context) {
	platform := c.Param("platform")
	provider, exists := mockProviders[platform]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_platform",
		})
		return
	}

	accessToken := c.GetString("Authorization")
	if accessToken == "" {
		accessToken = c.Query("access_token")
	}

	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing_access_token",
		})
		return
	}

	token, exists := provider.Tokens[accessToken]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid_access_token",
		})
		return
	}

	// Check if token is expired (simplified)
	if time.Since(token.IssuedAt) > time.Duration(token.ExpiresIn)*time.Second {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token_expired",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":        true,
		"account_id":   token.AccountID,
		"account_name": token.AccountName,
		"scope":        token.Scope,
		"issued_at":    token.IssuedAt,
		"expires_at":   token.IssuedAt.Add(time.Duration(token.ExpiresIn) * time.Second),
	})
}

// mockOAuthRevoke simulates revoking a token
// Used for: POST /mock/oauth/{platform}/revoke
// Returns: Success status
func mockOAuthRevoke(c *gin.Context) {
	platform := c.Param("platform")
	provider, exists := mockProviders[platform]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_platform",
		})
		return
	}

	var req struct {
		AccessToken string `json:"access_token" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Remove token
	delete(provider.Tokens, req.AccessToken)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token revoked successfully",
	})
}

// ============================================================
// MOCK OAUTH CALLBACK SIMULATION
// ============================================================

// mockOAuthCallback simulates the callback from OAuth provider
// Used for: GET /mock/oauth/{platform}/callback
// Returns: Redirect to app with code and state
func mockOAuthCallback(c *gin.Context) {
	platform := c.Param("platform")
	code := c.Query("code")
	state := c.Query("state")
	redirectTo := c.Query("redirect_to")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing_code_or_state",
		})
		return
	}

	if redirectTo == "" {
		redirectTo = "http://localhost:3000"
	}

	// In real scenario, this would redirect the browser
	// For API testing, we return the callback data
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"code":     code,
		"state":    state,
		"platform": platform,
		"message":  "In browser, you would be redirected to: " + redirectTo + "?code=" + code + "&state=" + state,
	})
}

// ============================================================
// TEST DATA ENDPOINTS
// ============================================================

// mockOAuthTestData returns predefined test tokens and configs
// Used for: GET /mock/oauth/test-data
// Returns: Test tokens and configuration
func mockOAuthTestData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"purpose": "Use these test credentials for development without real OAuth",
		"google": gin.H{
			"client_id":         mockProviders["google"].ClientID,
			"client_secret":     mockProviders["google"].ClientSecret,
			"redirect_uri":      mockProviders["google"].RedirectURI,
			"test_account_id":   "123-456-7890",
			"test_account_name": "Mock Google Ads Account",
		},
		"meta": gin.H{
			"client_id":         mockProviders["meta"].ClientID,
			"client_secret":     mockProviders["meta"].ClientSecret,
			"redirect_uri":      mockProviders["meta"].RedirectURI,
			"test_account_id":   "act_1234567890",
			"test_account_name": "Mock Meta Ads Account",
		},
		"endpoints": gin.H{
			"authorize": "/mock/oauth/{platform}/authorize",
			"token":     "/mock/oauth/{platform}/token",
			"refresh":   "/mock/oauth/{platform}/refresh",
			"validate":  "/mock/oauth/{platform}/validate",
			"revoke":    "/mock/oauth/{platform}/revoke",
			"callback":  "/mock/oauth/{platform}/callback",
		},
		"documentation": "See /mock/oauth/docs for full documentation",
	})
}

// mockOAuthDocs returns API documentation for mock OAuth
// Used for: GET /mock/oauth/docs
// Returns: Complete API documentation
func mockOAuthDocs(c *gin.Context) {
	docs := gin.H{
		"title":       "Mock OAuth API Documentation",
		"description": "Development OAuth endpoints for testing without real authentication",
		"base_url":    "http://localhost:8080/mock/oauth",
		"endpoints": []gin.H{
			{
				"method":      "POST",
				"path":        "/{platform}/authorize",
				"description": "Start OAuth flow",
				"params": gin.H{
					"platform": "google or meta",
				},
				"request": gin.H{
					"client_id":    "from test-data",
					"redirect_uri": "your app uri",
					"scope":        "permissions needed",
					"state":        "unique state string",
				},
				"response": gin.H{
					"code":  "authorization code",
					"state": "echo of state param",
				},
			},
			{
				"method":      "POST",
				"path":        "/{platform}/token",
				"description": "Exchange code for tokens",
				"request": gin.H{
					"code":          "from authorize response",
					"client_id":     "from test-data",
					"client_secret": "from test-data",
					"redirect_uri":  "same as authorize",
					"grant_type":    "authorization_code",
				},
				"response": gin.H{
					"access_token":  "JWT token",
					"refresh_token": "for refresh",
					"expires_in":    3600,
					"account_id":    "mock account id",
					"account_name":  "mock account name",
				},
			},
			{
				"method":      "POST",
				"path":        "/{platform}/refresh",
				"description": "Get new access token",
				"request": gin.H{
					"refresh_token": "from token response",
					"client_id":     "from test-data",
					"client_secret": "from test-data",
				},
				"response": gin.H{
					"access_token": "new JWT token",
					"expires_in":   3600,
				},
			},
			{
				"method":      "GET",
				"path":        "/{platform}/validate",
				"description": "Validate token",
				"headers": gin.H{
					"Authorization": "Bearer {access_token}",
				},
				"response": gin.H{
					"valid":      true,
					"account_id": "mock account id",
				},
			},
			{
				"method":      "POST",
				"path":        "/{platform}/revoke",
				"description": "Revoke token",
				"request": gin.H{
					"access_token": "token to revoke",
				},
				"response": gin.H{
					"success": true,
				},
			},
		},
		"test_flow": []string{
			"1. GET /mock/oauth/test-data to get credentials",
			"2. POST /mock/oauth/{platform}/authorize with client_id, redirect_uri, scope, state",
			"3. Copy the code from response",
			"4. POST /mock/oauth/{platform}/token with code, client_id, client_secret",
			"5. Use access_token for API calls",
			"6. POST /mock/oauth/{platform}/refresh to get new access_token",
			"7. GET /mock/oauth/{platform}/validate to check token validity",
		},
	}
	c.JSON(http.StatusOK, docs)
}

// ============================================================
// REGISTER MOCK OAUTH ROUTES
// ============================================================

// RegisterMockOAuthRoutes adds mock OAuth routes to the router
func RegisterMockOAuthRoutes(router *gin.Engine) {
	mockOAuth := router.Group("/mock/oauth")
	{
		// Public endpoints
		mockOAuth.GET("/test-data", mockOAuthTestData)
		mockOAuth.GET("/docs", mockOAuthDocs)

		// OAuth flow endpoints
		mockOAuth.POST("/:platform/authorize", mockOAuthAuthorize)
		mockOAuth.POST("/:platform/token", mockOAuthToken)
		mockOAuth.POST("/:platform/refresh", mockOAuthRefresh)
		mockOAuth.GET("/:platform/validate", mockOAuthValidate)
		mockOAuth.POST("/:platform/revoke", mockOAuthRevoke)
		mockOAuth.GET("/:platform/callback", mockOAuthCallback)
	}
}
