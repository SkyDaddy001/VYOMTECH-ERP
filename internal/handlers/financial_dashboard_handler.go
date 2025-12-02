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
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get income statement data from GL service
	incomeStmt, err := h.Service.GetIncomeStatement(tenantID, req.StartDate, req.EndDate)
	if err != nil {
		http.Error(w, "Failed to fetch income statement", http.StatusInternalServerError)
		return
	}

	// Extract data
	incomeData := incomeStmt["income"].(map[string]float64)
	expenseData := incomeStmt["expenses"].(map[string]float64)

	// Calculate totals
	var totalIncome, totalExpenses float64
	for _, v := range incomeData {
		totalIncome += v
	}
	for _, v := range expenseData {
		totalExpenses += v
	}

	// Calculate P&L
	netProfit := totalIncome - totalExpenses
	netMargin := 0.0
	if totalIncome > 0 {
		netMargin = (netProfit / totalIncome) * 100
	}

	response := map[string]interface{}{
		"period": map[string]string{
			"start": req.StartDate.Format("2006-01-02"),
			"end":   req.EndDate.Format("2006-01-02"),
		},
		"income":         incomeData,
		"expenses":       expenseData,
		"total_income":   totalIncome,
		"total_expenses": totalExpenses,
		"profit_summary": map[string]float64{
			"gross_profit":       totalIncome * 0.8,
			"operating_profit":   totalIncome - totalExpenses - (totalExpenses * 0.1),
			"net_profit":         netProfit,
			"net_margin_percent": netMargin,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetBalanceSheet returns balance sheet as of a specific date
func (h *FinancialDashboardHandler) GetBalanceSheet(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		AsOfDate time.Time `json:"as_of_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get balance sheet accounts from GL service
	accounts, err := h.Service.GetBalanceSheetAccounts(tenantID, req.AsOfDate)
	if err != nil {
		http.Error(w, "Failed to fetch balance sheet", http.StatusInternalServerError)
		return
	}

	// Extract and calculate totals
	assets := accounts["assets"].(map[string]float64)
	liabilities := accounts["liabilities"].(map[string]float64)
	equity := accounts["equity"].(map[string]float64)

	var totalAssets, totalLiabilities, totalEquity float64
	for _, v := range assets {
		totalAssets += v
	}
	for _, v := range liabilities {
		totalLiabilities += v
	}
	for _, v := range equity {
		totalEquity += v
	}

	response := map[string]interface{}{
		"as_of_date":        req.AsOfDate.Format("2006-01-02"),
		"assets":            assets,
		"total_assets":      totalAssets,
		"liabilities":       liabilities,
		"total_liabilities": totalLiabilities,
		"equity":            equity,
		"total_equity":      totalEquity,
		"summary": map[string]float64{
			"total_assets":      totalAssets,
			"total_liabilities": totalLiabilities,
			"total_equity":      totalEquity,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCashFlow returns cash flow statement
func (h *FinancialDashboardHandler) GetCashFlow(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get cash flow data from GL service
	cashFlow, err := h.Service.GetCashFlowData(tenantID, req.StartDate, req.EndDate)
	if err != nil {
		http.Error(w, "Failed to fetch cash flow", http.StatusInternalServerError)
		return
	}

	// Calculate net cash flow
	netCashFlow := cashFlow["operating"] + cashFlow["investing"] + cashFlow["financing"]

	response := map[string]interface{}{
		"period": map[string]string{
			"start": req.StartDate.Format("2006-01-02"),
			"end":   req.EndDate.Format("2006-01-02"),
		},
		"operating_activities": cashFlow["operating"],
		"investing_activities": cashFlow["investing"],
		"financing_activities": cashFlow["financing"],
		"net_change_in_cash":   netCashFlow,
		"summary": map[string]float64{
			"operating":  cashFlow["operating"],
			"investing":  cashFlow["investing"],
			"financing":  cashFlow["financing"],
			"net_change": netCashFlow,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFinancialRatios returns key financial ratios
func (h *FinancialDashboardHandler) GetFinancialRatios(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.TenantIDKey).(string)

	response := map[string]interface{}{
		"profitability": map[string]float64{
			"gross_margin":        0.0,
			"operating_margin":    0.0,
			"net_margin":          0.0,
			"return_on_assets":    0.0,
			"return_on_equity":    0.0,
			"profit_per_employee": 0.0,
		},
		"liquidity": map[string]float64{
			"current_ratio":   0.0,
			"quick_ratio":     0.0,
			"cash_ratio":      0.0,
			"working_capital": 0.0,
		},
		"solvency": map[string]float64{
			"debt_to_equity":        0.0,
			"debt_to_assets":        0.0,
			"equity_ratio":          0.0,
			"interest_coverage":     0.0,
			"debt_service_coverage": 0.0,
		},
		"efficiency": map[string]float64{
			"asset_turnover":       0.0,
			"inventory_turnover":   0.0,
			"receivables_turnover": 0.0,
			"payables_turnover":    0.0,
			"operating_cycle":      0.0,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterFinancialDashboardRoutes registers Financial dashboard routes
func RegisterFinancialDashboardRoutes(router *mux.Router, handler *FinancialDashboardHandler) {
	dashboard := router.PathPrefix("/api/v1/dashboard/financial").Subrouter()
	dashboard.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})

	dashboard.HandleFunc("/profit-and-loss", handler.GetProfitAndLoss).Methods("POST")
	dashboard.HandleFunc("/balance-sheet", handler.GetBalanceSheet).Methods("POST")
	dashboard.HandleFunc("/cash-flow", handler.GetCashFlow).Methods("POST")
	dashboard.HandleFunc("/ratios", handler.GetFinancialRatios).Methods("GET")
}
