package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// CommunicationProvider represents different communication providers
type CommunicationProviderType string

const (
	ProviderTypeSMS       CommunicationProviderType = "sms"
	ProviderTypeEmail     CommunicationProviderType = "email"
	ProviderTypeWhatsApp  CommunicationProviderType = "whatsapp"
	ProviderTypeSlack     CommunicationProviderType = "slack"
	ProviderTypePushNotif CommunicationProviderType = "push_notification"
)

// ProviderCredential stores API credentials for communication providers
type ProviderCredential struct {
	ID           int64                     `json:"id"`
	TenantID     string                    `json:"tenant_id"`
	ProviderType CommunicationProviderType `json:"provider_type"`
	APIKey       string                    `json:"api_key"`
	APISecret    string                    `json:"api_secret,omitempty"`
	BaseURL      string                    `json:"base_url,omitempty"`
	Config       map[string]interface{}    `json:"config,omitempty"`
	IsActive     bool                      `json:"is_active"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

// MessageTemplate defines a reusable message template
type MessageTemplate struct {
	ID           int64                     `json:"id"`
	TenantID     string                    `json:"tenant_id"`
	Name         string                    `json:"name"`
	Category     string                    `json:"category"` // welcome, reminder, follow_up, promotion
	ProviderType CommunicationProviderType `json:"provider_type"`
	Subject      string                    `json:"subject,omitempty"`
	Body         string                    `json:"body"`
	Variables    []string                  `json:"variables,omitempty"` // {{name}}, {{email}}, etc
	IsActive     bool                      `json:"is_active"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

// Message represents a message to be sent
type Message struct {
	ID           int64                     `json:"id"`
	TenantID     string                    `json:"tenant_id"`
	Recipient    string                    `json:"recipient"` // email, phone, etc
	ProviderType CommunicationProviderType `json:"provider_type"`
	TemplateID   *int64                    `json:"template_id,omitempty"`
	Subject      string                    `json:"subject,omitempty"`
	Body         string                    `json:"body"`
	Status       string                    `json:"status"` // pending, sent, failed, bounced
	ProviderID   string                    `json:"provider_id,omitempty"`
	Error        string                    `json:"error,omitempty"`
	SentAt       *time.Time                `json:"sent_at,omitempty"`
	DeliveredAt  *time.Time                `json:"delivered_at,omitempty"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

// CommunicationService handles all communication with external providers
type CommunicationService struct {
	db        *sql.DB
	providers map[string]ProviderCredential
}

// NewCommunicationService creates a new CommunicationService
func NewCommunicationService(db *sql.DB) *CommunicationService {
	return &CommunicationService{
		db:        db,
		providers: make(map[string]ProviderCredential),
	}
}

// RegisterProvider registers a new communication provider for a tenant
func (cs *CommunicationService) RegisterProvider(ctx context.Context, cred *ProviderCredential) error {
	query := `
		INSERT INTO communication_provider (tenant_id, provider_type, api_key, api_secret, base_url, config, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := cs.db.ExecContext(ctx, query,
		cred.TenantID, cred.ProviderType, cred.APIKey, cred.APISecret,
		cred.BaseURL, fmt.Sprintf("%v", cred.Config), cred.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	cred.ID = id

	// Cache the provider
	key := fmt.Sprintf("%s-%s", cred.TenantID, cred.ProviderType)
	cs.providers[key] = *cred

	return err
}

// GetProvider retrieves a provider credential
func (cs *CommunicationService) GetProvider(ctx context.Context, tenantID string, providerType CommunicationProviderType) (*ProviderCredential, error) {
	key := fmt.Sprintf("%s-%s", tenantID, providerType)

	// Check cache first
	if cred, ok := cs.providers[key]; ok {
		return &cred, nil
	}

	// Query database
	query := `
		SELECT id, tenant_id, provider_type, api_key, api_secret, base_url, config, is_active, created_at, updated_at
		FROM communication_provider
		WHERE tenant_id = ? AND provider_type = ? AND is_active = true
		LIMIT 1
	`

	cred := &ProviderCredential{}
	var configStr string

	err := cs.db.QueryRowContext(ctx, query, tenantID, providerType).
		Scan(&cred.ID, &cred.TenantID, &cred.ProviderType, &cred.APIKey, &cred.APISecret,
			&cred.BaseURL, &configStr, &cred.IsActive, &cred.CreatedAt, &cred.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Cache it
	cs.providers[key] = *cred
	return cred, nil
}

// CreateTemplate creates a message template
func (cs *CommunicationService) CreateTemplate(ctx context.Context, template *MessageTemplate) error {
	query := `
		INSERT INTO message_template (tenant_id, name, category, provider_type, subject, body, variables, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := cs.db.ExecContext(ctx, query,
		template.TenantID, template.Name, template.Category, template.ProviderType,
		template.Subject, template.Body, fmt.Sprintf("%v", template.Variables), template.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	template.ID = id
	return err
}

// GetTemplate retrieves a message template
func (cs *CommunicationService) GetTemplate(ctx context.Context, tenantID string, templateID int64) (*MessageTemplate, error) {
	query := `
		SELECT id, tenant_id, name, category, provider_type, subject, body, variables, is_active, created_at, updated_at
		FROM message_template
		WHERE id = ? AND tenant_id = ?
	`

	template := &MessageTemplate{}
	var variablesStr string

	err := cs.db.QueryRowContext(ctx, query, templateID, tenantID).
		Scan(&template.ID, &template.TenantID, &template.Name, &template.Category, &template.ProviderType,
			&template.Subject, &template.Body, &variablesStr, &template.IsActive, &template.CreatedAt, &template.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return template, nil
}

// SendMessage sends a message through the appropriate provider
func (cs *CommunicationService) SendMessage(ctx context.Context, msg *Message) error {
	// Get provider credentials
	provider, err := cs.GetProvider(ctx, msg.TenantID, msg.ProviderType)
	if err != nil {
		msg.Status = "failed"
		msg.Error = fmt.Sprintf("provider not configured: %v", err)
		return err
	}

	if !provider.IsActive {
		msg.Status = "failed"
		msg.Error = "provider is inactive"
		return fmt.Errorf("provider is inactive")
	}

	// Save message record
	query := `
		INSERT INTO message (tenant_id, recipient, provider_type, template_id, subject, body, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := cs.db.ExecContext(ctx, query,
		msg.TenantID, msg.Recipient, msg.ProviderType, msg.TemplateID,
		msg.Subject, msg.Body, "pending")
	if err != nil {
		return err
	}

	msgID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	msg.ID = msgID

	// Send through provider (in real implementation, call actual provider APIs)
	err = cs.sendViaProvider(ctx, provider, msg)
	if err != nil {
		msg.Status = "failed"
		msg.Error = err.Error()
		cs.updateMessageStatus(ctx, msgID, "failed", err.Error())
		return err
	}

	msg.Status = "sent"
	msg.SentAt = timePtr(time.Now())
	cs.updateMessageStatus(ctx, msgID, "sent", "")

	return nil
}

// sendViaProvider sends message through the specific provider
func (cs *CommunicationService) sendViaProvider(ctx context.Context, provider *ProviderCredential, msg *Message) error {
	// Simulate sending - in real implementation, call actual provider APIs
	switch provider.ProviderType {
	case ProviderTypeSMS:
		return cs.sendSMS(ctx, provider, msg)
	case ProviderTypeEmail:
		return cs.sendEmail(ctx, provider, msg)
	case ProviderTypeWhatsApp:
		return cs.sendWhatsApp(ctx, provider, msg)
	case ProviderTypeSlack:
		return cs.sendSlack(ctx, provider, msg)
	default:
		return fmt.Errorf("unsupported provider type: %s", provider.ProviderType)
	}
}

// sendSMS simulates sending an SMS (would integrate with Twilio, etc)
func (cs *CommunicationService) sendSMS(_ context.Context, _ *ProviderCredential, _ *Message) error {
	// In real implementation:
	// - Use Twilio SDK or similar
	// - msg.Recipient would be a phone number
	// - Return provider_id from SMS service
	return nil
}

// sendEmail simulates sending an email (would integrate with SendGrid, etc)
func (cs *CommunicationService) sendEmail(_ context.Context, _ *ProviderCredential, _ *Message) error {
	// In real implementation:
	// - Use SendGrid SDK or similar
	// - msg.Recipient would be an email address
	// - msg.Subject and msg.Body contain email content
	// - Return provider_id from email service
	return nil
}

// sendWhatsApp simulates sending a WhatsApp message
func (cs *CommunicationService) sendWhatsApp(_ context.Context, _ *ProviderCredential, _ *Message) error {
	// In real implementation:
	// - Use WhatsApp Business API
	// - msg.Recipient would be a phone number
	// - Return provider_id from WhatsApp service
	return nil
}

// sendSlack simulates sending a Slack message
func (cs *CommunicationService) sendSlack(_ context.Context, _ *ProviderCredential, _ *Message) error {
	// In real implementation:
	// - Use Slack SDK
	// - msg.Recipient would be a Slack user ID or channel
	// - Return provider_id from Slack service
	return nil
}

// updateMessageStatus updates the status of a sent message
func (cs *CommunicationService) updateMessageStatus(ctx context.Context, msgID int64, status string, errorMsg string) error {
	query := `
		UPDATE message SET status = ?, error = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := cs.db.ExecContext(ctx, query, status, errorMsg, msgID)
	return err
}

// GetMessageStatus retrieves the status of a sent message
func (cs *CommunicationService) GetMessageStatus(ctx context.Context, tenantID string, msgID int64) (*Message, error) {
	query := `
		SELECT id, tenant_id, recipient, provider_type, template_id, subject, body, status, provider_id, error, sent_at, delivered_at, created_at, updated_at
		FROM message
		WHERE id = ? AND tenant_id = ?
	`

	msg := &Message{}
	var templateID sql.NullInt64
	var subject, providerID, error_ sql.NullString
	var sentAt, deliveredAt sql.NullTime

	err := cs.db.QueryRowContext(ctx, query, msgID, tenantID).
		Scan(&msg.ID, &msg.TenantID, &msg.Recipient, &msg.ProviderType, &templateID,
			&subject, &msg.Body, &msg.Status, &providerID, &error_, &sentAt, &deliveredAt,
			&msg.CreatedAt, &msg.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if templateID.Valid {
		msg.TemplateID = &templateID.Int64
	}
	if subject.Valid {
		msg.Subject = subject.String
	}
	if providerID.Valid {
		msg.ProviderID = providerID.String
	}
	if error_.Valid {
		msg.Error = error_.String
	}
	if sentAt.Valid {
		msg.SentAt = &sentAt.Time
	}
	if deliveredAt.Valid {
		msg.DeliveredAt = &deliveredAt.Time
	}

	return msg, nil
}

// GetMessageStats retrieves communication statistics
func (cs *CommunicationService) GetMessageStats(ctx context.Context, tenantID string, startDate, endDate time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total messages by status
	query := `
		SELECT provider_type, status, COUNT(*) as count
		FROM message
		WHERE tenant_id = ? AND created_at BETWEEN ? AND ?
		GROUP BY provider_type, status
	`

	rows, err := cs.db.QueryContext(ctx, query, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providerStats := make(map[string]map[string]int)
	for rows.Next() {
		var provider, status string
		var count int
		if err := rows.Scan(&provider, &status, &count); err != nil {
			continue
		}

		if providerStats[provider] == nil {
			providerStats[provider] = make(map[string]int)
		}
		providerStats[provider][status] = count
	}

	stats["by_provider"] = providerStats

	return stats, nil
}

// timePtr is a helper function to get a pointer to a time
func timePtr(t time.Time) *time.Time {
	return &t
}
