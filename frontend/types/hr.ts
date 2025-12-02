export interface Employee {
  id?: string
  tenant_id?: string
  company_id?: string
  email: string
  first_name: string
  last_name: string
  department_id?: string
  designation?: string
  employment_type: 'full_time' | 'part_time' | 'contract' | 'intern'
  status: 'active' | 'inactive' | 'on_leave' | 'terminated'
  date_of_joining?: string
  joining_date?: string
  date_of_birth?: string
  phone?: string
  address?: string
  city?: string
  state?: string
  zip_code?: string
  salary_grade?: string
  cost_center?: string
  manager_id?: string
  created_at?: string
  updated_at?: string
}

export interface Department {
  id?: string
  tenant_id?: string
  company_id?: string
  name: string
  description?: string
  head_id?: string
  budget?: number
  budget_allocated?: number
  budget_spent?: number
  employee_count?: number
  status?: 'active' | 'inactive'
  created_at?: string
  updated_at?: string
}

export interface Attendance {
  id?: string
  employee_id: string
  tenant_id?: string
  date: string
  check_in?: string
  check_out?: string
  check_in_time?: string
  check_out_time?: string
  status: 'present' | 'absent' | 'half_day' | 'sick_leave' | 'leave' | 'holiday'
  location?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface Leave {
  id?: string
  employee_id?: string
  tenant_id?: string
  leave_type: 'annual' | 'sick' | 'personal' | 'maternity' | 'paternity' | 'unpaid'
  from_date: string
  to_date: string
  start_date?: string
  end_date?: string
  number_of_days?: number
  duration_days?: number
  reason: string
  status: 'pending' | 'approved' | 'rejected' | 'cancelled'
  approved_by?: string
  approval_date?: string
  created_at?: string
  updated_at?: string
}

export interface Payroll {
  id?: string
  employee_id?: string
  tenant_id?: string
  payroll_period?: string
  period_start?: string
  period_end?: string
  basic_salary?: number
  allowances?: number
  deductions?: number
  net_salary?: number
  status: 'draft' | 'processed' | 'paid' | 'pending'
  payment_date?: string
  created_at?: string
  updated_at?: string
}

export interface HRSummary {
  total_employees: number
  active_employees: number
  on_leave: number
  departments: number
  avg_salary: number
  pending_approvals: number
}
