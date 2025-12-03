# 🚀 Getting Started with VYOM ERP

**Last Updated:** December 3, 2025

---

## Quick Navigation

### 📖 Documentation

| Document | Purpose | When to Use |
|----------|---------|------------|
| **README.md** | Project overview | First time? Start here |
| **API_DOCUMENTATION.md** | Complete API reference | Building integrations |
| **QUICK_START.md** | 5-minute setup | Getting dev environment running |
| **QUICK_REFERENCE.md** | API endpoints cheat sheet | Quick lookup of endpoints |
| **DEVELOPMENT.md** | Dev setup guide | Setting up local development |
| **QUICK_START_TESTING.md** | Testing guide | Running tests |

---

## 5-Minute Quick Start

### 1. Clone & Install
```bash
cd d:/VYOMTECH-ERP
cd frontend && npm install
```

### 2. Start Development
```bash
npm run dev
# Frontend: http://localhost:3000
```

### 3. Start Backend
```bash
cd d:/VYOMTECH-ERP
go run cmd/main.go
# API: http://localhost:8080
```

---

## Common API Calls

### Authentication
```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Response includes: access_token, refresh_token
```

### Using API with Token
```bash
# Get current user
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>" \
  -H "X-Tenant-ID: <tenant_id>"
```

### Common Headers
```
Authorization: Bearer <jwt_token>
X-Tenant-ID: <tenant_uuid>
Content-Type: application/json
```

---

## Project Structure Quick Reference

```
d:/VYOMTECH-ERP/
├── frontend/                 # Next.js React app (PORT 3000)
│   ├── app/                  # 24 pages
│   ├── components/           # 30+ reusable components
│   ├── services/api.ts       # 65+ API methods
│   ├── hooks/                # 9 custom React hooks
│   └── package.json          # Dependencies
│
├── cmd/main.go               # Backend entry point
├── internal/
│   ├── handlers/             # 25+ HTTP handlers
│   ├── services/             # Business logic
│   ├── models/               # Data models
│   ├── middleware/           # Auth, logging
│   └── db/                   # Database config
│
├── migrations/               # 22 SQL migration files
├── go.mod                    # Go dependencies
└── Dockerfile               # Containerization
```

---

## Key Technologies

| Layer | Technology |
|-------|-----------|
| **Frontend** | Next.js 16, React 19, TypeScript, Tailwind CSS |
| **Backend** | Go 1.19+, GORM, Gorilla Mux |
| **Database** | MySQL/PostgreSQL, 22 migrations |
| **Auth** | JWT + OAuth2 |
| **Real-time** | WebSocket (socket.io) |
| **State** | Zustand, React Query |

---

## Frontend Pages (24 Total)

**Dashboard & Management:**
- Dashboard, Accounts, Agents, Audit Logs, Categories
- Cost Centers, Crm Accounts, Crm Deals, Crm Interactions
- Custom Fields, Customers, Departments, Employees
- GL Accounts, HR, Inventory, Marketing, Presales
- Projects, Purchase, Real Estate, Reports, Sales
- Scheduled Tasks, Tenants, Units, Users, Workflows

---

## Backend Endpoints Overview

**15+ Categories with 200+ Endpoints:**
- Authentication, Users, Customers
- Sales, Purchase, Inventory
- Accounting, Banking, HR/Payroll
- Projects, Call Center, Communication
- Reports, Workflows, RBAC, Settings

See **API_DOCUMENTATION.md** for complete reference.

---

## Environment Setup

### Required Environment Variables
```bash
# .env file
DATABASE_URL=postgres://user:pass@localhost:5432/vyomerp
JWT_SECRET=your_secret_key
REDIS_URL=redis://localhost:6379
PORT=8080
```

### Database Setup
```bash
# Run migrations
go run cmd/main.go migrate

# Or with Docker
docker-compose up
```

---

## Common Commands

### Frontend
```bash
npm install              # Install dependencies
npm run dev             # Start dev server
npm run build           # Build for production
npm run lint            # Run linter
npm test                # Run tests
```

### Backend
```bash
go run cmd/main.go      # Run server
go test ./...           # Run tests
go build -o ./bin/main  # Build binary
```

### Docker
```bash
docker-compose up       # Start all services
docker-compose down     # Stop all services
docker-compose logs -f  # View logs
```

---

## API Testing

### Using cURL
```bash
# Create customer
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "ABC Corp",
    "email": "contact@abc.com",
    "phone": "+1-800-123-4567"
  }'
```

### Using Postman
1. Import API_DOCUMENTATION.md endpoints
2. Set up variables: `{{base_url}}`, `{{token}}`, `{{tenant_id}}`
3. Use pre-built request collections

### Frontend API Client
```typescript
// frontend/services/api.ts
import { api } from '@/services/api';

// Create customer
const customer = await api.customers.create({
  name: 'ABC Corp',
  email: 'contact@abc.com'
});

// Get customers
const customers = await api.customers.list({
  page: 1,
  per_page: 50
});
```

---

## React Hooks

**Available Custom Hooks:**
```typescript
// In frontend/hooks/
- useAuth()               // Authentication
- useCustomers()          // Customer data
- useSales()              // Sales orders
- useProjects()           // Projects
- useEmployees()          // HR data
- useInventory()          // Inventory
- useAccounting()         // GL & accounts
- usePagination()         // Pagination
- useNotification()       // Toast/alerts
```

---

## Common Issues & Solutions

### Frontend Won't Build
```bash
# Clear cache and reinstall
rm -rf node_modules .next
npm install
npm run build
```

### Backend Won't Start
```bash
# Check database connection
go run cmd/main.go

# Check port availability
netstat -an | grep 8080
```

### Migrations Not Applied
```bash
# Run migrations manually
go run cmd/main.go migrate

# Check migration status
go run cmd/main.go migrate status
```

---

## Debugging

### Frontend Debugging
```typescript
// Browser console
localStorage.getItem('auth_token')
sessionStorage.getItem('user')

// Next.js dev tools
npm run dev -- --debug
```

### Backend Debugging
```go
// Enable debug logging
export LOG_LEVEL=debug
go run cmd/main.go

// Using delve debugger
dlv debug ./cmd/main.go
```

---

## Performance Tips

1. **Frontend:**
   - Use React.memo for expensive components
   - Implement code splitting for large pages
   - Cache API responses with React Query

2. **Backend:**
   - Use indexes on frequently queried columns
   - Implement pagination for large datasets
   - Cache DB queries with Redis

3. **Database:**
   - Monitor slow queries
   - Optimize N+1 queries
   - Use connection pooling

---

## Deployment

### Docker
```bash
# Build image
docker build -t vyomerp:latest .

# Run container
docker run -p 8080:8080 vyomerp:latest
```

### Kubernetes
```bash
# Apply manifests
kubectl apply -f k8s/

# Check status
kubectl get pods -n default
```

### Environment Variables for Production
```bash
DATABASE_URL=prod_db_url
JWT_SECRET=prod_secret
REDIS_URL=prod_redis
NODE_ENV=production
```

---

## Support & Resources

- **Documentation:** See `README.md`, `API_DOCUMENTATION.md`
- **Code Examples:** Check `frontend/services/api.ts`
- **Migrations:** Review `migrations/` folder
- **Settings:** Configure in `internal/config/`

---

## Next Steps

1. Read **README.md** for project overview
2. Follow **QUICK_START.md** for setup
3. Check **API_DOCUMENTATION.md** for endpoints
4. Browse **frontend/app/** for UI examples
5. Review **internal/handlers/** for API patterns

---

**Happy coding! 🚀**

For issues or questions, refer to the comprehensive documentation files in the root directory.
