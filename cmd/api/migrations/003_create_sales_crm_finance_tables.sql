-- Phase 4: Sales, CRM, and Finance Module Tables
-- Migration: 003_create_sales_crm_finance_tables.sql

-- ===== SALES & PRESALES MODULE =====

CREATE TABLE IF NOT EXISTS leads (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    assigned_to_id VARCHAR(50),
    assigned_to_name VARCHAR(255),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    mobile_phone VARCHAR(20),
    website VARCHAR(255),
    industry VARCHAR(100),
    company_size VARCHAR(50),
    location TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    source VARCHAR(50) NOT NULL, -- website, referral, cold_call, email, trade, partner, social, advertising
    status VARCHAR(50) DEFAULT 'new', -- new, qualified, nurturing, converted, lost, on_hold
    budget DECIMAL(15,2),
    currency VARCHAR(3) DEFAULT 'INR',
    description TEXT,
    rating INTEGER DEFAULT 0,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    converted_at TIMESTAMP WITH TIME ZONE,
    last_contacted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_leads_tenant_id ON leads(tenant_id);
CREATE INDEX idx_leads_assigned_to_id ON leads(assigned_to_id);
CREATE INDEX idx_leads_status ON leads(status);
CREATE INDEX idx_leads_source ON leads(source);
CREATE INDEX idx_leads_created_at ON leads(created_at DESC);

CREATE TABLE IF NOT EXISTS opportunities (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    lead_id VARCHAR(50) NOT NULL,
    account_id VARCHAR(50),
    assigned_to_id VARCHAR(50),
    assigned_to_name VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    stage VARCHAR(50) DEFAULT 'lead', -- lead, qualified, demonstration, proposal, negotiation, closed, closed_won, closed_lost
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'INR',
    close_date DATE NOT NULL,
    expected_revenue DECIMAL(15,2),
    probability INTEGER DEFAULT 50, -- 0-100
    source VARCHAR(50),
    competitor_info TEXT,
    next_action VARCHAR(255),
    next_action_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    won_at TIMESTAMP WITH TIME ZONE,
    lost_at TIMESTAMP WITH TIME ZONE,
    lost_reason TEXT,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (lead_id) REFERENCES leads(id) ON DELETE CASCADE
);

CREATE INDEX idx_opportunities_tenant_id ON opportunities(tenant_id);
CREATE INDEX idx_opportunities_lead_id ON opportunities(lead_id);
CREATE INDEX idx_opportunities_stage ON opportunities(stage);
CREATE INDEX idx_opportunities_assigned_to_id ON opportunities(assigned_to_id);
CREATE INDEX idx_opportunities_close_date ON opportunities(close_date);

CREATE TABLE IF NOT EXISTS activities (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    lead_id VARCHAR(50),
    opportunity_id VARCHAR(50),
    account_id VARCHAR(50),
    activity_type VARCHAR(50) NOT NULL, -- call, email, meeting, task
    subject VARCHAR(255) NOT NULL,
    description TEXT,
    created_by_id VARCHAR(50) NOT NULL,
    created_by_name VARCHAR(255),
    assigned_to_id VARCHAR(50),
    assigned_to_name VARCHAR(255),
    activity_date TIMESTAMP WITH TIME ZONE NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'pending', -- completed, pending, cancelled
    priority VARCHAR(20) DEFAULT 'medium', -- high, medium, low
    duration INTEGER, -- in minutes
    outcome TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (lead_id) REFERENCES leads(id) ON DELETE SET NULL,
    FOREIGN KEY (opportunity_id) REFERENCES opportunities(id) ON DELETE SET NULL
);

CREATE INDEX idx_activities_tenant_id ON activities(tenant_id);
CREATE INDEX idx_activities_lead_id ON activities(lead_id);
CREATE INDEX idx_activities_opportunity_id ON activities(opportunity_id);
CREATE INDEX idx_activities_assigned_to_id ON activities(assigned_to_id);
CREATE INDEX idx_activities_activity_date ON activities(activity_date);

-- ===== CRM MODULE =====

CREATE TABLE IF NOT EXISTS accounts (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    website VARCHAR(255),
    industry VARCHAR(100),
    company_size VARCHAR(50), -- small, medium, large, enterprise
    billing_address TEXT,
    shipping_address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    account_manager_id VARCHAR(50),
    annual_revenue DECIMAL(15,2),
    employees INTEGER,
    rating INTEGER DEFAULT 0,
    type VARCHAR(50) DEFAULT 'prospect', -- prospect, customer, partner
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_accounts_tenant_id ON accounts(tenant_id);
CREATE INDEX idx_accounts_type ON accounts(type);
CREATE INDEX idx_accounts_account_manager_id ON accounts(account_manager_id);

CREATE TABLE IF NOT EXISTS contacts (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    account_id VARCHAR(50) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    title VARCHAR(100),
    department VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20),
    mobile_phone VARCHAR(20),
    role VARCHAR(50), -- Decision maker, influencer, user
    is_primary BOOLEAN DEFAULT FALSE,
    linkedin VARCHAR(255),
    twitter VARCHAR(255),
    preferred_language VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE INDEX idx_contacts_tenant_id ON contacts(tenant_id);
CREATE INDEX idx_contacts_account_id ON contacts(account_id);
CREATE INDEX idx_contacts_email ON contacts(email);

CREATE TABLE IF NOT EXISTS interactions (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    account_id VARCHAR(50) NOT NULL,
    contact_id VARCHAR(50),
    interaction_type VARCHAR(50) NOT NULL, -- support, feedback, complaint
    channel VARCHAR(50) NOT NULL, -- email, phone, chat, social
    subject VARCHAR(255) NOT NULL,
    message TEXT,
    status VARCHAR(50) DEFAULT 'open', -- open, resolved, closed
    priority VARCHAR(50) DEFAULT 'medium', -- high, medium, low
    assigned_to_id VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE SET NULL
);

CREATE INDEX idx_interactions_tenant_id ON interactions(tenant_id);
CREATE INDEX idx_interactions_account_id ON interactions(account_id);
CREATE INDEX idx_interactions_status ON interactions(status);
CREATE INDEX idx_interactions_created_at ON interactions(created_at DESC);

-- ===== FINANCE/ACCOUNTING MODULE =====

CREATE TABLE IF NOT EXISTS chart_of_accounts (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    account_code VARCHAR(50) NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) NOT NULL, -- asset, liability, equity, revenue, expense
    sub_type VARCHAR(100),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    opening_balance DECIMAL(15,2) DEFAULT 0,
    current_balance DECIMAL(15,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'INR',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE(tenant_id, account_code)
);

CREATE INDEX idx_chart_of_accounts_tenant_id ON chart_of_accounts(tenant_id);
CREATE INDEX idx_chart_of_accounts_type ON chart_of_accounts(account_type);
CREATE INDEX idx_chart_of_accounts_code ON chart_of_accounts(account_code);

CREATE TABLE IF NOT EXISTS journal_entries (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    entry_number VARCHAR(50) NOT NULL,
    entry_date DATE NOT NULL,
    reference VARCHAR(255),
    description TEXT,
    status VARCHAR(50) DEFAULT 'draft', -- draft, posted, canceled
    prepared_by_id VARCHAR(50) NOT NULL,
    approved_by_id VARCHAR(50),
    total_debit DECIMAL(15,2) DEFAULT 0,
    total_credit DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    posted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE(tenant_id, entry_number)
);

CREATE INDEX idx_journal_entries_tenant_id ON journal_entries(tenant_id);
CREATE INDEX idx_journal_entries_entry_date ON journal_entries(entry_date);
CREATE INDEX idx_journal_entries_status ON journal_entries(status);
CREATE INDEX idx_journal_entries_entry_number ON journal_entries(entry_number);

CREATE TABLE IF NOT EXISTS journal_entry_lines (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    journal_entry_id VARCHAR(50) NOT NULL,
    account_id VARCHAR(50) NOT NULL,
    account_code VARCHAR(50) NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    debit_amount DECIMAL(15,2) DEFAULT 0,
    credit_amount DECIMAL(15,2) DEFAULT 0,
    description TEXT,
    reference VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (journal_entry_id) REFERENCES journal_entries(id) ON DELETE CASCADE,
    FOREIGN KEY (account_id) REFERENCES chart_of_accounts(id) ON DELETE RESTRICT
);

CREATE INDEX idx_journal_entry_lines_journal_entry_id ON journal_entry_lines(journal_entry_id);
CREATE INDEX idx_journal_entry_lines_account_id ON journal_entry_lines(account_id);

CREATE TABLE IF NOT EXISTS general_ledger (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    account_id VARCHAR(50) NOT NULL,
    account_code VARCHAR(50) NOT NULL,
    journal_entry_id VARCHAR(50),
    entry_date DATE NOT NULL,
    reference VARCHAR(255),
    debit_amount DECIMAL(15,2) DEFAULT 0,
    credit_amount DECIMAL(15,2) DEFAULT 0,
    running_balance DECIMAL(15,2),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (account_id) REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,
    FOREIGN KEY (journal_entry_id) REFERENCES journal_entries(id) ON DELETE SET NULL
);

CREATE INDEX idx_general_ledger_tenant_id ON general_ledger(tenant_id);
CREATE INDEX idx_general_ledger_account_id ON general_ledger(account_id);
CREATE INDEX idx_general_ledger_entry_date ON general_ledger(entry_date);
CREATE INDEX idx_general_ledger_journal_entry_id ON general_ledger(journal_entry_id);
CREATE INDEX idx_general_ledger_composite ON general_ledger(account_id, entry_date DESC);

CREATE TABLE IF NOT EXISTS accounts_receivable (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    customer_id VARCHAR(50) NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    invoice_id VARCHAR(50),
    invoice_number VARCHAR(50) NOT NULL,
    invoice_amount DECIMAL(15,2) NOT NULL,
    amount_paid DECIMAL(15,2) DEFAULT 0,
    outstanding_amount DECIMAL(15,2),
    invoice_date DATE NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'invoiced', -- invoiced, partial_paid, paid, overdue, written_off
    currency VARCHAR(3) DEFAULT 'INR',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE(tenant_id, invoice_number)
);

CREATE INDEX idx_accounts_receivable_tenant_id ON accounts_receivable(tenant_id);
CREATE INDEX idx_accounts_receivable_status ON accounts_receivable(status);
CREATE INDEX idx_accounts_receivable_due_date ON accounts_receivable(due_date);

CREATE TABLE IF NOT EXISTS accounts_payable (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    vendor_id VARCHAR(50) NOT NULL,
    vendor_name VARCHAR(255) NOT NULL,
    bill_id VARCHAR(50),
    bill_number VARCHAR(50) NOT NULL,
    bill_amount DECIMAL(15,2) NOT NULL,
    amount_paid DECIMAL(15,2) DEFAULT 0,
    outstanding_amount DECIMAL(15,2),
    bill_date DATE NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'received', -- received, partial_paid, paid, overdue, written_off
    currency VARCHAR(3) DEFAULT 'INR',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE(tenant_id, bill_number)
);

CREATE INDEX idx_accounts_payable_tenant_id ON accounts_payable(tenant_id);
CREATE INDEX idx_accounts_payable_status ON accounts_payable(status);
CREATE INDEX idx_accounts_payable_due_date ON accounts_payable(due_date);

-- ===== TRIGGERS FOR AUTOMATIC TIMESTAMPS =====

CREATE OR REPLACE FUNCTION update_timestamp() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER leads_update_timestamp BEFORE UPDATE ON leads FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER opportunities_update_timestamp BEFORE UPDATE ON opportunities FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER activities_update_timestamp BEFORE UPDATE ON activities FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER accounts_update_timestamp BEFORE UPDATE ON accounts FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER contacts_update_timestamp BEFORE UPDATE ON contacts FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER interactions_update_timestamp BEFORE UPDATE ON interactions FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER chart_of_accounts_update_timestamp BEFORE UPDATE ON chart_of_accounts FOR EACH ROW EXECUTE FUNCTION update_timestamp();
CREATE TRIGGER journal_entries_update_timestamp BEFORE UPDATE ON journal_entries FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- ===== VIEWS FOR COMMON QUERIES =====

CREATE OR REPLACE VIEW v_sales_pipeline AS
SELECT 
    o.id,
    o.tenant_id,
    o.name,
    l.first_name || ' ' || l.last_name as lead_name,
    o.stage,
    o.amount,
    o.expected_revenue,
    o.probability,
    o.assigned_to_name,
    o.close_date,
    o.created_at
FROM opportunities o
LEFT JOIN leads l ON o.lead_id = l.id
WHERE o.stage NOT IN ('closed_won', 'closed_lost')
ORDER BY o.close_date;

CREATE OR REPLACE VIEW v_accounts_aging AS
SELECT 
    ar.customer_id,
    ar.customer_name,
    SUM(CASE WHEN ar.due_date < CURRENT_DATE THEN ar.outstanding_amount ELSE 0 END) as overdue_amount,
    SUM(CASE WHEN ar.due_date >= CURRENT_DATE THEN ar.outstanding_amount ELSE 0 END) as current_amount,
    SUM(ar.outstanding_amount) as total_outstanding,
    COUNT(DISTINCT ar.invoice_number) as invoice_count
FROM accounts_receivable ar
WHERE ar.status != 'paid'
GROUP BY ar.customer_id, ar.customer_name;

CREATE OR REPLACE VIEW v_financial_summary AS
SELECT 
    coa.tenant_id,
    coa.account_type,
    COUNT(DISTINCT coa.id) as account_count,
    SUM(coa.current_balance) as total_balance
FROM chart_of_accounts coa
WHERE coa.is_active = true
GROUP BY coa.tenant_id, coa.account_type;
