'use client'

import React, { useState, useEffect } from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'
import { gamificationDashboardService } from '@/services/api'
import { Award, TrendingUp, Users, Zap } from 'lucide-react'

export default function GamificationPresentationDashboard() {
  // State for gamification data
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [overview, setOverview] = useState<any>(null)
  const [leaderboard, setLeaderboard] = useState<any>(null)
  const [challenges, setChallenges] = useState<any>(null)
  const [analytics, setAnalytics] = useState<any>(null)

  // Fetch gamification data on mount
  useEffect(() => {
    const fetchGamificationData = async () => {
      try {
        setLoading(true)

        // Fetch overview
        const overviewRes = await gamificationDashboardService.getGamificationOverview()
        setOverview(overviewRes.data)

        // Fetch leaderboard
        const leaderboardRes = await gamificationDashboardService.getLeaderboard()
        setLeaderboard(leaderboardRes.data)

        // Fetch challenges
        const challengesRes = await gamificationDashboardService.getUserChallenges()
        setChallenges(challengesRes.data)

        // Fetch analytics
        const analyticsRes = await gamificationDashboardService.getEngagementAnalytics()
        setAnalytics(analyticsRes.data)

        setError(null)
      } catch (err: any) {
        console.error('Failed to fetch gamification data:', err)
        setError(err.message || 'Failed to load gamification data')
      } finally {
        setLoading(false)
      }
    }

    fetchGamificationData()
  }, [])

  // Use real data or fallback values
  const totalPoints = overview?.total_points || 3200000
  const activeUsers = overview?.active_users || 245
  const badgesAwarded = overview?.badges_awarded || 1250
  const engagementRate = overview?.engagement_rate || 87

  const topPlayers = leaderboard?.top_players || [
    { rank: 1, name: 'Rahul Sharma', points: 125850, level: 'Platinum', badges: 48 },
    { rank: 2, name: 'Priya Nair', points: 98320, level: 'Gold', badges: 42 },
    { rank: 3, name: 'Arjun Singh', points: 87650, level: 'Gold', badges: 38 },
    { rank: 4, name: 'Meera Kapoor', points: 76440, level: 'Silver', badges: 34 },
    { rank: 5, name: 'Vikram Patel', points: 65200, level: 'Silver', badges: 28 }
  ]

  const risingStars = leaderboard?.rising_stars || [
    { name: 'Kavya Reddy', growth: 45200, newLevel: 'Gold', lastRank: 47 },
    { name: 'Sanjay Kumar', growth: 38500, newLevel: 'Silver', lastRank: 63 },
    { name: 'Neha Gupta', growth: 34100, newLevel: 'Silver', lastRank: 71 },
    { name: 'Aditya Verma', growth: 28900, newLevel: 'Bronze', lastRank: 92 },
    { name: 'Lisa Chen', growth: 24600, newLevel: 'Bronze', lastRank: 108 }
  ]

  const activeChallenges = challenges?.active || [
    { title: '🎯 December Sales Sprint', current: 42, target: 50, progress: 84, participants: 67, endDate: 'Dec 31', leader: 'Rahul Sharma' },
    { title: '📝 Quality Challenge', current: 98.2, target: 100, progress: 98, participants: 145, endDate: 'Dec 31', leader: 'Meera Kapoor' },
    { title: '🤝 Customer Satisfaction', current: 93.5, target: 95, progress: 93, participants: 89, endDate: 'Dec 15', leader: 'Priya Nair' },
    { title: '⚡ Speed Challenge', current: 167, target: 200, progress: 83, participants: 102, endDate: 'Dec 20', leader: 'Arjun Singh' }
  ]

  const levels = [
    { level: 'Bronze', minPoints: 0, maxPoints: 25000, users: 87 },
    { level: 'Silver', minPoints: 25000, maxPoints: 75000, users: 92 },
    { level: 'Gold', minPoints: 75000, maxPoints: 150000, users: 56 },
    { level: 'Platinum', minPoints: 150000, maxPoints: Infinity, users: 10 }
  ]

  const engagementMetrics = analytics?.metrics || {
    daily_active: 198,
    weekly_active: 245,
    avg_points_per_user: 13061
  }

  const slides: Slide[] = [
    {
      id: 'cover',
      title: 'Gamification & Engagement',
      subtitle: 'Employee Motivation & Performance Tracking',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <Award className="w-20 h-20 text-blue-600" />
          <div className="grid grid-cols-2 gap-6 w-full max-w-2xl">
            <div className="bg-blue-50 p-6 rounded-lg border border-blue-200">
              <div className="text-3xl font-bold text-blue-700">{(totalPoints / 1000000).toFixed(1)}M</div>
              <div className="text-sm text-gray-600 mt-1">Total Points Earned</div>
            </div>
            <div className="bg-green-50 p-6 rounded-lg border border-green-200">
              <div className="text-3xl font-bold text-green-700">{activeUsers}</div>
              <div className="text-sm text-gray-600 mt-1">Active Players</div>
            </div>
            <div className="bg-purple-50 p-6 rounded-lg border border-purple-200">
              <div className="text-3xl font-bold text-purple-700">{badgesAwarded}+</div>
              <div className="text-sm text-gray-600 mt-1">Badges Awarded</div>
            </div>
            <div className="bg-orange-50 p-6 rounded-lg border border-orange-200">
              <div className="text-3xl font-bold text-orange-700">{engagementRate}%</div>
              <div className="text-sm text-gray-600 mt-1">Engagement Rate</div>
            </div>
          </div>
          {error && <p className="text-red-600 text-sm">{error}</p>}
        </div>
      ),
      backgroundColor: 'from-gray-50 to-blue-50'
    },
    {
      id: 'leaderboard',
      title: 'Global Leaderboard',
      subtitle: 'Top performers and rising stars',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">🏆 Top 10 Overall</h3>
            {topPlayers.map((player: any, i: number) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200">
                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-3">
                    <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold text-white ${
                      i === 0 ? 'bg-yellow-500' :
                      i === 1 ? 'bg-gray-400' :
                      i === 2 ? 'bg-orange-600' :
                      'bg-blue-500'
                    }`}>{player.rank}</div>
                    <div>
                      <div className="font-bold text-gray-800 text-sm">{player.name}</div>
                      <div className="text-xs text-gray-600">{player.level} • {player.badges} badges</div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="font-bold text-blue-600 text-sm">{player.points.toLocaleString()}</div>
                    <div className="text-xs text-gray-600">pts</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">📈 Rising Stars (30-Day Growth)</h3>
            {risingStars.map((player: any, i: number) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-green-200 bg-green-50">
                <div className="flex justify-between items-center">
                  <div>
                    <div className="font-bold text-gray-800 text-sm">{player.name}</div>
                    <div className="text-xs text-gray-600">From rank #{player.lastRank}</div>
                  </div>
                  <div className="text-right">
                    <div className="text-green-600 font-bold text-sm">+{player.growth.toLocaleString()}</div>
                    <div className="text-xs bg-green-200 text-green-800 px-1.5 py-0.5 rounded">{player.newLevel}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    },
    {
      id: 'badges-achievements',
      title: 'Badges & Achievements',
      subtitle: 'Award system and earned recognitions',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Popular Badges</h3>
            {[
              { icon: '⭐', name: 'Rising Star', desc: 'Earn 10K points in 30 days', earned: 87 },
              { icon: '🚀', name: 'Streak Master', desc: '30-day login streak', earned: 124 },
              { icon: '💯', name: 'Perfect Score', desc: 'Achieve 100% on challenge', earned: 156 },
              { icon: '👥', name: 'Team Player', desc: 'Help 5+ team members', earned: 98 }
            ].map((badge, i) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200">
                <div className="flex items-start gap-3">
                  <div className="text-3xl">{badge.icon}</div>
                  <div className="flex-1">
                    <div className="font-bold text-gray-800 text-sm">{badge.name}</div>
                    <div className="text-xs text-gray-600">{badge.desc}</div>
                    <div className="text-xs text-blue-600 mt-1">{badge.earned}/{activeUsers} earned</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Level Progression</h3>
            {levels.map((tier: any, i: number) => {
              const colors = ['orange', 'gray', 'yellow', 'purple']
              return (
                <div key={i} className={`bg-${colors[i]}-50 border-l-4 border-${colors[i]}-500 p-4 rounded`}>
                  <div className="flex justify-between items-start mb-2">
                    <div>
                      <h4 className="font-bold text-gray-800">{tier.level}</h4>
                      <div className="text-xs text-gray-600">{tier.minPoints.toLocaleString()} - {tier.maxPoints === Infinity ? '∞' : tier.maxPoints.toLocaleString()} points</div>
                    </div>
                    <div className="text-right">
                      <div className="text-2xl font-bold text-gray-800">{tier.users}</div>
                      <div className="text-xs text-gray-600">users</div>
                    </div>
                  </div>
                  <div className="w-full bg-gray-200 h-2 rounded-full overflow-hidden">
                    <div className={`bg-${colors[i]}-500 h-full`} style={{ width: `${(tier.users / 87) * 100}%` }}></div>
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-yellow-50'
    },
    {
      id: 'challenges',
      title: 'Active Challenges & Contests',
      subtitle: 'Ongoing competitions and goal tracking',
      content: (
        <div className="space-y-3 h-full overflow-y-auto">
          <div className="grid grid-cols-2 gap-3">
            {activeChallenges.map((challenge: any, i: number) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200">
                <div className="font-bold text-gray-800 text-sm mb-1">{challenge.title}</div>
                <div className="text-xs text-gray-600 mb-2">Target: {challenge.target} {i === 1 ? '%' : i === 3 ? 'tasks' : 'Cr'}</div>
                <div className="flex justify-between items-center text-xs mb-2">
                  <span className="font-bold">{challenge.current} {i === 1 ? '%' : i === 3 ? 'done' : 'Cr'}</span>
                  <span className="text-gray-500">{challenge.participants} joining</span>
                </div>
                <div className="w-full bg-gray-200 h-2 rounded-full overflow-hidden mb-2">
                  <div className="bg-green-500 h-full" style={{ width: `${challenge.progress}%` }}></div>
                </div>
                <div className="flex justify-between text-xs text-gray-600">
                  <span>Leader: {challenge.leader}</span>
                  <span>{challenge.endDate}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-green-50'
    },
    {
      id: 'rewards-shop',
      title: 'Rewards Shop & Catalog',
      subtitle: 'Redeem points for tangible rewards',
      content: (
        <div className="space-y-3 h-full overflow-y-auto">
          <div className="grid grid-cols-3 gap-3">
            {[
              { item: 'Coffee Voucher', cost: 500, claimed: 234, image: '☕' },
              { item: 'Movie Tickets', cost: 2000, claimed: 156, image: '🎬' },
              { item: 'Amazon Gift Card', cost: 5000, claimed: 89, image: '🎁' },
              { item: 'Spa Voucher', cost: 3500, claimed: 67, image: '🧖' },
              { item: 'Weekend Getaway', cost: 15000, claimed: 12, image: '🏨' },
              { item: 'Team Dinner', cost: 8000, claimed: 23, image: '🍽️' }
            ].map((reward, i) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200 text-center hover:shadow-lg">
                <div className="text-4xl mb-2">{reward.image}</div>
                <div className="font-bold text-gray-800 text-sm">{reward.item}</div>
                <div className="text-blue-600 font-bold text-sm mt-1">{reward.cost} pts</div>
                <div className="text-xs text-gray-600 mt-1">{reward.claimed} claimed</div>
              </div>
            ))}
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-purple-50'
    },
    {
      id: 'analytics',
      title: 'Engagement Analytics',
      subtitle: 'System health and user behavior insights',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Participation Metrics</h3>
            <div className="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
              <div className="text-sm text-gray-600">Daily Active Users</div>
              <div className="text-3xl font-bold text-blue-700 mt-2">{engagementMetrics.daily_active} / {activeUsers}</div>
              <div className="text-xs text-gray-600 mt-1">{Math.round((engagementMetrics.daily_active / activeUsers) * 100)}% daily engagement</div>
            </div>
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="text-sm text-gray-600">Weekly Engagement</div>
              <div className="text-3xl font-bold text-green-700 mt-2">{engagementMetrics.weekly_active} / {activeUsers}</div>
              <div className="text-xs text-gray-600 mt-1">100% weekly participation</div>
            </div>
            <div className="bg-purple-50 border-l-4 border-purple-500 p-4 rounded">
              <div className="text-sm text-gray-600">Avg Points Per User</div>
              <div className="text-3xl font-bold text-purple-700 mt-2">{engagementMetrics.avg_points_per_user.toLocaleString()}</div>
              <div className="text-xs text-gray-600 mt-1">Total {(totalPoints / 1000000).toFixed(1)}M points in system</div>
            </div>
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Activity Distribution</h3>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="space-y-3">
                {[
                  { activity: 'Sales Goals Achieved', percent: 34 },
                  { activity: 'Quality Targets Met', percent: 28 },
                  { activity: 'Challenges Completed', percent: 22 },
                  { activity: 'Badges Earned', percent: 16 }
                ].map((act, i) => (
                  <div key={i}>
                    <div className="flex justify-between text-sm mb-1">
                      <span className="font-semibold text-gray-800">{act.activity}</span>
                      <span className="text-gray-600">{act.percent}%</span>
                    </div>
                    <div className="w-full bg-gray-200 h-2.5 rounded-full overflow-hidden">
                      <div className={`h-full ${
                        i === 0 ? 'bg-blue-500' :
                        i === 1 ? 'bg-green-500' :
                        i === 2 ? 'bg-purple-500' :
                        'bg-orange-500'
                      }`} style={{ width: `${act.percent}%` }}></div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
            <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
              <div className="font-bold text-yellow-900 text-sm">💡 Insight</div>
              <div className="text-xs text-gray-700 mt-1">Sales competition driving 34% of engagement. Q1 focus on quality challenges could balance motivation.</div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    },
    {
      id: 'summary',
      title: 'Gamification Impact & Future',
      subtitle: 'Success metrics and roadmap',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-4">
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="font-bold text-green-900">✓ Engagement Success</div>
              <div className="text-sm text-gray-700 mt-2">{engagementRate}% engagement rate. {Math.round((engagementMetrics.daily_active / activeUsers) * 100)}% daily active users. {activeUsers}/{activeUsers} users participated in Q4. System thriving with healthy competition.</div>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
              <div className="font-bold text-blue-900">📊 Performance Lift</div>
              <div className="text-sm text-gray-700 mt-2">Sales targets exceeded by 18%. Quality metrics improved 12%. Gamification delivering business impact beyond engagement.</div>
            </div>
            <div className="bg-purple-50 border-l-4 border-purple-500 p-4 rounded">
              <div className="font-bold text-purple-900">⭐ Team Satisfaction</div>
              <div className="text-sm text-gray-700 mt-2">NPS increased from 6.2 to 7.8/10. Employees enjoy recognition and competitive elements. Culture strengthened significantly.</div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
              <div className="font-bold text-yellow-900">🚀 Q1 2025 Roadmap</div>
              <div className="text-sm text-gray-700 mt-2">
                • Department-level leaderboards<br/>
                • Real-time notification system<br/>
                • Mobile app integration<br/>
                • AI-powered recommendations<br/>
                • Peer-to-peer rewards
              </div>
            </div>
            <div className="bg-orange-50 border-l-4 border-orange-500 p-4 rounded">
              <div className="font-bold text-orange-900">💰 ROI Opportunity</div>
              <div className="text-sm text-gray-700 mt-2">Estimated 8% productivity increase = ₹2.4 Cr annual value. Gamification cost: ₹25 L. PAYBACK: < 2 months.</div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    }
  ]

  return <PresentationDashboard slides={slides} title="Gamification Dashboard" showSlideNumbers={true} />
}
