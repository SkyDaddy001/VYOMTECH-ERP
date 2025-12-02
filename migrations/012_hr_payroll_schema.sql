-- ============================================================================
-- HR & PAYROLL MODULE SCHEMA
-- ============================================================================
-- This migration creates tables for Human Resources and Payroll management

-- Employees Table
CREATE TABLE IF NOT EXISTS employees (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20),
    date_of_birth DATE,
    gender ENUM('M', 'F', 'Other') DEFAULT 'M',
    nationality VARCHAR(50),
    address VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    
    -- Employment Details
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    designation VARCHAR(100) NOT NULL,
    department VARCHAR(100) NOT NULL,
    report_to VARCHAR(36),
    employment_type ENUM('Full-time', 'Part-time', 'Contract', 'Intern') DEFAULT 'Full-time',
    joining_date DATE NOT NULL,
    exit_date DATE,
    status ENUM('active', 'inactive', 'on_leave', 'terminated') DEFAULT 'active',
    
    -- Bank Details
    bank_account_number VARCHAR(50),
    bank_ifsc_code VARCHAR(20),
    bank_name VARCHAR(100),
    account_holder_name VARCHAR(100),
    
    -- Salary Structure
    base_salary DECIMAL(12, 2),
    dearness_allowance DECIMAL(12, 2) DEFAULT 0,
    house_rent_allowance DECIMAL(12, 2) DEFAULT 0,
    special_allowance DECIMAL(12, 2) DEFAULT 0,
    conveyance_allowance DECIMAL(12, 2) DEFAULT 0,
    medical_allowance DECIMAL(12, 2) DEFAULT 0,
    other_allowances DECIMAL(12, 2) DEFAULT 0,
    
    -- Deductions
    epf_deduction DECIMAL(12, 2) DEFAULT 0,
    esi_deduction DECIMAL(12, 2) DEFAULT 0,
    professional_tax DECIMAL(12, 2) DEFAULT 0,
    income_tax DECIMAL(12, 2) DEFAULT 0,
    loan_deduction DECIMAL(12, 2) DEFAULT 0,
    advance_deduction DECIMAL(12, 2) DEFAULT 0,
    other_deductions DECIMAL(12, 2) DEFAULT 0,
    
    -- Meta
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_email (email),
    INDEX idx_employee_id (employee_id),
    INDEX idx_department (department),
    INDEX idx_status (status)
);

-- Attendance Table
CREATE TABLE IF NOT EXISTS attendance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    attendance_date DATE NOT NULL,
    check_in_time TIME,
    check_out_time TIME,
    working_hours DECIMAL(5, 2),
    status ENUM('present', 'absent', 'half_day', 'on_leave', 'holiday') DEFAULT 'absent',
    notes VARCHAR(255),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    UNIQUE KEY unique_attendance (employee_id, attendance_date),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_employee_id (employee_id),
    INDEX idx_attendance_date (attendance_date)
);

-- Leave Types
CREATE TABLE IF NOT EXISTS leave_types (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    leave_type_name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    annual_entitlement INT DEFAULT 0,
    is_paid BOOLEAN DEFAULT TRUE,
    carry_forward_limit INT DEFAULT 5,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    INDEX idx_tenant_id (tenant_id)
);

-- Leave Requests
CREATE TABLE IF NOT EXISTS leave_requests (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    leave_type_id VARCHAR(36) NOT NULL,
    from_date DATE NOT NULL,
    to_date DATE NOT NULL,
    number_of_days INT NOT NULL,
    reason VARCHAR(500),
    status ENUM('pending', 'approved', 'rejected', 'cancelled') DEFAULT 'pending',
    approved_by VARCHAR(36),
    approval_date TIMESTAMP NULL,
    rejection_reason VARCHAR(500),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (leave_type_id) REFERENCES leave_types(id),
    FOREIGN KEY (approved_by) REFERENCES employees(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_employee_id (employee_id),
    INDEX idx_status (status)
);

-- Payroll Records
CREATE TABLE IF NOT EXISTS payroll (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    payroll_month DATE NOT NULL,
    payroll_status ENUM('draft', 'generated', 'approved', 'processed', 'paid') DEFAULT 'draft',
    
    -- Earnings
    basic_salary DECIMAL(12, 2),
    dearness_allowance DECIMAL(12, 2) DEFAULT 0,
    house_rent_allowance DECIMAL(12, 2) DEFAULT 0,
    special_allowance DECIMAL(12, 2) DEFAULT 0,
    conveyance_allowance DECIMAL(12, 2) DEFAULT 0,
    medical_allowance DECIMAL(12, 2) DEFAULT 0,
    other_allowances DECIMAL(12, 2) DEFAULT 0,
    total_earnings DECIMAL(12, 2),
    
    -- Deductions
    epf_deduction DECIMAL(12, 2) DEFAULT 0,
    esi_deduction DECIMAL(12, 2) DEFAULT 0,
    professional_tax DECIMAL(12, 2) DEFAULT 0,
    income_tax DECIMAL(12, 2) DEFAULT 0,
    loan_deduction DECIMAL(12, 2) DEFAULT 0,
    advance_deduction DECIMAL(12, 2) DEFAULT 0,
    other_deductions DECIMAL(12, 2) DEFAULT 0,
    total_deductions DECIMAL(12, 2),
    
    -- Net Salary
    net_salary DECIMAL(12, 2),
    
    -- Additional Info
    working_days INT,
    leave_days INT DEFAULT 0,
    paid_days DECIMAL(5, 2),
    notes VARCHAR(500),
    
    -- Meta
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    processed_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    UNIQUE KEY unique_payroll (tenant_id, employee_id, payroll_month),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_employee_id (employee_id),
    INDEX idx_payroll_month (payroll_month),
    INDEX idx_payroll_status (payroll_status)
);

-- Employee Loans
CREATE TABLE IF NOT EXISTS employee_loans (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    loan_amount DECIMAL(12, 2) NOT NULL,
    loan_type ENUM('salary_advance', 'personal_loan', 'vehicle_loan', 'home_loan', 'other') DEFAULT 'personal_loan',
    interest_rate DECIMAL(5, 2) DEFAULT 0,
    tenure_months INT,
    emi_amount DECIMAL(12, 2),
    loan_status ENUM('approved', 'disbursed', 'active', 'closed') DEFAULT 'approved',
    disbursement_date DATE,
    closure_date DATE,
    remaining_amount DECIMAL(12, 2),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_employee_id (employee_id),
    INDEX idx_loan_status (loan_status)
);

-- Loan Repayment Schedule
CREATE TABLE IF NOT EXISTS loan_repayments (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    loan_id VARCHAR(36) NOT NULL,
    installment_number INT,
    emi_amount DECIMAL(12, 2),
    principal_amount DECIMAL(12, 2),
    interest_amount DECIMAL(12, 2),
    due_date DATE,
    paid_date DATE,
    payment_status ENUM('pending', 'paid', 'overdue') DEFAULT 'pending',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (loan_id) REFERENCES employee_loans(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_loan_id (loan_id),
    INDEX idx_payment_status (payment_status)
);

-- Performance Appraisals
CREATE TABLE IF NOT EXISTS performance_appraisals (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    appraisal_period_from DATE NOT NULL,
    appraisal_period_to DATE NOT NULL,
    appraiser_id VARCHAR(36) NOT NULL,
    overall_rating DECIMAL(3, 1),
    communication_rating DECIMAL(3, 1),
    performance_rating DECIMAL(3, 1),
    attendance_rating DECIMAL(3, 1),
    teamwork_rating DECIMAL(3, 1),
    initiative_rating DECIMAL(3, 1),
    comments VARCHAR(1000),
    appraisal_status ENUM('draft', 'submitted', 'reviewed', 'finalized') DEFAULT 'draft',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (appraiser_id) REFERENCES employees(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_employee_id (employee_id)
);

-- Employee Benefits
CREATE TABLE IF NOT EXISTS employee_benefits (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    benefit_type VARCHAR(100),
    provider VARCHAR(100),
    policy_number VARCHAR(100),
    sum_insured DECIMAL(12, 2),
    effective_date DATE,
    expiry_date DATE,
    status ENUM('active', 'inactive', 'expired') DEFAULT 'active',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_employee_id (employee_id)
);

-- Employee Documents
CREATE TABLE IF NOT EXISTS employee_documents (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    document_type ENUM('aadhar', 'pan', 'passport', 'driving_license', 'bank_passbook', 'degree', 'certificate', 'other') NOT NULL,
    document_number VARCHAR(100),
    issue_date DATE,
    expiry_date DATE,
    document_path VARCHAR(255),
    status ENUM('verified', 'pending', 'rejected') DEFAULT 'pending',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_employee_id (employee_id)
);

-- HR Audit Log
CREATE TABLE IF NOT EXISTS hr_audit_log (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    entity_type VARCHAR(100),
    entity_id VARCHAR(36),
    action VARCHAR(50),
    old_values JSON,
    new_values JSON,
    changed_by VARCHAR(36),
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_entity_type (entity_type),
    INDEX idx_changed_at (changed_at)
);

-- Create indexes for better query performance
CREATE INDEX idx_employees_department ON employees(department);
CREATE INDEX idx_employees_joining_date ON employees(joining_date);
CREATE INDEX idx_attendance_employee_date ON attendance(employee_id, attendance_date);
CREATE INDEX idx_leave_requests_employee ON leave_requests(employee_id);
CREATE INDEX idx_payroll_employee_month ON payroll(employee_id, payroll_month);
