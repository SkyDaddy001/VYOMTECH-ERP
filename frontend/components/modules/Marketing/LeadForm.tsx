'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import { Lead } from '@/types/marketing'

interface LeadFormProps {
  lead?: Lead | null
  onSubmit: (data: Partial<Lead>) => Promise<void>
  onCancel: () => void
}

export default function LeadForm({ lead, onSubmit, onCancel }: LeadFormProps) {
  const [formData, setFormData] = useState<Partial<Lead>>({
    prospect_name: '',
    email: '',
    phone: '',
    property_type_interest: 'residential',
    budget_range_min: 0,
    budget_range_max: 0,
    source: 'website',
    lead_stage: 'inquiry',
    lead_quality: 'warm',
  })
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (lead) {
      setFormData({
        id: lead.id,
        prospect_name: lead.prospect_name,
        email: lead.email,
        phone: lead.phone,
        alternate_phone: lead.alternate_phone,
        property_type_interest: lead.property_type_interest,
        budget_range_min: lead.budget_range_min,
        budget_range_max: lead.budget_range_max,
        required_bhk: lead.required_bhk,
        source: lead.source,
        lead_stage: lead.lead_stage,
        lead_quality: lead.lead_quality,
        notes: lead.notes,
      })
    }
  }, [lead])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    const numericFields = ['budget_range_min', 'budget_range_max', 'required_bhk']
    setFormData((prev) => ({
      ...prev,
      [name]: numericFields.includes(name) ? (value ? parseFloat(value) : 0) : value,
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!formData.prospect_name || !formData.email || !formData.phone) {
      toast.error('Please fill all required fields')
      return
    }

    setLoading(true)
    try {
      await onSubmit(formData)
      toast.success(lead ? 'Lead updated!' : 'Lead created!')
      onCancel()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error saving lead')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Contact Information */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Contact Information</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Prospect Name *</label>
            <input
              type="text"
              name="prospect_name"
              value={formData.prospect_name || ''}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Email *</label>
            <input
              type="email"
              name="email"
              value={formData.email || ''}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4 mt-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Phone *</label>
            <input
              type="tel"
              name="phone"
              value={formData.phone || ''}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Alternate Phone</label>
            <input
              type="tel"
              name="alternate_phone"
              value={formData.alternate_phone || ''}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      </div>

      {/* Property Interest */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Property Interest</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Property Type *</label>
            <select
              name="property_type_interest"
              value={formData.property_type_interest || 'residential'}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
              <option value="residential">Residential</option>
              <option value="commercial">Commercial</option>
              <option value="mixed_use">Mixed-use</option>
              <option value="any">Any</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Required BHK</label>
            <select
              name="required_bhk"
              value={formData.required_bhk || 2}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value={1}>1 BHK</option>
              <option value={2}>2 BHK</option>
              <option value={3}>3 BHK</option>
              <option value={4}>4+ BHK</option>
            </select>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4 mt-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Budget Range Min (₹) *</label>
            <input
              type="number"
              name="budget_range_min"
              value={formData.budget_range_min || 0}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
              min="0"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Budget Range Max (₹) *</label>
            <input
              type="number"
              name="budget_range_max"
              value={formData.budget_range_max || 0}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
              min="0"
            />
          </div>
        </div>
      </div>

      {/* Lead Details */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Lead Details</h3>
        <div className="grid grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Source *</label>
            <select
              name="source"
              value={formData.source || 'website'}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
              <option value="website">Website</option>
              <option value="referral">Referral</option>
              <option value="campaign">Campaign</option>
              <option value="broker">Broker</option>
              <option value="social">Social</option>
              <option value="event">Event</option>
              <option value="walk_in">Walk-in</option>
              <option value="cold_call">Cold Call</option>
              <option value="other">Other</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Lead Stage *</label>
            <select
              name="lead_stage"
              value={formData.lead_stage || 'inquiry'}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
              <option value="inquiry">Inquiry</option>
              <option value="interested">Interested</option>
              <option value="site_visit_scheduled">Site Visit Scheduled</option>
              <option value="site_visit_done">Site Visit Done</option>
              <option value="qualified">Qualified</option>
              <option value="negotiation">Negotiation</option>
              <option value="lost">Lost</option>
              <option value="converted">Converted</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Lead Quality *</label>
            <select
              name="lead_quality"
              value={formData.lead_quality || 'warm'}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
              <option value="hot">🔥 Hot</option>
              <option value="warm">🌤️ Warm</option>
              <option value="cold">❄️ Cold</option>
            </select>
          </div>
        </div>
      </div>

      {/* Notes */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
        <textarea
          name="notes"
          value={formData.notes || ''}
          onChange={handleChange}
          rows={3}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      {/* Actions */}
      <div className="flex gap-4 pt-4 border-t border-gray-200">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
        >
          {loading ? 'Saving...' : lead ? 'Update Lead' : 'Create Lead'}
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
