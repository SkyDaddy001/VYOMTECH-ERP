import { useState, useCallback } from 'react'
import { callService } from '@/services/api'

export function useCalls() {
  const [calls, setCalls] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [stats, setStats] = useState<any>(null)

  const fetchCalls = useCallback(async (page?: number, limit?: number) => {
    setLoading(true)
    setError(null)
    try {
      const data = await callService.listCalls(page, limit)
      setCalls(data as any[])
    } catch (err: any) {
      setError(err.message || 'Failed to fetch calls')
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchCall = useCallback(async (id: number | string) => {
    setLoading(true)
    setError(null)
    try {
      const data = await callService.getCall(id)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch call')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const createCall = useCallback(async (callData: any) => {
    setLoading(true)
    setError(null)
    try {
      const data = await callService.createCall(callData)
      setCalls((prev) => [...prev, data])
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to create call')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const endCall = useCallback(async (id: number | string) => {
    setLoading(true)
    setError(null)
    try {
      const data = await callService.endCall(id)
      setCalls((prev) => prev.map((call: any) => (call.id === id ? { ...call, status: 'ended' } : call)))
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to end call')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchStats = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await callService.getCallStats()
      setStats(data)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch stats')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    calls,
    loading,
    error,
    stats,
    fetchCalls,
    fetchCall,
    createCall,
    endCall,
    fetchStats,
  }
}
