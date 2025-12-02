'use client'

import { useState, useEffect } from 'react'
import { StatCard } from '@/components/ui/stat-card'
import { CourseCard } from '@/components/ui/course-card'
import { SectionCard } from '@/components/ui/section-card'

interface KPIMetric {
  title: string
  value: string | number
  icon: string
  trend?: 'up' | 'down'
  trendValue?: string
}

export default function DashboardContent() {
  const [kpis] = useState<KPIMetric[]>([
    {
      title: 'Total Revenue',
      value: '₹2.45M',
      icon: '💰',
      trend: 'up',
      trendValue: '8.5% from last month',
    },
    {
      title: 'Active Employees',
      value: 245,
      icon: '👥',
      trend: 'up',
      trendValue: '12 new hires',
    },
    {
      title: 'Pending Invoices',
      value: 34,
      icon: '📄',
      trend: 'down',
      trendValue: '6 less than yesterday',
    },
    {
      title: 'Open POs',
      value: 127,
      icon: '📋',
      trend: 'up',
      trendValue: '3% increase',
    },
  ])

  return (
    <div className="space-y-6 md:space-y-8">
      {/* KPIs Overview */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {kpis.map((kpi, index) => (
          <StatCard
            key={index}
            label={kpi.title}
            value={kpi.value}
            icon={kpi.icon}
            trend={kpi.trend}
            trendValue={kpi.trendValue}
          />
        ))}
      </div>

      {/* Executive Summary */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Department Performance */}
        <SectionCard title="Department Performance" action={<a href="#" className="text-sm text-blue-600 hover:text-blue-700">View All</a>}>
          <div className="space-y-3">
            {[
              { dept: 'Sales', target: 2400000, actual: 2450000, status: 'On Track' },
              { dept: 'HR', target: 850000, actual: 820000, status: 'Below Target' },
              { dept: 'Finance', target: 450000, actual: 475000, status: 'Exceeded' },
              { dept: 'Operations', target: 620000, actual: 605000, status: 'On Track' },
            ].map((dept, i) => (
              <div key={i} className="px-3 py-2 hover:bg-gray-50 rounded transition border-l-4 border-gray-200">
                <div className="flex items-center justify-between mb-1">
                  <p className="text-sm font-medium text-gray-900">{dept.dept}</p>
                  <span className={`text-xs font-medium px-2 py-1 rounded ${
                    dept.status === 'On Track' ? 'bg-green-100 text-green-800' :
                    dept.status === 'Exceeded' ? 'bg-blue-100 text-blue-800' :
                    'bg-yellow-100 text-yellow-800'
                  }`}>
                    {dept.status}
                  </span>
                </div>
                <div className="flex items-center justify-between text-xs text-gray-600">
                  <span>₹{(dept.actual / 1000).toFixed(0)}K / ₹{(dept.target / 1000).toFixed(0)}K</span>
                  <span>{((dept.actual / dept.target) * 100).toFixed(1)}%</span>
                </div>
              </div>
            ))}
          </div>
        </SectionCard>

        {/* Pending Approvals */}
        <SectionCard title="Pending Approvals" action={<a href="#" className="text-sm text-blue-600 hover:text-blue-700">Manage</a>}>
          <div className="space-y-2">
            {[
              { type: 'Leave Request', count: 5, priority: 'high' },
              { type: 'Purchase Orders', count: 12, priority: 'medium' },
              { type: 'Expense Claims', count: 8, priority: 'medium' },
              { type: 'Salary Revisions', count: 3, priority: 'high' },
            ].map((item, i) => (
              <div
                key={i}
                className="flex items-center justify-between px-3 py-2 hover:bg-gray-50 rounded transition"
              >
                <div className="flex-1">
                  <p className="text-sm font-medium text-gray-900">{item.type}</p>
                </div>
                <div className="flex items-center gap-2">
                  <span className={`text-xs font-semibold px-2 py-1 rounded ${
                    item.priority === 'high' ? 'bg-red-100 text-red-800' : 'bg-orange-100 text-orange-800'
                  }`}>
                    {item.count}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </SectionCard>
      </div>

      {/* Financial Summary & HR Metrics */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Transactions */}
        <SectionCard title="Recent Transactions" action={<a href="#" className="text-sm text-blue-600 hover:text-blue-700">View All</a>}>
          <div className="space-y-2">
            {[
              { ref: 'INV-2025-001', type: 'Invoice', amount: '₹45,250', status: 'Paid' },
              { ref: 'PO-2025-127', type: 'Purchase Order', amount: '₹12,500', status: 'Pending' },
              { ref: 'EXP-2025-034', type: 'Expense', amount: '₹2,350', status: 'Approved' },
              { ref: 'SAL-2025-245', type: 'Salary', amount: '₹125,000', status: 'Processed' },
            ].map((tx, i) => (
              <div
                key={i}
                className="flex items-center justify-between px-3 py-2 hover:bg-gray-50 rounded transition"
              >
                <div className="flex-1">
                  <p className="text-sm font-medium text-gray-900">{tx.ref}</p>
                  <p className="text-xs text-gray-500">{tx.type}</p>
                </div>
                <div className="text-right">
                  <p className="text-sm font-semibold text-gray-900">{tx.amount}</p>
                  <span className={`text-xs font-medium px-2 py-1 rounded inline-block mt-1 ${
                    tx.status === 'Paid' ? 'bg-green-100 text-green-800' :
                    tx.status === 'Processed' ? 'bg-blue-100 text-blue-800' :
                    tx.status === 'Approved' ? 'bg-green-50 text-green-700' :
                    'bg-yellow-100 text-yellow-800'
                  }`}>
                    {tx.status}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </SectionCard>

        {/* HR Analytics */}
        <SectionCard title="HR Analytics" action={<a href="#" className="text-sm text-blue-600 hover:text-blue-700">Details</a>}>
          <div className="space-y-3">
            {[
              { metric: 'Attendance Rate', value: '94.2%', color: 'bg-green-100 text-green-800' },
              { metric: 'Avg Salary', value: '₹4,200', color: 'bg-blue-100 text-blue-800' },
              { metric: 'Open Positions', value: '8', color: 'bg-orange-100 text-orange-800' },
              { metric: 'Employee Turnover', value: '2.1%', color: 'bg-purple-100 text-purple-800' },
            ].map((item, i) => (
              <div key={i} className="flex items-center justify-between px-3 py-2 border-l-4 border-gray-200 hover:border-gray-300 transition">
                <span className="text-sm font-medium text-gray-900">{item.metric}</span>
                <span className={`text-sm font-bold px-3 py-1 rounded ${item.color}`}>{item.value}</span>
              </div>
            ))}
          </div>
        </SectionCard>
      </div>

      {/* Quick Actions */}
      <SectionCard title="Quick Actions">
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-3">
          {[
            { icon: '📊', label: 'Reports' },
            { icon: '📝', label: 'Approvals' },
            { icon: '💼', label: 'Projects' },
            { icon: '👤', label: 'Employees' },
            { icon: '📦', label: 'Purchase' },
            { icon: '💳', label: 'Finance' },
          ].map((action, i) => (
            <button
              key={i}
              className="flex flex-col items-center justify-center gap-2 px-3 py-3 rounded-lg border border-gray-200 text-gray-900 hover:bg-gray-50 transition"
            >
              <span className="text-2xl">{action.icon}</span>
              <span className="text-xs font-medium text-center">{action.label}</span>
            </button>
          ))}
        </div>
      </SectionCard>
    </div>
  )
}
