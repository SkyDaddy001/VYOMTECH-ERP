'use client'

import { Employee } from '@/types/hr'

interface EmployeeListProps {
  employees: Employee[]
  loading: boolean
  onEdit: (employee: Employee) => void
  onDelete: (employee: Employee) => void
}

export default function EmployeeList({ employees, loading, onEdit, onDelete }: EmployeeListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading employees...</p>
        </div>
      </div>
    )
  }

  if (employees.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No employees yet. Add your first employee to get started.</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{employees.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Active</p>
          <p className="text-2xl font-bold text-green-600 mt-1">
            {employees.filter((e) => e.status === 'active').length}
          </p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">On Leave</p>
          <p className="text-2xl font-bold text-yellow-600 mt-1">
            {employees.filter((e) => e.status === 'on_leave').length}
          </p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Full Time</p>
          <p className="text-2xl font-bold text-blue-600 mt-1">
            {employees.filter((e) => e.employment_type === 'full_time').length}
          </p>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-lg shadow overflow-x-auto border border-gray-200">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Designation</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {employees.map((emp) => (
              <tr key={emp.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {emp.first_name} {emp.last_name}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{emp.email}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{emp.designation}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                  <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs">{emp.employment_type}</span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <span
                    className={`px-2 py-1 rounded text-xs font-medium ${
                      emp.status === 'active'
                        ? 'bg-green-100 text-green-800'
                        : emp.status === 'on_leave'
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {emp.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                  <button onClick={() => onEdit(emp)} className="text-green-600 hover:text-green-900 font-medium">
                    Edit
                  </button>
                  <button onClick={() => onDelete(emp)} className="text-red-600 hover:text-red-900 font-medium">
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
