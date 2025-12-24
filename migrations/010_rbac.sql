-- ============================================================
-- MIGRATION 010: RBAC (ROLE-BASED ACCESS CONTROL) MODULE
-- Date: December 3, 2025
-- Purpose: Create roles, permissions, and access control tables
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- ROLE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `role` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `role_name` VARCHAR(100) NOT NULL,
    `description` TEXT,
    `is_system_role` BOOLEAN DEFAULT FALSE,
    `is_active` BOOLEAN DEFAULT TRUE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    UNIQUE KEY `unique_role` (`tenant_id`, `role_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PERMISSION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `permission` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `permission_name` VARCHAR(100) NOT NULL,
    `description` TEXT,
    `resource` VARCHAR(100),
    `action` VARCHAR(50),
    `is_system_permission` BOOLEAN DEFAULT FALSE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_resource` (`resource`),
    UNIQUE KEY `unique_permission` (`tenant_id`, `permission_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ROLE PERMISSION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `role_permission` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `role_id` VARCHAR(36) NOT NULL,
    `permission_id` VARCHAR(36) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`role_id`) REFERENCES `role`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`permission_id`) REFERENCES `permission`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_role` (`role_id`),
    KEY `idx_permission` (`permission_id`),
    UNIQUE KEY `unique_role_permission` (`role_id`, `permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- USER ROLE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `user_role` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `role_id` VARCHAR(36) NOT NULL,
    `assigned_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `assigned_by` VARCHAR(36),
    `expires_at` TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`role_id`) REFERENCES `role`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_user` (`user_id`),
    KEY `idx_role` (`role_id`),
    UNIQUE KEY `unique_user_role` (`user_id`, `role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- RESOURCE TABLE (API ENDPOINTS & FEATURES)
-- ============================================================
CREATE TABLE IF NOT EXISTS `resource` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `resource_name` VARCHAR(100) NOT NULL,
    `resource_type` VARCHAR(50),
    `endpoint` VARCHAR(255),
    `description` TEXT,
    `is_protected` BOOLEAN DEFAULT TRUE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_type` (`resource_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ACCESS LOG TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `access_log` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT,
    `resource_id` VARCHAR(36),
    `action` VARCHAR(50),
    `status` VARCHAR(50),
    `ip_address` VARCHAR(45),
    `user_agent` TEXT,
    `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`resource_id`) REFERENCES `resource`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_user` (`user_id`),
    KEY `idx_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
