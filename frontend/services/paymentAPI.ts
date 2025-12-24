// API service for payment operations

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export interface PaymentInitRequest {
  amount: number
  currency: string
  provider: string
  payment_method: string
  description: string
  customer_name: string
  customer_email: string
  customer_phone: string
  billing_address?: {
    street: string
    city: string
    state: string
    postal_code: string
    country: string
  }
}

export interface PaymentResponse {
  payment_id: string
  order_id: string
  amount: number
  currency: string
  status: string
  provider: string
  payment_method: string
  gateway_order_id: string
  payment_url: string
  expires_at: string
  created_at: string
}

export interface Payment {
  id: string
  tenant_id: string
  order_id: string
  amount: number
  currency: string
  status: string
  provider: string
  payment_method: string
  gateway_order_id: string
  gateway_payment_id: string
  transaction_id: string
  customer_name: string
  customer_email: string
  customer_phone: string
  billing_address?: any
  error_message?: string
  receipt_url?: string
  created_at: string
  updated_at: string
  processed_at?: string
}

export interface PaymentSummary {
  total_payments: number
  total_successful: number
  total_failed: number
  total_pending: number
  successful_count: number
  failed_count: number
  pending_count: number
  success_rate: number
}

export interface RefundRequest {
  amount: string
  reason: string
  notes?: string
}

export interface Refund {
  id: string
  payment_id: string
  amount: number
  status: string
  gateway_refund_id: string
  reason: string
  created_at: string
  processed_at?: string
}

// Get auth token
const getToken = (): string => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('authToken') || ''
  }
  return ''
}

// API request helper
const apiRequest = async (
  endpoint: string,
  options: RequestInit = {}
): Promise<Response> => {
  const token = getToken()
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
    ...options.headers,
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  })

  return response
}

// Payment API methods
export const paymentAPI = {
  // Initiate payment
  initiatePayment: async (data: PaymentInitRequest): Promise<PaymentResponse> => {
    const response = await apiRequest('/api/v1/payments/initiate', {
      method: 'POST',
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to initiate payment')
    }

    return response.json()
  },

  // Get payment status
  getPaymentStatus: async (paymentId: string): Promise<Payment> => {
    const response = await apiRequest(`/api/v1/payments/${paymentId}`)

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to fetch payment')
    }

    return response.json()
  },

  // Verify payment
  verifyPayment: async (
    paymentId: string,
    paymentGatewayId: string,
    signature: string
  ): Promise<Payment> => {
    const response = await apiRequest(`/api/v1/payments/${paymentId}/verify`, {
      method: 'POST',
      body: JSON.stringify({
        payment_gateway_id: paymentGatewayId,
        signature,
      }),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to verify payment')
    }

    return response.json()
  },

  // Create refund
  createRefund: async (paymentId: string, data: RefundRequest): Promise<Refund> => {
    const response = await apiRequest(`/api/v1/payments/${paymentId}/refund`, {
      method: 'POST',
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to create refund')
    }

    return response.json()
  },

  // List payments
  listPayments: async (
    status?: string,
    provider?: string,
    limit: number = 50,
    offset: number = 0
  ): Promise<{ payments: Payment[]; total: number }> => {
    const params = new URLSearchParams()
    if (status) params.append('status', status)
    if (provider) params.append('provider', provider)
    params.append('limit', limit.toString())
    params.append('offset', offset.toString())

    const response = await apiRequest(`/api/v1/payments/list?${params.toString()}`)

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to fetch payments')
    }

    return response.json()
  },

  // Get payment summary
  getPaymentSummary: async (): Promise<PaymentSummary> => {
    const response = await apiRequest('/api/v1/payments/summary')

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to fetch summary')
    }

    return response.json()
  },

  // Get available payment methods
  getPaymentMethods: async (): Promise<any[]> => {
    const response = await apiRequest('/api/v1/payments/methods/available')

    if (!response.ok) {
      throw new Error('Failed to fetch payment methods')
    }

    const data = await response.json()
    return data.payment_methods || []
  },

  // Get available banks
  getBanks: async (): Promise<any[]> => {
    const response = await apiRequest('/api/v1/payments/methods/banks')

    if (!response.ok) {
      throw new Error('Failed to fetch banks')
    }

    const data = await response.json()
    return data.banks || []
  },

  // Get payment gateways
  getPaymentGateways: async (): Promise<any[]> => {
    const response = await apiRequest('/api/v1/payment-config/gateways')

    if (!response.ok) {
      throw new Error('Failed to fetch gateways')
    }

    const data = await response.json()
    return data.gateways || []
  },

  // Configure payment gateway
  configureGateway: async (data: {
    provider: string
    api_key: string
    api_secret: string
    settings?: Record<string, any>
  }): Promise<any> => {
    const response = await apiRequest('/api/v1/payment-config/gateways/configure', {
      method: 'POST',
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to configure gateway')
    }

    return response.json()
  },

  // Update payment gateway
  updateGateway: async (
    gatewayId: string,
    data: {
      api_key?: string
      api_secret?: string
      is_active?: boolean
      settings?: Record<string, any>
    }
  ): Promise<any> => {
    const response = await apiRequest(
      `/api/v1/payment-config/gateways/${gatewayId}`,
      {
        method: 'PUT',
        body: JSON.stringify(data),
      }
    )

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to update gateway')
    }

    return response.json()
  },

  // Delete payment gateway
  deleteGateway: async (gatewayId: string): Promise<void> => {
    const response = await apiRequest(
      `/api/v1/payment-config/gateways/${gatewayId}`,
      {
        method: 'DELETE',
      }
    )

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to delete gateway')
    }
  },
}

export default paymentAPI
