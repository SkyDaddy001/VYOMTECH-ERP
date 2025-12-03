-- Migration: 022_project_management_extensions.sql
-- Description: Extensions to real estate module (migration 008) with enhanced customer management, detailed cost sheets, and project milestones
-- Date: 2025-12-03
-- Note: Integrates with existing property_project, property_unit, property_block tables from migration 008

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- ALTER EXISTING TABLES TO ADD MISSING FIELDS
-- ============================================================

-- Extend property_project with additional fields
ALTER TABLE `property_project` 
ADD COLUMN IF NOT EXISTS `project_manager_id` VARCHAR(36),
ADD COLUMN IF NOT EXISTS `total_estimated_cost` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `total_actual_cost` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `project_margin` DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS `bank_loan_amount` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `equity_amount` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `financial_status` VARCHAR(50) DEFAULT 'PLANNING',
ADD COLUMN IF NOT EXISTS `brochure_url` VARCHAR(500),
ADD COLUMN IF NOT EXISTS `master_plan_url` VARCHAR(500),
ADD COLUMN IF NOT EXISTS `legal_status` VARCHAR(50) DEFAULT 'PENDING';

-- ============================================================
-- ENHANCED UNIT COST SHEET (EXTENDS EXISTING TABLE)
-- ============================================================
-- The existing unit_cost_sheet from migration 008 already has most fields
-- Adding detailed fields for comprehensive costing with project-wise charge configuration

ALTER TABLE `unit_cost_sheet`
ADD COLUMN IF NOT EXISTS `block_name` VARCHAR(100),
ADD COLUMN IF NOT EXISTS `sbua` DECIMAL(15,2),
ADD COLUMN IF NOT EXISTS `rate_per_sqft` DECIMAL(15,2),
ADD COLUMN IF NOT EXISTS `car_parking_cost` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `plc` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `statutory_approval_charge` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `legal_documentation_charge` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `amenities_equipment_charge` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `other_charges_1` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `other_charges_1_name` VARCHAR(100),
ADD COLUMN IF NOT EXISTS `other_charges_1_type` VARCHAR(20),
ADD COLUMN IF NOT EXISTS `other_charges_2` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `other_charges_2_name` VARCHAR(100),
ADD COLUMN IF NOT EXISTS `other_charges_2_type` VARCHAR(20),
ADD COLUMN IF NOT EXISTS `other_charges_3` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `other_charges_3_name` VARCHAR(100),
ADD COLUMN IF NOT EXISTS `other_charges_3_type` VARCHAR(20),
ADD COLUMN IF NOT EXISTS `other_charges_4` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `other_charges_4_name` VARCHAR(100),
ADD COLUMN IF NOT EXISTS `other_charges_4_type` VARCHAR(20),
ADD COLUMN IF NOT EXISTS `other_charges_5` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `other_charges_5_name` VARCHAR(100),
ADD COLUMN IF NOT EXISTS `other_charges_5_type` VARCHAR(20),
ADD COLUMN IF NOT EXISTS `apartment_cost_excluding_govt` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `actual_sold_price_excluding_govt` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `gst_applicable` TINYINT(1) DEFAULT 0,
ADD COLUMN IF NOT EXISTS `gst_percentage` DECIMAL(5,2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS `gst_amount` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `grand_total` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `club_membership` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `registration_charge` DECIMAL(18,2),
ADD COLUMN IF NOT EXISTS `other_charges_json` JSON,
ADD COLUMN IF NOT EXISTS `effective_date` DATE,
ADD COLUMN IF NOT EXISTS `valid_until` DATE,
ADD COLUMN IF NOT EXISTS `created_by` VARCHAR(36);

-- ============================================================
-- AREA STATEMENT / AREA BREAKUP (NEW TABLE)
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_unit_area_statement` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `block_id` VARCHAR(36),
    `unit_id` VARCHAR(36) NOT NULL,
    -- UNIT IDENTIFICATION
    `apt_no` VARCHAR(50) COMMENT 'Apartment Number',
    `floor` VARCHAR(50) COMMENT 'Floor Number',
    `unit_type` VARCHAR(50) COMMENT 'Unit Type: 1BHK, 2BHK, 3BHK, Studio, etc.',
    `facing` VARCHAR(50) COMMENT 'Facing Direction: NORTH, SOUTH, EAST, WEST, NE, NW, SE, SW',
    -- AREA MEASUREMENTS - CARPET AREA
    `rera_carpet_area_sqft` DECIMAL(15,2) COMMENT 'RERA Carpet Area in Sq.Ft.',
    `rera_carpet_area_sqm` DECIMAL(15,2) COMMENT 'RERA Carpet Area in Sq.Mtrs.',
    -- AREA MEASUREMENTS - WITH BALCONY & UTILITY
    `carpet_area_with_balcony_sqft` DECIMAL(15,2) COMMENT 'Carpet Area with Balcony and Utility in Sq.Ft.',
    `carpet_area_with_balcony_sqm` DECIMAL(15,2) COMMENT 'Carpet Area with Balcony and Utility in Sq.Mtrs.',
    -- PLINTH AREA
    `plinth_area_sqft` DECIMAL(15,2) COMMENT 'Plinth Area in Sq.Ft.',
    `plinth_area_sqm` DECIMAL(15,2) COMMENT 'Plinth Area in Sq.Mtrs.',
    -- SUPER BUILT-UP AREA (SBUA)
    `sbua_sqft` DECIMAL(15,2) COMMENT 'Super Built-Up Area (SBUA) in Sq.Ft.',
    `sbua_sqm` DECIMAL(15,2) COMMENT 'Super Built-Up Area (SBUA) in Sq.Mtrs.',
    -- UNDIVIDED SHARE
    `uds_per_sqft` DECIMAL(15,2) COMMENT 'Undivided Share per Sq.Ft.',
    `uds_total_sqft` DECIMAL(15,2) COMMENT 'Total Undivided Share in Sq.Ft.',
    -- ADDITIONAL AREAS
    `balcony_area_sqft` DECIMAL(15,2),
    `balcony_area_sqm` DECIMAL(15,2),
    `utility_area_sqft` DECIMAL(15,2),
    `utility_area_sqm` DECIMAL(15,2),
    `garden_area_sqft` DECIMAL(15,2),
    `garden_area_sqm` DECIMAL(15,2),
    `terrace_area_sqft` DECIMAL(15,2),
    `terrace_area_sqm` DECIMAL(15,2),
    `parking_area_sqft` DECIMAL(15,2),
    `parking_area_sqm` DECIMAL(15,2),
    `common_area_sqft` DECIMAL(15,2),
    `common_area_sqm` DECIMAL(15,2),
    -- OWNERSHIP & ALLOCATION
    `alloted_to` VARCHAR(200) COMMENT 'Name of primary allottee',
    `key_holder` VARCHAR(200) COMMENT 'Key holder details',
    `percentage_allocation` DECIMAL(5,2) COMMENT 'Percentage share of common area',
    -- NOC & COMPLIANCE
    `noc_taken` VARCHAR(50) COMMENT 'NOC Status: YES, NO, PENDING, NA',
    `noc_date` DATE,
    `noc_document_url` VARCHAR(500),
    -- OTHER DETAILS
    `area_type` VARCHAR(50) COMMENT 'CARPET_AREA, BUILTUP_AREA, SUPER_AREA, COMMON_AREA, GARDEN_AREA, TERRACE_AREA, BALCONY_AREA',
    `description` TEXT,
    `active` TINYINT(1) DEFAULT 1,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_unit` (`unit_id`),
    KEY `idx_apt_no` (`apt_no`),
    KEY `idx_floor` (`floor`),
    KEY `idx_area_type` (`area_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Area Statement: Complete area breakup for units (carpet, plinth, SBUA, balcony, utility, etc.) with RERA compliance';

-- ============================================================
-- PROJECT COST CONFIGURATION (NEW TABLE - COST STRUCTURE SETUP)
-- ============================================================
CREATE TABLE IF NOT EXISTS `project_cost_configuration` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `config_name` VARCHAR(100) NOT NULL COMMENT 'e.g., CMWSSB, Water Tax, Electricity Deposit, etc.',
    `config_type` VARCHAR(50) NOT NULL COMMENT 'OTHER_CHARGE_1, OTHER_CHARGE_2, etc. or CUSTOM',
    `display_order` INT DEFAULT 0,
    `is_mandatory` TINYINT(1) DEFAULT 0 COMMENT 'Is this charge mandatory for all units',
    `applicable_for_unit_type` VARCHAR(100) COMMENT 'Comma-separated unit types (1BHK, 2BHK, etc.) or null for all',
    `description` TEXT,
    `active` TINYINT(1) DEFAULT 1,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_config_type` (`config_type`),
    UNIQUE KEY `unique_project_config` (`tenant_id`, `project_id`, `config_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Project-wise cost charge configuration for flexible other charges';

-- ============================================================
-- ENHANCED CUSTOMER PROFILE (NEW TABLE - EXTENDS BOOKING CONCEPT)
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_customer_profile` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `customer_code` VARCHAR(50) NOT NULL,
    `unit_id` VARCHAR(36),
    -- PRIMARY APPLICANT INFORMATION
    `first_name` VARCHAR(100) NOT NULL,
    `middle_name` VARCHAR(100),
    `last_name` VARCHAR(100),
    `email` VARCHAR(100),
    `phone_primary` VARCHAR(20),
    `phone_secondary` VARCHAR(20),
    `alternate_phone` VARCHAR(20),
    `company_name` VARCHAR(200),
    `designation` VARCHAR(100),
    `pan_number` VARCHAR(20),
    `aadhar_number` VARCHAR(20),
    `pan_copy_url` VARCHAR(500),
    `aadhar_copy_url` VARCHAR(500),
    `poa_document_no` VARCHAR(100),
    `care_of` VARCHAR(200),
    -- ADDRESSES
    `customer_type` VARCHAR(50) DEFAULT 'INDIVIDUAL' COMMENT 'INDIVIDUAL, JOINT, CORPORATE, NRI, HUF',
    `communication_address_line1` TEXT,
    `communication_address_line2` TEXT,
    `communication_city` VARCHAR(100),
    `communication_state` VARCHAR(100),
    `communication_country` VARCHAR(100),
    `communication_zip` VARCHAR(20),
    `permanent_address_line1` TEXT,
    `permanent_address_line2` TEXT,
    `permanent_city` VARCHAR(100),
    `permanent_state` VARCHAR(100),
    `permanent_country` VARCHAR(100),
    `permanent_zip` VARCHAR(20),
    -- EMPLOYMENT & FINANCIAL
    `profession` VARCHAR(100),
    `employer_name` VARCHAR(200),
    `employment_type` VARCHAR(50),
    `monthly_income` DECIMAL(18,2),
    -- CO-APPLICANT 1
    `co_applicant_1_name` VARCHAR(100),
    `co_applicant_1_number` VARCHAR(20),
    `co_applicant_1_alternate_number` VARCHAR(20),
    `co_applicant_1_email` VARCHAR(100),
    `co_applicant_1_communication_address` TEXT,
    `co_applicant_1_permanent_address` TEXT,
    `co_applicant_1_aadhar` VARCHAR(20),
    `co_applicant_1_pan` VARCHAR(20),
    `co_applicant_1_care_of` VARCHAR(200),
    `co_applicant_1_relation` VARCHAR(50),
    -- CO-APPLICANT 2
    `co_applicant_2_name` VARCHAR(100),
    `co_applicant_2_number` VARCHAR(20),
    `co_applicant_2_alternate_number` VARCHAR(20),
    `co_applicant_2_email` VARCHAR(100),
    `co_applicant_2_communication_address` TEXT,
    `co_applicant_2_permanent_address` TEXT,
    `co_applicant_2_aadhar` VARCHAR(20),
    `co_applicant_2_pan` VARCHAR(20),
    `co_applicant_2_care_of` VARCHAR(200),
    `co_applicant_2_relation` VARCHAR(50),
    -- CO-APPLICANT 3
    `co_applicant_3_name` VARCHAR(100),
    `co_applicant_3_number` VARCHAR(20),
    `co_applicant_3_alternate_number` VARCHAR(20),
    `co_applicant_3_email` VARCHAR(100),
    `co_applicant_3_communication_address` TEXT,
    `co_applicant_3_permanent_address` TEXT,
    `co_applicant_3_aadhar` VARCHAR(20),
    `co_applicant_3_pan` VARCHAR(20),
    `co_applicant_3_care_of` VARCHAR(200),
    `co_applicant_3_relation` VARCHAR(50),
    -- LOAN & FINANCING
    `loan_required` TINYINT(1) DEFAULT 0,
    `loan_amount` DECIMAL(18,2),
    `loan_sanction_date` DATE,
    `bank_name` VARCHAR(200),
    `bank_branch` VARCHAR(200),
    `bank_contact_person` VARCHAR(100),
    `bank_contact_number` VARCHAR(20),
    -- SALES & CRM
    `connector_code_number` VARCHAR(50),
    `lead_id` VARCHAR(50),
    `sales_executive_id` VARCHAR(36),
    `sales_executive_name` VARCHAR(100),
    `sales_head_id` VARCHAR(36),
    `sales_head_name` VARCHAR(100),
    `booking_source` VARCHAR(100),
    `life_certificate` VARCHAR(500),
    -- STATUS & DATES
    `customer_status` VARCHAR(50) DEFAULT 'INQUIRY' COMMENT 'INQUIRY, REGISTERED, BOOKING_CONFIRMED, AGREEMENT_SIGNED, REGISTERED_PROPERTY, HANDED_OVER, DEFAULTER',
    `booking_date` DATE,
    `welcome_date` DATE,
    `allotment_date` DATE,
    `agreement_date` DATE,
    `registration_date` DATE,
    `handover_date` DATE,
    `noc_received_date` DATE,
    -- CHARGES & MAINTENANCE
    `rate_per_sqft` DECIMAL(10,2),
    `composite_guideline_value` DECIMAL(18,2),
    `car_parking_type` VARCHAR(50),
    `maintenance_charges` DECIMAL(18,2),
    `other_works_charges` DECIMAL(18,2),
    `corpus_charges` DECIMAL(18,2),
    `eb_deposit` DECIMAL(18,2),
    -- NOTES & METADATA
    `notes` TEXT,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_code` (`customer_code`),
    KEY `idx_unit` (`unit_id`),
    KEY `idx_status` (`customer_status`),
    KEY `idx_email` (`email`),
    KEY `idx_pan` (`pan_number`),
    KEY `idx_aadhar` (`aadhar_number`),
    UNIQUE KEY `unique_customer_code` (`tenant_id`, `customer_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Enhanced customer profile with comprehensive personal, financial, employment, and co-applicant information';

-- ============================================================
-- CUSTOMER UNIT LINK (LINKS CUSTOMER TO PROPERTY UNIT AND BOOKING)
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_customer_unit_link` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `customer_id` VARCHAR(36) NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `booking_id` VARCHAR(36),
    `booking_status` VARCHAR(50) DEFAULT 'CONFIRMED',
    `booking_date` DATE,
    `agreement_date` DATE,
    `possession_date` DATE,
    `handover_date` DATE,
    `primary_customer` TINYINT(1) DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`customer_id`) REFERENCES `property_customer_profile`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`booking_id`) REFERENCES `property_booking`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_customer` (`customer_id`),
    KEY `idx_unit` (`unit_id`),
    UNIQUE KEY `unique_customer_unit` (`tenant_id`, `customer_id`, `unit_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Link customers to units and bookings';

-- ============================================================
-- ENHANCED PAYMENT TRACKING (EXTENDS INSTALLMENT)
-- ============================================================
-- The existing installment table from migration 008 already tracks payments
-- Creating a more detailed payment receipt table for comprehensive tracking

CREATE TABLE IF NOT EXISTS `property_payment_receipt` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `customer_id` VARCHAR(36) NOT NULL,
    `customer_name` VARCHAR(200),
    `unit_id` VARCHAR(36) NOT NULL,
    `installment_id` VARCHAR(36),
    `receipt_number` VARCHAR(50) NOT NULL,
    `receipt_date` DATE,
    `payment_date` DATE NOT NULL,
    `payment_mode` VARCHAR(50) NOT NULL COMMENT 'CASH, CHEQUE, NEFT, RTGS, ONLINE, DD',
    `payment_amount` DECIMAL(18,2),
    `installment_amount_due` DECIMAL(18,2),
    `shortfall_amount` DECIMAL(18,2),
    `excess_amount` DECIMAL(18,2),
    `payment_status` VARCHAR(50) DEFAULT 'PENDING' COMMENT 'PENDING, RECEIVED, PROCESSED, CLEARED, BOUNCED, CANCELLED',
    `bank_name` VARCHAR(200),
    `cheque_number` VARCHAR(50),
    `cheque_date` DATE,
    `transaction_id` VARCHAR(100),
    `account_number` VARCHAR(50),
    `towards_description` VARCHAR(200) COMMENT 'APARTMENT_COST, PARKING, MAINTENANCE, CORPUS, REGISTRATION, etc.',
    `received_in_bank_account` VARCHAR(50),
    `paid_by` VARCHAR(100) COMMENT 'Customer, Agent, Representative, etc.',
    `remarks` TEXT,
    `gl_account_id` VARCHAR(36),
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`customer_id`) REFERENCES `property_customer_profile`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`installment_id`) REFERENCES `installment`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_customer` (`customer_id`),
    KEY `idx_unit` (`unit_id`),
    KEY `idx_date` (`payment_date`),
    KEY `idx_status` (`payment_status`),
    UNIQUE KEY `unique_receipt` (`tenant_id`, `receipt_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Payment receipts with detailed transaction tracking';

-- ============================================================
-- PROJECT CONSTRUCTION MILESTONES (NEW TABLE)
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_project_milestone` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `block_id` VARCHAR(36),
    `milestone_name` VARCHAR(255) NOT NULL,
    `milestone_type` VARCHAR(50) NOT NULL COMMENT 'PLANNING, APPROVALS, LAND_ACQUISITION, DESIGN, CONSTRUCTION, TESTING, COMPLETION, HANDOVER',
    `milestone_description` TEXT,
    `start_date` DATE,
    `planned_completion_date` DATE,
    `actual_completion_date` DATE,
    `completion_status` VARCHAR(50) DEFAULT 'NOT_STARTED' COMMENT 'NOT_STARTED, IN_PROGRESS, ON_HOLD, COMPLETED, DELAYED, CANCELLED',
    `percentage_completion` INT DEFAULT 0,
    `responsible_party_id` VARCHAR(36),
    `responsible_party_type` VARCHAR(50) DEFAULT 'CONTRACTOR' COMMENT 'CONTRACTOR, CONSULTANT, ARCHITECT, ENGINEER, INTERNAL',
    `budget_allocated` DECIMAL(18,2),
    `budget_spent` DECIMAL(18,2),
    `budget_variance` DECIMAL(18,2),
    `priority` INT DEFAULT 0,
    `notes` TEXT,
    `documents_url` VARCHAR(500),
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`block_id`) REFERENCES `property_block`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_status` (`completion_status`),
    KEY `idx_type` (`milestone_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Construction project milestones and phases';

-- ============================================================
-- PROJECT ACTIVITY LOG (NEW TABLE)
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_project_activity` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `milestone_id` VARCHAR(36),
    `activity_type` VARCHAR(50) NOT NULL COMMENT 'WORK_START, WORK_COMPLETION, APPROVAL, INSPECTION, PAYMENT_DISBURSED, ISSUE_RAISED, ISSUE_RESOLVED, DELAY_REPORTED, DOCUMENTATION',
    `activity_description` TEXT NOT NULL,
    `activity_date` DATE,
    `activity_time` TIME,
    `assigned_to` VARCHAR(36),
    `status` VARCHAR(50) DEFAULT 'PENDING',
    `completion_date` DATE,
    `completion_percentage` INT DEFAULT 0,
    `attachments_url` VARCHAR(500),
    `notes` TEXT,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`milestone_id`) REFERENCES `property_project_milestone`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_date` (`activity_date`),
    KEY `idx_type` (`activity_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Project activity and work logs';

-- ============================================================
-- PROJECT DOCUMENTS (NEW TABLE)
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_project_document` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `document_type` VARCHAR(50) NOT NULL COMMENT 'APPROVAL, NOC, PLAN, SPECIFICATION, COMPLIANCE, CERTIFICATE, INSURANCE, AGREEMENT, TENDER, CONTRACT',
    `document_name` VARCHAR(255),
    `document_description` TEXT,
    `document_url` VARCHAR(500),
    `file_size_bytes` INT,
    `file_type` VARCHAR(100),
    `uploaded_by` VARCHAR(36),
    `upload_date` TIMESTAMP,
    `expiry_date` DATE,
    `is_active` TINYINT(1) DEFAULT 1,
    `approval_status` VARCHAR(50) DEFAULT 'PENDING',
    `approved_by` VARCHAR(36),
    `approval_date` DATE,
    `version_number` INT DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_type` (`document_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Project approvals and compliance documents';

-- ============================================================
-- PROJECT SUMMARY DASHBOARD (NEW TABLE)
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_project_summary` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `summary_date` DATE,
    `total_units` INT,
    `units_available` INT,
    `units_reserved` INT,
    `units_sold` INT,
    `units_handed_over` INT,
    `total_revenue_booked` DECIMAL(18,2),
    `total_revenue_received` DECIMAL(18,2),
    `total_construction_cost` DECIMAL(18,2),
    `total_cost_incurred` DECIMAL(18,2),
    `gross_profit` DECIMAL(18,2),
    `margin_percentage` DECIMAL(5,2),
    `project_completion_percentage` INT,
    `expected_completion_date` DATE,
    `number_of_active_milestones` INT,
    `number_of_delayed_milestones` INT,
    `customer_satisfaction_score` DECIMAL(3,1),
    `last_updated_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_date` (`summary_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Project KPI and financial summary dashboard';

-- ============================================================
-- BANK FINANCING REPORT (NEW TABLE)
-- Tracks banker's sanction, disbursement, and collection details per unit
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_bank_financing` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `block_id` VARCHAR(36),
    `unit_id` VARCHAR(36) NOT NULL,
    `customer_id` VARCHAR(36) NOT NULL,
    -- UNIT & APARTMENT DETAILS
    `apt_no` VARCHAR(50),
    `block_name` VARCHAR(100),
    `apartment_cost` DECIMAL(18,2) COMMENT 'Base apartment cost',
    -- BANK DETAILS
    `bank_name` VARCHAR(200),
    `banker_reference_no` VARCHAR(100),
    -- SANCTION DETAILS
    `sanctioned_amount` DECIMAL(18,2) COMMENT 'Total sanctioned loan amount',
    `sanctioned_date` DATE,
    -- DISBURSEMENT TRACKING
    `total_disbursed_amount` DECIMAL(18,2) COMMENT 'Total amount disbursed till date',
    `disbursement_status` VARCHAR(50) DEFAULT 'PENDING' COMMENT 'PENDING, PARTIAL, COMPLETED, PENDING_DOCUMENTS',
    `last_disbursement_date` DATE,
    `remaining_disbursement` DECIMAL(18,2) COMMENT 'Amount yet to be disbursed',
    -- COLLECTION FROM UNIT OWNER
    `total_collection_from_unit` DECIMAL(18,2) COMMENT 'Total amount collected from unit owner (non-bank)',
    `collection_status` VARCHAR(50) DEFAULT 'PENDING' COMMENT 'PENDING, PARTIAL, COMPLETED',
    -- FINANCIAL SUMMARY
    `total_commitment` DECIMAL(18,2) COMMENT 'Total project cost commitment',
    `outstanding_amount` DECIMAL(18,2) COMMENT 'Amount still to be collected/received',
    -- COMPLIANCE & DOCUMENTATION
    `noc_required` TINYINT(1) DEFAULT 0,
    `noc_received` TINYINT(1) DEFAULT 0,
    `noc_date` DATE,
    `documents_status` VARCHAR(50) COMMENT 'COMPLETE, PENDING, SUBMITTED_FOR_VERIFICATION',
    `documents_url` VARCHAR(500),
    -- METADATA
    `remarks` TEXT,
    `active` TINYINT(1) DEFAULT 1,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`customer_id`) REFERENCES `property_customer_profile`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_unit` (`unit_id`),
    KEY `idx_customer` (`customer_id`),
    KEY `idx_bank` (`bank_name`),
    KEY `idx_disbursement` (`disbursement_status`),
    KEY `idx_collection` (`collection_status`),
    UNIQUE KEY `unique_financing` (`tenant_id`, `unit_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Bank financing details: sanction, disbursement, and collection tracking per unit';

-- ============================================================
-- DISBURSEMENT SCHEDULE (NEW TABLE)
-- Tracks expected vs actual disbursements linked to payment milestones
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_disbursement_schedule` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `financing_id` VARCHAR(36) NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `customer_id` VARCHAR(36) NOT NULL,
    -- DISBURSEMENT DETAILS
    `disbursement_no` INT COMMENT 'Sequential disbursement number (1st, 2nd, 3rd, etc.)',
    `expected_disbursement_date` DATE,
    `actual_disbursement_date` DATE,
    `expected_disbursement_amount` DECIMAL(18,2),
    `actual_disbursement_amount` DECIMAL(18,2),
    `disbursement_percentage` DECIMAL(5,2) COMMENT 'Percentage of total sanctioned amount',
    -- MILESTONE LINKAGE
    `linked_milestone_id` VARCHAR(36),
    `milestone_stage` VARCHAR(100) COMMENT 'e.g., FOUNDATION, STRUCTURE, FINISHING, HANDOVER',
    -- BANK DOCUMENTATION
    `cheque_no` VARCHAR(50),
    `bank_reference_no` VARCHAR(100),
    `disbursement_status` VARCHAR(50) DEFAULT 'PENDING' COMMENT 'PENDING, CLEARED, FAILED, CANCELLED',
    `neft_ref_id` VARCHAR(100),
    -- VARIANCE TRACKING
    `variance_days` INT COMMENT 'Difference between expected and actual date',
    `variance_amount` DECIMAL(18,2) COMMENT 'Difference between expected and actual amount',
    `variance_reason` TEXT,
    -- METADATA
    `remarks` TEXT,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`financing_id`) REFERENCES `property_bank_financing`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`customer_id`) REFERENCES `property_customer_profile`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_financing` (`financing_id`),
    KEY `idx_unit` (`unit_id`),
    KEY `idx_status` (`disbursement_status`),
    KEY `idx_schedule_date` (`expected_disbursement_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Detailed disbursement schedule linked to project milestones';

-- ============================================================
-- PAYMENT STAGE TRACKING (NEW TABLE)
-- Maps installment stages to percentage of cost and collections
-- ============================================================
CREATE TABLE IF NOT EXISTS `property_payment_stage` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `customer_id` VARCHAR(36) NOT NULL,
    -- STAGE DEFINITION
    `stage_name` VARCHAR(100) COMMENT 'e.g., BOOKING, FOUNDATION, STRUCTURE, FINISHING, HANDOVER',
    `stage_number` INT,
    `stage_description` TEXT,
    -- PERCENTAGE & COST
    `stage_percentage` DECIMAL(5,2) COMMENT 'Percentage of total apartment cost due at this stage',
    `stage_due_amount` DECIMAL(18,2) COMMENT 'Calculated from apartment_cost * stage_percentage',
    `apartment_cost` DECIMAL(18,2),
    -- COLLECTION DETAILS
    `amount_due` DECIMAL(18,2),
    `amount_received` DECIMAL(18,2),
    `amount_pending` DECIMAL(18,2),
    `collection_status` VARCHAR(50) DEFAULT 'PENDING' COMMENT 'PENDING, PARTIAL, COMPLETED, OVERDUE',
    -- TIMELINE
    `due_date` DATE,
    `expected_collection_date` DATE,
    `actual_collection_date` DATE,
    `days_overdue` INT COMMENT 'Number of days past due date',
    -- PAYMENT DETAILS
    `payment_received_date` DATE,
    `payment_mode` VARCHAR(50) COMMENT 'CASH, CHEQUE, NEFT, RTGS, ONLINE, DD',
    `reference_no` VARCHAR(100) COMMENT 'Cheque no, transaction ID, etc.',
    -- VARIANCE
    `variance_days` INT COMMENT 'Difference from due date',
    `variance_amount` DECIMAL(18,2) COMMENT 'Shortfall or excess',
    -- METADATA
    `remarks` TEXT,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`customer_id`) REFERENCES `property_customer_profile`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_project` (`project_id`),
    KEY `idx_unit` (`unit_id`),
    KEY `idx_customer` (`customer_id`),
    KEY `idx_stage` (`stage_number`),
    KEY `idx_due_date` (`due_date`),
    KEY `idx_collection_status` (`collection_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Payment stage tracking: installment schedule with percentage, due amounts, and collections';

-- ============================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================
CREATE INDEX idx_customer_search ON `property_customer_profile`(`tenant_id`, `customer_status`, `created_at` DESC);
CREATE INDEX idx_payment_summary ON `property_payment_receipt`(`tenant_id`, `payment_date` DESC);
CREATE INDEX idx_milestone_progress ON `property_project_milestone`(`project_id`, `completion_status`);
CREATE INDEX idx_activity_timeline ON `property_project_activity`(`project_id`, `activity_date` DESC);

SET FOREIGN_KEY_CHECKS = 1;
