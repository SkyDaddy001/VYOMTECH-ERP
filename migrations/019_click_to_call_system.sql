-- Migration: 019_click_to_call_system.sql
-- Description: Click-to-call system supporting Asterisk, SIP, mCube, Exotel
-- Date: 2025-12-03

-- Table for VoIP provider configuration
CREATE TABLE IF NOT EXISTS `voip_provider` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    provider_name VARCHAR(100) NOT NULL COMMENT 'asterisk, sip, mcube, exotel',
    provider_type ENUM('ASTERISK', 'SIP', 'MCUBE', 'EXOTEL', 'TWILIO', 'VONAGE', 'CUSTOM') NOT NULL,
    api_key VARCHAR(500),
    api_secret VARCHAR(500),
    api_url VARCHAR(255),
    webhook_url VARCHAR(500),
    callback_url VARCHAR(500),
    auth_token VARCHAR(500),
    phone_number VARCHAR(20),
    caller_id VARCHAR(20),
    dial_plan_prefix VARCHAR(10),
    is_active TINYINT(1) DEFAULT 1,
    retry_count INT DEFAULT 3,
    timeout_seconds INT DEFAULT 30,
    priority INT DEFAULT 0,
    config_json JSON,
    notes TEXT,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_provider (tenant_id, provider_type),
    INDEX idx_active (is_active),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='VoIP provider configuration for click-to-call';

-- Table for click-to-call session tracking
CREATE TABLE IF NOT EXISTS `click_to_call_session` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    initiated_by CHAR(36),
    from_phone VARCHAR(20),
    to_phone VARCHAR(20) NOT NULL,
    phone_type ENUM('AGENT', 'LEAD', 'CUSTOMER', 'INTERNAL', 'EXTERNAL') NOT NULL,
    contact_name VARCHAR(100),
    contact_email VARCHAR(100),
    contact_id CHAR(36),
    lead_id CHAR(36),
    agent_id CHAR(36),
    account_id CHAR(36),
    campaign_id CHAR(36),
    provider_id CHAR(36),
    provider_type ENUM('ASTERISK', 'SIP', 'MCUBE', 'EXOTEL', 'TWILIO', 'VONAGE', 'CUSTOM'),
    session_id VARCHAR(100),
    correlation_id VARCHAR(100),
    status ENUM('INITIATED', 'CONNECTING', 'RINGING', 'CONNECTED', 'DISCONNECTED', 'FAILED', 'COMPLETED') DEFAULT 'INITIATED',
    direction ENUM('INBOUND', 'OUTBOUND', 'INTERNAL') DEFAULT 'OUTBOUND',
    call_started_at TIMESTAMP NULL,
    call_ended_at TIMESTAMP NULL,
    duration_seconds INT DEFAULT 0,
    ring_time_seconds INT DEFAULT 0,
    answer_time_seconds INT DEFAULT 0,
    disconnect_reason VARCHAR(255),
    error_code VARCHAR(50),
    error_message TEXT,
    recording_url VARCHAR(500),
    transcript_url VARCHAR(500),
    call_quality_score DECIMAL(3,1),
    is_recorded TINYINT(1) DEFAULT 0,
    is_transferred TINYINT(1) DEFAULT 0,
    transfer_to_agent CHAR(36),
    notes TEXT,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_from_to_phone (from_phone, to_phone),
    INDEX idx_provider (provider_id),
    INDEX idx_session_id (session_id),
    INDEX idx_correlation_id (correlation_id),
    INDEX idx_agent_lead (agent_id, lead_id),
    INDEX idx_call_dates (call_started_at, call_ended_at),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES voip_provider(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Click-to-call session tracking and call details';

-- Table for call routing rules
CREATE TABLE IF NOT EXISTS `call_routing_rule` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    rule_name VARCHAR(100) NOT NULL,
    rule_type ENUM('AGENT_AVAILABILITY', 'SKILL_BASED', 'LOAD_BALANCING', 'TIME_BASED', 'PRIORITY', 'FAILOVER') NOT NULL,
    priority INT DEFAULT 0,
    is_active TINYINT(1) DEFAULT 1,
    condition_json JSON,
    action_json JSON,
    provider_id CHAR(36),
    fallback_provider_id CHAR(36),
    retry_on_failure TINYINT(1) DEFAULT 1,
    retry_count INT DEFAULT 3,
    retry_delay_seconds INT DEFAULT 5,
    notes TEXT,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_type (tenant_id, rule_type),
    INDEX idx_active (is_active),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES voip_provider(id) ON DELETE SET NULL,
    FOREIGN KEY (fallback_provider_id) REFERENCES voip_provider(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Call routing rules for intelligent call distribution';

-- Table for DTMF (Dual Tone Multi-Frequency) interactions
CREATE TABLE IF NOT EXISTS `call_dtmf_interaction` (
    id CHAR(36) PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    dtmf_digit VARCHAR(1) NOT NULL,
    dtmf_sequence VARCHAR(100),
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    action_triggered VARCHAR(100),
    notes TEXT,
    INDEX idx_session_dtmf (session_id),
    INDEX idx_tenant_session (tenant_id, session_id),
    FOREIGN KEY (session_id) REFERENCES click_to_call_session(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='DTMF (keypad press) interactions during calls';

-- Table for IVR (Interactive Voice Response) menu items
CREATE TABLE IF NOT EXISTS `ivr_menu` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    menu_name VARCHAR(100) NOT NULL,
    menu_code VARCHAR(20) NOT NULL,
    prompt_text TEXT,
    prompt_file_url VARCHAR(500),
    timeout_seconds INT DEFAULT 5,
    max_retries INT DEFAULT 3,
    parent_menu_id CHAR(36),
    is_active TINYINT(1) DEFAULT 1,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_tenant_menu_code (tenant_id, menu_code),
    INDEX idx_active (is_active),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_menu_id) REFERENCES ivr_menu(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='IVR menu definitions';

-- Table for IVR menu options
CREATE TABLE IF NOT EXISTS `ivr_menu_option` (
    id CHAR(36) PRIMARY KEY,
    menu_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    option_digit VARCHAR(1) NOT NULL,
    option_label VARCHAR(100),
    next_menu_id CHAR(36),
    action_type ENUM('ROUTE_AGENT', 'ROUTE_DEPARTMENT', 'PLAY_MESSAGE', 'COLLECT_DTMF', 'HANGUP', 'TRANSFER', 'VOICEMAIL') NOT NULL,
    action_target VARCHAR(100),
    priority INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_menu_digit (menu_id, option_digit),
    INDEX idx_menu_options (menu_id),
    FOREIGN KEY (menu_id) REFERENCES ivr_menu(id) ON DELETE CASCADE,
    FOREIGN KEY (next_menu_id) REFERENCES ivr_menu(id) ON DELETE SET NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='IVR menu option mappings';

-- Table for call recording configuration
CREATE TABLE IF NOT EXISTS `call_recording_config` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    record_all_calls TINYINT(1) DEFAULT 1,
    record_on_consent_only TINYINT(1) DEFAULT 0,
    retention_days INT DEFAULT 90,
    storage_location VARCHAR(255) COMMENT 'local, s3, gcs, azure',
    storage_bucket VARCHAR(255),
    encryption_enabled TINYINT(1) DEFAULT 1,
    encryption_method VARCHAR(50),
    quality_format ENUM('MP3', 'WAV', 'OGG', 'M4A') DEFAULT 'MP3',
    bitrate_kbps INT DEFAULT 128,
    sample_rate_hz INT DEFAULT 16000,
    auto_transcription TINYINT(1) DEFAULT 1,
    transcription_provider VARCHAR(50) COMMENT 'google, aws, azure, deepgram',
    transcription_lang VARCHAR(10) DEFAULT 'en-US',
    sentiment_analysis TINYINT(1) DEFAULT 1,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant (tenant_id),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Call recording configuration per tenant';

-- Table for call quality metrics
CREATE TABLE IF NOT EXISTS `call_quality_metric` (
    id CHAR(36) PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(8,2),
    unit VARCHAR(50),
    metric_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_session_metric (session_id, metric_name),
    INDEX idx_tenant_timestamp (tenant_id, metric_timestamp),
    FOREIGN KEY (session_id) REFERENCES click_to_call_session(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Call quality metrics (latency, jitter, packet loss, MOS)';

-- Table for call transfer history
CREATE TABLE IF NOT EXISTS `call_transfer` (
    id CHAR(36) PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    transferred_from_agent CHAR(36),
    transferred_to_agent CHAR(36),
    transfer_type ENUM('BLIND', 'ATTENDED', 'WARM', 'COLD') DEFAULT 'BLIND',
    transfer_reason VARCHAR(255),
    transfer_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    new_session_id CHAR(36),
    transfer_success TINYINT(1) DEFAULT 1,
    notes TEXT,
    INDEX idx_session_transfer (session_id),
    INDEX idx_tenant_session (tenant_id, session_id),
    FOREIGN KEY (session_id) REFERENCES click_to_call_session(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Call transfer and handoff history';

-- Table for agent activity log
CREATE TABLE IF NOT EXISTS `agent_activity_log` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    agent_id CHAR(36) NOT NULL,
    activity_type ENUM('LOGIN', 'LOGOUT', 'ON_CALL', 'ON_BREAK', 'ON_ADMIN', 'IDLE', 'AWAY', 'READY') NOT NULL,
    status_value VARCHAR(100),
    session_id CHAR(36),
    is_available TINYINT(1) DEFAULT 1,
    activity_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    duration_seconds INT,
    notes TEXT,
    INDEX idx_agent_activity (agent_id, activity_timestamp),
    INDEX idx_tenant_agent (tenant_id, agent_id),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Agent activity and status tracking for call center';

-- Table for phone number whitelist/blacklist
CREATE TABLE IF NOT EXISTS `phone_number_list` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    list_type ENUM('WHITELIST', 'BLACKLIST', 'VIP', 'SPAM') NOT NULL,
    contact_name VARCHAR(100),
    contact_email VARCHAR(100),
    reason TEXT,
    priority INT DEFAULT 0,
    is_active TINYINT(1) DEFAULT 1,
    expires_at TIMESTAMP NULL,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_phone (tenant_id, phone_number),
    INDEX idx_list_type (list_type),
    INDEX idx_active (is_active),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Phone number whitelist/blacklist management';

-- Table for caller ID management
CREATE TABLE IF NOT EXISTS `caller_id_profile` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    profile_name VARCHAR(100) NOT NULL,
    display_name VARCHAR(100),
    phone_number VARCHAR(20) NOT NULL,
    phone_country_code VARCHAR(5),
    is_default TINYINT(1) DEFAULT 0,
    is_active TINYINT(1) DEFAULT 1,
    verification_status ENUM('UNVERIFIED', 'PENDING', 'VERIFIED', 'REJECTED') DEFAULT 'UNVERIFIED',
    verified_at TIMESTAMP NULL,
    verified_by CHAR(36),
    notes TEXT,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_tenant_phone (tenant_id, phone_number),
    INDEX idx_default (is_default),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Caller ID profiles for outbound calls';

-- Table for click-to-call webhook logs
CREATE TABLE IF NOT EXISTS `click_to_call_webhook_log` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    provider_id CHAR(36),
    webhook_event_type VARCHAR(100),
    webhook_payload JSON NOT NULL,
    webhook_signature VARCHAR(500),
    is_valid TINYINT(1) DEFAULT 1,
    processing_status ENUM('PENDING', 'PROCESSED', 'FAILED', 'SKIPPED') DEFAULT 'PENDING',
    error_message TEXT,
    processed_at TIMESTAMP NULL,
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_event (tenant_id, webhook_event_type),
    INDEX idx_processing_status (processing_status),
    INDEX idx_received_at (received_at),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES voip_provider(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Webhook event logs from VoIP providers';

-- Table for call rules and compliance
CREATE TABLE IF NOT EXISTS `call_compliance_rule` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    rule_name VARCHAR(100) NOT NULL,
    rule_type ENUM('RECORDING_CONSENT', 'RECORDING_RETENTION', 'GDPR_DPO', 'CCPA_COMPLIANCE', 'TELEMARKETING', 'INDUSTRY_SPECIFIC') NOT NULL,
    rule_description TEXT,
    is_mandatory TINYINT(1) DEFAULT 1,
    is_active TINYINT(1) DEFAULT 1,
    enforced_since TIMESTAMP,
    rules_json JSON,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_type (tenant_id, rule_type),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Call compliance and regulatory rules';

-- Table for call rate configuration
CREATE TABLE IF NOT EXISTS `call_rate_config` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    destination_country VARCHAR(3),
    destination_carrier VARCHAR(100),
    destination_number_type ENUM('MOBILE', 'LANDLINE', 'INTERNATIONAL', 'SPECIAL') NOT NULL,
    rate_per_minute DECIMAL(8,4),
    setup_fee DECIMAL(8,2),
    minimum_charge DECIMAL(8,2),
    currency_code VARCHAR(3) DEFAULT 'USD',
    effective_from TIMESTAMP,
    effective_to TIMESTAMP NULL,
    is_active TINYINT(1) DEFAULT 1,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_country (tenant_id, destination_country),
    INDEX idx_active (is_active),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Call rate configuration for billing';

-- Table for call usage and billing
CREATE TABLE IF NOT EXISTS `call_usage_billing` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    billing_period_start DATE,
    billing_period_end DATE,
    total_calls INT DEFAULT 0,
    total_duration_minutes INT DEFAULT 0,
    total_charge DECIMAL(12,2) DEFAULT 0,
    inbound_calls INT DEFAULT 0,
    outbound_calls INT DEFAULT 0,
    setup_fees DECIMAL(12,2) DEFAULT 0,
    other_charges DECIMAL(12,2) DEFAULT 0,
    currency_code VARCHAR(3) DEFAULT 'USD',
    billing_status ENUM('DRAFT', 'PENDING', 'INVOICED', 'PAID', 'OVERDUE') DEFAULT 'DRAFT',
    invoice_id CHAR(36),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_period (tenant_id, billing_period_start),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Call usage and billing tracking';

-- Create indexes for common queries
CREATE INDEX idx_click_to_call_search ON click_to_call_session(tenant_id, created_at DESC);
CREATE INDEX idx_click_to_call_agent_date ON click_to_call_session(agent_id, created_at DESC);
CREATE INDEX idx_click_to_call_lead_date ON click_to_call_session(lead_id, created_at DESC);
