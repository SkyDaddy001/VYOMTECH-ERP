-- Migration: 020_multi_channel_communication.sql
-- Description: Multi-channel communication system (Telegram, WhatsApp, SMS, Email)
-- Date: 2025-12-03

-- Table for communication channel configuration
CREATE TABLE IF NOT EXISTS `communication_channel` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    channel_type ENUM('TELEGRAM', 'WHATSAPP', 'SMS', 'EMAIL') NOT NULL,
    channel_name VARCHAR(100) NOT NULL,
    api_provider VARCHAR(100) NOT NULL COMMENT 'twilio, sendgrid, mailgun, telegram, whatsapp_business, etc',
    api_key VARCHAR(500),
    api_secret VARCHAR(500),
    api_url VARCHAR(255),
    webhook_url VARCHAR(500),
    callback_url VARCHAR(500),
    auth_token VARCHAR(500),
    account_id VARCHAR(100),
    sender_id VARCHAR(100) COMMENT 'Phone number, email, telegram bot id, whatsapp number',
    is_active TINYINT(1) DEFAULT 1,
    retry_count INT DEFAULT 3,
    timeout_seconds INT DEFAULT 30,
    priority INT DEFAULT 0,
    config_json JSON,
    notes TEXT,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_channel (tenant_id, channel_type),
    INDEX idx_active (is_active),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Communication channel configuration';

-- Table for message templates
CREATE TABLE IF NOT EXISTS `message_template` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    template_name VARCHAR(100) NOT NULL,
    channel_type ENUM('TELEGRAM', 'WHATSAPP', 'SMS', 'EMAIL') NOT NULL,
    template_body TEXT NOT NULL,
    template_variables JSON,
    subject VARCHAR(255) COMMENT 'For email templates',
    language VARCHAR(10) DEFAULT 'en',
    is_active TINYINT(1) DEFAULT 1,
    usage_count INT DEFAULT 0,
    last_used_at TIMESTAMP NULL,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_tenant_template (tenant_id, template_name, channel_type),
    INDEX idx_active (is_active),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Message templates for different channels';

-- Table for communication sessions/conversations
CREATE TABLE IF NOT EXISTS `communication_session` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    initiated_by CHAR(36),
    channel_type ENUM('TELEGRAM', 'WHATSAPP', 'SMS', 'EMAIL') NOT NULL,
    channel_id CHAR(36),
    sender_id VARCHAR(255) COMMENT 'Email, phone, telegram_id, whatsapp_number',
    recipient_id VARCHAR(255) NOT NULL COMMENT 'Email, phone, telegram_id, whatsapp_number',
    recipient_name VARCHAR(100),
    recipient_email VARCHAR(100),
    contact_id CHAR(36),
    lead_id CHAR(36),
    agent_id CHAR(36),
    account_id CHAR(36),
    campaign_id CHAR(36),
    conversation_id VARCHAR(100),
    correlation_id VARCHAR(100),
    status ENUM('INITIATED', 'SENT', 'DELIVERED', 'READ', 'FAILED', 'COMPLETED') DEFAULT 'INITIATED',
    direction ENUM('INBOUND', 'OUTBOUND', 'INTERNAL') DEFAULT 'OUTBOUND',
    message_count INT DEFAULT 0,
    last_message_at TIMESTAMP NULL,
    first_message_at TIMESTAMP NULL,
    last_response_at TIMESTAMP NULL,
    response_time_minutes INT,
    is_archived TINYINT(1) DEFAULT 0,
    is_starred TINYINT(1) DEFAULT 0,
    priority INT DEFAULT 0,
    notes TEXT,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_sender_recipient (sender_id, recipient_id),
    INDEX idx_channel (channel_type),
    INDEX idx_agent_lead (agent_id, lead_id),
    INDEX idx_session_dates (first_message_at, last_message_at),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (channel_id) REFERENCES communication_channel(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Communication sessions/conversations';

-- Table for individual messages
CREATE TABLE IF NOT EXISTS `communication_message` (
    id CHAR(36) PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    channel_type ENUM('TELEGRAM', 'WHATSAPP', 'SMS', 'EMAIL') NOT NULL,
    from_address VARCHAR(255) NOT NULL,
    to_address VARCHAR(255) NOT NULL,
    message_type ENUM('TEXT', 'IMAGE', 'VIDEO', 'AUDIO', 'FILE', 'LOCATION', 'TEMPLATE') DEFAULT 'TEXT',
    message_subject VARCHAR(255) COMMENT 'For email messages',
    message_body TEXT NOT NULL,
    message_html TEXT COMMENT 'For email messages with HTML',
    media_url VARCHAR(500),
    media_size_bytes INT,
    media_type VARCHAR(100),
    attachments JSON,
    template_id CHAR(36),
    template_variables JSON,
    external_message_id VARCHAR(100),
    status ENUM('QUEUED', 'SENT', 'DELIVERED', 'READ', 'FAILED', 'BOUNCED') DEFAULT 'QUEUED',
    error_code VARCHAR(50),
    error_message TEXT,
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    sent_at TIMESTAMP NULL,
    delivered_at TIMESTAMP NULL,
    read_at TIMESTAMP NULL,
    read_by VARCHAR(255),
    cost DECIMAL(10,4) DEFAULT 0,
    cost_currency VARCHAR(3) DEFAULT 'USD',
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_session_message (session_id),
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_from_to (from_address, to_address),
    INDEX idx_message_dates (sent_at, delivered_at),
    FOREIGN KEY (session_id) REFERENCES communication_session(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES message_template(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Individual messages in conversations';

-- Table for contact preferences
CREATE TABLE IF NOT EXISTS `contact_communication_preference` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    contact_id CHAR(36),
    lead_id CHAR(36),
    email_address VARCHAR(100),
    phone_number VARCHAR(20),
    telegram_id VARCHAR(100),
    whatsapp_number VARCHAR(20),
    preferred_channel ENUM('TELEGRAM', 'WHATSAPP', 'SMS', 'EMAIL') DEFAULT 'EMAIL',
    allow_telegram TINYINT(1) DEFAULT 1,
    allow_whatsapp TINYINT(1) DEFAULT 1,
    allow_sms TINYINT(1) DEFAULT 1,
    allow_email TINYINT(1) DEFAULT 1,
    opt_in_sms TINYINT(1) DEFAULT 0,
    opt_in_marketing TINYINT(1) DEFAULT 0,
    opt_out_date TIMESTAMP NULL,
    do_not_contact TINYINT(1) DEFAULT 0,
    do_not_contact_until TIMESTAMP NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_contact (tenant_id, contact_id),
    INDEX idx_lead (lead_id),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Contact communication preferences';

-- Table for webhook logs
CREATE TABLE IF NOT EXISTS `communication_webhook_log` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    channel_id CHAR(36),
    channel_type ENUM('TELEGRAM', 'WHATSAPP', 'SMS', 'EMAIL'),
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
    FOREIGN KEY (channel_id) REFERENCES communication_channel(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Webhook event logs from communication providers';

-- Table for campaign/bulk messaging
CREATE TABLE IF NOT EXISTS `bulk_message_campaign` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    campaign_name VARCHAR(100) NOT NULL,
    channel_type ENUM('TELEGRAM', 'WHATSAPP', 'SMS', 'EMAIL') NOT NULL,
    campaign_description TEXT,
    template_id CHAR(36),
    recipient_filter JSON,
    total_recipients INT DEFAULT 0,
    sent_count INT DEFAULT 0,
    delivered_count INT DEFAULT 0,
    failed_count INT DEFAULT 0,
    campaign_status ENUM('DRAFT', 'SCHEDULED', 'RUNNING', 'COMPLETED', 'CANCELLED') DEFAULT 'DRAFT',
    scheduled_at TIMESTAMP NULL,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    estimated_cost DECIMAL(12,2),
    actual_cost DECIMAL(12,2) DEFAULT 0,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_status (tenant_id, campaign_status),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES message_template(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Bulk messaging campaigns';

-- Table for campaign recipients
CREATE TABLE IF NOT EXISTS `bulk_message_recipient` (
    id CHAR(36) PRIMARY KEY,
    campaign_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    recipient_address VARCHAR(255) NOT NULL,
    recipient_name VARCHAR(100),
    contact_id CHAR(36),
    lead_id CHAR(36),
    template_variables JSON,
    send_status ENUM('PENDING', 'SENT', 'DELIVERED', 'FAILED', 'SKIPPED') DEFAULT 'PENDING',
    error_code VARCHAR(50),
    error_message TEXT,
    message_id CHAR(36),
    sent_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_campaign_status (campaign_id, send_status),
    FOREIGN KEY (campaign_id) REFERENCES bulk_message_campaign(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Bulk campaign recipients';

-- Table for message automation rules
CREATE TABLE IF NOT EXISTS `message_automation_rule` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    rule_name VARCHAR(100) NOT NULL,
    rule_type ENUM('TRIGGER_ON_EVENT', 'SCHEDULED', 'WORKFLOW', 'DRIP_CAMPAIGN') NOT NULL,
    trigger_event VARCHAR(100),
    channel_type ENUM('TELEGRAM', 'WHATSAPP', 'SMS', 'EMAIL') NOT NULL,
    template_id CHAR(36),
    condition_json JSON,
    action_json JSON,
    is_active TINYINT(1) DEFAULT 1,
    priority INT DEFAULT 0,
    execution_count INT DEFAULT 0,
    notes TEXT,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_active (tenant_id, is_active),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES message_template(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Message automation rules';

-- Table for communication analytics
CREATE TABLE IF NOT EXISTS `communication_analytics` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    channel_type ENUM('TELEGRAM', 'WHATSAPP', 'SMS', 'EMAIL'),
    metric_date DATE,
    total_messages INT DEFAULT 0,
    sent_messages INT DEFAULT 0,
    delivered_messages INT DEFAULT 0,
    failed_messages INT DEFAULT 0,
    read_messages INT DEFAULT 0,
    total_cost DECIMAL(12,2) DEFAULT 0,
    avg_delivery_time_seconds INT,
    avg_response_time_minutes INT,
    unique_recipients INT DEFAULT 0,
    engagement_rate DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_date (tenant_id, metric_date),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Communication channel analytics';

-- Table for message scheduling
CREATE TABLE IF NOT EXISTS `scheduled_message` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    channel_type ENUM('TELEGRAM', 'WHATSAPP', 'SMS', 'EMAIL') NOT NULL,
    from_address VARCHAR(255),
    to_address VARCHAR(255) NOT NULL,
    template_id CHAR(36),
    message_body TEXT,
    scheduled_for TIMESTAMP NOT NULL,
    scheduled_timezone VARCHAR(50),
    recurrence_pattern VARCHAR(100) COMMENT 'daily, weekly, monthly, once',
    recurrence_end_date TIMESTAMP NULL,
    status ENUM('SCHEDULED', 'SENT', 'CANCELLED', 'FAILED') DEFAULT 'SCHEDULED',
    last_sent_at TIMESTAMP NULL,
    next_send_at TIMESTAMP NULL,
    notes TEXT,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_scheduled_for (scheduled_for),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES message_template(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Scheduled messages for later delivery';

-- Table for conversation attachments
CREATE TABLE IF NOT EXISTS `communication_attachment` (
    id CHAR(36) PRIMARY KEY,
    message_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    attachment_type ENUM('IMAGE', 'VIDEO', 'AUDIO', 'FILE', 'DOCUMENT') NOT NULL,
    file_name VARCHAR(255),
    file_size_bytes INT,
    mime_type VARCHAR(100),
    file_url VARCHAR(500),
    storage_path VARCHAR(500),
    storage_location ENUM('local', 's3', 'gcs', 'azure') DEFAULT 'local',
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_message_attachment (message_id),
    FOREIGN KEY (message_id) REFERENCES communication_message(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Message attachments';

-- Table for user communication permissions
CREATE TABLE IF NOT EXISTS `user_communication_permission` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    can_send_telegram TINYINT(1) DEFAULT 0,
    can_send_whatsapp TINYINT(1) DEFAULT 0,
    can_send_sms TINYINT(1) DEFAULT 0,
    can_send_email TINYINT(1) DEFAULT 1,
    can_view_conversations TINYINT(1) DEFAULT 1,
    can_manage_templates TINYINT(1) DEFAULT 0,
    can_manage_campaigns TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_perms (tenant_id, user_id),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='User permissions for communication channels';

-- Create indexes for common queries
CREATE INDEX idx_communication_search ON communication_session(tenant_id, created_at DESC);
CREATE INDEX idx_communication_agent_date ON communication_session(agent_id, created_at DESC);
CREATE INDEX idx_communication_lead_date ON communication_session(lead_id, created_at DESC);
CREATE INDEX idx_message_session_date ON communication_message(session_id, created_at DESC);
