'use client'

import { useState } from 'react'
import { SalesQuota } from '@/types/sales'

interface SalesQuotaFormProps {
  quota?: SalesQuota | null
  onSubmit: (data: Partial<SalesQuota>) => Promise<void>
  onCancel: () => void
}

export default function SalesQuotaForm({ quota, onSubmit, onCancel }: SalesQuotaFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<SalesQuota>>(
    quota || {
      sales_executive_id: '',
      quarter: 'Q1',
      year: new Date().getFullYear(),
      quota_bookings: 0,
      quota_amount: 0,
      commission_rate_per_booking: 0,
      commission_rate_percentage: 0,
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
    setFormData({ ...formData, [field]: value })
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
          <label className="block text-sm font-medium text-gray-700 mb-2">Quarter *</label>
          <select
            required
            value={formData.quarter || 'Q1'}
            onChange={(e) => handleChange('quarter', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="Q1">Q1</option>
            <option value="Q2">Q2</option>
            <option value="Q3">Q3</option>
            <option value="Q4">Q4</option>
          </select>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Year *</label>
          <input
            type="number"
            required
            value={formData.year || new Date().getFullYear()}
            onChange={(e) => handleChange('year', Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            min="2020"
            max="2099"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Target Bookings *</label>
          <input
            type="number"
            required
            value={formData.quota_bookings || 0}
            onChange={(e) => handleChange('quota_bookings', Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0"
            min="0"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Quota Amount (₹) *</label>
          <input
            type="number"
            required
            value={formData.quota_amount || 0}
            onChange={(e) => handleChange('quota_amount', Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0"
            min="0"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Commission per Booking (₹) *</label>
          <input
            type="number"
            required
            value={formData.commission_rate_per_booking || 0}
            onChange={(e) => handleChange('commission_rate_per_booking', Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0"
            min="0"
          />
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Commission Percentage (%)</label>
        <input
          type="number"
          value={formData.commission_rate_percentage || 0}
          onChange={(e) => handleChange('commission_rate_percentage', Number(e.target.value))}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="0.00"
          step="0.01"
          min="0"
        />
      </div>

      <div className="flex gap-3">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
        >
          {loading ? 'Saving...' : 'Save Quota'}
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
