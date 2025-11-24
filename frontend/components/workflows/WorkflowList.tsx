// Workflow List Component - Phase 3B

'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { useWorkflow } from '@/contexts/WorkflowContext';
import { WorkflowDefinition, WorkflowFilter } from '@/types/workflow';

export const WorkflowList: React.FC = () => {
  const {
    workflows,
    loading,
    error,
    loadWorkflows,
    deleteWorkflow,
    clearError,
  } = useWorkflow();

  const [filter, setFilter] = useState<WorkflowFilter>({
    searchTerm: '',
    sortBy: 'updated_at',
    sortOrder: 'desc',
    page: 1,
    limit: 10,
  });

  const [confirmDelete, setConfirmDelete] = useState<string | null>(null);

  useEffect(() => {
    loadWorkflows();
  }, [loadWorkflows]);

  const handleDelete = async (id: string) => {
    try {
      await deleteWorkflow(id);
      setConfirmDelete(null);
    } catch (err) {
      console.error('Failed to delete workflow:', err);
    }
  };

  const handleSearch = (term: string) => {
    setFilter({ ...filter, searchTerm: term, page: 1 });
  };

  const handleSort = (sortBy: 'name' | 'created_at' | 'updated_at') => {
    const newOrder = filter.sortBy === sortBy && filter.sortOrder === 'asc' ? 'desc' : 'asc';
    setFilter({ ...filter, sortBy, sortOrder: newOrder });
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      active: 'bg-green-100 text-green-800',
      inactive: 'bg-gray-100 text-gray-800',
      draft: 'bg-yellow-100 text-yellow-800',
    };
    return styles[status] || 'bg-gray-100 text-gray-800';
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Workflows</h1>
          <p className="text-gray-600 mt-1">Manage automation workflows for your business</p>
        </div>
        <Link
          href="/dashboard/workflows/create"
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
        >
          + New Workflow
        </Link>
      </div>

      {/* Search and Filter */}
      <div className="bg-white rounded-lg shadow p-4">
        <div className="flex gap-4 items-center">
          <input
            type="text"
            placeholder="Search workflows..."
            value={filter.searchTerm || ''}
            onChange={(e) => handleSearch(e.target.value)}
            className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <select
            value={`${filter.sortBy}-${filter.sortOrder}`}
            onChange={(e) => {
              const [sortBy, sortOrder] = e.target.value.split('-');
              setFilter({
                ...filter,
                sortBy: sortBy as any,
                sortOrder: sortOrder as any,
              });
            }}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="updated_at-desc">Latest Updated</option>
            <option value="created_at-desc">Newest First</option>
            <option value="name-asc">Name (A-Z)</option>
          </select>
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg flex justify-between items-center">
          <span>{error}</span>
          <button onClick={clearError} className="text-red-600 hover:text-red-800">
            ✕
          </button>
        </div>
      )}

      {/* Workflows Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-gray-500">
            <div className="inline-block animate-spin text-2xl">⟳</div>
            <p className="mt-2">Loading workflows...</p>
          </div>
        ) : workflows.length === 0 ? (
          <div className="p-8 text-center text-gray-500">
            <p className="mb-2">No workflows found</p>
            <Link
              href="/dashboard/workflows/create"
              className="text-blue-600 hover:text-blue-700 font-medium"
            >
              Create your first workflow
            </Link>
          </div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                  Name
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                  Triggers
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                  Actions
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                  Updated
                </th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-700 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {workflows.map((workflow) => (
                <tr key={workflow.id} className="hover:bg-gray-50 transition">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div>
                      <p className="font-medium text-gray-900">{workflow.name}</p>
                      <p className="text-sm text-gray-500">{workflow.description}</p>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusBadge(workflow.status)}`}>
                      {workflow.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    {workflow.triggers?.length || 0}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    {workflow.actions?.length || 0}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    {new Date(workflow.updated_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right">
                    <div className="flex gap-2 justify-end">
                      <Link
                        href={`/dashboard/workflows/${workflow.id}`}
                        className="text-blue-600 hover:text-blue-700 text-sm font-medium"
                      >
                        Edit
                      </Link>
                      <Link
                        href={`/dashboard/workflows/${workflow.id}/executions`}
                        className="text-indigo-600 hover:text-indigo-700 text-sm font-medium"
                      >
                        Executions
                      </Link>
                      <button
                        onClick={() => setConfirmDelete(workflow.id)}
                        className="text-red-600 hover:text-red-700 text-sm font-medium"
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Delete Confirmation */}
      {confirmDelete && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-sm">
            <h3 className="text-lg font-bold text-gray-900">Delete Workflow?</h3>
            <p className="text-gray-600 mt-2">This action cannot be undone.</p>
            <div className="flex gap-3 mt-6">
              <button
                onClick={() => setConfirmDelete(null)}
                className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                onClick={() => handleDelete(confirmDelete)}
                className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default WorkflowList;
