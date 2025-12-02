-- ============================================================================
-- PURCHASE MODULE SCHEMA
-- ============================================================================
-- Multi-tenant Purchase Management with GL Integration
-- Supports: Purchase Requisitions, POs, Invoices, Payments, Account Payable (AP)

-- ============================================================================
-- VENDORS/SUPPLIERS
-- ============================================================================

CREATE TABLE IF NOT EXISTS vendors (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    vendor_code VARCHAR(50) NOT NULL UNIQUE,
    vendor_name VARCHAR(255) NOT NULL,
    vendor_type ENUM('Supplier', 'Service Provider', 'Contractor') NOT NULL,
    contact_person VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    gst_in VARCHAR(15),
    pan_number VARCHAR(10),
    bank_account_number VARCHAR(50),
    bank_ifsc_code VARCHAR(11),
    bank_name VARCHAR(255),
    payment_terms VARCHAR(100),
    credit_limit DECIMAL(15, 2),
    payment_status ENUM('Active', 'Inactive', 'On Hold') DEFAULT 'Active',
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_vendor (tenant_id, vendor_code),
    INDEX idx_created_at (created_at),
    CONSTRAINT fk_vendor_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- ============================================================================
-- PURCHASE REQUISITIONS
-- ============================================================================

CREATE TABLE IF NOT EXISTS purchase_requisitions (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    requisition_number VARCHAR(50) NOT NULL UNIQUE,
    requisition_date DATE NOT NULL,
    required_by_date DATE,
    requested_by VARCHAR(255),
    department VARCHAR(100),
    description TEXT,
    requisition_status ENUM('Draft', 'Submitted', 'Approved', 'Rejected', 'Converted to PO') DEFAULT 'Draft',
    total_amount DECIMAL(15, 2),
    approval_status ENUM('Pending', 'Approved', 'Rejected') DEFAULT 'Pending',
    approved_by VARCHAR(255),
    approved_at TIMESTAMP NULL,
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_requisition (tenant_id, requisition_number),
    INDEX idx_status (requisition_status),
    INDEX idx_created_at (created_at),
    CONSTRAINT fk_req_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- ============================================================================
-- PURCHASE REQUISITION DETAILS
-- ============================================================================

CREATE TABLE IF NOT EXISTS purchase_requisition_details (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    requisition_id VARCHAR(36) NOT NULL,
    item_code VARCHAR(100),
    item_name VARCHAR(255) NOT NULL,
    description TEXT,
    quantity DECIMAL(10, 2) NOT NULL,
    unit_of_measure VARCHAR(50),
    estimated_unit_price DECIMAL(15, 2),
    estimated_amount DECIMAL(15, 2),
    line_number INT,
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_requisition (tenant_id, requisition_id),
    CONSTRAINT fk_req_detail_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_req_detail_req FOREIGN KEY (requisition_id) REFERENCES purchase_requisitions(id)
);

-- ============================================================================
-- PURCHASE ORDERS
-- ============================================================================

CREATE TABLE IF NOT EXISTS purchase_orders (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    po_number VARCHAR(50) NOT NULL UNIQUE,
    po_date DATE NOT NULL,
    vendor_id VARCHAR(36) NOT NULL,
    requisition_id VARCHAR(36),
    delivery_address TEXT,
    delivery_date DATE,
    po_status ENUM('Draft', 'Sent', 'Confirmed', 'Partially Received', 'Received', 'Invoiced', 'Cancelled') DEFAULT 'Draft',
    total_quantity DECIMAL(10, 2),
    subtotal DECIMAL(15, 2),
    tax_rate DECIMAL(5, 2),
    tax_amount DECIMAL(15, 2),
    shipping_charge DECIMAL(15, 2),
    total_amount DECIMAL(15, 2),
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_po (tenant_id, po_number),
    INDEX idx_vendor (vendor_id),
    INDEX idx_status (po_status),
    INDEX idx_created_at (created_at),
    CONSTRAINT fk_po_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_po_vendor FOREIGN KEY (vendor_id) REFERENCES vendors(id),
    CONSTRAINT fk_po_req FOREIGN KEY (requisition_id) REFERENCES purchase_requisitions(id)
);

-- ============================================================================
-- PURCHASE ORDER DETAILS
-- ============================================================================

CREATE TABLE IF NOT EXISTS purchase_order_details (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    po_id VARCHAR(36) NOT NULL,
    item_code VARCHAR(100),
    item_name VARCHAR(255) NOT NULL,
    description TEXT,
    quantity_ordered DECIMAL(10, 2) NOT NULL,
    quantity_received DECIMAL(10, 2) DEFAULT 0,
    unit_of_measure VARCHAR(50),
    unit_price DECIMAL(15, 2) NOT NULL,
    line_amount DECIMAL(15, 2),
    line_number INT,
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_po (tenant_id, po_id),
    CONSTRAINT fk_po_detail_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_po_detail_po FOREIGN KEY (po_id) REFERENCES purchase_orders(id)
);

-- ============================================================================
-- GOODS RECEIPT NOTES / INWARD GOODS
-- ============================================================================

CREATE TABLE IF NOT EXISTS goods_receipt_notes (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    grn_number VARCHAR(50) NOT NULL UNIQUE,
    po_id VARCHAR(36) NOT NULL,
    receipt_date DATE NOT NULL,
    received_by VARCHAR(255),
    warehouse VARCHAR(100),
    receipt_status ENUM('Pending Inspection', 'Inspected', 'Accepted', 'Partially Accepted', 'Rejected') DEFAULT 'Pending Inspection',
    inspector_name VARCHAR(255),
    inspection_date DATE,
    discrepancies TEXT,
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_grn (tenant_id, grn_number),
    INDEX idx_po (po_id),
    INDEX idx_created_at (created_at),
    CONSTRAINT fk_grn_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_grn_po FOREIGN KEY (po_id) REFERENCES purchase_orders(id)
);

-- ============================================================================
-- GRN DETAILS
-- ============================================================================

CREATE TABLE IF NOT EXISTS goods_receipt_note_details (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    grn_id VARCHAR(36) NOT NULL,
    po_detail_id VARCHAR(36),
    item_code VARCHAR(100),
    item_name VARCHAR(255) NOT NULL,
    quantity_received DECIMAL(10, 2) NOT NULL,
    quantity_accepted DECIMAL(10, 2),
    quantity_rejected DECIMAL(10, 2),
    unit_price DECIMAL(15, 2),
    receipt_amount DECIMAL(15, 2),
    batch_number VARCHAR(100),
    expiry_date DATE,
    line_number INT,
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_grn (tenant_id, grn_id),
    CONSTRAINT fk_grn_detail_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_grn_detail_grn FOREIGN KEY (grn_id) REFERENCES goods_receipt_notes(id),
    CONSTRAINT fk_grn_detail_po_detail FOREIGN KEY (po_detail_id) REFERENCES purchase_order_details(id)
);

-- ============================================================================
-- PURCHASE INVOICES (BILLS)
-- ============================================================================

CREATE TABLE IF NOT EXISTS purchase_invoices (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    po_id VARCHAR(36),
    vendor_id VARCHAR(36) NOT NULL,
    invoice_date DATE NOT NULL,
    due_date DATE,
    invoice_status ENUM('Draft', 'Received', 'Matched', 'Approved', 'Paid', 'Partially Paid', 'Cancelled') DEFAULT 'Draft',
    subtotal DECIMAL(15, 2),
    tax_rate DECIMAL(5, 2),
    tax_amount DECIMAL(15, 2),
    shipping_charge DECIMAL(15, 2),
    discount_amount DECIMAL(15, 2),
    invoice_amount DECIMAL(15, 2),
    amount_paid DECIMAL(15, 2) DEFAULT 0,
    amount_due DECIMAL(15, 2),
    approval_status ENUM('Pending', 'Approved', 'Rejected') DEFAULT 'Pending',
    approved_by VARCHAR(255),
    approved_at TIMESTAMP NULL,
    three_way_match BOOLEAN DEFAULT FALSE,
    match_status ENUM('Not Matched', 'Partially Matched', 'Matched', 'Variance') DEFAULT 'Not Matched',
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_invoice (tenant_id, invoice_number),
    INDEX idx_vendor (vendor_id),
    INDEX idx_status (invoice_status),
    INDEX idx_due_date (due_date),
    INDEX idx_created_at (created_at),
    CONSTRAINT fk_invoice_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_invoice_vendor FOREIGN KEY (vendor_id) REFERENCES vendors(id),
    CONSTRAINT fk_invoice_po FOREIGN KEY (po_id) REFERENCES purchase_orders(id)
);

-- ============================================================================
-- PURCHASE INVOICE DETAILS
-- ============================================================================

CREATE TABLE IF NOT EXISTS purchase_invoice_details (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    invoice_id VARCHAR(36) NOT NULL,
    po_detail_id VARCHAR(36),
    item_code VARCHAR(100),
    item_name VARCHAR(255) NOT NULL,
    description TEXT,
    quantity DECIMAL(10, 2) NOT NULL,
    unit_price DECIMAL(15, 2) NOT NULL,
    line_amount DECIMAL(15, 2),
    tax_rate DECIMAL(5, 2),
    tax_amount DECIMAL(15, 2),
    line_number INT,
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_invoice (tenant_id, invoice_id),
    CONSTRAINT fk_invoice_detail_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_invoice_detail_invoice FOREIGN KEY (invoice_id) REFERENCES purchase_invoices(id),
    CONSTRAINT fk_invoice_detail_po_detail FOREIGN KEY (po_detail_id) REFERENCES purchase_order_details(id)
);

-- ============================================================================
-- PURCHASE PAYMENTS
-- ============================================================================

CREATE TABLE IF NOT EXISTS purchase_payments (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    invoice_id VARCHAR(36) NOT NULL,
    payment_number VARCHAR(50) NOT NULL UNIQUE,
    payment_date DATE NOT NULL,
    payment_method ENUM('Check', 'Bank Transfer', 'Credit Card', 'Cash', 'Other') NOT NULL,
    payment_amount DECIMAL(15, 2) NOT NULL,
    payment_status ENUM('Pending', 'Cleared', 'Bounced', 'Cancelled') DEFAULT 'Pending',
    reference_number VARCHAR(100),
    check_date DATE,
    check_number VARCHAR(50),
    bank_name VARCHAR(255),
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_payment (tenant_id, payment_number),
    INDEX idx_invoice (invoice_id),
    INDEX idx_payment_status (payment_status),
    INDEX idx_created_at (created_at),
    CONSTRAINT fk_payment_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_payment_invoice FOREIGN KEY (invoice_id) REFERENCES purchase_invoices(id)
);

-- ============================================================================
-- ACCOUNTS PAYABLE (AP) AGING
-- ============================================================================

CREATE TABLE IF NOT EXISTS ap_aging_analysis (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    vendor_id VARCHAR(36) NOT NULL,
    invoice_id VARCHAR(36) NOT NULL,
    invoice_number VARCHAR(50),
    invoice_date DATE,
    due_date DATE,
    invoice_amount DECIMAL(15, 2),
    amount_paid DECIMAL(15, 2),
    amount_due DECIMAL(15, 2),
    days_overdue INT,
    aging_bucket ENUM('Current', '1-30 days', '31-60 days', '61-90 days', '90+ days') DEFAULT 'Current',
    aging_days INT,
    payment_status ENUM('Unpaid', 'Partially Paid', 'Paid', 'Overdue') DEFAULT 'Unpaid',
    as_of_date DATE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_vendor (tenant_id, vendor_id),
    INDEX idx_aging_bucket (aging_bucket),
    INDEX idx_as_of_date (as_of_date),
    CONSTRAINT fk_ap_aging_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_ap_aging_vendor FOREIGN KEY (vendor_id) REFERENCES vendors(id),
    CONSTRAINT fk_ap_aging_invoice FOREIGN KEY (invoice_id) REFERENCES purchase_invoices(id)
);

-- ============================================================================
-- PURCHASE AUDIT LOG
-- ============================================================================

CREATE TABLE IF NOT EXISTS purchase_audit_log (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    document_type VARCHAR(100),
    document_id VARCHAR(36),
    document_number VARCHAR(50),
    action VARCHAR(100),
    changed_from TEXT,
    changed_to TEXT,
    changed_by VARCHAR(255),
    change_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(50),
    notes TEXT,
    
    INDEX idx_tenant_document (tenant_id, document_id),
    INDEX idx_timestamp (change_timestamp),
    CONSTRAINT fk_audit_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- ============================================================================
-- PURCHASE SETTINGS / CONFIGURATION
-- ============================================================================

CREATE TABLE IF NOT EXISTS purchase_settings (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL UNIQUE,
    po_prefix VARCHAR(10) DEFAULT 'PO',
    po_next_number INT DEFAULT 1000,
    grn_prefix VARCHAR(10) DEFAULT 'GRN',
    grn_next_number INT DEFAULT 1000,
    invoice_prefix VARCHAR(10) DEFAULT 'PI',
    invoice_next_number INT DEFAULT 1000,
    require_grn_for_invoice BOOLEAN DEFAULT TRUE,
    require_po_for_grn BOOLEAN DEFAULT TRUE,
    require_approval_for_po BOOLEAN DEFAULT TRUE,
    three_way_match_enabled BOOLEAN DEFAULT TRUE,
    default_tax_rate DECIMAL(5, 2) DEFAULT 18,
    default_payment_terms VARCHAR(100),
    enable_order_tracking BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_settings_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

CREATE INDEX idx_purchase_invoice_vendor_date ON purchase_invoices(tenant_id, vendor_id, invoice_date);
CREATE INDEX idx_purchase_order_vendor_date ON purchase_orders(tenant_id, vendor_id, po_date);
CREATE INDEX idx_purchase_order_status_date ON purchase_orders(tenant_id, po_status, po_date);
CREATE INDEX idx_grn_po_date ON goods_receipt_notes(tenant_id, po_id, receipt_date);
CREATE INDEX idx_ap_aging_due_date ON ap_aging_analysis(tenant_id, due_date, payment_status);
