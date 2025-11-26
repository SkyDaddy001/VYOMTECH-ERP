'use client'

import { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, Search, MapPin, Phone, Mail } from 'lucide-react'

interface Customer {
  id: string
  customer_code: string
  customer_name: string
  business_name: string
  business_type: string
  industry: string
  primary_contact_name: string
  primary_email: string
  primary_phone: string
  billing_city: string
  gst_number: string
  credit_limit: number
  status: string
  current_balance: number
  created_at: string
}

export function CustomerManagement() {
  const [customers, setCustomers] = useState<Customer[]>([])
  const [filteredCustomers, setFilteredCustomers] = useState<Customer[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState('all')
  const [showForm, setShowForm] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)

  const [formData, setFormData] = useState({
    customer_name: '',
    business_name: '',
    business_type: 'proprietorship',
    industry: '',
    primary_contact_name: '',
    primary_email: '',
    primary_phone: '',
    billing_address: '',
    billing_city: '',
    billing_state: '',
    billing_country: 'India',
    billing_zip: '',
    gst_number: '',
    credit_limit: 0,
    credit_days: 30,
    payment_terms: 'Net 30',
    customer_category: 'regular',
  })

  useEffect(() => {
    fetchCustomers()
  }, [])

  useEffect(() => {
    filterCustomers()
  }, [customers, searchTerm, filterStatus])

  const fetchCustomers = async () => {
    try {
      setLoading(true)
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
    } finally {
      setLoading(false)
    }
  }

  const filterCustomers = () => {
    let filtered = customers

    if (searchTerm) {
      filtered = filtered.filter(
        customer =>
          customer.customer_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          customer.business_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          customer.primary_email.toLowerCase().includes(searchTerm.toLowerCase()) ||
          customer.gst_number.toLowerCase().includes(searchTerm.toLowerCase())
      )
    }

    if (filterStatus !== 'all') {
      filtered = filtered.filter(customer => customer.status === filterStatus)
    }

    setFilteredCustomers(filtered)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      const method = editingId ? 'PUT' : 'POST'
      const url = editingId ? `/api/v1/sales/customers/${editingId}` : '/api/v1/sales/customers'

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
          customer_name: '',
          business_name: '',
          business_type: 'proprietorship',
          industry: '',
          primary_contact_name: '',
          primary_email: '',
          primary_phone: '',
          billing_address: '',
          billing_city: '',
          billing_state: '',
          billing_country: 'India',
          billing_zip: '',
          gst_number: '',
          credit_limit: 0,
          credit_days: 30,
          payment_terms: 'Net 30',
          customer_category: 'regular',
        })
        setEditingId(null)
        setShowForm(false)
        fetchCustomers()
      }
    } catch (error) {
      console.error('Failed to submit form:', error)
    }
  }

  const handleEdit = (customer: Customer) => {
    setFormData({
      customer_name: customer.customer_name,
      business_name: customer.business_name,
      business_type: customer.business_type,
      industry: customer.industry,
      primary_contact_name: customer.primary_contact_name,
      primary_email: customer.primary_email,
      primary_phone: customer.primary_phone,
      billing_address: '',
      billing_city: customer.billing_city,
      billing_state: '',
      billing_country: 'India',
      billing_zip: '',
      gst_number: customer.gst_number,
      credit_limit: customer.credit_limit,
      credit_days: 30,
      payment_terms: 'Net 30',
      customer_category: 'regular',
    })
    setEditingId(customer.id)
    setShowForm(true)
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this customer?')) return

    try {
      const response = await fetch(`/api/v1/sales/customers/${id}`, {
        method: 'DELETE',
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      })

      if (response.ok) {
        fetchCustomers()
      }
    } catch (error) {
      console.error('Failed to delete customer:', error)
    }
  }

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      active: 'bg-green-100 text-green-800',
      inactive: 'bg-gray-100 text-gray-800',
      blocked: 'bg-red-100 text-red-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const getCategoryColor = (category: string) => {
    const colors: Record<string, string> = {
      gold: 'text-yellow-600',
      silver: 'text-gray-400',
      bronze: 'text-orange-600',
      regular: 'text-blue-600',
    }
    return colors[category] || 'text-gray-600'
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-900">Customers</h2>
        <button
          onClick={() => {
            setShowForm(!showForm)
            setEditingId(null)
            setFormData({
              customer_name: '',
              business_name: '',
              business_type: 'proprietorship',
              industry: '',
              primary_contact_name: '',
              primary_email: '',
              primary_phone: '',
              billing_address: '',
              billing_city: '',
              billing_state: '',
              billing_country: 'India',
              billing_zip: '',
              gst_number: '',
              credit_limit: 0,
              credit_days: 30,
              payment_terms: 'Net 30',
              customer_category: 'regular',
            })
          }}
          className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition"
        >
          <Plus size={18} />
          New Customer
        </button>
      </div>

      {/* Form */}
      {showForm && (
        <div className="bg-white p-6 rounded-lg shadow-md border border-gray-200">
          <h3 className="text-lg font-semibold mb-4 text-gray-900">
            {editingId ? 'Edit Customer' : 'Create New Customer'}
          </h3>
          <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <input
              type="text"
              placeholder="Customer Name"
              value={formData.customer_name}
              onChange={e => setFormData({ ...formData, customer_name: e.target.value })}
              required
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="text"
              placeholder="Business Name"
              value={formData.business_name}
              onChange={e => setFormData({ ...formData, business_name: e.target.value })}
              required
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <select
              value={formData.business_type}
              onChange={e => setFormData({ ...formData, business_type: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="individual">Individual</option>
              <option value="proprietorship">Proprietorship</option>
              <option value="partnership">Partnership</option>
              <option value="pvt_ltd">Private Limited</option>
              <option value="public_ltd">Public Limited</option>
            </select>
            <input
              type="text"
              placeholder="Industry"
              value={formData.industry}
              onChange={e => setFormData({ ...formData, industry: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="text"
              placeholder="Contact Name"
              value={formData.primary_contact_name}
              onChange={e => setFormData({ ...formData, primary_contact_name: e.target.value })}
              required
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="email"
              placeholder="Email"
              value={formData.primary_email}
              onChange={e => setFormData({ ...formData, primary_email: e.target.value })}
              required
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="tel"
              placeholder="Phone"
              value={formData.primary_phone}
              onChange={e => setFormData({ ...formData, primary_phone: e.target.value })}
              required
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="text"
              placeholder="GST Number"
              value={formData.gst_number}
              onChange={e => setFormData({ ...formData, gst_number: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="text"
              placeholder="Billing City"
              value={formData.billing_city}
              onChange={e => setFormData({ ...formData, billing_city: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="number"
              placeholder="Credit Limit"
              value={formData.credit_limit}
              onChange={e => setFormData({ ...formData, credit_limit: parseFloat(e.target.value) })}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <div className="flex gap-2 md:col-span-2">
              <button
                type="submit"
                className="flex-1 bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition"
              >
                {editingId ? 'Update Customer' : 'Create Customer'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setShowForm(false)
                  setEditingId(null)
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
            placeholder="Search customers..."
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
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
          <option value="blocked">Blocked</option>
        </select>
      </div>

      {/* Grid View */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {loading ? (
          <div className="col-span-full p-8 text-center text-gray-500">Loading customers...</div>
        ) : filteredCustomers.length === 0 ? (
          <div className="col-span-full p-8 text-center text-gray-500">No customers found</div>
        ) : (
          filteredCustomers.map(customer => (
            <div key={customer.id} className="bg-white p-5 rounded-lg shadow-md border border-gray-200 hover:shadow-lg transition">
              <div className="flex items-start justify-between mb-3">
                <div>
                  <h3 className="font-semibold text-gray-900">{customer.customer_name}</h3>
                  <p className="text-sm text-gray-600">{customer.business_name}</p>
                </div>
                <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(customer.status)}`}>
                  {customer.status}
                </span>
              </div>

              <div className="space-y-2 mb-4 text-sm text-gray-600">
                <div className="flex items-center gap-2">
                  <Mail size={16} />
                  <span>{customer.primary_email}</span>
                </div>
                <div className="flex items-center gap-2">
                  <Phone size={16} />
                  <span>{customer.primary_phone}</span>
                </div>
                <div className="flex items-center gap-2">
                  <MapPin size={16} />
                  <span>{customer.billing_city}</span>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3 mb-4 text-sm">
                <div className="bg-blue-50 p-2 rounded">
                  <div className="text-xs text-gray-600">Credit Limit</div>
                  <div className="font-semibold text-blue-600">₹{customer.credit_limit.toLocaleString()}</div>
                </div>
                <div className="bg-red-50 p-2 rounded">
                  <div className="text-xs text-gray-600">Balance</div>
                  <div className="font-semibold text-red-600">₹{customer.current_balance.toLocaleString()}</div>
                </div>
              </div>

              <div className="flex gap-2">
                <button
                  onClick={() => handleEdit(customer)}
                  className="flex-1 p-2 text-blue-600 border border-blue-300 rounded hover:bg-blue-50 transition text-sm font-medium"
                >
                  <Edit2 size={14} className="inline mr-1" />
                  Edit
                </button>
                <button
                  onClick={() => handleDelete(customer.id)}
                  className="flex-1 p-2 text-red-600 border border-red-300 rounded hover:bg-red-50 transition text-sm font-medium"
                >
                  <Trash2 size={14} className="inline mr-1" />
                  Delete
                </button>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-green-50 p-4 rounded-lg border border-green-200">
          <div className="text-2xl font-bold text-green-600">{customers.length}</div>
          <div className="text-sm text-gray-600">Total Customers</div>
        </div>
        <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
          <div className="text-2xl font-bold text-blue-600">{customers.filter(c => c.status === 'active').length}</div>
          <div className="text-sm text-gray-600">Active</div>
        </div>
        <div className="bg-yellow-50 p-4 rounded-lg border border-yellow-200">
          <div className="text-2xl font-bold text-yellow-600">
            ₹{customers.reduce((sum, c) => sum + c.credit_limit, 0).toLocaleString('en-IN', { maximumFractionDigits: 0 })}
          </div>
          <div className="text-sm text-gray-600">Total Credit</div>
        </div>
        <div className="bg-red-50 p-4 rounded-lg border border-red-200">
          <div className="text-2xl font-bold text-red-600">
            ₹{customers.reduce((sum, c) => sum + c.current_balance, 0).toLocaleString('en-IN', { maximumFractionDigits: 0 })}
          </div>
          <div className="text-sm text-gray-600">Total Balance</div>
        </div>
      </div>
    </div>
  )
}
