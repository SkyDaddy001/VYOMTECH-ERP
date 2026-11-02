# Development guide

## Prerequisites

- Node.js 18+
- npm 9+
- Go 1.24+
- Docker and Docker Compose for local stack orchestration

## Install dependencies

```bash
npm install
npm run prisma:generate
```

## Run backend

```bash
npm run dev:backend
```

## Run frontend

```bash
npm run dev:frontend
```

## Run the full local stack

```bash
npm run dev
```

## Validation

```bash
npm run validate
```

## Build checks

```bash
npm run build:backend
npm run build:frontend
```

## Operational tips

- Keep environment-specific secrets in `.env` files and never commit them.
- Run validation before merging structural refactors.
- Update docs whenever the service boundaries or module structure change.
