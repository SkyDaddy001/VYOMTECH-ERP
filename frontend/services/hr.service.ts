import { apiClient } from './api'
import { Employee, Department, Attendance, Leave, Payroll } from '@/types/hr'

export const hrService = {
  // Employees
  async getEmployees(): Promise<Employee[]> {
    return apiClient.get<Employee[]>('/api/v1/employees')
  },

  async getEmployee(id: string): Promise<Employee> {
    return apiClient.get<Employee>(`/api/v1/employees/${id}`)
  },

  async createEmployee(data: Partial<Employee>): Promise<Employee> {
    return apiClient.post<Employee>('/api/v1/employees', data)
  },

  async updateEmployee(id: string, data: Partial<Employee>): Promise<Employee> {
    return apiClient.put<Employee>(`/api/v1/employees/${id}`, data)
  },

  async deleteEmployee(id: string): Promise<void> {
    return apiClient.post(`/api/v1/employees/${id}/delete`, {})
  },

  // Departments
  async getDepartments(): Promise<Department[]> {
    return apiClient.get<Department[]>('/api/v1/departments')
  },

  async getDepartment(id: string): Promise<Department> {
    return apiClient.get<Department>(`/api/v1/departments/${id}`)
  },

  async createDepartment(data: Partial<Department>): Promise<Department> {
    return apiClient.post<Department>('/api/v1/departments', data)
  },

  async updateDepartment(id: string, data: Partial<Department>): Promise<Department> {
    return apiClient.put<Department>(`/api/v1/departments/${id}`, data)
  },

  async deleteDepartment(id: string): Promise<void> {
    return apiClient.post(`/api/v1/departments/${id}/delete`, {})
  },

  // Attendance
  async getAttendance(employeeId?: string, fromDate?: string, toDate?: string): Promise<Attendance[]> {
    let url = '/api/v1/attendance'
    const params = []
    if (employeeId) params.push(`employee_id=${employeeId}`)
    if (fromDate) params.push(`from_date=${fromDate}`)
    if (toDate) params.push(`to_date=${toDate}`)
    if (params.length) url += '?' + params.join('&')
    return apiClient.get<Attendance[]>(url)
  },

  async markAttendance(data: Partial<Attendance>): Promise<Attendance> {
    return apiClient.post<Attendance>('/api/v1/attendance', data)
  },

  async updateAttendance(id: string, data: Partial<Attendance>): Promise<Attendance> {
    return apiClient.put<Attendance>(`/api/v1/attendance/${id}`, data)
  },

  // Leave
  async getLeaves(employeeId?: string, status?: string): Promise<Leave[]> {
    let url = '/api/v1/leaves'
    const params = []
    if (employeeId) params.push(`employee_id=${employeeId}`)
    if (status) params.push(`status=${status}`)
    if (params.length) url += '?' + params.join('&')
    return apiClient.get<Leave[]>(url)
  },

  async createLeave(data: Partial<Leave>): Promise<Leave> {
    return apiClient.post<Leave>('/api/v1/leaves', data)
  },

  async approveLeave(id: string): Promise<Leave> {
    return apiClient.post<Leave>(`/api/v1/leaves/${id}/approve`, {})
  },

  async rejectLeave(id: string): Promise<Leave> {
    return apiClient.post<Leave>(`/api/v1/leaves/${id}/reject`, {})
  },

  // Payroll
  async getPayrolls(employeeId?: string, period?: string): Promise<Payroll[]> {
    let url = '/api/v1/payroll'
    const params = []
    if (employeeId) params.push(`employee_id=${employeeId}`)
    if (period) params.push(`period=${period}`)
    if (params.length) url += '?' + params.join('&')
    return apiClient.get<Payroll[]>(url)
  },

  async createPayroll(data: Partial<Payroll>): Promise<Payroll> {
    return apiClient.post<Payroll>('/api/v1/payroll', data)
  },

  async processPayroll(id: string): Promise<Payroll> {
    return apiClient.post<Payroll>(`/api/v1/payroll/${id}/process`, {})
  },

  async getPayroll(id: string): Promise<Payroll> {
    return apiClient.get<Payroll>(`/api/v1/payroll/${id}`)
  },
}
