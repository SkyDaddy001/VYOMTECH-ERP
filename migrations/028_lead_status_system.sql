-- ============================================================
-- MIGRATION 028: LEAD STATUS SYSTEM ENHANCEMENT
-- Date: December 8, 2025
-- Purpose: Add detailed 30+ lead status system with pipeline stages
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- CREATE LEAD_STATUS_LOG TABLE - Track all status changes
-- ============================================================
CREATE TABLE IF NOT EXISTS `lead_status_log` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `lead_id` VARCHAR(36) NOT NULL,
    `old_status` VARCHAR(100),
    `new_status` VARCHAR(100) NOT NULL,
    `old_pipeline_stage` VARCHAR(50),
    `new_pipeline_stage` VARCHAR(50),
    `changed_by` VARCHAR(36),
    `change_reason` TEXT,
    `capture_date_type` VARCHAR(50),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`lead_id`) REFERENCES `sales_lead`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_lead` (`lead_id`),
    KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- CREATE LEAD_PIPELINE_CONFIG TABLE - Store pipeline configurations
-- ============================================================
CREATE TABLE IF NOT EXISTS `lead_pipeline_config` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `status` VARCHAR(100) NOT NULL,
    `pipeline_stage` VARCHAR(50) NOT NULL,
    `phase` VARCHAR(50),
    `color_code` VARCHAR(20),
    `icon` VARCHAR(50),
    `description` TEXT,
    `is_active` BOOLEAN DEFAULT TRUE,
    `sort_order` INT DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    UNIQUE KEY `uq_tenant_status` (`tenant_id`, `status`),
    KEY `idx_pipeline_stage` (`pipeline_stage`),
    KEY `idx_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- INSERT DEFAULT LEAD STATUSES FOR NEW TENANTS
-- ============================================================

-- Initial Contact Phase
INSERT IGNORE INTO `lead_pipeline_config` 
VALUES 
('lpc-fresh-lead', 'default-tenant', 'Fresh Lead', 'New', 'Initial Contact', '#3498db', 'star', 'Newly acquired lead', TRUE, 1, NOW(), NOW()),
('lpc-re-engaged', 'default-tenant', 'Re Engaged', 'Connected', 'Initial Contact', '#9b59b6', 'repeat', 'Previously contacted lead', TRUE, 2, NOW(), NOW()),
('lpc-not-connected', 'default-tenant', 'Not Connected', 'New', 'Initial Contact', '#e74c3c', 'x', 'Unable to connect', TRUE, 3, NOW(), NOW()),
('lpc-dead', 'default-tenant', 'Dead', 'New', 'Initial Contact', '#2c3e50', 'slash', 'No longer valid', TRUE, 4, NOW(), NOW()),
('lpc-follow-cold', 'default-tenant', 'Follow Up - Cold', 'Connected', 'Initial Contact', '#3498db', 'snowflake', 'Cold follow-up needed', TRUE, 5, NOW(), NOW()),
('lpc-follow-warm', 'default-tenant', 'Follow Up - Warm', 'Connected', 'Initial Contact', '#f39c12', 'flame', 'Warm follow-up needed', TRUE, 6, NOW(), NOW()),
('lpc-follow-hot', 'default-tenant', 'Follow Up - Hot', 'Interested', 'Initial Contact', '#e74c3c', 'fire', 'Hot follow-up needed', TRUE, 7, NOW(), NOW()),
('lpc-lost', 'default-tenant', 'Lost', 'Connected', 'Initial Contact', '#95a5a6', 'target', 'Lost opportunity', TRUE, 8, NOW(), NOW()),
('lpc-unq-location', 'default-tenant', 'Unqualified (Location)', 'Connected', 'Initial Contact', '#95a5a6', 'map-pin', 'Unqualified - location', TRUE, 9, NOW(), NOW()),
('lpc-unq-budget', 'default-tenant', 'Unqualified (Budget)', 'Connected', 'Initial Contact', '#95a5a6', 'credit', 'Unqualified - budget', TRUE, 10, NOW(), NOW()),
('lpc-unq-profile', 'default-tenant', 'Unqualified - Client Profile', 'Connected', 'Initial Contact', '#95a5a6', 'user-x', 'Unqualified - profile', TRUE, 11, NOW(), NOW()),

-- Site Visit Phase
('lpc-sv-scheduled', 'default-tenant', 'SV - Scheduled', 'Interested', 'Site Visit', '#3498db', 'calendar', 'Site visit scheduled', TRUE, 12, NOW(), NOW()),
('lpc-sv-done', 'default-tenant', 'SV - Done', 'Analysis', 'Site Visit', '#27ae60', 'check', 'Site visit completed', TRUE, 13, NOW(), NOW()),
('lpc-sv-cold', 'default-tenant', 'SV - Cold', 'Analysis', 'Site Visit', '#3498db', 'snowflake', 'Cold response to SV', TRUE, 14, NOW(), NOW()),
('lpc-sv-warm', 'default-tenant', 'SV - Warm', 'Analysis', 'Site Visit', '#f39c12', 'flame', 'Warm response to SV', TRUE, 15, NOW(), NOW()),
('lpc-sv-revisit-scheduled', 'default-tenant', 'SV - Revisit Scheduled', 'Analysis', 'Site Visit', '#3498db', 'calendar', 'Revisit scheduled', TRUE, 16, NOW(), NOW()),
('lpc-sv-revisit-done', 'default-tenant', 'SV - Revisit Done', 'Analysis', 'Site Visit', '#27ae60', 'check', 'Revisit completed', TRUE, 17, NOW(), NOW()),
('lpc-sv-hot', 'default-tenant', 'SV - Hot', 'Negotiation', 'Site Visit', '#e74c3c', 'fire', 'Hot response to SV', TRUE, 18, NOW(), NOW()),
('lpc-sv-lost-nr', 'default-tenant', 'SV - Lost (No Response)', 'Analysis', 'Site Visit', '#95a5a6', 'target', 'Lost - no response', TRUE, 19, NOW(), NOW()),
('lpc-sv-lost-budget', 'default-tenant', 'SV - Lost (Budget)', 'Analysis', 'Site Visit', '#95a5a6', 'credit', 'Lost - budget', TRUE, 20, NOW(), NOW()),
('lpc-sv-lost-plan', 'default-tenant', 'SV - Lost (Plan Dropped)', 'Analysis', 'Site Visit', '#95a5a6', 'x', 'Lost - plan dropped', TRUE, 21, NOW(), NOW()),
('lpc-sv-lost-location', 'default-tenant', 'SV - Lost (Location)', 'Analysis', 'Site Visit', '#95a5a6', 'map-pin', 'Lost - location', TRUE, 22, NOW(), NOW()),
('lpc-sv-lost-availability', 'default-tenant', 'SV - Lost (Availability)', 'Analysis', 'Site Visit', '#95a5a6', 'clock', 'Lost - availability', TRUE, 23, NOW(), NOW()),

-- Face-to-Face Phase
('lpc-f2f-scheduled', 'default-tenant', 'F2F - Scheduled', 'Interested', 'Face-to-Face', '#3498db', 'calendar', 'F2F meeting scheduled', TRUE, 24, NOW(), NOW()),
('lpc-f2f-done', 'default-tenant', 'F2F - Done', 'Analysis', 'Face-to-Face', '#27ae60', 'check', 'F2F meeting completed', TRUE, 25, NOW(), NOW()),
('lpc-f2f-followup', 'default-tenant', 'F2F - Follow Up', 'Analysis', 'Face-to-Face', '#3498db', 'repeat', 'F2F follow-up needed', TRUE, 26, NOW(), NOW()),
('lpc-f2f-warm', 'default-tenant', 'F2F - Warm', 'Analysis', 'Face-to-Face', '#f39c12', 'flame', 'Warm response to F2F', TRUE, 27, NOW(), NOW()),
('lpc-f2f-hot', 'default-tenant', 'F2F - Hot', 'Negotiation', 'Face-to-Face', '#e74c3c', 'fire', 'Hot response to F2F', TRUE, 28, NOW(), NOW()),
('lpc-f2f-lost-nr', 'default-tenant', 'F2F - Lost (No Response)', 'Analysis', 'Face-to-Face', '#95a5a6', 'target', 'Lost - no response', TRUE, 29, NOW(), NOW()),
('lpc-f2f-lost-budget', 'default-tenant', 'F2F - Lost (Budget)', 'Analysis', 'Face-to-Face', '#95a5a6', 'credit', 'Lost - budget', TRUE, 30, NOW(), NOW()),
('lpc-f2f-lost-plan', 'default-tenant', 'F2F - Lost (Plan Dropped)', 'Analysis', 'Face-to-Face', '#95a5a6', 'x', 'Lost - plan dropped', TRUE, 31, NOW(), NOW()),
('lpc-f2f-lost-location', 'default-tenant', 'F2F - Lost (Location)', 'Analysis', 'Face-to-Face', '#95a5a6', 'map-pin', 'Lost - location', TRUE, 32, NOW(), NOW()),
('lpc-f2f-lost-availability', 'default-tenant', 'F2F - Lost (Availability)', 'Analysis', 'Face-to-Face', '#95a5a6', 'clock', 'Lost - availability', TRUE, 33, NOW(), NOW()),

-- Booking Phase
('lpc-booking-progress', 'default-tenant', 'Booking-In-Progress', 'Pre Booking', 'Booking', '#3498db', 'bookmark', 'Booking in progress', TRUE, 34, NOW(), NOW()),
('lpc-booking-lost', 'default-tenant', 'Booking Progress-Lost', 'Pre Booking', 'Booking', '#95a5a6', 'target', 'Booking lost', TRUE, 35, NOW(), NOW()),
('lpc-booking-done', 'default-tenant', 'Booking Done', 'Booking', 'Booking', '#27ae60', 'check-circle', 'Booking completed', TRUE, 36, NOW(), NOW());

SET FOREIGN_KEY_CHECKS = 1;
