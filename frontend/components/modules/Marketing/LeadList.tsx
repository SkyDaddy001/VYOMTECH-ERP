'use client'

import { Lead } from '@/types/marketing'

interface LeadListProps {
  leads: Lead[]
  loading: boolean
  onEdit: (lead: Lead) => void
  onDelete: (lead: Lead) => void
  onQualify: (lead: Lead) => void
  onConvert: (lead: Lead) => void
}

export default function LeadList({ leads, loading, onEdit, onDelete, onQualify, onConvert }: LeadListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading leads...</p>
        </div>
      </div>
    )
  }

  if (leads.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No leads yet.</p>
      </div>
    )
  }

  const newLeads = leads.filter((l) => l.lead_stage === 'inquiry').length
  const qualifiedLeads = leads.filter((l) => l.lead_stage === 'qualified').length
  const convertedLeads = leads.filter((l) => l.lead_stage === 'converted').length

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{leads.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">New</p>
          <p className="text-2xl font-bold text-blue-600 mt-1">{newLeads}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Qualified</p>
          <p className="text-2xl font-bold text-orange-600 mt-1">{qualifiedLeads}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Converted</p>
          <p className="text-2xl font-bold text-green-600 mt-1">{convertedLeads}</p>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-lg shadow overflow-x-auto border border-gray-200">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Phone</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Budget Range</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Source</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Lead Stage</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Quality</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {leads.map((lead) => (
              <tr key={lead.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{lead.prospect_name}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{lead.email}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{lead.phone}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                  ₹{lead.budget_range_min?.toLocaleString()} - ₹{lead.budget_range_max?.toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                  <span className="px-2 py-1 bg-cyan-100 text-cyan-800 rounded text-xs">{lead.source}</span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <span
                    className={`px-2 py-1 rounded text-xs font-medium ${
                      lead.lead_stage === 'inquiry'
                        ? 'bg-blue-100 text-blue-800'
                        : lead.lead_stage === 'qualified'
                        ? 'bg-orange-100 text-orange-800'
                        : lead.lead_stage === 'converted'
                        ? 'bg-green-100 text-green-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {lead.lead_stage?.replace('_', ' ')}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <span className={`px-2 py-1 rounded text-xs ${
                    lead.lead_quality === 'hot' ? 'bg-red-100 text-red-800' :
                    lead.lead_quality === 'warm' ? 'bg-orange-100 text-orange-800' :
                    'bg-blue-100 text-blue-800'
                  }`}>
                    {lead.lead_quality}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                  <button onClick={() => onEdit(lead)} className="text-blue-600 hover:text-blue-900 font-medium text-xs">
                    Edit
                  </button>
                  <button onClick={() => onDelete(lead)} className="text-red-600 hover:text-red-900 font-medium text-xs">
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
