'use client'

import { useState, useEffect } from 'react'
import CompanyList from '@/components/modules/Company/CompanyList'
import CompanyForm from '@/components/modules/Company/CompanyForm'
import ProjectList from '@/components/modules/Company/ProjectList'
import ProjectForm from '@/components/modules/Company/ProjectForm'
import { Company, Project } from '@/types/company'
import { companyService } from '@/services/company.service'
import toast from 'react-hot-toast'

type View = 'companies' | 'projects'

export default function CompanyPage() {
  const [view, setView] = useState<View>('companies')
  const [companies, setCompanies] = useState<Company[]>([])
  const [projects, setProjects] = useState<Project[]>([])
  const [selectedCompany, setSelectedCompany] = useState<Company | undefined>()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showCompanyForm, setShowCompanyForm] = useState(false)
  const [editingCompany, setEditingCompany] = useState<Company | undefined>()
  const [showProjectForm, setShowProjectForm] = useState(false)
  const [editingProject, setEditingProject] = useState<Project | undefined>()

  useEffect(() => {
    loadCompanies()
  }, [])

  const loadCompanies = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await companyService.getCompanies()
      setCompanies(data)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to load companies'
      setError(errorMsg)
      toast.error(errorMsg)
      console.error('Error loading companies:', err)
    } finally {
      setLoading(false)
    }
  }

  const loadProjects = async (companyId: string) => {
    try {
      setLoading(true)
      const data = await companyService.getProjects(companyId)
      setProjects(data)
    } catch (err) {
      toast.error('Failed to load projects')
      console.error('Error loading projects:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateCompany = async (data: Partial<Company>) => {
    try {
      await companyService.createCompany(data)
      toast.success('Company created successfully!')
      await loadCompanies()
      setShowCompanyForm(false)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to create company'
      toast.error(errorMsg)
      throw err
    }
  }

  const handleEditCompany = async (data: Partial<Company>) => {
    try {
      if (!editingCompany) return
      await companyService.updateCompany(editingCompany.id, data)
      toast.success('Company updated successfully!')
      await loadCompanies()
      setShowCompanyForm(false)
      setEditingCompany(undefined)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to update company'
      toast.error(errorMsg)
      throw err
    }
  }

  const handleDeleteCompany = async (company: Company) => {
    try {
      if (confirm(`Are you sure you want to delete ${company.name}?`)) {
        await companyService.deleteCompany(company.id)
        toast.success('Company deleted successfully!')
        await loadCompanies()
      }
    } catch (err) {
      toast.error('Failed to delete company')
    }
  }

  const handleViewProjects = async (company: Company) => {
    setSelectedCompany(company)
    await loadProjects(company.id)
    setView('projects')
  }

  const handleCreateProject = async (data: Partial<Project>) => {
    try {
      if (!selectedCompany) return
      await companyService.createProject(selectedCompany.id, data)
      toast.success('Project created successfully!')
      await loadProjects(selectedCompany.id)
      setShowProjectForm(false)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to create project'
      toast.error(errorMsg)
      throw err
    }
  }

  const handleEditProject = async (data: Partial<Project>) => {
    try {
      if (!editingProject) return
      await companyService.updateProject(editingProject.id, data)
      toast.success('Project updated successfully!')
      if (selectedCompany) {
        await loadProjects(selectedCompany.id)
      }
      setShowProjectForm(false)
      setEditingProject(undefined)
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to update project'
      toast.error(errorMsg)
      throw err
    }
  }

  const handleDeleteProject = async (project: Project) => {
    try {
      if (confirm(`Are you sure you want to delete ${project.name}?`)) {
        await companyService.deleteProject(project.id)
        toast.success('Project deleted successfully!')
        if (selectedCompany) {
          await loadProjects(selectedCompany.id)
        }
      }
    } catch (err) {
      toast.error('Failed to delete project')
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            {view === 'companies' ? 'Companies' : selectedCompany?.name}
          </h1>
          <p className="mt-2 text-gray-600">
            {view === 'companies' ? 'Manage your organizations and companies' : 'Manage company projects'}
          </p>
        </div>
        {view === 'companies' && (
          <button
            onClick={() => {
              setEditingCompany(undefined)
              setShowCompanyForm(true)
            }}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700"
          >
            + New Company
          </button>
        )}
        {view === 'projects' && (
          <div className="space-x-3">
            <button
              onClick={() => setView('companies')}
              className="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50"
            >
              ← Back
            </button>
            <button
              onClick={() => {
                setEditingProject(undefined)
                setShowProjectForm(true)
              }}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700"
            >
              + New Project
            </button>
          </div>
        )}
      </div>

      {error && (
        <div className="rounded-md bg-red-50 p-4">
          <p className="text-sm font-medium text-red-800">{error}</p>
        </div>
      )}

      {view === 'companies' && (
        <CompanyList
          companies={companies}
          loading={loading}
          onEdit={(company) => {
            setEditingCompany(company)
            setShowCompanyForm(true)
          }}
          onDelete={handleDeleteCompany}
          onViewProjects={handleViewProjects}
        />
      )}

      {view === 'projects' && (
        <ProjectList
          projects={projects}
          loading={loading}
          onEdit={(project) => {
            setEditingProject(project)
            setShowProjectForm(true)
          }}
          onDelete={handleDeleteProject}
          onViewMembers={(project) => {
            toast.success('Project members feature coming soon!')
          }}
        />
      )}

      {showCompanyForm && (
        <CompanyForm
          company={editingCompany}
          onSubmit={editingCompany ? handleEditCompany : handleCreateCompany}
          onCancel={() => {
            setShowCompanyForm(false)
            setEditingCompany(undefined)
          }}
        />
      )}

      {showProjectForm && (
        <ProjectForm
          project={editingProject}
          onSubmit={editingProject ? handleEditProject : handleCreateProject}
          onCancel={() => {
            setShowProjectForm(false)
            setEditingProject(undefined)
          }}
        />
      )}
    </div>
  )
}
