-- ============================================================
-- MIGRATION 031: ROLE TEMPLATES SYSTEM
-- Date: December 23, 2025
-- Purpose: Support role templates and bulk role creation
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- ROLE TEMPLATE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `role_template` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `description` TEXT,
    `category` VARCHAR(50),
    `is_system_template` BOOLEAN DEFAULT FALSE,
    `is_active` BOOLEAN DEFAULT TRUE,
    `metadata` LONGTEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_category` (`category`),
    KEY `idx_is_system` (`is_system_template`),
    KEY `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- TEMPLATE INSTANCE TABLE (tracks role creation from templates)
-- ============================================================
CREATE TABLE IF NOT EXISTS `template_instance` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `template_id` VARCHAR(36) NOT NULL,
    `role_id` VARCHAR(36) NOT NULL,
    `created_by` INT NOT NULL,
    `customizations` LONGTEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`template_id`) REFERENCES `role_template`(`id`) ON DELETE SET NULL,
    FOREIGN KEY (`role_id`) REFERENCES `role`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`created_by`) REFERENCES `user`(`id`) ON DELETE SET NULL,
    
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_template` (`template_id`),
    KEY `idx_role` (`role_id`),
    KEY `idx_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
