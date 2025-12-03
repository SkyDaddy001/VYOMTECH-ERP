package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

type CompanyHandler struct {
	companyService *services.CompanyService
	logger         *logger.Logger
}

func NewCompanyHandler(companyService *services.CompanyService, logger *logger.Logger) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
		logger:         logger,
	}
}

// CreateCompany creates a new company under a tenant
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TenantID       string  `json:"tenant_id"`
		Name           string  `json:"name"`
		Description    string  `json:"description"`
		IndustryType   string  `json:"industry_type"`
		EmployeeCount  *int    `json:"employee_count"`
		Website        *string `json:"website"`
		MaxProjects    int     `json:"max_projects"`
		MaxUsers       int     `json:"max_users"`
		BillingEmail   string  `json:"billing_email"`
		BillingAddress string  `json:"billing_address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	company := &models.Company{
		ID:             fmt.Sprintf("comp_%d", time.Now().UnixNano()),
		TenantID:       req.TenantID,
		Name:           req.Name,
		Description:    req.Description,
		Status:         "active",
		IndustryType:   req.IndustryType,
		EmployeeCount:  req.EmployeeCount,
		Website:        req.Website,
		MaxProjects:    req.MaxProjects,
		MaxUsers:       req.MaxUsers,
		BillingEmail:   req.BillingEmail,
		BillingAddress: req.BillingAddress,
	}

	if err := h.companyService.CreateCompany(company); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create company: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(company)
}

// GetCompany retrieves a company
func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	companyID := r.URL.Query().Get("company_id")
	if companyID == "" {
		http.Error(w, "company_id is required", http.StatusBadRequest)
		return
	}

	company, err := h.companyService.GetCompany(companyID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get company: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// ListCompanies lists companies for a tenant
func (h *CompanyHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "tenant_id is required", http.StatusBadRequest)
		return
	}

	companies, err := h.companyService.ListCompaniesByTenant(tenantID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list companies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

// UpdateCompany updates company details
func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var company models.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.companyService.UpdateCompany(&company); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update company: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Company updated successfully"})
}

// CreateProject creates a new project
func (h *CompanyHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CompanyID       string     `json:"company_id"`
		TenantID        string     `json:"tenant_id"`
		Name            string     `json:"name"`
		Description     string     `json:"description"`
		ProjectType     string     `json:"project_type"`
		MaxUsers        int        `json:"max_users"`
		BudgetAllocated float64    `json:"budget_allocated"`
		StartDate       time.Time  `json:"start_date"`
		EndDate         *time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	project := &models.Project{
		ID:              fmt.Sprintf("proj_%d", time.Now().UnixNano()),
		CompanyID:       req.CompanyID,
		TenantID:        req.TenantID,
		Name:            req.Name,
		Description:     req.Description,
		Status:          "active",
		ProjectType:     req.ProjectType,
		MaxUsers:        req.MaxUsers,
		BudgetAllocated: req.BudgetAllocated,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
	}

	if err := h.companyService.CreateProject(project); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create project: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// GetProject retrieves a project
func (h *CompanyHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		http.Error(w, "project_id is required", http.StatusBadRequest)
		return
	}

	project, err := h.companyService.GetProject(projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get project: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// ListProjects lists projects for a company
func (h *CompanyHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	companyID := r.URL.Query().Get("company_id")
	if companyID == "" {
		http.Error(w, "company_id is required", http.StatusBadRequest)
		return
	}

	projects, err := h.companyService.ListProjectsByCompany(companyID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list projects: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// AddMemberToCompany adds a user to a company
func (h *CompanyHandler) AddMemberToCompany(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CompanyID  string `json:"company_id"`
		UserID     int    `json:"user_id"`
		TenantID   string `json:"tenant_id"`
		Role       string `json:"role"`
		Department string `json:"department"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	member := &models.CompanyMember{
		ID:         fmt.Sprintf("cmpmbr_%d", time.Now().UnixNano()),
		CompanyID:  req.CompanyID,
		UserID:     req.UserID,
		TenantID:   req.TenantID,
		Role:       req.Role,
		Department: req.Department,
		IsActive:   true,
	}

	if err := h.companyService.AddMemberToCompany(member); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add member: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)
}

// AddMemberToProject adds a user to a project
func (h *CompanyHandler) AddMemberToProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProjectID string `json:"project_id"`
		UserID    int    `json:"user_id"`
		CompanyID string `json:"company_id"`
		TenantID  string `json:"tenant_id"`
		Role      string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	member := &models.ProjectMember{
		ID:        fmt.Sprintf("projmbr_%d", time.Now().UnixNano()),
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
		CompanyID: req.CompanyID,
		TenantID:  req.TenantID,
		Role:      req.Role,
		JoinedAt:  time.Now(),
		IsActive:  true,
	}

	if err := h.companyService.AddMemberToProject(member); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add project member: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(member)
}

// GetCompanyMembers lists members of a company
func (h *CompanyHandler) GetCompanyMembers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	companyID := r.URL.Query().Get("company_id")
	if companyID == "" {
		http.Error(w, "company_id is required", http.StatusBadRequest)
		return
	}

	members, err := h.companyService.GetCompanyMembers(companyID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get members: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// GetProjectMembers lists members of a project
func (h *CompanyHandler) GetProjectMembers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		http.Error(w, "project_id is required", http.StatusBadRequest)
		return
	}

	members, err := h.companyService.GetProjectMembers(projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get members: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// RemoveProjectMember removes a user from a project
func (h *CompanyHandler) RemoveProjectMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project_id")
	userIDStr := r.URL.Query().Get("user_id")

	if projectID == "" || userIDStr == "" {
		http.Error(w, "project_id and user_id are required", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)
	if err := h.companyService.RemoveMemberFromProject(projectID, userID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove member: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "Member removed successfully"})
}
