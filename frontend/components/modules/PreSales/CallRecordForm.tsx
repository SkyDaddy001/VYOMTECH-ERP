'use client'

import { useState } from 'react'
import { CallRecord } from '@/types/presales'

interface CallRecordFormProps {
  onSubmit: (data: Partial<CallRecord>) => Promise<void>
  onCancel: () => void
}

export default function CallRecordForm({ onSubmit, onCancel }: CallRecordFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<CallRecord>>({
    call_date: new Date().toISOString().split('T')[0],
    call_time: new Date().toTimeString().slice(0, 5),
    customer_name: '',
    customer_phone: '',
    email: '',
    property_preference: [],
    property_status: [],
    configuration: [],
    area_required: 0,
    area_specified: false,
    budget_range: '',
    budget_specified: false,
    purchase_purpose: '',
    occupation_type: '',
    occupation_details: '',
    funding_source: '',
    interested_projects: [],
    other_requirements: '',
    call_summary: '',
    call_outcome: 'interested',
    follow_up_required: false,
    follow_up_date: '',
    presales_agent_id: '',
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await onSubmit(formData)
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (field: string, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  const toggleArrayValue = (field: string, value: string) => {
    setFormData(prev => {
      const current = (prev[field as keyof CallRecord] as string[]) || []
      if (current.includes(value)) {
        return { ...prev, [field]: current.filter(v => v !== value) }
      } else {
        return { ...prev, [field]: [...current, value] }
      }
    })
  }

  const isArraySelected = (field: string, value: string) => {
    const current = (formData[field as keyof CallRecord] as string[]) || []
    return current.includes(value)
  }

  const projects = ['League One', 'Luxe One', 'Sky Living', 'General Inquiry']

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      {/* Header Information */}
      <div>
        <h2 className="text-lg font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <span className="text-2xl">📞</span> Call Record
        </h2>
      </div>

      {/* Call Details */}
      <div className="bg-blue-50 rounded-lg p-4 border border-blue-200">
        <h3 className="text-sm font-semibold text-blue-900 mb-3 flex items-center gap-2">
          <span>🕐</span> Call Details
        </h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Call Date *</label>
            <input
              type="date"
              required
              value={formData.call_date || ''}
              onChange={(e) => handleChange('call_date', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Call Time *</label>
            <input
              type="time"
              required
              value={formData.call_time || ''}
              onChange={(e) => handleChange('call_time', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
        <div className="mt-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Presales Agent ID *</label>
          <input
            type="text"
            required
            value={formData.presales_agent_id || ''}
            onChange={(e) => handleChange('presales_agent_id', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Agent ID"
          />
        </div>
      </div>

      {/* Customer Contact Information */}
      <div className="bg-purple-50 rounded-lg p-4 border border-purple-200">
        <h3 className="text-sm font-semibold text-purple-900 mb-3 flex items-center gap-2">
          <span>👤</span> Customer Contact Information
        </h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Customer Name *</label>
            <input
              type="text"
              required
              value={formData.customer_name || ''}
              onChange={(e) => handleChange('customer_name', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
              placeholder="Enter name"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Phone Number *</label>
            <input
              type="tel"
              required
              value={formData.customer_phone || ''}
              onChange={(e) => handleChange('customer_phone', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
              placeholder="Phone number"
            />
          </div>
        </div>
        <div className="mt-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Email</label>
          <input
            type="email"
            value={formData.email || ''}
            onChange={(e) => handleChange('email', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
            placeholder="Email address"
          />
        </div>
      </div>

      {/* Property Preferences */}
      <div className="bg-green-50 rounded-lg p-4 border border-green-200">
        <h3 className="text-sm font-semibold text-green-900 mb-3 flex items-center gap-2">
          <span>🏠</span> Property Preferences
        </h3>
        
        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Property Type</label>
          <div className="grid grid-cols-3 gap-2">
            {['House', 'Apartment', 'Plot'].map(type => (
              <button
                key={type}
                type="button"
                onClick={() => toggleArrayValue('property_preference', type)}
                className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isArraySelected('property_preference', type)
                    ? 'bg-blue-600 text-white'
                    : 'bg-white border border-gray-300 text-gray-700 hover:border-gray-400'
                }`}
              >
                {type}
              </button>
            ))}
          </div>
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Property Status</label>
          <div className="grid grid-cols-3 gap-2">
            {['Under Construction', 'Ready to Move', 'Prelaunch'].map(status => (
              <button
                key={status}
                type="button"
                onClick={() => toggleArrayValue('property_status', status)}
                className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isArraySelected('property_status', status)
                    ? 'bg-blue-600 text-white'
                    : 'bg-white border border-gray-300 text-gray-700 hover:border-gray-400'
                }`}
              >
                {status}
              </button>
            ))}
          </div>
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Configuration</label>
          <div className="grid grid-cols-3 gap-2">
            {['2 BHK', '3 BHK', '4 BHK'].map(config => (
              <button
                key={config}
                type="button"
                onClick={() => toggleArrayValue('configuration', config)}
                className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isArraySelected('configuration', config)
                    ? 'bg-blue-600 text-white'
                    : 'bg-white border border-gray-300 text-gray-700 hover:border-gray-400'
                }`}
              >
                {config}
              </button>
            ))}
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Area Required (Sq. Ft.)</label>
          <div className="flex items-center gap-2">
            <input
              type="range"
              min="500"
              max="5000"
              step="100"
              value={formData.area_required || 500}
              onChange={(e) => {
                handleChange('area_required', parseInt(e.target.value))
                handleChange('area_specified', true)
              }}
              className="flex-1 h-2 bg-gray-300 rounded-lg appearance-none cursor-pointer"
            />
            <input
              type="number"
              value={formData.area_required || 500}
              onChange={(e) => {
                handleChange('area_required', parseInt(e.target.value))
                handleChange('area_specified', true)
              }}
              className="w-20 px-2 py-1 border border-gray-300 rounded-lg text-sm"
              min="500"
              max="5000"
            />
            <button
              type="button"
              onClick={() => {
                handleChange('area_required', 0)
                handleChange('area_specified', false)
              }}
              className="px-3 py-1 bg-gray-200 text-gray-700 rounded-lg text-sm hover:bg-gray-300"
            >
              NA
            </button>
          </div>
        </div>
      </div>

      {/* Customer Profile */}
      <div className="bg-orange-50 rounded-lg p-4 border border-orange-200">
        <h3 className="text-sm font-semibold text-orange-900 mb-3 flex items-center gap-2">
          <span>💰</span> Customer Profile
        </h3>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Budget Range</label>
          <select
            value={formData.budget_range || ''}
            onChange={(e) => {
              handleChange('budget_range', e.target.value)
              handleChange('budget_specified', !!e.target.value)
            }}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500"
          >
            <option value="">Select Budget Range</option>
            <option value="below_50_lakhs">Below 50 Lakhs</option>
            <option value="50_to_100_lakhs">50 - 100 Lakhs</option>
            <option value="1_to_2_crores">1 - 2 Crores</option>
            <option value="2_to_5_crores">2 - 5 Crores</option>
            <option value="above_5_crores">Above 5 Crores</option>
          </select>
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Purpose of Purchase</label>
          <div className="grid grid-cols-3 gap-2">
            {['Investment', 'Self Occupation', 'Rental'].map(purpose => (
              <button
                key={purpose}
                type="button"
                onClick={() => handleChange('purchase_purpose', purpose === formData.purchase_purpose ? '' : purpose)}
                className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  formData.purchase_purpose === purpose
                    ? 'bg-blue-600 text-white'
                    : 'bg-white border border-gray-300 text-gray-700 hover:border-gray-400'
                }`}
              >
                {purpose}
              </button>
            ))}
          </div>
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Occupation Type</label>
          <div className="grid grid-cols-3 gap-2">
            {['Salaried', 'Business', 'Self-Employed'].map(type => (
              <button
                key={type}
                type="button"
                onClick={() => handleChange('occupation_type', type === formData.occupation_type ? '' : type)}
                className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  formData.occupation_type === type
                    ? 'bg-blue-600 text-white'
                    : 'bg-white border border-gray-300 text-gray-700 hover:border-gray-400'
                }`}
              >
                {type}
              </button>
            ))}
          </div>
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Occupation Details</label>
          <input
            type="text"
            value={formData.occupation_details || ''}
            onChange={(e) => handleChange('occupation_details', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500"
            placeholder="e.g., IT Professional, Manufacturing Business"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Funding Source</label>
          <div className="grid grid-cols-2 gap-2">
            {['Own Funds', 'Bank Finance'].map(source => (
              <button
                key={source}
                type="button"
                onClick={() => handleChange('funding_source', source === formData.funding_source ? '' : source)}
                className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  formData.funding_source === source
                    ? 'bg-blue-600 text-white'
                    : 'bg-white border border-gray-300 text-gray-700 hover:border-gray-400'
                }`}
              >
                {source}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Projects & Notes */}
      <div className="bg-indigo-50 rounded-lg p-4 border border-indigo-200">
        <h3 className="text-sm font-semibold text-indigo-900 mb-3 flex items-center gap-2">
          <span>📋</span> Projects & Notes
        </h3>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Interested Projects</label>
          <div className="grid grid-cols-2 gap-2 mb-3">
            {projects.map(project => (
              <button
                key={project}
                type="button"
                onClick={() => toggleArrayValue('interested_projects', project)}
                className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isArraySelected('interested_projects', project)
                    ? 'bg-blue-600 text-white'
                    : 'bg-white border border-gray-300 text-gray-700 hover:border-gray-400'
                }`}
              >
                {project}
              </button>
            ))}
          </div>
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Other Requirements</label>
          <textarea
            value={formData.other_requirements || ''}
            onChange={(e) => handleChange('other_requirements', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="Any other requirements or preferences..."
            rows={3}
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Call Summary</label>
          <textarea
            value={formData.call_summary || ''}
            onChange={(e) => handleChange('call_summary', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="Summary of the call, key discussion points, customer feedback..."
            rows={4}
          />
        </div>
      </div>

      {/* Call Outcome & Follow-up */}
      <div className="bg-red-50 rounded-lg p-4 border border-red-200">
        <h3 className="text-sm font-semibold text-red-900 mb-3 flex items-center gap-2">
          <span>✅</span> Outcome & Follow-up
        </h3>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Call Outcome *</label>
          <select
            value={formData.call_outcome || 'interested'}
            onChange={(e) => handleChange('call_outcome', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500"
            required
          >
            <option value="interested">Interested</option>
            <option value="not_interested">Not Interested</option>
            <option value="maybe_later">Maybe Later</option>
            <option value="follow_up_needed">Follow-up Needed</option>
            <option value="converted">Converted to Booking</option>
          </select>
        </div>

        <div className="flex items-center gap-3 mb-4">
          <input
            type="checkbox"
            id="followUp"
            checked={formData.follow_up_required || false}
            onChange={(e) => handleChange('follow_up_required', e.target.checked)}
            className="w-4 h-4 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
          />
          <label htmlFor="followUp" className="text-sm font-medium text-gray-700">
            Follow-up Required
          </label>
        </div>

        {formData.follow_up_required && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Follow-up Date</label>
            <input
              type="date"
              value={formData.follow_up_date || ''}
              onChange={(e) => handleChange('follow_up_date', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500"
            />
          </div>
        )}
      </div>

      {/* Form Actions */}
      <div className="flex gap-3 pt-4">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium transition-colors"
        >
          {loading ? 'Saving...' : 'Save Call Record'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 bg-gray-200 text-gray-900 py-2 rounded-lg hover:bg-gray-300 font-medium transition-colors"
        >
          Cancel
        </button>
      </div>
    </form>
  )
}
