# Database Migration Index - VYOMTECH ERP

**Last Updated**: December 3, 2025  
**Total Migrations**: 13  
**Total Tables**: 90+  
**Database**: MySQL 8.0

---

## Migration Overview

| # | Migration Name | Tables | Key Purpose | Size |
|---|---|---|---|---|
| 001 | Foundation | 7 | Multi-tenant core infrastructure, users, authentication | 5.8K |
| 002 | Civil Engineering | 4 | Sites, safety incidents, compliance, permits | 4.4K |
| 003 | Construction | 5 | Projects, BOQ, progress tracking, equipment | 5.7K |
| 004 | HR & Payroll | 5 | Employees, attendance, leaves, payroll | 6.5K |
| 005 | Accounts & GL | 7 | Chart of accounts, journals, GL, financial reports | 9.1K |
| 006 | Purchase | 7 | Vendors, POs, GRNs, requisitions, receiving | 8.7K |
| 007 | Sales | 7 | Leads, customers, quotations, orders, invoices | 9.4K |
| 008 | Real Estate | 7 | Properties, blocks, units, cost sheets, bookings | 8.1K |
| 009 | Call Center & AI | 7 | Agents, calls, campaigns, AI models, communications | 7.9K |
| 010 | RBAC | 6 | Roles, permissions, access control, audit logs | 5.4K |
| 011 | Compliance & Tax | 6 | Compliance records, tax calculations, documents | 5.9K |
| 012 | Analytics & Billing | 8 | Dashboards, subscriptions, invoices, gamification | 8.0K |
| 013 | HR Compliance (ESI/PF) | 9 | EPF/ESI registration, contributions, claims, statutory | 11K |

---

## Detailed Migration Structure

### 001_foundation.sql (7 Tables)
**Purpose**: Multi-tenant base infrastructure and authentication  
**Tables**:
- `tenant` - Organization/tenant master
- `user` - User accounts with multi-tenant support
- `password_reset_token` - Password reset management
- `team` - Team/department structure
- `system_config` - Configuration settings
- `auth_token` - API token management
- `audit_log` - Change audit trail

**Key Features**:
- Multi-tenant isolation via `tenant_id`
- Soft deletes support
- Audit timestamps (created_at, updated_at)
- Indexed for performance

---

### 002_civil.sql (4 Tables)
**Purpose**: Civil engineering and construction site management  
**Tables**:
- `sites` - Construction sites
- `safety_incidents` - Safety incident tracking
- `compliance_records` - Compliance tracking
- `permits` - Permit management

**Key Features**:
- Site manager assignment
- Incident severity tracking
- Permit expiry monitoring
- Multi-tenant filtering

---

### 003_construction.sql (5 Tables)
**Purpose**: Construction project and progress management  
**Tables**:
- `construction_projects` - Project master
- `bill_of_quantities` - BOQ items
- `progress_tracking` - Daily/weekly progress
- `quality_control` - Quality checks
- `construction_equipment` - Equipment inventory

**Key Features**:
- Project status workflow
- BOQ cost calculations
- Quality metrics tracking
- Equipment utilization

---

### 004_hr_payroll.sql (5 Tables)
**Purpose**: Human Resources and payroll management  
**Tables**:
- `employee` - Employee master with bank details
- `attendance` - Daily attendance tracking
- `leave_type` - Leave type definitions
- `leave_request` - Leave applications
- `payroll_record` - Monthly payroll processing

**Key Features**:
- Bank account management
- Salary structure (earnings & deductions)
- Leave management
- ESI & EPF deduction tracking in payroll
- **Note**: Basic ESI/EPF tracked here, detailed tracking in migration 013

---

### 005_accounts_gl.sql (7 Tables)
**Purpose**: Accounting and General Ledger operations  
**Tables**:
- `chart_of_account` - Chart of accounts
- `financial_period` - Accounting periods
- `journal_entry` - Journal entries
- `journal_entry_detail` - Debit/credit lines
- `gl_account_balance` - Account balances by period
- `trial_balance` - Trial balance reports
- `income_statement` - P&L statements
- `balance_sheet` - Balance sheet

**Key Features**:
- Hierarchical chart of accounts
- Multi-period GL tracking
- Automated balance calculations
- Financial statement generation

---

### 006_purchase.sql (7 Tables)
**Purpose**: Vendor and purchase management  
**Tables**:
- `vendor` - Vendor master
- `vendor_contact` - Vendor contacts
- `vendor_address` - Vendor addresses
- `purchase_requisition` - PR creation
- `purchase_order` - PO processing
- `po_line_item` - PO details
- `goods_receipt` - GRN processing
- `grn_line_item` - GRN details

**Key Features**:
- Multi-vendor management
- Requisition to PO workflow
- Quality acceptance tracking
- Rejection handling

---

### 007_sales.sql (7 Tables)
**Purpose**: Sales and customer management  
**Tables**:
- `sales_lead` - Sales leads
- `sales_customer` - Customer master
- `sales_quotation` - Quotations
- `sales_quotation_item` - Quotation items
- `sales_order` - Sales orders
- `sales_order_item` - Order items
- `sales_invoice` - Invoices

**Key Features**:
- Lead to customer conversion
- Quote to order workflow
- Invoice generation
- Credit limit management
- Multi-currency support ready

---

### 008_real_estate.sql (7 Tables)
**Purpose**: Real estate project and property management  
**Tables**:
- `property_project` - Real estate projects
- `property_block` - Project blocks/wings
- `property_unit` - Individual units/apartments
- `unit_cost_sheet` - Unit pricing
- `property_booking` - Unit bookings
- `payment_plan` - Payment schedules
- `installment` - Installment tracking

**Key Features**:
- Project-block-unit hierarchy
- Multi-dimensional pricing (carpet, super built-up, etc.)
- Payment plan flexibility
- Installment tracking with status

---

### 009_call_center_ai.sql (7 Tables)
**Purpose**: Call center operations and AI integration  
**Tables**:
- `agent` - Call center agents
- `call` - Call records
- `call_log` - Call event logs
- `campaign` - Outbound campaigns
- `ai_model` - AI models
- `ai_interaction` - AI conversation logs
- `communication` - Multi-channel communications
- `agent_performance` - Agent KPIs

**Key Features**:
- Agent skill-based routing ready
- Call recording URL storage
- AI model versioning
- Campaign tracking
- Performance analytics

---

### 010_rbac.sql (6 Tables)
**Purpose**: Role-Based Access Control  
**Tables**:
- `role` - Role definitions
- `permission` - Permissions
- `role_permission` - Role-permission mapping
- `user_role` - User-role assignment
- `resource` - Protected resources
- `access_log` - Access audit trail

**Key Features**:
- Fine-grained permissions
- System roles support
- Time-based role expiry
- Full access audit trail
- Resource protection metadata

---

### 011_compliance_tax.sql (6 Tables)
**Purpose**: Compliance and tax management  
**Tables**:
- `compliance_record` - Compliance items
- `compliance_checklist` - Compliance checklists
- `tax_calculation` - Tax computations
- `audit_trail` - Entity audit trail
- `regulatory_requirement` - Regulatory items
- `document` - Document storage

**Key Features**:
- Compliance tracking by type
- Tax period-based calculations
- Document version control
- Expiry date monitoring
- Complete audit history

---

### 012_analytics_billing_gamification.sql (8 Tables)
**Purpose**: Analytics, billing, and user engagement  
**Tables**:
- `analytics` - Metrics and KPIs
- `dashboard_widget` - Widget configuration
- `billing_subscription` - Subscription plans
- `billing_invoice` - Billing invoices
- `payment` - Payment records
- `gamification_level` - User levels
- `gamification_achievement` - Achievements
- `user_gamification` - User points/levels
- `leaderboard` - Ranking system

**Key Features**:
- Flexible metrics tracking
- Multi-plan subscription support
- Points and achievement system
- Leaderboard by period
- Custom widget configuration

---

### 013_hr_compliance_esipf.sql (9 Tables)
**Purpose**: Detailed ESI and PF (Employee State Insurance & Provident Fund) compliance  
**Tables**:
- `epf_configuration` - EPF rules and rates
- `esi_configuration` - ESI rules and rates
- `employee_epf_registration` - EPF enrollment
- `employee_esi_registration` - ESI enrollment
- `epf_contribution` - Monthly EPF contributions
- `esi_contribution` - Monthly ESI contributions
- `epf_passbook` - EPF member balance
- `esi_claim` - ESI claim processing
- `statutory_compliance_record` - Compliance submission tracking

**Key Features**:
- Configurable contribution rates
- UAN and ESI number management
- Monthly contribution tracking
- Exemption status management
- Passbook balance maintenance
- Claim processing workflow
- Statutory filing tracking (Form 5, etc.)

---

## Dependency Graph

```
001_foundation.sql (base)
│
├── 002_civil.sql
│   └── 003_construction.sql
│
├── 004_hr_payroll.sql
│   └── 013_hr_compliance_esipf.sql
│
├── 005_accounts_gl.sql
│
├── 006_purchase.sql
│
├── 007_sales.sql
│   └── 008_real_estate.sql (can reference sales_customer)
│
├── 009_call_center_ai.sql (references sales_lead)
│
├── 010_rbac.sql (references user table)
│
├── 011_compliance_tax.sql
│
└── 012_analytics_billing_gamification.sql
```

---

## Table Count by Module

| Module | Tables | Purpose |
|--------|--------|---------|
| Foundation | 7 | Core infrastructure |
| Civil & Construction | 9 | Project management |
| HR & Compliance | 14 | Human resources |
| Accounting | 7 | Financial management |
| Purchase | 7 | Vendor management |
| Sales | 7 | Customer management |
| Real Estate | 7 | Property management |
| Call Center & AI | 7 | Operations |
| RBAC & Security | 6 | Access control |
| Compliance & Audit | 6 | Compliance |
| Analytics & Billing | 8 | Business intelligence |
| **TOTAL** | **91** | **Comprehensive ERP** |

---

## Key Database Features

### Multi-Tenancy
- All tables include `tenant_id` foreign key
- Automatic tenant isolation
- Tenant-scoped data

### Audit & Compliance
- Created_at, Updated_at timestamps
- Deleted_at for soft deletes
- Comprehensive audit logs
- Access tracking

### Performance
- Indexed foreign keys
- Indexed search columns
- Unique constraints for duplicates
- Composite keys where needed

### Data Integrity
- Foreign key constraints
- Referential integrity
- Unique constraints
- NOT NULL constraints on critical fields

### Scalability
- VARCHAR(36) for UUID primary keys
- BIGINT for auto-incrementing IDs where needed
- DECIMAL(18,2) for financial data
- JSON columns for flexible data

---

## Migration Execution Order

**Critical**: Migrations MUST be executed in order (001 → 013)

```bash
# Automatic via Docker
docker-compose up mysql -d

# Manual execution
mysql -u callcenter_user -psecure_app_pass callcenter < 001_foundation.sql
mysql -u callcenter_user -psecure_app_pass callcenter < 002_civil.sql
# ... and so on
```

---

## Verification Commands

```sql
-- Count all tables
SELECT COUNT(*) as total_tables FROM information_schema.TABLES 
WHERE TABLE_SCHEMA = 'callcenter';

-- List all tables by size
SELECT TABLE_NAME, ROUND((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024, 2) as size_mb
FROM INFORMATION_SCHEMA.TABLES 
WHERE TABLE_SCHEMA = 'callcenter'
ORDER BY size_mb DESC;

-- Check foreign keys
SELECT * FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
WHERE TABLE_SCHEMA = 'callcenter' AND REFERENCED_TABLE_NAME IS NOT NULL;

-- Verify tenant isolation
SELECT DISTINCT tenant_id FROM tenant;
```

---

## Future Migrations (Beyond 013)

Planned enhancements:
- 014: Inventory & Stock Management
- 015: Quality Assurance & Testing
- 016: Supply Chain & Logistics
- 017: Fixed Assets & Depreciation
- 018: Project Resource Planning
- 019: Time & Attendance (Advanced)
- 020: Multi-currency & Exchange Rates

---

## Docker Integration

**Location**: `docker-compose.yml`  
**Service**: MySQL 8.0  
**Database**: callcenter  
**User**: callcenter_user  

Migration files are automatically loaded on container startup via:
```yaml
volumes:
  - ./migrations/001_foundation.sql:/docker-entrypoint-initdb.d/01-foundation.sql
  - ./migrations/002_civil.sql:/docker-entrypoint-initdb.d/02-civil.sql
  # ... (all 13 migrations)
```

---

## Notes

- ✅ All 13 migrations completed
- ✅ 90+ tables created across all modules
- ✅ Foreign key constraints in place
- ✅ Multi-tenant isolation implemented
- ✅ ESI & PF compliance fully implemented
- ✅ RBAC ready for implementation
- ⏳ Ready for backend API integration
- ⏳ Ready for frontend connection

---

**Status**: ✅ **PRODUCTION READY**  
**Last Validated**: December 3, 2025
