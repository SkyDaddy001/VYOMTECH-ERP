'use client'

import { useState } from 'react'
import { Trophy, Star, Zap, TrendingUp } from 'lucide-react'
import type { UserAchievement, Leaderboard, Challenge, Badge, GamificationDashboard } from '@/types/gamification'

type TabType = 'dashboard' | 'leaderboard' | 'challenges' | 'badges' | 'achievements'

export default function GamificationPage() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard')

  // Mock data
  const dashboard: GamificationDashboard = {
    total_users_participated: 156,
    total_points_distributed: 125680,
    active_challenges: 5,
    completed_challenges: 23,
    total_badges: 42,
    top_performer_name: 'Alex Johnson',
    top_performer_points: 8540,
    average_user_points: 805,
    engagement_rate_percentage: 78,
  }

  const leaderboard: Leaderboard[] = [
    {
      id: '1',
      user_id: 'U001',
      rank: 1,
      user_name: 'Alex Johnson',
      total_points: 8540,
      achievements_count: 18,
      badges_count: 12,
      current_streak: 45,
      last_activity: '2024-11-29 14:30',
    },
    {
      id: '2',
      user_id: 'U002',
      rank: 2,
      user_name: 'Sarah Williams',
      total_points: 7820,
      achievements_count: 16,
      badges_count: 11,
      current_streak: 38,
      last_activity: '2024-11-29 13:45',
    },
    {
      id: '3',
      user_id: 'U003',
      rank: 3,
      user_name: 'Michael Chen',
      total_points: 7340,
      achievements_count: 15,
      badges_count: 10,
      current_streak: 32,
      last_activity: '2024-11-29 15:20',
    },
    {
      id: '4',
      user_id: 'U004',
      rank: 4,
      user_name: 'Emma Davis',
      total_points: 6890,
      achievements_count: 13,
      badges_count: 9,
      current_streak: 28,
      last_activity: '2024-11-29 12:15',
    },
    {
      id: '5',
      user_id: 'U005',
      rank: 5,
      user_name: 'James Wilson',
      total_points: 6320,
      achievements_count: 11,
      badges_count: 8,
      current_streak: 20,
      last_activity: '2024-11-29 11:50',
    },
  ]

  const challenges: Challenge[] = [
    {
      id: '1',
      challenge_name: 'Call Master',
      description: 'Complete 100 calls with 90%+ customer satisfaction',
      start_date: '2024-11-20',
      end_date: '2024-12-04',
      challenge_type: 'individual',
      objective: '100 calls',
      reward_points: 500,
      participation_count: 42,
      status: 'active',
    },
    {
      id: '2',
      challenge_name: 'Team Synergy',
      description: 'Achieve 95% team target completion',
      start_date: '2024-11-25',
      end_date: '2024-12-09',
      challenge_type: 'team',
      objective: 'Target Completion',
      reward_points: 1000,
      participation_count: 156,
      status: 'active',
    },
    {
      id: '3',
      challenge_name: 'Week Champion',
      description: 'Top performer in your department for the week',
      start_date: '2024-11-25',
      end_date: '2024-12-01',
      challenge_type: 'department',
      objective: 'Highest Score',
      reward_points: 300,
      participation_count: 87,
      status: 'active',
    },
  ]

  const badges: Badge[] = [
    {
      id: '1',
      badge_name: 'First Steps',
      description: 'Complete your first call',
      icon_url: '🎯',
      requirement: '1 call completed',
      points_value: 10,
      rarity: 'common',
      total_earned: 156,
    },
    {
      id: '2',
      badge_name: 'Century Club',
      description: 'Complete 100 calls',
      icon_url: '💯',
      requirement: '100 calls',
      points_value: 100,
      rarity: 'uncommon',
      total_earned: 78,
    },
    {
      id: '3',
      badge_name: 'Perfect Score',
      description: 'Achieve 100% customer satisfaction',
      icon_url: '⭐',
      requirement: '100% satisfaction',
      points_value: 250,
      rarity: 'rare',
      total_earned: 23,
    },
    {
      id: '4',
      badge_name: 'Legendary Agent',
      description: 'Reach 50000 points',
      icon_url: '👑',
      requirement: '50000 points',
      points_value: 1000,
      rarity: 'legendary',
      total_earned: 2,
    },
  ]

  const achievements: UserAchievement[] = [
    {
      id: '1',
      user_id: 'U001',
      user_name: 'Alex Johnson',
      achievement_id: 'A001',
      achievement_name: 'Call Warrior - 500 calls completed',
      achievement_date: '2024-11-28',
      points_earned: 500,
      badge_unlocked: true,
      badge_name: 'Call Warrior',
    },
    {
      id: '2',
      user_id: 'U002',
      user_name: 'Sarah Williams',
      achievement_id: 'A002',
      achievement_name: '7-Day Streak',
      achievement_date: '2024-11-29',
      points_earned: 100,
      badge_unlocked: true,
      badge_name: 'On Fire',
    },
  ]

  const tabs: Array<{ id: TabType; label: string }> = [
    { id: 'dashboard', label: 'Dashboard' },
    { id: 'leaderboard', label: 'Leaderboard' },
    { id: 'challenges', label: 'Challenges' },
    { id: 'badges', label: 'Badges' },
    { id: 'achievements', label: 'Achievements' },
  ]

  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-purple-600 to-purple-800 rounded-lg p-6 text-white">
        <h1 className="text-3xl font-bold">🎮 Gamification Center</h1>
        <p className="text-purple-100 mt-2">Earn points, unlock badges, compete on leaderboards and win rewards</p>
      </div>

      <div className="flex gap-2 border-b border-gray-200 overflow-x-auto">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`px-4 py-3 font-medium border-b-2 transition whitespace-nowrap ${
              activeTab === tab.id
                ? 'border-purple-600 text-purple-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'dashboard' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Active Challenges</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.active_challenges}</p>
                </div>
                <Zap className="text-yellow-600" size={32} />
              </div>
              <p className="text-gray-600 text-sm mt-2">Ongoing competitions</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Total Badges</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.total_badges}</p>
                </div>
                <Star className="text-blue-600" size={32} />
              </div>
              <p className="text-gray-600 text-sm mt-2">Collectible</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Engagement Rate</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.engagement_rate_percentage}%</p>
                </div>
                <TrendingUp className="text-green-600" size={32} />
              </div>
              <p className="text-green-600 text-sm mt-2">Excellent</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Avg User Points</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.average_user_points}</p>
                </div>
                <Trophy className="text-purple-600" size={32} />
              </div>
              <p className="text-gray-600 text-sm mt-2">Per user</p>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">🏆 Top Performer</h3>
            <div className="flex items-center gap-4">
              <div className="w-16 h-16 rounded-full bg-gradient-to-r from-yellow-400 to-yellow-600 flex items-center justify-center text-white text-2xl font-bold">
                1
              </div>
              <div>
                <p className="text-xl font-semibold text-gray-900">{dashboard.top_performer_name}</p>
                <p className="text-gray-600">{dashboard.top_performer_points} points · {dashboard.total_users_participated} competitors</p>
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'leaderboard' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Rank</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Name</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Points</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Streak</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Badges</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Achievements</th>
              </tr>
            </thead>
            <tbody>
              {leaderboard.map((entry) => (
                <tr key={entry.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <span className="text-lg font-bold text-gray-900">
                      {entry.rank === 1 ? '🥇' : entry.rank === 2 ? '🥈' : entry.rank === 3 ? '🥉' : entry.rank}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{entry.user_name}</td>
                  <td className="px-6 py-4 text-sm text-gray-900 font-semibold">{entry.total_points}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">🔥 {entry.current_streak}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">⭐ {entry.badges_count}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">✓ {entry.achievements_count}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeTab === 'challenges' && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {challenges.map((challenge) => (
            <div key={challenge.id} className="bg-white rounded-lg shadow p-6">
              <div className="flex justify-between items-start mb-4">
                <h3 className="text-lg font-semibold text-gray-900">{challenge.challenge_name}</h3>
                <span className="px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                  {challenge.status}
                </span>
              </div>
              <p className="text-gray-600 text-sm mb-4">{challenge.description}</p>
              <div className="space-y-2 mb-4">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Type:</span>
                  <span className="font-medium text-gray-900">{challenge.challenge_type}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Objective:</span>
                  <span className="font-medium text-gray-900">{challenge.objective}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Reward:</span>
                  <span className="font-semibold text-yellow-600">{challenge.reward_points} pts</span>
                </div>
              </div>
              <div className="pt-4 border-t">
                <p className="text-sm text-gray-600">{challenge.participation_count} participants</p>
              </div>
            </div>
          ))}
        </div>
      )}

      {activeTab === 'badges' && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {badges.map((badge) => (
            <div key={badge.id} className="bg-white rounded-lg shadow p-6 text-center">
              <p className="text-6xl mb-2">{badge.icon_url}</p>
              <h3 className="font-semibold text-gray-900 mb-2">{badge.badge_name}</h3>
              <p className="text-xs text-gray-600 mb-3">{badge.description}</p>
              <div className="space-y-1 mb-3">
                <p className="text-xs text-gray-600">{badge.requirement}</p>
                <p className={`text-xs font-semibold ${
                  badge.rarity === 'common' ? 'text-gray-600' :
                  badge.rarity === 'uncommon' ? 'text-green-600' :
                  badge.rarity === 'rare' ? 'text-blue-600' :
                  badge.rarity === 'epic' ? 'text-purple-600' :
                  'text-yellow-600'
                }`}>
                  {badge.rarity.toUpperCase()}
                </p>
              </div>
              <p className="text-sm text-gray-600 border-t pt-2">Earned by {badge.total_earned}</p>
            </div>
          ))}
        </div>
      )}

      {activeTab === 'achievements' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">User</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Achievement</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Points</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Date</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Badge</th>
              </tr>
            </thead>
            <tbody>
              {achievements.map((achievement) => (
                <tr key={achievement.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{achievement.user_name}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{achievement.achievement_name}</td>
                  <td className="px-6 py-4 text-sm font-semibold text-yellow-600">+{achievement.points_earned}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{achievement.achievement_date}</td>
                  <td className="px-6 py-4">
                    <span className="px-3 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                      {achievement.badge_name || 'N/A'}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

