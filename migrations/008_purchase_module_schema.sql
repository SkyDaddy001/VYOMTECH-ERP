-- ============================================================================
-- PHASE 3E: PURCHASE MODULE SCHEMA
-- Weeks 9-10: Complete purchase management with GRN, MRN, and contracts
-- ============================================================================

-- ============================================================================
-- 1. VENDOR MANAGEMENT TABLES
-- ============================================================================

CREATE TABLE IF NOT EXISTS `vendors` (
  `id` CHAR(26) PRIMARY KEY COMMENT 'ULID primary key',
  `tenant_id` CHAR(26) NOT NULL COMMENT 'Multi-tenant isolation',
  `vendor_code` VARCHAR(50) NOT NULL UNIQUE COMMENT 'Unique vendor identifier',
  `name` VARCHAR(255) NOT NULL COMMENT 'Vendor name',
  `email` VARCHAR(255) COMMENT 'Primary email',
  `phone` VARCHAR(20) COMMENT 'Contact phone',
  `address` TEXT COMMENT 'Full address',
  `city` VARCHAR(100) COMMENT 'City',
  `state` VARCHAR(100) COMMENT 'State/Province',
  `country` VARCHAR(100) COMMENT 'Country',
  `postal_code` VARCHAR(20) COMMENT 'Postal code',
  `tax_id` VARCHAR(50) COMMENT 'GST/Tax ID',
  `payment_terms` ENUM('COD', 'NET15', 'NET30', 'NET45', 'NET60') DEFAULT 'NET30' COMMENT 'Default payment terms',
  `vendor_type` ENUM('Manufacturer', 'Distributor', 'Retailer', 'Service_Provider', 'Contractor') COMMENT 'Type of vendor',
  `rating` DECIMAL(3,2) DEFAULT 0 COMMENT 'Vendor rating (0-5)',
  `is_active` TINYINT(1) DEFAULT 1 COMMENT 'Active status',
  `is_blocked` TINYINT(1) DEFAULT 0 COMMENT 'Blocked from purchase',
  `created_by` CHAR(26) COMMENT 'User who created',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL COMMENT 'Soft delete',
  `status` VARCHAR(20) DEFAULT 'active' COMMENT 'active, inactive, blocked, deleted',
  INDEX idx_tenant_id (`tenant_id`),
  INDEX idx_vendor_code (`vendor_code`),
  INDEX idx_status (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `vendor_contacts` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `vendor_id` CHAR(26) NOT NULL,
  `contact_name` VARCHAR(255) NOT NULL,
  `title` VARCHAR(100) COMMENT 'Contact title/designation',
  `phone` VARCHAR(20),
  `email` VARCHAR(255),
  `is_primary` TINYINT(1) DEFAULT 0,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`vendor_id`) REFERENCES `vendors`(`id`),
  INDEX idx_vendor_id (`vendor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `vendor_addresses` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `vendor_id` CHAR(26) NOT NULL,
  `address_type` ENUM('Billing', 'Shipping', 'Factory') DEFAULT 'Shipping',
  `address_line1` VARCHAR(255) NOT NULL,
  `address_line2` VARCHAR(255),
  `city` VARCHAR(100),
  `state` VARCHAR(100),
  `country` VARCHAR(100),
  `postal_code` VARCHAR(20),
  `is_primary` TINYINT(1) DEFAULT 0,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`vendor_id`) REFERENCES `vendors`(`id`),
  INDEX idx_vendor_id (`vendor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 2. PURCHASE REQUISITION & ORDER TABLES
-- ============================================================================

CREATE TABLE IF NOT EXISTS `purchase_requisitions` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `requisition_number` VARCHAR(50) NOT NULL UNIQUE,
  `requisitioner_id` CHAR(26) NOT NULL COMMENT 'Employee requesting',
  `department` VARCHAR(100) COMMENT 'Department name',
  `request_date` DATE NOT NULL,
  `required_by_date` DATE COMMENT 'Expected delivery date',
  `purpose` TEXT COMMENT 'Purpose of requisition',
  `status` ENUM('Draft', 'Submitted', 'Approved', 'Rejected', 'Converted_to_PO') DEFAULT 'Draft',
  `approved_by` CHAR(26) COMMENT 'User who approved',
  `approved_at` TIMESTAMP NULL,
  `rejection_reason` TEXT,
  `created_by` CHAR(26),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  INDEX idx_tenant_id (`tenant_id`),
  INDEX idx_status (`status`),
  INDEX idx_requisition_number (`requisition_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `purchase_orders` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `po_number` VARCHAR(50) NOT NULL UNIQUE COMMENT 'Purchase order number',
  `vendor_id` CHAR(26) NOT NULL,
  `requisition_id` CHAR(26) COMMENT 'Linked requisition',
  `po_date` DATE NOT NULL,
  `delivery_date` DATE COMMENT 'Expected delivery',
  `total_amount` DECIMAL(15,2) NOT NULL,
  `tax_amount` DECIMAL(15,2) DEFAULT 0,
  `shipping_amount` DECIMAL(15,2) DEFAULT 0,
  `discount_amount` DECIMAL(15,2) DEFAULT 0,
  `net_amount` DECIMAL(15,2) NOT NULL,
  `payment_terms` VARCHAR(50),
  `delivery_location` TEXT COMMENT 'Where goods should be delivered',
  `special_instructions` TEXT,
  `status` ENUM('Draft', 'Sent', 'Acknowledged', 'Partial_Received', 'Fully_Received', 'Cancelled', 'Closed') DEFAULT 'Draft',
  `sent_to_vendor_at` TIMESTAMP NULL,
  `acknowledged_at` TIMESTAMP NULL,
  `created_by` CHAR(26),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  FOREIGN KEY (`vendor_id`) REFERENCES `vendors`(`id`),
  FOREIGN KEY (`requisition_id`) REFERENCES `purchase_requisitions`(`id`),
  INDEX idx_tenant_id (`tenant_id`),
  INDEX idx_po_number (`po_number`),
  INDEX idx_vendor_id (`vendor_id`),
  INDEX idx_status (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `po_line_items` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `po_id` CHAR(26) NOT NULL,
  `line_number` INT NOT NULL COMMENT 'Line item sequence',
  `product_code` VARCHAR(100),
  `description` TEXT NOT NULL,
  `quantity` DECIMAL(15,4) NOT NULL,
  `unit` VARCHAR(20) COMMENT 'Unit of measurement (pcs, kg, meter, etc)',
  `unit_price` DECIMAL(15,4) NOT NULL,
  `line_total` DECIMAL(15,2) NOT NULL,
  `hsn_code` VARCHAR(50) COMMENT 'HSN/SAC code for tax',
  `tax_rate` DECIMAL(5,2) DEFAULT 0 COMMENT 'Tax percentage',
  `tax_amount` DECIMAL(15,2) DEFAULT 0,
  `specification` TEXT COMMENT 'Product specification/notes',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`po_id`) REFERENCES `purchase_orders`(`id`),
  INDEX idx_po_id (`po_id`),
  INDEX idx_line_number (`line_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 3. GOODS RECEIPT & QUALITY INSPECTION TABLES
-- ============================================================================

CREATE TABLE IF NOT EXISTS `goods_receipts` (
  `id` CHAR(26) PRIMARY KEY COMMENT 'GRN ID',
  `tenant_id` CHAR(26) NOT NULL,
  `grn_number` VARCHAR(50) NOT NULL UNIQUE COMMENT 'Goods Receipt Note number',
  `po_id` CHAR(26) NOT NULL,
  `receipt_date` DATE NOT NULL,
  `received_by` CHAR(26) NOT NULL COMMENT 'Employee who received goods',
  `total_quantity_received` DECIMAL(15,4),
  `total_quantity_accepted` DECIMAL(15,4),
  `total_quantity_rejected` DECIMAL(15,4) DEFAULT 0,
  `delivery_note_number` VARCHAR(100) COMMENT 'Vendor delivery note number',
  `vehicle_number` VARCHAR(50) COMMENT 'Transport vehicle number',
  `driver_name` VARCHAR(100),
  `driver_phone` VARCHAR(20),
  `remarks` TEXT,
  `status` ENUM('Received', 'QC_In_Progress', 'QC_Passed', 'QC_Failed', 'Partial_Accepted', 'Rejected') DEFAULT 'Received',
  `qc_status` ENUM('Pending', 'In_Progress', 'Passed', 'Failed', 'Partial') DEFAULT 'Pending',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  FOREIGN KEY (`po_id`) REFERENCES `purchase_orders`(`id`),
  INDEX idx_tenant_id (`tenant_id`),
  INDEX idx_grn_number (`grn_number`),
  INDEX idx_po_id (`po_id`),
  INDEX idx_status (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `grn_line_items` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `grn_id` CHAR(26) NOT NULL,
  `po_line_item_id` CHAR(26),
  `line_number` INT,
  `product_code` VARCHAR(100),
  `description` TEXT,
  `po_quantity` DECIMAL(15,4) COMMENT 'Quantity ordered',
  `received_quantity` DECIMAL(15,4) COMMENT 'Quantity received',
  `accepted_quantity` DECIMAL(15,4) COMMENT 'Quantity accepted after QC',
  `rejected_quantity` DECIMAL(15,4) DEFAULT 0,
  `unit` VARCHAR(20),
  `rejection_reason` TEXT,
  `batch_number` VARCHAR(100) COMMENT 'Product batch/lot number',
  `expiry_date` DATE COMMENT 'Product expiry date',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`grn_id`) REFERENCES `goods_receipts`(`id`),
  FOREIGN KEY (`po_line_item_id`) REFERENCES `po_line_items`(`id`),
  INDEX idx_grn_id (`grn_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `quality_inspections` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `grn_id` CHAR(26) NOT NULL,
  `grn_line_item_id` CHAR(26),
  `inspection_date` DATE NOT NULL,
  `inspected_by` CHAR(26) NOT NULL COMMENT 'QC inspector',
  `inspection_type` ENUM('Visual', 'Functional', 'Dimensional', 'Lab_Test', 'Batch_Test') DEFAULT 'Visual',
  `status` ENUM('Passed', 'Failed', 'Partial_Pass', 'Conditional_Pass') DEFAULT 'Passed',
  `quantity_inspected` DECIMAL(15,4),
  `quantity_passed` DECIMAL(15,4),
  `quantity_failed` DECIMAL(15,4),
  `defects_found` TEXT COMMENT 'List of defects',
  `quality_score` DECIMAL(5,2) COMMENT 'Score out of 100',
  `notes` TEXT,
  `certificate_number` VARCHAR(100) COMMENT 'Lab test certificate',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`grn_id`) REFERENCES `goods_receipts`(`id`),
  INDEX idx_grn_id (`grn_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 4. MATERIAL RECEIPT NOTE (MRN) - WAREHOUSE ACCEPTANCE
-- ============================================================================

CREATE TABLE IF NOT EXISTS `material_receipt_notes` (
  `id` CHAR(26) PRIMARY KEY COMMENT 'MRN ID',
  `tenant_id` CHAR(26) NOT NULL,
  `mrn_number` VARCHAR(50) NOT NULL UNIQUE COMMENT 'Material Receipt Note number',
  `grn_id` CHAR(26) NOT NULL COMMENT 'Linked GRN',
  `warehouse_id` CHAR(26) COMMENT 'Receiving warehouse',
  `receipt_date` DATE NOT NULL,
  `accepted_by` CHAR(26) NOT NULL COMMENT 'Warehouse in-charge',
  `storage_location` VARCHAR(100) COMMENT 'Bin/shelf location',
  `total_quantity_accepted` DECIMAL(15,4),
  `total_quantity_in_stock` DECIMAL(15,4),
  `status` ENUM('Received', 'Stored', 'Available') DEFAULT 'Received',
  `remarks` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (`grn_id`) REFERENCES `goods_receipts`(`id`),
  INDEX idx_tenant_id (`tenant_id`),
  INDEX idx_mrn_number (`mrn_number`),
  INDEX idx_grn_id (`grn_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `mrn_line_items` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `mrn_id` CHAR(26) NOT NULL,
  `grn_line_item_id` CHAR(26),
  `line_number` INT,
  `product_code` VARCHAR(100),
  `description` TEXT,
  `quantity_accepted` DECIMAL(15,4),
  `unit` VARCHAR(20),
  `storage_location` VARCHAR(100),
  `batch_number` VARCHAR(100),
  `expiry_date` DATE,
  `warehouse_notes` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`mrn_id`) REFERENCES `material_receipt_notes`(`id`),
  INDEX idx_mrn_id (`mrn_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 5. CONTRACT MANAGEMENT - AGAINST BOQ
-- ============================================================================

CREATE TABLE IF NOT EXISTS `contracts` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `contract_number` VARCHAR(50) NOT NULL UNIQUE COMMENT 'Contract reference number',
  `vendor_id` CHAR(26) NOT NULL,
  `contract_date` DATE NOT NULL,
  `start_date` DATE,
  `end_date` DATE COMMENT 'Contract validity period',
  `contract_type` ENUM('Material', 'Labour', 'Service', 'Hybrid') DEFAULT 'Material' COMMENT 'Type of contract',
  `boq_id` CHAR(26) COMMENT 'Linked BOQ (if from Construction)',
  `total_contract_value` DECIMAL(15,2) NOT NULL,
  `currency` VARCHAR(3) DEFAULT 'INR',
  `payment_terms` VARCHAR(100),
  `delivery_schedule` TEXT COMMENT 'Delivery schedule details',
  `contract_status` ENUM('Draft', 'Sent', 'Accepted', 'Active', 'Completed', 'Cancelled') DEFAULT 'Draft',
  `contract_file_path` VARCHAR(255) COMMENT 'Path to contract document',
  `signed_date` DATE,
  `signed_by_vendor` TINYINT(1) DEFAULT 0,
  `signed_by_company` TINYINT(1) DEFAULT 0,
  `company_signatory` CHAR(26) COMMENT 'Who signed for company',
  `vendor_signatory` VARCHAR(100) COMMENT 'Vendor signatory name',
  `created_by` CHAR(26),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (`vendor_id`) REFERENCES `vendors`(`id`),
  INDEX idx_tenant_id (`tenant_id`),
  INDEX idx_contract_number (`contract_number`),
  INDEX idx_vendor_id (`vendor_id`),
  INDEX idx_contract_status (`contract_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `contract_line_items` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `contract_id` CHAR(26) NOT NULL,
  `line_number` INT,
  `item_code` VARCHAR(100),
  `item_description` TEXT NOT NULL,
  `item_type` ENUM('Material', 'Labour', 'Service', 'Hybrid') COMMENT 'Type within contract',
  `quantity` DECIMAL(15,4),
  `unit` VARCHAR(20),
  `unit_price` DECIMAL(15,4),
  `line_total` DECIMAL(15,2),
  `specification` TEXT COMMENT 'Specifications/standards',
  `delivery_location` VARCHAR(255) COMMENT 'Where to be delivered',
  `delivery_date` DATE COMMENT 'Scheduled delivery',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`contract_id`) REFERENCES `contracts`(`id`),
  INDEX idx_contract_id (`contract_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 6. MATERIAL LABOUR SERVICE CONTRACTS (HYBRID)
-- ============================================================================

CREATE TABLE IF NOT EXISTS `contract_materials` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `contract_id` CHAR(26) NOT NULL,
  `material_code` VARCHAR(100) NOT NULL,
  `material_description` TEXT,
  `quantity` DECIMAL(15,4) NOT NULL,
  `unit` VARCHAR(20),
  `unit_price` DECIMAL(15,4),
  `total_price` DECIMAL(15,2),
  `hsn_code` VARCHAR(50),
  `supplier_notes` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`contract_id`) REFERENCES `contracts`(`id`),
  INDEX idx_contract_id (`contract_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `contract_labour` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `contract_id` CHAR(26) NOT NULL,
  `skill_type` VARCHAR(100) NOT NULL COMMENT 'Skill category',
  `labour_category` ENUM('Skilled', 'Semi-Skilled', 'Unskilled') DEFAULT 'Skilled',
  `number_of_workers` INT,
  `duration_days` INT,
  `daily_rate` DECIMAL(12,2),
  `total_labour_cost` DECIMAL(15,2),
  `work_description` TEXT COMMENT 'Work to be performed',
  `labour_notes` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`contract_id`) REFERENCES `contracts`(`id`),
  INDEX idx_contract_id (`contract_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `contract_services` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `contract_id` CHAR(26) NOT NULL,
  `service_code` VARCHAR(100),
  `service_description` TEXT NOT NULL,
  `service_type` VARCHAR(100) COMMENT 'Type of service',
  `unit_of_service` VARCHAR(50) COMMENT 'hour, day, project, etc',
  `quantity` DECIMAL(12,2),
  `unit_price` DECIMAL(15,4),
  `total_service_cost` DECIMAL(15,2),
  `service_level_agreement` TEXT COMMENT 'SLA terms',
  `service_notes` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`contract_id`) REFERENCES `contracts`(`id`),
  INDEX idx_contract_id (`contract_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 7. VENDOR INVOICES & PAYMENT
-- ============================================================================

CREATE TABLE IF NOT EXISTS `vendor_invoices` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `invoice_number` VARCHAR(50) NOT NULL UNIQUE,
  `vendor_id` CHAR(26) NOT NULL,
  `po_id` CHAR(26),
  `grn_id` CHAR(26),
  `invoice_date` DATE NOT NULL,
  `due_date` DATE,
  `invoice_amount` DECIMAL(15,2) NOT NULL,
  `tax_amount` DECIMAL(15,2) DEFAULT 0,
  `discount_amount` DECIMAL(15,2) DEFAULT 0,
  `total_payable` DECIMAL(15,2) NOT NULL,
  `status` ENUM('Received', 'Approved', 'Rejected', 'Paid', 'Partially_Paid') DEFAULT 'Received',
  `matched_status` ENUM('Not_Matched', 'PO_Matched', 'GRN_Matched', 'Three_Way_Match') DEFAULT 'Not_Matched',
  `three_way_match` TINYINT(1) DEFAULT 0 COMMENT 'PO-GRN-Invoice matched',
  `received_at` TIMESTAMP,
  `approved_at` TIMESTAMP NULL,
  `approved_by` CHAR(26),
  `rejection_reason` TEXT,
  `created_by` CHAR(26),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (`vendor_id`) REFERENCES `vendors`(`id`),
  FOREIGN KEY (`po_id`) REFERENCES `purchase_orders`(`id`),
  FOREIGN KEY (`grn_id`) REFERENCES `goods_receipts`(`id`),
  INDEX idx_tenant_id (`tenant_id`),
  INDEX idx_invoice_number (`invoice_number`),
  INDEX idx_status (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `invoice_line_items` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `invoice_id` CHAR(26) NOT NULL,
  `line_number` INT,
  `description` TEXT,
  `quantity` DECIMAL(15,4),
  `unit_price` DECIMAL(15,4),
  `line_total` DECIMAL(15,2),
  `hsn_code` VARCHAR(50),
  `tax_rate` DECIMAL(5,2),
  `tax_amount` DECIMAL(15,2),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`invoice_id`) REFERENCES `vendor_invoices`(`id`),
  INDEX idx_invoice_id (`invoice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `vendor_payments` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `invoice_id` CHAR(26) NOT NULL,
  `payment_date` DATE NOT NULL,
  `payment_amount` DECIMAL(15,2) NOT NULL,
  `payment_method` ENUM('Cheque', 'Bank_Transfer', 'Cash', 'Credit_Card', 'Digital_Payment') DEFAULT 'Bank_Transfer',
  `reference_number` VARCHAR(100) COMMENT 'Cheque/transfer ref',
  `payment_status` ENUM('Initiated', 'Processed', 'Confirmed', 'Cancelled') DEFAULT 'Initiated',
  `paid_by` CHAR(26),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`invoice_id`) REFERENCES `vendor_invoices`(`id`),
  INDEX idx_invoice_id (`invoice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 8. VENDOR PERFORMANCE & METRICS
-- ============================================================================

CREATE TABLE IF NOT EXISTS `vendor_performance_metrics` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `vendor_id` CHAR(26) NOT NULL,
  `metric_month` DATE COMMENT 'Month for which metric calculated',
  `total_orders` INT DEFAULT 0,
  `orders_delivered_on_time` INT DEFAULT 0,
  `on_time_delivery_rate` DECIMAL(5,2) DEFAULT 0,
  `total_grn` INT DEFAULT 0,
  `grn_accepted_first_time` INT DEFAULT 0,
  `quality_acceptance_rate` DECIMAL(5,2) DEFAULT 0,
  `total_invoices` INT DEFAULT 0,
  `invoice_discrepancies` INT DEFAULT 0,
  `average_response_time_hours` DECIMAL(8,2),
  `overall_rating` DECIMAL(3,2) DEFAULT 0,
  `notes` TEXT,
  `calculated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`vendor_id`) REFERENCES `vendors`(`id`),
  INDEX idx_vendor_id (`vendor_id`),
  INDEX idx_metric_month (`metric_month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `purchase_approvals` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `po_id` CHAR(26) COMMENT 'PO being approved',
  `approval_level` INT COMMENT 'Approval hierarchy level',
  `approver_id` CHAR(26) NOT NULL,
  `approval_status` ENUM('Pending', 'Approved', 'Rejected', 'On_Hold') DEFAULT 'Pending',
  `approval_date` TIMESTAMP NULL,
  `comments` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`po_id`) REFERENCES `purchase_orders`(`id`),
  INDEX idx_po_id (`po_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 9. AUDIT LOG & HISTORY
-- ============================================================================

CREATE TABLE IF NOT EXISTS `purchase_audit_log` (
  `id` CHAR(26) PRIMARY KEY,
  `tenant_id` CHAR(26) NOT NULL,
  `entity_type` VARCHAR(50) COMMENT 'PO, GRN, Invoice, etc',
  `entity_id` CHAR(26),
  `action_type` VARCHAR(50) COMMENT 'Create, Update, Delete, Approve, etc',
  `old_values` JSON COMMENT 'Previous values',
  `new_values` JSON COMMENT 'New values',
  `changed_by` CHAR(26),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_tenant_id (`tenant_id`),
  INDEX idx_entity_id (`entity_id`),
  INDEX idx_created_at (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 10. CREATE INDEXES FOR PERFORMANCE
-- ============================================================================

CREATE INDEX idx_vendors_tenant_status ON vendors(tenant_id, status);
CREATE INDEX idx_purchase_orders_tenant_status ON purchase_orders(tenant_id, status);
CREATE INDEX idx_goods_receipts_tenant_status ON goods_receipts(tenant_id, status);
CREATE INDEX idx_vendor_invoices_tenant_status ON vendor_invoices(tenant_id, status);
CREATE INDEX idx_contracts_tenant_status ON contracts(tenant_id, contract_status);

-- ============================================================================
-- 11. INSERTION OF AUDIT TRIGGER
-- ============================================================================

DELIMITER //

CREATE TRIGGER purchase_orders_audit AFTER UPDATE ON purchase_orders
FOR EACH ROW
BEGIN
  INSERT INTO purchase_audit_log (
    id, tenant_id, entity_type, entity_id, action_type, old_values, new_values, changed_by, created_at
  ) VALUES (
    UNHEX(REPLACE(UUID(), '-', '')),
    NEW.tenant_id,
    'PurchaseOrder',
    NEW.id,
    'Update',
    JSON_OBJECT(
      'status', OLD.status,
      'total_amount', OLD.total_amount,
      'updated_at', OLD.updated_at
    ),
    JSON_OBJECT(
      'status', NEW.status,
      'total_amount', NEW.total_amount,
      'updated_at', NEW.updated_at
    ),
    NEW.updated_by,
    NOW()
  );
END //

CREATE TRIGGER goods_receipts_audit AFTER UPDATE ON goods_receipts
FOR EACH ROW
BEGIN
  INSERT INTO purchase_audit_log (
    id, tenant_id, entity_type, entity_id, action_type, old_values, new_values, changed_by, created_at
  ) VALUES (
    UNHEX(REPLACE(UUID(), '-', '')),
    NEW.tenant_id,
    'GoodsReceipt',
    NEW.id,
    'Update',
    JSON_OBJECT(
      'status', OLD.status,
      'qc_status', OLD.qc_status,
      'updated_at', OLD.updated_at
    ),
    JSON_OBJECT(
      'status', NEW.status,
      'qc_status', NEW.qc_status,
      'updated_at', NEW.updated_at
    ),
    NULL,
    NOW()
  );
END //

CREATE TRIGGER vendor_invoices_audit AFTER UPDATE ON vendor_invoices
FOR EACH ROW
BEGIN
  INSERT INTO purchase_audit_log (
    id, tenant_id, entity_type, entity_id, action_type, old_values, new_values, changed_by, created_at
  ) VALUES (
    UNHEX(REPLACE(UUID(), '-', '')),
    NEW.tenant_id,
    'VendorInvoice',
    NEW.id,
    'Update',
    JSON_OBJECT(
      'status', OLD.status,
      'matched_status', OLD.matched_status,
      'updated_at', OLD.updated_at
    ),
    JSON_OBJECT(
      'status', NEW.status,
      'matched_status', NEW.matched_status,
      'updated_at', NEW.updated_at
    ),
    NEW.approved_by,
    NOW()
  );
END //

DELIMITER ;

-- ============================================================================
-- SCHEMA COMPLETE
-- Total Tables: 18
-- Total Views: 0
-- Status: Ready for application integration
-- ============================================================================
