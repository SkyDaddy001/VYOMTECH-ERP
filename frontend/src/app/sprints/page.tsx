'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth';
import { apiClient } from '@/lib/api-client';

interface Sprint {
  id: string;
  name: string;
  description?: string;
  status: string;
  start_date?: string;
  end_date?: string;
  goal?: string;
  created_at: string;
  updated_at: string;
}

export default function SprintsPage() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const [sprints, setSprints] = useState<Sprint[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showForm, setShowForm] = useState(false);
  const [newSprint, setNewSprint] = useState({
    name: '',
    description: '',
    goal: '',
    status: 'planning',
    start_date: '',
    end_date: '',
  });

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login');
    } else {
      fetchSprints();
    }
  }, [isAuthenticated, router]);

  const fetchSprints = async () => {
    try {
      setLoading(true);
      const data = await apiClient.listSprints();
      setSprints(Array.isArray(data) ? data : []);
      setError('');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load sprints');
      setSprints([]);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateSprint = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const sprintData = {
        name: newSprint.name,
        description: newSprint.description || undefined,
        goal: newSprint.goal || undefined,
        status: newSprint.status,
        start_date: newSprint.start_date ? new Date(newSprint.start_date) : undefined,
        end_date: newSprint.end_date ? new Date(newSprint.end_date) : undefined,
      };

      await apiClient.createSprint(sprintData);
      setNewSprint({
        name: '',
        description: '',
        goal: '',
        status: 'planning',
        start_date: '',
        end_date: '',
      });
      setShowForm(false);
      await fetchSprints();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create sprint');
    }
  };

  const handleDeleteSprint = async (id: string) => {
    if (!confirm('Are you sure you want to delete this sprint?')) return;

    try {
      await apiClient.deleteSprint(id);
      await fetchSprints();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to delete sprint');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800';
      case 'completed':
        return 'bg-blue-100 text-blue-800';
      case 'planning':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  if (!isAuthenticated) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Sprints</h1>
            <p className="text-gray-600 mt-2">Manage and track project sprints</p>
          </div>
          <button
            onClick={() => setShowForm(!showForm)}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
          >
            {showForm ? 'Cancel' : 'New Sprint'}
          </button>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        {/* Create Sprint Form */}
        {showForm && (
          <div className="mb-8 bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Create New Sprint</h2>
            <form onSubmit={handleCreateSprint}>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Sprint Name *
                  </label>
                  <input
                    type="text"
                    required
                    value={newSprint.name}
                    onChange={(e) => setNewSprint({ ...newSprint, name: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="e.g., Sprint 1 - Core Features"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
                  <select
                    value={newSprint.status}
                    onChange={(e) => setNewSprint({ ...newSprint, status: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="planning">Planning</option>
                    <option value="active">Active</option>
                    <option value="completed">Completed</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Start Date
                  </label>
                  <input
                    type="date"
                    value={newSprint.start_date}
                    onChange={(e) => setNewSprint({ ...newSprint, start_date: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    End Date
                  </label>
                  <input
                    type="date"
                    value={newSprint.end_date}
                    onChange={(e) => setNewSprint({ ...newSprint, end_date: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
              </div>

              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <textarea
                  value={newSprint.description}
                  onChange={(e) => setNewSprint({ ...newSprint, description: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Sprint description..."
                  rows={3}
                />
              </div>

              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-700 mb-1">Goal</label>
                <textarea
                  value={newSprint.goal}
                  onChange={(e) => setNewSprint({ ...newSprint, goal: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="What does this sprint aim to achieve?"
                  rows={2}
                />
              </div>

              <button
                type="submit"
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
              >
                Create Sprint
              </button>
            </form>
          </div>
        )}

        {/* Loading State */}
        {loading && (
          <div className="flex items-center justify-center py-12">
            <p className="text-gray-600">Loading sprints...</p>
          </div>
        )}

        {/* Sprints List */}
        {!loading && sprints.length > 0 && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {sprints.map((sprint) => (
              <div key={sprint.id} className="bg-white rounded-lg shadow hover:shadow-lg transition p-6">
                <div className="flex justify-between items-start mb-4">
                  <div className="flex-1">
                    <h3 className="text-xl font-semibold text-gray-900">{sprint.name}</h3>
                    <span className={`inline-block mt-2 px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(sprint.status)}`}>
                      {sprint.status.charAt(0).toUpperCase() + sprint.status.slice(1)}
                    </span>
                  </div>
                  <button
                    onClick={() => handleDeleteSprint(sprint.id)}
                    className="text-red-600 hover:text-red-800 text-sm font-medium"
                  >
                    Delete
                  </button>
                </div>

                {sprint.description && (
                  <p className="text-gray-600 text-sm mb-4">{sprint.description}</p>
                )}

                {sprint.goal && (
                  <div className="mb-4 p-3 bg-blue-50 rounded">
                    <p className="text-sm text-gray-700">
                      <strong>Goal:</strong> {sprint.goal}
                    </p>
                  </div>
                )}

                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <p className="text-gray-500">Start Date</p>
                    <p className="font-medium text-gray-900">{formatDate(sprint.start_date)}</p>
                  </div>
                  <div>
                    <p className="text-gray-500">End Date</p>
                    <p className="font-medium text-gray-900">{formatDate(sprint.end_date)}</p>
                  </div>
                </div>

                <div className="mt-4 pt-4 border-t border-gray-200">
                  <p className="text-xs text-gray-500">
                    Created: {formatDate(sprint.created_at)} | Updated: {formatDate(sprint.updated_at)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Empty State */}
        {!loading && sprints.length === 0 && (
          <div className="bg-white rounded-lg shadow p-12 text-center">
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No sprints yet</h3>
            <p className="text-gray-600 mb-6">Create your first sprint to get started</p>
            <button
              onClick={() => setShowForm(true)}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
            >
              Create First Sprint
            </button>
          </div>
        )}

        {/* Back Button */}
        <div className="mt-8">
          <a href="/" className="text-blue-600 hover:text-blue-700 font-medium">
            ← Back to Dashboard
          </a>
        </div>
      </div>
    </div>
  );
}
