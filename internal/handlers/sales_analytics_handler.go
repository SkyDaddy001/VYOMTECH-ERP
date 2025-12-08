package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/services"
)

// SalesAnalyticsHandler handles sales analytics and reporting endpoints
type SalesAnalyticsHandler struct {
	analyticsService *services.SalesAnalyticsService
	logger           *log.Logger
}

// NewSalesAnalyticsHandler creates a new sales analytics handler
func NewSalesAnalyticsHandler(analyticsService *services.SalesAnalyticsService, logger *log.Logger) *SalesAnalyticsHandler {
	return &SalesAnalyticsHandler{
		analyticsService: analyticsService,
		logger:           logger,
	}
}

// GetMonthlySales retrieves monthly sales summary
// GET /api/v1/sales/analytics/monthly-summary
func (h *SalesAnalyticsHandler) GetMonthlySales(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	// Optional date filtering
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}

	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		}
	}

	results, err := h.analyticsService.GetMonthlySalesAnalysis(r.Context(), tenantID, startDate, endDate)
	if err != nil {
		h.logger.Printf("Error fetching monthly sales: %v", err)
		http.Error(w, "Failed to fetch sales data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    results,
	})
}

// GetCollectionReport retrieves quarterly collection breakdown
// GET /api/v1/sales/analytics/collection-report
func (h *SalesAnalyticsHandler) GetCollectionReport(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	quarterStr := r.URL.Query().Get("quarter")
	if quarterStr == "" {
		quarterStr = "1"
	}

	quarter, err := strconv.Atoi(quarterStr)
	if err != nil || quarter < 1 || quarter > 4 {
		http.Error(w, "Invalid quarter. Must be 1-4", http.StatusBadRequest)
		return
	}

	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		yearStr = strconv.Itoa(time.Now().Year())
	}

	results, err := h.analyticsService.GetCollectionReport(r.Context(), tenantID, quarter, yearStr)
	if err != nil {
		h.logger.Printf("Error fetching collection report: %v", err)
		http.Error(w, "Failed to fetch collection data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    results,
	})
}

// GetBankPaymentAnalysis retrieves bank vs own payment analysis
// GET /api/v1/sales/analytics/bank-payment-analysis
func (h *SalesAnalyticsHandler) GetBankPaymentAnalysis(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	results, err := h.analyticsService.GetBankOwnPaymentAnalysis(r.Context(), tenantID)
	if err != nil {
		h.logger.Printf("Error fetching bank payment analysis: %v", err)
		http.Error(w, "Failed to fetch bank payment data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    results,
	})
}

// GetDashboardSummary retrieves overall sales dashboard summary
// GET /api/v1/sales/analytics/dashboard-summary
func (h *SalesAnalyticsHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	result, err := h.analyticsService.GetDashboardSummary(r.Context(), tenantID)
	if err != nil {
		h.logger.Printf("Error fetching dashboard summary: %v", err)
		http.Error(w, "Failed to fetch dashboard data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// GetAgreementStatus retrieves agreement signed/pending status
// GET /api/v1/sales/analytics/agreement-status
func (h *SalesAnalyticsHandler) GetAgreementStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	// Query v_agreement_status view
	query := `SELECT * FROM v_agreement_status WHERE tenant_id = ? ORDER BY booking_date DESC`

	rows, err := h.analyticsService.GetDB().QueryContext(r.Context(), query, tenantID)
	if err != nil {
		h.logger.Printf("Error fetching agreement status: %v", err)
		http.Error(w, "Failed to fetch agreement status", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Agreement status retrieved successfully",
	})
}

// GetSoldUnitsTracking retrieves individual sold units tracking
// GET /api/v1/sales/analytics/sold-units
func (h *SalesAnalyticsHandler) GetSoldUnitsTracking(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		http.Error(w, "Tenant ID not found", http.StatusUnauthorized)
		return
	}

	// Query v_sold_units_tracking view
	query := `SELECT * FROM v_sold_units_tracking WHERE tenant_id = ? ORDER BY sold_date DESC`

	rows, err := h.analyticsService.GetDB().QueryContext(r.Context(), query, tenantID)
	if err != nil {
		h.logger.Printf("Error fetching sold units tracking: %v", err)
		http.Error(w, "Failed to fetch sold units data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Sold units tracking retrieved successfully",
	})
}
