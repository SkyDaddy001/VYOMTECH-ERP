'use client'

import { Payment } from '@/types/accounts'

interface PaymentListProps {
  payments: Payment[]
  loading: boolean
  onDelete: (payment: Payment) => void
}

export default function PaymentList({ payments, loading, onDelete }: PaymentListProps) {
  if (loading) return <div className="text-center py-8 text-gray-500">Loading payments...</div>

  const totalPayments = payments.reduce((sum, p) => sum + p.amount, 0)
  const processedPayments = payments.filter(p => p.status === 'processed').length

  const methodIcons: Record<string, string> = {
    cash: '💵',
    check: '📋',
    card: '💳',
    bank_transfer: '🏦',
  }

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-gradient-to-br from-emerald-50 to-emerald-100 rounded-lg p-4 border border-emerald-200">
          <p className="text-gray-600 text-xs font-medium">Total Payments</p>
          <p className="text-2xl font-bold text-emerald-600 mt-1">₹{totalPayments.toLocaleString()}</p>
        </div>
        <div className="bg-gradient-to-br from-cyan-50 to-cyan-100 rounded-lg p-4 border border-cyan-200">
          <p className="text-gray-600 text-xs font-medium">Processed</p>
          <p className="text-2xl font-bold text-cyan-600 mt-1">{processedPayments}</p>
        </div>
        <div className="bg-gradient-to-br from-sky-50 to-sky-100 rounded-lg p-4 border border-sky-200">
          <p className="text-gray-600 text-xs font-medium">Total Count</p>
          <p className="text-2xl font-bold text-sky-600 mt-1">{payments.length}</p>
        </div>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Invoice #</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Amount</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Payment Date</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Method</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Status</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {payments.map((payment) => (
              <tr key={payment.id} className="hover:bg-gray-50 transition">
                <td className="px-6 py-4 text-sm font-medium text-gray-900">{payment.invoice_number || payment.invoice_id}</td>
                <td className="px-6 py-4 text-sm font-semibold text-gray-900">₹{payment.amount.toLocaleString()}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{new Date(payment.payment_date).toLocaleDateString()}</td>
                <td className="px-6 py-4 text-sm">
                  <span className="text-lg">{methodIcons[payment.payment_method] || '💰'}</span>
                  {payment.payment_method.replace('_', ' ').charAt(0).toUpperCase() + payment.payment_method.slice(1)}
                </td>
                <td className="px-6 py-4">
                  <span
                    className={`text-xs font-semibold px-3 py-1 rounded-full ${
                      payment.status === 'processed'
                        ? 'bg-green-100 text-green-800'
                        : payment.status === 'pending'
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {payment.status.charAt(0).toUpperCase() + payment.status.slice(1)}
                  </span>
                </td>
                <td className="px-6 py-4 text-sm">
                  <button
                    onClick={() => onDelete(payment)}
                    className="text-red-600 hover:text-red-900 font-medium"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
