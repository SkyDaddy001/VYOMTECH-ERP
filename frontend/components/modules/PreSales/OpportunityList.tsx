'use client'

import { PreSalesOpportunity } from '@/types/presales'

interface OpportunityListProps {
  opportunities: PreSalesOpportunity[]
  loading: boolean
  onEdit: (opportunity: PreSalesOpportunity) => void
  onDelete: (opportunity: PreSalesOpportunity) => void
  onMove: (opportunity: PreSalesOpportunity, stage: string) => void
}

export default function OpportunityList({ opportunities, loading, onEdit, onDelete, onMove }: OpportunityListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading opportunities...</p>
        </div>
      </div>
    )
  }

  if (opportunities.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No opportunities yet. Create your first opportunity to get started.</p>
      </div>
    )
  }

  const totalValue = opportunities.reduce((sum, o) => sum + (o.value || 0), 0)
  const wonCount = opportunities.filter((o) => o.stage === 'closed_won').length

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{opportunities.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Pipeline Value</p>
          <p className="text-2xl font-bold text-blue-600 mt-1">₹{totalValue.toLocaleString()}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Won</p>
          <p className="text-2xl font-bold text-green-600 mt-1">{wonCount}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Win Rate</p>
          <p className="text-2xl font-bold text-purple-600 mt-1">{((wonCount / opportunities.length) * 100).toFixed(0)}%</p>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-lg shadow overflow-x-auto border border-gray-200">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Opportunity</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Value</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Probability</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Stage</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Close Date</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {opportunities.map((opp) => (
              <tr key={opp.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{opp.name}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">₹{(opp.value || 0).toLocaleString()}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <div className="flex items-center gap-2">
                    <div className="w-16 bg-gray-200 rounded h-2">
                      <div className="bg-blue-600 h-2 rounded" style={{ width: `${opp.probability || 0}%` }}></div>
                    </div>
                    <span className="text-xs text-gray-600">{opp.probability || 0}%</span>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <select
                    value={opp.stage}
                    onChange={(e) => onMove(opp, e.target.value)}
                    className={`px-2 py-1 rounded text-xs font-medium border-0 ${
                      opp.stage === 'closed_won'
                        ? 'bg-green-100 text-green-800'
                        : opp.stage === 'closed_lost'
                        ? 'bg-red-100 text-red-800'
                        : 'bg-yellow-100 text-yellow-800'
                    }`}
                  >
                    <option value="prospecting">Prospecting</option>
                    <option value="qualification">Qualification</option>
                    <option value="proposal">Proposal</option>
                    <option value="negotiation">Negotiation</option>
                    <option value="closed_won">Closed Won</option>
                    <option value="closed_lost">Closed Lost</option>
                  </select>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                  {new Date(opp.expected_close_date).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                  <button onClick={() => onEdit(opp)} className="text-green-600 hover:text-green-900 font-medium">
                    Edit
                  </button>
                  <button onClick={() => onDelete(opp)} className="text-red-600 hover:text-red-900 font-medium">
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
