-- Phase 3.4: Integration Hub (Third-party integrations, webhooks, API gateway)

-- Integration Providers Configuration
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS integration_providers (
    id CHAR(26) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    type ENUM('CRM', 'ERP', 'ACCOUNTING', 'COMMUNICATION', 'ANALYTICS', 'PAYMENT', 'LOGISTICS') NOT NULL,
    api_base_url VARCHAR(500) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    rate_limit INT DEFAULT 1000,
    retry_count INT DEFAULT 3,
    timeout_seconds INT DEFAULT 30,
    webhook_secret VARCHAR(255),
    oauth_client_id VARCHAR(255),
    oauth_client_secret VARCHAR(255),
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    KEY idx_tenant_id (tenant_id),
    KEY idx_type (type),
    KEY idx_active (is_active),
    UNIQUE KEY uk_tenant_provider (tenant_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- API Credentials Management
CREATE TABLE IF NOT EXISTS integration_credentials (
    id CHAR(26) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    provider_id VARCHAR(36) NOT NULL,
    user_id CHAR(26) NOT NULL,
    api_key VARCHAR(500),
    api_secret VARCHAR(500),
    access_token VARCHAR(1000),
    refresh_token VARCHAR(1000),
    token_expires_at TIMESTAMP NULL,
    is_valid BOOLEAN DEFAULT TRUE,
    last_verified_at TIMESTAMP NULL,
    verification_error VARCHAR(500),
    encryption_key_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    KEY idx_tenant_id (tenant_id),
    KEY idx_provider_id (provider_id),
    KEY idx_user_id (user_id),
    FOREIGN KEY (provider_id) REFERENCES integration_providers(id),
    UNIQUE KEY uk_tenant_provider_user (tenant_id, provider_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Webhook Configurations
CREATE TABLE IF NOT EXISTS integration_webhooks (
    id CHAR(26) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    provider_id VARCHAR(36) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    webhook_url VARCHAR(500) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    retry_policy ENUM('LINEAR', 'EXPONENTIAL', 'FIBONACCI') DEFAULT 'EXPONENTIAL',
    max_retries INT DEFAULT 5,
    timeout_seconds INT DEFAULT 30,
    headers JSON,
    filter_conditions JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    KEY idx_tenant_id (tenant_id),
    KEY idx_provider_id (provider_id),
    KEY idx_event_type (event_type),
    FOREIGN KEY (provider_id) REFERENCES integration_providers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Webhook Events Log
CREATE TABLE IF NOT EXISTS integration_webhook_events (
    id CHAR(26) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    webhook_id VARCHAR(36) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_data JSON NOT NULL,
    status ENUM('PENDING', 'DELIVERED', 'FAILED', 'RETRYING') DEFAULT 'PENDING',
    delivery_attempts INT DEFAULT 0,
    error_message VARCHAR(1000),
    delivered_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_webhook_id (webhook_id),
    KEY idx_status (status),
    FOREIGN KEY (webhook_id) REFERENCES integration_webhooks(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data Synchronization Mappings
CREATE TABLE IF NOT EXISTS integration_field_mappings (
    id CHAR(26) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    provider_id VARCHAR(36) NOT NULL,
    source_entity VARCHAR(100) NOT NULL,
    source_field VARCHAR(100) NOT NULL,
    target_entity VARCHAR(100) NOT NULL,
    target_field VARCHAR(100) NOT NULL,
    transformation_type ENUM('DIRECT', 'CUSTOM_FUNCTION', 'LOOKUP', 'AGGREGATE') DEFAULT 'DIRECT',
    transformation_config JSON,
    is_bidirectional BOOLEAN DEFAULT FALSE,
    sync_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    KEY idx_tenant_id (tenant_id),
    KEY idx_provider_id (provider_id),
    FOREIGN KEY (provider_id) REFERENCES integration_providers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Integration Sync Jobs
CREATE TABLE IF NOT EXISTS integration_sync_jobs (
    id CHAR(26) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    provider_id VARCHAR(36) NOT NULL,
    sync_type ENUM('FULL', 'INCREMENTAL', 'DELTA') DEFAULT 'INCREMENTAL',
    status ENUM('SCHEDULED', 'RUNNING', 'COMPLETED', 'FAILED') DEFAULT 'SCHEDULED',
    last_sync_at TIMESTAMP NULL,
    next_sync_at TIMESTAMP NULL,
    records_synced INT DEFAULT 0,
    records_failed INT DEFAULT 0,
    sync_duration_seconds INT,
    error_log TEXT,
    sync_config JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_provider_id (provider_id),
    KEY idx_status (status),
    FOREIGN KEY (provider_id) REFERENCES integration_providers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Integration API Rate Limiting
CREATE TABLE IF NOT EXISTS integration_rate_limits (
    id CHAR(26) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    provider_id VARCHAR(36) NOT NULL,
    user_id CHAR(26),
    requests_count INT DEFAULT 0,
    limit_window_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    limit_window_end TIMESTAMP NULL,
    is_rate_limited BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_provider_id (provider_id),
    FOREIGN KEY (provider_id) REFERENCES integration_providers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Integration Errors & Logs
CREATE TABLE IF NOT EXISTS integration_error_logs (
    id CHAR(26) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    provider_id VARCHAR(36) NOT NULL,
    error_code VARCHAR(50),
    error_message TEXT NOT NULL,
    error_details JSON,
    endpoint VARCHAR(500),
    request_payload JSON,
    response_payload JSON,
    severity ENUM('INFO', 'WARNING', 'ERROR', 'CRITICAL') DEFAULT 'ERROR',
    resolved BOOLEAN DEFAULT FALSE,
    resolved_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_provider_id (provider_id),
    KEY idx_severity (severity),
    FOREIGN KEY (provider_id) REFERENCES integration_providers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Integration Audit Trail
CREATE TABLE IF NOT EXISTS integration_audit_logs (
    id CHAR(26) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    provider_id VARCHAR(36) NOT NULL,
    user_id CHAR(26),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100),
    resource_id VARCHAR(36),
    old_values JSON,
    new_values JSON,
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY idx_tenant_id (tenant_id),
    KEY idx_provider_id (provider_id),
    KEY idx_user_id (user_id),
    KEY idx_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;




