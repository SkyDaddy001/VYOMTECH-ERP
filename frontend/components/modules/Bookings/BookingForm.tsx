'use client'

import React, { useState } from 'react'
import { Booking } from '@/types/bookings'

interface BookingFormProps {
  booking?: Booking
  onSubmit: (data: Partial<Booking>) => Promise<void>
  onCancel: () => void
}

export default function BookingForm({ booking, onSubmit, onCancel }: BookingFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<Booking>>(
    booking || {
      unit_id: '',
      customer_id: '',
      booking_date: '',
      booking_reference: '',
      booking_status: 'inquired',
      rate_per_sqft: 0,
      booking_amount: 0,
      total_cost: 0,
      paid_amount: 0,
      balance_amount: 0,
    }
  )

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      formData.balance_amount = (formData.total_cost || 0) - (formData.paid_amount || 0)
      await onSubmit(formData)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6 max-w-2xl">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Unit ID *</label>
          <input
            type="text"
            required
            value={formData.unit_id || ''}
            onChange={(e) => setFormData({ ...formData, unit_id: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Customer ID *</label>
          <input
            type="text"
            required
            value={formData.customer_id || ''}
            onChange={(e) => setFormData({ ...formData, customer_id: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Booking Reference *</label>
          <input
            type="text"
            required
            value={formData.booking_reference || ''}
            onChange={(e) => setFormData({ ...formData, booking_reference: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Booking Date *</label>
          <input
            type="date"
            required
            value={formData.booking_date || ''}
            onChange={(e) => setFormData({ ...formData, booking_date: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Status *</label>
          <select
            required
            value={formData.booking_status || 'inquired'}
            onChange={(e) => setFormData({ ...formData, booking_status: e.target.value as any })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="inquired">Inquired</option>
            <option value="interested">Interested</option>
            <option value="booked">Booked</option>
            <option value="registered">Registered</option>
            <option value="handed_over">Handed Over</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Rate Per Sqft *</label>
          <input
            type="number"
            required
            value={formData.rate_per_sqft || 0}
            onChange={(e) => setFormData({ ...formData, rate_per_sqft: parseFloat(e.target.value) })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Booking Amount *</label>
          <input
            type="number"
            required
            value={formData.booking_amount || 0}
            onChange={(e) => setFormData({ ...formData, booking_amount: parseFloat(e.target.value) })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Total Cost *</label>
          <input
            type="number"
            required
            value={formData.total_cost || 0}
            onChange={(e) => setFormData({ ...formData, total_cost: parseFloat(e.target.value) })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Paid Amount *</label>
          <input
            type="number"
            required
            value={formData.paid_amount || 0}
            onChange={(e) => setFormData({ ...formData, paid_amount: parseFloat(e.target.value) })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <div className="flex gap-3">
        <button
          type="submit"
          disabled={loading}
          className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium disabled:opacity-50"
        >
          {loading ? 'Saving...' : booking ? 'Update Booking' : 'Create Booking'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 font-medium"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
