-- ============================================================
-- MIGRATION 003: CONSTRUCTION MODULE
-- Date: December 3, 2025
-- Purpose: Create construction project management tables
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- CONSTRUCTION PROJECTS
-- ============================================================
CREATE TABLE IF NOT EXISTS `construction_projects` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_name` VARCHAR(255) NOT NULL,
    `project_code` VARCHAR(100) UNIQUE NOT NULL,
    `location` VARCHAR(255),
    `client` VARCHAR(255),
    `contract_value` DECIMAL(15, 2),
    `start_date` TIMESTAMP,
    `expected_completion` TIMESTAMP,
    `current_progress_percent` INT DEFAULT 0,
    `status` VARCHAR(50) DEFAULT 'planning',
    `project_manager` VARCHAR(255),
    `gl_project_code` VARCHAR(50),
    `project_cost_account_id` VARCHAR(36),
    `revenue_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_cost_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`revenue_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_status` (`status`),
    KEY `idx_code` (`project_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- BILL OF QUANTITIES
-- ============================================================
CREATE TABLE IF NOT EXISTS `bill_of_quantities` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` BIGINT NOT NULL,
    `boq_number` VARCHAR(100),
    `item_description` TEXT,
    `unit` VARCHAR(50),
    `quantity` DECIMAL(15, 4),
    `unit_rate` DECIMAL(15, 2),
    `total_amount` DECIMAL(15, 2),
    `category` VARCHAR(100),
    `status` VARCHAR(50) DEFAULT 'planned',
    `gl_expense_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`gl_expense_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`project_id`) REFERENCES `construction_projects`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PROGRESS TRACKING
-- ============================================================
CREATE TABLE IF NOT EXISTS `progress_tracking` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` BIGINT NOT NULL,
    `date` TIMESTAMP NOT NULL,
    `activity_desc` LONGTEXT,
    `quantity_completed` DECIMAL(15, 4),
    `unit` VARCHAR(50),
    `percent_complete` INT,
    `workforce_deployed` INT,
    `notes` LONGTEXT,
    `photo_url` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `construction_projects`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_date` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- QUALITY CONTROL
-- ============================================================
CREATE TABLE IF NOT EXISTS `quality_control` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` BIGINT NOT NULL,
    `boq_item_id` BIGINT,
    `inspection_date` TIMESTAMP NOT NULL,
    `inspector_name` VARCHAR(255),
    `quality_status` VARCHAR(50) DEFAULT 'pending',
    `observations` LONGTEXT,
    `corrective_action` LONGTEXT,
    `follow_up_date` TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `construction_projects`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`boq_item_id`) REFERENCES `bill_of_quantities`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_status` (`quality_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- CONSTRUCTION EQUIPMENT
-- ============================================================
CREATE TABLE IF NOT EXISTS `construction_equipment` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` BIGINT NOT NULL,
    `equipment_name` VARCHAR(255),
    `equipment_type` VARCHAR(100),
    `serial_number` VARCHAR(100),
    `status` VARCHAR(50) DEFAULT 'available',
    `deployment_date` TIMESTAMP,
    `retirement_date` TIMESTAMP,
    `cost_per_day` DECIMAL(15, 2),
    `notes` LONGTEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `construction_projects`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
