'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/lib/hooks/useAuth';
import userService, { type UserActivity } from '@/services/user-service';
import Link from 'next/link';

export default function ActivityDashboard() {
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [activities, setActivities] = useState<UserActivity[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (authLoading || !isAuthenticated) return;
    loadData();
  }, [isAuthenticated, authLoading]);

  const loadData = async () => {
    try {
      setIsLoading(true);
      const data = await userService.getUserActivity({ days: 7, limit: 100 });
      setActivities(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load activities');
    } finally {
      setIsLoading(false);
    }
  };

  const getActivityStats = () => {
    const stats = {
      total: activities.length,
      uniqueUsers: new Set(activities.map((a) => a.userId)).size,
      activityTypes: new Set(activities.map((a) => a.activityType)).size,
      logins: activities.filter((a) => a.activityType === 'login').length,
      logouts: activities.filter((a) => a.activityType === 'logout').length,
      creates: activities.filter((a) => a.activityType === 'create').length,
      updates: activities.filter((a) => a.activityType === 'update').length,
      deletes: activities.filter((a) => a.activityType === 'delete').length,
    };
    return stats;
  };

  const getActivityTypeColor = (type: string) => {
    const colors: Record<string, string> = {
      login: 'bg-blue-100 text-blue-800',
      logout: 'bg-gray-100 text-gray-800',
      create: 'bg-green-100 text-green-800',
      update: 'bg-yellow-100 text-yellow-800',
      delete: 'bg-red-100 text-red-800',
      export: 'bg-purple-100 text-purple-800',
      import: 'bg-indigo-100 text-indigo-800',
    };
    return colors[type] || 'bg-gray-100 text-gray-800';
  };

  const getActivityIcon = (type: string) => {
    const icons: Record<string, string> = {
      login: '🔓',
      logout: '🔐',
      create: '➕',
      update: '✏️',
      delete: '🗑️',
      export: '📤',
      import: '📥',
    };
    return icons[type] || '📝';
  };

  if (authLoading || isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  const stats = getActivityStats();

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Activity Dashboard</h1>
            <p className="text-gray-600 mt-2">Last 7 days activity overview</p>
          </div>
          <Link
            href="/activity-logs"
            className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-6 rounded-lg transition"
          >
            View All Logs
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
                <p className="text-gray-600 text-sm font-medium">Total Activities</p>
                <p className="text-3xl font-bold text-blue-600 mt-2">{stats.total}</p>
              </div>
              <div className="text-4xl">📊</div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-green-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Unique Users</p>
                <p className="text-3xl font-bold text-green-600 mt-2">{stats.uniqueUsers}</p>
              </div>
              <div className="text-4xl">👥</div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-purple-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Activity Types</p>
                <p className="text-3xl font-bold text-purple-600 mt-2">{stats.activityTypes}</p>
              </div>
              <div className="text-4xl">🏷️</div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-orange-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Avg/Day</p>
                <p className="text-3xl font-bold text-orange-600 mt-2">{Math.round(stats.total / 7)}</p>
              </div>
              <div className="text-4xl">📈</div>
            </div>
          </div>
        </div>

        {/* Activity Breakdown */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Activity Types */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-6">Activity by Type</h3>
            <div className="space-y-4">
              {[
                { type: 'login', count: stats.logins, icon: '🔓' },
                { type: 'logout', count: stats.logouts, icon: '🔐' },
                { type: 'create', count: stats.creates, icon: '➕' },
                { type: 'update', count: stats.updates, icon: '✏️' },
                { type: 'delete', count: stats.deletes, icon: '🗑️' },
              ].map(({ type, count, icon }) => (
                <div key={type}>
                  <div className="flex justify-between items-center mb-2">
                    <span className="font-medium text-gray-700">{icon} {type.charAt(0).toUpperCase() + type.slice(1)}</span>
                    <span className="font-bold text-gray-900">{count}</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className={`h-2 rounded-full ${
                        type === 'create'
                          ? 'bg-green-500'
                          : type === 'update'
                          ? 'bg-yellow-500'
                          : type === 'delete'
                          ? 'bg-red-500'
                          : 'bg-blue-500'
                      } transition`}
                      style={{ width: `${stats.total ? (count / stats.total) * 100 : 0}%` }}
                    />
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Top Active Users */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-6">Top Active Users</h3>
            <div className="space-y-3">
              {Object.entries(
                activities.reduce(
                  (acc, activity) => {
                    acc[activity.email] = (acc[activity.email] || 0) + 1;
                    return acc;
                  },
                  {} as Record<string, number>
                )
              )
                .sort((a, b) => b[1] - a[1])
                .slice(0, 5)
                .map(([email, count]) => (
                  <div key={email} className="flex justify-between items-center p-3 bg-gray-50 rounded-lg">
                    <span className="text-gray-700">{email}</span>
                    <span className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">
                      {count} activities
                    </span>
                  </div>
                ))}
            </div>
          </div>
        </div>

        {/* Recent Activities Preview */}
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <div className="p-6 border-b">
            <h3 className="text-lg font-semibold text-gray-900">Recent Activities</h3>
          </div>

          {activities.length === 0 ? (
            <div className="p-8 text-center text-gray-600">No activities found</div>
          ) : (
            <div className="divide-y max-h-96 overflow-y-auto">
              {activities.slice(0, 10).map((activity, index) => (
                <div key={index} className="p-4 hover:bg-gray-50 transition">
                  <div className="flex items-start gap-4">
                    <div className="text-xl">{getActivityIcon(activity.activityType)}</div>
                    <div className="flex-1">
                      <div className="flex justify-between items-start">
                        <div>
                          <p className="font-medium text-gray-900">{activity.email}</p>
                          <span
                            className={`inline-block mt-1 px-2 py-1 rounded text-xs font-medium ${getActivityTypeColor(
                              activity.activityType
                            )}`}
                          >
                            {activity.activityType}
                          </span>
                        </div>
                        <time className="text-xs text-gray-500">
                          {new Date(activity.lastActivityAt).toLocaleString()}
                        </time>
                      </div>
                      {activity.ipAddress && (
                        <p className="text-xs text-gray-500 mt-2">IP: {activity.ipAddress}</p>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}

          <div className="p-4 bg-gray-50 border-t text-center">
            <Link href="/activity-logs" className="text-blue-600 hover:text-blue-800 font-medium">
              View All Activities →
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
