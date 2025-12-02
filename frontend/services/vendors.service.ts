import { apiClient } from './api'
import { Vendor, PurchaseOrder, VendorMetrics, VendorPayment } from '@/types/vendors'

export const vendorsService = {
  // Vendors
  async getVendors(): Promise<Vendor[]> {
    return apiClient.get<Vendor[]>('/api/v1/vendors')
  },

  async getVendor(id: string): Promise<Vendor> {
    return apiClient.get<Vendor>(`/api/v1/vendors/${id}`)
  },

  async createVendor(data: Partial<Vendor>): Promise<Vendor> {
    return apiClient.post<Vendor>('/api/v1/vendors', data)
  },

  async updateVendor(id: string, data: Partial<Vendor>): Promise<Vendor> {
    return apiClient.put<Vendor>(`/api/v1/vendors/${id}`, data)
  },

  async deleteVendor(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/vendors/${id}`)
  },

  // Purchase Orders
  async getPurchaseOrders(): Promise<PurchaseOrder[]> {
    return apiClient.get<PurchaseOrder[]>('/api/v1/purchase-orders')
  },

  async getPurchaseOrder(id: string): Promise<PurchaseOrder> {
    return apiClient.get<PurchaseOrder>(`/api/v1/purchase-orders/${id}`)
  },

  async createPurchaseOrder(data: Partial<PurchaseOrder>): Promise<PurchaseOrder> {
    return apiClient.post<PurchaseOrder>('/api/v1/purchase-orders', data)
  },

  async updatePurchaseOrder(id: string, data: Partial<PurchaseOrder>): Promise<PurchaseOrder> {
    return apiClient.put<PurchaseOrder>(`/api/v1/purchase-orders/${id}`, data)
  },

  async deletePurchaseOrder(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/purchase-orders/${id}`)
  },

  async updatePOStatus(id: string, status: string): Promise<PurchaseOrder> {
    return apiClient.put<PurchaseOrder>(`/api/v1/purchase-orders/${id}/status`, { status })
  },

  // Payments
  async getVendorPayments(vendorId: string): Promise<VendorPayment[]> {
    return apiClient.get<VendorPayment[]>(`/api/v1/vendors/${vendorId}/payments`)
  },

  async recordPayment(data: Partial<VendorPayment>): Promise<VendorPayment> {
    return apiClient.post<VendorPayment>('/api/v1/vendor-payments', data)
  },

  // Metrics
  async getMetrics(): Promise<VendorMetrics> {
    return apiClient.get<VendorMetrics>('/api/v1/vendors/metrics')
  },
}
