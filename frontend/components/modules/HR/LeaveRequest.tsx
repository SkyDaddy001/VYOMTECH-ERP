'use client'

import { useState } from 'react'
import toast from 'react-hot-toast'
import { Leave } from '@/types/hr'

interface LeaveRequestProps {
  onSubmit: (data: Partial<Leave>) => Promise<void>
  onCancel: () => void
}

export default function LeaveRequest({ onSubmit, onCancel }: LeaveRequestProps) {
  const [formData, setFormData] = useState<Partial<Leave>>({
    from_date: '',
    to_date: '',
    reason: '',
    leave_type: 'annual',
    status: 'pending',
  })
  const [loading, setLoading] = useState(false)

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: name === 'number_of_days' ? (value ? parseInt(value) : 0) : value,
    }))
  }

  const calculateDays = () => {
    if (formData.from_date && formData.to_date) {
      const from = new Date(formData.from_date)
      const to = new Date(formData.to_date)
      const days = Math.ceil((to.getTime() - from.getTime()) / (1000 * 60 * 60 * 24)) + 1
      setFormData((prev) => ({
        ...prev,
        number_of_days: Math.max(0, days),
      }))
    }
  }

  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleChange(e)
    calculateDays()
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!formData.from_date || !formData.to_date || !formData.reason) {
      toast.error('Please fill all required fields')
      return
    }

    setLoading(true)
    try {
      await onSubmit(formData)
      toast.success('Leave request submitted!')
      onCancel()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error submitting leave request')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Leave Dates */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Leave Dates</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">From Date *</label>
            <input
              type="date"
              name="from_date"
              value={formData.from_date || ''}
              onChange={handleDateChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">To Date *</label>
            <input
              type="date"
              name="to_date"
              value={formData.to_date || ''}
              onChange={handleDateChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
        </div>
      </div>

      {/* Number of Days */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Number of Days</label>
        <input
          type="number"
          value={formData.number_of_days || 0}
          disabled
          className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-600"
        />
      </div>

      {/* Leave Type */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Leave Details</h3>
        <label className="block text-sm font-medium text-gray-700 mb-1">Leave Type *</label>
        <select
          name="leave_type"
          value={formData.leave_type || 'annual'}
          onChange={handleChange}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          required
        >
          <option value="annual">Annual Leave</option>
          <option value="sick">Sick Leave</option>
          <option value="personal">Personal Leave</option>
          <option value="maternity">Maternity Leave</option>
          <option value="paternity">Paternity Leave</option>
          <option value="unpaid">Unpaid Leave</option>
        </select>
      </div>

      {/* Reason */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Reason for Leave *</label>
        <textarea
          name="reason"
          value={formData.reason || ''}
          onChange={handleChange}
          rows={4}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Please provide the reason for your leave..."
          required
        />
      </div>

      {/* Actions */}
      <div className="flex gap-4 pt-4 border-t border-gray-200">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
        >
          {loading ? 'Submitting...' : 'Submit Leave Request'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 bg-gray-200 text-gray-900 py-2 rounded-lg hover:bg-gray-300 font-medium"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
