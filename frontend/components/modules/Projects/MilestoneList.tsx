'use client'

import React, { useState } from 'react'
import toast from 'react-hot-toast'
import { ProjectMilestone } from '@/types/projects'

interface MilestoneListProps {
  milestones: ProjectMilestone[]
  loading: boolean
  onEdit: (milestone: ProjectMilestone) => void
  onDelete: (milestone: ProjectMilestone) => void
}

export default function MilestoneList({ milestones, loading, onEdit, onDelete }: MilestoneListProps) {
  if (loading) {
    return <div className="text-center py-8 text-gray-600">Loading milestones...</div>
  }

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      'pending': 'bg-gray-100 text-gray-800',
      'in_progress': 'bg-blue-100 text-blue-800',
      'completed': 'bg-green-100 text-green-800',
      'delayed': 'bg-red-100 text-red-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const getTypeColor = (type: string) => {
    const colors: Record<string, string> = {
      'foundation': 'bg-orange-100 text-orange-800',
      'structure': 'bg-blue-100 text-blue-800',
      'finishing': 'bg-purple-100 text-purple-800',
      'handover': 'bg-green-100 text-green-800',
    }
    return colors[type] || 'bg-gray-100 text-gray-800'
  }

  return (
    <div className="space-y-4">
      {milestones.map((milestone) => (
        <div key={milestone.id} className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <div className="flex justify-between items-start mb-3">
            <div>
              <h3 className="font-semibold text-gray-900">{milestone.milestone_name}</h3>
              <p className="text-xs text-gray-600">Planned: {new Date(milestone.planned_date).toLocaleDateString()}</p>
            </div>
            <div className="flex gap-2">
              <span className={`px-2 py-1 text-xs font-medium rounded ${getTypeColor(milestone.milestone_type)}`}>
                {milestone.milestone_type}
              </span>
              <span className={`px-2 py-1 text-xs font-medium rounded ${getStatusBadge(milestone.status)}`}>
                {milestone.status.replace('_', ' ')}
              </span>
            </div>
          </div>

          <div className="mb-3">
            <div className="flex justify-between items-center mb-1">
              <span className="text-xs text-gray-600">Completion</span>
              <span className="text-xs font-medium">{milestone.completion_percentage}%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div
                className="bg-blue-600 h-2 rounded-full transition-all"
                style={{ width: `${milestone.completion_percentage}%` }}
              />
            </div>
          </div>

          {milestone.description && (
            <p className="text-sm text-gray-600 mb-3">{milestone.description}</p>
          )}

          <div className="flex gap-2">
            <button
              onClick={() => onEdit(milestone)}
              className="flex-1 px-3 py-2 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition"
            >
              Edit
            </button>
            <button
              onClick={() => {
                if (confirm('Delete this milestone?')) {
                  onDelete(milestone)
                }
              }}
              className="flex-1 px-3 py-2 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition"
            >
              Delete
            </button>
          </div>
        </div>
      ))}

      {milestones.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No milestones found</p>
        </div>
      )}
    </div>
  )
}
