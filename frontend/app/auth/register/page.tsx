'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import RegisterForm from '@/components/auth/RegisterForm'
import { tenantService } from '@/services/api'
import toast from 'react-hot-toast'

type TenantMode = 'create' | 'join'

interface TenantData {
  id: string
  name: string
  domain?: string
  status: string
  max_users: number
  max_concurrent_calls: number
  ai_budget_monthly: number
  created_at: string
  updated_at: string
}

export default function RegisterPage() {
  const router = useRouter()
  const { register } = useAuth()
  const [loading, setLoading] = useState(false)

  const handleRegister = async (data: {
    email: string
    password: string
    name: string
    tenantMode: TenantMode
    tenantName?: string
    tenantDomain?: string
    tenantCode?: string
  }) => {
    try {
      setLoading(true)
      
      let tenantId = 'default-tenant'
      
      if (data.tenantMode === 'create') {
        // Create new tenant via API using tenantService
        try {
          console.log('Creating tenant:', { name: data.tenantName, domain: data.tenantDomain })
          const tenantData = await tenantService.createTenant(
            data.tenantName || 'Default Tenant',
            data.tenantDomain || ''
          )
          console.log('Tenant created:', tenantData)
          tenantId = (tenantData as TenantData).id
        } catch (err: any) {
          console.error('Tenant creation error:', err)
          const errorMsg = err.response?.data?.message || err.message || 'Failed to create tenant'
          throw new Error(errorMsg)
        }
      } else if (data.tenantMode === 'join') {
        // Join existing tenant using invite code
        tenantId = data.tenantCode || 'default-tenant'
      }
      
      // Register user with the tenant
      await register(data.email, data.password, 'user', tenantId, data.name)
      
      toast.success('Registration successful! Please login.')
      router.push('/auth/login')
    } catch (error) {
      console.error('Registration error:', error)
      toast.error(error instanceof Error ? error.message : 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 px-4">
      <div className="w-full max-w-md">
        <RegisterForm onSubmit={handleRegister} loading={loading} />
      </div>
    </div>
  )
}
