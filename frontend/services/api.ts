import axios, { AxiosInstance, AxiosError } from 'axios'
import { AuthResponse, LoginRequest, RegisterRequest, User } from '@/types'

// Custom API Error class for better error handling
export class ApiError extends Error {
  constructor(
    public status: number | null,
    public code: string,
    public userMessage: string,
    message: string
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

// Auth event emitter for handling 401 responses
class AuthEventEmitter {
  private listeners: Set<(event: string, data?: any) => void> = new Set()

  on(callback: (event: string, data?: any) => void) {
    this.listeners.add(callback)
  }

  off(callback: (event: string, data?: any) => void) {
    this.listeners.delete(callback)
  }

  emit(event: string, data?: any) {
    this.listeners.forEach(listener => listener(event, data))
  }
}

export const authEventEmitter = new AuthEventEmitter()

// Get API URL - resolves dynamically at runtime
const getApiUrl = (): string => {
  const configuredUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
  
  // In browser context, convert Docker network URL to localhost
  if (typeof window !== 'undefined' && configuredUrl.includes('http://app:')) {
    console.log('Converting Docker URL to localhost:', configuredUrl, '-> http://localhost:8080')
    return 'http://localhost:8080'
  }

  console.log('Using API URL:', configuredUrl)
  return configuredUrl
}

class ApiClient {
  private client: AxiosInstance
  private token: string | null = null
  private apiUrl: string

  constructor() {
    // Resolve API URL at instance creation time
    this.apiUrl = getApiUrl()
    
    this.client = axios.create({
      baseURL: this.apiUrl,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    // Load token from localStorage
    if (typeof window !== 'undefined') {
      this.token = localStorage.getItem('auth_token')
    }

    // Add token and tenant ID to requests
    this.client.interceptors.request.use((config) => {
      console.log('API Request:', config.method?.toUpperCase(), config.url)
      if (this.token) {
        config.headers.Authorization = `Bearer ${this.token}`
      }
      
      // Add X-Tenant-ID header from localStorage
      if (typeof window !== 'undefined') {
        const user = localStorage.getItem('user')
        if (user) {
          try {
            const userData = JSON.parse(user)
            if (userData.tenant_id) {
              config.headers['X-Tenant-ID'] = userData.tenant_id
              console.log('Added X-Tenant-ID header:', userData.tenant_id)
            }
          } catch (e) {
            console.warn('Failed to parse user data for tenant ID')
          }
        }
      }
      
      return config
    })

    // Handle errors
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        // Distinguish between network errors and HTTP errors
        if (!error.response) {
          // Network error or timeout
          const userMessage = 
            error.code === 'ECONNABORTED' ? 'Request timeout. Please try again.' :
            error.code === 'ERR_NETWORK' || error.message === 'Network Error' ? 'Unable to connect to server. Please check your connection and try again.' :
            error.code === 'ERR_CANCELED' ? 'Request was cancelled.' :
            'Network error. Please try again.'
          
          const apiError = new ApiError(
            null,
            error.code || 'NETWORK_ERROR',
            userMessage,
            error.message
          )
          console.error('Network Error:', error.code, error.message)
          return Promise.reject(apiError)
        }

        // HTTP error response
        const status = error.response.status
        console.error('HTTP Error:', status, error.response.data)

        if (status === 401) {
          console.warn('Received 401, clearing auth state and emitting logout event')
          localStorage.removeItem('auth_token')
          localStorage.removeItem('user')
          // Emit auth logout event to be caught by AuthProvider
          authEventEmitter.emit('logout', { reason: 'unauthorized' })
        }

        const userMessage = 
          status === 400 ? 'Invalid request. Please check your input.' :
          status === 401 ? 'Unauthorized. Please log in again.' :
          status === 403 ? 'Access denied.' :
          status === 404 ? 'Resource not found.' :
          status === 409 ? 'This resource already exists or there is a conflict.' :
          status === 422 ? 'Validation error. Please check your input.' :
          status >= 500 ? 'Server error. Please try again later.' :
          'An error occurred. Please try again.'

        const apiError = new ApiError(
          status,
          error.code || `HTTP_${status}`,
          userMessage,
          error.message
        )
        return Promise.reject(apiError)
      }
    )
  }

  setToken(token: string) {
    this.token = token
    localStorage.setItem('auth_token', token)
  }

  getToken() {
    return this.token
  }

  clearToken() {
    this.token = null
    localStorage.removeItem('auth_token')
  }

  async get<T>(url: string) {
    const response = await this.client.get<T>(url)
    return response.data
  }

  async post<T>(url: string, data?: any) {
    const response = await this.client.post<T>(url, data)
    return response.data
  }

  async put<T>(url: string, data?: any) {
    const response = await this.client.put<T>(url, data)
    return response.data
  }

  async delete<T>(url: string) {
    const response = await this.client.delete<T>(url)
    return response.data
  }
}

export const apiClient = new ApiClient()

// Auth Service
export const authService = {
  async login(email: string, password: string): Promise<AuthResponse> {
    const data = await apiClient.post<AuthResponse>('/api/v1/auth/login', {
      email,
      password,
    })
    if (data.token) {
      apiClient.setToken(data.token)
    }
    return data
  },

  async register(
    email: string,
    password: string,
    role: string,
    tenant_id: string,
    name?: string
  ): Promise<AuthResponse> {
    const data = await apiClient.post<AuthResponse>('/api/v1/auth/register', {
      email,
      password,
      role,
      tenant_id,
      name: name || email.split('@')[0], // Use first part of email as fallback
    })
    if (data.token) {
      apiClient.setToken(data.token)
    }
    return data
  },

  async validateToken(): Promise<{ valid: boolean; user_id: number }> {
    return apiClient.post('/api/v1/auth/validate', {})
  },

  logout() {
    apiClient.clearToken()
  },
}

// Agent Service
export const agentService = {
  async listAgents() {
    return apiClient.get('/api/v1/agents')
  },

  async getAgent(id: number) {
    return apiClient.get(`/api/v1/agents/${id}`)
  },

  async createAgent(data: any) {
    return apiClient.post('/api/v1/agents', data)
  },

  async updateAgent(id: number, data: any) {
    return apiClient.put(`/api/v1/agents/${id}`, data)
  },

  async updateAvailability(id: number, availability: string) {
    return apiClient.put(`/api/v1/agents/${id}/availability`, {
      availability,
    })
  },

  async getAgentStats(id: number) {
    return apiClient.get(`/api/v1/agents/${id}/stats`)
  },
}

// Tenant Service
export const tenantService = {
  async createTenant(name: string, domain: string = '') {
    return apiClient.post('/api/v1/tenants', {
      name,
      domain,
    })
  },

  async getTenantInfo() {
    return apiClient.get('/api/v1/tenant')
  },

  async getTenantUserCount() {
    return apiClient.get('/api/v1/tenant/users/count')
  },

  async listTenants() {
    return apiClient.get('/api/v1/tenants')
  },

  async getUserTenants() {
    return apiClient.get('/api/v1/tenants')
  },

  async switchTenant(tenantId: string) {
    return apiClient.post(`/api/v1/tenants/${tenantId}/switch`, {})
  },

  async addTenantMember(tenantId: string, email: string, role: string) {
    return apiClient.post(`/api/v1/tenants/${tenantId}/members`, {
      email,
      role,
    })
  },

  async removeTenantMember(tenantId: string, email: string) {
    return apiClient.delete(`/api/v1/tenants/${tenantId}/members/${email}`)
  },
}

// Admin Service - Tenant Management
export const adminTenantService = {
  async listTenants() {
    return apiClient.get('/api/v1/tenants')
  },

  async getTenant(tenantId: string) {
    return apiClient.get(`/api/v1/tenants/${tenantId}`)
  },

  async createTenant(data: {
    name: string
    domain: string
    status?: string
    max_users?: number
    max_concurrent_calls?: number
    ai_budget_monthly?: number
  }) {
    return apiClient.post('/api/v1/tenants', data)
  },

  async updateTenant(tenantId: string, data: any) {
    return apiClient.put(`/api/v1/tenants/${tenantId}`, data)
  },

  async deleteTenant(tenantId: string) {
    return apiClient.delete(`/api/v1/tenants/${tenantId}`)
  },

  async getTenantUsers(tenantId: string) {
    return apiClient.get(`/api/v1/tenants/${tenantId}/users`)
  },
}

// Admin Service - User Management
export const adminUserService = {
  async listUsers() {
    return apiClient.get('/api/v1/users')
  },

  async getUser(userId: number | string) {
    return apiClient.get(`/api/v1/users/${userId}`)
  },

  async createUser(data: {
    email: string
    password: string
    name?: string
    role: string
    tenant_id: string
  }) {
    return apiClient.post('/api/v1/users', data)
  },

  async updateUser(userId: number | string, data: any) {
    return apiClient.put(`/api/v1/users/${userId}`, data)
  },

  async deleteUser(userId: number | string) {
    return apiClient.delete(`/api/v1/users/${userId}`)
  },

  async updateUserRole(userId: number | string, role: string) {
    return apiClient.put(`/api/v1/users/${userId}/role`, { role })
  },

  async resetPassword(userId: number | string, newPassword: string) {
    return apiClient.post(`/api/v1/users/${userId}/reset-password`, { password: newPassword })
  },
}

// Lead Service - Feature: Lead Management
export const leadService = {
  async listLeads(page?: number, limit?: number) {
    const params = new URLSearchParams()
    if (page) params.append('page', page.toString())
    if (limit) params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/leads${params.size > 0 ? '?' + params : ''}`)
  },

  async getLead(id: number | string) {
    return apiClient.get(`/api/v1/leads?id=${id}`)
  },

  async createLead(data: any) {
    return apiClient.post('/api/v1/leads', data)
  },

  async updateLead(id: number | string, data: any) {
    return apiClient.put(`/api/v1/leads?id=${id}`, data)
  },

  async deleteLead(id: number | string) {
    return apiClient.delete(`/api/v1/leads?id=${id}`)
  },

  async getLeadStats() {
    return apiClient.get('/api/v1/leads/stats')
  },
}

// Call Service - Feature: Call Management
export const callService = {
  async listCalls(page?: number, limit?: number) {
    const params = new URLSearchParams()
    if (page) params.append('page', page.toString())
    if (limit) params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/calls${params.size > 0 ? '?' + params : ''}`)
  },

  async getCall(id: number | string) {
    return apiClient.get(`/api/v1/calls?id=${id}`)
  },

  async createCall(data: any) {
    return apiClient.post('/api/v1/calls', data)
  },

  async endCall(id: number | string) {
    return apiClient.post(`/api/v1/calls?id=${id}/end`, {})
  },

  async getCallStats() {
    return apiClient.get('/api/v1/calls/stats')
  },
}

// Campaign Service - Feature: Campaign Management
export const campaignService = {
  async listCampaigns(page?: number, limit?: number) {
    const params = new URLSearchParams()
    if (page) params.append('page', page.toString())
    if (limit) params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/campaigns${params.size > 0 ? '?' + params : ''}`)
  },

  async getCampaign(id: number | string) {
    return apiClient.get(`/api/v1/campaigns?id=${id}`)
  },

  async createCampaign(data: any) {
    return apiClient.post('/api/v1/campaigns', data)
  },

  async updateCampaign(id: number | string, data: any) {
    return apiClient.put(`/api/v1/campaigns?id=${id}`, data)
  },

  async deleteCampaign(id: number | string) {
    return apiClient.delete(`/api/v1/campaigns?id=${id}`)
  },

  async getCampaignStats() {
    return apiClient.get('/api/v1/campaigns/stats')
  },
}

// WebSocket Service - Feature 1: Real-time Communication
export const webSocketService = {
  getWebSocketUrl(): string {
    const apiUrl = getApiUrl()
    return apiUrl.replace('http://', 'ws://').replace('https://', 'wss://')
  },

  connect(token: string) {
    const wsUrl = `${this.getWebSocketUrl()}/api/v1/ws`
    return new WebSocket(wsUrl)
  },

  async getConnectionStats() {
    return apiClient.get('/api/v1/ws/stats')
  },
}

// Analytics Service - Feature 2: Advanced Analytics
export const analyticsService = {
  async generateReport(type: string, startDate: string, endDate: string) {
    return apiClient.post('/api/v1/analytics/reports', {
      type,
      start_date: startDate,
      end_date: endDate,
    })
  },

  async exportReport(reportId: string, format: 'csv' | 'json' | 'pdf' = 'json') {
    return apiClient.post('/api/v1/analytics/export', {
      report_id: reportId,
      format,
    })
  },

  async getTrends(metric: string, startDate: string, endDate: string) {
    const params = new URLSearchParams({
      metric,
      start_date: startDate,
      end_date: endDate,
    })
    return apiClient.get(`/api/v1/analytics/trends?${params}`)
  },

  async getCustomMetrics(metric: string, filters?: Record<string, any>) {
    const params = new URLSearchParams({ metric })
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        params.append(`filter_${key}`, String(value))
      })
    }
    return apiClient.get(`/api/v1/analytics/metrics?${params}`)
  },
}

// Automation Service - Feature 3: Automation & Routing
export const automationService = {
  async calculateLeadScore(leadId: number | string) {
    return apiClient.post('/api/v1/automation/leads/score', {
      lead_id: leadId,
    })
  },

  async rankLeads(limit: number = 100) {
    return apiClient.get(`/api/v1/automation/leads/ranked?limit=${limit}`)
  },

  async routeLeadToAgent(leadId: number | string) {
    return apiClient.post('/api/v1/automation/leads/route', {
      lead_id: leadId,
    })
  },

  async createRoutingRule(data: any) {
    return apiClient.post('/api/v1/automation/routing-rules', data)
  },

  async scheduleCampaign(campaignId: number | string, scheduledTime: string) {
    return apiClient.post('/api/v1/automation/schedule-campaign', {
      campaign_id: campaignId,
      scheduled_time: scheduledTime,
    })
  },

  async getLeadScoringMetrics() {
    return apiClient.get('/api/v1/automation/metrics')
  },
}

// Communication Service - Feature 4: Multi-channel Communication
export const communicationService = {
  async registerProvider(type: string, credentials: Record<string, any>) {
    return apiClient.post('/api/v1/communication/providers', {
      type,
      credentials,
    })
  },

  async createTemplate(name: string, type: string, content: string) {
    return apiClient.post('/api/v1/communication/templates', {
      name,
      type,
      content,
    })
  },

  async sendMessage(recipient: string, type: string, templateId?: string, body?: string) {
    return apiClient.post('/api/v1/communication/messages', {
      recipient,
      type,
      template_id: templateId,
      body,
    })
  },

  async getMessageStatus(messageId: string) {
    return apiClient.get(`/api/v1/communication/messages/status?id=${messageId}`)
  },

  async getMessageStats() {
    return apiClient.get('/api/v1/communication/stats')
  },
}

// Gamification Service - Feature 5: Advanced Gamification
export const gamificationService = {
  // Basic gamification endpoints
  async getUserPoints() {
    return apiClient.get('/api/v1/gamification/points')
  },

  async awardPoints(userId: number | string, points: number, reason: string) {
    return apiClient.post('/api/v1/gamification/points/award', {
      user_id: userId,
      points,
      reason,
    })
  },

  async revokePoints(userId: number | string, points: number, reason: string) {
    return apiClient.post('/api/v1/gamification/points/revoke', {
      user_id: userId,
      points,
      reason,
    })
  },

  async getUserBadges() {
    return apiClient.get('/api/v1/gamification/badges')
  },

  async createBadge(name: string, description: string, icon: string) {
    return apiClient.post('/api/v1/gamification/badges', {
      name,
      description,
      icon,
    })
  },

  async awardBadge(userId: number | string, badgeId: number | string) {
    return apiClient.post('/api/v1/gamification/badges/award', {
      user_id: userId,
      badge_id: badgeId,
    })
  },

  async getUserChallenges() {
    return apiClient.get('/api/v1/gamification/challenges')
  },

  async getActiveChallenges() {
    return apiClient.get('/api/v1/gamification/challenges/active')
  },

  async createChallenge(name: string, description: string, targetScore: number) {
    return apiClient.post('/api/v1/gamification/challenges', {
      name,
      description,
      target_score: targetScore,
    })
  },

  async getLeaderboard(limit: number = 50) {
    return apiClient.get(`/api/v1/gamification/leaderboard?limit=${limit}`)
  },

  async getGamificationProfile() {
    return apiClient.get('/api/v1/gamification/profile')
  },

  // Advanced gamification endpoints
  async createCompetition(name: string, description: string, startDate: string, endDate: string) {
    return apiClient.post('/api/v1/gamification-advanced/competitions', {
      name,
      description,
      start_date: startDate,
      end_date: endDate,
    })
  },

  async getTeamLeaderboard(competitionId: string) {
    return apiClient.get(`/api/v1/gamification-advanced/competitions/leaderboard?competition_id=${competitionId}`)
  },

  async createAdvancedChallenge(name: string, description: string, config: any) {
    return apiClient.post('/api/v1/gamification-advanced/challenges', {
      name,
      description,
      config,
    })
  },

  async getAvailableRewards() {
    return apiClient.get('/api/v1/gamification-advanced/rewards')
  },

  async createReward(name: string, pointsCost: number, description: string) {
    return apiClient.post('/api/v1/gamification-advanced/rewards', {
      name,
      points_cost: pointsCost,
      description,
    })
  },

  async redeemReward(rewardId: number | string) {
    return apiClient.post('/api/v1/gamification-advanced/redeem', {
      reward_id: rewardId,
    })
  },

  async getAdvancedLeaderboard() {
    return apiClient.get('/api/v1/gamification-advanced/leaderboard')
  },

  async getGamificationStats() {
    return apiClient.get('/api/v1/gamification-advanced/stats')
  },
}

// Compliance & Security Service - Feature 6: RBAC, Audit, Encryption, GDPR
export const complianceService = {
  // RBAC endpoints
  async createRole(name: string, description: string) {
    return apiClient.post('/api/v1/compliance/roles', {
      name,
      description,
    })
  },

  async getRoles() {
    return apiClient.get('/api/v1/compliance/roles')
  },

  // Audit endpoints
  async getAuditLogs(userId?: string, action?: string, limit: number = 50, offset: number = 0) {
    const params = new URLSearchParams()
    if (userId) params.append('user_id', userId)
    if (action) params.append('action', action)
    params.append('limit', limit.toString())
    params.append('offset', offset.toString())
    return apiClient.get(`/api/v1/compliance/audit-logs?${params}`)
  },

  async getAuditSummary(startDate?: string, endDate?: string) {
    const params = new URLSearchParams()
    if (startDate) params.append('start_date', startDate)
    if (endDate) params.append('end_date', endDate)
    return apiClient.get(`/api/v1/compliance/audit-summary${params.size > 0 ? '?' + params : ''}`)
  },

  async getSecurityEvents(status?: string, limit: number = 50) {
    const params = new URLSearchParams()
    if (status) params.append('status', status)
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/compliance/security-events?${params}`)
  },

  async getComplianceReport() {
    return apiClient.get('/api/v1/compliance/report')
  },

  // GDPR endpoints
  async requestDataAccess() {
    return apiClient.post('/api/v1/compliance/gdpr/request-access', {})
  },

  async exportUserData() {
    return apiClient.post('/api/v1/compliance/gdpr/export', {})
  },

  async requestDataDeletion(reason: string) {
    return apiClient.post('/api/v1/compliance/gdpr/request-deletion', {
      reason,
    })
  },

  async getUserConsents() {
    return apiClient.get('/api/v1/compliance/gdpr/consents')
  },

  async recordConsent(type: string, consentValue: boolean) {
    return apiClient.post('/api/v1/compliance/gdpr/consents', {
      type,
      value: consentValue,
    })
  },
}

// AI Service - AI Query Processing
export const aiService = {
  async processQuery(query: string, context?: Record<string, any>) {
    return apiClient.post('/api/v1/ai/query', {
      query,
      context,
    })
  },

  async listProviders() {
    return apiClient.get('/api/v1/ai/providers')
  },
}

// ============================================================================
// DASHBOARD SERVICES - Phase 6: Real Data Integration
// ============================================================================

// Financial Dashboard Service
export const financialDashboardService = {
  async getProfitAndLoss(startDate?: Date, endDate?: Date) {
    const params = new URLSearchParams()
    if (startDate) params.append('start_date', startDate.toISOString())
    if (endDate) params.append('end_date', endDate.toISOString())
    return apiClient.post('/api/v1/dashboard/profit-and-loss', {
      start_date: startDate,
      end_date: endDate,
    })
  },

  async getBalanceSheet(asOfDate?: Date) {
    return apiClient.post('/api/v1/dashboard/balance-sheet', {
      as_of_date: asOfDate,
    })
  },

  async getCashFlow(startDate?: Date, endDate?: Date) {
    return apiClient.post('/api/v1/dashboard/cash-flow', {
      start_date: startDate,
      end_date: endDate,
    })
  },

  async getFinancialRatios() {
    return apiClient.get('/api/v1/dashboard/ratios')
  },
}

// Sales Dashboard Service
export const salesDashboardService = {
  async getSalesOverview(startDate?: Date, endDate?: Date) {
    const params = new URLSearchParams()
    if (startDate) params.append('start_date', startDate.toISOString())
    if (endDate) params.append('end_date', endDate.toISOString())
    return apiClient.get(`/api/v1/dashboard/sales/overview${params.size > 0 ? '?' + params : ''}`)
  },

  async getPipelineAnalysis() {
    return apiClient.get('/api/v1/dashboard/sales/pipeline')
  },

  async getSalesMetrics(startDate?: Date, endDate?: Date) {
    const params = new URLSearchParams()
    if (startDate) params.append('start_date', startDate.toISOString())
    if (endDate) params.append('end_date', endDate.toISOString())
    return apiClient.get(`/api/v1/dashboard/sales/metrics${params.size > 0 ? '?' + params : ''}`)
  },

  async getInvoiceStatus() {
    return apiClient.get('/api/v1/dashboard/sales/invoice-status')
  },

  async getSalesForecast() {
    return apiClient.get('/api/v1/dashboard/sales/forecast')
  },

  async getCompetitionAnalysis() {
    return apiClient.get('/api/v1/dashboard/sales/competition')
  },

  async getTopCustomers(limit: number = 10) {
    const params = new URLSearchParams()
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/dashboard/sales/top-customers?${params}`)
  },
}

// HR Dashboard Service
export const hrDashboardService = {
  async getHROverview() {
    return apiClient.get('/api/v1/dashboard/hr/overview')
  },

  async getPayrollSummary(startDate?: Date, endDate?: Date) {
    return apiClient.post('/api/v1/dashboard/hr/payroll', {
      start_date: startDate,
      end_date: endDate,
    })
  },

  async getAttendanceDashboard(month?: Date) {
    return apiClient.post('/api/v1/dashboard/hr/attendance', {
      month: month || new Date(),
    })
  },

  async getLeaveDashboard() {
    return apiClient.get('/api/v1/dashboard/hr/leaves')
  },

  async getComplianceDashboard() {
    return apiClient.get('/api/v1/dashboard/hr/compliance')
  },

  async getHeadcountByDepartment() {
    return apiClient.get('/api/v1/dashboard/hr/headcount')
  },

  async getPerformanceMetrics() {
    return apiClient.get('/api/v1/dashboard/hr/performance')
  },
}

// Purchase Dashboard Service
export const purchaseDashboardService = {
  async getPurchaseSummary() {
    return apiClient.get('/api/v1/purchase/summary')
  },

  async getVendorList(page: number = 1, limit: number = 50) {
    const params = new URLSearchParams()
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/purchase/vendors?${params}`)
  },

  async getVendorScorecard() {
    return apiClient.get('/api/v1/purchase/vendors/scorecard')
  },

  async getPOStatus() {
    return apiClient.get('/api/v1/purchase/po-status')
  },

  async getCostAnalysis() {
    return apiClient.get('/api/v1/purchase/cost-analysis')
  },
}

// Project Management Dashboard Service
export const projectDashboardService = {
  async getProjectSummary() {
    return apiClient.get('/api/v1/projects/summary')
  },

  async getProjectList(page: number = 1, limit: number = 50) {
    const params = new URLSearchParams()
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/projects/list?${params}`)
  },

  async getProjectPortfolio() {
    return apiClient.get('/api/v1/projects/portfolio')
  },

  async getProjectTimeline(projectId?: string) {
    const url = projectId 
      ? `/api/v1/projects/${projectId}/timeline`
      : '/api/v1/projects/timeline'
    return apiClient.get(url)
  },

  async getProjectStats() {
    return apiClient.get('/api/v1/projects/stats')
  },
}

// Pre-Sales (Opportunities) Dashboard Service
export const prealesDashboardService = {
  async getSalesPipeline() {
    return apiClient.get('/api/v1/sales/pipeline')
  },

  async getOpportunities(stage?: string, page: number = 1, limit: number = 50) {
    const params = new URLSearchParams()
    if (stage) params.append('stage', stage)
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/sales/opportunities?${params}`)
  },

  async getTopDeals(limit: number = 10) {
    const params = new URLSearchParams()
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/sales/top-deals?${params}`)
  },

  async getConversionMetrics() {
    return apiClient.get('/api/v1/sales/conversion-metrics')
  },
}

// Inventory Dashboard Service
export const inventoryDashboardService = {
  async getInventorySummary() {
    return apiClient.get('/api/v1/inventory/summary')
  },

  async getWarehouseDistribution() {
    return apiClient.get('/api/v1/inventory/warehouses')
  },

  async getRealEstateSummary() {
    return apiClient.get('/api/v1/real-estate/summary')
  },

  async getRealEstateProperties(page: number = 1, limit: number = 50) {
    const params = new URLSearchParams()
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/real-estate/properties?${params}`)
  },

  async getInventoryByWarehouse(warehouseId?: string) {
    const url = warehouseId
      ? `/api/v1/inventory/warehouse/${warehouseId}`
      : '/api/v1/inventory/by-warehouse'
    return apiClient.get(url)
  },
}

// Gamification Dashboard Service
export const gamificationDashboardService = {
  async getGamificationOverview() {
    return apiClient.get('/api/v1/dashboard/gamification/overview')
  },

  async getLeaderboard(limit: number = 100, page: number = 1) {
    const params = new URLSearchParams()
    params.append('limit', limit.toString())
    params.append('page', page.toString())
    return apiClient.get(`/api/v1/gamification-advanced/leaderboard?${params}`)
  },

  async getUserChallenges(userId?: string) {
    const url = userId
      ? `/api/v1/gamification-advanced/user/${userId}/challenges`
      : '/api/v1/gamification-advanced/challenges'
    return apiClient.get(url)
  },

  async getRewardsShop(page: number = 1, limit: number = 50) {
    const params = new URLSearchParams()
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/gamification-advanced/rewards?${params}`)
  },

  async getEngagementAnalytics() {
    return apiClient.get('/api/v1/gamification-advanced/stats')
  },
}

// Construction Dashboard Service
export const constructionDashboardService = {
  async getConstructionProjects(page: number = 1, limit: number = 50) {
    const params = new URLSearchParams()
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/construction/projects?${params}`)
  },

  async getBoqSummary() {
    return apiClient.get('/api/v1/boq/summary')
  },

  async getProjectTimeline(projectId: string) {
    return apiClient.get(`/api/v1/construction/projects/${projectId}/timeline`)
  },

  async getSafetyMetrics() {
    return apiClient.get('/api/v1/construction/safety')
  },

  async getWorkerAllocation() {
    return apiClient.get('/api/v1/construction/workers')
  },
}

// General Ledger (Accounting) Service
export const generalLedgerService = {
  async getLedgerEntries(
    accountCode: string,
    startDate?: Date,
    endDate?: Date,
    page: number = 1,
    limit: number = 100
  ) {
    const params = new URLSearchParams()
    params.append('account_code', accountCode)
    if (startDate) params.append('start_date', startDate.toISOString())
    if (endDate) params.append('end_date', endDate.toISOString())
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/gl/entries?${params}`)
  },

  async getVouchers(type?: string, startDate?: Date, endDate?: Date, page: number = 1, limit: number = 50) {
    const params = new URLSearchParams()
    if (type) params.append('type', type)
    if (startDate) params.append('start_date', startDate.toISOString())
    if (endDate) params.append('end_date', endDate.toISOString())
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/gl/vouchers?${params}`)
  },

  async getTrialBalance(asOfDate?: Date) {
    const params = new URLSearchParams()
    if (asOfDate) params.append('as_of_date', asOfDate.toISOString())
    return apiClient.get(`/api/v1/gl/trial-balance${params.size > 0 ? '?' + params : ''}`)
  },

  async getReceiptVouchers(page: number = 1, limit: number = 50) {
    const params = new URLSearchParams()
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/gl/receipt-vouchers?${params}`)
  },

  async getPaymentVouchers(page: number = 1, limit: number = 50) {
    const params = new URLSearchParams()
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/gl/payment-vouchers?${params}`)
  },

  async getJournalVouchers(page: number = 1, limit: number = 50) {
    const params = new URLSearchParams()
    params.append('page', page.toString())
    params.append('limit', limit.toString())
    return apiClient.get(`/api/v1/gl/journal-vouchers?${params}`)
  },
}
