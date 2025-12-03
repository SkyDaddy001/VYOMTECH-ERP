-- ============================================================
-- MIGRATION 013: HR COMPLIANCE - ESI & PF (PROVIDENT FUND)
-- Date: December 3, 2025
-- Purpose: Create detailed ESI and PF tracking tables with compliance
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- EPF (EMPLOYEE PROVIDENT FUND) CONFIGURATION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `epf_configuration` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `organization_name` VARCHAR(255),
    `epf_number` VARCHAR(50) UNIQUE,
    `employer_contribution_rate` DECIMAL(5, 2) DEFAULT 12.00,
    `employee_contribution_rate` DECIMAL(5, 2) DEFAULT 12.00,
    `pension_contribution_rate` DECIMAL(5, 2) DEFAULT 8.33,
    `min_salary_limit` DECIMAL(18, 2) DEFAULT 0,
    `max_salary_limit` DECIMAL(18, 2),
    `effective_from` DATE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    UNIQUE KEY `unique_epf_config` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ESI (EMPLOYEE STATE INSURANCE) CONFIGURATION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `esi_configuration` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `esi_number` VARCHAR(50) UNIQUE,
    `employer_contribution_rate` DECIMAL(5, 2) DEFAULT 3.25,
    `employee_contribution_rate` DECIMAL(5, 2) DEFAULT 0.75,
    `wage_ceiling` DECIMAL(18, 2) DEFAULT 21000,
    `min_salary_limit` DECIMAL(18, 2) DEFAULT 0,
    `effective_from` DATE,
    `registration_date` DATE,
    `inspector_name` VARCHAR(100),
    `inspector_contact` VARCHAR(20),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    UNIQUE KEY `unique_esi_config` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- EMPLOYEE EPF REGISTRATION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `employee_epf_registration` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` VARCHAR(36) NOT NULL,
    `uan_number` VARCHAR(50) UNIQUE NOT NULL,
    `member_id` VARCHAR(50) UNIQUE,
    `aadhar_number` VARCHAR(50),
    `date_of_joining` DATE NOT NULL,
    `previous_employer_uan` VARCHAR(50),
    `previous_epf_balance` DECIMAL(18, 2) DEFAULT 0,
    `exemption_status` BOOLEAN DEFAULT FALSE,
    `exemption_reason` TEXT,
    `status` VARCHAR(50) DEFAULT 'active',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_employee` (`employee_id`),
    KEY `idx_uan` (`uan_number`),
    UNIQUE KEY `unique_epf_reg` (`employee_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- EMPLOYEE ESI REGISTRATION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `employee_esi_registration` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` VARCHAR(36) NOT NULL,
    `esi_number` VARCHAR(50) UNIQUE NOT NULL,
    `aadhar_number` VARCHAR(50),
    `date_of_joining` DATE NOT NULL,
    `status` VARCHAR(50) DEFAULT 'active',
    `coverage_status` VARCHAR(50) DEFAULT 'covered',
    `exemption_status` BOOLEAN DEFAULT FALSE,
    `exemption_from_date` DATE,
    `exemption_to_date` DATE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_employee` (`employee_id`),
    KEY `idx_esi_number` (`esi_number`),
    UNIQUE KEY `unique_esi_reg` (`employee_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- EPF CONTRIBUTION RECORD TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `epf_contribution` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` VARCHAR(36) NOT NULL,
    `payroll_month` DATE NOT NULL,
    `wages_for_contribution` DECIMAL(18, 2),
    `employer_contribution` DECIMAL(18, 2),
    `employee_contribution` DECIMAL(18, 2),
    `pension_fund_contribution` DECIMAL(18, 2),
    `total_contribution` DECIMAL(18, 2),
    `contribution_status` VARCHAR(50) DEFAULT 'pending',
    `payment_date` DATE,
    `challan_number` VARCHAR(50),
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_employee` (`employee_id`),
    KEY `idx_month` (`payroll_month`),
    KEY `idx_status` (`contribution_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ESI CONTRIBUTION RECORD TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `esi_contribution` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` VARCHAR(36) NOT NULL,
    `payroll_month` DATE NOT NULL,
    `wages_for_contribution` DECIMAL(18, 2),
    `employer_contribution` DECIMAL(18, 2),
    `employee_contribution` DECIMAL(18, 2),
    `total_contribution` DECIMAL(18, 2),
    `contribution_status` VARCHAR(50) DEFAULT 'pending',
    `payment_date` DATE,
    `challan_number` VARCHAR(50),
    `form_5_submitted` BOOLEAN DEFAULT FALSE,
    `form_5_submission_date` DATE,
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_employee` (`employee_id`),
    KEY `idx_month` (`payroll_month`),
    KEY `idx_status` (`contribution_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- EPF MEMBER PASSBOOK TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `epf_passbook` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` VARCHAR(36) NOT NULL,
    `opening_balance` DECIMAL(18, 2) DEFAULT 0,
    `employer_contribution_total` DECIMAL(18, 2) DEFAULT 0,
    `employee_contribution_total` DECIMAL(18, 2) DEFAULT 0,
    `pension_fund_total` DECIMAL(18, 2) DEFAULT 0,
    `interest_credited` DECIMAL(18, 2) DEFAULT 0,
    `closing_balance` DECIMAL(18, 2) DEFAULT 0,
    `last_updated` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_employee` (`employee_id`),
    UNIQUE KEY `unique_passbook` (`employee_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ESI CLAIM TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `esi_claim` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` VARCHAR(36) NOT NULL,
    `claim_number` VARCHAR(50) UNIQUE NOT NULL,
    `claim_date` DATE NOT NULL,
    `claim_type` VARCHAR(50),
    `claim_amount` DECIMAL(18, 2),
    `claim_status` VARCHAR(50) DEFAULT 'pending',
    `approved_amount` DECIMAL(18, 2),
    `approval_date` DATE,
    `claim_reason` TEXT,
    `supporting_documents` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_employee` (`employee_id`),
    KEY `idx_status` (`claim_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- STATUTORY COMPLIANCE RECORD TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `statutory_compliance_record` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `compliance_type` VARCHAR(50),
    `compliance_month` DATE NOT NULL,
    `due_date` DATE,
    `submission_date` DATE,
    `reference_number` VARCHAR(100),
    `status` VARCHAR(50) DEFAULT 'pending',
    `documents_submitted` TEXT,
    `officer_name` VARCHAR(100),
    `officer_contact` VARCHAR(20),
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_type` (`compliance_type`),
    KEY `idx_month` (`compliance_month`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
