export interface User {
  id: number
  email: string
  role: 'admin' | 'agent' | 'supervisor' | 'user'
  tenant_id: string
  first_name?: string
  last_name?: string
  created_at?: string
  updated_at?: string
}

export interface Agent extends User {
  status: 'active' | 'inactive'
  availability: 'online' | 'offline' | 'busy'
  skills: string[]
  max_concurrent_calls: number
  current_calls: number
  total_calls: number
  avg_handle_time: number
  satisfaction_score: number
  last_active: string
}

export interface AuthResponse {
  token: string
  user: User
  message?: string
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
}

export interface Lead {
  id: number
  tenant_id: string
  first_name: string
  last_name: string
  email: string
  phone: string
  status: 'new' | 'contacted' | 'converted' | 'lost'
  created_at: string
  updated_at: string
}

export interface Call {
  id: number
  tenant_id: string
  agent_id: number
  lead_id: number
  duration: number
  status: 'ongoing' | 'completed' | 'missed' | 'failed'
  recording_url?: string
  notes?: string
  created_at: string
  ended_at?: string
}

export interface Campaign {
  id: number
  tenant_id: string
  name: string
  status: 'active' | 'paused' | 'completed'
  target_leads: number
  completed_calls: number
  conversion_rate: number
  created_at: string
}

export interface DashboardStats {
  total_agents: number
  online_agents: number
  total_calls_today: number
  average_handle_time: number
  call_completion_rate: number
  customer_satisfaction: number
  leads_in_queue: number
  revenue_today: number
}
