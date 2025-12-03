# Multi-Tenant Quick Reference

Quick lookup guide for common tasks and endpoints.

## Quick Start

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -d '{"email":"user@example.com","password":"pass"}' \
  -H "Content-Type: application/json"
```

### Save Token
```bash
TOKEN="eyJhbGc..." # From login response
```

## Core Endpoints

### Tenant Info
```bash
# Get current tenant
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/tenant

# Get user count in tenant
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/tenant/users/count
```

### List & Switch
```bash
# List user's tenants
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/tenants

# Switch to tenant 2
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/tenants/2/switch

# List all tenants (admin)
curl http://localhost:8080/api/v1/tenants
```

### Members
```bash
# Add member
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -d '{"email":"user@example.com","role":"member"}' \
  http://localhost:8080/api/v1/tenants/1/members

# Remove member
curl -X DELETE -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/tenants/1/members/user@example.com
```

## Frontend Hooks

### TenantContext
```typescript
const { tenant, tenants, loading } = useTenant()
// Current tenant info and list of tenants
```

### TenantManagementContext
```typescript
const { 
  userTenants, 
  switchTenant, 
  addTenantMember, 
  removeTenantMember 
} = useTenantManagement()
```

## Components

### TenantSwitcher
```tsx
<TenantSwitcher />
// Shows dropdown with user's tenants
```

### TenantInfo
```tsx
<TenantInfo />
// Displays current tenant details
```

## Database Tables

### tenants
```sql
id | name | domain | max_users | status | created_at
```

### tenant_users
```sql
id | tenant_id | user_id | role | joined_at
```

### users (extended)
```sql
current_tenant_id | role | ...
```

## Common Workflows

### Tenant Switching Flow
```
1. GET /api/v1/tenants (list user's tenants)
2. POST /api/v1/tenants/{id}/switch (switch)
3. GET /api/v1/tenant (verify switch)
4. Refresh UI with new context
```

### Adding Members Flow
```
1. User enters email
2. POST /api/v1/tenants/{id}/members
3. User receives invitation
4. User joins tenant
5. User can access tenant
```

### Role Management
```
"admin"  - Full access, can manage members
"member" - Standard access
"viewer" - Read-only access
```

## Error Codes

| Code | HTTP | Meaning |
|------|------|---------|
| TENANT_NOT_FOUND | 404 | Tenant doesn't exist |
| ACCESS_DENIED | 403 | User doesn't have access |
| USER_NOT_FOUND | 404 | User not in system |
| USER_EXISTS | 400 | Already in tenant |
| UNAUTHORIZED | 401 | Token invalid/expired |

## Environment Setup

```bash
# Backend
export DATABASE_URL=postgresql://user:pass@localhost/ai_call_center
export JWT_SECRET=your-secret
export API_PORT=8080

# Frontend
export NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Common Commands

```bash
# Backend setup
go mod download
go run cmd/main.go

# Frontend setup
cd frontend && npm install
npm run dev

# Run tests
go test ./...
npm test

# Build
go build
npm run build
```

## Database Setup

```bash
# Create database
createdb ai_call_center

# Run migrations
psql -U postgres -d ai_call_center < migrations/001_initial_schema.sql

# Check tables
psql -d ai_call_center -c "\dt"
```

## Useful Queries

```sql
-- Check tenants
SELECT * FROM tenants;

-- Check user in tenant
SELECT * FROM tenant_users WHERE user_id = 1;

-- List members of tenant
SELECT u.*, tu.role FROM users u
JOIN tenant_users tu ON u.id = tu.user_id
WHERE tu.tenant_id = 1;

-- Find user's tenants
SELECT t.* FROM tenants t
JOIN tenant_users tu ON t.id = tu.tenant_id
WHERE tu.user_id = 1;
```

## Frontend Structure

```
frontend/
├── contexts/
│   ├── TenantContext.tsx       # State provider
│   └── TenantManagementContext.tsx # Operations
├── components/dashboard/
│   ├── TenantSwitcher.tsx      # Switch UI
│   └── TenantInfo.tsx          # Info display
├── services/
│   └── api.ts                  # API calls
└── hooks/
    └── useAuth.ts              # Auth hook
```

## Backend Structure

```
internal/
├── handlers/
│   └── tenant.go              # HTTP handlers
├── services/
│   └── tenant.go              # Business logic
├── models/
│   └── tenant.go              # Data models
└── middleware/
    └── auth.go                # Authentication
```

## Key Files to Review

1. **Backend Implementation**
   - `internal/handlers/tenant.go` - All handlers
   - `internal/services/tenant.go` - Service logic
   - `pkg/router/router.go` - Route setup

2. **Frontend Implementation**
   - `frontend/contexts/TenantContext.tsx` - State
   - `frontend/services/api.ts` - API methods
   - `frontend/components/dashboard/TenantSwitcher.tsx` - UI

3. **Database**
   - `migrations/001_initial_schema.sql` - Schema

4. **Documentation**
   - `MULTI_TENANT_FEATURES.md` - Full reference
   - `MULTI_TENANT_API_TESTING.md` - Testing
   - `MULTI_TENANT_INTEGRATION_CHECKLIST.md` - Checklist

## Testing Endpoints

### Test Script
```bash
#!/bin/bash
API=http://localhost:8080
TOKEN="your-token"

echo "Testing tenant endpoints..."
curl -s -H "Authorization: Bearer $TOKEN" $API/api/v1/tenant | jq '.'
curl -s -H "Authorization: Bearer $TOKEN" $API/api/v1/tenants | jq '.'
curl -s -H "Authorization: Bearer $TOKEN" $API/api/v1/tenant/users/count | jq '.'
```

## Performance Tips

1. **Cache tenant context** - Don't refetch unless needed
2. **Use indexes** - Ensure DB indexes on tenant_id
3. **Lazy load** - Load tenant data on demand
4. **Batch operations** - Group member operations
5. **Monitor queries** - Check slow query logs

## Debugging

### Check Token Claims
```bash
# Decode JWT (jq required)
TOKEN="eyJhbGc..."
echo $TOKEN | cut -d'.' -f2 | base64 -d | jq '.'
```

### Check Database Connection
```bash
psql postgresql://user:pass@localhost/ai_call_center -c "SELECT version();"
```

### Enable Debug Logging
```bash
export LOG_LEVEL=debug
go run cmd/main.go
```

## Useful Links

- **Go Docs**: https://golang.org/doc/
- **Next.js Docs**: https://nextjs.org/docs
- **PostgreSQL Docs**: https://www.postgresql.org/docs/
- **JWT.io**: https://jwt.io/
- **Gorilla Mux**: https://github.com/gorilla/mux

## Checklists

### Before Deployment
- [ ] Database migrated
- [ ] Environment variables set
- [ ] Tests passing
- [ ] Endpoints responding
- [ ] Tenant isolation verified
- [ ] JWT tokens valid
- [ ] CORS configured
- [ ] SSL configured

### After Deployment
- [ ] Health check passes
- [ ] Can login
- [ ] Can switch tenant
- [ ] Can add/remove members
- [ ] Data properly isolated
- [ ] Logs monitoring active

## Support

For detailed information:
- Features → MULTI_TENANT_FEATURES.md
- Testing → MULTI_TENANT_API_TESTING.md
- Integration → MULTI_TENANT_INTEGRATION_CHECKLIST.md
- Summary → MULTI_TENANT_IMPLEMENTATION_SUMMARY.md

---

**Version**: 2.0
**Last Updated**: 2024
**Status**: Production Ready
4. Submit
5. Tenant created, user is admin
6. Login → Dashboard
```

### Flow 2: Join Existing Tenant
```
1. Register page
2. Select "Join Existing Tenant"
3. Enter tenant code
4. Submit
5. User added to existing tenant
6. Login → Dashboard
```

## 🔄 Tenant Switching

```
Dashboard → Sidebar → "Switch Tenant" → Select → API call → Dashboard updates
```

## 📱 UI Layout

```
┌─────────────────────────────────────┐
│ ≡ Menu                              │
├─────────────────────────────────────┤
│ Current Tenant: Acme Corp           │
│ ▼ Switch Tenant (3 available)       │
│   [Acme Corp] ✓                     │
│   [Tech Inc]                        │
│   [Startup Co]                      │
├─────────────────────────────────────┤
│ 📊 Dashboard                        │
│ 🏢 Tenants                          │
│ ⚙️ Settings                         │
│ 🚪 Logout                           │
└─────────────────────────────────────┘
```

## 🔐 Data Persistence

### localStorage
- `auth_token` - JWT with tenant_id
- `current_tenant_id` - Active tenant ID

### Context
- Current tenant details
- List of all user's tenants
- Current tenant ID

## 🛠️ API Endpoints (Frontend Calls)

| Method | Endpoint | Purpose | Status |
|--------|----------|---------|--------|
| POST | `/api/v1/tenants` | Create tenant | ✅ Frontend ready |
| GET | `/api/v1/tenants` | List tenants | ✅ Frontend ready |
| POST | `/api/v1/tenants/{id}/switch` | Switch tenant | ✅ Frontend ready |
| POST | `/api/v1/tenants/{id}/members` | Add member | ✅ Frontend ready |
| DELETE | `/api/v1/tenants/{id}/members/{email}` | Remove member | ✅ Frontend ready |

**Note:** Frontend is ready for these endpoints. Backend needs to implement them.

## 🚀 How to Use

### As a New User
1. Go to `/auth/register`
2. Choose how to setup tenant
3. Enter details and register
4. Login and start using dashboard

### As Admin
1. Go to `/dashboard/tenants`
2. Click "+ Create Tenant" to add more
3. Invite team members (when backend ready)
4. Switch between tenants anytime

### Switching Tenants
1. Look for "Current Tenant" in sidebar
2. Click "Switch Tenant"
3. Select from dropdown
4. Dashboard updates automatically

## ✨ Features

✅ Multi-tenant registration
✅ Tenant switching
✅ Tenant management page
✅ Create new tenant
✅ View tenant details
✅ Member management (UI ready)
✅ Role badges (Admin/Member/Viewer)
✅ Error handling
✅ Loading states
✅ Toast notifications

## 🔌 Backend Integration Status

| Feature | Frontend | Backend |
|---------|----------|---------|
| Create tenant | ✅ | 🔴 |
| List tenants | ✅ | ✅ |
| Switch tenant | ✅ | 🔴 |
| Add member | ✅ | 🔴 |
| Remove member | ✅ | 🔴 |
| Get tenant | ✅ | ✅ |

## 📋 Build Status

```
✅ Frontend builds successfully
✅ No TypeScript errors
✅ All components working
✅ Ready for testing
🔴 Waiting for backend endpoints
```

## 🐛 Troubleshooting

### Can't switch tenants?
- Check if backend has switch endpoint
- Verify JWT token has tenant_id
- Check browser console for errors

### Registration fails?
- Verify tenant code is correct (join flow)
- Check tenant name not empty (create flow)
- Look for toast notification with error

### Can't see other tenants?
- Verify user is member of multiple tenants
- Check API returns all user's tenants
- Refresh page to reload

## 📚 Related Documentation

- `MULTI_TENANT_USER_GUIDE.md` - Complete user guide
- `BACKEND_TENANT_UPDATES.md` - Backend implementation guide
- `FRONTEND_TENANT_UI_GUIDE.md` - UI component guide
- `FRONTEND_TENANT_IMPLEMENTATION.md` - Full implementation details

## 🎓 Code Examples

### Using Tenant Context
```tsx
import { useTenantManagement } from '@/contexts/TenantManagementContext'

export function MyComponent() {
  const { userTenants, currentTenantId, switchTenant } = useTenantManagement()
  
  return (
    <div>
      Current: {currentTenantId}
      Total: {userTenants.length}
    </div>
  )
}
```

### Switching Tenant
```tsx
const { switchTenant } = useTenantManagement()

const handleSwitch = () => {
  switchTenant('tenant-123')
  // After API call succeeds
}
```

### Creating Tenant
```tsx
const { createTenant } = useTenantManagement()

const handleCreate = async () => {
  const newTenant = await createTenant('My Company', 'mycompany.com')
  console.log('Created:', newTenant.id)
}
```

## 📞 Next Actions

### For Frontend Developers
1. Test registration flows
2. Verify all UI works
3. Test error cases
4. Review documentation

### For Backend Developers
1. Implement switch endpoint
2. Implement member endpoints
3. Update database schema
4. Add authorization checks
5. Test with frontend

### For QA/Testing
1. Test new tenant registration
2. Test join existing tenant
3. Test tenant switching
4. Test member operations
5. Test error handling

---

**Last Updated:** After complete multi-tenant frontend implementation
**Status:** ✅ Frontend Complete, 🔴 Backend Pending
**Ready for:** Backend integration and end-to-end testing
