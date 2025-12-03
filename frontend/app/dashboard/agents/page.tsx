'use client'

import { useState } from 'react'
import { SectionCard } from '@/components/ui/section-card'
import { StatCard } from '@/components/ui/stat-card'
import Link from 'next/link'

interface Agent {
  id: string
  name: string
  email: string
  phone: string
  status: 'online' | 'offline' | 'busy'
  callsToday: number
  avgDuration: number
  successRate: number
  joinDate: string
}

export default function AgentsPage() {
  const [agents] = useState<Agent[]>([
    { id: '1', name: 'John Doe', email: 'john.doe@company.com', phone: '+1-555-0101', status: 'online', callsToday: 24, avgDuration: 8.5, successRate: 92, joinDate: '2025-01-15' },
    { id: '2', name: 'Jane Smith', email: 'jane.smith@company.com', phone: '+1-555-0102', status: 'online', callsToday: 31, avgDuration: 9.2, successRate: 95, joinDate: '2025-02-20' },
    { id: '3', name: 'Mike Johnson', email: 'mike.johnson@company.com', phone: '+1-555-0103', status: 'busy', callsToday: 18, avgDuration: 7.8, successRate: 88, joinDate: '2025-03-10' },
    { id: '4', name: 'Sarah Williams', email: 'sarah.williams@company.com', phone: '+1-555-0104', status: 'offline', callsToday: 0, avgDuration: 0, successRate: 91, joinDate: '2025-04-05' },
    { id: '5', name: 'Tom Brown', email: 'tom.brown@company.com', phone: '+1-555-0105', status: 'online', callsToday: 27, avgDuration: 8.1, successRate: 89, joinDate: '2025-05-12' },
  ])

  const stats = [
    { label: 'Total Agents', value: agents.length, icon: '👥' },
    { label: 'Online Now', value: agents.filter(a => a.status === 'online').length, icon: '🟢' },
    { label: 'Total Calls Today', value: agents.reduce((sum, a) => sum + a.callsToday, 0), icon: '☎️' },
    { label: 'Avg Success Rate', value: `${(agents.reduce((sum, a) => sum + a.successRate, 0) / agents.length).toFixed(1)}%`, icon: '✅' },
  ]

  const statusColors = {
    'online': 'bg-green-100 text-green-800',
    'offline': 'bg-gray-100 text-gray-800',
    'busy': 'bg-yellow-100 text-yellow-800',
  }

  const statusIcon = {
    'online': '🟢',
    'offline': '⚫',
    'busy': '🟡',
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-800 rounded-lg p-6 text-white">
        <h1 className="text-3xl font-bold">Call Center Agents</h1>
        <p className="text-blue-100 mt-2">Manage and monitor your call center team</p>
      </div>

      {/* KPIs */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, i) => (
          <StatCard key={i} label={stat.label} value={stat.value} icon={stat.icon} />
        ))}
      </div>

      {/* Agents Table */}
      <SectionCard title="All Agents" action={<Link href="/dashboard/agents/new"><button className="px-3 py-1 bg-blue-600 text-white text-sm rounded hover:bg-blue-700">+ Add Agent</button></Link>}>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Name</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Email</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Phone</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Status</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Calls Today</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Avg Duration</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Success Rate</th>
                <th className="px-4 py-2 text-left font-medium text-gray-700">Actions</th>
              </tr>
            </thead>
            <tbody>
              {agents.map((agent) => (
                <tr key={agent.id} className="border-b hover:bg-gray-50">
                  <td className="px-4 py-3 font-medium text-gray-900">{agent.name}</td>
                  <td className="px-4 py-3 text-gray-600">{agent.email}</td>
                  <td className="px-4 py-3 text-gray-600">{agent.phone}</td>
                  <td className="px-4 py-3">
                    <span className={`text-xs font-medium px-2 py-1 rounded ${statusColors[agent.status]}`}>
                      {statusIcon[agent.status]} {agent.status.charAt(0).toUpperCase() + agent.status.slice(1)}
                    </span>
                  </td>
                  <td className="px-4 py-3 font-medium text-gray-900">{agent.callsToday}</td>
                  <td className="px-4 py-3 text-gray-600">{agent.avgDuration}m</td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2">
                      <div className="w-16 bg-gray-200 rounded-full h-2">
                        <div className="bg-green-600 h-2 rounded-full" style={{ width: `${agent.successRate}%` }} />
                      </div>
                      <span className="font-medium text-gray-900">{agent.successRate}%</span>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex gap-2">
                      <Link href={`/dashboard/agents/${agent.id}`}><button className="text-blue-600 hover:text-blue-800 font-medium text-xs">View</button></Link>
                      <button className="text-orange-600 hover:text-orange-800 font-medium text-xs">Edit</button>
                    </div>
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
