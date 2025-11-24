// Workflow Hooks - Phase 3B

'use client';

import { useCallback, useEffect, useState } from 'react';
import { workflowAPI, workflowExecutionAPI } from '@/services/workflowAPI';
import { WorkflowDefinition, WorkflowInstance, WorkflowStats } from '@/types/workflow';

/**
 * Hook for managing workflow data
 */
export const useWorkflowData = (workflowId?: string) => {
  const [workflow, setWorkflow] = useState<WorkflowDefinition | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    if (!workflowId) return;

    setLoading(true);
    try {
      const data = await workflowAPI.getWorkflow(workflowId);
      setWorkflow(data);
      setError(null);
    } catch (err) {
      setError('Failed to load workflow');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [workflowId]);

  useEffect(() => {
    load();
  }, [load]);

  return { workflow, loading, error, reload: load };
};

/**
 * Hook for managing workflow instances
 */
export const useWorkflowInstances = (workflowId?: string, autoRefresh = false) => {
  const [instances, setInstances] = useState<WorkflowInstance[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const data = await workflowExecutionAPI.listWorkflowInstances(workflowId);
      setInstances(data);
      setError(null);
    } catch (err) {
      setError('Failed to load workflow instances');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [workflowId]);

  useEffect(() => {
    load();

    if (autoRefresh) {
      const interval = setInterval(load, 5000);
      return () => clearInterval(interval);
    }
  }, [load, autoRefresh]);

  const trigger = useCallback(
    async (triggerData?: any) => {
      if (!workflowId) throw new Error('Workflow ID is required');

      try {
        const instance = await workflowExecutionAPI.triggerWorkflow({
          workflow_id: workflowId,
          trigger_data: triggerData,
        });
        setInstances((prev) => [instance, ...prev]);
        return instance;
      } catch (err) {
        setError('Failed to trigger workflow');
        throw err;
      }
    },
    [workflowId]
  );

  const cancel = useCallback(async (instanceId: string) => {
    try {
      await workflowExecutionAPI.cancelWorkflowInstance(instanceId);
      setInstances((prev) =>
        prev.map((inst) =>
          inst.id === instanceId ? { ...inst, status: 'cancelled' } : inst
        )
      );
    } catch (err) {
      setError('Failed to cancel instance');
      throw err;
    }
  }, []);

  return { instances, loading, error, reload: load, trigger, cancel };
};

/**
 * Hook for workflow statistics
 */
export const useWorkflowStats = (autoRefresh = false) => {
  const [stats, setStats] = useState<WorkflowStats | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const data = await workflowAPI.getWorkflowStats();
      setStats(data);
      setError(null);
    } catch (err) {
      setError('Failed to load workflow stats');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();

    if (autoRefresh) {
      const interval = setInterval(load, 10000);
      return () => clearInterval(interval);
    }
  }, [load, autoRefresh]);

  return { stats, loading, error, reload: load };
};

/**
 * Hook for getting execution status updates
 */
export const useWorkflowExecutionStatus = (instanceId: string, pollInterval = 2000) => {
  const [instance, setInstance] = useState<WorkflowInstance | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const data = await workflowExecutionAPI.getWorkflowInstance(instanceId);
      setInstance(data);
      setError(null);
    } catch (err) {
      setError('Failed to load instance');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [instanceId]);

  useEffect(() => {
    load();

    // Poll while still running
    const interval = setInterval(async () => {
      const data = await workflowExecutionAPI.getWorkflowInstance(instanceId);
      setInstance(data);

      // Stop polling when complete
      if (data.status !== 'running' && data.status !== 'pending') {
        clearInterval(interval);
      }
    }, pollInterval);

    return () => clearInterval(interval);
  }, [instanceId, pollInterval, load]);

  return { instance, loading, error };
};
