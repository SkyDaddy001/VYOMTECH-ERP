/**
 * OAuth Types - Derived from Prisma AuthToken model
 * Ensures all OAuth flows use Prisma-defined structures
 */

export interface OAuthConfig {
  provider: 'google' | 'meta';
  clientId: string;
  clientSecret: string;
  redirectUri: string;
}

export interface OAuthToken {
  accessToken: string;
  refreshToken?: string;
  expiresAt: Date;
  tokenType: string;
  scope: string;
}

export interface OAuthUser {
  id: string;
  email: string;
  name?: string;
  picture?: string;
  provider: 'google' | 'meta';
}

export interface AuthTokenResponse {
  id: string;
  userId: string;
  token: string;
  expiresAt: Date;
  createdAt: Date;
}

/**
 * Campaign Types - From Prisma Campaign model
 */
export interface CampaignInput {
  name: string;
  platform: 'google_ads' | 'meta_ads';
  status: 'active' | 'paused' | 'stopped';
  budget: number;
  startDate: Date;
  endDate?: Date;
}

export interface CampaignResponse {
  id: string;
  name: string;
  platform: string;
  status: string;
  budget: number;
  spend: number;
  impressions: number;
  clicks: number;
  conversions: number;
  startDate: Date;
  endDate?: Date;
  createdAt: Date;
  updatedAt: Date;
}

/**
 * Metrics Response - From Prisma database
 */
export interface MetricsResponse {
  campaignId: string;
  impressions: number;
  clicks: number;
  conversions: number;
  spend: number;
  revenue: number;
  roi: number;
  ctr: number;
  cpc: number;
  conversionRate: number;
  lastSync: Date;
}
