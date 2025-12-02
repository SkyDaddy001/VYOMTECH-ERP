// Core Post-Sales Stages
export interface PostSalesStage {
  stage_number: number
  stage_name: string
  description: string
}

// Customer Interactions - Calls, Emails, Meetings
export interface CustomerInteraction {
  id?: string
  booking_id: string
  customer_id: string
  interaction_date: string
  interaction_type: 'call' | 'email' | 'meeting' | 'document' | 'payment' | 'query' | 'escalation'
  subject: string
  description: string
  status: 'pending' | 'resolved' | 'escalated'
  priority: 'low' | 'medium' | 'high' | 'critical'
  assigned_to?: string
  assigned_to_name?: string
  notes?: string
  follow_up_date?: string
  created_at?: string
  updated_at?: string
}

// Document Tracking - Allotment, Agreement, Tax Invoice, etc.
export interface DocumentTracker {
  id?: string
  booking_id: string
  customer_id: string
  document_type: 'allotment_letter' | 'agreement' | 'receipt' | 'tax_invoice' | 'tds_certificate' | 'oc_cc' | 'possession_letter' | 'noc'
  document_name: string
  status: 'pending' | 'generated' | 'sent' | 'received' | 'executed' | 'registered'
  issued_date?: string
  received_date?: string
  due_date?: string
  completion_date?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

// Snag List - Post-Possession Issues
export interface SnagList {
  id?: string
  booking_id: string
  property_id?: string
  snag_description: string
  snag_category: 'structural' | 'finishing' | 'utilities' | 'paintwork' | 'defect' | 'plumbing' | 'electrical' | 'other'
  severity: 'low' | 'medium' | 'high'
  reported_date: string
  target_completion_date: string
  actual_completion_date?: string
  status: 'open' | 'in_progress' | 'resolved' | 'pending_inspection'
  assigned_to?: string
  assigned_to_name?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

// Change Requests - Unit Changes, Parking, Floor
export interface ChangeRequest {
  id?: string
  booking_id: string
  customer_id: string
  crm_number: string
  request_type: 'unit_modification' | 'floor_change' | 'parking_choice' | 'amenity_upgrade' | 'specification_change' | 'other'
  description: string
  impact: 'no_cost_change' | 'additional_cost' | 'cost_reduction'
  cost_difference: number
  status: 'submitted' | 'under_review' | 'approved' | 'rejected' | 'implemented'
  request_date: string
  approval_date?: string
  completion_date?: string
  assigned_to?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

// Post-Sales KPI Dashboard
export interface PostSalesMetrics {
  total_active_bookings: number
  payment_collection_percentage: number
  avg_document_tat_hours: number
  pending_interactions: number
  resolved_interactions: number
  escalated_issues: number
  pending_snags: number
  resolved_snags: number
  nps_score: number
  agreement_signing_tat: number // in days
  possession_handovers_completed: number
  customer_satisfaction: number // 1-5
}

export interface KPIDashboard {
  period: string
  payment_collection_pct: number
  document_tat_hours: number
  snag_resolution_tat: number
  nps: number
  agreement_signing_conversion_tat: number
  possession_completion_pct: number
  customer_satisfaction: number
  interaction_resolution_rate: number
}
