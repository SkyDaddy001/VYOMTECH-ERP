'use client'

import { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, Search, ArrowRight, DollarSign, FileText, Calendar } from 'lucide-react'
import { formatDateToDDMMMYYYY, formatDateToInput } from '@/lib/dateFormat'

interface QuotationItem {
  id: string
  item_name: string
  description: string
  quantity: number
  unit_price: number
  tax_rate: number
  line_total: number
}

interface Quotation {
  id: string
  quotation_code: string
  customer_id: string
  customer_name: string
  quotation_date: string
  validity_date: string
  status: string
  sub_total: number
  tax_amount: number
  total_amount: number
  discount_percent: number
  discount_amount: number
  items: QuotationItem[]
  notes: string
  created_at: string
}

interface Customer {
  id: string
  customer_name: string
  business_name: string
}

export function QuotationManagement() {
  const [quotations, setQuotations] = useState<Quotation[]>([])
  const [filteredQuotations, setFilteredQuotations] = useState<Quotation[]>([])
  const [customers, setCustomers] = useState<Customer[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState('all')
  const [showForm, setShowForm] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)

  const [formData, setFormData] = useState({
    customer_id: '',
    quotation_date: new Date().toISOString().split('T')[0],
    validity_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    discount_percent: 0,
    notes: '',
    items: [{ item_name: '', description: '', quantity: 1, unit_price: 0, tax_rate: 18 }],
  })

  useEffect(() => {
    fetchQuotations()
    fetchCustomers()
  }, [])

  useEffect(() => {
    filterQuotations()
  }, [quotations, searchTerm, filterStatus])

  const fetchQuotations = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/v1/sales/quotations', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      })
      if (response.ok) {
        const data = await response.json()
        setQuotations(data.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch quotations:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchCustomers = async () => {
    try {
      const response = await fetch('/api/v1/sales/customers', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      })
      if (response.ok) {
        const data = await response.json()
        setCustomers(data.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch customers:', error)
    }
  }

  const filterQuotations = () => {
    let filtered = quotations

    if (searchTerm) {
      filtered = filtered.filter(
        q =>
          q.quotation_code.toLowerCase().includes(searchTerm.toLowerCase()) ||
          q.customer_name.toLowerCase().includes(searchTerm.toLowerCase())
      )
    }

    if (filterStatus !== 'all') {
      filtered = filtered.filter(q => q.status === filterStatus)
    }

    setFilteredQuotations(filtered)
  }

  const calculateTotals = () => {
    const subtotal = formData.items.reduce((sum, item) => sum + item.quantity * item.unit_price, 0)
    const tax = formData.items.reduce((sum, item) => sum + item.quantity * item.unit_price * (item.tax_rate / 100), 0)
    const discountAmount = (subtotal * formData.discount_percent) / 100
    const total = subtotal + tax - discountAmount

    return { subtotal, tax, discountAmount, total }
  }

  const handleAddItem = () => {
    setFormData({
      ...formData,
      items: [...formData.items, { item_name: '', description: '', quantity: 1, unit_price: 0, tax_rate: 18 }],
    })
  }

  const handleRemoveItem = (index: number) => {
    setFormData({
      ...formData,
      items: formData.items.filter((_, i) => i !== index),
    })
  }

  const handleItemChange = (index: number, field: string, value: any) => {
    const updatedItems = [...formData.items]
    updatedItems[index] = { ...updatedItems[index], [field]: value }
    setFormData({ ...formData, items: updatedItems })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      const method = editingId ? 'PUT' : 'POST'
      const url = editingId ? `/api/v1/sales/quotations/${editingId}` : '/api/v1/sales/quotations'

      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
          'X-User-ID': localStorage.getItem('userId') || '',
        },
        body: JSON.stringify(formData),
      })

      if (response.ok) {
        setFormData({
          customer_id: '',
          quotation_date: new Date().toISOString().split('T')[0],
          validity_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
          discount_percent: 0,
          notes: '',
          items: [{ item_name: '', description: '', quantity: 1, unit_price: 0, tax_rate: 18 }],
        })
        setEditingId(null)
        setShowForm(false)
        fetchQuotations()
      }
    } catch (error) {
      console.error('Failed to submit form:', error)
    }
  }

  const handleConvertToOrder = async (quotationId: string) => {
    if (!confirm('Convert this quotation to sales order?')) return

    try {
      const response = await fetch(`/api/v1/sales/quotations/${quotationId}/convert-to-order`, {
        method: 'POST',
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
          'X-User-ID': localStorage.getItem('userId') || '',
        },
      })

      if (response.ok) {
        fetchQuotations()
      }
    } catch (error) {
      console.error('Failed to convert to order:', error)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this quotation?')) return

    try {
      const response = await fetch(`/api/v1/sales/quotations/${id}`, {
        method: 'DELETE',
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      })

      if (response.ok) {
        fetchQuotations()
      }
    } catch (error) {
      console.error('Failed to delete quotation:', error)
    }
  }

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      draft: 'bg-gray-100 text-gray-800',
      sent: 'bg-blue-100 text-blue-800',
      accepted: 'bg-green-100 text-green-800',
      rejected: 'bg-red-100 text-red-800',
      converted: 'bg-purple-100 text-purple-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const totals = calculateTotals()

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-900">Quotations</h2>
        <button
          onClick={() => {
            setShowForm(!showForm)
            setEditingId(null)
            setFormData({
              customer_id: '',
              quotation_date: new Date().toISOString().split('T')[0],
              validity_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
              discount_percent: 0,
              notes: '',
              items: [{ item_name: '', description: '', quantity: 1, unit_price: 0, tax_rate: 18 }],
            })
          }}
          className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition"
        >
          <Plus size={18} />
          New Quotation
        </button>
      </div>

      {/* Form */}
      {showForm && (
        <div className="bg-white p-6 rounded-lg shadow-md border border-gray-200">
          <h3 className="text-lg font-semibold mb-4 text-gray-900">Create New Quotation</h3>
          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Basic Info */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Customer</label>
                <select
                  value={formData.customer_id}
                  onChange={e => setFormData({ ...formData, customer_id: e.target.value })}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Select Customer</option>
                  {customers.map(c => (
                    <option key={c.id} value={c.id}>
                      {c.customer_name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Quotation Date</label>
                <input
                  type="date"
                  value={formData.quotation_date}
                  onChange={e => setFormData({ ...formData, quotation_date: e.target.value })}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Validity Date</label>
                <input
                  type="date"
                  value={formData.validity_date}
                  onChange={e => setFormData({ ...formData, validity_date: e.target.value })}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Items */}
            <div className="border-t pt-4">
              <div className="flex items-center justify-between mb-3">
                <h4 className="font-semibold text-gray-900">Line Items</h4>
                <button type="button" onClick={handleAddItem} className="text-blue-600 text-sm hover:underline">
                  + Add Item
                </button>
              </div>

              <div className="space-y-3 max-h-96 overflow-y-auto">
                {formData.items.map((item, index) => (
                  <div key={index} className="grid grid-cols-1 md:grid-cols-6 gap-2 p-3 bg-gray-50 rounded">
                    <input
                      type="text"
                      placeholder="Item Name"
                      value={item.item_name}
                      onChange={e => handleItemChange(index, 'item_name', e.target.value)}
                      required
                      className="px-2 py-1 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <input
                      type="number"
                      placeholder="Qty"
                      value={item.quantity}
                      onChange={e => handleItemChange(index, 'quantity', parseFloat(e.target.value))}
                      required
                      className="px-2 py-1 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <input
                      type="number"
                      placeholder="Unit Price"
                      value={item.unit_price}
                      onChange={e => handleItemChange(index, 'unit_price', parseFloat(e.target.value))}
                      required
                      className="px-2 py-1 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <select
                      value={item.tax_rate}
                      onChange={e => handleItemChange(index, 'tax_rate', parseFloat(e.target.value))}
                      className="px-2 py-1 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      <option value="0">0%</option>
                      <option value="5">5%</option>
                      <option value="12">12%</option>
                      <option value="18">18%</option>
                      <option value="28">28%</option>
                    </select>
                    <div className="px-2 py-1 bg-white rounded text-sm border border-gray-300">
                      ₹{(item.quantity * item.unit_price).toFixed(2)}
                    </div>
                    <button
                      type="button"
                      onClick={() => handleRemoveItem(index)}
                      className="px-2 py-1 text-red-600 border border-red-300 rounded text-sm hover:bg-red-50"
                    >
                      Remove
                    </button>
                  </div>
                ))}
              </div>
            </div>

            {/* Totals */}
            <div className="border-t pt-4 grid grid-cols-2 md:grid-cols-4 gap-3">
              <div className="bg-blue-50 p-3 rounded">
                <div className="text-xs text-gray-600">Subtotal</div>
                <div className="text-lg font-semibold text-blue-600">₹{totals.subtotal.toFixed(2)}</div>
              </div>
              <div className="bg-green-50 p-3 rounded">
                <div className="text-xs text-gray-600">Tax</div>
                <div className="text-lg font-semibold text-green-600">₹{totals.tax.toFixed(2)}</div>
              </div>
              <div className="bg-yellow-50 p-3 rounded">
                <div className="text-xs text-gray-600">Discount %</div>
                <input
                  type="number"
                  value={formData.discount_percent}
                  onChange={e => setFormData({ ...formData, discount_percent: parseFloat(e.target.value) })}
                  min="0"
                  max="100"
                  className="w-full px-2 py-1 border border-gray-300 rounded text-sm"
                />
              </div>
              <div className="bg-red-50 p-3 rounded">
                <div className="text-xs text-gray-600">Total</div>
                <div className="text-lg font-semibold text-red-600">₹{totals.total.toFixed(2)}</div>
              </div>
            </div>

            {/* Notes */}
            <textarea
              placeholder="Additional notes..."
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
                Create Quotation
              </button>
              <button
                type="button"
                onClick={() => setShowForm(false)}
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
            placeholder="Search quotations..."
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
          <option value="draft">Draft</option>
          <option value="sent">Sent</option>
          <option value="accepted">Accepted</option>
          <option value="rejected">Rejected</option>
          <option value="converted">Converted</option>
        </select>
      </div>

      {/* Table View */}
      <div className="bg-white rounded-lg shadow-md border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-gray-100 border-b border-gray-200">
              <tr>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Quote #</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Customer</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Amount</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Date</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Status</th>
                <th className="px-4 py-3 text-center font-semibold text-gray-900">Actions</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={6} className="px-4 py-8 text-center text-gray-500">
                    Loading...
                  </td>
                </tr>
              ) : filteredQuotations.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-4 py-8 text-center text-gray-500">
                    No quotations found
                  </td>
                </tr>
              ) : (
                filteredQuotations.map(quotation => (
                  <tr key={quotation.id} className="border-b border-gray-200 hover:bg-gray-50">
                    <td className="px-4 py-3 font-semibold text-gray-900">{quotation.quotation_code}</td>
                    <td className="px-4 py-3 text-gray-700">{quotation.customer_name}</td>
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-1">
                        <DollarSign size={14} className="text-blue-600" />
                        <span className="font-semibold text-blue-600">₹{quotation.total_amount.toFixed(2)}</span>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-gray-700">{formatDateToDDMMMYYYY(quotation.quotation_date)}</td>
                    <td className="px-4 py-3">
                      <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(quotation.status)}`}>
                        {quotation.status}
                      </span>
                    </td>
                    <td className="px-4 py-3 flex items-center justify-center gap-2">
                      {quotation.status !== 'converted' && quotation.status === 'accepted' && (
                        <button
                          onClick={() => handleConvertToOrder(quotation.id)}
                          title="Convert to Order"
                          className="p-2 text-green-600 border border-green-300 rounded hover:bg-green-50"
                        >
                          <ArrowRight size={16} />
                        </button>
                      )}
                      <button
                        onClick={() => handleDelete(quotation.id)}
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
          <div className="text-2xl font-bold text-blue-600">{quotations.length}</div>
          <div className="text-sm text-gray-600">Total Quotations</div>
        </div>
        <div className="bg-green-50 p-4 rounded-lg border border-green-200">
          <div className="text-2xl font-bold text-green-600">{quotations.filter(q => q.status === 'accepted').length}</div>
          <div className="text-sm text-gray-600">Accepted</div>
        </div>
        <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
          <div className="text-2xl font-bold text-purple-600">
            ₹{quotations.reduce((sum, q) => sum + q.total_amount, 0).toLocaleString('en-IN', { maximumFractionDigits: 0 })}
          </div>
          <div className="text-sm text-gray-600">Total Value</div>
        </div>
        <div className="bg-orange-50 p-4 rounded-lg border border-orange-200">
          <div className="text-2xl font-bold text-orange-600">{quotations.filter(q => q.status === 'converted').length}</div>
          <div className="text-sm text-gray-600">Converted</div>
        </div>
      </div>
    </div>
  )
}
