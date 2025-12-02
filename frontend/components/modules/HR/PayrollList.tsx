'use client'

import { Payroll, Employee } from '@/types/hr'

interface PayrollListProps {
  payrolls: Payroll[]
  employees: Employee[]
  loading: boolean
  onProcess: (payrollId: string) => Promise<void>
}

export default function PayrollList({ payrolls, employees, loading, onProcess }: PayrollListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading payroll...</p>
        </div>
      </div>
    )
  }

  if (payrolls.length === 0) {
    return (
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
        <p className="text-gray-600 mb-4">No payroll records yet.</p>
      </div>
    )
  }

  const processedCount = payrolls.filter((p) => p.status === 'processed').length
  const totalGross = payrolls.reduce((sum, p) => sum + (p.basic_salary || 0) + (p.allowances || 0), 0)
  const totalDeductions = payrolls.reduce((sum, p) => sum + (p.deductions || 0), 0)
  const totalNet = payrolls.reduce((sum, p) => sum + (p.net_salary || 0), 0)

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total Records</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{payrolls.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Processed</p>
          <p className="text-2xl font-bold text-green-600 mt-1">{processedCount}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total Gross</p>
          <p className="text-2xl font-bold text-blue-600 mt-1">₹{totalGross.toLocaleString()}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
          <p className="text-gray-600 text-xs font-medium">Total Net</p>
          <p className="text-2xl font-bold text-purple-600 mt-1">₹{totalNet.toLocaleString()}</p>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white rounded-lg shadow overflow-x-auto border border-gray-200">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Employee</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Period</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Basic</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Allowances</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Deductions</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Net</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {payrolls.map((payroll) => {
              const emp = employees.find((e) => e.id === payroll.employee_id)
              return (
                <tr key={payroll.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {emp ? `${emp.first_name} ${emp.last_name}` : 'Unknown'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{payroll.payroll_period}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    ₹{(payroll.basic_salary || 0).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    ₹{(payroll.allowances || 0).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-red-600">
                    ₹{(payroll.deductions || 0).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-900">
                    ₹{(payroll.net_salary || 0).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <span
                      className={`px-2 py-1 rounded text-xs font-medium ${
                        payroll.status === 'processed'
                          ? 'bg-green-100 text-green-800'
                          : payroll.status === 'pending'
                          ? 'bg-yellow-100 text-yellow-800'
                          : 'bg-red-100 text-red-800'
                      }`}
                    >
                      {payroll.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    {payroll.status === 'pending' && (
                      <button
                        onClick={() => onProcess(payroll.id || '')}
                        className="text-blue-600 hover:text-blue-900 font-medium"
                      >
                        Process
                      </button>
                    )}
                    {payroll.status === 'processed' && <span className="text-gray-400 text-xs">Done</span>}
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
