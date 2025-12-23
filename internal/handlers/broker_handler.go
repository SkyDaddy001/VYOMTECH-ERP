package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

// ============================================================================
// BrokerHandler - API Endpoints for Broker Management
// ============================================================================
type BrokerHandler struct {
	brokerService *services.BrokerService
}

// NewBrokerHandler creates a new broker handler
func NewBrokerHandler(brokerService *services.BrokerService) *BrokerHandler {
	return &BrokerHandler{
		brokerService: brokerService,
	}
}

// ============================================================================
// Broker Profile Handlers
// ============================================================================

// CreateBroker handles POST /api/v1/brokers
func (h *BrokerHandler) CreateBroker(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	userID, _ := r.Context().Value("user_id").(int64)

	var req models.CreateBrokerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	broker, err := h.brokerService.CreateBroker(tenantID, &req, &userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create broker: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(broker)
}

// GetBroker handles GET /api/v1/brokers/{brokerId}
func (h *BrokerHandler) GetBroker(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	brokerID, err := strconv.ParseInt(vars["brokerId"], 10, 64)
	if err != nil {
		http.Error(w, "invalid broker id", http.StatusBadRequest)
		return
	}

	broker, err := h.brokerService.GetBroker(tenantID, brokerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get broker: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(broker)
}

// ListBrokers handles GET /api/v1/brokers
func (h *BrokerHandler) ListBrokers(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	status := r.URL.Query().Get("status")

	brokers, total, err := h.brokerService.ListBrokers(tenantID, &status, offset, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list brokers: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  brokers,
		"total": total,
	})
}

// UpdateBroker handles PUT /api/v1/brokers/{brokerId}
func (h *BrokerHandler) UpdateBroker(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	userID, _ := r.Context().Value("user_id").(int64)

	vars := mux.Vars(r)
	brokerID, err := strconv.ParseInt(vars["brokerId"], 10, 64)
	if err != nil {
		http.Error(w, "invalid broker id", http.StatusBadRequest)
		return
	}

	var req models.UpdateBrokerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	broker, err := h.brokerService.UpdateBroker(tenantID, brokerID, &req, &userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update broker: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(broker)
}

// DeleteBroker handles DELETE /api/v1/brokers/{brokerId}
func (h *BrokerHandler) DeleteBroker(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	userID, _ := r.Context().Value("user_id").(int64)

	vars := mux.Vars(r)
	brokerID, err := strconv.ParseInt(vars["brokerId"], 10, 64)
	if err != nil {
		http.Error(w, "invalid broker id", http.StatusBadRequest)
		return
	}

	if err := h.brokerService.DeleteBroker(tenantID, brokerID, &userID); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete broker: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================================
// Commission Structure Handlers
// ============================================================================

// CreateCommissionStructure handles POST /api/v1/brokers/{brokerId}/commission-structure
func (h *BrokerHandler) CreateCommissionStructure(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	userID, _ := r.Context().Value("user_id").(int64)

	var req models.CreateCommissionStructureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	structure, err := h.brokerService.CreateCommissionStructure(tenantID, &req, &userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create commission structure: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(structure)
}

// ListCommissionStructures handles GET /api/v1/brokers/{brokerId}/commission-structure
func (h *BrokerHandler) ListCommissionStructures(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	brokerID, err := strconv.ParseInt(vars["brokerId"], 10, 64)
	if err != nil {
		http.Error(w, "invalid broker id", http.StatusBadRequest)
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	structures, total, err := h.brokerService.ListCommissionStructures(tenantID, brokerID, offset, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list structures: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  structures,
		"total": total,
	})
}

// ============================================================================
// Booking Link Handlers
// ============================================================================

// CreateBookingLink handles POST /api/v1/brokers/booking-link
func (h *BrokerHandler) CreateBookingLink(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	userID, _ := r.Context().Value("user_id").(int64)

	var req models.CreateBookingLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	link, err := h.brokerService.CreateBookingLink(tenantID, &req, &userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create booking link: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(link)
}

// ListBookingLinks handles GET /api/v1/brokers/{brokerId}/bookings
func (h *BrokerHandler) ListBookingLinks(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	brokerID, err := strconv.ParseInt(vars["brokerId"], 10, 64)
	if err != nil {
		http.Error(w, "invalid broker id", http.StatusBadRequest)
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	status := r.URL.Query().Get("status")

	links, total, err := h.brokerService.ListBookingLinks(tenantID, brokerID, &status, offset, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list links: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  links,
		"total": total,
	})
}

// UpdateBookingLinkStatus handles PATCH /api/v1/brokers/{brokerId}/bookings/{linkId}/status
func (h *BrokerHandler) UpdateBookingLinkStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	userID, _ := r.Context().Value("user_id").(int64)

	vars := mux.Vars(r)
	brokerID, err := strconv.ParseInt(vars["brokerId"], 10, 64)
	if err != nil {
		http.Error(w, "invalid broker id", http.StatusBadRequest)
		return
	}

	linkID, err := strconv.ParseInt(vars["linkId"], 10, 64)
	if err != nil {
		http.Error(w, "invalid link id", http.StatusBadRequest)
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	status, ok := req["status"]
	if !ok {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	link, err := h.brokerService.UpdateBookingLinkStatus(tenantID, brokerID, linkID, status, &userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

// ============================================================================
// Commission Payout Handlers
// ============================================================================

// CreatePayout handles POST /api/v1/brokers/payout
func (h *BrokerHandler) CreatePayout(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	userID, _ := r.Context().Value("user_id").(int64)

	var req models.CreatePayoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	payout, err := h.brokerService.CreatePayout(tenantID, &req, &userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create payout: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payout)
}

// ListPayouts handles GET /api/v1/brokers/payouts
func (h *BrokerHandler) ListPayouts(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	brokerIDStr := r.URL.Query().Get("broker_id")
	var brokerID *int64
	if brokerIDStr != "" {
		id, err := strconv.ParseInt(brokerIDStr, 10, 64)
		if err == nil {
			brokerID = &id
		}
	}

	status := r.URL.Query().Get("status")

	payouts, total, err := h.brokerService.ListPayouts(tenantID, brokerID, &status, offset, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list payouts: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  payouts,
		"total": total,
	})
}

// UpdatePayoutStatus handles PATCH /api/v1/brokers/payouts/{payoutId}/status
func (h *BrokerHandler) UpdatePayoutStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	userID, _ := r.Context().Value("user_id").(int64)

	vars := mux.Vars(r)
	payoutID, err := strconv.ParseInt(vars["payoutId"], 10, 64)
	if err != nil {
		http.Error(w, "invalid payout id", http.StatusBadRequest)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	status, ok := req["status"].(string)
	if !ok {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	var paymentDate *string
	if pd, ok := req["payment_date"].(string); ok {
		paymentDate = &pd
	}

	payout, err := h.brokerService.UpdatePayoutStatus(tenantID, payoutID, status, paymentDate, &userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to update payout: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payout)
}

// ============================================================================
// Reporting & Analytics Handlers
// ============================================================================

// GetBrokerPerformance handles GET /api/v1/brokers/{brokerId}/performance
func (h *BrokerHandler) GetBrokerPerformance(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	brokerID, err := strconv.ParseInt(vars["brokerId"], 10, 64)
	if err != nil {
		http.Error(w, "invalid broker id", http.StatusBadRequest)
		return
	}

	performance, err := h.brokerService.GetBrokerPerformance(tenantID, brokerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get performance: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(performance)
}

// GetTopPerformers handles GET /api/v1/brokers/reports/top-performers
func (h *BrokerHandler) GetTopPerformers(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	brokers, err := h.brokerService.GetTopPerformingBrokers(tenantID, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get brokers: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brokers)
}

// GetCommissionDueReport handles GET /api/v1/brokers/reports/commission-due
func (h *BrokerHandler) GetCommissionDueReport(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value("tenant_id").(int64)
	if !ok {
		http.Error(w, "invalid tenant", http.StatusUnauthorized)
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	report, total, err := h.brokerService.GetCommissionDueReport(tenantID, offset, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get report: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  report,
		"total": total,
	})
}
