'use client'

import { useState, useEffect, useCallback } from 'react'
import { api } from '@/lib/api'

interface UseDataOptions {
  limit?: number
  offset?: number
  autoFetch?: boolean
  [key: string]: any
}

interface UseDataState<T> {
  data: T[]
  total: number
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
}

/**
 * Generic hook for fetching paginated data from API
 */
export function useData<T>(
  endpoint: () => Promise<any>,
  options: UseDataOptions = {}
): UseDataState<T> {
  const [data, setData] = useState<T[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refetch = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const result = await endpoint()
      
      if (Array.isArray(result)) {
        setData(result as T[])
        setTotal(result.length)
      } else if (result && typeof result === 'object') {
        setData(result.data || result)
        setTotal(result.total || result.data?.length || 0)
      } else {
        setData([])
        setTotal(0)
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch data'
      setError(message)
      console.error('Data fetch error:', err)
    } finally {
      setLoading(false)
    }
  }, [endpoint])

  useEffect(() => {
    if (options.autoFetch !== false) {
      refetch()
    }
  }, [refetch, options.autoFetch])

  return { data, total, loading, error, refetch }
}

/**
 * Hook for fetching leads
 */
export function useLeads(options: UseDataOptions = {}) {
  const { limit = 50, offset = 0, ...rest } = options
  return useData(() => api.leads.list({ limit, offset, ...rest }), options)
}

/**
 * Hook for fetching agents
 */
export function useAgents(options: UseDataOptions = {}) {
  const { limit = 50, offset = 0, ...rest } = options
  return useData(() => api.agents.list({ limit, offset, ...rest }), options)
}

/**
 * Hook for fetching calls
 */
export function useCalls(options: UseDataOptions = {}) {
  const { limit = 50, offset = 0, ...rest } = options
  return useData(() => api.calls.list({ limit, offset, ...rest }), options)
}

/**
 * Hook for fetching campaigns
 */
export function useCampaigns(options: UseDataOptions = {}) {
  const { limit = 50, offset = 0, ...rest } = options
  return useData(() => api.campaigns.list({ limit, offset, ...rest }), options)
}

/**
 * Hook for fetching properties
 */
export function useProperties(options: UseDataOptions = {}) {
  const { limit = 50, offset = 0, ...rest } = options
  return useData(() => api.realEstate.listProperties({ limit, offset, ...rest }), options)
}

/**
 * Hook for fetching opportunities
 */
export function useOpportunities(options: UseDataOptions = {}) {
  const { limit = 50, offset = 0, ...rest } = options
  return useData(() => api.sales.listOpportunities({ limit, offset, ...rest }), options)
}

/**
 * Hook for dashboard stats
 */
export function useDashboardStats() {
  const [stats, setStats] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refetch = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const result = await api.dashboard.getStats()
      setStats(result)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch stats'
      setError(message)
      console.error('Stats fetch error:', err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refetch()
  }, [refetch])

  return { stats, loading, error, refetch }
}

/**
 * Hook for single item
 */
export function useItem<T>(
  endpoint: () => Promise<any>,
  autoFetch = true
) {
  const [item, setItem] = useState<T | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refetch = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const result = await endpoint()
      setItem(result as T)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch item'
      setError(message)
      console.error('Item fetch error:', err)
    } finally {
      setLoading(false)
    }
  }, [endpoint])

  useEffect(() => {
    if (autoFetch) {
      refetch()
    }
  }, [refetch, autoFetch])

  return { item, loading, error, refetch }
}
