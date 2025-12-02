'use client'

import React, { useState } from 'react'
import { DocumentTracker, SnagList, ChangeRequest } from '@/types/postsales'

interface DocumentListProps {
  documents: DocumentTracker[]
  loading: boolean
  onEdit: (doc: DocumentTracker) => void
  onStatusChange: (doc: DocumentTracker, status: string) => void
}

export default function DocumentList({ documents, loading, onEdit, onStatusChange }: DocumentListProps) {
  const [filterStatus, setFilterStatus] = useState<string>('all')

  const filteredDocs = filterStatus === 'all'
    ? documents
    : documents.filter(d => d.status === filterStatus)

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      'pending': 'bg-gray-100 text-gray-800',
      'generated': 'bg-blue-100 text-blue-800',
      'sent': 'bg-yellow-100 text-yellow-800',
      'received': 'bg-green-100 text-green-800',
      'executed': 'bg-purple-100 text-purple-800',
      'registered': 'bg-green-100 text-green-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-600">Loading documents...</div>
  }

  return (
    <div className="space-y-4">
      <select
        value={filterStatus}
        onChange={(e) => setFilterStatus(e.target.value)}
        className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="all">All Documents</option>
        <option value="pending">Pending</option>
        <option value="generated">Generated</option>
        <option value="sent">Sent</option>
        <option value="received">Received</option>
        <option value="executed">Executed</option>
        <option value="registered">Registered</option>
      </select>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {filteredDocs.map((doc) => (
          <div key={doc.id} className="bg-white rounded-lg shadow p-4 border border-gray-200">
            <div className="flex justify-between items-start mb-2">
              <h3 className="font-semibold text-gray-900 text-sm">{doc.document_type.replace(/_/g, ' ')}</h3>
              <span className={`px-2 py-1 text-xs font-medium rounded ${getStatusBadge(doc.status)}`}>
                {doc.status}
              </span>
            </div>
            <p className="text-xs text-gray-600 mb-3">{doc.document_name}</p>
            <div className="grid grid-cols-2 gap-2 mb-3 text-xs">
              {doc.issued_date && (
                <div>
                  <p className="text-gray-600">Issued</p>
                  <p className="font-medium">{new Date(doc.issued_date).toLocaleDateString()}</p>
                </div>
              )}
              {doc.received_date && (
                <div>
                  <p className="text-gray-600">Received</p>
                  <p className="font-medium">{new Date(doc.received_date).toLocaleDateString()}</p>
                </div>
              )}
              {doc.due_date && (
                <div>
                  <p className="text-gray-600">Due</p>
                  <p className="font-medium text-orange-600">{new Date(doc.due_date).toLocaleDateString()}</p>
                </div>
              )}
            </div>
            <div className="flex gap-2">
              <select
                value={doc.status}
                onChange={(e) => onStatusChange(doc, e.target.value)}
                className="flex-1 px-3 py-2 text-xs font-medium border border-gray-300 rounded hover:bg-gray-50 transition"
              >
                <option value="pending">Pending</option>
                <option value="generated">Generated</option>
                <option value="sent">Sent</option>
                <option value="received">Received</option>
                <option value="executed">Executed</option>
                <option value="registered">Registered</option>
              </select>
              <button
                onClick={() => onEdit(doc)}
                className="px-3 py-2 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition"
              >
                Edit
              </button>
            </div>
          </div>
        ))}
      </div>

      {filteredDocs.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No documents found</p>
        </div>
      )}
    </div>
  )
}
