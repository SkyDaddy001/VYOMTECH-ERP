'use client'

import React, { useState } from 'react'
import { ChangeRequest } from '@/types/postsales'

interface ChangeRequestListProps {
  changeRequests: ChangeRequest[]
  loading: boolean
  onEdit: (cr: ChangeRequest) => void
  onDelete: (cr: ChangeRequest) => void
}

export default function ChangeRequestList({ changeRequests, loading, onEdit, onDelete }: ChangeRequestListProps) {
  const [filterStatus, setFilterStatus] = useState<string>('all')

  const filteredCRs = filterStatus === 'all'
    ? changeRequests
    : changeRequests.filter(cr => cr.status === filterStatus)

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      'submitted': 'bg-blue-100 text-blue-800',
      'under_review': 'bg-yellow-100 text-yellow-800',
      'approved': 'bg-green-100 text-green-800',
      'rejected': 'bg-red-100 text-red-800',
      'implemented': 'bg-purple-100 text-purple-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const getImpactBadge = (impact: string) => {
    const colors: Record<string, string> = {
      'no_cost_change': 'bg-green-100 text-green-800',
      'additional_cost': 'bg-red-100 text-red-800',
      'cost_reduction': 'bg-blue-100 text-blue-800',
    }
    return colors[impact] || 'bg-gray-100 text-gray-800'
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-600">Loading change requests...</div>
  }

  return (
    <div className="space-y-4">
      <select
        value={filterStatus}
        onChange={(e) => setFilterStatus(e.target.value)}
        className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="all">All Status</option>
        <option value="submitted">Submitted</option>
        <option value="under_review">Under Review</option>
        <option value="approved">Approved</option>
        <option value="rejected">Rejected</option>
        <option value="implemented">Implemented</option>
      </select>

      <div className="space-y-3">
        {filteredCRs.map((cr) => (
          <div key={cr.id} className="bg-white rounded-lg shadow p-4 border border-gray-200">
            <div className="flex justify-between items-start mb-2">
              <div className="flex-1">
                <p className="font-semibold text-gray-900">{cr.crm_number}</p>
                <p className="text-xs text-gray-600">{cr.request_type.replace(/_/g, ' ')}</p>
              </div>
              <div className="flex gap-2">
                <span className={`px-2 py-1 text-xs font-medium rounded ${getImpactBadge(cr.impact)}`}>
                  {cr.impact.replace(/_/g, ' ')}
                </span>
                <span className={`px-2 py-1 text-xs font-medium rounded ${getStatusBadge(cr.status)}`}>
                  {cr.status.replace('_', ' ')}
                </span>
              </div>
            </div>

            <p className="text-sm text-gray-700 mb-3">{cr.description}</p>

            <div className="grid grid-cols-2 md:grid-cols-3 gap-2 mb-3 text-xs">
              <div>
                <p className="text-gray-600">Cost Impact</p>
                <p className={`font-medium ${cr.cost_difference > 0 ? 'text-red-600' : cr.cost_difference < 0 ? 'text-green-600' : 'text-gray-900'}`}>
                  ₹{(cr.cost_difference / 100000).toFixed(2)}L
                </p>
              </div>
              <div>
                <p className="text-gray-600">Requested</p>
                <p className="font-medium">{new Date(cr.request_date).toLocaleDateString()}</p>
              </div>
              {cr.approval_date && (
                <div>
                  <p className="text-gray-600">Approved</p>
                  <p className="font-medium">{new Date(cr.approval_date).toLocaleDateString()}</p>
                </div>
              )}
            </div>

            <div className="flex gap-2">
              <button
                onClick={() => onEdit(cr)}
                className="flex-1 px-3 py-2 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition"
              >
                Edit
              </button>
              <button
                onClick={() => {
                  if (confirm('Delete this CRM?')) {
                    onDelete(cr)
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

      {filteredCRs.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No change requests found</p>
        </div>
      )}
    </div>
  )
}
