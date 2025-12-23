export interface User {
  id: number
  email: string
  name: string
  role: string
  tenant_id: string
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  token: string
  user: User
  message: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  role: string
  tenant_id: string
  name?: string
}

// ============================================================================
// SALES MODULE TYPES
// ============================================================================

export interface SalesLead {
  id: string
  tenant_id: string
  lead_code: string
  first_name: string
  last_name: string
  email: string
  phone: string
  company_name?: string
  industry?: string
  status: string // new, contacted, qualified, negotiation, converted, lost
  detailed_status?: string
  pipeline_stage?: string
  probability?: number
  source: string // website, email, phone, referral, event, social
  campaign_id?: string
  assigned_to?: string
  assigned_date?: string
  converted_to_customer: boolean
  customer_id?: string
  next_action_date?: string
  next_action_notes?: string
  capture_date_a?: string
  capture_date_b?: string
  capture_date_c?: string
  capture_date_d?: string
  last_status_change?: string
  created_by?: string
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface CreateLeadRequest {
  first_name: string
  last_name: string
  email: string
  phone: string
  company_name?: string
  industry?: string
  status?: string
  detailed_status?: string
  pipeline_stage?: string
  probability?: number
  source: string
  assigned_to?: string
  next_action_date?: string
  next_action_notes?: string
}

export interface UpdateLeadRequest {
  first_name?: string
  last_name?: string
  email?: string
  phone?: string
  company_name?: string
  industry?: string
  status?: string
  detailed_status?: string
  pipeline_stage?: string
  probability?: number
  source?: string
  assigned_to?: string
  next_action_date?: string
  next_action_notes?: string
}

export interface UpdateLeadStatusRequest {
  status: string
  detailed_status?: string
  pipeline_stage?: string
  notes?: string
  capture_date?: string
  capture_date_type?: string // a, b, c, or d
}

export interface LeadStatusLog {
  id: string
  tenant_id: string
  lead_id: string
  old_status: string
  new_status: string
  old_pipeline_stage?: string
  new_pipeline_stage?: string
  changed_by?: string
  change_reason?: string
  capture_date_type?: string
  created_at: string
}

export interface LeadPipelineConfig {
  id: string
  tenant_id: string
  status: string
  pipeline_stage: string
  phase?: string
  color_code?: string
  icon?: string
  description?: string
  is_active: boolean
  sort_order: number
  created_at: string
  updated_at: string
}
