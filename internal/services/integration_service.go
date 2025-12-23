package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"vyomtech-backend/internal/models"
)

// IntegrationService handles third-party integration operations
type IntegrationService struct {
	db *sql.DB
}

// NewIntegrationService creates a new integration service
func NewIntegrationService(db *sql.DB) *IntegrationService {
	return &IntegrationService{db: db}
}

// CreateProvider creates a new integration provider
func (s *IntegrationService) CreateProvider(ctx context.Context, tenantID string, req *models.CreateProviderRequest) (*models.IntegrationProvider, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	query := `
		INSERT INTO integration_providers (
			tenant_id, name, type, api_base_url, is_active, rate_limit, retry_count,
			timeout_seconds, webhook_secret, oauth_client_id, oauth_client_secret, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var metadataStr *string
	if req.Metadata != nil {
		metaStr := string(*req.Metadata)
		metadataStr = &metaStr
	}

	result, err := s.db.ExecContext(ctx, query,
		tenantID, req.Name, req.Type, req.APIBaseURL, true,
		req.RateLimit, req.RetryCount, req.TimeoutSeconds,
		req.WebhookSecret, req.OAuthClientID, req.OAuthClientSecret, metadataStr,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert ID: %w", err)
	}

	return s.GetProvider(ctx, tenantID, id)
}

// GetProvider retrieves a provider by ID
func (s *IntegrationService) GetProvider(ctx context.Context, tenantID string, providerID int64) (*models.IntegrationProvider, error) {
	query := `
		SELECT id, tenant_id, name, type, api_base_url, is_active, rate_limit, retry_count,
			   timeout_seconds, webhook_secret, oauth_client_id, oauth_client_secret, metadata,
			   created_at, updated_at, deleted_at
		FROM integration_providers
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	var provider models.IntegrationProvider
	var metadata *string
	err := s.db.QueryRowContext(ctx, query, providerID, tenantID).Scan(
		&provider.ID, &provider.TenantID, &provider.Name, &provider.Type, &provider.APIBaseURL,
		&provider.IsActive, &provider.RateLimit, &provider.RetryCount, &provider.TimeoutSeconds,
		&provider.WebhookSecret, &provider.OAuthClientID, &provider.OAuthClientSecret, &metadata,
		&provider.CreatedAt, &provider.UpdatedAt, &provider.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("provider not found")
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	if metadata != nil {
		raw := json.RawMessage(*metadata)
		provider.Metadata = &raw
	}

	return &provider, nil
}

// ListProviders lists providers with pagination
func (s *IntegrationService) ListProviders(ctx context.Context, tenantID string, limit, offset int) ([]models.IntegrationProvider, error) {
	query := `
		SELECT id, tenant_id, name, type, api_base_url, is_active, rate_limit, retry_count,
			   timeout_seconds, webhook_secret, oauth_client_id, oauth_client_secret, metadata,
			   created_at, updated_at, deleted_at
		FROM integration_providers
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}
	defer rows.Close()

	var providers []models.IntegrationProvider
	for rows.Next() {
		var provider models.IntegrationProvider
		var metadata *string
		err := rows.Scan(
			&provider.ID, &provider.TenantID, &provider.Name, &provider.Type, &provider.APIBaseURL,
			&provider.IsActive, &provider.RateLimit, &provider.RetryCount, &provider.TimeoutSeconds,
			&provider.WebhookSecret, &provider.OAuthClientID, &provider.OAuthClientSecret, &metadata,
			&provider.CreatedAt, &provider.UpdatedAt, &provider.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan provider: %w", err)
		}

		if metadata != nil {
			raw := json.RawMessage(*metadata)
			provider.Metadata = &raw
		}

		providers = append(providers, provider)
	}

	return providers, nil
}

// UpdateProvider updates a provider
func (s *IntegrationService) UpdateProvider(ctx context.Context, tenantID string, providerID int64, req *models.UpdateProviderRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	query := `UPDATE integration_providers SET `
	args := []interface{}{}
	updates := 0

	if req.Name != nil {
		if updates > 0 {
			query += ", "
		}
		query += "name = ?"
		args = append(args, *req.Name)
		updates++
	}
	if req.APIBaseURL != nil {
		if updates > 0 {
			query += ", "
		}
		query += "api_base_url = ?"
		args = append(args, *req.APIBaseURL)
		updates++
	}
	if req.IsActive != nil {
		if updates > 0 {
			query += ", "
		}
		query += "is_active = ?"
		args = append(args, *req.IsActive)
		updates++
	}
	if req.RateLimit != nil {
		if updates > 0 {
			query += ", "
		}
		query += "rate_limit = ?"
		args = append(args, *req.RateLimit)
		updates++
	}
	if req.RetryCount != nil {
		if updates > 0 {
			query += ", "
		}
		query += "retry_count = ?"
		args = append(args, *req.RetryCount)
		updates++
	}
	if req.TimeoutSeconds != nil {
		if updates > 0 {
			query += ", "
		}
		query += "timeout_seconds = ?"
		args = append(args, *req.TimeoutSeconds)
		updates++
	}
	if req.WebhookSecret != nil {
		if updates > 0 {
			query += ", "
		}
		query += "webhook_secret = ?"
		args = append(args, *req.WebhookSecret)
		updates++
	}

	if updates == 0 {
		return errors.New("no fields to update")
	}

	query += ", updated_at = NOW() WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL"
	args = append(args, providerID, tenantID)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update provider: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("provider not found")
	}

	return nil
}

// DeleteProvider soft-deletes a provider
func (s *IntegrationService) DeleteProvider(ctx context.Context, tenantID string, providerID int64) error {
	query := `
		UPDATE integration_providers
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, providerID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("provider not found")
	}

	return nil
}

// CreateWebhook creates a new webhook configuration
func (s *IntegrationService) CreateWebhook(ctx context.Context, tenantID string, req *models.CreateWebhookRequest) (*models.IntegrationWebhook, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	query := `
		INSERT INTO integration_webhooks (
			tenant_id, provider_id, event_type, webhook_url, is_active, retry_policy,
			max_retries, timeout_seconds, headers, filter_conditions
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var headersStr *string
	if req.Headers != nil {
		h := string(*req.Headers)
		headersStr = &h
	}

	var filterStr *string
	if req.FilterConditions != nil {
		f := string(*req.FilterConditions)
		filterStr = &f
	}

	result, err := s.db.ExecContext(ctx, query,
		tenantID, req.ProviderID, req.EventType, req.WebhookURL, true,
		req.RetryPolicy, req.MaxRetries, req.TimeoutSeconds, headersStr, filterStr,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert ID: %w", err)
	}

	return s.GetWebhook(ctx, tenantID, id)
}

// GetWebhook retrieves a webhook by ID
func (s *IntegrationService) GetWebhook(ctx context.Context, tenantID string, webhookID int64) (*models.IntegrationWebhook, error) {
	query := `
		SELECT id, tenant_id, provider_id, event_type, webhook_url, is_active, retry_policy,
			   max_retries, timeout_seconds, headers, filter_conditions, created_at, updated_at, deleted_at
		FROM integration_webhooks
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	var webhook models.IntegrationWebhook
	var headers, filters *string
	err := s.db.QueryRowContext(ctx, query, webhookID, tenantID).Scan(
		&webhook.ID, &webhook.TenantID, &webhook.ProviderID, &webhook.EventType, &webhook.WebhookURL,
		&webhook.IsActive, &webhook.RetryPolicy, &webhook.MaxRetries, &webhook.TimeoutSeconds,
		&headers, &filters, &webhook.CreatedAt, &webhook.UpdatedAt, &webhook.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("webhook not found")
		}
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}

	if headers != nil {
		raw := json.RawMessage(*headers)
		webhook.Headers = &raw
	}
	if filters != nil {
		raw := json.RawMessage(*filters)
		webhook.FilterConditions = &raw
	}

	return &webhook, nil
}

// ListWebhooks lists webhooks for a provider
func (s *IntegrationService) ListWebhooks(ctx context.Context, tenantID string, providerID int64) ([]models.IntegrationWebhook, error) {
	query := `
		SELECT id, tenant_id, provider_id, event_type, webhook_url, is_active, retry_policy,
			   max_retries, timeout_seconds, headers, filter_conditions, created_at, updated_at, deleted_at
		FROM integration_webhooks
		WHERE tenant_id = ? AND provider_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}
	defer rows.Close()

	var webhooks []models.IntegrationWebhook
	for rows.Next() {
		var webhook models.IntegrationWebhook
		var headers, filters *string
		err := rows.Scan(
			&webhook.ID, &webhook.TenantID, &webhook.ProviderID, &webhook.EventType, &webhook.WebhookURL,
			&webhook.IsActive, &webhook.RetryPolicy, &webhook.MaxRetries, &webhook.TimeoutSeconds,
			&headers, &filters, &webhook.CreatedAt, &webhook.UpdatedAt, &webhook.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan webhook: %w", err)
		}

		if headers != nil {
			raw := json.RawMessage(*headers)
			webhook.Headers = &raw
		}
		if filters != nil {
			raw := json.RawMessage(*filters)
			webhook.FilterConditions = &raw
		}

		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

// CreateSyncJob creates a new sync job
func (s *IntegrationService) CreateSyncJob(ctx context.Context, tenantID string, req *models.TriggerSyncRequest) (*models.IntegrationSyncJob, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	query := `
		INSERT INTO integration_sync_jobs (
			tenant_id, provider_id, sync_type, status, sync_config
		) VALUES (?, ?, ?, 'SCHEDULED', ?)
	`

	var configStr *string
	if req.SyncConfig != nil {
		c := string(*req.SyncConfig)
		configStr = &c
	}

	result, err := s.db.ExecContext(ctx, query,
		tenantID, req.ProviderID, req.SyncType, configStr,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create sync job: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get insert ID: %w", err)
	}

	return s.GetSyncJob(ctx, tenantID, id)
}

// GetSyncJob retrieves a sync job by ID
func (s *IntegrationService) GetSyncJob(ctx context.Context, tenantID string, jobID int64) (*models.IntegrationSyncJob, error) {
	query := `
		SELECT id, tenant_id, provider_id, sync_type, status, last_sync_at, next_sync_at,
			   records_synced, records_failed, sync_duration_seconds, error_log, sync_config,
			   created_at, updated_at
		FROM integration_sync_jobs
		WHERE id = ? AND tenant_id = ?
	`

	var job models.IntegrationSyncJob
	var errorLog, config *string
	err := s.db.QueryRowContext(ctx, query, jobID, tenantID).Scan(
		&job.ID, &job.TenantID, &job.ProviderID, &job.SyncType, &job.Status,
		&job.LastSyncAt, &job.NextSyncAt, &job.RecordsSynced, &job.RecordsFailed,
		&job.SyncDurationSecs, &errorLog, &config, &job.CreatedAt, &job.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("sync job not found")
		}
		return nil, fmt.Errorf("failed to get sync job: %w", err)
	}

	if errorLog != nil {
		job.ErrorLog = errorLog
	}
	if config != nil {
		raw := json.RawMessage(*config)
		job.SyncConfig = &raw
	}

	return &job, nil
}

// ListSyncJobs lists sync jobs for a provider
func (s *IntegrationService) ListSyncJobs(ctx context.Context, tenantID string, providerID int64) ([]models.IntegrationSyncJob, error) {
	query := `
		SELECT id, tenant_id, provider_id, sync_type, status, last_sync_at, next_sync_at,
			   records_synced, records_failed, sync_duration_seconds, error_log, sync_config,
			   created_at, updated_at
		FROM integration_sync_jobs
		WHERE tenant_id = ? AND provider_id = ?
		ORDER BY created_at DESC
		LIMIT 50
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sync jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.IntegrationSyncJob
	for rows.Next() {
		var job models.IntegrationSyncJob
		var errorLog, config *string
		err := rows.Scan(
			&job.ID, &job.TenantID, &job.ProviderID, &job.SyncType, &job.Status,
			&job.LastSyncAt, &job.NextSyncAt, &job.RecordsSynced, &job.RecordsFailed,
			&job.SyncDurationSecs, &errorLog, &config, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sync job: %w", err)
		}

		if errorLog != nil {
			job.ErrorLog = errorLog
		}
		if config != nil {
			raw := json.RawMessage(*config)
			job.SyncConfig = &raw
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// CreateErrorLog logs an integration error
func (s *IntegrationService) CreateErrorLog(ctx context.Context, tenantID string, providerID int64, errorCode, message string, severity string) error {
	query := `
		INSERT INTO integration_error_logs (
			tenant_id, provider_id, error_code, error_message, severity
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query, tenantID, providerID, errorCode, message, severity)
	return err
}

// ListErrorLogs lists recent error logs
func (s *IntegrationService) ListErrorLogs(ctx context.Context, tenantID string, providerID int64, limit int) ([]models.IntegrationErrorLog, error) {
	query := `
		SELECT id, tenant_id, provider_id, error_code, error_message, error_details,
			   endpoint, request_payload, response_payload, severity, resolved, resolved_at, created_at
		FROM integration_error_logs
		WHERE tenant_id = ? AND provider_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, providerID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list error logs: %w", err)
	}
	defer rows.Close()

	var logs []models.IntegrationErrorLog
	for rows.Next() {
		var log models.IntegrationErrorLog
		var details, requestPayload, responsePayload *string
		err := rows.Scan(
			&log.ID, &log.TenantID, &log.ProviderID, &log.ErrorCode, &log.ErrorMessage,
			&details, &log.Endpoint, &requestPayload, &responsePayload, &log.Severity,
			&log.Resolved, &log.ResolvedAt, &log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan error log: %w", err)
		}

		if details != nil {
			raw := json.RawMessage(*details)
			log.ErrorDetails = &raw
		}
		if requestPayload != nil {
			raw := json.RawMessage(*requestPayload)
			log.RequestPayload = &raw
		}
		if responsePayload != nil {
			raw := json.RawMessage(*responsePayload)
			log.ResponsePayload = &raw
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// GetIntegrationStats returns integration statistics
func (s *IntegrationService) GetIntegrationStats(ctx context.Context, tenantID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count providers
	var providerCount int64
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM integration_providers
		WHERE tenant_id = ? AND deleted_at IS NULL
	`, tenantID).Scan(&providerCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count providers: %w", err)
	}
	stats["total_providers"] = providerCount

	// Count webhooks
	var webhookCount int64
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM integration_webhooks
		WHERE tenant_id = ? AND deleted_at IS NULL
	`, tenantID).Scan(&webhookCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count webhooks: %w", err)
	}
	stats["total_webhooks"] = webhookCount

	// Count recent sync jobs
	var recentSyncs int64
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM integration_sync_jobs
		WHERE tenant_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
	`, tenantID).Scan(&recentSyncs)
	if err != nil {
		return nil, fmt.Errorf("failed to count recent syncs: %w", err)
	}
	stats["recent_syncs_7d"] = recentSyncs

	// Count errors
	var errorCount int64
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM integration_error_logs
		WHERE tenant_id = ? AND resolved = FALSE
	`, tenantID).Scan(&errorCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count errors: %w", err)
	}
	stats["unresolved_errors"] = errorCount

	return stats, nil
}
