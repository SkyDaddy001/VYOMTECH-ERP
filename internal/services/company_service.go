package services

import (
	"database/sql"
	"fmt"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/pkg/logger"
)

type CompanyService struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewCompanyService(db *sql.DB, logger *logger.Logger) *CompanyService {
	return &CompanyService{
		db:     db,
		logger: logger,
	}
}

// CreateCompany creates a new company under a tenant
func (s *CompanyService) CreateCompany(company *models.Company) error {
	query := `
		INSERT INTO companies (
			id, tenant_id, name, description, status, industry_type,
			employee_count, website, max_projects, max_users,
			billing_email, billing_address
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		company.ID, company.TenantID, company.Name, company.Description,
		company.Status, company.IndustryType, company.EmployeeCount,
		company.Website, company.MaxProjects, company.MaxUsers,
		company.BillingEmail, company.BillingAddress,
	)

	if err != nil {
		s.logger.Error("Failed to create company", "error", err, "tenant_id", company.TenantID)
		return fmt.Errorf("failed to create company: %w", err)
	}

	s.logger.Info("Company created", "company_id", company.ID, "tenant_id", company.TenantID)
	return nil
}

// GetCompany retrieves a company by ID
func (s *CompanyService) GetCompany(companyID string) (*models.Company, error) {
	query := `
		SELECT id, tenant_id, name, description, status, industry_type,
		       employee_count, website, max_projects, max_users,
		       current_user_count, current_project_count,
		       billing_email, billing_address, created_at, updated_at
		FROM companies WHERE id = ?
	`

	company := &models.Company{}
	err := s.db.QueryRow(query, companyID).Scan(
		&company.ID, &company.TenantID, &company.Name, &company.Description,
		&company.Status, &company.IndustryType, &company.EmployeeCount,
		&company.Website, &company.MaxProjects, &company.MaxUsers,
		&company.CurrentUserCount, &company.CurrentProjectCount,
		&company.BillingEmail, &company.BillingAddress,
		&company.CreatedAt, &company.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("company not found: %s", companyID)
		}
		s.logger.Error("Failed to get company", "error", err, "company_id", companyID)
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	return company, nil
}

// ListCompaniesByTenant retrieves all companies for a tenant
func (s *CompanyService) ListCompaniesByTenant(tenantID string) ([]*models.Company, error) {
	query := `
		SELECT id, tenant_id, name, description, status, industry_type,
		       employee_count, website, max_projects, max_users,
		       current_user_count, current_project_count,
		       billing_email, billing_address, created_at, updated_at
		FROM companies
		WHERE tenant_id = ? AND status != 'archived'
		ORDER BY name
	`

	rows, err := s.db.Query(query, tenantID)
	if err != nil {
		s.logger.Error("Failed to list companies", "error", err)
		return nil, fmt.Errorf("failed to list companies: %w", err)
	}
	defer rows.Close()

	var companies []*models.Company
	for rows.Next() {
		company := &models.Company{}
		err := rows.Scan(
			&company.ID, &company.TenantID, &company.Name, &company.Description,
			&company.Status, &company.IndustryType, &company.EmployeeCount,
			&company.Website, &company.MaxProjects, &company.MaxUsers,
			&company.CurrentUserCount, &company.CurrentProjectCount,
			&company.BillingEmail, &company.BillingAddress,
			&company.CreatedAt, &company.UpdatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan company", "error", err)
			continue
		}
		companies = append(companies, company)
	}

	return companies, nil
}

// UpdateCompany updates company details
func (s *CompanyService) UpdateCompany(company *models.Company) error {
	query := `
		UPDATE companies SET
			name = ?, description = ?, status = ?,
			industry_type = ?, employee_count = ?, website = ?,
			max_projects = ?, max_users = ?,
			billing_email = ?, billing_address = ?, updated_at = NOW()
		WHERE id = ?
	`

	result, err := s.db.Exec(
		query,
		company.Name, company.Description, company.Status,
		company.IndustryType, company.EmployeeCount, company.Website,
		company.MaxProjects, company.MaxUsers,
		company.BillingEmail, company.BillingAddress, company.ID,
	)

	if err != nil {
		s.logger.Error("Failed to update company", "error", err)
		return fmt.Errorf("failed to update company: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("company not found: %s", company.ID)
	}

	return nil
}

// CreateProject creates a new project under a company
func (s *CompanyService) CreateProject(project *models.Project) error {
	query := `
		INSERT INTO projects (
			id, company_id, tenant_id, name, description, status,
			project_type, max_users, budget_allocated, start_date, end_date
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		project.ID, project.CompanyID, project.TenantID, project.Name,
		project.Description, project.Status, project.ProjectType,
		project.MaxUsers, project.BudgetAllocated, project.StartDate, project.EndDate,
	)

	if err != nil {
		s.logger.Error("Failed to create project", "error", err, "company_id", project.CompanyID)
		return fmt.Errorf("failed to create project: %w", err)
	}

	// Update company's project count
	s.updateCompanyProjectCount(project.CompanyID)

	s.logger.Info("Project created", "project_id", project.ID, "company_id", project.CompanyID)
	return nil
}

// GetProject retrieves a project by ID
func (s *CompanyService) GetProject(projectID string) (*models.Project, error) {
	query := `
		SELECT id, company_id, tenant_id, name, description, status,
		       project_type, max_users, current_user_count,
		       budget_allocated, budget_spent, start_date, end_date,
		       created_at, updated_at
		FROM projects WHERE id = ?
	`

	project := &models.Project{}
	err := s.db.QueryRow(query, projectID).Scan(
		&project.ID, &project.CompanyID, &project.TenantID, &project.Name,
		&project.Description, &project.Status, &project.ProjectType,
		&project.MaxUsers, &project.CurrentUserCount,
		&project.BudgetAllocated, &project.BudgetSpent, &project.StartDate,
		&project.EndDate, &project.CreatedAt, &project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found: %s", projectID)
		}
		s.logger.Error("Failed to get project", "error", err)
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

// ListProjectsByCompany retrieves all projects for a company
func (s *CompanyService) ListProjectsByCompany(companyID string) ([]*models.Project, error) {
	query := `
		SELECT id, company_id, tenant_id, name, description, status,
		       project_type, max_users, current_user_count,
		       budget_allocated, budget_spent, start_date, end_date,
		       created_at, updated_at
		FROM projects
		WHERE company_id = ? AND status != 'archived'
		ORDER BY name
	`

	rows, err := s.db.Query(query, companyID)
	if err != nil {
		s.logger.Error("Failed to list projects", "error", err)
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(
			&project.ID, &project.CompanyID, &project.TenantID, &project.Name,
			&project.Description, &project.Status, &project.ProjectType,
			&project.MaxUsers, &project.CurrentUserCount,
			&project.BudgetAllocated, &project.BudgetSpent, &project.StartDate,
			&project.EndDate, &project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan project", "error", err)
			continue
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// AddMemberToCompany adds a user to a company
func (s *CompanyService) AddMemberToCompany(member *models.CompanyMember) error {
	query := `
		INSERT INTO company_members (
			id, company_id, user_id, tenant_id, role, department, is_active
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		member.ID, member.CompanyID, member.UserID, member.TenantID,
		member.Role, member.Department, true,
	)

	if err != nil {
		s.logger.Error("Failed to add company member", "error", err)
		return fmt.Errorf("failed to add company member: %w", err)
	}

	// Update company user count
	s.updateCompanyUserCount(member.CompanyID)

	return nil
}

// AddMemberToProject adds a user to a project
func (s *CompanyService) AddMemberToProject(member *models.ProjectMember) error {
	query := `
		INSERT INTO project_members (
			id, project_id, user_id, company_id, tenant_id, role, is_active
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		member.ID, member.ProjectID, member.UserID, member.CompanyID,
		member.TenantID, member.Role, true,
	)

	if err != nil {
		s.logger.Error("Failed to add project member", "error", err)
		return fmt.Errorf("failed to add project member: %w", err)
	}

	// Update project user count
	s.updateProjectUserCount(member.ProjectID)

	return nil
}

// GetCompanyMembers retrieves all members of a company
func (s *CompanyService) GetCompanyMembers(companyID string) ([]*models.CompanyMember, error) {
	query := `
		SELECT id, company_id, user_id, tenant_id, role, department,
		       is_active, created_at, updated_at
		FROM company_members
		WHERE company_id = ? AND is_active = TRUE
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, companyID)
	if err != nil {
		s.logger.Error("Failed to get company members", "error", err)
		return nil, fmt.Errorf("failed to get company members: %w", err)
	}
	defer rows.Close()

	var members []*models.CompanyMember
	for rows.Next() {
		member := &models.CompanyMember{}
		err := rows.Scan(
			&member.ID, &member.CompanyID, &member.UserID, &member.TenantID,
			&member.Role, &member.Department, &member.IsActive,
			&member.CreatedAt, &member.UpdatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan member", "error", err)
			continue
		}
		members = append(members, member)
	}

	return members, nil
}

// GetProjectMembers retrieves all members of a project
func (s *CompanyService) GetProjectMembers(projectID string) ([]*models.ProjectMember, error) {
	query := `
		SELECT id, project_id, user_id, company_id, tenant_id, role,
		       joined_at, is_active, created_at, updated_at
		FROM project_members
		WHERE project_id = ? AND is_active = TRUE
		ORDER BY joined_at DESC
	`

	rows, err := s.db.Query(query, projectID)
	if err != nil {
		s.logger.Error("Failed to get project members", "error", err)
		return nil, fmt.Errorf("failed to get project members: %w", err)
	}
	defer rows.Close()

	var members []*models.ProjectMember
	for rows.Next() {
		member := &models.ProjectMember{}
		err := rows.Scan(
			&member.ID, &member.ProjectID, &member.UserID, &member.CompanyID,
			&member.TenantID, &member.Role, &member.JoinedAt, &member.IsActive,
			&member.CreatedAt, &member.UpdatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan project member", "error", err)
			continue
		}
		members = append(members, member)
	}

	return members, nil
}

// RemoveMemberFromProject removes a user from a project
func (s *CompanyService) RemoveMemberFromProject(projectID string, userID int) error {
	query := `UPDATE project_members SET is_active = FALSE WHERE project_id = ? AND user_id = ?`

	result, err := s.db.Exec(query, projectID, userID)
	if err != nil {
		s.logger.Error("Failed to remove project member", "error", err)
		return fmt.Errorf("failed to remove project member: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("project member not found")
	}

	// Update project user count
	s.updateProjectUserCount(projectID)

	return nil
}

// Helper functions
func (s *CompanyService) updateCompanyUserCount(companyID string) {
	query := `
		UPDATE companies SET current_user_count = (
			SELECT COUNT(DISTINCT user_id) FROM company_members
			WHERE company_id = ? AND is_active = TRUE
		) WHERE id = ?
	`
	s.db.Exec(query, companyID, companyID)
}

func (s *CompanyService) updateCompanyProjectCount(companyID string) {
	query := `
		UPDATE companies SET current_project_count = (
			SELECT COUNT(*) FROM projects
			WHERE company_id = ? AND status != 'archived'
		) WHERE id = ?
	`
	s.db.Exec(query, companyID, companyID)
}

func (s *CompanyService) updateProjectUserCount(projectID string) {
	query := `
		UPDATE projects SET current_user_count = (
			SELECT COUNT(DISTINCT user_id) FROM project_members
			WHERE project_id = ? AND is_active = TRUE
		) WHERE id = ?
	`
	s.db.Exec(query, projectID, projectID)
}
