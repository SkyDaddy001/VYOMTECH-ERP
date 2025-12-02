'use client'

import { useState } from 'react'
import { Vendor } from '@/types/purchase'

interface VendorManagementFormProps {
  vendor?: Vendor | null
  onSubmit: (data: Partial<Vendor>) => Promise<void>
  onCancel: () => void
}

export default function VendorManagementForm({ vendor, onSubmit, onCancel }: VendorManagementFormProps) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Partial<Vendor>>(
    vendor || {
      vendor_code: '',
      vendor_name: '',
      contact_person: '',
      email: '',
      phone: '',
      alternate_phone: '',
      address: '',
      city: '',
      state: '',
      zip_code: '',
      gst_number: '',
      pan_number: '',
      bank_account_number: '',
      bank_ifsc_code: '',
      bank_name: '',
      account_holder_name: '',
      payment_terms: 'net_30',
      payment_method: 'bank_transfer',
      status: 'active',
      rating: 0,
      notes: '',
    }
  )

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

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      <div>
        <h2 className="text-xl font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <span className="text-2xl">🏢</span> {vendor ? 'Edit Vendor' : 'Add Vendor'}
        </h2>
      </div>

      {/* Basic Information */}
      <div className="bg-blue-50 rounded-lg p-4 border border-blue-200">
        <h3 className="text-sm font-semibold text-blue-900 mb-4">Basic Information</h3>
        <div className="grid grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Vendor Code *</label>
            <input
              type="text"
              required
              value={formData.vendor_code || ''}
              onChange={(e) => handleChange('vendor_code', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="VND-001"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Vendor Name *</label>
            <input
              type="text"
              required
              value={formData.vendor_name || ''}
              onChange={(e) => handleChange('vendor_name', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Company name"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
            <select
              value={formData.status || 'active'}
              onChange={(e) => handleChange('status', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
              <option value="suspended">Suspended</option>
            </select>
          </div>
        </div>
      </div>

      {/* Contact Information */}
      <div className="bg-green-50 rounded-lg p-4 border border-green-200">
        <h3 className="text-sm font-semibold text-green-900 mb-4">Contact Information</h3>
        <div className="grid grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Contact Person *</label>
            <input
              type="text"
              required
              value={formData.contact_person || ''}
              onChange={(e) => handleChange('contact_person', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              placeholder="Name"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Email *</label>
            <input
              type="email"
              required
              value={formData.email || ''}
              onChange={(e) => handleChange('email', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              placeholder="email@company.com"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Phone *</label>
            <input
              type="tel"
              required
              value={formData.phone || ''}
              onChange={(e) => handleChange('phone', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              placeholder="+91-98765-43210"
            />
          </div>
        </div>
        <div className="grid grid-cols-2 gap-4 mt-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Alternate Phone</label>
            <input
              type="tel"
              value={formData.alternate_phone || ''}
              onChange={(e) => handleChange('alternate_phone', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              placeholder="+91-98765-43211"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Rating (0-5)</label>
            <input
              type="number"
              value={formData.rating || 0}
              onChange={(e) => handleChange('rating', parseFloat(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              min="0"
              max="5"
              step="0.5"
            />
          </div>
        </div>
      </div>

      {/* Address Information */}
      <div className="bg-yellow-50 rounded-lg p-4 border border-yellow-200">
        <h3 className="text-sm font-semibold text-yellow-900 mb-4">Address</h3>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Address *</label>
          <textarea
            required
            value={formData.address || ''}
            onChange={(e) => handleChange('address', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500"
            placeholder="Street address"
            rows={2}
          />
        </div>
        <div className="grid grid-cols-3 gap-4 mt-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">City</label>
            <input
              type="text"
              value={formData.city || ''}
              onChange={(e) => handleChange('city', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500"
              placeholder="City"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">State</label>
            <input
              type="text"
              value={formData.state || ''}
              onChange={(e) => handleChange('state', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500"
              placeholder="State"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Zip Code</label>
            <input
              type="text"
              value={formData.zip_code || ''}
              onChange={(e) => handleChange('zip_code', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-yellow-500"
              placeholder="Zip code"
            />
          </div>
        </div>
      </div>

      {/* Tax & Compliance */}
      <div className="bg-purple-50 rounded-lg p-4 border border-purple-200">
        <h3 className="text-sm font-semibold text-purple-900 mb-4">Tax & Compliance</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">GST Number</label>
            <input
              type="text"
              value={formData.gst_number || ''}
              onChange={(e) => handleChange('gst_number', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
              placeholder="GST ID"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">PAN Number</label>
            <input
              type="text"
              value={formData.pan_number || ''}
              onChange={(e) => handleChange('pan_number', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
              placeholder="PAN ID"
            />
          </div>
        </div>
      </div>

      {/* Bank Information */}
      <div className="bg-indigo-50 rounded-lg p-4 border border-indigo-200">
        <h3 className="text-sm font-semibold text-indigo-900 mb-4">Bank Information</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Bank Name</label>
            <input
              type="text"
              value={formData.bank_name || ''}
              onChange={(e) => handleChange('bank_name', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Bank name"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Account Holder Name</label>
            <input
              type="text"
              value={formData.account_holder_name || ''}
              onChange={(e) => handleChange('account_holder_name', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Account holder"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Account Number</label>
            <input
              type="text"
              value={formData.bank_account_number || ''}
              onChange={(e) => handleChange('bank_account_number', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Account number"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">IFSC Code</label>
            <input
              type="text"
              value={formData.bank_ifsc_code || ''}
              onChange={(e) => handleChange('bank_ifsc_code', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="IFSC code"
            />
          </div>
        </div>
      </div>

      {/* Payment Terms */}
      <div className="bg-pink-50 rounded-lg p-4 border border-pink-200">
        <h3 className="text-sm font-semibold text-pink-900 mb-4">Payment Terms</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Payment Terms</label>
            <select
              value={formData.payment_terms || 'net_30'}
              onChange={(e) => handleChange('payment_terms', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-pink-500"
            >
              <option value="net_15">Net 15</option>
              <option value="net_30">Net 30</option>
              <option value="net_45">Net 45</option>
              <option value="net_60">Net 60</option>
              <option value="advance">Advance</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Payment Method</label>
            <select
              value={formData.payment_method || 'bank_transfer'}
              onChange={(e) => handleChange('payment_method', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-pink-500"
            >
              <option value="bank_transfer">Bank Transfer</option>
              <option value="cheque">Cheque</option>
              <option value="credit_card">Credit Card</option>
              <option value="cash">Cash</option>
            </select>
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
          placeholder="Additional notes about the vendor..."
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
          {loading ? 'Saving...' : vendor ? 'Update Vendor' : 'Add Vendor'}
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
