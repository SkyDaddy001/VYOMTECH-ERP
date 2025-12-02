'use client'

import { useState } from 'react'
import { Project, ProjectMilestone } from '@/types/projects'

interface ProjectTimelineProps {
  project?: Project | null
  milestones?: ProjectMilestone[]
  onUpdateMilestone?: (milestone: ProjectMilestone) => Promise<void>
}

export default function ProjectTimeline({ project, milestones = [], onUpdateMilestone }: ProjectTimelineProps) {
  const [selectedMilestone, setSelectedMilestone] = useState<ProjectMilestone | null>(null)
  const [editingStatus, setEditingStatus] = useState<Record<string, string>>({})

  const calculateProjectProgress = () => {
    if (milestones.length === 0) return 0
    const totalProgress = milestones.reduce((sum, m) => sum + (m.completion_percentage || 0), 0)
    return Math.round(totalProgress / milestones.length)
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800 border-green-300'
      case 'in_progress':
        return 'bg-blue-100 text-blue-800 border-blue-300'
      case 'pending':
        return 'bg-yellow-100 text-yellow-800 border-yellow-300'
      case 'delayed':
        return 'bg-red-100 text-red-800 border-red-300'
      default:
        return 'bg-gray-100 text-gray-800 border-gray-300'
    }
  }

  const getProgressColor = (progress: number) => {
    if (progress === 100) return 'bg-green-500'
    if (progress >= 75) return 'bg-blue-500'
    if (progress >= 50) return 'bg-yellow-500'
    return 'bg-orange-500'
  }

  const calculateDaysRemaining = (endDate: string) => {
    const end = new Date(endDate)
    const today = new Date()
    const daysRemaining = Math.ceil((end.getTime() - today.getTime()) / (1000 * 60 * 60 * 24))
    return daysRemaining
  }

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('en-IN', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    })
  }

  const projectProgress = calculateProjectProgress()

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      {/* Project Overview */}
      <div>
        <h2 className="text-2xl font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <span className="text-3xl">🏗️</span> Project Timeline
        </h2>
      </div>

      {/* Overall Progress */}
      {project && (
        <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg p-6 border border-blue-200">
          <div className="grid grid-cols-4 gap-4 mb-4">
            <div>
              <p className="text-sm text-gray-600 mb-1">Project Name</p>
              <p className="text-lg font-semibold text-gray-900">{project.project_name}</p>
            </div>
            <div>
              <p className="text-sm text-gray-600 mb-1">Status</p>
              <span className={`inline-block px-3 py-1 rounded-full text-sm font-medium border ${getStatusColor(project.status || 'on_track')}`}>
                {(project.status || 'on_track').replace(/_/g, ' ')}
              </span>
            </div>
            <div>
              <p className="text-sm text-gray-600 mb-1">Start Date</p>
              <p className="font-semibold text-gray-900">{formatDate(project.launch_date)}</p>
            </div>
            <div>
              <p className="text-sm text-gray-600 mb-1">End Date</p>
              <p className="font-semibold text-gray-900">{formatDate(project.expected_completion)}</p>
            </div>
          </div>

          {/* Overall Progress Bar */}
          <div>
            <div className="flex justify-between items-center mb-2">
              <p className="text-sm font-medium text-gray-700">Overall Progress</p>
              <p className="text-lg font-bold text-blue-600">{projectProgress}%</p>
            </div>
            <div className="w-full bg-gray-300 rounded-full h-3">
              <div
                className={`h-3 rounded-full ${getProgressColor(projectProgress)} transition-all duration-500`}
                style={{ width: `${projectProgress}%` }}
              />
            </div>
          </div>
        </div>
      )}

      {/* Milestones Timeline */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-gray-800">Milestones</h3>

        {milestones.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <p>No milestones added yet</p>
          </div>
        ) : (
          <div className="relative">
            {/* Timeline Line */}
            <div className="absolute left-8 top-0 bottom-0 w-1 bg-gradient-to-b from-blue-400 via-purple-400 to-pink-400" />

            {/* Milestones */}
            <div className="space-y-4">
              {milestones.map((milestone, index) => {
                return (
                  <div key={milestone.id || index} className="relative pl-24">
                    {/* Timeline Dot */}
                    <div className="absolute left-2 top-2 w-12 h-12 bg-white border-4 border-blue-500 rounded-full flex items-center justify-center font-bold text-blue-600">
                      {index + 1}
                    </div>

                    {/* Card */}
                    <div className="bg-white border-2 border-gray-200 rounded-lg p-4 hover:shadow-lg transition-shadow cursor-pointer"
                      onClick={() => setSelectedMilestone(milestone)}>
                      <div className="flex justify-between items-start mb-3">
                        <div className="flex-1">
                          <h4 className="text-lg font-semibold text-gray-900">{milestone.milestone_name}</h4>
                          <p className="text-sm text-gray-600 mt-1">{milestone.description || milestone.milestone_type}</p>
                        </div>
                        <span className={`inline-block px-3 py-1 rounded-full text-sm font-medium border ${getStatusColor(milestone.status)}`}>
                          {milestone.status.replace(/_/g, ' ')}
                        </span>
                      </div>

                      {/* Milestone Details Grid */}
                      <div className="grid grid-cols-3 gap-3 mb-3 text-sm">
                        <div>
                          <p className="text-gray-600 text-xs">Planned Date</p>
                          <p className="font-semibold text-gray-900">{formatDate(milestone.planned_date)}</p>
                        </div>
                        <div>
                          <p className="text-gray-600 text-xs">Type</p>
                          <p className="font-semibold text-gray-900">{milestone.milestone_type}</p>
                        </div>
                        <div>
                          <p className="text-gray-600 text-xs">Completion</p>
                          <p className="font-semibold text-gray-900">{milestone.completion_percentage}%</p>
                        </div>
                      </div>

                      {/* Progress Bar */}
                      <div>
                        <div className="flex justify-between items-center mb-1">
                          <p className="text-xs font-medium text-gray-700">Progress</p>
                          <p className="text-xs font-bold text-blue-600">{milestone.completion_percentage}%</p>
                        </div>
                        <div className="w-full bg-gray-300 rounded-full h-2">
                          <div
                            className={`h-2 rounded-full ${getProgressColor(milestone.completion_percentage)} transition-all duration-300`}
                            style={{ width: `${milestone.completion_percentage}%` }}
                          />
                        </div>
                      </div>

                      {/* Actual Date */}
                      {milestone.actual_date && (
                        <div className="mt-3 pt-3 border-t border-gray-200">
                          <p className="text-xs font-semibold text-gray-700">Actual Date:</p>
                          <p className="text-sm text-gray-900">{formatDate(milestone.actual_date)}</p>
                        </div>
                      )}
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        )}
      </div>

      {/* Legend */}
      <div className="bg-gray-50 rounded-lg p-4 border border-gray-200">
        <p className="text-sm font-semibold text-gray-700 mb-3">Status Legend</p>
        <div className="grid grid-cols-5 gap-3 text-xs">
          {[
            { status: 'Completed', color: 'bg-green-100 text-green-800' },
            { status: 'In Progress', color: 'bg-blue-100 text-blue-800' },
            { status: 'On Track', color: 'bg-yellow-100 text-yellow-800' },
            { status: 'At Risk', color: 'bg-orange-100 text-orange-800' },
            { status: 'Delayed', color: 'bg-red-100 text-red-800' }
          ].map(({ status, color }) => (
            <div key={status} className="flex items-center gap-2">
              <div className={`w-3 h-3 rounded-full ${color}`} />
              <span className="text-gray-700">{status}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
