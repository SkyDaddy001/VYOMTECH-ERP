-- Real Estate Property Management Schema
-- Includes: Area Statement, Cost Sheet, Control Sheet, Customer Details, Payment Details

-- ============================================
-- AREA STATEMENT TABLES
-- ============================================

CREATE TABLE IF NOT EXISTS property_projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    project_name VARCHAR(255) NOT NULL,
    project_code VARCHAR(50) UNIQUE NOT NULL,
    location VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    total_units INT,
    total_area DECIMAL(12, 2),
    project_type VARCHAR(50), -- 'residential', 'commercial', 'mixed'
    status VARCHAR(50), -- 'planning', 'under_construction', 'ready', 'sold_out'
    launch_date DATE,
    expected_completion DATE,
    actual_completion DATE,
    noc_status VARCHAR(50), -- 'pending', 'approved', 'applied'
    noc_date DATE,
    developer_name VARCHAR(255),
    architect_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

CREATE TABLE IF NOT EXISTS property_blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    project_id UUID NOT NULL,
    block_name VARCHAR(100) NOT NULL,
    block_code VARCHAR(50),
    wing_name VARCHAR(100),
    total_units INT,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (project_id) REFERENCES property_projects(id)
);

CREATE TABLE IF NOT EXISTS property_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    project_id UUID NOT NULL,
    block_id UUID NOT NULL,
    unit_number VARCHAR(100) NOT NULL,
    floor INT,
    unit_type VARCHAR(50), -- '1BHK', '2BHK', '3BHK', 'shop', 'office', etc.
    facing VARCHAR(50), -- 'north', 'south', 'east', 'west', 'corner'
    
    -- Area Details
    carpet_area DECIMAL(10, 2),
    carpet_area_with_balcony DECIMAL(10, 2),
    utility_area DECIMAL(10, 2),
    plinth_area DECIMAL(10, 2),
    sbua DECIMAL(10, 2), -- Super Built Up Area
    uds_sqft DECIMAL(10, 2), -- Undivided Share per sqft
    
    -- Status
    status VARCHAR(50), -- 'available', 'booked', 'sold', 'reserved'
    
    -- Allocation
    alloted_to VARCHAR(255),
    allotment_date DATE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (project_id) REFERENCES property_projects(id),
    FOREIGN KEY (block_id) REFERENCES property_blocks(id),
    UNIQUE(tenant_id, project_id, unit_number)
);

-- ============================================
-- COST SHEET TABLES
-- ============================================

CREATE TABLE IF NOT EXISTS unit_cost_sheets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    unit_id UUID NOT NULL,
    
    -- Basic Pricing
    rate_per_sqft DECIMAL(12, 4),
    sbua_rate DECIMAL(12, 2),
    base_price DECIMAL(14, 2),
    
    -- Additional Charges
    frc DECIMAL(14, 2), -- Floor Rise Charge
    car_parking_cost DECIMAL(14, 2),
    plc DECIMAL(14, 2), -- Preferential Location Charges
    statutory_charges DECIMAL(14, 2),
    other_charges DECIMAL(14, 2),
    legal_charges DECIMAL(14, 2),
    
    -- Final Prices
    apartment_cost_exc_govt DECIMAL(14, 2),
    apartment_cost_inc_govt DECIMAL(14, 2),
    composite_guideline_value DECIMAL(14, 2),
    actual_sold_price DECIMAL(14, 2),
    
    -- Parking
    car_parking_type VARCHAR(50), -- 'covered', 'open', 'none', 'tandem', 'hybrid'
    parking_location VARCHAR(255),
    
    effective_date DATE,
    validity_date DATE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (unit_id) REFERENCES property_units(id)
);

-- ============================================
-- CUSTOMER BOOKING DETAILS
-- ============================================

CREATE TABLE IF NOT EXISTS customer_bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    unit_id UUID NOT NULL,
    lead_id UUID,
    customer_id UUID,
    
    -- Booking Information
    booking_date DATE NOT NULL,
    booking_reference VARCHAR(100) UNIQUE,
    booking_status VARCHAR(50), -- 'active', 'cancelled', 'completed'
    welcome_date DATE,
    
    -- Progress Stages
    allotment_date DATE,
    agreement_date DATE,
    registration_date DATE,
    handover_date DATE,
    possession_date DATE,
    
    -- Financial
    rate_per_sqft DECIMAL(12, 4),
    composite_guideline_value DECIMAL(14, 2),
    
    -- Parking Details
    car_parking_type VARCHAR(50),
    parking_location VARCHAR(255),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (unit_id) REFERENCES property_units(id)
);

CREATE TABLE IF NOT EXISTS customer_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    booking_id UUID NOT NULL,
    
    -- Primary Applicant
    primary_name VARCHAR(255),
    primary_phone VARCHAR(20),
    primary_alternate_phone VARCHAR(20),
    primary_email VARCHAR(255),
    primary_communication_address TEXT,
    primary_permanent_address TEXT,
    primary_aadhar_no VARCHAR(50),
    primary_pan_no VARCHAR(50),
    
    -- Co-Applicant 1
    coapplicant1_name VARCHAR(255),
    coapplicant1_phone VARCHAR(20),
    coapplicant1_alternate_phone VARCHAR(20),
    coapplicant1_email VARCHAR(255),
    coapplicant1_communication_address TEXT,
    coapplicant1_permanent_address TEXT,
    coapplicant1_aadhar_no VARCHAR(50),
    coapplicant1_pan_no VARCHAR(50),
    coapplicant1_relation VARCHAR(100),
    care_of_co1 VARCHAR(255),
    
    -- Co-Applicant 2
    coapplicant2_name VARCHAR(255),
    coapplicant2_phone VARCHAR(20),
    coapplicant2_alternate_phone VARCHAR(20),
    coapplicant2_email VARCHAR(255),
    coapplicant2_communication_address TEXT,
    coapplicant2_permanent_address TEXT,
    coapplicant2_aadhar_no VARCHAR(50),
    coapplicant2_pan_no VARCHAR(50),
    coapplicant2_relation VARCHAR(100),
    care_of_co2 VARCHAR(255),
    
    -- Power of Attorney Holder
    poa_holder_name VARCHAR(255),
    poa_document_no VARCHAR(100),
    poa_relation VARCHAR(100),
    life_certificate_no VARCHAR(100),
    
    -- Bank Loan Details
    bank_name VARCHAR(255),
    loan_contact_person VARCHAR(255),
    loan_contact_number VARCHAR(20),
    loan_sanction_date DATE,
    connector_code VARCHAR(100),
    sanction_amount DECIMAL(14, 2),
    
    -- Sales Information
    sales_executive_id UUID,
    sales_executive_name VARCHAR(255),
    sales_head_name VARCHAR(255),
    booking_source VARCHAR(100), -- 'direct', 'broker', 'partner', 'site_visit'
    sales_head_id UUID,
    
    -- Other Details
    other_works TEXT,
    maintenance_charges DECIMAL(10, 2),
    corpus_charges DECIMAL(10, 2),
    eb_deposit DECIMAL(10, 2),
    noc_received_date DATE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (booking_id) REFERENCES customer_bookings(id)
);

-- ============================================
-- PAYMENT DETAILS TABLES
-- ============================================

CREATE TABLE IF NOT EXISTS booking_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    booking_id UUID NOT NULL,
    
    payment_date DATE NOT NULL,
    payment_mode VARCHAR(50), -- 'cash', 'cheque', 'transfer', 'neft', 'rtgs', 'demand_draft'
    paid_by VARCHAR(255),
    receipt_number VARCHAR(100) UNIQUE NOT NULL,
    receipt_date DATE,
    
    towards VARCHAR(100), -- APARTMENT COST, GST, REGISTRATION COST, MODT, OTHER WORKS, MAINTAINANCE, CORPUS, EB DEPOSIT, TDS, 
    amount DECIMAL(14, 2),
    
    -- Payment Details
    cheque_number VARCHAR(50),
    cheque_date DATE,
    bank_name VARCHAR(255),
    transaction_id VARCHAR(255),
    
    status VARCHAR(50), -- 'pending', 'cleared', 'bounced', 'cancelled'
    remarks TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (booking_id) REFERENCES customer_bookings(id)
);

CREATE TABLE IF NOT EXISTS payment_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    booking_id UUID NOT NULL,
    
    schedule_name VARCHAR(100),
    payment_stage VARCHAR(100), -- 'booking', 'agreement', 'possession', 'handover'
    payment_percentage DECIMAL(5, 2),
    payment_amount DECIMAL(14, 2),
    due_date DATE,
    
    amount_paid DECIMAL(14, 2) DEFAULT 0,
    outstanding DECIMAL(14, 2),
    status VARCHAR(50), -- 'pending', 'partial', 'completed', 'overdue'
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (booking_id) REFERENCES customer_bookings(id)
);

-- ============================================
-- ACCOUNT LEDGER FOR CUSTOMERS
-- ============================================

CREATE TABLE IF NOT EXISTS customer_account_ledgers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    booking_id UUID NOT NULL,
    customer_id UUID,
    
    -- Ledger Entry
    transaction_date DATE NOT NULL,
    transaction_type VARCHAR(50), -- 'credit', 'debit', 'adjustment'
    description VARCHAR(255),
    
    -- Amount
    debit_amount DECIMAL(14, 2) DEFAULT 0,
    credit_amount DECIMAL(14, 2) DEFAULT 0,
    opening_balance DECIMAL(14, 2),
    closing_balance DECIMAL(14, 2),
    
    -- Reference
    payment_id UUID,
    invoice_id UUID,
    reference_number VARCHAR(100),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (booking_id) REFERENCES customer_bookings(id)
);

-- ============================================
-- CONTROL SHEET / CONFIGURATION
-- ============================================

CREATE TABLE IF NOT EXISTS project_control_sheet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    project_id UUID NOT NULL,
    
    attribute_name VARCHAR(100) NOT NULL,
    attribute_value VARCHAR(255),
    attribute_type VARCHAR(50), -- 'text', 'number', 'date', 'percentage', 'boolean'
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (project_id) REFERENCES property_projects(id)
);

-- ============================================
-- MILESTONE & CAMPAIGN TRACKING
-- ============================================

CREATE TABLE IF NOT EXISTS property_milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    booking_id UUID NOT NULL,
    
    -- Campaign Information
    campaign_id UUID,
    campaign_name VARCHAR(255),
    source VARCHAR(100), -- 'direct', 'site_visit', 'broker', 'referral', 'digital', 'exhibition'
    subsource VARCHAR(100),
    
    -- Milestone Dates
    lead_generated_date DATE,
    re_engaged_date DATE,
    site_visit_date DATE,
    revisit_date DATE,
    booking_date DATE,
    cancelled_date DATE,
    
    -- Status
    status VARCHAR(50),
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (booking_id) REFERENCES customer_bookings(id)
);

-- ============================================
-- INDEXES FOR PERFORMANCE
-- ============================================

CREATE INDEX idx_property_projects_tenant ON property_projects(tenant_id);
CREATE INDEX idx_property_blocks_project ON property_blocks(project_id);
CREATE INDEX idx_property_units_project ON property_units(project_id);
CREATE INDEX idx_property_units_status ON property_units(status);
CREATE INDEX idx_unit_cost_sheets_unit ON unit_cost_sheets(unit_id);
CREATE INDEX idx_customer_bookings_tenant ON customer_bookings(tenant_id);
CREATE INDEX idx_customer_bookings_unit ON customer_bookings(unit_id);
CREATE INDEX idx_customer_details_booking ON customer_details(booking_id);
CREATE INDEX idx_booking_payments_booking ON booking_payments(booking_id);
CREATE INDEX idx_payment_schedules_booking ON payment_schedules(booking_id);
CREATE INDEX idx_account_ledgers_booking ON customer_account_ledgers(booking_id);
CREATE INDEX idx_property_milestones_booking ON property_milestones(booking_id);
CREATE INDEX idx_property_milestones_source ON property_milestones(source);
