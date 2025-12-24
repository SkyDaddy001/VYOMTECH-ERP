'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api-client';
import {
  BarChart,
  Bar,
  LineChart,
  Line,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

interface UserCountData {
  tenantId: string;
  totalActiveUsers: number;
  totalInactiveUsers: number;
  totalUsers: number;
  maxUsersAllowed: number;
  seatUtilization: number;
  isOverLimit: boolean;
  overageCount: number;
  seatsRemaining: number;
  lastCountedAt: string;
}

interface UserBreakdownData {
  tenantId: string;
  breakdown: {
    adminUsers: number;
    managerUsers: number;
    standardUsers: number;
    guestUsers: number;
  };
  totalByStatus: {
    active: number;
    inactive: number;
    suspended: number;
  };
}

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6'];

export default function UserCountDashboard() {
  const [userCount, setUserCount] = useState<UserCountData | null>(null);
  const [breakdown, setBreakdown] = useState<UserBreakdownData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadUserCountData();
    const interval = setInterval(loadUserCountData, 30000); // Refresh every 30 seconds
    return () => clearInterval(interval);
  }, []);

  const loadUserCountData = async () => {
    try {
      setLoading(true);
      setError(null);

      const [countData, breakdownData] = await Promise.all([
        fetch('/api/v1/tenant/users/count', {
          headers: { Authorization: `Bearer ${localStorage.getItem('authToken')}` },
        }).then((res) => res.json()),
        fetch('/api/v1/tenant/users/breakdown', {
          headers: { Authorization: `Bearer ${localStorage.getItem('authToken')}` },
        }).then((res) => res.json()),
      ]);

      setUserCount(countData);
      setBreakdown(breakdownData);
    } catch (err) {
      setError('Failed to load user count data');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
        {error}
      </div>
    );
  }

  if (!userCount || !breakdown) {
    return <div>No data available</div>;
  }

  const utilizationPercentage = Math.round(userCount.seatUtilization);
  const isHighUtilization = utilizationPercentage > 80;
  const isOverLimit = userCount.isOverLimit;

  const breakdownChartData = [
    { name: 'Admin', value: breakdown.breakdown.adminUsers, fill: '#3b82f6' },
    { name: 'Manager', value: breakdown.breakdown.managerUsers, fill: '#10b981' },
    { name: 'Standard', value: breakdown.breakdown.standardUsers, fill: '#f59e0b' },
    { name: 'Guest', value: breakdown.breakdown.guestUsers, fill: '#8b5cf6' },
  ];

  const statusChartData = [
    { name: 'Active', value: breakdown.totalByStatus.active },
    { name: 'Inactive', value: breakdown.totalByStatus.inactive },
    { name: 'Suspended', value: breakdown.totalByStatus.suspended },
  ];

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">User Count & Seat Management</h1>
          <p className="text-gray-600 mt-2">Real-time user metrics and seat utilization</p>
        </div>

        {/* Key Metrics Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          {/* Total Users Card */}
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Total Users</p>
                <p className="text-3xl font-bold text-gray-900 mt-2">{userCount.totalUsers}</p>
              </div>
              <div className="bg-blue-100 rounded-full p-3">
                <svg
                  className="w-6 h-6 text-blue-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M17 20h5v-2a3 3 0 00-5.856-1.487M15 10a3 3 0 11-6 0 3 3 0 016 0zM6 20a9 9 0 0118 0v2h2v-2a11 11 0 00-22 0v2h2v-2z"
                  />
                </svg>
              </div>
            </div>
          </div>

          {/* Active Users Card */}
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Active Users</p>
                <p className="text-3xl font-bold text-green-600 mt-2">
                  {breakdown.totalByStatus.active}
                </p>
              </div>
              <div className="bg-green-100 rounded-full p-3">
                <svg
                  className="w-6 h-6 text-green-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              </div>
            </div>
          </div>

          {/* Seat Utilization Card */}
          <div
            className={`rounded-lg shadow p-6 ${
              isHighUtilization || isOverLimit
                ? 'bg-yellow-50 border-2 border-yellow-200'
                : 'bg-white'
            }`}
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Seat Utilization</p>
                <p
                  className={`text-3xl font-bold mt-2 ${
                    isHighUtilization || isOverLimit ? 'text-yellow-600' : 'text-gray-900'
                  }`}
                >
                  {utilizationPercentage}%
                </p>
                {isHighUtilization && !isOverLimit && (
                  <p className="text-yellow-600 text-xs mt-1">⚠️ High utilization</p>
                )}
              </div>
              <div
                className={`rounded-full p-3 ${
                  isHighUtilization || isOverLimit ? 'bg-yellow-100' : 'bg-blue-100'
                }`}
              >
                <svg
                  className={`w-6 h-6 ${
                    isHighUtilization || isOverLimit ? 'text-yellow-600' : 'text-blue-600'
                  }`}
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M13 10V3L4 14h7v7l9-11h-7z"
                  />
                </svg>
              </div>
            </div>
          </div>

          {/* Seats Remaining Card */}
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Seats Remaining</p>
                <p
                  className={`text-3xl font-bold mt-2 ${
                    userCount.seatsRemaining <= 10 ? 'text-red-600' : 'text-gray-900'
                  }`}
                >
                  {userCount.seatsRemaining}
                </p>
                <p className="text-gray-500 text-xs mt-1">
                  of {userCount.maxUsersAllowed} total
                </p>
              </div>
              <div
                className={`rounded-full p-3 ${
                  userCount.seatsRemaining <= 10 ? 'bg-red-100' : 'bg-gray-100'
                }`}
              >
                <svg
                  className={`w-6 h-6 ${
                    userCount.seatsRemaining <= 10 ? 'text-red-600' : 'text-gray-600'
                  }`}
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                  />
                </svg>
              </div>
            </div>
          </div>
        </div>

        {/* Alerts */}
        {isOverLimit && (
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
              <h3 className="text-sm font-medium text-red-800">Over Seat Limit</h3>
              <p className="text-sm text-red-700 mt-1">
                You have {userCount.overageCount} users exceeding your seat limit. Overage charges may apply.
              </p>
            </div>
          </div>
        )}

        {/* Charts */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* User Breakdown by Role */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">User Breakdown by Role</h2>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={breakdownChartData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, value }) => `${name}: ${value}`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {breakdownChartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.fill} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>

          {/* User Status Distribution */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">User Status Distribution</h2>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={statusChartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="value" fill="#3b82f6" radius={[8, 8, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Seat Utilization Progress */}
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Seat Utilization Progress</h2>
          <div className="space-y-4">
            <div>
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium text-gray-700">Used Seats</span>
                <span className="text-sm font-semibold text-gray-900">
                  {userCount.totalUsers} / {userCount.maxUsersAllowed}
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-3">
                <div
                  className={`h-3 rounded-full transition-all ${
                    isOverLimit
                      ? 'bg-red-500'
                      : isHighUtilization
                        ? 'bg-yellow-500'
                        : 'bg-green-500'
                  }`}
                  style={{ width: `${Math.min(utilizationPercentage, 100)}%` }}
                />
              </div>
              {isOverLimit && (
                <p className="text-xs text-red-600 mt-1">
                  {userCount.overageCount} users over limit
                </p>
              )}
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <button className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg transition">
            View Activity Log
          </button>
          <button className="bg-gray-600 hover:bg-gray-700 text-white font-semibold py-2 px-4 rounded-lg transition">
            Manage Seat Limit
          </button>
          <button className="bg-green-600 hover:bg-green-700 text-white font-semibold py-2 px-4 rounded-lg transition">
            Export Report
          </button>
        </div>
      </div>
    </div>
  );
}
