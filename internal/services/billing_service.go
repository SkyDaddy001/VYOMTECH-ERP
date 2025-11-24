package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/pkg/logger"
)

type BillingService struct {
	db            *sql.DB
	logger        *logger.Logger
	moduleService *ModuleService
}

func NewBillingService(db *sql.DB, logger *logger.Logger, moduleService *ModuleService) *BillingService {
	return &BillingService{
		db:            db,
		logger:        logger,
		moduleService: moduleService,
	}
}

// CreatePricingPlan creates a new pricing plan
func (s *BillingService) CreatePricingPlan(plan *models.PricingPlan) error {
	query := `
		INSERT INTO pricing_plans (
			id, name, description, monthly_price, annual_price,
			max_users, max_companies, max_projects_per_company,
			included_modules, additional_modules, features,
			sort_order, is_active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		plan.ID, plan.Name, plan.Description, plan.MonthlyPrice, plan.AnnualPrice,
		plan.MaxUsers, plan.MaxCompanies, plan.MaxProjectsPerCompany,
		plan.IncludedModules, plan.AdditionalModules, plan.Features,
		plan.SortOrder, plan.IsActive,
	)

	if err != nil {
		s.logger.Error("Failed to create pricing plan", "error", err)
		return fmt.Errorf("failed to create pricing plan: %w", err)
	}

	s.logger.Info("Pricing plan created", "plan_id", plan.ID)
	return nil
}

// GetPricingPlan retrieves a pricing plan
func (s *BillingService) GetPricingPlan(planID string) (*models.PricingPlan, error) {
	query := `
		SELECT id, name, description, monthly_price, annual_price,
		       max_users, max_companies, max_projects_per_company,
		       included_modules, additional_modules, features,
		       sort_order, is_active, created_at, updated_at
		FROM pricing_plans WHERE id = ?
	`

	plan := &models.PricingPlan{}
	err := s.db.QueryRow(query, planID).Scan(
		&plan.ID, &plan.Name, &plan.Description, &plan.MonthlyPrice, &plan.AnnualPrice,
		&plan.MaxUsers, &plan.MaxCompanies, &plan.MaxProjectsPerCompany,
		&plan.IncludedModules, &plan.AdditionalModules, &plan.Features,
		&plan.SortOrder, &plan.IsActive, &plan.CreatedAt, &plan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pricing plan not found")
		}
		return nil, fmt.Errorf("failed to get pricing plan: %w", err)
	}

	return plan, nil
}

// ListActivePricingPlans retrieves all active pricing plans
func (s *BillingService) ListActivePricingPlans() ([]*models.PricingPlan, error) {
	query := `
		SELECT id, name, description, monthly_price, annual_price,
		       max_users, max_companies, max_projects_per_company,
		       included_modules, additional_modules, features,
		       sort_order, is_active, created_at, updated_at
		FROM pricing_plans
		WHERE is_active = TRUE
		ORDER BY sort_order, name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list pricing plans: %w", err)
	}
	defer rows.Close()

	var plans []*models.PricingPlan
	for rows.Next() {
		plan := &models.PricingPlan{}
		err := rows.Scan(
			&plan.ID, &plan.Name, &plan.Description, &plan.MonthlyPrice, &plan.AnnualPrice,
			&plan.MaxUsers, &plan.MaxCompanies, &plan.MaxProjectsPerCompany,
			&plan.IncludedModules, &plan.AdditionalModules, &plan.Features,
			&plan.SortOrder, &plan.IsActive, &plan.CreatedAt, &plan.UpdatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan pricing plan", "error", err)
			continue
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

// SubscribeToPlan subscribes a tenant to a pricing plan
func (s *BillingService) SubscribeToPlan(tenantID, planID string, billingCycle models.BillingCycleType) (*models.TenantPlanSubscription, error) {
	plan, err := s.GetPricingPlan(planID)
	if err != nil {
		return nil, err
	}

	subscription := &models.TenantPlanSubscription{
		ID:            fmt.Sprintf("sub_%d", time.Now().UnixNano()),
		TenantID:      tenantID,
		PricingPlanID: planID,
		StartDate:     time.Now(),
		BillingCycle:  billingCycle,
		IsAutoRenew:   true,
		Status:        "active",
	}

	// Calculate next billing date
	switch billingCycle {
	case models.BillingCycleMonthly:
		subscription.NextBillingDate = subscription.StartDate.AddDate(0, 1, 0)
	case models.BillingCycleQuarterly:
		subscription.NextBillingDate = subscription.StartDate.AddDate(0, 3, 0)
	case models.BillingCycleAnnual:
		subscription.NextBillingDate = subscription.StartDate.AddDate(1, 0, 0)
	}

	query := `
		INSERT INTO tenant_plan_subscriptions (
			id, tenant_id, pricing_plan_id, start_date, next_billing_date,
			status, billing_cycle, is_auto_renew
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = s.db.Exec(
		query,
		subscription.ID, subscription.TenantID, subscription.PricingPlanID,
		subscription.StartDate, subscription.NextBillingDate,
		subscription.Status, subscription.BillingCycle, subscription.IsAutoRenew,
	)

	if err != nil {
		s.logger.Error("Failed to subscribe to plan", "error", err)
		return nil, fmt.Errorf("failed to subscribe to plan: %w", err)
	}

	// Subscribe to all included modules
	var includedModules []string
	json.Unmarshal([]byte(plan.IncludedModules), &includedModules)

	for _, moduleID := range includedModules {
		subscription := &models.ModuleSubscription{
			ID:       fmt.Sprintf("modsub_%d_%s", time.Now().UnixNano(), moduleID),
			TenantID: tenantID,
			ModuleID: moduleID,
			Status:   "active",
		}
		s.moduleService.SubscribeToModule(subscription)
	}

	s.logger.Info("Tenant subscribed to plan", "tenant_id", tenantID, "plan_id", planID)
	return subscription, nil
}

// CreateInvoice generates an invoice for a tenant
func (s *BillingService) CreateInvoice(invoice *models.Invoice) error {
	query := `
		INSERT INTO invoices (
			id, tenant_id, invoice_number, billing_period_start, billing_period_end,
			subtotal_amount, tax_amount, discount_amount, total_amount,
			status, payment_method, issued_at, due_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		invoice.ID, invoice.TenantID, invoice.InvoiceNumber,
		invoice.BillingPeriodStart, invoice.BillingPeriodEnd,
		invoice.SubtotalAmount, invoice.TaxAmount, invoice.DiscountAmount,
		invoice.TotalAmount, invoice.Status, invoice.PaymentMethod,
		invoice.IssuedAt, invoice.DueAt,
	)

	if err != nil {
		s.logger.Error("Failed to create invoice", "error", err)
		return fmt.Errorf("failed to create invoice: %w", err)
	}

	s.logger.Info("Invoice created", "invoice_id", invoice.ID, "invoice_number", invoice.InvoiceNumber)
	return nil
}

// AddLineItem adds a line item to an invoice
func (s *BillingService) AddLineItem(lineItem *models.InvoiceLineItem) error {
	query := `
		INSERT INTO invoice_line_items (
			id, invoice_id, module_id, description, quantity,
			unit_price, total_price, tax_rate
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		lineItem.ID, lineItem.InvoiceID, lineItem.ModuleID, lineItem.Description,
		lineItem.Quantity, lineItem.UnitPrice, lineItem.TotalPrice, lineItem.TaxRate,
	)

	if err != nil {
		s.logger.Error("Failed to add line item", "error", err)
		return fmt.Errorf("failed to add line item: %w", err)
	}

	return nil
}

// GetInvoice retrieves an invoice
func (s *BillingService) GetInvoice(invoiceID string) (*models.Invoice, error) {
	query := `
		SELECT id, tenant_id, invoice_number, billing_period_start, billing_period_end,
		       subtotal_amount, tax_amount, discount_amount, total_amount,
		       status, payment_method, issued_at, due_at, paid_at,
		       created_at, updated_at
		FROM invoices WHERE id = ?
	`

	invoice := &models.Invoice{}
	err := s.db.QueryRow(query, invoiceID).Scan(
		&invoice.ID, &invoice.TenantID, &invoice.InvoiceNumber,
		&invoice.BillingPeriodStart, &invoice.BillingPeriodEnd,
		&invoice.SubtotalAmount, &invoice.TaxAmount, &invoice.DiscountAmount,
		&invoice.TotalAmount, &invoice.Status, &invoice.PaymentMethod,
		&invoice.IssuedAt, &invoice.DueAt, &invoice.PaidAt,
		&invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	return invoice, nil
}

// ListInvoicesByTenant retrieves invoices for a tenant
func (s *BillingService) ListInvoicesByTenant(tenantID string, limit, offset int) ([]*models.Invoice, error) {
	query := `
		SELECT id, tenant_id, invoice_number, billing_period_start, billing_period_end,
		       subtotal_amount, tax_amount, discount_amount, total_amount,
		       status, payment_method, issued_at, due_at, paid_at,
		       created_at, updated_at
		FROM invoices
		WHERE tenant_id = ?
		ORDER BY issued_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*models.Invoice
	for rows.Next() {
		invoice := &models.Invoice{}
		err := rows.Scan(
			&invoice.ID, &invoice.TenantID, &invoice.InvoiceNumber,
			&invoice.BillingPeriodStart, &invoice.BillingPeriodEnd,
			&invoice.SubtotalAmount, &invoice.TaxAmount, &invoice.DiscountAmount,
			&invoice.TotalAmount, &invoice.Status, &invoice.PaymentMethod,
			&invoice.IssuedAt, &invoice.DueAt, &invoice.PaidAt,
			&invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan invoice", "error", err)
			continue
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

// MarkInvoiceAsPaid marks an invoice as paid
func (s *BillingService) MarkInvoiceAsPaid(invoiceID string, paidAt time.Time) error {
	query := `UPDATE invoices SET status = 'paid', paid_at = ? WHERE id = ?`

	result, err := s.db.Exec(query, paidAt, invoiceID)
	if err != nil {
		s.logger.Error("Failed to mark invoice as paid", "error", err)
		return fmt.Errorf("failed to mark invoice as paid: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("invoice not found")
	}

	return nil
}

// RecordUsageMetrics records usage metrics for a tenant
func (s *BillingService) RecordUsageMetrics(metrics *models.UsageMetrics) error {
	query := `
		INSERT INTO usage_metrics (
			id, tenant_id, company_id, project_id, date,
			active_users, new_users, api_calls_used, storage_used_mb,
			module_usage_data
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			active_users = VALUES(active_users),
			new_users = VALUES(new_users),
			api_calls_used = VALUES(api_calls_used),
			storage_used_mb = VALUES(storage_used_mb),
			module_usage_data = VALUES(module_usage_data)
	`

	_, err := s.db.Exec(
		query,
		metrics.ID, metrics.TenantID, metrics.CompanyID, metrics.ProjectID,
		metrics.Date, metrics.ActiveUsers, metrics.NewUsers, metrics.APICallsUsed,
		metrics.StorageUsedMB, metrics.ModuleUsageData,
	)

	if err != nil {
		s.logger.Error("Failed to record usage metrics", "error", err)
		return fmt.Errorf("failed to record usage metrics: %w", err)
	}

	return nil
}

// GetUsageMetrics retrieves usage metrics
func (s *BillingService) GetUsageMetrics(tenantID string, startDate, endDate time.Time) ([]*models.UsageMetrics, error) {
	query := `
		SELECT id, tenant_id, company_id, project_id, date,
		       active_users, new_users, api_calls_used, storage_used_mb,
		       module_usage_data, created_at
		FROM usage_metrics
		WHERE tenant_id = ? AND date BETWEEN ? AND ?
		ORDER BY date DESC
	`

	rows, err := s.db.Query(query, tenantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage metrics: %w", err)
	}
	defer rows.Close()

	var metrics []*models.UsageMetrics
	for rows.Next() {
		metric := &models.UsageMetrics{}
		err := rows.Scan(
			&metric.ID, &metric.TenantID, &metric.CompanyID, &metric.ProjectID,
			&metric.Date, &metric.ActiveUsers, &metric.NewUsers, &metric.APICallsUsed,
			&metric.StorageUsedMB, &metric.ModuleUsageData, &metric.CreatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan usage metric", "error", err)
			continue
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// CalculateMonthlyCharges calculates charges for a billing period
func (s *BillingService) CalculateMonthlyCharges(tenantID string) (float64, error) {
	// Get all active subscriptions
	query := `
		SELECT SUM(amount_spent_this_month)
		FROM module_subscriptions
		WHERE tenant_id = ? AND is_enabled = TRUE
	`

	var totalCost float64
	err := s.db.QueryRow(query, tenantID).Scan(&totalCost)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to calculate monthly charges: %w", err)
	}

	return totalCost, nil
}
