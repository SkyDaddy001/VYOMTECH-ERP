# Modular Monetization System - Phase 3C

## Overview

This document outlines the comprehensive modular monetization system that allows:
- **Modules**: Individual features that can be turned on/off per tenant/company/project
- **Multi-Company/Project Structure**: Each tenant can manage multiple companies, each with multiple projects
- **Monetization**: Per-user, per-project, per-company, or flat-rate pricing models
- **Cross-Company Administration**: Accounts team works across all companies under a tenant
- **Flexible User Roles**: Users can be part of multiple subsets with different roles

---

## Architecture

### 1. Modular System

#### Core Components

**Module** (`models/module.go`)
- Represents a feature module in the system (e.g., "gamification", "workflow", "analytics")
- Each module can be independent or dependent on other modules
- Supports multiple pricing models

**ModuleSubscription**
- Tracks which entities (tenant/company/project) have subscribed to which modules
- Can be enabled/disabled independently
- Tracks usage and costs per subscription

**ModuleUsage**
- Records daily usage metrics for each module subscription
- Used for calculating charges under per-user/per-project pricing models

**ModuleLicense**
- Master license for a tenant
- Defines which modules are enabled/disabled
- Specifies limits (max companies, max projects per company, max users)

#### Module Categories
- **core** - Always available (gamification, leads, agents, campaigns)
- **analytics** - Analytics dashboards
- **automation** - Workflow automation, scheduling
- **communication** - Email, SMS, push notifications
- **ai** - AI-powered features

#### Pricing Models

1. **Free** - No cost
2. **Per User** - Cost multiplied by number of active users
3. **Per Project** - Cost multiplied by number of projects using module
4. **Per Company** - Cost per company subscription
5. **Flat** - Fixed monthly cost
6. **Tiered** - Price tiers based on usage brackets

---

### 2. Multi-Company/Project Structure

#### Company (`models/company.go`)
```
Tenant
  ├── Company A (Sales)
  │   ├── Project 1 (Sales Pipeline)
  │   ├── Project 2 (Lead Scoring)
  │   └── Project 3 (Outreach)
  ├── Company B (Support)
  │   ├── Project 1 (Support Tickets)
  │   └── Project 2 (Customer Feedback)
  └── Company C (Marketing)
      ├── Project 1 (Campaign A)
      └── Project 2 (Campaign B)
```

**Company Structure**
- Each tenant can create multiple companies
- Each company has independent configurations
- Companies can share resources based on settings
- Billing can be aggregated or separated per company

**Project Structure**
- Each company can have multiple projects
- Projects have independent budgets and user allocations
- Users can belong to multiple projects
- Each project can have different module enablement

#### Membership Models

**CompanyMember**
- Links users to companies
- Roles: owner, admin, manager, member, viewer
- Department assignment for organizational structure

**ProjectMember**
- Links users to projects under a company
- Roles: lead, member, viewer, analyst
- Tracks join date and active status

**UserRole** (Custom Roles)
- Tenant-specific role definitions
- Can grant permissions across company/project boundaries
- Enables "Accounts Team" to work cross-company

---

### 3. Monetization Tracking

#### Billing Cycle
- **Monthly** - Monthly recurring billing
- **Quarterly** - Every 3 months
- **Annual** - Once per year
- **Custom** - Custom billing periods

#### PricingPlan
Predefined packages customers can subscribe to:
```
- Startup: $99/month (5 users, 1 company, 5 projects)
  - Includes: Leads, Agents, Campaigns, Calls, Reports
  
- Professional: $499/month (50 users, 3 companies, 25 projects)
  - Includes: Startup + Gamification + Analytics
  
- Enterprise: Custom pricing
  - Includes: All modules + Custom configurations
```

#### Invoice
- Generated at the end of billing cycle
- Contains line items for each module/resource
- Supports tax calculations and discounts
- Tracks payment status and method

#### UsageMetrics
- Daily snapshots of usage per tenant/company/project
- Tracks active users, new users, API calls, storage
- Module-specific metrics in JSON format
- Used for post-billing reconciliation

---

## Database Schema

### Core Tables

```
┌─────────────────────────────────────────────────────────────┐
│ Modules Structure                                           │
├─────────────────────────────────────────────────────────────┤
│ modules                    - Feature definitions
│ module_subscriptions       - Who has which modules
│ module_usage              - Usage tracking per subscription
│ module_licenses           - Master licenses per tenant
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ Organization Structure                                      │
├─────────────────────────────────────────────────────────────┤
│ companies                 - Organizations under tenant
│ projects                  - Projects under company
│ company_members           - User memberships in company
│ project_members           - User memberships in project
│ user_roles               - Custom role definitions
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ Billing Structure                                           │
├─────────────────────────────────────────────────────────────┤
│ billing                   - Tenant billing configuration
│ pricing_plans            - Predefined packages
│ tenant_plan_subscriptions - Tenant plan subscriptions
│ invoices                 - Generated invoices
│ invoice_line_items       - Invoice details
│ usage_metrics            - Daily usage tracking
└─────────────────────────────────────────────────────────────┘
```

---

## API Endpoints

### Module Management

#### Register Module (Admin)
```
POST /api/modules/register
{
  "id": "gamification",
  "name": "Gamification",
  "category": "core",
  "pricing_model": "per_user",
  "base_cost": 0,
  "cost_per_user": 5,
  "is_core": false,
  "trial_days_allowed": 30
}
```

#### List Modules
```
GET /api/modules?status=active
Response: [Module, ...]
```

#### Subscribe to Module
```
POST /api/modules/subscribe
{
  "tenant_id": "tenant_123",
  "company_id": "comp_123",          // Optional
  "project_id": "proj_123",          // Optional
  "module_id": "gamification",
  "max_users_allowed": 100,
  "monthly_budget": 500
}
```

#### Toggle Module
```
PUT /api/modules/toggle
{
  "subscription_id": "modsub_123",
  "enabled": false
}
```

#### Get Module Usage
```
GET /api/modules/usage?subscription_id=modsub_123&start_date=2025-01-01&end_date=2025-01-31
Response: [ModuleUsage, ...]
```

### Company Management

#### Create Company
```
POST /api/companies/create
{
  "tenant_id": "tenant_123",
  "name": "Sales Department",
  "industry_type": "SaaS",
  "max_projects": 10,
  "max_users": 100,
  "billing_email": "billing@company.com"
}
```

#### List Companies
```
GET /api/companies?tenant_id=tenant_123
Response: [Company, ...]
```

#### Create Project
```
POST /api/projects/create
{
  "company_id": "comp_123",
  "tenant_id": "tenant_123",
  "name": "Sales Pipeline",
  "project_type": "sales",
  "max_users": 50,
  "budget_allocated": 5000
}
```

#### List Projects
```
GET /api/projects?company_id=comp_123
Response: [Project, ...]
```

### User Management

#### Add Member to Company
```
POST /api/companies/members/add
{
  "company_id": "comp_123",
  "user_id": 123,
  "tenant_id": "tenant_123",
  "role": "manager",
  "department": "Sales"
}
```

#### Add Member to Project
```
POST /api/projects/members/add
{
  "project_id": "proj_123",
  "user_id": 123,
  "company_id": "comp_123",
  "tenant_id": "tenant_123",
  "role": "lead"
}
```

#### Get Company Members
```
GET /api/companies/members?company_id=comp_123
Response: [CompanyMember, ...]
```

#### Get Project Members
```
GET /api/projects/members?project_id=proj_123
Response: [ProjectMember, ...]
```

#### Remove Project Member
```
DELETE /api/projects/members?project_id=proj_123&user_id=123
```

### Billing Management

#### Create Pricing Plan
```
POST /api/billing/plans/create
{
  "name": "Professional",
  "monthly_price": 499,
  "annual_price": 4990,
  "max_users": 50,
  "max_companies": 3,
  "max_projects_per_company": 10,
  "included_modules": ["leads", "agents", "campaigns", "gamification"],
  "features": ["advanced-analytics", "priority-support"]
}
```

#### List Pricing Plans
```
GET /api/billing/plans
Response: [PricingPlan, ...]
```

#### Subscribe to Plan
```
POST /api/billing/subscribe
{
  "tenant_id": "tenant_123",
  "pricing_plan_id": "plan_456",
  "billing_cycle": "monthly"
}
```

#### Record Usage Metrics
```
POST /api/billing/usage/record
{
  "tenant_id": "tenant_123",
  "company_id": "comp_123",
  "active_users": 42,
  "new_users": 5,
  "api_calls_used": 10500,
  "storage_used_mb": 256.5
}
```

#### Get Usage Metrics
```
GET /api/billing/usage?tenant_id=tenant_123&start_date=2025-01-01&end_date=2025-01-31
Response: [UsageMetrics, ...]
```

#### List Invoices
```
GET /api/billing/invoices?tenant_id=tenant_123&limit=50&offset=0
Response: [Invoice, ...]
```

#### Calculate Monthly Charges
```
GET /api/billing/charges?tenant_id=tenant_123
Response: {"total_monthly_charges": 1245.50}
```

#### Mark Invoice as Paid
```
PUT /api/billing/invoices/mark-paid
{
  "invoice_id": "inv_123"
}
```

---

## Service Layers

### ModuleService (`internal/services/module_service.go`)
- Module registration and management
- Subscription handling
- Usage tracking and calculation
- Dependency checking
- Cost calculation for different pricing models

### CompanyService (`internal/services/company_service.go`)
- Company CRUD operations
- Project management
- User membership (company and project level)
- Cross-company operations (for admin/accounts team)
- Auto-update counts (users, projects)

### BillingService (`internal/services/billing_service.go`)
- Pricing plan management
- Subscription lifecycle
- Invoice generation and management
- Usage metrics recording
- Monthly charge calculation

---

## Key Features

### 1. Module On/Off Control
```go
// Disable gamification for specific project
moduleService.ToggleModule(subscriptionID, false)

// Check if module is available
subscription, _ := moduleService.GetSubscription(subscriptionID)
if subscription.IsEnabled && subscription.Status == "active" {
    // Use gamification features
}
```

### 2. Monetization Per Module
```go
// Define pricing for each module
module := &models.Module{
    ID: "analytics",
    PricingModel: models.PricingModelPerUser,
    CostPerUser: 2.50,
}

// Calculate monthly cost automatically
cost := moduleService.CalculateModuleCost(module, userCount, 0, 0)
```

### 3. Multi-Company Structure
```go
// Create multiple companies under tenant
company1 := &models.Company{TenantID: "t1", Name: "Sales"}
company2 := &models.Company{TenantID: "t1", Name: "Support"}
companyService.CreateCompany(company1)
companyService.CreateCompany(company2)

// Each has independent configurations
```

### 4. Cross-Company Administration
```go
// Accounts team member can access both companies
// Via UserRole with cross-company permissions
roles, _ := userService.GetUserRoles(userID, tenantID)
// roles contains: ["accounts-manager"] with cross-company permissions
```

### 5. User Subset Management
```go
// User can be part of multiple company/project combinations
// Each with different roles

// User 123 is:
// - Manager in Company A
// - Member in Company B
// - Lead in Project A1
// - Viewer in Project B2

// Permissions aggregated across all roles
```

---

## Usage Example

### Scenario: Tech Startup with Multiple Teams

```go
// 1. Create tenant (during signup)
tenant := &models.Tenant{
    ID: "tech-startup-001",
    Name: "TechStartup Inc",
}

// 2. Create companies for different departments
salesCo := &models.Company{
    TenantID: "tech-startup-001",
    Name: "Sales Department",
    MaxUsers: 50,
    MaxProjects: 10,
}
companyService.CreateCompany(salesCo)

supportCo := &models.Company{
    TenantID: "tech-startup-001",
    Name: "Support Team",
    MaxUsers: 30,
    MaxProjects: 5,
}
companyService.CreateCompany(supportCo)

// 3. Create projects within each company
salesProject := &models.Project{
    CompanyID: salesCo.ID,
    Name: "Sales Pipeline",
    ProjectType: "sales",
    MaxUsers: 25,
}
companyService.CreateProject(salesProject)

// 4. Add team members to companies/projects
member := &models.CompanyMember{
    CompanyID: salesCo.ID,
    UserID: 1,
    Role: "manager",
}
companyService.AddMemberToCompany(member)

// 5. Subscribe to pricing plan
billingService.SubscribeToPlan(
    "tech-startup-001",
    "plan_professional",
    models.BillingCycleMonthly,
)

// 6. Enable/disable specific modules per scope
moduleService.SubscribeToModule(&models.ModuleSubscription{
    TenantID: "tech-startup-001",
    CompanyID: &salesCo.ID,
    ProjectID: &salesProject.ID,
    ModuleID: "gamification",
    MaxUsersAllowed: &25,
})

// 7. Track usage
usage := &models.UsageMetrics{
    TenantID: "tech-startup-001",
    CompanyID: &salesCo.ID,
    ActiveUsers: 18,
    NewUsers: 2,
}
billingService.RecordUsageMetrics(usage)

// 8. Calculate charges
charges, _ := billingService.CalculateMonthlyCharges("tech-startup-001")
```

---

## Admin Dashboard Requirements

The admin dashboard should provide:

### 1. Module Management
- [ ] Module registry (add, edit, delete)
- [ ] Enable/disable modules globally
- [ ] Set pricing models and costs
- [ ] Define dependencies
- [ ] Manage feature flags

### 2. Tenant Management
- [ ] List all tenants
- [ ] Manage company structure per tenant
- [ ] Manage project hierarchy
- [ ] Manage user assignments
- [ ] View license details

### 3. Billing Management
- [ ] Create and manage pricing plans
- [ ] View active subscriptions
- [ ] Monitor usage metrics
- [ ] Generate and send invoices
- [ ] Track payment status
- [ ] Apply discounts/refunds

### 4. Accounts Team Tools
- [ ] Cross-company user management
- [ ] Permission assignment across boundaries
- [ ] User role definitions
- [ ] Activity logging
- [ ] Audit trails

### 5. Reporting
- [ ] Module adoption rates
- [ ] Revenue by module/plan
- [ ] Usage trends
- [ ] Customer segmentation
- [ ] Churn analysis

---

## Migration Path

### Existing Tenants
1. Auto-migrate to "Default Company" structure
2. Map existing users to company members
3. Preserve all existing functionality
4. Enable gradual adoption of multi-company features

### Implementation Steps
1. Deploy schema changes
2. Deploy service layer
3. Deploy API handlers
4. Deploy admin dashboard
5. Deploy frontend components (Phase 3D)
6. Migrate existing data
7. Enable for new tenants

---

## Security Considerations

### Multi-Tenancy
- Strict tenant isolation in all queries
- Row-level security for company/project data
- Cross-tenant access prevention

### User Roles
- Role-based access control (RBAC)
- Permission inheritance rules
- Cross-scope permission handling
- Audit logging

### Billing
- Secure payment processing
- PCI compliance for card storage
- Encrypted sensitive data
- Audit trails for all financial transactions

---

## Future Enhancements

### Phase 3D (Frontend)
- [ ] Multi-company dashboard
- [ ] Project management UI
- [ ] Team member management interface
- [ ] Billing dashboard
- [ ] Module marketplace

### Phase 4 (Advanced)
- [ ] Custom pricing rules engine
- [ ] Usage-based auto-scaling
- [ ] Reservation pricing
- [ ] Seat-based licensing
- [ ] Volume discounts
- [ ] White-label support

---

## Glossary

| Term | Definition |
|------|-----------|
| **Module** | An individual feature that can be subscribed to independently |
| **Subscription** | A binding of a module to a tenant/company/project with status and limits |
| **Company** | An organizational unit under a tenant |
| **Project** | A work item under a company |
| **Pricing Model** | How a module's cost is calculated (per-user, per-project, etc.) |
| **Usage Metrics** | Daily snapshots of consumption for billing purposes |
| **Billing Cycle** | The frequency of invoicing (monthly, quarterly, annual) |
| **Accounts Team** | Users with cross-company/cross-project permissions |

