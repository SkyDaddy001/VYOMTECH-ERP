-- Phase 1 Critical Features Migration
-- Agent Availability, Lead Scoring, Audit Trail
-- Date: November 22, 2025

-- ============================================================================
-- 1. AGENT AVAILABILITY MANAGEMENT
-- ============================================================================

CREATE TABLE IF NOT EXISTS agent_availability (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL UNIQUE,
    tenant_id VARCHAR(36) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'available',
    -- Status enum: available, busy, on_break, offline, in_meeting, away
    break_reason VARCHAR(255),
    is_accepting_leads BOOLEAN DEFAULT TRUE,
    total_calls_today INT DEFAULT 0,
    current_call_duration_seconds INT DEFAULT 0,
    last_status_change TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_accepting_leads (is_accepting_leads),
    INDEX idx_last_activity (last_activity)
);

-- ============================================================================
-- 2. LEAD SCORING SYSTEM
-- ============================================================================

CREATE TABLE IF NOT EXISTS lead_scores (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    lead_id BIGINT NOT NULL UNIQUE,
    tenant_id VARCHAR(36) NOT NULL,
    source_quality_score DECIMAL(5,2) DEFAULT 0.0,
    engagement_score DECIMAL(5,2) DEFAULT 0.0,
    conversion_probability DECIMAL(5,2) DEFAULT 0.0,
    urgency_score DECIMAL(5,2) DEFAULT 0.0,
    overall_score DECIMAL(7,2) DEFAULT 0.0,
    score_category VARCHAR(50), -- hot, warm, cold, nurture
    previous_score DECIMAL(7,2),
    score_change DECIMAL(7,2),
    reason_text TEXT,
    calculation_method VARCHAR(100) DEFAULT 'weighted',
    last_calculated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (lead_id) REFERENCES `lead`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    INDEX idx_tenant_overall_score (tenant_id, overall_score DESC),
    INDEX idx_tenant_category (tenant_id, score_category),
    INDEX idx_last_calculated (last_calculated)
);

-- ============================================================================
-- 3. AUDIT TRAIL SYSTEM
-- ============================================================================

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT,
    tenant_id VARCHAR(36) NOT NULL,
    action VARCHAR(255) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    old_values JSON,
    new_values JSON,
    ip_address VARCHAR(45),
    user_agent TEXT,
    status VARCHAR(50) DEFAULT 'success', -- success, failure
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE SET NULL,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_tenant_entity (tenant_id, entity_type, entity_id),
    INDEX idx_created_at (created_at DESC),
    INDEX idx_action (action),
    INDEX idx_status (status)
);

-- ============================================================================
-- 4. LEAD ACTIVITY TIMELINE
-- ============================================================================

CREATE TABLE IF NOT EXISTS lead_activities (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    lead_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    activity_type VARCHAR(50) NOT NULL,
    -- call, email, meeting, follow_up, note, status_change, assignment
    description TEXT NOT NULL,
    created_by BIGINT,
    activity_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    duration_minutes INT,
    outcome VARCHAR(100),
    next_action VARCHAR(255),
    next_follow_up TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (lead_id) REFERENCES `lead`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES `user`(id) ON DELETE SET NULL,
    INDEX idx_lead_activity (lead_id, activity_date),
    INDEX idx_tenant_activity (tenant_id, activity_type),
    INDEX idx_activity_date (activity_date DESC),
    INDEX idx_next_follow_up (next_follow_up)
);

-- ============================================================================
-- 5. TASK MANAGEMENT SYSTEM
-- ============================================================================

CREATE TABLE IF NOT EXISTS tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    assigned_to BIGINT NOT NULL,
    created_by BIGINT,
    lead_id BIGINT,
    tenant_id VARCHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    priority VARCHAR(20) DEFAULT 'normal', -- critical, high, normal, low
    status VARCHAR(50) DEFAULT 'pending', -- pending, in_progress, completed, overdue, cancelled
    due_date TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (assigned_to) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES `user`(id) ON DELETE SET NULL,
    FOREIGN KEY (lead_id) REFERENCES `lead`(id) ON DELETE SET NULL,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    INDEX idx_assigned_to_status (assigned_to, status),
    INDEX idx_tenant_due_date (tenant_id, due_date),
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_lead_id (lead_id)
);

-- ============================================================================
-- 6. NOTIFICATION SYSTEM
-- ============================================================================

CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    -- lead_assigned, call_missed, deadline_reminder, task_completed, lead_scored_hot
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    related_entity_id VARCHAR(255),
    is_read BOOLEAN DEFAULT FALSE,
    priority VARCHAR(20) DEFAULT 'normal', -- critical, high, normal, low
    read_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    INDEX idx_user_is_read (user_id, is_read),
    INDEX idx_tenant_type (tenant_id, type),
    INDEX idx_priority (priority),
    INDEX idx_created_at (created_at DESC)
);

-- ============================================================================
-- 7. COMMUNICATION TEMPLATES
-- ============================================================================

CREATE TABLE IF NOT EXISTS communication_templates (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    channel VARCHAR(50) NOT NULL, -- email, sms, whatsapp, slack
    content TEXT NOT NULL,
    variables JSON, -- placeholders like {{name}}, {{lead_id}}
    created_by BIGINT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES `user`(id) ON DELETE SET NULL,
    INDEX idx_tenant_channel (tenant_id, channel),
    INDEX idx_tenant_active (tenant_id, is_active)
);

-- ============================================================================
-- 8. COMMUNICATION LOGS
-- ============================================================================

CREATE TABLE IF NOT EXISTS communication_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    lead_id BIGINT,
    user_id BIGINT,
    tenant_id VARCHAR(36) NOT NULL,
    channel VARCHAR(50) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    template_id BIGINT,
    status VARCHAR(50) DEFAULT 'sent', -- sent, failed, delivered, read, bounced
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivered_at TIMESTAMP,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (lead_id) REFERENCES `lead`(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE SET NULL,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES communication_templates(id) ON DELETE SET NULL,
    INDEX idx_lead_channel (lead_id, channel),
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_sent_at (sent_at DESC)
);

-- ============================================================================
-- 9. ANALYTICS DAILY METRICS
-- ============================================================================

CREATE TABLE IF NOT EXISTS analytics_daily (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    date DATE NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    team_id BIGINT,
    total_leads_created INT DEFAULT 0,
    total_calls_made INT DEFAULT 0,
    total_conversions INT DEFAULT 0,
    average_call_duration INT DEFAULT 0,
    conversion_rate DECIMAL(5,2) DEFAULT 0.0,
    average_lead_score DECIMAL(7,2) DEFAULT 0.0,
    hot_leads_count INT DEFAULT 0,
    warm_leads_count INT DEFAULT 0,
    cold_leads_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES team(id) ON DELETE SET NULL,
    UNIQUE KEY unique_date_tenant_team (date, tenant_id, team_id),
    INDEX idx_tenant_date (tenant_id, date DESC)
);

-- ============================================================================
-- 10. TWO-FACTOR AUTHENTICATION
-- ============================================================================

CREATE TABLE IF NOT EXISTS two_factor_codes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    INDEX idx_user_created_at (user_id, created_at),
    INDEX idx_expires_at (expires_at)
);

CREATE TABLE IF NOT EXISTS user_2fa_settings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL UNIQUE,
    is_enabled BOOLEAN DEFAULT FALSE,
    method VARCHAR(50), -- sms, email, authenticator
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE
);

-- ============================================================================
-- CREATE INDEXES AND KEYS
-- ============================================================================

-- Ensure tenant_id is indexed on all new tables for multi-tenancy

-- =================================================================
-- MIGRATION COMPLETE
-- ============================================================================
-- Phase 1 tables created successfully
-- Tables: 10
-- Records: agent_availability, lead_scores, audit_logs, lead_activities, tasks
--          notifications, communication_templates, communication_logs, analytics_daily
--          two_factor_codes, user_2fa_settings
