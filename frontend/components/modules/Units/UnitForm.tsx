'use client'

import { useState } from 'react'
import { PropertyUnit } from '@/types/unit'

interface UnitFormProps {
  unit?: PropertyUnit
  onSubmit: (data: Partial<PropertyUnit>) => Promise<void>
  onCancel: () => void
  loading?: boolean
}

export default function UnitForm({ unit, onSubmit, onCancel, loading = false }: UnitFormProps) {
  const [formData, setFormData] = useState({
    unit_number: unit?.unit_number || '',
    floor: unit?.floor || 1,
    unit_type: unit?.unit_type || '1BHK',
    facing: unit?.facing || 'north',
    carpet_area: unit?.carpet_area || 0,
    carpet_area_with_balcony: unit?.carpet_area_with_balcony || 0,
    utility_area: unit?.utility_area || 0,
    plinth_area: unit?.plinth_area || 0,
    sbua: unit?.sbua || 0,
    uds_sqft: unit?.uds_sqft || 0,
    status: unit?.status || 'available',
  })
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: name === 'floor' ? parseInt(value) : name.includes('area') || name === 'sbua' || name === 'uds_sqft' ? parseFloat(value) : value,
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!formData.unit_number.trim()) {
      setError('Unit number is required')
      return
    }

    try {
      setIsSubmitting(true)
      await onSubmit(formData)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save unit')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-lg max-w-2xl w-full p-6 max-h-[90vh] overflow-y-auto">
        <h2 className="text-2xl font-bold text-gray-800 mb-4">{unit ? 'Edit Unit' : 'Create Unit'}</h2>

        {error && <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded text-red-800 text-sm">{error}</div>}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Unit Number *</label>
              <input
                type="text"
                name="unit_number"
                value={formData.unit_number}
                onChange={handleChange}
                placeholder="e.g., A-101"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Floor</label>
              <input
                type="number"
                name="floor"
                value={formData.floor}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                disabled={isSubmitting}
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Unit Type</label>
              <select
                name="unit_type"
                value={formData.unit_type}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                disabled={isSubmitting}
              >
                <option>1BHK</option>
                <option>2BHK</option>
                <option>3BHK</option>
                <option>Shop</option>
                <option>Office</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Facing</label>
              <select
                name="facing"
                value={formData.facing}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                disabled={isSubmitting}
              >
                <option>North</option>
                <option>South</option>
                <option>East</option>
                <option>West</option>
                <option>Corner</option>
              </select>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Carpet Area (sq.ft)</label>
              <input
                type="number"
                name="carpet_area"
                value={formData.carpet_area}
                onChange={handleChange}
                step="0.01"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">With Balcony (sq.ft)</label>
              <input
                type="number"
                name="carpet_area_with_balcony"
                value={formData.carpet_area_with_balcony}
                onChange={handleChange}
                step="0.01"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                disabled={isSubmitting}
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Utility Area (sq.ft)</label>
              <input
                type="number"
                name="utility_area"
                value={formData.utility_area}
                onChange={handleChange}
                step="0.01"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
              <select
                name="status"
                value={formData.status}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                disabled={isSubmitting}
              >
                <option value="available">Available</option>
                <option value="booked">Booked</option>
                <option value="sold">Sold</option>
                <option value="reserved">Reserved</option>
              </select>
            </div>
          </div>

          <div className="pt-4 flex gap-3">
            <button
              type="button"
              onClick={onCancel}
              disabled={isSubmitting}
              className="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded-lg transition disabled:opacity-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition disabled:opacity-50"
            >
              {isSubmitting ? 'Saving...' : 'Save'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
