'use client'

import React, { useState } from 'react'
import { Booking } from '@/types/bookings'

interface BookingListProps {
  bookings: Booking[]
  loading: boolean
  onEdit: (booking: Booking) => void
  onDelete: (booking: Booking) => void
  onUpdateStatus: (booking: Booking, status: string) => void
}

export default function BookingList({ bookings, loading, onEdit, onDelete, onUpdateStatus }: BookingListProps) {
  const [filterStatus, setFilterStatus] = useState<string>('all')

  const filteredBookings = filterStatus === 'all' 
    ? bookings 
    : bookings.filter(b => b.booking_status === filterStatus)

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      'inquired': 'bg-gray-100 text-gray-800',
      'interested': 'bg-blue-100 text-blue-800',
      'booked': 'bg-yellow-100 text-yellow-800',
      'registered': 'bg-green-100 text-green-800',
      'handed_over': 'bg-purple-100 text-purple-800',
      'cancelled': 'bg-red-100 text-red-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-600">Loading bookings...</div>
  }

  return (
    <div className="space-y-4">
      <select
        value={filterStatus}
        onChange={(e) => setFilterStatus(e.target.value)}
        className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="all">All Bookings</option>
        <option value="inquired">Inquired</option>
        <option value="interested">Interested</option>
        <option value="booked">Booked</option>
        <option value="registered">Registered</option>
        <option value="handed_over">Handed Over</option>
        <option value="cancelled">Cancelled</option>
      </select>

      <div className="space-y-3">
        {filteredBookings.map((booking) => (
          <div key={booking.id} className="bg-white rounded-lg shadow p-4 border border-gray-200 hover:shadow-lg transition">
            <div className="flex justify-between items-start mb-3">
              <div>
                <h3 className="font-semibold text-gray-900">{booking.customer_name}</h3>
                <p className="text-xs text-gray-600">{booking.booking_reference}</p>
              </div>
              <span className={`px-2 py-1 text-xs font-medium rounded ${getStatusBadge(booking.booking_status)}`}>
                {booking.booking_status.replace('_', ' ')}
              </span>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-4 text-sm">
              <div>
                <p className="text-gray-600 text-xs">Unit</p>
                <p className="font-medium">{booking.unit_number || 'N/A'}</p>
              </div>
              <div>
                <p className="text-gray-600 text-xs">Booking Date</p>
                <p className="font-medium text-xs">{new Date(booking.booking_date).toLocaleDateString()}</p>
              </div>
              <div>
                <p className="text-gray-600 text-xs">Total Cost</p>
                <p className="font-medium">₹{(booking.total_cost / 1000000).toFixed(2)}L</p>
              </div>
              <div>
                <p className="text-gray-600 text-xs">Balance</p>
                <p className="font-medium text-orange-600">₹{(booking.balance_amount / 1000000).toFixed(2)}L</p>
              </div>
            </div>

            {/* Progress Bar */}
            <div className="mb-3">
              <div className="flex justify-between items-center mb-1">
                <span className="text-xs text-gray-600">Payment Progress</span>
                <span className="text-xs font-medium">{((booking.paid_amount / booking.total_cost) * 100).toFixed(0)}%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className="bg-green-600 h-2 rounded-full transition-all"
                  style={{ width: `${(booking.paid_amount / booking.total_cost) * 100}%` }}
                />
              </div>
            </div>

            <div className="flex gap-2 flex-wrap">
              <button
                onClick={() => onEdit(booking)}
                className="px-3 py-2 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition"
              >
                Edit
              </button>
              <select
                value={booking.booking_status}
                onChange={(e) => onUpdateStatus(booking, e.target.value)}
                className="px-3 py-2 text-xs font-medium border border-gray-300 rounded hover:bg-gray-50 transition"
              >
                <option value="inquired">Inquired</option>
                <option value="interested">Interested</option>
                <option value="booked">Booked</option>
                <option value="registered">Registered</option>
                <option value="handed_over">Handed Over</option>
                <option value="cancelled">Cancelled</option>
              </select>
              <button
                onClick={() => {
                  if (confirm('Delete this booking?')) {
                    onDelete(booking)
                  }
                }}
                className="px-3 py-2 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition"
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>

      {filteredBookings.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No bookings found</p>
        </div>
      )}
    </div>
  )
}
