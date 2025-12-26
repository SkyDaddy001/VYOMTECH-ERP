/**
 * User Service
 * Manages user-related operations
 */

'use client';

import { apiClient } from '@/lib/api-client';

export interface User {
  id: string;
  email: string;
  fullName: string;
  tenantId: string;
  roleId: string;
  status: 'active' | 'inactive' | 'suspended';
  emailVerified: boolean;
  lastLogin?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateUserRequest {
  email: string;
  fullName: string;
  roleId: string;
  sendInvitation?: boolean;
}

export interface UpdateUserRequest {
  fullName?: string;
  roleId?: string;
  status?: 'active' | 'inactive' | 'suspended';
}

export interface UserCount {
  tenantId: string;
  totalUsers: number;
  activeUsers: number;
  inactiveUsers: number;
  suspendedUsers: number;
  seats: {
    total: number;
    used: number;
    available: number;
  };
}

export interface UserActivity {
  userId: string;
  email: string;
  lastActivityAt: string;
  activityType: string;
  ipAddress?: string;
}

class UserService {
  /**
   * Get current user
   */
  async getCurrentUser(): Promise<User> {
    return apiClient.get('/users/me');
  }

  /**
   * Get user by ID
   */
  async getUser(userId: string): Promise<User> {
    return apiClient.get(`/users/${userId}`);
  }

  /**
   * List users for current tenant
   */
  async listUsers(filters?: {
    status?: string;
    roleId?: string;
    skip?: number;
    take?: number;
  }): Promise<{
    data: User[];
    total: number;
  }> {
    return apiClient.get('/users', { params: filters });
  }

  /**
   * Create new user
   */
  async createUser(data: CreateUserRequest): Promise<User> {
    return apiClient.post('/users', data);
  }

  /**
   * Update user
   */
  async updateUser(userId: string, data: UpdateUserRequest): Promise<User> {
    return apiClient.put(`/users/${userId}`, data);
  }

  /**
   * Delete user
   */
  async deleteUser(userId: string): Promise<{ message: string }> {
    return apiClient.delete(`/users/${userId}`);
  }

  /**
   * Get user count summary
   */
  async getUserCount(): Promise<UserCount> {
    return apiClient.get('/users/count');
  }

  /**
   * Get user activity list
   */
  async getUserActivity(filters?: { days?: number; limit?: number }): Promise<UserActivity[]> {
    return apiClient.get('/users/activity', { params: filters });
  }

  /**
   * Resend user invitation
   */
  async resendInvitation(userId: string): Promise<{ message: string }> {
    return apiClient.post(`/users/${userId}/resend-invitation`, {});
  }

  /**
   * Deactivate user
   */
  async deactivateUser(userId: string): Promise<{ message: string }> {
    return apiClient.post(`/users/${userId}/deactivate`, {});
  }

  /**
   * Reactivate user
   */
  async reactivateUser(userId: string): Promise<{ message: string }> {
    return apiClient.post(`/users/${userId}/reactivate`, {});
  }

  /**
   * Suspend user
   */
  async suspendUser(userId: string, reason?: string): Promise<{ message: string }> {
    return apiClient.post(`/users/${userId}/suspend`, { reason });
  }

  /**
   * Update user profile
   */
  async updateProfile(data: { fullName?: string; avatar?: string }): Promise<User> {
    return apiClient.put('/users/me', data);
  }

  /**
   * Get users by role
   */
  async getUsersByRole(roleId: string): Promise<User[]> {
    return apiClient.get(`/users/role/${roleId}`);
  }
}

export const userService = new UserService();
export default userService;
