'use client'

import { useState } from 'react'
import { Tenant } from '@/types/tenant'
import { Button } from '@/components/ui/button'
import { StatCard } from '@/components/ui/stat-card'

interface TenantListProps {
  tenants: Tenant[]
  loading: boolean
  onEdit: (tenant: Tenant) => void
  onDelete: (tenant: Tenant) => void
  onSwitch: (tenant: Tenant) => Promise<void>
}

export default function TenantList({
  tenants,
  loading,
  onEdit,
  onDelete,
  onSwitch,
}: TenantListProps) {
  const [selectedTenant, setSelectedTenant] = useState<string | null>(null)

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading tenants...</p>
        </div>
      </div>
    )
  }

  if (tenants.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No tenants yet. Create your first tenant to get started.</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Summary Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard label="Total Tenants" value={tenants.length} icon="🏢" trend="up" trendValue={`${tenants.length}`} />
        <StatCard label="Active Tenants" value={tenants.filter(t => !t.deleted_at).length} icon="✅" trend="up" trendValue="All active" />
      </div>

      {/* Tenants Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {tenants.map((tenant) => (
          <div
            key={tenant.id}
            className={`bg-white rounded-lg shadow-md hover:shadow-lg transition overflow-hidden border-2 cursor-pointer ${
              selectedTenant === tenant.id ? 'border-blue-500' : 'border-gray-200'
            }`}
            onClick={() => setSelectedTenant(tenant.id)}
          >
            <div className="p-6">
              {/* Tenant Name */}
              <h3 className="text-xl font-bold text-gray-800 mb-2">{tenant.name}</h3>

              {/* Tenant Domain */}
              {tenant.domain && (
                <p className="text-sm text-gray-600 mb-4">
                  <span className="font-semibold">Domain:</span> {tenant.domain}
                </p>
              )}

              {/* Stats */}
              <div className="space-y-3 mb-6 bg-gray-50 p-4 rounded-lg">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Max Users:</span>
                  <span className="font-semibold text-gray-800">{tenant.max_users}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Concurrent Calls:</span>
                  <span className="font-semibold text-gray-800">{tenant.max_concurrent_calls}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Status:</span>
                  <span className="font-semibold text-green-600">Active</span>
                </div>
              </div>

              {/* Actions */}
              <div className="space-y-2">
                <button
                  onClick={() => onSwitch(tenant)}
                  className="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition text-sm"
                >
                  Switch to Tenant
                </button>
                <div className="flex gap-2">
                  <button
                    onClick={() => onEdit(tenant)}
                    className="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded-lg transition text-sm"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(tenant)}
                    className="flex-1 bg-red-100 hover:bg-red-200 text-red-800 font-bold py-2 px-4 rounded-lg transition text-sm"
                  >
                    Delete
                  </button>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
