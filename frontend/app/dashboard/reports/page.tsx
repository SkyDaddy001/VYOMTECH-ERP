'use client'

import { useState } from 'react'
import { ProtectedRoute } from '@/hooks/use-auth'
import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { FiDownload, FiPlus, FiFilter, FiBarChart3, FiCalendar } from 'react-icons/fi'
import { format } from 'date-fns'

interface Report {
  id: string | number
  name: string
  category: string
  description: string
  frequency: string
  last_generated?: string
  file_size?: string
  format: string
  icon: string
}

const mockReports: Report[] = [
  {
    id: 1,
    name: 'Sales Pipeline Report',
    category: 'Sales',
    description: 'Detailed analysis of sales opportunities and pipeline status',
    frequency: 'Weekly',
    last_generated: new Date(Date.now() - 86400000).toISOString(),
    file_size: '2.4 MB',
    format: 'PDF',
    icon: '📊',
  },
  {
    id: 2,
    name: 'Financial Summary',
    category: 'Finance',
    description: 'Monthly financial performance and balance sheet summary',
    frequency: 'Monthly',
    last_generated: new Date(Date.now() - 604800000).toISOString(),
    file_size: '1.8 MB',
    format: 'Excel',
    icon: '💰',
  },
  {
    id: 3,
    name: 'Lead Conversion Analysis',
    category: 'CRM',
    description: 'Lead sources, conversion rates, and funnel analysis',
    frequency: 'Bi-weekly',
    last_generated: new Date(Date.now() - 1209600000).toISOString(),
    file_size: '3.2 MB',
    format: 'PDF',
    icon: '👥',
  },
  {
    id: 4,
    name: 'Agent Performance Metrics',
    category: 'Operations',
    description: 'Call center agent KPIs, productivity, and quality metrics',
    frequency: 'Daily',
    last_generated: new Date(Date.now() - 3600000).toISOString(),
    file_size: '1.2 MB',
    format: 'Excel',
    icon: '📈',
  },
  {
    id: 5,
    name: 'Property Inventory Report',
    category: 'Real Estate',
    description: 'Current property listings, occupancy, and market analysis',
    frequency: 'Weekly',
    last_generated: new Date(Date.now() - 172800000).toISOString(),
    file_size: '2.8 MB',
    format: 'PDF',
    icon: '🏠',
  },
  {
    id: 6,
    name: 'HR Analytics Dashboard',
    category: 'HR',
    description: 'Employee data, payroll summary, and compliance reports',
    frequency: 'Monthly',
    last_generated: new Date(Date.now() - 1209600000).toISOString(),
    file_size: '1.5 MB',
    format: 'Excel',
    icon: '👔',
  },
]

const categories = ['All', 'Sales', 'Finance', 'CRM', 'Operations', 'Real Estate', 'HR']

export default function ReportsPage() {
  const [filterCategory, setFilterCategory] = useState('All')

  const filteredReports = filterCategory === 'All'
    ? mockReports
    : mockReports.filter(r => r.category === filterCategory)

  return (
    <ProtectedRoute>
      <div className="flex h-screen bg-gray-50">
        <Sidebar />
        <div className="flex-1 flex flex-col lg:ml-64">
          <Header />
          <main className="flex-1 overflow-auto pt-20 pb-6">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
              {/* Header */}
              <div className="mb-8 flex items-center justify-between">
                <div>
                  <h1 className="text-3xl font-bold text-gray-900">Reports</h1>
                  <p className="text-gray-600 mt-2">View, download, and schedule reports</p>
                </div>
                <button className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg flex items-center gap-2 transition">
                  <FiPlus className="text-lg" />
                  Create Report
                </button>
              </div>

              {/* Category Filter */}
              <div className="mb-6 flex flex-wrap gap-2">
                {categories.map((cat) => (
                  <button
                    key={cat}
                    onClick={() => setFilterCategory(cat)}
                    className={`px-4 py-2 rounded-full font-medium transition ${
                      filterCategory === cat
                        ? 'bg-blue-600 text-white'
                        : 'bg-white text-gray-700 border border-gray-300 hover:border-gray-400'
                    }`}
                  >
                    {cat}
                  </button>
                ))}
              </div>

              {/* Reports Grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {filteredReports.map((report) => (
                  <div key={report.id} className="bg-white rounded-lg shadow hover:shadow-md transition p-6 flex flex-col">
                    {/* Icon and Category */}
                    <div className="flex items-start justify-between mb-4">
                      <div className="text-4xl">{report.icon}</div>
                      <span className="px-3 py-1 bg-blue-100 text-blue-800 text-xs font-semibold rounded-full">
                        {report.category}
                      </span>
                    </div>

                    {/* Title and Description */}
                    <h3 className="text-lg font-semibold text-gray-900 mb-2">{report.name}</h3>
                    <p className="text-sm text-gray-600 mb-4 flex-1">{report.description}</p>

                    {/* Metadata */}
                    <div className="space-y-2 mb-6 py-4 border-t border-gray-200">
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-gray-600">Frequency</span>
                        <span className="font-semibold text-gray-900">{report.frequency}</span>
                      </div>
                      {report.last_generated && (
                        <div className="flex items-center justify-between text-sm">
                          <span className="text-gray-600">Last Generated</span>
                          <span className="font-semibold text-gray-900">
                            {format(new Date(report.last_generated), 'MMM d, yyyy')}
                          </span>
                        </div>
                      )}
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-gray-600">Format</span>
                        <span className="px-2 py-1 bg-gray-100 text-gray-900 rounded text-xs font-semibold">
                          {report.format}
                        </span>
                      </div>
                      {report.file_size && (
                        <div className="flex items-center justify-between text-sm">
                          <span className="text-gray-600">Size</span>
                          <span className="font-semibold text-gray-900">{report.file_size}</span>
                        </div>
                      )}
                    </div>

                    {/* Action Buttons */}
                    <div className="flex gap-2 pt-4 border-t border-gray-200">
                      <button className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium text-sm transition flex items-center justify-center gap-2">
                        <FiDownload className="text-lg" />
                        Download
                      </button>
                      <button className="flex-1 px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg font-medium text-sm transition">
                        Preview
                      </button>
                    </div>
                  </div>
                ))}
              </div>

              {/* Quick Insights Section */}
              <div className="mt-12">
                <h2 className="text-2xl font-bold text-gray-900 mb-6">Quick Insights</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                  <div className="bg-white rounded-lg shadow p-6">
                    <p className="text-gray-600 text-sm font-medium">Total Reports Generated</p>
                    <p className="text-3xl font-bold text-gray-900 mt-2">24</p>
                    <p className="text-xs text-gray-500 mt-2">This month</p>
                  </div>
                  <div className="bg-white rounded-lg shadow p-6">
                    <p className="text-gray-600 text-sm font-medium">Average File Size</p>
                    <p className="text-3xl font-bold text-gray-900 mt-2">2.1 MB</p>
                    <p className="text-xs text-gray-500 mt-2">Across all reports</p>
                  </div>
                  <div className="bg-white rounded-lg shadow p-6">
                    <p className="text-gray-600 text-sm font-medium">Most Popular Category</p>
                    <p className="text-3xl font-bold text-gray-900 mt-2">Sales</p>
                    <p className="text-xs text-gray-500 mt-2">8 reports generated</p>
                  </div>
                  <div className="bg-white rounded-lg shadow p-6">
                    <p className="text-gray-600 text-sm font-medium">Scheduled Reports</p>
                    <p className="text-3xl font-bold text-gray-900 mt-2">12</p>
                    <p className="text-xs text-gray-500 mt-2">Auto-generated</p>
                  </div>
                </div>
              </div>
            </div>
          </main>
        </div>
      </div>
    </ProtectedRoute>
  )
}
