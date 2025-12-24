-- ============================================================
-- MIGRATION 030: AUDIT LOGGING FOR RBAC
-- Date: December 23, 2025
-- Purpose: Track all RBAC operations for compliance
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- AUDIT LOG TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `audit_log` (
    `id` CHAR(26) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` CHAR(36),
    `action` VARCHAR(100) NOT NULL,
    `resource_id` VARCHAR(36),
    `resource_type` VARCHAR(50),
    `changes` LONGTEXT,
    `status` VARCHAR(20),
    `ip_address` VARCHAR(45),
    `user_agent` TEXT,
    `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE SET NULL,
    
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_user` (`user_id`),
    KEY `idx_action` (`action`),
    KEY `idx_status` (`status`),
    KEY `idx_timestamp` (`timestamp`),
    KEY `idx_tenant_timestamp` (`tenant_id`, `timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- AUDIT SUMMARY TABLE (for dashboards)
-- ============================================================
CREATE TABLE IF NOT EXISTS `audit_summary` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `date` DATE NOT NULL,
    `total_actions` INT DEFAULT 0,
    `successful_actions` INT DEFAULT 0,
    `failed_actions` INT DEFAULT 0,
    `unique_users` INT DEFAULT 0,
    `last_updated` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_date` (`date`),
    UNIQUE KEY `unique_tenant_date` (`tenant_id`, `date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
