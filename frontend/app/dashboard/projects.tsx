'use client';

import React, { useState, useCallback } from 'react';
import SpreadsheetGrid, { Column } from '@/components/SpreadsheetGrid';
import { Plus } from 'lucide-react';

interface Project {
  id: string;
  name: string;
  code: string;
  location: string;
  client: string;
  value: number;
  status: string;
  manager: string;
  progress: number;
}

const ProjectsDashboard = () => {
  const [projects, setProjects] = useState<Project[]>([
    {
      id: '1',
      name: 'Downtown Metro Station',
      code: 'PRJ-001',
      location: 'New York, NY',
      client: 'City Transit Authority',
      value: 45000000,
      status: 'active',
      manager: 'John Smith',
      progress: 65,
    },
    {
      id: '2',
      name: 'Residential Complex A',
      code: 'PRJ-002',
      location: 'Boston, MA',
      client: 'Metropolitan Homes Inc',
      value: 28500000,
      status: 'active',
      manager: 'Sarah Johnson',
      progress: 42,
    },
    {
      id: '3',
      name: 'Commercial Plaza',
      code: 'PRJ-003',
      location: 'Chicago, IL',
      client: 'Urban Development LLC',
      value: 35000000,
      status: 'planning',
      manager: 'Mike Davis',
      progress: 15,
    },
    {
      id: '4',
      name: 'Highway Expansion',
      code: 'PRJ-004',
      location: 'Los Angeles, CA',
      client: 'State Highway Dept',
      value: 125000000,
      status: 'active',
      manager: 'Emily Wilson',
      progress: 78,
    },
  ]);

  const columns: Column[] = [
    {
      id: 'code',
      header: 'Project Code',
      accessor: 'code',
      type: 'text',
      width: 120,
      editable: true,
      sortable: true,
    },
    {
      id: 'name',
      header: 'Project Name',
      accessor: 'name',
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
      width: 180,
      editable: true,
      sortable: true,
    },
    {
      id: 'client',
      header: 'Client',
      accessor: 'client',
      type: 'text',
      width: 200,
      editable: true,
      sortable: true,
    },
    {
      id: 'manager',
      header: 'Manager',
      accessor: 'manager',
      type: 'text',
      width: 150,
      editable: true,
      sortable: true,
    },
    {
      id: 'value',
      header: 'Contract Value',
      accessor: 'value',
      type: 'number',
      width: 140,
      editable: true,
      sortable: true,
    },
    {
      id: 'progress',
      header: 'Progress %',
      accessor: 'progress',
      type: 'number',
      width: 110,
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
        { label: 'Active', value: 'active' },
        { label: 'Planning', value: 'planning' },
        { label: 'Completed', value: 'completed' },
        { label: 'Suspended', value: 'suspended' },
      ],
    },
  ];

  const handleDataChange = useCallback((updatedData: any[]) => {
    setProjects(updatedData);
    // TODO: Call API to save changes
    console.log('Projects updated:', updatedData);
  }, []);

  const handleAddRow = useCallback(() => {
    const newProject: Project = {
      id: String(Date.now()),
      name: 'New Project',
      code: `PRJ-${String(projects.length + 1).padStart(3, '0')}`,
      location: '',
      client: '',
      value: 0,
      status: 'planning',
      manager: '',
      progress: 0,
    };
    setProjects([...projects, newProject]);
  }, [projects]);

  const handleDeleteRow = useCallback((rowIndex: number) => {
    setProjects(projects.filter((_, i) => i !== rowIndex));
  }, [projects]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Page Header */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Projects</h1>
          <p className="text-gray-600 mt-2">
            Click any cell to edit. Use filters to search. Sort by clicking column headers.
          </p>
        </div>

        {/* Quick Stats */}
        <div className="grid grid-cols-4 gap-4 mb-6">
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Total Projects</div>
            <div className="text-2xl font-bold text-gray-900">{projects.length}</div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Active</div>
            <div className="text-2xl font-bold text-green-600">
              {projects.filter(p => p.status === 'active').length}
            </div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Total Value</div>
            <div className="text-2xl font-bold text-gray-900">
              ${(projects.reduce((sum, p) => sum + p.value, 0) / 1000000).toFixed(1)}M
            </div>
          </div>
          <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div className="text-sm text-gray-600">Avg Progress</div>
            <div className="text-2xl font-bold text-blue-600">
              {Math.round(projects.reduce((sum, p) => sum + p.progress, 0) / projects.length)}%
            </div>
          </div>
        </div>

        {/* Main Grid */}
        <div className="h-[600px]">
          <SpreadsheetGrid
            title="Projects List"
            columns={columns}
            data={projects}
            onDataChange={handleDataChange}
            onAddRow={handleAddRow}
            onDeleteRow={handleDeleteRow}
            showRowNumbers={true}
          />
        </div>

        {/* Instructions */}
        <div className="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
          <h3 className="font-semibold text-blue-900 mb-2">How to Use:</h3>
          <ul className="text-sm text-blue-800 space-y-1">
            <li>✓ <strong>Edit:</strong> Click any cell to edit inline</li>
            <li>✓ <strong>Sort:</strong> Click column headers to sort ascending/descending</li>
            <li>✓ <strong>Filter:</strong> Type in filter row below headers</li>
            <li>✓ <strong>Add:</strong> Click "Add Row" button to create new entry</li>
            <li>✓ <strong>Delete:</strong> Click trash icon to remove row</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default ProjectsDashboard;
