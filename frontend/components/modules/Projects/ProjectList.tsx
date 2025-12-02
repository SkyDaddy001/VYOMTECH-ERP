'use client'

import React, { useState } from 'react'
import toast from 'react-hot-toast'
import { Project } from '@/types/projects'

interface ProjectListProps {
  projects: Project[]
  loading: boolean
  onEdit: (project: Project) => void
  onDelete: (project: Project) => void
}

export default function ProjectList({ projects, loading, onEdit, onDelete }: ProjectListProps) {
  const [filterStatus, setFilterStatus] = useState<string>('all')

  const filteredProjects = filterStatus === 'all' 
    ? projects 
    : projects.filter(p => p.status === filterStatus)

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      'planning': 'bg-blue-100 text-blue-800',
      'approved': 'bg-green-100 text-green-800',
      'under_construction': 'bg-yellow-100 text-yellow-800',
      'completed': 'bg-purple-100 text-purple-800',
      'stalled': 'bg-red-100 text-red-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-600">Loading projects...</div>
  }

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <select
          value={filterStatus}
          onChange={(e) => setFilterStatus(e.target.value)}
          className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="all">All Status</option>
          <option value="planning">Planning</option>
          <option value="approved">Approved</option>
          <option value="under_construction">Under Construction</option>
          <option value="completed">Completed</option>
          <option value="stalled">Stalled</option>
        </select>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {filteredProjects.map((project) => (
          <div key={project.id} className="bg-white rounded-lg shadow p-4 border border-gray-200 hover:shadow-lg transition">
            <div className="flex justify-between items-start mb-3">
              <div>
                <h3 className="font-semibold text-gray-900">{project.project_name}</h3>
                <p className="text-xs text-gray-600">{project.project_code}</p>
              </div>
              <span className={`px-2 py-1 text-xs font-medium rounded ${getStatusBadge(project.status)}`}>
                {project.status.replace('_', ' ')}
              </span>
            </div>

            <div className="space-y-2 mb-4 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">Type:</span>
                <span className="font-medium">{project.project_type}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Units:</span>
                <span className="font-medium">{project.total_units}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Area:</span>
                <span className="font-medium">{(project.total_area / 1000).toFixed(1)}K sqft</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Location:</span>
                <span className="font-medium text-xs">{project.city}, {project.state}</span>
              </div>
            </div>

            <div className="flex gap-2">
              <button
                onClick={() => onEdit(project)}
                className="flex-1 px-3 py-2 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition"
              >
                Edit
              </button>
              <button
                onClick={() => {
                  if (confirm('Delete this project?')) {
                    onDelete(project)
                  }
                }}
                className="flex-1 px-3 py-2 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition"
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>

      {filteredProjects.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No projects found</p>
        </div>
      )}
    </div>
  )
}
