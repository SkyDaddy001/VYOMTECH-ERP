package handlers

import (
	"vyomtech-backend/internal/middleware"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"

	"github.com/gorilla/mux"
)

// TaskHandler handles task management endpoints
type TaskHandler struct {
	taskService services.TaskService
	log         *logger.Logger
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(taskService services.TaskService, log *logger.Logger) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		log:         log,
	}
}

// RegisterRoutes registers task routes
func (h *TaskHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/tasks", h.CreateTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks", h.ListTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", h.GetTask).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", h.UpdateTask).Methods(http.MethodPut)
	router.HandleFunc("/tasks/{id}", h.DeleteTask).Methods(http.MethodDelete)
	router.HandleFunc("/tasks/{id}/complete", h.CompleteTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks/user/{userID}", h.GetTasksByUser).Methods(http.MethodGet)
	router.HandleFunc("/tasks/stats", h.GetTaskStats).Methods(http.MethodGet)
}

// CreateTask creates a new task
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	var task services.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	task.TenantID = tenantID
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if ok {
		task.CreatedBy = userID
	}

	created, err := h.taskService.CreateTask(ctx, tenantID, &task)
	if err != nil {
		h.log.Error("Failed to create task", "error", err)
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetTask retrieves a specific task
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	taskID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetTask(ctx, tenantID, taskID)
	if err != nil {
		h.log.Error("Failed to get task", "error", err)
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// ListTasks lists tasks with optional filters
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	assignedTo := r.URL.Query().Get("assigned_to")

	var tasks []services.Task
	var err error

	if assignedTo != "" {
		userID, parseErr := strconv.ParseInt(assignedTo, 10, 64)
		if parseErr != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}
		tasks, err = h.taskService.GetTasksByUser(ctx, tenantID, userID, status, 100)
	} else {
		tasks = []services.Task{}
	}

	if err != nil {
		h.log.Error("Failed to list tasks", "error", err)
		http.Error(w, "failed to list tasks", http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []services.Task{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(tasks),
		"tasks": tasks,
	})
}

// UpdateTask updates an existing task
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	taskID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	var task services.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	task.ID = taskID
	task.TenantID = tenantID

	updated, err := h.taskService.UpdateTask(ctx, tenantID, &task)
	if err != nil {
		h.log.Error("Failed to update task", "error", err)
		http.Error(w, "failed to update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteTask deletes a task
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	taskID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	if err := h.taskService.DeleteTask(ctx, tenantID, taskID); err != nil {
		h.log.Error("Failed to delete task", "error", err)
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "task deleted"})
}

// CompleteTask marks a task as completed
func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	taskID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	if err := h.taskService.CompleteTask(ctx, tenantID, taskID); err != nil {
		h.log.Error("Failed to complete task", "error", err)
		http.Error(w, "failed to complete task", http.StatusInternalServerError)
		return
	}

	// Fetch updated task
	task, err := h.taskService.GetTask(ctx, tenantID, taskID)
	if err != nil {
		h.log.Error("Failed to fetch completed task", "error", err)
		http.Error(w, "failed to fetch task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// GetTasksByUser gets tasks assigned to a specific user
func (h *TaskHandler) GetTasksByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userID"], 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	tasks, err := h.taskService.GetTasksByUser(ctx, tenantID, userID, "", 100)
	if err != nil {
		h.log.Error("Failed to get user tasks", "error", err)
		http.Error(w, "failed to get tasks", http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []services.Task{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(tasks),
		"tasks": tasks,
	})
}

// GetTaskStats gets task statistics for current user
func (h *TaskHandler) GetTaskStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		userID = 0
	}

	stats, err := h.taskService.GetTaskStats(ctx, tenantID, userID)
	if err != nil {
		h.log.Error("Failed to get task stats", "error", err)
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
