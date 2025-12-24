'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth';
import { apiClient } from '@/lib/api-client';

interface Sale {
  id: string;
  lead_id: string;
  property_name: string;
  sale_amount: number;
  commission?: number;
  sale_date: string;
  status: string;
  payment_status: string;
  buyer_name?: string;
  buyer_email?: string;
  buyer_phone?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export default function SalesPage() {
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();
  const [sales, setSales] = useState<Sale[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showForm, setShowForm] = useState(false);
  const [newSale, setNewSale] = useState({
    lead_id: '',
    property_name: '',
    sale_amount: '',
    commission: '',
    sale_date: '',
    status: 'pipeline',
    payment_status: 'pending',
    buyer_name: '',
    buyer_email: '',
    buyer_phone: '',
    notes: '',
  });

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login');
    } else {
      fetchSales();
    }
  }, [isAuthenticated, router]);

  const fetchSales = async () => {
    try {
      setLoading(true);
      const data = await apiClient.listSales();
      setSales(Array.isArray(data) ? data : []);
      setError('');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load sales');
      setSales([]);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateSale = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const saleData = {
        lead_id: newSale.lead_id,
        property_name: newSale.property_name,
        sale_amount: parseFloat(newSale.sale_amount),
        commission: newSale.commission ? parseFloat(newSale.commission) : undefined,
        sale_date: new Date(newSale.sale_date),
        status: newSale.status,
        payment_status: newSale.payment_status,
        buyer_name: newSale.buyer_name || undefined,
        buyer_email: newSale.buyer_email || undefined,
        buyer_phone: newSale.buyer_phone || undefined,
        notes: newSale.notes || undefined,
      };

      await apiClient.createSale(saleData);
      setNewSale({
        lead_id: '',
        property_name: '',
        sale_amount: '',
        commission: '',
        sale_date: '',
        status: 'pipeline',
        payment_status: 'pending',
        buyer_name: '',
        buyer_email: '',
        buyer_phone: '',
        notes: '',
      });
      setShowForm(false);
      await fetchSales();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create sale');
    }
  };

  const handleDeleteSale = async (id: string) => {
    if (!confirm('Are you sure you want to delete this sale?')) return;

    try {
      await apiClient.deleteSale(id);
      await fetchSales();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to delete sale');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'closed':
        return 'bg-green-100 text-green-800';
      case 'negotiation':
        return 'bg-blue-100 text-blue-800';
      case 'pipeline':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getPaymentStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'overdue':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'AED',
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const calculateTotalSales = () => {
    return sales.reduce((total, sale) => total + sale.sale_amount, 0);
  };

  const calculateTotalCommission = () => {
    return sales.reduce((total, sale) => total + (sale.commission || 0), 0);
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
            <h1 className="text-3xl font-bold text-gray-900">Sales Dashboard</h1>
            <p className="text-gray-600 mt-2">Track and manage property sales</p>
          </div>
          <button
            onClick={() => setShowForm(!showForm)}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
          >
            {showForm ? 'Cancel' : 'New Sale'}
          </button>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        {/* Key Metrics */}
        {!loading && sales.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
            <div className="bg-white rounded-lg shadow p-6">
              <p className="text-gray-500 text-sm">Total Sales Value</p>
              <p className="text-2xl font-bold text-gray-900 mt-2">{formatCurrency(calculateTotalSales())}</p>
              <p className="text-xs text-gray-500 mt-2">{sales.length} properties</p>
            </div>
            <div className="bg-white rounded-lg shadow p-6">
              <p className="text-gray-500 text-sm">Total Commission</p>
              <p className="text-2xl font-bold text-green-600 mt-2">{formatCurrency(calculateTotalCommission())}</p>
            </div>
            <div className="bg-white rounded-lg shadow p-6">
              <p className="text-gray-500 text-sm">Closed Deals</p>
              <p className="text-2xl font-bold text-blue-600 mt-2">
                {sales.filter((s) => s.status === 'closed').length}
              </p>
            </div>
          </div>
        )}

        {/* Create Sale Form */}
        {showForm && (
          <div className="mb-8 bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Record New Sale</h2>
            <form onSubmit={handleCreateSale}>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Lead ID *
                  </label>
                  <input
                    type="text"
                    required
                    value={newSale.lead_id}
                    onChange={(e) => setNewSale({ ...newSale, lead_id: e.target.value })}
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
                    value={newSale.property_name}
                    onChange={(e) => setNewSale({ ...newSale, property_name: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Luxury Apartment Downtown"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Sale Amount (AED) *
                  </label>
                  <input
                    type="number"
                    required
                    value={newSale.sale_amount}
                    onChange={(e) => setNewSale({ ...newSale, sale_amount: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="2500000"
                    step="100000"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Commission (AED)</label>
                  <input
                    type="number"
                    value={newSale.commission}
                    onChange={(e) => setNewSale({ ...newSale, commission: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="75000"
                    step="10000"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Sale Date *
                  </label>
                  <input
                    type="date"
                    required
                    value={newSale.sale_date}
                    onChange={(e) => setNewSale({ ...newSale, sale_date: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
                  <select
                    value={newSale.status}
                    onChange={(e) => setNewSale({ ...newSale, status: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="pipeline">Pipeline</option>
                    <option value="negotiation">Negotiation</option>
                    <option value="closed">Closed</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Payment Status</label>
                  <select
                    value={newSale.payment_status}
                    onChange={(e) => setNewSale({ ...newSale, payment_status: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="pending">Pending</option>
                    <option value="completed">Completed</option>
                    <option value="overdue">Overdue</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Buyer Name</label>
                  <input
                    type="text"
                    value={newSale.buyer_name}
                    onChange={(e) => setNewSale({ ...newSale, buyer_name: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Ahmed Hassan"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Buyer Email</label>
                  <input
                    type="email"
                    value={newSale.buyer_email}
                    onChange={(e) => setNewSale({ ...newSale, buyer_email: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="buyer@example.com"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Buyer Phone</label>
                  <input
                    type="tel"
                    value={newSale.buyer_phone}
                    onChange={(e) => setNewSale({ ...newSale, buyer_phone: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="+971 50 123 4567"
                  />
                </div>
              </div>

              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
                <textarea
                  value={newSale.notes}
                  onChange={(e) => setNewSale({ ...newSale, notes: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Sale details and notes..."
                  rows={3}
                />
              </div>

              <button
                type="submit"
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
              >
                Record Sale
              </button>
            </form>
          </div>
        )}

        {/* Loading State */}
        {loading && (
          <div className="flex items-center justify-center py-12">
            <p className="text-gray-600">Loading sales...</p>
          </div>
        )}

        {/* Sales Grid */}
        {!loading && sales.length > 0 && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {sales.map((sale) => (
              <div key={sale.id} className="bg-white rounded-lg shadow hover:shadow-lg transition p-6">
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900">{sale.property_name}</h3>
                    <p className="text-sm text-gray-500">Lead: {sale.lead_id}</p>
                  </div>
                  <button
                    onClick={() => handleDeleteSale(sale.id)}
                    className="text-red-600 hover:text-red-800 text-sm font-medium"
                  >
                    ✕
                  </button>
                </div>

                <div className="flex gap-2 mb-4">
                  <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(sale.status)}`}>
                    {sale.status.charAt(0).toUpperCase() + sale.status.slice(1)}
                  </span>
                  <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getPaymentStatusColor(sale.payment_status)}`}>
                    {sale.payment_status.charAt(0).toUpperCase() + sale.payment_status.slice(1)}
                  </span>
                </div>

                <div className="bg-gradient-to-r from-blue-50 to-blue-100 rounded-lg p-4 mb-4">
                  <p className="text-gray-600 text-sm">Sale Amount</p>
                  <p className="text-2xl font-bold text-blue-600">{formatCurrency(sale.sale_amount)}</p>
                  {sale.commission && (
                    <p className="text-sm text-gray-600 mt-2">Commission: {formatCurrency(sale.commission)}</p>
                  )}
                </div>

                <div className="grid grid-cols-2 gap-4 mb-4 text-sm">
                  <div>
                    <p className="text-gray-500 text-xs">Sale Date</p>
                    <p className="font-medium text-gray-900">{formatDate(sale.sale_date)}</p>
                  </div>
                </div>

                {sale.buyer_name && (
                  <div className="mb-3 p-3 bg-gray-50 rounded">
                    <p className="text-xs text-gray-500">Buyer</p>
                    <p className="text-sm font-medium text-gray-900">{sale.buyer_name}</p>
                    {sale.buyer_email && (
                      <p className="text-xs text-gray-600">{sale.buyer_email}</p>
                    )}
                    {sale.buyer_phone && (
                      <p className="text-xs text-gray-600">{sale.buyer_phone}</p>
                    )}
                  </div>
                )}

                {sale.notes && (
                  <div className="mb-3 p-3 bg-yellow-50 rounded">
                    <p className="text-xs text-gray-500">Notes</p>
                    <p className="text-sm text-gray-700">{sale.notes}</p>
                  </div>
                )}

                <div className="pt-3 border-t border-gray-200">
                  <p className="text-xs text-gray-500">
                    Created: {formatDate(sale.created_at)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Empty State */}
        {!loading && sales.length === 0 && (
          <div className="bg-white rounded-lg shadow p-12 text-center">
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No sales recorded yet</h3>
            <p className="text-gray-600 mb-6">Record your first sale to get started</p>
            <button
              onClick={() => setShowForm(true)}
              className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
            >
              Record First Sale
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
