'use client';

/**
 * Sales Dashboard
 * Excel-like sales order and invoice tracking
 */

import React, { useState } from 'react';
import SimpleSpreadsheet, { SimpleColumn } from '@/components/SimpleSpreadsheet';
import ExcelDashboard, { DashboardTab } from '@/components/ExcelDashboard';
import { FileText, ShoppingCart } from 'lucide-react';

const SalesDashboardPage = () => {
  // Sample invoice data
  const [invoices, setInvoices] = useState([
    {
      id: '1',
      invoiceNumber: 'INV-001',
      date: '2024-01-15',
      customer: 'ACME Corp',
      amount: 45000,
      tax: 8100,
      total: 53100,
      status: 'PAID',
    },
    {
      id: '2',
      invoiceNumber: 'INV-002',
      date: '2024-01-16',
      customer: 'TechStart Inc',
      amount: 32000,
      tax: 5760,
      total: 37760,
      status: 'SENT',
    },
    {
      id: '3',
      invoiceNumber: 'INV-003',
      date: '2024-01-17',
      customer: 'Global Industries',
      amount: 75000,
      tax: 13500,
      total: 88500,
      status: 'DRAFT',
    },
  ]);

  // Sample sales order data
  const [orders, setOrders] = useState([
    {
      id: '1',
      orderNumber: 'SO-001',
      date: '2024-01-15',
      customer: 'ACME Corp',
      amount: 45000,
      status: 'CONFIRMED',
      dueDate: '2024-02-15',
    },
    {
      id: '2',
      orderNumber: 'SO-002',
      date: '2024-01-16',
      customer: 'TechStart Inc',
      amount: 32000,
      status: 'PROCESSING',
      dueDate: '2024-02-20',
    },
    {
      id: '3',
      orderNumber: 'SO-003',
      date: '2024-01-17',
      customer: 'Global Industries',
      amount: 75000,
      status: 'READY',
      dueDate: '2024-02-10',
    },
  ]);

  // Invoice columns
  const invoiceColumns: SimpleColumn[] = [
    { id: 'invoiceNumber', label: 'Invoice #', type: 'text', width: 120, editable: false },
    { id: 'date', label: 'Date', type: 'date', width: 120, editable: true },
    { id: 'customer', label: 'Customer', type: 'text', width: 180, editable: true },
    { id: 'amount', label: 'Amount', type: 'currency', width: 140, editable: true },
    { id: 'tax', label: 'Tax', type: 'currency', width: 120, editable: true },
    { id: 'total', label: 'Total', type: 'currency', width: 140, editable: false },
    { id: 'status', label: 'Status', type: 'select', width: 120, editable: true },
  ];

  // Sales order columns
  const orderColumns: SimpleColumn[] = [
    { id: 'orderNumber', label: 'Order #', type: 'text', width: 120, editable: false },
    { id: 'date', label: 'Date', type: 'date', width: 120, editable: true },
    { id: 'customer', label: 'Customer', type: 'text', width: 180, editable: true },
    { id: 'amount', label: 'Amount', type: 'currency', width: 140, editable: true },
    { id: 'dueDate', label: 'Due Date', type: 'date', width: 120, editable: true },
    { id: 'status', label: 'Status', type: 'select', width: 140, editable: true },
  ];

  // Dashboard tabs
  const tabs: DashboardTab[] = [
    {
      id: 'invoices',
      label: 'Invoices',
      icon: <FileText size={20} />,
      content: (
        <div>
          <div className="mb-4">
            <h3 className="text-lg font-semibold text-gray-800 mb-2">Sales Invoices</h3>
            <div className="grid grid-cols-4 gap-4 mb-4">
              <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                <p className="text-sm text-blue-600 font-semibold">Total Invoices</p>
                <p className="text-2xl font-bold text-blue-900">{invoices.length}</p>
              </div>
              <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                <p className="text-sm text-green-600 font-semibold">Paid</p>
                <p className="text-2xl font-bold text-green-900">{invoices.filter(i => i.status === 'PAID').length}</p>
              </div>
              <div className="bg-yellow-50 p-4 rounded-lg border border-yellow-200">
                <p className="text-sm text-yellow-600 font-semibold">Pending</p>
                <p className="text-2xl font-bold text-yellow-900">{invoices.filter(i => i.status === 'SENT').length}</p>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
                <p className="text-sm text-purple-600 font-semibold">Revenue</p>
                <p className="text-2xl font-bold text-purple-900">₹ {invoices.reduce((sum, i) => sum + i.total, 0).toLocaleString('en-IN')}</p>
              </div>
            </div>
          </div>
          <SimpleSpreadsheet
            columns={invoiceColumns}
            data={invoices}
            onDataChange={setInvoices}
            onDeleteRow={(idx) => setInvoices(invoices.filter((_, i) => i !== idx))}
            onAddRow={() => setInvoices([...invoices, {
              id: `inv-${Date.now()}`,
              invoiceNumber: '',
              date: new Date().toISOString().split('T')[0],
              customer: '',
              amount: 0,
              tax: 0,
              total: 0,
              status: 'DRAFT',
            }])}
            title="Invoice Register"
            showSearch={true}
            allowExport={true}
          />
        </div>
      ),
    },
    {
      id: 'orders',
      label: 'Sales Orders',
      icon: <ShoppingCart size={20} />,
      content: (
        <div>
          <div className="mb-4">
            <h3 className="text-lg font-semibold text-gray-800 mb-2">Open Sales Orders</h3>
            <div className="grid grid-cols-4 gap-4 mb-4">
              <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                <p className="text-sm text-blue-600 font-semibold">Total Orders</p>
                <p className="text-2xl font-bold text-blue-900">{orders.length}</p>
              </div>
              <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                <p className="text-sm text-green-600 font-semibold">Ready to Dispatch</p>
                <p className="text-2xl font-bold text-green-900">{orders.filter(o => o.status === 'READY').length}</p>
              </div>
              <div className="bg-orange-50 p-4 rounded-lg border border-orange-200">
                <p className="text-sm text-orange-600 font-semibold">In Progress</p>
                <p className="text-2xl font-bold text-orange-900">{orders.filter(o => o.status === 'PROCESSING').length}</p>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
                <p className="text-sm text-purple-600 font-semibold">Order Value</p>
                <p className="text-2xl font-bold text-purple-900">₹ {orders.reduce((sum, o) => sum + o.amount, 0).toLocaleString('en-IN')}</p>
              </div>
            </div>
          </div>
          <SimpleSpreadsheet
            columns={orderColumns}
            data={orders}
            onDataChange={setOrders}
            onDeleteRow={(idx) => setOrders(orders.filter((_, i) => i !== idx))}
            onAddRow={() => setOrders([...orders, {
              id: `so-${Date.now()}`,
              orderNumber: '',
              date: new Date().toISOString().split('T')[0],
              customer: '',
              amount: 0,
              dueDate: '',
              status: 'DRAFT',
            }])}
            title="Sales Order Register"
            showSearch={true}
            allowExport={true}
          />
        </div>
      ),
    },
  ];

  return (
    <div className="h-screen flex flex-col">
      <ExcelDashboard
        tabs={tabs}
        title="Sales Dashboard"
        subtitle="Invoices and sales orders • Simple Excel-like interface"
      />
    </div>
  );
};

export default SalesDashboardPage;
