'use client'

import { useEffect, useState } from 'react'

interface AttributionData {
  sessionId: string
  firstTouchAt: string
  utmSource?: string
  utmMedium?: string
  utmCampaign?: string
  utmContent?: string
  utmTerm?: string
  referrer?: string
  landingPage?: string
  device?: string
  touchCount: number
  lastPageVisited?: string
  lastVisitAt?: string
}

const ATTRIBUTION_COOKIE_NAME = 'vyom_attribution'

function getCookie(name: string): string | null {
  if (typeof document === 'undefined') return null
  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) return parts.pop()?.split(';').shift() || null
  return null
}

/**
 * Attribution Hook - Get attribution data captured by middleware
 * Usage in components: useAttribution()
 */
export function useAttribution() {
  const [attributionData, setAttributionData] = useState<AttributionData | null>(null)

  useEffect(() => {
    // Load attribution data from cookie
    const cookieValue = getCookie(ATTRIBUTION_COOKIE_NAME)
    if (cookieValue) {
      try {
        const data = JSON.parse(cookieValue) as AttributionData
        setAttributionData(data)
      } catch (e) {
        console.debug('Failed to parse attribution cookie')
      }
    }
  }, [])

  return attributionData
}
