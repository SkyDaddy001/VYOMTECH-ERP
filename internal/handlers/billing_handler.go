package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

type BillingHandler struct {
	billingService *services.BillingService
	logger         *logger.Logger
	RBACService    *services.RBACService
}

func NewBillingHandler(billingService *services.BillingService, logger *logger.Logger, rbacService *services.RBACService) *BillingHandler {
	return &BillingHandler{
		billingService: billingService,
		logger:         logger,
		RBACService:    rbacService,
	}
}

// CreatePricingPlan creates a new pricing plan (admin only)
func (h *BillingHandler) CreatePricingPlan(w http.ResponseWriter, r *http.Request) {
	// Extract user and tenant from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	tenant, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenant == "" {
		http.Error(w, `{"error": "Tenant ID not found in context"}`, http.StatusForbidden)
		return
	}

	// Verify permission - billing requires special privileges
	if err := h.RBACService.VerifyPermission(r.Context(), tenant, userID, "billing.plans.create"); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Permission denied: %s"}`, err.Error()), http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name                  string   `json:"name"`
		Description           string   `json:"description"`
		MonthlyPrice          float64  `json:"monthly_price"`
		AnnualPrice           float64  `json:"annual_price"`
		MaxUsers              int      `json:"max_users"`
		MaxCompanies          int      `json:"max_companies"`
		MaxProjectsPerCompany int      `json:"max_projects_per_company"`
		IncludedModules       []string `json:"included_modules"`
		AdditionalModules     []string `json:"additional_modules"`
		Features              []string `json:"features"`
		SortOrder             int      `json:"sort_order"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	includedJSON, _ := json.Marshal(req.IncludedModules)
	additionalJSON, _ := json.Marshal(req.AdditionalModules)
	featuresJSON, _ := json.Marshal(req.Features)

	plan := &models.PricingPlan{
		ID:                    fmt.Sprintf("plan_%d", time.Now().UnixNano()),
		Name:                  req.Name,
		Description:           req.Description,
		MonthlyPrice:          req.MonthlyPrice,
		AnnualPrice:           req.AnnualPrice,
		MaxUsers:              req.MaxUsers,
		MaxCompanies:          req.MaxCompanies,
		MaxProjectsPerCompany: req.MaxProjectsPerCompany,
		IncludedModules:       string(includedJSON),
		AdditionalModules:     string(additionalJSON),
		Features:              string(featuresJSON),
		SortOrder:             req.SortOrder,
		IsActive:              true,
	}

	if err := h.billingService.CreatePricingPlan(plan); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create pricing plan: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(plan)
}

// ListPricingPlans lists all active pricing plans
func (h *BillingHandler) ListPricingPlans(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	plans, err := h.billingService.ListActivePricingPlans()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list pricing plans: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plans)
}

// SubscribeToPlan subscribes a tenant to a pricing plan
func (h *BillingHandler) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TenantID      string `json:"tenant_id"`
		PricingPlanID string `json:"pricing_plan_id"`
		BillingCycle  string `json:"billing_cycle"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	subscription, err := h.billingService.SubscribeToPlan(req.TenantID, req.PricingPlanID, models.BillingCycleType(req.BillingCycle))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to subscribe: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}

// GetInvoice retrieves an invoice
func (h *BillingHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	invoiceID := r.URL.Query().Get("invoice_id")
	if invoiceID == "" {
		http.Error(w, "invoice_id is required", http.StatusBadRequest)
		return
	}

	invoice, err := h.billingService.GetInvoice(invoiceID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get invoice: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

// ListInvoices lists invoices for a tenant
func (h *BillingHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "tenant_id is required", http.StatusBadRequest)
		return
	}

	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	invoices, err := h.billingService.ListInvoicesByTenant(tenantID, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list invoices: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

// RecordUsageMetrics records usage metrics for a tenant
func (h *BillingHandler) RecordUsageMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TenantID      string  `json:"tenant_id"`
		CompanyID     *string `json:"company_id"`
		ProjectID     *string `json:"project_id"`
		ActiveUsers   int     `json:"active_users"`
		NewUsers      int     `json:"new_users"`
		APICallsUsed  int     `json:"api_calls_used"`
		StorageUsedMB float64 `json:"storage_used_mb"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	metrics := &models.UsageMetrics{
		ID:            fmt.Sprintf("usage_%d", time.Now().UnixNano()),
		TenantID:      req.TenantID,
		CompanyID:     req.CompanyID,
		ProjectID:     req.ProjectID,
		Date:          time.Now(),
		ActiveUsers:   req.ActiveUsers,
		NewUsers:      req.NewUsers,
		APICallsUsed:  req.APICallsUsed,
		StorageUsedMB: req.StorageUsedMB,
	}

	if err := h.billingService.RecordUsageMetrics(metrics); err != nil {
		http.Error(w, fmt.Sprintf("Failed to record metrics: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(metrics)
}

// GetUsageMetrics retrieves usage metrics
func (h *BillingHandler) GetUsageMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "tenant_id is required", http.StatusBadRequest)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	metrics, err := h.billingService.GetUsageMetrics(tenantID, startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get metrics: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// CalculateMonthlyCharges calculates monthly charges
func (h *BillingHandler) CalculateMonthlyCharges(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "tenant_id is required", http.StatusBadRequest)
		return
	}

	totalCost, err := h.billingService.CalculateMonthlyCharges(tenantID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to calculate charges: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"total_monthly_charges": totalCost})
}

// MarkInvoiceAsPaid marks an invoice as paid
func (h *BillingHandler) MarkInvoiceAsPaid(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		InvoiceID string `json:"invoice_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.billingService.MarkInvoiceAsPaid(req.InvoiceID, time.Now()); err != nil {
		http.Error(w, fmt.Sprintf("Failed to mark invoice as paid: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Invoice marked as paid"})
}
