'use client'

import { useCampaigns } from '@/hooks/use-dashboard'
import { formatCurrency, formatDate } from '@/lib/utils'

export const CampaignsOverview = ({ limit = 5 }: { limit?: number }) => {
  const { campaigns, loading } = useCampaigns()

  if (loading) {
    return (
      <div className="bg-white rounded-sm p-6">
        <h3 className="text-sm font-semibold uppercase tracking-wide mb-6">Active Campaigns</h3>
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="bg-gray-100 h-16 rounded-sm animate-pulse"></div>
          ))}
        </div>
      </div>
    )
  }

  if (!campaigns || campaigns.length === 0) {
    return (
      <div className="bg-white rounded-sm p-6">
        <h3 className="text-sm font-semibold uppercase tracking-wide mb-6">Active Campaigns</h3>
        <p className="text-gray-500 text-center py-8 text-sm">No campaigns found</p>
      </div>
    )
  }

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      active: 'bg-gray-100 text-gray-900',
      draft: 'bg-gray-100 text-gray-900',
      completed: 'bg-gray-100 text-gray-900',
      paused: 'bg-gray-50 text-gray-600',
    }
    return colors[status?.toLowerCase()] || 'bg-gray-100 text-gray-900'
  }

  return (
    <div className="bg-white rounded-sm p-6">
      <h3 className="text-sm font-semibold uppercase tracking-wide mb-6">Active Campaigns</h3>
      <div className="space-y-3">
        {campaigns.slice(0, limit).map((campaign: any) => (
          <div key={campaign.id} className="border-b border-gray-100 p-4 hover:bg-gray-50 transition-colors">
            <div className="flex items-start justify-between mb-2">
              <div className="flex-1">
                <h4 className="font-semibold text-gray-900 text-sm">{campaign.name}</h4>
                <p className="text-xs text-gray-600 mt-1">
                  {campaign.leads} leads • Budget: {formatCurrency(campaign.budget)}
                </p>
              </div>
              <span
                className={`px-2 py-0.5 rounded-sm text-xs font-medium ${getStatusColor(
                  campaign.status
                )}`}
              >
                {campaign.status}
              </span>
            </div>
            <div className="grid grid-cols-3 gap-4 text-xs mt-3">
              <div>
                <p className="text-gray-600">Spent</p>
                <p className="font-semibold text-gray-900 mt-0.5">
                  {formatCurrency(campaign.spent)}
                </p>
              </div>
              <div>
                <p className="text-gray-600">ROI</p>
                <p className="font-semibold text-gray-900 mt-0.5">{campaign.roi}%</p>
              </div>
              <div>
                <p className="text-gray-600">End Date</p>
                <p className="font-semibold text-gray-900 mt-0.5">
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
