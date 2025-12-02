'use client'

import { useState } from 'react'
import { Tenant } from '@/types/tenant'

interface TenantFormProps {
  tenant?: Tenant
  onSubmit: (data: Partial<Tenant>) => Promise<void>
  onCancel: () => void
  loading?: boolean
}

export default function TenantForm({ tenant, onSubmit, onCancel, loading = false }: TenantFormProps) {
  const [formData, setFormData] = useState({
    name: tenant?.name || '',
    domain: tenant?.domain || '',
    max_users: tenant?.max_users || 10,
    max_concurrent_calls: tenant?.max_concurrent_calls || 5,
  })
  const [error, setError] = useState('')

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: name.includes('max') ? parseInt(value) : value,
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!formData.name.trim()) {
      setError('Tenant name is required')
      return
    }

    try {
      await onSubmit(formData)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save tenant')
    }
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-lg max-w-md w-full p-6">
        <h2 className="text-2xl font-bold text-gray-800 mb-4">
          {tenant ? 'Edit Tenant' : 'Create New Tenant'}
        </h2>

        {error && <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-4 text-red-800">{error}</div>}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Tenant Name *
            </label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Enter tenant name"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Domain
            </label>
            <input
              type="text"
              name="domain"
              value={formData.domain}
              onChange={handleChange}
              placeholder="e.g., tenant.example.com"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Max Users
            </label>
            <input
              type="number"
              name="max_users"
              value={formData.max_users}
              onChange={handleChange}
              min="1"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Max Concurrent Calls
            </label>
            <input
              type="number"
              name="max_concurrent_calls"
              value={formData.max_concurrent_calls}
              onChange={handleChange}
              min="1"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            />
          </div>

          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onCancel}
              className="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded-lg transition disabled:opacity-50"
              disabled={loading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition disabled:opacity-50"
              disabled={loading}
            >
              {loading ? 'Saving...' : 'Save'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
