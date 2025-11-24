# Multi-Tenant Features Implementation Guide

## Overview

This document outlines all the multi-tenant features implemented in the AI Call Center application, including tenant management, user management, and API endpoints.

## Table of Contents

1. [Architecture](#architecture)
2. [Database Schema](#database-schema)
3. [API Endpoints](#api-endpoints)
4. [Frontend Components](#frontend-components)
5. [Context Providers](#context-providers)
6. [Services](#services)
7. [User Workflows](#user-workflows)
8. [Configuration](#configuration)

## Architecture

The multi-tenant system is built on a **database-per-tenant or shared-database-with-tenant-isolation** model:

- **Backend (Go)**: Handles tenant isolation at the database level and validates user permissions
- **Frontend (Next.js)**: Manages tenant context and switching
- **Database (PostgreSQL)**: Stores tenant information and enforces relationships

### Key Components

```
Frontend
  ├── TenantContext (provides current tenant)
  ├── TenantManagementContext (provides tenant operations)
  ├── TenantSwitcher (UI for switching tenants)
  └── TenantInfo (displays current tenant info)

Backend
  ├── TenantService (business logic)
  ├── TenantHandler (HTTP handlers)
  └── Middleware (authentication & tenant validation)

Database
  ├── tenants table
  ├── tenant_users table
  └── tenant_configs table
```

## Database Schema

### Tenants Table

```sql
CREATE TABLE tenants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE,
    description TEXT,
    max_users INT DEFAULT 10,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB
);
```

### Tenant Users Table

```sql
CREATE TABLE tenant_users (
    id SERIAL PRIMARY KEY,
    tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, user_id)
);
```

### Tenant Configs Table

```sql
CREATE TABLE tenant_configs (
    id SERIAL PRIMARY KEY,
    tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
    key VARCHAR(255),
    value TEXT,
    UNIQUE(tenant_id, key)
);
```

## API Endpoints

### 1. Tenant Information Endpoints

#### Get Current Tenant Info
```
GET /api/v1/tenant
Authorization: Bearer {token}

Response:
{
  "id": 1,
  "name": "ACME Corp",
  "domain": "acme.example.com",
  "description": "Main tenant",
  "max_users": 50,
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Get Tenant User Count
```
GET /api/v1/tenant/users/count
Authorization: Bearer {token}

Response:
{
  "count": 15,
  "max_users": 50
}
```

### 2. Tenant Listing and Management

#### List All Tenants (Admin)
```
GET /api/v1/tenants

Response:
[
  {
    "id": 1,
    "name": "ACME Corp",
    "domain": "acme.example.com",
    "status": "active"
  },
  ...
]
```

#### Get User's Tenants
```
GET /api/v1/tenants
Authorization: Bearer {token}

Response:
[
  {
    "id": 1,
    "name": "ACME Corp",
    "role": "admin",
    "joined_at": "2024-01-01T00:00:00Z"
  },
  ...
]
```

### 3. Tenant Switching

#### Switch Tenant
```
POST /api/v1/tenants/{id}/switch
Authorization: Bearer {token}
Content-Type: application/json

Request Body:
{}

Response:
{
  "success": true,
  "message": "Tenant switched successfully",
  "tenant_id": 1
}
```

**Note**: This endpoint updates the user's active tenant session.

### 4. Tenant Member Management

#### Add Tenant Member
```
POST /api/v1/tenants/{id}/members
Authorization: Bearer {token}
Content-Type: application/json

Request Body:
{
  "email": "user@example.com",
  "role": "user"  // or "admin", "viewer"
}

Response:
{
  "success": true,
  "message": "User added to tenant",
  "user_email": "user@example.com"
}
```

#### Remove Tenant Member
```
DELETE /api/v1/tenants/{id}/members/{email}
Authorization: Bearer {token}

Response:
{
  "success": true,
  "message": "User removed from tenant"
}
```

## Frontend Components

### TenantContext

Provides the current tenant information to the entire application.

```typescript
interface TenantContextType {
  tenant: Tenant | null
  tenants: Tenant[]
  loadTenant: (tenantId: number) => Promise<void>
  loading: boolean
}

// Usage
const { tenant } = useTenant()
```

### TenantManagementContext

Provides tenant management operations (create, switch, add/remove members).

```typescript
interface TenantManagementContextType {
  userTenants: UserTenant[]
  currentTenantId: string | null
  switchTenant: (tenantId: string) => void
  createTenant: (name: string, domain?: string) => Promise<Tenant>
  addTenantMember: (tenantId: string, email: string, role: string) => Promise<void>
  removeTenantMember: (tenantId: string, email: string) => Promise<void>
  loading: boolean
  error: string | null
}

// Usage
const { switchTenant, addTenantMember } = useTenantManagement()
```

### TenantSwitcher Component

Dropdown UI component that allows users to switch between their tenants.

```typescript
<TenantSwitcher />

// Features:
// - Shows current tenant
// - Lists all user's tenants
// - Handles tenant switching
// - Shows user count / max users
```

### TenantInfo Component

Displays detailed information about the current tenant.

```typescript
<TenantInfo />

// Displays:
// - Tenant name
// - Domain
// - User count
// - Creation date
```

## Context Providers

### Setup in Layout

```typescript
// app/layout.tsx
import { TenantProvider } from '@/contexts/TenantContext'
import { TenantManagementProvider } from '@/contexts/TenantManagementContext'
import { AuthProvider } from '@/components/providers/AuthProvider'

export default function RootLayout({ children }) {
  return (
    <AuthProvider>
      <TenantProvider>
        <TenantManagementProvider>
          {children}
        </TenantManagementProvider>
      </TenantProvider>
    </AuthProvider>
  )
}
```

### Provider Nesting Order

1. **AuthProvider**: Must be first (handles authentication)
2. **TenantProvider**: Second (provides tenant context)
3. **TenantManagementProvider**: Third (depends on auth context)

## Services

### Tenant Service (Backend)

Location: `internal/services/tenant.go`

```go
type TenantService interface {
    GetTenant(ctx context.Context, tenantID int) (*models.Tenant, error)
    ListTenants(ctx context.Context) ([]models.Tenant, error)
    CreateTenant(ctx context.Context, tenant *models.Tenant) error
    UpdateTenant(ctx context.Context, tenant *models.Tenant) error
    DeleteTenant(ctx context.Context, tenantID int) error
    
    // User management
    AddTenantUser(ctx context.Context, tenantID, userID int, role string) error
    RemoveTenantUser(ctx context.Context, tenantID, userID int) error
    GetTenantUsers(ctx context.Context, tenantID int) ([]models.User, error)
    
    // User's tenants
    GetUserTenants(ctx context.Context, userID int) ([]models.Tenant, error)
    SwitchTenant(ctx context.Context, userID, tenantID int) error
}
```

### Tenant Service (Frontend)

Location: `frontend/services/api.ts`

```typescript
export const tenantService = {
  async getTenantInfo(): Promise<Tenant>,
  async getTenantUserCount(): Promise<{ count: number; max_users: number }>,
  async listTenants(): Promise<Tenant[]>,
  async getUserTenants(): Promise<UserTenant[]>,
  async switchTenant(tenantId: number): Promise<{ success: boolean }>,
  async addTenantMember(tenantId: number, email: string, role: string): Promise<void>,
  async removeTenantMember(tenantId: number, email: string): Promise<void>,
}
```

## User Workflows

### 1. User Registration

```
1. User goes to /auth/register
2. Enters email, password, and selects/creates tenant
3. Backend creates user and adds to tenant with "owner" role
4. User is logged in and redirected to /dashboard
```

### 2. Switching Tenants

```
1. User clicks on TenantSwitcher dropdown
2. Selects desired tenant from list
3. Frontend calls POST /api/v1/tenants/{id}/switch
4. Backend updates session tenant_id
5. Frontend refreshes to reload tenant context
6. User sees new tenant's data
```

### 3. Adding Team Members

```
1. Tenant admin goes to tenant settings
2. Clicks "Add Member"
3. Enters member email and selects role
4. Frontend calls POST /api/v1/tenants/{id}/members
5. Backend invites user to tenant
6. User receives invitation (email)
7. User can join tenant from invitation
```

### 4. Removing Team Members

```
1. Tenant admin goes to members list
2. Clicks remove on member
3. Frontend calls DELETE /api/v1/tenants/{id}/members/{email}
4. Backend removes user from tenant
5. User no longer has access to that tenant
```

## Configuration

### Environment Variables

Required environment variables for multi-tenant support:

```bash
# Backend
DATABASE_URL=postgresql://user:password@localhost/dbname
JWT_SECRET=your-secret-key
API_PORT=8080

# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### Database Migration

Run migrations to set up tenant tables:

```bash
# Using provided migration scripts
./scripts/setup-dev.sh

# Or manually
psql -U postgres -d ai_call_center -f migrations/001_initial_schema.sql
```

## Security Considerations

### Tenant Isolation

1. **Row-Level Security (RLS)**: Implement at database level
   ```sql
   CREATE POLICY tenant_isolation ON tenants
   USING (tenant_id = current_user_id's_tenant)
   ```

2. **Middleware Validation**: All endpoints validate user belongs to tenant
   ```go
   // in auth.go middleware
   validateTenantAccess(userID, tenantID)
   ```

3. **JWT Claims**: Include tenant_id in token
   ```go
   claims := jwt.MapClaims{
       "sub": user.ID,
       "tenant_id": user.CurrentTenantID,
   }
   ```

### Rate Limiting

Consider implementing rate limiting per tenant:
```go
// Rate limit by tenant
rateLimiter := newRateLimiter(tenantID)
```

### Audit Logging

Log all tenant operations:
```go
log.WithFields(logrus.Fields{
    "tenant_id": tenantID,
    "user_id": userID,
    "action": "switch_tenant",
}).Info("Tenant operation")
```

## Testing

### Backend Tests

```bash
# Run tenant service tests
go test -v ./internal/services -run TestTenant*

# Run tenant handler tests
go test -v ./internal/handlers -run TestTenant*
```

### Frontend Tests

```bash
# Run tenant context tests
npm test -- TenantContext.test.ts

# Run tenant component tests
npm test -- TenantSwitcher.test.tsx
```

## Troubleshooting

### Tenant not switching

1. Check JWT token contains tenant_id
2. Verify user belongs to target tenant
3. Check database constraint on tenant_users

### Missing tenant in list

1. Verify user is added to tenant_users table
2. Check tenant is marked as "active"
3. Verify tenant creation completed successfully

### "Tenant not found" error

1. Verify tenant exists in database
2. Check tenant_id format (should be integer)
3. Verify user has access to tenant

## Future Enhancements

1. **Tenant Customization**
   - Custom branding per tenant
   - Tenant-specific features
   - Custom domain support

2. **Advanced Permissions**
   - More granular role-based access
   - Resource-level permissions
   - Delegation chains

3. **Analytics**
   - Per-tenant usage statistics
   - Audit trail per tenant
   - Tenant health dashboard

4. **Multi-region Support**
   - Tenant data localization
   - Region-specific databases
   - Data residency compliance

## Support & Documentation

For more information:
- Backend Implementation: See `BACKEND_TENANT_UPDATES.md`
- Frontend Implementation: See `FRONTEND_TENANT_IMPLEMENTATION.md`
- API Reference: See `QUICK_REFERENCE_TENANT.md`
- User Guide: See `MULTI_TENANT_USER_GUIDE.md`

