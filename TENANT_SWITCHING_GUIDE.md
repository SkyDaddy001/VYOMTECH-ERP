# How to Switch Tenants or Register Users Under Different Tenants

Complete guide for managing multi-tenant registration and switching in your system.

---

## 📋 Table of Contents
1. [User Registration Workflows](#user-registration-workflows)
2. [Tenant Switching Guide](#tenant-switching-guide)
3. [API Endpoints](#api-endpoints)
4. [Frontend Components](#frontend-components)
5. [Data Flow](#data-flow)
6. [Code Examples](#code-examples)
7. [Troubleshooting](#troubleshooting)

---

## User Registration Workflows

### ✨ Workflow 1: Register User & Create New Tenant

**Scenario**: First-time user creating their own company/workspace

**Steps**:
1. Navigate to `http://localhost:3000/auth/register`
2. Fill in registration form:
   - Name: Your full name
   - Email: Your email address
   - Password: Secure password
3. Select **"Create New Tenant"** option
4. Enter tenant details:
   - **Tenant Name**: Your company/workspace name (required)
   - **Domain**: Optional domain name (e.g., company.com)
5. Click **"Register"**
6. You'll be redirected to `/auth/login`
7. Log in with your credentials
8. Dashboard loads with your new tenant as active

**Backend Flow**:
```
1. POST /api/v1/tenants (create tenant)
   - Returns: { id, name, domain, max_users, ... }
2. POST /api/v1/auth/register (create user)
   - Params: { email, password, tenant_id, role: 'admin' }
   - Returns: JWT token with tenant_id embedded
3. User becomes tenant admin automatically
```

**What You Get**:
- ✓ New tenant created
- ✓ User account linked to tenant
- ✓ Full admin access to new tenant
- ✓ Ready to invite team members

---

### 👥 Workflow 2: Register User & Join Existing Tenant

**Scenario**: Team member joining an existing company workspace

**Prerequisites**:
- Tenant must already exist
- Must have tenant invitation code

**Steps**:
1. Navigate to `http://localhost:3000/auth/register`
2. Fill in registration form:
   - Name: Your full name
   - Email: Your email address
   - Password: Secure password
3. Select **"Join Existing Tenant"** option
4. Enter tenant details:
   - **Tenant Code**: The invitation code provided by admin
5. Click **"Register"**
6. You'll be redirected to `/auth/login`
7. Log in with your credentials
8. Dashboard loads with the tenant you joined

**Backend Flow**:
```
1. POST /api/v1/auth/register (create user)
   - Params: { email, password, tenant_id/code, role: 'user' }
   - Returns: JWT token with tenant_id embedded
2. User is added to tenant_members table
3. User joins with 'user' role (not admin)
```

**What You Get**:
- ✓ Access to existing tenant
- ✓ User account created
- ✓ Immediate access to tenant data (calls, agents, etc.)
- ✓ Can't modify tenant settings (unless promoted)

---

### 🔄 Workflow 3: Register User Under Specific Tenant (API)

**Scenario**: Programmatic registration using API

**API Endpoint**:
```bash
POST /api/v1/auth/register
Content-Type: application/json
Authorization: Bearer {token}

{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe",
  "tenant_id": "123",  # or tenant_code for join mode
  "role": "user"       # admin, user, agent, etc.
}
```

**Response**:
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "user-123",
    "email": "user@example.com",
    "tenant_id": "123",
    "role": "user"
  }
}
```

---

## Tenant Switching Guide

### 🔁 How to Switch Tenants (UI)

**Prerequisite**: User must be a member of multiple tenants

**Steps**:
1. Log in to dashboard (any tenant)
2. Look in the left sidebar for **"Current Tenant"** section
3. Click **"Switch Tenant"** button
4. A dropdown appears showing all your tenants
5. Click on desired tenant
6. Dashboard updates automatically
7. All subsequent API calls use new tenant context

**Visual Example**:
```
┌─────────────────────────┐
│ Current Tenant          │
│ Acme Corp               │
│                         │
│ Users: 5 / 100          │
│                         │
│ ▾ Switch Tenant (3)     │  ← Click here
└─────────────────────────┘
  
After clicking:
┌─────────────────────────┐
│ ▸ Acme Corp ✓           │  ← Current
│ TechStart Inc           │  ← Available
│ Green Solutions         │  ← Available
└─────────────────────────┘
```

**What Happens Behind The Scenes**:
1. Frontend calls: `POST /api/v1/tenants/{tenantId}/switch`
2. Backend updates: `user.current_tenant_id = {tenantId}`
3. Frontend updates: `TenantManagementContext.currentTenantId`
4. Frontend saves: `localStorage.current_tenant_id = {tenantId}`
5. Dashboard refreshes with new tenant context
6. All API calls now include new tenant in JWT

---

### 🎯 Switch Tenants (API)

**Endpoint**:
```bash
POST /api/v1/tenants/{tenantId}/switch
Authorization: Bearer {token}
```

**Example using curl**:
```bash
curl -X POST http://localhost:8080/api/v1/tenants/2/switch \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Content-Type: application/json"
```

**Response**:
```json
{
  "success": true,
  "message": "Tenant switched successfully",
  "tenant": {
    "id": 2,
    "name": "TechStart Inc",
    "domain": "techstart.com",
    "max_users": 50,
    "current_users": 3
  }
}
```

**Error Responses**:
```json
// User not member of tenant
{
  "error": "User is not a member of this tenant"
}

// Tenant not found
{
  "error": "Tenant not found"
}

// Not authenticated
{
  "error": "Authentication required"
}
```

---

## API Endpoints

### Authentication & Registration

```bash
# Register new user (with tenant creation)
POST /api/v1/auth/register
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secure123",
  "tenant_id": "new-tenant-id"  # or tenant_code for joining
}

# Login
POST /api/v1/auth/login
{
  "email": "john@example.com",
  "password": "secure123"
}
Response: { token: "jwt...", tenant_id: "123" }

# Refresh token
POST /api/v1/auth/refresh
Authorization: Bearer {token}
```

### Tenant Operations

```bash
# Create new tenant
POST /api/v1/tenants
Authorization: Bearer {token}
{
  "name": "My Company",
  "domain": "mycompany.com",
  "max_users": 100
}

# Get current tenant info
GET /api/v1/tenants
Authorization: Bearer {token}

# Switch to tenant
POST /api/v1/tenants/{id}/switch
Authorization: Bearer {token}

# List all user's tenants
GET /api/v1/tenants/list
Authorization: Bearer {token}
Response: [
  { id: 1, name: "Acme", role: "admin" },
  { id: 2, name: "TechStart", role: "user" }
]

# Add member to tenant
POST /api/v1/tenants/{id}/members
Authorization: Bearer {token}
{
  "email": "newuser@example.com",
  "role": "user"
}

# Remove member from tenant
DELETE /api/v1/tenants/{id}/members/{userId}
Authorization: Bearer {token}

# Get tenant members
GET /api/v1/tenants/{id}/members
Authorization: Bearer {token}
```

---

## Frontend Components

### TenantSwitcher Component

**Location**: `frontend/components/dashboard/TenantSwitcher.tsx`

**Features**:
- Shows current active tenant
- Displays user count vs max users
- Dropdown showing all available tenants
- One-click switching
- Loading states and error handling

**Usage**:
```tsx
import { TenantSwitcher } from '@/components/dashboard/TenantSwitcher'

export function Sidebar() {
  return (
    <div className="sidebar">
      <TenantSwitcher />
      {/* other sidebar items */}
    </div>
  )
}
```

**Props**: None (uses context directly)

**Output**:
```
┌────────────────────────┐
│ Current Tenant         │
│ Acme Corp              │
│                        │
│ Users: 3 / 100         │
│                        │
│ ▾ Switch Tenant (3)    │
└────────────────────────┘
```

---

### RegisterForm Component

**Location**: `frontend/components/auth/RegisterForm.tsx`

**Features**:
- User registration form
- Tenant mode selection (Create New / Join Existing)
- Conditional tenant input fields
- Form validation
- Loading states

**Usage**:
```tsx
import { RegisterForm } from '@/components/auth/RegisterForm'

export default function RegisterPage() {
  return (
    <div>
      <h1>Create Account</h1>
      <RegisterForm />
    </div>
  )
}
```

**Form States**:
```
Create New Tenant:
├─ Name (required)
├─ Email (required)
├─ Password (required)
├─ Tenant Name (required)
└─ Domain (optional)

Join Existing Tenant:
├─ Name (required)
├─ Email (required)
├─ Password (required)
└─ Tenant Code (required)
```

---

### TenantManagementContext

**Location**: `frontend/contexts/TenantManagementContext.tsx`

**Provides**:
```tsx
interface TenantManagementContextType {
  userTenants: UserTenant[]           // All user's tenants
  currentTenantId: string | null      // Currently active tenant
  switchTenant: (id: string) => void  // Switch active tenant
  createTenant: (name: string, domain?: string) => Promise<Tenant>
  addTenantMember: (tenantId: string, email: string, role: string) => Promise<void>
  removeTenantMember: (tenantId: string, email: string) => Promise<void>
  loading: boolean
  error: string | null
}
```

**Usage**:
```tsx
import { useTenantManagement } from '@/contexts/TenantManagementContext'

export function MyComponent() {
  const { userTenants, currentTenantId, switchTenant } = useTenantManagement()
  
  return (
    <div>
      <p>Current: {currentTenantId}</p>
      <p>Total Tenants: {userTenants.length}</p>
      <button onClick={() => switchTenant('123')}>
        Switch to Tenant 123
      </button>
    </div>
  )
}
```

---

## Data Flow

### Registration with New Tenant
```
┌─────────────────┐
│ RegisterForm    │
└────────┬────────┘
         │ User fills form & selects "Create New"
         ↓
┌─────────────────────────────────────┐
│ Frontend: handleRegister()           │
└────────┬────────────────────────────┘
         │
         ├─→ POST /api/v1/tenants (Create Tenant)
         │   ├─→ Database: INSERT INTO tenants
         │   └─→ Response: { id: 1, name, domain, ... }
         │
         ├─→ AuthProvider.register(email, pwd, tenantId)
         │   ├─→ POST /api/v1/auth/register
         │   ├─→ Database: INSERT INTO users
         │   ├─→ Database: INSERT INTO tenant_members
         │   └─→ Response: { token: "jwt..." }
         │
         └─→ localStorage.setItem('auth_token', token)
         └─→ router.push('/auth/login')
         
         ↓
┌─────────────────┐
│ LoginPage       │
└────────┬────────┘
         │ User logs in
         ↓
┌─────────────────────────────┐
│ JWT includes tenant_id      │
│ Redirect to Dashboard       │
└────────┬────────────────────┘
         │
         ├─→ TenantContext loads current tenant
         ├─→ TenantManagementContext loads all user tenants
         └─→ Dashboard displays with new tenant active
```

### Tenant Switching
```
┌──────────────────┐
│ TenantSwitcher   │
│ (Click button)   │
└────────┬─────────┘
         │
         ↓
┌─────────────────────────────┐
│ User selects tenant         │
│ (Dropdown opens)            │
└────────┬────────────────────┘
         │
         ↓
┌──────────────────────────────────────────┐
│ handleTenantSwitch(tenantId)             │
│                                          │
│ 1. POST /api/v1/tenants/{id}/switch     │
│    └─→ Backend: UPDATE users             │
│        SET current_tenant_id = {id}     │
│                                          │
│ 2. switchTenant(tenantId)               │
│    └─→ TenantManagementContext updated  │
│                                          │
│ 3. localStorage updated                  │
│                                          │
│ 4. router.refresh()                      │
│    └─→ Dashboard reloads                │
│    └─→ All API calls use new tenant    │
└──────────────────────────────────────────┘
         │
         ↓
┌─────────────────────────────┐
│ Dashboard refreshed         │
│ Shows new tenant data       │
│ TenantSwitcher shows update │
└─────────────────────────────┘
```

### Multi-Tenant API Calls
```
After switching tenant:

┌────────────────────────────────┐
│ Any API Call (GET /api/v1/...)│
└────────┬───────────────────────┘
         │
         ↓
┌────────────────────────────────────┐
│ Add Authorization Header           │
│ Bearer {token with tenant_id}       │
└────────┬───────────────────────────┘
         │
         ↓
┌────────────────────────────────────┐
│ Backend Middleware                 │
│ 1. Extract JWT                     │
│ 2. Verify tenant_id in JWT         │
│ 3. Validate user access            │
│ 4. Query using tenant_id context   │
└────────┬───────────────────────────┘
         │
         ↓
┌────────────────────────────────────┐
│ Database Query                     │
│ WHERE tenant_id = {from JWT}       │
│ AND user_id = {from JWT}           │
└────────┬───────────────────────────┘
         │
         ↓
┌────────────────────────────────────┐
│ Return tenant-specific data        │
└────────────────────────────────────┘
```

---

## Code Examples

### Example 1: Registering User via Frontend

```tsx
// frontend/app/auth/register/page.tsx

async function handleRegister(data: RegisterFormData) {
  try {
    // Step 1: Create tenant if needed
    let tenantId: string
    
    if (data.tenantMode === 'create') {
      const tenantData = await tenantService.createTenant(
        data.tenantName || 'Default Tenant',
        data.tenantDomain || ''
      )
      tenantId = tenantData.id
    } else {
      // Join existing
      tenantId = data.tenantCode
    }
    
    // Step 2: Register user with tenant
    await register(
      data.email,
      data.password,
      'user',
      tenantId,
      data.name
    )
    
    // Step 3: Redirect to login
    toast.success('Registration successful!')
    router.push('/auth/login')
    
  } catch (error) {
    toast.error(error.message)
  }
}
```

### Example 2: Switching Tenants

```tsx
// In a React component

import { useTenantManagement } from '@/contexts/TenantManagementContext'
import { useRouter } from 'next/navigation'

export function TenantSwitcherExample() {
  const router = useRouter()
  const { userTenants, currentTenantId, switchTenant } = useTenantManagement()
  
  const handleSwitch = async (tenantId: string) => {
    try {
      const token = localStorage.getItem('auth_token')
      
      // Call backend to switch
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/tenants/${tenantId}/switch`,
        {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
          }
        }
      )
      
      if (response.ok) {
        // Update frontend state
        switchTenant(tenantId)
        localStorage.setItem('current_tenant_id', tenantId)
        
        // Refresh to reload tenant context
        router.refresh()
      }
    } catch (error) {
      console.error('Switch failed:', error)
    }
  }
  
  return (
    <div>
      <h3>Current Tenant: {currentTenantId}</h3>
      <ul>
        {userTenants.map(tenant => (
          <li key={tenant.id}>
            {tenant.name}
            {tenant.id !== currentTenantId && (
              <button onClick={() => handleSwitch(tenant.id)}>
                Switch to {tenant.name}
              </button>
            )}
          </li>
        ))}
      </ul>
    </div>
  )
}
```

### Example 3: Creating Tenant via API

```bash
# Register with tenant creation

curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@acme.com",
    "password": "MySecure123!",
    "tenant_name": "Acme Corporation",
    "tenant_domain": "acme.com"
  }'

# Response
{
  "success": true,
  "message": "User registered successfully",
  "user": {
    "id": "user-123",
    "email": "john@acme.com",
    "tenant_id": "tenant-456",
    "role": "admin"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Example 4: Accessing TenantManagement Hook

```tsx
// Any component in dashboard

import { useTenantManagement } from '@/contexts/TenantManagementContext'

export function UserProfile() {
  const {
    userTenants,
    currentTenantId,
    switchTenant,
    createTenant,
    addTenantMember,
    removeTenantMember,
    loading,
    error
  } = useTenantManagement()
  
  // Create new tenant
  const handleCreateTenant = async () => {
    try {
      const newTenant = await createTenant('New Company', 'newcompany.com')
      console.log('Created tenant:', newTenant.id)
      switchTenant(newTenant.id)
    } catch (err) {
      console.error('Failed to create tenant:', err)
    }
  }
  
  // Add member to current tenant
  const handleAddMember = async () => {
    if (!currentTenantId) return
    try {
      await addTenantMember(currentTenantId, 'user@example.com', 'user')
    } catch (err) {
      console.error('Failed to add member:', err)
    }
  }
  
  return (
    <div>
      <h2>My Tenants ({userTenants.length})</h2>
      {userTenants.map(tenant => (
        <div key={tenant.id}>
          {tenant.name}
          {tenant.id === currentTenantId && <span> (Active)</span>}
        </div>
      ))}
      
      <button onClick={handleCreateTenant} disabled={loading}>
        Create New Tenant
      </button>
      
      {error && <p className="error">{error}</p>}
    </div>
  )
}
```

---

## Troubleshooting

### Issue: Can't see tenant switch option

**Causes**:
- User only belongs to one tenant
- TenantManagementContext not loaded
- Browser cache issue

**Solutions**:
1. Create another tenant to be a member of
2. Check browser console for errors
3. Clear localStorage: `localStorage.clear()`
4. Refresh page

---

### Issue: "User is not a member of this tenant"

**Causes**:
- Trying to switch to tenant you're not member of
- Tenant doesn't exist
- User deleted from tenant

**Solutions**:
1. Check you're registered with tenant
2. Ask admin to add you to tenant
3. Create new tenant if needed

---

### Issue: Tenant switch works but data doesn't update

**Causes**:
- Router refresh didn't complete
- JWT token not updated
- Cache not cleared

**Solutions**:
```tsx
// Force hard refresh
router.refresh()

// Or manual page refresh
window.location.reload()

// Or clear cache and switch
localStorage.removeItem('current_tenant_id')
router.refresh()
```

---

### Issue: Registration fails with "Tenant not found"

**Causes**:
- Tenant code is invalid
- Tenant doesn't exist
- Typo in tenant code

**Solutions**:
1. Verify tenant code from admin
2. Ask admin to provide correct code
3. Select "Create New Tenant" option instead

---

### Issue: Can't create new tenant

**Causes**:
- Not authenticated
- JWT expired
- API URL configuration wrong

**Solutions**:
1. Check you're logged in
2. Check browser console for errors
3. Verify `NEXT_PUBLIC_API_URL` env var: `http://app:8080` (Docker)
4. Check backend logs: `podman logs callcenter-app`

---

## Summary

| Task | Where | Method |
|------|-------|--------|
| **Register with new tenant** | `/auth/register` | Select "Create New Tenant" |
| **Register & join tenant** | `/auth/register` | Select "Join Existing Tenant" + code |
| **Switch tenants** | Dashboard sidebar | Click "Switch Tenant" dropdown |
| **Create new tenant** | `/dashboard/tenants` | Click "+ Create Tenant" |
| **Add team member** | `/dashboard/tenants` | Click tenant → Add member |
| **Remove member** | `/dashboard/tenants` | Click tenant → Remove member |
| **API registration** | `POST /api/v1/auth/register` | Send with tenant_id |
| **API switch** | `POST /api/v1/tenants/{id}/switch` | Authenticated request |

---

## Related Documentation

- `MULTI_TENANT_USER_GUIDE.md` - Detailed user workflows
- `FRONTEND_TENANT_UI_GUIDE.md` - UI component reference
- `FRONTEND_TENANT_IMPLEMENTATION.md` - Technical implementation
- `MULTI_TENANT_FEATURES.md` - Feature specifications
- `QUICK_REFERENCE_TENANT.md` - Quick code examples
