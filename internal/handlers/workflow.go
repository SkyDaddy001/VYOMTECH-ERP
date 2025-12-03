package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

// ==================== WORKFLOW HANDLER ====================

type WorkflowHandler struct {
	service *services.WorkflowService
}

// NewWorkflowHandler creates a new workflow handler
func NewWorkflowHandler(service *services.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		service: service,
	}
}

// RegisterRoutes registers all workflow routes
func (h *WorkflowHandler) RegisterRoutes(router *mux.Router) {
	// Workflow definition routes
	router.HandleFunc("/api/v1/workflows", h.CreateWorkflow).Methods("POST")
	router.HandleFunc("/api/v1/workflows", h.ListWorkflows).Methods("GET")
	router.HandleFunc("/api/v1/workflows/{id}", h.GetWorkflow).Methods("GET")
	router.HandleFunc("/api/v1/workflows/{id}", h.UpdateWorkflow).Methods("PUT")
	router.HandleFunc("/api/v1/workflows/{id}", h.DeleteWorkflow).Methods("DELETE")
	router.HandleFunc("/api/v1/workflows/{id}/enable", h.EnableWorkflow).Methods("PATCH")

	// Workflow trigger routes
	router.HandleFunc("/api/v1/workflows/{workflowId}/triggers", h.CreateTrigger).Methods("POST")
	router.HandleFunc("/api/v1/workflows/{workflowId}/triggers", h.ListTriggers).Methods("GET")
	router.HandleFunc("/api/v1/workflows/{workflowId}/triggers/{triggerId}", h.UpdateTrigger).Methods("PUT")
	router.HandleFunc("/api/v1/workflows/{workflowId}/triggers/{triggerId}", h.DeleteTrigger).Methods("DELETE")

	// Workflow action routes
	router.HandleFunc("/api/v1/workflows/{workflowId}/actions", h.CreateAction).Methods("POST")
	router.HandleFunc("/api/v1/workflows/{workflowId}/actions", h.ListActions).Methods("GET")
	router.HandleFunc("/api/v1/workflows/{workflowId}/actions/{actionId}", h.UpdateAction).Methods("PUT")
	router.HandleFunc("/api/v1/workflows/{workflowId}/actions/{actionId}", h.DeleteAction).Methods("DELETE")

	// Workflow execution routes
	router.HandleFunc("/api/v1/workflows/{workflowId}/trigger", h.TriggerWorkflow).Methods("POST")
	router.HandleFunc("/api/v1/workflows/{workflowId}/instances", h.ListInstances).Methods("GET")
	router.HandleFunc("/api/v1/workflows/{workflowId}/instances/{instanceId}", h.GetInstance).Methods("GET")

	// Workflow statistics
	router.HandleFunc("/api/v1/workflows/{workflowId}/stats", h.GetWorkflowStats).Methods("GET")

	// Scheduled tasks routes
	router.HandleFunc("/api/v1/scheduled-tasks", h.CreateScheduledTask).Methods("POST")
	router.HandleFunc("/api/v1/scheduled-tasks", h.ListScheduledTasks).Methods("GET")
	router.HandleFunc("/api/v1/scheduled-tasks/{taskId}", h.GetScheduledTask).Methods("GET")
	router.HandleFunc("/api/v1/scheduled-tasks/{taskId}", h.UpdateScheduledTask).Methods("PUT")
	router.HandleFunc("/api/v1/scheduled-tasks/{taskId}", h.DeleteScheduledTask).Methods("DELETE")
}

// ==================== WORKFLOW DEFINITION HANDLERS ====================

// CreateWorkflow creates a new workflow
// POST /api/v1/workflows
func (h *WorkflowHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	userID := r.Header.Get("X-User-ID")
	if tenantID == "" || userID == "" {
		http.Error(w, "Missing tenant or user ID", http.StatusBadRequest)
		return
	}

	var req models.WorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	workflow, err := h.service.CreateWorkflow(tenantID, &req, userIDInt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create workflow: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workflow)
}

// ListWorkflows lists all workflows for a tenant
// GET /api/v1/workflows?limit=10&offset=0
func (h *WorkflowHandler) ListWorkflows(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	limit := 10
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	workflows, err := h.service.ListWorkflows(tenantID, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list workflows: %v", err), http.StatusInternalServerError)
		return
	}

	if workflows == nil {
		workflows = []models.WorkflowDefinition{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  workflows,
		"count": len(workflows),
	})
}

// GetWorkflow retrieves a specific workflow
// GET /api/v1/workflows/{id}
func (h *WorkflowHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	workflowID, _ := strconv.ParseInt(vars["id"], 10, 64)

	workflow, err := h.service.GetWorkflow(tenantID, workflowID)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

// UpdateWorkflow updates a workflow
// PUT /api/v1/workflows/{id}
func (h *WorkflowHandler) UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	workflowID, _ := strconv.ParseInt(vars["id"], 10, 64)

	var req models.WorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workflow, err := h.service.UpdateWorkflow(tenantID, workflowID, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update workflow: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

// DeleteWorkflow deletes a workflow
// DELETE /api/v1/workflows/{id}
func (h *WorkflowHandler) DeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	workflowID, _ := strconv.ParseInt(vars["id"], 10, 64)

	if err := h.service.DeleteWorkflow(tenantID, workflowID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete workflow: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// EnableWorkflow enables or disables a workflow
// PATCH /api/v1/workflows/{id}/enable?enabled=true
func (h *WorkflowHandler) EnableWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	workflowID, _ := strconv.ParseInt(vars["id"], 10, 64)

	enabled := r.URL.Query().Get("enabled") == "true"

	if err := h.service.EnableWorkflow(tenantID, workflowID, enabled); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update workflow: %v", err), http.StatusInternalServerError)
		return
	}

	workflow, _ := h.service.GetWorkflow(tenantID, workflowID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

// ==================== WORKFLOW TRIGGER HANDLERS ====================

// CreateTrigger creates a trigger for a workflow
// POST /api/v1/workflows/{workflowId}/triggers
func (h *WorkflowHandler) CreateTrigger(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	workflowID, _ := strconv.ParseInt(vars["workflowId"], 10, 64)

	var trigger models.WorkflowTrigger
	if err := json.NewDecoder(r.Body).Decode(&trigger); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.CreateWorkflowTrigger(tenantID, workflowID, &trigger)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create trigger: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// ListTriggers lists all triggers for a workflow
// GET /api/v1/workflows/{workflowId}/triggers
func (h *WorkflowHandler) ListTriggers(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	workflowID, _ := strconv.ParseInt(vars["workflowId"], 10, 64)

	triggers, err := h.service.GetWorkflowTriggers(tenantID, workflowID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list triggers: %v", err), http.StatusInternalServerError)
		return
	}

	if triggers == nil {
		triggers = []models.WorkflowTrigger{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  triggers,
		"count": len(triggers),
	})
}

// UpdateTrigger updates a trigger
// PUT /api/v1/workflows/{workflowId}/triggers/{triggerId}
func (h *WorkflowHandler) UpdateTrigger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	triggerID, _ := strconv.ParseInt(vars["triggerId"], 10, 64)

	var trigger models.WorkflowTrigger
	if err := json.NewDecoder(r.Body).Decode(&trigger); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateWorkflowTrigger(triggerID, &trigger); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update trigger: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trigger)
}

// DeleteTrigger deletes a trigger
// DELETE /api/v1/workflows/{workflowId}/triggers/{triggerId}
func (h *WorkflowHandler) DeleteTrigger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	triggerID, _ := strconv.ParseInt(vars["triggerId"], 10, 64)

	if err := h.service.DeleteWorkflowTrigger(triggerID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete trigger: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// ==================== WORKFLOW ACTION HANDLERS ====================

// CreateAction creates an action for a workflow
// POST /api/v1/workflows/{workflowId}/actions
func (h *WorkflowHandler) CreateAction(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	workflowID, _ := strconv.ParseInt(vars["workflowId"], 10, 64)

	var action models.WorkflowAction
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.CreateWorkflowAction(tenantID, workflowID, &action)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create action: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// ListActions lists all actions for a workflow
// GET /api/v1/workflows/{workflowId}/actions
func (h *WorkflowHandler) ListActions(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	workflowID, _ := strconv.ParseInt(vars["workflowId"], 10, 64)

	actions, err := h.service.GetWorkflowActions(tenantID, workflowID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list actions: %v", err), http.StatusInternalServerError)
		return
	}

	if actions == nil {
		actions = []models.WorkflowAction{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  actions,
		"count": len(actions),
	})
}

// UpdateAction updates an action
// PUT /api/v1/workflows/{workflowId}/actions/{actionId}
func (h *WorkflowHandler) UpdateAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	actionID, _ := strconv.ParseInt(vars["actionId"], 10, 64)

	var action models.WorkflowAction
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateWorkflowAction(actionID, &action); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update action: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(action)
}

// DeleteAction deletes an action
// DELETE /api/v1/workflows/{workflowId}/actions/{actionId}
func (h *WorkflowHandler) DeleteAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	actionID, _ := strconv.ParseInt(vars["actionId"], 10, 64)

	if err := h.service.DeleteWorkflowAction(actionID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete action: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// ==================== WORKFLOW EXECUTION HANDLERS ====================

// TriggerWorkflow triggers a workflow execution
// POST /api/v1/workflows/{workflowId}/trigger
func (h *WorkflowHandler) TriggerWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	workflowID, _ := strconv.ParseInt(vars["workflowId"], 10, 64)

	var req models.WorkflowInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.WorkflowID = workflowID
	instance, err := h.service.TriggerWorkflowInstance(tenantID, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to trigger workflow: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(instance)
}

// ListInstances lists execution instances of a workflow
// GET /api/v1/workflows/{workflowId}/instances?limit=10&offset=0
func (h *WorkflowHandler) ListInstances(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	workflowID, _ := strconv.ParseInt(vars["workflowId"], 10, 64)

	limit := 10
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	instances, err := h.service.ListWorkflowInstances(tenantID, workflowID, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list instances: %v", err), http.StatusInternalServerError)
		return
	}

	if instances == nil {
		instances = []models.WorkflowInstance{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  instances,
		"count": len(instances),
	})
}

// GetInstance retrieves a specific workflow execution
// GET /api/v1/workflows/{workflowId}/instances/{instanceId}
func (h *WorkflowHandler) GetInstance(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	instanceID, _ := strconv.ParseInt(vars["instanceId"], 10, 64)

	instance, err := h.service.GetWorkflowInstance(tenantID, instanceID)
	if err != nil {
		http.Error(w, "Instance not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instance)
}

// ==================== STATISTICS HANDLER ====================

// GetWorkflowStats returns workflow execution statistics
// GET /api/v1/workflows/{workflowId}/stats?days=30
func (h *WorkflowHandler) GetWorkflowStats(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	workflowID, _ := strconv.ParseInt(vars["workflowId"], 10, 64)

	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		days, _ = strconv.Atoi(d)
	}

	stats, err := h.service.GetWorkflowStats(tenantID, workflowID, days)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ==================== SCHEDULED TASK HANDLERS ====================

// CreateScheduledTask creates a scheduled task
// POST /api/v1/scheduled-tasks
func (h *WorkflowHandler) CreateScheduledTask(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	userID := r.Header.Get("X-User-ID")
	if tenantID == "" || userID == "" {
		http.Error(w, "Missing tenant or user ID", http.StatusBadRequest)
		return
	}

	var task models.ScheduledTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task.TenantID = tenantID
	userIDInt, _ := strconv.ParseInt(userID, 10, 64)
	task.CreatedBy = userIDInt

	result, err := h.service.CreateScheduledTask(tenantID, &task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create scheduled task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// ListScheduledTasks lists all scheduled tasks
// GET /api/v1/scheduled-tasks?limit=10&offset=0
func (h *WorkflowHandler) ListScheduledTasks(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	limit := 10
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	tasks, err := h.service.ListScheduledTasks(tenantID, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list tasks: %v", err), http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []models.ScheduledTask{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  tasks,
		"count": len(tasks),
	})
}

// GetScheduledTask retrieves a specific scheduled task
// GET /api/v1/scheduled-tasks/{taskId}
func (h *WorkflowHandler) GetScheduledTask(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing tenant ID", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	taskID, _ := strconv.ParseInt(vars["taskId"], 10, 64)

	task, err := h.service.GetScheduledTask(tenantID, taskID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// UpdateScheduledTask updates a scheduled task
// PUT /api/v1/scheduled-tasks/{taskId}
func (h *WorkflowHandler) UpdateScheduledTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, _ := strconv.ParseInt(vars["taskId"], 10, 64)

	var task models.ScheduledTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateScheduledTask(taskID, &task); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteScheduledTask deletes a scheduled task
// DELETE /api/v1/scheduled-tasks/{taskId}
func (h *WorkflowHandler) DeleteScheduledTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, _ := strconv.ParseInt(vars["taskId"], 10, 64)

	if err := h.service.DeleteScheduledTask(taskID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
