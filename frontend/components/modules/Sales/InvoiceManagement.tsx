'use client'

import { useState, useEffect } from 'react'
import { Plus, Eye, Trash2, Search, DollarSign, Calendar, FileText, Download } from 'lucide-react'
import { formatDateToDDMMMYYYY, formatDateToInput } from '@/lib/dateFormat'

interface InvoiceItem {
  id: string
  item_name: string
  description: string
  quantity: number
  unit_price: number
  cgst_percent: number
  cgst_amount: number
  sgst_percent: number
  sgst_amount: number
  igst_percent: number
  igst_amount: number
  line_total: number
}

interface Invoice {
  id: string
  invoice_code: string
  invoice_number: string
  order_id: string
  customer_id: string
  customer_name: string
  invoice_date: string
  due_date: string
  payment_status: string
  sub_total: number
  cgst_total: number
  sgst_total: number
  igst_total: number
  total_tax: number
  total_amount: number
  discount_percent: number
  discount_amount: number
  paid_amount: number
  items: InvoiceItem[]
  notes: string
  ar_posting_status: string
  gl_reference_number: string
  created_at: string
}

interface SalesOrder {
  id: string
  order_code: string
  customer_name: string
  total_amount: number
  items: any[]
}

export function InvoiceManagement() {
  const [invoices, setInvoices] = useState<Invoice[]>([])
  const [filteredInvoices, setFilteredInvoices] = useState<Invoice[]>([])
  const [orders, setOrders] = useState<SalesOrder[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState('all')
  const [showForm, setShowForm] = useState(false)
  const [viewingId, setViewingId] = useState<string | null>(null)
  const [selectedOrder, setSelectedOrder] = useState<SalesOrder | null>(null)

  const [formData, setFormData] = useState({
    order_id: '',
    invoice_date: new Date().toISOString().split('T')[0],
    due_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    discount_percent: 0,
    notes: '',
  })

  useEffect(() => {
    fetchInvoices()
    fetchOrders()
  }, [])

  useEffect(() => {
    filterInvoices()
  }, [invoices, searchTerm, filterStatus])

  const fetchInvoices = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/v1/sales/invoices', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      })
      if (response.ok) {
        const data = await response.json()
        setInvoices(data.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch invoices:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchOrders = async () => {
    try {
      const response = await fetch('/api/v1/sales/orders', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      })
      if (response.ok) {
        const data = await response.json()
        setOrders(data.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch orders:', error)
    }
  }

  const filterInvoices = () => {
    let filtered = invoices

    if (searchTerm) {
      filtered = filtered.filter(
        i =>
          i.invoice_code.toLowerCase().includes(searchTerm.toLowerCase()) ||
          i.customer_name.toLowerCase().includes(searchTerm.toLowerCase())
      )
    }

    if (filterStatus !== 'all') {
      filtered = filtered.filter(i => i.payment_status === filterStatus)
    }

    setFilteredInvoices(filtered)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      const response = await fetch('/api/v1/sales/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
          'X-User-ID': localStorage.getItem('userId') || '',
        },
        body: JSON.stringify(formData),
      })

      if (response.ok) {
        setFormData({
          order_id: '',
          invoice_date: new Date().toISOString().split('T')[0],
          due_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
          discount_percent: 0,
          notes: '',
        })
        setShowForm(false)
        setSelectedOrder(null)
        fetchInvoices()
      }
    } catch (error) {
      console.error('Failed to create invoice:', error)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this invoice?')) return

    try {
      const response = await fetch(`/api/v1/sales/invoices/${id}`, {
        method: 'DELETE',
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      })

      if (response.ok) {
        fetchInvoices()
      }
    } catch (error) {
      console.error('Failed to delete invoice:', error)
    }
  }

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      unpaid: 'bg-red-100 text-red-800',
      partially_paid: 'bg-yellow-100 text-yellow-800',
      paid: 'bg-green-100 text-green-800',
      cancelled: 'bg-gray-100 text-gray-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const getARStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'bg-orange-100 text-orange-800',
      posted: 'bg-green-100 text-green-800',
      failed: 'bg-red-100 text-red-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const viewing = invoices.find(i => i.id === viewingId)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-900">Invoices</h2>
        <button
          onClick={() => {
            setShowForm(!showForm)
            setSelectedOrder(null)
            setFormData({
              order_id: '',
              invoice_date: new Date().toISOString().split('T')[0],
              due_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
              discount_percent: 0,
              notes: '',
            })
          }}
          className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition"
        >
          <Plus size={18} />
          New Invoice
        </button>
      </div>

      {/* Form */}
      {showForm && (
        <div className="bg-white p-6 rounded-lg shadow-md border border-gray-200">
          <h3 className="text-lg font-semibold mb-4 text-gray-900">Create Invoice from Order</h3>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Select Order</label>
                <select
                  value={formData.order_id}
                  onChange={e => {
                    setFormData({ ...formData, order_id: e.target.value })
                    const order = orders.find(o => o.id === e.target.value)
                    setSelectedOrder(order || null)
                  }}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Select Order</option>
                  {orders
                    .filter(o => o.id)
                    .map(o => (
                      <option key={o.id} value={o.id}>
                        {o.order_code} - {o.customer_name}
                      </option>
                    ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Invoice Date</label>
                <input
                  type="date"
                  value={formData.invoice_date}
                  onChange={e => setFormData({ ...formData, invoice_date: e.target.value })}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Due Date</label>
                <input
                  type="date"
                  value={formData.due_date}
                  onChange={e => setFormData({ ...formData, due_date: e.target.value })}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {selectedOrder && (
              <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                <h4 className="font-semibold text-gray-900 mb-3">Order Items</h4>
                <div className="space-y-2">
                  {selectedOrder.items?.map((item, idx) => (
                    <div key={idx} className="flex justify-between text-sm text-gray-700 pb-2 border-b">
                      <span>
                        {item.item_name} x {item.quantity}
                      </span>
                      <span>₹{(item.quantity * item.unit_price).toFixed(2)}</span>
                    </div>
                  ))}
                  <div className="flex justify-between text-lg font-semibold text-blue-600 pt-2">
                    <span>Total Order Value</span>
                    <span>₹{selectedOrder.total_amount?.toFixed(2)}</span>
                  </div>
                </div>
              </div>
            )}

            <textarea
              placeholder="Invoice notes..."
              value={formData.notes}
              onChange={e => setFormData({ ...formData, notes: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
              rows={3}
            />

            <div className="flex gap-2">
              <button
                type="submit"
                className="flex-1 bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition"
              >
                Create Invoice
              </button>
              <button
                type="button"
                onClick={() => {
                  setShowForm(false)
                  setSelectedOrder(null)
                }}
                className="flex-1 bg-gray-400 text-white px-4 py-2 rounded-lg hover:bg-gray-500 transition"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Filters */}
      <div className="flex flex-col md:flex-row gap-4">
        <div className="flex-1 relative">
          <Search className="absolute left-3 top-3 text-gray-400" size={18} />
          <input
            type="text"
            placeholder="Search invoices..."
            value={searchTerm}
            onChange={e => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <select
          value={filterStatus}
          onChange={e => setFilterStatus(e.target.value)}
          className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="all">All Status</option>
          <option value="unpaid">Unpaid</option>
          <option value="partially_paid">Partially Paid</option>
          <option value="paid">Paid</option>
          <option value="cancelled">Cancelled</option>
        </select>
      </div>

      {/* View Invoice Detail */}
      {viewing && (
        <div className="bg-white p-6 rounded-lg shadow-md border border-gray-200">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-gray-900">Invoice {viewing.invoice_code}</h3>
            <button
              onClick={() => setViewingId(null)}
              className="text-gray-500 hover:text-gray-700 text-xl"
            >
              ✕
            </button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <div className="text-sm text-gray-600">Customer</div>
              <div className="text-lg font-semibold text-gray-900">{viewing.customer_name}</div>
            </div>
            <div>
              <div className="text-sm text-gray-600">Invoice Date</div>
              <div className="text-lg font-semibold text-gray-900">{formatDateToDDMMMYYYY(viewing.invoice_date)}</div>
            </div>
            <div>
              <div className="text-sm text-gray-600">Due Date</div>
              <div className="text-lg font-semibold text-gray-900">{formatDateToDDMMMYYYY(viewing.due_date)}</div>
            </div>
            <div>
              <div className="text-sm text-gray-600">Payment Status</div>
              <span className={`inline-block px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(viewing.payment_status)}`}>
                {viewing.payment_status}
              </span>
            </div>
          </div>

          <div className="border-t mb-4 pt-4">
            <h4 className="font-semibold text-gray-900 mb-3">Tax Breakdown</h4>
            <div className="grid grid-cols-3 md:grid-cols-6 gap-3">
              <div className="bg-blue-50 p-3 rounded">
                <div className="text-xs text-gray-600">Subtotal</div>
                <div className="font-semibold text-blue-600">₹{viewing.sub_total.toFixed(2)}</div>
              </div>
              <div className="bg-orange-50 p-3 rounded">
                <div className="text-xs text-gray-600">CGST</div>
                <div className="font-semibold text-orange-600">₹{viewing.cgst_total.toFixed(2)}</div>
              </div>
              <div className="bg-orange-50 p-3 rounded">
                <div className="text-xs text-gray-600">SGST</div>
                <div className="font-semibold text-orange-600">₹{viewing.sgst_total.toFixed(2)}</div>
              </div>
              <div className="bg-orange-50 p-3 rounded">
                <div className="text-xs text-gray-600">IGST</div>
                <div className="font-semibold text-orange-600">₹{viewing.igst_total.toFixed(2)}</div>
              </div>
              <div className="bg-green-50 p-3 rounded">
                <div className="text-xs text-gray-600">Paid</div>
                <div className="font-semibold text-green-600">₹{viewing.paid_amount.toFixed(2)}</div>
              </div>
              <div className="bg-red-50 p-3 rounded">
                <div className="text-xs text-gray-600">Total</div>
                <div className="font-semibold text-red-600">₹{viewing.total_amount.toFixed(2)}</div>
              </div>
            </div>
          </div>

          {viewing.gl_reference_number && (
            <div className="bg-purple-50 p-3 rounded border border-purple-200 mb-4">
              <div className="flex items-center justify-between">
                <div>
                  <div className="text-sm text-gray-600">AR Posting</div>
                  <div className="font-semibold text-purple-600">{viewing.gl_reference_number}</div>
                </div>
                <span className={`px-3 py-1 rounded-full text-xs font-medium ${getARStatusColor(viewing.ar_posting_status)}`}>
                  {viewing.ar_posting_status}
                </span>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Table View */}
      <div className="bg-white rounded-lg shadow-md border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-gray-100 border-b border-gray-200">
              <tr>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Invoice #</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Customer</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Amount</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Paid</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Due Date</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Status</th>
                <th className="px-4 py-3 text-center font-semibold text-gray-900">Actions</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={7} className="px-4 py-8 text-center text-gray-500">
                    Loading...
                  </td>
                </tr>
              ) : filteredInvoices.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-4 py-8 text-center text-gray-500">
                    No invoices found
                  </td>
                </tr>
              ) : (
                filteredInvoices.map(invoice => (
                  <tr key={invoice.id} className="border-b border-gray-200 hover:bg-gray-50">
                    <td className="px-4 py-3 font-semibold text-gray-900">{invoice.invoice_code}</td>
                    <td className="px-4 py-3 text-gray-700">{invoice.customer_name}</td>
                    <td className="px-4 py-3">
                      <span className="font-semibold text-blue-600">₹{invoice.total_amount.toFixed(2)}</span>
                    </td>
                    <td className="px-4 py-3">
                      <span className="font-semibold text-green-600">₹{invoice.paid_amount.toFixed(2)}</span>
                    </td>
                    <td className="px-4 py-3 text-gray-700">{formatDateToDDMMMYYYY(invoice.due_date)}</td>
                    <td className="px-4 py-3">
                      <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(invoice.payment_status)}`}>
                        {invoice.payment_status}
                      </span>
                    </td>
                    <td className="px-4 py-3 flex items-center justify-center gap-2">
                      <button
                        onClick={() => setViewingId(invoice.id)}
                        title="View Details"
                        className="p-2 text-blue-600 border border-blue-300 rounded hover:bg-blue-50"
                      >
                        <Eye size={16} />
                      </button>
                      <button
                        onClick={() => handleDelete(invoice.id)}
                        title="Delete"
                        className="p-2 text-red-600 border border-red-300 rounded hover:bg-red-50"
                      >
                        <Trash2 size={16} />
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
          <div className="text-2xl font-bold text-blue-600">{invoices.length}</div>
          <div className="text-sm text-gray-600">Total Invoices</div>
        </div>
        <div className="bg-yellow-50 p-4 rounded-lg border border-yellow-200">
          <div className="text-2xl font-bold text-yellow-600">{invoices.filter(i => i.payment_status === 'partially_paid').length}</div>
          <div className="text-sm text-gray-600">Partially Paid</div>
        </div>
        <div className="bg-red-50 p-4 rounded-lg border border-red-200">
          <div className="text-2xl font-bold text-red-600">
            ₹{invoices
              .filter(i => i.payment_status === 'unpaid')
              .reduce((sum, i) => sum + (i.total_amount - i.paid_amount), 0)
              .toLocaleString('en-IN', { maximumFractionDigits: 0 })}
          </div>
          <div className="text-sm text-gray-600">Unpaid Amount</div>
        </div>
        <div className="bg-green-50 p-4 rounded-lg border border-green-200">
          <div className="text-2xl font-bold text-green-600">
            ₹{invoices
              .filter(i => i.payment_status === 'paid')
              .reduce((sum, i) => sum + i.total_amount, 0)
              .toLocaleString('en-IN', { maximumFractionDigits: 0 })}
          </div>
          <div className="text-sm text-gray-600">Paid Total</div>
        </div>
      </div>
    </div>
  )
}
