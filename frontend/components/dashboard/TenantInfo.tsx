'use client'

import { useTenant } from '@/contexts/TenantContext'
import { useEffect, useState } from 'react'

export function TenantInfo() {
  const { tenant, loading } = useTenant()
  const [userCount, setUserCount] = useState(0)

  useEffect(() => {
    const fetchUserCount = async () => {
      if (!tenant) return
      try {
        const token = localStorage.getItem('auth_token')
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/api/v1/tenant/users/count`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        )
        if (response.ok) {
          const data = await response.json()
          setUserCount(data.user_count || 0)
        }
      } catch (error) {
        console.error('Failed to fetch user count:', error)
      }
    }

    fetchUserCount()
  }, [tenant])

  if (loading || !tenant) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
          <div className="h-8 bg-gray-200 rounded w-3/4"></div>
        </div>
      </div>
    )
  }

  const userUsagePercent = (userCount / tenant.max_users) * 100
  const budgetUsagePercent = 40 // Placeholder

  return (
    <div className="space-y-4">
      {/* Tenant Overview */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold text-gray-800 mb-4">Tenant Overview</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Basic Info */}
          <div>
            <p className="text-sm text-gray-600 mb-1">Tenant Name</p>
            <p className="text-lg font-semibold text-gray-900">{tenant.name}</p>
          </div>

          {/* Status */}
          <div>
            <p className="text-sm text-gray-600 mb-1">Status</p>
            <p className={`text-lg font-semibold ${
              tenant.status === 'active' ? 'text-green-600' :
              tenant.status === 'suspended' ? 'text-red-600' :
              'text-gray-600'
            }`}>
              {tenant.status.charAt(0).toUpperCase() + tenant.status.slice(1)}
            </p>
          </div>
        </div>
      </div>

      {/* Resource Usage */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-bold text-gray-800 mb-4">Resource Usage</h3>

        {/* User Limit */}
        <div className="mb-4">
          <div className="flex justify-between mb-2">
            <p className="text-sm font-semibold text-gray-700">Users</p>
            <p className="text-sm text-gray-600">
              {userCount} / {tenant.max_users}
            </p>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div
              className={`h-2 rounded-full transition-all ${
                userUsagePercent > 80 ? 'bg-red-500' :
                userUsagePercent > 50 ? 'bg-yellow-500' :
                'bg-green-500'
              }`}
              style={{ width: `${Math.min(userUsagePercent, 100)}%` }}
            ></div>
          </div>
        </div>

        {/* Concurrent Calls */}
        <div className="mb-4">
          <div className="flex justify-between mb-2">
            <p className="text-sm font-semibold text-gray-700">Max Concurrent Calls</p>
            <p className="text-sm text-gray-600">{tenant.max_concurrent_calls}</p>
          </div>
          <p className="text-xs text-gray-500">Current: 0 / {tenant.max_concurrent_calls}</p>
        </div>

        {/* AI Budget */}
        <div>
          <div className="flex justify-between mb-2">
            <p className="text-sm font-semibold text-gray-700">AI Monthly Budget</p>
            <p className="text-sm text-gray-600">${tenant.ai_budget_monthly.toFixed(2)}</p>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div
              className={`h-2 rounded-full transition-all ${
                budgetUsagePercent > 80 ? 'bg-red-500' :
                budgetUsagePercent > 50 ? 'bg-yellow-500' :
                'bg-blue-500'
              }`}
              style={{ width: `${budgetUsagePercent}%` }}
            ></div>
          </div>
          <p className="text-xs text-gray-500 mt-1">Used: ${(tenant.ai_budget_monthly * budgetUsagePercent / 100).toFixed(2)}</p>
        </div>
      </div>
    </div>
  )
}
