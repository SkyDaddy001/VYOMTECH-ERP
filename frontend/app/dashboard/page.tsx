'use client'

import { useEffect, useState } from 'react'
import { useDashboardStats, useLeads, useAgents, useCalls } from '@/hooks/use-api'
import { ProtectedRoute } from '@/hooks/use-auth'
import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { format, formatDistance } from 'date-fns'
import { FiTrendingUp, FiPhone, FiUsers, FiTarget, FiMoreVertical } from 'react-icons/fi'

export default function DashboardPage() {
  const { stats, loading: statsLoading } = useDashboardStats()
  const { data: leads, loading: leadsLoading } = useLeads({ limit: 5 })
  const { data: agents, loading: agentsLoading } = useAgents({ limit: 5 })
  const { data: calls, loading: callsLoading } = useCalls({ limit: 5 })

  const loading = statsLoading || leadsLoading || agentsLoading || callsLoading

  // Stats Cards
  const statsCards = [
    {
      title: 'Active Leads',
      value: stats?.totalLeads || 0,
      icon: FiTarget,
      color: 'bg-blue-500',
      change: '+12.5%',
    },
    {
      title: 'Total Calls',
      value: stats?.totalCalls || 0,
      icon: FiPhone,
      color: 'bg-green-500',
      change: '+8.3%',
    },
    {
      title: 'Active Agents',
      value: stats?.agents || 0,
      icon: FiUsers,
      color: 'bg-purple-500',
      change: '+2.1%',
    },
    {
      title: 'Conversion Rate',
      value: `${Math.round(stats?.conversionRate || 0)}%`,
      icon: FiTrendingUp,
      color: 'bg-orange-500',
      change: '+4.2%',
    },
  ]

  return (
    <ProtectedRoute>
      <div className="flex h-screen bg-gray-50">
        <Sidebar />
        <div className="flex-1 flex flex-col lg:ml-64">
          <Header />
          <main className="flex-1 overflow-auto pt-20 pb-6">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
              {/* Page Title */}
              <div className="mb-8">
                <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
                <p className="text-gray-600 mt-2">
                  Welcome back! Here's your business performance overview.
                </p>
              </div>

              {/* Stats Cards */}
              <section className="mb-8">
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                  {statsCards.map((card, index) => {
                    const Icon = card.icon
                    return (
                      <div key={index} className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition">
                        <div className="flex items-start justify-between mb-4">
                          <div className={`${card.color} p-3 rounded-lg`}>
                            <Icon className="text-white text-xl" />
                          </div>
                          <button className="text-gray-400 hover:text-gray-600">
                            <FiMoreVertical />
                          </button>
                        </div>
                        <p className="text-gray-600 text-sm font-medium">{card.title}</p>
                        <p className="text-3xl font-bold text-gray-900 mt-2">{loading ? '—' : card.value}</p>
                        {!loading && (
                          <p className="text-green-600 text-xs mt-2 font-medium">{card.change} from last month</p>
                        )}
                      </div>
                    )
                  })}
                </div>
              </section>

              {/* Main Grid */}
              <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Left Column - Main Content */}
                <div className="lg:col-span-2 space-y-6">
                  {/* Recent Leads */}
                  <div className="bg-white rounded-lg shadow p-6">
                    <div className="flex items-center justify-between mb-4">
                      <h2 className="text-lg font-semibold text-gray-900">Recent Leads</h2>
                      <a href="/dashboard/leads" className="text-blue-600 hover:text-blue-700 text-sm font-medium">
                        View All
                      </a>
                    </div>

                    {leadsLoading ? (
                      <div className="space-y-3">
                        {[...Array(3)].map((_, i) => (
                          <div key={i} className="h-12 bg-gray-200 rounded animate-pulse"></div>
                        ))}
                      </div>
                    ) : leads.length === 0 ? (
                      <div className="text-center py-8 text-gray-500">No leads found</div>
                    ) : (
                      <div className="overflow-x-auto">
                        <table className="w-full text-sm">
                          <thead className="text-gray-600 border-b">
                            <tr>
                              <th className="text-left py-2 px-2 font-medium">Name</th>
                              <th className="text-left py-2 px-2 font-medium">Status</th>
                              <th className="text-left py-2 px-2 font-medium">Value</th>
                              <th className="text-left py-2 px-2 font-medium">Date</th>
                            </tr>
                          </thead>
                          <tbody>
                            {leads.slice(0, 5).map((lead: any) => (
                              <tr key={lead.id} className="border-b hover:bg-gray-50">
                                <td className="py-3 px-2 font-medium text-gray-900">{lead.name}</td>
                                <td className="py-3 px-2">
                                  <span className="px-2 py-1 rounded-full text-xs font-semibold bg-blue-100 text-blue-800">
                                    {lead.status || 'Unknown'}
                                  </span>
                                </td>
                                <td className="py-3 px-2 text-gray-600">
                                  {lead.value ? `$${lead.value.toLocaleString()}` : '—'}
                                </td>
                                <td className="py-3 px-2 text-gray-600 text-xs">
                                  {format(new Date(lead.created_at), 'MMM d, yyyy')}
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    )}
                  </div>

                  {/* Recent Calls */}
                  <div className="bg-white rounded-lg shadow p-6">
                    <div className="flex items-center justify-between mb-4">
                      <h2 className="text-lg font-semibold text-gray-900">Recent Calls</h2>
                      <a href="/dashboard/calls" className="text-blue-600 hover:text-blue-700 text-sm font-medium">
                        View All
                      </a>
                    </div>

                    {callsLoading ? (
                      <div className="space-y-3">
                        {[...Array(3)].map((_, i) => (
                          <div key={i} className="h-12 bg-gray-200 rounded animate-pulse"></div>
                        ))}
                      </div>
                    ) : calls.length === 0 ? (
                      <div className="text-center py-8 text-gray-500">No calls found</div>
                    ) : (
                      <div className="space-y-3">
                        {calls.slice(0, 5).map((call: any) => (
                          <div key={call.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition">
                            <div className="flex-1">
                              <p className="font-medium text-gray-900">{call.lead_name || 'Unknown'}</p>
                              <p className="text-xs text-gray-600">
                                {call.agent_name || 'Unknown Agent'} • {formatDistance(new Date(call.created_at), new Date(), { addSuffix: true })}
                              </p>
                            </div>
                            <div className="text-right">
                              <span className="inline-block px-2 py-1 rounded text-xs font-semibold bg-green-100 text-green-800">
                                {call.status || 'unknown'}
                              </span>
                              <p className="text-xs text-gray-600 mt-1">{call.duration ? Math.round(call.duration / 60) + 'm' : '—'}</p>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>

                {/* Right Column - Sidebar */}
                <div className="space-y-6">
                  {/* Active Agents */}
                  <div className="bg-white rounded-lg shadow p-6">
                    <div className="flex items-center justify-between mb-4">
                      <h2 className="text-lg font-semibold text-gray-900">Active Agents</h2>
                      <a href="/dashboard/agents" className="text-blue-600 hover:text-blue-700 text-sm font-medium">
                        View All
                      </a>
                    </div>

                    {agentsLoading ? (
                      <div className="space-y-3">
                        {[...Array(3)].map((_, i) => (
                          <div key={i} className="h-12 bg-gray-200 rounded animate-pulse"></div>
                        ))}
                      </div>
                    ) : agents.length === 0 ? (
                      <div className="text-center py-8 text-gray-500">No agents found</div>
                    ) : (
                      <div className="space-y-3">
                        {agents.slice(0, 5).map((agent: any) => (
                          <div key={agent.id} className="flex items-start gap-3 p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition">
                            <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0">
                              <span className="text-sm font-bold text-blue-600">{agent.name?.charAt(0)}</span>
                            </div>
                            <div className="flex-1 min-w-0">
                              <p className="font-medium text-gray-900 text-sm truncate">{agent.name}</p>
                              <p className="text-xs text-gray-600">
                                <span className={`inline-block px-2 py-1 rounded-full text-xs font-semibold ${
                                  agent.status === 'active'
                                    ? 'bg-green-100 text-green-800'
                                    : 'bg-gray-100 text-gray-800'
                                }`}>
                                  {agent.status || 'offline'}
                                </span>
                              </p>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* Quick Actions */}
                  <div className="bg-white rounded-lg shadow p-6">
                    <h2 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h2>
                    <div className="space-y-2">
                      <a href="/dashboard/leads?action=create" className="block w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition text-center text-sm">
                        Add Lead
                      </a>
                      <a href="/dashboard/agents?action=create" className="block w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition text-center text-sm">
                        Add Agent
                      </a>
                      <a href="/dashboard/campaigns?action=create" className="block w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition text-center text-sm">
                        Create Campaign
                      </a>
                    </div>
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
