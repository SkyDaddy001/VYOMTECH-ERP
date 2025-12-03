'use client'

import React, { useState } from 'react'

interface ReceiptVoucherProps {
  receiptNo: string
  date: string
  receivedFrom: string
  description: string
  amount: number
  paymentMode: 'Cash' | 'Cheque' | 'Bank Transfer' | 'Card'
  chequeNo?: string
  chequeDate?: string
  bankName?: string
  accountName?: string
  accountNumber?: string
  createdBy?: string
  approvedBy?: string
  remarks?: string
}

export default function ReceiptVoucher({
  receiptNo,
  date,
  receivedFrom,
  description,
  amount,
  paymentMode,
  chequeNo,
  chequeDate,
  bankName,
  accountName,
  accountNumber,
  createdBy,
  approvedBy,
  remarks
}: ReceiptVoucherProps) {
  return (
    <div className="font-serif bg-white p-8 border-4 border-black max-w-4xl mx-auto">
      {/* Header - Official Receipt Style */}
      <div className="border-b-4 border-black pb-3 mb-4">
        <div className="text-center mb-2">
          <h1 className="text-2xl font-bold text-black">RECEIPT VOUCHER</h1>
          <p className="text-xs text-gray-600">Original</p>
        </div>
        <div className="flex justify-between items-center">
          <div>
            <p className="text-sm font-semibold">Receipt No: <span className="font-mono font-bold text-lg">{receiptNo}</span></p>
            <p className="text-sm font-semibold">Date: <span className="font-mono">{new Date(date).toLocaleDateString('en-IN')}</span></p>
          </div>
          <div className="text-right">
            <p className="text-xs text-gray-600">For Official Use</p>
            <p className="border-t border-black mt-1 pt-1 font-mono text-sm">Ref: __________</p>
          </div>
        </div>
      </div>

      {/* Main Receipt Content */}
      <div className="mb-4">
        {/* Received From */}
        <div className="grid grid-cols-2 gap-4 mb-4 text-sm">
          <div>
            <p className="font-semibold text-gray-700 mb-1">RECEIVED FROM:</p>
            <p className="border-b-2 border-black pb-2 font-bold text-lg text-gray-900">{receivedFrom}</p>
          </div>
          <div>
            <p className="font-semibold text-gray-700 mb-1">AMOUNT RECEIVED:</p>
            <p className="border-b-2 border-black pb-2 font-bold text-2xl text-green-700 font-mono">₹ {amount.toFixed(2)}</p>
          </div>
        </div>

        {/* Amount in Words */}
        <div className="mb-4 pb-3 border-b border-gray-400">
          <p className="font-semibold text-xs text-gray-700 mb-1">AMOUNT IN WORDS:</p>
          <p className="text-sm text-gray-900 font-semibold italic">
            {numberToWords(amount)} Only
          </p>
        </div>

        {/* Description/Purpose */}
        <div className="mb-4 pb-3 border-b border-gray-400">
          <p className="font-semibold text-xs text-gray-700 mb-1">TOWARDS/DESCRIPTION:</p>
          <p className="text-sm text-gray-800 pl-2">{description}</p>
        </div>

        {/* Payment Details */}
        <div className="mb-4">
          <p className="font-semibold text-xs text-gray-700 mb-2">MODE OF PAYMENT:</p>
          <div className="grid grid-cols-2 gap-4">
            <div className="text-sm">
              <p className="mb-1">Payment Method: <span className="font-semibold text-gray-900">{paymentMode}</span></p>
              
              {paymentMode === 'Cheque' && (
                <>
                  <p className="mb-1">Cheque No: <span className="font-mono font-semibold">{chequeNo || '___________'}</span></p>
                  <p className="mb-1">Cheque Date: <span className="font-mono">{chequeDate ? new Date(chequeDate).toLocaleDateString('en-IN') : '___________'}</span></p>
                  <p className="mb-1">Bank: <span className="font-semibold">{bankName || '___________'}</span></p>
                </>
              )}
              
              {paymentMode === 'Bank Transfer' && (
                <>
                  <p className="mb-1">Account Name: <span className="font-semibold">{accountName || '___________'}</span></p>
                  <p className="mb-1">Account No: <span className="font-mono">{accountNumber || '___________'}</span></p>
                  <p className="mb-1">Bank: <span className="font-semibold">{bankName || '___________'}</span></p>
                </>
              )}
            </div>

            {/* Account Coding */}
            <div className="border-l-2 border-gray-400 pl-4">
              <p className="text-xs font-semibold text-gray-700 mb-2">ACCOUNT CODING:</p>
              <div className="space-y-1">
                <p className="text-xs">Debit A/c: <span className="font-mono border-b border-gray-300">_____________</span></p>
                <p className="text-xs">Credit A/c: <span className="font-mono border-b border-gray-300">_____________</span></p>
              </div>
            </div>
          </div>
        </div>

        {/* Remarks */}
        {remarks && (
          <div className="mb-4 pb-3 border-b border-gray-400">
            <p className="font-semibold text-xs text-gray-700 mb-1">REMARKS:</p>
            <p className="text-xs text-gray-600">{remarks}</p>
          </div>
        )}
      </div>

      {/* Signature Section */}
      <div className="border-t-2 border-black pt-4">
        <div className="grid grid-cols-3 gap-6 text-center">
          <div>
            <p className="text-xs font-semibold text-gray-700 mb-8">Prepared By</p>
            <p className="border-t-2 border-black pt-2 font-semibold text-sm">{createdBy || '_______________'}</p>
          </div>
          <div>
            <p className="text-xs font-semibold text-gray-700 mb-8">Checked By</p>
            <p className="border-t-2 border-black pt-2 font-semibold text-sm">_______________</p>
          </div>
          <div>
            <p className="text-xs font-semibold text-gray-700 mb-8">Approved By</p>
            <p className="border-t-2 border-black pt-2 font-semibold text-sm">{approvedBy || '_______________'}</p>
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="mt-4 pt-4 border-t border-gray-400 text-center">
        <p className="text-xs text-gray-500">This is a computer generated receipt and requires no signature</p>
        <p className="text-xs text-gray-500 mt-1">Date Printed: {new Date().toLocaleDateString('en-IN')}</p>
      </div>
    </div>
  )
}

// Helper function to convert numbers to words
function numberToWords(num: number): string {
  const ones = ['', 'One', 'Two', 'Three', 'Four', 'Five', 'Six', 'Seven', 'Eight', 'Nine']
  const tens = ['', '', 'Twenty', 'Thirty', 'Forty', 'Fifty', 'Sixty', 'Seventy', 'Eighty', 'Ninety']
  const teens = ['Ten', 'Eleven', 'Twelve', 'Thirteen', 'Fourteen', 'Fifteen', 'Sixteen', 'Seventeen', 'Eighteen', 'Nineteen']
  const scales = ['', 'Thousand', 'Lakh', 'Crore']

  if (num === 0) return 'Zero'

  let result = ''
  let scaleIndex = 0

  while (num > 0) {
    const remainder = num % (scaleIndex === 0 ? 1000 : 100)
    if (remainder !== 0) {
      result = convertHundreds(remainder, ones, tens, teens) + (scales[scaleIndex] ? ' ' + scales[scaleIndex] : '') + (result ? ' ' + result : '')
    }
    num = Math.floor(num / (scaleIndex === 0 ? 1000 : 100))
    scaleIndex++
  }

  return result.trim()
}

function convertHundreds(num: number, ones: string[], tens: string[], teens: string[]): string {
  let result = ''
  const hundreds = Math.floor(num / 100)
  const remainder = num % 100

  if (hundreds > 0) {
    result += ones[hundreds] + ' Hundred'
  }

  if (remainder >= 20) {
    if (result) result += ' '
    result += tens[Math.floor(remainder / 10)]
    if (remainder % 10 > 0) {
      result += ' ' + ones[remainder % 10]
    }
  } else if (remainder >= 10) {
    if (result) result += ' '
    result += teens[remainder - 10]
  } else if (remainder > 0) {
    if (result) result += ' '
    result += ones[remainder]
  }

  return result
}
