import { useState, useCallback } from 'react'
import { automationService } from '@/services/api'

export function useAutomation() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [metrics, setMetrics] = useState<any>(null)
  const [rankedLeads, setRankedLeads] = useState<any[]>([])

  const calculateLeadScore = useCallback(async (leadId: number | string) => {
    setLoading(true)
    setError(null)
    try {
      const data = await automationService.calculateLeadScore(leadId)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to calculate lead score')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const rankLeads = useCallback(async (limit?: number) => {
    setLoading(true)
    setError(null)
    try {
      const data = await automationService.rankLeads(limit)
      setRankedLeads(data as any[])
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to rank leads')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const routeLeadToAgent = useCallback(async (leadId: number | string) => {
    setLoading(true)
    setError(null)
    try {
      const data = await automationService.routeLeadToAgent(leadId)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to route lead to agent')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const createRoutingRule = useCallback(async (ruleData: any) => {
    setLoading(true)
    setError(null)
    try {
      const data = await automationService.createRoutingRule(ruleData)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to create routing rule')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const scheduleCampaign = useCallback(
    async (campaignId: number | string, scheduledTime: string) => {
      setLoading(true)
      setError(null)
      try {
        const data = await automationService.scheduleCampaign(campaignId, scheduledTime)
        return data
      } catch (err: any) {
        setError(err.message || 'Failed to schedule campaign')
        return null
      } finally {
        setLoading(false)
      }
    },
    []
  )

  const fetchMetrics = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await automationService.getLeadScoringMetrics()
      setMetrics(data)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch metrics')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    loading,
    error,
    metrics,
    rankedLeads,
    calculateLeadScore,
    rankLeads,
    routeLeadToAgent,
    createRoutingRule,
    scheduleCampaign,
    fetchMetrics,
  }
}
