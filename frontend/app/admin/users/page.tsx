'use client'

import { useEffect, useState } from 'react'
import { adminUserService, tenantService } from '@/services/api'

interface User {
  id: number
  email: string
  name: string
  role: string
  tenant_id: string
  created_at: string
  updated_at: string
}

interface Tenant {
  id: string
  name: string
}

interface FormData {
  email: string
  name: string
  password: string
  role: string
  tenant_id: string
}

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [formData, setFormData] = useState<FormData>({
    email: '',
    name: '',
    password: '',
    role: 'partner_admin',
    tenant_id: '',
  })
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      setLoading(true)
      const [usersRes, tenantsRes] = await Promise.all([
        adminUserService.listUsers().catch(() => []) as any,
        tenantService.listTenants().catch(() => []) as any,
      ])
      setUsers(Array.isArray(usersRes) ? usersRes : usersRes?.data || [])
      setTenants(Array.isArray(tenantsRes) ? tenantsRes : tenantsRes?.data || [])
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch data')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async () => {
    if (!formData.email || !formData.password || !formData.tenant_id) {
      setError('Email, password, and tenant are required')
      return
    }
    try {
      setSaving(true)
      await adminUserService.createUser(formData)
      setShowForm(false)
      setFormData({
        email: '',
        name: '',
        password: '',
        role: 'partner_admin',
        tenant_id: '',
      })
      await fetchData()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create user')
    } finally {
      setSaving(false)
    }
  }

  const handleUpdate = async () => {
    if (!editingId) return
    try {
      setSaving(true)
      await adminUserService.updateUser(editingId, {
        name: formData.name,
        role: formData.role,
      })
      setEditingId(null)
      setShowForm(false)
      setFormData({
        email: '',
        name: '',
        password: '',
        role: 'partner_admin',
        tenant_id: '',
      })
      await fetchData()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update user')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: number) => {
    const user = users.find(u => u.id === id)
    if (!confirm(`Delete user ${user?.email}? This cannot be undone.`)) return
    try {
      setSaving(true)
      await adminUserService.deleteUser(id)
      await fetchData()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete user')
    } finally {
      setSaving(false)
    }
  }

  const handleEdit = (user: User) => {
    setEditingId(user.id)
    setFormData({
      email: user.email,
      name: user.name,
      password: '',
      role: user.role,
      tenant_id: user.tenant_id,
    })
    setShowForm(true)
  }

  const handleCancel = () => {
    setShowForm(false)
    setEditingId(null)
    setFormData({
      email: '',
      name: '',
      password: '',
      role: 'partner_admin',
      tenant_id: '',
    })
  }

  const getTenantName = (tenantId: string) => {
    return tenants.find(t => t.id === tenantId)?.name || tenantId
  }

  // Calculate stats
  const stats = [
    { label: 'Total Users', value: users.length, icon: '👥', color: 'bg-gradient-to-br from-blue-500 to-blue-600' },
    { label: 'Admins', value: users.filter(u => u.role === 'admin').length, icon: '👑', color: 'bg-gradient-to-br from-purple-500 to-purple-600' },
    { label: 'Managers', value: users.filter(u => u.role === 'manager').length, icon: '📊', color: 'bg-gradient-to-br from-cyan-500 to-cyan-600' },
    { label: 'Agents', value: users.filter(u => u.role === 'agent').length, icon: '☎️', color: 'bg-gradient-to-br from-teal-500 to-teal-600' },
  ]

  return (
    <div className="min-h-screen bg-slate-50">
      {/* Header */}
      <div className="sticky top-0 z-30 bg-white border-b border-slate-200 shadow-sm">
        <div className="px-8 py-6 flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-slate-900">Users Management</h1>
            <p className="text-slate-600 text-sm mt-1">Manage system users, roles, and permissions</p>
          </div>
          <button onClick={() => { setEditingId(null); setShowForm(true) }} className="px-6 py-2.5 bg-gradient-to-r from-blue-600 to-blue-700 text-white rounded-lg hover:shadow-lg transition font-semibold text-sm">
            + Create User
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
              <h2 className="text-2xl font-bold text-slate-900 mb-6">{editingId ? 'Edit User' : 'Create New User'}</h2>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-semibold text-slate-900 mb-2">Email *</label>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({...formData, email: e.target.value})}
                  disabled={editingId !== null}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white disabled:bg-slate-100 disabled:text-slate-500"
                  placeholder="user@example.com"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-slate-900 mb-2">Name</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                  placeholder="John Doe"
                />
              </div>

              {!editingId && (
                <div>
                  <label className="block text-sm font-semibold text-slate-900 mb-2">Password *</label>
                  <input
                    type="password"
                    value={formData.password}
                    onChange={(e) => setFormData({...formData, password: e.target.value})}
                    className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                    placeholder="••••••••"
                  />
                </div>
              )}

              <div>
                <label className="block text-sm font-semibold text-slate-900 mb-2">Role</label>
                <select
                  value={formData.role}
                  onChange={(e) => setFormData({...formData, role: e.target.value})}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                >
                  <option value="admin">Admin</option>
                  <option value="partner_admin">Partner Admin</option>
                  <option value="manager">Manager</option>
                  <option value="agent">Agent</option>
                  <option value="user">User</option>
                </select>
              </div>

              {!editingId && (
                <div>
                  <label className="block text-sm font-semibold text-slate-900 mb-2">Tenant *</label>
                  <select
                    value={formData.tenant_id}
                    onChange={(e) => setFormData({...formData, tenant_id: e.target.value})}
                    className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
                  >
                    <option value="">Select a tenant</option>
                    {tenants.map(t => (
                      <option key={t.id} value={t.id}>{t.name}</option>
                    ))}
                  </select>
                </div>
              )}
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

      {/* Users Table */}
      {!loading && users.length > 0 && (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
          <table className="w-full">
            <thead className="bg-gradient-to-r from-slate-50 to-slate-100 border-b border-slate-200">
              <tr>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Email</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Name</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Role</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Tenant</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Created</th>
                <th className="px-6 py-4 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200">
              {users.map((user) => (
                <tr key={user.id} className="hover:bg-slate-50 transition">
                  <td className="px-6 py-4 text-sm font-medium text-slate-900">{user.email}</td>
                  <td className="px-6 py-4 text-sm text-slate-600">{user.name || '—'}</td>
                  <td className="px-6 py-4">
                    <span className="inline-flex items-center px-3 py-1.5 bg-blue-100 text-blue-800 rounded-full text-xs font-semibold">
                      {user.role}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-slate-600">{getTenantName(user.tenant_id)}</td>
                  <td className="px-6 py-4 text-sm text-slate-600">
                    {new Date(user.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 text-sm space-x-2">
                    <button onClick={() => handleEdit(user)} className="text-blue-600 hover:text-blue-800 hover:font-semibold transition">Edit</button>
                    <span className="text-slate-300">|</span>
                    <button onClick={() => handleDelete(user.id)} className="text-red-600 hover:text-red-800 hover:font-semibold transition">Delete</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Empty State */}
      {!loading && users.length === 0 && !error && (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-12 text-center">
          <div className="w-16 h-16 rounded-full bg-slate-100 flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
            </svg>
          </div>
          <p className="text-slate-600 text-lg font-medium mt-2">No users yet</p>
          <p className="text-slate-500 text-sm mt-1">Get started by creating your first user</p>
          <button onClick={() => { setEditingId(null); setShowForm(true) }} className="mt-6 px-6 py-2.5 bg-gradient-to-r from-blue-600 to-blue-700 text-white rounded-lg hover:shadow-lg transition font-semibold text-sm">
            Create First User
          </button>
        </div>
      )}
      </div>
    </div>
  )
}
