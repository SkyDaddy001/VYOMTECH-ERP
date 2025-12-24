-- CreateTable
CREATE TABLE `tenant` (
    `id` CHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `domain` VARCHAR(255) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `maxUsers` INTEGER NOT NULL DEFAULT 100,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    UNIQUE INDEX `tenant_domain_key`(`domain`),
    INDEX `tenant_domain_idx`(`domain`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `user` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `passwordHash` VARCHAR(500) NOT NULL,
    `firstName` VARCHAR(100) NOT NULL,
    `lastName` VARCHAR(100) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `lastLogin` DATETIME(3) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    INDEX `user_tenantId_idx`(`tenantId`),
    UNIQUE INDEX `user_tenantId_email_key`(`tenantId`, `email`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `role` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `description` TEXT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `role_tenantId_idx`(`tenantId`),
    UNIQUE INDEX `role_tenantId_name_key`(`tenantId`, `name`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `permission` (
    `id` CHAR(36) NOT NULL,
    `roleId` CHAR(36) NOT NULL,
    `action` VARCHAR(100) NOT NULL,
    `resource` VARCHAR(100) NOT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `permission_roleId_idx`(`roleId`),
    UNIQUE INDEX `permission_roleId_action_resource_key`(`roleId`, `action`, `resource`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `user_role` (
    `id` CHAR(36) NOT NULL,
    `userId` CHAR(36) NOT NULL,
    `roleId` CHAR(36) NOT NULL,
    `assignedAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `user_role_userId_idx`(`userId`),
    INDEX `user_role_roleId_idx`(`roleId`),
    UNIQUE INDEX `user_role_userId_roleId_key`(`userId`, `roleId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `auth_token` (
    `id` CHAR(36) NOT NULL,
    `userId` CHAR(36) NOT NULL,
    `token` VARCHAR(500) NOT NULL,
    `tokenType` VARCHAR(50) NOT NULL DEFAULT 'bearer',
    `expiresAt` DATETIME(3) NOT NULL,
    `revokedAt` DATETIME(3) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `auth_token_userId_idx`(`userId`),
    INDEX `auth_token_token_idx`(`token`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `sales_lead` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `leadCode` VARCHAR(50) NOT NULL,
    `leadName` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NULL,
    `phone` VARCHAR(20) NULL,
    `company` VARCHAR(255) NULL,
    `leadSource` VARCHAR(100) NOT NULL DEFAULT 'manual',
    `leadValue` DECIMAL(15, 2) NOT NULL DEFAULT 0,
    `description` TEXT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'new',
    `assignedTo` CHAR(36) NULL,
    `deletedAt` DATETIME(3) NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    UNIQUE INDEX `sales_lead_leadCode_key`(`leadCode`),
    INDEX `sales_lead_tenantId_idx`(`tenantId`),
    INDEX `sales_lead_status_idx`(`status`),
    INDEX `sales_lead_assignedTo_idx`(`assignedTo`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `lead_status_history` (
    `id` CHAR(36) NOT NULL,
    `leadId` CHAR(36) NOT NULL,
    `oldStatus` VARCHAR(50) NOT NULL,
    `newStatus` VARCHAR(50) NOT NULL,
    `changedBy` CHAR(36) NULL,
    `notes` TEXT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `lead_status_history_leadId_idx`(`leadId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `quotation` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `leadId` CHAR(36) NOT NULL,
    `quotationNumber` VARCHAR(50) NOT NULL,
    `amount` DECIMAL(15, 2) NOT NULL,
    `validUntil` DATE NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updatedAt` DATETIME(3) NOT NULL,

    UNIQUE INDEX `quotation_quotationNumber_key`(`quotationNumber`),
    INDEX `quotation_tenantId_idx`(`tenantId`),
    INDEX `quotation_leadId_idx`(`leadId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `audit_log` (
    `id` CHAR(36) NOT NULL,
    `tenantId` CHAR(36) NOT NULL,
    `userId` CHAR(36) NULL,
    `entityType` VARCHAR(100) NOT NULL,
    `entityId` VARCHAR(36) NULL,
    `action` VARCHAR(50) NOT NULL,
    `changes` LONGTEXT NULL,
    `metadata` TEXT NULL,
    `createdAt` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    INDEX `audit_log_tenantId_idx`(`tenantId`),
    INDEX `audit_log_userId_idx`(`userId`),
    INDEX `audit_log_entityType_entityId_idx`(`entityType`, `entityId`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- AddForeignKey
ALTER TABLE `user` ADD CONSTRAINT `user_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `role` ADD CONSTRAINT `role_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `permission` ADD CONSTRAINT `permission_roleId_fkey` FOREIGN KEY (`roleId`) REFERENCES `role`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `user_role` ADD CONSTRAINT `user_role_userId_fkey` FOREIGN KEY (`userId`) REFERENCES `user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `user_role` ADD CONSTRAINT `user_role_roleId_fkey` FOREIGN KEY (`roleId`) REFERENCES `role`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `auth_token` ADD CONSTRAINT `auth_token_userId_fkey` FOREIGN KEY (`userId`) REFERENCES `user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `sales_lead` ADD CONSTRAINT `sales_lead_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `lead_status_history` ADD CONSTRAINT `lead_status_history_leadId_fkey` FOREIGN KEY (`leadId`) REFERENCES `sales_lead`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `quotation` ADD CONSTRAINT `quotation_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `quotation` ADD CONSTRAINT `quotation_leadId_fkey` FOREIGN KEY (`leadId`) REFERENCES `sales_lead`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `audit_log` ADD CONSTRAINT `audit_log_tenantId_fkey` FOREIGN KEY (`tenantId`) REFERENCES `tenant`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `audit_log` ADD CONSTRAINT `audit_log_userId_fkey` FOREIGN KEY (`userId`) REFERENCES `user`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `audit_log` ADD CONSTRAINT `audit_log_entityId_fkey` FOREIGN KEY (`entityId`) REFERENCES `sales_lead`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;
