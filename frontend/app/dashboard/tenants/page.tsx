'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import { useTenantManagement } from '@/contexts/TenantManagementContext'
import DashboardLayout from '@/components/layouts/DashboardLayout'
import toast from 'react-hot-toast'

export default function TenantsPage() {
  const router = useRouter()
  const { user } = useAuth()
  const { userTenants, switchTenant, createTenant, loading } = useTenantManagement()
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    domain: '',
  })
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    if (!user) {
      router.push('/auth/login')
    }
  }, [user, router])

  const handleCreateTenant = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.name.trim()) {
      toast.error('Tenant name is required')
      return
    }

    try {
      setIsSubmitting(true)
      await createTenant(formData.name, formData.domain)
      toast.success('Tenant created successfully!')
      setFormData({ name: '', domain: '' })
      setShowCreateModal(false)
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to create tenant')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleSwitchTenant = (tenantId: string) => {
    switchTenant(tenantId)
    toast.success('Switched to tenant')
    router.push('/dashboard')
  }

  if (loading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">Loading tenants...</p>
          </div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-800">My Tenants</h1>
            <p className="text-gray-600 mt-2">
              Manage all tenants you are a member of
            </p>
          </div>
          <button
            onClick={() => setShowCreateModal(true)}
            className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg transition"
          >
            + Create Tenant
          </button>
        </div>

        {/* Tenants Grid */}
        {userTenants.length === 0 ? (
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
            <p className="text-gray-600 mb-4">
              You are not a member of any tenants yet.
            </p>
            <button
              onClick={() => setShowCreateModal(true)}
              className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg transition"
            >
              Create Your First Tenant
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {userTenants.map((tenant) => (
              <div
                key={tenant.id}
                className="bg-white rounded-lg shadow-md hover:shadow-lg transition overflow-hidden border border-gray-200"
              >
                <div className="p-6">
                  {/* Tenant Name */}
                  <h3 className="text-xl font-bold text-gray-800 mb-2">
                    {tenant.name}
                  </h3>

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
                      <span className="font-semibold text-gray-800">
                        {tenant.max_users}
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">Max Concurrent Calls:</span>
                      <span className="font-semibold text-gray-800">
                        {tenant.max_concurrent_calls}
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">AI Budget:</span>
                      <span className="font-semibold text-gray-800">
                        ${tenant.ai_budget_monthly}
                      </span>
                    </div>
                  </div>

                  {/* Role */}
                  {tenant.role && (
                    <div className="mb-4">
                      <span className="inline-block bg-blue-100 text-blue-800 text-xs font-semibold px-3 py-1 rounded-full">
                        {tenant.role.charAt(0).toUpperCase() + tenant.role.slice(1)}
                      </span>
                    </div>
                  )}

                  {/* Actions */}
                  <div className="space-y-2">
                    <button
                      onClick={() => handleSwitchTenant(tenant.id)}
                      className="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition text-sm"
                    >
                      Switch to This Tenant
                    </button>
                    <button
                      onClick={() => router.push(`/dashboard/tenants/${tenant.id}/settings`)}
                      className="w-full bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded-lg transition text-sm"
                    >
                      Manage
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Create Tenant Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-lg max-w-md w-full p-6">
            <h2 className="text-2xl font-bold text-gray-800 mb-4">
              Create New Tenant
            </h2>

            <form onSubmit={handleCreateTenant} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Tenant Name *
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) =>
                    setFormData((prev) => ({ ...prev, name: e.target.value }))
                  }
                  placeholder="e.g., Acme Corp"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Domain (Optional)
                </label>
                <input
                  type="text"
                  value={formData.domain}
                  onChange={(e) =>
                    setFormData((prev) => ({ ...prev, domain: e.target.value }))
                  }
                  placeholder="e.g., acme.callcenter.com"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowCreateModal(false)}
                  disabled={isSubmitting}
                  className="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded-lg transition disabled:opacity-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition disabled:opacity-50"
                >
                  {isSubmitting ? 'Creating...' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </DashboardLayout>
  )
}
