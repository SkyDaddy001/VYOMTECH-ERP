-- ============================================================
-- MIGRATION 027: SALES ANALYTICS & REPORTING VIEWS
-- Date: December 8, 2025
-- Purpose: Create views for sales analysis and reporting (read-only)
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- BOOKING DETAILS VIEW - Complete booking information with unit details
-- ============================================================
CREATE OR REPLACE VIEW v_booking_details AS
SELECT
    b.id,
    b.tenant_id,
    b.booking_code,
    b.booking_date,
    b.allotment_date,
    b.agreement_date,
    b.registration_date,
    b.handover_date,
    b.status,
    u.unit_code,
    u.block,
    u.block_wing,
    u.apt_no,
    u.unit_type,
    u.rera_carpet_area,
    u.sbua,
    u.uds_per_sqft,
    sl.first_name,
    sl.last_name,
    sl.email,
    sl.phone,
    ucs.frc,
    ucs.car_parking_type,
    ucs.car_parking_cost,
    ucs.plc,
    ucs.apartment_cost_excluding_govt,
    rd.gst_applicable,
    rd.gst_percentage,
    rd.gst_cost,
    rd.apartment_cost_including_gst,
    rd.registration_type,
    rd.registration_cost,
    ac.maintenance_charge,
    ac.corpus_charge,
    ac.eb_deposit,
    ac.other_works_charge
FROM booking b
LEFT JOIN unit u ON b.unit_id = u.id
LEFT JOIN sales_lead sl ON b.lead_id = sl.id
LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id
LEFT JOIN registration_details rd ON b.id = rd.booking_id
LEFT JOIN additional_charges ac ON b.id = ac.booking_id;

-- ============================================================
-- MONTHLY SALES ANALYSIS VIEW
-- ============================================================
CREATE OR REPLACE VIEW v_monthly_sales_analysis AS
SELECT
    b.tenant_id,
    DATE_FORMAT(b.booking_date, '%Y-%m-01') as financial_month,
    QUARTER(b.booking_date) as financial_quarter,
    YEAR(b.booking_date) as financial_year,
    COUNT(DISTINCT CASE WHEN b.status = 'active' THEN b.id END) as units_sold,
    COUNT(DISTINCT CASE WHEN u.status = 'available' THEN u.id END) as units_unsold,
    COUNT(DISTINCT u.id) as total_units,
    SUM(u.uds_per_sqft) as total_uds,
    SUM(CASE WHEN u.status = 'available' THEN u.uds_per_sqft ELSE 0 END) as unsold_uds,
    SUM(CASE WHEN b.status = 'active' THEN u.uds_per_sqft ELSE 0 END) as sold_uds,
    SUM(u.sbua) as total_sbua,
    SUM(CASE WHEN u.status = 'available' THEN u.sbua ELSE 0 END) as unsold_sbua,
    SUM(CASE WHEN b.status = 'active' THEN u.sbua ELSE 0 END) as sold_sbua,
    SUM(CASE WHEN b.status = 'active' THEN ucs.apartment_cost_excluding_govt ELSE 0 END) as sold_value,
    SUM(CASE WHEN b.status = 'active' THEN rd.gst_cost ELSE 0 END) as gst_total,
    SUM(CASE WHEN b.status = 'active' THEN (ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)) ELSE 0 END) as sold_value_with_gst,
    SUM(COALESCE(p.amount, 0)) as collections_done,
    SUM(CASE WHEN b.status = 'active' THEN (ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)) ELSE 0 END) - SUM(COALESCE(p.amount, 0)) as pending_due
FROM booking b
LEFT JOIN property_unit u ON b.unit_id = u.id
LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id
LEFT JOIN registration_details rd ON b.id = rd.booking_id
LEFT JOIN payment p ON b.id = p.booking_id
GROUP BY b.tenant_id, DATE_FORMAT(b.booking_date, '%Y-%m-01'), QUARTER(b.booking_date), YEAR(b.booking_date);

-- ============================================================
-- COLLECTION REPORT VIEW - Monthly collections breakdown
-- ============================================================
CREATE OR REPLACE VIEW v_collection_report AS
SELECT
    p.tenant_id,
    DATE_FORMAT(p.payment_date, '%Y-%m-01') as collection_month,
    QUARTER(p.payment_date) as financial_quarter,
    YEAR(p.payment_date) as financial_year,
    SUM(p.amount) as overall_collections,
    SUM(CASE WHEN p.towards LIKE '%apartment%' THEN p.amount ELSE 0 END) as apartment_cost,
    SUM(CASE WHEN p.towards LIKE '%gst%' THEN p.amount ELSE 0 END) as gst_collected,
    SUM(CASE WHEN p.towards LIKE '%tds%' THEN p.amount ELSE 0 END) as tds_collected,
    SUM(CASE WHEN p.towards NOT IN ('apartment', 'gst', 'tds') THEN p.amount ELSE 0 END) as others_collected
FROM payment p
GROUP BY p.tenant_id, DATE_FORMAT(p.payment_date, '%Y-%m-01'), QUARTER(p.payment_date), YEAR(p.payment_date);

-- ============================================================
-- SOLD UNITS VIEW - Individual unit sales tracking
-- ============================================================
CREATE OR REPLACE VIEW v_sold_units_tracking AS
SELECT
    b.id as booking_id,
    b.tenant_id,
    u.id as unit_id,
    u.unit_code,
    u.block,
    u.apt_no,
    b.booking_date as sold_date,
    MONTH(b.booking_date) as sold_month,
    ucs.apartment_cost_excluding_govt as sold_value_excl_tax,
    (ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)) as apartment_cost_with_gst,
    COALESCE(rd.gst_cost, 0) as gst_amount,
    0 as tds_amount,
    COALESCE(rd.registration_cost, 0) as registration_cost,
    0 as modt_amount,
    COALESCE(ac.other_works_charge, 0) as other_works,
    COALESCE(ac.maintenance_charge, 0) as maintenance_charge,
    COALESCE(ac.corpus_charge, 0) as corpus_amount,
    (ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0) + COALESCE(rd.registration_cost, 0) + COALESCE(ac.other_works_charge, 0) + COALESCE(ac.maintenance_charge, 0) + COALESCE(ac.corpus_charge, 0)) as total_receivable
FROM booking b
LEFT JOIN property_unit u ON b.unit_id = u.id
LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id
LEFT JOIN registration_details rd ON b.id = rd.booking_id
LEFT JOIN additional_charges ac ON b.id = ac.booking_id
WHERE b.status = 'active';

-- ============================================================
-- PAYMENT SCHEDULE VIEW - Stage-wise payment tracking
-- ============================================================
CREATE OR REPLACE VIEW v_payment_stage_tracking AS
SELECT
    b.id as booking_id,
    b.tenant_id,
    u.block,
    u.apt_no,
    sl.first_name as customer_name,
    ucs.apartment_cost_excluding_govt as apartment_cost,
    ps.payment_stage,
    ps.construction_stage,
    ps.scheduled_date,
    ps.amount_due,
    COALESCE(SUM(p.amount), 0) as amount_received,
    ps.amount_due - COALESCE(SUM(p.amount), 0) as amount_pending,
    ROUND((COALESCE(SUM(p.amount), 0) / ps.amount_due) * 100, 2) as payment_percentage
FROM booking b
LEFT JOIN property_unit u ON b.unit_id = u.id
LEFT JOIN sales_lead sl ON b.lead_id = sl.id
LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id
LEFT JOIN payment_schedule ps ON b.id = ps.booking_id
LEFT JOIN payment p ON ps.id = p.payment_schedule_id
GROUP BY b.id, b.tenant_id, u.block, u.apt_no, sl.first_name, ucs.apartment_cost_excluding_govt, ps.payment_stage, ps.construction_stage, ps.scheduled_date, ps.amount_due;

-- ============================================================
-- BANK vs OWN PAYMENT VIEW
-- Bank Loan Disbursement vs Customer Own Contributions
-- ============================================================
CREATE OR REPLACE VIEW v_bank_own_payment_analysis AS
SELECT
    b.id as booking_id,
    b.tenant_id,
    sl.first_name as customer_name,
    bl.bank_name,
    bl.sanction_amount as bank_sanctioned,
    COALESCE(bl.disbursed_amount, 0) as bank_loan_disbursed,
    bl.sanction_amount - COALESCE(bl.disbursed_amount, 0) as loan_available_for_disbursement,
    SUM(CASE WHEN p.payment_mode IN ('cash', 'cheque', 'bank_transfer') THEN p.amount ELSE 0 END) as customer_own_paid,
    (ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)) - COALESCE(bl.disbursed_amount, 0) - SUM(CASE WHEN p.payment_mode IN ('cash', 'cheque', 'bank_transfer') THEN p.amount ELSE 0 END) as customer_own_due,
    u.block,
    ucs.apartment_cost_excluding_govt,
    ps.payment_stage,
    ROUND((ps.payment_stage / 13) * 100, 2) as stage_percentage,
    ps.amount_due as stage_due,
    COALESCE(SUM(p.amount), 0) as stage_received,
    CURDATE() as as_on_date
FROM booking b
LEFT JOIN property_unit u ON b.unit_id = u.id
LEFT JOIN sales_lead sl ON b.lead_id = sl.id
LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id
LEFT JOIN registration_details rd ON b.id = rd.booking_id
LEFT JOIN bank_loan bl ON b.id = bl.booking_id
LEFT JOIN payment_schedule ps ON b.id = ps.booking_id
LEFT JOIN payment p ON b.id = p.booking_id
GROUP BY b.id, b.tenant_id, sl.first_name, bl.bank_name, bl.sanction_amount, COALESCE(bl.disbursed_amount, 0), ucs.apartment_cost_excluding_govt, COALESCE(rd.gst_cost, 0), u.block, ps.payment_stage, ps.amount_due;

-- ============================================================
-- AGREEMENT STATUS VIEW
-- ============================================================
CREATE OR REPLACE VIEW v_agreement_status AS
SELECT
    b.id as booking_id,
    b.tenant_id,
    u.unit_code,
    u.block,
    u.apt_no,
    sl.first_name as customer_name,
    CASE WHEN b.agreement_date IS NOT NULL THEN 'Signed' ELSE 'Pending' END as agreement_status,
    b.agreement_date,
    DATEDIFF(CURDATE(), b.booking_date) as days_pending,
    ucs.apartment_cost_excluding_govt as sold_value,
    COALESCE(rd.gst_cost, 0) as gst_amount,
    COALESCE(SUM(p.amount), 0) as collected_amount,
    (ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)) - COALESCE(SUM(p.amount), 0) as pending_due
FROM booking b
LEFT JOIN property_unit u ON b.unit_id = u.id
LEFT JOIN sales_lead sl ON b.lead_id = sl.id
LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id
LEFT JOIN registration_details rd ON b.id = rd.booking_id
LEFT JOIN payment p ON b.id = p.booking_id
GROUP BY b.id, b.tenant_id, u.unit_code, u.block, u.apt_no, sl.first_name, b.agreement_date, b.booking_date, ucs.apartment_cost_excluding_govt, COALESCE(rd.gst_cost, 0);

-- ============================================================
-- DASHBOARD SUMMARY VIEW
-- ============================================================
CREATE OR REPLACE VIEW v_sales_dashboard_summary AS
SELECT
    u.tenant_id,
    CURDATE() as report_date,
    COUNT(DISTINCT u.id) as total_units,
    COUNT(DISTINCT CASE WHEN b.status = 'active' THEN u.id END) as sold_units,
    COUNT(DISTINCT CASE WHEN u.status = 'available' THEN u.id END) as unsold_units,
    SUM(ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)) as total_value,
    SUM(COALESCE(p.amount, 0)) as total_collected,
    SUM(ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)) - SUM(COALESCE(p.amount, 0)) as total_pending,
    SUM(u.uds_per_sqft) as total_uds,
    SUM(u.sbua) as total_sbua,
    COUNT(DISTINCT CASE WHEN b.agreement_date IS NOT NULL THEN b.id END) as agreements_signed,
    COUNT(DISTINCT CASE WHEN b.agreement_date IS NULL THEN b.id END) as agreements_pending,
    ROUND((SUM(COALESCE(p.amount, 0)) / SUM(ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0))) * 100, 2) as collection_percentage,
    ROUND((COUNT(DISTINCT CASE WHEN b.status = 'active' THEN u.id END) / COUNT(DISTINCT u.id)) * 100, 2) as occupancy_percentage
FROM property_unit u
LEFT JOIN booking b ON u.id = b.unit_id
LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id
LEFT JOIN registration_details rd ON b.id = rd.booking_id
LEFT JOIN payment p ON b.id = p.booking_id
GROUP BY u.tenant_id;

SET FOREIGN_KEY_CHECKS = 1;
