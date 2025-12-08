-- ============================================================
-- MIGRATION 008: REAL ESTATE MODULE
-- Date: December 3, 2025
-- Purpose: Create property projects, units, blocks, and cost sheets
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- PROPERTY PROJECT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_project` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_name` VARCHAR(255) NOT NULL,
    `project_code` VARCHAR(50) UNIQUE NOT NULL,
    `location` VARCHAR(255),
    `city` VARCHAR(100),
    `state` VARCHAR(100),
    `postal_code` VARCHAR(20),
    `total_units` INT,
    `total_area` DECIMAL(18, 2),
    `project_type` VARCHAR(50),
    `status` VARCHAR(50) DEFAULT 'planning',
    `launch_date` DATE,
    `expected_completion` DATE,
    `actual_completion` DATE,
    `noc_status` VARCHAR(50),
    `noc_date` DATE,
    `developer_name` VARCHAR(255),
    `architect_name` VARCHAR(255),
    `created_by` VARCHAR(36),
    `gl_asset_account_id` VARCHAR(36),
    `gl_revenue_account_id` VARCHAR(36),
    `project_manager_id` VARCHAR(36),
    `total_estimated_cost` DECIMAL(18,2),
    `total_actual_cost` DECIMAL(18,2),
    `project_margin` DECIMAL(5,2),
    `bank_loan_amount` DECIMAL(18,2),
    `equity_amount` DECIMAL(18,2),
    `financial_status` VARCHAR(50) DEFAULT 'PLANNING',
    `brochure_url` VARCHAR(500),
    `master_plan_url` VARCHAR(500),
    `legal_status` VARCHAR(50) DEFAULT 'PENDING',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`gl_asset_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`gl_revenue_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_code` (`project_code`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PROPERTY BLOCK TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_block` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `block_name` VARCHAR(100) NOT NULL,
    `block_code` VARCHAR(50) NOT NULL,
    `wing_name` VARCHAR(100),
    `total_units` INT,
    `status` VARCHAR(50) DEFAULT 'planning',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    UNIQUE KEY `unique_block` (`project_id`, `block_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PROPERTY UNIT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_unit` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `block_id` VARCHAR(36) NOT NULL,
    `unit_number` VARCHAR(50) NOT NULL,
    `floor` INT,
    `unit_type` VARCHAR(50),
    `facing` VARCHAR(50),
    `carpet_area` DECIMAL(15, 2),
    `carpet_area_with_balcony` DECIMAL(15, 2),
    `utility_area` DECIMAL(15, 2),
    `plinth_area` DECIMAL(15, 2),
    `sbua` DECIMAL(15, 2),
    `uds_sqft` DECIMAL(15, 2),
    `status` VARCHAR(50) DEFAULT 'available',
    `alloted_to` VARCHAR(255),
    `allotment_date` DATE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`block_id`) REFERENCES `property_block`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_block` (`block_id`),
    KEY `idx_status` (`status`),
    UNIQUE KEY `unique_unit` (`block_id`, `unit_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- UNIT COST SHEET TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `unit_cost_sheet` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `rate_per_sqft` DECIMAL(18, 2),
    `sbua_rate` DECIMAL(18, 2),
    `base_price` DECIMAL(18, 2),
    `frc` DECIMAL(18, 2) DEFAULT 0,
    `car_parking_cost` DECIMAL(18, 2) DEFAULT 0,
    `plc` DECIMAL(18, 2) DEFAULT 0,
    `statutory_charges` DECIMAL(18, 2) DEFAULT 0,
    `other_charges` DECIMAL(18, 2) DEFAULT 0,
    `legal_charges` DECIMAL(18, 2) DEFAULT 0,
    `apartment_cost_exc_govt` DECIMAL(18, 2),
    `apartment_cost_inc_govt` DECIMAL(18, 2),
    `composite_guideline_value` DECIMAL(18, 2),
    `actual_sold_price` DECIMAL(18, 2),
    `car_parking_type` VARCHAR(50),
    `parking_location` VARCHAR(255),
    `block_name` VARCHAR(100),
    `sbua` DECIMAL(15,2),
    `statutory_approval_charge` DECIMAL(18,2),
    `legal_documentation_charge` DECIMAL(18,2),
    `amenities_equipment_charge` DECIMAL(18,2),
    `other_charges_1` DECIMAL(18,2),
    `other_charges_1_name` VARCHAR(100),
    `other_charges_1_type` VARCHAR(20),
    `other_charges_2` DECIMAL(18,2),
    `other_charges_2_name` VARCHAR(100),
    `other_charges_2_type` VARCHAR(20),
    `other_charges_3` DECIMAL(18,2),
    `other_charges_3_name` VARCHAR(100),
    `other_charges_3_type` VARCHAR(20),
    `other_charges_4` DECIMAL(18,2),
    `other_charges_4_name` VARCHAR(100),
    `other_charges_4_type` VARCHAR(20),
    `other_charges_5` DECIMAL(18,2),
    `other_charges_5_name` VARCHAR(100),
    `other_charges_5_type` VARCHAR(20),
    `apartment_cost_excluding_govt` DECIMAL(18,2),
    `actual_sold_price_excluding_govt` DECIMAL(18,2),
    `gst_applicable` TINYINT(1) DEFAULT 0,
    `gst_percentage` DECIMAL(5,2) DEFAULT 0,
    `gst_amount` DECIMAL(18,2),
    `grand_total` DECIMAL(18,2),
    `club_membership` DECIMAL(18,2),
    `registration_charge` DECIMAL(18,2),
    `other_charges_json` JSON,
    `effective_date` DATE,
    `valid_until` DATE,
    `cost_sheet_created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_unit` (`unit_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PROPERTY BOOKING TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_booking` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `customer_id` VARCHAR(36),
    `booking_date` DATE NOT NULL,
    `booking_status` VARCHAR(50) DEFAULT 'confirmed',
    `booking_amount` DECIMAL(18, 2),
    `agreement_value` DECIMAL(18, 2),
    `booking_notes` TEXT,
    `created_by` VARCHAR(36),
    `gl_receivable_account_id` VARCHAR(36),
    `gl_revenue_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`),
    FOREIGN KEY (`gl_receivable_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`gl_revenue_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_unit` (`unit_id`),
    KEY `idx_date` (`booking_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PAYMENT PLAN TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `payment_plan` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `booking_id` VARCHAR(36) NOT NULL,
    `plan_name` VARCHAR(100) NOT NULL,
    `total_amount` DECIMAL(18, 2),
    `number_of_installments` INT,
    `plan_description` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`booking_id`) REFERENCES `property_booking`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_booking` (`booking_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- INSTALLMENT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `installment` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `payment_plan_id` VARCHAR(36) NOT NULL,
    `installment_number` INT NOT NULL,
    `due_date` DATE NOT NULL,
    `amount_due` DECIMAL(18, 2),
    `amount_paid` DECIMAL(18, 2) DEFAULT 0,
    `status` VARCHAR(50) DEFAULT 'pending',
    `payment_date` DATE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`payment_plan_id`) REFERENCES `payment_plan`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_plan` (`payment_plan_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
