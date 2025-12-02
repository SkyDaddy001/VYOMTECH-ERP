'use client'

import { Proposal } from '@/types/presales'

interface ProposalListProps {
  proposals: Proposal[]
  loading: boolean
  onEdit: (proposal: Proposal) => void
  onDelete: (proposal: Proposal) => void
  onSend: (proposal: Proposal) => void
  onApprove: (proposal: Proposal) => void
}

export default function ProposalList({ proposals, loading, onEdit, onDelete, onSend, onApprove }: ProposalListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading proposals...</p>
        </div>
      </div>
    )
  }

  if (proposals.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No proposals yet.</p>
      </div>
    )
  }

  const totalAmount = proposals.reduce((sum, p) => sum + p.amount, 0)
  const sentCount = proposals.filter((p) => p.status !== 'draft').length
  const acceptedCount = proposals.filter((p) => p.status === 'accepted').length

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{proposals.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total Amount</p>
          <p className="text-2xl font-bold text-blue-600 mt-1">₹{totalAmount.toLocaleString()}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Sent</p>
          <p className="text-2xl font-bold text-orange-600 mt-1">{sentCount}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Accepted</p>
          <p className="text-2xl font-bold text-green-600 mt-1">{acceptedCount}</p>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-lg shadow overflow-x-auto border border-gray-200">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Proposal</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Amount</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Validity</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {proposals.map((proposal) => (
              <tr key={proposal.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{proposal.title}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-900">₹{proposal.amount.toLocaleString()}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                  {new Date(proposal.validity_date).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <span
                    className={`px-2 py-1 rounded text-xs font-medium ${
                      proposal.status === 'accepted'
                        ? 'bg-green-100 text-green-800'
                        : proposal.status === 'rejected'
                        ? 'bg-red-100 text-red-800'
                        : proposal.status === 'draft'
                        ? 'bg-gray-100 text-gray-800'
                        : 'bg-blue-100 text-blue-800'
                    }`}
                  >
                    {proposal.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                  {proposal.status === 'draft' && (
                    <button onClick={() => onSend(proposal)} className="text-blue-600 hover:text-blue-900 font-medium text-xs">
                      Send
                    </button>
                  )}
                  {proposal.status === 'sent' && (
                    <button onClick={() => onApprove(proposal)} className="text-green-600 hover:text-green-900 font-medium text-xs">
                      Approve
                    </button>
                  )}
                  <button onClick={() => onEdit(proposal)} className="text-orange-600 hover:text-orange-900 font-medium text-xs">
                    Edit
                  </button>
                  <button onClick={() => onDelete(proposal)} className="text-red-600 hover:text-red-900 font-medium text-xs">
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
