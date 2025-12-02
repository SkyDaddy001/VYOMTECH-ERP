'use client'

import React, { useState } from 'react'
import { Ledger } from '@/types/ledgers'

interface LedgerListProps {
  entries: Ledger[]
  loading: boolean
  onEdit: (entry: Ledger) => void
  onDelete: (entry: Ledger) => void
}

export default function LedgerList({ entries, loading, onEdit, onDelete }: LedgerListProps) {
  const [filterType, setFilterType] = useState<string>('all')

  const filteredEntries = filterType === 'all'
    ? entries
    : entries.filter(e => e.transaction_type === filterType)

  if (loading) {
    return <div className="text-center py-8 text-gray-600">Loading ledger entries...</div>
  }

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <select
          value={filterType}
          onChange={(e) => setFilterType(e.target.value)}
          className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="all">All Transactions</option>
          <option value="debit">Debits Only</option>
          <option value="credit">Credits Only</option>
        </select>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b-2 border-gray-300">
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Date</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Reference</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">Description</th>
              <th className="px-4 py-3 text-right font-semibold text-gray-700">Debit</th>
              <th className="px-4 py-3 text-right font-semibold text-gray-700">Credit</th>
              <th className="px-4 py-3 text-right font-semibold text-gray-700">Balance</th>
              <th className="px-4 py-3 text-center font-semibold text-gray-700">Action</th>
            </tr>
          </thead>
          <tbody>
            {filteredEntries.map((entry) => (
              <tr key={entry.id} className="border-b border-gray-200 hover:bg-gray-50">
                <td className="px-4 py-3 text-xs">{new Date(entry.transaction_date).toLocaleDateString()}</td>
                <td className="px-4 py-3 text-xs font-medium">{entry.reference_number}</td>
                <td className="px-4 py-3 text-sm">{entry.description}</td>
                <td className="px-4 py-3 text-right font-semibold text-red-600">
                  {entry.debit_amount > 0 ? `₹${(entry.debit_amount / 100000).toFixed(2)}L` : '-'}
                </td>
                <td className="px-4 py-3 text-right font-semibold text-green-600">
                  {entry.credit_amount > 0 ? `₹${(entry.credit_amount / 100000).toFixed(2)}L` : '-'}
                </td>
                <td className="px-4 py-3 text-right font-bold text-blue-600">
                  ₹{(entry.closing_balance / 100000).toFixed(2)}L
                </td>
                <td className="px-4 py-3 text-center flex gap-1 justify-center">
                  <button
                    onClick={() => onEdit(entry)}
                    className="px-2 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => {
                      if (confirm('Delete this entry?')) {
                        onDelete(entry)
                      }
                    }}
                    className="px-2 py-1 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {filteredEntries.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No ledger entries found</p>
        </div>
      )}
    </div>
  )
}
