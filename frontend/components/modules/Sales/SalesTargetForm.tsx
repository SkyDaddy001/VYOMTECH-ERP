'use client'

import { useState } from 'react'
import { SalesTarget } from '@/types/sales'

interface SalesTargetFormProps {
  target?: SalesTarget | null
  onSubmit: (data: Partial<SalesTarget>) => Promise<void>
  onCancel: () => void
}

export default function SalesTargetForm({ target, onSubmit, onCancel }: SalesTargetFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<SalesTarget>>(
    target || {
      sales_executive_id: '',
      period: new Date().toISOString().slice(0, 7),
      target_bookings: 0,
      target_amount: 0,
      achieved_bookings: 0,
      achieved_amount: 0,
      achievement_percentage: 0,
      status: 'not_started',
    }
  )

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await onSubmit(formData)
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (field: string, value: any) => {
    const updated = { ...formData, [field]: value }
    // Calculate achievement percentage
    if (['target_amount', 'achieved_amount'].includes(field)) {
      updated.achievement_percentage =
        (updated.target_amount || 0) > 0 ? ((updated.achieved_amount || 0) / (updated.target_amount || 1)) * 100 : 0
    }
    setFormData(updated)
  }

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Sales Executive ID *</label>
          <input
            type="text"
            required
            value={formData.sales_executive_id || ''}
            onChange={(e) => handleChange('sales_executive_id', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="SE-001"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Period (YYYY-MM) *</label>
          <input
            type="month"
            required
            value={formData.period || ''}
            onChange={(e) => handleChange('period', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Target Bookings *</label>
          <input
            type="number"
            required
            value={formData.target_bookings || 0}
            onChange={(e) => handleChange('target_bookings', Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0"
            min="0"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Achieved Bookings</label>
          <input
            type="number"
            value={formData.achieved_bookings || 0}
            onChange={(e) => handleChange('achieved_bookings', Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0"
            min="0"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Target Amount (₹) *</label>
          <input
            type="number"
            required
            value={formData.target_amount || 0}
            onChange={(e) => handleChange('target_amount', Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0"
            min="0"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Achieved Amount (₹)</label>
          <input
            type="number"
            value={formData.achieved_amount || 0}
            onChange={(e) => handleChange('achieved_amount', Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0"
            min="0"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Achievement % (Auto)</label>
          <input
            type="number"
            disabled
            value={(formData.achievement_percentage || 0).toFixed(2)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
          <select
            value={formData.status || 'not_started'}
            onChange={(e) => handleChange('status', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="not_started">Not Started</option>
            <option value="in_progress">In Progress</option>
            <option value="completed">Completed</option>
            <option value="exceeded">Exceeded</option>
          </select>
        </div>
      </div>

      <div className="flex gap-3">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
        >
          {loading ? 'Saving...' : 'Save Target'}
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
