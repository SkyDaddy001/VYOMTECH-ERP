export interface Company {
  id: string
  tenant_id: string
  name: string
  description: string
  status: 'active' | 'inactive' | 'suspended'
  industry_type: string
  employee_count?: number
  website?: string
  max_projects: number
  max_users: number
  current_user_count: number
  current_project_count: number
  billing_email: string
  billing_address: string
  created_at: string
  updated_at: string
}

export interface Project {
  id: string
  company_id: string
  tenant_id: string
  name: string
  description: string
  status: 'active' | 'inactive' | 'archived'
  project_type: 'sales' | 'support' | 'marketing' | 'custom'
  max_users: number
  current_user_count: number
  budget_allocated: number
  budget_spent: number
  start_date: string
  end_date?: string
  created_at: string
  updated_at: string
}

export interface CompanyMember {
  id: string
  company_id: string
  user_id: number
  tenant_id: string
  role: 'owner' | 'admin' | 'manager' | 'member' | 'viewer'
  department: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface ProjectMember {
  id: string
  project_id: string
  user_id: number
  company_id: string
  tenant_id: string
  role: 'lead' | 'member' | 'viewer' | 'analyst'
  joined_at: string
  is_active: boolean
  created_at: string
  updated_at: string
}
