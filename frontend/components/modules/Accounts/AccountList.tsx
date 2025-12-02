'use client'

import { Account } from '@/types/accounts'

interface AccountListProps {
  accounts: Account[]
  loading: boolean
  onEdit: (account: Account) => void
  onDelete: (account: Account) => void
}

export default function AccountList({ accounts, loading, onEdit, onDelete }: AccountListProps) {
  if (loading) return <div className="text-center py-8 text-gray-500">Loading accounts...</div>

  const totalBalance = accounts.reduce((sum, a) => sum + a.balance, 0)
  const activeAccounts = accounts.filter(a => a.status === 'active').length

  const accountTypeIcons: Record<string, string> = {
    checking: '💼',
    savings: '🏦',
    credit: '💳',
  }

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-gradient-to-br from-violet-50 to-violet-100 rounded-lg p-4 border border-violet-200">
          <p className="text-gray-600 text-xs font-medium">Total Balance</p>
          <p className="text-2xl font-bold text-violet-600 mt-1">₹{totalBalance.toLocaleString()}</p>
        </div>
        <div className="bg-gradient-to-br from-fuchsia-50 to-fuchsia-100 rounded-lg p-4 border border-fuchsia-200">
          <p className="text-gray-600 text-xs font-medium">Active Accounts</p>
          <p className="text-2xl font-bold text-fuchsia-600 mt-1">{activeAccounts}</p>
        </div>
        <div className="bg-gradient-to-br from-pink-50 to-pink-100 rounded-lg p-4 border border-pink-200">
          <p className="text-gray-600 text-xs font-medium">Total Accounts</p>
          <p className="text-2xl font-bold text-pink-600 mt-1">{accounts.length}</p>
        </div>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Account Name</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Type</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Bank</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Number</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Balance</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Status</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {accounts.map((account) => (
              <tr key={account.id} className="hover:bg-gray-50 transition">
                <td className="px-6 py-4 text-sm font-medium text-gray-900">{account.account_name}</td>
                <td className="px-6 py-4 text-sm">
                  <span className="text-lg">{accountTypeIcons[account.account_type]}</span>
                  {account.account_type.replace('_', ' ').charAt(0).toUpperCase() + account.account_type.slice(1)}
                </td>
                <td className="px-6 py-4 text-sm text-gray-600">{account.bank_name}</td>
                <td className="px-6 py-4 text-sm text-gray-600 font-mono">***{account.account_number.slice(-4)}</td>
                <td className="px-6 py-4 text-sm font-semibold text-gray-900">₹{account.balance.toLocaleString()}</td>
                <td className="px-6 py-4">
                  <span
                    className={`text-xs font-semibold px-3 py-1 rounded-full ${
                      account.status === 'active'
                        ? 'bg-green-100 text-green-800'
                        : account.status === 'inactive'
                        ? 'bg-gray-100 text-gray-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {account.status.charAt(0).toUpperCase() + account.status.slice(1)}
                  </span>
                </td>
                <td className="px-6 py-4 text-sm space-x-2">
                  <button
                    onClick={() => onEdit(account)}
                    className="text-blue-600 hover:text-blue-900 font-medium"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(account)}
                    className="text-red-600 hover:text-red-900 font-medium"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
