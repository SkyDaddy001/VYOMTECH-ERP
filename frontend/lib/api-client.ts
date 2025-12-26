/**
 * API Client with Prisma-aware types
 * All requests use types derived from Prisma schema
 */

'use client';

import axios, { AxiosInstance } from 'axios';
import { getStoredToken } from './auth-storage';
import type {
  CampaignResponse,
  MetricsResponse,
  OAuthToken,
} from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

class APIClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 10000,
    });

    // Add token to every request
    this.client.interceptors.request.use(async (config) => {
      const token = await getStoredToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Handle 401 responses
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Token expired or invalid
          if (typeof window !== 'undefined') {
            window.location.href = '/auth/login';
          }
        }
        return Promise.reject(error);
      }
    );
  }

  // ===== CAMPAIGN ENDPOINTS =====

  /**
   * Create campaign (Prisma Campaign model)
   */
  async createCampaign(platform: 'google' | 'meta', data: any): Promise<CampaignResponse> {
    const endpoint = platform === 'google' 
      ? '/google-ads/campaigns' 
      : '/meta-ads/campaigns';
    const response = await this.client.post(endpoint, data);
    return response.data.data;
  }

  /**
   * List campaigns for tenant
   */
  async listCampaigns(platform?: 'google' | 'meta'): Promise<CampaignResponse[]> {
    if (platform === 'google') {
      const response = await this.client.get('/google-ads/campaigns');
      return response.data.data;
    } else if (platform === 'meta') {
      const response = await this.client.get('/meta-ads/campaigns');
      return response.data.data;
    } else {
      // Get both platforms
      const [googleRes, metaRes] = await Promise.all([
        this.client.get('/google-ads/campaigns'),
        this.client.get('/meta-ads/campaigns'),
      ]);
      return [...googleRes.data.data, ...metaRes.data.data];
    }
  }

  /**
   * Get single campaign by ID
   */
  async getCampaign(campaignId: string, platform: 'google' | 'meta'): Promise<CampaignResponse> {
    const endpoint = platform === 'google' 
      ? `/google-ads/campaigns/${campaignId}` 
      : `/meta-ads/campaigns/${campaignId}`;
    const response = await this.client.get(endpoint);
    return response.data.data;
  }

  /**
   * Update campaign
   */
  async updateCampaign(
    campaignId: string,
    platform: 'google' | 'meta',
    data: any
  ): Promise<CampaignResponse> {
    const endpoint = platform === 'google' 
      ? `/google-ads/campaigns/${campaignId}` 
      : `/meta-ads/campaigns/${campaignId}`;
    const response = await this.client.patch(endpoint, data);
    return response.data.data;
  }

  /**
   * Pause campaign
   */
  async pauseCampaign(campaignId: string, platform: 'google' | 'meta'): Promise<void> {
    const endpoint = platform === 'google' 
      ? `/google-ads/campaigns/${campaignId}/pause` 
      : `/meta-ads/campaigns/${campaignId}/pause`;
    await this.client.post(endpoint);
  }

  /**
   * Resume campaign
   */
  async resumeCampaign(campaignId: string, platform: 'google' | 'meta'): Promise<void> {
    const endpoint = platform === 'google' 
      ? `/google-ads/campaigns/${campaignId}/resume` 
      : `/meta-ads/campaigns/${campaignId}/resume`;
    await this.client.post(endpoint);
  }

  /**
   * Update campaign budget
   */
  async updateBudget(
    campaignId: string,
    platform: 'google' | 'meta',
    dailyBudget: number
  ): Promise<void> {
    const endpoint = platform === 'google' 
      ? `/google-ads/campaigns/${campaignId}/budget` 
      : `/meta-ads/campaigns/${campaignId}/budget`;
    await this.client.patch(endpoint, { dailyBudget });
  }

  // ===== METRICS ENDPOINTS =====

  /**
   * Get campaign metrics (from Prisma database)
   */
  async getCampaignMetrics(campaignId: string, platform: 'google' | 'meta'): Promise<MetricsResponse> {
    const endpoint = platform === 'google' 
      ? `/google-ads/campaigns/${campaignId}/metrics` 
      : `/meta-ads/campaigns/${campaignId}/metrics`;
    const response = await this.client.get(endpoint);
    return response.data.data;
  }

  /**
   * Sync metrics with Prisma database
   */
  async syncMetrics(campaignId: string, platform: 'google' | 'meta'): Promise<MetricsResponse> {
    const endpoint = platform === 'google' 
      ? `/google-ads/campaigns/${campaignId}/sync-metrics` 
      : `/meta-ads/campaigns/${campaignId}/sync-metrics`;
    const response = await this.client.post(endpoint);
    return response.data.data;
  }

  // ===== ROI ENDPOINTS =====

  /**
   * Get campaign ROI
   */
  async getCampaignROI(campaignId: string): Promise<any> {
    const response = await this.client.get(`/roi/campaigns/${campaignId}`);
    return response.data.data;
  }

  /**
   * Get portfolio ROI
   */
  async getPortfolioROI(): Promise<any> {
    const response = await this.client.get('/roi/portfolio');
    return response.data.data;
  }

  /**
   * Get platform comparison ROI
   */
  async getPlatformROI(): Promise<any> {
    const response = await this.client.get('/roi/platform');
    return response.data.data;
  }

  // ===== OAUTH ENDPOINTS =====

  /**
   * Get OAuth authorization URL
   */
  async getAuthorizationUrl(provider: 'google' | 'meta'): Promise<string> {
    const response = await this.client.get(`/oauth/${provider}/authorize`);
    return response.data.authorizationUrl;
  }

  /**
   * Exchange code for token
   */
  async exchangeOAuthCode(provider: 'google' | 'meta', code: string): Promise<OAuthToken> {
    const response = await this.client.post(`/oauth/${provider}/callback`, { code });
    return response.data.data;
  }

  /**
   * Refresh OAuth token
   */
  async refreshOAuthToken(provider: 'google' | 'meta'): Promise<OAuthToken> {
    const response = await this.client.post(`/oauth/${provider}/refresh`);
    return response.data.data;
  }

  /**
   * Validate OAuth token
   */
  async validateOAuthToken(provider: 'google' | 'meta'): Promise<boolean> {
    try {
      const response = await this.client.get(`/oauth/${provider}/validate`);
      return response.data.valid;
    } catch {
      return false;
    }
  }

  /**
   * Revoke OAuth token
   */
  async revokeOAuthToken(provider: 'google' | 'meta'): Promise<void> {
    await this.client.post(`/oauth/${provider}/revoke`);
  }

  // ===== SALES ENDPOINTS =====

  /**
   * List sales
   */
  async listSales(): Promise<any[]> {
    const response = await this.client.get('/sales');
    return response.data.data || [];
  }

  /**
   * Create sale
   */
  async createSale(data: any): Promise<any> {
    const response = await this.client.post('/sales', data);
    return response.data.data;
  }

  /**
   * Delete sale
   */
  async deleteSale(id: string): Promise<void> {
    await this.client.delete(`/sales/${id}`);
  }

  // ===== LEADS ENDPOINTS =====

  /**
   * List leads
   */
  async listLeads(): Promise<any[]> {
    const response = await this.client.get('/leads');
    return response.data.data || [];
  }

  /**
   * Create lead
   */
  async createLead(data: any): Promise<any> {
    const response = await this.client.post('/leads', data);
    return response.data.data;
  }

  /**
   * Delete lead
   */
  async deleteLead(id: string): Promise<void> {
    await this.client.delete(`/leads/${id}`);
  }

  // ===== MARKETING ENDPOINTS =====

  /**
   * List marketing campaigns
   */
  async listMarketingCampaigns(): Promise<any[]> {
    const response = await this.client.get('/marketing-campaigns');
    return response.data.data || [];
  }

  /**
   * Create marketing campaign
   */
  async createMarketingCampaign(data: any): Promise<any> {
    const response = await this.client.post('/marketing-campaigns', data);
    return response.data.data;
  }

  /**
   * Delete marketing campaign
   */
  async deleteMarketingCampaign(id: string): Promise<void> {
    await this.client.delete(`/marketing-campaigns/${id}`);
  }

  // ===== SITE VISITS ENDPOINTS =====

  /**
   * List site visits
   */
  async listSiteVisits(): Promise<any[]> {
    const response = await this.client.get('/site-visits');
    return response.data.data || [];
  }

  /**
   * Create site visit
   */
  async createSiteVisit(data: any): Promise<any> {
    const response = await this.client.post('/site-visits', data);
    return response.data.data;
  }

  /**
   * Delete site visit
   */
  async deleteSiteVisit(id: string): Promise<void> {
    await this.client.delete(`/site-visits/${id}`);
  }

  // ===== SPRINTS ENDPOINTS =====

  /**
   * List sprints
   */
  async listSprints(): Promise<any[]> {
    const response = await this.client.get('/sprints');
    return response.data.data || [];
  }

  /**
   * Create sprint
   */
  async createSprint(data: any): Promise<any> {
    const response = await this.client.post('/sprints', data);
    return response.data.data;
  }

  /**
   * Delete sprint
   */
  async deleteSprint(id: string): Promise<void> {
    await this.client.delete(`/sprints/${id}`);
  }

  // ===== AUTH ENDPOINTS =====

  /**
   * Login user
   */
  async login(email: string, password: string): Promise<any> {
    const response = await this.client.post('/auth/login', { email, password });
    return response.data.data;
  }

  /**
   * Get current user
   */
  async getCurrentUser(): Promise<any> {
    const response = await this.client.get('/auth/me');
    return response.data.data;
  }

  // ===== GENERIC HTTP METHODS =====
  /**
   * Generic GET request
   */
  async get<T = any>(endpoint: string, config?: any): Promise<T> {
    const response = await this.client.get<T>(endpoint, config);
    return (response as any).data?.data || (response as any).data;
  }

  /**
   * Generic POST request
   */
  async post<T = any>(endpoint: string, data?: any, config?: any): Promise<T> {
    const response = await this.client.post<T>(endpoint, data, config);
    return (response as any).data?.data || (response as any).data;
  }

  /**
   * Generic PATCH request
   */
  async patch<T = any>(endpoint: string, data?: any, config?: any): Promise<T> {
    const response = await this.client.patch<T>(endpoint, data, config);
    return (response as any).data?.data || (response as any).data;
  }

  /**
   * Generic PUT request
   */
  async put<T = any>(endpoint: string, data?: any, config?: any): Promise<T> {
    const response = await this.client.put<T>(endpoint, data, config);
    return (response as any).data?.data || (response as any).data;
  }

  /**
   * Generic DELETE request
   */
  async delete<T = any>(endpoint: string, config?: any): Promise<T> {
    const response = await this.client.delete<T>(endpoint, config);
    return (response as any).data?.data || (response as any).data;
  }

  /**
   * Set auth token
   */
  setAuthToken(token: string): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem('authToken', token);
    }
  }

  /**
   * Clear auth token
   */
  clearAuthToken(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('authToken');
    }
  }
}

export const apiClient = new APIClient();

