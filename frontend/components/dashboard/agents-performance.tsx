'use client'

import { useAgents } from '@/hooks/use-dashboard'
import { FiPhone, FiAward } from 'react-icons/fi'

export const AgentsPerformance = () => {
  const { agents, loading } = useAgents()

  if (loading) {
    return (
      <div className="bg-white rounded-sm p-6">
        <h3 className="text-sm font-semibold uppercase tracking-wide mb-6">Top Agents</h3>
        <div className="space-y-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="bg-gray-100 h-12 rounded-sm animate-pulse"></div>
          ))}
        </div>
      </div>
    )
  }

  if (!agents || agents.length === 0) {
    return (
      <div className="bg-white rounded-sm p-6">
        <h3 className="text-sm font-semibold uppercase tracking-wide mb-6">Top Agents</h3>
        <p className="text-gray-500 text-center py-8 text-sm">No agents found</p>
      </div>
    )
  }

  const sortedAgents = [...agents].sort((a: any, b: any) => b.rating - a.rating)

  return (
    <div className="bg-white rounded-sm p-6">
      <h3 className="text-sm font-semibold uppercase tracking-wide mb-6">Top Agents</h3>
      <div className="space-y-3">
        {sortedAgents.slice(0, 5).map((agent: any, index: number) => (
          <div key={agent.id} className="flex items-center justify-between p-3 border-b border-gray-100 hover:bg-gray-50 transition-colors">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 bg-gray-900 rounded-sm flex items-center justify-center text-white font-semibold text-xs">
                {index + 1}
              </div>
              <div>
                <p className="font-medium text-gray-900 text-sm">{agent.name}</p>
                <p className="text-xs text-gray-600">{agent.email}</p>
              </div>
            </div>
            <div className="flex items-center gap-4">
              <div className="text-right">
                <div className="flex items-center gap-1 text-xs">
                  <FiPhone className="text-gray-600" size={12} />
                  <span className="font-semibold text-gray-900">
                    {agent.totalCalls}
                  </span>
                </div>
                <p className="text-xs text-gray-500">calls</p>
              </div>
              <div className="text-right">
                <div className="flex items-center gap-1 text-xs">
                  <FiAward className="text-gray-600" size={12} />
                  <span className="font-semibold text-gray-900">
                    {agent.rating.toFixed(1)}
                  </span>
                </div>
                <p className="text-xs text-gray-500">rating</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
