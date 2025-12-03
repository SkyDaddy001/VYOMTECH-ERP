'use client'

import React from 'react'

interface NotebookLineProps {
  entries: { date: string; description: string; amount: number }[]
  title: string
  pageNo?: number
}

/**
 * Notebook-style UI component - Like writing on actual notebook paper with lines
 * Perfect for traditional accounting entry
 */
export default function NotebookEntry({ entries, title, pageNo }: NotebookLineProps) {
  return (
    <div className="max-w-4xl mx-auto">
      {/* Paper Background with Lines */}
      <div
        className="bg-amber-50 p-8 relative min-h-screen"
        style={{
          backgroundImage: `
            repeating-linear-gradient(
              transparent,
              transparent calc(2rem - 1px),
              #e5d5be 2rem
            )
          `
        }}
      >
        {/* Margin Line */}
        <div
          className="absolute left-0 top-0 bottom-0 w-12 border-r-4 border-red-900"
          style={{ width: '60px' }}
        ></div>

        {/* Content Area */}
        <div className="pl-20">
          {/* Header */}
          <div className="mb-8 pb-4 border-b-2 border-gray-400">
            <h1 className="text-2xl font-serif font-bold text-gray-900">{title}</h1>
            {pageNo && <p className="text-xs text-gray-600 mt-1">Page No: {pageNo}</p>}
          </div>

          {/* Entries in Notebook Style */}
          <div className="space-y-8">
            {entries.map((entry, idx) => (
              <div key={idx} className="relative">
                {/* Date on left margin */}
                <div className="float-left w-16 text-xs font-mono text-gray-600 mr-4">
                  {entry.date}
                </div>

                {/* Entry content */}
                <div className="overflow-hidden">
                  <p className="text-sm text-gray-800 font-serif">
                    <span className="font-semibold">{entry.description}</span>
                    {entry.amount > 0 && (
                      <span className="float-right font-mono font-bold text-gray-900">
                        ₹ {entry.amount.toFixed(2)}
                      </span>
                    )}
                  </p>
                </div>
              </div>
            ))}
          </div>

          {/* Empty lines for writing */}
          <div className="mt-12 pt-8 text-gray-300">
            {Array.from({ length: 5 }).map((_, idx) => (
              <div
                key={`empty-${idx}`}
                className="h-8 border-b border-gray-300 mb-4 relative"
              >
                <span className="absolute left-0 text-xs text-gray-400 -top-3">_________________________________________________________</span>
              </div>
            ))}
          </div>

          {/* Footer */}
          <div className="mt-20 pt-8 border-t-2 border-gray-400">
            <div className="grid grid-cols-3 gap-8 text-xs text-gray-600">
              <div>
                <p className="font-semibold">Recorded By</p>
                <div className="h-8 border-b border-gray-400 mt-2"></div>
              </div>
              <div className="text-center">
                <p className="font-semibold">Date</p>
                <div className="h-8 border-b border-gray-400 mt-2"></div>
              </div>
              <div className="text-right">
                <p className="font-semibold">Verified By</p>
                <div className="h-8 border-b border-gray-400 mt-2"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Print Styles */}
      <style jsx>{`
        @media print {
          div {
            color-adjust: exact !important;
            -webkit-print-color-adjust: exact !important;
            print-color-adjust: exact !important;
          }
        }
      `}</style>
    </div>
  )
}
