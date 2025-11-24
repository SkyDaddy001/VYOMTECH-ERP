import { useState, useCallback, useEffect } from 'react'
import { analyticsService } from '@/services/api'

export function useAnalytics() {
  const [reports, setReports] = useState<any[]>([])
  const [trends, setTrends] = useState<any[]>([])
  const [metrics, setMetrics] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const generateReport = useCallback(
    async (type: string, startDate: string, endDate: string) => {
      setLoading(true)
      setError(null)
      try {
        const data = await analyticsService.generateReport(type, startDate, endDate)
        setReports((prev) => [...prev, data])
        return data
      } catch (err: any) {
        setError(err.message || 'Failed to generate report')
        return null
      } finally {
        setLoading(false)
      }
    },
    []
  )

  const exportReport = useCallback(
    async (reportId: string, format: 'csv' | 'json' | 'pdf' = 'json') => {
      setLoading(true)
      setError(null)
      try {
        const data = await analyticsService.exportReport(reportId, format)
        return data
      } catch (err: any) {
        setError(err.message || 'Failed to export report')
        return null
      } finally {
        setLoading(false)
      }
    },
    []
  )

  const fetchTrends = useCallback(
    async (metric: string, startDate: string, endDate: string) => {
      setLoading(true)
      setError(null)
      try {
        const data = await analyticsService.getTrends(metric, startDate, endDate)
        setTrends(data as any[])
        return data
      } catch (err: any) {
        setError(err.message || 'Failed to fetch trends')
        return null
      } finally {
        setLoading(false)
      }
    },
    []
  )

  const fetchMetrics = useCallback(
    async (metric: string, filters?: Record<string, any>) => {
      setLoading(true)
      setError(null)
      try {
        const data = await analyticsService.getCustomMetrics(metric, filters)
        setMetrics(data)
        return data
      } catch (err: any) {
        setError(err.message || 'Failed to fetch metrics')
        return null
      } finally {
        setLoading(false)
      }
    },
    []
  )

  return {
    reports,
    trends,
    metrics,
    loading,
    error,
    generateReport,
    exportReport,
    fetchTrends,
    fetchMetrics,
  }
}
