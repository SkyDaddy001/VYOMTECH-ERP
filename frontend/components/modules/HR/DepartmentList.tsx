'use client'

import { Department } from '@/types/hr'

interface DepartmentListProps {
  departments: Department[]
  loading: boolean
  onEdit: (department: Department) => void
  onDelete: (department: Department) => void
}

export default function DepartmentList({ departments, loading, onEdit, onDelete }: DepartmentListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading departments...</p>
        </div>
      </div>
    )
  }

  if (departments.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No departments yet. Add your first department to get started.</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{departments.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total Budget</p>
          <p className="text-2xl font-bold text-blue-600 mt-1">
            ₹{departments.reduce((sum, d) => sum + (d.budget || 0), 0).toLocaleString()}
          </p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Active</p>
          <p className="text-2xl font-bold text-green-600 mt-1">
            {departments.filter((d) => d.status === 'active').length}
          </p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Avg Budget</p>
          <p className="text-2xl font-bold text-purple-600 mt-1">
            ₹{Math.round(departments.reduce((sum, d) => sum + (d.budget || 0), 0) / departments.length).toLocaleString()}
          </p>
        </div>
      </div>

      {/* Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {departments.map((dept) => (
          <div key={dept.id} className="bg-white rounded-lg shadow border border-gray-200 p-6 hover:shadow-lg transition">
            <div className="flex items-start justify-between mb-4">
              <div>
                <h3 className="text-lg font-semibold text-gray-900">{dept.name}</h3>
                <p className="text-sm text-gray-600">{dept.description}</p>
              </div>
              <span
                className={`px-3 py-1 rounded text-xs font-medium ${
                  dept.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                }`}
              >
                {dept.status}
              </span>
            </div>

            <div className="space-y-2 mb-4">
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Budget</span>
                <span className="font-semibold text-gray-900">₹{(dept.budget || 0).toLocaleString()}</span>
              </div>
              {dept.head_id && (
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Head ID</span>
                  <span className="font-semibold text-gray-900">{dept.head_id}</span>
                </div>
              )}
              {dept.employee_count !== undefined && (
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Employees</span>
                  <span className="font-semibold text-blue-600">{dept.employee_count}</span>
                </div>
              )}
            </div>

            <div className="flex gap-2 pt-4 border-t border-gray-200">
              <button
                onClick={() => onEdit(dept)}
                className="flex-1 text-green-600 hover:text-green-900 hover:bg-green-50 py-2 rounded font-medium text-sm transition"
              >
                Edit
              </button>
              <button
                onClick={() => onDelete(dept)}
                className="flex-1 text-red-600 hover:text-red-900 hover:bg-red-50 py-2 rounded font-medium text-sm transition"
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
