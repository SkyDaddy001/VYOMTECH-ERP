/**
 * OAuth Callback Handler
 * Exchanges authorization code for token and stores in Prisma database
 * Route: /api/auth/callback?provider=google&code=xxx
 */

import { NextRequest, NextResponse } from 'next/server';
import { storeAuthToken } from '@/lib/auth-storage';
import { apiClient } from '@/lib/api-client';

// Disable static generation for this route - it's a dynamic API route
export const dynamic = 'force-dynamic';

export async function GET(request: NextRequest) {
  try {
    const searchParams = request.nextUrl.searchParams;
    const provider = searchParams.get('provider') as 'google' | 'meta';
    const code = searchParams.get('code');
    const error = searchParams.get('error');

    // Handle OAuth error response
    if (error) {
      return NextResponse.redirect(
        new URL(`/auth/error?error=${error}&provider=${provider}`, request.url)
      );
    }

    if (!provider || !code) {
      return NextResponse.redirect(new URL('/auth/error?error=missing_params', request.url));
    }

    // Exchange code for token via backend API
    const oauthResponse = await apiClient.exchangeOAuthCode(provider, code);

    if (!oauthResponse?.accessToken) {
      throw new Error('Invalid OAuth response');
    }

    // TODO: Get user info from OAuth provider
    // For now, create/update user in Prisma database
    const userId = await getOrCreateUser(provider, code);

    // Store token in localStorage with expiration
    const expiresAt = oauthResponse.expiresAt || new Date(Date.now() + 3600 * 1000);
    await storeAuthToken(userId, oauthResponse.accessToken, expiresAt);

    // Redirect to dashboard
    return NextResponse.redirect(new URL('/dashboard', request.url));
  } catch (error) {
    console.error('OAuth callback error:', error);
    return NextResponse.redirect(
      new URL(`/auth/error?error=callback_failed`, request.url)
    );
  }
}

/**
 * Get or create user in Prisma database
 * Should be called during OAuth callback
 */
async function getOrCreateUser(provider: string, code: string): Promise<string> {
  // TODO: Implement proper user lookup from OAuth provider info
  // For now, return a temporary user ID
  // In production: fetch user info from Google/Meta API using the code
  
  // Generate a consistent ID based on code
  return `user-${provider}-${Buffer.from(code).toString('base64').substring(0, 12)}`;
}