# Project structure

## Root-level organization

```text
VYOMTECH-ERP/
├── apps/                  # platform applications
│   ├── backend/          # Go service layer and API entry points
│   └── frontend/         # Next.js application
├── docs/                 # architecture and project guidance
├── scripts/              # automation and validation helpers
├── prisma/               # Prisma schema and tooling
├── cmd/                  # backend service entrypoints
├── internal/             # application internals
├── pkg/                 # reusable libraries and shared packages
├── frontend/             # frontend app (legacy compatibility path)
├── docker-compose.yml    # local stack orchestration
├── go.mod                # Go module definition
├── package.json          # repo-level scripts
├── .env.example          # environment template
└── README.md             # repo overview
```

## Naming conventions

- use lowercase, hyphenated directory names for operational groupings
- use PascalCase for Go types and service objects where appropriate
- keep domain names aligned with business modules rather than random technical labels
- avoid mixing transient scripts, docs, and runtime artifacts in core application folders

## Module ownership

- `cmd/` hosts executable entrypoints
- `internal/` contains backend implementation details and app logic
- `pkg/` contains reusable utilities and cross-cutting components
- `prisma/` owns database modeling and migration tooling
- `apps/frontend/` owns UI implementation and client-side state
- `docs/` is the canonical place for architecture and process guidance
