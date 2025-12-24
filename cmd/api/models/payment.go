package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// PaymentProvider defines supported payment providers
type PaymentProvider string

const (
	ProviderRazorpay PaymentProvider = "razorpay"
	ProviderBilldesk PaymentProvider = "billdesk"
)

// PaymentMethod defines payment methods
type PaymentMethod string

const (
	MethodNetbanking PaymentMethod = "netbanking"
	MethodCreditCard PaymentMethod = "credit_card"
	MethodDebitCard  PaymentMethod = "debit_card"
	MethodUPI        PaymentMethod = "upi"
	MethodWallet     PaymentMethod = "wallet"
)

// PaymentStatus defines payment statuses
type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "pending"
	PaymentInitiated  PaymentStatus = "initiated"
	PaymentProcessing PaymentStatus = "processing"
	PaymentSuccessful PaymentStatus = "successful"
	PaymentFailed     PaymentStatus = "failed"
	PaymentCancelled  PaymentStatus = "cancelled"
	PaymentRefunded   PaymentStatus = "refunded"
)

// Payment represents a payment transaction
type Payment struct {
	ID               string          `db:"id" json:"id"`
	TenantID         string          `db:"tenant_id" json:"tenant_id"`
	OrderID          string          `db:"order_id" json:"order_id"`
	Amount           float64         `db:"amount" json:"amount"`
	Currency         string          `db:"currency" json:"currency"`
	Status           PaymentStatus   `db:"status" json:"status"`
	Provider         PaymentProvider `db:"provider" json:"provider"`
	PaymentMethod    PaymentMethod   `db:"payment_method" json:"payment_method"`
	TransactionID    string          `db:"transaction_id" json:"transaction_id"`
	GatewayOrderID   string          `db:"gateway_order_id" json:"gateway_order_id"`
	GatewayPaymentID string          `db:"gateway_payment_id" json:"gateway_payment_id"`
	Description      string          `db:"description" json:"description"`
	CustomerName     string          `db:"customer_name" json:"customer_name"`
	CustomerEmail    string          `db:"customer_email" json:"customer_email"`
	CustomerPhone    string          `db:"customer_phone" json:"customer_phone"`
	BillingAddress   *Address        `db:"billing_address" json:"billing_address"`
	Metadata         json.RawMessage `db:"metadata" json:"metadata"`
	ErrorMessage     string          `db:"error_message" json:"error_message"`
	ErrorCode        string          `db:"error_code" json:"error_code"`
	ReceiptURL       string          `db:"receipt_url" json:"receipt_url"`
	RefundID         string          `db:"refund_id" json:"refund_id"`
	RefundAmount     float64         `db:"refund_amount" json:"refund_amount"`
	CreatedAt        time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at" json:"updated_at"`
	ProcessedAt      *time.Time      `db:"processed_at" json:"processed_at"`
	ExpiresAt        *time.Time      `db:"expires_at" json:"expires_at"`
}

// Address represents billing/shipping address
type Address struct {
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phone_number"`
}

// Scan implements sql.Scanner interface
func (a *Address) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion failed")
	}
	return json.Unmarshal(bytes, &a)
}

// Value implements driver.Valuer interface
func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	Amount         float64                `json:"amount" binding:"required,gt=0"`
	Currency       string                 `json:"currency" binding:"required" default:"INR"`
	Provider       PaymentProvider        `json:"provider" binding:"required,oneof=razorpay billdesk"`
	PaymentMethod  PaymentMethod          `json:"payment_method" binding:"required,oneof=netbanking credit_card debit_card upi wallet"`
	Description    string                 `json:"description" binding:"required"`
	CustomerName   string                 `json:"customer_name" binding:"required"`
	CustomerEmail  string                 `json:"customer_email" binding:"required,email"`
	CustomerPhone  string                 `json:"customer_phone" binding:"required"`
	BillingAddress *Address               `json:"billing_address"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	PaymentID      string          `json:"payment_id"`
	OrderID        string          `json:"order_id"`
	Amount         float64         `json:"amount"`
	Currency       string          `json:"currency"`
	Status         PaymentStatus   `json:"status"`
	Provider       PaymentProvider `json:"provider"`
	PaymentMethod  PaymentMethod   `json:"payment_method"`
	GatewayOrderID string          `json:"gateway_order_id"`
	PaymentURL     string          `json:"payment_url"`
	ExpiresAt      *time.Time      `json:"expires_at"`
	CreatedAt      time.Time       `json:"created_at"`
}

// WebhookPayload represents webhook payload from payment gateway
type WebhookPayload struct {
	EventID   string          `json:"event_id"`
	EventType string          `json:"event"`
	CreatedAt int64           `json:"created_at"`
	Provider  PaymentProvider `json:"provider"`
	Data      json.RawMessage `json:"data"`
	Signature string          `json:"signature"`
}

// Refund represents a refund transaction
type Refund struct {
	ID              string        `db:"id" json:"id"`
	PaymentID       string        `db:"payment_id" json:"payment_id"`
	Amount          float64       `db:"amount" json:"amount"`
	Status          PaymentStatus `db:"status" json:"status"`
	GatewayRefundID string        `db:"gateway_refund_id" json:"gateway_refund_id"`
	Reason          string        `db:"reason" json:"reason"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at" json:"updated_at"`
	ProcessedAt     *time.Time    `db:"processed_at" json:"processed_at"`
}

// RefundRequest represents a refund request
type RefundRequest struct {
	Amount string `json:"amount" binding:"required"`
	Reason string `json:"reason" binding:"required"`
	Notes  string `json:"notes"`
}

// PaymentGatewayConfig represents payment gateway configuration
type PaymentGatewayConfig struct {
	ID        string          `db:"id" json:"id"`
	TenantID  string          `db:"tenant_id" json:"tenant_id"`
	Provider  PaymentProvider `db:"provider" json:"provider"`
	IsActive  bool            `db:"is_active" json:"is_active"`
	ApiKey    string          `db:"api_key" json:"api_key"`
	ApiSecret string          `db:"api_secret" json:"-"`
	Settings  json.RawMessage `db:"settings" json:"settings"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

// Bank represents bank details
type Bank struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// AvailableBanks returns list of banks supporting netbanking
func AvailableBanks() []Bank {
	return []Bank{
		{Code: "HDFC", Name: "HDFC Bank"},
		{Code: "ICIC", Name: "ICICI Bank"},
		{Code: "AXIS", Name: "Axis Bank"},
		{Code: "SBIN", Name: "State Bank of India"},
		{Code: "UTIB", Name: "Axis Bank"},
		{Code: "KOBA", Name: "Kotak Mahindra Bank"},
		{Code: "IDFB", Name: "IDFC Bank"},
		{Code: "AUBL", Name: "Aurobindo Bank"},
	}
}
