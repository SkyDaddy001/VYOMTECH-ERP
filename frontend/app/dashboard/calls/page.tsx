'use client'

import { useEffect, useState } from 'react'
import { FiSearch, FiPlus, FiPhone, FiClock, FiUser, FiDownload } from 'react-icons/fi'
import { format, formatDistance } from 'date-fns'

interface Call {
  id: string
  leadName: string
  leadPhone: string
  agentName: string
  duration: number
  status: 'completed' | 'missed' | 'voicemail'
  outcome: 'positive' | 'neutral' | 'negative'
  notes?: string
  recordingUrl?: string
  timestamp: string
}

export default function CallsPage() {
  const [calls, setCalls] = useState<Call[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState<string>('all')

  useEffect(() => {
    fetchCalls()
  }, [])

  const fetchCalls = async () => {
    try {
      setLoading(true)
      // Mock data
      const mockCalls: Call[] = [
        {
          id: '1',
          leadName: 'John Smith',
          leadPhone: '+1 234 567 8900',
          agentName: 'Alice Johnson',
          duration: 540,
          status: 'completed',
          outcome: 'positive',
          notes: 'Lead interested in enterprise plan',
          recordingUrl: '#',
          timestamp: new Date(Date.now() - 3600000).toISOString()
        },
        {
          id: '2',
          leadName: 'Sarah Davis',
          leadPhone: '+1 234 567 8901',
          agentName: 'Bob Wilson',
          duration: 180,
          status: 'completed',
          outcome: 'neutral',
          notes: 'Requested callback next week',
          recordingUrl: '#',
          timestamp: new Date(Date.now() - 7200000).toISOString()
        },
        {
          id: '3',
          leadName: 'Mike Brown',
          leadPhone: '+1 234 567 8902',
          agentName: 'Carol Martinez',
          duration: 0,
          status: 'missed',
          outcome: 'negative',
          recordingUrl: '#',
          timestamp: new Date(Date.now() - 10800000).toISOString()
        },
        {
          id: '4',
          leadName: 'Emma Wilson',
          leadPhone: '+1 234 567 8903',
          agentName: 'Alice Johnson',
          duration: 900,
          status: 'completed',
          outcome: 'positive',
          notes: 'Demo scheduled for tomorrow',
          recordingUrl: '#',
          timestamp: new Date(Date.now() - 86400000).toISOString()
        }
      ]
      setCalls(mockCalls)
    } catch (error) {
      console.error('Error fetching calls:', error)
    } finally {
      setLoading(false)
    }
  }

  const filteredCalls = calls.filter(call => {
    const matchesSearch = call.leadName.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         call.leadPhone.includes(searchTerm)
    const matchesStatus = filterStatus === 'all' || call.status === filterStatus
    return matchesSearch && matchesStatus
  })

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      completed: 'bg-green-100 text-green-800',
      missed: 'bg-red-100 text-red-800',
      voicemail: 'bg-yellow-100 text-yellow-800'
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const getOutcomeColor = (outcome: string) => {
    const colors: Record<string, string> = {
      positive: 'text-green-600',
      neutral: 'text-yellow-600',
      negative: 'text-red-600'
    }
    return colors[outcome] || 'text-gray-600'
  }

  const formatDuration = (seconds: number) => {
    if (seconds === 0) return 'Missed'
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins}:${secs.toString().padStart(2, '0')}`
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Calls</h1>
          <p className="text-gray-600 mt-2">View and manage all call records</p>
        </div>
        <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition flex items-center">
          <FiPlus className="w-4 h-4 mr-2" />
          Make Call
        </button>
      </div>

      {/* Search and Filters */}
      <div className="bg-white rounded-lg shadow p-4 mb-6">
        <div className="flex flex-col md:flex-row gap-4">
          <div className="flex-1 relative">
            <FiSearch className="absolute left-3 top-3 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search by name or phone..."
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
            <option value="completed">Completed</option>
            <option value="missed">Missed</option>
            <option value="voicemail">Voicemail</option>
          </select>
        </div>
      </div>

      {/* Calls Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        {loading ? (
          <div className="p-8 text-center">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        ) : filteredCalls.length === 0 ? (
          <div className="p-8 text-center">
            <p className="text-gray-600">No calls found</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Lead</th>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Agent</th>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Duration</th>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Outcome</th>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Time</th>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Recording</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {filteredCalls.map((call) => (
                  <tr key={call.id} className="hover:bg-gray-50 transition">
                    <td className="px-6 py-4">
                      <div>
                        <p className="font-medium text-gray-900">{call.leadName}</p>
                        <p className="text-sm text-gray-600 flex items-center mt-1">
                          <FiPhone className="w-3 h-3 mr-1" />
                          {call.leadPhone}
                        </p>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center">
                        <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center mr-2">
                          <FiUser className="w-4 h-4 text-blue-600" />
                        </div>
                        <p className="text-sm text-gray-900">{call.agentName}</p>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center text-sm text-gray-600">
                        <FiClock className="w-4 h-4 mr-2" />
                        {formatDuration(call.duration)}
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(call.status)}`}>
                        {call.status.charAt(0).toUpperCase() + call.status.slice(1)}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <div className={`flex items-center ${getOutcomeColor(call.outcome)}`}>
                        <div className="w-2 h-2 rounded-full mr-2 bg-current"></div>
                        <p className="text-sm font-medium">{call.outcome.charAt(0).toUpperCase() + call.outcome.slice(1)}</p>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {formatDistance(new Date(call.timestamp), new Date(), { addSuffix: true })}
                    </td>
                    <td className="px-6 py-4">
                      <button className="p-2 hover:bg-gray-100 rounded-lg transition">
                        <FiDownload className="w-4 h-4 text-gray-600" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Summary Stats */}
      {!loading && filteredCalls.length > 0 && (
        <div className="mt-8 grid grid-cols-1 md:grid-cols-4 gap-6">
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-sm text-gray-600">Total Calls</p>
            <p className="text-3xl font-bold text-gray-900 mt-2">{filteredCalls.length}</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-sm text-gray-600">Completed</p>
            <p className="text-3xl font-bold text-green-600 mt-2">
              {filteredCalls.filter(c => c.status === 'completed').length}
            </p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-sm text-gray-600">Avg Duration</p>
            <p className="text-3xl font-bold text-gray-900 mt-2">
              {formatDuration(Math.round(filteredCalls.reduce((sum, c) => sum + c.duration, 0) / filteredCalls.length))}
            </p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-sm text-gray-600">Positive Outcome</p>
            <p className="text-3xl font-bold text-green-600 mt-2">
              {Math.round((filteredCalls.filter(c => c.outcome === 'positive').length / filteredCalls.length) * 100)}%
            </p>
          </div>
        </div>
      )}
    </div>
  )
}
