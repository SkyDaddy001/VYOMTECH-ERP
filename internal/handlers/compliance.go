package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/internal/services"
	"multi-tenant-ai-callcenter/pkg/logger"
)

// ComplianceHandler handles compliance and security endpoints
type ComplianceHandler struct {
	rbacService       *services.RBACService
	auditService      *services.AuditService
	encryptionService *services.EncryptionService
	gdprService       *services.GDPRService
	logger            *logger.Logger
}

// NewComplianceHandler creates a new compliance handler
func NewComplianceHandler(
	rbacService *services.RBACService,
	auditService *services.AuditService,
	encryptionService *services.EncryptionService,
	gdprService *services.GDPRService,
	logger *logger.Logger,
) *ComplianceHandler {
	return &ComplianceHandler{
		rbacService:       rbacService,
		auditService:      auditService,
		encryptionService: encryptionService,
		gdprService:       gdprService,
		logger:            logger,
	}
}

// CreateRoleRequest is the request body for creating a role
type CreateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// CreateRole creates a new role
// POST /api/v1/compliance/roles
func (ch *ComplianceHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	role := &models.Role{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		IsActive:    true,
	}

	err := ch.rbacService.CreateRole(ctx, role)
	if err != nil {
		ch.logger.Error("Failed to create role", "error", err)
		http.Error(w, "failed to create role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

// GetRoles retrieves all roles for a tenant
// GET /api/v1/compliance/roles
func (ch *ComplianceHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	roles, err := ch.rbacService.ListRoles(ctx, tenantID)
	if err != nil {
		ch.logger.Error("Failed to list roles", "error", err)
		http.Error(w, "failed to list roles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total": len(roles),
		"roles": roles,
	})
}

// GetAuditLogs retrieves audit logs
// GET /api/v1/compliance/audit-logs?limit=50&offset=0
func (ch *ComplianceHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	filters := make(map[string]interface{})

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		if parsed, err := strconv.ParseInt(userID, 10, 64); err == nil {
			filters["user_id"] = parsed
		}
	}

	if action := r.URL.Query().Get("action"); action != "" {
		filters["action"] = action
	}

	if resource := r.URL.Query().Get("resource"); resource != "" {
		filters["resource"] = resource
	}

	logs, err := ch.auditService.GetAuditLogs(ctx, tenantID, filters, limit, offset)
	if err != nil {
		ch.logger.Error("Failed to get audit logs", "error", err)
		http.Error(w, "failed to get audit logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":  len(logs),
		"limit":  limit,
		"offset": offset,
		"logs":   logs,
	})
}

// GetAuditSummary retrieves audit summary statistics
// GET /api/v1/compliance/audit-summary?days=30
func (ch *ComplianceHandler) GetAuditSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			days = parsed
		}
	}

	summary, err := ch.auditService.GetAuditSummary(ctx, tenantID, days)
	if err != nil {
		ch.logger.Error("Failed to get audit summary", "error", err)
		http.Error(w, "failed to get audit summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)
}

// GetSecurityEvents retrieves security events
// GET /api/v1/compliance/security-events
func (ch *ComplianceHandler) GetSecurityEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	filters := make(map[string]interface{})

	if eventType := r.URL.Query().Get("event_type"); eventType != "" {
		filters["event_type"] = eventType
	}

	if severity := r.URL.Query().Get("severity"); severity != "" {
		filters["severity"] = severity
	}

	events, err := ch.auditService.GetSecurityEvents(ctx, tenantID, filters, limit, offset)
	if err != nil {
		ch.logger.Error("Failed to get security events", "error", err)
		http.Error(w, "failed to get security events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":  len(events),
		"limit":  limit,
		"offset": offset,
		"events": events,
	})
}

// GetComplianceReport retrieves a compliance report
// GET /api/v1/compliance/report?start_date=2024-01-01&end_date=2024-12-31
func (ch *ComplianceHandler) GetComplianceReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = time.Now().AddDate(0, -1, 0) // Default to last month
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now()
	}

	report, err := ch.auditService.GetComplianceReport(ctx, tenantID, startDate, endDate)
	if err != nil {
		ch.logger.Error("Failed to get compliance report", "error", err)
		http.Error(w, "failed to get compliance report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}

// RequestDataAccess creates a GDPR data access request
// POST /api/v1/compliance/gdpr/access
func (ch *ComplianceHandler) RequestDataAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		http.Error(w, "user id not found", http.StatusBadRequest)
		return
	}

	request, err := ch.gdprService.CreateDataAccessRequest(ctx, tenantID, userID)
	if err != nil {
		ch.logger.Error("Failed to create data access request", "error", err)
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request)
}

// ExportUserData exports all user data
// GET /api/v1/compliance/gdpr/export
func (ch *ComplianceHandler) ExportUserData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		http.Error(w, "user id not found", http.StatusBadRequest)
		return
	}

	userData, err := ch.gdprService.ExportUserData(ctx, tenantID, userID)
	if err != nil {
		ch.logger.Error("Failed to export user data", "error", err)
		http.Error(w, "failed to export data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=user_data.json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userData)
}

// RequestDataDeletion creates a GDPR data deletion request
// POST /api/v1/compliance/gdpr/deletion
type DataDeletionRequest struct {
	Reason string `json:"reason"`
}

func (ch *ComplianceHandler) RequestDataDeletion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		http.Error(w, "user id not found", http.StatusBadRequest)
		return
	}

	var req DataDeletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	request, err := ch.gdprService.CreateDataDeletionRequest(ctx, tenantID, userID, req.Reason)
	if err != nil {
		ch.logger.Error("Failed to create deletion request", "error", err)
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request)
}

// GetUserConsent retrieves user consent records
// GET /api/v1/compliance/gdpr/consent
func (ch *ComplianceHandler) GetUserConsent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		http.Error(w, "user id not found", http.StatusBadRequest)
		return
	}

	consents, err := ch.gdprService.GetUserConsents(ctx, tenantID, userID)
	if err != nil {
		ch.logger.Error("Failed to get consents", "error", err)
		http.Error(w, "failed to get consents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":    len(consents),
		"consents": consents,
	})
}

// RecordConsent records user consent
// POST /api/v1/compliance/gdpr/consent
type ConsentRequest struct {
	Type  string `json:"type"` // marketing, analytics, third_party
	Given bool   `json:"given"`
}

func (ch *ComplianceHandler) RecordConsent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		http.Error(w, "user id not found", http.StatusBadRequest)
		return
	}

	var req ConsentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := ch.gdprService.RecordConsent(ctx, tenantID, userID, req.Type, req.Given)
	if err != nil {
		ch.logger.Error("Failed to record consent", "error", err)
		http.Error(w, "failed to record consent", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Consent recorded successfully",
	})
}
