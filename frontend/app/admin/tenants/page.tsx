'use client'

import { useEffect, useState } from 'react'
import { adminTenantService } from '@/services/api'

interface Tenant {
  id: string
  name: string
  domain: string
  status: string
  max_users: number
  max_concurrent_calls: number
  ai_budget_monthly: number
  created_at: string
  updated_at: string
}

interface FormData {
  name: string
  domain: string
  status: string
  max_users: number
  max_concurrent_calls: number
  ai_budget_monthly: number
}

export default function TenantsPage() {
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [formData, setFormData] = useState<FormData>({
    name: '',
    domain: '',
    status: 'active',
    max_users: 100,
    max_concurrent_calls: 50,
    ai_budget_monthly: 10000,
  })
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    fetchTenants()
  }, [])

  const fetchTenants = async () => {
    try {
      setLoading(true)
      const response = await adminTenantService.listTenants() as any
      const tenantsData = Array.isArray(response) ? response : response?.data || []
      setTenants(tenantsData)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch tenants')
      setTenants([])
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async () => {
    try {
      setSaving(true)
      await adminTenantService.createTenant(formData)
      setShowForm(false)
      setFormData({
        name: '',
        domain: '',
        status: 'active',
        max_users: 100,
        max_concurrent_calls: 50,
        ai_budget_monthly: 10000,
      })
      await fetchTenants()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create tenant')
    } finally {
      setSaving(false)
    }
  }

  const handleUpdate = async () => {
    if (!editingId) return
    try {
      setSaving(true)
      await adminTenantService.updateTenant(editingId, formData)
      setEditingId(null)
      setShowForm(false)
      setFormData({
        name: '',
        domain: '',
        status: 'active',
        max_users: 100,
        max_concurrent_calls: 50,
        ai_budget_monthly: 10000,
      })
      await fetchTenants()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update tenant')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this tenant?')) return
    try {
      setSaving(true)
      await adminTenantService.deleteTenant(id)
      await fetchTenants()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete tenant')
    } finally {
      setSaving(false)
    }
  }

  const handleEdit = (tenant: Tenant) => {
    setEditingId(tenant.id)
    setFormData({
      name: tenant.name,
      domain: tenant.domain,
      status: tenant.status,
      max_users: tenant.max_users,
      max_concurrent_calls: tenant.max_concurrent_calls,
      ai_budget_monthly: tenant.ai_budget_monthly,
    })
    setShowForm(true)
  }

  const handleCancel = () => {
    setShowForm(false)
    setEditingId(null)
    setFormData({
      name: '',
      domain: '',
      status: 'active',
      max_users: 100,
      max_concurrent_calls: 50,
      ai_budget_monthly: 10000,
    })
  }

  // Calculate stats
  const stats = [
    { label: 'Total Tenants', value: tenants.length, icon: '🏢', color: 'bg-gradient-to-br from-blue-500 to-blue-600' },
    { label: 'Active', value: tenants.filter(t => t.status === 'active').length, icon: '✓', color: 'bg-gradient-to-br from-green-500 to-green-600' },
    { label: 'Inactive', value: tenants.filter(t => t.status === 'inactive').length, icon: '⊗', color: 'bg-gradient-to-br from-yellow-500 to-yellow-600' },
    { label: 'Suspended', value: tenants.filter(t => t.status === 'suspended').length, icon: '⛔', color: 'bg-gradient-to-br from-red-500 to-red-600' },
  ]

  return (
    <div className="min-h-screen bg-slate-50">
      {/* Header */}
      <div className="sticky top-0 z-30 bg-white border-b border-slate-200 shadow-sm">
        <div className="px-8 py-6 flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-slate-900">Tenants Management</h1>
            <p className="text-slate-600 text-sm mt-1">Manage all system tenants and their configurations</p>
          </div>
          <button onClick={() => { setEditingId(null); setShowForm(true) }} className="px-6 py-2.5 bg-gradient-to-r from-blue-600 to-blue-700 text-white rounded-lg hover:shadow-lg transition font-semibold text-sm">
            + Create Tenant
          </button>
        </div>
      </div>

      <div className="p-8">
        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {stats.map((stat, i) => (
            <div key={i} className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 hover:shadow-md transition">
              <div className={`w-12 h-12 rounded-lg ${stat.color} text-white flex items-center justify-center text-lg font-bold mb-4`}>
                {stat.icon}
              </div>
              <p className="text-slate-600 text-sm font-medium">{stat.label}</p>
              <p className="text-4xl font-bold text-slate-900 mt-2">{stat.value}</p>
            </div>
          ))}
        </div>

        {/* Error State */}
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-xl p-4 mb-8 flex items-start gap-3">
            <span className="text-red-600 font-bold text-lg">⚠</span>
            <p className="text-red-800 text-sm">{error}</p>
          </div>
        )}

        {/* Form Modal */}
        {showForm && (
          <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 backdrop-blur-sm">
            <div className="bg-white rounded-2xl shadow-2xl p-8 max-w-md w-full mx-4">
              <h2 className="text-2xl font-bold text-slate-900 mb-6">{editingId ? 'Edit Tenant' : 'Create New Tenant'}</h2>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-semibold text-slate-900 mb-2">Tenant Name *</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                  placeholder="e.g., Acme Corp"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-slate-900 mb-2">Domain</label>
                <input
                  type="text"
                  value={formData.domain}
                  onChange={(e) => setFormData({...formData, domain: e.target.value})}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                  placeholder="acme.example.com"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-slate-900 mb-2">Status</label>
                <select
                  value={formData.status}
                  onChange={(e) => setFormData({...formData, status: e.target.value})}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                >
                  <option value="active">Active</option>
                  <option value="inactive">Inactive</option>
                  <option value="suspended">Suspended</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-semibold text-slate-900 mb-2">Max Users</label>
                <input
                  type="number"
                  value={formData.max_users}
                  onChange={(e) => setFormData({...formData, max_users: parseInt(e.target.value)})}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-slate-900 mb-2">Max Concurrent Calls</label>
                <input
                  type="number"
                  value={formData.max_concurrent_calls}
                  onChange={(e) => setFormData({...formData, max_concurrent_calls: parseInt(e.target.value)})}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-slate-900 mb-2">AI Budget (Monthly)</label>
                <input
                  type="number"
                  value={formData.ai_budget_monthly}
                  onChange={(e) => setFormData({...formData, ai_budget_monthly: parseInt(e.target.value)})}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                />
              </div>
            </div>

            <div className="mt-6 flex gap-3">
              <button
                onClick={editingId ? handleUpdate : handleCreate}
                disabled={saving}
                className="flex-1 px-4 py-2.5 bg-gradient-to-r from-blue-600 to-blue-700 text-white rounded-lg hover:shadow-lg disabled:opacity-50 transition font-semibold text-sm"
              >
                {saving ? 'Saving...' : editingId ? 'Update' : 'Create'}
              </button>
              <button
                onClick={handleCancel}
                disabled={saving}
                className="flex-1 px-4 py-2.5 bg-slate-200 text-slate-900 rounded-lg hover:bg-slate-300 disabled:opacity-50 transition font-semibold text-sm"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Loading State */}
      {loading && (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
          <div className="p-6 space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-16 bg-gradient-to-r from-slate-200 to-slate-100 rounded-lg animate-pulse"></div>
            ))}
          </div>
        </div>
      )}

      {/* Tenants Table */}
      {!loading && tenants.length > 0 && (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
          <table className="w-full">
            <thead className="bg-gradient-to-r from-slate-50 to-slate-100 border-b border-slate-200">
              <tr>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Tenant ID</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Tenant Name</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Domain</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Status</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Max Users</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Max Calls</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">AI Budget</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Created</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200">
              {tenants.map((tenant) => (
                <tr key={tenant.id} className="hover:bg-slate-50 transition">
                  <td className="px-6 py-4 text-sm font-mono text-slate-900 font-medium">{tenant.id.slice(0, 8)}</td>
                  <td className="px-6 py-4 text-sm font-semibold text-slate-900">{tenant.name}</td>
                  <td className="px-6 py-4 text-sm text-slate-600">{tenant.domain || '—'}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center px-3 py-1.5 rounded-full text-xs font-semibold ${
                      tenant.status === 'active'
                        ? 'bg-green-100 text-green-800'
                        : tenant.status === 'inactive'
                        ? 'bg-amber-100 text-amber-800'
                        : 'bg-red-100 text-red-800'
                    }`}>
                      {tenant.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-slate-900 font-medium">{tenant.max_users}</td>
                  <td className="px-6 py-4 text-sm text-slate-900 font-medium">{tenant.max_concurrent_calls}</td>
                  <td className="px-6 py-4 text-sm text-slate-900 font-medium">₹{tenant.ai_budget_monthly?.toLocaleString()}</td>
                  <td className="px-6 py-4 text-sm text-slate-600">
                    {new Date(tenant.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 text-sm space-x-2">
                    <button onClick={() => handleEdit(tenant)} className="text-blue-600 hover:text-blue-800 hover:font-semibold transition">Edit</button>
                    <span className="text-slate-300">|</span>
                    <button onClick={() => handleDelete(tenant.id)} className="text-red-600 hover:text-red-800 hover:font-semibold transition">Delete</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Empty State */}
      {!loading && tenants.length === 0 && !error && (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-12 text-center">
          <div className="w-16 h-16 rounded-full bg-slate-100 flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
            </svg>
          </div>
          <p className="text-slate-600 text-lg font-medium mt-2">No tenants yet</p>
          <p className="text-slate-500 text-sm mt-1">Get started by creating your first tenant</p>
          <button onClick={() => { setEditingId(null); setShowForm(true) }} className="mt-6 px-6 py-2.5 bg-gradient-to-r from-blue-600 to-blue-700 text-white rounded-lg hover:shadow-lg transition font-semibold text-sm">
            Create First Tenant
          </button>
        </div>
      )}
      </div>
    </div>
  )
}
