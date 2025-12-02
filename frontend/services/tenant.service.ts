import { apiClient } from './api'
import { Tenant, TenantStats } from '@/types/tenant'

export const tenantService = {
  // Get all tenants
  async getTenants(): Promise<Tenant[]> {
    return apiClient.get<Tenant[]>('/api/v1/tenants')
  },

  // Get single tenant
  async getTenant(id: string): Promise<Tenant> {
    return apiClient.get<Tenant>(`/api/v1/tenants/${id}`)
  },

  // Create tenant
  async createTenant(data: Partial<Tenant>): Promise<Tenant> {
    return apiClient.post<Tenant>('/api/v1/tenants', data)
  },

  // Update tenant
  async updateTenant(id: string, data: Partial<Tenant>): Promise<Tenant> {
    return apiClient.put<Tenant>(`/api/v1/tenants/${id}`, data)
  },

  // Delete tenant
  async deleteTenant(id: string): Promise<void> {
    return apiClient.post(`/api/v1/tenants/${id}/delete`, {})
  },

  // Get tenant stats
  async getTenantStats(): Promise<TenantStats> {
    return apiClient.get<TenantStats>('/api/v1/tenants/stats/overview')
  },

  // Switch tenant
  async switchTenant(id: string): Promise<Tenant> {
    return apiClient.post<Tenant>('/api/v1/tenants/switch', { tenant_id: id })
  },
}
