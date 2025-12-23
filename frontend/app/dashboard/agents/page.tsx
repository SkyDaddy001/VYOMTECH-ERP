'use client'

import { useState, useMemo } from 'react'
import { useAgents } from '@/hooks/use-api'
import { ProtectedRoute } from '@/hooks/use-auth'
import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { FiSearch, FiPlus, FiPhone, FiMail, FiActivity, FiStar } from 'react-icons/fi'
import { format } from 'date-fns'

interface Agent {
  id: string | number
  name: string
  email: string
  phone?: string
  status: 'active' | 'inactive' | 'on-break' | 'offline'
  availability: string
  total_calls?: number
  successful_calls?: number
  rating?: number
  created_at: string
  updated_at: string
}

export default function AgentsPage() {
  const { data: agents, loading, error, refetch } = useAgents({ limit: 100 })
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState<string>('all')
  const [sortBy, setSortBy] = useState<'name' | 'calls' | 'rating' | 'created'>('created')

  // Filter and search
  const filteredAgents = useMemo(() => {
    let result = agents as Agent[]

    // Filter by status
    if (filterStatus !== 'all') {
      result = result.filter(agent => 
        (agent.status || '').toLowerCase() === filterStatus.toLowerCase()
      )
    }

    // Search by name, email, or phone
    if (searchTerm) {
      const query = searchTerm.toLowerCase()
      result = result.filter(agent =>
        (agent.name || '').toLowerCase().includes(query) ||
        (agent.email || '').toLowerCase().includes(query) ||
        (agent.phone || '').toLowerCase().includes(query)
      )
    }

    // Sort
    result.sort((a, b) => {
      switch (sortBy) {
        case 'name':
          return (a.name || '').localeCompare(b.name || '')
        case 'calls':
          return (b.total_calls || 0) - (a.total_calls || 0)
        case 'rating':
          return (b.rating || 0) - (a.rating || 0)
        case 'created':
        default:
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      }
    })

    return result
  }, [agents, filterStatus, searchTerm, sortBy])

  const getStatusColor = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'active':
        return 'bg-green-100 text-green-800'
      case 'on-break':
        return 'bg-yellow-100 text-yellow-800'
      case 'inactive':
      case 'offline':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getStatusBgColor = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'active':
        return 'bg-green-50'
      case 'on-break':
        return 'bg-yellow-50'
      case 'inactive':
      case 'offline':
        return 'bg-red-50'
      default:
        return 'bg-gray-50'
    }
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
                  <h1 className="text-3xl font-bold text-gray-900">Agents</h1>
                  <p className="text-gray-600 mt-2">Manage and monitor your team</p>
                </div>
                <button className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg flex items-center gap-2 transition">
                  <FiPlus className="text-lg" />
                  Add Agent
                </button>
              </div>

              {/* Error Message */}
              {error && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
                  <p className="text-red-800 text-sm">{error}</p>
                </div>
              )}

              {/* Filters & Search */}
              <div className="bg-white rounded-lg shadow p-4 mb-6">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  {/* Search */}
                  <div className="relative">
                    <FiSearch className="absolute left-3 top-3 text-gray-400" />
                    <input
                      type="text"
                      placeholder="Search by name, email, phone..."
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
                    <option value="on-break">On Break</option>
                    <option value="inactive">Inactive</option>
                    <option value="offline">Offline</option>
                  </select>

                  {/* Sort */}
                  <select
                    value={sortBy}
                    onChange={(e) => setSortBy(e.target.value as any)}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="created">Newest First</option>
                    <option value="name">By Name</option>
                    <option value="calls">By Calls</option>
                    <option value="rating">By Rating</option>
                  </select>
                </div>

                {/* Results Count */}
                <div className="mt-4 text-sm text-gray-600">
                  Showing {filteredAgents.length} of {agents.length} agents
                </div>
              </div>

              {/* Loading State */}
              {loading && (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {[...Array(6)].map((_, i) => (
                    <div key={i} className="bg-white rounded-lg shadow h-64 animate-pulse"></div>
                  ))}
                </div>
              )}

              {/* Empty State */}
              {!loading && filteredAgents.length === 0 && (
                <div className="bg-white rounded-lg shadow p-12 text-center">
                  <div className="text-gray-500 mb-4 text-5xl">👥</div>
                  <p className="text-gray-600 text-lg font-medium">No agents found</p>
                  <p className="text-gray-500 mt-1">Try adjusting your filters or add a new agent</p>
                </div>
              )}

              {/* Agents Grid */}
              {!loading && filteredAgents.length > 0 && (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {filteredAgents.map((agent) => (
                    <div key={agent.id} className={`bg-white rounded-lg shadow hover:shadow-lg transition p-6 border-t-4 border-blue-600 ${getStatusBgColor(agent.status)}`}>
                      {/* Header */}
                      <div className="flex items-start justify-between mb-4">
                        <div className="flex-1">
                          <h3 className="text-lg font-semibold text-gray-900">{agent.name}</h3>
                          <p className="text-sm text-gray-600 mt-1">{agent.email}</p>
                        </div>
                        <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(agent.status)}`}>
                          {agent.status || 'Unknown'}
                        </span>
                      </div>

                      {/* Contact Info */}
                      {agent.phone && (
                        <div className="flex items-center gap-2 text-sm text-gray-600 mb-4">
                          <FiPhone className="text-gray-400" />
                          {agent.phone}
                        </div>
                      )}

                      {/* Stats */}
                      <div className="bg-gray-50 rounded-lg p-4 mb-4 space-y-3">
                        <div className="flex items-center justify-between">
                          <span className="text-sm text-gray-600 flex items-center gap-2">
                            <FiActivity className="text-blue-600" />
                            Total Calls
                          </span>
                          <span className="font-bold text-gray-900">{agent.total_calls || 0}</span>
                        </div>
                        <div className="flex items-center justify-between">
                          <span className="text-sm text-gray-600">Success Rate</span>
                          <span className="font-bold text-gray-900">
                            {agent.total_calls && agent.successful_calls
                              ? `${Math.round((agent.successful_calls / agent.total_calls) * 100)}%`
                              : '—'}
                          </span>
                        </div>
                        <div className="flex items-center justify-between">
                          <span className="text-sm text-gray-600 flex items-center gap-2">
                            <FiStar className="text-yellow-500" />
                            Rating
                          </span>
                          <span className="font-bold text-gray-900">{agent.rating?.toFixed(1) || '—'}</span>
                        </div>
                      </div>

                      {/* Footer */}
                      <div className="text-xs text-gray-500 pt-4 border-t">
                        Joined {format(new Date(agent.created_at), 'MMM d, yyyy')}
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
