package routes

import (
	"log"
	"net/http"

	"vyomtech-backend/internal/handlers"

	"github.com/gorilla/mux"
)

// RegisterSalesRoutes registers all sales and analytics routes
func RegisterSalesRoutes(router *mux.Router, salesHandler *handlers.SalesHandler, analyticsHandler *handlers.SalesAnalyticsHandler, logger *log.Logger) {
	// Sales module routes - using direct HandleFunc (middleware applied at router level)
	salesRoutes := router.PathPrefix("/api/v1/sales").Subrouter()

	// Sales Leads endpoints
	salesRoutes.HandleFunc("/leads", salesHandler.CreateSalesLead).Methods(http.MethodPost)
	salesRoutes.HandleFunc("/leads", salesHandler.ListSalesLeads).Methods(http.MethodGet)
	salesRoutes.HandleFunc("/leads/{id}", salesHandler.GetSalesLead).Methods(http.MethodGet)
	salesRoutes.HandleFunc("/leads/{id}", salesHandler.UpdateSalesLead).Methods(http.MethodPut)
	salesRoutes.HandleFunc("/leads/{id}", salesHandler.DeleteSalesLead).Methods(http.MethodDelete)

	// Sales Customers endpoints
	salesRoutes.HandleFunc("/customers", salesHandler.CreateSalesCustomer).Methods(http.MethodPost)
	salesRoutes.HandleFunc("/customers", salesHandler.ListSalesCustomers).Methods(http.MethodGet)
	salesRoutes.HandleFunc("/customers/{id}", salesHandler.GetSalesCustomer).Methods(http.MethodGet)
	salesRoutes.HandleFunc("/customers/{id}", salesHandler.UpdateSalesCustomer).Methods(http.MethodPut)

	// Analytics endpoints
	analyticsRoutes := router.PathPrefix("/api/v1/sales/analytics").Subrouter()

	// Monthly Sales Analytics
	analyticsRoutes.HandleFunc("/monthly-summary", analyticsHandler.GetMonthlySales).Methods(http.MethodGet)

	// Collection Report
	analyticsRoutes.HandleFunc("/collection-report", analyticsHandler.GetCollectionReport).Methods(http.MethodGet)

	// Bank vs Own Payment Analysis
	analyticsRoutes.HandleFunc("/bank-payment-analysis", analyticsHandler.GetBankPaymentAnalysis).Methods(http.MethodGet)

	// Dashboard Summary
	analyticsRoutes.HandleFunc("/dashboard-summary", analyticsHandler.GetDashboardSummary).Methods(http.MethodGet)

	// Agreement Status
	analyticsRoutes.HandleFunc("/agreement-status", analyticsHandler.GetAgreementStatus).Methods(http.MethodGet)

	// Sold Units Tracking
	analyticsRoutes.HandleFunc("/sold-units", analyticsHandler.GetSoldUnitsTracking).Methods(http.MethodGet)

	logger.Println("Sales and Analytics routes registered successfully")
}
