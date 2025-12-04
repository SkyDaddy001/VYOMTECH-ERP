import { apiClient } from '@/lib/api-client'

export interface DashboardStats {
  totalLeads: number
  activeLeads: number
  totalCalls: number
  successfulCalls: number
  totalRevenue: number
  conversionRate: number
  agents: number
  campaigns: number
}

export interface Lead {
  id: string
  name: string
  email: string
  phone: string
  status: string
  source: string
  score: number
  createdAt: string
  updatedAt: string
}

export interface Campaign {
  id: string
  name: string
  status: string
  startDate: string
  endDate: string
  leads: number
  budget: number
  spent: number
  roi: number
}

export interface Call {
  id: string
  agentId: string
  leadId: string
  duration: number
  status: string
  recordingUrl?: string
  createdAt: string
}

export interface Agent {
  id: string
  name: string
  email: string
  status: string
  totalCalls: number
  successfulCalls: number
  rating: number
}

export interface Gamification {
  userId: string
  points: number
  level: number
  badges: string[]
  leaderboardPosition: number
  weeklyPoints: number
}

const dashboardService = {
  // Dashboard Stats
  getDashboardStats: async () => {
    try {
      const response = await apiClient.get<DashboardStats>('/api/v1/dashboard/stats')
      return response
    } catch (error) {
      console.error('Failed to fetch dashboard stats:', error)
      throw error
    }
  },

  // Leads
  getLeads: async (params?: any) => {
    try {
      const response = await apiClient.get<Lead[]>('/api/v1/leads', { params })
      return response
    } catch (error) {
      console.error('Failed to fetch leads:', error)
      throw error
    }
  },

  getLead: async (id: string) => {
    try {
      const response = await apiClient.get<Lead>(`/api/v1/leads/${id}`)
      return response
    } catch (error) {
      console.error('Failed to fetch lead:', error)
      throw error
    }
  },

  createLead: async (data: any) => {
    try {
      const response = await apiClient.post<Lead>('/api/v1/leads', data)
      return response
    } catch (error) {
      console.error('Failed to create lead:', error)
      throw error
    }
  },

  updateLead: async (id: string, data: any) => {
    try {
      const response = await apiClient.put<Lead>(`/api/v1/leads/${id}`, data)
      return response
    } catch (error) {
      console.error('Failed to update lead:', error)
      throw error
    }
  },

  // Campaigns
  getCampaigns: async (params?: any) => {
    try {
      const response = await apiClient.get<Campaign[]>('/api/v1/campaigns', { params })
      return response
    } catch (error) {
      console.error('Failed to fetch campaigns:', error)
      throw error
    }
  },

  getCampaign: async (id: string) => {
    try {
      const response = await apiClient.get<Campaign>(`/api/v1/campaigns/${id}`)
      return response
    } catch (error) {
      console.error('Failed to fetch campaign:', error)
      throw error
    }
  },

  createCampaign: async (data: any) => {
    try {
      const response = await apiClient.post<Campaign>('/api/v1/campaigns', data)
      return response
    } catch (error) {
      console.error('Failed to create campaign:', error)
      throw error
    }
  },

  // Calls
  getCalls: async (params?: any) => {
    try {
      const response = await apiClient.get<Call[]>('/api/v1/calls', { params })
      return response
    } catch (error) {
      console.error('Failed to fetch calls:', error)
      throw error
    }
  },

  // Agents
  getAgents: async () => {
    try {
      const response = await apiClient.get<Agent[]>('/api/v1/agents')
      return response
    } catch (error) {
      console.error('Failed to fetch agents:', error)
      throw error
    }
  },

  // Gamification
  getGamificationStats: async (userId: string) => {
    try {
      const response = await apiClient.get<Gamification>(
        `/api/v1/gamification/${userId}`
      )
      return response
    } catch (error) {
      console.error('Failed to fetch gamification stats:', error)
      throw error
    }
  },
}

export default dashboardService
