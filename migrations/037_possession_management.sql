-- ====================================
-- Phase 2.1: Possession Management
-- ====================================

-- Possession status tracking table
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS possession_statuses (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  booking_id VARCHAR(36) NOT NULL,
  status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, in_progress, completed, cancelled
  possession_date DATETIME,
  estimated_possession_date DATETIME,
  possession_reason VARCHAR(500),
  possession_type VARCHAR(50) NOT NULL, -- normal, partial, interim, final
  is_complete BOOLEAN DEFAULT FALSE,
  completion_percentage DECIMAL(5, 2) DEFAULT 0,
  notes LONGTEXT,
  created_by CHAR(36),
  updated_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  
  UNIQUE KEY unique_booking_possession (tenant_id, booking_id),
  KEY idx_tenant_status (tenant_id, status),
  KEY idx_booking_status (booking_id, status),
  KEY idx_possession_type (tenant_id, possession_type),
  KEY idx_created_at (created_at),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (booking_id) REFERENCES booking(id)
);

-- Possession documents tracking table
CREATE TABLE IF NOT EXISTS possession_documents (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  possession_id VARCHAR(36) NOT NULL,
  document_type VARCHAR(100) NOT NULL, -- possession_letter, handover_checklist, keys, utilities_list, final_statement, insurance_doc
  document_name VARCHAR(255) NOT NULL,
  document_url VARCHAR(500),
  file_name VARCHAR(255),
  file_size BIGINT,
  file_format VARCHAR(10),
  s3_bucket VARCHAR(255),
  s3_key VARCHAR(500),
  document_status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, submitted, verified, rejected, approved
  verification_notes LONGTEXT,
  verified_by BIGINT,
  verified_at DATETIME,
  is_mandatory BOOLEAN DEFAULT FALSE,
  uploaded_by BIGINT,
  uploaded_at DATETIME,
  metadata JSON,
  created_by CHAR(36),
  updated_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  
  KEY idx_tenant_possession (tenant_id, possession_id),
  KEY idx_document_type (tenant_id, document_type),
  KEY idx_document_status (document_status),
  KEY idx_verified_by (verified_by),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (possession_id) REFERENCES possession_statuses(id),
  FOREIGN KEY (verified_by) REFERENCES `user`(id)
);

-- Registration process tracking table
CREATE TABLE IF NOT EXISTS possession_registrations (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  possession_id VARCHAR(36) NOT NULL,
  registration_type VARCHAR(50) NOT NULL, -- registration, name_transfer, title_transfer, mortgage_release
  registration_number VARCHAR(100),
  registration_office VARCHAR(255),
  registration_date DATETIME,
  registration_status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, in_progress, completed, rejected, appealed
  amount_paid DECIMAL(15, 2) DEFAULT 0,
  amount_pending DECIMAL(15, 2) DEFAULT 0,
  payment_mode VARCHAR(50), -- online, cheque, demand_draft, cash, bank_transfer
  reference_number VARCHAR(100),
  submission_date DATETIME,
  expected_completion_date DATETIME,
  actual_completion_date DATETIME,
  remarks LONGTEXT,
  approved_by BIGINT,
  approved_at DATETIME,
  created_by CHAR(36),
  updated_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  
  KEY idx_tenant_registration (tenant_id, registration_type),
  KEY idx_possession_registration (possession_id),
  KEY idx_registration_status (registration_status),
  KEY idx_registration_number (registration_number),
  KEY idx_approved_by (approved_by),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (possession_id) REFERENCES possession_statuses(id),
  FOREIGN KEY (approved_by) REFERENCES `user`(id)
);

-- Possession certificates table
CREATE TABLE IF NOT EXISTS possession_certificates (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  possession_id VARCHAR(36) NOT NULL,
  certificate_type VARCHAR(50) NOT NULL, -- possession_certificate, occupancy_certificate, completion_certificate, no_dues_certificate
  certificate_number VARCHAR(100),
  issuing_authority VARCHAR(255),
  issue_date DATETIME,
  validity_date DATETIME,
  certificate_url VARCHAR(500),
  file_name VARCHAR(255),
  file_size BIGINT,
  s3_bucket VARCHAR(255),
  s3_key VARCHAR(500),
  certificate_status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, issued, verified, expired, cancelled
  verification_status VARCHAR(50) DEFAULT 'pending', -- pending, verified, rejected
  verified_by BIGINT,
  verified_at DATETIME,
  verification_notes LONGTEXT,
  metadata JSON,
  created_by CHAR(36),
  updated_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  
  UNIQUE KEY unique_cert_number (tenant_id, certificate_number),
  KEY idx_tenant_certificate (tenant_id, certificate_type),
  KEY idx_possession_certificate (possession_id),
  KEY idx_certificate_status (certificate_status),
  KEY idx_verified_by (verified_by),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (possession_id) REFERENCES possession_statuses(id),
  FOREIGN KEY (verified_by) REFERENCES `user`(id)
);

-- Possession approvals and sign-offs table
CREATE TABLE IF NOT EXISTS possession_approvals (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  possession_id VARCHAR(36) NOT NULL,
  approval_type VARCHAR(50) NOT NULL, -- possession_approval, document_approval, registration_approval, final_approval
  approver_id CHAR(26) NOT NULL,
  approver_role VARCHAR(100),
  approval_status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, approved, rejected, conditional
  approval_notes LONGTEXT,
  conditional_remarks LONGTEXT,
  approval_date DATETIME,
  valid_from DATETIME,
  valid_till DATETIME,
  sequence_order INT,
  is_final_approval BOOLEAN DEFAULT FALSE,
  metadata JSON,
  created_by CHAR(36),
  updated_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  
  KEY idx_tenant_approval (tenant_id, approval_type),
  KEY idx_possession_approval (possession_id),
  KEY idx_approver_id (approver_id),
  KEY idx_approval_status (approval_status),
  KEY idx_sequence_order (possession_id, sequence_order),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (possession_id) REFERENCES possession_statuses(id),
  FOREIGN KEY (approver_id) REFERENCES `user`(id)
);

-- Possession audit log table
CREATE TABLE IF NOT EXISTS possession_audit_log (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  possession_id VARCHAR(36) NOT NULL,
  action VARCHAR(100) NOT NULL,
  entity_type VARCHAR(50) NOT NULL, -- possession, document, registration, certificate, approval
  entity_id VARCHAR(36),
  old_value LONGTEXT,
  new_value LONGTEXT,
  change_reason VARCHAR(500),
  performed_by BIGINT,
  user_ip_address VARCHAR(45),
  user_agent VARCHAR(500),
  metadata JSON,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  KEY idx_tenant_audit (tenant_id, created_at),
  KEY idx_possession_audit (possession_id),
  KEY idx_action (action),
  KEY idx_entity (entity_type, entity_id),
  KEY idx_performed_by (performed_by),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (possession_id) REFERENCES possession_statuses(id),
  FOREIGN KEY (performed_by) REFERENCES `user`(id)
);

-- Create indexes for improved query performance
CREATE INDEX idx_pos_dates ON possession_statuses(estimated_possession_date, possession_date);
CREATE INDEX idx_pos_type_status ON possession_statuses(tenant_id, possession_type, status);
CREATE INDEX idx_reg_dates ON possession_registrations(submission_date, expected_completion_date);
CREATE INDEX idx_cert_dates ON possession_certificates(issue_date, validity_date);
CREATE INDEX idx_approval_dates ON possession_approvals(approval_date, valid_till);

SET FOREIGN_KEY_CHECKS = 1;





