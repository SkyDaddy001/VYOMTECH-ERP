'use client'

import React, { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import { formatMoney, formatDate } from '@/utils/formatters'
import { COLORS, GRID } from '@/styles/designTokens'

interface Reward {
  id: number
  name: string
  description: string
  category: string
  pointsCost: number
  stock: number
  imageUrl: string
  featured: boolean
  expiryDate: string
}

interface UserRedemption {
  id: number
  rewardName: string
  status: 'pending' | 'approved' | 'completed'
  pointsSpent: number
  redemptionDate: string
  expiryDate: string
}

const CATEGORIES = [
  { id: 'all', label: 'All Rewards' },
  { id: 'discount', label: 'Discounts' },
  { id: 'digital', label: 'Digital' },
  { id: 'experience', label: 'Experiences' },
  { id: 'physical', label: 'Physical' },
]

export function RewardsShop() {
  const [rewards, setRewards] = useState<Reward[]>([])
  const [redemptions, setRedemptions] = useState<UserRedemption[]>([])
  const [userPoints, setUserPoints] = useState(0)
  const [loading, setLoading] = useState(true)
  const [selectedCategory, setSelectedCategory] = useState('all')
  const [redeeming, setRedeeming] = useState<number | null>(null)

  useEffect(() => {
    loadRewardsData()
    const interval = setInterval(loadRewardsData, 60000)
    return () => clearInterval(interval)
  }, [])

  const loadRewardsData = async () => {
    try {
      const token = localStorage.getItem('auth_token')
      if (!token) return

      // Load profile for points
      const profileResponse = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/gamification/profile`,
        { headers: { Authorization: `Bearer ${token}` } }
      )
      if (profileResponse.ok) {
        const data = await profileResponse.json()
        setUserPoints(data.user_points?.current_points || 0)
      }

      // Load rewards
      const rewardsResponse = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/gamification/rewards`,
        { headers: { Authorization: `Bearer ${token}` } }
      )
      if (rewardsResponse.ok) {
        const data = await rewardsResponse.json()
        setRewards(data.rewards || [])
      }

      // Load redemption history
      const historyResponse = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/gamification/redemptions`,
        { headers: { Authorization: `Bearer ${token}` } }
      )
      if (historyResponse.ok) {
        const data = await historyResponse.json()
        setRedemptions(data.redemptions || [])
      }
    } catch (error) {
      console.error('Failed to load rewards:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleRedeem = async (reward: Reward) => {
    if (userPoints < reward.pointsCost) {
      toast.error(`You need ${reward.pointsCost - userPoints} more points`)
      return
    }

    if (reward.stock === 0) {
      toast.error('This reward is out of stock')
      return
    }

    setRedeeming(reward.id)
    try {
      const token = localStorage.getItem('auth_token')
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/gamification/redeem`,
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ reward_id: reward.id, quantity: 1 }),
        }
      )

      if (response.ok) {
        toast.success('Reward redeemed! Check your email for details.')
        setUserPoints(userPoints - reward.pointsCost)
        loadRewardsData()
      } else {
        toast.error('Failed to redeem reward')
      }
    } catch (error) {
      toast.error('Error redeeming reward')
      console.error(error)
    } finally {
      setRedeeming(null)
    }
  }

  const filteredRewards =
    selectedCategory === 'all'
      ? rewards
      : rewards.filter((r) => r.category === selectedCategory)

  const featuredRewards = filteredRewards.filter((r) => r.featured)

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4" />
          <p className="text-gray-600">Loading rewards...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* POINTS BALANCE - Tier 1: Critical Info (Tip 8: Hierarchy) */}
      <div className="bg-gradient-to-r from-emerald-500 to-teal-600 text-white rounded-lg p-8 shadow-lg">
        <div className="flex justify-between items-center">
          <div>
            <div className="text-sm font-medium text-emerald-100">Your Points Balance</div>
            <div className="text-5xl font-bold mt-2">{userPoints.toLocaleString()}</div>
          </div>
          <div className="text-6xl opacity-20">⭐</div>
        </div>
        <div className="mt-6 pt-6 border-t border-emerald-400">
          <p className="text-emerald-100 text-sm">
            Earn more points by completing calls, challenges, and maintaining your streak
          </p>
        </div>
      </div>

      {/* Featured Rewards Carousel - Eye-catching */}
      {featuredRewards.length > 0 && (
        <section>
          <h2 className="text-2xl font-bold text-gray-900 mb-4">⭐ Featured Rewards</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {featuredRewards.slice(0, 2).map((reward) => (
              <div
                key={reward.id}
                className="relative bg-white border-2 border-yellow-300 rounded-lg overflow-hidden hover:shadow-lg transition-shadow"
              >
                <div className="absolute top-3 right-3 bg-yellow-300 text-yellow-900 px-3 py-1 rounded-full text-sm font-bold">
                  Featured
                </div>
                {reward.imageUrl && (
                  <div className="h-40 bg-gray-100 overflow-hidden">
                    <img
                      src={reward.imageUrl}
                      alt={reward.name}
                      className="w-full h-full object-cover"
                      onError={(e) => {
                        ;(e.target as HTMLImageElement).src =
                          'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="100" height="100"%3E%3Crect fill="%23f3f4f6" width="100" height="100"/%3E%3C/svg%3E'
                      }}
                    />
                  </div>
                )}
                <div className="p-6">
                  <h3 className="text-lg font-bold text-gray-900 mb-2">{reward.name}</h3>
                  <p className="text-gray-600 text-sm mb-4">{reward.description}</p>
                  <div className="flex justify-between items-center mb-4">
                    <div className="text-2xl font-bold text-emerald-600">
                      {reward.pointsCost.toLocaleString()} pts
                    </div>
                    {reward.stock > 0 ? (
                      <div className="text-sm text-gray-500">
                        {reward.stock === -1 ? 'Unlimited' : `${reward.stock} left`}
                      </div>
                    ) : (
                      <div className="text-sm text-red-500 font-semibold">Out of Stock</div>
                    )}
                  </div>
                  <button
                    onClick={() => handleRedeem(reward)}
                    disabled={
                      redeeming === reward.id ||
                      userPoints < reward.pointsCost ||
                      reward.stock === 0
                    }
                    className="w-full bg-emerald-500 hover:bg-emerald-600 disabled:bg-gray-300 text-white font-semibold py-3 rounded-lg transition-colors"
                  >
                    {redeeming === reward.id ? 'Redeeming...' : 'Redeem Now'}
                  </button>
                </div>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* Category Filter - Clear navigation (Tip 10: Clear Labels) */}
      <section>
        <h2 className="text-2xl font-bold text-gray-900 mb-4">Browse Rewards</h2>
        <div className="flex gap-2 overflow-x-auto pb-4">
          {CATEGORIES.map((cat) => (
            <button
              key={cat.id}
              onClick={() => setSelectedCategory(cat.id)}
              className={`px-4 py-2 rounded-full font-medium text-sm transition-all whitespace-nowrap ${
                selectedCategory === cat.id
                  ? 'bg-blue-600 text-white shadow-md'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              {cat.label}
            </button>
          ))}
        </div>
      </section>

      {/* Rewards Grid - Consistent layout */}
      {filteredRewards.length > 0 ? (
        <div className={`${GRID.tier2} gap-4`}>
          {filteredRewards.map((reward) => (
            <div
              key={reward.id}
              className="bg-white border border-gray-200 rounded-lg overflow-hidden hover:shadow-md transition-shadow flex flex-col"
            >
              {/* Image */}
              {reward.imageUrl && (
                <div className="h-32 bg-gray-100 overflow-hidden flex-shrink-0">
                  <img
                    src={reward.imageUrl}
                    alt={reward.name}
                    className="w-full h-full object-cover"
                    onError={(e) => {
                      ;(e.target as HTMLImageElement).src =
                        'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="100" height="100"%3E%3Crect fill="%23f3f4f6" width="100" height="100"/%3E%3C/svg%3E'
                    }}
                  />
                </div>
              )}

              {/* Content */}
              <div className="p-4 flex-1 flex flex-col">
                <h3 className="font-semibold text-gray-900 text-sm mb-2 line-clamp-2">
                  {reward.name}
                </h3>
                <p className="text-gray-600 text-xs mb-3 line-clamp-2">{reward.description}</p>

                {/* Points Cost */}
                <div className="mb-3 mt-auto">
                  <div className="text-xl font-bold text-emerald-600">
                    {reward.pointsCost.toLocaleString()} pts
                  </div>
                </div>

                {/* Stock Status */}
                {reward.stock === 0 ? (
                  <div className="bg-red-50 border border-red-200 rounded px-3 py-2 text-center mb-3">
                    <div className="text-xs font-semibold text-red-600">Out of Stock</div>
                  </div>
                ) : (
                  <div className="text-xs text-gray-500 text-center mb-3">
                    {reward.stock === -1 ? '✓ Always Available' : `${reward.stock} left`}
                  </div>
                )}

                {/* Redeem Button */}
                <button
                  onClick={() => handleRedeem(reward)}
                  disabled={
                    redeeming === reward.id ||
                    userPoints < reward.pointsCost ||
                    reward.stock === 0
                  }
                  className={`w-full py-2 px-3 rounded font-medium text-sm transition-all ${
                    userPoints < reward.pointsCost
                      ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                      : reward.stock === 0
                      ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                      : 'bg-emerald-500 hover:bg-emerald-600 text-white cursor-pointer'
                  }`}
                >
                  {redeeming === reward.id ? '...' : 'Redeem'}
                </button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="bg-white border border-gray-200 rounded-lg p-12 text-center">
          <div className="text-4xl mb-4">🎁</div>
          <p className="text-gray-600 font-medium">No rewards available in this category</p>
        </div>
      )}

      {/* Recent Redemptions - Supporting info (Tip 6: Group Related Metrics) */}
      {redemptions.length > 0 && (
        <section>
          <h2 className="text-xl font-bold text-gray-900 mb-4">Recent Redemptions</h2>
          <div className="bg-white border border-gray-200 rounded-lg overflow-hidden">
            <div className="divide-y">
              {redemptions.slice(0, 5).map((redemption) => (
                <div
                  key={redemption.id}
                  className="p-4 flex justify-between items-center hover:bg-gray-50"
                >
                  <div className="flex-1">
                    <div className="font-medium text-gray-900">{redemption.rewardName}</div>
                    <div className="text-xs text-gray-500">
                      {formatDate(redemption.redemptionDate)}
                    </div>
                  </div>
                  <div className="text-right mr-4">
                    <div className="font-semibold text-gray-900">
                      -{redemption.pointsSpent.toLocaleString()} pts
                    </div>
                  </div>
                  <div className="text-xs">
                    {redemption.status === 'pending' && (
                      <span className="bg-yellow-100 text-yellow-800 px-2 py-1 rounded">
                        Pending
                      </span>
                    )}
                    {redemption.status === 'approved' && (
                      <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded">
                        Approved
                      </span>
                    )}
                    {redemption.status === 'completed' && (
                      <span className="bg-green-100 text-green-800 px-2 py-1 rounded">
                        Completed
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>
      )}
    </div>
  )
}
