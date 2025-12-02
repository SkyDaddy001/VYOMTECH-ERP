import { apiClient } from './api'
import { SalesOrder, SalesTarget, SalesQuota, SalesMetrics, SalesForecast } from '@/types/sales'

export const salesService = {
  // Sales Orders
  async getOrders(): Promise<SalesOrder[]> {
    return apiClient.get<SalesOrder[]>('/api/v1/sales/orders')
  },

  async getOrder(id: string): Promise<SalesOrder> {
    return apiClient.get<SalesOrder>(`/api/v1/sales/orders/${id}`)
  },

  async createOrder(data: Partial<SalesOrder>): Promise<SalesOrder> {
    return apiClient.post<SalesOrder>('/api/v1/sales/orders', data)
  },

  async updateOrder(id: string, data: Partial<SalesOrder>): Promise<SalesOrder> {
    return apiClient.put<SalesOrder>(`/api/v1/sales/orders/${id}`, data)
  },

  async deleteOrder(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/sales/orders/${id}`)
  },

  async updateOrderStatus(id: string, status: string): Promise<SalesOrder> {
    return apiClient.put<SalesOrder>(`/api/v1/sales/orders/${id}/status`, { status })
  },

  // Sales Targets
  async getTargets(): Promise<SalesTarget[]> {
    return apiClient.get<SalesTarget[]>('/api/v1/sales/targets')
  },

  async createTarget(data: Partial<SalesTarget>): Promise<SalesTarget> {
    return apiClient.post<SalesTarget>('/api/v1/sales/targets', data)
  },

  async updateTarget(id: string, data: Partial<SalesTarget>): Promise<SalesTarget> {
    return apiClient.put<SalesTarget>(`/api/v1/sales/targets/${id}`, data)
  },

  async deleteTarget(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/sales/targets/${id}`)
  },

  // Sales Quotas
  async getQuotas(): Promise<SalesQuota[]> {
    return apiClient.get<SalesQuota[]>('/api/v1/sales/quotas')
  },

  async createQuota(data: Partial<SalesQuota>): Promise<SalesQuota> {
    return apiClient.post<SalesQuota>('/api/v1/sales/quotas', data)
  },

  async updateQuota(id: string, data: Partial<SalesQuota>): Promise<SalesQuota> {
    return apiClient.put<SalesQuota>(`/api/v1/sales/quotas/${id}`, data)
  },

  async deleteQuota(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/sales/quotas/${id}`)
  },

  // Metrics
  async getMetrics(): Promise<SalesMetrics> {
    return apiClient.get<SalesMetrics>('/api/v1/sales/metrics')
  },

  async getForecast(): Promise<SalesForecast[]> {
    return apiClient.get<SalesForecast[]>('/api/v1/sales/forecast')
  },
}
