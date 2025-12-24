-- ============================================================
-- VYOMTECH ERP - COMPLETE DATABASE SCHEMA
-- Generated from Prisma Schema
-- All tables with ULID/UUID support, multi-tenant isolation, and soft deletes
-- ============================================================

-- ============================================================
-- CORE MULTI-TENANT TABLES
-- ============================================================

CREATE TABLE IF NOT EXISTS `tenant` (
  `id` CHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `domain` VARCHAR(255) NOT NULL,
  `status` VARCHAR(50) NOT NULL DEFAULT 'active',
  `max_users` INT NOT NULL DEFAULT 100,
  `max_concurrent_calls` INT NOT NULL DEFAULT 50,
  `ai_budget_monthly` DECIMAL(15,2) NOT NULL DEFAULT 1000.00,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  UNIQUE INDEX `tenant_domain_key`(`domain`),
  INDEX `tenant_domain_idx`(`domain`),
  INDEX `tenant_status_idx`(`status`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `user` (
  `id` CHAR(36) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `role` VARCHAR(50) NOT NULL DEFAULT 'user',
  `tenant_id` VARCHAR(36) NOT NULL,
  `current_tenant_id` VARCHAR(36),
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  UNIQUE INDEX `user_email_key`(`email`),
  INDEX `user_email_idx`(`email`),
  INDEX `user_tenant_id_idx`(`tenant_id`),
  INDEX `user_current_tenant_id_idx`(`current_tenant_id`),
  INDEX `user_role_idx`(`role`),
  PRIMARY KEY (`id`),
  CONSTRAINT `user_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `user_current_tenant_id_fk` FOREIGN KEY (`current_tenant_id`) REFERENCES `tenant` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `team` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `description` TEXT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `team_tenant_id_idx`(`tenant_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `team_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `system_config` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `config_key` VARCHAR(255) NOT NULL,
  `config_value` LONGTEXT NOT NULL,
  `data_type` VARCHAR(50) NOT NULL,
  `description` TEXT,
  `is_global` BOOLEAN NOT NULL DEFAULT false,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  UNIQUE INDEX `system_config_tenant_id_config_key_key`(`tenant_id`, `config_key`),
  INDEX `system_config_tenant_id_idx`(`tenant_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `system_config_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `audit_log` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `user_id` CHAR(36),
  `action` VARCHAR(50) NOT NULL,
  `resource_type` VARCHAR(100) NOT NULL,
  `resource_id` VARCHAR(255),
  `changes` JSON,
  `ip_address` VARCHAR(45),
  `user_agent` VARCHAR(255),
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  INDEX `audit_log_tenant_id_idx`(`tenant_id`),
  INDEX `audit_log_user_id_idx`(`user_id`),
  INDEX `audit_log_action_idx`(`action`),
  INDEX `audit_log_resource_type_idx`(`resource_type`),
  PRIMARY KEY (`id`),
  CONSTRAINT `audit_log_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `audit_log_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- AUTHENTICATION & SECURITY
-- ============================================================

CREATE TABLE IF NOT EXISTS `password_reset_token` (
  `id` CHAR(36) NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `token` VARCHAR(255) NOT NULL,
  `expires_at` DATETIME(3) NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  UNIQUE INDEX `password_reset_token_token_key`(`token`),
  INDEX `password_reset_token_user_id_idx`(`user_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `password_reset_token_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `auth_token` (
  `id` CHAR(36) NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `token` VARCHAR(255) NOT NULL,
  `expires_at` DATETIME(3) NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  UNIQUE INDEX `auth_token_token_key`(`token`),
  INDEX `auth_token_user_id_idx`(`user_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `auth_token_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- SALES & CRM
-- ============================================================

CREATE TABLE IF NOT EXISTS `sales_lead` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `owner_id` CHAR(36) NOT NULL,
  `lead_name` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255),
  `phone` VARCHAR(20),
  `source_id` VARCHAR(36),
  `sub_source_id` VARCHAR(36),
  `status` VARCHAR(50) NOT NULL DEFAULT 'new',
  `budget` DECIMAL(15,2),
  `property_type` VARCHAR(100),
  `location` VARCHAR(255),
  `notes` LONGTEXT,
  `next_follow_up` DATETIME,
  `last_activity` DATETIME,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME,

  INDEX `sales_lead_tenant_id_idx`(`tenant_id`),
  INDEX `sales_lead_owner_id_idx`(`owner_id`),
  INDEX `sales_lead_source_id_idx`(`source_id`),
  INDEX `sales_lead_status_idx`(`status`),
  INDEX `sales_lead_created_at_idx`(`created_at`),
  PRIMARY KEY (`id`),
  CONSTRAINT `sales_lead_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `sales_lead_owner_id_fk` FOREIGN KEY (`owner_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `booking` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `lead_id` CHAR(26) NOT NULL,
  `booking_date` DATETIME NOT NULL,
  `amount` DECIMAL(15,2) NOT NULL,
  `status` VARCHAR(50) NOT NULL DEFAULT 'confirmed',
  `payment_status` VARCHAR(50) NOT NULL DEFAULT 'pending',
  `notes` LONGTEXT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME,

  INDEX `booking_tenant_id_idx`(`tenant_id`),
  INDEX `booking_lead_id_idx`(`lead_id`),
  INDEX `booking_status_idx`(`status`),
  PRIMARY KEY (`id`),
  CONSTRAINT `booking_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `booking_lead_id_fk` FOREIGN KEY (`lead_id`) REFERENCES `sales_lead` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `lead_activity` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `lead_id` CHAR(26) NOT NULL,
  `activity_type` VARCHAR(50) NOT NULL,
  `description` LONGTEXT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  INDEX `lead_activity_tenant_id_idx`(`tenant_id`),
  INDEX `lead_activity_lead_id_idx`(`lead_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `lead_activity_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `lead_activity_lead_id_fk` FOREIGN KEY (`lead_id`) REFERENCES `sales_lead` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `call_log` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `lead_id` CHAR(26),
  `duration` INT NOT NULL DEFAULT 0,
  `call_type` VARCHAR(50) NOT NULL,
  `recording_url` VARCHAR(500),
  `transcription` LONGTEXT,
  `ai_summary` LONGTEXT,
  `sentiment` VARCHAR(50),
  `status` VARCHAR(50) NOT NULL DEFAULT 'completed',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  INDEX `call_log_tenant_id_idx`(`tenant_id`),
  INDEX `call_log_lead_id_idx`(`lead_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `call_log_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `call_log_lead_id_fk` FOREIGN KEY (`lead_id`) REFERENCES `sales_lead` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- REAL ESTATE
-- ============================================================

CREATE TABLE IF NOT EXISTS `property` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `property_code` VARCHAR(100) NOT NULL,
  `property_name` VARCHAR(255) NOT NULL,
  `property_type` VARCHAR(100),
  `location` VARCHAR(255),
  `city` VARCHAR(100),
  `state` VARCHAR(100),
  `country` VARCHAR(100),
  `zip_code` VARCHAR(20),
  `area_sqft` DECIMAL(12,2),
  `area_sqm` DECIMAL(12,2),
  `bedrooms` INT,
  `bathrooms` INT,
  `parking_slots` INT,
  `price` DECIMAL(15,2),
  `currency` VARCHAR(10),
  `status` VARCHAR(50) NOT NULL DEFAULT 'available',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `property_tenant_id_idx`(`tenant_id`),
  INDEX `property_status_idx`(`status`),
  UNIQUE INDEX `property_code_idx`(`property_code`),
  PRIMARY KEY (`id`),
  CONSTRAINT `property_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `site_visit` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `lead_id` CHAR(26) NOT NULL,
  `property_id` CHAR(26),
  `visit_date` DATETIME NOT NULL,
  `feedback` LONGTEXT,
  `status` VARCHAR(50),
  `created_by` CHAR(36),
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `site_visit_tenant_id_idx`(`tenant_id`),
  INDEX `site_visit_lead_id_idx`(`lead_id`),
  INDEX `site_visit_property_id_idx`(`property_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `site_visit_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `site_visit_lead_id_fk` FOREIGN KEY (`lead_id`) REFERENCES `sales_lead` (`id`),
  CONSTRAINT `site_visit_property_id_fk` FOREIGN KEY (`property_id`) REFERENCES `property` (`id`) ON DELETE SET NULL,
  CONSTRAINT `site_visit_created_by_fk` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `possession` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `property_id` CHAR(26) NOT NULL,
  `lead_id` CHAR(26) NOT NULL,
  `handover_date` DATETIME,
  `possession_date` DATETIME,
  `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
  `documents` JSON,
  `created_by` CHAR(36),
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `possession_tenant_id_idx`(`tenant_id`),
  INDEX `possession_property_id_idx`(`property_id`),
  INDEX `possession_lead_id_idx`(`lead_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `possession_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `possession_property_id_fk` FOREIGN KEY (`property_id`) REFERENCES `property` (`id`),
  CONSTRAINT `possession_lead_id_fk` FOREIGN KEY (`lead_id`) REFERENCES `sales_lead` (`id`),
  CONSTRAINT `possession_created_by_fk` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `title_clearance` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `property_id` CHAR(26) NOT NULL,
  `lead_id` CHAR(26) NOT NULL,
  `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
  `clearance_date` DATETIME,
  `remarks` LONGTEXT,
  `created_by` CHAR(36),
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `title_clearance_tenant_id_idx`(`tenant_id`),
  INDEX `title_clearance_property_id_idx`(`property_id`),
  INDEX `title_clearance_lead_id_idx`(`lead_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `title_clearance_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `title_clearance_property_id_fk` FOREIGN KEY (`property_id`) REFERENCES `property` (`id`),
  CONSTRAINT `title_clearance_lead_id_fk` FOREIGN KEY (`lead_id`) REFERENCES `sales_lead` (`id`),
  CONSTRAINT `title_clearance_created_by_fk` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- DOCUMENTS
-- ============================================================

CREATE TABLE IF NOT EXISTS `document_category` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `category_name` VARCHAR(255) NOT NULL,
  `category_code` VARCHAR(50) NOT NULL,
  `description` TEXT,
  `display_order` INT DEFAULT 0,
  `is_active` BOOLEAN NOT NULL DEFAULT true,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `document_category_tenant_id_idx`(`tenant_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `document_category_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `document_type` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `category_id` CHAR(26) NOT NULL,
  `type_name` VARCHAR(255) NOT NULL,
  `type_code` VARCHAR(50) NOT NULL,
  `description` TEXT,
  `is_mandatory` BOOLEAN NOT NULL DEFAULT false,
  `is_identity_proof` BOOLEAN NOT NULL DEFAULT false,
  `file_formats` VARCHAR(255),
  `max_file_size` BIGINT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `document_type_tenant_id_idx`(`tenant_id`),
  INDEX `document_type_category_id_idx`(`category_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `document_type_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `document_type_category_id_fk` FOREIGN KEY (`category_id`) REFERENCES `document_category` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `document` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `document_type_id` CHAR(26) NOT NULL,
  `entity_type` VARCHAR(100) NOT NULL,
  `entity_id` VARCHAR(255) NOT NULL,
  `document_name` VARCHAR(255) NOT NULL,
  `document_url` VARCHAR(500) NOT NULL,
  `file_name` VARCHAR(255),
  `file_size` BIGINT,
  `file_format` VARCHAR(50),
  `s3_bucket` VARCHAR(255),
  `s3_key` VARCHAR(500),
  `verifier_id` CHAR(36),
  `is_primary` BOOLEAN NOT NULL DEFAULT false,
  `metadata` JSON,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `document_tenant_id_idx`(`tenant_id`),
  INDEX `document_type_id_idx`(`document_type_id`),
  INDEX `document_entity_type_entity_id_idx`(`entity_type`, `entity_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `document_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `document_document_type_id_fk` FOREIGN KEY (`document_type_id`) REFERENCES `document_type` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- FINANCE & ACCOUNTING
-- ============================================================

CREATE TABLE IF NOT EXISTS `loan` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `loan_type` VARCHAR(100) NOT NULL,
  `loan_amount` DECIMAL(15,2) NOT NULL,
  `interest_rate` DECIMAL(5,2),
  `tenure` INT,
  `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
  `loan_officer_id` CHAR(36),
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `loan_tenant_id_idx`(`tenant_id`),
  INDEX `loan_status_idx`(`status`),
  INDEX `loan_loan_officer_id_idx`(`loan_officer_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `loan_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `loan_loan_officer_id_fk` FOREIGN KEY (`loan_officer_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `bank_financing` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `loan_id` CHAR(26) NOT NULL,
  `bank_name` VARCHAR(255) NOT NULL,
  `bank_code` VARCHAR(50) NOT NULL,
  `sanction_amount` DECIMAL(15,2) NOT NULL,
  `disbursed_amount` DECIMAL(15,2) NOT NULL DEFAULT 0,
  `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `bank_financing_tenant_id_idx`(`tenant_id`),
  INDEX `bank_financing_loan_id_idx`(`loan_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `bank_financing_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `bank_financing_loan_id_fk` FOREIGN KEY (`loan_id`) REFERENCES `loan` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `payment_transaction` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `amount` DECIMAL(15,2) NOT NULL,
  `payment_method` VARCHAR(50),
  `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
  `transaction_date` DATETIME,
  `reference_number` VARCHAR(255),
  `description` TEXT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `payment_transaction_tenant_id_idx`(`tenant_id`),
  INDEX `payment_transaction_status_idx`(`status`),
  PRIMARY KEY (`id`),
  CONSTRAINT `payment_transaction_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `account` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `account_code` VARCHAR(50) NOT NULL,
  `account_name` VARCHAR(255) NOT NULL,
  `account_type` VARCHAR(50) NOT NULL,
  `sub_type` VARCHAR(50),
  `balance` DECIMAL(15,2) NOT NULL DEFAULT 0,
  `status` VARCHAR(50) NOT NULL DEFAULT 'active',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `account_tenant_id_idx`(`tenant_id`),
  INDEX `account_type_idx`(`account_type`),
  UNIQUE INDEX `account_code_idx`(`account_code`),
  PRIMARY KEY (`id`),
  CONSTRAINT `account_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `voucher` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `voucher_number` VARCHAR(50) NOT NULL,
  `voucher_date` DATETIME NOT NULL,
  `voucher_type` VARCHAR(50) NOT NULL,
  `amount` DECIMAL(15,2) NOT NULL,
  `description` TEXT,
  `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `voucher_tenant_id_idx`(`tenant_id`),
  INDEX `voucher_status_idx`(`status`),
  UNIQUE INDEX `voucher_number_idx`(`voucher_number`),
  PRIMARY KEY (`id`),
  CONSTRAINT `voucher_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- HR & PAYROLL
-- ============================================================

CREATE TABLE IF NOT EXISTS `employee` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `first_name` VARCHAR(100) NOT NULL,
  `last_name` VARCHAR(100) NOT NULL,
  `email` VARCHAR(255),
  `phone` VARCHAR(20),
  `designation` VARCHAR(100),
  `department` VARCHAR(100),
  `status` VARCHAR(50) NOT NULL DEFAULT 'active',
  `joining_date` DATETIME,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `employee_tenant_id_idx`(`tenant_id`),
  INDEX `employee_user_id_idx`(`user_id`),
  INDEX `employee_status_idx`(`status`),
  UNIQUE INDEX `employee_user_id_key`(`user_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `employee_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `employee_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `leave` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `employee_id` CHAR(26) NOT NULL,
  `leave_type` VARCHAR(50) NOT NULL,
  `from_date` DATETIME NOT NULL,
  `to_date` DATETIME NOT NULL,
  `days` DECIMAL(5,2) NOT NULL,
  `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
  `reason` TEXT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `leave_tenant_id_idx`(`tenant_id`),
  INDEX `leave_employee_id_idx`(`employee_id`),
  INDEX `leave_status_idx`(`status`),
  PRIMARY KEY (`id`),
  CONSTRAINT `leave_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `leave_employee_id_fk` FOREIGN KEY (`employee_id`) REFERENCES `employee` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `salary` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `employee_id` CHAR(26) NOT NULL,
  `month` VARCHAR(7) NOT NULL,
  `amount` DECIMAL(15,2) NOT NULL,
  `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `salary_tenant_id_idx`(`tenant_id`),
  INDEX `salary_employee_id_idx`(`employee_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `salary_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `salary_employee_id_fk` FOREIGN KEY (`employee_id`) REFERENCES `employee` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `attendance` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `employee_id` CHAR(26) NOT NULL,
  `attendance_date` DATETIME NOT NULL,
  `status` VARCHAR(50) NOT NULL,
  `check_in_time` DATETIME,
  `check_out_time` DATETIME,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `attendance_tenant_id_idx`(`tenant_id`),
  INDEX `attendance_employee_id_idx`(`employee_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `attendance_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `attendance_employee_id_fk` FOREIGN KEY (`employee_id`) REFERENCES `employee` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `payroll` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `employee_id` CHAR(26) NOT NULL,
  `payroll_month` VARCHAR(7) NOT NULL,
  `salary` DECIMAL(15,2),
  `allowances` DECIMAL(15,2) DEFAULT 0,
  `deductions` DECIMAL(15,2) DEFAULT 0,
  `net_amount` DECIMAL(15,2),
  `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `payroll_tenant_id_idx`(`tenant_id`),
  INDEX `payroll_employee_id_idx`(`employee_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `payroll_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `payroll_employee_id_fk` FOREIGN KEY (`employee_id`) REFERENCES `employee` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PARTNERS & BROKERS
-- ============================================================

CREATE TABLE IF NOT EXISTS `partner` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `user_id` CHAR(36),
  `name` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255),
  `phone` VARCHAR(20),
  `status` VARCHAR(50) NOT NULL DEFAULT 'active',
  `commission_percent` DECIMAL(5,2),
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `partner_tenant_id_idx`(`tenant_id`),
  INDEX `partner_status_idx`(`status`),
  PRIMARY KEY (`id`),
  CONSTRAINT `partner_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `partner_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `broker` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `broker_code` VARCHAR(100) NOT NULL,
  `broker_name` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255),
  `phone` VARCHAR(20),
  `status` VARCHAR(50) NOT NULL DEFAULT 'active',
  `commission_percent` DECIMAL(5,2),
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `broker_tenant_id_idx`(`tenant_id`),
  INDEX `broker_status_idx`(`status`),
  UNIQUE INDEX `broker_code_idx`(`broker_code`),
  PRIMARY KEY (`id`),
  CONSTRAINT `broker_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `joint_applicant` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `lead_id` CHAR(26) NOT NULL,
  `applicant_name` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255),
  `phone` VARCHAR(20),
  `relationship` VARCHAR(50),
  `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `joint_applicant_tenant_id_idx`(`tenant_id`),
  INDEX `joint_applicant_lead_id_idx`(`lead_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `joint_applicant_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `joint_applicant_lead_id_fk` FOREIGN KEY (`lead_id`) REFERENCES `sales_lead` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ADDITIONAL ENTITIES
-- ============================================================

CREATE TABLE IF NOT EXISTS `lead` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `lead_code` VARCHAR(50) NOT NULL,
  `first_name` VARCHAR(100),
  `last_name` VARCHAR(100),
  `email` VARCHAR(255),
  `phone` VARCHAR(20),
  `company_name` VARCHAR(255),
  `industry` VARCHAR(100),
  `status` VARCHAR(50) NOT NULL DEFAULT 'new',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `lead_tenant_id_idx`(`tenant_id`),
  INDEX `lead_status_idx`(`status`),
  UNIQUE INDEX `lead_code_idx`(`lead_code`),
  PRIMARY KEY (`id`),
  CONSTRAINT `lead_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `applicant` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `lead_id` CHAR(26),
  `applicant_name` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255),
  `phone` VARCHAR(20),
  `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `applicant_tenant_id_idx`(`tenant_id`),
  INDEX `applicant_lead_id_idx`(`lead_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `applicant_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `applicant_lead_id_fk` FOREIGN KEY (`lead_id`) REFERENCES `lead` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `project` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `project_name` VARCHAR(255) NOT NULL,
  `description` TEXT,
  `status` VARCHAR(50) NOT NULL DEFAULT 'active',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `project_tenant_id_idx`(`tenant_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `project_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `integration` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `integration_name` VARCHAR(255) NOT NULL,
  `integration_type` VARCHAR(100),
  `is_active` BOOLEAN NOT NULL DEFAULT true,
  `config` JSON,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL ON UPDATE CURRENT_TIMESTAMP(3),

  INDEX `integration_tenant_id_idx`(`tenant_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `integration_tenant_id_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- MIGRATION HISTORY (managed by Prisma)
-- ============================================================

CREATE TABLE IF NOT EXISTS `_prisma_migrations` (
  `id`                    VARCHAR(36) PRIMARY KEY,
  `checksum`              VARCHAR(64) NOT NULL,
  `finished_at`           DATETIME,
  `migration_name`        VARCHAR(255) NOT NULL,
  `logs`                  TEXT,
  `rolled_back_at`        DATETIME,
  `started_at`            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `applied_steps_count`   INT UNSIGNED NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- END OF SCHEMA CREATION
-- ============================================================
