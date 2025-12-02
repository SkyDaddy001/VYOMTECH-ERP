'use client'

import React, { useState, useEffect } from 'react'
import { CampaignIdea, MarketingTag } from '@/types/marketing'
import { marketingService } from '@/services/marketing.service'
import toast from 'react-hot-toast'

interface IdeaFormProps {
  idea?: CampaignIdea
  projectId: string
  onSuccess: () => void
  onCancel: () => void
}

const CHANNELS = ['email', 'sms', 'social', 'website', 'offline', 'events', 'broker_network'] as const
const SEGMENTS = ['residential', 'commercial', 'nri', 'corporate', 'investor', 'all'] as const
const PRIORITIES = ['low', 'medium', 'high', 'critical'] as const

export function IdeaForm({ idea, projectId, onSuccess, onCancel }: IdeaFormProps) {
  const [loading, setLoading] = useState(false)
  const [tags, setTags] = useState<MarketingTag[]>([])
  const [formData, setFormData] = useState<Partial<CampaignIdea>>(
    idea || {
      title: '',
      description: '',
      target_segment: 'residential',
      budget_estimate: 0,
      expected_leads: 0,
      channels: [],
      tags: [],
      priority: 'medium',
    }
  )

  useEffect(() => {
    loadTags()
  }, [])

  const loadTags = async () => {
    try {
      const data = await marketingService.getTags()
      setTags(data)
    } catch (error) {
      console.error('Failed to load tags:', error)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      if (idea?.id) {
        await marketingService.updateIdea(idea.id, formData)
        toast.success('Campaign idea updated successfully')
      } else {
        await marketingService.createIdea({
          ...formData,
          project_id: projectId,
        })
        toast.success('Campaign idea created successfully')
      }
      onSuccess()
    } catch (error) {
      toast.error('Failed to save campaign idea')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleChannelChange = (channel: typeof CHANNELS[number]) => {
    const channels = formData.channels || []
    if (channels.includes(channel)) {
      setFormData({
        ...formData,
        channels: channels.filter((c) => c !== channel),
      })
    } else {
      setFormData({
        ...formData,
        channels: [...channels, channel],
      })
    }
  }

  const handleTagChange = (tagId: string) => {
    const currentTags = formData.tags || []
    if (currentTags.includes(tagId)) {
      setFormData({
        ...formData,
        tags: currentTags.filter((t) => t !== tagId),
      })
    } else {
      setFormData({
        ...formData,
        tags: [...currentTags, tagId],
      })
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6 bg-white p-6 rounded-lg shadow">
      {/* Title */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Campaign Title *</label>
        <input
          type="text"
          required
          value={formData.title || ''}
          onChange={(e) => setFormData({ ...formData, title: e.target.value })}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
          placeholder="e.g., Summer 2024 Residential Launch"
        />
      </div>

      {/* Description */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Description</label>
        <textarea
          value={formData.description || ''}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          rows={3}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
          placeholder="Campaign objectives and strategy"
        />
      </div>

      {/* Target Segment */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Target Segment *</label>
        <select
          value={formData.target_segment || 'residential'}
          onChange={(e) =>
            setFormData({
              ...formData,
              target_segment: e.target.value as typeof SEGMENTS[number],
            })
          }
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
        >
          {SEGMENTS.map((seg) => (
            <option key={seg} value={seg}>
              {seg.charAt(0).toUpperCase() + seg.slice(1)}
            </option>
          ))}
        </select>
      </div>

      {/* Budget & Expected Leads Grid */}
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Budget Estimate (₹) *</label>
          <input
            type="number"
            required
            value={formData.budget_estimate || 0}
            onChange={(e) => setFormData({ ...formData, budget_estimate: parseFloat(e.target.value) })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
            min="0"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Expected Leads *</label>
          <input
            type="number"
            required
            value={formData.expected_leads || 0}
            onChange={(e) => setFormData({ ...formData, expected_leads: parseFloat(e.target.value) })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
            min="0"
          />
        </div>
      </div>

      {/* Channels */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Channels *</label>
        <div className="space-y-2">
          {CHANNELS.map((channel) => (
            <label key={channel} className="flex items-center">
              <input
                type="checkbox"
                checked={(formData.channels || []).includes(channel)}
                onChange={() => handleChannelChange(channel)}
                className="rounded border-gray-300 w-4 h-4"
              />
              <span className="ml-2 text-sm text-gray-700">
                {channel.charAt(0).toUpperCase() + channel.slice(1)}
              </span>
            </label>
          ))}
        </div>
      </div>

      {/* Tags */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Tags</label>
        <div className="space-y-2 max-h-40 overflow-y-auto">
          {tags.map((tag) => (
            <label key={tag.id} className="flex items-center">
              <input
                type="checkbox"
                checked={(formData.tags || []).includes(tag.id || '')}
                onChange={() => handleTagChange(tag.id || '')}
                className="rounded border-gray-300 w-4 h-4"
              />
              <span className="ml-2 text-sm text-gray-700">{tag.name}</span>
              {tag.color && (
                <span
                  className="ml-2 w-3 h-3 rounded-full"
                  style={{ backgroundColor: tag.color }}
                />
              )}
            </label>
          ))}
        </div>
      </div>

      {/* Priority */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Priority</label>
        <select
          value={formData.priority || 'medium'}
          onChange={(e) =>
            setFormData({
              ...formData,
              priority: e.target.value as typeof PRIORITIES[number],
            })
          }
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500"
        >
          {PRIORITIES.map((priority) => (
            <option key={priority} value={priority}>
              {priority.charAt(0).toUpperCase() + priority.slice(1)}
            </option>
          ))}
        </select>
      </div>

      {/* Actions */}
      <div className="flex gap-3 pt-4 border-t">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-2 px-4 rounded-md transition"
        >
          {loading ? 'Saving...' : idea ? 'Update Idea' : 'Create Idea'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium py-2 px-4 rounded-md transition"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
