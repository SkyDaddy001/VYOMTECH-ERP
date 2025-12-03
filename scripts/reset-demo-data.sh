#!/bin/bash

# Demo Data Reset Script
# Resets Vyomtech demo data every 30 days
# Scheduled via cron: 0 0 1 * * /path/to/reset-demo-data.sh

DB_HOST="${DB_HOST:-mysql}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-callcenter_user}"
DB_PASSWORD="${DB_PASSWORD:-secure_app_pass}"
DB_NAME="${DB_NAME:-callcenter}"
DEMO_TENANT_ID="demo_vyomtech_001"

echo "=========================================="
echo "Vyomtech Demo Data Reset ($(date))"
echo "=========================================="

# Function to run MySQL command
run_mysql() {
    MYSQL_PWD="$DB_PASSWORD" mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME"
}

# Delete old demo data (keep last backup for 30 days)
echo "Clearing previous demo data..."
run_mysql << 'SQL'
SET FOREIGN_KEY_CHECKS=0;

DELETE FROM partner_lead_credits WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM partner_leads WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM partner_payouts WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM partner_payout_details WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM partner_activities WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM partner_users WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM partners WHERE tenant_id = 'demo_vyomtech_001';

DELETE FROM gamification_points_history WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM gamification_badges WHERE tenant_id = 'demo_vyomtech_001';

DELETE FROM progress_tracking WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM compliance_records WHERE tenant_id = 'demo_vyomtech_001';

DELETE FROM task WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM call WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM campaign WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM campaign_recipient WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM lead WHERE tenant_id = 'demo_vyomtech_001';

DELETE FROM agent WHERE tenant_id = 'demo_vyomtech_001';
DELETE FROM construction_projects WHERE tenant_id = 'demo_vyomtech_001';

SET FOREIGN_KEY_CHECKS=1;
SQL

echo "✓ Previous data cleared"

# Reload fresh demo data
echo "Reloading fresh demo data..."
MYSQL_PWD="$DB_PASSWORD" mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" << 'SQL'
SET FOREIGN_KEY_CHECKS=0;

-- Demo Partners
INSERT IGNORE INTO partners (tenant_id, partner_code, organization_name, partner_type, status, contact_email, contact_person, created_by) 
VALUES 
('demo_vyomtech_001', 'DEMO_PORTAL', 'Vyomtech Portal Demo', 'portal', 'active', 'demo@vyomtech.com', 'Demo Admin', 1),
('demo_vyomtech_001', 'DEMO_CHANNEL', 'Demo Channel Partner', 'channel_partner', 'active', 'channel@demo.vyomtech.com', 'Channel Manager', 1),
('demo_vyomtech_001', 'DEMO_VENDOR', 'Demo Vendor Solutions', 'vendor', 'active', 'vendor@demo.vyomtech.com', 'Vendor Lead', 1),
('demo_vyomtech_001', 'DEMO_CUSTOMER', 'Demo Customer Account', 'customer', 'active', 'customer@demo.vyomtech.com', 'Customer Manager', 1);

-- Demo Partner Users
INSERT IGNORE INTO partner_users (partner_id, tenant_id, email, first_name, last_name, password_hash, role, is_active) 
VALUES 
(1, 'demo_vyomtech_001', 'demo@vyomtech.com', 'Demo', 'Admin', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'admin', TRUE),
(2, 'demo_vyomtech_001', 'channel@demo.vyomtech.com', 'Channel', 'Manager', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'admin', TRUE),
(3, 'demo_vyomtech_001', 'vendor@demo.vyomtech.com', 'Vendor', 'Lead', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'admin', TRUE),
(4, 'demo_vyomtech_001', 'customer@demo.vyomtech.com', 'Customer', 'Manager', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'admin', TRUE);

-- Demo Agents
INSERT IGNORE INTO agent (tenant_id, agent_name, email, phone, status, performance_score, created_at, updated_at) 
VALUES 
('demo_vyomtech_001', 'Rajesh Kumar', 'rajesh@demo.vyomtech.com', '+91-9876543210', 'active', 95, NOW(), NOW()),
('demo_vyomtech_001', 'Priya Singh', 'priya@demo.vyomtech.com', '+91-8765432109', 'active', 92, NOW(), NOW()),
('demo_vyomtech_001', 'Arun Patel', 'arun@demo.vyomtech.com', '+91-7654321098', 'active', 88, NOW(), NOW()),
('demo_vyomtech_001', 'Neha Sharma', 'neha@demo.vyomtech.com', '+91-6543210987', 'active', 90, NOW(), NOW());

-- Demo Leads
INSERT IGNORE INTO lead (tenant_id, lead_title, lead_value, contact_name, contact_email, contact_phone, property_type, location, status, created_at, updated_at) 
VALUES 
('demo_vyomtech_001', 'High Value Residential Project', 5000000, 'Amit Kumar', 'amit@example.com', '+91-9876543211', 'residential', 'Mumbai, Maharashtra', 'active', NOW(), NOW()),
('demo_vyomtech_001', 'Commercial Space Inquiry', 3500000, 'Sneha Desai', 'sneha@example.com', '+91-9876543212', 'commercial', 'Bangalore, Karnataka', 'active', NOW(), NOW()),
('demo_vyomtech_001', 'Plot Purchase Interest', 2000000, 'Vikram Singh', 'vikram@example.com', '+91-9876543213', 'plot', 'Delhi, Delhi', 'active', NOW(), NOW()),
('demo_vyomtech_001', 'Rental Inquiry', 1500000, 'Meera Nair', 'meera@example.com', '+91-9876543214', 'rental', 'Hyderabad, Telangana', 'active', NOW(), NOW()),
('demo_vyomtech_001', 'Apartment Pre-booking', 4000000, 'Rohan Gupta', 'rohan@example.com', '+91-9876543215', 'apartment', 'Pune, Maharashtra', 'active', NOW(), NOW());

-- Demo Campaigns
INSERT IGNORE INTO campaign (tenant_id, campaign_name, campaign_type, status, budget, target_audience, created_at, updated_at) 
VALUES 
('demo_vyomtech_001', 'Summer Residential Drive 2025', 'email', 'active', 500000, 'Homebuyers', NOW(), NOW()),
('demo_vyomtech_001', 'Commercial Real Estate Expo', 'event', 'active', 1000000, 'Business Owners', NOW(), NOW()),
('demo_vyomtech_001', 'Digital Marketing Campaign', 'social_media', 'active', 300000, 'Young Professionals', NOW(), NOW());

-- Demo Projects
INSERT IGNORE INTO construction_projects (tenant_id, project_name, project_type, location, total_cost, status, start_date, completion_date) 
VALUES 
('demo_vyomtech_001', 'Skyrise Towers Mumbai', 'residential', 'Mumbai, Maharashtra', 50000000, 'active', '2024-01-15', '2026-12-31'),
('demo_vyomtech_001', 'Tech Park Bangalore', 'commercial', 'Bangalore, Karnataka', 100000000, 'active', '2023-06-01', '2025-12-31');

SET FOREIGN_KEY_CHECKS=1;
SQL

echo "✓ Fresh demo data reloaded"

echo ""
echo "=========================================="
echo "Demo Reset Completed!"
echo "=========================================="
echo "Demo Credentials:"
echo "  Portal Admin: demo@vyomtech.com (password: demo123)"
echo "  Channel: channel@demo.vyomtech.com (password: demo123)"
echo "  Vendor: vendor@demo.vyomtech.com (password: demo123)"
echo "  Customer: customer@demo.vyomtech.com (password: demo123)"
echo "=========================================="
