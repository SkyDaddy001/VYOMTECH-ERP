-- ============================================================
-- MIGRATION 027: SALES ANALYTICS & REPORTING VIEWS
-- Date: December 23, 2025
-- Purpose: Create views for sales analysis and reporting (read-only)
-- NOTE: Uses actual columns from real_estate schema
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- BOOKING DETAILS VIEW
-- ============================================================
CREATE OR REPLACE VIEW v_booking_details AS
SELECT
    b.id,
    b.tenant_id,
    b.booking_date,
    b.allotment_date,
    b.agreement_date,
    b.booking_status,
    u.id as unit_id,
    u.unit_number,
    u.unit_type,
    u.carpet_area,
    u.sbua,
    u.uds_sqft,
    u.status as unit_status,
    pb.block_name,
    pp.project_name,
    ucs.base_price,
    ucs.frc,
    ucs.car_parking_cost
FROM property_booking b
LEFT JOIN property_unit u ON b.unit_id = u.id
LEFT JOIN property_block pb ON u.block_id = pb.id
LEFT JOIN property_project pp ON u.project_id = pp.id
LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id;

-- ============================================================
-- UNIT AVAILABILITY VIEW
-- ============================================================
CREATE OR REPLACE VIEW v_unit_availability_summary AS
SELECT
    pp.id as project_id,
    pp.project_name,
    pp.project_code,
    COUNT(DISTINCT u.id) as total_units,
    COUNT(DISTINCT CASE WHEN u.status = 'available' THEN u.id END) as available_units,
    COUNT(DISTINCT CASE WHEN u.status = 'booked' THEN u.id END) as booked_units,
    COUNT(DISTINCT CASE WHEN u.status = 'sold' THEN u.id END) as sold_units,
    SUM(u.carpet_area) as total_carpet_area,
    SUM(u.sbua) as total_sbua,
    ROUND(COUNT(DISTINCT CASE WHEN u.status = 'booked' THEN u.id END) / COUNT(DISTINCT u.id) * 100, 2) as booking_percentage
FROM property_project pp
LEFT JOIN property_block pb ON pp.id = pb.project_id
LEFT JOIN property_unit u ON pb.id = u.block_id
GROUP BY pp.id, pp.project_name, pp.project_code;

-- ============================================================
-- BOOKING PIPELINE VIEW
-- ============================================================
CREATE OR REPLACE VIEW v_booking_pipeline AS
SELECT
    b.tenant_id,
    b.id as booking_id,
    u.unit_number,
    u.unit_type,
    pp.project_name,
    pb.block_name,
    b.booking_date,
    b.booking_amount,
    b.booking_status,
    CASE 
        WHEN b.agreement_date IS NOT NULL THEN 'Registered'
        WHEN b.allotment_date IS NOT NULL THEN 'Allotted'
        ELSE 'Booked'
    END as booking_stage,
    DATEDIFF(CURDATE(), b.booking_date) as days_since_booking
FROM property_booking b
LEFT JOIN property_unit u ON b.unit_id = u.id
LEFT JOIN property_block pb ON u.block_id = pb.id
LEFT JOIN property_project pp ON u.project_id = pp.id;

-- ============================================================
-- PROJECT PERFORMANCE VIEW
-- ============================================================
CREATE OR REPLACE VIEW v_project_performance AS
SELECT
    pp.id,
    pp.tenant_id,
    pp.project_name,
    pp.project_code,
    pp.location,
    pp.city,
    pp.total_units,
    pp.launch_date,
    pp.expected_completion,
    pp.status,
    COUNT(DISTINCT u.id) as units_created,
    COUNT(DISTINCT b.id) as units_booked,
    ROUND(COUNT(DISTINCT b.id) / COUNT(DISTINCT u.id) * 100, 2) as booking_percentage,
    COUNT(DISTINCT CASE WHEN b.agreement_date IS NOT NULL THEN b.id END) as registered_count,
    SUM(COALESCE(ucs.base_price, 0)) as total_inventory_value
FROM property_project pp
LEFT JOIN property_block pb ON pp.id = pb.project_id
LEFT JOIN property_unit u ON pb.id = u.block_id
LEFT JOIN property_booking b ON u.id = b.unit_id
LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id
GROUP BY pp.id, pp.tenant_id, pp.project_name, pp.project_code, pp.location, pp.city, pp.total_units, pp.launch_date, pp.expected_completion, pp.status;

-- ============================================================
-- UNIT COST ANALYSIS VIEW
-- ============================================================
CREATE OR REPLACE VIEW v_unit_cost_analysis AS
SELECT
    u.id as unit_id,
    u.unit_number,
    u.unit_type,
    u.carpet_area,
    u.sbua,
    pp.project_name,
    pb.block_name,
    ucs.base_price,
    ucs.frc,
    ucs.car_parking_cost,
    ucs.plc,
    ucs.statutory_charges,
    ucs.legal_charges,
    (ucs.base_price + COALESCE(ucs.frc, 0) + COALESCE(ucs.car_parking_cost, 0) + COALESCE(ucs.plc, 0) + COALESCE(ucs.statutory_charges, 0) + COALESCE(ucs.legal_charges, 0)) as total_cost,
    ROUND(ucs.base_price / u.carpet_area, 2) as rate_per_sqft
FROM property_unit u
LEFT JOIN property_block pb ON u.block_id = pb.id
LEFT JOIN property_project pp ON u.project_id = pp.id
LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id;

-- ============================================================
-- PAYMENT TRACKING VIEW
-- ============================================================
CREATE OR REPLACE VIEW v_payment_tracking AS
SELECT
    b.tenant_id,
    b.id as booking_id,
    u.unit_number,
    pp.project_name,
    b.booking_date,
    b.booking_amount,
    COALESCE(SUM(p.amount_paid), 0) as total_paid,
    b.booking_amount - COALESCE(SUM(p.amount_paid), 0) as outstanding_amount,
    ROUND((COALESCE(SUM(p.amount_paid), 0) / b.booking_amount) * 100, 2) as payment_percentage
FROM property_booking b
LEFT JOIN property_unit u ON b.unit_id = u.id
LEFT JOIN property_project pp ON u.project_id = pp.id
LEFT JOIN installment inst ON b.id = (SELECT booking_id FROM property_booking WHERE id = b.id LIMIT 1)
LEFT JOIN installment p ON inst.id = p.id
GROUP BY b.id, b.tenant_id, u.unit_number, pp.project_name, b.booking_date, b.booking_amount;

SET FOREIGN_KEY_CHECKS = 1;
