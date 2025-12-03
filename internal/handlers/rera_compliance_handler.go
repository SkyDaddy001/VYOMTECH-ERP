package handlers

import (
	"encoding/json"
	"net/http"

	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

// ============================================================================
// RERA COMPLIANCE HANDLERS
// ============================================================================

type RERAComplianceHandler struct {
	Service *services.RERAComplianceService
}

func NewRERAComplianceHandler(service *services.RERAComplianceService) *RERAComplianceHandler {
	return &RERAComplianceHandler{Service: service}
}

// CreateProjectCollectionAccount creates a new collection account for a project
func (h *RERAComplianceHandler) CreateProjectCollectionAccount(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.CreateProjectCollectionAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	account, err := h.Service.CreateProjectCollectionAccount(
		tenantID,
		req.ProjectID,
		req.AccountName,
		req.BankName,
		req.AccountNumber,
		req.IFSCCode,
		req.MinimumBalance,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// RecordCollection records a customer collection
func (h *RERAComplianceHandler) RecordCollection(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		ProjectID           string  `json:"project_id" validate:"required"`
		CollectionAccountID string  `json:"collection_account_id" validate:"required"`
		BookingID           string  `json:"booking_id" validate:"required"`
		UnitID              string  `json:"unit_id" validate:"required"`
		PaymentMode         string  `json:"payment_mode" validate:"required"`
		AmountCollected     float64 `json:"amount_collected" validate:"required,gt=0"`
		PaidBy              string  `json:"paid_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	collection, err := h.Service.RecordCollection(
		tenantID,
		req.ProjectID,
		req.CollectionAccountID,
		req.BookingID,
		req.UnitID,
		req.PaymentMode,
		req.AmountCollected,
		req.PaidBy,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(collection)
}

// RecordFundUtilization records fund utilization
func (h *RERAComplianceHandler) RecordFundUtilization(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		ProjectID           string  `json:"project_id" validate:"required"`
		CollectionAccountID string  `json:"collection_account_id" validate:"required"`
		UtilizationType     string  `json:"utilization_type" validate:"required"`
		Description         string  `json:"description"`
		AmountUtilized      float64 `json:"amount_utilized" validate:"required,gt=0"`
		BillNumber          string  `json:"bill_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	utilization, err := h.Service.RecordFundUtilization(
		tenantID,
		req.ProjectID,
		req.CollectionAccountID,
		req.UtilizationType,
		req.Description,
		req.AmountUtilized,
		req.BillNumber,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(utilization)
}

// CheckBorrowingLimit validates borrowing eligibility
func (h *RERAComplianceHandler) CheckBorrowingLimit(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		ProjectID            string  `json:"project_id" validate:"required"`
		ProposedBorrowAmount float64 `json:"proposed_borrow_amount" validate:"required,gt=0"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	canBorrow, maxAllowed, err := h.Service.CheckBorrowingLimit(
		tenantID,
		req.ProjectID,
		req.ProposedBorrowAmount,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"can_borrow":       canBorrow,
		"max_allowed":      maxAllowed,
		"requested_amount": req.ProposedBorrowAmount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProjectCollectionSummary returns collection summary for a project
func (h *RERAComplianceHandler) GetProjectCollectionSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	params := mux.Vars(r)
	projectID := params["project_id"]

	summary, err := h.Service.GetProjectCollectionSummary(tenantID, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// PerformMonthlyReconciliation performs month-end reconciliation
func (h *RERAComplianceHandler) PerformMonthlyReconciliation(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req struct {
		ProjectID           string `json:"project_id" validate:"required"`
		CollectionAccountID string `json:"collection_account_id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reconciliation, err := h.Service.PerformMonthlyReconciliation(
		tenantID,
		req.ProjectID,
		req.CollectionAccountID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reconciliation)
}

// RegisterRERARoutes registers all RERA compliance routes
func RegisterRERARoutes(router *mux.Router, handler *RERAComplianceHandler) {
	rera := router.PathPrefix("/api/v1/rera-compliance").Subrouter()

	// Collection Account Management
	rera.HandleFunc("/project-collection-account", handler.CreateProjectCollectionAccount).Methods("POST")
	rera.HandleFunc("/project-collection-account/{project_id}", handler.GetProjectCollectionSummary).Methods("GET")

	// Collection Recording
	rera.HandleFunc("/collection", handler.RecordCollection).Methods("POST")

	// Fund Utilization
	rera.HandleFunc("/fund-utilization", handler.RecordFundUtilization).Methods("POST")

	// Borrowing Validation
	rera.HandleFunc("/check-borrowing-limit", handler.CheckBorrowingLimit).Methods("POST")

	// Reconciliation
	rera.HandleFunc("/monthly-reconciliation", handler.PerformMonthlyReconciliation).Methods("POST")
}
