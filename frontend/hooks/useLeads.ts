import { useState, useCallback } from 'react'
import { leadService } from '@/services/api'

export function useLeads() {
  const [leads, setLeads] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [stats, setStats] = useState<any>(null)

  const fetchLeads = useCallback(async (page?: number, limit?: number) => {
    setLoading(true)
    setError(null)
    try {
      const data = await leadService.listLeads(page, limit)
      setLeads(data as any[])
    } catch (err: any) {
      setError(err.message || 'Failed to fetch leads')
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchLead = useCallback(async (id: number | string) => {
    setLoading(true)
    setError(null)
    try {
      const data = await leadService.getLead(id)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch lead')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const createLead = useCallback(async (leadData: any) => {
    setLoading(true)
    setError(null)
    try {
      const data = await leadService.createLead(leadData)
      setLeads((prev) => [...prev, data])
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to create lead')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const updateLead = useCallback(async (id: number | string, leadData: any) => {
    setLoading(true)
    setError(null)
    try {
      const data = await leadService.updateLead(id, leadData)
      setLeads((prev) => prev.map((lead: any) => (lead.id === id ? data : lead)))
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to update lead')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const deleteLead = useCallback(async (id: number | string) => {
    setLoading(true)
    setError(null)
    try {
      await leadService.deleteLead(id)
      setLeads((prev) => prev.filter((lead: any) => lead.id !== id))
      return true
    } catch (err: any) {
      setError(err.message || 'Failed to delete lead')
      return false
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchStats = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await leadService.getLeadStats()
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
    leads,
    loading,
    error,
    stats,
    fetchLeads,
    fetchLead,
    createLead,
    updateLead,
    deleteLead,
    fetchStats,
  }
}
