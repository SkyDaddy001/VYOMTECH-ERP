package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TaskService defines task management operations
type TaskService interface {
	// Task CRUD operations
	CreateTask(ctx context.Context, tenantID string, task *Task) (*Task, error)
	GetTask(ctx context.Context, tenantID string, taskID int64) (*Task, error)
	GetTasksByUser(ctx context.Context, tenantID string, userID int64, status string, limit int) ([]Task, error)
	GetTasksByLead(ctx context.Context, tenantID string, leadID int64) ([]Task, error)
	GetOverdueTasks(ctx context.Context, tenantID string) ([]Task, error)
	UpdateTask(ctx context.Context, tenantID string, task *Task) (*Task, error)
	CompleteTask(ctx context.Context, tenantID string, taskID int64) error
	CancelTask(ctx context.Context, tenantID string, taskID int64) error
	DeleteTask(ctx context.Context, tenantID string, taskID int64) error

	// Task comments
	AddComment(ctx context.Context, tenantID string, comment *TaskComment) (*TaskComment, error)
	GetTaskComments(ctx context.Context, tenantID string, taskID int64) ([]TaskComment, error)

	// Task assignments
	AssignTask(ctx context.Context, tenantID string, assignment *TaskAssignment) (*TaskAssignment, error)
	RemoveAssignment(ctx context.Context, tenantID string, taskID, userID int64) error
	GetTaskAssignments(ctx context.Context, tenantID string, taskID int64) ([]TaskAssignment, error)

	// Statistics
	GetTaskStats(ctx context.Context, tenantID string, userID int64) (*TaskStats, error)
}

// Task represents a task in the system
type Task struct {
	ID                   int64      `json:"id"`
	AssignedTo           int64      `json:"assigned_to"`
	CreatedBy            int64      `json:"created_by"`
	LeadID               *int64     `json:"lead_id,omitempty"`
	TenantID             string     `json:"tenant_id"`
	Title                string     `json:"title"`
	Description          *string    `json:"description,omitempty"`
	Priority             string     `json:"priority"` // critical, high, normal, low
	Status               string     `json:"status"`   // pending, in_progress, completed, overdue, cancelled
	TaskType             *string    `json:"task_type,omitempty"`
	DueDate              time.Time  `json:"due_date"`
	ScheduledAt          *time.Time `json:"scheduled_at,omitempty"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	ParentTaskID         *int64     `json:"parent_task_id,omitempty"`
	RelatedEntityType    *string    `json:"related_entity_type,omitempty"`
	RelatedEntityID      *int64     `json:"related_entity_id,omitempty"`
	EstimatedDurationMin *int       `json:"estimated_duration_minutes,omitempty"`
	ActualDurationMin    *int       `json:"actual_duration_minutes,omitempty"`
	ProgressPercentage   int        `json:"progress_percentage"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// TaskComment represents a comment on a task
type TaskComment struct {
	ID            int64     `json:"id"`
	TaskID        int64     `json:"task_id"`
	UserID        int64     `json:"user_id"`
	TenantID      string    `json:"tenant_id"`
	CommentText   string    `json:"comment_text"`
	AttachmentURL *string   `json:"attachment_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TaskAssignment represents a task assignment
type TaskAssignment struct {
	ID                       int64      `json:"id"`
	TaskID                   int64      `json:"task_id"`
	UserID                   int64      `json:"user_id"`
	TenantID                 string     `json:"tenant_id"`
	AssignmentRole           string     `json:"assignment_role"` // assignee, reviewer, watcher
	ResponsibilityPercentage int        `json:"responsibility_percentage"`
	AssignedAt               time.Time  `json:"assigned_at"`
	AcceptedAt               *time.Time `json:"accepted_at,omitempty"`
	CompletedAt              *time.Time `json:"completed_at,omitempty"`
}

// TaskStats represents task statistics
type TaskStats struct {
	TotalTasks            int     `json:"total_tasks"`
	PendingTasks          int     `json:"pending_tasks"`
	InProgressTasks       int     `json:"in_progress_tasks"`
	CompletedTasks        int     `json:"completed_tasks"`
	OverdueTasks          int     `json:"overdue_tasks"`
	CancelledTasks        int     `json:"cancelled_tasks"`
	CompletionRate        float64 `json:"completion_rate"`
	AverageCompletionTime float64 `json:"average_completion_time_hours"`
	HighPriorityTasks     int     `json:"high_priority_tasks"`
	CriticalTasks         int     `json:"critical_tasks"`
}

// Task status constants
const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
	TaskStatusOverdue    = "overdue"
	TaskStatusCancelled  = "cancelled"
)

// Task priority constants
const (
	TaskPriorityCritical = "critical"
	TaskPriorityHigh     = "high"
	TaskPriorityNormal   = "normal"
	TaskPriorityLow      = "low"
)

// taskService implements TaskService
type taskService struct {
	db *sql.DB
}

// NewTaskService creates a new task service
func NewTaskService(db *sql.DB) TaskService {
	return &taskService{db: db}
}

// CreateTask creates a new task
func (s *taskService) CreateTask(ctx context.Context, tenantID string, task *Task) (*Task, error) {
	query := `
		INSERT INTO tasks (assigned_to, created_by, lead_id, tenant_id, title, description, 
		                   priority, status, task_type, due_date, scheduled_at, parent_task_id,
		                   related_entity_type, related_entity_id, estimated_duration_minutes, progress_percentage)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		task.AssignedTo, task.CreatedBy, task.LeadID, tenantID, task.Title, task.Description,
		task.Priority, TaskStatusPending, task.TaskType, task.DueDate, task.ScheduledAt,
		task.ParentTaskID, task.RelatedEntityType, task.RelatedEntityID,
		task.EstimatedDurationMin, 0)

	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get task ID: %w", err)
	}

	task.ID = id
	task.Status = TaskStatusPending
	task.ProgressPercentage = 0
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	return task, nil
}

// GetTask retrieves a task by ID
func (s *taskService) GetTask(ctx context.Context, tenantID string, taskID int64) (*Task, error) {
	query := `
		SELECT id, assigned_to, created_by, lead_id, tenant_id, title, description, 
		       priority, status, task_type, due_date, scheduled_at, completed_at,
		       parent_task_id, related_entity_type, related_entity_id, 
		       estimated_duration_minutes, actual_duration_minutes, progress_percentage,
		       created_at, updated_at
		FROM tasks
		WHERE id = ? AND tenant_id = ?
	`

	task := &Task{}
	err := s.db.QueryRowContext(ctx, query, taskID, tenantID).Scan(
		&task.ID, &task.AssignedTo, &task.CreatedBy, &task.LeadID, &task.TenantID,
		&task.Title, &task.Description, &task.Priority, &task.Status, &task.TaskType,
		&task.DueDate, &task.ScheduledAt, &task.CompletedAt, &task.ParentTaskID,
		&task.RelatedEntityType, &task.RelatedEntityID, &task.EstimatedDurationMin,
		&task.ActualDurationMin, &task.ProgressPercentage, &task.CreatedAt, &task.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("task not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// GetTasksByUser retrieves tasks assigned to a user
func (s *taskService) GetTasksByUser(ctx context.Context, tenantID string, userID int64, status string, limit int) ([]Task, error) {
	query := `
		SELECT id, assigned_to, created_by, lead_id, tenant_id, title, description, 
		       priority, status, task_type, due_date, scheduled_at, completed_at,
		       parent_task_id, related_entity_type, related_entity_id, 
		       estimated_duration_minutes, actual_duration_minutes, progress_percentage,
		       created_at, updated_at
		FROM tasks
		WHERE assigned_to = ? AND tenant_id = ?
	`
	args := []interface{}{userID, tenantID}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY due_date ASC LIMIT ?"
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		task := Task{}
		err := rows.Scan(
			&task.ID, &task.AssignedTo, &task.CreatedBy, &task.LeadID, &task.TenantID,
			&task.Title, &task.Description, &task.Priority, &task.Status, &task.TaskType,
			&task.DueDate, &task.ScheduledAt, &task.CompletedAt, &task.ParentTaskID,
			&task.RelatedEntityType, &task.RelatedEntityID, &task.EstimatedDurationMin,
			&task.ActualDurationMin, &task.ProgressPercentage, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTasksByLead retrieves tasks for a lead
func (s *taskService) GetTasksByLead(ctx context.Context, tenantID string, leadID int64) ([]Task, error) {
	query := `
		SELECT id, assigned_to, created_by, lead_id, tenant_id, title, description, 
		       priority, status, task_type, due_date, scheduled_at, completed_at,
		       parent_task_id, related_entity_type, related_entity_id, 
		       estimated_duration_minutes, actual_duration_minutes, progress_percentage,
		       created_at, updated_at
		FROM tasks
		WHERE lead_id = ? AND tenant_id = ?
		ORDER BY due_date ASC
	`

	rows, err := s.db.QueryContext(ctx, query, leadID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		task := Task{}
		err := rows.Scan(
			&task.ID, &task.AssignedTo, &task.CreatedBy, &task.LeadID, &task.TenantID,
			&task.Title, &task.Description, &task.Priority, &task.Status, &task.TaskType,
			&task.DueDate, &task.ScheduledAt, &task.CompletedAt, &task.ParentTaskID,
			&task.RelatedEntityType, &task.RelatedEntityID, &task.EstimatedDurationMin,
			&task.ActualDurationMin, &task.ProgressPercentage, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetOverdueTasks retrieves overdue tasks
func (s *taskService) GetOverdueTasks(ctx context.Context, tenantID string) ([]Task, error) {
	query := `
		SELECT id, assigned_to, created_by, lead_id, tenant_id, title, description, 
		       priority, status, task_type, due_date, scheduled_at, completed_at,
		       parent_task_id, related_entity_type, related_entity_id, 
		       estimated_duration_minutes, actual_duration_minutes, progress_percentage,
		       created_at, updated_at
		FROM tasks
		WHERE due_date < NOW() AND status IN (?, ?) AND tenant_id = ?
		ORDER BY due_date ASC
	`

	rows, err := s.db.QueryContext(ctx, query, TaskStatusPending, TaskStatusInProgress, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		task := Task{}
		err := rows.Scan(
			&task.ID, &task.AssignedTo, &task.CreatedBy, &task.LeadID, &task.TenantID,
			&task.Title, &task.Description, &task.Priority, &task.Status, &task.TaskType,
			&task.DueDate, &task.ScheduledAt, &task.CompletedAt, &task.ParentTaskID,
			&task.RelatedEntityType, &task.RelatedEntityID, &task.EstimatedDurationMin,
			&task.ActualDurationMin, &task.ProgressPercentage, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateTask updates a task
func (s *taskService) UpdateTask(ctx context.Context, tenantID string, task *Task) (*Task, error) {
	query := `
		UPDATE tasks
		SET title = ?, description = ?, priority = ?, status = ?, 
		    due_date = ?, progress_percentage = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query,
		task.Title, task.Description, task.Priority, task.Status,
		task.DueDate, task.ProgressPercentage, task.ID, tenantID)

	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	task.UpdatedAt = time.Now()
	return task, nil
}

// CompleteTask marks a task as completed
func (s *taskService) CompleteTask(ctx context.Context, tenantID string, taskID int64) error {
	now := time.Now()
	query := `
		UPDATE tasks
		SET status = ?, completed_at = ?, progress_percentage = 100, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query, TaskStatusCompleted, now, taskID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}

	return nil
}

// CancelTask cancels a task
func (s *taskService) CancelTask(ctx context.Context, tenantID string, taskID int64) error {
	query := `
		UPDATE tasks
		SET status = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?
	`

	_, err := s.db.ExecContext(ctx, query, TaskStatusCancelled, taskID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to cancel task: %w", err)
	}

	return nil
}

// DeleteTask deletes a task
func (s *taskService) DeleteTask(ctx context.Context, tenantID string, taskID int64) error {
	query := `DELETE FROM tasks WHERE id = ? AND tenant_id = ?`

	_, err := s.db.ExecContext(ctx, query, taskID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// AddComment adds a comment to a task
func (s *taskService) AddComment(ctx context.Context, tenantID string, comment *TaskComment) (*TaskComment, error) {
	query := `
		INSERT INTO task_comments (task_id, user_id, tenant_id, comment_text, attachment_url)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		comment.TaskID, comment.UserID, tenantID, comment.CommentText, comment.AttachmentURL)

	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get comment ID: %w", err)
	}

	comment.ID = id
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	return comment, nil
}

// GetTaskComments retrieves comments for a task
func (s *taskService) GetTaskComments(ctx context.Context, tenantID string, taskID int64) ([]TaskComment, error) {
	query := `
		SELECT id, task_id, user_id, tenant_id, comment_text, attachment_url, created_at, updated_at
		FROM task_comments
		WHERE task_id = ? AND tenant_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, taskID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []TaskComment
	for rows.Next() {
		comment := TaskComment{}
		err := rows.Scan(&comment.ID, &comment.TaskID, &comment.UserID, &comment.TenantID,
			&comment.CommentText, &comment.AttachmentURL, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// AssignTask assigns a task to a user
func (s *taskService) AssignTask(ctx context.Context, tenantID string, assignment *TaskAssignment) (*TaskAssignment, error) {
	query := `
		INSERT INTO task_assignments (task_id, user_id, tenant_id, assignment_role, responsibility_percentage)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		assignment.TaskID, assignment.UserID, tenantID, assignment.AssignmentRole, assignment.ResponsibilityPercentage)

	if err != nil {
		return nil, fmt.Errorf("failed to assign task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment ID: %w", err)
	}

	assignment.ID = id
	assignment.AssignedAt = time.Now()

	return assignment, nil
}

// RemoveAssignment removes a task assignment
func (s *taskService) RemoveAssignment(ctx context.Context, tenantID string, taskID, userID int64) error {
	query := `DELETE FROM task_assignments WHERE task_id = ? AND user_id = ? AND tenant_id = ?`

	_, err := s.db.ExecContext(ctx, query, taskID, userID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to remove assignment: %w", err)
	}

	return nil
}

// GetTaskAssignments retrieves assignments for a task
func (s *taskService) GetTaskAssignments(ctx context.Context, tenantID string, taskID int64) ([]TaskAssignment, error) {
	query := `
		SELECT id, task_id, user_id, tenant_id, assignment_role, responsibility_percentage, 
		       assigned_at, accepted_at, completed_at
		FROM task_assignments
		WHERE task_id = ? AND tenant_id = ?
	`

	rows, err := s.db.QueryContext(ctx, query, taskID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments: %w", err)
	}
	defer rows.Close()

	var assignments []TaskAssignment
	for rows.Next() {
		assignment := TaskAssignment{}
		err := rows.Scan(&assignment.ID, &assignment.TaskID, &assignment.UserID, &assignment.TenantID,
			&assignment.AssignmentRole, &assignment.ResponsibilityPercentage,
			&assignment.AssignedAt, &assignment.AcceptedAt, &assignment.CompletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment: %w", err)
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

// GetTaskStats retrieves task statistics for a user
func (s *taskService) GetTaskStats(ctx context.Context, tenantID string, userID int64) (*TaskStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as in_progress,
			SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN status IN (?, ?) AND due_date < NOW() THEN 1 ELSE 0 END) as overdue,
			SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as cancelled,
			SUM(CASE WHEN priority = ? THEN 1 ELSE 0 END) as high_priority,
			SUM(CASE WHEN priority = ? THEN 1 ELSE 0 END) as critical
		FROM tasks
		WHERE assigned_to = ? AND tenant_id = ?
	`

	stats := &TaskStats{}
	var total, pending, inProgress, completed, overdue, cancelled, highPriority, critical sql.NullInt64

	err := s.db.QueryRowContext(ctx, query,
		TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted,
		TaskStatusPending, TaskStatusInProgress, TaskStatusCancelled,
		TaskPriorityHigh, TaskPriorityCritical, userID, tenantID).Scan(
		&total, &pending, &inProgress, &completed, &overdue, &cancelled, &highPriority, &critical)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get task stats: %w", err)
	}

	if total.Valid {
		stats.TotalTasks = int(total.Int64)
	}
	if pending.Valid {
		stats.PendingTasks = int(pending.Int64)
	}
	if inProgress.Valid {
		stats.InProgressTasks = int(inProgress.Int64)
	}
	if completed.Valid {
		stats.CompletedTasks = int(completed.Int64)
	}
	if overdue.Valid {
		stats.OverdueTasks = int(overdue.Int64)
	}
	if cancelled.Valid {
		stats.CancelledTasks = int(cancelled.Int64)
	}
	if highPriority.Valid {
		stats.HighPriorityTasks = int(highPriority.Int64)
	}
	if critical.Valid {
		stats.CriticalTasks = int(critical.Int64)
	}

	if stats.TotalTasks > 0 {
		stats.CompletionRate = float64(stats.CompletedTasks) / float64(stats.TotalTasks) * 100
	}

	return stats, nil
}
