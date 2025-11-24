package models

import "time"

// ==================== WORKFLOW MODELS ====================

// WorkflowDefinition represents a workflow template
type WorkflowDefinition struct {
	ID          int64             `db:"id" json:"id"`
	TenantID    string            `db:"tenant_id" json:"tenant_id"`
	Name        string            `db:"name" json:"name"`
	Description string            `db:"description" json:"description"`
	Enabled     bool              `db:"enabled" json:"enabled"`
	Triggers    []WorkflowTrigger `json:"triggers"` // NOT in DB - loaded separately
	Actions     []WorkflowAction  `json:"actions"`  // NOT in DB - loaded separately
	CreatedBy   int64             `db:"created_by" json:"created_by"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
}

// WorkflowTrigger defines when a workflow should execute
type WorkflowTrigger struct {
	ID            int64              `db:"id" json:"id"`
	WorkflowID    int64              `db:"workflow_id" json:"workflow_id"`
	TriggerType   string             `db:"trigger_type" json:"trigger_type"`     // lead_created, lead_scored, task_completed, etc
	TriggerConfig string             `db:"trigger_config" json:"trigger_config"` // JSON string with condition details
	Conditions    []TriggerCondition `json:"conditions"`                         // NOT in DB - loaded separately
	CreatedAt     time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `db:"updated_at" json:"updated_at"`
}

// TriggerCondition represents a condition that must be met for trigger to fire
type TriggerCondition struct {
	ID        int64     `db:"id" json:"id"`
	TriggerID int64     `db:"trigger_id" json:"trigger_id"`
	Field     string    `db:"field" json:"field"`       // e.g., "lead_score", "lead_status"
	Operator  string    `db:"operator" json:"operator"` // e.g., "equals", "greater_than", "less_than", "contains"
	Value     string    `db:"value" json:"value"`       // Condition value
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// WorkflowAction represents an action to be taken when workflow triggers
type WorkflowAction struct {
	ID           int64     `db:"id" json:"id"`
	WorkflowID   int64     `db:"workflow_id" json:"workflow_id"`
	ActionType   string    `db:"action_type" json:"action_type"`     // send_email, send_sms, create_task, update_lead, etc
	ActionConfig string    `db:"action_config" json:"action_config"` // JSON string with action parameters
	Order        int       `db:"action_order" json:"order"`          // Execution order
	DelaySeconds int       `db:"delay_seconds" json:"delay_seconds"` // Delay before execution
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// WorkflowInstance represents an execution instance of a workflow
type WorkflowInstance struct {
	ID               int64      `db:"id" json:"id"`
	TenantID         string     `db:"tenant_id" json:"tenant_id"`
	WorkflowID       int64      `db:"workflow_id" json:"workflow_id"`
	TriggeredBy      string     `db:"triggered_by" json:"triggered_by"` // lead_id, task_id, etc
	TriggeredByValue string     `db:"triggered_by_value" json:"triggered_by_value"`
	Status           string     `db:"status" json:"status"` // pending, running, completed, failed, cancelled
	Progress         int        `db:"progress" json:"progress"`
	ExecutedActions  int        `db:"executed_actions" json:"executed_actions"`
	FailedActions    int        `db:"failed_actions" json:"failed_actions"`
	ErrorMessage     string     `db:"error_message" json:"error_message"`
	StartedAt        *time.Time `db:"started_at" json:"started_at"`
	CompletedAt      *time.Time `db:"completed_at" json:"completed_at"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}

// WorkflowActionExecution tracks execution of individual actions
type WorkflowActionExecution struct {
	ID           int64      `db:"id" json:"id"`
	WorkflowID   int64      `db:"workflow_id" json:"workflow_id"`
	InstanceID   int64      `db:"instance_id" json:"instance_id"`
	ActionID     int64      `db:"action_id" json:"action_id"`
	Status       string     `db:"status" json:"status"` // pending, executing, completed, failed
	Result       string     `db:"result" json:"result"` // JSON result from action
	ErrorMessage string     `db:"error_message" json:"error_message"`
	RetryCount   int        `db:"retry_count" json:"retry_count"`
	StartedAt    *time.Time `db:"started_at" json:"started_at"`
	CompletedAt  *time.Time `db:"completed_at" json:"completed_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

// ScheduledTask represents a scheduled workflow or action
type ScheduledTask struct {
	ID         int64      `db:"id" json:"id"`
	TenantID   string     `db:"tenant_id" json:"tenant_id"`
	Name       string     `db:"name" json:"name"`
	Type       string     `db:"type" json:"type"`         // workflow, action, report, cleanup
	Config     string     `db:"config" json:"config"`     // JSON configuration
	Schedule   string     `db:"schedule" json:"schedule"` // Cron expression
	LastRunAt  *time.Time `db:"last_run_at" json:"last_run_at"`
	NextRunAt  time.Time  `db:"next_run_at" json:"next_run_at"`
	Enabled    bool       `db:"enabled" json:"enabled"`
	MaxRetries int        `db:"max_retries" json:"max_retries"`
	CreatedBy  int64      `db:"created_by" json:"created_by"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

// ScheduledTaskExecution tracks execution history
type ScheduledTaskExecution struct {
	ID           int64      `db:"id" json:"id"`
	TaskID       int64      `db:"task_id" json:"task_id"`
	TenantID     string     `db:"tenant_id" json:"tenant_id"`
	Status       string     `db:"status" json:"status"` // success, failed, running, cancelled
	Output       string     `db:"output" json:"output"` // JSON output
	ErrorMessage string     `db:"error_message" json:"error_message"`
	Duration     int        `db:"duration" json:"duration"` // milliseconds
	StartedAt    time.Time  `db:"started_at" json:"started_at"`
	CompletedAt  *time.Time `db:"completed_at" json:"completed_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
}

// WorkflowTemplate predefined workflow templates
type WorkflowTemplate struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Category    string    `db:"category" json:"category"` // sales, support, onboarding, etc
	Description string    `db:"description" json:"description"`
	Definition  string    `db:"definition" json:"definition"` // JSON workflow definition
	IsPublic    bool      `db:"is_public" json:"is_public"`
	CreatedBy   int64     `db:"created_by" json:"created_by"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// WorkflowRequest request to create/update workflow
type WorkflowRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Enabled     bool              `json:"enabled"`
	Triggers    []WorkflowTrigger `json:"triggers"`
	Actions     []WorkflowAction  `json:"actions"`
}

// WorkflowInstanceRequest request to trigger workflow
type WorkflowInstanceRequest struct {
	WorkflowID       int64                  `json:"workflow_id"`
	TriggeredBy      string                 `json:"triggered_by"`
	TriggeredByValue string                 `json:"triggered_by_value"`
	AdditionalData   map[string]interface{} `json:"additional_data"`
}
