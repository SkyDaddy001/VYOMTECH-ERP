import { apiClient } from './api'

export interface Module {
  id: string
  name: string
  description: string
  category: string
  status: 'active' | 'inactive' | 'beta' | 'deprecated'
  pricing_model: 'free' | 'per_user' | 'per_project' | 'per_company' | 'flat' | 'tiered'
  base_cost: number
  cost_per_user: number
  cost_per_project: number
  cost_per_company: number
  max_users?: number | null
  max_projects?: number | null
  max_companies?: number | null
  is_core: boolean
  trial_days_allowed: number
}

export interface Company {
  id: string
  tenant_id: string
  name: string
  description: string
  status: 'active' | 'inactive'
  industry_type?: string
  employee_count?: number
  website?: string
  max_projects: number
  max_users: number
  current_user_count: number
  current_project_count: number
  billing_email?: string
  billing_address?: string
  created_at: string
  updated_at: string
}

export interface Project {
  id: string
  company_id: string
  name: string
  description: string
  status: 'active' | 'inactive'
  budget_allocated: number
  created_at: string
  updated_at: string
}

export interface PricingPlan {
  id: string
  name: string
  description: string
  price: number
  billing_cycle: 'monthly' | 'quarterly' | 'annual'
  included_modules: string[]
  created_at: string
  updated_at: string
}

export interface Invoice {
  id: string
  tenant_id: string
  invoice_number: string
  total_amount: number
  status: 'draft' | 'sent' | 'paid' | 'overdue'
  paid_at?: string
  created_at: string
  updated_at: string
}

class Phase3CApi {
  private apiClient = apiClient

  // ==================== MODULE ENDPOINTS ====================

  async registerModule(module: Partial<Module>) {
    return this.apiClient.post('/modules/register', module)
  }

  async listModules(status?: string) {
    const params = status ? `?status=${status}` : ''
    return this.apiClient.get(`/modules${params}`)
  }

  async subscribeToModule(subscriptionData: {
    module_id: string
    scope_level: 'tenant' | 'company' | 'project'
    scope_id: string
  }) {
    return this.apiClient.post('/modules/subscribe', subscriptionData)
  }

  async toggleModule(toggleData: {
    subscription_id: string
    enabled: boolean
  }) {
    return this.apiClient.put('/modules/toggle', toggleData)
  }

  async getModuleUsage(params?: {
    subscription_id?: string
    start_date?: string
    end_date?: string
  }) {
    let query = '/modules/usage'
    if (params) {
      const queryParams = new URLSearchParams()
      if (params.subscription_id) queryParams.append('subscription_id', params.subscription_id)
      if (params.start_date) queryParams.append('start_date', params.start_date)
      if (params.end_date) queryParams.append('end_date', params.end_date)
      query += '?' + queryParams.toString()
    }
    return this.apiClient.get(query)
  }

  async listModuleSubscriptions(params?: {
    company_id?: string
    project_id?: string
  }) {
    let query = '/modules/subscriptions'
    if (params) {
      const queryParams = new URLSearchParams()
      if (params.company_id) queryParams.append('company_id', params.company_id)
      if (params.project_id) queryParams.append('project_id', params.project_id)
      query += '?' + queryParams.toString()
    }
    return this.apiClient.get(query)
  }

  // ==================== COMPANY ENDPOINTS ====================

  async createCompany(company: Omit<Company, 'id' | 'tenant_id' | 'current_user_count' | 'current_project_count' | 'created_at' | 'updated_at'>) {
    return this.apiClient.post('/companies', company)
  }

  async listCompanies() {
    return this.apiClient.get('/companies')
  }

  async getCompany(companyId: string) {
    return this.apiClient.get(`/companies/${companyId}`)
  }

  async updateCompany(companyId: string, updates: Partial<Company>) {
    return this.apiClient.put(`/companies/${companyId}`, updates)
  }

  async createProject(companyId: string, project: Omit<Project, 'id' | 'company_id' | 'created_at' | 'updated_at'>) {
    return this.apiClient.post(`/companies/${companyId}/projects`, project)
  }

  async listProjects(companyId: string) {
    return this.apiClient.get(`/companies/${companyId}/projects`)
  }

  async getProject(companyId: string, projectId: string) {
    return this.apiClient.get(`/companies/${companyId}/projects/${projectId}`)
  }

  async getCompanyMembers(companyId: string) {
    return this.apiClient.get(`/companies/${companyId}/members`)
  }

  async addMemberToCompany(companyId: string, memberData: {
    user_id: string
    role: string
    department?: string
  }) {
    return this.apiClient.post(`/companies/${companyId}/members`, memberData)
  }

  async addMemberToProject(companyId: string, projectId: string, memberData: {
    user_id: string
    role: string
  }) {
    return this.apiClient.post(
      `/companies/${companyId}/projects/${projectId}/members`,
      memberData
    )
  }

  async getProjectMembers(companyId: string, projectId: string) {
    return this.apiClient.get(
      `/companies/${companyId}/projects/${projectId}/members`
    )
  }

  async removeProjectMember(companyId: string, projectId: string, userId: string) {
    return this.apiClient.delete(
      `/companies/${companyId}/projects/${projectId}/members/${userId}`
    )
  }

  // ==================== BILLING ENDPOINTS ====================

  async createPricingPlan(plan: Omit<PricingPlan, 'id' | 'created_at' | 'updated_at'>) {
    return this.apiClient.post('/billing/plans', plan)
  }

  async listPricingPlans() {
    return this.apiClient.get('/billing/plans')
  }

  async subscribeToPlan(subscriptionData: {
    plan_id: string
    billing_cycle: string
  }) {
    return this.apiClient.post('/billing/subscribe', subscriptionData)
  }

  async recordUsageMetrics(usageData: {
    metric_type: string
    value: number
    timestamp?: string
  }) {
    return this.apiClient.post('/billing/usage', usageData)
  }

  async getUsageMetrics(params?: {
    start_date?: string
    end_date?: string
    metric_type?: string
  }) {
    let query = '/billing/usage'
    if (params) {
      const queryParams = new URLSearchParams()
      if (params.start_date) queryParams.append('start_date', params.start_date)
      if (params.end_date) queryParams.append('end_date', params.end_date)
      if (params.metric_type) queryParams.append('metric_type', params.metric_type)
      query += '?' + queryParams.toString()
    }
    return this.apiClient.get(query)
  }

  async listInvoices() {
    return this.apiClient.get('/billing/invoices')
  }

  async getInvoice(invoiceId: string) {
    return this.apiClient.get(`/billing/invoices/${invoiceId}`)
  }

  async markInvoiceAsPaid(invoiceId: string, paymentData: {
    payment_method: string
    transaction_id?: string
  }) {
    return this.apiClient.put(`/billing/invoices/${invoiceId}/pay`, paymentData)
  }

  async calculateMonthlyCharges(params?: {
    month?: string
    year?: number
  }) {
    let query = '/billing/charges'
    if (params) {
      const queryParams = new URLSearchParams()
      if (params.month) queryParams.append('month', params.month)
      if (params.year) queryParams.append('year', params.year.toString())
      query += '?' + queryParams.toString()
    }
    return this.apiClient.get(query)
  }
}

export const phase3cApi = new Phase3CApi()
