package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"multi-tenant-ai-callcenter/internal/services"
	"multi-tenant-ai-callcenter/pkg/logger"
)

// AnalyticsHandler handles analytics and reporting requests
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	logger           *logger.Logger
}

// NewAnalyticsHandler creates a new AnalyticsHandler
func NewAnalyticsHandler(analyticsService *services.AnalyticsService, logger *logger.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		logger:           logger,
	}
}

// GenerateReportRequest is the request body for generating a report
type GenerateReportRequest struct {
	Type      string                 `json:"type"`
	StartDate string                 `json:"start_date"`
	EndDate   string                 `json:"end_date"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	Format    string                 `json:"format,omitempty"`
}

// GenerateReport generates a custom report
// POST /api/v1/analytics/reports
func (ah *AnalyticsHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ah.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req GenerateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		http.Error(w, "invalid end_date format", http.StatusBadRequest)
		return
	}

	if req.Format == "" {
		req.Format = "json"
	}

	reportReq := &services.ReportRequest{
		TenantID:  tenantID,
		Type:      services.ReportType(req.Type),
		StartDate: startDate,
		EndDate:   endDate,
		Filters:   req.Filters,
		Format:    req.Format,
	}

	report, err := ah.analyticsService.GenerateReport(ctx, reportReq)
	if err != nil {
		ah.logger.Error("Failed to generate report", "error", err)
		http.Error(w, "failed to generate report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}

// GetTrends retrieves trend data
// GET /api/v1/analytics/trends?metric=leads&start_date=2024-01-01&end_date=2024-01-31
func (ah *AnalyticsHandler) GetTrends(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ah.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	metric := r.URL.Query().Get("metric")
	if metric == "" {
		http.Error(w, "metric parameter required", http.StatusBadRequest)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "invalid start_date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "invalid end_date format", http.StatusBadRequest)
		return
	}

	trends, err := ah.analyticsService.GetTrends(ctx, tenantID, metric, startDate, endDate)
	if err != nil {
		ah.logger.Error("Failed to get trends", "error", err)
		http.Error(w, "failed to get trends", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metric": metric,
		"data":   trends,
	})
}

// ExportReport exports a report in the specified format
// POST /api/v1/analytics/export
func (ah *AnalyticsHandler) ExportReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ah.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req GenerateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		http.Error(w, "invalid end_date format", http.StatusBadRequest)
		return
	}

	reportReq := &services.ReportRequest{
		TenantID:  tenantID,
		Type:      services.ReportType(req.Type),
		StartDate: startDate,
		EndDate:   endDate,
		Filters:   req.Filters,
		Format:    req.Format,
	}

	report, err := ah.analyticsService.GenerateReport(ctx, reportReq)
	if err != nil {
		ah.logger.Error("Failed to generate report", "error", err)
		http.Error(w, "failed to generate report", http.StatusInternalServerError)
		return
	}

	data, err := ah.analyticsService.ExportReportData(report, req.Format)
	if err != nil {
		ah.logger.Error("Failed to export report", "error", err)
		http.Error(w, "failed to export report", http.StatusInternalServerError)
		return
	}

	// Set appropriate content type based on format
	contentType := "application/json"
	fileExtension := "json"
	switch req.Format {
	case "csv":
		contentType = "text/csv"
		fileExtension = "csv"
	case "pdf":
		contentType = "application/pdf"
		fileExtension = "pdf"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=report."+fileExtension)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetCustomMetrics retrieves custom metrics
// GET /api/v1/analytics/metrics?metric=lead_source_distribution
func (ah *AnalyticsHandler) GetCustomMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		ah.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	metric := r.URL.Query().Get("metric")
	if metric == "" {
		http.Error(w, "metric parameter required", http.StatusBadRequest)
		return
	}

	metrics, err := ah.analyticsService.GetCustomMetrics(ctx, tenantID, metric, nil)
	if err != nil {
		ah.logger.Error("Failed to get metrics", "error", err)
		http.Error(w, "failed to get metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}
