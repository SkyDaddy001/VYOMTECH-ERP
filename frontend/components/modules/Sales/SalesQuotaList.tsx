'use client'

import { SalesQuota } from '@/types/sales'

interface SalesQuotaListProps {
  quotas: SalesQuota[]
  loading: boolean
  onEdit: (quota: SalesQuota) => void
  onDelete: (quota: SalesQuota) => void
}

export default function SalesQuotaList({ quotas, loading, onEdit, onDelete }: SalesQuotaListProps) {
  if (loading) return <div className="text-center py-8 text-gray-500">Loading quotas...</div>

  const totalQuota = quotas.reduce((sum, q) => sum + q.quota_amount, 0)
  const avgCommissionPct = quotas.length > 0 ? quotas.reduce((sum, q) => sum + (q.commission_rate_percentage || 0), 0) / quotas.length : 0

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-gradient-to-br from-purple-50 to-purple-100 rounded-lg p-4 border border-purple-200">
          <p className="text-gray-600 text-xs font-medium">Total Quota</p>
          <p className="text-2xl font-bold text-purple-600 mt-1">₹{totalQuota.toLocaleString()}</p>
        </div>
        <div className="bg-gradient-to-br from-pink-50 to-pink-100 rounded-lg p-4 border border-pink-200">
          <p className="text-gray-600 text-xs font-medium">Avg Commission Rate</p>
          <p className="text-2xl font-bold text-pink-600 mt-1">{avgCommissionPct.toFixed(2)}%</p>
        </div>
        <div className="bg-gradient-to-br from-indigo-50 to-indigo-100 rounded-lg p-4 border border-indigo-200">
          <p className="text-gray-600 text-xs font-medium">Total Quotas</p>
          <p className="text-2xl font-bold text-indigo-600 mt-1">{quotas.length}</p>
        </div>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Sales Executive</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Quarter</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Year</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Target Bookings</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Quota Amount</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Commission</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {quotas.map((quota) => (
              <tr key={quota.id} className="hover:bg-gray-50 transition">
                <td className="px-6 py-4 text-sm font-medium text-gray-900">{quota.sales_executive_name || quota.sales_executive_id}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{quota.quarter}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{quota.year}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{quota.quota_bookings}</td>
                <td className="px-6 py-4 text-sm font-semibold text-gray-900">₹{quota.quota_amount.toLocaleString()}</td>
                <td className="px-6 py-4 text-sm text-blue-600">
                  ₹{quota.commission_rate_per_booking.toLocaleString()} 
                  {quota.commission_rate_percentage ? ` + ${quota.commission_rate_percentage}%` : ''}
                </td>
                <td className="px-6 py-4 text-sm space-x-2">
                  <button
                    onClick={() => onEdit(quota)}
                    className="text-blue-600 hover:text-blue-900 font-medium"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(quota)}
                    className="text-red-600 hover:text-red-900 font-medium"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
