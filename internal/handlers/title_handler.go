package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
)

// TitleHandler handles title clearance operations
type TitleHandler struct {
	service *services.TitleService
}

// NewTitleHandler creates a new title handler
func NewTitleHandler(service *services.TitleService) *TitleHandler {
	return &TitleHandler{service: service}
}

// Helper functions
func getTitleTenantID(r *http.Request) (int64, error) {
	tenantID := r.Context().Value("tenant_id")
	if tenantID == nil {
		return 0, errors.New("tenant_id not found in context")
	}
	return tenantID.(int64), nil
}

func getTitleUserID(r *http.Request) (int64, error) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		return 0, errors.New("user_id not found in context")
	}
	return userID.(int64), nil
}

func respondWithTitleJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getTitlePagination(r *http.Request) (int, int) {
	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}

// CreateTitleClearance creates a new title clearance
func (h *TitleHandler) CreateTitleClearance(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	var req models.CreateTitleClearanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.BookingID == 0 || req.ClearanceType == "" {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "booking_id and clearance_type are required"})
		return
	}

	clearance, err := h.service.CreateTitleClearance(tenantID, req.BookingID, req.ClearanceType, &req)
	if err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusCreated, clearance)
}

// GetTitleClearance retrieves a title clearance
func (h *TitleHandler) GetTitleClearance(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	clearanceID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid clearance ID"})
		return
	}

	clearance, err := h.service.GetTitleClearance(tenantID, clearanceID)
	if err != nil {
		if err.Error() == "clearance not found" {
			respondWithTitleJSON(w, http.StatusNotFound, map[string]string{"error": "Clearance not found"})
			return
		}
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, clearance)
}

// UpdateTitleClearance updates a title clearance
func (h *TitleHandler) UpdateTitleClearance(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	clearanceID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid clearance ID"})
		return
	}

	var req models.UpdateTitleClearanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.service.UpdateTitleClearance(tenantID, clearanceID, &req); err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, map[string]string{"message": "Clearance updated successfully"})
}

// ListTitleClearances lists all title clearances
func (h *TitleHandler) ListTitleClearances(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	limit, offset := getTitlePagination(r)

	clearances, total, err := h.service.ListTitleClearances(tenantID, limit, offset)
	if err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithTitleJSON(w, http.StatusOK, map[string]interface{}{
		"data":  clearances,
		"total": total,
	})
}

// CreateTitleIssue creates a title issue
func (h *TitleHandler) CreateTitleIssue(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	clearanceID, err := strconv.ParseInt(r.PathValue("clearance_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid clearance ID"})
		return
	}

	var req models.CreateTitleIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.IssueType == "" || req.IssueTitle == "" || req.Severity == "" {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "issue_type, issue_title, and severity are required"})
		return
	}

	issue, err := h.service.CreateTitleIssue(tenantID, clearanceID, &req)
	if err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusCreated, issue)
}

// GetTitleIssue retrieves a title issue
func (h *TitleHandler) GetTitleIssue(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	issueID, err := strconv.ParseInt(r.PathValue("issue_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid issue ID"})
		return
	}

	issue, err := h.service.GetTitleIssue(tenantID, issueID)
	if err != nil {
		if err.Error() == "issue not found" {
			respondWithTitleJSON(w, http.StatusNotFound, map[string]string{"error": "Issue not found"})
			return
		}
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, issue)
}

// UpdateTitleIssue updates a title issue
func (h *TitleHandler) UpdateTitleIssue(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	issueID, err := strconv.ParseInt(r.PathValue("issue_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid issue ID"})
		return
	}

	var req models.UpdateTitleIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.service.UpdateTitleIssue(tenantID, issueID, &req); err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, map[string]string{"message": "Issue updated successfully"})
}

// ResolveTitleIssue resolves a title issue
func (h *TitleHandler) ResolveTitleIssue(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getTitleUserID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	issueID, err := strconv.ParseInt(r.PathValue("issue_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid issue ID"})
		return
	}

	var req models.ResolveTitleIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Status == "" || req.ResolutionMethod == "" {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "status and resolution_method are required"})
		return
	}

	if err := h.service.ResolveTitleIssue(tenantID, issueID, userID, &req); err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, map[string]string{"message": "Issue resolved successfully"})
}

// ListTitleIssues lists issues for a clearance
func (h *TitleHandler) ListTitleIssues(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	clearanceID, err := strconv.ParseInt(r.PathValue("clearance_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid clearance ID"})
		return
	}

	limit, offset := getTitlePagination(r)

	issues, total, err := h.service.ListTitleIssues(tenantID, clearanceID, limit, offset)
	if err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithTitleJSON(w, http.StatusOK, map[string]interface{}{
		"data":  issues,
		"total": total,
	})
}

// CreateSearchReport creates a search report
func (h *TitleHandler) CreateSearchReport(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	clearanceID, err := strconv.ParseInt(r.PathValue("clearance_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid clearance ID"})
		return
	}

	var req models.CreateSearchReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.SearchType == "" {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "search_type is required"})
		return
	}

	report, err := h.service.CreateTitleSearchReport(tenantID, clearanceID, &req)
	if err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusCreated, report)
}

// GetSearchReport retrieves a search report
func (h *TitleHandler) GetSearchReport(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	reportID, err := strconv.ParseInt(r.PathValue("report_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid report ID"})
		return
	}

	report, err := h.service.GetTitleSearchReport(tenantID, reportID)
	if err != nil {
		if err.Error() == "search report not found" {
			respondWithTitleJSON(w, http.StatusNotFound, map[string]string{"error": "Report not found"})
			return
		}
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, report)
}

// VerifySearchReport verifies a search report
func (h *TitleHandler) VerifySearchReport(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getTitleUserID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	reportID, err := strconv.ParseInt(r.PathValue("report_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid report ID"})
		return
	}

	var req models.VerifySearchReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.SearchStatus == "" {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "search_status is required"})
		return
	}

	if err := h.service.VerifySearchReport(tenantID, reportID, userID, &req); err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, map[string]string{"message": "Report verified successfully"})
}

// ListSearchReports lists search reports
func (h *TitleHandler) ListSearchReports(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	clearanceID, err := strconv.ParseInt(r.PathValue("clearance_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid clearance ID"})
		return
	}

	limit, offset := getTitlePagination(r)

	reports, total, err := h.service.ListTitleSearchReports(tenantID, clearanceID, limit, offset)
	if err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithTitleJSON(w, http.StatusOK, map[string]interface{}{
		"data":  reports,
		"total": total,
	})
}

// CreateLegalOpinion creates a legal opinion
func (h *TitleHandler) CreateLegalOpinion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	clearanceID, err := strconv.ParseInt(r.PathValue("clearance_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid clearance ID"})
		return
	}

	var req models.CreateLegalOpinionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.OpinionType == "" {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "opinion_type is required"})
		return
	}

	opinion, err := h.service.CreateTitleLegalOpinion(tenantID, clearanceID, &req)
	if err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusCreated, opinion)
}

// GetLegalOpinion retrieves a legal opinion
func (h *TitleHandler) GetLegalOpinion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	opinionID, err := strconv.ParseInt(r.PathValue("opinion_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid opinion ID"})
		return
	}

	opinion, err := h.service.GetTitleLegalOpinion(tenantID, opinionID)
	if err != nil {
		if err.Error() == "legal opinion not found" {
			respondWithTitleJSON(w, http.StatusNotFound, map[string]string{"error": "Opinion not found"})
			return
		}
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, opinion)
}

// ReviewLegalOpinion reviews a legal opinion
func (h *TitleHandler) ReviewLegalOpinion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getTitleUserID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	opinionID, err := strconv.ParseInt(r.PathValue("opinion_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid opinion ID"})
		return
	}

	var req models.ReviewLegalOpinionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.OpinionStatus == "" || req.RiskAssessment == "" {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "opinion_status and risk_assessment are required"})
		return
	}

	if err := h.service.ReviewLegalOpinion(tenantID, opinionID, userID, &req); err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, map[string]string{"message": "Opinion reviewed successfully"})
}

// ListLegalOpinions lists legal opinions
func (h *TitleHandler) ListLegalOpinions(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	clearanceID, err := strconv.ParseInt(r.PathValue("clearance_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid clearance ID"})
		return
	}

	limit, offset := getTitlePagination(r)

	opinions, total, err := h.service.ListTitleLegalOpinions(tenantID, clearanceID, limit, offset)
	if err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithTitleJSON(w, http.StatusOK, map[string]interface{}{
		"data":  opinions,
		"total": total,
	})
}

// VerifyChecklistItem marks a checklist item as verified
func (h *TitleHandler) VerifyChecklistItem(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getTitleUserID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	checklistID, err := strconv.ParseInt(r.PathValue("checklist_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid checklist ID"})
		return
	}

	var payload map[string]interface{}
	json.NewDecoder(r.Body).Decode(&payload)

	var notes *string
	if n, ok := payload["verification_notes"].(string); ok {
		notes = &n
	}

	if err := h.service.VerifyChecklistItem(tenantID, checklistID, userID, notes); err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, map[string]string{"message": "Checklist item verified successfully"})
}

// ListVerificationChecklists lists checklist items
func (h *TitleHandler) ListVerificationChecklists(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	clearanceID, err := strconv.ParseInt(r.PathValue("clearance_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid clearance ID"})
		return
	}

	limit, offset := getTitlePagination(r)

	items, total, err := h.service.ListVerificationChecklists(tenantID, clearanceID, limit, offset)
	if err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithTitleJSON(w, http.StatusOK, map[string]interface{}{
		"data":  items,
		"total": total,
	})
}

// ApproveClearance approves a clearance
func (h *TitleHandler) ApproveClearance(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	approvalID, err := strconv.ParseInt(r.PathValue("approval_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid approval ID"})
		return
	}

	var req models.ApproveClearanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.ApprovalStatus == "" {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "approval_status is required"})
		return
	}

	if err := h.service.ApproveClearance(tenantID, approvalID, &req); err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, map[string]string{"message": "Clearance approved successfully"})
}

// ListClearanceApprovals lists approvals
func (h *TitleHandler) ListClearanceApprovals(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	clearanceID, err := strconv.ParseInt(r.PathValue("clearance_id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid clearance ID"})
		return
	}

	limit, offset := getTitlePagination(r)

	approvals, total, err := h.service.ListClearanceApprovals(tenantID, clearanceID, limit, offset)
	if err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithTitleJSON(w, http.StatusOK, map[string]interface{}{
		"data":  approvals,
		"total": total,
	})
}

// GetClearanceSummary gets a clearance summary
func (h *TitleHandler) GetClearanceSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getTitleTenantID(r)
	if err != nil {
		respondWithTitleJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	clearanceID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondWithTitleJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid clearance ID"})
		return
	}

	summary, err := h.service.GetClearanceSummary(tenantID, clearanceID)
	if err != nil {
		respondWithTitleJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithTitleJSON(w, http.StatusOK, summary)
}
