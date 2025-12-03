-- ============================================================
-- MIGRATION 002: CIVIL ENGINEERING MODULE
-- Date: December 3, 2025
-- Purpose: Create civil engineering and construction site tables
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- SITES
-- ============================================================
CREATE TABLE IF NOT EXISTS `sites` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `site_name` VARCHAR(255) NOT NULL,
    `location` VARCHAR(255),
    `project_id` VARCHAR(36),
    `site_manager` VARCHAR(255),
    `start_date` TIMESTAMP,
    `expected_end_date` TIMESTAMP,
    `current_status` VARCHAR(50) DEFAULT 'planning',
    `site_area_sqm` DECIMAL(15, 2),
    `workforce_count` INT DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_status` (`current_status`),
    KEY `idx_project` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- SAFETY INCIDENTS
-- ============================================================
CREATE TABLE IF NOT EXISTS `safety_incidents` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `site_id` BIGINT,
    `incident_type` VARCHAR(100),
    `severity` VARCHAR(50),
    `incident_date` TIMESTAMP,
    `description` LONGTEXT,
    `reported_by` VARCHAR(255),
    `status` VARCHAR(50) DEFAULT 'open',
    `incident_number` VARCHAR(50) UNIQUE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`site_id`) REFERENCES `sites`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_site` (`site_id`),
    KEY `idx_severity` (`severity`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- COMPLIANCE RECORDS
-- ============================================================
CREATE TABLE IF NOT EXISTS `compliance_records` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `site_id` BIGINT,
    `compliance_type` VARCHAR(100),
    `requirement` VARCHAR(255),
    `due_date` TIMESTAMP,
    `status` VARCHAR(50) DEFAULT 'in_progress',
    `last_audit_date` TIMESTAMP,
    `audit_result` VARCHAR(50),
    `notes` LONGTEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`site_id`) REFERENCES `sites`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_site` (`site_id`),
    KEY `idx_type` (`compliance_type`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PERMITS
-- ============================================================
CREATE TABLE IF NOT EXISTS `permits` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `site_id` BIGINT,
    `permit_type` VARCHAR(100),
    `permit_number` VARCHAR(50) UNIQUE,
    `issued_date` TIMESTAMP,
    `expiry_date` TIMESTAMP,
    `issuing_authority` VARCHAR(255),
    `status` VARCHAR(50) DEFAULT 'pending',
    `document_url` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`site_id`) REFERENCES `sites`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_site` (`site_id`),
    KEY `idx_status` (`status`),
    KEY `idx_expiry` (`expiry_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
