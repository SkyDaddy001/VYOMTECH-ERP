'use client'

import React, { useState, useEffect } from 'react'
import { SourceIntegration } from '@/types/marketing'
import toast from 'react-hot-toast'

interface IntegrationListProps {
  projectId: string
  onEdit: (integration: SourceIntegration) => void
  onCreateNew: () => void
}

const SOURCE_ICONS: Record<string, string> = {
  google: '🔍',
  meta: 'f',
  portal: '🏠',
  website: '🌐',
  landing_page: '📄',
  email: '📧',
  sms: '💬',
  whatsapp: '💬',
  referral: '🔗',
  offline: '📞',
  custom: '⚙️',
}

const STATUS_COLORS: Record<string, string> = {
  configured: 'bg-green-100 text-green-800',
  syncing: 'bg-blue-100 text-blue-800',
  paused: 'bg-yellow-100 text-yellow-800',
  error: 'bg-red-100 text-red-800',
}

export function IntegrationList({ projectId, onEdit, onCreateNew }: IntegrationListProps) {
  const [integrations, setIntegrations] = useState<SourceIntegration[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<string>('all')

  useEffect(() => {
    loadIntegrations()
  }, [])

  const loadIntegrations = async () => {
    setLoading(true)
    try {
      // TODO: Call API to fetch integrations
      // For now, mock data
      setIntegrations([
        {
          id: '1',
          source_type: 'google',
          source_name: 'Google Ads - Residential',
          is_active: true,
          integration_status: 'configured',
          sync_frequency: 'hourly',
          lead_count: 245,
        },
        {
          id: '2',
          source_type: 'meta',
          source_name: 'Facebook & Instagram',
          is_active: true,
          integration_status: 'configured',
          sync_frequency: 'real_time',
          lead_count: 189,
        },
        {
          id: '3',
          source_type: 'portal',
          source_name: 'MagicBricks',
          sub_source: 'magicbricks',
          is_active: true,
          integration_status: 'configured',
          sync_frequency: 'daily',
          lead_count: 126,
        },
      ] as SourceIntegration[])
    } catch (error) {
      toast.error('Failed to load integrations')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleToggleActive = async (integration: SourceIntegration) => {
    try {
      // TODO: Call API to toggle active status
      toast.success(integration.is_active ? 'Integration paused' : 'Integration resumed')
      loadIntegrations()
    } catch (error) {
      toast.error('Failed to update integration')
      console.error(error)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure? This will stop syncing leads from this source.')) return

    try {
      // TODO: Call API to delete
      toast.success('Integration deleted')
      loadIntegrations()
    } catch (error) {
      toast.error('Failed to delete integration')
      console.error(error)
    }
  }

  const handleTestSync = async (id: string) => {
    try {
      // TODO: Call API to test sync
      toast.success('Test sync initiated. Check back in a few seconds.')
    } catch (error) {
      toast.error('Test sync failed')
      console.error(error)
    }
  }

  const filteredIntegrations =
    filter === 'all'
      ? integrations
      : integrations.filter((i) => (filter === 'active' ? i.is_active : !i.is_active))

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-gray-500">Loading integrations...</div>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {/* Header with filters */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex gap-2">
          <button
            onClick={() => setFilter('all')}
            className={`px-3 py-1 rounded text-sm ${
              filter === 'all'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            All ({integrations.length})
          </button>
          <button
            onClick={() => setFilter('active')}
            className={`px-3 py-1 rounded text-sm ${
              filter === 'active'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            Active ({integrations.filter((i) => i.is_active).length})
          </button>
        </div>

        <button
          onClick={onCreateNew}
          className="bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded-md transition"
        >
          + Add Integration
        </button>
      </div>

      {/* Integrations Grid */}
      {filteredIntegrations.length === 0 ? (
        <div className="bg-gray-50 p-8 rounded-lg text-center">
          <p className="text-gray-600 mb-4">No integrations configured yet</p>
          <button
            onClick={onCreateNew}
            className="text-blue-600 hover:text-blue-700 font-medium"
          >
            Set up your first integration →
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {filteredIntegrations.map((integration) => (
            <div key={integration.id} className="bg-white p-4 rounded-lg shadow hover:shadow-md transition border-l-4 border-blue-500">
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-3">
                  <span className="text-2xl">{SOURCE_ICONS[integration.source_type]}</span>
                  <div>
                    <h3 className="font-semibold text-gray-900">{integration.source_name}</h3>
                    <p className="text-xs text-gray-600">{integration.sub_source || integration.source_type}</p>
                  </div>
                </div>

                <div className="flex flex-col gap-2">
                  <span className={`px-2 py-1 rounded text-xs font-medium whitespace-nowrap ${STATUS_COLORS[integration.integration_status]}`}>
                    {integration.integration_status?.replace('_', ' ').toUpperCase()}
                  </span>
                </div>
              </div>

              {/* Stats */}
              <div className="space-y-2 mb-4 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">Leads Synced:</span>
                  <span className="font-semibold text-gray-900">{integration.lead_count || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Sync Frequency:</span>
                  <span className="font-semibold text-gray-900">{integration.sync_frequency?.replace('_', ' ')}</span>
                </div>
                {integration.last_sync && (
                  <div className="flex justify-between">
                    <span className="text-gray-600">Last Sync:</span>
                    <span className="font-semibold text-gray-900">
                      {new Date(integration.last_sync).toLocaleString()}
                    </span>
                  </div>
                )}
              </div>

              {/* Actions */}
              <div className="flex gap-2 pt-3 border-t">
                <button
                  onClick={() => handleTestSync(integration.id || '')}
                  className="flex-1 text-sm bg-gray-100 hover:bg-gray-200 text-gray-700 font-medium py-2 px-3 rounded transition"
                >
                  Test Sync
                </button>
                <button
                  onClick={() => onEdit(integration)}
                  className="flex-1 text-sm bg-blue-100 hover:bg-blue-200 text-blue-700 font-medium py-2 px-3 rounded transition"
                >
                  Edit
                </button>
                <button
                  onClick={() => handleToggleActive(integration)}
                  className={`flex-1 text-sm font-medium py-2 px-3 rounded transition ${
                    integration.is_active
                      ? 'bg-yellow-100 hover:bg-yellow-200 text-yellow-700'
                      : 'bg-green-100 hover:bg-green-200 text-green-700'
                  }`}
                >
                  {integration.is_active ? 'Pause' : 'Resume'}
                </button>
                <button
                  onClick={() => handleDelete(integration.id || '')}
                  className="flex-1 text-sm bg-red-100 hover:bg-red-200 text-red-700 font-medium py-2 px-3 rounded transition"
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Info Section */}
      <div className="bg-blue-50 p-4 rounded-lg border border-blue-200 mt-6">
        <h4 className="font-semibold text-blue-900 mb-2">💡 Supported Integrations</h4>
        <div className="text-sm text-blue-800 space-y-1">
          <p>• <strong>Google Ads:</strong> Real-time lead import & conversion tracking</p>
          <p>• <strong>Meta:</strong> Facebook & Instagram lead forms & pixel tracking</p>
          <p>• <strong>Portals:</strong> MagicBricks, 99 Acres, Housing.com, NoBroker, PropTiger</p>
          <p>• <strong>Website:</strong> Form submissions, chat messages, call tracking</p>
          <p>• <strong>Email & SMS:</strong> Campaign responses & click tracking</p>
          <p>• <strong>Referral:</strong> Broker network & customer referral tracking</p>
        </div>
      </div>
    </div>
  )
}
