/**
 * Authentication Service
 * Handles login, signup, and token management
 */

'use client';

import { apiClient } from '@/lib/api-client';

export interface LoginRequest {
  email: string;
  password: string;
  tenantId?: string;
}

export interface SignupRequest {
  email: string;
  password: string;
  fullName: string;
  tenantName: string;
}

export interface AuthResponse {
  token: string;
  refreshToken?: string;
  user: {
    id: string;
    email: string;
    fullName: string;
    tenantId: string;
    role: string;
  };
}

export interface VerifyTokenResponse {
  valid: boolean;
  user?: {
    id: string;
    email: string;
    tenantId: string;
  };
}

class AuthService {
  /**
   * Login user
   */
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    try {
      const data = await apiClient.login(credentials.email, credentials.password);
      return {
        token: data.token,
        refreshToken: data.refreshToken,
        user: {
          id: data.user.id,
          email: data.user.email,
          fullName: data.user.fullName,
          tenantId: data.user.tenantId,
          role: data.user.role,
        },
      };
    } catch (error) {
      throw new Error(error instanceof Error ? error.message : 'Login failed');
    }
  }

  /**
   * Signup new user/tenant
   */
  async signup(data: SignupRequest): Promise<AuthResponse> {
    // Implement signup via API call
    const response = await fetch('/api/v1/auth/signup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Signup failed');
    return response.json();
  }

  /**
   * Refresh access token
   */
  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    const response = await fetch('/api/v1/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refreshToken }),
    });
    if (!response.ok) throw new Error('Token refresh failed');
    return response.json();
  }

  /**
   * Verify current token
   */
  async verifyToken(): Promise<VerifyTokenResponse> {
    try {
      const response = await fetch('/api/v1/auth/verify', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('authToken')}` },
      });
      if (!response.ok) return { valid: false };
      const data = await response.json();
      return { valid: true, user: data.data };
    } catch (error) {
      return { valid: false };
    }
  }

  /**
   * Logout user
   */
  async logout(): Promise<void> {
    try {
      await fetch('/api/v1/auth/logout', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('authToken')}` },
      });
    } finally {
      localStorage.removeItem('authToken');
    }
  }

  /**
   * Change password
   */
  async changePassword(
    currentPassword: string,
    newPassword: string
  ): Promise<{ message: string }> {
    const response = await fetch('/api/v1/auth/change-password', {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
      },
      body: JSON.stringify({ currentPassword, newPassword }),
    });
    if (!response.ok) throw new Error('Password change failed');
    return response.json();
  }

  /**
   * Request password reset
   */
  async requestPasswordReset(email: string): Promise<{ message: string }> {
    const response = await fetch('/api/v1/auth/forgot-password', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email }),
    });
    if (!response.ok) throw new Error('Password reset request failed');
    return response.json();
  }

  /**
   * Reset password with token
   */
  async resetPassword(token: string, newPassword: string): Promise<{ message: string }> {
    const response = await fetch('/api/v1/auth/reset-password', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token, newPassword }),
    });
    if (!response.ok) throw new Error('Password reset failed');
    return response.json();
  }
}

export const authService = new AuthService();
export default authService;

