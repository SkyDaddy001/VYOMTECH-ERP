// Workflow API Service - Phase 3B

import axios from 'axios';
import {
  WorkflowDefinition,
  WorkflowInstance,
  ScheduledTask,
  ScheduledTaskExecution,
  CreateWorkflowRequest,
  UpdateWorkflowRequest,
  WorkflowInstanceRequest,
  WorkflowResponse,
  WorkflowFilter,
  WorkflowStats,
} from '@/types/workflow';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000/api';

// Create axios instance with auth
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
});

// Add token to requests
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Workflow Management
export const workflowAPI = {
  // List workflows
  async listWorkflows(filter?: WorkflowFilter): Promise<WorkflowDefinition[]> {
    try {
      const params = new URLSearchParams();
      if (filter?.status) params.append('status', filter.status);
      if (filter?.searchTerm) params.append('search', filter.searchTerm);
      if (filter?.sortBy) params.append('sort_by', filter.sortBy);
      if (filter?.sortOrder) params.append('sort_order', filter.sortOrder);
      if (filter?.page) params.append('page', filter.page.toString());
      if (filter?.limit) params.append('limit', filter.limit.toString());

      const response = await apiClient.get<WorkflowResponse>('/workflows', { params });
      return response.data.data || [];
    } catch (error) {
      console.error('Failed to list workflows:', error);
      throw error;
    }
  },

  // Get single workflow
  async getWorkflow(workflowId: string): Promise<WorkflowDefinition> {
    try {
      const response = await apiClient.get<WorkflowResponse>(`/workflows/${workflowId}`);
      return response.data.data;
    } catch (error) {
      console.error('Failed to get workflow:', error);
      throw error;
    }
  },

  // Create workflow
  async createWorkflow(request: CreateWorkflowRequest): Promise<WorkflowDefinition> {
    try {
      const response = await apiClient.post<WorkflowResponse>('/workflows', request);
      return response.data.data;
    } catch (error) {
      console.error('Failed to create workflow:', error);
      throw error;
    }
  },

  // Update workflow
  async updateWorkflow(
    workflowId: string,
    request: UpdateWorkflowRequest
  ): Promise<WorkflowDefinition> {
    try {
      const response = await apiClient.put<WorkflowResponse>(`/workflows/${workflowId}`, request);
      return response.data.data;
    } catch (error) {
      console.error('Failed to update workflow:', error);
      throw error;
    }
  },

  // Delete workflow
  async deleteWorkflow(workflowId: string): Promise<void> {
    try {
      await apiClient.delete(`/workflows/${workflowId}`);
    } catch (error) {
      console.error('Failed to delete workflow:', error);
      throw error;
    }
  },

  // Enable/Disable workflow
  async toggleWorkflow(workflowId: string, enabled: boolean): Promise<WorkflowDefinition> {
    try {
      const response = await apiClient.patch<WorkflowResponse>(`/workflows/${workflowId}/toggle`, {
        is_enabled: enabled,
      });
      return response.data.data;
    } catch (error) {
      console.error('Failed to toggle workflow:', error);
      throw error;
    }
  },

  // Get workflow statistics
  async getWorkflowStats(): Promise<WorkflowStats> {
    try {
      const response = await apiClient.get<WorkflowResponse>('/workflows/stats');
      return response.data.data;
    } catch (error) {
      console.error('Failed to get workflow stats:', error);
      throw error;
    }
  },
};

// Workflow Execution
export const workflowExecutionAPI = {
  // Trigger workflow
  async triggerWorkflow(request: WorkflowInstanceRequest): Promise<WorkflowInstance> {
    try {
      const response = await apiClient.post<WorkflowResponse>('/workflow-instances', request);
      return response.data.data;
    } catch (error) {
      console.error('Failed to trigger workflow:', error);
      throw error;
    }
  },

  // Get workflow instance
  async getWorkflowInstance(instanceId: string): Promise<WorkflowInstance> {
    try {
      const response = await apiClient.get<WorkflowResponse>(`/workflow-instances/${instanceId}`);
      return response.data.data;
    } catch (error) {
      console.error('Failed to get workflow instance:', error);
      throw error;
    }
  },

  // List workflow instances
  async listWorkflowInstances(
    workflowId?: string,
    page?: number,
    limit?: number
  ): Promise<WorkflowInstance[]> {
    try {
      const params = new URLSearchParams();
      if (workflowId) params.append('workflow_id', workflowId);
      if (page) params.append('page', page.toString());
      if (limit) params.append('limit', limit.toString());

      const response = await apiClient.get<WorkflowResponse>('/workflow-instances', { params });
      return response.data.data || [];
    } catch (error) {
      console.error('Failed to list workflow instances:', error);
      throw error;
    }
  },

  // Cancel workflow instance
  async cancelWorkflowInstance(instanceId: string): Promise<WorkflowInstance> {
    try {
      const response = await apiClient.post<WorkflowResponse>(
        `/workflow-instances/${instanceId}/cancel`
      );
      return response.data.data;
    } catch (error) {
      console.error('Failed to cancel workflow instance:', error);
      throw error;
    }
  },
};

// Scheduled Tasks
export const scheduledTasksAPI = {
  // List scheduled tasks
  async listScheduledTasks(page?: number, limit?: number): Promise<ScheduledTask[]> {
    try {
      const params = new URLSearchParams();
      if (page) params.append('page', page.toString());
      if (limit) params.append('limit', limit.toString());

      const response = await apiClient.get<WorkflowResponse>('/scheduled-tasks', { params });
      return response.data.data || [];
    } catch (error) {
      console.error('Failed to list scheduled tasks:', error);
      throw error;
    }
  },

  // Get scheduled task
  async getScheduledTask(taskId: string): Promise<ScheduledTask> {
    try {
      const response = await apiClient.get<WorkflowResponse>(`/scheduled-tasks/${taskId}`);
      return response.data.data;
    } catch (error) {
      console.error('Failed to get scheduled task:', error);
      throw error;
    }
  },

  // Create scheduled task
  async createScheduledTask(request: any): Promise<ScheduledTask> {
    try {
      const response = await apiClient.post<WorkflowResponse>('/scheduled-tasks', request);
      return response.data.data;
    } catch (error) {
      console.error('Failed to create scheduled task:', error);
      throw error;
    }
  },

  // Update scheduled task
  async updateScheduledTask(taskId: string, request: any): Promise<ScheduledTask> {
    try {
      const response = await apiClient.put<WorkflowResponse>(`/scheduled-tasks/${taskId}`, request);
      return response.data.data;
    } catch (error) {
      console.error('Failed to update scheduled task:', error);
      throw error;
    }
  },

  // Delete scheduled task
  async deleteScheduledTask(taskId: string): Promise<void> {
    try {
      await apiClient.delete(`/scheduled-tasks/${taskId}`);
    } catch (error) {
      console.error('Failed to delete scheduled task:', error);
      throw error;
    }
  },

  // Enable/Disable scheduled task
  async toggleScheduledTask(taskId: string, enabled: boolean): Promise<ScheduledTask> {
    try {
      const response = await apiClient.patch<WorkflowResponse>(
        `/scheduled-tasks/${taskId}/toggle`,
        { is_enabled: enabled }
      );
      return response.data.data;
    } catch (error) {
      console.error('Failed to toggle scheduled task:', error);
      throw error;
    }
  },

  // Get scheduled task executions
  async getTaskExecutions(
    taskId: string,
    page?: number,
    limit?: number
  ): Promise<ScheduledTaskExecution[]> {
    try {
      const params = new URLSearchParams();
      if (page) params.append('page', page.toString());
      if (limit) params.append('limit', limit.toString());

      const response = await apiClient.get<WorkflowResponse>(
        `/scheduled-tasks/${taskId}/executions`,
        { params }
      );
      return response.data.data || [];
    } catch (error) {
      console.error('Failed to get task executions:', error);
      throw error;
    }
  },
};

export default {
  workflowAPI,
  workflowExecutionAPI,
  scheduledTasksAPI,
};
