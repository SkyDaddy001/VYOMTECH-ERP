'use client'

import { useState } from 'react'
import { SectionCard } from '@/components/ui/section-card'
import { StatCard } from '@/components/ui/stat-card'

interface Call {
  id: string
  direction: 'inbound' | 'outbound'
  agent: string
  customer: string
  duration: number
  status: 'completed' | 'missed' | 'ongoing'
  date: string
  outcome: string
}

export default function CallsPage() {
  const [calls] = useState<Call[]>([
    { id: '1', direction: 'inbound', agent: 'John Doe', customer: 'Acme Corp', duration: 420, status: 'completed', date: '2025-11-26', outcome: 'Lead Qualified' },
    { id: '2', direction: 'outbound', agent: 'Jane Smith', customer: 'Tech Solutions', duration: 180, status: 'completed', date: '2025-11-26', outcome: 'Meeting Scheduled' },
    { id: '3', direction: 'inbound', agent: 'Mike Johnson', customer: 'Global Ltd', duration: 0, status: 'missed', date: '2025-11-26', outcome: 'Callback Required' },
    { id: '4', direction: 'outbound', agent: 'Sarah Williams', customer: 'Enterprise Inc', duration: 720, status: 'completed', date: '2025-11-26', outcome: 'Contract Sent' },
    { id: '5', direction: 'inbound', agent: 'Tom Brown', customer: 'StartUp Inc', duration: 240, status: 'completed', date: '2025-11-25', outcome: 'Sale Closed' },
  ])

  const stats = [
    { label: 'Total Calls', value: calls.length, icon: '📞', trend: 'up' as const, trendValue: '23 today' },
    { label: 'Completed', value: calls.filter(c => c.status === 'completed').length, icon: '✅', trend: 'up' as const, trendValue: '4 missed' },
    { label: 'Avg Duration', value: `${Math.round(calls.filter(c => c.status === 'completed').reduce((sum, c) => sum + c.duration, 0) / calls.filter(c => c.status === 'completed').length / 60)}m`, icon: '⏱️', trend: 'down' as const, trendValue: '2m improvement' },
    { label: 'Success Rate', value: `${((calls.filter(c => c.status === 'completed').length / calls.length) * 100).toFixed(0)}%`, icon: '📊', trend: 'up' as const, trendValue: '+5%' },
  ]

  const directionIcon = (dir: string) => dir === 'inbound' ? '📥' : '📤'

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-gradient-to-r from-green-600 to-green-800 rounded-lg p-6 text-white">
        <h1 className="text-3xl font-bold">Call Management</h1>
        <p className="text-green-100 mt-2">Track inbound and outbound calls with outcomes</p>
      </div>

      {/* KPIs */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, i) => (
          <StatCard key={i} label={stat.label} value={stat.value} icon={stat.icon} trend={stat.trend} trendValue={stat.trendValue} />
        ))}
      </div>

      {/* Calls Table */}
      <SectionCard title="Recent Calls" action={<button className="px-3 py-1 bg-green-600 text-white text-sm rounded hover:bg-green-700">+ New Call</button>}>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Type</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Agent</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Customer</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Duration</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Outcome</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Status</th>
              </tr>
            </thead>
            <tbody>
              {calls.map((call) => (
                <tr key={call.id} className="border-b hover:bg-gray-50">
                  <td className="px-4 py-3 text-2xl">{directionIcon(call.direction)}</td>
                  <td className="px-4 py-3 font-medium text-gray-900">{call.agent}</td>
                  <td className="px-4 py-3 text-gray-600">{call.customer}</td>
                  <td className="px-4 py-3 text-gray-600">{call.status === 'completed' ? `${Math.floor(call.duration / 60)}m ${call.duration % 60}s` : 'N/A'}</td>
                  <td className="px-4 py-3 font-medium text-gray-900">{call.outcome}</td>
                  <td className="px-4 py-3">
                    <span className={`text-xs font-medium px-2 py-1 rounded ${
                      call.status === 'completed' ? 'bg-green-100 text-green-800' :
                      call.status === 'missed' ? 'bg-red-100 text-red-800' :
                      'bg-blue-100 text-blue-800'
                    }`}>
                      {call.status.charAt(0).toUpperCase() + call.status.slice(1)}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </SectionCard>
    </div>
  )
}
