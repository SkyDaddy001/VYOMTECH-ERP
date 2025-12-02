import { apiClient } from './api'
import { Company, Project, CompanyMember, ProjectMember } from '@/types/company'

export const companyService = {
  // Companies
  async getCompanies(): Promise<Company[]> {
    return apiClient.get<Company[]>('/api/v1/companies')
  },

  async getCompany(id: string): Promise<Company> {
    return apiClient.get<Company>(`/api/v1/companies/${id}`)
  },

  async createCompany(data: Partial<Company>): Promise<Company> {
    return apiClient.post<Company>('/api/v1/companies', data)
  },

  async updateCompany(id: string, data: Partial<Company>): Promise<Company> {
    return apiClient.put<Company>(`/api/v1/companies/${id}`, data)
  },

  async deleteCompany(id: string): Promise<void> {
    return apiClient.post(`/api/v1/companies/${id}/delete`, {})
  },

  // Projects
  async getProjects(companyId?: string): Promise<Project[]> {
    const url = companyId ? `/api/v1/companies/${companyId}/projects` : '/api/v1/projects'
    return apiClient.get<Project[]>(url)
  },

  async getProject(id: string): Promise<Project> {
    return apiClient.get<Project>(`/api/v1/projects/${id}`)
  },

  async createProject(companyId: string, data: Partial<Project>): Promise<Project> {
    return apiClient.post<Project>(`/api/v1/companies/${companyId}/projects`, data)
  },

  async updateProject(id: string, data: Partial<Project>): Promise<Project> {
    return apiClient.put<Project>(`/api/v1/projects/${id}`, data)
  },

  async deleteProject(id: string): Promise<void> {
    return apiClient.post(`/api/v1/projects/${id}/delete`, {})
  },

  // Company Members
  async getCompanyMembers(companyId: string): Promise<CompanyMember[]> {
    return apiClient.get<CompanyMember[]>(`/api/v1/companies/${companyId}/members`)
  },

  async addCompanyMember(companyId: string, data: Partial<CompanyMember>): Promise<CompanyMember> {
    return apiClient.post<CompanyMember>(`/api/v1/companies/${companyId}/members`, data)
  },

  async updateCompanyMember(companyId: string, userId: number, data: Partial<CompanyMember>): Promise<CompanyMember> {
    return apiClient.put<CompanyMember>(`/api/v1/companies/${companyId}/members/${userId}`, data)
  },

  async removeCompanyMember(companyId: string, userId: number): Promise<void> {
    return apiClient.post(`/api/v1/companies/${companyId}/members/${userId}/remove`, {})
  },

  // Project Members
  async getProjectMembers(projectId: string): Promise<ProjectMember[]> {
    return apiClient.get<ProjectMember[]>(`/api/v1/projects/${projectId}/members`)
  },

  async addProjectMember(projectId: string, data: Partial<ProjectMember>): Promise<ProjectMember> {
    return apiClient.post<ProjectMember>(`/api/v1/projects/${projectId}/members`, data)
  },

  async updateProjectMember(projectId: string, userId: number, data: Partial<ProjectMember>): Promise<ProjectMember> {
    return apiClient.put<ProjectMember>(`/api/v1/projects/${projectId}/members/${userId}`, data)
  },

  async removeProjectMember(projectId: string, userId: number): Promise<void> {
    return apiClient.post(`/api/v1/projects/${projectId}/members/${userId}/remove`, {})
  },
}
