-- ============================================================
-- MIGRATION 007: SALES MODULE
-- Date: December 3, 2025
-- Purpose: Create sales leads, customers, quotations, orders, and invoices
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- SALES LEAD TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `sales_lead` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `lead_code` VARCHAR(50) UNIQUE NOT NULL,
    `first_name` VARCHAR(100),
    `last_name` VARCHAR(100),
    `email` VARCHAR(255),
    `phone` VARCHAR(20),
    `company_name` VARCHAR(255),
    `industry` VARCHAR(100),
    `status` VARCHAR(50) DEFAULT 'new',
    `probability` DECIMAL(5, 2) DEFAULT 0,
    `source` VARCHAR(50),
    `campaign_id` VARCHAR(36),
    `assigned_to` VARCHAR(36),
    `assigned_date` TIMESTAMP,
    `converted_to_customer` BOOLEAN DEFAULT FALSE,
    `customer_id` VARCHAR(36),
    `next_action_date` DATE,
    `next_action_notes` TEXT,
    `created_by` VARCHAR(36),
    `gl_customer_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`gl_customer_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_code` (`lead_code`),
    KEY `idx_status` (`status`),
    KEY `idx_assigned` (`assigned_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- SALES CUSTOMER TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `sales_customer` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `customer_code` VARCHAR(50) UNIQUE NOT NULL,
    `customer_name` VARCHAR(255) NOT NULL,
    `business_name` VARCHAR(255),
    `business_type` VARCHAR(50),
    `industry` VARCHAR(100),
    `primary_contact_name` VARCHAR(100),
    `primary_email` VARCHAR(255),
    `primary_phone` VARCHAR(20),
    `billing_address` TEXT,
    `billing_city` VARCHAR(100),
    `billing_state` VARCHAR(100),
    `billing_country` VARCHAR(100),
    `billing_zip` VARCHAR(20),
    `shipping_address` TEXT,
    `shipping_city` VARCHAR(100),
    `shipping_state` VARCHAR(100),
    `shipping_country` VARCHAR(100),
    `shipping_zip` VARCHAR(20),
    `pan_number` VARCHAR(50),
    `gst_number` VARCHAR(50),
    `credit_limit` DECIMAL(18, 2),
    `credit_days` INT DEFAULT 0,
    `payment_terms` VARCHAR(100),
    `customer_category` VARCHAR(50),
    `status` VARCHAR(50) DEFAULT 'active',
    `current_balance` DECIMAL(18, 2) DEFAULT 0,
    `created_by` VARCHAR(36),
    `gl_receivable_account_id` VARCHAR(36),
    `gl_advance_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`gl_receivable_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`gl_advance_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_code` (`customer_code`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- SALES QUOTATION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `sales_quotation` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `quotation_number` VARCHAR(50) UNIQUE NOT NULL,
    `customer_id` VARCHAR(36) NOT NULL,
    `quotation_date` DATE NOT NULL,
    `valid_until` DATE,
    `subtotal_amount` DECIMAL(18, 2),
    `discount_amount` DECIMAL(18, 2) DEFAULT 0,
    `tax_amount` DECIMAL(18, 2) DEFAULT 0,
    `total_amount` DECIMAL(18, 2),
    `status` VARCHAR(50) DEFAULT 'draft',
    `converted_to_order` BOOLEAN DEFAULT FALSE,
    `sales_order_id` VARCHAR(36),
    `notes` TEXT,
    `terms_and_conditions` TEXT,
    `created_by` VARCHAR(36),
    `gl_revenue_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`customer_id`) REFERENCES `sales_customer`(`id`),
    FOREIGN KEY (`gl_revenue_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_number` (`quotation_number`),
    KEY `idx_customer` (`customer_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- SALES QUOTATION ITEM TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `sales_quotation_item` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `quotation_id` VARCHAR(36) NOT NULL,
    `line_number` INT NOT NULL,
    `product_code` VARCHAR(100),
    `description` TEXT,
    `quantity` DECIMAL(15, 2) NOT NULL,
    `unit` VARCHAR(50),
    `unit_price` DECIMAL(18, 2) NOT NULL,
    `line_total` DECIMAL(18, 2),
    `tax_rate` DECIMAL(5, 2) DEFAULT 0,
    `tax_amount` DECIMAL(18, 2) DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`quotation_id`) REFERENCES `sales_quotation`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_quotation` (`quotation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- SALES ORDER TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `sales_order` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `order_number` VARCHAR(50) UNIQUE NOT NULL,
    `customer_id` VARCHAR(36) NOT NULL,
    `quotation_id` VARCHAR(36),
    `order_date` DATE NOT NULL,
    `delivery_date` DATE,
    `subtotal_amount` DECIMAL(18, 2),
    `discount_amount` DECIMAL(18, 2) DEFAULT 0,
    `tax_amount` DECIMAL(18, 2) DEFAULT 0,
    `shipping_amount` DECIMAL(18, 2) DEFAULT 0,
    `total_amount` DECIMAL(18, 2),
    `payment_terms` VARCHAR(100),
    `delivery_location` VARCHAR(255),
    `shipping_address` TEXT,
    `order_status` VARCHAR(50) DEFAULT 'pending',
    `payment_status` VARCHAR(50) DEFAULT 'unpaid',
    `created_by` VARCHAR(36),
    `gl_revenue_account_id` VARCHAR(36),
    `gl_receivable_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`customer_id`) REFERENCES `sales_customer`(`id`),
    FOREIGN KEY (`gl_revenue_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`gl_receivable_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`quotation_id`) REFERENCES `sales_quotation`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_number` (`order_number`),
    KEY `idx_customer` (`customer_id`),
    KEY `idx_status` (`order_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- SALES ORDER ITEM TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `sales_order_item` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `order_id` VARCHAR(36) NOT NULL,
    `line_number` INT NOT NULL,
    `product_code` VARCHAR(100),
    `description` TEXT,
    `quantity` DECIMAL(15, 2) NOT NULL,
    `unit` VARCHAR(50),
    `unit_price` DECIMAL(18, 2) NOT NULL,
    `line_total` DECIMAL(18, 2),
    `tax_rate` DECIMAL(5, 2) DEFAULT 0,
    `tax_amount` DECIMAL(18, 2) DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`order_id`) REFERENCES `sales_order`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- SALES INVOICE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `sales_invoice` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `invoice_number` VARCHAR(50) UNIQUE NOT NULL,
    `customer_id` VARCHAR(36) NOT NULL,
    `order_id` VARCHAR(36),
    `invoice_date` DATE NOT NULL,
    `due_date` DATE,
    `subtotal_amount` DECIMAL(18, 2),
    `discount_amount` DECIMAL(18, 2) DEFAULT 0,
    `tax_amount` DECIMAL(18, 2) DEFAULT 0,
    `shipping_amount` DECIMAL(18, 2) DEFAULT 0,
    `total_amount` DECIMAL(18, 2),
    `amount_paid` DECIMAL(18, 2) DEFAULT 0,
    `balance_due` DECIMAL(18, 2),
    `invoice_status` VARCHAR(50) DEFAULT 'draft',
    `payment_status` VARCHAR(50) DEFAULT 'unpaid',
    `created_by` VARCHAR(36),
    `gl_revenue_account_id` VARCHAR(36),
    `gl_receivable_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`customer_id`) REFERENCES `sales_customer`(`id`),
    FOREIGN KEY (`order_id`) REFERENCES `sales_order`(`id`),
    FOREIGN KEY (`gl_revenue_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`gl_receivable_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_number` (`invoice_number`),
    KEY `idx_customer` (`customer_id`),
    KEY `idx_status` (`invoice_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
