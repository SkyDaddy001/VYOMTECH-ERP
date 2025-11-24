package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/internal/services"
	"multi-tenant-ai-callcenter/pkg/logger"
)

type ModuleHandler struct {
	moduleService *services.ModuleService
	logger        *logger.Logger
}

func NewModuleHandler(moduleService *services.ModuleService, logger *logger.Logger) *ModuleHandler {
	return &ModuleHandler{
		moduleService: moduleService,
		logger:        logger,
	}
}

// RegisterModule creates a new module (admin only)
func (h *ModuleHandler) RegisterModule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID               string              `json:"id"`
		Name             string              `json:"name"`
		Description      string              `json:"description"`
		Category         string              `json:"category"`
		Version          string              `json:"version"`
		PricingModel     models.PricingModel `json:"pricing_model"`
		BaseCost         float64             `json:"base_cost"`
		CostPerUser      float64             `json:"cost_per_user"`
		CostPerProject   float64             `json:"cost_per_project"`
		CostPerCompany   float64             `json:"cost_per_company"`
		MaxUsers         *int                `json:"max_users"`
		MaxProjects      *int                `json:"max_projects"`
		MaxCompanies     *int                `json:"max_companies"`
		IsDependentOn    []string            `json:"is_dependent_on"`
		IsCore           bool                `json:"is_core"`
		RequiresApproval bool                `json:"requires_approval"`
		TrialDaysAllowed int                 `json:"trial_days_allowed"`
		Features         []string            `json:"features"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dependentOnJSON, _ := json.Marshal(req.IsDependentOn)
	featuresJSON, _ := json.Marshal(req.Features)

	module := &models.Module{
		ID:               req.ID,
		Name:             req.Name,
		Description:      req.Description,
		Category:         req.Category,
		Status:           models.ModuleStatusActive,
		Version:          req.Version,
		PricingModel:     req.PricingModel,
		BaseCost:         req.BaseCost,
		CostPerUser:      req.CostPerUser,
		CostPerProject:   req.CostPerProject,
		CostPerCompany:   req.CostPerCompany,
		MaxUsers:         req.MaxUsers,
		MaxProjects:      req.MaxProjects,
		MaxCompanies:     req.MaxCompanies,
		IsDependentOn:    dependentOnJSON,
		IsCore:           req.IsCore,
		RequiresApproval: req.RequiresApproval,
		TrialDaysAllowed: req.TrialDaysAllowed,
		Features:         featuresJSON,
	}

	if err := h.moduleService.RegisterModule(module); err != nil {
		http.Error(w, fmt.Sprintf("Failed to register module: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "Module registered successfully", "module_id": req.ID})
}

// ListModules lists all available modules
func (h *ModuleHandler) ListModules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	modules, err := h.moduleService.ListModules(statusPtr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list modules: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}

// SubscribeToModule subscribes a tenant to a module
func (h *ModuleHandler) SubscribeToModule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TenantID        string   `json:"tenant_id"`
		CompanyID       *string  `json:"company_id"`
		ProjectID       *string  `json:"project_id"`
		ModuleID        string   `json:"module_id"`
		MaxUsersAllowed *int     `json:"max_users_allowed"`
		MonthlyBudget   *float64 `json:"monthly_budget"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	subscription := &models.ModuleSubscription{
		ID:              fmt.Sprintf("modsub_%d", time.Now().UnixNano()),
		TenantID:        req.TenantID,
		CompanyID:       req.CompanyID,
		ProjectID:       req.ProjectID,
		ModuleID:        req.ModuleID,
		MaxUsersAllowed: req.MaxUsersAllowed,
		MonthlyBudget:   req.MonthlyBudget,
	}

	if err := h.moduleService.SubscribeToModule(subscription); err != nil {
		http.Error(w, fmt.Sprintf("Failed to subscribe: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}

// GetModuleUsage retrieves usage metrics for a module
func (h *ModuleHandler) GetModuleUsage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	subscriptionID := r.URL.Query().Get("subscription_id")
	if subscriptionID == "" {
		http.Error(w, "subscription_id is required", http.StatusBadRequest)
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	if start.IsZero() {
		start = time.Now().AddDate(0, -1, 0)
	}
	if end.IsZero() {
		end = time.Now()
	}

	usages, err := h.moduleService.GetModuleUsage(subscriptionID, start, end)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get usage: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usages)
}

// ToggleModule enables or disables a module subscription
func (h *ModuleHandler) ToggleModule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SubscriptionID string `json:"subscription_id"`
		Enabled        bool   `json:"enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.moduleService.ToggleModule(req.SubscriptionID, req.Enabled); err != nil {
		http.Error(w, fmt.Sprintf("Failed to toggle module: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Module toggled successfully"})
}

// ListSubscriptions lists all subscriptions for a tenant
func (h *ModuleHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "tenant_id is required", http.StatusBadRequest)
		return
	}

	companyID := r.URL.Query().Get("company_id")
	projectID := r.URL.Query().Get("project_id")

	var companyIDPtr, projectIDPtr *string
	if companyID != "" {
		companyIDPtr = &companyID
	}
	if projectID != "" {
		projectIDPtr = &projectID
	}

	subscriptions, err := h.moduleService.ListSubscriptions(tenantID, companyIDPtr, projectIDPtr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list subscriptions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}
