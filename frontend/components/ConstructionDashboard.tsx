'use client';

/**
 * Construction Dashboard
 * Excel-like BOQ and project tracking
 */

import React, { useState } from 'react';
import SimpleSpreadsheet, { SimpleColumn } from '@/components/SimpleSpreadsheet';
import ExcelDashboard, { DashboardTab } from '@/components/ExcelDashboard';
import { Building2, BarChart3 } from 'lucide-react';

const ConstructionDashboard = () => {
  // Sample BOQ data
  const [boqItems, setBoqItems] = useState([
    {
      id: '1',
      code: 'B001',
      description: 'Excavation',
      quantity: 500,
      unit: 'cum',
      rate: 2500,
      amount: 1250000,
      progress: 100,
    },
    {
      id: '2',
      code: 'B002',
      description: 'Concrete foundation',
      quantity: 200,
      unit: 'cum',
      rate: 8500,
      amount: 1700000,
      progress: 85,
    },
    {
      id: '3',
      code: 'B003',
      description: 'Structural steel',
      quantity: 50,
      unit: 't',
      rate: 75000,
      amount: 3750000,
      progress: 60,
    },
    {
      id: '4',
      code: 'B004',
      description: 'Brick masonry',
      quantity: 800,
      unit: 'sqm',
      rate: 1200,
      amount: 960000,
      progress: 40,
    },
  ]);

  // Sample project data
  const [projects, setProjects] = useState([
    {
      id: '1',
      projectCode: 'PRJ-001',
      name: 'Commercial Complex Phase 1',
      contractor: 'BuildCorp Ltd',
      location: 'Mumbai',
      boqValue: 75000000,
      progress: 65,
      status: 'ACTIVE',
    },
    {
      id: '2',
      projectCode: 'PRJ-002',
      name: 'Residential Towers',
      contractor: 'ConstroTech Inc',
      location: 'Bangalore',
      boqValue: 125000000,
      progress: 40,
      status: 'ACTIVE',
    },
    {
      id: '3',
      projectCode: 'PRJ-003',
      name: 'Infrastructure Highway',
      contractor: 'Roads & More',
      location: 'Delhi-NCR',
      boqValue: 250000000,
      progress: 85,
      status: 'NEARING_COMPLETION',
    },
  ]);

  // BOQ columns
  const boqColumns: SimpleColumn[] = [
    { id: 'code', label: 'Code', type: 'text', width: 80, editable: false },
    { id: 'description', label: 'Description', type: 'text', width: 180, editable: true },
    { id: 'quantity', label: 'Qty', type: 'number', width: 100, editable: true },
    { id: 'unit', label: 'Unit', type: 'text', width: 80, editable: true },
    { id: 'rate', label: 'Rate', type: 'currency', width: 120, editable: true },
    { id: 'amount', label: 'Amount', type: 'currency', width: 140, editable: false },
    { id: 'progress', label: 'Progress %', type: 'percentage', width: 120, editable: true },
  ];

  // Project columns
  const projectColumns: SimpleColumn[] = [
    { id: 'projectCode', label: 'Project Code', type: 'text', width: 120, editable: false },
    { id: 'name', label: 'Project Name', type: 'text', width: 200, editable: true },
    { id: 'contractor', label: 'Contractor', type: 'text', width: 160, editable: true },
    { id: 'location', label: 'Location', type: 'text', width: 140, editable: true },
    { id: 'boqValue', label: 'BOQ Value', type: 'currency', width: 140, editable: false },
    { id: 'progress', label: 'Progress %', type: 'percentage', width: 120, editable: true },
    { id: 'status', label: 'Status', type: 'select', width: 150, editable: true },
  ];

  // Get total BOQ value
  const totalBoqValue = boqItems.reduce((sum, item) => sum + item.amount, 0);
  const avgProgress = Math.round(boqItems.reduce((sum, item) => sum + item.progress, 0) / boqItems.length);

  // Get project stats
  const totalProjectValue = projects.reduce((sum, p) => sum + p.boqValue, 0);
  const activeProjects = projects.filter(p => p.status === 'ACTIVE').length;

  // Dashboard tabs
  const tabs: DashboardTab[] = [
    {
      id: 'boq',
      label: 'BOQ Items',
      icon: <BarChart3 size={20} />,
      content: (
        <div>
          <div className="mb-4">
            <h3 className="text-lg font-semibold text-gray-800 mb-2">Bill of Quantities</h3>
            <div className="grid grid-cols-4 gap-4 mb-4">
              <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                <p className="text-sm text-blue-600 font-semibold">Total Items</p>
                <p className="text-2xl font-bold text-blue-900">{boqItems.length}</p>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
                <p className="text-sm text-purple-600 font-semibold">BOQ Value</p>
                <p className="text-2xl font-bold text-purple-900">₹ {(totalBoqValue / 10000000).toFixed(1)}Cr</p>
              </div>
              <div className="bg-orange-50 p-4 rounded-lg border border-orange-200">
                <p className="text-sm text-orange-600 font-semibold">Avg Progress</p>
                <p className="text-2xl font-bold text-orange-900">{avgProgress}%</p>
              </div>
              <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                <p className="text-sm text-green-600 font-semibold">Executed Value</p>
                <p className="text-2xl font-bold text-green-900">₹ {(totalBoqValue * (avgProgress / 100) / 10000000).toFixed(1)}Cr</p>
              </div>
            </div>
          </div>
          <SimpleSpreadsheet
            columns={boqColumns}
            data={boqItems}
            onDataChange={setBoqItems}
            onDeleteRow={(idx) => setBoqItems(boqItems.filter((_, i) => i !== idx))}
            onAddRow={() => setBoqItems([...boqItems, {
              id: `boq-${Date.now()}`,
              code: `B${String(boqItems.length + 1).padStart(3, '0')}`,
              description: '',
              quantity: 1,
              unit: 'nos',
              rate: 0,
              amount: 0,
              progress: 0,
            }])}
            title="BOQ Register"
            showSearch={true}
            allowExport={true}
          />
        </div>
      ),
    },
    {
      id: 'projects',
      label: 'Projects',
      icon: <Building2 size={20} />,
      content: (
        <div>
          <div className="mb-4">
            <h3 className="text-lg font-semibold text-gray-800 mb-2">Construction Projects</h3>
            <div className="grid grid-cols-4 gap-4 mb-4">
              <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                <p className="text-sm text-blue-600 font-semibold">Total Projects</p>
                <p className="text-2xl font-bold text-blue-900">{projects.length}</p>
              </div>
              <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                <p className="text-sm text-green-600 font-semibold">Active</p>
                <p className="text-2xl font-bold text-green-900">{activeProjects}</p>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
                <p className="text-sm text-purple-600 font-semibold">Total Value</p>
                <p className="text-2xl font-bold text-purple-900">₹ {(totalProjectValue / 10000000).toFixed(1)}Cr</p>
              </div>
              <div className="bg-orange-50 p-4 rounded-lg border border-orange-200">
                <p className="text-sm text-orange-600 font-semibold">Avg Progress</p>
                <p className="text-2xl font-bold text-orange-900">
                  {Math.round(projects.reduce((sum, p) => sum + p.progress, 0) / projects.length)}%
                </p>
              </div>
            </div>
          </div>
          <SimpleSpreadsheet
            columns={projectColumns}
            data={projects}
            onDataChange={setProjects}
            onDeleteRow={(idx) => setProjects(projects.filter((_, i) => i !== idx))}
            onAddRow={() => setProjects([...projects, {
              id: `prj-${Date.now()}`,
              projectCode: '',
              name: '',
              contractor: '',
              location: '',
              boqValue: 0,
              progress: 0,
              status: 'PLANNING',
            }])}
            title="Projects Register"
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
        title="Construction Dashboard"
        subtitle="BOQ items and project tracking • Simple Excel-like interface"
      />
    </div>
  );
};

export default ConstructionDashboard;
