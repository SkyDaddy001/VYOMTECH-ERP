'use client'

import React, { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'

type PaymentProvider = 'razorpay' | 'billdesk'
type PaymentMethod = 'netbanking' | 'credit_card' | 'debit_card' | 'upi' | 'wallet'

interface PaymentMethod {
  name: string
  code: string
  description: string
  icon: string
}

interface Bank {
  code: string
  name: string
}

export default function BillingPayment() {
  const router = useRouter()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [payments, setPayments] = useState<any[]>([])
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([])
  const [banks, setBanks] = useState<Bank[]>([])
  const [summary, setSummary] = useState<any>(null)

  const [formData, setFormData] = useState({
    amount: '',
    currency: 'INR',
    provider: 'razorpay' as PaymentProvider,
    paymentMethod: 'upi' as PaymentMethod,
    description: 'VYOM ERP Subscription',
    customerName: '',
    customerEmail: '',
    customerPhone: '',
    selectedBank: '',
  })

  useEffect(() => {
    const token = localStorage.getItem('authToken')
    if (!token) {
      router.push('/login')
      return
    }

    fetchPaymentMethods()
    fetchBanks()
    fetchPayments()
    fetchPaymentSummary()
  }, [])

  const fetchPaymentMethods = async () => {
    try {
      const response = await fetch('/api/v1/payments/methods/available', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('authToken')}` },
      })
      if (response.ok) {
        const data = await response.json()
        setPaymentMethods(data.payment_methods || [])
      }
    } catch (err) {
      console.error('Failed to fetch payment methods')
    }
  }

  const fetchBanks = async () => {
    try {
      const response = await fetch('/api/v1/payments/methods/banks', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('authToken')}` },
      })
      if (response.ok) {
        const data = await response.json()
        setBanks(data.banks || [])
      }
    } catch (err) {
      console.error('Failed to fetch banks')
    }
  }

  const fetchPayments = async () => {
    try {
      const response = await fetch('/api/v1/payments/list', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('authToken')}` },
      })
      if (response.ok) {
        const data = await response.json()
        setPayments(data.payments || [])
      }
    } catch (err) {
      console.error('Failed to fetch payments')
    }
  }

  const fetchPaymentSummary = async () => {
    try {
      const response = await fetch('/api/v1/payments/summary', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('authToken')}` },
      })
      if (response.ok) {
        const data = await response.json()
        setSummary(data)
      }
    } catch (err) {
      console.error('Failed to fetch summary')
    }
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }))
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')
    setSuccess('')
    setLoading(true)

    try {
      if (!formData.amount || parseFloat(formData.amount) <= 0) {
        setError('Please enter a valid amount')
        setLoading(false)
        return
      }

      const response = await fetch('/api/v1/payments/initiate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
        },
        body: JSON.stringify({
          amount: parseFloat(formData.amount),
          currency: formData.currency,
          provider: formData.provider,
          payment_method: formData.paymentMethod,
          description: formData.description,
          customer_name: formData.customerName,
          customer_email: formData.customerEmail,
          customer_phone: formData.customerPhone,
        }),
      })

      const data = await response.json()

      if (!response.ok) {
        setError(data.error || 'Failed to initiate payment')
        setLoading(false)
        return
      }

      // Redirect to payment gateway
      if (data.payment_url) {
        window.location.href = data.payment_url
      }

      setSuccess('Payment initiated successfully')
      fetchPayments()
      fetchPaymentSummary()
    } catch (err: any) {
      setError(err.message || 'Failed to process payment')
    } finally {
      setLoading(false)
    }
  }

  const getStatusBadge = (status: string) => {
    const statusStyles: { [key: string]: string } = {
      'successful': 'bg-green-100 text-green-800',
      'pending': 'bg-yellow-100 text-yellow-800',
      'failed': 'bg-red-100 text-red-800',
      'processing': 'bg-blue-100 text-blue-800',
    }
    return statusStyles[status] || 'bg-gray-100 text-gray-800'
  }

  const getMethodIcon = (method: string) => {
    const icons: { [key: string]: string } = {
      'netbanking': '🏦',
      'credit_card': '💳',
      'debit_card': '💳',
      'upi': '📱',
      'wallet': '👛',
    }
    return icons[method] || '💰'
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 py-12 px-4">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-slate-900 mb-2">💳 Payment Gateway</h1>
          <p className="text-slate-600">Process payments securely using Razorpay or Billdesk</p>
        </div>

        {/* Summary Cards */}
        {summary && (
          <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-8">
            <div className="bg-white rounded-lg shadow-sm border border-slate-200 p-6">
              <p className="text-slate-600 text-sm font-medium">Total Payments</p>
              <p className="text-3xl font-bold text-slate-900 mt-2">{summary.total_payments}</p>
            </div>
            <div className="bg-white rounded-lg shadow-sm border border-slate-200 p-6">
              <p className="text-slate-600 text-sm font-medium">Total Successful</p>
              <p className="text-3xl font-bold text-green-600 mt-2">₹{summary.total_successful?.toFixed(2)}</p>
            </div>
            <div className="bg-white rounded-lg shadow-sm border border-slate-200 p-6">
              <p className="text-slate-600 text-sm font-medium">Pending</p>
              <p className="text-3xl font-bold text-yellow-600 mt-2">₹{summary.total_pending?.toFixed(2)}</p>
            </div>
            <div className="bg-white rounded-lg shadow-sm border border-slate-200 p-6">
              <p className="text-slate-600 text-sm font-medium">Failed</p>
              <p className="text-3xl font-bold text-red-600 mt-2">₹{summary.total_failed?.toFixed(2)}</p>
            </div>
            <div className="bg-white rounded-lg shadow-sm border border-slate-200 p-6">
              <p className="text-slate-600 text-sm font-medium">Success Rate</p>
              <p className="text-3xl font-bold text-blue-600 mt-2">{summary.success_rate?.toFixed(1)}%</p>
            </div>
          </div>
        )}

        {/* Main Content */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Payment Form */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow-md border border-slate-200 p-6">
              <h2 className="text-2xl font-bold text-slate-900 mb-6">Make a Payment</h2>

              {error && (
                <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
                  {error}
                </div>
              )}

              {success && (
                <div className="mb-4 p-4 bg-green-50 border border-green-200 rounded-lg text-green-700 text-sm">
                  {success}
                </div>
              )}

              <form onSubmit={handleSubmit} className="space-y-4">
                {/* Amount */}
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-2">Amount</label>
                  <div className="flex">
                    <input
                      type="number"
                      name="amount"
                      step="0.01"
                      min="1"
                      value={formData.amount}
                      onChange={handleInputChange}
                      placeholder="0.00"
                      className="flex-1 px-3 py-2 border border-slate-300 rounded-l-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                    />
                    <select
                      name="currency"
                      value={formData.currency}
                      onChange={handleInputChange}
                      className="px-3 py-2 border border-l-0 border-slate-300 rounded-r-lg bg-slate-50 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                    >
                      <option value="INR">INR</option>
                      <option value="USD">USD</option>
                    </select>
                  </div>
                </div>

                {/* Provider */}
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-2">Payment Provider</label>
                  <select
                    name="provider"
                    value={formData.provider}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                  >
                    <option value="razorpay">Razorpay</option>
                    <option value="billdesk">Billdesk</option>
                  </select>
                </div>

                {/* Payment Method */}
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-2">Payment Method</label>
                  <select
                    name="paymentMethod"
                    value={formData.paymentMethod}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                  >
                    <option value="upi">UPI</option>
                    <option value="netbanking">Netbanking</option>
                    <option value="credit_card">Credit Card</option>
                    <option value="debit_card">Debit Card</option>
                    <option value="wallet">Digital Wallet</option>
                  </select>
                </div>

                {/* Bank Selection (for Netbanking) */}
                {formData.paymentMethod === 'netbanking' && banks.length > 0 && (
                  <div>
                    <label className="block text-sm font-medium text-slate-700 mb-2">Select Bank</label>
                    <select
                      name="selectedBank"
                      value={formData.selectedBank}
                      onChange={handleInputChange}
                      className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                    >
                      <option value="">-- Select Bank --</option>
                      {banks.map(bank => (
                        <option key={bank.code} value={bank.code}>{bank.name}</option>
                      ))}
                    </select>
                  </div>
                )}

                {/* Customer Details */}
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-2">Name</label>
                  <input
                    type="text"
                    name="customerName"
                    value={formData.customerName}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-2">Email</label>
                  <input
                    type="email"
                    name="customerEmail"
                    value={formData.customerEmail}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-2">Phone</label>
                  <input
                    type="tel"
                    name="customerPhone"
                    value={formData.customerPhone}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                  />
                </div>

                {/* Submit Button */}
                <button
                  type="submit"
                  disabled={loading}
                  className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-bold py-2 px-4 rounded-lg transition-colors duration-200"
                >
                  {loading ? 'Processing...' : 'Proceed to Payment'}
                </button>
              </form>

              {/* Payment Methods Info */}
              <div className="mt-8 pt-6 border-t border-slate-200">
                <h3 className="font-semibold text-slate-900 mb-4">Available Payment Methods</h3>
                <div className="space-y-3">
                  {paymentMethods.map(method => (
                    <div key={method.code} className="flex items-start">
                      <span className="text-2xl mr-3">{getMethodIcon(method.code)}</span>
                      <div>
                        <p className="font-medium text-slate-900">{method.name}</p>
                        <p className="text-xs text-slate-500">{method.description}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>

          {/* Payment History */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow-md border border-slate-200 p-6">
              <h2 className="text-2xl font-bold text-slate-900 mb-6">Payment History</h2>

              {payments.length === 0 ? (
                <div className="text-center py-12">
                  <p className="text-slate-500">No payments yet. Make your first payment above.</p>
                </div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-slate-50 border-b border-slate-200">
                      <tr>
                        <th className="px-6 py-3 text-left text-sm font-semibold text-slate-700">Order ID</th>
                        <th className="px-6 py-3 text-left text-sm font-semibold text-slate-700">Amount</th>
                        <th className="px-6 py-3 text-left text-sm font-semibold text-slate-700">Method</th>
                        <th className="px-6 py-3 text-left text-sm font-semibold text-slate-700">Provider</th>
                        <th className="px-6 py-3 text-left text-sm font-semibold text-slate-700">Status</th>
                        <th className="px-6 py-3 text-left text-sm font-semibold text-slate-700">Date</th>
                      </tr>
                    </thead>
                    <tbody>
                      {payments.map(payment => (
                        <tr key={payment.id} className="border-b border-slate-200 hover:bg-slate-50 transition-colors">
                          <td className="px-6 py-3 text-sm font-medium text-slate-900">{payment.order_id}</td>
                          <td className="px-6 py-3 text-sm font-medium text-slate-900">₹{payment.amount?.toFixed(2)}</td>
                          <td className="px-6 py-3 text-sm text-slate-600">{getMethodIcon(payment.payment_method)} {payment.payment_method.replace('_', ' ')}</td>
                          <td className="px-6 py-3 text-sm text-slate-600 capitalize">{payment.provider}</td>
                          <td className="px-6 py-3 text-sm">
                            <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getStatusBadge(payment.status)}`}>
                              {payment.status}
                            </span>
                          </td>
                          <td className="px-6 py-3 text-sm text-slate-600">{new Date(payment.created_at).toLocaleDateString()}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Additional Info Section */}
        <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Razorpay Info */}
          <div className="bg-gradient-to-br from-blue-50 to-blue-100 rounded-lg border border-blue-200 p-6">
            <div className="flex items-center mb-4">
              <span className="text-4xl mr-3">₹</span>
              <h3 className="text-xl font-bold text-blue-900">Razorpay</h3>
            </div>
            <ul className="text-blue-900 text-sm space-y-2">
              <li>✅ All major banks (Netbanking)</li>
              <li>✅ Credit & Debit Cards</li>
              <li>✅ UPI & Digital Wallets</li>
              <li>✅ Recurring Payments</li>
              <li>✅ Instant Settlement</li>
            </ul>
          </div>

          {/* Billdesk Info */}
          <div className="bg-gradient-to-br from-orange-50 to-orange-100 rounded-lg border border-orange-200 p-6">
            <div className="flex items-center mb-4">
              <span className="text-4xl mr-3">🏧</span>
              <h3 className="text-xl font-bold text-orange-900">Billdesk</h3>
            </div>
            <ul className="text-orange-900 text-sm space-y-2">
              <li>✅ 100+ Banks</li>
              <li>✅ Government Payments</li>
              <li>✅ Multi-currency Support</li>
              <li>✅ High Success Rate</li>
              <li>✅ Bulk Payment Features</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  )
}
