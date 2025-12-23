'use client'

import { useState, useMemo } from 'react'
import { useCampaigns } from '@/hooks/use-api'
import { ProtectedRoute } from '@/hooks/use-auth'
import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { FiSearch, FiPlus, FiBarChart2, FiCalendar, FiDollarSign, FiTrendingUp } from 'react-icons/fi'
import { format } from 'date-fns'

interface Campaign {
  id: string | number
  name: string
  description?: string
  status: string
  target_leads?: number
  generated_leads?: number
  converted_leads?: number
  budget?: number
  spent_budget?: number
  start_date: string
  end_date?: string
  created_at: string
}

export default function CampaignsPage() {
  const { data: campaigns, loading, error, refetch } = useCampaigns({ limit: 100 })
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState<string>('all')
  const [sortBy, setSortBy] = useState<'created' | 'name' | 'budget' | 'leads'>('created')

  const filteredCampaigns = useMemo(() => {
    let result = [...(campaigns || [])]

    if (filterStatus !== 'all') {
      result = result.filter(c => c.status === filterStatus)
    }

    if (searchTerm) {
      const query = searchTerm.toLowerCase()
      result = result.filter(c =>
        c.name.toLowerCase().includes(query) ||
        (c.description || '').toLowerCase().includes(query)
      )
    }

    result.sort((a, b) => {
      switch (sortBy) {
        case 'name':
          return a.name.localeCompare(b.name)
        case 'budget':
          return (b.budget || 0) - (a.budget || 0)
        case 'leads':
          return (b.generated_leads || 0) - (a.generated_leads || 0)
        case 'created':
        default:
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      }
    })

    return result
  }, [campaigns, searchTerm, filterStatus, sortBy])

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'active':
        return 'bg-green-100 text-green-800'
      case 'completed':
        return 'bg-blue-100 text-blue-800'
      case 'paused':
        return 'bg-yellow-100 text-yellow-800'
      case 'planned':
        return 'bg-gray-100 text-gray-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getConversionRate = (generated: number, converted: number) => {
    if (generated === 0) return 0
    return ((converted / generated) * 100).toFixed(1)
  }

  const getBudgetUsage = (spent: number, budget: number) => {
    if (budget === 0) return 0
    return ((spent / budget) * 100).toFixed(1)
  }

  return (
    <ProtectedRoute>
      <div className="flex h-screen bg-gray-50">
        <Sidebar />
        <div className="flex-1 flex flex-col lg:ml-64">
          <Header />
          <main className="flex-1 overflow-auto pt-20 pb-6">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
              {/* Header */}
              <div className="mb-8 flex items-center justify-between">
                <div>
                  <h1 className="text-3xl font-bold text-gray-900">Campaigns</h1>
                  <p className="text-gray-600 mt-2">Create and manage marketing campaigns</p>
                </div>
                <button className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg flex items-center gap-2 transition">
                  <FiPlus className="text-lg" />
                  Create Campaign
                </button>
              </div>

              {/* Filters & Search */}
              <div className="bg-white rounded-lg shadow p-4 mb-6">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  {/* Search */}
                  <div className="relative">
                    <FiSearch className="absolute left-3 top-3 text-gray-400" />
                    <input
                      type="text"
                      placeholder="Search campaigns..."
                      value={searchTerm}
                      onChange={(e) => setSearchTerm(e.target.value)}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>

                  {/* Status Filter */}
                  <select
                    value={filterStatus}
                    onChange={(e) => setFilterStatus(e.target.value)}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="all">All Status</option>
                    <option value="active">Active</option>
                    <option value="completed">Completed</option>
                    <option value="paused">Paused</option>
                    <option value="planned">Planned</option>
                  </select>

                  {/* Sort */}
                  <select
                    value={sortBy}
                    onChange={(e) => setSortBy(e.target.value as any)}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="created">Newest First</option>
                    <option value="name">By Name</option>
                    <option value="budget">By Budget</option>
                    <option value="leads">By Leads</option>
                  </select>
                </div>

                <div className="mt-4 text-sm text-gray-600">
                  Showing {filteredCampaigns.length} of {(campaigns || []).length} campaigns
                </div>
              </div>

              {/* Campaigns Grid */}
              {filteredCampaigns.length === 0 ? (
                <div className="bg-white rounded-lg shadow p-12 text-center">
                  <div className="text-gray-500 mb-4 text-5xl">🎯</div>
                  <p className="text-gray-600 text-lg font-medium">No campaigns found</p>
                  <p className="text-gray-500 mt-1">Create a new campaign to get started</p>
                </div>
              ) : (
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                  {filteredCampaigns.map((campaign) => (
                    <div key={campaign.id} className="bg-white rounded-lg shadow hover:shadow-md transition p-6">
                      {/* Header */}
                      <div className="flex items-start justify-between mb-4">
                        <div>
                          <h3 className="text-lg font-semibold text-gray-900">{campaign.name}</h3>
                          <p className="text-sm text-gray-600 mt-1">{campaign.description}</p>
                        </div>
                        <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(campaign.status)}`}>
                          {campaign.status}
                        </span>
                      </div>

                      {/* Timeline */}
                      <div className="flex items-center gap-2 text-sm text-gray-600 mb-6">
                        <FiCalendar className="text-gray-400" />
                        {format(new Date(campaign.start_date), 'MMM d, yyyy')}
                        {campaign.end_date && ` - ${format(new Date(campaign.end_date), 'MMM d, yyyy')}`}
                      </div>

                      {/* Metrics */}
                      <div className="grid grid-cols-2 gap-4 mb-6">
                        {/* Leads */}
                        <div className="bg-gray-50 rounded-lg p-4">
                          <div className="flex items-center justify-between mb-2">
                            <span className="text-sm text-gray-600">Leads Generated</span>
                            <FiBarChart2 className="text-blue-500" />
                          </div>
                          <div className="text-2xl font-bold text-gray-900">
                            {campaign.generated_leads || 0}
                          </div>
                          <div className="text-xs text-gray-500 mt-1">
                            Target: {campaign.target_leads || 0}
                          </div>
                        </div>

                        {/* Conversion */}
                        <div className="bg-gray-50 rounded-lg p-4">
                          <div className="flex items-center justify-between mb-2">
                            <span className="text-sm text-gray-600">Conversion Rate</span>
                            <FiTrendingUp className="text-green-500" />
                          </div>
                          <div className="text-2xl font-bold text-gray-900">
                            {getConversionRate(campaign.generated_leads || 0, campaign.converted_leads || 0)}%
                          </div>
                          <div className="text-xs text-gray-500 mt-1">
                            {campaign.converted_leads || 0} converted
                          </div>
                        </div>

                        {/* Budget */}
                        <div className="bg-gray-50 rounded-lg p-4">
                          <div className="flex items-center justify-between mb-2">
                            <span className="text-sm text-gray-600">Budget Spent</span>
                            <FiDollarSign className="text-purple-500" />
                          </div>
                          <div className="text-2xl font-bold text-gray-900">
                            ${(campaign.spent_budget || 0).toLocaleString()}
                          </div>
                          <div className="text-xs text-gray-500 mt-1">
                            of ${(campaign.budget || 0).toLocaleString()}
                          </div>
                        </div>

                        {/* Cost Per Lead */}
                        <div className="bg-gray-50 rounded-lg p-4">
                          <div className="flex items-center justify-between mb-2">
                            <span className="text-sm text-gray-600">Cost per Lead</span>
                            <FiDollarSign className="text-orange-500" />
                          </div>
                          <div className="text-2xl font-bold text-gray-900">
                            ${(campaign.generated_leads && campaign.spent_budget ? (campaign.spent_budget / campaign.generated_leads).toFixed(0) : 0)}
                          </div>
                          <div className="text-xs text-gray-500 mt-1">
                            Average cost
                          </div>
                        </div>
                      </div>

                      {/* Budget Progress Bar */}
                      <div className="mb-6">
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-sm font-medium text-gray-700">Budget Usage</span>
                          <span className="text-sm text-gray-600">
                            {getBudgetUsage(campaign.spent_budget || 0, campaign.budget || 0)}%
                          </span>
                        </div>
                        <div className="w-full bg-gray-200 rounded-full h-2">
                          <div
                            className="bg-blue-600 h-2 rounded-full transition-all"
                            style={{ width: `${Math.min(parseFloat(getBudgetUsage(campaign.spent_budget || 0, campaign.budget || 0)), 100)}%` }}
                          ></div>
                        </div>
                      </div>

                      {/* Lead Progress Bar */}
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-sm font-medium text-gray-700">Lead Generation</span>
                          <span className="text-sm text-gray-600">
                            {campaign.generated_leads && campaign.target_leads ? Math.round((campaign.generated_leads / campaign.target_leads) * 100) : 0}%
                          </span>
                        </div>
                        <div className="w-full bg-gray-200 rounded-full h-2">
                          <div
                            className="bg-green-600 h-2 rounded-full transition-all"
                            style={{
                              width: `${Math.min(campaign.generated_leads && campaign.target_leads ? (campaign.generated_leads / campaign.target_leads) * 100 : 0, 100)}%`,
                            }}
                          ></div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </main>
        </div>
      </div>
    </ProtectedRoute>
  )
}
