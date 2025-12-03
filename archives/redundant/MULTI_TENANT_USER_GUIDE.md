# Multi-Tenant User Management Guide

## Overview
This guide explains how users can register, manage, and switch between multiple tenants in the AI Call Center system.

## Features Implemented

### 1. Multi-Tenant Registration
Users can now register with the ability to:
- **Create a new tenant** during registration
- **Join an existing tenant** using an invitation code

#### Registration Flow
1. Go to `/auth/register`
2. Enter user details (Name, Email, Password)
3. Select tenant mode:
   - **Create New Tenant**: Enter tenant name and optional domain
   - **Join Existing Tenant**: Enter tenant invitation code
4. Complete registration

### 2. Tenant Switching
Once logged in, users with multiple tenants can:
- Switch between tenants using the sidebar TenantSwitcher component
- View all available tenants on `/dashboard/tenants` page
- Each tenant switch updates the current context and JWT token

#### Switching Tenants
1. Look for "Current Tenant" in the sidebar
2. Click "Switch Tenant" to see all available tenants
3. Select a tenant to switch to
4. Dashboard updates to show that tenant's data

### 3. Tenant Management Page
Navigate to `/dashboard/tenants` to:
- View all tenants you're a member of
- See tenant details (domain, limits, budget)
- Create new tenants
- Manage tenant settings
- Switch tenants

### 4. Dynamic API Integration
The system now supports:
- `POST /api/v1/tenants` - Create new tenant
- `GET /api/v1/tenants` - List all user's tenants
- `GET /api/v1/tenants/{id}` - Get specific tenant
- `POST /api/v1/tenants/{id}/switch` - Switch current tenant
- `POST /api/v1/tenants/{id}/members` - Add tenant member
- `DELETE /api/v1/tenants/{id}/members/{email}` - Remove tenant member

## Frontend Components

### TenantManagementContext
Location: `frontend/contexts/TenantManagementContext.tsx`

Provides:
- `userTenants`: Array of all tenants user belongs to
- `currentTenantId`: Currently active tenant
- `switchTenant()`: Switch active tenant
- `createTenant()`: Create new tenant
- `addTenantMember()`: Add user to tenant
- `removeTenantMember()`: Remove user from tenant

### Updated Components
- **RegisterForm.tsx**: Now includes tenant mode selection
- **TenantSwitcher.tsx**: Functional tenant switching with API calls
- **TenantsPage**: New page for managing all tenants

## Backend Requirements

The following backend endpoints need to be available:

```go
// Tenant Operations
POST   /api/v1/tenants          - Create tenant
GET    /api/v1/tenants          - List user's tenants
GET    /api/v1/tenants/{id}     - Get tenant
POST   /api/v1/tenants/{id}/switch - Switch current tenant

// Tenant Members
POST   /api/v1/tenants/{id}/members           - Add member
DELETE /api/v1/tenants/{id}/members/{email}   - Remove member
```

## User Workflows

### Workflow 1: Register with New Tenant
```
1. User visits /auth/register
2. Fills in user details
3. Selects "Create New Tenant"
4. Enters tenant name
5. Submits registration
6. Tenant is created and user becomes admin
7. Redirect to login
8. User logs in and goes to dashboard
```

### Workflow 2: Register and Join Existing Tenant
```
1. User visits /auth/register
2. Fills in user details
3. Selects "Join Existing Tenant"
4. Enters tenant invitation code
5. Submits registration
6. User is added to existing tenant
7. Redirect to login
8. User logs in to dashboard of that tenant
```

### Workflow 3: Switch Between Tenants
```
1. User logs in with tenant A active
2. In sidebar, click "Switch Tenant"
3. Select tenant B from dropdown
4. System calls /api/v1/tenants/{tenantB}/switch
5. Dashboard refreshes and shows tenant B data
6. All subsequent API calls use tenant B context
```

### Workflow 4: Create Additional Tenant
```
1. User goes to /dashboard/tenants
2. Clicks "+ Create Tenant"
3. Enters tenant details
4. New tenant created
5. Can switch to it immediately
```

## Data Flow

### Registration with New Tenant
```
RegisterForm (submit)
  ↓
register page (handleRegister)
  ↓
POST /api/v1/tenants (create tenant)
  ↓
register() auth method (with tenantId)
  ↓
JWT token includes tenant_id
  ↓
Redirect to login
```

### Tenant Switching
```
TenantSwitcher (user clicks switch)
  ↓
POST /api/v1/tenants/{id}/switch
  ↓
TenantManagementContext.switchTenant()
  ↓
Update localStorage: current_tenant_id
  ↓
router.refresh() 
  ↓
Dashboard refreshes with new tenant context
```

## Storage

### LocalStorage Keys
- `auth_token`: JWT authentication token
- `current_tenant_id`: Currently active tenant ID

### Context State
- TenantContext: Current tenant from API
- TenantManagementContext: All user tenants list

## Error Handling

All tenant operations include:
- Loading states
- Error messages (toast notifications)
- Validation before submission
- API error handling

## Next Steps

### To Complete Backend Support:
1. Add `GET /api/v1/tenants/{id}/switch` endpoint
2. Add `POST /api/v1/tenants/{id}/members` endpoint
3. Add `DELETE /api/v1/tenants/{id}/members/{email}` endpoint
4. Implement tenant validation in all endpoints
5. Add proper authorization checks

### To Enhance Frontend:
1. Add tenant settings page
2. Add team member management
3. Add invite user by email
4. Add tenant billing/usage page
5. Add tenant activity log
6. Add user role management within tenant

## Testing Checklist

- [ ] Register with new tenant
- [ ] Register and join existing tenant
- [ ] Switch between multiple tenants
- [ ] Create tenant from tenants page
- [ ] Verify API calls include correct tenant context
- [ ] Test error handling for invalid tenants
- [ ] Test loading states
- [ ] Verify localStorage persists tenant selection
