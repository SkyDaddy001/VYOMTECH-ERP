package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

// PossessionHandler handles possession management HTTP requests
type PossessionHandler struct {
	service *services.PossessionService
}

// NewPossessionHandler creates a new possession handler
func NewPossessionHandler(service *services.PossessionService) *PossessionHandler {
	return &PossessionHandler{
		service: service,
	}
}

// Helper functions to extract tenant and user IDs
func getPossTenantID(r *http.Request) int64 {
	tenantID := r.Header.Get("X-Tenant-ID")
	if id, err := strconv.ParseInt(tenantID, 10, 64); err == nil {
		return id
	}
	return 0
}

func getPossUserID(r *http.Request) *int64 {
	userID := r.Header.Get("X-User-ID")
	if id, err := strconv.ParseInt(userID, 10, 64); err == nil {
		return &id
	}
	return nil
}

// Helper to respond with JSON
func respondWithPossJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Helper to get offset and limit from query parameters
func getPossOffsetFromQuery(r *http.Request) int {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	return offset
}

func getPossLimitFromQuery(r *http.Request) int {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return limit
}

// ============================================
// Possession Status Endpoints
// ============================================

// CreatePossession creates a new possession status
func (h *PossessionHandler) CreatePossession(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePossessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	tenantID := getPossTenantID(r)
	userID := getPossUserID(r)

	ps, err := h.service.CreatePossessionStatus(tenantID, req.BookingID, req.Status, req.PossessionType, req.Notes, userID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusCreated, ps)
}

// GetPossession retrieves a possession status by ID
func (h *PossessionHandler) GetPossession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	ps, err := h.service.GetPossessionStatus(id)
	if err != nil {
		respondWithPossJSON(w, http.StatusNotFound, map[string]string{"error": "Possession not found"})
		return
	}

	respondWithPossJSON(w, http.StatusOK, ps)
}

// ListPossessions lists all possession statuses for a tenant
func (h *PossessionHandler) ListPossessions(w http.ResponseWriter, r *http.Request) {
	tenantID := getPossTenantID(r)
	limit := getPossLimitFromQuery(r)
	offset := getPossOffsetFromQuery(r)

	statuses, total, err := h.service.ListPossessionStatuses(tenantID, limit, offset)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithPossJSON(w, http.StatusOK, statuses)
}

// UpdatePossession updates a possession status
func (h *PossessionHandler) UpdatePossession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	var req models.UpdatePossessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	userID := getPossUserID(r)
	updates := make(map[string]interface{})

	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.PossessionDate != nil {
		updates["possession_date"] = *req.PossessionDate
	}
	if req.EstimatedPossessionDate != nil {
		updates["estimated_possession_date"] = *req.EstimatedPossessionDate
	}
	if req.IsComplete != nil {
		updates["is_complete"] = *req.IsComplete
	}
	if req.CompletionPercentage != nil {
		updates["completion_percentage"] = *req.CompletionPercentage
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}

	ps, err := h.service.UpdatePossessionStatus(id, updates, userID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusOK, ps)
}

// DeletePossession deletes a possession status
func (h *PossessionHandler) DeletePossession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	if err := h.service.DeletePossessionStatus(id); err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusOK, map[string]string{"message": "Possession deleted successfully"})
}

// ============================================
// Document Endpoints
// ============================================

// CreatePossessionDocument creates a new possession document
func (h *PossessionHandler) CreatePossessionDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	possessionID, _ := strconv.ParseInt(vars["possession_id"], 10, 64)

	var req models.CreatePossessionDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	tenantID := getPossTenantID(r)
	userID := getPossUserID(r)

	doc, err := h.service.CreatePossessionDocument(tenantID, possessionID, req.DocumentType, req.DocumentName, req.IsMandatory, req.Metadata, userID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusCreated, doc)
}

// ListPossessionDocuments lists documents for a possession
func (h *PossessionHandler) ListPossessionDocuments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	possessionID, _ := strconv.ParseInt(vars["possession_id"], 10, 64)

	limit := getPossLimitFromQuery(r)
	offset := getPossOffsetFromQuery(r)

	docs, total, err := h.service.ListPossessionDocuments(possessionID, limit, offset)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithPossJSON(w, http.StatusOK, docs)
}

// GetPossessionDocument retrieves a document
func (h *PossessionHandler) GetPossessionDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID, _ := strconv.ParseInt(vars["doc_id"], 10, 64)

	doc, err := h.service.GetPossessionDocument(docID)
	if err != nil {
		respondWithPossJSON(w, http.StatusNotFound, map[string]string{"error": "Document not found"})
		return
	}

	respondWithPossJSON(w, http.StatusOK, doc)
}

// UpdatePossessionDocument updates a document
func (h *PossessionHandler) UpdatePossessionDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID, _ := strconv.ParseInt(vars["doc_id"], 10, 64)

	var req models.UpdatePossessionDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	userID := getPossUserID(r)
	updates := make(map[string]interface{})

	if req.DocumentStatus != nil {
		updates["document_status"] = *req.DocumentStatus
	}
	if req.VerificationNotes != nil {
		updates["verification_notes"] = *req.VerificationNotes
	}

	doc, err := h.service.UpdatePossessionDocument(docID, updates, userID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusOK, doc)
}

// VerifyPossessionDocument verifies a document
func (h *PossessionHandler) VerifyPossessionDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID, _ := strconv.ParseInt(vars["doc_id"], 10, 64)

	var req models.VerifyPossessionDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	verifierID := getPossUserID(r)

	doc, err := h.service.VerifyPossessionDocument(docID, req.DocumentStatus, req.VerificationNotes, verifierID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusOK, doc)
}

// ============================================
// Registration Endpoints
// ============================================

// CreateRegistration creates a new registration
func (h *PossessionHandler) CreateRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	possessionID, _ := strconv.ParseInt(vars["possession_id"], 10, 64)

	var req models.CreatePossessionRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	tenantID := getPossTenantID(r)
	userID := getPossUserID(r)

	reg, err := h.service.CreatePossessionRegistration(tenantID, possessionID, req.RegistrationType, userID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusCreated, reg)
}

// ListRegistrations lists registrations for a possession
func (h *PossessionHandler) ListRegistrations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	possessionID, _ := strconv.ParseInt(vars["possession_id"], 10, 64)

	limit := getPossLimitFromQuery(r)
	offset := getPossOffsetFromQuery(r)

	regs, total, err := h.service.ListPossessionRegistrations(possessionID, limit, offset)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithPossJSON(w, http.StatusOK, regs)
}

// GetRegistration retrieves a registration
func (h *PossessionHandler) GetRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regID, _ := strconv.ParseInt(vars["reg_id"], 10, 64)

	reg, err := h.service.GetPossessionRegistration(regID)
	if err != nil {
		respondWithPossJSON(w, http.StatusNotFound, map[string]string{"error": "Registration not found"})
		return
	}

	respondWithPossJSON(w, http.StatusOK, reg)
}

// UpdateRegistration updates a registration
func (h *PossessionHandler) UpdateRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regID, _ := strconv.ParseInt(vars["reg_id"], 10, 64)

	var req models.UpdatePossessionRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	userID := getPossUserID(r)
	updates := make(map[string]interface{})

	if req.RegistrationStatus != nil {
		updates["registration_status"] = *req.RegistrationStatus
	}
	if req.AmountPaid != nil {
		updates["amount_paid"] = *req.AmountPaid
	}
	if req.AmountPending != nil {
		updates["amount_pending"] = *req.AmountPending
	}
	if req.ExpectedCompletionDate != nil {
		updates["expected_completion_date"] = *req.ExpectedCompletionDate
	}
	if req.Remarks != nil {
		updates["remarks"] = *req.Remarks
	}

	reg, err := h.service.UpdatePossessionRegistration(regID, updates, userID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusOK, reg)
}

// ApproveRegistration approves a registration
func (h *PossessionHandler) ApproveRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regID, _ := strconv.ParseInt(vars["reg_id"], 10, 64)

	var req models.ApproveRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	approverID := getPossUserID(r)

	reg, err := h.service.ApproveRegistration(regID, req.ApprovalStatus, req.ApprovalNotes, approverID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusOK, reg)
}

// ============================================
// Certificate Endpoints
// ============================================

// CreateCertificate creates a new certificate
func (h *PossessionHandler) CreateCertificate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	possessionID, _ := strconv.ParseInt(vars["possession_id"], 10, 64)

	var certType struct {
		CertificateType string `json:"certificate_type" binding:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&certType); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	tenantID := getPossTenantID(r)
	userID := getPossUserID(r)

	cert, err := h.service.CreatePossessionCertificate(tenantID, possessionID, certType.CertificateType, userID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusCreated, cert)
}

// ListCertificates lists certificates for a possession
func (h *PossessionHandler) ListCertificates(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	possessionID, _ := strconv.ParseInt(vars["possession_id"], 10, 64)

	limit := getPossLimitFromQuery(r)
	offset := getPossOffsetFromQuery(r)

	certs, total, err := h.service.ListPossessionCertificates(possessionID, limit, offset)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithPossJSON(w, http.StatusOK, certs)
}

// GetCertificate retrieves a certificate
func (h *PossessionHandler) GetCertificate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	certID, _ := strconv.ParseInt(vars["cert_id"], 10, 64)

	cert, err := h.service.GetPossessionCertificate(certID)
	if err != nil {
		respondWithPossJSON(w, http.StatusNotFound, map[string]string{"error": "Certificate not found"})
		return
	}

	respondWithPossJSON(w, http.StatusOK, cert)
}

// VerifyCertificate verifies a certificate
func (h *PossessionHandler) VerifyCertificate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	certID, _ := strconv.ParseInt(vars["cert_id"], 10, 64)

	var req struct {
		VerificationStatus string `json:"verification_status" binding:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	verifierID := getPossUserID(r)

	cert, err := h.service.VerifyCertificate(certID, req.VerificationStatus, verifierID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusOK, cert)
}

// ============================================
// Approval Endpoints
// ============================================

// CreateApproval creates a new approval
func (h *PossessionHandler) CreateApproval(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	possessionID, _ := strconv.ParseInt(vars["possession_id"], 10, 64)

	var req models.CreatePossessionApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	tenantID := getPossTenantID(r)
	userID := getPossUserID(r)

	approval, err := h.service.CreatePossessionApproval(tenantID, possessionID, req.ApproverID, req.ApprovalType, req.SequenceOrder, userID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusCreated, approval)
}

// ListApprovals lists approvals for a possession
func (h *PossessionHandler) ListApprovals(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	possessionID, _ := strconv.ParseInt(vars["possession_id"], 10, 64)

	limit := getPossLimitFromQuery(r)
	offset := getPossOffsetFromQuery(r)

	approvals, total, err := h.service.ListPossessionApprovals(possessionID, limit, offset)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithPossJSON(w, http.StatusOK, approvals)
}

// GetApproval retrieves an approval
func (h *PossessionHandler) GetApproval(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	approvalID, _ := strconv.ParseInt(vars["approval_id"], 10, 64)

	approval, err := h.service.GetPossessionApproval(approvalID)
	if err != nil {
		respondWithPossJSON(w, http.StatusNotFound, map[string]string{"error": "Approval not found"})
		return
	}

	respondWithPossJSON(w, http.StatusOK, approval)
}

// ApprovePossession approves a possession
func (h *PossessionHandler) ApprovePossession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	approvalID, _ := strconv.ParseInt(vars["approval_id"], 10, 64)

	var req models.ApprovePossessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPossJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	approverID := getPossUserID(r)

	approval, err := h.service.ApprovePossession(approvalID, req.ApprovalStatus, req.ApprovalNotes, approverID, req.IsFinalApproval)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusOK, approval)
}

// ============================================
// Summary Endpoint
// ============================================

// GetPossessionSummary retrieves possession summary
func (h *PossessionHandler) GetPossessionSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	possessionID, _ := strconv.ParseInt(vars["possession_id"], 10, 64)

	summary, err := h.service.GetPossessionSummary(possessionID)
	if err != nil {
		respondWithPossJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPossJSON(w, http.StatusOK, summary)
}
