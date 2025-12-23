'use client'

import { useState, useMemo } from 'react'
import { useLeads } from '@/hooks/use-api'
import { ProtectedRoute } from '@/hooks/use-auth'
import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { FiSearch, FiPlus, FiPhone, FiMail, FiTrendingUp, FiFilter } from 'react-icons/fi'
import { format } from 'date-fns'

interface Lead {
  id: string | number
  name: string
  email: string
  phone?: string
  company?: string
  status: string
  source?: string
  value?: number
  score?: number
  created_at: string
  updated_at: string
}

export default function LeadsPage() {
  const { data: leads, loading, error, refetch } = useLeads({ limit: 100 })
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState<string>('all')
  const [sortBy, setSortBy] = useState<'name' | 'value' | 'created' | 'score'>('created')

  // Filter and search
  const filteredLeads = useMemo(() => {
    let result = leads as Lead[]

    // Filter by status
    if (filterStatus !== 'all') {
      result = result.filter(lead => 
        (lead.status || '').toLowerCase().includes(filterStatus.toLowerCase())
      )
    }

    // Search by name, email, phone, or company
    if (searchTerm) {
      const query = searchTerm.toLowerCase()
      result = result.filter(lead =>
        (lead.name || '').toLowerCase().includes(query) ||
        (lead.email || '').toLowerCase().includes(query) ||
        (lead.phone || '').toLowerCase().includes(query) ||
        (lead.company || '').toLowerCase().includes(query)
      )
    }

    // Sort
    result.sort((a, b) => {
      switch (sortBy) {
        case 'name':
          return (a.name || '').localeCompare(b.name || '')
        case 'value':
          return (b.value || 0) - (a.value || 0)
        case 'score':
          return (b.score || 0) - (a.score || 0)
        case 'created':
        default:
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      }
    })

    return result
  }, [leads, filterStatus, searchTerm, sortBy])

  const getStatusColor = (status: string) => {
    const statusLower = (status || '').toLowerCase()
    if (statusLower.includes('done')) return 'bg-green-100 text-green-800'
    if (statusLower.includes('warm')) return 'bg-yellow-100 text-yellow-800'
    if (statusLower.includes('cold')) return 'bg-blue-100 text-blue-800'
    if (statusLower.includes('lost')) return 'bg-red-100 text-red-800'
    return 'bg-gray-100 text-gray-800'
  }

  const getStatusBgColor = (status: string) => {
    const statusLower = (status || '').toLowerCase()
    if (statusLower.includes('done')) return 'bg-green-50'
    if (statusLower.includes('warm')) return 'bg-yellow-50'
    if (statusLower.includes('cold')) return 'bg-blue-50'
    if (statusLower.includes('lost')) return 'bg-red-50'
    return 'bg-gray-50'
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
                  <h1 className="text-3xl font-bold text-gray-900">Leads</h1>
                  <p className="text-gray-600 mt-2">Manage and track all your leads</p>
                </div>
                <button className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg flex items-center gap-2 transition">
                  <FiPlus className="text-lg" />
                  Add Lead
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
                    <option value="done">Done</option>
                    <option value="warm">Warm</option>
                    <option value="cold">Cold</option>
                    <option value="lost">Lost</option>
                  </select>

                  {/* Sort */}
                  <select
                    value={sortBy}
                    onChange={(e) => setSortBy(e.target.value as any)}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="created">Newest First</option>
                    <option value="name">By Name</option>
                    <option value="value">By Value</option>
                    <option value="score">By Score</option>
                  </select>
                </div>

                {/* Results Count */}
                <div className="mt-4 text-sm text-gray-600">
                  Showing {filteredLeads.length} of {leads.length} leads
                </div>
              </div>

              {/* Loading State */}
              {loading && (
                <div className="space-y-4">
                  {[...Array(5)].map((_, i) => (
                    <div key={i} className="bg-white rounded-lg shadow h-20 animate-pulse"></div>
                  ))}
                </div>
              )}

              {/* Empty State */}
              {!loading && filteredLeads.length === 0 && (
                <div className="bg-white rounded-lg shadow p-12 text-center">
                  <div className="text-gray-500 mb-4 text-5xl">📭</div>
                  <p className="text-gray-600 text-lg font-medium">No leads found</p>
                  <p className="text-gray-500 mt-1">Try adjusting your filters or create a new lead</p>
                </div>
              )}

              {/* Leads Table */}
              {!loading && filteredLeads.length > 0 && (
                <div className="overflow-x-auto bg-white rounded-lg shadow">
                  <table className="w-full">
                    <thead className="bg-gray-50 border-b">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Name</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Contact</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Company</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Status</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Value</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Source</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">Created</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y">
                      {filteredLeads.map((lead) => (
                        <tr key={lead.id} className={`hover:bg-gray-50 transition ${getStatusBgColor(lead.status)}`}>
                          <td className="px-6 py-4">
                            <div className="font-medium text-gray-900">{lead.name}</div>
                          </td>
                          <td className="px-6 py-4">
                            <div className="space-y-1">
                              {lead.email && (
                                <div className="flex items-center gap-2 text-sm text-gray-600">
                                  <FiMail className="text-gray-400" />
                                  {lead.email}
                                </div>
                              )}
                              {lead.phone && (
                                <div className="flex items-center gap-2 text-sm text-gray-600">
                                  <FiPhone className="text-gray-400" />
                                  {lead.phone}
                                </div>
                              )}
                            </div>
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-600">{lead.company || '—'}</td>
                          <td className="px-6 py-4">
                            <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(lead.status)}`}>
                              {lead.status || 'Unknown'}
                            </span>
                          </td>
                          <td className="px-6 py-4 text-sm font-medium text-gray-900">
                            {lead.value ? `$${lead.value.toLocaleString()}` : '—'}
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-600">{lead.source || '—'}</td>
                          <td className="px-6 py-4 text-sm text-gray-600">
                            {format(new Date(lead.created_at), 'MMM d, yyyy')}
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
