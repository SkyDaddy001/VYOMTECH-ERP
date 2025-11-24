// Workflow Types for Phase 3B

export type WorkflowStatus = 'active' | 'inactive' | 'draft';
export type TriggerType = 'lead_created' | 'lead_scored' | 'task_completed' | 'custom_event';
export type ActionType = 'send_email' | 'send_sms' | 'create_task' | 'update_lead' | 'send_notification' | 'add_tag';
export type InstanceStatus = 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';
export type ScheduledTaskType = 'workflow' | 'action' | 'report' | 'cleanup';
export type OperatorType = 'equals' | 'greater_than' | 'less_than' | 'contains' | 'in_list';

// Core Models
export interface WorkflowDefinition {
  id: string;
  tenant_id: string;
  name: string;
  description?: string;
  status: WorkflowStatus;
  is_enabled: boolean;
  created_at: string;
  updated_at: string;
  triggers?: WorkflowTrigger[];
  actions?: WorkflowAction[];
}

export interface WorkflowTrigger {
  id: string;
  workflow_id: string;
  trigger_type: TriggerType;
  trigger_config: Record<string, any>;
  conditions?: TriggerCondition[];
  created_at: string;
  updated_at: string;
}

export interface TriggerCondition {
  id: string;
  trigger_id: string;
  field_name: string;
  operator: OperatorType;
  value: string;
  logic_operator?: 'AND' | 'OR';
}

export interface WorkflowAction {
  id: string;
  workflow_id: string;
  action_type: ActionType;
  action_config: Record<string, any>;
  action_order: number;
  delay_seconds?: number;
  created_at: string;
  updated_at: string;
}

export interface WorkflowInstance {
  id: string;
  workflow_id: string;
  tenant_id: string;
  status: InstanceStatus;
  progress_percentage: number;
  executed_actions: number;
  failed_actions: number;
  started_at?: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
  action_executions?: WorkflowActionExecution[];
}

export interface WorkflowActionExecution {
  id: string;
  workflow_instance_id: string;
  action_id: string;
  status: InstanceStatus;
  retry_count: number;
  error_message?: string;
  result_data?: Record<string, any>;
  executed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface ScheduledTask {
  id: string;
  tenant_id: string;
  task_name: string;
  task_type: ScheduledTaskType;
  cron_expression: string;
  is_enabled: boolean;
  max_retries: number;
  last_run_at?: string;
  next_run_at?: string;
  created_at: string;
  updated_at: string;
}

export interface ScheduledTaskExecution {
  id: string;
  scheduled_task_id: string;
  status: 'success' | 'failed' | 'skipped';
  duration_ms: number;
  output?: Record<string, any>;
  error_message?: string;
  executed_at: string;
  created_at: string;
}

export interface WorkflowTemplate {
  id: string;
  tenant_id: string;
  name: string;
  description?: string;
  category: 'sales' | 'support' | 'onboarding' | 'operations';
  workflow_definition: WorkflowDefinition;
  is_public: boolean;
  created_by: string;
  created_at: string;
  updated_at: string;
}

// Request/Response DTOs
export interface CreateWorkflowRequest {
  name: string;
  description?: string;
  triggers: CreateTriggerRequest[];
  actions: CreateActionRequest[];
}

export interface CreateTriggerRequest {
  trigger_type: TriggerType;
  trigger_config: Record<string, any>;
  conditions?: CreateConditionRequest[];
}

export interface CreateConditionRequest {
  field_name: string;
  operator: OperatorType;
  value: string;
  logic_operator?: 'AND' | 'OR';
}

export interface CreateActionRequest {
  action_type: ActionType;
  action_config: Record<string, any>;
  action_order: number;
  delay_seconds?: number;
}

export interface UpdateWorkflowRequest {
  name?: string;
  description?: string;
  status?: WorkflowStatus;
  is_enabled?: boolean;
  triggers?: CreateTriggerRequest[];
  actions?: CreateActionRequest[];
}

export interface WorkflowInstanceRequest {
  workflow_id: string;
  trigger_data?: Record<string, any>;
}

export interface WorkflowResponse {
  success: boolean;
  data?: any;
  error?: string;
  message?: string;
}

// UI State Models
export interface WorkflowFormState {
  name: string;
  description: string;
  triggers: WorkflowTrigger[];
  actions: WorkflowAction[];
  errors: Record<string, string>;
}

export interface WorkflowFilter {
  status?: WorkflowStatus;
  searchTerm?: string;
  sortBy?: 'name' | 'created_at' | 'updated_at';
  sortOrder?: 'asc' | 'desc';
  page?: number;
  limit?: number;
}

export interface WorkflowStats {
  totalWorkflows: number;
  activeWorkflows: number;
  totalInstances: number;
  completedInstances: number;
  failedInstances: number;
  averageExecutionTime: number;
}
