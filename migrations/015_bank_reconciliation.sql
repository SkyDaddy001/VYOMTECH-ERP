-- ============================================================
-- MIGRATION 015: BANK RECONCILIATION & CASH MANAGEMENT
-- Date: December 3, 2025
-- Purpose: Bank statement matching, reconciliation status, cash flow
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- BANK STATEMENT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `bank_statement` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `bank_account_id` VARCHAR(36) NOT NULL,
    `statement_date` DATE NOT NULL,
    `statement_period_start` DATE NOT NULL,
    `statement_period_end` DATE NOT NULL,
    `opening_balance` DECIMAL(18, 2),
    `closing_balance` DECIMAL(18, 2),
    `total_deposits` DECIMAL(18, 2) DEFAULT 0,
    `total_withdrawals` DECIMAL(18, 2) DEFAULT 0,
    `statement_reference` VARCHAR(100),
    `currency` VARCHAR(10) DEFAULT 'INR',
    `reconciliation_status` VARCHAR(50) DEFAULT 'pending',
    `reconciled_by` VARCHAR(36),
    `reconciled_at` TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`bank_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_bank_account` (`bank_account_id`),
    KEY `idx_statement_date` (`statement_date`),
    KEY `idx_reconciliation_status` (`reconciliation_status`),
    UNIQUE KEY `unique_statement` (`bank_account_id`, `statement_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- BANK TRANSACTION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `bank_transaction` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `bank_statement_id` VARCHAR(36) NOT NULL,
    `transaction_date` DATE NOT NULL,
    `cheque_number` VARCHAR(50),
    `utr_number` VARCHAR(100),
    `description` TEXT,
    `debit_amount` DECIMAL(18, 2) DEFAULT 0,
    `credit_amount` DECIMAL(18, 2) DEFAULT 0,
    `balance_after_transaction` DECIMAL(18, 2),
    `transaction_type` VARCHAR(50),
    `remarks` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`bank_statement_id`) REFERENCES `bank_statement`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_statement` (`bank_statement_id`),
    KEY `idx_date` (`transaction_date`),
    KEY `idx_cheque` (`cheque_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- BANK RECONCILIATION MATCHING TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `bank_reconciliation_match` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `bank_statement_id` VARCHAR(36) NOT NULL,
    `bank_transaction_id` VARCHAR(36) NOT NULL,
    `journal_entry_detail_id` VARCHAR(36),
    `cheque_id` VARCHAR(36),
    `matched_amount` DECIMAL(18, 2),
    `match_date` DATE,
    `match_status` VARCHAR(50) DEFAULT 'matched',
    `variance_amount` DECIMAL(18, 2) DEFAULT 0,
    `remarks` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`bank_statement_id`) REFERENCES `bank_statement`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`bank_transaction_id`) REFERENCES `bank_transaction`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_detail_id`) REFERENCES `journal_entry_detail`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_statement` (`bank_statement_id`),
    KEY `idx_match_status` (`match_status`),
    UNIQUE KEY `unique_match` (`bank_transaction_id`, `journal_entry_detail_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- UNCLEARED ITEMS TABLE (Outstanding cheques, pending deposits)
-- ============================================================
CREATE TABLE IF NOT EXISTS `uncleared_item` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `bank_account_id` VARCHAR(36) NOT NULL,
    `item_type` VARCHAR(50),
    `cheque_number` VARCHAR(50),
    `utr_number` VARCHAR(100),
    `amount` DECIMAL(18, 2),
    `issued_date` DATE,
    `expected_clearance_date` DATE,
    `journal_entry_id` VARCHAR(36),
    `status` VARCHAR(50) DEFAULT 'outstanding',
    `cleared_date` DATE,
    `remarks` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`bank_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_bank_account` (`bank_account_id`),
    KEY `idx_status` (`status`),
    KEY `idx_cheque` (`cheque_number`),
    KEY `idx_item_type` (`item_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- CASH FLOW FORECAST TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `cash_flow_forecast` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `forecast_date` DATE NOT NULL,
    `forecast_period` VARCHAR(50),
    `opening_balance` DECIMAL(18, 2),
    `total_inflow` DECIMAL(18, 2) DEFAULT 0,
    `total_outflow` DECIMAL(18, 2) DEFAULT 0,
    `closing_balance` DECIMAL(18, 2),
    `forecast_type` VARCHAR(50),
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_forecast_date` (`forecast_date`),
    KEY `idx_period` (`forecast_period`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- CASH FLOW ITEM TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `cash_flow_item` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `cash_flow_forecast_id` VARCHAR(36) NOT NULL,
    `item_description` VARCHAR(255),
    `amount` DECIMAL(18, 2),
    `item_type` VARCHAR(50),
    `source_id` VARCHAR(36),
    `source_type` VARCHAR(50),
    `actual_amount` DECIMAL(18, 2),
    `variance` DECIMAL(18, 2),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`cash_flow_forecast_id`) REFERENCES `cash_flow_forecast`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_forecast` (`cash_flow_forecast_id`),
    KEY `idx_item_type` (`item_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
