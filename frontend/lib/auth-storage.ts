'use server';

/**
 * Secure Token Management with httpOnly Cookies
 * Stores OAuth tokens securely server-side via cookies
 * Note: Database storage via Prisma can be added later
 */

import { cookies } from 'next/headers';
import type { AuthTokenResponse } from './types';

const SECURE_COOKIE_OPTIONS = {
  httpOnly: true,
  secure: process.env.NODE_ENV === 'production',
  sameSite: 'strict' as const,
  maxAge: 7 * 24 * 60 * 60, // 7 days
};

/**
 * Store token in secure httpOnly cookie
 * Note: Database storage can be implemented in backend API
 */
export async function storeAuthToken(
  userId: string,
  token: string,
  expiresAt: Date
): Promise<AuthTokenResponse> {
  // Set secure httpOnly cookie
  const cookieStore = await cookies();
  cookieStore.set('auth_token', token, SECURE_COOKIE_OPTIONS);
  cookieStore.set('user_id', userId, SECURE_COOKIE_OPTIONS);

  // Return token response
  return {
    userId,
    token,
    expiresAt,
  } as AuthTokenResponse;
}

/**
 * Retrieve token from secure cookie
 */
export async function getStoredToken(): Promise<string | null> {
  const cookieStore = await cookies();
  return cookieStore.get('auth_token')?.value || null;
}

/**
 * Verify token from cookie
 * Note: Full validation should be done on backend API
 */
export async function verifyStoredToken(token: string): Promise<boolean> {
  const storedToken = await getStoredToken();
  return storedToken === token;
}

/**
 * Refresh token - update cookie with new token
 * Note: Token refresh logic should be done via backend API
 */
export async function refreshAuthToken(oldToken: string, newToken: string, newExpiresAt: Date) {
  // Update secure cookie
  const cookieStore = await cookies();
  cookieStore.set('auth_token', newToken, SECURE_COOKIE_OPTIONS);

  return newToken;
}

/**
 * Revoke token by clearing cookie
 */
export async function revokeAuthToken(token: string) {
  // Clear cookie
  const cookieStore = await cookies();
  cookieStore.delete('auth_token');
  cookieStore.delete('user_id');
}

/**
 * Clean up expired tokens (can be called from API endpoint)
 * Note: For now, tokens are managed by browser session and cookies
 */
export async function cleanupExpiredTokens() {
  // TODO: Implement backend API endpoint to cleanup expired tokens
  // For now, rely on cookie expiration
}

