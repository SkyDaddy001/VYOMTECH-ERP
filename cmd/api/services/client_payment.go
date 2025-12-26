package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"lms/cmd/api/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// TenantAccountService handles tenant account operations
type TenantAccountService struct {
	db *sql.DB
}

// NewTenantAccountService creates a new tenant account service
func NewTenantAccountService(db *sql.DB) *TenantAccountService {
	return &TenantAccountService{db: db}
}

// CreateTenantAccount creates a new tenant account for collecting payments
func (s *TenantAccountService) CreateTenantAccount(tenantID string, req *models.TenantAccountRequest) (*models.TenantAccount, error) {
	account := &models.TenantAccount{
		ID:              uuid.New().String(),
		TenantID:        tenantID,
		ChargeType:      req.ChargeType,
		ChargeTypeName:  req.ChargeTypeName,
		Description:     req.Description,
		BankAccountName: req.BankAccountName,
		BankAccountNo:   req.BankAccountNo,
		IFSCCode:        req.IFSCCode,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	query := `
		INSERT INTO tenant_accounts (
			id, tenant_id, charge_type, charge_type_name, description,
			bank_account_name, bank_account_no, ifsc_code, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(query,
		account.ID, account.TenantID, account.ChargeType, account.ChargeTypeName,
		account.Description, account.BankAccountName, account.BankAccountNo, account.IFSCCode,
		account.IsActive, account.CreatedAt, account.UpdatedAt,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" { // Unique violation
			return nil, fmt.Errorf("account for this charge type already exists")
		}
		return nil, err
	}

	return account, nil
}

// GetTenantAccount retrieves a tenant account by ID
func (s *TenantAccountService) GetTenantAccount(accountID string) (*models.TenantAccount, error) {
	account := &models.TenantAccount{}

	query := `
		SELECT id, tenant_id, charge_type, charge_type_name, description,
		       razorpay_account_id, billdesk_account_id, bank_account_name,
		       bank_account_no, ifsc_code, is_active, total_collected,
		       total_refunded, created_at, updated_at
		FROM tenant_accounts
		WHERE id = $1
	`

	err := s.db.QueryRow(query, accountID).Scan(
		&account.ID, &account.TenantID, &account.ChargeType, &account.ChargeTypeName,
		&account.Description, &account.RazorpayID, &account.BilldeskID,
		&account.BankAccountName, &account.BankAccountNo, &account.IFSCCode,
		&account.IsActive, &account.TotalCollected, &account.TotalRefunded,
		&account.CreatedAt, &account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetTenantAccountByChargeType retrieves tenant account by charge type
func (s *TenantAccountService) GetTenantAccountByChargeType(tenantID string, chargeType models.ChargeType) (*models.TenantAccount, error) {
	account := &models.TenantAccount{}

	query := `
		SELECT id, tenant_id, charge_type, charge_type_name, description,
		       razorpay_account_id, billdesk_account_id, bank_account_name,
		       bank_account_no, ifsc_code, is_active, total_collected,
		       total_refunded, created_at, updated_at
		FROM tenant_accounts
		WHERE tenant_id = $1 AND charge_type = $2
	`

	err := s.db.QueryRow(query, tenantID, chargeType).Scan(
		&account.ID, &account.TenantID, &account.ChargeType, &account.ChargeTypeName,
		&account.Description, &account.RazorpayID, &account.BilldeskID,
		&account.BankAccountName, &account.BankAccountNo, &account.IFSCCode,
		&account.IsActive, &account.TotalCollected, &account.TotalRefunded,
		&account.CreatedAt, &account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found for charge type: %s", chargeType)
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

// ListTenantAccounts lists all accounts for a tenant
func (s *TenantAccountService) ListTenantAccounts(tenantID string) ([]models.TenantAccount, error) {
	query := `
		SELECT id, tenant_id, charge_type, charge_type_name, description,
		       razorpay_account_id, billdesk_account_id, bank_account_name,
		       bank_account_no, ifsc_code, is_active, total_collected,
		       total_refunded, created_at, updated_at
		FROM tenant_accounts
		WHERE tenant_id = $1
		ORDER BY charge_type
	`

	rows, err := s.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.TenantAccount

	for rows.Next() {
		account := models.TenantAccount{}
		err := rows.Scan(
			&account.ID, &account.TenantID, &account.ChargeType, &account.ChargeTypeName,
			&account.Description, &account.RazorpayID, &account.BilldeskID,
			&account.BankAccountName, &account.BankAccountNo, &account.IFSCCode,
			&account.IsActive, &account.TotalCollected, &account.TotalRefunded,
			&account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

// UpdateTenantAccount updates a tenant account
func (s *TenantAccountService) UpdateTenantAccount(accountID string, req *models.TenantAccountRequest) (*models.TenantAccount, error) {
	query := `
		UPDATE tenant_accounts
		SET charge_type_name = $1, description = $2,
		    bank_account_name = $3, bank_account_no = $4, ifsc_code = $5
		WHERE id = $6
		RETURNING id, tenant_id, charge_type, charge_type_name, description,
		          razorpay_account_id, billdesk_account_id, bank_account_name,
		          bank_account_no, ifsc_code, is_active, total_collected,
		          total_refunded, created_at, updated_at
	`

	account := &models.TenantAccount{}
	err := s.db.QueryRow(query,
		req.ChargeTypeName, req.Description,
		req.BankAccountName, req.BankAccountNo, req.IFSCCode, accountID,
	).Scan(
		&account.ID, &account.TenantID, &account.ChargeType, &account.ChargeTypeName,
		&account.Description, &account.RazorpayID, &account.BilldeskID,
		&account.BankAccountName, &account.BankAccountNo, &account.IFSCCode,
		&account.IsActive, &account.TotalCollected, &account.TotalRefunded,
		&account.CreatedAt, &account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

// UpdateGatewayCredentials updates payment gateway credentials for an account
func (s *TenantAccountService) UpdateGatewayCredentials(accountID, razorpayID, billdeskID string) error {
	query := `
		UPDATE tenant_accounts
		SET razorpay_account_id = COALESCE(NULLIF($1, ''), razorpay_account_id),
		    billdesk_account_id = COALESCE(NULLIF($2, ''), billdesk_account_id)
		WHERE id = $3
	`

	result, err := s.db.Exec(query, razorpayID, billdeskID, accountID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

// ClientInvoiceService handles client invoice operations
type ClientInvoiceService struct {
	db *sql.DB
}

// NewClientInvoiceService creates a new client invoice service
func NewClientInvoiceService(db *sql.DB) *ClientInvoiceService {
	return &ClientInvoiceService{db: db}
}

// CreateInvoice creates a new client invoice
func (s *ClientInvoiceService) CreateInvoice(tenantID string, req *models.CreateInvoiceRequest) (*models.ClientInvoice, error) {
	invoice := &models.ClientInvoice{
		ID:                uuid.New().String(),
		TenantID:          tenantID,
		ClientID:          req.ClientID,
		ClientName:        req.ClientName,
		ClientEmail:       req.ClientEmail,
		ClientPhone:       req.ClientPhone,
		ChargeType:        req.ChargeType,
		InvoiceNumber:     fmt.Sprintf("INV-%d-%s", time.Now().Unix(), uuid.New().String()[:8]),
		Amount:            req.Amount,
		OutstandingAmount: req.Amount,
		Currency:          "INR",
		Description:       req.Description,
		InvoiceDate:       time.Now(),
		DueDate:           req.DueDate,
		Status:            "issued",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	query := `
		INSERT INTO client_invoices (
			id, tenant_id, client_id, client_name, client_email, client_phone,
			charge_type, invoice_number, amount, amount_paid, outstanding_amount,
			currency, description, invoice_date, due_date, status, metadata, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id, created_at, updated_at
	`

	var metadata interface{}
	if req.Metadata != nil {
		data, _ := json.Marshal(req.Metadata)
		metadata = data
	}

	err := s.db.QueryRow(query,
		invoice.ID, invoice.TenantID, invoice.ClientID, invoice.ClientName,
		invoice.ClientEmail, invoice.ClientPhone, invoice.ChargeType, invoice.InvoiceNumber,
		invoice.Amount, 0, invoice.OutstandingAmount, invoice.Currency, invoice.Description,
		invoice.InvoiceDate, invoice.DueDate, invoice.Status, metadata, invoice.CreatedAt, invoice.UpdatedAt,
	).Scan(&invoice.ID, &invoice.CreatedAt, &invoice.UpdatedAt)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" { // Unique violation
			return nil, fmt.Errorf("invoice number already exists")
		}
		return nil, err
	}

	return invoice, nil
}

// CreateBulkInvoices creates multiple invoices for multiple clients
func (s *ClientInvoiceService) CreateBulkInvoices(tenantID string, req *models.BulkInvoiceRequest) ([]models.ClientInvoice, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var invoices []models.ClientInvoice

	query := `
		INSERT INTO client_invoices (
			id, tenant_id, client_id, client_name, client_email, client_phone,
			charge_type, invoice_number, amount, amount_paid, outstanding_amount,
			currency, description, invoice_date, due_date, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, created_at, updated_at
	`

	for _, item := range req.Invoices {
		invoice := models.ClientInvoice{
			ID:                uuid.New().String(),
			TenantID:          tenantID,
			ClientID:          item.ClientID,
			ClientName:        item.ClientName,
			ClientEmail:       item.ClientEmail,
			ClientPhone:       item.ClientPhone,
			ChargeType:        req.ChargeType,
			InvoiceNumber:     fmt.Sprintf("INV-%d-%s", time.Now().UnixNano(), uuid.New().String()[:8]),
			Amount:            item.Amount,
			OutstandingAmount: item.Amount,
			Currency:          "INR",
			Description:       req.Description,
			InvoiceDate:       time.Now(),
			DueDate:           req.DueDate,
			Status:            "issued",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		err := tx.QueryRow(query,
			invoice.ID, invoice.TenantID, invoice.ClientID, invoice.ClientName,
			invoice.ClientEmail, invoice.ClientPhone, invoice.ChargeType, invoice.InvoiceNumber,
			invoice.Amount, 0, invoice.OutstandingAmount, invoice.Currency, invoice.Description,
			invoice.InvoiceDate, invoice.DueDate, invoice.Status, invoice.CreatedAt, invoice.UpdatedAt,
		).Scan(&invoice.ID, &invoice.CreatedAt, &invoice.UpdatedAt)

		if err != nil {
			return nil, err
		}

		invoices = append(invoices, invoice)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return invoices, nil
}

// GetInvoice retrieves a specific invoice
func (s *ClientInvoiceService) GetInvoice(invoiceID string) (*models.ClientInvoice, error) {
	invoice := &models.ClientInvoice{}

	query := `
		SELECT id, tenant_id, client_id, client_name, client_email, client_phone,
		       charge_type, invoice_number, amount, amount_paid, outstanding_amount,
		       currency, description, invoice_date, due_date, status, metadata,
		       created_at, updated_at
		FROM client_invoices
		WHERE id = $1
	`

	err := s.db.QueryRow(query, invoiceID).Scan(
		&invoice.ID, &invoice.TenantID, &invoice.ClientID, &invoice.ClientName,
		&invoice.ClientEmail, &invoice.ClientPhone, &invoice.ChargeType,
		&invoice.InvoiceNumber, &invoice.Amount, &invoice.AmountPaid, &invoice.OutstandingAmount,
		&invoice.Currency, &invoice.Description, &invoice.InvoiceDate, &invoice.DueDate,
		&invoice.Status, &invoice.Metadata, &invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invoice not found")
	}
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

// GetClientInvoices retrieves all invoices for a client
func (s *ClientInvoiceService) GetClientInvoices(tenantID, clientID string) ([]models.ClientInvoice, error) {
	query := `
		SELECT id, tenant_id, client_id, client_name, client_email, client_phone,
		       charge_type, invoice_number, amount, amount_paid, outstanding_amount,
		       currency, description, invoice_date, due_date, status, metadata,
		       created_at, updated_at
		FROM client_invoices
		WHERE tenant_id = $1 AND client_id = $2
		ORDER BY invoice_date DESC
	`

	rows, err := s.db.Query(query, tenantID, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.ClientInvoice

	for rows.Next() {
		invoice := models.ClientInvoice{}
		err := rows.Scan(
			&invoice.ID, &invoice.TenantID, &invoice.ClientID, &invoice.ClientName,
			&invoice.ClientEmail, &invoice.ClientPhone, &invoice.ChargeType,
			&invoice.InvoiceNumber, &invoice.Amount, &invoice.AmountPaid, &invoice.OutstandingAmount,
			&invoice.Currency, &invoice.Description, &invoice.InvoiceDate, &invoice.DueDate,
			&invoice.Status, &invoice.Metadata, &invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return invoices, nil
}

// GetOutstandingInvoices retrieves unpaid/partially paid invoices
func (s *ClientInvoiceService) GetOutstandingInvoices(tenantID, clientID string) ([]models.ClientInvoice, error) {
	query := `
		SELECT id, tenant_id, client_id, client_name, client_email, client_phone,
		       charge_type, invoice_number, amount, amount_paid, outstanding_amount,
		       currency, description, invoice_date, due_date, status, metadata,
		       created_at, updated_at
		FROM client_invoices
		WHERE tenant_id = $1 AND client_id = $2 AND outstanding_amount > 0
		ORDER BY due_date ASC
	`

	rows, err := s.db.Query(query, tenantID, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.ClientInvoice

	for rows.Next() {
		invoice := models.ClientInvoice{}
		err := rows.Scan(
			&invoice.ID, &invoice.TenantID, &invoice.ClientID, &invoice.ClientName,
			&invoice.ClientEmail, &invoice.ClientPhone, &invoice.ChargeType,
			&invoice.InvoiceNumber, &invoice.Amount, &invoice.AmountPaid, &invoice.OutstandingAmount,
			&invoice.Currency, &invoice.Description, &invoice.InvoiceDate, &invoice.DueDate,
			&invoice.Status, &invoice.Metadata, &invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return invoices, nil
}

// GetInvoicesByChargeType retrieves invoices by charge type
func (s *ClientInvoiceService) GetInvoicesByChargeType(tenantID string, chargeType models.ChargeType) ([]models.ClientInvoice, error) {
	query := `
		SELECT id, tenant_id, client_id, client_name, client_email, client_phone,
		       charge_type, invoice_number, amount, amount_paid, outstanding_amount,
		       currency, description, invoice_date, due_date, status, metadata,
		       created_at, updated_at
		FROM client_invoices
		WHERE tenant_id = $1 AND charge_type = $2
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, tenantID, chargeType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.ClientInvoice

	for rows.Next() {
		invoice := models.ClientInvoice{}
		err := rows.Scan(
			&invoice.ID, &invoice.TenantID, &invoice.ClientID, &invoice.ClientName,
			&invoice.ClientEmail, &invoice.ClientPhone, &invoice.ChargeType,
			&invoice.InvoiceNumber, &invoice.Amount, &invoice.AmountPaid, &invoice.OutstandingAmount,
			&invoice.Currency, &invoice.Description, &invoice.InvoiceDate, &invoice.DueDate,
			&invoice.Status, &invoice.Metadata, &invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return invoices, nil
}

// UpdateInvoiceStatus updates invoice payment status
func (s *ClientInvoiceService) UpdateInvoiceStatus(invoiceID string, amountPaid float64, status string) error {
	query := `
		UPDATE client_invoices
		SET amount_paid = $1,
		    outstanding_amount = GREATEST(0, amount - $1),
		    status = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := s.db.Exec(query, amountPaid, status, invoiceID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invoice not found")
	}

	return nil
}

// ClientPaymentService handles client payment operations
type ClientPaymentService struct {
	db *sql.DB
}

// NewClientPaymentService creates a new client payment service
func NewClientPaymentService(db *sql.DB) *ClientPaymentService {
	return &ClientPaymentService{db: db}
}

// CreateClientPayment creates a new client payment record
func (s *ClientPaymentService) CreateClientPayment(payment *models.ClientPayment) error {
	query := `
		INSERT INTO client_payments (
			id, tenant_id, client_id, invoice_id, tenant_account_id,
			charge_type, order_id, amount, currency, status, payment_type,
			provider, payment_method, client_name, client_email, client_phone,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := s.db.Exec(query,
		payment.ID, payment.TenantID, payment.ClientID, payment.InvoiceID, payment.TenantAccountID,
		payment.ChargeType, payment.OrderID, payment.Amount, payment.Currency, payment.Status,
		payment.PaymentType, payment.Provider, payment.PaymentMethod, payment.ClientName,
		payment.ClientEmail, payment.ClientPhone, payment.CreatedAt, payment.UpdatedAt,
	)

	return err
}

// GetClientPayment retrieves a client payment by ID
func (s *ClientPaymentService) GetClientPayment(paymentID string) (*models.ClientPayment, error) {
	payment := &models.ClientPayment{}

	query := `
		SELECT id, tenant_id, client_id, invoice_id, tenant_account_id,
		       charge_type, order_id, amount, currency, status, payment_type,
		       provider, payment_method, gateway_order_id, gateway_payment_id,
		       transaction_id, client_name, client_email, client_phone,
		       receipt_url, refund_id, refund_amount, error_message,
		       created_at, updated_at, processed_at
		FROM client_payments
		WHERE id = $1
	`

	err := s.db.QueryRow(query, paymentID).Scan(
		&payment.ID, &payment.TenantID, &payment.ClientID, &payment.InvoiceID, &payment.TenantAccountID,
		&payment.ChargeType, &payment.OrderID, &payment.Amount, &payment.Currency, &payment.Status,
		&payment.PaymentType, &payment.Provider, &payment.PaymentMethod, &payment.GatewayOrderID,
		&payment.GatewayPaymentID, &payment.TransactionID, &payment.ClientName, &payment.ClientEmail,
		&payment.ClientPhone, &payment.ReceiptURL, &payment.RefundID, &payment.RefundAmount,
		&payment.ErrorMessage, &payment.CreatedAt, &payment.UpdatedAt, &payment.ProcessedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetClientPaymentsByInvoice retrieves all payments for an invoice
func (s *ClientPaymentService) GetClientPaymentsByInvoice(invoiceID string) ([]models.ClientPayment, error) {
	query := `
		SELECT id, tenant_id, client_id, invoice_id, tenant_account_id,
		       charge_type, order_id, amount, currency, status, payment_type,
		       provider, payment_method, gateway_order_id, gateway_payment_id,
		       transaction_id, client_name, client_email, client_phone,
		       receipt_url, refund_id, refund_amount, error_message,
		       created_at, updated_at, processed_at
		FROM client_payments
		WHERE invoice_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.ClientPayment

	for rows.Next() {
		payment := models.ClientPayment{}
		err := rows.Scan(
			&payment.ID, &payment.TenantID, &payment.ClientID, &payment.InvoiceID, &payment.TenantAccountID,
			&payment.ChargeType, &payment.OrderID, &payment.Amount, &payment.Currency, &payment.Status,
			&payment.PaymentType, &payment.Provider, &payment.PaymentMethod, &payment.GatewayOrderID,
			&payment.GatewayPaymentID, &payment.TransactionID, &payment.ClientName, &payment.ClientEmail,
			&payment.ClientPhone, &payment.ReceiptURL, &payment.RefundID, &payment.RefundAmount,
			&payment.ErrorMessage, &payment.CreatedAt, &payment.UpdatedAt, &payment.ProcessedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

// ListClientPayments retrieves payments for a client within a date range
func (s *ClientPaymentService) ListClientPayments(tenantID, clientID string, limit, offset int) ([]models.ClientPayment, int64, error) {
	var total int64

	// Get total count
	countQuery := `SELECT COUNT(*) FROM client_payments WHERE tenant_id = $1 AND client_id = $2`
	err := s.db.QueryRow(countQuery, tenantID, clientID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT id, tenant_id, client_id, invoice_id, tenant_account_id,
		       charge_type, order_id, amount, currency, status, payment_type,
		       provider, payment_method, gateway_order_id, gateway_payment_id,
		       transaction_id, client_name, client_email, client_phone,
		       receipt_url, refund_id, refund_amount, error_message,
		       created_at, updated_at, processed_at
		FROM client_payments
		WHERE tenant_id = $1 AND client_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := s.db.Query(query, tenantID, clientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var payments []models.ClientPayment

	for rows.Next() {
		payment := models.ClientPayment{}
		err := rows.Scan(
			&payment.ID, &payment.TenantID, &payment.ClientID, &payment.InvoiceID, &payment.TenantAccountID,
			&payment.ChargeType, &payment.OrderID, &payment.Amount, &payment.Currency, &payment.Status,
			&payment.PaymentType, &payment.Provider, &payment.PaymentMethod, &payment.GatewayOrderID,
			&payment.GatewayPaymentID, &payment.TransactionID, &payment.ClientName, &payment.ClientEmail,
			&payment.ClientPhone, &payment.ReceiptURL, &payment.RefundID, &payment.RefundAmount,
			&payment.ErrorMessage, &payment.CreatedAt, &payment.UpdatedAt, &payment.ProcessedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// UpdateClientPaymentStatus updates payment status after gateway verification
func (s *ClientPaymentService) UpdateClientPaymentStatus(paymentID string, status models.PaymentStatus, gatewayPaymentID string) error {
	query := `
		UPDATE client_payments
		SET status = $1, gateway_payment_id = $2, processed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := s.db.Exec(query, status, gatewayPaymentID, paymentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// GetClientOutstanding retrieves outstanding balance summary for a client
func (s *ClientPaymentService) GetClientOutstanding(tenantID, clientID string) (*models.ClientOutstandingResponse, error) {
	response := &models.ClientOutstandingResponse{
		ClientID:     clientID,
		ByChargeType: make(map[models.ChargeType]models.OutstandingByType),
	}

	// Get invoices
	invoiceService := NewClientInvoiceService(s.db)
	invoices, err := invoiceService.GetClientInvoices(tenantID, clientID)
	if err != nil {
		return nil, err
	}

	if len(invoices) == 0 {
		return response, nil
	}

	response.ClientName = invoices[0].ClientName
	response.ClientEmail = invoices[0].ClientEmail
	response.Invoices = invoices

	// Calculate totals and by-charge-type breakdown
	byType := make(map[models.ChargeType]models.OutstandingByType)

	for _, inv := range invoices {
		response.TotalPaid += inv.AmountPaid
		response.TotalOutstanding += inv.OutstandingAmount

		chargeType := inv.ChargeType
		stats := byType[chargeType]
		stats.ChargeType = chargeType
		stats.Total += inv.Amount
		stats.Paid += inv.AmountPaid
		stats.Outstanding += inv.OutstandingAmount
		byType[chargeType] = stats
	}

	response.ByChargeType = byType

	return response, nil
}

// GetCollectionDashboard retrieves collection statistics for a tenant
func (s *ClientPaymentService) GetCollectionDashboard(tenantID string, limit int) (*models.TenantCollectionDashboard, error) {
	dashboard := &models.TenantCollectionDashboard{
		TenantID:         tenantID,
		CollectionByType: make(map[models.ChargeType]models.CollectionStats),
	}

	// Get recent payments
	query := `
		SELECT id, tenant_id, client_id, invoice_id, tenant_account_id,
		       charge_type, order_id, amount, currency, status, payment_type,
		       provider, payment_method, gateway_order_id, gateway_payment_id,
		       transaction_id, client_name, client_email, client_phone,
		       receipt_url, refund_id, refund_amount, error_message,
		       created_at, updated_at, processed_at
		FROM client_payments
		WHERE tenant_id = $1 AND status = 'completed'
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(query, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		payment := models.ClientPayment{}
		err := rows.Scan(
			&payment.ID, &payment.TenantID, &payment.ClientID, &payment.InvoiceID, &payment.TenantAccountID,
			&payment.ChargeType, &payment.OrderID, &payment.Amount, &payment.Currency, &payment.Status,
			&payment.PaymentType, &payment.Provider, &payment.PaymentMethod, &payment.GatewayOrderID,
			&payment.GatewayPaymentID, &payment.TransactionID, &payment.ClientName, &payment.ClientEmail,
			&payment.ClientPhone, &payment.ReceiptURL, &payment.RefundID, &payment.RefundAmount,
			&payment.ErrorMessage, &payment.CreatedAt, &payment.UpdatedAt, &payment.ProcessedAt,
		)
		if err != nil {
			return nil, err
		}

		dashboard.RecentPayments = append(dashboard.RecentPayments, payment)
		dashboard.TotalCollected += payment.Amount
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Get collection statistics by type from view
	statsQuery := `
		SELECT charge_type, total_billed, COALESCE(total_collected, 0),
		       COALESCE(outstanding_amount, 0), collection_rate, 
		       COALESCE(total_invoices, 0), COALESCE(paid_invoices, 0)
		FROM v_tenant_collection_summary
		WHERE tenant_id = $1
	`

	statsRows, err := s.db.Query(statsQuery, tenantID)
	if err != nil {
		return nil, err
	}
	defer statsRows.Close()

	for statsRows.Next() {
		var chargeType models.ChargeType
		stats := models.CollectionStats{}

		err := statsRows.Scan(
			&chargeType, &stats.TotalBilled, &stats.TotalCollected,
			&stats.Outstanding, &stats.CollectionRate, &stats.InvoiceCount, &stats.PaidInvoices,
		)
		if err != nil {
			return nil, err
		}

		stats.ChargeType = chargeType
		stats.OverdueAmount = 0 // TODO: Calculate from overdue invoices

		dashboard.CollectionByType[chargeType] = stats
		dashboard.TotalOutstanding += stats.Outstanding
	}

	if err = statsRows.Err(); err != nil {
		return nil, err
	}

	return dashboard, nil
}
