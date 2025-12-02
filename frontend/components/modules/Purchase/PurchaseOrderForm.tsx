'use client'

import { useState } from 'react'
import { PurchaseOrder, PurchaseOrderItem } from '@/types/purchase'

interface PurchaseOrderFormProps {
  purchaseOrder?: PurchaseOrder | null
  onSubmit: (data: Partial<PurchaseOrder>) => Promise<void>
  onCancel: () => void
}

export default function PurchaseOrderForm({ purchaseOrder, onSubmit, onCancel }: PurchaseOrderFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<PurchaseOrder>>(
    purchaseOrder || {
      po_number: '',
      vendor_id: '',
      vendor_name: '',
      po_date: new Date().toISOString().split('T')[0],
      due_date: '',
      delivery_date: '',
      po_status: 'draft',
      items: [],
      subtotal_amount: 0,
      tax_percentage: 18,
      tax_amount: 0,
      shipping_amount: 0,
      discount_amount: 0,
      total_amount: 0,
      payment_status: 'pending',
      notes: '',
    }
  )

  const handleChange = (field: string, value: any) => {
    setFormData(prev => {
      const updated = { ...prev, [field]: value }
      // Auto-calculate totals
      if (field === 'subtotal_amount' || field === 'tax_percentage' || field === 'shipping_amount' || field === 'discount_amount') {
        const subtotal = field === 'subtotal_amount' ? value : updated.subtotal_amount || 0
        const taxPercent = field === 'tax_percentage' ? value : updated.tax_percentage || 0
        const shipping = field === 'shipping_amount' ? value : updated.shipping_amount || 0
        const discount = field === 'discount_amount' ? value : updated.discount_amount || 0

        const tax = (subtotal * taxPercent) / 100
        const total = subtotal + tax + shipping - discount

        return { ...updated, tax_amount: tax, total_amount: total }
      }
      return updated
    })
  }

  const addLineItem = () => {
    setFormData(prev => ({
      ...prev,
      items: [...(prev.items || []), { item_id: '', description: '', quantity: 1, unit_price: 0, amount: 0 }]
    }))
  }

  const updateLineItem = (index: number, field: string, value: any) => {
    const items = [...(formData.items || [])]
    items[index] = { ...items[index], [field]: value }

    if (field === 'quantity' || field === 'unit_price') {
      const quantity = field === 'quantity' ? value : items[index].quantity
      const unit_price = field === 'unit_price' ? value : items[index].unit_price
      items[index].amount = quantity * unit_price
    }

    const subtotal = items.reduce((sum, item) => sum + (item.amount || 0), 0)
    setFormData(prev => ({
      ...prev,
      items,
      subtotal_amount: subtotal,
      tax_amount: (subtotal * (prev.tax_percentage || 0)) / 100,
      total_amount: subtotal + ((subtotal * (prev.tax_percentage || 0)) / 100) + (prev.shipping_amount || 0) - (prev.discount_amount || 0)
    }))
  }

  const removeLineItem = (index: number) => {
    const items = formData.items?.filter((_, i) => i !== index) || []
    const subtotal = items.reduce((sum, item) => sum + (item.amount || 0), 0)
    setFormData(prev => ({
      ...prev,
      items,
      subtotal_amount: subtotal,
      tax_amount: (subtotal * (prev.tax_percentage || 0)) / 100,
      total_amount: subtotal + ((subtotal * (prev.tax_percentage || 0)) / 100) + (prev.shipping_amount || 0) - (prev.discount_amount || 0)
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await onSubmit(formData)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <span className="text-2xl">📦</span> {purchaseOrder ? 'Edit Purchase Order' : 'Create Purchase Order'}
        </h2>
      </div>

      {/* PO Header */}
      <div className="bg-blue-50 rounded-lg p-4 border border-blue-200 grid grid-cols-4 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">PO Number *</label>
          <input
            type="text"
            required
            value={formData.po_number || ''}
            onChange={(e) => handleChange('po_number', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="PO-001"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Vendor *</label>
          <input
            type="text"
            required
            value={formData.vendor_name || ''}
            onChange={(e) => handleChange('vendor_name', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Vendor name"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">PO Date *</label>
          <input
            type="date"
            required
            value={formData.po_date || ''}
            onChange={(e) => handleChange('po_date', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
          <select
            value={formData.po_status || 'draft'}
            onChange={(e) => handleChange('po_status', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="draft">Draft</option>
            <option value="sent">Sent</option>
            <option value="confirmed">Confirmed</option>
            <option value="partial_received">Partial Received</option>
            <option value="received">Received</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>
      </div>

      {/* Dates */}
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Due Date</label>
          <input
            type="date"
            value={formData.due_date || ''}
            onChange={(e) => handleChange('due_date', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Delivery Date</label>
          <input
            type="date"
            value={formData.delivery_date || ''}
            onChange={(e) => handleChange('delivery_date', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      {/* Line Items */}
      <div className="bg-gray-50 rounded-lg p-4 border border-gray-200">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-sm font-semibold text-gray-800">Line Items</h3>
          <button
            type="button"
            onClick={addLineItem}
            className="bg-green-600 text-white px-3 py-1 rounded text-sm hover:bg-green-700 transition-colors"
          >
            + Add Item
          </button>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="bg-gray-200 border-b-2 border-gray-300">
                <th className="px-2 py-2 text-left">Item ID</th>
                <th className="px-2 py-2 text-left">Description</th>
                <th className="px-2 py-2 text-right">Quantity</th>
                <th className="px-2 py-2 text-right">Unit Price</th>
                <th className="px-2 py-2 text-right">Amount</th>
                <th className="px-2 py-2 text-center">Action</th>
              </tr>
            </thead>
            <tbody>
              {formData.items?.map((item, index) => (
                <tr key={index} className="border-b border-gray-300 hover:bg-white transition-colors">
                  <td className="px-2 py-2">
                    <input
                      type="text"
                      value={item.item_id || ''}
                      onChange={(e) => updateLineItem(index, 'item_id', e.target.value)}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-xs"
                      placeholder="Item"
                    />
                  </td>
                  <td className="px-2 py-2">
                    <input
                      type="text"
                      value={item.description || ''}
                      onChange={(e) => updateLineItem(index, 'description', e.target.value)}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-xs"
                      placeholder="Description"
                    />
                  </td>
                  <td className="px-2 py-2">
                    <input
                      type="number"
                      value={item.quantity || 1}
                      onChange={(e) => updateLineItem(index, 'quantity', parseFloat(e.target.value))}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-xs text-right"
                      step="0.01"
                    />
                  </td>
                  <td className="px-2 py-2">
                    <input
                      type="number"
                      value={item.unit_price || 0}
                      onChange={(e) => updateLineItem(index, 'unit_price', parseFloat(e.target.value))}
                      className="w-full px-2 py-1 border border-gray-300 rounded text-xs text-right"
                      step="0.01"
                    />
                  </td>
                  <td className="px-2 py-2 text-right text-gray-700 font-semibold">
                    ₹{((item.amount || 0) * 100) / 100}
                  </td>
                  <td className="px-2 py-2 text-center">
                    <button
                      type="button"
                      onClick={() => removeLineItem(index)}
                      className="bg-red-600 text-white px-2 py-1 rounded text-xs hover:bg-red-700 transition-colors"
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
      <div className="bg-yellow-50 rounded-lg p-4 border border-yellow-200 grid grid-cols-2 gap-8">
        <div>
          <div className="text-sm text-gray-600 mb-2">Notes</div>
          <textarea
            value={formData.notes || ''}
            onChange={(e) => handleChange('notes', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500 text-sm"
            placeholder="Order notes..."
            rows={3}
          />
        </div>

        <div className="space-y-2">
          <div className="flex justify-between text-sm">
            <span className="text-gray-700">Subtotal:</span>
            <span className="font-semibold">₹{((formData.subtotal_amount || 0) * 100) / 100}</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-gray-700">Tax ({formData.tax_percentage || 0}%):</span>
            <span className="font-semibold">₹{((formData.tax_amount || 0) * 100) / 100}</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-gray-700">Shipping:</span>
            <input
              type="number"
              value={formData.shipping_amount || 0}
              onChange={(e) => handleChange('shipping_amount', parseFloat(e.target.value))}
              className="w-20 px-2 py-1 border border-gray-300 rounded text-sm text-right"
              step="0.01"
            />
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-gray-700">Discount:</span>
            <input
              type="number"
              value={formData.discount_amount || 0}
              onChange={(e) => handleChange('discount_amount', parseFloat(e.target.value))}
              className="w-20 px-2 py-1 border border-gray-300 rounded text-sm text-right"
              step="0.01"
            />
          </div>
          <div className="flex justify-between text-lg font-bold pt-2 border-t-2 border-gray-300 mt-2">
            <span>Total:</span>
            <span className="text-blue-600">₹{((formData.total_amount || 0) * 100) / 100}</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-gray-700">Payment Status:</span>
            <select
              value={formData.payment_status || 'pending'}
              onChange={(e) => handleChange('payment_status', e.target.value)}
              className="px-2 py-1 border border-gray-300 rounded text-sm"
            >
              <option value="pending">Pending</option>
              <option value="partial">Partial</option>
              <option value="paid">Paid</option>
            </select>
          </div>
        </div>
      </div>

      {/* Form Actions */}
      <div className="flex gap-3 pt-4">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium transition-colors"
        >
          {loading ? 'Saving...' : purchaseOrder ? 'Update Purchase Order' : 'Create Purchase Order'}
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
