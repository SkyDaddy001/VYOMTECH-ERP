-- Create Modules table
CREATE TABLE IF NOT EXISTS modules (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    version VARCHAR(20) NOT NULL,
    pricing_model VARCHAR(50) NOT NULL,
    base_cost DECIMAL(10, 2) DEFAULT 0,
    cost_per_user DECIMAL(10, 2) DEFAULT 0,
    cost_per_project DECIMAL(10, 2) DEFAULT 0,
    cost_per_company DECIMAL(10, 2) DEFAULT 0,
    max_users INT DEFAULT NULL,
    max_projects INT DEFAULT NULL,
    max_companies INT DEFAULT NULL,
    is_dependent_on JSON DEFAULT NULL,
    is_core BOOLEAN DEFAULT FALSE,
    requires_approval BOOLEAN DEFAULT FALSE,
    trial_days_allowed INT DEFAULT 0,
    features JSON DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Companies table
CREATE TABLE IF NOT EXISTS companies (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    industry_type VARCHAR(100),
    employee_count INT DEFAULT NULL,
    website VARCHAR(255) DEFAULT NULL,
    max_projects INT NOT NULL DEFAULT 10,
    max_users INT NOT NULL DEFAULT 100,
    current_user_count INT DEFAULT 0,
    current_project_count INT DEFAULT 0,
    billing_email VARCHAR(255) NOT NULL,
    billing_address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    UNIQUE KEY unique_tenant_company_name (tenant_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Projects table
CREATE TABLE IF NOT EXISTS projects (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL,
    tenant_id VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    project_type VARCHAR(50) NOT NULL,
    max_users INT NOT NULL DEFAULT 50,
    current_user_count INT DEFAULT 0,
    budget_allocated DECIMAL(12, 2) DEFAULT 0,
    budget_spent DECIMAL(12, 2) DEFAULT 0,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    INDEX idx_company_id (company_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Company Members table
CREATE TABLE IF NOT EXISTS company_members (
    id VARCHAR(50) PRIMARY KEY,
    company_id VARCHAR(50) NOT NULL,
    user_id INT NOT NULL,
    tenant_id VARCHAR(50) NOT NULL,
    role VARCHAR(50) NOT NULL,
    department VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    INDEX idx_company_id (company_id),
    INDEX idx_user_id (user_id),
    INDEX idx_tenant_id (tenant_id),
    UNIQUE KEY unique_company_user (company_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Project Members table
CREATE TABLE IF NOT EXISTS project_members (
    id VARCHAR(50) PRIMARY KEY,
    project_id VARCHAR(50) NOT NULL,
    user_id INT NOT NULL,
    company_id VARCHAR(50) NOT NULL,
    tenant_id VARCHAR(50) NOT NULL,
    role VARCHAR(50) NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    INDEX idx_project_id (project_id),
    INDEX idx_user_id (user_id),
    INDEX idx_company_id (company_id),
    UNIQUE KEY unique_project_user (project_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create User Roles table (for custom roles per tenant)
CREATE TABLE IF NOT EXISTS user_roles (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    permissions JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    INDEX idx_tenant_id (tenant_id),
    UNIQUE KEY unique_tenant_role (tenant_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Module Subscriptions table
CREATE TABLE IF NOT EXISTS module_subscriptions (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    company_id VARCHAR(50) DEFAULT NULL,
    project_id VARCHAR(50) DEFAULT NULL,
    module_id VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    subscription_started TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    subscription_ended TIMESTAMP NULL,
    trial_started_at TIMESTAMP NULL,
    trial_ends_at TIMESTAMP NULL,
    max_users_allowed INT DEFAULT NULL,
    current_user_count INT DEFAULT 0,
    monthly_budget DECIMAL(12, 2) DEFAULT NULL,
    amount_spent_this_month DECIMAL(12, 2) DEFAULT 0,
    configuration JSON DEFAULT NULL,
    is_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL,
    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_module_id (module_id),
    INDEX idx_company_id (company_id),
    INDEX idx_project_id (project_id),
    INDEX idx_status (status),
    UNIQUE KEY unique_tenant_module_scope (tenant_id, company_id, project_id, module_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Module Usage table
CREATE TABLE IF NOT EXISTS module_usage (
    id VARCHAR(50) PRIMARY KEY,
    subscription_id VARCHAR(50) NOT NULL,
    tenant_id VARCHAR(50) NOT NULL,
    company_id VARCHAR(50) DEFAULT NULL,
    project_id VARCHAR(50) DEFAULT NULL,
    module_id VARCHAR(50) NOT NULL,
    user_count INT DEFAULT 0,
    project_count INT DEFAULT 0,
    company_count INT DEFAULT 0,
    custom_metrics JSON DEFAULT NULL,
    usage_date DATE NOT NULL,
    estimated_cost DECIMAL(12, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (subscription_id) REFERENCES module_subscriptions(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE,
    INDEX idx_subscription_id (subscription_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_usage_date (usage_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Module Licenses table
CREATE TABLE IF NOT EXISTS module_licenses (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    license_key VARCHAR(255) NOT NULL UNIQUE,
    license_type VARCHAR(50) NOT NULL,
    max_companies INT NOT NULL,
    max_projects_per_company INT NOT NULL,
    max_users_per_project INT NOT NULL,
    total_max_users INT NOT NULL,
    enabled_modules JSON NOT NULL,
    disabled_modules JSON DEFAULT NULL,
    expires_at TIMESTAMP NOT NULL,
    issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    issued_by VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Billing table
CREATE TABLE IF NOT EXISTS billing (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL UNIQUE,
    billing_email VARCHAR(255) NOT NULL,
    billing_cycle VARCHAR(50) NOT NULL,
    next_billing_date TIMESTAMP NOT NULL,
    automatic_payment BOOLEAN DEFAULT TRUE,
    payment_method_id VARCHAR(255) DEFAULT NULL,
    tax_rate DECIMAL(5, 2) DEFAULT 0,
    tax_id VARCHAR(100) DEFAULT NULL,
    billing_address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Pricing Plans table
CREATE TABLE IF NOT EXISTS pricing_plans (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    monthly_price DECIMAL(10, 2) NOT NULL,
    annual_price DECIMAL(10, 2) NOT NULL,
    max_users INT NOT NULL,
    max_companies INT NOT NULL,
    max_projects_per_company INT NOT NULL,
    included_modules JSON NOT NULL,
    additional_modules JSON DEFAULT NULL,
    features JSON NOT NULL,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_is_active (is_active),
    INDEX idx_sort_order (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Tenant Plan Subscriptions table
CREATE TABLE IF NOT EXISTS tenant_plan_subscriptions (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    pricing_plan_id VARCHAR(50) NOT NULL,
    start_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    next_billing_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    billing_cycle VARCHAR(50) NOT NULL,
    is_auto_renew BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (pricing_plan_id) REFERENCES pricing_plans(id) ON DELETE RESTRICT,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Invoices table
CREATE TABLE IF NOT EXISTS invoices (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    billing_period_start DATE NOT NULL,
    billing_period_end DATE NOT NULL,
    subtotal_amount DECIMAL(12, 2) NOT NULL,
    tax_amount DECIMAL(12, 2) DEFAULT 0,
    discount_amount DECIMAL(12, 2) DEFAULT 0,
    total_amount DECIMAL(12, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    payment_method VARCHAR(50),
    issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_at TIMESTAMP NOT NULL,
    paid_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_issued_at (issued_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Invoice Line Items table
CREATE TABLE IF NOT EXISTS invoice_line_items (
    id VARCHAR(50) PRIMARY KEY,
    invoice_id VARCHAR(50) NOT NULL,
    module_id VARCHAR(50) DEFAULT NULL,
    description VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(10, 2) NOT NULL,
    total_price DECIMAL(12, 2) NOT NULL,
    tax_rate DECIMAL(5, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE SET NULL,
    INDEX idx_invoice_id (invoice_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create Usage Metrics table
CREATE TABLE IF NOT EXISTS usage_metrics (
    id VARCHAR(50) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    company_id VARCHAR(50) DEFAULT NULL,
    project_id VARCHAR(50) DEFAULT NULL,
    date DATE NOT NULL,
    active_users INT DEFAULT 0,
    new_users INT DEFAULT 0,
    api_calls_used INT DEFAULT 0,
    storage_used_mb DECIMAL(12, 2) DEFAULT 0,
    module_usage_data JSON DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_date (date),
    UNIQUE KEY unique_tenant_scope_date (tenant_id, company_id, project_id, date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
