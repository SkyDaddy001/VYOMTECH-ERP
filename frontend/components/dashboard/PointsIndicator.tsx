'use client'

import React, { useState, useEffect } from 'react'
import Link from 'next/link'
import { calculateLevel, formatLargeNumber } from '@/utils/formatters'
import { COLORS } from '@/styles/designTokens'

interface PointsData {
  currentPoints: number
  dailyPoints: number
  streakDays: number
  nextBadge?: {
    name: string
    pointsNeeded: number
  }
}

export function PointsIndicator() {
  const [points, setPoints] = useState<PointsData | null>(null)
  const [loading, setLoading] = useState(true)
  const [isExpanded, setIsExpanded] = useState(false)

  useEffect(() => {
    loadPoints()
    const interval = setInterval(loadPoints, 60000)
    return () => clearInterval(interval)
  }, [])

  const loadPoints = async () => {
    try {
      const token = localStorage.getItem('auth_token')
      if (!token) return

      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/gamification/profile`,
        { headers: { Authorization: `Bearer ${token}` } }
      )

      if (response.ok) {
        const data = await response.json()
        setPoints({
          currentPoints: data.user_points?.current_points || 0,
          dailyPoints: data.user_points?.daily_points || 0,
          streakDays: data.user_points?.streak_days || 0,
          nextBadge: {
            name: 'Badge Name',
            pointsNeeded: 150,
          },
        })
      }
    } catch (error) {
      console.error('Failed to load points:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading || !points) {
    return null
  }

  const level = calculateLevel(points.currentPoints)

  return (
    <div className="relative">
      {/* Compact Button - Always visible */}
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="flex items-center gap-2 px-4 py-2 rounded-lg bg-gradient-to-r from-yellow-400 to-yellow-500 text-yellow-900 font-semibold hover:shadow-lg transition-all hover:scale-105"
      >
        <span className="text-xl">⭐</span>
        <div className="text-right">
          <div className="text-sm font-bold">{formatLargeNumber(points.currentPoints)}</div>
          <div className="text-xs opacity-75">Pts</div>
        </div>
      </button>

      {/* Expanded Dropdown - Context rich (Tip 9: Give Numbers Context) */}
      {isExpanded && (
        <div className="absolute top-full right-0 mt-2 w-80 bg-white border border-gray-200 rounded-lg shadow-xl z-50 overflow-hidden">
          {/* Header */}
          <div className="bg-gradient-to-r from-yellow-400 to-yellow-500 text-yellow-900 px-6 py-4">
            <div className="flex items-center justify-between mb-2">
              <div className="text-sm font-medium opacity-75">Current Balance</div>
              <button
                onClick={() => setIsExpanded(false)}
                className="text-lg font-bold opacity-50 hover:opacity-100"
              >
                ×
              </button>
            </div>
            <div className="text-4xl font-bold">{formatLargeNumber(points.currentPoints)}</div>
          </div>

          {/* Quick Stats - Scannable (Tip 2: Include Only What's Important) */}
          <div className="px-6 py-4 space-y-3 border-b border-gray-100">
            {/* Today's Progress */}
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Today's Earned</span>
              <span className="font-bold text-emerald-600">{points.dailyPoints} pts</span>
            </div>

            {/* Streak */}
            {points.streakDays > 0 && (
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600">Current Streak</span>
                <span className="font-bold flex items-center gap-1">
                  <span className="text-orange-500">🔥</span>
                  {points.streakDays} days
                </span>
              </div>
            )}

            {/* Level */}
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Level</span>
              <span className="font-bold text-blue-600">Level {level}</span>
            </div>
          </div>

          {/* Next Badge Preview */}
          {points.nextBadge && (
            <div className="px-6 py-4 bg-blue-50 border-t border-blue-100">
              <div className="text-xs text-blue-600 font-semibold mb-2">Next Achievement</div>
              <div className="text-sm font-medium text-gray-900 mb-1">{points.nextBadge.name}</div>
              <div className="text-xs text-gray-500">
                {points.nextBadge.pointsNeeded} points to unlock
              </div>
            </div>
          )}

          {/* Action Buttons */}
          <div className="px-6 py-4 bg-gray-50 border-t border-gray-100 flex gap-2">
            <Link
              href="/dashboard/gamification"
              className="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-3 rounded text-center text-sm transition-colors"
              onClick={() => setIsExpanded(false)}
            >
              View Dashboard
            </Link>
            <Link
              href="/dashboard/rewards"
              className="flex-1 bg-emerald-600 hover:bg-emerald-700 text-white font-medium py-2 px-3 rounded text-center text-sm transition-colors"
              onClick={() => setIsExpanded(false)}
            >
              Rewards Shop
            </Link>
          </div>
        </div>
      )}

      {/* Click anywhere outside to close */}
      {isExpanded && (
        <div
          className="fixed inset-0 z-40"
          onClick={() => setIsExpanded(false)}
        />
      )}
    </div>
  )
}
