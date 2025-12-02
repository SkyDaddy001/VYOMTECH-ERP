-- Civil Module Schema Migration
-- Creates tables for Sites, Safety Incidents, Compliance Records, and Permits

-- Sites Table
CREATE TABLE IF NOT EXISTS sites (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    site_name VARCHAR(255) NOT NULL,
    location VARCHAR(500),
    project_id VARCHAR(255),
    site_manager VARCHAR(255),
    start_date DATETIME,
    expected_end_date DATETIME,
    current_status VARCHAR(50), -- planning, active, paused, completed, closed
    site_area_sqm DECIMAL(12, 2),
    workforce_count INT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_project_id (project_id),
    INDEX idx_status (current_status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Safety Incidents Table
CREATE TABLE IF NOT EXISTS safety_incidents (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    site_id BIGINT NOT NULL,
    incident_type VARCHAR(50) NOT NULL, -- accident, near_miss, hazard, violation
    severity VARCHAR(50) NOT NULL, -- low, medium, high, critical
    incident_date DATETIME NOT NULL,
    description LONGTEXT,
    reported_by VARCHAR(255),
    status VARCHAR(50) NOT NULL, -- open, investigating, resolved, closed
    incident_number VARCHAR(100) NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_site_id (site_id),
    INDEX idx_severity (severity),
    INDEX idx_status (status),
    INDEX idx_incident_date (incident_date),
    FOREIGN KEY (site_id) REFERENCES sites(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Compliance Records Table
CREATE TABLE IF NOT EXISTS compliance_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    site_id BIGINT NOT NULL,
    compliance_type VARCHAR(100) NOT NULL, -- safety, environmental, labor, regulatory
    requirement VARCHAR(500),
    due_date DATETIME,
    status VARCHAR(50) NOT NULL, -- compliant, non_compliant, in_progress, not_applicable
    last_audit_date DATETIME,
    audit_result VARCHAR(50), -- pass, fail, pending
    notes LONGTEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_site_id (site_id),
    INDEX idx_compliance_type (compliance_type),
    INDEX idx_status (status),
    INDEX idx_due_date (due_date),
    FOREIGN KEY (site_id) REFERENCES sites(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Permits Table
CREATE TABLE IF NOT EXISTS permits (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    site_id BIGINT NOT NULL,
    permit_type VARCHAR(100),
    permit_number VARCHAR(100) NOT NULL UNIQUE,
    issued_date DATETIME,
    expiry_date DATETIME,
    issuing_authority VARCHAR(255),
    status VARCHAR(50) NOT NULL, -- active, expired, cancelled, pending
    document_url VARCHAR(500),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_site_id (site_id),
    INDEX idx_status (status),
    INDEX idx_expiry_date (expiry_date),
    FOREIGN KEY (site_id) REFERENCES sites(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
