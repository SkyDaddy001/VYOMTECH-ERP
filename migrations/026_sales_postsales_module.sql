-- ============================================================
-- MIGRATION 026: SALES & POST-SALES MODULES
-- Date: December 8, 2025
-- Purpose: Create tables for real estate sales and post-sales management
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;
-- ============================================================
-- SALES & POST-SALES SPECIFIC TABLES
-- ============================================================
-- Note: sales_lead, unit, and unit_cost_sheet are defined in migrations 007, 003, and 008
-- This migration extends them with booking and client management features

-- Booking/Unit Allocation
CREATE TABLE IF NOT EXISTS `booking` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `booking_code` VARCHAR(50) UNIQUE NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `lead_id` VARCHAR(36) NOT NULL,
    `booking_date` DATE NOT NULL,
    `allotment_date` DATE,
    `agreement_date` DATE,
    `registration_date` DATE,
    `handover_date` DATE,
    `status` VARCHAR(50) DEFAULT 'active',
    `sales_executive_id` INT,
    `sales_head_id` INT,
    `booking_source` VARCHAR(100),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`lead_id`) REFERENCES `sales_lead`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_booking` (`tenant_id`, `status`),
    KEY `idx_unit` (`unit_id`),
    KEY `idx_booking_date` (`booking_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Applicant/Client Details
CREATE TABLE IF NOT EXISTS `client` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `booking_id` VARCHAR(36) NOT NULL,
    `applicant_type` VARCHAR(50),
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100),
    `phone` VARCHAR(20),
    `alternate_phone` VARCHAR(20),
    `email` VARCHAR(255),
    `communication_address` TEXT,
    `permanent_address` TEXT,
    `aadhar_no` VARCHAR(50),
    `pan_no` VARCHAR(50),
    `care_of` VARCHAR(100),
    `relation` VARCHAR(100),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`booking_id`) REFERENCES `booking`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_client` (`tenant_id`, `booking_id`),
    KEY `idx_aadhar` (`aadhar_no`),
    KEY `idx_pan` (`pan_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Power of Attorney (PoA) Details
CREATE TABLE IF NOT EXISTS `power_of_attorney` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `booking_id` VARCHAR(36) NOT NULL,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100),
    `phone` VARCHAR(20),
    `alternate_phone` VARCHAR(20),
    `email` VARCHAR(255),
    `communication_address` TEXT,
    `permanent_address` TEXT,
    `aadhar_no` VARCHAR(50),
    `pan_no` VARCHAR(50),
    `care_of` VARCHAR(100),
    `relation_to_applicant` VARCHAR(100),
    `poa_document_no` VARCHAR(100) UNIQUE,
    `life_certificate` VARCHAR(500),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`booking_id`) REFERENCES `booking`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_poa` (`tenant_id`, `booking_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- POST-SALES MODULE TABLES
-- ============================================================

-- Payment Schedule & Records
CREATE TABLE IF NOT EXISTS `payment_schedule` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `booking_id` VARCHAR(36) NOT NULL,
    `payment_stage` INT,
    `construction_stage` VARCHAR(100),
    `scheduled_date` DATE,
    `amount_due` DECIMAL(15, 2),
    `payment_type` VARCHAR(50),
    `status` VARCHAR(50) DEFAULT 'pending',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`booking_id`) REFERENCES `booking`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_payment_schedule` (`tenant_id`, `status`),
    KEY `idx_booking_payment` (`booking_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Payment Details/Records
CREATE TABLE IF NOT EXISTS `payment` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `booking_id` VARCHAR(36) NOT NULL,
    `payment_schedule_id` VARCHAR(36),
    `customer_name` VARCHAR(255),
    `unit_id` VARCHAR(36),
    `receipt_no` VARCHAR(100) UNIQUE,
    `received_on` DATE,
    `cleared_on` DATE,
    `payment_date` DATE NOT NULL,
    `payment_mode` VARCHAR(50),
    `paid_by` VARCHAR(255),
    `towards` VARCHAR(500),
    `amount` DECIMAL(15, 2) NOT NULL,
    `status` VARCHAR(50) DEFAULT 'received',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`booking_id`) REFERENCES `booking`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`payment_schedule_id`) REFERENCES `payment_schedule`(`id`) ON DELETE SET NULL,
    FOREIGN KEY (`unit_id`) REFERENCES `unit`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant_payment` (`tenant_id`, `payment_date`),
    KEY `idx_booking_payment` (`booking_id`),
    KEY `idx_receipt` (`receipt_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Bank Loan Details
CREATE TABLE IF NOT EXISTS `bank_loan` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `booking_id` VARCHAR(36) NOT NULL,
    `bank_name` VARCHAR(255) NOT NULL,
    `contact_person` VARCHAR(100),
    `phone` VARCHAR(20),
    `loan_sanction_date` DATE,
    `connector_code` VARCHAR(100),
    `sanction_amount` DECIMAL(15, 2),
    `disbursed_amount` DECIMAL(15, 2) DEFAULT 0,
    `disbursement_date` DATE,
    `disbursement_status` VARCHAR(50) DEFAULT 'pending' COMMENT 'pending, partial, completed',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`booking_id`) REFERENCES `booking`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_loan` (`tenant_id`, `booking_id`),
    KEY `idx_disbursement_status` (`disbursement_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Car Parking Details
CREATE TABLE IF NOT EXISTS `car_parking` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `booking_id` VARCHAR(36) NOT NULL,
    `parking_type` VARCHAR(50),
    `parking_number` VARCHAR(50),
    `parking_cost` DECIMAL(15, 2),
    `status` VARCHAR(50) DEFAULT 'allotted',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`booking_id`) REFERENCES `booking`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_parking` (`tenant_id`, `booking_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- GST/Registration Details
CREATE TABLE IF NOT EXISTS `registration_details` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `booking_id` VARCHAR(36) NOT NULL,
    `gst_applicable` BOOLEAN DEFAULT FALSE,
    `gst_percentage` DECIMAL(5, 2),
    `gst_cost` DECIMAL(15, 2),
    `apartment_cost_including_gst` DECIMAL(15, 2),
    `registration_type` VARCHAR(100),
    `registration_cost` DECIMAL(15, 2),
    `noc_received_date` DATE,
    `status` VARCHAR(50) DEFAULT 'pending',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`booking_id`) REFERENCES `booking`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_registration` (`tenant_id`, `booking_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Additional Charges/Maintenance
CREATE TABLE IF NOT EXISTS `additional_charges` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `booking_id` VARCHAR(36) NOT NULL,
    `maintenance_charge` DECIMAL(15, 2),
    `corpus_charge` DECIMAL(15, 2),
    `eb_deposit` DECIMAL(15, 2),
    `other_works_charge` DECIMAL(15, 2),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`booking_id`) REFERENCES `booking`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_charges` (`tenant_id`, `booking_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
