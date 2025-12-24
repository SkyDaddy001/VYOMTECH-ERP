-- ============================================================
-- MIGRATION 018: INVENTORY MANAGEMENT MODULE
-- Date: December 3, 2025
-- Purpose: Complete inventory management with stock tracking, warehouses, valuation
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- WAREHOUSE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `warehouse` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `warehouse_code` VARCHAR(50) NOT NULL,
    `warehouse_name` VARCHAR(255) NOT NULL,
    `warehouse_type` VARCHAR(50),
    `address` TEXT,
    `city` VARCHAR(100),
    `state` VARCHAR(100),
    `country` VARCHAR(100),
    `postal_code` VARCHAR(20),
    `manager_id` VARCHAR(36),
    `capacity` DECIMAL(18, 2),
    `current_utilization` DECIMAL(18, 2) DEFAULT 0,
    `is_active` BOOLEAN DEFAULT TRUE,
    `gl_inventory_account_id` VARCHAR(36),
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`gl_inventory_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_code` (`warehouse_code`),
    KEY `idx_is_active` (`is_active`),
    UNIQUE KEY `unique_warehouse` (`tenant_id`, `warehouse_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- INVENTORY ITEM (STOCK KEEPING UNIT - SKU) TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `inventory_item` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `sku` VARCHAR(50) NOT NULL,
    `item_name` VARCHAR(255) NOT NULL,
    `item_description` TEXT,
    `item_category` VARCHAR(100),
    `item_type` VARCHAR(50),
    `unit_of_measure` VARCHAR(20),
    `reorder_level` DECIMAL(18, 4),
    `reorder_quantity` DECIMAL(18, 4),
    `safety_stock` DECIMAL(18, 4),
    `lead_time_days` INT,
    `hsn_code` VARCHAR(50),
    `is_serialized` BOOLEAN DEFAULT FALSE,
    `is_batch_tracked` BOOLEAN DEFAULT FALSE,
    `item_status` VARCHAR(50) DEFAULT 'active',
    `gl_inventory_account_id` VARCHAR(36),
    `gl_expense_account_id` VARCHAR(36),
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`gl_inventory_account_id`) REFERENCES `chart_of_account`(`id`),
    FOREIGN KEY (`gl_expense_account_id`) REFERENCES `chart_of_account`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_sku` (`sku`),
    KEY `idx_category` (`item_category`),
    KEY `idx_status` (`item_status`),
    UNIQUE KEY `unique_sku` (`tenant_id`, `sku`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- INVENTORY ITEM VENDOR TABLE (Multiple vendors per item)
-- ============================================================
CREATE TABLE IF NOT EXISTS `inventory_item_vendor` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `vendor_id` VARCHAR(36) NOT NULL,
    `vendor_sku` VARCHAR(50),
    `vendor_part_number` VARCHAR(100),
    `lead_time_days` INT,
    `minimum_order_quantity` DECIMAL(18, 4),
    `unit_price` DECIMAL(18, 4),
    `last_price_date` DATE,
    `preferred_vendor` BOOLEAN DEFAULT FALSE,
    `quality_rating` DECIMAL(3, 1),
    `is_active` BOOLEAN DEFAULT TRUE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`vendor_id`) REFERENCES `vendor`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_item` (`inventory_item_id`),
    KEY `idx_vendor` (`vendor_id`),
    UNIQUE KEY `unique_item_vendor` (`inventory_item_id`, `vendor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- STOCK LEVEL TABLE (Current inventory quantity by warehouse)
-- ============================================================
CREATE TABLE IF NOT EXISTS `stock_level` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `warehouse_id` VARCHAR(36) NOT NULL,
    `quantity_on_hand` DECIMAL(18, 4) DEFAULT 0,
    `quantity_reserved` DECIMAL(18, 4) DEFAULT 0,
    `quantity_available` DECIMAL(18, 4) DEFAULT 0,
    `quantity_in_transit` DECIMAL(18, 4) DEFAULT 0,
    `last_counted_date` DATE,
    `recount_required` BOOLEAN DEFAULT FALSE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_item` (`inventory_item_id`),
    KEY `idx_warehouse` (`warehouse_id`),
    KEY `idx_recount_required` (`recount_required`),
    UNIQUE KEY `unique_stock` (`inventory_item_id`, `warehouse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- STOCK MOVEMENT TABLE (Detailed transaction log)
-- ============================================================
CREATE TABLE IF NOT EXISTS `stock_movement` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `warehouse_id` VARCHAR(36) NOT NULL,
    `movement_type` VARCHAR(50),
    `movement_date` DATE NOT NULL,
    `quantity_change` DECIMAL(18, 4),
    `reference_type` VARCHAR(50),
    `reference_id` VARCHAR(36),
    `from_location` VARCHAR(100),
    `to_location` VARCHAR(100),
    `batch_number` VARCHAR(100),
    `serial_numbers` TEXT,
    `unit_price` DECIMAL(18, 4),
    `total_value` DECIMAL(18, 2),
    `reason_code` VARCHAR(50),
    `notes` TEXT,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_item` (`inventory_item_id`),
    KEY `idx_warehouse` (`warehouse_id`),
    KEY `idx_date` (`movement_date`),
    KEY `idx_movement_type` (`movement_type`),
    KEY `idx_reference` (`reference_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- BATCH/LOT TRACKING TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `inventory_batch` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `batch_number` VARCHAR(100) NOT NULL,
    `manufacture_date` DATE,
    `expiry_date` DATE,
    `quantity_received` DECIMAL(18, 4),
    `quantity_remaining` DECIMAL(18, 4),
    `purchase_order_id` VARCHAR(36),
    `supplier_batch_number` VARCHAR(100),
    `quality_status` VARCHAR(50),
    `storage_location` VARCHAR(100),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`purchase_order_id`) REFERENCES `purchase_order`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_item` (`inventory_item_id`),
    KEY `idx_batch_number` (`batch_number`),
    KEY `idx_expiry_date` (`expiry_date`),
    KEY `idx_quality_status` (`quality_status`),
    UNIQUE KEY `unique_batch` (`tenant_id`, `batch_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- SERIAL NUMBER TRACKING TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `inventory_serial` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `serial_number` VARCHAR(100) NOT NULL,
    `batch_id` VARCHAR(36),
    `purchase_order_id` VARCHAR(36),
    `warranty_start_date` DATE,
    `warranty_end_date` DATE,
    `status` VARCHAR(50),
    `current_location` VARCHAR(100),
    `asset_id` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`batch_id`) REFERENCES `inventory_batch`(`id`),
    FOREIGN KEY (`purchase_order_id`) REFERENCES `purchase_order`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_item` (`inventory_item_id`),
    KEY `idx_serial_number` (`serial_number`),
    KEY `idx_status` (`status`),
    UNIQUE KEY `unique_serial` (`tenant_id`, `serial_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- INVENTORY VALUATION METHOD TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `inventory_valuation` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `valuation_method` VARCHAR(50),
    `weighted_average_cost` DECIMAL(18, 4),
    `fifo_cost` DECIMAL(18, 4),
    `lifo_cost` DECIMAL(18, 4),
    `standard_cost` DECIMAL(18, 4),
    `last_cost` DECIMAL(18, 4),
    `market_value` DECIMAL(18, 4),
    `replacement_cost` DECIMAL(18, 4),
    `valuation_date` DATE,
    `total_inventory_value` DECIMAL(18, 2),
    `last_valuation_date` DATE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_item` (`inventory_item_id`),
    KEY `idx_valuation_method` (`valuation_method`),
    UNIQUE KEY `unique_valuation` (`inventory_item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- STOCK ADJUSTMENT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `stock_adjustment` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `adjustment_number` VARCHAR(50) NOT NULL,
    `adjustment_date` DATE NOT NULL,
    `warehouse_id` VARCHAR(36),
    `adjustment_reason` VARCHAR(100),
    `adjustment_type` VARCHAR(50),
    `total_adjustment_value` DECIMAL(18, 2),
    `notes` TEXT,
    `created_by` VARCHAR(36),
    `approved_by` VARCHAR(36),
    `approved_at` TIMESTAMP,
    `journal_entry_id` VARCHAR(36),
    `status` VARCHAR(50) DEFAULT 'draft',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`),
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_number` (`adjustment_number`),
    KEY `idx_date` (`adjustment_date`),
    KEY `idx_status` (`status`),
    UNIQUE KEY `unique_adjustment` (`tenant_id`, `adjustment_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- STOCK ADJUSTMENT LINE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `stock_adjustment_line` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `adjustment_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `line_number` INT,
    `quantity_variance` DECIMAL(18, 4),
    `old_quantity` DECIMAL(18, 4),
    `new_quantity` DECIMAL(18, 4),
    `unit_cost` DECIMAL(18, 4),
    `variance_value` DECIMAL(18, 2),
    `reason_code` VARCHAR(50),
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`adjustment_id`) REFERENCES `stock_adjustment`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_adjustment` (`adjustment_id`),
    KEY `idx_item` (`inventory_item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PHYSICAL INVENTORY COUNT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `physical_inventory` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `count_number` VARCHAR(50) NOT NULL,
    `warehouse_id` VARCHAR(36) NOT NULL,
    `count_date` DATE NOT NULL,
    `count_start_time` TIME,
    `count_end_time` TIME,
    `total_items_counted` INT,
    `total_variance` DECIMAL(18, 2),
    `variance_percentage` DECIMAL(5, 2),
    `count_status` VARCHAR(50) DEFAULT 'in_progress',
    `counted_by_id` VARCHAR(36),
    `verified_by_id` VARCHAR(36),
    `verified_at` TIMESTAMP,
    `journal_entry_id` VARCHAR(36),
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_warehouse` (`warehouse_id`),
    KEY `idx_count_date` (`count_date`),
    KEY `idx_status` (`count_status`),
    UNIQUE KEY `unique_count` (`tenant_id`, `count_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PHYSICAL INVENTORY COUNT DETAIL TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `physical_inventory_detail` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `physical_inventory_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `system_quantity` DECIMAL(18, 4),
    `counted_quantity` DECIMAL(18, 4),
    `variance_quantity` DECIMAL(18, 4),
    `count_status` VARCHAR(50),
    `counted_by_id` VARCHAR(36),
    `count_time` TIMESTAMP,
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`physical_inventory_id`) REFERENCES `physical_inventory`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_inventory_count` (`physical_inventory_id`),
    KEY `idx_item` (`inventory_item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- INVENTORY TRANSFER TABLE (Inter-warehouse transfers)
-- ============================================================
CREATE TABLE IF NOT EXISTS `inventory_transfer` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `transfer_number` VARCHAR(50) NOT NULL,
    `from_warehouse_id` VARCHAR(36) NOT NULL,
    `to_warehouse_id` VARCHAR(36) NOT NULL,
    `transfer_date` DATE NOT NULL,
    `expected_receipt_date` DATE,
    `actual_receipt_date` DATE,
    `transfer_status` VARCHAR(50) DEFAULT 'draft',
    `total_items` INT,
    `total_quantity` DECIMAL(18, 4),
    `transfer_cost` DECIMAL(18, 2),
    `created_by` VARCHAR(36),
    `approved_by` VARCHAR(36),
    `received_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`from_warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`to_warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_from_warehouse` (`from_warehouse_id`),
    KEY `idx_to_warehouse` (`to_warehouse_id`),
    KEY `idx_transfer_date` (`transfer_date`),
    KEY `idx_status` (`transfer_status`),
    UNIQUE KEY `unique_transfer` (`tenant_id`, `transfer_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- INVENTORY TRANSFER LINE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `inventory_transfer_line` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `transfer_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `line_number` INT,
    `quantity_transferred` DECIMAL(18, 4),
    `quantity_received` DECIMAL(18, 4),
    `unit_cost` DECIMAL(18, 4),
    `line_status` VARCHAR(50),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`transfer_id`) REFERENCES `inventory_transfer`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_transfer` (`transfer_id`),
    KEY `idx_item` (`inventory_item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- MINIMUM STOCK ALERT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `min_stock_alert` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `warehouse_id` VARCHAR(36) NOT NULL,
    `alert_date` DATE NOT NULL,
    `current_stock` DECIMAL(18, 4),
    `reorder_level` DECIMAL(18, 4),
    `suggested_order_quantity` DECIMAL(18, 4),
    `alert_status` VARCHAR(50) DEFAULT 'active',
    `purchase_order_id` VARCHAR(36),
    `acknowledged_by` VARCHAR(36),
    `acknowledged_at` TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`purchase_order_id`) REFERENCES `purchase_order`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_item` (`inventory_item_id`),
    KEY `idx_warehouse` (`warehouse_id`),
    KEY `idx_alert_status` (`alert_status`),
    KEY `idx_alert_date` (`alert_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- INVENTORY DAMAGE/OBSOLESCENCE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `inventory_damage` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `damage_number` VARCHAR(50) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `warehouse_id` VARCHAR(36) NOT NULL,
    `damage_date` DATE NOT NULL,
    `damage_type` VARCHAR(50),
    `quantity_damaged` DECIMAL(18, 4),
    `unit_cost` DECIMAL(18, 4),
    `total_loss_value` DECIMAL(18, 2),
    `damage_reason` TEXT,
    `responsibility` VARCHAR(100),
    `insurance_claim` BOOLEAN DEFAULT FALSE,
    `claim_number` VARCHAR(50),
    `claim_amount` DECIMAL(18, 2),
    `journal_entry_id` VARCHAR(36),
    `status` VARCHAR(50) DEFAULT 'reported',
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_damage_date` (`damage_date`),
    KEY `idx_status` (`status`),
    UNIQUE KEY `unique_damage` (`tenant_id`, `damage_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
