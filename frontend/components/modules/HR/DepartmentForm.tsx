'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import { Department } from '@/types/hr'

interface DepartmentFormProps {
  department?: Department | null
  employees: Array<{ id?: string; first_name: string; last_name: string }>
  onSubmit: (data: Partial<Department>) => Promise<void>
  onCancel: () => void
}

export default function DepartmentForm({ department, employees, onSubmit, onCancel }: DepartmentFormProps) {
  const [formData, setFormData] = useState<Partial<Department>>({
    name: '',
    description: '',
    budget: 0,
    head_id: '',
    status: 'active',
  })
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (department) {
      setFormData({
        id: department.id,
        name: department.name,
        description: department.description,
        budget: department.budget,
        head_id: department.head_id,
        status: department.status,
      })
    }
  }, [department])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: name === 'budget' ? (value ? parseInt(value) : 0) : value,
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!formData.name) {
      toast.error('Please fill all required fields')
      return
    }

    setLoading(true)
    try {
      await onSubmit(formData)
      toast.success(department ? 'Department updated!' : 'Department created!')
      onCancel()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error saving department')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Basic Information */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Department Information</h3>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Department Name *</label>
          <input
            type="text"
            name="name"
            value={formData.name || ''}
            onChange={handleChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
      </div>

      {/* Description */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
        <textarea
          name="description"
          value={formData.description || ''}
          onChange={handleChange}
          rows={3}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      {/* Budget & Head */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Management</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Budget</label>
            <input
              type="number"
              name="budget"
              value={formData.budget || 0}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Department Head</label>
            <select
              name="head_id"
              value={formData.head_id || ''}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Select Department Head</option>
              {employees.map((emp) => (
                <option key={emp.id} value={emp.id}>
                  {emp.first_name} {emp.last_name}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* Status */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Status</h3>
        <select
          name="status"
          value={formData.status || 'active'}
          onChange={handleChange}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
        </select>
      </div>

      {/* Actions */}
      <div className="flex gap-4 pt-4 border-t border-gray-200">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
        >
          {loading ? 'Saving...' : department ? 'Update Department' : 'Create Department'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 bg-gray-200 text-gray-900 py-2 rounded-lg hover:bg-gray-300 font-medium"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
