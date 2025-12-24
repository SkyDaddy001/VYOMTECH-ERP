// ============================================================
// NEXT.JS EDGE MIDDLEWARE - Attribution Web Capture
// ============================================================
// Location: frontend/src/middleware.ts
// Runs on every request, captures UTM params and custom fields
// Stores in cookies for 90 days, sends async to backend

import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// ============================================================
// CONFIGURATION
// ============================================================

const ATTRIBUTION_COOKIE_NAME = 'vyom_attribution'
const ATTRIBUTION_COOKIE_MAX_AGE = 90 * 24 * 60 * 60 // 90 days in seconds
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
const ATTRIBUTION_ENDPOINT = `${API_BASE_URL}/api/v1/attribution/events`

// ============================================================
// TYPES
// ============================================================

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

interface IngestedAttributionData {
  leadId?: string
  touchType: 'click' // or 'form' when form is submitted
  source: string
  subSource?: string
  medium?: string
  campaign?: string
  adId?: string
  creativeId?: string
  landingPage: string
  referrer?: string
  utm_source?: string
  utm_medium?: string
  utm_campaign?: string
  utm_content?: string
  utm_term?: string
  device?: string
  userAgent?: string
  ipAddress?: string
  country?: string
  timezone?: string
  sessionId: string
  occurredAt: string // ISO timestamp
  idempotencyKey: string // Unique: sessionId + timestamp + utm_source
}

// ============================================================
// MIDDLEWARE
// ============================================================

export function middleware(request: NextRequest) {
  // Skip attribution tracking for:
  // - API routes
  // - Static assets
  // - Admin pages
  // - Already tracked pages
  if (
    request.nextUrl.pathname.startsWith('/api') ||
    request.nextUrl.pathname.startsWith('/_next') ||
    request.nextUrl.pathname.match(/\.(ico|png|svg|jpg|jpeg|gif|webp)$/)
  ) {
    return NextResponse.next()
  }

  // Parse request
  const response = NextResponse.next()
  const url = request.nextUrl
  const requestHeaders = request.headers

  // 1. Extract UTM parameters from URL
  const utm = extractUTMParameters(url.searchParams)

  // 2. Extract or create session ID
  const sessionId = getOrCreateSessionId(request, response)

  // 3. Detect device type
  const device = detectDevice(requestHeaders.get('user-agent') || '')

  // 4. Build attribution data
  const attributionData: AttributionData = {
    sessionId,
    firstTouchAt: new Date().toISOString(),
    utmSource: utm.utmSource ?? undefined,
    utmMedium: utm.utmMedium ?? undefined,
    utmCampaign: utm.utmCampaign ?? undefined,
    utmContent: utm.utmContent ?? undefined,
    utmTerm: utm.utmTerm ?? undefined,
    referrer: requestHeaders.get('referer') || undefined,
    landingPage: url.pathname,
    device,
    touchCount: 1,
    lastPageVisited: url.pathname,
    lastVisitAt: new Date().toISOString(),
  }

  // 5. Store in cookie (for use in client-side forms)
  const cookieValue = JSON.stringify(attributionData)
  response.cookies.set(ATTRIBUTION_COOKIE_NAME, cookieValue, {
    maxAge: ATTRIBUTION_COOKIE_MAX_AGE,
    secure: process.env.NODE_ENV === 'production',
    httpOnly: false, // Needed for client-side access
    sameSite: 'lax',
    path: '/',
  })

  // 6. Ingest attribution event asynchronously (non-blocking)
  // This happens in the background via next request
  ingestAttributionEvent({
    touchType: 'click',
    source: utm.source || 'direct',
    subSource: utm.subSource,
    medium: utm.medium || 'organic',
    campaign: utm.utmCampaign ?? undefined,
    adId: url.searchParams.get('ad_id') || undefined,
    creativeId: url.searchParams.get('creative_id') || undefined,
    landingPage: url.pathname,
    referrer: requestHeaders.get('referer') || undefined,
    utm_source: utm.utmSource ?? undefined,
    utm_medium: utm.utmMedium ?? undefined,
    utm_campaign: utm.utmCampaign ?? undefined,
    utm_content: utm.utmContent ?? undefined,
    utm_term: utm.utmTerm ?? undefined,
    device,
    userAgent: requestHeaders.get('user-agent') || undefined,
    ipAddress: (requestHeaders.get('x-forwarded-for') || requestHeaders.get('x-real-ip')) ?? undefined,
    country: requestHeaders.get('cloudflare-country-code') || undefined,
    timezone: undefined,
    sessionId,
    occurredAt: new Date().toISOString(),
    idempotencyKey: generateIdempotencyKey(sessionId, utm.utmSource || 'direct'),
  }).catch(error => {
    console.error('Failed to ingest attribution event:', error)
  })

  return response
}

// ============================================================
// UTILITY FUNCTIONS
// ============================================================

/**
 * Extract UTM parameters from URL
 */
function extractUTMParameters(searchParams: URLSearchParams) {
  const utmSource = searchParams.get('utm_source')
  const utmMedium = searchParams.get('utm_medium')
  const utmCampaign = searchParams.get('utm_campaign')
  const utmContent = searchParams.get('utm_content')
  const utmTerm = searchParams.get('utm_term')

  // Infer source from UTM
  let source = 'direct'
  let subSource: string | undefined
  let medium: string | undefined

  if (utmSource) {
    source = utmSource
    if (utmSource === 'google') {
      subSource = utmMedium === 'cpc' ? 'google_ads' : 'google_organic'
    } else if (utmSource === 'facebook') {
      subSource = 'facebook_ads'
    } else if (utmSource === 'email') {
      medium = 'email'
      source = 'email'
    }
  }

  return {
    source,
    subSource,
    medium: medium || utmMedium,
    utmSource,
    utmMedium,
    utmCampaign,
    utmContent,
    utmTerm,
  }
}

/**
 * Get or create session ID (stored in cookie)
 */
function getOrCreateSessionId(request: NextRequest, response: NextResponse): string {
  const existingSession = request.cookies.get('vyom_session')?.value
  if (existingSession) {
    return existingSession
  }

  // Generate new session ID
  const newSessionId = `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  response.cookies.set('vyom_session', newSessionId, {
    maxAge: 180 * 24 * 60 * 60, // 180 days
    secure: process.env.NODE_ENV === 'production',
    httpOnly: false,
    sameSite: 'lax',
    path: '/',
  })

  return newSessionId
}

/**
 * Detect device type from user agent
 */
function detectDevice(userAgent: string): string {
  if (/mobile|android|iphone|ipod/i.test(userAgent)) {
    return 'mobile'
  }
  if (/tablet|ipad/i.test(userAgent)) {
    return 'tablet'
  }
  return 'desktop'
}

/**
 * Generate unique idempotency key
 * Format: sessionId_timestamp_source_hash
 * Ensures deduplication on backend
 */
function generateIdempotencyKey(sessionId: string, source: string): string {
  const timestamp = Date.now()
  const randomHash = Math.random().toString(36).substr(2, 9)
  return `${sessionId}_${timestamp}_${source}_${randomHash}`
}

/**
 * Ingest attribution event to backend asynchronously
 * This is non-blocking and doesn't slow down the request
 */
async function ingestAttributionEvent(event: Partial<IngestedAttributionData>) {
  try {
    // Fire and forget - no await, no blocking
    if (typeof window === 'undefined') {
      // Server-side: use fetch API with no-blocking pattern
      fetch(ATTRIBUTION_ENDPOINT, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Ingestion-Source': 'nextjs-middleware',
        },
        body: JSON.stringify({
          tenant_id: process.env.NEXT_PUBLIC_TENANT_ID || 'default',
          // lead_id will be set when user logs in or fills a form
          ...event,
          occurredAt: event.occurredAt || new Date().toISOString(),
        }),
      }).catch(error => {
        // Fail silently - don't interrupt user experience
        console.debug('Attribution ingestion error (non-blocking):', error)
      })
    }
  } catch (error) {
    console.debug('Attribution ingest failed:', error)
  }
}

// ============================================================
// MATCHER - Which routes to apply middleware to
// ============================================================

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - api (API routes - handled separately)
     */
    '/((?!_next/static|_next/image|favicon.ico|api).*)',
  ],
}

// ============================================================
// CLIENT-SIDE ATTRIBUTION HOOK MOVED
// ============================================================
// The useAttribution hook has been moved to:
// frontend/src/hooks/useAttribution.ts
// Import it in your client components with:
// import { useAttribution } from '@/hooks/useAttribution'
