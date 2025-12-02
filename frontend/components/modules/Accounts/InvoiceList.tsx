'use client'

import { Invoice } from '@/types/accounts'

interface InvoiceListProps {
  invoices: Invoice[]
  loading: boolean
  onEdit: (invoice: Invoice) => void
  onDelete: (invoice: Invoice) => void
  onStatusChange: (invoice: Invoice, status: string) => void
}

export default function InvoiceList({ invoices, loading, onEdit, onDelete, onStatusChange }: InvoiceListProps) {
  const statuses = ['draft', 'sent', 'paid', 'overdue', 'cancelled']
  const statusColors: Record<string, string> = {
    draft: 'bg-gray-100 text-gray-800',
    sent: 'bg-blue-100 text-blue-800',
    paid: 'bg-green-100 text-green-800',
    overdue: 'bg-red-100 text-red-800',
    cancelled: 'bg-yellow-100 text-yellow-800',
  }

  if (loading) return <div className="text-center py-8 text-gray-500">Loading invoices...</div>

  const totalAmount = invoices.reduce((sum, i) => sum + i.net_amount, 0)
  const paidAmount = invoices.filter(i => i.status === 'paid').reduce((sum, i) => sum + i.net_amount, 0)

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-gradient-to-br from-blue-50 to-blue-100 rounded-lg p-4 border border-blue-200">
          <p className="text-gray-600 text-xs font-medium">Total Invoices</p>
          <p className="text-2xl font-bold text-blue-600 mt-1">{invoices.length}</p>
        </div>
        <div className="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-4 border border-green-200">
          <p className="text-gray-600 text-xs font-medium">Total Amount</p>
          <p className="text-2xl font-bold text-green-600 mt-1">₹{totalAmount.toLocaleString()}</p>
        </div>
        <div className="bg-gradient-to-br from-purple-50 to-purple-100 rounded-lg p-4 border border-purple-200">
          <p className="text-gray-600 text-xs font-medium">Paid</p>
          <p className="text-2xl font-bold text-purple-600 mt-1">₹{paidAmount.toLocaleString()}</p>
        </div>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Invoice #</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Customer</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Date</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Amount</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Due Date</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Status</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {invoices.map((invoice) => (
              <tr key={invoice.id} className="hover:bg-gray-50 transition">
                <td className="px-6 py-4 text-sm font-medium text-gray-900">{invoice.invoice_number}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{invoice.customer_name || invoice.customer_id}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{new Date(invoice.invoice_date).toLocaleDateString()}</td>
                <td className="px-6 py-4 text-sm font-semibold text-gray-900">₹{invoice.net_amount.toLocaleString()}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{new Date(invoice.due_date).toLocaleDateString()}</td>
                <td className="px-6 py-4">
                  <select
                    value={invoice.status}
                    onChange={(e) => onStatusChange(invoice, e.target.value)}
                    className={`text-xs font-medium px-3 py-1 rounded-full cursor-pointer ${statusColors[invoice.status]}`}
                  >
                    {statuses.map((s) => (
                      <option key={s} value={s}>
                        {s.charAt(0).toUpperCase() + s.slice(1)}
                      </option>
                    ))}
                  </select>
                </td>
                <td className="px-6 py-4 text-sm space-x-2">
                  <button
                    onClick={() => onEdit(invoice)}
                    className="text-blue-600 hover:text-blue-900 font-medium"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(invoice)}
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
