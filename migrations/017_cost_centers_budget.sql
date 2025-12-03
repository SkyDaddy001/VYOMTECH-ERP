-- ============================================================
-- MIGRATION 017: COST CENTERS & PROFIT CENTER ACCOUNTING
-- Date: December 3, 2025
-- Purpose: Cost center allocation, profit center wise P&L, budget tracking
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- COST CENTER MASTER TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `cost_center` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `cost_center_code` VARCHAR(50) NOT NULL,
    `cost_center_name` VARCHAR(255) NOT NULL,
    `cost_center_type` VARCHAR(50),
    `parent_cost_center_id` VARCHAR(36),
    `description` TEXT,
    `cost_center_manager_id` VARCHAR(36),
    `budget_amount` DECIMAL(18, 2),
    `is_active` BOOLEAN DEFAULT TRUE,
    `is_profit_center` BOOLEAN DEFAULT FALSE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`parent_cost_center_id`) REFERENCES `cost_center`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_code` (`cost_center_code`),
    KEY `idx_type` (`cost_center_type`),
    KEY `idx_is_active` (`is_active`),
    UNIQUE KEY `unique_cost_center` (`tenant_id`, `cost_center_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- COST ALLOCATION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `cost_allocation` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `cost_center_id` VARCHAR(36) NOT NULL,
    `allocation_date` DATE NOT NULL,
    `fiscal_period` VARCHAR(10),
    `allocated_amount` DECIMAL(18, 2),
    `allocation_basis` VARCHAR(50),
    `allocation_percentage` DECIMAL(5, 2),
    `expense_category` VARCHAR(100),
    `journal_entry_id` VARCHAR(36),
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`cost_center_id`) REFERENCES `cost_center`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_cost_center` (`cost_center_id`),
    KEY `idx_allocation_date` (`allocation_date`),
    KEY `idx_fiscal_period` (`fiscal_period`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- COST CENTER EXPENSE DISTRIBUTION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `cost_distribution` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `source_cost_center_id` VARCHAR(36) NOT NULL,
    `target_cost_center_id` VARCHAR(36) NOT NULL,
    `distribution_date` DATE NOT NULL,
    `amount` DECIMAL(18, 2),
    `distribution_basis` VARCHAR(50),
    `distribution_percentage` DECIMAL(5, 2),
    `journal_entry_id` VARCHAR(36),
    `fiscal_period` VARCHAR(10),
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`source_cost_center_id`) REFERENCES `cost_center`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`target_cost_center_id`) REFERENCES `cost_center`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_source` (`source_cost_center_id`),
    KEY `idx_target` (`target_cost_center_id`),
    KEY `idx_distribution_date` (`distribution_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- COST CENTER WISE PROFIT & LOSS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `cost_center_pl` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `cost_center_id` VARCHAR(36) NOT NULL,
    `fiscal_period` VARCHAR(10),
    `revenue_amount` DECIMAL(18, 2) DEFAULT 0,
    `cost_of_goods_sold` DECIMAL(18, 2) DEFAULT 0,
    `gross_profit` DECIMAL(18, 2) DEFAULT 0,
    `operating_expenses` DECIMAL(18, 2) DEFAULT 0,
    `operating_profit` DECIMAL(18, 2) DEFAULT 0,
    `other_income` DECIMAL(18, 2) DEFAULT 0,
    `other_expenses` DECIMAL(18, 2) DEFAULT 0,
    `net_profit` DECIMAL(18, 2) DEFAULT 0,
    `profit_margin_percentage` DECIMAL(5, 2),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`cost_center_id`) REFERENCES `cost_center`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_cost_center` (`cost_center_id`),
    KEY `idx_fiscal_period` (`fiscal_period`),
    UNIQUE KEY `unique_pl` (`cost_center_id`, `fiscal_period`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- BUDGET MASTER TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `budget` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `budget_name` VARCHAR(255) NOT NULL,
    `budget_code` VARCHAR(50),
    `fiscal_year` VARCHAR(10),
    `fiscal_period` VARCHAR(10),
    `cost_center_id` VARCHAR(36),
    `budget_type` VARCHAR(50),
    `start_date` DATE,
    `end_date` DATE,
    `total_budget_amount` DECIMAL(18, 2),
    `budget_status` VARCHAR(50) DEFAULT 'draft',
    `approved_by` VARCHAR(36),
    `approved_at` TIMESTAMP,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`cost_center_id`) REFERENCES `cost_center`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_fiscal_year` (`fiscal_year`),
    KEY `idx_cost_center` (`cost_center_id`),
    KEY `idx_budget_status` (`budget_status`),
    UNIQUE KEY `unique_budget` (`tenant_id`, `budget_code`, `fiscal_year`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- BUDGET LINE ITEMS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `budget_line` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `budget_id` VARCHAR(36) NOT NULL,
    `account_id` VARCHAR(36),
    `line_number` INT,
    `account_code` VARCHAR(50),
    `account_name` VARCHAR(255),
    `account_type` VARCHAR(50),
    `budgeted_amount` DECIMAL(18, 2),
    `actual_amount` DECIMAL(18, 2) DEFAULT 0,
    `variance` DECIMAL(18, 2),
    `variance_percentage` DECIMAL(5, 2),
    `remarks` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`budget_id`) REFERENCES `budget`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_budget` (`budget_id`),
    KEY `idx_account` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- BUDGET vs ACTUAL VARIANCE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `budget_variance` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `budget_id` VARCHAR(36) NOT NULL,
    `budget_line_id` VARCHAR(36),
    `variance_date` DATE,
    `account_id` VARCHAR(36),
    `budgeted_amount` DECIMAL(18, 2),
    `actual_amount` DECIMAL(18, 2),
    `variance_amount` DECIMAL(18, 2),
    `variance_percentage` DECIMAL(5, 2),
    `variance_type` VARCHAR(50),
    `variance_reason` TEXT,
    `approved_by` VARCHAR(36),
    `action_taken` TEXT,
    `fiscal_period` VARCHAR(10),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`budget_id`) REFERENCES `budget`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`budget_line_id`) REFERENCES `budget_line`(`id`),
    FOREIGN KEY (`account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_budget` (`budget_id`),
    KEY `idx_variance_type` (`variance_type`),
    KEY `idx_fiscal_period` (`fiscal_period`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
