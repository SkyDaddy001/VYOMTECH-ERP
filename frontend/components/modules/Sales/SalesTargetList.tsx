'use client'

import { SalesTarget } from '@/types/sales'

interface SalesTargetListProps {
  targets: SalesTarget[]
  loading: boolean
  onEdit: (target: SalesTarget) => void
  onDelete: (target: SalesTarget) => void
}

export default function SalesTargetList({ targets, loading, onEdit, onDelete }: SalesTargetListProps) {
  if (loading) return <div className="text-center py-8 text-gray-500">Loading targets...</div>

  const totalTarget = targets.reduce((sum, t) => sum + t.target_amount, 0)
  const totalAchieved = targets.reduce((sum, t) => sum + t.achieved_amount, 0)

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-gradient-to-br from-orange-50 to-orange-100 rounded-lg p-4 border border-orange-200">
          <p className="text-gray-600 text-xs font-medium">Total Target</p>
          <p className="text-2xl font-bold text-orange-600 mt-1">₹{totalTarget.toLocaleString()}</p>
        </div>
        <div className="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-4 border border-green-200">
          <p className="text-gray-600 text-xs font-medium">Total Achieved</p>
          <p className="text-2xl font-bold text-green-600 mt-1">₹{totalAchieved.toLocaleString()}</p>
        </div>
        <div className="bg-gradient-to-br from-blue-50 to-blue-100 rounded-lg p-4 border border-blue-200">
          <p className="text-gray-600 text-xs font-medium">Avg Achievement</p>
          <p className="text-2xl font-bold text-blue-600 mt-1">
            {targets.length > 0 ? (targets.reduce((sum, t) => sum + t.achievement_percentage, 0) / targets.length).toFixed(1) : 0}%
          </p>
        </div>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Sales Executive</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Period</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Target Bookings</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Achieved</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Target Amount</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Achievement %</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Status</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {targets.map((target) => (
              <tr key={target.id} className="hover:bg-gray-50 transition">
                <td className="px-6 py-4 text-sm font-medium text-gray-900">{target.sales_executive_name || target.sales_executive_id}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{target.period}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{target.target_bookings}</td>
                <td className="px-6 py-4 text-sm font-semibold text-green-600">{target.achieved_bookings}</td>
                <td className="px-6 py-4 text-sm font-semibold text-gray-900">₹{target.target_amount.toLocaleString()}</td>
                <td className="px-6 py-4">
                  <div className="w-24">
                    <div className="bg-gray-200 rounded-full h-2 overflow-hidden">
                      <div
                        className={`h-full ${target.achievement_percentage >= 100 ? 'bg-green-500' : 'bg-blue-500'}`}
                        style={{ width: `${Math.min(target.achievement_percentage, 100)}%` }}
                      />
                    </div>
                    <p className="text-xs text-gray-600 mt-1">{target.achievement_percentage.toFixed(1)}%</p>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <span
                    className={`text-xs font-semibold px-3 py-1 rounded-full ${
                      target.status === 'completed'
                        ? 'bg-green-100 text-green-800'
                        : target.status === 'exceeded'
                        ? 'bg-blue-100 text-blue-800'
                        : target.status === 'in_progress'
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-gray-100 text-gray-800'
                    }`}
                  >
                    {target.status.replace('_', ' ').charAt(0).toUpperCase() + target.status.slice(1)}
                  </span>
                </td>
                <td className="px-6 py-4 text-sm space-x-2">
                  <button
                    onClick={() => onEdit(target)}
                    className="text-blue-600 hover:text-blue-900 font-medium"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(target)}
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
