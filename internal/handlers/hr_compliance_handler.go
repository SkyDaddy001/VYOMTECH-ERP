package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"multi-tenant-ai-callcenter/internal/middleware"
	"multi-tenant-ai-callcenter/internal/services"

	"github.com/gorilla/mux"
)

// ============================================================================
// HR COMPLIANCE HANDLERS
// ============================================================================

type HRComplianceHandler struct {
	Service *services.HRComplianceService
}

func NewHRComplianceHandler(service *services.HRComplianceService) *HRComplianceHandler {
	return &HRComplianceHandler{Service: service}
}

// CreateESICompliance registers an employee for ESI
func (h *HRComplianceHandler) CreateESICompliance(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		EmployeeID string `json:"employee_id" validate:"required"`
		ESINumber  string `json:"esi_number" validate:"required"`
		ESIOffice  string `json:"esi_office" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	esiCompliance, err := h.Service.CreateESICompliance(
		tenantID,
		req.EmployeeID,
		req.ESINumber,
		req.ESIOffice,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(esiCompliance)
}

// CreateEPFCompliance registers an employee for EPF
func (h *HRComplianceHandler) CreateEPFCompliance(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		EmployeeID      string `json:"employee_id" validate:"required"`
		PFAccountNumber string `json:"pf_account_number" validate:"required"`
		PFOffice        string `json:"pf_office" validate:"required"`
		PFExempt        bool   `json:"pf_exempt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	epfCompliance, err := h.Service.CreateEPFCompliance(
		tenantID,
		req.EmployeeID,
		req.PFAccountNumber,
		req.PFOffice,
		req.PFExempt,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(epfCompliance)
}

// CreateProfessionalTaxCompliance sets up Professional Tax for employee
func (h *HRComplianceHandler) CreateProfessionalTaxCompliance(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		EmployeeID             string  `json:"employee_id" validate:"required"`
		PTState                string  `json:"pt_state" validate:"required"`
		MonthlySalaryThreshold float64 `json:"monthly_salary_threshold" validate:"required"`
		PTAmount               float64 `json:"pt_amount" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ptCompliance, err := h.Service.CreateProfessionalTaxCompliance(
		tenantID,
		req.EmployeeID,
		req.PTState,
		req.MonthlySalaryThreshold,
		req.PTAmount,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ptCompliance)
}

// CreateGratuityCompliance initializes gratuity tracking
func (h *HRComplianceHandler) CreateGratuityCompliance(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		EmployeeID  string    `json:"employee_id" validate:"required"`
		JoiningDate time.Time `json:"joining_date" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gratuityCompliance, err := h.Service.CreateGratuityCompliance(
		tenantID,
		req.EmployeeID,
		req.JoiningDate,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(gratuityCompliance)
}

// CheckGratuityEligibility checks and updates gratuity eligibility
func (h *HRComplianceHandler) CheckGratuityEligibility(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		EmployeeID      string  `json:"employee_id" validate:"required"`
		LastDrawnSalary float64 `json:"last_drawn_salary" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Service.CheckAndUpdateGratuityEligibility(
		tenantID,
		req.EmployeeID,
		req.LastDrawnSalary,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Gratuity eligibility updated",
		"employee_id": req.EmployeeID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RecordBonusPayment records bonus payment
func (h *HRComplianceHandler) RecordBonusPayment(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		EmployeeID  string  `json:"employee_id" validate:"required"`
		BonusType   string  `json:"bonus_type" validate:"required"`
		BonusAmount float64 `json:"bonus_amount" validate:"required,gt=0"`
		BonusYear   int     `json:"bonus_year" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bonus, err := h.Service.RecordBonusPayment(
		tenantID,
		req.EmployeeID,
		req.BonusType,
		req.BonusAmount,
		req.BonusYear,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bonus)
}

// InitializeLeaveCompliance sets up leave entitlements
func (h *HRComplianceHandler) InitializeLeaveCompliance(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		EmployeeID        string    `json:"employee_id" validate:"required"`
		LeaveTypeID       string    `json:"leave_type_id" validate:"required"`
		FiscalYearStart   time.Time `json:"fiscal_year_start" validate:"required"`
		FiscalYearEnd     time.Time `json:"fiscal_year_end" validate:"required"`
		AnnualEntitlement int       `json:"annual_entitlement" validate:"required,gt=0"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	leaveCompliance, err := h.Service.InitializeLeaveCompliance(
		tenantID,
		req.EmployeeID,
		req.LeaveTypeID,
		req.FiscalYearStart,
		req.FiscalYearEnd,
		req.AnnualEntitlement,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(leaveCompliance)
}

// LogComplianceAudit logs a compliance audit entry
func (h *HRComplianceHandler) LogComplianceAudit(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		EmployeeID     string `json:"employee_id" validate:"required"`
		ComplianceType string `json:"compliance_type" validate:"required"`
		ComplianceItem string `json:"compliance_item"`
		IsCompliant    bool   `json:"is_compliant"`
		ViolationFound string `json:"violation_found"`
		Severity       string `json:"severity" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	audit, err := h.Service.LogComplianceAudit(
		tenantID,
		req.EmployeeID,
		req.ComplianceType,
		req.ComplianceItem,
		req.IsCompliant,
		req.ViolationFound,
		req.Severity,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(audit)
}

// GetEmployeeComplianceStatus returns employee compliance status
func (h *HRComplianceHandler) GetEmployeeComplianceStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	params := mux.Vars(r)
	employeeID := params["employee_id"]

	status, err := h.Service.GetEmployeeComplianceStatus(tenantID, employeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// RegisterHRComplianceRoutes registers all HR compliance routes
func RegisterHRComplianceRoutes(router *mux.Router, handler *HRComplianceHandler) {
	hr := router.PathPrefix("/api/v1/hr-compliance").Subrouter()

	// ESI Compliance
	hr.HandleFunc("/esi", handler.CreateESICompliance).Methods("POST")

	// EPF Compliance
	hr.HandleFunc("/epf", handler.CreateEPFCompliance).Methods("POST")

	// Professional Tax
	hr.HandleFunc("/professional-tax", handler.CreateProfessionalTaxCompliance).Methods("POST")

	// Gratuity
	hr.HandleFunc("/gratuity", handler.CreateGratuityCompliance).Methods("POST")
	hr.HandleFunc("/gratuity/check-eligibility", handler.CheckGratuityEligibility).Methods("POST")

	// Bonus
	hr.HandleFunc("/bonus", handler.RecordBonusPayment).Methods("POST")

	// Leave
	hr.HandleFunc("/leave", handler.InitializeLeaveCompliance).Methods("POST")

	// Compliance Audit
	hr.HandleFunc("/audit", handler.LogComplianceAudit).Methods("POST")

	// Compliance Status
	hr.HandleFunc("/employee/{employee_id}/status", handler.GetEmployeeComplianceStatus).Methods("GET")
}
