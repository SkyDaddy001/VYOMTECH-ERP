'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/lib/hooks/useAuth';
import tenantService, { type Tenant, type TenantSettings } from '@/services/tenant-service';

export default function SettingsPage() {
  const { user, isLoading: authLoading } = useAuth();
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [settings, setSettings] = useState<TenantSettings | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'general' | 'credentials' | 'sso'>('general');

  useEffect(() => {
    if (authLoading || !user) return;
    loadData();
  }, [authLoading, user]);

  const loadData = async () => {
    try {
      setIsLoading(true);
      const [tenantData, settingsData] = await Promise.all([
        tenantService.getCurrentTenant(),
        tenantService.getTenantSettings(user?.tenantId || ''),
      ]);
      setTenant(tenantData);
      setSettings(settingsData);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load settings');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSaveGeneral = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!tenant || !user) return;

    try {
      setIsSaving(true);
      setError(null);
      setSuccess(null);

      await tenantService.updateTenant(user.tenantId, {
        name: tenant.name,
        email: tenant.email,
        phone: tenant.phone,
        maxUsers: tenant.maxUsers,
      });

      setSuccess('Tenant settings updated successfully');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save settings');
    } finally {
      setIsSaving(false);
    }
  };

  const handleSaveCredential = async (
    credentialType: 'google' | 'meta' | 'aws' | 'email' | 'razorpay' | 'billdesk',
    credentialData: Record<string, any>
  ) => {
    if (!user) return;

    try {
      setIsSaving(true);
      setError(null);
      await tenantService.storeCredential(user.tenantId, credentialType, credentialData);
      setSuccess(`${credentialType} credential saved successfully`);
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save credential');
    } finally {
      setIsSaving(false);
    }
  };

  if (authLoading || isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Tenant Settings</h1>
          <p className="text-gray-600 mt-2">Configure your account and integrations</p>
        </div>

        {/* Alerts */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4">
            {error}
          </div>
        )}
        {success && (
          <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg mb-4">
            {success}
          </div>
        )}

        {/* Tabs */}
        <div className="bg-white rounded-lg shadow mb-6">
          <div className="flex border-b">
            {(['general', 'credentials', 'sso'] as const).map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`flex-1 py-4 px-6 font-medium text-center transition ${
                  activeTab === tab
                    ? 'text-blue-600 border-b-2 border-blue-600'
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                {tab === 'general' && 'General'}
                {tab === 'credentials' && 'Credentials'}
                {tab === 'sso' && 'SSO'}
              </button>
            ))}
          </div>

          <div className="p-8">
            {/* General Tab */}
            {activeTab === 'general' && tenant && (
              <form onSubmit={handleSaveGeneral} className="space-y-6 max-w-2xl">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Tenant Name
                  </label>
                  <input
                    type="text"
                    value={tenant.name}
                    onChange={(e) => setTenant({ ...tenant, name: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Email
                  </label>
                  <input
                    type="email"
                    value={tenant.email}
                    onChange={(e) => setTenant({ ...tenant, email: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Phone
                  </label>
                  <input
                    type="tel"
                    value={tenant.phone || ''}
                    onChange={(e) => setTenant({ ...tenant, phone: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Max Users
                  </label>
                  <input
                    type="number"
                    value={tenant.maxUsers}
                    onChange={(e) => setTenant({ ...tenant, maxUsers: parseInt(e.target.value) })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>

                <button
                  type="submit"
                  disabled={isSaving}
                  className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-semibold py-2 px-6 rounded-lg transition"
                >
                  {isSaving ? 'Saving...' : 'Save Changes'}
                </button>
              </form>
            )}

            {/* Credentials Tab */}
            {activeTab === 'credentials' && (
              <div className="space-y-8 max-w-2xl">
                <CredentialForm
                  type="google"
                  title="Google OAuth"
                  fields={['client_id', 'client_secret', 'redirect_uri']}
                  onSave={(data) => handleSaveCredential('google', data)}
                  isSaving={isSaving}
                />
                <CredentialForm
                  type="meta"
                  title="Meta (Facebook) OAuth"
                  fields={['app_id', 'app_secret', 'redirect_uri']}
                  onSave={(data) => handleSaveCredential('meta', data)}
                  isSaving={isSaving}
                />
                <CredentialForm
                  type="aws"
                  title="AWS S3"
                  fields={['access_key_id', 'secret_access_key', 'region', 'bucket_name']}
                  onSave={(data) => handleSaveCredential('aws', data)}
                  isSaving={isSaving}
                />
                <CredentialForm
                  type="email"
                  title="Email (SMTP)"
                  fields={['host', 'port', 'username', 'password', 'from_email']}
                  onSave={(data) => handleSaveCredential('email', data)}
                  isSaving={isSaving}
                />
              </div>
            )}

            {/* SSO Tab */}
            {activeTab === 'sso' && settings && (
              <form
                onSubmit={(e) => {
                  e.preventDefault();
                  if (!user || !settings) return;
                  setIsSaving(true);
                  tenantService
                    .updateTenantSettings(user.tenantId, settings)
                    .then(() => {
                      setSuccess('SSO settings updated');
                      setTimeout(() => setSuccess(null), 3000);
                    })
                    .catch((err) => {
                      setError(err instanceof Error ? err.message : 'Failed to save');
                    })
                    .finally(() => setIsSaving(false));
                }}
                className="space-y-6 max-w-2xl"
              >
                <div className="flex items-center gap-4">
                  <input
                    type="checkbox"
                    id="sso"
                    checked={settings.enableSSO}
                    onChange={(e) => setSettings({ ...settings, enableSSO: e.target.checked })}
                    className="h-4 w-4 rounded border-gray-300"
                  />
                  <label htmlFor="sso" className="font-medium text-gray-700">
                    Enable Single Sign-On (SSO)
                  </label>
                </div>

                {settings.enableSSO && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      SSO Provider
                    </label>
                    <select
                      value={settings.ssoProvider || 'saml'}
                      onChange={(e) => setSettings({ ...settings, ssoProvider: e.target.value })}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    >
                      <option value="saml">SAML 2.0</option>
                      <option value="oidc">OpenID Connect</option>
                      <option value="azure">Azure AD</option>
                      <option value="okta">Okta</option>
                    </select>
                  </div>
                )}

                <button
                  type="submit"
                  disabled={isSaving}
                  className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-semibold py-2 px-6 rounded-lg transition"
                >
                  {isSaving ? 'Saving...' : 'Save SSO Settings'}
                </button>
              </form>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

interface CredentialFormProps {
  type: string;
  title: string;
  fields: string[];
  onSave: (data: Record<string, string>) => void;
  isSaving: boolean;
}

function CredentialForm({ type, title, fields, onSave, isSaving }: CredentialFormProps) {
  const [data, setData] = useState<Record<string, string>>(
    Object.fromEntries(fields.map((f) => [f, '']))
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(data);
  };

  return (
    <form onSubmit={handleSubmit} className="border rounded-lg p-6">
      <h3 className="font-semibold text-gray-900 mb-4">{title}</h3>
      <div className="space-y-4">
        {fields.map((field) => (
          <div key={field}>
            <label className="block text-sm font-medium text-gray-700 mb-2 capitalize">
              {field.replace(/_/g, ' ')}
            </label>
            <input
              type={field.includes('secret') || field.includes('password') ? 'password' : 'text'}
              value={data[field] || ''}
              onChange={(e) => setData({ ...data, [field]: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        ))}
      </div>
      <button
        type="submit"
        disabled={isSaving}
        className="mt-6 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-semibold py-2 px-6 rounded-lg transition"
      >
        {isSaving ? 'Saving...' : `Save ${title}`}
      </button>
    </form>
  );
}
