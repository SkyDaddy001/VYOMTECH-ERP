'use client'

import React, { useState } from 'react'
import { SnagList } from '@/types/postsales'

interface SnagListComponentProps {
  snags: SnagList[]
  loading: boolean
  onEdit: (snag: SnagList) => void
  onDelete: (snag: SnagList) => void
  onStatusChange: (snag: SnagList, status: string) => void
}

export default function SnagListComponent({ snags, loading, onEdit, onDelete, onStatusChange }: SnagListComponentProps) {
  const [filterStatus, setFilterStatus] = useState<string>('all')

  const filteredSnags = filterStatus === 'all'
    ? snags
    : snags.filter(s => s.status === filterStatus)

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      'open': 'bg-red-100 text-red-800',
      'in_progress': 'bg-yellow-100 text-yellow-800',
      'resolved': 'bg-green-100 text-green-800',
      'pending_inspection': 'bg-blue-100 text-blue-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const getSeverityBadge = (severity: string) => {
    const colors: Record<string, string> = {
      'low': 'bg-green-100 text-green-800',
      'medium': 'bg-yellow-100 text-yellow-800',
      'high': 'bg-red-100 text-red-800',
    }
    return colors[severity] || 'bg-gray-100 text-gray-800'
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-600">Loading snags...</div>
  }

  return (
    <div className="space-y-4">
      <select
        value={filterStatus}
        onChange={(e) => setFilterStatus(e.target.value)}
        className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="all">All Snags</option>
        <option value="open">Open</option>
        <option value="in_progress">In Progress</option>
        <option value="resolved">Resolved</option>
        <option value="pending_inspection">Pending Inspection</option>
      </select>

      <div className="space-y-3">
        {filteredSnags.map((snag) => (
          <div key={snag.id} className="bg-white rounded-lg shadow p-4 border border-gray-200">
            <div className="flex justify-between items-start mb-2">
              <div className="flex-1">
                <p className="font-semibold text-gray-900">{snag.snag_description}</p>
                <p className="text-xs text-gray-600">{snag.snag_category.replace(/_/g, ' ')}</p>
              </div>
              <div className="flex gap-2">
                <span className={`px-2 py-1 text-xs font-medium rounded ${getSeverityBadge(snag.severity)}`}>
                  {snag.severity}
                </span>
                <span className={`px-2 py-1 text-xs font-medium rounded ${getStatusBadge(snag.status)}`}>
                  {snag.status.replace('_', ' ')}
                </span>
              </div>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-3 gap-2 mb-3 text-xs">
              <div>
                <p className="text-gray-600">Reported</p>
                <p className="font-medium">{new Date(snag.reported_date).toLocaleDateString()}</p>
              </div>
              <div>
                <p className="text-gray-600">Target</p>
                <p className="font-medium text-orange-600">{new Date(snag.target_completion_date).toLocaleDateString()}</p>
              </div>
              {snag.actual_completion_date && (
                <div>
                  <p className="text-gray-600">Completed</p>
                  <p className="font-medium text-green-600">{new Date(snag.actual_completion_date).toLocaleDateString()}</p>
                </div>
              )}
            </div>

            <div className="flex gap-2">
              <select
                value={snag.status}
                onChange={(e) => onStatusChange(snag, e.target.value)}
                className="flex-1 px-3 py-2 text-xs font-medium border border-gray-300 rounded hover:bg-gray-50 transition"
              >
                <option value="open">Open</option>
                <option value="in_progress">In Progress</option>
                <option value="resolved">Resolved</option>
                <option value="pending_inspection">Pending Inspection</option>
              </select>
              <button
                onClick={() => onEdit(snag)}
                className="flex-1 px-3 py-2 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition"
              >
                Edit
              </button>
              <button
                onClick={() => {
                  if (confirm('Delete this snag?')) {
                    onDelete(snag)
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

      {filteredSnags.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No snags found</p>
        </div>
      )}
    </div>
  )
}
