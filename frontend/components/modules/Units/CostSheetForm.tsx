'use client'

import { useState } from 'react'
import { UnitCostSheet } from '@/types/unit'

interface CostSheetFormProps {
  costSheet?: UnitCostSheet
  onSubmit: (data: Partial<UnitCostSheet>) => Promise<void>
  onCancel: () => void
  loading?: boolean
}

export default function CostSheetForm({ costSheet, onSubmit, onCancel, loading = false }: CostSheetFormProps) {
  const [formData, setFormData] = useState({
    rate_per_sqft: costSheet?.rate_per_sqft || 0,
    sbua_rate: costSheet?.sbua_rate || 0,
    base_price: costSheet?.base_price || 0,
    frc: costSheet?.frc || 0,
    car_parking_cost: costSheet?.car_parking_cost || 0,
    plc: costSheet?.plc || 0,
    statutory_charges: costSheet?.statutory_charges || 0,
    other_charges: costSheet?.other_charges || 0,
    legal_charges: costSheet?.legal_charges || 0,
    composite_guideline_value: costSheet?.composite_guideline_value || 0,
    actual_sold_price: costSheet?.actual_sold_price || 0,
    car_parking_type: costSheet?.car_parking_type || 'open',
    parking_location: costSheet?.parking_location || '',
  })
  const [error, setError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: name === 'car_parking_type' || name === 'parking_location' ? value : parseFloat(value) || 0,
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    try {
      setIsSubmitting(true)
      await onSubmit(formData)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save cost sheet')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-lg max-w-3xl w-full p-6 max-h-[90vh] overflow-y-auto">
        <h2 className="text-2xl font-bold text-gray-800 mb-4">Unit Cost Sheet</h2>

        {error && <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded text-red-800 text-sm">{error}</div>}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="border-b pb-4">
            <h3 className="font-semibold text-gray-800 mb-3">Rates</h3>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Rate/sq.ft</label>
                <input
                  type="number"
                  name="rate_per_sqft"
                  value={formData.rate_per_sqft}
                  onChange={handleChange}
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">SBUA Rate</label>
                <input
                  type="number"
                  name="sbua_rate"
                  value={formData.sbua_rate}
                  onChange={handleChange}
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Base Price</label>
                <input
                  type="number"
                  name="base_price"
                  value={formData.base_price}
                  onChange={handleChange}
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
            </div>
          </div>

          <div className="border-b pb-4">
            <h3 className="font-semibold text-gray-800 mb-3">Charges</h3>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">FRC</label>
                <input
                  type="number"
                  name="frc"
                  value={formData.frc}
                  onChange={handleChange}
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Car Parking</label>
                <input
                  type="number"
                  name="car_parking_cost"
                  value={formData.car_parking_cost}
                  onChange={handleChange}
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">PLC</label>
                <input
                  type="number"
                  name="plc"
                  value={formData.plc}
                  onChange={handleChange}
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Statutory Charges</label>
                <input
                  type="number"
                  name="statutory_charges"
                  value={formData.statutory_charges}
                  onChange={handleChange}
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Legal Charges</label>
                <input
                  type="number"
                  name="legal_charges"
                  value={formData.legal_charges}
                  onChange={handleChange}
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Other Charges</label>
                <input
                  type="number"
                  name="other_charges"
                  value={formData.other_charges}
                  onChange={handleChange}
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
            </div>
          </div>

          <div className="border-b pb-4">
            <h3 className="font-semibold text-gray-800 mb-3">Prices</h3>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Composite Guideline Value</label>
                <input
                  type="number"
                  name="composite_guideline_value"
                  value={formData.composite_guideline_value}
                  onChange={handleChange}
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Actual Sold Price</label>
                <input
                  type="number"
                  name="actual_sold_price"
                  value={formData.actual_sold_price}
                  onChange={handleChange}
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
            </div>
          </div>

          <div className="pb-4">
            <h3 className="font-semibold text-gray-800 mb-3">Parking</h3>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Parking Type</label>
                <select
                  name="car_parking_type"
                  value={formData.car_parking_type}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                >
                  <option value="covered">Covered</option>
                  <option value="open">Open</option>
                  <option value="none">None</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Parking Location</label>
                <input
                  type="text"
                  name="parking_location"
                  value={formData.parking_location}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                />
              </div>
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
