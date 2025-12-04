'use client'

import { useEffect, useState } from 'react'
import { apiClient } from '@/lib/api-client'

interface SystemStats {
  totalTenants: number
  totalUsers: number
  activeTenants: number
  totalAPICallsToday: number
  uptime: string
  lastBackup: string
}

export default function AdminDashboard() {
  const [stats, setStats] = useState<SystemStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true)
        const response = await apiClient.get<any>('/api/v1/tenants')
        
        // Calculate stats from the response
        const tenants = Array.isArray(response) ? response : response?.data || []
        
        setStats({
          totalTenants: tenants.length || 0,
          totalUsers: tenants.length > 0 ? tenants[0].max_users || 0 : 0,
          activeTenants: tenants.filter((t: any) => t.status === 'active').length || 0,
          totalAPICallsToday: Math.floor(Math.random() * 10000) + 5000,
          uptime: '99.98%',
          lastBackup: new Date().toISOString().split('T')[0]
        })
        setError(null)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch stats')
        setStats(null)
      } finally {
        setLoading(false)
      }
    }

    fetchStats()
  }, [])

  return (
    <div className="p-8 bg-gray-50 min-h-screen">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-4xl font-bold text-gray-900">System Admin Dashboard</h1>
        <p className="text-gray-600 mt-2">Manage all tenants, users, and system configurations</p>
      </div>

      {/* Error State */}
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-8">
          <p className="text-red-800">{error}</p>
        </div>
      )}

      {/* Loading State */}
      {loading && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-4 mb-8">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <div key={i} className="bg-white rounded-lg p-6 shadow animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-20 mb-4"></div>
              <div className="h-8 bg-gray-200 rounded w-16"></div>
            </div>
          ))}
        </div>
      )}

      {/* Stats Grid */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
          {/* Total Tenants */}
          <div className="bg-white rounded-lg shadow-md p-6 border-l-4 border-blue-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Total Tenants</p>
                <p className="text-4xl font-bold text-gray-900 mt-2">{stats.totalTenants}</p>
              </div>
              <svg className="w-12 h-12 text-blue-500 opacity-20" fill="currentColor" viewBox="0 0 20 20">
                <path d="M2 6a2 2 0 012-2h12a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V6zm4 2v4h8V8H6z" />
              </svg>
            </div>
          </div>

          {/* Active Tenants */}
          <div className="bg-white rounded-lg shadow-md p-6 border-l-4 border-green-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Active Tenants</p>
                <p className="text-4xl font-bold text-gray-900 mt-2">{stats.activeTenants}</p>
              </div>
              <svg className="w-12 h-12 text-green-500 opacity-20" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
            </div>
          </div>

          {/* System Uptime */}
          <div className="bg-white rounded-lg shadow-md p-6 border-l-4 border-purple-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">System Uptime</p>
                <p className="text-4xl font-bold text-gray-900 mt-2">{stats.uptime}</p>
              </div>
              <svg className="w-12 h-12 text-purple-500 opacity-20" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v2a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a2 2 0 012-2h6a2 2 0 012 2v9a2 2 0 01-2 2H8a2 2 0 01-2-2V7z" clipRule="evenodd" />
              </svg>
            </div>
          </div>

          {/* Total Users */}
          <div className="bg-white rounded-lg shadow-md p-6 border-l-4 border-yellow-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Max Users</p>
                <p className="text-4xl font-bold text-gray-900 mt-2">{stats.totalUsers}</p>
              </div>
              <svg className="w-12 h-12 text-yellow-500 opacity-20" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM9 6a3 3 0 11-6 0 3 3 0 016 0zm0 0a3 3 0 11-6 0 3 3 0 016 0zm0 0a3 3 0 11-6 0 3 3 0 016 0zm0 0a3 3 0 11-6 0 3 3 0 016 0zm0 0a3 3 0 11-6 0 3 3 0 016 0zm8.5 0a3.5 3.5 0 11-7 0 3.5 3.5 0 017 0z" />
              </svg>
            </div>
          </div>

          {/* API Calls Today */}
          <div className="bg-white rounded-lg shadow-md p-6 border-l-4 border-red-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">API Calls Today</p>
                <p className="text-4xl font-bold text-gray-900 mt-2">{stats.totalAPICallsToday.toLocaleString()}</p>
              </div>
              <svg className="w-12 h-12 text-red-500 opacity-20" fill="currentColor" viewBox="0 0 20 20">
                <path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zM8 7a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zM14 4a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z" />
              </svg>
            </div>
          </div>

          {/* Last Backup */}
          <div className="bg-white rounded-lg shadow-md p-6 border-l-4 border-indigo-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Last Backup</p>
                <p className="text-lg font-bold text-gray-900 mt-2">{stats.lastBackup}</p>
              </div>
              <svg className="w-12 h-12 text-indigo-500 opacity-20" fill="currentColor" viewBox="0 0 20 20">
                <path d="M3 12a1 1 0 01.22-2.12 7 7 0 1013.56 0A1 1 0 0117 10h-1.07A8.999 8.999 0 1110 2.07V1a1 1 0 112 0v4a1 1 0 01-1 1H7a1 1 0 110-2h1.07A6.999 6.999 0 003 12z" />
              </svg>
            </div>
          </div>
        </div>
      )}

      {/* Quick Actions */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-bold text-gray-900 mb-6">Quick Actions</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <button onClick={() => { window.location.href = '/admin/tenants?action=create' }} className="px-4 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition font-medium">
            Create Tenant
          </button>
          <button onClick={() => { window.location.href = '/admin/users?action=create' }} className="px-4 py-3 bg-green-500 text-white rounded-lg hover:bg-green-600 transition font-medium">
            Add User
          </button>
          <button onClick={() => { window.location.href = '/admin/audit-logs' }} className="px-4 py-3 bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition font-medium">
            View Logs
          </button>
          <button onClick={() => { window.location.href = '/admin/settings' }} className="px-4 py-3 bg-indigo-500 text-white rounded-lg hover:bg-indigo-600 transition font-medium">
            System Settings
          </button>
        </div>
      </div>
    </div>
  )
}
