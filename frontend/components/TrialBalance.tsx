'use client'

import React, { useState } from 'react'

interface TrialBalanceEntry {
  accountName: string
  accountCode: string
  debitBalance: number
  creditBalance: number
  accountType: 'Asset' | 'Liability' | 'Equity' | 'Income' | 'Expense'
}

interface TrialBalanceProps {
  date: string
  entries: TrialBalanceEntry[]
  preparedBy?: string
  approvedBy?: string
}

export default function TrialBalance({
  date,
  entries,
  preparedBy,
  approvedBy
}: TrialBalanceProps) {
  const totalDebits = entries.reduce((sum, e) => sum + e.debitBalance, 0)
  const totalCredits = entries.reduce((sum, e) => sum + e.creditBalance, 0)
  const isBalanced = Math.abs(totalDebits - totalCredits) < 0.01

  return (
    <div className="font-serif bg-white p-8 border-4 border-black max-w-5xl mx-auto">
      {/* Header */}
      <div className="text-center border-b-4 border-black pb-3 mb-4">
        <h1 className="text-2xl font-bold text-black mb-1">TRIAL BALANCE</h1>
        <p className="text-sm text-gray-700">As on {new Date(date).toLocaleDateString('en-IN')}</p>
      </div>

      {/* Prepared By Section */}
      <div className="grid grid-cols-3 gap-4 mb-4 text-xs border-b border-gray-400 pb-2">
        <div>
          <span className="font-semibold">Prepared By:</span>
          <p className="mt-1">{preparedBy || '_______________'}</p>
        </div>
        <div className="text-center">
          <span className="font-semibold">Date: {new Date().toLocaleDateString('en-IN')}</span>
        </div>
        <div className="text-right">
          <span className="font-semibold">Approved By:</span>
          <p className="mt-1">{approvedBy || '_______________'}</p>
        </div>
      </div>

      {/* Trial Balance Table */}
      <div className="mb-4">
        {/* Header Row */}
        <div className="grid grid-cols-12 gap-1 font-bold text-xs bg-gray-300 border-2 border-black">
          <div className="col-span-1 px-1 py-2 border-r border-black text-center">Sr.</div>
          <div className="col-span-2 px-1 py-2 border-r border-black">Account Code</div>
          <div className="col-span-5 px-1 py-2 border-r border-black">Account Name</div>
          <div className="col-span-2 px-1 py-2 border-r border-black text-right">Debit (₹)</div>
          <div className="col-span-2 px-1 py-2 text-right">Credit (₹)</div>
        </div>

        {/* Entries grouped by type */}
        {['Asset', 'Liability', 'Equity', 'Income', 'Expense'].map((type, typeIdx) => {
          const typeEntries = entries.filter(e => e.accountType === type)
          if (typeEntries.length === 0) return null

          return (
            <div key={type}>
              {/* Category Header */}
              <div className="grid grid-cols-12 gap-1 bg-gray-100 border-b border-gray-400 py-1 mt-1">
                <div className="col-span-12 px-2 font-semibold text-sm text-gray-800">
                  {type.toUpperCase()}S
                </div>
              </div>

              {/* Category Entries */}
              {typeEntries.map((entry, idx) => (
                <div
                  key={`${type}-${idx}`}
                  className={`grid grid-cols-12 gap-1 border-b border-gray-200 py-1 ${
                    idx % 2 === 0 ? 'bg-white' : 'bg-gray-50'
                  }`}
                >
                  <div className="col-span-1 px-1 py-1 border-r border-gray-300 text-center text-xs text-gray-600">
                    {idx + 1}
                  </div>
                  <div className="col-span-2 px-1 py-1 border-r border-gray-300 font-mono text-xs text-gray-700">
                    {entry.accountCode}
                  </div>
                  <div className="col-span-5 px-1 py-1 border-r border-gray-300 text-xs text-gray-800">
                    {entry.accountName}
                  </div>
                  <div className="col-span-2 px-1 py-1 border-r border-gray-300 text-right font-mono text-xs">
                    {entry.debitBalance > 0 ? entry.debitBalance.toFixed(2) : '-'}
                  </div>
                  <div className="col-span-2 px-1 py-1 text-right font-mono text-xs">
                    {entry.creditBalance > 0 ? entry.creditBalance.toFixed(2) : '-'}
                  </div>
                </div>
              ))}
            </div>
          )
        })}

        {/* Total Row */}
        <div className="grid grid-cols-12 gap-1 border-t-4 border-b-4 border-black font-bold text-sm py-2 bg-gray-100">
          <div className="col-span-8 px-1 py-1 text-right">TOTAL</div>
          <div className="col-span-2 px-1 py-1 border-r border-black text-right font-mono">
            {totalDebits.toFixed(2)}
          </div>
          <div className="col-span-2 px-1 py-1 text-right font-mono">
            {totalCredits.toFixed(2)}
          </div>
        </div>
      </div>

      {/* Balance Verification */}
      <div className="grid grid-cols-2 gap-4 mb-4 mt-4 text-sm">
        <div>
          <p className="font-semibold text-gray-700">Total Debits:</p>
          <p className="font-mono font-bold text-lg">₹ {totalDebits.toFixed(2)}</p>
        </div>
        <div>
          <p className="font-semibold text-gray-700">Total Credits:</p>
          <p className="font-mono font-bold text-lg">₹ {totalCredits.toFixed(2)}</p>
        </div>
      </div>

      {/* Balance Status */}
      <div className={`text-center py-3 border-2 ${isBalanced ? 'border-green-700 bg-green-50' : 'border-red-700 bg-red-50'}`}>
        <p className={`text-lg font-bold ${isBalanced ? 'text-green-700' : 'text-red-700'}`}>
          {isBalanced ? '✓ TRIAL BALANCE BALANCED' : '✗ TRIAL BALANCE NOT BALANCED'}
        </p>
        {!isBalanced && (
          <p className="text-xs text-red-600 mt-1">
            Difference: ₹ {Math.abs(totalDebits - totalCredits).toFixed(2)}
          </p>
        )}
      </div>

      {/* Accounting Equation */}
      <div className="mt-6 pt-4 border-t-2 border-black">
        <p className="text-xs font-semibold text-gray-700 mb-2">ACCOUNTING EQUATION VERIFICATION:</p>
        <div className="grid grid-cols-3 gap-4 text-sm">
          <div className="text-center border border-gray-400 p-2">
            <p className="text-xs text-gray-600">Assets</p>
            <p className="font-mono font-bold text-lg">
              ₹ {entries
                .filter(e => e.accountType === 'Asset')
                .reduce((sum, e) => sum + e.debitBalance, 0)
                .toFixed(2)}
            </p>
          </div>
          <div className="text-center py-2">
            <p className="text-2xl font-bold">=</p>
          </div>
          <div>
            <div className="text-center border border-gray-400 p-2 mb-2">
              <p className="text-xs text-gray-600">Liabilities + Equity</p>
              <p className="font-mono font-bold text-lg">
                ₹ {(
                  entries
                    .filter(e => e.accountType === 'Liability')
                    .reduce((sum, e) => sum + e.creditBalance, 0) +
                  entries
                    .filter(e => e.accountType === 'Equity')
                    .reduce((sum, e) => sum + e.creditBalance, 0)
                ).toFixed(2)}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="mt-6 pt-4 border-t-2 border-black">
        <div className="grid grid-cols-3 gap-6 text-center">
          <div>
            <p className="text-xs font-semibold text-gray-700 mb-6">Prepared By</p>
            <p className="border-t-2 border-black pt-2 h-8"></p>
          </div>
          <div>
            <p className="text-xs font-semibold text-gray-700 mb-6">Verified By</p>
            <p className="border-t-2 border-black pt-2 h-8"></p>
          </div>
          <div>
            <p className="text-xs font-semibold text-gray-700 mb-6">Approved By</p>
            <p className="border-t-2 border-black pt-2 h-8"></p>
          </div>
        </div>
      </div>
    </div>
  )
}
