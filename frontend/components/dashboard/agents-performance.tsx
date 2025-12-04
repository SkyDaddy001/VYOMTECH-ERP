'use client'

import { useAgents } from '@/hooks/use-dashboard'
import { FiPhone, FiAward } from 'react-icons/fi'

export const AgentsPerformance = () => {
  const { agents, loading } = useAgents()

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold mb-4">Top Agents</h3>
        <div className="space-y-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="bg-gray-200 h-12 rounded animate-pulse"></div>
          ))}
        </div>
      </div>
    )
  }

  if (!agents || agents.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold mb-4">Top Agents</h3>
        <p className="text-gray-500 text-center py-8">No agents found</p>
      </div>
    )
  }

  const sortedAgents = [...agents].sort((a: any, b: any) => b.rating - a.rating)

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h3 className="text-lg font-semibold mb-4">Top Agents</h3>
      <div className="space-y-4">
        {sortedAgents.slice(0, 5).map((agent: any, index: number) => (
          <div key={agent.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-gradient-to-br from-blue-400 to-blue-600 rounded-full flex items-center justify-center text-white font-semibold text-sm">
                {index + 1}
              </div>
              <div>
                <p className="font-semibold text-gray-900">{agent.name}</p>
                <p className="text-xs text-gray-600">{agent.email}</p>
              </div>
            </div>
            <div className="flex items-center gap-4">
              <div className="text-right">
                <div className="flex items-center gap-1 text-sm">
                  <FiPhone className="text-blue-600" size={14} />
                  <span className="font-semibold text-gray-900">
                    {agent.totalCalls}
                  </span>
                </div>
                <p className="text-xs text-gray-600">calls</p>
              </div>
              <div className="text-right">
                <div className="flex items-center gap-1 text-sm">
                  <FiAward className="text-yellow-600" size={14} />
                  <span className="font-semibold text-gray-900">
                    {agent.rating.toFixed(1)}
                  </span>
                </div>
                <p className="text-xs text-gray-600">rating</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
