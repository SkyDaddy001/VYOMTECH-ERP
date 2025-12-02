import { apiClient } from './api'
import { Ledger, LedgerSummary, LedgerMetrics, LedgerReport } from '@/types/ledgers'

export const ledgersService = {
  // Ledger entries
  async getLedgerEntries(customerId?: string): Promise<Ledger[]> {
    if (customerId) {
      return apiClient.get<Ledger[]>(`/api/v1/ledgers?customer_id=${customerId}`)
    }
    return apiClient.get<Ledger[]>('/api/v1/ledgers')
  },

  async getLedgerEntry(id: string): Promise<Ledger> {
    return apiClient.get<Ledger>(`/api/v1/ledgers/${id}`)
  },

  async createEntry(data: Partial<Ledger>): Promise<Ledger> {
    return apiClient.post<Ledger>('/api/v1/ledgers', data)
  },

  async updateEntry(id: string, data: Partial<Ledger>): Promise<Ledger> {
    return apiClient.put<Ledger>(`/api/v1/ledgers/${id}`, data)
  },

  async deleteEntry(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/ledgers/${id}`)
  },

  // Customer summary
  async getCustomerSummary(customerId: string): Promise<LedgerSummary> {
    return apiClient.get<LedgerSummary>(`/api/v1/ledgers/customer/${customerId}/summary`)
  },

  async getAllCustomerSummaries(): Promise<LedgerSummary[]> {
    return apiClient.get<LedgerSummary[]>('/api/v1/ledgers/customers/summaries')
  },

  // Metrics
  async getMetrics(): Promise<LedgerMetrics> {
    return apiClient.get<LedgerMetrics>('/api/v1/ledgers/metrics')
  },

  // Reports
  async getReport(period: string): Promise<LedgerReport> {
    return apiClient.get<LedgerReport>(`/api/v1/ledgers/reports/${period}`)
  },

  async getOutstandingCustomers(): Promise<LedgerSummary[]> {
    return apiClient.get<LedgerSummary[]>('/api/v1/ledgers/outstanding')
  },
}
