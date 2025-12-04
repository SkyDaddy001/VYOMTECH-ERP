'use client'

import React, { useState, useEffect } from 'react'
import { generalLedgerService } from '@/services/api'
import LedgerBook from '@/components/LedgerBook'
import TraditionalVoucher from '@/components/TraditionalVoucher'
import ReceiptVoucher from '@/components/ReceiptVoucher'
import TrialBalance from '@/components/TrialBalance'

type TabType = 'ledger' | 'vouchers' | 'receipts' | 'trial-balance'

export default function TraditionalAccountingDashboard() {
  // State for accounting data
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState<TabType>('ledger')

  // State for different data types
  const [ledgerEntries, setLedgerEntries] = useState<any[]>([])
  const [voucherEntries, setVoucherEntries] = useState<any[]>([])
  const [receiptVouchers, setReceiptVouchers] = useState<any[]>([])
  const [trialBalanceEntries, setTrialBalanceEntries] = useState<any[]>([])

  // Fetch accounting data on mount
  useEffect(() => {
    const fetchAccountingData = async () => {
      try {
        setLoading(true)

        // Fetch ledger entries
        try {
          const ledgerRes = await generalLedgerService.getLedgerEntries()
          setLedgerEntries(ledgerRes.data?.entries || [])
        } catch (err) {
          console.warn('Failed to fetch ledger entries:', err)
          setLedgerEntries([])
        }

        // Fetch journal vouchers
        try {
          const voucherRes = await generalLedgerService.getJournalVouchers()
          setVoucherEntries(voucherRes.data?.vouchers || [])
        } catch (err) {
          console.warn('Failed to fetch journal vouchers:', err)
          setVoucherEntries([])
        }

        // Fetch receipt vouchers
        try {
          const receiptRes = await generalLedgerService.getReceiptVouchers()
          setReceiptVouchers(receiptRes.data?.vouchers || [])
        } catch (err) {
          console.warn('Failed to fetch receipt vouchers:', err)
          setReceiptVouchers([])
        }

        // Fetch trial balance
        try {
          const trialRes = await generalLedgerService.getTrialBalance()
          setTrialBalanceEntries(trialRes.data?.accounts || [])
        } catch (err) {
          console.warn('Failed to fetch trial balance:', err)
          setTrialBalanceEntries([])
        }

        setError(null)
      } catch (err: any) {
        console.error('Failed to fetch accounting data:', err)
        setError(err.message || 'Failed to load accounting data')
      } finally {
        setLoading(false)
      }
    }

    fetchAccountingData()
  }, [])

  // Use real data or fallback values
  const displayLedgerEntries = ledgerEntries.length > 0 ? ledgerEntries : [
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

  const displayVoucherEntries = voucherEntries.length > 0 ? voucherEntries : [
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

  const displayReceiptVouchers = receiptVouchers.length > 0 ? receiptVouchers : [
    {
      date: '2025-12-03',
      account: 'Cash',
      description: 'Payment received',
      amount: 5000,
      reference: 'CHQ-101'
    }
  ]

  const displayTrialBalance = trialBalanceEntries.length > 0 ? trialBalanceEntries : [
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
        <p className="text-center text-gray-700 mb-6">Traditional Double-Entry Bookkeeping</p>

        {error && <div className="bg-red-100 text-red-800 p-4 rounded-lg mb-4">{error}</div>}
        {loading && <div className="bg-blue-100 text-blue-800 p-4 rounded-lg mb-4">Loading accounting data...</div>}

        {/* Tab Navigation */}
        <div className="flex gap-4 mb-8 flex-wrap justify-center">
          {[
            { id: 'ledger', label: '📘 Ledger Book', icon: '📚' },
            { id: 'vouchers', label: '📋 Journal Vouchers', icon: '📄' },
            { id: 'receipts', label: '💵 Receipt Vouchers', icon: '💳' },
            { id: 'trial-balance', label: '⚖️ Trial Balance', icon: '📊' }
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as TabType)}
              className={`px-6 py-3 rounded-lg font-semibold transition-all ${
                activeTab === tab.id
                  ? 'bg-yellow-600 text-white shadow-lg'
                  : 'bg-white text-yellow-900 border-2 border-yellow-200 hover:bg-yellow-50'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Content Area */}
      <div className="max-w-7xl mx-auto">
        {activeTab === 'ledger' && <LedgerBook entries={displayLedgerEntries} />}
        {activeTab === 'vouchers' && <TraditionalVoucher entries={displayVoucherEntries} />}
        {activeTab === 'receipts' && <ReceiptVoucher entries={displayReceiptVouchers} />}
        {activeTab === 'trial-balance' && <TrialBalance entries={displayTrialBalance} />}
      </div>
    </div>
  )
}
