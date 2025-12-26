package credentials

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CredentialHandler handles HTTP requests for credential management
type CredentialHandler struct {
	service *CredentialService
}

// NewCredentialHandler creates a new credential handler
func NewCredentialHandler(service *CredentialService) *CredentialHandler {
	return &CredentialHandler{
		service: service,
	}
}

// StoreGoogleOAuthRequest request body for storing Google OAuth credential
type StoreGoogleOAuthRequest struct {
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	RedirectURI  string `json:"redirect_uri" binding:"required"`
	Description  string `json:"description"`
}

// StoreMetaOAuthRequest request body for storing Meta OAuth credential
type StoreMetaOAuthRequest struct {
	AppID       string `json:"app_id" binding:"required"`
	AppSecret   string `json:"app_secret" binding:"required"`
	RedirectURI string `json:"redirect_uri" binding:"required"`
	Description string `json:"description"`
}

// StoreEmailSMTPRequest request body for storing Email SMTP credential
type StoreEmailSMTPRequest struct {
	Host        string `json:"host" binding:"required"`
	Port        int    `json:"port" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	FromName    string `json:"from_name" binding:"required"`
	FromEmail   string `json:"from_email" binding:"required,email"`
	TLS         bool   `json:"tls"`
	Description string `json:"description"`
}

// StoreAWSS3Request request body for storing AWS S3 credential
type StoreAWSS3Request struct {
	AccessKeyID     string `json:"access_key_id" binding:"required"`
	SecretAccessKey string `json:"secret_access_key" binding:"required"`
	Region          string `json:"region" binding:"required"`
	Bucket          string `json:"bucket" binding:"required"`
	BucketPrefix    string `json:"bucket_prefix"`
	Description     string `json:"description"`
}

// StoreRazorpayRequest request body for storing Razorpay credential
type StoreRazorpayRequest struct {
	KeyID       string `json:"key_id" binding:"required"`
	KeySecret   string `json:"key_secret" binding:"required"`
	Description string `json:"description"`
}

// StoreBilldeskRequest request body for storing Billdesk credential
type StoreBilldeskRequest struct {
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	MerchantID   string `json:"merchant_id" binding:"required"`
	Sandbox      bool   `json:"sandbox"`
	Description  string `json:"description"`
}

// ErrorResponse standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse standard success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// StoreGoogleOAuth stores Google OAuth credentials for a tenant
// POST /api/v1/credentials/google-oauth
func (h *CredentialHandler) StoreGoogleOAuth(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "tenant_id not found in context",
		})
		return
	}

	var req StoreGoogleOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cred := GoogleOAuthCredential{
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		RedirectURI:  req.RedirectURI,
	}

	err := h.service.StoreCredential(
		c.Request.Context(),
		tenantID,
		CredentialTypeGoogleOAuth,
		cred,
		req.Description,
		nil,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "STORE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Google OAuth credential stored successfully",
	})
}

// StoreMetaOAuth stores Meta OAuth credentials for a tenant
// POST /api/v1/credentials/meta-oauth
func (h *CredentialHandler) StoreMetaOAuth(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "tenant_id not found in context",
		})
		return
	}

	var req StoreMetaOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cred := MetaOAuthCredential{
		AppID:       req.AppID,
		AppSecret:   req.AppSecret,
		RedirectURI: req.RedirectURI,
	}

	err := h.service.StoreCredential(
		c.Request.Context(),
		tenantID,
		CredentialTypeMetaOAuth,
		cred,
		req.Description,
		nil,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "STORE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Meta OAuth credential stored successfully",
	})
}

// StoreEmailSMTP stores Email SMTP credentials for a tenant
// POST /api/v1/credentials/email-smtp
func (h *CredentialHandler) StoreEmailSMTP(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "tenant_id not found in context",
		})
		return
	}

	var req StoreEmailSMTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cred := EmailSMTPCredential{
		Host:      req.Host,
		Port:      req.Port,
		Username:  req.Username,
		Password:  req.Password,
		FromName:  req.FromName,
		FromEmail: req.FromEmail,
		TLS:       req.TLS,
	}

	err := h.service.StoreCredential(
		c.Request.Context(),
		tenantID,
		CredentialTypeEmailSMTP,
		cred,
		req.Description,
		nil,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "STORE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Email SMTP credential stored successfully",
	})
}

// StoreAWSS3 stores AWS S3 credentials for a tenant
// POST /api/v1/credentials/aws-s3
func (h *CredentialHandler) StoreAWSS3(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "tenant_id not found in context",
		})
		return
	}

	var req StoreAWSS3Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cred := AWSS3Credential{
		AccessKeyID:     req.AccessKeyID,
		SecretAccessKey: req.SecretAccessKey,
		Region:          req.Region,
		Bucket:          req.Bucket,
		BucketPrefix:    req.BucketPrefix,
	}

	err := h.service.StoreCredential(
		c.Request.Context(),
		tenantID,
		CredentialTypeAWSS3,
		cred,
		req.Description,
		nil,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "STORE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "AWS S3 credential stored successfully",
	})
}

// StoreRazorpay stores Razorpay credentials for a tenant
// POST /api/v1/credentials/razorpay
func (h *CredentialHandler) StoreRazorpay(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "tenant_id not found in context",
		})
		return
	}

	var req StoreRazorpayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cred := RazorpayCredential{
		KeyID:     req.KeyID,
		KeySecret: req.KeySecret,
	}

	err := h.service.StoreCredential(
		c.Request.Context(),
		tenantID,
		CredentialTypeRazorpay,
		cred,
		req.Description,
		nil,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "STORE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Razorpay credential stored successfully",
	})
}

// StoreBilldesk stores Billdesk credentials for a tenant
// POST /api/v1/credentials/billdesk
func (h *CredentialHandler) StoreBilldesk(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "tenant_id not found in context",
		})
		return
	}

	var req StoreBilldeskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cred := BilldeskCredential{
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		MerchantID:   req.MerchantID,
		Sandbox:      req.Sandbox,
	}

	err := h.service.StoreCredential(
		c.Request.Context(),
		tenantID,
		CredentialTypeBilldesk,
		cred,
		req.Description,
		nil,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "STORE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Billdesk credential stored successfully",
	})
}

// ListCredentials lists all credentials for a tenant
// GET /api/v1/credentials
func (h *CredentialHandler) ListCredentials(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "tenant_id not found in context",
		})
		return
	}

	creds, err := h.service.ListCredentials(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "LIST_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Credentials retrieved successfully",
		Data:    creds,
	})
}

// DeleteCredential deletes (deactivates) a credential
// DELETE /api/v1/credentials/:credentialType
func (h *CredentialHandler) DeleteCredential(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "tenant_id not found in context",
		})
		return
	}

	credType := c.Param("credentialType")
	if credType == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "credentialType is required",
		})
		return
	}

	err := h.service.DeleteCredential(c.Request.Context(), tenantID, CredentialType(credType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "DELETE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Credential deleted successfully",
	})
}

// RotateCredential rotates a credential to a new value
// POST /api/v1/credentials/:credentialType/rotate
func (h *CredentialHandler) RotateCredential(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "tenant_id not found in context",
		})
		return
	}

	credType := c.Param("credentialType")
	if credType == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "credentialType is required",
		})
		return
	}

	// Parse the request body based on credential type
	var credentialData interface{}
	var err error

	switch CredentialType(credType) {
	case CredentialTypeGoogleOAuth:
		var req StoreGoogleOAuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "INVALID_REQUEST",
				Message: err.Error(),
			})
			return
		}
		credentialData = GoogleOAuthCredential{
			ClientID:     req.ClientID,
			ClientSecret: req.ClientSecret,
			RedirectURI:  req.RedirectURI,
		}

	case CredentialTypeMetaOAuth:
		var req StoreMetaOAuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "INVALID_REQUEST",
				Message: err.Error(),
			})
			return
		}
		credentialData = MetaOAuthCredential{
			AppID:       req.AppID,
			AppSecret:   req.AppSecret,
			RedirectURI: req.RedirectURI,
		}

	case CredentialTypeRazorpay:
		var req StoreRazorpayRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "INVALID_REQUEST",
				Message: err.Error(),
			})
			return
		}
		credentialData = RazorpayCredential{
			KeyID:     req.KeyID,
			KeySecret: req.KeySecret,
		}

	default:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "UNSUPPORTED_TYPE",
			Message: "Credential type not supported for rotation",
		})
		return
	}

	err = h.service.RotateCredential(c.Request.Context(), tenantID, CredentialType(credType), credentialData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "ROTATE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Credential rotated successfully",
	})
}
