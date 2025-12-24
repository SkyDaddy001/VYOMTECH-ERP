-- ============================================================
-- MIGRATION 011: COMPLIANCE & TAX MODULE
-- Date: December 3, 2025
-- Purpose: Create compliance, tax, audit, and regulatory tables
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- COMPLIANCE RECORD TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `compliance_record` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `compliance_type` VARCHAR(100) NOT NULL,
    `title` VARCHAR(255),
    `description` TEXT,
    `status` VARCHAR(50) DEFAULT 'pending',
    `due_date` DATE,
    `completed_date` DATE,
    `assigned_to` VARCHAR(36),
    `document_url` TEXT,
    `notes` TEXT,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_type` (`compliance_type`),
    KEY `idx_status` (`status`),
    KEY `idx_due_date` (`due_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- COMPLIANCE CHECKLIST TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `compliance_checklist` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `compliance_record_id` VARCHAR(36) NOT NULL,
    `item_number` INT,
    `item_description` TEXT,
    `is_completed` BOOLEAN DEFAULT FALSE,
    `completion_date` TIMESTAMP,
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`compliance_record_id`) REFERENCES `compliance_record`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_record` (`compliance_record_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- TAX CALCULATION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `tax_calculation` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `tax_period` DATE NOT NULL,
    `tax_type` VARCHAR(50),
    `taxable_amount` DECIMAL(18, 2),
    `tax_rate` DECIMAL(5, 2),
    `tax_amount` DECIMAL(18, 2),
    `deductions` DECIMAL(18, 2) DEFAULT 0,
    `net_tax_payable` DECIMAL(18, 2),
    `payment_status` VARCHAR(50) DEFAULT 'pending',
    `payment_date` DATE,
    `reference_number` VARCHAR(100),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_period` (`tax_period`),
    KEY `idx_type` (`tax_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- AUDIT TRAIL TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `audit_trail` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `entity_type` VARCHAR(100),
    `entity_id` VARCHAR(36),
    `action_type` VARCHAR(50),
    `old_values` JSON,
    `new_values` JSON,
    `changed_by` VARCHAR(36),
    `change_reason` TEXT,
    `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_entity` (`entity_type`, `entity_id`),
    KEY `idx_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- REGULATORY REQUIREMENT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `regulatory_requirement` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `requirement_name` VARCHAR(255) NOT NULL,
    `requirement_type` VARCHAR(100),
    `jurisdiction` VARCHAR(100),
    `description` TEXT,
    `frequency` VARCHAR(50),
    `due_date_pattern` VARCHAR(100),
    `penalty_for_non_compliance` TEXT,
    `status` VARCHAR(50) DEFAULT 'active',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_type` (`requirement_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- DOCUMENT MANAGEMENT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `document` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `document_type` VARCHAR(50),
    `document_name` VARCHAR(255) NOT NULL,
    `file_url` TEXT,
    `file_size` INT,
    `mime_type` VARCHAR(100),
    `uploaded_by` VARCHAR(36),
    `version` INT DEFAULT 1,
    `is_latest` BOOLEAN DEFAULT TRUE,
    `expiry_date` DATE,
    `status` VARCHAR(50) DEFAULT 'active',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_type` (`document_type`),
    KEY `idx_expiry` (`expiry_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
