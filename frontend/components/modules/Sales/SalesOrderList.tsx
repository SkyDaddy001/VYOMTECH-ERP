'use client'

import { SalesOrder } from '@/types/sales'

interface SalesOrderListProps {
  orders: SalesOrder[]
  loading: boolean
  onEdit: (order: SalesOrder) => void
  onDelete: (order: SalesOrder) => void
  onStatusChange: (order: SalesOrder, status: string) => void
}

export default function SalesOrderList({ orders, loading, onEdit, onDelete, onStatusChange }: SalesOrderListProps) {
  const stages = ['inquiry', 'quote', 'booking', 'agreement', 'completed', 'cancelled']
  const stageColors: Record<string, string> = {
    inquiry: 'bg-blue-100 text-blue-800',
    quote: 'bg-yellow-100 text-yellow-800',
    booking: 'bg-purple-100 text-purple-800',
    agreement: 'bg-indigo-100 text-indigo-800',
    completed: 'bg-green-100 text-green-800',
    cancelled: 'bg-red-100 text-red-800',
  }

  if (loading) return <div className="text-center py-8 text-gray-500">Loading bookings...</div>

  const totalAmount = orders.reduce((sum, o) => sum + o.net_amount, 0)
  const averageBooking = orders.length > 0 ? totalAmount / orders.length : 0

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-gradient-to-br from-blue-50 to-blue-100 rounded-lg p-4 border border-blue-200">
          <p className="text-gray-600 text-xs font-medium">Total Bookings</p>
          <p className="text-2xl font-bold text-blue-600 mt-1">{orders.length}</p>
        </div>
        <div className="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-4 border border-green-200">
          <p className="text-gray-600 text-xs font-medium">Total Booking Value</p>
          <p className="text-2xl font-bold text-green-600 mt-1">₹{totalAmount.toLocaleString()}</p>
        </div>
        <div className="bg-gradient-to-br from-purple-50 to-purple-100 rounded-lg p-4 border border-purple-200">
          <p className="text-gray-600 text-xs font-medium">Avg Booking Value</p>
          <p className="text-2xl font-bold text-purple-600 mt-1">₹{Math.round(averageBooking).toLocaleString()}</p>
        </div>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Booking #</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Customer</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Booking Date</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Amount</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Stage</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {orders.map((order) => (
              <tr key={order.id} className="hover:bg-gray-50 transition">
                <td className="px-6 py-4 text-sm font-medium text-gray-900">{order.booking_number}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{order.customer_name || order.customer_id}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{new Date(order.booking_date).toLocaleDateString()}</td>
                <td className="px-6 py-4 text-sm font-semibold text-gray-900">₹{order.net_amount.toLocaleString()}</td>
                <td className="px-6 py-4">
                  <select
                    value={order.booking_stage}
                    onChange={(e) => onStatusChange(order, e.target.value)}
                    className={`text-xs font-medium px-3 py-1 rounded-full cursor-pointer ${stageColors[order.booking_stage]}`}
                  >
                    {stages.map((s) => (
                      <option key={s} value={s}>
                        {s.charAt(0).toUpperCase() + s.slice(1)}
                      </option>
                    ))}
                  </select>
                </td>
                <td className="px-6 py-4 text-sm space-x-2">
                  <button
                    onClick={() => onEdit(order)}
                    className="text-blue-600 hover:text-blue-900 font-medium"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(order)}
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
