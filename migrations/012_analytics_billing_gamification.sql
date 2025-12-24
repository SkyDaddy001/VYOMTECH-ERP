-- ============================================================
-- MIGRATION 012: ANALYTICS, BILLING & GAMIFICATION MODULE
-- Date: December 3, 2025
-- Purpose: Create analytics, billing, subscription, and gamification tables
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- ANALYTICS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `analytics` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `metric_name` VARCHAR(100) NOT NULL,
    `metric_category` VARCHAR(50),
    `metric_value` DECIMAL(18, 2),
    `metric_date` DATE NOT NULL,
    `dimensions` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_metric` (`metric_name`),
    KEY `idx_date` (`metric_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- DASHBOARD WIDGET TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `dashboard_widget` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT,
    `widget_type` VARCHAR(50),
    `widget_title` VARCHAR(255),
    `metric_id` VARCHAR(36),
    `position` INT,
    `is_visible` BOOLEAN DEFAULT TRUE,
    `config` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- BILLING SUBSCRIPTION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `billing_subscription` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `subscription_plan` VARCHAR(100) NOT NULL,
    `plan_type` VARCHAR(50),
    `billing_cycle` VARCHAR(50),
    `amount` DECIMAL(18, 2),
    `currency` VARCHAR(10) DEFAULT 'INR',
    `start_date` DATE NOT NULL,
    `end_date` DATE,
    `auto_renew` BOOLEAN DEFAULT TRUE,
    `status` VARCHAR(50) DEFAULT 'active',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- BILLING INVOICE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `billing_invoice` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `invoice_number` VARCHAR(50) UNIQUE NOT NULL,
    `subscription_id` VARCHAR(36),
    `invoice_date` DATE NOT NULL,
    `due_date` DATE,
    `amount` DECIMAL(18, 2),
    `tax_amount` DECIMAL(18, 2) DEFAULT 0,
    `total_amount` DECIMAL(18, 2),
    `amount_paid` DECIMAL(18, 2) DEFAULT 0,
    `status` VARCHAR(50) DEFAULT 'pending',
    `payment_method` VARCHAR(50),
    `transaction_id` VARCHAR(100),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`subscription_id`) REFERENCES `billing_subscription`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_number` (`invoice_number`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PAYMENT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `payment` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `invoice_id` VARCHAR(36),
    `amount` DECIMAL(18, 2),
    `payment_method` VARCHAR(50),
    `transaction_id` VARCHAR(100),
    `status` VARCHAR(50) DEFAULT 'pending',
    `payment_date` TIMESTAMP,
    `reference_number` VARCHAR(100),
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`invoice_id`) REFERENCES `billing_invoice`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- GAMIFICATION LEVEL TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `gamification_level` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `level_name` VARCHAR(100) NOT NULL,
    `level_number` INT,
    `min_points` INT DEFAULT 0,
    `max_points` INT,
    `badge_icon_url` TEXT,
    `description` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- GAMIFICATION ACHIEVEMENT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `gamification_achievement` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `achievement_name` VARCHAR(255) NOT NULL,
    `achievement_type` VARCHAR(50),
    `points_reward` INT,
    `badge_icon_url` TEXT,
    `description` TEXT,
    `criteria` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- USER GAMIFICATION PROFILE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `user_gamification` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `total_points` INT DEFAULT 0,
    `current_level_id` VARCHAR(36),
    `badges_earned` JSON,
    `achievements_unlocked` JSON,
    `last_updated` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`current_level_id`) REFERENCES `gamification_level`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_user` (`user_id`),
    UNIQUE KEY `unique_user_gamification` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- LEADERBOARD TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `leaderboard` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `rank` INT,
    `points` INT,
    `period` VARCHAR(50),
    `period_start_date` DATE,
    `period_end_date` DATE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_period` (`period`, `period_start_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
