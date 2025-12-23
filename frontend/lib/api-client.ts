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
    try {
      const response = await this.client.get<any>(url, config)
      // Handle null response body
      if (!response.data) {
        return null
      }
      // Return nested data if available, otherwise return the response data itself
      return response.data?.data ?? response.data
    } catch (error) {
      console.error('API GET error:', error)
      throw error
    }
  }

  async post<T>(url: string, data?: any, config = {}) {
    try {
      const response = await this.client.post<any>(url, data, config)
      if (!response.data) {
        return null
      }
      return response.data?.data ?? response.data
    } catch (error) {
      console.error('API POST error:', error)
      throw error
    }
  }

  async put<T>(url: string, data?: any, config = {}) {
    try {
      const response = await this.client.put<any>(url, data, config)
      if (!response.data) {
        return null
      }
      return response.data?.data ?? response.data
    } catch (error) {
      console.error('API PUT error:', error)
      throw error
    }
  }

  async patch<T>(url: string, data?: any, config = {}) {
    try {
      const response = await this.client.patch<any>(url, data, config)
      if (!response.data) {
        return null
      }
      return response.data?.data ?? response.data
    } catch (error) {
      console.error('API PATCH error:', error)
      throw error
    }
  }

  async delete<T>(url: string, config = {}) {
    try {
      const response = await this.client.delete<any>(url, config)
      if (!response.data) {
        return null
      }
      return response.data?.data ?? response.data
    } catch (error) {
      console.error('API DELETE error:', error)
      throw error
    }
  }
}

export const apiClient = new ApiClient()
