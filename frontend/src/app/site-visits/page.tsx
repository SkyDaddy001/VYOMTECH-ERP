'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth';
import { apiClient } from '@/lib/api-client';

interface SiteVisit {
  id: string;
  lead_id: string;
  property_name: string;
  visit_date: string;
  visit_time?: string;
  duration?: number;
  visited_by?: string;
  notes?: string;
  rating?: number;
  status: string;
  followup_required: boolean;
  followup_date?: string;
  created_at: string;
  updated_at: string;
}

export default function SiteVisitsPage() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const [visits, setVisits] = useState<SiteVisit[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showForm, setShowForm] = useState(false);
  const [newVisit, setNewVisit] = useState({
    lead_id: '',
    property_name: '',
    visit_date: '',
    visit_time: '',
    duration: '',
    visited_by: '',
    notes: '',
    rating: '',
    status: 'scheduled',
    followup_required: false,
    followup_date: '',
  });

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login');
    } else {
      fetchVisits();
    }
  }, [isAuthenticated, router]);

  const fetchVisits = async () => {
    try {
      setLoading(true);
      const data = await apiClient.listSiteVisits();
      setVisits(Array.isArray(data) ? data : []);
      setError('');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load site visits');
      setVisits([]);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateVisit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const visitData = {
        lead_id: newVisit.lead_id,
        property_name: newVisit.property_name,
        visit_date: new Date(newVisit.visit_date),
        visit_time: newVisit.visit_time || undefined,
        duration: newVisit.duration ? parseInt(newVisit.duration) : undefined,
        visited_by: newVisit.visited_by || undefined,
        notes: newVisit.notes || undefined,
        rating: newVisit.rating ? parseInt(newVisit.rating) : undefined,
        status: newVisit.status,
        followup_required: newVisit.followup_required,
        followup_date: newVisit.followup_date ? new Date(newVisit.followup_date) : undefined,
      };

      await apiClient.createSiteVisit(visitData);
      setNewVisit({
        lead_id: '',
        property_name: '',
        visit_date: '',
        visit_time: '',
        duration: '',
        visited_by: '',
        notes: '',
        rating: '',
        status: 'scheduled',
        followup_required: false,
        followup_date: '',
      });
      setShowForm(false);
      await fetchVisits();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create visit');
    }
  };

  const handleDeleteVisit = async (id: string) => {
    if (!confirm('Are you sure you want to delete this visit?')) return;

    try {
      await apiClient.deleteSiteVisit(id);
      await fetchVisits();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to delete visit');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'scheduled':
        return 'bg-blue-100 text-blue-800';
      case 'cancelled':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getRatingStars = (rating?: number) => {
    if (!rating) return '-';
    return '⭐'.repeat(rating) + ` (${rating}/5)`;
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
            <h1 className="text-3xl font-bold text-gray-900">Site Visits</h1>
            <p className="text-gray-600 mt-2">Schedule and track property site visits</p>
          </div>
          <button
            onClick={() => setShowForm(!showForm)}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
          >
            {showForm ? 'Cancel' : 'Schedule Visit'}
          </button>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        {/* Create Visit Form */}
        {showForm && (
          <div className="mb-8 bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Schedule New Site Visit</h2>
            <form onSubmit={handleCreateVisit}>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Lead ID *
                  </label>
                  <input
                    type="text"
                    required
                    value={newVisit.lead_id}
                    onChange={(e) => setNewVisit({ ...newVisit, lead_id: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="lead-123"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Property Name *
                  </label>
                  <input
                    type="text"
                    required
                    value={newVisit.property_name}
                    onChange={(e) => setNewVisit({ ...newVisit, property_name: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Luxury Apartment Downtown"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Visit Date *
                  </label>
                  <input
                    type="date"
                    required
                    value={newVisit.visit_date}
                    onChange={(e) => setNewVisit({ ...newVisit, visit_date: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Visit Time</label>
                  <input
                    type="time"
                    value={newVisit.visit_time}
                    onChange={(e) => setNewVisit({ ...newVisit, visit_time: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Duration (minutes)</label>
                  <input
                    type="number"
                    value={newVisit.duration}
                    onChange={(e) => setNewVisit({ ...newVisit, duration: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="45"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Visited By</label>
                  <input
                    type="text"
                    value={newVisit.visited_by}
                    onChange={(e) => setNewVisit({ ...newVisit, visited_by: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Sales Agent Name"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
                  <select
                    value={newVisit.status}
                    onChange={(e) => setNewVisit({ ...newVisit, status: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="scheduled">Scheduled</option>
                    <option value="completed">Completed</option>
                    <option value="cancelled">Cancelled</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Rating (1-5)</label>
                  <input
                    type="number"
                    min="1"
                    max="5"
                    value={newVisit.rating}
                    onChange={(e) => setNewVisit({ ...newVisit, rating: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="5"
                  />
                </div>
              </div>

              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
                <textarea
                  value={newVisit.notes}
                  onChange={(e) => setNewVisit({ ...newVisit, notes: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Visit observations and feedback..."
                  rows={3}
                />
              </div>

              <div className="mb-4 flex items-center">
                <input
                  type="checkbox"
                  id="followup"
                  checked={newVisit.followup_required}
                  onChange={(e) => setNewVisit({ ...newVisit, followup_required: e.target.checked })}
                  className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                />
                <label htmlFor="followup" className="ml-2 block text-sm text-gray-700">
                  Followup Required
                </label>
              </div>

              {newVisit.followup_required && (
                <div className="mb-6">
                  <label className="block text-sm font-medium text-gray-700 mb-1">Followup Date</label>
                  <input
                    type="date"
                    value={newVisit.followup_date}
                    onChange={(e) => setNewVisit({ ...newVisit, followup_date: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
              )}

              <button
                type="submit"
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
              >
                Schedule Visit
              </button>
            </form>
          </div>
        )}

        {/* Loading State */}
        {loading && (
          <div className="flex items-center justify-center py-12">
            <p className="text-gray-600">Loading site visits...</p>
          </div>
        )}

        {/* Visits Grid */}
        {!loading && visits.length > 0 && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {visits.map((visit) => (
              <div key={visit.id} className="bg-white rounded-lg shadow hover:shadow-lg transition p-6">
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900">{visit.property_name}</h3>
                    <p className="text-sm text-gray-500">Lead: {visit.lead_id}</p>
                  </div>
                  <button
                    onClick={() => handleDeleteVisit(visit.id)}
                    className="text-red-600 hover:text-red-800 text-sm font-medium"
                  >
                    ✕
                  </button>
                </div>

                <span className={`inline-block mb-4 px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(visit.status)}`}>
                  {visit.status.charAt(0).toUpperCase() + visit.status.slice(1)}
                </span>

                <div className="grid grid-cols-2 gap-4 mb-4 text-sm">
                  <div>
                    <p className="text-gray-500 text-xs">Visit Date</p>
                    <p className="font-medium text-gray-900">{formatDate(visit.visit_date)}</p>
                  </div>
                  <div>
                    <p className="text-gray-500 text-xs">Visit Time</p>
                    <p className="font-medium text-gray-900">{visit.visit_time || '-'}</p>
                  </div>
                  <div>
                    <p className="text-gray-500 text-xs">Duration</p>
                    <p className="font-medium text-gray-900">{visit.duration ? `${visit.duration} min` : '-'}</p>
                  </div>
                  <div>
                    <p className="text-gray-500 text-xs">Rating</p>
                    <p className="font-medium text-gray-900">{getRatingStars(visit.rating)}</p>
                  </div>
                </div>

                {visit.visited_by && (
                  <div className="mb-3 p-3 bg-blue-50 rounded">
                    <p className="text-xs text-gray-500">Visited By</p>
                    <p className="text-sm font-medium text-gray-900">{visit.visited_by}</p>
                  </div>
                )}

                {visit.notes && (
                  <div className="mb-3 p-3 bg-gray-50 rounded">
                    <p className="text-xs text-gray-500">Notes</p>
                    <p className="text-sm text-gray-700">{visit.notes}</p>
                  </div>
                )}

                {visit.followup_required && (
                  <div className="mb-3 p-3 bg-yellow-50 rounded border border-yellow-200">
                    <p className="text-xs text-gray-500">Followup Date</p>
                    <p className="text-sm font-medium text-yellow-900">{formatDate(visit.followup_date)}</p>
                  </div>
                )}

                <div className="pt-3 border-t border-gray-200">
                  <p className="text-xs text-gray-500">
                    Created: {formatDate(visit.created_at)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Empty State */}
        {!loading && visits.length === 0 && (
          <div className="bg-white rounded-lg shadow p-12 text-center">
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No site visits scheduled</h3>
            <p className="text-gray-600 mb-6">Schedule your first site visit to get started</p>
            <button
              onClick={() => setShowForm(true)}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
            >
              Schedule First Visit
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
