'use client'

import { useLeads } from '@/hooks/use-dashboard'
import { formatDate } from '@/lib/utils'

export const RecentLeads = ({ limit = 5 }: { limit?: number }) => {
  const { leads, loading } = useLeads({ limit })

  if (loading) {
    return (
      <div className="bg-white rounded-sm p-6">
        <h3 className="text-sm font-semibold uppercase tracking-wide mb-6">Recent Leads</h3>
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="bg-gray-100 h-12 rounded-sm animate-pulse"></div>
          ))}
        </div>
      </div>
    )
  }

  if (!leads || leads.length === 0) {
    return (
      <div className="bg-white rounded-sm p-6">
        <h3 className="text-sm font-semibold uppercase tracking-wide mb-6">Recent Leads</h3>
        <p className="text-gray-500 text-center py-8 text-sm">No leads found</p>
      </div>
    )
  }

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      new: 'bg-blue-100 text-blue-800',
      qualified: 'bg-green-100 text-green-800',
      contacted: 'bg-yellow-100 text-yellow-800',
      converted: 'bg-purple-100 text-purple-800',
      lost: 'bg-red-100 text-red-800',
    }
    return colors[status?.toLowerCase()] || 'bg-gray-100 text-gray-800'
  }

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h3 className="text-lg font-semibold mb-4">Recent Leads</h3>
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="text-gray-600 border-b">
            <tr>
              <th className="text-left py-3 px-2 font-medium">Name</th>
              <th className="text-left py-3 px-2 font-medium">Email</th>
              <th className="text-left py-3 px-2 font-medium">Score</th>
              <th className="text-left py-3 px-2 font-medium">Status</th>
              <th className="text-left py-3 px-2 font-medium">Date</th>
            </tr>
          </thead>
          <tbody>
            {leads.slice(0, limit).map((lead: any) => (
              <tr key={lead.id} className="border-b hover:bg-gray-50">
                <td className="py-3 px-2 font-medium text-gray-900">{lead.name}</td>
                <td className="py-3 px-2 text-gray-600">{lead.email}</td>
                <td className="py-3 px-2">
                  <div className="w-16 bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-blue-600 h-2 rounded-full"
                      style={{ width: `${lead.score}%` }}
                    ></div>
                  </div>
                </td>
                <td className="py-3 px-2">
                  <span
                    className={`px-2 py-1 rounded-full text-xs font-semibold ${getStatusColor(
                      lead.status
                    )}`}
                  >
                    {lead.status}
                  </span>
                </td>
                <td className="py-3 px-2 text-gray-600">
                  {formatDate(lead.createdAt)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
