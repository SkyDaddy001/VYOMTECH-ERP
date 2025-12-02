export interface Site {
  id: string
  site_name: string
  location: string
  project_id: string
  site_manager: string
  start_date: string
  expected_end_date: string
  current_status: 'planning' | 'active' | 'paused' | 'completed' | 'closed'
  site_area_sqm: number
  workforce_count: number
  created_at?: string
  updated_at?: string
}

export interface SafetyIncident {
  id: string
  site_id: string
  incident_type: 'accident' | 'near_miss' | 'hazard' | 'violation'
  severity: 'low' | 'medium' | 'high' | 'critical'
  incident_date: string
  description: string
  reported_by: string
  status: 'open' | 'investigating' | 'resolved' | 'closed'
  incident_number: string
  created_at?: string
}

export interface Compliance {
  id: string
  site_id: string
  compliance_type: 'safety' | 'environmental' | 'labor' | 'regulatory'
  requirement: string
  due_date: string
  status: 'compliant' | 'non_compliant' | 'in_progress' | 'not_applicable'
  last_audit_date: string
  audit_result: 'pass' | 'fail' | 'pending'
  notes: string
  created_at?: string
}

export interface Permit {
  id: string
  site_id: string
  permit_type: string
  permit_number: string
  issued_date: string
  expiry_date: string
  issuing_authority: string
  status: 'active' | 'expired' | 'cancelled' | 'pending'
  document_url?: string
  created_at?: string
}

export interface CivilDashboard {
  total_sites: number
  active_sites: number
  total_incidents: number
  critical_incidents: number
  compliance_status: {
    compliant: number
    non_compliant: number
    in_progress: number
  }
  pending_permits: number
  workforce_total: number
  safety_score: number
}
