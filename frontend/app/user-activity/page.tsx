'use client';

import { useEffect, useState } from 'react';

interface UserActivity {
  userId: string;
  userName: string;
  email: string;
  lastLogin: string;
  lastLogout: string;
  currentStatus: 'online' | 'offline' | 'idle';
  sessionDuration: number;
  device: string;
  ipAddress: string;
}

export default function UserActivityLog() {
  const [activities, setActivities] = useState<UserActivity[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<'all' | 'online' | 'offline'>('all');

  useEffect(() => {
    loadActivityLog();
    const interval = setInterval(loadActivityLog, 10000); // Refresh every 10 seconds
    return () => clearInterval(interval);
  }, []);

  const loadActivityLog = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch('/api/v1/tenant/users/activity', {
        headers: { Authorization: `Bearer ${localStorage.getItem('authToken')}` },
      });

      if (!response.ok) throw new Error('Failed to load activity log');

      const data = await response.json();
      setActivities(data.activities || []);
    } catch (err) {
      setError('Failed to load activity log');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const formatDuration = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
  };

  const filteredActivities =
    filter === 'all'
      ? activities
      : activities.filter((a) => (filter === 'online' ? a.currentStatus === 'online' : a.currentStatus === 'offline'));

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">User Activity Log</h1>
          <p className="text-gray-600 mt-2">Real-time user login/logout activity</p>
        </div>

        {/* Filters */}
        <div className="bg-white rounded-lg shadow p-4 mb-6">
          <div className="flex gap-4">
            <button
              onClick={() => setFilter('all')}
              className={`px-4 py-2 rounded-lg font-medium transition ${
                filter === 'all'
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              All Users ({activities.length})
            </button>
            <button
              onClick={() => setFilter('online')}
              className={`px-4 py-2 rounded-lg font-medium transition ${
                filter === 'online'
                  ? 'bg-green-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              Online ({activities.filter((a) => a.currentStatus === 'online').length})
            </button>
            <button
              onClick={() => setFilter('offline')}
              className={`px-4 py-2 rounded-lg font-medium transition ${
                filter === 'offline'
                  ? 'bg-gray-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              Offline ({activities.filter((a) => a.currentStatus === 'offline').length})
            </button>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        {/* Activity Table */}
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    User
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Last Login
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Last Logout
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Duration
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Device
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    IP Address
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {filteredActivities.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="px-6 py-4 text-center text-gray-500">
                      No user activity found
                    </td>
                  </tr>
                ) : (
                  filteredActivities.map((activity, idx) => (
                    <tr key={idx} className="hover:bg-gray-50 transition">
                      <td className="px-6 py-4">
                        <div>
                          <p className="text-sm font-medium text-gray-900">{activity.userName}</p>
                          <p className="text-xs text-gray-500">{activity.email}</p>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-2">
                          <span
                            className={`w-2 h-2 rounded-full ${
                              activity.currentStatus === 'online'
                                ? 'bg-green-500'
                                : 'bg-gray-400'
                            }`}
                          />
                          <span className="text-sm font-medium text-gray-700 capitalize">
                            {activity.currentStatus}
                          </span>
                        </div>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-600">
                        <div>
                          <p>{formatDate(activity.lastLogin)}</p>
                          <p className="text-xs text-gray-500">{formatTime(activity.lastLogin)}</p>
                        </div>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-600">
                        {activity.lastLogout ? (
                          <div>
                            <p>{formatDate(activity.lastLogout)}</p>
                            <p className="text-xs text-gray-500">{formatTime(activity.lastLogout)}</p>
                          </div>
                        ) : (
                          <span className="text-gray-500">-</span>
                        )}
                      </td>
                      <td className="px-6 py-4 text-sm font-medium text-gray-900">
                        {formatDuration(activity.sessionDuration)}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-600">{activity.device}</td>
                      <td className="px-6 py-4 text-sm text-gray-600 font-mono text-xs">
                        {activity.ipAddress}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        {/* Footer Info */}
        <div className="mt-6 text-sm text-gray-600">
          <p>Last updated: {new Date().toLocaleTimeString()}</p>
        </div>
      </div>
    </div>
  );
}
