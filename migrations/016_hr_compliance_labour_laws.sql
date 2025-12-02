-- ============================================================================
-- HR COMPLIANCE MODULE - Labour Laws Compliance
-- ============================================================================
-- Covers: ESI, EPF, PF, Labour Laws, Statutory Compliance, Leaves, etc.

CREATE TABLE IF NOT EXISTS hr_compliance_rules (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Rule Configuration
    rule_name VARCHAR(255) NOT NULL,
    rule_type ENUM(
        'ESI_Threshold', 'EPF_Threshold', 'Bonus_Calculation', 'Gratuity_Calculation',
        'Leave_Accrual', 'Overtime', 'WorkingHours', 'Weekly_Off', 'Attendance', 'Other'
    ) NOT NULL,
    
    rule_description TEXT,
    
    -- Applicable Parameters
    applicable_from TIMESTAMP NOT NULL,
    applicable_to TIMESTAMP,
    
    -- Rule Values (JSON for flexibility)
    rule_parameters JSON, -- e.g., {"wage_limit": 21000, "applicable": true}
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by VARCHAR(36),
    
    INDEX idx_tenant_rule_type (tenant_id, rule_type),
    INDEX idx_rule_active (is_active),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- ESI (EMPLOYEE STATE INSURANCE) COMPLIANCE
-- ============================================================================

CREATE TABLE IF NOT EXISTS esi_compliance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    
    -- ESI Registration
    esi_number VARCHAR(50) UNIQUE NOT NULL,
    esi_registration_date TIMESTAMP,
    esi_office VARCHAR(255),
    esi_state VARCHAR(100),
    
    -- Contribution Details
    effective_from TIMESTAMP NOT NULL,
    effective_to TIMESTAMP,
    
    -- Monthly Contribution (Usually 0.75% employee, 3.25% employer on wages up to ₹21,000)
    employee_contribution_rate DECIMAL(5, 2) DEFAULT 0.75,
    employer_contribution_rate DECIMAL(5, 2) DEFAULT 3.25,
    wage_limit DECIMAL(15, 2) DEFAULT 21000,
    
    -- Benefits Availed
    sick_leave_balance INT DEFAULT 0,
    disability_benefit_availed BOOLEAN DEFAULT false,
    medical_benefit_availed BOOLEAN DEFAULT false,
    
    -- Compliance Status
    is_compliant BOOLEAN DEFAULT true,
    non_compliance_reason TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_employee (tenant_id, employee_id),
    INDEX idx_esi_number (esi_number),
    INDEX idx_compliance_status (is_compliant),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- EPF (EMPLOYEE PROVIDENT FUND) & PF COMPLIANCE
-- ============================================================================

CREATE TABLE IF NOT EXISTS epf_compliance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    
    -- PF Account Details
    pf_account_number VARCHAR(50) UNIQUE NOT NULL,
    pf_registration_date TIMESTAMP,
    pf_office VARCHAR(255),
    pf_state VARCHAR(100),
    
    -- Employee & Employer Details
    effective_from TIMESTAMP NOT NULL,
    effective_to TIMESTAMP,
    
    -- Contribution Rates
    employee_contribution_rate DECIMAL(5, 2) DEFAULT 12, -- 12% of basic salary
    employer_contribution_rate DECIMAL(5, 2) DEFAULT 12, -- 12% of basic salary
    vps_contribution_rate DECIMAL(5, 2) DEFAULT 0,
    
    -- Non-Exempt / Exempt Status
    pf_exempt BOOLEAN DEFAULT false,
    exempt_reason VARCHAR(255),
    
    -- Account Details
    accumulated_balance DECIMAL(15, 2) DEFAULT 0,
    current_balance DECIMAL(15, 2) DEFAULT 0,
    
    -- Withdrawal Status
    partial_withdrawal_allowed BOOLEAN DEFAULT true,
    final_settlement_requested BOOLEAN DEFAULT false,
    settlement_date TIMESTAMP NULL,
    settlement_amount DECIMAL(15, 2) DEFAULT 0,
    
    -- Compliance
    is_compliant BOOLEAN DEFAULT true,
    non_compliance_reason TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_employee (tenant_id, employee_id),
    INDEX idx_pf_account (pf_account_number),
    INDEX idx_pf_exempt (pf_exempt),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- PROFESSIONAL TAX COMPLIANCE (State-wise)
-- ============================================================================

CREATE TABLE IF NOT EXISTS professional_tax_compliance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    
    -- PT Registration
    pt_number VARCHAR(50),
    pt_registration_state VARCHAR(100) NOT NULL,
    pt_registration_date TIMESTAMP,
    
    -- PT Slab Configuration (State-wise)
    monthly_salary_threshold DECIMAL(15, 2),
    pt_amount DECIMAL(10, 2),
    applicable_from TIMESTAMP NOT NULL,
    applicable_to TIMESTAMP,
    
    -- Compliance Status
    is_compliant BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_employee (tenant_id, employee_id),
    INDEX idx_pt_state (pt_registration_state),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- GRATUITY COMPLIANCE (Payment of Gratuity Act, 1972)
-- ============================================================================

CREATE TABLE IF NOT EXISTS gratuity_compliance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    
    -- Gratuity Eligibility
    joining_date TIMESTAMP NOT NULL,
    gratuity_eligible BOOLEAN DEFAULT false,
    years_of_service INT DEFAULT 0,
    
    -- Gratuity Calculation
    -- 15 days' wages per year of service (for first 5 years) + 30 days for subsequent years
    -- Applicable only if service > 5 years
    last_drawn_salary DECIMAL(15, 2),
    gratuity_accrued DECIMAL(15, 2) DEFAULT 0,
    gratuity_liability DECIMAL(15, 2) DEFAULT 0,
    
    -- Payment Details
    gratuity_paid BOOLEAN DEFAULT false,
    gratuity_paid_date TIMESTAMP,
    gratuity_paid_amount DECIMAL(15, 2),
    
    -- Gratuity Fund / Insurance
    fund_type ENUM('Gratuity_Fund', 'Insurance_Policy', 'Direct_Payment') DEFAULT 'Gratuity_Fund',
    policy_number VARCHAR(100),
    fund_balance DECIMAL(15, 2) DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_employee (tenant_id, employee_id),
    INDEX idx_gratuity_eligible (gratuity_eligible),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- BONUS & VARIABLE PAY COMPLIANCE
-- ============================================================================
-- Diwali Bonus, Festival Bonus, Annual Bonus, etc.

CREATE TABLE IF NOT EXISTS bonus_compliance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    
    -- Bonus Configuration
    bonus_type ENUM('Festival_Bonus', 'Annual_Bonus', 'Performance_Bonus', 'Diwali_Bonus', 'Other') NOT NULL,
    bonus_name VARCHAR(255),
    
    -- Eligibility (Usually 30 days continuous service)
    eligible BOOLEAN DEFAULT true,
    days_worked INT,
    minimum_days_required INT DEFAULT 30,
    
    -- Calculation
    applicable_salary_component ENUM('Basic', 'Basic_DA', 'CTC', 'Gross') DEFAULT 'Basic_DA',
    bonus_percentage DECIMAL(5, 2), -- e.g., 8.33% for one month bonus
    bonus_amount DECIMAL(15, 2),
    
    -- Payment
    bonus_year INT,
    due_date TIMESTAMP,
    payment_date TIMESTAMP,
    payment_status ENUM('Not_Applicable', 'Pending', 'Paid', 'Waived') DEFAULT 'Pending',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_employee (tenant_id, employee_id),
    INDEX idx_bonus_type (bonus_type),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- LEAVE COMPLIANCE (Weekly Holidays, Earned Leave, Sick Leave, etc.)
-- ============================================================================

CREATE TABLE IF NOT EXISTS leave_compliance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    leave_type_id VARCHAR(36) NOT NULL,
    
    -- Fiscal Year
    fiscal_year_from TIMESTAMP NOT NULL,
    fiscal_year_to TIMESTAMP NOT NULL,
    
    -- Entitlement (varies by state and company policy)
    -- Earned Leave: 1 day per month (12 days/year) or as per state law
    -- Sick Leave: 7 days/year
    -- Casual Leave: 5-7 days/year
    annual_entitlement INT NOT NULL,
    
    -- Opening Balance (from previous year - carry forward allowed)
    opening_balance INT DEFAULT 0,
    carry_forward_limit INT DEFAULT 5,
    
    -- Current Year Usage
    utilized INT DEFAULT 0,
    available INT,
    
    -- Encashment (if allowed)
    encashment_allowed BOOLEAN DEFAULT true,
    encashment_amount DECIMAL(15, 2) DEFAULT 0,
    
    -- Compliance Status
    is_compliant BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_employee (tenant_id, employee_id),
    INDEX idx_fiscal_year (fiscal_year_from, fiscal_year_to),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (leave_type_id) REFERENCES leave_types(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- STATUTORY COMPLIANCE AUDIT LOG
-- ============================================================================

CREATE TABLE IF NOT EXISTS statutory_compliance_audit (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    
    -- Compliance Check
    compliance_type VARCHAR(100) NOT NULL, -- 'ESI', 'EPF', 'PT', 'Leave', 'Gratuity', etc.
    compliance_item VARCHAR(255),
    compliance_period TIMESTAMP,
    
    -- Check Result
    is_compliant BOOLEAN NOT NULL,
    violation_found TEXT,
    severity ENUM('Critical', 'High', 'Medium', 'Low') DEFAULT 'Low',
    
    -- Action Taken
    action_required TEXT,
    action_taken TEXT,
    action_date TIMESTAMP,
    action_taken_by VARCHAR(36),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_employee (tenant_id, employee_id),
    INDEX idx_compliance_type (compliance_type),
    INDEX idx_severity (severity),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- LABOUR COMPLIANCE DOCUMENTS
-- ============================================================================

CREATE TABLE IF NOT EXISTS labour_compliance_documents (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    
    -- Document Details
    document_type ENUM(
        'Offer_Letter', 'Appointment_Letter', 'Salary_Slip', 'Form_12AA', 'Form_1',
        'Form_5', 'Gratuity_Certificate', 'ESI_Certificate', 'EPF_Certificate', 'Other'
    ) NOT NULL,
    
    document_name VARCHAR(255) NOT NULL,
    document_date TIMESTAMP,
    document_number VARCHAR(100),
    
    -- Document Details
    issue_date TIMESTAMP,
    expiry_date TIMESTAMP,
    document_path VARCHAR(500), -- Storage path
    
    -- Verification
    verified_by VARCHAR(36),
    verified_date TIMESTAMP,
    verification_status ENUM('Not_Verified', 'Verified', 'Expired', 'Invalid') DEFAULT 'Not_Verified',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_employee (tenant_id, employee_id),
    INDEX idx_document_type (document_type),
    INDEX idx_verification_status (verification_status),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- WORKING HOURS & OVERTIME COMPLIANCE
-- ============================================================================

CREATE TABLE IF NOT EXISTS working_hours_compliance (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_id VARCHAR(36) NOT NULL,
    
    -- Working Hours Configuration
    fiscal_year_from TIMESTAMP NOT NULL,
    fiscal_year_to TIMESTAMP NOT NULL,
    
    -- Normal Working Hours
    normal_working_hours_per_day DECIMAL(5, 2) DEFAULT 8,
    working_days_per_week INT DEFAULT 5,
    
    -- Overtime Policy
    overtime_allowed BOOLEAN DEFAULT true,
    overtime_rate_multiplier DECIMAL(5, 2) DEFAULT 1.5, -- 1.5x for normal hours, 2x for holidays
    max_overtime_hours_per_month INT DEFAULT 50,
    
    -- Actual Hours Worked
    total_hours_worked DECIMAL(10, 2) DEFAULT 0,
    total_overtime_hours DECIMAL(10, 2) DEFAULT 0,
    total_undtime_hours DECIMAL(10, 2) DEFAULT 0, -- Negative OT
    
    -- Compliance Status
    is_compliant BOOLEAN DEFAULT true,
    non_compliance_reason TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_employee (tenant_id, employee_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- INTERNAL COMPLAINT REDRESSAL (POSH Act - Prevention of Sexual Harassment)
-- ============================================================================

CREATE TABLE IF NOT EXISTS internal_complaints (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Complaint Details
    complaint_number VARCHAR(50) UNIQUE NOT NULL,
    complaint_date TIMESTAMP NOT NULL,
    
    -- Complainant & Accused
    complainant_id VARCHAR(36) NOT NULL,
    accused_id VARCHAR(36) NOT NULL,
    
    -- Complaint Details
    complaint_type ENUM('Sexual_Harassment', 'Discrimination', 'Harassment', 'Other') NOT NULL,
    description TEXT NOT NULL,
    
    -- Investigation
    investigation_initiated BOOLEAN DEFAULT false,
    investigation_start_date TIMESTAMP,
    investigation_completion_date TIMESTAMP,
    investigation_report_path VARCHAR(500),
    
    -- Resolution
    status ENUM('Open', 'Under_Investigation', 'Resolved', 'Closed', 'Escalated') DEFAULT 'Open',
    resolution_date TIMESTAMP,
    resolution_action TEXT,
    action_taken_by VARCHAR(36),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_complaint (tenant_id),
    INDEX idx_complainant_accused (complainant_id, accused_id),
    INDEX idx_status (status),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (complainant_id) REFERENCES employees(id),
    FOREIGN KEY (accused_id) REFERENCES employees(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- STATUTORY FORMS & FILINGS (Form 5, Form 1, Form 12AA, etc.)
-- ============================================================================

CREATE TABLE IF NOT EXISTS statutory_forms_filings (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    
    -- Form Details
    form_name VARCHAR(100) NOT NULL, -- 'Form 5', 'Form 1', 'Form 12AA', 'Quarter Return'
    form_type ENUM('EPF_Form', 'ESI_Form', 'Income_Tax_Form', 'Labour_Form', 'Other') NOT NULL,
    
    -- Applicable Period
    applicable_from TIMESTAMP NOT NULL,
    applicable_to TIMESTAMP NOT NULL,
    
    -- Filing Details
    filing_date TIMESTAMP,
    filing_reference_number VARCHAR(100),
    filed_by VARCHAR(36),
    
    -- Submission Status
    submission_status ENUM('Not_Due', 'Draft', 'Submitted', 'Accepted', 'Rejected', 'Pending_Response') DEFAULT 'Draft',
    submission_date TIMESTAMP,
    
    -- Document Path
    form_document_path VARCHAR(500),
    acknowledgement_document_path VARCHAR(500),
    
    -- Notes
    remarks TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_tenant_form (tenant_id, form_type),
    INDEX idx_filing_status (submission_status),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ============================================================================
-- CREATE INDEXES FOR PERFORMANCE
-- ============================================================================

CREATE INDEX idx_esi_compliance_tenant ON esi_compliance(tenant_id);
CREATE INDEX idx_epf_compliance_tenant ON epf_compliance(tenant_id);
CREATE INDEX idx_gratuity_compliance_eligible ON gratuity_compliance(gratuity_eligible);
CREATE INDEX idx_leave_compliance_fiscal ON leave_compliance(fiscal_year_from, fiscal_year_to);
CREATE INDEX idx_statutory_compliance_severity ON statutory_compliance_audit(severity);
CREATE INDEX idx_internal_complaints_status ON internal_complaints(status);
CREATE INDEX idx_statutory_forms_due ON statutory_forms_filings(applicable_to);
