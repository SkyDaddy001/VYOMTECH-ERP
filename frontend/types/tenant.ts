export interface Tenant {
  id: string
  name: string
  domain?: string
  max_users: number
  max_concurrent_calls: number
  created_at: string
  updated_at: string
  deleted_at?: string
  status: 'active' | 'inactive' | 'suspended'
}

export interface TenantStats {
  total_tenants: number
  active_tenants: number
  total_users: number
  active_calls: number
}
