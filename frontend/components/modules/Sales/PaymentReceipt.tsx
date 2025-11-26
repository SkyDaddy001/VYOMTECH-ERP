'use client'

import { useState, useEffect } from 'react'
import { Plus, Check, Trash2, Search, DollarSign, Calendar, Printer } from 'lucide-react'
import { formatDateToDDMMMYYYY, formatDateToInput } from '@/lib/dateFormat'

interface Payment {
  id: string
  invoice_id: string
  invoice_code: string
  customer_id: string
  customer_name: string
  payment_date: string
  payment_method: string
  amount_paid: number
  reference_number: string
  remarks: string
  receipt_generated: boolean
  created_at: string
}

interface Invoice {
  id: string
  invoice_code: string
  customer_name: string
  total_amount: number
  paid_amount: number
  payment_status: string
}

export function PaymentReceipt() {
  const [payments, setPayments] = useState<Payment[]>([])
  const [filteredPayments, setFilteredPayments] = useState<Payment[]>([])
  const [invoices, setInvoices] = useState<Invoice[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [filterMethod, setFilterMethod] = useState('all')
  const [showForm, setShowForm] = useState(false)
  const [viewingId, setViewingId] = useState<string | null>(null)

  const [formData, setFormData] = useState({
    invoice_id: '',
    payment_date: new Date().toISOString().split('T')[0],
    payment_method: 'bank_transfer',
    amount_paid: 0,
    reference_number: '',
    remarks: '',
  })

  useEffect(() => {
    fetchPayments()
    fetchInvoices()
  }, [])

  useEffect(() => {
    filterPayments()
  }, [payments, searchTerm, filterMethod])

  const fetchPayments = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/v1/sales/payments', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      })
      if (response.ok) {
        const data = await response.json()
        setPayments(data.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch payments:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchInvoices = async () => {
    try {
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
    }
  }

  const filterPayments = () => {
    let filtered = payments

    if (searchTerm) {
      filtered = filtered.filter(
        p =>
          p.invoice_code.toLowerCase().includes(searchTerm.toLowerCase()) ||
          p.customer_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          p.reference_number.toLowerCase().includes(searchTerm.toLowerCase())
      )
    }

    if (filterMethod !== 'all') {
      filtered = filtered.filter(p => p.payment_method === filterMethod)
    }

    setFilteredPayments(filtered)
  }

  const getPaymentMethodLabel = (method: string) => {
    const labels: Record<string, string> = {
      bank_transfer: '🏦 Bank Transfer',
      cheque: '📄 Cheque',
      cash: '💵 Cash',
      credit_card: '💳 Credit Card',
      digital_payment: '📱 Digital Payment',
    }
    return labels[method] || method
  }

  const getPaymentMethodColor = (method: string) => {
    const colors: Record<string, string> = {
      bank_transfer: 'bg-blue-100 text-blue-800',
      cheque: 'bg-yellow-100 text-yellow-800',
      cash: 'bg-green-100 text-green-800',
      credit_card: 'bg-purple-100 text-purple-800',
      digital_payment: 'bg-pink-100 text-pink-800',
    }
    return colors[method] || 'bg-gray-100 text-gray-800'
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      const response = await fetch('/api/v1/sales/payments', {
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
          invoice_id: '',
          payment_date: new Date().toISOString().split('T')[0],
          payment_method: 'bank_transfer',
          amount_paid: 0,
          reference_number: '',
          remarks: '',
        })
        setShowForm(false)
        fetchPayments()
        fetchInvoices()
      }
    } catch (error) {
      console.error('Failed to record payment:', error)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this payment record?')) return

    try {
      const response = await fetch(`/api/v1/sales/payments/${id}`, {
        method: 'DELETE',
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenantId') || '',
        },
      })

      if (response.ok) {
        fetchPayments()
      }
    } catch (error) {
      console.error('Failed to delete payment:', error)
    }
  }

  const handleGenerateReceipt = async (paymentId: string) => {
    const payment = payments.find(p => p.id === paymentId)
    if (!payment) return

    // Generate receipt content
    const receipt = `
PAYMENT RECEIPT
===============================
Receipt #: ${payment.id.substring(0, 8).toUpperCase()}
Date: ${formatDateToDDMMMYYYY(payment.payment_date)}

Customer: ${payment.customer_name}
Invoice: ${payment.invoice_code}

Payment Details:
Amount Paid: ₹${payment.amount_paid.toFixed(2)}
Method: ${getPaymentMethodLabel(payment.payment_method)}
Reference: ${payment.reference_number}
Remarks: ${payment.remarks}

This is a computer-generated receipt.
===============================
    `

    // Create downloadable file
    const element = document.createElement('a')
    element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(receipt))
    element.setAttribute('download', `receipt_${payment.invoice_code}_${payment.id.substring(0, 8)}.txt`)
    element.style.display = 'none'
    document.body.appendChild(element)
    element.click()
    document.body.removeChild(element)
  }

  const viewing = payments.find(p => p.id === viewingId)
  const selectedInvoice = formData.invoice_id ? invoices.find(i => i.id === formData.invoice_id) : null

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-900">Payment Receipts</h2>
        <button
          onClick={() => {
            setShowForm(!showForm)
            setFormData({
              invoice_id: '',
              payment_date: new Date().toISOString().split('T')[0],
              payment_method: 'bank_transfer',
              amount_paid: 0,
              reference_number: '',
              remarks: '',
            })
          }}
          className="flex items-center gap-2 bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition"
        >
          <Plus size={18} />
          Record Payment
        </button>
      </div>

      {/* Form */}
      {showForm && (
        <div className="bg-white p-6 rounded-lg shadow-md border border-gray-200">
          <h3 className="text-lg font-semibold mb-4 text-gray-900">Record Payment</h3>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Select Invoice</label>
                <select
                  value={formData.invoice_id}
                  onChange={e => {
                    setFormData({ ...formData, invoice_id: e.target.value })
                  }}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
                >
                  <option value="">Select Invoice</option>
                  {invoices
                    .filter(inv => inv.payment_status !== 'paid')
                    .map(inv => (
                      <option key={inv.id} value={inv.id}>
                        {inv.invoice_code} - {inv.customer_name} (Due: ₹{(inv.total_amount - inv.paid_amount).toFixed(2)})
                      </option>
                    ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Payment Date</label>
                <input
                  type="date"
                  value={formData.payment_date}
                  onChange={e => setFormData({ ...formData, payment_date: e.target.value })}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Payment Method</label>
                <select
                  value={formData.payment_method}
                  onChange={e => setFormData({ ...formData, payment_method: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
                >
                  <option value="bank_transfer">Bank Transfer</option>
                  <option value="cheque">Cheque</option>
                  <option value="cash">Cash</option>
                  <option value="credit_card">Credit Card</option>
                  <option value="digital_payment">Digital Payment</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Amount Paid</label>
                <input
                  type="number"
                  step="0.01"
                  placeholder="0.00"
                  value={formData.amount_paid}
                  onChange={e => setFormData({ ...formData, amount_paid: parseFloat(e.target.value) })}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Reference Number</label>
                <input
                  type="text"
                  placeholder="Cheque #, Transaction ID, etc."
                  value={formData.reference_number}
                  onChange={e => setFormData({ ...formData, reference_number: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Remarks</label>
                <input
                  type="text"
                  placeholder="Additional notes..."
                  value={formData.remarks}
                  onChange={e => setFormData({ ...formData, remarks: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
                />
              </div>
            </div>

            {selectedInvoice && (
              <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                  <div>
                    <div className="text-xs text-gray-600">Invoice Amount</div>
                    <div className="font-semibold text-green-600">₹{selectedInvoice.total_amount.toFixed(2)}</div>
                  </div>
                  <div>
                    <div className="text-xs text-gray-600">Already Paid</div>
                    <div className="font-semibold text-blue-600">₹{selectedInvoice.paid_amount.toFixed(2)}</div>
                  </div>
                  <div>
                    <div className="text-xs text-gray-600">Balance Due</div>
                    <div className="font-semibold text-red-600">₹{(selectedInvoice.total_amount - selectedInvoice.paid_amount).toFixed(2)}</div>
                  </div>
                </div>
              </div>
            )}

            <div className="flex gap-2">
              <button
                type="submit"
                className="flex-1 bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition font-medium"
              >
                Record Payment
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
            placeholder="Search by invoice or customer..."
            value={searchTerm}
            onChange={e => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>
        <select
          value={filterMethod}
          onChange={e => setFilterMethod(e.target.value)}
          className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
        >
          <option value="all">All Methods</option>
          <option value="bank_transfer">Bank Transfer</option>
          <option value="cheque">Cheque</option>
          <option value="cash">Cash</option>
          <option value="credit_card">Credit Card</option>
          <option value="digital_payment">Digital Payment</option>
        </select>
      </div>

      {/* View Receipt Detail */}
      {viewing && (
        <div className="bg-white p-6 rounded-lg shadow-md border border-gray-200">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-gray-900">Payment Receipt</h3>
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
              <div className="text-sm text-gray-600">Invoice</div>
              <div className="text-lg font-semibold text-gray-900">{viewing.invoice_code}</div>
            </div>
            <div>
              <div className="text-sm text-gray-600">Payment Date</div>
              <div className="text-lg font-semibold text-gray-900">{formatDateToDDMMMYYYY(viewing.payment_date)}</div>
            </div>
            <div>
              <div className="text-sm text-gray-600">Payment Method</div>
              <span className={`inline-block px-3 py-1 rounded-full text-xs font-medium ${getPaymentMethodColor(viewing.payment_method)}`}>
                {getPaymentMethodLabel(viewing.payment_method)}
              </span>
            </div>
          </div>

          <div className="bg-blue-50 p-4 rounded-lg border border-blue-200 mb-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
              <div className="text-center">
                <div className="text-sm text-gray-600">Amount Paid</div>
                <div className="text-2xl font-bold text-blue-600">₹{viewing.amount_paid.toFixed(2)}</div>
              </div>
              {viewing.reference_number && (
                <div className="text-center">
                  <div className="text-sm text-gray-600">Reference</div>
                  <div className="text-lg font-semibold text-gray-900">{viewing.reference_number}</div>
                </div>
              )}
              {viewing.remarks && (
                <div className="text-center">
                  <div className="text-sm text-gray-600">Remarks</div>
                  <div className="text-lg font-semibold text-gray-900">{viewing.remarks}</div>
                </div>
              )}
            </div>
          </div>

          <button
            onClick={() => handleGenerateReceipt(viewing.id)}
            className="w-full flex items-center justify-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition"
          >
            <Printer size={18} />
            Download Receipt
          </button>
        </div>
      )}

      {/* Table View */}
      <div className="bg-white rounded-lg shadow-md border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-gray-100 border-b border-gray-200">
              <tr>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Invoice</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Customer</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Amount</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Method</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Date</th>
                <th className="px-4 py-3 text-left font-semibold text-gray-900">Reference</th>
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
              ) : filteredPayments.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-4 py-8 text-center text-gray-500">
                    No payments found
                  </td>
                </tr>
              ) : (
                filteredPayments.map(payment => (
                  <tr key={payment.id} className="border-b border-gray-200 hover:bg-gray-50">
                    <td className="px-4 py-3 font-semibold text-gray-900">{payment.invoice_code}</td>
                    <td className="px-4 py-3 text-gray-700">{payment.customer_name}</td>
                    <td className="px-4 py-3">
                      <span className="font-semibold text-green-600">₹{payment.amount_paid.toFixed(2)}</span>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`inline-block px-3 py-1 rounded-full text-xs font-medium ${getPaymentMethodColor(payment.payment_method)}`}>
                        {getPaymentMethodLabel(payment.payment_method)}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-gray-700">{formatDateToDDMMMYYYY(payment.payment_date)}</td>
                    <td className="px-4 py-3 text-gray-700">{payment.reference_number || '-'}</td>
                    <td className="px-4 py-3 flex items-center justify-center gap-2">
                      <button
                        onClick={() => {
                          setViewingId(payment.id)
                        }}
                        title="View Receipt"
                        className="p-2 text-blue-600 border border-blue-300 rounded hover:bg-blue-50"
                      >
                        <Check size={16} />
                      </button>
                      <button
                        onClick={() => handleGenerateReceipt(payment.id)}
                        title="Download Receipt"
                        className="p-2 text-green-600 border border-green-300 rounded hover:bg-green-50"
                      >
                        <Printer size={16} />
                      </button>
                      <button
                        onClick={() => handleDelete(payment.id)}
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
          <div className="text-2xl font-bold text-blue-600">{payments.length}</div>
          <div className="text-sm text-gray-600">Total Payments</div>
        </div>
        <div className="bg-green-50 p-4 rounded-lg border border-green-200">
          <div className="text-2xl font-bold text-green-600">
            ₹{payments.reduce((sum, p) => sum + p.amount_paid, 0).toLocaleString('en-IN', { maximumFractionDigits: 0 })}
          </div>
          <div className="text-sm text-gray-600">Total Collected</div>
        </div>
        <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
          <div className="text-2xl font-bold text-purple-600">
            {payments.filter(p => p.payment_method === 'bank_transfer').length}
          </div>
          <div className="text-sm text-gray-600">Bank Transfers</div>
        </div>
        <div className="bg-yellow-50 p-4 rounded-lg border border-yellow-200">
          <div className="text-2xl font-bold text-yellow-600">{payments.filter(p => p.payment_method === 'cheque').length}</div>
          <div className="text-sm text-gray-600">Cheques Received</div>
        </div>
      </div>
    </div>
  )
}
