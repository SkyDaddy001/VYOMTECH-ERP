'use client'

import React, { useState } from 'react'
import { CustomerInteraction } from '@/types/postsales'

interface InteractionFormProps {
  interaction?: CustomerInteraction
  onSubmit: (data: Partial<CustomerInteraction>) => Promise<void>
  onCancel: () => void
}

export default function InteractionForm({ interaction, onSubmit, onCancel }: InteractionFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<CustomerInteraction>>(
    interaction || {
      booking_id: '',
      customer_id: '',
      interaction_date: new Date().toISOString().split('T')[0],
      interaction_type: 'call',
      subject: '',
      description: '',
      status: 'pending',
      priority: 'medium',
    }
  )

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await onSubmit(formData)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4 max-w-2xl">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Booking ID *</label>
          <input
            type="text"
            required
            value={formData.booking_id || ''}
            onChange={(e) => setFormData({ ...formData, booking_id: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Customer ID *</label>
          <input
            type="text"
            required
            value={formData.customer_id || ''}
            onChange={(e) => setFormData({ ...formData, customer_id: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Interaction Type *</label>
          <select
            required
            value={formData.interaction_type || 'call'}
            onChange={(e) => setFormData({ ...formData, interaction_type: e.target.value as any })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="call">Call</option>
            <option value="email">Email</option>
            <option value="meeting">Meeting</option>
            <option value="document">Document</option>
            <option value="payment">Payment</option>
            <option value="query">Query</option>
            <option value="escalation">Escalation</option>
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Interaction Date *</label>
          <input
            type="date"
            required
            value={formData.interaction_date || ''}
            onChange={(e) => setFormData({ ...formData, interaction_date: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-700 mb-1">Subject *</label>
          <input
            type="text"
            required
            value={formData.subject || ''}
            onChange={(e) => setFormData({ ...formData, subject: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-700 mb-1">Description *</label>
          <textarea
            required
            value={formData.description || ''}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            rows={3}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Status *</label>
          <select
            required
            value={formData.status || 'pending'}
            onChange={(e) => setFormData({ ...formData, status: e.target.value as any })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="pending">Pending</option>
            <option value="resolved">Resolved</option>
            <option value="escalated">Escalated</option>
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Priority *</label>
          <select
            required
            value={formData.priority || 'medium'}
            onChange={(e) => setFormData({ ...formData, priority: e.target.value as any })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="critical">Critical</option>
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Assigned To</label>
          <input
            type="text"
            value={formData.assigned_to_name || ''}
            onChange={(e) => setFormData({ ...formData, assigned_to_name: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Follow-up Date</label>
          <input
            type="date"
            value={formData.follow_up_date || ''}
            onChange={(e) => setFormData({ ...formData, follow_up_date: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
          <textarea
            value={formData.notes || ''}
            onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
            rows={2}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <div className="flex gap-3">
        <button
          type="submit"
          disabled={loading}
          className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium disabled:opacity-50"
        >
          {loading ? 'Saving...' : interaction ? 'Update Interaction' : 'Log Interaction'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 font-medium"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
