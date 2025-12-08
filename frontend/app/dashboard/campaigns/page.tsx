'use client'

import { useEffect, useState } from 'react'
import { Search, Plus, MoreVertical, TrendingUp, Calendar } from 'lucide-react'
import { format } from 'date-fns'

interface Campaign {
  id: string
  name: string
  description?: string
  type: 'email' | 'sms' | 'call' | 'social'
  status: 'draft' | 'active' | 'paused' | 'completed'
  startDate: string
  endDate?: string
  leads: number
  conversions: number
  conversionRate: number
  budget?: number
  spent?: number
}

export default function CampaignsPage() {
  const [campaigns, setCampaigns] = useState<Campaign[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState<string>('all')

  useEffect(() => {
    fetchCampaigns()
  }, [])

  const fetchCampaigns = async () => {
    try {
      setLoading(true)
      // Mock data
      const mockCampaigns: Campaign[] = [
        {
          id: '1',
          name: 'Q1 Email Campaign',
          description: 'Quarterly email outreach',
          type: 'email',
          status: 'active',
          startDate: new Date(Date.now() - 604800000).toISOString(),
          endDate: new Date(Date.now() + 1209600000).toISOString(),
          leads: 1250,
          conversions: 87,
          conversionRate: 6.96,
          budget: 5000,
          spent: 2500
        },
        {
          id: '2',
          name: 'Social Media Blitz',
          description: 'LinkedIn and Facebook campaign',
          type: 'social',
          status: 'active',
          startDate: new Date(Date.now() - 259200000).toISOString(),
          leads: 890,
          conversions: 45,
          conversionRate: 5.06,
          budget: 3000,
          spent: 1500
        },
        {
          id: '3',
          name: 'Holiday Special',
          description: 'End of year promotions',
          type: 'email',
          status: 'completed',
          startDate: new Date(Date.now() - 2592000000).toISOString(),
          endDate: new Date(Date.now() - 1209600000).toISOString(),
          leads: 2100,
          conversions: 312,
          conversionRate: 14.86,
          budget: 8000,
          spent: 8000
        }
      ]
      setCampaigns(mockCampaigns)
    } catch (error) {
      console.error('Error fetching campaigns:', error)
    } finally {
      setLoading(false)
    }
  }

  const filteredCampaigns = campaigns.filter(campaign => {
    const matchesSearch = campaign.name.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesStatus = filterStatus === 'all' || campaign.status === filterStatus
    return matchesSearch && matchesStatus
  })

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      draft: 'bg-gray-100 text-gray-800',
      active: 'bg-green-100 text-green-800',
      paused: 'bg-yellow-100 text-yellow-800',
      completed: 'bg-blue-100 text-blue-800'
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const getTypeColor = (type: string) => {
    const colors: Record<string, string> = {
      email: 'bg-blue-500',
      sms: 'bg-green-500',
      call: 'bg-purple-500',
      social: 'bg-pink-500'
    }
    return colors[type] || 'bg-gray-500'
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Campaigns</h1>
          <p className="text-gray-600 mt-2">Create and manage your marketing campaigns</p>
        </div>
        <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition flex items-center">
          <Plus className="w-4 h-4 mr-2" />
          New Campaign
        </button>
      </div>

      {/* Search and Filters */}
      <div className="bg-white rounded-lg shadow p-4 mb-6">
        <div className="flex flex-col md:flex-row gap-4">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-3 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search campaigns..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="all">All Statuses</option>
            <option value="draft">Draft</option>
            <option value="active">Active</option>
            <option value="paused">Paused</option>
            <option value="completed">Completed</option>
          </select>
        </div>
      </div>

      {/* Campaigns Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {loading ? (
          <div className="col-span-full text-center py-12">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        ) : filteredCampaigns.length === 0 ? (
          <div className="col-span-full text-center py-12">
            <p className="text-gray-600">No campaigns found</p>
          </div>
        ) : (
          filteredCampaigns.map((campaign) => (
            <div key={campaign.id} className="bg-white rounded-lg shadow hover:shadow-md transition">
              {/* Card Header */}
              <div className="p-6 border-b">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex-1">
                    <h3 className="text-lg font-bold text-gray-900">{campaign.name}</h3>
                    <p className="text-sm text-gray-600 mt-1">{campaign.description}</p>
                  </div>
                  <button className="p-2 hover:bg-gray-100 rounded-lg transition">
                    <MoreVertical className="w-4 h-4 text-gray-600" />
                  </button>
                </div>
                <div className="flex items-center gap-2">
                  <span className={`px-2 py-1 rounded text-xs font-medium text-white ${getTypeColor(campaign.type)}`}>
                    {campaign.type.toUpperCase()}
                  </span>
                  <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(campaign.status)}`}>
                    {campaign.status.charAt(0).toUpperCase() + campaign.status.slice(1)}
                  </span>
                </div>
              </div>

              {/* Card Body */}
              <div className="p-6 space-y-4">
                {/* Timeline */}
                <div>
                  <div className="flex items-center text-sm text-gray-600 mb-1">
                    <Calendar className="w-4 h-4 mr-2" />
                    {format(new Date(campaign.startDate), 'MMM dd, yyyy')} - {campaign.endDate ? format(new Date(campaign.endDate), 'MMM dd, yyyy') : 'Ongoing'}
                  </div>
                </div>

                {/* Metrics */}
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-xs text-gray-600 uppercase font-semibold">Leads</p>
                    <p className="text-2xl font-bold text-gray-900 mt-1">{campaign.leads}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-600 uppercase font-semibold">Conversions</p>
                    <p className="text-2xl font-bold text-gray-900 mt-1">{campaign.conversions}</p>
                  </div>
                </div>

                {/* Conversion Rate */}
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <p className="text-sm text-gray-600 uppercase font-semibold flex items-center">
                      <TrendingUp className="w-4 h-4 mr-1" />
                      Conversion Rate
                    </p>
                    <p className="text-sm font-bold text-gray-900">{campaign.conversionRate.toFixed(2)}%</p>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-green-500 h-2 rounded-full"
                      style={{ width: `${Math.min(campaign.conversionRate, 100)}%` }}
                    ></div>
                  </div>
                </div>

                {/* Budget */}
                {campaign.budget && (
                  <div>
                    <p className="text-xs text-gray-600 uppercase font-semibold mb-2">Budget Used</p>
                    <div className="flex items-center justify-between mb-1">
                      <p className="text-sm text-gray-600">
                        ${campaign.spent?.toLocaleString()} / ${campaign.budget?.toLocaleString()}
                      </p>
                      <p className="text-sm font-semibold text-gray-900">
                        {campaign.spent && campaign.budget ? Math.round((campaign.spent / campaign.budget) * 100) : 0}%
                      </p>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div
                        className="bg-blue-500 h-2 rounded-full"
                        style={{ width: `${campaign.spent && campaign.budget ? Math.min((campaign.spent / campaign.budget) * 100, 100) : 0}%` }}
                      ></div>
                    </div>
                  </div>
                )}
              </div>
            </div>
          ))
        )}
      </div>

      {/* Summary */}
      {!loading && filteredCampaigns.length > 0 && (
        <div className="mt-8 bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-bold text-gray-900 mb-4">Summary</h3>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            <div>
              <p className="text-sm text-gray-600">Total Campaigns</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">{filteredCampaigns.length}</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Total Leads</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">
                {filteredCampaigns.reduce((sum, c) => sum + c.leads, 0).toLocaleString()}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Total Conversions</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">
                {filteredCampaigns.reduce((sum, c) => sum + c.conversions, 0).toLocaleString()}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Avg Conversion Rate</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">
                {(filteredCampaigns.reduce((sum, c) => sum + c.conversionRate, 0) / filteredCampaigns.length).toFixed(2)}%
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
