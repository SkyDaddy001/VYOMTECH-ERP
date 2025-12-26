package routes

import (
	"database/sql"
	"net/http"
	"time"

	"lms/cmd/api/models"
	"lms/cmd/api/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentRoutes sets up payment-related routes
func SetupPaymentRoutes(router *gin.Engine, db *sql.DB, paymentService *services.PaymentService) {
	paymentGroup := router.Group("/api/v1/payments")
	paymentGroup.Use(AuthMiddleware())

	// Payment endpoints
	paymentGroup.POST("/initiate", initiatePayment(paymentService))
	paymentGroup.GET("/:payment_id", getPaymentStatus(paymentService))
	paymentGroup.POST("/:payment_id/verify", verifyPayment(paymentService))
	paymentGroup.POST("/:payment_id/refund", refundPayment(paymentService))
	paymentGroup.GET("/list", listPayments(paymentService))
	paymentGroup.GET("/summary", getPaymentSummary(paymentService))

	// Payment methods endpoints
	paymentGroup.GET("/methods/banks", getBanks())
	paymentGroup.GET("/methods/available", getAvailablePaymentMethods())

	// Webhook endpoints (no auth required)
	webhookGroup := router.Group("/api/v1/webhooks")
	webhookGroup.POST("/razorpay", handleRazorpayWebhook(paymentService))
	webhookGroup.POST("/billdesk", handleBilldeskWebhook(paymentService))

	// Payment configuration endpoints
	configGroup := router.Group("/api/v1/payment-config")
	configGroup.Use(AuthMiddleware(), AdminMiddleware())

	configGroup.GET("/gateways", getPaymentGateways(paymentService))
	configGroup.POST("/gateways/configure", configurePaymentGateway(paymentService))
	configGroup.PUT("/gateways/:gateway_id", updatePaymentGateway(paymentService))
	configGroup.DELETE("/gateways/:gateway_id", deletePaymentGateway(paymentService))
}

// initiatePayment initiates a payment request
func initiatePayment(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
			return
		}

		var req models.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create payment record
		payment := &models.Payment{
			ID:             uuid.New().String(),
			TenantID:       tenantID,
			OrderID:        uuid.New().String(),
			Amount:         req.Amount,
			Currency:       req.Currency,
			Status:         models.PaymentPending,
			Provider:       req.Provider,
			PaymentMethod:  req.PaymentMethod,
			Description:    req.Description,
			CustomerName:   req.CustomerName,
			CustomerEmail:  req.CustomerEmail,
			CustomerPhone:  req.CustomerPhone,
			BillingAddress: req.BillingAddress,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Store payment in database
		if err := paymentService.CreatePayment(payment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
			return
		}

		// Initialize payment with gateway
		paymentURL, gatewayOrderID, expiresAt, err := paymentService.InitiatePayment(payment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update payment with gateway details
		payment.Status = models.PaymentInitiated
		payment.GatewayOrderID = gatewayOrderID
		payment.ExpiresAt = expiresAt
		paymentService.UpdatePayment(payment)

		response := models.PaymentResponse{
			PaymentID:      payment.ID,
			OrderID:        payment.OrderID,
			Amount:         payment.Amount,
			Currency:       payment.Currency,
			Status:         payment.Status,
			Provider:       payment.Provider,
			PaymentMethod:  payment.PaymentMethod,
			GatewayOrderID: gatewayOrderID,
			PaymentURL:     paymentURL,
			ExpiresAt:      expiresAt,
			CreatedAt:      payment.CreatedAt,
		}

		c.JSON(http.StatusCreated, response)
	}
}

// getPaymentStatus retrieves payment status
func getPaymentStatus(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		paymentID := c.Param("payment_id")
		tenantID := c.GetString("tenant_id")

		payment, err := paymentService.GetPayment(paymentID, tenantID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}

		c.JSON(http.StatusOK, payment)
	}
}

// verifyPayment verifies payment
func verifyPayment(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		paymentID := c.Param("payment_id")
		tenantID := c.GetString("tenant_id")

		var req struct {
			PaymentGatewayID string `json:"payment_gateway_id" binding:"required"`
			Signature        string `json:"signature" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		payment, err := paymentService.GetPayment(paymentID, tenantID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}

		// Verify payment with gateway
		isValid, err := paymentService.VerifyPayment(payment, req.PaymentGatewayID, req.Signature)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !isValid {
			payment.Status = models.PaymentFailed
			payment.ErrorMessage = "Payment verification failed"
			paymentService.UpdatePayment(payment)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Payment verification failed"})
			return
		}

		payment.Status = models.PaymentSuccessful
		now := time.Now()
		payment.ProcessedAt = &now
		paymentService.UpdatePayment(payment)

		c.JSON(http.StatusOK, gin.H{
			"status":  "verified",
			"payment": payment,
		})
	}
}

// refundPayment initiates a refund
func refundPayment(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		paymentID := c.Param("payment_id")
		tenantID := c.GetString("tenant_id")

		var req models.RefundRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		refund, err := paymentService.CreateRefund(paymentID, tenantID, req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, refund)
	}
}

// listPayments lists all payments for tenant
func listPayments(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		status := c.Query("status")
		provider := c.Query("provider")
		limit := 50
		offset := 0

		payments, total, err := paymentService.ListPayments(tenantID, status, provider, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"payments": payments,
			"total":    total,
		})
	}
}

// getPaymentSummary gets payment statistics
func getPaymentSummary(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")

		summary, err := paymentService.GetPaymentSummary(tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summary"})
			return
		}

		c.JSON(http.StatusOK, summary)
	}
}

// getBanks returns list of supported banks
func getBanks() gin.HandlerFunc {
	return func(c *gin.Context) {
		banks := models.AvailableBanks()
		c.JSON(http.StatusOK, gin.H{"banks": banks})
	}
}

// getAvailablePaymentMethods returns available payment methods
func getAvailablePaymentMethods() gin.HandlerFunc {
	return func(c *gin.Context) {
		methods := []map[string]interface{}{
			{
				"name":        "Netbanking",
				"code":        "netbanking",
				"description": "Pay using your bank account",
				"icon":        "bank",
			},
			{
				"name":        "Credit Card",
				"code":        "credit_card",
				"description": "Pay using credit card",
				"icon":        "credit-card",
			},
			{
				"name":        "Debit Card",
				"code":        "debit_card",
				"description": "Pay using debit card",
				"icon":        "credit-card",
			},
			{
				"name":        "UPI",
				"code":        "upi",
				"description": "Unified Payments Interface",
				"icon":        "mobile-phone",
			},
			{
				"name":        "Digital Wallet",
				"code":        "wallet",
				"description": "Use digital wallets like Google Pay, Apple Pay",
				"icon":        "wallet",
			},
		}

		c.JSON(http.StatusOK, gin.H{"payment_methods": methods})
	}
}

// handleRazorpayWebhook handles Razorpay webhooks
func handleRazorpayWebhook(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = c.GetHeader("X-Razorpay-Signature")

		var webhook map[string]interface{}
		if err := c.ShouldBindJSON(&webhook); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "received"})
	}
}

// handleBilldeskWebhook handles Billdesk webhooks
func handleBilldeskWebhook(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = c.GetHeader("X-Billdesk-Signature")

		var webhook map[string]interface{}
		if err := c.ShouldBindJSON(&webhook); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, _ = c.GetRawData()

		c.JSON(http.StatusOK, gin.H{"status": "received"})
	}
}

// getPaymentGateways retrieves payment gateway configurations
func getPaymentGateways(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")

		gateways, err := paymentService.GetPaymentGateways(tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch gateways"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"gateways": gateways})
	}
}

// configurePaymentGateway configures a payment gateway
func configurePaymentGateway(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")

		var req struct {
			Provider  string                 `json:"provider" binding:"required,oneof=razorpay billdesk"`
			ApiKey    string                 `json:"api_key" binding:"required"`
			ApiSecret string                 `json:"api_secret" binding:"required"`
			Settings  map[string]interface{} `json:"settings"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		config, err := paymentService.ConfigureGateway(tenantID, req.Provider, req.ApiKey, req.ApiSecret, req.Settings)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, config)
	}
}

// updatePaymentGateway updates payment gateway configuration
func updatePaymentGateway(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		gatewayID := c.Param("gateway_id")
		tenantID := c.GetString("tenant_id")

		var req struct {
			ApiKey    string                 `json:"api_key"`
			ApiSecret string                 `json:"api_secret"`
			IsActive  *bool                  `json:"is_active"`
			Settings  map[string]interface{} `json:"settings"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		config, err := paymentService.UpdateGateway(gatewayID, tenantID, req.ApiKey, req.ApiSecret, req.IsActive, req.Settings)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, config)
	}
}

// deletePaymentGateway deletes payment gateway configuration
func deletePaymentGateway(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		gatewayID := c.Param("gateway_id")
		tenantID := c.GetString("tenant_id")

		if err := paymentService.DeleteGateway(gatewayID, tenantID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}
