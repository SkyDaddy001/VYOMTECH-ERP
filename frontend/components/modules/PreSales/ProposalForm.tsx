'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import { Proposal } from '@/types/presales'

interface ProposalFormProps {
  proposal?: Proposal | null
  opportunities: Array<{ id?: string; name: string }>
  onSubmit: (data: Partial<Proposal>) => Promise<void>
  onCancel: () => void
}

export default function ProposalForm({ proposal, opportunities, onSubmit, onCancel }: ProposalFormProps) {
  const [formData, setFormData] = useState<Partial<Proposal>>({
    title: '',
    description: '',
    amount: 0,
    validity_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    status: 'draft',
  })
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (proposal) {
      setFormData({
        id: proposal.id,
        title: proposal.title,
        description: proposal.description,
        amount: proposal.amount,
        validity_date: proposal.validity_date?.split('T')[0] || '',
        status: proposal.status,
        opportunity_id: proposal.opportunity_id,
      })
    }
  }, [proposal])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: name === 'amount' ? (value ? parseInt(value) : 0) : value,
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!formData.title || !formData.amount) {
      toast.error('Please fill all required fields')
      return
    }

    setLoading(true)
    try {
      await onSubmit(formData)
      toast.success(proposal ? 'Proposal updated!' : 'Proposal created!')
      onCancel()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error saving proposal')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Proposal Info */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Proposal Information</h3>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Title *</label>
          <input
            type="text"
            name="title"
            value={formData.title || ''}
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

      {/* Amount & Opportunity */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Proposal Details</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Amount (₹) *</label>
            <input
              type="number"
              name="amount"
              value={formData.amount || 0}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Opportunity</label>
            <select
              name="opportunity_id"
              value={formData.opportunity_id || ''}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Select Opportunity</option>
              {opportunities.map((opp) => (
                <option key={opp.id} value={opp.id}>
                  {opp.name}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* Validity Date & Status */}
      <div>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Validity Date</label>
            <input
              type="date"
              name="validity_date"
              value={formData.validity_date || ''}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
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
              <option value="sent">Sent</option>
              <option value="viewed">Viewed</option>
              <option value="accepted">Accepted</option>
              <option value="rejected">Rejected</option>
            </select>
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
          {loading ? 'Saving...' : proposal ? 'Update Proposal' : 'Create Proposal'}
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
