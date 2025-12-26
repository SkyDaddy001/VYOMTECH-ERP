'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/lib/hooks/useAuth';
import userService, { type UserCount, type UserActivity } from '@/services/user-service';
import tenantService, { type Tenant } from '@/services/tenant-service';
import Link from 'next/link';

// Sidebar Component
function Sidebar() {
  return (
    <div className="fixed left-0 top-0 h-full w-64 bg-gray-900 text-white shadow-2xl overflow-y-auto">
      {/* Logo Area */}
      <div className="p-6 border-b border-gray-800">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg flex items-center justify-center font-bold text-lg">
            V
          </div>
          <div>
            <h1 className="font-bold text-lg leading-tight">VYOM</h1>
            <p className="text-xs text-gray-400">Enterprise</p>
          </div>
        </div>
      </div>

      {/* Navigation Menu */}
      <nav className="p-6 space-y-2">
        <SidebarItem href="/dashboard" icon="📊" label="Dashboard" active />
        <SidebarItem href="/users/dashboard" icon="👥" label="Users" />
        <SidebarItem href="/billing/dashboard" icon="💰" label="Billing" />
        <SidebarItem href="/activity-logs/dashboard" icon="📋" label="Activity Logs" />
        <SidebarItem href="/settings/dashboard" icon="⚙️" label="Settings" />
      </nav>

      {/* Footer Info */}
      <div className="absolute bottom-0 left-0 right-0 p-6 border-t border-gray-800 bg-gray-800 bg-opacity-50">
        <div className="text-xs text-gray-400 space-y-1">
          <p>Version 1.0.0</p>
          <p className="text-gray-500">© 2025 VYOM ERP</p>
        </div>
      </div>
    </div>
  );
}

function SidebarItem({ href, icon, label, active = false }: { href: string; icon: string; label: string; active?: boolean }) {
  return (
    <Link
      href={href}
      className={`flex items-center gap-3 px-4 py-3 rounded-lg transition ${
        active ? 'bg-blue-600 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'
      }`}
    >
      <span className="text-lg">{icon}</span>
      <span className="font-medium">{label}</span>
    </Link>
  );
}

// Top Header Component
function TopHeader({ userName, tenantName }: { userName: string; tenantName: string }) {
  return (
    <div className="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-8 shadow-sm">
      <div>
        <h2 className="text-xl font-bold text-gray-900">Dashboard</h2>
        <p className="text-xs text-gray-500">{tenantName}</p>
      </div>
      <div className="flex items-center gap-4">
        <div className="text-right">
          <p className="text-sm font-medium text-gray-900">{userName}</p>
          <p className="text-xs text-gray-500">Admin</p>
        </div>
        <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-blue-600 rounded-full flex items-center justify-center text-white font-bold">
          {userName.charAt(0).toUpperCase()}
        </div>
      </div>
    </div>
  );
}

// Stat Card Component
function StatCard({
  title,
  value,
  icon,
  color,
  trend,
  href,
}: {
  title: string;
  value: string | number;
  icon: string;
  color: string;
  trend?: string;
  href: string;
}) {
  return (
    <Link href={href}>
      <div className={`bg-white rounded-xl shadow-md hover:shadow-lg transition p-6 border-l-4 ${color} cursor-pointer`}>
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <p className="text-gray-600 text-sm font-medium">{title}</p>
            <p className="text-3xl font-bold text-gray-900 mt-2">{value}</p>
            {trend && <p className="text-xs text-green-600 mt-2">{trend}</p>}
          </div>
          <div className="text-4xl opacity-80">{icon}</div>
        </div>
      </div>
    </Link>
  );
}

// Chart Card Component (Placeholder)
function ChartCard({ title, children, action }: { title: string; children: React.ReactNode; action?: React.ReactNode }) {
  return (
    <div className="bg-white rounded-xl shadow-md p-6">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-bold text-gray-900">{title}</h3>
        {action}
      </div>
      <div>{children}</div>
    </div>
  );
}

export default function Dashboard() {
  const { user, isLoading: authLoading } = useAuth();
  const [userStats, setUserStats] = useState<UserCount | null>(null);
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [activities, setActivities] = useState<UserActivity[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (authLoading || !user) return;
    loadData();
  }, [authLoading, user]);

  const loadData = async () => {
    try {
      setIsLoading(true);
      const [stats, tenantData, activityData] = await Promise.all([
        userService.getUserCount(),
        tenantService.getCurrentTenant(),
        userService.getUserActivity({ days: 7, limit: 10 }),
      ]);
      setUserStats(stats);
      setTenant(tenantData);
      setActivities(activityData);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load dashboard data');
    } finally {
      setIsLoading(false);
    }
  };

  if (authLoading || isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="text-center">
          <div className="inline-block animate-spin">
            <div className="w-12 h-12 rounded-full border-4 border-gray-200 border-t-blue-600" />
          </div>
          <p className="text-gray-600 mt-4">Loading dashboard...</p>
        </div>
      </div>
    );
  }

  const seatUsagePercent = userStats
    ? Math.round((userStats.activeUsers / userStats.seats.total) * 100)
    : 0;
  const userName = user?.email?.split('@')[0] || 'User';
  const tenantName = tenant?.name || 'Organization';

  return (
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar */}
      <Sidebar />

      {/* Main Content */}
      <div className="flex-1 flex flex-col ml-64 overflow-hidden">
        {/* Top Header */}
        <TopHeader userName={userName} tenantName={tenantName} />

        {/* Main Content Area */}
        <div className="flex-1 overflow-y-auto">
          <div className="p-8">
            {/* Error Alert */}
            {error && (
              <div className="mb-6 bg-red-50 border-l-4 border-red-500 p-4 rounded">
                <p className="text-red-800 font-medium">Error</p>
                <p className="text-red-700 text-sm">{error}</p>
              </div>
            )}

            {/* Welcome Section */}
            <div className="mb-8">
              <h1 className="text-3xl font-bold text-gray-900">Welcome back! 👋</h1>
              <p className="text-gray-600 text-sm mt-1">
                Here's what's happening in your organization today.
              </p>
            </div>

            {/* KPI Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
              <StatCard
                title="Total Users"
                value={userStats?.totalUsers || 0}
                icon="👥"
                color="border-blue-500"
                trend={`↑ ${userStats?.activeUsers || 0} active`}
                href="/users/dashboard"
              />
              <StatCard
                title="Seat Usage"
                value={`${userStats?.seats.used || 0}/${userStats?.seats.total || 0}`}
                icon="💺"
                color="border-green-500"
                trend={`${seatUsagePercent}% capacity`}
                href="/billing/dashboard"
              />
              <StatCard
                title="Activities (7d)"
                value={activities.length}
                icon="📊"
                color="border-purple-500"
                trend={`${Math.round(activities.length / 7)}/day avg`}
                href="/activity-logs/dashboard"
              />
              <StatCard
                title="Monthly Spend"
                value="$299.99"
                icon="💰"
                color="border-orange-500"
                trend="Professional Plan"
                href="/billing/dashboard"
              />
            </div>

            {/* Main Grid */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
              {/* Users Module - Spans 2 columns */}
              <div className="lg:col-span-2">
                <ChartCard
                  title="User Management"
                  action={
                    <Link href="/users/dashboard" className="text-blue-600 hover:text-blue-800 text-sm font-medium">
                      View All →
                    </Link>
                  }
                >
                  <div className="grid grid-cols-3 gap-4">
                    <div className="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-4">
                      <p className="text-gray-600 text-sm">Active Users</p>
                      <p className="text-2xl font-bold text-green-600 mt-1">{userStats?.activeUsers || 0}</p>
                    </div>
                    <div className="bg-gradient-to-br from-yellow-50 to-yellow-100 rounded-lg p-4">
                      <p className="text-gray-600 text-sm">Inactive</p>
                      <p className="text-2xl font-bold text-yellow-600 mt-1">{userStats?.inactiveUsers || 0}</p>
                    </div>
                    <div className="bg-gradient-to-br from-red-50 to-red-100 rounded-lg p-4">
                      <p className="text-gray-600 text-sm">Suspended</p>
                      <p className="text-2xl font-bold text-red-600 mt-1">{userStats?.suspendedUsers || 0}</p>
                    </div>
                  </div>
                </ChartCard>
              </div>

              {/* Quick Actions */}
              <div>
                <ChartCard title="Quick Actions">
                  <div className="space-y-3">
                    <Link
                      href="/users/create"
                      className="block w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg transition text-center"
                    >
                      + Add User
                    </Link>
                    <Link
                      href="/settings/dashboard"
                      className="block w-full bg-purple-600 hover:bg-purple-700 text-white font-medium py-2 px-4 rounded-lg transition text-center"
                    >
                      ⚙️ Settings
                    </Link>
                    <Link
                      href="/billing/dashboard"
                      className="block w-full bg-orange-600 hover:bg-orange-700 text-white font-medium py-2 px-4 rounded-lg transition text-center"
                    >
                      💳 Billing
                    </Link>
                  </div>
                </ChartCard>
              </div>
            </div>

            {/* Recent Activity & Billing Info */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Recent Activity - Spans 2 columns */}
              <div className="lg:col-span-2">
                <ChartCard
                  title="Recent Activity"
                  action={
                    <Link href="/activity-logs/dashboard" className="text-blue-600 hover:text-blue-800 text-sm font-medium">
                      View All →
                    </Link>
                  }
                >
                  {activities.length === 0 ? (
                    <div className="text-center py-12">
                      <p className="text-gray-500 text-sm">No recent activities</p>
                    </div>
                  ) : (
                    <div className="space-y-0 divide-y divide-gray-200">
                      {activities.map((activity, index) => (
                        <div key={index} className="py-4 hover:bg-gray-50 px-4 -mx-4 transition">
                          <div className="flex items-center gap-4">
                            <div className="text-2xl flex-shrink-0">
                              {activity.activityType === 'login'
                                ? '🔓'
                                : activity.activityType === 'logout'
                                ? '🔐'
                                : activity.activityType === 'create'
                                ? '➕'
                                : activity.activityType === 'update'
                                ? '✏️'
                                : activity.activityType === 'delete'
                                ? '🗑️'
                                : '📝'}
                            </div>
                            <div className="flex-1 min-w-0">
                              <p className="font-medium text-gray-900 text-sm">{activity.email}</p>
                              <p className="text-gray-600 text-sm">
                                {activity.activityType.charAt(0).toUpperCase() + activity.activityType.slice(1)}
                              </p>
                            </div>
                            <div className="text-right flex-shrink-0">
                              <time className="text-xs text-gray-500">
                                {new Date(activity.lastActivityAt).toLocaleDateString()}
                              </time>
                              <p className="text-xs text-gray-400">
                                {new Date(activity.lastActivityAt).toLocaleTimeString()}
                              </p>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </ChartCard>
              </div>

              {/* Billing & System Status */}
              <div className="space-y-6">
                {/* Billing Card */}
                <ChartCard
                  title="Billing Status"
                  action={
                    <Link href="/billing/dashboard" className="text-orange-600 hover:text-orange-800 text-sm font-medium">
                      Manage →
                    </Link>
                  }
                >
                  <div className="space-y-4">
                    <div>
                      <div className="flex justify-between mb-2">
                        <span className="text-sm text-gray-700">Seat Capacity</span>
                        <span className="text-sm font-bold text-gray-900">{seatUsagePercent}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-3">
                        <div
                          className={`h-3 rounded-full transition ${
                            seatUsagePercent > 90
                              ? 'bg-red-500'
                              : seatUsagePercent > 80
                              ? 'bg-yellow-500'
                              : 'bg-green-500'
                          }`}
                          style={{ width: `${Math.min(seatUsagePercent, 100)}%` }}
                        />
                      </div>
                    </div>
                    <div className="border-t border-gray-200 pt-4">
                      <p className="text-gray-600 text-sm mb-1">Monthly Cost</p>
                      <p className="text-2xl font-bold text-gray-900">$299.99</p>
                      <p className="text-xs text-gray-500 mt-1">Professional Plan</p>
                    </div>
                  </div>
                </ChartCard>

                {/* System Status */}
                <ChartCard title="System Status">
                  <div className="space-y-3">
                    <div className="flex items-center justify-between p-3 bg-green-50 rounded-lg">
                      <span className="text-sm font-medium text-gray-900">API</span>
                      <span className="text-xs font-bold text-green-700 bg-green-100 px-2 py-1 rounded">✅ Healthy</span>
                    </div>
                    <div className="flex items-center justify-between p-3 bg-green-50 rounded-lg">
                      <span className="text-sm font-medium text-gray-900">Database</span>
                      <span className="text-xs font-bold text-green-700 bg-green-100 px-2 py-1 rounded">✅ Healthy</span>
                    </div>
                    <div
                      className={`flex items-center justify-between p-3 rounded-lg ${
                        seatUsagePercent > 90 ? 'bg-red-50' : seatUsagePercent > 80 ? 'bg-yellow-50' : 'bg-green-50'
                      }`}
                    >
                      <span className="text-sm font-medium text-gray-900">Capacity</span>
                      <span
                        className={`text-xs font-bold px-2 py-1 rounded ${
                          seatUsagePercent > 90
                            ? 'text-red-700 bg-red-100'
                            : seatUsagePercent > 80
                            ? 'text-yellow-700 bg-yellow-100'
                            : 'text-green-700 bg-green-100'
                        }`}
                      >
                        {seatUsagePercent > 90 ? '🔴 Critical' : seatUsagePercent > 80 ? '🟡 Warning' : '✅ Good'}
                      </span>
                    </div>
                  </div>
                </ChartCard>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
