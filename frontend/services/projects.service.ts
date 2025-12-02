import { apiClient } from './api'
import { Project, ProjectBlock, ProjectMilestone, ProjectMetrics, ProjectTimeline } from '@/types/projects'

export const projectsService = {
  // Projects
  async getProjects(): Promise<Project[]> {
    return apiClient.get<Project[]>('/api/v1/projects')
  },

  async getProject(id: string): Promise<Project> {
    return apiClient.get<Project>(`/api/v1/projects/${id}`)
  },

  async createProject(data: Partial<Project>): Promise<Project> {
    return apiClient.post<Project>('/api/v1/projects', data)
  },

  async updateProject(id: string, data: Partial<Project>): Promise<Project> {
    return apiClient.put<Project>(`/api/v1/projects/${id}`, data)
  },

  async deleteProject(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/projects/${id}`)
  },

  // Blocks
  async getProjectBlocks(projectId: string): Promise<ProjectBlock[]> {
    return apiClient.get<ProjectBlock[]>(`/api/v1/projects/${projectId}/blocks`)
  },

  async createBlock(projectId: string, data: Partial<ProjectBlock>): Promise<ProjectBlock> {
    return apiClient.post<ProjectBlock>(`/api/v1/projects/${projectId}/blocks`, data)
  },

  async updateBlock(projectId: string, blockId: string, data: Partial<ProjectBlock>): Promise<ProjectBlock> {
    return apiClient.put<ProjectBlock>(`/api/v1/projects/${projectId}/blocks/${blockId}`, data)
  },

  async deleteBlock(projectId: string, blockId: string): Promise<void> {
    return apiClient.delete(`/api/v1/projects/${projectId}/blocks/${blockId}`)
  },

  // Milestones
  async getMilestones(projectId: string): Promise<ProjectMilestone[]> {
    return apiClient.get<ProjectMilestone[]>(`/api/v1/projects/${projectId}/milestones`)
  },

  async createMilestone(projectId: string, data: Partial<ProjectMilestone>): Promise<ProjectMilestone> {
    return apiClient.post<ProjectMilestone>(`/api/v1/projects/${projectId}/milestones`, data)
  },

  async updateMilestone(projectId: string, milestoneId: string, data: Partial<ProjectMilestone>): Promise<ProjectMilestone> {
    return apiClient.put<ProjectMilestone>(`/api/v1/projects/${projectId}/milestones/${milestoneId}`, data)
  },

  async deleteMilestone(projectId: string, milestoneId: string): Promise<void> {
    return apiClient.delete(`/api/v1/projects/${projectId}/milestones/${milestoneId}`)
  },

  // Metrics
  async getMetrics(): Promise<ProjectMetrics> {
    return apiClient.get<ProjectMetrics>('/api/v1/projects/metrics/overview')
  },

  async getTimeline(projectId: string): Promise<ProjectTimeline[]> {
    return apiClient.get<ProjectTimeline[]>(`/api/v1/projects/${projectId}/timeline`)
  },
}
