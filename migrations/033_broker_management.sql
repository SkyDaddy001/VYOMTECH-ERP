-- Migration 033: Broker Management System
-- Created: 2025-12-23
-- Purpose: Track brokers, commission structures, and payouts for real estate sales

-- ============================================================================
-- TABLE: broker_profile
-- Description: Master list of brokers/agents
-- ============================================================================
CREATE TABLE IF NOT EXISTS broker_profile (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    
    -- Broker Details
    broker_name VARCHAR(100) NOT NULL,
    broker_email VARCHAR(100),
    broker_phone VARCHAR(20),
    broker_license_no VARCHAR(50) UNIQUE,
    rera_registration_no VARCHAR(50) UNIQUE,
    pan_no VARCHAR(20),
    gst_no VARCHAR(20),
    
    -- Business Details
    broker_firm_name VARCHAR(100),
    firm_registration_no VARCHAR(50),
    business_address TEXT,
    business_city VARCHAR(50),
    business_state VARCHAR(50),
    business_postal_code VARCHAR(10),
    
    -- Financial Details
    bank_account_no VARCHAR(20),
    bank_name VARCHAR(100),
    ifsc_code VARCHAR(20),
    beneficiary_name VARCHAR(100),
    
    -- Status & Compliance
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, inactive, suspended
    kyc_verified BOOLEAN DEFAULT FALSE,
    kyc_verified_date DATETIME,
    compliance_status VARCHAR(50) DEFAULT 'pending',
    
    -- Performance Tracking
    total_bookings INT DEFAULT 0,
    total_commission_earned DECIMAL(15, 2) DEFAULT 0,
    total_commission_paid DECIMAL(15, 2) DEFAULT 0,
    
    -- Metadata
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    KEY idx_tenant_id (tenant_id),
    KEY idx_status (status),
    KEY idx_broker_email (broker_email),
    CONSTRAINT fk_broker_tenant FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLE: broker_commission_structure
-- Description: Commission rates and slabs for brokers
-- ============================================================================
CREATE TABLE IF NOT EXISTS broker_commission_structure (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    broker_id BIGINT NOT NULL,
    
    -- Commission Type
    commission_type VARCHAR(50) NOT NULL, -- percentage, fixed_amount, slab_based
    
    -- For Percentage-based
    commission_percentage DECIMAL(5, 2),
    
    -- For Fixed Amount
    fixed_amount DECIMAL(12, 2),
    
    -- Slab-based Commission
    slab_1_max_amount DECIMAL(15, 2),
    slab_1_commission_percentage DECIMAL(5, 2),
    slab_2_min_amount DECIMAL(15, 2),
    slab_2_max_amount DECIMAL(15, 2),
    slab_2_commission_percentage DECIMAL(5, 2),
    slab_3_min_amount DECIMAL(15, 2),
    slab_3_commission_percentage DECIMAL(5, 2),
    
    -- Applicability
    applicable_for VARCHAR(50), -- residential, commercial, all
    min_unit_price DECIMAL(12, 2),
    max_unit_price DECIMAL(12, 2),
    
    -- Validity
    effective_from DATE NOT NULL,
    effective_till DATE,
    
    -- Status
    status VARCHAR(50) DEFAULT 'active', -- active, inactive
    
    -- Metadata
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    KEY idx_tenant_id (tenant_id),
    KEY idx_broker_id (broker_id),
    KEY idx_effective_dates (effective_from, effective_till),
    CONSTRAINT fk_commission_broker FOREIGN KEY (broker_id) REFERENCES broker_profile(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLE: broker_booking_link
-- Description: Links bookings to brokers for commission tracking
-- ============================================================================
CREATE TABLE IF NOT EXISTS broker_booking_link (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    broker_id BIGINT NOT NULL,
    booking_id BIGINT NOT NULL,
    
    -- Booking Details
    unit_price DECIMAL(15, 2),
    booking_amount DECIMAL(15, 2),
    
    -- Commission Calculation
    commission_structure_id BIGINT,
    commission_percentage DECIMAL(5, 2),
    commission_amount DECIMAL(12, 2),
    
    -- Status
    booking_status VARCHAR(50) DEFAULT 'active', -- active, cancelled, completed
    commission_status VARCHAR(50) DEFAULT 'pending', -- pending, approved, paid, cancelled
    
    -- Dates
    booking_date DATETIME,
    approval_date DATETIME,
    payment_date DATETIME,
    
    -- Metadata
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    KEY idx_tenant_id (tenant_id),
    KEY idx_broker_id (broker_id),
    KEY idx_booking_id (booking_id),
    KEY idx_commission_status (commission_status),
    CONSTRAINT fk_broker_booking_broker FOREIGN KEY (broker_id) REFERENCES broker_profile(id),
    CONSTRAINT fk_broker_booking_commission FOREIGN KEY (commission_structure_id) REFERENCES broker_commission_structure(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- TABLE: broker_commission_payout
-- Description: Commission payout tracking and reconciliation
-- ============================================================================
CREATE TABLE IF NOT EXISTS broker_commission_payout (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    broker_id BIGINT NOT NULL,
    
    -- Payout Period
    payout_period_from DATE NOT NULL,
    payout_period_till DATE NOT NULL,
    payout_reference_no VARCHAR(100) UNIQUE,
    
    -- Commission Calculation
    total_commission DECIMAL(15, 2) NOT NULL,
    tds_amount DECIMAL(12, 2),
    gst_amount DECIMAL(12, 2),
    net_payable_amount DECIMAL(15, 2) NOT NULL,
    
    -- Payment Details
    payment_mode VARCHAR(50), -- bank_transfer, cheque, cash, neft, rtgs
    payment_date DATETIME,
    payment_reference_no VARCHAR(100),
    bank_name VARCHAR(100),
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, approved, paid, rejected, cancelled
    approval_by BIGINT,
    approval_date DATETIME,
    
    -- Documents
    calculation_sheet_url VARCHAR(500),
    payment_proof_url VARCHAR(500),
    
    -- Metadata
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by BIGINT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    KEY idx_tenant_id (tenant_id),
    KEY idx_broker_id (broker_id),
    KEY idx_status (status),
    KEY idx_payout_period (payout_period_from, payout_period_till),
    CONSTRAINT fk_payout_broker FOREIGN KEY (broker_id) REFERENCES broker_profile(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- INDEXES & CONSTRAINTS
-- ============================================================================

-- Performance indexes
CREATE INDEX idx_broker_active ON broker_profile(status, created_at DESC);
CREATE INDEX idx_commission_due ON broker_commission_payout(status, payment_date);
CREATE INDEX idx_broker_performance ON broker_profile(total_bookings, total_commission_earned);

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================
