'use client'

import { Campaign } from '@/types/marketing'

interface CampaignListProps {
  campaigns: Campaign[]
  loading: boolean
  onEdit: (campaign: Campaign) => void
  onDelete: (campaign: Campaign) => void
}

export default function CampaignList({ campaigns, loading, onEdit, onDelete }: CampaignListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading campaigns...</p>
        </div>
      </div>
    )
  }

  if (campaigns.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No campaigns yet. Create your first campaign to get started.</p>
      </div>
    )
  }

  const activeCampaigns = campaigns.filter((c) => c.status === 'active').length
  const totalBudget = campaigns.reduce((sum, c) => sum + (c.budget || 0), 0)
  const totalSpent = campaigns.reduce((sum, c) => sum + (c.spent || 0), 0)

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{campaigns.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Active</p>
          <p className="text-2xl font-bold text-green-600 mt-1">{activeCampaigns}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total Budget</p>
          <p className="text-2xl font-bold text-blue-600 mt-1">₹{totalBudget.toLocaleString()}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Spent</p>
          <p className="text-2xl font-bold text-red-600 mt-1">₹{totalSpent.toLocaleString()}</p>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-lg shadow overflow-x-auto border border-gray-200">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Campaign</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Duration</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Budget</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Spent</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {campaigns.map((campaign) => {
              const progress = campaign.budget ? ((campaign.spent || 0) / campaign.budget) * 100 : 0
              return (
                <tr key={campaign.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{campaign.name}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    <span className="px-2 py-1 bg-purple-100 text-purple-800 rounded text-xs">{campaign.campaign_type?.replace('_', ' ') || 'Campaign'}</span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    {new Date(campaign.start_date).toLocaleDateString()} - {new Date(campaign.end_date).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">₹{(campaign.budget || 0).toLocaleString()}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <div className="flex items-center gap-2">
                      <div className="w-16 bg-gray-200 rounded h-2">
                        <div className="bg-blue-600 h-2 rounded" style={{ width: `${Math.min(progress, 100)}%` }}></div>
                      </div>
                      <span className="text-xs text-gray-600">{Math.round(progress)}%</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <span
                      className={`px-2 py-1 rounded text-xs font-medium ${
                        campaign.status === 'active'
                          ? 'bg-green-100 text-green-800'
                          : campaign.status === 'draft'
                          ? 'bg-gray-100 text-gray-800'
                          : campaign.status === 'paused'
                          ? 'bg-yellow-100 text-yellow-800'
                          : 'bg-red-100 text-red-800'
                      }`}
                    >
                      {campaign.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                    <button onClick={() => onEdit(campaign)} className="text-green-600 hover:text-green-900 font-medium">
                      Edit
                    </button>
                    <button onClick={() => onDelete(campaign)} className="text-red-600 hover:text-red-900 font-medium">
                      Delete
                    </button>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}
