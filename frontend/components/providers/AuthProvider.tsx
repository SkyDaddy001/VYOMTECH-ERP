'use client'

import React, { createContext, useEffect, useState, ReactNode } from 'react'
import { User } from '@/types'
import { authService } from '@/services/api'

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
        const stored = localStorage.getItem('user')
        if (stored) {
          setUser(JSON.parse(stored))
        }
      } catch (err) {
        console.error('Auth check failed:', err)
      } finally {
        setLoading(false)
      }
    }

    checkAuth()
  }, [])

  const login = async (email: string, password: string) => {
    try {
      setLoading(true)
      setError(null)
      const response = await authService.login(email, password)
      setUser(response.user)
      localStorage.setItem('user', JSON.stringify(response.user))
    } catch (err: any) {
      const message = err.response?.data?.message || err.message || 'Login failed'
      setError(message)
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
      const message = err.response?.data?.message || err.message || 'Registration failed'
      setError(message)
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
