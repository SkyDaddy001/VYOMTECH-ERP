package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"vyom-erp/models"
	"vyom-erp/payment"

	"github.com/google/uuid"
)

// PaymentService handles payment operations
type PaymentService struct {
	db              *sql.DB
	razorpayClients map[string]*payment.RazorpayClient
	billdeskClients map[string]*payment.BilldeskClient
}

// NewPaymentService creates a new payment service
func NewPaymentService(db *sql.DB) *PaymentService {
	return &PaymentService{
		db:              db,
		razorpayClients: make(map[string]*payment.RazorpayClient),
		billdeskClients: make(map[string]*payment.BilldeskClient),
	}
}

// CreatePayment creates a new payment record
func (ps *PaymentService) CreatePayment(p *models.Payment) error {
	query := `
		INSERT INTO payments (
			id, tenant_id, order_id, amount, currency, status, provider, payment_method,
			description, customer_name, customer_email, customer_phone, billing_address,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`

	billingAddr, _ := json.Marshal(p.BillingAddress)

	_, err := ps.db.Exec(query,
		p.ID, p.TenantID, p.OrderID, p.Amount, p.Currency, p.Status, p.Provider,
		p.PaymentMethod, p.Description, p.CustomerName, p.CustomerEmail, p.CustomerPhone,
		billingAddr, p.CreatedAt, p.UpdatedAt,
	)

	return err
}

// GetPayment retrieves a payment by ID
func (ps *PaymentService) GetPayment(paymentID, tenantID string) (*models.Payment, error) {
	query := `
		SELECT id, tenant_id, order_id, amount, currency, status, provider, payment_method,
		       transaction_id, gateway_order_id, gateway_payment_id, description, customer_name,
		       customer_email, customer_phone, billing_address, error_message, error_code,
		       receipt_url, refund_id, refund_amount, created_at, updated_at, processed_at, expires_at
		FROM payments
		WHERE id = $1 AND tenant_id = $2
	`

	p := &models.Payment{}
	var billingAddr json.RawMessage

	err := ps.db.QueryRow(query, paymentID, tenantID).Scan(
		&p.ID, &p.TenantID, &p.OrderID, &p.Amount, &p.Currency, &p.Status, &p.Provider,
		&p.PaymentMethod, &p.TransactionID, &p.GatewayOrderID, &p.GatewayPaymentID, &p.Description,
		&p.CustomerName, &p.CustomerEmail, &p.CustomerPhone, &billingAddr, &p.ErrorMessage,
		&p.ErrorCode, &p.ReceiptURL, &p.RefundID, &p.RefundAmount, &p.CreatedAt, &p.UpdatedAt,
		&p.ProcessedAt, &p.ExpiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	if billingAddr != nil {
		json.Unmarshal(billingAddr, &p.BillingAddress)
	}

	return p, nil
}

// UpdatePayment updates a payment record
func (ps *PaymentService) UpdatePayment(p *models.Payment) error {
	query := `
		UPDATE payments
		SET status = $1, transaction_id = $2, gateway_order_id = $3, gateway_payment_id = $4,
		    error_message = $5, error_code = $6, receipt_url = $7, refund_id = $8,
		    refund_amount = $9, processed_at = $10, expires_at = $11, updated_at = $12
		WHERE id = $13
	`

	_, err := ps.db.Exec(query,
		p.Status, p.TransactionID, p.GatewayOrderID, p.GatewayPaymentID,
		p.ErrorMessage, p.ErrorCode, p.ReceiptURL, p.RefundID,
		p.RefundAmount, p.ProcessedAt, p.ExpiresAt, time.Now(), p.ID,
	)

	return err
}

// InitiatePayment initiates a payment with the gateway
func (ps *PaymentService) InitiatePayment(p *models.Payment) (string, string, *time.Time, error) {
	// Get gateway configuration
	config, err := ps.GetActiveGatewayConfig(p.TenantID, string(p.Provider))
	if err != nil {
		return "", "", nil, err
	}

	switch p.Provider {
	case models.ProviderRazorpay:
		return ps.initiateRazorpayPayment(p, config)
	case models.ProviderBilldesk:
		return ps.initiateBilldeskPayment(p, config)
	default:
		return "", "", nil, fmt.Errorf("unsupported provider: %s", p.Provider)
	}
}

// initiateRazorpayPayment initiates Razorpay payment
func (ps *PaymentService) initiateRazorpayPayment(p *models.Payment, config *models.PaymentGatewayConfig) (string, string, *time.Time, error) {
	// Get or create Razorpay client
	client, exists := ps.razorpayClients[config.ID]
	if !exists {
		var c *payment.RazorpayClient
		var err error
		c, err = payment.NewRazorpayClient(config.ApiKey, config.ApiSecret)
		if err != nil {
			return "", "", nil, err
		}
		client = c
		ps.razorpayClients[config.ID] = client
	}

	// Create order
	order, err := client.CreateOrder(p.Amount, p.Currency, p.OrderID, p.Description)
	if err != nil {
		return "", "", nil, err
	}

	// Generate payment URL
	paymentURL := fmt.Sprintf("https://checkout.razorpay.com/?key_id=%s&order_id=%s&customer_name=%s&customer_email=%s&customer_contact=%s",
		config.ApiKey, order.ID, p.CustomerName, p.CustomerEmail, p.CustomerPhone)

	expiresAt := time.Now().Add(24 * time.Hour)

	return paymentURL, order.ID, &expiresAt, nil
}

// initiateBilldeskPayment initiates Billdesk payment
func (ps *PaymentService) initiateBilldeskPayment(p *models.Payment, config *models.PaymentGatewayConfig) (string, string, *time.Time, error) {
	// Parse settings
	var settings map[string]interface{}
	if err := json.Unmarshal(config.Settings, &settings); err != nil {
		return "", "", nil, err
	}

	merchantID := config.ApiKey
	clientID := ""
	environment := "sandbox"

	if v, ok := settings["client_id"].(string); ok {
		clientID = v
	}
	if v, ok := settings["environment"].(string); ok {
		environment = v
	}

	// Get or create Billdesk client
	key := config.ID
	client, exists := ps.billdeskClients[key]
	if !exists {
		var c *payment.BilldeskClient
		var err error
		c, err = payment.NewBilldeskClient(config.ApiKey, config.ApiSecret, merchantID, clientID, environment)
		if err != nil {
			return "", "", nil, err
		}
		client = c
		ps.billdeskClients[key] = client
	}

	// Map payment method
	paymentMethods := []string{string(p.PaymentMethod)}

	// Create order
	order, err := client.CreateOrder(p.Amount, p.Currency, p.OrderID, p.Description, paymentMethods)
	if err != nil {
		return "", "", nil, err
	}

	return order.PaymentURL, order.BdOrderID, order.ExpiresAt, nil
}

// VerifyPayment verifies a payment
func (ps *PaymentService) VerifyPayment(p *models.Payment, gatewayID, signature string) (bool, error) {
	config, err := ps.GetGatewayConfig(gatewayID)
	if err != nil {
		return false, err
	}

	switch p.Provider {
	case models.ProviderRazorpay:
		client, exists := ps.razorpayClients[config.ID]
		if !exists {
			var c *payment.RazorpayClient
			c, err = payment.NewRazorpayClient(config.ApiKey, config.ApiSecret)
			if err != nil {
				return false, err
			}
			client = c
			ps.razorpayClients[config.ID] = client
		}

		return client.VerifyPaymentSignature(p.GatewayOrderID, p.GatewayPaymentID, signature)

	case models.ProviderBilldesk:
		client, exists := ps.billdeskClients[config.ID]
		if !exists {
			var c *payment.BilldeskClient
			var settings map[string]interface{}
			json.Unmarshal(config.Settings, &settings)

			merchantID := config.ApiKey
			clientID := ""
			environment := "sandbox"

			if v, ok := settings["client_id"].(string); ok {
				clientID = v
			}
			if v, ok := settings["environment"].(string); ok {
				environment = v
			}

			c, err = payment.NewBilldeskClient(config.ApiKey, config.ApiSecret, merchantID, clientID, environment)
			if err != nil {
				return false, err
			}
			client = c
			ps.billdeskClients[config.ID] = client
		}

		return client.VerifyPaymentSignature(signature, string(p.Metadata)), nil

	default:
		return false, fmt.Errorf("unsupported provider: %s", p.Provider)
	}
}

// CreateRefund creates a refund
func (ps *PaymentService) CreateRefund(paymentID, tenantID string, req models.RefundRequest) (*models.Refund, error) {
	p, err := ps.GetPayment(paymentID, tenantID)
	if err != nil {
		return nil, err
	}

	if p.Status != models.PaymentSuccessful {
		return nil, errors.New("only successful payments can be refunded")
	}

	// Get gateway configuration
	config, err := ps.GetActiveGatewayConfig(tenantID, string(p.Provider))
	if err != nil {
		return nil, err
	}

	amount := 0.0
	fmt.Sscanf(req.Amount, "%f", &amount)
	if amount <= 0 || amount > p.Amount {
		return nil, errors.New("invalid refund amount")
	}

	refund := &models.Refund{
		ID:        uuid.New().String(),
		PaymentID: paymentID,
		Amount:    amount,
		Status:    models.PaymentPending,
		Reason:    req.Reason,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Process refund with gateway
	switch p.Provider {
	case models.ProviderRazorpay:
		if err := ps.processRazorpayRefund(refund, config, p); err != nil {
			return nil, err
		}

	case models.ProviderBilldesk:
		if err := ps.processBilldeskRefund(refund, config, p); err != nil {
			return nil, err
		}
	}

	// Store refund in database
	if err := ps.storeRefund(refund); err != nil {
		return nil, err
	}

	// Update payment
	p.RefundID = refund.ID
	p.RefundAmount = amount
	p.Status = models.PaymentRefunded
	ps.UpdatePayment(p)

	return refund, nil
}

// processRazorpayRefund processes Razorpay refund
func (ps *PaymentService) processRazorpayRefund(refund *models.Refund, config *models.PaymentGatewayConfig, p *models.Payment) error {
	client, exists := ps.razorpayClients[config.ID]
	if !exists {
		var err error
		client, err = payment.NewRazorpayClient(config.ApiKey, config.ApiSecret)
		if err != nil {
			return err
		}
		ps.razorpayClients[config.ID] = client
	}

	result, err := client.CreateRefund(p.GatewayPaymentID, refund.Amount, refund.Reason)
	if err != nil {
		return err
	}

	refund.GatewayRefundID = result.ID
	refund.Status = models.PaymentStatus(result.Status)
	refund.ProcessedAt = time.Now()

	return nil
}

// processBilldeskRefund processes Billdesk refund
func (ps *PaymentService) processBilldeskRefund(refund *models.Refund, config *models.PaymentGatewayConfig, p *models.Payment) error {
	var settings map[string]interface{}
	json.Unmarshal(config.Settings, &settings)

	merchantID := config.ApiKey
	clientID := ""
	environment := "sandbox"

	if v, ok := settings["client_id"].(string); ok {
		clientID = v
	}
	if v, ok := settings["environment"].(string); ok {
		environment = v
	}

	client, exists := ps.billdeskClients[config.ID]
	if !exists {
		var err error
		client, err = payment.NewBilldeskClient(config.ApiKey, config.ApiSecret, merchantID, clientID, environment)
		if err != nil {
			return err
		}
		ps.billdeskClients[config.ID] = client
	}

	result, err := client.CreateRefund(p.GatewayOrderID, refund.Amount, refund.Reason)
	if err != nil {
		return err
	}

	refund.GatewayRefundID = result.RefundID
	refund.Status = models.PaymentStatus(result.Status)
	refund.ProcessedAt = time.Now()

	return nil
}

// storeRefund stores refund in database
func (ps *PaymentService) storeRefund(refund *models.Refund) error {
	query := `
		INSERT INTO refunds (id, payment_id, amount, status, gateway_refund_id, reason, created_at, updated_at, processed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := ps.db.Exec(query,
		refund.ID, refund.PaymentID, refund.Amount, refund.Status, refund.GatewayRefundID,
		refund.Reason, refund.CreatedAt, refund.UpdatedAt, refund.ProcessedAt,
	)

	return err
}

// ListPayments lists payments
func (ps *PaymentService) ListPayments(tenantID, status, provider string, limit, offset int) ([]models.Payment, int, error) {
	query := `
		SELECT id, tenant_id, order_id, amount, currency, status, provider, payment_method,
		       transaction_id, gateway_order_id, gateway_payment_id, description, customer_name,
		       customer_email, customer_phone, billing_address, error_message, error_code,
		       receipt_url, refund_id, refund_amount, created_at, updated_at, processed_at, expires_at
		FROM payments
		WHERE tenant_id = $1
	`

	args := []interface{}{tenantID}
	argNum := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argNum)
		args = append(args, status)
		argNum++
	}

	if provider != "" {
		query += fmt.Sprintf(" AND provider = $%d", argNum)
		args = append(args, provider)
		argNum++
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM payments WHERE tenant_id = $1"
	if status != "" {
		countQuery += " AND status = $2"
	}
	if provider != "" {
		if status != "" {
			countQuery += " AND provider = $3"
		} else {
			countQuery += " AND provider = $2"
		}
	}

	var total int
	ps.db.QueryRow(countQuery, args...).Scan(&total)

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, limit, offset)

	rows, err := ps.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	payments := []models.Payment{}
	for rows.Next() {
		p := models.Payment{}
		var billingAddr json.RawMessage

		err := rows.Scan(
			&p.ID, &p.TenantID, &p.OrderID, &p.Amount, &p.Currency, &p.Status, &p.Provider,
			&p.PaymentMethod, &p.TransactionID, &p.GatewayOrderID, &p.GatewayPaymentID, &p.Description,
			&p.CustomerName, &p.CustomerEmail, &p.CustomerPhone, &billingAddr, &p.ErrorMessage,
			&p.ErrorCode, &p.ReceiptURL, &p.RefundID, &p.RefundAmount, &p.CreatedAt, &p.UpdatedAt,
			&p.ProcessedAt, &p.ExpiresAt,
		)
		if err != nil {
			continue
		}

		if billingAddr != nil {
			json.Unmarshal(billingAddr, &p.BillingAddress)
		}

		payments = append(payments, p)
	}

	return payments, total, nil
}

// GetPaymentSummary returns payment summary statistics
func (ps *PaymentService) GetPaymentSummary(tenantID string) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_payments,
			SUM(CASE WHEN status = 'successful' THEN amount ELSE 0 END) as total_successful,
			SUM(CASE WHEN status = 'failed' THEN amount ELSE 0 END) as total_failed,
			SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END) as total_pending,
			COUNT(CASE WHEN status = 'successful' THEN 1 END) as successful_count,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_count
		FROM payments
		WHERE tenant_id = $1
	`

	var totalPayments int
	var totalSuccessful, totalFailed, totalPending sql.NullFloat64
	var successfulCount, failedCount, pendingCount int

	err := ps.db.QueryRow(query, tenantID).Scan(
		&totalPayments, &totalSuccessful, &totalFailed, &totalPending,
		&successfulCount, &failedCount, &pendingCount,
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_payments":   totalPayments,
		"total_successful": totalSuccessful.Float64,
		"total_failed":     totalFailed.Float64,
		"total_pending":    totalPending.Float64,
		"successful_count": successfulCount,
		"failed_count":     failedCount,
		"pending_count":    pendingCount,
		"success_rate":     calculateSuccessRate(successfulCount, totalPayments),
	}, nil
}

// GetActiveGatewayConfig retrieves active gateway configuration
func (ps *PaymentService) GetActiveGatewayConfig(tenantID, provider string) (*models.PaymentGatewayConfig, error) {
	query := `
		SELECT id, tenant_id, provider, is_active, api_key, api_secret, settings, created_at, updated_at
		FROM payment_gateway_config
		WHERE tenant_id = $1 AND provider = $2 AND is_active = true
		LIMIT 1
	`

	config := &models.PaymentGatewayConfig{}
	var settings json.RawMessage

	err := ps.db.QueryRow(query, tenantID, provider).Scan(
		&config.ID, &config.TenantID, &config.Provider, &config.IsActive,
		&config.ApiKey, &config.ApiSecret, &settings, &config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no active %s configuration found", provider)
		}
		return nil, err
	}

	config.Settings = settings
	return config, nil
}

// GetGatewayConfig retrieves gateway configuration by ID
func (ps *PaymentService) GetGatewayConfig(configID string) (*models.PaymentGatewayConfig, error) {
	query := `
		SELECT id, tenant_id, provider, is_active, api_key, api_secret, settings, created_at, updated_at
		FROM payment_gateway_config
		WHERE id = $1
	`

	config := &models.PaymentGatewayConfig{}
	var settings json.RawMessage

	err := ps.db.QueryRow(query, configID).Scan(
		&config.ID, &config.TenantID, &config.Provider, &config.IsActive,
		&config.ApiKey, &config.ApiSecret, &settings, &config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	config.Settings = settings
	return config, nil
}

// GetPaymentGateways retrieves all gateway configurations
func (ps *PaymentService) GetPaymentGateways(tenantID string) ([]models.PaymentGatewayConfig, error) {
	query := `
		SELECT id, tenant_id, provider, is_active, api_key, api_secret, settings, created_at, updated_at
		FROM payment_gateway_config
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := ps.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	configs := []models.PaymentGatewayConfig{}
	for rows.Next() {
		config := models.PaymentGatewayConfig{}
		var settings json.RawMessage

		err := rows.Scan(
			&config.ID, &config.TenantID, &config.Provider, &config.IsActive,
			&config.ApiKey, &config.ApiSecret, &settings, &config.CreatedAt, &config.UpdatedAt,
		)
		if err != nil {
			continue
		}

		config.Settings = settings
		configs = append(configs, config)
	}

	return configs, nil
}

// ConfigureGateway configures a payment gateway
func (ps *PaymentService) ConfigureGateway(tenantID, provider, apiKey, apiSecret string, settings map[string]interface{}) (*models.PaymentGatewayConfig, error) {
	config := &models.PaymentGatewayConfig{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		Provider:  models.PaymentProvider(provider),
		IsActive:  true,
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	settingsJSON, _ := json.Marshal(settings)
	config.Settings = settingsJSON

	query := `
		INSERT INTO payment_gateway_config (id, tenant_id, provider, is_active, api_key, api_secret, settings, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := ps.db.Exec(query,
		config.ID, config.TenantID, config.Provider, config.IsActive,
		config.ApiKey, config.ApiSecret, config.Settings, config.CreatedAt, config.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// UpdateGateway updates a gateway configuration
func (ps *PaymentService) UpdateGateway(configID, tenantID, apiKey, apiSecret string, isActive *bool, settings map[string]interface{}) (*models.PaymentGatewayConfig, error) {
	config, err := ps.GetGatewayConfig(configID)
	if err != nil {
		return nil, err
	}

	if config.TenantID != tenantID {
		return nil, errors.New("unauthorized")
	}

	if apiKey != "" {
		config.ApiKey = apiKey
	}
	if apiSecret != "" {
		config.ApiSecret = apiSecret
	}
	if isActive != nil {
		config.IsActive = *isActive
	}
	if settings != nil {
		settingsJSON, _ := json.Marshal(settings)
		config.Settings = settingsJSON
	}

	config.UpdatedAt = time.Now()

	query := `
		UPDATE payment_gateway_config
		SET api_key = $1, api_secret = $2, is_active = $3, settings = $4, updated_at = $5
		WHERE id = $6
	`

	_, err = ps.db.Exec(query,
		config.ApiKey, config.ApiSecret, config.IsActive, config.Settings, config.UpdatedAt, config.ID,
	)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// DeleteGateway deletes a gateway configuration
func (ps *PaymentService) DeleteGateway(configID, tenantID string) error {
	config, err := ps.GetGatewayConfig(configID)
	if err != nil {
		return err
	}

	if config.TenantID != tenantID {
		return errors.New("unauthorized")
	}

	query := "DELETE FROM payment_gateway_config WHERE id = $1"
	_, err = ps.db.Exec(query, configID)
	return err
}

// HandleRazorpayWebhook handles Razorpay webhook
func (ps *PaymentService) HandleRazorpayWebhook(webhook models.RazorpayWebhook, signature string) error {
	// Verify signature and process webhook
	// Implementation details depend on your webhook verification strategy
	return nil
}

// HandleBilldeskWebhook handles Billdesk webhook
func (ps *PaymentService) HandleBilldeskWebhook(webhook models.BilldeskWebhook, body, signature string) error {
	// Verify signature and process webhook
	// Implementation details depend on your webhook verification strategy
	return nil
}

// Helper functions
func calculateSuccessRate(successCount, totalCount int) float64 {
	if totalCount == 0 {
		return 0.0
	}
	return float64(successCount) / float64(totalCount) * 100
}
