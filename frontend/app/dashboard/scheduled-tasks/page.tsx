'use client'

import { useState } from 'react'
import { Clock, CheckCircle, AlertCircle, PlayCircle, Pause } from 'lucide-react'
import type { ScheduledTask, TaskExecution, TaskTemplate, ScheduledTasksDashboard } from '@/types/scheduledTasks'

type TabType = 'dashboard' | 'tasks' | 'executions' | 'templates' | 'notifications'

export default function ScheduledTasksPage() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard')

  // Mock data
  const dashboard: ScheduledTasksDashboard = {
    total_tasks: 28,
    active_tasks: 22,
    pending_tasks: 4,
    completed_tasks_today: 12,
    failed_tasks: 1,
    overdue_tasks: 2,
    completion_rate_percentage: 94,
    average_execution_time_minutes: 8,
    upcoming_due_in_24hrs: 6,
  }

  const tasks: ScheduledTask[] = [
    {
      id: '1',
      task_name: 'Daily Report Generation',
      task_description: 'Generate daily call center performance reports',
      frequency: 'daily',
      scheduled_date: '2024-11-29',
      scheduled_time: '09:00',
      next_execution: '2024-11-30 09:00',
      last_execution: '2024-11-29 09:15',
      execution_status: 'completed',
      assigned_to: 'System',
      priority: 'high',
      category: 'reporting',
      is_active: true,
    },
    {
      id: '2',
      task_name: 'Database Backup',
      task_description: 'Automated daily backup of production database',
      frequency: 'daily',
      scheduled_date: '2024-11-29',
      scheduled_time: '23:00',
      next_execution: '2024-11-29 23:00',
      last_execution: '2024-11-28 23:05',
      execution_status: 'pending',
      assigned_to: 'System',
      priority: 'critical',
      category: 'maintenance',
      is_active: true,
    },
    {
      id: '3',
      task_name: 'Agent Performance Review',
      task_description: 'Weekly agent performance metrics calculation',
      frequency: 'weekly',
      scheduled_date: '2024-12-02',
      scheduled_time: '18:00',
      next_execution: '2024-12-02 18:00',
      last_execution: '2024-11-25 18:20',
      execution_status: 'pending',
      assigned_to: 'HR Manager',
      priority: 'high',
      category: 'hr',
      is_active: true,
    },
  ]

  const executions: TaskExecution[] = [
    {
      id: '1',
      task_id: '1',
      task_name: 'Daily Report Generation',
      execution_start: '2024-11-29 09:00',
      execution_end: '2024-11-29 09:15',
      duration_minutes: 15,
      status: 'completed',
      execution_result: 'Successfully generated 28 reports',
      executed_by: 'System',
      created_at: '2024-11-29 09:00',
    },
    {
      id: '2',
      task_id: '3',
      task_name: 'Agent Performance Review',
      execution_start: '2024-11-25 18:00',
      execution_end: '2024-11-25 18:20',
      duration_minutes: 20,
      status: 'completed',
      execution_result: 'Performance metrics updated for 156 agents',
      executed_by: 'HR Manager',
      created_at: '2024-11-25 18:00',
    },
    {
      id: '3',
      task_id: '2',
      task_name: 'Database Backup',
      execution_start: '2024-11-28 23:00',
      execution_end: '2024-11-28 23:45',
      duration_minutes: 45,
      status: 'failed',
      execution_result: 'Failed',
      error_message: 'Database lock timeout - retrying at next scheduled time',
      executed_by: 'System',
      created_at: '2024-11-28 23:00',
    },
  ]

  const templates: TaskTemplate[] = [
    {
      id: '1',
      template_name: 'Daily Report Template',
      description: 'Template for generating daily performance reports',
      task_type: 'automated',
      default_assignee: 'System',
      default_priority: 'high',
      estimated_duration_minutes: 15,
      checklist_items: ['Collect metrics', 'Generate charts', 'Send notifications'],
      created_at: '2024-01-15',
    },
    {
      id: '2',
      template_name: 'Weekly Team Meeting',
      description: 'Template for scheduling weekly team meetings',
      task_type: 'manual',
      default_assignee: 'Team Lead',
      default_priority: 'medium',
      estimated_duration_minutes: 60,
      created_at: '2024-02-01',
    },
  ]

  const tabs: Array<{ id: TabType; label: string }> = [
    { id: 'dashboard', label: 'Dashboard' },
    { id: 'tasks', label: 'Scheduled Tasks' },
    { id: 'executions', label: 'Execution History' },
    { id: 'templates', label: 'Templates' },
  ]

  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-indigo-600 to-indigo-800 rounded-lg p-6 text-white">
        <h1 className="text-3xl font-bold">Scheduled Tasks</h1>
        <p className="text-indigo-100 mt-2">Manage background jobs, automation workflows, and recurring tasks</p>
      </div>

      <div className="flex gap-2 border-b border-gray-200 overflow-x-auto">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`px-4 py-3 font-medium border-b-2 transition whitespace-nowrap ${
              activeTab === tab.id
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'dashboard' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Active Tasks</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.active_tasks}</p>
                </div>
                <PlayCircle className="text-green-600" size={32} />
              </div>
              <p className="text-gray-600 text-sm mt-2">of {dashboard.total_tasks} total</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Completion Rate</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.completion_rate_percentage}%</p>
                </div>
                <CheckCircle className="text-blue-600" size={32} />
              </div>
              <p className="text-blue-600 text-sm mt-2">Excellent</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Overdue Tasks</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.overdue_tasks}</p>
                </div>
                <AlertCircle className="text-red-600" size={32} />
              </div>
              <p className="text-red-600 text-sm mt-2">Requires attention</p>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-gray-500 text-sm">Avg Execution Time</p>
                  <p className="text-3xl font-bold text-gray-900">{dashboard.average_execution_time_minutes}m</p>
                </div>
                <Clock className="text-purple-600" size={32} />
              </div>
              <p className="text-gray-600 text-sm mt-2">Per task</p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="bg-white rounded-lg shadow p-6">
              <h3 className="font-semibold text-gray-900 mb-4">Status Breakdown</h3>
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-gray-600">Completed Today</span>
                  <span className="font-semibold text-green-600">{dashboard.completed_tasks_today}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-600">Pending Execution</span>
                  <span className="font-semibold text-yellow-600">{dashboard.pending_tasks}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-600">Failed</span>
                  <span className="font-semibold text-red-600">{dashboard.failed_tasks}</span>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <h3 className="font-semibold text-gray-900 mb-4">Upcoming Executions</h3>
              <p className="text-3xl font-bold text-indigo-600 mb-2">{dashboard.upcoming_due_in_24hrs}</p>
              <p className="text-gray-600">Tasks due in next 24 hours</p>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'tasks' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Task Name</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Frequency</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Next Execution</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Priority</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
              </tr>
            </thead>
            <tbody>
              {tasks.map((task) => (
                <tr key={task.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm text-gray-900">{task.task_name}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{task.frequency}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{task.next_execution}</td>
                  <td className="px-6 py-4">
                    <span
                      className={`px-3 py-1 rounded-full text-xs font-medium ${
                        task.priority === 'critical'
                          ? 'bg-red-100 text-red-800'
                          : task.priority === 'high'
                            ? 'bg-orange-100 text-orange-800'
                            : 'bg-yellow-100 text-yellow-800'
                      }`}
                    >
                      {task.priority}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <span
                      className={`px-3 py-1 rounded-full text-xs font-medium ${
                        task.execution_status === 'completed'
                          ? 'bg-green-100 text-green-800'
                          : task.execution_status === 'pending'
                            ? 'bg-blue-100 text-blue-800'
                            : 'bg-red-100 text-red-800'
                      }`}
                    >
                      {task.execution_status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeTab === 'executions' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Task Name</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Start Time</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Duration</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Result</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
              </tr>
            </thead>
            <tbody>
              {executions.map((execution) => (
                <tr key={execution.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm text-gray-900">{execution.task_name}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{execution.execution_start}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{execution.duration_minutes || '--'} min</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{execution.execution_result}</td>
                  <td className="px-6 py-4">
                    <span
                      className={`px-3 py-1 rounded-full text-xs font-medium ${
                        execution.status === 'completed'
                          ? 'bg-green-100 text-green-800'
                          : execution.status === 'running'
                            ? 'bg-blue-100 text-blue-800'
                            : 'bg-red-100 text-red-800'
                      }`}
                    >
                      {execution.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeTab === 'templates' && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {templates.map((template) => (
            <div key={template.id} className="bg-white rounded-lg shadow p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-2">{template.template_name}</h3>
              <p className="text-gray-600 text-sm mb-4">{template.description}</p>
              <div className="space-y-2 mb-4">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Type:</span>
                  <span className="font-medium text-gray-900">{template.task_type}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Est. Duration:</span>
                  <span className="font-medium text-gray-900">{template.estimated_duration_minutes} min</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Default Priority:</span>
                  <span className="font-medium text-gray-900">{template.default_priority}</span>
                </div>
              </div>
              {template.checklist_items && template.checklist_items.length > 0 && (
                <div className="pt-4 border-t">
                  <p className="text-sm font-medium text-gray-900 mb-2">Checklist:</p>
                  <ul className="space-y-1">
                    {template.checklist_items.map((item, index) => (
                      <li key={index} className="text-sm text-gray-600 flex items-center gap-2">
                        <span className="text-green-600">✓</span> {item}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
