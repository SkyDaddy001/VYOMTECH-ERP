import { useState, useCallback } from 'react'
import { complianceService } from '@/services/api'

export function useCompliance() {
  const [auditLogs, setAuditLogs] = useState<any[]>([])
  const [securityEvents, setSecurityEvents] = useState<any[]>([])
  const [complianceReport, setComplianceReport] = useState<any>(null)
  const [auditSummary, setAuditSummary] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchAuditLogs = useCallback(
    async (userId?: string, action?: string, limit?: number, offset?: number) => {
      setLoading(true)
      setError(null)
      try {
        const data = await complianceService.getAuditLogs(userId, action, limit, offset)
        setAuditLogs(data as any[])
        return data
      } catch (err: any) {
        setError(err.message || 'Failed to fetch audit logs')
        return null
      } finally {
        setLoading(false)
      }
    },
    []
  )

  const fetchAuditSummary = useCallback(async (startDate?: string, endDate?: string) => {
    setLoading(true)
    setError(null)
    try {
      const data = await complianceService.getAuditSummary(startDate, endDate)
      setAuditSummary(data)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch audit summary')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchSecurityEvents = useCallback(async (status?: string, limit?: number) => {
    setLoading(true)
    setError(null)
    try {
      const data = await complianceService.getSecurityEvents(status, limit)
      setSecurityEvents(data as any[])
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch security events')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchComplianceReport = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await complianceService.getComplianceReport()
      setComplianceReport(data)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch compliance report')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const createRole = useCallback(async (name: string, description: string) => {
    setLoading(true)
    setError(null)
    try {
      const data = await complianceService.createRole(name, description)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to create role')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const getRoles = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await complianceService.getRoles()
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch roles')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const requestDataAccess = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await complianceService.requestDataAccess()
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to request data access')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const exportUserData = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await complianceService.exportUserData()
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to export user data')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const requestDataDeletion = useCallback(async (reason: string) => {
    setLoading(true)
    setError(null)
    try {
      const data = await complianceService.requestDataDeletion(reason)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to request data deletion')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const getUserConsents = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await complianceService.getUserConsents()
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch user consents')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const recordConsent = useCallback(async (type: string, consentValue: boolean) => {
    setLoading(true)
    setError(null)
    try {
      const data = await complianceService.recordConsent(type, consentValue)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to record consent')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    auditLogs,
    securityEvents,
    complianceReport,
    auditSummary,
    loading,
    error,
    fetchAuditLogs,
    fetchAuditSummary,
    fetchSecurityEvents,
    fetchComplianceReport,
    createRole,
    getRoles,
    requestDataAccess,
    exportUserData,
    requestDataDeletion,
    getUserConsents,
    recordConsent,
  }
}
