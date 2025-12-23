'use client'

import { useState } from 'react'
import { ProtectedRoute } from '@/hooks/use-auth'
import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { FiCheck, FiX, FiRefreshCw, FiPlus, FiChevronRight, FiActivity, FiClock, FiAlertCircle } from 'react-icons/fi'
import { format } from 'date-fns'

interface Provider {
  id: string | number
  name: string
  type: string
  status: 'connected' | 'disconnected' | 'error'
  last_sync?: string
  next_sync?: string
  sync_frequency: string
  records_synced: number
  error_message?: string
}

interface SyncJob {
  id: string | number
  provider: string
  status: 'completed' | 'in-progress' | 'failed'
  records_count: number
  started_at: string
  completed_at?: string
  error?: string
}

// Mock data
const mockProviders: Provider[] = [
  {
    id: 1,
    name: 'Salesforce CRM',
    type: 'CRM',
    status: 'connected',
    last_sync: new Date(Date.now() - 3600000).toISOString(),
    next_sync: new Date(Date.now() + 86400000).toISOString(),
    sync_frequency: 'Daily',
    records_synced: 1250,
  },
  {
    id: 2,
    name: 'QuickBooks Online',
    type: 'Accounting',
    status: 'connected',
    last_sync: new Date(Date.now() - 7200000).toISOString(),
    next_sync: new Date(Date.now() + 43200000).toISOString(),
    sync_frequency: 'Twice Daily',
    records_synced: 890,
  },
  {
    id: 3,
    name: 'Google Workspace',
    type: 'Collaboration',
    status: 'error',
    last_sync: new Date(Date.now() - 259200000).toISOString(),
    sync_frequency: 'Hourly',
    records_synced: 450,
    error_message: 'Authentication token expired',
  },
  {
    id: 4,
    name: 'Zapier',
    type: 'Automation',
    status: 'disconnected',
    sync_frequency: 'On Demand',
    records_synced: 0,
  },
]

const mockSyncJobs: SyncJob[] = [
  {
    id: 1,
    provider: 'Salesforce CRM',
    status: 'completed',
    records_count: 245,
    started_at: new Date(Date.now() - 3600000).toISOString(),
    completed_at: new Date(Date.now() - 3540000).toISOString(),
  },
  {
    id: 2,
    provider: 'QuickBooks Online',
    status: 'in-progress',
    records_count: 120,
    started_at: new Date(Date.now() - 1200000).toISOString(),
  },
  {
    id: 3,
    provider: 'Google Workspace',
    status: 'failed',
    records_count: 0,
    started_at: new Date(Date.now() - 259200000).toISOString(),
    completed_at: new Date(Date.now() - 259140000).toISOString(),
    error: 'Authentication failed: Invalid token',
  },
]

export default function IntegrationPage() {
  const [providers, setProviders] = useState<Provider[]>(mockProviders)
  const [syncJobs] = useState<SyncJob[]>(mockSyncJobs)

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'connected':
        return 'text-green-600'
      case 'disconnected':
        return 'text-gray-600'
      case 'error':
        return 'text-red-600'
      case 'in-progress':
        return 'text-blue-600'
      case 'completed':
        return 'text-green-600'
      case 'failed':
        return 'text-red-600'
      default:
        return 'text-gray-600'
    }
  }

  const getStatusBadgeColor = (status: string) => {
    switch (status) {
      case 'connected':
        return 'bg-green-100 text-green-800'
      case 'disconnected':
        return 'bg-gray-100 text-gray-800'
      case 'error':
        return 'bg-red-100 text-red-800'
      case 'in-progress':
        return 'bg-blue-100 text-blue-800'
      case 'completed':
        return 'bg-green-100 text-green-800'
      case 'failed':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'CRM':
        return 'bg-blue-500'
      case 'Accounting':
        return 'bg-green-500'
      case 'Collaboration':
        return 'bg-purple-500'
      case 'Automation':
        return 'bg-orange-500'
      default:
        return 'bg-gray-500'
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
                  <h1 className="text-3xl font-bold text-gray-900">Integration Hub</h1>
                  <p className="text-gray-600 mt-2">Manage third-party integrations and data synchronization</p>
                </div>
                <button className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg flex items-center gap-2 transition">
                  <FiPlus className="text-lg" />
                  Add Integration
                </button>
              </div>

              {/* Tabs */}
              <div className="mb-6 border-b border-gray-200">
                <div className="flex gap-8">
                  <button className="pb-4 px-2 border-b-2 border-blue-600 text-blue-600 font-semibold">
                    <FiActivity className="inline mr-2" />
                    Providers
                  </button>
                  <button className="pb-4 px-2 text-gray-600 hover:text-gray-900 font-semibold">
                    <FiClock className="inline mr-2" />
                    Sync Jobs
                  </button>
                </div>
              </div>

              {/* Providers Grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-12">
                {providers.map((provider) => (
                  <div key={provider.id} className={`bg-white rounded-lg shadow hover:shadow-md transition p-6 border-l-4 ${provider.status === 'error' ? 'border-l-red-500' : provider.status === 'connected' ? 'border-l-green-500' : 'border-l-gray-300'}`}>
                    {/* Header */}
                    <div className="flex items-start justify-between mb-4">
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-2">
                          <span className={`px-3 py-1 rounded text-xs font-semibold text-white ${getTypeColor(provider.type)}`}>
                            {provider.type}
                          </span>
                          <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getStatusBadgeColor(provider.status)}`}>
                            {provider.status.charAt(0).toUpperCase() + provider.status.slice(1)}
                          </span>
                        </div>
                        <h3 className="text-lg font-semibold text-gray-900">{provider.name}</h3>
                      </div>
                      {provider.status === 'error' && (
                        <FiAlertCircle className="text-red-600 text-xl" />
                      )}
                    </div>

                    {/* Error Message */}
                    {provider.error_message && (
                      <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
                        <p className="text-sm text-red-700">{provider.error_message}</p>
                      </div>
                    )}

                    {/* Sync Info */}
                    {provider.status === 'connected' && (
                      <div className="mb-4 space-y-2 text-sm text-gray-600">
                        <div>
                          <span className="font-medium">Last Sync:</span> {provider.last_sync ? format(new Date(provider.last_sync), 'MMM d, yyyy HH:mm') : 'Never'}
                        </div>
                        <div>
                          <span className="font-medium">Next Sync:</span> {provider.next_sync ? format(new Date(provider.next_sync), 'MMM d, yyyy HH:mm') : '-'}
                        </div>
                        <div>
                          <span className="font-medium">Frequency:</span> {provider.sync_frequency}
                        </div>
                      </div>
                    )}

                    {/* Records Synced */}
                    <div className="mb-6 p-4 bg-gray-50 rounded-lg">
                      <p className="text-sm text-gray-600 mb-1">Records Synced</p>
                      <p className="text-2xl font-bold text-gray-900">{provider.records_synced.toLocaleString()}</p>
                    </div>

                    {/* Actions */}
                    <div className="flex gap-2">
                      {provider.status === 'connected' ? (
                        <>
                          <button className="flex-1 px-4 py-2 bg-blue-50 hover:bg-blue-100 text-blue-600 rounded-lg font-medium text-sm transition flex items-center justify-center gap-2">
                            <FiRefreshCw className="text-lg" />
                            Sync Now
                          </button>
                          <button className="flex-1 px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg font-medium text-sm transition">
                            Configure
                          </button>
                        </>
                      ) : provider.status === 'error' ? (
                        <>
                          <button className="flex-1 px-4 py-2 bg-red-50 hover:bg-red-100 text-red-600 rounded-lg font-medium text-sm transition">
                            Fix Issue
                          </button>
                          <button className="flex-1 px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg font-medium text-sm transition">
                            Retry
                          </button>
                        </>
                      ) : (
                        <button className="flex-1 px-4 py-2 bg-blue-50 hover:bg-blue-100 text-blue-600 rounded-lg font-medium text-sm transition">
                          Connect
                        </button>
                      )}
                    </div>
                  </div>
                ))}
              </div>

              {/* Sync Jobs Section */}
              <div>
                <h2 className="text-2xl font-bold text-gray-900 mb-6">Recent Sync Jobs</h2>
                <div className="bg-white rounded-lg shadow overflow-hidden">
                  <table className="w-full">
                    <thead className="bg-gray-50 border-b border-gray-200">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Provider</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Status</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Records</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Started</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Duration</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Details</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                      {syncJobs.map((job) => (
                        <tr key={job.id} className="hover:bg-gray-50 transition">
                          <td className="px-6 py-4 font-semibold text-gray-900">{job.provider}</td>
                          <td className="px-6 py-4">
                            <span className={`px-3 py-1 rounded-full text-xs font-semibold flex items-center gap-2 w-fit ${getStatusBadgeColor(job.status)}`}>
                              {job.status === 'completed' && <FiCheck className="text-lg" />}
                              {job.status === 'failed' && <FiX className="text-lg" />}
                              {job.status === 'in-progress' && <FiRefreshCw className="text-lg animate-spin" />}
                              {job.status.charAt(0).toUpperCase() + job.status.slice(1)}
                            </span>
                          </td>
                          <td className="px-6 py-4 text-gray-700">{job.records_count}</td>
                          <td className="px-6 py-4 text-sm text-gray-600">
                            {format(new Date(job.started_at), 'MMM d, yyyy HH:mm')}
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-600">
                            {job.completed_at
                              ? `${Math.round((new Date(job.completed_at).getTime() - new Date(job.started_at).getTime()) / 1000)}s`
                              : 'In progress...'}
                          </td>
                          <td className="px-6 py-4">
                            {job.error ? (
                              <div className="text-sm text-red-600">
                                <p className="font-medium">Error:</p>
                                <p>{job.error}</p>
                              </div>
                            ) : (
                              <button className="text-blue-600 hover:text-blue-800 font-medium text-sm flex items-center gap-1">
                                View <FiChevronRight />
                              </button>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </main>
        </div>
      </div>
    </ProtectedRoute>
  )
}
