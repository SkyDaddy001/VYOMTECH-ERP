# Multi-Tenant Integration Checklist

Use this checklist to verify all multi-tenant components are properly integrated and functional.

## Database Setup

- [ ] PostgreSQL database created
- [ ] Migration 001_initial_schema.sql executed
- [ ] Tables created:
  - [ ] `tenants` table
  - [ ] `tenant_users` table
  - [ ] `tenant_configs` table
  - [ ] `users` table has `current_tenant_id` field
  - [ ] `users` table has `role` field
- [ ] Indexes created for performance:
  - [ ] Index on `tenant_users(tenant_id, user_id)`
  - [ ] Index on `users(current_tenant_id)`

## Backend (Go) Implementation

### Services
- [ ] TenantService interface defined in `internal/services/tenant.go`
- [ ] TenantService implementation includes:
  - [ ] GetTenant()
  - [ ] ListTenants()
  - [ ] CreateTenant()
  - [ ] UpdateTenant()
  - [ ] DeleteTenant()
  - [ ] AddTenantUser()
  - [ ] RemoveTenantUser()
  - [ ] GetTenantUsers()
  - [ ] GetUserTenants()
  - [ ] SwitchTenant()

### Handlers
- [ ] TenantHandler created in `internal/handlers/tenant.go`
- [ ] Handler methods implemented:
  - [ ] ListTenants (GET /api/v1/tenants)
  - [ ] GetTenantInfo (GET /api/v1/tenant)
  - [ ] GetTenantUserCount (GET /api/v1/tenant/users/count)
  - [ ] GetUserTenants (GET /api/v1/tenants)
  - [ ] SwitchTenant (POST /api/v1/tenants/{id}/switch)
  - [ ] AddTenantMember (POST /api/v1/tenants/{id}/members)
  - [ ] RemoveTenantMember (DELETE /api/v1/tenants/{id}/members/{email})

### Middleware
- [ ] AuthMiddleware validates JWT token
- [ ] AuthMiddleware extracts tenant_id from token
- [ ] TenantMiddleware validates user access to tenant
- [ ] Proper error handling for unauthorized access

### Router Configuration
- [ ] Routes registered in `pkg/router/router.go`
- [ ] Admin tenant routes configured
- [ ] Protected tenant routes configured
- [ ] Multi-tenant routes configured with auth middleware

## Frontend (Next.js) Implementation

### Context Providers
- [ ] TenantContext created in `frontend/contexts/TenantContext.tsx`
  - [ ] Provides `tenant` and `tenants` state
  - [ ] Implements `loadTenant()` function
  - [ ] Handles loading and error states
- [ ] TenantManagementContext created in `frontend/contexts/TenantManagementContext.tsx`
  - [ ] Provides `userTenants` state
  - [ ] Implements `switchTenant()` function
  - [ ] Implements `addTenantMember()` function
  - [ ] Implements `removeTenantMember()` function
  - [ ] Handles loading and error states

### Components
- [ ] TenantSwitcher component created in `frontend/components/dashboard/TenantSwitcher.tsx`
  - [ ] Shows current tenant
  - [ ] Lists user's tenants
  - [ ] Handles tenant switching with API call
  - [ ] Shows user count / max users
  - [ ] Has proper loading states
  - [ ] Shows success/error messages
- [ ] TenantInfo component created in `frontend/components/dashboard/TenantInfo.tsx`
  - [ ] Displays tenant details
  - [ ] Shows user count
  - [ ] Shows creation date

### Services
- [ ] API service methods added in `frontend/services/api.ts`
  - [ ] getTenantInfo()
  - [ ] getTenantUserCount()
  - [ ] listTenants()
  - [ ] getUserTenants()
  - [ ] switchTenant()
  - [ ] addTenantMember()
  - [ ] removeTenantMember()

### Layout & Providers
- [ ] RootLayout properly nests providers:
  - [ ] AuthProvider (outermost)
  - [ ] TenantProvider
  - [ ] TenantManagementProvider (innermost)
- [ ] All children components can access tenant context

### Pages
- [ ] Dashboard page shows tenant information
- [ ] Dashboard includes TenantSwitcher component
- [ ] Authenticated pages check for current tenant
- [ ] Tenant info displays in header/sidebar

## Integration Tests

### Backend
- [ ] Test: Create tenant
- [ ] Test: List tenants
- [ ] Test: Add user to tenant
- [ ] Test: Remove user from tenant
- [ ] Test: Switch tenant
- [ ] Test: Get tenant users
- [ ] Test: Verify tenant isolation (user can't access other tenants)
- [ ] Test: Verify role-based access control

### Frontend
- [ ] Test: Load tenant information on dashboard
- [ ] Test: Display TenantSwitcher with user's tenants
- [ ] Test: Switch between tenants
- [ ] Test: Show success message on switch
- [ ] Test: Handle switching errors gracefully
- [ ] Test: Verify tenant context updates after switch
- [ ] Test: Test with single tenant (no switcher)
- [ ] Test: Test with multiple tenants

### E2E
- [ ] User can register with tenant
- [ ] User can login and see their tenant
- [ ] User can switch to different tenant
- [ ] User can see tenant-specific data
- [ ] Admin can add members to tenant
- [ ] Admin can remove members from tenant
- [ ] Removed user can't access tenant

## Configuration

### Environment Variables
- [ ] `DATABASE_URL` set for PostgreSQL
- [ ] `JWT_SECRET` configured
- [ ] `API_PORT` configured (backend)
- [ ] `NEXT_PUBLIC_API_URL` configured (frontend)
- [ ] Database name includes multi-tenant tables

### Database Secrets
- [ ] Database credentials stored securely
- [ ] Connection string properly formatted
- [ ] SSL/TLS enabled if needed

## Performance Optimization

- [ ] Database indexes created on:
  - [ ] tenant_users(tenant_id, user_id)
  - [ ] users(current_tenant_id)
  - [ ] tenants(domain)
- [ ] Query optimization:
  - [ ] TenantService uses efficient queries
  - [ ] N+1 query problems avoided
  - [ ] Connection pooling configured
- [ ] Frontend caching:
  - [ ] Tenant data cached in context
  - [ ] API responses cached appropriately
  - [ ] LocalStorage used for tenant preference

## Security

- [ ] JWT token includes tenant_id
- [ ] All endpoints validate user belongs to tenant
- [ ] Row-level security implemented (if using RLS)
- [ ] API validates request ownership
- [ ] CORS properly configured for frontend
- [ ] Rate limiting implemented
- [ ] Audit logging for tenant operations
- [ ] Proper error messages (no data leakage)

## Documentation

- [ ] MULTI_TENANT_FEATURES.md created and complete
- [ ] API endpoint documentation updated
- [ ] Component documentation written
- [ ] Setup guide created
- [ ] Troubleshooting guide included
- [ ] Code comments added for complex logic

## Deployment

- [ ] Database migrations applied to production
- [ ] Environment variables set in production
- [ ] Backend compiled and deployed
- [ ] Frontend built and deployed
- [ ] CORS headers configured
- [ ] SSL certificates configured
- [ ] Database backups configured
- [ ] Monitoring alerts set up

## Post-Deployment Verification

- [ ] Health check endpoint responds
- [ ] Tenant list endpoint accessible
- [ ] Can create tenant via API
- [ ] Can switch tenant via UI
- [ ] Can add/remove members
- [ ] Tenant isolation verified
- [ ] Error handling tested
- [ ] Performance acceptable

## Monitoring & Logging

- [ ] Application logs capture tenant operations
- [ ] Error logs include tenant context
- [ ] Metrics tracked per tenant:
  - [ ] API response times
  - [ ] Error rates
  - [ ] User counts
- [ ] Alerts configured for:
  - [ ] Failed tenant operations
  - [ ] Unusual access patterns
  - [ ] Database connection issues

## Sign-Off

- [ ] Developer testing completed: ___________
- [ ] QA testing completed: ___________
- [ ] Deployment approved: ___________
- [ ] Production verified: ___________

---

**Notes:**
- Review this checklist before deploying multi-tenant features
- Check off items as they're completed
- Address any unchecked items before going to production
- Keep this checklist updated as new features are added

