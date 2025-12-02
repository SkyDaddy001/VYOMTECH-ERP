'use client'

import { useState } from 'react'
import { SectionCard } from '@/components/ui/section-card'
import { StatCard } from '@/components/ui/stat-card'

interface Report {
  id: string
  name: string
  type: string
  lastGenerated: string
  frequency: string
  format: string
}

export default function ReportsPage() {
  const [reports] = useState<Report[]>([
    { id: '1', name: 'Sales Performance', type: 'Revenue', lastGenerated: '2025-11-26', frequency: 'Weekly', format: 'PDF' },
    { id: '2', name: 'Lead Pipeline', type: 'Sales', lastGenerated: '2025-11-26', frequency: 'Daily', format: 'Excel' },
    { id: '3', name: 'Employee Attendance', type: 'HR', lastGenerated: '2025-11-25', frequency: 'Monthly', format: 'PDF' },
    { id: '4', name: 'Purchase Summary', type: 'Finance', lastGenerated: '2025-11-24', frequency: 'Monthly', format: 'Excel' },
    { id: '5', name: 'Call Analytics', type: 'Operations', lastGenerated: '2025-11-26', frequency: 'Daily', format: 'PDF' },
  ])

  const stats = [
    { label: 'Total Reports', value: reports.length, icon: '📊' },
    { label: 'This Month', value: '12', icon: '📈' },
    { label: 'Scheduled', value: '8', icon: '📅' },
    { label: 'Custom', value: '4', icon: '⚙️' },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-gradient-to-r from-indigo-600 to-indigo-800 rounded-lg p-6 text-white">
        <h1 className="text-3xl font-bold">Reports & Analytics</h1>
        <p className="text-indigo-100 mt-2">Generate and manage business reports</p>
      </div>

      {/* KPIs */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, i) => (
          <StatCard key={i} label={stat.label} value={stat.value} icon={stat.icon} />
        ))}
      </div>

      {/* Reports Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {reports.map((report) => (
          <SectionCard key={report.id} title={report.name} action={
            <div className="flex gap-2">
              <button className="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded hover:bg-blue-200">View</button>
              <button className="px-2 py-1 bg-green-100 text-green-700 text-xs rounded hover:bg-green-200">Download</button>
            </div>
          }>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">Type:</span>
                <span className="font-medium">{report.type}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Format:</span>
                <span className="font-medium">{report.format}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Frequency:</span>
                <span className="font-medium">{report.frequency}</span>
              </div>
              <div className="flex justify-between pt-2 border-t">
                <span className="text-gray-600">Last Generated:</span>
                <span className="font-medium">{report.lastGenerated}</span>
              </div>
            </div>
          </SectionCard>
        ))}
      </div>

      {/* Generate New Report Section */}
      <SectionCard title="Generate Custom Report">
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Report Type</label>
            <select className="w-full px-3 py-2 border border-gray-300 rounded-lg">
              <option>Sales</option>
              <option>HR</option>
              <option>Finance</option>
              <option>Operations</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Date Range</label>
            <div className="flex gap-2">
              <input type="date" className="flex-1 px-3 py-2 border border-gray-300 rounded-lg" />
              <input type="date" className="flex-1 px-3 py-2 border border-gray-300 rounded-lg" />
            </div>
          </div>
          <button className="w-full px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 font-medium">Generate Report</button>
        </div>
      </SectionCard>
    </div>
  )
}
