-- ============================================================================
-- REAL ESTATE PROPERTY MANAGEMENT SYSTEM - SAMPLE DATA
-- ============================================================================
-- This file contains comprehensive sample data for testing and demonstration
-- of the Real Estate Property Management module with milestone tracking
-- ============================================================================

-- Disable foreign key checks temporarily to avoid constraint issues during insertion
SET session_replication_role = 'replica';

-- ============================================================================
-- CLEAN EXISTING DATA (Optional - uncomment if needed)
-- ============================================================================
-- DELETE FROM customer_account_ledgers;
-- DELETE FROM payment_schedules;
-- DELETE FROM booking_payments;
-- DELETE FROM property_milestones;
-- DELETE FROM customer_details;
-- DELETE FROM customer_bookings;
-- DELETE FROM unit_cost_sheets;
-- DELETE FROM property_units;
-- DELETE FROM property_blocks;
-- DELETE FROM project_control_sheet;
-- DELETE FROM property_projects;

-- ============================================================================
-- 1. INSERT TENANT DATA (Multi-Tenant System)
-- ============================================================================

INSERT INTO tenants (id, name, domain, status, tier, created_at) 
VALUES 
  (1, 'Mahindra Realty', 'mahindra-realty.com', 'active', 'premium', NOW()),
  (2, 'DLF Limited', 'dlf-limited.com', 'active', 'premium', NOW()),
  (3, 'Godrej Properties', 'godrej-properties.com', 'active', 'standard', NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- 2. INSERT PROPERTY PROJECTS
-- ============================================================================

INSERT INTO property_projects 
(tenant_id, project_name, project_code, location, city, state, country, zip_code, 
 total_units, project_type, developer_name, architect_name, noc_status, noc_date, 
 project_start_date, expected_completion, status, created_at, created_by)
VALUES
-- Tenant 1: Mahindra Realty
(1, 'Mahindra Lifespace - Green Valley', 'MH-GV-001', 'Pune - Hinjewadi', 'Pune', 'Maharashtra', 'India', '411057', 250, 'residential', 'Mahindra Lifespaces', 'Hafeez Contractor', true, '2023-06-15', '2023-01-15', '2026-12-31', 'active', NOW(), 'admin'),
(1, 'Mahindra Lifespace - Urban Edge', 'MH-UE-002', 'Mumbai - Lower Parel', 'Mumbai', 'Maharashtra', 'India', '400013', 180, 'residential', 'Mahindra Lifespaces', 'Vastu Shilpa', true, '2023-08-20', '2023-03-01', '2025-09-30', 'active', NOW(), 'admin'),

-- Tenant 2: DLF Limited
(2, 'DLF The Rise', 'DLF-TR-003', 'Gurgaon - Golf Course Road', 'Gurgaon', 'Haryana', 'India', '122002', 320, 'commercial', 'DLF Limited', 'Studio Lotus', true, '2023-07-10', '2023-02-01', '2027-06-30', 'active', NOW(), 'admin'),
(2, 'DLF Camellias', 'DLF-CM-004', 'Delhi - Sector 26 Dwarka', 'New Delhi', 'Delhi', 'India', '110077', 200, 'residential', 'DLF Limited', 'Morphogenesis', true, '2023-05-15', '2023-01-10', '2025-12-31', 'active', NOW(), 'admin'),

-- Tenant 3: Godrej Properties
(3, 'Godrej Alive - Upcoming', 'GOD-AL-005', 'Bangalore - Sarjapur Road', 'Bangalore', 'Karnataka', 'India', '560034', 150, 'residential', 'Godrej Properties', 'Nila Architects', true, '2024-01-20', '2024-01-01', '2026-06-30', 'planning', NOW(), 'admin');

-- ============================================================================
-- 3. INSERT PROPERTY BLOCKS
-- ============================================================================

INSERT INTO property_blocks 
(tenant_id, project_id, block_name, block_number, total_units, ground_level, block_description, created_at, created_by)
VALUES
-- Mahindra Realty - Project 1 Blocks
(1, 1, 'Tower A', 'A', 60, true, 'Premium North-facing tower with modern amenities', NOW(), 'admin'),
(1, 1, 'Tower B', 'B', 65, false, 'Twin tower with skyline view', NOW(), 'admin'),
(1, 1, 'Tower C', 'C', 70, false, 'Contemporary design tower', NOW(), 'admin'),
(1, 1, 'Tower D', 'D', 55, true, 'Community block with shops on ground floor', NOW(), 'admin'),

-- Mahindra Realty - Project 2 Blocks
(1, 2, 'North Wing', 'N', 90, true, 'North-facing residential block', NOW(), 'admin'),
(1, 2, 'South Wing', 'S', 90, true, 'South-facing residential block', NOW(), 'admin'),

-- DLF Limited - Project 3 Blocks
(2, 3, 'Office Tower 1', 'OT-1', 160, true, 'High-end commercial office space', NOW(), 'admin'),
(2, 3, 'Office Tower 2', 'OT-2', 160, false, 'LEED certified commercial tower', NOW(), 'admin'),

-- DLF Limited - Project 4 Blocks
(2, 4, 'Residential Block 1', 'RB-1', 100, true, 'Luxury apartments', NOW(), 'admin'),
(2, 4, 'Residential Block 2', 'RB-2', 100, true, 'Premium residences', NOW(), 'admin'),

-- Godrej Properties - Project 5 Blocks
(3, 5, 'Phase 1 - Block A', 'P1-A', 75, true, 'First phase development', NOW(), 'admin'),
(3, 5, 'Phase 1 - Block B', 'P1-B', 75, true, 'First phase development', NOW(), 'admin');

-- ============================================================================
-- 4. INSERT PROPERTY UNITS
-- ============================================================================

INSERT INTO property_units 
(tenant_id, project_id, block_id, unit_number, unit_type, unit_size, carpet_area, utility_area, sbua, plinth_area, uds_sqft,
 facing, floor_number, status, created_at, created_by)
VALUES
-- Tower A Units (1BHK)
(1, 1, 1, 'A-101', '1BHK', 'apartment', 500.00, 150.00, 650.00, 750.00, 0.50, 'North', 1, 'available', NOW(), 'admin'),
(1, 1, 1, 'A-102', '1BHK', 'apartment', 500.00, 150.00, 650.00, 750.00, 0.50, 'North', 1, 'available', NOW(), 'admin'),
(1, 1, 1, 'A-103', '1BHK', 'apartment', 510.00, 155.00, 665.00, 765.00, 0.50, 'East', 1, 'available', NOW(), 'admin'),
(1, 1, 1, 'A-201', '1BHK', 'apartment', 500.00, 150.00, 650.00, 750.00, 0.50, 'North', 2, 'booked', NOW(), 'admin'),
(1, 1, 1, 'A-202', '2BHK', 'apartment', 750.00, 225.00, 975.00, 1120.00, 0.75, 'North', 2, 'booked', NOW(), 'admin'),
(1, 1, 1, 'A-203', '2BHK', 'apartment', 750.00, 225.00, 975.00, 1120.00, 0.75, 'East', 2, 'available', NOW(), 'admin'),

-- Tower B Units (2BHK)
(1, 1, 2, 'B-301', '2BHK', 'apartment', 800.00, 240.00, 1040.00, 1200.00, 0.80, 'South', 3, 'available', NOW(), 'admin'),
(1, 1, 2, 'B-302', '2BHK', 'apartment', 800.00, 240.00, 1040.00, 1200.00, 0.80, 'West', 3, 'available', NOW(), 'admin'),
(1, 1, 2, 'B-303', '3BHK', 'apartment', 1200.00, 360.00, 1560.00, 1800.00, 1.20, 'South', 3, 'booked', NOW(), 'admin'),
(1, 1, 2, 'B-401', '2BHK', 'apartment', 800.00, 240.00, 1040.00, 1200.00, 0.80, 'South', 4, 'available', NOW(), 'admin'),
(1, 1, 2, 'B-402', '3BHK', 'apartment', 1200.00, 360.00, 1560.00, 1800.00, 1.20, 'West', 4, 'available', NOW(), 'admin'),

-- Tower C Units (3BHK)
(1, 1, 3, 'C-501', '3BHK', 'apartment', 1250.00, 375.00, 1625.00, 1875.00, 1.25, 'North', 5, 'available', NOW(), 'admin'),
(1, 1, 3, 'C-502', '3BHK', 'apartment', 1250.00, 375.00, 1625.00, 1875.00, 1.25, 'East', 5, 'booked', NOW(), 'admin'),
(1, 1, 3, 'C-503', '4BHK', 'apartment', 1600.00, 480.00, 2080.00, 2400.00, 1.60, 'North', 5, 'available', NOW(), 'admin'),

-- Mahindra Realty Project 2 Units
(1, 2, 5, 'N-101', '2BHK', 'apartment', 750.00, 225.00, 975.00, 1120.00, 0.75, 'North', 1, 'available', NOW(), 'admin'),
(1, 2, 5, 'N-102', '2BHK', 'apartment', 750.00, 225.00, 975.00, 1120.00, 0.75, 'North', 1, 'booked', NOW(), 'admin'),
(1, 2, 5, 'N-201', '3BHK', 'apartment', 1100.00, 330.00, 1430.00, 1650.00, 1.10, 'North', 2, 'available', NOW(), 'admin'),

-- DLF Project 3 Commercial Units
(2, 3, 7, 'OT1-101', 'Office', 'commercial', 1000.00, 200.00, 1200.00, 1400.00, 1.00, 'North', 1, 'available', NOW(), 'admin'),
(2, 3, 7, 'OT1-102', 'Office', 'commercial', 1000.00, 200.00, 1200.00, 1400.00, 1.00, 'South', 1, 'booked', NOW(), 'admin'),
(2, 3, 7, 'OT1-201', 'Office', 'commercial', 1500.00, 300.00, 1800.00, 2100.00, 1.50, 'North', 2, 'available', NOW(), 'admin'),

-- DLF Project 4 Residential
(2, 4, 9, 'RB1-101', '3BHK', 'apartment', 1300.00, 390.00, 1690.00, 1950.00, 1.30, 'North', 1, 'available', NOW(), 'admin'),
(2, 4, 9, 'RB1-102', '3BHK', 'apartment', 1300.00, 390.00, 1690.00, 1950.00, 1.30, 'East', 1, 'available', NOW(), 'admin'),
(2, 4, 9, 'RB1-201', '4BHK', 'apartment', 1800.00, 540.00, 2340.00, 2700.00, 1.80, 'North', 2, 'booked', NOW(), 'admin'),

-- Godrej Properties
(3, 5, 11, 'P1A-101', '2BHK', 'apartment', 700.00, 210.00, 910.00, 1050.00, 0.70, 'North', 1, 'available', NOW(), 'admin'),
(3, 5, 11, 'P1A-102', '3BHK', 'apartment', 1100.00, 330.00, 1430.00, 1650.00, 1.10, 'North', 1, 'available', NOW(), 'admin');

-- ============================================================================
-- 5. INSERT UNIT COST SHEETS
-- ============================================================================

INSERT INTO unit_cost_sheets 
(tenant_id, unit_id, base_price, carpet_area_rate, sbua_rate, parking_type, parking_rate, frc_applicable, frc_amount, 
 statutory_charges, guideline_value_rate, total_project_cost, effective_price, discount_percent, discount_amount, 
 final_price, gst_rate, gst_amount, created_at, created_by)
VALUES
-- Tower A Unit A-101
(1, 1, 3250000.00, 6500.00, 5000.00, 'open', 500000.00, true, 50000.00, 
 150000.00, 4800.00, 3250000.00, 3250000.00, 0.00, 0.00, 3250000.00, 5.00, 162500.00, NOW(), 'admin'),

-- Tower A Unit A-102
(1, 2, 3250000.00, 6500.00, 5000.00, 'open', 500000.00, true, 50000.00, 
 150000.00, 4800.00, 3250000.00, 3250000.00, 0.00, 0.00, 3250000.00, 5.00, 162500.00, NOW(), 'admin'),

-- Tower A Unit A-103
(1, 3, 3300000.00, 6500.00, 5000.00, 'open', 500000.00, true, 55000.00, 
 150000.00, 4800.00, 3300000.00, 3300000.00, 5.00, 165000.00, 3135000.00, 5.00, 156750.00, NOW(), 'admin'),

-- Tower A Unit A-201 (Booked)
(1, 4, 3250000.00, 6500.00, 5000.00, 'covered', 650000.00, true, 50000.00, 
 150000.00, 4800.00, 3250000.00, 3250000.00, 0.00, 0.00, 3250000.00, 5.00, 162500.00, NOW(), 'admin'),

-- Tower A Unit A-202 (Booked)
(1, 5, 4950000.00, 6600.00, 5000.00, 'covered', 650000.00, true, 75000.00, 
 200000.00, 4800.00, 4950000.00, 4950000.00, 0.00, 0.00, 4950000.00, 5.00, 247500.00, NOW(), 'admin'),

-- Tower A Unit A-203
(1, 6, 4950000.00, 6600.00, 5000.00, 'open', 500000.00, true, 75000.00, 
 200000.00, 4800.00, 4950000.00, 4950000.00, 3.00, 148500.00, 4801500.00, 5.00, 240075.00, NOW(), 'admin'),

-- Tower B Unit B-301
(1, 7, 5200000.00, 6500.00, 5000.00, 'covered', 650000.00, true, 80000.00, 
 220000.00, 4800.00, 5200000.00, 5200000.00, 0.00, 0.00, 5200000.00, 5.00, 260000.00, NOW(), 'admin'),

-- Tower B Unit B-302
(1, 8, 5200000.00, 6500.00, 5000.00, 'covered', 650000.00, true, 80000.00, 
 220000.00, 4800.00, 5200000.00, 5200000.00, 0.00, 0.00, 5200000.00, 5.00, 260000.00, NOW(), 'admin'),

-- Tower B Unit B-303 (Booked)
(1, 9, 7800000.00, 6500.00, 5000.00, 'covered', 700000.00, true, 120000.00, 
 300000.00, 4800.00, 7800000.00, 7800000.00, 0.00, 0.00, 7800000.00, 5.00, 390000.00, NOW(), 'admin'),

-- Tower B Unit B-401
(1, 10, 5200000.00, 6500.00, 5000.00, 'open', 500000.00, true, 80000.00, 
 220000.00, 4800.00, 5200000.00, 5200000.00, 2.00, 104000.00, 5096000.00, 5.00, 254800.00, NOW(), 'admin'),

-- DLF Commercial Units
(2, 19, 10000000.00, 10000.00, 8000.00, 'covered', 1000000.00, true, 150000.00, 
 400000.00, 8000.00, 10000000.00, 10000000.00, 0.00, 0.00, 10000000.00, 18.00, 1800000.00, NOW(), 'admin'),

-- DLF Residential
(2, 23, 6500000.00, 5000.00, 4000.00, 'covered', 750000.00, true, 100000.00, 
 250000.00, 5000.00, 6500000.00, 6500000.00, 0.00, 0.00, 6500000.00, 5.00, 325000.00, NOW(), 'admin');

-- ============================================================================
-- 6. INSERT CUSTOMER DETAILS
-- ============================================================================

INSERT INTO customer_details 
(tenant_id, primary_name, primary_email, primary_mobile, primary_aadhar, primary_pan, primary_occupation,
 co_applicant1_name, co_applicant1_mobile, co_applicant1_pan,
 co_applicant2_name, co_applicant2_mobile, co_applicant2_pan,
 bank_name, account_holder_name, account_number, ifsc_code, relationship_manager, created_at, created_by)
VALUES
-- Customer 1: Rajesh Kumar
(1, 'Rajesh Kumar', 'rajesh.kumar@email.com', '+919876543210', '123456789012', 'ABCDE1234K', 'Software Engineer',
 'Priya Kumar', '+919876543211', 'DEFGH5678L',
 NULL, NULL, NULL,
 'HDFC Bank', 'Rajesh Kumar', '0001234567891', 'HDFC0000001', 'Amit Sharma', NOW(), 'admin'),

-- Customer 2: Vikram Singh
(1, 'Vikram Singh', 'vikram.singh@email.com', '+919988776655', '234567890123', 'IJKLM9876N', 'Business Owner',
 'Neha Singh', '+919988776656', 'NOPQR0123O',
 'Arjun Singh', '+919988776657', 'QRSTU4567P',
 'ICIC Bank', 'Vikram Singh', '0002345678902', 'ICIC0000002', 'Ravi Patel', NOW(), 'admin'),

-- Customer 3: Meena Verma
(1, 'Meena Verma', 'meena.verma@email.com', '+919765432109', '345678901234', 'UVWXY6789Q', 'Consultant',
 'Arun Verma', '+919765432110', 'XYZAB0123R',
 NULL, NULL, NULL,
 'Axis Bank', 'Meena Verma', '0003456789013', 'AXIS0000003', 'Deepak Nair', NOW(), 'admin'),

-- Customer 4: Anand Patel
(2, 'Anand Patel', 'anand.patel@email.com', '+919654321098', '456789012345', 'CDEFG2345S', 'Entrepreneur',
 'Sneha Patel', '+919654321099', 'GHIJK6789T',
 NULL, NULL, NULL,
 'Kotak Bank', 'Anand Patel', '0004567890124', 'KOTAK0000004', 'Sanjay Kumar', NOW(), 'admin'),

-- Customer 5: Priya Sharma
(2, 'Priya Sharma', 'priya.sharma@email.com', '+919543210987', '567890123456', 'KLMNO0123U', 'Doctor',
 'Rohit Sharma', '+919543210988', 'PQRST4567V',
 NULL, NULL, NULL,
 'SBI Bank', 'Priya Sharma', '0005678901235', 'SBIN0000005', 'Mahesh Desai', NOW(), 'admin'),

-- Customer 6: Suresh Nair
(3, 'Suresh Nair', 'suresh.nair@email.com', '+919432109876', '678901234567', 'UVWXY8901W', 'Architect',
 'Anjali Nair', '+919432109877', 'ZABCD2345X',
 NULL, NULL, NULL,
 'Yes Bank', 'Suresh Nair', '0006789012346', 'YESB0000006', 'Vikram Singh', NOW(), 'admin');

-- ============================================================================
-- 7. INSERT CUSTOMER BOOKINGS
-- ============================================================================

INSERT INTO customer_bookings 
(tenant_id, unit_id, customer_id, booking_date, booking_value, rate_per_sqft, parking_included, 
 buyer_category, possession_date, expected_completion, status, created_at, created_by)
VALUES
-- Booking 1: Rajesh Kumar - Unit A-201
(1, 4, 1, '2024-11-01', 3250000.00, 6500.00, true, 'end_user', '2025-06-30', '2026-12-31', 'confirmed', NOW(), 'admin'),

-- Booking 2: Vikram Singh - Unit A-202
(1, 5, 2, '2024-10-15', 4950000.00, 6600.00, true, 'end_user', '2025-07-31', '2026-12-31', 'confirmed', NOW(), 'admin'),

-- Booking 3: Meena Verma - Unit B-303
(1, 9, 3, '2024-09-20', 7800000.00, 6500.00, true, 'end_user', '2026-06-30', '2027-06-30', 'confirmed', NOW(), 'admin'),

-- Booking 4: Anand Patel - Unit OT1-102
(2, 20, 4, '2024-08-10', 10000000.00, 10000.00, true, 'investor', '2025-05-31', '2026-06-30', 'confirmed', NOW(), 'admin'),

-- Booking 5: Priya Sharma - Unit RB1-201
(2, 25, 5, '2024-07-05', 6500000.00, 5000.00, true, 'end_user', '2025-09-30', '2027-06-30', 'confirmed', NOW(), 'admin');

-- ============================================================================
-- 8. INSERT BOOKING PAYMENTS
-- ============================================================================

INSERT INTO booking_payments 
(tenant_id, booking_id, payment_date, amount, payment_mode, receipt_number, transaction_id, 
 towards, notes, created_at, created_by)
VALUES
-- Booking 1 Payments (Rajesh Kumar)
(1, 1, '2024-11-01', 650000.00, 'bank_transfer', 'RCP001', 'TXN001', 'booking_amount', 'Initial booking amount', NOW(), 'admin'),
(1, 1, '2024-11-15', 650000.00, 'cheque', 'RCP002', 'CHQ001', 'advance', '1st installment', NOW(), 'admin'),
(1, 1, '2024-12-01', 325000.00, 'dd', 'RCP003', 'DD001', 'installment_1', 'Construction linked payment', NOW(), 'admin'),

-- Booking 2 Payments (Vikram Singh)
(1, 2, '2024-10-15', 990000.00, 'cash', 'RCP004', 'CASH001', 'booking_amount', 'Booking amount paid in cash', NOW(), 'admin'),
(1, 2, '2024-11-01', 990000.00, 'neft', 'RCP005', 'NEFT001', 'advance', 'Advance payment via NEFT', NOW(), 'admin'),
(1, 2, '2024-12-01', 742500.00, 'rtgs', 'RCP006', 'RTGS001', 'installment_1', 'Phase 1 completion payment', NOW(), 'admin'),
(1, 2, '2025-01-15', 742500.00, 'bank_transfer', 'RCP007', 'TXN002', 'installment_2', 'Phase 2 completion payment', NOW(), 'admin'),

-- Booking 3 Payments (Meena Verma)
(1, 3, '2024-09-20', 1560000.00, 'bank_transfer', 'RCP008', 'TXN003', 'booking_amount', 'Booking locked', NOW(), 'admin'),
(1, 3, '2024-10-10', 1560000.00, 'cheque', 'RCP009', 'CHQ002', 'advance', 'Advance payment', NOW(), 'admin'),

-- Booking 4 Payments (Anand Patel)
(2, 4, '2024-08-10', 2000000.00, 'bank_transfer', 'RCP010', 'TXN004', 'booking_amount', 'Commercial booking', NOW(), 'admin'),
(2, 4, '2024-09-01', 2000000.00, 'rtgs', 'RCP011', 'RTGS002', 'advance', 'Advance for commercial property', NOW(), 'admin'),

-- Booking 5 Payments (Priya Sharma)
(2, 5, '2024-07-05', 1300000.00, 'bank_transfer', 'RCP012', 'TXN005', 'booking_amount', 'Booking amount', NOW(), 'admin'),
(2, 5, '2024-08-01', 1300000.00, 'neft', 'RCP013', 'NEFT002', 'advance', 'First advance', NOW(), 'admin'),
(2, 5, '2024-09-15', 975000.00, 'bank_transfer', 'RCP014', 'TXN006', 'installment_1', 'Foundation stage payment', NOW(), 'admin');

-- ============================================================================
-- 9. INSERT PAYMENT SCHEDULES
-- ============================================================================

INSERT INTO payment_schedules 
(tenant_id, booking_id, installment_number, due_date, amount, payment_purpose, status, created_at, created_by)
VALUES
-- Booking 1 Schedule
(1, 1, 1, '2024-11-01', 650000.00, 'Booking Amount', 'completed', NOW(), 'admin'),
(1, 1, 2, '2024-11-15', 650000.00, 'Advance', 'completed', NOW(), 'admin'),
(1, 1, 3, '2024-12-01', 325000.00, 'Foundation Completion', 'completed', NOW(), 'admin'),
(1, 1, 4, '2025-03-01', 325000.00, 'Plinth Completion', 'pending', NOW(), 'admin'),
(1, 1, 5, '2025-06-01', 325000.00, 'Structure Completion', 'pending', NOW(), 'admin'),

-- Booking 2 Schedule
(1, 2, 1, '2024-10-15', 990000.00, 'Booking Amount', 'completed', NOW(), 'admin'),
(1, 2, 2, '2024-11-01', 990000.00, 'Advance', 'completed', NOW(), 'admin'),
(1, 2, 3, '2024-12-01', 742500.00, 'Foundation Completion', 'completed', NOW(), 'admin'),
(1, 2, 4, '2025-03-01', 742500.00, 'Plinth Completion', 'completed', NOW(), 'admin'),
(1, 2, 5, '2025-06-01', 742500.00, 'Structure Completion', 'pending', NOW(), 'admin'),

-- Booking 3 Schedule
(1, 3, 1, '2024-09-20', 1560000.00, 'Booking Amount', 'completed', NOW(), 'admin'),
(1, 3, 2, '2024-10-10', 1560000.00, 'Advance', 'completed', NOW(), 'admin'),
(1, 3, 3, '2024-12-01', 1560000.00, 'Foundation Completion', 'pending', NOW(), 'admin'),
(1, 3, 4, '2025-03-01', 1560000.00, 'Plinth Completion', 'pending', NOW(), 'admin'),

-- Booking 4 Schedule
(2, 4, 1, '2024-08-10', 2000000.00, 'Booking Amount', 'completed', NOW(), 'admin'),
(2, 4, 2, '2024-09-01', 2000000.00, 'Advance', 'completed', NOW(), 'admin'),
(2, 4, 3, '2024-12-01', 2000000.00, 'Framework Completion', 'pending', NOW(), 'admin'),

-- Booking 5 Schedule
(2, 5, 1, '2024-07-05', 1300000.00, 'Booking Amount', 'completed', NOW(), 'admin'),
(2, 5, 2, '2024-08-01', 1300000.00, 'Advance', 'completed', NOW(), 'admin'),
(2, 5, 3, '2024-09-15', 975000.00, 'Foundation Completion', 'completed', NOW(), 'admin'),
(2, 5, 4, '2024-12-01', 975000.00, 'Plinth Completion', 'pending', NOW(), 'admin');

-- ============================================================================
-- 10. INSERT PROPERTY MILESTONES (Campaign Tracking)
-- ============================================================================

INSERT INTO property_milestones 
(tenant_id, booking_id, campaign_name, source, subsource, lead_generated_date, re_engaged_date, 
 site_visit_date, revisit_date, booking_date, cancelled_date, notes, created_at, created_by)
VALUES
-- Milestone 1: Rajesh Kumar
(1, 1, 'Diwali Special Campaign 2024', 'direct', 'reference', '2024-09-15', '2024-10-01', 
 '2024-10-20', '2024-11-01', '2024-11-01', NULL, 'Customer was referred by existing client', NOW(), 'admin'),

-- Milestone 2: Vikram Singh
(1, 2, 'Premium Business Program', 'site_visit', 'hoardings', '2024-08-20', '2024-09-05', 
 '2024-09-15', '2024-09-28', '2024-10-15', NULL, 'Walked-in from site hoarding, converted after revisit', NOW(), 'admin'),

-- Milestone 3: Meena Verma
(1, 3, 'Digital Marketing Push Q3 2024', 'digital', 'facebook', '2024-08-01', '2024-08-25', 
 '2024-09-10', '2024-09-18', '2024-09-20', NULL, 'Led from Facebook campaign, high engagement', NOW(), 'admin'),

-- Milestone 4: Anand Patel
(2, 4, 'Corporate Investor Program', 'broker', 'commercial_broker', '2024-07-15', '2024-07-28', 
 '2024-08-05', '2024-08-08', '2024-08-10', NULL, 'Commercial property broker referral', NOW(), 'admin'),

-- Milestone 5: Priya Sharma
(2, 5, 'Medical Professional Scheme', 'referral', 'colleague', '2024-05-30', '2024-06-15', 
 '2024-06-25', '2024-07-02', '2024-07-05', NULL, 'Referred by hospital colleague', NOW(), 'admin'),

-- Additional Milestone for Lost Customer (Example)
(1, NULL, 'Exhibition Drive November 2024', 'exhibition', 'trade_show', '2024-10-25', '2024-11-05', 
 '2024-11-10', NULL, NULL, '2024-11-25', 'Customer visited exhibition but cancelled after budget concerns', NOW(), 'admin');

-- ============================================================================
-- 11. INSERT ACCOUNT LEDGERS (Auto-Generated on Payment)
-- ============================================================================

INSERT INTO customer_account_ledgers 
(tenant_id, booking_id, transaction_date, transaction_type, opening_balance, credit_amount, debit_amount, 
 closing_balance, reference_number, description, created_at, created_by)
VALUES
-- Ledger Entries for Booking 1 (Rajesh Kumar)
(1, 1, '2024-11-01', 'credit', 0.00, 650000.00, 0.00, 650000.00, 'RCP001', 'Initial booking amount received', NOW(), 'admin'),
(1, 1, '2024-11-15', 'credit', 650000.00, 650000.00, 0.00, 1300000.00, 'RCP002', '1st installment received', NOW(), 'admin'),
(1, 1, '2024-12-01', 'credit', 1300000.00, 325000.00, 0.00, 1625000.00, 'RCP003', 'Construction linked payment received', NOW(), 'admin'),

-- Ledger Entries for Booking 2 (Vikram Singh)
(1, 2, '2024-10-15', 'credit', 0.00, 990000.00, 0.00, 990000.00, 'RCP004', 'Booking amount received', NOW(), 'admin'),
(1, 2, '2024-11-01', 'credit', 990000.00, 990000.00, 0.00, 1980000.00, 'RCP005', 'Advance payment received', NOW(), 'admin'),
(1, 2, '2024-12-01', 'credit', 1980000.00, 742500.00, 0.00, 2722500.00, 'RCP006', 'Phase 1 completion payment', NOW(), 'admin'),
(1, 2, '2025-01-15', 'credit', 2722500.00, 742500.00, 0.00, 3465000.00, 'RCP007', 'Phase 2 completion payment', NOW(), 'admin'),

-- Ledger Entries for Booking 3 (Meena Verma)
(1, 3, '2024-09-20', 'credit', 0.00, 1560000.00, 0.00, 1560000.00, 'RCP008', 'Booking locked', NOW(), 'admin'),
(1, 3, '2024-10-10', 'credit', 1560000.00, 1560000.00, 0.00, 3120000.00, 'RCP009', 'Advance payment', NOW(), 'admin'),

-- Ledger Entries for Booking 4 (Anand Patel)
(2, 4, '2024-08-10', 'credit', 0.00, 2000000.00, 0.00, 2000000.00, 'RCP010', 'Commercial booking amount', NOW(), 'admin'),
(2, 4, '2024-09-01', 'credit', 2000000.00, 2000000.00, 0.00, 4000000.00, 'RCP011', 'Commercial advance', NOW(), 'admin'),

-- Ledger Entries for Booking 5 (Priya Sharma)
(2, 5, '2024-07-05', 'credit', 0.00, 1300000.00, 0.00, 1300000.00, 'RCP012', 'Booking amount', NOW(), 'admin'),
(2, 5, '2024-08-01', 'credit', 1300000.00, 1300000.00, 0.00, 2600000.00, 'RCP013', 'First advance', NOW(), 'admin'),
(2, 5, '2024-09-15', 'credit', 2600000.00, 975000.00, 0.00, 3575000.00, 'RCP014', 'Foundation stage payment', NOW(), 'admin');

-- ============================================================================
-- 12. INSERT PROJECT CONTROL SHEET (Configuration)
-- ============================================================================

INSERT INTO project_control_sheet 
(tenant_id, project_id, parameter_name, parameter_value, effective_date, created_at, created_by)
VALUES
-- Mahindra Realty - Project 1 Control Settings
(1, 1, 'max_discount_percent', '10.00', NOW(), NOW(), 'admin'),
(1, 1, 'escalation_percent_annual', '8.00', NOW(), NOW(), 'admin'),
(1, 1, 'registration_charges_percent', '2.50', NOW(), NOW(), 'admin'),
(1, 1, 'frc_percent_of_price', '1.50', NOW(), NOW(), 'admin'),
(1, 1, 'possession_grace_period_days', '30', NOW(), NOW(), 'admin'),
(1, 1, 'penalty_delay_percent_monthly', '0.50', NOW(), NOW(), 'admin'),
(1, 1, 'maintenance_charge_multiplier', '1.25', NOW(), NOW(), 'admin'),

-- DLF Limited - Project 3 Control Settings
(2, 3, 'max_discount_percent', '15.00', NOW(), NOW(), 'admin'),
(2, 3, 'escalation_percent_annual', '7.50', NOW(), NOW(), 'admin'),
(2, 3, 'registration_charges_percent', '3.00', NOW(), NOW(), 'admin'),
(2, 3, 'frc_percent_of_price', '2.00', NOW(), NOW(), 'admin'),
(2, 3, 'possession_grace_period_days', '45', NOW(), NOW(), 'admin'),
(2, 3, 'penalty_delay_percent_monthly', '0.75', NOW(), NOW(), 'admin'),
(2, 3, 'maintenance_charge_multiplier', '1.50', NOW(), NOW(), 'admin'),

-- Godrej Properties - Project 5 Control Settings
(3, 5, 'max_discount_percent', '5.00', NOW(), NOW(), 'admin'),
(3, 5, 'escalation_percent_annual', '6.00', NOW(), NOW(), 'admin'),
(3, 5, 'registration_charges_percent', '2.00', NOW(), NOW(), 'admin'),
(3, 5, 'frc_percent_of_price', '1.00', NOW(), NOW(), 'admin'),
(3, 5, 'possession_grace_period_days', '60', NOW(), NOW(), 'admin'),
(3, 5, 'penalty_delay_percent_monthly', '0.25', NOW(), NOW(), 'admin'),
(3, 5, 'maintenance_charge_multiplier', '1.00', NOW(), NOW(), 'admin');

-- ============================================================================
-- RE-ENABLE FOREIGN KEY CHECKS
-- ============================================================================

SET session_replication_role = 'default';

-- ============================================================================
-- SUMMARY OF SAMPLE DATA
-- ============================================================================
-- Tenants: 3
-- Projects: 5
-- Blocks: 12
-- Units: 24
-- Cost Sheets: 12
-- Customers: 6
-- Bookings: 5 (All with full payment history)
-- Payments: 13 total transactions
-- Payment Schedules: 21 installments
-- Milestones: 6 (5 booked + 1 cancelled)
-- Ledger Entries: 15 auto-generated ledger transactions
-- Control Settings: 21 configuration parameters
--
-- Key Features Demonstrated:
-- • Multi-tenant isolation (3 distinct tenants)
-- • Complete booking lifecycle (lead → booking → payment → ledger)
-- • All 6 milestone dates captured
-- • 6 campaign sources and multiple subsources
-- • Auto-generated ledger on payment recording
-- • Payment schedule with status tracking
-- • Cost sheet with discount calculations
-- • 5-stage booking progress
-- • Multiple payment modes (cash, cheque, bank transfer, NEFT, RTGS, DD)
-- ============================================================================
