package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

// ProjectManagementHandler handles all project management operations
type ProjectManagementHandler struct {
	Service *services.ProjectManagementService
	DB      *sql.DB
}

// NewProjectManagementHandler creates a new project management handler
func NewProjectManagementHandler(service *services.ProjectManagementService, db *sql.DB) *ProjectManagementHandler {
	return &ProjectManagementHandler{
		Service: service,
		DB:      db,
	}
}

// ============================================================
// CUSTOMER PROFILE ENDPOINTS
// ============================================================

// CreateCustomerProfile handles POST /api/v1/project-management/customers
func (h *ProjectManagementHandler) CreateCustomerProfile(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.PropertyCustomerProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	customer, err := h.Service.CreateCustomerProfile(tenantID, &req)
	if err != nil {
		log.Printf("Error creating customer profile: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create customer profile")
		return
	}

	h.respondJSON(w, http.StatusCreated, customer)
}

// GetCustomerProfile handles GET /api/v1/project-management/customers/{id}
func (h *ProjectManagementHandler) GetCustomerProfile(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	customerID := mux.Vars(r)["id"]

	customer, err := h.Service.GetCustomerProfile(tenantID, customerID)
	if err != nil {
		if err.Error() == "customer not found" {
			h.respondError(w, http.StatusNotFound, "Customer not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch customer profile")
		return
	}

	h.respondJSON(w, http.StatusOK, customer)
}

// ============================================================
// AREA STATEMENT ENDPOINTS
// ============================================================

// CreateAreaStatement handles POST /api/v1/project-management/area-statements
func (h *ProjectManagementHandler) CreateAreaStatement(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.CreateAreaStatementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	statement, err := h.Service.CreateAreaStatement(tenantID, &req)
	if err != nil {
		log.Printf("Error creating area statement: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create area statement")
		return
	}

	h.respondJSON(w, http.StatusCreated, statement)
}

// ============================================================
// COST SHEET ENDPOINTS
// ============================================================

// UpdateCostSheet handles POST /api/v1/project-management/cost-sheets
func (h *ProjectManagementHandler) UpdateCostSheet(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.UpdateCostSheetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.Service.UpdateCostSheet(tenantID, &req)
	if err != nil {
		log.Printf("Error updating cost sheet: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update cost sheet")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "Cost sheet updated successfully",
		"unit_id": req.UnitID,
	})
}

// ============================================================
// PROJECT COST CONFIGURATION ENDPOINTS
// ============================================================

// CreateProjectCostConfiguration handles POST /api/v1/project-management/cost-configurations
func (h *ProjectManagementHandler) CreateProjectCostConfiguration(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.CreateProjectCostConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	config, err := h.Service.CreateProjectCostConfiguration(tenantID, &req)
	if err != nil {
		log.Printf("Error creating cost configuration: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create cost configuration")
		return
	}

	h.respondJSON(w, http.StatusCreated, config)
}

// ============================================================
// BANK FINANCING ENDPOINTS
// ============================================================

// CreateBankFinancing handles POST /api/v1/project-management/bank-financing
func (h *ProjectManagementHandler) CreateBankFinancing(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.CreateBankFinancingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	financing, err := h.Service.CreateBankFinancing(tenantID, &req)
	if err != nil {
		log.Printf("Error creating bank financing: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create bank financing record")
		return
	}

	h.respondJSON(w, http.StatusCreated, financing)
}

// ============================================================
// DISBURSEMENT SCHEDULE ENDPOINTS
// ============================================================

// CreateDisbursementSchedule handles POST /api/v1/project-management/disbursement-schedule
func (h *ProjectManagementHandler) CreateDisbursementSchedule(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.CreateDisbursementScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	schedule, err := h.Service.CreateDisbursementSchedule(tenantID, &req)
	if err != nil {
		log.Printf("Error creating disbursement schedule: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create disbursement schedule")
		return
	}

	h.respondJSON(w, http.StatusCreated, schedule)
}

// UpdateDisbursement handles PUT /api/v1/project-management/disbursement/{id}
func (h *ProjectManagementHandler) UpdateDisbursement(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	disbursementID := mux.Vars(r)["id"]

	var req models.UpdateDisbursementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.DisbursementID = disbursementID

	// Update disbursement in service
	actualDate := (*string)(nil)
	if req.ActualDisbursementDate != "" {
		actualDate = &req.ActualDisbursementDate
	}

	query := `UPDATE property_disbursement_schedule 
		SET actual_disbursement_date = $1, actual_disbursement_amount = $2,
		 disbursement_status = $3, cheque_no = $4, bank_reference_no = $5,
		 neft_ref_id = $6, updated_at = NOW()
		WHERE id = $7 AND tenant_id = $8`

	_, err := h.DB.Exec(query,
		actualDate, req.ActualDisbursementAmount, req.DisbursementStatus,
		req.ChequeNo, req.BankReferenceNo, req.NEFTRefID, disbursementID, tenantID,
	)

	if err != nil {
		log.Printf("Error updating disbursement: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update disbursement")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message":         "Disbursement updated successfully",
		"disbursement_id": disbursementID,
	})
}

// ============================================================
// PAYMENT STAGE ENDPOINTS
// ============================================================

// CreatePaymentStage handles POST /api/v1/project-management/payment-stages
func (h *ProjectManagementHandler) CreatePaymentStage(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.CreatePaymentStageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	stage, err := h.Service.CreatePaymentStage(tenantID, &req)
	if err != nil {
		log.Printf("Error creating payment stage: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create payment stage")
		return
	}

	h.respondJSON(w, http.StatusCreated, stage)
}

// RecordPaymentCollection handles PUT /api/v1/project-management/payment-stages/{id}/collection
func (h *ProjectManagementHandler) RecordPaymentCollection(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	stageID := mux.Vars(r)["id"]

	var req models.UpdatePaymentStageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.Service.UpdatePaymentStageCollection(tenantID, stageID, &req)
	if err != nil {
		log.Printf("Error recording payment collection: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to record payment collection")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message":  "Payment collected successfully",
		"stage_id": stageID,
	})
}

// ============================================================
// REPORTING ENDPOINTS
// ============================================================

// GetBankFinancingReport handles GET /api/v1/project-management/reports/bank-financing
func (h *ProjectManagementHandler) GetBankFinancingReport(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	projectID := r.URL.Query().Get("project_id")

	query := `SELECT id, unit_id, apt_no, apartment_cost, sanctioned_amount, 
		total_disbursed_amount, remaining_disbursement, total_collection_from_unit,
		disbursement_status, collection_status, noc_received
		FROM property_bank_financing 
		WHERE tenant_id = $1`

	args := []interface{}{tenantID}
	if projectID != "" {
		query += " AND project_id = $2"
		args = append(args, projectID)
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch bank financing report")
		return
	}
	defer rows.Close()

	var reportData []map[string]interface{}
	for rows.Next() {
		var id, unitID, aptNo, disbursementStatus, collectionStatus string
		var apartmentCost, sanctionedAmount, totalDisbursed, remaining, totalCollection float64
		var nocReceived bool

		if err := rows.Scan(&id, &unitID, &aptNo, &apartmentCost, &sanctionedAmount,
			&totalDisbursed, &remaining, &totalCollection, &disbursementStatus, &collectionStatus, &nocReceived); err != nil {
			continue
		}

		reportData = append(reportData, map[string]interface{}{
			"id":                         id,
			"unit_id":                    unitID,
			"apt_no":                     aptNo,
			"apartment_cost":             apartmentCost,
			"sanctioned_amount":          sanctionedAmount,
			"total_disbursed_amount":     totalDisbursed,
			"remaining_disbursement":     remaining,
			"total_collection_from_unit": totalCollection,
			"disbursement_status":        disbursementStatus,
			"collection_status":          collectionStatus,
			"noc_received":               nocReceived,
		})
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"count": len(reportData),
		"data":  reportData,
	})
}

// GetPaymentStageReport handles GET /api/v1/project-management/reports/payment-stages
func (h *ProjectManagementHandler) GetPaymentStageReport(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	projectID := r.URL.Query().Get("project_id")
	unitID := r.URL.Query().Get("unit_id")

	query := `SELECT id, stage_name, stage_number, stage_percentage, stage_due_amount,
		amount_due, amount_received, amount_pending, collection_status, due_date,
		payment_received_date, payment_mode
		FROM property_payment_stage 
		WHERE tenant_id = $1`

	args := []interface{}{tenantID}
	idx := 2

	if projectID != "" {
		query += fmt.Sprintf(" AND project_id = $%d", idx)
		args = append(args, projectID)
		idx++
	}

	if unitID != "" {
		query += fmt.Sprintf(" AND unit_id = $%d", idx)
		args = append(args, unitID)
	}

	query += " ORDER BY stage_number"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch payment stage report")
		return
	}
	defer rows.Close()

	var reportData []map[string]interface{}
	for rows.Next() {
		var id, stageName, collectionStatus, paymentMode string
		var stageNumber int
		var stagePercentage, stageDueAmount, amountDue, amountReceived, amountPending float64
		var dueDate, paymentReceivedDate sql.NullTime

		if err := rows.Scan(&id, &stageName, &stageNumber, &stagePercentage, &stageDueAmount,
			&amountDue, &amountReceived, &amountPending, &collectionStatus, &dueDate,
			&paymentReceivedDate, &paymentMode); err != nil {
			continue
		}

		reportData = append(reportData, map[string]interface{}{
			"id":                    id,
			"stage_name":            stageName,
			"stage_number":          stageNumber,
			"stage_percentage":      stagePercentage,
			"stage_due_amount":      stageDueAmount,
			"amount_due":            amountDue,
			"amount_received":       amountReceived,
			"amount_pending":        amountPending,
			"collection_status":     collectionStatus,
			"due_date":              dueDate.Time,
			"payment_received_date": paymentReceivedDate.Time,
			"payment_mode":          paymentMode,
		})
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"count": len(reportData),
		"data":  reportData,
	})
}

// ============================================================
// LIST ENDPOINTS
// ============================================================

// ListCustomers handles GET /api/v1/project-management/customers
func (h *ProjectManagementHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	customers, count, err := h.Service.ListCustomerProfiles(tenantID, limit, offset)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch customers")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"count":  count,
		"limit":  limit,
		"offset": offset,
		"data":   customers,
	})
}

// ListAreaStatements handles GET /api/v1/project-management/area-statements
func (h *ProjectManagementHandler) ListAreaStatements(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	projectID := r.URL.Query().Get("project_id")

	if projectID == "" {
		h.respondError(w, http.StatusBadRequest, "project_id query parameter required")
		return
	}

	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	statements, count, err := h.Service.ListAreaStatements(tenantID, projectID, limit, offset)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch area statements")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"count":  count,
		"limit":  limit,
		"offset": offset,
		"data":   statements,
	})
}

// GetAreaStatement handles GET /api/v1/project-management/area-statements/{unit_id}
func (h *ProjectManagementHandler) GetAreaStatement(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	unitID := mux.Vars(r)["unit_id"]

	statement, err := h.Service.GetAreaStatement(tenantID, unitID)
	if err != nil {
		if err.Error() == "area statement not found" {
			h.respondError(w, http.StatusNotFound, "Area statement not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch area statement")
		return
	}

	h.respondJSON(w, http.StatusOK, statement)
}

// ListPaymentStages handles GET /api/v1/project-management/payment-stages
func (h *ProjectManagementHandler) ListPaymentStages(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	unitID := r.URL.Query().Get("unit_id")

	if unitID == "" {
		h.respondError(w, http.StatusBadRequest, "unit_id query parameter required")
		return
	}

	stages, err := h.Service.ListPaymentStages(tenantID, unitID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch payment stages")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"count": len(stages),
		"data":  stages,
	})
}

// GetBankFinancing handles GET /api/v1/project-management/bank-financing/{id}
func (h *ProjectManagementHandler) GetBankFinancing(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	financingID := mux.Vars(r)["id"]

	financing, err := h.Service.GetBankFinancing(tenantID, financingID)
	if err != nil {
		if err.Error() == "financing record not found" {
			h.respondError(w, http.StatusNotFound, "Financing record not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch financing record")
		return
	}

	h.respondJSON(w, http.StatusOK, financing)
}

// ============================================================
// UPDATE ENDPOINTS
// ============================================================

// UpdateCustomerProfile handles PUT /api/v1/project-management/customers/{id}
func (h *ProjectManagementHandler) UpdateCustomerProfile(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	customerID := mux.Vars(r)["id"]

	var req models.PropertyCustomerProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.Service.UpdateCustomerProfile(tenantID, customerID, &req)
	if err != nil {
		log.Printf("Error updating customer profile: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update customer profile")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "Customer profile updated successfully",
		"id":      customerID,
	})
}

// UpdateBankFinancing handles PUT /api/v1/project-management/bank-financing/{id}
func (h *ProjectManagementHandler) UpdateBankFinancing(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	financingID := mux.Vars(r)["id"]

	var req models.CreateBankFinancingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.Service.UpdateBankFinancing(tenantID, financingID, &req)
	if err != nil {
		log.Printf("Error updating bank financing: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update bank financing")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "Bank financing updated successfully",
		"id":      financingID,
	})
}

// DeleteAreaStatement handles DELETE /api/v1/project-management/area-statements/{unit_id}
func (h *ProjectManagementHandler) DeleteAreaStatement(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	unitID := mux.Vars(r)["unit_id"]

	err := h.Service.DeleteAreaStatement(tenantID, unitID)
	if err != nil {
		log.Printf("Error deleting area statement: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete area statement")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "Area statement deleted successfully",
		"unit_id": unitID,
	})
}

// ============================================================
// CALCULATION & DASHBOARD ENDPOINTS
// ============================================================

// GetCostBreakdown handles GET /api/v1/project-management/cost-breakdown/{unit_id}
func (h *ProjectManagementHandler) GetCostBreakdown(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	unitID := mux.Vars(r)["unit_id"]

	breakdown, err := h.Service.CalculateCostBreakdown(tenantID, unitID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"unit_id":   unitID,
		"breakdown": breakdown,
	})
}

// GetProjectSummary handles GET /api/v1/project-management/summary/{project_id}
func (h *ProjectManagementHandler) GetProjectSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	projectID := mux.Vars(r)["project_id"]

	summary, err := h.Service.GetProjectSummary(tenantID, projectID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch project summary")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"project_id": projectID,
		"summary":    summary,
	})
}

// GetCollectionStatus handles GET /api/v1/project-management/collection-status
func (h *ProjectManagementHandler) GetCollectionStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	projectID := r.URL.Query().Get("project_id")

	query := `SELECT collection_status, COUNT(*) as count, SUM(amount_due) as total_due, 
		SUM(amount_received) as total_received
		FROM property_payment_stage 
		WHERE tenant_id = $1`

	args := []interface{}{tenantID}
	if projectID != "" {
		query += " AND project_id = $2"
		args = append(args, projectID)
	}

	query += " GROUP BY collection_status"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch collection status")
		return
	}
	defer rows.Close()

	status := make(map[string]interface{})
	for rows.Next() {
		var collectionStatus string
		var count int
		var totalDue, totalReceived sql.NullFloat64

		if err := rows.Scan(&collectionStatus, &count, &totalDue, &totalReceived); err != nil {
			continue
		}

		status[collectionStatus] = map[string]interface{}{
			"count":          count,
			"total_due":      totalDue.Float64,
			"total_received": totalReceived.Float64,
		}
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"collection_status": status,
	})
}

// GetDisbursementStatus handles GET /api/v1/project-management/disbursement-status
func (h *ProjectManagementHandler) GetDisbursementStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	projectID := r.URL.Query().Get("project_id")

	query := `SELECT disbursement_status, COUNT(*) as count, 
		SUM(expected_disbursement_amount) as expected, 
		SUM(COALESCE(actual_disbursement_amount, 0)) as actual
		FROM property_disbursement_schedule 
		WHERE tenant_id = $1`

	args := []interface{}{tenantID}
	if projectID != "" {
		query += " AND financing_id IN (SELECT id FROM property_bank_financing WHERE project_id = $2)"
		args = append(args, projectID)
	}

	query += " GROUP BY disbursement_status"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch disbursement status")
		return
	}
	defer rows.Close()

	status := make(map[string]interface{})
	for rows.Next() {
		var disbursementStatus string
		var count int
		var expected, actual sql.NullFloat64

		if err := rows.Scan(&disbursementStatus, &count, &expected, &actual); err != nil {
			continue
		}

		status[disbursementStatus] = map[string]interface{}{
			"count":           count,
			"expected_amount": expected.Float64,
			"actual_amount":   actual.Float64,
		}
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"disbursement_status": status,
	})
}

// ============================================================
// RESPONSE HELPERS
// ============================================================

func (h *ProjectManagementHandler) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *ProjectManagementHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{
		"error": message,
	})
}
