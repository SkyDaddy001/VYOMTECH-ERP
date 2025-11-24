package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/pkg/logger"
)

type ModuleService struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewModuleService(db *sql.DB, logger *logger.Logger) *ModuleService {
	return &ModuleService{
		db:     db,
		logger: logger,
	}
}

// RegisterModule creates a new module in the system
func (s *ModuleService) RegisterModule(module *models.Module) error {
	query := `
		INSERT INTO modules (
			id, name, description, category, status, version, pricing_model,
			base_cost, cost_per_user, cost_per_project, cost_per_company,
			max_users, max_projects, max_companies, is_dependent_on,
			is_core, requires_approval, trial_days_allowed, features
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		module.ID, module.Name, module.Description, module.Category,
		module.Status, module.Version, module.PricingModel,
		module.BaseCost, module.CostPerUser, module.CostPerProject, module.CostPerCompany,
		module.MaxUsers, module.MaxProjects, module.MaxCompanies,
		module.IsDependentOn, module.IsCore, module.RequiresApproval,
		module.TrialDaysAllowed, module.Features,
	)

	if err != nil {
		s.logger.Error("Failed to register module", "error", err, "module_id", module.ID)
		return fmt.Errorf("failed to register module: %w", err)
	}

	s.logger.Info("Module registered successfully", "module_id", module.ID, "name", module.Name)
	return nil
}

// GetModule retrieves a module by ID
func (s *ModuleService) GetModule(moduleID string) (*models.Module, error) {
	query := `
		SELECT id, name, description, category, status, version, pricing_model,
		       base_cost, cost_per_user, cost_per_project, cost_per_company,
		       max_users, max_projects, max_companies, is_dependent_on,
		       is_core, requires_approval, trial_days_allowed, features,
		       created_at, updated_at
		FROM modules WHERE id = ?
	`

	module := &models.Module{}
	err := s.db.QueryRow(query, moduleID).Scan(
		&module.ID, &module.Name, &module.Description, &module.Category,
		&module.Status, &module.Version, &module.PricingModel,
		&module.BaseCost, &module.CostPerUser, &module.CostPerProject, &module.CostPerCompany,
		&module.MaxUsers, &module.MaxProjects, &module.MaxCompanies,
		&module.IsDependentOn, &module.IsCore, &module.RequiresApproval,
		&module.TrialDaysAllowed, &module.Features,
		&module.CreatedAt, &module.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("module not found: %s", moduleID)
		}
		s.logger.Error("Failed to get module", "error", err, "module_id", moduleID)
		return nil, fmt.Errorf("failed to get module: %w", err)
	}

	return module, nil
}

// ListModules returns all modules, optionally filtered by status
func (s *ModuleService) ListModules(status *string) ([]*models.Module, error) {
	query := `
		SELECT id, name, description, category, status, version, pricing_model,
		       base_cost, cost_per_user, cost_per_project, cost_per_company,
		       max_users, max_projects, max_companies, is_dependent_on,
		       is_core, requires_approval, trial_days_allowed, features,
		       created_at, updated_at
		FROM modules
	`

	if status != nil {
		query += " WHERE status = ?"
	}

	query += " ORDER BY category, name"

	var rows *sql.Rows
	var err error

	if status != nil {
		rows, err = s.db.Query(query, *status)
	} else {
		rows, err = s.db.Query(query)
	}

	if err != nil {
		s.logger.Error("Failed to list modules", "error", err)
		return nil, fmt.Errorf("failed to list modules: %w", err)
	}
	defer rows.Close()

	var modules []*models.Module
	for rows.Next() {
		module := &models.Module{}
		err := rows.Scan(
			&module.ID, &module.Name, &module.Description, &module.Category,
			&module.Status, &module.Version, &module.PricingModel,
			&module.BaseCost, &module.CostPerUser, &module.CostPerProject, &module.CostPerCompany,
			&module.MaxUsers, &module.MaxProjects, &module.MaxCompanies,
			&module.IsDependentOn, &module.IsCore, &module.RequiresApproval,
			&module.TrialDaysAllowed, &module.Features,
			&module.CreatedAt, &module.UpdatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan module", "error", err)
			continue
		}
		modules = append(modules, module)
	}

	return modules, nil
}

// SubscribeToModule subscribes a tenant/company/project to a module
func (s *ModuleService) SubscribeToModule(subscription *models.ModuleSubscription) error {
	// Check if module exists
	module, err := s.GetModule(subscription.ModuleID)
	if err != nil {
		return fmt.Errorf("module does not exist: %w", err)
	}

	// Check for existing subscription
	existingQuery := `
		SELECT id FROM module_subscriptions
		WHERE tenant_id = ? AND module_id = ? 
		  AND (company_id <=> ?) AND (project_id <=> ?)
	`
	var existingID string
	err = s.db.QueryRow(
		existingQuery,
		subscription.TenantID,
		subscription.ModuleID,
		subscription.CompanyID,
		subscription.ProjectID,
	).Scan(&existingID)

	if err == nil {
		return fmt.Errorf("subscription already exists")
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing subscription: %w", err)
	}

	// Create subscription
	query := `
		INSERT INTO module_subscriptions (
			id, tenant_id, company_id, project_id, module_id, status,
			subscription_started, trial_started_at, trial_ends_at,
			max_users_allowed, monthly_budget, configuration, is_enabled
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	trialStarted := time.Now()
	trialEnds := trialStarted.AddDate(0, 0, module.TrialDaysAllowed)

	_, err = s.db.Exec(
		query,
		subscription.ID, subscription.TenantID, subscription.CompanyID, subscription.ProjectID,
		subscription.ModuleID, "trial",
		time.Now(), trialStarted, trialEnds,
		subscription.MaxUsersAllowed, subscription.MonthlyBudget,
		subscription.Configuration, true,
	)

	if err != nil {
		s.logger.Error("Failed to subscribe to module", "error", err, "module_id", subscription.ModuleID)
		return fmt.Errorf("failed to subscribe to module: %w", err)
	}

	s.logger.Info("Subscribed to module", "subscription_id", subscription.ID, "module_id", subscription.ModuleID)
	return nil
}

// GetSubscription retrieves a module subscription
func (s *ModuleService) GetSubscription(subscriptionID string) (*models.ModuleSubscription, error) {
	query := `
		SELECT id, tenant_id, company_id, project_id, module_id, status,
		       subscription_started, subscription_ended, trial_started_at, trial_ends_at,
		       max_users_allowed, current_user_count, monthly_budget, amount_spent_this_month,
		       configuration, is_enabled, created_at, updated_at
		FROM module_subscriptions WHERE id = ?
	`

	sub := &models.ModuleSubscription{}
	err := s.db.QueryRow(query, subscriptionID).Scan(
		&sub.ID, &sub.TenantID, &sub.CompanyID, &sub.ProjectID, &sub.ModuleID,
		&sub.Status, &sub.SubscriptionStarted, &sub.SubscriptionEnded,
		&sub.TrialStartedAt, &sub.TrialEndsAt,
		&sub.MaxUsersAllowed, &sub.CurrentUserCount, &sub.MonthlyBudget,
		&sub.AmountSpentThisMonth, &sub.Configuration, &sub.IsEnabled,
		&sub.CreatedAt, &sub.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription not found: %s", subscriptionID)
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return sub, nil
}

// ListSubscriptions retrieves all subscriptions for a tenant
func (s *ModuleService) ListSubscriptions(tenantID string, companyID *string, projectID *string) ([]*models.ModuleSubscription, error) {
	query := `
		SELECT id, tenant_id, company_id, project_id, module_id, status,
		       subscription_started, subscription_ended, trial_started_at, trial_ends_at,
		       max_users_allowed, current_user_count, monthly_budget, amount_spent_this_month,
		       configuration, is_enabled, created_at, updated_at
		FROM module_subscriptions
		WHERE tenant_id = ?
	`

	args := []interface{}{tenantID}

	if companyID != nil {
		query += " AND company_id = ?"
		args = append(args, *companyID)
	} else {
		query += " AND company_id IS NULL"
	}

	if projectID != nil {
		query += " AND project_id = ?"
		args = append(args, *projectID)
	} else {
		query += " AND project_id IS NULL"
	}

	query += " ORDER BY module_id"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		s.logger.Error("Failed to list subscriptions", "error", err)
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.ModuleSubscription
	for rows.Next() {
		sub := &models.ModuleSubscription{}
		err := rows.Scan(
			&sub.ID, &sub.TenantID, &sub.CompanyID, &sub.ProjectID, &sub.ModuleID,
			&sub.Status, &sub.SubscriptionStarted, &sub.SubscriptionEnded,
			&sub.TrialStartedAt, &sub.TrialEndsAt,
			&sub.MaxUsersAllowed, &sub.CurrentUserCount, &sub.MonthlyBudget,
			&sub.AmountSpentThisMonth, &sub.Configuration, &sub.IsEnabled,
			&sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan subscription", "error", err)
			continue
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

// ToggleModule enables or disables a subscription
func (s *ModuleService) ToggleModule(subscriptionID string, enabled bool) error {
	query := `UPDATE module_subscriptions SET is_enabled = ? WHERE id = ?`

	result, err := s.db.Exec(query, enabled, subscriptionID)
	if err != nil {
		s.logger.Error("Failed to toggle module", "error", err, "subscription_id", subscriptionID)
		return fmt.Errorf("failed to toggle module: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found: %s", subscriptionID)
	}

	s.logger.Info("Module toggled", "subscription_id", subscriptionID, "enabled", enabled)
	return nil
}

// GetModuleUsage retrieves usage metrics for a module
func (s *ModuleService) GetModuleUsage(subscriptionID string, startDate, endDate time.Time) ([]*models.ModuleUsage, error) {
	query := `
		SELECT id, subscription_id, tenant_id, company_id, project_id, module_id,
		       user_count, project_count, company_count, custom_metrics,
		       usage_date, estimated_cost, created_at
		FROM module_usage
		WHERE subscription_id = ? AND usage_date BETWEEN ? AND ?
		ORDER BY usage_date DESC
	`

	rows, err := s.db.Query(query, subscriptionID, startDate, endDate)
	if err != nil {
		s.logger.Error("Failed to get module usage", "error", err)
		return nil, fmt.Errorf("failed to get module usage: %w", err)
	}
	defer rows.Close()

	var usages []*models.ModuleUsage
	for rows.Next() {
		usage := &models.ModuleUsage{}
		err := rows.Scan(
			&usage.ID, &usage.SubscriptionID, &usage.TenantID, &usage.CompanyID, &usage.ProjectID,
			&usage.ModuleID, &usage.UserCount, &usage.ProjectCount, &usage.CompanyCount,
			&usage.CustomMetrics, &usage.UsageDate, &usage.EstimatedCost, &usage.CreatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan usage", "error", err)
			continue
		}
		usages = append(usages, usage)
	}

	return usages, nil
}

// RecordUsage records usage metrics for a subscription
func (s *ModuleService) RecordUsage(usage *models.ModuleUsage) error {
	query := `
		INSERT INTO module_usage (
			id, subscription_id, tenant_id, company_id, project_id, module_id,
			user_count, project_count, company_count, custom_metrics,
			usage_date, estimated_cost
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			user_count = VALUES(user_count),
			project_count = VALUES(project_count),
			company_count = VALUES(company_count),
			custom_metrics = VALUES(custom_metrics),
			estimated_cost = VALUES(estimated_cost)
	`

	_, err := s.db.Exec(
		query,
		usage.ID, usage.SubscriptionID, usage.TenantID, usage.CompanyID, usage.ProjectID,
		usage.ModuleID, usage.UserCount, usage.ProjectCount, usage.CompanyCount,
		usage.CustomMetrics, usage.UsageDate, usage.EstimatedCost,
	)

	if err != nil {
		s.logger.Error("Failed to record usage", "error", err)
		return fmt.Errorf("failed to record usage: %w", err)
	}

	return nil
}

// CalculateModuleCost calculates the cost for a module subscription
func (s *ModuleService) CalculateModuleCost(module *models.Module, userCount, projectCount, companyCount int) float64 {
	var cost float64 = module.BaseCost

	switch module.PricingModel {
	case models.PricingModelPerUser:
		cost += float64(userCount) * module.CostPerUser
	case models.PricingModelPerProject:
		cost += float64(projectCount) * module.CostPerProject
	case models.PricingModelPerCompany:
		cost += float64(companyCount) * module.CostPerCompany
	case models.PricingModelFlat:
		// Just base cost
	case models.PricingModelTiered:
		// Implement tiered pricing logic based on usage
		if userCount > 100 {
			cost = module.BaseCost * 1.5
		} else if userCount > 50 {
			cost = module.BaseCost * 1.25
		}
	}

	return cost
}

// CheckModuleDependencies verifies all dependencies are enabled
func (s *ModuleService) CheckModuleDependencies(tenantID string, moduleID string) (bool, []string, error) {
	module, err := s.GetModule(moduleID)
	if err != nil {
		return false, nil, err
	}

	if module.IsDependentOn == nil {
		return true, []string{}, nil
	}

	var dependencies []string
	err = json.Unmarshal(module.IsDependentOn, &dependencies)
	if err != nil {
		s.logger.Error("Failed to parse dependencies", "error", err, "module_id", moduleID)
		return false, nil, err
	}

	var missing []string
	for _, depID := range dependencies {
		query := `SELECT id FROM module_subscriptions WHERE tenant_id = ? AND module_id = ? AND is_enabled = TRUE`
		var id string
		err := s.db.QueryRow(query, tenantID, depID).Scan(&id)
		if err == sql.ErrNoRows {
			missing = append(missing, depID)
		}
	}

	return len(missing) == 0, missing, nil
}
