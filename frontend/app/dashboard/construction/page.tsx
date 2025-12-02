'use client'

import { useState } from 'react'
import { TrendingUp, AlertCircle, CheckCircle, Hammer } from 'lucide-react'
import type { ConstructionProject, BillOfQuantities, ProgressTracking, QualityControl } from '@/types/construction'

type TabType = 'dashboard' | 'projects' | 'boq' | 'progress' | 'qc'

export default function ConstructionPage() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard')

  // Mock data
  const projects: ConstructionProject[] = [
    {
      id: '1',
      project_name: 'Commercial Tower Development',
      project_code: 'PROJ-2024-001',
      location: 'Central Business District',
      client: 'ABC Developers Ltd',
      contract_value: 5000000,
      start_date: '2024-01-15',
      expected_completion: '2025-12-31',
      current_progress_percentage: 45,
      status: 'active',
      project_manager: 'Mike Wilson',
    },
    {
      id: '2',
      project_name: 'Residential Complex Phase 1',
      project_code: 'PROJ-2024-002',
      location: 'Suburban Area',
      client: 'XYZ Properties',
      contract_value: 3500000,
      start_date: '2024-03-01',
      expected_completion: '2025-09-30',
      current_progress_percentage: 62,
      status: 'active',
      project_manager: 'Lisa Chen',
    },
  ]

  const boqItems: BillOfQuantities[] = [
    {
      id: '1',
      project_id: '1',
      boq_number: 'BOQ-001-001',
      item_description: 'Foundation Concrete Works',
      unit: 'M3',
      quantity: 2000,
      unit_rate: 150,
      total_amount: 300000,
      category: 'structural',
      status: 'completed',
    },
    {
      id: '2',
      project_id: '1',
      boq_number: 'BOQ-001-002',
      item_description: 'Steel Reinforcement',
      unit: 'TON',
      quantity: 450,
      unit_rate: 1200,
      total_amount: 540000,
      category: 'structural',
      status: 'in_progress',
    },
    {
      id: '3',
      project_id: '1',
      boq_number: 'BOQ-001-003',
      item_description: 'Electrical Installation',
      unit: 'LUMPSUM',
      quantity: 1,
      unit_rate: 280000,
      total_amount: 280000,
      category: 'electrical',
      status: 'planned',
    },
  ]

  const progressItems: ProgressTracking[] = [
    {
      id: '1',
      project_id: '1',
      date: '2024-11-28',
      activity_description: 'Column casting for floors 5-8',
      quantity_completed: 100,
      unit: 'M3',
      percentage_complete: 85,
      workforce_deployed: 35,
      notes: 'On schedule, excellent work quality',
    },
    {
      id: '2',
      project_id: '1',
      date: '2024-11-29',
      activity_description: 'Facade installation - East Wing',
      quantity_completed: 450,
      unit: 'SQM',
      percentage_complete: 40,
      workforce_deployed: 42,
      notes: 'Weather delays managed well',
    },
  ]

  const qcItems: QualityControl[] = [
    {
      id: '1',
      project_id: '1',
      boq_item_id: '1',
      inspection_date: '2024-11-27',
      inspector_name: 'Robert Johnson',
      quality_status: 'passed',
      observations: 'Concrete strength and finish excellent',
      follow_up_date: undefined,
    },
    {
      id: '2',
      project_id: '1',
      boq_item_id: '2',
      inspection_date: '2024-11-29',
      inspector_name: 'Diana Martinez',
      quality_status: 'partial',
      observations: 'Minor rework needed on sections B-5 to B-7',
      corrective_actions: 'Reinspection scheduled',
      follow_up_date: '2024-12-05',
    },
  ]

  const dashboard = {
    total_projects: projects.length,
    active_projects: projects.filter((p) => p.status === 'active').length,
    avg_progress_percentage:
      Math.round(projects.reduce((sum, p) => sum + p.current_progress_percentage, 0) / projects.length),
    completed_projects: 3,
    boq_items_total: boqItems.length,
    boq_items_completed: boqItems.filter((b) => b.status === 'completed').length,
    quality_pass_rate: 85,
    equipment_deployed: 24,
    workforce_deployed: 247,
    project_timeline_status: 'on_track' as const,
  }

  const tabs: Array<{ id: TabType; label: string }> = [
    { id: 'dashboard', label: 'Dashboard' },
    { id: 'projects', label: 'Projects' },
    { id: 'boq', label: 'Bill of Quantities' },
    { id: 'progress', label: 'Progress Tracking' },
    { id: 'qc', label: 'Quality Control' },
  ]

  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-red-600 to-red-800 rounded-lg p-6 text-white">
        <h1 className="text-3xl font-bold">Construction Module</h1>
        <p className="text-red-100 mt-2">Manage construction projects, BOQ, progress, and quality control</p>
      </div>

      <div className="flex gap-2 border-b border-gray-200 overflow-x-auto">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`px-4 py-3 font-medium border-b-2 transition whitespace-nowrap ${
              activeTab === tab.id
                ? 'border-red-600 text-red-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'dashboard' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Active Projects</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.active_projects}</p>
                </div>
                <Hammer className="text-red-600" size={32} />
              </div>
              <p className="text-gray-600 text-sm mt-2">of {dashboard.total_projects} total</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Avg Progress</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.avg_progress_percentage}%</p>
                </div>
                <TrendingUp className="text-blue-600" size={32} />
              </div>
              <p className="text-blue-600 text-sm mt-2">on track</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Quality Pass Rate</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.quality_pass_rate}%</p>
                </div>
                <CheckCircle className="text-green-600" size={32} />
              </div>
              <p className="text-green-600 text-sm mt-2">Excellent</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Workforce Deployed</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.workforce_deployed}</p>
                </div>
                <AlertCircle className="text-orange-600" size={32} />
              </div>
              <p className="text-gray-600 text-sm mt-2">Personnel</p>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'projects' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Project Name</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Code</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Progress</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Manager</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
              </tr>
            </thead>
            <tbody>
              {projects.map((project) => (
                <tr key={project.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm text-gray-900">{project.project_name}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{project.project_code}</td>
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-2">
                      <div className="w-24 bg-gray-200 rounded-full h-2">
                        <div
                          className="bg-blue-600 h-2 rounded-full"
                          style={{ width: `${project.current_progress_percentage}%` }}
                        ></div>
                      </div>
                      <span className="text-sm text-gray-600">{project.current_progress_percentage}%</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">{project.project_manager}</td>
                  <td className="px-6 py-4">
                    <span className="px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      {project.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeTab === 'boq' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Item Description</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Category</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Quantity</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Total Amount</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
              </tr>
            </thead>
            <tbody>
              {boqItems.map((item) => (
                <tr key={item.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm text-gray-900">{item.item_description}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{item.category}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {item.quantity} {item.unit}
                  </td>
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">${item.total_amount.toLocaleString()}</td>
                  <td className="px-6 py-4">
                    <span
                      className={`px-3 py-1 rounded-full text-xs font-medium ${
                        item.status === 'completed'
                          ? 'bg-green-100 text-green-800'
                          : item.status === 'in_progress'
                            ? 'bg-blue-100 text-blue-800'
                            : 'bg-gray-100 text-gray-800'
                      }`}
                    >
                      {item.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeTab === 'progress' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Activity</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Date</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Completion</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Workforce</th>
              </tr>
            </thead>
            <tbody>
              {progressItems.map((item) => (
                <tr key={item.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm text-gray-900">{item.activity_description}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{item.date}</td>
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-2">
                      <div className="w-16 bg-gray-200 rounded-full h-2">
                        <div
                          className="bg-green-600 h-2 rounded-full"
                          style={{ width: `${item.percentage_complete}%` }}
                        ></div>
                      </div>
                      <span className="text-sm text-gray-600">{item.percentage_complete}%</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">{item.workforce_deployed}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeTab === 'qc' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Inspection Date</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Inspector</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Observations</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
              </tr>
            </thead>
            <tbody>
              {qcItems.map((item) => (
                <tr key={item.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm text-gray-600">{item.inspection_date}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{item.inspector_name}</td>
                  <td className="px-6 py-4 text-sm text-gray-900">{item.observations}</td>
                  <td className="px-6 py-4">
                    <span
                      className={`px-3 py-1 rounded-full text-xs font-medium ${
                        item.quality_status === 'passed'
                          ? 'bg-green-100 text-green-800'
                          : item.quality_status === 'partial'
                            ? 'bg-yellow-100 text-yellow-800'
                            : 'bg-red-100 text-red-800'
                      }`}
                    >
                      {item.quality_status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
