import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

interface ApiResponse<T> {
  data?: T;
  message?: string;
  code?: string;
  total?: number;
  page?: number;
  page_size?: number;
  total_pages?: number;
}

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add token to requests
    this.client.interceptors.request.use((config) => {
      const token = typeof window !== 'undefined' ? localStorage.getItem('access_token') : null;
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Handle errors
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Token expired, redirect to login
          if (typeof window !== 'undefined') {
            localStorage.removeItem('access_token');
            window.location.href = '/login';
          }
        }
        return Promise.reject(error);
      }
    );
  }

  // Auth endpoints
  async login(email: string, password: string) {
    const response = await this.client.post<{ access_token: string; token_type: string; expires_in: number; user: any }>(
      '/auth/login',
      { email, password }
    );
    return response.data;
  }

  async logout() {
    return this.client.post<ApiResponse<void>>('/auth/logout');
  }

  async getCurrentUser() {
    const response = await this.client.get<ApiResponse<any>>('/auth/me');
    return response.data;
  }

  // Lead endpoints
  async createLead(data: any) {
    const response = await this.client.post<ApiResponse<any>>('/leads', data);
    return response.data;
  }

  async getLead(id: string) {
    const response = await this.client.get<ApiResponse<any>>(`/leads/${id}`);
    return response.data;
  }

  async listLeads(page = 1, pageSize = 10, filters: any = {}) {
    const params = {
      page,
      page_size: pageSize,
      ...filters,
    };
    const response = await this.client.get<ApiResponse<any[]>>('/leads', { params });
    return response.data;
  }

  async updateLead(id: string, data: any) {
    const response = await this.client.patch<ApiResponse<any>>(`/leads/${id}`, data);
    return response.data;
  }

  async deleteLead(id: string) {
    const response = await this.client.delete<ApiResponse<void>>(`/leads/${id}`);
    return response.data;
  }

  async updateLeadStatus(id: string, status: string, reason?: string) {
    const response = await this.client.patch<ApiResponse<any>>(`/leads/${id}/status`, {
      status,
      reason,
    });
    return response.data;
  }

  async getLeadStatusHistory(id: string) {
    const response = await this.client.get<ApiResponse<any[]>>(`/leads/${id}/status-history`);
    return response.data;
  }

  // Sprint endpoints
  async createSprint(data: any) {
    const response = await this.client.post<any>('/sprints', data);
    return response.data;
  }

  async getSprint(id: string) {
    const response = await this.client.get<any>(`/sprints/${id}`);
    return response.data;
  }

  async listSprints() {
    const response = await this.client.get<any[]>('/sprints');
    return response.data;
  }

  async updateSprint(id: string, data: any) {
    const response = await this.client.put<any>(`/sprints/${id}`, data);
    return response.data;
  }

  async deleteSprint(id: string) {
    const response = await this.client.delete<void>(`/sprints/${id}`);
    return response.data;
  }

  // Marketing endpoints
  async createMarketingCampaign(data: any) {
    const response = await this.client.post<any>('/marketing/campaigns', data);
    return response.data;
  }

  async listMarketingCampaigns() {
    const response = await this.client.get<any[]>('/marketing/campaigns');
    return response.data;
  }

  async getMarketingCampaign(id: string) {
    const response = await this.client.get<any>(`/marketing/campaigns/${id}`);
    return response.data;
  }

  async updateMarketingCampaign(id: string, data: any) {
    const response = await this.client.put<any>(`/marketing/campaigns/${id}`, data);
    return response.data;
  }

  async deleteMarketingCampaign(id: string) {
    const response = await this.client.delete<void>(`/marketing/campaigns/${id}`);
    return response.data;
  }

  // Site visits endpoints
  async createSiteVisit(data: any) {
    const response = await this.client.post<any>('/site-visits', data);
    return response.data;
  }

  async listSiteVisits() {
    const response = await this.client.get<any[]>('/site-visits');
    return response.data;
  }

  async getSiteVisit(id: string) {
    const response = await this.client.get<any>(`/site-visits/${id}`);
    return response.data;
  }

  async updateSiteVisit(id: string, data: any) {
    const response = await this.client.put<any>(`/site-visits/${id}`, data);
    return response.data;
  }

  async deleteSiteVisit(id: string) {
    const response = await this.client.delete<void>(`/site-visits/${id}`);
    return response.data;
  }

  // Sales endpoints
  async createSale(data: any) {
    const response = await this.client.post<any>('/sales', data);
    return response.data;
  }

  async listSales() {
    const response = await this.client.get<any[]>('/sales');
    return response.data;
  }

  async getSale(id: string) {
    const response = await this.client.get<any>(`/sales/${id}`);
    return response.data;
  }

  async updateSale(id: string, data: any) {
    const response = await this.client.put<any>(`/sales/${id}`, data);
    return response.data;
  }

  async deleteSale(id: string) {
    const response = await this.client.delete<void>(`/sales/${id}`);
    return response.data;
  }
}

export const apiClient = new ApiClient();
