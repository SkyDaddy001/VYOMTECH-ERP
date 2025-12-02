import { apiClient } from './api'
import { CustomerInteraction, DocumentTracker, SnagList, ChangeRequest, PostSalesMetrics, KPIDashboard } from '@/types/postsales'

export const postSalesService = {
  // Customer Interactions
  async getInteractions(bookingId?: string): Promise<CustomerInteraction[]> {
    if (bookingId) {
      return apiClient.get<CustomerInteraction[]>(`/api/v1/postsales/interactions?booking_id=${bookingId}`)
    }
    return apiClient.get<CustomerInteraction[]>('/api/v1/postsales/interactions')
  },

  async getInteraction(id: string): Promise<CustomerInteraction> {
    return apiClient.get<CustomerInteraction>(`/api/v1/postsales/interactions/${id}`)
  },

  async createInteraction(data: Partial<CustomerInteraction>): Promise<CustomerInteraction> {
    return apiClient.post<CustomerInteraction>('/api/v1/postsales/interactions', data)
  },

  async updateInteraction(id: string, data: Partial<CustomerInteraction>): Promise<CustomerInteraction> {
    return apiClient.put<CustomerInteraction>(`/api/v1/postsales/interactions/${id}`, data)
  },

  async deleteInteraction(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/postsales/interactions/${id}`)
  },

  async updateInteractionStatus(id: string, status: string): Promise<CustomerInteraction> {
    return apiClient.put<CustomerInteraction>(`/api/v1/postsales/interactions/${id}/status`, { status })
  },

  // Document Tracker
  async getDocuments(bookingId?: string): Promise<DocumentTracker[]> {
    if (bookingId) {
      return apiClient.get<DocumentTracker[]>(`/api/v1/postsales/documents?booking_id=${bookingId}`)
    }
    return apiClient.get<DocumentTracker[]>('/api/v1/postsales/documents')
  },

  async createDocument(data: Partial<DocumentTracker>): Promise<DocumentTracker> {
    return apiClient.post<DocumentTracker>('/api/v1/postsales/documents', data)
  },

  async updateDocument(id: string, data: Partial<DocumentTracker>): Promise<DocumentTracker> {
    return apiClient.put<DocumentTracker>(`/api/v1/postsales/documents/${id}`, data)
  },

  async deleteDocument(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/postsales/documents/${id}`)
  },

  async updateDocumentStatus(id: string, status: string): Promise<DocumentTracker> {
    return apiClient.put<DocumentTracker>(`/api/v1/postsales/documents/${id}/status`, { status })
  },

  // Snag List
  async getSnags(bookingId?: string): Promise<SnagList[]> {
    if (bookingId) {
      return apiClient.get<SnagList[]>(`/api/v1/postsales/snags?booking_id=${bookingId}`)
    }
    return apiClient.get<SnagList[]>('/api/v1/postsales/snags')
  },

  async createSnag(data: Partial<SnagList>): Promise<SnagList> {
    return apiClient.post<SnagList>('/api/v1/postsales/snags', data)
  },

  async updateSnag(id: string, data: Partial<SnagList>): Promise<SnagList> {
    return apiClient.put<SnagList>(`/api/v1/postsales/snags/${id}`, data)
  },

  async deleteSnag(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/postsales/snags/${id}`)
  },

  async updateSnagStatus(id: string, status: string): Promise<SnagList> {
    return apiClient.put<SnagList>(`/api/v1/postsales/snags/${id}/status`, { status })
  },

  // Change Requests
  async getChangeRequests(bookingId?: string): Promise<ChangeRequest[]> {
    if (bookingId) {
      return apiClient.get<ChangeRequest[]>(`/api/v1/postsales/change-requests?booking_id=${bookingId}`)
    }
    return apiClient.get<ChangeRequest[]>('/api/v1/postsales/change-requests')
  },

  async createChangeRequest(data: Partial<ChangeRequest>): Promise<ChangeRequest> {
    return apiClient.post<ChangeRequest>('/api/v1/postsales/change-requests', data)
  },

  async updateChangeRequest(id: string, data: Partial<ChangeRequest>): Promise<ChangeRequest> {
    return apiClient.put<ChangeRequest>(`/api/v1/postsales/change-requests/${id}`, data)
  },

  async deleteChangeRequest(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/postsales/change-requests/${id}`)
  },

  // Metrics & KPI
  async getMetrics(): Promise<PostSalesMetrics> {
    return apiClient.get<PostSalesMetrics>('/api/v1/postsales/metrics/overview')
  },

  async getKPIDashboard(period: string): Promise<KPIDashboard> {
    return apiClient.get<KPIDashboard>(`/api/v1/postsales/kpi/${period}`)
  },
}
