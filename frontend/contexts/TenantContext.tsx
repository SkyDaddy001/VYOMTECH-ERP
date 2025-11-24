'use client'

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'

export interface Tenant {
  id: string
  name: string
  domain?: string
  status: 'active' | 'inactive' | 'suspended'
  max_users: number
  max_concurrent_calls: number
  ai_budget_monthly: number
  created_at: string
  updated_at: string
}

interface TenantContextType {
  tenant: Tenant | null
  tenants: Tenant[]
  loading: boolean
  error: string | null
  setTenant: (tenant: Tenant) => void
  refreshTenant: () => Promise<void>
  refreshTenants: () => Promise<void>
}

const TenantContext = createContext<TenantContextType | undefined>(undefined)

export function TenantProvider({ children }: { children: ReactNode }) {
  const [tenant, setTenant] = useState<Tenant | null>(null)
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Get current tenant info
  const refreshTenant = async () => {
    try {
      setLoading(true)
      setError(null)
      const token = localStorage.getItem('auth_token')
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/tenant`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      if (!response.ok) throw new Error('Failed to fetch tenant')
      const data = await response.json()
      setTenant(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
    } finally {
      setLoading(false)
    }
  }

  // Get all tenants (admin)
  const refreshTenants = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/tenants`)
      if (!response.ok) throw new Error('Failed to fetch tenants')
      const data = await response.json()
      setTenants(Array.isArray(data) ? data : [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
    } finally {
      setLoading(false)
    }
  }

  // Load tenant info on mount
  useEffect(() => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      refreshTenant()
    }
  }, [])

  return (
    <TenantContext.Provider
      value={{
        tenant,
        tenants,
        loading,
        error,
        setTenant,
        refreshTenant,
        refreshTenants,
      }}
    >
      {children}
    </TenantContext.Provider>
  )
}

export function useTenant() {
  const context = useContext(TenantContext)
  if (!context) {
    throw new Error('useTenant must be used within TenantProvider')
  }
  return context
}
