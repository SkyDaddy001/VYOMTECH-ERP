'use client'

import { useState, useMemo } from 'react'

export default function AuditLogsPage() {
  const [userFilter, setUserFilter] = useState('')
  const [actionFilter, setActionFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [dateFilter, setDateFilter] = useState('')

  const logs = [
    { id: 1, timestamp: '2025-12-04 14:32:10', user: 'master.admin@vyomtech.com', action: 'LOGIN', resource: 'Admin Dashboard', details: 'Successful login', status: 'success' },
    { id: 2, timestamp: '2025-12-04 14:25:45', user: 'demo@vyomtech.com', action: 'VIEW_TENANT', resource: 'Vyomtech Demo', details: 'Viewed tenant details', status: 'success' },
    { id: 3, timestamp: '2025-12-04 13:18:22', user: 'master.admin@vyomtech.com', action: 'UPDATE_SETTINGS', resource: 'System Settings', details: 'Updated API rate limit', status: 'success' },
    { id: 4, timestamp: '2025-12-04 12:45:33', user: 'admin@example.com', action: 'LOGIN', resource: 'Admin Dashboard', details: 'Failed login - invalid credentials', status: 'failed' },
    { id: 5, timestamp: '2025-12-04 12:30:10', user: 'master.admin@vyomtech.com', action: 'CREATE_USER', resource: 'User Management', details: 'Created new user account', status: 'success' },
    { id: 6, timestamp: '2025-12-04 11:15:20', user: 'master.admin@vyomtech.com', action: 'DELETE_TENANT', resource: 'Old Test Tenant', details: 'Deleted inactive tenant', status: 'success' },
    { id: 7, timestamp: '2025-12-04 10:45:50', user: 'vendor@demo.vyomtech.com', action: 'EXPORT_DATA', resource: 'Reports', details: 'Exported monthly report', status: 'success' },
    { id: 8, timestamp: '2025-12-04 09:30:15', user: 'master.admin@vyomtech.com', action: 'RUN_BACKUP', resource: 'Database', details: 'Full backup completed successfully', status: 'success' },
  ]

  const filteredLogs = useMemo(() => {
    return logs.filter(log => {
      const userMatch = log.user.toLowerCase().includes(userFilter.toLowerCase())
      const actionMatch = log.action.toLowerCase().includes(actionFilter.toLowerCase())
      const statusMatch = !statusFilter || log.status === statusFilter
      const dateMatch = !dateFilter || log.timestamp.startsWith(dateFilter)
      return userMatch && actionMatch && statusMatch && dateMatch
    })
  }, [userFilter, actionFilter, statusFilter, dateFilter])

  return (
    <div className="p-8 bg-gray-50 min-h-screen">
      <div className="mb-8">
        <h1 className="text-4xl font-bold text-gray-900">Audit Logs</h1>
        <p className="text-gray-600 mt-2">View complete audit trail of all system activities</p>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-lg shadow-md p-6 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-900 mb-2">User</label>
            <input
              type="text"
              placeholder="Filter by user..."
              value={userFilter}
              onChange={(e) => setUserFilter(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-900 mb-2">Action</label>
            <input
              type="text"
              placeholder="Filter by action..."
              value={actionFilter}
              onChange={(e) => setActionFilter(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-900 mb-2">Status</label>
            <select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)} className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
              <option value="">All</option>
              <option value="success">Success</option>
              <option value="failed">Failed</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-900 mb-2">Date Range</label>
            <input
              type="date"
              value={dateFilter}
              onChange={(e) => setDateFilter(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      </div>

      {/* Logs Table */}
      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-100 border-b border-gray-200">
            <tr>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-900">Timestamp</th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-900">User</th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-900">Action</th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-900">Resource</th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-900">Details</th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-900">Status</th>
            </tr>
          </thead>
          <tbody>
            {filteredLogs.map((log) => (
              <tr key={log.id} className="border-b border-gray-200 hover:bg-gray-50 transition">
                <td className="px-6 py-4 text-sm text-gray-600 font-mono">{log.timestamp}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{log.user}</td>
                <td className="px-6 py-4">
                  <span className="inline-block px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-xs font-semibold">
                    {log.action}
                  </span>
                </td>
                <td className="px-6 py-4 text-sm text-gray-900">{log.resource}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{log.details}</td>
                <td className="px-6 py-4">
                  <span className={`inline-block px-3 py-1 rounded-full text-sm font-medium ${
                    log.status === 'success'
                      ? 'bg-green-100 text-green-800'
                      : 'bg-red-100 text-red-800'
                  }`}>
                    {log.status}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      <div className="mt-6 flex justify-between items-center">
        <p className="text-sm text-gray-600">Showing {filteredLogs.length > 0 ? 1 : 0}-{Math.min(8, filteredLogs.length)} of {filteredLogs.length} logs</p>
        <div className="flex gap-2">
          <button onClick={() => alert('Pagination coming soon')} className="px-4 py-2 border border-gray-300 rounded-lg text-gray-900 hover:bg-gray-50">Previous</button>
          <button onClick={() => alert('Pagination coming soon')} className="px-4 py-2 border border-gray-300 rounded-lg text-gray-900 hover:bg-gray-50">1</button>
          <button onClick={() => alert('Pagination coming soon')} className="px-4 py-2 border border-gray-300 rounded-lg bg-blue-500 text-white">2</button>
          <button onClick={() => alert('Pagination coming soon')} className="px-4 py-2 border border-gray-300 rounded-lg text-gray-900 hover:bg-gray-50">3</button>
          <button onClick={() => alert('Pagination coming soon')} className="px-4 py-2 border border-gray-300 rounded-lg text-gray-900 hover:bg-gray-50">Next</button>
        </div>
      </div>
    </div>
  )
}
