'use client'

import React, { useState, useEffect } from 'react'
import { BarChart, LineChart, PieChart, TrendingUp, Users, BookOpen, DollarSign } from 'lucide-react'

interface FunnelData {
  month: string
  leads_new: number
  leads_contacted: number
  leads_qualified: number
  leads_negotiation: number
  leads_converted: number
  leads_lost: number
  total_leads: number
  conversion_rate: number
}

interface SourceData {
  source_type: string
  total_leads: number
  converted_leads: number
  conversion_rate: number
  lost_leads: number
  qualified_leads: number
}

interface BookingData {
  status: string
  booking_count: number
  total_booking_amount: number
  total_units_booked: number
  avg_booking_amount: number
}

interface DashboardMetrics {
  total_leads: number
  new_leads: number
  qualified_leads: number
  converted_leads: number
  conversion_rate: number
  active_customers: number
  total_bookings: number
  booked_amount: number
  outstanding_balance: number
  engagement_this_month: number
  pending_follow_ups: number
}

export function ReportingDashboard() {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null)
  const [funnelData, setFunnelData] = useState<FunnelData[]>([])
  const [sourceData, setSourceData] = useState<SourceData[]>([])
  const [bookingData, setBookingData] = useState<BookingData[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedReport, setSelectedReport] = useState<'dashboard' | 'funnel' | 'source' | 'bookings'>('dashboard')

  const tenantId = localStorage.getItem('tenantId') || ''

  useEffect(() => {
    fetchDashboardMetrics()
    fetchReports()
  }, [])

  const fetchDashboardMetrics = async () => {
    try {
      const response = await fetch('/api/v1/sales/reports/dashboard', {
        headers: {
          'X-Tenant-ID': tenantId,
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      })
      if (response.ok) {
        const data = await response.json()
        setMetrics(data)
      }
    } catch (error) {
      console.error('Failed to fetch metrics:', error)
    }
  }

  const fetchReports = async () => {
    setLoading(true)
    try {
      const [funnelRes, sourceRes, bookingRes] = await Promise.all([
        fetch('/api/v1/sales/reports/funnel', {
          headers: {
            'X-Tenant-ID': tenantId,
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
        }),
        fetch('/api/v1/sales/reports/source-performance', {
          headers: {
            'X-Tenant-ID': tenantId,
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
        }),
        fetch('/api/v1/sales/reports/bookings', {
          headers: {
            'X-Tenant-ID': tenantId,
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
        }),
      ])

      if (funnelRes.ok) {
        const data = await funnelRes.json()
        setFunnelData(data || [])
      }
      if (sourceRes.ok) {
        const data = await sourceRes.json()
        setSourceData(data || [])
      }
      if (bookingRes.ok) {
        const data = await bookingRes.json()
        setBookingData(data || [])
      }
    } catch (error) {
      console.error('Failed to fetch reports:', error)
    } finally {
      setLoading(false)
    }
  }

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: 'INR',
      maximumFractionDigits: 0,
    }).format(value)
  }

  return (
    <div className="space-y-6">
      {/* Tabs */}
      <div className="flex gap-2 border-b border-gray-200 overflow-x-auto">
        {[
          { id: 'dashboard', label: 'Dashboard', icon: <TrendingUp className="w-4 h-4" /> },
          { id: 'funnel', label: 'Lead Funnel', icon: <BarChart className="w-4 h-4" /> },
          { id: 'source', label: 'Source Performance', icon: <PieChart className="w-4 h-4" /> },
          { id: 'bookings', label: 'Bookings', icon: <BookOpen className="w-4 h-4" /> },
        ].map((tab) => (
          <button
            key={tab.id}
            onClick={() => setSelectedReport(tab.id as any)}
            className={`flex items-center gap-2 px-4 py-2 font-medium transition-colors whitespace-nowrap ${
              selectedReport === tab.id
                ? 'text-blue-600 border-b-2 border-blue-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            {tab.icon}
            {tab.label}
          </button>
        ))}
      </div>

      {/* Dashboard Overview */}
      {selectedReport === 'dashboard' && metrics && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {/* Total Leads Card */}
            <div className="bg-white p-6 rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600 mb-1">Total Leads</p>
                  <p className="text-3xl font-bold text-gray-900">{metrics.total_leads}</p>
                  <p className="text-xs text-gray-500 mt-2">
                    New: {metrics.new_leads} | Qualified: {metrics.qualified_leads}
                  </p>
                </div>
                <Users className="w-12 h-12 text-blue-500 opacity-20" />
              </div>
            </div>

            {/* Conversion Rate Card */}
            <div className="bg-white p-6 rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600 mb-1">Conversion Rate</p>
                  <p className="text-3xl font-bold text-gray-900">{metrics.conversion_rate.toFixed(1)}%</p>
                  <p className="text-xs text-gray-500 mt-2">
                    Converted: {metrics.converted_leads} leads
                  </p>
                </div>
                <TrendingUp className="w-12 h-12 text-green-500 opacity-20" />
              </div>
            </div>

            {/* Active Customers Card */}
            <div className="bg-white p-6 rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600 mb-1">Active Customers</p>
                  <p className="text-3xl font-bold text-gray-900">{metrics.active_customers}</p>
                  <p className="text-xs text-gray-500 mt-2">
                    Bookings: {metrics.total_bookings}
                  </p>
                </div>
                <Users className="w-12 h-12 text-purple-500 opacity-20" />
              </div>
            </div>

            {/* Outstanding Balance Card */}
            <div className="bg-white p-6 rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600 mb-1">Outstanding Balance</p>
                  <p className="text-2xl font-bold text-gray-900">{formatCurrency(metrics.outstanding_balance)}</p>
                  <p className="text-xs text-gray-500 mt-2">
                    Booked: {formatCurrency(metrics.booked_amount)}
                  </p>
                </div>
                <DollarSign className="w-12 h-12 text-orange-500 opacity-20" />
              </div>
            </div>

            {/* Engagement This Month Card */}
            <div className="bg-white p-6 rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition">
              <div>
                <p className="text-sm text-gray-600 mb-1">Engagement This Month</p>
                <p className="text-3xl font-bold text-gray-900">{metrics.engagement_this_month}</p>
              </div>
            </div>

            {/* Pending Follow-ups Card */}
            <div className="bg-white p-6 rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition">
              <div>
                <p className="text-sm text-gray-600 mb-1">Pending Follow-ups</p>
                <p className="text-3xl font-bold text-red-600">{metrics.pending_follow_ups}</p>
              </div>
            </div>
          </div>

          {/* Quick Stats */}
          <div className="bg-white p-6 rounded-lg border border-gray-200">
            <h3 className="font-semibold mb-4">Key Metrics Summary</h3>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              <div className="space-y-1">
                <p className="text-sm text-gray-600">Lead to Customer Ratio</p>
                <p className="text-lg font-bold">
                  {metrics.total_leads > 0 ? (metrics.converted_leads / metrics.total_leads * 100).toFixed(1) : 0}%
                </p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-gray-600">Avg Booking Value</p>
                <p className="text-lg font-bold">
                  {formatCurrency(metrics.total_bookings > 0 ? metrics.booked_amount / metrics.total_bookings : 0)}
                </p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-gray-600">Lead Quality</p>
                <p className="text-lg font-bold">
                  {metrics.total_leads > 0 ? (metrics.qualified_leads / metrics.total_leads * 100).toFixed(1) : 0}%
                </p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Lead Funnel */}
      {selectedReport === 'funnel' && (
        <div className="space-y-4">
          <div className="bg-white p-6 rounded-lg border border-gray-200">
            <h3 className="font-semibold mb-4">Lead Funnel Analysis (Last 12 Months)</h3>
            {loading ? (
              <p className="text-gray-500">Loading data...</p>
            ) : funnelData.length === 0 ? (
              <p className="text-gray-500">No data available</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead className="border-b">
                    <tr className="text-left">
                      <th className="pb-2 font-medium text-gray-700">Month</th>
                      <th className="pb-2 font-medium text-gray-700">New</th>
                      <th className="pb-2 font-medium text-gray-700">Contacted</th>
                      <th className="pb-2 font-medium text-gray-700">Qualified</th>
                      <th className="pb-2 font-medium text-gray-700">Negotiation</th>
                      <th className="pb-2 font-medium text-gray-700">Converted</th>
                      <th className="pb-2 font-medium text-gray-700">Lost</th>
                      <th className="pb-2 font-medium text-gray-700">Conv. Rate</th>
                    </tr>
                  </thead>
                  <tbody>
                    {funnelData.map((row) => (
                      <tr key={row.month} className="border-b hover:bg-gray-50">
                        <td className="py-2">
                          {new Date(row.month).toLocaleDateString('en-US', { month: 'short', year: 'numeric' })}
                        </td>
                        <td className="py-2">{row.leads_new}</td>
                        <td className="py-2">{row.leads_contacted}</td>
                        <td className="py-2">{row.leads_qualified}</td>
                        <td className="py-2">{row.leads_negotiation}</td>
                        <td className="py-2 font-bold text-green-600">{row.leads_converted}</td>
                        <td className="py-2 text-red-600">{row.leads_lost}</td>
                        <td className="py-2 font-bold">{row.conversion_rate}%</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Source Performance */}
      {selectedReport === 'source' && (
        <div className="space-y-4">
          <div className="bg-white p-6 rounded-lg border border-gray-200">
            <h3 className="font-semibold mb-4">Lead Source Performance</h3>
            {loading ? (
              <p className="text-gray-500">Loading data...</p>
            ) : sourceData.length === 0 ? (
              <p className="text-gray-500">No data available</p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {sourceData.map((source) => (
                  <div key={source.source_type} className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <h4 className="font-medium text-gray-900 capitalize mb-3">{source.source_type}</h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-600">Total Leads:</span>
                        <span className="font-bold">{source.total_leads}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Converted:</span>
                        <span className="font-bold text-green-600">{source.converted_leads}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Conversion Rate:</span>
                        <span className="font-bold">{source.conversion_rate}%</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Lost:</span>
                        <span className="font-bold text-red-600">{source.lost_leads}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Qualified:</span>
                        <span className="font-bold text-blue-600">{source.qualified_leads}</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}

      {/* Booking Summary */}
      {selectedReport === 'bookings' && (
        <div className="space-y-4">
          <div className="bg-white p-6 rounded-lg border border-gray-200">
            <h3 className="font-semibold mb-4">Booking Status Summary</h3>
            {loading ? (
              <p className="text-gray-500">Loading data...</p>
            ) : bookingData.length === 0 ? (
              <p className="text-gray-500">No booking data available</p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {bookingData.map((booking) => (
                  <div key={booking.status} className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <h4 className="font-medium text-gray-900 capitalize mb-3">{booking.status}</h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-600">Count:</span>
                        <span className="font-bold">{booking.booking_count}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Total Amount:</span>
                        <span className="font-bold">{formatCurrency(booking.total_booking_amount)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Units Booked:</span>
                        <span className="font-bold">{booking.total_units_booked}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Avg Amount:</span>
                        <span className="font-bold">{formatCurrency(booking.avg_booking_amount)}</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
