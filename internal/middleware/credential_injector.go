package middleware

import (
	"fmt"

	"lms/internal/credentials"

	"github.com/gin-gonic/gin"
)

// CredentialMiddleware middleware to inject tenant credentials into context
type CredentialMiddleware struct {
	credentialService *credentials.CredentialService
}

// NewCredentialMiddleware creates a new credential middleware
func NewCredentialMiddleware(service *credentials.CredentialService) *CredentialMiddleware {
	return &CredentialMiddleware{
		credentialService: service,
	}
}

// InjectGoogleOAuth injects Google OAuth credentials into context for the tenant
func (cm *CredentialMiddleware) InjectGoogleOAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.Next()
			return
		}

		var googleCred credentials.GoogleOAuthCredential
		err := cm.credentialService.GetCredential(
			c.Request.Context(),
			tenantID,
			credentials.CredentialTypeGoogleOAuth,
			&googleCred,
		)

		if err != nil {
			// Log error but don't block the request
			fmt.Printf("Failed to load Google OAuth credential for tenant %s: %v\n", tenantID, err)
			c.Next()
			return
		}

		// Inject into context
		c.Set("google_oauth_credential", googleCred)
		c.Next()
	}
}

// InjectMetaOAuth injects Meta OAuth credentials into context for the tenant
func (cm *CredentialMiddleware) InjectMetaOAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.Next()
			return
		}

		var metaCred credentials.MetaOAuthCredential
		err := cm.credentialService.GetCredential(
			c.Request.Context(),
			tenantID,
			credentials.CredentialTypeMetaOAuth,
			&metaCred,
		)

		if err != nil {
			// Log error but don't block the request
			fmt.Printf("Failed to load Meta OAuth credential for tenant %s: %v\n", tenantID, err)
			c.Next()
			return
		}

		// Inject into context
		c.Set("meta_oauth_credential", metaCred)
		c.Next()
	}
}

// InjectEmailSMTP injects Email SMTP credentials into context for the tenant
func (cm *CredentialMiddleware) InjectEmailSMTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.Next()
			return
		}

		var emailCred credentials.EmailSMTPCredential
		err := cm.credentialService.GetCredential(
			c.Request.Context(),
			tenantID,
			credentials.CredentialTypeEmailSMTP,
			&emailCred,
		)

		if err != nil {
			// Log error but don't block the request
			fmt.Printf("Failed to load Email SMTP credential for tenant %s: %v\n", tenantID, err)
			c.Next()
			return
		}

		// Inject into context
		c.Set("email_smtp_credential", emailCred)
		c.Next()
	}
}

// InjectAWSS3 injects AWS S3 credentials into context for the tenant
func (cm *CredentialMiddleware) InjectAWSS3() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.Next()
			return
		}

		var awsCred credentials.AWSS3Credential
		err := cm.credentialService.GetCredential(
			c.Request.Context(),
			tenantID,
			credentials.CredentialTypeAWSS3,
			&awsCred,
		)

		if err != nil {
			// Log error but don't block the request
			fmt.Printf("Failed to load AWS S3 credential for tenant %s: %v\n", tenantID, err)
			c.Next()
			return
		}

		// Inject into context
		c.Set("aws_s3_credential", awsCred)
		c.Next()
	}
}

// InjectRazorpay injects Razorpay credentials into context for the tenant
func (cm *CredentialMiddleware) InjectRazorpay() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.Next()
			return
		}

		var razorpayCred credentials.RazorpayCredential
		err := cm.credentialService.GetCredential(
			c.Request.Context(),
			tenantID,
			credentials.CredentialTypeRazorpay,
			&razorpayCred,
		)

		if err != nil {
			// Log error but don't block the request
			fmt.Printf("Failed to load Razorpay credential for tenant %s: %v\n", tenantID, err)
			c.Next()
			return
		}

		// Inject into context
		c.Set("razorpay_credential", razorpayCred)
		c.Next()
	}
}

// InjectAllCredentials injects all available credentials for the tenant
func (cm *CredentialMiddleware) InjectAllCredentials() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.Next()
			return
		}

		credentialTypes := []credentials.CredentialType{
			credentials.CredentialTypeGoogleOAuth,
			credentials.CredentialTypeMetaOAuth,
			credentials.CredentialTypeEmailSMTP,
			credentials.CredentialTypeAWSS3,
			credentials.CredentialTypeRazorpay,
			credentials.CredentialTypeBilldesk,
			credentials.CredentialTypeGoogleAds,
			credentials.CredentialTypeMetaAds,
		}

		for _, credType := range credentialTypes {
			// Create a generic placeholder based on type
			switch credType {
			case credentials.CredentialTypeGoogleOAuth:
				var cred credentials.GoogleOAuthCredential
				_ = cm.credentialService.GetCredential(c.Request.Context(), tenantID, credType, &cred)
				c.Set("google_oauth_credential", cred)

			case credentials.CredentialTypeMetaOAuth:
				var cred credentials.MetaOAuthCredential
				_ = cm.credentialService.GetCredential(c.Request.Context(), tenantID, credType, &cred)
				c.Set("meta_oauth_credential", cred)

			case credentials.CredentialTypeEmailSMTP:
				var cred credentials.EmailSMTPCredential
				_ = cm.credentialService.GetCredential(c.Request.Context(), tenantID, credType, &cred)
				c.Set("email_smtp_credential", cred)

			case credentials.CredentialTypeAWSS3:
				var cred credentials.AWSS3Credential
				_ = cm.credentialService.GetCredential(c.Request.Context(), tenantID, credType, &cred)
				c.Set("aws_s3_credential", cred)

			case credentials.CredentialTypeRazorpay:
				var cred credentials.RazorpayCredential
				_ = cm.credentialService.GetCredential(c.Request.Context(), tenantID, credType, &cred)
				c.Set("razorpay_credential", cred)

			case credentials.CredentialTypeBilldesk:
				var cred credentials.BilldeskCredential
				_ = cm.credentialService.GetCredential(c.Request.Context(), tenantID, credType, &cred)
				c.Set("billdesk_credential", cred)

			case credentials.CredentialTypeGoogleAds:
				var cred credentials.GoogleAdsCredential
				_ = cm.credentialService.GetCredential(c.Request.Context(), tenantID, credType, &cred)
				c.Set("google_ads_credential", cred)

			case credentials.CredentialTypeMetaAds:
				var cred credentials.MetaAdsCredential
				_ = cm.credentialService.GetCredential(c.Request.Context(), tenantID, credType, &cred)
				c.Set("meta_ads_credential", cred)
			}
		}

		c.Next()
	}
}
