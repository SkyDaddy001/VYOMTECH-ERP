'use client'

import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { StatsOverview } from '@/components/dashboard/stats-overview'
import { RecentLeads } from '@/components/dashboard/recent-leads'
import { CampaignsOverview } from '@/components/dashboard/campaigns-overview'
import { AgentsPerformance } from '@/components/dashboard/agents-performance'
import { ProtectedRoute } from '@/hooks/use-auth'

export default function DashboardPage() {
  return (
    <ProtectedRoute>
      <div className="flex h-screen bg-gray-50">
      <Sidebar />
      <div className="flex-1 flex flex-col lg:ml-64">
        <Header />
        <main className="flex-1 overflow-auto pt-20 pb-6">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            {/* Page Title */}
            <div className="mb-8">
              <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
              <p className="text-gray-600 mt-2">
                Welcome back! Here's your business performance overview.
              </p>
            </div>

            {/* Stats Overview */}
            <section className="mb-8">
              <StatsOverview />
            </section>

            {/* Main Grid */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Left Column - Large */}
              <div className="lg:col-span-2 space-y-6">
                <RecentLeads limit={8} />
                <CampaignsOverview limit={3} />
              </div>

              {/* Right Column - Sidebar */}
              <div className="space-y-6">
                <AgentsPerformance />
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
    </ProtectedRoute>
  )
}
