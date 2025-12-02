'use client'

import { PropertyUnit } from '@/types/unit'

interface UnitListProps {
  units: PropertyUnit[]
  loading: boolean
  onEdit: (unit: PropertyUnit) => void
  onDelete: (unit: PropertyUnit) => void
  onViewCostSheet: (unit: PropertyUnit) => void
}

export default function UnitList({
  units,
  loading,
  onEdit,
  onDelete,
  onViewCostSheet,
}: UnitListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading units...</p>
        </div>
      </div>
    )
  }

  if (units.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No units found.</p>
      </div>
    )
  }

  const stats = {
    total: units.length,
    available: units.filter((u) => u.status === 'available').length,
    booked: units.filter((u) => u.status === 'booked').length,
    sold: units.filter((u) => u.status === 'sold').length,
    reserved: units.filter((u) => u.status === 'reserved').length,
  }

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{stats.total}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Available</p>
          <p className="text-2xl font-bold text-green-600 mt-1">{stats.available}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Booked</p>
          <p className="text-2xl font-bold text-yellow-600 mt-1">{stats.booked}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Sold</p>
          <p className="text-2xl font-bold text-blue-600 mt-1">{stats.sold}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Reserved</p>
          <p className="text-2xl font-bold text-purple-600 mt-1">{stats.reserved}</p>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow overflow-x-auto border border-gray-200">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Unit</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Floor</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Area</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {units.map((unit) => (
              <tr key={unit.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{unit.unit_number}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{unit.floor}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{unit.unit_type}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{unit.carpet_area} sq.ft</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 py-1 rounded-full text-xs font-medium ${
                      unit.status === 'available'
                        ? 'bg-green-100 text-green-800'
                        : unit.status === 'booked'
                        ? 'bg-yellow-100 text-yellow-800'
                        : unit.status === 'sold'
                        ? 'bg-blue-100 text-blue-800'
                        : 'bg-purple-100 text-purple-800'
                    }`}
                  >
                    {unit.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                  <button
                    onClick={() => onViewCostSheet(unit)}
                    className="text-blue-600 hover:text-blue-900 font-medium"
                  >
                    Cost
                  </button>
                  <button onClick={() => onEdit(unit)} className="text-green-600 hover:text-green-900 font-medium">
                    Edit
                  </button>
                  <button onClick={() => onDelete(unit)} className="text-red-600 hover:text-red-900 font-medium">
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
