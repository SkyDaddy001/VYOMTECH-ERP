'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/lib/hooks/useAuth';
import tenantService, { type Tenant } from '@/services/tenant-service';
import Link from 'next/link';

export default function BillingDashboard() {
  const { user, isLoading: authLoading } = useAuth();
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

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
      setError(err instanceof Error ? err.message : 'Failed to load data');
    } finally {
      setIsLoading(false);
    }
  };

  if (authLoading || isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  const monthlyCount = 12; // Sample data
  const usagePercent = tenant ? Math.round((tenant.activeUsers / tenant.maxUsers) * 100) : 0;

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Billing Dashboard</h1>
            <p className="text-gray-600 mt-2">Subscription and payment overview</p>
          </div>
          <Link
            href="/billing"
            className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-6 rounded-lg transition"
          >
            View Invoices
          </Link>
        </div>

        {/* Error Alert */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
            {error}
          </div>
        )}

        {/* Key Metrics */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-blue-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Monthly Cost</p>
                <p className="text-3xl font-bold text-blue-600 mt-2">$299.99</p>
              </div>
              <div className="text-4xl">💰</div>
            </div>
            <p className="text-xs text-gray-600 mt-3">Professional Plan</p>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-green-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Status</p>
                <p className="text-2xl font-bold text-green-600 mt-2">Active</p>
              </div>
              <div className="text-4xl">✅</div>
            </div>
            <p className="text-xs text-gray-600 mt-3">Next renewal: Jan 1, 2026</p>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-purple-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Invoices</p>
                <p className="text-3xl font-bold text-purple-600 mt-2">12</p>
              </div>
              <div className="text-4xl">📄</div>
            </div>
            <p className="text-xs text-gray-600 mt-3">Total: $3,599.88</p>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-orange-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Payment Method</p>
                <p className="text-lg font-bold text-orange-600 mt-2">Visa ••4242</p>
              </div>
              <div className="text-4xl">💳</div>
            </div>
            <p className="text-xs text-gray-600 mt-3">Exp: 12/26</p>
          </div>
        </div>

        {/* Plan & Usage */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Plan Details */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-6">Current Plan</h3>
            <div className="space-y-4">
              <div className="flex justify-between items-center pb-3 border-b">
                <span className="text-gray-700">Plan Type</span>
                <span className="font-bold text-gray-900">Professional</span>
              </div>
              <div className="flex justify-between items-center pb-3 border-b">
                <span className="text-gray-700">Users Included</span>
                <span className="font-bold text-gray-900">{tenant?.maxUsers}</span>
              </div>
              <div className="flex justify-between items-center pb-3 border-b">
                <span className="text-gray-700">Monthly Cost</span>
                <span className="font-bold text-gray-900">$299.99</span>
              </div>
              <div className="flex justify-between items-center pb-3 border-b">
                <span className="text-gray-700">Billing Cycle</span>
                <span className="font-bold text-gray-900">Monthly</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-gray-700">Next Billing Date</span>
                <span className="font-bold text-gray-900">Jan 1, 2026</span>
              </div>

              <button className="w-full mt-6 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg transition">
                Upgrade Plan
              </button>
            </div>
          </div>

          {/* Seat Usage */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-6">Seat Usage</h3>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between mb-3">
                  <span className="text-gray-700 font-medium">Usage</span>
                  <span className="font-bold text-gray-900">
                    {tenant?.activeUsers} / {tenant?.maxUsers}
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-4">
                  <div
                    className={`h-4 rounded-full transition-all ${
                      usagePercent > 90
                        ? 'bg-red-500'
                        : usagePercent > 80
                        ? 'bg-yellow-500'
                        : 'bg-green-500'
                    }`}
                    style={{ width: `${usagePercent}%` }}
                  />
                </div>
                <p className="text-sm text-gray-600 mt-3">
                  {usagePercent}% of available seats in use
                </p>
              </div>

              <div className="mt-6 pt-4 border-t">
                <div className="grid grid-cols-2 gap-4">
                  <div className="bg-green-50 rounded-lg p-4">
                    <p className="text-xs text-gray-600 font-medium">Available</p>
                    <p className="text-2xl font-bold text-green-600 mt-2">
                      {(tenant?.maxUsers || 0) - (tenant?.activeUsers || 0)}
                    </p>
                  </div>
                  <div className="bg-blue-50 rounded-lg p-4">
                    <p className="text-xs text-gray-600 font-medium">Total Capacity</p>
                    <p className="text-2xl font-bold text-blue-600 mt-2">{tenant?.maxUsers}</p>
                  </div>
                </div>
              </div>

              {usagePercent > 80 && (
                <div className="mt-4 bg-yellow-50 border border-yellow-200 rounded-lg p-3">
                  <p className="text-sm text-yellow-800 font-medium">
                    ⚠️ Capacity Warning: {usagePercent}% used
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Recent Transactions */}
        <div className="bg-white rounded-lg shadow overflow-hidden mb-8">
          <div className="p-6 border-b">
            <h3 className="text-lg font-semibold text-gray-900">Recent Invoices</h3>
          </div>

          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Invoice</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Date</th>
                <th className="px-6 py-3 text-right text-sm font-semibold text-gray-900">Amount</th>
                <th className="px-6 py-3 text-center text-sm font-semibold text-gray-900">Status</th>
                <th className="px-6 py-3 text-right text-sm font-semibold text-gray-900">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {[
                { id: 'INV-001', date: '2025-12-01', amount: 299.99, status: 'paid' },
                { id: 'INV-002', date: '2025-11-01', amount: 299.99, status: 'paid' },
                { id: 'INV-003', date: '2025-10-01', amount: 299.99, status: 'paid' },
              ].map((invoice) => (
                <tr key={invoice.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 font-medium text-gray-900">{invoice.id}</td>
                  <td className="px-6 py-4 text-gray-600">
                    {new Date(invoice.date).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 text-right font-medium text-gray-900">
                    ${invoice.amount.toFixed(2)}
                  </td>
                  <td className="px-6 py-4 text-center">
                    <span className="px-3 py-1 bg-green-100 text-green-800 text-sm font-medium rounded-full">
                      Paid
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button className="text-blue-600 hover:text-blue-800 font-medium text-sm">
                      Download
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          <div className="p-4 bg-gray-50 border-t text-center">
            <Link href="/billing" className="text-blue-600 hover:text-blue-800 font-medium">
              View All Invoices →
            </Link>
          </div>
        </div>

        {/* Billing History */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Billing Summary</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="bg-blue-50 rounded-lg p-4">
              <p className="text-sm text-gray-600 font-medium">Year-to-Date</p>
              <p className="text-2xl font-bold text-blue-600 mt-2">$3,599.88</p>
              <p className="text-xs text-gray-600 mt-1">12 months (Professional plan)</p>
            </div>
            <div className="bg-green-50 rounded-lg p-4">
              <p className="text-sm text-gray-600 font-medium">Average/Month</p>
              <p className="text-2xl font-bold text-green-600 mt-2">$299.99</p>
              <p className="text-xs text-gray-600 mt-1">Fixed monthly cost</p>
            </div>
            <div className="bg-purple-50 rounded-lg p-4">
              <p className="text-sm text-gray-600 font-medium">Cost Per User</p>
              <p className="text-2xl font-bold text-purple-600 mt-2">$3.00</p>
              <p className="text-xs text-gray-600 mt-1">Monthly average</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
