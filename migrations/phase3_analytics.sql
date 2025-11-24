-- =====================================================
-- Phase 3 Analytics Tables Migration
-- =====================================================
-- Date: November 24, 2025
-- Status: Ready for execution
-- Tables: 8 new analytics tables
-- =====================================================

-- 1. Analytics Events Table
CREATE TABLE IF NOT EXISTS analytics_events (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    user_id BIGINT NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_data LONGTEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_event (tenant_id, event_type),
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    CONSTRAINT fk_ae_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. Conversion Funnels Table
CREATE TABLE IF NOT EXISTS conversion_funnels (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    lead_id BIGINT NOT NULL,
    campaign_id BIGINT NOT NULL,
    stage VARCHAR(100) NOT NULL,
    conversion_rate DECIMAL(5, 2),
    time_in_stage BIGINT DEFAULT 0,
    entered_at TIMESTAMP,
    exited_at TIMESTAMP,
    converted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_campaign (tenant_id, campaign_id),
    INDEX idx_lead_id (lead_id),
    INDEX idx_stage (stage),
    INDEX idx_entered_at (entered_at),
    CONSTRAINT fk_cf_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_cf_lead FOREIGN KEY (lead_id) REFERENCES leads(id) ON DELETE CASCADE,
    CONSTRAINT fk_cf_campaign FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. Agent Metrics Table
CREATE TABLE IF NOT EXISTS agent_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    agent_id BIGINT NOT NULL,
    metric_date DATE NOT NULL,
    calls_handled BIGINT DEFAULT 0,
    average_call_time BIGINT DEFAULT 0,
    leads_converted BIGINT DEFAULT 0,
    conversion_rate DECIMAL(5, 2) DEFAULT 0,
    customer_rating DECIMAL(3, 1),
    tasks_completed BIGINT DEFAULT 0,
    available_time BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_tenant_agent_date (tenant_id, agent_id, metric_date),
    INDEX idx_tenant_agent (tenant_id, agent_id),
    INDEX idx_metric_date (metric_date),
    CONSTRAINT fk_agm_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_agm_agent FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. Campaign Metrics Table
CREATE TABLE IF NOT EXISTS campaign_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    campaign_id BIGINT NOT NULL,
    metric_date DATE NOT NULL,
    leads_generated BIGINT DEFAULT 0,
    leads_contacted BIGINT DEFAULT 0,
    leads_converted BIGINT DEFAULT 0,
    conversion_rate DECIMAL(5, 2) DEFAULT 0,
    average_lead_value DECIMAL(10, 2) DEFAULT 0,
    total_revenue DECIMAL(12, 2) DEFAULT 0,
    roi DECIMAL(6, 2) DEFAULT 0,
    cost_per_lead DECIMAL(10, 2) DEFAULT 0,
    cost DECIMAL(12, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_tenant_campaign_date (tenant_id, campaign_id, metric_date),
    INDEX idx_tenant_campaign (tenant_id, campaign_id),
    INDEX idx_metric_date (metric_date),
    CONSTRAINT fk_cm_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_cm_campaign FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. Daily Reports Table
CREATE TABLE IF NOT EXISTS daily_reports (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    report_date DATE NOT NULL,
    total_calls BIGINT DEFAULT 0,
    total_leads BIGINT DEFAULT 0,
    converted_leads BIGINT DEFAULT 0,
    conversion_rate DECIMAL(5, 2) DEFAULT 0,
    average_call_time BIGINT DEFAULT 0,
    tasks_completed BIGINT DEFAULT 0,
    active_agents BIGINT DEFAULT 0,
    total_revenue DECIMAL(12, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_tenant_date (tenant_id, report_date),
    INDEX idx_tenant (tenant_id),
    INDEX idx_report_date (report_date),
    CONSTRAINT fk_dr_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 6. Custom Reports Table
CREATE TABLE IF NOT EXISTS custom_reports (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    created_by BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description LONGTEXT,
    metrics LONGTEXT,
    filters LONGTEXT,
    schedule VARCHAR(100),
    format VARCHAR(50),
    recipients LONGTEXT,
    enabled BOOLEAN DEFAULT TRUE,
    last_run_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant (tenant_id),
    INDEX idx_created_by (created_by),
    INDEX idx_enabled (enabled),
    CONSTRAINT fk_cr_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_cr_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 7. Report Executions Table
CREATE TABLE IF NOT EXISTS report_executions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    report_id BIGINT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    result_url VARCHAR(255),
    error_message LONGTEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant (tenant_id),
    INDEX idx_report_id (report_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    CONSTRAINT fk_re_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_re_report FOREIGN KEY (report_id) REFERENCES custom_reports(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 8. Dashboard Widgets Table
CREATE TABLE IF NOT EXISTS dashboard_widgets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    widget_type VARCHAR(100) NOT NULL,
    config LONGTEXT,
    position INT DEFAULT 0,
    size VARCHAR(50),
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_enabled (enabled),
    CONSTRAINT fk_dw_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_dw_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- Migration Complete
-- =====================================================
-- Tables Created: 8
-- Total Fields: 120+
-- Indexes: 30+
-- Foreign Keys: 16
-- =====================================================
