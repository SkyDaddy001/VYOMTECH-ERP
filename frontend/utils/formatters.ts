/**
 * Formatting Utilities
 * Tip 4: Round Your Numbers - Consistent formatting across dashboard
 */

import { NUMBER_FORMAT, THRESHOLDS, COLORS } from '@/styles/designTokens'

// ============================================================================
// DURATION FORMATTING
// ============================================================================

export function formatDuration(seconds: number | null): string {
  if (seconds === null || seconds === undefined) return 'N/A'
  
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = Math.floor(seconds % 60)
  
  if (hours > 0) return `${hours}h ${minutes}m`
  if (minutes > 0) return `${minutes}m ${secs}s`
  return `${secs}s`
}

// ============================================================================
// LARGE NUMBER FORMATTING
// ============================================================================

export function formatLargeNumber(num: number | null): string {
  if (num === null || num === undefined) return 'N/A'
  
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

// ============================================================================
// PERCENTAGE FORMATTING
// ============================================================================

export function formatPercentage(num: number | null, decimals = 0): string {
  if (num === null || num === undefined) return 'N/A'
  return num.toFixed(decimals) + '%'
}

// ============================================================================
// RATING FORMATTING
// ============================================================================

export function formatRating(num: number | null, maxRating = 5): string {
  if (num === null || num === undefined) return 'N/A'
  return `${num.toFixed(1)} / ${maxRating}`
}

// ============================================================================
// MONEY FORMATTING
// ============================================================================

export function formatMoney(num: number | null): string {
  if (num === null || num === undefined) return 'N/A'
  return '₹' + num.toLocaleString('en-US', { maximumFractionDigits: 0 })
}

// ============================================================================
// STATUS CONTEXT - Give numbers meaning (Tip 9)
// ============================================================================

export interface StatusContext {
  status: 'good' | 'warning' | 'critical'
  color: string
  icon: string
  message: string
}

export function getQueueStatus(waitSeconds: number | null): StatusContext {
  if (waitSeconds === null) {
    return {
      status: 'good',
      color: COLORS.success,
      icon: '✓',
      message: 'Within target',
    }
  }

  if (waitSeconds > THRESHOLDS.queueWaitTime.critical) {
    return {
      status: 'critical',
      color: COLORS.critical,
      icon: '✗',
      message: 'Critical - Exceeds 10 minutes',
    }
  }

  if (waitSeconds > THRESHOLDS.queueWaitTime.warning) {
    return {
      status: 'warning',
      color: COLORS.warning,
      icon: '⚠',
      message: 'Warning - Exceeds 5 minutes',
    }
  }

  return {
    status: 'good',
    color: COLORS.success,
    icon: '✓',
    message: 'Good - Within target',
  }
}

export function getCSATStatus(rating: number | null): StatusContext {
  if (rating === null) {
    return {
      status: 'good',
      color: COLORS.success,
      icon: '✓',
      message: 'No data',
    }
  }

  if (rating < THRESHOLDS.customerSatisfaction.critical) {
    return {
      status: 'critical',
      color: COLORS.critical,
      icon: '✗',
      message: `Below ${THRESHOLDS.customerSatisfaction.critical} - Critical`,
    }
  }

  if (rating < THRESHOLDS.customerSatisfaction.warning) {
    return {
      status: 'warning',
      color: COLORS.warning,
      icon: '⚠',
      message: `Below ${THRESHOLDS.customerSatisfaction.warning} - Warning`,
    }
  }

  if (rating >= THRESHOLDS.customerSatisfaction.target) {
    return {
      status: 'good',
      color: COLORS.success,
      icon: '✓',
      message: `Exceeds target (${THRESHOLDS.customerSatisfaction.target})`,
    }
  }

  return {
    status: 'good',
    color: COLORS.success,
    icon: '✓',
    message: 'Within target',
  }
}

export function getCompletionRateStatus(rate: number | null): StatusContext {
  if (rate === null) {
    return {
      status: 'good',
      color: COLORS.success,
      icon: '✓',
      message: 'No data',
    }
  }

  if (rate < THRESHOLDS.callCompletionRate.critical) {
    return {
      status: 'critical',
      color: COLORS.critical,
      icon: '✗',
      message: `Below ${THRESHOLDS.callCompletionRate.critical}% - Critical`,
    }
  }

  if (rate < THRESHOLDS.callCompletionRate.warning) {
    return {
      status: 'warning',
      color: COLORS.warning,
      icon: '⚠',
      message: `Below ${THRESHOLDS.callCompletionRate.warning}% - Warning`,
    }
  }

  if (rate >= THRESHOLDS.callCompletionRate.target) {
    return {
      status: 'good',
      color: COLORS.success,
      icon: '✓',
      message: `Exceeds target (${THRESHOLDS.callCompletionRate.target}%)`,
    }
  }

  return {
    status: 'good',
    color: COLORS.success,
    icon: '✓',
    message: 'Within target',
  }
}

// ============================================================================
// TREND FORMATTING
// ============================================================================

export interface TrendInfo {
  direction: 'up' | 'down' | 'stable'
  icon: string
  percent: number
  color: string
}

export function getTrendInfo(current: number, previous: number | null): TrendInfo {
  if (previous === null || previous === 0) {
    return {
      direction: 'stable',
      icon: '→',
      percent: 0,
      color: COLORS.neutral,
    }
  }

  const percentChange = ((current - previous) / previous) * 100

  if (percentChange > 5) {
    return {
      direction: 'up',
      icon: '↑',
      percent: percentChange,
      color: COLORS.success,
    }
  }

  if (percentChange < -5) {
    return {
      direction: 'down',
      icon: '↓',
      percent: Math.abs(percentChange),
      color: COLORS.critical,
    }
  }

  return {
    direction: 'stable',
    icon: '→',
    percent: 0,
    color: COLORS.neutral,
  }
}

export function formatTrendLabel(trend: TrendInfo, metricName: string): string {
  if (trend.direction === 'stable') {
    return `No change vs yesterday`
  }

  const symbol = trend.direction === 'up' ? '+' : '-'
  return `${symbol}${trend.percent.toFixed(1)}% vs yesterday`
}

// ============================================================================
// DATE FORMATTING
// ============================================================================

export function formatDate(date: Date | string | null): string {
  if (date === null) return 'N/A'

  const d = typeof date === 'string' ? new Date(date) : date
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

export function formatDateTime(date: Date | string | null): string {
  if (date === null) return 'N/A'

  const d = typeof date === 'string' ? new Date(date) : date
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// ============================================================================
// LEVEL CALCULATION
// ============================================================================

export function calculateLevel(points: number): number {
  return Math.floor(points / 1000) + 1
}

export function calculateProgressToNextLevel(points: number): {
  current: number
  total: number
  percent: number
} {
  const currentLevel = calculateLevel(points)
  const pointsForCurrentLevel = (currentLevel - 1) * 1000
  const pointsForNextLevel = currentLevel * 1000
  const current = points - pointsForCurrentLevel
  const total = pointsForNextLevel - pointsForCurrentLevel

  return {
    current,
    total,
    percent: (current / total) * 100,
  }
}

// ============================================================================
// RARITY BADGE
// ============================================================================

export function getRarityBadgeColor(rarity: string): string {
  const rarityMap: { [key: string]: string } = {
    common: 'bg-gray-100 text-gray-800 border-gray-300',
    uncommon: 'bg-green-100 text-green-800 border-green-300',
    rare: 'bg-blue-100 text-blue-800 border-blue-300',
    epic: 'bg-purple-100 text-purple-800 border-purple-300',
    legendary: 'bg-yellow-100 text-yellow-800 border-yellow-300',
  }
  return rarityMap[rarity] || rarityMap.common
}

// ============================================================================
// COMPARISON TEXT
// ============================================================================

export function getComparisonText(
  current: number,
  goal: number,
  format?: (n: number) => string
): string {
  const formatter = format || ((n) => n.toString())
  const currentStr = formatter(current)
  const goalStr = formatter(goal)

  if (current >= goal) {
    return `${currentStr} (Goal: ${goalStr}) ✓ Exceeded`
  }
  return `${currentStr} (Goal: ${goalStr})`
}

// ============================================================================
// RANGE TEXT
// ============================================================================

export function getRangeText(
  current: number,
  min: number,
  max: number,
  format?: (n: number) => string
): string {
  const formatter = format || ((n) => n.toString())
  const currentStr = formatter(current)
  const minStr = formatter(min)
  const maxStr = formatter(max)

  return `${currentStr} (Min: ${minStr}, Max: ${maxStr})`
}

export default {
  formatDuration,
  formatLargeNumber,
  formatPercentage,
  formatRating,
  formatMoney,
  getQueueStatus,
  getCSATStatus,
  getCompletionRateStatus,
  getTrendInfo,
  formatTrendLabel,
  formatDate,
  formatDateTime,
  calculateLevel,
  calculateProgressToNextLevel,
  getRarityBadgeColor,
  getComparisonText,
  getRangeText,
}
