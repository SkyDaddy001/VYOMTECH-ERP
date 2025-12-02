'use client'

import { useState } from 'react'
import { CallRecord } from '@/types/presales'
import CallRecordForm from './CallRecordForm'

interface CallRecordListProps {
  records?: CallRecord[]
  onCreateNew?: () => void
}

export default function CallRecordList({ records = [], onCreateNew }: CallRecordListProps) {
  const [selectedRecord, setSelectedRecord] = useState<CallRecord | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [filterOutcome, setFilterOutcome] = useState<string>('')

  const filteredRecords = filterOutcome
    ? records.filter(r => r.call_outcome === filterOutcome)
    : records

  const getOutcomeColor = (outcome: string) => {
    switch (outcome) {
      case 'interested':
        return 'bg-green-100 text-green-800'
      case 'not_interested':
        return 'bg-red-100 text-red-800'
      case 'maybe_later':
        return 'bg-yellow-100 text-yellow-800'
      case 'follow_up_needed':
        return 'bg-blue-100 text-blue-800'
      case 'converted':
        return 'bg-purple-100 text-purple-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString('en-IN', { month: 'short', day: 'numeric', year: 'numeric' })
  }

  if (showForm) {
    return (
      <CallRecordForm
        onSubmit={async (data) => {
          // Handle form submission
          console.log('New call record:', data)
          setShowForm(false)
        }}
        onCancel={() => setShowForm(false)}
      />
    )
  }

  if (selectedRecord) {
    return (
      <div className="bg-white rounded-lg border border-gray-200 p-6">
        <button
          onClick={() => setSelectedRecord(null)}
          className="mb-4 text-blue-600 hover:text-blue-700 font-medium"
        >
          ← Back to List
        </button>
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm text-gray-600">Customer Name</label>
              <p className="text-lg font-semibold">{selectedRecord.customer_name}</p>
            </div>
            <div>
              <label className="text-sm text-gray-600">Phone</label>
              <p className="text-lg font-semibold">{selectedRecord.customer_phone}</p>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm text-gray-600">Call Date & Time</label>
              <p className="text-lg font-semibold">
                {formatDate(selectedRecord.call_date)} at {selectedRecord.call_time}
              </p>
            </div>
            <div>
              <label className="text-sm text-gray-600">Call Outcome</label>
              <div className="mt-1">
                <span className={`inline-block px-3 py-1 rounded-full text-sm font-medium ${getOutcomeColor(selectedRecord.call_outcome)}`}>
                  {selectedRecord.call_outcome.replace(/_/g, ' ')}
                </span>
              </div>
            </div>
          </div>

          <div>
            <label className="text-sm text-gray-600 block mb-2">Property Preferences</label>
            <div className="flex flex-wrap gap-2">
              {selectedRecord.property_preference.map(pref => (
                <span key={pref} className="px-2 py-1 bg-green-100 text-green-800 rounded text-sm">
                  {pref}
                </span>
              ))}
            </div>
          </div>

          <div>
            <label className="text-sm text-gray-600 block mb-2">Budget Range</label>
            <p className="text-sm">{selectedRecord.budget_specified ? selectedRecord.budget_range : 'Not Specified'}</p>
          </div>

          <div>
            <label className="text-sm text-gray-600 block mb-2">Interested Projects</label>
            <div className="flex flex-wrap gap-2">
              {selectedRecord.interested_projects.map(proj => (
                <span key={proj} className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-sm">
                  {proj}
                </span>
              ))}
            </div>
          </div>

          <div>
            <label className="text-sm text-gray-600">Call Summary</label>
            <p className="text-sm text-gray-700 mt-1 bg-gray-50 p-3 rounded">
              {selectedRecord.call_summary}
            </p>
          </div>

          {selectedRecord.follow_up_required && (
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-3">
              <p className="text-sm font-medium text-yellow-900">
                Follow-up Required: {selectedRecord.follow_up_date ? formatDate(selectedRecord.follow_up_date) : 'Date not set'}
              </p>
            </div>
          )}
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-semibold text-gray-800">Call Records</h2>
        <button
          onClick={() => setShowForm(true)}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium transition-colors"
        >
          + New Call Record
        </button>
      </div>

      {/* Filters */}
      <div className="mb-6">
        <label className="text-sm text-gray-600 block mb-2">Filter by Outcome</label>
        <div className="flex gap-2 flex-wrap">
          <button
            onClick={() => setFilterOutcome('')}
            className={`px-3 py-1 rounded-lg text-sm font-medium transition-colors ${
              filterOutcome === ''
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            All
          </button>
          {['interested', 'not_interested', 'maybe_later', 'follow_up_needed', 'converted'].map(outcome => (
            <button
              key={outcome}
              onClick={() => setFilterOutcome(outcome)}
              className={`px-3 py-1 rounded-lg text-sm font-medium transition-colors ${
                filterOutcome === outcome
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              {outcome.replace(/_/g, ' ')}
            </button>
          ))}
        </div>
      </div>

      {/* Records Table */}
      {filteredRecords.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-600 mb-4">No call records found</p>
          <button
            onClick={() => setShowForm(true)}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium"
          >
            Create First Call Record
          </button>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-300 bg-gray-50">
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Date & Time</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Customer</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Phone</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Properties</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Budget</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Outcome</th>
                <th className="text-left px-4 py-3 text-sm font-semibold text-gray-700">Action</th>
              </tr>
            </thead>
            <tbody>
              {filteredRecords.map(record => (
                <tr key={record.id} className="border-b border-gray-200 hover:bg-gray-50 transition-colors">
                  <td className="px-4 py-3 text-sm">
                    <div className="font-medium">{formatDate(record.call_date)}</div>
                    <div className="text-xs text-gray-600">{record.call_time}</div>
                  </td>
                  <td className="px-4 py-3 text-sm font-medium">{record.customer_name}</td>
                  <td className="px-4 py-3 text-sm">{record.customer_phone}</td>
                  <td className="px-4 py-3 text-sm">
                    <div className="flex flex-wrap gap-1">
                      {record.property_preference.slice(0, 2).map(pref => (
                        <span key={pref} className="px-2 py-0.5 bg-green-100 text-green-800 rounded text-xs">
                          {pref}
                        </span>
                      ))}
                      {record.property_preference.length > 2 && (
                        <span className="px-2 py-0.5 bg-gray-100 text-gray-800 rounded text-xs">
                          +{record.property_preference.length - 2}
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="px-4 py-3 text-sm">
                    {record.budget_specified ? record.budget_range.replace(/_/g, ' ') : 'Not Specified'}
                  </td>
                  <td className="px-4 py-3 text-sm">
                    <span className={`inline-block px-2 py-1 rounded text-xs font-medium ${getOutcomeColor(record.call_outcome)}`}>
                      {record.call_outcome.replace(/_/g, ' ')}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm">
                    <button
                      onClick={() => setSelectedRecord(record)}
                      className="text-blue-600 hover:text-blue-700 font-medium"
                    >
                      View Details
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Summary Stats */}
      <div className="mt-6 grid grid-cols-4 gap-4 pt-6 border-t border-gray-200">
        <div className="bg-blue-50 rounded-lg p-4">
          <p className="text-sm text-blue-600 mb-1">Total Calls</p>
          <p className="text-2xl font-semibold text-blue-900">{records.length}</p>
        </div>
        <div className="bg-green-50 rounded-lg p-4">
          <p className="text-sm text-green-600 mb-1">Interested</p>
          <p className="text-2xl font-semibold text-green-900">
            {records.filter(r => r.call_outcome === 'interested').length}
          </p>
        </div>
        <div className="bg-purple-50 rounded-lg p-4">
          <p className="text-sm text-purple-600 mb-1">Converted</p>
          <p className="text-2xl font-semibold text-purple-900">
            {records.filter(r => r.call_outcome === 'converted').length}
          </p>
        </div>
        <div className="bg-orange-50 rounded-lg p-4">
          <p className="text-sm text-orange-600 mb-1">Follow-ups Pending</p>
          <p className="text-2xl font-semibold text-orange-900">
            {records.filter(r => r.follow_up_required && r.call_outcome !== 'converted').length}
          </p>
        </div>
      </div>
    </div>
  )
}
