/**
 * OAuth Callback Handler
 * Exchanges authorization code for token and stores in Prisma database
 * Route: /api/auth/callback?provider=google&code=xxx
 */

import { NextRequest, NextResponse } from 'next/server';
import { prisma } from '@/lib/prisma';
import { storeAuthToken } from '@/lib/auth-storage';
import { apiClient } from '@/lib/api-client';

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
    const oauthToken = await apiClient.exchangeOAuthCode(provider, code);

    if (!oauthToken?.accessToken) {
      throw new Error('Invalid OAuth response');
    }

    // TODO: Get user info from OAuth provider
    // For now, create/update user in Prisma database
    const userId = await getOrCreateUser(provider, code);

    // Store token in Prisma database with secure cookie
    await storeAuthToken(userId, oauthToken.accessToken, oauthToken.expiresAt);

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
  // For now, create a temporary user
  // In production: fetch user info from Google/Meta API using the code

  const email = `oauth-${provider}-${Date.now()}@vyom.local`;

  const user = await prisma.user.upsert({
    where: { email },
    update: { updatedAt: new Date() },
    create: {
      email,
      passwordHash: '', // OAuth users don't have passwords
      role: 'user',
      tenantId: process.env.DEFAULT_TENANT_ID || 'default-tenant',
    },
  });

  return user.id;
}
