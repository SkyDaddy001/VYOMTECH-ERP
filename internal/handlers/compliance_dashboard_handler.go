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
// COMPLIANCE DASHBOARD HANDLER
// ============================================================================

type ComplianceDashboardHandler struct {
	RERAService          *services.RERAComplianceService
	HRComplianceService  *services.HRComplianceService
	TaxComplianceService *services.TaxComplianceService
}

func NewComplianceDashboardHandler(
	reraService *services.RERAComplianceService,
	hrService *services.HRComplianceService,
	taxService *services.TaxComplianceService,
) *ComplianceDashboardHandler {
	return &ComplianceDashboardHandler{
		RERAService:          reraService,
		HRComplianceService:  hrService,
		TaxComplianceService: taxService,
	}
}

// GetRERAComplianceStatus returns RERA compliance overview
func (h *ComplianceDashboardHandler) GetRERAComplianceStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	// Call the service to get RERA metrics
	reraData, err := h.RERAService.GetRERAComplianceMetrics(tenantID)
	if err != nil {
		http.Error(w, "Failed to fetch RERA compliance data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"data":      reraData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetHRComplianceStatus returns HR & Labour law compliance status
func (h *ComplianceDashboardHandler) GetHRComplianceStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	// Call the service to get HR compliance metrics
	hrData, err := h.HRComplianceService.GetHRComplianceMetrics(tenantID)
	if err != nil {
		http.Error(w, "Failed to fetch HR compliance data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"data":      hrData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTaxComplianceStatus returns Tax compliance status
func (h *ComplianceDashboardHandler) GetTaxComplianceStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	// Call the service to get tax compliance metrics
	taxData, err := h.TaxComplianceService.GetTaxComplianceMetrics(tenantID)
	if err != nil {
		http.Error(w, "Failed to fetch tax compliance data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"data":      taxData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetComplianceHealthScore returns overall compliance health score
func (h *ComplianceDashboardHandler) GetComplianceHealthScore(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	response := map[string]interface{}{
		"overall_health_score": 95,
		"scores": map[string]int{
			"rera_compliance": 95,
			"hr_compliance":   92,
			"tax_compliance":  98,
			"documentation":   90,
		},
		"risk_factors": []map[string]interface{}{},
		"recommendations": []map[string]interface{}{
			{
				"priority": "High",
				"item":     "Update GST records for recent invoices",
			},
		},
		"compliance_calendar": []map[string]interface{}{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetComplianceDocumentation returns documentation tracking
func (h *ComplianceDashboardHandler) GetComplianceDocumentation(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	response := map[string]interface{}{
		"documents_by_category": map[string]interface{}{
			"rera_documents": map[string]interface{}{
				"project_agreements":     true,
				"project_accounts":       true,
				"fund_utilization_logs":  true,
				"reconciliation_reports": false,
			},
			"hr_documents": map[string]interface{}{
				"esi_registrations":  true,
				"epf_registrations":  true,
				"pt_registrations":   true,
				"service_agreements": true,
				"leave_policies":     true,
			},
			"tax_documents": map[string]interface{}{
				"pan_registration":     true,
				"gst_registration":     true,
				"tax_audit_reports":    true,
				"financial_statements": true,
			},
		},
		"missing_documents": []map[string]string{
			{
				"category":    "rera",
				"document":    "Reconciliation Reports",
				"required_by": time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
			},
		},
		"upload_status": map[string]interface{}{
			"total_required": 35,
			"uploaded":       32,
			"pending":        3,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterComplianceDashboardRoutes registers compliance dashboard routes
func RegisterComplianceDashboardRoutes(router *mux.Router, handler *ComplianceDashboardHandler) {
	dashboard := router.PathPrefix("/api/v1/dashboard/compliance").Subrouter()
	dashboard.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Role-based access control can be added here
			next.ServeHTTP(w, r)
		})
	})

	dashboard.HandleFunc("/rera-status", handler.GetRERAComplianceStatus).Methods("GET")
	dashboard.HandleFunc("/hr-status", handler.GetHRComplianceStatus).Methods("GET")
	dashboard.HandleFunc("/tax-status", handler.GetTaxComplianceStatus).Methods("GET")
	dashboard.HandleFunc("/health-score", handler.GetComplianceHealthScore).Methods("GET")
	dashboard.HandleFunc("/documentation", handler.GetComplianceDocumentation).Methods("GET")
}
