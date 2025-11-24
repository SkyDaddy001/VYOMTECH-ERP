'use client'

import { useAuth } from '@/hooks/useAuth'
import DashboardLayout from '@/components/layouts/DashboardLayout'

export default function AgentsPage() {
  const { user } = useAuth()

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <h1 className="text-3xl font-bold text-gray-800">Agents Management</h1>

        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-bold text-gray-800">All Agents</h2>
            <button className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition">
              ➕ Add New Agent
            </button>
          </div>

          <table className="w-full text-left">
            <thead className="border-b-2 border-gray-200">
              <tr>
                <th className="pb-3 font-semibold text-gray-700">Name</th>
                <th className="pb-3 font-semibold text-gray-700">Email</th>
                <th className="pb-3 font-semibold text-gray-700">Status</th>
                <th className="pb-3 font-semibold text-gray-700">Calls Today</th>
                <th className="pb-3 font-semibold text-gray-700">Actions</th>
              </tr>
            </thead>
            <tbody>
              {[1, 2, 3, 4, 5].map((i) => (
                <tr key={i} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-4">Agent #{i}</td>
                  <td className="py-4">agent{i}@callcenter.com</td>
                  <td className="py-4">
                    <span
                      className={`px-3 py-1 rounded-full text-sm font-medium ${
                        i % 2 === 0 ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                      }`}
                    >
                      {i % 2 === 0 ? 'Online' : 'Offline'}
                    </span>
                  </td>
                  <td className="py-4">{Math.floor(Math.random() * 50) + 5}</td>
                  <td className="py-4">
                    <button className="text-blue-600 hover:text-blue-800 font-medium">
                      View
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </DashboardLayout>
  )
}
