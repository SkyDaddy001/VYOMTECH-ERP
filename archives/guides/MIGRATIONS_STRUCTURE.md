# VYOMTECH ERP - Migration Structure

## Overview
Database migrations have been consolidated into **12 module-based migration files** in dependency order. This structure ensures clean, organized, and predictable database initialization.

---

## Consolidated Migration Modules

### 1. **001_core_foundation.sql** ⭐ FOUNDATION
**Merged from:** 001_initial_schema.sql + 003_multi_tenant_users.sql + 005_team_table.sql

**Core Tables:**
- `tenants` - Multi-tenant isolation
- `users` - User management with tenant awareness
- `organizations` - Organization hierarchy
- `auth_tokens` - JWT and OAuth tokens
- `system_config` - Global system configuration
- `team` - Team management

**Dependencies:** None (First to run)

---

### 2. **002_modules_base.sql** 
**Merged from:** 002_civil_schema.sql + 004_construction_schema.sql + 006_gamification_schema.sql + 007_gamification_system.sql + 008_modular_monetization_schema.sql

**Modules:**
- Civil Engineering (sites, safety incidents, compliance)
- Construction (projects, BOQ, progress tracking)
- Gamification (points, badges, achievements, leaderboards)
- Monetization (module management, licensing)

**Dependencies:** 001_core_foundation

---

### 3. **003_features_tasks.sql**
**Merged from:** 009_phase1_features.sql + 010_scheduled_tasks_schema.sql + 011_phase2_tasks_notifications.sql + 018_milestone_tracking_and_reporting.sql

**Features:**
- Agent availability tracking
- Lead scoring system
- Scheduled tasks & reminders
- Task management & notifications
- Milestone tracking
- Communication templates
- Audit trails

**Dependencies:** 001_core_foundation

---

### 4. **004_financial_accounting.sql**
**Merged from:** 013_accounts_gl_schema.sql + 017_project_collection_accounts_rera.sql + 021_tax_compliance_income_tax_gst.sql

**Modules:**
- General Ledger (GL) accounting
- Income Tax compliance
- GST compliance
- Project-specific collection accounts
- RERA regulatory requirements

**Dependencies:** 001_core_foundation + 004_financial_accounting

---

### 5. **005_procurement_inventory.sql**
**Merged from:** 015_purchase_module_schema.sql

**Features:**
- Purchase orders (PO)
- Goods Receipt Note (GRN)
- Material Receipt Note (MRN)
- Vendor management
- Contract management
- Inventory tracking

**Dependencies:** 004_financial_accounting

---

### 6. **006_crm_sales.sql**
**Merged from:** 016_sales_module_schema.sql

**Features:**
- CRM Lead Management
- Sales Orders (SO)
- Order fulfillment
- Customer management
- Sales analytics
- Deal pipeline tracking

**Dependencies:** 004_financial_accounting

---

### 7. **007_hr_payroll.sql**
**Merged from:** 019_hr_compliance_labour_laws.sql + 022_hr_payroll_schema.sql

**Features:**
- Employee management
- Payroll processing
- Labour law compliance (ESI, EPF, PF)
- Leave management
- Statutory compliance
- Salary structures

**Dependencies:** 001_core_foundation

---

### 8. **008_real_estate.sql**
**Merged from:** 020_real_estate_property_management.sql

**Features:**
- Property management
- Area statements
- Cost sheets
- Control sheets
- Customer details
- Payment tracking

**Dependencies:** 001_core_foundation

---

### 9. **009_customization_rbac.sql**
**Merged from:** 014_tenant_customization.sql + 024_comprehensive_customization.sql + 025_roles_permissions.sql

**Features:**
- Role-Based Access Control (RBAC)
- Custom roles & permissions
- Lead source customization
- Milestone customization
- Campaign customization
- Workflow stages

**Dependencies:** 001_core_foundation

---

### 10. **010_partner_system.sql**
**Merged from:** 026_external_partner_system.sql + 027_partner_sources_and_credit_policies.sql + 028_sample_partner_logins.sql

**Features:**
- External partner management
- Partner types (portals, channels, vendors, customers)
- Credit policies (time-based, project-based, campaign-based)
- Partner source tracking
- Sample login credentials

**Dependencies:** 001_core_foundation

---

### 11. **011_analytics_workflows.sql**
**Merged from:** 029_phase3_analytics.sql + 030_phase3_workflows.sql

**Features:**
- Analytics tracking tables
- Workflow automation
- Workflow definitions & triggers
- Action execution tracking
- Workflow scheduling

**Dependencies:** 001_core_foundation

---

### 12. **012_sample_data.sql** 🧪 TEST DATA
**Merged from:** 012_sample_data.sql + 023_comprehensive_test_data.sql

**Content:**
- Sample sites & projects
- Test users & credentials
- Sample data across all modules
- Demo configurations

**Dependencies:** All previous migrations

---

## Execution Order

```
001_core_foundation.sql
    ├─→ 002_modules_base.sql
    ├─→ 003_features_tasks.sql
    ├─→ 004_financial_accounting.sql
    │    ├─→ 005_procurement_inventory.sql
    │    └─→ 006_crm_sales.sql
    ├─→ 007_hr_payroll.sql
    ├─→ 008_real_estate.sql
    ├─→ 009_customization_rbac.sql
    ├─→ 010_partner_system.sql
    ├─→ 011_analytics_workflows.sql
    └─→ 012_sample_data.sql (last)
```

---

## Docker Compose Integration

The `docker-compose.yml` file automatically executes migrations in order:

```yaml
volumes:
  - ./migrations/001_core_foundation.sql:/docker-entrypoint-initdb.d/01-core-foundation.sql
  - ./migrations/002_modules_base.sql:/docker-entrypoint-initdb.d/02-modules-base.sql
  - ./migrations/003_features_tasks.sql:/docker-entrypoint-initdb.d/03-features-tasks.sql
  - ./migrations/004_financial_accounting.sql:/docker-entrypoint-initdb.d/04-financial-accounting.sql
  - ./migrations/005_procurement_inventory.sql:/docker-entrypoint-initdb.d/05-procurement-inventory.sql
  - ./migrations/006_crm_sales.sql:/docker-entrypoint-initdb.d/06-crm-sales.sql
  - ./migrations/007_hr_payroll.sql:/docker-entrypoint-initdb.d/07-hr-payroll.sql
  - ./migrations/008_real_estate.sql:/docker-entrypoint-initdb.d/08-real_estate.sql
  - ./migrations/009_customization_rbac.sql:/docker-entrypoint-initdb.d/09-customization-rbac.sql
  - ./migrations/010_partner_system.sql:/docker-entrypoint-initdb.d/10-partner-system.sql
  - ./migrations/011_analytics_workflows.sql:/docker-entrypoint-initdb.d/11-analytics-workflows.sql
  - ./migrations/012_sample_data.sql:/docker-entrypoint-initdb.d/12-sample-data.sql
```

---

## First-Run Setup

### Start fresh database:
```bash
docker-compose down -v  # Remove all volumes
docker-compose up mysql # Wait for initialization
```

### Access database:
```bash
docker exec -it callcenter-mysql mysql -u callcenter_user -p callcenter
```

### Check migration status:
```bash
docker logs callcenter-mysql
```

---

## Benefits

✅ **Modular Organization** - Related tables grouped by business domain
✅ **Clear Dependencies** - Explicit execution order prevents constraint violations
✅ **Easy Maintenance** - Single file per module simplifies updates
✅ **Faster Initialization** - Consolidated structure optimizes boot time
✅ **Scalability** - New modules can be added without disrupting existing structure
✅ **Testing** - Load specific modules independently for unit testing

---

## Key Database Users

| User | Password | Database | Purpose |
|------|----------|----------|---------|
| `root` | `rootpass` | All | Admin (do not use in app) |
| `callcenter_user` | `secure_app_pass` | `callcenter` | Application access |

---

## Notes

- All migrations use `IF NOT EXISTS` to ensure idempotency
- Multi-tenancy is enforced at the application level
- Foreign key constraints are enabled by default
- Sample data (migration 012) is optional for production
