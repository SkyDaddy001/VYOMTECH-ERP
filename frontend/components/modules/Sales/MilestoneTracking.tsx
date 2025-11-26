'use client'

import React, { useState, useEffect } from 'react'
import { Plus, Trash2, Calendar, MapPin, User, Clock, AlertCircle, CheckCircle } from 'lucide-react'
import { formatDateToDDMMMYYYY, formatDateToInput } from '@/lib/dateFormat'

interface Milestone {
  id: string
  lead_id: string
  milestone_type: string
  milestone_date: string
  milestone_time?: string
  notes?: string
  location_name?: string
  visited_by?: string
  outcome?: string
  status_before?: string
  status_after?: string
  follow_up_required?: boolean
}

interface Engagement {
  id: string
  lead_id: string
  engagement_type: string
  engagement_date: string
  engagement_channel?: string
  subject?: string
  status: string
  response_received: boolean
  response_date?: string
}

export function MilestoneTracking() {
  const [activeTab, setActiveTab] = useState<'milestones' | 'engagement'>('milestones')
  const [leadId, setLeadId] = useState('')
  const [milestones, setMilestones] = useState<Milestone[]>([])
  const [engagements, setEngagements] = useState<Engagement[]>([])
  const [loading, setLoading] = useState(false)
  const [showNewMilestone, setShowNewMilestone] = useState(false)
  const [showNewEngagement, setShowNewEngagement] = useState(false)

  const [newMilestone, setNewMilestone] = useState({
    lead_id: '',
    milestone_type: 'contacted',
    milestone_date: new Date().toISOString().split('T')[0],
    milestone_time: '',
    notes: '',
    location_name: '',
    visited_by: '',
    outcome: 'neutral',
    follow_up_required: false,
  })

  const [newEngagement, setNewEngagement] = useState({
    lead_id: '',
    engagement_type: 'email_sent',
    engagement_channel: 'email',
    subject: '',
    notes: '',
    status: 'completed',
  })

  const tenantId = localStorage.getItem('tenantId') || ''
  const userId = localStorage.getItem('userId') || ''

  const fetchMilestones = async (id: string) => {
    if (!id || !tenantId) return

    setLoading(true)
    try {
      const response = await fetch(`/api/v1/sales/milestones/lead/${id}`, {
        headers: {
          'X-Tenant-ID': tenantId,
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      })
      if (response.ok) {
        const data = await response.json()
        setMilestones(data || [])
      }
    } catch (error) {
      console.error('Failed to fetch milestones:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchEngagements = async (id: string) => {
    if (!id || !tenantId) return

    setLoading(true)
    try {
      const response = await fetch(`/api/v1/sales/engagement/${id}`, {
        headers: {
          'X-Tenant-ID': tenantId,
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      })
      if (response.ok) {
        const data = await response.json()
        setEngagements(data || [])
      }
    } catch (error) {
      console.error('Failed to fetch engagements:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleLeadIdSubmit = () => {
    if (leadId) {
      setNewMilestone({ ...newMilestone, lead_id: leadId })
      setNewEngagement({ ...newEngagement, lead_id: leadId })
      fetchMilestones(leadId)
      fetchEngagements(leadId)
    }
  }

  const handleCreateMilestone = async () => {
    if (!newMilestone.lead_id) {
      alert('Please enter a Lead ID first')
      return
    }

    try {
      const response = await fetch('/api/v1/sales/milestones/lead', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': tenantId,
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify({
          ...newMilestone,
          tenant_id: tenantId,
          created_by: userId,
        }),
      })

      if (response.ok) {
        alert('Milestone created successfully')
        setShowNewMilestone(false)
        setNewMilestone({
          lead_id: leadId,
          milestone_type: 'contacted',
          milestone_date: new Date().toISOString().split('T')[0],
          milestone_time: '',
          notes: '',
          location_name: '',
          visited_by: '',
          outcome: 'neutral',
          follow_up_required: false,
        })
        fetchMilestones(leadId)
      }
    } catch (error) {
      console.error('Failed to create milestone:', error)
      alert('Failed to create milestone')
    }
  }

  const handleCreateEngagement = async () => {
    if (!newEngagement.lead_id) {
      alert('Please enter a Lead ID first')
      return
    }

    try {
      const response = await fetch('/api/v1/sales/engagement', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': tenantId,
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify({
          ...newEngagement,
          tenant_id: tenantId,
          created_by: userId,
        }),
      })

      if (response.ok) {
        alert('Engagement recorded successfully')
        setShowNewEngagement(false)
        setNewEngagement({
          lead_id: leadId,
          engagement_type: 'email_sent',
          engagement_channel: 'email',
          subject: '',
          notes: '',
          status: 'completed',
        })
        fetchEngagements(leadId)
      }
    } catch (error) {
      console.error('Failed to create engagement:', error)
      alert('Failed to create engagement')
    }
  }

  const getMilestoneIcon = (type: string) => {
    const icons: Record<string, React.ReactNode> = {
      lead_generated: <CheckCircle className="w-4 h-4 text-blue-500" />,
      contacted: <User className="w-4 h-4 text-green-500" />,
      site_visit: <MapPin className="w-4 h-4 text-purple-500" />,
      demo: <Calendar className="w-4 h-4 text-orange-500" />,
      booking: <CheckCircle className="w-4 h-4 text-red-500" />,
      cancellation: <AlertCircle className="w-4 h-4 text-red-600" />,
    }
    return icons[type] || <Calendar className="w-4 h-4" />
  }

  return (
    <div className="space-y-6">
      {/* Lead ID Input */}
      <div className="bg-white p-4 rounded-lg border border-gray-200">
        <label className="block text-sm font-medium text-gray-700 mb-2">Enter Lead ID</label>
        <div className="flex gap-2">
          <input
            type="text"
            value={leadId}
            onChange={(e) => setLeadId(e.target.value)}
            placeholder="Enter Lead ID to view milestones and engagement"
            className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
          />
          <button
            onClick={handleLeadIdSubmit}
            className="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600"
          >
            Load
          </button>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex gap-4 border-b border-gray-200">
        <button
          onClick={() => setActiveTab('milestones')}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === 'milestones'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          Milestones
        </button>
        <button
          onClick={() => setActiveTab('engagement')}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === 'engagement'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          Engagement
        </button>
      </div>

      {/* Milestones Section */}
      {activeTab === 'milestones' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="font-semibold text-lg">Milestone Timeline</h3>
            <button
              onClick={() => setShowNewMilestone(!showNewMilestone)}
              className="flex items-center gap-2 bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600"
            >
              <Plus className="w-4 h-4" /> New Milestone
            </button>
          </div>

          {showNewMilestone && (
            <div className="bg-gray-50 p-4 rounded-lg border border-gray-200 space-y-3">
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Milestone Type</label>
                  <select
                    value={newMilestone.milestone_type}
                    onChange={(e) => setNewMilestone({ ...newMilestone, milestone_type: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  >
                    <option value="lead_generated">Lead Generated</option>
                    <option value="contacted">Contacted</option>
                    <option value="site_visit">Site Visit</option>
                    <option value="revisit">Re-visit</option>
                    <option value="demo">Demo</option>
                    <option value="proposal">Proposal</option>
                    <option value="negotiation">Negotiation</option>
                    <option value="booking">Booking</option>
                    <option value="cancellation">Cancellation</option>
                    <option value="reengaged">Re-engaged</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Milestone Date</label>
                  <input
                    type="date"
                    value={newMilestone.milestone_date}
                    onChange={(e) => setNewMilestone({ ...newMilestone, milestone_date: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Location Name</label>
                <input
                  type="text"
                  value={newMilestone.location_name}
                  onChange={(e) => setNewMilestone({ ...newMilestone, location_name: e.target.value })}
                  placeholder="e.g., Bangalore Office, Site A"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Visit Outcome</label>
                <select
                  value={newMilestone.outcome}
                  onChange={(e) => setNewMilestone({ ...newMilestone, outcome: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                >
                  <option value="positive">Positive</option>
                  <option value="neutral">Neutral</option>
                  <option value="negative">Negative</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
                <textarea
                  value={newMilestone.notes}
                  onChange={(e) => setNewMilestone({ ...newMilestone, notes: e.target.value })}
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                />
              </div>

              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={newMilestone.follow_up_required}
                  onChange={(e) => setNewMilestone({ ...newMilestone, follow_up_required: e.target.checked })}
                  className="rounded"
                />
                <label className="text-sm text-gray-700">Follow-up Required</label>
              </div>

              <div className="flex gap-2">
                <button
                  onClick={handleCreateMilestone}
                  className="bg-green-500 text-white px-4 py-2 rounded-lg hover:bg-green-600"
                >
                  Create Milestone
                </button>
                <button
                  onClick={() => setShowNewMilestone(false)}
                  className="bg-gray-300 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-400"
                >
                  Cancel
                </button>
              </div>
            </div>
          )}

          <div className="space-y-3">
            {loading ? (
              <p className="text-gray-500">Loading milestones...</p>
            ) : milestones.length === 0 ? (
              <p className="text-gray-500">No milestones found. Create one to get started.</p>
            ) : (
              milestones.map((milestone) => (
                <div key={milestone.id} className="bg-white p-4 rounded-lg border border-gray-200 hover:shadow-md transition">
                  <div className="flex items-start gap-3">
                    {getMilestoneIcon(milestone.milestone_type)}
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <h4 className="font-medium text-gray-900 capitalize">
                          {milestone.milestone_type.replace(/_/g, ' ')}
                        </h4>
                        <span className="text-sm text-gray-500">
                          {formatDateToDDMMMYYYY(milestone.milestone_date)}
                        </span>
                      </div>
                      {milestone.location_name && (
                        <p className="text-sm text-gray-600 flex items-center gap-1 mt-1">
                          <MapPin className="w-3 h-3" /> {milestone.location_name}
                        </p>
                      )}
                      {milestone.outcome && (
                        <p className="text-sm text-gray-600 capitalize">Outcome: {milestone.outcome}</p>
                      )}
                      {milestone.notes && <p className="text-sm text-gray-700 mt-2">{milestone.notes}</p>}
                      {milestone.follow_up_required && (
                        <p className="text-sm text-red-600 mt-2 flex items-center gap-1">
                          <AlertCircle className="w-3 h-3" /> Follow-up required
                        </p>
                      )}
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      )}

      {/* Engagement Section */}
      {activeTab === 'engagement' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="font-semibold text-lg">Engagement History</h3>
            <button
              onClick={() => setShowNewEngagement(!showNewEngagement)}
              className="flex items-center gap-2 bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600"
            >
              <Plus className="w-4 h-4" /> New Engagement
            </button>
          </div>

          {showNewEngagement && (
            <div className="bg-gray-50 p-4 rounded-lg border border-gray-200 space-y-3">
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Engagement Type</label>
                  <select
                    value={newEngagement.engagement_type}
                    onChange={(e) => setNewEngagement({ ...newEngagement, engagement_type: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  >
                    <option value="email_sent">Email Sent</option>
                    <option value="call_made">Call Made</option>
                    <option value="message_sent">Message Sent</option>
                    <option value="meeting_scheduled">Meeting Scheduled</option>
                    <option value="proposal_sent">Proposal Sent</option>
                    <option value="quote_sent">Quote Sent</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Channel</label>
                  <select
                    value={newEngagement.engagement_channel}
                    onChange={(e) => setNewEngagement({ ...newEngagement, engagement_channel: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  >
                    <option value="email">Email</option>
                    <option value="phone">Phone</option>
                    <option value="sms">SMS</option>
                    <option value="whatsapp">WhatsApp</option>
                    <option value="in_person">In Person</option>
                    <option value="video">Video Call</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Subject</label>
                <input
                  type="text"
                  value={newEngagement.subject}
                  onChange={(e) => setNewEngagement({ ...newEngagement, subject: e.target.value })}
                  placeholder="What was the engagement about?"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
                <textarea
                  value={newEngagement.notes}
                  onChange={(e) => setNewEngagement({ ...newEngagement, notes: e.target.value })}
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                />
              </div>

              <div className="flex gap-2">
                <button
                  onClick={handleCreateEngagement}
                  className="bg-green-500 text-white px-4 py-2 rounded-lg hover:bg-green-600"
                >
                  Record Engagement
                </button>
                <button
                  onClick={() => setShowNewEngagement(false)}
                  className="bg-gray-300 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-400"
                >
                  Cancel
                </button>
              </div>
            </div>
          )}

          <div className="space-y-3">
            {loading ? (
              <p className="text-gray-500">Loading engagements...</p>
            ) : engagements.length === 0 ? (
              <p className="text-gray-500">No engagement records found. Record your first engagement.</p>
            ) : (
              engagements.map((engagement) => (
                <div key={engagement.id} className="bg-white p-4 rounded-lg border border-gray-200 hover:shadow-md transition">
                  <div className="flex items-start gap-3">
                    <Clock className="w-4 h-4 text-blue-500 mt-1" />
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <h4 className="font-medium text-gray-900 capitalize">
                          {engagement.engagement_type.replace(/_/g, ' ')}
                        </h4>
                        <span className={`text-xs px-2 py-1 rounded-full ${
                          engagement.status === 'completed' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                        }`}>
                          {engagement.status}
                        </span>
                      </div>
                      <p className="text-sm text-gray-600 capitalize mt-1">
                        Channel: {engagement.engagement_channel}
                      </p>
                      {engagement.subject && <p className="text-sm font-medium mt-2">{engagement.subject}</p>}
                      <p className="text-xs text-gray-500 mt-2">
                        {new Date(engagement.engagement_date).toLocaleString()}
                      </p>
                      {engagement.response_received && (
                        <p className="text-sm text-green-600 mt-2 flex items-center gap-1">
                          <CheckCircle className="w-3 h-3" /> Response received
                        </p>
                      )}
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}
