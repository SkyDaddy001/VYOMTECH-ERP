'use client'

import { useState } from 'react'
import { ChangeRequest } from '@/types/postsales'

interface ChangeRequestFormProps {
  changeRequest?: ChangeRequest | null
  onSubmit: (data: Partial<ChangeRequest>) => Promise<void>
  onCancel: () => void
}

export default function ChangeRequestForm({ changeRequest, onSubmit, onCancel }: ChangeRequestFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<ChangeRequest>>(
    changeRequest || {
      crm_number: '',
      booking_id: '',
      customer_id: '',
      request_type: 'unit_modification',
      description: '',
      impact: 'no_cost_change',
      cost_difference: 0,
      status: 'submitted',
      request_date: new Date().toISOString().split('T')[0],
      approval_date: '',
      completion_date: '',
      notes: '',
    }
  )

  const handleChange = (field: string, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await onSubmit(formData)
    } finally {
      setLoading(false)
    }
  }

  const requestTypes = [
    { value: 'unit_modification', label: 'Unit Modification' },
    { value: 'floor_change', label: 'Floor Change' },
    { value: 'parking_choice', label: 'Parking Choice' },
    { value: 'amenity_upgrade', label: 'Amenity Upgrade' },
    { value: 'specification_change', label: 'Specification Change' },
    { value: 'other', label: 'Other' },
  ]

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <span className="text-2xl">📝</span> {changeRequest ? 'Edit Change Request' : 'Submit Change Request'}
        </h2>
      </div>

      {/* Header Information */}
      <div className="bg-blue-50 rounded-lg p-4 border border-blue-200">
        <h3 className="text-sm font-semibold text-blue-900 mb-4">Header Information</h3>
        <div className="grid grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">CRM Number *</label>
            <input
              type="text"
              required
              value={formData.crm_number || ''}
              onChange={(e) => handleChange('crm_number', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="CRM-001"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Booking ID</label>
            <input
              type="text"
              value={formData.booking_id || ''}
              onChange={(e) => handleChange('booking_id', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="BOOK-001"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Customer ID</label>
            <input
              type="text"
              value={formData.customer_id || ''}
              onChange={(e) => handleChange('customer_id', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="CUS-001"
            />
          </div>
        </div>
      </div>

      {/* Request Type & Status */}
      <div className="grid grid-cols-3 gap-4">
        <div className="bg-green-50 rounded-lg p-4 border border-green-200">
          <label className="block text-sm font-medium text-gray-700 mb-2">Request Type *</label>
          <select
            required
            value={formData.request_type || 'unit_modification'}
            onChange={(e) => handleChange('request_type', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
          >
            {requestTypes.map(t => (
              <option key={t.value} value={t.value}>{t.label}</option>
            ))}
          </select>
        </div>
        <div className="bg-yellow-50 rounded-lg p-4 border border-yellow-200">
          <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
          <select
            value={formData.status || 'submitted'}
            onChange={(e) => handleChange('status', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500"
          >
            <option value="submitted">Submitted</option>
            <option value="under_review">Under Review</option>
            <option value="approved">Approved</option>
            <option value="rejected">Rejected</option>
            <option value="implemented">Implemented</option>
          </select>
        </div>
        <div className="bg-purple-50 rounded-lg p-4 border border-purple-200">
          <label className="block text-sm font-medium text-gray-700 mb-2">Impact Type</label>
          <select
            value={formData.impact || 'no_cost_change'}
            onChange={(e) => handleChange('impact', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
          >
            <option value="no_cost_change">No Cost Change</option>
            <option value="additional_cost">Additional Cost</option>
            <option value="cost_reduction">Cost Reduction</option>
          </select>
        </div>
      </div>

      {/* Change Details */}
      <div className="bg-indigo-50 rounded-lg p-4 border border-indigo-200">
        <h3 className="text-sm font-semibold text-indigo-900 mb-4">Change Details</h3>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Description *</label>
          <textarea
            required
            value={formData.description || ''}
            onChange={(e) => handleChange('description', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="Detailed description of the change..."
            rows={4}
          />
        </div>
      </div>

      {/* Financial Impact */}
      <div className="bg-orange-50 rounded-lg p-4 border border-orange-200">
        <label className="block text-sm font-medium text-gray-700 mb-2">Cost Difference (₹)</label>
        <input
          type="number"
          value={formData.cost_difference || ''}
          onChange={(e) => handleChange('cost_difference', parseFloat(e.target.value))}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500"
          step="0.01"
        />
      </div>

      {/* Timeline */}
      <div className="grid grid-cols-3 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Request Date *</label>
          <input
            type="date"
            required
            value={formData.request_date || ''}
            onChange={(e) => handleChange('request_date', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Approval Date</label>
          <input
            type="date"
            value={formData.approval_date || ''}
            onChange={(e) => handleChange('approval_date', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Completion Date</label>
          <input
            type="date"
            value={formData.completion_date || ''}
            onChange={(e) => handleChange('completion_date', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      {/* Notes */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Additional Notes</label>
        <textarea
          value={formData.notes || ''}
          onChange={(e) => handleChange('notes', e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Additional notes and remarks..."
          rows={3}
        />
      </div>

      {/* Form Actions */}
      <div className="flex gap-3 pt-4">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium transition-colors"
        >
          {loading ? 'Saving...' : changeRequest ? 'Update Change Request' : 'Submit Change Request'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 bg-gray-200 text-gray-900 py-2 rounded-lg hover:bg-gray-300 font-medium transition-colors"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
