'use client'

import React, { createContext, useEffect, useState, ReactNode } from 'react'
import { User } from '@/types'
import { authService, authEventEmitter } from '@/services/api'

export interface AuthContextType {
  user: User | null
  loading: boolean
  error: string | null
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, role: string, tenant_id: string, name?: string) => Promise<void>
  logout: () => void
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Check if user is already logged in
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const storedUser = localStorage.getItem('user')
        const storedToken = localStorage.getItem('auth_token')

        if (storedUser && storedToken) {
          // Restore user from storage
          const user = JSON.parse(storedUser)
          setUser(user)
          console.log('User restored from localStorage')
          // Note: Token validation with backend happens on first API call via interceptor
          // If token is invalid, the 401 response will trigger logout
        }
      } catch (err) {
        console.error('Auth check failed:', err)
      } finally {
        setLoading(false)
      }
    }

    checkAuth()
  }, [])

  // Listen for auth events (e.g., 401 unauthorized)
  useEffect(() => {
    const handleAuthEvent = (event: string, data?: any) => {
      if (event === 'logout') {
        console.log('Auth event: logout triggered', data)
        // Silently clear user state to trigger redirect in layout
        setUser(null)
        setError(data?.reason === 'unauthorized' ? 'Your session has expired. Please log in again.' : null)
      }
    }

    authEventEmitter.on(handleAuthEvent)
    return () => authEventEmitter.off(handleAuthEvent)
  }, [])

  const login = async (email: string, password: string) => {
    try {
      setLoading(true)
      setError(null)
      const response = await authService.login(email, password)
      setUser(response.user)
      localStorage.setItem('user', JSON.stringify(response.user))
    } catch (err: any) {
      // Handle custom ApiError with userMessage, or fallback to generic message
      const message = err.userMessage || err.response?.data?.message || err.message || 'Login failed'
      setError(message)
      console.error('Login error:', message)
      throw new Error(message)
    } finally {
      setLoading(false)
    }
  }

  const register = async (
    email: string,
    password: string,
    role: string,
    tenant_id: string,
    name?: string
  ) => {
    try {
      setLoading(true)
      setError(null)
      const response = await authService.register(email, password, role, tenant_id, name)
      setUser(response.user)
      localStorage.setItem('user', JSON.stringify(response.user))
    } catch (err: any) {
      const message = err.userMessage || err.response?.data?.message || err.message || 'Registration failed'
      setError(message)
      console.error('Register error:', message)
      throw new Error(message)
    } finally {
      setLoading(false)
    }
  }

  const logout = () => {
    authService.logout()
    setUser(null)
    localStorage.removeItem('user')
    localStorage.removeItem('auth_token')
  }

  return (
    <AuthContext.Provider value={{ user, loading, error, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  )
}
