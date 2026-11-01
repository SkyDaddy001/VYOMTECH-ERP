# VYOMTECH PR Execution Sequence

This document defines the ordered implementation path for VYOMTECH-ERP as a merged, production-ready platform.

## PR order

### PR-01 — Foundation & Platform Baseline
Status: Ready for implementation

Scope:
- repo standardization
- environment and configuration baseline
- shared tenant model
- application shell and admin layout
- observability and logging baseline
- deployment pipeline

Acceptance criteria:
- app boots in dev mode
- tenant model exists and is documented
- config is centralized
- all modules share a common environment contract
- audit log and observability basic setup is complete

Audit gate:
- verify no broken references
- verify development environment reproducibility
- verify deployment config is valid

### PR-02 — Tenant Core & RBAC
Status: Ready for implementation

Scope:
- users, roles, permissions
- tenant and company model
- department, location, branch configuration
- access guards and authorization middleware
- audit trail for access changes

Acceptance criteria:
- tenant isolation is enforced in all business tables
- RBAC works at module and action levels
- admin can create tenants and assign roles
- audit log records critical access events

Audit gate:
- verify permission checks for all sample roles
- verify no role grants bypass tenant boundaries

### PR-03 — Lead Management & CRM
Status: Ready for implementation

Scope:
- leads and lead sources
- campaigns and attribution
- contacts and accounts
- opportunities and sales pipeline
- tasks, notes, and activities
- automation triggers and WhatsApp/email integration hooks

Acceptance criteria:
- lead capture API works from forms and webhooks
- lead scoring and assignment rules run
- opportunity conversion works end-to-end
- campaign metrics are reportable
- sales workflow is visible in dashboard

Audit gate:
- verify deduplication logic
- verify lead stage transitions
- verify automation triggers reflect real business flow

### PR-04 — Finance, HR & Procurement
Status: Ready for implementation

Scope:
- finance core and chart of accounts
- vendors, customers, invoices, payments
- employee management and payroll
- procurement and purchase orders
- expense approvals and accounting entries

Acceptance criteria:
- invoice and payment lifecycle is operational
- payroll generation produces payslips and reports
- purchase order and inventory intake flows are connected
- financial, payroll, and vendor records are tenant-scoped

Audit gate:
- verify transaction consistency
- verify all financial records link to the correct tenant and accounting period

### PR-05 — Inventory & Projects
Status: Ready for implementation

Scope:
- warehouses and stock levels
- stock movement logs
- reorder rules and inventory valuation
- projects and milestones
- assets and property-related workflows

Acceptance criteria:
- inventory movement history is accurate
- stock level updates trigger reorder rules
- project and milestone tracking works
- asset lifecycle tracking is operational

Audit gate:
- verify stock movement integrity
- verify project financial and operational linkage

### PR-06 — LMS & Academy
Status: Ready for implementation

Scope:
- courses, modules, lessons
- cohorts and enrollments
- progress tracking and assessments
- certificates and completion logs
- LMS reporting and learner analytics

Acceptance criteria:
- learners can enroll and progress through content
- assessments can be scored
- certificates are generated correctly
- completion data is queryable and auditable

Audit gate:
- verify progression logic
- verify certificate issuance rules and learner data access

### PR-07 — Support, Analytics & AI
Status: Ready for implementation

Scope:
- tickets and service desk
- knowledge base and SLA tracking
- dashboards and report generation
- AI assistant and recommendation layer
- workflow orchestration and integrations

Acceptance criteria:
- support lifecycle is handled end-to-end
- executive dashboards show key business metrics
- AI recommendations are based on module data
- workflows integrate across modules without breaking boundaries

Audit gate:
- verify dashboard metric correctness
- verify AI recommendations are based on valid tenant data
- verify workflow triggers do not create loops

## Merge order rules

1. Foundation merges first.
2. Tenant/RBAC merges second.
3. CRM and sales merge third.
4. Finance and HR merge next.
5. Operations and inventory merge after that.
6. LMS merges before intelligence and support expansion.
7. Support, analytics, and AI merge last.

This order ensures no module depends on a tenant or data model that is not already established.

## Merge discipline

Each PR requires:
- implementation completion
- audit checklist sealed
- dependency validation
- change review against shared data model
- merge only after the previous PR is stable

## Final release condition

Once all PRs are merged and audited, VYOMTECH-ERP becomes a complete modular enterprise platform for ERP + LMS + CRM + operations.
