-- ==================== PHASE 3B: WORKFLOW AUTOMATION MIGRATION ====================
-- This migration creates the workflow management and automation infrastructure
-- Includes workflow definitions, triggers, actions, executions, and scheduling

-- ==================== WORKFLOWS TABLE ====================
CREATE TABLE IF NOT EXISTS workflows (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    enabled BOOLEAN DEFAULT true,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tenant (tenant_id),
    INDEX idx_enabled (enabled),
    INDEX idx_created (created_at),
    UNIQUE KEY unique_tenant_name (tenant_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==================== WORKFLOW TRIGGERS TABLE ====================
-- Stores trigger configurations that activate workflows
CREATE TABLE IF NOT EXISTS workflow_triggers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    workflow_id BIGINT NOT NULL,
    trigger_type VARCHAR(100) NOT NULL,  -- lead_created, lead_scored, task_completed, etc
    trigger_config JSON,  -- Condition details
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_workflow (workflow_id),
    INDEX idx_trigger_type (trigger_type),
    FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==================== TRIGGER CONDITIONS TABLE ====================
-- Stores detailed conditions for trigger evaluation
CREATE TABLE IF NOT EXISTS trigger_conditions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    trigger_id BIGINT NOT NULL,
    field VARCHAR(255) NOT NULL,  -- lead_score, lead_status, etc
    operator VARCHAR(50) NOT NULL,  -- equals, greater_than, less_than, contains
    value VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_trigger (trigger_id),
    FOREIGN KEY (trigger_id) REFERENCES workflow_triggers(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==================== WORKFLOW ACTIONS TABLE ====================
-- Stores actions to be executed when workflow triggers
CREATE TABLE IF NOT EXISTS workflow_actions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    workflow_id BIGINT NOT NULL,
    action_type VARCHAR(100) NOT NULL,  -- send_email, send_sms, create_task, update_lead
    action_config JSON,  -- Action parameters
    action_order INT DEFAULT 0,  -- Execution order
    delay_seconds INT DEFAULT 0,  -- Delay before execution
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_workflow (workflow_id),
    INDEX idx_action_type (action_type),
    INDEX idx_order (action_order),
    FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==================== WORKFLOW INSTANCES TABLE ====================
-- Tracks execution instances of workflows
CREATE TABLE IF NOT EXISTS workflow_instances (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    workflow_id BIGINT NOT NULL,
    triggered_by VARCHAR(100),  -- lead_id, task_id, etc
    triggered_by_value VARCHAR(255),  -- The actual ID value
    status VARCHAR(50),  -- pending, running, completed, failed, cancelled
    progress INT DEFAULT 0,  -- 0-100
    executed_actions INT DEFAULT 0,
    failed_actions INT DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tenant (tenant_id),
    INDEX idx_workflow (workflow_id),
    INDEX idx_status (status),
    INDEX idx_created (created_at),
    FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==================== WORKFLOW ACTION EXECUTIONS TABLE ====================
-- Tracks execution of individual actions within a workflow
CREATE TABLE IF NOT EXISTS workflow_action_executions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    workflow_id BIGINT NOT NULL,
    instance_id BIGINT NOT NULL,
    action_id BIGINT NOT NULL,
    status VARCHAR(50),  -- pending, executing, completed, failed
    result JSON,  -- Execution result
    error_message TEXT,
    retry_count INT DEFAULT 0,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_workflow (workflow_id),
    INDEX idx_instance (instance_id),
    INDEX idx_action (action_id),
    INDEX idx_status (status),
    FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE,
    FOREIGN KEY (instance_id) REFERENCES workflow_instances(id) ON DELETE CASCADE,
    FOREIGN KEY (action_id) REFERENCES workflow_actions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==================== SCHEDULED TASKS TABLE ====================
-- Stores scheduled workflow and action executions
CREATE TABLE IF NOT EXISTS scheduled_tasks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100),  -- workflow, action, report, cleanup
    config JSON,  -- Task configuration
    schedule VARCHAR(255),  -- Cron expression
    last_run_at TIMESTAMP NULL,
    next_run_at TIMESTAMP,
    enabled BOOLEAN DEFAULT true,
    max_retries INT DEFAULT 3,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tenant (tenant_id),
    INDEX idx_type (type),
    INDEX idx_enabled (enabled),
    INDEX idx_next_run (next_run_at),
    UNIQUE KEY unique_tenant_task (tenant_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==================== SCHEDULED TASK EXECUTIONS TABLE ====================
-- Tracks execution history of scheduled tasks
CREATE TABLE IF NOT EXISTS scheduled_task_executions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    task_id BIGINT NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    status VARCHAR(50),  -- success, failed, running, cancelled
    output JSON,  -- Execution output
    error_message TEXT,
    duration INT,  -- milliseconds
    started_at TIMESTAMP,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_task (task_id),
    INDEX idx_tenant (tenant_id),
    INDEX idx_status (status),
    INDEX idx_started (started_at),
    FOREIGN KEY (task_id) REFERENCES scheduled_tasks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==================== WORKFLOW TEMPLATES TABLE ====================
-- Predefined workflow templates for quick setup
CREATE TABLE IF NOT EXISTS workflow_templates (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),  -- sales, support, onboarding, etc
    description TEXT,
    definition JSON,  -- Complete workflow definition
    is_public BOOLEAN DEFAULT false,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_category (category),
    INDEX idx_public (is_public),
    UNIQUE KEY unique_template_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==================== WORKFLOW EXECUTION AUDIT LOG TABLE ====================
-- Comprehensive audit trail for all workflow executions
CREATE TABLE IF NOT EXISTS workflow_execution_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    workflow_id BIGINT NOT NULL,
    instance_id BIGINT NOT NULL,
    action VARCHAR(255),  -- create, execute, complete, fail
    details JSON,  -- Additional details
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_tenant (tenant_id),
    INDEX idx_workflow (workflow_id),
    INDEX idx_instance (instance_id),
    INDEX idx_action (action),
    INDEX idx_created (created_at),
    FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE,
    FOREIGN KEY (instance_id) REFERENCES workflow_instances(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==================== INDEXES FOR PERFORMANCE ====================
-- Additional indexes for common queries

-- Workflow query optimization
CREATE INDEX idx_workflow_tenant_enabled ON workflows(tenant_id, enabled);

-- Instance query optimization
CREATE INDEX idx_instance_workflow_status ON workflow_instances(workflow_id, status);
CREATE INDEX idx_instance_tenant_status ON workflow_instances(tenant_id, status);

-- Action execution query optimization
CREATE INDEX idx_action_exec_instance_status ON workflow_action_executions(instance_id, status);

-- Scheduled task query optimization
CREATE INDEX idx_scheduled_tenant_enabled ON scheduled_tasks(tenant_id, enabled);
CREATE INDEX idx_scheduled_next_run_enabled ON scheduled_tasks(next_run_at, enabled);

-- Sample seed data for workflow templates

-- INSERT INTO workflow_templates (name, category, description, is_public, created_at) VALUES
-- ('Lead Scoring Workflow', 'sales', 'Automatically score leads based on engagement', true, NOW()),
-- ('Task Assignment Workflow', 'operations', 'Assign tasks to agents based on workload', true, NOW()),
-- ('Lead Nurture Workflow', 'sales', 'Send automated follow-ups to cold leads', true, NOW()),
-- ('Escalation Workflow', 'support', 'Escalate high-priority issues automatically', true, NOW());
