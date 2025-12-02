'use client'

import React, { useState } from 'react'
import { UTMTracker, SourceIntegration } from '@/types/marketing'
import { marketingService } from '@/services/marketing.service'
import toast from 'react-hot-toast'

interface UTMGeneratorProps {
  campaignId: string
  campaignName: string
  integrations: SourceIntegration[]
  onSuccess: (utm: UTMTracker) => void
}

const SOURCE_PRESETS: Record<string, { medium: string; defaults: Record<string, string> }> = {
  google: { medium: 'cpc', defaults: { term: 'property keywords' } },
  facebook: { medium: 'social', defaults: { content: 'feed_ad' } },
  instagram: { medium: 'social', defaults: { content: 'story_ad' } },
  email: { medium: 'email', defaults: { content: 'newsletter' } },
  sms: { medium: 'sms', defaults: { content: 'campaign' } },
  whatsapp: { medium: 'sms', defaults: { content: 'broadcast' } },
  magicbricks: { medium: 'cpc', defaults: {} },
  '99acres': { medium: 'cpc', defaults: {} },
  housing: { medium: 'cpc', defaults: {} },
  nobroker: { medium: 'cpc', defaults: {} },
  proptiger: { medium: 'cpc', defaults: {} },
  website: { medium: 'organic', defaults: {} },
  website_form: { medium: 'form', defaults: { content: 'contact_form' } },
  landing_page: { medium: 'form', defaults: { content: 'lp_form' } },
  referral: { medium: 'referral', defaults: {} },
  offline: { medium: 'offline', defaults: { content: 'event' } },
}

export function UTMGenerator({ campaignId, campaignName, integrations, onSuccess }: UTMGeneratorProps) {
  const [loading, setLoading] = useState(false)
  const [step, setStep] = useState(1)
  const [selectedSource, setSelectedSource] = useState<string>('')
  const [selectedIntegration, setSelectedIntegration] = useState<SourceIntegration | null>(null)
  const [formData, setFormData] = useState({
    utm_source: '',
    utm_medium: '',
    utm_campaign: campaignName.toLowerCase().replace(/\s+/g, '_'),
    utm_content: '',
    utm_term: '',
    full_url: '',
  })
  const [generatedTrackers, setGeneratedTrackers] = useState<UTMTracker[]>([])

  const handleSourceSelect = (source: string, integration?: SourceIntegration) => {
    setSelectedSource(source)
    setSelectedIntegration(integration || null)

    const preset = SOURCE_PRESETS[source]
    if (preset) {
      setFormData((prev) => ({
        ...prev,
        utm_source: source,
        utm_medium: preset.medium,
        utm_content: preset.defaults.content || '',
        utm_term: preset.defaults.term || '',
      }))
    }
    setStep(2)
  }

  const handleGenerateURL = (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.full_url) {
      toast.error('Please enter a base URL')
      return
    }

    const trackingURL = marketingService.buildTrackingURL(formData.full_url, formData)
    setFormData((prev) => ({ ...prev, full_url: trackingURL }))
    setStep(3)
  }

  const handleCreateTracker = async () => {
    setLoading(true)
    try {
      const tracker = await marketingService.createUTMTracker({
        campaign_id: campaignId,
        integration_id: selectedIntegration?.id,
        utm_source: formData.utm_source,
        utm_medium: formData.utm_medium,
        utm_campaign: formData.utm_campaign,
        utm_content: formData.utm_content,
        utm_term: formData.utm_term,
        full_url: formData.full_url,
        source_type: (selectedSource as any) || 'organic',
        portal_name: selectedIntegration?.source_type === 'portal' ? selectedIntegration?.sub_source : undefined,
      })

      toast.success('UTM tracker created successfully')
      setGeneratedTrackers([...generatedTrackers, tracker])
      onSuccess(tracker)

      // Reset form
      setStep(1)
      setSelectedSource('')
      setSelectedIntegration(null)
      setFormData({
        utm_source: '',
        utm_medium: '',
        utm_campaign: campaignName.toLowerCase().replace(/\s+/g, '_'),
        utm_content: '',
        utm_term: '',
        full_url: '',
      })
    } catch (error) {
      toast.error('Failed to create UTM tracker')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleCopyURL = () => {
    navigator.clipboard.writeText(formData.full_url)
    toast.success('URL copied to clipboard!')
  }

  return (
    <div className="space-y-6">
      {/* Step 1: Select Source */}
      {step === 1 && (
        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-gray-900">Select Source Type</h3>

          {/* Paid Channels */}
          <div>
            <h4 className="text-sm font-medium text-gray-700 mb-2">🔥 Paid Advertising</h4>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
              <button
                onClick={() => handleSourceSelect('google')}
                className="p-3 border-2 border-gray-200 hover:border-blue-500 rounded-lg text-center transition"
              >
                <span className="text-2xl">🔍</span>
                <p className="text-sm font-medium mt-1">Google Ads</p>
              </button>
              <button
                onClick={() => handleSourceSelect('facebook')}
                className="p-3 border-2 border-gray-200 hover:border-blue-500 rounded-lg text-center transition"
              >
                <span className="text-2xl">f</span>
                <p className="text-sm font-medium mt-1">Facebook</p>
              </button>
              <button
                onClick={() => handleSourceSelect('instagram')}
                className="p-3 border-2 border-gray-200 hover:border-blue-500 rounded-lg text-center transition"
              >
                <span className="text-2xl">📷</span>
                <p className="text-sm font-medium mt-1">Instagram</p>
              </button>
            </div>
          </div>

          {/* Portal Integrations */}
          {integrations.filter((i) => i.source_type === 'portal').length > 0 && (
            <div>
              <h4 className="text-sm font-medium text-gray-700 mb-2">🏠 Real Estate Portals</h4>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
                {integrations
                  .filter((i) => i.source_type === 'portal')
                  .map((integration) => (
                    <button
                      key={integration.id}
                      onClick={() => handleSourceSelect(integration.sub_source || 'portal', integration)}
                      className="p-3 border-2 border-gray-200 hover:border-green-500 rounded-lg text-center transition"
                    >
                      <span className="text-2xl">📍</span>
                      <p className="text-sm font-medium mt-1">{integration.source_name}</p>
                    </button>
                  ))}
              </div>
            </div>
          )}

          {/* Website & Landing Pages */}
          <div>
            <h4 className="text-sm font-medium text-gray-700 mb-2">🌐 Website & Landing Pages</h4>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
              <button
                onClick={() => handleSourceSelect('website')}
                className="p-3 border-2 border-gray-200 hover:border-purple-500 rounded-lg text-center transition"
              >
                <span className="text-2xl">🌐</span>
                <p className="text-sm font-medium mt-1">Website</p>
              </button>
              <button
                onClick={() => handleSourceSelect('website_form')}
                className="p-3 border-2 border-gray-200 hover:border-purple-500 rounded-lg text-center transition"
              >
                <span className="text-2xl">📝</span>
                <p className="text-sm font-medium mt-1">Web Form</p>
              </button>
              <button
                onClick={() => handleSourceSelect('landing_page')}
                className="p-3 border-2 border-gray-200 hover:border-purple-500 rounded-lg text-center transition"
              >
                <span className="text-2xl">📄</span>
                <p className="text-sm font-medium mt-1">Landing Page</p>
              </button>
            </div>
          </div>

          {/* Owned Channels */}
          <div>
            <h4 className="text-sm font-medium text-gray-700 mb-2">📧 Owned Channels</h4>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
              <button
                onClick={() => handleSourceSelect('email')}
                className="p-3 border-2 border-gray-200 hover:border-orange-500 rounded-lg text-center transition"
              >
                <span className="text-2xl">📧</span>
                <p className="text-sm font-medium mt-1">Email</p>
              </button>
              <button
                onClick={() => handleSourceSelect('sms')}
                className="p-3 border-2 border-gray-200 hover:border-orange-500 rounded-lg text-center transition"
              >
                <span className="text-2xl">💬</span>
                <p className="text-sm font-medium mt-1">SMS</p>
              </button>
              <button
                onClick={() => handleSourceSelect('whatsapp')}
                className="p-3 border-2 border-gray-200 hover:border-orange-500 rounded-lg text-center transition"
              >
                <span className="text-2xl">💬</span>
                <p className="text-sm font-medium mt-1">WhatsApp</p>
              </button>
            </div>
          </div>

          {/* Other Sources */}
          <div>
            <h4 className="text-sm font-medium text-gray-700 mb-2">📞 Other Sources</h4>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
              <button
                onClick={() => handleSourceSelect('referral')}
                className="p-3 border-2 border-gray-200 hover:border-red-500 rounded-lg text-center transition"
              >
                <span className="text-2xl">🔗</span>
                <p className="text-sm font-medium mt-1">Referral</p>
              </button>
              <button
                onClick={() => handleSourceSelect('offline')}
                className="p-3 border-2 border-gray-200 hover:border-red-500 rounded-lg text-center transition"
              >
                <span className="text-2xl">📞</span>
                <p className="text-sm font-medium mt-1">Offline</p>
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Step 2: Customize URL */}
      {step === 2 && (
        <form onSubmit={handleGenerateURL} className="space-y-4 bg-gray-50 p-4 rounded-lg">
          <h3 className="text-lg font-semibold text-gray-900">Customize Tracking Parameters</h3>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">UTM Source *</label>
              <input
                type="text"
                value={formData.utm_source}
                onChange={(e) => setFormData({ ...formData, utm_source: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 text-sm"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">UTM Medium *</label>
              <input
                type="text"
                value={formData.utm_medium}
                onChange={(e) => setFormData({ ...formData, utm_medium: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 text-sm"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">UTM Campaign *</label>
              <input
                type="text"
                value={formData.utm_campaign}
                onChange={(e) => setFormData({ ...formData, utm_campaign: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 text-sm"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">UTM Content</label>
              <input
                type="text"
                value={formData.utm_content}
                onChange={(e) => setFormData({ ...formData, utm_content: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 text-sm"
                placeholder="Ad variant, email subject, etc."
              />
            </div>

            <div className="col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">UTM Term</label>
              <input
                type="text"
                value={formData.utm_term}
                onChange={(e) => setFormData({ ...formData, utm_term: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 text-sm"
                placeholder="Keywords, search terms, etc."
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Base URL (without UTM) *</label>
            <input
              type="url"
              value={formData.full_url}
              onChange={(e) => setFormData({ ...formData, full_url: e.target.value })}
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
              placeholder="https://yoursite.com/property/project-name"
            />
          </div>

          <div className="flex gap-2 pt-2">
            <button
              type="button"
              onClick={() => setStep(1)}
              className="flex-1 bg-gray-300 hover:bg-gray-400 text-gray-800 font-medium py-2 px-4 rounded-md transition"
            >
              Back
            </button>
            <button
              type="submit"
              className="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition"
            >
              Generate URL
            </button>
          </div>
        </form>
      )}

      {/* Step 3: Review & Save */}
      {step === 3 && (
        <div className="space-y-4 bg-green-50 p-4 rounded-lg border border-green-200">
          <h3 className="text-lg font-semibold text-gray-900">✓ Tracking URL Ready</h3>

          <div className="bg-white p-4 rounded border border-gray-300">
            <p className="text-xs text-gray-600 mb-2">Full Tracking URL:</p>
            <p className="text-sm font-mono text-gray-900 break-all">{formData.full_url}</p>
            <button
              onClick={handleCopyURL}
              className="mt-2 text-sm text-blue-600 hover:text-blue-700 font-medium"
            >
              📋 Copy URL
            </button>
          </div>

          <div className="flex gap-2 pt-2">
            <button
              type="button"
              onClick={() => setStep(2)}
              className="flex-1 bg-gray-300 hover:bg-gray-400 text-gray-800 font-medium py-2 px-4 rounded-md transition"
            >
              Edit
            </button>
            <button
              onClick={handleCreateTracker}
              disabled={loading}
              className="flex-1 bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white font-medium py-2 px-4 rounded-md transition"
            >
              {loading ? 'Saving...' : 'Save & Create Another'}
            </button>
          </div>
        </div>
      )}

      {/* Generated Trackers */}
      {generatedTrackers.length > 0 && (
        <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
          <h4 className="font-semibold text-blue-900 mb-2">Generated Trackers ({generatedTrackers.length})</h4>
          <div className="space-y-2 text-sm">
            {generatedTrackers.map((tracker, idx) => (
              <div key={tracker.id} className="bg-white p-2 rounded">
                <p className="font-medium text-gray-900">
                  {idx + 1}. {tracker.utm_source} - {tracker.utm_medium}
                </p>
                <p className="text-xs text-gray-600 text-gray-600 truncate">{tracker.full_url}</p>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
