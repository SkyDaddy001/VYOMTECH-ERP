'use client'

import { useState } from 'react'
import { Invoice, InvoiceItem } from '@/types/accounts'

interface InvoiceFormProps {
  invoice?: Invoice | null
  onSubmit: (data: Partial<Invoice>) => Promise<void>
  onCancel: () => void
}

export default function InvoiceFormComponent({ invoice, onSubmit, onCancel }: InvoiceFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<Invoice>>(
    invoice || {
      invoice_number: '',
      customer_id: '',
      customer_name: '',
      invoice_date: new Date().toISOString().split('T')[0],
      due_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      items: [{ description: '', quantity: 1, unit_price: 0, line_total: 0 }],
      discount_amount: 0,
      tax_amount: 0,
      net_amount: 0,
      total_amount: 0,
      status: 'draft',
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
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  const updateItem = (index: number, field: string, value: any) => {
    const updated = [...(formData.items || [])]
    updated[index] = { ...updated[index], [field]: value }
    
    // Recalculate line total for this item
    if (field === 'quantity' || field === 'unit_price') {
      const quantity = field === 'quantity' ? value : updated[index].quantity
      const unit_price = field === 'unit_price' ? value : updated[index].unit_price
      updated[index].line_total = quantity * unit_price
    }
    
    // Recalculate totals
    const subtotal = updated.reduce((sum, item) => sum + (item.line_total || 0), 0)
    
    setFormData(prev => ({
      ...prev,
      items: updated,
      net_amount: subtotal - (prev.discount_amount || 0),
      total_amount: subtotal
    }))
  }

  const addItem = () => {
    const updated = [...(formData.items || [])]
    updated.push({ description: '', quantity: 1, unit_price: 0, line_total: 0 })
    setFormData(prev => ({ ...prev, items: updated }))
  }

  const removeItem = (index: number) => {
    const updated = (formData.items || []).filter((_, i) => i !== index)
    const subtotal = updated.reduce((sum, item) => sum + (item.line_total || 0), 0)
    setFormData(prev => ({
      ...prev,
      items: updated,
      net_amount: subtotal - (prev.discount_amount || 0),
      total_amount: subtotal
    }))
  }

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <span className="text-2xl">📄</span> Invoice
        </h2>
      </div>

      {/* Invoice Header */}
      <div className="grid grid-cols-3 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Invoice Number *</label>
          <input
            type="text"
            required
            value={formData.invoice_number || ''}
            onChange={(e) => handleChange('invoice_number', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="INV-2024-001"
          />
        </div>
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

      {/* Customer Information */}
      <div className="bg-blue-50 rounded-lg p-4 border border-blue-200">
        <h3 className="text-sm font-semibold text-blue-900 mb-3">Bill To</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Customer ID</label>
            <input
              type="text"
              value={formData.customer_id || ''}
              onChange={(e) => handleChange('customer_id', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="CUS-001"
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
              placeholder="Customer name"
            />
          </div>
        </div>
      </div>

      {/* Invoice Items */}
      <div className="space-y-3">
        <div className="flex justify-between items-center">
          <h3 className="text-lg font-semibold text-gray-800">Line Items</h3>
          <button
            type="button"
            onClick={addItem}
            className="px-3 py-1 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 font-medium"
          >
            + Add Item
          </button>
        </div>

        <div className="overflow-x-auto border border-gray-200 rounded-lg">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Description</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Quantity</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Unit Price</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Line Total</th>
                <th className="px-4 py-3 text-center text-sm font-semibold text-gray-700">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {(formData.items || []).map((item, index) => (
                <tr key={index}>
                  <td className="px-4 py-3">
                    <input
                      type="text"
                      value={item.description || ''}
                      onChange={(e) => updateItem(index, 'description', e.target.value)}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-sm"
                      placeholder="Item description"
                    />
                  </td>
                  <td className="px-4 py-3">
                    <input
                      type="number"
                      value={item.quantity || 1}
                      onChange={(e) => updateItem(index, 'quantity', parseFloat(e.target.value))}
                      className="w-20 px-2 py-1 border border-gray-300 rounded text-sm"
                      min="0"
                      step="0.01"
                    />
                  </td>
                  <td className="px-4 py-3">
                    <input
                      type="number"
                      value={item.unit_price || 0}
                      onChange={(e) => updateItem(index, 'unit_price', parseFloat(e.target.value))}
                      className="w-24 px-2 py-1 border border-gray-300 rounded text-sm"
                      min="0"
                      step="0.01"
                    />
                  </td>
                  <td className="px-4 py-3 text-sm font-semibold text-gray-900">
                    ₹{((item.line_total || 0) * 100) / 100}
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

      {/* Totals */}
      <div className="flex justify-end">
        <div className="w-80 space-y-2 bg-gray-50 p-4 rounded-lg border border-gray-200">
          <div className="flex justify-between text-sm">
            <span className="text-gray-600">Subtotal:</span>
            <span className="font-semibold">₹{(formData.total_amount || 0).toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-gray-600">Discount:</span>
            <input
              type="number"
              value={formData.discount_amount || 0}
              onChange={(e) => {
                const discount = parseFloat(e.target.value) || 0
                const subtotal = formData.total_amount || 0
                setFormData(prev => ({
                  ...prev,
                  discount_amount: discount,
                  net_amount: subtotal - discount
                }))
              }}
              className="w-24 px-2 py-1 border border-gray-300 rounded text-sm text-right"
              step="0.01"
            />
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-gray-600">Tax Amount:</span>
            <span className="font-semibold">₹{(formData.tax_amount || 0).toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-lg border-t border-gray-300 pt-2">
            <span className="font-semibold">Net Amount:</span>
            <span className="font-bold text-blue-600">₹{(formData.net_amount || 0).toLocaleString()}</span>
          </div>
        </div>
      </div>

      {/* Notes & Status */}
      <div className="grid grid-cols-2 gap-4">
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
          <input
            type="text"
            value={formData.notes || ''}
            onChange={(e) => handleChange('notes', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Additional notes..."
          />
        </div>
      </div>

      {/* Form Actions */}
      <div className="flex gap-3 pt-4">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium transition-colors"
        >
          {loading ? 'Saving...' : 'Save Invoice'}
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
