'use client'

import { Breadcrumbs } from '@/components/ui/Breadcrumbs'
import DashboardContent from '@/components/dashboard/DashboardContent'
import { SiteNavigation } from '@/components/navigation/SiteNavigation'
import { useState } from 'react'

export default function DashboardPage() {
  const [showNav, setShowNav] = useState(false)

  return (
    <div className="space-y-6">
      <Breadcrumbs />
      
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <button
          onClick={() => setShowNav(!showNav)}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
        >
          {showNav ? '📖 Hide Map' : '📖 View Map'}
        </button>
      </div>

      {showNav && <SiteNavigation />}
      
      <DashboardContent />
    </div>
  )
}
