Plan: Phase One Deployment and Verification

Goal
- Run Phase 1 DB migration, deploy backend and frontend, and verify the 4 API endpoints.

Assumptions
- You have `podman`/`podman-compose` available and containers are defined in `docker-compose.yml`.
- Services run in the repository root.
- DB credentials come from `docker-compose.yml` (MYSQL_ROOT_PASSWORD=rootpass, DB user `callcenter_user`).
- Endpoints require `X-Tenant-ID` header and valid Authorization bearer token.

High-level steps
1. Start containers (db, redis, app, frontend) with podman-compose
2. Wait for MySQL to be healthy
3. Run the Phase 1 migration SQL file
4. Verify tables exist and basic schema sanity checks
5. Build & restart backend service (inside container or by rebuilding image)
6. Build & restart frontend (inside container or by rebuilding image)
7. Test the 4 API endpoints with curl (include tenant + auth headers)
8. Collect logs and verify no errors; if errors, gather DB/APP logs
9. Rollback plan if migration or deploy fails

Commands (copy/paste)

# 1) Start containers (from repo root)
```bash
cd "c:\Users\Skydaddy\Desktop\Developement"
podman-compose up -d
```

# 2) Wait for MySQL healthy state (poll / check)
```bash
# wait 10s then show health/state
sleep 10
podman ps --format "table {{.ID}}\t{{.Image}}\t{{.Names}}\t{{.Status}}"
```

# 3) Run migration (inside MySQL container)
```bash
# run migration file from host into container
podman exec -i callcenter-mysql bash -c "mysql -u root -prootpass callcenter" < migrations/004_phase1_features.sql
```

# 4) Verify tables exist
```bash
podman exec callcenter-mysql mysql -u root -prootpass callcenter -e "SHOW TABLES;"
# or check specific tables
podman exec callcenter-mysql mysql -u root -prootpass callcenter -e "SHOW TABLES LIKE 'lead_%' OR 'agent_%' OR 'audit_%';"
```

# 5) (Optional) Build/restart backend container
```bash
# if code changed and you want to rebuild the container image
podman-compose build app
podman-compose up -d app
# or restart container if image already baked
podman restart callcenter-app
```

# 6) (Optional) Build/restart frontend container
```bash
podman-compose build frontend
podman-compose up -d frontend
# or restart
podman restart callcenter-frontend
```

# 7) Test API endpoints (replace <token> and tenant)
```bash
# 1) Get single lead score
curl -i -H "X-Tenant-ID: test-tenant" -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/leads/1/score

# 2) Calculate single lead score
curl -i -X POST -H "X-Tenant-ID: test-tenant" -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/leads/1/score/calculate

# 3) Get leads by category
curl -i -H "X-Tenant-ID: test-tenant" -H "Authorization: Bearer <token>" \
  "http://localhost:8080/api/v1/leads/scores/category/hot?limit=50"

# 4) Batch calculate scores
curl -i -X POST -H "X-Tenant-ID: test-tenant" -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/leads/scores/batch-calculate
```

Validation checklist
- Migration completed without SQL errors
- Tables created and FK/indexes present
- Backend container logs show successful DB connection and service started on port 8080
- Frontend container serves app on port 3000 (or as configured)
- Each API endpoint returns expected HTTP status codes (200/202/404 etc.)
- Batch job returns accepted/processing status

Troubleshooting hints
- If mysql client not present in host, run migrations inside container (we do that above). If that fails, copy SQL into container and run from /tmp.
- To view container logs:
```bash
podman logs -f callcenter-mysql
podman logs -f callcenter-app
podman logs -f callcenter-frontend
```
- If migration fails due to existing schema, capture the failing statement from logs and run it manually with `-v` verbosity.

Rollback plan
- Stop app/frontend containers immediately: `podman stop callcenter-app callcenter-frontend`
- Restore DB from backup (if available). If no backup, try to reverse changes by analyzing migration file and issuing DROP TABLE statements for newly created tables (risky) — recommended to restore from snapshot.

Monitoring & next steps
- Monitor `callcenter-app` logs and `audit_logs` for migration effects
- Run E2E tests for lead flows (create lead, assign, calculate score)
- Add smoke-check endpoint script to run periodically after deploy

Notes
- The DB credentials used by containers are set in `docker-compose.yml`. If your production environment differs, adapt the commands accordingly.
- If you want me to run these commands now, say "Run migration and deploy" and I will proceed to execute them sequentially and report results.
