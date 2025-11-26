'use client'

import { useState } from 'react'
import DashboardLayout from '@/components/layouts/DashboardLayout'
import { 
  PurchaseModule,
  DashboardStats as PurchaseDashboard,
  VendorManagement,
  PurchaseOrderManagement,
  GRNManagement
} from '@/components/modules/Purchase/PurchaseModule'

type TabType = 'dashboard' | 'vendors' | 'orders' | 'grn'

export default function PurchasePage() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard')

  const tabs: Array<{ id: TabType; label: string }> = [
    { id: 'dashboard', label: 'Dashboard' },
    { id: 'vendors', label: 'Vendors' },
    { id: 'orders', label: 'Purchase Orders' },
    { id: 'grn', label: 'GRN/MRN' },
  ]

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="bg-gradient-to-r from-blue-600 to-blue-800 rounded-lg p-6 text-white">
          <h1 className="text-3xl font-bold">Purchase Module</h1>
          <p className="text-blue-100 mt-2">Manage vendors, purchase orders, and GRN/MRN</p>
        </div>

        {/* Tabs Navigation */}
        <div className="flex gap-2 border-b border-gray-200">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-3 font-medium border-b-2 transition ${
                activeTab === tab.id
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Tab Content */}
        <div className="bg-white rounded-lg shadow">
          {activeTab === 'dashboard' && <PurchaseDashboard />}
          {activeTab === 'vendors' && <VendorManagement />}
          {activeTab === 'orders' && <PurchaseOrderManagement />}
          {activeTab === 'grn' && <GRNManagement />}
        </div>
      </div>
    </DashboardLayout>
  )
}
