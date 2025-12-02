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
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	response := map[string]interface{}{
		"revenue_metrics": map[string]interface{}{
			"ytd_revenue":           0,
			"current_month":         0,
			"previous_month":        0,
			"month_on_month_growth": 0.0,
			"target_achievement":    0.0,
		},
		"sales_pipeline": map[string]interface{}{
			"total_opportunities": 0,
			"pipeline_value":      0,
			"expected_revenue":    0,
			"conversion_rate":     0.0,
		},
		"performance_by_rep": []map[string]interface{}{},
		"top_customers":      []map[string]interface{}{},
		"recent_activity":    []map[string]interface{}{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPipelineAnalysis returns sales pipeline analysis
func (h *SalesDashboardHandler) GetPipelineAnalysis(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	response := map[string]interface{}{
		"pipeline_summary": map[string]interface{}{
			"total_opportunities": 0,
			"total_value":         0,
			"weighted_value":      0,
			"average_deal_size":   0,
		},
		"by_stage": []map[string]interface{}{
			{
				"stage":         "Prospecting",
				"opportunities": 0,
				"value":         0,
				"count_change":  0,
				"value_change":  0,
			},
			{
				"stage":         "Qualification",
				"opportunities": 0,
				"value":         0,
				"count_change":  0,
				"value_change":  0,
			},
			{
				"stage":         "Proposal",
				"opportunities": 0,
				"value":         0,
				"count_change":  0,
				"value_change":  0,
			},
			{
				"stage":         "Negotiation",
				"opportunities": 0,
				"value":         0,
				"count_change":  0,
				"value_change":  0,
			},
			{
				"stage":         "Closed Won",
				"opportunities": 0,
				"value":         0,
				"count_change":  0,
				"value_change":  0,
			},
		},
		"by_region":  map[string]interface{}{},
		"by_product": map[string]interface{}{},
		"aged_pipeline": map[string]interface{}{
			"0_to_30_days":  0,
			"31_to_60_days": 0,
			"61_to_90_days": 0,
			"over_90_days":  0,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSalesMetrics returns detailed sales metrics
func (h *SalesDashboardHandler) GetSalesMetrics(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		Period string `json:"period"` // "monthly", "quarterly", "annual"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"period": req.Period,
		"revenue": map[string]interface{}{
			"total_revenue":       0,
			"new_revenue":         0,
			"recurring_revenue":   0,
			"average_order_value": 0,
		},
		"invoices": map[string]interface{}{
			"total_invoices": 0,
			"invoiced_value": 0,
			"paid_value":     0,
			"outstanding":    0,
		},
		"by_sales_rep":    []map[string]interface{}{},
		"by_product_line": []map[string]interface{}{},
		"by_customer_segment": map[string]interface{}{
			"enterprise": 0,
			"mid_market": 0,
			"smb":        0,
			"startup":    0,
		},
		"trend": []map[string]interface{}{},
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
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	response := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_invoices": 0,
			"total_value":    0,
			"paid":           0,
			"outstanding":    0,
			"overdue":        0,
		},
		"by_status": map[string]interface{}{
			"issued":    0,
			"delivered": 0,
			"paid":      0,
			"overdue":   0,
			"cancelled": 0,
		},
		"aging_report": map[string]interface{}{
			"0_to_30_days":  0,
			"31_to_60_days": 0,
			"61_to_90_days": 0,
			"over_90_days":  0,
		},
		"top_overdue_customers": []map[string]interface{}{},
		"collection_pipeline": []map[string]interface{}{
			{
				"period":   "This Week",
				"expected": 0,
			},
			{
				"period":   "Next 7-14 Days",
				"expected": 0,
			},
		},
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
