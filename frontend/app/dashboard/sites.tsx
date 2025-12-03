'use client';

import React, { useState, useCallback } from 'react';
import SpreadsheetGrid, { Column } from '@/components/SpreadsheetGrid';
import SpreadsheetToolbar from '@/components/SpreadsheetToolbar';

interface Site {
  id: string;
  siteName: string;
  location: string;
  projectId: string;
  manager: string;
  status: string;
  area: number;
  workforce: number;
  startDate: string;
  endDate: string;
}

const SitesDashboard = () => {
  const [sites, setSites] = useState<Site[]>([
    {
      id: '1',
      siteName: 'Downtown Metro - Phase 1',
      location: 'Manhattan, NY',
      projectId: 'PRJ-001',
      manager: 'John Smith',
      status: 'active',
      area: 5000,
      workforce: 45,
      startDate: '2024-01-15',
      endDate: '2025-01-15',
    },
    {
      id: '2',
      siteName: 'Downtown Metro - Phase 2',
      location: 'Manhattan, NY',
      projectId: 'PRJ-001',
      manager: 'Sarah Johnson',
      status: 'active',
      area: 3500,
      workforce: 32,
      startDate: '2024-06-01',
      endDate: '2025-06-01',
    },
    {
      id: '3',
      siteName: 'Residential Complex A - Foundation',
      location: 'Boston, MA',
      projectId: 'PRJ-002',
      manager: 'Mike Davis',
      status: 'active',
      area: 8000,
      workforce: 78,
      startDate: '2024-03-10',
      endDate: '2024-12-10',
    },
    {
      id: '4',
      siteName: 'Commercial Plaza - Structural',
      location: 'Chicago, IL',
      projectId: 'PRJ-003',
      manager: 'Emily Wilson',
      status: 'planning',
      area: 12000,
      workforce: 0,
      startDate: '2025-02-01',
      endDate: '2026-02-01',
    },
  ]);

  const [searchTerm, setSearchTerm] = useState('');

  const columns: Column[] = [
    {
      id: 'siteName',
      header: 'Site Name',
      accessor: 'siteName',
      type: 'text',
      width: 250,
      editable: true,
      sortable: true,
    },
    {
      id: 'location',
      header: 'Location',
      accessor: 'location',
      type: 'text',
      width: 200,
      editable: true,
      sortable: true,
    },
    {
      id: 'projectId',
      header: 'Project ID',
      accessor: 'projectId',
      type: 'text',
      width: 120,
      editable: true,
      sortable: true,
    },
    {
      id: 'manager',
      header: 'Site Manager',
      accessor: 'manager',
      type: 'text',
      width: 180,
      editable: true,
      sortable: true,
    },
    {
      id: 'status',
      header: 'Status',
      accessor: 'status',
      type: 'select',
      width: 120,
      editable: true,
      sortable: true,
      filterOptions: [
        { label: 'Planning', value: 'planning' },
        { label: 'Active', value: 'active' },
        { label: 'Paused', value: 'paused' },
        { label: 'Completed', value: 'completed' },
      ],
    },
    {
      id: 'area',
      header: 'Area (m²)',
      accessor: 'area',
      type: 'number',
      width: 120,
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
      id: 'startDate',
      header: 'Start Date',
      accessor: 'startDate',
      type: 'date',
      width: 130,
      editable: true,
      sortable: true,
    },
    {
      id: 'endDate',
      header: 'End Date',
      accessor: 'endDate',
      type: 'date',
      width: 130,
      editable: true,
      sortable: true,
    },
  ];

  const handleDataChange = useCallback((updatedData: any[]) => {
    setSites(updatedData);
    console.log('Sites updated:', updatedData);
  }, []);

  const handleAddRow = useCallback(() => {
    const newSite: Site = {
      id: String(Date.now()),
      siteName: 'New Site',
      location: '',
      projectId: '',
      manager: '',
      status: 'planning',
      area: 0,
      workforce: 0,
      startDate: new Date().toISOString().split('T')[0],
      endDate: '',
    };
    setSites([...sites, newSite]);
  }, [sites]);

  const handleDeleteRow = useCallback((rowIndex: number) => {
    setSites(sites.filter((_, i) => i !== rowIndex));
  }, [sites]);

  const filteredSites = sites.filter(
    site =>
      site.siteName.toLowerCase().includes(searchTerm.toLowerCase()) ||
      site.location.toLowerCase().includes(searchTerm.toLowerCase()) ||
      site.manager.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const activeSites = sites.filter(s => s.status === 'active').length;
  const totalArea = sites.reduce((sum, s) => sum + s.area, 0);
  const totalWorkforce = sites.reduce((sum, s) => sum + s.workforce, 0);

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Page Header */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Construction Sites</h1>
          <p className="text-gray-600 mt-2">
            Manage all construction sites with real-time data entry and updates.
          </p>
        </div>

        {/* Quick Stats */}
        <div className="grid grid-cols-2 gap-4 mb-6 lg:grid-cols-4">
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Total Sites</div>
            <div className="text-2xl font-bold text-gray-900">{sites.length}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Active Sites</div>
            <div className="text-2xl font-bold text-green-600">{activeSites}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Total Area</div>
            <div className="text-2xl font-bold text-gray-900">{totalArea.toLocaleString()} m²</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Total Workforce</div>
            <div className="text-2xl font-bold text-blue-600">{totalWorkforce}</div>
          </div>
        </div>

        {/* Toolbar */}
        <SpreadsheetToolbar
          title="Sites List"
          onSearch={setSearchTerm}
          onExport={() => console.log('Export')}
          onImport={() => console.log('Import')}
        />

        {/* Main Grid */}
        <div className="h-[600px] bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <SpreadsheetGrid
            columns={columns}
            data={filteredSites}
            onDataChange={handleDataChange}
            onAddRow={handleAddRow}
            onDeleteRow={handleDeleteRow}
            showRowNumbers={true}
            densePacking={true}
          />
        </div>

        {/* Footer Help */}
        <div className="mt-4 text-xs text-gray-500">
          Tip: Use the filter row to quickly find sites. Click any cell to edit.
        </div>
      </div>
    </div>
  );
};

export default SitesDashboard;
