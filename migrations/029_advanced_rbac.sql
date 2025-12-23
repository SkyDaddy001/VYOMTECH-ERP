-- ============================================================
-- MIGRATION 029: ADVANCED RBAC (RESOURCE, TIME-BASED, FIELD-LEVEL PERMISSIONS)
-- Date: December 23, 2025
-- Purpose: Add resource-level, time-based, and field-level permission controls
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- RESOURCE ACCESS TABLE (Resource-Level Permissions)
-- ============================================================
CREATE TABLE IF NOT EXISTS `resource_access` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `resource_type` VARCHAR(50) NOT NULL,
    `resource_id` VARCHAR(36) NOT NULL,
    `access_level` VARCHAR(50) NOT NULL COMMENT 'view, edit, delete, admin',
    `expires_at` TIMESTAMP NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_user` (`user_id`),
    KEY `idx_resource` (`resource_type`, `resource_id`),
    KEY `idx_access_level` (`access_level`),
    UNIQUE KEY `unique_user_resource` (`user_id`, `resource_type`, `resource_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- TIME-BASED PERMISSION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `time_based_permission` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `role_id` VARCHAR(36) NOT NULL,
    `permission_id` VARCHAR(36) NOT NULL,
    `effective_from` TIMESTAMP NOT NULL,
    `expires_at` TIMESTAMP NOT NULL,
    `is_active` BOOLEAN DEFAULT TRUE,
    `reason` VARCHAR(255),
    `approved_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`role_id`) REFERENCES `role`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_role` (`role_id`),
    KEY `idx_expires` (`expires_at`),
    KEY `idx_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- FIELD-LEVEL PERMISSION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `field_level_permission` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `role_id` VARCHAR(36) NOT NULL,
    `module_name` VARCHAR(50) NOT NULL,
    `entity_name` VARCHAR(100) NOT NULL,
    `field_name` VARCHAR(100) NOT NULL,
    `can_view` BOOLEAN DEFAULT TRUE,
    `can_edit` BOOLEAN DEFAULT FALSE,
    `is_masked` BOOLEAN DEFAULT FALSE,
    `mask_pattern` VARCHAR(50),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`role_id`) REFERENCES `role`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_role` (`role_id`),
    KEY `idx_field` (`module_name`, `entity_name`, `field_name`),
    UNIQUE KEY `unique_field_permission` (`role_id`, `module_name`, `entity_name`, `field_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ROLE DELEGATION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `role_delegation` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `parent_role_id` VARCHAR(36) NOT NULL,
    `sub_role_id` VARCHAR(36) NOT NULL,
    `permission_bound` VARCHAR(50) COMMENT 'max permissions delegator can assign',
    `delegated_by` VARCHAR(36) NOT NULL,
    `is_active` BOOLEAN DEFAULT TRUE,
    `expires_at` TIMESTAMP NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`parent_role_id`) REFERENCES `role`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`sub_role_id`) REFERENCES `role`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_parent_role` (`parent_role_id`),
    KEY `idx_sub_role` (`sub_role_id`),
    UNIQUE KEY `unique_delegation` (`parent_role_id`, `sub_role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- BULK PERMISSION ASSIGNMENT LOG TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `bulk_permission_log` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `assignment_type` VARCHAR(50) NOT NULL COMMENT 'assign, revoke, update',
    `target_type` VARCHAR(50) NOT NULL COMMENT 'user, role, resource',
    `total_count` INT DEFAULT 0,
    `success_count` INT DEFAULT 0,
    `failed_count` INT DEFAULT 0,
    `executed_by` VARCHAR(36) NOT NULL,
    `execution_date` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `details` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_assignment_type` (`assignment_type`),
    KEY `idx_execution_date` (`execution_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
