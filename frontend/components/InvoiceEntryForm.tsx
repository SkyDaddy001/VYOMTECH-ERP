'use client';

/**
 * InvoiceEntryForm Component
 * Excel-like invoice entry interface
 * Simple 2-column form (Label | Value) with auto-calculation
 */

import React, { useState, useCallback } from 'react';
import { Plus, Trash2, Calculator } from 'lucide-react';

export interface InvoiceItem {
  id: string;
  description: string;
  quantity: number;
  unitPrice: number;
  taxRate: number;
  amount: number;
}

export interface InvoiceFormData {
  invoiceNumber: string;
  date: string;
  customerName: string;
  customerEmail: string;
  taxId: string;
  items: InvoiceItem[];
  subtotal: number;
  taxAmount: number;
  total: number;
  paymentTerms: string;
  notes: string;
}

interface InvoiceEntryFormProps {
  onSave: (data: InvoiceFormData) => void;
  initialData?: Partial<InvoiceFormData>;
  tenantId?: string;
}

const InvoiceEntryForm: React.FC<InvoiceEntryFormProps> = ({
  onSave,
  initialData = {},
  tenantId,
}) => {
  const [formData, setFormData] = useState<InvoiceFormData>({
    invoiceNumber: '',
    date: new Date().toISOString().split('T')[0],
    customerName: '',
    customerEmail: '',
    taxId: '',
    items: [],
    subtotal: 0,
    taxAmount: 0,
    total: 0,
    paymentTerms: 'NET30',
    notes: '',
    ...initialData,
  });

  // Calculate totals
  const calculateTotals = useCallback((items: InvoiceItem[]) => {
    const subtotal = items.reduce((sum, item) => sum + item.amount, 0);
    const taxAmount = items.reduce(
      (sum, item) => sum + (item.amount * (item.taxRate / 100)),
      0
    );
    return { subtotal, taxAmount, total: subtotal + taxAmount };
  }, []);

  // Update form field
  const updateField = (field: keyof Omit<InvoiceFormData, 'items'>, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  // Update invoice item
  const updateItem = (index: number, field: keyof InvoiceItem, value: any) => {
    const newItems = [...formData.items];
    const item = newItems[index];

    if (field === 'quantity' || field === 'unitPrice' || field === 'taxRate') {
      const qty = field === 'quantity' ? parseFloat(value) : item.quantity;
      const price = field === 'unitPrice' ? parseFloat(value) : item.unitPrice;
      item[field] = field === 'description' ? value : parseFloat(value) || 0;
      item.amount = qty * price;
    } else {
      item[field] = value;
    }

    newItems[index] = item;
    const totals = calculateTotals(newItems);
    setFormData(prev => ({
      ...prev,
      items: newItems,
      ...totals,
    }));
  };

  // Add item
  const addItem = () => {
    const newItem: InvoiceItem = {
      id: `item-${Date.now()}`,
      description: '',
      quantity: 1,
      unitPrice: 0,
      taxRate: 18, // Default 18% GST
      amount: 0,
    };
    const newItems = [...formData.items, newItem];
    const totals = calculateTotals(newItems);
    setFormData(prev => ({
      ...prev,
      items: newItems,
      ...totals,
    }));
  };

  // Remove item
  const removeItem = (index: number) => {
    const newItems = formData.items.filter((_, i) => i !== index);
    const totals = calculateTotals(newItems);
    setFormData(prev => ({
      ...prev,
      items: newItems,
      ...totals,
    }));
  };

  // Handle save
  const handleSave = () => {
    if (!formData.invoiceNumber || !formData.customerName) {
      alert('Please fill in invoice number and customer name');
      return;
    }
    onSave(formData);
  };

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200">
      {/* Header */}
      <div className="px-6 py-4 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-t-lg">
        <h2 className="text-2xl font-bold">Invoice Entry</h2>
        <p className="text-blue-100 text-sm mt-1">Create new invoice • All required fields marked *</p>
      </div>

      <div className="p-6">
        {/* Invoice Header Section */}
        <div className="grid grid-cols-2 gap-6 mb-8 pb-8 border-b border-gray-200">
          {/* Left column */}
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Invoice Number *
              </label>
              <input
                type="text"
                value={formData.invoiceNumber}
                onChange={(e) => updateField('invoiceNumber', e.target.value)}
                placeholder="INV-001"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Invoice Date
              </label>
              <input
                type="date"
                value={formData.date}
                onChange={(e) => updateField('date', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Payment Terms
              </label>
              <select
                value={formData.paymentTerms}
                onChange={(e) => updateField('paymentTerms', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="IMMEDIATE">Immediate</option>
                <option value="NET7">Net 7 Days</option>
                <option value="NET15">Net 15 Days</option>
                <option value="NET30">Net 30 Days</option>
                <option value="NET60">Net 60 Days</option>
              </select>
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
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Customer Email
              </label>
              <input
                type="email"
                value={formData.customerEmail}
                onChange={(e) => updateField('customerEmail', e.target.value)}
                placeholder="contact@example.com"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Tax ID / GST Number
              </label>
              <input
                type="text"
                value={formData.taxId}
                onChange={(e) => updateField('taxId', e.target.value)}
                placeholder="27ABCDE1234F2Z5"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>
        </div>

        {/* Line Items Section */}
        <div className="mb-8">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Line Items</h3>
          
          {/* Items table */}
          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border-collapse">
              <thead>
                <tr className="bg-gray-100 border border-gray-300">
                  <th className="px-3 py-2 text-left font-semibold text-gray-700 min-w-48">Description</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-20">Qty</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-24">Unit Price</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-16">Tax %</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-24">Amount</th>
                  <th className="px-3 py-2 text-center font-semibold text-gray-700 w-12"></th>
                </tr>
              </thead>
              <tbody>
                {formData.items.map((item, idx) => (
                  <tr key={item.id} className="border border-gray-300 hover:bg-blue-50">
                    <td className="px-3 py-2 border-r border-gray-300">
                      <input
                        type="text"
                        value={item.description}
                        onChange={(e) => updateItem(idx, 'description', e.target.value)}
                        placeholder="e.g., Web Development Services"
                        className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300 text-right">
                      <input
                        type="number"
                        value={item.quantity}
                        onChange={(e) => updateItem(idx, 'quantity', e.target.value)}
                        min="0.01"
                        step="0.01"
                        className="w-full px-2 py-1 border border-gray-300 rounded text-right focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300 text-right">
                      <input
                        type="number"
                        value={item.unitPrice}
                        onChange={(e) => updateItem(idx, 'unitPrice', e.target.value)}
                        min="0"
                        step="0.01"
                        className="w-full px-2 py-1 border border-gray-300 rounded text-right focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300 text-right">
                      <input
                        type="number"
                        value={item.taxRate}
                        onChange={(e) => updateItem(idx, 'taxRate', e.target.value)}
                        min="0"
                        max="100"
                        step="0.01"
                        className="w-full px-2 py-1 border border-gray-300 rounded text-right focus:outline-none focus:ring-2 focus:ring-blue-500"
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
          <div className="grid grid-cols-3 gap-4 text-sm">
            <div>
              <p className="text-gray-600 mb-1">Subtotal</p>
              <p className="text-2xl font-bold text-gray-900">₹ {formData.subtotal.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-gray-600 mb-1">Tax Amount</p>
              <p className="text-2xl font-bold text-gray-900">₹ {formData.taxAmount.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-gray-600 mb-1">Total Due</p>
              <p className="text-2xl font-bold text-blue-600">₹ {formData.total.toFixed(2)}</p>
            </div>
          </div>
        </div>

        {/* Notes Section */}
        <div className="mb-8">
          <label className="block text-sm font-semibold text-gray-700 mb-2">
            Notes / Terms
          </label>
          <textarea
            value={formData.notes}
            onChange={(e) => updateField('notes', e.target.value)}
            placeholder="e.g., Thank you for your business. Payment due by..."
            rows={4}
            className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Actions */}
        <div className="flex gap-4 justify-end">
          <button
            onClick={handleSave}
            className="flex items-center gap-2 px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors font-semibold"
          >
            <Calculator size={18} />
            Save Invoice
          </button>
        </div>
      </div>
    </div>
  );
};

export default InvoiceEntryForm;
