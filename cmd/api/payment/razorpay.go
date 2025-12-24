package payment

import (
	"errors"
	"fmt"
	"time"

	"github.com/razorpay/razorpay-go/requests"
)

// RazorpayClient handles Razorpay payments
type RazorpayClient struct {
	KeyID     string
	KeySecret string
	Client    *razorpay.Client
}

// NewRazorpayClient creates a new Razorpay client
func NewRazorpayClient(keyID, keySecret string) (*RazorpayClient, error) {
	if keyID == "" || keySecret == "" {
		return nil, errors.New("razorpay key ID and secret required")
	}

	client := razorpay.NewClient(keyID, keySecret)
	return &RazorpayClient{
		KeyID:     keyID,
		KeySecret: keySecret,
		Client:    client,
	}, nil
}

// CreateOrder creates a Razorpay order
func (rc *RazorpayClient) CreateOrder(amount float64, currency, receipt, description string) (*RazorpayOrder, error) {
	// Amount in paise (multiply by 100 for INR)
	amountInPaise := int64(amount * 100)

	orderRequest := requests.OrderRequest{
		Amount:      amountInPaise,
		Currency:    currency,
		Receipt:     receipt,
		Description: description,
		Notes: map[string]interface{}{
			"created_at": time.Now().Unix(),
		},
	}

	order, err := rc.Client.Order.Create(orderRequest, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create razorpay order: %w", err)
	}

	return &RazorpayOrder{
		ID:     order.ID,
		Amount: float64(order.Amount) / 100,
		Status: order.Status,
	}, nil
}

// VerifyPaymentSignature verifies Razorpay payment signature
func (rc *RazorpayClient) VerifyPaymentSignature(orderID, paymentID, signature string) (bool, error) {
	attributes := map[string]interface{}{
		"order_id":   orderID,
		"payment_id": paymentID,
	}

	err := rc.Client.Payment.ValidateSignature(attributes, signature)
	return err == nil, err
}

// FetchOrder fetches order details
func (rc *RazorpayClient) FetchOrder(orderID string) (*RazorpayOrder, error) {
	order, err := rc.Client.Order.Get(orderID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	return &RazorpayOrder{
		ID:     order.ID,
		Amount: float64(order.Amount) / 100,
		Status: order.Status,
	}, nil
}

// FetchPayment fetches payment details
func (rc *RazorpayClient) FetchPayment(paymentID string) (*RazorpayPayment, error) {
	payment, err := rc.Client.Payment.Get(paymentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment: %w", err)
	}

	return &RazorpayPayment{
		ID:               payment.ID,
		OrderID:          payment.OrderID,
		Amount:           float64(payment.Amount) / 100,
		Status:           payment.Status,
		Method:           payment.Method,
		Acquired:         payment.Acquired,
		Description:      payment.Description,
		Email:            payment.Email,
		Contact:          payment.Contact,
		Fee:              int64(payment.Fee),
		Tax:              int64(payment.Tax),
		ErrorCode:        payment.ErrorCode,
		ErrorDescription: payment.ErrorDescription,
		VPA:              payment.VPA,
		CreatedAt:        payment.CreatedAt,
	}, nil
}

// CreateRefund creates a refund
func (rc *RazorpayClient) CreateRefund(paymentID string, amount float64, notes string) (*RazorpayRefund, error) {
	amountInPaise := int64(amount * 100)

	refundRequest := requests.RefundRequest{
		Amount: &amountInPaise,
		Notes: map[string]interface{}{
			"reason": notes,
		},
	}

	refund, err := rc.Client.Payment.Refund(paymentID, refundRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	return &RazorpayRefund{
		ID:        refund.ID,
		PaymentID: refund.PaymentID,
		Amount:    float64(refund.Amount) / 100,
		Status:    refund.Status,
		CreatedAt: refund.CreatedAt,
	}, nil
}

// RazorpayOrder represents Razorpay order
type RazorpayOrder struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

// RazorpayPayment represents Razorpay payment
type RazorpayPayment struct {
	ID               string  `json:"id"`
	OrderID          string  `json:"order_id"`
	Amount           float64 `json:"amount"`
	Status           string  `json:"status"`
	Method           string  `json:"method"`
	Acquired         bool    `json:"acquired"`
	Description      string  `json:"description"`
	Email            string  `json:"email"`
	Contact          string  `json:"contact"`
	Fee              int64   `json:"fee"`
	Tax              int64   `json:"tax"`
	ErrorCode        string  `json:"error_code"`
	ErrorDescription string  `json:"error_description"`
	VPA              string  `json:"vpa"`
	CreatedAt        int64   `json:"created_at"`
}

// RazorpayRefund represents Razorpay refund
type RazorpayRefund struct {
	ID        string  `json:"id"`
	PaymentID string  `json:"payment_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	CreatedAt int64   `json:"created_at"`
}

// RazorpayWebhook represents Razorpay webhook payload
type RazorpayWebhook struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

// RazorpayPaymentData represents payment data in webhook
type RazorpayPaymentData struct {
	Payment map[string]interface{} `json:"payment"`
	Order   map[string]interface{} `json:"order"`
}
