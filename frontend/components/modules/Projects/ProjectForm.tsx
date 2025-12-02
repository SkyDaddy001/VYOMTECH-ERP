'use client'

import React, { useState } from 'react'
import { Project } from '@/types/projects'

interface ProjectFormProps {
  project?: Project
  onSubmit: (data: Partial<Project>) => Promise<void>
  onCancel: () => void
}

export default function ProjectForm({ project, onSubmit, onCancel }: ProjectFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<Project>>(
    project || {
      project_name: '',
      project_code: '',
      location: '',
      city: '',
      state: '',
      postal_code: '',
      project_type: 'residential',
      total_units: 0,
      total_area: 0,
      launch_date: '',
      expected_completion: '',
      status: 'planning',
      developer_name: '',
      noc_status: 'pending',
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
    <form onSubmit={handleSubmit} className="space-y-6 max-w-2xl">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Project Name *</label>
          <input
            type="text"
            required
            value={formData.project_name || ''}
            onChange={(e) => setFormData({ ...formData, project_name: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Project Code *</label>
          <input
            type="text"
            required
            value={formData.project_code || ''}
            onChange={(e) => setFormData({ ...formData, project_code: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Type *</label>
          <select
            required
            value={formData.project_type || 'residential'}
            onChange={(e) => setFormData({ ...formData, project_type: e.target.value as any })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="residential">Residential</option>
            <option value="commercial">Commercial</option>
            <option value="mixed">Mixed</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Developer Name *</label>
          <input
            type="text"
            required
            value={formData.developer_name || ''}
            onChange={(e) => setFormData({ ...formData, developer_name: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Location *</label>
          <input
            type="text"
            required
            value={formData.location || ''}
            onChange={(e) => setFormData({ ...formData, location: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">City *</label>
          <input
            type="text"
            required
            value={formData.city || ''}
            onChange={(e) => setFormData({ ...formData, city: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">State *</label>
          <input
            type="text"
            required
            value={formData.state || ''}
            onChange={(e) => setFormData({ ...formData, state: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Postal Code</label>
          <input
            type="text"
            value={formData.postal_code || ''}
            onChange={(e) => setFormData({ ...formData, postal_code: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Total Units *</label>
          <input
            type="number"
            required
            value={formData.total_units || 0}
            onChange={(e) => setFormData({ ...formData, total_units: parseInt(e.target.value) })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Total Area (sqft) *</label>
          <input
            type="number"
            required
            value={formData.total_area || 0}
            onChange={(e) => setFormData({ ...formData, total_area: parseInt(e.target.value) })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Launch Date *</label>
          <input
            type="date"
            required
            value={formData.launch_date || ''}
            onChange={(e) => setFormData({ ...formData, launch_date: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Expected Completion *</label>
          <input
            type="date"
            required
            value={formData.expected_completion || ''}
            onChange={(e) => setFormData({ ...formData, expected_completion: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Status *</label>
          <select
            required
            value={formData.status || 'planning'}
            onChange={(e) => setFormData({ ...formData, status: e.target.value as any })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="planning">Planning</option>
            <option value="approved">Approved</option>
            <option value="under_construction">Under Construction</option>
            <option value="completed">Completed</option>
            <option value="stalled">Stalled</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">NOC Status *</label>
          <select
            required
            value={formData.noc_status || 'pending'}
            onChange={(e) => setFormData({ ...formData, noc_status: e.target.value as any })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="pending">Pending</option>
            <option value="approved">Approved</option>
            <option value="rejected">Rejected</option>
          </select>
        </div>
      </div>

      <div className="flex gap-3">
        <button
          type="submit"
          disabled={loading}
          className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium disabled:opacity-50"
        >
          {loading ? 'Saving...' : project ? 'Update Project' : 'Create Project'}
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
