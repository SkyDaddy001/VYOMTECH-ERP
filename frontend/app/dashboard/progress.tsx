'use client';

import React, { useState, useCallback } from 'react';
import SpreadsheetGrid, { Column } from '@/components/SpreadsheetGrid';

interface ProgressRecord {
  id: string;
  date: string;
  projectId: string;
  activity: string;
  quantityCompleted: number;
  unit: string;
  percentComplete: number;
  workforce: number;
  notes: string;
}

const ProgressDashboard = () => {
  const [progressRecords, setProgressRecords] = useState<ProgressRecord[]>([
    {
      id: '1',
      date: '2024-12-01',
      projectId: 'PRJ-001',
      activity: 'Excavation Phase 1 Completion',
      quantityCompleted: 5000,
      unit: 'cum',
      percentComplete: 100,
      workforce: 45,
      notes: 'Ahead of schedule',
    },
    {
      id: '2',
      date: '2024-12-02',
      projectId: 'PRJ-001',
      activity: 'Foundation Concrete Pouring',
      quantityCompleted: 200,
      unit: 'cum',
      percentComplete: 45,
      workforce: 38,
      notes: 'Weather delays by 2 days',
    },
    {
      id: '3',
      date: '2024-12-03',
      projectId: 'PRJ-002',
      activity: 'Brick Masonry Layer 1',
      quantityCompleted: 2500,
      unit: 'sqm',
      percentComplete: 65,
      workforce: 25,
      notes: 'Normal progress',
    },
    {
      id: '4',
      date: '2024-12-04',
      projectId: 'PRJ-002',
      activity: 'Plaster on Block Level',
      quantityCompleted: 1200,
      unit: 'sqm',
      percentComplete: 30,
      workforce: 18,
      notes: 'Pending material delivery',
    },
    {
      id: '5',
      date: '2024-12-05',
      projectId: 'PRJ-003',
      activity: 'Site Preparation',
      quantityCompleted: 8000,
      unit: 'sqm',
      percentComplete: 85,
      workforce: 32,
      notes: 'On track',
    },
  ]);

  const columns: Column[] = [
    {
      id: 'date',
      header: 'Date',
      accessor: 'date',
      type: 'date',
      width: 110,
      editable: true,
      sortable: true,
    },
    {
      id: 'projectId',
      header: 'Project ID',
      accessor: 'projectId',
      type: 'text',
      width: 110,
      editable: true,
      sortable: true,
    },
    {
      id: 'activity',
      header: 'Activity',
      accessor: 'activity',
      type: 'text',
      width: 280,
      editable: true,
      sortable: true,
    },
    {
      id: 'quantityCompleted',
      header: 'Qty Completed',
      accessor: 'quantityCompleted',
      type: 'number',
      width: 130,
      editable: true,
      sortable: true,
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
      id: 'percentComplete',
      header: '% Complete',
      accessor: 'percentComplete',
      type: 'number',
      width: 110,
      editable: true,
      sortable: true,
    },
    {
      id: 'workforce',
      header: 'Workforce',
      accessor: 'workforce',
      type: 'number',
      width: 100,
      editable: true,
      sortable: true,
    },
    {
      id: 'notes',
      header: 'Notes',
      accessor: 'notes',
      type: 'text',
      width: 250,
      editable: true,
      sortable: false,
    },
  ];

  const handleDataChange = useCallback((updatedData: any[]) => {
    setProgressRecords(updatedData);
    console.log('Progress updated:', updatedData);
  }, []);

  const handleAddRow = useCallback(() => {
    const newRecord: ProgressRecord = {
      id: String(Date.now()),
      date: new Date().toISOString().split('T')[0],
      projectId: '',
      activity: '',
      quantityCompleted: 0,
      unit: 'sqm',
      percentComplete: 0,
      workforce: 0,
      notes: '',
    };
    setProgressRecords([...progressRecords, newRecord]);
  }, [progressRecords]);

  const handleDeleteRow = useCallback((rowIndex: number) => {
    setProgressRecords(progressRecords.filter((_, i) => i !== rowIndex));
  }, [progressRecords]);

  // Calculate statistics
  const stats = {
    totalRecords: progressRecords.length,
    avgProgress: Math.round(
      progressRecords.reduce((sum, p) => sum + p.percentComplete, 0) / progressRecords.length
    ),
    totalWorkforce: progressRecords.reduce((sum, p) => sum + p.workforce, 0),
    projects: new Set(progressRecords.map(p => p.projectId)).size,
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Page Header */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Progress Tracking</h1>
          <p className="text-gray-600 mt-2">
            Daily progress entries for all projects and activities.
          </p>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-2 gap-4 mb-6 lg:grid-cols-4">
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Total Entries</div>
            <div className="text-2xl font-bold text-gray-900">{stats.totalRecords}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Avg Progress</div>
            <div className="text-2xl font-bold text-blue-600">{stats.avgProgress}%</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Projects</div>
            <div className="text-2xl font-bold text-green-600">{stats.projects}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Today's Workforce</div>
            <div className="text-2xl font-bold text-orange-600">
              {progressRecords[progressRecords.length - 1]?.workforce || 0}
            </div>
          </div>
        </div>

        {/* Chart Placeholder */}
        <div className="mb-6 bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="font-semibold text-gray-800 mb-4">Progress Timeline</h3>
          <div className="h-48 bg-gradient-to-r from-blue-50 to-green-50 rounded-lg flex items-center justify-center text-gray-400">
            <div className="text-center">
              <p>Progress chart will be displayed here</p>
              <p className="text-sm mt-2">Integration with Chart.js</p>
            </div>
          </div>
        </div>

        {/* Main Grid */}
        <div className="h-[600px] bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <SpreadsheetGrid
            title="Progress Records"
            columns={columns}
            data={progressRecords}
            onDataChange={handleDataChange}
            onAddRow={handleAddRow}
            onDeleteRow={handleDeleteRow}
            showRowNumbers={true}
          />
        </div>

        {/* Best Practices */}
        <div className="mt-6 bg-green-50 border border-green-200 rounded-lg p-4">
          <h3 className="font-semibold text-green-900 mb-2">Best Practices:</h3>
          <ul className="text-sm text-green-800 space-y-1">
            <li>✓ Update daily to maintain accurate project status</li>
            <li>✓ Use consistent units across similar activities</li>
            <li>✓ Add notes for any delays or issues</li>
            <li>✓ Track workforce to monitor resource utilization</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default ProgressDashboard;
