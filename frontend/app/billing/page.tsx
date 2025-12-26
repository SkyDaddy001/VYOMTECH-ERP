'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/lib/hooks/useAuth';
import tenantService, { type Tenant } from '@/services/tenant-service';

interface Invoice {
  id: string;
  date: string;
  amount: number;
  status: 'paid' | 'pending' | 'failed';
  description: string;
}

export default function BillingPage() {
  const { user, isLoading: authLoading } = useAuth();
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [invoices, setInvoices] = useState<Invoice[]>([
    {
      id: 'INV-001',
      date: '2025-12-01',
      amount: 299.99,
      status: 'paid',
      description: 'Monthly subscription - 100 users',
    },
    {
      id: 'INV-002',
      date: '2025-11-01',
      amount: 299.99,
      status: 'paid',
      description: 'Monthly subscription - 100 users',
    },
  ]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showPaymentModal, setShowPaymentModal] = useState(false);

  useEffect(() => {
    if (authLoading || !user) return;
    loadData();
  }, [authLoading, user]);

  const loadData = async () => {
    try {
      setIsLoading(true);
      const tenantData = await tenantService.getCurrentTenant();
      setTenant(tenantData);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load billing data');
    } finally {
      setIsLoading(false);
    }
  };

  const statusColor = (status: string) => {
    switch (status) {
      case 'paid':
        return 'bg-green-100 text-green-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  if (authLoading || isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Billing & Invoices</h1>
          <p className="text-gray-600 mt-2">Manage your subscription and payments</p>
        </div>

        {/* Error Alert */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
            {error}
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
          {/* Current Plan */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="font-semibold text-gray-900 mb-4">Current Plan</h3>
            <div className="space-y-3">
              <div>
                <p className="text-sm text-gray-600">Plan Type</p>
                <p className="text-xl font-bold text-gray-900">Professional</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Users Included</p>
                <p className="text-xl font-bold text-gray-900">{tenant?.maxUsers}</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Monthly Cost</p>
                <p className="text-xl font-bold text-blue-600">$299.99</p>
              </div>
              <button className="w-full mt-4 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg transition">
                Upgrade Plan
              </button>
            </div>
          </div>

          {/* Usage */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="font-semibold text-gray-900 mb-4">Usage</h3>
            <div className="space-y-3">
              <div>
                <p className="text-sm text-gray-600">Active Users</p>
                <p className="text-xl font-bold text-gray-900">{tenant?.activeUsers}</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Seats Available</p>
                <p className="text-xl font-bold text-green-600">
                  {(tenant?.maxUsers || 0) - (tenant?.activeUsers || 0)}
                </p>
              </div>
              <div className="mt-4">
                <div className="flex justify-between mb-2">
                  <p className="text-sm text-gray-600">Usage</p>
                  <p className="text-sm font-medium text-gray-900">
                    {Math.round(((tenant?.activeUsers || 0) / (tenant?.maxUsers || 1)) * 100)}%
                  </p>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-blue-600 h-2 rounded-full transition"
                    style={{
                      width: `${Math.round(((tenant?.activeUsers || 0) / (tenant?.maxUsers || 1)) * 100)}%`,
                    }}
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Payment Method */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="font-semibold text-gray-900 mb-4">Payment Method</h3>
            <div className="space-y-3">
              <div className="bg-gray-100 rounded-lg p-4 flex items-center gap-3">
                <div className="text-2xl">💳</div>
                <div>
                  <p className="font-medium text-gray-900">Visa</p>
                  <p className="text-sm text-gray-600">**** **** **** 4242</p>
                </div>
              </div>
              <button className="w-full mt-4 bg-gray-200 hover:bg-gray-300 text-gray-900 font-semibold py-2 px-4 rounded-lg transition">
                Update Method
              </button>
            </div>
          </div>
        </div>

        {/* Recent Invoices */}
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <div className="p-6 border-b">
            <h3 className="text-lg font-semibold text-gray-900">Recent Invoices</h3>
          </div>

          {invoices.length === 0 ? (
            <div className="p-8 text-center text-gray-600">No invoices found</div>
          ) : (
            <table className="w-full">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                    Invoice
                  </th>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                    Date
                  </th>
                  <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">
                    Description
                  </th>
                  <th className="px-6 py-3 text-right text-sm font-semibold text-gray-900">
                    Amount
                  </th>
                  <th className="px-6 py-3 text-center text-sm font-semibold text-gray-900">
                    Status
                  </th>
                  <th className="px-6 py-3 text-right text-sm font-semibold text-gray-900">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody>
                {invoices.map((invoice) => (
                  <tr key={invoice.id} className="border-b hover:bg-gray-50">
                    <td className="px-6 py-4 font-medium text-gray-900">{invoice.id}</td>
                    <td className="px-6 py-4 text-gray-600">
                      {new Date(invoice.date).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4 text-gray-600">{invoice.description}</td>
                    <td className="px-6 py-4 text-right font-medium text-gray-900">
                      ${invoice.amount.toFixed(2)}
                    </td>
                    <td className="px-6 py-4 text-center">
                      <span className={`px-3 py-1 rounded-full text-sm font-medium ${statusColor(invoice.status)}`}>
                        {invoice.status.charAt(0).toUpperCase() + invoice.status.slice(1)}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-right space-x-2">
                      <button className="text-blue-600 hover:text-blue-800 font-medium">
                        Download
                      </button>
                      {invoice.status === 'pending' && (
                        <button
                          onClick={() => setShowPaymentModal(true)}
                          className="text-green-600 hover:text-green-800 font-medium"
                        >
                          Pay
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>

      {/* Payment Modal */}
      {showPaymentModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg shadow-xl p-8 max-w-md w-full">
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Make a Payment</h2>

            <form
              onSubmit={(e) => {
                e.preventDefault();
                setShowPaymentModal(false);
              }}
              className="space-y-6"
            >
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Card Number
                </label>
                <input
                  type="text"
                  placeholder="4242 4242 4242 4242"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Expiry
                  </label>
                  <input
                    type="text"
                    placeholder="MM/YY"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    CVC
                  </label>
                  <input
                    type="text"
                    placeholder="123"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  type="submit"
                  className="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg transition"
                >
                  Pay Now
                </button>
                <button
                  type="button"
                  onClick={() => setShowPaymentModal(false)}
                  className="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-900 font-semibold py-2 px-4 rounded-lg transition"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

