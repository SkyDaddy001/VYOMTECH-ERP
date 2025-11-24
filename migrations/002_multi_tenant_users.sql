-- Multi-Tenant User Management Migration
-- Migration: 002_multi_tenant_users.sql

-- Add current_tenant_id column to users table
ALTER TABLE `user` ADD COLUMN current_tenant_id VARCHAR(36) AFTER tenant_id;
ALTER TABLE `user` ADD INDEX idx_current_tenant (current_tenant_id);

-- Add foreign key constraint for current_tenant_id
ALTER TABLE `user` ADD CONSTRAINT fk_user_current_tenant 
  FOREIGN KEY (current_tenant_id) REFERENCES tenant(id) ON DELETE SET NULL;

-- Create tenant_members table for multi-tenant membership
CREATE TABLE IF NOT EXISTS tenant_members (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    user_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL,
    role ENUM('admin', 'member', 'viewer') DEFAULT 'member',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    UNIQUE KEY unique_tenant_member (tenant_id, user_id),
    INDEX idx_email (email),
    INDEX idx_tenant (tenant_id),
    INDEX idx_user (user_id),
    INDEX idx_role (role)
);

-- Populate tenant_members table with existing users (initial data)
INSERT INTO tenant_members (id, tenant_id, user_id, email, role, created_at, updated_at)
SELECT 
    UUID(),
    tenant_id,
    id,
    email,
    'admin' as role,
    created_at,
    updated_at
FROM `user`
ON DUPLICATE KEY UPDATE
    updated_at = NOW();
