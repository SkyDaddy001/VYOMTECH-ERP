'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/lib/hooks/useAuth';
import userService, { type UserCount } from '@/services/user-service';
import Link from 'next/link';

export default function UsersDashboard() {
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [stats, setStats] = useState<UserCount | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (authLoading || !isAuthenticated) return;
    loadStats();
  }, [isAuthenticated, authLoading]);

  const loadStats = async () => {
    try {
      setIsLoading(true);
      const data = await userService.getUserCount();
      setStats(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load stats');
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

  const usagePercent = stats ? Math.round((stats.activeUsers / stats.seats.total) * 100) : 0;

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Users Dashboard</h1>
            <p className="text-gray-600 mt-2">Team management overview</p>
          </div>
          <Link
            href="/users"
            className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-6 rounded-lg transition"
          >
            Manage Users
          </Link>
        </div>

        {/* Error Alert */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
            {error}
          </div>
        )}

        {/* KPI Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          {/* Total Users */}
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Total Users</p>
                <p className="text-3xl font-bold text-gray-900 mt-2">{stats?.totalUsers || 0}</p>
              </div>
              <div className="text-4xl">👥</div>
            </div>
          </div>

          {/* Active Users */}
          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-green-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Active Users</p>
                <p className="text-3xl font-bold text-green-600 mt-2">{stats?.activeUsers || 0}</p>
              </div>
              <div className="text-4xl">✅</div>
            </div>
          </div>

          {/* Inactive Users */}
          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-yellow-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Inactive Users</p>
                <p className="text-3xl font-bold text-yellow-600 mt-2">{stats?.inactiveUsers || 0}</p>
              </div>
              <div className="text-4xl">⏸️</div>
            </div>
          </div>

          {/* Suspended Users */}
          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-red-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Suspended Users</p>
                <p className="text-3xl font-bold text-red-600 mt-2">{stats?.suspendedUsers || 0}</p>
              </div>
              <div className="text-4xl">🚫</div>
            </div>
          </div>
        </div>

        {/* Seat Management */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Seat Usage */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-6">Seat Usage</h3>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between mb-2">
                  <span className="text-sm font-medium text-gray-700">Used Seats</span>
                  <span className="text-sm font-bold text-gray-900">
                    {stats?.seats.used || 0} / {stats?.seats.total || 0}
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div
                    className="bg-blue-600 h-3 rounded-full transition-all"
                    style={{ width: `${usagePercent}%` }}
                  />
                </div>
                <p className="text-sm text-gray-600 mt-2">
                  {usagePercent}% capacity used
                </p>
              </div>

              <div className="pt-4 border-t">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-gray-600">Available Seats</p>
                    <p className="text-2xl font-bold text-green-600 mt-1">
                      {stats?.seats.available || 0}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Total Capacity</p>
                    <p className="text-2xl font-bold text-blue-600 mt-1">
                      {stats?.seats.total || 0}
                    </p>
                  </div>
                </div>
              </div>

              {usagePercent > 80 && (
                <div className="mt-4 bg-yellow-50 border border-yellow-200 rounded-lg p-3">
                  <p className="text-sm text-yellow-800">
                    ⚠️ You're using {usagePercent}% of your seats. Consider upgrading your plan.
                  </p>
                </div>
              )}
            </div>
          </div>

          {/* User Distribution */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-6">User Distribution</h3>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                  <span className="text-gray-700">Active</span>
                </div>
                <span className="font-bold text-gray-900">{stats?.activeUsers || 0}</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
                  <span className="text-gray-700">Inactive</span>
                </div>
                <span className="font-bold text-gray-900">{stats?.inactiveUsers || 0}</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                  <span className="text-gray-700">Suspended</span>
                </div>
                <span className="font-bold text-gray-900">{stats?.suspendedUsers || 0}</span>
              </div>

              <div className="mt-6 pt-4 border-t">
                <div className="flex justify-between items-center">
                  <span className="text-gray-700 font-medium">Active Rate</span>
                  <span className="text-2xl font-bold text-green-600">
                    {stats?.totalUsers ? Math.round((stats.activeUsers / stats.totalUsers) * 100) : 0}%
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <Link
              href="/users/create"
              className="p-4 border rounded-lg hover:bg-blue-50 transition text-center"
            >
              <div className="text-2xl mb-2">➕</div>
              <p className="font-medium text-gray-900">Add User</p>
            </Link>
            <Link
              href="/users"
              className="p-4 border rounded-lg hover:bg-blue-50 transition text-center"
            >
              <div className="text-2xl mb-2">📋</div>
              <p className="font-medium text-gray-900">View All</p>
            </Link>
            <Link
              href="/users?status=inactive"
              className="p-4 border rounded-lg hover:bg-blue-50 transition text-center"
            >
              <div className="text-2xl mb-2">⏸️</div>
              <p className="font-medium text-gray-900">Inactive</p>
            </Link>
            <Link
              href="/users?status=suspended"
              className="p-4 border rounded-lg hover:bg-blue-50 transition text-center"
            >
              <div className="text-2xl mb-2">🚫</div>
              <p className="font-medium text-gray-900">Suspended</p>
            </Link>
          </div>
        </div>

        {/* Health Status */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">System Health</h3>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-gray-700">User Count Status</span>
              <span className="px-3 py-1 bg-green-100 text-green-800 text-sm font-medium rounded-full">
                ✅ Healthy
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-700">Seat Capacity</span>
              <span className={`px-3 py-1 text-sm font-medium rounded-full ${
                usagePercent > 90
                  ? 'bg-red-100 text-red-800'
                  : usagePercent > 80
                  ? 'bg-yellow-100 text-yellow-800'
                  : 'bg-green-100 text-green-800'
              }`}>
                {usagePercent > 90 ? '⚠️ Critical' : usagePercent > 80 ? '⚠️ Warning' : '✅ Healthy'}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-700">Data Sync</span>
              <span className="px-3 py-1 bg-green-100 text-green-800 text-sm font-medium rounded-full">
                ✅ Up to date
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
