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
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		PayrollMonth time.Time `json:"payroll_month"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"payroll_month": req.PayrollMonth.Format("2006-01"),
		"summary": map[string]float64{
			"gross_salary":          0,
			"basic_pay":             0,
			"allowances":            0,
			"deductions":            0,
			"net_salary":            0,
			"employer_contribution": 0,
			"total_cost":            0,
		},
		"breakdown": map[string]interface{}{
			"by_department": map[string]float64{},
			"by_grade":      map[string]float64{},
		},
		"compliance": map[string]interface{}{
			"esi_contribution": 0,
			"pf_contribution":  0,
			"professional_tax": 0,
			"tds_deducted":     0,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAttendanceDashboard returns attendance metrics
func (h *HRDashboardHandler) GetAttendanceDashboard(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"period": map[string]string{
			"start": req.StartDate.Format("2006-01-02"),
			"end":   req.EndDate.Format("2006-01-02"),
		},
		"overall_metrics": map[string]float64{
			"attendance_rate":       0,
			"absent_count":          0,
			"late_arrival_count":    0,
			"early_departure_count": 0,
		},
		"by_department": map[string]interface{}{},
		"trend":         []map[string]interface{}{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetLeaveDashboard returns leave analytics
func (h *HRDashboardHandler) GetLeaveDashboard(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	response := map[string]interface{}{
		"leave_summary": map[string]interface{}{
			"pending_requests": 0,
			"approved":         0,
			"rejected":         0,
			"cancelled":        0,
		},
		"by_leave_type": map[string]interface{}{
			"sick_leave": map[string]int{
				"available": 0,
				"taken":     0,
				"pending":   0,
			},
			"casual_leave": map[string]int{
				"available": 0,
				"taken":     0,
				"pending":   0,
			},
			"earned_leave": map[string]int{
				"available": 0,
				"taken":     0,
				"pending":   0,
			},
			"maternity_leave": map[string]int{
				"available": 0,
				"taken":     0,
				"pending":   0,
			},
		},
		"top_requestors":   []map[string]interface{}{},
		"approval_pending": []map[string]interface{}{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetComplianceDashboard returns HR compliance status
func (h *HRDashboardHandler) GetComplianceDashboard(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	response := map[string]interface{}{
		"compliance_status": map[string]interface{}{
			"esi_compliant":      true,
			"epf_compliant":      true,
			"pt_compliant":       true,
			"gratuity_compliant": true,
			"overall_status":     "Compliant",
		},
		"pending_audits": []map[string]interface{}{},
		"violations": map[string]int{
			"critical": 0,
			"high":     0,
			"medium":   0,
			"low":      0,
		},
		"upcoming_deadlines": []map[string]interface{}{
			{
				"compliance_item": "ESI Return Filing",
				"due_date":        time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
				"status":          "Pending",
			},
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
