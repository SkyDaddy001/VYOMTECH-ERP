-- Migration: Partner Management System
-- Version: 024
-- Description: Create partners and partner_users tables for multi-tenant partner ecosystem

-- ============================================================================
-- PARTNERS TABLE
-- ============================================================================
-- Stores partner organizations (portals, channel partners, vendors, customers)

CREATE TABLE IF NOT EXISTS partners (
    id CHAR(26) PRIMARY KEY,
    tenant_id VARCHAR(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    partner_code VARCHAR(50) NOT NULL UNIQUE,
    organization_name VARCHAR(150) NOT NULL,
    partner_type ENUM('portal', 'channel_partner', 'vendor', 'customer') NOT NULL DEFAULT 'customer',
    status ENUM('pending', 'active', 'inactive', 'suspended', 'rejected') NOT NULL DEFAULT 'pending',
    
    -- Contact Information
    contact_email VARCHAR(100) NOT NULL UNIQUE,
    contact_phone VARCHAR(20),
    contact_person VARCHAR(100),
    website VARCHAR(200),
    description TEXT,
    
    -- Address Information
    address TEXT,
    city VARCHAR(50),
    state VARCHAR(50),
    country VARCHAR(50) DEFAULT 'India',
    zip_code VARCHAR(10),
    
    -- Tax and Banking
    tax_id VARCHAR(50),
    banking_details JSON,
    
    -- Commission and Lead Pricing
    commission_percentage DECIMAL(5,2) DEFAULT 0.00,
    lead_price DECIMAL(10,2) DEFAULT 0.00,
    monthly_quota INT DEFAULT 100,
    current_month_leads INT DEFAULT 0,
    
    -- Lead and Conversion Tracking
    total_leads_submitted BIGINT DEFAULT 0,
    approved_leads BIGINT DEFAULT 0,
    rejected_leads BIGINT DEFAULT 0,
    converted_leads BIGINT DEFAULT 0,
    
    -- Financial Tracking
    total_earnings DECIMAL(15,2) DEFAULT 0.00,
    pending_payout_amount DECIMAL(15,2) DEFAULT 0.00,
    withdrawn_amount DECIMAL(15,2) DEFAULT 0.00,
    available_balance DECIMAL(15,2) DEFAULT 0.00,
    
    -- Approval and Status Management
    approved_by BIGINT,
    approved_at TIMESTAMP NULL,
    rejection_reason TEXT,
    suspension_reason TEXT,
    suspended_at TIMESTAMP NULL,
    
    -- Documents
    document_urls JSON,
    
    -- Audit Fields
    created_by CHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    KEY idx_tenant_id (tenant_id),
    KEY idx_partner_code (partner_code),
    KEY idx_status (status),
    KEY idx_partner_type (partner_type),
    KEY idx_contact_email (contact_email),
    KEY idx_created_at (created_at),
    
    -- Foreign Keys
    CONSTRAINT fk_partners_tenant FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================================================
-- PARTNER_USERS TABLE
-- ============================================================================
-- Stores user accounts for partner organizations

CREATE TABLE IF NOT EXISTS partner_users (
    id CHAR(26) PRIMARY KEY,
    partner_id VARCHAR(36) NOT NULL,
    tenant_id VARCHAR(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    
    -- Authentication
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    
    -- User Information
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    
    -- Access Control
    role ENUM('admin', 'manager', 'user', 'viewer') NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Security
    last_login TIMESTAMP NULL,
    password_changed_at TIMESTAMP NULL,
    failed_login_attempts INT DEFAULT 0,
    locked_until TIMESTAMP NULL,
    
    -- Audit Fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    KEY idx_partner_id (partner_id),
    KEY idx_tenant_id (tenant_id),
    KEY idx_email (email),
    KEY idx_is_active (is_active),
    KEY idx_created_at (created_at),
    
    -- Foreign Keys
    CONSTRAINT fk_partner_users_partner FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    CONSTRAINT fk_partner_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

-- Partner search and filtering
CREATE INDEX idx_partners_tenant_status ON partners(tenant_id, status);
CREATE INDEX idx_partners_tenant_type ON partners(tenant_id, partner_type);
CREATE INDEX idx_partners_tenant_created ON partners(tenant_id, created_at DESC);

-- Partner user authentication
CREATE INDEX idx_partner_users_partner_active ON partner_users(partner_id, is_active);
CREATE INDEX idx_partner_users_email_tenant ON partner_users(email, tenant_id);

-- ============================================================================
-- AUDIT LOGGING (Optional - for tracking partner changes)
-- ============================================================================

CREATE TABLE IF NOT EXISTS partner_audit_log (
    id CHAR(26) PRIMARY KEY,
    partner_id VARCHAR(36) NOT NULL,
    tenant_id VARCHAR(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    action VARCHAR(50) NOT NULL, -- 'created', 'updated', 'status_changed', 'suspended', 'deleted'
    changed_fields JSON,
    changed_by BIGINT,
    change_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    KEY idx_partner_id (partner_id),
    KEY idx_tenant_id (tenant_id),
    KEY idx_action (action),
    KEY idx_created_at (created_at),
    
    CONSTRAINT fk_audit_partner FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    CONSTRAINT fk_audit_tenant FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================================================
-- PARTNER PAYOUT TRACKING (Optional - for commission management)
-- ============================================================================

CREATE TABLE IF NOT EXISTS partner_payouts (
    id CHAR(26) PRIMARY KEY,
    partner_id VARCHAR(36) NOT NULL,
    tenant_id VARCHAR(36) COLLATE utf8mb4_unicode_ci NOT NULL,
    
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    
    leads_approved INT DEFAULT 0,
    conversion_count INT DEFAULT 0,
    total_commission DECIMAL(15,2) DEFAULT 0.00,
    status ENUM('pending', 'approved', 'processed', 'paid', 'cancelled') DEFAULT 'pending',
    
    payment_method VARCHAR(50), -- 'bank_transfer', 'check', 'paypal', 'wire'
    payout_date DATE,
    payout_reference VARCHAR(100),
    
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    KEY idx_partner_id (partner_id),
    KEY idx_tenant_id (tenant_id),
    KEY idx_status (status),
    KEY idx_period (period_start, period_end),
    
    CONSTRAINT fk_payouts_partner FOREIGN KEY (partner_id) REFERENCES partners(id) ON DELETE CASCADE,
    CONSTRAINT fk_payouts_tenant FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



