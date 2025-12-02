'use client'

import React, { useState } from 'react'
import PropertyManagement from '@/components/modules/RealEstate/PropertyManagement'
import CustomerBookingTracker from '@/components/modules/RealEstate/CustomerBookingTracker'
import MilestoneAndPaymentTracking from '@/components/modules/RealEstate/MilestoneAndPaymentTracking'

type TabType = 'properties' | 'bookings' | 'milestones'

interface Tab {
  id: TabType
  label: string
}

export default function RealEstatePage() {
  const [activeTab, setActiveTab] = useState<TabType>('properties')

  const tabs: Tab[] = [
    { id: 'properties', label: 'Property Management' },
    { id: 'bookings', label: 'Customer Bookings' },
    { id: 'milestones', label: 'Milestones & Payments' }
  ]

  const renderContent = () => {
    switch (activeTab) {
      case 'properties':
        return <PropertyManagement />
      case 'bookings':
        return <CustomerBookingTracker />
      case 'milestones':
        return <MilestoneAndPaymentTracking />
      default:
        return null
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-4xl font-bold text-gray-900">Real Estate Management System</h1>
        <p className="text-gray-600 mt-2">Track properties, bookings, milestones, and payments in one unified platform</p>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <div className="flex gap-8">
          {tabs.map(tab => (
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
      </div>

      {/* Tab Content */}
      <div className="bg-white rounded-lg">
        {renderContent()}
      </div>
    </div>
  )
}
