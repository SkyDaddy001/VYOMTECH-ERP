'use client'

import React from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'
import { Award, TrendingUp, Users, Zap } from 'lucide-react'

export default function GamificationPresentationDashboard() {
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
              <div className="text-3xl font-bold text-blue-700">3.2M</div>
              <div className="text-sm text-gray-600 mt-1">Total Points Earned</div>
            </div>
            <div className="bg-green-50 p-6 rounded-lg border border-green-200">
              <div className="text-3xl font-bold text-green-700">245</div>
              <div className="text-sm text-gray-600 mt-1">Active Players</div>
            </div>
            <div className="bg-purple-50 p-6 rounded-lg border border-purple-200">
              <div className="text-3xl font-bold text-purple-700">1,250+</div>
              <div className="text-sm text-gray-600 mt-1">Badges Awarded</div>
            </div>
            <div className="bg-orange-50 p-6 rounded-lg border border-orange-200">
              <div className="text-3xl font-bold text-orange-700">87%</div>
              <div className="text-sm text-gray-600 mt-1">Engagement Rate</div>
            </div>
          </div>
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
            {[
              { rank: 1, name: 'Rahul Sharma', points: '125,850', level: 'Platinum', badges: 48 },
              { rank: 2, name: 'Priya Nair', points: '98,320', level: 'Gold', badges: 42 },
              { rank: 3, name: 'Arjun Singh', points: '87,650', level: 'Gold', badges: 38 },
              { rank: 4, name: 'Meera Kapoor', points: '76,440', level: 'Silver', badges: 34 },
              { rank: 5, name: 'Vikram Patel', points: '65,200', level: 'Silver', badges: 28 }
            ].map((player, i) => (
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
                    <div className="font-bold text-blue-600 text-sm">{player.points}</div>
                    <div className="text-xs text-gray-600">pts</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">📈 Rising Stars (30-Day Growth)</h3>
            {[
              { name: 'Kavya Reddy', growth: '+45,200 pts', newLevel: 'Gold', lastRank: 47 },
              { name: 'Sanjay Kumar', growth: '+38,500 pts', newLevel: 'Silver', lastRank: 63 },
              { name: 'Neha Gupta', growth: '+34,100 pts', newLevel: 'Silver', lastRank: 71 },
              { name: 'Aditya Verma', growth: '+28,900 pts', newLevel: 'Bronze', lastRank: 92 },
              { name: 'Lisa Chen', growth: '+24,600 pts', newLevel: 'Bronze', lastRank: 108 }
            ].map((player, i) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-green-200 bg-green-50">
                <div className="flex justify-between items-center">
                  <div>
                    <div className="font-bold text-gray-800 text-sm">{player.name}</div>
                    <div className="text-xs text-gray-600">From rank #{player.lastRank}</div>
                  </div>
                  <div className="text-right">
                    <div className="text-green-600 font-bold text-sm">{player.growth}</div>
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
              { icon: '⭐', name: 'Rising Star', desc: 'Earn 10K points in 30 days', earned: 87, totalUsers: 245 },
              { icon: '🚀', name: 'Streak Master', desc: '30-day login streak', earned: 124, totalUsers: 245 },
              { icon: '💯', name: 'Perfect Score', desc: 'Achieve 100% on challenge', earned: 156, totalUsers: 245 },
              { icon: '👥', name: 'Team Player', desc: 'Help 5+ team members', earned: 98, totalUsers: 245 }
            ].map((badge, i) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200">
                <div className="flex items-start gap-3">
                  <div className="text-3xl">{badge.icon}</div>
                  <div className="flex-1">
                    <div className="font-bold text-gray-800 text-sm">{badge.name}</div>
                    <div className="text-xs text-gray-600">{badge.desc}</div>
                    <div className="text-xs text-blue-600 mt-1">{badge.earned}/{badge.totalUsers} earned</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Level Progression</h3>
            {[
              { level: 'Bronze', minPoints: '0', maxPoints: '25K', users: 87, color: 'orange' },
              { level: 'Silver', minPoints: '25K', maxPoints: '75K', users: 92, color: 'gray' },
              { level: 'Gold', minPoints: '75K', maxPoints: '150K', users: 56, color: 'yellow' },
              { level: 'Platinum', minPoints: '150K+', maxPoints: '—', users: 10, color: 'purple' }
            ].map((tier, i) => (
              <div key={i} className={`bg-${tier.color}-50 border-l-4 border-${tier.color}-500 p-4 rounded`}>
                <div className="flex justify-between items-start mb-2">
                  <div>
                    <h4 className="font-bold text-gray-800">{tier.level}</h4>
                    <div className="text-xs text-gray-600">{tier.minPoints} - {tier.maxPoints} points</div>
                  </div>
                  <div className="text-right">
                    <div className="text-2xl font-bold text-gray-800">{tier.users}</div>
                    <div className="text-xs text-gray-600">users</div>
                  </div>
                </div>
                <div className="w-full bg-gray-200 h-2 rounded-full overflow-hidden">
                  <div className={`bg-${tier.color}-500 h-full`} style={{ width: `${(tier.users / 87) * 100}%` }}></div>
                </div>
              </div>
            ))}
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
            {[
              {
                title: '🎯 December Sales Sprint',
                desc: 'Target: ₹50 Cr revenue',
                current: '₹42 Cr',
                progress: 84,
                participants: 67,
                endDate: 'Dec 31',
                topContributor: 'Rahul Sharma'
              },
              {
                title: '📝 Quality Challenge',
                desc: 'Zero defects this month',
                current: '98.2% quality',
                progress: 98,
                participants: 145,
                endDate: 'Dec 31',
                topContributor: 'Meera Kapoor'
              },
              {
                title: '🤝 Customer Satisfaction',
                desc: 'Target: 95% CSAT',
                current: '93.5%',
                progress: 93,
                participants: 89,
                endDate: 'Dec 15',
                topContributor: 'Priya Nair'
              },
              {
                title: '⚡ Speed Challenge',
                desc: 'Task completion in 24hrs',
                current: '167/200 tasks',
                progress: 83,
                participants: 102,
                endDate: 'Dec 20',
                topContributor: 'Arjun Singh'
              }
            ].map((challenge, i) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200">
                <div className="font-bold text-gray-800 text-sm mb-1">{challenge.title}</div>
                <div className="text-xs text-gray-600 mb-2">{challenge.desc}</div>
                <div className="flex justify-between items-center text-xs mb-2">
                  <span className="font-bold">{challenge.current}</span>
                  <span className="text-gray-500">{challenge.participants} joining</span>
                </div>
                <div className="w-full bg-gray-200 h-2 rounded-full overflow-hidden mb-2">
                  <div className="bg-green-500 h-full" style={{ width: `${challenge.progress}%` }}></div>
                </div>
                <div className="flex justify-between text-xs text-gray-600">
                  <span>Leader: {challenge.topContributor}</span>
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
              { item: 'Coffee Voucher', cost: '500 pts', claimed: 234, image: '☕' },
              { item: 'Movie Tickets', cost: '2K pts', claimed: 156, image: '🎬' },
              { item: 'Amazon Gift Card', cost: '5K pts', claimed: 89, image: '🎁' },
              { item: 'Spa Voucher', cost: '3.5K pts', claimed: 67, image: '🧖' },
              { item: 'Weekend Getaway', cost: '15K pts', claimed: 12, image: '🏨' },
              { item: 'Team Dinner', cost: '8K pts', claimed: 23, image: '🍽️' }
            ].map((reward, i) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200 text-center hover:shadow-lg">
                <div className="text-4xl mb-2">{reward.image}</div>
                <div className="font-bold text-gray-800 text-sm">{reward.item}</div>
                <div className="text-blue-600 font-bold text-sm mt-1">{reward.cost}</div>
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
              <div className="text-3xl font-bold text-blue-700 mt-2">198 / 245</div>
              <div className="text-xs text-gray-600 mt-1">81% daily engagement</div>
            </div>
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="text-sm text-gray-600">Weekly Engagement</div>
              <div className="text-3xl font-bold text-green-700 mt-2">245 / 245</div>
              <div className="text-xs text-gray-600 mt-1">100% weekly participation</div>
            </div>
            <div className="bg-purple-50 border-l-4 border-purple-500 p-4 rounded">
              <div className="text-sm text-gray-600">Avg Points Per User</div>
              <div className="text-3xl font-bold text-purple-700 mt-2">13,061</div>
              <div className="text-xs text-gray-600 mt-1">Total 3.2M points in system</div>
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
              <div className="text-sm text-gray-700 mt-2">87% engagement rate. 81% daily active users. 245/245 users participated in Q4. System thriving with healthy competition.</div>
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
