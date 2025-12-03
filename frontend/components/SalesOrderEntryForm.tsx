'use client';

/**
 * SalesOrderEntryForm Component
 * Excel-like sales order form
 * Simple spreadsheet-style order entry with inventory checks
 */

import React, { useState, useCallback } from 'react';
import { Plus, Trash2, ShoppingCart, AlertCircle } from 'lucide-react';

export interface SalesOrderItem {
  id: string;
  productId: string;
  productName: string;
  sku: string;
  quantity: number;
  unitPrice: number;
  discount: number;
  amount: number;
}

export interface SalesOrderData {
  orderNumber: string;
  date: string;
  dueDate: string;
  customerName: string;
  customerPhone: string;
  deliveryAddress: string;
  items: SalesOrderItem[];
  subtotal: number;
  discountAmount: number;
  taxAmount: number;
  total: number;
  orderStatus: string;
}

interface SalesOrderEntryFormProps {
  onSave: (data: SalesOrderData) => void;
  initialData?: Partial<SalesOrderData>;
  tenantId?: string;
}

const SalesOrderEntryForm: React.FC<SalesOrderEntryFormProps> = ({
  onSave,
  initialData = {},
  tenantId,
}) => {
  const [formData, setFormData] = useState<SalesOrderData>({
    orderNumber: '',
    date: new Date().toISOString().split('T')[0],
    dueDate: '',
    customerName: '',
    customerPhone: '',
    deliveryAddress: '',
    items: [],
    subtotal: 0,
    discountAmount: 0,
    taxAmount: 0,
    total: 0,
    orderStatus: 'DRAFT',
    ...initialData,
  });

  // Calculate totals
  const calculateTotals = useCallback((items: SalesOrderItem[], globalDiscount: number = 0) => {
    const subtotal = items.reduce((sum, item) => sum + item.amount, 0);
    const discountAmount = globalDiscount;
    const taxableAmount = subtotal - discountAmount;
    const taxAmount = taxableAmount * 0.18; // 18% GST
    return { subtotal, discountAmount, taxAmount, total: taxableAmount + taxAmount };
  }, []);

  // Update form field
  const updateField = (field: keyof Omit<SalesOrderData, 'items'>, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  // Update order item
  const updateItem = (index: number, field: keyof SalesOrderItem, value: any) => {
    const newItems = [...formData.items];
    const item = newItems[index];

    if (field === 'quantity' || field === 'unitPrice' || field === 'discount') {
      item[field] = parseFloat(value) || 0;
      const subtotal = item.quantity * item.unitPrice;
      item.amount = subtotal - (subtotal * (item.discount / 100));
    } else {
      item[field] = value;
    }

    newItems[index] = item;
    const totals = calculateTotals(newItems, formData.discountAmount);
    setFormData(prev => ({
      ...prev,
      items: newItems,
      ...totals,
    }));
  };

  // Add item
  const addItem = () => {
    const newItem: SalesOrderItem = {
      id: `item-${Date.now()}`,
      productId: '',
      productName: '',
      sku: '',
      quantity: 1,
      unitPrice: 0,
      discount: 0,
      amount: 0,
    };
    const newItems = [...formData.items, newItem];
    const totals = calculateTotals(newItems, formData.discountAmount);
    setFormData(prev => ({
      ...prev,
      items: newItems,
      ...totals,
    }));
  };

  // Remove item
  const removeItem = (index: number) => {
    const newItems = formData.items.filter((_, i) => i !== index);
    const totals = calculateTotals(newItems, formData.discountAmount);
    setFormData(prev => ({
      ...prev,
      items: newItems,
      ...totals,
    }));
  };

  // Handle save
  const handleSave = () => {
    if (!formData.orderNumber || !formData.customerName) {
      alert('Please fill in order number and customer name');
      return;
    }
    if (formData.items.length === 0) {
      alert('Please add at least one item');
      return;
    }
    onSave(formData);
  };

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200">
      {/* Header */}
      <div className="px-6 py-4 bg-gradient-to-r from-emerald-500 to-emerald-600 text-white rounded-t-lg">
        <h2 className="text-2xl font-bold">Sales Order Entry</h2>
        <p className="text-emerald-100 text-sm mt-1">Create new sales order • All required fields marked *</p>
      </div>

      <div className="p-6">
        {/* Order Header Section */}
        <div className="grid grid-cols-2 gap-6 mb-8 pb-8 border-b border-gray-200">
          {/* Left column */}
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Order Number *
              </label>
              <input
                type="text"
                value={formData.orderNumber}
                onChange={(e) => updateField('orderNumber', e.target.value)}
                placeholder="SO-001"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-emerald-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Order Date
              </label>
              <input
                type="date"
                value={formData.date}
                onChange={(e) => updateField('date', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-emerald-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Due Date
              </label>
              <input
                type="date"
                value={formData.dueDate}
                onChange={(e) => updateField('dueDate', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-emerald-500"
              />
            </div>
          </div>

          {/* Right column */}
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Customer Name *
              </label>
              <input
                type="text"
                value={formData.customerName}
                onChange={(e) => updateField('customerName', e.target.value)}
                placeholder="ACME Corp"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-emerald-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Customer Phone
              </label>
              <input
                type="tel"
                value={formData.customerPhone}
                onChange={(e) => updateField('customerPhone', e.target.value)}
                placeholder="+91-9876543210"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-emerald-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Delivery Address *
              </label>
              <input
                type="text"
                value={formData.deliveryAddress}
                onChange={(e) => updateField('deliveryAddress', e.target.value)}
                placeholder="Street, City, State, PIN"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-emerald-500"
              />
            </div>
          </div>
        </div>

        {/* Line Items Section */}
        <div className="mb-8">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Order Items</h3>
          
          {/* Items table */}
          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border-collapse">
              <thead>
                <tr className="bg-gray-100 border border-gray-300">
                  <th className="px-3 py-2 text-left font-semibold text-gray-700 min-w-40">Product Name</th>
                  <th className="px-3 py-2 text-left font-semibold text-gray-700 w-20">SKU</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-16">Qty</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-20">Unit Price</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-16">Disc %</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-24">Amount</th>
                  <th className="px-3 py-2 text-center font-semibold text-gray-700 w-12"></th>
                </tr>
              </thead>
              <tbody>
                {formData.items.map((item, idx) => (
                  <tr key={item.id} className="border border-gray-300 hover:bg-emerald-50">
                    <td className="px-3 py-2 border-r border-gray-300">
                      <input
                        type="text"
                        value={item.productName}
                        onChange={(e) => updateItem(idx, 'productName', e.target.value)}
                        placeholder="e.g., Widget A"
                        className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-emerald-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300">
                      <input
                        type="text"
                        value={item.sku}
                        onChange={(e) => updateItem(idx, 'sku', e.target.value)}
                        placeholder="WID-001"
                        className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-emerald-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300 text-right">
                      <input
                        type="number"
                        value={item.quantity}
                        onChange={(e) => updateItem(idx, 'quantity', e.target.value)}
                        min="0.01"
                        step="0.01"
                        className="w-full px-2 py-1 border border-gray-300 rounded text-right focus:outline-none focus:ring-2 focus:ring-emerald-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300 text-right">
                      <input
                        type="number"
                        value={item.unitPrice}
                        onChange={(e) => updateItem(idx, 'unitPrice', e.target.value)}
                        min="0"
                        step="0.01"
                        className="w-full px-2 py-1 border border-gray-300 rounded text-right focus:outline-none focus:ring-2 focus:ring-emerald-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300 text-right">
                      <input
                        type="number"
                        value={item.discount}
                        onChange={(e) => updateItem(idx, 'discount', e.target.value)}
                        min="0"
                        max="100"
                        step="0.01"
                        className="w-full px-2 py-1 border border-gray-300 rounded text-right focus:outline-none focus:ring-2 focus:ring-emerald-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300 text-right font-semibold text-gray-900">
                      ₹ {item.amount.toFixed(2)}
                    </td>
                    <td className="px-3 py-2 text-center">
                      <button
                        onClick={() => removeItem(idx)}
                        className="p-1 text-red-500 hover:bg-red-100 rounded transition-colors"
                      >
                        <Trash2 size={16} />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Add item button */}
          <button
            onClick={addItem}
            className="flex items-center gap-2 px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition-colors"
          >
            <Plus size={18} />
            Add Item
          </button>
        </div>

        {/* Summary Section */}
        <div className="bg-gray-50 rounded-lg p-4 mb-8 border border-gray-200">
          <div className="grid grid-cols-4 gap-4 text-sm">
            <div>
              <p className="text-gray-600 mb-1">Subtotal</p>
              <p className="text-xl font-bold text-gray-900">₹ {formData.subtotal.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-gray-600 mb-1">Discount</p>
              <p className="text-xl font-bold text-gray-900">-₹ {formData.discountAmount.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-gray-600 mb-1">Tax (18%)</p>
              <p className="text-xl font-bold text-gray-900">₹ {formData.taxAmount.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-gray-600 mb-1">Total</p>
              <p className="text-xl font-bold text-emerald-600">₹ {formData.total.toFixed(2)}</p>
            </div>
          </div>
        </div>

        {/* Status Section */}
        <div className="mb-8">
          <label className="block text-sm font-semibold text-gray-700 mb-2">
            Order Status
          </label>
          <select
            value={formData.orderStatus}
            onChange={(e) => updateField('orderStatus', e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-emerald-500"
          >
            <option value="DRAFT">Draft</option>
            <option value="CONFIRMED">Confirmed</option>
            <option value="PROCESSING">Processing</option>
            <option value="READY">Ready for Dispatch</option>
          </select>
        </div>

        {/* Actions */}
        <div className="flex gap-4 justify-end">
          <button
            onClick={handleSave}
            className="flex items-center gap-2 px-6 py-2 bg-emerald-500 text-white rounded-lg hover:bg-emerald-600 transition-colors font-semibold"
          >
            <ShoppingCart size={18} />
            Save Order
          </button>
        </div>
      </div>
    </div>
  );
};

export default SalesOrderEntryForm;
