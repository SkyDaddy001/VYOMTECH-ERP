# Modular Monetization System - Quick Reference

## What's Implemented

### 1. Core Models ✅
- **Module** - Feature definitions with pricing models
- **ModuleSubscription** - Module assignments to scope
- **ModuleUsage** - Usage tracking for billing
- **ModuleLicense** - Master tenant licenses

### 2. Organization Models ✅
- **Company** - Multiple companies per tenant
- **Project** - Multiple projects per company
- **CompanyMember** - User-to-company relationships
- **ProjectMember** - User-to-project relationships
- **UserRole** - Custom role definitions

### 3. Billing Models ✅
- **Billing** - Tenant billing configuration
- **PricingPlan** - Predefined packages
- **TenantPlanSubscription** - Plan subscriptions
- **Invoice** - Billing documents
- **InvoiceLineItem** - Invoice details
- **UsageMetrics** - Daily usage tracking

### 4. Database Schema ✅
- 15 new tables created
- Proper foreign key relationships
- Indexes on critical fields
- Support for NULL values (optional associations)

### 5. Services ✅
- **ModuleService** (450+ lines)
  - Module registration and lifecycle
  - Subscription management
  - Usage tracking
  - Cost calculation
  - Dependency checking

- **CompanyService** (350+ lines)
  - Company CRUD
  - Project management
  - Member management
  - Cross-company operations
  - Auto-count updates

- **BillingService** (380+ lines)
  - Pricing plan management
  - Subscription lifecycle
  - Invoice generation
  - Usage metrics
  - Charge calculation

### 6. API Handlers ✅
- **ModuleHandler** - Module endpoints
- **CompanyHandler** - Company/project endpoints
- **BillingHandler** - Billing endpoints

### 7. Documentation ✅
- Comprehensive architecture guide (2,000+ lines)
- API endpoint reference
- Usage examples
- Database schema documentation

---

## Quick Start

### 1. Register a Module
```bash
curl -X POST http://localhost:8080/api/modules/register \
  -H "Content-Type: application/json" \
  -d '{
    "id": "advanced-analytics",
    "name": "Advanced Analytics",
    "category": "analytics",
    "pricing_model": "per_user",
    "cost_per_user": 5,
    "base_cost": 50
  }'
```

### 2. Create a Company
```bash
curl -X POST http://localhost:8080/api/companies/create \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "tenant_123",
    "name": "Sales Department",
    "max_users": 100,
    "max_projects": 10,
    "billing_email": "billing@company.com"
  }'
```

### 3. Create a Project
```bash
curl -X POST http://localhost:8080/api/projects/create \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "comp_123",
    "tenant_id": "tenant_123",
    "name": "Q1 Sales Pipeline",
    "project_type": "sales",
    "max_users": 25,
    "budget_allocated": 5000
  }'
```

### 4. Subscribe to Module
```bash
curl -X POST http://localhost:8080/api/modules/subscribe \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "tenant_123",
    "company_id": "comp_123",
    "module_id": "advanced-analytics",
    "max_users_allowed": 25
  }'
```

### 5. Add User to Company
```bash
curl -X POST http://localhost:8080/api/companies/members/add \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "comp_123",
    "user_id": 123,
    "tenant_id": "tenant_123",
    "role": "manager"
  }'
```

### 6. Add User to Project
```bash
curl -X POST http://localhost:8080/api/projects/members/add \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "proj_123",
    "user_id": 123,
    "company_id": "comp_123",
    "tenant_id": "tenant_123",
    "role": "lead"
  }'
```

### 7. Record Usage
```bash
curl -X POST http://localhost:8080/api/billing/usage/record \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "tenant_123",
    "company_id": "comp_123",
    "active_users": 18,
    "new_users": 2,
    "api_calls_used": 5000
  }'
```

### 8. Get Usage Metrics
```bash
curl -X GET "http://localhost:8080/api/billing/usage?tenant_id=tenant_123&start_date=2025-01-01&end_date=2025-01-31"
```

### 9. List Subscriptions
```bash
curl -X GET "http://localhost:8080/api/modules/subscriptions?tenant_id=tenant_123"
```

### 10. Toggle Module
```bash
curl -X PUT http://localhost:8080/api/modules/toggle \
  -H "Content-Type: application/json" \
  -d '{
    "subscription_id": "modsub_123",
    "enabled": false
  }'
```

---

## Scope Levels

### Tenant-Level Module
```json
{
  "tenant_id": "tenant_123",
  "company_id": null,
  "project_id": null,
  "module_id": "enterprise-sso"
}
```
Effect: Available to entire tenant

### Company-Level Module
```json
{
  "tenant_id": "tenant_123",
  "company_id": "comp_123",
  "project_id": null,
  "module_id": "advanced-analytics"
}
```
Effect: Available to entire company

### Project-Level Module
```json
{
  "tenant_id": "tenant_123",
  "company_id": "comp_123",
  "project_id": "proj_123",
  "module_id": "gamification"
}
```
Effect: Available only to project

---

## Pricing Models

### 1. Free
- No cost regardless of usage
- Good for core features

### 2. Per User
- Cost = base + (cost_per_user × active_users)
- Example: $50 base + $5/user = $50 + (5 × 20) = $150

### 3. Per Project
- Cost = base + (cost_per_project × projects_using_module)
- Example: $100 base + $25/project = $100 + (25 × 3) = $175

### 4. Per Company
- Cost = base + (cost_per_company × companies_using_module)
- Example: $200 base + $50/company = $200 + (50 × 2) = $300

### 5. Flat Rate
- Fixed monthly cost regardless of usage
- Example: $999/month

### 6. Tiered
- Cost varies by usage levels
- Example: $50 (0-50 users), $75 (51-100 users), $100 (100+ users)

---

## Roles & Permissions

### Company Roles
- **owner** - Full control over company
- **admin** - Administrative access
- **manager** - Management capabilities
- **member** - Standard member access
- **viewer** - Read-only access

### Project Roles
- **lead** - Project leadership
- **member** - Active project member
- **viewer** - Read-only access
- **analyst** - Analytics-specific access

### Cross-Scope Roles (Accounts Team)
- **accounts-admin** - All accounting operations across all companies
- **accounts-manager** - Manage billing across companies
- **super-admin** - Complete system access

---

## Key Metrics

### Usage Tracking
- Active users per day
- New users added
- API calls consumed
- Storage used (MB)
- Custom module metrics (JSON)

### Billing Metrics
- Monthly base cost
- Per-unit costs
- Overage charges
- Total monthly bill
- Lifetime value

### Module Metrics
- Adoption rate (% of users)
- Utilization rate (% of capacity)
- Cost per user
- Revenue per module
- Churn rate

---

## Database Tables Summary

| Table | Purpose | Key Fields |
|-------|---------|-----------|
| modules | Module definitions | id, pricing_model, costs |
| module_subscriptions | Active subscriptions | tenant_id, company_id, project_id, module_id |
| module_usage | Usage tracking | subscription_id, user_count, date, cost |
| module_licenses | Tenant licenses | tenant_id, enabled_modules, max_limits |
| companies | Organization units | tenant_id, name, status, limits |
| projects | Work items | company_id, name, budget_allocated |
| company_members | Company users | company_id, user_id, role |
| project_members | Project users | project_id, user_id, role |
| user_roles | Custom roles | tenant_id, name, permissions |
| billing | Billing config | tenant_id, cycle, next_billing_date |
| pricing_plans | Package definitions | name, price, included_modules |
| tenant_plan_subscriptions | Plan subscriptions | tenant_id, plan_id, status |
| invoices | Billing documents | tenant_id, amount, status, paid_at |
| invoice_line_items | Invoice details | invoice_id, module_id, quantity, price |
| usage_metrics | Daily snapshots | tenant_id, date, active_users, api_calls |

---

## Next Steps (Phase 3D)

1. **Frontend Dashboard**
   - Multi-company view
   - Project management UI
   - Team member interface
   - Billing portal

2. **Admin Console**
   - Module marketplace
   - Pricing configuration
   - License management
   - Usage analytics

3. **Reporting**
   - Revenue reports
   - Usage analytics
   - Customer segmentation
   - Churn analysis

4. **Integrations**
   - Stripe/payment processor
   - Email/notification system
   - Data warehouse export
   - Analytics tools

---

## Code Structure

```
internal/
├── models/
│   ├── module.go          (Module, Subscription, Usage, License)
│   ├── company.go         (Company, Project, Members, Roles)
│   ├── billing.go         (Plans, Invoices, Metrics)
│
├── services/
│   ├── module_service.go     (450+ lines, Module operations)
│   ├── company_service.go    (350+ lines, Organization operations)
│   ├── billing_service.go    (380+ lines, Billing operations)
│
├── handlers/
│   ├── module_handler.go     (API endpoints for modules)
│   ├── company_handler.go    (API endpoints for organization)
│   ├── billing_handler.go    (API endpoints for billing)
│
├── migrations/
│   └── 004_modular_monetization_schema.sql (15 tables)
```

---

## File References

| File | Lines | Purpose |
|------|-------|---------|
| models/module.go | 180 | Module data models |
| models/company.go | 100 | Organization data models |
| models/billing.go | 150 | Billing data models |
| services/module_service.go | 450 | Module business logic |
| services/company_service.go | 350 | Company business logic |
| services/billing_service.go | 380 | Billing business logic |
| handlers/module_handler.go | 250 | Module API endpoints |
| handlers/company_handler.go | 350 | Company API endpoints |
| handlers/billing_handler.go | 300 | Billing API endpoints |
| migrations/004_modular_monetization_schema.sql | 450 | Database schema |
| MODULAR_MONETIZATION_GUIDE.md | 2000+ | Comprehensive guide |

---

## Testing Checklist

- [ ] Module registration works
- [ ] Module subscription creation works
- [ ] Company creation works
- [ ] Project creation works
- [ ] User addition to company works
- [ ] User addition to project works
- [ ] Usage metrics recording works
- [ ] Charge calculation works
- [ ] Module toggling works
- [ ] Cross-scope module access works
- [ ] Multi-company isolation works
- [ ] Billing cycle calculations work
- [ ] Invoice generation works
- [ ] User permission aggregation works

---

## Known Limitations & Future Work

### Current Limitations
1. Invoice generation is triggered manually (needs automation)
2. Payment processing not integrated (needs Stripe/PayPal)
3. Usage metrics daily snapshots only (needs real-time tracking)
4. No automatic usage alerts (needs threshold alerts)
5. No usage forecasting (needs ML model)

### Future Enhancements
1. Automated invoice generation based on billing cycles
2. Payment processor integration
3. Real-time usage dashboards
4. Overage warnings and alerts
5. Usage forecasting and recommendations
6. Custom pricing rules engine
7. Volume discount calculations
8. White-label support
9. Multi-currency support
10. Tax compliance automation

