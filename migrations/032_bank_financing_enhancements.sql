-- Migration 032: Bank Financing & Disbursement Management
-- Created: 2025-12-23
-- Purpose: Track bank loans, disbursements, NOC management for real estate projects

-- ============================================================================
-- TABLE: bank_financing
-- Description: Main bank financing records linked to bookings
-- ============================================================================
CREATE TABLE IF NOT EXISTS bank_financing (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    booking_id VARCHAR(36) NOT NULL,
    bank_id VARCHAR(36),
    
    -- Financing Details
    loan_amount DECIMAL(15, 2) NOT NULL,
    sanctioned_amount DECIMAL(15, 2) NOT NULL,
    disbursed_amount DECIMAL(15, 2) DEFAULT 0,
    outstanding_amount DECIMAL(15, 2) DEFAULT 0,
    
    -- Loan Details
    loan_type VARCHAR(50) NOT NULL, -- 'Home Loan', 'Construction Loan', 'Bridge Loan'
    interest_rate DECIMAL(5, 2),
    tenure_months INT,
    emi_amount DECIMAL(12, 2),
    
    -- Status & Dates
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, approved, sanctioned, disbursing, completed, rejected
    application_date DATETIME,
    approval_date DATETIME,
    sanction_date DATETIME,
    expected_completion_date DATETIME,
    
    -- Documents
    application_ref_no VARCHAR(100) UNIQUE,
    sanction_letter_url VARCHAR(500),
    
    -- Metadata
    created_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    KEY idx_tenant_id (tenant_id),
    KEY idx_booking_id (booking_id),
    KEY idx_bank_id (bank_id),
    KEY idx_status (status),
    KEY idx_application_ref (application_ref_no),
    CONSTRAINT fk_financing_booking FOREIGN KEY (booking_id) REFERENCES booking(id),
    CONSTRAINT fk_financing_bank FOREIGN KEY (bank_id) REFERENCES bank(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLE: bank_disbursement
-- Description: Tracks disbursement schedule and actual disbursements
-- ============================================================================
CREATE TABLE IF NOT EXISTS bank_disbursement (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    financing_id VARCHAR(36) NOT NULL,
    
    -- Disbursement Details
    disbursement_number INT NOT NULL,
    scheduled_amount DECIMAL(15, 2) NOT NULL,
    actual_amount DECIMAL(15, 2),
    
    -- Milestone Linked
    milestone_id VARCHAR(36), -- Construction milestone for which disbursement is scheduled
    milestone_percentage INT, -- e.g., 25%, 50%, 75%, 100%
    
    -- Status & Dates
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, released, credited, delayed, cancelled
    scheduled_date DATE NOT NULL,
    actual_date DATE,
    bank_reference_no VARCHAR(100),
    
    -- Documentation
    claim_document_url VARCHAR(500),
    release_approval_by VARCHAR(36),
    release_approval_date DATETIME,
    
    -- Metadata
    created_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    KEY idx_tenant_id (tenant_id),
    KEY idx_financing_id (financing_id),
    KEY idx_milestone_id (milestone_id),
    KEY idx_status (status),
    KEY idx_scheduled_date (scheduled_date),
    CONSTRAINT fk_disbursement_financing FOREIGN KEY (financing_id) REFERENCES bank_financing(id),
    CONSTRAINT fk_disbursement_milestone FOREIGN KEY (milestone_id) REFERENCES project_milestone(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLE: bank_noc
-- Description: No Objection Certificate (NOC) tracking from banks
-- ============================================================================
CREATE TABLE IF NOT EXISTS bank_noc (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    financing_id VARCHAR(36) NOT NULL,
    
    -- NOC Details
    noc_type VARCHAR(50) NOT NULL, -- 'Pre-sanction', 'Post-completion', 'Full-settlement'
    noc_request_date DATETIME NOT NULL,
    noc_received_date DATETIME,
    
    -- Documents
    noc_document_url VARCHAR(500),
    noc_amount DECIMAL(15, 2), -- Amount covered by NOC
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'requested', -- requested, issued, expired, cancelled
    issued_by_bank VARCHAR(100),
    valid_till_date DATE,
    
    -- Remarks
    remarks TEXT,
    
    -- Metadata
    created_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    KEY idx_tenant_id (tenant_id),
    KEY idx_financing_id (financing_id),
    KEY idx_status (status),
    KEY idx_noc_type (noc_type),
    CONSTRAINT fk_noc_financing FOREIGN KEY (financing_id) REFERENCES bank_financing(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLE: bank_collection_tracking
-- Description: Collection of loan amounts with EMI and prepayment tracking
-- ============================================================================
CREATE TABLE IF NOT EXISTS bank_collection_tracking (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    financing_id VARCHAR(36) NOT NULL,
    
    -- Collection Details
    collection_type VARCHAR(50) NOT NULL, -- 'EMI', 'Prepayment', 'Partial', 'Full-Settlement'
    collection_amount DECIMAL(15, 2) NOT NULL,
    collection_date DATE NOT NULL,
    
    -- Payment Mode
    payment_mode VARCHAR(50), -- 'Bank Transfer', 'Cheque', 'NEFT', 'RTGS'
    payment_reference_no VARCHAR(100),
    
    -- EMI Details (if applicable)
    emi_month VARCHAR(10), -- e.g., 'Jan-2025'
    emi_number INT,
    principal_amount DECIMAL(12, 2),
    interest_amount DECIMAL(12, 2),
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, verified, credited, failed
    bank_confirmation_date DATETIME,
    
    -- Metadata
    created_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    KEY idx_tenant_id (tenant_id),
    KEY idx_financing_id (financing_id),
    KEY idx_collection_date (collection_date),
    KEY idx_status (status),
    CONSTRAINT fk_collection_financing FOREIGN KEY (financing_id) REFERENCES bank_financing(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLE: bank
-- Description: Bank master list (if not already exists)
-- ============================================================================
CREATE TABLE IF NOT EXISTS bank (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    bank_name VARCHAR(100) NOT NULL,
    branch_name VARCHAR(100),
    ifsc_code VARCHAR(20),
    branch_contact VARCHAR(20),
    branch_email VARCHAR(100),
    
    -- Contact Person
    relationship_manager_name VARCHAR(100),
    relationship_manager_phone VARCHAR(20),
    relationship_manager_email VARCHAR(100),
    
    status VARCHAR(50) DEFAULT 'active',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    KEY idx_tenant_id (tenant_id),
    UNIQUE KEY uk_bank_branch (tenant_id, bank_name, branch_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- INDEXES & CONSTRAINTS
-- ============================================================================

-- Add indexes for common queries
CREATE INDEX idx_financing_created_at ON bank_financing(created_at DESC);
CREATE INDEX idx_financing_status_date ON bank_financing(status, sanction_date);
CREATE INDEX idx_disbursement_scheduled_date ON bank_disbursement(scheduled_date);
CREATE INDEX idx_collection_month ON bank_collection_tracking(emi_month, financing_id);

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================
