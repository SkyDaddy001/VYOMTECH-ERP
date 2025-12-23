'use client'

import { useState, useMemo } from 'react'
import { ProtectedRoute } from '@/hooks/use-auth'
import { Sidebar } from '@/components/layout/sidebar'
import { Header } from '@/components/layout/header'
import { FiSearch, FiPlus, FiFilter, FiCalendar, FiDollarSign } from 'react-icons/fi'
import { format } from 'date-fns'

interface GLEntry {
  id: string | number
  account_code: string
  account_name: string
  description: string
  debit: number
  credit: number
  journal_ref: string
  posting_date: string
  status: 'posted' | 'draft' | 'reversed'
  created_by: string
  created_at: string
}

// Mock GL entries for demo
const mockGLEntries: GLEntry[] = [
  {
    id: 1,
    account_code: '1010',
    account_name: 'Cash in Bank',
    description: 'Bank deposit - Customer advance',
    debit: 50000,
    credit: 0,
    journal_ref: 'JNL-2025-001',
    posting_date: '2025-12-20',
    status: 'posted',
    created_by: 'admin',
    created_at: '2025-12-20T10:30:00Z',
  },
  {
    id: 2,
    account_code: '2010',
    account_name: 'Accounts Payable',
    description: 'Supplier invoice payment',
    debit: 0,
    credit: 35000,
    journal_ref: 'JNL-2025-002',
    posting_date: '2025-12-19',
    status: 'posted',
    created_by: 'admin',
    created_at: '2025-12-19T14:20:00Z',
  },
  {
    id: 3,
    account_code: '4010',
    account_name: 'Sales Revenue',
    description: 'Project sales invoice #INV-001',
    debit: 0,
    credit: 100000,
    journal_ref: 'JNL-2025-003',
    posting_date: '2025-12-18',
    status: 'posted',
    created_by: 'admin',
    created_at: '2025-12-18T09:15:00Z',
  },
  {
    id: 4,
    account_code: '6010',
    account_name: 'Construction Expenses',
    description: 'Material and labor costs',
    debit: 75000,
    credit: 0,
    journal_ref: 'JNL-2025-004',
    posting_date: '2025-12-17',
    status: 'posted',
    created_by: 'admin',
    created_at: '2025-12-17T11:45:00Z',
  },
  {
    id: 5,
    account_code: '1020',
    account_name: 'Accounts Receivable',
    description: 'Customer invoice outstanding',
    debit: 60000,
    credit: 0,
    journal_ref: 'JNL-2025-005',
    posting_date: '2025-12-16',
    status: 'draft',
    created_by: 'admin',
    created_at: '2025-12-16T16:30:00Z',
  },
]

export default function GLEntriesPage() {
  const [entries, setEntries] = useState<GLEntry[]>(mockGLEntries)
  const [searchTerm, setSearchTerm] = useState('')
  const [filterStatus, setFilterStatus] = useState<string>('all')
  const [sortBy, setSortBy] = useState<'date' | 'amount' | 'account' | 'ref'>('date')

  const filteredEntries = useMemo(() => {
    let result = [...entries]

    if (filterStatus !== 'all') {
      result = result.filter(e => e.status === filterStatus)
    }

    if (searchTerm) {
      const query = searchTerm.toLowerCase()
      result = result.filter(e =>
        e.account_code.toLowerCase().includes(query) ||
        e.account_name.toLowerCase().includes(query) ||
        e.description.toLowerCase().includes(query) ||
        e.journal_ref.toLowerCase().includes(query)
      )
    }

    result.sort((a, b) => {
      switch (sortBy) {
        case 'amount':
          return (b.debit + b.credit) - (a.debit + a.credit)
        case 'account':
          return a.account_code.localeCompare(b.account_code)
        case 'ref':
          return a.journal_ref.localeCompare(b.journal_ref)
        case 'date':
        default:
          return new Date(b.posting_date).getTime() - new Date(a.posting_date).getTime()
      }
    })

    return result
  }, [entries, searchTerm, filterStatus, sortBy])

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'posted':
        return 'bg-green-100 text-green-800'
      case 'draft':
        return 'bg-yellow-100 text-yellow-800'
      case 'reversed':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getTotalDebits = () => filteredEntries.reduce((sum, e) => sum + e.debit, 0)
  const getTotalCredits = () => filteredEntries.reduce((sum, e) => sum + e.credit, 0)
  const getBalance = () => getTotalDebits() - getTotalCredits()

  return (
    <ProtectedRoute>
      <div className="flex h-screen bg-gray-50">
        <Sidebar />
        <div className="flex-1 flex flex-col lg:ml-64">
          <Header />
          <main className="flex-1 overflow-auto pt-20 pb-6">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
              {/* Header */}
              <div className="mb-8 flex items-center justify-between">
                <div>
                  <h1 className="text-3xl font-bold text-gray-900">GL Entries</h1>
                  <p className="text-gray-600 mt-2">General ledger journal entries and postings</p>
                </div>
                <button className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg flex items-center gap-2 transition">
                  <FiPlus className="text-lg" />
                  New Entry
                </button>
              </div>

              {/* Summary Cards */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
                <div className="bg-white rounded-lg shadow p-6">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-gray-600 text-sm font-medium">Total Debits</p>
                      <p className="text-3xl font-bold text-gray-900 mt-2">
                        ${getTotalDebits().toLocaleString()}
                      </p>
                    </div>
                    <FiDollarSign className="text-blue-500 text-4xl opacity-20" />
                  </div>
                </div>
                <div className="bg-white rounded-lg shadow p-6">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-gray-600 text-sm font-medium">Total Credits</p>
                      <p className="text-3xl font-bold text-gray-900 mt-2">
                        ${getTotalCredits().toLocaleString()}
                      </p>
                    </div>
                    <FiDollarSign className="text-green-500 text-4xl opacity-20" />
                  </div>
                </div>
                <div className={`bg-white rounded-lg shadow p-6 border-t-4 ${getBalance() === 0 ? 'border-green-500' : 'border-orange-500'}`}>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-gray-600 text-sm font-medium">Balance</p>
                      <p className={`text-3xl font-bold mt-2 ${getBalance() === 0 ? 'text-green-600' : 'text-orange-600'}`}>
                        ${getBalance().toLocaleString()}
                      </p>
                    </div>
                    <div className={`text-4xl opacity-20 ${getBalance() === 0 ? 'text-green-500' : 'text-orange-500'}`}>
                      {getBalance() === 0 ? '✓' : '!'}
                    </div>
                  </div>
                </div>
              </div>

              {/* Filters & Search */}
              <div className="bg-white rounded-lg shadow p-4 mb-6">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  {/* Search */}
                  <div className="relative">
                    <FiSearch className="absolute left-3 top-3 text-gray-400" />
                    <input
                      type="text"
                      placeholder="Search entries..."
                      value={searchTerm}
                      onChange={(e) => setSearchTerm(e.target.value)}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>

                  {/* Status Filter */}
                  <select
                    value={filterStatus}
                    onChange={(e) => setFilterStatus(e.target.value)}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="all">All Status</option>
                    <option value="posted">Posted</option>
                    <option value="draft">Draft</option>
                    <option value="reversed">Reversed</option>
                  </select>

                  {/* Sort */}
                  <select
                    value={sortBy}
                    onChange={(e) => setSortBy(e.target.value as any)}
                    className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="date">Latest First</option>
                    <option value="amount">By Amount</option>
                    <option value="account">By Account</option>
                    <option value="ref">By Reference</option>
                  </select>
                </div>

                <div className="mt-4 text-sm text-gray-600">
                  Showing {filteredEntries.length} of {entries.length} entries
                </div>
              </div>

              {/* Entries Table */}
              {filteredEntries.length === 0 ? (
                <div className="bg-white rounded-lg shadow p-12 text-center">
                  <div className="text-gray-500 mb-4 text-5xl">📊</div>
                  <p className="text-gray-600 text-lg font-medium">No GL entries found</p>
                  <p className="text-gray-500 mt-1">Create a new journal entry to get started</p>
                </div>
              ) : (
                <div className="bg-white rounded-lg shadow overflow-hidden">
                  <table className="w-full">
                    <thead className="bg-gray-50 border-b border-gray-200">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Date</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Reference</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Account</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Description</th>
                        <th className="px-6 py-3 text-right text-xs font-semibold text-gray-700 uppercase tracking-wider">Debit</th>
                        <th className="px-6 py-3 text-right text-xs font-semibold text-gray-700 uppercase tracking-wider">Credit</th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Status</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                      {filteredEntries.map((entry) => (
                        <tr key={entry.id} className="hover:bg-gray-50 transition">
                          <td className="px-6 py-4 text-sm text-gray-700">
                            {format(new Date(entry.posting_date), 'MMM dd, yyyy')}
                          </td>
                          <td className="px-6 py-4 text-sm font-semibold text-gray-900">
                            {entry.journal_ref}
                          </td>
                          <td className="px-6 py-4">
                            <div>
                              <p className="text-sm font-semibold text-gray-900">{entry.account_code}</p>
                              <p className="text-xs text-gray-600">{entry.account_name}</p>
                            </div>
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-700 max-w-xs truncate">
                            {entry.description}
                          </td>
                          <td className="px-6 py-4 text-right">
                            {entry.debit > 0 ? (
                              <span className="text-sm font-semibold text-gray-900">
                                ${entry.debit.toLocaleString()}
                              </span>
                            ) : (
                              <span className="text-sm text-gray-500">-</span>
                            )}
                          </td>
                          <td className="px-6 py-4 text-right">
                            {entry.credit > 0 ? (
                              <span className="text-sm font-semibold text-gray-900">
                                ${entry.credit.toLocaleString()}
                              </span>
                            ) : (
                              <span className="text-sm text-gray-500">-</span>
                            )}
                          </td>
                          <td className="px-6 py-4">
                            <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getStatusColor(entry.status)}`}>
                              {entry.status.charAt(0).toUpperCase() + entry.status.slice(1)}
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>

                  {/* Table Footer with Totals */}
                  <div className="bg-gray-50 border-t border-gray-200 px-6 py-4 flex justify-between font-semibold text-gray-900">
                    <span>Totals:</span>
                    <div className="flex gap-12">
                      <span>Debit: ${getTotalDebits().toLocaleString()}</span>
                      <span>Credit: ${getTotalCredits().toLocaleString()}</span>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </main>
        </div>
      </div>
    </ProtectedRoute>
  )
}
