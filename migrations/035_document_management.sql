-- Document Management System (Phase 1.4 - Real Estate)
-- Handles all document types: property docs, identity proofs, financial documents, legal agreements

-- Document Categories and Types
CREATE TABLE IF NOT EXISTS document_categories (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    category_name VARCHAR(100) NOT NULL,
    category_code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    icon_url VARCHAR(500),
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    INDEX idx_tenant_category (tenant_id, category_code),
    INDEX idx_active_category (is_active)
);

-- Document Types (e.g., Aadhaar, PAN, Property Deed, etc.)
CREATE TABLE IF NOT EXISTS document_types (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    type_name VARCHAR(100) NOT NULL,
    type_code VARCHAR(50) NOT NULL,
    description TEXT,
    is_mandatory BOOLEAN DEFAULT false,
    is_identity_proof BOOLEAN DEFAULT false,
    is_property_doc BOOLEAN DEFAULT false,
    is_financial_doc BOOLEAN DEFAULT false,
    expiry_required BOOLEAN DEFAULT false,
    file_formats VARCHAR(500),
    max_file_size INT DEFAULT 10485760,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (category_id) REFERENCES document_categories(id),
    INDEX idx_tenant_type (tenant_id, type_code),
    UNIQUE KEY uk_tenant_type (tenant_id, type_code)
);

-- Main Documents Table
CREATE TABLE IF NOT EXISTS documents (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    document_type_id BIGINT NOT NULL,
    entity_type VARCHAR(50) NOT NULL COMMENT 'booking, property, applicant, etc.',
    entity_id BIGINT NOT NULL COMMENT 'ID of the entity (booking_id, property_id, etc.)',
    document_name VARCHAR(255) NOT NULL,
    document_description TEXT,
    document_url VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT,
    file_format VARCHAR(20),
    s3_bucket VARCHAR(255),
    s3_key VARCHAR(500),
    thumbnail_url VARCHAR(500),
    document_status ENUM('pending', 'verified', 'rejected', 'expired') DEFAULT 'pending',
    verification_status ENUM('pending', 'verified', 'rejected') DEFAULT 'pending',
    verification_notes TEXT,
    verifier_id BIGINT,
    verified_at TIMESTAMP NULL,
    issue_date DATE,
    expiry_date DATE,
    is_primary BOOLEAN DEFAULT false,
    metadata JSON COMMENT 'Additional metadata (OCR data, signatures, etc.)',
    tags VARCHAR(500),
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (document_type_id) REFERENCES document_types(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    FOREIGN KEY (verifier_id) REFERENCES users(id),
    INDEX idx_tenant_entity (tenant_id, entity_type, entity_id),
    INDEX idx_document_status (document_status),
    INDEX idx_verification_status (verification_status),
    INDEX idx_expiry_date (expiry_date),
    INDEX idx_document_type (document_type_id)
);

-- Document Templates (Standard templates for properties/projects)
CREATE TABLE IF NOT EXISTS document_templates (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    template_name VARCHAR(255) NOT NULL,
    template_description TEXT,
    template_type VARCHAR(50) NOT NULL COMMENT 'agreement, receipt, notice, etc.',
    template_content LONGTEXT,
    required_documents JSON COMMENT 'List of document type IDs required for this template',
    is_active BOOLEAN DEFAULT true,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    INDEX idx_tenant_template (tenant_id, is_active)
);

-- Document Verification Workflow
CREATE TABLE IF NOT EXISTS document_verifications (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    document_id BIGINT NOT NULL,
    verification_type VARCHAR(50) NOT NULL COMMENT 'manual, ocr, ai_based, third_party',
    verifier_id BIGINT,
    verification_notes TEXT,
    verification_result ENUM('pending', 'approved', 'rejected', 'requires_clarification') DEFAULT 'pending',
    rejection_reason TEXT,
    ai_confidence_score DECIMAL(5,2),
    ai_verification_data JSON COMMENT 'OCR results, face matching scores, etc.',
    third_party_result JSON,
    verification_duration_minutes INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (document_id) REFERENCES documents(id),
    FOREIGN KEY (verifier_id) REFERENCES users(id),
    INDEX idx_tenant_document (tenant_id, document_id),
    INDEX idx_verification_result (verification_result)
);

-- Document Compliance Status (Track regulatory compliance)
CREATE TABLE IF NOT EXISTS document_compliance (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    document_id BIGINT NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id BIGINT NOT NULL,
    compliance_requirement VARCHAR(100) NOT NULL COMMENT 'rera, itax, gst, etc.',
    is_compliant BOOLEAN DEFAULT false,
    last_checked_at TIMESTAMP NULL,
    compliance_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (document_id) REFERENCES documents(id),
    INDEX idx_tenant_compliance (tenant_id, entity_type, entity_id),
    INDEX idx_compliance_requirement (compliance_requirement)
);

-- Document Audit Trail
CREATE TABLE IF NOT EXISTS document_audit_log (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    document_id BIGINT NOT NULL,
    action VARCHAR(50) NOT NULL COMMENT 'upload, verify, reject, download, share, etc.',
    action_by BIGINT,
    action_details TEXT,
    ip_address VARCHAR(50),
    user_agent VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (document_id) REFERENCES documents(id),
    FOREIGN KEY (action_by) REFERENCES users(id),
    INDEX idx_tenant_document_action (tenant_id, document_id, action),
    INDEX idx_action_timestamp (created_at)
);

-- Document Sharing & Permissions
CREATE TABLE IF NOT EXISTS document_shares (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    document_id BIGINT NOT NULL,
    shared_with_user_id BIGINT,
    shared_with_role VARCHAR(50),
    share_permission VARCHAR(50) NOT NULL COMMENT 'view, download, verify, edit',
    expiry_date TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT true,
    shared_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (document_id) REFERENCES documents(id),
    FOREIGN KEY (shared_with_user_id) REFERENCES users(id),
    FOREIGN KEY (shared_by) REFERENCES users(id),
    INDEX idx_tenant_document_share (tenant_id, document_id),
    INDEX idx_shared_with_user (shared_with_user_id)
);

-- Document Collections (Bundles of documents for a purpose)
CREATE TABLE IF NOT EXISTS document_collections (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    collection_name VARCHAR(255) NOT NULL,
    collection_description TEXT,
    collection_type VARCHAR(50) NOT NULL COMMENT 'booking, property, project, etc.',
    entity_type VARCHAR(50) NOT NULL,
    entity_id BIGINT NOT NULL,
    status ENUM('incomplete', 'complete', 'verified') DEFAULT 'incomplete',
    completion_percentage INT DEFAULT 0,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    INDEX idx_tenant_collection (tenant_id, entity_type, entity_id),
    INDEX idx_collection_status (status)
);

-- Document Collection Items (Documents in a collection)
CREATE TABLE IF NOT EXISTS document_collection_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    collection_id BIGINT NOT NULL,
    document_id BIGINT NOT NULL,
    is_mandatory BOOLEAN DEFAULT false,
    status ENUM('pending', 'submitted', 'verified', 'rejected') DEFAULT 'pending',
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (collection_id) REFERENCES document_collections(id) ON DELETE CASCADE,
    FOREIGN KEY (document_id) REFERENCES documents(id),
    UNIQUE KEY uk_collection_document (collection_id, document_id)
);
