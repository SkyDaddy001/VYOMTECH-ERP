'use client'

import React, { createContext, useContext, useState, useCallback } from 'react'
import { Tenant } from '@/contexts/TenantContext'

interface UserTenant extends Tenant {
  role?: string // User's role in this tenant
}

interface TenantManagementContextType {
  userTenants: UserTenant[]
  currentTenantId: string | null
  switchTenant: (tenantId: string) => void
  createTenant: (name: string, domain?: string) => Promise<Tenant>
  addTenantMember: (tenantId: string, email: string, role: string) => Promise<void>
  removeTenantMember: (tenantId: string, email: string) => Promise<void>
  loading: boolean
  error: string | null
}

const TenantManagementContext = createContext<TenantManagementContextType | undefined>(undefined)

export function TenantManagementProvider({ children }: { children: React.ReactNode }) {
  const [userTenants, setUserTenants] = useState<UserTenant[]>([])
  const [currentTenantId, setCurrentTenantId] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Load user's tenants from backend
  const loadUserTenants = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const token = localStorage.getItem('auth_token')
      if (!token) return

      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/tenants`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (response.ok) {
        const data = await response.json()
        setUserTenants(Array.isArray(data) ? data : [])
        if (data.length > 0 && !currentTenantId) {
          setCurrentTenantId(data[0].id)
          localStorage.setItem('current_tenant_id', data[0].id)
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load tenants')
    } finally {
      setLoading(false)
    }
  }, [currentTenantId])

  // Switch to a different tenant
  const switchTenant = useCallback((tenantId: string) => {
    const tenant = userTenants.find(t => t.id === tenantId)
    if (tenant) {
      setCurrentTenantId(tenantId)
      localStorage.setItem('current_tenant_id', tenantId)
    } else {
      setError('Tenant not found')
    }
  }, [userTenants])

  // Create a new tenant
  const createTenant = useCallback(async (name: string, domain?: string) => {
    try {
      setLoading(true)
      setError(null)
      const token = localStorage.getItem('auth_token')
      if (!token) throw new Error('Not authenticated')

      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/tenants`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ name, domain }),
      })

      if (!response.ok) {
        throw new Error('Failed to create tenant')
      }

      const tenant = await response.json()
      setUserTenants(prev => [...prev, tenant])
      return tenant
    } finally {
      setLoading(false)
    }
  }, [])

  // Add a tenant member
  const addTenantMember = useCallback(async (tenantId: string, email: string, role: string) => {
    try {
      setLoading(true)
      setError(null)
      const token = localStorage.getItem('auth_token')
      if (!token) throw new Error('Not authenticated')

      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/tenants/${tenantId}/members`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ email, role }),
        }
      )

      if (!response.ok) {
        throw new Error('Failed to add member')
      }
    } finally {
      setLoading(false)
    }
  }, [])

  // Remove a tenant member
  const removeTenantMember = useCallback(async (tenantId: string, email: string) => {
    try {
      setLoading(true)
      setError(null)
      const token = localStorage.getItem('auth_token')
      if (!token) throw new Error('Not authenticated')

      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/tenants/${tenantId}/members/${email}`,
        {
          method: 'DELETE',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )

      if (!response.ok) {
        throw new Error('Failed to remove member')
      }
    } finally {
      setLoading(false)
    }
  }, [])

  return (
    <TenantManagementContext.Provider
      value={{
        userTenants,
        currentTenantId,
        switchTenant,
        createTenant,
        addTenantMember,
        removeTenantMember,
        loading,
        error,
      }}
    >
      {children}
    </TenantManagementContext.Provider>
  )
}

export function useTenantManagement() {
  const context = useContext(TenantManagementContext)
  if (!context) {
    throw new Error('useTenantManagement must be used within TenantManagementProvider')
  }
  return context
}
