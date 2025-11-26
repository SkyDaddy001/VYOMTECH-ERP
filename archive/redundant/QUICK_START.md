# 🚀 Quick Start Guide

## Current Status: ✅ Application Running

The **Multi-Tenant AI Call Center** is fully operational!

### See It In Action
```bash
# In a new terminal, from the project directory:
curl http://localhost:8080/health
# Response: {"status":"healthy"}
```

---

## Quick Commands

### Start/Stop Application
```bash
./startup.sh start   # Start MySQL container + application
./startup.sh stop    # Stop both containers
./startup.sh clean   # Remove containers completely
./startup.sh status  # Show running containers
```

### Test Authentication
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@test.com",
    "password": "Demo123!",
    "role": "user",
    "tenant_id": "default-tenant"
  }'

# Login  
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@test.com",
    "password": "Demo123!"
  }'
```

### Database Access
```bash
# Connect to MySQL
podman exec -it mysql-callcenter mysql -u callcenter_user -psecure_app_pass callcenter

# View logs
podman logs mysql-callcenter
```

### Rebuild Application
```bash
go build -o bin/main ./cmd/main.go
```

---

## Architecture Overview

```
┌─────────────────────────────────────────────┐
│  Multi-Tenant AI Call Center Application    │
│          (Go 1.24 on Port 8080)            │
├─────────────────────────────────────────────┤
│                                             │
│  ┌──────────────────────────────────────┐  │
│  │  HTTP Router (Gorilla Mux)           │  │
│  │  - /health                           │  │
│  │  - /api/v1/auth/*                    │  │
│  │  - /api/v1/agents/*                  │  │
│  └──────────────────────────────────────┘  │
│                   │                        │
│  ┌────────────────┴────────────────────┐  │
│  │  Service Layer                      │  │
│  │  - AuthService                      │  │
│  │  - AgentService                     │  │
│  │  - PasswordResetService             │  │
│  │  - EmailService                     │  │
│  │  - AIOrchestrator                   │  │
│  └────────────────┬────────────────────┘  │
│                   │                        │
│  ┌────────────────┴────────────────────┐  │
│  │  Middleware                         │  │
│  │  - JWT Authentication               │  │
│  │  - CORS Support                     │  │
│  │  - Error Recovery                   │  │
│  │  - Request Logging                  │  │
│  └────────────────┬────────────────────┘  │
│                   │                        │
└───────────────────┼────────────────────────┘
                    │
        ┌───────────┴───────────┐
        │                       │
    ┌───▼─────┐         ┌──────▼──────┐
    │  MySQL  │         │   Redis    │
    │ 3306    │         │  6379      │
    └─────────┘         └────────────┘
```

---

## Key Information

**Database Credentials**
- User: `callcenter_user`
- Password: `secure_app_pass`
- Database: `callcenter`
- Host: `127.0.0.1` (not localhost!)

**Default Tenant**
- ID: `default-tenant`
- Name: Default Tenant

**API Base URL**
- `http://localhost:8080`

**JWT Configuration**
- Algorithm: HS256
- Duration: 24 hours
- Secret: Configured via `JWT_SECRET` env var

---

## Files You Need to Know About

| File | Purpose |
|------|---------|
| `cmd/main.go` | Application entry point |
| `internal/config/config.go` | Configuration management |
| `internal/db/db.go` | Database connection pooling |
| `pkg/router/router.go` | HTTP route definitions |
| `internal/services/` | Business logic layer |
| `internal/handlers/` | HTTP request handlers |
| `migrations/001_initial_schema.sql` | Database schema |
| `startup.sh` | Automation script |

---

## Troubleshooting

### Application won't start
```bash
# Kill any stuck processes
pkill -9 main

# Rebuild
go build -o bin/main ./cmd/main.go

# Try again
./bin/main
```

### Port 8080 already in use
```bash
# Try different port
export SERVER_PORT=8081
./bin/main
```

### Database connection fails
```bash
# Make sure MySQL is running
podman ps | grep mysql

# Try connecting directly
podman exec mysql-callcenter mysql -u callcenter_user -psecure_app_pass -e "SELECT 1;"

# Verify host is 127.0.0.1
echo $DB_HOST
```

---

## What's Working

✅ User Registration & Login  
✅ JWT Authentication  
✅ Database Connection & Pooling  
✅ API Routing  
✅ CORS Support  
✅ Error Handling & Logging  
✅ Multi-tenant Architecture  
✅ Container Orchestration  

---

## What's Next

- [ ] Lead Management endpoints
- [ ] Call Management endpoints  
- [ ] Campaign Management endpoints
- [ ] AI Provider Adapters (OpenAI, Claude, Gemini)
- [ ] WebSocket Support
- [ ] Unit Tests
- [ ] Performance Testing
- [ ] Production Deployment

---

**Documentation**: See `APPLICATION_RUNNING.md` for full details
