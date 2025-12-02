'use client'

import { useState } from 'react'
import { Property } from '@/types/realEstate'

interface PropertyManagementFormProps {
  property?: Property | null
  onSubmit: (data: Partial<Property>) => Promise<void>
  onCancel: () => void
}

export default function PropertyManagementForm({ property, onSubmit, onCancel }: PropertyManagementFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<Property>>(
    property || {
      property_id: '',
      project_id: '',
      unit_type: 'residential',
      wing: '',
      floor: '',
      unit_number: '',
      property_status: 'available',
      super_area: 0,
      carpet_area: 0,
      builtup_area: 0,
      terrace_area: 0,
      parking_count: 0,
      facing: 'north',
      configuration: '2bhk',
      base_price: 0,
      base_price_per_sqft: 0,
      possession_date_expected: '',
      possession_date_actual: '',
      ownership_type: 'freehold',
      notes: '',
    }
  )

  const handleChange = (field: string, value: any) => {
    setFormData(prev => {
      const updated = { ...prev, [field]: value }
      // Auto-calculate base_price_per_sqft if super_area or base_price changes
      if (field === 'super_area' || field === 'base_price') {
        const superArea = field === 'super_area' ? value : updated.super_area || 0
        const basePrice = field === 'base_price' ? value : updated.base_price || 0
        if (superArea > 0) {
          updated.base_price_per_sqft = basePrice / superArea
        }
      }
      return updated
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await onSubmit(formData)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <span className="text-2xl">🏠</span> {property ? 'Edit Property' : 'Add Property'}
        </h2>
      </div>

      {/* Property Identity */}
      <div className="bg-blue-50 rounded-lg p-4 border border-blue-200">
        <h3 className="text-sm font-semibold text-blue-900 mb-4">Property Identity</h3>
        <div className="grid grid-cols-4 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Property ID *</label>
            <input
              type="text"
              required
              value={formData.property_id || ''}
              onChange={(e) => handleChange('property_id', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="PROP-001"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Project ID</label>
            <input
              type="text"
              value={formData.project_id || ''}
              onChange={(e) => handleChange('project_id', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="PRJ-001"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Unit Type</label>
            <select
              value={formData.unit_type || 'residential'}
              onChange={(e) => handleChange('unit_type', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="residential">Residential</option>
              <option value="commercial">Commercial</option>
              <option value="parking">Parking</option>
              <option value="retail">Retail</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
            <select
              value={formData.property_status || 'available'}
              onChange={(e) => handleChange('property_status', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="available">Available</option>
              <option value="booked">Booked</option>
              <option value="sold">Sold</option>
              <option value="under_construction">Under Construction</option>
              <option value="ready">Ready to Possess</option>
            </select>
          </div>
        </div>
      </div>

      {/* Location Details */}
      <div className="bg-green-50 rounded-lg p-4 border border-green-200">
        <h3 className="text-sm font-semibold text-green-900 mb-4">Location Details</h3>
        <div className="grid grid-cols-4 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Wing</label>
            <input
              type="text"
              value={formData.wing || ''}
              onChange={(e) => handleChange('wing', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              placeholder="A, B, C"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Floor</label>
            <input
              type="number"
              value={formData.floor || ''}
              onChange={(e) => handleChange('floor', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              placeholder="Floor"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Unit Number</label>
            <input
              type="text"
              value={formData.unit_number || ''}
              onChange={(e) => handleChange('unit_number', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              placeholder="101, 102"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Facing</label>
            <select
              value={formData.facing || 'north'}
              onChange={(e) => handleChange('facing', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
            >
              <option value="north">North</option>
              <option value="south">South</option>
              <option value="east">East</option>
              <option value="west">West</option>
              <option value="northeast">Northeast</option>
              <option value="northwest">Northwest</option>
              <option value="southeast">Southeast</option>
              <option value="southwest">Southwest</option>
            </select>
          </div>
        </div>
      </div>

      {/* Configuration & Areas */}
      <div className="bg-yellow-50 rounded-lg p-4 border border-yellow-200">
        <h3 className="text-sm font-semibold text-yellow-900 mb-4">Configuration & Areas</h3>
        <div className="grid grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Configuration</label>
            <select
              value={formData.configuration || '2bhk'}
              onChange={(e) => handleChange('configuration', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500"
            >
              <option value="studio">Studio</option>
              <option value="1bhk">1 BHK</option>
              <option value="2bhk">2 BHK</option>
              <option value="3bhk">3 BHK</option>
              <option value="4bhk">4 BHK</option>
              <option value="5bhk">5 BHK</option>
              <option value="penthouse">Penthouse</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Super Area (sq.ft) *</label>
            <input
              type="number"
              required
              value={formData.super_area || ''}
              onChange={(e) => handleChange('super_area', parseFloat(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500"
              step="0.01"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Carpet Area (sq.ft)</label>
            <input
              type="number"
              value={formData.carpet_area || ''}
              onChange={(e) => handleChange('carpet_area', parseFloat(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500"
              step="0.01"
            />
          </div>
        </div>
        <div className="grid grid-cols-3 gap-4 mt-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Builtup Area (sq.ft)</label>
            <input
              type="number"
              value={formData.builtup_area || ''}
              onChange={(e) => handleChange('builtup_area', parseFloat(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500"
              step="0.01"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Terrace Area (sq.ft)</label>
            <input
              type="number"
              value={formData.terrace_area || ''}
              onChange={(e) => handleChange('terrace_area', parseFloat(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500"
              step="0.01"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Parking Count</label>
            <input
              type="number"
              value={formData.parking_count || ''}
              onChange={(e) => handleChange('parking_count', parseInt(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500"
              min="0"
            />
          </div>
        </div>
      </div>

      {/* Pricing */}
      <div className="bg-purple-50 rounded-lg p-4 border border-purple-200">
        <h3 className="text-sm font-semibold text-purple-900 mb-4">Pricing</h3>
        <div className="grid grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Base Price *</label>
            <input
              type="number"
              required
              value={formData.base_price || ''}
              onChange={(e) => handleChange('base_price', parseFloat(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
              step="0.01"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Price per sq.ft</label>
            <input
              type="number"
              value={formData.base_price_per_sqft || ''}
              onChange={(e) => handleChange('base_price_per_sqft', parseFloat(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 bg-gray-100"
              step="0.01"
              readOnly
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Ownership Type</label>
            <select
              value={formData.ownership_type || 'freehold'}
              onChange={(e) => handleChange('ownership_type', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
            >
              <option value="freehold">Freehold</option>
              <option value="leasehold">Leasehold</option>
            </select>
          </div>
        </div>
      </div>

      {/* Possession Dates */}
      <div className="bg-indigo-50 rounded-lg p-4 border border-indigo-200">
        <h3 className="text-sm font-semibold text-indigo-900 mb-4">Possession Dates</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Expected Possession Date</label>
            <input
              type="date"
              value={formData.possession_date_expected || ''}
              onChange={(e) => handleChange('possession_date_expected', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Actual Possession Date</label>
            <input
              type="date"
              value={formData.possession_date_actual || ''}
              onChange={(e) => handleChange('possession_date_actual', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>
        </div>
      </div>

      {/* Notes */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Notes</label>
        <textarea
          value={formData.notes || ''}
          onChange={(e) => handleChange('notes', e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Property details and notes..."
          rows={3}
        />
      </div>

      {/* Form Actions */}
      <div className="flex gap-3 pt-4">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium transition-colors"
        >
          {loading ? 'Saving...' : property ? 'Update Property' : 'Add Property'}
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
