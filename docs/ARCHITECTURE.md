# System architecture

## 1. Platform intent

VYOMTECH ERP is designed as a single multi-tenant enterprise platform that combines ERP, CRM, operations, and learning workflows under one backend system-of-record.

## 2. Core architectural layer model

### Application layer

- backend API services for business logic
- tenant-aware authorization and policy enforcement
- workflow orchestration for modules and cross-domain automation

### Data layer

- domain models and persistence stores
- Prisma-based schema tooling and database migrations
- shared access patterns with tenant scoping

### Presentation layer

- Next.js frontend for business users and administrators
- dashboard-driven workflows for operational visibility and actions

## 3. Domain grouping

The system is organized by operational domain rather than by team-only folders. Domain grouping helps maintain boundaries around:

- tenant and admin
- sales and CRM
- finance and invoicing
- HR and payroll
- procurement and inventory
- projects and property operations
- LMS and learning paths
- analytics and AI

## 4. Infrastructure and deployment model

The project is structured to support:

- local development via Docker and Go runtime
- frontend development in Next.js
- database and schema generation via Prisma
- production deployment with environment-specific configuration

## 5. Design principles

- single source of truth for business logic in backend services
- tenant isolation by default in all domain logic
- clear separation between platform core and application modules
- documentation-first onboarding and architecture management
- modular growth without duplicate entity models
