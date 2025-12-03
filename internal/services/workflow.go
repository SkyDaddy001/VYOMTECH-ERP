package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"vyomtech-backend/internal/models"
)

// ==================== WORKFLOW SERVICE ====================

type WorkflowService struct {
	db *sql.DB
}

// NewWorkflowService creates a new workflow service
func NewWorkflowService(db *sql.DB) *WorkflowService {
	return &WorkflowService{
		db: db,
	}
}

// ==================== WORKFLOW DEFINITION METHODS ====================

// CreateWorkflow creates a new workflow definition
func (s *WorkflowService) CreateWorkflow(tenantID string, req *models.WorkflowRequest, userID int64) (*models.WorkflowDefinition, error) {
	workflow := &models.WorkflowDefinition{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `
		INSERT INTO workflows (tenant_id, name, description, enabled, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, workflow.TenantID, workflow.Name, workflow.Description, workflow.Enabled, workflow.CreatedBy, workflow.CreatedAt, workflow.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow id: %w", err)
	}
	workflow.ID = id

	// Add triggers
	for _, trigger := range req.Triggers {
		if _, err := s.CreateWorkflowTrigger(tenantID, id, &trigger); err != nil {
			return nil, fmt.Errorf("failed to create trigger: %w", err)
		}
	}

	// Add actions
	for _, action := range req.Actions {
		if _, err := s.CreateWorkflowAction(tenantID, id, &action); err != nil {
			return nil, fmt.Errorf("failed to create action: %w", err)
		}
	}

	return workflow, nil
}

// GetWorkflow retrieves a workflow by ID
func (s *WorkflowService) GetWorkflow(tenantID string, workflowID int64) (*models.WorkflowDefinition, error) {
	query := `
		SELECT id, tenant_id, name, description, enabled, created_by, created_at, updated_at
		FROM workflows
		WHERE id = ? AND tenant_id = ?
	`
	workflow := &models.WorkflowDefinition{}
	err := s.db.QueryRow(query, workflowID, tenantID).Scan(
		&workflow.ID, &workflow.TenantID, &workflow.Name, &workflow.Description,
		&workflow.Enabled, &workflow.CreatedBy, &workflow.CreatedAt, &workflow.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	// Load triggers and actions
	triggers, _ := s.GetWorkflowTriggers(tenantID, workflowID)
	actions, _ := s.GetWorkflowActions(tenantID, workflowID)
	workflow.Triggers = triggers
	workflow.Actions = actions

	return workflow, nil
}

// ListWorkflows lists all workflows for a tenant
func (s *WorkflowService) ListWorkflows(tenantID string, limit int, offset int) ([]models.WorkflowDefinition, error) {
	query := `
		SELECT id, tenant_id, name, description, enabled, created_by, created_at, updated_at
		FROM workflows
		WHERE tenant_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := s.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}
	defer rows.Close()

	var workflows []models.WorkflowDefinition
	for rows.Next() {
		var workflow models.WorkflowDefinition
		err := rows.Scan(
			&workflow.ID, &workflow.TenantID, &workflow.Name, &workflow.Description,
			&workflow.Enabled, &workflow.CreatedBy, &workflow.CreatedAt, &workflow.UpdatedAt,
		)
		if err != nil {
			continue
		}
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// UpdateWorkflow updates a workflow definition
func (s *WorkflowService) UpdateWorkflow(tenantID string, workflowID int64, req *models.WorkflowRequest) (*models.WorkflowDefinition, error) {
	query := `
		UPDATE workflows
		SET name = ?, description = ?, enabled = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`
	_, err := s.db.Exec(query, req.Name, req.Description, req.Enabled, time.Now(), workflowID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to update workflow: %w", err)
	}

	// Delete existing triggers and actions, then add new ones
	s.db.Exec("DELETE FROM workflow_triggers WHERE workflow_id = ?", workflowID)
	s.db.Exec("DELETE FROM workflow_actions WHERE workflow_id = ?", workflowID)

	for _, trigger := range req.Triggers {
		s.CreateWorkflowTrigger(tenantID, workflowID, &trigger)
	}

	for _, action := range req.Actions {
		s.CreateWorkflowAction(tenantID, workflowID, &action)
	}

	return s.GetWorkflow(tenantID, workflowID)
}

// DeleteWorkflow deletes a workflow
func (s *WorkflowService) DeleteWorkflow(tenantID string, workflowID int64) error {
	// Delete related records
	s.db.Exec("DELETE FROM workflow_triggers WHERE workflow_id = ?", workflowID)
	s.db.Exec("DELETE FROM workflow_actions WHERE workflow_id = ?", workflowID)
	s.db.Exec("DELETE FROM workflow_instances WHERE workflow_id = ?", workflowID)

	query := `DELETE FROM workflows WHERE id = ? AND tenant_id = ?`
	_, err := s.db.Exec(query, workflowID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete workflow: %w", err)
	}

	return nil
}

// ==================== WORKFLOW TRIGGER METHODS ====================

// CreateWorkflowTrigger creates a trigger for a workflow
func (s *WorkflowService) CreateWorkflowTrigger(tenantID string, workflowID int64, trigger *models.WorkflowTrigger) (*models.WorkflowTrigger, error) {
	configJSON, _ := json.Marshal(trigger.TriggerConfig)

	query := `
		INSERT INTO workflow_triggers (workflow_id, trigger_type, trigger_config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, workflowID, trigger.TriggerType, string(configJSON), time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow trigger: %w", err)
	}

	id, _ := result.LastInsertId()
	trigger.ID = id
	trigger.WorkflowID = workflowID

	return trigger, nil
}

// GetWorkflowTriggers retrieves all triggers for a workflow
func (s *WorkflowService) GetWorkflowTriggers(tenantID string, workflowID int64) ([]models.WorkflowTrigger, error) {
	query := `
		SELECT id, workflow_id, trigger_type, trigger_config, created_at, updated_at
		FROM workflow_triggers
		WHERE workflow_id = ?
	`
	rows, err := s.db.Query(query, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow triggers: %w", err)
	}
	defer rows.Close()

	var triggers []models.WorkflowTrigger
	for rows.Next() {
		var trigger models.WorkflowTrigger
		var configStr string
		err := rows.Scan(&trigger.ID, &trigger.WorkflowID, &trigger.TriggerType, &configStr, &trigger.CreatedAt, &trigger.UpdatedAt)
		if err != nil {
			continue
		}
		json.Unmarshal([]byte(configStr), &trigger.TriggerConfig)
		triggers = append(triggers, trigger)
	}

	return triggers, nil
}

// UpdateWorkflowTrigger updates a trigger
func (s *WorkflowService) UpdateWorkflowTrigger(triggerID int64, trigger *models.WorkflowTrigger) error {
	configJSON, _ := json.Marshal(trigger.TriggerConfig)

	query := `
		UPDATE workflow_triggers
		SET trigger_type = ?, trigger_config = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := s.db.Exec(query, trigger.TriggerType, string(configJSON), time.Now(), triggerID)
	if err != nil {
		return fmt.Errorf("failed to update workflow trigger: %w", err)
	}

	return nil
}

// DeleteWorkflowTrigger deletes a trigger
func (s *WorkflowService) DeleteWorkflowTrigger(triggerID int64) error {
	_, err := s.db.Exec("DELETE FROM workflow_triggers WHERE id = ?", triggerID)
	return err
}

// ==================== WORKFLOW ACTION METHODS ====================

// CreateWorkflowAction creates an action for a workflow
func (s *WorkflowService) CreateWorkflowAction(tenantID string, workflowID int64, action *models.WorkflowAction) (*models.WorkflowAction, error) {
	configJSON, _ := json.Marshal(action.ActionConfig)

	query := `
		INSERT INTO workflow_actions (workflow_id, action_type, action_config, action_order, delay_seconds, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, workflowID, action.ActionType, string(configJSON), action.Order, action.DelaySeconds, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow action: %w", err)
	}

	id, _ := result.LastInsertId()
	action.ID = id
	action.WorkflowID = workflowID

	return action, nil
}

// GetWorkflowActions retrieves all actions for a workflow
func (s *WorkflowService) GetWorkflowActions(tenantID string, workflowID int64) ([]models.WorkflowAction, error) {
	query := `
		SELECT id, workflow_id, action_type, action_config, action_order, delay_seconds, created_at, updated_at
		FROM workflow_actions
		WHERE workflow_id = ?
		ORDER BY action_order ASC
	`
	rows, err := s.db.Query(query, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow actions: %w", err)
	}
	defer rows.Close()

	var actions []models.WorkflowAction
	for rows.Next() {
		var action models.WorkflowAction
		var configStr string
		err := rows.Scan(&action.ID, &action.WorkflowID, &action.ActionType, &configStr, &action.Order, &action.DelaySeconds, &action.CreatedAt, &action.UpdatedAt)
		if err != nil {
			continue
		}
		json.Unmarshal([]byte(configStr), &action.ActionConfig)
		actions = append(actions, action)
	}

	return actions, nil
}

// UpdateWorkflowAction updates an action
func (s *WorkflowService) UpdateWorkflowAction(actionID int64, action *models.WorkflowAction) error {
	configJSON, _ := json.Marshal(action.ActionConfig)

	query := `
		UPDATE workflow_actions
		SET action_type = ?, action_config = ?, action_order = ?, delay_seconds = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := s.db.Exec(query, action.ActionType, string(configJSON), action.Order, action.DelaySeconds, time.Now(), actionID)
	if err != nil {
		return fmt.Errorf("failed to update workflow action: %w", err)
	}

	return nil
}

// DeleteWorkflowAction deletes an action
func (s *WorkflowService) DeleteWorkflowAction(actionID int64) error {
	_, err := s.db.Exec("DELETE FROM workflow_actions WHERE id = ?", actionID)
	return err
}

// ==================== WORKFLOW INSTANCE/EXECUTION METHODS ====================

// TriggerWorkflowInstance creates and starts a workflow execution
func (s *WorkflowService) TriggerWorkflowInstance(tenantID string, req *models.WorkflowInstanceRequest) (*models.WorkflowInstance, error) {
	instance := &models.WorkflowInstance{
		TenantID:         tenantID,
		WorkflowID:       req.WorkflowID,
		TriggeredBy:      req.TriggeredBy,
		TriggeredByValue: req.TriggeredByValue,
		Status:           "pending",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	query := `
		INSERT INTO workflow_instances (tenant_id, workflow_id, triggered_by, triggered_by_value, status, progress, executed_actions, failed_actions, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, instance.TenantID, instance.WorkflowID, instance.TriggeredBy, instance.TriggeredByValue,
		instance.Status, 0, 0, 0, instance.CreatedAt, instance.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow instance: %w", err)
	}

	id, _ := result.LastInsertId()
	instance.ID = id

	// Start execution asynchronously
	go s.executeWorkflowInstance(tenantID, instance)

	return instance, nil
}

// GetWorkflowInstance retrieves a workflow execution
func (s *WorkflowService) GetWorkflowInstance(tenantID string, instanceID int64) (*models.WorkflowInstance, error) {
	query := `
		SELECT id, tenant_id, workflow_id, triggered_by, triggered_by_value, status, progress, executed_actions, failed_actions, error_message, started_at, completed_at, created_at, updated_at
		FROM workflow_instances
		WHERE id = ? AND tenant_id = ?
	`
	instance := &models.WorkflowInstance{}
	err := s.db.QueryRow(query, instanceID, tenantID).Scan(
		&instance.ID, &instance.TenantID, &instance.WorkflowID, &instance.TriggeredBy, &instance.TriggeredByValue,
		&instance.Status, &instance.Progress, &instance.ExecutedActions, &instance.FailedActions, &instance.ErrorMessage,
		&instance.StartedAt, &instance.CompletedAt, &instance.CreatedAt, &instance.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow instance: %w", err)
	}

	return instance, nil
}

// ListWorkflowInstances lists execution history
func (s *WorkflowService) ListWorkflowInstances(tenantID string, workflowID int64, limit int, offset int) ([]models.WorkflowInstance, error) {
	query := `
		SELECT id, tenant_id, workflow_id, triggered_by, triggered_by_value, status, progress, executed_actions, failed_actions, error_message, started_at, completed_at, created_at, updated_at
		FROM workflow_instances
		WHERE tenant_id = ? AND workflow_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := s.db.Query(query, tenantID, workflowID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflow instances: %w", err)
	}
	defer rows.Close()

	var instances []models.WorkflowInstance
	for rows.Next() {
		var instance models.WorkflowInstance
		err := rows.Scan(
			&instance.ID, &instance.TenantID, &instance.WorkflowID, &instance.TriggeredBy, &instance.TriggeredByValue,
			&instance.Status, &instance.Progress, &instance.ExecutedActions, &instance.FailedActions, &instance.ErrorMessage,
			&instance.StartedAt, &instance.CompletedAt, &instance.CreatedAt, &instance.UpdatedAt,
		)
		if err != nil {
			continue
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// executeWorkflowInstance executes all actions in a workflow
func (s *WorkflowService) executeWorkflowInstance(tenantID string, instance *models.WorkflowInstance) error {
	// Mark as running
	startTime := time.Now()
	s.db.Exec("UPDATE workflow_instances SET status = ?, started_at = ?, updated_at = ? WHERE id = ?",
		"running", startTime, startTime, instance.ID)

	// Get workflow with actions
	workflow, err := s.GetWorkflow(tenantID, instance.WorkflowID)
	if err != nil {
		s.db.Exec("UPDATE workflow_instances SET status = ?, error_message = ?, updated_at = ? WHERE id = ?",
			"failed", err.Error(), time.Now(), instance.ID)
		return err
	}

	// Execute actions in order
	successCount := 0
	failureCount := 0

	for i, action := range workflow.Actions {
		progress := (i / len(workflow.Actions)) * 100

		// Create action execution record
		actionExec := &models.WorkflowActionExecution{
			WorkflowID: workflow.ID,
			InstanceID: instance.ID,
			ActionID:   action.ID,
			Status:     "executing",
			CreatedAt:  time.Now(),
		}

		if err := s.executeWorkflowAction(tenantID, instance, &action, actionExec); err != nil {
			failureCount++
			actionExec.Status = "failed"
			actionExec.ErrorMessage = err.Error()
		} else {
			successCount++
			actionExec.Status = "completed"
		}

		// Save action execution
		s.recordActionExecution(actionExec)

		// Update progress
		s.db.Exec("UPDATE workflow_instances SET progress = ?, executed_actions = ?, failed_actions = ?, updated_at = ? WHERE id = ?",
			progress, successCount, failureCount, time.Now(), instance.ID)
	}

	// Mark as completed
	completedTime := time.Now()
	finalStatus := "completed"
	if failureCount > 0 {
		finalStatus = "failed"
	}

	s.db.Exec("UPDATE workflow_instances SET status = ?, progress = ?, completed_at = ?, updated_at = ? WHERE id = ?",
		finalStatus, 100, completedTime, completedTime, instance.ID)

	return nil
}

// executeWorkflowAction executes a single action
func (s *WorkflowService) executeWorkflowAction(tenantID string, instance *models.WorkflowInstance, action *models.WorkflowAction, _ *models.WorkflowActionExecution) error {
	// Parse action config
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(action.ActionConfig), &config); err != nil {
		return fmt.Errorf("invalid action config: %w", err)
	}

	switch action.ActionType {
	case "create_task":
		return s.executeCreateTask(tenantID, instance, config)
	case "send_notification":
		return s.executeSendNotification(tenantID, instance, config)
	case "update_lead":
		return s.executeUpdateLead(tenantID, instance, config)
	case "add_tag":
		return s.executeAddTag(tenantID, instance, config)
	case "send_sms":
		return s.executeSendSMS(tenantID, instance, config)
	case "send_email":
		return s.executeSendEmail(tenantID, instance, config)
	default:
		return fmt.Errorf("unknown action type: %s", action.ActionType)
	}
}

// Helper methods for different action types
func (s *WorkflowService) executeCreateTask(tenantID string, _ *models.WorkflowInstance, config map[string]interface{}) error {
	// Extract task details from config
	title, _ := config["title"].(string)
	description, _ := config["description"].(string)
	assignedToID, _ := config["assigned_to_id"].(float64)

	query := `
		INSERT INTO tasks (tenant_id, title, description, assigned_to, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, tenantID, title, description, int64(assignedToID), "pending", time.Now(), time.Now())
	return err
}

func (s *WorkflowService) executeSendNotification(tenantID string, _ *models.WorkflowInstance, config map[string]interface{}) error {
	userID, _ := config["user_id"].(float64)
	message, _ := config["message"].(string)

	query := `
		INSERT INTO notifications (tenant_id, user_id, message, type, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, tenantID, int64(userID), message, "workflow", "unread", time.Now(), time.Now())
	return err
}

func (s *WorkflowService) executeUpdateLead(tenantID string, _ *models.WorkflowInstance, config map[string]interface{}) error {
	leadID, _ := config["lead_id"].(float64)
	status, _ := config["status"].(string)
	score, _ := config["score"].(float64)

	query := `UPDATE leads SET status = ?, lead_score = ?, updated_at = ? WHERE id = ? AND tenant_id = ?`
	_, err := s.db.Exec(query, status, int(score), time.Now(), int64(leadID), tenantID)
	return err
}

func (s *WorkflowService) executeAddTag(_ string, _ *models.WorkflowInstance, _ map[string]interface{}) error {
	// Implementation depends on your tag system
	// This is a placeholder
	return nil
}

func (s *WorkflowService) executeSendSMS(_ string, _ *models.WorkflowInstance, config map[string]interface{}) error {
	phoneNumber, _ := config["phone_number"].(string)
	message, _ := config["message"].(string)

	// Would integrate with SMS provider here
	fmt.Printf("Sending SMS to %s: %s\n", phoneNumber, message)
	return nil
}

func (s *WorkflowService) executeSendEmail(_ string, _ *models.WorkflowInstance, config map[string]interface{}) error {
	email, _ := config["email"].(string)
	subject, _ := config["subject"].(string)
	_, _ = config["body"].(string)

	// Would integrate with email provider here
	fmt.Printf("Sending email to %s: %s\n", email, subject)
	return nil
}

// recordActionExecution saves action execution record
func (s *WorkflowService) recordActionExecution(exec *models.WorkflowActionExecution) error {
	query := `
		INSERT INTO workflow_action_executions (workflow_id, instance_id, action_id, status, error_message, retry_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, exec.WorkflowID, exec.InstanceID, exec.ActionID, exec.Status, exec.ErrorMessage, exec.RetryCount, exec.CreatedAt, time.Now())
	return err
}

// ==================== SCHEDULED TASK METHODS ====================

// CreateScheduledTask creates a scheduled workflow/action
func (s *WorkflowService) CreateScheduledTask(tenantID string, task *models.ScheduledTask) (*models.ScheduledTask, error) {
	configJSON, _ := json.Marshal(task.Config)

	query := `
		INSERT INTO scheduled_tasks (tenant_id, name, type, config, schedule, next_run_at, enabled, max_retries, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, tenantID, task.Name, task.Type, string(configJSON), task.Schedule,
		task.NextRunAt, task.Enabled, task.MaxRetries, task.CreatedBy, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduled task: %w", err)
	}

	id, _ := result.LastInsertId()
	task.ID = id

	return task, nil
}

// GetScheduledTask retrieves a scheduled task
func (s *WorkflowService) GetScheduledTask(tenantID string, taskID int64) (*models.ScheduledTask, error) {
	query := `
		SELECT id, tenant_id, name, type, config, schedule, last_run_at, next_run_at, enabled, max_retries, created_by, created_at, updated_at
		FROM scheduled_tasks
		WHERE id = ? AND tenant_id = ?
	`
	task := &models.ScheduledTask{}
	var configStr string
	err := s.db.QueryRow(query, taskID, tenantID).Scan(
		&task.ID, &task.TenantID, &task.Name, &task.Type, &configStr, &task.Schedule,
		&task.LastRunAt, &task.NextRunAt, &task.Enabled, &task.MaxRetries, &task.CreatedBy, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get scheduled task: %w", err)
	}

	json.Unmarshal([]byte(configStr), &task.Config)
	return task, nil
}

// ListScheduledTasks lists all scheduled tasks
func (s *WorkflowService) ListScheduledTasks(tenantID string, limit int, offset int) ([]models.ScheduledTask, error) {
	query := `
		SELECT id, tenant_id, name, type, config, schedule, last_run_at, next_run_at, enabled, max_retries, created_by, created_at, updated_at
		FROM scheduled_tasks
		WHERE tenant_id = ?
		ORDER BY next_run_at ASC
		LIMIT ? OFFSET ?
	`
	rows, err := s.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list scheduled tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.ScheduledTask
	for rows.Next() {
		var task models.ScheduledTask
		var configStr string
		err := rows.Scan(
			&task.ID, &task.TenantID, &task.Name, &task.Type, &configStr, &task.Schedule,
			&task.LastRunAt, &task.NextRunAt, &task.Enabled, &task.MaxRetries, &task.CreatedBy, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			continue
		}
		json.Unmarshal([]byte(configStr), &task.Config)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateScheduledTask updates a scheduled task
func (s *WorkflowService) UpdateScheduledTask(taskID int64, task *models.ScheduledTask) error {
	configJSON, _ := json.Marshal(task.Config)

	query := `
		UPDATE scheduled_tasks
		SET name = ?, type = ?, config = ?, schedule = ?, enabled = ?, max_retries = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := s.db.Exec(query, task.Name, task.Type, string(configJSON), task.Schedule, task.Enabled, task.MaxRetries, time.Now(), taskID)
	return err
}

// DeleteScheduledTask deletes a scheduled task
func (s *WorkflowService) DeleteScheduledTask(taskID int64) error {
	_, err := s.db.Exec("DELETE FROM scheduled_tasks WHERE id = ?", taskID)
	return err
}

// GetWorkflowStats returns workflow statistics
func (s *WorkflowService) GetWorkflowStats(tenantID string, workflowID int64, days int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total executions
	query := `
		SELECT COUNT(*) as total, 
		       SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as successful,
		       SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed
		FROM workflow_instances
		WHERE tenant_id = ? AND workflow_id = ? AND created_at > DATE_SUB(NOW(), INTERVAL ? DAY)
	`
	var total, successful, failed int
	err := s.db.QueryRow(query, tenantID, workflowID, days).Scan(
		&total, &successful, &failed,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow stats: %w", err)
	}

	stats["total"] = total
	stats["successful"] = successful
	stats["failed"] = failed

	return stats, nil
}

// CountWorkflowsForTenant counts workflows for a tenant
func (s *WorkflowService) CountWorkflowsForTenant(tenantID string) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM workflows WHERE tenant_id = ?", tenantID).Scan(&count)
	return count, err
}

// EnableWorkflow enables/disables a workflow
func (s *WorkflowService) EnableWorkflow(tenantID string, workflowID int64, enabled bool) error {
	query := `UPDATE workflows SET enabled = ?, updated_at = ? WHERE id = ? AND tenant_id = ?`
	_, err := s.db.Exec(query, enabled, time.Now(), workflowID, tenantID)
	return err
}

// EvaluateTrigger checks if a trigger condition is met
func (s *WorkflowService) EvaluateTrigger(trigger *models.WorkflowTrigger, data map[string]interface{}) bool {
	// Parse trigger config and evaluate against data
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(trigger.TriggerConfig), &config); err != nil {
		return false
	}

	// Implementation depends on your condition format
	// This is a placeholder
	return true
}

// GetWorkflowByTriggerType gets workflows triggered by a specific event
func (s *WorkflowService) GetWorkflowByTriggerType(tenantID string, triggerType string) ([]models.WorkflowDefinition, error) {
	query := `
		SELECT DISTINCT w.id, w.tenant_id, w.name, w.description, w.enabled, w.created_by, w.created_at, w.updated_at
		FROM workflows w
		JOIN workflow_triggers wt ON w.id = wt.workflow_id
		WHERE w.tenant_id = ? AND wt.trigger_type = ? AND w.enabled = true
	`
	rows, err := s.db.Query(query, tenantID, triggerType)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflows by trigger type: %w", err)
	}
	defer rows.Close()

	var workflows []models.WorkflowDefinition
	for rows.Next() {
		var workflow models.WorkflowDefinition
		err := rows.Scan(
			&workflow.ID, &workflow.TenantID, &workflow.Name, &workflow.Description,
			&workflow.Enabled, &workflow.CreatedBy, &workflow.CreatedAt, &workflow.UpdatedAt,
		)
		if err != nil {
			continue
		}
		// Load triggers and actions
		triggers, _ := s.GetWorkflowTriggers(tenantID, workflow.ID)
		actions, _ := s.GetWorkflowActions(tenantID, workflow.ID)
		workflow.Triggers = triggers
		workflow.Actions = actions
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}
