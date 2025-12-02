'use client'

import { useState, useEffect } from 'react'
import TenantList from '@/components/modules/Tenants/TenantList'
import TenantForm from '@/components/modules/Tenants/TenantForm'
import { Tenant } from '@/types/tenant'
import { tenantService } from '@/services/tenant.service'
import toast from 'react-hot-toast'

export default function TenantsPage() {
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [editingTenant, setEditingTenant] = useState<Tenant | undefined>()

  useEffect(() => {
    loadTenants()
  }, [])

  const loadTenants = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await tenantService.getTenants()
      setTenants(data)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to load tenants'
      setError(errorMsg)
      toast.error(errorMsg)
      console.error('Error loading tenants:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (data: Partial<Tenant>) => {
    try {
      await tenantService.createTenant(data)
      toast.success('Tenant created successfully!')
      await loadTenants()
      setShowForm(false)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to create tenant'
      toast.error(errorMsg)
      console.error('Error creating tenant:', err)
      throw err
    }
  }

  const handleEdit = async (tenant: Tenant, data: Partial<Tenant>) => {
    try {
      await tenantService.updateTenant(tenant.id, data)
      toast.success('Tenant updated successfully!')
      await loadTenants()
      setEditingTenant(undefined)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to update tenant'
      toast.error(errorMsg)
      console.error('Error updating tenant:', err)
      throw err
    }
  }

  const handleDelete = async (tenant: Tenant) => {
    try {
      if (confirm('Are you sure you want to delete this tenant?')) {
        await tenantService.deleteTenant(tenant.id)
        toast.success('Tenant deleted successfully!')
        await loadTenants()
      }
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to delete tenant'
      toast.error(errorMsg)
      console.error('Error deleting tenant:', err)
    }
  }

  const handleSwitch = async (tenant: Tenant) => {
    try {
      await tenantService.switchTenant(tenant.id)
      toast.success('Switched tenant successfully!')
      // Optionally redirect or refresh context
      window.location.reload()
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to switch tenant'
      toast.error(errorMsg)
      console.error('Error switching tenant:', err)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Tenants</h1>
          <p className="mt-2 text-gray-600">Manage your organization's tenants</p>
        </div>
        <button
          onClick={() => {
            setEditingTenant(undefined)
            setShowForm(true)
          }}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700"
        >
          + New Tenant
        </button>
      </div>

      {error && (
        <div className="rounded-md bg-red-50 p-4">
          <p className="text-sm font-medium text-red-800">{error}</p>
        </div>
      )}

      <TenantList
        tenants={tenants}
        loading={loading}
        onEdit={(tenant) => {
          setEditingTenant(tenant)
          setShowForm(true)
        }}
        onDelete={handleDelete}
        onSwitch={handleSwitch}
      />

      {showForm && (
        <TenantForm
          tenant={editingTenant}
          onSubmit={editingTenant ? (data) => handleEdit(editingTenant, data) : handleCreate}
          onCancel={() => {
            setShowForm(false)
            setEditingTenant(undefined)
          }}
        />
      )}
    </div>
  )
}
