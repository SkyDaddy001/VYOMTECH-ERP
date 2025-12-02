'use client'

import { useState } from 'react'
import { Booking, BookingPaymentSchedule } from '@/types/bookings'

interface PaymentPlanFormProps {
  booking?: Booking | null
  onSubmit: (data: Partial<BookingPaymentSchedule>[]) => Promise<void>
  onCancel: () => void
}

export default function PaymentPlanForm({ booking, onSubmit, onCancel }: PaymentPlanFormProps) {
  const [loading, setLoading] = useState(false)
  const [paymentSchedule, setPaymentSchedule] = useState<Partial<BookingPaymentSchedule>[]>([
    {
      booking_id: booking?.id || '',
      number_of_installments: 6,
      payment_stage: 'booking',
      installment_date: '',
      installment_amount: 0,
      status: 'pending',
    },
  ])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await onSubmit(paymentSchedule)
    } finally {
      setLoading(false)
    }
  }

  const updateItem = (index: number, field: string, value: any) => {
    const updated = [...paymentSchedule]
    updated[index] = { ...updated[index], [field]: value }
    setPaymentSchedule(updated)
  }

  const addItem = () => {
    setPaymentSchedule([...paymentSchedule, {
      booking_id: booking?.id || '',
      number_of_installments: 6,
      payment_stage: 'agreement',
      installment_date: '',
      installment_amount: 0,
      status: 'pending',
    }])
  }

  const removeItem = (index: number) => {
    setPaymentSchedule(paymentSchedule.filter((_, i) => i !== index))
  }

  const stages: Array<BookingPaymentSchedule['payment_stage']> = [
    'booking', 'agreement', 'foundation', 'structure', 'finishing', 'possession'
  ]

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <span className="text-2xl">í²°</span> Payment Plan
        </h2>
      </div>

      {/* Payment Schedule */}
      <div className="space-y-3">
        <div className="flex justify-between items-center">
          <h3 className="text-lg font-semibold text-gray-800">Payment Schedule</h3>
          <button
            type="button"
            onClick={addItem}
            className="px-3 py-1 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 font-medium"
          >
            + Add Installment
          </button>
        </div>

        <div className="overflow-x-auto border border-gray-200 rounded-lg">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Payment Stage</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Installment Date</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Amount</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Status</th>
                <th className="px-4 py-3 text-center text-sm font-semibold text-gray-700">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {paymentSchedule.map((item, index) => (
                <tr key={index}>
                  <td className="px-4 py-3">
                    <select
                      value={item.payment_stage || 'booking'}
                      onChange={(e) => updateItem(index, 'payment_stage', e.target.value)}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-sm"
                    >
                      {stages.map(s => (
                        <option key={s} value={s}>{s.charAt(0).toUpperCase() + s.slice(1)}</option>
                      ))}
                    </select>
                  </td>
                  <td className="px-4 py-3">
                    <input
                      type="date"
                      value={item.installment_date || ''}
                      onChange={(e) => updateItem(index, 'installment_date', e.target.value)}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-sm"
                    />
                  </td>
                  <td className="px-4 py-3">
                    <input
                      type="number"
                      value={item.installment_amount || 0}
                      onChange={(e) => updateItem(index, 'installment_amount', parseFloat(e.target.value))}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-sm"
                      step="0.01"
                    />
                  </td>
                  <td className="px-4 py-3">
                    <select
                      value={item.status || 'pending'}
                      onChange={(e) => updateItem(index, 'status', e.target.value)}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-sm"
                    >
                      <option value="pending">Pending</option>
                      <option value="received">Received</option>
                      <option value="overdue">Overdue</option>
                      <option value="bounced">Bounced</option>
                    </select>
                  </td>
                  <td className="px-4 py-3 text-center">
                    <button
                      type="button"
                      onClick={() => removeItem(index)}
                      className="text-red-600 hover:text-red-900 font-medium text-sm"
                    >
                      Remove
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Form Actions */}
      <div className="flex gap-3 pt-4">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium transition-colors"
        >
          {loading ? 'Saving...' : 'Save Payment Plan'}
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
