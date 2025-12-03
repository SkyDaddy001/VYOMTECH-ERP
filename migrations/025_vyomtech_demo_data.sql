-- Migration: Vyomtech Demo Data
-- Version: 025
-- Description: Comprehensive dummy data for demo tenant

SET FOREIGN_KEY_CHECKS=0;

-- ============================================================================
-- DEMO TENANT SETUP
-- ============================================================================

INSERT IGNORE INTO tenant (id, name, domain, status, max_users, max_concurrent_calls, ai_budget_monthly, created_at, updated_at) 
VALUES ('demo_vyomtech_001', 'Vyomtech Demo Company', 'demo.vyomtech.com', 'active', 100, 50, 1000.00, NOW(), NOW());

-- ============================================================================
-- DEMO SYSTEM AND PARTNER USERS (FOR AUTHENTICATION)
-- ============================================================================

INSERT IGNORE INTO `user` (email, password_hash, role, tenant_id, current_tenant_id, created_at, updated_at)
VALUES 
('master.admin@vyomtech.com', '$2a$10$AMsTEjlqMasSrxBTjiL3FuojehMAIIK.cqREZWlGfBmKZ6fPZUcym', 'admin', 'demo_vyomtech_001', 'demo_vyomtech_001', NOW(), NOW()),
('rajesh@demo.vyomtech.com', '$2a$10$AMsTEjlqMasSrxBTjiL3FuojehMAIIK.cqREZWlGfBmKZ6fPZUcym', 'agent', 'demo_vyomtech_001', 'demo_vyomtech_001', NOW(), NOW()),
('priya@demo.vyomtech.com', '$2a$10$AMsTEjlqMasSrxBTjiL3FuojehMAIIK.cqREZWlGfBmKZ6fPZUcym', 'agent', 'demo_vyomtech_001', 'demo_vyomtech_001', NOW(), NOW()),
('arun@demo.vyomtech.com', '$2a$10$AMsTEjlqMasSrxBTjiL3FuojehMAIIK.cqREZWlGfBmKZ6fPZUcym', 'agent', 'demo_vyomtech_001', 'demo_vyomtech_001', NOW(), NOW()),
('neha@demo.vyomtech.com', '$2a$10$AMsTEjlqMasSrxBTjiL3FuojehMAIIK.cqREZWlGfBmKZ6fPZUcym', 'agent', 'demo_vyomtech_001', 'demo_vyomtech_001', NOW(), NOW()),
('demo@vyomtech.com', '$2a$10$AMsTEjlqMasSrxBTjiL3FuojehMAIIK.cqREZWlGfBmKZ6fPZUcym', 'partner_admin', 'demo_vyomtech_001', 'demo_vyomtech_001', NOW(), NOW()),
('channel@demo.vyomtech.com', '$2a$10$AMsTEjlqMasSrxBTjiL3FuojehMAIIK.cqREZWlGfBmKZ6fPZUcym', 'partner_admin', 'demo_vyomtech_001', 'demo_vyomtech_001', NOW(), NOW()),
('vendor@demo.vyomtech.com', '$2a$10$AMsTEjlqMasSrxBTjiL3FuojehMAIIK.cqREZWlGfBmKZ6fPZUcym', 'partner_admin', 'demo_vyomtech_001', 'demo_vyomtech_001', NOW(), NOW()),
('customer@demo.vyomtech.com', '$2a$10$AMsTEjlqMasSrxBTjiL3FuojehMAIIK.cqREZWlGfBmKZ6fPZUcym', 'partner_admin', 'demo_vyomtech_001', 'demo_vyomtech_001', NOW(), NOW());

-- ============================================================================
-- DEMO PARTNERS
-- ============================================================================

INSERT IGNORE INTO partners (
    tenant_id, partner_code, organization_name, partner_type, status,
    contact_email, contact_phone, contact_person, website, description,
    address, city, state, country, zip_code, tax_id,
    commission_percentage, lead_price, monthly_quota,
    created_by, created_at, updated_at
) VALUES 
('demo_vyomtech_001', 'DEMO_PORTAL', 'Vyomtech White-Label Portal', 'portal', 'active',
'demo@vyomtech.com', '+91-9000000001', 'Demo Admin', 'https://portal.demo.vyomtech.com', 'Demonstration white-label portal for resellers',
'101 Tech Park', 'Bangalore', 'Karnataka', 'India', '560001', 'GST00DEM001',
10.00, 0.00, 100,
1, NOW(), NOW()),

('demo_vyomtech_001', 'DEMO_CHANNEL', 'Demo Channel Partner Network', 'channel_partner', 'active',
'channel@demo.vyomtech.com', '+91-9000000002', 'Channel Manager', 'https://partner.demo.vyomtech.com', 'Channel partner for lead distribution and sales',
'202 Commerce Hub', 'Mumbai', 'Maharashtra', 'India', '400001', 'GST00DEM002',
15.00, 500.00, 150,
1, NOW(), NOW()),

('demo_vyomtech_001', 'DEMO_VENDOR', 'Demo Vendor Solutions Ltd', 'vendor', 'active',
'vendor@demo.vyomtech.com', '+91-9000000003', 'Vendor Lead', 'https://vendor.demo.vyomtech.com', 'Vendor partner for supply chain and services',
'303 Industrial Zone', 'Delhi', 'Delhi', 'India', '110001', 'GST00DEM003',
5.00, 1000.00, 200,
1, NOW(), NOW()),

('demo_vyomtech_001', 'DEMO_CUSTOMER', 'Demo Customer Organization', 'customer', 'active',
'customer@demo.vyomtech.com', '+91-9000000004', 'Customer Manager', 'https://customer.demo.vyomtech.com', 'Direct customer account for bulk lead purchase',
'404 Business Plaza', 'Pune', 'Maharashtra', 'India', '411001', 'GST00DEM004',
0.00, 2000.00, 500,
1, NOW(), NOW());

-- ============================================================================
-- DEMO PARTNER USERS
-- ============================================================================

INSERT IGNORE INTO partner_users (
    partner_id, tenant_id, email, password_hash, first_name, last_name, phone, role, is_active, created_at, updated_at
) VALUES 
(1, 'demo_vyomtech_001', 'demo@vyomtech.com', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'Demo', 'Admin', '+91-9000000001', 'admin', TRUE, NOW(), NOW()),
(2, 'demo_vyomtech_001', 'channel@demo.vyomtech.com', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'Channel', 'Manager', '+91-9000000002', 'admin', TRUE, NOW(), NOW()),
(3, 'demo_vyomtech_001', 'vendor@demo.vyomtech.com', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'Vendor', 'Lead', '+91-9000000003', 'manager', TRUE, NOW(), NOW()),
(4, 'demo_vyomtech_001', 'customer@demo.vyomtech.com', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'Customer', 'Manager', '+91-9000000004', 'admin', TRUE, NOW(), NOW());

-- ============================================================================
-- DEMO AGENTS
-- ============================================================================

INSERT IGNORE INTO agent (id, tenant_id, agent_code, first_name, last_name, email, phone, status, agent_type, skills, available, created_at, updated_at) 
VALUES 
(UUID(), 'demo_vyomtech_001', 'AGENT001', 'Rajesh', 'Kumar', 'rajesh@demo.vyomtech.com', '+91-9876543210', 'available', 'inbound', '["Customer Support","Sales"]', TRUE, NOW(), NOW()),
(UUID(), 'demo_vyomtech_001', 'AGENT002', 'Priya', 'Singh', 'priya@demo.vyomtech.com', '+91-8765432109', 'available', 'inbound', '["Technical Support","Billing"]', TRUE, NOW(), NOW()),
(UUID(), 'demo_vyomtech_001', 'AGENT003', 'Arun', 'Patel', 'arun@demo.vyomtech.com', '+91-7654321098', 'available', 'inbound', '["Sales","Lead Management"]', TRUE, NOW(), NOW()),
(UUID(), 'demo_vyomtech_001', 'AGENT004', 'Neha', 'Sharma', 'neha@demo.vyomtech.com', '+91-6543210987', 'available', 'inbound', '["Customer Support"]', TRUE, NOW(), NOW());

-- ============================================================================
-- DEMO LEADS
-- ============================================================================

INSERT IGNORE INTO sales_lead (id, tenant_id, lead_code, first_name, last_name, email, phone, company_name, status, source, created_by, created_at, updated_at) 
VALUES 
(UUID(), 'demo_vyomtech_001', 'LEAD001', 'Amit', 'Kumar', 'amit@example.com', '+91-9876543211', 'Residential Projects', 'new', 'Direct', 1, NOW(), NOW()),
(UUID(), 'demo_vyomtech_001', 'LEAD002', 'Sneha', 'Desai', 'sneha@example.com', '+91-9876543212', 'Commercial Enterprises', 'new', 'Digital', 1, NOW(), NOW()),
(UUID(), 'demo_vyomtech_001', 'LEAD003', 'Vikram', 'Singh', 'vikram@example.com', '+91-9876543213', 'Plot Developers', 'new', 'Referral', 1, NOW(), NOW()),
(UUID(), 'demo_vyomtech_001', 'LEAD004', 'Meera', 'Nair', 'meera@example.com', '+91-9876543214', 'Rental Seekers', 'new', 'Portal', 1, NOW(), NOW()),
(UUID(), 'demo_vyomtech_001', 'LEAD005', 'Rohan', 'Gupta', 'rohan@example.com', '+91-9876543215', 'Investment Partners', 'new', 'Event', 1, NOW(), NOW());

-- ============================================================================
-- DEMO CAMPAIGNS
-- ============================================================================

INSERT IGNORE INTO campaign (id, tenant_id, campaign_name, campaign_type, description, status, start_date, target_leads, budget, assigned_agents, created_by, created_at, updated_at) 
VALUES 
(UUID(), 'demo_vyomtech_001', 'Summer Residential Drive 2025', 'email', 'Summer campaign for residential properties', 'planning', '2025-06-01', 500, 500000, 2, 1, NOW(), NOW()),
(UUID(), 'demo_vyomtech_001', 'Commercial Real Estate Expo', 'event', 'Commercial real estate expo and networking event', 'planning', '2025-07-15', 1000, 1000000, 4, 1, NOW(), NOW()),
(UUID(), 'demo_vyomtech_001', 'Digital Marketing Campaign', 'social_media', 'Digital marketing campaign across social media platforms', 'planning', '2025-05-01', 300, 300000, 2, 1, NOW(), NOW()),
(UUID(), 'demo_vyomtech_001', 'Corporate Bulk Purchase', 'direct_sales', 'Corporate bulk purchase campaign for B2B customers', 'planning', '2025-08-01', 800, 800000, 3, 1, NOW(), NOW());

SET FOREIGN_KEY_CHECKS=1;
