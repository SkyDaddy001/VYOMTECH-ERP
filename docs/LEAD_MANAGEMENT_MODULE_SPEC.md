# VYOMTECH ERP + LMS Architecture Blueprint

## 1. Purpose

This document defines the architecture for a complete VYOMTECH ERP + LMS platform that combines business operations and learning operations inside a single multi-tenant system.

The architecture is built around the principle that VYOMTECH is not just an ERP or just an LMS. It is one integrated operating system for a business, its customers, employees, learning programs, and workflows.

## 2. Strategic direction

### 2.1 Product vision
VYOMTECH will support:
- multi-tenant organizations
- internal employees and external customers
- business operations across departments
- digital learning, training, assessments, and certifications
- sales pipeline and customer lifecycle management
- automated operational workflows
- analytics and reporting for leadership

### 2.2 Design principles
- single source of truth for users, accounts, entities, and roles
- tenant-aware data isolation
- API-first module architecture
- workflow-driven automation
- role-based access control across all domains
- auditability for all critical transactions
- scalable modules that can be enabled independently

## 3. System overview

### 3.1 Core layers

1. Client Layer
   - admin web app
   - employee portal
   - customer portal
   - LMS learner portal
   - mobile app (future)

2. API Layer
   - REST API
   - authentication and RBAC
   - service orchestration
   - event triggers and webhooks

3. Domain Layer
   - tenant core
   - CRM and sales
   - HRMS and payroll
   - finance and accounting
   - procurement and inventory
   - projects and assets
   - support and service desk
   - LMS and academy
   - analytics and AI

4. Data Layer
   - PostgreSQL primary database
   - Redis for caching and queues
   - object storage for files and media
   - event logs and audit trail storage

5. Workflow & Integration Layer
   - automation engine
   - webhook gateway
   - notifications
   - external connectors

## 4. Modules

### 4.1 Platform Core
Purpose: shared system foundation for all modules.

Core entities:
- tenants
- companies
- users
- roles
- permissions
- departments
- branches
- locations
- audit_logs
- settings
- notifications
- integrations

Core workflows:
- onboarding tenant
- user invitation and activation
- role assignment
- department and branch configuration
- global settings and branding
- audit tracking for all business actions

### 4.2 CRM & Sales
Purpose: manage leads, contacts, accounts, opportunities, and forecasting.

Core entities:
- leads
- lead_sources
- lead_statuses
- campaigns
- contacts
- accounts
- opportunities
- deals
- activities
- tasks
- notes
- products
- quotes
- contracts

Core workflows:
- lead capture from website or API
- lead assignment and qualification
- nurture automation via email/WhatsApp
- opportunity progression
- quote and contract creation
- deal closure and revenue recognition

### 4.3 Finance & Accounting
Purpose: manage accounting, invoices, payment flows, taxes, and budgeting.

Core entities:
- chart_of_accounts
- accounting_periods
- journal_entries
- journal_entry_lines
- invoices
- invoice_lines
- payments
- payment_terms
- vendors
- customers
- budgets
- tax_rules
- bank_accounts
- expense_claims

Core workflows:
- invoice creation and posting
- payment collection and reconciliation
- expense claims and approvals
- journal entry posting
- account aging and cashflow reporting
- tax and GST handling

### 4.4 HRMS & Payroll
Purpose: manage people, attendance, time, performance, and compensation.

Core entities:
- employees
- employee_documents
- departments
- job_positions
- employment_history
- leave_requests
- attendance
- payroll_runs
- payslips
- benefits
- performance_reviews
- goals
- training_assignments

Core workflows:
- employee onboarding
- attendance capture
- leave submission and approval
- payroll processing
- payslip distribution
- performance review cycle
- training assignment tied to role growth

### 4.5 Procurement & Inventory
Purpose: manage suppliers, purchase orders, stock, pricing, and warehouse operations.

Core entities:
- vendors
- purchase_orders
- purchase_order_lines
- warehouses
- inventory_items
- stock_movements
- stock_levels
- suppliers
- reorder_rules
- goods_receipts
- returns

Core workflows:
- purchase request to purchase order
- goods receipt and inspection
- stock updates and inventory valuation
- reorder alerts
- stock transfers
- supplier payment matching

### 4.6 Projects, Assets & Property Operations
Purpose: support project-based operations and asset management.

Core entities:
- projects
- project_tasks
- project_members
- milestones
- assets
- asset_assignments
- property_projects
- property_units
- bookings
- payment_plans
- service_requests

Core workflows:
- project kickoff
- task tracking and milestones
- asset lifecycle management
- property booking and installment tracking
- service request resolution

### 4.7 LMS / Academy
Purpose: deliver learning content, track learner progress, and manage certifications.

Core entities:
- courses
- course_modules
- lessons
- quizzes
- assignments
- learning_paths
- enrollments
- completion_records
- certificates
- assessments
- cohort_groups
- instructors
- learning_resources

Core workflows:
- course creation and publishing
- enrollment and cohort assignment
- learning path progression
- lesson completion tracking
- assessment scoring
- certification issuance
- reminder and re-enrollment flows

### 4.8 Support & Service Desk
Purpose: manage tickets, requests, and customer support workflows.

Core entities:
- tickets
- ticket_comments
- ticket_categories
- SLAs
- service_requests
- knowledge_articles
- escalation_rules

Core workflows:
- ticket intake
- triage and assignment
- SLA enforcement
- resolution and closure
- knowledge base updates

### 4.9 Analytics & Business Intelligence
Purpose: provide leadership and department-level analytics.

Core entities:
- dashboards
- reports
- report_templates
- metrics
- calculations
- alerts
- forecast_models

Core workflows:
- KPI aggregation
- executive dashboard updates
- anomaly detection
- pipeline forecasting
- learning completion analytics
- financial summary reporting

### 4.10 Workflow Automation & Integrations
Purpose: connect modules and automate business events.

Core entities:
- automation_rules
- triggers
- actions
- webhooks
- integrations
- connector_configs
- jobs
- scheduled_tasks

Core workflows:
- lead-to-opportunity automation
- payment failed notification
- invoice due alert
- course completion trigger
- approval chain routing
- webhook-based external sync

## 5. Exact tables by domain

Below is the database structure direction for the core platform.

### 5.1 shared
- tenants
- companies
- users
- user_profiles
- roles
- permissions
- role_permissions
- user_roles
- departments
- branches
- locations
- settings
- notifications
- audit_logs
- integrations
- integration_configs

### 5.2 crm
- leads
- lead_sources
- lead_statuses
- lead_notes
- campaigns
- campaign_members
- contacts
- accounts
- opportunities
- opportunity_stages
- tasks
- activities
- notes
- products
- quotes
- quote_lines
- contracts

### 5.3 finance
- chart_of_accounts
- accounting_periods
- journal_entries
- journal_entry_lines
- invoices
- invoice_lines
- payments
- payment_terms
- tax_rules
- vendors
- customers
- bank_accounts
- expense_claims
- budgets

### 5.4 hr
- employees
- employee_documents
- job_positions
- employment_history
- attendance
- leave_requests
- payroll_runs
- payslips
- benefits
- performance_reviews
- goals

### 5.5 procurement_inventory
- vendors
- purchase_orders
- purchase_order_lines
- goods_receipts
- warehouses
- inventory_items
- stock_levels
- stock_movements
- reorder_rules
- returns

### 5.6 projects_assets
- projects
- project_tasks
- project_members
- milestones
- assets
- asset_assignments
- property_projects
- property_units
- bookings
- payment_plans

### 5.7 lms
- courses
- course_modules
- lessons
- lesson_contents
- quizzes
- quiz_questions
- assignments
- enrollments
- completion_records
- certificates
- cohorts
- cohort_members
- instructors
- learning_paths

### 5.8 support
- tickets
- ticket_comments
- ticket_categories
- slas
- knowledge_articles

### 5.9 workflow
- automation_rules
- trigger_events
- rule_actions
- webhook_endpoints
- scheduled_jobs

## 6. Core workflows

### 6.1 Lead to cash workflow
1. website form creates lead
2. lead is scored and assigned
3. marketing automation sends initial response
4. sales rep qualifies lead
5. opportunity created
6. quote is generated
7. contract approved
8. invoice created
9. payment received
10. revenue recorded and reporting updates

### 6.2 Employee lifecycle workflow
1. employee record created
2. manager assigns role and department
3. documents and onboarding tasks created
4. attendance and leave configured
5. payroll run generated
6. review and performance cycle executed

### 6.3 Course learning workflow
1. learner enrolled in cohort or course
2. course content is unlocked
3. progress is tracked per lesson
4. quiz and assignment are evaluated
5. completion status is updated
6. certificate generated
7. compliance or skill tracking is reported

### 6.4 Procurement workflow
1. requisition is submitted
2. vendor selected
3. purchase order generated
4. goods received
5. invoice matched
6. payment issued
7. inventory updated

### 6.5 Support workflow
1. ticket created
2. category and priority assigned
3. SLA timer begins
4. owner assigned
5. resolution and closure
6. knowledge article updated if recurring issue

## 7. Technology recommendation

### 7.1 Backend
- Go or Node.js API services
- Prisma ORM for data modeling and migrations
- PostgreSQL for transactional data
- Redis for caching and queues
- BullMQ or similar for async jobs

### 7.2 Frontend
- Next.js for admin and customer portals
- Tailwind or design system for UI
- componentized architecture

### 7.3 Authentication & access
- JWT or session-based auth
- SSO support as enterprise feature
- RBAC with scoped tenant permissions

### 7.4 Workflow automation
- n8n or custom workflow engine
- centralized job runner
- webhook orchestration

### 7.5 Monitoring & operations
- error tracking
- app metrics
- audit trail logs
- performance tracing
- observability dashboards

## 8. Recommended repo structure

```text
VYOMTECH-ERP/
├── apps/
│   ├── admin/
│   ├── web/
│   ├── customer-portal/
│   └── lms/
├── api/
│   ├── routes/
│   ├── handlers/
│   ├── middleware/
│   └── services/
├── internal/
│   ├── auth/
│   ├── config/
│   ├── events/
│   └── workflow/
├── prisma/
│   ├── schema.prisma
│   ├── migrations/
│   └── seed.ts
├── modules/
│   ├── crm/
│   ├── finance/
│   ├── hr/
│   ├── inventory/
│   ├── lms/
│   ├── support/
│   ├── analytics/
│   └── automation/
├── docs/
├── docker-compose.yml
├── Makefile
├── go.mod
├── package.json
└── README.md
```

## 9. Final direction

The VYOMTECH platform should be treated as a modular enterprise operating system.

The ERP and LMS are not competing products; they are connected business layers within the same platform.

The next step is to implement the core foundation and the lead management module first, then expand into the rest of the business domains.
