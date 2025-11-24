import { useState, useCallback } from 'react'
import { gamificationService } from '@/services/api'

export function useGamification() {
  const [profile, setProfile] = useState<any>(null)
  const [badges, setBadges] = useState<any[]>([])
  const [challenges, setChallenges] = useState<any[]>([])
  const [leaderboard, setLeaderboard] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchProfile = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await gamificationService.getGamificationProfile()
      setProfile(data)
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch gamification profile')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchBadges = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await gamificationService.getUserBadges()
      setBadges(data as any[])
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch badges')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const awardBadge = useCallback(
    async (userId: number | string, badgeId: number | string) => {
      setLoading(true)
      setError(null)
      try {
        const data = await gamificationService.awardBadge(userId, badgeId)
        return data
      } catch (err: any) {
        setError(err.message || 'Failed to award badge')
        return null
      } finally {
        setLoading(false)
      }
    },
    []
  )

  const fetchChallenges = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await gamificationService.getUserChallenges()
      setChallenges(data as any[])
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch challenges')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchActiveChallenges = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await gamificationService.getActiveChallenges()
      setChallenges(data as any[])
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch active challenges')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const fetchLeaderboard = useCallback(async (limit?: number) => {
    setLoading(true)
    setError(null)
    try {
      const data = await gamificationService.getLeaderboard(limit)
      setLeaderboard(data as any[])
      return data
    } catch (err: any) {
      setError(err.message || 'Failed to fetch leaderboard')
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  const awardPoints = useCallback(
    async (userId: number | string, points: number, reason: string) => {
      setLoading(true)
      setError(null)
      try {
        const data = await gamificationService.awardPoints(userId, points, reason)
        return data
      } catch (err: any) {
        setError(err.message || 'Failed to award points')
        return null
      } finally {
        setLoading(false)
      }
    },
    []
  )

  return {
    profile,
    badges,
    challenges,
    leaderboard,
    loading,
    error,
    fetchProfile,
    fetchBadges,
    awardBadge,
    fetchChallenges,
    fetchActiveChallenges,
    fetchLeaderboard,
    awardPoints,
  }
}
