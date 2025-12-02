'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import EmployeeList from '@/components/modules/HR/EmployeeList'
import EmployeeForm from '@/components/modules/HR/EmployeeForm'
import DepartmentList from '@/components/modules/HR/DepartmentList'
import DepartmentForm from '@/components/modules/HR/DepartmentForm'
import AttendanceList from '@/components/modules/HR/AttendanceList'
import LeaveRequest from '@/components/modules/HR/LeaveRequest'
import PayrollList from '@/components/modules/HR/PayrollList'
import { hrService } from '@/services/hr.service'
import { Employee, Department, Attendance, Leave, Payroll } from '@/types/hr'

type TabType = 'employees' | 'departments' | 'attendance' | 'leaves' | 'payroll'
type FormType = 'employee' | 'department' | 'leave' | null

export default function HRPage() {
  const [activeTab, setActiveTab] = useState<TabType>('employees')
  const [showForm, setShowForm] = useState(false)
  const [editingItem, setEditingItem] = useState<any>(null)
  const [formType, setFormType] = useState<FormType>(null)

  // Data states
  const [employees, setEmployees] = useState<Employee[]>([])
  const [departments, setDepartments] = useState<Department[]>([])
  const [attendances, setAttendances] = useState<Attendance[]>([])
  const [leaves, setLeaves] = useState<Leave[]>([])
  const [payrolls, setPayrolls] = useState<Payroll[]>([])

  // Loading states
  const [employeesLoading, setEmployeesLoading] = useState(false)
  const [departmentsLoading, setDepartmentsLoading] = useState(false)
  const [attendancesLoading, setAttendancesLoading] = useState(false)
  const [leavesLoading, setLeavesLoading] = useState(false)
  const [payrollsLoading, setPayrollsLoading] = useState(false)

  // Load employees
  const loadEmployees = async () => {
    setEmployeesLoading(true)
    try {
      const data = await hrService.getEmployees()
      setEmployees(data)
    } catch (error) {
      toast.error('Failed to load employees')
    } finally {
      setEmployeesLoading(false)
    }
  }

  // Load departments
  const loadDepartments = async () => {
    setDepartmentsLoading(true)
    try {
      const data = await hrService.getDepartments()
      setDepartments(data)
    } catch (error) {
      toast.error('Failed to load departments')
    } finally {
      setDepartmentsLoading(false)
    }
  }

  // Load attendance
  const loadAttendance = async () => {
    setAttendancesLoading(true)
    try {
      const data = await hrService.getAttendance()
      setAttendances(data)
    } catch (error) {
      toast.error('Failed to load attendance')
    } finally {
      setAttendancesLoading(false)
    }
  }

  // Load leaves
  const loadLeaves = async () => {
    setLeavesLoading(true)
    try {
      const data = await hrService.getLeaves()
      setLeaves(data)
    } catch (error) {
      toast.error('Failed to load leaves')
    } finally {
      setLeavesLoading(false)
    }
  }

  // Load payroll
  const loadPayroll = async () => {
    setPayrollsLoading(true)
    try {
      const data = await hrService.getPayrolls()
      setPayrolls(data)
    } catch (error) {
      toast.error('Failed to load payroll')
    } finally {
      setPayrollsLoading(false)
    }
  }

  // Load data on tab change
  useEffect(() => {
    switch (activeTab) {
      case 'employees':
        loadEmployees()
        break
      case 'departments':
        loadDepartments()
        break
      case 'attendance':
        loadAttendance()
        break
      case 'leaves':
        loadLeaves()
        break
      case 'payroll':
        loadPayroll()
        break
    }
  }, [activeTab])

  // Employee CRUD
  const handleCreateEmployee = () => {
    setEditingItem(null)
    setFormType('employee')
    setShowForm(true)
  }

  const handleEditEmployee = (employee: Employee) => {
    setEditingItem(employee)
    setFormType('employee')
    setShowForm(true)
  }

  const handleDeleteEmployee = async (employee: Employee) => {
    if (!confirm('Are you sure?')) return
    try {
      await hrService.deleteEmployee(employee.id || '')
      toast.success('Employee deleted!')
      loadEmployees()
    } catch (error) {
      toast.error('Failed to delete employee')
    }
  }

  const handleSubmitEmployee = async (data: Partial<Employee>) => {
    try {
      if (editingItem) {
        await hrService.updateEmployee(editingItem.id, data)
      } else {
        await hrService.createEmployee(data)
      }
      setShowForm(false)
      loadEmployees()
    } catch (error) {
      throw error
    }
  }

  // Department CRUD
  const handleCreateDepartment = () => {
    setEditingItem(null)
    setFormType('department')
    setShowForm(true)
  }

  const handleEditDepartment = (department: Department) => {
    setEditingItem(department)
    setFormType('department')
    setShowForm(true)
  }

  const handleDeleteDepartment = async (department: Department) => {
    if (!confirm('Are you sure?')) return
    try {
      await hrService.deleteDepartment(department.id || '')
      toast.success('Department deleted!')
      loadDepartments()
    } catch (error) {
      toast.error('Failed to delete department')
    }
  }

  const handleSubmitDepartment = async (data: Partial<Department>) => {
    try {
      if (editingItem) {
        await hrService.updateDepartment(editingItem.id, data)
      } else {
        await hrService.createDepartment(data)
      }
      setShowForm(false)
      loadDepartments()
    } catch (error) {
      throw error
    }
  }

  // Attendance
  const handleMarkAttendance = async (employeeId: string, status: 'present' | 'absent' | 'half_day' | 'sick_leave') => {
    try {
      await hrService.markAttendance({ employee_id: employeeId, date: new Date().toISOString(), status })
      toast.success('Attendance marked!')
      loadAttendance()
    } catch (error) {
      toast.error('Failed to mark attendance')
    }
  }

  // Leave
  const handleSubmitLeave = async (data: Partial<Leave>) => {
    try {
      await hrService.createLeave(data)
      setShowForm(false)
      loadLeaves()
    } catch (error) {
      throw error
    }
  }

  // Payroll
  const handleProcessPayroll = async (payrollId: string) => {
    try {
      await hrService.processPayroll(payrollId)
      toast.success('Payroll processed!')
      loadPayroll()
    } catch (error) {
      toast.error('Failed to process payroll')
    }
  }

  const closeForm = () => {
    setShowForm(false)
    setEditingItem(null)
    setFormType(null)
  }

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">HR Management</h1>
          <p className="text-gray-600">Manage employees, departments, attendance, leaves, and payroll</p>
        </div>

        {/* Tabs */}
        <div className="bg-white rounded-lg shadow mb-6 border border-gray-200">
          <div className="flex flex-wrap border-b border-gray-200">
            <button
              onClick={() => {
                setActiveTab('employees')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'employees'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Employees
            </button>
            <button
              onClick={() => {
                setActiveTab('departments')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'departments'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Departments
            </button>
            <button
              onClick={() => {
                setActiveTab('attendance')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'attendance'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Attendance
            </button>
            <button
              onClick={() => {
                setActiveTab('leaves')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'leaves'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Leaves
            </button>
            <button
              onClick={() => {
                setActiveTab('payroll')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'payroll'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Payroll
            </button>
          </div>

          {/* Tab Content */}
          <div className="p-6">
            {/* Employees Tab */}
            {activeTab === 'employees' && (
              <div>
                {!showForm ? (
                  <div>
                    <button
                      onClick={handleCreateEmployee}
                      className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                    >
                      + Add Employee
                    </button>
                    <EmployeeList
                      employees={employees}
                      loading={employeesLoading}
                      onEdit={handleEditEmployee}
                      onDelete={handleDeleteEmployee}
                    />
                  </div>
                ) : (
                  <EmployeeForm
                    employee={editingItem}
                    departments={departments}
                    onSubmit={handleSubmitEmployee}
                    onCancel={closeForm}
                  />
                )}
              </div>
            )}

            {/* Departments Tab */}
            {activeTab === 'departments' && (
              <div>
                {!showForm ? (
                  <div>
                    <button
                      onClick={handleCreateDepartment}
                      className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                    >
                      + Add Department
                    </button>
                    <DepartmentList
                      departments={departments}
                      loading={departmentsLoading}
                      onEdit={handleEditDepartment}
                      onDelete={handleDeleteDepartment}
                    />
                  </div>
                ) : (
                  <DepartmentForm
                    department={editingItem}
                    employees={employees}
                    onSubmit={handleSubmitDepartment}
                    onCancel={closeForm}
                  />
                )}
              </div>
            )}

            {/* Attendance Tab */}
            {activeTab === 'attendance' && (
              <AttendanceList
                attendances={attendances}
                employees={employees}
                loading={attendancesLoading}
                onMarkAttendance={handleMarkAttendance}
              />
            )}

            {/* Leaves Tab */}
            {activeTab === 'leaves' && (
              <div>
                {!showForm ? (
                  <div>
                    <button
                      onClick={() => {
                        setFormType('leave')
                        setShowForm(true)
                      }}
                      className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                    >
                      + Request Leave
                    </button>
                    <div className="space-y-4">
                      {leaves.length === 0 ? (
                        <div className="bg-blue-50 border border-blue-200 rounded-lg p-8 text-center">
                          <p className="text-gray-600">No leave requests yet.</p>
                        </div>
                      ) : (
                        <div className="bg-white rounded-lg shadow border border-gray-200 overflow-x-auto">
                          <table className="min-w-full divide-y divide-gray-200">
                            <thead className="bg-gray-50">
                              <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                  From
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">To</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                  Type
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                  Days
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                  Status
                                </th>
                              </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200">
                              {leaves.map((leave) => (
                                <tr key={leave.id}>
                                  <td className="px-6 py-4 text-sm text-gray-900">
                                    {new Date(leave.from_date).toLocaleDateString()}
                                  </td>
                                  <td className="px-6 py-4 text-sm text-gray-900">
                                    {new Date(leave.to_date).toLocaleDateString()}
                                  </td>
                                  <td className="px-6 py-4 text-sm text-gray-600">{leave.leave_type}</td>
                                  <td className="px-6 py-4 text-sm text-gray-900">{leave.number_of_days}</td>
                                  <td className="px-6 py-4 text-sm">
                                    <span
                                      className={`px-2 py-1 rounded text-xs font-medium ${
                                        leave.status === 'approved'
                                          ? 'bg-green-100 text-green-800'
                                          : leave.status === 'rejected'
                                          ? 'bg-red-100 text-red-800'
                                          : 'bg-yellow-100 text-yellow-800'
                                      }`}
                                    >
                                      {leave.status}
                                    </span>
                                  </td>
                                </tr>
                              ))}
                            </tbody>
                          </table>
                        </div>
                      )}
                    </div>
                  </div>
                ) : (
                  <LeaveRequest onSubmit={handleSubmitLeave} onCancel={closeForm} />
                )}
              </div>
            )}

            {/* Payroll Tab */}
            {activeTab === 'payroll' && (
              <PayrollList
                payrolls={payrolls}
                employees={employees}
                loading={payrollsLoading}
                onProcess={handleProcessPayroll}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
