export interface ConstructionProject {
  id: string
  project_name: string
  project_code: string
  location: string
  client: string
  contract_value: number
  start_date: string
  expected_completion: string
  current_progress_percentage: number
  status: 'planning' | 'active' | 'suspended' | 'completed' | 'on_hold'
  project_manager: string
  created_at?: string
}

export interface BillOfQuantities {
  id: string
  project_id: string
  boq_number: string
  item_description: string
  unit: string
  quantity: number
  unit_rate: number
  total_amount: number
  category: 'civil' | 'structural' | 'electrical' | 'plumbing' | 'finishing' | 'other'
  status: 'planned' | 'in_progress' | 'completed' | 'on_hold'
  created_at?: string
}

export interface ProgressTracking {
  id: string
  project_id: string
  date: string
  activity_description: string
  quantity_completed: number
  unit: string
  percentage_complete: number
  workforce_deployed: number
  notes: string
  photo_url?: string
  created_at?: string
}

export interface QualityControl {
  id: string
  project_id: string
  boq_item_id: string
  inspection_date: string
  inspector_name: string
  quality_status: 'passed' | 'failed' | 'partial' | 'pending'
  observations: string
  corrective_actions?: string
  follow_up_date?: string
  created_at?: string
}

export interface ConstructionEquipment {
  id: string
  project_id: string
  equipment_type: string
  equipment_name: string
  quantity: number
  deployment_date: string
  removal_date?: string
  operator_name: string
  status: 'active' | 'idle' | 'maintenance' | 'removed'
  utilization_hours: number
}

export interface ConstructionDashboard {
  total_projects: number
  active_projects: number
  avg_progress_percentage: number
  completed_projects: number
  boq_items_total: number
  boq_items_completed: number
  quality_pass_rate: number
  equipment_deployed: number
  workforce_deployed: number
  project_timeline_status: 'on_track' | 'at_risk' | 'delayed'
}
