import { apiClient } from './api'
import { Invoice, Payment, Expense, Account, AccountingMetrics, FinancialReport } from '@/types/accounts'

export const accountsService = {
  // Invoices
  async getInvoices(): Promise<Invoice[]> {
    return apiClient.get<Invoice[]>('/api/v1/accounts/invoices')
  },

  async getInvoice(id: string): Promise<Invoice> {
    return apiClient.get<Invoice>(`/api/v1/accounts/invoices/${id}`)
  },

  async createInvoice(data: Partial<Invoice>): Promise<Invoice> {
    return apiClient.post<Invoice>('/api/v1/accounts/invoices', data)
  },

  async updateInvoice(id: string, data: Partial<Invoice>): Promise<Invoice> {
    return apiClient.put<Invoice>(`/api/v1/accounts/invoices/${id}`, data)
  },

  async deleteInvoice(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/accounts/invoices/${id}`)
  },

  async updateInvoiceStatus(id: string, status: string): Promise<Invoice> {
    return apiClient.put<Invoice>(`/api/v1/accounts/invoices/${id}/status`, { status })
  },

  // Payments
  async getPayments(): Promise<Payment[]> {
    return apiClient.get<Payment[]>('/api/v1/accounts/payments')
  },

  async createPayment(data: Partial<Payment>): Promise<Payment> {
    return apiClient.post<Payment>('/api/v1/accounts/payments', data)
  },

  async updatePayment(id: string, data: Partial<Payment>): Promise<Payment> {
    return apiClient.put<Payment>(`/api/v1/accounts/payments/${id}`, data)
  },

  async deletePayment(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/accounts/payments/${id}`)
  },

  // Expenses
  async getExpenses(): Promise<Expense[]> {
    return apiClient.get<Expense[]>('/api/v1/accounts/expenses')
  },

  async createExpense(data: Partial<Expense>): Promise<Expense> {
    return apiClient.post<Expense>('/api/v1/accounts/expenses', data)
  },

  async updateExpense(id: string, data: Partial<Expense>): Promise<Expense> {
    return apiClient.put<Expense>(`/api/v1/accounts/expenses/${id}`, data)
  },

  async deleteExpense(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/accounts/expenses/${id}`)
  },

  async updateExpenseStatus(id: string, status: string): Promise<Expense> {
    return apiClient.put<Expense>(`/api/v1/accounts/expenses/${id}/status`, { status })
  },

  // Accounts
  async getAccounts(): Promise<Account[]> {
    return apiClient.get<Account[]>('/api/v1/accounts/bank-accounts')
  },

  async createAccount(data: Partial<Account>): Promise<Account> {
    return apiClient.post<Account>('/api/v1/accounts/bank-accounts', data)
  },

  async updateAccount(id: string, data: Partial<Account>): Promise<Account> {
    return apiClient.put<Account>(`/api/v1/accounts/bank-accounts/${id}`, data)
  },

  async deleteAccount(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/accounts/bank-accounts/${id}`)
  },

  // Metrics
  async getMetrics(): Promise<AccountingMetrics> {
    return apiClient.get<AccountingMetrics>('/api/v1/accounts/metrics')
  },

  async getFinancialReport(): Promise<FinancialReport[]> {
    return apiClient.get<FinancialReport[]>('/api/v1/accounts/reports')
  },
}
