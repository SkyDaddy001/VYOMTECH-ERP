'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import ProtectedLayout from '@/components/ProtectedLayout';
import { BarChart, Bar, PieChart, Pie, Cell, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, LineChart, Line } from 'recharts';
import axios from 'axios';

interface Invoice {
  id: string;
  invoice_number: string;
  amount: number;
  amount_paid: number;
  outstanding_amount: number;
  charge_type: string;
  status: string;
  due_date: string;
  created_at: string;
}

interface OutstandingByType {
  charge_type: string;
  total: number;
  paid: number;
  outstanding: number;
}

interface OutstandingData {
  client_id: string;
  client_name: string;
  client_email: string;
  total_outstanding: number;
  total_paid: number;
  invoices: Invoice[];
  by_charge_type: Record<string, OutstandingByType>;
}

interface PaymentHistory {
  id: string;
  invoice_id: string;
  amount: number;
  status: string;
  provider: string;
  payment_method: string;
  gateway_payment_id: string;
  created_at: string;
  processed_at: string;
}

const CHARGE_TYPES: Record<string, { label: string; color: string; description: string }> = {
  'apartment_cost': { label: 'Apartment Cost', color: '#3b82f6', description: 'Property maintenance and utilities' },
  'maintenance': { label: 'Maintenance', color: '#8b5cf6', description: 'Building maintenance charges' },
  'other_charges': { label: 'Other Charges', color: '#ec4899', description: 'Additional miscellaneous charges' },
};

export default function ClientPaymentPage() {
  const [outstanding, setOutstanding] = useState<OutstandingData | null>(null);
  const [paymentHistory, setPaymentHistory] = useState<PaymentHistory[]>([]);
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState('overview');
  const [paymentMethod, setPaymentMethod] = useState('card');
  const [processingPayment, setProcessingPayment] = useState(false);

  useEffect(() => {
    fetchOutstandingBalance();
    fetchPaymentHistory();
  }, []);

  const fetchOutstandingBalance = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/client-payments/outstanding');
      setOutstanding(response.data);
      setError(null);
    } catch (err) {
      setError('Failed to fetch outstanding balance');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const fetchPaymentHistory = async () => {
    try {
      const response = await axios.get('/api/client-payments/history?limit=10');
      setPaymentHistory(response.data.payments || []);
    } catch (err) {
      console.error('Failed to fetch payment history', err);
    }
  };

  const handlePaymentInitiation = async (invoice: Invoice) => {
    if (!outstanding) return;

    try {
      setProcessingPayment(true);
      const payload = {
        invoice_id: invoice.id,
        amount: invoice.outstanding_amount,
        currency: 'INR',
        provider: 'razorpay',
        payment_method: paymentMethod,
        charge_type: invoice.charge_type,
        client_name: outstanding.client_name,
        client_email: outstanding.client_email,
        client_phone: '',
      };

      const response = await axios.post('/api/client-payments/initiate', payload);
      
      // Open payment gateway in new window or modal
      if (response.data.gateway_order) {
        // Handle gateway-specific payment flow
        window.location.href = response.data.gateway_order.payment_url;
      }
    } catch (err) {
      setError('Failed to initiate payment');
      console.error(err);
    } finally {
      setProcessingPayment(false);
    }
  };

  const getStatusColor = (status: string): string => {
    switch (status) {
      case 'paid':
        return '#10b981';
      case 'partial_paid':
        return '#f59e0b';
      case 'issued':
      case 'overdue':
        return '#ef4444';
      default:
        return '#6b7280';
    }
  };

  const getStatusLabel = (status: string): string => {
    switch (status) {
      case 'partial_paid':
        return 'Partially Paid';
      case 'apartment_cost':
      case 'maintenance':
      case 'other_charges':
        return CHARGE_TYPES[status]?.label || status;
      default:
        return status.charAt(0).toUpperCase() + status.slice(1);
    }
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

  if (!outstanding) {
    return (
      <ProtectedLayout>
        <div className="container mx-auto p-4">
          <Card>
            <CardContent className="pt-6">
              <p className="text-center text-gray-500">No outstanding invoices found</p>
            </CardContent>
          </Card>
        </div>
      </ProtectedLayout>
    );
  }

  const chargeTypeData = Object.values(outstanding.by_charge_type || {}).map(type => ({
    name: getStatusLabel(type.charge_type),
    outstanding: type.outstanding,
    paid: type.paid,
    fill: CHARGE_TYPES[type.charge_type]?.color || '#6b7280',
  }));

  const totalByChargeType = Object.entries(outstanding.by_charge_type || {}).map(([type, data]) => ({
    name: getStatusLabel(type),
    value: data.outstanding,
  }));

  return (
    <ProtectedLayout>
      <div className="container mx-auto p-4 space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Payment & Billing</h1>
          <p className="text-gray-600 mt-2">Manage your outstanding invoices and payment history</p>
        </div>

        {error && (
          <Card className="bg-red-50 border-red-200">
            <CardContent className="pt-6">
              <p className="text-red-800">{error}</p>
            </CardContent>
          </Card>
        )}

        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Total Outstanding</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-red-600">₹{outstanding.total_outstanding.toFixed(2)}</div>
              <p className="text-xs text-gray-500 mt-1">Across all charge types</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Total Paid</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">₹{outstanding.total_paid.toFixed(2)}</div>
              <p className="text-xs text-gray-500 mt-1">Total payments made</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-gray-600">Total Billed</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-blue-600">
                ₹{(outstanding.total_outstanding + outstanding.total_paid).toFixed(2)}
              </div>
              <p className="text-xs text-gray-500 mt-1">Total invoices issued</p>
            </CardContent>
          </Card>
        </div>

        {/* Charts */}
        {chargeTypeData.length > 0 && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <Card>
              <CardHeader>
                <CardTitle>Outstanding by Charge Type</CardTitle>
                <CardDescription>Amount due for each category</CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <BarChart data={chargeTypeData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip />
                    <Bar dataKey="outstanding" fill="#ef4444" name="Outstanding" />
                    <Bar dataKey="paid" fill="#10b981" name="Paid" />
                  </BarChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>

            {totalByChargeType.length > 0 && (
              <Card>
                <CardHeader>
                  <CardTitle>Outstanding Distribution</CardTitle>
                  <CardDescription>Percentage breakdown</CardDescription>
                </CardHeader>
                <CardContent>
                  <ResponsiveContainer width="100%" height={300}>
                    <PieChart>
                      <Pie
                        data={totalByChargeType}
                        cx="50%"
                        cy="50%"
                        labelLine={false}
                        label={({ name, value }) => `${name}: ₹${value}`}
                        outerRadius={80}
                        fill="#8884d8"
                        dataKey="value"
                      >
                        {totalByChargeType.map((entry, index) => (
                          <Cell key={`cell-${index}`} fill={chargeTypeData[index]?.fill || '#6b7280'} />
                        ))}
                      </Pie>
                      <Tooltip formatter={(value) => `₹${value}`} />
                    </PieChart>
                  </ResponsiveContainer>
                </CardContent>
              </Card>
            )}
          </div>
        )}

        {/* Tabs */}
        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="overview">Outstanding Invoices</TabsTrigger>
            <TabsTrigger value="history">Payment History</TabsTrigger>
            <TabsTrigger value="charges">By Charge Type</TabsTrigger>
          </TabsList>

          {/* Outstanding Invoices Tab */}
          <TabsContent value="overview" className="space-y-4">
            {outstanding.invoices && outstanding.invoices.length > 0 ? (
              <div className="space-y-3">
                {outstanding.invoices.map((invoice) => (
                  <Card key={invoice.id} className="hover:shadow-md transition-shadow">
                    <CardContent className="pt-6">
                      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                        <div className="space-y-2 flex-1">
                          <div className="flex items-center gap-2">
                            <h3 className="font-semibold text-lg">{invoice.invoice_number}</h3>
                            <Badge style={{ backgroundColor: CHARGE_TYPES[invoice.charge_type]?.color }}>
                              {getStatusLabel(invoice.charge_type)}
                            </Badge>
                            <Badge variant={invoice.status === 'paid' ? 'default' : 'secondary'}>
                              {getStatusLabel(invoice.status)}
                            </Badge>
                          </div>
                          <p className="text-sm text-gray-600">
                            Due Date: {new Date(invoice.due_date).toLocaleDateString()}
                          </p>
                          <div className="flex gap-6 text-sm">
                            <div>
                              <span className="text-gray-600">Total: </span>
                              <span className="font-semibold">₹{invoice.amount.toFixed(2)}</span>
                            </div>
                            <div>
                              <span className="text-gray-600">Paid: </span>
                              <span className="font-semibold text-green-600">₹{invoice.amount_paid.toFixed(2)}</span>
                            </div>
                            <div>
                              <span className="text-gray-600">Outstanding: </span>
                              <span className="font-semibold text-red-600">₹{invoice.outstanding_amount.toFixed(2)}</span>
                            </div>
                          </div>
                          <div className="w-full bg-gray-200 rounded-full h-2">
                            <div
                              className="bg-green-500 h-2 rounded-full"
                              style={{ width: `${(invoice.amount_paid / invoice.amount) * 100}%` }}
                            ></div>
                          </div>
                        </div>
                        <div className="flex flex-col gap-2 md:w-auto">
                          {invoice.outstanding_amount > 0 && (
                            <>
                              <select
                                value={paymentMethod}
                                onChange={(e) => setPaymentMethod(e.target.value)}
                                className="px-3 py-2 border border-gray-300 rounded-md text-sm"
                              >
                                <option value="card">Credit/Debit Card</option>
                                <option value="netbanking">Net Banking</option>
                                <option value="upi">UPI</option>
                                <option value="wallet">Digital Wallet</option>
                              </select>
                              <Button
                                onClick={() => handlePaymentInitiation(invoice)}
                                disabled={processingPayment}
                                className="w-full"
                              >
                                {processingPayment ? 'Processing...' : 'Pay Now'}
                              </Button>
                            </>
                          )}
                          {invoice.outstanding_amount === 0 && (
                            <Button disabled className="w-full">
                              Paid
                            </Button>
                          )}
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            ) : (
              <Card>
                <CardContent className="pt-6">
                  <p className="text-center text-gray-500">No outstanding invoices</p>
                </CardContent>
              </Card>
            )}
          </TabsContent>

          {/* Payment History Tab */}
          <TabsContent value="history" className="space-y-4">
            {paymentHistory && paymentHistory.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b">
                      <th className="text-left py-2 px-3">Date</th>
                      <th className="text-left py-2 px-3">Amount</th>
                      <th className="text-left py-2 px-3">Method</th>
                      <th className="text-left py-2 px-3">Status</th>
                      <th className="text-left py-2 px-3">Reference</th>
                    </tr>
                  </thead>
                  <tbody>
                    {paymentHistory.map((payment) => (
                      <tr key={payment.id} className="border-b hover:bg-gray-50">
                        <td className="py-2 px-3">
                          {new Date(payment.created_at).toLocaleDateString()}
                        </td>
                        <td className="py-2 px-3 font-semibold">₹{payment.amount.toFixed(2)}</td>
                        <td className="py-2 px-3">
                          <Badge variant="outline">{getStatusLabel(payment.payment_method)}</Badge>
                        </td>
                        <td className="py-2 px-3">
                          <Badge style={{ backgroundColor: getStatusColor(payment.status) }}>
                            {getStatusLabel(payment.status)}
                          </Badge>
                        </td>
                        <td className="py-2 px-3 font-mono text-xs">
                          {payment.gateway_payment_id?.substring(0, 12)}...
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <Card>
                <CardContent className="pt-6">
                  <p className="text-center text-gray-500">No payment history</p>
                </CardContent>
              </Card>
            )}
          </TabsContent>

          {/* By Charge Type Tab */}
          <TabsContent value="charges" className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {Object.entries(outstanding.by_charge_type || {}).map(([type, data]) => (
                <Card key={type}>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <div
                        className="w-4 h-4 rounded"
                        style={{ backgroundColor: CHARGE_TYPES[type]?.color }}
                      ></div>
                      {getStatusLabel(type)}
                    </CardTitle>
                    <CardDescription>{CHARGE_TYPES[type]?.description}</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <p className="text-sm text-gray-600">Total</p>
                      <p className="text-2xl font-bold">₹{data.total.toFixed(2)}</p>
                    </div>
                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div>
                        <p className="text-gray-600">Paid</p>
                        <p className="font-semibold text-green-600">₹{data.paid.toFixed(2)}</p>
                      </div>
                      <div>
                        <p className="text-gray-600">Outstanding</p>
                        <p className="font-semibold text-red-600">₹{data.outstanding.toFixed(2)}</p>
                      </div>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div
                        className="bg-blue-500 h-2 rounded-full"
                        style={{ width: `${(data.paid / data.total) * 100}%` }}
                      ></div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </ProtectedLayout>
  );
}
