// Workflow Editor Component - Phase 3B

'use client';

import React, { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useWorkflow } from '@/contexts/WorkflowContext';
import {
  WorkflowDefinition,
  CreateWorkflowRequest,
  TriggerType,
  ActionType,
} from '@/types/workflow';

interface EditorState {
  name: string;
  description: string;
  triggers: any[];
  actions: any[];
  errors: Record<string, string>;
}

export const WorkflowEditor: React.FC<{ workflowId?: string }> = ({ workflowId }) => {
  const router = useRouter();
  const { createWorkflow, updateWorkflow, selectedWorkflow, loadWorkflow } = useWorkflow();

  const [state, setState] = useState<EditorState>({
    name: '',
    description: '',
    triggers: [],
    actions: [],
    errors: {},
  });

  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState<'basic' | 'triggers' | 'actions'>('basic');

  useEffect(() => {
    if (workflowId) {
      loadWorkflow(workflowId);
    }
  }, [workflowId, loadWorkflow]);

  useEffect(() => {
    if (selectedWorkflow && workflowId) {
      setState({
        name: selectedWorkflow.name,
        description: selectedWorkflow.description || '',
        triggers: selectedWorkflow.triggers || [],
        actions: selectedWorkflow.actions || [],
        errors: {},
      });
    }
  }, [selectedWorkflow, workflowId]);

  const handleFieldChange = (field: string, value: any) => {
    setState((prev) => ({
      ...prev,
      [field]: value,
      errors: { ...prev.errors, [field]: '' },
    }));
  };

  const addTrigger = () => {
    setState((prev) => ({
      ...prev,
      triggers: [
        ...prev.triggers,
        {
          id: Date.now().toString(),
          trigger_type: 'lead_created' as TriggerType,
          trigger_config: {},
          conditions: [],
        },
      ],
    }));
  };

  const removeTrigger = (index: number) => {
    setState((prev) => ({
      ...prev,
      triggers: prev.triggers.filter((_, i) => i !== index),
    }));
  };

  const updateTrigger = (index: number, field: string, value: any) => {
    setState((prev) => ({
      ...prev,
      triggers: prev.triggers.map((t, i) => (i === index ? { ...t, [field]: value } : t)),
    }));
  };

  const addAction = () => {
    setState((prev) => ({
      ...prev,
      actions: [
        ...prev.actions,
        {
          id: Date.now().toString(),
          action_type: 'create_task' as ActionType,
          action_config: {},
          action_order: prev.actions.length + 1,
          delay_seconds: 0,
        },
      ],
    }));
  };

  const removeAction = (index: number) => {
    setState((prev) => ({
      ...prev,
      actions: prev.actions.filter((_, i) => i !== index),
    }));
  };

  const updateAction = (index: number, field: string, value: any) => {
    setState((prev) => ({
      ...prev,
      actions: prev.actions.map((a, i) => (i === index ? { ...a, [field]: value } : a)),
    }));
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!state.name.trim()) {
      errors.name = 'Workflow name is required';
    }

    if (state.triggers.length === 0) {
      errors.triggers = 'At least one trigger is required';
    }

    if (state.actions.length === 0) {
      errors.actions = 'At least one action is required';
    }

    setState((prev) => ({ ...prev, errors }));
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) {
      return;
    }

    setLoading(true);
    try {
      const request: CreateWorkflowRequest = {
        name: state.name,
        description: state.description,
        triggers: state.triggers.map((t) => ({
          trigger_type: t.trigger_type,
          trigger_config: t.trigger_config,
          conditions: t.conditions,
        })),
        actions: state.actions.map((a) => ({
          action_type: a.action_type,
          action_config: a.action_config,
          action_order: a.action_order,
          delay_seconds: a.delay_seconds,
        })),
      };

      if (workflowId) {
        await updateWorkflow(workflowId, request);
      } else {
        await createWorkflow(request);
      }

      router.push('/dashboard/workflows');
    } catch (error) {
      setState((prev) => ({
        ...prev,
        errors: { submit: 'Failed to save workflow' },
      }));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">
          {workflowId ? 'Edit Workflow' : 'Create Workflow'}
        </h1>
        <p className="text-gray-600 mt-1">
          {workflowId
            ? 'Update your workflow automation rules'
            : 'Build a new workflow automation'}
        </p>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-lg shadow border-b border-gray-200">
        <div className="flex">
          {(['basic', 'triggers', 'actions'] as const).map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`flex-1 py-4 px-6 font-medium text-center border-b-2 transition ${
                activeTab === tab
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900'
              }`}
            >
              {tab === 'basic' && 'Basic Info'}
              {tab === 'triggers' && 'Triggers'}
              {tab === 'actions' && 'Actions'}
            </button>
          ))}
        </div>
      </div>

      {/* Form */}
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Basic Info Tab */}
        {activeTab === 'basic' && (
          <div className="bg-white rounded-lg shadow p-6 space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-900">Workflow Name</label>
              <input
                type="text"
                value={state.name}
                onChange={(e) => handleFieldChange('name', e.target.value)}
                placeholder="e.g., Lead Scoring Automation"
                className={`mt-2 w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 ${
                  state.errors.name ? 'border-red-500 focus:ring-red-500' : 'border-gray-300 focus:ring-blue-500'
                }`}
              />
              {state.errors.name && (
                <p className="text-red-600 text-sm mt-1">{state.errors.name}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-900">Description</label>
              <textarea
                value={state.description}
                onChange={(e) => handleFieldChange('description', e.target.value)}
                placeholder="Describe what this workflow does..."
                rows={4}
                className="mt-2 w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>
        )}

        {/* Triggers Tab */}
        {activeTab === 'triggers' && (
          <div className="space-y-4">
            {state.errors.triggers && (
              <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg">
                {state.errors.triggers}
              </div>
            )}

            {state.triggers.map((trigger, index) => (
              <div key={trigger.id} className="bg-white rounded-lg shadow p-6">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="font-medium text-gray-900">Trigger {index + 1}</h3>
                  <button
                    type="button"
                    onClick={() => removeTrigger(index)}
                    className="text-red-600 hover:text-red-700 text-sm"
                  >
                    Remove
                  </button>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-900">Event Type</label>
                  <select
                    value={trigger.trigger_type}
                    onChange={(e) => updateTrigger(index, 'trigger_type', e.target.value)}
                    className="mt-2 w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="lead_created">Lead Created</option>
                    <option value="lead_scored">Lead Scored</option>
                    <option value="task_completed">Task Completed</option>
                    <option value="custom_event">Custom Event</option>
                  </select>
                </div>
              </div>
            ))}

            <button
              type="button"
              onClick={addTrigger}
              className="w-full py-2 border-2 border-dashed border-gray-300 text-gray-600 rounded-lg hover:border-gray-400 hover:text-gray-700 transition"
            >
              + Add Trigger
            </button>
          </div>
        )}

        {/* Actions Tab */}
        {activeTab === 'actions' && (
          <div className="space-y-4">
            {state.errors.actions && (
              <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg">
                {state.errors.actions}
              </div>
            )}

            {state.actions.map((action, index) => (
              <div key={action.id} className="bg-white rounded-lg shadow p-6">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="font-medium text-gray-900">Action {index + 1}</h3>
                  <button
                    type="button"
                    onClick={() => removeAction(index)}
                    className="text-red-600 hover:text-red-700 text-sm"
                  >
                    Remove
                  </button>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-900">Action Type</label>
                    <select
                      value={action.action_type}
                      onChange={(e) => updateAction(index, 'action_type', e.target.value)}
                      className="mt-2 w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      <option value="create_task">Create Task</option>
                      <option value="send_email">Send Email</option>
                      <option value="send_sms">Send SMS</option>
                      <option value="send_notification">Send Notification</option>
                      <option value="update_lead">Update Lead</option>
                      <option value="add_tag">Add Tag</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-900">Delay (seconds)</label>
                    <input
                      type="number"
                      min="0"
                      value={action.delay_seconds || 0}
                      onChange={(e) =>
                        updateAction(index, 'delay_seconds', parseInt(e.target.value))
                      }
                      className="mt-2 w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                </div>
              </div>
            ))}

            <button
              type="button"
              onClick={addAction}
              className="w-full py-2 border-2 border-dashed border-gray-300 text-gray-600 rounded-lg hover:border-gray-400 hover:text-gray-700 transition"
            >
              + Add Action
            </button>
          </div>
        )}

        {/* Submit */}
        <div className="flex gap-3">
          <button
            type="button"
            onClick={() => router.back()}
            className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={loading}
            className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
          >
            {loading ? 'Saving...' : workflowId ? 'Update Workflow' : 'Create Workflow'}
          </button>
        </div>
      </form>
    </div>
  );
};

export default WorkflowEditor;
