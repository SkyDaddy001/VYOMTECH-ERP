'use client'

import { useState } from 'react'
import { SectionCard } from '@/components/ui/section-card'
import { StatCard } from '@/components/ui/stat-card'

interface Lead {
  id: string
  name: string
  email: string
  phone: string
  company: string
  status: 'new' | 'contacted' | 'qualified' | 'converted'
  source: string
  value: number
}

export default function LeadsPage() {
  const [leads] = useState<Lead[]>([
    { id: '1', name: 'John Smith', email: 'john@company.com', phone: '555-1234', company: 'Tech Corp', status: 'qualified', source: 'Website', value: 50000 },
    { id: '2', name: 'Sarah Jones', email: 'sarah@enterprise.com', phone: '555-5678', company: 'Enterprise Ltd', status: 'contacted', source: 'LinkedIn', value: 75000 },
    { id: '3', name: 'Mike Brown', email: 'mike@startup.io', phone: '555-9012', company: 'Startup Inc', status: 'new', source: 'Referral', value: 30000 },
    { id: '4', name: 'Lisa Wong', email: 'lisa@global.com', phone: '555-3456', company: 'Global Solutions', status: 'converted', source: 'Event', value: 120000 },
  ])

  const stats = [
    { label: 'Total Leads', value: leads.length, icon: '📋', trend: 'up' as const, trendValue: '12 this month' },
    { label: 'Qualified', value: leads.filter(l => l.status === 'qualified').length, icon: '✅', trend: 'up' as const, trendValue: '3 new' },
    { label: 'Converted', value: leads.filter(l => l.status === 'converted').length, icon: '🎯', trend: 'up' as const, trendValue: '2 this week' },
    { label: 'Avg Deal Value', value: '$44.2K', icon: '💰', trend: 'up' as const, trendValue: '+8%' },
  ]

  const statusColors = {
    'new': 'bg-blue-100 text-blue-800',
    'contacted': 'bg-yellow-100 text-yellow-800',
    'qualified': 'bg-green-100 text-green-800',
    'converted': 'bg-purple-100 text-purple-800',
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-800 rounded-lg p-6 text-white">
        <h1 className="text-3xl font-bold">Leads Management</h1>
        <p className="text-blue-100 mt-2">Track and manage sales leads throughout the pipeline</p>
      </div>

      {/* KPIs */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, i) => (
          <StatCard key={i} label={stat.label} value={stat.value} icon={stat.icon} trend={stat.trend} trendValue={stat.trendValue} />
        ))}
      </div>

      {/* Leads Table */}
      <SectionCard title="Active Leads" action={<button className="px-3 py-1 bg-blue-600 text-white text-sm rounded hover:bg-blue-700">+ New Lead</button>}>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Lead Name</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Company</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Email</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Phone</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Source</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Status</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Value</th>
              </tr>
            </thead>
            <tbody>
              {leads.map((lead) => (
                <tr key={lead.id} className="border-b hover:bg-gray-50">
                  <td className="px-4 py-3 font-medium text-gray-900">{lead.name}</td>
                  <td className="px-4 py-3 text-gray-600">{lead.company}</td>
                  <td className="px-4 py-3 text-gray-600">{lead.email}</td>
                  <td className="px-4 py-3 text-gray-600">{lead.phone}</td>
                  <td className="px-4 py-3 text-gray-600">{lead.source}</td>
                  <td className="px-4 py-3">
                    <span className={`text-xs font-medium px-2 py-1 rounded ${statusColors[lead.status]}`}>
                      {lead.status.charAt(0).toUpperCase() + lead.status.slice(1)}
                    </span>
                  </td>
                  <td className="px-4 py-3 font-semibold text-gray-900">${(lead.value / 1000).toFixed(0)}K</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </SectionCard>
    </div>
  )
}
