'use client'

import { Company } from '@/types/company'

interface CompanyListProps {
  companies: Company[]
  loading: boolean
  onEdit: (company: Company) => void
  onDelete: (company: Company) => void
  onViewProjects: (company: Company) => void
}

export default function CompanyList({
  companies,
  loading,
  onEdit,
  onDelete,
  onViewProjects,
}: CompanyListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading companies...</p>
        </div>
      </div>
    )
  }

  if (companies.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No companies yet. Create your first company to get started.</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Summary Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-6 border border-gray-200">
          <p className="text-gray-600 text-sm font-medium">Total Companies</p>
          <p className="text-3xl font-bold text-gray-900 mt-2">{companies.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-6 border border-gray-200">
          <p className="text-gray-600 text-sm font-medium">Active</p>
          <p className="text-3xl font-bold text-green-600 mt-2">
            {companies.filter((c) => c.status === 'active').length}
          </p>
        </div>
        <div className="bg-white rounded-lg shadow p-6 border border-gray-200">
          <p className="text-gray-600 text-sm font-medium">Total Projects</p>
          <p className="text-3xl font-bold text-blue-600 mt-2">
            {companies.reduce((sum, c) => sum + c.current_project_count, 0)}
          </p>
        </div>
        <div className="bg-white rounded-lg shadow p-6 border border-gray-200">
          <p className="text-gray-600 text-sm font-medium">Total Users</p>
          <p className="text-3xl font-bold text-purple-600 mt-2">
            {companies.reduce((sum, c) => sum + c.current_user_count, 0)}
          </p>
        </div>
      </div>

      {/* Companies Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {companies.map((company) => (
          <div key={company.id} className="bg-white rounded-lg shadow-md hover:shadow-lg transition border border-gray-200 p-6">
            <div className="mb-4">
              <h3 className="text-lg font-bold text-gray-900">{company.name}</h3>
              {company.description && <p className="text-sm text-gray-600 mt-1">{company.description}</p>}
            </div>

            {/* Status Badge */}
            <div className="mb-4">
              <span
                className={`inline-block text-xs font-semibold px-3 py-1 rounded-full ${
                  company.status === 'active'
                    ? 'bg-green-100 text-green-800'
                    : company.status === 'inactive'
                    ? 'bg-gray-100 text-gray-800'
                    : 'bg-red-100 text-red-800'
                }`}
              >
                {company.status.charAt(0).toUpperCase() + company.status.slice(1)}
              </span>
            </div>

            {/* Stats */}
            <div className="space-y-2 mb-6 bg-gray-50 p-4 rounded-lg">
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Industry:</span>
                <span className="font-medium text-gray-800">{company.industry_type}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Projects:</span>
                <span className="font-medium text-gray-800">{company.current_project_count}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Users:</span>
                <span className="font-medium text-gray-800">{company.current_user_count}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Max Users:</span>
                <span className="font-medium text-gray-800">{company.max_users}</span>
              </div>
            </div>

            {/* Actions */}
            <div className="space-y-2">
              <button
                onClick={() => onViewProjects(company)}
                className="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition text-sm"
              >
                View Projects
              </button>
              <div className="flex gap-2">
                <button
                  onClick={() => onEdit(company)}
                  className="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded-lg transition text-sm"
                >
                  Edit
                </button>
                <button
                  onClick={() => onDelete(company)}
                  className="flex-1 bg-red-100 hover:bg-red-200 text-red-800 font-bold py-2 px-4 rounded-lg transition text-sm"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
