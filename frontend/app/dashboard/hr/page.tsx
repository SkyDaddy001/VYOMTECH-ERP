'use client'

import { useState } from 'react'
import DashboardLayout from '@/components/layouts/DashboardLayout'

type TabType = 'dashboard' | 'employees' | 'attendance' | 'leave' | 'payroll'

export default function HRPage() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard')

  const tabs: Array<{ id: TabType; label: string }> = [
    { id: 'dashboard', label: 'Dashboard' },
    { id: 'employees', label: 'Employees' },
    { id: 'attendance', label: 'Attendance' },
    { id: 'leave', label: 'Leave Management' },
    { id: 'payroll', label: 'Payroll' },
  ]

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <div className="bg-gradient-to-r from-orange-600 to-orange-800 rounded-lg p-6 text-white">
          <h1 className="text-3xl font-bold">HR & Payroll Module</h1>
          <p className="text-orange-100 mt-2">Manage employees, attendance, leave, and payroll</p>
        </div>

        <div className="flex gap-2 border-b border-gray-200">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-3 font-medium border-b-2 transition ${
                activeTab === tab.id
                  ? 'border-orange-600 text-orange-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">HR Module - {tabs.find((t) => t.id === activeTab)?.label}</p>
            <p className="text-gray-400 mt-2">Detailed implementation coming soon...</p>
          </div>
        </div>
      </div>
    </DashboardLayout>
  )
}
