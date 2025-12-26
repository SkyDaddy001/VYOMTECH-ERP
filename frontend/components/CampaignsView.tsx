/**
 * Campaign Management UI Component
 * Uses Prisma Campaign model types
 */

'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api-client';
import type { CampaignResponse, MetricsResponse } from '@/lib/types';

interface CampaignsViewProps {
  platform?: 'google' | 'meta';
}

export default function CampaignsView({ platform = 'google' }: CampaignsViewProps) {
  const [campaigns, setCampaigns] = useState<CampaignResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedCampaign, setSelectedCampaign] = useState<CampaignResponse | null>(null);

  useEffect(() => {
    loadCampaigns();
  }, [platform]);

  const loadCampaigns = async () => {
    try {
      setLoading(true);
      const data = await apiClient.listCampaigns();
      if (Array.isArray(data)) {
        setCampaigns(data as CampaignResponse[]);
      }
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Failed to load campaigns'
      );
    } finally {
      setLoading(false);
    }
  };

  const handlePauseCampaign = async (campaignId: string) => {
    try {
      await apiClient.pauseCampaign(campaignId, platform);
      await loadCampaigns(); // Reload
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Failed to pause campaign'
      );
    }
  };

  const handleResumeCampaign = async (campaignId: string) => {
    try {
      await apiClient.resumeCampaign(campaignId, platform);
      await loadCampaigns(); // Reload
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Failed to resume campaign'
      );
    }
  };

  const handleUpdateBudget = async (campaignId: string, newBudget: number) => {
    try {
      await apiClient.updateBudget(campaignId, platform, newBudget);
      await loadCampaigns(); // Reload
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Failed to update budget'
      );
    }
  };

  if (loading) {
    return <div className="p-4 text-center">Loading campaigns...</div>;
  }

  if (error) {
    return <div className="p-4 bg-red-50 text-red-800 rounded">{error}</div>;
  }

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-bold text-gray-900 capitalize">
        {platform === 'google' ? 'Google Ads' : 'Meta Ads'} Campaigns
      </h2>

      {campaigns.length === 0 ? (
        <div className="p-4 bg-gray-50 rounded text-gray-600 text-center">
          No campaigns found
        </div>
      ) : (
        <div className="grid gap-4">
          {campaigns.map((campaign) => (
            <div
              key={campaign.id}
              className="p-4 bg-white border border-gray-200 rounded-lg hover:shadow-md transition"
            >
              <div className="flex items-center justify-between mb-3">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900">
                    {campaign.name}
                  </h3>
                  <p className="text-sm text-gray-600">
                    ID: {campaign.id}
                  </p>
                </div>
                <div className="text-right">
                  <span
                    className={`inline-block px-3 py-1 rounded text-sm font-medium ${
                      campaign.status === 'active'
                        ? 'bg-green-100 text-green-800'
                        : campaign.status === 'paused'
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {campaign.status.charAt(0).toUpperCase() +
                      campaign.status.slice(1)}
                  </span>
                </div>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                <div className="bg-gray-50 p-3 rounded">
                  <p className="text-xs text-gray-600">Daily Budget</p>
                  <p className="text-lg font-bold text-gray-900">
                    ${campaign.budget.toFixed(2)}
                  </p>
                </div>
                <div className="bg-gray-50 p-3 rounded">
                  <p className="text-xs text-gray-600">Impressions</p>
                  <p className="text-lg font-bold text-gray-900">
                    {campaign.impressions?.toLocaleString() || '0'}
                  </p>
                </div>
                <div className="bg-gray-50 p-3 rounded">
                  <p className="text-xs text-gray-600">Clicks</p>
                  <p className="text-lg font-bold text-gray-900">
                    {campaign.clicks?.toLocaleString() || '0'}
                  </p>
                </div>
                <div className="bg-gray-50 p-3 rounded">
                  <p className="text-xs text-gray-600">Conversions</p>
                  <p className="text-lg font-bold text-gray-900">
                    {campaign.conversions?.toLocaleString() || '0'}
                  </p>
                </div>
              </div>

              <div className="flex gap-2 flex-wrap">
                {campaign.status === 'active' ? (
                  <button
                    onClick={() => handlePauseCampaign(campaign.id)}
                    className="px-4 py-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded font-medium transition"
                  >
                    Pause
                  </button>
                ) : (
                  <button
                    onClick={() => handleResumeCampaign(campaign.id)}
                    className="px-4 py-2 bg-green-500 hover:bg-green-600 text-white rounded font-medium transition"
                  >
                    Resume
                  </button>
                )}
                <button
                  onClick={() => setSelectedCampaign(campaign)}
                  className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded font-medium transition"
                >
                  View Details
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {selectedCampaign && (
        <CampaignDetailsModal
          campaign={selectedCampaign}
          platform={platform}
          onClose={() => setSelectedCampaign(null)}
          onBudgetUpdate={(newBudget) => {
            handleUpdateBudget(selectedCampaign.id, newBudget);
            setSelectedCampaign(null);
          }}
        />
      )}
    </div>
  );
}

interface CampaignDetailsModalProps {
  campaign: CampaignResponse;
  platform: 'google' | 'meta';
  onClose: () => void;
  onBudgetUpdate: (newBudget: number) => void;
}

function CampaignDetailsModal({
  campaign,
  platform,
  onClose,
  onBudgetUpdate,
}: CampaignDetailsModalProps) {
  const [newBudget, setNewBudget] = useState(campaign.budget.toString());
  const [metrics, setMetrics] = useState<MetricsResponse | null>(null);
  const [loadingMetrics, setLoadingMetrics] = useState(false);

  useEffect(() => {
    loadMetrics();
  }, [campaign.id]);

  const loadMetrics = async () => {
    try {
      setLoadingMetrics(true);
      const response = await apiClient.getCampaignMetrics(campaign.id, platform);
      setMetrics(response as MetricsResponse);
    } catch (err) {
      console.error('Failed to load metrics:', err);
    } finally {
      setLoadingMetrics(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg max-w-2xl w-full max-h-96 overflow-y-auto p-6">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-xl font-bold text-gray-900">
            Campaign Details: {campaign.name}
          </h3>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 text-2xl"
          >
            ×
          </button>
        </div>

        <div className="grid grid-cols-2 gap-4 mb-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Campaign ID
            </label>
            <p className="text-gray-900 text-sm p-2 bg-gray-50 rounded">
              {campaign.id}
            </p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Platform
            </label>
            <p className="text-gray-900 text-sm p-2 bg-gray-50 rounded capitalize">
              {platform}
            </p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Daily Budget
            </label>
            <input
              type="number"
              value={newBudget}
              onChange={(e) => setNewBudget(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Status
            </label>
            <p className="text-gray-900 text-sm p-2 bg-gray-50 rounded capitalize">
              {campaign.status}
            </p>
          </div>
        </div>

        {loadingMetrics ? (
          <p className="text-gray-600">Loading metrics...</p>
        ) : metrics ? (
          <div className="grid grid-cols-2 md:grid-cols-3 gap-3 mb-6">
            <div className="bg-gray-50 p-3 rounded">
              <p className="text-xs text-gray-600">ROI</p>
              <p className="text-lg font-bold text-green-600">
                {metrics.roi.toFixed(1)}%
              </p>
            </div>
            <div className="bg-gray-50 p-3 rounded">
              <p className="text-xs text-gray-600">CTR</p>
              <p className="text-lg font-bold text-gray-900">
                {metrics.ctr.toFixed(2)}%
              </p>
            </div>
            <div className="bg-gray-50 p-3 rounded">
              <p className="text-xs text-gray-600">CPC</p>
              <p className="text-lg font-bold text-gray-900">
                ${metrics.cpc.toFixed(2)}
              </p>
            </div>
          </div>
        ) : null}

        <div className="flex gap-2">
          <button
            onClick={() => onBudgetUpdate(parseFloat(newBudget))}
            className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition"
          >
            Update Budget
          </button>
          <button
            onClick={onClose}
            className="flex-1 px-4 py-2 bg-gray-300 hover:bg-gray-400 text-gray-800 rounded-lg font-medium transition"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  );
}
