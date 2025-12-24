'use server';

/**
 * Secure Token Management with httpOnly Cookies
 * Stores OAuth tokens securely server-side via Prisma
 */

import { cookies } from 'next/headers';
import { prisma } from './prisma';
import type { AuthTokenResponse } from './types';

const SECURE_COOKIE_OPTIONS = {
  httpOnly: true,
  secure: process.env.NODE_ENV === 'production',
  sameSite: 'strict' as const,
  maxAge: 7 * 24 * 60 * 60, // 7 days
};

/**
 * Store token in database (Prisma) and set secure cookie
 */
export async function storeAuthToken(
  userId: string,
  token: string,
  expiresAt: Date
): Promise<AuthTokenResponse> {
  // Create in Prisma database
  const authToken = await prisma.authToken.create({
    data: {
      userId,
      token,
      expiresAt,
    },
  });

  // Also set secure httpOnly cookie
  const cookieStore = await cookies();
  cookieStore.set('auth_token', token, SECURE_COOKIE_OPTIONS);

  return authToken;
}

/**
 * Retrieve token from secure cookie
 */
export async function getStoredToken(): Promise<string | null> {
  const cookieStore = await cookies();
  return cookieStore.get('auth_token')?.value || null;
}

/**
 * Verify token from database
 */
export async function verifyStoredToken(token: string): Promise<boolean> {
  const authToken = await prisma.authToken.findUnique({
    where: { token },
  });

  if (!authToken) return false;

  // Check if expired
  if (new Date() > authToken.expiresAt) {
    await prisma.authToken.delete({
      where: { id: authToken.id },
    });
    return false;
  }

  return true;
}

/**
 * Refresh token if valid
 */
export async function refreshAuthToken(oldToken: string, newToken: string, newExpiresAt: Date) {
  const authToken = await prisma.authToken.findUnique({
    where: { token: oldToken },
  });

  if (!authToken) throw new Error('Token not found');

  // Update in database
  await prisma.authToken.update({
    where: { id: authToken.id },
    data: {
      token: newToken,
      expiresAt: newExpiresAt,
    },
  });

  // Update secure cookie
  const cookieStore = await cookies();
  cookieStore.set('auth_token', newToken, SECURE_COOKIE_OPTIONS);

  return newToken;
}

/**
 * Revoke token
 */
export async function revokeAuthToken(token: string) {
  await prisma.authToken.deleteMany({
    where: { token },
  });

  // Clear cookie
  const cookieStore = await cookies();
  cookieStore.delete('auth_token');
}

/**
 * Clean up expired tokens (run periodically)
 */
export async function cleanupExpiredTokens() {
  await prisma.authToken.deleteMany({
    where: {
      expiresAt: {
        lt: new Date(),
      },
    },
  });
}
