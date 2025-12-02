'use client'

import { useState } from 'react'
import { Invoice } from '@/types/accounts'

interface InvoiceFormProps {
  invoice?: Invoice | null
  onSubmit: (data: Partial<Invoice>) => Promise<void>
  onCancel: () => void
}

export default function InvoiceForm({ invoice, onSubmit, onCancel }: InvoiceFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<Invoice>>(
    invoice || {
      invoice_number: `INV-${Date.now().toString().slice(-6)}`,
      customer_id: '',
      invoice_date: new Date().toISOString().split('T')[0],
      due_date: '',
      total_amount: 0,
      tax_amount: 0,
      discount_amount: 0,
      net_amount: 0,
      status: 'draft',
      items: [],
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
    if (['total_amount', 'tax_amount', 'discount_amount'].includes(field)) {
      updated.net_amount = (updated.total_amount || 0) + (updated.tax_amount || 0) - (updated.discount_amount || 0)
    }
    setFormData(updated)
  }

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Invoice Number *</label>
          <input
            type="text"
            required
            disabled
            value={formData.invoice_number || ''}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Customer ID *</label>
          <input
            type="text"
            required
            value={formData.customer_id || ''}
            onChange={(e) => handleChange('customer_id', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="CUST-001"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Invoice Date *</label>
          <input
            type="date"
            required
            value={formData.invoice_date || ''}
            onChange={(e) => handleChange('invoice_date', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Due Date *</label>
          <input
            type="date"
            required
            value={formData.due_date || ''}
            onChange={(e) => handleChange('due_date', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Total Amount ₹ *</label>
          <input
            type="number"
            required
            value={formData.total_amount || 0}
            onChange={(e) => handleChange('total_amount', Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Tax Amount ₹</label>
          <input
            type="number"
            value={formData.tax_amount || 0}
            onChange={(e) => handleChange('tax_amount', Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Discount Amount ₹</label>
          <input
            type="number"
            value={formData.discount_amount || 0}
            onChange={(e) => handleChange('discount_amount', Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="0"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Net Amount ₹</label>
          <input
            type="number"
            disabled
            value={formData.net_amount || 0}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
        <select
          value={formData.status || 'draft'}
          onChange={(e) => handleChange('status', e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="draft">Draft</option>
          <option value="sent">Sent</option>
          <option value="paid">Paid</option>
          <option value="overdue">Overdue</option>
          <option value="cancelled">Cancelled</option>
        </select>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Notes</label>
        <textarea
          value={formData.notes || ''}
          onChange={(e) => handleChange('notes', e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Invoice notes..."
          rows={3}
        />
      </div>

      <div className="flex gap-3">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
        >
          {loading ? 'Saving...' : 'Save Invoice'}
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
