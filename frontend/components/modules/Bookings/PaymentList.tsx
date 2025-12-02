'use client'

import React, { useState } from 'react'
import toast from 'react-hot-toast'
import { BookingPayment } from '@/types/bookings'

interface PaymentListProps {
  payments: BookingPayment[]
  loading: boolean
  onEdit: (payment: BookingPayment) => void
  onDelete: (payment: BookingPayment) => void
}

export default function PaymentList({ payments, loading, onEdit, onDelete }: PaymentListProps) {
  const [filterStatus, setFilterStatus] = useState<string>('all')

  const filteredPayments = filterStatus === 'all'
    ? payments
    : payments.filter(p => p.status === filterStatus)

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      'pending': 'bg-yellow-100 text-yellow-800',
      'cleared': 'bg-green-100 text-green-800',
      'bounced': 'bg-red-100 text-red-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const getModeLabel = (mode: string) => {
    const labels: Record<string, string> = {
      'cash': 'Cash',
      'cheque': 'Cheque',
      'transfer': 'Bank Transfer',
      'neft': 'NEFT',
      'rtgs': 'RTGS',
      'dd': 'DD',
    }
    return labels[mode] || mode
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-600">Loading payments...</div>
  }

  return (
    <div className="space-y-4">
      <select
        value={filterStatus}
        onChange={(e) => setFilterStatus(e.target.value)}
        className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="all">All Payments</option>
        <option value="pending">Pending</option>
        <option value="cleared">Cleared</option>
        <option value="bounced">Bounced</option>
      </select>

      <div className="space-y-3">
        {filteredPayments.map((payment) => (
          <div key={payment.id} className="bg-white rounded-lg shadow p-4 border border-gray-200 hover:shadow-lg transition">
            <div className="flex justify-between items-start mb-3">
              <div>
                <p className="font-semibold text-gray-900">₹{(payment.amount / 100000).toFixed(2)}L</p>
                <p className="text-xs text-gray-600">{getModeLabel(payment.payment_mode)}</p>
              </div>
              <span className={`px-2 py-1 text-xs font-medium rounded ${getStatusBadge(payment.status)}`}>
                {payment.status}
              </span>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-3 gap-3 mb-3 text-sm">
              <div>
                <p className="text-gray-600 text-xs">Date</p>
                <p className="font-medium text-xs">{new Date(payment.payment_date).toLocaleDateString()}</p>
              </div>
              <div>
                <p className="text-gray-600 text-xs">Receipt #</p>
                <p className="font-medium text-xs">{payment.receipt_number}</p>
              </div>
              <div>
                <p className="text-gray-600 text-xs">Towards</p>
                <p className="font-medium text-xs">{payment.towards}</p>
              </div>
            </div>

            <div className="flex gap-2">
              <button
                onClick={() => onEdit(payment)}
                className="flex-1 px-3 py-2 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition"
              >
                Edit
              </button>
              <button
                onClick={() => {
                  if (confirm('Delete this payment?')) {
                    onDelete(payment)
                  }
                }}
                className="flex-1 px-3 py-2 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition"
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>

      {filteredPayments.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No payments recorded</p>
        </div>
      )}
    </div>
  )
}
