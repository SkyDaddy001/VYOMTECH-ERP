'use client'

import React, { useState } from 'react'
import { CustomerInteraction } from '@/types/postsales'

interface InteractionListProps {
  interactions: CustomerInteraction[]
  loading: boolean
  onEdit: (interaction: CustomerInteraction) => void
  onDelete: (interaction: CustomerInteraction) => void
  onUpdateStatus: (interaction: CustomerInteraction, status: string) => void
}

export default function InteractionList({ interactions, loading, onEdit, onDelete, onUpdateStatus }: InteractionListProps) {
  const [filterStatus, setFilterStatus] = useState<string>('all')

  const filteredInteractions = filterStatus === 'all'
    ? interactions
    : interactions.filter(i => i.status === filterStatus)

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      'pending': 'bg-yellow-100 text-yellow-800',
      'resolved': 'bg-green-100 text-green-800',
      'escalated': 'bg-red-100 text-red-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const getPriorityBadge = (priority: string) => {
    const colors: Record<string, string> = {
      'low': 'bg-blue-100 text-blue-800',
      'medium': 'bg-orange-100 text-orange-800',
      'high': 'bg-red-100 text-red-800',
      'critical': 'bg-purple-100 text-purple-800',
    }
    return colors[priority] || 'bg-gray-100 text-gray-800'
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-600">Loading interactions...</div>
  }

  return (
    <div className="space-y-4">
      <select
        value={filterStatus}
        onChange={(e) => setFilterStatus(e.target.value)}
        className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="all">All Status</option>
        <option value="pending">Pending</option>
        <option value="resolved">Resolved</option>
        <option value="escalated">Escalated</option>
      </select>

      <div className="space-y-3">
        {filteredInteractions.map((interaction) => (
          <div key={interaction.id} className="bg-white rounded-lg shadow p-4 border border-gray-200 hover:shadow-lg transition">
            <div className="flex justify-between items-start mb-3">
              <div>
                <h3 className="font-semibold text-gray-900">{interaction.subject}</h3>
                <p className="text-xs text-gray-600">{interaction.interaction_type.replace('_', ' ')} - {interaction.customer_id}</p>
              </div>
              <div className="flex gap-2">
                <span className={`px-2 py-1 text-xs font-medium rounded ${getPriorityBadge(interaction.priority)}`}>
                  {interaction.priority}
                </span>
                <span className={`px-2 py-1 text-xs font-medium rounded ${getStatusBadge(interaction.status)}`}>
                  {interaction.status}
                </span>
              </div>
            </div>

            <p className="text-sm text-gray-700 mb-3">{interaction.description}</p>

            <div className="grid grid-cols-2 gap-3 mb-4 text-sm">
              <div>
                <p className="text-gray-600 text-xs">Date</p>
                <p className="font-medium text-xs">{new Date(interaction.interaction_date).toLocaleDateString()}</p>
              </div>
              {interaction.follow_up_date && (
                <div>
                  <p className="text-gray-600 text-xs">Follow-up</p>
                  <p className="font-medium text-xs text-orange-600">{new Date(interaction.follow_up_date).toLocaleDateString()}</p>
                </div>
              )}
              {interaction.assigned_to_name && (
                <div>
                  <p className="text-gray-600 text-xs">Assigned To</p>
                  <p className="font-medium text-xs">{interaction.assigned_to_name}</p>
                </div>
              )}
            </div>

            <div className="flex gap-2 flex-wrap">
              <button
                onClick={() => onEdit(interaction)}
                className="px-3 py-2 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition"
              >
                Edit
              </button>
              <select
                value={interaction.status}
                onChange={(e) => onUpdateStatus(interaction, e.target.value)}
                className="px-3 py-2 text-xs font-medium border border-gray-300 rounded hover:bg-gray-50 transition"
              >
                <option value="pending">Pending</option>
                <option value="resolved">Resolved</option>
                <option value="escalated">Escalated</option>
              </select>
              <button
                onClick={() => {
                  if (confirm('Delete this interaction?')) {
                    onDelete(interaction)
                  }
                }}
                className="px-3 py-2 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition"
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>

      {filteredInteractions.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No customer interactions found</p>
        </div>
      )}
    </div>
  )
}
