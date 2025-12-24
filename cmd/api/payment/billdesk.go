package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// BilldeskClient handles Billdesk payments
type BilldeskClient struct {
	APIKey      string
	APISecret   string
	MerchantID  string
	ClientID    string
	Environment string // "sandbox" or "production"
}

// NewBilldeskClient creates a new Billdesk client
func NewBilldeskClient(apiKey, apiSecret, merchantID, clientID, environment string) (*BilldeskClient, error) {
	if apiKey == "" || apiSecret == "" || merchantID == "" {
		return nil, errors.New("billdesk api key, secret, and merchant id required")
	}

	return &BilldeskClient{
		APIKey:      apiKey,
		APISecret:   apiSecret,
		MerchantID:  merchantID,
		ClientID:    clientID,
		Environment: environment,
	}, nil
}

// GetBaseURL returns the base URL for Billdesk API
func (bc *BilldeskClient) GetBaseURL() string {
	if bc.Environment == "production" {
		return "https://api.billdesk.com"
	}
	return "https://sandbox.billdesk.com"
}

// CreateOrder creates a Billdesk order
func (bc *BilldeskClient) CreateOrder(amount float64, currency, orderID, description string, paymentMethods []string) (*BilldeskOrder, error) {
	url := fmt.Sprintf("%s/v1/orders", bc.GetBaseURL())

	// Amount in paise
	amountInPaise := int64(amount * 100)

	requestBody := map[string]interface{}{
		"merchant_id": bc.MerchantID,
		"customer_id": "",
		"order_id":    orderID,
		"amount":      amountInPaise,
		"currency":    currency,
		"description": description,
		"ip_address":  "127.0.0.1",
		"device_type": "web",
		"redirect_url": map[string]interface{}{
			"success_url": "https://yoursite.com/payment/success",
			"failure_url": "https://yoursite.com/payment/failure",
			"neutral_url": "https://yoursite.com/payment/neutral",
		},
		"payment_options": map[string]interface{}{
			"wallets":    contains(paymentMethods, "wallet"),
			"netbanking": contains(paymentMethods, "netbanking"),
			"cards":      contains(paymentMethods, "credit_card") || contains(paymentMethods, "debit_card"),
			"upi":        contains(paymentMethods, "upi"),
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request with authentication
	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	timestamp := time.Now().Unix()
	signature := bc.generateSignature(string(body), timestamp)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bc.APIKey))
	req.Header.Set("X-Billdesk-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Billdesk-Signature", signature)

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("api error: %s", string(respBody))
	}

	var response BilldeskOrderResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &BilldeskOrder{
		OrderID:    response.OrderID,
		BdOrderID:  response.BdOrderID,
		Amount:     amount,
		Status:     response.Status,
		PaymentURL: response.PaymentURL,
		ExpiresAt:  response.ExpiresAt,
	}, nil
}

// VerifyPaymentSignature verifies Billdesk payment signature
func (bc *BilldeskClient) VerifyPaymentSignature(signature string, body string) bool {
	expectedSignature := bc.generateSignature(body, time.Now().Unix())
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// FetchOrder fetches Billdesk order details
func (bc *BilldeskClient) FetchOrder(orderID string) (*BilldeskOrder, error) {
	url := fmt.Sprintf("%s/v1/orders/%s", bc.GetBaseURL(), orderID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	timestamp := time.Now().Unix()
	signature := bc.generateSignature("", timestamp)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bc.APIKey))
	req.Header.Set("X-Billdesk-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Billdesk-Signature", signature)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error: %s", string(respBody))
	}

	var response BilldeskOrderResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &BilldeskOrder{
		OrderID:    response.OrderID,
		BdOrderID:  response.BdOrderID,
		Amount:     float64(response.Amount) / 100,
		Status:     response.Status,
		PaymentURL: response.PaymentURL,
		ExpiresAt:  response.ExpiresAt,
	}, nil
}

// CreateRefund creates a refund
func (bc *BilldeskClient) CreateRefund(orderID string, amount float64, reason string) (*BilldeskRefund, error) {
	url := fmt.Sprintf("%s/v1/orders/%s/refund", bc.GetBaseURL(), orderID)

	amountInPaise := int64(amount * 100)

	requestBody := map[string]interface{}{
		"amount": amountInPaise,
		"reason": reason,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	timestamp := time.Now().Unix()
	signature := bc.generateSignature(string(body), timestamp)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bc.APIKey))
	req.Header.Set("X-Billdesk-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Billdesk-Signature", signature)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("api error: %s", string(respBody))
	}

	var response BilldeskRefundResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &BilldeskRefund{
		RefundID:  response.RefundID,
		OrderID:   response.OrderID,
		Amount:    amount,
		Status:    response.Status,
		CreatedAt: response.CreatedAt,
	}, nil
}

// generateSignature generates HMAC-SHA256 signature for Billdesk
func (bc *BilldeskClient) generateSignature(body string, timestamp int64) string {
	message := fmt.Sprintf("%s|%d", body, timestamp)
	h := hmac.New(sha256.New, []byte(bc.APISecret))
	h.Write([]byte(message))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// BilldeskOrder represents Billdesk order
type BilldeskOrder struct {
	OrderID    string     `json:"order_id"`
	BdOrderID  string     `json:"bd_order_id"`
	Amount     float64    `json:"amount"`
	Status     string     `json:"status"`
	PaymentURL string     `json:"payment_url"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

// BilldeskOrderResponse represents Billdesk API response
type BilldeskOrderResponse struct {
	OrderID    string     `json:"order_id"`
	BdOrderID  string     `json:"bd_order_id"`
	Amount     int64      `json:"amount"`
	Status     string     `json:"status"`
	PaymentURL string     `json:"payment_url"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

// BilldeskRefund represents Billdesk refund
type BilldeskRefund struct {
	RefundID  string    `json:"refund_id"`
	OrderID   string    `json:"order_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// BilldeskRefundResponse represents Billdesk refund response
type BilldeskRefundResponse struct {
	RefundID  string    `json:"refund_id"`
	OrderID   string    `json:"order_id"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// BilldeskWebhook represents Billdesk webhook
type BilldeskWebhook struct {
	OrderID       string `json:"order_id"`
	BdOrderID     string `json:"bd_order_id"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Amount        int64  `json:"amount"`
	Timestamp     int64  `json:"timestamp"`
}

// Helper function
func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
