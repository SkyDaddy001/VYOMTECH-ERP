'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth';
import { apiClient } from '@/lib/api-client';

interface MarketingCampaign {
  id: string;
  name: string;
  description?: string;
  channel: string;
  budget?: number;
  status: string;
  start_date?: string;
  end_date?: string;
  target_audience?: string;
  created_at: string;
  updated_at: string;
}

export default function MarketingPage() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const [campaigns, setCampaigns] = useState<MarketingCampaign[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showForm, setShowForm] = useState(false);
  const [newCampaign, setNewCampaign] = useState({
    name: '',
    description: '',
    channel: 'social',
    budget: '',
    status: 'active',
    start_date: '',
    end_date: '',
    target_audience: '',
  });

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login');
    } else {
      fetchCampaigns();
    }
  }, [isAuthenticated, router]);

  const fetchCampaigns = async () => {
    try {
      setLoading(true);
      const data = await apiClient.listMarketingCampaigns();
      setCampaigns(Array.isArray(data) ? data : []);
      setError('');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load campaigns');
      setCampaigns([]);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateCampaign = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const campaignData = {
        name: newCampaign.name,
        description: newCampaign.description || undefined,
        channel: newCampaign.channel,
        budget: newCampaign.budget ? parseFloat(newCampaign.budget) : undefined,
        status: newCampaign.status,
        start_date: newCampaign.start_date ? new Date(newCampaign.start_date) : undefined,
        end_date: newCampaign.end_date ? new Date(newCampaign.end_date) : undefined,
        target_audience: newCampaign.target_audience || undefined,
      };

      await apiClient.createMarketingCampaign(campaignData);
      setNewCampaign({
        name: '',
        description: '',
        channel: 'social',
        budget: '',
        status: 'active',
        start_date: '',
        end_date: '',
        target_audience: '',
      });
      setShowForm(false);
      await fetchCampaigns();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create campaign');
    }
  };

  const handleDeleteCampaign = async (id: string) => {
    if (!confirm('Are you sure you want to delete this campaign?')) return;

    try {
      await apiClient.deleteMarketingCampaign(id);
      await fetchCampaigns();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to delete campaign');
    }
  };

  const getChannelIcon = (channel: string) => {
    switch (channel) {
      case 'social':
        return '📱';
      case 'email':
        return '📧';
      case 'broadcast':
        return '📺';
      case 'print':
        return '📰';
      case 'digital':
        return '💻';
      default:
        return '📢';
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

  const formatCurrency = (amount?: number) => {
    if (!amount) return '-';
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'AED',
    }).format(amount);
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
            <h1 className="text-3xl font-bold text-gray-900">Marketing Campaigns</h1>
            <p className="text-gray-600 mt-2">Manage and track marketing campaigns to generate leads</p>
          </div>
          <button
            onClick={() => setShowForm(!showForm)}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
          >
            {showForm ? 'Cancel' : 'New Campaign'}
          </button>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        {/* Create Campaign Form */}
        {showForm && (
          <div className="mb-8 bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Create New Campaign</h2>
            <form onSubmit={handleCreateCampaign}>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Campaign Name *
                  </label>
                  <input
                    type="text"
                    required
                    value={newCampaign.name}
                    onChange={(e) => setNewCampaign({ ...newCampaign, name: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="e.g., Summer Property Campaign"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Channel *</label>
                  <select
                    required
                    value={newCampaign.channel}
                    onChange={(e) => setNewCampaign({ ...newCampaign, channel: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="social">Social Media</option>
                    <option value="email">Email</option>
                    <option value="broadcast">TV/Radio</option>
                    <option value="print">Print</option>
                    <option value="digital">Digital Ads</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Budget (AED)</label>
                  <input
                    type="number"
                    value={newCampaign.budget}
                    onChange={(e) => setNewCampaign({ ...newCampaign, budget: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="50000"
                    step="1000"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
                  <select
                    value={newCampaign.status}
                    onChange={(e) => setNewCampaign({ ...newCampaign, status: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="planning">Planning</option>
                    <option value="active">Active</option>
                    <option value="completed">Completed</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Start Date</label>
                  <input
                    type="date"
                    value={newCampaign.start_date}
                    onChange={(e) => setNewCampaign({ ...newCampaign, start_date: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">End Date</label>
                  <input
                    type="date"
                    value={newCampaign.end_date}
                    onChange={(e) => setNewCampaign({ ...newCampaign, end_date: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
              </div>

              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                <textarea
                  value={newCampaign.description}
                  onChange={(e) => setNewCampaign({ ...newCampaign, description: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Campaign description..."
                  rows={2}
                />
              </div>

              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-700 mb-1">Target Audience</label>
                <input
                  type="text"
                  value={newCampaign.target_audience}
                  onChange={(e) => setNewCampaign({ ...newCampaign, target_audience: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="e.g., Expats aged 25-45"
                />
              </div>

              <button
                type="submit"
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
              >
                Create Campaign
              </button>
            </form>
          </div>
        )}

        {/* Loading State */}
        {loading && (
          <div className="flex items-center justify-center py-12">
            <p className="text-gray-600">Loading campaigns...</p>
          </div>
        )}

        {/* Campaigns Grid */}
        {!loading && campaigns.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {campaigns.map((campaign) => (
              <div key={campaign.id} className="bg-white rounded-lg shadow hover:shadow-lg transition p-6">
                <div className="flex justify-between items-start mb-4">
                  <div className="flex items-center gap-3">
                    <span className="text-3xl">{getChannelIcon(campaign.channel)}</span>
                    <div>
                      <h3 className="text-lg font-semibold text-gray-900">{campaign.name}</h3>
                      <p className="text-xs text-gray-500 capitalize">{campaign.channel}</p>
                    </div>
                  </div>
                  <button
                    onClick={() => handleDeleteCampaign(campaign.id)}
                    className="text-red-600 hover:text-red-800 text-sm font-medium"
                  >
                    ✕
                  </button>
                </div>

                <span className={`inline-block mb-4 px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(campaign.status)}`}>
                  {campaign.status.charAt(0).toUpperCase() + campaign.status.slice(1)}
                </span>

                {campaign.description && (
                  <p className="text-gray-600 text-sm mb-4">{campaign.description}</p>
                )}

                <div className="bg-gray-50 rounded p-3 mb-4">
                  <p className="text-sm font-semibold text-gray-900">Budget: {formatCurrency(campaign.budget)}</p>
                </div>

                {campaign.target_audience && (
                  <div className="mb-4">
                    <p className="text-xs text-gray-500">Target Audience</p>
                    <p className="text-sm text-gray-700">{campaign.target_audience}</p>
                  </div>
                )}

                <div className="grid grid-cols-2 gap-2 text-sm mb-4">
                  <div>
                    <p className="text-gray-500 text-xs">Start</p>
                    <p className="font-medium text-gray-900">{formatDate(campaign.start_date)}</p>
                  </div>
                  <div>
                    <p className="text-gray-500 text-xs">End</p>
                    <p className="font-medium text-gray-900">{formatDate(campaign.end_date)}</p>
                  </div>
                </div>

                <div className="pt-3 border-t border-gray-200">
                  <p className="text-xs text-gray-500">
                    Updated: {formatDate(campaign.updated_at)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Empty State */}
        {!loading && campaigns.length === 0 && (
          <div className="bg-white rounded-lg shadow p-12 text-center">
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No campaigns yet</h3>
            <p className="text-gray-600 mb-6">Create your first marketing campaign to get started</p>
            <button
              onClick={() => setShowForm(true)}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
            >
              Create First Campaign
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
