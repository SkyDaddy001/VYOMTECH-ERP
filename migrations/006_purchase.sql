-- ============================================================
-- MIGRATION 006: PURCHASE MODULE
-- Date: December 3, 2025
-- Purpose: Create vendor, purchase requisition, PO, GRN, and receiving tables
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- VENDOR TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `vendor` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `vendor_code` VARCHAR(50) UNIQUE NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255),
    `phone` VARCHAR(20),
    `address` TEXT,
    `city` VARCHAR(100),
    `state` VARCHAR(100),
    `country` VARCHAR(100),
    `postal_code` VARCHAR(20),
    `tax_id` VARCHAR(50),
    `payment_terms` VARCHAR(100),
    `vendor_type` VARCHAR(50),
    `rating` DECIMAL(3, 1),
    `is_active` BOOLEAN DEFAULT TRUE,
    `is_blocked` BOOLEAN DEFAULT FALSE,
    `status` VARCHAR(50) DEFAULT 'active',
    `created_by` VARCHAR(36),
    `gl_vendor_payable_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`gl_vendor_payable_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_code` (`vendor_code`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- VENDOR CONTACT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `vendor_contact` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `vendor_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `title` VARCHAR(100),
    `phone` VARCHAR(20),
    `email` VARCHAR(255),
    `is_primary` BOOLEAN DEFAULT FALSE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`vendor_id`) REFERENCES `vendor`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_vendor` (`vendor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- VENDOR ADDRESS TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `vendor_address` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `vendor_id` VARCHAR(36) NOT NULL,
    `type` VARCHAR(50),
    `line1` VARCHAR(255),
    `line2` VARCHAR(255),
    `city` VARCHAR(100),
    `state` VARCHAR(100),
    `country` VARCHAR(100),
    `post_code` VARCHAR(20),
    `is_primary` BOOLEAN DEFAULT FALSE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`vendor_id`) REFERENCES `vendor`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_vendor` (`vendor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PURCHASE REQUISITION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `purchase_requisition` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `requisition_number` VARCHAR(50) UNIQUE NOT NULL,
    `requester_id` VARCHAR(36),
    `department` VARCHAR(100),
    `request_date` DATE NOT NULL,
    `required_by_date` DATE,
    `purpose` TEXT,
    `status` VARCHAR(50) DEFAULT 'draft',
    `approved_by` VARCHAR(36),
    `approved_at` TIMESTAMP,
    `rejection_reason` TEXT,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_number` (`requisition_number`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PURCHASE ORDER TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `purchase_order` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `po_number` VARCHAR(50) UNIQUE NOT NULL,
    `vendor_id` VARCHAR(36) NOT NULL,
    `requisition_id` VARCHAR(36),
    `po_date` DATE NOT NULL,
    `delivery_date` DATE,
    `total_amount` DECIMAL(18, 2),
    `tax_amount` DECIMAL(18, 2) DEFAULT 0,
    `shipping_amount` DECIMAL(18, 2) DEFAULT 0,
    `discount_amount` DECIMAL(18, 2) DEFAULT 0,
    `net_amount` DECIMAL(18, 2),
    `payment_terms` VARCHAR(100),
    `delivery_location` VARCHAR(255),
    `special_instructions` TEXT,
    `status` VARCHAR(50) DEFAULT 'draft',
    `sent_to_vendor_at` TIMESTAMP,
    `acknowledged_at` TIMESTAMP,
    `created_by` VARCHAR(36),
    `gl_inventory_account_id` VARCHAR(36),
    `gl_expense_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`vendor_id`) REFERENCES `vendor`(`id`),
    FOREIGN KEY (`gl_inventory_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`gl_expense_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`requisition_id`) REFERENCES `purchase_requisition`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_number` (`po_number`),
    KEY `idx_vendor` (`vendor_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PO LINE ITEM TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `po_line_item` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `po_id` VARCHAR(36) NOT NULL,
    `line_number` INT NOT NULL,
    `product_code` VARCHAR(100),
    `description` TEXT,
    `quantity` DECIMAL(15, 2) NOT NULL,
    `unit` VARCHAR(50),
    `unit_price` DECIMAL(18, 2) NOT NULL,
    `line_total` DECIMAL(18, 2),
    `hsn_code` VARCHAR(50),
    `tax_rate` DECIMAL(5, 2) DEFAULT 0,
    `tax_amount` DECIMAL(18, 2) DEFAULT 0,
    `specification` TEXT,
    `gl_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`po_id`) REFERENCES `purchase_order`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`gl_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_po` (`po_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- GOODS RECEIPT NOTE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `goods_receipt` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `grn_number` VARCHAR(50) UNIQUE NOT NULL,
    `po_id` VARCHAR(36) NOT NULL,
    `receipt_date` DATE NOT NULL,
    `received_by` VARCHAR(36),
    `total_quantity_received` DECIMAL(15, 2),
    `total_quantity_accepted` DECIMAL(15, 2),
    `total_quantity_rejected` DECIMAL(15, 2),
    `receipt_status` VARCHAR(50) DEFAULT 'pending',
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`po_id`) REFERENCES `purchase_order`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_grn` (`grn_number`),
    KEY `idx_po` (`po_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- GRN LINE ITEM TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `grn_line_item` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `grn_id` VARCHAR(36) NOT NULL,
    `po_line_item_id` VARCHAR(36),
    `quantity_received` DECIMAL(15, 2),
    `quantity_accepted` DECIMAL(15, 2),
    `quantity_rejected` DECIMAL(15, 2),
    `rejection_reason` TEXT,
    `gl_account_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`grn_id`) REFERENCES `goods_receipt`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`po_line_item_id`) REFERENCES `po_line_item`(`id`),
    FOREIGN KEY (`gl_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_grn` (`grn_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
