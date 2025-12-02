'use client'

import { useState } from 'react'
import { Lead } from '@/types/presales'

interface LeadListProps {
  leads?: Lead[]
  onSelectLead?: (lead: Lead) => void
}

export default function LeadList({ leads = [], onSelectLead }: LeadListProps) {
  const [filterStatus, setFilterStatus] = useState<string>('')
  const [filterPriority, setFilterPriority] = useState<string>('')
  const [searchTerm, setSearchTerm] = useState('')

  const filteredLeads = leads.filter(lead => {
    const statusMatch = !filterStatus || lead.status === filterStatus
    const priorityMatch = !filterPriority || lead.priority === filterPriority
    const searchMatch = !searchTerm || 
      lead.customer_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      lead.customer_phone.includes(searchTerm)
    return statusMatch && priorityMatch && searchMatch
  })

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'new':
        return 'bg-blue-100 text-blue-800'
      case 'contacted':
        return 'bg-yellow-100 text-yellow-800'
      case 'qualified':
        return 'bg-green-100 text-green-800'
      case 'converted':
        return 'bg-purple-100 text-purple-800'
      case 'rejected':
        return 'bg-red-100 text-red-800'
      case 'inactive':
        return 'bg-gray-100 text-gray-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'urgent':
        return 'bg-red-100 text-red-800 border-red-300'
      case 'high':
        return 'bg-orange-100 text-orange-800 border-orange-300'
      case 'medium':
        return 'bg-yellow-100 text-yellow-800 border-yellow-300'
      case 'low':
        return 'bg-gray-100 text-gray-800 border-gray-300'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getPriorityBadge = (priority: string) => {
    const icons = {
      urgent: '🔴',
      high: '🟠',
      medium: '🟡',
      low: '🔵'
    }
    return icons[priority as keyof typeof icons] || '⚪'
  }

  const formatDate = (dateStr?: string) => {
    if (!dateStr) return 'N/A'
    const date = new Date(dateStr)
    return date.toLocaleDateString('en-IN', { month: 'short', day: 'numeric' })
  }

  const getSourceIcon = (source: string) => {
    const icons = {
      call: '☎️',
      website: '🌐',
      referral: '👥',
      walk_in: '🚶',
      social_media: '📱'
    }
    return icons[source as keyof typeof icons] || '📌'
  }

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6">
      <h2 className="text-2xl font-semibold text-gray-800 mb-6">Lead Management</h2>

      {/* Search and Filters */}
      <div className="mb-6 space-y-4">
        <div>
          <input
            type="text"
            placeholder="Search by name or phone..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="text-sm text-gray-600 block mb-2">Filter by Status</label>
            <select
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
            >
              <option value="">All Status</option>
              <option value="new">New</option>
              <option value="contacted">Contacted</option>
              <option value="qualified">Qualified</option>
              <option value="converted">Converted</option>
              <option value="rejected">Rejected</option>
              <option value="inactive">Inactive</option>
            </select>
          </div>

          <div>
            <label className="text-sm text-gray-600 block mb-2">Filter by Priority</label>
            <select
              value={filterPriority}
              onChange={(e) => setFilterPriority(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
            >
              <option value="">All Priority</option>
              <option value="urgent">🔴 Urgent</option>
              <option value="high">🟠 High</option>
              <option value="medium">🟡 Medium</option>
              <option value="low">🔵 Low</option>
            </select>
          </div>
        </div>
      </div>

      {/* Leads Table */}
      {filteredLeads.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-600 text-lg">No leads found</p>
          <p className="text-gray-500 text-sm mt-2">Try adjusting your filters or search criteria</p>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-300 bg-gray-50">
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Priority</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Customer</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Contact</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Source</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Status</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Last Contact</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Next Follow-up</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Action</th>
              </tr>
            </thead>
            <tbody>
              {filteredLeads.map(lead => (
                <tr key={lead.id} className="border-b border-gray-200 hover:bg-gray-50 transition-colors">
                  <td className="px-4 py-3 text-center text-xl">
                    {getPriorityBadge(lead.priority)}
                  </td>
                  <td className="px-4 py-3">
                    <div className="font-medium text-gray-900">{lead.customer_name}</div>
                    {lead.email && <div className="text-xs text-gray-600">{lead.email}</div>}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-700">{lead.customer_phone}</td>
                  <td className="px-4 py-3">
                    <span className="text-lg">{getSourceIcon(lead.source)}</span>
                    <span className="text-xs text-gray-600 ml-1">{lead.source}</span>
                  </td>
                  <td className="px-4 py-3">
                    <span className={`inline-block px-2 py-1 rounded text-xs font-medium ${getStatusColor(lead.status)}`}>
                      {lead.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm">{formatDate(lead.last_contact_date)}</td>
                  <td className="px-4 py-3 text-sm">
                    {lead.next_follow_up_date ? (
                      <span className="text-sm font-medium text-blue-600">
                        {formatDate(lead.next_follow_up_date)}
                      </span>
                    ) : (
                      <span className="text-gray-500">Not Set</span>
                    )}
                  </td>
                  <td className="px-4 py-3">
                    <button
                      onClick={() => onSelectLead?.(lead)}
                      className="text-blue-600 hover:text-blue-700 font-medium text-sm"
                    >
                      Details
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Summary Stats */}
      <div className="mt-6 grid grid-cols-5 gap-4 pt-6 border-t border-gray-200">
        <div className="bg-blue-50 rounded-lg p-4 text-center">
          <p className="text-sm text-blue-600 mb-1">Total Leads</p>
          <p className="text-2xl font-semibold text-blue-900">{leads.length}</p>
        </div>
        <div className="bg-yellow-50 rounded-lg p-4 text-center">
          <p className="text-sm text-yellow-600 mb-1">New/Contacted</p>
          <p className="text-2xl font-semibold text-yellow-900">
            {leads.filter(l => ['new', 'contacted'].includes(l.status)).length}
          </p>
        </div>
        <div className="bg-green-50 rounded-lg p-4 text-center">
          <p className="text-sm text-green-600 mb-1">Qualified</p>
          <p className="text-2xl font-semibold text-green-900">
            {leads.filter(l => l.status === 'qualified').length}
          </p>
        </div>
        <div className="bg-purple-50 rounded-lg p-4 text-center">
          <p className="text-sm text-purple-600 mb-1">Converted</p>
          <p className="text-2xl font-semibold text-purple-900">
            {leads.filter(l => l.status === 'converted').length}
          </p>
        </div>
        <div className="bg-red-50 rounded-lg p-4 text-center">
          <p className="text-sm text-red-600 mb-1">Urgent Priority</p>
          <p className="text-2xl font-semibold text-red-900">
            {leads.filter(l => l.priority === 'urgent').length}
          </p>
        </div>
      </div>
    </div>
  )
}
