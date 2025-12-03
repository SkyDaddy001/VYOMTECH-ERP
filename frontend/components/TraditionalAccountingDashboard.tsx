'use client'

import React, { useState } from 'react'
import LedgerBook from '@/components/LedgerBook'
import TraditionalVoucher from '@/components/TraditionalVoucher'
import ReceiptVoucher from '@/components/ReceiptVoucher'
import TrialBalance from '@/components/TrialBalance'

type TabType = 'ledger' | 'vouchers' | 'receipts' | 'trial-balance'

export default function TraditionalAccountingDashboard() {
  const [activeTab, setActiveTab] = useState<TabType>('ledger')

  // Sample data
  const sampleLedgerEntries = [
    {
      id: '1',
      date: '2025-12-01',
      description: 'Opening Balance',
      account: 'Cash',
      debit: 10000,
      credit: 0,
      balance: 10000,
      reference: 'OB'
    },
    {
      id: '2',
      date: '2025-12-02',
      description: 'Sale to ABC Corp',
      account: 'Debtors',
      debit: 5000,
      credit: 0,
      balance: 15000,
      reference: 'INV-001'
    },
    {
      id: '3',
      date: '2025-12-03',
      description: 'Payment received',
      account: 'Cash',
      debit: 3000,
      credit: 0,
      balance: 18000,
      reference: 'CHQ-101'
    }
  ]

  const sampleVoucherEntries = [
    {
      account: 'Sales A/c',
      description: 'Sale of goods',
      debit: 5000,
      reference: 'INV-001'
    },
    {
      account: 'Debtors A/c',
      description: 'Customer XYZ',
      credit: 5000,
      reference: 'INV-001'
    }
  ]

  const sampleTrialBalanceEntries = [
    {
      accountName: 'Cash',
      accountCode: '1010',
      debitBalance: 50000,
      creditBalance: 0,
      accountType: 'Asset' as const
    },
    {
      accountName: 'Bank',
      accountCode: '1020',
      debitBalance: 25000,
      creditBalance: 0,
      accountType: 'Asset' as const
    },
    {
      accountName: 'Accounts Payable',
      accountCode: '2010',
      debitBalance: 0,
      creditBalance: 15000,
      accountType: 'Liability' as const
    },
    {
      accountName: 'Capital',
      accountCode: '3010',
      debitBalance: 0,
      creditBalance: 60000,
      accountType: 'Equity' as const
    }
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-amber-100 to-yellow-50 p-4">
      {/* Header */}
      <div className="max-w-7xl mx-auto mb-8">
        <h1 className="text-4xl font-serif font-bold text-yellow-900 text-center mb-2">
          Accounting System
        </h1>
        <p className="text-center text-gray-700 font-serif text-lg">
          Traditional Ledger & Voucher Management
        </p>
      </div>

      {/* Tab Navigation - Like a Notebook */}
      <div className="max-w-7xl mx-auto mb-6">
        <div className="flex gap-2 flex-wrap border-b-4 border-yellow-900 pb-0">
          {[
            { id: 'ledger', label: '📖 Ledger Books', icon: '📖' },
            { id: 'vouchers', label: '📝 Vouchers', icon: '📝' },
            { id: 'receipts', label: '🧾 Receipt Vouchers', icon: '🧾' },
            { id: 'trial-balance', label: '⚖️ Trial Balance', icon: '⚖️' }
          ].map(tab => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as TabType)}
              className={`px-6 py-3 font-serif font-semibold text-sm transition-all ${
                activeTab === tab.id
                  ? 'bg-yellow-900 text-white border-b-4 border-yellow-900 -mb-1'
                  : 'bg-white text-yellow-900 hover:bg-yellow-50 border border-yellow-200'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Content Area */}
      <div className="max-w-7xl mx-auto">
        {/* Ledger Tab */}
        {activeTab === 'ledger' && (
          <div className="space-y-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-2xl font-serif font-bold text-yellow-900">Ledger Books</h2>
              <button className="px-4 py-2 bg-yellow-900 text-white rounded hover:bg-yellow-800 text-sm font-semibold">
                ➕ New Ledger
              </button>
            </div>
            <LedgerBook
              title="VYOMTECH-ERP"
              accountName="Cash Account"
              entries={sampleLedgerEntries}
              openingBalance={10000}
            />
          </div>
        )}

        {/* Vouchers Tab */}
        {activeTab === 'vouchers' && (
          <div className="space-y-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-2xl font-serif font-bold text-yellow-900">Vouchers</h2>
              <button className="px-4 py-2 bg-yellow-900 text-white rounded hover:bg-yellow-800 text-sm font-semibold">
                ➕ New Voucher
              </button>
            </div>
            <TraditionalVoucher
              voucherNo="JV/2025/001"
              date="2025-12-01"
              voucherType="JV"
              entries={sampleVoucherEntries}
              createdBy="John Doe"
              approvedBy="Finance Manager"
              narration="Sale of goods as per Invoice INV-001"
            />
          </div>
        )}

        {/* Receipts Tab */}
        {activeTab === 'receipts' && (
          <div className="space-y-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-2xl font-serif font-bold text-yellow-900">Receipt Vouchers</h2>
              <button className="px-4 py-2 bg-yellow-900 text-white rounded hover:bg-yellow-800 text-sm font-semibold">
                ➕ New Receipt
              </button>
            </div>
            <ReceiptVoucher
              receiptNo="RCP/2025/001"
              date="2025-12-01"
              receivedFrom="ABC Corporation Pvt Ltd"
              description="Payment for Invoice INV-001"
              amount={5000}
              paymentMode="Cheque"
              chequeNo="123456"
              chequeDate="2025-12-01"
              bankName="State Bank of India"
              createdBy="Cashier"
              approvedBy="Accountant"
            />
          </div>
        )}

        {/* Trial Balance Tab */}
        {activeTab === 'trial-balance' && (
          <div className="space-y-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-2xl font-serif font-bold text-yellow-900">Trial Balance</h2>
              <button className="px-4 py-2 bg-yellow-900 text-white rounded hover:bg-yellow-800 text-sm font-semibold">
                🔄 Refresh
              </button>
            </div>
            <TrialBalance
              date="2025-12-04"
              entries={sampleTrialBalanceEntries}
              preparedBy="Accountant"
              approvedBy="Finance Manager"
            />
          </div>
        )}
      </div>

      {/* Action Buttons at Bottom */}
      <div className="max-w-7xl mx-auto mt-8 flex justify-center gap-4">
        <button className="px-6 py-3 bg-green-700 text-white rounded font-semibold hover:bg-green-800 flex items-center gap-2">
          🖨️ Print
        </button>
        <button className="px-6 py-3 bg-blue-700 text-white rounded font-semibold hover:bg-blue-800 flex items-center gap-2">
          💾 Save
        </button>
        <button className="px-6 py-3 bg-orange-700 text-white rounded font-semibold hover:bg-orange-800 flex items-center gap-2">
          📥 Export to PDF
        </button>
      </div>
    </div>
  )
}
