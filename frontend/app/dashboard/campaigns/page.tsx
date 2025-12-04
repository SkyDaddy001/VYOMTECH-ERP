'use client'

import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { useCampaigns } from '@/hooks/use-dashboard'
import { formatCurrency, formatDate } from '@/lib/utils'
import { FiPlus, FiEdit2, FiTrash2 } from 'react-icons/fi'

export default function CampaignsPage() {
  const { campaigns, loading } = useCampaigns()

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
    <div className="flex h-screen bg-gray-50">
      <Sidebar />
      <div className="flex-1 flex flex-col lg:ml-64">
        <Header />
        <main className="flex-1 overflow-auto pt-20 pb-6">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            {/* Page Header */}
            <div className="flex items-center justify-between mb-8">
              <div>
                <h1 className="text-3xl font-bold text-gray-900">Campaigns</h1>
                <p className="text-gray-600 mt-2">
                  Create and manage marketing campaigns
                </p>
              </div>
              <button className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition">
                <FiPlus size={20} />
                New Campaign
              </button>
            </div>

            {/* Campaigns Grid */}
            {loading ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {[...Array(6)].map((_, i) => (
                  <div
                    key={i}
                    className="bg-white rounded-lg shadow h-80 animate-pulse"
                  ></div>
                ))}
              </div>
            ) : campaigns.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-500">No campaigns found</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {campaigns.map((campaign: any) => (
                  <div
                    key={campaign.id}
                    className="bg-white rounded-lg shadow hover:shadow-lg transition p-6"
                  >
                    <div className="flex items-start justify-between mb-4">
                      <h3 className="text-lg font-semibold text-gray-900 flex-1">
                        {campaign.name}
                      </h3>
                      <span
                        className={`px-2 py-1 rounded text-xs font-semibold ${getStatusColor(
                          campaign.status
                        )}`}
                      >
                        {campaign.status}
                      </span>
                    </div>

                    <div className="space-y-4 mb-6">
                      <div className="flex justify-between text-sm">
                        <span className="text-gray-600">Budget</span>
                        <span className="font-semibold text-gray-900">
                          {formatCurrency(campaign.budget)}
                        </span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-gray-600">Spent</span>
                        <span className="font-semibold text-gray-900">
                          {formatCurrency(campaign.spent)}
                        </span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-gray-600">Leads</span>
                        <span className="font-semibold text-gray-900">
                          {campaign.leads}
                        </span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-gray-600">ROI</span>
                        <span className="font-semibold text-green-600">
                          {campaign.roi}%
                        </span>
                      </div>
                    </div>

                    <div className="border-t pt-4 mb-4">
                      <p className="text-xs text-gray-500 mb-2">Duration</p>
                      <p className="text-sm text-gray-600">
                        {formatDate(campaign.startDate)} to{' '}
                        {formatDate(campaign.endDate)}
                      </p>
                    </div>

                    <div className="flex gap-2">
                      <button className="flex-1 flex items-center justify-center gap-2 py-2 px-3 bg-blue-50 text-blue-600 rounded hover:bg-blue-100 transition text-sm font-medium">
                        <FiEdit2 size={16} />
                        Edit
                      </button>
                      <button className="py-2 px-3 bg-red-50 text-red-600 rounded hover:bg-red-100 transition">
                        <FiTrash2 size={16} />
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </main>
      </div>
    </div>
  )
}
