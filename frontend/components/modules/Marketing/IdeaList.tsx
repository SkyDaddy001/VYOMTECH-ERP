'use client'

import React, { useState, useEffect } from 'react'
import { CampaignIdea } from '@/types/marketing'
import { marketingService } from '@/services/marketing.service'
import toast from 'react-hot-toast'

interface IdeaListProps {
  projectId: string
  onEdit: (idea: CampaignIdea) => void
  onCreateNew: () => void
}

const STATUS_COLORS: Record<string, string> = {
  draft: 'bg-gray-100 text-gray-800',
  approved: 'bg-green-100 text-green-800',
  in_execution: 'bg-blue-100 text-blue-800',
  completed: 'bg-purple-100 text-purple-800',
  archived: 'bg-red-100 text-red-800',
}

const PRIORITY_COLORS: Record<string, string> = {
  low: 'text-gray-600',
  medium: 'text-yellow-600',
  high: 'text-orange-600',
  critical: 'text-red-600',
}

export function IdeaList({ projectId, onEdit, onCreateNew }: IdeaListProps) {
  const [ideas, setIdeas] = useState<CampaignIdea[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<string>('all')

  useEffect(() => {
    loadIdeas()
  }, [])

  const loadIdeas = async () => {
    setLoading(true)
    try {
      const data = await marketingService.getIdeas()
      setIdeas(data.filter((idea) => idea.project_id === projectId))
    } catch (error) {
      toast.error('Failed to load campaign ideas')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleStatusChange = async (idea: CampaignIdea, newStatus: string) => {
    try {
      if (idea.id) {
        await marketingService.updateIdeaStatus(idea.id, newStatus)
        toast.success(`Idea status updated to ${newStatus}`)
        loadIdeas()
      }
    } catch (error) {
      toast.error('Failed to update idea status')
      console.error(error)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this idea?')) return

    try {
      await marketingService.deleteIdea(id)
      toast.success('Idea deleted successfully')
      loadIdeas()
    } catch (error) {
      toast.error('Failed to delete idea')
      console.error(error)
    }
  }

  const filteredIdeas =
    filter === 'all' ? ideas : ideas.filter((idea) => idea.status === filter)

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-gray-500">Loading campaign ideas...</div>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {/* Header with filters */}
      <div className="flex items-center justify-between">
        <div className="flex gap-2">
          <button
            onClick={() => setFilter('all')}
            className={`px-3 py-1 rounded text-sm ${
              filter === 'all'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            All ({ideas.length})
          </button>
          <button
            onClick={() => setFilter('draft')}
            className={`px-3 py-1 rounded text-sm ${
              filter === 'draft'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            Draft ({ideas.filter((i) => i.status === 'draft').length})
          </button>
          <button
            onClick={() => setFilter('approved')}
            className={`px-3 py-1 rounded text-sm ${
              filter === 'approved'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            Approved ({ideas.filter((i) => i.status === 'approved').length})
          </button>
          <button
            onClick={() => setFilter('in_execution')}
            className={`px-3 py-1 rounded text-sm ${
              filter === 'in_execution'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            Executing ({ideas.filter((i) => i.status === 'in_execution').length})
          </button>
        </div>

        <button
          onClick={onCreateNew}
          className="bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded-md transition"
        >
          + New Idea
        </button>
      </div>

      {/* Ideas List */}
      {filteredIdeas.length === 0 ? (
        <div className="bg-gray-50 p-8 rounded-lg text-center">
          <p className="text-gray-600 mb-4">No campaign ideas yet</p>
          <button
            onClick={onCreateNew}
            className="text-blue-600 hover:text-blue-700 font-medium"
          >
            Create the first one →
          </button>
        </div>
      ) : (
        <div className="space-y-3">
          {filteredIdeas.map((idea) => (
            <div key={idea.id} className="bg-white p-4 rounded-lg shadow hover:shadow-md transition">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="font-semibold text-gray-900">{idea.title}</h3>
                    <span className={`px-2 py-1 rounded text-xs font-medium ${STATUS_COLORS[idea.status]}`}>
                      {idea.status?.replace('_', ' ').toUpperCase()}
                    </span>
                    <span className={`text-sm font-medium ${PRIORITY_COLORS[idea.priority]}`}>
                      {idea.priority?.toUpperCase()}
                    </span>
                  </div>

                  {idea.description && (
                    <p className="text-sm text-gray-600 mb-2">{idea.description}</p>
                  )}

                  <div className="flex flex-wrap gap-4 text-sm text-gray-600 mb-3">
                    <div>
                      <span className="font-medium">Segment:</span> {idea.target_segment}
                    </div>
                    <div>
                      <span className="font-medium">Budget:</span> ₹{idea.budget_estimate?.toLocaleString()}
                    </div>
                    <div>
                      <span className="font-medium">Expected Leads:</span> {idea.expected_leads}
                    </div>
                  </div>

                  {/* Channels */}
                  {idea.channels && idea.channels.length > 0 && (
                    <div className="mb-3">
                      <span className="text-sm font-medium text-gray-600 mr-2">Channels:</span>
                      <div className="flex flex-wrap gap-1">
                        {idea.channels.map((channel) => (
                          <span
                            key={channel}
                            className="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded"
                          >
                            {channel}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Tags */}
                  {idea.tags && idea.tags.length > 0 && (
                    <div>
                      <span className="text-sm font-medium text-gray-600 mr-2">Tags:</span>
                      <div className="flex flex-wrap gap-1">
                        {idea.tags.map((tag) => (
                          <span
                            key={tag}
                            className="inline-block bg-purple-100 text-purple-800 text-xs px-2 py-1 rounded"
                          >
                            {tag}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                </div>

                {/* Actions */}
                <div className="ml-4 flex gap-2">
                  <button
                    onClick={() => onEdit(idea)}
                    className="text-blue-600 hover:text-blue-700 font-medium text-sm"
                  >
                    Edit
                  </button>

                  <div className="relative group">
                    <button className="text-green-600 hover:text-green-700 font-medium text-sm">
                      Move to
                    </button>
                    <div className="absolute right-0 top-full hidden group-hover:block bg-white border border-gray-300 rounded shadow-lg z-10">
                      {['approved', 'in_execution', 'completed', 'archived'].map((status) => (
                        <button
                          key={status}
                          onClick={() => handleStatusChange(idea, status)}
                          className="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                        >
                          {status.replace('_', ' ')}
                        </button>
                      ))}
                    </div>
                  </div>

                  <button
                    onClick={() => handleDelete(idea.id || '')}
                    className="text-red-600 hover:text-red-700 font-medium text-sm"
                  >
                    Delete
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
