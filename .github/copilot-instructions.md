## Quick Context

This monorepo hosts a Next.js TypeScript frontend (`frontend/`) and a Go backend (`cmd/main.go`, `internal/`). It's a multi-tenant SaaS ERP + AI call-center platform. Key infra and deploy manifests live at `docker-compose.yml`, `k8s/`, and `migrations/`.

## Primary Goals For AI Coding Agents
- Make small, focused edits to handlers, services, migrations, or UI components.
- Preserve tenant isolation: never change tenant routing or DB selection without tests and reviewer sign-off.
- Prefer adding migrations and service-layer changes over ad-hoc DB edits.

## Where to Look (high-value files)
- Frontend app: `frontend/app/` (Next.js App Router pages). Example: `frontend/app/dashboard/page.tsx`.
- Frontend API client: `frontend/services/api.ts` (Axios wrapper used across UI).
- Backend entry: `cmd/main.go` (server bootstrap).
- Backend handlers: `internal/handlers/` (HTTP handlers wired into Gorilla Mux).
- Business logic/services: `internal/services/` (domain logic called from handlers).
- Models and DB: `internal/models/` (GORM models) and `migrations/*.sql` (schema files).
- Config and middleware: `internal/config/`, `internal/middleware/` (auth, tenant routing, CORS).
- Dev/infra: `docker-compose.yml`, `k8s/*.yaml`, `Makefile`.

## Important Patterns & Conventions
- Multi-tenancy: tenant selection is header-driven. Look for `X-Tenant-ID` usage and `Tenant Router` logic in `internal/middleware/` before changing data access code.
- DB access: uses GORM; models live in `internal/models/` and expect soft-deletes/audit fields. Add migrations to `migrations/` when altering schema.
- Auth: JWT + OAuth2. Handlers expect `Authorization: Bearer <token>` and middleware attaches `user`/`tenant` context.
- API shape: backend exposes REST endpoints under `/api/v1/*` (see top-level README for common routes). Frontend calls go through `frontend/services/api.ts`—update both sides when changing contract.
- Frontend routing: Next.js App Router. Pages/components follow `frontend/app/<route>/page.tsx` + server/client component boundaries.
- State & data fetching: uses React Context + Zustand + React Query; prefer updating `services/api.ts` + hooks in `frontend/hooks/` or `frontend/providers/`.
- Real-time: socket.io used for live dashboards; check `frontend` + server socket setup before touching real-time features.

## Developer Workflows (commands you can run)
- Frontend dev:
  - `cd frontend && npm install`
  - `cd frontend && npm run dev` (localhost:3000)
- Backend dev:
  - `go mod download`
  - Load migrations: `mysql -u root -p database_name < migrations/*.sql`
  - `go run cmd/main.go` (API: http://localhost:8080)
- Tests:
  - Frontend: `cd frontend && npm test` (or `npm run test:watch`)
  - Backend: `go test ./...`
- Docker/k8s:
  - `docker-compose up -d`
  - `kubectl apply -f k8s/deployment.yaml`

## How to make safe changes (checklist)
1. Search for `X-Tenant-ID` and `tenant` uses before editing data access.
2. If changing DB schema: add a new SQL migration file in `migrations/` and update any GORM model in `internal/models/`.
3. Update backend API contract in `internal/handlers/` and corresponding client calls in `frontend/services/api.ts`.
4. Run unit tests: `go test ./...` and `cd frontend && npm test`.
5. Run local dev: backend + frontend and exercise the UI flows affected.
6. Create a short PR with rationale and link to migration files; request a backend reviewer for DB/tenant changes.

## Examples (copyable snippets)
- Start frontend dev server:
  - `cd frontend && npm install && npm run dev`
- Run backend locally with migrations:
  - `go mod download && mysql -u root -p database_name < migrations/001_initial_schema.sql && go run cmd/main.go`
- Find tenant routing code:
  - `rg "X-Tenant-ID|Tenant" internal | sed -n '1,120p'`

## What not to do
- Do not modify tenant routing, global DB connection selection, or shared migration history without coordination.
- Avoid wide refactors across frontend and backend in a single PR — prefer small iterative changes.

## Useful searches for quick orientation
- `rg "X-Tenant-ID|multi-tenant|Tenant"`
- `rg "internal/handlers|internal/services|internal/models"`
- `rg "NEXT_PUBLIC_API_URL|services/api.ts" frontend -n`

## Follow-up
If anything in this file is unclear or you want more examples (e.g., common handler-to-UI change steps), reply and I'll iterate.
