import { apiClient } from './api'
import { User, TenantMember, Role, Permission } from '@/types/user'

export const userService = {
  // Get all users in tenant
  async getUsers(): Promise<User[]> {
    return apiClient.get<User[]>('/api/v1/users')
  },

  // Get single user
  async getUser(id: number): Promise<User> {
    return apiClient.get<User>(`/api/v1/users/${id}`)
  },

  // Create user
  async createUser(data: Partial<User> & { password: string }): Promise<User> {
    return apiClient.post<User>('/api/v1/users', data)
  },

  // Update user
  async updateUser(id: number, data: Partial<User>): Promise<User> {
    return apiClient.put<User>(`/api/v1/users/${id}`, data)
  },

  // Delete user
  async deleteUser(id: number): Promise<void> {
    return apiClient.post(`/api/v1/users/${id}/delete`, {})
  },

  // Invite user to tenant
  async inviteUser(email: string, role: string): Promise<TenantMember> {
    return apiClient.post<TenantMember>('/api/v1/users/invite', { email, role })
  },

  // Get tenant members
  async getTenantMembers(): Promise<TenantMember[]> {
    return apiClient.get<TenantMember[]>('/api/v1/tenant-members')
  },

  // Update member role
  async updateMemberRole(userId: number, role: string): Promise<TenantMember> {
    return apiClient.put<TenantMember>(`/api/v1/tenant-members/${userId}`, { role })
  },

  // Remove member from tenant
  async removeMember(userId: number): Promise<void> {
    return apiClient.post(`/api/v1/tenant-members/${userId}/remove`, {})
  },
}

export const rbacService = {
  // Get all roles
  async getRoles(): Promise<Role[]> {
    return apiClient.get<Role[]>('/api/v1/roles')
  },

  // Get single role
  async getRole(id: string): Promise<Role> {
    return apiClient.get<Role>(`/api/v1/roles/${id}`)
  },

  // Create role
  async createRole(data: Partial<Role>): Promise<Role> {
    return apiClient.post<Role>('/api/v1/roles', data)
  },

  // Update role
  async updateRole(id: string, data: Partial<Role>): Promise<Role> {
    return apiClient.put<Role>(`/api/v1/roles/${id}`, data)
  },

  // Delete role
  async deleteRole(id: string): Promise<void> {
    return apiClient.post(`/api/v1/roles/${id}/delete`, {})
  },

  // Get all permissions
  async getPermissions(): Promise<Permission[]> {
    return apiClient.get<Permission[]>('/api/v1/permissions')
  },

  // Assign permission to role
  async assignPermission(roleId: string, permissionId: string): Promise<void> {
    return apiClient.post(`/api/v1/roles/${roleId}/permissions`, { permission_id: permissionId })
  },

  // Remove permission from role
  async removePermission(roleId: string, permissionId: string): Promise<void> {
    return apiClient.post(`/api/v1/roles/${roleId}/permissions/${permissionId}/remove`, {})
  },

  // Assign role to user
  async assignRoleToUser(userId: number, roleId: string): Promise<void> {
    return apiClient.post(`/api/v1/users/${userId}/roles`, { role_id: roleId })
  },

  // Remove role from user
  async removeRoleFromUser(userId: number, roleId: string): Promise<void> {
    return apiClient.post(`/api/v1/users/${userId}/roles/${roleId}/remove`, {})
  },

  // Get user roles
  async getUserRoles(userId: number): Promise<Role[]> {
    return apiClient.get<Role[]>(`/api/v1/users/${userId}/roles`)
  },

  // Check permission
  async checkPermission(resource: string, action: string): Promise<boolean> {
    return apiClient.get<boolean>(`/api/v1/permissions/check?resource=${resource}&action=${action}`)
  },
}
