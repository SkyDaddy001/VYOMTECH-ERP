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
// SALES DASHBOARD HANDLER
// ============================================================================

type SalesDashboardHandler struct {
	SalesService *services.SalesService
}

func NewSalesDashboardHandler(salesService *services.SalesService) *SalesDashboardHandler {
	return &SalesDashboardHandler{
		SalesService: salesService,
	}
}

// GetSalesOverview returns sales department overview
func (h *SalesDashboardHandler) GetSalesOverview(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	// Call service to get overview metrics
	overviewData, err := h.SalesService.GetSalesOverviewMetrics(tenantID)
	if err != nil {
		http.Error(w, "Failed to fetch sales overview: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"data":      overviewData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPipelineAnalysis returns sales pipeline analysis
func (h *SalesDashboardHandler) GetPipelineAnalysis(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	// Call service to get pipeline analysis
	pipelineData, err := h.SalesService.GetPipelineAnalysisMetrics(tenantID)
	if err != nil {
		http.Error(w, "Failed to fetch pipeline data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"data":      pipelineData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSalesMetrics returns detailed sales metrics
func (h *SalesDashboardHandler) GetSalesMetrics(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call service to get period metrics
	metricsData, err := h.SalesService.GetSalesMetricsForPeriod(tenantID, req.StartDate, req.EndDate)
	if err != nil {
		http.Error(w, "Failed to fetch sales metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"period": map[string]string{
			"start": req.StartDate.Format("2006-01-02"),
			"end":   req.EndDate.Format("2006-01-02"),
		},
		"data": metricsData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetForecast returns sales forecast
func (h *SalesDashboardHandler) GetForecast(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	quarterStart := time.Now()
	// Align to quarter start
	quarter := (int(quarterStart.Month()) - 1) / 3
	quarterStart = time.Date(quarterStart.Year(), time.Month(quarter*3+1), 1, 0, 0, 0, 0, time.UTC)
	quarterEnd := quarterStart.AddDate(0, 3, -1)

	response := map[string]interface{}{
		"forecast_period": map[string]string{
			"start": quarterStart.Format("2006-01-02"),
			"end":   quarterEnd.Format("2006-01-02"),
		},
		"overall_forecast": map[string]interface{}{
			"target":           0,
			"forecast":         0,
			"confidence_level": 0,
		},
		"by_sales_rep": []map[string]interface{}{},
		"by_product":   []map[string]interface{}{},
		"risk_factors": []map[string]interface{}{
			{
				"factor": "Delayed decision making",
				"impact": "High",
			},
		},
		"opportunities_at_risk": []map[string]interface{}{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetInvoiceStatus returns invoice tracking dashboard
func (h *SalesDashboardHandler) GetInvoiceStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	// Call service to get invoice status metrics
	invoiceData, err := h.SalesService.GetInvoiceStatusMetrics(tenantID)
	if err != nil {
		http.Error(w, "Failed to fetch invoice data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"data":      invoiceData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCompetitionAnalysis returns competitive intelligence
func (h *SalesDashboardHandler) GetCompetitionAnalysis(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	response := map[string]interface{}{
		"market_position": map[string]interface{}{
			"market_share":       0,
			"win_rate":           0,
			"loss_rate":          0,
			"competitive_wins":   0,
			"competitive_losses": 0,
		},
		"competitive_intelligence": []map[string]interface{}{
			{
				"competitor":      "Competitor A",
				"market_activity": "High",
				"customer_impact": "Medium",
			},
		},
		"product_comparison": map[string]interface{}{},
		"pricing_analysis": map[string]interface{}{
			"our_average_price": 0,
			"market_average":    0,
			"price_difference":  0.0,
		},
		"win_loss_analysis": []map[string]interface{}{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterSalesDashboardRoutes registers sales dashboard routes
func RegisterSalesDashboardRoutes(router *mux.Router, handler *SalesDashboardHandler) {
	dashboard := router.PathPrefix("/api/v1/dashboard/sales").Subrouter()
	dashboard.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Can add role-based access control here
			next.ServeHTTP(w, r)
		})
	})

	dashboard.HandleFunc("/overview", handler.GetSalesOverview).Methods("GET")
	dashboard.HandleFunc("/pipeline", handler.GetPipelineAnalysis).Methods("GET")
	dashboard.HandleFunc("/metrics", handler.GetSalesMetrics).Methods("POST")
	dashboard.HandleFunc("/forecast", handler.GetForecast).Methods("GET")
	dashboard.HandleFunc("/invoices", handler.GetInvoiceStatus).Methods("GET")
	dashboard.HandleFunc("/competition", handler.GetCompetitionAnalysis).Methods("GET")
}
