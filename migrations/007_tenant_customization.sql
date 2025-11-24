-- Phase 2B: Tenant-Level Customization
-- Task Statuses, Stages, Notification Types, and Workflow Management
-- Date: November 24, 2025

-- ============================================================================
-- TENANT CUSTOMIZATION: TASK WORKFLOW
-- ============================================================================

-- Custom task statuses per tenant
CREATE TABLE IF NOT EXISTS tenant_task_statuses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Status definition
    status_code VARCHAR(50) NOT NULL,
    status_name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Status properties
    color_hex VARCHAR(7),
    icon VARCHAR(50),
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Status behavior
    is_initial_status BOOLEAN DEFAULT FALSE, -- First status for new tasks
    is_final_status BOOLEAN DEFAULT FALSE,   -- Task is done/closed
    is_blocking_status BOOLEAN DEFAULT FALSE, -- Blocks certain operations
    allows_editing BOOLEAN DEFAULT TRUE,
    allows_reassignment BOOLEAN DEFAULT TRUE,
    
    -- Audit
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES `user`(id) ON DELETE SET NULL,
    
    UNIQUE KEY unique_tenant_status (tenant_id, status_code),
    INDEX idx_tenant_active (tenant_id, is_active),
    INDEX idx_display_order (display_order)
);

-- Custom task stages/pipelines per tenant
CREATE TABLE IF NOT EXISTS tenant_task_stages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Stage definition
    stage_code VARCHAR(50) NOT NULL,
    stage_name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Stage properties
    color_hex VARCHAR(7),
    icon VARCHAR(50),
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Stage behavior
    min_duration_hours INT,
    max_duration_hours INT,
    sla_minutes INT, -- Service level agreement response time
    auto_advance_to_stage_id BIGINT, -- Auto-advance after condition
    
    -- Audit
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES `user`(id) ON DELETE SET NULL,
    FOREIGN KEY (auto_advance_to_stage_id) REFERENCES tenant_task_stages(id) ON DELETE SET NULL,
    
    UNIQUE KEY unique_tenant_stage (tenant_id, stage_code),
    INDEX idx_tenant_active (tenant_id, is_active),
    INDEX idx_display_order (display_order)
);

-- Status transitions allowed per tenant
CREATE TABLE IF NOT EXISTS tenant_status_transitions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Transition definition
    from_status_code VARCHAR(50) NOT NULL,
    to_status_code VARCHAR(50) NOT NULL,
    
    -- Transition properties
    is_allowed BOOLEAN DEFAULT TRUE,
    requires_comment BOOLEAN DEFAULT FALSE,
    requires_approval BOOLEAN DEFAULT FALSE,
    notification_on_transition BOOLEAN DEFAULT FALSE,
    
    -- Conditions
    requires_role VARCHAR(100), -- comma-separated roles that can make this transition
    requires_field_completion VARCHAR(500), -- JSON array of required fields
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    
    UNIQUE KEY unique_transition (tenant_id, from_status_code, to_status_code),
    INDEX idx_tenant_from_status (tenant_id, from_status_code)
);

-- ============================================================================
-- TENANT CUSTOMIZATION: TASK TYPES
-- ============================================================================

-- Custom task types per tenant
CREATE TABLE IF NOT EXISTS tenant_task_types (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Task type definition
    type_code VARCHAR(50) NOT NULL,
    type_name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Type properties
    icon VARCHAR(50),
    color_hex VARCHAR(7),
    default_priority VARCHAR(20) DEFAULT 'normal',
    default_due_days INT,
    
    -- Workflow
    required_statuses VARCHAR(500), -- JSON array of required statuses
    is_lead_related BOOLEAN DEFAULT TRUE,
    is_agent_assignable BOOLEAN DEFAULT TRUE,
    
    -- Audit
    is_active BOOLEAN DEFAULT TRUE,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES `user`(id) ON DELETE SET NULL,
    
    UNIQUE KEY unique_tenant_type (tenant_id, type_code),
    INDEX idx_tenant_active (tenant_id, is_active)
);

-- ============================================================================
-- TENANT CUSTOMIZATION: NOTIFICATION TYPES
-- ============================================================================

-- Custom notification types per tenant
CREATE TABLE IF NOT EXISTS tenant_notification_types (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Notification type definition
    type_code VARCHAR(50) NOT NULL,
    type_name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Type properties
    icon VARCHAR(50),
    color_hex VARCHAR(7),
    default_priority VARCHAR(20) DEFAULT 'normal',
    category VARCHAR(50), -- task, lead, call, message, system, custom
    
    -- Channel configuration
    supported_channels VARCHAR(200), -- JSON array: email, sms, push, in_app, slack
    default_channels VARCHAR(200), -- JSON array of channels enabled by default
    
    -- Behavior
    is_dismissable BOOLEAN DEFAULT TRUE,
    auto_archive_after_days INT,
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Audit
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES `user`(id) ON DELETE SET NULL,
    
    UNIQUE KEY unique_tenant_notif_type (tenant_id, type_code),
    INDEX idx_tenant_active (tenant_id, is_active)
);

-- ============================================================================
-- TENANT CUSTOMIZATION: PRIORITIES
-- ============================================================================

-- Custom priority levels per tenant
CREATE TABLE IF NOT EXISTS tenant_priority_levels (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Priority definition
    priority_code VARCHAR(50) NOT NULL,
    priority_name VARCHAR(100) NOT NULL,
    priority_value INT NOT NULL, -- 1-5, higher = more urgent
    
    -- Priority properties
    color_hex VARCHAR(7),
    icon VARCHAR(50),
    description TEXT,
    
    -- SLA configuration
    sla_response_hours INT,
    sla_resolution_hours INT,
    
    -- Notification rules
    notify_on_assignment BOOLEAN DEFAULT TRUE,
    notify_supervisors BOOLEAN DEFAULT FALSE,
    escalation_enabled BOOLEAN DEFAULT FALSE,
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    
    UNIQUE KEY unique_tenant_priority (tenant_id, priority_code),
    INDEX idx_tenant_value (tenant_id, priority_value),
    INDEX idx_tenant_active (tenant_id, is_active)
);

-- ============================================================================
-- TENANT CUSTOMIZATION: FIELDS & FORMS
-- ============================================================================

-- Custom fields for tasks per tenant
CREATE TABLE IF NOT EXISTS tenant_task_fields (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Field definition
    field_code VARCHAR(50) NOT NULL,
    field_name VARCHAR(100) NOT NULL,
    field_type VARCHAR(50), -- text, textarea, number, select, date, checkbox, user, lead
    
    -- Field properties
    is_required BOOLEAN DEFAULT FALSE,
    is_visible BOOLEAN DEFAULT TRUE,
    is_editable BOOLEAN DEFAULT TRUE,
    display_order INT DEFAULT 0,
    
    -- Validation
    validation_rules JSON, -- Min/max length, regex patterns, etc.
    default_value VARCHAR(500),
    
    -- Options for select fields
    field_options JSON, -- Array of {label, value} pairs
    
    -- Visibility rules
    visible_on_statuses VARCHAR(500), -- JSON array of status codes
    visible_on_task_types VARCHAR(500), -- JSON array of task type codes
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    
    UNIQUE KEY unique_tenant_field (tenant_id, field_code),
    INDEX idx_tenant_active (tenant_id, is_active),
    INDEX idx_display_order (display_order)
);

-- ============================================================================
-- TENANT CUSTOMIZATION: RULES & AUTOMATION
-- ============================================================================

-- Automation rules per tenant (Auto-status changes, notifications, assignments)
CREATE TABLE IF NOT EXISTS tenant_automation_rules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Rule definition
    rule_code VARCHAR(50) NOT NULL,
    rule_name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Trigger conditions
    trigger_event VARCHAR(50), -- task_created, status_changed, due_soon, overdue, assigned, commented
    trigger_conditions JSON, -- Complex conditions {field, operator, value}
    
    -- Actions
    action_type VARCHAR(50), -- auto_status_change, auto_assign, send_notification, escalate
    action_data JSON, -- Action-specific data
    
    -- Rule properties
    is_active BOOLEAN DEFAULT TRUE,
    priority INT DEFAULT 10, -- Lower number = higher priority
    run_once BOOLEAN DEFAULT FALSE, -- Run only once per task
    
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES `user`(id) ON DELETE SET NULL,
    
    UNIQUE KEY unique_tenant_rule (tenant_id, rule_code),
    INDEX idx_tenant_active (tenant_id, is_active),
    INDEX idx_priority (priority)
);

-- ============================================================================
-- TENANT CUSTOMIZATION AUDIT & VALIDATION
-- ============================================================================

-- Track customization changes for audit
CREATE TABLE IF NOT EXISTS tenant_customization_audit (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Change details
    entity_type VARCHAR(50), -- status, stage, task_type, priority, field, rule
    entity_code VARCHAR(50),
    action VARCHAR(20), -- created, updated, deleted, activated, deactivated
    
    -- Change data
    old_values JSON,
    new_values JSON,
    
    -- Audit
    changed_by BIGINT,
    change_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES `tenant`(id) ON DELETE CASCADE,
    FOREIGN KEY (changed_by) REFERENCES `user`(id) ON DELETE SET NULL,
    
    INDEX idx_tenant_created (tenant_id, created_at DESC),
    INDEX idx_entity_type (entity_type)
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

CREATE INDEX idx_task_statuses_tenant_code ON tenant_task_statuses(tenant_id, status_code);
CREATE INDEX idx_task_stages_tenant_code ON tenant_task_stages(tenant_id, stage_code);
CREATE INDEX idx_task_types_tenant_code ON tenant_task_types(tenant_id, type_code);
CREATE INDEX idx_notif_types_tenant_code ON tenant_notification_types(tenant_id, type_code);
CREATE INDEX idx_priority_tenant_code ON tenant_priority_levels(tenant_id, priority_code);

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================
-- Tenant-level customization tables created
-- Enables full tenant customization of task workflows, statuses, and notifications
-- Status: Ready for configuration management service implementation

SELECT 'Tenant Customization Migration Complete' as status;
