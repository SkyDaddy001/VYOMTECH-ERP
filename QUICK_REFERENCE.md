# Quick Reference - Commands & Workflows

## Build & Run

### Local Development
```bash
# Install dependencies
go mod download
go mod tidy

# Build
go build -o bin/main ./cmd/main.go

# Run
./bin/main

# Or using make
make build
make run
```

### With Docker Compose
```bash
# Start entire stack (MySQL, Redis, Prometheus, Grafana, App)
make dev-up

# Stop stack
make dev-down

# View logs
docker-compose logs -f app
```

### Docker Build
```bash
# Build image
docker build -t multi-tenant-ai-callcenter:latest .

# Run container
docker run -p 8080:8080 --env-file .env multi-tenant-ai-callcenter:latest
```

## Testing API

### Health Checks
```bash
# Application health
curl http://localhost:8080/health

# Readiness check
curl http://localhost:8080/ready
```

### Authentication

#### Register New User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "role": "agent",
    "tenant_id": "tenant-123"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

Response includes JWT token:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "role": "agent",
    "tenant_id": "tenant-123"
  },
  "message": "Login successful"
}
```

### Using JWT Token
```bash
# Store token
TOKEN="eyJhbGciOiJIUzI1NiIs..."

# Use in requests
curl -X GET http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer $TOKEN"
```

### Agent Management

#### Get All Agents for Tenant
```bash
curl -X GET http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer $TOKEN"
```

#### Get Specific Agent
```bash
curl -X GET http://localhost:8080/api/v1/agents/1 \
  -H "Authorization: Bearer $TOKEN"
```

#### Update Agent Availability
```bash
curl -X PATCH http://localhost:8080/api/v1/agents/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "availability": "online"
  }'
```

Valid values: `online`, `offline`, `busy`

#### Get Available Agents
```bash
curl -X GET http://localhost:8080/api/v1/agents/available \
  -H "Authorization: Bearer $TOKEN"
```

#### Get Agent Statistics
```bash
curl -X GET http://localhost:8080/api/v1/agents/stats \
  -H "Authorization: Bearer $TOKEN"
```

### Password Reset

#### Request Password Reset
```bash
curl -X POST http://localhost:8080/api/v1/password-reset/request \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'
```

#### Reset Password with Token
```bash
curl -X POST http://localhost:8080/api/v1/password-reset/reset \
  -H "Content-Type: application/json" \
  -d '{
    "token": "token-from-email",
    "new_password": "NewSecurePass123!"
  }'
```

## Development Commands

### Code Quality
```bash
# Format code
make format
# or
go fmt ./...

# Run linter
make lint
# or
golangci-lint run ./...

# Run tests
make test

# Run tests with coverage
make test-all
# View HTML coverage report
open coverage.html
```

### Database

#### Create Local MySQL Container
```bash
make db-docker
```

#### Stop MySQL Container
```bash
make db-stop
```

#### Run Migrations
```bash
mysql -u callcenter_user -p secure_app_pass -h localhost callcenter < migrations/001_initial_schema.sql
```

## Configuration

### Environment Variables (.env)
```bash
# Copy example to .env
cp .env.example .env

# Edit with your settings
nano .env
```

### Key Settings
- `SERVER_PORT=8080` - Application port
- `DB_HOST=localhost` - Database host
- `DB_USER=callcenter_user` - DB username
- `DB_PASSWORD=secure_app_pass` - DB password
- `JWT_SECRET=<your-secret-key>` - JWT signing key (min 32 chars)
- `SMTP_HOST=smtp.gmail.com` - Email server
- `DEBUG=true` - Debug mode

## Debugging

### Enable Debug Logging
```bash
DEBUG=true ./bin/main
```

### Database Connection Issues
```bash
# Test MySQL connection
mysql -h localhost -u callcenter_user -p -e "SELECT 1"

# Check if port 3306 is open
netstat -an | grep 3306
```

### Port Conflicts
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

## Project Structure

```
├── cmd/
│   └── main.go                    # Entry point
├── internal/
│   ├── config/                    # Configuration
│   ├── db/                        # Database setup
│   ├── handlers/                  # HTTP handlers
│   ├── middleware/                # HTTP middleware
│   ├── models/                    # Data models
│   └── services/                  # Business logic
├── pkg/
│   ├── auth/                      # JWT operations
│   ├── logger/                    # Logging
│   └── router/                    # Route setup
├── migrations/                    # Database migrations
├── monitoring/                    # Prometheus config
├── k8s/                          # Kubernetes manifests
├── go.mod & go.sum               # Modules
├── Dockerfile                     # Docker image
├── docker-compose.yml             # Local dev stack
├── Makefile                       # Development commands
└── .env.example                  # Config template
```

## File Locations Reference

| File | Purpose | Location |
|------|---------|----------|
| Main Application | Entry point | `cmd/main.go` |
| Config Template | Environment settings | `.env.example` |
| Database Schema | Table definitions | `migrations/001_initial_schema.sql` |
| Logger | Logging utility | `pkg/logger/logger.go` |
| Router | API routes | `pkg/router/router.go` |
| Auth Service | Authentication | `internal/services/auth.go` |
| Agent Service | Agent management | `internal/services/agent.go` |
| Email Service | Email sending | `internal/services/email.go` |
| Auth Handler | Auth endpoints | `internal/handlers/auth.go` |
| Middleware | HTTP middleware | `internal/middleware/auth.go` |
| Docker Compose | Local stack | `docker-compose.yml` |
| Kubernetes | K8s manifests | `k8s/deployment.yaml` |

## Kubernetes Deployment

### Prerequisites
```bash
# Install kubectl
# Install helm (optional)

# Check cluster
kubectl cluster-info
```

### Deploy to Kubernetes
```bash
# Apply configurations
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Check deployment status
kubectl get deployments
kubectl get pods
kubectl get services

# View logs
kubectl logs -f deployment/callcenter-app

# Port forward for local testing
kubectl port-forward svc/callcenter-service 8080:8080
```

### Scaling
```bash
# Scale replicas
kubectl scale deployment callcenter-app --replicas=3

# Auto-scaling (requires metrics server)
kubectl autoscale deployment callcenter-app --min=2 --max=10 --cpu-percent=80
```

## Troubleshooting

### "Connection refused" on port 8080
```bash
# Ensure database is running
make dev-up

# Check if process is running
ps aux | grep main

# Try different port
SERVER_PORT=9000 ./bin/main
```

### "Database connection failed"
```bash
# Check MySQL is running
docker ps | grep mysql

# Test credentials
mysql -h localhost -u callcenter_user -p

# Recreate database
make db-stop
make db-docker
```

### "Invalid JWT Token"
```bash
# Ensure JWT_SECRET is set in .env
# Regenerate token with login endpoint
# Check token hasn't expired (24 hours by default)
```

### "Unauthorized" on protected endpoints
```bash
# Verify Authorization header format
# Should be: "Authorization: Bearer <token>"

# Login to get new token
curl -X POST http://localhost:8080/api/v1/auth/login ...
```

## Performance Monitoring

### With Prometheus
```bash
# Access Prometheus dashboard
http://localhost:9090

# Query metrics
curl http://localhost:9090/api/v1/targets
```

### With Grafana
```bash
# Access Grafana dashboard
http://localhost:3000

# Default credentials
# Username: admin
# Password: admin

# Add Prometheus datasource: http://prometheus:9090
```

## Common Tasks

### Creating a New Database Migration
```bash
# Create migration file
touch migrations/002_add_new_table.sql

# Add SQL commands
# Apply it
mysql -u callcenter_user -p callcenter < migrations/002_add_new_table.sql
```

### Adding New Endpoint
1. Create handler in `internal/handlers/`
2. Add service method in `internal/services/`
3. Register route in `pkg/router/router.go`
4. Add middleware if needed
5. Test with curl

### Running Tests with Coverage
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

## Resources

- **Go Documentation**: https://golang.org/doc/
- **Gorilla Mux**: https://github.com/gorilla/mux
- **JWT Go**: https://github.com/golang-jwt/jwt
- **MySQL Go Driver**: https://github.com/go-sql-driver/mysql
- **Docker Docs**: https://docs.docker.com/
- **Kubernetes Docs**: https://kubernetes.io/docs/

---

**Last Updated:** November 21, 2025  
**Status:** Ready for use  
**Application:** Multi-Tenant AI Call Center
