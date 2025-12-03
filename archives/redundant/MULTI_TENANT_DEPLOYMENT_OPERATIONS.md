# Multi-Tenant Deployment & Operations Guide

Complete guide for deploying and operating the multi-tenant AI Call Center application.

## Table of Contents

1. [Pre-Deployment](#pre-deployment)
2. [Database Setup](#database-setup)
3. [Backend Deployment](#backend-deployment)
4. [Frontend Deployment](#frontend-deployment)
5. [Configuration](#configuration)
6. [Health Checks](#health-checks)
7. [Monitoring](#monitoring)
8. [Troubleshooting](#troubleshooting)
9. [Scaling](#scaling)

## Pre-Deployment

### System Requirements

**Backend**
- Go 1.18+
- PostgreSQL 12+
- 2GB RAM minimum
- 4GB disk space

**Frontend**
- Node.js 16+
- npm/yarn
- 2GB disk space (node_modules)

**Network**
- Port 8080 available (backend)
- Port 3000 available (frontend dev) / 80,443 (production)
- Outbound HTTPS access for email

### Pre-Deployment Checklist

- [ ] Review all documentation
- [ ] Verify all tests pass locally
- [ ] Create backup of any existing data
- [ ] Prepare environment variables
- [ ] Verify database credentials
- [ ] Check SSL certificates (if using HTTPS)
- [ ] Review firewall rules
- [ ] Plan rollback procedure

## Database Setup

### PostgreSQL Installation

#### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### macOS (with Homebrew)
```bash
brew install postgresql@14
brew services start postgresql@14
```

#### Docker
```bash
docker run --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=ai_call_center \
  -p 5432:5432 \
  -d postgres:14
```

### Database Creation

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE ai_call_center;

# Create user
CREATE USER app_user WITH PASSWORD 'secure_password';

# Grant permissions
GRANT ALL PRIVILEGES ON DATABASE ai_call_center TO app_user;
```

### Running Migrations

```bash
# Navigate to project root
cd /path/to/project

# Run migration
psql -U app_user -d ai_call_center -h localhost < migrations/001_initial_schema.sql

# Verify tables created
psql -U app_user -d ai_call_center -c "\dt"
```

Expected output:
```
                List of relations
 Schema |        Name         | Type  |  Owner
--------+---------------------+-------+---------
 public | agents              | table | app_user
 public | calls               | table | app_user
 public | leads               | table | app_user
 public | tenants             | table | app_user
 public | tenant_configs      | table | app_user
 public | tenant_users        | table | app_user
 public | users               | table | app_user
 public | password_resets     | table | app_user
(8 rows)
```

### Create Indexes

```bash
# Connect to database
psql -U app_user -d ai_call_center

# Create indexes
CREATE INDEX idx_tenant_users_tenant_id ON tenant_users(tenant_id);
CREATE INDEX idx_tenant_users_user_id ON tenant_users(user_id);
CREATE INDEX idx_users_current_tenant_id ON users(current_tenant_id);
CREATE INDEX idx_tenants_domain ON tenants(domain);

# Verify indexes
\di
```

## Backend Deployment

### Build Process

```bash
# Navigate to project root
cd /path/to/project

# Download dependencies
go mod download

# Build binary
go build -o bin/api-server cmd/main.go

# Verify binary
./bin/api-server --version
```

### Environment Configuration

Create `.env` file or set system environment variables:

```bash
# Database
DATABASE_URL=postgresql://app_user:secure_password@localhost:5432/ai_call_center

# Authentication
JWT_SECRET=your-very-secret-key-generate-with-openssl-rand-base64-32
JWT_EXPIRATION=24h

# Server
API_PORT=8080
API_HOST=0.0.0.0

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Email (for password reset)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@example.com

# Frontend
FRONTEND_URL=https://yourdomain.com
```

Generate JWT secret:
```bash
openssl rand -base64 32
```

### Running the Backend

#### Development
```bash
go run cmd/main.go
```

#### Production (using systemd)

Create `/etc/systemd/system/api-server.service`:

```ini
[Unit]
Description=AI Call Center API Server
After=network.target postgresql.service
Wants=postgresql.service

[Service]
Type=simple
User=apiserver
WorkingDirectory=/opt/api-server
EnvironmentFile=/opt/api-server/.env
ExecStart=/opt/api-server/bin/api-server
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable api-server
sudo systemctl start api-server
sudo systemctl status api-server
```

#### Production (using Docker)

```bash
docker build -t ai-callcenter-api:latest .

docker run -d \
  --name api-server \
  -p 8080:8080 \
  --env-file .env \
  --link postgres:postgres \
  ai-callcenter-api:latest
```

### Verify Backend

```bash
# Check if server is running
curl http://localhost:8080/health

# Expected response
# {"status":"healthy","timestamp":"2024-01-01T00:00:00Z"}
```

## Frontend Deployment

### Build Process

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Build for production
npm run build

# Verify build
ls -la .next
```

### Environment Configuration

Create `frontend/.env.local`:

```bash
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
NEXT_PUBLIC_APP_NAME=AI Call Center
NEXT_PUBLIC_APP_VERSION=1.0.0
```

### Deployment Options

#### Option 1: Static Hosting (Vercel)

```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
cd frontend
vercel

# Configure environment variables in Vercel dashboard
```

#### Option 2: Self-Hosted (Node.js)

Create `/etc/systemd/system/nextjs-app.service`:

```ini
[Unit]
Description=Next.js Frontend
After=network.target

[Service]
Type=simple
User=nextjs
WorkingDirectory=/opt/frontend
EnvironmentFile=/opt/frontend/.env.local
ExecStart=/usr/bin/npm start
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable nextjs-app
sudo systemctl start nextjs-app
```

#### Option 3: Docker Container

```bash
docker build -t ai-callcenter-frontend:latest -f Dockerfile.frontend .

docker run -d \
  --name frontend \
  -p 3000:3000 \
  --env-file frontend/.env.local \
  ai-callcenter-frontend:latest
```

#### Option 4: Reverse Proxy (Nginx)

Create `/etc/nginx/sites-available/ai-callcenter`:

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Static files cache
    location /_next/static {
        expires 365d;
        add_header Cache-Control "public, immutable";
    }
}

server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/ai-callcenter /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

## Configuration

### SSL/TLS Setup

#### Using Let's Encrypt

```bash
# Install Certbot
sudo apt-get install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot certonly --nginx -d yourdomain.com

# Auto-renewal
sudo certbot renew --dry-run
```

### CORS Configuration

Update backend CORS middleware in `internal/middleware/cors.go`:

```go
// Allow specific origins in production
AllowedOrigins: []string{
    "https://yourdomain.com",
    "https://www.yourdomain.com",
}
```

### Security Headers

Add to Nginx configuration:

```nginx
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "no-referrer-when-downgrade" always;
```

## Health Checks

### Backend Health Check

```bash
# Basic health check
curl -i http://localhost:8080/health

# Expected: 200 OK with JSON response
```

### Database Health Check

```bash
# Test database connection
psql -U app_user -d ai_call_center -c "SELECT NOW();"

# Or via backend
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/tenant
```

### Frontend Health Check

```bash
# Test frontend
curl -I http://localhost:3000

# Expected: 200 OK
```

### Full Workflow Test

```bash
#!/bin/bash
echo "Testing multi-tenant API..."

API="https://yourdomain.com"

# 1. Health check
echo "1. Health check..."
curl -s $API/health | jq '.'

# 2. List tenants
echo "2. List tenants..."
curl -s $API/api/v1/tenants | jq '.'

# 3. Login
echo "3. Login..."
TOKEN=$(curl -s -X POST $API/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' | jq -r '.token')

echo "Token: $TOKEN"

# 4. Get current tenant
echo "4. Get current tenant..."
curl -s -H "Authorization: Bearer $TOKEN" \
  $API/api/v1/tenant | jq '.'

echo "All tests completed!"
```

## Monitoring

### Logging

#### Verify Logs
```bash
# Backend logs (systemd)
journalctl -u api-server -f

# Frontend logs (Node.js)
journalctl -u nextjs-app -f

# Docker logs
docker logs -f api-server
docker logs -f frontend
```

### Key Metrics to Monitor

1. **Response Time**
   - API endpoint latency
   - Database query time
   - Frontend load time

2. **Error Rates**
   - 4xx errors (client errors)
   - 5xx errors (server errors)
   - Database errors

3. **Resource Usage**
   - CPU usage
   - Memory usage
   - Disk space
   - Database connections

4. **Business Metrics**
   - Active tenants
   - Active users
   - API calls per tenant
   - Tenant user count

### Alert Configuration

Set up alerts for:
- High error rate (> 1%)
- High response time (> 1000ms)
- Database connection pool exhausted
- Disk space < 10%
- CPU usage > 80%

## Troubleshooting

### Backend Issues

**Issue**: Server won't start
```bash
# Check port is available
sudo lsof -i :8080

# Check database connection
psql $DATABASE_URL -c "SELECT 1"

# Check environment variables
env | grep DATABASE_URL
```

**Issue**: "Connection refused" error
```bash
# Verify database is running
sudo systemctl status postgresql

# Check database credentials
psql -U app_user -d ai_call_center -c "\d"
```

**Issue**: High memory usage
```bash
# Check Go runtime metrics
curl http://localhost:8080/debug/pprof/heap

# Monitor continuously
watch -n 1 'top -b -n 1 | grep api-server'
```

### Frontend Issues

**Issue**: Build fails
```bash
# Clear cache and rebuild
rm -rf node_modules package-lock.json
npm install
npm run build
```

**Issue**: Cannot connect to API
```bash
# Check NEXT_PUBLIC_API_URL
grep NEXT_PUBLIC_API_URL .env.local

# Test API from browser console
fetch('https://yourdomain.com/api/v1/health')
```

### Database Issues

**Issue**: Slow queries
```sql
-- Check slow query log
SELECT query, mean_time, calls
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
```

**Issue**: Connection pool exhausted
```sql
-- Check active connections
SELECT datname, count(*) FROM pg_stat_activity
GROUP BY datname;

-- Close idle connections
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE state = 'idle' AND query_start < now() - interval '30 minutes';
```

## Scaling

### Horizontal Scaling

#### Multiple API Servers

```bash
# Server 1
API_PORT=8080 go run cmd/main.go

# Server 2
API_PORT=8081 go run cmd/main.go

# Load balancer (Nginx)
upstream api_backend {
    server localhost:8080;
    server localhost:8081;
}

location /api {
    proxy_pass http://api_backend;
}
```

#### Database Replication

For production:
```sql
-- Setup read replicas
-- Use PostgreSQL streaming replication
-- Configure connection pooling (PgBouncer)
```

### Caching Strategy

#### Redis Caching
```go
// Cache tenant data
const cacheKey = "tenant:" + tenantID
cache.Get(cacheKey)
```

#### CDN for Static Files
```nginx
# Cache Next.js static files
location /_next/static {
    add_header Cache-Control "public, immutable, max-age=31536000";
}
```

### Performance Optimization

1. **Database**
   - Add indexes on frequently queried columns
   - Use connection pooling
   - Implement query caching

2. **Application**
   - Use context caching
   - Implement request deduplication
   - Add gzip compression

3. **Frontend**
   - Use Code splitting
   - Implement lazy loading
   - Optimize images

## Backup & Recovery

### Database Backup

```bash
# Full backup
pg_dump -U app_user -d ai_call_center > backup.sql

# Automated daily backup
0 2 * * * pg_dump -U app_user -d ai_call_center > /backups/db-$(date +\%Y\%m\%d).sql
```

### Restore from Backup

```bash
# Restore database
psql -U app_user -d ai_call_center < backup.sql
```

## Support & Maintenance

### Regular Maintenance Tasks

- [ ] Weekly: Check log files for errors
- [ ] Weekly: Verify backups are working
- [ ] Monthly: Update dependencies
- [ ] Monthly: Review and optimize slow queries
- [ ] Quarterly: Security audit
- [ ] Annually: Capacity planning

### Emergency Contacts

- Database Administrator: [contact]
- Infrastructure Team: [contact]
- On-call Support: [contact]

---

**Version**: 1.0
**Last Updated**: 2024
**Status**: Production Ready

