-- ============================================================
-- MIGRATION 009: CALL CENTER & AI AGENTS MODULE
-- Date: December 3, 2025
-- Purpose: Create call center, agents, leads, calls, and AI-related tables
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- AGENT TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `agent` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `agent_code` VARCHAR(50) UNIQUE NOT NULL,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100),
    `email` VARCHAR(255) UNIQUE,
    `phone` VARCHAR(20),
    `status` VARCHAR(50) DEFAULT 'available',
    `agent_type` VARCHAR(50),
    `skills` VARCHAR(255),
    `max_concurrent_calls` INT DEFAULT 1,
    `available` BOOLEAN DEFAULT TRUE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- CALL TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `call` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `lead_id` VARCHAR(36),
    `agent_id` VARCHAR(36),
    `call_status` VARCHAR(50) DEFAULT 'initiated',
    `duration_seconds` INT DEFAULT 0,
    `recording_url` TEXT,
    `notes` TEXT,
    `outcome` VARCHAR(50),
    `started_at` TIMESTAMP,
    `ended_at` TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`lead_id`) REFERENCES `sales_lead`(`id`),
    FOREIGN KEY (`agent_id`) REFERENCES `agent`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_lead` (`lead_id`),
    KEY `idx_agent` (`agent_id`),
    KEY `idx_status` (`call_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- CALL LOG TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `call_log` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `call_id` VARCHAR(36) NOT NULL,
    `event` VARCHAR(50),
    `details` TEXT,
    `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`call_id`) REFERENCES `call`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_call` (`call_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- CAMPAIGN TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `campaign` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `campaign_name` VARCHAR(255) NOT NULL,
    `campaign_type` VARCHAR(50),
    `description` TEXT,
    `start_date` DATE NOT NULL,
    `end_date` DATE,
    `status` VARCHAR(50) DEFAULT 'planning',
    `target_leads` INT DEFAULT 0,
    `budget` DECIMAL(18, 2),
    `assigned_agents` INT DEFAULT 0,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- AI MODEL TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `ai_model` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `model_name` VARCHAR(255) NOT NULL,
    `model_type` VARCHAR(50),
    `description` TEXT,
    `version` VARCHAR(50),
    `status` VARCHAR(50) DEFAULT 'active',
    `accuracy` DECIMAL(5, 2),
    `training_data_count` INT DEFAULT 0,
    `last_trained_at` TIMESTAMP,
    `parameters` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_type` (`model_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- AI INTERACTION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `ai_interaction` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `model_id` VARCHAR(36),
    `lead_id` VARCHAR(36),
    `interaction_type` VARCHAR(50),
    `input_text` LONGTEXT,
    `output_text` LONGTEXT,
    `confidence_score` DECIMAL(5, 2),
    `interaction_duration` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`model_id`) REFERENCES `ai_model`(`id`),
    FOREIGN KEY (`lead_id`) REFERENCES `sales_lead`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_model` (`model_id`),
    KEY `idx_lead` (`lead_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- COMMUNICATION TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `communication` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `lead_id` VARCHAR(36),
    `communication_type` VARCHAR(50),
    `channel` VARCHAR(50),
    `subject` VARCHAR(255),
    `message` LONGTEXT,
    `from_address` VARCHAR(255),
    `to_address` VARCHAR(255),
    `status` VARCHAR(50),
    `sent_at` TIMESTAMP,
    `opened_at` TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`lead_id`) REFERENCES `sales_lead`(`id`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_lead` (`lead_id`),
    KEY `idx_channel` (`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- AGENT PERFORMANCE TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS `agent_performance` (
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `agent_id` VARCHAR(36) NOT NULL,
    `performance_date` DATE NOT NULL,
    `total_calls` INT DEFAULT 0,
    `successful_calls` INT DEFAULT 0,
    `average_duration` INT DEFAULT 0,
    `conversion_rate` DECIMAL(5, 2) DEFAULT 0,
    `customer_satisfaction` DECIMAL(5, 2),
    `sales_generated` DECIMAL(18, 2) DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`agent_id`) REFERENCES `agent`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_agent` (`agent_id`),
    KEY `idx_date` (`performance_date`),
    UNIQUE KEY `unique_performance` (`agent_id`, `performance_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
