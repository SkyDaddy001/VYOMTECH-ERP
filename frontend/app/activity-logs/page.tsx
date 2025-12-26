'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/lib/hooks/useAuth';
import userService, { type UserActivity } from '@/services/user-service';

export default function ActivityLogsPage() {
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [activities, setActivities] = useState<UserActivity[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filterDays, setFilterDays] = useState('7');

  useEffect(() => {
    if (authLoading) return;
    if (!isAuthenticated) return;
    
    loadActivity();
  }, [isAuthenticated, authLoading, filterDays]);

  const loadActivity = async () => {
    try {
      setIsLoading(true);
      const data = await userService.getUserActivity({
        days: parseInt(filterDays),
        limit: 500,
      });
      setActivities(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load activity');
    } finally {
      setIsLoading(false);
    }
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  const getActivityIcon = (activityType: string) => {
    const icons: Record<string, string> = {
      login: '🔓',
      logout: '🔐',
      create: '➕',
      update: '✏️',
      delete: '🗑️',
      export: '📤',
      import: '📥',
      default: '📝',
    };
    return icons[activityType] || icons.default;
  };

  if (authLoading) {
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
          <h1 className="text-3xl font-bold text-gray-900">Activity Logs</h1>
          <p className="text-gray-600 mt-2">User activity and system events</p>
        </div>

        {/* Error Alert */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
            {error}
          </div>
        )}

        {/* Filters */}
        <div className="bg-white rounded-lg shadow p-4 mb-6">
          <div className="flex gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Time Range</label>
              <select
                value={filterDays}
                onChange={(e) => setFilterDays(e.target.value)}
                className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="1">Last 24 hours</option>
                <option value="7">Last 7 days</option>
                <option value="30">Last 30 days</option>
                <option value="90">Last 90 days</option>
              </select>
            </div>
          </div>
        </div>

        {/* Activity Timeline */}
        <div className="bg-white rounded-lg shadow overflow-hidden">
          {isLoading ? (
            <div className="p-8 text-center text-gray-600">Loading activities...</div>
          ) : activities.length === 0 ? (
            <div className="p-8 text-center text-gray-600">No activities found</div>
          ) : (
            <div className="divide-y">
              {activities.map((activity, index) => (
                <div key={index} className="p-6 hover:bg-gray-50 transition">
                  <div className="flex items-start gap-4">
                    <div className="text-2xl">{getActivityIcon(activity.activityType)}</div>
                    <div className="flex-1">
                      <div className="flex justify-between items-start mb-2">
                        <div>
                          <p className="font-medium text-gray-900">{activity.email}</p>
                          <p className="text-sm text-gray-600 mt-1">
                            {activity.activityType.charAt(0).toUpperCase() + activity.activityType.slice(1)}
                          </p>
                        </div>
                        <time className="text-sm text-gray-600 whitespace-nowrap ml-4">
                          {formatTime(activity.lastActivityAt)}
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
        </div>

        {/* Stats */}
        {!isLoading && activities.length > 0 && (
          <div className="mt-6 grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="bg-white rounded-lg shadow p-6">
              <div className="text-3xl font-bold text-blue-600">{activities.length}</div>
              <div className="text-gray-600 text-sm mt-2">Total Activities</div>
            </div>
            <div className="bg-white rounded-lg shadow p-6">
              <div className="text-3xl font-bold text-green-600">
                {new Set(activities.map((a) => a.userId)).size}
              </div>
              <div className="text-gray-600 text-sm mt-2">Active Users</div>
            </div>
            <div className="bg-white rounded-lg shadow p-6">
              <div className="text-3xl font-bold text-purple-600">
                {new Set(activities.map((a) => a.activityType)).size}
              </div>
              <div className="text-gray-600 text-sm mt-2">Activity Types</div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
