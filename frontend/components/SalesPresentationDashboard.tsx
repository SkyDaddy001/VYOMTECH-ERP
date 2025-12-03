'use client'

import React, { useState } from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'

interface MetricCard {
  label: string
  value: string | number
  change?: string
  color?: 'blue' | 'green' | 'red' | 'yellow'
}

interface SalesData {
  invoices: number
  orders: number
  revenue: number
  pending: number
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
  const maxValue = Math.max(...data.map(d => d.value))

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-bold text-gray-800">{title}</h3>
      {data.map((item, idx) => (
        <div key={idx}>
          <div className="flex justify-between mb-1">
            <span className="text-sm font-semibold text-gray-700">{item.label}</span>
            <span className="text-sm font-bold text-gray-900">₹{item.value.toLocaleString()}</span>
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
  const maxValue = Math.max(...data.map(d => d.amount))
  const minValue = 0
  const range = maxValue - minValue

  return (
    <div className="flex items-end justify-between h-64 gap-2 bg-gradient-to-t from-blue-100 to-transparent p-4 rounded-lg">
      {data.map((item, idx) => {
        const height = ((item.amount - minValue) / range) * 100
        return (
          <div key={idx} className="flex-1 flex flex-col items-center gap-2">
            <div className="w-full bg-gradient-to-t from-blue-500 to-blue-400 rounded-t" style={{ height: `${height}%` }}></div>
            <span className="text-xs font-semibold text-gray-700">{item.month}</span>
            <span className="text-xs text-gray-600">₹{(item.amount / 1000).toFixed(0)}K</span>
          </div>
        )
      })}
    </div>
  )
}

const salesData: SalesData = {
  invoices: 142,
  orders: 89,
  revenue: 2450000,
  pending: 680000
}

const revenueData = [
  { month: 'Jan', amount: 150000 },
  { month: 'Feb', amount: 180000 },
  { month: 'Mar', amount: 220000 },
  { month: 'Apr', amount: 240000 },
  { month: 'May', amount: 200000 },
  { month: 'Jun', amount: 320000 }
]

const topCustomers = [
  { name: 'ABC Corp', revenue: 450000 },
  { name: 'XYZ Industries', revenue: 380000 },
  { name: 'Tech Solutions', revenue: 320000 },
  { name: 'Global Ltd', revenue: 290000 }
]

export default function SalesPresentationDashboard() {
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
            <p className="text-xl text-gray-600">Multi-Tenant ERP System</p>
          </div>
          <div className="grid grid-cols-2 gap-8 mt-8 w-full max-w-2xl">
            <MetricCard label="Total Invoices" value={salesData.invoices} color="blue" />
            <MetricCard label="Sales Orders" value={salesData.orders} color="green" />
          </div>
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
            value={`₹${(salesData.revenue / 100000).toFixed(1)}L`}
            change="+12% vs last month"
            color="green"
          />
          <MetricCard 
            label="Pending Amount" 
            value={`₹${(salesData.pending / 100000).toFixed(1)}L`}
            change="-5% vs last month"
            color="yellow"
          />
          <MetricCard 
            label="Invoices Processed" 
            value={salesData.invoices}
            change="+18 this month"
            color="blue"
          />
          <MetricCard 
            label="Avg Invoice Value" 
            value={`₹${(salesData.revenue / salesData.invoices / 10000).toFixed(1)}K`}
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
              <p className="text-2xl font-bold text-green-600">June - ₹3.2L</p>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-400 p-4 rounded">
              <p className="text-sm text-gray-600">Average Monthly</p>
              <p className="text-2xl font-bold text-blue-600">₹2.18L</p>
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
            data={topCustomers.map(c => ({ label: c.name, value: c.revenue }))}
          />
          <div className="space-y-4">
            <div className="bg-gradient-to-r from-blue-500 to-blue-600 text-white p-6 rounded-lg">
              <p className="text-sm font-semibold mb-1">Highest Revenue Customer</p>
              <p className="text-3xl font-bold">ABC Corp</p>
              <p className="text-blue-100 text-sm mt-2">₹4.5L (18.4% of total)</p>
            </div>
            <div className="bg-gray-50 p-6 rounded-lg border border-gray-200">
              <p className="text-sm font-semibold text-gray-700 mb-4">Growth Rate</p>
              <div className="space-y-2">
                <div>
                  <p className="text-xs text-gray-600 mb-1">Q1 2025</p>
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
              <p className="text-5xl font-bold text-yellow-600">45</p>
              <p className="text-gray-700 font-semibold mt-2">Draft Orders</p>
              <p className="text-xs text-gray-600 mt-1">Awaiting confirmation</p>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-400 p-6 rounded-lg text-center">
              <p className="text-5xl font-bold text-blue-600">32</p>
              <p className="text-gray-700 font-semibold mt-2">In Progress</p>
              <p className="text-xs text-gray-600 mt-1">Being processed</p>
            </div>
            <div className="bg-green-50 border-l-4 border-green-400 p-6 rounded-lg text-center">
              <p className="text-5xl font-bold text-green-600">12</p>
              <p className="text-gray-700 font-semibold mt-2">Ready to Ship</p>
              <p className="text-xs text-gray-600 mt-1">Awaiting dispatch</p>
            </div>
          </div>
          <div className="bg-gradient-to-r from-blue-50 to-blue-100 p-6 rounded-lg border border-blue-200">
            <p className="font-bold text-gray-900 mb-3">Order Fulfillment Rate</p>
            <div className="flex items-center gap-4">
              <div className="flex-1">
                <div className="bg-white h-6 rounded-full overflow-hidden">
                  <div className="bg-gradient-to-r from-green-400 to-green-600 h-full rounded-full" style={{ width: '78%' }}></div>
                </div>
              </div>
              <span className="text-2xl font-bold text-green-600">78%</span>
            </div>
            <p className="text-sm text-gray-600 mt-2">89 of 114 orders completed on time</p>
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
            <p className="text-gray-700">Revenue increased by 12% compared to last month, showing consistent growth.</p>
          </div>
          <div className="bg-blue-50 border-l-4 border-blue-500 p-6 rounded-lg">
            <p className="font-bold text-blue-900 mb-2">📊 Customer Concentration</p>
            <p className="text-gray-700">Top 4 customers account for 65% of revenue. Consider risk diversification strategy.</p>
          </div>
          <div className="bg-yellow-50 border-l-4 border-yellow-500 p-6 rounded-lg">
            <p className="font-bold text-yellow-900 mb-2">⚠️ Action Item</p>
            <p className="text-gray-700">45 draft orders pending confirmation. Follow up to convert into actual sales.</p>
          </div>
          <div className="bg-purple-50 border-l-4 border-purple-500 p-6 rounded-lg">
            <p className="font-bold text-purple-900 mb-2">🎯 Next Steps</p>
            <p className="text-gray-700">Launch Q2 2025 promotional campaign targeting new customer segments.</p>
          </div>
        </div>
      )
    }
  ]

  return <PresentationDashboard slides={slides} title="Sales Performance Dashboard" />
}
