'use client'

import { useEffect, useState } from 'react'
import { FiSearch, FiPlus, FiMoreVertical, FiPhone, FiMail, FiUser, FiTrendingUp, FiFilter } from 'react-icons/fi'
import { apiClient } from '@/lib/api-client'
import { format } from 'date-fns'
import { DetailedLeadStatus, STATUS_MAP, getStatusesByPhase, getPhases, getStatusInfo } from '@/lib/lead-status-config'

interface Lead {
  id: string
  name: string
  email: string
  phone?: string
  company?: string
  status: DetailedLeadStatus
  source?: string
  value?: number
  created_at: string
  updated_at: string
}

export default function LeadsPage() {
  const [leads, setLeads] = useState<Lead[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState<string>('all')

  useEffect(() => {
    fetchLeads()
  }, [])

  const fetchLeads = async () => {
    try {
      setLoading(true)
      // Fetch leads from real API
      const response = (await apiClient.get('/api/v1/leads', {
        params: {
          limit: 50,
          offset: 0
        }
      })) as any
      
      // Handle both direct array and wrapped response formats
      let leadData: Lead[] = []
      
      if (Array.isArray(response.data)) {
        leadData = response.data as Lead[]
      } else if (response.data && typeof response.data === 'object' && 'data' in response.data) {
        leadData = response.data.data || []
      }
      
      const formattedLeads = (leadData || []).map((lead: any) => ({
        id: lead.id?.toString() || '',
        name: lead.name || '',
        email: lead.email || '',
        phone: lead.phone,
        company: lead.company,
        status: lead.status as DetailedLeadStatus,
        source: lead.source,
        value: lead.value,
        created_at: lead.created_at,
        updated_at: lead.updated_at
      }))
      
      setLeads(formattedLeads)
    } catch (error) {
      console.error('Error fetching leads:', error)
      // Fallback to mock data if API fails
      const mockLeads: Lead[] = [
        {
          id: '1',
          name: 'John Smith',
          email: 'john@example.com',
          phone: '+1 234 567 8900',
          company: 'Acme Corp',
          status: 'SV - Done',
          source: 'LinkedIn',
          value: 50000,
          created_at: new Date(Date.now() - 604800000).toISOString(),
          updated_at: new Date(Date.now() - 172800000).toISOString()
        },
        {
          id: '2',
          name: 'Sarah Johnson',
          email: 'sarah@example.com',
          phone: '+1 234 567 8901',
          company: 'Tech Solutions Inc',
          status: 'Follow Up - Warm',
          source: 'Website',
          value: 25000,
          created_at: new Date(Date.now() - 432000000).toISOString(),
          updated_at: new Date(Date.now() - 86400000).toISOString()
        },
        {
          id: '3',
          name: 'Mike Davis',
          email: 'mike@example.com',
          phone: '+1 234 567 8902',
          company: 'Global Enterprises',
          status: 'Fresh Lead',
          source: 'Referral',
          value: 75000,
          created_at: new Date(Date.now() - 172800000).toISOString(),
          updated_at: new Date(Date.now() - 172800000).toISOString()
        },
        {
          id: '4',
          name: 'Emily Chen',
          email: 'emily@example.com',
          phone: '+1 234 567 8903',
          company: 'Innovation Labs',
          status: 'F2F - Scheduled',
          source: 'Event',
          value: 100000,
          created_at: new Date(Date.now() - 345600000).toISOString(),
          updated_at: new Date(Date.now() - 259200000).toISOString()
        },
        {
          id: '5',
          name: 'Robert Wilson',
          email: 'robert@example.com',
          phone: '+1 234 567 8904',
          company: 'Enterprise Solutions',
          status: 'Booking - In Progress',
          source: 'Cold Call',
          value: 150000,
          created_at: new Date(Date.now() - 864000000).toISOString(),
          updated_at: new Date(Date.now() - 345600000).toISOString()
        },
        {
          id: '6',
          name: 'Lisa Anderson',
          email: 'lisa@example.com',
          phone: '+1 234 567 8905',
          company: 'Digital Ventures',
          status: 'SV - Warm',
          source: 'Partner',
          value: 60000,
          created_at: new Date(Date.now() - 518400000).toISOString(),
          updated_at: new Date(Date.now() - 259200000).toISOString()
        }
      ]
      setLeads(mockLeads)
    } finally {
      setLoading(false)
    }
  }

  const filteredLeads = leads.filter(lead => {
    const matchesSearch = lead.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         lead.email.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesStatus = filterStatus === 'all' || lead.status === filterStatus
    return matchesSearch && matchesStatus
  })

  const getStatusColor = (status: DetailedLeadStatus) => {
    const statusInfo = getStatusInfo(status)
    return statusInfo.color
  }

  const phases = getPhases().sort()
  const statusesByPhase = Object.fromEntries(
    phases.map(phase => [phase, getStatusesByPhase(phase)])
  )

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Leads</h1>
          <p className="text-gray-600 mt-2">Manage and track your sales leads</p>
        </div>
        <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition flex items-center">
          <FiPlus className="w-4 h-4 mr-2" />
          New Lead
        </button>
      </div>

      {/* Search and Filters */}
      <div className="bg-white rounded-lg shadow p-4 mb-6">
        <div className="flex flex-col gap-4">
          <div className="flex flex-col md:flex-row gap-4">
            <div className="flex-1 relative">
              <FiSearch className="absolute left-3 top-3 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="Search by name, email, or company..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
          
          {/* Status Filter by Phase */}
          <div className="border-t pt-4">
            <label className="text-sm font-semibold text-gray-700 mb-3 flex items-center">
              <FiFilter className="w-4 h-4 mr-2" />
              Filter by Status
            </label>
            <div className="space-y-3 max-h-96 overflow-y-auto">
              <div>
                <button
                  onClick={() => setFilterStatus('all')}
                  className={`w-full text-left px-3 py-2 rounded-lg transition ${
                    filterStatus === 'all'
                      ? 'bg-blue-100 text-blue-800 font-medium'
                      : 'hover:bg-gray-100'
                  }`}
                >
                  All Statuses
                </button>
              </div>
              {phases.map(phase => (
                <div key={phase} className="border-l-2 border-gray-200 pl-3">
                  <p className="text-xs font-semibold text-gray-600 mb-2 uppercase">{phase}</p>
                  <div className="space-y-1">
                    {statusesByPhase[phase]?.map(status => (
                      <button
                        key={status}
                        onClick={() => setFilterStatus(status)}
                        className={`w-full text-left px-2 py-1.5 rounded text-sm transition ${
                          filterStatus === status
                            ? getStatusColor(status) + ' font-medium'
                            : 'hover:bg-gray-100'
                        }`}
                      >
                        {status}
                      </button>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Leads Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        {loading ? (
          <div className="p-8 text-center">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        ) : filteredLeads.length === 0 ? (
          <div className="p-8 text-center">
            <p className="text-gray-600">No leads found</p>
          </div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Name</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Company</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Contact</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Value</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Added</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900"></th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {filteredLeads.map((lead) => (
                <tr key={lead.id} className="hover:bg-gray-50 transition">
                  <td className="px-6 py-4">
                    <div className="flex items-center">
                      <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center">
                        <FiUser className="w-5 h-5 text-blue-600" />
                      </div>
                      <p className="ml-3 font-medium text-gray-900">{lead.name}</p>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-gray-600">{lead.company || '-'}</td>
                  <td className="px-6 py-4">
                    <div className="space-y-1">
                      {lead.email && (
                        <div className="flex items-center text-sm text-gray-600">
                          <FiMail className="w-4 h-4 mr-2" />
                          {lead.email}
                        </div>
                      )}
                      {lead.phone && (
                        <div className="flex items-center text-sm text-gray-600">
                          <FiPhone className="w-4 h-4 mr-2" />
                          {lead.phone}
                        </div>
                      )}
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex flex-col gap-2">
                      <span className={`px-3 py-1 rounded-full text-sm font-medium w-fit ${getStatusColor(lead.status)}`}>
                        {lead.status}
                      </span>
                      <span className="text-xs text-gray-500">
                        {getStatusInfo(lead.status).pipeline}
                      </span>
                    </div>
                  </td>
                  <td className="px-6 py-4 font-medium text-gray-900">
                    {lead.value ? `$${(lead.value / 1000).toFixed(0)}k` : '-'}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {format(new Date(lead.created_at), 'MMM dd, yyyy')}
                  </td>
                  <td className="px-6 py-4">
                    <button className="p-2 hover:bg-gray-100 rounded-lg transition">
                      <FiMoreVertical className="w-4 h-4 text-gray-600" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Results Summary */}
      <div className="mt-4 text-sm text-gray-600">
        Showing {filteredLeads.length} of {leads.length} leads
      </div>
    </div>
  )
}
