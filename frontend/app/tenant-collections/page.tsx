'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import ProtectedLayout from '@/components/ProtectedLayout';
import { BarChart, Bar, LineChart, Line, PieChart, Pie, Cell, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import axios from 'axios';

interface TenantAccount {
  id: string;
  charge_type: string;
  charge_type_name: string;
  description: string;
  is_active: boolean;
  total_collected: number;
  total_refunded: number;
  bank_account_name: string;
  bank_account_no: string;
  ifsc_code: string;
  created_at: string;
}

interface CollectionStats {
  charge_type: string;
  total_billed: number;
  total_collected: number;
  outstanding: number;
  collection_rate: number;
  invoice_count: number;
  paid_invoices: number;
  overdue_amount: number;
}

interface ClientPayment {
  id: string;
  client_id: string;
  client_name: string;
  amount: number;
  charge_type: string;
  status: string;
  payment_method: string;
  created_at: string;
  processed_at: string;
}

interface TenantDashboard {
  tenant_id: string;
  total_collected: number;
  total_outstanding: number;
  total_clients: number;
  partial_paid_invoices: number;
  overdue_invoices: number;
  collection_by_type: Record<string, CollectionStats>;
  recent_payments: ClientPayment[];
}

const CHARGE_TYPES: Record<string, { label: string; color: string; bgColor: string }> = {
  'apartment_cost': { label: 'Apartment Cost', color: '#3b82f6', bgColor: 'bg-blue-50' },
  'maintenance': { label: 'Maintenance', color: '#8b5cf6', bgColor: 'bg-purple-50' },
  'other_charges': { label: 'Other Charges', color: '#ec4899', bgColor: 'bg-pink-50' },
  'property_tax': { label: 'Property Tax', color: '#f59e0b', bgColor: 'bg-amber-50' },
  'water_charges': { label: 'Water Charges', color: '#06b6d4', bgColor: 'bg-cyan-50' },
  'electricity_tax': { label: 'Electricity Tax', color: '#10b981', bgColor: 'bg-emerald-50' },
};

export default function TenantCollectionDashboard() {
  const [dashboard, setDashboard] = useState<TenantDashboard | null>(null);
  const [accounts, setAccounts] = useState<TenantAccount[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState('overview');
  const [showNewAccountForm, setShowNewAccountForm] = useState(false);

  useEffect(() => {
    fetchDashboard();
    fetchAccounts();
  }, []);

  const fetchDashboard = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/tenant-collections/dashboard');
      setDashboard(response.data);
      setError(null);
    } catch (err) {
      setError('Failed to fetch collection dashboard');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const fetchAccounts = async () => {
    try {
      const response = await axios.get('/api/tenant-collections/accounts');
      setAccounts(response.data.accounts || []);
    } catch (err) {
      console.error('Failed to fetch accounts', err);
    }
  };

  const getStatusColor = (status: string): string => {
    switch (status) {
      case 'completed':
        return '#10b981';
      case 'pending':
        return '#f59e0b';
      case 'failed':
        return '#ef4444';
      default:
        return '#6b7280';
    }
  };

  const getChargeTypeLabel = (chargeType: string): string => {
    return CHARGE_TYPES[chargeType]?.label || chargeType;
  };

  if (loading) {
    return (
      <ProtectedLayout>
        <div className="flex items-center justify-center min-h-screen">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
        </div>
      </ProtectedLayout>
    );
  }

  if (!dashboard) {
    return (
      <ProtectedLayout>
        <div className="container mx-auto p-4">
          <Card>
            <CardContent className="pt-6">
              <p className="text-center text-gray-500">No collection data available</p>
            </CardContent>
          </Card>
        </div>
      </ProtectedLayout>
    );
  }

  const collectionData = Object.values(dashboard.collection_by_type || {}).map(stats => ({
    name: getChargeTypeLabel(stats.charge_type),
    billed: stats.total_billed,
    collected: stats.total_collected,
    outstanding: stats.outstanding,
    color: CHARGE_TYPES[stats.charge_type]?.color || '#6b7280',
  }));

  const collectionRateData = Object.values(dashboard.collection_by_type || {}).map(stats => ({
    name: getChargeTypeLabel(stats.charge_type),
    rate: parseFloat(stats.collection_rate.toFixed(1)),
  }));

  const recentPaymentsData = dashboard.recent_payments?.map((payment, index) => ({
    date: new Date(payment.created_at).toLocaleDateString(),
    amount: payment.amount,
    name: payment.client_name,
    chargeType: getChargeTypeLabel(payment.charge_type),
  })) || [];

  return (
    <ProtectedLayout>
      <div className="container mx-auto p-4 space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Collection Dashboard</h1>
            <p className="text-gray-600 mt-2">Manage client collections and payment accounts</p>
          </div>
          <Button onClick={() => setShowNewAccountForm(true)}>New Charge Account</Button>
        </div>

        {error && (
          <Card className="bg-red-50 border-red-200">
            <CardContent className="pt-6">
              <p className="text-red-800">{error}</p>
            </CardContent>
          </Card>
        )}

        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Total Collected</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">₹{dashboard.total_collected.toFixed(2)}</div>
              <p className="text-xs text-gray-500 mt-1">From all clients</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Outstanding</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-red-600">₹{dashboard.total_outstanding.toFixed(2)}</div>
              <p className="text-xs text-gray-500 mt-1">Pending payments</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Total Clients</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-blue-600">{dashboard.total_clients}</div>
              <p className="text-xs text-gray-500 mt-1">Active clients</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Overdue Invoices</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-orange-600">{dashboard.overdue_invoices}</div>
              <p className="text-xs text-gray-500 mt-1">Require attention</p>
            </CardContent>
          </Card>
        </div>

        {/* Charts */}
        {collectionData.length > 0 && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <Card>
              <CardHeader>
                <CardTitle>Collection by Charge Type</CardTitle>
                <CardDescription>Billed vs Collected amounts</CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <BarChart data={collectionData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip formatter={(value) => `₹${value}`} />
                    <Legend />
                    <Bar dataKey="billed" fill="#ccc" name="Billed" />
                    <Bar dataKey="collected" fill="#10b981" name="Collected" />
                    <Bar dataKey="outstanding" fill="#ef4444" name="Outstanding" />
                  </BarChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Collection Rate by Type</CardTitle>
                <CardDescription>Percentage collected</CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <BarChart data={collectionRateData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis domain={[0, 100]} />
                    <Tooltip formatter={(value) => `${value}%`} />
                    <Bar dataKey="rate" fill="#3b82f6" name="Collection Rate %" />
                  </BarChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Tabs */}
        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="accounts">Payment Accounts</TabsTrigger>
            <TabsTrigger value="recent">Recent Payments</TabsTrigger>
          </TabsList>

          {/* Overview Tab */}
          <TabsContent value="overview" className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {Object.entries(dashboard.collection_by_type || {}).map(([type, stats]) => {
                const chargeTypeConfig = CHARGE_TYPES[type];
                return (
                  <Card key={type} className={chargeTypeConfig?.bgColor || 'bg-gray-50'}>
                    <CardHeader>
                      <CardTitle className="flex items-center gap-2 text-lg">
                        <div
                          className="w-4 h-4 rounded"
                          style={{ backgroundColor: chargeTypeConfig?.color }}
                        ></div>
                        {getChargeTypeLabel(type)}
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div>
                        <p className="text-sm text-gray-600">Total Billed</p>
                        <p className="text-2xl font-bold">₹{stats.total_billed.toFixed(2)}</p>
                      </div>
                      <div className="grid grid-cols-2 gap-4 text-sm">
                        <div>
                          <p className="text-gray-600">Collected</p>
                          <p className="font-semibold text-green-600">₹{stats.total_collected.toFixed(2)}</p>
                        </div>
                        <div>
                          <p className="text-gray-600">Outstanding</p>
                          <p className="font-semibold text-red-600">₹{stats.outstanding.toFixed(2)}</p>
                        </div>
                      </div>
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span className="text-gray-600">Collection Rate</span>
                          <span className="font-semibold">{stats.collection_rate.toFixed(1)}%</span>
                        </div>
                        <div className="w-full bg-gray-200 rounded-full h-2">
                          <div
                            className="bg-blue-500 h-2 rounded-full"
                            style={{ width: `${stats.collection_rate}%` }}
                          ></div>
                        </div>
                      </div>
                      <div className="grid grid-cols-2 gap-2 text-xs">
                        <div>
                          <p className="text-gray-600">Total Invoices</p>
                          <p className="font-semibold">{stats.invoice_count}</p>
                        </div>
                        <div>
                          <p className="text-gray-600">Paid</p>
                          <p className="font-semibold text-green-600">{stats.paid_invoices}</p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                );
              })}
            </div>
          </TabsContent>

          {/* Payment Accounts Tab */}
          <TabsContent value="accounts" className="space-y-4">
            {accounts && accounts.length > 0 ? (
              <div className="space-y-3">
                {accounts.map((account) => (
                  <Card key={account.id}>
                    <CardContent className="pt-6">
                      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                        <div>
                          <p className="text-sm text-gray-600">Charge Type</p>
                          <p className="font-semibold text-lg">{getChargeTypeLabel(account.charge_type)}</p>
                          <p className="text-xs text-gray-500 mt-1">{account.charge_type_name}</p>
                        </div>
                        <div>
                          <p className="text-sm text-gray-600">Total Collected</p>
                          <p className="font-semibold text-lg text-green-600">₹{account.total_collected.toFixed(2)}</p>
                          <p className="text-xs text-gray-500 mt-1">Net amount</p>
                        </div>
                        <div>
                          <p className="text-sm text-gray-600">Bank Account</p>
                          <p className="font-semibold text-sm">{account.bank_account_name || 'Not configured'}</p>
                          <p className="text-xs text-gray-500 mt-1">
                            {account.bank_account_no ? `...${account.bank_account_no.slice(-4)}` : 'No account'}
                          </p>
                        </div>
                        <div className="flex justify-end items-center">
                          <div className="space-y-2">
                            <Badge variant={account.is_active ? 'default' : 'secondary'}>
                              {account.is_active ? 'Active' : 'Inactive'}
                            </Badge>
                            <Button variant="outline" size="sm" className="w-full">
                              Edit
                            </Button>
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            ) : (
              <Card>
                <CardContent className="pt-6">
                  <p className="text-center text-gray-500">No payment accounts configured. Create your first account to start collecting payments.</p>
                </CardContent>
              </Card>
            )}
          </TabsContent>

          {/* Recent Payments Tab */}
          <TabsContent value="recent" className="space-y-4">
            {recentPaymentsData && recentPaymentsData.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b">
                      <th className="text-left py-2 px-3">Date</th>
                      <th className="text-left py-2 px-3">Client Name</th>
                      <th className="text-left py-2 px-3">Charge Type</th>
                      <th className="text-left py-2 px-3">Amount</th>
                      <th className="text-left py-2 px-3">Status</th>
                    </tr>
                  </thead>
                  <tbody>
                    {dashboard.recent_payments?.map((payment) => (
                      <tr key={payment.id} className="border-b hover:bg-gray-50">
                        <td className="py-2 px-3 font-medium">
                          {new Date(payment.created_at).toLocaleDateString()}
                        </td>
                        <td className="py-2 px-3">{payment.client_name}</td>
                        <td className="py-2 px-3">
                          <Badge style={{ backgroundColor: CHARGE_TYPES[payment.charge_type]?.color }}>
                            {getChargeTypeLabel(payment.charge_type)}
                          </Badge>
                        </td>
                        <td className="py-2 px-3 font-semibold text-green-600">₹{payment.amount.toFixed(2)}</td>
                        <td className="py-2 px-3">
                          <Badge style={{ backgroundColor: getStatusColor(payment.status) }}>
                            {payment.status.charAt(0).toUpperCase() + payment.status.slice(1)}
                          </Badge>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <Card>
                <CardContent className="pt-6">
                  <p className="text-center text-gray-500">No recent payments</p>
                </CardContent>
              </Card>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </ProtectedLayout>
  );
}
