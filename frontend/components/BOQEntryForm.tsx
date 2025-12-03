'use client';

/**
 * BOQEntryForm Component
 * Excel-like Bill of Quantities form for construction projects
 * Precision calculations with 0.01 rupee tolerance
 */

import React, { useState, useCallback } from 'react';
import { Plus, Trash2, Building2, AlertCircle } from 'lucide-react';

export interface BOQItem {
  id: string;
  itemCode: string;
  description: string;
  specification: string;
  quantity: number;
  unit: string;
  ratePerUnit: number;
  amount: number;
  progress: number;
}

export interface BOQData {
  projectCode: string;
  projectName: string;
  contractorName: string;
  contractorContact: string;
  boqDate: string;
  projectLocation: string;
  items: BOQItem[];
  subtotal: number;
  contingency: number;
  contingencyAmount: number;
  total: number;
  progressPercentage: number;
}

interface BOQEntryFormProps {
  onSave: (data: BOQData) => void;
  initialData?: Partial<BOQData>;
  tenantId?: string;
}

const BOQEntryForm: React.FC<BOQEntryFormProps> = ({
  onSave,
  initialData = {},
  tenantId,
}) => {
  const [formData, setFormData] = useState<BOQData>({
    projectCode: '',
    projectName: '',
    contractorName: '',
    contractorContact: '',
    boqDate: new Date().toISOString().split('T')[0],
    projectLocation: '',
    items: [],
    subtotal: 0,
    contingency: 5, // Default 5% contingency
    contingencyAmount: 0,
    total: 0,
    progressPercentage: 0,
    ...initialData,
  });

  // Calculate totals (with 0.01 rupee precision)
  const calculateTotals = useCallback((items: BOQItem[], contingencyPercent: number) => {
    const subtotal = parseFloat(
      items.reduce((sum, item) => sum + item.amount, 0).toFixed(2)
    );
    const contingencyAmount = parseFloat(
      (subtotal * (contingencyPercent / 100)).toFixed(2)
    );
    const total = parseFloat((subtotal + contingencyAmount).toFixed(2));
    return { subtotal, contingencyAmount, total };
  }, []);

  // Update form field
  const updateField = (field: keyof Omit<BOQData, 'items'>, value: any) => {
    if (field === 'contingency') {
      const numValue = Math.max(0, Math.min(100, parseFloat(value) || 0));
      const totals = calculateTotals(formData.items, numValue);
      setFormData(prev => ({
        ...prev,
        [field]: numValue,
        ...totals,
      }));
    } else if (field === 'progressPercentage') {
      const numValue = Math.max(0, Math.min(100, parseFloat(value) || 0));
      setFormData(prev => ({ ...prev, [field]: numValue }));
    } else {
      setFormData(prev => ({ ...prev, [field]: value }));
    }
  };

  // Update BOQ item (with 0.01 precision)
  const updateItem = (index: number, field: keyof BOQItem, value: any) => {
    const newItems = [...formData.items];
    const item = newItems[index];

    if (field === 'quantity' || field === 'ratePerUnit') {
      item[field] = parseFloat(value) || 0;
      // Calculate with 2 decimal precision
      item.amount = parseFloat((item.quantity * item.ratePerUnit).toFixed(2));
    } else if (field === 'progress') {
      item[field] = Math.max(0, Math.min(100, parseFloat(value) || 0));
    } else {
      item[field] = value;
    }

    newItems[index] = item;
    const totals = calculateTotals(newItems, formData.contingency);
    setFormData(prev => ({
      ...prev,
      items: newItems,
      ...totals,
    }));
  };

  // Add item
  const addItem = () => {
    const newItem: BOQItem = {
      id: `boq-${Date.now()}`,
      itemCode: `B${String(formData.items.length + 1).padStart(3, '0')}`,
      description: '',
      specification: '',
      quantity: 1,
      unit: 'nos',
      ratePerUnit: 0,
      amount: 0,
      progress: 0,
    };
    const newItems = [...formData.items, newItem];
    const totals = calculateTotals(newItems, formData.contingency);
    setFormData(prev => ({
      ...prev,
      items: newItems,
      ...totals,
    }));
  };

  // Remove item
  const removeItem = (index: number) => {
    const newItems = formData.items.filter((_, i) => i !== index);
    const totals = calculateTotals(newItems, formData.contingency);
    setFormData(prev => ({
      ...prev,
      items: newItems,
      ...totals,
    }));
  };

  // Handle save
  const handleSave = () => {
    if (!formData.projectCode || !formData.projectName) {
      alert('Please fill in project code and name');
      return;
    }
    if (!formData.contractorName) {
      alert('Please enter contractor name');
      return;
    }
    if (formData.items.length === 0) {
      alert('Please add at least one BOQ item');
      return;
    }
    onSave(formData);
  };

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200">
      {/* Header */}
      <div className="px-6 py-4 bg-gradient-to-r from-orange-500 to-orange-600 text-white rounded-t-lg">
        <h2 className="text-2xl font-bold">Bill of Quantities (BOQ) Entry</h2>
        <p className="text-orange-100 text-sm mt-1">Create project BOQ with precision calculations • All required fields marked *</p>
      </div>

      <div className="p-6">
        {/* Project Header Section */}
        <div className="grid grid-cols-2 gap-6 mb-8 pb-8 border-b border-gray-200">
          {/* Left column */}
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Project Code *
              </label>
              <input
                type="text"
                value={formData.projectCode}
                onChange={(e) => updateField('projectCode', e.target.value)}
                placeholder="PRJ-001"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Project Name *
              </label>
              <input
                type="text"
                value={formData.projectName}
                onChange={(e) => updateField('projectName', e.target.value)}
                placeholder="e.g., Commercial Complex - Phase 1"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Project Location
              </label>
              <input
                type="text"
                value={formData.projectLocation}
                onChange={(e) => updateField('projectLocation', e.target.value)}
                placeholder="City, State"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
            </div>
          </div>

          {/* Right column */}
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Contractor Name *
              </label>
              <input
                type="text"
                value={formData.contractorName}
                onChange={(e) => updateField('contractorName', e.target.value)}
                placeholder="Construction Ltd."
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                Contact Number
              </label>
              <input
                type="tel"
                value={formData.contractorContact}
                onChange={(e) => updateField('contractorContact', e.target.value)}
                placeholder="+91-9876543210"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-1">
                BOQ Date
              </label>
              <input
                type="date"
                value={formData.boqDate}
                onChange={(e) => updateField('boqDate', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
            </div>
          </div>
        </div>

        {/* BOQ Items Section */}
        <div className="mb-8">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">Bill of Quantities Items</h3>
          
          {/* Items table */}
          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border-collapse">
              <thead>
                <tr className="bg-gray-100 border border-gray-300">
                  <th className="px-3 py-2 text-left font-semibold text-gray-700 w-16">Code</th>
                  <th className="px-3 py-2 text-left font-semibold text-gray-700 min-w-48">Description</th>
                  <th className="px-3 py-2 text-left font-semibold text-gray-700 w-32">Specification</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-16">Qty</th>
                  <th className="px-3 py-2 text-left font-semibold text-gray-700 w-12">Unit</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-20">Rate/Unit</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-20">Amount</th>
                  <th className="px-3 py-2 text-right font-semibold text-gray-700 w-16">Progress %</th>
                  <th className="px-3 py-2 text-center font-semibold text-gray-700 w-12"></th>
                </tr>
              </thead>
              <tbody>
                {formData.items.map((item, idx) => (
                  <tr key={item.id} className="border border-gray-300 hover:bg-orange-50">
                    <td className="px-3 py-2 border-r border-gray-300 font-mono text-gray-600">
                      {item.itemCode}
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300">
                      <input
                        type="text"
                        value={item.description}
                        onChange={(e) => updateItem(idx, 'description', e.target.value)}
                        placeholder="e.g., Excavation"
                        className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300">
                      <input
                        type="text"
                        value={item.specification}
                        onChange={(e) => updateItem(idx, 'specification', e.target.value)}
                        placeholder="e.g., 1m depth"
                        className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300 text-right">
                      <input
                        type="number"
                        value={item.quantity}
                        onChange={(e) => updateItem(idx, 'quantity', e.target.value)}
                        min="0.01"
                        step="0.01"
                        className="w-full px-2 py-1 border border-gray-300 rounded text-right focus:outline-none focus:ring-2 focus:ring-orange-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300">
                      <select
                        value={item.unit}
                        onChange={(e) => updateItem(idx, 'unit', e.target.value)}
                        className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
                      >
                        <option value="nos">nos</option>
                        <option value="sqm">sqm</option>
                        <option value="cum">cum</option>
                        <option value="lm">lm</option>
                        <option value="kg">kg</option>
                        <option value="t">t</option>
                      </select>
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300 text-right">
                      <input
                        type="number"
                        value={item.ratePerUnit}
                        onChange={(e) => updateItem(idx, 'ratePerUnit', e.target.value)}
                        min="0"
                        step="0.01"
                        className="w-full px-2 py-1 border border-gray-300 rounded text-right focus:outline-none focus:ring-2 focus:ring-orange-500"
                      />
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300 text-right font-semibold text-gray-900">
                      ₹ {item.amount.toFixed(2)}
                    </td>
                    <td className="px-3 py-2 border-r border-gray-300 text-right">
                      <input
                        type="number"
                        value={item.progress}
                        onChange={(e) => updateItem(idx, 'progress', e.target.value)}
                        min="0"
                        max="100"
                        step="0.01"
                        className="w-full px-2 py-1 border border-gray-300 rounded text-right focus:outline-none focus:ring-2 focus:ring-orange-500"
                      />
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
            Add BOQ Item
          </button>
        </div>

        {/* Summary Section */}
        <div className="bg-gray-50 rounded-lg p-4 mb-8 border border-gray-200">
          <div className="grid grid-cols-4 gap-4 mb-4">
            <div>
              <p className="text-gray-600 text-sm mb-1">Subtotal</p>
              <p className="text-2xl font-bold text-gray-900">₹ {formData.subtotal.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-gray-600 text-sm mb-1">Contingency ({formData.contingency}%)</p>
              <p className="text-2xl font-bold text-gray-900">₹ {formData.contingencyAmount.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-gray-600 text-sm mb-1">Contract Value</p>
              <p className="text-2xl font-bold text-orange-600">₹ {formData.total.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-gray-600 text-sm mb-1">Overall Progress</p>
              <p className="text-2xl font-bold text-gray-900">{formData.progressPercentage.toFixed(2)}%</p>
            </div>
          </div>

          {/* Contingency input */}
          <div className="flex items-center gap-4">
            <label className="text-sm font-semibold text-gray-700">
              Contingency %:
            </label>
            <input
              type="number"
              value={formData.contingency}
              onChange={(e) => updateField('contingency', e.target.value)}
              min="0"
              max="100"
              step="0.01"
              className="w-20 px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
            />
            <label className="text-sm font-semibold text-gray-700">
              Overall Progress %:
            </label>
            <input
              type="number"
              value={formData.progressPercentage}
              onChange={(e) => updateField('progressPercentage', e.target.value)}
              min="0"
              max="100"
              step="0.01"
              className="w-20 px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
            />
          </div>
        </div>

        {/* Notes */}
        <div className="mb-6 p-3 bg-blue-50 border border-blue-200 rounded-lg flex gap-3">
          <AlertCircle size={20} className="text-blue-600 flex-shrink-0 mt-0.5" />
          <div>
            <p className="text-sm font-semibold text-blue-900">Precision Calculations</p>
            <p className="text-sm text-blue-800 mt-1">
              All amounts calculated to 0.01 rupee precision. Progress percentage tracks completion from 0-100%.
            </p>
          </div>
        </div>

        {/* Actions */}
        <div className="flex gap-4 justify-end">
          <button
            onClick={handleSave}
            className="flex items-center gap-2 px-6 py-2 bg-orange-500 text-white rounded-lg hover:bg-orange-600 transition-colors font-semibold"
          >
            <Building2 size={18} />
            Save BOQ
          </button>
        </div>
      </div>
    </div>
  );
};

export default BOQEntryForm;
