'use client'

import { useState, useEffect } from 'react'
import axios from 'axios'
import toast from 'react-hot-toast'
import { formatDateToDDMMMYYYY, formatDateToInput } from '@/lib/dateFormat'

interface GRN {
  id: string
  grn_number: string
  po_id: string
  po_number?: string
  received_date: string
  total_received: number
  inspection_status: string
  notes?: string
}

export default function GRNManagement() {
  const [grns, setGrns] = useState<GRN[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [purchaseOrders, setPurchaseOrders] = useState<{ id: string; po_number: string }[]>([])
  const [formData, setFormData] = useState({
    po_id: '',
    received_date: new Date().toISOString().split('T')[0],
    total_received: 0,
    notes: '',
  })

  useEffect(() => {
    fetchGRNs()
    fetchPurchaseOrders()
  }, [])

  const fetchGRNs = async () => {
    try {
      const response = await axios.get('/api/v1/purchase/grn')
      setGrns(response.data || [])
    } catch (error) {
      toast.error('Failed to fetch GRN records')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const fetchPurchaseOrders = async () => {
    try {
      const response = await axios.get('/api/v1/purchase/orders')
      setPurchaseOrders(response.data || [])
    } catch (error) {
      console.error('Failed to fetch purchase orders:', error)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await axios.post('/api/v1/purchase/grn', formData)
      toast.success('GRN created successfully')
      setShowForm(false)
      setFormData({
        po_id: '',
        received_date: new Date().toISOString().split('T')[0],
        total_received: 0,
        notes: '',
      })
      fetchGRNs()
    } catch (error) {
      toast.error('Failed to create GRN')
      console.error(error)
    }
  }

  const getInspectionStatusColor = (status: string) => {
    const colors: { [key: string]: string } = {
      pending: 'bg-yellow-100 text-yellow-800',
      accepted: 'bg-green-100 text-green-800',
      rejected: 'bg-red-100 text-red-800',
      partial: 'bg-orange-100 text-orange-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  if (loading) {
    return <div className="p-6 text-center text-gray-500">Loading GRN records...</div>
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Goods Receipt Notes (GRN) / Material Receipt Notes (MRN)</h2>
        <button
          onClick={() => setShowForm(true)}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition"
        >
          + Create GRN/MRN
        </button>
      </div>

      {/* Info Box */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <p className="text-sm text-blue-800">
          <strong>GRN/MRN Management:</strong> Log material receipts with quality inspection. Track all incoming goods and perform quality checks.
        </p>
      </div>

      {/* Form Modal */}
      {showForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-8 max-w-md w-full">
            <h3 className="text-xl font-bold mb-4">Create GRN/MRN</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Purchase Order</label>
                <select
                  value={formData.po_id}
                  onChange={(e) => setFormData({ ...formData, po_id: e.target.value })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                >
                  <option value="">Select PO</option>
                  {purchaseOrders.map((po) => (
                    <option key={po.id} value={po.id}>
                      {po.po_number}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Received Date</label>
                <input
                  type="date"
                  value={formData.received_date}
                  onChange={(e) => setFormData({ ...formData, received_date: e.target.value })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Total Received Quantity</label>
                <input
                  type="number"
                  value={formData.total_received}
                  onChange={(e) => setFormData({ ...formData, total_received: parseFloat(e.target.value) })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
                <textarea
                  value={formData.notes}
                  onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                  rows={3}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Add any notes about the receipt..."
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowForm(false)}
                  className="flex-1 px-4 py-2 border rounded-lg hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                >
                  Create
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* GRN/MRN Table */}
      <div className="overflow-x-auto">
        <table className="w-full border-collapse">
          <thead>
            <tr className="bg-gray-100 border-b">
              <th className="px-4 py-3 text-left font-semibold text-gray-700">GRN Number</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">PO Number</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Received Date</th>
              <th className="px-4 py-3 text-right font-semibold text-gray-700">Qty Received</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">QC Status</th>
              <th className="px-4 py-3 text-center font-semibold text-gray-700">Actions</th>
            </tr>
          </thead>
          <tbody>
            {grns.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-4 py-3 text-center text-gray-500">
                  No GRN/MRN records found
                </td>
              </tr>
            ) : (
              grns.map((grn) => (
                <tr key={grn.id} className="border-b hover:bg-gray-50">
                  <td className="px-4 py-3 font-mono font-semibold">{grn.grn_number}</td>
                  <td className="px-4 py-3 font-mono">{grn.po_number || 'N/A'}</td>
                  <td className="px-4 py-3 text-sm">{formatDateToDDMMMYYYY(grn.received_date)}</td>
                  <td className="px-4 py-3 text-right font-semibold">{grn.total_received}</td>
                  <td className="px-4 py-3">
                    <span
                      className={`inline-block px-3 py-1 rounded text-sm font-medium ${getInspectionStatusColor(
                        grn.inspection_status
                      )}`}
                    >
                      {grn.inspection_status.charAt(0).toUpperCase() + grn.inspection_status.slice(1)}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-center space-x-2">
                    <button className="px-3 py-1 bg-blue-100 text-blue-600 rounded hover:bg-blue-200 text-sm">
                      Inspect
                    </button>
                    <button className="px-3 py-1 bg-purple-100 text-purple-600 rounded hover:bg-purple-200 text-sm">
                      Accept
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
