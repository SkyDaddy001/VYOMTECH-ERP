'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import { PreSalesOpportunity } from '@/types/presales'

interface OpportunityFormProps {
  opportunity?: PreSalesOpportunity | null
  onSubmit: (data: Partial<PreSalesOpportunity>) => Promise<void>
  onCancel: () => void
}

export default function OpportunityForm({ opportunity, onSubmit, onCancel }: OpportunityFormProps) {
  const [formData, setFormData] = useState<Partial<PreSalesOpportunity>>({
    name: '',
    description: '',
    stage: 'prospecting',
    value: 0,
    probability: 50,
    expected_close_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
  })
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (opportunity) {
      setFormData({
        id: opportunity.id,
        name: opportunity.name,
        description: opportunity.description,
        stage: opportunity.stage,
        value: opportunity.value,
        probability: opportunity.probability,
        expected_close_date: opportunity.expected_close_date?.split('T')[0] || '',
        notes: opportunity.notes,
      })
    }
  }, [opportunity])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: name === 'value' || name === 'probability' ? (value ? parseInt(value) : 0) : value,
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
      toast.success(opportunity ? 'Opportunity updated!' : 'Opportunity created!')
      onCancel()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error saving opportunity')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Opportunity Info */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Opportunity Information</h3>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Opportunity Name *</label>
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

      {/* Stage, Value, Probability */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Opportunity Details</h3>
        <div className="grid grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Stage</label>
            <select
              name="stage"
              value={formData.stage || 'prospecting'}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="prospecting">Prospecting</option>
              <option value="qualification">Qualification</option>
              <option value="proposal">Proposal</option>
              <option value="negotiation">Negotiation</option>
              <option value="closed_won">Closed Won</option>
              <option value="closed_lost">Closed Lost</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Deal Value (₹)</label>
            <input
              type="number"
              name="value"
              value={formData.value || 0}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Probability (%)</label>
            <input
              type="number"
              name="probability"
              value={formData.probability || 50}
              onChange={handleChange}
              min="0"
              max="100"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      </div>

      {/* Expected Close Date */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Expected Close Date</label>
        <input
          type="date"
          name="expected_close_date"
          value={formData.expected_close_date || ''}
          onChange={handleChange}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      {/* Notes */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
        <textarea
          name="notes"
          value={formData.notes || ''}
          onChange={handleChange}
          rows={2}
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
          {loading ? 'Saving...' : opportunity ? 'Update Opportunity' : 'Create Opportunity'}
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
