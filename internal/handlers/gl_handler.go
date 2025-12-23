package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"vyomtech-backend/internal/constants"
	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type GLHandler struct {
	Service     *services.GLService
	RBACService *services.RBACService
}

// NewGLHandler creates a new GL handler
func NewGLHandler(service *services.GLService, rbacService *services.RBACService) *GLHandler {
	return &GLHandler{
		Service:     service,
		RBACService: rbacService,
	}
}

// ============================================================================
// CHART OF ACCOUNTS ENDPOINTS
// ============================================================================

// CreateAccount - POST /api/v1/gl/accounts
func (h *GLHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Extract user and tenant from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	tenant, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenant == "" {
		http.Error(w, `{"error": "Tenant ID not found in context"}`, http.StatusForbidden)
		return
	}

	// Verify permission
	if err := h.RBACService.VerifyPermission(r.Context(), tenant, userID, constants.AccountCreate); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Permission denied: %s"}`, err.Error()), http.StatusForbidden)
		return
	}

	var account models.ChartOfAccount
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	account.ID = uuid.New().String()

	if err := h.Service.CreateAccount(tenant, &account); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to create account: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// GetAccount - GET /api/v1/gl/accounts/{id}
func (h *GLHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	account, err := h.Service.GetAccount(tenant, accountID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// ListAccounts - GET /api/v1/gl/accounts
func (h *GLHandler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	tenant := r.Header.Get("X-Tenant-ID")
	accountType := r.URL.Query().Get("type")

	accounts, err := h.Service.ListAccounts(tenant, accountType)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to list accounts: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"accounts": accounts,
		"total":    len(accounts),
	})
}

// ============================================================================
// JOURNAL ENTRY ENDPOINTS
// ============================================================================

// CreateJournalEntry - POST /api/v1/gl/journal-entries
func (h *GLHandler) CreateJournalEntry(w http.ResponseWriter, r *http.Request) {
	// Extract user and tenant from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	tenant, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenant == "" {
		http.Error(w, `{"error": "Tenant ID not found in context"}`, http.StatusForbidden)
		return
	}

	// Verify permission - entries need to be posted, so check EntryPost
	if err := h.RBACService.VerifyPermission(r.Context(), tenant, userID, constants.EntryPost); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Permission denied: %s"}`, err.Error()), http.StatusForbidden)
		return
	}

	var req models.JournalEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	entryID := uuid.New().String()

	entry := &models.JournalEntry{
		ID:              entryID,
		TenantID:        tenant,
		EntryDate:       req.EntryDate,
		ReferenceNumber: &req.ReferenceNumber,
		ReferenceType:   req.ReferenceType,
		Description:     req.Description,
		Narration:       req.Narration,
	}

	if err := h.Service.CreateJournalEntry(tenant, entry); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to create entry: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Add details
	for i, detail := range req.Details {
		entryDetail := &models.JournalEntryDetail{
			ID:             uuid.New().String(),
			TenantID:       tenant,
			JournalEntryID: entryID,
			AccountID:      detail.AccountID,
			DebitAmount:    detail.DebitAmount,
			CreditAmount:   detail.CreditAmount,
			Description:    detail.Description,
			LineNumber:     i + 1,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := h.Service.AddJournalEntryDetail(entryDetail); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Failed to add entry detail: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		entry.Details = append(entry.Details, *entryDetail)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

// PostJournalEntry - POST /api/v1/gl/journal-entries/{id}/post
func (h *GLHandler) PostJournalEntry(w http.ResponseWriter, r *http.Request) {
	// Extract user and tenant from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	tenant, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenant == "" {
		http.Error(w, `{"error": "Tenant ID not found in context"}`, http.StatusForbidden)
		return
	}

	// Verify permission
	if err := h.RBACService.VerifyPermission(r.Context(), tenant, userID, constants.EntryPost); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Permission denied: %s"}`, err.Error()), http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	entryID := vars["id"]

	var req models.PostJournalEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if err := h.Service.PostJournalEntry(tenant, entryID, req.PostedBy); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to post entry: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	entry, err := h.Service.GetJournalEntry(tenant, entryID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to retrieve entry: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

// GetJournalEntry - GET /api/v1/gl/journal-entries/{id}
func (h *GLHandler) GetJournalEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entryID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	entry, err := h.Service.GetJournalEntry(tenant, entryID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

// ListJournalEntries - GET /api/v1/gl/journal-entries
func (h *GLHandler) ListJournalEntries(w http.ResponseWriter, r *http.Request) {
	tenant := r.Header.Get("X-Tenant-ID")

	fromDateStr := r.URL.Query().Get("from_date")
	toDateStr := r.URL.Query().Get("to_date")

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		fromDate = time.Now().AddDate(0, -1, 0)
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		toDate = time.Now()
	}

	entries, err := h.Service.ListJournalEntries(tenant, fromDate, toDate)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to list entries: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"entries": entries,
		"total":   len(entries),
		"from":    fromDate.Format("2006-01-02"),
		"to":      toDate.Format("2006-01-02"),
	})
}

// ============================================================================
// REPORTING ENDPOINTS
// ============================================================================

// GetTrialBalance - GET /api/v1/gl/reports/trial-balance
func (h *GLHandler) GetTrialBalance(w http.ResponseWriter, r *http.Request) {
	tenant := r.Header.Get("X-Tenant-ID")

	fromDateStr := r.URL.Query().Get("from_date")
	toDateStr := r.URL.Query().Get("to_date")

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		fromDate = time.Now().AddDate(-1, 0, 0)
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		toDate = time.Now()
	}

	trialBalance, err := h.Service.GetTrialBalance(tenant, fromDate, toDate)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to generate trial balance: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	var totalDebit, totalCredit float64
	for _, tb := range trialBalance {
		totalDebit += tb.DebitBalance
		totalCredit += tb.CreditBalance
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"trial_balance": trialBalance,
		"total_debit":   totalDebit,
		"total_credit":  totalCredit,
		"is_balanced":   totalDebit == totalCredit,
		"from_date":     fromDate.Format("2006-01-02"),
		"to_date":       toDate.Format("2006-01-02"),
	})
}

// GetAccountLedger - GET /api/v1/gl/accounts/{id}/ledger
func (h *GLHandler) GetAccountLedger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	fromDateStr := r.URL.Query().Get("from_date")
	toDateStr := r.URL.Query().Get("to_date")

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		fromDate = time.Now().AddDate(-1, 0, 0)
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		toDate = time.Now()
	}

	ledger, err := h.Service.GetAccountLedger(tenant, accountID, fromDate, toDate)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to retrieve ledger: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"account_id": accountID,
		"ledger":     ledger,
		"entries":    len(ledger),
		"from_date":  fromDate.Format("2006-01-02"),
		"to_date":    toDate.Format("2006-01-02"),
	})
}

// ============================================================================
// FINANCIAL PERIOD ENDPOINTS
// ============================================================================

// CreateFinancialPeriod - POST /api/v1/gl/periods
func (h *GLHandler) CreateFinancialPeriod(w http.ResponseWriter, r *http.Request) {
	var period models.FinancialPeriod
	if err := json.NewDecoder(r.Body).Decode(&period); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	tenant := r.Header.Get("X-Tenant-ID")
	period.ID = uuid.New().String()

	if err := h.Service.CreateFinancialPeriod(tenant, &period); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to create period: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(period)
}

// GetFinancialPeriod - GET /api/v1/gl/periods/{id}
func (h *GLHandler) GetFinancialPeriod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	periodID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	period, err := h.Service.GetFinancialPeriod(tenant, periodID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(period)
}

// ClosePeriod - POST /api/v1/gl/periods/{id}/close
func (h *GLHandler) ClosePeriod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	periodID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	var req struct {
		ClosedBy string `json:"closed_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if err := h.Service.ClosePeriod(tenant, periodID, req.ClosedBy); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to close period: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Period closed successfully"})
}

// ============================================================================
// ROUTE REGISTRATION
// ============================================================================

// RegisterGLRoutes registers all GL routes
func RegisterGLRoutes(r *mux.Router, glService *services.GLService, rbacService *services.RBACService) {
	handler := NewGLHandler(glService, rbacService)

	// Chart of Accounts routes
	r.HandleFunc("/api/v1/gl/accounts", handler.CreateAccount).Methods("POST")
	r.HandleFunc("/api/v1/gl/accounts", handler.ListAccounts).Methods("GET")
	r.HandleFunc("/api/v1/gl/accounts/{id}", handler.GetAccount).Methods("GET")

	// Journal Entry routes
	r.HandleFunc("/api/v1/gl/journal-entries", handler.CreateJournalEntry).Methods("POST")
	r.HandleFunc("/api/v1/gl/journal-entries", handler.ListJournalEntries).Methods("GET")
	r.HandleFunc("/api/v1/gl/journal-entries/{id}", handler.GetJournalEntry).Methods("GET")
	r.HandleFunc("/api/v1/gl/journal-entries/{id}/post", handler.PostJournalEntry).Methods("POST")

	// Reporting routes
	r.HandleFunc("/api/v1/gl/reports/trial-balance", handler.GetTrialBalance).Methods("GET")
	r.HandleFunc("/api/v1/gl/accounts/{id}/ledger", handler.GetAccountLedger).Methods("GET")

	// Financial Period routes
	r.HandleFunc("/api/v1/gl/periods", handler.CreateFinancialPeriod).Methods("POST")
	r.HandleFunc("/api/v1/gl/periods/{id}", handler.GetFinancialPeriod).Methods("GET")
	r.HandleFunc("/api/v1/gl/periods/{id}/close", handler.ClosePeriod).Methods("POST")
}
