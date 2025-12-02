import { apiClient } from './api'
import {
  Campaign,
  Lead,
  MarketingMetrics,
  CampaignIdea,
  UTMTracker,
  LeadSourceTracking,
  MarketingTag,
  MarketingCampaignPerformance,
  SourceMetrics,
  SourceIntegration,
  GoogleIntegration,
  MetaIntegration,
  PortalIntegration,
  WebsiteIntegration,
} from '@/types/marketing'

export const marketingService = {
  // ========== CAMPAIGN IDEAS ==========
  async getIdeas(): Promise<CampaignIdea[]> {
    return apiClient.get<CampaignIdea[]>('/api/v1/marketing/ideas')
  },

  async getIdea(id: string): Promise<CampaignIdea> {
    return apiClient.get<CampaignIdea>(`/api/v1/marketing/ideas/${id}`)
  },

  async createIdea(data: Partial<CampaignIdea>): Promise<CampaignIdea> {
    return apiClient.post<CampaignIdea>('/api/v1/marketing/ideas', data)
  },

  async updateIdea(id: string, data: Partial<CampaignIdea>): Promise<CampaignIdea> {
    return apiClient.put<CampaignIdea>(`/api/v1/marketing/ideas/${id}`, data)
  },

  async deleteIdea(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/marketing/ideas/${id}`)
  },

  async updateIdeaStatus(id: string, status: string): Promise<CampaignIdea> {
    return apiClient.put<CampaignIdea>(`/api/v1/marketing/ideas/${id}/status`, { status })
  },

  async updateIdeaTags(id: string, tags: string[]): Promise<CampaignIdea> {
    return apiClient.put<CampaignIdea>(`/api/v1/marketing/ideas/${id}/tags`, { tags })
  },

  // ========== CAMPAIGNS ==========
  async getCampaigns(): Promise<Campaign[]> {
    return apiClient.get<Campaign[]>('/api/v1/marketing/campaigns')
  },

  async getCampaign(id: string): Promise<Campaign> {
    return apiClient.get<Campaign>(`/api/v1/marketing/campaigns/${id}`)
  },

  async createCampaign(campaign: Partial<Campaign>): Promise<Campaign> {
    return apiClient.post<Campaign>('/api/v1/marketing/campaigns', campaign)
  },

  async updateCampaign(id: string, campaign: Partial<Campaign>): Promise<Campaign> {
    return apiClient.put<Campaign>(`/api/v1/marketing/campaigns/${id}`, campaign)
  },

  async deleteCampaign(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/marketing/campaigns/${id}`)
  },

  // ========== UTM TRACKING ==========
  async createUTMTracker(data: Partial<UTMTracker>): Promise<UTMTracker> {
    return apiClient.post<UTMTracker>('/api/v1/marketing/utm-trackers', data)
  },

  async getUTMTrackers(campaignId?: string): Promise<UTMTracker[]> {
    if (campaignId) {
      return apiClient.get<UTMTracker[]>(`/api/v1/marketing/campaigns/${campaignId}/utm`)
    }
    return apiClient.get<UTMTracker[]>('/api/v1/marketing/utm-trackers')
  },

  async updateUTMTracker(id: string, data: Partial<UTMTracker>): Promise<UTMTracker> {
    return apiClient.put<UTMTracker>(`/api/v1/marketing/utm-trackers/${id}`, data)
  },

  async deleteUTMTracker(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/marketing/utm-trackers/${id}`)
  },

  async generateQRCode(id: string): Promise<{ qr_code: string }> {
    return apiClient.post<{ qr_code: string }>(`/api/v1/marketing/utm-trackers/${id}/qr`, {})
  },

  async shortenURL(id: string): Promise<{ short_url: string }> {
    return apiClient.post<{ short_url: string }>(`/api/v1/marketing/utm-trackers/${id}/shorten`, {})
  },

  // ========== LEAD SOURCE TRACKING ==========
  async setLeadSource(leadId: string, data: Partial<LeadSourceTracking>): Promise<LeadSourceTracking> {
    return apiClient.post<LeadSourceTracking>(`/api/v1/marketing/lead-source/${leadId}`, data)
  },

  async getLeadSource(leadId: string): Promise<LeadSourceTracking> {
    return apiClient.get<LeadSourceTracking>(`/api/v1/marketing/lead-source/${leadId}`)
  },

  async getLeadSourcesByCampaign(campaignId: string): Promise<LeadSourceTracking[]> {
    return apiClient.get<LeadSourceTracking[]>(
      `/api/v1/marketing/lead-sources/by-campaign/${campaignId}`
    )
  },

  // ========== MARKETING TAGS ==========
  async getTags(): Promise<MarketingTag[]> {
    return apiClient.get<MarketingTag[]>('/api/v1/marketing/tags')
  },

  async createTag(data: Partial<MarketingTag>): Promise<MarketingTag> {
    return apiClient.post<MarketingTag>('/api/v1/marketing/tags', data)
  },

  async updateTag(id: string, data: Partial<MarketingTag>): Promise<MarketingTag> {
    return apiClient.put<MarketingTag>(`/api/v1/marketing/tags/${id}`, data)
  },

  async deleteTag(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/marketing/tags/${id}`)
  },

  // ========== LEADS ==========
  async getLeads(status?: string): Promise<Lead[]> {
    const url = status ? `/api/v1/marketing/leads?status=${status}` : '/api/v1/marketing/leads'
    return apiClient.get<Lead[]>(url)
  },

  async getLead(id: string): Promise<Lead> {
    return apiClient.get<Lead>(`/api/v1/marketing/leads/${id}`)
  },

  async createLead(lead: Partial<Lead>): Promise<Lead> {
    return apiClient.post<Lead>('/api/v1/marketing/leads', lead)
  },

  async updateLead(id: string, lead: Partial<Lead>): Promise<Lead> {
    return apiClient.put<Lead>(`/api/v1/marketing/leads/${id}`, lead)
  },

  async deleteLead(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/marketing/leads/${id}`)
  },

  async qualifyLead(id: string): Promise<Lead> {
    return apiClient.post<Lead>(`/api/v1/marketing/leads/${id}/qualify`)
  },

  async convertLead(id: string): Promise<Lead> {
    return apiClient.post<Lead>(`/api/v1/marketing/leads/${id}/convert`)
  },

  // ========== ANALYTICS ==========
  async getMetrics(): Promise<MarketingMetrics> {
    return apiClient.get<MarketingMetrics>('/api/v1/marketing/metrics')
  },

  async getCampaignMetrics(campaignId: string): Promise<MarketingCampaignPerformance> {
    return apiClient.get<MarketingCampaignPerformance>(`/api/v1/marketing/campaigns/${campaignId}/metrics`)
  },

  async getCampaignAnalytics(campaignId: string): Promise<MarketingCampaignPerformance> {
    return apiClient.get<MarketingCampaignPerformance>(`/api/v1/marketing/analytics/${campaignId}`)
  },

  async getAnalyticsBySource(campaignId: string): Promise<Record<string, SourceMetrics>> {
    return apiClient.get<Record<string, SourceMetrics>>(
      `/api/v1/marketing/analytics/${campaignId}/by-source`
    )
  },

  async getAnalyticsByTag(campaignId: string): Promise<Record<string, SourceMetrics>> {
    return apiClient.get<Record<string, SourceMetrics>>(
      `/api/v1/marketing/analytics/${campaignId}/by-tag`
    )
  },

  // ========== INTEGRATIONS ==========
  async getIntegrations(): Promise<SourceIntegration[]> {
    return apiClient.get<SourceIntegration[]>('/api/v1/marketing/integrations')
  },

  async getIntegration(id: string): Promise<SourceIntegration> {
    return apiClient.get<SourceIntegration>(`/api/v1/marketing/integrations/${id}`)
  },

  async createIntegration(data: Partial<SourceIntegration>): Promise<SourceIntegration> {
    return apiClient.post<SourceIntegration>('/api/v1/marketing/integrations', data)
  },

  async updateIntegration(id: string, data: Partial<SourceIntegration>): Promise<SourceIntegration> {
    return apiClient.put<SourceIntegration>(`/api/v1/marketing/integrations/${id}`, data)
  },

  async deleteIntegration(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/marketing/integrations/${id}`)
  },

  async testIntegrationSync(id: string): Promise<{ status: string; leads_synced: number }> {
    return apiClient.post<{ status: string; leads_synced: number }>(
      `/api/v1/marketing/integrations/${id}/test-sync`,
      {}
    )
  },

  async toggleIntegrationActive(id: string, is_active: boolean): Promise<SourceIntegration> {
    return apiClient.put<SourceIntegration>(`/api/v1/marketing/integrations/${id}/toggle`, {
      is_active,
    })
  },

  // ========== INTEGRATION SPECIFIC ==========
  async createGoogleIntegration(data: Partial<GoogleIntegration>): Promise<GoogleIntegration> {
    return apiClient.post<GoogleIntegration>('/api/v1/marketing/integrations/google', data)
  },

  async createMetaIntegration(data: Partial<MetaIntegration>): Promise<MetaIntegration> {
    return apiClient.post<MetaIntegration>('/api/v1/marketing/integrations/meta', data)
  },

  async createPortalIntegration(data: Partial<PortalIntegration>): Promise<PortalIntegration> {
    return apiClient.post<PortalIntegration>('/api/v1/marketing/integrations/portal', data)
  },

  async createWebsiteIntegration(data: Partial<WebsiteIntegration>): Promise<WebsiteIntegration> {
    return apiClient.post<WebsiteIntegration>('/api/v1/marketing/integrations/website', data)
  },

  // ========== MULTI-CHANNEL SOURCE TRACKING ==========
  async getLeadsBySource(source: string): Promise<Lead[]> {
    return apiClient.get<Lead[]>(`/api/v1/marketing/leads/by-source/${source}`)
  },

  async getLeadsByPortal(portalName: string): Promise<Lead[]> {
    return apiClient.get<Lead[]>(`/api/v1/marketing/leads/by-portal/${portalName}`)
  },

  async getSourceAnalytics(campaignId?: string): Promise<Record<string, any>> {
    const url = campaignId
      ? `/api/v1/marketing/analytics/sources?campaign_id=${campaignId}`
      : '/api/v1/marketing/analytics/sources'
    return apiClient.get<Record<string, any>>(url)
  },

  // ========== UTM HELPERS ==========
  generateUTMParameters(source: string, medium: string, campaign: string, content?: string, term?: string) {
    return {
      utm_source: source,
      utm_medium: medium,
      utm_campaign: campaign,
      utm_content: content || '',
      utm_term: term || '',
    }
  },

  buildTrackingURL(baseURL: string, utmParams: Partial<UTMTracker>): string {
    const params = new URLSearchParams()
    if (utmParams.utm_source) params.append('utm_source', utmParams.utm_source)
    if (utmParams.utm_medium) params.append('utm_medium', utmParams.utm_medium)
    if (utmParams.utm_campaign) params.append('utm_campaign', utmParams.utm_campaign)
    if (utmParams.utm_content) params.append('utm_content', utmParams.utm_content)
    if (utmParams.utm_term) params.append('utm_term', utmParams.utm_term)

    return `${baseURL}?${params.toString()}`
  },

  parseUTMFromURL(url: string): Partial<UTMTracker> {
    try {
      const urlObj = new URL(url)
      return {
        utm_source: urlObj.searchParams.get('utm_source') || undefined,
        utm_medium: urlObj.searchParams.get('utm_medium') || undefined,
        utm_campaign: urlObj.searchParams.get('utm_campaign') || undefined,
        utm_content: urlObj.searchParams.get('utm_content') || undefined,
        utm_term: urlObj.searchParams.get('utm_term') || undefined,
      }
    } catch {
      return {}
    }
  },
}
