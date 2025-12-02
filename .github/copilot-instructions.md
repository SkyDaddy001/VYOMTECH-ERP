# GitHub Copilot Instructions for Contributors
## 🔍 Project Overview
This monorepo contains:
- **Next.js + TypeScript frontend** → `frontend/`
- **Go backend** → `cmd/main.go`, `internal/*`
- **Multi-tenant SaaS ERP + AI Call Center**
- Deployment + infra → `docker-compose.yml`, `k8s/`, `migrations/`

Your responsibility is to perform **narrow, safe edits** that preserve multi-tenancy, API contracts, and schema history.

---

## 🎯 Primary Objectives
- Make small targeted edits to handlers, services, UI components, or migrations.
- Never alter tenant routing or DB selection logic without explicit instructions + tests.
- Perform DB changes via **migrations + model updates**, never inline SQL.

---

## 📌 High-Impact File Map
| Area | Paths |
|------|-------|
| Frontend Pages | `frontend/app/**/page.tsx` |
| API Client | `frontend/services/api.ts` |
| Backend Entry | `cmd/main.go` |
| HTTP Handlers | `internal/handlers/**` |
| Business Logic | `internal/services/**` |
| DB Models | `internal/models/**` |
| Schema | `migrations/*.sql` |
| Middleware | `internal/middleware/**` |

---

## 🧩 Critical Architecture Rules

### Multi-tenancy
- Tenant derived from header: **`X-Tenant-ID`**
- Tenant routing and DB selection handled in `internal/middleware/`
- **Do not modify multi-tenant routing** unless an issue explicitly requests it

### DB & Models
- ORM: **GORM**
- Models include soft-deletes + audit fields
- Schema changes require:
  1. New migration file in `migrations/`
  2. Corresponding model updates in `internal/models/`

### Auth
- JWT + OAuth2
- Middleware attaches `user` and `tenant` context to handlers

### API Contract
- Backend REST → `/api/v1/*`
- If modifying a route:
  - Update handler + service
  - Update frontend call in `frontend/services/api.ts`
  - Update fetching logic (hooks/providers) if applicable

---

## 🖥 Frontend Conventions
- Next.js App Router structure
- Data: **React Query + Zustand + Context**
- Realtime: **socket.io**
- Place fetching logic in `frontend/services/api.ts` or custom hooks

---

## 🚀 Dev Commands
Frontend:
```sh
cd frontend && npm install && npm run dev
