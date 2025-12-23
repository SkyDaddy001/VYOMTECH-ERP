-- ============================================================
-- MIGRATION 034: SITE VISIT & LEAD MANAGEMENT
-- Date: [Current Date]
-- Purpose: Create tables for site visit scheduling and logging
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- SITE VISIT SCHEDULE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `site_visit_schedule` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `lead_id` VARCHAR(36) NULL,
    `visitor_name` VARCHAR(100) NOT NULL,
    `visitor_phone` VARCHAR(20),
    `visitor_email` VARCHAR(100),
    `scheduled_date` DATETIME NOT NULL,
    `scheduled_by` VARCHAR(36) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'scheduled' CHECK (`status` IN ('scheduled', 'completed', 'cancelled', 'no_show')),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`lead_id`) REFERENCES `leads`(`id`) ON DELETE SET NULL,
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_lead_id` (`lead_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_scheduled_date` (`scheduled_date`),
    INDEX `idx_tenant_status` (`tenant_id`, `status`),
    INDEX `idx_tenant_scheduled_date` (`tenant_id`, `scheduled_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- SITE VISIT LOG TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `site_visit_log` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `visit_schedule_id` VARCHAR(36) NOT NULL,
    `check_in_time` DATETIME,
    `check_out_time` DATETIME,
    `visited_by` VARCHAR(36) NOT NULL,
    `units_viewed` JSON,
    `feedback` TEXT,
    `follow_up_required` BOOLEAN DEFAULT FALSE,
    `next_followup_date` DATE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`visit_schedule_id`) REFERENCES `site_visit_schedule`(`id`) ON DELETE CASCADE,
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_visit_schedule_id` (`visit_schedule_id`),
    INDEX `idx_tenant_visit_schedule` (`tenant_id`, `visit_schedule_id`),
    INDEX `idx_check_in_time` (`check_in_time`),
    INDEX `idx_follow_up_required` (`follow_up_required`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
