-- Scheduled Tasks Module Schema Migration
-- Creates tables for Task Management, Scheduling, Reminders, and Notifications

-- Scheduled Tasks Table
CREATE TABLE IF NOT EXISTS scheduled_tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    assigned_to VARCHAR(255) NOT NULL,
    task_title VARCHAR(255) NOT NULL,
    description LONGTEXT,
    task_type VARCHAR(100), -- call, followup, meeting, reminder, other
    priority VARCHAR(50), -- low, medium, high, critical
    due_date DATETIME NOT NULL,
    scheduled_date DATETIME,
    reminder_before_minutes INT,
    status VARCHAR(50) NOT NULL, -- pending, in_progress, completed, cancelled, overdue
    related_entity_type VARCHAR(100), -- lead, contact, deal, account
    related_entity_id VARCHAR(255),
    notes LONGTEXT,
    created_by VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_assigned_to (assigned_to),
    INDEX idx_due_date (due_date),
    INDEX idx_scheduled_date (scheduled_date),
    INDEX idx_status (status),
    INDEX idx_priority (priority),
    INDEX idx_related_entity (related_entity_type, related_entity_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Task Reminders Table
CREATE TABLE IF NOT EXISTS task_reminders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    task_id BIGINT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    reminder_type VARCHAR(100), -- email, sms, push, in_app
    scheduled_time DATETIME NOT NULL,
    sent_time DATETIME,
    status VARCHAR(50), -- pending, sent, failed, skipped
    retry_count INT DEFAULT 0,
    error_message VARCHAR(500),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_task_id (task_id),
    INDEX idx_user_id (user_id),
    INDEX idx_scheduled_time (scheduled_time),
    INDEX idx_status (status),
    FOREIGN KEY (task_id) REFERENCES scheduled_tasks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Task Activity Log Table
CREATE TABLE IF NOT EXISTS task_activity_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    task_id BIGINT NOT NULL,
    activity_type VARCHAR(100), -- created, assigned, updated, completed, commented
    old_value VARCHAR(500),
    new_value VARCHAR(500),
    changed_by VARCHAR(255),
    notes LONGTEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_task_id (task_id),
    INDEX idx_activity_type (activity_type),
    INDEX idx_changed_by (changed_by),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (task_id) REFERENCES scheduled_tasks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Task Recurring Schedule Table
CREATE TABLE IF NOT EXISTS task_recurring_schedule (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    original_task_id BIGINT NOT NULL,
    recurrence_pattern VARCHAR(100), -- daily, weekly, biweekly, monthly, yearly
    frequency INT, -- every X days/weeks/months
    day_of_week VARCHAR(50), -- for weekly: mon,tue,wed,thu,fri,sat,sun
    day_of_month INT, -- for monthly: 1-31
    end_date DATETIME,
    max_occurrences INT,
    occurrences_created INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_original_task_id (original_task_id),
    INDEX idx_end_date (end_date),
    INDEX idx_is_active (is_active),
    FOREIGN KEY (original_task_id) REFERENCES scheduled_tasks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Task Dependencies Table
CREATE TABLE IF NOT EXISTS task_dependencies (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    task_id BIGINT NOT NULL,
    depends_on_task_id BIGINT NOT NULL,
    dependency_type VARCHAR(50), -- blocks, blocked_by, related_to
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_task_id (task_id),
    INDEX idx_depends_on_task_id (depends_on_task_id),
    UNIQUE KEY unique_dependency (task_id, depends_on_task_id),
    FOREIGN KEY (task_id) REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (depends_on_task_id) REFERENCES scheduled_tasks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Task Comments Table
CREATE TABLE IF NOT EXISTS task_comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    task_id BIGINT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    comment_text LONGTEXT NOT NULL,
    attachments_url VARCHAR(500),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_task_id (task_id),
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (task_id) REFERENCES scheduled_tasks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
