package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"vyomtech-backend/internal/models"
)

// MultiChannelCommunicationService handles communication across multiple channels
type MultiChannelCommunicationService struct {
	db     *gorm.DB
	logger Logger
	client *http.Client
}

// NewMultiChannelCommunicationService creates a new multi-channel service
func NewMultiChannelCommunicationService(db *gorm.DB, logger Logger) *MultiChannelCommunicationService {
	return &MultiChannelCommunicationService{
		db:     db,
		logger: logger,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// CreateCommunicationChannel creates a new communication channel configuration
func (s *MultiChannelCommunicationService) CreateCommunicationChannel(ctx context.Context, req *models.CreateChannelRequest, tenantID string) (*models.CommunicationChannel, error) {
	channel := &models.CommunicationChannel{
		ID:             uuid.New().String(),
		TenantID:       tenantID,
		ChannelType:    req.ChannelType,
		ChannelName:    req.ChannelName,
		APIProvider:    req.APIProvider,
		APIKey:         req.APIKey,
		APISecret:      req.APISecret,
		APIURL:         req.APIURL,
		WebhookURL:     req.WebhookURL,
		AuthToken:      req.AuthToken,
		AccountID:      req.AccountID,
		SenderID:       req.SenderID,
		ConfigJSON:     req.ConfigJSON,
		IsActive:       true,
		RetryCount:     3,
		TimeoutSeconds: 30,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(channel).Error; err != nil {
		s.logger.Error(fmt.Sprintf("Failed to create communication channel: %v", err))
		return nil, err
	}

	return channel, nil
}

// GetCommunicationChannel retrieves a channel by ID
func (s *MultiChannelCommunicationService) GetCommunicationChannel(ctx context.Context, tenantID, channelID string) (*models.CommunicationChannel, error) {
	var channel models.CommunicationChannel
	if err := s.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", channelID, tenantID).First(&channel).Error; err != nil {
		return nil, err
	}
	return &channel, nil
}

// ListCommunicationChannels lists all channels for a tenant
func (s *MultiChannelCommunicationService) ListCommunicationChannels(ctx context.Context, tenantID string, limit, offset int) ([]*models.CommunicationChannel, error) {
	var channels []*models.CommunicationChannel
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = 1", tenantID).
		Limit(limit).
		Offset(offset).
		Find(&channels).Error; err != nil {
		return nil, err
	}
	return channels, nil
}

// SendMessage sends a message via specified channel
func (s *MultiChannelCommunicationService) SendMessage(ctx context.Context, tenantID string, req *models.SendMessageRequest) (*models.CommunicationMessage, error) {
	// Get channel configuration
	var channel models.CommunicationChannel
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND channel_type = ? AND is_active = 1", tenantID, req.ChannelType).
		First(&channel).Error; err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	// Create or get session
	session, err := s.getOrCreateSession(ctx, tenantID, &channel, req)
	if err != nil {
		return nil, err
	}

	// Create message record
	message := &models.CommunicationMessage{
		ID:          uuid.New().String(),
		SessionID:   session.ID,
		TenantID:    tenantID,
		ChannelType: req.ChannelType,
		FromAddress: channel.SenderID,
		ToAddress:   req.ToAddress,
		MessageType: req.MessageType,
		MessageBody: req.MessageBody,
		Status:      "QUEUED",
		MaxRetries:  3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.Subject != "" {
		message.MessageSubject = req.Subject
	}

	// Save message to database
	if err := s.db.WithContext(ctx).Create(message).Error; err != nil {
		return nil, err
	}

	// Send message via provider
	externalID, err := s.sendViaProvider(ctx, &channel, message)
	if err != nil {
		// Update message status to failed
		s.db.WithContext(ctx).Model(message).Updates(map[string]interface{}{
			"status":        "FAILED",
			"error_message": err.Error(),
			"updated_at":    time.Now(),
		})
		return nil, err
	}

	// Update message with external ID and sent status
	now := time.Now()
	s.db.WithContext(ctx).Model(message).Updates(map[string]interface{}{
		"external_message_id": externalID,
		"status":              "SENT",
		"sent_at":             now,
		"updated_at":          now,
	})
	message.ExternalMessageID = externalID
	message.Status = "SENT"
	message.SentAt = &now

	return message, nil
}

// SendBulkMessages sends messages to multiple recipients
func (s *MultiChannelCommunicationService) SendBulkMessages(ctx context.Context, tenantID string, req *models.BulkSendRequest) (*models.BulkMessageCampaign, error) {
	campaign := &models.BulkMessageCampaign{
		ID:              uuid.New().String(),
		TenantID:        tenantID,
		CampaignName:    req.CampaignName,
		ChannelType:     req.ChannelType,
		TotalRecipients: len(req.Recipients),
		CampaignStatus:  "DRAFT",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(campaign).Error; err != nil {
		return nil, err
	}

	// Create recipients
	for _, recipient := range req.Recipients {
		bulkRecipient := &models.BulkMessageRecipient{
			ID:               uuid.New().String(),
			CampaignID:       campaign.ID,
			TenantID:         tenantID,
			RecipientAddress: recipient.Address,
			RecipientName:    recipient.Name,
			ContactID:        recipient.ContactID,
			LeadID:           recipient.LeadID,
			SendStatus:       "PENDING",
			CreatedAt:        time.Now(),
		}

		if recipient.TemplateVariables != nil {
			varBytes, _ := json.Marshal(recipient.TemplateVariables)
			bulkRecipient.TemplateVariables = varBytes
		}

		s.db.WithContext(ctx).Create(bulkRecipient)
	}

	return campaign, nil
}

// getOrCreateSession gets or creates a communication session
func (s *MultiChannelCommunicationService) getOrCreateSession(ctx context.Context, tenantID string, channel *models.CommunicationChannel, req *models.SendMessageRequest) (*models.CommunicationSession, error) {
	var session models.CommunicationSession

	// Try to find existing session
	if req.LeadID != "" {
		if err := s.db.WithContext(ctx).
			Where("tenant_id = ? AND lead_id = ? AND channel_type = ? AND status != ?",
				tenantID, req.LeadID, req.ChannelType, "COMPLETED").
			Order("updated_at DESC").
			First(&session).Error; err == nil {
			return &session, nil
		}
	}

	// Create new session
	session = models.CommunicationSession{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		ChannelID:   channel.ID,
		ChannelType: req.ChannelType,
		SenderID:    channel.SenderID,
		RecipientID: req.ToAddress,
		LeadID:      req.LeadID,
		ContactID:   req.ContactID,
		Status:      "INITIATED",
		Direction:   "OUTBOUND",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

// sendViaProvider sends message using the configured provider
func (s *MultiChannelCommunicationService) sendViaProvider(ctx context.Context, channel *models.CommunicationChannel, message *models.CommunicationMessage) (string, error) {
	switch strings.ToUpper(channel.APIProvider) {
	case "TWILIO":
		return s.sendViaTwilio(ctx, channel, message)
	case "SENDGRID":
		return s.sendViaSendGrid(ctx, channel, message)
	case "TELEGRAM":
		return s.sendViaTelegram(ctx, channel, message)
	case "WHATSAPP_BUSINESS":
		return s.sendViaWhatsAppBusiness(ctx, channel, message)
	case "AWS_SNS":
		return s.sendViaAWSSNS(ctx, channel, message)
	case "VONAGE":
		return s.sendViaVonage(ctx, channel, message)
	case "MAILGUN":
		return s.sendViaMailgun(ctx, channel, message)
	default:
		return "", fmt.Errorf("unsupported provider: %s", channel.APIProvider)
	}
}

// sendViaTwilio sends SMS via Twilio
func (s *MultiChannelCommunicationService) sendViaTwilio(ctx context.Context, channel *models.CommunicationChannel, message *models.CommunicationMessage) (string, error) {
	payload := map[string]string{
		"From": channel.SenderID,
		"To":   message.ToAddress,
		"Body": message.MessageBody,
	}

	return s.sendHTTPRequest(ctx, "POST", channel, payload)
}

// sendViaSendGrid sends email via SendGrid
func (s *MultiChannelCommunicationService) sendViaSendGrid(ctx context.Context, channel *models.CommunicationChannel, message *models.CommunicationMessage) (string, error) {
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{"email": message.ToAddress},
				},
			},
		},
		"from":    map[string]string{"email": channel.SenderID},
		"subject": message.MessageSubject,
		"content": []map[string]string{{"type": "text/html", "value": message.MessageHTML}},
	}

	return s.sendHTTPRequest(ctx, "POST", channel, payload)
}

// sendViaTelegram sends message via Telegram Bot API
func (s *MultiChannelCommunicationService) sendViaTelegram(ctx context.Context, channel *models.CommunicationChannel, message *models.CommunicationMessage) (string, error) {
	// Extract chat ID from ToAddress (should be Telegram user ID)
	payload := map[string]interface{}{
		"chat_id": message.ToAddress,
		"text":    message.MessageBody,
	}

	return s.sendHTTPRequest(ctx, "POST", channel, payload)
}

// sendViaWhatsAppBusiness sends message via WhatsApp Business API
func (s *MultiChannelCommunicationService) sendViaWhatsAppBusiness(ctx context.Context, channel *models.CommunicationChannel, message *models.CommunicationMessage) (string, error) {
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                message.ToAddress,
		"type":              "text",
		"text": map[string]string{
			"preview_url": "false",
			"body":        message.MessageBody,
		},
	}

	return s.sendHTTPRequest(ctx, "POST", channel, payload)
}

// sendViaAWSSNS sends SMS via AWS SNS
func (s *MultiChannelCommunicationService) sendViaAWSSNS(ctx context.Context, channel *models.CommunicationChannel, message *models.CommunicationMessage) (string, error) {
	payload := map[string]interface{}{
		"Message":     message.MessageBody,
		"PhoneNumber": message.ToAddress,
	}

	return s.sendHTTPRequest(ctx, "POST", channel, payload)
}

// sendViaVonage sends SMS via Vonage (Nexmo)
func (s *MultiChannelCommunicationService) sendViaVonage(ctx context.Context, channel *models.CommunicationChannel, message *models.CommunicationMessage) (string, error) {
	payload := map[string]string{
		"from": channel.SenderID,
		"to":   message.ToAddress,
		"text": message.MessageBody,
	}

	return s.sendHTTPRequest(ctx, "POST", channel, payload)
}

// sendViaMailgun sends email via Mailgun
func (s *MultiChannelCommunicationService) sendViaMailgun(ctx context.Context, channel *models.CommunicationChannel, message *models.CommunicationMessage) (string, error) {
	payload := map[string]interface{}{
		"from":    channel.SenderID,
		"to":      message.ToAddress,
		"subject": message.MessageSubject,
		"html":    message.MessageHTML,
	}

	return s.sendHTTPRequest(ctx, "POST", channel, payload)
}

// sendHTTPRequest sends HTTP request to provider API
func (s *MultiChannelCommunicationService) sendHTTPRequest(ctx context.Context, method string, channel *models.CommunicationChannel, payload interface{}) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, method, channel.APIURL, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	// Add authentication headers
	if channel.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+channel.AuthToken)
	}
	if channel.APIKey != "" {
		req.Header.Set("X-API-Key", channel.APIKey)
	}
	if channel.APISecret != "" {
		req.Header.Set("X-API-Secret", channel.APISecret)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("provider returned error: %d - %s", resp.StatusCode, string(respBody))
	}

	// Extract message ID from response
	var respData map[string]interface{}
	json.Unmarshal(respBody, &respData)

	if msgID, ok := respData["id"].(string); ok {
		return msgID, nil
	}
	if msgID, ok := respData["message_id"].(string); ok {
		return msgID, nil
	}
	if msgID, ok := respData["sid"].(string); ok {
		return msgID, nil
	}

	return fmt.Sprintf("msg_%d", time.Now().Unix()), nil
}

// CreateMessageTemplate creates a new message template
func (s *MultiChannelCommunicationService) CreateMessageTemplate(ctx context.Context, tenantID string, req *models.MessageTemplateRequest) (*models.MessageTemplate, error) {
	template := &models.MessageTemplate{
		ID:           uuid.New().String(),
		TenantID:     tenantID,
		TemplateName: req.TemplateName,
		ChannelType:  req.ChannelType,
		TemplateBody: req.TemplateBody,
		Subject:      req.Subject,
		Language:     req.Language,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if req.TemplateVariables != nil {
		varBytes, _ := json.Marshal(req.TemplateVariables)
		template.TemplateVariables = varBytes
	}

	if err := s.db.WithContext(ctx).Create(template).Error; err != nil {
		return nil, err
	}

	return template, nil
}

// GetMessageTemplate retrieves a template
func (s *MultiChannelCommunicationService) GetMessageTemplate(ctx context.Context, tenantID, templateID string) (*models.MessageTemplate, error) {
	var template models.MessageTemplate
	if err := s.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", templateID, tenantID).
		First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// ListMessageTemplates lists all templates
func (s *MultiChannelCommunicationService) ListMessageTemplates(ctx context.Context, tenantID string, limit, offset int) ([]*models.MessageTemplate, error) {
	var templates []*models.MessageTemplate
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = 1", tenantID).
		Limit(limit).
		Offset(offset).
		Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// UpdateContactPreference updates contact communication preferences
func (s *MultiChannelCommunicationService) UpdateContactPreference(ctx context.Context, tenantID string, contactID string, req *models.UpdateContactPreferenceRequest) (*models.ContactCommunicationPreference, error) {
	pref := &models.ContactCommunicationPreference{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		ContactID: contactID,
	}

	if req.EmailAddress != "" {
		pref.EmailAddress = req.EmailAddress
	}
	if req.PhoneNumber != "" {
		pref.PhoneNumber = req.PhoneNumber
	}
	if req.TelegramID != "" {
		pref.TelegramID = req.TelegramID
	}
	if req.WhatsappNumber != "" {
		pref.WhatsappNumber = req.WhatsappNumber
	}
	if req.PreferredChannel != "" {
		pref.PreferredChannel = req.PreferredChannel
	}
	if req.AllowTelegram != nil {
		pref.AllowTelegram = *req.AllowTelegram
	}
	if req.AllowWhatsapp != nil {
		pref.AllowWhatsapp = *req.AllowWhatsapp
	}
	if req.AllowSMS != nil {
		pref.AllowSMS = *req.AllowSMS
	}
	if req.AllowEmail != nil {
		pref.AllowEmail = *req.AllowEmail
	}

	// Use upsert pattern
	if err := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tenant_id"}, {Name: "contact_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"email_address", "phone_number", "telegram_id", "whatsapp_number", "preferred_channel", "allow_telegram", "allow_whatsapp", "allow_sms", "allow_email"}),
		}).
		Create(pref).Error; err != nil {
		return nil, err
	}

	return pref, nil
}

// GetContactPreference gets contact preferences
func (s *MultiChannelCommunicationService) GetContactPreference(ctx context.Context, tenantID, contactID string) (*models.ContactCommunicationPreference, error) {
	var pref models.ContactCommunicationPreference
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND contact_id = ?", tenantID, contactID).
		First(&pref).Error; err != nil {
		return nil, err
	}
	return &pref, nil
}

// ProcessWebhookEvent processes incoming webhook from providers
func (s *MultiChannelCommunicationService) ProcessWebhookEvent(ctx context.Context, tenantID, channelType string, signature string, payload []byte) error {
	// Verify webhook signature
	var channel models.CommunicationChannel
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND channel_type = ?", tenantID, channelType).
		First(&channel).Error; err != nil {
		return err
	}

	if !s.verifyWebhookSignature(&channel, signature, payload) {
		return fmt.Errorf("invalid webhook signature")
	}

	// Log webhook
	webhookLog := &models.CommunicationWebhookLog{
		ID:               uuid.New().String(),
		TenantID:         tenantID,
		ChannelID:        channel.ID,
		ChannelType:      channelType,
		WebhookPayload:   payload,
		WebhookSignature: signature,
		IsValid:          true,
		ProcessingStatus: "PENDING",
		ReceivedAt:       time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(webhookLog).Error; err != nil {
		return err
	}

	// Parse and process webhook
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}

	// Update message status based on webhook event
	if messageID, ok := data["message_id"].(string); ok {
		if eventType, ok := data["event"].(string); ok {
			s.updateMessageStatus(ctx, messageID, eventType)
		}
	}

	return nil
}

// verifyWebhookSignature verifies webhook signature
func (s *MultiChannelCommunicationService) verifyWebhookSignature(channel *models.CommunicationChannel, signature string, payload []byte) bool {
	if channel.APISecret == "" {
		return true
	}

	h := hmac.New(sha256.New, []byte(channel.APISecret))
	h.Write(payload)
	expectedSig := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSig))
}

// updateMessageStatus updates message status based on webhook event
func (s *MultiChannelCommunicationService) updateMessageStatus(ctx context.Context, messageID string, eventType string) error {
	statusMap := map[string]string{
		"delivered": "DELIVERED",
		"read":      "READ",
		"failed":    "FAILED",
		"bounced":   "BOUNCED",
	}

	status, ok := statusMap[strings.ToLower(eventType)]
	if !ok {
		return nil
	}

	return s.db.WithContext(ctx).
		Model(&models.CommunicationMessage{}).
		Where("external_message_id = ?", messageID).
		Update("status", status).Error
}

// GetCommunicationSession retrieves a session
func (s *MultiChannelCommunicationService) GetCommunicationSession(ctx context.Context, tenantID, sessionID string) (*models.CommunicationSession, error) {
	var session models.CommunicationSession
	if err := s.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", sessionID, tenantID).
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// ListCommunicationSessions lists sessions for a tenant
func (s *MultiChannelCommunicationService) ListCommunicationSessions(ctx context.Context, tenantID string, limit, offset int) ([]*models.CommunicationSession, error) {
	var sessions []*models.CommunicationSession
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// GetCommunicationAnalytics retrieves analytics for a date range
func (s *MultiChannelCommunicationService) GetCommunicationAnalytics(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*models.CommunicationAnalytics, error) {
	var analytics []*models.CommunicationAnalytics
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND metric_date BETWEEN ? AND ?", tenantID, startDate, endDate).
		Find(&analytics).Error; err != nil {
		return nil, err
	}
	return analytics, nil
}

// Logger interface for logging
type Logger interface {
	Error(string)
	Info(string)
	Debug(string)
}
