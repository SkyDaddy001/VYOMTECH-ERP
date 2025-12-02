export interface Project {
  id?: string
  project_name: string
  project_code: string
  location: string
  city: string
  state: string
  postal_code: string
  project_type: 'residential' | 'commercial' | 'mixed'
  total_units: number
  total_area: number
  launch_date: string
  expected_completion: string
  actual_completion?: string
  status: 'planning' | 'approved' | 'under_construction' | 'completed' | 'stalled'
  developer_name: string
  architect_name?: string
  noc_status: 'pending' | 'approved' | 'rejected'
  noc_date?: string
  created_at?: string
  updated_at?: string
}

export interface ProjectBlock {
  id?: string
  project_id: string
  block_name: string
  block_code: string
  wing_name?: string
  total_units: number
  status: 'planning' | 'under_construction' | 'completed'
  created_at?: string
  updated_at?: string
}

export interface ProjectMilestone {
  id?: string
  project_id: string
  milestone_name: string
  milestone_type: 'foundation' | 'structure' | 'finishing' | 'handover'
  planned_date: string
  actual_date?: string
  completion_percentage: number
  status: 'pending' | 'in_progress' | 'completed' | 'delayed'
  description?: string
  created_at?: string
  updated_at?: string
}

export interface ProjectMetrics {
  total_projects: number
  active_projects: number
  completed_projects: number
  total_units_planned: number
  total_units_sold: number
  avg_completion_percentage: number
  delayed_projects: number
}

export interface ProjectTimeline {
  month: string
  milestones_planned: number
  milestones_completed: number
  units_sold: number
}
