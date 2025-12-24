-- CreateTable tenant
CREATE TABLE `tenant` (
  `id` CHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `domain` VARCHAR(255) NOT NULL,
  `status` VARCHAR(50) NOT NULL DEFAULT 'active',
  `max_users` INT NOT NULL DEFAULT 100,
  `max_concurrent_calls` INT NOT NULL DEFAULT 50,
  `ai_budget_monthly` DECIMAL(15,2) NOT NULL DEFAULT 1000.00,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL,

  UNIQUE INDEX `tenant_domain_key`(`domain`),
  INDEX `tenant_domain_idx`(`domain`),
  INDEX `tenant_status_idx`(`status`),
  PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable user
CREATE TABLE `user` (
  `id` CHAR(36) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `role` VARCHAR(50) NOT NULL DEFAULT 'user',
  `tenant_id` VARCHAR(36) NOT NULL,
  `current_tenant_id` VARCHAR(36),
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL,

  UNIQUE INDEX `user_email_key`(`email`),
  INDEX `user_email_idx`(`email`),
  INDEX `user_tenantId_idx`(`tenant_id`),
  INDEX `user_role_idx`(`role`),
  PRIMARY KEY (`id`),
  CONSTRAINT `user_tenantId_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `user_currentTenantId_fk` FOREIGN KEY (`current_tenant_id`) REFERENCES `tenant` (`id`) ON DELETE SET NULL
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable password_reset_token
CREATE TABLE `password_reset_token` (
  `id` CHAR(36) NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `token` VARCHAR(255) NOT NULL,
  `expires_at` DATETIME(3) NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  UNIQUE INDEX `password_reset_token_token_key`(`token`),
  INDEX `password_reset_token_userId_idx`(`user_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `password_reset_token_userId_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable auth_token
CREATE TABLE `auth_token` (
  `id` CHAR(36) NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `token` VARCHAR(255) NOT NULL,
  `expires_at` DATETIME(3) NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  UNIQUE INDEX `auth_token_token_key`(`token`),
  INDEX `auth_token_userId_idx`(`user_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `auth_token_userId_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable audit_log
CREATE TABLE `audit_log` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `user_id` CHAR(36),
  `action` VARCHAR(50) NOT NULL,
  `resource_type` VARCHAR(100) NOT NULL,
  `resource_id` VARCHAR(255),
  `changes` JSON,
  `ip_address` VARCHAR(45),
  `user_agent` VARCHAR(255),
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  INDEX `audit_log_tenantId_idx`(`tenant_id`),
  INDEX `audit_log_userId_idx`(`user_id`),
  INDEX `audit_log_action_idx`(`action`),
  INDEX `audit_log_resourceType_idx`(`resource_type`),
  PRIMARY KEY (`id`),
  CONSTRAINT `audit_log_tenantId_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE,
  CONSTRAINT `audit_log_userId_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE SET NULL
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable team
CREATE TABLE `team` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `description` TEXT,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL,

  INDEX `team_tenantId_idx`(`tenant_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `team_tenantId_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable system_config
CREATE TABLE `system_config` (
  `id` CHAR(26) NOT NULL,
  `tenant_id` VARCHAR(36) NOT NULL,
  `config_key` VARCHAR(255) NOT NULL,
  `config_value` LONGTEXT NOT NULL,
  `data_type` VARCHAR(50) NOT NULL,
  `description` TEXT,
  `is_global` BOOLEAN NOT NULL DEFAULT false,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL,

  UNIQUE INDEX `system_config_tenantId_configKey_key`(`tenant_id`, `config_key`),
  INDEX `system_config_tenantId_idx`(`tenant_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `system_config_tenantId_fk` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Note: Additional tables (SalesLead, Booking, Property, Document, Loan, Employee, Account, etc.)
-- have been created by the migrations in the /migrations directory
-- This is the baseline migration for core tables
