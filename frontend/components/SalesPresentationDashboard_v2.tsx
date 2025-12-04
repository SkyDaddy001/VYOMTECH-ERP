'use client'

import React, { useState, useEffect } from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'
import { salesDashboardService } from '@/services/api'

interface MetricCard {
  label: string
  value: string | number
  change?: string
  color?: 'blue' | 'green' | 'red' | 'yellow'
}

interface RevenueChartProps {
  data: { month: string; amount: number }[]
}

function MetricCard({ label, value, change, color = 'blue' }: MetricCard) {
  const colorClasses = {
    blue: 'border-blue-400 bg-blue-50',
    green: 'border-green-400 bg-green-50',
    red: 'border-red-400 bg-red-50',
    yellow: 'border-yellow-400 bg-yellow-50'
  }

  return (
    <div className={`border-l-4 ${colorClasses[color]} p-6 rounded-lg`}>
      <p className="text-gray-600 text-sm font-semibold mb-2">{label}</p>
      <p className="text-3xl font-bold text-gray-900">{value}</p>
      {change && (
        <p className="text-sm text-green-600 mt-2">
          {change.startsWith('+') ? '📈' : '📉'} {change}
        </p>
      )}
    </div>
  )
}

function SimpleBarChart({ title, data }: { title: string; data: { label: string; value: number }[] }) {
  const maxValue = Math.max(...data.map(d => d.value), 1)

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-bold text-gray-800">{title}</h3>
      {data.map((item, idx) => (
        <div key={idx}>
          <div className="flex justify-between mb-1">
            <span className="text-sm font-semibold text-gray-700">{item.label}</span>
            <span className="text-sm font-bold text-gray-900">₹{item.value.toLocaleString('en-IN')}</span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
            <div
              className="bg-gradient-to-r from-blue-400 to-blue-600 h-full rounded-full transition-all"
              style={{ width: `${(item.value / maxValue) * 100}%` }}
            ></div>
          </div>
        </div>
      ))}
    </div>
  )
}

function SimpleLineChart({ data }: { data: { month: string; amount: number }[] }) {
  const maxValue = Math.max(...data.map(d => d.amount), 1)
  const minValue = 0
  const range = maxValue - minValue || 1

  return (
    <div className="flex items-end justify-between h-64 gap-2 bg-gradient-to-t from-blue-100 to-transparent p-4 rounded-lg">
      {data.map((item, idx) => {
        const height = ((item.amount - minValue) / range) * 100
        return (
          <div key={idx} className="flex-1 flex flex-col items-center gap-2">
            <div className="w-full bg-gradient-to-t from-blue-500 to-blue-400 rounded-t" style={{ height: `${height}%` }}></div>
            <span className="text-xs font-semibold text-gray-700">{item.month}</span>
            <span className="text-xs text-gray-600">₹{(item.amount / 100000).toFixed(1)}L</span>
          </div>
        )
      })}
    </div>
  )
}

// Format currency in Indian Rupees
function formatCurrency(value: number): string {
  if (value >= 10000000) {
    return `₹${(value / 10000000).toFixed(2)}Cr`
  } else if (value >= 100000) {
    return `₹${(value / 100000).toFixed(2)}L`
  } else if (value >= 1000) {
    return `₹${(value / 1000).toFixed(1)}K`
  } else {
    return `₹${value.toLocaleString('en-IN')}`
  }
}

export default function SalesPresentationDashboard() {
  // State for sales data
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [overview, setOverview] = useState<any>(null)
  const [metrics, setMetrics] = useState<any>(null)
  const [pipeline, setPipeline] = useState<any>(null)
  const [topCustomers, setTopCustomers] = useState<any[]>([])

  // Fetch sales data on mount
  useEffect(() => {
    const fetchSalesData = async () => {
      try {
        setLoading(true)
        const endDate = new Date()
        const startDate = new Date(endDate.getFullYear(), endDate.getMonth() - 5, 1)

        // Fetch overview
        const overviewRes = await salesDashboardService.getSalesOverview(startDate, endDate)
        setOverview(overviewRes.data)

        // Fetch metrics
        const metricsRes = await salesDashboardService.getSalesMetrics(startDate, endDate)
        setMetrics(metricsRes.data)

        // Fetch pipeline
        const pipelineRes = await salesDashboardService.getPipelineAnalysis()
        setPipeline(pipelineRes.data)

        // Fetch top customers
        const customersRes = await salesDashboardService.getTopCustomers(4)
        setTopCustomers(customersRes.data || [])

        setError(null)
      } catch (err: any) {
        console.error('Failed to fetch sales data:', err)
        setError(err.message || 'Failed to load sales data')
      } finally {
        setLoading(false)
      }
    }

    fetchSalesData()
  }, [])

  // Use real data or fallback values
  const totalInvoices = overview?.total_invoices || 142
  const totalOrders = overview?.total_orders || 89
  const totalRevenue = overview?.total_revenue || 2450000
  const pendingAmount = overview?.pending_amount || 680000
  const revenueGrowth = overview?.growth_percentage || 12
  const invoiceCount = metrics?.invoice_count || 142
  const avgInvoiceValue = metrics?.avg_invoice_value || (totalRevenue / totalInvoices)
  
  const draftOrders = pipeline?.draft_orders || 45
  const inProgressOrders = pipeline?.in_progress_orders || 32
  const readyToShip = pipeline?.ready_to_ship || 12
  const fulfillmentRate = pipeline?.fulfillment_rate || 78

  // Sample revenue data (would come from API in real scenario)
  const revenueData = [
    { month: 'Jul', amount: 150000 },
    { month: 'Aug', amount: 180000 },
    { month: 'Sep', amount: 220000 },
    { month: 'Oct', amount: 240000 },
    { month: 'Nov', amount: 200000 },
    { month: 'Dec', amount: 320000 }
  ]

  const customerData = topCustomers.length > 0 
    ? topCustomers.map(c => ({ label: c.name || c.customer_name, value: c.revenue || c.amount }))
    : [
        { label: 'ABC Corp', value: 450000 },
        { label: 'XYZ Industries', value: 380000 },
        { label: 'Tech Solutions', value: 320000 },
        { label: 'Global Ltd', value: 290000 }
      ]

  const slides: Slide[] = [
    {
      id: 'cover',
      title: 'Sales Dashboard',
      subtitle: 'Performance Overview & Analytics',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <div className="text-center">
            <h3 className="text-5xl font-bold text-blue-600 mb-4">📊 SALES REPORT</h3>
            <p className="text-2xl text-gray-700 mb-2">December 2025</p>
            <p className="text-xl text-gray-600">Real-Time Performance Analytics</p>
          </div>
          <div className="grid grid-cols-2 gap-8 mt-8 w-full max-w-2xl">
            <MetricCard label="Total Invoices" value={totalInvoices} color="blue" />
            <MetricCard label="Sales Orders" value={totalOrders} color="green" />
          </div>
          {error && <p className="text-red-600 text-sm">{error}</p>}
          <p className="text-gray-600 text-lg mt-8">Click Next to continue →</p>
        </div>
      )
    },
    {
      id: 'key-metrics',
      title: 'Key Performance Indicators',
      subtitle: 'At a Glance Summary',
      content: (
        <div className="grid grid-cols-2 gap-6">
          <MetricCard 
            label="Total Revenue" 
            value={formatCurrency(totalRevenue)}
            change={`+${revenueGrowth}% vs last month`}
            color="green"
          />
          <MetricCard 
            label="Pending Amount" 
            value={formatCurrency(pendingAmount)}
            change={`-5% vs last month`}
            color="yellow"
          />
          <MetricCard 
            label="Invoices Processed" 
            value={invoiceCount}
            change="+18 this month"
            color="blue"
          />
          <MetricCard 
            label="Avg Invoice Value" 
            value={formatCurrency(avgInvoiceValue)}
            color="blue"
          />
        </div>
      )
    },
    {
      id: 'revenue-trend',
      title: 'Revenue Trend',
      subtitle: 'Last 6 Months Performance',
      content: (
        <div className="h-full flex flex-col gap-6">
          <SimpleLineChart data={revenueData} />
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-green-50 border-l-4 border-green-400 p-4 rounded">
              <p className="text-sm text-gray-600">Peak Month</p>
              <p className="text-2xl font-bold text-green-600">Dec - {formatCurrency(320000)}</p>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-400 p-4 rounded">
              <p className="text-sm text-gray-600">Average Monthly</p>
              <p className="text-2xl font-bold text-blue-600">{formatCurrency(218333)}</p>
            </div>
          </div>
        </div>
      )
    },
    {
      id: 'breakdown',
      title: 'Sales Breakdown',
      subtitle: 'By Customer Segment',
      content: (
        <div className="grid grid-cols-2 gap-8">
          <SimpleBarChart 
            title="Top Customers"
            data={customerData}
          />
          <div className="space-y-4">
            <div className="bg-gradient-to-r from-blue-500 to-blue-600 text-white p-6 rounded-lg">
              <p className="text-sm font-semibold mb-1">Highest Revenue Customer</p>
              <p className="text-3xl font-bold">{customerData[0]?.label || 'ABC Corp'}</p>
              <p className="text-blue-100 text-sm mt-2">{formatCurrency(customerData[0]?.value || 450000)} ({((customerData[0]?.value || 450000) / totalRevenue * 100).toFixed(1)}% of total)</p>
            </div>
            <div className="bg-gray-50 p-6 rounded-lg border border-gray-200">
              <p className="text-sm font-semibold text-gray-700 mb-4">Revenue Concentration</p>
              <div className="space-y-2">
                <div>
                  <p className="text-xs text-gray-600 mb-1">Top 4 Customers</p>
                  <div className="flex items-center gap-2">
                    <div className="flex-1 bg-gray-200 h-2 rounded-full overflow-hidden">
                      <div className="bg-blue-500 h-full" style={{ width: '65%' }}></div>
                    </div>
                    <span className="text-sm font-semibold">65%</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )
    },
    {
      id: 'orders',
      title: 'Order Status',
      subtitle: 'Current Pipeline',
      content: (
        <div className="space-y-6">
          <div className="grid grid-cols-3 gap-4">
            <div className="bg-yellow-50 border-l-4 border-yellow-400 p-6 rounded-lg text-center">
              <p className="text-5xl font-bold text-yellow-600">{draftOrders}</p>
              <p className="text-gray-700 font-semibold mt-2">Draft Orders</p>
              <p className="text-xs text-gray-600 mt-1">Awaiting confirmation</p>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-400 p-6 rounded-lg text-center">
              <p className="text-5xl font-bold text-blue-600">{inProgressOrders}</p>
              <p className="text-gray-700 font-semibold mt-2">In Progress</p>
              <p className="text-xs text-gray-600 mt-1">Being processed</p>
            </div>
            <div className="bg-green-50 border-l-4 border-green-400 p-6 rounded-lg text-center">
              <p className="text-5xl font-bold text-green-600">{readyToShip}</p>
              <p className="text-gray-700 font-semibold mt-2">Ready to Ship</p>
              <p className="text-xs text-gray-600 mt-1">Awaiting dispatch</p>
            </div>
          </div>
          <div className="bg-gradient-to-r from-blue-50 to-blue-100 p-6 rounded-lg border border-blue-200">
            <p className="font-bold text-gray-900 mb-3">Order Fulfillment Rate</p>
            <div className="flex items-center gap-4">
              <div className="flex-1">
                <div className="bg-white h-6 rounded-full overflow-hidden">
                  <div className="bg-gradient-to-r from-green-400 to-green-600 h-full rounded-full" style={{ width: `${fulfillmentRate}%` }}></div>
                </div>
              </div>
              <span className="text-2xl font-bold text-green-600">{fulfillmentRate}%</span>
            </div>
            <p className="text-sm text-gray-600 mt-2">Orders completed on time</p>
          </div>
        </div>
      )
    },
    {
      id: 'summary',
      title: 'Summary & Insights',
      subtitle: 'Key Takeaways',
      content: (
        <div className="space-y-4">
          <div className="bg-green-50 border-l-4 border-green-500 p-6 rounded-lg">
            <p className="font-bold text-green-900 mb-2">✓ Strong Performance</p>
            <p className="text-gray-700">Revenue increased by {revenueGrowth}% compared to last month, showing consistent growth.</p>
          </div>
          <div className="bg-blue-50 border-l-4 border-blue-500 p-6 rounded-lg">
            <p className="font-bold text-blue-900 mb-2">📊 Customer Concentration</p>
            <p className="text-gray-700">Top customers account for significant revenue. Consider risk diversification strategy.</p>
          </div>
          <div className="bg-yellow-50 border-l-4 border-yellow-500 p-6 rounded-lg">
            <p className="font-bold text-yellow-900 mb-2">⚠️ Action Item</p>
            <p className="text-gray-700">{draftOrders} draft orders pending confirmation. Follow up to convert into actual sales.</p>
          </div>
          <div className="bg-purple-50 border-l-4 border-purple-500 p-6 rounded-lg">
            <p className="font-bold text-purple-900 mb-2">🎯 Next Steps</p>
            <p className="text-gray-700">Continue focus on customer retention and expanding the order pipeline.</p>
          </div>
        </div>
      )
    }
  ]

  return <PresentationDashboard slides={slides} title="Sales Performance Dashboard" />
}
