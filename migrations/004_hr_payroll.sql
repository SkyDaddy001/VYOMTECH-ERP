-- ============================================================
-- MIGRATION 004: HR & PAYROLL MODULE
-- Date: December 3, 2025
-- Purpose: Create HR, employee, attendance, and payroll tables
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- EMPLOYEE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `employee` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100),
    `email` VARCHAR(255) UNIQUE,
    `phone` VARCHAR(20),
    `date_of_birth` DATE,
    `gender` VARCHAR(20),
    `nationality` VARCHAR(100),
    `address` TEXT,
    `city` VARCHAR(100),
    `state` VARCHAR(100),
    `country` VARCHAR(100),
    `postal_code` VARCHAR(20),
    `employee_id` VARCHAR(50) UNIQUE NOT NULL,
    `designation` VARCHAR(100),
    `department` VARCHAR(100),
    `report_to` VARCHAR(36),
    `employment_type` VARCHAR(50),
    `joining_date` DATE,
    `exit_date` DATE,
    `status` VARCHAR(50) DEFAULT 'active',
    `bank_account_number` VARCHAR(50),
    `bank_ifsc_code` VARCHAR(20),
    `bank_name` VARCHAR(100),
    `account_holder_name` VARCHAR(100),
    `base_salary` DECIMAL(18, 2),
    `gl_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`gl_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_department` (`department`),
    KEY `idx_status` (`status`),
    KEY `idx_employee_id` (`employee_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ATTENDANCE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `attendance` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` VARCHAR(36) NOT NULL,
    `attendance_date` DATE NOT NULL,
    `check_in_time` TIMESTAMP,
    `check_out_time` TIMESTAMP,
    `working_hours` DECIMAL(5, 2),
    `status` VARCHAR(50),
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_employee` (`employee_id`),
    KEY `idx_date` (`attendance_date`),
    UNIQUE KEY `unique_attendance` (`employee_id`, `attendance_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- LEAVE TYPE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `leave_type` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `leave_type_name` VARCHAR(100) NOT NULL,
    `description` TEXT,
    `annual_entitlement` INT DEFAULT 0,
    `is_paid` BOOLEAN DEFAULT TRUE,
    `carry_forward_limit` INT DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- LEAVE REQUEST TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `leave_request` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` VARCHAR(36) NOT NULL,
    `leave_type_id` VARCHAR(36) NOT NULL,
    `from_date` DATE NOT NULL,
    `to_date` DATE NOT NULL,
    `number_of_days` INT,
    `reason` TEXT,
    `status` VARCHAR(50) DEFAULT 'pending',
    `approved_by` VARCHAR(36),
    `approval_date` TIMESTAMP,
    `rejection_reason` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`leave_type_id`) REFERENCES `leave_type`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_employee` (`employee_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PAYROLL RECORD TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `payroll_record` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` VARCHAR(36) NOT NULL,
    `payroll_month` DATE NOT NULL,
    `payroll_status` VARCHAR(50) DEFAULT 'draft',
    `basic_salary` DECIMAL(18, 2),
    `dearness_allowance` DECIMAL(18, 2) DEFAULT 0,
    `house_rent_allowance` DECIMAL(18, 2) DEFAULT 0,
    `special_allowance` DECIMAL(18, 2) DEFAULT 0,
    `conveyance_allowance` DECIMAL(18, 2) DEFAULT 0,
    `medical_allowance` DECIMAL(18, 2) DEFAULT 0,
    `other_allowances` DECIMAL(18, 2) DEFAULT 0,
    `total_earnings` DECIMAL(18, 2),
    `epf_deduction` DECIMAL(18, 2) DEFAULT 0,
    `esi_deduction` DECIMAL(18, 2) DEFAULT 0,
    `professional_tax` DECIMAL(18, 2) DEFAULT 0,
    `income_tax` DECIMAL(18, 2) DEFAULT 0,
    `loan_deduction` DECIMAL(18, 2) DEFAULT 0,
    `advance_deduction` DECIMAL(18, 2) DEFAULT 0,
    `other_deductions` DECIMAL(18, 2) DEFAULT 0,
    `total_deductions` DECIMAL(18, 2),
    `net_salary` DECIMAL(18, 2),
    `working_days` INT DEFAULT 0,
    `leave_days` INT DEFAULT 0,
    `salary_expense_account_id` VARCHAR(36),
    `epf_expense_account_id` VARCHAR(36),
    `esi_expense_account_id` VARCHAR(36),
    `payable_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`salary_expense_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`epf_expense_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`esi_expense_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`payable_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_employee` (`employee_id`),
    KEY `idx_month` (`payroll_month`),
    UNIQUE KEY `unique_payroll` (`employee_id`, `payroll_month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
