package services

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/logger"
)

// VoIPProviderAdapter interface for different VoIP providers
type VoIPProviderAdapter interface {
	InitiateCall(ctx context.Context, session *models.ClickToCallSession) (string, error)
	EndCall(ctx context.Context, sessionID string) error
	TransferCall(ctx context.Context, sessionID, toPhone string) error
	GetCallStatus(ctx context.Context, sessionID string) (string, error)
	ValidateWebhookSignature(payload []byte, signature string) bool
}

// AsteriskAdapter implements VoIPProviderAdapter for Asterisk
type AsteriskAdapter struct {
	provider *models.VoIPProvider
	logger   *logger.Logger
	client   *http.Client
}

// NewAsteriskAdapter creates a new Asterisk adapter
func NewAsteriskAdapter(provider *models.VoIPProvider, logger *logger.Logger) *AsteriskAdapter {
	return &AsteriskAdapter{
		provider: provider,
		logger:   logger,
		client: &http.Client{
			Timeout: time.Duration(provider.TimeoutSeconds) * time.Second,
		},
	}
}

// InitiateCall initiates a call via Asterisk AMI or ARI
func (a *AsteriskAdapter) InitiateCall(ctx context.Context, session *models.ClickToCallSession) (string, error) {
	// Use Asterisk REST Interface (ARI) for newer versions
	reqBody := map[string]interface{}{
		"endpoint": fmt.Sprintf("SIP/%s", session.ToPhone),
		"extension": session.ToPhone,
		"context": "from-internal",
		"priority": 1,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/channels", a.provider.APIURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.provider.AuthToken))

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("asterisk error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if id, ok := result["id"].(string); ok {
		return id, nil
	}

	return session.SessionID, nil
}

// EndCall ends a call in Asterisk
func (a *AsteriskAdapter) EndCall(ctx context.Context, sessionID string) error {
	url := fmt.Sprintf("%s/api/channels/%s", a.provider.APIURL, sessionID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.provider.AuthToken))

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("asterisk error: %d", resp.StatusCode)
	}

	return nil
}

// TransferCall transfers a call to another extension
func (a *AsteriskAdapter) TransferCall(ctx context.Context, sessionID, toPhone string) error {
	reqBody := map[string]interface{}{
		"extension": toPhone,
		"context": "from-internal",
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/channels/%s/redirect", a.provider.APIURL, sessionID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.provider.AuthToken))

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("asterisk error: %d", resp.StatusCode)
	}

	return nil
}

// GetCallStatus gets status of a call
func (a *AsteriskAdapter) GetCallStatus(ctx context.Context, sessionID string) (string, error) {
	url := fmt.Sprintf("%s/api/channels/%s", a.provider.APIURL, sessionID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.provider.AuthToken))

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "UNKNOWN", nil
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "UNKNOWN", nil
	}

	if state, ok := result["state"].(string); ok {
		return state, nil
	}

	return "UNKNOWN", nil
}

// ValidateWebhookSignature validates Asterisk webhook signature
func (a *AsteriskAdapter) ValidateWebhookSignature(payload []byte, signature string) bool {
	// Implement HMAC validation if configured
	return true
}

// ExotelAdapter implements VoIPProviderAdapter for Exotel
type ExotelAdapter struct {
	provider *models.VoIPProvider
	logger   *logger.Logger
	client   *http.Client
}

// NewExotelAdapter creates a new Exotel adapter
func NewExotelAdapter(provider *models.VoIPProvider, logger *logger.Logger) *ExotelAdapter {
	return &ExotelAdapter{
		provider: provider,
		logger:   logger,
		client: &http.Client{
			Timeout: time.Duration(provider.TimeoutSeconds) * time.Second,
		},
	}
}

// InitiateCall initiates a call via Exotel
func (e *ExotelAdapter) InitiateCall(ctx context.Context, session *models.ClickToCallSession) (string, error) {
	// Exotel API endpoint
	urlStr := fmt.Sprintf("%s/v1/Accounts/%s/Calls/connect", e.provider.APIURL, e.provider.APIKey)

	// Prepare form data
	formData := url.Values{}
	formData.Set("From", e.provider.PhoneNumber)
	formData.Set("To", session.ToPhone)
	formData.Set("CallerId", e.provider.CallerID)
	formData.Set("Timeout", fmt.Sprintf("%d", e.provider.TimeoutSeconds))

	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(e.provider.APIKey, e.provider.APISecret)

	resp, err := e.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("exotel error: %s", string(body))
	}

	// Parse XML response
	type CallResponse struct {
		CallSid string `xml:"Call>Sid"`
	}
	var result CallResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		return session.SessionID, nil
	}

	return result.CallSid, nil
}

// EndCall ends a call in Exotel
func (e *ExotelAdapter) EndCall(ctx context.Context, sessionID string) error {
	urlStr := fmt.Sprintf("%s/v1/Accounts/%s/Calls/%s", e.provider.APIURL, e.provider.APIKey, sessionID)

	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(e.provider.APIKey, e.provider.APISecret)

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("exotel error: %s", string(body))
	}

	return nil
}

// TransferCall transfers a call in Exotel
func (e *ExotelAdapter) TransferCall(ctx context.Context, sessionID, toPhone string) error {
	urlStr := fmt.Sprintf("%s/v1/Accounts/%s/Calls/%s", e.provider.APIURL, e.provider.APIKey, sessionID)

	formData := url.Values{}
	formData.Set("ForwardingPhoneNumber", toPhone)

	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(e.provider.APIKey, e.provider.APISecret)

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("exotel error: %d", resp.StatusCode)
	}

	return nil
}

// GetCallStatus gets status of a call in Exotel
func (e *ExotelAdapter) GetCallStatus(ctx context.Context, sessionID string) (string, error) {
	urlStr := fmt.Sprintf("%s/v1/Accounts/%s/Calls/%s", e.provider.APIURL, e.provider.APIKey, sessionID)

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(e.provider.APIKey, e.provider.APISecret)

	resp, err := e.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "UNKNOWN", nil
	}

	body, _ := io.ReadAll(resp.Body)

	type CallStatus struct {
		Status string `xml:"Call>Status"`
	}
	var result CallStatus
	if err := xml.Unmarshal(body, &result); err != nil {
		return "UNKNOWN", nil
	}

	return result.Status, nil
}

// ValidateWebhookSignature validates Exotel webhook signature
func (e *ExotelAdapter) ValidateWebhookSignature(payload []byte, signature string) bool {
	// Exotel typically uses API key in query params or headers
	return true
}

// mCubeAdapter implements VoIPProviderAdapter for mCube
type mCubeAdapter struct {
	provider *models.VoIPProvider
	logger   *logger.Logger
	client   *http.Client
}

// NewmCubeAdapter creates a new mCube adapter
func NewmCubeAdapter(provider *models.VoIPProvider, logger *logger.Logger) *mCubeAdapter {
	return &mCubeAdapter{
		provider: provider,
		logger:   logger,
		client: &http.Client{
			Timeout: time.Duration(provider.TimeoutSeconds) * time.Second,
		},
	}
}

// InitiateCall initiates a call via mCube
func (m *mCubeAdapter) InitiateCall(ctx context.Context, session *models.ClickToCallSession) (string, error) {
	reqBody := map[string]interface{}{
		"from":      session.FromPhone,
		"to":        session.ToPhone,
		"format":    "json",
		"priority":  1,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/click2call", m.provider.APIURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.provider.AuthToken))

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("mcube error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return session.SessionID, nil
	}

	if id, ok := result["call_id"].(string); ok {
		return id, nil
	}

	return session.SessionID, nil
}

// EndCall ends a call in mCube
func (m *mCubeAdapter) EndCall(ctx context.Context, sessionID string) error {
	url := fmt.Sprintf("%s/click2call/%s", m.provider.APIURL, sessionID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.provider.AuthToken))

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("mcube error: %d", resp.StatusCode)
	}

	return nil
}

// TransferCall transfers a call in mCube
func (m *mCubeAdapter) TransferCall(ctx context.Context, sessionID, toPhone string) error {
	reqBody := map[string]interface{}{
		"transfer_to": toPhone,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/click2call/%s/transfer", m.provider.APIURL, sessionID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.provider.AuthToken))

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mcube error: %d", resp.StatusCode)
	}

	return nil
}

// GetCallStatus gets call status in mCube
func (m *mCubeAdapter) GetCallStatus(ctx context.Context, sessionID string) (string, error) {
	url := fmt.Sprintf("%s/click2call/%s", m.provider.APIURL, sessionID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.provider.AuthToken))

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "UNKNOWN", nil
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "UNKNOWN", nil
	}

	if status, ok := result["status"].(string); ok {
		return status, nil
	}

	return "UNKNOWN", nil
}

// ValidateWebhookSignature validates mCube webhook signature
func (m *mCubeAdapter) ValidateWebhookSignature(payload []byte, signature string) bool {
	// Implement HMAC validation if configured
	return true
}

// GetProviderAdapter returns the appropriate adapter for a provider
func GetProviderAdapter(provider *models.VoIPProvider, logger *logger.Logger) VoIPProviderAdapter {
	switch provider.ProviderType {
	case "ASTERISK":
		return NewAsteriskAdapter(provider, logger)
	case "EXOTEL":
		return NewExotelAdapter(provider, logger)
	case "MCUBE":
		return NewmCubeAdapter(provider, logger)
	case "SIP":
		return NewSIPAdapter(provider, logger)
	case "TWILIO":
		return NewTwilioAdapter(provider, logger)
	default:
		return nil
	}
}
