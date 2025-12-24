package models

import (
	"encoding/json"
	"time"
)

// ChargeType defines types of charges/accounts for tenant collections
type ChargeType string

const (
	ChargeApartmentCost  ChargeType = "apartment_cost"
	ChargeMaintenance    ChargeType = "maintenance"
	ChargeOtherCharges   ChargeType = "other_charges"
	ChargePropertyTax    ChargeType = "property_tax"
	ChargeWaterCharges   ChargeType = "water_charges"
	ChargeElectricityTax ChargeType = "electricity_tax"
)

// PaymentType defines payment types (tenant admin vs client)
type PaymentType string

const (
	PaymentTypeAdmin  PaymentType = "admin"  // Admin/subscription payments
	PaymentTypeClient PaymentType = "client" // Client/resident payments
	PaymentTypeBulk   PaymentType = "bulk"   // Bulk collection payments
)

// TenantAccount represents a tenant's revenue account for a specific charge type
type TenantAccount struct {
	ID              string     `db:"id" json:"id"`
	TenantID        string     `db:"tenant_id" json:"tenant_id"`
	ChargeType      ChargeType `db:"charge_type" json:"charge_type"`
	ChargeTypeName  string     `db:"charge_type_name" json:"charge_type_name"`
	Description     string     `db:"description" json:"description"`
	RazorpayID      string     `db:"razorpay_account_id" json:"razorpay_account_id"`
	BilldeskID      string     `db:"billdesk_account_id" json:"billdesk_account_id"`
	BankAccountName string     `db:"bank_account_name" json:"bank_account_name"`
	BankAccountNo   string     `db:"bank_account_no" json:"bank_account_no"`
	IFSCCode        string     `db:"ifsc_code" json:"ifsc_code"`
	IsActive        bool       `db:"is_active" json:"is_active"`
	TotalCollected  float64    `db:"total_collected" json:"total_collected"`
	TotalRefunded   float64    `db:"total_refunded" json:"total_refunded"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

// ClientInvoice represents an invoice for a client (resident)
type ClientInvoice struct {
	ID                string          `db:"id" json:"id"`
	TenantID          string          `db:"tenant_id" json:"tenant_id"`
	ClientID          string          `db:"client_id" json:"client_id"`
	ClientName        string          `db:"client_name" json:"client_name"`
	ClientEmail       string          `db:"client_email" json:"client_email"`
	ClientPhone       string          `db:"client_phone" json:"client_phone"`
	ChargeType        ChargeType      `db:"charge_type" json:"charge_type"`
	InvoiceNumber     string          `db:"invoice_number" json:"invoice_number"`
	Amount            float64         `db:"amount" json:"amount"`
	AmountPaid        float64         `db:"amount_paid" json:"amount_paid"`
	OutstandingAmount float64         `db:"outstanding_amount" json:"outstanding_amount"`
	Currency          string          `db:"currency" json:"currency"`
	Description       string          `db:"description" json:"description"`
	InvoiceDate       time.Time       `db:"invoice_date" json:"invoice_date"`
	DueDate           time.Time       `db:"due_date" json:"due_date"`
	Status            string          `db:"status" json:"status"` // draft, issued, partial_paid, paid, overdue
	Metadata          json.RawMessage `db:"metadata" json:"metadata"`
	CreatedAt         time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time       `db:"updated_at" json:"updated_at"`
}

// ClientPayment represents payment from a client towards an invoice
type ClientPayment struct {
	ID               string          `db:"id" json:"id"`
	TenantID         string          `db:"tenant_id" json:"tenant_id"`
	ClientID         string          `db:"client_id" json:"client_id"`
	InvoiceID        string          `db:"invoice_id" json:"invoice_id"`
	TenantAccountID  string          `db:"tenant_account_id" json:"tenant_account_id"`
	ChargeType       ChargeType      `db:"charge_type" json:"charge_type"`
	OrderID          string          `db:"order_id" json:"order_id"`
	Amount           float64         `db:"amount" json:"amount"`
	Currency         string          `db:"currency" json:"currency"`
	Status           PaymentStatus   `db:"status" json:"status"`
	PaymentType      PaymentType     `db:"payment_type" json:"payment_type"`
	Provider         PaymentProvider `db:"provider" json:"provider"`
	PaymentMethod    PaymentMethod   `db:"payment_method" json:"payment_method"`
	GatewayOrderID   string          `db:"gateway_order_id" json:"gateway_order_id"`
	GatewayPaymentID string          `db:"gateway_payment_id" json:"gateway_payment_id"`
	TransactionID    string          `db:"transaction_id" json:"transaction_id"`
	ClientName       string          `db:"client_name" json:"client_name"`
	ClientEmail      string          `db:"client_email" json:"client_email"`
	ClientPhone      string          `db:"client_phone" json:"client_phone"`
	ReceiptURL       string          `db:"receipt_url" json:"receipt_url"`
	RefundID         string          `db:"refund_id" json:"refund_id"`
	RefundAmount     float64         `db:"refund_amount" json:"refund_amount"`
	ErrorMessage     string          `db:"error_message" json:"error_message"`
	CreatedAt        time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at" json:"updated_at"`
	ProcessedAt      *time.Time      `db:"processed_at" json:"processed_at"`
}

// ClientPaymentRequest represents a client payment request
type ClientPaymentRequest struct {
	InvoiceID     string          `json:"invoice_id" binding:"required"`
	Amount        float64         `json:"amount" binding:"required,gt=0"`
	Currency      string          `json:"currency" binding:"required" default:"INR"`
	Provider      PaymentProvider `json:"provider" binding:"required,oneof=razorpay billdesk"`
	PaymentMethod PaymentMethod   `json:"payment_method" binding:"required"`
	ChargeType    ChargeType      `json:"charge_type" binding:"required"`
	ClientName    string          `json:"client_name" binding:"required"`
	ClientEmail   string          `json:"client_email" binding:"required,email"`
	ClientPhone   string          `json:"client_phone" binding:"required"`
}

// TenantCollectionDashboard represents collection statistics for a tenant
type TenantCollectionDashboard struct {
	TenantID            string                         `json:"tenant_id"`
	TotalCollected      float64                        `json:"total_collected"`
	TotalOutstanding    float64                        `json:"total_outstanding"`
	TotalClients        int                            `json:"total_clients"`
	PartialPaidInvoices int                            `json:"partial_paid_invoices"`
	OverdueInvoices     int                            `json:"overdue_invoices"`
	CollectionByType    map[ChargeType]CollectionStats `json:"collection_by_type"`
	RecentPayments      []ClientPayment                `json:"recent_payments"`
}

// CollectionStats represents collection statistics for a charge type
type CollectionStats struct {
	ChargeType     ChargeType `json:"charge_type"`
	TotalBilled    float64    `json:"total_billed"`
	TotalCollected float64    `json:"total_collected"`
	Outstanding    float64    `json:"outstanding"`
	CollectionRate float64    `json:"collection_rate"` // percentage
	InvoiceCount   int        `json:"invoice_count"`
	PaidInvoices   int        `json:"paid_invoices"`
	OverdueAmount  float64    `json:"overdue_amount"`
}

// ClientOutstandingRequest represents client's outstanding balance request
type ClientOutstandingRequest struct {
	ClientID string `uri:"client_id" binding:"required"`
	TenantID string `uri:"tenant_id" binding:"required"`
}

// ClientOutstandingResponse represents client's outstanding amount summary
type ClientOutstandingResponse struct {
	ClientID         string                           `json:"client_id"`
	ClientName       string                           `json:"client_name"`
	ClientEmail      string                           `json:"client_email"`
	TotalOutstanding float64                          `json:"total_outstanding"`
	TotalPaid        float64                          `json:"total_paid"`
	Invoices         []ClientInvoice                  `json:"invoices"`
	ByChargeType     map[ChargeType]OutstandingByType `json:"by_charge_type"`
}

// OutstandingByType represents outstanding amount by charge type
type OutstandingByType struct {
	ChargeType  ChargeType `json:"charge_type"`
	Total       float64    `json:"total"`
	Paid        float64    `json:"paid"`
	Outstanding float64    `json:"outstanding"`
}

// InvoicePaymentStatusResponse represents payment status for invoice
type InvoicePaymentStatusResponse struct {
	InvoiceID         string          `json:"invoice_id"`
	InvoiceNumber     string          `json:"invoice_number"`
	ClientName        string          `json:"client_name"`
	Amount            float64         `json:"amount"`
	AmountPaid        float64         `json:"amount_paid"`
	OutstandingAmount float64         `json:"outstanding_amount"`
	Status            string          `json:"status"`
	DueDate           time.Time       `json:"due_date"`
	IsOverdue         bool            `json:"is_overdue"`
	PaymentHistory    []ClientPayment `json:"payment_history"`
}

// TenantAccountRequest represents tenant account creation request
type TenantAccountRequest struct {
	ChargeType      ChargeType `json:"charge_type" binding:"required"`
	ChargeTypeName  string     `json:"charge_type_name" binding:"required"`
	Description     string     `json:"description"`
	BankAccountName string     `json:"bank_account_name"`
	BankAccountNo   string     `json:"bank_account_no"`
	IFSCCode        string     `json:"ifsc_code"`
}

// CreateInvoiceRequest represents invoice creation request
type CreateInvoiceRequest struct {
	ClientID    string                 `json:"client_id" binding:"required"`
	ClientName  string                 `json:"client_name" binding:"required"`
	ClientEmail string                 `json:"client_email" binding:"required,email"`
	ClientPhone string                 `json:"client_phone" binding:"required"`
	ChargeType  ChargeType             `json:"charge_type" binding:"required"`
	Amount      float64                `json:"amount" binding:"required,gt=0"`
	Description string                 `json:"description"`
	DueDate     time.Time              `json:"due_date"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// BulkInvoiceRequest represents bulk invoice creation
type BulkInvoiceRequest struct {
	ChargeType  ChargeType        `json:"charge_type" binding:"required"`
	Description string            `json:"description"`
	DueDate     time.Time         `json:"due_date"`
	Invoices    []BulkInvoiceItem `json:"invoices" binding:"required"`
}

// BulkInvoiceItem represents single item in bulk invoice
type BulkInvoiceItem struct {
	ClientID    string  `json:"client_id" binding:"required"`
	ClientName  string  `json:"client_name" binding:"required"`
	ClientEmail string  `json:"client_email" binding:"required,email"`
	ClientPhone string  `json:"client_phone" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
}
