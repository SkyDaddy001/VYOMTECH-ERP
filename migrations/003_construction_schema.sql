-- Construction Module Schema Migration
-- Creates tables for Construction Projects, BOQ, Progress Tracking, Quality Control, and Equipment

-- Construction Projects Table
CREATE TABLE IF NOT EXISTS construction_projects (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    project_name VARCHAR(255) NOT NULL,
    project_code VARCHAR(100) NOT NULL UNIQUE,
    location VARCHAR(500),
    client VARCHAR(255),
    contract_value DECIMAL(15, 2),
    start_date DATETIME,
    expected_completion DATETIME,
    current_progress_percent INT DEFAULT 0,
    status VARCHAR(50) NOT NULL, -- planning, active, suspended, completed, on_hold
    project_manager VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_start_date (start_date),
    INDEX idx_expected_completion (expected_completion),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Bill of Quantities Table
CREATE TABLE IF NOT EXISTS bill_of_quantities (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    project_id BIGINT NOT NULL,
    boq_number VARCHAR(100),
    item_description VARCHAR(500),
    unit VARCHAR(50),
    quantity DECIMAL(12, 4),
    unit_rate DECIMAL(12, 2),
    total_amount DECIMAL(15, 2),
    category VARCHAR(100), -- civil, structural, electrical, plumbing, finishing, other
    status VARCHAR(50) NOT NULL, -- planned, in_progress, completed, on_hold
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_project_id (project_id),
    INDEX idx_category (category),
    INDEX idx_status (status),
    FOREIGN KEY (project_id) REFERENCES construction_projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Progress Tracking Table
CREATE TABLE IF NOT EXISTS progress_tracking (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    project_id BIGINT NOT NULL,
    date DATETIME NOT NULL,
    activity_desc LONGTEXT,
    quantity_completed DECIMAL(12, 4),
    unit VARCHAR(50),
    percent_complete INT,
    workforce_deployed INT,
    notes LONGTEXT,
    photo_url VARCHAR(500),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_project_id (project_id),
    INDEX idx_date (date),
    INDEX idx_percent_complete (percent_complete),
    FOREIGN KEY (project_id) REFERENCES construction_projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Quality Control Table
CREATE TABLE IF NOT EXISTS quality_control (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    project_id BIGINT NOT NULL,
    boq_item_id BIGINT,
    inspection_date DATETIME NOT NULL,
    inspector_name VARCHAR(255),
    quality_status VARCHAR(50) NOT NULL, -- passed, failed, partial, pending
    observations LONGTEXT,
    corrective_action LONGTEXT,
    follow_up_date DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_project_id (project_id),
    INDEX idx_boq_item_id (boq_item_id),
    INDEX idx_quality_status (quality_status),
    INDEX idx_inspection_date (inspection_date),
    FOREIGN KEY (project_id) REFERENCES construction_projects(id) ON DELETE CASCADE,
    FOREIGN KEY (boq_item_id) REFERENCES bill_of_quantities(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Construction Equipment Table
CREATE TABLE IF NOT EXISTS construction_equipment (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    project_id BIGINT NOT NULL,
    equipment_name VARCHAR(255),
    equipment_type VARCHAR(100),
    serial_number VARCHAR(100),
    status VARCHAR(50) NOT NULL, -- available, in_use, maintenance, retired
    deployment_date DATETIME,
    retirement_date DATETIME,
    cost_per_day DECIMAL(10, 2),
    notes LONGTEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_project_id (project_id),
    INDEX idx_status (status),
    INDEX idx_equipment_type (equipment_type),
    FOREIGN KEY (project_id) REFERENCES construction_projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
