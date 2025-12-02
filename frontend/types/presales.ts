// Call Records & Lead Management
export interface CallRecord {
  id?: string
  call_date: string
  call_time: string
  customer_name: string
  customer_phone: string
  email?: string
  property_preference: string[] // House, Apartment, Plot
  property_status: string[] // Under Construction, Ready to Move, Prelaunch
  configuration: string[] // 2 BHK, 3 BHK, 4 BHK, etc.
  area_required: number // Sq. Ft.
  area_specified: boolean
  budget_range: string // below_50_lakhs, 50_to_100_lakhs, etc.
  budget_specified: boolean
  purchase_purpose: string // Investment, Self Occupation, Rental
  occupation_type: string // Salaried, Business, Self-Employed
  occupation_details?: string
  funding_source: string // Own Funds, Bank Finance
  interested_projects: string[] // Project names
  other_requirements?: string
  call_summary: string
  call_outcome: 'interested' | 'not_interested' | 'maybe_later' | 'follow_up_needed' | 'converted'
  follow_up_required: boolean
  follow_up_date?: string
  presales_agent_id: string
  created_at?: string
  updated_at?: string
}

export interface Lead {
  id?: string
  customer_name: string
  customer_phone: string
  email?: string
  source: 'call' | 'website' | 'referral' | 'walk_in' | 'social_media'
  status: 'new' | 'contacted' | 'qualified' | 'converted' | 'rejected' | 'inactive'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  assigned_to: string // Presales agent ID
  first_contact_date?: string
  last_contact_date?: string
  next_follow_up_date?: string
  notes?: string
  call_records?: CallRecord[]
  created_at?: string
  updated_at?: string
}

// Traditional Pre-Sales Pipeline
export interface PreSalesOpportunity {
  id?: string
  tenant_id?: string
  company_id?: string
  name: string
  description?: string
  lead_id?: string
  stage: 'prospecting' | 'qualification' | 'proposal' | 'negotiation' | 'closed_won' | 'closed_lost'
  value?: number
  probability?: number
  expected_close_date: string
  assigned_to?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface Proposal {
  id?: string
  tenant_id?: string
  opportunity_id: string
  title: string
  description?: string
  amount: number
  validity_date: string
  status: 'draft' | 'sent' | 'viewed' | 'accepted' | 'rejected'
  created_at?: string
  updated_at?: string
}

export interface PreSalesMetrics {
  total_opportunities: number
  active_opportunities: number
  total_value: number
  win_rate: number
  avg_deal_size: number
  pipeline_value: number
  close_rate: number
}

export interface Pipeline {
  stage: 'prospecting' | 'qualification' | 'proposal' | 'negotiation' | 'closed_won' | 'closed_lost'
  count: number
  value: number
  average_days: number
}

// Presales Target & Performance
export interface PresalesTarget {
  id?: string
  presales_agent_id: string
  presales_agent_name?: string
  period: string // YYYY-MM
  target_calls: number
  target_qualified_leads: number
  target_conversions: number
  achieved_calls: number
  achieved_qualified_leads: number
  achieved_conversions: number
  achievement_percentage: number
  status: 'not_started' | 'in_progress' | 'completed' | 'exceeded'
  created_at?: string
  updated_at?: string
}

export interface PresalesPerformanceMetrics {
  total_calls: number
  total_leads_generated: number
  conversion_rate: number
  average_call_duration_minutes: number
  calls_this_month: number
  leads_this_month: number
  conversions_this_month: number
  target_achievement_percentage: number
  hottest_leads_count: number
}

export interface CallTranscript {
  id?: string
  call_record_id: string
  transcript_text: string
  key_points: string[]
  sentiment: 'positive' | 'neutral' | 'negative'
  created_at?: string
}

export interface LeadScore {
  id?: string
  lead_id: string
  score: number // 0-100
  factors: {
    budget_alignment: number
    timeline_readiness: number
    property_preference_match: number
    engagement_level: number
    decision_maker: boolean
  }
  updated_at?: string
}
