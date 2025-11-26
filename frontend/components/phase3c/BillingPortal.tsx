'use client'

import React, { useEffect, useState } from 'react'
import { useBillingStore } from '@/contexts/phase3cStore'
import toast from 'react-hot-toast'
import { formatDateToDDMMMYYYY } from '@/lib/dateFormat'

export function BillingPortal() {
  const { invoices, plans, charges, fetchInvoices, fetchPlans, fetchCharges, markAsPaid, loading, error } =
    useBillingStore()
  const [selectedInvoice, setSelectedInvoice] = useState<any>(null)
  const [showPaymentModal, setShowPaymentModal] = useState(false)

  useEffect(() => {
    fetchInvoices()
    fetchPlans()
    fetchCharges()
  }, [fetchInvoices, fetchPlans, fetchCharges])

  useEffect(() => {
    if (error) {
      toast.error(error)
    }
  }, [error])

  const handleMarkAsPaid = async (invoiceId: string) => {
    try {
      await markAsPaid(invoiceId, {
        payment_method: 'credit_card',
        transaction_id: `TXN-${Date.now()}`,
      })
      toast.success('Invoice marked as paid!')
      setShowPaymentModal(false)
      setSelectedInvoice(null)
    } catch (error: any) {
      toast.error(error.message || 'Failed to mark invoice as paid')
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'paid':
        return 'bg-green-100 text-green-800'
      case 'sent':
        return 'bg-blue-100 text-blue-800'
      case 'overdue':
        return 'bg-red-100 text-red-800'
      case 'draft':
        return 'bg-gray-100 text-gray-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">Billing Portal</h2>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <p className="text-gray-600 text-sm">Current Plan</p>
          <p className="text-xl font-bold text-blue-600">{plans.length > 0 ? plans[0].name : 'No Active Plan'}</p>
        </div>
        <div className="bg-green-50 border border-green-200 rounded-lg p-4">
          <p className="text-gray-600 text-sm">Monthly Charges</p>
          <p className="text-xl font-bold text-green-600">${charges?.total_amount || '0.00'}</p>
        </div>
        <div className="bg-orange-50 border border-orange-200 rounded-lg p-4">
          <p className="text-gray-600 text-sm">Outstanding</p>
          <p className="text-xl font-bold text-orange-600">
            ${invoices
              .filter((inv) => inv.status !== 'paid')
              .reduce((sum, inv) => sum + inv.total_amount, 0)
              .toFixed(2)}
          </p>
        </div>
      </div>

      {/* Invoices */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-4 border-b">
          <h3 className="text-lg font-bold">Invoices</h3>
        </div>

        {loading ? (
          <div className="p-8 text-center text-gray-500">Loading invoices...</div>
        ) : invoices.length === 0 ? (
          <div className="p-8 text-center text-gray-500">No invoices yet</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="px-4 py-2 text-left text-sm font-semibold">Invoice #</th>
                  <th className="px-4 py-2 text-left text-sm font-semibold">Amount</th>
                  <th className="px-4 py-2 text-left text-sm font-semibold">Status</th>
                  <th className="px-4 py-2 text-left text-sm font-semibold">Date</th>
                  <th className="px-4 py-2 text-left text-sm font-semibold">Actions</th>
                </tr>
              </thead>
              <tbody>
                {invoices.map((invoice: any) => (
                  <tr key={invoice.id} className="border-b hover:bg-gray-50">
                    <td className="px-4 py-3 text-sm font-mono">{invoice.invoice_number}</td>
                    <td className="px-4 py-3 text-sm font-bold">${invoice.total_amount.toFixed(2)}</td>
                    <td className="px-4 py-3 text-sm">
                      <span className={`px-2 py-1 text-xs font-semibold rounded ${getStatusColor(invoice.status)}`}>
                        {invoice.status}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {formatDateToDDMMMYYYY(invoice.created_at)}
                    </td>
                    <td className="px-4 py-3 text-sm space-x-2">
                      <button
                        onClick={() => {
                          setSelectedInvoice(invoice)
                          setShowPaymentModal(true)
                        }}
                        className="text-blue-600 hover:text-blue-800 font-semibold"
                      >
                        View
                      </button>
                      {invoice.status !== 'paid' && (
                        <button
                          onClick={() => handleMarkAsPaid(invoice.id)}
                          className="text-green-600 hover:text-green-800 font-semibold"
                        >
                          Pay
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Payment Modal */}
      {showPaymentModal && selectedInvoice && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <h3 className="text-lg font-bold mb-4">Invoice Details</h3>

            <div className="space-y-3 mb-6 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">Invoice Number:</span>
                <span className="font-mono font-bold">{selectedInvoice.invoice_number}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Amount:</span>
                <span className="font-bold text-lg text-blue-600">
                  ${selectedInvoice.total_amount.toFixed(2)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Status:</span>
                <span className={`px-2 py-1 text-xs font-semibold rounded ${getStatusColor(selectedInvoice.status)}`}>
                  {selectedInvoice.status}
                </span>
              </div>
            </div>

            <div className="flex gap-3">
              <button
                onClick={() => setShowPaymentModal(false)}
                className="flex-1 bg-gray-200 text-gray-800 px-4 py-2 rounded-lg hover:bg-gray-300"
              >
                Close
              </button>
              {selectedInvoice.status !== 'paid' && (
                <button
                  onClick={() => handleMarkAsPaid(selectedInvoice.id)}
                  disabled={loading}
                  className="flex-1 bg-green-500 text-white px-4 py-2 rounded-lg hover:bg-green-600 disabled:opacity-50"
                >
                  {loading ? 'Processing...' : 'Pay Now'}
                </button>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
