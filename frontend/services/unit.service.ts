import { apiClient } from './api'
import { PropertyUnit, UnitCostSheet } from '@/types/unit'

export const unitService = {
  async getUnits(projectId?: string): Promise<PropertyUnit[]> {
    const url = projectId ? `/api/v1/projects/${projectId}/units` : '/api/v1/units'
    return apiClient.get<PropertyUnit[]>(url)
  },

  async getUnit(id: string): Promise<PropertyUnit> {
    return apiClient.get<PropertyUnit>(`/api/v1/units/${id}`)
  },

  async createUnit(projectId: string, data: Partial<PropertyUnit>): Promise<PropertyUnit> {
    return apiClient.post<PropertyUnit>(`/api/v1/projects/${projectId}/units`, data)
  },

  async updateUnit(id: string, data: Partial<PropertyUnit>): Promise<PropertyUnit> {
    return apiClient.put<PropertyUnit>(`/api/v1/units/${id}`, data)
  },

  async deleteUnit(id: string): Promise<void> {
    return apiClient.post(`/api/v1/units/${id}/delete`, {})
  },

  async getCostSheet(unitId: string): Promise<UnitCostSheet> {
    return apiClient.get<UnitCostSheet>(`/api/v1/units/${unitId}/cost-sheet`)
  },

  async updateCostSheet(unitId: string, data: Partial<UnitCostSheet>): Promise<UnitCostSheet> {
    return apiClient.put<UnitCostSheet>(`/api/v1/units/${unitId}/cost-sheet`, data)
  },

  async getUnitsByStatus(projectId: string, status: string): Promise<PropertyUnit[]> {
    return apiClient.get<PropertyUnit[]>(`/api/v1/projects/${projectId}/units?status=${status}`)
  },
}
