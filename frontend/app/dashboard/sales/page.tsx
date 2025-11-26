'use client'

import { useState } from 'react'
import DashboardLayout from '@/components/layouts/DashboardLayout'
import { LeadManagement } from '@/components/modules/Sales/LeadManagement'
import { CustomerManagement } from '@/components/modules/Sales/CustomerManagement'
import { QuotationManagement } from '@/components/modules/Sales/QuotationManagement'
import { SalesOrderManagement } from '@/components/modules/Sales/SalesOrderManagement'
import { InvoiceManagement } from '@/components/modules/Sales/InvoiceManagement'
import { PaymentReceipt } from '@/components/modules/Sales/PaymentReceipt'
import { MilestoneTracking } from '@/components/modules/Sales/MilestoneTracking'
import { ReportingDashboard } from '@/components/modules/Sales/ReportingDashboard'

type TabType = 'leads' | 'customers' | 'quotations' | 'orders' | 'invoices' | 'payments' | 'milestones' | 'reports'

export default function SalesPage() {
  const [activeTab, setActiveTab] = useState<TabType>('leads')

  const tabs: Array<{ id: TabType; label: string }> = [
    { id: 'leads', label: 'Leads' },
    { id: 'customers', label: 'Customers' },
    { id: 'quotations', label: 'Quotations' },
    { id: 'orders', label: 'Sales Orders' },
    { id: 'invoices', label: 'Invoices' },
    { id: 'payments', label: 'Payments' },
    { id: 'milestones', label: 'Milestones & Tracking' },
    { id: 'reports', label: 'Reports & Analytics' },
  ]

  const renderContent = () => {
    switch (activeTab) {
      case 'leads':
        return <LeadManagement />
      case 'customers':
        return <CustomerManagement />
      case 'quotations':
        return <QuotationManagement />
      case 'orders':
        return <SalesOrderManagement />
      case 'invoices':
        return <InvoiceManagement />
      case 'payments':
        return <PaymentReceipt />
      case 'milestones':
        return <MilestoneTracking />
      case 'reports':
        return <ReportingDashboard />
      default:
        return null
    }
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="bg-gradient-to-r from-green-600 to-green-800 rounded-lg p-6 text-white">
          <h1 className="text-3xl font-bold">Sales Module</h1>
          <p className="text-green-100 mt-2">Manage leads, customers, quotations, orders, invoices, payments, milestones, and reporting</p>
        </div>

        {/* Tabs Navigation */}
        <div className="flex gap-2 border-b border-gray-200 overflow-x-auto">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-3 font-medium border-b-2 transition whitespace-nowrap ${
                activeTab === tab.id
                  ? 'border-green-600 text-green-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Tab Content */}
        <div className="bg-white rounded-lg shadow p-6">
          {renderContent()}
        </div>
      </div>
    </DashboardLayout>
  )
}
