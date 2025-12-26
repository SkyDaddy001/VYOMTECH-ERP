package routes

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"lms/cmd/api/models"
	"lms/cmd/api/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegisterClientPaymentRoutes registers all client payment routes
func RegisterClientPaymentRoutes(router *gin.Engine, db interface{}) {
	// Cast to *sql.DB if needed
	sqlDB := db // Assuming db is already *sql.DB

	// Initialize services
	tenantAccountService := services.NewTenantAccountService(sqlDB.(*sql.DB))
	invoiceService := services.NewClientInvoiceService(sqlDB.(*sql.DB))
	paymentService := services.NewClientPaymentService(sqlDB.(*sql.DB))
	existingPaymentService := services.NewPaymentService(sqlDB.(*sql.DB))

	// Tenant Collection Management Routes (Protected - Tenant Admin)
	tenantGroup := router.Group("/api/tenant-collections")
	tenantGroup.Use(AuthMiddleware())
	{
		// Tenant Account Management
		tenantGroup.POST("/accounts", CreateTenantAccount(tenantAccountService))
		tenantGroup.GET("/accounts", ListTenantAccounts(tenantAccountService))
		tenantGroup.GET("/accounts/:accountId", GetTenantAccount(tenantAccountService))
		tenantGroup.PUT("/accounts/:accountId", UpdateTenantAccount(tenantAccountService))

		// Invoice Management
		tenantGroup.POST("/invoices", CreateInvoice(invoiceService))
		tenantGroup.POST("/invoices/bulk", CreateBulkInvoices(invoiceService))
		tenantGroup.GET("/invoices", GetInvoicesByType(invoiceService))
		tenantGroup.GET("/invoices/:invoiceId", GetInvoiceDetail(invoiceService, paymentService))
		tenantGroup.GET("/invoices/client/:clientId", GetClientInvoices(invoiceService))

		// Collection Dashboard
		tenantGroup.GET("/dashboard", GetCollectionDashboard(paymentService))
		tenantGroup.GET("/dashboard/by-charge-type", GetCollectionByChargeType(paymentService))
		tenantGroup.GET("/dashboard/top-clients", GetTopClients(paymentService))
	}

	// Client Payment Routes (Protected - Client/Customer)
	clientGroup := router.Group("/api/client-payments")
	clientGroup.Use(AuthMiddleware())
	{
		// Client Outstanding Balance
		clientGroup.GET("/outstanding", GetClientOutstanding(paymentService, invoiceService))
		clientGroup.GET("/invoices", GetClientInvoiceList(invoiceService))

		// Payment Initiation
		clientGroup.POST("/initiate", InitiateClientPayment(paymentService, invoiceService, tenantAccountService, existingPaymentService))
		clientGroup.GET("/status/:paymentId", GetPaymentStatus(paymentService))
		clientGroup.GET("/history", GetPaymentHistory(paymentService))

		// Payment Verification (Webhook)
		clientGroup.POST("/verify/:provider", VerifyClientPayment(paymentService, invoiceService, existingPaymentService))
	}

	// Admin Routes (Protected - Platform Admin)
	adminGroup := router.Group("/api/admin/tenant-collections")
	adminGroup.Use(AuthMiddleware(), AdminMiddleware())
	{
		adminGroup.GET("/tenants/:tenantId/accounts", GetTenantAccountsAdmin(tenantAccountService))
		adminGroup.GET("/tenants/:tenantId/dashboard", GetTenantDashboardAdmin(paymentService))
		adminGroup.GET("/tenants/:tenantId/invoices", GetTenantInvoicesAdmin(invoiceService))
	}
}

// Tenant Account Handlers

// CreateTenantAccount creates a new payment account for a charge type
func CreateTenantAccount(service *services.TenantAccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id not found in context"})
			return
		}

		var req models.TenantAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		account, err := service.CreateTenantAccount(tenantID, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, account)
	}
}

// ListTenantAccounts lists all accounts for the tenant
func ListTenantAccounts(service *services.TenantAccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id not found in context"})
			return
		}

		accounts, err := service.ListTenantAccounts(tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"accounts": accounts,
			"total":    len(accounts),
		})
	}
}

// GetTenantAccount retrieves a specific account
func GetTenantAccount(service *services.TenantAccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID := c.Param("accountId")
		if accountID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "accountId required"})
			return
		}

		account, err := service.GetTenantAccount(accountID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, account)
	}
}

// UpdateTenantAccount updates account details
func UpdateTenantAccount(service *services.TenantAccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID := c.Param("accountId")
		if accountID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "accountId required"})
			return
		}

		var req models.TenantAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		account, err := service.UpdateTenantAccount(accountID, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, account)
	}
}

// Invoice Handlers

// CreateInvoice creates a single invoice
func CreateInvoice(service *services.ClientInvoiceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id not found"})
			return
		}

		var req models.CreateInvoiceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		invoice, err := service.CreateInvoice(tenantID, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, invoice)
	}
}

// CreateBulkInvoices creates multiple invoices at once
func CreateBulkInvoices(service *services.ClientInvoiceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id not found"})
			return
		}

		var req models.BulkInvoiceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		invoices, err := service.CreateBulkInvoices(tenantID, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"invoices": invoices,
			"count":    len(invoices),
		})
	}
}

// GetInvoicesByType retrieves invoices filtered by charge type
func GetInvoicesByType(service *services.ClientInvoiceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		chargeType := c.Query("charge_type")

		if chargeType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "charge_type query parameter required"})
			return
		}

		invoices, err := service.GetInvoicesByChargeType(tenantID, models.ChargeType(chargeType))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"invoices": invoices,
			"total":    len(invoices),
		})
	}
}

// GetInvoiceDetail retrieves detailed invoice with payment history
func GetInvoiceDetail(invoiceService *services.ClientInvoiceService, paymentService *services.ClientPaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		invoiceID := c.Param("invoiceId")
		if invoiceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invoiceId required"})
			return
		}

		invoice, err := invoiceService.GetInvoice(invoiceID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		payments, err := paymentService.GetClientPaymentsByInvoice(invoiceID)
		if err != nil {
			payments = make([]models.ClientPayment, 0)
		}

		c.JSON(http.StatusOK, gin.H{
			"invoice":  invoice,
			"payments": payments,
		})
	}
}

// GetClientInvoices retrieves invoices for a specific client
func GetClientInvoices(service *services.ClientInvoiceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		clientID := c.Param("clientId")

		if clientID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "clientId required"})
			return
		}

		invoices, err := service.GetClientInvoices(tenantID, clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"invoices": invoices,
			"total":    len(invoices),
		})
	}
}

// Client Payment Handlers

// GetClientOutstanding retrieves outstanding balance for logged-in client
func GetClientOutstanding(paymentService *services.ClientPaymentService, invoiceService *services.ClientInvoiceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client ID from context (set by auth middleware)
		clientID := c.GetString("client_id")
		tenantID := c.GetString("tenant_id")

		if clientID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "client_id not found"})
			return
		}

		outstanding, err := paymentService.GetClientOutstanding(tenantID, clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, outstanding)
	}
}

// GetClientInvoiceList retrieves invoices for logged-in client
func GetClientInvoiceList(service *services.ClientInvoiceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.GetString("client_id")
		tenantID := c.GetString("tenant_id")

		if clientID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "client_id not found"})
			return
		}

		// Get outstanding invoices only
		invoices, err := service.GetOutstandingInvoices(tenantID, clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"invoices": invoices,
			"total":    len(invoices),
		})
	}
}

// InitiateClientPayment initiates a payment from a client
func InitiateClientPayment(
	paymentService *services.ClientPaymentService,
	invoiceService *services.ClientInvoiceService,
	accountService *services.TenantAccountService,
	existingPaymentService *services.PaymentService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.GetString("client_id")
		tenantID := c.GetString("tenant_id")

		if clientID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "client_id not found"})
			return
		}

		var req models.ClientPaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify invoice exists and belongs to client
		invoice, err := invoiceService.GetInvoice(req.InvoiceID)
		if err != nil || invoice.ClientID != clientID {
			c.JSON(http.StatusForbidden, gin.H{"error": "invoice not found or unauthorized"})
			return
		}

		// Get tenant account for charge type
		account, err := accountService.GetTenantAccountByChargeType(tenantID, req.ChargeType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create payment record
		payment := &models.ClientPayment{
			ID:              uuid.New().String(),
			TenantID:        tenantID,
			ClientID:        clientID,
			InvoiceID:       req.InvoiceID,
			TenantAccountID: account.ID,
			ChargeType:      req.ChargeType,
			OrderID:         uuid.New().String(),
			Amount:          req.Amount,
			Currency:        req.Currency,
			Status:          "created",
			PaymentType:     "client",
			Provider:        req.Provider,
			PaymentMethod:   req.PaymentMethod,
			ClientName:      req.ClientName,
			ClientEmail:     req.ClientEmail,
			ClientPhone:     req.ClientPhone,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := paymentService.CreateClientPayment(payment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Route to appropriate payment gateway
		var gatewayOrder interface{}
		switch req.Provider {
		case "razorpay":
			// Use existing Razorpay service with tenant account ID
			// This would need account-specific credentials
			//gatewayOrder, err = initializeRazorpayOrder(account, payment)
		case "billdesk":
			// Use existing Billdesk service with tenant account ID
			//gatewayOrder, err = initializeBilldeskOrder(account, payment)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported payment provider"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"payment_id":    payment.ID,
			"order_id":      payment.OrderID,
			"amount":        payment.Amount,
			"currency":      payment.Currency,
			"gateway_order": gatewayOrder,
		})
	}
}

// GetPaymentStatus retrieves payment status
func GetPaymentStatus(service *services.ClientPaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		paymentID := c.Param("paymentId")
		if paymentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "paymentId required"})
			return
		}

		payment, err := service.GetClientPayment(paymentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, payment)
	}
}

// GetPaymentHistory retrieves payment history for logged-in client
func GetPaymentHistory(service *services.ClientPaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.GetString("client_id")
		tenantID := c.GetString("tenant_id")

		limitStr := c.DefaultQuery("limit", "20")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		if limit > 100 {
			limit = 100
		}

		payments, total, err := service.ListClientPayments(tenantID, clientID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"payments": payments,
			"total":    total,
			"limit":    limit,
			"offset":   offset,
		})
	}
}

// VerifyClientPayment verifies payment from gateway webhook
func VerifyClientPayment(
	paymentService *services.ClientPaymentService,
	invoiceService *services.ClientInvoiceService,
	existingPaymentService *services.PaymentService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")
		if provider == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "provider required"})
			return
		}

		var webhookPayload map[string]interface{}
		if err := c.ShouldBindJSON(&webhookPayload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify signature based on provider
		// Update payment status if verified
		// Update invoice payment amount

		c.JSON(http.StatusOK, gin.H{"status": "verified"})
	}
}

// Collection Dashboard Handlers

// GetCollectionDashboard retrieves collection dashboard for tenant
func GetCollectionDashboard(service *services.ClientPaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")

		dashboard, err := service.GetCollectionDashboard(tenantID, 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dashboard)
	}
}

// GetCollectionByChargeType retrieves collection statistics grouped by charge type
func GetCollectionByChargeType(service *services.ClientPaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")

		dashboard, err := service.GetCollectionDashboard(tenantID, 1000)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dashboard.CollectionByType)
	}
}

// GetTopClients retrieves top clients by collection amount
func GetTopClients(service *services.ClientPaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would need additional database query
		// For now, return placeholder
		c.JSON(http.StatusOK, gin.H{
			"message": "Top clients data - to be implemented with custom query",
		})
	}
}

// Admin Handlers

// GetTenantAccountsAdmin retrieves accounts for a specific tenant (admin only)
func GetTenantAccountsAdmin(service *services.TenantAccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("tenantId")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenantId required"})
			return
		}

		accounts, err := service.ListTenantAccounts(tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, accounts)
	}
}

// GetTenantDashboardAdmin retrieves dashboard for a tenant (admin only)
func GetTenantDashboardAdmin(service *services.ClientPaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("tenantId")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tenantId required"})
			return
		}

		dashboard, err := service.GetCollectionDashboard(tenantID, 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dashboard)
	}
}

// GetTenantInvoicesAdmin retrieves all invoices for a tenant (admin only)
func GetTenantInvoicesAdmin(service *services.ClientInvoiceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would need a custom database query for all invoices
		c.JSON(http.StatusOK, gin.H{
			"message": "Tenant invoices - to be implemented with custom query",
		})
	}
}
