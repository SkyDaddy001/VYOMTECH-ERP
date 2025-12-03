'use client'

import React, { useState } from 'react'

interface LedgerEntry {
  id: string
  date: string
  description: string
  account: string
  debit: number
  credit: number
  balance: number
  reference?: string
}

interface LedgerBookProps {
  title: string
  entries: LedgerEntry[]
  accountName: string
  openingBalance: number
  onAddEntry?: (entry: LedgerEntry) => void
  editable?: boolean
}

export default function LedgerBook({
  title,
  entries,
  accountName,
  openingBalance,
  onAddEntry,
  editable = false
}: LedgerBookProps) {
  const [balance, setBalance] = useState(openingBalance)

  // Calculate running balance
  const entriesWithBalance = entries.map((entry, idx) => {
    const newBalance = balance + (entry.debit - entry.credit)
    return {
      ...entry,
      balance: newBalance
    }
  })

  return (
    <div className="font-serif bg-amber-50 p-8 max-w-4xl mx-auto border-4 border-yellow-900">
      {/* Header - Traditional Ledger Style */}
      <div className="text-center mb-2 pb-2 border-b-2 border-yellow-900">
        <h1 className="text-xl font-bold text-yellow-900">{title}</h1>
        <p className="text-sm text-gray-700 mt-1">LEDGER ACCOUNT</p>
      </div>

      {/* Account Name Header */}
      <div className="grid grid-cols-2 gap-4 mb-4 text-sm">
        <div className="border-b border-gray-400 pb-1">
          <span className="font-semibold text-yellow-900">Account Name:</span>
          <p className="text-gray-800">{accountName}</p>
        </div>
        <div className="border-b border-gray-400 pb-1">
          <span className="font-semibold text-yellow-900">Opening Balance:</span>
          <p className="text-gray-800 font-mono">₹ {openingBalance.toFixed(2)}</p>
        </div>
      </div>

      {/* Ledger Table - Traditional Lined Format */}
      <div className="mb-4">
        {/* Header Row */}
        <div className="grid grid-cols-12 gap-1 mb-1 font-bold text-xs text-yellow-900 bg-yellow-100">
          <div className="col-span-2 border-b border-yellow-900 px-2 py-1">Date</div>
          <div className="col-span-4 border-b border-yellow-900 px-2 py-1">Description</div>
          <div className="col-span-2 border-b border-yellow-900 px-2 py-1 text-right">Debit (₹)</div>
          <div className="col-span-2 border-b border-yellow-900 px-2 py-1 text-right">Credit (₹)</div>
          <div className="col-span-2 border-b border-yellow-900 px-2 py-1 text-right">Balance (₹)</div>
        </div>

        {/* Entries - With alternating light lines */}
        {entriesWithBalance.map((entry, idx) => (
          <div
            key={entry.id}
            className={`grid grid-cols-12 gap-1 border-b border-gray-300 py-1 text-xs ${
              idx % 2 === 0 ? 'bg-white' : 'bg-amber-50'
            }`}
          >
            <div className="col-span-2 px-2 py-1 text-gray-700 font-mono">
              {new Date(entry.date).toLocaleDateString('en-IN', {
                day: '2-digit',
                month: '2-digit',
                year: '2-digit'
              })}
            </div>
            <div className="col-span-4 px-2 py-1 text-gray-800 truncate">
              {entry.description}
            </div>
            <div className="col-span-2 px-2 py-1 text-right text-gray-800 font-mono">
              {entry.debit > 0 ? entry.debit.toFixed(2) : '-'}
            </div>
            <div className="col-span-2 px-2 py-1 text-right text-gray-800 font-mono">
              {entry.credit > 0 ? entry.credit.toFixed(2) : '-'}
            </div>
            <div className="col-span-2 px-2 py-1 text-right text-yellow-900 font-semibold font-mono">
              {entry.balance.toFixed(2)}
            </div>
          </div>
        ))}

        {/* Empty lines for writing more entries */}
        {entries.length < 15 && (
          <>
            {Array.from({ length: Math.max(0, 15 - entries.length) }).map((_, idx) => (
              <div
                key={`empty-${idx}`}
                className={`grid grid-cols-12 gap-1 border-b border-gray-300 py-3 ${
                  idx % 2 === 0 ? 'bg-white' : 'bg-amber-50'
                }`}
              />
            ))}
          </>
        )}
      </div>

      {/* Footer Summary */}
      <div className="border-t-2 border-yellow-900 pt-2 mt-4">
        <div className="grid grid-cols-3 gap-4 text-sm font-semibold text-yellow-900">
          <div className="text-center">
            <p>Total Debits</p>
            <p className="font-mono text-lg">₹ {entriesWithBalance.reduce((sum, e) => sum + e.debit, 0).toFixed(2)}</p>
          </div>
          <div className="text-center">
            <p>Total Credits</p>
            <p className="font-mono text-lg">₹ {entriesWithBalance.reduce((sum, e) => sum + e.credit, 0).toFixed(2)}</p>
          </div>
          <div className="text-center border-t-2 border-yellow-900 pt-1">
            <p>Closing Balance</p>
            <p className="font-mono text-lg text-green-700">
              ₹ {entriesWithBalance.length > 0 ? entriesWithBalance[entriesWithBalance.length - 1].balance.toFixed(2) : openingBalance.toFixed(2)}
            </p>
          </div>
        </div>
      </div>

      {/* Footer Notes */}
      <div className="text-xs text-gray-600 mt-4 pt-4 border-t border-gray-400">
        <p>Entered By: _________________ | Verified By: _________________ | Date: _________________</p>
      </div>
    </div>
  )
}
