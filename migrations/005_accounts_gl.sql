-- ============================================================
-- MIGRATION 005: ACCOUNTS & GENERAL LEDGER MODULE
-- Date: December 3, 2025
-- Purpose: Create accounting, GL, journal entries, and financial reporting tables
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- CHART OF ACCOUNTS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `chart_of_account` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `account_code` VARCHAR(50) NOT NULL,
    `account_name` VARCHAR(255) NOT NULL,
    `account_type` VARCHAR(50),
    `sub_account_type` VARCHAR(50),
    `parent_account_id` VARCHAR(36),
    `description` TEXT,
    `opening_balance` DECIMAL(18, 2) DEFAULT 0,
    `current_balance` DECIMAL(18, 2) DEFAULT 0,
    `is_active` BOOLEAN DEFAULT TRUE,
    `is_header` BOOLEAN DEFAULT FALSE,
    `is_default` BOOLEAN DEFAULT FALSE,
    `currency` VARCHAR(10) DEFAULT 'INR',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`parent_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_code` (`account_code`),
    KEY `idx_type` (`account_type`),
    UNIQUE KEY `unique_account` (`tenant_id`, `account_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- FINANCIAL PERIOD TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `financial_period` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `period_name` VARCHAR(100) NOT NULL,
    `period_type` VARCHAR(50),
    `start_date` DATE NOT NULL,
    `end_date` DATE NOT NULL,
    `is_closed` BOOLEAN DEFAULT FALSE,
    `closed_by` VARCHAR(36),
    `closed_at` TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_period_name` (`period_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- JOURNAL ENTRY TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `journal_entry` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `entry_date` DATE NOT NULL,
    `reference_number` VARCHAR(100),
    `reference_type` VARCHAR(50),
    `reference_id` VARCHAR(36),
    `description` TEXT NOT NULL,
    `amount` DECIMAL(18, 2),
    `narration` TEXT,
    `entry_status` VARCHAR(50) DEFAULT 'draft',
    `posted_by` VARCHAR(36),
    `posted_at` TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_date` (`entry_date`),
    KEY `idx_status` (`entry_status`),
    KEY `idx_reference` (`reference_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- JOURNAL ENTRY DETAIL TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `journal_entry_detail` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `journal_entry_id` VARCHAR(36) NOT NULL,
    `account_id` VARCHAR(36) NOT NULL,
    `account_code` VARCHAR(50),
    `debit_amount` DECIMAL(18, 2) DEFAULT 0,
    `credit_amount` DECIMAL(18, 2) DEFAULT 0,
    `description` TEXT,
    `line_number` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_journal_entry` (`journal_entry_id`),
    KEY `idx_account` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- GL ACCOUNT BALANCE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `gl_account_balance` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `account_id` VARCHAR(36) NOT NULL,
    `fiscal_period` DATE NOT NULL,
    `opening_balance` DECIMAL(18, 2) DEFAULT 0,
    `total_debit` DECIMAL(18, 2) DEFAULT 0,
    `total_credit` DECIMAL(18, 2) DEFAULT 0,
    `closing_balance` DECIMAL(18, 2) DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_account` (`account_id`),
    KEY `idx_period` (`fiscal_period`),
    UNIQUE KEY `unique_balance` (`account_id`, `fiscal_period`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- TRIAL BALANCE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `trial_balance` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `period_id` VARCHAR(36) NOT NULL,
    `account_id` VARCHAR(36) NOT NULL,
    `account_code` VARCHAR(50),
    `account_name` VARCHAR(255),
    `debit_balance` DECIMAL(18, 2) DEFAULT 0,
    `credit_balance` DECIMAL(18, 2) DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`period_id`) REFERENCES `financial_period`(`id`),
    FOREIGN KEY (`account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_period` (`period_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- INCOME STATEMENT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `income_statement` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `period_id` VARCHAR(36) NOT NULL,
    `revenue_total` DECIMAL(18, 2) DEFAULT 0,
    `cost_of_goods_sold` DECIMAL(18, 2) DEFAULT 0,
    `gross_profit` DECIMAL(18, 2) DEFAULT 0,
    `operating_expenses` DECIMAL(18, 2) DEFAULT 0,
    `operating_income` DECIMAL(18, 2) DEFAULT 0,
    `other_income` DECIMAL(18, 2) DEFAULT 0,
    `other_expenses` DECIMAL(18, 2) DEFAULT 0,
    `income_before_tax` DECIMAL(18, 2) DEFAULT 0,
    `tax_expense` DECIMAL(18, 2) DEFAULT 0,
    `net_income` DECIMAL(18, 2) DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`period_id`) REFERENCES `financial_period`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_period` (`period_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- BALANCE SHEET TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `balance_sheet` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `period_id` VARCHAR(36) NOT NULL,
    `current_assets` DECIMAL(18, 2) DEFAULT 0,
    `fixed_assets` DECIMAL(18, 2) DEFAULT 0,
    `other_assets` DECIMAL(18, 2) DEFAULT 0,
    `total_assets` DECIMAL(18, 2) DEFAULT 0,
    `current_liabilities` DECIMAL(18, 2) DEFAULT 0,
    `long_term_liabilities` DECIMAL(18, 2) DEFAULT 0,
    `other_liabilities` DECIMAL(18, 2) DEFAULT 0,
    `total_liabilities` DECIMAL(18, 2) DEFAULT 0,
    `equity` DECIMAL(18, 2) DEFAULT 0,
    `total_equity_and_liabilities` DECIMAL(18, 2) DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`period_id`) REFERENCES `financial_period`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_period` (`period_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
