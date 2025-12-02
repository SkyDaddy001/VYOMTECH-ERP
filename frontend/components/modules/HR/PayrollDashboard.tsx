'use client'

import { useState } from 'react'
import { Employee, Payroll } from '@/types/hr'

interface PayrollDashboardProps {
  employees?: Employee[]
  payrolls?: Payroll[]
  onSelectEmployee?: (employee: Employee) => void
}

export default function PayrollDashboard({ employees = [], payrolls = [], onSelectEmployee }: PayrollDashboardProps) {
  const [filterStatus, setFilterStatus] = useState<string>('all')
  const [filterMonth, setFilterMonth] = useState<string>(new Date().toISOString().slice(0, 7))
  const [searchTerm, setSearchTerm] = useState('')

  const filteredPayrolls = payrolls.filter(p => {
    const monthMatch = p.payroll_period === filterMonth || filterMonth === 'all'
    const statusMatch = filterStatus === 'all' || p.status === filterStatus
    const employeeInfo = employees.find(e => e.id === p.employee_id)
    const employeeName = employeeInfo ? `${employeeInfo.first_name} ${employeeInfo.last_name}` : ''
    const searchMatch = !searchTerm || 
      employeeName.toLowerCase().includes(searchTerm.toLowerCase()) ||
      p.employee_id?.includes(searchTerm)
    return monthMatch && statusMatch && searchMatch
  })

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'paid':
        return 'bg-green-100 text-green-800'
      case 'pending':
        return 'bg-yellow-100 text-yellow-800'
      case 'processed':
        return 'bg-blue-100 text-blue-800'
      case 'draft':
        return 'bg-gray-100 text-gray-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getEmployeeName = (employeeId: string) => {
    const emp = employees.find(e => e.id === employeeId)
    return emp ? `${emp.first_name} ${emp.last_name}` : employeeId
  }

  const calculateStats = () => {
    const filteredForMonth = payrolls.filter(p => p.payroll_period === filterMonth)
    return {
      totalEmployees: new Set(filteredForMonth.map(p => p.employee_id)).size,
      totalSalary: filteredForMonth.reduce((sum, p) => sum + (p.basic_salary || 0), 0),
      totalDeductions: filteredForMonth.reduce((sum, p) => sum + (p.deductions || 0), 0),
      totalNet: filteredForMonth.reduce((sum, p) => sum + (p.net_salary || 0), 0),
      paidCount: filteredForMonth.filter(p => p.status === 'paid').length,
      pendingCount: filteredForMonth.filter(p => p.status === 'pending').length,
    }
  }

  const stats = calculateStats()

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6 space-y-6">
      <div>
        <h2 className="text-2xl font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <span className="text-3xl">💼</span> Payroll Management
        </h2>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-cols-6 gap-3">
        <div className="bg-blue-50 rounded-lg p-3 border border-blue-200">
          <p className="text-xs text-blue-600 mb-1">Total Employees</p>
          <p className="text-2xl font-bold text-blue-900">{stats.totalEmployees}</p>
        </div>
        <div className="bg-green-50 rounded-lg p-3 border border-green-200">
          <p className="text-xs text-green-600 mb-1">Total Salary</p>
          <p className="text-lg font-bold text-green-900">₹{(stats.totalSalary / 100000).toFixed(1)}L</p>
        </div>
        <div className="bg-orange-50 rounded-lg p-3 border border-orange-200">
          <p className="text-xs text-orange-600 mb-1">Deductions</p>
          <p className="text-lg font-bold text-orange-900">₹{(stats.totalDeductions / 100000).toFixed(1)}L</p>
        </div>
        <div className="bg-purple-50 rounded-lg p-3 border border-purple-200">
          <p className="text-xs text-purple-600 mb-1">Net Salary</p>
          <p className="text-lg font-bold text-purple-900">₹{(stats.totalNet / 100000).toFixed(1)}L</p>
        </div>
        <div className="bg-green-100 rounded-lg p-3 border border-green-300">
          <p className="text-xs text-green-700 mb-1">Paid</p>
          <p className="text-2xl font-bold text-green-700">{stats.paidCount}</p>
        </div>
        <div className="bg-yellow-100 rounded-lg p-3 border border-yellow-300">
          <p className="text-xs text-yellow-700 mb-1">Pending</p>
          <p className="text-2xl font-bold text-yellow-700">{stats.pendingCount}</p>
        </div>
      </div>

      {/* Filters */}
      <div className="flex gap-4 items-end">
        <div className="flex-1">
          <label className="block text-sm font-medium text-gray-700 mb-2">Search Employee</label>
          <input
            type="text"
            placeholder="Search by name or ID..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Month</label>
          <input
            type="month"
            value={filterMonth}
            onChange={(e) => setFilterMonth(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">All Status</option>
            <option value="pending">Pending</option>
            <option value="processed">Processed</option>
            <option value="paid">Paid</option>
            <option value="draft">Draft</option>
          </select>
        </div>
      </div>

      {/* Payroll Table */}
      {filteredPayrolls.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          <p>No payroll records found</p>
        </div>
      ) : (
        <div className="overflow-x-auto border border-gray-200 rounded-lg">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Employee</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Period</th>
                <th className="px-4 py-3 text-right text-sm font-semibold text-gray-700">Basic Salary</th>
                <th className="px-4 py-3 text-right text-sm font-semibold text-gray-700">Allowances</th>
                <th className="px-4 py-3 text-right text-sm font-semibold text-gray-700">Deductions</th>
                <th className="px-4 py-3 text-right text-sm font-semibold text-gray-700">Net Salary</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Status</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {filteredPayrolls.map((payroll) => (
                <tr key={payroll.id} className="hover:bg-gray-50 transition">
                  <td className="px-4 py-3">
                    <div className="font-medium text-gray-900">{getEmployeeName(payroll.employee_id || '')}</div>
                    <div className="text-xs text-gray-600">{payroll.employee_id}</div>
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-600">{payroll.payroll_period}</td>
                  <td className="px-4 py-3 text-right text-sm font-semibold text-gray-900">
                    ₹{((payroll.basic_salary || 0) * 100 / 100).toLocaleString()}
                  </td>
                  <td className="px-4 py-3 text-right text-sm font-semibold text-gray-900">
                    ₹{((payroll.allowances || 0) * 100 / 100).toLocaleString()}
                  </td>
                  <td className="px-4 py-3 text-right text-sm font-semibold text-red-600">
                    -₹{((payroll.deductions || 0) * 100 / 100).toLocaleString()}
                  </td>
                  <td className="px-4 py-3 text-right text-sm font-bold text-green-600">
                    ₹{((payroll.net_salary || 0) * 100 / 100).toLocaleString()}
                  </td>
                  <td className="px-4 py-3">
                    <span className={`inline-block px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(payroll.status)}`}>
                      {payroll.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm space-x-2">
                    <button className="text-blue-600 hover:text-blue-900 font-medium">View</button>
                    {payroll.status === 'draft' && (
                      <>
                        <button className="text-green-600 hover:text-green-900 font-medium">Process</button>
                        <button className="text-red-600 hover:text-red-900 font-medium">Reject</button>
                      </>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Salary Breakdown Chart Placeholder */}
      <div className="bg-gray-50 rounded-lg p-4 border border-gray-200">
        <h3 className="font-semibold text-gray-800 mb-4">Salary Structure Breakdown</h3>
        <div className="grid grid-cols-4 gap-4">
          <div className="text-center">
            <div className="text-2xl font-bold text-green-600">{Math.round((stats.totalSalary / (stats.totalSalary + stats.totalDeductions)) * 100)}%</div>
            <p className="text-sm text-gray-600">Net Salary</p>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-red-600">{Math.round((stats.totalDeductions / (stats.totalSalary + stats.totalDeductions)) * 100)}%</div>
            <p className="text-sm text-gray-600">Deductions</p>
          </div>
          <div className="text-center">
            <div className="text-sm text-gray-600">Total Employees</div>
            <p className="text-2xl font-bold text-blue-600">{stats.totalEmployees}</p>
          </div>
          <div className="text-center">
            <div className="text-sm text-gray-600">Avg Salary</div>
            <p className="text-2xl font-bold text-purple-600">₹{Math.round(stats.totalNet / (stats.totalEmployees || 1)).toLocaleString()}</p>
          </div>
        </div>
      </div>
    </div>
  )
}
