-- ============================================================================
-- PROJECT COLLECTION ACCOUNTS & RERA COMPLIANCE
-- ============================================================================
-- RERA Regulations require:
-- 1. Separate bank account for each project
-- 2. Segregated fund management
-- 3. Collection tracking against project units
-- 4. Borrowing restrictions
-- 5. Interest calculation on delayed collections

CREATE TABLE IF NOT EXISTS project_collection_accounts (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    project_id VARCHAR(36) NOT NULL,
    
    -- Account Details
    account_name VARCHAR(255) NOT NULL,
    account_code VARCHAR(50) UNIQUE NOT NULL,
    bank_name VARCHAR(255) NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    ifsc_code VARCHAR(20) NOT NULL,
    account_type ENUM('Current', 'Savings') DEFAULT 'Current',
    
    -- Opening Balance & Limits
    opening_balance DECIMAL(15, 2) DEFAULT 0,
    current_balance DECIMAL(15, 2) DEFAULT 0,
    minimum_balance DECIMAL(15, 2) DEFAULT 0,
    
    -- RERA Compliance
    rera_compliant BOOLEAN DEFAULT true,
    regulated_account BOOLEAN DEFAULT true, -- Must be segregated
    max_borrowing_allowed DECIMAL(15, 2), -- Max loan against account (10% of collections)
    current_borrowing DECIMAL(15, 2) DEFAULT 0,
    
    -- Status & Dates
    status ENUM('Active', 'Inactive', 'Closed') DEFAULT 'Active',
    opened_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    closed_date TIMESTAMP NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by VARCHAR(36),
    
    INDEX idx_tenant_project (tenant_id, project_id),
    INDEX idx_account_code (account_code),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (project_id) REFERENCES property_projects(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- PROJECT COLLECTION LEDGER
-- ============================================================================
-- Tracks all collections received for a project segregated by unit

CREATE TABLE IF NOT EXISTS project_collection_ledger (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    project_id VARCHAR(36) NOT NULL,
    collection_account_id VARCHAR(36) NOT NULL,
    
    -- Collection Details
    collection_date TIMESTAMP NOT NULL,
    collection_number VARCHAR(50) UNIQUE NOT NULL,
    booking_id VARCHAR(36),
    customer_id VARCHAR(36),
    unit_id VARCHAR(36),
    
    -- Payment Information
    payment_mode ENUM('Cash', 'Cheque', 'NEFT', 'RTGS', 'IFT', 'Demand_Draft') NOT NULL,
    amount_collected DECIMAL(15, 2) NOT NULL,
    
    -- Cheque Details (if applicable)
    cheque_number VARCHAR(50),
    cheque_date TIMESTAMP,
    cheque_status ENUM('Not_Applicable', 'Pending', 'Cleared', 'Bounced', 'Post_Dated'),
    
    -- Reference
    reference_number VARCHAR(100),
    paid_by VARCHAR(255),
    remarks TEXT,
    
    -- GL Integration
    gl_posted BOOLEAN DEFAULT false,
    journal_entry_id VARCHAR(36),
    posted_at TIMESTAMP NULL,
    posted_by VARCHAR(36),
    
    -- Status
    status ENUM('Pending', 'Approved', 'Reversed', 'Cancelled') DEFAULT 'Pending',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by VARCHAR(36),
    
    INDEX idx_tenant_project (tenant_id, project_id),
    INDEX idx_collection_account (collection_account_id),
    INDEX idx_booking_unit (booking_id, unit_id),
    INDEX idx_collection_date (collection_date),
    INDEX idx_cheque_status (cheque_status),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (project_id) REFERENCES property_projects(id),
    FOREIGN KEY (collection_account_id) REFERENCES project_collection_accounts(id),
    FOREIGN KEY (booking_id) REFERENCES customer_bookings(id),
    FOREIGN KEY (unit_id) REFERENCES property_units(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- COLLECTION AGAINST BOOKING/PAYMENT SCHEDULE
-- ============================================================================
-- Links collections to specific payment milestones

CREATE TABLE IF NOT EXISTS collection_against_milestone (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    project_id VARCHAR(36) NOT NULL,
    
    -- Reference
    collection_id VARCHAR(36) NOT NULL,
    booking_id VARCHAR(36) NOT NULL,
    payment_schedule_id VARCHAR(36) NOT NULL,
    
    -- Amount Mapping
    scheduled_amount DECIMAL(15, 2) NOT NULL,
    collected_amount DECIMAL(15, 2) NOT NULL,
    excess_amount DECIMAL(15, 2) DEFAULT 0,
    shortfall_amount DECIMAL(15, 2) DEFAULT 0,
    
    -- Dates
    scheduled_date TIMESTAMP,
    collection_date TIMESTAMP,
    days_delayed INT DEFAULT 0,
    
    -- Interest on Delayed Collection (if applicable)
    interest_applicable BOOLEAN DEFAULT false,
    interest_rate DECIMAL(5, 2) DEFAULT 0, -- % per month
    interest_amount DECIMAL(15, 2) DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_project (tenant_id, project_id),
    INDEX idx_collection_id (collection_id),
    INDEX idx_booking_id (booking_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (project_id) REFERENCES property_projects(id),
    FOREIGN KEY (collection_id) REFERENCES project_collection_ledger(id),
    FOREIGN KEY (booking_id) REFERENCES customer_bookings(id),
    FOREIGN KEY (payment_schedule_id) REFERENCES payment_schedules(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- PROJECT FUND UTILIZATION TRACKING
-- ============================================================================
-- RERA requires tracking how collection funds are used

CREATE TABLE IF NOT EXISTS project_fund_utilization (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    project_id VARCHAR(36) NOT NULL,
    collection_account_id VARCHAR(36) NOT NULL,
    
    -- Fund Utilization Details
    utilization_date TIMESTAMP NOT NULL,
    utilization_type ENUM('Construction', 'Land_Cost', 'Statutory_Approval', 'Admin', 'Interest', 'Other') NOT NULL,
    description TEXT,
    
    -- Amount
    amount_utilized DECIMAL(15, 2) NOT NULL,
    
    -- Supporting Documents
    bill_number VARCHAR(100),
    bill_date TIMESTAMP,
    invoice_id VARCHAR(36),
    bill_amount DECIMAL(15, 2),
    
    -- Approval
    approved_by VARCHAR(36),
    approval_date TIMESTAMP,
    approval_remarks TEXT,
    
    -- GL Integration
    gl_posted BOOLEAN DEFAULT false,
    journal_entry_id VARCHAR(36),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_project (tenant_id, project_id),
    INDEX idx_collection_account (collection_account_id),
    INDEX idx_utilization_type (utilization_type),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (project_id) REFERENCES property_projects(id),
    FOREIGN KEY (collection_account_id) REFERENCES project_collection_accounts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- BORROWING AGAINST PROJECT ACCOUNT (RERA - Limited to 10%)
-- ============================================================================

CREATE TABLE IF NOT EXISTS project_account_borrowings (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    project_id VARCHAR(36) NOT NULL,
    collection_account_id VARCHAR(36) NOT NULL,
    
    -- Loan Details
    loan_number VARCHAR(50) UNIQUE NOT NULL,
    loan_date TIMESTAMP NOT NULL,
    principal_amount DECIMAL(15, 2) NOT NULL,
    interest_rate DECIMAL(5, 2) NOT NULL,
    loan_tenure_months INT,
    
    -- Collections-Based Limits
    total_collections_till_date DECIMAL(15, 2) NOT NULL,
    max_borrowing_allowed DECIMAL(15, 2), -- 10% of total collections
    
    -- Repayment
    maturity_date TIMESTAMP,
    repaid_amount DECIMAL(15, 2) DEFAULT 0,
    outstanding_amount DECIMAL(15, 2),
    
    -- Status
    status ENUM('Active', 'Repaid', 'Overdue', 'Waived') DEFAULT 'Active',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_project (tenant_id, project_id),
    INDEX idx_collection_account (collection_account_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (project_id) REFERENCES property_projects(id),
    FOREIGN KEY (collection_account_id) REFERENCES project_collection_accounts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- PROJECT RECONCILIATION (Monthly/Quarterly)
-- ============================================================================

CREATE TABLE IF NOT EXISTS project_account_reconciliation (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    project_id VARCHAR(36) NOT NULL,
    collection_account_id VARCHAR(36) NOT NULL,
    
    -- Period
    reconciliation_date TIMESTAMP NOT NULL,
    period_from TIMESTAMP NOT NULL,
    period_to TIMESTAMP NOT NULL,
    
    -- Collections
    opening_balance DECIMAL(15, 2) NOT NULL,
    total_collections DECIMAL(15, 2) NOT NULL,
    total_utilizations DECIMAL(15, 2) NOT NULL,
    borrowings DECIMAL(15, 2) DEFAULT 0,
    interest_accrued DECIMAL(15, 2) DEFAULT 0,
    closing_balance DECIMAL(15, 2) NOT NULL,
    
    -- Bank Reconciliation
    bank_balance DECIMAL(15, 2),
    book_balance DECIMAL(15, 2),
    reconciled BOOLEAN DEFAULT false,
    reconciliation_variance DECIMAL(15, 2) DEFAULT 0,
    
    -- Approval
    verified_by VARCHAR(36),
    verified_date TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_project (tenant_id, project_id),
    INDEX idx_collection_account (collection_account_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (project_id) REFERENCES property_projects(id),
    FOREIGN KEY (collection_account_id) REFERENCES project_collection_accounts(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- CREATE INDEXES FOR PERFORMANCE
-- ============================================================================

CREATE INDEX idx_pca_tenant_status ON project_collection_accounts(tenant_id, status);
CREATE INDEX idx_pcl_date_status ON project_collection_ledger(collection_date, status);
CREATE INDEX idx_pfu_date_type ON project_fund_utilization(utilization_date, utilization_type);
CREATE INDEX idx_pab_status ON project_account_borrowings(status);
