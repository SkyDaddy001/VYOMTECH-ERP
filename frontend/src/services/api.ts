// frontend/src/services/api.ts
import axios, { AxiosInstance, AxiosError } from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

interface ApiResponse<T> {
  data: T;
  status: number;
  message?: string;
}

interface LoginPayload {
  email: string;
  password: string;
}

interface RegisterPayload {
  email: string;
  password: string;
  name: string;
  tenantId: string;
}

interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  user: any;
  expiresIn: number;
}

class ApiService {
  private client: AxiosInstance;
  private accessToken: string | null = null;
  private refreshToken: string | null = null;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Load tokens from localStorage if available
    if (typeof window !== 'undefined') {
      this.accessToken = localStorage.getItem('accessToken');
      this.refreshToken = localStorage.getItem('refreshToken');
      this.updateAuthHeader();
    }

    // Add request interceptor to attach auth token
    this.client.interceptors.request.use(
      (config) => {
        if (this.accessToken) {
          config.headers.Authorization = `Bearer ${this.accessToken}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Add response interceptor to handle token refresh
    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        if (error.response?.status === 401 && this.refreshToken) {
          try {
            const response = await this.refreshAccessToken();
            if (response.accessToken) {
              this.accessToken = response.accessToken;
              this.updateAuthHeader();
              return this.client.request(error.config!);
            }
          } catch (refreshError) {
            this.logout();
          }
        }
        return Promise.reject(error);
      }
    );
  }

  private updateAuthHeader() {
    if (this.accessToken) {
      this.client.defaults.headers.common['Authorization'] = `Bearer ${this.accessToken}`;
    } else {
      delete this.client.defaults.headers.common['Authorization'];
    }
  }

  // Auth endpoints
  async login(payload: LoginPayload): Promise<LoginResponse> {
    const response = await this.client.post<LoginResponse>('/auth/login', payload);
    this.setTokens(response.data.accessToken, response.data.refreshToken);
    return response.data;
  }

  async register(payload: RegisterPayload): Promise<any> {
    const response = await this.client.post('/auth/register', payload);
    return response.data;
  }

  async refreshAccessToken(): Promise<{ accessToken: string; expiresIn: number }> {
    const response = await this.client.post('/auth/refresh', {
      refreshToken: this.refreshToken,
    });
    return response.data;
  }

  setTokens(accessToken: string, refreshToken: string) {
    this.accessToken = accessToken;
    this.refreshToken = refreshToken;
    if (typeof window !== 'undefined') {
      localStorage.setItem('accessToken', accessToken);
      localStorage.setItem('refreshToken', refreshToken);
    }
    this.updateAuthHeader();
  }

  logout() {
    this.accessToken = null;
    this.refreshToken = null;
    if (typeof window !== 'undefined') {
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
    }
    this.updateAuthHeader();
  }

  // Tenant endpoints
  async getTenant(id: string) {
    const response = await this.client.get(`/tenants/${id}`);
    return response.data;
  }

  async updateTenant(id: string, data: any) {
    const response = await this.client.put(`/tenants/${id}`, data);
    return response.data;
  }

  async listTenants(limit: number = 10, offset: number = 0) {
    const response = await this.client.get('/tenants', {
      params: { limit, offset },
    });
    return response.data;
  }

  // User endpoints
  async getUser(id: string) {
    const response = await this.client.get(`/users/${id}`);
    return response.data;
  }

  async updateUser(id: string, data: any) {
    const response = await this.client.put(`/users/${id}`, data);
    return response.data;
  }

  async listUsers(tenantId: string, limit: number = 10, offset: number = 0) {
    const response = await this.client.get('/users', {
      params: { tenantId, limit, offset },
    });
    return response.data;
  }

  async deleteUser(id: string) {
    await this.client.delete(`/users/${id}`);
  }

  // Call Center endpoints
  async createCallCenter(data: any) {
    const response = await this.client.post('/call-centers', data);
    return response.data;
  }

  async getCallCenter(id: string) {
    const response = await this.client.get(`/call-centers/${id}`);
    return response.data;
  }

  async updateCallCenter(id: string, data: any) {
    const response = await this.client.put(`/call-centers/${id}`, data);
    return response.data;
  }

  async listCallCenters(tenantId: string) {
    const response = await this.client.get('/call-centers', {
      params: { tenantId },
    });
    return response.data;
  }

  // Agent endpoints
  async createAgent(data: any) {
    const response = await this.client.post('/agents', data);
    return response.data;
  }

  async getAgent(id: string) {
    const response = await this.client.get(`/agents/${id}`);
    return response.data;
  }

  async updateAgent(id: string, data: any) {
    const response = await this.client.put(`/agents/${id}`, data);
    return response.data;
  }

  async listAgents(callCenterId: string) {
    const response = await this.client.get('/agents', {
      params: { callCenterId },
    });
    return response.data;
  }

  // Call endpoints
  async createCall(data: any) {
    const response = await this.client.post('/calls', data);
    return response.data;
  }

  async getCall(id: string) {
    const response = await this.client.get(`/calls/${id}`);
    return response.data;
  }

  async updateCall(id: string, data: any) {
    const response = await this.client.put(`/calls/${id}`, data);
    return response.data;
  }

  async listCalls(callCenterId: string, limit: number = 10, offset: number = 0) {
    const response = await this.client.get('/calls', {
      params: { callCenterId, limit, offset },
    });
    return response.data;
  }

  // Campaign endpoints
  async createCampaign(data: any) {
    const response = await this.client.post('/campaigns', data);
    return response.data;
  }

  async getCampaign(id: string) {
    const response = await this.client.get(`/campaigns/${id}`);
    return response.data;
  }

  async updateCampaign(id: string, data: any) {
    const response = await this.client.put(`/campaigns/${id}`, data);
    return response.data;
  }

  async listCampaigns(callCenterId: string) {
    const response = await this.client.get('/campaigns', {
      params: { callCenterId },
    });
    return response.data;
  }

  // Sales Lead endpoints
  async createSalesLead(data: any) {
    const response = await this.client.post('/sales-leads', data);
    return response.data;
  }

  async getSalesLead(id: string) {
    const response = await this.client.get(`/sales-leads/${id}`);
    return response.data;
  }

  async updateSalesLead(id: string, data: any) {
    const response = await this.client.put(`/sales-leads/${id}`, data);
    return response.data;
  }

  async listSalesLeads(tenantId: string, limit: number = 10, offset: number = 0) {
    const response = await this.client.get('/sales-leads', {
      params: { tenantId, limit, offset },
    });
    return response.data;
  }

  // Sales Customer endpoints
  async createSalesCustomer(data: any) {
    const response = await this.client.post('/sales-customers', data);
    return response.data;
  }

  async getSalesCustomer(id: string) {
    const response = await this.client.get(`/sales-customers/${id}`);
    return response.data;
  }

  async updateSalesCustomer(id: string, data: any) {
    const response = await this.client.put(`/sales-customers/${id}`, data);
    return response.data;
  }

  async listSalesCustomers(tenantId: string) {
    const response = await this.client.get('/sales-customers', {
      params: { tenantId },
    });
    return response.data;
  }

  // Inventory Item endpoints
  async createInventoryItem(data: any) {
    const response = await this.client.post('/inventory-items', data);
    return response.data;
  }

  async getInventoryItem(id: string) {
    const response = await this.client.get(`/inventory-items/${id}`);
    return response.data;
  }

  async updateInventoryItem(id: string, data: any) {
    const response = await this.client.put(`/inventory-items/${id}`, data);
    return response.data;
  }

  async listInventoryItems(tenantId: string) {
    const response = await this.client.get('/inventory-items', {
      params: { tenantId },
    });
    return response.data;
  }

  // Stock Level endpoints
  async getStockLevel(id: string) {
    const response = await this.client.get(`/stock-levels/${id}`);
    return response.data;
  }

  async updateStockLevel(id: string, data: any) {
    const response = await this.client.put(`/stock-levels/${id}`, data);
    return response.data;
  }

  async listStockLevels(tenantId: string) {
    const response = await this.client.get('/stock-levels', {
      params: { tenantId },
    });
    return response.data;
  }
}

export const apiService = new ApiService();
export type { LoginPayload, RegisterPayload, LoginResponse };
