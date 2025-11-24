-- Multi-Tenant AI Call Center Database Schema
-- Migration: 001_initial_schema.sql

-- Enable foreign key constraints
SET FOREIGN_KEY_CHECKS = 1;

-- Tenants table
CREATE TABLE IF NOT EXISTS tenant (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE,
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
    max_users INT DEFAULT 100,
    max_concurrent_calls INT DEFAULT 50,
    ai_budget_monthly DECIMAL(10,2) DEFAULT 1000.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_domain (domain)
);

-- Users table
CREATE TABLE IF NOT EXISTS user (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('admin', 'agent', 'supervisor', 'user') NOT NULL DEFAULT 'user',
    status ENUM('active', 'inactive') DEFAULT 'active',
    tenant_id VARCHAR(36) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    INDEX idx_email (email),
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_role (role)
);

-- Password reset tokens
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    token VARCHAR(64) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    INDEX idx_token (token),
    INDEX idx_user_expires (user_id, expires_at)
);

-- Agents table (extends users)
CREATE TABLE IF NOT EXISTS agent (
    user_id BIGINT PRIMARY KEY,
    status ENUM('active', 'inactive') DEFAULT 'active',
    availability ENUM('online', 'offline', 'busy') DEFAULT 'offline',
    skills JSON,
    max_concurrent_calls INT DEFAULT 3,
    current_calls INT DEFAULT 0,
    total_calls INT DEFAULT 0,
    avg_handle_time DECIMAL(5,2) DEFAULT 0.00,
    satisfaction_score DECIMAL(3,2) DEFAULT 0.00,
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    tenant_id VARCHAR(36) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_availability (availability),
    INDEX idx_tenant_availability (tenant_id, availability)
);

-- Leads table
CREATE TABLE IF NOT EXISTS `lead` (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20),
    company VARCHAR(255),
    source VARCHAR(100),
    status ENUM('new', 'contacted', 'qualified', 'converted', 'lost') DEFAULT 'new',
    priority ENUM('low', 'medium', 'high') DEFAULT 'medium',
    assigned_agent_id BIGINT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_agent_id) REFERENCES `user`(id) ON DELETE SET NULL,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_assigned_agent (assigned_agent_id),
    INDEX idx_created_at (created_at),
    INDEX idx_source (source)
);

-- Calls table
CREATE TABLE IF NOT EXISTS `call` (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    lead_id BIGINT,
    direction ENUM('inbound', 'outbound') NOT NULL,
    status ENUM('ringing', 'connected', 'completed', 'failed', 'transferred') DEFAULT 'ringing',
    phone_number VARCHAR(20) NOT NULL,
    agent_id BIGINT,
    ai_used BOOLEAN DEFAULT FALSE,
    ai_provider VARCHAR(50),
    duration_seconds INT DEFAULT 0,
    recording_url VARCHAR(500),
    transcription TEXT,
    sentiment_score DECIMAL(3,2),
    notes TEXT,
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (lead_id) REFERENCES `lead`(id) ON DELETE SET NULL,
    FOREIGN KEY (agent_id) REFERENCES `user`(id) ON DELETE SET NULL,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_agent_id (agent_id),
    INDEX idx_started_at (started_at),
    INDEX idx_direction (direction)
);

-- AI Request Log
CREATE TABLE IF NOT EXISTS ai_request_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    query TEXT NOT NULL,
    provider VARCHAR(50) NOT NULL,
    tokens_used INT DEFAULT 0,
    processing_time_ms INT DEFAULT 0,
    cost DECIMAL(8,4) DEFAULT 0.0000,
    priority ENUM('low', 'medium', 'high') DEFAULT 'medium',
    cached BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    INDEX idx_tenant_created (tenant_id, created_at),
    INDEX idx_provider (provider),
    INDEX idx_priority (priority)
);

-- Campaigns table
CREATE TABLE IF NOT EXISTS campaign (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type ENUM('email', 'sms', 'call', 'social') NOT NULL,
    status ENUM('draft', 'scheduled', 'running', 'completed', 'paused') DEFAULT 'draft',
    target_audience JSON,
    content TEXT,
    scheduled_at TIMESTAMP,
    completed_at TIMESTAMP,
    total_recipients INT DEFAULT 0,
    sent_count INT DEFAULT 0,
    open_count INT DEFAULT 0,
    click_count INT DEFAULT 0,
    conversion_count INT DEFAULT 0,
    budget DECIMAL(10,2),
    cost DECIMAL(10,2) DEFAULT 0.00,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES `user`(id) ON DELETE CASCADE,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_type (type),
    INDEX idx_scheduled_at (scheduled_at)
);

-- Campaign Recipients
CREATE TABLE IF NOT EXISTS campaign_recipient (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    campaign_id BIGINT NOT NULL,
    lead_id BIGINT,
    email VARCHAR(255),
    phone VARCHAR(20),
    status ENUM('pending', 'sent', 'delivered', 'opened', 'clicked', 'converted', 'failed') DEFAULT 'pending',
    sent_at TIMESTAMP,
    opened_at TIMESTAMP,
    clicked_at TIMESTAMP,
    error_message TEXT,
    FOREIGN KEY (campaign_id) REFERENCES campaign(id) ON DELETE CASCADE,
    FOREIGN KEY (lead_id) REFERENCES `lead`(id) ON DELETE SET NULL,
    INDEX idx_campaign_status (campaign_id, status),
    INDEX idx_lead_id (lead_id)
);

-- Marketing Attribution
CREATE TABLE IF NOT EXISTS marketing_attribution (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    lead_id BIGINT,
    source VARCHAR(100) NOT NULL,
    campaign_id BIGINT,
    channel VARCHAR(50),
    cost DECIMAL(8,2) DEFAULT 0.00,
    conversions INT DEFAULT 0,
    revenue DECIMAL(10,2) DEFAULT 0.00,
    roi DECIMAL(5,2) DEFAULT 0.00,
    attribution_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (lead_id) REFERENCES `lead`(id) ON DELETE SET NULL,
    FOREIGN KEY (campaign_id) REFERENCES campaign(id) ON DELETE SET NULL,
    INDEX idx_tenant_date (tenant_id, attribution_date),
    INDEX idx_source (source),
    INDEX idx_campaign (campaign_id)
);

-- Settings table for tenant-specific configurations
CREATE TABLE IF NOT EXISTS tenant_settings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    setting_key VARCHAR(100) NOT NULL,
    setting_value JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    UNIQUE KEY unique_tenant_setting (tenant_id, setting_key),
    INDEX idx_tenant_key (tenant_id, setting_key)
);

-- Insert default tenant
INSERT IGNORE INTO tenant (id, name, domain, status) VALUES
('default-tenant', 'Default Tenant', 'default.callcenter.com', 'active');

-- Insert default admin user (password: admin123 - CHANGE THIS IN PRODUCTION!)
INSERT IGNORE INTO user (email, password_hash, role, tenant_id, first_name, last_name) VALUES
('admin@callcenter.com', '$2a$10$8K3VZ6Yx5Z8Yx5Z8Yx5Z8Yx5Z8Yx5Z8Yx5Z8Yx5Z8Yx5Z8Yx5Z8Y', 'admin', 'default-tenant', 'System', 'Administrator');
