'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/lib/hooks/useAuth';
import tenantService, { type Tenant } from '@/services/tenant-service';
import Link from 'next/link';

export default function SettingsDashboard() {
  const { user, isLoading: authLoading } = useAuth();
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [credentials, setCredentials] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (authLoading || !user) return;
    loadData();
  }, [authLoading, user]);

  const loadData = async () => {
    try {
      setIsLoading(true);
      const [tenantData, credentialsData] = await Promise.all([
        tenantService.getCurrentTenant(),
        tenantService.getTenantCredentials(user?.tenantId || ''),
      ]);
      setTenant(tenantData);
      setCredentials(credentialsData);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load data');
    } finally {
      setIsLoading(false);
    }
  };

  const credentialTypes = [
    { key: 'google', name: 'Google OAuth', icon: '🔵' },
    { key: 'meta', name: 'Meta OAuth', icon: '🟦' },
    { key: 'aws', name: 'AWS S3', icon: '🟨' },
    { key: 'email', name: 'Email (SMTP)', icon: '📧' },
    { key: 'payment', name: 'Payment Gateway', icon: '💳' },
  ];

  if (authLoading || isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  const configuredCount = credentials
    ? Object.values(credentials).filter((v: any) => v?.configured).length
    : 0;

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Settings Dashboard</h1>
            <p className="text-gray-600 mt-2">Tenant configuration and integrations</p>
          </div>
          <Link
            href="/settings"
            className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-6 rounded-lg transition"
          >
            Manage Settings
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
                <p className="text-gray-600 text-sm font-medium">Tenant Name</p>
                <p className="text-xl font-bold text-gray-900 mt-2 truncate">{tenant?.name}</p>
              </div>
              <div className="text-4xl">🏢</div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-green-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Credentials</p>
                <p className="text-3xl font-bold text-green-600 mt-2">{configuredCount}/5</p>
              </div>
              <div className="text-4xl">🔐</div>
            </div>
            <p className="text-xs text-gray-600 mt-2">Configured</p>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-purple-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Status</p>
                <p className="text-lg font-bold text-purple-600 mt-2">{tenant?.status}</p>
              </div>
              <div className="text-4xl">✅</div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-t-4 border-t-orange-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm font-medium">Configuration</p>
                <p className="text-2xl font-bold text-orange-600 mt-2">{Math.round((configuredCount / 5) * 100)}%</p>
              </div>
              <div className="text-4xl">⚙️</div>
            </div>
          </div>
        </div>

        {/* Tenant Information */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Basic Info */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-6">Tenant Information</h3>
            <div className="space-y-4">
              <div className="pb-4 border-b">
                <p className="text-sm text-gray-600 font-medium">Tenant Name</p>
                <p className="text-lg font-semibold text-gray-900 mt-1">{tenant?.name}</p>
              </div>
              <div className="pb-4 border-b">
                <p className="text-sm text-gray-600 font-medium">Email</p>
                <p className="text-lg font-semibold text-gray-900 mt-1">{tenant?.email}</p>
              </div>
              <div className="pb-4 border-b">
                <p className="text-sm text-gray-600 font-medium">Phone</p>
                <p className="text-lg font-semibold text-gray-900 mt-1">{tenant?.phone || 'Not set'}</p>
              </div>
              <div>
                <p className="text-sm text-gray-600 font-medium">Max Users</p>
                <p className="text-lg font-semibold text-gray-900 mt-1">{tenant?.maxUsers}</p>
              </div>
            </div>
          </div>

          {/* Configuration Completion */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-6">Configuration Progress</h3>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between mb-2">
                  <span className="text-sm font-medium text-gray-700">Overall Setup</span>
                  <span className="text-sm font-bold text-gray-900">{Math.round((configuredCount / 5) * 100)}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div
                    className="bg-blue-600 h-3 rounded-full transition-all"
                    style={{ width: `${(configuredCount / 5) * 100}%` }}
                  />
                </div>
              </div>

              <div className="pt-4 border-t space-y-2">
                {configuredCount < 5 && (
                  <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
                    <p className="text-sm text-blue-800 font-medium">
                      📝 {5 - configuredCount} more credential{5 - configuredCount !== 1 ? 's' : ''} to configure
                    </p>
                  </div>
                )}
                {configuredCount === 5 && (
                  <div className="bg-green-50 border border-green-200 rounded-lg p-3">
                    <p className="text-sm text-green-800 font-medium">✅ All credentials configured!</p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Credentials Status */}
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <div className="p-6 border-b">
            <h3 className="text-lg font-semibold text-gray-900">Integration Status</h3>
          </div>

          <div className="divide-y">
            {credentialTypes.map(({ key, name, icon }) => {
              const isConfigured = credentials?.[key]?.configured || false;
              return (
                <div key={key} className="p-6 hover:bg-gray-50 transition">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <span className="text-3xl">{icon}</span>
                      <div>
                        <p className="font-semibold text-gray-900">{name}</p>
                        <p className="text-sm text-gray-600">
                          {isConfigured ? 'Configured and active' : 'Not configured'}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      <span
                        className={`px-3 py-1 rounded-full text-sm font-medium ${
                          isConfigured
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {isConfigured ? '✅ Ready' : '⏳ Pending'}
                      </span>
                      <Link
                        href="/settings"
                        className="text-blue-600 hover:text-blue-800 font-medium"
                      >
                        {isConfigured ? 'Update' : 'Configure'}
                      </Link>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Security & Compliance */}
        <div className="mt-8 bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-6">Security</h3>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 bg-green-50 rounded-lg border border-green-200">
              <div className="flex items-center gap-3">
                <span className="text-2xl">🔒</span>
                <div>
                  <p className="font-medium text-gray-900">Encryption at Rest</p>
                  <p className="text-sm text-gray-600">All credentials encrypted with AES-256</p>
                </div>
              </div>
              <span className="text-green-600 font-medium">✅ Active</span>
            </div>

            <div className="flex items-center justify-between p-4 bg-green-50 rounded-lg border border-green-200">
              <div className="flex items-center gap-3">
                <span className="text-2xl">🛡️</span>
                <div>
                  <p className="font-medium text-gray-900">Multi-Tenant Isolation</p>
                  <p className="text-sm text-gray-600">Complete data isolation per tenant</p>
                </div>
              </div>
              <span className="text-green-600 font-medium">✅ Active</span>
            </div>

            <div className="flex items-center justify-between p-4 bg-green-50 rounded-lg border border-green-200">
              <div className="flex items-center gap-3">
                <span className="text-2xl">🔐</span>
                <div>
                  <p className="font-medium text-gray-900">Credential Rotation</p>
                  <p className="text-sm text-gray-600">Rotate credentials on demand</p>
                </div>
              </div>
              <span className="text-green-600 font-medium">✅ Available</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
