-- Sample Data for Dashboard Testing

-- Sites: 4 records (3 for tenant_1, 1 for tenant_2)
INSERT INTO sites (tenant_id, site_name, location, project_id, site_manager, start_date, expected_end_date, current_status, site_area_sqm, workforce_count, created_at, updated_at) VALUES
('tenant_1', 'Downtown Project', 'New York, NY', 'PRJ-001', 'John', DATE_SUB(NOW(), INTERVAL 60 DAY), DATE_ADD(NOW(), INTERVAL 300 DAY), 'active', 5000.00, 45, NOW(), NOW()),
('tenant_1', 'Harbor Site', 'Boston, MA', 'PRJ-002', 'Sarah', DATE_SUB(NOW(), INTERVAL 120 DAY), DATE_ADD(NOW(), INTERVAL 240 DAY), 'active', 8500.00, 60, NOW(), NOW()),
('tenant_1', 'Industrial Zone', 'Chicago, IL', 'PRJ-003', 'Mike', DATE_SUB(NOW(), INTERVAL 30 DAY), DATE_ADD(NOW(), INTERVAL 520 DAY), 'planning', 12000.00, 0, NOW(), NOW()),
('tenant_2', 'Tech Park', 'Austin, TX', 'PRJ-004', 'Alice', DATE_SUB(NOW(), INTERVAL 90 DAY), DATE_ADD(NOW(), INTERVAL 450 DAY), 'active', 15000.00, 80, NOW(), NOW());

-- Safety Incidents: 4 records
INSERT INTO safety_incidents (tenant_id, site_id, incident_type, severity, incident_date, description, reported_by, status, incident_number, created_at, updated_at) VALUES
('tenant_1', 1, 'accident', 'high', DATE_SUB(NOW(), INTERVAL 45 DAY), 'Fall', 'John', 'resolved', 'INC-001', NOW(), NOW()),
('tenant_1', 1, 'near_miss', 'medium', DATE_SUB(NOW(), INTERVAL 30 DAY), 'Equipment', 'Jane', 'investigating', 'INC-002', NOW(), NOW()),
('tenant_1', 2, 'accident', 'critical', DATE_SUB(NOW(), INTERVAL 15 DAY), 'Serious', 'Mike', 'investigating', 'INC-003', NOW(), NOW()),
('tenant_2', 4, 'near_miss', 'low', DATE_SUB(NOW(), INTERVAL 5 DAY), 'Minor', 'Alex', 'closed', 'INC-004', NOW(), NOW());

-- Compliance Records
INSERT INTO compliance_records (tenant_id, site_id, compliance_type, requirement, due_date, status, last_audit_date, audit_result, notes, created_at, updated_at) VALUES
('tenant_1', 1, 'osha', 'OSHA Standards', DATE_ADD(NOW(), INTERVAL 60 DAY), 'compliant', DATE_SUB(NOW(), INTERVAL 30 DAY), 'pass', 'OK', NOW(), NOW()),
('tenant_1', 1, 'environmental', 'EPA Rules', DATE_ADD(NOW(), INTERVAL 45 DAY), 'compliant', DATE_SUB(NOW(), INTERVAL 45 DAY), 'pass', 'OK', NOW(), NOW()),
('tenant_1', 2, 'safety', 'ISO 45001', DATE_ADD(NOW(), INTERVAL 15 DAY), 'non_compliant', DATE_SUB(NOW(), INTERVAL 90 DAY), 'fail', 'Needs work', NOW(), NOW()),
('tenant_1', 3, 'osha', 'OSHA Standards', DATE_ADD(NOW(), INTERVAL 70 DAY), 'compliant', DATE_SUB(NOW(), INTERVAL 20 DAY), 'pass', 'OK', NOW(), NOW());

-- Permits
INSERT INTO permits (tenant_id, site_id, permit_type, permit_number, issuing_authority, issue_date, expiry_date, status, created_at, updated_at) VALUES
('tenant_1', 1, 'building', 'BLD-001', 'NYC', DATE_SUB(NOW(), INTERVAL 60 DAY), DATE_ADD(NOW(), INTERVAL 304 DAY), 'active', NOW(), NOW()),
('tenant_1', 1, 'environmental', 'ENV-001', 'EPA', DATE_SUB(NOW(), INTERVAL 45 DAY), DATE_ADD(NOW(), INTERVAL 319 DAY), 'active', NOW(), NOW()),
('tenant_1', 2, 'building', 'BLD-002', 'Boston', DATE_SUB(NOW(), INTERVAL 30 DAY), DATE_ADD(NOW(), INTERVAL 10 DAY), 'active', NOW(), NOW()),
('tenant_1', 2, 'electrical', 'ELE-001', 'Boston Electric', DATE_SUB(NOW(), INTERVAL 90 DAY), DATE_SUB(NOW(), INTERVAL 10 DAY), 'expired', NOW(), NOW());
