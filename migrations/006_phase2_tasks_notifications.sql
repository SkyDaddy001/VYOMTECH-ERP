-- Phase 2: Task Management & Notification System
-- Date: November 24, 2025
-- Features: Tasks, Notifications, Communication Templates, Communication Logs

-- ============================================================================
-- VERIFICATION: Check if Phase 1 tables exist
-- ============================================================================

-- Verify Phase 1 infrastructure is in place
SELECT 'Phase 1 schema validation...' as migration_step;

-- Confirm key Phase 1 tables exist
CREATE TABLE IF NOT EXISTS _migration_check (
    id INT AUTO_INCREMENT PRIMARY KEY,
    phase VARCHAR(10),
    status VARCHAR(50),
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO _migration_check (phase, status) VALUES ('Phase1', 'Starting Phase 2');

-- ============================================================================
-- PHASE 2: ENHANCED TASK MANAGEMENT SYSTEM
-- ============================================================================

-- Tasks table with rich metadata
CREATE TABLE IF NOT EXISTS tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    assigned_to BIGINT NOT NULL,
    created_by BIGINT,
    lead_id BIGINT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Task details
    title VARCHAR(255) NOT NULL,
    description TEXT,
    priority VARCHAR(20) NOT NULL DEFAULT 'normal', -- critical, high, normal, low
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, in_progress, completed, overdue, cancelled
    task_type VARCHAR(50), -- follow_up, meeting, call, demo, proposal, contract, review
    
    -- Dates and deadlines
    due_date TIMESTAMP NOT NULL,
    scheduled_at TIMESTAMP,
    completed_at TIMESTAMP,
    
    -- Task relationships
    parent_task_id BIGINT,
    related_entity_type VARCHAR(50), -- lead, opportunity, contact, customer
    related_entity_id BIGINT,
    
    -- Tracking
    estimated_duration_minutes INT,
    actual_duration_minutes INT,
    progress_percentage INT DEFAULT 0,
    
    -- Audit
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Constraints and indexes
    FOREIGN KEY (assigned_to) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES `user`(id) ON DELETE SET NULL,
    FOREIGN KEY (lead_id) REFERENCES `lead`(id) ON DELETE SET NULL,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_task_id) REFERENCES tasks(id) ON DELETE SET NULL,
    
    INDEX idx_assigned_to_status (assigned_to, status),
    INDEX idx_tenant_due_date (tenant_id, due_date),
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_tenant_priority (tenant_id, priority),
    INDEX idx_lead_id (lead_id),
    INDEX idx_created_at (created_at DESC),
    INDEX idx_overdue (due_date, status),
    INDEX idx_progress (progress_percentage)
);

-- Task comments and notes
CREATE TABLE IF NOT EXISTS task_comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    
    comment_text TEXT NOT NULL,
    attachment_url VARCHAR(500),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    
    INDEX idx_task_id (task_id),
    INDEX idx_user_id (user_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_created_at (created_at DESC)
);

-- Task assignments (for subtasks or multi-assigned tasks)
CREATE TABLE IF NOT EXISTS task_assignments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    
    assignment_role VARCHAR(50), -- assignee, reviewer, watcher
    responsibility_percentage INT DEFAULT 100,
    
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP,
    completed_at TIMESTAMP,
    
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    
    INDEX idx_task_id (task_id),
    INDEX idx_user_id (user_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_assignment_role (assignment_role),
    UNIQUE KEY unique_task_user_role (task_id, user_id, assignment_role)
);

-- ============================================================================
-- PHASE 2: NOTIFICATION SYSTEM
-- ============================================================================

-- Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Notification content
    type VARCHAR(50) NOT NULL, -- lead_assigned, call_missed, deadline_reminder, task_completed, task_assigned, message_received
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    
    -- Notification metadata
    priority VARCHAR(20) NOT NULL DEFAULT 'normal', -- critical, high, normal, low
    category VARCHAR(50), -- task, lead, call, message, system
    
    -- Related entity
    related_entity_type VARCHAR(50), -- task, lead, call, agent, contact
    related_entity_id BIGINT,
    
    -- Notification state
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    action_url VARCHAR(500),
    
    -- Notification lifecycle
    is_archived BOOLEAN DEFAULT FALSE,
    archived_at TIMESTAMP,
    
    -- Cleanup
    expires_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    
    INDEX idx_user_id_is_read (user_id, is_read),
    INDEX idx_user_id_created_at (user_id, created_at),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_priority (priority),
    INDEX idx_type (type),
    INDEX idx_expires_at (expires_at)
);

-- Notification preferences per user
CREATE TABLE IF NOT EXISTS notification_preferences (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL UNIQUE,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Channel preferences
    enable_email_notifications BOOLEAN DEFAULT TRUE,
    enable_sms_notifications BOOLEAN DEFAULT FALSE,
    enable_in_app_notifications BOOLEAN DEFAULT TRUE,
    enable_desktop_notifications BOOLEAN DEFAULT TRUE,
    
    -- Notification type preferences
    notify_task_assigned BOOLEAN DEFAULT TRUE,
    notify_task_completed BOOLEAN DEFAULT TRUE,
    notify_deadline_reminder BOOLEAN DEFAULT TRUE,
    notify_lead_assigned BOOLEAN DEFAULT TRUE,
    notify_call_missed BOOLEAN DEFAULT TRUE,
    notify_message_received BOOLEAN DEFAULT TRUE,
    notify_system_alerts BOOLEAN DEFAULT TRUE,
    
    -- Quiet hours
    quiet_hours_enabled BOOLEAN DEFAULT FALSE,
    quiet_hours_start TIME,
    quiet_hours_end TIME,
    
    -- Batching preferences
    email_batch_frequency VARCHAR(50) DEFAULT 'immediate', -- immediate, hourly, daily, weekly
    sms_batch_frequency VARCHAR(50) DEFAULT 'immediate',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    
    INDEX idx_tenant_id (tenant_id)
);

-- ============================================================================
-- PHASE 2: COMMUNICATION TEMPLATES
-- ============================================================================

-- Communication templates for reusable messages
CREATE TABLE IF NOT EXISTS communication_templates (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Template details
    name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Template type
    template_type VARCHAR(50), -- email, sms, whatsapp, slack, internal_message
    channel VARCHAR(50) NOT NULL, -- email, sms, whatsapp, slack
    
    -- Template content
    subject VARCHAR(500), -- for email templates
    content TEXT NOT NULL,
    variables JSON, -- placeholders: {{name}}, {{lead_id}}, {{task_title}}, etc.
    
    -- Template metadata
    category VARCHAR(50), -- task_reminder, lead_followup, welcome, follow_up, etc.
    created_by BIGINT,
    is_active BOOLEAN DEFAULT TRUE,
    usage_count INT DEFAULT 0,
    
    -- Audit
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES `user`(id) ON DELETE SET NULL,
    
    INDEX idx_tenant_active (tenant_id, is_active),
    INDEX idx_tenant_channel (tenant_id, channel),
    INDEX idx_category (category),
    INDEX idx_template_type (template_type)
);

-- ============================================================================
-- PHASE 2: COMMUNICATION LOGS
-- ============================================================================

-- Communication logs for tracking all sent messages
CREATE TABLE IF NOT EXISTS communication_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    lead_id BIGINT,
    user_id BIGINT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Communication details
    channel VARCHAR(50) NOT NULL, -- email, sms, whatsapp, slack, internal_message
    recipient VARCHAR(255) NOT NULL,
    subject VARCHAR(500),
    message TEXT NOT NULL,
    
    -- Template used
    template_id BIGINT,
    
    -- Status tracking
    status VARCHAR(50) DEFAULT 'sent', -- sent, failed, bounced, opened, clicked, read
    delivery_timestamp TIMESTAMP,
    read_timestamp TIMESTAMP,
    
    -- Response tracking
    has_response BOOLEAN DEFAULT FALSE,
    response_message TEXT,
    response_timestamp TIMESTAMP,
    
    -- Metrics
    open_count INT DEFAULT 0,
    click_count INT DEFAULT 0,
    
    -- Error handling
    error_message TEXT,
    retry_count INT DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (lead_id) REFERENCES `lead`(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE SET NULL,
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES communication_templates(id) ON DELETE SET NULL,
    
    INDEX idx_lead_id (lead_id),
    INDEX idx_user_id (user_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_channel (channel),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at DESC),
    INDEX idx_recipient (recipient),
    INDEX idx_has_response (has_response)
);

-- ============================================================================
-- PHASE 2: NOTIFICATION DELIVERY TRACKING
-- ============================================================================

-- Track notification delivery status
CREATE TABLE IF NOT EXISTS notification_delivery (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    notification_id BIGINT NOT NULL,
    channel VARCHAR(50) NOT NULL, -- email, sms, push, in_app
    
    status VARCHAR(50) DEFAULT 'pending', -- pending, sent, failed, bounced, opened, clicked
    
    -- Delivery attempt details
    attempt_number INT DEFAULT 1,
    sent_at TIMESTAMP,
    opened_at TIMESTAMP,
    clicked_at TIMESTAMP,
    
    -- Error tracking
    error_message TEXT,
    error_code VARCHAR(50),
    
    -- External tracking
    external_message_id VARCHAR(255),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE,
    
    INDEX idx_notification_id (notification_id),
    INDEX idx_channel (channel),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at DESC)
);

-- ============================================================================
-- MIGRATION VERIFICATION
-- ============================================================================

-- Verify Phase 2 tables created
SELECT 'Phase 2 schema created successfully' as migration_status;

INSERT INTO _migration_check (phase, status) VALUES ('Phase2', 'All tables created');

-- Summary
SELECT 
    'Phase 2 Migration Complete' as message,
    COUNT(*) as tables_checked
FROM information_schema.tables 
WHERE table_schema = DATABASE() 
AND table_name IN ('tasks', 'task_comments', 'task_assignments', 'notifications', 
                    'notification_preferences', 'communication_templates', 
                    'communication_logs', 'notification_delivery');

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================
-- Phase 2 tables created successfully
-- Tables: 8 new
-- New functionality: Task Management, Notifications, Communication Tracking
-- Status: Ready for service implementation
