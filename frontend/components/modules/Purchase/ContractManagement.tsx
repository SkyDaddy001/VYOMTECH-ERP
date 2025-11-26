'use client'

import { useState, useEffect } from 'react'
import axios from 'axios'
import toast from 'react-hot-toast'
import { formatDateToDDMMMYYYY, formatDateToInput } from '@/lib/dateFormat'

interface Contract {
  id: string
  contract_number: string
  vendor_id: string
  vendor_name?: string
  contract_type: string
  boq_id?: string
  start_date: string
  end_date: string
  contract_value: number
  status: string
  description?: string
}

const CONTRACT_TYPES = [
  { value: 'material', label: 'Material Contract' },
  { value: 'labour', label: 'Labour Contract' },
  { value: 'service', label: 'Service Contract' },
  { value: 'hybrid', label: 'Hybrid (Material + Labour)' },
  { value: 'hybrid_service', label: 'Hybrid (Material + Service)' },
]

export default function ContractManagement() {
  const [contracts, setContracts] = useState<Contract[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [vendors, setVendors] = useState<{ id: string; name: string }[]>([])
  const [boqs, setBoqs] = useState<{ id: string; description: string }[]>([])
  const [formData, setFormData] = useState({
    vendor_id: '',
    contract_type: 'material',
    boq_id: '',
    start_date: new Date().toISOString().split('T')[0],
    end_date: '',
    contract_value: 0,
    description: '',
  })

  useEffect(() => {
    fetchContracts()
    fetchVendors()
    fetchBOQs()
  }, [])

  const fetchContracts = async () => {
    try {
      const response = await axios.get('/api/v1/purchase/contracts')
      setContracts(response.data || [])
    } catch (error) {
      toast.error('Failed to fetch contracts')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const fetchVendors = async () => {
    try {
      const response = await axios.get('/api/v1/purchase/vendors')
      setVendors(response.data || [])
    } catch (error) {
      console.error('Failed to fetch vendors:', error)
    }
  }

  const fetchBOQs = async () => {
    try {
      const response = await axios.get('/api/v1/construction/boq')
      setBoqs(response.data || [])
    } catch (error) {
      console.error('Failed to fetch BOQs:', error)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await axios.post('/api/v1/purchase/contracts', formData)
      toast.success('Contract created successfully')
      setShowForm(false)
      setFormData({
        vendor_id: '',
        contract_type: 'material',
        boq_id: '',
        start_date: new Date().toISOString().split('T')[0],
        end_date: '',
        contract_value: 0,
        description: '',
      })
      fetchContracts()
    } catch (error) {
      toast.error('Failed to create contract')
      console.error(error)
    }
  }

  const getContractTypeLabel = (type: string) => {
    const contract = CONTRACT_TYPES.find((c) => c.value === type)
    return contract ? contract.label : type
  }

  const getStatusColor = (status: string) => {
    const colors: { [key: string]: string } = {
      active: 'bg-green-100 text-green-800',
      pending: 'bg-yellow-100 text-yellow-800',
      completed: 'bg-blue-100 text-blue-800',
      terminated: 'bg-red-100 text-red-800',
      on_hold: 'bg-orange-100 text-orange-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  if (loading) {
    return <div className="p-6 text-center text-gray-500">Loading contracts...</div>
  }

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Contract Management</h2>
        <button
          onClick={() => setShowForm(true)}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition"
        >
          + New Contract
        </button>
      </div>

      {/* Info Box */}
      <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
        <p className="text-sm text-purple-800">
          <strong>Contract Types:</strong>
        </p>
        <ul className="text-sm text-purple-700 mt-2 space-y-1">
          <li>• <strong>Material:</strong> Supply of materials/goods</li>
          <li>• <strong>Labour:</strong> Contract labour services</li>
          <li>• <strong>Service:</strong> Professional services</li>
          <li>• <strong>Hybrid:</strong> Combination contracts (Material + Labour, Material + Service)</li>
        </ul>
      </div>

      {/* Form Modal */}
      {showForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-8 max-w-md w-full max-h-96 overflow-y-auto">
            <h3 className="text-xl font-bold mb-4">Create Contract Against BOQ</h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Vendor</label>
                <select
                  value={formData.vendor_id}
                  onChange={(e) => setFormData({ ...formData, vendor_id: e.target.value })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                >
                  <option value="">Select Vendor</option>
                  {vendors.map((vendor) => (
                    <option key={vendor.id} value={vendor.id}>
                      {vendor.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Contract Type</label>
                <select
                  value={formData.contract_type}
                  onChange={(e) => setFormData({ ...formData, contract_type: e.target.value })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                >
                  {CONTRACT_TYPES.map((type) => (
                    <option key={type.value} value={type.value}>
                      {type.label}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Bill of Quantities (BOQ)</label>
                <select
                  value={formData.boq_id}
                  onChange={(e) => setFormData({ ...formData, boq_id: e.target.value })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Select BOQ (Optional)</option>
                  {boqs.map((boq) => (
                    <option key={boq.id} value={boq.id}>
                      {boq.description}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Start Date</label>
                <input
                  type="date"
                  value={formData.start_date}
                  onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">End Date</label>
                <input
                  type="date"
                  value={formData.end_date}
                  onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Contract Value</label>
                <input
                  type="number"
                  value={formData.contract_value}
                  onChange={(e) => setFormData({ ...formData, contract_value: parseFloat(e.target.value) })}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                <textarea
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  rows={2}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Contract details..."
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowForm(false)}
                  className="flex-1 px-4 py-2 border rounded-lg hover:bg-gray-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                >
                  Create
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Contracts Table */}
      <div className="overflow-x-auto">
        <table className="w-full border-collapse">
          <thead>
            <tr className="bg-gray-100 border-b">
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Contract #</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Vendor</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Type</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Period</th>
              <th className="px-4 py-3 text-right font-semibold text-gray-700">Value</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Status</th>
              <th className="px-4 py-3 text-center font-semibold text-gray-700">Actions</th>
            </tr>
          </thead>
          <tbody>
            {contracts.length === 0 ? (
              <tr>
                <td colSpan={7} className="px-4 py-3 text-center text-gray-500">
                  No contracts found
                </td>
              </tr>
            ) : (
              contracts.map((contract) => (
                <tr key={contract.id} className="border-b hover:bg-gray-50">
                  <td className="px-4 py-3 font-mono font-semibold">{contract.contract_number}</td>
                  <td className="px-4 py-3">{contract.vendor_name || 'N/A'}</td>
                  <td className="px-4 py-3 text-sm">
                    <span className="inline-block bg-gray-100 px-2 py-1 rounded">
                      {getContractTypeLabel(contract.contract_type)}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm">
                    {formatDateToDDMMMYYYY(contract.start_date)} -{' '}
                    {formatDateToDDMMMYYYY(contract.end_date)}
                  </td>
                  <td className="px-4 py-3 text-right font-semibold">
                    ${contract.contract_value.toLocaleString('en-US', { minimumFractionDigits: 2 })}
                  </td>
                  <td className="px-4 py-3">
                    <span
                      className={`inline-block px-3 py-1 rounded text-sm font-medium ${getStatusColor(
                        contract.status
                      )}`}
                    >
                      {contract.status.charAt(0).toUpperCase() + contract.status.slice(1)}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-center space-x-2">
                    <button className="px-3 py-1 bg-blue-100 text-blue-600 rounded hover:bg-blue-200 text-sm">
                      View
                    </button>
                    <button className="px-3 py-1 bg-purple-100 text-purple-600 rounded hover:bg-purple-200 text-sm">
                      Edit
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
