import { apiClient } from './api-client'

/**
 * Enhanced API Service Layer for Frontend
 * Provides type-safe methods for all backend endpoints
 */

// ==================== LEADS ====================
export const leadsApi = {
  list: async (params?: { limit?: number; offset?: number; status?: string }) => {
    return apiClient.get('/api/v1/leads', { params })
  },
  
  get: async (id: string | number) => {
    return apiClient.get(`/api/v1/leads/${id}`)
  },
  
  create: async (data: any) => {
    return apiClient.post('/api/v1/leads', data)
  },
  
  update: async (id: string | number, data: any) => {
    return apiClient.put(`/api/v1/leads/${id}`, data)
  },
  
  delete: async (id: string | number) => {
    return apiClient.delete(`/api/v1/leads/${id}`)
  },

  search: async (query: string) => {
    return apiClient.get('/api/v1/leads/search', { params: { q: query } })
  },
}

// ==================== AGENTS ====================
export const agentsApi = {
  list: async (params?: { limit?: number; offset?: number }) => {
    return apiClient.get('/api/v1/agents', { params })
  },
  
  get: async (id: string | number) => {
    return apiClient.get(`/api/v1/agents/${id}`)
  },
  
  create: async (data: any) => {
    return apiClient.post('/api/v1/agents', data)
  },
  
  update: async (id: string | number, data: any) => {
    return apiClient.put(`/api/v1/agents/${id}`, data)
  },
  
  delete: async (id: string | number) => {
    return apiClient.delete(`/api/v1/agents/${id}`)
  },

  getStats: async (id: string | number) => {
    return apiClient.get(`/api/v1/agents/${id}/stats`)
  },

  updateAvailability: async (id: string | number, availability: string) => {
    return apiClient.put(`/api/v1/agents/${id}/availability`, { availability })
  },
}

// ==================== CALLS ====================
export const callsApi = {
  list: async (params?: { limit?: number; offset?: number; status?: string }) => {
    return apiClient.get('/api/v1/calls', { params })
  },
  
  get: async (id: string | number) => {
    return apiClient.get(`/api/v1/calls/${id}`)
  },
  
  create: async (data: any) => {
    return apiClient.post('/api/v1/calls', data)
  },
  
  update: async (id: string | number, data: any) => {
    return apiClient.put(`/api/v1/calls/${id}`, data)
  },

  getRecording: async (id: string | number) => {
    return apiClient.get(`/api/v1/calls/${id}/recording`)
  },

  getTranscription: async (id: string | number) => {
    return apiClient.get(`/api/v1/calls/${id}/transcription`)
  },
}

// ==================== CAMPAIGNS ====================
export const campaignsApi = {
  list: async (params?: { limit?: number; offset?: number; status?: string }) => {
    return apiClient.get('/api/v1/campaigns', { params })
  },
  
  get: async (id: string | number) => {
    return apiClient.get(`/api/v1/campaigns/${id}`)
  },
  
  create: async (data: any) => {
    return apiClient.post('/api/v1/campaigns', data)
  },
  
  update: async (id: string | number, data: any) => {
    return apiClient.put(`/api/v1/campaigns/${id}`, data)
  },
  
  delete: async (id: string | number) => {
    return apiClient.delete(`/api/v1/campaigns/${id}`)
  },

  getStats: async (id: string | number) => {
    return apiClient.get(`/api/v1/campaigns/${id}/stats`)
  },
}

// ==================== DASHBOARD ====================
export const dashboardApi = {
  getStats: async () => {
    return apiClient.get('/api/v1/dashboard/stats')
  },

  getMetrics: async () => {
    return apiClient.get('/api/v1/dashboard/metrics')
  },

  getCharts: async (period?: string) => {
    const params = period ? { period } : {}
    return apiClient.get('/api/v1/dashboard/charts', { params })
  },
}

// ==================== REAL ESTATE ====================
export const realEstateApi = {
  listProperties: async (params?: { limit?: number; offset?: number; status?: string }) => {
    return apiClient.get('/api/v1/real-estate/properties', { params })
  },
  
  getProperty: async (id: string | number) => {
    return apiClient.get(`/api/v1/real-estate/properties/${id}`)
  },
  
  createProperty: async (data: any) => {
    return apiClient.post('/api/v1/real-estate/properties', data)
  },
  
  updateProperty: async (id: string | number, data: any) => {
    return apiClient.put(`/api/v1/real-estate/properties/${id}`, data)
  },

  getProjectStats: async () => {
    return apiClient.get('/api/v1/real-estate/summary')
  },

  listBrokers: async (params?: { limit?: number; offset?: number }) => {
    return apiClient.get('/api/v1/real-estate/brokers', { params })
  },
}

// ==================== SALES ====================
export const salesApi = {
  listOpportunities: async (params?: { limit?: number; offset?: number; stage?: string }) => {
    return apiClient.get('/api/v1/sales/opportunities', { params })
  },
  
  getOpportunity: async (id: string | number) => {
    return apiClient.get(`/api/v1/sales/opportunities/${id}`)
  },
  
  createOpportunity: async (data: any) => {
    return apiClient.post('/api/v1/sales/opportunities', data)
  },
  
  updateOpportunity: async (id: string | number, data: any) => {
    return apiClient.put(`/api/v1/sales/opportunities/${id}`, data)
  },

  getPipeline: async () => {
    return apiClient.get('/api/v1/sales/pipeline')
  },

  getConversionMetrics: async () => {
    return apiClient.get('/api/v1/sales/conversion-metrics')
  },
}

// ==================== FINANCE ====================
export const financeApi = {
  getGLEntries: async (params?: { account_code?: string; start_date?: string; end_date?: string; limit?: number; offset?: number }) => {
    return apiClient.get('/api/v1/gl/entries', { params })
  },
  
  getTrialBalance: async (asOfDate?: string) => {
    const params = asOfDate ? { as_of_date: asOfDate } : {}
    return apiClient.get('/api/v1/gl/trial-balance', { params })
  },
  
  getVouchers: async (params?: { type?: string; start_date?: string; end_date?: string; limit?: number; offset?: number }) => {
    return apiClient.get('/api/v1/gl/vouchers', { params })
  },
}

// ==================== INTEGRATION ====================
export const integrationApi = {
  listProviders: async (params?: { limit?: number; offset?: number }) => {
    return apiClient.get('/api/v1/integration/providers', { params })
  },
  
  getProvider: async (id: string | number) => {
    return apiClient.get(`/api/v1/integration/providers/${id}`)
  },
  
  createProvider: async (data: any) => {
    return apiClient.post('/api/v1/integration/providers', data)
  },
  
  updateProvider: async (id: string | number, data: any) => {
    return apiClient.put(`/api/v1/integration/providers/${id}`, data)
  },
  
  deleteProvider: async (id: string | number) => {
    return apiClient.delete(`/api/v1/integration/providers/${id}`)
  },

  getStats: async () => {
    return apiClient.get('/api/v1/integration/stats')
  },

  triggerSync: async (data: any) => {
    return apiClient.post('/api/v1/integration/sync', data)
  },
}

// ==================== EXPORT ALL ====================
export const api = {
  leads: leadsApi,
  agents: agentsApi,
  calls: callsApi,
  campaigns: campaignsApi,
  dashboard: dashboardApi,
  realEstate: realEstateApi,
  sales: salesApi,
  finance: financeApi,
  integration: integrationApi,
}

export default api
