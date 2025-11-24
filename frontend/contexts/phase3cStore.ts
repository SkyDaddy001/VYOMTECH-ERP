import { create } from 'zustand'
import { phase3cApi, Module, Company, Project, PricingPlan, Invoice } from '@/services/phase3cAPI'

interface ModuleState {
  modules: Module[]
  selectedModule: Module | null
  subscriptions: any[]
  loading: boolean
  error: string | null

  // Actions
  fetchModules: (status?: string) => Promise<void>
  selectModule: (module: Module) => void
  registerModule: (module: Partial<Module>) => Promise<void>
  subscribeToModule: (subscriptionData: any) => Promise<void>
  toggleModule: (toggleData: any) => Promise<void>
  fetchModuleUsage: (params?: any) => Promise<any>
  fetchSubscriptions: (params?: any) => Promise<void>
}

export const useModuleStore = create<ModuleState>((set) => ({
  modules: [],
  selectedModule: null,
  subscriptions: [],
  loading: false,
  error: null,

  fetchModules: async (status?: string) => {
    set({ loading: true, error: null })
    try {
      const data = await phase3cApi.listModules(status)
      const response = data as any
      set({ modules: (response.modules || response) as Module[], loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch modules', loading: false })
    }
  },

  selectModule: (module: Module) => {
    set({ selectedModule: module })
  },

  registerModule: async (module: Partial<Module>) => {
    set({ loading: true, error: null })
    try {
      const newModule = await phase3cApi.registerModule(module)
      set((state) => ({
        modules: [...state.modules, newModule as Module],
        loading: false,
      }))
    } catch (error: any) {
      set({ error: error.message || 'Failed to register module', loading: false })
    }
  },

  subscribeToModule: async (subscriptionData: any) => {
    set({ loading: true, error: null })
    try {
      const subscription = await phase3cApi.subscribeToModule(subscriptionData)
      set((state) => ({
        subscriptions: [...state.subscriptions, subscription] as any[],
        loading: false,
      }))
    } catch (error: any) {
      set({ error: error.message || 'Failed to subscribe to module', loading: false })
    }
  },

  toggleModule: async (toggleData: any) => {
    set({ loading: true, error: null })
    try {
      await phase3cApi.toggleModule(toggleData)
      set({ loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to toggle module', loading: false })
    }
  },

  fetchModuleUsage: async (params?: any) => {
    set({ loading: true, error: null })
    try {
      const usage = await phase3cApi.getModuleUsage(params)
      set({ loading: false })
      return usage
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch usage', loading: false })
      throw error
    }
  },

  fetchSubscriptions: async (params?: any) => {
    set({ loading: true, error: null })
    try {
      const subscriptions = await phase3cApi.listModuleSubscriptions(params)
      const data = subscriptions as any
      set({ subscriptions: (data.subscriptions || data) as any[], loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch subscriptions', loading: false })
    }
  },
}))

interface CompanyState {
  companies: Company[]
  selectedCompany: Company | null
  projects: Project[]
  members: any[]
  loading: boolean
  error: string | null

  // Actions
  fetchCompanies: () => Promise<void>
  selectCompany: (company: Company) => void
  createCompany: (company: any) => Promise<void>
  updateCompany: (companyId: string, updates: any) => Promise<void>
  fetchProjects: (companyId: string) => Promise<void>
  createProject: (companyId: string, project: any) => Promise<void>
  fetchMembers: (companyId: string) => Promise<void>
  addMember: (companyId: string, memberData: any) => Promise<void>
  removeMember: (companyId: string, projectId: string, userId: string) => Promise<void>
}

export const useCompanyStore = create<CompanyState>((set) => ({
  companies: [],
  selectedCompany: null,
  projects: [],
  members: [],
  loading: false,
  error: null,

  fetchCompanies: async () => {
    set({ loading: true, error: null })
    try {
      const data = await phase3cApi.listCompanies()
      const response = data as any
      set({ companies: (response.companies || response) as Company[], loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch companies', loading: false })
    }
  },

  selectCompany: (company: Company) => {
    set({ selectedCompany: company })
  },

  createCompany: async (company: any) => {
    set({ loading: true, error: null })
    try {
      const newCompany = await phase3cApi.createCompany(company)
      set((state) => ({
        companies: [...state.companies, newCompany as Company],
        loading: false,
      }))
    } catch (error: any) {
      set({ error: error.message || 'Failed to create company', loading: false })
    }
  },

  updateCompany: async (companyId: string, updates: any) => {
    set({ loading: true, error: null })
    try {
      const updated = await phase3cApi.updateCompany(companyId, updates)
      set((state) => ({
        companies: state.companies.map((c) => (c.id === companyId ? (updated as Company) : c)),
        selectedCompany: state.selectedCompany?.id === companyId ? (updated as Company) : state.selectedCompany,
        loading: false,
      }))
    } catch (error: any) {
      set({ error: error.message || 'Failed to update company', loading: false })
    }
  },

  fetchProjects: async (companyId: string) => {
    set({ loading: true, error: null })
    try {
      const data = await phase3cApi.listProjects(companyId)
      const response = data as any
      set({ projects: (response.projects || response) as Project[], loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch projects', loading: false })
    }
  },

  createProject: async (companyId: string, project: any) => {
    set({ loading: true, error: null })
    try {
      const newProject = await phase3cApi.createProject(companyId, project)
      set((state) => ({
        projects: [...state.projects, newProject as Project],
        loading: false,
      }))
    } catch (error: any) {
      set({ error: error.message || 'Failed to create project', loading: false })
    }
  },

  fetchMembers: async (companyId: string) => {
    set({ loading: true, error: null })
    try {
      const data = await phase3cApi.getCompanyMembers(companyId)
      const response = data as any
      set({ members: (response.members || response) as any[], loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch members', loading: false })
    }
  },

  addMember: async (companyId: string, memberData: any) => {
    set({ loading: true, error: null })
    try {
      const member = await phase3cApi.addMemberToCompany(companyId, memberData)
      set((state) => ({
        members: [...state.members, member],
        loading: false,
      }))
    } catch (error: any) {
      set({ error: error.message || 'Failed to add member', loading: false })
    }
  },

  removeMember: async (companyId: string, projectId: string, userId: string) => {
    set({ loading: true, error: null })
    try {
      await phase3cApi.removeProjectMember(companyId, projectId, userId)
      set((state) => ({
        members: state.members.filter((m) => m.user_id !== userId),
        loading: false,
      }))
    } catch (error: any) {
      set({ error: error.message || 'Failed to remove member', loading: false })
    }
  },
}))

interface BillingState {
  plans: PricingPlan[]
  invoices: Invoice[]
  charges: any
  usageMetrics: any[]
  loading: boolean
  error: string | null

  // Actions
  fetchPlans: () => Promise<void>
  createPlan: (plan: any) => Promise<void>
  subscribeToPlan: (subscriptionData: any) => Promise<void>
  fetchInvoices: () => Promise<void>
  markAsPaid: (invoiceId: string, paymentData: any) => Promise<void>
  recordUsage: (usageData: any) => Promise<void>
  fetchUsageMetrics: (params?: any) => Promise<void>
  fetchCharges: (params?: any) => Promise<void>
}

export const useBillingStore = create<BillingState>((set) => ({
  plans: [],
  invoices: [],
  charges: null,
  usageMetrics: [],
  loading: false,
  error: null,

  fetchPlans: async () => {
    set({ loading: true, error: null })
    try {
      const data = await phase3cApi.listPricingPlans()
      const response = data as any
      set({ plans: (response.plans || response) as PricingPlan[], loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch plans', loading: false })
    }
  },

  createPlan: async (plan: any) => {
    set({ loading: true, error: null })
    try {
      const newPlan = await phase3cApi.createPricingPlan(plan)
      set((state) => ({
        plans: [...state.plans, newPlan as PricingPlan],
        loading: false,
      }))
    } catch (error: any) {
      set({ error: error.message || 'Failed to create plan', loading: false })
    }
  },

  subscribeToPlan: async (subscriptionData: any): Promise<void> => {
    set({ loading: true, error: null })
    try {
      await phase3cApi.subscribeToPlan(subscriptionData)
      set({ loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to subscribe', loading: false })
    }
  },

  fetchInvoices: async () => {
    set({ loading: true, error: null })
    try {
      const data = await phase3cApi.listInvoices()
      const response = data as any
      set({ invoices: (response.invoices || response) as Invoice[], loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch invoices', loading: false })
    }
  },

  markAsPaid: async (invoiceId: string, paymentData: any) => {
    set({ loading: true, error: null })
    try {
      const updated = await phase3cApi.markInvoiceAsPaid(invoiceId, paymentData)
      set((state) => ({
        invoices: state.invoices.map((inv) => (inv.id === invoiceId ? (updated as Invoice) : inv)),
        loading: false,
      }))
    } catch (error: any) {
      set({ error: error.message || 'Failed to mark as paid', loading: false })
    }
  },

  recordUsage: async (usageData: any) => {
    set({ loading: true, error: null })
    try {
      await phase3cApi.recordUsageMetrics(usageData)
      set({ loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to record usage', loading: false })
    }
  },

  fetchUsageMetrics: async (params?: any) => {
    set({ loading: true, error: null })
    try {
      const data = await phase3cApi.getUsageMetrics(params)
      const response = data as any
      set({ usageMetrics: (response.metrics || response) as any[], loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch usage metrics', loading: false })
    }
  },

  fetchCharges: async (params?: any) => {
    set({ loading: true, error: null })
    try {
      const data = await phase3cApi.calculateMonthlyCharges(params)
      set({ charges: data, loading: false })
    } catch (error: any) {
      set({ error: error.message || 'Failed to fetch charges', loading: false })
    }
  },
}))
