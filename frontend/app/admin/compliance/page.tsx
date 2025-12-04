'use client'

export default function CompliancePage() {
  return (
    <div className="p-8 bg-gray-50 min-h-screen">
      <div className="mb-8">
        <h1 className="text-4xl font-bold text-gray-900">Compliance & Audit</h1>
        <p className="text-gray-600 mt-2">Monitor system compliance, audit logs, and security events</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        {/* Compliance Status */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Compliance Status</h2>
          <div className="space-y-4">
            {[
              { name: 'GDPR Compliance', status: 'compliant' },
              { name: 'Data Protection', status: 'compliant' },
              { name: 'Access Control', status: 'compliant' },
              { name: 'Encryption', status: 'compliant' },
              { name: 'Backup & Recovery', status: 'warning' }
            ].map((item, i) => (
              <div key={i} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <span className="text-gray-900 font-medium">{item.name}</span>
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                  item.status === 'compliant'
                    ? 'bg-green-100 text-green-800'
                    : 'bg-yellow-100 text-yellow-800'
                }`}>
                  {item.status}
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Recent Security Events */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Recent Security Events</h2>
          <div className="space-y-4">
            {[
              { event: 'Failed Login Attempt', time: '2 hours ago', severity: 'warning' },
              { event: 'User Created', time: '5 hours ago', severity: 'info' },
              { event: 'API Key Generated', time: '1 day ago', severity: 'info' },
              { event: 'Backup Completed', time: '1 day ago', severity: 'success' },
              { event: 'Permission Changed', time: '2 days ago', severity: 'info' }
            ].map((item, i) => (
              <div key={i} className="flex items-start p-3 bg-gray-50 rounded-lg">
                <div className={`w-2 h-2 rounded-full mt-2 mr-3 flex-shrink-0 ${
                  item.severity === 'warning' ? 'bg-red-500' :
                  item.severity === 'success' ? 'bg-green-500' :
                  'bg-blue-500'
                }`}></div>
                <div className="flex-1">
                  <p className="text-gray-900 font-medium text-sm">{item.event}</p>
                  <p className="text-gray-600 text-xs">{item.time}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Audit Trails */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-bold text-gray-900 mb-6">Audit Trail</h2>
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">Timestamp</th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">User</th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">Action</th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">Resource</th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
            </tr>
          </thead>
          <tbody>
            {[
              { time: '2025-12-04 10:30:45', user: 'master.admin@vyomtech.com', action: 'View', resource: 'Tenants', status: 'success' },
              { time: '2025-12-04 10:28:12', user: 'master.admin@vyomtech.com', action: 'Update', resource: 'Settings', status: 'success' },
              { time: '2025-12-04 10:15:30', user: 'demo@vyomtech.com', action: 'Login', resource: 'Auth', status: 'success' },
              { time: '2025-12-04 09:45:22', user: 'admin@example.com', action: 'Delete', resource: 'User', status: 'failed' }
            ].map((log, i) => (
              <tr key={i} className="border-b border-gray-200 hover:bg-gray-50">
                <td className="px-4 py-3 text-sm text-gray-600">{log.time}</td>
                <td className="px-4 py-3 text-sm text-gray-900">{log.user}</td>
                <td className="px-4 py-3 text-sm text-gray-900">{log.action}</td>
                <td className="px-4 py-3 text-sm text-gray-900">{log.resource}</td>
                <td className="px-4 py-3">
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
    </div>
  )
}
