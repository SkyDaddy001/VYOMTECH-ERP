'use client'

import React, { useEffect } from 'react'
import { useGamification } from '@/hooks/useGamification'
import GamificationProfile from './GamificationProfile'
import Badges from './Badges'
import Challenges from './Challenges'
import Leaderboard from './Leaderboard'

export default function GamificationDashboard() {
  const {
    profile,
    badges,
    challenges,
    leaderboard,
    loading,
    error,
    fetchProfile,
    fetchBadges,
    fetchChallenges,
    fetchActiveChallenges,
    fetchLeaderboard,
  } = useGamification()

  useEffect(() => {
    fetchProfile()
    fetchBadges()
    fetchActiveChallenges()
    fetchLeaderboard()
  }, [])

  if (loading) {
    return <div>Loading gamification data...</div>
  }

  if (error) {
    return <div className="text-red-600">Error: {error}</div>
  }

  return (
    <div className="space-y-8">
      {profile && <GamificationProfile level={profile.currentLevel?.currentLevel || 0} points={profile.userPoints?.currentPoints || 0} streak={profile.userPoints?.streakDays || 0} />}
      <Badges badges={badges} />
      <Challenges challenges={challenges} />
      <Leaderboard leaderboard={leaderboard} />
    </div>
  )
}
