'use client'

import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'

export default function ReportsPage() {
  return (
    <div className="flex h-screen bg-gray-50">
      <Sidebar />
      <div className="flex-1 flex flex-col lg:ml-64">
        <Header />
        <main className="flex-1 overflow-auto pt-20 pb-6">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Reports</h1>
              <p className="text-gray-600 mt-2">
                Analytics and business intelligence
              </p>
            </div>

            <div className="mt-8 bg-white rounded-lg shadow p-8">
              <p className="text-center text-gray-500">
                Reports module coming soon...
              </p>
            </div>
          </div>
        </main>
      </div>
    </div>
  )
}
