package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/logger"
)

// ClickToCallService handles click-to-call operations
type ClickToCallService struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewClickToCallService creates a new ClickToCallService
func NewClickToCallService(db *sql.DB, logger *logger.Logger) *ClickToCallService {
	return &ClickToCallService{
		db:     db,
		logger: logger,
	}
}

// CreateClickToCallSession initiates a new click-to-call session
func (s *ClickToCallService) CreateClickToCallSession(ctx context.Context, tenantID string, req *models.CreateClickToCallRequest) (*models.ClickToCallSession, error) {
	session := &models.ClickToCallSession{
		ID:            generateID(),
		TenantID:      tenantID,
		ToPhone:       req.ToPhone,
		PhoneType:     req.PhoneType,
		LeadID:        req.LeadID,
		AgentID:       req.AgentID,
		AccountID:     req.AccountID,
		CampaignID:    req.CampaignID,
		ContactName:   req.ContactName,
		ContactEmail:  req.ContactEmail,
		Status:        "INITIATED",
		Direction:     req.Direction,
		SessionID:     generateSessionID(),
		CorrelationID: generateCorrelationID(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Determine provider based on routing rules
	provider, err := s.SelectProvider(ctx, tenantID, req)
	if err != nil {
		s.logger.Error("Failed to select provider", "error", err, "tenant", tenantID)
		session.Status = "FAILED"
		session.ErrorMessage = "No provider available"
		session.ErrorCode = "PROVIDER_NOT_FOUND"
	} else {
		session.ProviderID = provider.ID
		session.ProviderType = provider.ProviderType
		session.FromPhone = provider.PhoneNumber
	}

	// Store metadata
	if req.CallMetadata != nil {
		metadataJSON, _ := json.Marshal(req.CallMetadata)
		session.Metadata = string(metadataJSON)
	}

	// Insert into database
	query := `
		INSERT INTO click_to_call_session (
			id, tenant_id, initiated_by, from_phone, to_phone, phone_type, contact_name, contact_email,
			contact_id, lead_id, agent_id, account_id, campaign_id, provider_id, provider_type,
			session_id, correlation_id, status, direction, error_code, error_message, notes, metadata, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err = s.db.ExecContext(ctx, query,
		session.ID, session.TenantID, session.InitiatedBy, session.FromPhone, session.ToPhone,
		session.PhoneType, session.ContactName, session.ContactEmail, session.ContactID,
		session.LeadID, session.AgentID, session.AccountID, session.CampaignID, session.ProviderID,
		session.ProviderType, session.SessionID, session.CorrelationID, session.Status, session.Direction,
		session.ErrorCode, session.ErrorMessage, session.Notes, session.Metadata,
	)
	if err != nil {
		s.logger.Error("Failed to create session", "error", err, "tenant", tenantID)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// UpdateSessionStatus updates the status of a click-to-call session
func (s *ClickToCallService) UpdateSessionStatus(ctx context.Context, sessionID, tenantID, status string) error {
	query := `
		UPDATE click_to_call_session
		SET status = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	result, err := s.db.ExecContext(ctx, query, status, sessionID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// UpdateCallTiming updates call timing information
func (s *ClickToCallService) UpdateCallTiming(ctx context.Context, sessionID, tenantID string, timing map[string]interface{}) error {
	query := `
		UPDATE click_to_call_session
		SET call_started_at = ?, call_ended_at = ?, duration_seconds = ?, ring_time_seconds = ?, answer_time_seconds = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	result, err := s.db.ExecContext(ctx, query,
		timing["call_started_at"],
		timing["call_ended_at"],
		timing["duration_seconds"],
		timing["ring_time_seconds"],
		timing["answer_time_seconds"],
		sessionID,
		tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update call timing: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// EndSession ends a click-to-call session
func (s *ClickToCallService) EndSession(ctx context.Context, sessionID, tenantID string, reason, outcome string) error {
	now := time.Now()
	query := `
		UPDATE click_to_call_session
		SET status = 'COMPLETED', call_ended_at = ?, disconnect_reason = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	result, err := s.db.ExecContext(ctx, query, now, reason, sessionID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to end session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// GetSession retrieves a click-to-call session
func (s *ClickToCallService) GetSession(ctx context.Context, sessionID, tenantID string) (*models.ClickToCallSession, error) {
	query := `
		SELECT id, tenant_id, initiated_by, from_phone, to_phone, phone_type, contact_name, contact_email,
		       contact_id, lead_id, agent_id, account_id, campaign_id, provider_id, provider_type,
		       session_id, correlation_id, status, direction, call_started_at, call_ended_at, duration_seconds,
		       ring_time_seconds, answer_time_seconds, disconnect_reason, error_code, error_message,
		       recording_url, transcript_url, call_quality_score, is_recorded, is_transferred, transfer_to_agent,
		       notes, metadata, created_at, updated_at, deleted_at
		FROM click_to_call_session
		WHERE id = ? AND tenant_id = ?
	`

	session := &models.ClickToCallSession{}
	err := s.db.QueryRowContext(ctx, query, sessionID, tenantID).Scan(
		&session.ID, &session.TenantID, &session.InitiatedBy, &session.FromPhone, &session.ToPhone,
		&session.PhoneType, &session.ContactName, &session.ContactEmail, &session.ContactID,
		&session.LeadID, &session.AgentID, &session.AccountID, &session.CampaignID, &session.ProviderID,
		&session.ProviderType, &session.SessionID, &session.CorrelationID, &session.Status, &session.Direction,
		&session.CallStartedAt, &session.CallEndedAt, &session.DurationSeconds, &session.RingTimeSeconds,
		&session.AnswerTimeSeconds, &session.DisconnectReason, &session.ErrorCode, &session.ErrorMessage,
		&session.RecordingURL, &session.TranscriptURL, &session.CallQualityScore, &session.IsRecorded,
		&session.IsTransferred, &session.TransferToAgent, &session.Notes, &session.Metadata,
		&session.CreatedAt, &session.UpdatedAt, &session.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// ListSessions lists all sessions with filters
func (s *ClickToCallService) ListSessions(ctx context.Context, tenantID string, filters map[string]interface{}, limit, offset int) ([]models.ClickToCallSession, int64, error) {
	query := `
		SELECT id, tenant_id, initiated_by, from_phone, to_phone, phone_type, contact_name, contact_email,
		       contact_id, lead_id, agent_id, account_id, campaign_id, provider_id, provider_type,
		       session_id, correlation_id, status, direction, call_started_at, call_ended_at, duration_seconds,
		       ring_time_seconds, answer_time_seconds, disconnect_reason, error_code, error_message,
		       recording_url, transcript_url, call_quality_score, is_recorded, is_transferred, transfer_to_agent,
		       notes, metadata, created_at, updated_at, deleted_at
		FROM click_to_call_session
		WHERE tenant_id = ? AND deleted_at IS NULL
	`

	args := []interface{}{tenantID}

	// Apply filters
	if status, ok := filters["status"]; ok {
		query += " AND status = ?"
		args = append(args, status)
	}
	if agentID, ok := filters["agent_id"]; ok {
		query += " AND agent_id = ?"
		args = append(args, agentID)
	}
	if leadID, ok := filters["lead_id"]; ok {
		query += " AND lead_id = ?"
		args = append(args, leadID)
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", query)
	var total int64
	err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	// Add ordering and pagination
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.ClickToCallSession
	for rows.Next() {
		session := models.ClickToCallSession{}
		err := rows.Scan(
			&session.ID, &session.TenantID, &session.InitiatedBy, &session.FromPhone, &session.ToPhone,
			&session.PhoneType, &session.ContactName, &session.ContactEmail, &session.ContactID,
			&session.LeadID, &session.AgentID, &session.AccountID, &session.CampaignID, &session.ProviderID,
			&session.ProviderType, &session.SessionID, &session.CorrelationID, &session.Status, &session.Direction,
			&session.CallStartedAt, &session.CallEndedAt, &session.DurationSeconds, &session.RingTimeSeconds,
			&session.AnswerTimeSeconds, &session.DisconnectReason, &session.ErrorCode, &session.ErrorMessage,
			&session.RecordingURL, &session.TranscriptURL, &session.CallQualityScore, &session.IsRecorded,
			&session.IsTransferred, &session.TransferToAgent, &session.Notes, &session.Metadata,
			&session.CreatedAt, &session.UpdatedAt, &session.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, total, nil
}

// CreateVoIPProvider creates a new VoIP provider configuration
func (s *ClickToCallService) CreateVoIPProvider(ctx context.Context, tenantID string, provider *models.VoIPProvider) error {
	provider.ID = generateID()
	provider.TenantID = tenantID
	provider.CreatedAt = time.Now()
	provider.UpdatedAt = time.Now()

	query := `
		INSERT INTO voip_provider (
			id, tenant_id, provider_name, provider_type, api_key, api_secret, api_url, webhook_url,
			callback_url, auth_token, phone_number, caller_id, dial_plan_prefix, is_active,
			retry_count, timeout_seconds, priority, config_json, notes, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := s.db.ExecContext(ctx, query,
		provider.ID, provider.TenantID, provider.ProviderName, provider.ProviderType,
		provider.APIKey, provider.APISecret, provider.APIURL, provider.WebhookURL,
		provider.CallbackURL, provider.AuthToken, provider.PhoneNumber, provider.CallerID,
		provider.DialPlanPrefix, provider.IsActive, provider.RetryCount, provider.TimeoutSeconds,
		provider.Priority, provider.ConfigJSON, provider.Notes, provider.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	return nil
}

// GetVoIPProvider gets a VoIP provider by ID
func (s *ClickToCallService) GetVoIPProvider(ctx context.Context, providerID, tenantID string) (*models.VoIPProvider, error) {
	query := `
		SELECT id, tenant_id, provider_name, provider_type, api_key, api_secret, api_url, webhook_url,
		       callback_url, auth_token, phone_number, caller_id, dial_plan_prefix, is_active,
		       retry_count, timeout_seconds, priority, config_json, notes, created_by, created_at, updated_at
		FROM voip_provider
		WHERE id = ? AND tenant_id = ?
	`

	provider := &models.VoIPProvider{}
	err := s.db.QueryRowContext(ctx, query, providerID, tenantID).Scan(
		&provider.ID, &provider.TenantID, &provider.ProviderName, &provider.ProviderType,
		&provider.APIKey, &provider.APISecret, &provider.APIURL, &provider.WebhookURL,
		&provider.CallbackURL, &provider.AuthToken, &provider.PhoneNumber, &provider.CallerID,
		&provider.DialPlanPrefix, &provider.IsActive, &provider.RetryCount, &provider.TimeoutSeconds,
		&provider.Priority, &provider.ConfigJSON, &provider.Notes, &provider.CreatedBy,
		&provider.CreatedAt, &provider.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("provider not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return provider, nil
}

// ListVoIPProviders lists all VoIP providers for a tenant
func (s *ClickToCallService) ListVoIPProviders(ctx context.Context, tenantID string) ([]models.VoIPProvider, error) {
	query := `
		SELECT id, tenant_id, provider_name, provider_type, api_key, api_secret, api_url, webhook_url,
		       callback_url, auth_token, phone_number, caller_id, dial_plan_prefix, is_active,
		       retry_count, timeout_seconds, priority, config_json, notes, created_by, created_at, updated_at
		FROM voip_provider
		WHERE tenant_id = ? AND is_active = 1
		ORDER BY priority DESC, created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}
	defer rows.Close()

	var providers []models.VoIPProvider
	for rows.Next() {
		provider := models.VoIPProvider{}
		err := rows.Scan(
			&provider.ID, &provider.TenantID, &provider.ProviderName, &provider.ProviderType,
			&provider.APIKey, &provider.APISecret, &provider.APIURL, &provider.WebhookURL,
			&provider.CallbackURL, &provider.AuthToken, &provider.PhoneNumber, &provider.CallerID,
			&provider.DialPlanPrefix, &provider.IsActive, &provider.RetryCount, &provider.TimeoutSeconds,
			&provider.Priority, &provider.ConfigJSON, &provider.Notes, &provider.CreatedBy,
			&provider.CreatedAt, &provider.UpdatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan provider", "error", err)
			continue
		}
		providers = append(providers, provider)
	}

	return providers, nil
}

// SelectProvider selects the best provider based on routing rules
func (s *ClickToCallService) SelectProvider(ctx context.Context, tenantID string, req *models.CreateClickToCallRequest) (*models.VoIPProvider, error) {
	// If provider is specified, use it
	if req.ProviderID != "" {
		return s.GetVoIPProvider(ctx, req.ProviderID, tenantID)
	}

	// Otherwise, get the highest priority active provider
	providers, err := s.ListVoIPProviders(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("no active providers configured")
	}

	return &providers[0], nil
}

// SaveWebhookLog saves a webhook event log
func (s *ClickToCallService) SaveWebhookLog(ctx context.Context, tenantID string, log *models.CallWebhookLog) error {
	log.ID = generateID()
	log.TenantID = tenantID
	log.ReceivedAt = time.Now()
	log.ProcessingStatus = "PENDING"

	query := `
		INSERT INTO click_to_call_webhook_log (
			id, tenant_id, provider_id, webhook_event_type, webhook_payload, webhook_signature,
			is_valid, processing_status, error_message, received_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`

	_, err := s.db.ExecContext(ctx, query,
		log.ID, log.TenantID, log.ProviderID, log.WebhookEventType, log.WebhookPayload,
		log.WebhookSignature, log.IsValid, log.ProcessingStatus, log.ErrorMessage,
	)
	if err != nil {
		return fmt.Errorf("failed to save webhook log: %w", err)
	}

	return nil
}

// ProcessWebhookEvent processes an incoming webhook event
func (s *ClickToCallService) ProcessWebhookEvent(ctx context.Context, tenantID string, payload *models.CallWebhookPayload) error {
	if payload.EventType == "" || payload.SessionID == "" {
		return fmt.Errorf("invalid webhook payload: missing event_type or session_id")
	}

	// Update session based on event type
	switch payload.EventType {
	case "CALL_INITIATED":
		return s.UpdateSessionStatus(ctx, payload.SessionID, tenantID, "CONNECTING")
	case "CALL_RINGING":
		return s.UpdateSessionStatus(ctx, payload.SessionID, tenantID, "RINGING")
	case "CALL_ANSWERED":
		return s.UpdateSessionStatus(ctx, payload.SessionID, tenantID, "CONNECTED")
	case "CALL_ENDED":
		return s.EndSession(ctx, payload.SessionID, tenantID, payload.ErrorCode, "COMPLETED")
	case "CALL_FAILED":
		return s.UpdateSessionStatus(ctx, payload.SessionID, tenantID, "FAILED")
	default:
		s.logger.Warn("Unknown webhook event type", "event_type", payload.EventType)
		return nil
	}
}

// LogAgentActivity logs agent activity
func (s *ClickToCallService) LogAgentActivity(ctx context.Context, activity *models.AgentActivityLog) error {
	activity.ID = generateID()
	activity.ActivityTimestamp = time.Now()

	query := `
		INSERT INTO agent_activity_log (
			id, tenant_id, agent_id, activity_type, status_value, session_id, is_available,
			activity_timestamp, duration_seconds, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		activity.ID, activity.TenantID, activity.AgentID, activity.ActivityType, activity.StatusValue,
		activity.SessionID, activity.IsAvailable, activity.DurationSeconds, activity.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}

	return nil
}

// Helper functions
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func generateSessionID() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}

func generateCorrelationID() string {
	return fmt.Sprintf("corr_%d", time.Now().UnixNano())
}
