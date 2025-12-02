export interface User {
  id: number
  email: string
  role: 'admin' | 'manager' | 'member' | 'viewer'
  tenant_id: string
  current_tenant_id: string
  created_at: string
  updated_at: string
}

export interface TenantMember {
  id: string
  tenant_id: string
  user_id: number
  email: string
  role: 'admin' | 'member' | 'viewer'
  created_at: string
  updated_at: string
}

export interface Permission {
  id: string
  name: string
  description: string
  resource: string
  action: string
}

export interface Role {
  id: string
  name: string
  description: string
  permissions: Permission[]
  created_at: string
  updated_at: string
}

export interface RolePermission {
  role_id: string
  permission_id: string
  assigned_at: string
}

export interface UserRole {
  user_id: number
  role_id: string
  assigned_at: string
}
