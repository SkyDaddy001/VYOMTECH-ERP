package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
)

// ============================================================================
// BankFinancingHandler
// ============================================================================
type BankFinancingHandler struct {
	financingService *services.BankFinancingService
}

// NewBankFinancingHandler creates new bank financing handler
func NewBankFinancingHandler(financingService *services.BankFinancingService) *BankFinancingHandler {
	return &BankFinancingHandler{
		financingService: financingService,
	}
}

// Helper functions
func (h *BankFinancingHandler) getTenantID(r *http.Request) int64 {
	if id, ok := r.Context().Value("tenant_id").(int64); ok {
		return id
	}
	return 0
}

func (h *BankFinancingHandler) getUserID(r *http.Request) int64 {
	if id, ok := r.Context().Value("user_id").(int64); ok {
		return id
	}
	return 0
}

func (h *BankFinancingHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *BankFinancingHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

// ============================================================================
// Financing Endpoints
// ============================================================================

// CreateFinancing creates new financing record
func (h *BankFinancingHandler) CreateFinancing(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	userID := h.getUserID(r)

	var req models.CreateFinancingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	financing, err := h.financingService.CreateFinancing(tenantID, &req, userID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, financing)
}

// GetFinancing retrieves financing by ID
func (h *BankFinancingHandler) GetFinancing(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	financing, err := h.financingService.GetFinancingByID(tenantID, id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, financing)
}

// GetFinancingByBooking retrieves financing by booking ID
func (h *BankFinancingHandler) GetFinancingByBooking(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	bookingIDStr := r.PathValue("bookingId")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	financing, err := h.financingService.GetFinancingByBookingID(tenantID, bookingID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if financing == nil {
		h.writeError(w, http.StatusNotFound, "No financing found for this booking")
		return
	}

	h.writeJSON(w, http.StatusOK, financing)
}

// ListFinancing lists all financing records
func (h *BankFinancingHandler) ListFinancing(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	filters := make(map[string]interface{})
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}
	if loanType := r.URL.Query().Get("loan_type"); loanType != "" {
		filters["loan_type"] = loanType
	}

	financings, total, err := h.financingService.ListFinancing(tenantID, filters, limit, offset)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":   financings,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateFinancing updates financing record
func (h *BankFinancingHandler) UpdateFinancing(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	userID := h.getUserID(r)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req models.UpdateFinancingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	financing, err := h.financingService.UpdateFinancing(tenantID, id, &req, userID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, financing)
}

// DeleteFinancing soft deletes financing record
func (h *BankFinancingHandler) DeleteFinancing(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	userID := h.getUserID(r)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.financingService.DeleteFinancing(tenantID, id, userID); err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================================
// Disbursement Endpoints
// ============================================================================

// CreateDisbursement creates disbursement schedule
func (h *BankFinancingHandler) CreateDisbursement(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	userID := h.getUserID(r)

	var req models.CreateDisbursementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	disbursement, err := h.financingService.CreateDisbursement(tenantID, &req, userID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, disbursement)
}

// ListDisbursements lists disbursements for financing
func (h *BankFinancingHandler) ListDisbursements(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	financingIDStr := r.PathValue("financingId")
	financingID, err := strconv.ParseInt(financingIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid financing ID")
		return
	}

	disbursements, err := h.financingService.ListDisbursements(tenantID, financingID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, disbursements)
}

// UpdateDisbursementStatus updates disbursement status
func (h *BankFinancingHandler) UpdateDisbursementStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	userID := h.getUserID(r)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req struct {
		Status       string   `json:"status"`
		ActualAmount *float64 `json:"actual_amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	disbursement, err := h.financingService.UpdateDisbursementStatus(tenantID, id, req.Status, req.ActualAmount, userID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, disbursement)
}

// ============================================================================
// NOC Endpoints
// ============================================================================

// CreateNOC creates NOC record
func (h *BankFinancingHandler) CreateNOC(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	userID := h.getUserID(r)

	var req models.CreateNOCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	noc, err := h.financingService.CreateNOC(tenantID, &req, userID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, noc)
}

// ListNOCs lists NOCs for financing
func (h *BankFinancingHandler) ListNOCs(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	financingIDStr := r.PathValue("financingId")
	financingID, err := strconv.ParseInt(financingIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid financing ID")
		return
	}

	nocs, err := h.financingService.ListNOCs(tenantID, financingID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, nocs)
}

// ============================================================================
// Collection Endpoints
// ============================================================================

// CreateCollection creates collection record
func (h *BankFinancingHandler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	userID := h.getUserID(r)

	var req models.CreateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	collection, err := h.financingService.CreateCollection(tenantID, &req, userID)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, collection)
}

// ListCollections lists collections for financing
func (h *BankFinancingHandler) ListCollections(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	financingIDStr := r.PathValue("financingId")
	financingID, err := strconv.ParseInt(financingIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid financing ID")
		return
	}

	collections, err := h.financingService.ListCollections(tenantID, financingID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, collections)
}

// GetFinancingSummary retrieves comprehensive financing summary
func (h *BankFinancingHandler) GetFinancingSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantID(r)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	summary, err := h.financingService.GetFinancingSummary(tenantID, id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, summary)
}
