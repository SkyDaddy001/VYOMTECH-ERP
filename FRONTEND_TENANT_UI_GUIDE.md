# Multi-Tenant Frontend UI Guide

## Visual Overview

### 1. Registration Page (`/auth/register`)

```
┌─────────────────────────────────────────────┐
│         Create Account                       │
├─────────────────────────────────────────────┤
│ User Information                             │
│ ─────────────────────────────────────────── │
│ Full Name: [________________________]        │
│ Email:     [________________________]        │
│ Password:  [________________________]        │
│ Confirm:   [________________________]        │
│                                             │
│ Tenant Selection                            │
│ ─────────────────────────────────────────── │
│ ○ Create New Tenant                         │
│   Tenant Name: [____________________]      │
│   Domain:      [____________________]      │
│                                             │
│ ○ Join Existing Tenant                      │
│   Invite Code: [____________________]      │
│                                             │
│ [        Sign Up        ]                   │
│                                             │
│ Already have account? Login here            │
└─────────────────────────────────────────────┘
```

### 2. Dashboard Sidebar with Tenant Switcher

```
┌──────────────────────────┐
│ ≡  Menu                  │
├──────────────────────────┤
│ Current Tenant           │
│ ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄   │
│ Acme Corp         0/100  │
│ ▼ Switch Tenant          │
│   (3 available)          │
│                          │
│ ┌────────────────────┐   │
│ │ Acme Corp     ✓    │   │
│ │ ────────────────── │   │
│ │ Tech Startup       │   │
│ │ tech.domain.com    │   │
│ │ ────────────────── │   │
│ │ Consulting Inc     │   │
│ │ consulting.com     │   │
│ └────────────────────┘   │
│                          │
├──────────────────────────┤
│ 📊 Dashboard             │
│ 👥 Agents                │
│ 📞 Calls                 │
│ 📋 Campaigns             │
│ 🎯 Leads                 │
│ 📈 Reports               │
│ 🏢 Tenants               │
│                          │
├──────────────────────────┤
│ ⚙️  Settings              │
│ 🚪 Logout                │
└──────────────────────────┘
```

### 3. Tenant Info Card (Top of Dashboard)

```
┌──────────────────────────────────────────────┐
│ Current Tenant Information                   │
├──────────────────────────────────────────────┤
│ Acme Corp                  ✓ Active         │
│ acme.callcenter.com                          │
│                                              │
│ Users:              0 / 100                  │
│ ████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  │
│                                              │
│ Concurrent Calls:   12 / 50                  │
│ ████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  │
│                                              │
│ AI Budget Usage:    $250 / $1000             │
│ ████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  │
└──────────────────────────────────────────────┘
```

### 4. Tenants Management Page (`/dashboard/tenants`)

```
┌────────────────────────────────────────────────────────┐
│ My Tenants                                 [+ Create]  │
│ Manage all tenants you are a member of                 │
├────────────────────────────────────────────────────────┤
│
│ ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ │ Acme Corp    │  │ Tech Startup │  │ Consulting   │
│ │              │  │              │  │              │
│ │ Domain:      │  │ Domain:      │  │ Domain:      │
│ │ acme.cc.com  │  │ tech.cc.com  │  │ cons.cc.com  │
│ │              │  │              │  │              │
│ │ Max Users:   │  │ Max Users:   │  │ Max Users:   │
│ │ 100          │  │ 50           │  │ 200          │
│ │              │  │              │  │              │
│ │ Max Calls:   │  │ Max Calls:   │  │ Max Calls:   │
│ │ 50           │  │ 25           │  │ 100          │
│ │              │  │              │  │              │
│ │ Budget:      │  │ Budget:      │  │ Budget:      │
│ │ $1000        │  │ $500         │  │ $2000        │
│ │              │  │              │  │              │
│ │ Admin        │  │ Member       │  │ Admin        │
│ │              │  │              │  │              │
│ │ [Switch] [Manage]              │  │              │
│ │                                │  │              │
│ └──────────────┘  └──────────────┘  └──────────────┘
│
└────────────────────────────────────────────────────────┘
```

### 5. Create Tenant Modal

```
┌──────────────────────────────────────┐
│ Create New Tenant              × │
├──────────────────────────────────────┤
│                                      │
│ Tenant Name *                        │
│ [______________________________]     │
│                                      │
│ Domain (Optional)                    │
│ [______________________________]     │
│ Used for custom domain setup         │
│                                      │
│ [  Cancel  ]  [  Create  ]           │
│                                      │
└──────────────────────────────────────┘
```

## Component Hierarchy

```
App (layout.tsx)
├── AuthProvider
│   └── TenantWrapper (TenantProvider)
│       └── TenantManagementProvider
│           └── ToasterProvider
│               ├── /auth/register
│               │   └── RegisterForm (with tenant selection)
│               │
│               ├── /auth/login
│               │   └── LoginForm
│               │
│               └── /dashboard
│                   ├── DashboardLayout
│                   │   ├── TenantSwitcher (sidebar)
│                   │   │   └── Tenant dropdown
│                   │   │
│                   │   └── DashboardContent
│                   │       ├── TenantInfo (card)
│                   │       └── Main content
│                   │
│                   └── /dashboard/tenants
│                       └── TenantsPage (grid)
```

## Data Flow Diagram

```
┌─────────────────────────────────────────────────┐
│           User Registration Flow               │
├─────────────────────────────────────────────────┤
│                                                  │
│ [Register Page]                                 │
│       ↓                                          │
│ [RegisterForm] → Select Tenant Mode            │
│       ↓                                          │
│ If Create: POST /api/v1/tenants               │
│ If Join:   Use invite code                    │
│       ↓                                          │
│ [Auth Provider] → register()                   │
│       ↓                                          │
│ Backend: Create user + tenant_member           │
│       ↓                                          │
│ Redirect to /auth/login                        │
│       ↓                                          │
│ [Login] → JWT token + tenant_id                │
│       ↓                                          │
│ Store: localStorage + context                  │
│       ↓                                          │
│ [Dashboard] → Display current tenant           │
│                                                  │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│         Tenant Switching Flow                  │
├─────────────────────────────────────────────────┤
│                                                  │
│ [TenantSwitcher] → Click "Switch Tenant"      │
│       ↓                                          │
│ Show dropdown list of available tenants        │
│       ↓                                          │
│ Select tenant from list                        │
│       ↓                                          │
│ POST /api/v1/tenants/{id}/switch              │
│       ↓                                          │
│ Backend: Update current_tenant_id              │
│       ↓                                          │
│ Frontend: useTenantManagement.switchTenant()  │
│       ↓                                          │
│ Update localStorage: current_tenant_id         │
│       ↓                                          │
│ router.refresh()                               │
│       ↓                                          │
│ Dashboard refreshes with new tenant context    │
│       ↓                                          │
│ All API calls use new tenant context           │
│                                                  │
└─────────────────────────────────────────────────┘
```

## State Management

```
TenantContext (from API)
├── tenant: Current tenant details
├── tenants: List of tenants (legacy)
├── loading: Loading state
└── error: Error messages

TenantManagementContext
├── userTenants: Array of user's tenants
├── currentTenantId: Active tenant ID
├── switchTenant(tenantId): Switch active tenant
├── createTenant(name, domain): Create new tenant
├── addTenantMember(tenantId, email, role): Add user
├── removeTenantMember(tenantId, email): Remove user
├── loading: Operation loading state
└── error: Operation error messages

localStorage
├── auth_token: JWT with tenant_id
├── user: User details
└── current_tenant_id: Active tenant ID
```

## Feature Checklist

✅ Multi-tenant registration (create/join)
✅ Tenant switching from sidebar
✅ Tenant info display
✅ Tenants management page
✅ Create new tenant modal
✅ Tenant context provider
✅ Tenant management hooks
✅ API integration ready
✅ Error handling with toast
✅ Loading states
✅ TypeScript type safety

## User Interaction Flows

### First-Time User
1. Register → Create New Tenant
2. Redirected to login
3. Login → Dashboard
4. Sees new tenant as active
5. Can switch to other tenants on /dashboard/tenants

### Invited User
1. Register → Join Existing Tenant (with code)
2. Redirected to login
3. Login → Dashboard (on joined tenant)
4. Can switch between tenants if member of multiple

### Admin User
1. Can create new tenants
2. Can invite other users
3. Can remove team members
4. Can manage tenant settings

### Returning User
1. Login → Dashboard
2. Last active tenant automatically loaded
3. Can switch at any time
4. Can create/join new tenants
