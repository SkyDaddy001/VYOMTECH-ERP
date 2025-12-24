/**
 * Dashboard - Main application page
 * Shows overview of all campaigns and ROI metrics
 */

'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { apiClient } from '@/lib/api-client';
import { getStoredToken, revokeAuthToken } from '@/lib/auth-storage';
import { CampaignsView } from '@/components/CampaignsView';

export default function DashboardPage() {
  const [activeTab, setActiveTab] = useState<'google' | 'meta'>('google');
  const [portfolioROI, setPortfolioROI] = useState<any>(null);
  const [platformROI, setPlatformROI] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      setLoading(true);
      const [roiData, platformData] = await Promise.all([
        apiClient.getPortfolioROI(),
        apiClient.getPlatformROI(),
      ]);
      setPortfolioROI(roiData);
      setPlatformROI(platformData);
    } catch (err) {
      console.error('Failed to load dashboard data:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      const token = await getStoredToken();
      if (token) {
        // Try to revoke on backend
        try {
          await apiClient.revokeOAuthToken('google');
          await apiClient.revokeOAuthToken('meta');
        } catch {
          // Ignore errors during revocation
        }
      }
      await revokeAuthToken(token || '');
      window.location.href = '/auth/login';
    } catch (err) {
      console.error('Logout error:', err);
      window.location.href = '/auth/login';
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">VYOM ERP</h1>
            <p className="text-sm text-gray-600">Campaign Management & Analytics</p>
          </div>
          <button
            onClick={handleLogout}
            className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg font-medium transition"
          >
            Logout
          </button>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-8">
        {/* ROI Overview */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          {loading ? (
            <>
              <div className="h-32 bg-white rounded-lg animate-pulse"></div>
              <div className="h-32 bg-white rounded-lg animate-pulse"></div>
              <div className="h-32 bg-white rounded-lg animate-pulse"></div>
            </>
          ) : (
            <>
              <div className="bg-white p-6 rounded-lg shadow">
                <p className="text-gray-600 text-sm mb-2">Portfolio ROI</p>
                <p className="text-3xl font-bold text-green-600">
                  {portfolioROI?.roi?.toFixed(1) || '0'}%
                </p>
                <p className="text-xs text-gray-500 mt-2">
                  Revenue: ${portfolioROI?.revenue?.toLocaleString() || '0'}
                </p>
              </div>

              <div className="bg-white p-6 rounded-lg shadow">
                <p className="text-gray-600 text-sm mb-2">Total Spend</p>
                <p className="text-3xl font-bold text-blue-600">
                  ${portfolioROI?.spend?.toLocaleString() || '0'}
                </p>
                <p className="text-xs text-gray-500 mt-2">
                  Across all campaigns
                </p>
              </div>

              <div className="bg-white p-6 rounded-lg shadow">
                <p className="text-gray-600 text-sm mb-2">ROAS</p>
                <p className="text-3xl font-bold text-purple-600">
                  {portfolioROI?.roas?.toFixed(2) || '0'}x
                </p>
                <p className="text-xs text-gray-500 mt-2">
                  Return on Ad Spend
                </p>
              </div>
            </>
          )}
        </div>

        {/* Platform Comparison */}
        {platformROI && (
          <div className="bg-white p-6 rounded-lg shadow mb-8">
            <h2 className="text-xl font-bold text-gray-900 mb-4">
              Platform Comparison
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {platformROI.google && (
                <div className="border border-gray-200 p-4 rounded">
                  <h3 className="font-semibold text-gray-900 mb-3">
                    Google Ads
                  </h3>
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-gray-600">ROI:</span>
                      <span className="font-bold text-green-600">
                        {platformROI.google.roi?.toFixed(1)}%
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Spend:</span>
                      <span className="font-bold text-gray-900">
                        ${platformROI.google.spend?.toLocaleString() || '0'}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Revenue:</span>
                      <span className="font-bold text-gray-900">
                        ${platformROI.google.revenue?.toLocaleString() || '0'}
                      </span>
                    </div>
                  </div>
                </div>
              )}

              {platformROI.meta && (
                <div className="border border-gray-200 p-4 rounded">
                  <h3 className="font-semibold text-gray-900 mb-3">
                    Meta Ads
                  </h3>
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-gray-600">ROI:</span>
                      <span className="font-bold text-green-600">
                        {platformROI.meta.roi?.toFixed(1)}%
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Spend:</span>
                      <span className="font-bold text-gray-900">
                        ${platformROI.meta.spend?.toLocaleString() || '0'}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Revenue:</span>
                      <span className="font-bold text-gray-900">
                        ${platformROI.meta.revenue?.toLocaleString() || '0'}
                      </span>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Campaign Tabs */}
        <div className="bg-white rounded-lg shadow">
          <div className="flex border-b border-gray-200">
            <button
              onClick={() => setActiveTab('google')}
              className={`flex-1 px-6 py-4 font-medium transition ${
                activeTab === 'google'
                  ? 'text-blue-600 border-b-2 border-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Google Ads
            </button>
            <button
              onClick={() => setActiveTab('meta')}
              className={`flex-1 px-6 py-4 font-medium transition ${
                activeTab === 'meta'
                  ? 'text-blue-600 border-b-2 border-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Meta Ads
            </button>
          </div>

          <div className="p-6">
            <CampaignsView platform={activeTab} />
          </div>
        </div>
      </main>
    </div>
  );
}
