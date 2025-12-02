'use client'

import { useState } from 'react'
import { GoodsReceiptNote, GRNLineItem } from '@/types/purchase'

interface GRNFormProps {
  grn?: GoodsReceiptNote | null
  onSubmit: (data: Partial<GoodsReceiptNote>) => Promise<void>
  onCancel: () => void
}

export default function GRNForm({ grn, onSubmit, onCancel }: GRNFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<GoodsReceiptNote>>(
    grn || {
      grn_number: '',
      po_number: '',
      vendor_name: '',
      grn_date: new Date().toISOString().split('T')[0],
      received_date: new Date().toISOString().split('T')[0],
      grn_status: 'pending_inspection',
      items: [],
      notes: '',
      warehouse_location: '',
      received_by: '',
      inspected_by: '',
    }
  )

  const handleChange = (field: string, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  const addLineItem = () => {
    setFormData(prev => ({
      ...prev,
      items: [...(prev.items || []), { po_item_id: '', description: '', quantity_ordered: 0, quantity_received: 0, quality_status: 'accepted', remarks: '' }]
    }))
  }

  const updateLineItem = (index: number, field: string, value: any) => {
    const items = [...(formData.items || [])]
    items[index] = { ...items[index], [field]: value }
    setFormData(prev => ({ ...prev, items }))
  }

  const removeLineItem = (index: number) => {
    const items = formData.items?.filter((_, i) => i !== index) || []
    setFormData(prev => ({ ...prev, items }))
  }

  const getVarianceStatus = () => {
    const items = formData.items || []
    const totalOrdered = items.reduce((sum, item) => sum + (item.quantity_ordered || 0), 0)
    const totalReceived = items.reduce((sum, item) => sum + (item.quantity_received || 0), 0)

    if (totalReceived === totalOrdered) return { status: 'Complete', color: 'green' }
    if (totalReceived > 0 && totalReceived < totalOrdered) return { status: 'Partial', color: 'yellow' }
    if (totalReceived > totalOrdered) return { status: 'Over', color: 'red' }
    return { status: 'None', color: 'gray' }
  }

  const variance = getVarianceStatus()

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
          <span className="text-2xl">📋</span> {grn ? 'Edit GRN' : 'Create Goods Receipt Note'}
        </h2>
      </div>

      {/* GRN Header */}
      <div className="bg-blue-50 rounded-lg p-4 border border-blue-200 grid grid-cols-4 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">GRN Number *</label>
          <input
            type="text"
            required
            value={formData.grn_number || ''}
            onChange={(e) => handleChange('grn_number', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="GRN-001"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">PO Number</label>
          <input
            type="text"
            value={formData.po_number || ''}
            onChange={(e) => handleChange('po_number', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="PO-001"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Vendor Name</label>
          <input
            type="text"
            value={formData.vendor_name || ''}
            onChange={(e) => handleChange('vendor_name', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Vendor name"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
          <select
            value={formData.grn_status || 'pending_inspection'}
            onChange={(e) => handleChange('grn_status', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="pending_inspection">Pending Inspection</option>
            <option value="inspected">Inspected</option>
            <option value="accepted">Accepted</option>
            <option value="rejected">Rejected</option>
            <option value="partial_accepted">Partial Accepted</option>
          </select>
        </div>
      </div>

      {/* Dates and Personnel */}
      <div className="grid grid-cols-2 gap-4">
        <div className="bg-green-50 rounded-lg p-4 border border-green-200 space-y-3">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">GRN Date *</label>
            <input
              type="date"
              required
              value={formData.grn_date || ''}
              onChange={(e) => handleChange('grn_date', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Received Date</label>
            <input
              type="date"
              value={formData.received_date || ''}
              onChange={(e) => handleChange('received_date', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>
        </div>

        <div className="bg-purple-50 rounded-lg p-4 border border-purple-200 space-y-3">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Received By</label>
            <input
              type="text"
              value={formData.received_by || ''}
              onChange={(e) => handleChange('received_by', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
              placeholder="Name"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Inspected By</label>
            <input
              type="text"
              value={formData.inspected_by || ''}
              onChange={(e) => handleChange('inspected_by', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
              placeholder="Name"
            />
          </div>
        </div>
      </div>

      {/* Warehouse & Variance */}
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Warehouse Location</label>
          <input
            type="text"
            value={formData.warehouse_location || ''}
            onChange={(e) => handleChange('warehouse_location', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Warehouse area"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Variance Status</label>
          <div className={`px-3 py-2 rounded-lg font-semibold text-center bg-${variance.color}-100 text-${variance.color}-900 border border-${variance.color}-300`}>
            {variance.status}
          </div>
        </div>
      </div>

      {/* Line Items */}
      <div className="bg-gray-50 rounded-lg p-4 border border-gray-200">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-sm font-semibold text-gray-800">Received Items</h3>
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
                <th className="px-2 py-2 text-left">PO Item ID</th>
                <th className="px-2 py-2 text-left">Description</th>
                <th className="px-2 py-2 text-right">Ordered</th>
                <th className="px-2 py-2 text-right">Received</th>
                <th className="px-2 py-2 text-center">Variance</th>
                <th className="px-2 py-2 text-center">Quality</th>
                <th className="px-2 py-2 text-left">Remarks</th>
                <th className="px-2 py-2 text-center">Action</th>
              </tr>
            </thead>
            <tbody>
              {formData.items?.map((item, index) => {
                const variance_qty = (item.quantity_received || 0) - (item.quantity_ordered || 0)
                return (
                  <tr key={index} className="border-b border-gray-300 hover:bg-white transition-colors">
                    <td className="px-2 py-2">
                      <input
                        type="text"
                        value={item.po_item_id || ''}
                        onChange={(e) => updateLineItem(index, 'po_item_id', e.target.value)}
                        className="w-full px-2 py-1 border border-gray-300 rounded text-xs"
                        placeholder="Item ID"
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
                        value={item.quantity_ordered || 0}
                        onChange={(e) => updateLineItem(index, 'quantity_ordered', parseFloat(e.target.value))}
                        className="w-full px-2 py-1 border border-gray-300 rounded text-xs text-right"
                        step="0.01"
                      />
                    </td>
                    <td className="px-2 py-2">
                      <input
                        type="number"
                        value={item.quantity_received || 0}
                        onChange={(e) => updateLineItem(index, 'quantity_received', parseFloat(e.target.value))}
                        className="w-full px-2 py-1 border border-gray-300 rounded text-xs text-right"
                        step="0.01"
                      />
                    </td>
                    <td className={`px-2 py-2 text-center font-semibold ${variance_qty === 0 ? 'text-green-600' : variance_qty > 0 ? 'text-red-600' : 'text-yellow-600'}`}>
                      {variance_qty > 0 ? '+' : ''}{variance_qty}
                    </td>
                    <td className="px-2 py-2">
                      <select
                        value={item.quality_status || 'accepted'}
                        onChange={(e) => updateLineItem(index, 'quality_status', e.target.value)}
                        className="w-full px-2 py-1 border border-gray-300 rounded text-xs"
                      >
                        <option value="accepted">Accepted</option>
                        <option value="rejected">Rejected</option>
                        <option value="partial">Partial</option>
                      </select>
                    </td>
                    <td className="px-2 py-2">
                      <input
                        type="text"
                        value={item.remarks || ''}
                        onChange={(e) => updateLineItem(index, 'remarks', e.target.value)}
                        className="w-full px-2 py-1 border border-gray-300 rounded text-xs"
                        placeholder="Remarks"
                      />
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
                )
              })}
            </tbody>
          </table>
        </div>
      </div>

      {/* Notes */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Notes</label>
        <textarea
          value={formData.notes || ''}
          onChange={(e) => handleChange('notes', e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="GRN remarks..."
          rows={3}
        />
      </div>

      {/* Form Actions */}
      <div className="flex gap-3 pt-4">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium transition-colors"
        >
          {loading ? 'Saving...' : grn ? 'Update GRN' : 'Create GRN'}
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
