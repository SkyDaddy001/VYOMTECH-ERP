-- ============================================================
-- MIGRATION 040: MOBILE APP FEATURES
-- Date: December 23, 2025
-- Purpose: Multi-tenant mobile app infrastructure with push notifications,
--          device management, offline data caching, and crash reporting
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- MOBILE APP CONFIGURATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS `mobile_app` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `app_name` VARCHAR(255) NOT NULL,
    `app_type` VARCHAR(50) NOT NULL DEFAULT 'ios', -- ios, android, cross-platform
    `bundle_identifier` VARCHAR(255) UNIQUE NOT NULL,
    `version` VARCHAR(50) NOT NULL,
    `build_number` INT NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active', -- active, inactive, deprecated, beta
    `description` TEXT,
    `api_key` VARCHAR(255) UNIQUE NOT NULL,
    `api_secret` VARCHAR(255) NOT NULL,
    `min_supported_version` VARCHAR(50),
    `max_supported_version` VARCHAR(50),
    `store_url_ios` VARCHAR(500),
    `store_url_android` VARCHAR(500),
    `app_icon_url` VARCHAR(500),
    `app_banner_url` VARCHAR(500),
    `support_email` VARCHAR(255),
    `support_phone` VARCHAR(20),
    `support_chat_url` VARCHAR(500),
    `privacy_policy_url` VARCHAR(500),
    `terms_of_service_url` VARCHAR(500),
    `changelog_url` VARCHAR(500),
    `feature_flags` JSON,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    UNIQUE KEY `uk_tenant_bundle` (`tenant_id`, `bundle_identifier`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_status` (`status`),
    KEY `idx_api_key` (`api_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- MOBILE DEVICE REGISTRATION
-- ============================================================
CREATE TABLE IF NOT EXISTS `mobile_device` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `app_id` VARCHAR(36) NOT NULL,
    `device_id` VARCHAR(255) UNIQUE NOT NULL,
    `device_name` VARCHAR(255),
    `device_model` VARCHAR(255),
    `device_manufacturer` VARCHAR(255),
    `os_type` VARCHAR(50) NOT NULL, -- ios, android, windows
    `os_version` VARCHAR(50),
    `app_version` VARCHAR(50),
    `app_build` INT,
    `push_token` VARCHAR(500),
    `push_token_updated_at` TIMESTAMP NULL,
    `device_uuid` VARCHAR(36),
    `imei` VARCHAR(50),
    `device_serial` VARCHAR(255),
    `screen_resolution` VARCHAR(50),
    `screen_density` VARCHAR(50),
    `locale` VARCHAR(10),
    `timezone` VARCHAR(100),
    `battery_optimization_enabled` BOOLEAN DEFAULT 0,
    `biometric_enabled` BOOLEAN DEFAULT 0,
    `biometric_type` VARCHAR(50), -- fingerprint, face, iris
    `fcm_token` VARCHAR(500), -- Firebase Cloud Messaging token
    `apns_token` VARCHAR(500), -- Apple Push Notification service token
    `device_status` VARCHAR(50) NOT NULL DEFAULT 'active', -- active, inactive, suspended, lost
    `last_location` POINT,
    `last_activity_at` TIMESTAMP NULL,
    `first_seen_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `last_seen_at` TIMESTAMP NULL,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`app_id`) REFERENCES `mobile_app`(`id`) ON DELETE CASCADE,
    UNIQUE KEY `uk_tenant_device` (`tenant_id`, `device_id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_device_status` (`device_status`),
    KEY `idx_push_token` (`push_token`),
    KEY `idx_last_activity` (`last_activity_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- MOBILE USER SESSIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS `mobile_session` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `device_id` VARCHAR(36) NOT NULL,
    `app_id` VARCHAR(36) NOT NULL,
    `session_token` VARCHAR(500) UNIQUE NOT NULL,
    `refresh_token` VARCHAR(500) UNIQUE,
    `token_expires_at` TIMESTAMP NOT NULL,
    `refresh_token_expires_at` TIMESTAMP,
    `session_status` VARCHAR(50) NOT NULL DEFAULT 'active', -- active, inactive, expired, revoked
    `ip_address` VARCHAR(45),
    `user_agent` TEXT,
    `app_version` VARCHAR(50),
    `login_method` VARCHAR(50), -- password, sso, biometric, otp
    `login_timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `last_activity_timestamp` TIMESTAMP NULL,
    `logout_timestamp` TIMESTAMP NULL,
    `session_duration_seconds` INT,
    `device_trusted` BOOLEAN DEFAULT 0,
    `two_factor_verified` BOOLEAN DEFAULT 0,
    `geo_location` POINT,
    `network_type` VARCHAR(50), -- wifi, cellular, unknown
    `network_provider` VARCHAR(255),
    `app_in_foreground` BOOLEAN DEFAULT 1,
    `background_activity_allowed` BOOLEAN DEFAULT 1,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`device_id`) REFERENCES `mobile_device`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`app_id`) REFERENCES `mobile_app`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_device_id` (`device_id`),
    KEY `idx_session_token` (`session_token`),
    KEY `idx_session_status` (`session_status`),
    KEY `idx_login_timestamp` (`login_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- APP PUSH NOTIFICATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS `app_notification` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `app_id` VARCHAR(36) NOT NULL,
    `user_id` INT,
    `device_id` VARCHAR(36),
    `notification_type` VARCHAR(100) NOT NULL, -- lead_update, payment_reminder, booking_confirmation, etc.
    `notification_category` VARCHAR(50), -- system, promotional, transaction, reminder
    `title` VARCHAR(255) NOT NULL,
    `body` TEXT NOT NULL,
    `content_type` VARCHAR(50) DEFAULT 'text', -- text, html, rich_media
    `image_url` VARCHAR(500),
    `action_url` VARCHAR(500),
    `action_type` VARCHAR(50), -- deeplink, external, internal
    `priority` VARCHAR(50) DEFAULT 'normal', -- low, normal, high, critical
    `sound` VARCHAR(100),
    `badge_count` INT DEFAULT 1,
    `custom_data` JSON,
    `scheduled_time` TIMESTAMP,
    `send_time` TIMESTAMP NULL,
    `delivery_status` VARCHAR(50) DEFAULT 'pending', -- pending, sent, delivered, failed, read, deleted
    `delivery_attempts` INT DEFAULT 0,
    `next_retry_at` TIMESTAMP NULL,
    `read_at` TIMESTAMP NULL,
    `clicked_at` TIMESTAMP NULL,
    `dismissed_at` TIMESTAMP NULL,
    `expiry_time` TIMESTAMP NULL,
    `mute_until` TIMESTAMP NULL,
    `is_campaign` BOOLEAN DEFAULT 0,
    `campaign_id` VARCHAR(36),
    `segment_id` VARCHAR(36),
    `ab_test_variant` VARCHAR(100),
    `retry_count` INT DEFAULT 0,
    `max_retries` INT DEFAULT 3,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`app_id`) REFERENCES `mobile_app`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE SET NULL,
    FOREIGN KEY (`device_id`) REFERENCES `mobile_device`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_device_id` (`device_id`),
    KEY `idx_notification_type` (`notification_type`),
    KEY `idx_delivery_status` (`delivery_status`),
    KEY `idx_send_time` (`send_time`),
    KEY `idx_scheduled_time` (`scheduled_time`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- MOBILE APP FEATURES & FEATURE FLAGS
-- ============================================================
CREATE TABLE IF NOT EXISTS `app_feature` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `app_id` VARCHAR(36) NOT NULL,
    `feature_name` VARCHAR(255) NOT NULL,
    `feature_code` VARCHAR(100) NOT NULL,
    `feature_description` TEXT,
    `feature_category` VARCHAR(100), -- core, premium, beta, experimental
    `is_enabled` BOOLEAN DEFAULT 1,
    `min_app_version` VARCHAR(50),
    `max_app_version` VARCHAR(50),
    `min_os_version` VARCHAR(50),
    `required_permissions` JSON,
    `enabled_for_users` JSON, -- list of user IDs or roles
    `disabled_for_users` JSON,
    `rollout_percentage` INT DEFAULT 100, -- 0-100 for gradual rollout
    `ab_test_variant` VARCHAR(100),
    `analytics_track` BOOLEAN DEFAULT 1,
    `requires_consent` BOOLEAN DEFAULT 0,
    `consent_type` VARCHAR(50), -- gdpr, privacy, location
    `feature_url` VARCHAR(500),
    `documentation_url` VARCHAR(500),
    `support_email` VARCHAR(255),
    `launch_date` TIMESTAMP NULL,
    `sunset_date` TIMESTAMP NULL,
    `config` JSON,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`app_id`) REFERENCES `mobile_app`(`id`) ON DELETE CASCADE,
    UNIQUE KEY `uk_tenant_app_feature` (`tenant_id`, `app_id`, `feature_code`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_is_enabled` (`is_enabled`),
    KEY `idx_feature_category` (`feature_category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- MOBILE APP USER SETTINGS & PREFERENCES
-- ============================================================
CREATE TABLE IF NOT EXISTS `app_setting` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `app_id` VARCHAR(36) NOT NULL,
    `device_id` VARCHAR(36),
    `setting_key` VARCHAR(255) NOT NULL,
    `setting_value` LONGTEXT,
    `setting_type` VARCHAR(50), -- string, boolean, integer, json, array
    `display_name` VARCHAR(255),
    `description` TEXT,
    `category` VARCHAR(100), -- notification, privacy, security, ui, performance
    `is_user_editable` BOOLEAN DEFAULT 1,
    `is_device_specific` BOOLEAN DEFAULT 0,
    `default_value` LONGTEXT,
    `validation_rules` JSON,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`app_id`) REFERENCES `mobile_app`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`device_id`) REFERENCES `mobile_device`(`id`) ON DELETE SET NULL,
    UNIQUE KEY `uk_user_app_setting` (`user_id`, `app_id`, `setting_key`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_device_id` (`device_id`),
    KEY `idx_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- OFFLINE DATA CACHE
-- ============================================================
CREATE TABLE IF NOT EXISTS `offline_data` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `device_id` VARCHAR(36) NOT NULL,
    `app_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `data_type` VARCHAR(100) NOT NULL, -- leads, contacts, products, transactions
    `data_key` VARCHAR(255) NOT NULL,
    `data_source_table` VARCHAR(100), -- reference to original table
    `source_record_id` VARCHAR(36),
    `cached_data` LONGBLOB NOT NULL,
    `data_hash` VARCHAR(64), -- SHA-256 hash for integrity check
    `compression_type` VARCHAR(50), -- gzip, brotli, none
    `sync_status` VARCHAR(50) DEFAULT 'pending', -- pending, syncing, synced, failed, conflict
    `last_sync_at` TIMESTAMP NULL,
    `needs_sync` BOOLEAN DEFAULT 1,
    `local_only` BOOLEAN DEFAULT 0,
    `cache_priority` VARCHAR(50) DEFAULT 'normal', -- low, normal, high, critical
    `expiry_time` TIMESTAMP NULL,
    `size_bytes` BIGINT,
    `conflict_resolution_strategy` VARCHAR(50), -- last_write_wins, server_priority, user_choice
    `conflicts` JSON,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`device_id`) REFERENCES `mobile_device`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`app_id`) REFERENCES `mobile_app`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    UNIQUE KEY `uk_device_data_key` (`device_id`, `app_id`, `data_key`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_device_id` (`device_id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_sync_status` (`sync_status`),
    KEY `idx_needs_sync` (`needs_sync`),
    KEY `idx_data_type` (`data_type`),
    KEY `idx_expiry_time` (`expiry_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- MOBILE APP CRASH REPORTING
-- ============================================================
CREATE TABLE IF NOT EXISTS `app_crash` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `app_id` VARCHAR(36) NOT NULL,
    `device_id` VARCHAR(36) NOT NULL,
    `user_id` INT,
    `crash_timestamp` TIMESTAMP NOT NULL,
    `app_version` VARCHAR(50),
    `app_build` INT,
    `os_version` VARCHAR(50),
    `device_model` VARCHAR(255),
    `crash_type` VARCHAR(100), -- exception, anr, force_close, oom, segfault
    `crash_reason` VARCHAR(255),
    `exception_type` VARCHAR(255),
    `exception_message` TEXT,
    `stack_trace` LONGTEXT,
    `breadcrumbs` JSON,
    `user_report` TEXT,
    `severity` VARCHAR(50) DEFAULT 'medium', -- low, medium, high, critical
    `is_reproducible` BOOLEAN DEFAULT 0,
    `reproduction_steps` TEXT,
    `memory_used_mb` INT,
    `memory_available_mb` INT,
    `battery_level` INT,
    `battery_health` VARCHAR(50),
    `storage_available_mb` INT,
    `network_type` VARCHAR(50),
    `cpu_usage_percent` INT,
    `temperature_celsius` INT,
    `session_id` VARCHAR(36),
    `crash_status` VARCHAR(50) DEFAULT 'new', -- new, reviewing, acknowledged, fixed, won_t_fix, duplicate
    `assigned_to` INT,
    `assigned_at` TIMESTAMP NULL,
    `resolved_at` TIMESTAMP NULL,
    `resolution_notes` TEXT,
    `related_issue_id` VARCHAR(36),
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`app_id`) REFERENCES `mobile_app`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`device_id`) REFERENCES `mobile_device`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_device_id` (`device_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_crash_timestamp` (`crash_timestamp`),
    KEY `idx_crash_type` (`crash_type`),
    KEY `idx_severity` (`severity`),
    KEY `idx_crash_status` (`crash_status`),
    KEY `idx_assigned_to` (`assigned_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- MOBILE DEVICE ANALYTICS
-- ============================================================
CREATE TABLE IF NOT EXISTS `device_analytic` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `app_id` VARCHAR(36) NOT NULL,
    `device_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `event_type` VARCHAR(100) NOT NULL, -- screen_view, button_click, form_submit, etc.
    `event_name` VARCHAR(255),
    `event_category` VARCHAR(100),
    `event_action` VARCHAR(100),
    `event_label` VARCHAR(255),
    `event_value` DECIMAL(15, 2),
    `screen_name` VARCHAR(255),
    `screen_class` VARCHAR(255),
    `page_title` VARCHAR(255),
    `page_path` VARCHAR(500),
    `session_id` VARCHAR(36),
    `event_id` VARCHAR(36) UNIQUE,
    `event_timestamp` TIMESTAMP NOT NULL,
    `session_start_time` TIMESTAMP,
    `time_on_screen_seconds` INT,
    `engagement_time_seconds` INT,
    `scroll_depth_percent` INT,
    `click_count` INT,
    `form_completion_percent` INT,
    `errors_encountered` INT,
    `referrer` VARCHAR(500),
    `utm_source` VARCHAR(100),
    `utm_medium` VARCHAR(100),
    `utm_campaign` VARCHAR(100),
    `utm_content` VARCHAR(100),
    `utm_term` VARCHAR(100),
    `custom_params` JSON,
    `device_info` JSON,
    `performance_metrics` JSON,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`app_id`) REFERENCES `mobile_app`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`device_id`) REFERENCES `mobile_device`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_device_id` (`device_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_event_type` (`event_type`),
    KEY `idx_event_timestamp` (`event_timestamp`),
    KEY `idx_session_id` (`session_id`),
    KEY `idx_screen_name` (`screen_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- MOBILE APP UPDATES & VERSION TRACKING
-- ============================================================
CREATE TABLE IF NOT EXISTS `app_update` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `app_id` VARCHAR(36) NOT NULL,
    `version` VARCHAR(50) NOT NULL,
    `build_number` INT NOT NULL,
    `release_date` TIMESTAMP NOT NULL,
    `update_type` VARCHAR(50) NOT NULL DEFAULT 'minor', -- major, minor, patch, hotfix, beta
    `is_mandatory` BOOLEAN DEFAULT 0,
    `is_rollback_available` BOOLEAN DEFAULT 0,
    `rollback_version` VARCHAR(50),
    `update_title` VARCHAR(255),
    `update_description` TEXT,
    `changelog` LONGTEXT,
    `release_notes` LONGTEXT,
    `update_size_mb` DECIMAL(10, 2),
    `min_os_version` VARCHAR(50),
    `max_os_version` VARCHAR(50),
    `download_url_ios` VARCHAR(500),
    `download_url_android` VARCHAR(500),
    `download_url_windows` VARCHAR(500),
    `checksum_sha256` VARCHAR(64),
    `install_instructions` TEXT,
    `breaking_changes` TEXT,
    `deprecated_features` JSON,
    `new_features` JSON,
    `bug_fixes` JSON,
    `performance_improvements` JSON,
    `security_updates` JSON,
    `update_stage` VARCHAR(50) DEFAULT 'beta', -- beta, staged_rollout, production, deprecated
    `rollout_percentage` INT DEFAULT 100,
    `rollout_start_time` TIMESTAMP NULL,
    `rollout_end_time` TIMESTAMP NULL,
    `install_count` INT DEFAULT 0,
    `uninstall_count` INT DEFAULT 0,
    `rollback_count` INT DEFAULT 0,
    `update_status` VARCHAR(50) DEFAULT 'active', -- active, deprecated, archived
    `deprecation_reason` TEXT,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`app_id`) REFERENCES `mobile_app`(`id`) ON DELETE CASCADE,
    UNIQUE KEY `uk_tenant_app_version` (`tenant_id`, `app_id`, `version`, `build_number`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_release_date` (`release_date`),
    KEY `idx_update_type` (`update_type`),
    KEY `idx_update_stage` (`update_stage`),
    KEY `idx_update_status` (`update_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- DEVICE UPDATE HISTORY
-- ============================================================
CREATE TABLE IF NOT EXISTS `device_update_history` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `device_id` VARCHAR(36) NOT NULL,
    `app_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `update_id` VARCHAR(36) NOT NULL,
    `previous_version` VARCHAR(50),
    `new_version` VARCHAR(50) NOT NULL,
    `update_start_time` TIMESTAMP,
    `update_completion_time` TIMESTAMP NULL,
    `update_duration_seconds` INT,
    `update_method` VARCHAR(50), -- auto, manual, scheduled, forced
    `update_status` VARCHAR(50) DEFAULT 'pending', -- pending, downloading, installing, completed, failed, rolled_back
    `failure_reason` TEXT,
    `installation_log` LONGTEXT,
    `device_restarted` BOOLEAN DEFAULT 0,
    `restart_time` TIMESTAMP NULL,
    `post_install_errors` JSON,
    `user_feedback` TEXT,
    `user_rating` INT,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`device_id`) REFERENCES `mobile_device`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`app_id`) REFERENCES `mobile_app`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_device_id` (`device_id`),
    KEY `idx_app_id` (`app_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_update_id` (`update_id`),
    KEY `idx_update_status` (`update_status`),
    KEY `idx_update_start_time` (`update_start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
