'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import { Campaign } from '@/types/marketing'

interface CampaignFormProps {
  campaign?: Campaign | null
  onSubmit: (data: Partial<Campaign>) => Promise<void>
  onCancel: () => void
}

export default function CampaignForm({ campaign, onSubmit, onCancel }: CampaignFormProps) {
  const [formData, setFormData] = useState<Partial<Campaign>>({
    name: '',
    description: '',
    campaign_type: 'project_launch',
    status: 'draft',
    start_date: new Date().toISOString().split('T')[0],
    end_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    budget: 0,
    target_segment: 'residential',
    target_price_range_min: 0,
    target_price_range_max: 0,
    expected_leads: 0,
    expected_bookings: 0,
  })
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (campaign) {
      setFormData({
        id: campaign.id,
        name: campaign.name,
        description: campaign.description,
        campaign_type: campaign.campaign_type,
        status: campaign.status,
        start_date: campaign.start_date?.split('T')[0] || '',
        end_date: campaign.end_date?.split('T')[0] || '',
        budget: campaign.budget,
        target_segment: campaign.target_segment,
        target_price_range_min: campaign.target_price_range_min,
        target_price_range_max: campaign.target_price_range_max,
        expected_leads: campaign.expected_leads,
        expected_bookings: campaign.expected_bookings,
      })
    }
  }, [campaign])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    const numericFields = ['budget', 'expected_leads', 'expected_bookings', 'target_price_range_min', 'target_price_range_max']
    setFormData((prev) => ({
      ...prev,
      [name]: numericFields.includes(name) ? (value ? parseFloat(value) : 0) : value,
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!formData.name) {
      toast.error('Please fill all required fields')
      return
    }

    setLoading(true)
    try {
      await onSubmit(formData)
      toast.success(campaign ? 'Campaign updated!' : 'Campaign created!')
      onCancel()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error saving campaign')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Campaign Info */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Campaign Information</h3>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Campaign Name *</label>
          <input
            type="text"
            name="name"
            value={formData.name || ''}
            onChange={handleChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
      </div>

      {/* Description */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
        <textarea
          name="description"
          value={formData.description || ''}
          onChange={handleChange}
          rows={3}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      {/* Type & Status */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Campaign Details</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Campaign Type *</label>
            <select
              name="campaign_type"
              value={formData.campaign_type || 'project_launch'}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
              <option value="project_launch">Project Launch</option>
              <option value="phase_launch">Phase Launch</option>
              <option value="grand_opening">Grand Opening</option>
              <option value="season_promotion">Seasonal Promotion</option>
              <option value="inventory_clearance">Inventory Clearance</option>
              <option value="referral">Referral Campaign</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select
              name="status"
              value={formData.status || 'draft'}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="draft">Draft</option>
              <option value="active">Active</option>
              <option value="paused">Paused</option>
              <option value="completed">Completed</option>
              <option value="archived">Archived</option>
            </select>
          </div>
        </div>
      </div>

      {/* Dates */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Campaign Duration</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Start Date</label>
            <input
              type="date"
              name="start_date"
              value={formData.start_date || ''}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">End Date</label>
            <input
              type="date"
              name="end_date"
              value={formData.end_date || ''}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      </div>

      {/* Budget & Targeting */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Budget & Targeting</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Budget (₹) *</label>
            <input
              type="number"
              name="budget"
              value={formData.budget || 0}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
              min="0"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Target Segment *</label>
            <select
              name="target_segment"
              value={formData.target_segment || 'residential'}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
              <option value="residential">Residential</option>
              <option value="commercial">Commercial</option>
              <option value="nri">NRI</option>
              <option value="corporate">Corporate</option>
              <option value="investor">Investor</option>
              <option value="all">All Segments</option>
            </select>
          </div>
        </div>

        {/* Price Range */}
        <div className="grid grid-cols-2 gap-4 mt-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Min Price (₹)</label>
            <input
              type="number"
              name="target_price_range_min"
              value={formData.target_price_range_min || 0}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              min="0"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Max Price (₹)</label>
            <input
              type="number"
              name="target_price_range_max"
              value={formData.target_price_range_max || 0}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              min="0"
            />
          </div>
        </div>

        {/* Expected Leads & Bookings */}
        <div className="grid grid-cols-2 gap-4 mt-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Expected Leads</label>
            <input
              type="number"
              name="expected_leads"
              value={formData.expected_leads || 0}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              min="0"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Expected Bookings</label>
            <input
              type="number"
              name="expected_bookings"
              value={formData.expected_bookings || 0}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              min="0"
            />
          </div>
        </div>
      </div>

      {/* Actions */}
      <div className="flex gap-4 pt-4 border-t border-gray-200">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
        >
          {loading ? 'Saving...' : campaign ? 'Update Campaign' : 'Create Campaign'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 bg-gray-200 text-gray-900 py-2 rounded-lg hover:bg-gray-300 font-medium"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
