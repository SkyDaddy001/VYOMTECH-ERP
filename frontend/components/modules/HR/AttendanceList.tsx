'use client'

import { Attendance, Employee } from '@/types/hr'
import { useState } from 'react'

interface AttendanceListProps {
  attendances: Attendance[]
  employees: Employee[]
  loading: boolean
  onMarkAttendance: (employeeId: string, status: 'present' | 'absent' | 'half_day' | 'sick_leave') => Promise<void>
}

export default function AttendanceList({ attendances, employees, loading, onMarkAttendance }: AttendanceListProps) {
  const [markingLoading, setMarkingLoading] = useState<string | null>(null)

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading attendance...</p>
        </div>
      </div>
    )
  }

  const presentCount = attendances.filter((a) => a.status === 'present').length
  const absentCount = attendances.filter((a) => a.status === 'absent').length
  const halfDayCount = attendances.filter((a) => a.status === 'half_day').length

  const handleMark = async (employeeId: string, status: 'present' | 'absent' | 'half_day' | 'sick_leave') => {
    setMarkingLoading(employeeId)
    try {
      await onMarkAttendance(employeeId, status)
    } finally {
      setMarkingLoading(null)
    }
  }

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total Records</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{attendances.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Present</p>
          <p className="text-2xl font-bold text-green-600 mt-1">{presentCount}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Absent</p>
          <p className="text-2xl font-bold text-red-600 mt-1">{absentCount}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Half Day</p>
          <p className="text-2xl font-bold text-yellow-600 mt-1">{halfDayCount}</p>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-lg shadow overflow-x-auto border border-gray-200">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Employee</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Check In</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Check Out</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {attendances.map((att) => {
              const emp = employees.find((e) => e.id === att.employee_id)
              return (
                <tr key={att.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {emp ? `${emp.first_name} ${emp.last_name}` : 'Unknown'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    {new Date(att.date).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <span
                      className={`px-2 py-1 rounded text-xs font-medium ${
                        att.status === 'present'
                          ? 'bg-green-100 text-green-800'
                          : att.status === 'absent'
                          ? 'bg-red-100 text-red-800'
                          : att.status === 'half_day'
                          ? 'bg-yellow-100 text-yellow-800'
                          : 'bg-purple-100 text-purple-800'
                      }`}
                    >
                      {att.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    {att.check_in ? new Date(att.check_in).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : '-'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    {att.check_out ? new Date(att.check_out).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : '-'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm space-x-1">
                    <button
                      onClick={() => handleMark(att.employee_id, 'present')}
                      disabled={markingLoading === att.employee_id}
                      className="text-green-600 hover:text-green-900 font-medium disabled:opacity-50"
                    >
                      P
                    </button>
                    <button
                      onClick={() => handleMark(att.employee_id, 'absent')}
                      disabled={markingLoading === att.employee_id}
                      className="text-red-600 hover:text-red-900 font-medium disabled:opacity-50"
                    >
                      A
                    </button>
                    <button
                      onClick={() => handleMark(att.employee_id, 'half_day')}
                      disabled={markingLoading === att.employee_id}
                      className="text-yellow-600 hover:text-yellow-900 font-medium disabled:opacity-50"
                    >
                      H
                    </button>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}
