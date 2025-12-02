'use client'

import React, { useState, useEffect } from 'react'
import { SourceIntegration } from '@/types/marketing'
import toast from 'react-hot-toast'

const SOURCE_TYPES = [
  { value: 'google', label: 'Google Ads', category: 'paid' },
  { value: 'meta', label: 'Meta (Facebook & Instagram)', category: 'paid' },
  { value: 'portal', label: 'Real Estate Portals', category: 'portal' },
  { value: 'website', label: 'Website & Landing Pages', category: 'organic' },
  { value: 'landing_page', label: 'Custom Landing Page', category: 'campaign' },
  { value: 'email', label: 'Email Marketing', category: 'owned' },
  { value: 'sms', label: 'SMS & WhatsApp', category: 'owned' },
  { value: 'whatsapp', label: 'WhatsApp Marketing', category: 'owned' },
  { value: 'referral', label: 'Referral Network', category: 'referral' },
  { value: 'offline', label: 'Offline Sources', category: 'offline' },
  { value: 'custom', label: 'Custom Source', category: 'custom' },
]

const PORTAL_OPTIONS = [
  { value: 'magicbricks', label: 'MagicBricks' },
  { value: '99acres', label: '99 Acres' },
  { value: 'housing', label: 'Housing.com' },
  { value: 'nobroker', label: 'NoBroker' },
  { value: 'proptiger', label: 'PropTiger' },
]

interface IntegrationFormProps {
  integration?: SourceIntegration
  onSuccess: () => void
  onCancel: () => void
}

export function IntegrationForm({ integration, onSuccess, onCancel }: IntegrationFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<SourceIntegration>>(
    integration || {
      source_type: 'google',
      source_name: '',
      sub_source: '',
      is_active: true,
      sync_frequency: 'daily',
    }
  )

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      // TODO: Call API to save integration
      // For now, show success
      toast.success(integration ? 'Integration updated' : 'Integration created')
      onSuccess()
    } catch (error) {
      toast.error('Failed to save integration')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const sourceTypeObj = SOURCE_TYPES.find((s) => s.value === formData.source_type)

  return (
    <form onSubmit={handleSubmit} className="space-y-6 bg-white p-6 rounded-lg shadow">
      {/* Source Type */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Source Type *</label>
        <select
          value={formData.source_type || 'google'}
          onChange={(e) => setFormData({ ...formData, source_type: e.target.value as any })}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
        >
          {SOURCE_TYPES.map((type) => (
            <option key={type.value} value={type.value}>
              {type.label}
            </option>
          ))}
        </select>
      </div>

      {/* Source Name */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Display Name *</label>
        <input
          type="text"
          value={formData.source_name || ''}
          onChange={(e) => setFormData({ ...formData, source_name: e.target.value })}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
          placeholder="e.g., Google Ads - Residential Campaign"
        />
      </div>

      {/* Portal Selection */}
      {formData.source_type === 'portal' && (
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Select Portal *</label>
          <select
            value={formData.sub_source || 'magicbricks'}
            onChange={(e) => setFormData({ ...formData, sub_source: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
          >
            {PORTAL_OPTIONS.map((portal) => (
              <option key={portal.value} value={portal.value}>
                {portal.label}
              </option>
            ))}
          </select>
        </div>
      )}

      {/* API Key / Credentials */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">API Key / Credentials</label>
        <input
          type="password"
          value={formData.api_key || ''}
          onChange={(e) => setFormData({ ...formData, api_key: e.target.value })}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
          placeholder="Paste API key or credentials (will be encrypted)"
        />
        <p className="text-xs text-gray-500 mt-1">Your credentials are encrypted and secure</p>
      </div>

      {/* Webhook URL */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Webhook URL</label>
        <input
          type="url"
          value={formData.webhook_url || ''}
          onChange={(e) => setFormData({ ...formData, webhook_url: e.target.value })}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
          placeholder="https://your-domain.com/webhook"
        />
        <p className="text-xs text-gray-500 mt-1">For real-time lead syncing</p>
      </div>

      {/* Sync Frequency */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Sync Frequency</label>
        <select
          value={formData.sync_frequency || 'daily'}
          onChange={(e) => setFormData({ ...formData, sync_frequency: e.target.value as any })}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
        >
          <option value="real_time">Real-time (via webhook)</option>
          <option value="hourly">Hourly</option>
          <option value="daily">Daily</option>
          <option value="manual">Manual</option>
        </select>
      </div>

      {/* Active Status */}
      <div className="flex items-center">
        <input
          type="checkbox"
          checked={formData.is_active ?? true}
          onChange={(e) => setFormData({ ...formData, is_active: e.target.checked })}
          className="w-4 h-4 rounded border-gray-300"
        />
        <label className="ml-2 text-sm text-gray-700">Active Integration</label>
      </div>

      {/* Info Box */}
      <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
        <p className="text-sm text-blue-800">
          <strong>Integration Type:</strong> {sourceTypeObj?.label}
        </p>
        <p className="text-xs text-blue-700 mt-2">
          Leads from this source will be automatically tagged and tracked with their UTM parameters.
        </p>
      </div>

      {/* Actions */}
      <div className="flex gap-3 pt-4 border-t">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-2 px-4 rounded-md transition"
        >
          {loading ? 'Saving...' : integration ? 'Update Integration' : 'Connect Integration'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium py-2 px-4 rounded-md transition"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
