// Scheduled Tasks Component - Phase 3B

'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { scheduledTasksAPI } from '@/services/workflowAPI';
import { ScheduledTask, ScheduledTaskExecution } from '@/types/workflow';

export const ScheduledTasks: React.FC = () => {
  const [tasks, setTasks] = useState<ScheduledTask[]>([]);
  const [executions, setExecutions] = useState<ScheduledTaskExecution[]>([]);
  const [selectedTask, setSelectedTask] = useState<ScheduledTask | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  useEffect(() => {
    loadTasks();
  }, []);

  const loadTasks = async () => {
    setLoading(true);
    try {
      const data = await scheduledTasksAPI.listScheduledTasks();
      setTasks(data);
      setError(null);
    } catch (err) {
      setError('Failed to load scheduled tasks');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSelectTask = async (task: ScheduledTask) => {
    setSelectedTask(task);
    try {
      const data = await scheduledTasksAPI.getTaskExecutions(task.id);
      setExecutions(data);
    } catch (err) {
      console.error('Failed to load executions:', err);
    }
  };

  const handleToggle = async (task: ScheduledTask) => {
    try {
      const updated = await scheduledTasksAPI.toggleScheduledTask(task.id, !task.is_enabled);
      setTasks(tasks.map((t) => (t.id === task.id ? updated : t)));
      if (selectedTask?.id === task.id) {
        setSelectedTask(updated);
      }
    } catch (err) {
      console.error('Failed to toggle task:', err);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this task?')) return;

    try {
      await scheduledTasksAPI.deleteScheduledTask(id);
      setTasks(tasks.filter((t) => t.id !== id));
      if (selectedTask?.id === id) {
        setSelectedTask(null);
      }
    } catch (err) {
      console.error('Failed to delete task:', err);
    }
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      success: 'bg-green-100 text-green-800',
      failed: 'bg-red-100 text-red-800',
      skipped: 'bg-yellow-100 text-yellow-800',
    };
    return styles[status] || 'bg-gray-100 text-gray-800';
  };

  return (
    <div className="grid grid-cols-3 gap-6">
      {/* Tasks List */}
      <div className="col-span-1">
        <div className="bg-white rounded-lg shadow">
          <div className="p-6 border-b border-gray-200">
            <h2 className="text-lg font-bold text-gray-900">Scheduled Tasks</h2>
            <p className="text-sm text-gray-600 mt-1">{tasks.length} tasks</p>
          </div>

          {loading ? (
            <div className="p-6 text-center text-gray-500">
              <div className="inline-block animate-spin">⟳</div>
              <p className="mt-2">Loading...</p>
            </div>
          ) : tasks.length === 0 ? (
            <div className="p-6 text-center text-gray-500">
              <p>No scheduled tasks</p>
            </div>
          ) : (
            <div className="divide-y divide-gray-200">
              {tasks.map((task) => (
                <div
                  key={task.id}
                  onClick={() => handleSelectTask(task)}
                  className={`p-4 cursor-pointer transition ${
                    selectedTask?.id === task.id
                      ? 'bg-blue-50 border-l-4 border-blue-500'
                      : 'hover:bg-gray-50'
                  }`}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <p className="font-medium text-gray-900">{task.task_name}</p>
                      <p className="text-xs text-gray-600 mt-1">{task.task_type}</p>
                    </div>
                    <label
                      onClick={(e) => {
                        e.stopPropagation();
                        handleToggle(task);
                      }}
                      className="flex items-center"
                    >
                      <input
                        type="checkbox"
                        checked={task.is_enabled}
                        onChange={() => {}}
                        className="rounded"
                      />
                    </label>
                  </div>
                </div>
              ))}
            </div>
          )}

          {error && (
            <div className="p-4 bg-red-50 border-t border-red-200 text-red-800 text-sm">
              {error}
            </div>
          )}
        </div>
      </div>

      {/* Task Details and Executions */}
      <div className="col-span-2 space-y-6">
        {selectedTask ? (
          <>
            {/* Task Details */}
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h3 className="text-lg font-bold text-gray-900">{selectedTask.task_name}</h3>
                  <p className="text-gray-600">{selectedTask.task_type}</p>
                </div>
                <button
                  onClick={() => handleDelete(selectedTask.id)}
                  className="text-red-600 hover:text-red-700 text-sm"
                >
                  Delete
                </button>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-gray-600">Status</p>
                  <span
                    className={`mt-1 px-3 py-1 rounded-full text-xs font-medium inline-block ${
                      selectedTask.is_enabled
                        ? 'bg-green-100 text-green-800'
                        : 'bg-gray-100 text-gray-800'
                    }`}
                  >
                    {selectedTask.is_enabled ? 'Enabled' : 'Disabled'}
                  </span>
                </div>

                <div>
                  <p className="text-sm text-gray-600">Cron Expression</p>
                  <code className="mt-1 text-sm font-mono text-gray-900">{selectedTask.cron_expression}</code>
                </div>

                <div>
                  <p className="text-sm text-gray-600">Last Run</p>
                  <p className="mt-1 text-sm text-gray-900">
                    {selectedTask.last_run_at
                      ? new Date(selectedTask.last_run_at).toLocaleString()
                      : 'Never'}
                  </p>
                </div>

                <div>
                  <p className="text-sm text-gray-600">Next Run</p>
                  <p className="mt-1 text-sm text-gray-900">
                    {selectedTask.next_run_at
                      ? new Date(selectedTask.next_run_at).toLocaleString()
                      : 'Not scheduled'}
                  </p>
                </div>

                <div>
                  <p className="text-sm text-gray-600">Max Retries</p>
                  <p className="mt-1 text-sm text-gray-900">{selectedTask.max_retries}</p>
                </div>
              </div>
            </div>

            {/* Execution History */}
            <div className="bg-white rounded-lg shadow p-6">
              <h4 className="font-bold text-gray-900 mb-4">Recent Executions</h4>
              {executions.length === 0 ? (
                <p className="text-gray-600 text-sm">No executions yet</p>
              ) : (
                <div className="space-y-2">
                  {executions.map((exec) => (
                    <div key={exec.id} className="p-3 bg-gray-50 rounded-lg">
                      <div className="flex justify-between items-start">
                        <div>
                          <span className={`px-2 py-1 rounded text-xs font-medium ${getStatusBadge(exec.status)}`}>
                            {exec.status}
                          </span>
                          <p className="text-xs text-gray-600 mt-1">
                            {new Date(exec.executed_at).toLocaleString()}
                          </p>
                        </div>
                        <div className="text-right text-xs text-gray-600">
                          {exec.duration_ms}ms
                        </div>
                      </div>
                      {exec.error_message && (
                        <p className="text-xs text-red-600 mt-2">{exec.error_message}</p>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          </>
        ) : (
          <div className="bg-white rounded-lg shadow p-12 text-center text-gray-500">
            <p>Select a scheduled task to view details</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default ScheduledTasks;
