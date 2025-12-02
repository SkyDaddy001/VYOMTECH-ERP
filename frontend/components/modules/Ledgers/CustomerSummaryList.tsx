'use client'

import React, { useState } from 'react'
import { LedgerSummary } from '@/types/ledgers'

interface CustomerSummaryListProps {
  summaries: LedgerSummary[]
  loading: boolean
}

export default function CustomerSummaryList({ summaries, loading }: CustomerSummaryListProps) {
  const [filterOutstanding, setFilterOutstanding] = useState(false)

  const filteredSummaries = filterOutstanding
    ? summaries.filter(s => s.closing_balance < 0)
    : summaries

  if (loading) {
    return <div className="text-center py-8 text-gray-600">Loading summaries...</div>
  }

  return (
    <div className="space-y-4">
      <label className="flex items-center gap-2 cursor-pointer">
        <input
          type="checkbox"
          checked={filterOutstanding}
          onChange={(e) => setFilterOutstanding(e.target.checked)}
          className="w-4 h-4 rounded border-gray-300"
        />
        <span className="text-sm font-medium text-gray-700">Show Outstanding Only</span>
      </label>

      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b-2 border-gray-300 bg-gray-50">
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Customer</th>
              <th className="px-4 py-3 text-right font-semibold text-gray-700">Opening</th>
              <th className="px-4 py-3 text-right font-semibold text-gray-700">Debit</th>
              <th className="px-4 py-3 text-right font-semibold text-gray-700">Credit</th>
              <th className="px-4 py-3 text-right font-semibold text-gray-700">Closing</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Status</th>
            </tr>
          </thead>
          <tbody>
            {filteredSummaries.map((summary) => (
              <tr key={summary.customer_id} className="border-b border-gray-200 hover:bg-gray-50">
                <td className="px-4 py-3 font-medium text-gray-900">{summary.customer_name}</td>
                <td className="px-4 py-3 text-right">₹{(summary.opening_balance / 100000).toFixed(2)}L</td>
                <td className="px-4 py-3 text-right text-red-600 font-semibold">₹{(summary.total_debit / 100000).toFixed(2)}L</td>
                <td className="px-4 py-3 text-right text-green-600 font-semibold">₹{(summary.total_credit / 100000).toFixed(2)}L</td>
                <td className={`px-4 py-3 text-right font-bold ${summary.closing_balance < 0 ? 'text-orange-600' : 'text-blue-600'}`}>
                  ₹{(summary.closing_balance / 100000).toFixed(2)}L
                </td>
                <td className="px-4 py-3">
                  <span className={`px-2 py-1 text-xs font-medium rounded ${
                    summary.closing_balance < 0
                      ? 'bg-orange-100 text-orange-800'
                      : 'bg-green-100 text-green-800'
                  }`}>
                    {summary.closing_balance < 0 ? 'Outstanding' : 'Settled'}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {filteredSummaries.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No customer summaries found</p>
        </div>
      )}
    </div>
  )
}
