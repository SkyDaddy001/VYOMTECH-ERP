package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

// ============================================================================
// HR DASHBOARD HANDLER
// ============================================================================

type HRDashboardHandler struct {
	HRService         *services.HRService
	ComplianceService *services.HRComplianceService
}

func NewHRDashboardHandler(hrService *services.HRService, complianceService *services.HRComplianceService) *HRDashboardHandler {
	return &HRDashboardHandler{
		HRService:         hrService,
		ComplianceService: complianceService,
	}
}

// GetHROverview returns HR department overview
func (h *HRDashboardHandler) GetHROverview(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	response := map[string]interface{}{
		"workforce": map[string]int{
			"total_employees":  0,
			"active_employees": 0,
			"on_leave":         0,
			"inactive":         0,
			"contractors":      0,
		},
		"departments":     map[string]interface{}{},
		"positions":       map[string]int{},
		"headcount_trend": []map[string]interface{}{},
		"attrition": map[string]float64{
			"current_month":    0,
			"ytd_attrition":    0,
			"voluntary_exit":   0,
			"involuntary_exit": 0,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPayrollSummary returns payroll summary for a period
func (h *HRDashboardHandler) GetPayrollSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		PayrollMonth time.Time `json:"payroll_month"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service to get aggregated payroll data
	payrollData, err := h.HRService.GetPayrollSummary(tenantID, req.PayrollMonth)
	if err != nil {
		http.Error(w, "Failed to fetch payroll data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"payroll_month": req.PayrollMonth.Format("2006-01"),
		"data":          payrollData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAttendanceDashboard returns attendance metrics
func (h *HRDashboardHandler) GetAttendanceDashboard(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service to get aggregated attendance metrics
	attendanceData, err := h.HRService.GetAttendanceMetrics(tenantID, req.StartDate, req.EndDate)
	if err != nil {
		http.Error(w, "Failed to fetch attendance data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"period": map[string]string{
			"start": req.StartDate.Format("2006-01-02"),
			"end":   req.EndDate.Format("2006-01-02"),
		},
		"data": attendanceData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetLeaveDashboard returns leave analytics
func (h *HRDashboardHandler) GetLeaveDashboard(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	// Call the service to get leave analytics
	leaveData, err := h.HRService.GetLeaveAnalytics(tenantID)
	if err != nil {
		http.Error(w, "Failed to fetch leave data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"fiscal_year": time.Now().Year(),
		"data":        leaveData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetComplianceDashboard returns HR compliance status
func (h *HRDashboardHandler) GetComplianceDashboard(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	_ = tenantID

	response := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"compliance_status": map[string]interface{}{
			"esi_compliant":      true,
			"epf_compliant":      true,
			"pt_compliant":       true,
			"gratuity_compliant": true,
			"overall_status":     "Compliant",
		},
		"violations": map[string]int{
			"critical": 0,
			"high":     0,
			"medium":   0,
			"low":      0,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterHRDashboardRoutes registers HR dashboard routes
func RegisterHRDashboardRoutes(router *mux.Router, handler *HRDashboardHandler) {
	dashboard := router.PathPrefix("/api/v1/dashboard/hr").Subrouter()
	dashboard.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Can add role-based access control here if needed
			next.ServeHTTP(w, r)
		})
	})

	dashboard.HandleFunc("/overview", handler.GetHROverview).Methods("GET")
	dashboard.HandleFunc("/payroll", handler.GetPayrollSummary).Methods("POST")
	dashboard.HandleFunc("/attendance", handler.GetAttendanceDashboard).Methods("POST")
	dashboard.HandleFunc("/leaves", handler.GetLeaveDashboard).Methods("GET")
	dashboard.HandleFunc("/compliance", handler.GetComplianceDashboard).Methods("GET")
}
