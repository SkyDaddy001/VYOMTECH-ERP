'use client';

import React, { useState, useCallback, useMemo } from 'react';
import SpreadsheetGrid, { Column } from '@/components/SpreadsheetGrid';

interface BOQItem {
  id: string;
  itemNo: number;
  description: string;
  unit: string;
  quantity: number;
  unitRate: number;
  totalAmount: number;
  category: string;
  status: string;
}

const BOQDashboard = () => {
  const [items, setItems] = useState<BOQItem[]>([
    {
      id: '1',
      itemNo: 1,
      description: 'Excavation and Earth Work',
      unit: 'cum',
      quantity: 5000,
      unitRate: 150,
      totalAmount: 750000,
      category: 'Foundation',
      status: 'completed',
    },
    {
      id: '2',
      itemNo: 2,
      description: 'Reinforced Concrete (M30)',
      unit: 'cum',
      quantity: 1200,
      unitRate: 8500,
      totalAmount: 10200000,
      category: 'Structural',
      status: 'in_progress',
    },
    {
      id: '3',
      itemNo: 3,
      description: 'Brick Masonry (12" thick)',
      unit: 'sqm',
      quantity: 8500,
      unitRate: 450,
      totalAmount: 3825000,
      category: 'Masonry',
      status: 'planned',
    },
    {
      id: '4',
      itemNo: 4,
      description: 'Plaster and Finishing',
      unit: 'sqm',
      quantity: 12000,
      unitRate: 200,
      totalAmount: 2400000,
      category: 'Finishing',
      status: 'planned',
    },
    {
      id: '5',
      itemNo: 5,
      description: 'Electrical Wiring and Cabling',
      unit: 'rm',
      quantity: 5000,
      unitRate: 50,
      totalAmount: 250000,
      category: 'Electrical',
      status: 'planned',
    },
    {
      id: '6',
      itemNo: 6,
      description: 'Plumbing Pipes and Fittings',
      unit: 'rm',
      quantity: 3500,
      unitRate: 75,
      totalAmount: 262500,
      category: 'Plumbing',
      status: 'planned',
    },
  ]);

  const columns: Column[] = [
    {
      id: 'itemNo',
      header: 'Item #',
      accessor: 'itemNo',
      type: 'number',
      width: 80,
      editable: true,
      sortable: true,
    },
    {
      id: 'description',
      header: 'Description',
      accessor: 'description',
      type: 'text',
      width: 300,
      editable: true,
      sortable: true,
    },
    {
      id: 'category',
      header: 'Category',
      accessor: 'category',
      type: 'select',
      width: 140,
      editable: true,
      sortable: true,
      filterOptions: [
        { label: 'Foundation', value: 'Foundation' },
        { label: 'Structural', value: 'Structural' },
        { label: 'Masonry', value: 'Masonry' },
        { label: 'Finishing', value: 'Finishing' },
        { label: 'Electrical', value: 'Electrical' },
        { label: 'Plumbing', value: 'Plumbing' },
      ],
    },
    {
      id: 'unit',
      header: 'Unit',
      accessor: 'unit',
      type: 'text',
      width: 80,
      editable: true,
      sortable: true,
    },
    {
      id: 'quantity',
      header: 'Quantity',
      accessor: 'quantity',
      type: 'number',
      width: 110,
      editable: true,
      sortable: true,
    },
    {
      id: 'unitRate',
      header: 'Unit Rate',
      accessor: 'unitRate',
      type: 'number',
      width: 110,
      editable: true,
      sortable: true,
    },
    {
      id: 'totalAmount',
      header: 'Total Amount',
      accessor: 'totalAmount',
      type: 'number',
      width: 130,
      editable: false,
      sortable: true,
    },
    {
      id: 'status',
      header: 'Status',
      accessor: 'status',
      type: 'select',
      width: 130,
      editable: true,
      sortable: true,
      filterOptions: [
        { label: 'Planned', value: 'planned' },
        { label: 'In Progress', value: 'in_progress' },
        { label: 'Completed', value: 'completed' },
        { label: 'On Hold', value: 'on_hold' },
      ],
    },
  ];

  const handleDataChange = useCallback((updatedData: any[]) => {
    setItems(updatedData);
    console.log('BOQ updated:', updatedData);
  }, []);

  const handleAddRow = useCallback(() => {
    const newItem: BOQItem = {
      id: String(Date.now()),
      itemNo: Math.max(...items.map(i => i.itemNo), 0) + 1,
      description: 'New Item',
      unit: 'nos',
      quantity: 0,
      unitRate: 0,
      totalAmount: 0,
      category: 'Other',
      status: 'planned',
    };
    setItems([...items, newItem]);
  }, [items]);

  const handleDeleteRow = useCallback((rowIndex: number) => {
    setItems(items.filter((_, i) => i !== rowIndex));
  }, [items]);

  // Calculate totals
  const totals = useMemo(() => {
    return {
      quantity: items.reduce((sum, item) => sum + item.quantity, 0),
      totalAmount: items.reduce((sum, item) => sum + item.totalAmount, 0),
      completed: items.filter(i => i.status === 'completed').length,
      inProgress: items.filter(i => i.status === 'in_progress').length,
    };
  }, [items]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Page Header */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Bill of Quantities</h1>
          <p className="text-gray-600 mt-2">
            Project BOQ with automatic total calculations.
          </p>
        </div>

        {/* Summary Stats */}
        <div className="grid grid-cols-2 gap-4 mb-6 lg:grid-cols-4">
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Total Items</div>
            <div className="text-2xl font-bold text-gray-900">{items.length}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Total Amount</div>
            <div className="text-2xl font-bold text-green-600">
              ₹{(totals.totalAmount / 100000).toFixed(1)}L
            </div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Completed</div>
            <div className="text-2xl font-bold text-blue-600">{totals.completed}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">In Progress</div>
            <div className="text-2xl font-bold text-orange-600">{totals.inProgress}</div>
          </div>
        </div>

        {/* Category Summary */}
        <div className="mb-6 bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <h3 className="font-semibold text-gray-800 mb-3">Summary by Category</h3>
          <div className="grid grid-cols-2 gap-4 lg:grid-cols-3">
            {['Foundation', 'Structural', 'Masonry', 'Finishing', 'Electrical', 'Plumbing'].map(
              category => {
                const categoryItems = items.filter(i => i.category === category);
                const categoryTotal = categoryItems.reduce((sum, i) => sum + i.totalAmount, 0);
                return (
                  <div key={category} className="text-sm">
                    <div className="text-gray-600">{category}</div>
                    <div className="text-lg font-semibold text-gray-900">
                      ₹{(categoryTotal / 100000).toFixed(1)}L
                    </div>
                    <div className="text-xs text-gray-500">{categoryItems.length} items</div>
                  </div>
                );
              }
            )}
          </div>
        </div>

        {/* Main Grid */}
        <div className="h-[600px] bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <SpreadsheetGrid
            title="Bill of Quantities"
            columns={columns}
            data={items}
            onDataChange={handleDataChange}
            onAddRow={handleAddRow}
            onDeleteRow={handleDeleteRow}
            showRowNumbers={true}
          />
        </div>

        {/* Totals Row */}
        <div className="mt-4 bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div className="flex justify-between items-center">
            <span className="font-semibold text-gray-800">TOTAL</span>
            <div className="text-2xl font-bold text-green-600">
              ₹{totals.totalAmount.toLocaleString('en-IN')}
            </div>
          </div>
        </div>

        {/* Tips */}
        <div className="mt-6 bg-amber-50 border border-amber-200 rounded-lg p-4">
          <h3 className="font-semibold text-amber-900 mb-2">Tips:</h3>
          <ul className="text-sm text-amber-800 space-y-1">
            <li>✓ Edit Quantity or Unit Rate - Total Amount updates automatically</li>
            <li>✓ Filter by Category to see items by type</li>
            <li>✓ Sort by Total Amount to identify high-value items</li>
            <li>✓ Track Status to monitor project progress</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default BOQDashboard;
