-- Migration: 009_sales_module_schema.sql
-- Purpose: Create Sales (CRM + Order Management) module tables
-- Date: 2025-11-25

-- ============================================================================
-- LEAD MANAGEMENT TABLES
-- ============================================================================

CREATE TABLE sales_leads (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    lead_code VARCHAR(50) NOT NULL UNIQUE,
    
    -- Lead Information
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    company_name VARCHAR(150),
    industry VARCHAR(100),
    
    -- Lead Status
    status ENUM('new', 'contacted', 'qualified', 'negotiation', 'converted', 'lost') DEFAULT 'new',
    probability DECIMAL(5,2) DEFAULT 0,
    
    -- Source & Campaign
    source VARCHAR(50), -- 'website', 'email', 'phone', 'referral', 'event', 'social'
    campaign_id VARCHAR(36),
    
    -- Assignment
    assigned_to VARCHAR(36),
    assigned_date DATETIME,
    
    -- Conversion
    converted_to_customer BOOLEAN DEFAULT FALSE,
    customer_id VARCHAR(36),
    
    -- Tracking
    next_action_date DATETIME,
    next_action_notes TEXT,
    
    -- Timestamps
    created_by VARCHAR(36),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- ============================================================================
-- CUSTOMER MANAGEMENT TABLES
-- ============================================================================

CREATE TABLE sales_customers (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    customer_code VARCHAR(50) NOT NULL UNIQUE,
    
    -- Customer Information
    customer_name VARCHAR(150) NOT NULL,
    business_name VARCHAR(150),
    business_type VARCHAR(100), -- 'individual', 'proprietorship', 'partnership', 'pvt_ltd', 'public_ltd'
    industry VARCHAR(100),
    
    -- Contact Information
    primary_contact_name VARCHAR(100),
    primary_email VARCHAR(100),
    primary_phone VARCHAR(20),
    
    -- Address
    billing_address TEXT,
    billing_city VARCHAR(50),
    billing_state VARCHAR(50),
    billing_country VARCHAR(50) DEFAULT 'India',
    billing_zip VARCHAR(10),
    
    shipping_address TEXT,
    shipping_city VARCHAR(50),
    shipping_state VARCHAR(50),
    shipping_country VARCHAR(50) DEFAULT 'India',
    shipping_zip VARCHAR(10),
    
    -- Tax Information
    pan_number VARCHAR(20),
    gst_number VARCHAR(20),
    
    -- Commercial Terms
    credit_limit DECIMAL(15,2) DEFAULT 0,
    credit_days INT DEFAULT 0,
    payment_terms VARCHAR(100), -- '30 days net', 'COD', etc.
    
    -- Classification
    customer_category VARCHAR(50), -- 'gold', 'silver', 'bronze', 'regular'
    
    -- Status
    status ENUM('active', 'inactive', 'blocked') DEFAULT 'active',
    
    -- Account Balance
    current_balance DECIMAL(15,2) DEFAULT 0,
    
    -- Tracking
    created_by VARCHAR(36),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- ============================================================================
-- QUOTATION MANAGEMENT TABLES
-- ============================================================================

CREATE TABLE sales_quotations (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    quotation_number VARCHAR(50) NOT NULL UNIQUE,
    
    -- Customer Reference
    customer_id VARCHAR(36) NOT NULL,
    
    -- Quotation Details
    quotation_date DATETIME NOT NULL,
    valid_until DATETIME,
    
    -- Financial Summary
    subtotal_amount DECIMAL(15,2) DEFAULT 0,
    discount_amount DECIMAL(15,2) DEFAULT 0,
    tax_amount DECIMAL(15,2) DEFAULT 0,
    total_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Status
    status ENUM('draft', 'sent', 'accepted', 'rejected', 'expired', 'converted_to_order') DEFAULT 'draft',
    
    -- Conversion
    converted_to_order BOOLEAN DEFAULT FALSE,
    sales_order_id VARCHAR(36),
    
    -- Additional Info
    notes TEXT,
    terms_and_conditions TEXT,
    
    -- Tracking
    created_by VARCHAR(36),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (customer_id) REFERENCES sales_customers(id)
);

CREATE TABLE sales_quotation_items (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    quotation_id VARCHAR(36) NOT NULL,
    
    -- Line Item Details
    line_number INT NOT NULL,
    description TEXT NOT NULL,
    product_service_code VARCHAR(50),
    
    -- Quantity & Pricing
    quantity DECIMAL(12,2) NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    line_total DECIMAL(15,2) DEFAULT 0,
    
    -- Discount
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Tax
    hsn_code VARCHAR(20),
    tax_rate DECIMAL(5,2) DEFAULT 0,
    tax_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_quotation_id (quotation_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (quotation_id) REFERENCES sales_quotations(id) ON DELETE CASCADE
);

-- ============================================================================
-- SALES ORDER MANAGEMENT TABLES
-- ============================================================================

CREATE TABLE sales_orders (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    order_number VARCHAR(50) NOT NULL UNIQUE,
    
    -- Customer Reference
    customer_id VARCHAR(36) NOT NULL,
    quotation_id VARCHAR(36),
    
    -- Order Dates
    order_date DATETIME NOT NULL,
    required_by_date DATETIME,
    
    -- Delivery
    delivery_location VARCHAR(250),
    delivery_instructions TEXT,
    
    -- Financial Summary
    subtotal_amount DECIMAL(15,2) DEFAULT 0,
    discount_amount DECIMAL(15,2) DEFAULT 0,
    tax_amount DECIMAL(15,2) DEFAULT 0,
    total_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Invoicing
    invoiced_amount DECIMAL(15,2) DEFAULT 0,
    pending_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Status
    status ENUM('draft', 'confirmed', 'partially_invoiced', 'invoiced', 'partially_delivered', 'delivered', 'cancelled') DEFAULT 'draft',
    
    -- Additional Info
    notes TEXT,
    
    -- Tracking
    created_by VARCHAR(36),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (customer_id) REFERENCES sales_customers(id),
    FOREIGN KEY (quotation_id) REFERENCES sales_quotations(id)
);

CREATE TABLE sales_order_items (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    order_id VARCHAR(36) NOT NULL,
    
    -- Line Item Details
    line_number INT NOT NULL,
    description TEXT NOT NULL,
    product_service_code VARCHAR(50),
    
    -- Quantity & Pricing
    ordered_quantity DECIMAL(12,2) NOT NULL,
    invoiced_quantity DECIMAL(12,2) DEFAULT 0,
    unit_price DECIMAL(15,2) NOT NULL,
    line_total DECIMAL(15,2) DEFAULT 0,
    
    -- Discount
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Tax
    hsn_code VARCHAR(20),
    tax_rate DECIMAL(5,2) DEFAULT 0,
    tax_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_order_id (order_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (order_id) REFERENCES sales_orders(id) ON DELETE CASCADE
);

-- ============================================================================
-- INVOICE MANAGEMENT TABLES
-- ============================================================================

CREATE TABLE sales_invoices (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    
    -- Customer Reference
    customer_id VARCHAR(36) NOT NULL,
    sales_order_id VARCHAR(36),
    
    -- Invoice Dates
    invoice_date DATETIME NOT NULL,
    due_date DATETIME,
    
    -- Financial Summary
    subtotal_amount DECIMAL(15,2) DEFAULT 0,
    discount_amount DECIMAL(15,2) DEFAULT 0,
    cgst_amount DECIMAL(15,2) DEFAULT 0,
    sgst_amount DECIMAL(15,2) DEFAULT 0,
    igst_amount DECIMAL(15,2) DEFAULT 0,
    total_tax DECIMAL(15,2) DEFAULT 0,
    total_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Payment Status
    payment_status ENUM('unpaid', 'partially_paid', 'paid', 'overdue', 'cancelled') DEFAULT 'unpaid',
    paid_amount DECIMAL(15,2) DEFAULT 0,
    pending_amount DECIMAL(15,2) DEFAULT 0,
    
    -- GL Integration
    ar_posting_status ENUM('not_posted', 'posted', 'reversed') DEFAULT 'not_posted',
    gl_reference_number VARCHAR(50),
    
    -- Document Tracking
    document_status ENUM('draft', 'issued', 'cancelled') DEFAULT 'draft',
    
    -- Notes
    notes TEXT,
    
    -- Tracking
    created_by VARCHAR(36),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_payment_status (payment_status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (customer_id) REFERENCES sales_customers(id),
    FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id)
);

CREATE TABLE sales_invoice_items (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    invoice_id VARCHAR(36) NOT NULL,
    
    -- Line Item Details
    line_number INT NOT NULL,
    description TEXT NOT NULL,
    hsn_code VARCHAR(20),
    
    -- Quantity & Pricing
    quantity DECIMAL(12,2) NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    line_total DECIMAL(15,2) DEFAULT 0,
    
    -- Discount
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Tax Breakdown
    cgst_rate DECIMAL(5,2) DEFAULT 0,
    cgst_amount DECIMAL(15,2) DEFAULT 0,
    sgst_rate DECIMAL(5,2) DEFAULT 0,
    sgst_amount DECIMAL(15,2) DEFAULT 0,
    igst_rate DECIMAL(5,2) DEFAULT 0,
    igst_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_invoice_id (invoice_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (invoice_id) REFERENCES sales_invoices(id) ON DELETE CASCADE
);

-- ============================================================================
-- PAYMENT TRACKING
-- ============================================================================

CREATE TABLE sales_payments (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    invoice_id VARCHAR(36) NOT NULL,
    
    -- Payment Details
    payment_date DATETIME NOT NULL,
    payment_amount DECIMAL(15,2) NOT NULL,
    payment_method ENUM('cheque', 'bank_transfer', 'cash', 'credit_card', 'digital_payment') DEFAULT 'bank_transfer',
    reference_number VARCHAR(100),
    
    -- Status
    payment_status ENUM('initiated', 'processed', 'confirmed', 'failed', 'cancelled') DEFAULT 'initiated',
    
    -- Notes
    notes TEXT,
    
    -- Tracking
    created_by VARCHAR(36),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_invoice_id (invoice_id),
    INDEX idx_payment_date (payment_date),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (invoice_id) REFERENCES sales_invoices(id)
);

-- ============================================================================
-- PERFORMANCE TRACKING
-- ============================================================================

CREATE TABLE sales_performance_metrics (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    salesperson_id VARCHAR(36) NOT NULL,
    
    -- Time Period
    metric_month DATE,
    
    -- Lead Metrics
    leads_generated INT DEFAULT 0,
    leads_qualified INT DEFAULT 0,
    conversion_rate DECIMAL(5,2) DEFAULT 0,
    
    -- Sales Metrics
    total_orders INT DEFAULT 0,
    total_order_value DECIMAL(15,2) DEFAULT 0,
    average_order_value DECIMAL(15,2) DEFAULT 0,
    
    -- Performance
    target_amount DECIMAL(15,2) DEFAULT 0,
    actual_amount DECIMAL(15,2) DEFAULT 0,
    achievement_percent DECIMAL(5,2) DEFAULT 0,
    
    -- Commission
    commission_rate DECIMAL(5,2) DEFAULT 0,
    commission_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_salesperson_id (salesperson_id),
    INDEX idx_metric_month (metric_month),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- ============================================================================
-- CUSTOMER CONTACTS
-- ============================================================================

CREATE TABLE sales_customer_contacts (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    customer_id VARCHAR(36) NOT NULL,
    
    -- Contact Information
    contact_name VARCHAR(100) NOT NULL,
    contact_title VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    mobile VARCHAR(20),
    
    -- Role
    contact_role VARCHAR(50), -- 'decision_maker', 'finance', 'technical', 'other'
    is_primary_contact BOOLEAN DEFAULT FALSE,
    
    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_customer_id (customer_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (customer_id) REFERENCES sales_customers(id) ON DELETE CASCADE
);

-- ============================================================================
-- DELIVERY TRACKING
-- ============================================================================

CREATE TABLE sales_delivery_notes (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    order_id VARCHAR(36) NOT NULL,
    
    -- Delivery Details
    delivery_date DATETIME NOT NULL,
    delivery_location VARCHAR(250),
    delivered_by VARCHAR(100),
    vehicle_number VARCHAR(20),
    
    -- Quantity
    delivered_quantity DECIMAL(12,2),
    
    -- Status
    delivery_status ENUM('in_transit', 'delivered', 'pending_pod') DEFAULT 'in_transit',
    
    -- Proof of Delivery
    pod_received BOOLEAN DEFAULT FALSE,
    pod_date DATETIME,
    receiver_name VARCHAR(100),
    receiver_signature_url VARCHAR(255),
    
    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_order_id (order_id),
    INDEX idx_delivery_date (delivery_date),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (order_id) REFERENCES sales_orders(id)
);

-- ============================================================================
-- RETURN & CREDIT NOTES
-- ============================================================================

CREATE TABLE sales_credit_notes (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    invoice_id VARCHAR(36) NOT NULL,
    
    -- Credit Note Details
    credit_note_number VARCHAR(50) NOT NULL UNIQUE,
    credit_note_date DATETIME NOT NULL,
    
    -- Reason
    reason VARCHAR(255),
    
    -- Amount
    total_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Status
    status ENUM('draft', 'issued', 'cancelled') DEFAULT 'draft',
    
    -- GL Integration
    ar_posting_status ENUM('not_posted', 'posted', 'reversed') DEFAULT 'not_posted',
    
    -- Timestamps
    created_by VARCHAR(36),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_invoice_id (invoice_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (invoice_id) REFERENCES sales_invoices(id)
);

CREATE TABLE sales_credit_note_items (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    credit_note_id VARCHAR(36) NOT NULL,
    
    -- Item Details
    line_number INT NOT NULL,
    description TEXT NOT NULL,
    quantity DECIMAL(12,2) NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    line_total DECIMAL(15,2) DEFAULT 0,
    
    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_credit_note_id (credit_note_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (credit_note_id) REFERENCES sales_credit_notes(id) ON DELETE CASCADE
);
