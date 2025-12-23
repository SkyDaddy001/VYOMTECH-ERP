package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

// JointApplicantHandler handles joint applicant HTTP requests
type JointApplicantHandler struct {
	Service *services.JointApplicantService
}

// NewJointApplicantHandler creates a new joint applicant handler
func NewJointApplicantHandler(service *services.JointApplicantService) *JointApplicantHandler {
	return &JointApplicantHandler{Service: service}
}

// CreateJointApplicant creates a new joint applicant
func (h *JointApplicantHandler) CreateJointApplicant(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	userID := getJAUserID(r)

	var req models.CreateJointApplicantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	applicant, err := h.Service.CreateJointApplicant(tenantID, &req, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, applicant)
}

// GetJointApplicant gets a joint applicant by ID
func (h *JointApplicantHandler) GetJointApplicant(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	vars := mux.Vars(r)
	applicantID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid applicant ID")
		return
	}

	applicant, err := h.Service.GetJointApplicant(tenantID, applicantID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, applicant)
}

// ListJointApplicants lists joint applicants for a booking
func (h *JointApplicantHandler) ListJointApplicants(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	bookingID, err := strconv.ParseInt(r.URL.Query().Get("booking_id"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	offset := getOffsetFromQuery(r)
	limit := getLimitFromQuery(r)

	applicants, total, err := h.Service.ListJointApplicants(tenantID, bookingID, offset, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSONPaginated(w, http.StatusOK, applicants, total)
}

// UpdateJointApplicant updates a joint applicant
func (h *JointApplicantHandler) UpdateJointApplicant(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	userID := getJAUserID(r)
	vars := mux.Vars(r)
	applicantID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid applicant ID")
		return
	}

	var req models.UpdateJointApplicantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	applicant, err := h.Service.UpdateJointApplicant(tenantID, applicantID, &req, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, applicant)
}

// DeleteJointApplicant deletes a joint applicant
func (h *JointApplicantHandler) DeleteJointApplicant(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	userID := getJAUserID(r)
	vars := mux.Vars(r)
	applicantID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid applicant ID")
		return
	}

	err = h.Service.DeleteJointApplicant(tenantID, applicantID, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// CreateCoOwnershipAgreement creates a co-ownership agreement
func (h *JointApplicantHandler) CreateCoOwnershipAgreement(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	userID := getJAUserID(r)

	var req models.CreateCoOwnershipAgreementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	agreement, err := h.Service.CreateCoOwnershipAgreement(tenantID, &req, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, agreement)
}

// GetCoOwnershipAgreement gets a co-ownership agreement
func (h *JointApplicantHandler) GetCoOwnershipAgreement(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	vars := mux.Vars(r)
	agreementID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid agreement ID")
		return
	}

	agreement, err := h.Service.GetCoOwnershipAgreement(tenantID, agreementID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, agreement)
}

// ListCoOwnershipAgreements lists co-ownership agreements
func (h *JointApplicantHandler) ListCoOwnershipAgreements(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	status := r.URL.Query().Get("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	offset := getOffsetFromQuery(r)
	limit := getLimitFromQuery(r)

	agreements, total, err := h.Service.ListCoOwnershipAgreements(tenantID, statusPtr, offset, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSONPaginated(w, http.StatusOK, agreements, total)
}

// UpdateCoOwnershipAgreementStatus updates agreement status
func (h *JointApplicantHandler) UpdateCoOwnershipAgreementStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	userID := getJAUserID(r)
	vars := mux.Vars(r)
	agreementID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid agreement ID")
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	status, ok := req["status"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Status field required")
		return
	}

	agreement, err := h.Service.UpdateCoOwnershipAgreementStatus(tenantID, agreementID, status, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, agreement)
}

// UploadDocument uploads a document
func (h *JointApplicantHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	userID := getJAUserID(r)
	vars := mux.Vars(r)
	applicantID, err := strconv.ParseInt(vars["applicant_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid applicant ID")
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	documentType, ok := req["document_type"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Document type required")
		return
	}

	documentName := req["document_name"]
	documentURL := req["document_url"]

	doc, err := h.Service.UploadDocument(tenantID, applicantID, documentType, documentName, documentURL, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, doc)
}

// VerifyDocument verifies a document
func (h *JointApplicantHandler) VerifyDocument(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	userID := getJAUserID(r)
	vars := mux.Vars(r)
	documentID, err := strconv.ParseInt(vars["document_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	status, ok := req["status"].(string)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Status field required")
		return
	}

	var notes *string
	if n, ok := req["notes"].(string); ok {
		notes = &n
	}

	doc, err := h.Service.VerifyDocument(tenantID, documentID, status, notes, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, doc)
}

// CreateIncomeVerification creates income verification
func (h *JointApplicantHandler) CreateIncomeVerification(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	userID := getJAUserID(r)
	vars := mux.Vars(r)
	applicantID, err := strconv.ParseInt(vars["applicant_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid applicant ID")
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var annualIncome *float64
	if ai, ok := req["annual_income"].(float64); ok {
		annualIncome = &ai
	}

	verification, err := h.Service.CreateIncomeVerification(tenantID, applicantID, annualIncome, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, verification)
}

// VerifyIncome verifies applicant income
func (h *JointApplicantHandler) VerifyIncome(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	userID := getJAUserID(r)
	vars := mux.Vars(r)
	verificationID, err := strconv.ParseInt(vars["verification_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid verification ID")
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	status, ok := req["status"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Status field required")
		return
	}

	verification, err := h.Service.VerifyIncome(tenantID, verificationID, status, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, verification)
}

// AddLiability adds a liability
func (h *JointApplicantHandler) AddLiability(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	userID := getJAUserID(r)
	vars := mux.Vars(r)
	applicantID, err := strconv.ParseInt(vars["applicant_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid applicant ID")
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	liabilityType, _ := req["liability_type"].(string)
	creditorName, _ := req["creditor_name"].(string)

	var outstanding *float64
	if o, ok := req["outstanding_amount"].(float64); ok {
		outstanding = &o
	}

	liability, err := h.Service.AddLiability(tenantID, applicantID, liabilityType, creditorName, outstanding, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, liability)
}

// GetJointApplicantSummary gets applicant summary for a booking
func (h *JointApplicantHandler) GetJointApplicantSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := getJATenantID(r)
	vars := mux.Vars(r)
	bookingID, err := strconv.ParseInt(vars["booking_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	summary, err := h.Service.GetJointApplicantSummary(tenantID, bookingID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, summary)
}

// Helper functions
func getJATenantID(r *http.Request) int64 {
	tenantIDStr, ok := r.Context().Value("tenant_id").(string)
	if ok {
		tenantID, _ := strconv.ParseInt(tenantIDStr, 10, 64)
		return tenantID
	}
	if tenantID, ok := r.Context().Value("tenant_id").(int64); ok {
		return tenantID
	}
	return 0
}

func getJAUserID(r *http.Request) *int64 {
	if userID, ok := r.Context().Value("user_id").(int64); ok {
		return &userID
	}
	return nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func respondWithJSONPaginated(w http.ResponseWriter, code int, payload interface{}, total int64) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func getOffsetFromQuery(r *http.Request) int {
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		return 0
	}
	offset, _ := strconv.Atoi(offsetStr)
	return offset
}

func getLimitFromQuery(r *http.Request) int {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return 10
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit > 100 {
		limit = 100
	}
	return limit
}
