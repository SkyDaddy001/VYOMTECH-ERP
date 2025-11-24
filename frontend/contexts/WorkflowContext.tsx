// Workflow context for managing state - Phase 3B

'use client';

import React, { createContext, useContext, useReducer, ReactNode, useCallback } from 'react';
import { WorkflowDefinition, WorkflowInstance, WorkflowStats } from '@/types/workflow';
import { workflowAPI, workflowExecutionAPI } from '@/services/workflowAPI';

interface WorkflowState {
  workflows: WorkflowDefinition[];
  selectedWorkflow: WorkflowDefinition | null;
  instances: WorkflowInstance[];
  stats: WorkflowStats | null;
  loading: boolean;
  error: string | null;
}

type WorkflowAction =
  | { type: 'SET_WORKFLOWS'; payload: WorkflowDefinition[] }
  | { type: 'SET_SELECTED_WORKFLOW'; payload: WorkflowDefinition | null }
  | { type: 'SET_INSTANCES'; payload: WorkflowInstance[] }
  | { type: 'SET_STATS'; payload: WorkflowStats }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'ADD_WORKFLOW'; payload: WorkflowDefinition }
  | { type: 'UPDATE_WORKFLOW'; payload: WorkflowDefinition }
  | { type: 'DELETE_WORKFLOW'; payload: string }
  | { type: 'ADD_INSTANCE'; payload: WorkflowInstance }
  | { type: 'UPDATE_INSTANCE'; payload: WorkflowInstance }
  | { type: 'RESET' };

const initialState: WorkflowState = {
  workflows: [],
  selectedWorkflow: null,
  instances: [],
  stats: null,
  loading: false,
  error: null,
};

const workflowReducer = (state: WorkflowState, action: WorkflowAction): WorkflowState => {
  switch (action.type) {
    case 'SET_WORKFLOWS':
      return { ...state, workflows: action.payload };
    case 'SET_SELECTED_WORKFLOW':
      return { ...state, selectedWorkflow: action.payload };
    case 'SET_INSTANCES':
      return { ...state, instances: action.payload };
    case 'SET_STATS':
      return { ...state, stats: action.payload };
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    case 'SET_ERROR':
      return { ...state, error: action.payload };
    case 'ADD_WORKFLOW':
      return { ...state, workflows: [...state.workflows, action.payload] };
    case 'UPDATE_WORKFLOW':
      return {
        ...state,
        workflows: state.workflows.map((w) => (w.id === action.payload.id ? action.payload : w)),
        selectedWorkflow:
          state.selectedWorkflow?.id === action.payload.id ? action.payload : state.selectedWorkflow,
      };
    case 'DELETE_WORKFLOW':
      return {
        ...state,
        workflows: state.workflows.filter((w) => w.id !== action.payload),
        selectedWorkflow:
          state.selectedWorkflow?.id === action.payload ? null : state.selectedWorkflow,
      };
    case 'ADD_INSTANCE':
      return { ...state, instances: [...state.instances, action.payload] };
    case 'UPDATE_INSTANCE':
      return {
        ...state,
        instances: state.instances.map((i) =>
          i.id === action.payload.id ? action.payload : i
        ),
      };
    case 'RESET':
      return initialState;
    default:
      return state;
  }
};

interface WorkflowContextType extends WorkflowState {
  loadWorkflows: () => Promise<void>;
  loadWorkflow: (id: string) => Promise<void>;
  createWorkflow: (data: any) => Promise<void>;
  updateWorkflow: (id: string, data: any) => Promise<void>;
  deleteWorkflow: (id: string) => Promise<void>;
  triggerWorkflow: (workflowId: string, data?: any) => Promise<void>;
  loadWorkflowInstances: (workflowId?: string) => Promise<void>;
  loadStats: () => Promise<void>;
  clearError: () => void;
}

const WorkflowContext = createContext<WorkflowContextType | undefined>(undefined);

export const WorkflowProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [state, dispatch] = useReducer(workflowReducer, initialState);

  const loadWorkflows = useCallback(async () => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const workflows = await workflowAPI.listWorkflows();
      dispatch({ type: 'SET_WORKFLOWS', payload: workflows });
      dispatch({ type: 'SET_ERROR', payload: null });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to load workflows' });
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  }, []);

  const loadWorkflow = useCallback(async (id: string) => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const workflow = await workflowAPI.getWorkflow(id);
      dispatch({ type: 'SET_SELECTED_WORKFLOW', payload: workflow });
      dispatch({ type: 'SET_ERROR', payload: null });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to load workflow' });
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  }, []);

  const createWorkflow = useCallback(async (data: any) => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const workflow = await workflowAPI.createWorkflow(data);
      dispatch({ type: 'ADD_WORKFLOW', payload: workflow });
      dispatch({ type: 'SET_ERROR', payload: null });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to create workflow' });
      throw error;
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  }, []);

  const updateWorkflow = useCallback(async (id: string, data: any) => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const workflow = await workflowAPI.updateWorkflow(id, data);
      dispatch({ type: 'UPDATE_WORKFLOW', payload: workflow });
      dispatch({ type: 'SET_ERROR', payload: null });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to update workflow' });
      throw error;
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  }, []);

  const deleteWorkflow = useCallback(async (id: string) => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      await workflowAPI.deleteWorkflow(id);
      dispatch({ type: 'DELETE_WORKFLOW', payload: id });
      dispatch({ type: 'SET_ERROR', payload: null });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to delete workflow' });
      throw error;
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  }, []);

  const triggerWorkflow = useCallback(async (workflowId: string, data?: any) => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const instance = await workflowExecutionAPI.triggerWorkflow({
        workflow_id: workflowId,
        trigger_data: data,
      });
      dispatch({ type: 'ADD_INSTANCE', payload: instance });
      dispatch({ type: 'SET_ERROR', payload: null });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to trigger workflow' });
      throw error;
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  }, []);

  const loadWorkflowInstances = useCallback(async (workflowId?: string) => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const instances = await workflowExecutionAPI.listWorkflowInstances(workflowId);
      dispatch({ type: 'SET_INSTANCES', payload: instances });
      dispatch({ type: 'SET_ERROR', payload: null });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to load workflow instances' });
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  }, []);

  const loadStats = useCallback(async () => {
    try {
      const stats = await workflowAPI.getWorkflowStats();
      dispatch({ type: 'SET_STATS', payload: stats });
    } catch (error) {
      console.error('Failed to load workflow stats:', error);
    }
  }, []);

  const clearError = useCallback(() => {
    dispatch({ type: 'SET_ERROR', payload: null });
  }, []);

  const value: WorkflowContextType = {
    ...state,
    loadWorkflows,
    loadWorkflow,
    createWorkflow,
    updateWorkflow,
    deleteWorkflow,
    triggerWorkflow,
    loadWorkflowInstances,
    loadStats,
    clearError,
  };

  return (
    <WorkflowContext.Provider value={value}>
      {children}
    </WorkflowContext.Provider>
  );
};

export const useWorkflow = () => {
  const context = useContext(WorkflowContext);
  if (!context) {
    throw new Error('useWorkflow must be used within WorkflowProvider');
  }
  return context;
};
