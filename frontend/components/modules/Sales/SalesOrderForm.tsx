'use client'

import { useState } from 'react'
import { SalesOrder } from '@/types/sales'

interface SalesOrderFormProps {
  order?: SalesOrder | null
  onSubmit: (data: Partial<SalesOrder>) => Promise<void>
  onCancel: () => void
}

export default function SalesOrderForm({ order, onSubmit, onCancel }: SalesOrderFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<SalesOrder>>(
    order || {
      booking_number: '',
      customer_id: '',
      customer_name: '',
      project_id: '',
      property_id: '',
      unit_type: 'residential',
      super_area: 0,
      carpet_area: 0,
      base_price: 0,
      base_price_per_sqft: 0,
      total_amount: 0,
      gst_amount: 0,
      registration_amount: 0,
      discount_amount: 0,
      net_amount: 0,
      booking_stage: 'inquiry',
      booking_date: new Date().toISOString().split('T')[0],
      sales_executive_id: '',
      notes: '',
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
    // Recalculate net amount
    if (['total_amount', 'gst_amount', 'registration_amount', 'discount_amount'].includes(field)) {
      updated.net_amount = (updated.total_amount || 0) + (updated.gst_amount || 0) + (updated.registration_amount || 0) - (updated.discount_amount || 0)
    }
    setFormData(updated)
  }

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Booking Number *</label>
          <input
            type="text"
            required
            value={formData.booking_number || ''}
            onChange={(e) => handleChange('booking_number', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="BK-001"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Customer Name *</label>
          <input
            type="text"
            required
            value={formData.customer_name || ''}
            onChange={(e) => handleChange('customer_name', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter customer name"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Project *</label>
          <input
            type="text"
            required
            value={formData.project_id || ''}
            onChange={(e) => handleChange('project_id', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Project ID"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Property *</label>
          <input
            type="text"
            required
            value={formData.property_id || ''}
            onChange={(e) => handleChange('property_id', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Property ID"
          />
        </div>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Unit Type *</label>
          <select
            value={formData.unit_type || 'residential'}
            onChange={(e) => handleChange('unit_type', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          >
            <option value="residential">Residential</option>
            <option value="commercial">Commercial</option>
            <option value="parking">Parking</option>
            <option value="other">Other</option>
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Super Area (sqft) *</label>
          <input
            type="number"
            required
            value={formData.super_area || 0}
            onChange={(e) => handleChange('super_area', parseFloat(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            min="0"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Carpet Area (sqft) *</label>
          <input
            type="number"
            required
            value={formData.carpet_area || 0}
            onChange={(e) => handleChange('carpet_area', parseFloat(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            min="0"
          />
        </div>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Base Price (₹) *</label>
          <input
            type="number"
            required
            value={formData.base_price || 0}
            onChange={(e) => handleChange('base_price', parseFloat(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            min="0"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Price per Sqft (₹) *</label>
          <input
            type="number"
            required
            value={formData.base_price_per_sqft || 0}
            onChange={(e) => handleChange('base_price_per_sqft', parseFloat(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            min="0"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Booking Date *</label>
          <input
            type="date"
            required
            value={formData.booking_date || ''}
            onChange={(e) => handleChange('booking_date', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Total Amount (₹) *</label>
          <input
            type="number"
            required
            value={formData.total_amount || 0}
            onChange={(e) => handleChange('total_amount', parseFloat(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            min="0"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">GST Amount (₹)</label>
          <input
            type="number"
            value={formData.gst_amount || 0}
            onChange={(e) => handleChange('gst_amount', parseFloat(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            min="0"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Registration Amount (₹)</label>
          <input
            type="number"
            value={formData.registration_amount || 0}
            onChange={(e) => handleChange('registration_amount', parseFloat(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            min="0"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Discount Amount (₹)</label>
          <input
            type="number"
            value={formData.discount_amount || 0}
            onChange={(e) => handleChange('discount_amount', parseFloat(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            min="0"
          />
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Net Amount (₹)</label>
        <input
          type="number"
          disabled
          value={formData.net_amount || 0}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Booking Stage *</label>
        <select
          value={formData.booking_stage || 'inquiry'}
          onChange={(e) => handleChange('booking_stage', e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          required
        >
          <option value="inquiry">Inquiry</option>
          <option value="quote">Quote</option>
          <option value="booking">Booking</option>
          <option value="agreement">Agreement</option>
          <option value="completed">Completed</option>
          <option value="cancelled">Cancelled</option>
        </select>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Sales Executive ID *</label>
        <input
          type="text"
          required
          value={formData.sales_executive_id || ''}
          onChange={(e) => handleChange('sales_executive_id', e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Sales executive ID"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Notes</label>
        <textarea
          value={formData.notes || ''}
          onChange={(e) => handleChange('notes', e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Booking notes..."
          rows={3}
        />
      </div>

      <div className="flex gap-3">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
        >
          {loading ? 'Saving...' : 'Save Order'}
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
