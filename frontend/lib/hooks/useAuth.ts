/**
 * useAuth Hook
 * Provides authentication state and methods to components
 */

'use client';

import { useEffect, useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import authService, { type AuthResponse } from '@/services/auth-service';

export interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: AuthResponse['user'] | null;
  error: string | null;
}

export function useAuth() {
  const router = useRouter();
  const [state, setState] = useState<AuthState>({
    isAuthenticated: false,
    isLoading: true,
    user: null,
    error: null,
  });

  // Check token on mount
  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = useCallback(async () => {
    try {
      setState((prev) => ({ ...prev, isLoading: true }));
      const result = await authService.verifyToken();
      
      if (result.valid && result.user) {
        const user = result.user;
        setState((prev) => ({
          ...prev,
          isAuthenticated: true,
          user: {
            id: user.id,
            email: user.email,
            tenantId: user.tenantId,
            fullName: user.email?.split('@')[0] || 'User',
            role: 'member',
          },
          error: null,
        }));
      } else {
        setState((prev) => ({
          ...prev,
          isAuthenticated: false,
          user: null,
        }));
      }
    } catch (error) {
      setState((prev) => ({
        ...prev,
        isAuthenticated: false,
        user: null,
      }));
    } finally {
      setState((prev) => ({ ...prev, isLoading: false }));
    }
  }, []);

  const login = useCallback(
    async (email: string, password: string) => {
      try {
        setState((prev) => ({ ...prev, isLoading: true, error: null }));
        const response = await authService.login({ email, password });
        setState((prev) => ({
          ...prev,
          isAuthenticated: true,
          user: response.user,
          error: null,
        }));
        router.push('/dashboard');
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Login failed';
        setState((prev) => ({
          ...prev,
          isAuthenticated: false,
          user: null,
          error: errorMessage,
        }));
        throw error;
      } finally {
        setState((prev) => ({ ...prev, isLoading: false }));
      }
    },
    [router]
  );

  const signup = useCallback(
    async (email: string, password: string, fullName: string, tenantName: string) => {
      try {
        setState((prev) => ({ ...prev, isLoading: true, error: null }));
        const response = await authService.signup({
          email,
          password,
          fullName,
          tenantName,
        });
        setState((prev) => ({
          ...prev,
          isAuthenticated: true,
          user: response.user,
          error: null,
        }));
        router.push('/dashboard');
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Signup failed';
        setState((prev) => ({
          ...prev,
          isAuthenticated: false,
          user: null,
          error: errorMessage,
        }));
        throw error;
      } finally {
        setState((prev) => ({ ...prev, isLoading: false }));
      }
    },
    [router]
  );

  const logout = useCallback(async () => {
    try {
      setState((prev) => ({ ...prev, isLoading: true }));
      await authService.logout();
      setState({
        isAuthenticated: false,
        isLoading: false,
        user: null,
        error: null,
      });
      router.push('/login');
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      setState((prev) => ({ ...prev, isLoading: false }));
    }
  }, [router]);

  const clearError = useCallback(() => {
    setState((prev) => ({ ...prev, error: null }));
  }, []);

  return {
    ...state,
    login,
    signup,
    logout,
    checkAuth,
    clearError,
  };
}
