-- ============================================================================
-- Migration: 034_joint_applicants.sql
-- Phase 1.3: Joint Applicants Management for Real Estate
-- ============================================================================
-- This migration creates tables for managing joint applicants in real estate
-- bookings, including co-ownership details, income verification, and legal
-- documentation for multiple property owners.

CREATE TABLE IF NOT EXISTS joint_applicant (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  booking_id VARCHAR(36) NOT NULL,
  
  -- Primary Applicant Flag
  is_primary_applicant BOOLEAN DEFAULT FALSE,
  
  -- Applicant Details
  applicant_name VARCHAR(255) NOT NULL,
  date_of_birth DATE,
  gender ENUM('male', 'female', 'other'),
  marital_status VARCHAR(50),
  nationality VARCHAR(100),
  
  -- Contact Information
  email VARCHAR(255),
  phone_number VARCHAR(20),
  alternate_phone VARCHAR(20),
  address_line1 VARCHAR(255),
  address_line2 VARCHAR(255),
  city VARCHAR(100),
  state_province VARCHAR(100),
  postal_code VARCHAR(20),
  country VARCHAR(100),
  
  -- Identification
  id_type ENUM('aadhar', 'pan', 'passport', 'driving_license', 'voter_id'),
  id_number VARCHAR(50),
  id_issue_date DATE,
  id_expiry_date DATE,
  
  -- Occupation & Income
  occupation VARCHAR(100),
  employer_name VARCHAR(255),
  annual_income DECIMAL(12,2),
  income_source VARCHAR(100),
  bank_name VARCHAR(255),
  account_number VARCHAR(50),
  account_type ENUM('savings', 'current', 'nro', 'nre'),
  
  -- Co-Ownership Details
  ownership_share_percentage DECIMAL(5,2),
  ownership_type ENUM('joint_tenant', 'tenant_in_common', 'coparcenary'),
  
  -- Verification Status
  kyc_status VARCHAR(50), -- pending, verified, rejected
  kyc_verified_date TIMESTAMP NULL,
  income_verified BOOLEAN DEFAULT FALSE,
  income_verified_date TIMESTAMP NULL,
  document_verification_status VARCHAR(50), -- pending, approved, rejected
  
  -- Legal Consent
  consent_given BOOLEAN DEFAULT FALSE,
  consent_given_date TIMESTAMP NULL,
  legal_notice_acknowledged BOOLEAN DEFAULT FALSE,
  
  -- Metadata
  created_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  -- Indexes
  INDEX idx_tenant_id (tenant_id),
  INDEX idx_booking_id (booking_id),
  INDEX idx_kyc_status (kyc_status),
  INDEX idx_created_at (created_at),
  UNIQUE KEY unique_booking_applicant (booking_id, email, id_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- Joint Applicant Document Table
-- ============================================================================
CREATE TABLE IF NOT EXISTS joint_applicant_document (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  joint_applicant_id VARCHAR(36) NOT NULL,
  
  -- Document Details
  document_type VARCHAR(100), -- aadhar_copy, pan_copy, income_proof, bank_statement, etc.
  document_name VARCHAR(255) NOT NULL,
  document_url VARCHAR(500) NOT NULL,
  document_hash VARCHAR(255),
  
  -- File Information
  file_size BIGINT,
  file_type VARCHAR(50), -- pdf, jpg, png, etc.
  upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  -- Verification
  document_status VARCHAR(50), -- pending, approved, rejected, expired
  verification_date TIMESTAMP NULL,
  verification_notes TEXT,
  
  -- Metadata
  created_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  -- Indexes
  INDEX idx_tenant_id (tenant_id),
  INDEX idx_joint_applicant_id (joint_applicant_id),
  INDEX idx_document_type (document_type),
  INDEX idx_document_status (document_status),
  FOREIGN KEY (joint_applicant_id) REFERENCES joint_applicant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- Co-Ownership Agreement Table
-- ============================================================================
CREATE TABLE IF NOT EXISTS co_ownership_agreement (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  booking_id VARCHAR(36) NOT NULL,
  
  -- Agreement Details
  agreement_reference_no VARCHAR(100) UNIQUE NOT NULL,
  agreement_date DATE,
  agreement_type ENUM('joint_tenant', 'tenant_in_common', 'coparcenary'),
  
  -- Property Details in Agreement
  property_description TEXT,
  agreed_purchase_price DECIMAL(12,2),
  
  -- Ownership Distribution
  total_share_percentage DECIMAL(5,2),
  
  -- Status & Verification
  agreement_status VARCHAR(50), -- draft, pending_approval, approved, executed, cancelled
  notarized BOOLEAN DEFAULT FALSE,
  notarization_date TIMESTAMP NULL,
  registered_with_authorities BOOLEAN DEFAULT FALSE,
  registration_number VARCHAR(100),
  registration_date DATE,
  
  -- Approvals
  all_parties_agreed BOOLEAN DEFAULT FALSE,
  legal_review_completed BOOLEAN DEFAULT FALSE,
  legal_reviewer_notes TEXT,
  
  -- Document Storage
  agreement_document_url VARCHAR(500),
  notarized_copy_url VARCHAR(500),
  registration_certificate_url VARCHAR(500),
  
  -- Metadata
  created_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  -- Indexes
  INDEX idx_tenant_id (tenant_id),
  INDEX idx_booking_id (booking_id),
  INDEX idx_agreement_status (agreement_status),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- Joint Applicant Signatory (Legal Signatories Table)
-- ============================================================================
CREATE TABLE IF NOT EXISTS co_ownership_signatory (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  co_ownership_agreement_id VARCHAR(36) NOT NULL,
  joint_applicant_id VARCHAR(36) NOT NULL,
  
  -- Signatory Details
  signatory_name VARCHAR(255) NOT NULL,
  signatory_email VARCHAR(255),
  signatory_phone VARCHAR(20),
  signature_type ENUM('digital', 'physical', 'aadhaar_signed'),
  
  -- Signature Status
  signature_status VARCHAR(50), -- pending, signed, expired, cancelled
  signature_date TIMESTAMP NULL,
  signature_ip_address VARCHAR(50),
  signature_device_info VARCHAR(255),
  
  -- Signature Document
  signature_image_url VARCHAR(500),
  signed_document_url VARCHAR(500),
  
  -- Witness Information (if required)
  witness_name VARCHAR(255),
  witness_id_type VARCHAR(50),
  witness_id_number VARCHAR(50),
  witness_signature_url VARCHAR(500),
  
  -- Metadata
  created_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  -- Indexes
  INDEX idx_tenant_id (tenant_id),
  INDEX idx_co_ownership_agreement_id (co_ownership_agreement_id),
  INDEX idx_joint_applicant_id (joint_applicant_id),
  INDEX idx_signature_status (signature_status),
  FOREIGN KEY (co_ownership_agreement_id) REFERENCES co_ownership_agreement(id) ON DELETE CASCADE,
  FOREIGN KEY (joint_applicant_id) REFERENCES joint_applicant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- Joint Applicant Income Verification Table
-- ============================================================================
CREATE TABLE IF NOT EXISTS joint_applicant_income_verification (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  joint_applicant_id VARCHAR(36) NOT NULL,
  
  -- Income Details
  annual_income DECIMAL(12,2),
  monthly_income DECIMAL(12,2),
  income_year INT,
  income_tax_slab VARCHAR(50),
  
  -- Source of Income
  primary_occupation VARCHAR(100),
  primary_employer VARCHAR(255),
  employment_type ENUM('salaried', 'self_employed', 'business', 'professional', 'retired'),
  
  -- Financial Status
  total_assets DECIMAL(12,2),
  total_liabilities DECIMAL(12,2),
  net_worth DECIMAL(12,2),
  
  -- Bank Details
  primary_bank_name VARCHAR(255),
  primary_account_number VARCHAR(50),
  primary_account_type VARCHAR(50),
  bank_verification_status VARCHAR(50),
  
  -- Income Tax Information
  pan_number VARCHAR(50),
  aadhar_number VARCHAR(50),
  last_itr_filed_year INT,
  last_itr_amount DECIMAL(12,2),
  
  -- Verification Documents
  salary_slip_url VARCHAR(500),
  itr_document_url VARCHAR(500),
  bank_statement_url VARCHAR(500),
  employer_letter_url VARCHAR(500),
  
  -- Verification Status
  verification_status VARCHAR(50), -- pending, verified, rejected, expired
  verification_date TIMESTAMP NULL,
  verifier_id CHAR(26),
  verification_notes TEXT,
  
  -- Metadata
  created_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  -- Indexes
  INDEX idx_tenant_id (tenant_id),
  INDEX idx_joint_applicant_id (joint_applicant_id),
  INDEX idx_verification_status (verification_status),
  INDEX idx_pan_number (pan_number),
  FOREIGN KEY (joint_applicant_id) REFERENCES joint_applicant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- Joint Applicant Liability Tracking Table
-- ============================================================================
CREATE TABLE IF NOT EXISTS joint_applicant_liability (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  joint_applicant_id VARCHAR(36) NOT NULL,
  
  -- Liability Details
  liability_type ENUM('loan', 'credit_card', 'mortgage', 'personal_loan', 'business_loan'),
  creditor_name VARCHAR(255),
  
  -- Loan/Credit Details
  outstanding_amount DECIMAL(12,2),
  monthly_payment DECIMAL(12,2),
  interest_rate DECIMAL(5,2),
  loan_maturity_date DATE,
  
  -- Liability Status
  status VARCHAR(50), -- active, settled, defaulted
  
  -- Documents
  document_url VARCHAR(500),
  
  -- Metadata
  created_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_by CHAR(36),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  -- Indexes
  INDEX idx_tenant_id (tenant_id),
  INDEX idx_joint_applicant_id (joint_applicant_id),
  INDEX idx_liability_type (liability_type),
  FOREIGN KEY (joint_applicant_id) REFERENCES joint_applicant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- Audit Log for Joint Applicants
-- ============================================================================
CREATE TABLE IF NOT EXISTS joint_applicant_audit_log (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  joint_applicant_id VARCHAR(36) NOT NULL,
  
  -- Audit Details
  action VARCHAR(50), -- created, updated, verified, rejected, etc.
  action_description TEXT,
  field_changed VARCHAR(255),
  old_value TEXT,
  new_value TEXT,
  
  -- User Information
  user_id CHAR(26),
  user_email VARCHAR(255),
  ip_address VARCHAR(50),
  
  -- Metadata
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  -- Indexes
  INDEX idx_tenant_id (tenant_id),
  INDEX idx_joint_applicant_id (joint_applicant_id),
  INDEX idx_action (action),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;





