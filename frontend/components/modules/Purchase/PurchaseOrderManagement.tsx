'use client'

import { useState, useEffect } from 'react'
import axios from 'axios'
import toast from 'react-hot-toast'
import { formatDateToDDMMMYYYY, formatDateToInput } from '@/lib/dateFormat'

interface PurchaseOrder {
  id: string
  po_number: string
  vendor_id: string
  vendor_name?: string
  po_date: string
  delivery_date?: string
  total_amount: number
  status: string
}

export default function PurchaseOrderManagement() {
  const [orders, setOrders] = useState<PurchaseOrder[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [vendors, setVendors] = useState<{ id: string; name: string }[]>([])
  const [formData, setFormData] = useState({
    vendor_id: '',
    po_date: new Date().toISOString().split('T')[0],
    delivery_date: '',
    total_amount: 0,
  })

  useEffect(() => {
    fetchOrders()
    fetchVendors()
  }, [])

  const fetchOrders = async () => {
    try {
      const response = await axios.get('/api/v1/purchase/orders')
      setOrders(response.data || [])
    } catch (error) {
      toast.error('Failed to fetch purchase orders')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const fetchVendors = async () => {
    try {
      const response = await axios.get('/api/v1/purchase/vendors')
      setVendors(response.data || [])
    } catch (error) {
      console.error('Failed to fetch vendors:', error)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await axios.post('/api/v1/purchase/orders', formData)
      toast.success('Purchase order created successfully')
      setShowForm(false)
      setFormData({
        vendor_id: '',
        po_date: new Date().toISOString().split('T')[0],
        delivery_date: '',
        total_amount: 0,
      })
      fetchOrders()
    } catch (error) {
      toast.error('Failed to create purchase order')
      console.error(error)
    }
  }

  const getStatusColor = (status: string) => {
    const colors: { [key: string]: string } = {
      draft: 'bg-gray-100 text-gray-800',
      submitted: 'bg-blue-100 text-blue-800',
      approved: 'bg-green-100 text-green-800',
      received: 'bg-purple-100 text-purple-800',
      completed: 'bg-emerald-100 text-emerald-800',
      cancelled: 'bg-red-100 text-red-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  if (loading) {
    return <div className="p-6 text-center text-gray-500">Loading purchase orders...</div>
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Purchase Orders</h2>
        <button
          onClick={() => setShowForm(true)}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition"
        >
          + New PO
        </button>
      </div>

      {/* Form Modal */}
      {showForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-8 max-w-md w-full">
            <h3 className="text-xl font-bold mb-4">Create Purchase Order</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Vendor</label>
                <select
                  value={formData.vendor_id}
                  onChange={(e) => setFormData({ ...formData, vendor_id: e.target.value })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                >
                  <option value="">Select Vendor</option>
                  {vendors.map((vendor) => (
                    <option key={vendor.id} value={vendor.id}>
                      {vendor.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">PO Date</label>
                <input
                  type="date"
                  value={formData.po_date}
                  onChange={(e) => setFormData({ ...formData, po_date: e.target.value })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Delivery Date</label>
                <input
                  type="date"
                  value={formData.delivery_date}
                  onChange={(e) => setFormData({ ...formData, delivery_date: e.target.value })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Total Amount</label>
                <input
                  type="number"
                  value={formData.total_amount}
                  onChange={(e) => setFormData({ ...formData, total_amount: parseFloat(e.target.value) })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
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

      {/* Purchase Orders Table */}
      <div className="overflow-x-auto">
        <table className="w-full border-collapse">
          <thead>
            <tr className="bg-gray-100 border-b">
              <th className="px-4 py-3 text-left font-semibold text-gray-700">PO Number</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Vendor</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">PO Date</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Delivery Date</th>
              <th className="px-4 py-3 text-right font-semibold text-gray-700">Amount</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Status</th>
              <th className="px-4 py-3 text-center font-semibold text-gray-700">Actions</th>
            </tr>
          </thead>
          <tbody>
            {orders.length === 0 ? (
              <tr>
                <td colSpan={7} className="px-4 py-3 text-center text-gray-500">
                  No purchase orders found
                </td>
              </tr>
            ) : (
              orders.map((order) => (
                <tr key={order.id} className="border-b hover:bg-gray-50">
                  <td className="px-4 py-3 font-mono font-semibold">{order.po_number}</td>
                  <td className="px-4 py-3">{order.vendor_name || 'N/A'}</td>
                  <td className="px-4 py-3 text-sm">{formatDateToDDMMMYYYY(order.po_date)}</td>
                  <td className="px-4 py-3 text-sm">
                    {order.delivery_date ? formatDateToDDMMMYYYY(order.delivery_date) : 'N/A'}
                  </td>
                  <td className="px-4 py-3 text-right font-semibold">
                    ${order.total_amount.toLocaleString('en-US', { minimumFractionDigits: 2 })}
                  </td>
                  <td className="px-4 py-3">
                    <span className={`inline-block px-3 py-1 rounded text-sm font-medium ${getStatusColor(order.status)}`}>
                      {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-center space-x-2">
                    <button className="px-3 py-1 bg-blue-100 text-blue-600 rounded hover:bg-blue-200 text-sm">
                      View
                    </button>
                    <button className="px-3 py-1 bg-green-100 text-green-600 rounded hover:bg-green-200 text-sm">
                      GRN
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
