'use client'

import { useState } from 'react'
import DashboardLayout from '@/components/layouts/DashboardLayout'

type TabType = 'dashboard' | 'projects' | 'boq' | 'progress' | 'qc'

export default function ConstructionPage() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard')

  const tabs: Array<{ id: TabType; label: string }> = [
    { id: 'dashboard', label: 'Dashboard' },
    { id: 'projects', label: 'Projects' },
    { id: 'boq', label: 'Bill of Quantities' },
    { id: 'progress', label: 'Progress Tracking' },
    { id: 'qc', label: 'Quality Control' },
  ]

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <div className="bg-gradient-to-r from-red-600 to-red-800 rounded-lg p-6 text-white">
          <h1 className="text-3xl font-bold">Construction Module</h1>
          <p className="text-red-100 mt-2">Manage construction projects, BOQ, progress, and quality control</p>
        </div>

        <div className="flex gap-2 border-b border-gray-200">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-3 font-medium border-b-2 transition ${
                activeTab === tab.id
                  ? 'border-red-600 text-red-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">Construction Module - {tabs.find((t) => t.id === activeTab)?.label}</p>
            <p className="text-gray-400 mt-2">Detailed implementation coming soon...</p>
          </div>
        </div>
      </div>
    </DashboardLayout>
  )
}
