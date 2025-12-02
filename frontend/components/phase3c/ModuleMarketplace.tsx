'use client'

import React, { useEffect, useState } from 'react'
import { useModuleStore } from '@/contexts/phase3cStore'
import toast from 'react-hot-toast'

export function ModuleMarketplace() {
  const { modules, subscriptions, fetchModules, subscribeToModule, toggleModule, loading, error } =
    useModuleStore()
  const [filter, setFilter] = useState<'all' | 'active' | 'beta'>('all')

  useEffect(() => {
    const status = filter === 'all' ? undefined : filter
    fetchModules(status)
  }, [filter, fetchModules])

  useEffect(() => {
    if (error) {
      toast.error(error)
    }
  }, [error])

  const isSubscribed = (moduleId: string) => {
    return subscriptions.some((sub: any) => sub.module_id === moduleId)
  }

  const handleSubscribe = async (moduleId: string) => {
    try {
      await subscribeToModule({
        module_id: moduleId,
        scope_level: 'tenant',
        scope_id: 'current', // Would be replaced with actual tenant ID
      })
      toast.success('Module subscribed successfully!')
    } catch (error: any) {
      toast.error(error.message || 'Failed to subscribe')
    }
  }

  const getPricingDisplay = (module: any) => {
    switch (module.pricing_model) {
      case 'free':
        return 'Free'
      case 'per_user':
        return `₹${module.cost_per_user}/user`
      case 'per_project':
        return `₹${module.cost_per_project}/project`
      case 'per_company':
        return `₹${module.cost_per_company}/company`
      case 'flat':
        return `₹${module.base_cost}/month`
      case 'tiered':
        return 'Tiered Pricing'
      default:
        return 'Contact Sales'
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold mb-4">Module Marketplace</h2>

        <div className="flex gap-2 mb-6">
          {(['all', 'active', 'beta'] as const).map((status) => (
            <button
              key={status}
              onClick={() => setFilter(status)}
              className={`px-4 py-2 rounded-lg transition capitalize ${
                filter === status
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              {status}
            </button>
          ))}
        </div>
      </div>

      {loading ? (
        <div className="text-center py-8">Loading modules...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {modules.map((module: any) => (
            <div key={module.id} className="border rounded-lg p-4 hover:shadow-lg transition">
              <div className="flex justify-between items-start mb-2">
                <div>
                  <h3 className="font-bold text-lg">{module.name}</h3>
                  <p className="text-gray-600 text-sm">{module.category}</p>
                </div>
                <span
                  className={`px-2 py-1 text-xs rounded font-semibold ${
                    module.status === 'active'
                      ? 'bg-green-100 text-green-800'
                      : module.status === 'beta'
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-gray-100 text-gray-800'
                  }`}
                >
                  {module.status}
                </span>
              </div>

              <p className="text-gray-600 text-sm mb-4">{module.description}</p>

              <div className="bg-gray-50 p-3 rounded mb-4">
                <p className="text-lg font-bold text-blue-600">{getPricingDisplay(module)}</p>
                {module.trial_days_allowed > 0 && (
                  <p className="text-xs text-gray-600 mt-1">{module.trial_days_allowed} days trial</p>
                )}
              </div>

              {module.max_users && (
                <p className="text-xs text-gray-600 mb-3">Max Users: {module.max_users}</p>
              )}

              <button
                onClick={() => handleSubscribe(module.id)}
                disabled={loading || isSubscribed(module.id)}
                className={`w-full py-2 rounded-lg font-semibold transition ${
                  isSubscribed(module.id)
                    ? 'bg-gray-300 text-gray-600 cursor-not-allowed'
                    : 'bg-blue-500 hover:bg-blue-600 text-white'
                }`}
              >
                {isSubscribed(module.id) ? 'Subscribed' : 'Subscribe'}
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
