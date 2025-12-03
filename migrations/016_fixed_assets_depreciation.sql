-- ============================================================
-- MIGRATION 016: FIXED ASSETS & DEPRECIATION MANAGEMENT
-- Date: December 3, 2025
-- Purpose: Asset register, depreciation schedules, asset disposal
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- FIXED ASSET MASTER TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `fixed_asset` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `asset_code` VARCHAR(50) NOT NULL,
    `asset_name` VARCHAR(255) NOT NULL,
    `asset_category` VARCHAR(100),
    `asset_description` TEXT,
    `asset_location` VARCHAR(255),
    `purchase_date` DATE NOT NULL,
    `original_cost` DECIMAL(18, 2) NOT NULL,
    `salvage_value` DECIMAL(18, 2) DEFAULT 0,
    `useful_life_years` INT,
    `depreciation_method` VARCHAR(50),
    `depreciation_rate` DECIMAL(5, 2),
    `asset_status` VARCHAR(50) DEFAULT 'active',
    `gl_asset_account_id` VARCHAR(36),
    `accumulated_depreciation_account_id` VARCHAR(36),
    `depreciation_expense_account_id` VARCHAR(36),
    `vendor_id` VARCHAR(36),
    `invoice_number` VARCHAR(100),
    `serial_number` VARCHAR(100),
    `warranty_expiry_date` DATE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`gl_asset_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`accumulated_depreciation_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`depreciation_expense_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`vendor_id`) REFERENCES `vendor`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_code` (`asset_code`),
    KEY `idx_category` (`asset_category`),
    KEY `idx_status` (`asset_status`),
    KEY `idx_purchase_date` (`purchase_date`),
    UNIQUE KEY `unique_asset_code` (`tenant_id`, `asset_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- DEPRECIATION SCHEDULE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `depreciation_schedule` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `fixed_asset_id` VARCHAR(36) NOT NULL,
    `fiscal_year` VARCHAR(10),
    `opening_cost` DECIMAL(18, 2),
    `depreciation_rate` DECIMAL(5, 2),
    `depreciation_amount` DECIMAL(18, 2),
    `accumulated_depreciation` DECIMAL(18, 2),
    `net_book_value` DECIMAL(18, 2),
    `schedule_date` DATE,
    `is_posted` BOOLEAN DEFAULT FALSE,
    `journal_entry_id` VARCHAR(36),
    `posted_at` TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`fixed_asset_id`) REFERENCES `fixed_asset`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_asset` (`fixed_asset_id`),
    KEY `idx_fiscal_year` (`fiscal_year`),
    KEY `idx_is_posted` (`is_posted`),
    UNIQUE KEY `unique_schedule` (`fixed_asset_id`, `fiscal_year`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ASSET REVALUATION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `asset_revaluation` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `fixed_asset_id` VARCHAR(36) NOT NULL,
    `revaluation_date` DATE NOT NULL,
    `previous_cost` DECIMAL(18, 2),
    `new_cost` DECIMAL(18, 2),
    `revaluation_surplus` DECIMAL(18, 2),
    `reason` TEXT,
    `approved_by` VARCHAR(36),
    `journal_entry_id` VARCHAR(36),
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`fixed_asset_id`) REFERENCES `fixed_asset`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_asset` (`fixed_asset_id`),
    KEY `idx_revaluation_date` (`revaluation_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ASSET DISPOSAL TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `asset_disposal` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `fixed_asset_id` VARCHAR(36) NOT NULL,
    `disposal_date` DATE NOT NULL,
    `disposal_type` VARCHAR(50),
    `selling_price` DECIMAL(18, 2),
    `book_value` DECIMAL(18, 2),
    `gain_loss` DECIMAL(18, 2),
    `disposal_method` VARCHAR(50),
    `buyer_name` VARCHAR(255),
    `disposal_reference` VARCHAR(100),
    `remarks` TEXT,
    `disposal_gl_posting_id` VARCHAR(36),
    `journal_entry_id` VARCHAR(36),
    `approved_by` VARCHAR(36),
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`fixed_asset_id`) REFERENCES `fixed_asset`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_asset` (`fixed_asset_id`),
    KEY `idx_disposal_date` (`disposal_date`),
    KEY `idx_disposal_type` (`disposal_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ASSET MAINTENANCE LOG TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `asset_maintenance` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `fixed_asset_id` VARCHAR(36) NOT NULL,
    `maintenance_date` DATE NOT NULL,
    `maintenance_type` VARCHAR(50),
    `maintenance_description` TEXT,
    `maintenance_cost` DECIMAL(18, 2),
    `vendor_id` VARCHAR(36),
    `invoice_number` VARCHAR(100),
    `expense_account_id` VARCHAR(36),
    `journal_entry_id` VARCHAR(36),
    `next_maintenance_date` DATE,
    `notes` TEXT,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`fixed_asset_id`) REFERENCES `fixed_asset`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`vendor_id`) REFERENCES `vendor`(`id`),
    FOREIGN KEY (`expense_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_asset` (`fixed_asset_id`),
    KEY `idx_maintenance_date` (`maintenance_date`),
    KEY `idx_maintenance_type` (`maintenance_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ASSET TRANSFER TABLE (For asset movement between locations/departments)
-- ============================================================
CREATE TABLE IF NOT EXISTS `asset_transfer` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `fixed_asset_id` VARCHAR(36) NOT NULL,
    `transfer_date` DATE NOT NULL,
    `from_location` VARCHAR(255),
    `to_location` VARCHAR(255),
    `from_department` VARCHAR(255),
    `to_department` VARCHAR(255),
    `transfer_reason` TEXT,
    `approved_by` VARCHAR(36),
    `transferred_by` VARCHAR(36),
    `transfer_reference` VARCHAR(100),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`fixed_asset_id`) REFERENCES `fixed_asset`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_asset` (`fixed_asset_id`),
    KEY `idx_transfer_date` (`transfer_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
