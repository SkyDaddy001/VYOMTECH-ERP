'use client'

import { useState, useMemo } from 'react'
import { useCalls } from '@/hooks/use-api'
import { ProtectedRoute } from '@/hooks/use-auth'
import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { FiSearch, FiPhone, FiClock, FiDownload, FiVolume2 } from 'react-icons/fi'
import { format, formatDistance } from 'date-fns'

interface Call {
  id: string | number
  lead_name?: string
  agent_name?: string
  phone?: string
  duration?: number
  status: 'completed' | 'missed' | 'voicemail' | 'ongoing'
  outcome?: 'positive' | 'neutral' | 'negative'
  notes?: string
  recording_url?: string
  created_at: string
  updated_at: string
}

export default function CallsPage() {
  const { data: calls, loading, error, refetch } = useCalls({ limit: 100 })
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState<string>('all')
  const [sortBy, setSortBy] = useState<'newest' | 'duration' | 'oldest'>('newest')

  // Filter and search
  const filteredCalls = useMemo(() => {
    let result = calls as Call[]

    // Filter by status
    if (filterStatus !== 'all') {
      result = result.filter(call => 
        (call.status || '').toLowerCase() === filterStatus.toLowerCase()
      )
    }

    // Search by name, phone, or agent
    if (searchTerm) {
      const query = searchTerm.toLowerCase()
      result = result.filter(call =>
        (call.lead_name || '').toLowerCase().includes(query) ||
        (call.phone || '').toLowerCase().includes(query) ||
        (call.agent_name || '').toLowerCase().includes(query)
      )
    }

    // Sort
    result.sort((a, b) => {
      switch (sortBy) {
        case 'duration':
          return (b.duration || 0) - (a.duration || 0)
        case 'oldest':
          return new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
        case 'newest':
        default:
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      }
    })

    return result
  }, [calls, filterStatus, searchTerm, sortBy])

  const getStatusColor = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'completed':
        return 'bg-green-100 text-green-800'
      case 'missed':
        return 'bg-red-100 text-red-800'
      case 'voicemail':
        return 'bg-blue-100 text-blue-800'
      case 'ongoing':
        return 'bg-yellow-100 text-yellow-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getOutcomeColor = (outcome?: string) => {
    switch (outcome?.toLowerCase()) {
      case 'positive':
        return 'text-green-600'
      case 'negative':
        return 'text-red-600'
      case 'neutral':
        return 'text-gray-600'
      default:
        return 'text-gray-400'
    }
  }

  const formatDuration = (seconds?: number) => {
    if (!seconds) return '—'
    const minutes = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${minutes}m ${secs}s`
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
              <div className="mb-8">
                <h1 className="text-3xl font-bold text-gray-900">Calls</h1>
                <p className="text-gray-600 mt-2">Monitor and review all call activity</p>
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
                      placeholder="Search by name, phone, or agent..."
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
                    <option value="completed">Completed</option>
                    <option value="missed">Missed</option>
                    <option value="voicemail">Voicemail</option>
                    <option value="ongoing">Ongoing</option>
                  </select>

                  {/* Sort */}
                  <select
                    value={sortBy}
                    onChange={(e) => setSortBy(e.target.value as any)}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="newest">Newest First</option>
                    <option value="oldest">Oldest First</option>
                    <option value="duration">Longest Duration</option>
                  </select>
                </div>

                {/* Results Count */}
                <div className="mt-4 text-sm text-gray-600">
                  Showing {filteredCalls.length} of {calls.length} calls
                </div>
              </div>

              {/* Loading State */}
              {loading && (
                <div className="space-y-4">
                  {[...Array(8)].map((_, i) => (
                    <div key={i} className="bg-white rounded-lg shadow h-20 animate-pulse"></div>
                  ))}
                </div>
              )}

              {/* Empty State */}
              {!loading && filteredCalls.length === 0 && (
                <div className="bg-white rounded-lg shadow p-12 text-center">
                  <div className="text-gray-500 mb-4 text-5xl">📞</div>
                  <p className="text-gray-600 text-lg font-medium">No calls found</p>
                  <p className="text-gray-500 mt-1">Try adjusting your filters</p>
                </div>
              )}

              {/* Calls List */}
              {!loading && filteredCalls.length > 0 && (
                <div className="space-y-3">
                  {filteredCalls.map((call) => (
                    <div key={call.id} className="bg-white rounded-lg shadow hover:shadow-md transition p-4 border-l-4 border-blue-600">
                      <div className="flex items-start justify-between mb-3">
                        {/* Left - Info */}
                        <div className="flex-1">
                          <div className="flex items-center gap-3 mb-2">
                            <div className="p-2 bg-blue-100 rounded-lg">
                              <FiPhone className="text-blue-600" />
                            </div>
                            <div>
                              <h3 className="font-semibold text-gray-900">{call.lead_name || 'Unknown Lead'}</h3>
                              <p className="text-sm text-gray-600">{call.phone || '—'}</p>
                            </div>
                          </div>
                          <p className="text-sm text-gray-600 mt-2">
                            Agent: <span className="font-medium">{call.agent_name || 'Unknown'}</span>
                          </p>
                        </div>

                        {/* Right - Status & Duration */}
                        <div className="text-right">
                          <span className={`inline-block px-3 py-1 rounded-full text-xs font-semibold mb-2 ${getStatusColor(call.status)}`}>
                            {call.status || 'Unknown'}
                          </span>
                          <div className="text-sm text-gray-600 flex items-center justify-end gap-1">
                            <FiClock className="text-gray-400" />
                            {formatDuration(call.duration)}
                          </div>
                        </div>
                      </div>

                      {/* Bottom - Meta & Actions */}
                      <div className="flex items-center justify-between pt-3 border-t border-gray-100">
                        <div className="flex items-center gap-4 text-xs text-gray-500">
                          <span>{format(new Date(call.created_at), 'MMM d, yyyy HH:mm')}</span>
                          {call.outcome && (
                            <span className={`font-medium ${getOutcomeColor(call.outcome)}`}>
                              {call.outcome.charAt(0).toUpperCase() + call.outcome.slice(1)}
                            </span>
                          )}
                        </div>

                        {/* Actions */}
                        <div className="flex items-center gap-2">
                          {call.recording_url && (
                            <button className="p-2 hover:bg-gray-100 rounded-lg transition text-blue-600 hover:text-blue-700 flex items-center gap-1">
                              <FiVolume2 className="text-lg" />
                            </button>
                          )}
                          <button className="p-2 hover:bg-gray-100 rounded-lg transition text-gray-600 hover:text-gray-700">
                            <FiDownload className="text-lg" />
                          </button>
                        </div>
                      </div>

                      {/* Notes */}
                      {call.notes && (
                        <div className="mt-3 p-3 bg-gray-50 rounded-lg text-sm text-gray-700 border-l-2 border-yellow-400">
                          <span className="font-medium text-gray-900">Notes:</span> {call.notes}
                        </div>
                      )}
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
