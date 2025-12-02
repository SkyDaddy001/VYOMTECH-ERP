'use client'

import { useState } from 'react'

type TabType = 'dashboard' | 'accounts' | 'entries' | 'reports' | 'reconciliation'

export default function AccountsPage() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard')

  const tabs: Array<{ id: TabType; label: string }> = [
    { id: 'dashboard', label: 'Dashboard' },
    { id: 'accounts', label: 'Chart of Accounts' },
    { id: 'entries', label: 'Journal Entries' },
    { id: 'reports', label: 'Financial Reports' },
    { id: 'reconciliation', label: 'Reconciliation' },
  ]

  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-indigo-600 to-indigo-800 rounded-lg p-6 text-white">
        <h1 className="text-3xl font-bold">Accounts (GL) Module</h1>
        <p className="text-indigo-100 mt-2">Manage chart of accounts, journal entries, and financial reports</p>
      </div>

      <div className="flex gap-2 border-b border-gray-200">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`px-4 py-3 font-medium border-b-2 transition ${
              activeTab === tab.id
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <div className="text-center py-12">
          <p className="text-gray-500 text-lg">Accounts Module - {tabs.find((t) => t.id === activeTab)?.label}</p>
          <p className="text-gray-400 mt-2">Detailed implementation coming soon...</p>
        </div>
      </div>
    </div>
  )
}
