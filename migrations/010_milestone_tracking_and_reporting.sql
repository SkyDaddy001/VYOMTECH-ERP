-- ============================================================================
-- MILESTONE TRACKING AND REPORTING SCHEMA
-- Migration: 010_milestone_tracking_and_reporting
-- Purpose: Add comprehensive milestone tracking and lead lifecycle tagging
-- ============================================================================

-- ============================================================================
-- 1. CAMPAIGNS TABLE (for campaign tagging)
-- ============================================================================
CREATE TABLE IF NOT EXISTS sales_campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    campaign_code VARCHAR(50) NOT NULL,
    campaign_name VARCHAR(255) NOT NULL,
    campaign_type VARCHAR(50) NOT NULL, -- email, social, referral, event, digital, traditional, direct, outbound
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE,
    budget DECIMAL(15,2),
    expected_roi DECIMAL(5,2),
    assigned_to UUID REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, inactive, completed, paused
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(tenant_id, campaign_code),
    CONSTRAINT valid_status CHECK (status IN ('active', 'inactive', 'completed', 'paused'))
);

CREATE INDEX idx_sales_campaigns_tenant ON sales_campaigns(tenant_id);
CREATE INDEX idx_sales_campaigns_status ON sales_campaigns(status);
CREATE INDEX idx_sales_campaigns_type ON sales_campaigns(campaign_type);

-- ============================================================================
-- 2. LEAD SOURCE AND SUBSOURCE TAGGING TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS sales_lead_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    source_code VARCHAR(50) NOT NULL,
    source_name VARCHAR(100) NOT NULL,
    source_type VARCHAR(50) NOT NULL, -- website, email, phone, referral, event, social, direct, digital, traditional, partner
    subsource_name VARCHAR(100), -- Google Ads, Facebook, LinkedIn, etc.
    channel VARCHAR(100), -- Direct, Organic Search, Paid Search, Social Media, Referral, etc.
    is_active BOOLEAN DEFAULT true,
    description TEXT,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(tenant_id, source_code),
    CONSTRAINT valid_source_type CHECK (source_type IN ('website', 'email', 'phone', 'referral', 'event', 'social', 'direct', 'digital', 'traditional', 'partner'))
);

CREATE INDEX idx_sales_lead_sources_tenant ON sales_lead_sources(tenant_id);
CREATE INDEX idx_sales_lead_sources_type ON sales_lead_sources(source_type);

-- ============================================================================
-- 3. LEAD MILESTONE TRACKING TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS sales_lead_milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    lead_id UUID NOT NULL REFERENCES sales_leads(id) ON DELETE CASCADE,
    milestone_type VARCHAR(50) NOT NULL, -- lead_generated, contacted, site_visit, revisit, demo, proposal, negotiation, booking, cancellation, reengaged
    milestone_date DATE NOT NULL,
    milestone_time TIME,
    notes TEXT,
    location_latitude DECIMAL(10,8),
    location_longitude DECIMAL(11,8),
    location_name VARCHAR(255),
    status_before VARCHAR(50),
    status_after VARCHAR(50),
    visited_by UUID REFERENCES users(id),
    duration_minutes INT, -- for site visits/demos
    outcome VARCHAR(100), -- for site visits: positive, neutral, negative
    follow_up_date DATE,
    follow_up_required BOOLEAN DEFAULT false,
    metadata JSONB, -- flexible field for additional data
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_milestone_type CHECK (milestone_type IN ('lead_generated', 'contacted', 'site_visit', 'revisit', 'demo', 'proposal', 'negotiation', 'booking', 'cancellation', 'reengaged'))
);

CREATE INDEX idx_sales_lead_milestones_tenant ON sales_lead_milestones(tenant_id);
CREATE INDEX idx_sales_lead_milestones_lead ON sales_lead_milestones(lead_id);
CREATE INDEX idx_sales_lead_milestones_type ON sales_lead_milestones(milestone_type);
CREATE INDEX idx_sales_lead_milestones_date ON sales_lead_milestones(milestone_date);

-- ============================================================================
-- 4. LEAD ENGAGEMENT TRACKING TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS sales_lead_engagement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    lead_id UUID NOT NULL REFERENCES sales_leads(id) ON DELETE CASCADE,
    engagement_type VARCHAR(50) NOT NULL, -- email_sent, call_made, message_sent, meeting_scheduled, proposal_sent, quote_sent
    engagement_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    engagement_channel VARCHAR(50), -- email, phone, sms, whatsapp, in_person, video
    subject VARCHAR(255),
    notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'completed', -- completed, pending, failed
    response_received BOOLEAN DEFAULT false,
    response_date TIMESTAMP WITH TIME ZONE,
    response_notes TEXT,
    assigned_to UUID REFERENCES users(id),
    duration_seconds INT, -- for calls/meetings
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_engagement_type CHECK (engagement_type IN ('email_sent', 'call_made', 'message_sent', 'meeting_scheduled', 'proposal_sent', 'quote_sent')),
    CONSTRAINT valid_engagement_status CHECK (status IN ('completed', 'pending', 'failed'))
);

CREATE INDEX idx_sales_lead_engagement_tenant ON sales_lead_engagement(tenant_id);
CREATE INDEX idx_sales_lead_engagement_lead ON sales_lead_engagement(lead_id);
CREATE INDEX idx_sales_lead_engagement_type ON sales_lead_engagement(engagement_type);
CREATE INDEX idx_sales_lead_engagement_date ON sales_lead_engagement(engagement_date);

-- ============================================================================
-- 5. BOOKING AND UNITS MANAGEMENT TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS sales_bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES sales_customers(id) ON DELETE CASCADE,
    lead_id UUID REFERENCES sales_leads(id),
    booking_code VARCHAR(50) NOT NULL,
    booking_date DATE NOT NULL,
    booking_amount DECIMAL(15,2) NOT NULL,
    unit_type VARCHAR(100), -- property type, product type, etc.
    unit_count INT NOT NULL DEFAULT 1,
    units_booked INT NOT NULL DEFAULT 1,
    units_available INT,
    delivery_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed', -- confirmed, pending, cancelled, completed, on_hold
    cancellation_date DATE,
    cancellation_reason TEXT,
    cancellation_refund_amount DECIMAL(15,2),
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(tenant_id, booking_code),
    CONSTRAINT valid_booking_status CHECK (status IN ('confirmed', 'pending', 'cancelled', 'completed', 'on_hold'))
);

CREATE INDEX idx_sales_bookings_tenant ON sales_bookings(tenant_id);
CREATE INDEX idx_sales_bookings_customer ON sales_bookings(customer_id);
CREATE INDEX idx_sales_bookings_lead ON sales_bookings(lead_id);
CREATE INDEX idx_sales_bookings_date ON sales_bookings(booking_date);
CREATE INDEX idx_sales_bookings_status ON sales_bookings(status);

-- ============================================================================
-- 6. ACCOUNT LEDGER TABLE (for customer financial tracking)
-- ============================================================================
CREATE TABLE IF NOT EXISTS sales_account_ledgers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES sales_customers(id) ON DELETE CASCADE,
    ledger_code VARCHAR(50) NOT NULL,
    ledger_date DATE NOT NULL,
    transaction_type VARCHAR(50) NOT NULL, -- invoice, payment, credit_note, debit_note, adjustment
    reference_document_type VARCHAR(50), -- invoice, booking, quote, order
    reference_document_id UUID,
    debit_amount DECIMAL(15,2) DEFAULT 0,
    credit_amount DECIMAL(15,2) DEFAULT 0,
    balance_after DECIMAL(15,2),
    description VARCHAR(255),
    remarks TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, reversed, cancelled
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, ledger_code),
    CONSTRAINT valid_transaction_type CHECK (transaction_type IN ('invoice', 'payment', 'credit_note', 'debit_note', 'adjustment')),
    CONSTRAINT valid_ledger_status CHECK (status IN ('active', 'reversed', 'cancelled'))
);

CREATE INDEX idx_sales_account_ledgers_tenant ON sales_account_ledgers(tenant_id);
CREATE INDEX idx_sales_account_ledgers_customer ON sales_account_ledgers(customer_id);
CREATE INDEX idx_sales_account_ledgers_date ON sales_account_ledgers(ledger_date);
CREATE INDEX idx_sales_account_ledgers_type ON sales_account_ledgers(transaction_type);

-- ============================================================================
-- 7. REPORTING AND ANALYTICS VIEWS
-- ============================================================================

-- Lead Funnel Analysis
CREATE OR REPLACE VIEW v_lead_funnel_analysis AS
SELECT 
    tenant_id,
    DATE_TRUNC('month', created_at)::DATE as month,
    COUNT(CASE WHEN status = 'new' THEN 1 END) as leads_new,
    COUNT(CASE WHEN status = 'contacted' THEN 1 END) as leads_contacted,
    COUNT(CASE WHEN status = 'qualified' THEN 1 END) as leads_qualified,
    COUNT(CASE WHEN status = 'negotiation' THEN 1 END) as leads_negotiation,
    COUNT(CASE WHEN status = 'converted' THEN 1 END) as leads_converted,
    COUNT(CASE WHEN status = 'lost' THEN 1 END) as leads_lost,
    COUNT(*) as total_leads,
    ROUND(CAST(COUNT(CASE WHEN status = 'converted' THEN 1 END) AS NUMERIC) / NULLIF(COUNT(*), 0) * 100, 2) as conversion_rate
FROM sales_leads
WHERE deleted_at IS NULL
GROUP BY tenant_id, DATE_TRUNC('month', created_at)::DATE
ORDER BY month DESC;

-- Lead Source Performance
CREATE OR REPLACE VIEW v_lead_source_performance AS
SELECT 
    sl.tenant_id,
    sl.source as source_type,
    sls.source_name,
    sls.subsource_name,
    COUNT(sl.id) as total_leads,
    COUNT(CASE WHEN sl.converted_to_customer = true THEN 1 END) as converted_leads,
    ROUND(CAST(COUNT(CASE WHEN sl.converted_to_customer = true THEN 1 END) AS NUMERIC) / NULLIF(COUNT(sl.id), 0) * 100, 2) as conversion_rate,
    COUNT(CASE WHEN sl.status = 'lost' THEN 1 END) as lost_leads,
    SUM(CASE WHEN sl.probability > 0 THEN 1 ELSE 0 END) as qualified_leads
FROM sales_leads sl
LEFT JOIN sales_lead_sources sls ON sl.source = sls.source_code AND sl.tenant_id = sls.tenant_id
WHERE sl.deleted_at IS NULL
GROUP BY sl.tenant_id, sl.source, sls.source_name, sls.subsource_name
ORDER BY total_leads DESC;

-- Campaign Performance
CREATE OR REPLACE VIEW v_campaign_performance AS
SELECT 
    sc.tenant_id,
    sc.id as campaign_id,
    sc.campaign_code,
    sc.campaign_name,
    sc.campaign_type,
    COUNT(DISTINCT sl.id) as total_leads,
    COUNT(DISTINCT CASE WHEN sl.converted_to_customer = true THEN sl.id END) as converted_customers,
    ROUND(CAST(COUNT(DISTINCT CASE WHEN sl.converted_to_customer = true THEN sl.id END) AS NUMERIC) / 
          NULLIF(COUNT(DISTINCT sl.id), 0) * 100, 2) as conversion_rate,
    COUNT(DISTINCT CASE WHEN sl.status = 'lost' THEN sl.id END) as lost_leads,
    sc.budget,
    sc.expected_roi
FROM sales_campaigns sc
LEFT JOIN sales_leads sl ON sc.id = sl.campaign_id AND sc.tenant_id = sl.tenant_id
WHERE sc.deleted_at IS NULL AND sl.deleted_at IS NULL
GROUP BY sc.tenant_id, sc.id, sc.campaign_code, sc.campaign_name, sc.campaign_type, sc.budget, sc.expected_roi
ORDER BY total_leads DESC;

-- Booking Status Summary
CREATE OR REPLACE VIEW v_booking_summary AS
SELECT 
    tenant_id,
    status,
    COUNT(id) as booking_count,
    SUM(booking_amount) as total_booking_amount,
    SUM(units_booked) as total_units_booked,
    AVG(booking_amount) as avg_booking_amount,
    DATE_TRUNC('month', booking_date)::DATE as booking_month
FROM sales_bookings
WHERE deleted_at IS NULL
GROUP BY tenant_id, status, DATE_TRUNC('month', booking_date)::DATE
ORDER BY booking_month DESC, booking_count DESC;

-- Customer Ledger Summary
CREATE OR REPLACE VIEW v_customer_ledger_summary AS
SELECT 
    sal.tenant_id,
    sal.customer_id,
    sc.customer_name,
    SUM(CASE WHEN sal.transaction_type = 'invoice' THEN sal.debit_amount ELSE 0 END) as total_invoiced,
    SUM(CASE WHEN sal.transaction_type = 'payment' THEN sal.credit_amount ELSE 0 END) as total_paid,
    SUM(CASE WHEN sal.transaction_type = 'credit_note' THEN sal.credit_amount ELSE 0 END) as total_credit_notes,
    (SUM(CASE WHEN sal.transaction_type = 'invoice' THEN sal.debit_amount ELSE 0 END) - 
     SUM(CASE WHEN sal.transaction_type = 'payment' THEN sal.credit_amount ELSE 0 END) -
     SUM(CASE WHEN sal.transaction_type = 'credit_note' THEN sal.credit_amount ELSE 0 END)) as outstanding_balance,
    COUNT(CASE WHEN sal.transaction_type = 'invoice' THEN 1 END) as total_transactions
FROM sales_account_ledgers sal
JOIN sales_customers sc ON sal.customer_id = sc.id
WHERE sal.deleted_at IS NULL AND sal.status = 'active'
GROUP BY sal.tenant_id, sal.customer_id, sc.customer_name
ORDER BY outstanding_balance DESC;

-- Lead Milestone Timeline
CREATE OR REPLACE VIEW v_lead_milestone_timeline AS
SELECT 
    slm.tenant_id,
    slm.lead_id,
    sl.first_name,
    sl.last_name,
    sl.email,
    slm.milestone_type,
    slm.milestone_date,
    slm.milestone_time,
    slm.notes,
    slm.status_before,
    slm.status_after,
    LEAD(slm.milestone_date) OVER (PARTITION BY slm.lead_id ORDER BY slm.milestone_date) as next_milestone_date,
    (LEAD(slm.milestone_date) OVER (PARTITION BY slm.lead_id ORDER BY slm.milestone_date) - slm.milestone_date)::INT as days_to_next_milestone
FROM sales_lead_milestones slm
JOIN sales_leads sl ON slm.lead_id = sl.id
WHERE slm.deleted_at IS NULL
ORDER BY slm.tenant_id, slm.lead_id, slm.milestone_date;

-- ============================================================================
-- 8. GRANT PERMISSIONS
-- ============================================================================
-- Assuming role exists from previous migrations
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;

COMMIT;
