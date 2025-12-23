'use client'

import { useState, useMemo } from 'react'
import { useOpportunities } from '@/hooks/use-api'
import { ProtectedRoute } from '@/hooks/use-auth'
import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { FiSearch, FiPlus, FiTrendingUp, FiDollarSign, FiCalendar, FiUser } from 'react-icons/fi'
import { format } from 'date-fns'

interface Opportunity {
  id: string | number
  name: string
  lead_id?: string
  lead_name?: string
  company?: string
  amount?: number
  stage?: string
  probability?: number
  expected_close_date?: string
  description?: string
  owner?: string
  created_at: string
  updated_at: string
}

export default function OpportunitiesPage() {
  const { data: opportunities, loading, error } = useOpportunities({ limit: 100 })
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStage, setFilterStage] = useState<string>('all')
  const [sortBy, setSortBy] = useState<'created' | 'amount' | 'probability' | 'name'>('created')

  const filteredOpportunities = useMemo(() => {
    let result = [...(opportunities || [])]

    if (filterStage !== 'all') {
      result = result.filter(o => o.stage === filterStage)
    }

    if (searchTerm) {
      const query = searchTerm.toLowerCase()
      result = result.filter(o =>
        (o.name || '').toLowerCase().includes(query) ||
        (o.lead_name || '').toLowerCase().includes(query) ||
        (o.company || '').toLowerCase().includes(query)
      )
    }

    result.sort((a, b) => {
      switch (sortBy) {
        case 'amount':
          return (b.amount || 0) - (a.amount || 0)
        case 'probability':
          return (b.probability || 0) - (a.probability || 0)
        case 'name':
          return (a.name || '').localeCompare(b.name || '')
        case 'created':
        default:
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      }
    })

    return result
  }, [opportunities, searchTerm, filterStage, sortBy])

  const getStageColor = (stage?: string) => {
    switch (stage?.toLowerCase()) {
      case 'prospect':
        return 'bg-gray-100 text-gray-800'
      case 'qualification':
        return 'bg-blue-100 text-blue-800'
      case 'proposal':
        return 'bg-yellow-100 text-yellow-800'
      case 'negotiation':
        return 'bg-orange-100 text-orange-800'
      case 'won':
        return 'bg-green-100 text-green-800'
      case 'lost':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getProbabilityColor = (probability?: number) => {
    if (!probability) return 'bg-gray-100'
    if (probability < 25) return 'bg-red-100'
    if (probability < 50) return 'bg-orange-100'
    if (probability < 75) return 'bg-yellow-100'
    return 'bg-green-100'
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
                  <h1 className="text-3xl font-bold text-gray-900">Sales Opportunities</h1>
                  <p className="text-gray-600 mt-2">Track and manage sales pipeline</p>
                </div>
                <button className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg flex items-center gap-2 transition">
                  <FiPlus className="text-lg" />
                  New Opportunity
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
                      placeholder="Search opportunities..."
                      value={searchTerm}
                      onChange={(e) => setSearchTerm(e.target.value)}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>

                  {/* Stage Filter */}
                  <select
                    value={filterStage}
                    onChange={(e) => setFilterStage(e.target.value)}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="all">All Stages</option>
                    <option value="prospect">Prospect</option>
                    <option value="qualification">Qualification</option>
                    <option value="proposal">Proposal</option>
                    <option value="negotiation">Negotiation</option>
                    <option value="won">Won</option>
                    <option value="lost">Lost</option>
                  </select>

                  {/* Sort */}
                  <select
                    value={sortBy}
                    onChange={(e) => setSortBy(e.target.value as any)}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="created">Newest First</option>
                    <option value="amount">By Amount (High to Low)</option>
                    <option value="probability">By Win Probability</option>
                    <option value="name">By Name</option>
                  </select>
                </div>

                <div className="mt-4 text-sm text-gray-600">
                  Showing {filteredOpportunities.length} of {(opportunities || []).length} opportunities
                </div>
              </div>

              {/* Opportunities List */}
              {filteredOpportunities.length === 0 ? (
                <div className="bg-white rounded-lg shadow p-12 text-center">
                  <div className="text-gray-500 mb-4 text-5xl">📊</div>
                  <p className="text-gray-600 text-lg font-medium">No opportunities found</p>
                  <p className="text-gray-500 mt-1">Create a new opportunity to get started</p>
                </div>
              ) : (
                <div className="bg-white rounded-lg shadow overflow-hidden">
                  <table className="w-full">
                    <thead className="bg-gray-50 border-b border-gray-200">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Opportunity</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Company</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Stage</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Amount</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Probability</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Expected Close</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                      {filteredOpportunities.map((opp) => (
                        <tr key={opp.id} className="hover:bg-gray-50 transition">
                          <td className="px-6 py-4">
                            <div>
                              <p className="font-semibold text-gray-900">{opp.name}</p>
                              {opp.lead_name && <p className="text-sm text-gray-600">{opp.lead_name}</p>}
                            </div>
                          </td>
                          <td className="px-6 py-4 text-gray-700">{opp.company || '-'}</td>
                          <td className="px-6 py-4">
                            <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getStageColor(opp.stage)}`}>
                              {opp.stage || 'Unknown'}
                            </span>
                          </td>
                          <td className="px-6 py-4">
                            <div className="flex items-center gap-2 text-gray-900 font-semibold">
                              <FiDollarSign className="text-gray-400" />
                              ${(opp.amount || 0).toLocaleString()}
                            </div>
                          </td>
                          <td className="px-6 py-4">
                            <div className="flex items-center gap-2">
                              <div className={`px-2 py-1 rounded text-xs font-semibold text-gray-900 ${getProbabilityColor(opp.probability)}`}>
                                {opp.probability || 0}%
                              </div>
                              <div className="w-16 bg-gray-200 rounded-full h-2">
                                <div
                                  className="bg-blue-600 h-2 rounded-full"
                                  style={{ width: `${Math.min(opp.probability || 0, 100)}%` }}
                                ></div>
                              </div>
                            </div>
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-600">
                            {opp.expected_close_date ? format(new Date(opp.expected_close_date), 'MMM d, yyyy') : '-'}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          </main>
        </div>
      </div>
    </ProtectedRoute>
  )
}
