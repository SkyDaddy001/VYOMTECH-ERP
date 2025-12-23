-- Customer Portal System
-- Phase 2.3: Complete customer-facing portal for dashboard, notifications, messaging, and tracking

-- ============================================
-- CUSTOMER DASHBOARD & PROFILE
-- ============================================

CREATE TABLE IF NOT EXISTS customer_profiles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    booking_id BIGINT,
    phone_number VARCHAR(20),
    alternate_phone VARCHAR(20),
    email_address VARCHAR(255),
    alternate_email VARCHAR(255),
    date_of_birth DATE,
    gender ENUM('male', 'female', 'other'),
    
    -- Address Information
    current_address VARCHAR(500),
    permanent_address VARCHAR(500),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(10),
    country VARCHAR(100),
    
    -- Document Information
    id_proof_type VARCHAR(50),
    id_proof_number VARCHAR(100),
    id_proof_file_url VARCHAR(500),
    s3_bucket VARCHAR(100),
    s3_key VARCHAR(200),
    
    -- Preferences
    communication_preference ENUM('email', 'sms', 'phone', 'in_app'),
    language_preference VARCHAR(50),
    timezone VARCHAR(100),
    notification_enabled BOOLEAN DEFAULT TRUE,
    email_updates_enabled BOOLEAN DEFAULT TRUE,
    sms_updates_enabled BOOLEAN DEFAULT FALSE,
    push_notifications_enabled BOOLEAN DEFAULT TRUE,
    
    -- Profile Status
    profile_completion_percentage FLOAT DEFAULT 0,
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP NULL,
    verification_notes TEXT,
    
    -- Metadata
    metadata JSON,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    UNIQUE KEY uk_customer_user (tenant_id, user_id),
    INDEX idx_tenant_customer (tenant_id),
    INDEX idx_booking_id (booking_id),
    INDEX idx_profile_completion (profile_completion_percentage),
    CONSTRAINT fk_customer_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- CUSTOMER NOTIFICATIONS
-- ============================================

CREATE TABLE IF NOT EXISTS customer_notifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    notification_type VARCHAR(100),
    title VARCHAR(200),
    message TEXT NOT NULL,
    description TEXT,
    
    -- Notification Category
    category ENUM('booking', 'payment', 'document', 'possession', 'legal', 'document_upload', 'status_update', 'alert', 'announcement', 'reminder'),
    priority ENUM('low', 'normal', 'high', 'critical') DEFAULT 'normal',
    
    -- Related Entities
    related_entity_type VARCHAR(100),
    related_entity_id BIGINT,
    
    -- Status
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP NULL,
    is_archived BOOLEAN DEFAULT FALSE,
    archived_at TIMESTAMP NULL,
    
    -- Delivery Status
    delivery_status ENUM('pending', 'sent', 'delivered', 'read', 'failed') DEFAULT 'pending',
    
    -- Delivery Channels
    email_sent BOOLEAN DEFAULT FALSE,
    email_sent_at TIMESTAMP NULL,
    sms_sent BOOLEAN DEFAULT FALSE,
    sms_sent_at TIMESTAMP NULL,
    push_sent BOOLEAN DEFAULT FALSE,
    push_sent_at TIMESTAMP NULL,
    in_app_sent BOOLEAN DEFAULT TRUE,
    in_app_sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Action URL
    action_url VARCHAR(500),
    cta_text VARCHAR(100),
    
    -- Metadata & Expiry
    metadata JSON,
    expires_at TIMESTAMP NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_user_notification (tenant_id, user_id),
    INDEX idx_user_read_status (user_id, is_read),
    INDEX idx_notification_created (created_at DESC),
    INDEX idx_notification_category (category),
    INDEX idx_notification_priority (priority),
    INDEX idx_related_entity (related_entity_type, related_entity_id),
    CONSTRAINT fk_notif_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- CUSTOMER MESSAGES & COMMUNICATION
-- ============================================

CREATE TABLE IF NOT EXISTS customer_conversations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    customer_user_id BIGINT NOT NULL,
    support_user_id BIGINT,
    booking_id BIGINT,
    
    -- Conversation Details
    subject VARCHAR(255),
    conversation_type ENUM('support', 'inquiry', 'complaint', 'feedback', 'follow_up') DEFAULT 'support',
    status ENUM('open', 'in_progress', 'waiting_customer', 'waiting_support', 'resolved', 'closed') DEFAULT 'open',
    priority ENUM('low', 'normal', 'high', 'urgent') DEFAULT 'normal',
    
    -- Last Activity
    last_message_at TIMESTAMP NULL,
    last_message_from VARCHAR(50),
    
    -- Assignment
    assigned_to BIGINT,
    assigned_at TIMESTAMP NULL,
    
    -- Resolution
    resolution_notes TEXT,
    resolved_at TIMESTAMP NULL,
    resolved_by BIGINT,
    
    -- Metadata
    metadata JSON,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_conversation (tenant_id),
    INDEX idx_customer_conversation (customer_user_id),
    INDEX idx_conversation_status (status),
    INDEX idx_conversation_created (created_at DESC),
    CONSTRAINT fk_conv_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS customer_messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    conversation_id BIGINT NOT NULL,
    sender_user_id BIGINT NOT NULL,
    sender_type ENUM('customer', 'support', 'system'),
    
    message_text TEXT NOT NULL,
    message_type ENUM('text', 'image', 'document', 'system'),
    
    -- File Attachment
    attachment_url VARCHAR(500),
    attachment_file_name VARCHAR(200),
    attachment_file_size BIGINT,
    attachment_file_type VARCHAR(50),
    s3_bucket VARCHAR(100),
    s3_key VARCHAR(200),
    
    -- Message Status
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP NULL,
    
    -- Metadata
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_conversation_message (conversation_id),
    INDEX idx_sender_message (sender_user_id),
    INDEX idx_message_created (created_at DESC),
    CONSTRAINT fk_message_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_message_conversation FOREIGN KEY (conversation_id) REFERENCES customer_conversations(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- CUSTOMER DOCUMENT MANAGEMENT
-- ============================================

CREATE TABLE IF NOT EXISTS customer_document_uploads (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    booking_id BIGINT,
    
    -- Document Details
    document_type VARCHAR(100),
    document_name VARCHAR(255),
    document_description TEXT,
    
    -- File Information
    file_url VARCHAR(500),
    file_name VARCHAR(255),
    file_size BIGINT,
    file_extension VARCHAR(20),
    file_mime_type VARCHAR(100),
    
    -- Storage
    s3_bucket VARCHAR(100),
    s3_key VARCHAR(200),
    
    -- Upload Status
    upload_status ENUM('uploading', 'completed', 'failed', 'virus_detected') DEFAULT 'uploading',
    upload_progress_percentage FLOAT DEFAULT 0,
    
    -- Verification
    verification_status ENUM('pending', 'verified', 'rejected') DEFAULT 'pending',
    verification_notes TEXT,
    verified_by BIGINT,
    verified_at TIMESTAMP NULL,
    
    -- Required Status
    is_required_document BOOLEAN DEFAULT FALSE,
    is_verified_by_admin BOOLEAN DEFAULT FALSE,
    
    -- Metadata
    metadata JSON,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_uploads (tenant_id),
    INDEX idx_user_uploads (user_id),
    INDEX idx_booking_uploads (booking_id),
    INDEX idx_upload_status (upload_status),
    INDEX idx_verification_status (verification_status),
    CONSTRAINT fk_upload_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- CUSTOMER BOOKING & PROPERTY TRACKING
-- ============================================

CREATE TABLE IF NOT EXISTS customer_booking_tracking (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    booking_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    
    -- Booking Status Tracking
    current_status VARCHAR(100),
    previous_status VARCHAR(100),
    status_changed_at TIMESTAMP NULL,
    status_change_reason TEXT,
    
    -- Property Details
    property_id BIGINT,
    property_name VARCHAR(255),
    property_location VARCHAR(500),
    unit_number VARCHAR(100),
    
    -- Timeline
    booking_date DATE,
    possession_date DATE,
    estimated_handover_date DATE,
    
    -- Payment Tracking
    total_amount DECIMAL(15, 2),
    amount_paid DECIMAL(15, 2) DEFAULT 0,
    amount_pending DECIMAL(15, 2),
    payment_percentage FLOAT DEFAULT 0,
    
    -- Milestone Tracking
    total_milestones INT,
    completed_milestones INT DEFAULT 0,
    pending_milestones INT DEFAULT 0,
    
    -- Document Status
    required_documents_count INT,
    uploaded_documents_count INT DEFAULT 0,
    verified_documents_count INT DEFAULT 0,
    
    -- Last Update
    last_update_at TIMESTAMP NULL,
    last_update_type VARCHAR(100),
    
    -- Metadata
    metadata JSON,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    UNIQUE KEY uk_booking_tracking (tenant_id, booking_id),
    INDEX idx_tenant_tracking (tenant_id),
    INDEX idx_user_tracking (user_id),
    INDEX idx_booking_tracking (booking_id),
    INDEX idx_status_change (status_changed_at DESC),
    CONSTRAINT fk_tracking_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- PAYMENT & INVOICE TRACKING
-- ============================================

CREATE TABLE IF NOT EXISTS customer_payment_tracking (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    booking_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    
    -- Invoice Details
    invoice_number VARCHAR(100),
    invoice_date DATE,
    due_date DATE,
    
    -- Amount Details
    invoice_amount DECIMAL(15, 2),
    tax_amount DECIMAL(15, 2),
    total_amount DECIMAL(15, 2),
    
    -- Payment Status
    payment_status ENUM('pending', 'partial', 'paid', 'overdue', 'failed') DEFAULT 'pending',
    
    -- Payment Details
    payment_method ENUM('credit_card', 'debit_card', 'net_banking', 'upi', 'check', 'cash', 'ecs'),
    transaction_id VARCHAR(100),
    payment_date DATE,
    amount_paid DECIMAL(15, 2) DEFAULT 0,
    
    -- Refund Information
    is_refunded BOOLEAN DEFAULT FALSE,
    refund_amount DECIMAL(15, 2),
    refund_date DATE,
    refund_reason TEXT,
    
    -- Payment Reminder
    reminder_sent_count INT DEFAULT 0,
    last_reminder_sent_at TIMESTAMP NULL,
    
    -- Document Reference
    invoice_file_url VARCHAR(500),
    s3_bucket VARCHAR(100),
    s3_key VARCHAR(200),
    
    -- Metadata
    metadata JSON,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_payment (tenant_id),
    INDEX idx_booking_payment (booking_id),
    INDEX idx_user_payment (user_id),
    INDEX idx_payment_status (payment_status),
    INDEX idx_invoice_date (invoice_date DESC),
    CONSTRAINT fk_payment_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- CUSTOMER FEEDBACK & RATINGS
-- ============================================

CREATE TABLE IF NOT EXISTS customer_feedback (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    booking_id BIGINT,
    
    -- Feedback Details
    feedback_type ENUM('suggestion', 'complaint', 'praise', 'question', 'review'),
    subject VARCHAR(255),
    message TEXT NOT NULL,
    
    -- Rating
    overall_rating DECIMAL(2, 1),
    service_rating DECIMAL(2, 1),
    communication_rating DECIMAL(2, 1),
    documentation_rating DECIMAL(2, 1),
    
    -- Status
    feedback_status ENUM('open', 'acknowledged', 'in_progress', 'resolved', 'closed') DEFAULT 'open',
    response_message TEXT,
    responded_by BIGINT,
    responded_at TIMESTAMP NULL,
    
    -- Attachments
    attachments_count INT DEFAULT 0,
    
    -- Metadata
    metadata JSON,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_feedback (tenant_id),
    INDEX idx_user_feedback (user_id),
    INDEX idx_feedback_type (feedback_type),
    INDEX idx_feedback_status (feedback_status),
    INDEX idx_feedback_rating (overall_rating),
    CONSTRAINT fk_feedback_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- CUSTOMER ACTIVITY LOG
-- ============================================

CREATE TABLE IF NOT EXISTS customer_activity_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    
    -- Activity Details
    activity_type VARCHAR(100),
    activity_description TEXT,
    entity_type VARCHAR(100),
    entity_id BIGINT,
    
    -- Action Details
    action_taken VARCHAR(100),
    old_value TEXT,
    new_value TEXT,
    
    -- Session Information
    ip_address VARCHAR(45),
    user_agent TEXT,
    device_type VARCHAR(50),
    
    -- Metadata
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_tenant_activity (tenant_id),
    INDEX idx_user_activity (user_id),
    INDEX idx_activity_type (activity_type),
    INDEX idx_activity_created (created_at DESC),
    INDEX idx_entity_activity (entity_type, entity_id),
    CONSTRAINT fk_activity_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- CUSTOMER PREFERENCES & SETTINGS
-- ============================================

CREATE TABLE IF NOT EXISTS customer_preferences (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    
    -- Notification Preferences
    email_notifications BOOLEAN DEFAULT TRUE,
    sms_notifications BOOLEAN DEFAULT FALSE,
    push_notifications BOOLEAN DEFAULT TRUE,
    in_app_notifications BOOLEAN DEFAULT TRUE,
    
    -- Notification Categories
    receive_booking_updates BOOLEAN DEFAULT TRUE,
    receive_payment_reminders BOOLEAN DEFAULT TRUE,
    receive_document_requests BOOLEAN DEFAULT TRUE,
    receive_possession_updates BOOLEAN DEFAULT TRUE,
    receive_promotional_emails BOOLEAN DEFAULT FALSE,
    receive_newsletter BOOLEAN DEFAULT TRUE,
    
    -- Language & Localization
    preferred_language VARCHAR(50) DEFAULT 'en',
    preferred_timezone VARCHAR(100),
    
    -- Privacy Settings
    is_profile_public BOOLEAN DEFAULT FALSE,
    allow_contact_by_phone BOOLEAN DEFAULT TRUE,
    allow_contact_by_email BOOLEAN DEFAULT TRUE,
    allow_contact_by_sms BOOLEAN DEFAULT FALSE,
    allow_marketing BOOLEAN DEFAULT FALSE,
    
    -- Dashboard Customization
    dashboard_layout VARCHAR(50),
    theme_preference ENUM('light', 'dark', 'auto') DEFAULT 'auto',
    
    -- Metadata
    metadata JSON,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    UNIQUE KEY uk_customer_preferences (tenant_id, user_id),
    INDEX idx_tenant_preferences (tenant_id),
    CONSTRAINT fk_pref_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
