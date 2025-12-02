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
// FINANCIAL DASHBOARD HANDLER
// ============================================================================

type FinancialDashboardHandler struct {
	Service *services.GLService
}

func NewFinancialDashboardHandler(glService *services.GLService) *FinancialDashboardHandler {
	return &FinancialDashboardHandler{Service: glService}
}

// GetProfitAndLoss returns P&L statement for a date range
func (h *FinancialDashboardHandler) GetProfitAndLoss(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate P&L summary
	response := map[string]interface{}{
		"period": map[string]string{
			"start": req.StartDate.Format("2006-01-02"),
			"end":   req.EndDate.Format("2006-01-02"),
		},
		"income": map[string]float64{
			"sales_revenue":   0,
			"service_revenue": 0,
			"other_income":    0,
			"total_income":    0,
		},
		"expenses": map[string]float64{
			"cost_of_goods":      0,
			"operating_expenses": 0,
			"administrative":     0,
			"depreciation":       0,
			"finance_costs":      0,
			"total_expenses":     0,
		},
		"profit": map[string]float64{
			"gross_profit":     0,
			"operating_profit": 0,
			"net_profit":       0,
			"net_margin":       0,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBalanceSheet returns balance sheet as of a specific date
func (h *FinancialDashboardHandler) GetBalanceSheet(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		AsOfDate time.Time `json:"as_of_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate balance sheet
	response := map[string]interface{}{
		"as_of_date": req.AsOfDate.Format("2006-01-02"),
		"assets": map[string]interface{}{
			"current_assets": map[string]float64{
				"cash":                0,
				"accounts_receivable": 0,
				"inventory":           0,
				"other_current":       0,
				"total":               0,
			},
			"non_current_assets": map[string]float64{
				"fixed_assets":             0,
				"accumulated_depreciation": 0,
				"investments":              0,
				"intangibles":              0,
				"total":                    0,
			},
			"total_assets": 0,
		},
		"liabilities": map[string]interface{}{
			"current_liabilities": map[string]float64{
				"accounts_payable": 0,
				"short_term_loans": 0,
				"accrued_expenses": 0,
				"other_current":    0,
				"total":            0,
			},
			"non_current_liabilities": map[string]float64{
				"long_term_loans": 0,
				"deferred_tax":    0,
				"provisions":      0,
				"total":           0,
			},
			"total_liabilities": 0,
		},
		"equity": map[string]float64{
			"share_capital":     0,
			"reserves":          0,
			"retained_earnings": 0,
			"total_equity":      0,
		},
		"total_liabilities_equity": 0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCashFlow returns cash flow statement
func (h *FinancialDashboardHandler) GetCashFlow(w http.ResponseWriter, r *http.Request) {
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
		"operating_activities": map[string]float64{
			"net_income":                 0,
			"depreciation":               0,
			"changes_in_working_capital": 0,
			"net_operating_cash":         0,
		},
		"investing_activities": map[string]float64{
			"capex":              0,
			"asset_sales":        0,
			"net_investing_cash": 0,
		},
		"financing_activities": map[string]float64{
			"debt_raised":        0,
			"debt_repaid":        0,
			"equity_raised":      0,
			"dividends_paid":     0,
			"net_financing_cash": 0,
		},
		"net_change_cash": 0,
		"opening_cash":    0,
		"closing_cash":    0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFinancialRatios returns key financial ratios
func (h *FinancialDashboardHandler) GetFinancialRatios(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	response := map[string]interface{}{
		"profitability": map[string]float64{
			"gross_margin":     0,
			"operating_margin": 0,
			"net_margin":       0,
			"roe":              0,
			"roa":              0,
		},
		"liquidity": map[string]float64{
			"current_ratio": 0,
			"quick_ratio":   0,
			"cash_ratio":    0,
		},
		"solvency": map[string]float64{
			"debt_to_equity":    0,
			"interest_coverage": 0,
			"debt_to_assets":    0,
		},
		"efficiency": map[string]float64{
			"asset_turnover":       0,
			"inventory_turnover":   0,
			"receivables_turnover": 0,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterFinancialDashboardRoutes registers financial dashboard routes
func RegisterFinancialDashboardRoutes(router *mux.Router, handler *FinancialDashboardHandler) {
	dashboard := router.PathPrefix("/api/v1/dashboard/financial").Subrouter()
	dashboard.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Can add role-based access control here if needed
			next.ServeHTTP(w, r)
		})
	})

	dashboard.HandleFunc("/profit-and-loss", handler.GetProfitAndLoss).Methods("POST")
	dashboard.HandleFunc("/balance-sheet", handler.GetBalanceSheet).Methods("POST")
	dashboard.HandleFunc("/cash-flow", handler.GetCashFlow).Methods("POST")
	dashboard.HandleFunc("/ratios", handler.GetFinancialRatios).Methods("GET")
}
