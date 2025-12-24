-- ====================================
-- Phase 2.2: Title Clearance Management
-- ====================================

-- Title clearance status table
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS title_clearances (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  booking_id VARCHAR(36) NOT NULL,
  property_id VARCHAR(36),
  status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, in_progress, cleared, issues_found, rejected, expired
  clearance_type VARCHAR(50) NOT NULL, -- full_clearance, encumbrance_check, mutation, legal_opinion, boundary_verification
  start_date DATETIME,
  target_clearance_date DATETIME,
  actual_clearance_date DATETIME,
  clearance_percentage DECIMAL(5, 2) DEFAULT 0,
  issues_count INT DEFAULT 0,
  resolved_issues_count INT DEFAULT 0,
  priority VARCHAR(50) DEFAULT 'normal', -- low, normal, high, critical
  notes LONGTEXT,
  created_by CHAR(36),
  updated_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  
  UNIQUE KEY unique_booking_clearance (tenant_id, booking_id),
  KEY idx_tenant_status (tenant_id, status),
  KEY idx_booking_status (booking_id, status),
  KEY idx_clearance_type (tenant_id, clearance_type),
  KEY idx_priority (tenant_id, priority),
  KEY idx_created_at (created_at),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (booking_id) REFERENCES booking(id)
);

-- Title issues/encumbrances table
CREATE TABLE IF NOT EXISTS title_issues (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  clearance_id CHAR(26) NOT NULL,
  issue_type VARCHAR(100) NOT NULL, -- lien, encumbrance, dispute, boundary_issue, mortgage, legal_claim, tax_issue
  issue_title VARCHAR(255) NOT NULL,
  issue_description LONGTEXT,
  severity VARCHAR(50) NOT NULL DEFAULT 'medium', -- low, medium, high, critical
  status VARCHAR(50) NOT NULL DEFAULT 'open', -- open, under_review, escalated, resolved, deferred, invalid
  reported_date DATETIME,
  source_document VARCHAR(255),
  affected_parties LONGTEXT, -- JSON list of parties involved
  resolution_notes LONGTEXT,
  resolved_date DATETIME,
  resolved_by BIGINT,
  resolution_method VARCHAR(100), -- legal_opinion, court_order, mutual_agreement, rectification, insurance
  metadata JSON,
  created_by CHAR(36),
  updated_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  
  KEY idx_tenant_issue (tenant_id, issue_type),
  KEY idx_clearance_issue (clearance_id),
  KEY idx_issue_status (status),
  KEY idx_severity (severity),
  KEY idx_resolved_by (resolved_by),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (clearance_id) REFERENCES title_clearances(id),
  FOREIGN KEY (resolved_by) REFERENCES `user`(id)
);

-- Title search reports table
CREATE TABLE IF NOT EXISTS title_search_reports (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  clearance_id CHAR(26) NOT NULL,
  search_type VARCHAR(50) NOT NULL, -- government_records, registry_search, municipal_search, court_search, property_search
  search_date DATETIME,
  search_authority VARCHAR(255),
  search_reference_number VARCHAR(100),
  report_url VARCHAR(500),
  report_file_name VARCHAR(255),
  report_file_size BIGINT,
  s3_bucket VARCHAR(255),
  s3_key VARCHAR(500),
  encumbrances_found INT DEFAULT 0,
  search_status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, completed, issues_found, verified
  verified_by BIGINT,
  verified_at DATETIME,
  verification_notes LONGTEXT,
  search_cost DECIMAL(12, 2) DEFAULT 0,
  metadata JSON,
  created_by CHAR(36),
  updated_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  
  KEY idx_tenant_search (tenant_id, search_type),
  KEY idx_clearance_search (clearance_id),
  KEY idx_search_status (search_status),
  KEY idx_verified_by (verified_by),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (clearance_id) REFERENCES title_clearances(id),
  FOREIGN KEY (verified_by) REFERENCES `user`(id)
);

-- Title verification checklist table
CREATE TABLE IF NOT EXISTS title_verification_checklists (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  clearance_id CHAR(26) NOT NULL,
  item_name VARCHAR(255) NOT NULL,
  item_category VARCHAR(100), -- ownership, encumbrance, mutation, boundary, legal, tax, litigation
  description LONGTEXT,
  is_mandatory BOOLEAN DEFAULT FALSE,
  status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, verified, not_applicable, issue_found
  verified_by BIGINT,
  verified_at DATETIME,
  verification_notes LONGTEXT,
  sequence_order INT,
  metadata JSON,
  created_by CHAR(36),
  updated_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  
  KEY idx_tenant_checklist (tenant_id, item_category),
  KEY idx_clearance_checklist (clearance_id),
  KEY idx_checklist_status (status),
  KEY idx_verified_by (verified_by),
  KEY idx_sequence (clearance_id, sequence_order),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (clearance_id) REFERENCES title_clearances(id),
  FOREIGN KEY (verified_by) REFERENCES `user`(id)
);

-- Legal opinions and expert review table
CREATE TABLE IF NOT EXISTS title_legal_opinions (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  clearance_id CHAR(26) NOT NULL,
  opinion_type VARCHAR(50) NOT NULL, -- legal_opinion, expert_review, boundary_survey, environmental_review
  expert_name VARCHAR(255),
  expert_organization VARCHAR(255),
  expert_license_number VARCHAR(100),
  opinion_date DATETIME,
  opinion_status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, received, under_review, approved, concerns_noted, rejected
  opinion_url VARCHAR(500),
  opinion_file_name VARCHAR(255),
  file_size BIGINT,
  s3_bucket VARCHAR(255),
  s3_key VARCHAR(500),
  opinion_summary LONGTEXT,
  recommendations LONGTEXT,
  risk_assessment VARCHAR(50), -- low_risk, medium_risk, high_risk
  review_by_lawyer BIGINT,
  review_notes LONGTEXT,
  reviewed_at DATETIME,
  cost DECIMAL(12, 2) DEFAULT 0,
  metadata JSON,
  created_by CHAR(36),
  updated_by CHAR(36),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  
  KEY idx_tenant_opinion (tenant_id, opinion_type),
  KEY idx_clearance_opinion (clearance_id),
  KEY idx_opinion_status (opinion_status),
  KEY idx_risk_assessment (risk_assessment),
  KEY idx_reviewed_by (review_by_lawyer),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (clearance_id) REFERENCES title_clearances(id),
  FOREIGN KEY (review_by_lawyer) REFERENCES `user`(id)
);

-- Title clearance approvals table
CREATE TABLE IF NOT EXISTS title_clearance_approvals (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  clearance_id CHAR(26) NOT NULL,
  approval_type VARCHAR(50) NOT NULL, -- issue_resolution_approval, legal_opinion_approval, final_clearance_approval
  approver_id CHAR(26) NOT NULL,
  approver_role VARCHAR(100),
  approval_status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, approved, rejected, conditional
  approval_notes LONGTEXT,
  conditional_requirements LONGTEXT,
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
  KEY idx_clearance_approval (clearance_id),
  KEY idx_approver_id (approver_id),
  KEY idx_approval_status (approval_status),
  KEY idx_sequence_order (clearance_id, sequence_order),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (clearance_id) REFERENCES title_clearances(id),
  FOREIGN KEY (approver_id) REFERENCES `user`(id)
);

-- Title clearance audit log table
CREATE TABLE IF NOT EXISTS title_clearance_audit_log (
  id CHAR(26) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  clearance_id CHAR(26) NOT NULL,
  action VARCHAR(100) NOT NULL,
  entity_type VARCHAR(50) NOT NULL, -- clearance, issue, search_report, legal_opinion, approval
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
  KEY idx_clearance_audit (clearance_id),
  KEY idx_action (action),
  KEY idx_entity (entity_type, entity_id),
  KEY idx_performed_by (performed_by),
  FOREIGN KEY (tenant_id) REFERENCES `tenant`(id),
  FOREIGN KEY (clearance_id) REFERENCES title_clearances(id),
  FOREIGN KEY (performed_by) REFERENCES `user`(id)
);

-- Create indexes for improved query performance
CREATE INDEX idx_title_dates ON title_clearances(target_clearance_date, actual_clearance_date);
CREATE INDEX idx_title_type_status ON title_clearances(tenant_id, clearance_type, status);
CREATE INDEX idx_issue_dates ON title_issues(reported_date, resolved_date);
CREATE INDEX idx_search_dates ON title_search_reports(search_date, verified_at);
CREATE INDEX idx_opinion_dates ON title_legal_opinions(opinion_date, reviewed_at);
CREATE INDEX idx_approval_dates ON title_clearance_approvals(approval_date, valid_till);

SET FOREIGN_KEY_CHECKS = 1;





