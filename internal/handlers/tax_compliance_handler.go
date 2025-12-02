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
// TAX COMPLIANCE HANDLERS
// ============================================================================

type TaxComplianceHandler struct {
	Service *services.TaxComplianceService
}

func NewTaxComplianceHandler(service *services.TaxComplianceService) *TaxComplianceHandler {
	return &TaxComplianceHandler{Service: service}
}

// SetupTaxConfiguration initializes tax configuration for organization
func (h *TaxComplianceHandler) SetupTaxConfiguration(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		ITPAN                 string    `json:"it_pan" validate:"required"`
		GSTRegistrationNumber string    `json:"gst_registration_number" validate:"required"`
		GSTRegistrationDate   time.Time `json:"gst_registration_date" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config, err := h.Service.SetupTaxConfiguration(
		tenantID,
		req.ITPAN,
		req.GSTRegistrationNumber,
		req.GSTRegistrationDate,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

// InitializeIncomeTaxCompliance creates Income Tax compliance record for fiscal year
func (h *TaxComplianceHandler) InitializeIncomeTaxCompliance(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		FiscalYear int `json:"fiscal_year" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	compliance, err := h.Service.InitializeIncomeTaxCompliance(
		tenantID,
		req.FiscalYear,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(compliance)
}

// InitializeGSTCompliance creates GST compliance record for return period
func (h *TaxComplianceHandler) InitializeGSTCompliance(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		FiscalYear   int       `json:"fiscal_year" validate:"required"`
		ReturnPeriod string    `json:"return_period" validate:"required"`
		MonthYear    time.Time `json:"month_year" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	compliance, err := h.Service.InitializeGSTCompliance(
		tenantID,
		req.FiscalYear,
		req.ReturnPeriod,
		req.MonthYear,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(compliance)
}

// TrackGSTInvoice records GST on sales invoice
func (h *TaxComplianceHandler) TrackGSTInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		InvoiceID         string    `json:"invoice_id" validate:"required"`
		InvoiceNumber     string    `json:"invoice_number" validate:"required"`
		CustomerGSTIN     string    `json:"customer_gstin"`
		InvoiceAmount     float64   `json:"invoice_amount" validate:"required,gt=0"`
		GSTRate           float64   `json:"gst_rate" validate:"required"`
		GSTAmount         float64   `json:"gst_amount" validate:"required,gt=0"`
		InvoiceRaisedDate time.Time `json:"invoice_raised_date" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tracking, err := h.Service.TrackGSTInvoice(
		tenantID,
		req.InvoiceID,
		req.InvoiceNumber,
		req.CustomerGSTIN,
		req.InvoiceAmount,
		req.GSTRate,
		req.GSTAmount,
		req.InvoiceRaisedDate,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tracking)
}

// TrackGSTInputCredit records GST on purchase invoice
func (h *TaxComplianceHandler) TrackGSTInputCredit(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		PurchaseInvoiceID string    `json:"purchase_invoice_id" validate:"required"`
		VendorGSTIN       string    `json:"vendor_gstin" validate:"required"`
		InvoiceNumber     string    `json:"invoice_number" validate:"required"`
		InvoiceAmount     float64   `json:"invoice_amount" validate:"required,gt=0"`
		GSTRate           float64   `json:"gst_rate" validate:"required"`
		GSTAmount         float64   `json:"gst_amount" validate:"required,gt=0"`
		InvoiceDate       time.Time `json:"invoice_date" validate:"required"`
		ITCEligible       bool      `json:"itc_eligible"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	inputCredit, err := h.Service.TrackGSTInputCredit(
		tenantID,
		req.PurchaseInvoiceID,
		req.VendorGSTIN,
		req.InvoiceNumber,
		req.InvoiceAmount,
		req.GSTRate,
		req.GSTAmount,
		req.InvoiceDate,
		req.ITCEligible,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inputCredit)
}

// InitializeAdvanceTaxSchedule sets up quarterly advance tax schedule
func (h *TaxComplianceHandler) InitializeAdvanceTaxSchedule(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		FiscalYear            int     `json:"fiscal_year" validate:"required"`
		EstimatedTaxLiability float64 `json:"estimated_tax_liability" validate:"required,gt=0"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	schedule, err := h.Service.InitializeAdvanceTaxSchedule(
		tenantID,
		req.FiscalYear,
		req.EstimatedTaxLiability,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schedule)
}

// RecordAdvanceTaxPayment records quarterly advance tax payment
func (h *TaxComplianceHandler) RecordAdvanceTaxPayment(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		FiscalYear    int     `json:"fiscal_year" validate:"required"`
		Quarter       int     `json:"quarter" validate:"required,min=1,max=4"`
		AmountPaid    float64 `json:"amount_paid" validate:"required,gt=0"`
		ChallanNumber string  `json:"challan_number" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Service.RecordAdvanceTaxPayment(
		tenantID,
		req.FiscalYear,
		req.Quarter,
		req.AmountPaid,
		req.ChallanNumber,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":        "Advance tax payment recorded",
		"fiscal_year":    req.FiscalYear,
		"quarter":        req.Quarter,
		"amount_paid":    req.AmountPaid,
		"challan_number": req.ChallanNumber,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTaxComplianceStatus returns comprehensive tax compliance status
func (h *TaxComplianceHandler) GetTaxComplianceStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		FiscalYear int `json:"fiscal_year" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status, err := h.Service.GetTaxComplianceStatus(tenantID, req.FiscalYear)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// RegisterTaxComplianceRoutes registers all tax compliance routes
func RegisterTaxComplianceRoutes(router *mux.Router, handler *TaxComplianceHandler) {
	tax := router.PathPrefix("/api/v1/tax-compliance").Subrouter()

	// Tax Configuration
	tax.HandleFunc("/configuration", handler.SetupTaxConfiguration).Methods("POST")

	// Income Tax
	tax.HandleFunc("/income-tax", handler.InitializeIncomeTaxCompliance).Methods("POST")

	// GST Compliance
	tax.HandleFunc("/gst", handler.InitializeGSTCompliance).Methods("POST")
	tax.HandleFunc("/gst/invoice", handler.TrackGSTInvoice).Methods("POST")
	tax.HandleFunc("/gst/input-credit", handler.TrackGSTInputCredit).Methods("POST")

	// Advance Tax
	tax.HandleFunc("/advance-tax/schedule", handler.InitializeAdvanceTaxSchedule).Methods("POST")
	tax.HandleFunc("/advance-tax/payment", handler.RecordAdvanceTaxPayment).Methods("POST")

	// Tax Status
	tax.HandleFunc("/status", handler.GetTaxComplianceStatus).Methods("POST")
}
