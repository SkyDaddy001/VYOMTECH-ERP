'use client'

import { Expense } from '@/types/accounts'

interface ExpenseListProps {
  expenses: Expense[]
  loading: boolean
  onEdit: (expense: Expense) => void
  onDelete: (expense: Expense) => void
  onStatusChange: (expense: Expense, status: string) => void
}

export default function ExpenseList({ expenses, loading, onEdit, onDelete, onStatusChange }: ExpenseListProps) {
  const statuses = ['draft', 'submitted', 'approved', 'rejected', 'paid']
  const statusColors: Record<string, string> = {
    draft: 'bg-gray-100 text-gray-800',
    submitted: 'bg-blue-100 text-blue-800',
    approved: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800',
    paid: 'bg-purple-100 text-purple-800',
  }

  if (loading) return <div className="text-center py-8 text-gray-500">Loading expenses...</div>

  const totalExpenses = expenses.reduce((sum, e) => sum + e.amount, 0)
  const approvedExpenses = expenses.filter(e => e.status === 'approved').reduce((sum, e) => sum + e.amount, 0)

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-gradient-to-br from-red-50 to-red-100 rounded-lg p-4 border border-red-200">
          <p className="text-gray-600 text-xs font-medium">Total Expenses</p>
          <p className="text-2xl font-bold text-red-600 mt-1">₹{totalExpenses.toLocaleString()}</p>
        </div>
        <div className="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-4 border border-green-200">
          <p className="text-gray-600 text-xs font-medium">Approved</p>
          <p className="text-2xl font-bold text-green-600 mt-1">₹{approvedExpenses.toLocaleString()}</p>
        </div>
        <div className="bg-gradient-to-br from-orange-50 to-orange-100 rounded-lg p-4 border border-orange-200">
          <p className="text-gray-600 text-xs font-medium">Count</p>
          <p className="text-2xl font-bold text-orange-600 mt-1">{expenses.length}</p>
        </div>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Expense #</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Category</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Description</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Amount</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Date</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Status</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {expenses.map((expense) => (
              <tr key={expense.id} className="hover:bg-gray-50 transition">
                <td className="px-6 py-4 text-sm font-medium text-gray-900">{expense.expense_number}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{expense.category}</td>
                <td className="px-6 py-4 text-sm text-gray-600 line-clamp-1">{expense.description}</td>
                <td className="px-6 py-4 text-sm font-semibold text-gray-900">₹{expense.amount.toLocaleString()}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{new Date(expense.expense_date).toLocaleDateString()}</td>
                <td className="px-6 py-4">
                  <select
                    value={expense.status}
                    onChange={(e) => onStatusChange(expense, e.target.value)}
                    className={`text-xs font-medium px-3 py-1 rounded-full cursor-pointer ${statusColors[expense.status]}`}
                  >
                    {statuses.map((s) => (
                      <option key={s} value={s}>
                        {s.charAt(0).toUpperCase() + s.slice(1)}
                      </option>
                    ))}
                  </select>
                </td>
                <td className="px-6 py-4 text-sm space-x-2">
                  <button
                    onClick={() => onEdit(expense)}
                    className="text-blue-600 hover:text-blue-900 font-medium"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(expense)}
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
