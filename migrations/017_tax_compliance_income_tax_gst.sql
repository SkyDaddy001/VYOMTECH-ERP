-- ============================================================================
-- INCOME TAX & GST COMPLIANCE MODULE
-- ============================================================================
-- Comprehensive tax compliance for Income Tax and GST as per Indian law

CREATE TABLE IF NOT EXISTS tax_configuration (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Financial Year Configuration
    financial_year_start TIMESTAMP NOT NULL,
    financial_year_end TIMESTAMP NOT NULL,
    fiscal_year INT NOT NULL, -- e.g., 2024-2025
    
    -- Income Tax Configuration
    it_assessment_year INT,
    it_registration_number VARCHAR(50),
    it_pan VARCHAR(20),
    it_circle VARCHAR(100),
    it_assessing_officer VARCHAR(255),
    
    -- GST Configuration
    gst_registration_number VARCHAR(50) UNIQUE,
    gst_registration_date TIMESTAMP,
    gst_status ENUM('Active', 'Inactive', 'Cancelled', 'Suspended') DEFAULT 'Active',
    gst_return_filing_frequency ENUM('Monthly', 'Quarterly', 'Annual') DEFAULT 'Monthly',
    
    -- Tax Filing Status
    it_filing_status ENUM('Not_Filed', 'Filed', 'Processing', 'Assessed', 'Under_Appeal') DEFAULT 'Not_Filed',
    gst_filing_status ENUM('Not_Filed', 'Filed', 'Paid', 'Overdue') DEFAULT 'Not_Filed',
    
    -- Contact Information
    compliance_officer_name VARCHAR(255),
    compliance_officer_email VARCHAR(255),
    compliance_officer_phone VARCHAR(20),
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    last_audit_date TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_fy (tenant_id, fiscal_year),
    INDEX idx_pan_gstin (it_pan, gst_registration_number),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- INCOME TAX COMPLIANCE TRACKING
-- ============================================================================

CREATE TABLE IF NOT EXISTS income_tax_compliance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Financial Year
    fiscal_year INT NOT NULL,
    assessment_year INT NOT NULL,
    
    -- Income Calculation
    total_gross_income DECIMAL(15, 2) DEFAULT 0,
    total_deductions DECIMAL(15, 2) DEFAULT 0,
    taxable_income DECIMAL(15, 2) DEFAULT 0,
    
    -- Income Breakup
    salary_income DECIMAL(15, 2) DEFAULT 0,
    business_income DECIMAL(15, 2) DEFAULT 0,
    capital_gains_income DECIMAL(15, 2) DEFAULT 0,
    rental_income DECIMAL(15, 2) DEFAULT 0,
    other_income DECIMAL(15, 2) DEFAULT 0,
    
    -- Deductions (Section 80 Categories)
    section_80c_deduction DECIMAL(15, 2) DEFAULT 0, -- Life Insurance, PPF, Children Education, etc.
    section_80d_deduction DECIMAL(15, 2) DEFAULT 0, -- Medical Insurance
    section_80e_deduction DECIMAL(15, 2) DEFAULT 0, -- Education Loan Interest
    section_80g_deduction DECIMAL(15, 2) DEFAULT 0, -- Donations
    section_80ggc_deduction DECIMAL(15, 2) DEFAULT 0, -- Rent Deduction
    other_deductions DECIMAL(15, 2) DEFAULT 0,
    
    -- Tax Calculation
    tax_before_surcharge DECIMAL(15, 2) DEFAULT 0,
    surcharge_amount DECIMAL(15, 2) DEFAULT 0,
    cess_amount DECIMAL(15, 2) DEFAULT 0,
    total_tax_liability DECIMAL(15, 2) DEFAULT 0,
    
    -- Tax Payment
    advance_tax_paid DECIMAL(15, 2) DEFAULT 0,
    tds_credit DECIMAL(15, 2) DEFAULT 0,
    total_tax_paid DECIMAL(15, 2) DEFAULT 0,
    tax_payable_refundable DECIMAL(15, 2) DEFAULT 0,
    
    -- Tax Return Filing
    itr_form_type VARCHAR(20), -- ITR-1, ITR-2, ITR-3, ITR-4, ITR-5, ITR-6, ITR-7
    return_filing_date TIMESTAMP,
    return_filed_status ENUM('Not_Filed', 'Filed', 'Processed', 'Accepted', 'Rejected', 'Under_Scrutiny') DEFAULT 'Not_Filed',
    ack_number VARCHAR(100), -- ITR Acknowledgment Number
    filing_acknowledgement_date TIMESTAMP,
    
    -- Scrutiny/Assessment
    scrutiny_initiated BOOLEAN DEFAULT false,
    scrutiny_date TIMESTAMP,
    scrutiny_response_date TIMESTAMP,
    assessment_order_date TIMESTAMP,
    additional_tax_assessed DECIMAL(15, 2) DEFAULT 0,
    
    -- Compliance Status
    is_compliant BOOLEAN DEFAULT true,
    non_compliance_reason TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_fy (tenant_id, fiscal_year),
    INDEX idx_filing_status (return_filed_status),
    INDEX idx_compliance (is_compliant),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- TDS (TAX DEDUCTED AT SOURCE) COMPLIANCE
-- ============================================================================

CREATE TABLE IF NOT EXISTS tds_compliance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Financial Year
    fiscal_year INT NOT NULL,
    
    -- TDS Categories (Sections)
    tds_section VARCHAR(20) NOT NULL, -- 192 (Salary), 193 (Interest), 194A (Commission), 194C (Contractors), 194D, 194H, 194I, 194J, 194LA, 194LB
    tds_category VARCHAR(100),
    
    -- TDS Calculation
    total_payments_made DECIMAL(15, 2) DEFAULT 0,
    tds_rate DECIMAL(5, 2) NOT NULL,
    tds_amount DEDUCTED DECIMAL(15, 2) DEFAULT 0,
    
    -- TDS Deposit
    tds_deposit_date TIMESTAMP,
    challan_number VARCHAR(50),
    bank_name VARCHAR(255),
    
    -- TDS Return Filing (Quarterly)
    tds_return_period VARCHAR(20), -- Q1, Q2, Q3, Q4
    tds_return_filing_date TIMESTAMP,
    tds_return_status ENUM('Not_Filed', 'Filed', 'Accepted', 'Rejected') DEFAULT 'Not_Filed',
    
    -- Reconciliation
    payee_details_reconciled BOOLEAN DEFAULT false,
    reconciliation_date TIMESTAMP,
    reconciliation_variance DECIMAL(15, 2) DEFAULT 0,
    
    -- Annual TDS Summary
    annual_tds_summary_filed BOOLEAN DEFAULT false,
    annual_summary_filing_date TIMESTAMP,
    
    is_compliant BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_fy (tenant_id, fiscal_year),
    INDEX idx_tds_section (tds_section),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- GST (GOODS & SERVICES TAX) COMPLIANCE
-- ============================================================================

CREATE TABLE IF NOT EXISTS gst_compliance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Financial Year & Period
    fiscal_year INT NOT NULL,
    return_period ENUM('Monthly', 'Quarterly') NOT NULL,
    month_year TIMESTAMP NOT NULL,
    
    -- Outward Supplies (GST Payable)
    total_outward_supplies DECIMAL(15, 2) DEFAULT 0,
    
    -- Sales Breakup
    gst_5_sales DECIMAL(15, 2) DEFAULT 0,
    gst_12_sales DECIMAL(15, 2) DEFAULT 0,
    gst_18_sales DECIMAL(15, 2) DEFAULT 0,
    gst_28_sales DECIMAL(15, 2) DEFAULT 0,
    
    -- Intra-state vs Inter-state
    intra_state_sales DECIMAL(15, 2) DEFAULT 0,
    inter_state_sales DECIMAL(15, 2) DEFAULT 0,
    exports DECIMAL(15, 2) DEFAULT 0,
    
    -- Output GST
    output_gst_5 DECIMAL(15, 2) DEFAULT 0,
    output_gst_12 DECIMAL(15, 2) DEFAULT 0,
    output_gst_18 DECIMAL(15, 2) DEFAULT 0,
    output_gst_28 DECIMAL(15, 2) DEFAULT 0,
    output_gst_cess DECIMAL(15, 2) DEFAULT 0,
    total_output_gst DECIMAL(15, 2) DEFAULT 0,
    
    -- Inward Supplies (Input GST)
    total_inward_supplies DECIMAL(15, 2) DEFAULT 0,
    
    -- Input Credit
    input_gst_5 DECIMAL(15, 2) DEFAULT 0,
    input_gst_12 DECIMAL(15, 2) DEFAULT 0,
    input_gst_18 DECIMAL(15, 2) DEFAULT 0,
    input_gst_28 DECIMAL(15, 2) DEFAULT 0,
    input_gst_cess DECIMAL(15, 2) DEFAULT 0,
    total_input_gst DECIMAL(15, 2) DEFAULT 0,
    
    -- Interest-Free Credit
    interest_free_credit_carried DECIMAL(15, 2) DEFAULT 0,
    
    -- Net GST Payable / Refundable
    net_gst_payable DECIMAL(15, 2) DEFAULT 0,
    refund_claim DECIMAL(15, 2) DEFAULT 0,
    refund_status ENUM('Not_Claimed', 'Claimed', 'Processing', 'Processed', 'Rejected') DEFAULT 'Not_Claimed',
    
    -- Advance Tax to be Paid
    advance_tax_paid DECIMAL(15, 2) DEFAULT 0,
    
    -- GST Return Filing (GSTR)
    gstr_1_filed BOOLEAN DEFAULT false, -- Outward Supplies
    gstr_2_filed BOOLEAN DEFAULT false, -- Inward Supplies (Applicable upto 30th Sep 2020)
    gstr_3_filed BOOLEAN DEFAULT false, -- Returns
    gstr_4_filed BOOLEAN DEFAULT false, -- Composition Scheme
    gstr_5_filed BOOLEAN DEFAULT false, -- Non-resident Taxable Person
    gstr_6_filed BOOLEAN DEFAULT false, -- Input Service Distributor
    gstr_7_filed BOOLEAN DEFAULT false, -- Tax Deductor
    gstr_8_filed BOOLEAN DEFAULT false, -- E-commerce Operator
    gstr_9_filed BOOLEAN DEFAULT false, -- Annual Return
    gstr_10_filed BOOLEAN DEFAULT false, -- Cancelled Registration
    
    -- GSTR Filing Dates
    gstr_filing_date TIMESTAMP,
    gstr_filing_status ENUM('Not_Filed', 'Filed', 'Accepted', 'Rejected', 'Processed') DEFAULT 'Not_Filed',
    
    -- Reconciliation
    gstr_1_gstr_3_reconciled BOOLEAN DEFAULT false,
    reconciliation_notes TEXT,
    
    -- Compliance
    is_compliant BOOLEAN DEFAULT true,
    non_compliance_reason TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_fy (tenant_id, fiscal_year),
    INDEX idx_return_period (month_year),
    INDEX idx_filing_status (gstr_filing_status),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- GST INVOICE TRACKING (for GSTR-1 & GSTR-3)
-- ============================================================================

CREATE TABLE IF NOT EXISTS gst_invoice_tracking (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Invoice Details
    invoice_id VARCHAR(36) NOT NULL,
    invoice_number VARCHAR(100),
    invoice_date TIMESTAMP,
    invoice_amount DECIMAL(15, 2),
    
    -- Customer Details
    customer_id VARCHAR(36),
    customer_gstin VARCHAR(50),
    
    -- GST Details
    gst_rate DECIMAL(5, 2),
    gst_amount DECIMAL(15, 2),
    
    -- Invoice Status
    invoice_raised_date TIMESTAMP,
    invoice_cancelled BOOLEAN DEFAULT false,
    cancellation_date TIMESTAMP,
    
    -- GSTR-1 Filing Status
    gstr_1_reported BOOLEAN DEFAULT false,
    gstr_1_filing_month INT,
    gstr_1_filing_year INT,
    
    -- ITC Status (Input Tax Credit)
    itc_eligible BOOLEAN DEFAULT true,
    itc_claimed BOOLEAN DEFAULT false,
    
    -- Reconciliation
    reconciled BOOLEAN DEFAULT false,
    reconciliation_date TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_invoice (tenant_id, invoice_id),
    INDEX idx_gstr_1_status (gstr_1_reported),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (invoice_id) REFERENCES sales_invoices(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- GST INPUT CREDIT TRACKING
-- ============================================================================

CREATE TABLE IF NOT EXISTS gst_input_credit (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Reference
    purchase_invoice_id VARCHAR(36),
    vendor_id VARCHAR(36),
    vendor_gstin VARCHAR(50),
    
    -- Invoice Details
    invoice_number VARCHAR(100),
    invoice_date TIMESTAMP,
    invoice_amount DECIMAL(15, 2),
    
    -- GST Details
    gst_rate DECIMAL(5, 2),
    gst_amount DECIMAL(15, 2),
    
    -- ITC Eligibility
    itc_eligible BOOLEAN DEFAULT true,
    itc_ineligible_reason TEXT, -- e.g., exempt supplies, personal consumption
    
    -- ITC Status
    itc_claimed BOOLEAN DEFAULT false,
    itc_claim_month INT,
    itc_claim_year INT,
    
    -- Blocked Credit (if any)
    blocked_percentage DECIMAL(5, 2) DEFAULT 0, -- e.g., personal use 50%
    blocked_amount DECIMAL(15, 2) DEFAULT 0,
    
    -- GSTR-2 Filing Status
    gstr_2_reported BOOLEAN DEFAULT false,
    gstr_2_filing_month INT,
    gstr_2_filing_year INT,
    
    -- Reconciliation with Vendor GSTR-1
    vendor_gstr_1_reconciled BOOLEAN DEFAULT false,
    reconciliation_date TIMESTAMP,
    discrepancy_found BOOLEAN DEFAULT false,
    discrepancy_notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_vendor (tenant_id, vendor_id),
    INDEX idx_itc_claim_status (itc_claimed),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (vendor_id) REFERENCES vendors(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- ADVANCE TAX (QUARTERLY PAYMENT)
-- ============================================================================

CREATE TABLE IF NOT EXISTS advance_tax_schedule (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Financial Year
    fiscal_year INT NOT NULL,
    assessment_year INT NOT NULL,
    
    -- Quarterly Dues
    q1_due_date TIMESTAMP, -- 15th June
    q1_amount_due DECIMAL(15, 2) DEFAULT 0,
    q1_amount_paid DECIMAL(15, 2) DEFAULT 0,
    q1_payment_date TIMESTAMP,
    q1_challan_number VARCHAR(50),
    q1_status ENUM('Not_Due', 'Pending', 'Paid', 'Overdue') DEFAULT 'Not_Due',
    
    q2_due_date TIMESTAMP, -- 15th September
    q2_amount_due DECIMAL(15, 2) DEFAULT 0,
    q2_amount_paid DECIMAL(15, 2) DEFAULT 0,
    q2_payment_date TIMESTAMP,
    q2_challan_number VARCHAR(50),
    q2_status ENUM('Not_Due', 'Pending', 'Paid', 'Overdue') DEFAULT 'Not_Due',
    
    q3_due_date TIMESTAMP, -- 15th December
    q3_amount_due DECIMAL(15, 2) DEFAULT 0,
    q3_amount_paid DECIMAL(15, 2) DEFAULT 0,
    q3_payment_date TIMESTAMP,
    q3_challan_number VARCHAR(50),
    q3_status ENUM('Not_Due', 'Pending', 'Paid', 'Overdue') DEFAULT 'Not_Due',
    
    q4_due_date TIMESTAMP, -- 15th March
    q4_amount_due DECIMAL(15, 2) DEFAULT 0,
    q4_amount_paid DECIMAL(15, 2) DEFAULT 0,
    q4_payment_date TIMESTAMP,
    q4_challan_number VARCHAR(50),
    q4_status ENUM('Not_Due', 'Pending', 'Paid', 'Overdue') DEFAULT 'Not_Due',
    
    -- Summary
    total_advance_tax_due DECIMAL(15, 2) DEFAULT 0,
    total_advance_tax_paid DECIMAL(15, 2) DEFAULT 0,
    
    is_compliant BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_fy (tenant_id, fiscal_year),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- TAX AUDIT & DOCUMENTATION
-- ============================================================================

CREATE TABLE IF NOT EXISTS tax_audit_trail (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Audit Details
    audit_type ENUM('Income_Tax', 'GST', 'TDS', 'Statutory', 'Internal') NOT NULL,
    audit_date TIMESTAMP,
    audit_period VARCHAR(50), -- e.g., "2024-2025"
    
    -- Audit Information
    auditor_name VARCHAR(255),
    auditor_firm VARCHAR(255),
    audit_report_reference VARCHAR(100),
    
    -- Findings
    audit_report_path VARCHAR(500),
    total_adjustments DECIMAL(15, 2),
    recommendations TEXT,
    
    -- Status
    audit_status ENUM('In_Progress', 'Completed', 'Action_Pending', 'Resolved') DEFAULT 'In_Progress',
    follow_up_required BOOLEAN DEFAULT false,
    follow_up_date TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_audit_type (tenant_id, audit_type),
    INDEX idx_audit_status (audit_status),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- TAX COMPLIANCE DOCUMENTS
-- ============================================================================

CREATE TABLE IF NOT EXISTS tax_compliance_documents (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Document Type
    document_type ENUM(
        'ITR_Form', 'ITR_Acknowledgement', 'TDS_Challan', 'GSTR_1', 'GSTR_3',
        'GSTR_9', 'TDS_Return', 'Audit_Report', 'Certificate', 'Other'
    ) NOT NULL,
    
    -- Document Details
    document_name VARCHAR(255) NOT NULL,
    document_date TIMESTAMP,
    reference_number VARCHAR(100),
    
    -- Financial Year / Period
    fiscal_year INT,
    filing_period VARCHAR(50),
    
    -- File Details
    document_path VARCHAR(500),
    file_size INT,
    file_hash VARCHAR(255),
    
    -- Verification
    verified BOOLEAN DEFAULT false,
    verified_by VARCHAR(36),
    verified_date TIMESTAMP,
    
    -- Archival
    archived BOOLEAN DEFAULT false,
    archive_date TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_doc_type (tenant_id, document_type),
    INDEX idx_fiscal_year (fiscal_year),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- COMPLIANCE CHECKLIST & DEADLINES
-- ============================================================================

CREATE TABLE IF NOT EXISTS tax_compliance_checklist (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Compliance Item
    compliance_item VARCHAR(255) NOT NULL,
    compliance_type ENUM('Income_Tax', 'GST', 'TDS', 'Other') NOT NULL,
    
    -- Timeline
    due_date TIMESTAMP NOT NULL,
    estimated_completion_date TIMESTAMP,
    actual_completion_date TIMESTAMP,
    
    -- Status
    status ENUM('Not_Started', 'In_Progress', 'Completed', 'Overdue', 'N/A') DEFAULT 'Not_Started',
    
    -- Responsibility
    assigned_to VARCHAR(36),
    
    -- Details
    description TEXT,
    documents_required TEXT,
    remarks TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_type (tenant_id, compliance_type),
    INDEX idx_due_date (due_date),
    INDEX idx_status (status),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- CREATE INDEXES FOR PERFORMANCE
-- ============================================================================

CREATE INDEX idx_income_tax_fy ON income_tax_compliance(tenant_id, fiscal_year);
CREATE INDEX idx_income_tax_filing ON income_tax_compliance(return_filed_status);
CREATE INDEX idx_tds_compliance_section ON tds_compliance(tenant_id, tds_section);
CREATE INDEX idx_gst_compliance_period ON gst_compliance(tenant_id, month_year);
CREATE INDEX idx_gst_invoice_gstr ON gst_invoice_tracking(tenant_id, gstr_1_reported);
CREATE INDEX idx_gst_credit_claim ON gst_input_credit(tenant_id, itc_claimed);
CREATE INDEX idx_advance_tax_status ON advance_tax_schedule(tenant_id, q1_status, q2_status, q3_status, q4_status);
