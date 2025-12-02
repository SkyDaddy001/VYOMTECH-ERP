import { apiClient } from './api'
import { PreSalesOpportunity, Proposal, PreSalesMetrics, Pipeline } from '@/types/presales'

export const preSalesService = {
  // Opportunity operations
  async getOpportunities(stage?: string): Promise<PreSalesOpportunity[]> {
    const url = stage ? `/api/v1/presales/opportunities?stage=${stage}` : '/api/v1/presales/opportunities'
    return apiClient.get<PreSalesOpportunity[]>(url)
  },

  async getOpportunity(id: string): Promise<PreSalesOpportunity> {
    return apiClient.get<PreSalesOpportunity>(`/api/v1/presales/opportunities/${id}`)
  },

  async createOpportunity(opportunity: Partial<PreSalesOpportunity>): Promise<PreSalesOpportunity> {
    return apiClient.post<PreSalesOpportunity>('/api/v1/presales/opportunities', opportunity)
  },

  async updateOpportunity(id: string, opportunity: Partial<PreSalesOpportunity>): Promise<PreSalesOpportunity> {
    return apiClient.put<PreSalesOpportunity>(`/api/v1/presales/opportunities/${id}`, opportunity)
  },

  async deleteOpportunity(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/presales/opportunities/${id}`)
  },

  async moveOpportunity(id: string, stage: string): Promise<PreSalesOpportunity> {
    return apiClient.post<PreSalesOpportunity>(`/api/v1/presales/opportunities/${id}/move`, { stage })
  },

  // Proposal operations
  async getProposals(opportunityId?: string): Promise<Proposal[]> {
    const url = opportunityId ? `/api/v1/presales/proposals?opportunity_id=${opportunityId}` : '/api/v1/presales/proposals'
    return apiClient.get<Proposal[]>(url)
  },

  async getProposal(id: string): Promise<Proposal> {
    return apiClient.get<Proposal>(`/api/v1/presales/proposals/${id}`)
  },

  async createProposal(proposal: Partial<Proposal>): Promise<Proposal> {
    return apiClient.post<Proposal>('/api/v1/presales/proposals', proposal)
  },

  async updateProposal(id: string, proposal: Partial<Proposal>): Promise<Proposal> {
    return apiClient.put<Proposal>(`/api/v1/presales/proposals/${id}`, proposal)
  },

  async deleteProposal(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/presales/proposals/${id}`)
  },

  async sendProposal(id: string): Promise<Proposal> {
    return apiClient.post<Proposal>(`/api/v1/presales/proposals/${id}/send`)
  },

  async approveProposal(id: string): Promise<Proposal> {
    return apiClient.post<Proposal>(`/api/v1/presales/proposals/${id}/approve`)
  },

  // Pre-Sales metrics
  async getMetrics(): Promise<PreSalesMetrics> {
    return apiClient.get<PreSalesMetrics>('/api/v1/presales/metrics')
  },

  async getPipeline(): Promise<Pipeline[]> {
    return apiClient.get<Pipeline[]>('/api/v1/presales/pipeline')
  },
}
