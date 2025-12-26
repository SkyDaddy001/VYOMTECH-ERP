/*
  Warnings:

  - The primary key for the `audit_log` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `createdAt` on the `audit_log` table. All the data in the column will be lost.
  - You are about to drop the column `entityId` on the `audit_log` table. All the data in the column will be lost.
  - You are about to drop the column `entityType` on the `audit_log` table. All the data in the column will be lost.
  - You are about to drop the column `metadata` on the `audit_log` table. All the data in the column will be lost.
  - You are about to drop the column `tenantId` on the `audit_log` table. All the data in the column will be lost.
  - You are about to drop the column `userId` on the `audit_log` table. All the data in the column will be lost.
  - You are about to alter the column `id` on the `audit_log` table. The data in that column could be lost. The data in that column will be cast from `Char(36)` to `Char(26)`.
  - You are about to alter the column `changes` on the `audit_log` table. The data in that column could be lost. The data in that column will be cast from `LongText` to `Json`.
  - The primary key for the `auth_token` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `createdAt` on the `auth_token` table. All the data in the column will be lost.
  - You are about to drop the column `expiresAt` on the `auth_token` table. All the data in the column will be lost.
  - You are about to drop the column `revokedAt` on the `auth_token` table. All the data in the column will be lost.
  - You are about to drop the column `tokenType` on the `auth_token` table. All the data in the column will be lost.
  - You are about to drop the column `userId` on the `auth_token` table. All the data in the column will be lost.
  - You are about to alter the column `id` on the `auth_token` table. The data in that column could be lost. The data in that column will be cast from `Char(36)` to `Char(26)`.
  - You are about to drop the column `createdAt` on the `permission` table. All the data in the column will be lost.
  - You are about to drop the column `roleId` on the `permission` table. All the data in the column will be lost.
  - You are about to alter the column `action` on the `permission` table. The data in that column could be lost. The data in that column will be cast from `VarChar(100)` to `VarChar(50)`.
  - You are about to drop the column `amount` on the `quotation` table. All the data in the column will be lost.
  - You are about to drop the column `createdAt` on the `quotation` table. All the data in the column will be lost.
  - You are about to drop the column `leadId` on the `quotation` table. All the data in the column will be lost.
  - You are about to drop the column `quotationNumber` on the `quotation` table. All the data in the column will be lost.
  - You are about to drop the column `tenantId` on the `quotation` table. All the data in the column will be lost.
  - You are about to drop the column `updatedAt` on the `quotation` table. All the data in the column will be lost.
  - You are about to drop the column `validUntil` on the `quotation` table. All the data in the column will be lost.
  - You are about to drop the column `createdAt` on the `role` table. All the data in the column will be lost.
  - You are about to drop the column `name` on the `role` table. All the data in the column will be lost.
  - You are about to drop the column `tenantId` on the `role` table. All the data in the column will be lost.
  - You are about to drop the column `assignedTo` on the `sales_lead` table. All the data in the column will be lost.
  - You are about to drop the column `createdAt` on the `sales_lead` table. All the data in the column will be lost.
  - You are about to drop the column `deletedAt` on the `sales_lead` table. All the data in the column will be lost.
  - You are about to drop the column `leadCode` on the `sales_lead` table. All the data in the column will be lost.
  - You are about to drop the column `leadName` on the `sales_lead` table. All the data in the column will be lost.
  - You are about to drop the column `leadSource` on the `sales_lead` table. All the data in the column will be lost.
  - You are about to drop the column `leadValue` on the `sales_lead` table. All the data in the column will be lost.
  - You are about to drop the column `status` on the `sales_lead` table. All the data in the column will be lost.
  - You are about to drop the column `tenantId` on the `sales_lead` table. All the data in the column will be lost.
  - You are about to drop the column `updatedAt` on the `sales_lead` table. All the data in the column will be lost.
  - You are about to drop the column `createdAt` on the `tenant` table. All the data in the column will be lost.
  - You are about to drop the column `maxUsers` on the `tenant` table. All the data in the column will be lost.
  - You are about to drop the column `updatedAt` on the `tenant` table. All the data in the column will be lost.
  - You are about to drop the column `createdAt` on the `user` table. All the data in the column will be lost.
  - You are about to drop the column `firstName` on the `user` table. All the data in the column will be lost.
  - You are about to drop the column `lastLogin` on the `user` table. All the data in the column will be lost.
  - You are about to drop the column `lastName` on the `user` table. All the data in the column will be lost.
  - You are about to drop the column `passwordHash` on the `user` table. All the data in the column will be lost.
  - You are about to drop the column `status` on the `user` table. All the data in the column will be lost.
  - You are about to drop the column `tenantId` on the `user` table. All the data in the column will be lost.
  - You are about to drop the column `updatedAt` on the `user` table. All the data in the column will be lost.
  - You are about to drop the column `assignedAt` on the `user_role` table. All the data in the column will be lost.
  - You are about to drop the column `roleId` on the `user_role` table. All the data in the column will be lost.
  - You are about to drop the column `userId` on the `user_role` table. All the data in the column will be lost.
  - You are about to drop the `lead_status_history` table. If the table is not empty, all the data it contains will be lost.
  - A unique constraint covering the columns `[token]` on the table `auth_token` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[tenant_id,permission_name]` on the table `permission` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[quotation_number]` on the table `quotation` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[tenant_id,role_name]` on the table `role` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[lead_code]` on the table `sales_lead` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[email]` on the table `user` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[user_id,role_id]` on the table `user_role` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `entity_type` to the `audit_log` table without a default value. This is not possible if the table is not empty.
  - Added the required column `tenant_id` to the `audit_log` table without a default value. This is not possible if the table is not empty.
  - Added the required column `expires_at` to the `auth_token` table without a default value. This is not possible if the table is not empty.
  - Added the required column `user_id` to the `auth_token` table without a default value. This is not possible if the table is not empty.
  - Added the required column `permission_name` to the `permission` table without a default value. This is not possible if the table is not empty.
  - Added the required column `tenant_id` to the `permission` table without a default value. This is not possible if the table is not empty.
  - Added the required column `updated_at` to the `permission` table without a default value. This is not possible if the table is not empty.
  - Added the required column `quotation_date` to the `quotation` table without a default value. This is not possible if the table is not empty.
  - Added the required column `quotation_number` to the `quotation` table without a default value. This is not possible if the table is not empty.
  - Added the required column `tenant_id` to the `quotation` table without a default value. This is not possible if the table is not empty.
  - Added the required column `updated_at` to the `quotation` table without a default value. This is not possible if the table is not empty.
  - Added the required column `valid_till` to the `quotation` table without a default value. This is not possible if the table is not empty.
  - Added the required column `role_name` to the `role` table without a default value. This is not possible if the table is not empty.
  - Added the required column `tenant_id` to the `role` table without a default value. This is not possible if the table is not empty.
  - Added the required column `updated_at` to the `role` table without a default value. This is not possible if the table is not empty.
  - Added the required column `lead_code` to the `sales_lead` table without a default value. This is not possible if the table is not empty.
  - Added the required column `lead_name` to the `sales_lead` table without a default value. This is not possible if the table is not empty.
  - Added the required column `tenant_id` to the `sales_lead` table without a default value. This is not possible if the table is not empty.
  - Added the required column `updated_at` to the `sales_lead` table without a default value. This is not possible if the table is not empty.
  - Added the required column `updated_at` to the `tenant` table without a default value. This is not possible if the table is not empty.
  - Added the required column `password_hash` to the `user` table without a default value. This is not possible if the table is not empty.
  - Added the required column `tenant_id` to the `user` table without a default value. This is not possible if the table is not empty.
  - Added the required column `updated_at` to the `user` table without a default value. This is not possible if the table is not empty.
  - Added the required column `role_id` to the `user_role` table without a default value. This is not possible if the table is not empty.
  - Added the required column `tenant_id` to the `user_role` table without a default value. This is not possible if the table is not empty.
  - Added the required column `user_id` to the `user_role` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE `audit_log` DROP FOREIGN KEY `audit_log_entityId_fkey`;

-- DropForeignKey
ALTER TABLE `audit_log` DROP FOREIGN KEY `audit_log_tenantId_fkey`;

-- DropForeignKey
ALTER TABLE `audit_log` DROP FOREIGN KEY `audit_log_userId_fkey`;

-- DropForeignKey
ALTER TABLE `auth_token` DROP FOREIGN KEY `auth_token_userId_fkey`;

-- DropForeignKey
ALTER TABLE `lead_status_history` DROP FOREIGN KEY `lead_status_history_leadId_fkey`;

-- DropForeignKey
ALTER TABLE `permission` DROP FOREIGN KEY `permission_roleId_fkey`;

-- DropForeignKey
ALTER TABLE `quotation` DROP FOREIGN KEY `quotation_leadId_fkey`;

-- DropForeignKey
ALTER TABLE `quotation` DROP FOREIGN KEY `quotation_tenantId_fkey`;

-- DropForeignKey
ALTER TABLE `role` DROP FOREIGN KEY `role_tenantId_fkey`;

-- DropForeignKey
ALTER TABLE `sales_lead` DROP FOREIGN KEY `sales_lead_tenantId_fkey`;

-- DropForeignKey
ALTER TABLE `user` DROP FOREIGN KEY `user_tenantId_fkey`;

-- DropForeignKey
ALTER TABLE `user_role` DROP FOREIGN KEY `user_role_roleId_fkey`;

-- DropForeignKey
ALTER TABLE `user_role` DROP FOREIGN KEY `user_role_userId_fkey`;

-- DropIndex
DROP INDEX `audit_log_entityId_fkey` ON `audit_log`;

-- DropIndex
DROP INDEX `audit_log_entityType_entityId_idx` ON `audit_log`;

-- DropIndex
DROP INDEX `audit_log_tenantId_idx` ON `audit_log`;

-- DropIndex
DROP INDEX `audit_log_userId_idx` ON `audit_log`;

-- DropIndex
DROP INDEX `auth_token_userId_idx` ON `auth_token`;

-- DropIndex
DROP INDEX `permission_roleId_action_resource_key` ON `permission`;

-- DropIndex
DROP INDEX `permission_roleId_idx` ON `permission`;

-- DropIndex
DROP INDEX `quotation_leadId_idx` ON `quotation`;

-- DropIndex
DROP INDEX `quotation_quotationNumber_key` ON `quotation`;

-- DropIndex
DROP INDEX `quotation_tenantId_idx` ON `quotation`;

-- DropIndex
DROP INDEX `role_tenantId_idx` ON `role`;

-- DropIndex
DROP INDEX `role_tenantId_name_key` ON `role`;

-- DropIndex
DROP INDEX `sales_lead_assignedTo_idx` ON `sales_lead`;

-- DropIndex
DROP INDEX `sales_lead_leadCode_key` ON `sales_lead`;

-- DropIndex
DROP INDEX `sales_lead_status_idx` ON `sales_lead`;

-- DropIndex
DROP INDEX `sales_lead_tenantId_idx` ON `sales_lead`;

-- DropIndex
DROP INDEX `user_tenantId_email_key` ON `user`;

-- DropIndex
DROP INDEX `user_tenantId_idx` ON `user`;

-- DropIndex
DROP INDEX `user_role_roleId_idx` ON `user_role`;

-- DropIndex
DROP INDEX `user_role_userId_idx` ON `user_role`;

-- DropIndex
DROP INDEX `user_role_userId_roleId_key` ON `user_role`;

-- AlterTable
ALTER TABLE `audit_log` DROP PRIMARY KEY,
    DROP COLUMN `createdAt`,
    DROP COLUMN `entityId`,
    DROP COLUMN `entityType`,
    DROP COLUMN `metadata`,
    DROP COLUMN `tenantId`,
    DROP COLUMN `userId`,
    ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `entity_id` VARCHAR(36) NULL,
    ADD COLUMN `entity_type` VARCHAR(100) NOT NULL,
    ADD COLUMN `ip_address` VARCHAR(45) NULL,
    ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL,
    ADD COLUMN `user_id` CHAR(36) NULL,
    MODIFY `id` CHAR(26) NOT NULL,
    MODIFY `action` VARCHAR(100) NOT NULL,
    MODIFY `changes` JSON NULL,
    ADD PRIMARY KEY (`id`);

-- AlterTable
ALTER TABLE `auth_token` DROP PRIMARY KEY,
    DROP COLUMN `createdAt`,
    DROP COLUMN `expiresAt`,
    DROP COLUMN `revokedAt`,
    DROP COLUMN `tokenType`,
    DROP COLUMN `userId`,
    ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `expires_at` DATETIME(3) NOT NULL,
    ADD COLUMN `user_id` CHAR(36) NOT NULL,
    MODIFY `id` CHAR(26) NOT NULL,
    ADD PRIMARY KEY (`id`);

-- AlterTable
ALTER TABLE `permission` DROP COLUMN `createdAt`,
    DROP COLUMN `roleId`,
    ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `description` TEXT NULL,
    ADD COLUMN `is_system_permission` BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN `permission_name` VARCHAR(100) NOT NULL,
    ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL,
    ADD COLUMN `updated_at` DATETIME(3) NOT NULL,
    MODIFY `action` VARCHAR(50) NULL,
    MODIFY `resource` VARCHAR(100) NULL;

-- AlterTable
ALTER TABLE `quotation` DROP COLUMN `amount`,
    DROP COLUMN `createdAt`,
    DROP COLUMN `leadId`,
    DROP COLUMN `quotationNumber`,
    DROP COLUMN `tenantId`,
    DROP COLUMN `updatedAt`,
    DROP COLUMN `validUntil`,
    ADD COLUMN `approved_at` DATETIME(3) NULL,
    ADD COLUMN `approved_by` VARCHAR(36) NULL,
    ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `created_by` VARCHAR(36) NULL,
    ADD COLUMN `customer_id` VARCHAR(36) NULL,
    ADD COLUMN `net_amount` DECIMAL(18, 2) NULL,
    ADD COLUMN `quotation_date` DATE NOT NULL,
    ADD COLUMN `quotation_number` VARCHAR(50) NOT NULL,
    ADD COLUMN `tax` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL,
    ADD COLUMN `total_amount` DECIMAL(18, 2) NULL,
    ADD COLUMN `updated_at` DATETIME(3) NOT NULL,
    ADD COLUMN `valid_till` DATE NOT NULL;

-- AlterTable
ALTER TABLE `role` DROP COLUMN `createdAt`,
    DROP COLUMN `name`,
    DROP COLUMN `tenantId`,
    ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `is_active` BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN `is_system_role` BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN `role_name` VARCHAR(100) NOT NULL,
    ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL,
    ADD COLUMN `updated_at` DATETIME(3) NOT NULL;

-- AlterTable
ALTER TABLE `sales_lead` DROP COLUMN `assignedTo`,
    DROP COLUMN `createdAt`,
    DROP COLUMN `deletedAt`,
    DROP COLUMN `leadCode`,
    DROP COLUMN `leadName`,
    DROP COLUMN `leadSource`,
    DROP COLUMN `leadValue`,
    DROP COLUMN `status`,
    DROP COLUMN `tenantId`,
    DROP COLUMN `updatedAt`,
    ADD COLUMN `assigned_to` VARCHAR(36) NULL,
    ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `created_by` VARCHAR(36) NULL,
    ADD COLUMN `lead_code` VARCHAR(50) NOT NULL,
    ADD COLUMN `lead_name` VARCHAR(255) NOT NULL,
    ADD COLUMN `lead_source` VARCHAR(100) NULL,
    ADD COLUMN `lead_status` VARCHAR(50) NOT NULL DEFAULT 'new',
    ADD COLUMN `lead_value` DECIMAL(15, 2) NULL,
    ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL,
    ADD COLUMN `updated_at` DATETIME(3) NOT NULL;

-- AlterTable
ALTER TABLE `tenant` DROP COLUMN `createdAt`,
    DROP COLUMN `maxUsers`,
    DROP COLUMN `updatedAt`,
    ADD COLUMN `ai_budget_monthly` DECIMAL(15, 2) NOT NULL DEFAULT 1000.00,
    ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `max_concurrent_calls` INTEGER NOT NULL DEFAULT 50,
    ADD COLUMN `max_users` INTEGER NOT NULL DEFAULT 100,
    ADD COLUMN `updated_at` DATETIME(3) NOT NULL;

-- AlterTable
ALTER TABLE `user` DROP COLUMN `createdAt`,
    DROP COLUMN `firstName`,
    DROP COLUMN `lastLogin`,
    DROP COLUMN `lastName`,
    DROP COLUMN `passwordHash`,
    DROP COLUMN `status`,
    DROP COLUMN `tenantId`,
    DROP COLUMN `updatedAt`,
    ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `current_tenant_id` VARCHAR(36) NULL,
    ADD COLUMN `password_hash` VARCHAR(255) NOT NULL,
    ADD COLUMN `role` VARCHAR(50) NOT NULL DEFAULT 'user',
    ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL,
    ADD COLUMN `updated_at` DATETIME(3) NOT NULL;

-- AlterTable
ALTER TABLE `user_role` DROP COLUMN `assignedAt`,
    DROP COLUMN `roleId`,
    DROP COLUMN `userId`,
    ADD COLUMN `assigned_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `assigned_by` VARCHAR(36) NULL,
    ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `expires_at` DATETIME(3) NULL,
    ADD COLUMN `role_id` VARCHAR(36) NOT NULL,
    ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL,
    ADD COLUMN `user_id` CHAR(36) NOT NULL;

-- DropTable
DROP TABLE `lead_status_history`;

-- CreateTable
CREATE TABLE `password_reset_token` (
    `id` CHAR(26) NOT NULL,
    `user_id` CHAR(36) NOT NULL,
    `token` VARCHAR(255) NOT NULL,
    `expires_at` DATETIME(3) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `password_reset_token_token_key`(`token`),
    INDEX `password_reset_token_user_id_idx`(`user_id`),
    INDEX `password_reset_token_token_idx`(`token`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `team` (
    `id` CHAR(26) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `team_name` VARCHAR(255) NOT NULL,
    `description` LONGTEXT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `team_tenant_id_idx`(`tenant_id`),
    INDEX `team_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `system_config` (
    `id` CHAR(26) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `config_key` VARCHAR(255) NOT NULL,
    `config_value` LONGTEXT NULL,
    `config_type` VARCHAR(50) NOT NULL,
    `description` LONGTEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `system_config_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `system_config_tenant_id_config_key_key`(`tenant_id`, `config_key`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `call_centers` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `type` VARCHAR(50) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    `location` VARCHAR(255) NULL,
    `phone` VARCHAR(20) NULL,
    `email` VARCHAR(255) NULL,
    `timezone` VARCHAR(100) NOT NULL DEFAULT 'Asia/Kolkata',
    `operatingHoursStart` VARCHAR(10) NULL,
    `operatingHoursEnd` VARCHAR(10) NULL,
    `isHolidayOperating` BOOLEAN NOT NULL DEFAULT false,
    `maxAgents` INTEGER NOT NULL DEFAULT 50,
    `maxCampaigns` INTEGER NOT NULL DEFAULT 10,
    `maxQueues` INTEGER NOT NULL DEFAULT 20,
    `currencyCode` VARCHAR(10) NOT NULL DEFAULT 'INR',
    `costPerMinute` DECIMAL(10, 2) NOT NULL DEFAULT 0,
    `costPerCall` DECIMAL(10, 2) NOT NULL DEFAULT 0,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `call_centers_tenantId_idx`(`tenantId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `agents` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callCenterId` CHAR(36) NOT NULL,
    `employeeId` CHAR(36) NULL,
    `agentId` VARCHAR(50) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `phone` VARCHAR(20) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'OFFLINE',
    `currentCallId` CHAR(36) NULL,
    `lastStatusChangeTime` DATETIME(3) NULL,
    `lastCallEndTime` DATETIME(3) NULL,
    `assignmentType` VARCHAR(50) NOT NULL DEFAULT 'ANY',
    `stickyAgentPreference` BOOLEAN NOT NULL DEFAULT false,
    `totalCalls` INTEGER NOT NULL DEFAULT 0,
    `totalHandleTime` INTEGER NOT NULL DEFAULT 0,
    `averageHandleTime` INTEGER NOT NULL DEFAULT 0,
    `totalWrapUpTime` INTEGER NOT NULL DEFAULT 0,
    `startDate` DATETIME(3) NOT NULL,
    `endDate` DATETIME(3) NULL,
    `workDaysJson` TEXT NULL,
    `workingHoursStart` VARCHAR(10) NOT NULL DEFAULT '09:00',
    `workingHoursEnd` VARCHAR(10) NOT NULL DEFAULT '18:00',
    `breakDurationMinutes` INTEGER NOT NULL DEFAULT 15,
    `lunchDurationMinutes` INTEGER NOT NULL DEFAULT 60,
    `isTraining` BOOLEAN NOT NULL DEFAULT false,
    `trainingUntil` DATETIME(3) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    UNIQUE INDEX `agents_agentId_key`(`agentId`),
    UNIQUE INDEX `agents_currentCallId_key`(`currentCallId`),
    INDEX `agents_tenantId_idx`(`tenantId`),
    INDEX `agents_callCenterId_idx`(`callCenterId`),
    INDEX `agents_employeeId_idx`(`employeeId`),
    INDEX `agents_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ai_agents` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callCenterId` CHAR(36) NOT NULL,
    `agentId` VARCHAR(50) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `provider` VARCHAR(100) NOT NULL,
    `providerAgentId` VARCHAR(255) NOT NULL,
    `providerApiKey` TEXT NULL,
    `providerConfig` TEXT NULL,
    `language` VARCHAR(10) NOT NULL DEFAULT 'hi',
    `voiceType` VARCHAR(50) NOT NULL DEFAULT 'FEMALE',
    `voiceAccent` VARCHAR(100) NULL,
    `supportedLanguages` VARCHAR(255) NOT NULL DEFAULT 'en,hi',
    `sentiment` VARCHAR(50) NOT NULL DEFAULT 'PROFESSIONAL',
    `usesNLU` BOOLEAN NOT NULL DEFAULT true,
    `usesNLG` BOOLEAN NOT NULL DEFAULT true,
    `learningEnabled` BOOLEAN NOT NULL DEFAULT true,
    `emotionDetection` BOOLEAN NOT NULL DEFAULT true,
    `intentRecognition` BOOLEAN NOT NULL DEFAULT true,
    `canMakeOutbound` BOOLEAN NOT NULL DEFAULT true,
    `canReceiveInbound` BOOLEAN NOT NULL DEFAULT true,
    `canConferenceCall` BOOLEAN NOT NULL DEFAULT false,
    `canTransfer` BOOLEAN NOT NULL DEFAULT true,
    `handoverThreshold` DECIMAL(3, 2) NOT NULL DEFAULT 0.6,
    `handoverPreference` VARCHAR(50) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    `isAvailable` BOOLEAN NOT NULL DEFAULT true,
    `currentCallId` CHAR(36) NULL,
    `totalCalls` INTEGER NOT NULL DEFAULT 0,
    `successfulCalls` INTEGER NOT NULL DEFAULT 0,
    `transferredCalls` INTEGER NOT NULL DEFAULT 0,
    `averageConfidence` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `customerSatisfaction` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `scriptId` CHAR(36) NULL,
    `knowledgeBaseId` CHAR(36) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    UNIQUE INDEX `ai_agents_agentId_key`(`agentId`),
    UNIQUE INDEX `ai_agents_currentCallId_key`(`currentCallId`),
    INDEX `ai_agents_tenantId_idx`(`tenantId`),
    INDEX `ai_agents_callCenterId_idx`(`callCenterId`),
    INDEX `ai_agents_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `skills` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callCenterId` CHAR(36) NOT NULL,
    `skillName` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `skillCategory` VARCHAR(100) NOT NULL,
    `level` VARCHAR(50) NOT NULL DEFAULT 'BEGINNER',
    `routingPriority` INTEGER NOT NULL DEFAULT 1,
    `minimumProficiency` VARCHAR(50) NOT NULL DEFAULT 'BEGINNER',
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `skills_tenantId_idx`(`tenantId`),
    INDEX `skills_callCenterId_idx`(`callCenterId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `agent_skills` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `agentId` CHAR(36) NOT NULL,
    `skillId` CHAR(36) NOT NULL,
    `proficiency` VARCHAR(50) NOT NULL DEFAULT 'BEGINNER',
    `yearsOfExperience` INTEGER NOT NULL DEFAULT 0,
    `isVerified` BOOLEAN NOT NULL DEFAULT false,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `agent_skills_tenantId_idx`(`tenantId`),
    UNIQUE INDEX `agent_skills_agentId_skillId_key`(`agentId`, `skillId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `campaigns` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callCenterId` CHAR(36) NOT NULL,
    `campaignName` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `campaignType` VARCHAR(50) NOT NULL,
    `purpose` VARCHAR(100) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'DRAFT',
    `startDate` DATETIME(3) NULL,
    `endDate` DATETIME(3) NULL,
    `assignmentStrategy` VARCHAR(50) NOT NULL DEFAULT 'LEAST_BUSY',
    `stickyAgentEnabled` BOOLEAN NOT NULL DEFAULT false,
    `stickyAgentDuration` INTEGER NOT NULL DEFAULT 86400,
    `allowAIAgent` BOOLEAN NOT NULL DEFAULT false,
    `allowHumanAgent` BOOLEAN NOT NULL DEFAULT true,
    `aiAgentPriority` INTEGER NOT NULL DEFAULT 1,
    `routeBySource` BOOLEAN NOT NULL DEFAULT false,
    `routeBySubSource` BOOLEAN NOT NULL DEFAULT false,
    `routeByLeadSource` BOOLEAN NOT NULL DEFAULT true,
    `maxWaitTime` INTEGER NOT NULL DEFAULT 900,
    `abandonedCallHandling` VARCHAR(50) NOT NULL,
    `salesProjectId` CHAR(36) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `campaigns_tenantId_idx`(`tenantId`),
    INDEX `campaigns_callCenterId_idx`(`callCenterId`),
    INDEX `campaigns_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ai_campaigns` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `campaignId` CHAR(36) NOT NULL,
    `aiAgentId` CHAR(36) NOT NULL,
    `assignmentPercentage` INTEGER NOT NULL DEFAULT 50,
    `priority` INTEGER NOT NULL DEFAULT 1,
    `isActive` BOOLEAN NOT NULL DEFAULT true,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `ai_campaigns_tenantId_idx`(`tenantId`),
    UNIQUE INDEX `ai_campaigns_campaignId_aiAgentId_key`(`campaignId`, `aiAgentId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `queues` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callCenterId` CHAR(36) NOT NULL,
    `queueName` VARCHAR(255) NOT NULL,
    `queueNumber` VARCHAR(50) NOT NULL,
    `description` TEXT NULL,
    `maxCapacity` INTEGER NOT NULL DEFAULT 100,
    `currentSize` INTEGER NOT NULL DEFAULT 0,
    `routingStrategy` VARCHAR(50) NOT NULL DEFAULT 'FIFO',
    `priorityQueueEnabled` BOOLEAN NOT NULL DEFAULT false,
    `timeoutSeconds` INTEGER NOT NULL DEFAULT 300,
    `abandonmentThreshold` INTEGER NOT NULL DEFAULT 600,
    `welcomeMessageId` CHAR(36) NULL,
    `hasWelcomeMessage` BOOLEAN NOT NULL DEFAULT false,
    `holdMusicUrl` TEXT NULL,
    `hasHoldMusic` BOOLEAN NOT NULL DEFAULT false,
    `status` VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `queues_tenantId_idx`(`tenantId`),
    INDEX `queues_callCenterId_idx`(`callCenterId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `queue_skill_requirements` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `queueId` CHAR(36) NOT NULL,
    `skillId` CHAR(36) NOT NULL,
    `isRequired` BOOLEAN NOT NULL DEFAULT false,
    `minimumProficiency` VARCHAR(50) NOT NULL DEFAULT 'BEGINNER',
    `weight` INTEGER NOT NULL DEFAULT 1,

    INDEX `queue_skill_requirements_tenantId_idx`(`tenantId`),
    UNIQUE INDEX `queue_skill_requirements_queueId_skillId_key`(`queueId`, `skillId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `queue_assignments` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `agentId` CHAR(36) NOT NULL,
    `queueId` CHAR(36) NOT NULL,
    `assignmentType` VARCHAR(50) NOT NULL DEFAULT 'PRIMARY',
    `priority` INTEGER NOT NULL DEFAULT 1,
    `assignmentStartTime` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `assignmentEndTime` DATETIME(3) NULL,
    `isActive` BOOLEAN NOT NULL DEFAULT true,

    INDEX `queue_assignments_tenantId_idx`(`tenantId`),
    UNIQUE INDEX `queue_assignments_agentId_queueId_assignmentStartTime_key`(`agentId`, `queueId`, `assignmentStartTime`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `agent_campaigns` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `agentId` CHAR(36) NOT NULL,
    `campaignId` CHAR(36) NOT NULL,
    `assignmentType` VARCHAR(50) NOT NULL DEFAULT 'PRIMARY',
    `priority` INTEGER NOT NULL DEFAULT 1,
    `assignmentStartTime` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `assignmentEndTime` DATETIME(3) NULL,
    `isActive` BOOLEAN NOT NULL DEFAULT true,

    INDEX `agent_campaigns_tenantId_idx`(`tenantId`),
    UNIQUE INDEX `agent_campaigns_agentId_campaignId_assignmentStartTime_key`(`agentId`, `campaignId`, `assignmentStartTime`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `calls` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callCenterId` CHAR(36) NOT NULL,
    `callId` VARCHAR(50) NOT NULL,
    `callULID` CHAR(26) NOT NULL,
    `campaignId` CHAR(36) NULL,
    `queueId` CHAR(36) NULL,
    `callDirection` VARCHAR(20) NOT NULL,
    `callerId` VARCHAR(20) NOT NULL,
    `callerName` VARCHAR(255) NULL,
    `callerEmail` VARCHAR(255) NULL,
    `recipientId` VARCHAR(20) NULL,
    `recipientName` VARCHAR(255) NULL,
    `leadId` CHAR(36) NULL,
    `customerId` CHAR(36) NULL,
    `agentId` CHAR(36) NULL,
    `aiAgentId` CHAR(36) NULL,
    `callType` VARCHAR(50) NOT NULL,
    `source` VARCHAR(100) NOT NULL,
    `subSource` VARCHAR(100) NULL,
    `campaignSource` VARCHAR(100) NULL,
    `callStatus` VARCHAR(50) NOT NULL DEFAULT 'INITIATED',
    `callInitiatedAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `callConnectedAt` DATETIME(3) NULL,
    `callEndedAt` DATETIME(3) NULL,
    `callDurationSeconds` INTEGER NOT NULL DEFAULT 0,
    `dispositionId` CHAR(36) NULL,
    `notes` TEXT NULL,
    `transferredFromAgentId` CHAR(36) NULL,
    `transferredToAgentId` CHAR(36) NULL,
    `transferredToAIAgentId` CHAR(36) NULL,
    `transferReason` VARCHAR(255) NULL,
    `transferredAt` DATETIME(3) NULL,
    `aiHandlingDuration` INTEGER NULL,
    `aiConfidenceScore` DECIMAL(5, 2) NULL,
    `wasHandledByAI` BOOLEAN NOT NULL DEFAULT false,
    `wasTransferredFromAI` BOOLEAN NOT NULL DEFAULT false,
    `aiTranscript` TEXT NULL,
    `recordingUrl` TEXT NULL,
    `recordingId` VARCHAR(255) NULL,
    `recordingStatus` VARCHAR(50) NULL,
    `isCallRecorded` BOOLEAN NOT NULL DEFAULT false,
    `satScore` INTEGER NULL,
    `feedback` TEXT NULL,
    `isSensitive` BOOLEAN NOT NULL DEFAULT false,
    `requiresAudit` BOOLEAN NOT NULL DEFAULT false,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    UNIQUE INDEX `calls_callId_key`(`callId`),
    UNIQUE INDEX `calls_callULID_key`(`callULID`),
    INDEX `calls_tenantId_idx`(`tenantId`),
    INDEX `calls_callCenterId_idx`(`callCenterId`),
    INDEX `calls_campaignId_idx`(`campaignId`),
    INDEX `calls_queueId_idx`(`queueId`),
    INDEX `calls_agentId_idx`(`agentId`),
    INDEX `calls_aiAgentId_idx`(`aiAgentId`),
    INDEX `calls_leadId_idx`(`leadId`),
    INDEX `calls_customerId_idx`(`customerId`),
    INDEX `calls_callStatus_idx`(`callStatus`),
    INDEX `calls_callInitiatedAt_idx`(`callInitiatedAt`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `call_dispositions` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callCenterId` CHAR(36) NOT NULL,
    `dispositionCode` VARCHAR(50) NOT NULL,
    `dispositionName` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `dispositionCategory` VARCHAR(100) NOT NULL,
    `isPositive` BOOLEAN NOT NULL DEFAULT false,
    `allowsNotes` BOOLEAN NOT NULL DEFAULT true,
    `requiresFollowup` BOOLEAN NOT NULL DEFAULT false,
    `requiresCallback` BOOLEAN NOT NULL DEFAULT false,
    `allowsWrapupTime` BOOLEAN NOT NULL DEFAULT true,
    `sequenceOrder` INTEGER NOT NULL DEFAULT 1,
    `isActive` BOOLEAN NOT NULL DEFAULT true,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `call_dispositions_tenantId_idx`(`tenantId`),
    INDEX `call_dispositions_callCenterId_idx`(`callCenterId`),
    UNIQUE INDEX `call_dispositions_callCenterId_dispositionCode_key`(`callCenterId`, `dispositionCode`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `call_recordings` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callCenterId` CHAR(36) NOT NULL,
    `callId` CHAR(36) NOT NULL,
    `recordingUrl` TEXT NOT NULL,
    `recordingExternalId` VARCHAR(255) NOT NULL,
    `fileSize` BIGINT NULL,
    `fileDuration` INTEGER NULL,
    `fileFormat` VARCHAR(20) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'PROCESSING',
    `hasTranscript` BOOLEAN NOT NULL DEFAULT false,
    `transcriptUrl` TEXT NULL,
    `transcriptText` LONGTEXT NULL,
    `transcriptionAccuracy` DECIMAL(5, 2) NULL,
    `hasSentimentAnalysis` BOOLEAN NOT NULL DEFAULT false,
    `overallSentiment` VARCHAR(50) NULL,
    `sentimentScore` DECIMAL(5, 2) NULL,
    `isCompliant` BOOLEAN NULL,
    `complianceNotes` TEXT NULL,
    `retentionDays` INTEGER NULL,
    `deleteAfterDate` DATETIME(3) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `call_recordings_tenantId_idx`(`tenantId`),
    INDEX `call_recordings_callCenterId_idx`(`callCenterId`),
    INDEX `call_recordings_callId_idx`(`callId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ai_agent_dialogues` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callId` CHAR(36) NOT NULL,
    `aiAgentId` CHAR(36) NOT NULL,
    `turnNumber` INTEGER NOT NULL,
    `speaker` VARCHAR(20) NOT NULL,
    `rawText` TEXT NOT NULL,
    `cleanedText` TEXT NOT NULL,
    `intent` VARCHAR(255) NULL,
    `entities` TEXT NULL,
    `sentiment` VARCHAR(50) NULL,
    `sentimentScore` DECIMAL(5, 2) NULL,
    `confidence` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `generatedResponse` TEXT NULL,
    `responseType` VARCHAR(100) NULL,
    `turnStartTime` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `turnDurationSeconds` INTEGER NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `ai_agent_dialogues_tenantId_idx`(`tenantId`),
    INDEX `ai_agent_dialogues_callId_idx`(`callId`),
    INDEX `ai_agent_dialogues_aiAgentId_idx`(`aiAgentId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ai_conversation_logs` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callId` CHAR(36) NOT NULL,
    `aiAgentId` CHAR(36) NOT NULL,
    `totalTurns` INTEGER NOT NULL,
    `averageConfidence` DECIMAL(5, 2) NOT NULL,
    `dominantIntent` VARCHAR(255) NULL,
    `multipleIntents` BOOLEAN NOT NULL DEFAULT false,
    `overallSentiment` VARCHAR(50) NULL,
    `sentimentTrend` VARCHAR(20) NULL,
    `handoverReason` VARCHAR(255) NULL,
    `handoverConfidence` DECIMAL(5, 2) NULL,
    `flaggedForReview` BOOLEAN NOT NULL DEFAULT false,
    `reviewNotes` TEXT NULL,
    `trainingDataUsed` BOOLEAN NOT NULL DEFAULT true,
    `callerSatisfaction` DECIMAL(5, 2) NULL,
    `issueResolved` BOOLEAN NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `ai_conversation_logs_tenantId_idx`(`tenantId`),
    INDEX `ai_conversation_logs_callId_idx`(`callId`),
    INDEX `ai_conversation_logs_aiAgentId_idx`(`aiAgentId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `call_scripts` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callCenterId` CHAR(36) NOT NULL,
    `scriptName` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `scriptType` VARCHAR(50) NOT NULL,
    `purpose` VARCHAR(100) NOT NULL,
    `scriptContent` LONGTEXT NOT NULL,
    `scriptVariations` TEXT NULL,
    `isAIOptimized` BOOLEAN NOT NULL DEFAULT false,
    `aiTrainingScore` DECIMAL(5, 2) NULL,
    `complianceChecked` BOOLEAN NOT NULL DEFAULT false,
    `complianceNotes` TEXT NULL,
    `isActive` BOOLEAN NOT NULL DEFAULT true,
    `version` INTEGER NOT NULL DEFAULT 1,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `call_scripts_tenantId_idx`(`tenantId`),
    INDEX `call_scripts_callCenterId_idx`(`callCenterId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `external_providers` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `callCenterId` CHAR(36) NOT NULL,
    `providerName` VARCHAR(100) NOT NULL,
    `providerType` VARCHAR(50) NOT NULL,
    `apiEndpoint` TEXT NOT NULL,
    `apiKey` TEXT NOT NULL,
    `apiSecret` TEXT NULL,
    `accountId` VARCHAR(255) NULL,
    `country` VARCHAR(10) NOT NULL DEFAULT 'IN',
    `supportedCountries` TEXT NOT NULL,
    `supportsInbound` BOOLEAN NOT NULL DEFAULT true,
    `supportsOutbound` BOOLEAN NOT NULL DEFAULT true,
    `supportsConference` BOOLEAN NOT NULL DEFAULT false,
    `supportsRecording` BOOLEAN NOT NULL DEFAULT true,
    `supportsIVR` BOOLEAN NOT NULL DEFAULT true,
    `supportsTransfer` BOOLEAN NOT NULL DEFAULT true,
    `allocatedNumbers` TEXT NOT NULL,
    `numberFormat` VARCHAR(20) NOT NULL DEFAULT 'E.164',
    `costModel` VARCHAR(50) NOT NULL,
    `costPerMinute` DECIMAL(10, 4) NULL,
    `costPerCall` DECIMAL(10, 2) NULL,
    `setupCost` DECIMAL(10, 2) NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    `lastHealthCheck` DATETIME(3) NULL,
    `isHealthy` BOOLEAN NOT NULL DEFAULT true,
    `webhookUrl` TEXT NULL,
    `webhookSecret` TEXT NULL,
    `notes` TEXT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `external_providers_tenantId_idx`(`tenantId`),
    INDEX `external_providers_callCenterId_idx`(`callCenterId`),
    INDEX `external_providers_status_idx`(`status`),
    UNIQUE INDEX `external_providers_tenantId_callCenterId_providerName_key`(`tenantId`, `callCenterId`, `providerName`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `provider_connections` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `providerId` CHAR(36) NOT NULL,
    `callCenterId` CHAR(36) NOT NULL,
    `isPrimary` BOOLEAN NOT NULL DEFAULT false,
    `priority` INTEGER NOT NULL DEFAULT 1,
    `allocationPercentage` INTEGER NOT NULL DEFAULT 100,
    `isActive` BOOLEAN NOT NULL DEFAULT true,
    `connectionStatus` VARCHAR(50) NOT NULL DEFAULT 'CONNECTED',
    `connectedAt` DATETIME(3) NULL,
    `lastActivityAt` DATETIME(3) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `provider_connections_tenantId_idx`(`tenantId`),
    UNIQUE INDEX `provider_connections_providerId_callCenterId_key`(`providerId`, `callCenterId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `agent_notes` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `agentId` CHAR(36) NOT NULL,
    `noteType` VARCHAR(50) NOT NULL,
    `noteTitle` VARCHAR(255) NOT NULL,
    `noteContent` TEXT NOT NULL,
    `severity` VARCHAR(20) NULL,
    `createdBy` VARCHAR(255) NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `agent_notes_tenantId_idx`(`tenantId`),
    INDEX `agent_notes_agentId_idx`(`agentId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `agent_performance_metrics` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `agentId` CHAR(36) NOT NULL,
    `metricDate` DATETIME(3) NOT NULL,
    `metricPeriod` VARCHAR(20) NOT NULL DEFAULT 'DAILY',
    `totalCalls` INTEGER NOT NULL DEFAULT 0,
    `inboundCalls` INTEGER NOT NULL DEFAULT 0,
    `outboundCalls` INTEGER NOT NULL DEFAULT 0,
    `answeredCalls` INTEGER NOT NULL DEFAULT 0,
    `missedCalls` INTEGER NOT NULL DEFAULT 0,
    `totalTalkTime` INTEGER NOT NULL DEFAULT 0,
    `totalHoldTime` INTEGER NOT NULL DEFAULT 0,
    `totalWrapUpTime` INTEGER NOT NULL DEFAULT 0,
    `totalBreakTime` INTEGER NOT NULL DEFAULT 0,
    `answerRate` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `averageHandleTime` INTEGER NOT NULL DEFAULT 0,
    `averageHoldTime` DECIMAL(8, 2) NOT NULL DEFAULT 0,
    `averageCSAT` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `transferRate` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `salesConversions` INTEGER NOT NULL DEFAULT 0,
    `conversionRate` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `totalSalesValue` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `agent_performance_metrics_tenantId_idx`(`tenantId`),
    INDEX `agent_performance_metrics_agentId_idx`(`agentId`),
    INDEX `agent_performance_metrics_metricDate_idx`(`metricDate`),
    UNIQUE INDEX `agent_performance_metrics_agentId_metricDate_metricPeriod_key`(`agentId`, `metricDate`, `metricPeriod`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `attribution_event` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `event_id` CHAR(26) NOT NULL,
    `lead_id` VARCHAR(36) NULL,
    `visitor_id` VARCHAR(36) NULL,
    `session_id` VARCHAR(100) NULL,
    `attribution_bucket` ENUM('SEO', 'SEM', 'SMM', 'PORTALS', 'OWNED', 'OFFLINE', 'AFFILIATE', 'DIRECT') NOT NULL,
    `attribution_source` ENUM('google', 'bing', 'yahoo', 'other_search', 'meta', 'instagram', 'facebook', 'youtube', 'linkedin', 'twitter', 'portal_99acres', 'portal_magicbricks', 'portal_housing', 'portal_nobroker', 'portal_squareyards', 'website', 'whatsapp', 'email', 'sms', 'crm', 'walk_in', 'print_media', 'btl_activation', 'broker', 'referral', 'consultant', 'direct') NOT NULL,
    `attribution_subsource` ENUM('google_organic', 'bing_organic', 'yahoo_organic', 'discover', 'local_pack', 'google_search_ads', 'google_display_ads', 'youtube_ads', 'performance_max', 'bing_ads', 'facebook_ads', 'instagram_ads', 'linkedin_ads', 'twitter_ads', 'instagram_organic', 'facebook_organic', 'youtube_organic', 'linkedin_organic', 'portal_listing', 'web_form', 'whatsapp_chat', 'email_campaign', 'sms_campaign', 'crm_remarketing', 'walk_in_site', 'walk_in_office', 'newspaper', 'magazine', 'flyers', 'brochure', 'mall_stall', 'roadshow_stall', 'expo_stall', 'society_activation', 'channel_partner', 'referral_agent', 'property_consultant', 'direct_unknown', 'bookmark', 'typed_url') NOT NULL,
    `event_type` VARCHAR(100) NOT NULL,
    `event_value` DECIMAL(15, 2) NULL,
    `event_data` LONGTEXT NULL,
    `campaign` VARCHAR(255) NULL,
    `ad_set` VARCHAR(255) NULL,
    `ad` VARCHAR(255) NULL,
    `keyword` VARCHAR(255) NULL,
    `referrer` TEXT NULL,
    `user_agent` TEXT NULL,
    `ip_address` VARCHAR(45) NULL,
    `country` VARCHAR(100) NULL,
    `state` VARCHAR(100) NULL,
    `city` VARCHAR(100) NULL,
    `source` VARCHAR(100) NULL,
    `is_test` BOOLEAN NOT NULL DEFAULT false,
    `event_timestamp` DATETIME NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `attribution_event_tenant_id_attribution_bucket_idx`(`tenant_id`, `attribution_bucket`),
    INDEX `attribution_event_tenant_id_attribution_source_idx`(`tenant_id`, `attribution_source`),
    INDEX `attribution_event_tenant_id_attribution_subsource_idx`(`tenant_id`, `attribution_subsource`),
    INDEX `attribution_event_tenant_id_event_timestamp_idx`(`tenant_id`, `event_timestamp`),
    INDEX `attribution_event_tenant_id_lead_id_idx`(`tenant_id`, `lead_id`),
    INDEX `attribution_event_tenant_id_visitor_id_idx`(`tenant_id`, `visitor_id`),
    INDEX `attribution_event_tenant_id_campaign_idx`(`tenant_id`, `campaign`),
    INDEX `attribution_event_lead_id_idx`(`lead_id`),
    INDEX `attribution_event_visitor_id_idx`(`visitor_id`),
    UNIQUE INDEX `attribution_event_tenant_id_event_id_key`(`tenant_id`, `event_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `attribution_bucket_config` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `bucket` ENUM('SEO', 'SEM', 'SMM', 'PORTALS', 'OWNED', 'OFFLINE', 'AFFILIATE', 'DIRECT') NOT NULL,
    `bucket_label` VARCHAR(100) NOT NULL,
    `bucket_description` TEXT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `priority` INTEGER NOT NULL DEFAULT 0,
    `weight` DECIMAL(5, 2) NOT NULL DEFAULT 1.0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `attribution_bucket_config_tenant_id_is_active_idx`(`tenant_id`, `is_active`),
    UNIQUE INDEX `attribution_bucket_config_tenant_id_bucket_key`(`tenant_id`, `bucket`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `attribution_source_config` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `source` ENUM('google', 'bing', 'yahoo', 'other_search', 'meta', 'instagram', 'facebook', 'youtube', 'linkedin', 'twitter', 'portal_99acres', 'portal_magicbricks', 'portal_housing', 'portal_nobroker', 'portal_squareyards', 'website', 'whatsapp', 'email', 'sms', 'crm', 'walk_in', 'print_media', 'btl_activation', 'broker', 'referral', 'consultant', 'direct') NOT NULL,
    `source_label` VARCHAR(100) NOT NULL,
    `source_description` TEXT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `requires_authorization` BOOLEAN NOT NULL DEFAULT false,
    `api_key_required` BOOLEAN NOT NULL DEFAULT false,
    `integration_url` VARCHAR(500) NULL,
    `last_sync_date` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `attribution_source_config_tenant_id_is_active_idx`(`tenant_id`, `is_active`),
    UNIQUE INDEX `attribution_source_config_tenant_id_source_key`(`tenant_id`, `source`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `lead_attribution_snapshot` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `lead_id` VARCHAR(36) NOT NULL,
    `first_touch_bucket` ENUM('SEO', 'SEM', 'SMM', 'PORTALS', 'OWNED', 'OFFLINE', 'AFFILIATE', 'DIRECT') NULL,
    `first_touch_source` ENUM('google', 'bing', 'yahoo', 'other_search', 'meta', 'instagram', 'facebook', 'youtube', 'linkedin', 'twitter', 'portal_99acres', 'portal_magicbricks', 'portal_housing', 'portal_nobroker', 'portal_squareyards', 'website', 'whatsapp', 'email', 'sms', 'crm', 'walk_in', 'print_media', 'btl_activation', 'broker', 'referral', 'consultant', 'direct') NULL,
    `first_touch_subsource` ENUM('google_organic', 'bing_organic', 'yahoo_organic', 'discover', 'local_pack', 'google_search_ads', 'google_display_ads', 'youtube_ads', 'performance_max', 'bing_ads', 'facebook_ads', 'instagram_ads', 'linkedin_ads', 'twitter_ads', 'instagram_organic', 'facebook_organic', 'youtube_organic', 'linkedin_organic', 'portal_listing', 'web_form', 'whatsapp_chat', 'email_campaign', 'sms_campaign', 'crm_remarketing', 'walk_in_site', 'walk_in_office', 'newspaper', 'magazine', 'flyers', 'brochure', 'mall_stall', 'roadshow_stall', 'expo_stall', 'society_activation', 'channel_partner', 'referral_agent', 'property_consultant', 'direct_unknown', 'bookmark', 'typed_url') NULL,
    `first_touch_campaign` VARCHAR(255) NULL,
    `first_touch_date` DATETIME(3) NULL,
    `last_touch_bucket` ENUM('SEO', 'SEM', 'SMM', 'PORTALS', 'OWNED', 'OFFLINE', 'AFFILIATE', 'DIRECT') NULL,
    `last_touch_source` ENUM('google', 'bing', 'yahoo', 'other_search', 'meta', 'instagram', 'facebook', 'youtube', 'linkedin', 'twitter', 'portal_99acres', 'portal_magicbricks', 'portal_housing', 'portal_nobroker', 'portal_squareyards', 'website', 'whatsapp', 'email', 'sms', 'crm', 'walk_in', 'print_media', 'btl_activation', 'broker', 'referral', 'consultant', 'direct') NULL,
    `last_touch_subsource` ENUM('google_organic', 'bing_organic', 'yahoo_organic', 'discover', 'local_pack', 'google_search_ads', 'google_display_ads', 'youtube_ads', 'performance_max', 'bing_ads', 'facebook_ads', 'instagram_ads', 'linkedin_ads', 'twitter_ads', 'instagram_organic', 'facebook_organic', 'youtube_organic', 'linkedin_organic', 'portal_listing', 'web_form', 'whatsapp_chat', 'email_campaign', 'sms_campaign', 'crm_remarketing', 'walk_in_site', 'walk_in_office', 'newspaper', 'magazine', 'flyers', 'brochure', 'mall_stall', 'roadshow_stall', 'expo_stall', 'society_activation', 'channel_partner', 'referral_agent', 'property_consultant', 'direct_unknown', 'bookmark', 'typed_url') NULL,
    `last_touch_campaign` VARCHAR(255) NULL,
    `last_touch_date` DATETIME(3) NULL,
    `total_touches` INTEGER NOT NULL DEFAULT 1,
    `touches_by_bucket` LONGTEXT NULL,
    `touches_by_source` LONGTEXT NULL,
    `touch_sequence` LONGTEXT NULL,
    `conversion_bucket` ENUM('SEO', 'SEM', 'SMM', 'PORTALS', 'OWNED', 'OFFLINE', 'AFFILIATE', 'DIRECT') NULL,
    `conversion_source` ENUM('google', 'bing', 'yahoo', 'other_search', 'meta', 'instagram', 'facebook', 'youtube', 'linkedin', 'twitter', 'portal_99acres', 'portal_magicbricks', 'portal_housing', 'portal_nobroker', 'portal_squareyards', 'website', 'whatsapp', 'email', 'sms', 'crm', 'walk_in', 'print_media', 'btl_activation', 'broker', 'referral', 'consultant', 'direct') NULL,
    `conversion_subsource` ENUM('google_organic', 'bing_organic', 'yahoo_organic', 'discover', 'local_pack', 'google_search_ads', 'google_display_ads', 'youtube_ads', 'performance_max', 'bing_ads', 'facebook_ads', 'instagram_ads', 'linkedin_ads', 'twitter_ads', 'instagram_organic', 'facebook_organic', 'youtube_organic', 'linkedin_organic', 'portal_listing', 'web_form', 'whatsapp_chat', 'email_campaign', 'sms_campaign', 'crm_remarketing', 'walk_in_site', 'walk_in_office', 'newspaper', 'magazine', 'flyers', 'brochure', 'mall_stall', 'roadshow_stall', 'expo_stall', 'society_activation', 'channel_partner', 'referral_agent', 'property_consultant', 'direct_unknown', 'bookmark', 'typed_url') NULL,
    `conversion_date` DATETIME(3) NULL,
    `snapshot_date` DATE NOT NULL,
    `is_converted` BOOLEAN NOT NULL DEFAULT false,
    `conversion_value` DECIMAL(15, 2) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `lead_attribution_snapshot_tenant_id_lead_id_idx`(`tenant_id`, `lead_id`),
    INDEX `lead_attribution_snapshot_tenant_id_first_touch_bucket_idx`(`tenant_id`, `first_touch_bucket`),
    INDEX `lead_attribution_snapshot_tenant_id_last_touch_bucket_idx`(`tenant_id`, `last_touch_bucket`),
    INDEX `lead_attribution_snapshot_tenant_id_conversion_bucket_idx`(`tenant_id`, `conversion_bucket`),
    INDEX `lead_attribution_snapshot_tenant_id_is_converted_idx`(`tenant_id`, `is_converted`),
    INDEX `lead_attribution_snapshot_snapshot_date_idx`(`snapshot_date`),
    UNIQUE INDEX `lead_attribution_snapshot_tenant_id_lead_id_snapshot_date_key`(`tenant_id`, `lead_id`, `snapshot_date`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `attribution_rule` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `rule_name` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `rule_type` VARCHAR(50) NOT NULL,
    `conditions` LONGTEXT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `priority` INTEGER NOT NULL DEFAULT 0,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `attribution_rule_tenant_id_is_active_idx`(`tenant_id`, `is_active`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `attribution_exclusion` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `exclusion_type` VARCHAR(100) NOT NULL,
    `exclusion_name` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `match_type` VARCHAR(50) NOT NULL,
    `match_value` VARCHAR(500) NOT NULL,
    `field` VARCHAR(100) NOT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `applied_count` INTEGER NOT NULL DEFAULT 0,
    `last_applied_date` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `attribution_exclusion_tenant_id_field_idx`(`tenant_id`, `field`),
    INDEX `attribution_exclusion_tenant_id_is_active_idx`(`tenant_id`, `is_active`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `attribution_metric` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `metric_date` DATE NOT NULL,
    `bucket` ENUM('SEO', 'SEM', 'SMM', 'PORTALS', 'OWNED', 'OFFLINE', 'AFFILIATE', 'DIRECT') NOT NULL,
    `total_events` INTEGER NOT NULL,
    `unique_leads` INTEGER NOT NULL,
    `unique_visitors` INTEGER NOT NULL,
    `conversions` INTEGER NOT NULL,
    `conversion_rate` DECIMAL(5, 4) NOT NULL,
    `avg_touches_per_lead` DECIMAL(5, 2) NOT NULL,
    `total_spend` DECIMAL(15, 2) NULL,
    `total_revenue` DECIMAL(15, 2) NULL,
    `roi` DECIMAL(8, 2) NULL,
    `data_quality` VARCHAR(50) NOT NULL,
    `last_updated` DATETIME NOT NULL,

    INDEX `attribution_metric_tenant_id_metric_date_idx`(`tenant_id`, `metric_date`),
    INDEX `attribution_metric_tenant_id_bucket_idx`(`tenant_id`, `bucket`),
    UNIQUE INDEX `attribution_metric_tenant_id_metric_date_bucket_key`(`tenant_id`, `metric_date`, `bucket`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `attribution_audit_log` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `action` VARCHAR(100) NOT NULL,
    `entity_type` VARCHAR(100) NOT NULL,
    `entity_id` VARCHAR(36) NULL,
    `old_value` LONGTEXT NULL,
    `new_value` LONGTEXT NULL,
    `user_id` VARCHAR(36) NULL,
    `ip_address` VARCHAR(45) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `attribution_audit_log_tenant_id_action_idx`(`tenant_id`, `action`),
    INDEX `attribution_audit_log_tenant_id_entity_type_idx`(`tenant_id`, `entity_type`),
    INDEX `attribution_audit_log_tenant_id_created_at_idx`(`tenant_id`, `created_at`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `tenant_user_count` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` CHAR(36) NOT NULL,
    `total_active_users` INTEGER NOT NULL DEFAULT 0,
    `total_inactive_users` INTEGER NOT NULL DEFAULT 0,
    `total_suspended_users` INTEGER NOT NULL DEFAULT 0,
    `total_users` INTEGER NOT NULL DEFAULT 0,
    `max_users_allowed` INTEGER NOT NULL DEFAULT 100,
    `is_over_limit` BOOLEAN NOT NULL DEFAULT false,
    `overage_count` INTEGER NOT NULL DEFAULT 0,
    `seat_utilization` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `seats_remaining` INTEGER NOT NULL DEFAULT 0,
    `admin_users` INTEGER NOT NULL DEFAULT 0,
    `manager_users` INTEGER NOT NULL DEFAULT 0,
    `standard_users` INTEGER NOT NULL DEFAULT 0,
    `guest_users` INTEGER NOT NULL DEFAULT 0,
    `last_counted_at` DATETIME(3) NOT NULL,
    `last_overage_alert_sent` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `tenant_user_count_tenant_id_key`(`tenant_id`),
    INDEX `tenant_user_count_tenant_id_idx`(`tenant_id`),
    INDEX `tenant_user_count_is_over_limit_idx`(`is_over_limit`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `tenant_user_count_history` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` CHAR(36) NOT NULL,
    `user_count_id` CHAR(36) NOT NULL,
    `snapshot_date` DATE NOT NULL,
    `active_user_count` INTEGER NOT NULL,
    `inactive_user_count` INTEGER NOT NULL,
    `suspended_user_count` INTEGER NOT NULL,
    `total_user_count` INTEGER NOT NULL,
    `billable_user_count` INTEGER NOT NULL,
    `peak_concurrent_users` INTEGER NULL,
    `new_users_added` INTEGER NOT NULL DEFAULT 0,
    `users_removed` INTEGER NOT NULL DEFAULT 0,
    `users_reactivated` INTEGER NOT NULL DEFAULT 0,
    `users_suspended` INTEGER NOT NULL DEFAULT 0,
    `current_plan` VARCHAR(100) NULL,
    `max_allowed_users` INTEGER NOT NULL,
    `is_over_limit` BOOLEAN NOT NULL DEFAULT false,
    `overage_charge_applied` BOOLEAN NOT NULL DEFAULT false,
    `recorded_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `recorded_by` VARCHAR(100) NULL,
    `notes` TEXT NULL,

    INDEX `tenant_user_count_history_tenant_id_idx`(`tenant_id`),
    INDEX `tenant_user_count_history_snapshot_date_idx`(`snapshot_date`),
    UNIQUE INDEX `tenant_user_count_history_user_count_id_snapshot_date_key`(`user_count_id`, `snapshot_date`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `user_activity_log` (
    `id` CHAR(36) NOT NULL,
    `user_id` CHAR(36) NOT NULL,
    `tenant_id` CHAR(36) NOT NULL,
    `activity_type` VARCHAR(50) NOT NULL,
    `activity_timestamp` DATETIME(3) NOT NULL,
    `session_duration_seconds` INTEGER NULL,
    `session_id` VARCHAR(255) NULL,
    `ip_address` VARCHAR(45) NULL,
    `user_agent` TEXT NULL,
    `device_type` VARCHAR(50) NULL,
    `country` VARCHAR(100) NULL,
    `state` VARCHAR(100) NULL,
    `city` VARCHAR(100) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `user_activity_log_user_id_idx`(`user_id`),
    INDEX `user_activity_log_tenant_id_idx`(`tenant_id`),
    INDEX `user_activity_log_activity_timestamp_idx`(`activity_timestamp`),
    INDEX `user_activity_log_activity_type_idx`(`activity_type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `vendor` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `vendor_code` VARCHAR(50) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NULL,
    `phone` VARCHAR(20) NULL,
    `address` TEXT NULL,
    `city` VARCHAR(100) NULL,
    `state` VARCHAR(100) NULL,
    `country` VARCHAR(100) NULL,
    `postal_code` VARCHAR(20) NULL,
    `tax_id` VARCHAR(50) NULL,
    `payment_terms` VARCHAR(100) NULL,
    `vendor_type` VARCHAR(50) NULL,
    `rating` DECIMAL(3, 1) NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `is_blocked` BOOLEAN NOT NULL DEFAULT false,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_by` VARCHAR(36) NULL,
    `gl_vendor_payable_account_id` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `vendor_vendor_code_key`(`vendor_code`),
    INDEX `vendor_tenant_id_idx`(`tenant_id`),
    INDEX `vendor_vendor_code_idx`(`vendor_code`),
    INDEX `vendor_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `vendor_contact` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `vendor_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `title` VARCHAR(100) NULL,
    `phone` VARCHAR(20) NULL,
    `email` VARCHAR(255) NULL,
    `is_primary` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `vendor_contact_tenant_id_idx`(`tenant_id`),
    INDEX `vendor_contact_vendor_id_idx`(`vendor_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `vendor_address` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `vendor_id` VARCHAR(36) NOT NULL,
    `type` VARCHAR(50) NULL,
    `line1` VARCHAR(255) NULL,
    `line2` VARCHAR(255) NULL,
    `city` VARCHAR(100) NULL,
    `state` VARCHAR(100) NULL,
    `country` VARCHAR(100) NULL,
    `post_code` VARCHAR(20) NULL,
    `is_primary` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `vendor_address_tenant_id_idx`(`tenant_id`),
    INDEX `vendor_address_vendor_id_idx`(`vendor_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `purchase_requisition` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `requisition_number` VARCHAR(50) NOT NULL,
    `requester_id` VARCHAR(36) NULL,
    `department` VARCHAR(100) NULL,
    `request_date` DATE NOT NULL,
    `required_by_date` DATE NULL,
    `purpose` TEXT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `approved_by` VARCHAR(36) NULL,
    `approved_at` DATETIME(3) NULL,
    `rejection_reason` TEXT NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `purchase_requisition_requisition_number_key`(`requisition_number`),
    INDEX `purchase_requisition_tenant_id_idx`(`tenant_id`),
    INDEX `purchase_requisition_requisition_number_idx`(`requisition_number`),
    INDEX `purchase_requisition_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `purchase_order` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `po_number` VARCHAR(50) NOT NULL,
    `vendor_id` VARCHAR(36) NOT NULL,
    `requisition_id` VARCHAR(36) NULL,
    `po_date` DATE NOT NULL,
    `delivery_date` DATE NULL,
    `total_amount` DECIMAL(18, 2) NULL,
    `tax_amount` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `shipping_amount` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `discount_amount` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `net_amount` DECIMAL(18, 2) NULL,
    `payment_terms` VARCHAR(100) NULL,
    `delivery_location` VARCHAR(255) NULL,
    `special_instructions` TEXT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `sent_to_vendor_at` DATETIME(3) NULL,
    `acknowledged_at` DATETIME(3) NULL,
    `created_by` VARCHAR(36) NULL,
    `gl_inventory_account_id` VARCHAR(36) NULL,
    `gl_expense_account_id` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `purchase_order_po_number_key`(`po_number`),
    INDEX `purchase_order_tenant_id_idx`(`tenant_id`),
    INDEX `purchase_order_po_number_idx`(`po_number`),
    INDEX `purchase_order_vendor_id_idx`(`vendor_id`),
    INDEX `purchase_order_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `po_line_item` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `po_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NULL,
    `description` VARCHAR(500) NULL,
    `quantity` DECIMAL(18, 4) NOT NULL,
    `unit_price` DECIMAL(18, 4) NOT NULL,
    `line_amount` DECIMAL(18, 2) NOT NULL,
    `tax_percent` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `tax_amount` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `line_total` DECIMAL(18, 2) NOT NULL,
    `delivery_date` DATE NULL,
    `quantity_received` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `quantity_invoiced` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `po_line_item_tenant_id_idx`(`tenant_id`),
    INDEX `po_line_item_po_id_idx`(`po_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `goods_receipt` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `grn_number` VARCHAR(50) NOT NULL,
    `po_id` VARCHAR(36) NOT NULL,
    `vendor_id` VARCHAR(36) NULL,
    `receipt_date` DATE NOT NULL,
    `warehouse_id` VARCHAR(36) NULL,
    `received_by` VARCHAR(36) NULL,
    `inspection_status` VARCHAR(50) NULL,
    `acceptance_status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `notes` TEXT NULL,
    `reference_number` VARCHAR(100) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `goods_receipt_grn_number_key`(`grn_number`),
    INDEX `goods_receipt_tenant_id_idx`(`tenant_id`),
    INDEX `goods_receipt_grn_number_idx`(`grn_number`),
    INDEX `goods_receipt_po_id_idx`(`po_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `grn_line_item` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `grn_id` VARCHAR(36) NOT NULL,
    `po_line_item_id` VARCHAR(36) NULL,
    `inventory_item_id` VARCHAR(36) NULL,
    `quantity_received` DECIMAL(18, 4) NOT NULL,
    `quantity_accepted` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `quantity_rejected` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `rejection_reason` TEXT NULL,
    `batch_number` VARCHAR(100) NULL,
    `expiry_date` DATE NULL,
    `serial_numbers` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `grn_line_item_tenant_id_idx`(`tenant_id`),
    INDEX `grn_line_item_grn_id_idx`(`grn_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `vendor_invoice` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `vendor_id` VARCHAR(36) NOT NULL,
    `invoice_number` VARCHAR(100) NOT NULL,
    `invoice_date` DATE NOT NULL,
    `due_date` DATE NULL,
    `total_amount` DECIMAL(18, 2) NOT NULL,
    `tax_amount` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `discount_amount` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `payable_amount` DECIMAL(18, 2) NOT NULL,
    `currency` VARCHAR(3) NOT NULL DEFAULT 'USD',
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `description` TEXT NULL,
    `notes` TEXT NULL,
    `attachments` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `vendor_invoice_tenant_id_idx`(`tenant_id`),
    INDEX `vendor_invoice_vendor_id_idx`(`vendor_id`),
    INDEX `vendor_invoice_invoice_date_idx`(`invoice_date`),
    UNIQUE INDEX `vendor_invoice_tenant_id_invoice_number_key`(`tenant_id`, `invoice_number`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `warehouse` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `warehouse_code` VARCHAR(50) NOT NULL,
    `warehouse_name` VARCHAR(255) NOT NULL,
    `warehouse_type` VARCHAR(50) NULL,
    `address` TEXT NULL,
    `city` VARCHAR(100) NULL,
    `state` VARCHAR(100) NULL,
    `country` VARCHAR(100) NULL,
    `postal_code` VARCHAR(20) NULL,
    `manager_id` VARCHAR(36) NULL,
    `capacity` DECIMAL(18, 2) NULL,
    `current_utilization` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `gl_inventory_account_id` VARCHAR(36) NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `warehouse_tenant_id_idx`(`tenant_id`),
    INDEX `warehouse_warehouse_code_idx`(`warehouse_code`),
    INDEX `warehouse_is_active_idx`(`is_active`),
    UNIQUE INDEX `warehouse_tenant_id_warehouse_code_key`(`tenant_id`, `warehouse_code`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `inventory_item` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `sku` VARCHAR(50) NOT NULL,
    `item_name` VARCHAR(255) NOT NULL,
    `item_description` TEXT NULL,
    `item_category` VARCHAR(100) NULL,
    `item_type` VARCHAR(50) NULL,
    `unit_of_measure` VARCHAR(20) NULL,
    `reorder_level` DECIMAL(18, 4) NULL,
    `reorder_quantity` DECIMAL(18, 4) NULL,
    `safety_stock` DECIMAL(18, 4) NULL,
    `lead_time_days` INTEGER NULL,
    `hsn_code` VARCHAR(50) NULL,
    `is_serialized` BOOLEAN NOT NULL DEFAULT false,
    `is_batch_tracked` BOOLEAN NOT NULL DEFAULT false,
    `item_status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `gl_inventory_account_id` VARCHAR(36) NULL,
    `gl_expense_account_id` VARCHAR(36) NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `inventory_item_tenant_id_idx`(`tenant_id`),
    INDEX `inventory_item_sku_idx`(`sku`),
    INDEX `inventory_item_item_category_idx`(`item_category`),
    INDEX `inventory_item_item_status_idx`(`item_status`),
    UNIQUE INDEX `inventory_item_tenant_id_sku_key`(`tenant_id`, `sku`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `inventory_item_vendor` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `vendor_id` VARCHAR(36) NOT NULL,
    `vendor_sku` VARCHAR(50) NULL,
    `vendor_part_number` VARCHAR(100) NULL,
    `lead_time_days` INTEGER NULL,
    `minimum_order_quantity` DECIMAL(18, 4) NULL,
    `unit_price` DECIMAL(18, 4) NULL,
    `last_price_date` DATE NULL,
    `preferred_vendor` BOOLEAN NOT NULL DEFAULT false,
    `quality_rating` DECIMAL(3, 1) NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `inventory_item_vendor_tenant_id_idx`(`tenant_id`),
    INDEX `inventory_item_vendor_inventory_item_id_idx`(`inventory_item_id`),
    INDEX `inventory_item_vendor_vendor_id_idx`(`vendor_id`),
    UNIQUE INDEX `inventory_item_vendor_inventory_item_id_vendor_id_key`(`inventory_item_id`, `vendor_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `stock_level` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `warehouse_id` VARCHAR(36) NOT NULL,
    `quantity_on_hand` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `quantity_reserved` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `quantity_available` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `quantity_in_transit` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `last_counted_date` DATE NULL,
    `recount_required` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `stock_level_tenant_id_idx`(`tenant_id`),
    INDEX `stock_level_inventory_item_id_idx`(`inventory_item_id`),
    INDEX `stock_level_warehouse_id_idx`(`warehouse_id`),
    INDEX `stock_level_recount_required_idx`(`recount_required`),
    UNIQUE INDEX `stock_level_inventory_item_id_warehouse_id_key`(`inventory_item_id`, `warehouse_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `stock_movement` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `warehouse_id` VARCHAR(36) NOT NULL,
    `movement_type` VARCHAR(50) NULL,
    `movement_date` DATE NOT NULL,
    `quantity_change` DECIMAL(18, 4) NULL,
    `reference_type` VARCHAR(50) NULL,
    `reference_id` VARCHAR(36) NULL,
    `from_location` VARCHAR(100) NULL,
    `to_location` VARCHAR(100) NULL,
    `batch_number` VARCHAR(100) NULL,
    `serial_numbers` TEXT NULL,
    `unit_price` DECIMAL(18, 4) NULL,
    `total_value` DECIMAL(18, 2) NULL,
    `reason_code` VARCHAR(50) NULL,
    `notes` TEXT NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `stock_movement_tenant_id_idx`(`tenant_id`),
    INDEX `stock_movement_inventory_item_id_idx`(`inventory_item_id`),
    INDEX `stock_movement_warehouse_id_idx`(`warehouse_id`),
    INDEX `stock_movement_movement_date_idx`(`movement_date`),
    INDEX `stock_movement_movement_type_idx`(`movement_type`),
    INDEX `stock_movement_reference_id_idx`(`reference_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `inventory_batch` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `batch_number` VARCHAR(100) NOT NULL,
    `manufacture_date` DATE NULL,
    `expiry_date` DATE NULL,
    `quantity_received` DECIMAL(18, 4) NULL,
    `quantity_remaining` DECIMAL(18, 4) NULL,
    `purchase_order_id` VARCHAR(36) NULL,
    `supplier_batch_number` VARCHAR(100) NULL,
    `quality_status` VARCHAR(50) NULL,
    `storage_location` VARCHAR(100) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `inventory_batch_tenant_id_idx`(`tenant_id`),
    INDEX `inventory_batch_inventory_item_id_idx`(`inventory_item_id`),
    INDEX `inventory_batch_batch_number_idx`(`batch_number`),
    INDEX `inventory_batch_expiry_date_idx`(`expiry_date`),
    INDEX `inventory_batch_quality_status_idx`(`quality_status`),
    UNIQUE INDEX `inventory_batch_tenant_id_batch_number_key`(`tenant_id`, `batch_number`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `inventory_serial` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `serial_number` VARCHAR(100) NOT NULL,
    `batch_id` VARCHAR(36) NULL,
    `purchase_order_id` VARCHAR(36) NULL,
    `warranty_start_date` DATE NULL,
    `warranty_end_date` DATE NULL,
    `current_location` VARCHAR(100) NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'in_stock',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `inventory_serial_serial_number_key`(`serial_number`),
    INDEX `inventory_serial_tenant_id_idx`(`tenant_id`),
    INDEX `inventory_serial_inventory_item_id_idx`(`inventory_item_id`),
    INDEX `inventory_serial_serial_number_idx`(`serial_number`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `inventory_valuation` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `warehouse_id` VARCHAR(36) NOT NULL,
    `valuation_date` DATE NOT NULL,
    `valuation_method` VARCHAR(50) NOT NULL,
    `quantity_on_hand` DECIMAL(18, 4) NOT NULL,
    `unit_price` DECIMAL(18, 4) NOT NULL,
    `total_value` DECIMAL(18, 2) NOT NULL,
    `notes` TEXT NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `inventory_valuation_tenant_id_idx`(`tenant_id`),
    INDEX `inventory_valuation_inventory_item_id_idx`(`inventory_item_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `stock_adjustment` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `adjustment_number` VARCHAR(50) NOT NULL,
    `adjustment_date` DATE NOT NULL,
    `reason` VARCHAR(100) NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `approved_by` VARCHAR(36) NULL,
    `approval_date` DATETIME(3) NULL,
    `notes` TEXT NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `stock_adjustment_adjustment_number_key`(`adjustment_number`),
    INDEX `stock_adjustment_tenant_id_idx`(`tenant_id`),
    INDEX `stock_adjustment_adjustment_number_idx`(`adjustment_number`),
    INDEX `stock_adjustment_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `stock_adjustment_line` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `adjustment_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `warehouse_id` VARCHAR(36) NOT NULL,
    `quantity_variance` DECIMAL(18, 4) NOT NULL,
    `reason` VARCHAR(100) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `stock_adjustment_line_tenant_id_idx`(`tenant_id`),
    INDEX `stock_adjustment_line_adjustment_id_idx`(`adjustment_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `physical_inventory` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_number` VARCHAR(50) NOT NULL,
    `warehouse_id` VARCHAR(36) NULL,
    `count_date` DATE NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `counters_names` TEXT NULL,
    `completion_percentage` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `notes` TEXT NULL,
    `started_by` VARCHAR(36) NULL,
    `completed_by` VARCHAR(36) NULL,
    `completed_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `physical_inventory_inventory_number_key`(`inventory_number`),
    INDEX `physical_inventory_tenant_id_idx`(`tenant_id`),
    INDEX `physical_inventory_inventory_number_idx`(`inventory_number`),
    INDEX `physical_inventory_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `physical_inventory_detail` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `physical_inventory_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `quantity_expected` DECIMAL(18, 4) NULL,
    `quantity_counted` DECIMAL(18, 4) NULL,
    `variance` DECIMAL(18, 4) NULL,
    `variance_percent` DECIMAL(5, 2) NULL,
    `notes` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `physical_inventory_detail_tenant_id_idx`(`tenant_id`),
    INDEX `physical_inventory_detail_physical_inventory_id_idx`(`physical_inventory_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `inventory_transfer` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `transfer_number` VARCHAR(50) NOT NULL,
    `from_warehouse_id` VARCHAR(36) NULL,
    `to_warehouse_id` VARCHAR(36) NULL,
    `transfer_date` DATE NOT NULL,
    `expected_delivery_date` DATE NULL,
    `actual_delivery_date` DATE NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `notes` TEXT NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `inventory_transfer_transfer_number_key`(`transfer_number`),
    INDEX `inventory_transfer_tenant_id_idx`(`tenant_id`),
    INDEX `inventory_transfer_transfer_number_idx`(`transfer_number`),
    INDEX `inventory_transfer_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `inventory_transfer_line` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `transfer_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `quantity_transfered` DECIMAL(18, 4) NOT NULL,
    `quantity_received` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `quantity_shortage` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `inventory_transfer_line_tenant_id_idx`(`tenant_id`),
    INDEX `inventory_transfer_line_transfer_id_idx`(`transfer_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `min_stock_alert` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `warehouse_id` VARCHAR(36) NOT NULL,
    `min_stock_level` DECIMAL(18, 4) NOT NULL,
    `max_stock_level` DECIMAL(18, 4) NOT NULL,
    `current_quantity` DECIMAL(18, 4) NULL,
    `alert_status` VARCHAR(50) NOT NULL,
    `last_alert_date` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `min_stock_alert_tenant_id_idx`(`tenant_id`),
    INDEX `min_stock_alert_alert_status_idx`(`alert_status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `inventory_damage` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `inventory_item_id` VARCHAR(36) NOT NULL,
    `warehouse_id` VARCHAR(36) NULL,
    `damage_date` DATE NOT NULL,
    `quantity_damaged` DECIMAL(18, 4) NOT NULL,
    `damage_reason` VARCHAR(255) NULL,
    `estimated_value` DECIMAL(18, 2) NULL,
    `notes` TEXT NULL,
    `reported_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `inventory_damage_tenant_id_idx`(`tenant_id`),
    INDEX `inventory_damage_damage_date_idx`(`damage_date`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `role_permission` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `role_id` VARCHAR(36) NOT NULL,
    `permission_id` VARCHAR(36) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `role_permission_tenant_id_idx`(`tenant_id`),
    INDEX `role_permission_role_id_idx`(`role_id`),
    INDEX `role_permission_permission_id_idx`(`permission_id`),
    UNIQUE INDEX `role_permission_role_id_permission_id_key`(`role_id`, `permission_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `resource` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `resource_name` VARCHAR(100) NOT NULL,
    `resource_type` VARCHAR(50) NULL,
    `endpoint` VARCHAR(255) NULL,
    `description` TEXT NULL,
    `is_protected` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `resource_tenant_id_idx`(`tenant_id`),
    INDEX `resource_resource_type_idx`(`resource_type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `access_log` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` CHAR(36) NULL,
    `resource_id` VARCHAR(36) NULL,
    `action` VARCHAR(50) NULL,
    `status` VARCHAR(50) NULL,
    `ip_address` VARCHAR(45) NULL,
    `user_agent` TEXT NULL,
    `timestamp` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `access_log_tenant_id_idx`(`tenant_id`),
    INDEX `access_log_user_id_idx`(`user_id`),
    INDEX `access_log_timestamp_idx`(`timestamp`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `resource_access` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `role_id` VARCHAR(36) NOT NULL,
    `resource_id` VARCHAR(36) NOT NULL,
    `access_level` VARCHAR(50) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `resource_access_tenant_id_idx`(`tenant_id`),
    INDEX `resource_access_role_id_idx`(`role_id`),
    INDEX `resource_access_resource_id_idx`(`resource_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `time_based_permission` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `role_id` VARCHAR(36) NOT NULL,
    `valid_from` DATETIME(3) NOT NULL,
    `valid_to` DATETIME(3) NOT NULL,
    `days_of_week` VARCHAR(50) NULL,
    `time_from` VARCHAR(10) NULL,
    `time_to` VARCHAR(10) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `time_based_permission_tenant_id_idx`(`tenant_id`),
    INDEX `time_based_permission_role_id_idx`(`role_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `field_level_permission` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `permission_id` VARCHAR(36) NOT NULL,
    `entity_name` VARCHAR(100) NOT NULL,
    `field_name` VARCHAR(100) NOT NULL,
    `access_type` VARCHAR(50) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `field_level_permission_tenant_id_idx`(`tenant_id`),
    INDEX `field_level_permission_permission_id_idx`(`permission_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `role_delegation` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `role_id` VARCHAR(36) NOT NULL,
    `delegated_to` CHAR(36) NOT NULL,
    `delegated_by` CHAR(36) NOT NULL,
    `valid_from` DATETIME(3) NOT NULL,
    `valid_to` DATETIME(3) NOT NULL,
    `reason` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `role_delegation_tenant_id_idx`(`tenant_id`),
    INDEX `role_delegation_role_id_idx`(`role_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `bulk_permission_log` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `batch_number` VARCHAR(100) NOT NULL,
    `operation` VARCHAR(50) NOT NULL,
    `role_ids` TEXT NULL,
    `permission_ids` TEXT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `total_affected` INTEGER NOT NULL,
    `success_count` INTEGER NOT NULL DEFAULT 0,
    `failure_count` INTEGER NOT NULL DEFAULT 0,
    `created_by` VARCHAR(36) NULL,
    `completed_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `bulk_permission_log_batch_number_key`(`batch_number`),
    INDEX `bulk_permission_log_tenant_id_idx`(`tenant_id`),
    INDEX `bulk_permission_log_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `chart_of_account` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `account_code` VARCHAR(50) NOT NULL,
    `account_name` VARCHAR(255) NOT NULL,
    `account_type` VARCHAR(50) NOT NULL,
    `account_subtype` VARCHAR(50) NULL,
    `description` TEXT NULL,
    `balance` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `parent_account_id` VARCHAR(36) NULL,
    `is_header` BOOLEAN NOT NULL DEFAULT false,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `chart_of_account_account_code_key`(`account_code`),
    INDEX `chart_of_account_tenant_id_idx`(`tenant_id`),
    INDEX `chart_of_account_account_code_idx`(`account_code`),
    INDEX `chart_of_account_account_type_idx`(`account_type`),
    UNIQUE INDEX `chart_of_account_tenant_id_account_code_key`(`tenant_id`, `account_code`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `financial_period` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `period_name` VARCHAR(100) NOT NULL,
    `period_type` VARCHAR(50) NOT NULL,
    `start_date` DATE NOT NULL,
    `end_date` DATE NOT NULL,
    `is_closed` BOOLEAN NOT NULL DEFAULT false,
    `closed_by` VARCHAR(36) NULL,
    `closed_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `financial_period_tenant_id_idx`(`tenant_id`),
    INDEX `financial_period_start_date_end_date_idx`(`start_date`, `end_date`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `journal_entry` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `journal_number` VARCHAR(50) NOT NULL,
    `entry_date` DATE NOT NULL,
    `description` TEXT NULL,
    `reference_number` VARCHAR(100) NULL,
    `reference_type` VARCHAR(50) NULL,
    `total_debit` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `total_credit` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `approved_by` VARCHAR(36) NULL,
    `approved_at` DATETIME(3) NULL,
    `posted_by` VARCHAR(36) NULL,
    `posted_at` DATETIME(3) NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `journal_entry_journal_number_key`(`journal_number`),
    INDEX `journal_entry_tenant_id_idx`(`tenant_id`),
    INDEX `journal_entry_journal_number_idx`(`journal_number`),
    INDEX `journal_entry_entry_date_idx`(`entry_date`),
    INDEX `journal_entry_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `journal_entry_detail` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `journal_entry_id` VARCHAR(36) NOT NULL,
    `account_id` VARCHAR(36) NOT NULL,
    `debit_amount` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `credit_amount` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `description` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `journal_entry_detail_tenant_id_idx`(`tenant_id`),
    INDEX `journal_entry_detail_journal_entry_id_idx`(`journal_entry_id`),
    INDEX `journal_entry_detail_account_id_idx`(`account_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `gl_account_balance` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `account_id` VARCHAR(36) NOT NULL,
    `fiscal_year` INTEGER NOT NULL,
    `period_number` INTEGER NOT NULL,
    `opening_balance` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `debit_total` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `credit_total` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `closing_balance` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `gl_account_balance_tenant_id_idx`(`tenant_id`),
    INDEX `gl_account_balance_account_id_idx`(`account_id`),
    UNIQUE INDEX `gl_account_balance_tenant_id_account_id_fiscal_year_period_n_key`(`tenant_id`, `account_id`, `fiscal_year`, `period_number`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `trial_balance` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `trial_balance_date` DATE NOT NULL,
    `total_debits` DECIMAL(18, 2) NOT NULL,
    `total_credits` DECIMAL(18, 2) NOT NULL,
    `status` VARCHAR(50) NOT NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `trial_balance_tenant_id_idx`(`tenant_id`),
    INDEX `trial_balance_trial_balance_date_idx`(`trial_balance_date`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `income_statement` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `statement_date` DATE NOT NULL,
    `revenue` DECIMAL(18, 2) NOT NULL,
    `cost_of_sales` DECIMAL(18, 2) NOT NULL,
    `gross_profit` DECIMAL(18, 2) NOT NULL,
    `operating_expenses` DECIMAL(18, 2) NOT NULL,
    `operating_profit` DECIMAL(18, 2) NOT NULL,
    `net_income` DECIMAL(18, 2) NOT NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `income_statement_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `balance_sheet` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `sheet_date` DATE NOT NULL,
    `total_assets` DECIMAL(18, 2) NOT NULL,
    `total_liabilities` DECIMAL(18, 2) NOT NULL,
    `total_equity` DECIMAL(18, 2) NOT NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `balance_sheet_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `gl_posting_template` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `template_name` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `transaction_type` VARCHAR(50) NOT NULL,
    `source_account_id` VARCHAR(36) NOT NULL,
    `target_account_id` VARCHAR(36) NOT NULL,
    `debit_account_id` VARCHAR(36) NOT NULL,
    `credit_account_id` VARCHAR(36) NOT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `gl_posting_template_tenant_id_idx`(`tenant_id`),
    INDEX `gl_posting_template_transaction_type_idx`(`transaction_type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `payroll_gl_posting` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `payroll_id` VARCHAR(36) NOT NULL,
    `journal_entry_id` VARCHAR(36) NOT NULL,
    `salary_expense_account_id` VARCHAR(36) NOT NULL,
    `salary_payable_account_id` VARCHAR(36) NOT NULL,
    `total_salary` DECIMAL(15, 2) NOT NULL,
    `total_deductions` DECIMAL(15, 2) NOT NULL,
    `posting_date` DATE NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `payroll_gl_posting_tenant_id_idx`(`tenant_id`),
    INDEX `payroll_gl_posting_payroll_id_idx`(`payroll_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `purchase_gl_posting` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `po_id` VARCHAR(36) NOT NULL,
    `journal_entry_id` VARCHAR(36) NOT NULL,
    `inventory_account_id` VARCHAR(36) NOT NULL,
    `expense_account_id` VARCHAR(36) NOT NULL,
    `payable_account_id` VARCHAR(36) NOT NULL,
    `net_amount` DECIMAL(18, 2) NOT NULL,
    `posting_date` DATE NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `purchase_gl_posting_tenant_id_idx`(`tenant_id`),
    INDEX `purchase_gl_posting_po_id_idx`(`po_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sales_gl_posting` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `invoice_id` VARCHAR(36) NOT NULL,
    `journal_entry_id` VARCHAR(36) NOT NULL,
    `revenue_account_id` VARCHAR(36) NOT NULL,
    `receivable_account_id` VARCHAR(36) NOT NULL,
    `net_amount` DECIMAL(18, 2) NOT NULL,
    `posting_date` DATE NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `sales_gl_posting_tenant_id_idx`(`tenant_id`),
    INDEX `sales_gl_posting_invoice_id_idx`(`invoice_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `bank_statement` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `bank_account_id` VARCHAR(36) NOT NULL,
    `statement_date` DATE NOT NULL,
    `opening_balance` DECIMAL(18, 2) NOT NULL,
    `closing_balance` DECIMAL(18, 2) NOT NULL,
    `total_deposits` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `total_withdrawals` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `statement_file` VARCHAR(255) NULL,
    `uploaded_by` VARCHAR(36) NULL,
    `uploaded_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `bank_statement_tenant_id_idx`(`tenant_id`),
    INDEX `bank_statement_bank_account_id_idx`(`bank_account_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `bank_transaction` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `bank_statement_id` VARCHAR(36) NOT NULL,
    `transaction_date` DATE NOT NULL,
    `description` VARCHAR(255) NOT NULL,
    `amount` DECIMAL(18, 2) NOT NULL,
    `transaction_type` VARCHAR(50) NOT NULL,
    `reference_number` VARCHAR(100) NULL,
    `balance` DECIMAL(18, 2) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `bank_transaction_tenant_id_idx`(`tenant_id`),
    INDEX `bank_transaction_bank_statement_id_idx`(`bank_statement_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `bank_reconciliation_match` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `bank_transaction_id` VARCHAR(36) NOT NULL,
    `journal_entry_id` VARCHAR(36) NULL,
    `match_date` DATE NOT NULL,
    `amount` DECIMAL(18, 2) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'matched',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `bank_reconciliation_match_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `uncleared_item` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `bank_account_id` VARCHAR(36) NOT NULL,
    `journal_entry_id` VARCHAR(36) NULL,
    `item_date` DATE NOT NULL,
    `description` VARCHAR(255) NULL,
    `amount` DECIMAL(18, 2) NOT NULL,
    `item_type` VARCHAR(50) NOT NULL,
    `days_outstanding` INTEGER NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `uncleared_item_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `cash_flow_forecast` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `forecast_date` DATE NOT NULL,
    `forecast_period` VARCHAR(50) NOT NULL,
    `projected_inflows` DECIMAL(18, 2) NOT NULL,
    `projected_outflows` DECIMAL(18, 2) NOT NULL,
    `projected_balance` DECIMAL(18, 2) NOT NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `cash_flow_forecast_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `cash_flow_item` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `forecast_id` VARCHAR(36) NOT NULL,
    `item_type` VARCHAR(50) NOT NULL,
    `description` VARCHAR(255) NULL,
    `amount` DECIMAL(18, 2) NOT NULL,
    `due_date` DATE NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `cash_flow_item_tenant_id_idx`(`tenant_id`),
    INDEX `cash_flow_item_forecast_id_idx`(`forecast_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `employee` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_code` VARCHAR(50) NOT NULL,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100) NOT NULL,
    `email` VARCHAR(255) NULL,
    `phone` VARCHAR(20) NULL,
    `department_id` VARCHAR(36) NULL,
    `designation_id` VARCHAR(36) NULL,
    `joining_date` DATE NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `employee_employee_code_key`(`employee_code`),
    INDEX `employee_tenant_id_idx`(`tenant_id`),
    INDEX `employee_employee_code_idx`(`employee_code`),
    INDEX `employee_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `employee_details` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` CHAR(36) NOT NULL,
    `father_name` VARCHAR(255) NULL,
    `mother_name` VARCHAR(255) NULL,
    `spouse` VARCHAR(255) NULL,
    `children` INTEGER NOT NULL DEFAULT 0,
    `date_of_birth` DATE NOT NULL,
    `blood_group` VARCHAR(5) NULL,
    `aadhar` VARCHAR(50) NULL,
    `pan` VARCHAR(50) NULL,
    `pf_number` VARCHAR(50) NULL,
    `esi_number` VARCHAR(50) NULL,
    `bank_account_number` VARCHAR(50) NULL,
    `bank_name` VARCHAR(100) NULL,
    `ifsc_code` VARCHAR(20) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `employee_details_employee_id_key`(`employee_id`),
    INDEX `employee_details_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `leave_type` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `leave_type_name` VARCHAR(100) NOT NULL,
    `allowed_days` INTEGER NOT NULL,
    `accrual_rate` DECIMAL(5, 2) NULL,
    `is_carry_over` BOOLEAN NOT NULL DEFAULT false,
    `max_carry_over` INTEGER NOT NULL DEFAULT 0,
    `requires_approval` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `leave_type_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `leave_type_tenant_id_leave_type_name_key`(`tenant_id`, `leave_type_name`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `leave_balance` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` CHAR(36) NOT NULL,
    `leave_type_id` VARCHAR(36) NOT NULL,
    `year` INTEGER NOT NULL,
    `opening_balance` DECIMAL(5, 2) NOT NULL,
    `accrued` DECIMAL(5, 2) NOT NULL,
    `utilized` DECIMAL(5, 2) NOT NULL,
    `carried` DECIMAL(5, 2) NOT NULL,
    `closing_balance` DECIMAL(5, 2) NOT NULL,
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `leave_balance_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `leave_balance_tenant_id_employee_id_leave_type_id_year_key`(`tenant_id`, `employee_id`, `leave_type_id`, `year`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `designation` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `designation_name` VARCHAR(100) NOT NULL,
    `designation_level` INTEGER NULL,
    `description` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `designation_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `designation_tenant_id_designation_name_key`(`tenant_id`, `designation_name`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `department` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `department_name` VARCHAR(100) NOT NULL,
    `department_code` VARCHAR(50) NOT NULL,
    `head_id` VARCHAR(36) NULL,
    `cost_center` VARCHAR(50) NULL,
    `description` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `department_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `department_tenant_id_department_code_key`(`tenant_id`, `department_code`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `salary_component` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `component_name` VARCHAR(100) NOT NULL,
    `component_type` VARCHAR(50) NOT NULL,
    `taxable` BOOLEAN NOT NULL DEFAULT false,
    `formula` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `salary_component_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `salary_component_tenant_id_component_name_key`(`tenant_id`, `component_name`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `payroll_run` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `month` INTEGER NOT NULL,
    `year` INTEGER NOT NULL,
    `payroll_status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `total_employees` INTEGER NOT NULL,
    `total_salary` DECIMAL(15, 2) NOT NULL,
    `total_deductions` DECIMAL(15, 2) NOT NULL,
    `net_payable` DECIMAL(15, 2) NOT NULL,
    `freeze_date` DATETIME(3) NULL,
    `processed_by` VARCHAR(36) NULL,
    `processed_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `payroll_run_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `payroll_run_tenant_id_month_year_key`(`tenant_id`, `month`, `year`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `epf_configuration` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `epf_number` VARCHAR(50) NOT NULL,
    `epf_registration_date` DATE NOT NULL,
    `epf_establishment_code` VARCHAR(50) NOT NULL,
    `company_contribution_rate` DECIMAL(5, 2) NOT NULL,
    `employee_contribution_rate` DECIMAL(5, 2) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `epf_configuration_epf_number_key`(`epf_number`),
    INDEX `epf_configuration_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `esi_configuration` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `esi_number` VARCHAR(50) NOT NULL,
    `esi_registration_date` DATE NOT NULL,
    `company_contribution_rate` DECIMAL(5, 2) NOT NULL,
    `employee_contribution_rate` DECIMAL(5, 2) NOT NULL,
    `wage_limit` DECIMAL(10, 2) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `esi_configuration_esi_number_key`(`esi_number`),
    INDEX `esi_configuration_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `employee_epf_registration` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` CHAR(36) NOT NULL,
    `epf_number` VARCHAR(50) NOT NULL,
    `enrollment_date` DATE NOT NULL,
    `member_status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `employee_epf_registration_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `employee_esi_registration` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` CHAR(36) NOT NULL,
    `esi_number` VARCHAR(50) NOT NULL,
    `enrollment_date` DATE NOT NULL,
    `member_status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `employee_esi_registration_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `epf_contribution` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` CHAR(36) NOT NULL,
    `payroll_month` INTEGER NOT NULL,
    `payroll_year` INTEGER NOT NULL,
    `company_contribution` DECIMAL(12, 2) NOT NULL,
    `employee_contribution` DECIMAL(12, 2) NOT NULL,
    `submission_date` DATE NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `epf_contribution_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `esi_contribution` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` CHAR(36) NOT NULL,
    `payroll_month` INTEGER NOT NULL,
    `payroll_year` INTEGER NOT NULL,
    `company_contribution` DECIMAL(12, 2) NOT NULL,
    `employee_contribution` DECIMAL(12, 2) NOT NULL,
    `submission_date` DATE NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `esi_contribution_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `compliance_record` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `record_type` VARCHAR(100) NOT NULL,
    `record_value` TEXT NULL,
    `due_date` DATE NULL,
    `completion_date` DATE NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `compliance_record_tenant_id_idx`(`tenant_id`),
    INDEX `compliance_record_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `employee_compliance_checklist` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `employee_id` CHAR(36) NOT NULL,
    `item_description` VARCHAR(255) NOT NULL,
    `is_completed` BOOLEAN NOT NULL DEFAULT false,
    `completion_date` DATE NULL,
    `notes` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `employee_compliance_checklist_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sales_customer` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `customer_code` VARCHAR(50) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NULL,
    `phone` VARCHAR(20) NULL,
    `address` TEXT NULL,
    `city` VARCHAR(100) NULL,
    `state` VARCHAR(100) NULL,
    `country` VARCHAR(100) NULL,
    `postal_code` VARCHAR(20) NULL,
    `customer_type` VARCHAR(50) NULL,
    `credit_limit` DECIMAL(15, 2) NULL,
    `payment_terms` VARCHAR(100) NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `sales_customer_customer_code_key`(`customer_code`),
    INDEX `sales_customer_tenant_id_idx`(`tenant_id`),
    INDEX `sales_customer_customer_code_idx`(`customer_code`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `quotation_line_item` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `quotation_id` VARCHAR(36) NOT NULL,
    `description` VARCHAR(500) NULL,
    `quantity` DECIMAL(18, 4) NOT NULL,
    `unit_price` DECIMAL(18, 4) NOT NULL,
    `line_amount` DECIMAL(18, 2) NOT NULL,
    `tax` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `line_total` DECIMAL(18, 2) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `quotation_line_item_tenant_id_idx`(`tenant_id`),
    INDEX `quotation_line_item_quotation_id_idx`(`quotation_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sales_order` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `order_number` VARCHAR(50) NOT NULL,
    `customer_id` VARCHAR(36) NULL,
    `quotation_id` VARCHAR(36) NULL,
    `order_date` DATE NOT NULL,
    `delivery_date` DATE NULL,
    `total_amount` DECIMAL(18, 2) NULL,
    `tax` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `net_amount` DECIMAL(18, 2) NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `sales_order_order_number_key`(`order_number`),
    INDEX `sales_order_tenant_id_idx`(`tenant_id`),
    INDEX `sales_order_order_number_idx`(`order_number`),
    INDEX `sales_order_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sales_order_line_item` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `order_id` VARCHAR(36) NOT NULL,
    `description` VARCHAR(500) NULL,
    `quantity` DECIMAL(18, 4) NOT NULL,
    `unit_price` DECIMAL(18, 4) NOT NULL,
    `line_amount` DECIMAL(18, 2) NOT NULL,
    `tax` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `line_total` DECIMAL(18, 2) NOT NULL,
    `quantity_delivered` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `sales_order_line_item_tenant_id_idx`(`tenant_id`),
    INDEX `sales_order_line_item_order_id_idx`(`order_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sales_invoice` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `invoice_number` VARCHAR(50) NOT NULL,
    `order_id` VARCHAR(36) NULL,
    `customer_id` VARCHAR(36) NULL,
    `invoice_date` DATE NOT NULL,
    `due_date` DATE NULL,
    `total_amount` DECIMAL(18, 2) NULL,
    `tax` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `net_amount` DECIMAL(18, 2) NULL,
    `amount_paid` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `sales_invoice_invoice_number_key`(`invoice_number`),
    INDEX `sales_invoice_tenant_id_idx`(`tenant_id`),
    INDEX `sales_invoice_invoice_number_idx`(`invoice_number`),
    INDEX `sales_invoice_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sales_invoice_line_item` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `invoice_id` VARCHAR(36) NOT NULL,
    `description` VARCHAR(500) NULL,
    `quantity` DECIMAL(18, 4) NOT NULL,
    `unit_price` DECIMAL(18, 4) NOT NULL,
    `line_amount` DECIMAL(18, 2) NOT NULL,
    `tax` DECIMAL(18, 2) NOT NULL DEFAULT 0,
    `line_total` DECIMAL(18, 2) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `sales_invoice_line_item_tenant_id_idx`(`tenant_id`),
    INDEX `sales_invoice_line_item_invoice_id_idx`(`invoice_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sales_return` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `return_number` VARCHAR(50) NOT NULL,
    `invoice_id` VARCHAR(36) NULL,
    `return_date` DATE NOT NULL,
    `reason` TEXT NULL,
    `total_amount` DECIMAL(18, 2) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `sales_return_return_number_key`(`return_number`),
    INDEX `sales_return_tenant_id_idx`(`tenant_id`),
    INDEX `sales_return_return_number_idx`(`return_number`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sales_status_history` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `lead_id` VARCHAR(36) NOT NULL,
    `old_status` VARCHAR(50) NOT NULL,
    `new_status` VARCHAR(50) NOT NULL,
    `changed_by` VARCHAR(36) NULL,
    `remarks` TEXT NULL,
    `changed_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `sales_status_history_tenant_id_idx`(`tenant_id`),
    INDEX `sales_status_history_lead_id_idx`(`lead_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sales_follow_up` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `lead_id` CHAR(36) NOT NULL,
    `follow_up_date` DATE NOT NULL,
    `follow_up_type` VARCHAR(50) NOT NULL,
    `notes` TEXT NULL,
    `result` VARCHAR(100) NULL,
    `next_follow_up_date` DATE NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `sales_follow_up_tenant_id_idx`(`tenant_id`),
    INDEX `sales_follow_up_lead_id_idx`(`lead_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sales_target` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `period` VARCHAR(50) NOT NULL,
    `sales_person_id` VARCHAR(36) NOT NULL,
    `target_amount` DECIMAL(15, 2) NOT NULL,
    `achieved_amount` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `sales_target_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `lead_status_type` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `status_name` VARCHAR(100) NOT NULL,
    `status_code` VARCHAR(50) NOT NULL,
    `description` TEXT NULL,
    `color_code` VARCHAR(10) NULL,
    `sequence_order` INTEGER NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `lead_status_type_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `lead_status_type_tenant_id_status_code_key`(`tenant_id`, `status_code`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `lead_conversion_rate` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `source_id` VARCHAR(36) NOT NULL,
    `status_id` VARCHAR(36) NOT NULL,
    `conversion_percent` DECIMAL(5, 2) NOT NULL,
    `total_leads` INTEGER NOT NULL DEFAULT 0,
    `converted_leads` INTEGER NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `lead_conversion_rate_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `property_project` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_name` VARCHAR(255) NOT NULL,
    `project_code` VARCHAR(50) NOT NULL,
    `location` VARCHAR(255) NULL,
    `city` VARCHAR(100) NULL,
    `state` VARCHAR(100) NULL,
    `country` VARCHAR(100) NULL,
    `postal_code` VARCHAR(20) NULL,
    `project_type` VARCHAR(50) NULL,
    `total_area` DECIMAL(15, 2) NULL,
    `total_units` INTEGER NOT NULL DEFAULT 0,
    `launch_date` DATE NULL,
    `completion_date` DATE NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'planning',
    `developer_id` VARCHAR(36) NULL,
    `description` LONGTEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `property_project_project_code_key`(`project_code`),
    INDEX `property_project_tenant_id_idx`(`tenant_id`),
    INDEX `property_project_project_code_idx`(`project_code`),
    INDEX `property_project_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `property_block` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `block_name` VARCHAR(100) NOT NULL,
    `block_code` VARCHAR(50) NOT NULL,
    `total_units` INTEGER NOT NULL DEFAULT 0,
    `units_left` INTEGER NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `property_block_tenant_id_idx`(`tenant_id`),
    INDEX `property_block_project_id_idx`(`project_id`),
    UNIQUE INDEX `property_block_project_id_block_code_key`(`project_id`, `block_code`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `property_unit` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `block_id` VARCHAR(36) NULL,
    `unit_number` VARCHAR(50) NOT NULL,
    `unit_type` VARCHAR(50) NOT NULL,
    `built_up_area` DECIMAL(10, 2) NOT NULL,
    `carpet_area` DECIMAL(10, 2) NOT NULL,
    `super_built_up_area` DECIMAL(10, 2) NULL,
    `bedrooms` INTEGER NULL,
    `bathrooms` INTEGER NULL,
    `floor` INTEGER NULL,
    `base_price` DECIMAL(15, 2) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'available',
    `description` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `property_unit_tenant_id_idx`(`tenant_id`),
    INDEX `property_unit_project_id_idx`(`project_id`),
    INDEX `property_unit_status_idx`(`status`),
    UNIQUE INDEX `property_unit_project_id_unit_number_key`(`project_id`, `unit_number`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `unit_cost_sheet` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `basic_price` DECIMAL(15, 2) NOT NULL,
    `floor_rise_factor` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `corner_factor` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `other_charges` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `taxes` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `total_price` DECIMAL(15, 2) NOT NULL,
    `valid_from` DATE NOT NULL,
    `valid_till` DATE NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `unit_cost_sheet_tenant_id_idx`(`tenant_id`),
    INDEX `unit_cost_sheet_unit_id_idx`(`unit_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `payment_plan` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `plan_name` VARCHAR(100) NOT NULL,
    `total_price` DECIMAL(15, 2) NOT NULL,
    `upfront_amount` DECIMAL(15, 2) NOT NULL,
    `upfront_percent` DECIMAL(5, 2) NULL,
    `description` TEXT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `payment_plan_tenant_id_idx`(`tenant_id`),
    INDEX `payment_plan_unit_id_idx`(`unit_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `installment` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `payment_plan_id` VARCHAR(36) NOT NULL,
    `installment_number` INTEGER NOT NULL,
    `due_date` DATE NOT NULL,
    `amount` DECIMAL(15, 2) NOT NULL,
    `description` VARCHAR(255) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `installment_tenant_id_idx`(`tenant_id`),
    INDEX `installment_payment_plan_id_idx`(`payment_plan_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `booking` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `customer_id` VARCHAR(36) NULL,
    `booking_date` DATE NOT NULL,
    `booking_amount` DECIMAL(15, 2) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'booked',
    `notes` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `booking_tenant_id_idx`(`tenant_id`),
    INDEX `booking_customer_id_idx`(`customer_id`),
    UNIQUE INDEX `booking_unit_id_key`(`unit_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `property_cost_breakup` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `cost_category` VARCHAR(100) NOT NULL,
    `cost_amount` DECIMAL(15, 2) NOT NULL,
    `cost_percent` DECIMAL(5, 2) NULL,
    `description` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `property_cost_breakup_tenant_id_idx`(`tenant_id`),
    INDEX `property_cost_breakup_unit_id_idx`(`unit_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `property_amenity` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NULL,
    `amenity_name` VARCHAR(100) NOT NULL,
    `amenity_type` VARCHAR(50) NULL,
    `description` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `property_amenity_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `property_unit_area_statement` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `area_type` VARCHAR(100) NOT NULL,
    `area_value` DECIMAL(10, 2) NOT NULL,
    `area_unit` VARCHAR(20) NOT NULL,
    `description` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `property_unit_area_statement_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `project_cost_configuration` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `cost_name` VARCHAR(100) NOT NULL,
    `cost_formula` TEXT NULL,
    `applicable_units` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `project_cost_configuration_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `property_customer_profile` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `customer_id` VARCHAR(36) NOT NULL,
    `occupation_status` VARCHAR(100) NULL,
    `income_range` VARCHAR(50) NULL,
    `investment_goal` VARCHAR(255) NULL,
    `investment_timeline` VARCHAR(50) NULL,
    `preferences` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `property_customer_profile_tenant_id_idx`(`tenant_id`),
    INDEX `property_customer_profile_customer_id_idx`(`customer_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `property_comparable` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `comparable_type` VARCHAR(100) NOT NULL,
    `location` VARCHAR(255) NULL,
    `price_per_sqft` DECIMAL(10, 2) NOT NULL,
    `market_data` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `property_comparable_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `site_visit_feedback` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `site_visit_id` VARCHAR(36) NOT NULL,
    `feedback` TEXT NULL,
    `rating` DECIMAL(3, 1) NULL,
    `follow_up_required` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `site_visit_feedback_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `compliance_checklist` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `checklist_name` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `frequency` VARCHAR(50) NOT NULL,
    `due_date` DATE NOT NULL,
    `completion_date` DATE NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `completed_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `compliance_checklist_tenant_id_idx`(`tenant_id`),
    INDEX `compliance_checklist_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `tax_calculation` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `transaction_type` VARCHAR(100) NOT NULL,
    `transaction_id` VARCHAR(36) NULL,
    `tax_type` VARCHAR(50) NOT NULL,
    `taxable_amount` DECIMAL(15, 2) NOT NULL,
    `tax_rate` DECIMAL(5, 2) NOT NULL,
    `tax_amount` DECIMAL(15, 2) NOT NULL,
    `calculated_date` DATE NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `tax_calculation_tenant_id_idx`(`tenant_id`),
    INDEX `tax_calculation_transaction_type_idx`(`transaction_type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `audit_trail` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `audit_type` VARCHAR(100) NOT NULL,
    `entity_type` VARCHAR(100) NOT NULL,
    `entity_id` VARCHAR(36) NULL,
    `audit_date` DATE NOT NULL,
    `findings` LONGTEXT NULL,
    `status` VARCHAR(50) NOT NULL,
    `auditor` VARCHAR(255) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `audit_trail_tenant_id_idx`(`tenant_id`),
    INDEX `audit_trail_audit_type_idx`(`audit_type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `regulatory_requirement` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `requirement_name` VARCHAR(255) NOT NULL,
    `requirement_type` VARCHAR(100) NOT NULL,
    `jurisdiction` VARCHAR(100) NOT NULL,
    `due_date` DATE NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `owner_department` VARCHAR(100) NULL,
    `description` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `regulatory_requirement_tenant_id_idx`(`tenant_id`),
    INDEX `regulatory_requirement_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `compliance_document` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `compliance_type` VARCHAR(100) NOT NULL,
    `document_name` VARCHAR(255) NOT NULL,
    `document_path` VARCHAR(500) NULL,
    `upload_date` DATE NOT NULL,
    `expiry_date` DATE NULL,
    `uploaded_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `compliance_document_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `compliance_notification` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `notification_type` VARCHAR(100) NOT NULL,
    `subject` VARCHAR(255) NOT NULL,
    `message` LONGTEXT NULL,
    `due_date` DATE NOT NULL,
    `recipient_email` VARCHAR(255) NULL,
    `sent_date` DATETIME(3) NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `compliance_notification_tenant_id_idx`(`tenant_id`),
    INDEX `compliance_notification_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `communication_channel` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `channel_name` VARCHAR(100) NOT NULL,
    `channel_type` VARCHAR(50) NOT NULL,
    `description` TEXT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `config_data` JSON NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `communication_channel_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `communication_channel_tenant_id_channel_name_key`(`tenant_id`, `channel_name`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `message_template` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `template_name` VARCHAR(255) NOT NULL,
    `template_type` VARCHAR(50) NOT NULL,
    `channel_id` VARCHAR(36) NULL,
    `subject` VARCHAR(255) NULL,
    `body` LONGTEXT NULL,
    `variables` TEXT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `message_template_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `communication_session` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `session_id` VARCHAR(100) NOT NULL,
    `initiator_id` VARCHAR(36) NOT NULL,
    `participant_ids` TEXT NULL,
    `channel_type` VARCHAR(50) NOT NULL,
    `subject` VARCHAR(255) NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `start_time` DATETIME(3) NOT NULL,
    `end_time` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `communication_session_session_id_key`(`session_id`),
    INDEX `communication_session_tenant_id_idx`(`tenant_id`),
    INDEX `communication_session_session_id_idx`(`session_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `communication_message` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `session_id` VARCHAR(36) NOT NULL,
    `sender_id` VARCHAR(36) NOT NULL,
    `recipient_ids` TEXT NULL,
    `message_type` VARCHAR(50) NOT NULL,
    `subject` VARCHAR(255) NULL,
    `body` LONGTEXT NULL,
    `attachments` JSON NULL,
    `delivery_status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `read_at` DATETIME(3) NULL,
    `failure_reason` TEXT NULL,
    `sent_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `communication_message_tenant_id_idx`(`tenant_id`),
    INDEX `communication_message_session_id_idx`(`session_id`),
    INDEX `communication_message_sender_id_idx`(`sender_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `email_configuration` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `smtp_server` VARCHAR(255) NOT NULL,
    `smtp_port` INTEGER NOT NULL,
    `from_email` VARCHAR(255) NOT NULL,
    `from_name` VARCHAR(255) NULL,
    `username` VARCHAR(255) NULL,
    `password` VARCHAR(500) NULL,
    `use_ssl` BOOLEAN NOT NULL DEFAULT true,
    `is_active` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `email_configuration_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sms_configuration` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `provider` VARCHAR(100) NOT NULL,
    `api_key` VARCHAR(500) NOT NULL,
    `sender_id` VARCHAR(50) NOT NULL,
    `account_balance` DECIMAL(15, 2) NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `sms_configuration_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `webhook_endpoint` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `endpoint_url` VARCHAR(500) NOT NULL,
    `event_type` VARCHAR(100) NOT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `retry_count` INTEGER NOT NULL DEFAULT 3,
    `timeout` INTEGER NOT NULL DEFAULT 30,
    `last_triggered_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `webhook_endpoint_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `webhook_log` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `webhook_id` VARCHAR(36) NOT NULL,
    `event_type` VARCHAR(100) NOT NULL,
    `payload` JSON NULL,
    `status_code` INTEGER NULL,
    `response` LONGTEXT NULL,
    `attempt_number` INTEGER NOT NULL,
    `next_retry_at` DATETIME(3) NULL,
    `triggered_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `webhook_log_tenant_id_idx`(`tenant_id`),
    INDEX `webhook_log_webhook_id_idx`(`webhook_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `notification_preference` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` CHAR(36) NOT NULL,
    `notification_type` VARCHAR(100) NOT NULL,
    `channel_type` VARCHAR(50) NOT NULL,
    `is_enabled` BOOLEAN NOT NULL DEFAULT true,
    `frequency` VARCHAR(50) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `notification_preference_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `notification_preference_tenant_id_user_id_notification_type__key`(`tenant_id`, `user_id`, `notification_type`, `channel_type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `communication_log` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `communication_type` VARCHAR(50) NOT NULL,
    `from_id` VARCHAR(36) NOT NULL,
    `to_id` VARCHAR(36) NOT NULL,
    `subject` VARCHAR(255) NULL,
    `status` VARCHAR(50) NOT NULL,
    `sent_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `delivered_at` DATETIME(3) NULL,
    `read_at` DATETIME(3) NULL,

    INDEX `communication_log_tenant_id_idx`(`tenant_id`),
    INDEX `communication_log_from_id_to_id_idx`(`from_id`, `to_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `voip_provider` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `provider_name` VARCHAR(100) NOT NULL,
    `api_key` VARCHAR(500) NOT NULL,
    `api_secret` VARCHAR(500) NOT NULL,
    `endpoint` VARCHAR(255) NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `voip_provider_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `click_to_call_session` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `session_id` VARCHAR(100) NOT NULL,
    `initiator_id` VARCHAR(36) NOT NULL,
    `receiver_id` VARCHAR(36) NULL,
    `initiator_phone` VARCHAR(20) NOT NULL,
    `receiver_phone` VARCHAR(20) NULL,
    `status` VARCHAR(50) NOT NULL,
    `duration` INTEGER NOT NULL DEFAULT 0,
    `recording_url` VARCHAR(500) NULL,
    `start_time` DATETIME(3) NOT NULL,
    `end_time` DATETIME(3) NULL,
    `call_type` VARCHAR(50) NULL,
    `notes` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `click_to_call_session_session_id_key`(`session_id`),
    INDEX `click_to_call_session_tenant_id_idx`(`tenant_id`),
    INDEX `click_to_call_session_session_id_idx`(`session_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `call_routing_rule` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `rule_name` VARCHAR(100) NOT NULL,
    `priority` INTEGER NOT NULL,
    `source` VARCHAR(255) NULL,
    `destination` VARCHAR(255) NULL,
    `routing_type` VARCHAR(50) NOT NULL,
    `target_number` VARCHAR(20) NULL,
    `target_group_id` VARCHAR(36) NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `call_routing_rule_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `call_dtmf_interaction` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `session_id` VARCHAR(36) NOT NULL,
    `dtmf_sequence` VARCHAR(50) NOT NULL,
    `action` VARCHAR(100) NOT NULL,
    `timestamp` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `call_dtmf_interaction_tenant_id_idx`(`tenant_id`),
    INDEX `call_dtmf_interaction_session_id_idx`(`session_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ivr_menu` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `menu_name` VARCHAR(100) NOT NULL,
    `prompt` TEXT NOT NULL,
    `recording_url` VARCHAR(500) NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `ivr_menu_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ivr_menu_option` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `menu_id` VARCHAR(36) NOT NULL,
    `option_key` VARCHAR(10) NOT NULL,
    `option_label` VARCHAR(255) NOT NULL,
    `action_type` VARCHAR(50) NOT NULL,
    `action_value` VARCHAR(255) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `ivr_menu_option_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ai_model` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `model_name` VARCHAR(255) NOT NULL,
    `model_version` VARCHAR(50) NOT NULL,
    `model_type` VARCHAR(50) NOT NULL,
    `description` TEXT NULL,
    `accuracy` DECIMAL(5, 2) NULL,
    `training_date` DATE NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `ai_model_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ai_interaction` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `session_id` VARCHAR(36) NOT NULL,
    `model_id` VARCHAR(36) NULL,
    `input` LONGTEXT NULL,
    `output` LONGTEXT NULL,
    `confidence` DECIMAL(5, 2) NULL,
    `processing_time` INTEGER NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `ai_interaction_tenant_id_idx`(`tenant_id`),
    INDEX `ai_interaction_session_id_idx`(`session_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `team_chat_channel` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `channel_name` VARCHAR(100) NOT NULL,
    `channel_type` VARCHAR(50) NOT NULL,
    `is_public` BOOLEAN NOT NULL DEFAULT true,
    `description` TEXT NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `team_chat_channel_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `team_chat_member` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `channel_id` VARCHAR(36) NOT NULL,
    `user_id` CHAR(36) NOT NULL,
    `role` VARCHAR(50) NOT NULL,
    `joined_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `team_chat_member_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `team_chat_member_channel_id_user_id_key`(`channel_id`, `user_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `team_chat_message` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `channel_id` VARCHAR(36) NOT NULL,
    `sender_id` CHAR(36) NOT NULL,
    `message` LONGTEXT NULL,
    `attachments` JSON NULL,
    `edited_at` DATETIME(3) NULL,
    `deleted_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `team_chat_message_tenant_id_idx`(`tenant_id`),
    INDEX `team_chat_message_channel_id_idx`(`channel_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `team_chat_reaction` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `message_id` VARCHAR(36) NOT NULL,
    `user_id` CHAR(36) NOT NULL,
    `emoji` VARCHAR(10) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `team_chat_reaction_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `voice_video_call` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `session_id` VARCHAR(100) NOT NULL,
    `initiator_id` CHAR(36) NOT NULL,
    `participant_ids` TEXT NOT NULL,
    `call_type` VARCHAR(50) NOT NULL,
    `status` VARCHAR(50) NOT NULL,
    `start_time` DATETIME(3) NOT NULL,
    `end_time` DATETIME(3) NULL,
    `duration` INTEGER NOT NULL DEFAULT 0,
    `recording_url` VARCHAR(500) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `voice_video_call_session_id_key`(`session_id`),
    INDEX `voice_video_call_tenant_id_idx`(`tenant_id`),
    INDEX `voice_video_call_session_id_idx`(`session_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `meeting` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `meeting_title` VARCHAR(255) NOT NULL,
    `meeting_start_time` DATETIME(3) NOT NULL,
    `meeting_end_time` DATETIME(3) NULL,
    `organizer_id` CHAR(36) NOT NULL,
    `participant_ids` TEXT NULL,
    `description` TEXT NULL,
    `location` VARCHAR(255) NULL,
    `status` VARCHAR(50) NOT NULL,
    `recording_url` VARCHAR(500) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `meeting_tenant_id_idx`(`tenant_id`),
    INDEX `meeting_organizer_id_idx`(`organizer_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `work_item` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `description` LONGTEXT NULL,
    `assigned_to` CHAR(36) NULL,
    `status` VARCHAR(50) NOT NULL,
    `priority` VARCHAR(50) NOT NULL,
    `due_date` DATE NULL,
    `completion_date` DATETIME(3) NULL,
    `created_by` CHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `work_item_tenant_id_idx`(`tenant_id`),
    INDEX `work_item_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `work_item_comment` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `work_item_id` VARCHAR(36) NOT NULL,
    `user_id` CHAR(36) NOT NULL,
    `comment` LONGTEXT NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `work_item_comment_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `notification` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` CHAR(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `message` LONGTEXT NULL,
    `type` VARCHAR(50) NOT NULL,
    `entity_id` VARCHAR(36) NULL,
    `entity_type` VARCHAR(100) NULL,
    `is_read` BOOLEAN NOT NULL DEFAULT false,
    `read_at` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `notification_tenant_id_idx`(`tenant_id`),
    INDEX `notification_user_id_idx`(`user_id`),
    INDEX `notification_is_read_idx`(`is_read`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `task_assignment` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `work_item_id` VARCHAR(36) NOT NULL,
    `assigned_to` CHAR(36) NOT NULL,
    `assigned_by` CHAR(36) NOT NULL,
    `assignment_date` DATETIME(3) NOT NULL,
    `due_date` DATE NULL,
    `priority` VARCHAR(50) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `task_assignment_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `project_timeline` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `milestone` VARCHAR(255) NOT NULL,
    `target_date` DATE NOT NULL,
    `completion_date` DATE NULL,
    `status` VARCHAR(50) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `project_timeline_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `team_meeting_notes` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `meeting_id` VARCHAR(36) NOT NULL,
    `notes` LONGTEXT NULL,
    `action_items` TEXT NULL,
    `attendees` TEXT NULL,
    `created_by` CHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `team_meeting_notes_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `fixed_asset` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `asset_code` VARCHAR(50) NOT NULL,
    `asset_name` VARCHAR(255) NOT NULL,
    `asset_category` VARCHAR(100) NOT NULL,
    `acquisition_date` DATE NOT NULL,
    `acquisition_cost` DECIMAL(15, 2) NOT NULL,
    `current_value` DECIMAL(15, 2) NOT NULL,
    `location` VARCHAR(255) NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `useful_life` INTEGER NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `fixed_asset_asset_code_key`(`asset_code`),
    INDEX `fixed_asset_tenant_id_idx`(`tenant_id`),
    INDEX `fixed_asset_asset_code_idx`(`asset_code`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `depreciation_schedule` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `asset_id` VARCHAR(36) NOT NULL,
    `period` VARCHAR(50) NOT NULL,
    `opening_value` DECIMAL(15, 2) NOT NULL,
    `depreciation_rate` DECIMAL(5, 2) NOT NULL,
    `depreciation_amount` DECIMAL(15, 2) NOT NULL,
    `closing_value` DECIMAL(15, 2) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `depreciation_schedule_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `asset_revaluation` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `asset_id` VARCHAR(36) NOT NULL,
    `revaluation_date` DATE NOT NULL,
    `previous_value` DECIMAL(15, 2) NOT NULL,
    `new_value` DECIMAL(15, 2) NOT NULL,
    `revaluation_gain_loss` DECIMAL(15, 2) NOT NULL,
    `reason` TEXT NULL,
    `created_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `asset_revaluation_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `asset_disposal` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `asset_id` VARCHAR(36) NOT NULL,
    `disposal_date` DATE NOT NULL,
    `disposal_method` VARCHAR(50) NOT NULL,
    `disposal_price` DECIMAL(15, 2) NOT NULL,
    `book_value` DECIMAL(15, 2) NOT NULL,
    `gain_loss` DECIMAL(15, 2) NOT NULL,
    `approved_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `asset_disposal_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `asset_maintenance` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `asset_id` VARCHAR(36) NOT NULL,
    `maintenance_type` VARCHAR(50) NOT NULL,
    `maintenance_date` DATE NOT NULL,
    `cost` DECIMAL(15, 2) NOT NULL,
    `vendor` VARCHAR(255) NULL,
    `notes` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `asset_maintenance_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `asset_transfer` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `asset_id` VARCHAR(36) NOT NULL,
    `transfer_date` DATE NOT NULL,
    `from_location` VARCHAR(255) NOT NULL,
    `to_location` VARCHAR(255) NOT NULL,
    `reason` TEXT NULL,
    `approved_by` VARCHAR(36) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `asset_transfer_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `cost_center` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `cost_center_code` VARCHAR(50) NOT NULL,
    `cost_center_name` VARCHAR(255) NOT NULL,
    `department` VARCHAR(100) NULL,
    `parent_cost_center_id` VARCHAR(36) NULL,
    `manager` VARCHAR(36) NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `cost_center_cost_center_code_key`(`cost_center_code`),
    INDEX `cost_center_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `cost_allocation` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `cost_center_id` VARCHAR(36) NOT NULL,
    `gl_account_id` VARCHAR(36) NOT NULL,
    `amount` DECIMAL(15, 2) NOT NULL,
    `percentage` DECIMAL(5, 2) NOT NULL,
    `allocation_method` VARCHAR(50) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `cost_allocation_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `budget` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `budget_code` VARCHAR(50) NOT NULL,
    `budget_name` VARCHAR(255) NOT NULL,
    `cost_center_id` VARCHAR(36) NOT NULL,
    `fiscal_year` INTEGER NOT NULL,
    `budget_period` VARCHAR(20) NOT NULL,
    `budgeted_amount` DECIMAL(15, 2) NOT NULL,
    `actual_amount` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `committed_amount` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `available_amount` DECIMAL(15, 2) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `budget_budget_code_key`(`budget_code`),
    INDEX `budget_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `budget_variance` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `budget_id` VARCHAR(36) NOT NULL,
    `variance_period` VARCHAR(20) NOT NULL,
    `budgeted_amount` DECIMAL(15, 2) NOT NULL,
    `actual_amount` DECIMAL(15, 2) NOT NULL,
    `variance_amount` DECIMAL(15, 2) NOT NULL,
    `variance_percentage` DECIMAL(5, 2) NOT NULL,
    `variance_type` VARCHAR(20) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `budget_variance_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `budget_allocation` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `budget_id` VARCHAR(36) NOT NULL,
    `gl_account_id` VARCHAR(36) NOT NULL,
    `allocated_amount` DECIMAL(15, 2) NOT NULL,
    `spent_amount` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `remaining_amount` DECIMAL(15, 2) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `budget_allocation_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `budget_approval` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `budget_id` VARCHAR(36) NOT NULL,
    `approved_by` VARCHAR(36) NOT NULL,
    `approval_date` DATE NOT NULL,
    `comments` TEXT NULL,
    `status` VARCHAR(50) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `budget_approval_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `dashboard` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `dashboard_name` VARCHAR(255) NOT NULL,
    `dashboard_type` VARCHAR(100) NOT NULL,
    `created_by` VARCHAR(36) NOT NULL,
    `is_public` BOOLEAN NOT NULL DEFAULT false,
    `refresh_interval` INTEGER NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `dashboard_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `dashboard_widget` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `dashboard_id` VARCHAR(36) NOT NULL,
    `widget_name` VARCHAR(255) NOT NULL,
    `widget_type` VARCHAR(100) NOT NULL,
    `position` INTEGER NOT NULL,
    `size` VARCHAR(50) NOT NULL,
    `configuration` LONGTEXT NULL,
    `refresh_interval` INTEGER NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `dashboard_widget_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `kpi_definition` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `kpi_code` VARCHAR(100) NOT NULL,
    `kpi_name` VARCHAR(255) NOT NULL,
    `kpi_category` VARCHAR(100) NOT NULL,
    `description` TEXT NULL,
    `formula` LONGTEXT NOT NULL,
    `target_value` DECIMAL(15, 4) NULL,
    `min_value` DECIMAL(15, 4) NULL,
    `max_value` DECIMAL(15, 4) NULL,
    `unit` VARCHAR(50) NOT NULL,
    `frequency` VARCHAR(50) NOT NULL,
    `owner` VARCHAR(36) NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `kpi_definition_kpi_code_key`(`kpi_code`),
    INDEX `kpi_definition_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `kpi_value` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `kpi_id` VARCHAR(36) NOT NULL,
    `value_date` DATE NOT NULL,
    `actual_value` DECIMAL(15, 4) NOT NULL,
    `target_value` DECIMAL(15, 4) NULL,
    `variance` DECIMAL(15, 4) NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'normal',
    `trend` VARCHAR(50) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `kpi_value_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `report` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `report_code` VARCHAR(100) NOT NULL,
    `report_name` VARCHAR(255) NOT NULL,
    `report_type` VARCHAR(100) NOT NULL,
    `description` TEXT NULL,
    `query` LONGTEXT NOT NULL,
    `schedule` VARCHAR(100) NULL,
    `last_run_date` DATETIME(3) NULL,
    `next_run_date` DATETIME(3) NULL,
    `owner` VARCHAR(36) NULL,
    `is_published` BOOLEAN NOT NULL DEFAULT false,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `report_report_code_key`(`report_code`),
    INDEX `report_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `report_execution` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `report_id` VARCHAR(36) NOT NULL,
    `executed_by` VARCHAR(36) NOT NULL,
    `execution_date` DATETIME NOT NULL,
    `execution_time` INTEGER NOT NULL,
    `record_count` INTEGER NOT NULL DEFAULT 0,
    `output_file` VARCHAR(255) NULL,
    `status` VARCHAR(50) NOT NULL,
    `error_message` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `report_execution_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `analytics_data_cache` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `entity_type` VARCHAR(100) NOT NULL,
    `period` VARCHAR(50) NOT NULL,
    `dimension_values` LONGTEXT NULL,
    `metric_name` VARCHAR(255) NOT NULL,
    `metric_value` DECIMAL(15, 4) NOT NULL,
    `last_updated` DATETIME NOT NULL,
    `expiry_date` DATE NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `analytics_data_cache_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `custom_report` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `report_name` VARCHAR(255) NOT NULL,
    `created_by` VARCHAR(36) NOT NULL,
    `columns` LONGTEXT NOT NULL,
    `filters` LONGTEXT NULL,
    `group_by` LONGTEXT NULL,
    `order_by` LONGTEXT NULL,
    `format` VARCHAR(50) NOT NULL DEFAULT 'pdf',
    `is_public` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `custom_report_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `report_schedule` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `report_id` VARCHAR(36) NOT NULL,
    `frequency` VARCHAR(50) NOT NULL,
    `day_of_week` INTEGER NULL,
    `day_of_month` INTEGER NULL,
    `hour` INTEGER NOT NULL,
    `minute` INTEGER NOT NULL,
    `email_recipients` TEXT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `next_run_date` DATETIME NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `report_schedule_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `analytics_event` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `event_type` VARCHAR(100) NOT NULL,
    `entity_type` VARCHAR(100) NOT NULL,
    `entity_id` VARCHAR(36) NULL,
    `user_id` VARCHAR(36) NULL,
    `event_data` LONGTEXT NULL,
    `timestamp` DATETIME NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `analytics_event_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `mobile_device` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` VARCHAR(36) NOT NULL,
    `device_id` VARCHAR(255) NOT NULL,
    `device_type` VARCHAR(50) NOT NULL,
    `os_type` VARCHAR(50) NOT NULL,
    `os_version` VARCHAR(50) NULL,
    `app_version` VARCHAR(50) NULL,
    `push_token` TEXT NULL,
    `device_model` VARCHAR(100) NULL,
    `manufacturer` VARCHAR(100) NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `last_active_at` DATETIME NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `mobile_device_device_id_key`(`device_id`),
    INDEX `mobile_device_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `mobile_session` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `device_id` VARCHAR(36) NOT NULL,
    `user_id` VARCHAR(36) NOT NULL,
    `session_token` VARCHAR(500) NOT NULL,
    `login_time` DATETIME NOT NULL,
    `logout_time` DATETIME NULL,
    `ip_address` VARCHAR(45) NULL,
    `user_agent` TEXT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `idx_mobile_session_token`(`session_token`),
    INDEX `mobile_session_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `mobile_push_notification` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `device_id` VARCHAR(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `body` TEXT NOT NULL,
    `notification_type` VARCHAR(100) NOT NULL,
    `entity_type` VARCHAR(100) NULL,
    `entity_id` VARCHAR(36) NULL,
    `payload` LONGTEXT NULL,
    `sent_at` DATETIME NOT NULL,
    `delivered_at` DATETIME NULL,
    `opened_at` DATETIME NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `mobile_push_notification_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `mobile_app_update` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `version_number` VARCHAR(50) NOT NULL,
    `release_notes` LONGTEXT NULL,
    `download_url` TEXT NOT NULL,
    `minimum_os_version` VARCHAR(50) NULL,
    `force_update` BOOLEAN NOT NULL DEFAULT false,
    `release_date` DATE NOT NULL,
    `status` VARCHAR(50) NOT NULL,
    `download_count` INTEGER NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `mobile_app_update_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `mobile_offline_data` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `device_id` VARCHAR(36) NOT NULL,
    `entity_type` VARCHAR(100) NOT NULL,
    `entity_id` VARCHAR(36) NOT NULL,
    `data_snapshot` LONGTEXT NOT NULL,
    `last_synced_at` DATETIME NOT NULL,
    `sync_required` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `mobile_offline_data_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `mobile_app_config` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `config_key` VARCHAR(255) NOT NULL,
    `config_value` LONGTEXT NOT NULL,
    `config_type` VARCHAR(100) NOT NULL,
    `description` TEXT NULL,
    `updated_by` VARCHAR(36) NULL,
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `mobile_app_config_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `mobile_analytics` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `device_id` VARCHAR(36) NOT NULL,
    `screen_name` VARCHAR(255) NOT NULL,
    `event_name` VARCHAR(255) NOT NULL,
    `event_value` LONGTEXT NULL,
    `session_duration` INTEGER NULL,
    `crash_reported` BOOLEAN NOT NULL DEFAULT false,
    `timestamp` DATETIME NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `mobile_analytics_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `construction_project` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_code` VARCHAR(50) NOT NULL,
    `project_name` VARCHAR(255) NOT NULL,
    `location` TEXT NOT NULL,
    `project_manager` VARCHAR(36) NULL,
    `start_date` DATE NOT NULL,
    `completion_date` DATE NULL,
    `estimated_budget` DECIMAL(15, 2) NOT NULL,
    `actual_budget` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'planning',
    `percentage_complete` DECIMAL(5, 2) NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `construction_project_project_code_key`(`project_code`),
    INDEX `construction_project_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `construction_phase` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `phase_name` VARCHAR(255) NOT NULL,
    `description` TEXT NULL,
    `start_date` DATE NOT NULL,
    `expected_end_date` DATE NOT NULL,
    `actual_end_date` DATE NULL,
    `budgeted_amount` DECIMAL(15, 2) NOT NULL,
    `actual_amount` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'not_started',
    `sequence` INTEGER NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `construction_phase_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `construction_material` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `material_code` VARCHAR(100) NOT NULL,
    `material_name` VARCHAR(255) NOT NULL,
    `quantity` DECIMAL(18, 4) NOT NULL,
    `unit` VARCHAR(50) NOT NULL,
    `unit_cost` DECIMAL(15, 2) NOT NULL,
    `total_cost` DECIMAL(15, 2) NOT NULL,
    `supplier` VARCHAR(255) NULL,
    `delivery_date` DATE NOT NULL,
    `received_quantity` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `wastage_quantity` DECIMAL(18, 4) NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `construction_material_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `construction_labor` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `labor_date` DATE NOT NULL,
    `labor_type` VARCHAR(100) NOT NULL,
    `workers` INTEGER NOT NULL,
    `hours_per_worker` DECIMAL(5, 2) NOT NULL,
    `rate_per_hour` DECIMAL(10, 2) NOT NULL,
    `total_cost` DECIMAL(15, 2) NOT NULL,
    `supervisor` VARCHAR(36) NULL,
    `notes` TEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `construction_labor_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `construction_equipment` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `equipment_code` VARCHAR(100) NOT NULL,
    `equipment_name` VARCHAR(255) NOT NULL,
    `equipment_type` VARCHAR(100) NOT NULL,
    `allocation_date` DATE NOT NULL,
    `release_date` DATE NULL,
    `daily_cost` DECIMAL(15, 2) NOT NULL,
    `total_cost` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'allocated',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `construction_equipment_equipment_code_key`(`equipment_code`),
    INDEX `construction_equipment_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `construction_quality_check` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `check_date` DATE NOT NULL,
    `check_type` VARCHAR(100) NOT NULL,
    `checked_by` VARCHAR(36) NOT NULL,
    `area` VARCHAR(255) NOT NULL,
    `findings` LONGTEXT NULL,
    `status` VARCHAR(50) NOT NULL,
    `passed_criteria` INTEGER NOT NULL,
    `total_criteria` INTEGER NOT NULL,
    `next_check_date` DATE NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `construction_quality_check_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `construction_safety` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `incident_date` DATE NOT NULL,
    `incident_type` VARCHAR(100) NOT NULL,
    `severity` VARCHAR(50) NOT NULL,
    `description` LONGTEXT NOT NULL,
    `reported_by` VARCHAR(36) NOT NULL,
    `affected_workers` INTEGER NULL,
    `lost_days` INTEGER NULL,
    `investigation_required` BOOLEAN NOT NULL DEFAULT false,
    `status` VARCHAR(50) NOT NULL DEFAULT 'reported',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `construction_safety_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `site_visit` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_id` VARCHAR(36) NOT NULL,
    `visit_date` DATE NOT NULL,
    `visited_by` VARCHAR(36) NOT NULL,
    `purpose` TEXT NOT NULL,
    `observations` LONGTEXT NULL,
    `photographs` LONGTEXT NULL,
    `actions_required` LONGTEXT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'completed',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `site_visit_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `custom_field_definition` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `field_name` VARCHAR(255) NOT NULL,
    `field_code` VARCHAR(100) NOT NULL,
    `entity_type` VARCHAR(100) NOT NULL,
    `field_type` VARCHAR(50) NOT NULL,
    `field_length` INTEGER NULL,
    `is_required` BOOLEAN NOT NULL DEFAULT false,
    `is_unique` BOOLEAN NOT NULL DEFAULT false,
    `default_value` TEXT NULL,
    `validation_rule` TEXT NULL,
    `display_order` INTEGER NOT NULL,
    `help_text` TEXT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `custom_field_definition_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `custom_field_value` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `field_definition_id` VARCHAR(36) NOT NULL,
    `entity_type` VARCHAR(100) NOT NULL,
    `entity_id` VARCHAR(36) NOT NULL,
    `field_value` LONGTEXT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `custom_field_value_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `multi_language_resource` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `resource_key` VARCHAR(255) NOT NULL,
    `language` VARCHAR(10) NOT NULL,
    `resource_value` LONGTEXT NOT NULL,
    `resource_context` VARCHAR(100) NULL,
    `last_updated` DATETIME NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `multi_language_resource_tenant_id_idx`(`tenant_id`),
    UNIQUE INDEX `multi_language_resource_resource_key_language_tenant_id_key`(`resource_key`, `language`, `tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `data_import_template` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `template_name` VARCHAR(255) NOT NULL,
    `template_code` VARCHAR(100) NOT NULL,
    `entity_type` VARCHAR(100) NOT NULL,
    `file_format` VARCHAR(50) NOT NULL,
    `column_mapping` LONGTEXT NOT NULL,
    `validation_rules` LONGTEXT NULL,
    `transformations` LONGTEXT NULL,
    `created_by` VARCHAR(36) NOT NULL,
    `is_published` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `data_import_template_template_code_key`(`template_code`),
    INDEX `data_import_template_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `data_import_job` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `template_id` VARCHAR(36) NOT NULL,
    `file_name` VARCHAR(255) NOT NULL,
    `file_path` TEXT NOT NULL,
    `uploaded_by` VARCHAR(36) NOT NULL,
    `total_records` INTEGER NOT NULL DEFAULT 0,
    `success_count` INTEGER NOT NULL DEFAULT 0,
    `error_count` INTEGER NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `error_log` LONGTEXT NULL,
    `start_time` DATETIME NULL,
    `completion_time` DATETIME NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `data_import_job_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `data_export_template` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `template_name` VARCHAR(255) NOT NULL,
    `template_code` VARCHAR(100) NOT NULL,
    `entity_type` VARCHAR(100) NOT NULL,
    `columns` LONGTEXT NOT NULL,
    `filters` LONGTEXT NULL,
    `sort_order` LONGTEXT NULL,
    `file_format` VARCHAR(50) NOT NULL,
    `created_by` VARCHAR(36) NOT NULL,
    `is_published` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `data_export_template_template_code_key`(`template_code`),
    INDEX `data_export_template_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `data_export_job` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `template_id` VARCHAR(36) NOT NULL,
    `requested_by` VARCHAR(36) NOT NULL,
    `total_records` INTEGER NOT NULL DEFAULT 0,
    `exported_records` INTEGER NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `file_location` TEXT NULL,
    `download_url` TEXT NULL,
    `expiry_date` DATE NULL,
    `start_time` DATETIME NULL,
    `completion_time` DATETIME NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `data_export_job_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `plugin` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `plugin_name` VARCHAR(255) NOT NULL,
    `plugin_code` VARCHAR(100) NOT NULL,
    `plugin_version` VARCHAR(50) NOT NULL,
    `author` VARCHAR(255) NULL,
    `description` LONGTEXT NULL,
    `plugin_url` TEXT NULL,
    `settings` LONGTEXT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT false,
    `install_date` DATE NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `plugin_plugin_code_key`(`plugin_code`),
    INDEX `plugin_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `plugin_event` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `plugin_id` VARCHAR(36) NOT NULL,
    `event_type` VARCHAR(100) NOT NULL,
    `event_data` LONGTEXT NULL,
    `event_status` VARCHAR(50) NOT NULL,
    `execution_time` INTEGER NULL,
    `error_message` TEXT NULL,
    `timestamp` DATETIME NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `plugin_event_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `service_integration` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `integration_name` VARCHAR(255) NOT NULL,
    `integration_code` VARCHAR(100) NOT NULL,
    `service_type` VARCHAR(100) NOT NULL,
    `api_endpoint` TEXT NOT NULL,
    `api_key` TEXT NULL,
    `auth_type` VARCHAR(50) NOT NULL,
    `settings` LONGTEXT NULL,
    `is_active` BOOLEAN NOT NULL DEFAULT false,
    `last_sync_date` DATETIME NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    UNIQUE INDEX `service_integration_integration_code_key`(`integration_code`),
    INDEX `service_integration_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `integration_log` (
    `id` CHAR(36) NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `integration_id` VARCHAR(36) NOT NULL,
    `action` VARCHAR(100) NOT NULL,
    `request_data` LONGTEXT NULL,
    `response_data` LONGTEXT NULL,
    `status_code` INTEGER NULL,
    `error_message` TEXT NULL,
    `execution_time` INTEGER NULL,
    `timestamp` DATETIME NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `integration_log_tenant_id_idx`(`tenant_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ImportShipment` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `shipmentNumber` VARCHAR(191) NOT NULL,
    `containerNumber` VARCHAR(191) NULL,
    `vesselName` VARCHAR(191) NULL,
    `voyageNumber` VARCHAR(191) NULL,
    `exporterName` VARCHAR(191) NOT NULL,
    `exporterCountry` VARCHAR(191) NOT NULL,
    `importerName` VARCHAR(191) NOT NULL,
    `importerIEC` VARCHAR(191) NOT NULL,
    `importLicenseId` VARCHAR(191) NULL,
    `incoterm` VARCHAR(191) NOT NULL,
    `cifValue` DECIMAL(15, 2) NOT NULL,
    `fobValue` DECIMAL(15, 2) NOT NULL,
    `insuranceValue` DECIMAL(12, 2) NOT NULL,
    `freight` DECIMAL(12, 2) NOT NULL,
    `quantity` DECIMAL(12, 3) NOT NULL,
    `unit` VARCHAR(191) NOT NULL,
    `status` VARCHAR(191) NOT NULL,
    `estimatedArrival` DATETIME(3) NULL,
    `actualArrival` DATETIME(3) NULL,
    `clearanceDate` DATETIME(3) NULL,
    `warehouseId` VARCHAR(191) NULL,
    `vendorInvoiceId` VARCHAR(191) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,
    `userId` VARCHAR(191) NOT NULL,

    UNIQUE INDEX `ImportShipment_shipmentNumber_key`(`shipmentNumber`),
    INDEX `ImportShipment_tenantId_idx`(`tenantId`),
    INDEX `ImportShipment_status_idx`(`status`),
    INDEX `ImportShipment_estimatedArrival_idx`(`estimatedArrival`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ExportShipment` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `shipmentNumber` VARCHAR(191) NOT NULL,
    `containerNumber` VARCHAR(191) NULL,
    `vesselName` VARCHAR(191) NULL,
    `voyageNumber` VARCHAR(191) NULL,
    `exporterName` VARCHAR(191) NOT NULL,
    `exporterIEC` VARCHAR(191) NOT NULL,
    `importerName` VARCHAR(191) NOT NULL,
    `importerCountry` VARCHAR(191) NOT NULL,
    `salesOrderId` VARCHAR(191) NULL,
    `invoiceId` VARCHAR(191) NULL,
    `incoterm` VARCHAR(191) NOT NULL,
    `fobValue` DECIMAL(15, 2) NOT NULL,
    `quantity` DECIMAL(12, 3) NOT NULL,
    `unit` VARCHAR(191) NOT NULL,
    `status` VARCHAR(191) NOT NULL,
    `shipmentDate` DATETIME(3) NOT NULL,
    `deliveryDate` DATETIME(3) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,
    `userId` VARCHAR(191) NOT NULL,

    UNIQUE INDEX `ExportShipment_shipmentNumber_key`(`shipmentNumber`),
    INDEX `ExportShipment_tenantId_idx`(`tenantId`),
    INDEX `ExportShipment_status_idx`(`status`),
    INDEX `ExportShipment_shipmentDate_idx`(`shipmentDate`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ImportHSCodeLineItem` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `importShipmentId` VARCHAR(191) NOT NULL,
    `hsCode` CHAR(8) NOT NULL,
    `description` VARCHAR(191) NOT NULL,
    `quantity` DECIMAL(12, 3) NOT NULL,
    `unit` VARCHAR(191) NOT NULL,
    `unitPrice` DECIMAL(12, 4) NOT NULL,
    `cifValue` DECIMAL(15, 2) NOT NULL,
    `gstRate` DECIMAL(5, 2) NOT NULL,
    `gstAmount` DECIMAL(12, 2) NOT NULL,
    `stockLevelId` VARCHAR(191) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `ImportHSCodeLineItem_importShipmentId_idx`(`importShipmentId`),
    INDEX `ImportHSCodeLineItem_hsCode_idx`(`hsCode`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ExportHSCodeLineItem` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `exportShipmentId` VARCHAR(191) NOT NULL,
    `hsCode` CHAR(8) NOT NULL,
    `description` VARCHAR(191) NOT NULL,
    `quantity` DECIMAL(12, 3) NOT NULL,
    `unit` VARCHAR(191) NOT NULL,
    `unitPrice` DECIMAL(12, 4) NOT NULL,
    `fobValue` DECIMAL(15, 2) NOT NULL,
    `soLineItemId` VARCHAR(191) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `ExportHSCodeLineItem_exportShipmentId_idx`(`exportShipmentId`),
    INDEX `ExportHSCodeLineItem_hsCode_idx`(`hsCode`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `DutyCalculation` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `importShipmentId` VARCHAR(191) NOT NULL,
    `lineItemId` VARCHAR(191) NOT NULL,
    `basicCustomsDutyRate` DECIMAL(5, 2) NOT NULL,
    `basicCustomsDuty` DECIMAL(12, 2) NOT NULL,
    `additionalDutyRate` DECIMAL(5, 2) NOT NULL,
    `additionalDuty` DECIMAL(12, 2) NOT NULL,
    `safeguardDutyRate` DECIMAL(5, 2) NOT NULL,
    `safeguardDuty` DECIMAL(12, 2) NOT NULL,
    `antiDumpingDutyRate` DECIMAL(5, 2) NOT NULL,
    `antiDumpingDuty` DECIMAL(12, 2) NOT NULL,
    `totalCustomsDuty` DECIMAL(15, 2) NOT NULL,
    `iGSTRate` DECIMAL(5, 2) NOT NULL,
    `iGST` DECIMAL(12, 2) NOT NULL,
    `tradeAgreementId` VARCHAR(191) NULL,
    `preferentialDutyRate` DECIMAL(5, 2) NOT NULL,
    `preferentialDuty` DECIMAL(12, 2) NOT NULL,
    `totalPayable` DECIMAL(15, 2) NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    UNIQUE INDEX `DutyCalculation_lineItemId_key`(`lineItemId`),
    INDEX `DutyCalculation_importShipmentId_idx`(`importShipmentId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `CustomsDocumentation` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `importShipmentId` VARCHAR(191) NULL,
    `exportShipmentId` VARCHAR(191) NULL,
    `documentType` VARCHAR(191) NOT NULL,
    `documentNumber` VARCHAR(191) NOT NULL,
    `referenceNumber` VARCHAR(191) NOT NULL,
    `status` VARCHAR(191) NOT NULL,
    `generatedDate` DATETIME(3) NOT NULL,
    `filedDate` DATETIME(3) NULL,
    `approvedDate` DATETIME(3) NULL,
    `rejectionReason` VARCHAR(191) NULL,
    `documentFile` LONGBLOB NULL,
    `documentPath` VARCHAR(191) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,
    `userId` VARCHAR(191) NOT NULL,
    `approverUserId` VARCHAR(191) NULL,

    INDEX `CustomsDocumentation_tenantId_idx`(`tenantId`),
    INDEX `CustomsDocumentation_documentType_idx`(`documentType`),
    INDEX `CustomsDocumentation_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ShipmentTracking` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `importShipmentId` VARCHAR(191) NULL,
    `exportShipmentId` VARCHAR(191) NULL,
    `location` VARCHAR(191) NOT NULL,
    `status` VARCHAR(191) NOT NULL,
    `timestamp` DATETIME(3) NOT NULL,
    `notes` VARCHAR(191) NULL,
    `gpsCoordinates` VARCHAR(191) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `ShipmentTracking_tenantId_idx`(`tenantId`),
    INDEX `ShipmentTracking_importShipmentId_idx`(`importShipmentId`),
    INDEX `ShipmentTracking_exportShipmentId_idx`(`exportShipmentId`),
    INDEX `ShipmentTracking_timestamp_idx`(`timestamp`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ImportLicense` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `licenseNumber` VARCHAR(191) NOT NULL,
    `licenseType` VARCHAR(191) NOT NULL,
    `tradingCategory` VARCHAR(191) NOT NULL,
    `validityFrom` DATETIME(3) NOT NULL,
    `validityTo` DATETIME(3) NOT NULL,
    `status` VARCHAR(191) NOT NULL,
    `issuedBy` VARCHAR(191) NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    UNIQUE INDEX `ImportLicense_licenseNumber_key`(`licenseNumber`),
    INDEX `ImportLicense_tenantId_idx`(`tenantId`),
    INDEX `ImportLicense_licenseNumber_idx`(`licenseNumber`),
    INDEX `ImportLicense_status_idx`(`status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ComplianceCheckImEx` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `importShipmentId` VARCHAR(191) NULL,
    `exportShipmentId` VARCHAR(191) NULL,
    `gstIECMatch` BOOLEAN NOT NULL DEFAULT false,
    `gstEligibleForITC` BOOLEAN NOT NULL DEFAULT false,
    `gstRate` DECIMAL(5, 2) NOT NULL,
    `gstCompliance` VARCHAR(191) NOT NULL,
    `dgftRestricted` BOOLEAN NOT NULL DEFAULT false,
    `dgftLicenseRequired` BOOLEAN NOT NULL DEFAULT false,
    `dgftLicenseNumber` VARCHAR(191) NULL,
    `dgftStatus` VARCHAR(191) NOT NULL,
    `rbiEligible` BOOLEAN NOT NULL DEFAULT false,
    `femaCompliant` BOOLEAN NOT NULL DEFAULT false,
    `advanceRemittanceAllowed` BOOLEAN NOT NULL DEFAULT false,
    `femaStatus` VARCHAR(191) NOT NULL,
    `overallStatus` VARCHAR(191) NOT NULL,
    `issues` VARCHAR(191) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,
    `userId` VARCHAR(191) NOT NULL,

    UNIQUE INDEX `ComplianceCheckImEx_importShipmentId_key`(`importShipmentId`),
    UNIQUE INDEX `ComplianceCheckImEx_exportShipmentId_key`(`exportShipmentId`),
    INDEX `ComplianceCheckImEx_tenantId_idx`(`tenantId`),
    INDEX `ComplianceCheckImEx_overallStatus_idx`(`overallStatus`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ImportExportPartner` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `partnerType` VARCHAR(191) NOT NULL,
    `companyName` VARCHAR(191) NOT NULL,
    `licenseNumber` VARCHAR(191) NOT NULL,
    `licenseCategory` VARCHAR(191) NOT NULL,
    `validUntil` DATETIME(3) NOT NULL,
    `portsOfOperation` JSON NOT NULL,
    `countriesCovered` JSON NOT NULL,
    `contactName` VARCHAR(191) NOT NULL,
    `email` VARCHAR(191) NOT NULL,
    `phone` VARCHAR(191) NOT NULL,
    `website` VARCHAR(191) NULL,
    `reliabilityRating` DECIMAL(3, 1) NOT NULL DEFAULT 0,
    `turnaroundTime` INTEGER NOT NULL DEFAULT 0,
    `isActive` BOOLEAN NOT NULL DEFAULT true,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `ImportExportPartner_tenantId_idx`(`tenantId`),
    INDEX `ImportExportPartner_partnerType_idx`(`partnerType`),
    INDEX `ImportExportPartner_isActive_idx`(`isActive`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `HSCodeMaster` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `hsCode` CHAR(8) NOT NULL,
    `description` VARCHAR(191) NOT NULL,
    `longDescription` VARCHAR(191) NULL,
    `gstRate` DECIMAL(5, 2) NOT NULL,
    `gstCategory` VARCHAR(191) NOT NULL,
    `basicCustomsDutyRate` DECIMAL(5, 2) NOT NULL,
    `additionalDutyRate` DECIMAL(5, 2) NOT NULL,
    `isExempted` BOOLEAN NOT NULL DEFAULT false,
    `isSafeguardApplicable` BOOLEAN NOT NULL DEFAULT false,
    `isAntiDumpingApplicable` BOOLEAN NOT NULL DEFAULT false,
    `isRestricted` BOOLEAN NOT NULL DEFAULT false,
    `requiresLicense` BOOLEAN NOT NULL DEFAULT false,
    `restrictions` VARCHAR(191) NULL,
    `notes` VARCHAR(191) NULL,
    `lastUpdatedFrom` VARCHAR(191) NOT NULL,
    `lastUpdatedDate` DATETIME(3) NOT NULL,
    `validFrom` DATETIME(3) NOT NULL,
    `validTo` DATETIME(3) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    UNIQUE INDEX `HSCodeMaster_hsCode_key`(`hsCode`),
    INDEX `HSCodeMaster_tenantId_idx`(`tenantId`),
    INDEX `HSCodeMaster_hsCode_idx`(`hsCode`),
    INDEX `HSCodeMaster_gstRate_idx`(`gstRate`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `TradeAgreementPreference` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `agreementName` VARCHAR(191) NOT NULL,
    `agreementCode` VARCHAR(191) NOT NULL,
    `hsCode` CHAR(8) NOT NULL,
    `originCountry` VARCHAR(191) NOT NULL,
    `standardDutyRate` DECIMAL(5, 2) NOT NULL,
    `preferentialDutyRate` DECIMAL(5, 2) NOT NULL,
    `dutyReduction` DECIMAL(5, 2) NOT NULL,
    `requiresDocuments` JSON NOT NULL,
    `documentType` VARCHAR(191) NOT NULL,
    `validFrom` DATETIME(3) NOT NULL,
    `validTo` DATETIME(3) NULL,
    `isActive` BOOLEAN NOT NULL DEFAULT true,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `TradeAgreementPreference_tenantId_idx`(`tenantId`),
    INDEX `TradeAgreementPreference_agreementName_idx`(`agreementName`),
    INDEX `TradeAgreementPreference_hsCode_idx`(`hsCode`),
    INDEX `TradeAgreementPreference_originCountry_idx`(`originCountry`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `tenant_credentials` (
    `id` VARCHAR(191) NOT NULL,
    `tenantId` VARCHAR(191) NOT NULL,
    `credentialType` VARCHAR(50) NOT NULL,
    `encryptedValue` TEXT NOT NULL,
    `description` VARCHAR(255) NULL,
    `isActive` BOOLEAN NOT NULL DEFAULT true,
    `lastRotatedAt` DATETIME(3) NULL,
    `expiresAt` DATETIME(3) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,
    `createdBy` VARCHAR(191) NULL,
    `updatedBy` VARCHAR(191) NULL,

    INDEX `tenant_credentials_tenantId_idx`(`tenantId`),
    INDEX `tenant_credentials_credentialType_idx`(`credentialType`),
    INDEX `tenant_credentials_isActive_idx`(`isActive`),
    UNIQUE INDEX `tenant_credentials_tenantId_credentialType_key`(`tenantId`, `credentialType`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateIndex
CREATE INDEX `audit_log_tenant_id_idx` ON `audit_log`(`tenant_id`);

-- CreateIndex
CREATE INDEX `audit_log_user_id_idx` ON `audit_log`(`user_id`);

-- CreateIndex
CREATE INDEX `audit_log_entity_type_entity_id_idx` ON `audit_log`(`entity_type`, `entity_id`);

-- CreateIndex
CREATE INDEX `audit_log_created_at_idx` ON `audit_log`(`created_at`);

-- CreateIndex
CREATE UNIQUE INDEX `auth_token_token_key` ON `auth_token`(`token`);

-- CreateIndex
CREATE INDEX `auth_token_user_id_idx` ON `auth_token`(`user_id`);

-- CreateIndex
CREATE INDEX `permission_tenant_id_idx` ON `permission`(`tenant_id`);

-- CreateIndex
CREATE INDEX `permission_resource_idx` ON `permission`(`resource`);

-- CreateIndex
CREATE UNIQUE INDEX `permission_tenant_id_permission_name_key` ON `permission`(`tenant_id`, `permission_name`);

-- CreateIndex
CREATE UNIQUE INDEX `quotation_quotation_number_key` ON `quotation`(`quotation_number`);

-- CreateIndex
CREATE INDEX `quotation_tenant_id_idx` ON `quotation`(`tenant_id`);

-- CreateIndex
CREATE INDEX `quotation_quotation_number_idx` ON `quotation`(`quotation_number`);

-- CreateIndex
CREATE INDEX `quotation_status_idx` ON `quotation`(`status`);

-- CreateIndex
CREATE INDEX `role_tenant_id_idx` ON `role`(`tenant_id`);

-- CreateIndex
CREATE UNIQUE INDEX `role_tenant_id_role_name_key` ON `role`(`tenant_id`, `role_name`);

-- CreateIndex
CREATE UNIQUE INDEX `sales_lead_lead_code_key` ON `sales_lead`(`lead_code`);

-- CreateIndex
CREATE INDEX `sales_lead_tenant_id_idx` ON `sales_lead`(`tenant_id`);

-- CreateIndex
CREATE INDEX `sales_lead_lead_code_idx` ON `sales_lead`(`lead_code`);

-- CreateIndex
CREATE INDEX `sales_lead_lead_status_idx` ON `sales_lead`(`lead_status`);

-- CreateIndex
CREATE INDEX `tenant_status_idx` ON `tenant`(`status`);

-- CreateIndex
CREATE UNIQUE INDEX `user_email_key` ON `user`(`email`);

-- CreateIndex
CREATE INDEX `user_email_idx` ON `user`(`email`);

-- CreateIndex
CREATE INDEX `user_tenant_id_idx` ON `user`(`tenant_id`);

-- CreateIndex
CREATE INDEX `user_role_idx` ON `user`(`role`);

-- CreateIndex
CREATE INDEX `user_role_tenant_id_idx` ON `user_role`(`tenant_id`);

-- CreateIndex
CREATE INDEX `user_role_user_id_idx` ON `user_role`(`user_id`);

-- CreateIndex
CREATE INDEX `user_role_role_id_idx` ON `user_role`(`role_id`);

-- CreateIndex
CREATE UNIQUE INDEX `user_role_user_id_role_id_key` ON `user_role`(`user_id`, `role_id`);

-- AddForeignKey
ALTER TABLE `user` ADD CONSTRAINT `user_tenant_id_fkey` FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `user` ADD CONSTRAINT `user_current_tenant_id_fkey` FOREIGN KEY (`current_tenant_id`) REFERENCES `tenant`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `password_reset_token` ADD CONSTRAINT `password_reset_token_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `auth_token` ADD CONSTRAINT `auth_token_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `team` ADD CONSTRAINT `team_tenant_id_fkey` FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `system_config` ADD CONSTRAINT `system_config_tenant_id_fkey` FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `audit_log` ADD CONSTRAINT `audit_log_tenant_id_fkey` FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `audit_log` ADD CONSTRAINT `audit_log_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `call_centers` ADD CONSTRAINT `call_centers_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agents` ADD CONSTRAINT `agents_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agents` ADD CONSTRAINT `agents_callCenterId_fkey` FOREIGN KEY (`callCenterId`) REFERENCES `call_centers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agents` ADD CONSTRAINT `agents_currentCallId_fkey` FOREIGN KEY (`currentCallId`) REFERENCES `calls`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_agents` ADD CONSTRAINT `ai_agents_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_agents` ADD CONSTRAINT `ai_agents_callCenterId_fkey` FOREIGN KEY (`callCenterId`) REFERENCES `call_centers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_agents` ADD CONSTRAINT `ai_agents_currentCallId_fkey` FOREIGN KEY (`currentCallId`) REFERENCES `calls`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_agents` ADD CONSTRAINT `ai_agents_scriptId_fkey` FOREIGN KEY (`scriptId`) REFERENCES `call_scripts`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `skills` ADD CONSTRAINT `skills_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `skills` ADD CONSTRAINT `skills_callCenterId_fkey` FOREIGN KEY (`callCenterId`) REFERENCES `call_centers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agent_skills` ADD CONSTRAINT `agent_skills_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agent_skills` ADD CONSTRAINT `agent_skills_agentId_fkey` FOREIGN KEY (`agentId`) REFERENCES `agents`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agent_skills` ADD CONSTRAINT `agent_skills_skillId_fkey` FOREIGN KEY (`skillId`) REFERENCES `skills`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `campaigns` ADD CONSTRAINT `campaigns_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `campaigns` ADD CONSTRAINT `campaigns_callCenterId_fkey` FOREIGN KEY (`callCenterId`) REFERENCES `call_centers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_campaigns` ADD CONSTRAINT `ai_campaigns_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_campaigns` ADD CONSTRAINT `ai_campaigns_campaignId_fkey` FOREIGN KEY (`campaignId`) REFERENCES `campaigns`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_campaigns` ADD CONSTRAINT `ai_campaigns_aiAgentId_fkey` FOREIGN KEY (`aiAgentId`) REFERENCES `ai_agents`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `queues` ADD CONSTRAINT `queues_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `queues` ADD CONSTRAINT `queues_callCenterId_fkey` FOREIGN KEY (`callCenterId`) REFERENCES `call_centers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `queue_skill_requirements` ADD CONSTRAINT `queue_skill_requirements_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `queue_skill_requirements` ADD CONSTRAINT `queue_skill_requirements_queueId_fkey` FOREIGN KEY (`queueId`) REFERENCES `queues`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `queue_skill_requirements` ADD CONSTRAINT `queue_skill_requirements_skillId_fkey` FOREIGN KEY (`skillId`) REFERENCES `skills`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `queue_assignments` ADD CONSTRAINT `queue_assignments_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `queue_assignments` ADD CONSTRAINT `queue_assignments_agentId_fkey` FOREIGN KEY (`agentId`) REFERENCES `agents`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `queue_assignments` ADD CONSTRAINT `queue_assignments_queueId_fkey` FOREIGN KEY (`queueId`) REFERENCES `queues`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agent_campaigns` ADD CONSTRAINT `agent_campaigns_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agent_campaigns` ADD CONSTRAINT `agent_campaigns_agentId_fkey` FOREIGN KEY (`agentId`) REFERENCES `agents`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agent_campaigns` ADD CONSTRAINT `agent_campaigns_campaignId_fkey` FOREIGN KEY (`campaignId`) REFERENCES `campaigns`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `calls` ADD CONSTRAINT `calls_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `calls` ADD CONSTRAINT `calls_callCenterId_fkey` FOREIGN KEY (`callCenterId`) REFERENCES `call_centers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `calls` ADD CONSTRAINT `calls_campaignId_fkey` FOREIGN KEY (`campaignId`) REFERENCES `campaigns`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `calls` ADD CONSTRAINT `calls_queueId_fkey` FOREIGN KEY (`queueId`) REFERENCES `queues`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `calls` ADD CONSTRAINT `calls_customerId_fkey` FOREIGN KEY (`customerId`) REFERENCES `sales_customer`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `calls` ADD CONSTRAINT `calls_aiAgentId_fkey` FOREIGN KEY (`aiAgentId`) REFERENCES `ai_agents`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `calls` ADD CONSTRAINT `calls_dispositionId_fkey` FOREIGN KEY (`dispositionId`) REFERENCES `call_dispositions`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `call_dispositions` ADD CONSTRAINT `call_dispositions_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `call_dispositions` ADD CONSTRAINT `call_dispositions_callCenterId_fkey` FOREIGN KEY (`callCenterId`) REFERENCES `call_centers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `call_recordings` ADD CONSTRAINT `call_recordings_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `call_recordings` ADD CONSTRAINT `call_recordings_callCenterId_fkey` FOREIGN KEY (`callCenterId`) REFERENCES `call_centers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `call_recordings` ADD CONSTRAINT `call_recordings_callId_fkey` FOREIGN KEY (`callId`) REFERENCES `calls`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_agent_dialogues` ADD CONSTRAINT `ai_agent_dialogues_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_agent_dialogues` ADD CONSTRAINT `ai_agent_dialogues_callId_fkey` FOREIGN KEY (`callId`) REFERENCES `calls`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_agent_dialogues` ADD CONSTRAINT `ai_agent_dialogues_aiAgentId_fkey` FOREIGN KEY (`aiAgentId`) REFERENCES `ai_agents`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_conversation_logs` ADD CONSTRAINT `ai_conversation_logs_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_conversation_logs` ADD CONSTRAINT `ai_conversation_logs_callId_fkey` FOREIGN KEY (`callId`) REFERENCES `calls`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ai_conversation_logs` ADD CONSTRAINT `ai_conversation_logs_aiAgentId_fkey` FOREIGN KEY (`aiAgentId`) REFERENCES `ai_agents`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `call_scripts` ADD CONSTRAINT `call_scripts_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `call_scripts` ADD CONSTRAINT `call_scripts_callCenterId_fkey` FOREIGN KEY (`callCenterId`) REFERENCES `call_centers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `external_providers` ADD CONSTRAINT `external_providers_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `external_providers` ADD CONSTRAINT `external_providers_callCenterId_fkey` FOREIGN KEY (`callCenterId`) REFERENCES `call_centers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `provider_connections` ADD CONSTRAINT `provider_connections_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `provider_connections` ADD CONSTRAINT `provider_connections_providerId_fkey` FOREIGN KEY (`providerId`) REFERENCES `external_providers`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agent_notes` ADD CONSTRAINT `agent_notes_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agent_notes` ADD CONSTRAINT `agent_notes_agentId_fkey` FOREIGN KEY (`agentId`) REFERENCES `agents`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agent_performance_metrics` ADD CONSTRAINT `agent_performance_metrics_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `agent_performance_metrics` ADD CONSTRAINT `agent_performance_metrics_agentId_fkey` FOREIGN KEY (`agentId`) REFERENCES `agents`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `attribution_event` ADD CONSTRAINT `attribution_event_lead_id_fkey` FOREIGN KEY (`lead_id`) REFERENCES `sales_lead`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `lead_attribution_snapshot` ADD CONSTRAINT `lead_attribution_snapshot_lead_id_fkey` FOREIGN KEY (`lead_id`) REFERENCES `sales_lead`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `tenant_user_count` ADD CONSTRAINT `tenant_user_count_tenant_id_fkey` FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `tenant_user_count_history` ADD CONSTRAINT `tenant_user_count_history_user_count_id_fkey` FOREIGN KEY (`user_count_id`) REFERENCES `tenant_user_count`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `vendor_contact` ADD CONSTRAINT `vendor_contact_vendor_id_fkey` FOREIGN KEY (`vendor_id`) REFERENCES `vendor`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `vendor_address` ADD CONSTRAINT `vendor_address_vendor_id_fkey` FOREIGN KEY (`vendor_id`) REFERENCES `vendor`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `purchase_order` ADD CONSTRAINT `purchase_order_vendor_id_fkey` FOREIGN KEY (`vendor_id`) REFERENCES `vendor`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `purchase_order` ADD CONSTRAINT `purchase_order_requisition_id_fkey` FOREIGN KEY (`requisition_id`) REFERENCES `purchase_requisition`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `po_line_item` ADD CONSTRAINT `po_line_item_po_id_fkey` FOREIGN KEY (`po_id`) REFERENCES `purchase_order`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `goods_receipt` ADD CONSTRAINT `goods_receipt_po_id_fkey` FOREIGN KEY (`po_id`) REFERENCES `purchase_order`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `grn_line_item` ADD CONSTRAINT `grn_line_item_grn_id_fkey` FOREIGN KEY (`grn_id`) REFERENCES `goods_receipt`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `vendor_invoice` ADD CONSTRAINT `vendor_invoice_vendor_id_fkey` FOREIGN KEY (`vendor_id`) REFERENCES `vendor`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `inventory_item_vendor` ADD CONSTRAINT `inventory_item_vendor_inventory_item_id_fkey` FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `inventory_item_vendor` ADD CONSTRAINT `inventory_item_vendor_vendor_id_fkey` FOREIGN KEY (`vendor_id`) REFERENCES `vendor`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `stock_level` ADD CONSTRAINT `stock_level_inventory_item_id_fkey` FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `stock_level` ADD CONSTRAINT `stock_level_warehouse_id_fkey` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `stock_movement` ADD CONSTRAINT `stock_movement_inventory_item_id_fkey` FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `stock_movement` ADD CONSTRAINT `stock_movement_warehouse_id_fkey` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `inventory_batch` ADD CONSTRAINT `inventory_batch_inventory_item_id_fkey` FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `inventory_batch` ADD CONSTRAINT `inventory_batch_purchase_order_id_fkey` FOREIGN KEY (`purchase_order_id`) REFERENCES `purchase_order`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `inventory_serial` ADD CONSTRAINT `inventory_serial_inventory_item_id_fkey` FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `stock_adjustment_line` ADD CONSTRAINT `stock_adjustment_line_adjustment_id_fkey` FOREIGN KEY (`adjustment_id`) REFERENCES `stock_adjustment`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `stock_adjustment_line` ADD CONSTRAINT `stock_adjustment_line_inventory_item_id_fkey` FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `physical_inventory_detail` ADD CONSTRAINT `physical_inventory_detail_physical_inventory_id_fkey` FOREIGN KEY (`physical_inventory_id`) REFERENCES `physical_inventory`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `physical_inventory_detail` ADD CONSTRAINT `physical_inventory_detail_inventory_item_id_fkey` FOREIGN KEY (`inventory_item_id`) REFERENCES `inventory_item`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `inventory_transfer_line` ADD CONSTRAINT `inventory_transfer_line_transfer_id_fkey` FOREIGN KEY (`transfer_id`) REFERENCES `inventory_transfer`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `role_permission` ADD CONSTRAINT `role_permission_role_id_fkey` FOREIGN KEY (`role_id`) REFERENCES `role`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `role_permission` ADD CONSTRAINT `role_permission_permission_id_fkey` FOREIGN KEY (`permission_id`) REFERENCES `permission`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `user_role` ADD CONSTRAINT `user_role_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `user_role` ADD CONSTRAINT `user_role_role_id_fkey` FOREIGN KEY (`role_id`) REFERENCES `role`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `access_log` ADD CONSTRAINT `access_log_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `access_log` ADD CONSTRAINT `access_log_resource_id_fkey` FOREIGN KEY (`resource_id`) REFERENCES `resource`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `resource_access` ADD CONSTRAINT `resource_access_resource_id_fkey` FOREIGN KEY (`resource_id`) REFERENCES `resource`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `field_level_permission` ADD CONSTRAINT `field_level_permission_permission_id_fkey` FOREIGN KEY (`permission_id`) REFERENCES `permission`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `role_delegation` ADD CONSTRAINT `role_delegation_role_id_fkey` FOREIGN KEY (`role_id`) REFERENCES `role`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `journal_entry_detail` ADD CONSTRAINT `journal_entry_detail_journal_entry_id_fkey` FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entry`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `journal_entry_detail` ADD CONSTRAINT `journal_entry_detail_account_id_fkey` FOREIGN KEY (`account_id`) REFERENCES `chart_of_account`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `gl_posting_template` ADD CONSTRAINT `gl_posting_template_source_account_id_fkey` FOREIGN KEY (`source_account_id`) REFERENCES `chart_of_account`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `gl_posting_template` ADD CONSTRAINT `gl_posting_template_target_account_id_fkey` FOREIGN KEY (`target_account_id`) REFERENCES `chart_of_account`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `gl_posting_template` ADD CONSTRAINT `gl_posting_template_debit_account_id_fkey` FOREIGN KEY (`debit_account_id`) REFERENCES `chart_of_account`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `gl_posting_template` ADD CONSTRAINT `gl_posting_template_credit_account_id_fkey` FOREIGN KEY (`credit_account_id`) REFERENCES `chart_of_account`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `employee_details` ADD CONSTRAINT `employee_details_employee_id_fkey` FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `leave_balance` ADD CONSTRAINT `leave_balance_employee_id_fkey` FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `employee_epf_registration` ADD CONSTRAINT `employee_epf_registration_employee_id_fkey` FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `employee_esi_registration` ADD CONSTRAINT `employee_esi_registration_employee_id_fkey` FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `epf_contribution` ADD CONSTRAINT `epf_contribution_employee_id_fkey` FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `esi_contribution` ADD CONSTRAINT `esi_contribution_employee_id_fkey` FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `employee_compliance_checklist` ADD CONSTRAINT `employee_compliance_checklist_employee_id_fkey` FOREIGN KEY (`employee_id`) REFERENCES `employee`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `quotation_line_item` ADD CONSTRAINT `quotation_line_item_quotation_id_fkey` FOREIGN KEY (`quotation_id`) REFERENCES `quotation`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `sales_order_line_item` ADD CONSTRAINT `sales_order_line_item_order_id_fkey` FOREIGN KEY (`order_id`) REFERENCES `sales_order`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `sales_invoice` ADD CONSTRAINT `sales_invoice_order_id_fkey` FOREIGN KEY (`order_id`) REFERENCES `sales_order`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `sales_invoice_line_item` ADD CONSTRAINT `sales_invoice_line_item_invoice_id_fkey` FOREIGN KEY (`invoice_id`) REFERENCES `sales_invoice`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `sales_status_history` ADD CONSTRAINT `sales_status_history_lead_id_fkey` FOREIGN KEY (`lead_id`) REFERENCES `sales_lead`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `sales_follow_up` ADD CONSTRAINT `sales_follow_up_lead_id_fkey` FOREIGN KEY (`lead_id`) REFERENCES `sales_lead`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `property_block` ADD CONSTRAINT `property_block_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `property_unit` ADD CONSTRAINT `property_unit_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `property_project`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `property_unit` ADD CONSTRAINT `property_unit_block_id_fkey` FOREIGN KEY (`block_id`) REFERENCES `property_block`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `unit_cost_sheet` ADD CONSTRAINT `unit_cost_sheet_unit_id_fkey` FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `payment_plan` ADD CONSTRAINT `payment_plan_unit_id_fkey` FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `installment` ADD CONSTRAINT `installment_payment_plan_id_fkey` FOREIGN KEY (`payment_plan_id`) REFERENCES `payment_plan`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `booking` ADD CONSTRAINT `booking_unit_id_fkey` FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `communication_message` ADD CONSTRAINT `communication_message_session_id_fkey` FOREIGN KEY (`session_id`) REFERENCES `communication_session`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ivr_menu_option` ADD CONSTRAINT `ivr_menu_option_menu_id_fkey` FOREIGN KEY (`menu_id`) REFERENCES `ivr_menu`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `team_chat_member` ADD CONSTRAINT `team_chat_member_channel_id_fkey` FOREIGN KEY (`channel_id`) REFERENCES `team_chat_channel`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `team_chat_message` ADD CONSTRAINT `team_chat_message_channel_id_fkey` FOREIGN KEY (`channel_id`) REFERENCES `team_chat_channel`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `team_chat_reaction` ADD CONSTRAINT `team_chat_reaction_message_id_fkey` FOREIGN KEY (`message_id`) REFERENCES `team_chat_message`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `work_item_comment` ADD CONSTRAINT `work_item_comment_work_item_id_fkey` FOREIGN KEY (`work_item_id`) REFERENCES `work_item`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `depreciation_schedule` ADD CONSTRAINT `depreciation_schedule_asset_id_fkey` FOREIGN KEY (`asset_id`) REFERENCES `fixed_asset`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `cost_allocation` ADD CONSTRAINT `cost_allocation_cost_center_id_fkey` FOREIGN KEY (`cost_center_id`) REFERENCES `cost_center`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `budget` ADD CONSTRAINT `budget_cost_center_id_fkey` FOREIGN KEY (`cost_center_id`) REFERENCES `cost_center`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `dashboard_widget` ADD CONSTRAINT `dashboard_widget_dashboard_id_fkey` FOREIGN KEY (`dashboard_id`) REFERENCES `dashboard`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `kpi_value` ADD CONSTRAINT `kpi_value_kpi_id_fkey` FOREIGN KEY (`kpi_id`) REFERENCES `kpi_definition`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `report_execution` ADD CONSTRAINT `report_execution_report_id_fkey` FOREIGN KEY (`report_id`) REFERENCES `report`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `mobile_session` ADD CONSTRAINT `mobile_session_device_id_fkey` FOREIGN KEY (`device_id`) REFERENCES `mobile_device`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `mobile_push_notification` ADD CONSTRAINT `mobile_push_notification_device_id_fkey` FOREIGN KEY (`device_id`) REFERENCES `mobile_device`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `construction_phase` ADD CONSTRAINT `construction_phase_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `construction_project`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `construction_material` ADD CONSTRAINT `construction_material_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `construction_project`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `construction_labor` ADD CONSTRAINT `construction_labor_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `construction_project`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `construction_equipment` ADD CONSTRAINT `construction_equipment_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `construction_project`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `construction_quality_check` ADD CONSTRAINT `construction_quality_check_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `construction_project`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `construction_safety` ADD CONSTRAINT `construction_safety_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `construction_project`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `site_visit` ADD CONSTRAINT `site_visit_project_id_fkey` FOREIGN KEY (`project_id`) REFERENCES `construction_project`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ImportShipment` ADD CONSTRAINT `ImportShipment_importLicenseId_fkey` FOREIGN KEY (`importLicenseId`) REFERENCES `ImportLicense`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ImportHSCodeLineItem` ADD CONSTRAINT `ImportHSCodeLineItem_importShipmentId_fkey` FOREIGN KEY (`importShipmentId`) REFERENCES `ImportShipment`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ImportHSCodeLineItem` ADD CONSTRAINT `ImportHSCodeLineItem_stockLevelId_fkey` FOREIGN KEY (`stockLevelId`) REFERENCES `stock_level`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ExportHSCodeLineItem` ADD CONSTRAINT `ExportHSCodeLineItem_exportShipmentId_fkey` FOREIGN KEY (`exportShipmentId`) REFERENCES `ExportShipment`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ExportHSCodeLineItem` ADD CONSTRAINT `ExportHSCodeLineItem_soLineItemId_fkey` FOREIGN KEY (`soLineItemId`) REFERENCES `sales_order_line_item`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `DutyCalculation` ADD CONSTRAINT `DutyCalculation_importShipmentId_fkey` FOREIGN KEY (`importShipmentId`) REFERENCES `ImportShipment`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `DutyCalculation` ADD CONSTRAINT `DutyCalculation_lineItemId_fkey` FOREIGN KEY (`lineItemId`) REFERENCES `ImportHSCodeLineItem`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `DutyCalculation` ADD CONSTRAINT `DutyCalculation_tradeAgreementId_fkey` FOREIGN KEY (`tradeAgreementId`) REFERENCES `TradeAgreementPreference`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `CustomsDocumentation` ADD CONSTRAINT `CustomsDocumentation_importShipmentId_fkey` FOREIGN KEY (`importShipmentId`) REFERENCES `ImportShipment`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `CustomsDocumentation` ADD CONSTRAINT `CustomsDocumentation_exportShipmentId_fkey` FOREIGN KEY (`exportShipmentId`) REFERENCES `ExportShipment`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ShipmentTracking` ADD CONSTRAINT `ShipmentTracking_importShipmentId_fkey` FOREIGN KEY (`importShipmentId`) REFERENCES `ImportShipment`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ShipmentTracking` ADD CONSTRAINT `ShipmentTracking_exportShipmentId_fkey` FOREIGN KEY (`exportShipmentId`) REFERENCES `ExportShipment`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ComplianceCheckImEx` ADD CONSTRAINT `ComplianceCheckImEx_importShipmentId_fkey` FOREIGN KEY (`importShipmentId`) REFERENCES `ImportShipment`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ComplianceCheckImEx` ADD CONSTRAINT `ComplianceCheckImEx_exportShipmentId_fkey` FOREIGN KEY (`exportShipmentId`) REFERENCES `ExportShipment`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `tenant_credentials` ADD CONSTRAINT `tenant_credentials_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;
