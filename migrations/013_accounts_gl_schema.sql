-- ============================================================================
-- ACCOUNTS (GENERAL LEDGER) MODULE SCHEMA
-- ============================================================================
-- This migration creates tables for complete accounting and GL management
-- All modules (HR, Sales, Purchase, Construction) post transactions here

-- Chart of Accounts
CREATE TABLE IF NOT EXISTS chart_of_accounts (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    account_code VARCHAR(50) UNIQUE NOT NULL,
    account_name VARCHAR(200) NOT NULL,
    account_type ENUM('Asset', 'Liability', 'Equity', 'Revenue', 'Expense', 'Other') NOT NULL,
    sub_account_type VARCHAR(100),
    parent_account_id VARCHAR(36),
    description VARCHAR(500),
    opening_balance DECIMAL(16, 2) DEFAULT 0,
    current_balance DECIMAL(16, 2) DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    is_header BOOLEAN DEFAULT FALSE,
    is_default BOOLEAN DEFAULT FALSE,
    currency VARCHAR(3) DEFAULT 'INR',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (parent_account_id) REFERENCES chart_of_accounts(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_account_code (account_code),
    INDEX idx_account_type (account_type),
    INDEX idx_parent_account_id (parent_account_id)
);

-- Journal Entries (transactions)
CREATE TABLE IF NOT EXISTS journal_entries (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    entry_date DATE NOT NULL,
    reference_number VARCHAR(100),
    reference_type ENUM('Manual', 'HR_Payroll', 'Sales_Invoice', 'Purchase_Invoice', 'Expense', 'Transfer', 'Other') NOT NULL,
    reference_id VARCHAR(36),
    description VARCHAR(500),
    amount DECIMAL(16, 2) NOT NULL,
    narration VARCHAR(1000),
    entry_status ENUM('Draft', 'Posted', 'Cancelled') DEFAULT 'Draft',
    posted_by VARCHAR(36),
    posted_at TIMESTAMP NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (posted_by) REFERENCES users(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_entry_date (entry_date),
    INDEX idx_reference_number (reference_number),
    INDEX idx_entry_status (entry_status)
);

-- Journal Entry Details (individual debit/credit lines)
CREATE TABLE IF NOT EXISTS journal_entry_details (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    journal_entry_id VARCHAR(36) NOT NULL,
    account_id VARCHAR(36) NOT NULL,
    account_code VARCHAR(50),
    debit_amount DECIMAL(16, 2) DEFAULT 0,
    credit_amount DECIMAL(16, 2) DEFAULT 0,
    description VARCHAR(500),
    line_number INT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (journal_entry_id) REFERENCES journal_entries(id),
    FOREIGN KEY (account_id) REFERENCES chart_of_accounts(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_journal_entry_id (journal_entry_id),
    INDEX idx_account_id (account_id)
);

-- GL Account Balance (cached for performance)
CREATE TABLE IF NOT EXISTS gl_account_balance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    account_id VARCHAR(36) NOT NULL,
    fiscal_period DATE NOT NULL,
    opening_balance DECIMAL(16, 2) DEFAULT 0,
    total_debit DECIMAL(16, 2) DEFAULT 0,
    total_credit DECIMAL(16, 2) DEFAULT 0,
    closing_balance DECIMAL(16, 2) DEFAULT 0,
    
    UNIQUE KEY unique_account_period (tenant_id, account_id, fiscal_period),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (account_id) REFERENCES chart_of_accounts(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_account_id (account_id),
    INDEX idx_fiscal_period (fiscal_period)
);

-- Financial Periods (Year/Month definitions)
CREATE TABLE IF NOT EXISTS financial_periods (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    period_name VARCHAR(100) NOT NULL,
    period_type ENUM('Monthly', 'Quarterly', 'Annual') DEFAULT 'Monthly',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_closed BOOLEAN DEFAULT FALSE,
    closed_by VARCHAR(36),
    closed_at TIMESTAMP NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (closed_by) REFERENCES users(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_start_date (start_date),
    INDEX idx_end_date (end_date)
);

-- Trial Balance (debit/credit summary per account per period)
CREATE TABLE IF NOT EXISTS trial_balance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    period_id VARCHAR(36) NOT NULL,
    account_id VARCHAR(36) NOT NULL,
    account_code VARCHAR(50),
    account_name VARCHAR(200),
    debit_balance DECIMAL(16, 2) DEFAULT 0,
    credit_balance DECIMAL(16, 2) DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (period_id) REFERENCES financial_periods(id),
    FOREIGN KEY (account_id) REFERENCES chart_of_accounts(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_period_id (period_id),
    INDEX idx_account_id (account_id)
);

-- Income Statement (P&L Report)
CREATE TABLE IF NOT EXISTS income_statement (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    period_id VARCHAR(36) NOT NULL,
    revenue_total DECIMAL(16, 2) DEFAULT 0,
    cost_of_goods_sold DECIMAL(16, 2) DEFAULT 0,
    gross_profit DECIMAL(16, 2) DEFAULT 0,
    operating_expenses DECIMAL(16, 2) DEFAULT 0,
    operating_income DECIMAL(16, 2) DEFAULT 0,
    other_income DECIMAL(16, 2) DEFAULT 0,
    other_expenses DECIMAL(16, 2) DEFAULT 0,
    income_before_tax DECIMAL(16, 2) DEFAULT 0,
    tax_expense DECIMAL(16, 2) DEFAULT 0,
    net_income DECIMAL(16, 2) DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (period_id) REFERENCES financial_periods(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_period_id (period_id)
);

-- Balance Sheet
CREATE TABLE IF NOT EXISTS balance_sheet (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    period_id VARCHAR(36) NOT NULL,
    
    -- Assets
    current_assets DECIMAL(16, 2) DEFAULT 0,
    fixed_assets DECIMAL(16, 2) DEFAULT 0,
    other_assets DECIMAL(16, 2) DEFAULT 0,
    total_assets DECIMAL(16, 2) DEFAULT 0,
    
    -- Liabilities
    current_liabilities DECIMAL(16, 2) DEFAULT 0,
    long_term_liabilities DECIMAL(16, 2) DEFAULT 0,
    other_liabilities DECIMAL(16, 2) DEFAULT 0,
    total_liabilities DECIMAL(16, 2) DEFAULT 0,
    
    -- Equity
    paid_up_capital DECIMAL(16, 2) DEFAULT 0,
    retained_earnings DECIMAL(16, 2) DEFAULT 0,
    total_equity DECIMAL(16, 2) DEFAULT 0,
    
    -- Validation: Assets = Liabilities + Equity
    is_balanced BOOLEAN DEFAULT FALSE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (period_id) REFERENCES financial_periods(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_period_id (period_id)
);

-- Reconciliation Records
CREATE TABLE IF NOT EXISTS reconciliations (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    account_id VARCHAR(36) NOT NULL,
    reconciliation_type ENUM('Bank', 'Cash', 'Receivables', 'Payables', 'Inventory', 'Other') NOT NULL,
    period_from DATE NOT NULL,
    period_to DATE NOT NULL,
    system_balance DECIMAL(16, 2),
    actual_balance DECIMAL(16, 2),
    difference DECIMAL(16, 2),
    reconciliation_status ENUM('Pending', 'In Progress', 'Completed', 'Discrepancy') DEFAULT 'Pending',
    reconciled_by VARCHAR(36),
    reconciled_at TIMESTAMP NULL,
    notes VARCHAR(1000),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (account_id) REFERENCES chart_of_accounts(id),
    FOREIGN KEY (reconciled_by) REFERENCES users(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_account_id (account_id),
    INDEX idx_reconciliation_status (reconciliation_status)
);

-- Audit Trail for GL transactions
CREATE TABLE IF NOT EXISTS gl_audit_log (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    entity_type VARCHAR(100),
    entity_id VARCHAR(36),
    action VARCHAR(50),
    old_values JSON,
    new_values JSON,
    changed_by VARCHAR(36),
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_entity_type (entity_type),
    INDEX idx_changed_at (changed_at)
);

-- Create indexes for better query performance
CREATE INDEX idx_je_tenant_date ON journal_entries(tenant_id, entry_date);
CREATE INDEX idx_jed_entry_account ON journal_entry_details(journal_entry_id, account_id);
CREATE INDEX idx_coa_tenant_active ON chart_of_accounts(tenant_id, is_active);
