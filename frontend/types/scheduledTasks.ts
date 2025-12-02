export interface ScheduledTask {
  id: string
  task_name: string
  task_description: string
  frequency: 'once' | 'daily' | 'weekly' | 'monthly' | 'quarterly' | 'annually'
  scheduled_date: string
  scheduled_time: string
  next_execution?: string
  last_execution?: string
  execution_status: 'pending' | 'running' | 'completed' | 'failed' | 'skipped'
  assigned_to: string
  priority: 'low' | 'medium' | 'high' | 'critical'
  category: string
  is_active: boolean
  created_at?: string
}

export interface TaskExecution {
  id: string
  task_id: string
  task_name: string
  execution_start: string
  execution_end?: string
  duration_minutes?: number
  status: 'pending' | 'running' | 'completed' | 'failed' | 'skipped'
  execution_result: string
  error_message?: string
  executed_by: string
  notes?: string
  created_at?: string
}

export interface TaskTemplate {
  id: string
  template_name: string
  description: string
  task_type: 'manual' | 'automated' | 'recurring' | 'one_time'
  default_assignee: string
  default_priority: 'low' | 'medium' | 'high' | 'critical'
  estimated_duration_minutes: number
  checklist_items?: string[]
  created_at?: string
}

export interface TaskNotification {
  id: string
  task_id: string
  task_name: string
  notification_type: 'reminder' | 'due_soon' | 'overdue' | 'completed' | 'failed'
  recipient: string
  notification_date: string
  is_read: boolean
  action_url?: string
  created_at?: string
}

export interface ScheduledTasksDashboard {
  total_tasks: number
  active_tasks: number
  pending_tasks: number
  completed_tasks_today: number
  failed_tasks: number
  overdue_tasks: number
  completion_rate_percentage: number
  average_execution_time_minutes: number
  upcoming_due_in_24hrs: number
}
