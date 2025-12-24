'use client';

import { useEffect, useState } from 'react';

interface SeatManagementData {
  currentLimit: number;
  maxAllowed: number;
  currentUsers: number;
  costPerExtraUser: number;
  overageCost: number;
}

export default function SeatManagement() {
  const [seatData, setSeatData] = useState<SeatManagementData | null>(null);
  const [newLimit, setNewLimit] = useState<number>(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    loadSeatData();
  }, []);

  const loadSeatData = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch('/api/v1/tenant/users/seats', {
        headers: { Authorization: `Bearer ${localStorage.getItem('authToken')}` },
      });

      if (!response.ok) throw new Error('Failed to load seat data');

      const data = await response.json();
      setSeatData(data);
      setNewLimit(data.currentLimit);
    } catch (err) {
      setError('Failed to load seat management data');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateLimit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newLimit < seatData!.currentUsers) {
      setError(
        `New limit must be at least ${seatData!.currentUsers} (current number of users)`
      );
      return;
    }

    try {
      setSubmitting(true);
      setError(null);
      setSuccess(null);

      const response = await fetch('/api/v1/tenant/users/limit', {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('authToken')}`,
        },
        body: JSON.stringify({ maxUsersAllowed: newLimit }),
      });

      if (!response.ok) throw new Error('Failed to update seat limit');

      setSuccess('Seat limit updated successfully');
      loadSeatData();
    } catch (err) {
      setError('Failed to update seat limit');
      console.error(err);
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (!seatData) {
    return <div>No seat data available</div>;
  }

  const projectedMonthlyOverage = seatData.overageCost;
  const utilizationPercentage = Math.round(
    (seatData.currentUsers / seatData.currentLimit) * 100
  );

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Seat Management</h1>
          <p className="text-gray-600 mt-2">Manage your user seat limits and monitor overage costs</p>
        </div>

        {/* Current Status Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          {/* Current Users */}
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-gray-600 text-sm font-medium">Current Users</p>
            <p className="text-4xl font-bold text-blue-600 mt-2">{seatData.currentUsers}</p>
            <p className="text-gray-500 text-xs mt-2">Active user accounts</p>
          </div>

          {/* Current Limit */}
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-gray-600 text-sm font-medium">Seat Limit</p>
            <p className="text-4xl font-bold text-green-600 mt-2">{seatData.currentLimit}</p>
            <p className="text-gray-500 text-xs mt-2">Maximum allowed users</p>
          </div>

          {/* Utilization */}
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-gray-600 text-sm font-medium">Utilization</p>
            <p
              className={`text-4xl font-bold mt-2 ${
                utilizationPercentage > 80 ? 'text-yellow-600' : 'text-gray-900'
              }`}
            >
              {utilizationPercentage}%
            </p>
            <p className="text-gray-500 text-xs mt-2">Of available seats</p>
          </div>
        </div>

        {/* Alerts */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        {success && (
          <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded mb-6">
            {success}
          </div>
        )}

        {utilizationPercentage > 80 && (
          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-6 flex items-start">
            <div className="flex-shrink-0">
              <svg
                className="h-5 w-5 text-yellow-400"
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
              <h3 className="text-sm font-medium text-yellow-800">High Seat Utilization</h3>
              <p className="text-sm text-yellow-700 mt-1">
                You're using {utilizationPercentage}% of your available seats. Consider increasing your seat limit to avoid overage charges.
              </p>
            </div>
          </div>
        )}

        {/* Update Seat Limit Form */}
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h2 className="text-lg font-semibold text-gray-900 mb-6">Update Seat Limit</h2>

          <form onSubmit={handleUpdateLimit} className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                New Seat Limit
              </label>
              <div className="relative">
                <input
                  type="number"
                  min={seatData.currentUsers}
                  max={seatData.maxAllowed}
                  value={newLimit}
                  onChange={(e) => setNewLimit(parseInt(e.target.value) || 0)}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                <p className="text-xs text-gray-500 mt-1">
                  Minimum: {seatData.currentUsers} | Maximum: {seatData.maxAllowed}
                </p>
              </div>
            </div>

            {/* Cost Impact */}
            {newLimit !== seatData.currentLimit && (
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <p className="text-sm text-gray-700 mb-2">
                  <strong>Cost Impact:</strong>
                </p>
                {newLimit > seatData.currentLimit ? (
                  <p className="text-sm text-blue-700">
                    Increasing seats from {seatData.currentLimit} to {newLimit} (+
                    {newLimit - seatData.currentLimit} seats)
                  </p>
                ) : (
                  <p className="text-sm text-blue-700">
                    Decreasing seats from {seatData.currentLimit} to {newLimit} (−
                    {seatData.currentLimit - newLimit} seats)
                  </p>
                )}
              </div>
            )}

            <button
              type="submit"
              disabled={submitting || newLimit === seatData.currentLimit}
              className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg transition"
            >
              {submitting ? 'Updating...' : 'Update Seat Limit'}
            </button>
          </form>
        </div>

        {/* Overage Information */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-6">Overage Information</h2>

          <div className="space-y-4">
            <div className="border-b border-gray-200 pb-4">
              <div className="flex items-center justify-between">
                <span className="text-gray-700">Cost per Additional User</span>
                <span className="text-lg font-semibold text-gray-900">
                  ${seatData.costPerExtraUser}/month
                </span>
              </div>
            </div>

            <div className="border-b border-gray-200 pb-4">
              <div className="flex items-center justify-between">
                <span className="text-gray-700">Current Overage Cost (Projected)</span>
                <span className="text-lg font-semibold text-red-600">
                  ${projectedMonthlyOverage}/month
                </span>
              </div>
            </div>

            <div className="bg-gray-50 rounded-lg p-4">
              <h3 className="text-sm font-medium text-gray-900 mb-2">How Overage Works</h3>
              <ul className="text-sm text-gray-600 space-y-1">
                <li>• Overage charges apply when users exceed your seat limit</li>
                <li>• Charged at ${seatData.costPerExtraUser} per additional user per month</li>
                <li>• Charges are calculated daily and billed at month end</li>
                <li>• Increase your seat limit to eliminate overage charges</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
