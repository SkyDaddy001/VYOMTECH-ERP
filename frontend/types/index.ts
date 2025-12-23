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

// ============================================================================
// ROLE TEMPLATE TYPES
// ============================================================================

export interface RoleTemplate {
  id: string
  tenant_id: string
  name: string
  description?: string
  category: string
  is_system_template: boolean
  is_active: boolean
  permission_ids: string[]
  metadata?: Record<string, any>
  created_at: string
  updated_at: string
}

export interface TemplateInstance {
  id: string
  tenant_id: string
  template_id: string
  role_id: string
  created_by: number
  customizations?: Record<string, any>
  created_at: string
}

export interface CreateRoleTemplateRequest {
  name: string
  description?: string
  category: string
  is_system_template: boolean
  permission_ids: string[]
  metadata?: Record<string, any>
}

export interface CreateRoleFromTemplateRequest {
  template_id: string
  role_name: string
  customizations?: Record<string, any>
}

// ============================================================================
// BANK FINANCING MODULE TYPES
// ============================================================================

export interface BankFinancing {
  id: string
  tenant_id: string
  booking_id: string
  bank_id?: string
  loan_amount: number
  sanctioned_amount: number
  disbursed_amount: number
  outstanding_amount: number
  loan_type: string // Home Loan, Construction Loan, Bridge Loan
  interest_rate?: number
  tenure_months?: number
  emi_amount?: number
  status: string // pending, approved, sanctioned, disbursing, completed, rejected
  application_date?: string
  approval_date?: string
  sanction_date?: string
  expected_completion_date?: string
  application_ref_no?: string
  sanction_letter_url?: string
  created_by?: string
  created_at: string
  updated_by?: string
  updated_at: string
  deleted_at?: string
}

export interface BankDisbursement {
  id: string
  tenant_id: string
  financing_id: string
  disbursement_number: number
  scheduled_amount: number
  actual_amount?: number
  milestone_id?: string
  milestone_percentage?: number
  status: string // pending, released, credited, delayed, cancelled
  scheduled_date: string
  actual_date?: string
  bank_reference_no?: string
  claim_document_url?: string
  release_approval_by?: string
  release_approval_date?: string
  created_by?: string
  created_at: string
  updated_by?: string
  updated_at: string
  deleted_at?: string
}

export interface BankNOC {
  id: string
  tenant_id: string
  financing_id: string
  noc_type: string // Pre-sanction, Post-completion, Full-settlement
  noc_request_date: string
  noc_received_date?: string
  noc_document_url?: string
  noc_amount?: number
  status: string // requested, issued, expired, cancelled
  issued_by_bank?: string
  valid_till_date?: string
  remarks?: string
  created_by?: string
  created_at: string
  updated_by?: string
  updated_at: string
  deleted_at?: string
}

export interface BankCollectionTracking {
  id: string
  tenant_id: string
  financing_id: string
  collection_type: string // EMI, Prepayment, Partial, Full-Settlement
  collection_amount: number
  collection_date: string
  payment_mode?: string // Bank Transfer, Cheque, NEFT, RTGS
  payment_reference_no?: string
  emi_month?: string
  emi_number?: number
  principal_amount?: number
  interest_amount?: number
  status: string // pending, verified, credited, failed
  bank_confirmation_date?: string
  created_by?: string
  created_at: string
  updated_by?: string
  updated_at: string
  deleted_at?: string
}

export interface Bank {
  id: string
  tenant_id: string
  bank_name: string
  branch_name: string
  ifsc_code?: string
  branch_contact?: string
  branch_email?: string
  relationship_manager_name?: string
  relationship_manager_phone?: string
  relationship_manager_email?: string
  status: string // active, inactive
  created_at: string
  updated_at: string
}

export interface CreateBankFinancingRequest {
  booking_id: string
  bank_id?: string
  loan_amount: number
  sanctioned_amount: number
  loan_type: string
  interest_rate?: number
  tenure_months?: number
  application_ref_no?: string
}

export interface CreateBankDisbursementRequest {
  financing_id: string
  disbursement_number: number
  scheduled_amount: number
  scheduled_date: string
  milestone_id?: string
  milestone_percentage?: number
}

export interface CreateBankNOCRequest {
  financing_id: string
  noc_type: string
  noc_request_date: string
  noc_amount?: number
}

export interface CreateBankCollectionRequest {
  financing_id: string
  collection_type: string
  collection_amount: number
  collection_date: string
  payment_mode?: string
  payment_reference_no?: string
  emi_month?: string
  emi_number?: number
}

export interface CreateBankRequest {
  bank_name: string
  branch_name: string
  ifsc_code?: string
  branch_contact?: string
  branch_email?: string
  relationship_manager_name?: string
  relationship_manager_phone?: string
  relationship_manager_email?: string
}
