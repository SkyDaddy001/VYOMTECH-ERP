-- ============================================================
-- MIGRATION 014: GL POSTING & ACCOUNTING LINKS
-- Date: December 3, 2025
-- Purpose: Create GL posting templates and links for all modules
-- Ensures: Payroll, Purchase, Sales, Construction post to GL
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- GL POSTING TEMPLATE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `gl_posting_template` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `template_name` VARCHAR(255) NOT NULL,
    `template_type` VARCHAR(50),
    `module` VARCHAR(50),
    `description` TEXT,
    `is_active` BOOLEAN DEFAULT TRUE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_module` (`module`),
    UNIQUE KEY `unique_template` (`tenant_id`, `template_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- GL POSTING TEMPLATE LINE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `gl_posting_template_line` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `template_id` VARCHAR(36) NOT NULL,
    `line_number` INT,
    `account_id` VARCHAR(36) NOT NULL,
    `posting_type` VARCHAR(50),
    `amount_type` VARCHAR(50),
    `amount_field` VARCHAR(100),
    `is_mandatory` BOOLEAN DEFAULT FALSE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`template_id`) REFERENCES `gl_posting_template`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_template` (`template_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PAYROLL GL POSTING TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `payroll_gl_posting` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `payroll_record_id` VARCHAR(36) NOT NULL,
    `journal_entry_id` VARCHAR(36),
    `posting_date` DATE,
    `posting_status` VARCHAR(50) DEFAULT 'pending',
    `salary_expense_account_id` VARCHAR(36),
    `payable_account_id` VARCHAR(36),
    `bank_account_id` VARCHAR(36),
    `epf_payable_account_id` VARCHAR(36),
    `esi_payable_account_id` VARCHAR(36),
    `tax_payable_account_id` VARCHAR(36),
    `posted_at` TIMESTAMP,
    `posted_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`payroll_record_id`) REFERENCES `payroll_record`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    FOREIGN KEY (`salary_expense_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`payable_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`bank_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_payroll` (`payroll_record_id`),
    KEY `idx_status` (`posting_status`),
    UNIQUE KEY `unique_payroll_posting` (`payroll_record_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PURCHASE GL POSTING TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `purchase_gl_posting` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `purchase_order_id` VARCHAR(36) NOT NULL,
    `journal_entry_id` VARCHAR(36),
    `posting_date` DATE,
    `posting_status` VARCHAR(50) DEFAULT 'pending',
    `inventory_account_id` VARCHAR(36),
    `expense_account_id` VARCHAR(36),
    `payable_account_id` VARCHAR(36),
    `tax_payable_account_id` VARCHAR(36),
    `posted_at` TIMESTAMP,
    `posted_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`purchase_order_id`) REFERENCES `purchase_order`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    FOREIGN KEY (`inventory_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`expense_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`payable_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_po` (`purchase_order_id`),
    KEY `idx_status` (`posting_status`),
    UNIQUE KEY `unique_po_posting` (`purchase_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- SALES GL POSTING TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `sales_gl_posting` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `sales_invoice_id` VARCHAR(36) NOT NULL,
    `journal_entry_id` VARCHAR(36),
    `posting_date` DATE,
    `posting_status` VARCHAR(50) DEFAULT 'pending',
    `revenue_account_id` VARCHAR(36),
    `receivable_account_id` VARCHAR(36),
    `bank_account_id` VARCHAR(36),
    `tax_receivable_account_id` VARCHAR(36),
    `discount_account_id` VARCHAR(36),
    `posted_at` TIMESTAMP,
    `posted_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`sales_invoice_id`) REFERENCES `sales_invoice`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    FOREIGN KEY (`revenue_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`receivable_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`bank_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_invoice` (`sales_invoice_id`),
    KEY `idx_status` (`posting_status`),
    UNIQUE KEY `unique_invoice_posting` (`sales_invoice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- CONSTRUCTION GL POSTING TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `construction_gl_posting` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `boq_id` BIGINT NOT NULL,
    `journal_entry_id` VARCHAR(36),
    `posting_date` DATE,
    `posting_status` VARCHAR(50) DEFAULT 'pending',
    `wip_account_id` VARCHAR(36),
    `cost_account_id` VARCHAR(36),
    `payable_account_id` VARCHAR(36),
    `revenue_account_id` VARCHAR(36),
    `posted_at` TIMESTAMP,
    `posted_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`boq_id`) REFERENCES `bill_of_quantities`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    FOREIGN KEY (`wip_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`cost_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`payable_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`revenue_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_boq` (`boq_id`),
    KEY `idx_status` (`posting_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- REAL ESTATE GL POSTING TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `real_estate_gl_posting` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `booking_id` VARCHAR(36) NOT NULL,
    `journal_entry_id` VARCHAR(36),
    `posting_date` DATE,
    `posting_status` VARCHAR(50) DEFAULT 'pending',
    `asset_account_id` VARCHAR(36),
    `receivable_account_id` VARCHAR(36),
    `revenue_account_id` VARCHAR(36),
    `deferred_revenue_account_id` VARCHAR(36),
    `posted_at` TIMESTAMP,
    `posted_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`booking_id`) REFERENCES `property_booking`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    FOREIGN KEY (`asset_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`receivable_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`revenue_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_booking` (`booking_id`),
    KEY `idx_status` (`posting_status`),
    UNIQUE KEY `unique_booking_posting` (`booking_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- GL POSTING AUDIT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `gl_posting_audit` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `entity_type` VARCHAR(50),
    `entity_id` VARCHAR(36),
    `journal_entry_id` VARCHAR(36),
    `posting_amount` DECIMAL(18, 2),
    `action` VARCHAR(50),
    `reason` TEXT,
    `posted_by` VARCHAR(36),
    `posted_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_entity` (`entity_type`, `entity_id`),
    KEY `idx_timestamp` (`posted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ACCOUNT MAPPING TABLE (Module to GL Account mapping)
-- ============================================================
CREATE TABLE IF NOT EXISTS `account_mapping` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `module_name` VARCHAR(100) NOT NULL,
    `mapping_type` VARCHAR(100) NOT NULL,
    `account_id` VARCHAR(36) NOT NULL,
    `is_default` BOOLEAN DEFAULT FALSE,
    `is_active` BOOLEAN DEFAULT TRUE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_module` (`module_name`),
    KEY `idx_mapping_type` (`mapping_type`),
    UNIQUE KEY `unique_mapping` (`tenant_id`, `module_name`, `mapping_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
