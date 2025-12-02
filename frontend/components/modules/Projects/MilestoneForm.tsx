'use client'

import React, { useState } from 'react'
import { ProjectMilestone } from '@/types/projects'

interface MilestoneFormProps {
  milestone?: ProjectMilestone
  onSubmit: (data: Partial<ProjectMilestone>) => Promise<void>
  onCancel: () => void
}

export default function MilestoneForm({ milestone, onSubmit, onCancel }: MilestoneFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<ProjectMilestone>>(
    milestone || {
      milestone_name: '',
      milestone_type: 'foundation',
      planned_date: '',
      completion_percentage: 0,
      status: 'pending',
      description: '',
    }
  )

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await onSubmit(formData)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4 max-w-2xl">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Milestone Name *</label>
          <input
            type="text"
            required
            value={formData.milestone_name || ''}
            onChange={(e) => setFormData({ ...formData, milestone_name: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Type *</label>
          <select
            required
            value={formData.milestone_type || 'foundation'}
            onChange={(e) => setFormData({ ...formData, milestone_type: e.target.value as any })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="foundation">Foundation</option>
            <option value="structure">Structure</option>
            <option value="finishing">Finishing</option>
            <option value="handover">Handover</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Planned Date *</label>
          <input
            type="date"
            required
            value={formData.planned_date || ''}
            onChange={(e) => setFormData({ ...formData, planned_date: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Status *</label>
          <select
            required
            value={formData.status || 'pending'}
            onChange={(e) => setFormData({ ...formData, status: e.target.value as any })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="pending">Pending</option>
            <option value="in_progress">In Progress</option>
            <option value="completed">Completed</option>
            <option value="delayed">Delayed</option>
          </select>
        </div>

        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-700 mb-1">Completion % *</label>
          <input
            type="number"
            required
            min="0"
            max="100"
            value={formData.completion_percentage || 0}
            onChange={(e) => setFormData({ ...formData, completion_percentage: parseInt(e.target.value) })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
          <textarea
            value={formData.description || ''}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            rows={3}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <div className="flex gap-3">
        <button
          type="submit"
          disabled={loading}
          className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium disabled:opacity-50"
        >
          {loading ? 'Saving...' : milestone ? 'Update Milestone' : 'Create Milestone'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 font-medium"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
