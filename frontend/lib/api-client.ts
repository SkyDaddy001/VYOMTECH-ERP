import axios, { AxiosInstance, AxiosError } from 'axios'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

interface ApiResponse<T> {
  data: T
  message?: string
  success?: boolean
}

class ApiClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    // Add request interceptor for auth token and tenant
    this.client.interceptors.request.use((config) => {
      const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null
      const tenantId = typeof window !== 'undefined' ? localStorage.getItem('tenant_id') : null

      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      if (tenantId) {
        config.headers['X-Tenant-ID'] = tenantId
      }
      return config
    })

    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          // Handle unauthorized
          if (typeof window !== 'undefined') {
            localStorage.removeItem('auth_token')
            localStorage.removeItem('tenant_id')
            localStorage.removeItem('user_id')
            window.location.href = '/login'
          }
        }
        return Promise.reject(error)
      }
    )
  }

  async get<T>(url: string, config = {}) {
    const response = await this.client.get<ApiResponse<T>>(url, config)
    return response.data.data || response.data
  }

  async post<T>(url: string, data?: any, config = {}) {
    const response = await this.client.post<ApiResponse<T>>(url, data, config)
    return response.data.data || response.data
  }

  async put<T>(url: string, data?: any, config = {}) {
    const response = await this.client.put<ApiResponse<T>>(url, data, config)
    return response.data.data || response.data
  }

  async patch<T>(url: string, data?: any, config = {}) {
    const response = await this.client.patch<ApiResponse<T>>(url, data, config)
    return response.data.data || response.data
  }

  async delete<T>(url: string, config = {}) {
    const response = await this.client.delete<ApiResponse<T>>(url, config)
    return response.data.data || response.data
  }
}

export const apiClient = new ApiClient()
