package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

// CustomizationHandler handles tenant customization endpoints
type CustomizationHandler struct {
	customizationService services.TenantCustomizationService
}

// NewCustomizationHandler creates a new customization handler
func NewCustomizationHandler(customizationService services.TenantCustomizationService) *CustomizationHandler {
	return &CustomizationHandler{
		customizationService: customizationService,
	}
}

// RegisterRoutes registers customization routes
func (h *CustomizationHandler) RegisterRoutes(router *mux.Router) {
	// Task Status endpoints
	router.HandleFunc("/task-statuses", h.CreateTaskStatus).Methods(http.MethodPost)
	router.HandleFunc("/task-statuses", h.GetTaskStatuses).Methods(http.MethodGet)
	router.HandleFunc("/task-statuses/{statusCode}", h.GetTaskStatus).Methods(http.MethodGet)
	router.HandleFunc("/task-statuses/{statusCode}", h.UpdateTaskStatus).Methods(http.MethodPut)
	router.HandleFunc("/task-statuses/{statusCode}", h.DeactivateTaskStatus).Methods(http.MethodDelete)

	// Task Stage endpoints
	router.HandleFunc("/task-stages", h.CreateTaskStage).Methods(http.MethodPost)
	router.HandleFunc("/task-stages", h.GetTaskStages).Methods(http.MethodGet)
	router.HandleFunc("/task-stages/{stageCode}", h.UpdateTaskStage).Methods(http.MethodPut)

	// Status Transition endpoints
	router.HandleFunc("/status-transitions", h.CreateStatusTransition).Methods(http.MethodPost)
	router.HandleFunc("/status-transitions", h.GetStatusTransitions).Methods(http.MethodGet)
	router.HandleFunc("/status-transitions/check", h.CheckStatusTransition).Methods(http.MethodGet)

	// Task Type endpoints
	router.HandleFunc("/task-types", h.CreateTaskType).Methods(http.MethodPost)
	router.HandleFunc("/task-types", h.GetTaskTypes).Methods(http.MethodGet)
	router.HandleFunc("/task-types/{typeCode}", h.UpdateTaskType).Methods(http.MethodPut)

	// Priority Level endpoints
	router.HandleFunc("/priority-levels", h.CreatePriorityLevel).Methods(http.MethodPost)
	router.HandleFunc("/priority-levels", h.GetPriorityLevels).Methods(http.MethodGet)
	router.HandleFunc("/priority-levels/{priorityCode}", h.UpdatePriorityLevel).Methods(http.MethodPut)

	// Notification Type endpoints
	router.HandleFunc("/notification-types", h.CreateNotificationType).Methods(http.MethodPost)
	router.HandleFunc("/notification-types", h.GetNotificationTypes).Methods(http.MethodGet)
	router.HandleFunc("/notification-types/{typeCode}", h.UpdateNotificationType).Methods(http.MethodPut)

	// Custom Field endpoints
	router.HandleFunc("/custom-fields", h.CreateCustomField).Methods(http.MethodPost)
	router.HandleFunc("/custom-fields", h.GetCustomFields).Methods(http.MethodGet)
	router.HandleFunc("/custom-fields/{fieldCode}", h.UpdateCustomField).Methods(http.MethodPut)

	// Automation Rule endpoints
	router.HandleFunc("/automation-rules", h.CreateAutomationRule).Methods(http.MethodPost)
	router.HandleFunc("/automation-rules", h.GetAutomationRules).Methods(http.MethodGet)
	router.HandleFunc("/automation-rules/{ruleCode}", h.UpdateAutomationRule).Methods(http.MethodPut)

	// Complete Configuration endpoint
	router.HandleFunc("/all", h.GetTenantConfiguration).Methods(http.MethodGet)
}

// Helper function to get tenant from context
func getTenantID(r *http.Request) string {
	tenantID, ok := r.Context().Value("tenantID").(string)
	if !ok {
		return ""
	}
	return tenantID
}

// Helper function to write JSON response
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Helper function to write error response
func writeError(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]string{"error": message}
	writeJSON(w, statusCode, response)
}

// ============================================================================
// TASK STATUS HANDLERS
// ============================================================================

// CreateTaskStatus creates a new task status
func (h *CustomizationHandler) CreateTaskStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	var status services.TenantTaskStatus
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	created, err := h.customizationService.CreateTaskStatus(r.Context(), tenantID, &status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create status: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// GetTaskStatus retrieves a specific task status
func (h *CustomizationHandler) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	vars := mux.Vars(r)
	statusCode := vars["statusCode"]

	status, err := h.customizationService.GetTaskStatus(r.Context(), tenantID, statusCode)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("status not found: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, status)
}

// GetTaskStatuses retrieves all task statuses
func (h *CustomizationHandler) GetTaskStatuses(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	statuses, err := h.customizationService.GetTaskStatuses(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get statuses: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, statuses)
}

// UpdateTaskStatus updates a task status
func (h *CustomizationHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	vars := mux.Vars(r)
	statusCode := vars["statusCode"]

	var status services.TenantTaskStatus
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	// Verify the status code matches
	if status.StatusCode != statusCode {
		writeError(w, http.StatusBadRequest, "status code mismatch")
		return
	}

	updated, err := h.customizationService.UpdateTaskStatus(r.Context(), tenantID, &status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to update status: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// DeactivateTaskStatus deactivates a task status
func (h *CustomizationHandler) DeactivateTaskStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	vars := mux.Vars(r)
	statusCode := vars["statusCode"]

	err := h.customizationService.DeactivateTaskStatus(r.Context(), tenantID, statusCode)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to deactivate status: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "status deactivated successfully"})
}

// ============================================================================
// TASK STAGE HANDLERS
// ============================================================================

// CreateTaskStage creates a new task stage
func (h *CustomizationHandler) CreateTaskStage(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	var stage services.TenantTaskStage
	if err := json.NewDecoder(r.Body).Decode(&stage); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	created, err := h.customizationService.CreateTaskStage(r.Context(), tenantID, &stage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create stage: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// GetTaskStages retrieves all task stages
func (h *CustomizationHandler) GetTaskStages(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	stages, err := h.customizationService.GetTaskStages(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get stages: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, stages)
}

// UpdateTaskStage updates a task stage
func (h *CustomizationHandler) UpdateTaskStage(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	vars := mux.Vars(r)
	stageCode := vars["stageCode"]

	var stage services.TenantTaskStage
	if err := json.NewDecoder(r.Body).Decode(&stage); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	if stage.StageCode != stageCode {
		writeError(w, http.StatusBadRequest, "stage code mismatch")
		return
	}

	updated, err := h.customizationService.UpdateTaskStage(r.Context(), tenantID, &stage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to update stage: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// ============================================================================
// STATUS TRANSITION HANDLERS
// ============================================================================

// CreateStatusTransition creates a status transition rule
func (h *CustomizationHandler) CreateStatusTransition(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	var transition services.TenantStatusTransition
	if err := json.NewDecoder(r.Body).Decode(&transition); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	created, err := h.customizationService.CreateStatusTransition(r.Context(), tenantID, &transition)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create transition: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// GetStatusTransitions retrieves allowed transitions
func (h *CustomizationHandler) GetStatusTransitions(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	fromStatus := r.URL.Query().Get("from_status")
	if fromStatus == "" {
		writeError(w, http.StatusBadRequest, "from_status query parameter required")
		return
	}

	transitions, err := h.customizationService.GetAllowedTransitions(r.Context(), tenantID, fromStatus)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get transitions: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"from_status": fromStatus,
		"allowed_to":  transitions,
	})
}

// CheckStatusTransition checks if a transition is allowed
func (h *CustomizationHandler) CheckStatusTransition(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	fromStatus := r.URL.Query().Get("from_status")
	toStatus := r.URL.Query().Get("to_status")

	if fromStatus == "" || toStatus == "" {
		writeError(w, http.StatusBadRequest, "from_status and to_status query parameters required")
		return
	}

	allowed, err := h.customizationService.IsTransitionAllowed(r.Context(), tenantID, fromStatus, toStatus)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to check transition: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"from_status": fromStatus,
		"to_status":   toStatus,
		"is_allowed":  allowed,
	})
}

// ============================================================================
// TASK TYPE HANDLERS
// ============================================================================

// CreateTaskType creates a new task type
func (h *CustomizationHandler) CreateTaskType(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	var taskType services.TenantTaskType
	if err := json.NewDecoder(r.Body).Decode(&taskType); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	created, err := h.customizationService.CreateTaskType(r.Context(), tenantID, &taskType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create task type: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// GetTaskTypes retrieves all task types
func (h *CustomizationHandler) GetTaskTypes(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	types, err := h.customizationService.GetTaskTypes(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get task types: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, types)
}

// UpdateTaskType updates a task type
func (h *CustomizationHandler) UpdateTaskType(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	vars := mux.Vars(r)
	typeCode := vars["typeCode"]

	var taskType services.TenantTaskType
	if err := json.NewDecoder(r.Body).Decode(&taskType); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	if taskType.TypeCode != typeCode {
		writeError(w, http.StatusBadRequest, "type code mismatch")
		return
	}

	updated, err := h.customizationService.UpdateTaskType(r.Context(), tenantID, &taskType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to update task type: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// ============================================================================
// PRIORITY LEVEL HANDLERS
// ============================================================================

// CreatePriorityLevel creates a new priority level
func (h *CustomizationHandler) CreatePriorityLevel(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	var priority services.TenantPriorityLevel
	if err := json.NewDecoder(r.Body).Decode(&priority); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	created, err := h.customizationService.CreatePriorityLevel(r.Context(), tenantID, &priority)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create priority: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// GetPriorityLevels retrieves all priority levels
func (h *CustomizationHandler) GetPriorityLevels(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	levels, err := h.customizationService.GetPriorityLevels(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get priorities: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, levels)
}

// UpdatePriorityLevel updates a priority level
func (h *CustomizationHandler) UpdatePriorityLevel(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	vars := mux.Vars(r)
	priorityCode := vars["priorityCode"]

	var priority services.TenantPriorityLevel
	if err := json.NewDecoder(r.Body).Decode(&priority); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	if priority.PriorityCode != priorityCode {
		writeError(w, http.StatusBadRequest, "priority code mismatch")
		return
	}

	updated, err := h.customizationService.UpdatePriorityLevel(r.Context(), tenantID, &priority)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to update priority: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// ============================================================================
// NOTIFICATION TYPE HANDLERS
// ============================================================================

// CreateNotificationType creates a new notification type
func (h *CustomizationHandler) CreateNotificationType(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	var notifType services.TenantNotificationType
	if err := json.NewDecoder(r.Body).Decode(&notifType); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	created, err := h.customizationService.CreateNotificationType(r.Context(), tenantID, &notifType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create notification type: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// GetNotificationTypes retrieves all notification types
func (h *CustomizationHandler) GetNotificationTypes(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	types, err := h.customizationService.GetNotificationTypes(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get notification types: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, types)
}

// UpdateNotificationType updates a notification type
func (h *CustomizationHandler) UpdateNotificationType(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	vars := mux.Vars(r)
	typeCode := vars["typeCode"]

	var notifType services.TenantNotificationType
	if err := json.NewDecoder(r.Body).Decode(&notifType); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	if notifType.TypeCode != typeCode {
		writeError(w, http.StatusBadRequest, "type code mismatch")
		return
	}

	updated, err := h.customizationService.UpdateNotificationType(r.Context(), tenantID, &notifType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to update notification type: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// ============================================================================
// CUSTOM FIELD HANDLERS
// ============================================================================

// CreateCustomField creates a new custom field
func (h *CustomizationHandler) CreateCustomField(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	var field services.TenantTaskField
	if err := json.NewDecoder(r.Body).Decode(&field); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	created, err := h.customizationService.CreateCustomField(r.Context(), tenantID, &field)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create field: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// GetCustomFields retrieves all custom fields
func (h *CustomizationHandler) GetCustomFields(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	fields, err := h.customizationService.GetCustomFields(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get fields: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, fields)
}

// UpdateCustomField updates a custom field
func (h *CustomizationHandler) UpdateCustomField(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	vars := mux.Vars(r)
	fieldCode := vars["fieldCode"]

	var field services.TenantTaskField
	if err := json.NewDecoder(r.Body).Decode(&field); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	if field.FieldCode != fieldCode {
		writeError(w, http.StatusBadRequest, "field code mismatch")
		return
	}

	updated, err := h.customizationService.UpdateCustomField(r.Context(), tenantID, &field)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to update field: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// ============================================================================
// AUTOMATION RULE HANDLERS
// ============================================================================

// CreateAutomationRule creates a new automation rule
func (h *CustomizationHandler) CreateAutomationRule(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	var rule services.TenantAutomationRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	created, err := h.customizationService.CreateAutomationRule(r.Context(), tenantID, &rule)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create rule: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// GetAutomationRules retrieves all automation rules
func (h *CustomizationHandler) GetAutomationRules(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	rules, err := h.customizationService.GetAutomationRules(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get rules: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, rules)
}

// UpdateAutomationRule updates an automation rule
func (h *CustomizationHandler) UpdateAutomationRule(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	vars := mux.Vars(r)
	ruleCode := vars["ruleCode"]

	var rule services.TenantAutomationRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	if rule.RuleCode != ruleCode {
		writeError(w, http.StatusBadRequest, "rule code mismatch")
		return
	}

	updated, err := h.customizationService.UpdateAutomationRule(r.Context(), tenantID, &rule)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to update rule: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// ============================================================================
// AGGREGATE CONFIGURATION HANDLER
// ============================================================================

// GetTenantConfiguration retrieves complete tenant configuration
func (h *CustomizationHandler) GetTenantConfiguration(w http.ResponseWriter, r *http.Request) {
	tenantID := getTenantID(r)
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, "tenant_id required")
		return
	}

	config, err := h.customizationService.GetTenantConfiguration(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get configuration: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, config)
}
