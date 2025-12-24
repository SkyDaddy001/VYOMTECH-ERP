-- Migration: Create Tenant Collection Tables
-- This migration creates tables for tenant-to-client payment collections
-- Tables: tenant_accounts, client_invoices, client_payments

-- Create tenant_accounts table
CREATE TABLE IF NOT EXISTS tenant_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    charge_type VARCHAR(50) NOT NULL,
    charge_type_name VARCHAR(100) NOT NULL,
    description TEXT,
    razorpay_account_id VARCHAR(255),
    billdesk_account_id VARCHAR(255),
    bank_account_name VARCHAR(255),
    bank_account_no VARCHAR(50),
    ifsc_code VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    total_collected DECIMAL(15, 2) DEFAULT 0,
    total_refunded DECIMAL(15, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_tenant_charge_type UNIQUE(tenant_id, charge_type)
);

-- Create indexes for tenant_accounts
CREATE INDEX idx_tenant_accounts_tenant_id ON tenant_accounts(tenant_id);
CREATE INDEX idx_tenant_accounts_charge_type ON tenant_accounts(charge_type);
CREATE INDEX idx_tenant_accounts_is_active ON tenant_accounts(is_active);

-- Create client_invoices table
CREATE TABLE IF NOT EXISTS client_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id VARCHAR(100) NOT NULL,
    client_name VARCHAR(255) NOT NULL,
    client_email VARCHAR(255),
    client_phone VARCHAR(20),
    charge_type VARCHAR(50) NOT NULL,
    invoice_number VARCHAR(100) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    amount_paid DECIMAL(15, 2) DEFAULT 0,
    outstanding_amount DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(5) DEFAULT 'INR',
    description TEXT,
    invoice_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP,
    status VARCHAR(50) DEFAULT 'issued',
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT check_amount_positive CHECK(amount > 0),
    CONSTRAINT check_amount_paid_non_negative CHECK(amount_paid >= 0),
    CONSTRAINT check_outstanding_amount_non_negative CHECK(outstanding_amount >= 0),
    CONSTRAINT unique_invoice_number UNIQUE(tenant_id, invoice_number)
);

-- Create indexes for client_invoices
CREATE INDEX idx_client_invoices_tenant_id ON client_invoices(tenant_id);
CREATE INDEX idx_client_invoices_client_id ON client_invoices(client_id);
CREATE INDEX idx_client_invoices_charge_type ON client_invoices(charge_type);
CREATE INDEX idx_client_invoices_status ON client_invoices(status);
CREATE INDEX idx_client_invoices_due_date ON client_invoices(due_date);
CREATE INDEX idx_client_invoices_tenant_client ON client_invoices(tenant_id, client_id);

-- Create client_payments table
CREATE TABLE IF NOT EXISTS client_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id VARCHAR(100) NOT NULL,
    invoice_id UUID REFERENCES client_invoices(id) ON DELETE SET NULL,
    tenant_account_id UUID REFERENCES tenant_accounts(id) ON DELETE RESTRICT,
    charge_type VARCHAR(50) NOT NULL,
    order_id VARCHAR(100) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(5) DEFAULT 'INR',
    status VARCHAR(50) DEFAULT 'created',
    payment_type VARCHAR(50) DEFAULT 'client',
    provider VARCHAR(50) NOT NULL,
    payment_method VARCHAR(50),
    gateway_order_id VARCHAR(255),
    gateway_payment_id VARCHAR(255),
    transaction_id VARCHAR(255),
    client_name VARCHAR(255) NOT NULL,
    client_email VARCHAR(255),
    client_phone VARCHAR(20),
    receipt_url TEXT,
    refund_id VARCHAR(255),
    refund_amount DECIMAL(15, 2) DEFAULT 0,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    
    CONSTRAINT check_amount_gt_zero CHECK(amount > 0),
    CONSTRAINT unique_gateway_payment UNIQUE(provider, gateway_payment_id)
);

-- Create indexes for client_payments
CREATE INDEX idx_client_payments_tenant_id ON client_payments(tenant_id);
CREATE INDEX idx_client_payments_client_id ON client_payments(client_id);
CREATE INDEX idx_client_payments_invoice_id ON client_payments(invoice_id);
CREATE INDEX idx_client_payments_status ON client_payments(status);
CREATE INDEX idx_client_payments_provider ON client_payments(provider);
CREATE INDEX idx_client_payments_charge_type ON client_payments(charge_type);
CREATE INDEX idx_client_payments_created_at ON client_payments(created_at);
CREATE INDEX idx_client_payments_tenant_charge ON client_payments(tenant_id, charge_type);
CREATE INDEX idx_client_payments_gateway_id ON client_payments(gateway_payment_id);

-- Create trigger for updated_at on tenant_accounts
CREATE OR REPLACE FUNCTION update_tenant_accounts_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tenant_accounts_updated_at
BEFORE UPDATE ON tenant_accounts
FOR EACH ROW
EXECUTE FUNCTION update_tenant_accounts_timestamp();

-- Create trigger for updated_at on client_invoices
CREATE OR REPLACE FUNCTION update_client_invoices_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER client_invoices_updated_at
BEFORE UPDATE ON client_invoices
FOR EACH ROW
EXECUTE FUNCTION update_client_invoices_timestamp();

-- Create trigger for updated_at on client_payments
CREATE OR REPLACE FUNCTION update_client_payments_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER client_payments_updated_at
BEFORE UPDATE ON client_payments
FOR EACH ROW
EXECUTE FUNCTION update_client_payments_timestamp();

-- Create view for invoice summary
CREATE OR REPLACE VIEW v_invoice_summary AS
SELECT 
    ci.id,
    ci.tenant_id,
    ci.client_id,
    ci.invoice_number,
    ci.amount,
    ci.amount_paid,
    ci.outstanding_amount,
    ci.charge_type,
    ci.status,
    ci.due_date,
    CASE WHEN ci.due_date < CURRENT_TIMESTAMP AND ci.status != 'paid' THEN true ELSE false END as is_overdue,
    COUNT(cp.id) as payment_count,
    COALESCE(SUM(cp.amount), 0) as total_paid
FROM client_invoices ci
LEFT JOIN client_payments cp ON ci.id = cp.invoice_id AND cp.status = 'completed'
GROUP BY ci.id, ci.tenant_id, ci.client_id, ci.invoice_number, ci.amount, ci.amount_paid, 
         ci.outstanding_amount, ci.charge_type, ci.status, ci.due_date;

-- Create view for tenant collection summary
CREATE OR REPLACE VIEW v_tenant_collection_summary AS
SELECT 
    ta.tenant_id,
    ta.id as account_id,
    ta.charge_type,
    COUNT(DISTINCT ci.id) as total_invoices,
    COUNT(DISTINCT CASE WHEN ci.status = 'paid' THEN ci.id END) as paid_invoices,
    COUNT(DISTINCT CASE WHEN ci.status IN ('issued', 'partial_paid', 'overdue') THEN ci.id END) as outstanding_invoices,
    COALESCE(SUM(ci.amount), 0) as total_billed,
    COALESCE(SUM(ci.amount_paid), 0) as total_collected,
    COALESCE(SUM(ci.outstanding_amount), 0) as outstanding_amount,
    CASE 
        WHEN COALESCE(SUM(ci.amount), 0) > 0 
        THEN ROUND((COALESCE(SUM(ci.amount_paid), 0) / COALESCE(SUM(ci.amount), 0) * 100)::numeric, 2)
        ELSE 0
    END as collection_rate
FROM tenant_accounts ta
LEFT JOIN client_invoices ci ON ta.tenant_id = ci.tenant_id AND ta.charge_type = ci.charge_type
WHERE ta.is_active = true
GROUP BY ta.tenant_id, ta.id, ta.charge_type;

-- Create view for client outstanding summary
CREATE OR REPLACE VIEW v_client_outstanding_summary AS
SELECT 
    ci.tenant_id,
    ci.client_id,
    ci.client_name,
    ci.client_email,
    SUM(ci.amount) as total_billed,
    SUM(ci.amount_paid) as total_paid,
    SUM(ci.outstanding_amount) as total_outstanding,
    COUNT(DISTINCT ci.id) as total_invoices,
    COUNT(DISTINCT CASE WHEN ci.status = 'paid' THEN ci.id END) as paid_invoices,
    COUNT(DISTINCT CASE WHEN ci.due_date < CURRENT_TIMESTAMP AND ci.status != 'paid' THEN ci.id END) as overdue_invoices
FROM client_invoices ci
WHERE ci.outstanding_amount > 0
GROUP BY ci.tenant_id, ci.client_id, ci.client_name, ci.client_email;
