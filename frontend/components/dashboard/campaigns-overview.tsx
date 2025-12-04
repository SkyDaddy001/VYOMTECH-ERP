'use client'

import { useCampaigns } from '@/hooks/use-dashboard'
import { formatCurrency, formatDate } from '@/lib/utils'

export const CampaignsOverview = ({ limit = 5 }: { limit?: number }) => {
  const { campaigns, loading } = useCampaigns()

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold mb-4">Active Campaigns</h3>
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="bg-gray-200 h-16 rounded animate-pulse"></div>
          ))}
        </div>
      </div>
    )
  }

  if (!campaigns || campaigns.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold mb-4">Active Campaigns</h3>
        <p className="text-gray-500 text-center py-8">No campaigns found</p>
      </div>
    )
  }

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      active: 'bg-green-100 text-green-800',
      draft: 'bg-gray-100 text-gray-800',
      completed: 'bg-blue-100 text-blue-800',
      paused: 'bg-yellow-100 text-yellow-800',
    }
    return colors[status?.toLowerCase()] || 'bg-gray-100 text-gray-800'
  }

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h3 className="text-lg font-semibold mb-4">Active Campaigns</h3>
      <div className="space-y-4">
        {campaigns.slice(0, limit).map((campaign: any) => (
          <div key={campaign.id} className="border rounded-lg p-4 hover:shadow transition">
            <div className="flex items-start justify-between mb-2">
              <div className="flex-1">
                <h4 className="font-semibold text-gray-900">{campaign.name}</h4>
                <p className="text-sm text-gray-600 mt-1">
                  {campaign.leads} leads • Budget: {formatCurrency(campaign.budget)}
                </p>
              </div>
              <span
                className={`px-2 py-1 rounded text-xs font-semibold ${getStatusColor(
                  campaign.status
                )}`}
              >
                {campaign.status}
              </span>
            </div>
            <div className="grid grid-cols-3 gap-4 text-sm mt-3">
              <div>
                <p className="text-gray-600">Spent</p>
                <p className="font-semibold text-gray-900">
                  {formatCurrency(campaign.spent)}
                </p>
              </div>
              <div>
                <p className="text-gray-600">ROI</p>
                <p className="font-semibold text-green-600">{campaign.roi}%</p>
              </div>
              <div>
                <p className="text-gray-600">End Date</p>
                <p className="font-semibold text-gray-900">
                  {formatDate(campaign.endDate)}
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
