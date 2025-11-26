/**
 * PurchaseModule.tsx - Fixed Version (antd dependencies removed)
 * Replaced with Tailwind CSS components
 * Last Updated: November 25, 2025
 */

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import toast from 'react-hot-toast';
import { format } from 'date-fns';

const API_BASE = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// ============================================================================
// PURCHASE MODULE MAIN COMPONENT
// ============================================================================

export const PurchaseModule = () => {
  const [activeTab, setActiveTab] = useState('dashboard');

  return (
    <div className="p-6 bg-gray-50 min-h-screen">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">📦 Purchase Module - Phase 3E</h1>
        <p className="text-gray-600 mt-2">Manage vendors, purchase orders, GRN, and invoices</p>
      </div>

      {/* Tab Navigation */}
      <div className="flex gap-4 mb-6 border-b border-gray-200">
        <button
          onClick={() => setActiveTab('dashboard')}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === 'dashboard'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          📊 Dashboard
        </button>
        <button
          onClick={() => setActiveTab('vendors')}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === 'vendors'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          👥 Vendors
        </button>
        <button
          onClick={() => setActiveTab('purchase-orders')}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === 'purchase-orders'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          📋 Purchase Orders
        </button>
        <button
          onClick={() => setActiveTab('grn')}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === 'grn'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          📦 GRN
        </button>
        <button
          onClick={() => setActiveTab('invoices')}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === 'invoices'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          💰 Invoices
        </button>
      </div>

      {/* Tab Content */}
      <div className="bg-white rounded-lg shadow-md p-6">
        {activeTab === 'dashboard' && <DashboardStats />}
        {activeTab === 'vendors' && <VendorManagement />}
        {activeTab === 'purchase-orders' && <PurchaseOrderManagement />}
        {activeTab === 'grn' && <GRNManagement />}
        {activeTab === 'invoices' && <InvoiceManagement />}
      </div>
    </div>
  );
};

// ============================================================================
// DASHBOARD STATS COMPONENT
// ============================================================================

const DashboardStats = () => {
  const [stats, setStats] = useState({
    totalVendors: 0,
    activePOs: 0,
    pendingGRN: 0,
    pendingInvoices: 0,
  });

  useEffect(() => {
    // Fetch stats from API
    // setStats(...);
  }, []);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">Purchase Dashboard</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          title="Total Vendors"
          value={stats.totalVendors}
          icon="👥"
          color="bg-blue-50"
        />
        <StatCard
          title="Active POs"
          value={stats.activePOs}
          icon="📋"
          color="bg-green-50"
        />
        <StatCard
          title="Pending GRN"
          value={stats.pendingGRN}
          icon="📦"
          color="bg-yellow-50"
        />
        <StatCard
          title="Pending Invoices"
          value={stats.pendingInvoices}
          icon="💰"
          color="bg-purple-50"
        />
      </div>
    </div>
  );
};

interface StatCardProps {
  title: string;
  value: number;
  icon: string;
  color: string;
}

const StatCard: React.FC<StatCardProps> = ({ title, value, icon, color }) => {
  return (
    <div className={`${color} rounded-lg p-6 border border-gray-200`}>
      <div className="flex items-center justify-between">
        <div>
          <p className="text-gray-600 text-sm font-medium">{title}</p>
          <p className="text-3xl font-bold text-gray-900 mt-2">{value}</p>
        </div>
        <div className="text-4xl">{icon}</div>
      </div>
    </div>
  );
};

// ============================================================================
// VENDOR MANAGEMENT COMPONENT
// ============================================================================

const VendorManagement = () => {
  const [vendors, setVendors] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    address: '',
    vendor_type: 'supplier',
    payment_terms: '30',
  });

  useEffect(() => {
    fetchVendors();
  }, []);

  const fetchVendors = async () => {
    setLoading(true);
    try {
      const response = await axios.get(`${API_BASE}/api/v1/purchase/vendors`, {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenant_id'),
        },
      });
      setVendors(response.data || []);
    } catch (error) {
      toast.error('Failed to fetch vendors');
    }
    setLoading(false);
  };

  const handleCreateOrUpdate = async () => {
    try {
      if (editingId) {
        await axios.put(`${API_BASE}/api/v1/purchase/vendors/${editingId}`, formData, {
          headers: { 'X-Tenant-ID': localStorage.getItem('tenant_id') },
        });
        toast.success('Vendor updated');
      } else {
        await axios.post(`${API_BASE}/api/v1/purchase/vendors`, formData, {
          headers: { 'X-Tenant-ID': localStorage.getItem('tenant_id') },
        });
        toast.success('Vendor created');
      }
      setShowForm(false);
      setFormData({
        name: '',
        email: '',
        phone: '',
        address: '',
        vendor_type: 'supplier',
        payment_terms: '30',
      });
      setEditingId(null);
      fetchVendors();
    } catch (error) {
      toast.error('Failed to save vendor');
    }
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this vendor?')) {
      try {
        await axios.delete(`${API_BASE}/api/v1/purchase/vendors/${id}`, {
          headers: { 'X-Tenant-ID': localStorage.getItem('tenant_id') },
        });
        toast.success('Vendor deleted');
        fetchVendors();
      } catch (error) {
        toast.error('Failed to delete vendor');
      }
    }
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Vendor Management</h2>
        <button
          onClick={() => {
            setShowForm(!showForm);
            setEditingId(null);
          }}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition"
        >
          + Add Vendor
        </button>
      </div>

      {showForm && (
        <div className="bg-gray-50 rounded-lg p-6 mb-6 border border-gray-200">
          <h3 className="text-lg font-semibold mb-4">
            {editingId ? 'Edit Vendor' : 'New Vendor'}
          </h3>
          <div className="grid grid-cols-2 gap-4">
            <input
              type="text"
              placeholder="Vendor Name"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="email"
              placeholder="Email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="text"
              placeholder="Phone"
              value={formData.phone}
              onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="text"
              placeholder="Address"
              value={formData.address}
              onChange={(e) => setFormData({ ...formData, address: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <select
              value={formData.vendor_type}
              onChange={(e) => setFormData({ ...formData, vendor_type: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="supplier">Supplier</option>
              <option value="service_provider">Service Provider</option>
              <option value="manufacturer">Manufacturer</option>
            </select>
            <input
              type="text"
              placeholder="Payment Terms (days)"
              value={formData.payment_terms}
              onChange={(e) => setFormData({ ...formData, payment_terms: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div className="flex gap-2 mt-4">
            <button
              onClick={handleCreateOrUpdate}
              className="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition"
            >
              Save
            </button>
            <button
              onClick={() => setShowForm(false)}
              className="bg-gray-500 text-white px-4 py-2 rounded-lg hover:bg-gray-600 transition"
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      <div className="overflow-x-auto">
        <table className="w-full border-collapse">
          <thead>
            <tr className="bg-gray-100 border-b border-gray-200">
              <th className="px-4 py-2 text-left font-semibold">Name</th>
              <th className="px-4 py-2 text-left font-semibold">Email</th>
              <th className="px-4 py-2 text-left font-semibold">Phone</th>
              <th className="px-4 py-2 text-left font-semibold">Type</th>
              <th className="px-4 py-2 text-left font-semibold">Actions</th>
            </tr>
          </thead>
          <tbody>
            {vendors.map((vendor) => (
              <tr key={vendor.id} className="border-b border-gray-200 hover:bg-gray-50">
                <td className="px-4 py-2">{vendor.name}</td>
                <td className="px-4 py-2">{vendor.email}</td>
                <td className="px-4 py-2">{vendor.phone}</td>
                <td className="px-4 py-2">{vendor.vendor_type}</td>
                <td className="px-4 py-2">
                  <button
                    onClick={() => {
                      setFormData(vendor);
                      setEditingId(vendor.id);
                      setShowForm(true);
                    }}
                    className="text-blue-600 hover:text-blue-800 mr-2"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => handleDelete(vendor.id)}
                    className="text-red-600 hover:text-red-800"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {vendors.length === 0 && !loading && (
          <div className="text-center py-8 text-gray-500">No vendors found</div>
        )}
      </div>
    </div>
  );
};

// ============================================================================
// PURCHASE ORDER MANAGEMENT COMPONENT
// ============================================================================

const PurchaseOrderManagement = () => {
  return (
    <div className="text-center py-12 text-gray-500">
      <p className="text-lg">Purchase Order Management - Coming Soon</p>
    </div>
  );
};

// ============================================================================
// GRN MANAGEMENT COMPONENT
// ============================================================================

const GRNManagement = () => {
  return (
    <div className="text-center py-12 text-gray-500">
      <p className="text-lg">Goods Receipt Note (GRN) Management - Coming Soon</p>
    </div>
  );
};

// ============================================================================
// INVOICE MANAGEMENT COMPONENT
// ============================================================================

const InvoiceManagement = () => {
  return (
    <div className="text-center py-12 text-gray-500">
      <p className="text-lg">Invoice Management - Coming Soon</p>
    </div>
  );
};

export default PurchaseModule;

// Export individual components for use in other parts of the app
export { DashboardStats, VendorManagement, PurchaseOrderManagement, GRNManagement, InvoiceManagement };
