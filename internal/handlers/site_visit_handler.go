package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

// SiteVisitHandler handles site visit HTTP requests
type SiteVisitHandler struct {
	Service *services.SiteVisitService
	Logger  *log.Logger
}

// NewSiteVisitHandler creates a new site visit handler
func NewSiteVisitHandler(service *services.SiteVisitService, logger *log.Logger) *SiteVisitHandler {
	return &SiteVisitHandler{Service: service, Logger: logger}
}

// Schedule Endpoints

// CreateSchedule creates a new site visit schedule
func (h *SiteVisitHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	userID, err := getSiteVisitUserID(r)
	if err != nil {
		h.Logger.Printf("Error getting user ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	var req models.CreateSiteVisitScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Printf("Invalid request body: %v", err)
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	serviceReq := &services.ScheduleVisitRequest{
		TenantID:      tenantID,
		LeadID:        req.LeadID,
		VisitorName:   req.VisitorName,
		VisitorPhone:  req.VisitorPhone,
		VisitorEmail:  req.VisitorEmail,
		ScheduledDate: req.ScheduledDate,
		ScheduledBy:   userID,
		Status:        req.Status,
	}

	schedule, err := h.Service.CreateSchedule(r.Context(), serviceReq)
	if err != nil {
		h.Logger.Printf("Error creating schedule: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusCreated, schedule)
}

// GetSchedule gets a site visit schedule by ID
func (h *SiteVisitHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	scheduleID := vars["scheduleId"]
	if scheduleID == "" {
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Schedule ID required"})
		return
	}

	schedule, err := h.Service.GetSchedule(r.Context(), scheduleID, tenantID)
	if err != nil {
		h.Logger.Printf("Error getting schedule: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, schedule)
}

// ListSchedules lists site visit schedules with filters
func (h *SiteVisitHandler) ListSchedules(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	limit, offset := getSiteVisitPagination(r)

	filter := services.ScheduleFilter{}

	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = &status
	}
	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}
	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}
	if leadID := r.URL.Query().Get("lead_id"); leadID != "" {
		filter.LeadID = &leadID
	}
	if scheduledBy := r.URL.Query().Get("scheduled_by"); scheduledBy != "" {
		filter.ScheduledBy = &scheduledBy
	}

	schedules, err := h.Service.ListSchedules(r.Context(), tenantID, filter, limit, offset)
	if err != nil {
		h.Logger.Printf("Error listing schedules: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, schedules)
}

// UpdateScheduleStatus updates the status of a site visit schedule
func (h *SiteVisitHandler) UpdateScheduleStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	scheduleID := vars["scheduleId"]
	if scheduleID == "" {
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Schedule ID required"})
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Printf("Invalid request body: %v", err)
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	status, ok := req["status"]
	if !ok {
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Status required"})
		return
	}

	err = h.Service.UpdateScheduleStatus(r.Context(), scheduleID, tenantID, status)
	if err != nil {
		h.Logger.Printf("Error updating schedule status: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, map[string]string{"message": "Status updated successfully"})
}

// CancelSchedule cancels a site visit schedule
func (h *SiteVisitHandler) CancelSchedule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	scheduleID := vars["scheduleId"]
	if scheduleID == "" {
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Schedule ID required"})
		return
	}

	err = h.Service.CancelSchedule(r.Context(), scheduleID, tenantID)
	if err != nil {
		h.Logger.Printf("Error canceling schedule: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, map[string]string{"message": "Schedule canceled successfully"})
}

// Visit Log Endpoints

// CheckInVisit checks in a visit
func (h *SiteVisitHandler) CheckInVisit(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	userID, err := getSiteVisitUserID(r)
	if err != nil {
		h.Logger.Printf("Error getting user ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	var req models.CreateSiteVisitLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Printf("Invalid request body: %v", err)
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	now := time.Now()
	req.CheckInTime = &now
	req.VisitedBy = userID

	var unitsViewed []string
	if req.UnitsViewed != nil {
		unitsViewed = *req.UnitsViewed
	}

	serviceReq := &services.CreateVisitLogRequest{
		TenantID:         tenantID,
		VisitScheduleID:  req.VisitScheduleID,
		CheckInTime:      req.CheckInTime,
		CheckOutTime:     req.CheckOutTime,
		VisitedBy:        req.VisitedBy,
		UnitsViewed:      unitsViewed,
		Feedback:         req.Feedback,
		FollowUpRequired: req.FollowUpRequired,
		NextFollowupDate: req.NextFollowupDate,
	}

	logEntry, err := h.Service.CreateVisitLog(r.Context(), serviceReq)
	if err != nil {
		h.Logger.Printf("Error checking in visit: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusCreated, logEntry)
}

// CheckOutVisit checks out a visit
func (h *SiteVisitHandler) CheckOutVisit(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	logID := vars["logId"]
	if logID == "" {
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Log ID required"})
		return
	}

	var req models.UpdateSiteVisitLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Printf("Invalid request body: %v", err)
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	now := time.Now()
	req.CheckOutTime = &now

	var unitsViewed []string
	if req.UnitsViewed != nil {
		unitsViewed = *req.UnitsViewed
	}

	err = h.Service.UpdateVisitLog(r.Context(), logID, tenantID, &services.UpdateVisitLogRequest{
		CheckInTime:      req.CheckInTime,
		CheckOutTime:     req.CheckOutTime,
		UnitsViewed:      unitsViewed,
		Feedback:         req.Feedback,
		FollowUpRequired: req.FollowUpRequired,
		NextFollowupDate: req.NextFollowupDate,
	})
	if err != nil {
		h.Logger.Printf("Error checking out visit: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, map[string]string{"message": "Checked out successfully"})
}

// GetVisitLog gets a visit log by ID
func (h *SiteVisitHandler) GetVisitLog(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	logID := vars["logId"]
	if logID == "" {
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Log ID required"})
		return
	}

	logEntry, err := h.Service.GetVisitLog(r.Context(), logID, tenantID)
	if err != nil {
		h.Logger.Printf("Error getting visit log: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, logEntry)
}

// ListVisitLogs lists visit logs for a schedule
func (h *SiteVisitHandler) ListVisitLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	scheduleID := vars["scheduleId"]
	if scheduleID == "" {
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Schedule ID required"})
		return
	}

	logs, err := h.Service.ListVisitLogs(r.Context(), scheduleID, tenantID)
	if err != nil {
		h.Logger.Printf("Error listing visit logs: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, logs)
}

// SubmitFeedback submits feedback for a visit
func (h *SiteVisitHandler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	logID := vars["logId"]
	if logID == "" {
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Log ID required"})
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Printf("Invalid request body: %v", err)
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	feedback, _ := req["feedback"].(string)
	followUp, _ := req["follow_up_required"].(bool)

	updateReq := &services.UpdateVisitLogRequest{
		Feedback:         &feedback,
		FollowUpRequired: &followUp,
	}

	err = h.Service.UpdateVisitLog(r.Context(), logID, tenantID, updateReq)
	if err != nil {
		h.Logger.Printf("Error submitting feedback: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, map[string]string{"message": "Feedback submitted successfully"})
}

// Analytics Endpoints

// GetLeadVisitHistory gets visit history for a lead
func (h *SiteVisitHandler) GetLeadVisitHistory(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	leadID := vars["leadId"]
	if leadID == "" {
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Lead ID required"})
		return
	}

	history, err := h.Service.GetLeadVisitHistory(r.Context(), leadID, tenantID)
	if err != nil {
		h.Logger.Printf("Error getting lead visit history: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, history)
}

// GetUserStats gets visit statistics for a user
func (h *SiteVisitHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	userID, err := getSiteVisitUserID(r)
	if err != nil {
		h.Logger.Printf("Error getting user ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	startDate := time.Now().AddDate(0, -1, 0) // Last month
	endDate := time.Now()

	stats, err := h.Service.GetUserVisitStats(r.Context(), userID, tenantID, startDate, endDate)
	if err != nil {
		h.Logger.Printf("Error getting user stats: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, stats)
}

// GetProjectSummary gets a summary of visits for a project
func (h *SiteVisitHandler) GetProjectSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	vars := mux.Vars(r)
	projectID := vars["projectId"]
	if projectID == "" {
		respondWithSiteVisitJSON(w, http.StatusBadRequest, map[string]string{"error": "Project ID required"})
		return
	}

	summary, err := h.Service.GetProjectVisitSummary(r.Context(), projectID, tenantID)
	if err != nil {
		h.Logger.Printf("Error getting project summary: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, summary)
}

// ListPendingFollowups lists pending follow-ups
func (h *SiteVisitHandler) ListPendingFollowups(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getSiteVisitTenantID(r)
	if err != nil {
		h.Logger.Printf("Error getting tenant ID: %v", err)
		respondWithSiteVisitJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	followups, err := h.Service.ListPendingFollowups(r.Context(), tenantID)
	if err != nil {
		h.Logger.Printf("Error listing pending followups: %v", err)
		respondWithSiteVisitJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithSiteVisitJSON(w, http.StatusOK, followups)
}

// Helper functions

func getSiteVisitTenantID(r *http.Request) (string, error) {
	tenantID := r.Context().Value("tenant_id")
	if tenantID == nil {
		return "", http.ErrNoCookie
	}
	if tid, ok := tenantID.(string); ok {
		return tid, nil
	}
	return "", http.ErrNoCookie
}

func getSiteVisitUserID(r *http.Request) (string, error) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		return "", http.ErrNoCookie
	}
	if uid, ok := userID.(string); ok {
		return uid, nil
	}
	return "", http.ErrNoCookie
}

func respondWithSiteVisitJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func getSiteVisitPagination(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limit = 10
	} else {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		} else {
			limit = 10
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offset = 0
	} else {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		} else {
			offset = 0
		}
	}

	return limit, offset
}
