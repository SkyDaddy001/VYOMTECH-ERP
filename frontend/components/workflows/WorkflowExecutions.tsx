// Workflow Executions Viewer Component - Phase 3B

'use client';

import React, { useEffect, useState } from 'react';
import { useWorkflow } from '@/contexts/WorkflowContext';
import { WorkflowInstance } from '@/types/workflow';
import { workflowExecutionAPI } from '@/services/workflowAPI';

export const WorkflowExecutions: React.FC<{ workflowId: string }> = ({ workflowId }) => {
  const { loadWorkflowInstances, instances, loading } = useWorkflow();
  const [selectedInstance, setSelectedInstance] = useState<WorkflowInstance | null>(null);
  const [autoRefresh, setAutoRefresh] = useState(true);

  useEffect(() => {
    loadWorkflowInstances(workflowId);

    if (autoRefresh) {
      const interval = setInterval(() => {
        loadWorkflowInstances(workflowId);
      }, 5000);

      return () => clearInterval(interval);
    }
  }, [workflowId, loadWorkflowInstances, autoRefresh]);

  const handleCancel = async (instanceId: string) => {
    try {
      await workflowExecutionAPI.cancelWorkflowInstance(instanceId);
      loadWorkflowInstances(workflowId);
    } catch (error) {
      console.error('Failed to cancel instance:', error);
    }
  };

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-800',
      running: 'bg-blue-100 text-blue-800',
      completed: 'bg-green-100 text-green-800',
      failed: 'bg-red-100 text-red-800',
      cancelled: 'bg-gray-100 text-gray-800',
    };
    return colors[status] || 'bg-gray-100 text-gray-800';
  };

  const getProgressColor = (percentage: number) => {
    if (percentage === 100) return 'bg-green-500';
    if (percentage >= 75) return 'bg-blue-500';
    if (percentage >= 50) return 'bg-yellow-500';
    return 'bg-orange-500';
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Workflow Executions</h2>
          <p className="text-gray-600 mt-1">Monitor workflow runs and their status</p>
        </div>
        <label className="flex items-center gap-2 text-gray-600">
          <input
            type="checkbox"
            checked={autoRefresh}
            onChange={(e) => setAutoRefresh(e.target.checked)}
            className="rounded"
          />
          Auto-refresh
        </label>
      </div>

      {/* Executions Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-gray-500">
            <div className="inline-block animate-spin text-2xl">⟳</div>
            <p className="mt-2">Loading executions...</p>
          </div>
        ) : instances.length === 0 ? (
          <div className="p-8 text-center text-gray-500">
            <p>No executions found for this workflow</p>
          </div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">
                  Instance ID
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">
                  Progress
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">
                  Actions
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase">
                  Started
                </th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-700 uppercase">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {instances.map((instance) => (
                <tr
                  key={instance.id}
                  className="hover:bg-gray-50 transition cursor-pointer"
                  onClick={() => setSelectedInstance(instance)}
                >
                  <td className="px-6 py-4">
                    <code className="text-sm text-gray-600 font-mono">{instance.id.substring(0, 12)}...</code>
                  </td>
                  <td className="px-6 py-4">
                    <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(instance.status)}`}>
                      {instance.status}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <div className="w-32">
                      <div className="flex items-center justify-between mb-1">
                        <span className="text-sm font-medium text-gray-600">
                          {instance.progress_percentage}%
                        </span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div
                          className={`h-2 rounded-full transition-all ${getProgressColor(
                            instance.progress_percentage
                          )}`}
                          style={{ width: `${instance.progress_percentage}%` }}
                        />
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <div className="flex gap-4">
                      <div>
                        <span className="font-medium text-gray-900">{instance.executed_actions}</span>
                        <span className="text-gray-600 ml-1">executed</span>
                      </div>
                      {instance.failed_actions > 0 && (
                        <div>
                          <span className="font-medium text-red-600">{instance.failed_actions}</span>
                          <span className="text-gray-600 ml-1">failed</span>
                        </div>
                      )}
                    </div>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {instance.started_at
                      ? new Date(instance.started_at).toLocaleString()
                      : 'Not started'}
                  </td>
                  <td className="px-6 py-4 text-right">
                    {instance.status === 'running' && (
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleCancel(instance.id);
                        }}
                        className="text-red-600 hover:text-red-700 text-sm font-medium"
                      >
                        Cancel
                      </button>
                    )}
                    {instance.status === 'completed' && (
                      <span className="text-green-600 text-sm font-medium">✓ Done</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Detail Panel */}
      {selectedInstance && (
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex justify-between items-start mb-6">
            <div>
              <h3 className="text-lg font-bold text-gray-900">Execution Details</h3>
              <code className="text-sm text-gray-600 font-mono">{selectedInstance.id}</code>
            </div>
            <button
              onClick={() => setSelectedInstance(null)}
              className="text-gray-400 hover:text-gray-600"
            >
              ✕
            </button>
          </div>

          <div className="grid grid-cols-2 gap-6">
            <div>
              <p className="text-sm text-gray-600">Status</p>
              <span
                className={`mt-1 px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(
                  selectedInstance.status
                )}`}
              >
                {selectedInstance.status}
              </span>
            </div>

            <div>
              <p className="text-sm text-gray-600">Progress</p>
              <p className="mt-1 text-lg font-semibold text-gray-900">{selectedInstance.progress_percentage}%</p>
            </div>

            <div>
              <p className="text-sm text-gray-600">Actions Executed</p>
              <p className="mt-1 text-lg font-semibold text-gray-900">{selectedInstance.executed_actions}</p>
            </div>

            <div>
              <p className="text-sm text-gray-600">Actions Failed</p>
              <p className={`mt-1 text-lg font-semibold ${selectedInstance.failed_actions > 0 ? 'text-red-600' : 'text-gray-900'}`}>
                {selectedInstance.failed_actions}
              </p>
            </div>

            <div>
              <p className="text-sm text-gray-600">Started At</p>
              <p className="mt-1 text-sm text-gray-900">
                {selectedInstance.started_at
                  ? new Date(selectedInstance.started_at).toLocaleString()
                  : 'Not started'}
              </p>
            </div>

            <div>
              <p className="text-sm text-gray-600">Completed At</p>
              <p className="mt-1 text-sm text-gray-900">
                {selectedInstance.completed_at
                  ? new Date(selectedInstance.completed_at).toLocaleString()
                  : 'In progress'}
              </p>
            </div>
          </div>

          {/* Action Executions */}
          {selectedInstance.action_executions && selectedInstance.action_executions.length > 0 && (
            <div className="mt-6 pt-6 border-t border-gray-200">
              <h4 className="font-medium text-gray-900 mb-4">Action Executions</h4>
              <div className="space-y-3">
                {selectedInstance.action_executions.map((action) => (
                  <div key={action.id} className="p-3 bg-gray-50 rounded-lg">
                    <div className="flex justify-between items-start">
                      <div>
                        <p className="font-medium text-gray-900">{action.action_id}</p>
                        <p className="text-sm text-gray-600">
                          <span className={`px-2 py-0.5 rounded text-xs font-medium ${getStatusColor(action.status)}`}>
                            {action.status}
                          </span>
                        </p>
                      </div>
                      <div className="text-right">
                        <p className="text-sm text-gray-600">Retries: {action.retry_count}</p>
                        {action.error_message && (
                          <p className="text-sm text-red-600 mt-1">{action.error_message}</p>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default WorkflowExecutions;
