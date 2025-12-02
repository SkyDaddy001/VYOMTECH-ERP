'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import ProjectList from '@/components/modules/Projects/ProjectList'
import ProjectForm from '@/components/modules/Projects/ProjectForm'
import MilestoneList from '@/components/modules/Projects/MilestoneList'
import MilestoneForm from '@/components/modules/Projects/MilestoneForm'
import { projectsService } from '@/services/projects.service'
import { Project, ProjectMilestone, ProjectMetrics } from '@/types/projects'

type TabType = 'projects' | 'milestones'
type FormType = 'project' | 'milestone' | null

export default function ProjectsPage() {
  const [activeTab, setActiveTab] = useState<TabType>('projects')
  const [showForm, setShowForm] = useState(false)
  const [editingItem, setEditingItem] = useState<any>(null)
  const [formType, setFormType] = useState<FormType>(null)

  // Data states
  const [projects, setProjects] = useState<Project[]>([])
  const [milestones, setMilestones] = useState<ProjectMilestone[]>([])
  const [metrics, setMetrics] = useState<ProjectMetrics | null>(null)
  const [selectedProject, setSelectedProject] = useState<Project | null>(null)

  // Loading states
  const [projectsLoading, setProjectsLoading] = useState(false)
  const [milestonesLoading, setMilestonesLoading] = useState(false)

  // Load projects
  const loadProjects = async () => {
    setProjectsLoading(true)
    try {
      const data = await projectsService.getProjects()
      setProjects(data)
      if (data.length > 0 && !selectedProject) {
        setSelectedProject(data[0])
      }
    } catch (error) {
      toast.error('Failed to load projects')
    } finally {
      setProjectsLoading(false)
    }
  }

  // Load milestones for selected project
  const loadMilestones = async () => {
    if (!selectedProject?.id) return
    setMilestonesLoading(true)
    try {
      const data = await projectsService.getMilestones(selectedProject.id)
      setMilestones(data)
    } catch (error) {
      toast.error('Failed to load milestones')
    } finally {
      setMilestonesLoading(false)
    }
  }

  // Load metrics
  const loadMetrics = async () => {
    try {
      const data = await projectsService.getMetrics()
      setMetrics(data)
    } catch (error) {
      // Metrics are optional
    }
  }

  // Load data on mount and tab change
  useEffect(() => {
    loadMetrics()
    if (activeTab === 'projects') {
      loadProjects()
    } else if (activeTab === 'milestones') {
      loadMilestones()
    }
  }, [activeTab])

  // Project CRUD
  const handleCreateProject = () => {
    setEditingItem(null)
    setFormType('project')
    setShowForm(true)
  }

  const handleEditProject = (project: Project) => {
    setEditingItem(project)
    setFormType('project')
    setShowForm(true)
  }

  const handleDeleteProject = async (project: Project) => {
    if (!confirm('Are you sure?')) return
    try {
      await projectsService.deleteProject(project.id || '')
      toast.success('Project deleted!')
      loadProjects()
    } catch (error) {
      toast.error('Failed to delete project')
    }
  }

  const handleSubmitProject = async (data: Partial<Project>) => {
    try {
      if (editingItem) {
        await projectsService.updateProject(editingItem.id, data)
      } else {
        await projectsService.createProject(data)
      }
      setShowForm(false)
      loadProjects()
    } catch (error) {
      throw error
    }
  }

  // Milestone CRUD
  const handleCreateMilestone = () => {
    setEditingItem(null)
    setFormType('milestone')
    setShowForm(true)
  }

  const handleEditMilestone = (milestone: ProjectMilestone) => {
    setEditingItem(milestone)
    setFormType('milestone')
    setShowForm(true)
  }

  const handleDeleteMilestone = async (milestone: ProjectMilestone) => {
    if (!confirm('Are you sure?')) return
    try {
      await projectsService.deleteMilestone(selectedProject?.id || '', milestone.id || '')
      toast.success('Milestone deleted!')
      loadMilestones()
    } catch (error) {
      toast.error('Failed to delete milestone')
    }
  }

  const handleSubmitMilestone = async (data: Partial<ProjectMilestone>) => {
    try {
      if (editingItem) {
        await projectsService.updateMilestone(selectedProject?.id || '', editingItem.id, data)
      } else {
        await projectsService.createMilestone(selectedProject?.id || '', data)
      }
      setShowForm(false)
      loadMilestones()
    } catch (error) {
      throw error
    }
  }

  const closeForm = () => {
    setShowForm(false)
    setEditingItem(null)
    setFormType(null)
  }

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">Real Estate Projects</h1>
          <p className="text-gray-600">Manage construction projects, milestones, and timelines</p>
        </div>

        {/* Metrics */}
        {metrics && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Total Projects</p>
              <p className="text-2xl font-bold text-blue-600 mt-1">{metrics.total_projects}</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Active</p>
              <p className="text-2xl font-bold text-green-600 mt-1">{metrics.active_projects}</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Completed</p>
              <p className="text-2xl font-bold text-purple-600 mt-1">{metrics.completed_projects}</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Delayed</p>
              <p className="text-2xl font-bold text-red-600 mt-1">{metrics.delayed_projects}</p>
            </div>
          </div>
        )}

        {/* Tabs */}
        <div className="bg-white rounded-lg shadow mb-6 border border-gray-200">
          <div className="flex flex-wrap border-b border-gray-200">
            <button
              onClick={() => {
                setActiveTab('projects')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'projects'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Projects
            </button>
            <button
              onClick={() => {
                setActiveTab('milestones')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'milestones'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Milestones
            </button>
          </div>

          {/* Tab Content */}
          <div className="p-6">
            {/* Projects Tab */}
            {activeTab === 'projects' && (
              <div>
                {!showForm ? (
                  <div>
                    <button
                      onClick={handleCreateProject}
                      className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                    >
                      + New Project
                    </button>
                    <ProjectList
                      projects={projects}
                      loading={projectsLoading}
                      onEdit={handleEditProject}
                      onDelete={handleDeleteProject}
                    />
                  </div>
                ) : (
                  <ProjectForm
                    project={editingItem}
                    onSubmit={handleSubmitProject}
                    onCancel={closeForm}
                  />
                )}
              </div>
            )}

            {/* Milestones Tab */}
            {activeTab === 'milestones' && (
              <div>
                {selectedProject ? (
                  <>
                    <div className="mb-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
                      <p className="text-sm text-gray-600">Selected Project: <span className="font-semibold text-blue-700">{selectedProject.project_name}</span></p>
                    </div>
                    {!showForm ? (
                      <div>
                        <button
                          onClick={handleCreateMilestone}
                          className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                        >
                          + New Milestone
                        </button>
                        <MilestoneList
                          milestones={milestones}
                          loading={milestonesLoading}
                          onEdit={handleEditMilestone}
                          onDelete={handleDeleteMilestone}
                        />
                      </div>
                    ) : (
                      <MilestoneForm
                        milestone={editingItem}
                        onSubmit={handleSubmitMilestone}
                        onCancel={closeForm}
                      />
                    )}
                  </>
                ) : (
                  <div className="text-center py-12 text-gray-500">
                    <p>Please create or select a project first</p>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
