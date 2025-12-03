'use client'

import React, { useState } from 'react'

interface VoucherLine {
  account: string
  description: string
  debit?: number
  credit?: number
  reference?: string
}

interface VoucherProps {
  voucherNo: string
  date: string
  voucherType: 'JV' | 'CV' | 'PV' | 'RV' // Journal, Credit, Payment, Receipt
  entries: VoucherLine[]
  approvedBy?: string
  createdBy?: string
  narration?: string
  editable?: boolean
}

const voucherTypeLabels = {
  'JV': 'JOURNAL VOUCHER',
  'CV': 'CREDIT VOUCHER',
  'PV': 'PAYMENT VOUCHER',
  'RV': 'RECEIPT VOUCHER'
}

export default function TraditionalVoucher({
  voucherNo,
  date,
  voucherType,
  entries,
  approvedBy,
  createdBy,
  narration,
  editable = false
}: VoucherProps) {
  const totalDebit = entries.reduce((sum, e) => sum + (e.debit || 0), 0)
  const totalCredit = entries.reduce((sum, e) => sum + (e.credit || 0), 0)
  const isBalanced = Math.abs(totalDebit - totalCredit) < 0.01

  return (
    <div className="font-serif bg-white p-6 border-4 border-black max-w-5xl mx-auto">
      {/* Header */}
      <div className="border-b-2 border-black pb-2 mb-3">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-lg font-bold text-black">{voucherTypeLabels[voucherType]}</h1>
            <p className="text-xs text-gray-600">Original</p>
          </div>
          <div className="text-right">
            <p className="font-mono font-bold text-sm">Voucher No: {voucherNo}</p>
            <p className="text-sm">Date: {new Date(date).toLocaleDateString('en-IN')}</p>
          </div>
        </div>
      </div>

      {/* Voucher Details */}
      <div className="grid grid-cols-3 gap-2 text-xs mb-3">
        <div className="border-b border-gray-400 pb-1">
          <span className="font-semibold">Prepared By:</span>
          <p className="h-5">{createdBy || '_______________'}</p>
        </div>
        <div className="border-b border-gray-400 pb-1">
          <span className="font-semibold">Approved By:</span>
          <p className="h-5">{approvedBy || '_______________'}</p>
        </div>
        <div className="border-b border-gray-400 pb-1">
          <span className="font-semibold">Date:</span>
          <p className="h-5">_______________</p>
        </div>
      </div>

      {/* Main Entry Table */}
      <div className="mb-3">
        {/* Header Row */}
        <div className="grid grid-cols-12 gap-1 font-bold text-xs bg-gray-200 border-2 border-black">
          <div className="col-span-1 px-1 py-1 border-r border-gray-400 text-center">Sr.</div>
          <div className="col-span-4 px-1 py-1 border-r border-gray-400">Account / Description</div>
          <div className="col-span-2 px-1 py-1 border-r border-gray-400">Reference</div>
          <div className="col-span-2 px-1 py-1 border-r border-gray-400 text-right">Debit (₹)</div>
          <div className="col-span-2 px-1 py-1 text-right">Credit (₹)</div>
          <div className="col-span-1 px-1 py-1 text-center">Chk</div>
        </div>

        {/* Entry Lines */}
        {entries.map((entry, idx) => (
          <div key={idx} className="grid grid-cols-12 gap-1 border-b border-gray-300 py-2 text-xs">
            <div className="col-span-1 px-1 py-1 border-r border-gray-300 text-center text-gray-600">
              {idx + 1}
            </div>
            <div className="col-span-4 px-1 py-1 border-r border-gray-300">
              <p className="font-semibold text-gray-800">{entry.account}</p>
              {entry.description && <p className="text-gray-600 italic text-xs">{entry.description}</p>}
            </div>
            <div className="col-span-2 px-1 py-1 border-r border-gray-300 text-xs text-gray-600">
              {entry.reference || ''}
            </div>
            <div className="col-span-2 px-1 py-1 border-r border-gray-300 text-right font-mono">
              {entry.debit ? entry.debit.toFixed(2) : ''}
            </div>
            <div className="col-span-2 px-1 py-1 text-right font-mono">
              {entry.credit ? entry.credit.toFixed(2) : ''}
            </div>
            <div className="col-span-1 px-1 py-1 text-center">☐</div>
          </div>
        ))}

        {/* Empty rows for manual entries */}
        {entries.length < 12 && (
          <>
            {Array.from({ length: 12 - entries.length }).map((_, idx) => (
              <div key={`empty-${idx}`} className="grid grid-cols-12 gap-1 border-b border-gray-200 py-2">
                <div className="col-span-1 px-1 border-r border-gray-300"></div>
                <div className="col-span-4 px-1 border-r border-gray-300"></div>
                <div className="col-span-2 px-1 border-r border-gray-300"></div>
                <div className="col-span-2 px-1 border-r border-gray-300"></div>
                <div className="col-span-2 px-1"></div>
                <div className="col-span-1 px-1 text-center">☐</div>
              </div>
            ))}
          </>
        )}

        {/* Total Row */}
        <div className="grid grid-cols-12 gap-1 border-t-2 border-b-2 border-black font-bold text-xs py-1">
          <div className="col-span-7"></div>
          <div className="col-span-2 px-1 py-1 border-r border-black text-right font-mono">
            {totalDebit.toFixed(2)}
          </div>
          <div className="col-span-2 px-1 py-1 text-right font-mono">
            {totalCredit.toFixed(2)}
          </div>
          <div className="col-span-1"></div>
        </div>
      </div>

      {/* Narration */}
      {narration && (
        <div className="mb-3 pb-3 border-b-2 border-black">
          <p className="text-xs font-semibold mb-1">NARRATION:</p>
          <p className="text-xs text-gray-700 pl-4 border-l-2 border-gray-400">{narration}</p>
        </div>
      )}

      {/* Balance Status */}
      <div className="grid grid-cols-2 gap-4 text-sm mb-3">
        <div>
          <p className="text-xs font-semibold text-gray-600">Debit Total: ₹ {totalDebit.toFixed(2)}</p>
          <p className="text-xs font-semibold text-gray-600">Credit Total: ₹ {totalCredit.toFixed(2)}</p>
        </div>
        <div className="text-right">
          <p className={`text-xs font-bold ${isBalanced ? 'text-green-700' : 'text-red-700'}`}>
            {isBalanced ? '✓ VOUCHER BALANCED' : '✗ VOUCHER NOT BALANCED'}
          </p>
        </div>
      </div>

      {/* Verification Section */}
      <div className="border-t-2 border-black pt-2">
        <div className="grid grid-cols-3 gap-4 text-xs">
          <div>
            <p className="font-semibold mb-4">Prepared By</p>
            <p className="border-t border-black pt-1">_______________</p>
          </div>
          <div>
            <p className="font-semibold mb-4">Checked By</p>
            <p className="border-t border-black pt-1">_______________</p>
          </div>
          <div>
            <p className="font-semibold mb-4">Approved By</p>
            <p className="border-t border-black pt-1">_______________</p>
          </div>
        </div>
      </div>
    </div>
  )
}
