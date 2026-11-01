# VYOMTECH Production Roadmap

## 1. Objective

The goal is to finish VYOMTECH-ERP as a production-ready enterprise platform that includes the ERP core, the lead management engine, and the LMS layer under one unified multi-tenant operating system.

This roadmap is based on the actual repo state and the strategic decision to consolidate all business logic into VYOMTECH-ERP.

## 2. Delivery strategy

### Phase 0: stabilization and foundation
Goal: make the repo coherent, production-safe, and architecture-aligned.

Tasks:
- define canonical repo architecture for monorepo or modular service layout
- finalize tenant, user, auth, RBAC model
- clean and standardize documentation
- agree on stack and service boundaries
- define environment variables and deployment model

Deliverables:
- single architecture decision record
- updated repo README and modules docs
- production environment template
- security baseline and secret management

### Phase 1: ERP core MVP
Goal: deploy the minimum viable enterprise operations foundation.

Modules:
- tenant management
- company and branch setup
- users, roles, permissions
- departments and locations
- audit log and notifications
- settings and branding
- dashboard and executive overview

Deliverables:
- auth and multi-tenant foundation
- admin dashboard shell
- RBAC enforcement
- system config and tenant onboarding

### Phase 2: lead management and CRM
Goal: operationalize the proven lead engine based on propertyscouts.

Modules:
- leads
- campaigns
- sources and attribution
- lead status flow
- assignment and tasks
- notes and activities
- opportunity creation
- sales pipeline

Deliverables:
- lead capture APIs
- assignment logic and workflows
- email/WhatsApp automation hooks
- pipeline dashboard
- reporting on conversion and sources

### Phase 3: finance and accounting
Goal: close the operational and revenue loop.

Modules:
- chart of accounts
- vendors and customers
- invoices and payments
- journal entries
- budgets and expenses
- reports and aging

Deliverables:
- invoice lifecycle
- payment reconciliation
- GL entries
- financial dashboard

### Phase 4: HR and payroll
Goal: make the platform useful to people operations.

Modules:
- employees
- attendance and leave
- payroll runs
- payslips
- performance review
- onboarding and documents

Deliverables:
- employee directory
- leave management
- payroll engine
- HR dashboard

### Phase 5: procurement and inventory
Goal: add supply chain operations.

Modules:
- vendors
- purchase orders
- warehouses
- stock items and levels
- stock movements
- reorder rules

Deliverables:
- purchase workflow
- inventory system
- warehouse reporting
- cost tracking

### Phase 6: projects, assets, and property workflows
Goal: support real-world operational use cases.

Modules:
- projects
- milestones
- tasks
- assets
- property bookings and payments
- service requests

Deliverables:
- project tracking
- asset lifecycle management
- property-related operations workflows

### Phase 7: LMS and academy layer
Goal: build the learning system into the same platform.

Modules:
- courses
- modules and lessons
- assessment engine
- enrollments
- cohort management
- certificates
- learning paths

Deliverables:
- course catalog and authoring
- learner portal
- progress tracking
- certificate generation
- reporting by learner and cohort

### Phase 8: support and service desk
Goal: add customer support and internal help desk operations.

Modules:
- tickets
- categories and priorities
- SLA tracking
- knowledge base
- escalation workflows

Deliverables:
- support dashboard
- service resolution flow
- ticket analytics

### Phase 9: analytics and AI
Goal: produce executive intelligence and automation.

Modules:
- dashboards
- KPI reporting
- forecasting
- AI recommendations
- workflow triggers
- anomaly detection

Deliverables:
- executive dashboard
- pipeline and revenue forecasting
- operational insights
- AI assistant layer

## 3. Milestones by quarter

### Quarter 1: foundation
- repo architecture freeze
- tenant/auth foundation
- admin dashboard shell
- deployment pipeline

### Quarter 2: growth and operations
- CRM and lead management
- sales pipeline and automation
- finance baseline

### Quarter 3: people and execution
- HR/payroll
- procurement/inventory
- project and asset tracking

### Quarter 4: learning and intelligence
- LMS module
- support and knowledge base
- analytics and AI workflow layer

## 4. Development order of modules

Recommended order:
1. platform core
2. CRM/leads
3. finance
4. HR/payroll
5. inventory/procurement
6. projects/assets
7. LMS
8. support
9. analytics
10. AI automation layer

This order minimizes rework and ensures the platform reaches revenue-generating use cases early.

## 5. Release strategy

### MVP release
- tenant and admin foundation
- lead management
- basic CRM pipeline
- dashboard and reports
- basic invoicing

### Beta release
- HR and payroll
- procurement and inventory
- LMS pilot cohort
- workflow automation

### Production release
- full multi-module ERP + LMS platform
- enterprise security hardening
- analytics and alerting
- SSO and integrations
- documentation and support process

## 6. Quality gates

Every module must pass:
- schema validation
- RBAC test coverage
- payment or transaction audit
- API contract review
- tenant isolation checks
- audit log verification
- data retention policy validation

## 7. Risks and mitigations

### Risk: too many modules too early
Mitigation:
- implement by domain in sequence
- keep shared platform stable
- avoid building duplicate user models

### Risk: inconsistent tenancy model
Mitigation:
- enforce tenant_id on all business tables
- centralize auth and authorization

### Risk: disconnected workflow logic
Mitigation:
- use one workflow engine for cross-module automation

### Risk: poor analytics quality
Mitigation:
- define metrics and source-of-truth tables early

## 8. Success criteria

The platform is ready when:
- admin can create a tenant and configure a company
- sales teams can manage leads and opportunities
- finance team can issue invoices and track payments
- employees can be onboarded and paid
- learners can enroll and complete courses
- support can track and fulfill service requests
- executive dashboards show real business KPI information

## 9. Final recommendation

The correct path is to finish VYOMTECH-ERP as the master platform, then extend it as a modular enterprise system with integrated LMS and CRM capabilities. This approach is more resilient than maintaining separate products that duplicate customer, user, tenant, and workflow logic.

VYOMTECH should become the system of record for business operations and learning.
