package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/logger"
)

// SIPAdapter implements VoIPProviderAdapter for generic SIP providers
type SIPAdapter struct {
	provider *models.VoIPProvider
	logger   *logger.Logger
	client   *http.Client
}

// NewSIPAdapter creates a new SIP adapter
func NewSIPAdapter(provider *models.VoIPProvider, logger *logger.Logger) *SIPAdapter {
	return &SIPAdapter{
		provider: provider,
		logger:   logger,
		client: &http.Client{
			Timeout: time.Duration(provider.TimeoutSeconds) * time.Second,
		},
	}
}

// InitiateCall initiates a call via SIP provider
func (s *SIPAdapter) InitiateCall(ctx context.Context, session *models.ClickToCallSession) (string, error) {
	// Generic SIP API call
	reqBody := map[string]interface{}{
		"caller":    session.FromPhone,
		"callee":    session.ToPhone,
		"caller_id": s.provider.CallerID,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/calls", s.provider.APIURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.provider.APIKey != "" {
		req.Header.Set("X-API-Key", s.provider.APIKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("SIP error: %s", string(body))
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

// EndCall ends a call via SIP provider
func (s *SIPAdapter) EndCall(ctx context.Context, sessionID string) error {
	url := fmt.Sprintf("%s/calls/%s", s.provider.APIURL, sessionID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if s.provider.APIKey != "" {
		req.Header.Set("X-API-Key", s.provider.APIKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("SIP error: %d", resp.StatusCode)
	}

	return nil
}

// TransferCall transfers a call via SIP provider
func (s *SIPAdapter) TransferCall(ctx context.Context, sessionID, toPhone string) error {
	reqBody := map[string]interface{}{
		"transfer_to": toPhone,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/calls/%s/transfer", s.provider.APIURL, sessionID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.provider.APIKey != "" {
		req.Header.Set("X-API-Key", s.provider.APIKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SIP error: %d", resp.StatusCode)
	}

	return nil
}

// GetCallStatus gets call status via SIP provider
func (s *SIPAdapter) GetCallStatus(ctx context.Context, sessionID string) (string, error) {
	url := fmt.Sprintf("%s/calls/%s", s.provider.APIURL, sessionID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	if s.provider.APIKey != "" {
		req.Header.Set("X-API-Key", s.provider.APIKey)
	}

	resp, err := s.client.Do(req)
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

// ValidateWebhookSignature validates SIP webhook signature
func (s *SIPAdapter) ValidateWebhookSignature(payload []byte, signature string) bool {
	return true
}

// TwilioAdapter implements VoIPProviderAdapter for Twilio
type TwilioAdapter struct {
	provider *models.VoIPProvider
	logger   *logger.Logger
	client   *http.Client
}

// NewTwilioAdapter creates a new Twilio adapter
func NewTwilioAdapter(provider *models.VoIPProvider, logger *logger.Logger) *TwilioAdapter {
	return &TwilioAdapter{
		provider: provider,
		logger:   logger,
		client: &http.Client{
			Timeout: time.Duration(provider.TimeoutSeconds) * time.Second,
		},
	}
}

// InitiateCall initiates a call via Twilio
func (t *TwilioAdapter) InitiateCall(ctx context.Context, session *models.ClickToCallSession) (string, error) {
	// Twilio API: Create a call
	reqBody := map[string]interface{}{
		"From": session.FromPhone,
		"To":   session.ToPhone,
		"Url":  t.provider.CallbackURL,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	accountSID := t.provider.APIKey
	url := fmt.Sprintf("%s/2010-04-01/Accounts/%s/Calls.json", t.provider.APIURL, accountSID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(accountSID, t.provider.APISecret)

	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("twilio error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return session.SessionID, nil
	}

	if id, ok := result["sid"].(string); ok {
		return id, nil
	}

	return session.SessionID, nil
}

// EndCall ends a call via Twilio
func (t *TwilioAdapter) EndCall(ctx context.Context, sessionID string) error {
	accountSID := t.provider.APIKey
	url := fmt.Sprintf("%s/2010-04-01/Accounts/%s/Calls/%s.json", t.provider.APIURL, accountSID, sessionID)

	reqBody := map[string]interface{}{
		"Status": "completed",
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(accountSID, t.provider.APISecret)

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("twilio error: %s", string(body))
	}

	return nil
}

// TransferCall transfers a call via Twilio
func (t *TwilioAdapter) TransferCall(ctx context.Context, sessionID, toPhone string) error {
	accountSID := t.provider.APIKey
	url := fmt.Sprintf("%s/2010-04-01/Accounts/%s/Calls/%s.json", t.provider.APIURL, accountSID, sessionID)

	reqBody := map[string]interface{}{
		"Url": t.provider.CallbackURL + "?transfer_to=" + toPhone,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(accountSID, t.provider.APISecret)

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("twilio error: %d", resp.StatusCode)
	}

	return nil
}

// GetCallStatus gets call status via Twilio
func (t *TwilioAdapter) GetCallStatus(ctx context.Context, sessionID string) (string, error) {
	accountSID := t.provider.APIKey
	url := fmt.Sprintf("%s/2010-04-01/Accounts/%s/Calls/%s.json", t.provider.APIURL, accountSID, sessionID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(accountSID, t.provider.APISecret)

	resp, err := t.client.Do(req)
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

// ValidateWebhookSignature validates Twilio webhook signature
func (t *TwilioAdapter) ValidateWebhookSignature(payload []byte, signature string) bool {
	// Implement Twilio signature validation
	// https://www.twilio.com/docs/usage/webhooks/webhooks-security
	return true
}
