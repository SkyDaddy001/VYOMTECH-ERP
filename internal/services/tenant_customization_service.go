package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TenantCustomizationService defines tenant-level customization operations
type TenantCustomizationService interface {
	// Task Status Management
	CreateTaskStatus(ctx context.Context, tenantID string, status *TenantTaskStatus) (*TenantTaskStatus, error)
	GetTaskStatus(ctx context.Context, tenantID, statusCode string) (*TenantTaskStatus, error)
	GetTaskStatuses(ctx context.Context, tenantID string) ([]TenantTaskStatus, error)
	UpdateTaskStatus(ctx context.Context, tenantID string, status *TenantTaskStatus) (*TenantTaskStatus, error)
	DeactivateTaskStatus(ctx context.Context, tenantID, statusCode string) error

	// Task Stage Management
	CreateTaskStage(ctx context.Context, tenantID string, stage *TenantTaskStage) (*TenantTaskStage, error)
	GetTaskStages(ctx context.Context, tenantID string) ([]TenantTaskStage, error)
	UpdateTaskStage(ctx context.Context, tenantID string, stage *TenantTaskStage) (*TenantTaskStage, error)

	// Status Transitions
	CreateStatusTransition(ctx context.Context, tenantID string, transition *TenantStatusTransition) (*TenantStatusTransition, error)
	GetAllowedTransitions(ctx context.Context, tenantID, fromStatus string) ([]string, error)
	IsTransitionAllowed(ctx context.Context, tenantID, fromStatus, toStatus string) (bool, error)

	// Task Types
	CreateTaskType(ctx context.Context, tenantID string, taskType *TenantTaskType) (*TenantTaskType, error)
	GetTaskTypes(ctx context.Context, tenantID string) ([]TenantTaskType, error)
	UpdateTaskType(ctx context.Context, tenantID string, taskType *TenantTaskType) (*TenantTaskType, error)

	// Priority Levels
	CreatePriorityLevel(ctx context.Context, tenantID string, priority *TenantPriorityLevel) (*TenantPriorityLevel, error)
	GetPriorityLevels(ctx context.Context, tenantID string) ([]TenantPriorityLevel, error)
	UpdatePriorityLevel(ctx context.Context, tenantID string, priority *TenantPriorityLevel) (*TenantPriorityLevel, error)

	// Notification Types
	CreateNotificationType(ctx context.Context, tenantID string, notifType *TenantNotificationType) (*TenantNotificationType, error)
	GetNotificationTypes(ctx context.Context, tenantID string) ([]TenantNotificationType, error)
	UpdateNotificationType(ctx context.Context, tenantID string, notifType *TenantNotificationType) (*TenantNotificationType, error)

	// Custom Fields
	CreateCustomField(ctx context.Context, tenantID string, field *TenantTaskField) (*TenantTaskField, error)
	GetCustomFields(ctx context.Context, tenantID string) ([]TenantTaskField, error)
	UpdateCustomField(ctx context.Context, tenantID string, field *TenantTaskField) (*TenantTaskField, error)

	// Automation Rules
	CreateAutomationRule(ctx context.Context, tenantID string, rule *TenantAutomationRule) (*TenantAutomationRule, error)
	GetAutomationRules(ctx context.Context, tenantID string) ([]TenantAutomationRule, error)
	UpdateAutomationRule(ctx context.Context, tenantID string, rule *TenantAutomationRule) (*TenantAutomationRule, error)

	// Get complete configuration for a tenant
	GetTenantConfiguration(ctx context.Context, tenantID string) (*TenantConfiguration, error)
}

// ============================================================================
// DATA MODELS
// ============================================================================

// TenantTaskStatus represents a custom task status
type TenantTaskStatus struct {
	ID                 int64     `json:"id"`
	TenantID           string    `json:"tenant_id"`
	StatusCode         string    `json:"status_code"`
	StatusName         string    `json:"status_name"`
	Description        *string   `json:"description,omitempty"`
	ColorHex           *string   `json:"color_hex,omitempty"`
	Icon               *string   `json:"icon,omitempty"`
	DisplayOrder       int       `json:"display_order"`
	IsActive           bool      `json:"is_active"`
	IsInitialStatus    bool      `json:"is_initial_status"`
	IsFinalStatus      bool      `json:"is_final_status"`
	IsBlockingStatus   bool      `json:"is_blocking_status"`
	AllowsEditing      bool      `json:"allows_editing"`
	AllowsReassignment bool      `json:"allows_reassignment"`
	CreatedBy          *int64    `json:"created_by,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// TenantTaskStage represents a custom task stage
type TenantTaskStage struct {
	ID                   int64     `json:"id"`
	TenantID             string    `json:"tenant_id"`
	StageCode            string    `json:"stage_code"`
	StageName            string    `json:"stage_name"`
	Description          *string   `json:"description,omitempty"`
	ColorHex             *string   `json:"color_hex,omitempty"`
	Icon                 *string   `json:"icon,omitempty"`
	DisplayOrder         int       `json:"display_order"`
	IsActive             bool      `json:"is_active"`
	MinDurationHours     *int      `json:"min_duration_hours,omitempty"`
	MaxDurationHours     *int      `json:"max_duration_hours,omitempty"`
	SLAMinutes           *int      `json:"sla_minutes,omitempty"`
	AutoAdvanceToStageID *int64    `json:"auto_advance_to_stage_id,omitempty"`
	CreatedBy            *int64    `json:"created_by,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// TenantStatusTransition defines allowed status transitions
type TenantStatusTransition struct {
	ID                       int64     `json:"id"`
	TenantID                 string    `json:"tenant_id"`
	FromStatusCode           string    `json:"from_status_code"`
	ToStatusCode             string    `json:"to_status_code"`
	IsAllowed                bool      `json:"is_allowed"`
	RequiresComment          bool      `json:"requires_comment"`
	RequiresApproval         bool      `json:"requires_approval"`
	NotificationOnTransition bool      `json:"notification_on_transition"`
	RequiresRole             *string   `json:"requires_role,omitempty"`
	RequiresFieldCompletion  *string   `json:"requires_field_completion,omitempty"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

// TenantTaskType represents a custom task type
type TenantTaskType struct {
	ID                int64     `json:"id"`
	TenantID          string    `json:"tenant_id"`
	TypeCode          string    `json:"type_code"`
	TypeName          string    `json:"type_name"`
	Description       *string   `json:"description,omitempty"`
	Icon              *string   `json:"icon,omitempty"`
	ColorHex          *string   `json:"color_hex,omitempty"`
	DefaultPriority   string    `json:"default_priority"`
	DefaultDueDays    *int      `json:"default_due_days,omitempty"`
	RequiredStatuses  *string   `json:"required_statuses,omitempty"`
	IsLeadRelated     bool      `json:"is_lead_related"`
	IsAgentAssignable bool      `json:"is_agent_assignable"`
	IsActive          bool      `json:"is_active"`
	CreatedBy         *int64    `json:"created_by,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TenantPriorityLevel represents a custom priority level
type TenantPriorityLevel struct {
	ID                 int64     `json:"id"`
	TenantID           string    `json:"tenant_id"`
	PriorityCode       string    `json:"priority_code"`
	PriorityName       string    `json:"priority_name"`
	PriorityValue      int       `json:"priority_value"`
	ColorHex           *string   `json:"color_hex,omitempty"`
	Icon               *string   `json:"icon,omitempty"`
	Description        *string   `json:"description,omitempty"`
	SLAResponseHours   *int      `json:"sla_response_hours,omitempty"`
	SLAResolutionHours *int      `json:"sla_resolution_hours,omitempty"`
	NotifyOnAssignment bool      `json:"notify_on_assignment"`
	NotifySupervisors  bool      `json:"notify_supervisors"`
	EscalationEnabled  bool      `json:"escalation_enabled"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// TenantNotificationType represents a custom notification type
type TenantNotificationType struct {
	ID                int64     `json:"id"`
	TenantID          string    `json:"tenant_id"`
	TypeCode          string    `json:"type_code"`
	TypeName          string    `json:"type_name"`
	Description       *string   `json:"description,omitempty"`
	Icon              *string   `json:"icon,omitempty"`
	ColorHex          *string   `json:"color_hex,omitempty"`
	DefaultPriority   string    `json:"default_priority"`
	Category          string    `json:"category"`
	SupportedChannels *string   `json:"supported_channels,omitempty"`
	DefaultChannels   *string   `json:"default_channels,omitempty"`
	IsDismissable     bool      `json:"is_dismissable"`
	AutoArchiveDays   *int      `json:"auto_archive_days,omitempty"`
	IsActive          bool      `json:"is_active"`
	CreatedBy         *int64    `json:"created_by,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TenantTaskField represents a custom task field
type TenantTaskField struct {
	ID                 int64     `json:"id"`
	TenantID           string    `json:"tenant_id"`
	FieldCode          string    `json:"field_code"`
	FieldName          string    `json:"field_name"`
	FieldType          string    `json:"field_type"`
	IsRequired         bool      `json:"is_required"`
	IsVisible          bool      `json:"is_visible"`
	IsEditable         bool      `json:"is_editable"`
	DisplayOrder       int       `json:"display_order"`
	ValidationRules    *string   `json:"validation_rules,omitempty"`
	DefaultValue       *string   `json:"default_value,omitempty"`
	FieldOptions       *string   `json:"field_options,omitempty"`
	VisibleOnStatuses  *string   `json:"visible_on_statuses,omitempty"`
	VisibleOnTaskTypes *string   `json:"visible_on_task_types,omitempty"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// TenantAutomationRule represents an automation rule
type TenantAutomationRule struct {
	ID                int64     `json:"id"`
	TenantID          string    `json:"tenant_id"`
	RuleCode          string    `json:"rule_code"`
	RuleName          string    `json:"rule_name"`
	Description       *string   `json:"description,omitempty"`
	TriggerEvent      string    `json:"trigger_event"`
	TriggerConditions *string   `json:"trigger_conditions,omitempty"`
	ActionType        string    `json:"action_type"`
	ActionData        *string   `json:"action_data,omitempty"`
	IsActive          bool      `json:"is_active"`
	Priority          int       `json:"priority"`
	RunOnce           bool      `json:"run_once"`
	CreatedBy         *int64    `json:"created_by,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TenantConfiguration aggregates all tenant customization
type TenantConfiguration struct {
	TenantID          string                   `json:"tenant_id"`
	Statuses          []TenantTaskStatus       `json:"statuses"`
	Stages            []TenantTaskStage        `json:"stages"`
	TaskTypes         []TenantTaskType         `json:"task_types"`
	PriorityLevels    []TenantPriorityLevel    `json:"priority_levels"`
	NotificationTypes []TenantNotificationType `json:"notification_types"`
	CustomFields      []TenantTaskField        `json:"custom_fields"`
	AutomationRules   []TenantAutomationRule   `json:"automation_rules"`
	StatusTransitions []TenantStatusTransition `json:"status_transitions"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

// ============================================================================
// SERVICE IMPLEMENTATION
// ============================================================================

type tenantCustomizationService struct {
	db *sql.DB
}

// NewTenantCustomizationService creates a new customization service
func NewTenantCustomizationService(db *sql.DB) TenantCustomizationService {
	return &tenantCustomizationService{db: db}
}

// ============================================================================
// TASK STATUS METHODS
// ============================================================================

// CreateTaskStatus creates a new custom status
func (s *tenantCustomizationService) CreateTaskStatus(ctx context.Context, tenantID string, status *TenantTaskStatus) (*TenantTaskStatus, error) {
	query := `
		INSERT INTO tenant_task_statuses (tenant_id, status_code, status_name, description, color_hex, icon,
		                                 display_order, is_active, is_initial_status, is_final_status,
		                                 is_blocking_status, allows_editing, allows_reassignment)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		tenantID, status.StatusCode, status.StatusName, status.Description, status.ColorHex, status.Icon,
		status.DisplayOrder, status.IsActive, status.IsInitialStatus, status.IsFinalStatus,
		status.IsBlockingStatus, status.AllowsEditing, status.AllowsReassignment)

	if err != nil {
		return nil, fmt.Errorf("failed to create task status: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get status ID: %w", err)
	}

	status.ID = id
	status.TenantID = tenantID
	status.CreatedAt = time.Now()
	status.UpdatedAt = time.Now()

	return status, nil
}

// GetTaskStatus retrieves a specific status
func (s *tenantCustomizationService) GetTaskStatus(ctx context.Context, tenantID, statusCode string) (*TenantTaskStatus, error) {
	query := `
		SELECT id, tenant_id, status_code, status_name, description, color_hex, icon, display_order,
		       is_active, is_initial_status, is_final_status, is_blocking_status, allows_editing,
		       allows_reassignment, created_by, created_at, updated_at
		FROM tenant_task_statuses
		WHERE tenant_id = ? AND status_code = ?
	`

	status := &TenantTaskStatus{}
	err := s.db.QueryRowContext(ctx, query, tenantID, statusCode).Scan(
		&status.ID, &status.TenantID, &status.StatusCode, &status.StatusName, &status.Description,
		&status.ColorHex, &status.Icon, &status.DisplayOrder, &status.IsActive, &status.IsInitialStatus,
		&status.IsFinalStatus, &status.IsBlockingStatus, &status.AllowsEditing, &status.AllowsReassignment,
		&status.CreatedBy, &status.CreatedAt, &status.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("status not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	return status, nil
}

// GetTaskStatuses retrieves all statuses for a tenant
func (s *tenantCustomizationService) GetTaskStatuses(ctx context.Context, tenantID string) ([]TenantTaskStatus, error) {
	query := `
		SELECT id, tenant_id, status_code, status_name, description, color_hex, icon, display_order,
		       is_active, is_initial_status, is_final_status, is_blocking_status, allows_editing,
		       allows_reassignment, created_by, created_at, updated_at
		FROM tenant_task_statuses
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY display_order ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statuses: %w", err)
	}
	defer rows.Close()

	var statuses []TenantTaskStatus
	for rows.Next() {
		status := TenantTaskStatus{}
		err := rows.Scan(
			&status.ID, &status.TenantID, &status.StatusCode, &status.StatusName, &status.Description,
			&status.ColorHex, &status.Icon, &status.DisplayOrder, &status.IsActive, &status.IsInitialStatus,
			&status.IsFinalStatus, &status.IsBlockingStatus, &status.AllowsEditing, &status.AllowsReassignment,
			&status.CreatedBy, &status.CreatedAt, &status.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan status: %w", err)
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

// UpdateTaskStatus updates a status
func (s *tenantCustomizationService) UpdateTaskStatus(ctx context.Context, tenantID string, status *TenantTaskStatus) (*TenantTaskStatus, error) {
	query := `
		UPDATE tenant_task_statuses
		SET status_name = ?, description = ?, color_hex = ?, icon = ?, display_order = ?,
		    is_active = ?, is_initial_status = ?, is_final_status = ?, is_blocking_status = ?,
		    allows_editing = ?, allows_reassignment = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		status.StatusName, status.Description, status.ColorHex, status.Icon, status.DisplayOrder,
		status.IsActive, status.IsInitialStatus, status.IsFinalStatus, status.IsBlockingStatus,
		status.AllowsEditing, status.AllowsReassignment, status.ID, tenantID)

	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	status.UpdatedAt = time.Now()
	return status, nil
}

// DeactivateTaskStatus deactivates a status
func (s *tenantCustomizationService) DeactivateTaskStatus(ctx context.Context, tenantID, statusCode string) error {
	query := `
		UPDATE tenant_task_statuses
		SET is_active = FALSE, updated_at = NOW()
		WHERE tenant_id = ? AND status_code = ?
	`

	_, err := s.db.ExecContext(ctx, query, tenantID, statusCode)
	if err != nil {
		return fmt.Errorf("failed to deactivate status: %w", err)
	}

	return nil
}

// ============================================================================
// TASK STAGE METHODS
// ============================================================================

// CreateTaskStage creates a new task stage
func (s *tenantCustomizationService) CreateTaskStage(ctx context.Context, tenantID string, stage *TenantTaskStage) (*TenantTaskStage, error) {
	query := `
		INSERT INTO tenant_task_stages (tenant_id, stage_code, stage_name, description, color_hex, icon,
		                               display_order, is_active, min_duration_hours, max_duration_hours,
		                               sla_minutes, auto_advance_to_stage_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		tenantID, stage.StageCode, stage.StageName, stage.Description, stage.ColorHex, stage.Icon,
		stage.DisplayOrder, stage.IsActive, stage.MinDurationHours, stage.MaxDurationHours,
		stage.SLAMinutes, stage.AutoAdvanceToStageID)

	if err != nil {
		return nil, fmt.Errorf("failed to create stage: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get stage ID: %w", err)
	}

	stage.ID = id
	stage.TenantID = tenantID
	stage.CreatedAt = time.Now()
	stage.UpdatedAt = time.Now()

	return stage, nil
}

// GetTaskStages retrieves all stages for a tenant
func (s *tenantCustomizationService) GetTaskStages(ctx context.Context, tenantID string) ([]TenantTaskStage, error) {
	query := `
		SELECT id, tenant_id, stage_code, stage_name, description, color_hex, icon, display_order,
		       is_active, min_duration_hours, max_duration_hours, sla_minutes, auto_advance_to_stage_id,
		       created_by, created_at, updated_at
		FROM tenant_task_stages
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY display_order ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stages: %w", err)
	}
	defer rows.Close()

	var stages []TenantTaskStage
	for rows.Next() {
		stage := TenantTaskStage{}
		err := rows.Scan(
			&stage.ID, &stage.TenantID, &stage.StageCode, &stage.StageName, &stage.Description,
			&stage.ColorHex, &stage.Icon, &stage.DisplayOrder, &stage.IsActive, &stage.MinDurationHours,
			&stage.MaxDurationHours, &stage.SLAMinutes, &stage.AutoAdvanceToStageID,
			&stage.CreatedBy, &stage.CreatedAt, &stage.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stage: %w", err)
		}
		stages = append(stages, stage)
	}

	return stages, nil
}

// UpdateTaskStage updates a stage
func (s *tenantCustomizationService) UpdateTaskStage(ctx context.Context, tenantID string, stage *TenantTaskStage) (*TenantTaskStage, error) {
	query := `
		UPDATE tenant_task_stages
		SET stage_name = ?, description = ?, color_hex = ?, icon = ?, display_order = ?,
		    is_active = ?, min_duration_hours = ?, max_duration_hours = ?, sla_minutes = ?,
		    auto_advance_to_stage_id = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		stage.StageName, stage.Description, stage.ColorHex, stage.Icon, stage.DisplayOrder,
		stage.IsActive, stage.MinDurationHours, stage.MaxDurationHours, stage.SLAMinutes,
		stage.AutoAdvanceToStageID, stage.ID, tenantID)

	if err != nil {
		return nil, fmt.Errorf("failed to update stage: %w", err)
	}

	stage.UpdatedAt = time.Now()
	return stage, nil
}

// ============================================================================
// STATUS TRANSITION METHODS
// ============================================================================

// CreateStatusTransition creates a status transition rule
func (s *tenantCustomizationService) CreateStatusTransition(ctx context.Context, tenantID string, transition *TenantStatusTransition) (*TenantStatusTransition, error) {
	query := `
		INSERT INTO tenant_status_transitions (tenant_id, from_status_code, to_status_code, is_allowed,
		                                       requires_comment, requires_approval, notification_on_transition,
		                                       requires_role, requires_field_completion)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		tenantID, transition.FromStatusCode, transition.ToStatusCode, transition.IsAllowed,
		transition.RequiresComment, transition.RequiresApproval, transition.NotificationOnTransition,
		transition.RequiresRole, transition.RequiresFieldCompletion)

	if err != nil {
		return nil, fmt.Errorf("failed to create transition: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get transition ID: %w", err)
	}

	transition.ID = id
	transition.TenantID = tenantID
	transition.CreatedAt = time.Now()
	transition.UpdatedAt = time.Now()

	return transition, nil
}

// GetAllowedTransitions gets allowed transitions from a status
func (s *tenantCustomizationService) GetAllowedTransitions(ctx context.Context, tenantID, fromStatus string) ([]string, error) {
	query := `
		SELECT to_status_code
		FROM tenant_status_transitions
		WHERE tenant_id = ? AND from_status_code = ? AND is_allowed = TRUE
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, fromStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to get transitions: %w", err)
	}
	defer rows.Close()

	var transitions []string
	for rows.Next() {
		var toStatus string
		err := rows.Scan(&toStatus)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transition: %w", err)
		}
		transitions = append(transitions, toStatus)
	}

	return transitions, nil
}

// IsTransitionAllowed checks if a status transition is allowed
func (s *tenantCustomizationService) IsTransitionAllowed(ctx context.Context, tenantID, fromStatus, toStatus string) (bool, error) {
	query := `
		SELECT is_allowed
		FROM tenant_status_transitions
		WHERE tenant_id = ? AND from_status_code = ? AND to_status_code = ?
	`

	var isAllowed bool
	err := s.db.QueryRowContext(ctx, query, tenantID, fromStatus, toStatus).Scan(&isAllowed)
	if err == sql.ErrNoRows {
		// If no explicit rule, allow transition by default
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check transition: %w", err)
	}

	return isAllowed, nil
}

// ============================================================================
// TASK TYPE METHODS
// ============================================================================

// CreateTaskType creates a new task type
func (s *tenantCustomizationService) CreateTaskType(ctx context.Context, tenantID string, taskType *TenantTaskType) (*TenantTaskType, error) {
	query := `
		INSERT INTO tenant_task_types (tenant_id, type_code, type_name, description, icon, color_hex,
		                              default_priority, default_due_days, required_statuses,
		                              is_lead_related, is_agent_assignable, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		tenantID, taskType.TypeCode, taskType.TypeName, taskType.Description, taskType.Icon, taskType.ColorHex,
		taskType.DefaultPriority, taskType.DefaultDueDays, taskType.RequiredStatuses,
		taskType.IsLeadRelated, taskType.IsAgentAssignable, taskType.IsActive)

	if err != nil {
		return nil, fmt.Errorf("failed to create task type: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get type ID: %w", err)
	}

	taskType.ID = id
	taskType.TenantID = tenantID
	taskType.CreatedAt = time.Now()
	taskType.UpdatedAt = time.Now()

	return taskType, nil
}

// GetTaskTypes retrieves all task types for a tenant
func (s *tenantCustomizationService) GetTaskTypes(ctx context.Context, tenantID string) ([]TenantTaskType, error) {
	query := `
		SELECT id, tenant_id, type_code, type_name, description, icon, color_hex, default_priority,
		       default_due_days, required_statuses, is_lead_related, is_agent_assignable, is_active,
		       created_by, created_at, updated_at
		FROM tenant_task_types
		WHERE tenant_id = ? AND is_active = TRUE
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task types: %w", err)
	}
	defer rows.Close()

	var types []TenantTaskType
	for rows.Next() {
		tt := TenantTaskType{}
		err := rows.Scan(
			&tt.ID, &tt.TenantID, &tt.TypeCode, &tt.TypeName, &tt.Description, &tt.Icon, &tt.ColorHex,
			&tt.DefaultPriority, &tt.DefaultDueDays, &tt.RequiredStatuses, &tt.IsLeadRelated,
			&tt.IsAgentAssignable, &tt.IsActive, &tt.CreatedBy, &tt.CreatedAt, &tt.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan type: %w", err)
		}
		types = append(types, tt)
	}

	return types, nil
}

// UpdateTaskType updates a task type
func (s *tenantCustomizationService) UpdateTaskType(ctx context.Context, tenantID string, taskType *TenantTaskType) (*TenantTaskType, error) {
	query := `
		UPDATE tenant_task_types
		SET type_name = ?, description = ?, icon = ?, color_hex = ?, default_priority = ?,
		    default_due_days = ?, required_statuses = ?, is_lead_related = ?, is_agent_assignable = ?,
		    is_active = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		taskType.TypeName, taskType.Description, taskType.Icon, taskType.ColorHex, taskType.DefaultPriority,
		taskType.DefaultDueDays, taskType.RequiredStatuses, taskType.IsLeadRelated, taskType.IsAgentAssignable,
		taskType.IsActive, taskType.ID, tenantID)

	if err != nil {
		return nil, fmt.Errorf("failed to update task type: %w", err)
	}

	taskType.UpdatedAt = time.Now()
	return taskType, nil
}

// ============================================================================
// PRIORITY LEVEL METHODS
// ============================================================================

// CreatePriorityLevel creates a new priority level
func (s *tenantCustomizationService) CreatePriorityLevel(ctx context.Context, tenantID string, priority *TenantPriorityLevel) (*TenantPriorityLevel, error) {
	query := `
		INSERT INTO tenant_priority_levels (tenant_id, priority_code, priority_name, priority_value,
		                                   color_hex, icon, description, sla_response_hours,
		                                   sla_resolution_hours, notify_on_assignment, notify_supervisors,
		                                   escalation_enabled, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		tenantID, priority.PriorityCode, priority.PriorityName, priority.PriorityValue,
		priority.ColorHex, priority.Icon, priority.Description, priority.SLAResponseHours,
		priority.SLAResolutionHours, priority.NotifyOnAssignment, priority.NotifySupervisors,
		priority.EscalationEnabled, priority.IsActive)

	if err != nil {
		return nil, fmt.Errorf("failed to create priority: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get priority ID: %w", err)
	}

	priority.ID = id
	priority.TenantID = tenantID
	priority.CreatedAt = time.Now()
	priority.UpdatedAt = time.Now()

	return priority, nil
}

// GetPriorityLevels retrieves all priority levels for a tenant
func (s *tenantCustomizationService) GetPriorityLevels(ctx context.Context, tenantID string) ([]TenantPriorityLevel, error) {
	query := `
		SELECT id, tenant_id, priority_code, priority_name, priority_value, color_hex, icon, description,
		       sla_response_hours, sla_resolution_hours, notify_on_assignment, notify_supervisors,
		       escalation_enabled, is_active, created_at, updated_at
		FROM tenant_priority_levels
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY priority_value DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get priorities: %w", err)
	}
	defer rows.Close()

	var priorities []TenantPriorityLevel
	for rows.Next() {
		priority := TenantPriorityLevel{}
		err := rows.Scan(
			&priority.ID, &priority.TenantID, &priority.PriorityCode, &priority.PriorityName, &priority.PriorityValue,
			&priority.ColorHex, &priority.Icon, &priority.Description, &priority.SLAResponseHours,
			&priority.SLAResolutionHours, &priority.NotifyOnAssignment, &priority.NotifySupervisors,
			&priority.EscalationEnabled, &priority.IsActive, &priority.CreatedAt, &priority.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan priority: %w", err)
		}
		priorities = append(priorities, priority)
	}

	return priorities, nil
}

// UpdatePriorityLevel updates a priority level
func (s *tenantCustomizationService) UpdatePriorityLevel(ctx context.Context, tenantID string, priority *TenantPriorityLevel) (*TenantPriorityLevel, error) {
	query := `
		UPDATE tenant_priority_levels
		SET priority_name = ?, priority_value = ?, color_hex = ?, icon = ?, description = ?,
		    sla_response_hours = ?, sla_resolution_hours = ?, notify_on_assignment = ?,
		    notify_supervisors = ?, escalation_enabled = ?, is_active = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		priority.PriorityName, priority.PriorityValue, priority.ColorHex, priority.Icon, priority.Description,
		priority.SLAResponseHours, priority.SLAResolutionHours, priority.NotifyOnAssignment,
		priority.NotifySupervisors, priority.EscalationEnabled, priority.IsActive, priority.ID, tenantID)

	if err != nil {
		return nil, fmt.Errorf("failed to update priority: %w", err)
	}

	priority.UpdatedAt = time.Now()
	return priority, nil
}

// ============================================================================
// NOTIFICATION TYPE METHODS
// ============================================================================

// CreateNotificationType creates a new notification type
func (s *tenantCustomizationService) CreateNotificationType(ctx context.Context, tenantID string, notifType *TenantNotificationType) (*TenantNotificationType, error) {
	query := `
		INSERT INTO tenant_notification_types (tenant_id, type_code, type_name, description, icon, color_hex,
		                                       default_priority, category, supported_channels, default_channels,
		                                       is_dismissable, auto_archive_after_days, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		tenantID, notifType.TypeCode, notifType.TypeName, notifType.Description, notifType.Icon, notifType.ColorHex,
		notifType.DefaultPriority, notifType.Category, notifType.SupportedChannels, notifType.DefaultChannels,
		notifType.IsDismissable, notifType.AutoArchiveDays, notifType.IsActive)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification type: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get type ID: %w", err)
	}

	notifType.ID = id
	notifType.TenantID = tenantID
	notifType.CreatedAt = time.Now()
	notifType.UpdatedAt = time.Now()

	return notifType, nil
}

// GetNotificationTypes retrieves all notification types for a tenant
func (s *tenantCustomizationService) GetNotificationTypes(ctx context.Context, tenantID string) ([]TenantNotificationType, error) {
	query := `
		SELECT id, tenant_id, type_code, type_name, description, icon, color_hex, default_priority,
		       category, supported_channels, default_channels, is_dismissable, auto_archive_after_days,
		       is_active, created_by, created_at, updated_at
		FROM tenant_notification_types
		WHERE tenant_id = ? AND is_active = TRUE
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification types: %w", err)
	}
	defer rows.Close()

	var types []TenantNotificationType
	for rows.Next() {
		nt := TenantNotificationType{}
		err := rows.Scan(
			&nt.ID, &nt.TenantID, &nt.TypeCode, &nt.TypeName, &nt.Description, &nt.Icon, &nt.ColorHex,
			&nt.DefaultPriority, &nt.Category, &nt.SupportedChannels, &nt.DefaultChannels,
			&nt.IsDismissable, &nt.AutoArchiveDays, &nt.IsActive, &nt.CreatedBy, &nt.CreatedAt, &nt.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan type: %w", err)
		}
		types = append(types, nt)
	}

	return types, nil
}

// UpdateNotificationType updates a notification type
func (s *tenantCustomizationService) UpdateNotificationType(ctx context.Context, tenantID string, notifType *TenantNotificationType) (*TenantNotificationType, error) {
	query := `
		UPDATE tenant_notification_types
		SET type_name = ?, description = ?, icon = ?, color_hex = ?, default_priority = ?,
		    category = ?, supported_channels = ?, default_channels = ?, is_dismissable = ?,
		    auto_archive_after_days = ?, is_active = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		notifType.TypeName, notifType.Description, notifType.Icon, notifType.ColorHex, notifType.DefaultPriority,
		notifType.Category, notifType.SupportedChannels, notifType.DefaultChannels, notifType.IsDismissable,
		notifType.AutoArchiveDays, notifType.IsActive, notifType.ID, tenantID)

	if err != nil {
		return nil, fmt.Errorf("failed to update notification type: %w", err)
	}

	notifType.UpdatedAt = time.Now()
	return notifType, nil
}

// ============================================================================
// CUSTOM FIELD METHODS
// ============================================================================

// CreateCustomField creates a new custom field
func (s *tenantCustomizationService) CreateCustomField(ctx context.Context, tenantID string, field *TenantTaskField) (*TenantTaskField, error) {
	query := `
		INSERT INTO tenant_task_fields (tenant_id, field_code, field_name, field_type, is_required, is_visible,
		                               is_editable, display_order, validation_rules, default_value, field_options,
		                               visible_on_statuses, visible_on_task_types, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		tenantID, field.FieldCode, field.FieldName, field.FieldType, field.IsRequired, field.IsVisible,
		field.IsEditable, field.DisplayOrder, field.ValidationRules, field.DefaultValue, field.FieldOptions,
		field.VisibleOnStatuses, field.VisibleOnTaskTypes, field.IsActive)

	if err != nil {
		return nil, fmt.Errorf("failed to create field: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get field ID: %w", err)
	}

	field.ID = id
	field.TenantID = tenantID
	field.CreatedAt = time.Now()
	field.UpdatedAt = time.Now()

	return field, nil
}

// GetCustomFields retrieves all custom fields for a tenant
func (s *tenantCustomizationService) GetCustomFields(ctx context.Context, tenantID string) ([]TenantTaskField, error) {
	query := `
		SELECT id, tenant_id, field_code, field_name, field_type, is_required, is_visible, is_editable,
		       display_order, validation_rules, default_value, field_options, visible_on_statuses,
		       visible_on_task_types, is_active, created_at, updated_at
		FROM tenant_task_fields
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY display_order ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fields: %w", err)
	}
	defer rows.Close()

	var fields []TenantTaskField
	for rows.Next() {
		field := TenantTaskField{}
		err := rows.Scan(
			&field.ID, &field.TenantID, &field.FieldCode, &field.FieldName, &field.FieldType, &field.IsRequired,
			&field.IsVisible, &field.IsEditable, &field.DisplayOrder, &field.ValidationRules, &field.DefaultValue,
			&field.FieldOptions, &field.VisibleOnStatuses, &field.VisibleOnTaskTypes, &field.IsActive,
			&field.CreatedAt, &field.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan field: %w", err)
		}
		fields = append(fields, field)
	}

	return fields, nil
}

// UpdateCustomField updates a custom field
func (s *tenantCustomizationService) UpdateCustomField(ctx context.Context, tenantID string, field *TenantTaskField) (*TenantTaskField, error) {
	query := `
		UPDATE tenant_task_fields
		SET field_name = ?, field_type = ?, is_required = ?, is_visible = ?, is_editable = ?,
		    display_order = ?, validation_rules = ?, default_value = ?, field_options = ?,
		    visible_on_statuses = ?, visible_on_task_types = ?, is_active = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		field.FieldName, field.FieldType, field.IsRequired, field.IsVisible, field.IsEditable,
		field.DisplayOrder, field.ValidationRules, field.DefaultValue, field.FieldOptions,
		field.VisibleOnStatuses, field.VisibleOnTaskTypes, field.IsActive, field.ID, tenantID)

	if err != nil {
		return nil, fmt.Errorf("failed to update field: %w", err)
	}

	field.UpdatedAt = time.Now()
	return field, nil
}

// ============================================================================
// AUTOMATION RULE METHODS
// ============================================================================

// CreateAutomationRule creates a new automation rule
func (s *tenantCustomizationService) CreateAutomationRule(ctx context.Context, tenantID string, rule *TenantAutomationRule) (*TenantAutomationRule, error) {
	query := `
		INSERT INTO tenant_automation_rules (tenant_id, rule_code, rule_name, description, trigger_event,
		                                    trigger_conditions, action_type, action_data, is_active, priority, run_once)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		tenantID, rule.RuleCode, rule.RuleName, rule.Description, rule.TriggerEvent,
		rule.TriggerConditions, rule.ActionType, rule.ActionData, rule.IsActive, rule.Priority, rule.RunOnce)

	if err != nil {
		return nil, fmt.Errorf("failed to create rule: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get rule ID: %w", err)
	}

	rule.ID = id
	rule.TenantID = tenantID
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	return rule, nil
}

// GetAutomationRules retrieves all automation rules for a tenant
func (s *tenantCustomizationService) GetAutomationRules(ctx context.Context, tenantID string) ([]TenantAutomationRule, error) {
	query := `
		SELECT id, tenant_id, rule_code, rule_name, description, trigger_event, trigger_conditions,
		       action_type, action_data, is_active, priority, run_once, created_by, created_at, updated_at
		FROM tenant_automation_rules
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY priority ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rules: %w", err)
	}
	defer rows.Close()

	var rules []TenantAutomationRule
	for rows.Next() {
		rule := TenantAutomationRule{}
		err := rows.Scan(
			&rule.ID, &rule.TenantID, &rule.RuleCode, &rule.RuleName, &rule.Description, &rule.TriggerEvent,
			&rule.TriggerConditions, &rule.ActionType, &rule.ActionData, &rule.IsActive, &rule.Priority,
			&rule.RunOnce, &rule.CreatedBy, &rule.CreatedAt, &rule.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rule: %w", err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// UpdateAutomationRule updates an automation rule
func (s *tenantCustomizationService) UpdateAutomationRule(ctx context.Context, tenantID string, rule *TenantAutomationRule) (*TenantAutomationRule, error) {
	query := `
		UPDATE tenant_automation_rules
		SET rule_name = ?, description = ?, trigger_event = ?, trigger_conditions = ?,
		    action_type = ?, action_data = ?, is_active = ?, priority = ?, run_once = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		rule.RuleName, rule.Description, rule.TriggerEvent, rule.TriggerConditions,
		rule.ActionType, rule.ActionData, rule.IsActive, rule.Priority, rule.RunOnce, rule.ID, tenantID)

	if err != nil {
		return nil, fmt.Errorf("failed to update rule: %w", err)
	}

	rule.UpdatedAt = time.Now()
	return rule, nil
}

// ============================================================================
// AGGREGATE CONFIGURATION METHOD
// ============================================================================

// GetTenantConfiguration retrieves complete tenant configuration
func (s *tenantCustomizationService) GetTenantConfiguration(ctx context.Context, tenantID string) (*TenantConfiguration, error) {
	config := &TenantConfiguration{
		TenantID: tenantID,
	}

	var err error

	// Get all components in parallel if possible
	config.Statuses, err = s.GetTaskStatuses(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statuses: %w", err)
	}

	config.Stages, err = s.GetTaskStages(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stages: %w", err)
	}

	config.TaskTypes, err = s.GetTaskTypes(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task types: %w", err)
	}

	config.PriorityLevels, err = s.GetPriorityLevels(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get priorities: %w", err)
	}

	config.NotificationTypes, err = s.GetNotificationTypes(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification types: %w", err)
	}

	config.CustomFields, err = s.GetCustomFields(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fields: %w", err)
	}

	config.AutomationRules, err = s.GetAutomationRules(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rules: %w", err)
	}

	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	return config, nil
}
