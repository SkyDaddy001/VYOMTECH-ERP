'use client'

import { useEffect, useState } from 'react'
import { FiBarChart2, FiDownload, FiFilter, FiPlus } from 'react-icons/fi'
import { format } from 'date-fns'

interface Report {
  id: string
  name: string
  type: 'sales' | 'leads' | 'campaigns' | 'calls' | 'revenue'
  period: string
  generatedAt: string
  size: string
  format: 'pdf' | 'csv' | 'xlsx'
}

export default function ReportsPage() {
  const [reports, setReports] = useState<Report[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchReports()
  }, [])

  const fetchReports = async () => {
    try {
      setLoading(true)
      const mockReports: Report[] = [
        {
          id: '1',
          name: 'Q1 Sales Report',
          type: 'sales',
          period: 'Jan - Mar 2024',
          generatedAt: new Date(Date.now() - 604800000).toISOString(),
          size: '2.5 MB',
          format: 'pdf'
        },
        {
          id: '2',
          name: 'Lead Generation Report',
          type: 'leads',
          period: 'Mar 2024',
          generatedAt: new Date(Date.now() - 259200000).toISOString(),
          size: '1.2 MB',
          format: 'xlsx'
        },
        {
          id: '3',
          name: 'Campaign Performance',
          type: 'campaigns',
          period: 'Last 30 Days',
          generatedAt: new Date(Date.now() - 86400000).toISOString(),
          size: '3.1 MB',
          format: 'csv'
        },
        {
          id: '4',
          name: 'Call Center Analytics',
          type: 'calls',
          period: 'Mar 2024',
          generatedAt: new Date(Date.now() - 172800000).toISOString(),
          size: '1.8 MB',
          format: 'pdf'
        },
        {
          id: '5',
          name: 'Revenue Summary',
          type: 'revenue',
          period: 'Q1 2024',
          generatedAt: new Date(Date.now() - 1209600000).toISOString(),
          size: '0.9 MB',
          format: 'xlsx'
        }
      ]
      setReports(mockReports)
    } catch (error) {
      console.error('Error fetching reports:', error)
    } finally {
      setLoading(false)
    }
  }

  const getTypeColor = (type: string) => {
    const colors: Record<string, string> = {
      sales: 'bg-blue-100 text-blue-800',
      leads: 'bg-green-100 text-green-800',
      campaigns: 'bg-purple-100 text-purple-800',
      calls: 'bg-yellow-100 text-yellow-800',
      revenue: 'bg-indigo-100 text-indigo-800'
    }
    return colors[type] || 'bg-gray-100 text-gray-800'
  }

  const getFormatIcon = (format: string) => {
    const icons: Record<string, string> = {
      pdf: '📄',
      csv: '📊',
      xlsx: '📈'
    }
    return icons[format] || '📁'
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Reports</h1>
          <p className="text-gray-600 mt-2">View and download your business reports</p>
        </div>
        <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition flex items-center">
          <FiPlus className="w-4 h-4 mr-2" />
          Generate Report
        </button>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Total Reports</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">{reports.length}</p>
            </div>
            <FiBarChart2 className="w-10 h-10 text-blue-500 opacity-20" />
          </div>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">This Month</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">
                {reports.filter(r => new Date(r.generatedAt).getMonth() === new Date().getMonth()).length}
              </p>
            </div>
            <FiBarChart2 className="w-10 h-10 text-green-500 opacity-20" />
          </div>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Total Size</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">
                {(reports.reduce((sum, r) => sum + parseFloat(r.size), 0)).toFixed(1)} MB
              </p>
            </div>
            <FiBarChart2 className="w-10 h-10 text-purple-500 opacity-20" />
          </div>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Last Generated</p>
              <p className="text-lg font-bold text-gray-900 mt-2">
                {reports.length > 0 ? format(new Date(reports[0].generatedAt), 'MMM dd') : 'N/A'}
              </p>
            </div>
            <FiBarChart2 className="w-10 h-10 text-orange-500 opacity-20" />
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-4 mb-6">
        <div className="flex items-center gap-2">
          <FiFilter className="w-4 h-4 text-gray-600" />
          <button className="px-3 py-1 text-sm border border-gray-300 rounded-lg hover:bg-gray-50">
            All Types
          </button>
          <button className="px-3 py-1 text-sm border border-gray-300 rounded-lg hover:bg-gray-50">
            PDF
          </button>
          <button className="px-3 py-1 text-sm border border-gray-300 rounded-lg hover:bg-gray-50">
            XLSX
          </button>
          <button className="px-3 py-1 text-sm border border-gray-300 rounded-lg hover:bg-gray-50">
            CSV
          </button>
        </div>
      </div>

      {/* Reports List */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        {loading ? (
          <div className="p-8 text-center">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        ) : reports.length === 0 ? (
          <div className="p-8 text-center">
            <p className="text-gray-600">No reports available</p>
          </div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Report Name</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Type</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Period</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Generated</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Size</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Format</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900"></th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {reports.map((report) => (
                <tr key={report.id} className="hover:bg-gray-50 transition">
                  <td className="px-6 py-4">
                    <p className="font-medium text-gray-900">{report.name}</p>
                  </td>
                  <td className="px-6 py-4">
                    <span className={`px-3 py-1 rounded-full text-sm font-medium ${getTypeColor(report.type)}`}>
                      {report.type.charAt(0).toUpperCase() + report.type.slice(1)}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {report.period}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {format(new Date(report.generatedAt), 'MMM dd, yyyy')}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {report.size}
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center">
                      <span className="mr-2">{getFormatIcon(report.format)}</span>
                      <p className="text-sm font-medium text-gray-900">{report.format.toUpperCase()}</p>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <button className="p-2 hover:bg-gray-100 rounded-lg transition">
                      <FiDownload className="w-4 h-4 text-gray-600" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
