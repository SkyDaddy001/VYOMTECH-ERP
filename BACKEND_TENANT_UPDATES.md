# Backend Updates Required for Multi-Tenant User Management

## Overview
The frontend now fully supports multi-tenant user registration and tenant switching. This document outlines the backend endpoints and changes needed to complete the implementation.

## Required Backend Endpoints

### 1. Tenant Switch Endpoint
**Endpoint:** `POST /api/v1/tenants/{id}/switch`

Purpose: Switch user's active tenant

Request:
```json
Header: Authorization: Bearer <jwt_token>
```

Response:
```json
{
  "success": true,
  "message": "Tenant switched successfully",
  "tenant": {
    "id": "tenant-123",
    "name": "Acme Corp",
    "domain": "acme.callcenter.com",
    "max_users": 100,
    "max_concurrent_calls": 50,
    "ai_budget_monthly": 1000
  }
}
```

Implementation Notes:
- Validate that user is member of the tenant
- Update user's current_tenant_id in database
- Return new JWT token with updated tenant_id (optional but recommended)

### 2. Add Tenant Member Endpoint
**Endpoint:** `POST /api/v1/tenants/{id}/members`

Purpose: Add a user to a tenant (invite/add team member)

Request:
```json
{
  "email": "user@example.com",
  "role": "admin|member|viewer"
}
Header: Authorization: Bearer <jwt_token>
```

Response:
```json
{
  "success": true,
  "message": "User added to tenant",
  "member": {
    "email": "user@example.com",
    "role": "member",
    "joined_at": "2024-01-15T10:30:00Z"
  }
}
```

Implementation Notes:
- Check request user is tenant admin
- Create tenant_member record
- Optionally send invitation email
- Handle existing members gracefully

### 3. Remove Tenant Member Endpoint
**Endpoint:** `DELETE /api/v1/tenants/{id}/members/{email}`

Purpose: Remove a user from a tenant

Request:
```
Header: Authorization: Bearer <jwt_token>
```

Response:
```json
{
  "success": true,
  "message": "User removed from tenant"
}
```

Implementation Notes:
- Check request user is tenant admin
- Prevent removing the last admin
- Delete tenant_member record
- Clean up user sessions if needed

### 4. Enhanced List Tenants Endpoint
**Endpoint:** `GET /api/v1/tenants`

Current implementation should return:
```json
[
  {
    "id": "tenant-123",
    "name": "Acme Corp",
    "domain": "acme.callcenter.com",
    "max_users": 100,
    "max_concurrent_calls": 50,
    "ai_budget_monthly": 1000,
    "role": "admin",  // User's role in this tenant
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### 5. Enhanced Create Tenant Endpoint
**Endpoint:** `POST /api/v1/tenants`

Current implementation should handle:

Request:
```json
{
  "name": "New Company",
  "domain": "newcompany.callcenter.com"
}
Header: Authorization: Bearer <jwt_token>
```

Response:
```json
{
  "id": "tenant-123",
  "name": "New Company",
  "domain": "newcompany.callcenter.com",
  "max_users": 100,
  "max_concurrent_calls": 50,
  "ai_budget_monthly": 1000,
  "created_at": "2024-01-15T10:30:00Z"
}
```

Implementation Notes:
- Auto-generate tenant ID (UUID or formatted string)
- Make request user admin of new tenant
- Create default tenant_member record with admin role
- Return tenant details immediately

## Database Schema Updates

### User Table Addition
```sql
ALTER TABLE users ADD COLUMN current_tenant_id VARCHAR(36) AFTER tenant_id;
ALTER TABLE users ADD INDEX idx_current_tenant (current_tenant_id);
```

### Tenant Members Table
```sql
CREATE TABLE tenant_members (
  id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
  tenant_id VARCHAR(36) NOT NULL,
  user_id VARCHAR(36) NOT NULL,
  email VARCHAR(255) NOT NULL,
  role ENUM('admin', 'member', 'viewer') DEFAULT 'member',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE KEY unique_tenant_member (tenant_id, user_id),
  INDEX idx_email (email),
  INDEX idx_tenant (tenant_id)
);
```

## Go Code Structure

### Models Update
```go
// Update User model
type User struct {
  ID              string `db:"id"`
  Email           string `db:"email"`
  Name            string `db:"name"`
  PasswordHash    string `db:"password_hash"`
  Role            string `db:"role"`
  TenantID        string `db:"tenant_id"`
  CurrentTenantID string `db:"current_tenant_id"` // New field
  CreatedAt       time.Time `db:"created_at"`
  UpdatedAt       time.Time `db:"updated_at"`
}

// New TenantMember model
type TenantMember struct {
  ID        string `db:"id"`
  TenantID  string `db:"tenant_id"`
  UserID    string `db:"user_id"`
  Email     string `db:"email"`
  Role      string `db:"role"`
  CreatedAt time.Time `db:"created_at"`
  UpdatedAt time.Time `db:"updated_at"`
}
```

### Service Updates
```go
// TenantService additions
func (ts *TenantService) SwitchUserTenant(userID, tenantID string) error
func (ts *TenantService) AddTenantMember(tenantID, email, role string) error
func (ts *TenantService) RemoveTenantMember(tenantID, userID string) error
func (ts *TenantService) GetTenantMembers(tenantID string) ([]TenantMember, error)
func (ts *TenantService) UserIsTenantAdmin(userID, tenantID string) (bool, error)
```

### Handler Updates
```go
// TenantHandler additions
func (th *TenantHandler) SwitchTenant(w http.ResponseWriter, r *http.Request)
func (th *TenantHandler) AddMember(w http.ResponseWriter, r *http.Request)
func (th *TenantHandler) RemoveMember(w http.ResponseWriter, r *http.Request)
func (th *TenantHandler) ListMembers(w http.ResponseWriter, r *http.Request)
```

### Router Updates
```go
// In pkg/router/router.go
// Add new routes
router.HandleFunc("POST /api/v1/tenants/{tenantID}/switch", 
  tenantHandler.SwitchTenant).Methods("POST", "OPTIONS")
router.HandleFunc("POST /api/v1/tenants/{tenantID}/members", 
  tenantHandler.AddMember).Methods("POST", "OPTIONS")
router.HandleFunc("DELETE /api/v1/tenants/{tenantID}/members/{email}", 
  tenantHandler.RemoveMember).Methods("DELETE", "OPTIONS")
router.HandleFunc("GET /api/v1/tenants/{tenantID}/members", 
  tenantHandler.ListMembers).Methods("GET", "OPTIONS")
```

### Middleware Update
Update auth middleware to handle tenant switching:
```go
// Extract tenant_id from JWT and set in context
// Validate user is member of that tenant
// All subsequent operations use that tenant context
```

## Frontend Integration Points

### API Service (`frontend/services/api.ts`)
Already prepared with endpoints:
- `getTenantList()` - Fetches user's tenants
- `switchTenant(tenantId)` - Calls switch endpoint
- `createTenant(name, domain)` - Creates new tenant
- `addTenantMember(tenantId, email, role)` - Adds member
- `removeTenantMember(tenantId, email)` - Removes member

### Components Using Endpoints
- `RegisterForm.tsx` - Creates tenant during registration
- `TenantSwitcher.tsx` - Switches active tenant
- `TenantsPage.tsx` - Lists and manages tenants

## Implementation Checklist

### Database
- [ ] Add `current_tenant_id` column to users table
- [ ] Create `tenant_members` table
- [ ] Add indexes
- [ ] Create migration files

### Go Models
- [ ] Update User struct with current_tenant_id
- [ ] Create TenantMember struct
- [ ] Update serialization/deserialization

### TenantService
- [ ] Implement SwitchUserTenant()
- [ ] Implement AddTenantMember()
- [ ] Implement RemoveTenantMember()
- [ ] Implement GetTenantMembers()
- [ ] Implement UserIsTenantAdmin()
- [ ] Add validation for tenant operations
- [ ] Handle authorization checks

### TenantHandler
- [ ] Implement SwitchTenant handler
- [ ] Implement AddMember handler
- [ ] Implement RemoveMember handler
- [ ] Implement ListMembers handler
- [ ] Add proper error handling
- [ ] Add validation

### Router
- [ ] Add new routes for switch, members endpoints
- [ ] Ensure CORS headers on OPTIONS

### Middleware
- [ ] Update auth middleware to validate tenant membership
- [ ] Ensure tenant_id in JWT matches user's current tenant
- [ ] Add tenant context to request

### Testing
- [ ] Test user can switch tenants
- [ ] Test user can only switch to owned tenants
- [ ] Test add member to tenant
- [ ] Test remove member from tenant
- [ ] Test authorization checks
- [ ] Test invalid tenant IDs
- [ ] Test expired tokens

## Security Considerations

1. **Tenant Isolation**: Ensure users cannot access other tenants' data
2. **Authorization**: Validate user permissions before operations
3. **Role-Based**: Check user role (admin vs member) for operations
4. **Token Validation**: Verify JWT includes valid tenant_id
5. **Rate Limiting**: Consider rate limiting for member operations

## Migration Path

1. Deploy database schema changes
2. Update Go models and services
3. Add new handlers and routes
4. Test with postman/curl
5. Frontend components already support it
6. Test full flow: register → switch → manage

## Troubleshooting

### User Cannot Switch Tenants
- Check tenant_member record exists
- Verify JWT includes tenant_id
- Check current_tenant_id is being updated

### Members Cannot Be Added
- Verify auth token has admin role
- Check email format validation
- Ensure role is valid enum

### Registration with Join Tenant Code Fails
- Verify tenant code format
- Check tenant exists with that code
- Ensure user not already member

## Support

For questions about frontend implementation, refer to:
- `MULTI_TENANT_USER_GUIDE.md` - User workflows
- `frontend/contexts/TenantManagementContext.tsx` - Frontend logic
- `frontend/app/dashboard/tenants/page.tsx` - Tenant management UI
