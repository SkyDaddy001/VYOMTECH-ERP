# Development Environment Setup

## Prerequisites
- Go 1.24+
- MySQL 8.0+
- Git

## Quick Start

### 1. Install Dependencies
```bash
go mod download
go mod tidy
```

### 2. Setup Database
```bash
# Create database
mysql -u root -p < migrations/001_initial_schema.sql

# Or use Docker
docker run --name mysql-callcenter \
  -e MYSQL_ROOT_PASSWORD=rootpass \
  -e MYSQL_DATABASE=callcenter \
  -e MYSQL_USER=callcenter_user \
  -e MYSQL_PASSWORD=secure_app_pass \
  -p 3306:3306 \
  mysql:8.0
```

### 3. Configure Environment
```bash
cp .env.example .env
# Edit .env with your local settings
```

### 4. Run Application
```bash
go run ./cmd/main.go
```

### 5. Test Endpoints
```bash
# Health check
curl http://localhost:8080/health

# Readiness check
curl http://localhost:8080/ready

# Password reset request
curl -X POST http://localhost:8080/api/v1/password-reset/request \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'
```

## Development Commands

### Build
```bash
go build -o main ./cmd/main.go
```

### Run Tests
```bash
go test ./...
```

### Run with Hot Reload
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with air
air
```

### Format Code
```bash
go fmt ./...
```

### Lint
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

## Project Structure
```
.
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── db/                     # Database initialization
│   ├── handlers/               # HTTP handlers
│   ├── models/                 # Data models
│   └── services/               # Business logic
├── pkg/
│   ├── auth/                   # Authentication (JWT)
│   ├── logger/                 # Logging
│   └── router/                 # Route setup
├── migrations/                 # Database migrations
├── go.mod                      # Go modules
├── Dockerfile                  # Docker configuration
└── .env.example               # Environment template
```

## Troubleshooting

### Database Connection Error
- Ensure MySQL is running
- Check DB credentials in .env
- Verify database name exists

### Port Already in Use
- Change SERVER_PORT in .env
- Or kill process: `lsof -ti:8080 | xargs kill -9`

### Missing Dependencies
```bash
go mod tidy
go mod download
```

## Next Development Steps

1. **Implement Agent Management**
   - Create agent service methods
   - Add agent handlers
   - Write tests

2. **Complete AI Orchestrator**
   - Implement provider adapters (OpenAI, Claude, Gemini)
   - Add caching layer
   - Rate limiting

3. **WebSocket Support**
   - Real-time agent status
   - Call notifications
   - Live chat

4. **Authentication Middleware**
   - JWT validation
   - Tenant isolation
   - RBAC implementation

5. **Testing**
   - Unit tests
   - Integration tests
   - API tests

## Contributing
- Follow Go code standards
- Write tests for new features
- Update documentation
- Use meaningful commit messages
