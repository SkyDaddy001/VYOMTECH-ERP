'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useTenant } from '@/contexts/TenantContext'
import { useTenantManagement } from '@/contexts/TenantManagementContext'
import toast from 'react-hot-toast'

export function TenantSwitcher() {
  const router = useRouter()
  const { tenant, tenants } = useTenant()
  const { userTenants, switchTenant } = useTenantManagement()
  const [isOpen, setIsOpen] = useState(false)
  const [isSwitching, setIsSwitching] = useState(false)

  if (!tenant) return null

  const handleTenantSwitch = async (tenantId: string) => {
    if (tenantId === tenant.id) return

    try {
      setIsSwitching(true)
      
      // Call API to switch tenant
      const token = localStorage.getItem('auth_token')
      if (!token) throw new Error('Not authenticated')

      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/tenants/${tenantId}/switch`,
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )

      if (!response.ok) {
        throw new Error('Failed to switch tenant')
      }

      // Update local state
      switchTenant(tenantId)
      setIsOpen(false)
      
      // Refresh page to reload tenant context
      router.refresh()
      toast.success('Tenant switched successfully')
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to switch tenant')
    } finally {
      setIsSwitching(false)
    }
  }

  const displayTenants = userTenants.length > 0 ? userTenants : tenants

  return (
    <div className="px-4 py-3 border-b border-gray-200">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-xs text-gray-500">Current Tenant</p>
          <p className="font-semibold text-sm text-gray-800">{tenant.name}</p>
        </div>
        <div className="text-right">
          <p className="text-xs text-gray-500">Users</p>
          <p className="font-semibold text-sm text-gray-800">
            0 / {tenant.max_users}
          </p>
        </div>
      </div>
      
      {displayTenants.length > 1 && (
        <div className="mt-3 relative">
          <button
            onClick={() => setIsOpen(!isOpen)}
            disabled={isSwitching}
            className="text-xs text-blue-600 hover:text-blue-700 font-medium disabled:opacity-50"
          >
            {isOpen ? '▾' : '▸'} Switch Tenant ({displayTenants.length} available)
          </button>
          
          {isOpen && (
            <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-lg shadow-lg z-50">
              <div className="max-h-60 overflow-y-auto">
                {displayTenants.map((t) => (
                  <button
                    key={t.id}
                    onClick={() => handleTenantSwitch(t.id)}
                    disabled={isSwitching}
                    className={`w-full text-left px-4 py-3 text-sm border-b last:border-b-0 transition ${
                      t.id === tenant.id
                        ? 'bg-blue-50 text-blue-700 font-semibold'
                        : 'text-gray-700 hover:bg-gray-50'
                    } disabled:opacity-50`}
                  >
                    <div className="flex items-center justify-between">
                      <span>{t.name}</span>
                      {t.id === tenant.id && <span className="text-blue-600">✓</span>}
                    </div>
                    {t.domain && (
                      <p className="text-xs text-gray-500 mt-1">{t.domain}</p>
                    )}
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
