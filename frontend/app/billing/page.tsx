'use client';

import { useEffect, useState } from 'react';

interface BillingData {
  currentOverages: number;
  costPerOverage: number;
  totalOverageCost: number;
  projectedMonthlyOverage: number;
  overageHistory: Array<{
    date: string;
    overageCount: number;
    cost: number;
  }>;
}

export default function BillingOverview() {
  const [billingData, setBillingData] = useState<BillingData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadBillingData();
    const interval = setInterval(loadBillingData, 60000); // Refresh every minute
    return () => clearInterval(interval);
  }, []);

  const loadBillingData = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch('/api/v1/tenant/users/billing/check', {
        headers: { Authorization: `Bearer ${localStorage.getItem('authToken')}` },
      });

      if (!response.ok) throw new Error('Failed to load billing data');

      const data = await response.json();
      setBillingData(data);
    } catch (err) {
      setError('Failed to load billing information');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (!billingData) {
    return <div>No billing data available</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Billing & Overage Charges</h1>
          <p className="text-gray-600 mt-2">Monitor your current and projected overage costs</p>
        </div>

        {/* Error Message */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        {/* Key Metrics Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          {/* Current Overages */}
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-gray-600 text-sm font-medium">Current Overages</p>
            <p className="text-3xl font-bold text-red-600 mt-2">{billingData.currentOverages}</p>
            <p className="text-gray-500 text-xs mt-2">Users exceeding limit</p>
          </div>

          {/* Cost per Overage */}
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-gray-600 text-sm font-medium">Cost per Overage User</p>
            <p className="text-3xl font-bold text-blue-600 mt-2">
              {formatCurrency(billingData.costPerOverage)}
            </p>
            <p className="text-gray-500 text-xs mt-2">Monthly charge per user</p>
          </div>

          {/* Total Overage Cost */}
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-gray-600 text-sm font-medium">Total Overage Cost (YTD)</p>
            <p className="text-3xl font-bold text-orange-600 mt-2">
              {formatCurrency(billingData.totalOverageCost)}
            </p>
            <p className="text-gray-500 text-xs mt-2">Year-to-date charges</p>
          </div>

          {/* Projected Monthly */}
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-gray-600 text-sm font-medium">Projected This Month</p>
            <p className="text-3xl font-bold text-red-500 mt-2">
              {formatCurrency(billingData.projectedMonthlyOverage)}
            </p>
            <p className="text-gray-500 text-xs mt-2">Based on current overages</p>
          </div>
        </div>

        {/* Alerts */}
        {billingData.currentOverages > 0 && (
          <div className="mb-6 bg-red-50 border border-red-200 rounded-lg p-4 flex items-start">
            <div className="flex-shrink-0">
              <svg
                className="h-5 w-5 text-red-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 9v2m0 4v2m0 4v2M6.228 6.228a9 9 0 1112.544 12.544M6.228 6.228l12.544 12.544"
                />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-red-800">Active Overage Charges</h3>
              <p className="text-sm text-red-700 mt-1">
                You have {billingData.currentOverages} users exceeding your seat limit.
                {billingData.projectedMonthlyOverage > 0 && (
                  <> Your projected overage charge this month is{' '}
                  <strong>{formatCurrency(billingData.projectedMonthlyOverage)}</strong>.</>
                )}
              </p>
              <a
                href="/seat-management"
                className="text-sm font-medium text-red-700 hover:text-red-600 mt-2 inline-block"
              >
                Increase Seat Limit →
              </a>
            </div>
          </div>
        )}

        {/* Overage History */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-6">Overage History</h2>

          {billingData.overageHistory.length === 0 ? (
            <div className="text-center py-8">
              <p className="text-gray-500">No overage history</p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b border-gray-200">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Date
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Overage Users
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Charge
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {billingData.overageHistory.map((entry, idx) => (
                    <tr key={idx} className="hover:bg-gray-50 transition">
                      <td className="px-6 py-4 text-sm text-gray-600">
                        {formatDate(entry.date)}
                      </td>
                      <td className="px-6 py-4 text-sm font-medium text-red-600">
                        {entry.overageCount} users
                      </td>
                      <td className="px-6 py-4 text-sm font-semibold text-gray-900">
                        {formatCurrency(entry.cost)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>

        {/* Billing Information */}
        <div className="mt-8 grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* How Billing Works */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">How Billing Works</h2>
            <ul className="space-y-3 text-sm text-gray-600">
              <li className="flex items-start">
                <span className="flex-shrink-0 h-5 w-5 text-blue-500 mr-3">✓</span>
                <span>Monthly billing cycle runs from the 1st to last day of the month</span>
              </li>
              <li className="flex items-start">
                <span className="flex-shrink-0 h-5 w-5 text-blue-500 mr-3">✓</span>
                <span>Overage charges are calculated daily based on exceeding your seat limit</span>
              </li>
              <li className="flex items-start">
                <span className="flex-shrink-0 h-5 w-5 text-blue-500 mr-3">✓</span>
                <span>Invoices are generated on the last day of the month</span>
              </li>
              <li className="flex items-start">
                <span className="flex-shrink-0 h-5 w-5 text-blue-500 mr-3">✓</span>
                <span>Payment is due within 30 days of invoice date</span>
              </li>
            </ul>
          </div>

          {/* Ways to Reduce Charges */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Reduce Overage Charges</h2>
            <ul className="space-y-3 text-sm text-gray-600">
              <li className="flex items-start">
                <span className="flex-shrink-0 h-5 w-5 text-green-500 mr-3">→</span>
                <span>Increase your seat limit to accommodate more users</span>
              </li>
              <li className="flex items-start">
                <span className="flex-shrink-0 h-5 w-5 text-green-500 mr-3">→</span>
                <span>Remove inactive users from your account</span>
              </li>
              <li className="flex items-start">
                <span className="flex-shrink-0 h-5 w-5 text-green-500 mr-3">→</span>
                <span>Upgrade your subscription tier for higher limits</span>
              </li>
              <li className="flex items-start">
                <span className="flex-shrink-0 h-5 w-5 text-green-500 mr-3">→</span>
                <span>Contact support to discuss volume discounts</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}
