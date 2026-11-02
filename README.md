# VYOMTECH ERP

A multi-tenant enterprise platform combining ERP operations, CRM, HR, finance, real-estate operations, and learning workflows in a single application foundation.

## Repository purpose

This repository is organized around a clean delivery model:

- `apps/backend` contains the Go API and service layer.
- `apps/frontend` contains the Next.js user-facing application.
- `docs/` contains architecture, roadmap, runbooks, and contributor guidance.
- `scripts/` contains automation for setup, validation, and operational tasks.
- `prisma/` contains database schema and Prisma tooling.

## Current platform scope

The platform currently includes the following operational domains:

- tenant and user management
- role-based access control
- sales and lead lifecycle management
- finance and accounting workflows
- HR and payroll capabilities
- procurement and inventory support
- project and property operations
- LMS and learning workflows
- analytics and AI orchestration

## Standardized repo layout

```text
VYOMTECH-ERP/
├── apps/
│   ├── backend/
│   └── frontend/
├── docs/
│   ├── ARCHITECTURE.md
│   ├── DEVELOPMENT.md
│   ├── PROJECT_STRUCTURE.md
│   ├── ROADMAP.md
│   └── ...
├── scripts/
│   ├── bootstrap.sh
│   └── validate.sh
├── prisma/
├── .github/
├── cmd/
├── internal/
├── pkg/
├── frontend/
├── README.md
├── package.json
├── go.mod
├── docker-compose.yml
├── .env.example
└── ...
```

## Quick start

### Backend

```bash
npm run dev:backend
```

### Frontend

```bash
npm run dev:frontend
```

### Full stack

```bash
npm run dev
```

## Commands

```bash
npm run setup
npm run validate
npm run build:backend
npm run build:frontend
npm run health-check
```

## Operating model

- Treat the current branch as a baseline history marker.
- Use the working branch for structural cleanup and future refactors.
- Keep architecture documentation in sync with service boundaries.
- Prefer domain-first grouping for new modules.
- Avoid duplicate business logic across frontend and backend packages.

## Development notes

- Go backend is the system-of-record for business logic and APIs.
- Prisma is used for schema and database tooling when the project needs formal data modeling.
- Frontend uses Next.js and is expected to consume the backend through shared API contracts.
- Documentation should remain the single source of truth for architecture and onboarding.

## Documentation

See:

- `docs/ARCHITECTURE.md`
- `docs/PROJECT_STRUCTURE.md`
- `docs/DEVELOPMENT.md`
- `docs/ROADMAP.md`

## License

This project is managed under the repository's current project license and operational policy.
