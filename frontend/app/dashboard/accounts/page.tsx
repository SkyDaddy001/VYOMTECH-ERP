'use client'

import { useState } from 'react'
import { SectionCard } from '@/components/ui/section-card'
import { StatCard } from '@/components/ui/stat-card'

type TabType = 'dashboard' | 'accounts' | 'entries' | 'reports' | 'reconciliation'

interface Account {
  code: string
  name: string
  type: string
  balance: number
  status: 'active' | 'inactive'
}

interface JournalEntry {
  id: string
  date: string
  description: string
  debitAmount: number
  creditAmount: number
  status: 'draft' | 'posted'
}

export default function AccountsPage() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard')

  const [accounts] = useState<Account[]>([
    { code: '1000', name: 'Cash at Bank', type: 'Asset', balance: 250000, status: 'active' },
    { code: '1010', name: 'Petty Cash', type: 'Asset', balance: 5000, status: 'active' },
    { code: '1100', name: 'Accounts Receivable', type: 'Asset', balance: 150000, status: 'active' },
    { code: '2000', name: 'Accounts Payable', type: 'Liability', balance: -75000, status: 'active' },
    { code: '3000', name: 'Capital Stock', type: 'Equity', balance: -500000, status: 'active' },
    { code: '4000', name: 'Sales Revenue', type: 'Income', balance: -450000, status: 'active' },
    { code: '5000', name: 'Operating Expenses', type: 'Expense', balance: 200000, status: 'active' },
  ])

  const [journalEntries] = useState<JournalEntry[]>([
    { id: 'JE-001', date: '2025-11-26', description: 'Cash sale to customer ABC', debitAmount: 50000, creditAmount: 50000, status: 'posted' },
    { id: 'JE-002', date: '2025-11-25', description: 'Payment to vendor XYZ', debitAmount: 25000, creditAmount: 25000, status: 'posted' },
    { id: 'JE-003', date: '2025-11-24', description: 'Accrued salary expense', debitAmount: 15000, creditAmount: 15000, status: 'draft' },
  ])

  const stats = [
    { label: 'Total Assets', value: '$405K', icon: '📊' },
    { label: 'Total Liabilities', value: '$75K', icon: '📉' },
    { label: 'Total Equity', value: '$330K', icon: '💼' },
    { label: 'Active Accounts', value: accounts.filter(a => a.status === 'active').length, icon: '✅' },
  ]

  const tabs: Array<{ id: TabType; label: string }> = [
    { id: 'dashboard', label: 'Dashboard' },
    { id: 'accounts', label: 'Chart of Accounts' },
    { id: 'entries', label: 'Journal Entries' },
    { id: 'reports', label: 'Financial Reports' },
    { id: 'reconciliation', label: 'Reconciliation' },
  ]

  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-indigo-600 to-indigo-800 rounded-lg p-6 text-white">
        <h1 className="text-3xl font-bold">Accounts (GL) Module</h1>
        <p className="text-indigo-100 mt-2">Manage chart of accounts, journal entries, and financial reports</p>
      </div>

      <div className="flex gap-2 border-b border-gray-200 overflow-x-auto">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`px-4 py-3 font-medium border-b-2 transition whitespace-nowrap ${
              activeTab === tab.id
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Dashboard Tab */}
      {activeTab === 'dashboard' && (
        <div className="space-y-6">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {stats.map((stat, i) => (
              <StatCard key={i} label={stat.label} value={stat.value} icon={stat.icon} />
            ))}
          </div>

          <SectionCard title="Account Summary">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="bg-blue-50 p-4 rounded-lg">
                <h3 className="font-semibold text-blue-900 mb-2">Assets</h3>
                <p className="text-2xl font-bold text-blue-600">$405K</p>
                <p className="text-xs text-blue-700 mt-1">3 accounts active</p>
              </div>
              <div className="bg-red-50 p-4 rounded-lg">
                <h3 className="font-semibold text-red-900 mb-2">Liabilities</h3>
                <p className="text-2xl font-bold text-red-600">$75K</p>
                <p className="text-xs text-red-700 mt-1">1 account active</p>
              </div>
              <div className="bg-green-50 p-4 rounded-lg">
                <h3 className="font-semibold text-green-900 mb-2">Equity</h3>
                <p className="text-2xl font-bold text-green-600">$330K</p>
                <p className="text-xs text-green-700 mt-1">1 account active</p>
              </div>
            </div>
          </SectionCard>
        </div>
      )}

      {/* Chart of Accounts Tab */}
      {activeTab === 'accounts' && (
        <SectionCard title="Chart of Accounts" action={<button className="px-3 py-1 bg-indigo-600 text-white text-sm rounded hover:bg-indigo-700">+ New Account</button>}>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="px-4 py-2 text-left font-medium text-gray-700">Code</th>
                  <th className="px-4 py-2 text-left font-medium text-gray-700">Name</th>
                  <th className="px-4 py-2 text-left font-medium text-gray-700">Type</th>
                  <th className="px-4 py-2 text-left font-medium text-gray-700">Balance</th>
                  <th className="px-4 py-2 text-left font-medium text-gray-700">Status</th>
                  <th className="px-4 py-2 text-left font-medium text-gray-700">Actions</th>
                </tr>
              </thead>
              <tbody>
                {accounts.map((account) => (
                  <tr key={account.code} className="border-b hover:bg-gray-50">
                    <td className="px-4 py-3 font-medium text-gray-900">{account.code}</td>
                    <td className="px-4 py-3 text-gray-600">{account.name}</td>
                    <td className="px-4 py-3 text-gray-600">{account.type}</td>
                    <td className="px-4 py-3 font-medium text-gray-900">${(account.balance / 1000).toFixed(0)}K</td>
                    <td className="px-4 py-3">
                      <span className={`text-xs font-medium px-2 py-1 rounded ${account.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>
                        {account.status.charAt(0).toUpperCase() + account.status.slice(1)}
                      </span>
                    </td>
                    <td className="px-4 py-3 flex gap-2">
                      <button className="text-blue-600 hover:text-blue-800 text-xs font-medium">Edit</button>
                      <button className="text-red-600 hover:text-red-800 text-xs font-medium">Delete</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </SectionCard>
      )}

      {/* Journal Entries Tab */}
      {activeTab === 'entries' && (
        <SectionCard title="Journal Entries" action={<button className="px-3 py-1 bg-indigo-600 text-white text-sm rounded hover:bg-indigo-700">+ New Entry</button>}>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 border-b">
                <tr>
                  <th className="px-4 py-2 text-left font-medium text-gray-700">ID</th>
                  <th className="px-4 py-2 text-left font-medium text-gray-700">Date</th>
                  <th className="px-4 py-2 text-left font-medium text-gray-700">Description</th>
                  <th className="px-4 py-2 text-right font-medium text-gray-700">Debit</th>
                  <th className="px-4 py-2 text-right font-medium text-gray-700">Credit</th>
                  <th className="px-4 py-2 text-left font-medium text-gray-700">Status</th>
                </tr>
              </thead>
              <tbody>
                {journalEntries.map((entry) => (
                  <tr key={entry.id} className="border-b hover:bg-gray-50">
                    <td className="px-4 py-3 font-medium text-gray-900">{entry.id}</td>
                    <td className="px-4 py-3 text-gray-600">{entry.date}</td>
                    <td className="px-4 py-3 text-gray-600">{entry.description}</td>
                    <td className="px-4 py-3 text-right font-medium text-gray-900">${(entry.debitAmount / 1000).toFixed(0)}K</td>
                    <td className="px-4 py-3 text-right font-medium text-gray-900">${(entry.creditAmount / 1000).toFixed(0)}K</td>
                    <td className="px-4 py-3">
                      <span className={`text-xs font-medium px-2 py-1 rounded ${entry.status === 'posted' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}`}>
                        {entry.status.charAt(0).toUpperCase() + entry.status.slice(1)}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </SectionCard>
      )}

      {/* Reports Tab */}
      {activeTab === 'reports' && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <SectionCard title="Balance Sheet">
            <p className="text-gray-600 text-sm mb-4">As of November 26, 2025</p>
            <div className="space-y-3">
              <div className="flex justify-between border-b pb-2">
                <span className="font-semibold">Assets</span>
                <span className="font-semibold">$405,000</span>
              </div>
              <div className="flex justify-between border-b pb-2">
                <span className="font-semibold">Liabilities</span>
                <span className="font-semibold">$75,000</span>
              </div>
              <div className="flex justify-between border-b pb-2">
                <span className="font-semibold">Equity</span>
                <span className="font-semibold">$330,000</span>
              </div>
            </div>
          </SectionCard>
          <SectionCard title="Income Statement">
            <p className="text-gray-600 text-sm mb-4">For November 2025</p>
            <div className="space-y-3">
              <div className="flex justify-between border-b pb-2">
                <span className="font-semibold">Revenue</span>
                <span className="font-semibold text-green-600">$450,000</span>
              </div>
              <div className="flex justify-between border-b pb-2">
                <span className="font-semibold">Expenses</span>
                <span className="font-semibold text-red-600">$200,000</span>
              </div>
              <div className="flex justify-between border-b pb-2">
                <span className="font-semibold">Net Income</span>
                <span className="font-semibold text-green-700">$250,000</span>
              </div>
            </div>
          </SectionCard>
        </div>
      )}

      {/* Reconciliation Tab */}
      {activeTab === 'reconciliation' && (
        <SectionCard title="Bank Reconciliation">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-blue-50 p-4 rounded-lg">
              <h3 className="font-semibold text-blue-900 mb-3">Bank Statement</h3>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span>Opening Balance:</span>
                  <span className="font-medium">$200,000</span>
                </div>
                <div className="flex justify-between">
                  <span>Deposits:</span>
                  <span className="font-medium text-green-600">+$100,000</span>
                </div>
                <div className="flex justify-between">
                  <span>Withdrawals:</span>
                  <span className="font-medium text-red-600">-$50,000</span>
                </div>
                <div className="flex justify-between border-t pt-2 font-bold">
                  <span>Closing Balance:</span>
                  <span>$250,000</span>
                </div>
              </div>
            </div>
            <div className="bg-green-50 p-4 rounded-lg">
              <h3 className="font-semibold text-green-900 mb-3">Cash Book</h3>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span>Opening Balance:</span>
                  <span className="font-medium">$200,000</span>
                </div>
                <div className="flex justify-between">
                  <span>Deposits:</span>
                  <span className="font-medium text-green-600">+$100,000</span>
                </div>
                <div className="flex justify-between">
                  <span>Withdrawals:</span>
                  <span className="font-medium text-red-600">-$50,000</span>
                </div>
                <div className="flex justify-between border-t pt-2 font-bold">
                  <span>Closing Balance:</span>
                  <span>$250,000</span>
                </div>
              </div>
            </div>
          </div>
          <div className="mt-6 p-4 bg-green-100 border-l-4 border-green-600 rounded">
            <p className="text-green-800 font-semibold">✅ Reconciliation Matched: Bank and Cash Book balances are in agreement</p>
          </div>
        </SectionCard>
      )}
    </div>
  )
}
