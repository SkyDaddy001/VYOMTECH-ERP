'use client'

import { Project } from '@/types/company'

interface ProjectListProps {
  projects: Project[]
  loading: boolean
  onEdit: (project: Project) => void
  onDelete: (project: Project) => void
  onViewMembers: (project: Project) => void
}

export default function ProjectList({
  projects,
  loading,
  onEdit,
  onDelete,
  onViewMembers,
}: ProjectListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading projects...</p>
        </div>
      </div>
    )
  }

  if (projects.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No projects yet. Create your first project to get started.</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Summary Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="bg-white rounded-lg shadow p-6 border border-gray-200">
          <p className="text-gray-600 text-sm font-medium">Total Projects</p>
          <p className="text-3xl font-bold text-gray-900 mt-2">{projects.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-6 border border-gray-200">
          <p className="text-gray-600 text-sm font-medium">Active</p>
          <p className="text-3xl font-bold text-green-600 mt-2">
            {projects.filter((p) => p.status === 'active').length}
          </p>
        </div>
        <div className="bg-white rounded-lg shadow p-6 border border-gray-200">
          <p className="text-gray-600 text-sm font-medium">Total Budget</p>
          <p className="text-3xl font-bold text-blue-600 mt-2">
            ${projects.reduce((sum, p) => sum + p.budget_allocated, 0).toLocaleString()}
          </p>
        </div>
      </div>

      {/* Projects Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden border border-gray-200">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Project Name
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Type
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Budget
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Users
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {projects.map((project) => {
              const budgetPercentage = (project.budget_spent / project.budget_allocated) * 100
              return (
                <tr key={project.id} className="hover:bg-gray-50 transition">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <p className="text-sm font-medium text-gray-900">{project.name}</p>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded-full text-xs font-medium">
                      {project.project_type}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span
                      className={`px-2 py-1 rounded-full text-xs font-medium ${
                        project.status === 'active'
                          ? 'bg-green-100 text-green-800'
                          : project.status === 'archived'
                          ? 'bg-gray-100 text-gray-800'
                          : 'bg-yellow-100 text-yellow-800'
                      }`}
                    >
                      {project.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    <div className="flex items-center space-x-2">
                      <span>${project.budget_spent.toLocaleString()} / ${project.budget_allocated.toLocaleString()}</span>
                      <div className="w-16 bg-gray-200 rounded-full h-2">
                        <div
                          className={`h-2 rounded-full ${
                            budgetPercentage > 90 ? 'bg-red-500' : budgetPercentage > 70 ? 'bg-yellow-500' : 'bg-green-500'
                          }`}
                          style={{ width: `${Math.min(budgetPercentage, 100)}%` }}
                        />
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    {project.current_user_count} / {project.max_users}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                    <button
                      onClick={() => onViewMembers(project)}
                      className="text-blue-600 hover:text-blue-900 font-medium"
                    >
                      Members
                    </button>
                    <button
                      onClick={() => onEdit(project)}
                      className="text-green-600 hover:text-green-900 font-medium"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => onDelete(project)}
                      className="text-red-600 hover:text-red-900 font-medium"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}
