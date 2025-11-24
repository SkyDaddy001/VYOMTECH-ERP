'use client'

import React, { useEffect, useState } from 'react'
import { useCompanyStore } from '@/contexts/phase3cStore'
import toast from 'react-hot-toast'

export function CompanyDashboard() {
  const { companies, selectedCompany, selectCompany, fetchCompanies, loading, error } =
    useCompanyStore()
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    industry_type: '',
    employee_count: 0,
    website: '',
  })

  useEffect(() => {
    fetchCompanies()
  }, [fetchCompanies])

  useEffect(() => {
    if (error) {
      toast.error(error)
    }
  }, [error])

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: name === 'employee_count' ? parseInt(value) : value,
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await useCompanyStore.getState().createCompany(formData)
      toast.success('Company created successfully!')
      setFormData({
        name: '',
        description: '',
        industry_type: '',
        employee_count: 0,
        website: '',
      })
      setShowForm(false)
    } catch (error: any) {
      toast.error(error.message || 'Failed to create company')
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Companies</h2>
        <button
          onClick={() => setShowForm(!showForm)}
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg transition"
        >
          {showForm ? 'Cancel' : 'New Company'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <input
              type="text"
              name="name"
              placeholder="Company Name"
              value={formData.name}
              onChange={handleInputChange}
              required
              className="border rounded-lg px-3 py-2 w-full"
            />
            <input
              type="text"
              name="industry_type"
              placeholder="Industry Type"
              value={formData.industry_type}
              onChange={handleInputChange}
              className="border rounded-lg px-3 py-2 w-full"
            />
            <textarea
              name="description"
              placeholder="Description"
              value={formData.description}
              onChange={handleInputChange}
              className="border rounded-lg px-3 py-2 w-full col-span-1 md:col-span-2"
              rows={3}
            />
            <input
              type="number"
              name="employee_count"
              placeholder="Employee Count"
              value={formData.employee_count}
              onChange={handleInputChange}
              className="border rounded-lg px-3 py-2 w-full"
            />
            <input
              type="url"
              name="website"
              placeholder="Website"
              value={formData.website}
              onChange={handleInputChange}
              className="border rounded-lg px-3 py-2 w-full"
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="mt-4 bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-lg transition disabled:opacity-50"
          >
            {loading ? 'Creating...' : 'Create Company'}
          </button>
        </form>
      )}

      {loading ? (
        <div className="text-center py-8">Loading companies...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {companies.map((company) => (
            <div
              key={company.id}
              onClick={() => selectCompany(company)}
              className={`p-4 rounded-lg cursor-pointer transition border-2 ${
                selectedCompany?.id === company.id
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <h3 className="font-bold text-lg">{company.name}</h3>
              <p className="text-gray-600 text-sm">{company.industry_type}</p>
              <p className="text-gray-500 text-xs mt-2">Users: {company.current_user_count}</p>
              <p className="text-gray-500 text-xs">Projects: {company.current_project_count}</p>
              <span
                className={`inline-block mt-2 px-2 py-1 text-xs rounded ${
                  company.status === 'active'
                    ? 'bg-green-100 text-green-800'
                    : 'bg-gray-100 text-gray-800'
                }`}
              >
                {company.status}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
