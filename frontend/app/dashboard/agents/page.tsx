'use client'

import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { useAgents } from '@/hooks/use-dashboard'
import { FiPhone, FiAward, FiTrendingUp } from 'react-icons/fi'

export default function AgentsPage() {
  const { agents, loading } = useAgents()

  return (
    <div className="flex h-screen bg-gray-50">
      <Sidebar />
      <div className="flex-1 flex flex-col lg:ml-64">
        <Header />
        <main className="flex-1 overflow-auto pt-20 pb-6">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div className="mb-8">
              <h1 className="text-3xl font-bold text-gray-900">Agents</h1>
              <p className="text-gray-600 mt-2">
                Manage and monitor agent performance
              </p>
            </div>

            {loading ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {[...Array(6)].map((_, i) => (
                  <div
                    key={i}
                    className="bg-white rounded-lg shadow h-64 animate-pulse"
                  ></div>
                ))}
              </div>
            ) : !agents || agents.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-500">No agents found</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {agents.map((agent: any) => (
                  <div
                    key={agent.id}
                    className="bg-white rounded-lg shadow hover:shadow-lg transition p-6"
                  >
                    <div className="flex items-start justify-between mb-4">
                      <h3 className="text-lg font-semibold text-gray-900">
                        {agent.name}
                      </h3>
                      <span
                        className={`px-2 py-1 rounded text-xs font-semibold ${
                          agent.status === 'active'
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {agent.status}
                      </span>
                    </div>

                    <p className="text-sm text-gray-600 mb-4">{agent.email}</p>

                    <div className="space-y-3 border-t pt-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2 text-gray-600">
                          <FiPhone size={16} />
                          <span className="text-sm">Total Calls</span>
                        </div>
                        <span className="font-semibold text-gray-900">
                          {agent.totalCalls}
                        </span>
                      </div>
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2 text-gray-600">
                          <FiTrendingUp size={16} />
                          <span className="text-sm">Successful</span>
                        </div>
                        <span className="font-semibold text-green-600">
                          {agent.successfulCalls}
                        </span>
                      </div>
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2 text-gray-600">
                          <FiAward size={16} />
                          <span className="text-sm">Rating</span>
                        </div>
                        <span className="font-semibold text-yellow-600">
                          {agent.rating.toFixed(1)}/5
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </main>
      </div>
    </div>
  )
}
