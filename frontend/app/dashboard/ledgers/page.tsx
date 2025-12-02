'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import LedgerList from '@/components/modules/Ledgers/LedgerList'
import CustomerSummaryList from '@/components/modules/Ledgers/CustomerSummaryList'
import { ledgersService } from '@/services/ledgers.service'
import { Ledger, LedgerSummary, LedgerMetrics } from '@/types/ledgers'

type TabType = 'transactions' | 'summaries'

export default function LedgersPage() {
  const [activeTab, setActiveTab] = useState<TabType>('transactions')

  // Data states
  const [entries, setEntries] = useState<Ledger[]>([])
  const [summaries, setSummaries] = useState<LedgerSummary[]>([])
  const [metrics, setMetrics] = useState<LedgerMetrics | null>(null)

  // Loading states
  const [entriesLoading, setEntriesLoading] = useState(false)
  const [summariesLoading, setSummariesLoading] = useState(false)

  // Load ledger entries
  const loadEntries = async () => {
    setEntriesLoading(true)
    try {
      const data = await ledgersService.getLedgerEntries()
      setEntries(data.sort((a, b) => new Date(b.transaction_date).getTime() - new Date(a.transaction_date).getTime()))
    } catch (error) {
      toast.error('Failed to load ledger entries')
    } finally {
      setEntriesLoading(false)
    }
  }

  // Load customer summaries
  const loadSummaries = async () => {
    setSummariesLoading(true)
    try {
      const data = await ledgersService.getAllCustomerSummaries()
      setSummaries(data)
    } catch (error) {
      toast.error('Failed to load summaries')
    } finally {
      setSummariesLoading(false)
    }
  }

  // Load metrics
  const loadMetrics = async () => {
    try {
      const data = await ledgersService.getMetrics()
      setMetrics(data)
    } catch (error) {
      // Metrics are optional
    }
  }

  // Load data on tab change
  useEffect(() => {
    loadMetrics()
    if (activeTab === 'transactions') {
      loadEntries()
    } else if (activeTab === 'summaries') {
      loadSummaries()
    }
  }, [activeTab])

  const handleDeleteEntry = async (entry: Ledger) => {
    if (!confirm('Are you sure?')) return
    try {
      await ledgersService.deleteEntry(entry.id || '')
      toast.success('Entry deleted!')
      loadEntries()
    } catch (error) {
      toast.error('Failed to delete entry')
    }
  }

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">Account Ledger</h1>
          <p className="text-gray-600">Track customer account transactions and balances</p>
        </div>

        {/* Metrics */}
        {metrics && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Total Transactions</p>
              <p className="text-2xl font-bold text-blue-600 mt-1">{metrics.total_transactions}</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Total Debits</p>
              <p className="text-2xl font-bold text-red-600 mt-1">₹{(metrics.total_debit / 10000000).toFixed(0)}Cr</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Total Credits</p>
              <p className="text-2xl font-bold text-green-600 mt-1">₹{(metrics.total_credit / 10000000).toFixed(0)}Cr</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Outstanding</p>
              <p className="text-2xl font-bold text-orange-600 mt-1">₹{(metrics.total_outstanding / 10000000).toFixed(0)}Cr</p>
            </div>
          </div>
        )}

        {/* Tabs */}
        <div className="bg-white rounded-lg shadow mb-6 border border-gray-200">
          <div className="flex flex-wrap border-b border-gray-200">
            <button
              onClick={() => setActiveTab('transactions')}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'transactions'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Transactions
            </button>
            <button
              onClick={() => setActiveTab('summaries')}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'summaries'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Customer Summaries
            </button>
          </div>

          {/* Tab Content */}
          <div className="p-6">
            {/* Transactions Tab */}
            {activeTab === 'transactions' && (
              <LedgerList
                entries={entries}
                loading={entriesLoading}
                onEdit={(entry) => {
                  // Edit would require a form modal
                  toast('Edit functionality coming soon', { icon: 'ℹ️' })
                }}
                onDelete={handleDeleteEntry}
              />
            )}

            {/* Summaries Tab */}
            {activeTab === 'summaries' && (
              <CustomerSummaryList
                summaries={summaries}
                loading={summariesLoading}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
