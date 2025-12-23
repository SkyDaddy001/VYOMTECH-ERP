'use client'

import { useState, useEffect, useCallback, useMemo } from 'react'
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
        // Handle different API response formats
        // Format 1: { data: [...], total: number }
        if (result.data) {
          setData(Array.isArray(result.data) ? result.data : [result.data])
          setTotal(result.total || (Array.isArray(result.data) ? result.data.length : 1))
        }
        // Format 2: { agents: [...], count: number } or similar
        else if (result.agents || result.leads || result.calls || result.campaigns || result.properties || result.opportunities) {
          const key = Object.keys(result).find(k => Array.isArray(result[k]))
          const items = key ? result[key] : []
          setData(items as T[])
          setTotal(result.count || result.total || items.length)
        }
        // Format 3: Direct array in response
        else if (Object.keys(result).length > 0) {
          const values = Object.values(result).find(v => Array.isArray(v))
          if (Array.isArray(values)) {
            setData(values as T[])
            setTotal(values.length)
          } else {
            setData([result] as T[])
            setTotal(1)
          }
        } else {
          setData([])
          setTotal(0)
        }
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
    refetch()
  }, [endpoint])

  return { data, total, loading, error, refetch }
}

/**
 * Hook for fetching leads
 */
export function useLeads(options: UseDataOptions = {}) {
  const limit = options.limit || 50
  const offset = options.offset || 0
  const endpoint = useCallback(() => api.leads.list({ limit, offset }), [limit, offset])
  return useData(endpoint, options)
}

/**
 * Hook for fetching agents
 */
export function useAgents(options: UseDataOptions = {}) {
  const limit = options.limit || 50
  const offset = options.offset || 0
  const endpoint = useCallback(() => api.agents.list({ limit, offset }), [limit, offset])
  return useData(endpoint, options)
}

/**
 * Hook for fetching calls
 */
export function useCalls(options: UseDataOptions = {}) {
  const limit = options.limit || 50
  const offset = options.offset || 0
  const endpoint = useCallback(() => api.calls.list({ limit, offset }), [limit, offset])
  return useData(endpoint, options)
}

/**
 * Hook for fetching campaigns
 */
export function useCampaigns(options: UseDataOptions = {}) {
  const limit = options.limit || 50
  const offset = options.offset || 0
  const endpoint = useCallback(() => api.campaigns.list({ limit, offset }), [limit, offset])
  return useData(endpoint, options)
}

/**
 * Hook for fetching properties
 */
export function useProperties(options: UseDataOptions = {}) {
  const limit = options.limit || 50
  const offset = options.offset || 0
  const endpoint = useCallback(() => api.realEstate.listProperties({ limit, offset }), [limit, offset])
  return useData(endpoint, options)
}

/**
 * Hook for fetching opportunities
 */
export function useOpportunities(options: UseDataOptions = {}) {
  const limit = options.limit || 50
  const offset = options.offset || 0
  const endpoint = useCallback(() => api.sales.listOpportunities({ limit, offset }), [limit, offset])
  return useData(endpoint, options)
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
