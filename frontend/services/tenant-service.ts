/**
 * Tenant Service
 * Manages tenant-related operations
 */

'use client';

import { apiClient } from '@/lib/api-client';

export interface Tenant {
  id: string;
  name: string;
  email: string;
  phone?: string;
  logo?: string;
  maxUsers: number;
  activeUsers: number;
  status: 'active' | 'inactive' | 'suspended';
  createdAt: string;
  updatedAt: string;
}

export interface CreateTenantRequest {
  name: string;
  email: string;
  phone?: string;
  maxUsers: number;
}

export interface UpdateTenantRequest {
  name?: string;
  email?: string;
  phone?: string;
  maxUsers?: number;
  logo?: string;
}

export interface TenantSettings {
  tenantId: string;
  logoUrl?: string;
  primaryColor?: string;
  secondaryColor?: string;
  customDomain?: string;
  allowMultipleRoles: boolean;
  enableSSO: boolean;
  ssoProvider?: string;
}

class TenantService {
  /**
   * Get current tenant
   */
  async getCurrentTenant(): Promise<Tenant> {
    return apiClient.get('/tenants/current');
  }

  /**
   * Get tenant by ID
   */
  async getTenant(tenantId: string): Promise<Tenant> {
    return apiClient.get(`/tenants/${tenantId}`);
  }

  /**
   * List all tenants (admin only)
   */
  async listTenants(filters?: { status?: string; skip?: number; take?: number }): Promise<{
    data: Tenant[];
    total: number;
  }> {
    return apiClient.get('/tenants', { params: filters });
  }

  /**
   * Create new tenant
   */
  async createTenant(data: CreateTenantRequest): Promise<Tenant> {
    return apiClient.post('/tenants', data);
  }

  /**
   * Update tenant
   */
  async updateTenant(tenantId: string, data: UpdateTenantRequest): Promise<Tenant> {
    return apiClient.put(`/tenants/${tenantId}`, data);
  }

  /**
   * Delete tenant
   */
  async deleteTenant(tenantId: string): Promise<{ message: string }> {
    return apiClient.delete(`/tenants/${tenantId}`);
  }

  /**
   * Get tenant settings
   */
  async getTenantSettings(tenantId: string): Promise<TenantSettings> {
    return apiClient.get(`/tenants/${tenantId}/settings`);
  }

  /**
   * Update tenant settings
   */
  async updateTenantSettings(tenantId: string, settings: Partial<TenantSettings>): Promise<TenantSettings> {
    return apiClient.put(`/tenants/${tenantId}/settings`, settings);
  }

  /**
   * Get tenant credentials
   */
  async getTenantCredentials(tenantId: string): Promise<{
    google?: { configured: boolean };
    meta?: { configured: boolean };
    aws?: { configured: boolean };
    email?: { configured: boolean };
    payment?: { configured: boolean };
  }> {
    return apiClient.get(`/tenants/${tenantId}/credentials`);
  }

  /**
   * Store credential for tenant
   */
  async storeCredential(
    tenantId: string,
    credentialType: 'google' | 'meta' | 'aws' | 'email' | 'razorpay' | 'billdesk',
    credential: any
  ): Promise<{ message: string; credentialId: string }> {
    return apiClient.post(`/tenants/${tenantId}/credentials/${credentialType}`, credential);
  }

  /**
   * Delete tenant credential
   */
  async deleteCredential(
    tenantId: string,
    credentialType: string
  ): Promise<{ message: string }> {
    return apiClient.delete(`/tenants/${tenantId}/credentials/${credentialType}`);
  }

  /**
   * Rotate credential
   */
  async rotateCredential(
    tenantId: string,
    credentialType: string,
    newCredential: any
  ): Promise<{ message: string }> {
    return apiClient.post(
      `/tenants/${tenantId}/credentials/${credentialType}/rotate`,
      newCredential
    );
  }
}

export const tenantService = new TenantService();
export default tenantService;
